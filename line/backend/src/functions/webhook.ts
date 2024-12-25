import { config } from "@/config.js";
import { messagingApi } from "@line/bot-sdk";
import type { Request, Response } from "express";
import { createClient } from '@connectrpc/connect';
import { createGrpcTransport } from '@connectrpc/connect-node';
import { TalkRequestSchema, BackendService } from "@/gen/grpc/backend/v1/backend_pb.js";
import { create } from "@bufbuild/protobuf";

const { MessagingApiClient } = messagingApi;

const client = new MessagingApiClient(config.line.messagingApiClient);

const transport = createGrpcTransport({
  baseUrl: 'http://localhost:50051',
});
console.log("transport:", transport);
export const BackendClient = createClient(BackendService, transport);


type User = {
  user_id: string;
  name: string;
};

type TalkHistories = {
date: string;
user: User;
message: string;
};

async function sendTalkRequest(talkHistories: TalkHistories[]): Promise<string[]> {


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

function parseTalkHistories(talk: string, hostUserName:string): TalkHistories[] {
const rows = talk.split("\n");
const TalkHistories: TalkHistories[] = [];
let talkDate: string | null = null;

rows.forEach(row => {
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
      const dateTime = `${talkDate}T${time}:00+0900`; // ISO 8601形式

      const name = userName || "Unknown";
      const user_id =
      name === hostUserName
        ? `${hostUserName}01`
        : name === "Unknown"
        ? "Unknown"
        : `${name}02`;

      TalkHistories.push({
      date: dateTime,
      user: { name, user_id },
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
      const eventPromises = req.body.events.map(async (e: any) => {
        if (e.type === "message" && e.message.type === "text") {
          console.log("Replying to message:", e.message.text);
          await client.replyMessage({
            replyToken: e.replyToken,
            messages: [
              { type: "text", text: e.message.text },
            ],
          });
        }


      if (e.type === 'message' && e.message.type === 'file') {
        console.log('res:', e);
        try {
            const endpoint = `https://api-data.line.me/v2/bot/message/${e.message.id}/content`
            console.log("env token", config.line.messagingApiClient.channelAccessToken)
            const response = await fetch(endpoint, {
                method: 'GET',
                headers: {
                    Authorization: `Bearer ${config.line.messagingApiClient.channelAccessToken}`
                }
            });
            if (!response.ok) {
                throw new Error(`Failed to fetch content: ${response.statusText}`);
            }
            const buffer = await response.arrayBuffer();
            console.log('fileRes:', response);

            const decoder = new TextDecoder('utf-8');
            const talk = decoder.decode(buffer);
            console.log('file contents:', talk);
            const match = talk.match(/\[LINE\] (.*?)とのトーク履歴/);
            let hostUserName = "noName";
            if (match && match[1]) {
                hostUserName = match[1];
                console.log("Hostname:",hostUserName);
            }

            const TalkHistories = parseTalkHistories(talk, hostUserName);
            console.log('TalkHistories:' ,TalkHistories);
            let message:string[] = [];
            if (TalkHistories){
                message = await sendTalkRequest(TalkHistories);
            }
            const messages: any[] = [];
            for (let i = 0; i < message.length; i++){
                messages.push({ type: "text", text: message[i] });
            }
            await client.replyMessage({
              replyToken: e.replyToken,
              messages: messages,
            });
        } catch (e) {
            console.log('Error', e);
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
