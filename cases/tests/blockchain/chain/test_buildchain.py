import pytest
import allure
from copy import copy
import time
from common.log import log
from tests.lib import Genesis
from dacite import from_dict


@pytest.fixture(scope="function", autouse=True)
def stop_env(global_test_env):
    global_test_env.stop_all()


@pytest.fixture(scope="module", autouse=True)
def reset_env(global_test_env):
    cfg = copy(global_test_env.cfg)
    yield
    log.info("reset deploy.................")
    global_test_env.cfg = cfg
    global_test_env.deploy_all()


@allure.title("Node interconnects that do not initialize boot nodes and different creation files")
@pytest.mark.P0
def test_CH_IN_006(global_test_env):
    pass
    # global_test_env.deploy_all()
    # test_node = copy(global_test_env.get_a_normal_node())
    # test_node.clean()
    # new_cfg = copy(global_test_env.cfg)
    # new_cfg.init_chain = False
    # new_cfg.append_cmd = "--nodiscover"
    # test_node.cfg = new_cfg
    # log.info("start deploy {}".format(test_node.node_mark))
    # log.info("is init:{}".format(test_node.cfg.init_chain))
    # test_node.deploy_me(genesis_file=None)
    # log.info("deploy end")
    # test_node.admin.addPeer(global_test_env.get_rand_node().enode)
    # time.sleep(5)
    # assert test_node.web3.net.peerCount == 0, "the number of connected nodes has increased"


@allure.title("Test deployment of a single-node private chain, synchronization of single-node blocks")
@pytest.mark.P0
def test_CH_IN_005(global_test_env):
    test_node = copy(global_test_env.get_a_normal_node())
    log.info("test node :{}".format(test_node.node_mark))
    genesis_data = global_test_env.genesis_config
    genesis = from_dict(data_class=Genesis, data=genesis_data)
    genesis.config.cbft.initialNodes = [{"node": test_node.enode, "blsPubKey": test_node.blspubkey}]
    file = test_node.local_node_tmp + "/genesis_0.13.0.json"
    genesis.to_file(file)
    test_node.deploy_me(file)
    time.sleep(5)
    assert test_node.block_number > 0, "block height has not increased"
    join_node = copy(global_test_env.get_rand_node())
    log.info("join node:{}".format(join_node.node_mark))
    join_node.deploy_me(file)
    join_node.admin.addPeer(test_node.enode)
    time.sleep(5)
    assert join_node.block_number > 0, "block height has not increased"


@allure.title("Test node interconnections between different initnode founding files")
@pytest.mark.P0
def test_CH_IN_009(global_test_env):
    global_test_env.deploy_all()
    test_node = copy(global_test_env.get_a_normal_node())
    log.info("test node :{}".format(test_node.node_mark))
    genesis = from_dict(data_class=Genesis, data=global_test_env.genesis_config)
    genesis.config.cbft.initialNodes = [{"node": test_node.enode, "blsPubKey": test_node.blspubkey}]
    file = test_node.local_node_tmp + "/genesis_0.13.0.json"
    genesis.to_file(file)
    test_node.deploy_me(file)
    log.info(test_node.running)
    test_node.admin.addPeer(global_test_env.get_rand_node().enode)
    time.sleep(5)
    assert test_node.web3.net.peerCount == 0
