import { config } from "@/config.js";
import {
  BackendService,
  SubmitUserChoiceRequestSchema,
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
console.log("transport:", transport);
export const BackendClient = createClient(BackendService, transport);

async function sendTalkRequest(
  talkHistories: TalkHistory[],
): Promise<string[]> {
  try {
    const request = create(TalkRequestSchema, {
      histories: talkHistories,
      actionKind: 1,
    });

    const response = await BackendClient.talk(request);
    console.log("Response:", response.message);
    return response.message;
  } catch (error) {
    console.error("Error:", error);
    return [];
  }
}

// biome-ignore lint/suspicious/noExplicitAny: <explanation>
async function sendQuestionnaire(messageNumber: string): Promise<any> {
  try {
    const request = create(SubmitUserChoiceRequestSchema, {
      choice: messageNumber,
    });

    const response = await BackendClient.submitUserChoice(request);
    console.log("Response:", response);
    return response;
  } catch (error) {
    console.error("Error:", error);
    return [];
  }
}

export function getYourName(talk: string): string {
  const match = talk.match(
    /(?:\[LINE\] |Chat history with )(.+?)(?:とのトーク|\.txt)/,
  );
  return match ? match[1] : "noName";
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
  yourName: string,
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
      console.log("name:", userName);
      const dateTime = `${talkDate}T${time}:00+0900`; // ISO 8601形式

      const name = userName || "Unknown";
      const userId =
        name === yourName
          ? `${yourName}02`
          : name === "Unknown"
            ? "Unknown"
            : `${name}01`;
      const role = (() => {
        switch (name) {
          case yourName:
            return User_UserRole.YOU;
          case "Unknown":
            return User_UserRole.UNSPECIFIED;
          default:
            return User_UserRole.SELF;
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

export const webhookHandler = async (
  req: Request,
  res: Response,
): Promise<void> => {
  try {
    if (req.body.events && req.body.events.length > 0) {
      // biome-ignore lint/suspicious/noExplicitAny: <explanation>
      const eventPromises = req.body.events.map(async (e: any) => {
        if (e.type === "postback") {
          console.log("Postback data:", e.postback.data);
          const messageNumber = e.postback.data;

          await sendQuestionnaire(messageNumber);
          console.log("Questionnaire sent");
        }

        if (e.type === "message" && e.message.type === "file") {
          console.log("res:", e);
          try {
            const endpoint = `https://api-data.line.me/v2/bot/message/${e.message.id}/content`;
            console.log(
              "env token",
              config.line.messagingApiClient.channelAccessToken,
            );
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
            console.log("fileRes:", response);

            const decoder = new TextDecoder("utf-8");
            const talk = decoder.decode(buffer);
            console.log("file contents:", talk);
            // TODO: こちらに関して、多言語に対応する必要がある
            const yourName = getYourName(talk);
            console.log("Hostname:", yourName);

            const TalkHistories = parseTalkHistories(talk, yourName);
            console.log("TalkHistories:", TalkHistories);
            let message: string[] = [];
            if (TalkHistories) {
              message = await sendTalkRequest(TalkHistories);
            }
            // biome-ignore lint/suspicious/noExplicitAny: <explanation>
            const messages: any[] = [];
            // biome-ignore lint/suspicious/noExplicitAny: <explanation>
            const choices: any[] = [];
            for (let i = 0; i < message.length; i++) {
              const index = i;
              choices.push({
                type: "text",
                text: `${index + 1}: ${message[i]}`,
              });
              messages.push({
                type: "text",
                text: message[i],
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
                title: "Suggested messages",
                text: "Which message do you want to copy?",
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
            // biome-ignore lint/suspicious/noExplicitAny: <explanation>
            const buttonTemplateQuestionnaire: any = {
              type: "flex",
              altText: "どのメッセージがよかったですか？",
              contents: {
                type: "bubble",
                body: {
                  type: "box",
                  layout: "vertical",
                  contents: [
                    {
                      type: "text",
                      text: "どのメッセージがよかったですか？",
                      wrap: true,
                      weight: "regular",
                      size: "md",
                      color: "#222222",
                      margin: "none",
                    },
                  ],
                  spacing: "sm",
                },
                footer: {
                  type: "box",
                  layout: "horizontal",
                  contents: [
                    {
                      type: "button",
                      style: "primary",
                      action: {
                        type: "postback",
                        label: "1",
                        data: "1",
                      },
                      color: "#0E71EB",
                      height: "sm",
                    },
                    {
                      type: "button",
                      style: "primary",
                      action: {
                        type: "postback",
                        label: "2",
                        data: "2",
                      },
                      color: "#0E71EB",
                      height: "sm",
                    },
                    {
                      type: "button",
                      style: "primary",
                      action: {
                        type: "postback",
                        label: "3",
                        data: "3",
                      },
                      color: "#0E71EB",
                      height: "sm",
                    },
                  ],
                  spacing: "sm",
                },
              },
            };
            choices.push(buttonTemplateMessage);
            choices.push(buttonTemplateQuestionnaire);
            await client.replyMessage({
              replyToken: e.replyToken,
              messages: choices,
            });
          } catch (e) {
            console.log("Error", e);
          }
        }
      });

      await Promise.all(eventPromises);
    }

    res.status(200).send("OK");
  } catch (err) {
    console.error("Error in webhookHandler:", err);
    res.status(500).send("Error");
  }
};
