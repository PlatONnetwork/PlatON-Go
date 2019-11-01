import os

from conf.settings import BASE_DIR


def abspath(path):
    """
    Based on the project root directory stitching path,
    The path format is spliced ​​when the relative path is ./path or path/path2.
     When the path format is absolute path/path, the original path is directly returned.
    :param path:
    :return:
    """
    if os.path.isabs(path):
        return path
    path = path.lstrip("./")
    return os.path.abspath(BASE_DIR + "/" + path)
