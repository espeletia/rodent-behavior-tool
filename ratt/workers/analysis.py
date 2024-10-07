import sys
import json

from natsqueue.client import NatsClient
from utils.logger import log_error, log_message
from utils.gracefulTermination import GracefulTermination
from detection.detection import Detection


class AnalystWorker:

    def __init__(self, detector: Detection):
        self.detector = detector

    async def start(self):
        nats_client = NatsClient()

        await nats_client.connect()
        killer = GracefulTermination()
        while not killer.terminate:
            try:
                msgs = await nats_client.pull_analyst_job()
                log_message(msgs)
                for msg in msgs:
                    parsed_msg = json.loads(msg.data)
                    job_id = parsed_msg["job_id"]
                    url = parsed_msg["url"]
                    bucket, key = parse_s3_url(url)
                    results, center_points = self.detector.S3VideoDetection(
                        key)
                    encoding_job_message = {
                        "Message": {
                            "id": job_id,
                            "url": f"S3://{bucket}/{results}",
                        },
                        "Err": None,
                    }

                    await nats_client.publish_video_encoding_job(encoding_job_message)
                    response_object = {
                        "job_id": job_id,
                        "results": results,
                    }
                    await nats_client.publish_analyst_result(response_object)
                    await msg.ack()
            except Exception as err:
                if str(err) != "nats: timeout":
                    log_error(err)
                if str(err) == "nats: connection closed":
                    sys.exit(1)


def parse_s3_url(s3_url):
    # Remove the "s3://" prefix
    stripped_url = s3_url.replace("s3://", "")

    # Split into bucket and key
    parts = stripped_url.split('/', 1)

    # parts[0] is the bucket name, parts[1] is the key
    bucket = parts[0]
    key = parts[1] if len(parts) > 1 else None

    return bucket, key
