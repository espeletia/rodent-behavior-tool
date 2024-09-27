from detection.detection import Detection
from downloader.download import downloader
from handlers.app import Routes
from file_manager.file_manager import FileManager


if __name__ == "__main__":
    s3manager = FileManager('http://minio:9000',
                            'minio123', 'minio123', 'test')
    detect = Detection('./models/best.pt')
    download = downloader()
    router = Routes(detect, download, s3manager)
    router.run()
