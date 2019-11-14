import time

import allure
import pytest
from client_sdk_python.eth import Eth
from eth_utils import is_integer

from common.log import log
from common.connect import run_ssh_cmd
from client_sdk_python.admin import Admin

#作用域设置为module，自动运行
from conf.settings import NODE_FILE
# from environment import t1est_env_impl


# py.test tests/cmd/account/t1est_account.py -s --nodeFile "deploy/4_node.yml" --accountFile "deploy/accounts.yml" --initChain --startAll
from environment import Node


class AccountEnv:
    __slots__ = ('remote_pwd_file', 'remote_account_address', 'remote_account_file', 'remote_key_file')

@pytest.fixture(scope='module',autouse=False)
def account_env(global_test_env)->(Node, AccountEnv):
    log.info("module account begin.................................")

    env = global_test_env
    node = env.get_rand_node()
    log.info("Node::::::::::::::::::::::::::::::{}".format(node))

    remote_pwd_file = node.remote_node_path + "/password.txt"
    node.upload_file("./deploy/keystore/password.txt", remote_pwd_file)

    remote_account_file = node.remote_keystore_dir +"/UTC--2019-10-15T10-27-31.520865283Z--c198603d3793c11e5362c8564a65d3880bae341b"
    node.upload_file("./deploy/keystore/UTC--2019-10-15T10-27-31.520865283Z--c198603d3793c11e5362c8564a65d3880bae341b", remote_account_file)

    remote_key_file = node.remote_keystore_dir + "/key.pri"
    node.upload_file("./deploy/key.pri", remote_key_file)

    account_env = AccountEnv()
    account_env.remote_pwd_file = remote_pwd_file
    account_env.remote_account_file = remote_account_file
    account_env.remote_key_file = remote_key_file
    account_env.remote_account_address = "c198603d3793c11e5362c8564a65d3880bae341b"

    yield node, account_env

    log.info("module account end.................................")
    # node.deleteRemoteFile(remote_pwd_file)
    # node.deleteRemoteFile(remote_pwd_file)
    # node.deleteRemoteFile(remote_pwd_file)


@allure.title("指定datadir和keystore路径，通过输入密码创建新账号")
@pytest.mark.P1
@pytest.mark.SYNC
def test_account_new(account_env):
    node, env = account_env
    returnList = run_ssh_cmd(node.ssh, "{} account list --datadir {}".format(node.remote_bin_file, node.remote_data_dir))
    oldCounts = len(returnList) - 1

    run_ssh_cmd(node.ssh, "{} account new --datadir {}  --keystore {}".format(node.remote_bin_file, node.remote_data_dir, node.remote_keystore_dir), "88888888", "88888888")
    time.sleep(0.2)
    returnList2 = run_ssh_cmd(node.ssh, "{} account list --datadir {}".format(node.remote_bin_file, node.remote_data_dir))
    newCounts = len(returnList2) - 1
    assert oldCounts + 1 == newCounts

@allure.title("指定datadir，在缺省的datadir/keystore下，通过输入密码创建新账号")
@pytest.mark.P1
@pytest.mark.SYNC
def test_account_new_defualt_keystore_dir(account_env):
    node, env = account_env
    returnList = run_ssh_cmd(node.ssh, "{} account list --datadir {}".format(node.remote_bin_file, node.remote_data_dir))
    oldCounts = len(returnList) - 1

    run_ssh_cmd(node.ssh, "{} account new --datadir {}".format(node.remote_bin_file, node.remote_data_dir), "88888888", "88888888")
    time.sleep(0.2)
    returnList2 = run_ssh_cmd(node.ssh, "{} account list --datadir {}".format(node.remote_bin_file, node.remote_data_dir))
    newCounts = len(returnList2) - 1
    assert oldCounts + 1 == newCounts


@allure.title("指定keystore下，通过输入密码创建新账号")
@pytest.mark.P1
@pytest.mark.SYNC
def test_account_new_keystore_dir(account_env):
    node, env = account_env
    returnList = run_ssh_cmd(node.ssh, "{} account list --keystore {}".format(node.remote_bin_file, node.remote_keystore_dir))
    oldCounts = len(returnList) - 1

    run_ssh_cmd(node.ssh, "{} account new --keystore {}".format(node.remote_bin_file, node.remote_keystore_dir), "88888888", "88888888")
    time.sleep(0.2)
    returnList2 = run_ssh_cmd(node.ssh, "{} account list --keystore {}".format(node.remote_bin_file, node.remote_keystore_dir))
    newCounts = len(returnList2) - 1
    assert oldCounts + 1 == newCounts


