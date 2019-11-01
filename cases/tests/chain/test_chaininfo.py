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


@allure.title("区块连续性，验证hash")
def test_hash_continuous(global_test_env):
    """
    测试区块的连续性，验证一定数量的区块，区块哈希必须是连续的
    """
    global_test_env.check_block(100, 2)
    node = global_test_env.get_rand_node()
    block_hash = HexBytes(node.eth.getBlock(1).get("hash")).hex()
    for i in range(2, 100):
        block = node.eth.getBlock(i)
        parent_hash = HexBytes(block.get("parentHash")).hex()
        assert block_hash == parent_hash, "父区块哈希值错误"
        block_hash = HexBytes(block.get("hash")).hex()


@allure.title("测试platon文件的版本号")
@pytest.mark.P3
def test_platon_versions(global_test_env):
    node = global_test_env.get_rand_node()
    cmd_list = node.run_ssh("{} version".format(node.remote_bin_file))
    assert global_test_env.version in cmd_list[1], "版本号不正确"


@allure.title("测试重启所有共识节点")
@pytest.mark.P0
def test_restart_all(global_test_env):
    current_block = max(global_test_env.block_numbers().values())
    log.info("重启前块高:{}".format(current_block))
    global_test_env.reset_all()
    log.info("重启所有共识节点成功")
    time.sleep(20)
    after_block = max(global_test_env.block_numbers().values())
    log.info("重启后块高为:{}".format(after_block))
    assert after_block - current_block > 0, "重启后区块没有正常增长"


@allure.title("测试fast模式同步")
@pytest.mark.P2
def test_syncmode(global_test_env):
    """
    同步
    :return:
    """
    test_node = copy(global_test_env.get_a_normal_node())
    test_node.clean()
    new_cfg = copy(global_test_env.cfg)
    new_cfg.syncmode = "fast"
    test_node.cfg = new_cfg
    log.info(global_test_env.cfg.syncmode)
    test_node.deploy_me(global_test_env.cfg.genesis_tmp)
    test_node.admin.addPeer(global_test_env.get_rand_node().enode)
    time.sleep(5)
    log.info("{}".format(test_node.web3.net.peerCount))
    assert test_node.web3.net.peerCount > 0, "加入链失败"
    global_test_env.check_block(200, 2)
    time.sleep(5)
    assert test_node.eth.blockNumber >= 200, "区块同步失败,当前块高{}".format(test_node.eth.blockNumber)


@allure.title("测试区块同步")
@pytest.mark.P0
def test_deconsensus_block_synchronize(global_test_env):
    """
    非共识节点块高同步
    :param global_test_env:
    :return:
    """
    test_node = global_test_env.get_a_normal_node()
    test_node.clean()
    new_cfg = global_test_env.cfg
    new_cfg.syncmode = "full"
    test_node.cfg = new_cfg
    test_node.deploy_me(global_test_env.cfg.genesis_tmp)
    test_node.admin.addPeer(global_test_env.get_rand_node().enode)
    time.sleep(5)
    assert test_node.web3.net.peerCount > 0, "加入链失败"
    global_test_env.check_block()
    assert test_node.block_number > 0, "非共识节点同步区块失败，块高：{}".format(test_node.block_number)