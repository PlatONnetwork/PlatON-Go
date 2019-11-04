# -*- coding: utf-8 -*-
from tests.lib.utils import *
import pytest


def test_AS_031_032_39(client_new_node_obj, get_generate_account):
    """
    Normal overweight
    The verifier initiates the overweight with the amount of free account, meeting the minimum threshold
    :param client_new_node_obj:
    :param get_generate_account:
    :return:
    """
    address, _ = get_generate_account
    result = client_new_node_obj.staking.create_staking(0, address, address)
    assert result.get('Code') == 0
    result = client_new_node_obj.staking.increase_staking(0, address)
    assert result.get('Code') == 0
    result = client_new_node_obj.ppos.getCandidateInfo(client_new_node_obj.node.node_id)
    log.info(result)
    staking_amount = client_new_node_obj.economic.create_staking_limit
    add_staking_amount = client_new_node_obj.economic.add_staking_limit
    assert result["Data"]["Shares"] == staking_amount + add_staking_amount


def test_AS_033(client_new_node_obj, get_generate_account):
    """
    The verifier is not on the verifier and candidate list
    :param client_new_node_obj:
    :param get_generate_account:
    :return:
    """
    address, _ = get_generate_account
    result = client_new_node_obj.staking.increase_staking(0, address)
    log.info(result)
    assert result.get('Code') == 301102


def test_AS_034(client_new_node_obj, get_generate_account):
    """
    Undersupply of gas
    :param client_new_node_obj:
    :param get_generate_account:
    :return:
    """
    address, _ = get_generate_account
    client_new_node_obj.staking.create_staking(0, address, address)
    fig = {"gas": 1}
    status = 0
    try:
        result = client_new_node_obj.staking.increase_staking(0, address, transaction_cfg=fig)
        log.info(result)
    except:
        status = 1
    assert status == 1


def test_AS_035(client_new_node_obj):
    """
    Insufficient balance initiated overweight
    :param client_new_node_obj:
    :return:
    """
    account = client_new_node_obj.economic.account
    node = client_new_node_obj.node
    address, _ = account.generate_account(node.web3, 10)
    status = 0
    try:
        result = client_new_node_obj.staking.increase_staking(0, address)
        log.info(result)
    except:
        status = 1
    assert status == 1


def test_AS_037(client_new_node_obj, get_generate_account):
    """
    (hesitation period) holdings less than the minimum threshold
    :param client_new_node_obj:
    :param get_generate_account:
    :return:
    """
    address, _ = get_generate_account
    client_new_node_obj.staking.create_staking(0, address, address)
    add_staking_amount = client_new_node_obj.economic.add_staking_limit
    result = client_new_node_obj.staking.increase_staking(0, address, amount=add_staking_amount - 1)
    log.info(result)
    assert result.get('Code') == 301104


def test_AS_038(client_new_node_obj, get_generate_account):
    """
    (hesitation period) when the verifier revoks the pledge, he/she shall apply for adding the pledge
    :param client_new_node_obj:
    :param get_generate_account:
    :return:
    """
    address, _ = get_generate_account
    client_new_node_obj.staking.create_staking(0, address, address)
    log.info("进入下个周期")
    client_new_node_obj.economic.wait_settlement_blocknum(client_new_node_obj.node)
    client_new_node_obj.staking.withdrew_staking(address)
    result = client_new_node_obj.staking.increase_staking(0, address)
    log.info(result)
    assert result.get('Code') == 301103


def test_AS_041_042_043_044(client_new_node_obj, get_generate_account):
    """
    (lockup period) normal increase
    (lockup period) overweight meets the minimum threshold
    (lock-up period) gas underinitiation overweight
    (lockup period) insufficient balance to initiate overweight pledge
    :param client_new_node_obj:
    :param get_generate_account:
    :return:
    """
    address, _ = get_generate_account
    client_new_node_obj.staking.create_staking(0, address, address)
    log.info("进入下个周期")
    client_new_node_obj.economic.wait_settlement_blocknum(client_new_node_obj.node)
    result = client_new_node_obj.staking.increase_staking(0, address)
    log.info(result)
    assert result.get('Code') == 0
    result = client_new_node_obj.ppos.getCandidateInfo(client_new_node_obj.node.node_id)
    log.info(result)
    staking_amount = client_new_node_obj.economic.create_staking_limit
    add_staking_amount = client_new_node_obj.economic.add_staking_limit
    assert result["Data"]["Shares"] == staking_amount + add_staking_amount
    assert result["Data"]["Released"] == staking_amount
    assert result["Data"]["ReleasedHes"] == add_staking_amount
    fig = {"gas": 1}
    status = 0
    try:
        result = client_new_node_obj.staking.increase_staking(0, address, transaction_cfg=fig)
        log.info(result)
    except:
        status = 1
    assert status == 1
    account = client_new_node_obj.economic.account
    node = client_new_node_obj.node
    address, _ = account.generate_account(node.web3, 10)
    status = 0
    try:
        result = client_new_node_obj.staking.increase_staking(0, address, address)
        log.info(result)
    except:
        status = 1
    assert status == 1


