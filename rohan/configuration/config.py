import os

API_URL = os.environ.get("API_URL", "http://192.168.0.21:8081")
DURATION = os.environ.get("DURATION", 10)
CONFIG_FILE = os.environ.get("CONFIG_FILE", "cage_config.json")
