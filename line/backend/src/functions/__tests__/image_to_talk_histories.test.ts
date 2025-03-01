import fs from "node:fs";
import { imageToTalkHistories } from "@/functions/image_to_talk_histories.js";
import { describe, it } from "vitest";

describe("imageToTalkHistories", () => {
  it("LINEのスクショを送信し、OCRが成功する", async () => {
    const img = fs.readFileSync("tests/fixtures/tmp_1738475450801.jpg");
    const result = await imageToTalkHistories(img);
    console.log(result);
  });
});
