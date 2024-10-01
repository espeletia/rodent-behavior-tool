import logging

from pythonjsonlogger import jsonlogger

logger = logging.getLogger()

logHandler = logging.StreamHandler()
formatter = jsonlogger.JsonFormatter()
logHandler.setFormatter(formatter)
logger.addHandler(logHandler)
logger.setLevel("INFO")


def log_message(msg):
    logger.info(msg)


def log_error(err):
    logger.error({"err": err})
