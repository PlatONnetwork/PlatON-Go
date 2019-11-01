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


@allure.title("正常启动所有节点,逐渐关闭f+1个")
@pytest.mark.P0
def test_start_all_node_close_f_add_1(global_test_env):
    """
    启动所有节点后，逐渐关闭f+1个，那么关闭后将不会出块
    """
    global_test_env.deploy_all()
    test_nodes = global_test_env.consensus_node_list[:global_test_env.max_byzantium + 1]
    global_test_env.stop_nodes(test_nodes)
    running_node = global_test_env.consensus_node_list[global_test_env.max_byzantium + 1:]
    time.sleep(5)
    start = max(global_test_env.block_numbers(running_node).values())
    time.sleep(5)
    end = max(global_test_env.block_numbers(running_node).values())
    assert start == end


@allure.title("先启动2f个节点，间隔{t}秒后再启动一个")
@pytest.mark.P2
@pytest.mark.parametrize('t', [50, 150])
def test_start_2f_after_one(t, global_test_env):
    """
    先启动2f个节点，间隔一定时间之后再启动一个节点，查看出块情况
    """
    num = int(2 * global_test_env.max_byzantium)
    log.info("先启动{}个节点".format(num))
    test_nodes = global_test_env.consensus_node_list[0:num + 1]
    global_test_env.deploy_nodes(global_test_env.consensus_node_list[0:num], genesis_file=global_test_env.cfg.genesis_tmp)
    time.sleep(t)
    global_test_env.deploy_nodes(global_test_env.consensus_node_list[num:num + 1], genesis_file=global_test_env.cfg.genesis_tmp)
    global_test_env.check_block(node_list=test_nodes, multiple=num)


@allure.title("先启动2f个节点，间隔{t}秒后启动所有节点")
@pytest.mark.P2
@pytest.mark.parametrize('t', [50, 150])
def test_start_2f_after_all(t, global_test_env):
    """
    先启动2f个节点，间隔一定时间之后再启动未启动的所有共识节点，查看出块情况
    """
    num = int(2 * global_test_env.max_byzantium)
    log.info("先启动{}个节点".format(num))
    global_test_env.deploy_nodes(global_test_env.consensus_node_list[0:num], genesis_file=global_test_env.cfg.genesis_tmp)
    time.sleep(int(t))
    log.info("在启动另外所有共识节点")
    global_test_env.deploy_nodes(global_test_env.consensus_node_list[num:], genesis_file=global_test_env.cfg.genesis_tmp)
    global_test_env.check_block(node_list=global_test_env.consensus_node_list, multiple=5)


@allure.title("先启动2f个节点，30秒内不停重启另外节点")
@pytest.mark.P2
def test_up2f_after_other(global_test_env):
    """
    用例id:61,62
    """
    num = int(2 * global_test_env.max_byzantium)
    log.info("先启动{}个节点".format(num))
    global_test_env.deploy_nodes(global_test_env.consensus_node_list[0:num], genesis_file=global_test_env.cfg.genesis_tmp)
    i = 0
    global_test_env.deploy_nodes(global_test_env.consensus_node_list[num:], genesis_file=global_test_env.cfg.genesis_tmp)
    while i <= 30:
        global_test_env.reset_nodes(global_test_env.consensus_node_list[num:])
        i += 1
    global_test_env.check_block(node_list=global_test_env.consensus_node_list)


@allure.title("正常启动所有节点,逐渐关闭f+1个,再逐步启动所有")
@pytest.mark.P0
def test_start_all_node_close_f_add_1_and_all(global_test_env):
    """
    用例id
    启动所有节点后，逐渐关闭f+1个，那么关闭后将不会出块,等待重启后出块
    """
    global_test_env.deploy_all()
    test_nodes = global_test_env.consensus_node_list[0:global_test_env.max_byzantium + 1]
    global_test_env.stop_nodes(test_nodes)
    global_test_env.start_nodes(test_nodes, False)
    global_test_env.check_block(need_number=20, multiple=2)


@allure.title("正常启动所有节点,逐渐关闭f+1个,再逐步启动一个")
@pytest.mark.P0
def test_start_all_node_close_f_add_1_and_one(global_test_env):
    """
    启动所有节点后，逐渐关闭f+1个，那么关闭后将不会出块,等待重启后出块
    """
    global_test_env.deploy_all()
    test_nodes = global_test_env.consensus_node_list[0:global_test_env.max_byzantium + 1]
    global_test_env.stop_nodes(test_nodes)
    global_test_env.start_nodes(test_nodes[global_test_env.max_byzantium:global_test_env.max_byzantium + 1], False)
    global_test_env.check_block(multiple=5, node_list=global_test_env.consensus_node_list[global_test_env.max_byzantium:])


@allure.title("正常启动所有节点,等待出块一段时间后，关闭一个，并删除数据库，用fast模式启动")
@pytest.mark.P0
def test_start_all_node_close_f_add_1_and_fast_one(global_test_env):
    """
    用例id
    正常启动所有节点,等待出块一段时间后，关闭一个，并删除数据库，用fast模式启动
    """
    global_test_env.deploy_all()
    time.sleep(100)
    test_node = copy(global_test_env.get_rand_node())
    test_node.clean()
    new_cfg = copy(global_test_env.cfg)
    new_cfg.syncmode = "fast"
    test_node.cfg = new_cfg
    test_node.deploy_me(genesis_file=new_cfg.genesis_tmp)
    time.sleep(10)
    assert test_node.block_number > 10