@allure.title("指定datadir和keystore路径，通过密码文件创建新账号")
@pytest.mark.P1
@pytest.mark.SYNC
def test_account_new_with_pwd_file(account_env):
    node, env = account_env
    returnList = run_ssh_cmd(node.ssh, "{} account list --datadir {}".format(node.remote_bin_file, node.remote_data_dir))
    oldCounts = len(returnList) - 1

    run_ssh_cmd(node.ssh, "{} account new --datadir {} --keystore {} --password {}".format(node.remote_bin_file,
                                                                                            node.remote_data_dir,
                                                                                            node.remote_keystore_dir,
                                                                                            env.remote_pwd_file))
    time.sleep(0.2)
    returnList2 = run_ssh_cmd(node.ssh, "{} account list --datadir {}".format(node.remote_bin_file, node.remote_data_dir))
    newCounts = len(returnList2) - 1
    assert oldCounts + 1 == newCounts


@allure.title("指定datadir，在缺省的datadir/keystore下，通过密码文件创建新账号")
@pytest.mark.P1
@pytest.mark.SYNC
def test_account_new_defualt_keystore_dir(account_env):
    node, env = account_env
    returnList = run_ssh_cmd(node.ssh, "{} account list --datadir {}".format(node.remote_bin_file, node.remote_data_dir))
    oldCounts = len(returnList) - 1

    run_ssh_cmd(node.ssh, "{} account new --datadir {} --password {}".format(node.remote_bin_file,
                                                                                           node.remote_data_dir,
                                                                                           env.remote_pwd_file))
    time.sleep(0.2)
    returnList2 = run_ssh_cmd(node.ssh, "{} account list --datadir {}".format(node.remote_bin_file, node.remote_data_dir))
    newCounts = len(returnList2) - 1
    assert oldCounts + 1 == newCounts

@allure.title("指定keystore下，通过输入密码创建新账号")
@pytest.mark.P1
@pytest.mark.SYNC
def test_account_new_with_pwd_file_just_keystore_dir(account_env):
    node, env = account_env
    returnList = run_ssh_cmd(node.ssh, "{} account list --keystore {}".format(node.remote_bin_file, node.remote_keystore_dir))
    oldCounts = len(returnList) - 1

    run_ssh_cmd(node.ssh, "{} account new --keystore {} --password {}".format(node.remote_bin_file, node.remote_keystore_dir, env.remote_pwd_file))
    time.sleep(0.2)
    returnList2 = run_ssh_cmd(node.ssh, "{} account list --keystore {}".format(node.remote_bin_file, node.remote_keystore_dir))
    newCounts = len(returnList2) - 1
    assert oldCounts + 1 == newCounts



@allure.title("修改账号密码，指定datadir")
@pytest.mark.P1
@pytest.mark.SYNC
def test_account_update_with_data_dir(account_env):
    node, env = account_env
    run_ssh_cmd(node.ssh, "{} account update {} --datadir {}".format(node.remote_bin_file, env.remote_account_address, node.remote_data_dir), "88888888", "88888888", "88888888")
    pass


@allure.title("修改账号密码，指定keystore")
@pytest.mark.P1
@pytest.mark.SYNC
def test_account_update_with_keystore_dir(account_env):
    node, env = account_env
    run_ssh_cmd(node.ssh, "{} account update {} --keystore {}".format(node.remote_bin_file, env.remote_account_address, node.remote_keystore_dir), "88888888", "88888888", "88888888")
    pass


@allure.title("导入账号，不指定密码文件，指定datadir")
@pytest.mark.P1
@pytest.mark.SYNC
def test_account_import(account_env):
    node, env = account_env

    returnList = run_ssh_cmd(node.ssh, "{} account list --datadir {}".format(node.remote_bin_file, node.remote_data_dir))
    oldCounts = len(returnList) - 1

    remote_key_file = node.remote_keystore_dir + "/key.pri"
    node.upload_file("./deploy/key.pri", remote_key_file)

    run_ssh_cmd(node.ssh, "{} account import {} --datadir {}".format(node.remote_bin_file, remote_key_file, node.remote_data_dir), "88888888", "88888888")
    time.sleep(0.2)
    returnList2 = run_ssh_cmd(node.ssh, "{} account list --datadir {}".format(node.remote_bin_file, node.remote_data_dir))

    newCounts = len(returnList2) - 1

    assert oldCounts + 1 == newCounts

