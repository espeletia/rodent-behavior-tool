import time
import adafruit_dht
import board


class TempSensor():
    def __init__(self):
        self.dht_device = adafruit_dht.DHT22(board.D27)

    def fetch(self):
        try:
            temperature_c = self.dht_device.temperature
            humidity = self.dht_device.humidity

            if temperature_c is not None and humidity is not None:
                return temperature_c, int(humidity)
                temperature_f = temperature_c * (9 / 5) + 32
                print("Temp: {:.1f} C / {:.1f} F    Humidity: {:.1f}%".format(
                    temperature_c, temperature_f, humidity))
            else:
                print("Sensor reading failed, retrying...")
                return 0, 0

        except RuntimeError as err:
            # Common error for reading failures
            print(f"RuntimeError: {err.args[0]}")

        except OverflowError:
            print("OverflowError: Sensor data invalid, skipping this read.")

        except Exception as e:
            print(f"Unexpected error: {e}")
            print("Reinitializing sensor...")
