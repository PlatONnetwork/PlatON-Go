# -*- coding: utf-8 -*-
"""
@Author: wuyiqin
@Date: 2019/10/30 11:55
@Description:
"""

from common.log import log
from tests.lib.utils import *


def test_staking(staking_obj):
    account = staking_obj.economic.account
    node = staking_obj.node
    address, _ = account.generate_account(node.web3, 10**18 * 10000000)
    log.info("Generate address:{}".format(address))
    result = staking_obj.create_staking(0, address, address)
    log.info("Staking result:{}".format(result))
    assert result["Code"] == 0
    assert result["ErrMsg"] == "ok"


def test_IV_001_002(global_test_env,staking_consensus_obj):
    node_info = staking_consensus_obj.ppos.getValidatorList()
    log.info(node_info)
    nodeid_list = []
    for node in node_info.get("Data"):
        nodeid_list.append(node.get("NodeId"))
        StakingAddress = node.get("StakingAddress")
        log.info(StakingAddress)
        assert staking_consensus_obj.node.web3.toChecksumAddress(StakingAddress) == \
               staking_consensus_obj.economic.cfg.DEVELOPER_FOUNDATAION_ADDRESS
    log.info(nodeid_list)
    consensus_node_list = global_test_env.consensus_node_list
    nodeid_list_ = [node.node_id for node in consensus_node_list]
    log.info(nodeid_list_)
    for nodeid in nodeid_list_:
        assert nodeid in nodeid_list


def test_IV_003(staking_consensus_obj):
    StakingAddress = staking_consensus_obj.economic.cfg.DEVELOPER_FOUNDATAION_ADDRESS
    result = staking_consensus_obj.create_staking(0, StakingAddress, StakingAddress)
    log.info("Staking result:{}".format(result))
    assert result["Code"] == 301101


def test_IV_004(get_generate_account,delegate_obj):
    address, _ = get_generate_account
    msg = delegate_obj.delegate(0, address)
    log.info(msg)
    assert msg["Code"] == 301107


def test_IV_005(staking_consensus_obj):
    StakingAddress = staking_consensus_obj.economic.cfg.DEVELOPER_FOUNDATAION_ADDRESS
    msg = staking_consensus_obj.increase_staking(0,StakingAddress)
    assert msg["Code"] == 0


def test_IV_006(staking_consensus_obj):
    StakingAddress = staking_consensus_obj.economic.cfg.DEVELOPER_FOUNDATAION_ADDRESS
    msg = staking_consensus_obj.withdrew_staking(StakingAddress)
    log.info(msg)
    assert msg["Code"] == 0