@allure.title("导入账号，不指定密码文件，指定keystore")
@pytest.mark.P1
@pytest.mark.SYNC
def test_account_import_2(account_env):
    node, env = account_env

    returnList = run_ssh_cmd(node.ssh, "{} account list --keystore {}".format(node.remote_bin_file, node.remote_keystore_dir))
    oldCounts = len(returnList) - 1

    remote_key_file = node.remote_keystore_dir + "/key.pri_2"
    node.upload_file("./deploy/key.pri_2", remote_key_file)

    run_ssh_cmd(node.ssh, "{} account import {} --keystore {}".format(node.remote_bin_file, remote_key_file, node.remote_keystore_dir), "88888888", "88888888")
    time.sleep(0.2)
    returnList2 = run_ssh_cmd(node.ssh, "{} account list --keystore {}".format(node.remote_bin_file, node.remote_keystore_dir))

    newCounts = len(returnList2) - 1

    assert oldCounts + 1 == newCounts


@allure.title("导入账号，指定密码文件，指定datadir")
@pytest.mark.P1
@pytest.mark.SYNC
def test_account_import_3(account_env):
    node, env = account_env

    returnList = run_ssh_cmd(node.ssh, "{} account list --datadir {}".format(node.remote_bin_file, node.remote_data_dir))
    oldCounts = len(returnList) - 1

    remote_key_file = node.remote_keystore_dir + "/key.pri_3"
    node.upload_file("./deploy/key.pri_3", remote_key_file)

    run_ssh_cmd(node.ssh, "{} account import {} --datadir {} --password {}".format(node.remote_bin_file, remote_key_file, node.remote_data_dir, env.remote_pwd_file))
    time.sleep(0.2)
    returnList2 = run_ssh_cmd(node.ssh, "{} account list --datadir {}".format(node.remote_bin_file, node.remote_data_dir))

    newCounts = len(returnList2) - 1

    assert oldCounts + 1 == newCounts

@allure.title("导入账号，指定密码文件，指定keystore")
@pytest.mark.P1
@pytest.mark.SYNC
def test_account_import_4(account_env):
    node, env = account_env

    returnList = run_ssh_cmd(node.ssh, "{} account list --keystore {}".format(node.remote_bin_file, node.remote_keystore_dir))
    oldCounts = len(returnList) - 1

    remote_key_file = node.remote_keystore_dir + "/key.pri_4"
    node.upload_file("./deploy/key.pri_4", remote_key_file)

    run_ssh_cmd(node.ssh, "{} account import {} --keystore {}  --password {}".format(node.remote_bin_file, remote_key_file, node.remote_keystore_dir, env.remote_pwd_file))
    time.sleep(0.2)
    returnList2 = run_ssh_cmd(node.ssh, "{} account list --keystore {}".format(node.remote_bin_file, node.remote_keystore_dir))

    newCounts = len(returnList2) - 1

    assert oldCounts + 1 == newCounts


@allure.title("列出账号")
@pytest.mark.P1
@pytest.mark.SYNC
def test_account_list(account_env):
    node, env = account_env

    returnList1 = run_ssh_cmd(node.ssh, "{} account list --datadir {}".format(node.remote_bin_file, node.remote_data_dir))

    counts1 = len(returnList1) - 1

    returnList2 = run_ssh_cmd(node.ssh, "{} account list --keystore {}".format(node.remote_bin_file, node.remote_keystore_dir))
    counts2 = len(returnList2) - 1

    assert  counts1 == counts2


'''
platon attach http / ws
'''
def test_attach_http(account_env):
    node, env = account_env

    print("node.remote_bin_file:::", node.remote_bin_file)
    print("node.url:::", node.url)

    blockNumber = node.run_ssh("{} attach {} --exec platon.blockNumber".format(node.remote_bin_file, node.url))

    bn = int(blockNumber[0])

    assert is_integer(bn)
    assert bn > 0


'''
platon attach http / ws
'''
def test_copydb(global_test_env):
    globalEnv = global_test_env

    node = globalEnv.collusion_node_list[0]

    log.info("test copydb on host: {}".format(node.host))

    node.stop()

    # copy deploy data to bak
    bakremote_data_dir = node.remote_node_path + "/data_bak"

    run_ssh_cmd(node.ssh, "sudo -S -p '' cp -r {} {}".format(node.remote_data_dir, bakremote_data_dir), node.password)

    run_ssh_cmd(node.ssh, "sudo -S -p '' rm -rf {}/platon".format(node.remote_data_dir), node.password)
    #run_ssh_cmd(node.ssh, "sudo -S -p '' rm -rf {}/platon/chaindata".format(node.remote_data_dir), node.password)

    # re-init
    run_ssh_cmd(node.ssh, "sudo -S -p '' {} init {} --datadir {}".format(node.remote_bin_file, node.remote_genesis_file, node.remote_data_dir), node.password)

    time.sleep(10)

    # copyDb from bak
    run_ssh_cmd(node.ssh, "sudo -S -p '' {} copydb {}/platon/chaindata/ {}/platon/snapshotdb/ --datadir {}".format(node.remote_bin_file, bakremote_data_dir, bakremote_data_dir, node.remote_data_dir), node.password)
    time.sleep(10)

    node.start(False)

    time.sleep(5)

    blockNumber = node.run_ssh("{} attach {} --exec platon.blockNumber".format(node.remote_bin_file, node.url))

    for i in range(len(blockNumber)):
        print("序号：{}".format(i), "结果：{}".format(blockNumber[i]))

    bn = int(blockNumber[0])

    assert is_integer(bn)
    assert bn > 0
    #pass

