import logging
import os
import time
from logging import handlers

from conf.settings import BASE_DIR, RUN_LOG_LEVEL


class Logger(object):
    level_relations = {
        'debug': logging.DEBUG,
        'info': logging.INFO,
        'warning': logging.WARNING,
        'error': logging.ERROR,
        'crit': logging.CRITICAL
    }

    def __init__(self, filename, level='info', fmt="[%(asctime)s]-%(filename)s[line:%(lineno)d] - %(levelname)s: %(message)s"):
        self.logger = logging.getLogger(filename)
        log_format = logging.Formatter(fmt)
        self.logger.setLevel(self.level_relations.get(level))
        sh = logging.StreamHandler()
        sh.setFormatter(log_format)
        file = handlers.WatchedFileHandler(filename, encoding='UTF-8')
        file.setFormatter(log_format)
        self.logger.addHandler(sh)
        self.logger.addHandler(file)


def setup_logger(logfile, loglevel):
    log = Logger(logfile, level=loglevel)
    return log


if not os.path.exists('{}/log'.format(BASE_DIR)):
    os.makedirs('{}/log'.format(BASE_DIR))
log = setup_logger('{}/log/{}.log'.format(BASE_DIR, str(
    time.strftime("%Y-%m-%d", time.localtime()))), RUN_LOG_LEVEL).logger
