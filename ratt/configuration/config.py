import os

FLASK_PORT = os.environ.get("FLASK_PORT", 8080)

S3_SHARED_CREDENTIALS_FILE = os.environ.get(
    "S3_SHARED_CREDENTIALS_FILE", "/app/configuration/creds")
S3_ENDPOINT = os.environ.get("S3_ENDPOINT", "http://localhost:9000")
S3_ACCESS_KEY = os.environ.get("S3_ACCESS_KEY", "minio123")
S3_SECRET_ACCESS_KEY = os.environ.get("S3_SECRET_ACCESS_KEY", "minio123")

WHEIGHTS_PATH = os.environ.get("WHEIGHTS_PATH", "./models/best.pt")

NATS_URL = os.environ.get("NATS_URL", "nats://localhost:4222")
NATS_STREAM = os.environ.get("NATS_STREAM", "ANALYST")
NATS_JOB_STATUS_SUBJECT = os.environ.get(
    "NATS_JOB_STATUS_SUBJECT", "analysisJobStatus")
NATS_BATCH_SIZE = os.environ.get("NATS_BATCH_SIZE", 1)

NATS_ANALYSIS_GROUP = os.environ.get("NATS_JOB_GROUP", "analystWorker")
NATS_ANALYSIS_SUBJECT = os.environ.get("NATS_JOB_SUBJECT", "analystJob")
NATS_ANALYSIS_RESULT_SUBJECT = os.environ.get(
    "NATS_JOB_RESULT_SUBJECT", "ANALYST.analystJobFinished")

NATS_ENCODER_STREAM = os.environ.get("NATS_ENCODER_STREAM", "ENCODER")
NATS_ENCODER_VIDEO_JOB_SUBJECT = os.environ.get(
    "NATS_ENCODER_VIDEO_JOB_SUBJECT", "ENCODER.mediaVideoJobs")
