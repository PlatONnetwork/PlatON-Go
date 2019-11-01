import time

import allure
import pytest
from copy import copy
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


@allure.title("正常启动所有节点")
@pytest.mark.P0
def test_start_all_node(global_test_env):
    """
    用例id：50
    用于测试启动所有共识节点后，检查出块情况
    """
    log.info("部署{}节点".format(len(global_test_env.consensus_node_config_list)))
    global_test_env.deploy_all()
    global_test_env.check_block()


@allure.title("启动共识节点2f+1开始出块")
@pytest.mark.P0
def test_start_mini_node(global_test_env):
    """
    用例id:51
    测试启动共识节点达到最低共识节点数量时，开始出块
    """
    num = int(2 * global_test_env.max_byzantium + 1)
    log.info("部署{}个节点".format(num))
    global_test_env.deploy_nodes(global_test_env.consensus_node_list[0:num], genesis_file=global_test_env.cfg.genesis_tmp)
    global_test_env.check_block(multiple=num, node_list=global_test_env.consensus_node_list[0:num])