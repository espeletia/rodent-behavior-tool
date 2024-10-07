import boto3
import subprocess
import os
from utils import logger
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
        logger.log_message(key)
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

    def upload_file(self, key: str, src: str):
        try:
            # Upload the file to the specified bucket and key
            self.client.upload_file(src, self.bucket, key, ExtraArgs={
                                    'ContentType': 'video/mp4'})
            print(
                f"File '{src}' uploaded to bucket '{self.bucket}' with key '{key}'.")
            os.remove(src)
            return True
        except ClientError as e:
            print(f"An error occurred while uploading the file: {e}")
            return False
        except Exception as e:
            print(f"An unexpected error occurred: {e}")
            return False

    def get_output_path(self, key: str, output_prefix: str = "outputs", output_suffix: str = "boxes_"):
        # Split the path to get the directory and filename
        directory, filename = os.path.split(key)
        # Get the base name of the file without the extension
        base_name, extension = os.path.splitext(filename)
        # Construct the new path
        output_path = os.path.join(
            directory, output_prefix, f"{output_suffix}{base_name}{extension}")
        return output_path
