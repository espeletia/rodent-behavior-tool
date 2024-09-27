from flask import Flask, request, jsonify


class Routes:

    def __init__(self, detector, downloader, file_manager):
        self.app = Flask(__name__)
        self.detector = detector
        self.downloader = downloader
        self.file_manager = file_manager
        self.setup_routes()

    def setup_routes(self):
        @self.app.route('/', methods=['GET'])
        def home():
            return "Hello world"

        @self.app.route('/detect', methods=['POST'])
        def detect():
            data = request.get_json()

            if not data:
                return jsonify({"error": "Invalid or missing JSON"}), 400

            url = data.get('url')

            video = self.downloader.downloadVideoFromUrl(url)

            self.detector.videoDetection(video)
            locations = self.detector.analyzeResults()

            return jsonify(
                {
                    "message": f"you've sent me this url: {url}",
                    "locations": locations
                }
            ), 200

        @self.app.route('/detect-s3', methods=['POST'])
        def detectFromS3():
            data = request.get_json()
            if not data:
                return jsonify({"error": "Invalid or missing JSON"})
            key = data.get('key')
            video = self.file_manager.download_file(key)
            self.detector.videoDetection(video)
            locations = self.detector.analyzeResults()

            return jsonify(
                {
                    "message": f"you've sent me this url: {key}",
                    "locations": locations
                }
            ), 200

    def run(self):
        self.app.run(host='0.0.0.0', port=8080, debug=True)
