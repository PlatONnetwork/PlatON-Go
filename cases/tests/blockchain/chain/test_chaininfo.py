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


@allure.title("Whether the block information is consistent")
def test_CH_IN_021(global_test_env):
    """
    Test whether all node block information is consistent
    """
    nodes = global_test_env.get_all_nodes()
    rand_node = global_test_env.get_rand_node()
    block_number = min(global_test_env.block_numbers().values())
    block_info = rand_node.eth.getBlock(block_number)
    for node in nodes:
        assert block_info == node.eth.getBlock(block_number), "The same block height information of different nodes is inconsistent, block number: {}".format(
            block_number)


@allure.title("Block continuity, verify hash")
def test_CH_IN_020(global_test_env):
    """
    Test block continuity, verify a certain number of blocks, block hash must be continuous
    """
    global_test_env.check_block(100, 2)
    node = global_test_env.get_rand_node()
    block_hash = HexBytes(node.eth.getBlock(1).get("hash")).hex()
    for i in range(2, 100):
        block = node.eth.getBlock(i)
        parent_hash = HexBytes(block.get("parentHash")).hex()
        assert block_hash == parent_hash, "Parent block hash value error"
        block_hash = HexBytes(block.get("hash")).hex()


@allure.title("Test the version number of the platon file")
@pytest.mark.P3
def test_platon_versions(global_test_env):
    node = global_test_env.get_rand_node()
    cmd_list = node.run_ssh("{} version".format(node.remote_bin_file))
    assert global_test_env.version in cmd_list[1], "The version number is incorrect"


@allure.title("Test restart all consensus nodes")
@pytest.mark.P0
def test_CH_IN_019_SC_RC_003(global_test_env):
    current_block = max(global_test_env.block_numbers().values())
    log.info("Block height before restart:{}".format(current_block))
    global_test_env.reset_all()
    log.info("Restart all consensus nodes successfully")
    time.sleep(20)
    after_block = max(global_test_env.block_numbers().values())
    log.info("After restarting, the block height is: {}".format(after_block))
    assert after_block - current_block > 0, "After the restart, the block did not grow normally."


@allure.title("Test fast mode synchronization")
@pytest.mark.P2
def test_CH_IN_017(global_test_env):
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
    assert test_node.web3.net.peerCount > 0, "Joining the chain failed"
    global_test_env.check_block(200, 2)
    time.sleep(5)
    assert test_node.eth.blockNumber >= 200, "Block sync failed, current block is high{}".format(test_node.eth.blockNumber)


@allure.title("Test block synchronization")
@pytest.mark.P0
def test_CH_IN_018_CH_IN_007_CMD_039(global_test_env):
    """
    Non-consensus node block high synchronization
    """
    test_node = global_test_env.get_a_normal_node()
    test_node.clean()
    new_cfg = global_test_env.cfg
    new_cfg.syncmode = "full"
    test_node.cfg = new_cfg
    test_node.deploy_me(global_test_env.cfg.genesis_tmp)
    test_node.admin.addPeer(global_test_env.get_rand_node().enode)
    time.sleep(5)
    assert test_node.web3.net.peerCount > 0, "Joining the chain failed"
    global_test_env.check_block()
    assert test_node.block_number > 0, "Non-consensus node sync block failed, block height: {}".format(test_node.block_number)


@allure.title("Test node interconnection between identical founding files")
@pytest.mark.P0
def test_CH_IN_008_CH_IN_011(global_test_env):
    test_node = global_test_env.normal_node_list[0]
    assert test_node.web3.net.peerCount > 0, "Joining the chain failed"
