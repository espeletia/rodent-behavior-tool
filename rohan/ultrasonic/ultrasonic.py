from gpiozero import DistanceSensor


class UltrasonicSensor():

    def __init__(self, echo_pin, trigger_pin):
        self.sensor = DistanceSensor(echo=echo_pin, trigger=trigger_pin)

    def get_distance_cm(self):
        return round(self.sensor.distance * 100, 2)
