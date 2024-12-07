import { messagingApi } from '@line/bot-sdk';
import type { Request, Response } from 'express';
import { config } from '../config.js';

const { MessagingApiClient } = messagingApi;

const client = new MessagingApiClient(config.line.messagingApiClient);

type User = {
    // user_id: string; //TODO: 【要検討】talk履歴からだと相手のuser_idは取得できない
    name: string; // ユーザー名
};

type TalkHistory = {
date: string;
user: User;
message: string;
};

function parseTalkHistory(talk: string): TalkHistory[] {
const lines = talk.split("\n");
const talkHistory: TalkHistory[] = [];
let talkDate: string | null = null;

lines.forEach(line => {
    const trimmedLine = line.trim();
    // 日付
    const dateMatch = trimmedLine.match(/^(\d{4}\/\d{2}\/\d{2})/);
    if (dateMatch) {
    talkDate = dateMatch[1].replace(/\//g, "-"); // YYYY-MM-DD
    return;
    }

    // メッセージ（例: "22:07   Test    おはよう"）
    const messageMatch = trimmedLine.match(/^(\d{2}:\d{2})\s+([^ ]+)?\s+(.+)$/);
    if (messageMatch && talkDate) {
        const [_, time, userName, message] = messageMatch;
        const dateTime = `${talkDate}T${time}:00Z`; // ISO 8601形式

        talkHistory.push({
        date: dateTime,
        user: { name: userName || "Unknown" },
        message,
        });
    }else {
        console.log("cannot parse");
    }
});

return talkHistory;
}

export const webhookHandler = async (req: Request, res: Response): Promise<void> => {
    try {
        if (req.body.events && req.body.events.length > 0) {
            const eventPromises = req.body.events.map(async (e:any) => {
                if (e.type === 'message' && e.message.type === 'text') {
                    console.log('Replying to message:', e.message.text);
                    await client.replyMessage({
                        replyToken: e.replyToken,
                        messages: [
                            { type: 'text', text: e.message.text },
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

                        const talkHistory = parseTalkHistory(talk);
                        console.log('talkHistory:' ,talkHistory);
                    } catch (e) {
                        console.log('Error', e);
                    }
                }
            });

            await Promise.all(eventPromises);
        }

        res.status(200).send('OK');
    } catch (err) {
        console.error('Error in webhookHandler:', err);
        res.status(500).send('Error');
    }
};
