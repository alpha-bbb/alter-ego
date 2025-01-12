import express, { type Request, type Response } from "express";

const app = express();
app.use(express.json());

type CopyHistory = {
  text: string;
  count: number;
};

const copyHistories: Record<string, CopyHistory> = {};

export const copyHandler = async (
  req: Request,
  res: Response,
): Promise<void> => {
  try {
    // クエリ文字列からコピーするテキストを取得
    // const textToCopy = req.query.text as string;
    const textToCopy = "fugafuga";

    if (!textToCopy) {
      res.status(400).json({ error: "Missing 'text' query parameter" });
      return;
    }

    // HTMLを返してページ読み込み時にクリップボードにコピーさせる

    res.send(`
      <!DOCTYPE html>
      <html lang="en">
      <head>
        <meta charset="UTF-8">
        <meta name="viewport" content="width=device-width, initial-scale=1.0">
        <title>Clipboard Copy</title>
      </head>
      <body>
        <p>Text to copy: <span id="text-to-copy">${textToCopy}</span></p>
        
        <script>
          // コピーするテキスト
          const textToCopy = "${textToCopy.replace(/"/g, '\\"')}";

          // ページ読み込み時に自動でコピー
          window.addEventListener('load', () => {
            navigator.clipboard.writeText(textToCopy)
              .then(() => {
                console.log("Copied to clipboard:", textToCopy);
                alert("Copied to clipboard: " + textToCopy);
              })
              .catch(err => {
                console.error("Failed to copy:", err);
                alert("Failed to copy: " + err);
              });
          });
        </script>
      </body>
      </html>
    `);

    // 履歴を更新
    if (copyHistories[textToCopy]) {
      copyHistories[textToCopy].count += 1;
    } else {
      copyHistories[textToCopy] = { text: textToCopy, count: 1 };
    }
  } catch (error) {
    console.error("Error in copyHandler:", error);
    res.status(500).json({ error: "Internal Server Error" });
  }
};
