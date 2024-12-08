import sys
import time
from detection.detection import Detection
from downloader.download import downloader
from handlers.app import Routes
from file_manager.file_manager import FileManager
import configuration.config as config
import workers.analysis
import asyncio
import utils.logger as logger
from multiprocessing import Process


def main():
    running_processes = []

    s3manager = FileManager(config.S3_ENDPOINT,
                            config.S3_ACCESS_KEY,
                            config.S3_SECRET_ACCESS_KEY,
                            'test',
                            )
    detect = Detection(config.WHEIGHTS_PATH, s3manager)
    download = downloader()
    router = Routes(detect, download, s3manager)
    analysis_worker = workers.analysis.AnalystWorker(detect)

    def run_api():
        router.run()

    def run_analysis_worker():
        asyncio.run(analysis_worker.start())

    def create_analysis_worker():
        logger.log_message("Starting analysis worker")
        p = Process(target=run_analysis_worker,
                    name="analysis_worker", daemon=True)
        running_processes.append(p)
        p.start()

    run_analysis_worker()

    logger.log_message("starting RATT service!")
    api = Process(target=run_api, name="api", daemon=True)
    running_processes.append(api)
    api.start()

    while True:
        for process in running_processes:
            if process.exitcode is not None:
                logger.log_error(
                    f"Process: {process.name} EXITED with code: {process.exitcode}")
                sys.exit(1)
        if len(running_processes) == 0:
            logger.log_error("No worker is running anymore exiting")
            sys.exit(1)
        time.sleep(2)
        pass


if __name__ == "__main__":
    # asyncio.run(main())
    main()
