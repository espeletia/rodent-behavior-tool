import requests
import time
import json
import os

from lcd.lcd import SmallDisplay
from ultrasonic.ultrasonic import UltrasonicSensor

API_URL = "http://192.168.0.111:8081"
CONFIG_FILE = "cage_config.json"


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


def init_cage():
    if os.path.exists(CONFIG_FILE):
        with open(CONFIG_FILE, 'r') as file:
            config = json.load(file)
            return config.get('activation_code'), config.get('secret_token')

    try:
        response = requests.post(API_URL + "/cages")
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


# Example usage:
if __name__ == "__main__":
    food = UltrasonicSensor(24, 18)
    water = UltrasonicSensor(23, 17)
    display = SmallDisplay()
    url = "/"
    response = make_request(url)
    activation, secret_token = init_cage()
    if activation and secret_token:
        print(activation, secret_token)
        display.draw_text(activation)
        user_id = poll_cage_status(secret_token)
        while user_id is None:
            time.sleep(0.5)
            user_id = poll_cage_status(secret_token)
        if user_id is not None:
            print(user_id)
            display.draw_success(user_id)

    else:
        print("Failed to initialize cage.")
