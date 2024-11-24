import express from 'express';
import { webhookHandler } from './functions/webhook';
import { config } from './config';

export const app = express();

app.use(express.json());

app.get('/webhook', webhookHandler);
app.post('/webhook', webhookHandler);

app.listen(config.port, () => {
    console.log(`http://localhost:${config.port}/`);
});
