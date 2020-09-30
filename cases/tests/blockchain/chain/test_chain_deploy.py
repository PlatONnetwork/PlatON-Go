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


@allure.title("start all nodes normally")
@pytest.mark.P0
def test_SC_ST_001(global_test_env):
    """
    Used to test the start of all consensus nodes, check out the block situation
    """
    log.info("Deploy {} node".format(len(global_test_env.consensus_node_config_list)))
    global_test_env.deploy_all()
    global_test_env.check_block()


@allure.title("Start consensus node 2f+1 starts to block")
@pytest.mark.P0
def test_SC_ST_002(global_test_env):
    """
    When the test start consensus node reaches the minimum consensus node number, it starts to pop out.
    """
    num = int(2 * global_test_env.max_byzantium + 1)
    log.info("Deploy {} nodes".format(num))
    global_test_env.deploy_nodes(global_test_env.consensus_node_list[0:num], genesis_file=global_test_env.cfg.genesis_tmp)
    global_test_env.check_block(multiple=num, node_list=global_test_env.consensus_node_list[0:num])


@allure.title("Start all nodes normally, and gradually close f")
@pytest.mark.P0
def test_SC_CL_001(global_test_env):
    """
    After starting n nodes, gradually close f, then the window of the closed node does not come out.
    """
    global_test_env.deploy_all()
    global_test_env.check_block()
    close_nodes = global_test_env.get_all_nodes()[0:global_test_env.max_byzantium]
    global_test_env.stop_nodes(close_nodes)
    global_test_env.check_block(need_number=30, multiple=2,
                                node_list=global_test_env.get_all_nodes()[global_test_env.max_byzantium:])


@allure.title("Start 2f+1 nodes normally, start one after 50 seconds")
@pytest.mark.P2
def test_SC_IV_001(global_test_env):
    """
    Start 2f+1 first, start one after 50 seconds
    """
    num = int(2 * global_test_env.max_byzantium + 1)
    log.info("Deploy {} nodes".format(num))
    test_nodes = global_test_env.consensus_node_list[0:num]
    global_test_env.deploy_nodes(node_list=test_nodes, genesis_file=global_test_env.cfg.genesis_tmp)
    time.sleep(50)
    start = max(global_test_env.block_numbers(node_list=test_nodes).values())
    global_test_env.deploy_nodes(global_test_env.consensus_node_list[num:num + 1], genesis_file=global_test_env.cfg.genesis_tmp)
    global_test_env.check_block(need_number=start + 10, multiple=2, node_list=global_test_env.consensus_node_list[0:num + 1])


@allure.title("Only start 2f nodes")
@pytest.mark.P0
def test_SC_ST_003(global_test_env):
    """
    Start 2f nodes
    """
    num = int(2 * global_test_env.max_byzantium)
    log.info("Deploy {} nodes".format(num))
    test_nodes = global_test_env.consensus_node_list[0:num]
    global_test_env.deploy_nodes(test_nodes, genesis_file=global_test_env.cfg.genesis_tmp)
    time.sleep(10)
    block = max([node.block_number for node in test_nodes])
    assert block == 0


