# -*- coding: utf-8 -*-
'''
@Description: rpc用例
'''
import time

import allure
import pytest
from client_sdk_python import Web3
from hexbytes import HexBytes

from common.connect import connect_web3
from common.load_file import get_node_list, LoadFile
from conf import setting as conf

node_yml = conf.NODE_YML
collusion_list, nocollusion_list = get_node_list(node_yml)
# datadir: /home/platon/dark_test/node-16789/data
preDatadir = "/home/platon/dark_test/"
startApi = "eth,web3,net,txpool,platon,admin,personal"

# genesis config
genesis_path = conf.GENESIS_TMP
genesis_dict = LoadFile(genesis_path).get_data()
chainid = int(genesis_dict["config"]["chainId"])
amount = int(genesis_dict["config"]["cbft"]["amount"])
period = int(genesis_dict["config"]["cbft"]["period"])
validatorMode = genesis_dict["config"]["cbft"]["validatorMode"]

w3 = None
ws = None
host = ""
http_url = ""
http_port = 6789
ws_url = ""
ws_port = 5789
if len(collusion_list) > 0:
    try:
        host = collusion_list[0]["host"]
        http_url = collusion_list[0]["url"]
        http_port = collusion_list[0]["rpcport"]
        ws_url = collusion_list[0]["wsurl"]
        ws_port = collusion_list[0]["wsport"]
        w3 = connect_web3(http_url)
        ws = connect_web3(ws_url)
    except Exception as e:
        w3 = None
        ws = None

# start websocket rpc service
if w3 != None:
    try:
        assert True == w3.admin.startWS(host, ws_port, "*", startApi)
        print("\n启动websocket rpc服务成功:{}".format(host + ":" + str(ws_port)))
    except Exception as e:
        print("\n启动websocket rpc服务失败, error message:{}".format(e))

@allure.title("获取进程保存数据的目录:admin.datadir")
@pytest.mark.P1
def test_admin_datadir():
    if w3 != None:
        datadir = preDatadir + "node-" + str(collusion_list[0]["port"]) + "/data"
        try:
            assert w3.admin.datadir == datadir
            print("\n当前进程保存数据目录datadir:{}".format(w3.admin.datadir))
        except Exception as e:
            print("\n error message:{}".format(e))

@allure.title("获取程序的版本号和签名:admin.getProgramVersion()")
@pytest.mark.P1
def test_admin_getProgramVersion():
    if w3 != None:
        msg = w3.admin.getProgramVersion()
        ProgramVersionSign = msg["Sign"]
        ProgramVersion = msg["Version"]
        assert len(ProgramVersionSign) == 132
        assert ProgramVersion == 1794
        print("\n获取当前程序的版本号:【{}】, 版本签名:【{}】成功".format(ProgramVersion, ProgramVersionSign))

@allure.title("获取零知识证明信息:admin.getSchnorrNIZKProve()")
@pytest.mark.P1
def test_admin_getSchnorrNIZKProve():
    if w3 != None:
        blsproof = w3.admin.getSchnorrNIZKProve()
        assert len(blsproof) == 128
        print("\n获取零知识证明信息成功:{}".format(blsproof))

@allure.title("校验节点信息:admin.nodeInfo")
@pytest.mark.P1
def test_admin_nodeInfo():
    if w3 != None:
        try:
            genHash = HexBytes(w3.platon.getBlock(0)["hash"]).hex()
            nodeInfo = w3.admin.nodeInfo
            # node id
            id = nodeInfo["id"]
            # listen port
            listener = nodeInfo["ports"]["listener"]
            # discovery port
            discovery = nodeInfo["ports"]["discovery"]
            # config
            config = nodeInfo["protocols"]["eth"]["config"]
            assert id == collusion_list[0]["id"]
            assert listener == collusion_list[0]["port"]
            assert discovery == collusion_list[0]["port"]
            assert chainid == config["chainId"]
            assert amount == config["cbft"]["amount"]
            assert period == config["cbft"]["period"]
            assert validatorMode == config["cbft"]["validatorMode"]
            assert genHash == nodeInfo["protocols"]["eth"]["genesis"]

            print("\n校验节点信息成功,节点信息：{}".format(nodeInfo))
        except Exception as e:
            print("\nerror:{}".format(e))

