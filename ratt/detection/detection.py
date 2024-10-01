from ultralytics import YOLO
from file_manager.file_manager import FileManager


class Detection:

    def __init__(self, modelPath: str, s3_manager: FileManager):
        self.model = YOLO(modelPath)
        self.s3_manager = s3_manager

    def imageDetection(self, imagePath):
        self.results = self.model(imagePath)

    def videoDetection(self, videoPath):
        self.results = self.model(videoPath, stream=True)

    def analyzeResults(self):
        if not self.results:
            raise ("results not defined")
        centerPoints = []
        for result in self.results:
            boxes = result.boxes  # Get the bounding boxes
            for box in boxes:
                # Each box has xywh format (center x, center y, width, height)
                x1, y1, x2, y2 = box.xyxy[0]  # Get the bounding box corners

                # You now have the coordinates
                # print(f"Top-left: ({x1}, {y1}), Bottom-right: ({x2}, {y2})")
                center_x = (x1 + x2) / 2
                center_y = (y1 + y2) / 2
                # print(f"Center: [{center_x}, {center_y}]")
                centerPoints.append((center_x.item(), center_y.item()))
        return centerPoints

    def S3VideoDetection(self, key):
        video = self.s3_manager.download_file(key)
        self.videoDetection(video)
        return self.analyzeResults()
