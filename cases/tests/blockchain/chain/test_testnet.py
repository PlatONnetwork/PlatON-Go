import pytest
import allure
from copy import copy
import time
from common.log import log
from tests.lib import Genesis
from dacite import from_dict

def test_testnet_fast(global_test_env):
    test_node = copy(global_test_env.get_a_normal_node())
    test_node.clean()
    new_cfg = copy(global_test_env.cfg)
    new_cfg.init_chain = False
    new_cfg.append_cmd = "--testnet"
    new_cfg.syncmode = "fast"
    test_node.cfg = new_cfg
    log.info("start deploy {}".format(test_node.node_mark))
    log.info("is init:{}".format(test_node.cfg.init_chain))
    test_node.deploy_me(genesis_file=None)
    log.info("deploy end")
    time.sleep(5)
    assert test_node.web3.net.peerCount > 10
    time.sleep(10)
    t = 1000
    while t:
        print(test_node.block_number)
        time.sleep(10)
        t -= 10
    assert test_node.block_number >= 10000

def test_testnet_full(global_test_env):
    test_node = copy(global_test_env.get_a_normal_node())
    test_node.clean()
    new_cfg = copy(global_test_env.cfg)
    new_cfg.init_chain = False
    new_cfg.append_cmd = "--testnet"
    test_node.cfg = new_cfg
    log.info("start deploy {}".format(test_node.node_mark))
    log.info("is init:{}".format(test_node.cfg.init_chain))
    test_node.deploy_me(genesis_file=None)
    log.info("deploy end")
    time.sleep(5)
    assert test_node.web3.net.peerCount >= 1
    time.sleep(10)
    t = 18000
    while t:
        print(test_node.block_number)
        time.sleep(10)
        t -= 10
    assert test_node.block_number >= 1000