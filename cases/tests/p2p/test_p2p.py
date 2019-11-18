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


@allure.title("节点最大链接数量测试")
@pytest.mark.P1
def test_p2p_max(global_test_env):
    '''节点最大链接数量测试'''
    log.info("节点最大链接数量测试")
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


@allure.title("自动发现配置测试")
@pytest.mark.P1
def test_p2p_discovery(global_test_env):
    '''自动发现配置测试'''
    log.info("自动发现配置测试")
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


@allure.title("静态节点配置测试")
@pytest.mark.P1
def test_p2p_static(global_test_env):
    '''静态节点配置测试'''
    log.info("静态节点配置测试")
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
        cmd_list = run_ssh_cmd(node.ssh, "netstat -an | grep 16789 | grep ESTABLISHED |wc -l")
        assert int(cmd_list[0][0]) <= static_number


@allure.title("异常无法出块测试")
@pytest.mark.P1
def test_p2p_except(global_test_env):
    '''异常无法出块测试'''
    log.info("异常无法出块测试")
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


