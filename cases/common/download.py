import os
import shutil
import tarfile
import requests
from conf.settings import PLATON_BIN_FILE


def download_platon(download_url: 'str', path=PLATON_BIN_FILE):
    """
    :param download_url: new package download address
    :param path: platon relative path
    :return:
    """

    packge_name = download_url.split('/')[-1][:-7]
    platon_path = os.path.abspath(path)
    platon_path = os.path.join(platon_path, "../")
    platon_path = os.path.abspath(platon_path)
    platon_tar_file = os.path.join(platon_path, 'platon.tar.gz')
    extractall_path = os.path.join(platon_path, packge_name)

    if not os.path.exists(platon_path):
        os.makedirs(platon_path)

    # download
    resp = requests.get(url=download_url)
    with open(platon_tar_file, 'wb') as f:
        f.write(resp.content)
        f.close()

    # Extract files
    tar = tarfile.open(platon_tar_file)
    tar.extractall(path=platon_path)
    tar.close()

    # copy file
    shutil.copy(os.path.join(extractall_path, 'platon'), platon_path)

    # remove directory and file
    for root, dirs, files in os.walk(extractall_path, topdown=False):
        for name in files:
            os.remove(os.path.join(root, name))
        for name in dirs:
            os.rmdir(os.path.join(root, name))
    os.rmdir(extractall_path)
    os.remove(platon_tar_file)
