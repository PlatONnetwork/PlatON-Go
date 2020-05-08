from concurrent.futures.process import ProcessPoolExecutor
from concurrent.futures.thread import ThreadPoolExecutor

from common.log import log

_global_dict = {}


def initGlobal():
    global _global_dict
    _global_dict = {}
    _global_dict["threadPoolExecutor"] = ThreadPoolExecutor(max_workers=40)


def set_value(name, value):
    _global_dict[name] = value


def get_value(name, defValue=None):
    try:
        return _global_dict[name]
    except KeyError:
        return defValue


def getThreadPoolExecutor(defValue=None):
    try:
        return _global_dict["threadPoolExecutor"]
    except KeyError:
        return defValue


def default_thread_pool_callback(worker):
    worker_exception = worker.exception()
    if worker_exception:
        log.exception("Thread return exception: {}".format(worker_exception))


initGlobal()
