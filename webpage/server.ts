import cors from "cors";
import express, { type Express, type Request, type Response } from "express";

const app: Express = express();
const port = 4001;

app.use(cors());

app.get("/", (req: Request, res: Response) => {
  res.send("Hello World!");
});

app.listen(port, () => {
  console.log(`listening on port ${port}`);
});
