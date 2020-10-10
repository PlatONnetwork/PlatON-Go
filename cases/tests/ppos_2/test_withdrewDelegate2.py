# -*- coding: utf-8 -*-
from tests.lib.utils import *
import pytest


# Undo delegate use cases from 031 to 049


@pytest.mark.P2
def test_ROE_031(staking_delegate_client):
    """
    :param client_new_node:
    :return:
    """
    client = staking_delegate_client
    delegate_address = client.delegate_address
    node = client.node
    economic = client.economic
    # Return a pledge
    client.staking.withdrew_staking(client.staking_address)
    # The next cycle
    client.economic.wait_settlement(node)
    result = client.delegate.delegate(0, delegate_address)
    log.info(result)
    # The next two cycle
    client.economic.wait_settlement(node, 2)
    balance1 = client.node.eth.getBalance(delegate_address)
    log.info("The wallet balance:{}".format(balance1))

    result = client.delegate.withdrew_delegate(client.staking_blocknum, delegate_address)
    assert_code(result, 0)
    balance2 = client.node.eth.getBalance(delegate_address)
    log.info("The wallet balance:{}".format(balance2))
    delegate_limit = economic.delegate_limit
    assert delegate_limit - (balance2 - balance1) < node.web3.toWei(1, "ether")


@pytest.mark.P2
def test_ROE_032_035(staking_delegate_client):
    client = staking_delegate_client
    delegate_address = client.delegate_address
    node = client.node
    economic = client.economic
    log.info("The next cycle")
    economic.wait_settlement(node)
    result = client.delegate.delegate(0, delegate_address, amount=client.delegate_amount * 2)
    assert_code(result, 0)
    balance1 = node.eth.getBalance(delegate_address)
    log.info("Wallet balance{}".format(balance1))
    amount = client.delegate_amount * 2 + node.web3.toWei(1, "ether")
    result = client.delegate.withdrew_delegate(client.staking_blocknum, delegate_address,
                                               amount=amount)
    assert_code(result, 0)
    balance2 = node.eth.getBalance(delegate_address)
    log.info("Wallet balance{}".format(balance2))
    assert amount - (balance2 - balance1) < node.web3.toWei(1, "ether")


@pytest.mark.P2
def test_ROE_033_034(staking_delegate_client):
    client = staking_delegate_client
    delegate_address = client.delegate_address
    node = client.node
    economic = client.economic
    log.info("The next cycle")
    economic.wait_settlement(node)
    result = client.delegate.delegate(0, delegate_address)
    assert_code(result, 0)
    balance1 = node.eth.getBalance(delegate_address)
    log.info("Wallet balance{}".format(balance1))
    result = client.delegate.withdrew_delegate(client.staking_blocknum, delegate_address)
    assert_code(result, 0)
    balance2 = node.eth.getBalance(delegate_address)
    log.info("Wallet balance{}".format(balance2))
    assert client.delegate_amount - (balance2 - balance1) < node.web3.toWei(1, "ether")


@pytest.mark.P2
def test_ROE_038(staking_delegate_client):
    client = staking_delegate_client
    delegate_address = client.delegate_address
    node = client.node
    economic = client.economic
    log.info("The next cycle")
    economic.wait_settlement(node)
    lockup_amount = client.node.web3.toWei(20, "ether")
    plan = [{'Epoch': 1, 'Amount': lockup_amount}]
    # Create a lock plan
    result = client.restricting.createRestrictingPlan(delegate_address, plan, delegate_address)
    assert_code(result, 0)
    result = client.delegate.delegate(1, delegate_address)
    assert_code(result, 0)
    result = client.delegate.delegate(0, delegate_address)
    assert_code(result, 0)
    balance1 = node.eth.getBalance(delegate_address)
    log.info("Wallet balance{}".format(balance1))
    result = client.delegate.withdrew_delegate(client.staking_blocknum, delegate_address)
    assert_code(result, 0)
    balance2 = node.eth.getBalance(delegate_address)
    log.info("Wallet balance{}".format(balance2))
    msg = client.ppos.getRestrictingInfo(delegate_address)
    log.info(msg)
    delegate_limit = client.delegate_amount
    assert delegate_limit - (balance2 - balance1) < node.web3.toWei(1, "ether")
    assert msg["Ret"]["Pledge"] == delegate_limit


