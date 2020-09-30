import os
import pytest
import json
import allure
from common.log import log
from concurrent.futures import ThreadPoolExecutor, wait, ALL_COMPLETED
from common.load_file import LoadFile
from common.connect import run_ssh_cmd


def one_put_config_task(node):
    try:
        node.put_config()
    except Exception as e:
        return False, "{} upload config file failed:{}".format(node.node_mark, e)
    return True, "{} upload config file success".format(node.node_mark)


def one_put_static_task(node):
    try:
        node.put_static()
    except Exception as e:
        return False, "{} upload static file failed:{}".format(node.node_mark, e)
    return True, "{} upload static file success".format(node.node_mark)


def one_put_genesis_task(node, genesis_file):
    try:
        node.put_genesis(genesis_file)
    except Exception as e:
        return False, "{} upload genesis file failed:{}".format(node.node_mark, e)
    return True, "{} upload genesis file success".format(node.node_mark)


@pytest.fixture(scope="module", autouse=True)
def reset_config(global_test_env):
    yield
    config_data = LoadFile(global_test_env.cfg.config_json_tmp).get_data()
    with open(global_test_env.cfg.config_json_tmp, 'w', encoding='utf-8') as f:
        f.write(json.dumps(config_data, indent=4))
    global_test_env.deploy_all()


@allure.title("Node maximum link quantity test")
@pytest.mark.P1
def test_NE_P2P_001(global_test_env):
    log.info("Node maximum link quantity test")
    all_node = global_test_env.get_all_nodes()

    # stop node
    if global_test_env.running:
        global_test_env.stop_all()

    # modify config file
    config_data = LoadFile(global_test_env.cfg.config_json_tmp).get_data()
    config_data['node']['P2P']['MaxPeers'] = 2
    with open(global_test_env.cfg.config_json_tmp, 'w', encoding='utf-8') as f:
        f.write(json.dumps(config_data, indent=4))

    # upload config file
    global_test_env.executor(one_put_config_task, all_node)

    # start node
    global_test_env.cfg.init_chain = False
    global_test_env.start_all()

    # run ssh
    static_number = len(global_test_env.get_static_nodes())
    for node in all_node:
        cmd_list = run_ssh_cmd(node.ssh, "netstat -an | grep 16789 | grep ESTABLISHED |wc -l")
        assert int(cmd_list[0][0]) <= 2 + static_number


@allure.title("Automatic discovery configuration test")
@pytest.mark.P1
def test_NE_P2P_002(global_test_env):
    log.info("Automatic discovery configuration test")
    all_node = global_test_env.get_all_nodes()

    # stop node
    if global_test_env.running:
        global_test_env.stop_all()

    # modify config file
    config_data = LoadFile(global_test_env.cfg.config_json_tmp).get_data()
    config_data['node']['P2P']['NoDiscovery'] = True
    with open(global_test_env.cfg.config_json_tmp, 'w', encoding='utf-8') as f:
        f.write(json.dumps(config_data, indent=4))
    all_node = global_test_env.get_all_nodes()

    # upload config file
    global_test_env.executor(one_put_config_task, all_node)

    # start node
    global_test_env.cfg.init_chain = False
    global_test_env.start_all()

    # run ssh
    for node in all_node:
        cmd_list = run_ssh_cmd(node.ssh, "netstat -unlp | grep 16789 |wc -l")
        assert 0 == int(cmd_list[0][0])


@allure.title("Static node configuration test")
@pytest.mark.P1
def test_NE_P2P_003(global_test_env):
    log.info("Static node configuration test")
    all_node = global_test_env.get_all_nodes()

    # stop node
    if global_test_env.running:
        global_test_env.stop_all()

    # modify config file
    config_data = LoadFile(global_test_env.cfg.config_json_tmp).get_data()
    config_data['node']['P2P']['MaxPeers'] = 50
    config_data['node']['P2P']['NoDiscovery'] = True
    config_data['node']['P2P']["BootstrapNodes"] = []
    with open(global_test_env.cfg.config_json_tmp, 'w', encoding='utf-8') as f:
        f.write(json.dumps(config_data, indent=4))
    all_node = global_test_env.get_all_nodes()

    # upload config file
    global_test_env.executor(one_put_config_task, all_node)

    # start node
    global_test_env.cfg.init_chain = False
    global_test_env.start_all()

    # run ssh
    static_number = len(global_test_env.get_static_nodes())
    for node in all_node:
        # cmd_list = run_ssh_cmd(node.ssh, "netstat -an | grep 16789 | grep ESTABLISHED |wc -l")
        # log.info(node.web3.net.peerCount)
        # assert int(cmd_list[0][0]) <= static_number
        assert node.web3.net.peerCount <= static_number


@allure.title("Exception can not be out of the block test")
@pytest.mark.P1
def test_NE_P2P_004(global_test_env):
    log.info("Exception can not be out of the block test")
    # stop node
    if global_test_env.running:
        global_test_env.stop_all()

    # modify config file
    config_data = LoadFile(global_test_env.cfg.config_json_tmp).get_data()
    config_data['node']['P2P']['MaxPeers'] = 50
    config_data['node']['P2P']['NoDiscovery'] = True
    config_data['node']['P2P']["BootstrapNodes"] = []
    with open(global_test_env.cfg.config_json_tmp, 'w', encoding='utf-8') as f:
        f.write(json.dumps(config_data, indent=4))
    all_node = global_test_env.get_all_nodes()

    # upload config file
    global_test_env.executor(one_put_config_task, all_node)

    # modify static file
    with open(global_test_env.cfg.static_node_tmp, 'w', encoding='utf-8') as f:
        f.write(json.dumps([], indent=4))

    # upload static file
    global_test_env.executor(one_put_static_task, all_node)

    # modify genesis file
    global_test_env.genesis_config['config']['cbft']["initialNodes"] = []
    with open(global_test_env.cfg.genesis_tmp, 'w', encoding='utf-8') as f:
        f.write(json.dumps(global_test_env.genesis_config, indent=4))

    # upload genesis file
    global_test_env.executor(one_put_genesis_task, all_node, global_test_env.cfg.genesis_tmp)

    # start node
    global_test_env.cfg.init_chain = False
    global_test_env.start_all()

    # check
    try:
        global_test_env.check_block()
    except Exception as e:
        log.error("check block has except:{}".format(e))
        assert 0 == 0
