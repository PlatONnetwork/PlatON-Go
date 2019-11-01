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


@allure.title("正常启动所有节点,逐渐关闭f个")
@pytest.mark.P0
def test_start_all_node_close_f(global_test_env):
    """
    用例id：52
    启动n个节点后，逐渐关闭f个，那么关闭节点的窗口期不出块
    """
    global_test_env.deploy_all()
    global_test_env.check_block()
    close_nodes = global_test_env.get_all_nodes()[0:global_test_env.max_byzantium]
    global_test_env.stop_nodes(close_nodes)
    global_test_env.check_block(need_number=30, multiple=2,
                                node_list=global_test_env.get_all_nodes()[global_test_env.max_byzantium:])


@allure.title("正常启动2f+1个节点,50秒后在启动一个")
@pytest.mark.P2
def test_start_2f1_node_and_start_one(global_test_env):
    """
    先启动2f+1个，50秒后在启动一个
    """
    num = int(2 * global_test_env.max_byzantium + 1)
    log.info("部署{}个节点".format(num))
    test_nodes = global_test_env.consensus_node_list[0:num]
    global_test_env.deploy_nodes(node_list=test_nodes, genesis_file=global_test_env.cfg.genesis_tmp)
    time.sleep(50)
    start = max(global_test_env.block_numbers(node_list=test_nodes).values())
    global_test_env.deploy_nodes(global_test_env.consensus_node_list[num:num + 1], genesis_file=global_test_env.cfg.genesis_tmp)
    global_test_env.check_block(need_number=start + 10, multiple=2, node_list=global_test_env.consensus_node_list[0:num + 1])


@allure.title("只启动2f个节点")
@pytest.mark.P0
def test_start_2f(global_test_env):
    """
    启动2f个节点
    """
    num = int(2 * global_test_env.max_byzantium)
    log.info("部署{}个节点".format(num))
    test_nodes = global_test_env.consensus_node_list[0:num]
    global_test_env.deploy_nodes(test_nodes, genesis_file=global_test_env.cfg.genesis_tmp)
    time.sleep(10)
    block = max([node.block_number for node in test_nodes])
    assert block == 0
