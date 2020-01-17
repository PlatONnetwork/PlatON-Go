import socket
import time
from copy import copy
import pytest
import allure

from common.connect import connect_web3, run_ssh


def isConnection(ip, port):
    """
    Checks whether the specified port is open
    :param ip: ip address
    :param port: port
    :return: boole
    """
    s = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
    try:
        s.connect((ip, port))
        s.shutdown(2)
        return True
    except BaseException:
        return False


def file_is_exist(ssh, path, file_name):
    cmd_list = run_ssh(
        ssh, "find {} -name {}".format(path, file_name))
    if len(cmd_list) > 0:
        return file_name in cmd_list[0]
    return False


def append_cmd_restart(global_test_env, cmd, node=None):
    if node is None:
        node = global_test_env.get_rand_node()
    test_node = copy(node)
    test_node.clean()
    new_cfg = copy(global_test_env.cfg)
    new_cfg.append_cmd = cmd
    test_node.cfg = new_cfg
    test_node.deploy_me(genesis_file=new_cfg.genesis_tmp)
    return test_node


class TestStartParam(object):

    @pytest.mark.compatibility
    @allure.title("Test access rpcapi")
    @pytest.mark.P3
    def test_CMD_077(self, global_test_env):
        env = global_test_env
        node = env.get_rand_node()
        modules = node.web3.manager.request_blocking("rpc_modules", [])
        api_method = "debug"
        assert modules.get(api_method) is not None

    @pytest.mark.compatibility
    @allure.title("Test to enable ws function")
    @pytest.mark.P3
    def test_CMD_078_CMD_081(self, global_test_env):
        env = global_test_env
        node = env.get_rand_node()
        ws_url = "ws://{}:".format(node.host)
        ws_port = node.wsport
        if ws_port is None:
            ws_port = 16000
            node.wsport = ws_port
            node.wsurl = "{}{}".format(ws_url, ws_port)
            append_cmd_restart(global_test_env, None, node)
        assert node.ws_web3.isConnected()

    @allure.title("Test to enable wsapi function")
    @pytest.mark.P3
    def test_CMD_082(self, global_test_env):
        env = global_test_env
        node = env.get_rand_node()
        if node.wsport is None:
            node.wsport = 16000
            node.wsurl = "ws://{}:{}".format(node.host, node.wsport)
            append_cmd_restart(global_test_env, None, node)
        modules = node.ws_web3.manager.request_blocking("rpc_modules", [])
        api_method = "debug"
        assert modules.get(api_method) is not None

    @allure.title("Test to enable ipc function")
    @pytest.mark.P3
    @pytest.mark.compatibility
    def test_enable_ipc(self, global_test_env):
        env = global_test_env
        node = env.get_rand_node()
        assert file_is_exist(node.ssh, node.remote_data_dir, "platon.ipc")

    @allure.title("Test off ipc function")
    @pytest.mark.P3
    def test_CMD_085(self, global_test_env):
        test_node = append_cmd_restart(global_test_env, "--ipcdisable")
        assert bool(1 - file_is_exist(test_node.ssh, test_node.remote_data_dir, "platon.ipc"))

    @allure.title("Test configuration ipc file name")
    @pytest.mark.P3
    def test_CMD_086(self, global_test_env):
        test_node = append_cmd_restart(global_test_env, "--ipcpath platon_test.ipc")
        time.sleep(10)
        assert file_is_exist(test_node.ssh, test_node.remote_data_dir, "platon_test.ipc")

    @pytest.mark.compatibility
    @allure.title("Test enable seed node")
    @pytest.mark.P3
    def test_CMD_089(self, global_test_env):
        global_test_env.deploy_all()
        env = global_test_env
        normal_node = env.get_a_normal_node()
        collusion_node = env.get_rand_node()
        new_cfg = copy(env.cfg)
        test_node = copy(normal_node)
        test_node.clean()
        new_cfg.is_need_static = False
        new_cfg.append_cmd = "--bootnodes \"{}\"".format(collusion_node.enode)
        test_node.cfg = new_cfg
        test_node.deploy_me(genesis_file=new_cfg.genesis_tmp)
        time.sleep(20)
        node_peers = test_node.admin.peers
        assert len(node_peers) == 1
        assert node_peers[0]["id"] == collusion_node.node_id

    @pytest.mark.compatibility
    @allure.title("Test open p2p port")
    @pytest.mark.P3
    def test_CMD_090(self, global_test_env):
        env = global_test_env
        node = env.get_rand_node()
        assert isConnection(node.host, int(node.p2p_port))

    @pytest.mark.compatibility
    @allure.title("Test to enable the discovery function")
    @pytest.mark.P3
    def test_CMD_097(self, global_test_env):
        env = global_test_env
        node = env.get_rand_node()
        node_info = node.admin.nodeInfo
        discovery = node_info["ports"]["discovery"]
        assert discovery != 0

    @allure.title("Test off the discovery function")
    @pytest.mark.P3
    def test_CMD_098(self, global_test_env):
        test_node = append_cmd_restart(global_test_env, "--nodiscover")
        node_info = test_node.admin.nodeInfo
        discovery = node_info["ports"]["discovery"]
        assert discovery == 0

    @allure.title("Test to enable pprof function")
    @pytest.mark.P3
    def test_CMD_115(self, global_test_env):
        test_node = global_test_env.get_rand_node()
        pprof = 6060
        test_node = append_cmd_restart(global_test_env,
                                       "--pprof --pprofaddr {} --pprofport {}".format(test_node.host, pprof), test_node)
        assert isConnection(test_node.host, pprof)

    @allure.title("Test to enable trace information file output")
    @pytest.mark.P3
    def test_CMD_119(self, global_test_env):
        test_node = global_test_env.get_rand_node()
        append_cmd_restart(global_test_env, "--trace {}/tracefile".format(test_node.remote_node_path), test_node)
        time.sleep(10)
        assert file_is_exist(test_node.ssh, test_node.remote_node_path, "tracefile")

    @allure.title("Test open output cpufile content")
    @pytest.mark.P3
    def test_CMD_118(self, global_test_env):
        test_node = global_test_env.get_rand_node()
        append_cmd_restart(global_test_env, "--cpuprofile {}/cpufile".format(test_node.remote_node_path), test_node)
        time.sleep(10)
        assert file_is_exist(test_node.ssh, test_node.remote_node_path, "cpufile")

    @allure.title("Test open indicator monitoring function")
    @pytest.mark.P3
    def test_CMD_121(self, global_test_env):
        test_node = append_cmd_restart(global_test_env, "--metrics")
        time.sleep(10)
        metrics = test_node.debug.web3.manager.request_blocking("debug_metrics", [True])
        assert metrics.cbft.gauage.block.number > 0
