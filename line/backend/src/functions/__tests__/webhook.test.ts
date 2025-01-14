import { expect, it } from "vitest";
import { getYourName, parseTalkHistories } from "../webhook.js";

it("トーク履歴が正しくparseされるか (日本語)", () => {
  const text = `[LINE] 太郎とのトーク履歴
保存日時：2025/1/13 1:21

2025/1/12(日)
0:33	太郎	こんにちは！
0:34	太郎	元気？
`;
  const res = parseTalkHistories(text, "ぽ");
  expect(res).toEqual([
    {
      $typeName: "backend.v1.TalkHistory",
      date: "2025-01-12T00:33:00+0900",
      user: {
        $typeName: "backend.v1.User",
        name: "太郎",
        userId: "太郎01",
        role: 1,
      },
      message: "こんにちは！",
    },
    {
      $typeName: "backend.v1.TalkHistory",
      date: "2025-01-12T00:34:00+0900",
      user: {
        $typeName: "backend.v1.User",
        name: "太郎",
        userId: "太郎01",
        role: 1,
      },
      message: "元気？",
    },
  ]);
});

it("トーク履歴が正しくparseされるか (英語)", () => {
  const text = `Chat history with 太郎
Saved on: 1/13/2025, 1:21

Sun, 1/12/2025
00:33	太郎	こんにちは！
00:34	太郎	元気？
`;
  const res = parseTalkHistories(text, "ぽ");
  expect(res).toEqual([
    {
      $typeName: "backend.v1.TalkHistory",
      date: "2025-01-12T00:33:00+0900",
      user: {
        $typeName: "backend.v1.User",
        name: "太郎",
        userId: "太郎01",
        role: 1,
      },
      message: "こんにちは！",
    },
    {
      $typeName: "backend.v1.TalkHistory",
      date: "2025-01-12T00:34:00+0900",
      user: {
        $typeName: "backend.v1.User",
        name: "太郎",
        userId: "太郎01",
        role: 1,
      },
      message: "元気？",
    },
  ]);
});

it("ファイル名からユーザー名を抽出 (日本語)", () => {
  const filename = "[LINE] 太郎とのトーク.txt";
  const res = getYourName(filename);
  expect(res).toEqual("太郎");
});

it("ファイル名からユーザー名を抽出 (英語)", () => {
  const filename = "Chat history with 太郎.txt";
  const res = getYourName(filename);
  expect(res).toEqual("太郎");
});
