# -*- coding: utf-8 -*-
'''
@Description: rpc用例
'''
import json

import allure
import pytest
import rlp
from client_sdk_python import Web3
from client_sdk_python.eth import Eth


@allure.title("Check if the version of the protocol is 63")
@pytest.mark.P1
def test_platon_protocolVersion(global_running_env):
    node = global_running_env.get_rand_node()
    assert node.eth.protocolVersion == '63'


@allure.title("Get the amount of the account")
@pytest.mark.P1
def test_platon_GetBalance(global_running_env):
    node = global_running_env.get_rand_node()
    platon = Eth(node.web3)
    account = global_running_env.account
    addr = account.account_with_money["address"]
    from_addr = Web3.toChecksumAddress(addr)
    # balance = platon.getBalance(from_addr)
    balance = platon.getBalance(node.web3.pipAddress)
    assert balance == 0


def platon_call(platon, from_addr, to_addr="0x1000000000000000000000000000000000000002", data=""):
    recive = platon.call({
        "from": from_addr,
        "to": to_addr,
        "data": data
    })
    recive = str(recive, encoding="utf8")
    recive = recive.replace('\\', '').replace('"[', '[').replace(']"', ']')
    recive = json.loads(recive)
    return recive


@allure.title("Call the built-in contract interface with platon.call")
@pytest.mark.P1
def test_platon_call(global_running_env):
    node = global_running_env.get_rand_node()
    platon = Eth(node.web3)
    account = global_running_env.account
    addr = account.account_with_money["address"]
    from_addr = Web3.toChecksumAddress(addr)

    to_addr = node.web3.stakingAddress
    data = rlp.encode([rlp.encode(int(1100))])
    recive = platon_call(platon, from_addr, to_addr, data)
    assert recive != "0x"
    # not exist interface on staking contract
    data = rlp.encode([rlp.encode(int(2222))])

    status = 0
    try:
        recive = platon_call(platon, from_addr, to_addr, data)
        assert recive == "0x"
        status = 1
    except Exception as e:
        print("\nQuery the built-in contract interface that does not exist and return an exception.:{}".format(e))
    assert status == 0


@allure.title("Get node double-sign evidence")
@pytest.mark.P1
def test_platon_evidences(global_running_env):
    node = global_running_env.get_rand_node()
    platon = Eth(node.web3)
    ret = platon.evidences
    assert ret is not None


@allure.title("Get the aggregate signature of any block")
@pytest.mark.P1
@pytest.mark.compatibility
def test_platon_getPrepareQC(global_running_env):
    node = global_running_env.get_rand_node()
    platon = Eth(node.web3)
    blockNumber = platon.blockNumber
    qc = platon.getPrepareQC(blockNumber)
    assert qc is not None


if __name__ == '__main__':
    pytest.main(['-v', 'test_rpc_platon.py'])
