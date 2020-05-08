# -*- coding: utf-8 -*-
'''
@Description: rpc用例
'''

import allure
import pytest
from hexbytes import HexBytes

startApi = "eth,web3,net,txpool,platon,admin,personal"


@allure.title("获取进程保存数据的目录:admin.datadir")
@pytest.mark.P1
def test_admin_datadir(global_running_env):
    node = global_running_env.get_rand_node()
    dataDir = node.remote_data_dir
    assert node.admin.datadir == dataDir
    print("\n当前进程保存数据目录datadir=================:{}".format(node.admin.datadir))


@allure.title("获取程序的版本号和签名:admin.getProgramVersion()")
@pytest.mark.P1
def test_admin_getProgramVersion(global_running_env):
    node = global_running_env.get_rand_node()
    msg = node.admin.getProgramVersion()
    ProgramVersionSign = msg["Sign"]
    ProgramVersion = msg["Version"]
    assert len(ProgramVersionSign) == 132
    assert ProgramVersion >= 1794
    print("\n获取当前程序的版本号:【{}】, 版本签名:【{}】成功".format(ProgramVersion, ProgramVersionSign))


@allure.title("获取零知识证明信息:admin.getSchnorrNIZKProve()")
@pytest.mark.P1
def test_admin_getSchnorrNIZKProve(global_running_env):
    node = global_running_env.get_rand_node()
    blsproof = node.admin.getSchnorrNIZKProve()
    assert len(blsproof) == 128
    print("\n获取零知识证明信息成功:{}".format(blsproof))


@allure.title("校验节点信息:admin.nodeInfo")
@pytest.mark.P1
def test_admin_nodeInfo(global_running_env):
    node = global_running_env.get_rand_node()
    genHash = HexBytes(node.eth.getBlock(0)["hash"]).hex()
    nodeInfo = node.admin.nodeInfo

    # config
    config = nodeInfo["protocols"]["platon"]["config"]

    # node id
    assert node.node_id == nodeInfo["id"]
    # listen port
    assert node.p2p_port == str(nodeInfo["ports"]["listener"])
    # discovery port
    assert node.p2p_port == str(nodeInfo["ports"]["discovery"])

    assert global_running_env.amount == config["cbft"]["amount"]
    assert global_running_env.period == config["cbft"]["period"]
    assert global_running_env.validatorMode == config["cbft"]["validatorMode"]

    assert global_running_env.chain_id == config["chainId"]
    assert genHash == nodeInfo["protocols"]["platon"]["genesis"]

    print("\n校验节点信息成功,节点信息：{}".format(nodeInfo))


@allure.title("和本节点的连接信息:admin.peers")
@pytest.mark.P1
def test_admin_peers(global_running_env):
    node = global_running_env.get_rand_node()
    lenPeers = len(node.admin.peers)
    assert lenPeers >= 0


@allure.title("导出区块数据:admin.exportChain()")
@pytest.mark.P1
@pytest.fixture()
def test_admin_exportChain(global_running_env):
    node = global_running_env.get_rand_node()
    filePath = node.admin.datadir + "chainData.txt"
    assert True == node.admin.exportChain(filePath)
    print("\n区块数据导出成功,文件路径:【{}】".format(filePath))
    yield node


@allure.title("导入区块数据:admin.importChain()")
@pytest.mark.P1
def test_admin_importChain(test_admin_exportChain):
    filePath = test_admin_exportChain.admin.datadir + "chainData.txt"
    assert True == test_admin_exportChain.admin.importChain(filePath), "区块数据导入失败！"
    print("\n区块数据导入成功,文件路径:【{}】".format(filePath))


@allure.title("移除peer连接:admin.removePeer()")
@pytest.mark.P1
def test_admin_removePeer(global_running_env):
    node = global_running_env.get_rand_node()
    peers = node.admin.peers
    if len(peers) > 0:
        node_url = "enode://" + peers[0]["id"] + "@" + peers[0]["network"]["remoteAddress"]
        assert True == node.admin.removePeer(node_url)
        print("\n移除peer成功:{}".format(node_url))


@allure.title("停止websocket rpc服务:admin.stopWS()")
@pytest.fixture()
def admin_stopWS(global_running_env):
    node = global_running_env.get_rand_node()
    try:
        ws = node.ws_web3
        assert True == ws.admin.stopWS()
    except Exception as e:
        print("web socket服务没有启动===================>")

    yield node


@allure.title("启动websocket rpc服务:admin.startWS()")
@pytest.mark.P1
@pytest.fixture()
def test_admin_startWS(admin_stopWS):
    node = admin_stopWS
    if None == node.wsport:
        node.wsport = 5789
    if None == node.wsurl:
        node.wsurl = "ws://" + str(node.host) + ":" + str(node.wsport)
    assert True == node.admin.startWS(node.host, int(node.wsport), "*", startApi)
    print("web socket启动成功===================>")

    ws = node.ws_web3
    print("blockNumber:{}===================>".format(ws.eth.blockNumber))


@allure.title("停止http rpc服务:admin.stopRPC()")
@pytest.fixture()
def admin_stopRPC(global_running_env):
    node = global_running_env.get_rand_node()
    try:
        ws = node.ws_web3
        assert True == ws.admin.stopRPC()
        print("\n停止http rpc服务成功")
    except Exception as e:
        pass

    yield node


'''
@allure.title("启动http rpc服务:admin.startRPC()")
@pytest.mark.P0
def test_admin_startRPC(admin_stopRPC):
    node = admin_stopRPC
    ws = node.ws_web3
    assert True == ws.admin.startRPC(admin_stopRPC.host, int(admin_stopRPC.rpc_port))
    print("\n启动http rpc服务成功.")

    print("blockNumber:{}===================>".format(node.eth.blockNumber))
'''

if __name__ == '__main__':
    pytest.main(['-v', 'test_rpc_admin.py'])
