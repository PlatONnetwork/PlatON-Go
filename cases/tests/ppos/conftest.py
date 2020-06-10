#!/usr/bin/env python3
# -*- coding: utf-8 -*-
#   @Time    : 2020/1/2 18:06
#   @Author  : PlatON-Developer
#   @Site    : https://github.com/PlatONnetwork/

import pytest
from decimal import Decimal
from tests.lib.utils import assert_code


def calculate(big_int, mul):
    return int(Decimal(str(big_int)) * Decimal(mul))


def create_stakings(clients, reward):
    for client in clients:
        create_staking(client, reward)


def create_staking(client, reward):
    amount = calculate(client.economic.create_staking_limit, 5)
    staking_amount = calculate(client.economic.create_staking_limit, 2)
    staking_address, _ = client.economic.account.generate_account(client.node.web3, amount)
    delegate_address, _ = client.economic.account.generate_account(client.node.web3, client.economic.add_staking_limit * 5)
    result = client.staking.create_staking(0, staking_address, staking_address, amount=staking_amount, reward_per=reward)
    assert_code(result, 0)
    setattr(client, "staking_amount", staking_amount)
    return staking_address, delegate_address


@pytest.fixture()
def staking_node_client(client_new_node):
    reward = 2000
    staking_address, delegate_address = create_staking(client_new_node, reward)
    setattr(client_new_node, "staking_address", staking_address)
    setattr(client_new_node, "delegate_address", delegate_address)
    setattr(client_new_node, "reward", reward)
    yield client_new_node
    client_new_node.economic.env.deploy_all()


@pytest.fixture()
def delegate_node_client(client_new_node):
    reward = 1000
    staking_address, delegate_address = create_staking(client_new_node, reward)
    result = client_new_node.delegate.delegate(0, delegate_address, amount=client_new_node.economic.add_staking_limit * 2)
    assert_code(result, 0)
    setattr(client_new_node, "staking_address", staking_address)
    setattr(client_new_node, "delegate_address", delegate_address)
    setattr(client_new_node, "reward", reward)
    yield client_new_node
    client_new_node.economic.env.deploy_all()