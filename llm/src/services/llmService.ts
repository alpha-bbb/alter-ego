import type { ServiceImpl } from "@connectrpc/connect";
import type { LlmService } from "../gen/grpc/llm/v1/llm_pb";
import { callDifyTalk } from "./external_api/dify";

export const llmServiceImpl: ServiceImpl<typeof LlmService> = {
  async talk(req, _context) {
    try {
      const input = req.histories;

      const messages = await callDifyTalk(input);

      return {
        message: messages.replies,
      };
    } catch (error) {
      console.error("Error in talk:", error);
      throw new Error("Internal server error");
    }
  },
};
