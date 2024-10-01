import sys
import utils.logger as logger
import json
from nats.aio.client import Client
import configuration.config as config


class NatsClient:

    def __init__(self):
        self.analysis_sub = None
        self.conn = Client()
        self.js = None

    async def connect(self):
        try:
            await self.conn.connect(config.NATS_URL, max_reconnect_attempts=5)
            self.js = self.conn.jetstream()

        except Exception as err:
            print(err)
            sys.exit(1)

    async def close(self):
        try:
            await self.conn.close()
        except Exception as err:
            print(err)

    async def pull_analyst_job(self):
        if self.analysis_sub is None:
            self.analysis_sub = await self.js.pull_subscribe(durable=config.NATS_ANALYSIS_GROUP,
                                                             subject=config.NATS_ANALYSIS_SUBJECT,
                                                             stream=config.NATS_STREAM)
        return await self.analysis_sub.fetch(config.NATS_BATCH_SIZE)

    async def publish_analyst_result(self, message):
        logger.log_message(config.NATS_ANALYSIS_RESULT_SUBJECT)
        logger.log_message(config.NATS_STREAM)
        ack = await self.js.publish(f"{config.NATS_ANALYSIS_RESULT_SUBJECT}",
                                    json.dumps(message).encode())
