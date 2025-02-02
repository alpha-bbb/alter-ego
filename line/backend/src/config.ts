import dotenv from "dotenv";

dotenv.config();

export const config = {
  port: process.env.PORT || 3000,
  line: {
    messagingApiClient: {
      channelAccessToken: process.env.LINE_CHANNEL_ACCESS_TOKEN || "",
      channelSecret: process.env.LINE_CHANNEL_SECRET || "",
    },
  },
  backend: {
    url: process.env.BACKEND_URL || "",
  },
  gcp: {
    keyFilename: process.env.GOOGLE_APPLICATION_CREDENTIALS || "",
  },
};
