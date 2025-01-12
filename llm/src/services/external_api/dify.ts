import axios from "axios";
import type { TalkHistory } from "../../gen/grpc/llm/v1/llm_pb";
import dotenv from "dotenv";

dotenv.config();

const DIFY_API_URL = process.env.DIFY_API_URL;
const DIFY_API_KEY = process.env.DIFY_API_KEY;

export const difyClient = axios.create({
  baseURL: DIFY_API_URL,
  headers: {
    Authorization: `Bearer ${DIFY_API_KEY}`,
    "Content-Type": "application/json",
  },
});

export const callDifyTalk = async (
  input: TalkHistory[],
): Promise<{ replies: string[] }> => {
  try {
    const response = await difyClient.post("/chat-messages", {
      query: input,
      inputs: {},
      response_mode: "blocking",
      user: "abc-123",
      conversation_id: "",
      files: [],
      auto_generate_name: true,
    });

    return JSON.parse(response.data.answer);
  } catch (error) {
    console.error("Error calling Dify API:", error);
    throw new Error("Dify API call failed");
  }
};
