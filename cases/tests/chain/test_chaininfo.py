import pytest
import allure
from environment.env import TestEnvironment
from copy import copy
from hexbytes import HexBytes
import time
from common.log import log


@pytest.fixture(scope="function", autouse=True)
def restart_env(global_test_env):
    if not global_test_env.running:
        global_test_env.deploy_all()
    global_test_env.check_block(multiple=3)


@allure.title("区块信息是否一致")
def test_block_info_synchronize(global_test_env):
    """
    测试所有节点区块信息是否一致
    :param global_test_env:
    :return:
    """
    nodes = global_test_env.get_all_nodes()
    rand_node = global_test_env.get_rand_node()
    block_number = min(global_test_env.block_numbers().values())
    block_info = rand_node.eth.getBlock(block_number)
    for node in nodes:
        assert block_info == node.eth.getBlock(block_number), "不同节点的相同块高信息不一致,区块号：{}".format(
            block_number)