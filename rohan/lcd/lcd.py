import time
import board
import busio
import gpiozero

from PIL import Image, ImageDraw, ImageFont
import adafruit_ssd1306

import subprocess

# Use gpiozero to control the reset pin

# Display Parameters
WIDTH = 128
HEIGHT = 64
BORDER = 5

# Display Refresh
LOOPTIME = 1.0

# Use I2C for communication


class SmallDisplay():

    def __init__(self):
        oled_reset_pin = gpiozero.OutputDevice(
            4, active_high=False)  # GPIO 4 for reset, active low
        i2c = board.I2C()

        # Manually reset the display (high -> low -> high for reset pulse)
        oled_reset_pin.on()
        time.sleep(0.1)  # Delay for a brief moment
        oled_reset_pin.off()  # Toggle reset pin low
        time.sleep(0.1)  # Wait for reset
        oled_reset_pin.on()  # Turn reset pin back high

        self.oled = adafruit_ssd1306.SSD1306_I2C(WIDTH, HEIGHT, i2c, addr=0x3C)

        self.oled.fill(0)
        self.oled.show()

        image = Image.new("1", (self.oled.width, self.oled.height))

        self.draw = ImageDraw.Draw(image)

        self.draw.rectangle((0, 0, self.oled.width, self.oled.height),
                            outline=255, fill=255)

        self.font = ImageFont.truetype('PixelOperator.ttf', 16)

    def drawText(self, text):
        self.draw.text((0, 0), "ACTIVATION TOKEN:", font=self.font, fill=255)
        self.draw.text((0, 16), text, font=self.font, fill=255)

    def clearDisplay(self):
        self.draw.rectangle(
            (0, 0, self.oled.width, self.oled.height), outline=0, fill=0)


# while True:
#     # Draw a black filled box to clear the image
#     draw.rectangle((0, 0, oled.width, oled.height), outline=0, fill=0)
#
#     # Shell scripts for system monitoring from here : https://unix.stackexchange.com/questions/119126/command-to-display-memory-usage-disk-usage-and-cpu-load
#     cmd = "hostname -I | cut -d\' \' -f1"
#     IP = subprocess.check_output(cmd, shell=True)
#     cmd = "top -bn1 | grep load | awk '{printf \"CPU: %.2f\", $(NF-2)}'"
#     CPU = subprocess.check_output(cmd, shell=True)
#     cmd = "free -m | awk 'NR==2{printf \"%.1f %.1f %.1f\", $3/1024,$2/1024,($3/$2)*100}'"
#     MemUsage = subprocess.check_output(cmd, shell=True)
#     mem_parts = MemUsage.decode('utf-8').strip().split()
#     mem_used_gb = mem_parts[0]
#     mem_total_gb = mem_parts[1]
#     mem_percent = mem_parts[2]
#     mem_display = f"Mem: {mem_used_gb}/{mem_total_gb}GB {mem_percent}%"
#     cmd = "df -h | awk '$NF==\"/\"{printf \"Disk: %d/%dGB %s\", $3,$2,$5}'"
#     Disk = subprocess.check_output(cmd, shell=True)
#     cmd = "vcgencmd measure_temp |cut -f 2 -d '='"
#     Temp = subprocess.check_output(cmd, shell=True)
#
#     # Pi Stats Display
#     draw.text((0, 0), "IP: " + str(IP, 'utf-8'), font=font, fill=255)
#     draw.text((0, 16), str(CPU, 'utf-8') + "LA", font=font, fill=255)
#     draw.text((80, 16), str(Temp, 'utf-8'), font=font, fill=255)
#     draw.text((0, 32), mem_display, font=font, fill=255)
#     draw.text((0, 48), str(Disk, 'utf-8'), font=font, fill=255)
#
#     # Display the image
#     oled.image(image)
#     oled.show()
#
#     # Wait for the next loop
#     time.sleep(LOOPTIME)
