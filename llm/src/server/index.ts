import type { ConnectRouter } from "@connectrpc/connect";
import { fastifyConnectPlugin } from "@connectrpc/connect-fastify";
import fastify from "fastify";
import { LlmService } from "../gen/grpc/llm/v1/llm_pb";
import { llmServiceImpl } from "../services/llmService";

(async () => {
  const server = fastify({
    http2: true,
    logger: true,
  });

  // サービス定義を実装を紐づける
  const routes = (router: ConnectRouter) => {
    router.service(LlmService, llmServiceImpl);
  };

  await server.register(fastifyConnectPlugin, {
    routes,
  });

  await server.ready();

  await server.listen({
    host: "0.0.0.0",
    port: 50052,
  });
})();