@allure.title("和本节点的连接信息:admin.peers")
@pytest.mark.P1
def test_admin_peers():
    if w3 != None:
        try:
            print("\n本节点连接信息为:{}".format(w3.admin.peers))
        except Exception as e:
            print("\nerror message:{}".format(e))


#@allure.title("添加节点的连接:admin.addpeer()")
#@pytest.mark.P0
#def test_admin_addpeer():
#    test_init_diff_genesis_join_chain()
#    test_init_same_genesis_join_chain()

@allure.title("导出区块数据:admin.exportChain()")
@pytest.mark.P1
def test_admin_exportChain():
    if w3 != None:
        filePath = preDatadir + "chainData.txt"
        try:
            assert True == w3.admin.exportChain(filePath)
            print("\n区块数据导出成功,文件路径:【{}】".format(filePath))
        except Exception as e:
            print("\n区块数据导出失败,error message:【{}】".format(e))

@allure.title("导入区块数据:admin.importChain()")
@pytest.mark.P1
def test_admin_importChain():
    if w3 != None:
        try:
            filePath = preDatadir + "chainData.txt"
            assert True == w3.admin.importChain(filePath), "区块数据导入失败！"
            print("\n区块数据导入成功,文件路径:【{}】".format(filePath))
        except Exception as e:
            print("\n区块数据导入失败,error message:【{}】".format(e))

@allure.title("移除peer连接:admin.removePeer()")
@pytest.mark.P1
def test_admin_removePeer():
    if w3 != None:
        peers = w3.admin.peers
        if len(peers) > 0 :
            node_url = "enode://" + peers[0]["id"] + "@" + peers[0]["network"]["remoteAddress"]
            assert True == w3.admin.removePeer(node_url)
            print("\n移除peer成功:{}".format(node_url))


@allure.title("停止websocket rpc服务:admin.stopWS()")
@pytest.fixture()
def admin_stopWS():
    if ws != None:
        try:
            assert True == ws.admin.stopWS()
            print("\n停止websocket rpc服务成功:{}".format(ws_url))
        except Exception as e:
            print("\n停止websocket rpc服务失败, error message:{}".format(e))

@allure.title("启动websocket rpc服务:admin.startWS()")
@pytest.mark.P1
def test_admin_startWS(admin_stopWS):
    if w3 != None:
        try:
        #    new_w3 = connect_web3(collusion_list[0]["url"])
            assert True == w3.admin.startWS(host, ws_port, "*", startApi)
            print("\n启动websocket rpc服务成功:{}".format(ws_url))
        except Exception as e:
            print("\n启动websocket rpc服务失败, error message:{}".format(e))


@allure.title("停止http rpc服务:admin.stopRPC()")
@pytest.fixture()
def admin_stopRPC():
    if w3 != None:
        try:
            assert True == w3.admin.stopRPC()
            print("\n停止http rpc服务成功:{}".format(http_url))
        except Exception as e:
            print("\n停止http rpc服务失败, error message:{}".format(e))

@allure.title("启动http rpc服务:admin.startRPC()")
@pytest.mark.P0
def test_admin_startRPC(admin_stopRPC):
    if ws != None:
        try:
            assert True == ws.admin.startRPC(host, http_port, "*", startApi)
            print("\n启动http rpc服务成功:{}".format(http_url))
        except Exception as e:
            print("\n启动http rpc服务失败,error message:{}".format(e))

if __name__ == '__main__':
    pytest.main(['-v','test_rpc_admin.py'])