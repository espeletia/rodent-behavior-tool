import requests
import tempfile
import os


class downloader():

    def downloadVideoFromUrl(self, url):
        try:
            with tempfile.NamedTemporaryFile(delete=False, suffix=".mp4") as temp_file:
                print(f"Temporary file created: {temp_file.name}")
                response = requests.get(url, stream=True)
                if response.status_code == 200:
                    for chunk in response.iter_content(chunk_size=1024):
                        temp_file.write(chunk)
                    print("Video stored successfully")
                    return temp_file.name
                else:
                    print(
                        f"Failed to download video. Status code: {response.status_code}")
                    os.remove(temp_file.name)
                    return None
        except Exception as e:
            print(f"an error occurred: {e}")
            if os.path.exists(temp_file.name):
                os.remove(temp_file.name)
            return None
