import { sendSubscribeRequest } from "@/adapter/backend/subscribe.js";
import { config } from "@/config.js";
import { imageToTalkHistories } from "@/functions/image_to_talk_histories.js";
import {
  Account_PlatformType,
  BackendService,
  type TalkHistory,
  TalkRequestSchema,
  User_UserRole,
} from "@/gen/grpc/backend/v1/backend_pb.js";
import { create } from "@bufbuild/protobuf";
import { createClient } from "@connectrpc/connect";
import { createGrpcTransport } from "@connectrpc/connect-node";
import { messagingApi } from "@line/bot-sdk";
import type { Request, Response } from "express";
import express from "express";

const app = express();
app.use(express.json());

const { MessagingApiClient } = messagingApi;

const client = new MessagingApiClient(config.line.messagingApiClient);

const transport = createGrpcTransport({
  baseUrl: config.backend.url,
});
export const BackendClient = createClient(BackendService, transport);

async function sendTalkRequest(
  talkHistories: TalkHistory[],
  accountId: string,
): Promise<{ messages: string[]; status: string } | null> {
  try {
    const request = create(TalkRequestSchema, {
      histories: talkHistories,
      actionKind: 1,
      account: {
        platformType: Account_PlatformType.PLATFORM_LINE,
        accountId: accountId,
      },
    });

    const response = await BackendClient.talk(request);
    let status = "error";
    switch (response.status) {
      case 1:
        status = "success";
        break;
      case 2:
        status = "error";
        break;
      case 3:
        status = "limit";
        break;
      default:
        status = "error";
        break;
    }
    return {
      messages: response.message,
      status: status,
    };
  } catch (error) {
    console.error("Error:", error);
    return null;
  }
}

export function getDate(row: string): string | undefined {
  let talkDate: string | undefined = undefined;
  const newDateMatch = row.match(
    /^(\d{4})\/(\d{1,2})\/(\d{1,2})|^(Sun|Mon|Tue|Wed|Thu|Fri|Sat), (\d{1,2})\/(\d{1,2})\/(\d{4})/,
  );

  if (newDateMatch) {
    if (newDateMatch[1] && newDateMatch[2] && newDateMatch[3]) {
      // フォーマット: 2025/1/12(日)
      const year = newDateMatch[1];
      const month = newDateMatch[2].padStart(2, "0");
      const day = newDateMatch[3].padStart(2, "0");
      talkDate = `${year}-${month}-${day}`;
    } else if (
      newDateMatch[4] &&
      newDateMatch[5] &&
      newDateMatch[6] &&
      newDateMatch[7]
    ) {
      // フォーマット: Sun, 1/12/2025
      const month = newDateMatch[5].padStart(2, "0");
      const day = newDateMatch[6].padStart(2, "0");
      const year = newDateMatch[7];
      talkDate = `${year}-${month}-${day}`;
    } else {
      talkDate = undefined;
    }
    return talkDate;
  }
}

export function parseTalkHistories(
  talk: string,
  selfName: string,
): TalkHistory[] {
  const rows = talk.split("\n");
  const TalkHistories: TalkHistory[] = [];
  let talkDate: string | undefined = undefined;

  for (const row of rows) {
    const trimmedRow = row.trim();
    // 日付
    const newTalkDate = getDate(trimmedRow);
    if (newTalkDate) {
      talkDate = newTalkDate;
    }
    if (talkDate === undefined) {
      continue;
    }
    // メッセージ（例: "22:07   Test    おはよう"）
    const messageMatch = trimmedRow.match(
      /^(\d{1,2}):(\d{2})\t+([^\t]+)?\t+(.+)$/,
    );
    if (messageMatch && talkDate) {
      const [_, hour, minutes, userName, message] = messageMatch;
      const time = `${hour.padStart(2, "0")}:${minutes}`;
      const dateTime = `${talkDate}T${time}:00+09:00`; // ISO 8601形式

      const name = userName || "Unknown";
      const userId =
        name === selfName
          ? `${selfName}01`
          : name === "Unknown"
            ? "Unknown"
            : `${name}02`;
      const role = (() => {
        switch (name) {
          case selfName:
            return User_UserRole.SELF;
          case "Unknown":
            return User_UserRole.UNSPECIFIED;
          default:
            return User_UserRole.YOU;
        }
      })();

      TalkHistories.push({
        $typeName: "backend.v1.TalkHistory",
        date: dateTime,
        user: {
          $typeName: "backend.v1.User",
          name,
          userId,
          role: role,
        },
        message,
      });
    }
  }

  return TalkHistories;
}

