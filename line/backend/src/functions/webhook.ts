import { messagingApi } from '@line/bot-sdk';
import type { Request, Response } from 'express';
import { config } from '../config.js';

const { MessagingApiClient } = messagingApi;

const client = new MessagingApiClient(config.line.messagingApiClient);


// type User = {
// // user_id: string; //TODO: 【確認】talk履歴からだと相手のuser_idは取得できない？こちらでuuidを割り振る？
// name: string;
// };

type User = {
    // user_id: string; //TODO: 【確認】talk履歴からだと相手のuser_idは取得できない？こちらでuuidを割り振る？
    name: string; // ユーザー名
};

type TalkHistory = {
date: string;
user: User;
message: string;
};

function parseLineTalkHistory(talk: string): TalkHistory[] {
const lines = talk.split("\n");
const talkHistory: TalkHistory[] = [];
let talkDate: string | null = null; // 現在のメッセージの日付

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
                            { type: 'text', text: 'hoge'},
                        ],
                    });
                }
                if (e.type === 'message' && e.message.type === 'file') {
                    console.log('res:', e);
                    try {
                        const endpoint = `https://api-data.line.me/v2/bot/message/${e.message.id}/content`
                        const accessToken = 'o4t50dLq7Gkhg300mo6qQxSpRCk9af2Pt4alDsj80y7+BgoSnpKRXAnU53Y5ok2A7QyabweFvERSJWr2q9Iy/CM/iq2YRN/t2RpHedz0CS377cJyvh9ehc+MqcmaK1f1hB9jpobeSlOzbYwnVlUjLgdB04t89/1O/w1cDnyilFU=';

                        const response = await fetch(endpoint, {
                            method: 'GET',
                            headers: {
                                Authorization: `Bearer ${accessToken}`
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

                        const talkHistorys = parseLineTalkHistory(talk);
                        console.log('talkHistory:' ,talkHistorys);
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
