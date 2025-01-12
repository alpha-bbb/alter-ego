import { config } from "@/config.js";
import { copyHandler } from "@/functions/copy.js";
import { webhookHandler } from "@/functions/webhook.js";
import express from "express";

export const app = express();

app.use(express.json());

app.get("/", (_, res) => {
  res.status(200).send("Hello, World!");
});
app.get("/webhook", webhookHandler);
app.post("/webhook", webhookHandler);
app.get("/copy", copyHandler);

app.listen(config.port, () => {
  console.log(`http://localhost:${config.port}/`);
});
