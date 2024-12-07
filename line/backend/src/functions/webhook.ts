import { config } from "@/config.js";
import { messagingApi } from "@line/bot-sdk";
import type { Request, Response } from "express";

const { MessagingApiClient } = messagingApi;

const client = new MessagingApiClient(config.line.messagingApiClient);

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
              {
                type: "text",
                text: "https://liff.line.me/2006595464-xg40E95E",
              },
            ],
          });
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
