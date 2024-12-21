import { config } from "@/config.js";
import { messagingApi } from "@line/bot-sdk";
import type { Request, Response } from "express";

const { MessagingApiClient } = messagingApi;

const client = new MessagingApiClient(config.line.messagingApiClient);

type User = {
  user_id: string;
  name: string;
};

type TalkHistories = {
date: string;
user: User;
message: string;
};

async function sendTalkRequest(talkHistories: TalkHistories[]): Promise<void> {



  const requestMessage = {
    talkHistories,
    actionKind: 1,
  };

  try {
    const response = await fetch("http://localhost:3000/backend/v1/BackendService/Talk", {
      method: "POST",
      headers: {
        "Content-Type": "application/protobuf",
      },
      body: JSON.stringify(requestMessage),
    });

    if (!response.ok) {
      throw new Error(`HTTP error! status: ${response.status}`);
    }

    console.log("Response:", response);
  } catch (error) {
    console.error("Error:", error);
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
  return;
  }

  // メッセージ（例: "22:07   Test    おはよう"）
  const messageMatch = trimmedRow.match(/^(\d{2}:\d{2})\t+([^\t]+)?\t+(.+)$/);
  if (messageMatch && talkDate) {
      const [_, time, userName, message] = messageMatch;
      const dateTime = `${talkDate}T${time}::00+0900`; // ISO 8601形式

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
  }else {
      console.log("cannot parse");
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
            if (TalkHistories){
                sendTalkRequest(TalkHistories);
            }
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