@pytest.mark.P2
def test_ROE_039(staking_delegate_client):
    client = staking_delegate_client
    delegate_address = client.delegate_address
    node = client.node
    economic = client.economic
    log.info("The next cycle")
    economic.wait_settlement(node)
    lockup_amount = client.node.web3.toWei(20, "ether")
    plan = [{'Epoch': 1, 'Amount': lockup_amount}]
    # Create a lock plan
    result = client.restricting.createRestrictingPlan(delegate_address, plan, delegate_address)
    assert_code(result, 0)
    result = client.delegate.delegate(1, delegate_address, amount=client.delegate_amount * 2)
    assert_code(result, 0)
    result = client.delegate.delegate(0, delegate_address)
    assert_code(result, 0)
    msg = client.ppos.getDelegateInfo(client.staking_blocknum, delegate_address, node.node_id)
    log.info(msg)
    balance1 = node.eth.getBalance(delegate_address)
    log.info("Wallet balance{}".format(balance1))
    amount = client.delegate_amount * 2
    result = client.delegate.withdrew_delegate(client.staking_blocknum, delegate_address,
                                               amount=amount)
    assert_code(result, 0)
    balance2 = node.eth.getBalance(delegate_address)
    log.info("Wallet balance{}".format(balance2))
    msg = client.ppos.getRestrictingInfo(delegate_address)
    log.info(msg)
    assert client.delegate_amount - (balance2 - balance1) < node.web3.toWei(1, "ether")
    assert msg["Ret"]["Pledge"] == client.delegate_amount


@pytest.mark.P2
def test_ROE_040(free_locked_delegate_client):
    client = free_locked_delegate_client
    delegate_address = client.delegate_address
    node = client.node
    balance1 = node.eth.getBalance(delegate_address)
    log.info("Wallet balance{}".format(balance1))
    amount = client.delegate_amount * 2 + client.node.web3.toWei(1, "ether")
    result = client.delegate.withdrew_delegate(client.staking_blocknum, delegate_address,
                                               amount=amount)
    assert_code(result, 0)
    balance2 = node.eth.getBalance(delegate_address)
    log.info("Wallet balance{}".format(balance2))
    assert client.delegate_amount * 2 - (balance2 - balance1) < node.web3.toWei(1, "ether")
    msg = client.ppos.getRestrictingInfo(delegate_address)
    log.info(msg)
    assert msg["Ret"]["Pledge"] == 0


@pytest.mark.P2
def test_ROE_041(free_locked_delegate_client):
    client = free_locked_delegate_client
    delegate_address = client.delegate_address
    node = client.node
    balance1 = node.eth.getBalance(delegate_address)
    log.info("Wallet balance{}".format(balance1))
    amount = client.delegate_amount * 3
    result = client.delegate.withdrew_delegate(client.staking_blocknum, delegate_address,
                                               amount=amount)
    assert_code(result, 0)
    balance2 = node.eth.getBalance(delegate_address)
    log.info("Wallet balance{}".format(balance2))
    msg = client.ppos.getRestrictingInfo(delegate_address)
    log.info(msg)
    assert client.delegate_amount * 2 - (balance2 - balance1) < node.web3.toWei(1, "ether")


@pytest.mark.P2
def test_ROE_042(free_locked_delegate_client):
    client = free_locked_delegate_client
    delegate_address = client.delegate_address
    node = client.node
    economic = client.economic
    log.info("The next cycle")
    economic.wait_settlement(node)
    result = client.delegate.delegate(0, delegate_address)
    assert_code(result, 0)
    balance1 = node.eth.getBalance(delegate_address)
    log.info("Wallet balance{}".format(balance1))
    amount = client.delegate_amount * 3
    result = client.delegate.withdrew_delegate(client.staking_blocknum, delegate_address,
                                               amount=amount)
    assert_code(result, 0)
    balance2 = node.eth.getBalance(delegate_address)
    log.info("Wallet balance{}".format(balance2))
    msg = client.ppos.getRestrictingInfo(delegate_address)
    log.info(msg)
    assert client.delegate_amount * 3 - (balance2 - balance1) < node.web3.toWei(1, "ether")
    assert msg["Ret"]["Pledge"] == 0


@pytest.mark.P2
def test_ROE_043(free_locked_delegate_client):
    client = free_locked_delegate_client
    delegate_address = client.delegate_address
    node = client.node
    economic = client.economic
    log.info("The next cycle")
    economic.wait_settlement(node)
    result = client.delegate.delegate(0, delegate_address)
    assert_code(result, 0)
    result = client.delegate.delegate(1, delegate_address)
    assert_code(result, 0)
    balance1 = node.eth.getBalance(delegate_address)
    log.info("Wallet balance{}".format(balance1))
    amount = client.delegate_amount * 4
    result = client.delegate.withdrew_delegate(client.staking_blocknum, delegate_address,
                                               amount=amount)
    assert_code(result, 0)
    balance2 = node.eth.getBalance(delegate_address)
    log.info("Wallet balance{}".format(balance2))
    msg = client.ppos.getRestrictingInfo(delegate_address)
    log.info(msg)
    assert client.delegate_amount * 3 - (balance2 - balance1) < node.web3.toWei(1, "ether")
    assert msg["Ret"]["Pledge"] == client.delegate_amount