def test_dump_block(global_test_env):
    globalEnv = global_test_env

    node = globalEnv.collusion_node_list[0]
    node.stop()

    # dump
    returnList = run_ssh_cmd(node.ssh, "sudo -S -p '' {} --datadir {} dump 0".format(node.remote_bin_file, node.remote_data_dir), node.password)

    node.start(False)

    assert len(returnList) > 0 and "root" in returnList[1]


def test_dump_config(global_test_env):
    globalEnv = global_test_env

    node = globalEnv.collusion_node_list[0]
    # dump
    returnList = run_ssh_cmd(node.ssh, "{} --nodekey {} --cbft.blskey {} dumpconfig".format(node.remote_bin_file, node.remote_nodekey_file, node.remote_blskey_file))
    assert returnList[0].strip()=='[Eth]'

def test_update_dumped_config(global_test_env):
    globalEnv = global_test_env

    node = globalEnv.collusion_node_list[0]
    # dump
    returnList = run_ssh_cmd(node.ssh, "{} --nodekey {} --cbft.blskey {} dumpconfig --networkid 1500".format(node.remote_bin_file, node.remote_nodekey_file, node.remote_blskey_file))

    assert returnList[1].strip()=='NetworkId = 1500'


def test_export_import_preimages(global_test_env):
    globalEnv = global_test_env

    node = globalEnv.collusion_node_list[0]
    node.stop()

    # dump
    exportList = run_ssh_cmd(node.ssh, "sudo -S -p '' {} export-preimages exportPreImage --datadir {}".format(node.remote_bin_file, node.remote_data_dir), node.password)
    for i in range(len(exportList)):
        log.info("序号：{}   结果：{}".format(i, exportList[i]))

    time.sleep(1)

    importList = run_ssh_cmd(node.ssh, "sudo -S -p '' {} import-preimages exportPreImage --datadir {}".format(node.remote_bin_file, node.remote_data_dir), node.password)
    node.start(False)

    for i in range(len(importList)):
        log.info("序号：{}   结果：{}".format(i, importList[i]))

    assert len(exportList) == 1
    assert len(importList) == 1


def test_license(global_test_env):
    globalEnv = global_test_env

    node = globalEnv.collusion_node_list[0]

    returnList = run_ssh_cmd(node.ssh, "{} license".format(node.remote_bin_file))
    # for i in range(len(returnList)):
    #     log.info("序号：{}   结果：{}".format(i, returnList[i]))

    assert returnList[0].strip()=="platon is free software: you can redistribute it and/or modify"

def test_version(global_test_env):
    globalEnv = global_test_env

    node = globalEnv.collusion_node_list[0]

    returnList = run_ssh_cmd(node.ssh, "{} version".format(node.remote_bin_file))
    # for i in range(len(returnList)):
    #     log.info("序号：{}   结果：{}".format(i, returnList[i]))

    assert returnList[0].strip()=="Platon"
    assert "Version:" in returnList[1]


def test_config(global_test_env):
    globalEnv = global_test_env

    node = globalEnv.collusion_node_list[0]
    node.stop()

    returnList = run_ssh_cmd(node.ssh, "sed -i 's/\"NetworkId\": 1/\"NetworkId\": 111/g' {}".format(node.remote_config_file))

    node.start(False)

    time.sleep(2)

    ret = node.admin.nodeInfo
    #print(ret)
    assert  ret["protocols"]["platon"]["network"] == 111



def no_t1est_removedb(global_test_env):
    globalEnv = global_test_env

    node = globalEnv.collusion_node_list[0]
    node.stop()

    returnList = run_ssh_cmd(node.ssh, "{} removedb --datadir {}".format(node.remote_bin_file, node.remote_data_dir, "y", "y"))
    for i in range(len(returnList)):
        log.info("序号：{}   结果：{}".format(i, returnList[i]))

    node.start(False)
    #assert returnList[0].strip()=="platon is free software: you can redistribute it and/or modify"