@allure.title("Start all nodes normally and gradually close f+1")
@pytest.mark.P0
def test_SC_CL_002(global_test_env):
    """
    After starting all nodes, gradually close f+1, then it will not be blocked after closing.
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


@allure.title("Start 2f nodes first, then start one after {t} seconds")
@pytest.mark.P2
@pytest.mark.parametrize('t', [50, 150])
def test_SC_IV_002_to_003(t, global_test_env):
    """
    Start 2f nodes first, then start a node after a certain interval to see the block situation.
    """
    num = int(2 * global_test_env.max_byzantium)
    log.info("Start {} nodes first".format(num))
    test_nodes = global_test_env.consensus_node_list[0:num + 1]
    global_test_env.deploy_nodes(global_test_env.consensus_node_list[0:num], genesis_file=global_test_env.cfg.genesis_tmp)
    time.sleep(t)
    global_test_env.deploy_nodes(global_test_env.consensus_node_list[num:num + 1], genesis_file=global_test_env.cfg.genesis_tmp)
    global_test_env.check_block(node_list=test_nodes, multiple=num)


@allure.title("Start 2f nodes first, start all nodes after {t} seconds")
@pytest.mark.P2
@pytest.mark.parametrize('t', [50, 150])
def test_SC_IV_004_to_005(t, global_test_env):
    """
    Start 2f nodes first, and then start all the consensus nodes that are not
    started after a certain interval, and check the block status.
    """
    num = int(2 * global_test_env.max_byzantium)
    log.info("Start {} nodes first".format(num))
    global_test_env.deploy_nodes(global_test_env.consensus_node_list[0:num], genesis_file=global_test_env.cfg.genesis_tmp)
    time.sleep(int(t))
    log.info("Start all other consensus nodes")
    global_test_env.deploy_nodes(global_test_env.consensus_node_list[num:], genesis_file=global_test_env.cfg.genesis_tmp)
    global_test_env.check_block(node_list=global_test_env.consensus_node_list, multiple=5)


@allure.title("Start 2f nodes first, and restart other nodes within 30 seconds")
@pytest.mark.P2
def test_SC_IR_001(global_test_env):
    num = int(2 * global_test_env.max_byzantium)
    log.info("Start {} nodes first".format(num))
    global_test_env.deploy_nodes(global_test_env.consensus_node_list[0:num], genesis_file=global_test_env.cfg.genesis_tmp)
    i = 0
    global_test_env.deploy_nodes(global_test_env.consensus_node_list[num:], genesis_file=global_test_env.cfg.genesis_tmp)
    while i <= 30:
        global_test_env.reset_nodes(global_test_env.consensus_node_list[num:])
        i += 1
    global_test_env.check_block(node_list=global_test_env.consensus_node_list)


@allure.title("Start all nodes normally, gradually close f+1, and then start all gradually")
@pytest.mark.P0
def test_SC_IR_002(global_test_env):
    """
    After starting all nodes, gradually close f+1, then there will
    be no block after closing, waiting for the block to be restarted.
    """
    global_test_env.deploy_all()
    test_nodes = global_test_env.consensus_node_list[0:global_test_env.max_byzantium + 1]
    global_test_env.stop_nodes(test_nodes)
    global_test_env.start_nodes(test_nodes, False)
    global_test_env.check_block(need_number=20, multiple=2)


@allure.title("Start all nodes normally, gradually close f+1, and then start one step by step.")
@pytest.mark.P0
def test_SC_RC_001(global_test_env):
    """
    After starting all nodes, gradually close f+1, then there will be no
     block after closing, waiting for the block to be restarted.
    """
    global_test_env.deploy_all()
    test_nodes = global_test_env.consensus_node_list[0:global_test_env.max_byzantium + 1]
    global_test_env.stop_nodes(test_nodes)
    global_test_env.start_nodes(test_nodes[global_test_env.max_byzantium:global_test_env.max_byzantium + 1], False)
    global_test_env.check_block(multiple=5, node_list=global_test_env.consensus_node_list[global_test_env.max_byzantium:])


@allure.title("Start all nodes normally, wait for a block of time, close one, and delete the database, start with fast mode")
@pytest.mark.P0
def test_SC_FT_001(global_test_env):
    """
    Start all nodes normally, wait for a block of time, close one, and delete the database, start with fast mode
    """
    global_test_env.deploy_all()
    time.sleep(100)
    test_node = copy(global_test_env.get_rand_node())
    test_node.stop()
    test_node.clean()
    test_node.run_ssh("cd {};ls".format(test_node.remote_node_path))
    new_cfg = copy(global_test_env.cfg)
    new_cfg.syncmode = "fast"
    test_node.cfg = new_cfg
    is_success, msg = test_node.deploy_me(genesis_file=new_cfg.genesis_tmp)
    if not is_success:
        raise Exception(msg)
    time.sleep(100)
    assert test_node.block_number > 100
