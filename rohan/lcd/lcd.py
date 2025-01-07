import time
import board
import gpiozero
from PIL import Image, ImageDraw, ImageFont
import adafruit_ssd1306

# Display Parameters
WIDTH = 128
HEIGHT = 64
BORDER = 5


class SmallDisplay:
    def __init__(self):
        # GPIO for reset pin
        self.oled_reset_pin = gpiozero.OutputDevice(
            4, active_high=False)  # GPIO 4 for reset, active low
        i2c = board.I2C()

        # Manually reset the display
        self.oled_reset_pin.on()
        time.sleep(0.1)
        self.oled_reset_pin.off()
        time.sleep(0.1)
        self.oled_reset_pin.on()

        # Initialize the OLED display
        self.oled = adafruit_ssd1306.SSD1306_I2C(WIDTH, HEIGHT, i2c, addr=0x3C)
        self.oled.fill(0)
        self.oled.show()

        # Create image and drawing objects
        self.image = Image.new("1", (self.oled.width, self.oled.height))
        self.draw = ImageDraw.Draw(self.image)

        # Load font
        self.font = ImageFont.truetype("PixelOperator.ttf", 16)

    def draw_text(self, text):
        # Clear the image before drawing
        self.clear_display()

        # Draw the text
        self.draw.text((0, 0), "ACTIVATION TOKEN:", font=self.font, fill=255)
        self.draw.text((0, 16), text, font=self.font, fill=255)

        # Update the OLED display
        self.oled.image(self.image)
        self.oled.show()

    def draw_success(self, text):
        # Clear the image before drawing
        self.clear_display()

        # Draw the text
        self.draw.text((0, 0), "Cage Active!", font=self.font, fill=255)
        self.draw.text((0, 16), f"{text}", font=self.font, fill=255)

        # Update the OLED display
        self.oled.image(self.image)
        self.oled.show()

    def clear_display(self):
        # Clear the image
        self.draw.rectangle(
            (0, 0, self.oled.width, self.oled.height), outline=0, fill=0)

        # Update the OLED display
        self.oled.image(self.image)
        self.oled.show()

    def draw_smiley_face(self):
        # Clear the image before drawing
        self.clear_display()

        # Draw the face using individual pixels
        for x in range(32, 97):
            for y in range(16, 81):
                if (x - 64)**2 + (y - 48)**2 <= 32**2:  # Circle equation
                    self.draw.point((x, y), fill=255)

        # Draw the eyes
        for x in range(50, 59):
            for y in range(32, 41):
                self.draw.point((x, y), fill=255)
        for x in range(70, 79):
            for y in range(32, 41):
                self.draw.point((x, y), fill=255)

        # Draw the smile
        for x in range(48, 81):
            y = int(72 - ((32**2 - (x-64)**2)**0.5))  # Semi-circle for smile
            self.draw.point((x, y), fill=255)

        # Update the OLED display
        self.oled.image(self.image)
        self.oled.show()
