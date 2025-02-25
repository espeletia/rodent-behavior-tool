import subprocess
import time
import requests
import os


def capture_video(duration=10, output_file='/home/pi/video.mp4'):
    try:
        # libcamera-vid -t 10000 --autofocus-mode continuous --width 1920 --height 1080 -o video.mp4
        # libcamera-vid -t 10000 --autofocus-mode continuous --width 1920 --height 1080 --roi 0.125,0.125,0.75,0.75 -o video_test.mp4
        command = [
            'libcamera-vid',
            '-t', str(duration * 1000),
            '--width', '1920',
            '--height', '1080',
            '--roi', '0.125,0.125,0.75,0.75'  # adjustement for v2 cage
            '--output', output_file
        ]

        subprocess.run(command, check=True)

        print(f"Video captured successfully: {output_file}")
        return output_file
    except subprocess.CalledProcessError as e:
        print(f"Error capturing video: {e}")
        return None


def send_video_to_api(video_path, api_endpoint):
    try:
        with open(video_path, 'rb') as file:
            files = {'file': (os.path.basename(video_path), file, 'video/mp4')}
            headers = {
                'Accept': 'application/json',
            }
            response = requests.put(api_endpoint, files=files, headers=headers)
            response.raise_for_status()
            print(f"Video uploaded successfully: {response.status_code}")
            return response.json()
    except requests.exceptions.RequestException as e:
        print(f"Error uploading video: {e}")
        return None


def main_loop(api_url, duration=10):
    while True:
        try:
            current_timestamp = int(time.time())
            output_file = f"./videos/video_{current_timestamp}.mp4"
            video_file = capture_video(duration, output_file)

            if video_file:
                send_video_to_api(video_file, api_url)
                os.remove(video_file)
                print(f"Deleted file: {video_file}")

        except KeyboardInterrupt:
            print("Loop interrupted by user.")
            break
        except Exception as e:
            print(f"An error occurred: {e}")


if __name__ == "__main__":
    api_url = "http://192.168.1.106:8081/v1/upload"
    main_loop(api_url)
