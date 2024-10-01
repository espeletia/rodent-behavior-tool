import signal

from utils import logger


class GracefulTermination:
    terminate = False

    def __init__(self):
        signal.signal(signal.SIGINT, self.exit_gracefully)
        signal.signal(signal.SIGTERM, self.exit_gracefully)
        signal.signal(signal.SIGABRT, self.exit_gracefully)

    def exit_gracefully(self, *args):
        self.terminate = True
        logger.log_error("Gracefully shutting down worker")
