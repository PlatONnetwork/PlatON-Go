# -*- coding: utf-8 -*-
'''
@Description: rpc用例
'''
import allure
import pytest
from common.connect import connect_web3
from hexbytes import HexBytes

dataDir = "/home/platon/"
startApi = "eth,web3,net,txpool,platon,admin,personal"

w3 = None
ws = None
id = ""
host = ""
http_url = ""
http_port = "6789"
p2p_prot = "16789"
ws_url = ""
ws_port = "5789"
chainid = 100
amount = ""
period = ""
validatorMode = ""

@pytest.fixture()
def setNodeInfo(global_test_env):
    collusion_list = global_test_env.consensus_node_list
    if len(collusion_list) > 0:
        try:
            global w3
            global ws
            global id
            global host
            global http_url
            global http_port
            global p2p_prot
            global ws_url
            global ws_port
            global chainid
            global amount
            global period
            global validatorMode
            global dataDir

            test_node = collusion_list[0]
            id = test_node.node_id

            host = test_node.host
            http_url = test_node.url
            http_port = test_node.rpc_port
            p2p_prot = test_node.p2p_port
            ws_port = test_node.wsport
            if ws_port == "":
                ws_port = "5789"
            ws_url = test_node.wsurl
            if ws_url == "":
                ws_url = "ws://" + host + ":" + ws_port
            chainid = global_test_env.chain_id

            w3 = test_node.web3
        #    ws = test_node.ws_web3
            amount = global_test_env.amount
            period = global_test_env.period
            validatorMode = global_test_env.validatorMode

            dataDir = dataDir + test_node.cfg.deploy_path + "/node-" + p2p_prot + "/data"

            # start websocket rpc service
            if w3 != None:
                try:
                    assert True == w3.admin.startWS(host, ws_port, "*", startApi)
                    print("\n启动websocket rpc服务成功================:{}".format(host + ":" + str(ws_port)))
                    ws = connect_web3(ws_url, chainid)
                except Exception as e:
                    print("\n启动websocket rpc服务失败, error message:{}==================".format(e))
        except Exception as e:
            print("setNodeInfo error:{}>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>".format(e))
            w3 = None
            ws = None

@allure.title("获取进程保存数据的目录:admin.datadir")
@pytest.mark.P1
def test_admin_datadir(setNodeInfo):
    if w3 != None:
        try:
            assert w3.admin.datadir == dataDir
            print("\n当前进程保存数据目录datadir=================:{}".format(w3.admin.datadir))
        except Exception as e:
            print("\n error message:{}".format(e))

@allure.title("获取程序的版本号和签名:admin.getProgramVersion()")
@pytest.mark.P1
@pytest.mark.compatibility
def test_admin_getProgramVersion(setNodeInfo):
    if w3 != None:
        msg = w3.admin.getProgramVersion()
        ProgramVersionSign = msg["Sign"]
        ProgramVersion = msg["Version"]
        try:
            assert len(ProgramVersionSign) == 132
            assert ProgramVersion >= 1794
            print("\n获取当前程序的版本号:【{}】, 版本签名:【{}】成功".format(ProgramVersion, ProgramVersionSign))
        except Exception as e:
            print("\nerror message:{}".format(e))


@allure.title("获取零知识证明信息:admin.getSchnorrNIZKProve()")
@pytest.mark.P1
@pytest.mark.compatibility
def test_admin_getSchnorrNIZKProve(setNodeInfo):
    if w3 != None:
        blsproof = w3.admin.getSchnorrNIZKProve()
        assert len(blsproof) == 128
        print("\n获取零知识证明信息成功:{}".format(blsproof))


@allure.title("校验节点信息:admin.nodeInfo")
@pytest.mark.P1
@pytest.mark.compatibility
def test_admin_nodeInfo(setNodeInfo):
    if w3 != None:
        try:
            genHash = HexBytes(w3.platon.getBlock(0)["hash"]).hex()
            nodeInfo = w3.admin.nodeInfo

            # config
            config = nodeInfo["protocols"]["platon"]["config"]

            # node id
            assert id == nodeInfo["id"]
            # listen port
            assert p2p_prot == str(nodeInfo["ports"]["listener"])
            # discovery port
            assert p2p_prot == str(nodeInfo["ports"]["discovery"])

            assert amount == config["cbft"]["amount"]
            assert period == config["cbft"]["period"]
            assert validatorMode == config["cbft"]["validatorMode"]

            assert chainid == config["chainId"]
            assert genHash == nodeInfo["protocols"]["eth"]["genesis"]

            print("\n校验节点信息成功,节点信息：{}".format(nodeInfo))
        except Exception as e:
            print("\nerror:{}".format(e))


