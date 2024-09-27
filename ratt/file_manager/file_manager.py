import boto3
import tempfile
import os
from botocore.exceptions import ClientError
from botocore.client import Config


class FileManager:

    # access_key, secret_access_key, bucket
    def __init__(self, url, access_key, secret_access_key, bucket):
        self.client = boto3.client(
            's3',
            endpoint_url=url,
            aws_access_key_id=access_key,
            aws_secret_access_key=secret_access_key,
            config=Config(signature_version='s3v4'),
        )
        self.bucket = bucket

    def list_buckets(self):
        try:
            buckets = self.client.list_buckets()
            for bucket in buckets['Buckets']:
                print(f"Bucket Name: {bucket['Name']}")
        except ClientError as e:
            print(e)
            return

    def download_file(self, key: str):
        _, suffix = os.path.splitext(key)
        print(suffix)
        try:
            with tempfile.NamedTemporaryFile(delete=False, suffix=suffix) as temp_file:
                self.client.download_file(self.bucket, key, temp_file.name)
                return temp_file.name
        except Exception as e:
            print(f"an error occurred: {e}")
            if os.path.exists(temp_file.name):
                os.remove(temp_file.name)
            return None