@pytest.mark.P2
def test_ROE_044(free_locked_delegate_client):
    client = free_locked_delegate_client
    delegate_address = client.delegate_address
    node = client.node
    economic = client.economic
    log.info("The next cycle")
    economic.wait_settlement(node)
    result = client.delegate.delegate(0, delegate_address)
    assert_code(result, 0)
    result = client.delegate.delegate(1, delegate_address)
    assert_code(result, 0)
    balance1 = node.eth.getBalance(delegate_address)
    log.info("Wallet balance{}".format(balance1))
    amount = client.delegate_amount * 4
    result = client.delegate.withdrew_delegate(client.staking_blocknum, delegate_address,
                                               amount=amount)
    assert_code(result, 0)
    balance2 = node.eth.getBalance(delegate_address)
    log.info("Wallet balance{}".format(balance2))
    msg = client.ppos.getRestrictingInfo(delegate_address)
    log.info(msg)
    log.info("Wallet balance{}".format(balance2))
    assert client.delegate_amount * 3 - (balance2 - balance1) < node.web3.toWei(1, "ether")
    assert msg["Ret"]["Pledge"] == 0


@pytest.mark.P2
def test_ROE_045(staking_delegate_client):
    client = staking_delegate_client
    delegate_address = client.delegate_address
    node = client.node
    economic = client.economic
    log.info("The next cycle")
    economic.wait_settlement(node)
    result = client.delegate.delegate(0, delegate_address)
    assert_code(result, 0)
    balance1 = node.eth.getBalance(delegate_address)
    log.info("Wallet balance{}".format(balance1))
    amount = client.delegate_amount * 3 - node.web3.toWei(1, "ether")
    result = client.delegate.withdrew_delegate(client.staking_blocknum, delegate_address,
                                               amount=amount)
    assert_code(result, 0)
    balance2 = node.eth.getBalance(delegate_address)
    log.info("Wallet balance{}".format(balance2))
    assert client.delegate_amount * 3 - (balance2 - balance1) < node.web3.toWei(1, "ether")


@pytest.mark.P2
def test_ROE_048(staking_delegate_client):
    client = staking_delegate_client
    delegate_address = client.delegate_address
    node = client.node
    economic = client.economic
    log.info("The next cycle")
    economic.wait_settlement(node)
    result = client.delegate.delegate(0, delegate_address)
    assert_code(result, 0)
    balance1 = node.eth.getBalance(delegate_address)
    log.info("Wallet balance{}".format(balance1))
    amount = client.delegate_amount * 3 - node.web3.toWei(1, "ether")
    result = client.delegate.withdrew_delegate(client.staking_blocknum, delegate_address,
                                               amount=amount)
    assert_code(result, 0)
    balance2 = node.eth.getBalance(delegate_address)
    log.info("Wallet balance{}".format(balance2))
    msg = client.ppos.getRestrictingInfo(delegate_address)
    log.info(msg)
    assert client.delegate_amount * 3 - (balance2 - balance1) < node.web3.toWei(1, "ether")


@pytest.mark.P2
def test_ROE_049(staking_delegate_client):
    client = staking_delegate_client
    delegate_address = client.delegate_address
    node = client.node
    economic = client.economic
    log.info("The next cycle")
    economic.wait_settlement(node)
    lockup_amount = client.node.web3.toWei(20, "ether")
    plan = [{'Epoch': 1, 'Amount': lockup_amount}]
    # Create a lock plan
    result = client.restricting.createRestrictingPlan(delegate_address, plan, delegate_address)
    assert_code(result, 0)
    result = client.delegate.delegate(1, delegate_address)
    assert_code(result, 0)
    result = client.delegate.delegate(0, delegate_address)
    assert_code(result, 0)
    balance1 = node.eth.getBalance(delegate_address)
    log.info("Wallet balance{}".format(balance1))
    amount = client.delegate_amount * 4
    result = client.delegate.withdrew_delegate(client.staking_blocknum, delegate_address,
                                               amount=amount)
    assert_code(result, 0)
    balance2 = node.eth.getBalance(delegate_address)
    log.info("Wallet balance{}".format(balance2))
    msg = client.ppos.getRestrictingInfo(delegate_address)
    log.info(msg)
    assert client.delegate_amount * 3 - (balance2 - balance1) < node.web3.toWei(1, "ether")
    assert msg["Ret"]["Pledge"] == 0