/**
 * チャットにローディングアニメーションを表示する
 *
 * @param chatId
 */
async function loading_animation(chatId: string): Promise<void> {
  const endpoint = "https://api.line.me/v2/bot/chat/loading/start";

  const response = await fetch(endpoint, {
    method: "POST",
    headers: {
      Authorization: `Bearer ${config.line.messagingApiClient.channelAccessToken}`,
      "Content-Type": "application/json",
    },
    body: JSON.stringify({
      chatId: chatId,
      loadingSeconds: 60,
    }),
  });
  if (!response.ok) {
    throw new Error(`Failed to fetch content: ${response.statusText}`);
  }
}

export const webhookHandler = async (
  req: Request,
  res: Response,
): Promise<void> => {
  try {
    if (req.body.events && req.body.events.length > 0) {
      // biome-ignore lint/suspicious/noExplicitAny: <explanation>
      const eventPromises = req.body.events.map(async (e: any) => {
        console.log("invoking");
        if (e.source?.userId == null) {
          await client.replyMessage({
            replyToken: e.replyToken,
            messages: [
              {
                type: "text",
                text: "エラーが発生しました。もう一度やり直してください",
              },
            ],
          });
          console.error("ユーザーIDが取得できませんでした。");
          return;
        }
        // LINEユーザーのプロフィールを取得
        const profile = await client.getProfile(e.source.userId);
        const selfName = profile.displayName;
        const userId = profile.userId;

        let talkHistories: TalkHistory[] | null = null;

        if (e.type === "message" && e.message.type === "text") {
          loading_animation(userId);
          switch (e.message.text) {
            case "subscribe":
              {
                const res = await sendSubscribeRequest(userId, 1);
                if (res?.status === "ACTIVE") {
                  const expired_date = new Date(
                    res.expire_at * 1000,
                  ).toLocaleDateString();
                  await client.replyMessage({
                    replyToken: e.replyToken,
                    messages: [
                      {
                        type: "text",
                        text: `サブスク中です。\n有効期限: ${expired_date}`,
                      },
                    ],
                  });
                } else {
                  await client.replyMessage({
                    replyToken: e.replyToken,
                    messages: [
                      {
                        type: "text",
                        text: `こちらからサブスクリプションを購入できます。\n${res?.redirect_url}`,
                      },
                    ],
                  });
                }
              }
              return;

            case "unsubscribe":
              {
                const res = await sendSubscribeRequest(userId, 2);
                if (res === null) {
                  await client.replyMessage({
                    replyToken: e.replyToken,
                    messages: [
                      {
                        type: "text",
                        text: "エラーが発生しました。もう一度やり直してください",
                      },
                    ],
                  });
                  return;
                }

                if (res?.status === "UNSPECIFIED") {
                  await client.replyMessage({
                    replyToken: e.replyToken,
                    messages: [
                      {
                        type: "text",
                        text: "サブスクリプションはありません。",
                      },
                    ],
                  });
                }
                const expired_date = new Date(
                  res.expire_at * 1000,
                ).toLocaleDateString();
                await client.replyMessage({
                  replyToken: e.replyToken,
                  messages: [
                    {
                      type: "text",
                      text: `状態: ${res?.messages}\n有効期限: ${expired_date}`,
                    },
                  ],
                });
              }
              return;
            case "check":
              {
                const res = await sendSubscribeRequest(userId, 3);
                if (res?.status === "ACTIVE") {
                  const expired_date = new Date(
                    res.expire_at * 1000,
                  ).toLocaleDateString();
                  await client.replyMessage({
                    replyToken: e.replyToken,
                    messages: [
                      {
                        type: "text",
                        text: `サブスク中です。\n有効期限: ${expired_date}`,
                      },
                    ],
                  });
                } else {
                  await client.replyMessage({
                    replyToken: e.replyToken,
                    messages: [
                      {
                        type: "text",
                        text: "サブスク中ではありません。",
                      },
                    ],
                  });
                }
              }
              return;
            default:
              await client.replyMessage({
                replyToken: e.replyToken,
                messages: [
                  {
                    type: "text",
                    text:
                      "ヘルプメッセージ\n" +
                      "無料ユーザーは1日3回まで使用できます。\n" +
                      "サブスクは月額1000円です。\n" +
                      "\n" +
                      "コマンド\n" +
                      "subscribe: サブスクリプションを購入します。\n" +
                      "unsubscribe: サブスクリプションを解約します。\n" +
                      "check: サブスクリプションの状態を確認します。\n",
                  },
                ],
              });
              break;
          }
        }

        // 画像メッセージの場合（例: LINEの画像メッセージは type が "image"）
        if (e.type === "message" && e.message.type === "image") {
          try {
            loading_animation(userId);
            const endpoint = `https://api-data.line.me/v2/bot/message/${e.message.id}/content`;
            const response = await fetch(endpoint, {
              method: "GET",
              headers: {
                Authorization: `Bearer ${config.line.messagingApiClient.channelAccessToken}`,
              },
            });
            if (!response.ok) {
              throw new Error(
                `Failed to fetch image content: ${response.statusText}`,
              );
            }
            // 画像のバイナリデータを ArrayBuffer として取得
            const buffer = await response.arrayBuffer();
            // OCR を実施して TalkHistory 配列を取得する
            talkHistories = await imageToTalkHistories(buffer);

            if (talkHistories == null) {
              await client.replyMessage({
                replyToken: e.replyToken,
                messages: [
                  {
                    type: "text",
                    text: "エラーが発生しました。もう一度やり直してください",
                  },
                ],
              });
              return;
            }
          } catch (err) {
            console.error("画像処理エラー:", err);
          }
        }

        if (e.type === "message" && e.message.type === "file") {
          try {
            loading_animation(userId);
            const endpoint = `https://api-data.line.me/v2/bot/message/${e.message.id}/content`;
            const response = await fetch(endpoint, {
              method: "GET",
              headers: {
                Authorization: `Bearer ${config.line.messagingApiClient.channelAccessToken}`,
              },
            });
            if (!response.ok) {
              throw new Error(
                `Failed to fetch content: ${response.statusText}`,
              );
            }
            const buffer = await response.arrayBuffer();

            const decoder = new TextDecoder("utf-8");
            const talk = decoder.decode(buffer);
            // TODO: こちらに関して、多言語に対応する必要がある

            talkHistories = parseTalkHistories(talk, selfName);

            if (talkHistories == null) {
              await client.replyMessage({
                replyToken: e.replyToken,
                messages: [
                  {
                    type: "text",
                    text: "エラーが発生しました。もう一度やり直してください",
                  },
                ],
              });
              return;
            }
          } catch (e) {
            console.error("Error:", e);
          }
        }

        let message: string[] = [];
        let status = "error";
        if (talkHistories === null) {
          return;
        }
        const talkResponse = await sendTalkRequest(talkHistories, userId);
        if (talkResponse) {
          message = talkResponse.messages;
          status = talkResponse.status;
        }
        if (status === "limit") {
          await client.replyMessage({
            replyToken: e.replyToken,
            messages: [
              {
                type: "text",
                text: "無料ユーザーは1日3回まで使用できます。",
              },
            ],
          });
          return;
        }
        if (status === "error") {
          await client.replyMessage({
            replyToken: e.replyToken,
            messages: [
              {
                type: "text",
                text: "エラーが発生しました。もう一度やり直してください",
              },
            ],
          });
          return;
        }
        // biome-ignore lint/suspicious/noExplicitAny: <explanation>
        const messages: any[] = [];
        // biome-ignore lint/suspicious/noExplicitAny: <explanation>
        const choices: any[] = [];
        for (let i = 0; i < message.length; i++) {
          const index = i;
          const noQuotationMessage = message[i]
            .replace(/\「|\」/g, "")
            .replace(/\n+$/, "");
          choices.push({
            type: "text",
            text: `${index + 1}: ${noQuotationMessage}`,
          });
          messages.push({
            type: "text",
            text: noQuotationMessage,
          });
        }
        // biome-ignore lint/suspicious/noExplicitAny: <explanation>
        const buttonTemplateMessage: any = {
          type: "template",
          altText: "This is a buttons template",
          template: {
            type: "buttons",
            imageAspectRatio: "rectangle",
            imageSize: "cover",
            title: "どのメッセージをコピーしますか？",
            text: "番号を選んでください",
            actions: [
              {
                type: "clipboard",
                label: "1",
                clipboardText: messages[0].text,
              },
              {
                type: "clipboard",
                label: "2",
                clipboardText: messages[1].text,
              },
              {
                type: "clipboard",
                label: "3",
                clipboardText: messages[2].text,
              },
            ],
          },
        };
        choices.push(buttonTemplateMessage);
        await client.replyMessage({
          replyToken: e.replyToken,
          messages: choices,
        });
      });

      await Promise.all(eventPromises);
    }

    res.status(200).send("OK");
  } catch (err) {
    console.error("Error in webhookHandler:", err);
    res.status(500).send("Error");
  }
};
