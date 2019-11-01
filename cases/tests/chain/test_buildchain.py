import pytest
import allure
from copy import copy
import time
from common.log import log


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


@allure.title("不初始化启动节点和不同创世文件的节点互连")
@pytest.mark.P0
def test_no_init_no_join_chain(global_test_env):
    global_test_env.deploy_all()
    test_node = copy(global_test_env.get_a_normal_node())
    test_node.clean()
    new_cfg = copy(global_test_env.cfg)
    new_cfg.init_chain = False
    new_cfg.append_cmd = "--nodiscover"
    test_node.cfg = new_cfg
    log.info("start deploy {}".format(test_node.node_mark))
    log.info("is init:{}".format(test_node.cfg.init_chain))
    test_node.deploy_me(genesis_file=None)
    log.info("deploy end")
    test_node.admin.addPeer(global_test_env.get_rand_node().enode)
    time.sleep(5)
    assert test_node.web3.net.peerCount == 0, "连接节点数有增加"