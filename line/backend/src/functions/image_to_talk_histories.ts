import {
  type TalkHistory,
  User_UserRole,
} from "@/gen/grpc/backend/v1/backend_pb.js";
import vision from "@google-cloud/vision";
import { imageSize } from "image-size";
import { config } from "../config.js";

/**
 * 画像から OCR を実施し、検出結果を「相手」と「自分」の発言を含む TalkHistory の配列として返す
 * 発言の順序はバウンディングボックスの上端 (y座標) の昇順（時系列順）とする
 *
 * 発言者の判定ルール:
 * - 「相手」: バウンディングボックスの左上の x 座標が **imageWidth * 0.09 ～ imageWidth * 0.11**
 * - 「自分」: バウンディングボックスの右上の x 座標が **imageWidth * 0.85 ～ imageWidth * 0.87**
 *
 * @param img 画像の ArrayBuffer
 * @returns TalkHistory[] 発言の時系列順の配列
 */
export async function imageToTalkHistories(
  img: ArrayBuffer,
): Promise<TalkHistory[]> {
  const imageBuffer = Buffer.from(img);
  // 画像の横幅（width）を取得
  const dimensions = imageSize(imageBuffer);
  const imageWidth = dimensions.width;
  if (!imageWidth) {
    throw new Error("画像の横幅を取得できませんでした。");
  }

  const imgBase64 = imageBuffer.toString("base64");

  const client = new vision.ImageAnnotatorClient({
    keyFilename: config.gcp.keyFilename,
  });

  const [result] = await client.documentTextDetection({
    image: { content: imgBase64 },
  });

  if (!result.fullTextAnnotation) {
    console.error("OCR結果が取得できませんでした。");
    return [];
  }

  const otherMessages: { y: number; history: TalkHistory }[] = [];
  const selfMessages: { y: number; history: TalkHistory }[] = [];

  for (const page of result.fullTextAnnotation.pages || []) {
    for (const block of page.blocks || []) {
      const blockText = block.paragraphs
        ?.map(
          (paragraph) =>
            paragraph.words
              ?.map(
                (word) =>
                  word.symbols?.map((symbol) => symbol.text).join("") ?? "",
              )
              .join("") ?? "",
        )
        .join("\n")
        .trim();

      if (!blockText) continue;

      // バウンディングボックスの頂点を取得
      const vertices = block.boundingBox?.vertices;
      // 左上の頂点（相手判定用）
      const topLeftX = vertices?.[0]?.x ?? 0;
      // 上端の y 座標（ソート用）
      const topY = vertices?.[0]?.y ?? 0;
      // 右上の頂点（自分判定用）
      const topRightX = vertices?.[1]?.x ?? 0;

      // 発言者の判定（画像の横幅から算出した値を使用）
      if (topLeftX >= 110 && topLeftX <= 130) {
        // 「相手」の発言
        otherMessages.push({
          y: topY,
          history: {
            $typeName: "backend.v1.TalkHistory",
            date: new Date().toISOString(),
            user: {
              $typeName: "backend.v1.User",
              userId: "other",
              name: "相手",
              role: User_UserRole.YOU,
            },
            message: blockText,
          },
        });
      } else if (topRightX >= imageWidth - 66 && topRightX <= imageWidth - 44) {
        // 「自分」の発言
        selfMessages.push({
          y: topY,
          history: {
            $typeName: "backend.v1.TalkHistory",
            date: new Date().toISOString(),
            user: {
              $typeName: "backend.v1.User",
              userId: "self",
              name: "自分",
              role: User_UserRole.SELF,
            },
            message: blockText,
          },
        });
      } else {
        // 指定範囲外の場合は警告を出して無視する
        // Debug用で残している
        // console.warn(
        //   `バウンディングボックスの頂点座標が期待範囲外です。topLeftX: ${topLeftX}, topRightX: ${topRightX}`,
        // );
      }
    }
  }

  // 他のメッセージと自分のメッセージを一つの配列に結合し、y 座標（上端）で昇順にソート（時系列順）
  const combinedMessages = [...otherMessages, ...selfMessages];
  combinedMessages.sort((a, b) => a.y - b.y);
  return combinedMessages.map((item) => item.history);
}
