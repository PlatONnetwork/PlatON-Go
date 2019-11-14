import socket
import time
from copy import copy

import allure

from common.connect import connect_web3, run_ssh


def isConnection(url, port):
    """
    检测是否开启了某个端口
    :param url: ip地址
    :param port: 端口号
    :return: 是否连接成功
    """
    s = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
    try:
        s.connect((url, port))
        s.shutdown(2)
        return True
    except:
        return False


def isWsConnection(url):
    ws_web3 = connect_web3(url)
    try:
        return ws_web3.isConnected()
    except:
        return False


def file_is_exist(ssh, path, file_name):
    cmd_list = run_ssh(
        ssh, "find {} -name {}".format(path, file_name))
    if len(cmd_list) > 0:
        return file_name in cmd_list[0]
    else:
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

    def test_rpc_api(self, global_test_env):
        env = global_test_env
        node = env.get_rand_node()
        modules = node.web3.manager.request_blocking("rpc_modules", [])
        api_method = "debug"
        assert modules.get(api_method) is not None

    def test_open_ws_connection(self, global_test_env):
        env = global_test_env
        node = env.get_rand_node()
        assert isWsConnection(node.wsurl)

    def test_ws_api(self, global_test_env):
        env = global_test_env
        node = env.get_rand_node()
        modules = node.web3.manager.request_blocking("rpc_modules", [])
        api_method = "debug"
        assert modules.get(api_method) is not None

    def test_enable_ipc(self, global_test_env):
        env = global_test_env
        node = env.get_rand_node()
        assert file_is_exist(node.ssh, node.remote_data_dir, "platon.ipc")

    def test_disable_ipc(self, global_test_env):
        test_node = append_cmd_restart(global_test_env, "--ipcdisable")
        assert bool(1 - file_is_exist(test_node.ssh, test_node.remote_data_dir, "platon.ipc"))

    def test_enable_ipc_config_name(self, global_test_env):
        test_node = append_cmd_restart(global_test_env, "--ipcpath platon_test.ipc")
        assert file_is_exist(test_node.ssh, test_node.remote_data_dir, "platon_test.ipc")

    ## 待处理
    def test_open_bootnodes(self, global_test_env):
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
        time.sleep(10)
        node_peers = test_node.admin.peers
        assert len(node_peers) == 1
        assert node_peers[0]["id"] == collusion_node.node_id

    def test_open_p2p_connection(self, global_test_env):
        env = global_test_env
        node = env.get_rand_node()
        assert isConnection(node.host, int(node.p2p_port))

    def test_open_discovery(self, global_test_env):
        env = global_test_env
        node = env.get_rand_node()
        node_info = node.admin.nodeInfo
        discovery = node_info["ports"]["discovery"]
        assert discovery != 0

    def test_close_discovery(self, global_test_env):
        test_node = append_cmd_restart(global_test_env, "--nodiscover")
        node_info = test_node.admin.nodeInfo
        discovery = node_info["ports"]["discovery"]
        assert discovery == 0

    def test_open_pprof(self, global_test_env):
        test_node = global_test_env.get_rand_node()
        pprof = 6060
        test_node = append_cmd_restart(global_test_env,
                                       "--pprof --pprofaddr {} --pprofport {}".format(test_node.host, pprof), test_node)
        assert isConnection(test_node.host, pprof)

    def test_enable_trace(self, global_test_env):
        test_node = global_test_env.get_rand_node()
        test_node = append_cmd_restart(global_test_env, "--trace {}/tracefile".format(test_node.remote_node_path))
        assert file_is_exist(test_node.ssh, test_node.remote_node_path, "tracefile")

    def test_enable_cpufile(self, global_test_env):
        test_node = global_test_env.get_rand_node()
        test_node = append_cmd_restart(global_test_env, "--cpuprofile {}/cpufile".format(test_node.remote_node_path))
        assert file_is_exist(test_node.ssh, test_node.remote_node_path, "cpufile")

    @allure.title("测试开启指标监控功能")
    def test_enable_metrics(self, global_test_env):
        test_node = append_cmd_restart(global_test_env, "--metrics")
        time.sleep(10)
        metrics = test_node.debug.web3.manager.request_blocking("debug_metrics", [True])
        assert metrics.cbft.gauage.block.number > 0
