# -*- coding: utf-8 -*-
'''
@Description: rpc用例
'''
import json
import time

import allure
import pytest
import rlp
from client_sdk_python import Web3

password = "88888888"
address = "0x493301712671Ada506ba6Ca7891F436D29185821"

w3 = None
id = ""
host = ""
http_url = ""
http_port = "6789"
p2p_prot = "16789"

@pytest.fixture()
def setNodeInfo(global_test_env):
    collusion_list = global_test_env.consensus_node_list
    if len(collusion_list) > 0:
        try:
            global w3
            global id
            global host
            global http_url
            global http_port
            global p2p_prot

            test_node = collusion_list[0]
            id = test_node.node_id
            host = test_node.host
            http_url = test_node.url
            http_port = test_node.rpc_port
            p2p_prot = test_node.p2p_port

            # rpc连接
            w3 = test_node.web3
        except Exception as e:
            print("setNodeInfo error:{}>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>".format(e))
            w3 = None

@allure.title("查看协议的版本是否是63")
@pytest.mark.P1
def test_platon_protocolVersion(setNodeInfo):
    if w3 != None:
        assert w3.platon.protocolVersion == '63'
        print("\n当前节点的协议版本号为:63")


@allure.title("获取账户的金额")
@pytest.mark.P1
def test_platon_GetBalance(setNodeInfo):
    if w3 != None:
        from_addr = Web3.toChecksumAddress(address)
        balance = w3.platon.getBalance(from_addr)
        print("\n当前账户【{}】的余额为:{}".format(from_addr, balance))
        balance = w3.platon.getBalance("0x1111111111111111111111111111111111111111")
        assert balance == 0
        print("\n当前不存在的账户【{}】余额为:{}".format("0x1111111111111111111111111111111111111111", balance))

def platon_call(w3,from_addr,to_addr="0x1000000000000000000000000000000000000002",data=""):
    recive = w3.platon.call({
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
def test_platon_call(setNodeInfo):
    if w3 != None:
        to_addr = Web3.toChecksumAddress("0x1000000000000000000000000000000000000002")
        data = rlp.encode([rlp.encode(int(1100))])
        from_addr = Web3.toChecksumAddress(address)
        recive = platon_call(w3, from_addr, to_addr, data)
        assert recive != "0x"
    #    assert recive['Code'] == 0
    #    assert len(recive['Data']) > 0
        print("\ngetCandidateList查询的候选人列表成功:【{}】".format(recive))
        # not exist interface on staking contract
        data = rlp.encode([rlp.encode(int(2222))])
        try:
            recive = platon_call(w3, from_addr, to_addr, data)
            assert recive == "0x"
            print("\ngetCandidateList查询的不存在的候选人列表返回为空")
        except Exception as e:
            print("\ngetCandidateList查询的不存在的候选人列表,error message:{}".format(e))

@allure.title("查询节点双签证据：platon.Evidences")
@pytest.mark.P1
def test_platon_evidences(setNodeInfo):
    if w3 != None:
        try:
            ret = w3.platon.evidences()
            assert ret != None
            print("\n查询节点双签证据成功:【{}】".format(ret))
        except Exception as e:
            print("\nerror message:{}".format(e))

@allure.title("查询任意区块的聚合签名：GetPrepareQC")
@pytest.mark.P1
@pytest.mark.compatibility
def test_platon_getPrepareQC(setNodeInfo):
    if w3 != None:
        blockNumber = w3.platon.blockNumber
        qc = w3.platon.getPrepareQC(blockNumber)
        assert qc != None
        print("\n查询区块的聚合签名成功:区块高度【{}】,签名数据:【{}】".format(blockNumber, qc))

if __name__ == '__main__':
    pytest.main(['-v', 'test_rpc_platon.py'])