def test_AS_045(client_new_node_obj, get_generate_account):
    """
    The free amount is insufficient, the lock position is sufficient, and the free amount is added
    :param client_new_node_obj:
    :param get_generate_account:
    :return:
    """
    address, _ = get_generate_account
    result = client_new_node_obj.staking.create_staking(0, address, address)
    log.info(result)
    node = client_new_node_obj.node
    amount = node.eth.getBalance(address)
    log.info("Wallet balance{}".format(amount))
    locked_amount = amount - node.web3.toWei(1, "ether")

    plan = [{'Epoch': 1, 'Amount': locked_amount}]
    log.info("The balance of the wallet is used as a lock")
    result = client_new_node_obj.restricting.createRestrictingPlan(address, plan, address)
    log.info(result)
    amount = node.eth.getBalance(address)
    log.info("Check your wallet balance{}".format(amount))
    result = client_new_node_obj.staking.increase_staking(0, address)
    assert result.get('Code') == 302021
    locked_info = client_new_node_obj.ppos.getRestrictingInfo(address)
    log.info(locked_info)


def test_AS_046(client_new_node_obj, get_generate_account):
    """
    The free amount is insufficient, the lock position is sufficient, and the free amount is added
    :param client_new_node_obj:
    :param get_generate_account:
    :return:
    """
    address, _ = get_generate_account
    result = client_new_node_obj.staking.create_staking(0, address, address)
    log.info(result)
    node = client_new_node_obj.node
    amount = node.eth.getBalance(address)
    log.info("Wallet balance{}".format(amount))
    locked_amount = node.web3.toWei(100, "ether")
    plan = [{'Epoch': 1, 'Amount': locked_amount}]
    result = client_new_node_obj.restricting.createRestrictingPlan(address, plan, address)
    log.info(result)
    assert result.get('Code') == 0
    value = 101
    result = client_new_node_obj.staking.increase_staking(0, address, amount=value)
    log.info(result)
    assert result.get('Code') == 301104


def test_AS_047(client_new_node_obj, get_generate_account):
    """
    The amount of the increase is less than the threshold
    :param client_new_node_obj:
    :param get_generate_account:
    :return:
    """
    address, _ = get_generate_account
    result = client_new_node_obj.staking.create_staking(0, address, address)
    log.info(result)
    add_staking_amount = client_new_node_obj.economic.add_staking_limit
    result = client_new_node_obj.staking.increase_staking(0, address, amount=add_staking_amount - 1)
    log.info(result)
    assert result.get('Code') == 301104


def test_AS_048_049(client_new_node_obj, get_generate_account):
    """
    Increase the number of active withdrawal but still in the freeze period of the candidate
    Candidates whose holdings have been actively withdrawn and who have passed the freeze period
    :param client_new_node_obj:
    :param get_generate_account:
    :return:
    """
    address, _ = get_generate_account
    result = client_new_node_obj.staking.create_staking(0, address, address)
    assert result.get('Code') == 0
    log.info("Next settlement period")
    client_new_node_obj.economic.wait_settlement_blocknum(client_new_node_obj.node)
    result = client_new_node_obj.staking.withdrew_staking(address)
    assert result.get('Code') == 0
    result = client_new_node_obj.staking.increase_staking(0, address)
    log.info(result)
    assert result.get('Code') == 301103
    log.info("Next settlement period")
    client_new_node_obj.economic.wait_settlement_blocknum(client_new_node_obj.node)
    result = client_new_node_obj.staking.increase_staking(0, address)
    log.info(result)
    assert result.get('Code') == 301103


def test_AS_050_051(client_new_node_obj, get_generate_account,client_consensus_obj,greater_than_staking_amount):
    """
    Add to the list of candidates who have been penalized and are still in the freeze period
    A candidate whose holdings have been penalized has passed the freeze period
    :param client_new_node_obj:
    :param get_generate_account:
    :return:
    """
    address, _ = get_generate_account
    result = client_new_node_obj.staking.create_staking(0, address, address,amount=greater_than_staking_amount)
    assert result.get('Code') == 0
    log.info("Close one node")
    client_new_node_obj.node.stop()
    node = client_consensus_obj.node
    log.info("The next two periods")
    client_new_node_obj.economic.wait_settlement_blocknum(node,number=1)
    log.info("Restart the node")
    client_new_node_obj.node.start()
    result = client_new_node_obj.staking.increase_staking(0, address)
    log.info(result)
    log.info("Next settlement period")
    client_new_node_obj.economic.wait_settlement_blocknum(node)
    result = client_new_node_obj.staking.increase_staking(0, address)
    log.info(result)








