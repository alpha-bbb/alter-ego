import { config } from "@/config.js";
import {
  BackendService,
  SubmitUserChoiceRequestSchema,
  TalkRequestSchema,
} from "@/gen/grpc/backend/v1/backend_pb.js";
import { create } from "@bufbuild/protobuf";
import { createClient } from "@connectrpc/connect";
import { createGrpcTransport } from "@connectrpc/connect-node";
import { messagingApi } from "@line/bot-sdk";
import type { Request, Response } from "express";
import express from "express";
import sharp from "sharp";
// import { createWorker } from "tesseract.js";

const app = express();
app.use(express.json());

const { MessagingApiClient } = messagingApi;

const client = new MessagingApiClient(config.line.messagingApiClient);

const transport = createGrpcTransport({
  baseUrl: config.backend.url,
});
console.log("transport:", transport);
export const BackendClient = createClient(BackendService, transport);

type User = {
  user_id: string;
  role: string;
  name: string;
};

type TalkHistories = {
  date: string;
  user: User;
  message: string;
};

async function sendTalkRequest(
  talkHistories: TalkHistories[],
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

function parseTalkHistories(talk: string, yourName: string): TalkHistories[] {
  const rows = talk.split("\n");
  const TalkHistories: TalkHistories[] = [];
  let talkDate: string | null = null;

  // biome-ignore lint/complexity/noForEach: <explanation>
  rows.forEach((row) => {
    const trimmedRow = row.trim();
    // 日付
    const dateMatch = trimmedRow.match(/^(\d{4}\/\d{2}\/\d{2})/);
    if (dateMatch) {
      talkDate = dateMatch[1].replace(/\//g, "-"); // YYYY-MM-DD
      return [];
    }

    // メッセージ（例: "22:07   Test    おはよう"）
    const messageMatch = trimmedRow.match(/^(\d{2}:\d{2})\t+([^\t]+)?\t+(.+)$/);
    if (messageMatch && talkDate) {
      const [_, time, userName, message] = messageMatch;
      console.log("name:", userName);
      const dateTime = `${talkDate}T${time}:00+0900`; // ISO 8601形式

      const name = userName || "Unknown";
      const user_id =
        name === yourName
          ? `${yourName}02`
          : name === "Unknown"
            ? "Unknown"
            : `${name}01`;
      const role =
        name === yourName ? "YOU" : name === "Unknown" ? "UNSPECIFIED" : "SELF";

      TalkHistories.push({
        date: dateTime,
        user: { name, user_id, role },
        message,
      });
    }
  });

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
        // if (e.type === "message" && e.message.type === "text") {
        //   console.log("Replying to message:", e.message.text);
        //   await client.replyMessage({
        //     replyToken: e.replyToken,
        //     messages: [{ type: "text", text: e.message.text }],
        //   });
        // }
        if (e.type === "message" && e.message.type === "image") {
          console.log("res:", e);
          try {
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

            const imageBuffer = await response.arrayBuffer();
            const imageJpeg = await sharp(imageBuffer)
              .jpeg({ quality: 90 }) // 画像の品質を指定 (90%)
              .toBuffer();
            // (async () => {
            //   const worker = await createWorker("jpn");
            //   const ret = await worker.recognize(imageJpeg);
            //   console.log("image text!:", ret.data.text);
            //   await worker.terminate();
            // })();

            const DIFY_API_URL = "http://dify.alter-ego.jtj.jp/v1";
            const DIFY_API_KEY = "app-UXq7yOxpjZBRKBQoBljCFaTP";

            const formData = new FormData();
            const mimeType = "image/jpeg";
            const fileName = "image.jpg";
            const blob = new Blob([imageJpeg], { type: mimeType });
            const file = new File([blob], fileName, { type: mimeType });
            formData.append("file", file, file.name);
            formData.append("user", "abc-123");
            const uploadImage = await fetch(`${DIFY_API_URL}/files/upload`, {
              method: "POST",
              headers: {
                Authorization: `Bearer ${DIFY_API_KEY}`,
              },
              body: formData,
            });
            // const uploadedRes = JSON.stringify(uploadImage);
            console.log("uploadedRes:", uploadImage);
            if (!response.ok) {
              throw new Error(
                `Failed to fetch content: ${response.statusText}`,
              );
            }
            // const talkResponse = await fetch(`${DIFY_API_URL}/chat-messages`, {
            //   method: "POST",
            //   headers: {
            //     Authorization: `Bearer ${DIFY_API_KEY}`,
            //     "Content-Type": "application/json",
            //   },
            //   body: JSON.stringify({
            //     files: [
            //       {
            //         type: "image",
            //         transfer_method: "local_file",
            //         upload_file_id: uploadedRes,
            //       },
            //     ],
            //   }),
            // });
            if (!response.ok) {
              throw new Error(
                `Failed to fetch content: ${response.statusText}`,
              );
            }

            await client.replyMessage({
              replyToken: e.replyToken,
              messages: [
                {
                  type: "text",
                  text: "画像を受け取りました。処理中です。",
                },
              ],
            });
          } catch (e) {
            console.log("Error", e);
          }
        }
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
            const match = talk.match(/\[LINE\] (.*?)とのトーク履歴/);
            let yourName = "noName";
            // biome-ignore lint/complexity/useOptionalChain: <explanation>
            if (match && match[1]) {
              yourName = match[1];
              console.log("Hostname:", yourName);
            }

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
