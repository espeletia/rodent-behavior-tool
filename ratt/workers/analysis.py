import sys
import json

from natsqueue.client import NatsClient
from utils.logger import log_error, log_message
from utils.gracefulTermination import GracefulTermination
from detection.detection import Detection
from urllib.parse import urlparse


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
                    job_id = parsed_msg["Message"]["job_id"]
                    video_id = parsed_msg["Message"]["video_id"]
                    media_id = parsed_msg["Message"]["media_id"]
                    url = parsed_msg["Message"]["url"]
                    bucket, key = parse_s3_url(url)
                    results, center_points = self.detector.S3VideoDetection(
                        key)
                    analyst_job_result_message = {
                        "Message": {
                            "job_id": job_id,
                            "video_id": video_id,
                            "media_id": media_id,
                            # TODO: Make it fully configurable
                            "url": f"http://minio:9000/{bucket}/{results}",
                        },
                        "Err": None,
                    }

                    await nats_client.publish_analyst_result(analyst_job_result_message)
                    await msg.ack()
            except Exception as err:
                if str(err) != "nats: timeout":
                    log_error(err)
                if str(err) == "nats: connection closed":
                    sys.exit(1)


def parse_s3_url(s3_url):
    # Parse the URL using urlparse
    parsed_url = urlparse(s3_url)

    # The path part of the URL starts with '/', so we strip it
    stripped_path = parsed_url.path.lstrip('/')

    # Split the path into bucket and key
    parts = stripped_path.split('/', 1)

    # parts[0] is the bucket name, parts[1] is the key
    bucket = parts[0]
    key = parts[1] if len(parts) > 1 else None

    return bucket, key