@allure.title("和本节点的连接信息:admin.peers")
@pytest.mark.P1
@pytest.mark.compatibility
def test_admin_peers(setNodeInfo):
    if w3 != None:
        try:
            print("\n本节点连接信息为:{}".format(w3.admin.peers))
        except Exception as e:
            print("\nerror message:{}".format(e))


# @allure.title("添加节点的连接:admin.addpeer()")
# @pytest.mark.P0
# def test_admin_addpeer():
#    test_init_diff_genesis_join_chain()
#    test_init_same_genesis_join_chain()

@allure.title("导出区块数据:admin.exportChain()")
@pytest.mark.P1
def test_admin_exportChain(setNodeInfo):
    if w3 != None:
        filePath = dataDir + "chainData.txt"
        try:
            assert True == w3.admin.exportChain(filePath)
            print("\n区块数据导出成功,文件路径:【{}】".format(filePath))
        except Exception as e:
            print("\n区块数据导出失败,error message:【{}】".format(e))


@allure.title("导入区块数据:admin.importChain()")
@pytest.mark.P1
def test_admin_importChain(setNodeInfo):
    if w3 != None:
        try:
            filePath = dataDir + "chainData.txt"
            assert True == w3.admin.importChain(filePath), "区块数据导入失败！"
            print("\n区块数据导入成功,文件路径:【{}】".format(filePath))
        except Exception as e:
            print("\n区块数据导入失败,error message:【{}】".format(e))


@allure.title("移除peer连接:admin.removePeer()")
@pytest.mark.P1
def test_admin_removePeer(setNodeInfo):
    if w3 != None:
        peers = w3.admin.peers
        if len(peers) > 0:
            node_url = "enode://" + peers[0]["id"] + "@" + peers[0]["network"]["remoteAddress"]
            assert True == w3.admin.removePeer(node_url)
            print("\n移除peer成功:{}".format(node_url))


@allure.title("停止websocket rpc服务:admin.stopWS()")
@pytest.fixture()
def admin_stopWS(setNodeInfo):
    if ws != None:
        try:
            assert True == ws.admin.stopWS()
            print("\n停止websocket rpc服务成功:{}".format(ws_url))
        except Exception as e:
            print("\n停止websocket rpc服务失败, error message:{}".format(e))


@allure.title("启动websocket rpc服务:admin.startWS()")
@pytest.mark.P1
@pytest.mark.compatibility
def test_admin_startWS(admin_stopWS):
    if w3 != None:
        try:
            assert True == w3.admin.startWS(host, ws_port, "*", startApi)
            print("\n启动websocket rpc服务成功:{}".format(ws_url))
        except Exception as e:
            print("\n启动websocket rpc服务失败, error message:{}".format(e))


@allure.title("停止http rpc服务:admin.stopRPC()")
@pytest.fixture()
def admin_stopRPC(setNodeInfo):
    if w3 != None:
        try:
            assert True == w3.admin.stopRPC()
            print("\n停止http rpc服务成功:{}".format(http_url))
        except Exception as e:
            print("\n停止http rpc服务失败, error message:{}".format(e))


@allure.title("启动http rpc服务:admin.startRPC()")
@pytest.mark.P0
@pytest.mark.compatibility
def test_admin_startRPC(admin_stopRPC):
    if ws != None:
        try:
            assert True == ws.admin.startRPC(host, http_port, "*", startApi)
            print("\n启动http rpc服务成功:{}".format(http_url))
        except Exception as e:
            print("\n启动http rpc服务失败,error message:{}".format(e))

@pytest.mark.P0
def test_deploy_all(global_test_env):
    global_test_env.deploy_all()

if __name__ == '__main__':
    pytest.main(['-v', 'test_rpc_admin.py'])