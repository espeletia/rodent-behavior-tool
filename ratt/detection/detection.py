from ultralytics import YOLO
import os
import cv2
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

    def analyzeWithHeatmapAndBoundingBoxes(self, videoPath):
        cap = cv2.VideoCapture(videoPath)
        frame_width = int(cap.get(cv2.CAP_PROP_FRAME_WIDTH))
        frame_height = int(cap.get(cv2.CAP_PROP_FRAME_HEIGHT))
        fps = int(cap.get(cv2.CAP_PROP_FPS))
        outputPath = "/tmp/output.mp4"
        fourcc = cv2.VideoWriter_fourcc(*'mp4v')  # Codec for MP4
        out = cv2.VideoWriter(outputPath, fourcc, fps,
                              (frame_width, frame_height))

        centerPoints = []
        for result in self.results:
            ret, frame = cap.read()
            if not ret:
                break

            boxes = result.boxes

            for box in boxes:
                # bounding box
                x1, y1, x2, y2 = map(int, box.xyxy[0])
                cv2.rectangle(frame, (x1, y1), (x2, y2), (0, 255, 0), 2)

                # center point
                center_x = (x1 + x2) // 2
                center_y = (y1 + y2) // 2
                centerPoints.append((center_x, center_y))
                cv2.circle(frame, (center_x, center_y), 5, (225, 0, 0), -1)

                # confidence
                confidence = box.conf[0]
                label = f'Conf: {confidence:.2f}'
                cv2.putText(frame, label, (x1, y1 - 10),
                            cv2.FONT_HERSHEY_SIMPLEX, 0.5, (0, 255, 0), 2)

            out.write(frame)

        cap.release()
        out.release()
        print(f"Processed video saved at: {outputPath}")
        return outputPath, centerPoints

    def S3VideoDetection(self, key):
        video = self.s3_manager.download_file(key)
        self.videoDetection(video)
        output, center_points = self.analyzeWithHeatmapAndBoundingBoxes(video)
        output_key = self.s3_manager.get_output_path(key)
        os.remove(video)
        self.s3_manager.upload_file(output_key, output)
        return output_key, center_points
