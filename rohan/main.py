import requests
import time
import json
import os
import asyncio

from lcd.lcd import SmallDisplay
from configuration.config import API_URL, DURATION, CONFIG_FILE
from temp.temp import TempSensor
from camera.camera import capture_video, send_video_to_api
from ultrasonic.ultrasonic import UltrasonicSensor
from light.light import LightSensor


def make_request(url, method='GET', data=None, headers=None):
    try:
        if method.upper() == 'GET':
            response = requests.get(url, headers=headers)
        elif method.upper() == 'POST':
            response = requests.post(url, json=data, headers=headers)
        elif method.upper() == 'PUT':
            response = requests.put(url, json=data, headers=headers)
        elif method.upper() == 'DELETE':
            response = requests.delete(url, headers=headers)
        else:
            raise ValueError("Unsupported HTTP method")

        response.raise_for_status()
        return response.json() if 'application/json' in response.headers.get('Content-Type', '') else response.text

    except requests.exceptions.RequestException as e:
        return f"An error occurred: {e}"


async def async_video_task(video_file, lux, foodDistance, waterDistance, current_timestamp, secret_token, temp):
    video_upload_response = await send_video_to_api(video_file, f"{API_URL}/v1/upload")
    if video_upload_response:
        print(f"Lux (calculated): {lux}")
        send_cage_message(
            secret_token,
            int(lux),
            int(foodDistance),
            int(waterDistance),
            1,
            int(temp),
            current_timestamp,
            video_upload_response.get('upload_url')
        )
    # Delete the video file asynchronously
    await asyncio.to_thread(os.remove, video_file)
    print(f"Deleted file: {video_file}")


def run_async_video_task(video_file, lux, foodDistance, waterDistance, current_timestamp, secret_token, temp):
    # Run the asynchronous video task in the background
    asyncio.create_task(async_video_task(
        video_file, lux, foodDistance, waterDistance, current_timestamp, secret_token, temp))


def init_cage():
    if os.path.exists(CONFIG_FILE):
        with open(CONFIG_FILE, 'r') as file:
            config = json.load(file)
            return config.get('activation_code'), config.get('secret_token')

    try:
        response = requests.post(API_URL + "/v1/cages")
        response.raise_for_status()
        if 'application/json' in response.headers.get('Content-Type', ''):
            response_json = response.json()
            activation_code = response_json.get('activation_code')
            secret_token = response_json.get('secret_token')
            with open(CONFIG_FILE, 'w') as file:
                json.dump(response_json, file)
            return activation_code, secret_token
        else:
            raise ValueError("Response is not JSON")
    except requests.exceptions.RequestException as e:
        print(f"Error initializing cage: {e}")
        return None, None


def poll_cage_status(secret_token):
    try:
        headers = {"Authorization": f"Bearer {secret_token}"}
        response = requests.get(f"{API_URL}/internal/cage", headers=headers)
        response.raise_for_status()
        if 'application/json' in response.headers.get('Content-Type', ''):
            response_json = response.json()
            user_id = response_json.get('user_id')
            return user_id
        else:
            raise ValueError("Response is not JSON")
    except requests.exceptions.RequestException as e:
        return f"Error polling cage status: {e}"


def send_cage_message(secret_token, light, food, water, revision, temp, timestamp, videoUrl):
    message = {
        "food": food,
        "humidity": 0,
        "light": light,
        "revision": revision,
        "temp": temp,
        "timestamp": timestamp,
        "video_url": videoUrl,
        "water": water
    }
    try:
        headers = {
            "Authorization": f"Bearer {secret_token}",
            "Content-Type": "application/json"
        }
        response = requests.post(
            f"{API_URL}/internal/cage/message", json=message, headers=headers)
        response.raise_for_status()
        if 'application/json' in response.headers.get('Content-Type', ''):
            response_json = response.json()
            user_id = response_json.get('user_id')
            return user_id
        else:
            raise ValueError("Response is not JSON")
    except requests.exceptions.RequestException as e:
        return f"Error sending cage message: {e}"


async def main():
    food = UltrasonicSensor(24, 18)
    water = UltrasonicSensor(23, 17)
    light = LightSensor()
    display = SmallDisplay()
    dth = TempSensor()
    activation, secret_token = init_cage()
    if activation and secret_token:
        print(activation, secret_token)
        display.draw_text(activation)
        user_id = poll_cage_status(secret_token)
        while user_id is None:
            await asyncio.sleep(0.5)  # Use asyncio.sleep instead of time.sleep
            user_id = poll_cage_status(secret_token)
        if user_id is not None:
            print(user_id)
            display.draw_success(user_id)
            while True:
                current_timestamp = int(time.time())
                output_file = f"./videos/video_{current_timestamp}.mp4"
                video_file = capture_video(DURATION, output_file)
                foodDistance = food.get_distance_cm()
                waterDistance = water.get_distance_cm()
                lux = light.read_tsl2591()
                print(f"food left: {foodDistance}cm")
                print(f"water left: {waterDistance}")
                print(f"video: {video_file}")
                temp, hum = dth.fetch()
                run_async_video_task(
                    video_file, lux, foodDistance, waterDistance, current_timestamp, secret_token, temp)
                # Add a small delay to avoid busy-waiting
                await asyncio.sleep(1)

    else:
        print("Failed to initialize cage.")


# Run the main function in the event loop
if __name__ == "__main__":
    asyncio.run(main())
