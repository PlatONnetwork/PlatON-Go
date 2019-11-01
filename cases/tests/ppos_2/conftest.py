# -*- coding: utf-8 -*-
"""
@Author: wuyiqin
@Date: 2019/10/30 16:12
@Description:
"""

import pytest
from tests.lib import Delegate
from tests.lib import Staking
from tests.lib import StakingConfig
from tests.lib.client import Client
from copy import copy
from tests.lib.utils import *



"""获取所有节点对象"""
@pytest.fixture()
def node_list_obj(global_test_env):
    node_list_obj = global_test_env.get_all_nodes()
    return node_list_obj


"""获取非共识节点的质押对象"""
@pytest.fixture()
def staking_normal_obj(global_test_env):
    cfg = StakingConfig("11111", "faker", "www.baidu.com", "how much")
    node = global_test_env.get_a_normal_node()
    return Staking(global_test_env, node, cfg)


"""获取共识节点的质押对象"""
@pytest.fixture()
def staking_consensus_obj(global_test_env):
    cfg = StakingConfig("11111", "faker", "www.baidu.com", "how much")
    node = global_test_env.get_rand_node()
    return Staking(global_test_env, node, cfg)



"""获取非共识节点的委托、赎回对象"""
@pytest.fixture()
def delegate_obj(global_test_env):
    node = global_test_env.get_rand_node()
    return Delegate(global_test_env,node)


"""获取共识节点的ppos对象"""
@pytest.fixture()
def ppos_consensus_obj(global_test_env):
    cfg = StakingConfig("11111", "faker", "www.baidu.com", "how much")
    node = global_test_env.get_rand_node()
    return Client(global_test_env, node, cfg)


"""获取非共识节点的ppos对象"""
@pytest.fixture()
def ppos_consensus_obj(global_test_env):
    cfg = StakingConfig("11111", "faker", "www.baidu.com", "how much")
    node = global_test_env.get_a_normal_node()
    return Client(global_test_env, node, cfg)


"""获取所有节点的ppos对象组成列表"""
@pytest.fixture()
def ppos_obj_list(global_test_env, node_list_obj):
    cfg = StakingConfig("11111", "faker", "www.baidu.com", "how much")
    ppos_obj_list = []
    for node in node_list_obj:
        ppos_obj_list.append(Client(global_test_env, node, cfg))
    return ppos_obj_list


"""获取非共识节点未质押的节点"""
@pytest.fixture()
def get_node_nostaking(global_test_env, ppos_obj):
    normal_node_obj_list = global_test_env.normal_node_list
    """查询实时验证人列表"""
    candidate_list = get_pledge_list(ppos_obj.node.ppos.getCandidateList)
    log.info('candidatelist{}'.format(candidate_list))
    for noconsensus_node_obj in normal_node_obj_list:
        if noconsensus_node_obj.node_id not in candidate_list:
            return get_pip_obj(noconsensus_node_obj.node_id, ppos_obj)
    log.info('非共识节点已全部质押，重新启链')
    global_test_env.deploy_all()
    return get_pip_obj(normal_node_obj_list[0].node_id, ppos_obj)


"""获取一个新的global_test_env"""
@pytest.fixture()
def new_env(global_test_env):
    cfg_copy = copy(global_test_env.cfg)
    yield cfg_copy
    # global_test_env.set_cfg(cfg_copy)
    # cfg_copy.deploy_all()


"""获取新的钱包和私钥"""
@pytest.fixture()
def get_generate_account(staking_normal_obj):
    account = staking_normal_obj.economic.account
    node = staking_normal_obj.node
    address, prikey = account.generate_account(node.web3, 10 ** 18 * 10000000)
    return address, prikey
