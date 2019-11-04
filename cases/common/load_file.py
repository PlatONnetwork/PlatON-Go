import json
import os
import hashlib
import yaml


class LoadFile(object):
    """
    Convert json or yaml files to python dictionary or list dictionary
    """

    def __init__(self, file):
        if file.split('.')[-1] != 'yaml' and file.split('.')[-1] != 'json' and file.split('.')[-1] != 'yml':
            raise Exception("the file format must be yaml or json")
        self.file = file

    def get_data(self):
        """
        call this method to get the result
        """
        if self.file.split('.')[-1] == "json":
            return self.load_json()
        return self.load_yaml()

    def load_json(self):
        """
        Convert json file to dictionary
        """
        try:
            with open(self.file, encoding="utf-8") as f:
                result = json.load(f)
                if isinstance(result, list):
                    result = [i for i in result if i != '']
                return result
        except FileNotFoundError as e:
            raise e

    def load_yaml(self):
        """
        Convert yaml file to dictionary
        """
        try:
            with open(self.file, encoding="utf-8")as f:
                result = yaml.load(f)
                if isinstance(result, list):
                    result = [i for i in result if i != '']
                return result
        except FileNotFoundError as e:
            raise e


def get_all_file(path):
    """
    Get all yaml or json files
    :param path:
    :return:
    """
    try:
        result = [os.path.abspath(os.path.join(path, filename)) for filename in os.listdir(
            path) if filename.endswith(".json") or filename.endswith(".yml") or filename.endswith(".yaml")]
        return result
    except FileNotFoundError as e:
        raise e


def get_file(path):
    """
    Get all yaml or json files
    :param path:
    :return:
    """
    try:
        result = []
        for x, _, _ in os.walk(path):
            if os.listdir(x):
                result += get_all_file(x)
            else:
                result += x
        return result
    except FileNotFoundError as e:
        raise e


def get_f(collsion_list):
    """
    get the maximum number of cheat nodes
    :param collsion_list:
    :return:
    """
    num = len(collsion_list)
    if num < 3:
        raise Exception("the number of consensus nodes is less than 3")
    if num == 3:
        return 0
    f = (num - 1) / 3
    return int(f)


def get_f_for_n(n):
    """
    Get the maximum number of cheat nodes based on the total number of consensus nodes
    :param n:
    :return:
    """
    num = n
    if num < 3:
        raise Exception("the number of consensus nodes is less than 3")
    if num == 3:
        return 0
    f = (num - 1) / 3
    return int(f)


def get_file_time(file) -> int:
    if not os.path.exists(file):
        raise Exception("file is not exists")
    return int(os.path.getctime(file))


def calc_hash(file):
    with open(file, 'rb') as f:
        sha1obj = hashlib.sha1()
        sha1obj.update(f.read())
        result_hash = sha1obj.hexdigest()
        return result_hash


if __name__ == "__main__":
    calc_hash(r"D:\python\PlatON-Tests\deploy\bin\platon")
