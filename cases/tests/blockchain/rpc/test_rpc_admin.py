# -*- coding: utf-8 -*-
'''
@Description: rpc用例
'''

import allure
import pytest
from hexbytes import HexBytes

startApi = "eth,web3,net,txpool,platon,admin,personal"


@allure.title("Get current process datadir")
@pytest.mark.P1
def test_admin_datadir(global_running_env):
    node = global_running_env.get_rand_node()
    dataDir = node.remote_data_dir
    assert node.admin.datadir == dataDir


@allure.title("Get program version")
@pytest.mark.P1
def test_admin_getProgramVersion(global_running_env):
    node = global_running_env.get_rand_node()
    msg = node.admin.getProgramVersion()
    ProgramVersionSign = msg["Sign"]
    ProgramVersion = msg["Version"]
    assert len(ProgramVersionSign) == 132
    assert ProgramVersion >= 1794


@allure.title("Get schnorrNIZKProve")
@pytest.mark.P1
def test_admin_getSchnorrNIZKProve(global_running_env):
    node = global_running_env.get_rand_node()
    blsproof = node.admin.getSchnorrNIZKProve()
    assert len(blsproof) == 128


@allure.title("get node info")
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


@allure.title("get node peers")
@pytest.mark.P1
def test_admin_peers(global_running_env):
    node = global_running_env.get_rand_node()
    lenPeers = len(node.admin.peers)
    assert lenPeers >= 0


@allure.title("export chain")
@pytest.mark.P1
@pytest.fixture()
def test_admin_exportChain(global_running_env):
    node = global_running_env.get_rand_node()
    filePath = node.admin.datadir + "chainData.txt"
    assert True == node.admin.exportChain(filePath)
    yield node


@allure.title("import chain")
@pytest.mark.P1
def test_admin_importChain(test_admin_exportChain):
    filePath = test_admin_exportChain.admin.datadir + "chainData.txt"
    assert True == test_admin_exportChain.admin.importChain(filePath), "import chain failed！"


@allure.title("remove peer")
@pytest.mark.P1
def test_admin_removePeer(global_running_env):
    node = global_running_env.get_rand_node()
    peers = node.admin.peers
    if len(peers) > 0:
        node_url = "enode://" + peers[0]["id"] + "@" + peers[0]["network"]["remoteAddress"]
        assert True == node.admin.removePeer(node_url)


@allure.title("stop websocket rpc service")
@pytest.fixture()
def admin_stopWS(global_running_env):
    node = global_running_env.get_rand_node()
    try:
        ws = node.ws_web3
        assert True == ws.admin.stopWS()
    except Exception as e:
        print("websocket service not started===================>")

    yield node


@allure.title("Start websocket rpc service")
@pytest.mark.P1
@pytest.fixture()
def test_admin_startWS(admin_stopWS):
    node = admin_stopWS
    if None == node.wsport:
        node.wsport = 5789
    if None == node.wsurl:
        node.wsurl = "ws://" + str(node.host) + ":" + str(node.wsport)
    assert True == node.admin.startWS(node.host, int(node.wsport), "*", startApi)

    ws = node.ws_web3
    assert ws.eth.blockNumber >= 0


@allure.title("stop http rpc service")
@pytest.fixture()
def admin_stopRPC(global_running_env):
    node = global_running_env.get_rand_node()
    try:
        ws = node.ws_web3
        assert True == ws.admin.stopRPC()
    except Exception as e:
        pass

    yield node


'''
@allure.title("start http rpc service")
@pytest.mark.P0
def test_admin_startRPC(admin_stopRPC):
    node = admin_stopRPC
    ws = node.ws_web3
    assert True == ws.admin.startRPC(admin_stopRPC.host, int(admin_stopRPC.rpc_port))
'''

if __name__ == '__main__':
    pytest.main(['-v', 'test_rpc_admin.py'])
