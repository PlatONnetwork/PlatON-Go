# -*- coding: utf-8 -*-

from tests.lib.utils import *
import pytest
import allure


def test_RV_069(client_new_node_obj, get_generate_account, get_amount_limit):
    """
     Withdraw pledge in hesitation period
    :param client_new_node_obj:
    :param get_generate_account:
    :return:
    """
    address, pri_key = get_generate_account
    result = client_new_node_obj.staking.create_staking(0, address, address)
    assert result.get('Code') == 0
    amount = client_new_node_obj.node.eth.getBalance(address)
    log.info("The wallet balance:{}".format(amount))
    result = client_new_node_obj.staking.withdrew_staking(address)
    assert result.get('Code') == 0
    amount1 = client_new_node_obj.node.eth.getBalance(address)
    log.info("The wallet balance:{}".format(amount1))
    staking_amount, _ = get_amount_limit
    amount_dill = amount1 - amount
    log.info(amount_dill)
    assert staking_amount - amount_dill < client_new_node_obj.node.web3.toWei(1, "ether")


"""To debug"""
def test_RV_070(client_new_node_obj, get_generate_account, get_amount_limit):
    """
    :param client_new_node_obj:
    :param get_generate_account:
    :param get_amount_limit:
    :return:
    """
    address, pri_key = get_generate_account
    result = client_new_node_obj.staking.create_staking(0, address, address)
    assert result.get('Code') == 0
    amount = client_new_node_obj.node.eth.getBalance(address)
    log.info("The wallet balance:{}".format(amount))
    log.info("Enter lockup period")
    client_new_node_obj.economic.wait_settlement_blocknum(client_new_node_obj.node)
    result = client_new_node_obj.staking.withdrew_staking(address)
    assert result.get('Code') == 0
    client_new_node_obj.economic.wait_settlement_blocknum(client_new_node_obj.node, number=1)
    amount1 = client_new_node_obj.node.eth.getBalance(address)
    log.info("The wallet balance:{}".format(amount1))
    staking_amount, _ = get_amount_limit
    amount_dill = amount1 - amount
    assert staking_amount - amount_dill < client_new_node_obj.node.web3.toWei(1, "ether")


"""To debug"""
def test_RV_071(client_new_node_obj, get_generate_account, get_amount_limit):
    """
     Do revocation pledge when there is lockup period or hesitation period
    :param client_new_node_obj:
    :param get_generate_account:
    :param get_amount_limit:
    :return:
    """
    address, pri_key = get_generate_account
    staking_amount, addstaking_amount = get_amount_limit
    result = client_new_node_obj.staking.create_staking(0, address, address)
    assert result.get('Code') == 0
    log.info("Enter lockup period")
    client_new_node_obj.economic.wait_settlement_blocknum(client_new_node_obj.node)
    result = client_new_node_obj.staking.increase_staking(0, address)
    assert result.get('Code') == 0

    # block_reward, staking_reward = get_first_year_block_reward()
    # block_reward = Web3.fromWei(block_reward, "ether")

    amount = client_new_node_obj.node.eth.getBalance(address)
    log.info("The wallet balance:{}".format(amount))
    result = client_new_node_obj.staking.withdrew_staking(address)
    assert result.get('Code') == 0
    amount1 = client_new_node_obj.node.eth.getBalance(address)
    log.info("The wallet balance:{}".format(amount1))
    amount_dill = amount1 - amount
    log.info("The amount returned during the hesitation period:{}".format(amount_dill))
    assert addstaking_amount - amount_dill < client_new_node_obj.node.web3.toWei(1, "ether")
    log.info("Enter two lockup period")


    client_new_node_obj.economic.wait_settlement_blocknum(client_new_node_obj.node, number=1)
    amount2 = client_new_node_obj.node.eth.getBalance(address)
    log.info("The balance after the lockup period ends:{}".format(amount2))



