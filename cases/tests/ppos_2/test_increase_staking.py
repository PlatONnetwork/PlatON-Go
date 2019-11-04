# -*- coding: utf-8 -*-
from tests.lib.utils import *
import pytest
from copy import copy


def test_AS_031_032_39(client_new_node_obj, get_generate_account):
    """
    正常增持
    验证人用自由账户金额发起增持，满足最低门槛
    :param client_new_node_obj:
    :param get_generate_account:
    :return:
    """
    address, _ = get_generate_account
    msg = client_new_node_obj.staking.create_staking(0, address, address)
    assert msg["Code"] == 0
    msg = client_new_node_obj.staking.increase_staking(0, address)
    assert msg["Code"] == 0
    msg = client_new_node_obj.ppos.getCandidateInfo(client_new_node_obj.node.node_id)
    log.info(msg)
    staking_amount = client_new_node_obj.economic.create_staking_limit
    add_staking_amount = client_new_node_obj.economic.add_staking_limit
    assert msg["Data"]["Shares"] == staking_amount + add_staking_amount


def test_AS_033(client_new_node_obj, get_generate_account):
    """
    验证人不在验证人与候选人名单
    :param client_new_node_obj:
    :param get_generate_account:
    :return:
    """
    address, _ = get_generate_account
    msg = client_new_node_obj.staking.increase_staking(0, address)
    log.info(msg)
    assert msg["Code"] == 301102


def test_AS_034(client_new_node_obj, get_generate_account):
    address, _ = get_generate_account
    client_new_node_obj.staking.create_staking(0, address, address)
    fig = {"gas": 1}
    status = 0
    try:
        msg = client_new_node_obj.staking.increase_staking(0, address, transaction_cfg=fig)
        log.info(msg)
    except:
        status = 1
    assert status == 1


def test_AS_035(client_new_node_obj):
    account = client_new_node_obj.economic.account
    node = client_new_node_obj.node
    address, _ = account.generate_account(node.web3, 10)
    status = 0
    try:
        msg = client_new_node_obj.staking.increase_staking(0, address)
        log.info(msg)
    except:
        status = 1
    assert status == 1


def test_AS_037(client_new_node_obj,get_generate_account):
    """
    小于增持最低门槛
    :param client_new_node_obj:
    :param get_generate_account:
    :return:
    """
    address, _ = get_generate_account
    client_new_node_obj.staking.create_staking(0, address, address)
    add_staking_amount = client_new_node_obj.economic.add_staking_limit
    msg = client_new_node_obj.staking.increase_staking(0, address,amount = add_staking_amount-1)
    log.info(msg)
    assert msg["Code"] == 301104


def test_AS_038(client_new_node_obj,get_generate_account):
    """
    验证人撤销质押中，申请增持质押
    :param client_new_node_obj:
    :param get_generate_account:
    :return:
    """
    address, _ = get_generate_account
    client_new_node_obj.staking.create_staking(0, address, address)
    log.info("进入下个周期")
    client_new_node_obj.economic.wait_settlement_blocknum(client_new_node_obj.node)
    client_new_node_obj.staking.withdrew_staking(address)
    msg = client_new_node_obj.staking.increase_staking(0, address)
    log.info(msg)
    assert msg["Code"] == 301103


def test_AS_041_042_043_044(client_new_node_obj,get_generate_account):
    """
    锁定期正常增持
    满足增持最低门槛
    gas不足发起增持质押
    余额不足发起增持质押
    :param client_new_node_obj:
    :param get_generate_account:
    :return:
    """
    address, _ = get_generate_account
    client_new_node_obj.staking.create_staking(0, address, address)
    log.info("进入下个周期")
    client_new_node_obj.economic.wait_settlement_blocknum(client_new_node_obj.node)
    msg = client_new_node_obj.staking.increase_staking(0, address)
    log.info(msg)
    assert msg["Code"] == 0
    msg = client_new_node_obj.ppos.getCandidateInfo(client_new_node_obj.node.node_id)
    log.info(msg)
    staking_amount = client_new_node_obj.economic.create_staking_limit
    add_staking_amount = client_new_node_obj.economic.add_staking_limit
    assert msg["Data"]["Shares"] == staking_amount + add_staking_amount
    assert msg["Data"]["Released"] == staking_amount
    assert msg["Data"]["ReleasedHes"] == add_staking_amount
    fig = {"gas": 1}
    status = 0
    try:
        msg = client_new_node_obj.staking.increase_staking(0, address, transaction_cfg=fig)
        log.info(msg)
    except:
        status = 1
    assert status == 1
    account = client_new_node_obj.economic.account
    node = client_new_node_obj.node
    address, _ = account.generate_account(node.web3, 10)
    status = 0
    try:
        msg = client_new_node_obj.staking.increase_staking(0, address, address)
        log.info(msg)
    except:
        status = 1
    assert status == 1


def test_test_AS_045(client_new_node_obj,get_generate_account):
    """
    自由金额不足，锁仓金额充足，使用自由金额增持
    :param client_new_node_obj:
    :param get_generate_account:
    :return:
    """
    address, _ = get_generate_account
    msg = client_new_node_obj.staking.create_staking(0, address, address)
    log.info(msg)
    node = client_new_node_obj.node
    amount = node.eth.getBalance(address)
    log.info("钱包的余额{}".format(amount))
    locked_amount = amount - node.web3.toWei(1,"ether")

    plan = [{'Epoch': 1, 'Amount': locked_amount}]
    log.info("钱包的余额拿去做锁仓")
    msg = client_new_node_obj.restricting.CreateRestrictingPlan(address, plan)
    log.info(msg)
    amount = node.eth.getBalance(address)
    log.info("查看钱包还剩下的余额{}".format(amount))
    msg = client_new_node_obj.staking.increase_staking(0, address, address)
    log.info(msg)
    locked_info = client_new_node_obj.restricting.GetRestrictingInfo(address)
    log.info(locked_info)


















