import requests
from lcd import SmallDisplay

API_URL = "http://192.168.0.111:8081"


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
    try:
        response = requests.post(API_URL + "/cages")
        response.raise_for_status()
        if 'application/json' in response.headers.get('Content-Type', ''):
            response_json = response.json()
            activation_code = response_json.get('activation_code')
            secret_token = response_json.get('secret_token')
            return activation_code, secret_token
        else:
            raise ValueError("Response is not JSON")
    except requests.exceptions.RequestException as e:
        print(f"Error initializing cage: {e}")
        return None, None


# Example usage:
if __name__ == "__main__":
    display = SmallDisplay()
    url = "/"
    response = make_request(url)
    activation, secret_token = init_cage()
    if activation and secret_token:
        display.DrawText(activation)
        print(activation, secret_token)
    else:
        print("Failed to initialize cage.")
