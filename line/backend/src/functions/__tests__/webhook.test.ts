import { it } from "vitest";
import { parseTalkHistories } from "../webhook.js";

it("トーク履歴が正しくparseされるか", () => {
  const text = `[LINE] 太郎とのトーク履歴
保存日時：2025/1/13 1:21

2025/1/12(日)
0:33	太郎	こんにちは！
0:34	太郎	元気？
`;
  const res = parseTalkHistories(text, "ぽ");
  console.log(res);
});
