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


@allure.title("查看协议的版本是否是63")
@pytest.mark.P1
def test_platon_protocolVersion(global_running_env):
    node = global_running_env.get_rand_node()
    assert node.eth.protocolVersion == '63'
    print("\n当前节点的协议版本号为:63")


@allure.title("获取账户的金额")
@pytest.mark.P1
def test_platon_GetBalance(global_running_env):
    node = global_running_env.get_rand_node()
    platon = Eth(node.web3)
    account = global_running_env.account
    addr = account.account_with_money["address"]
    from_addr = Web3.toChecksumAddress(addr)
    balance = platon.getBalance(from_addr)
    print("\n当前账户【{}】的余额为:{}".format(from_addr, balance))
    balance = platon.getBalance("0x1111111111111111111111111111111111111111")
    assert balance == 0
    print("\n当前不存在的账户【{}】余额为:{}".format("0x1111111111111111111111111111111111111111", balance))


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


@allure.title("使用platon.call调用内置合约接口,如:getCandidateList,或不存在的接口")
@pytest.mark.P1
def test_platon_call(global_running_env):
    node = global_running_env.get_rand_node()
    platon = Eth(node.web3)
    account = global_running_env.account
    addr = account.account_with_money["address"]
    from_addr = Web3.toChecksumAddress(addr)

    to_addr = Web3.toChecksumAddress("0x1000000000000000000000000000000000000002")
    data = rlp.encode([rlp.encode(int(1100))])
    recive = platon_call(platon, from_addr, to_addr, data)
    assert recive != "0x"
    # not exist interface on staking contract
    data = rlp.encode([rlp.encode(int(2222))])

    # 异常测试场景
    status = 0
    try:
        recive = platon_call(platon, from_addr, to_addr, data)
        assert recive == "0x"
        status = 1
    except Exception as e:
        print("\n查询不存在的内置合约接口返回异常:{}".format(e))
    assert status == 0


@allure.title("查询节点双签证据：platon.Evidences")
@pytest.mark.P1
def test_platon_evidences(global_running_env):
    node = global_running_env.get_rand_node()
    platon = Eth(node.web3)
    ret = platon.evidences
    assert ret is not None


@allure.title("查询任意区块的聚合签名：GetPrepareQC")
@pytest.mark.P1
@pytest.mark.compatibility
def test_platon_getPrepareQC(global_running_env):
    node = global_running_env.get_rand_node()
    platon = Eth(node.web3)
    blockNumber = platon.blockNumber
    qc = platon.getPrepareQC(blockNumber)
    assert qc is not None
    print("\n查询区块的聚合签名成功:区块高度【{}】,签名数据:【{}】".format(blockNumber, qc))


if __name__ == '__main__':
    pytest.main(['-v', 'test_rpc_platon.py'])
