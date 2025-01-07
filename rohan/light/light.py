import board
import busio
import time
import adafruit_tsl2591


class LightSensor():
    def __init__(self):
        # Initialize I2C bus
        self.i2c = busio.I2C(board.SCL, board.SDA)

        # Initialize the TSL2591 sensor
        self.sensor = adafruit_tsl2591.TSL2591(self.i2c)

    def read_tsl2591(self):
        try:
            # Get raw and calculated lux values
            full_spectrum = self.sensor.full_spectrum
            infrared = self.sensor.infrared
            visible = self.sensor.visible
            lux = self.sensor.lux

        # Print the readings
            print(f"Full Spectrum (IR + Visible): {full_spectrum}")
            print(f"Infrared light level: {infrared}")
            print(f"Visible light level: {visible}")
            print(f"Lux (calculated): {lux}")
        except Exception as e:
            print(f"Error reading sensor: {e}")

# Read and display sensor data
