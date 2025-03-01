import { SubscribeRequestSchema } from "@/gen/grpc/backend/v1/backend_pb.js";

import { config } from "@/config.js";
import {
  Account_PlatformType,
  BackendService,
} from "@/gen/grpc/backend/v1/backend_pb.js";
import { create } from "@bufbuild/protobuf";
import { createClient } from "@connectrpc/connect";
import { createGrpcTransport } from "@connectrpc/connect-node";
import type {} from "express";

export async function sendSubscribeRequest(
  accountId: string,
  action: number,
): Promise<{
  status: "UNSPECIFIED" | "NOT_SUBSCRIBED" | "PROCESSING" | "ACTIVE";
  messages: string;
  redirect_url: string;
  expire_at: number;
} | null> {
  try {
    const request = create(SubscribeRequestSchema, {
      account: {
        platformType: Account_PlatformType.PLATFORM_LINE,
        accountId: accountId,
      },
      action: action,
    });

    const transport = createGrpcTransport({
      baseUrl: config.backend.url,
    });
    const BackendClient = createClient(BackendService, transport);

    const response = await BackendClient.subscribe(request);
    let status: "UNSPECIFIED" | "NOT_SUBSCRIBED" | "PROCESSING" | "ACTIVE";
    console.log(response);
    switch (response.status) {
      case 1:
        status = "NOT_SUBSCRIBED";
        break;
      case 2:
        status = "PROCESSING";
        break;
      case 3:
        status = "ACTIVE";
        break;
      default:
        status = "UNSPECIFIED";
        break;
    }
    return {
      status: status,
      messages: response.message,
      redirect_url: response.redirectUrl,
      expire_at: Number(response.expireAt),
    };
  } catch (error) {
    console.error("Error:", error);
    return null;
  }
}
