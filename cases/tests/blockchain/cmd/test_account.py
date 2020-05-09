import time

import allure
import pytest
from client_sdk_python.eth import Eth
from eth_utils import is_integer

from common.log import log
from common.connect import run_ssh_cmd
from client_sdk_python.admin import Admin

from conf.settings import NODE_FILE
# from environment import t1est_env_impl


# py.test tests/cmd/account/t1est_account.py -s --nodeFile "deploy/4_node.yml" --accountFile "deploy/accounts.yml" --initChain --startAll
from environment import Node


class AccountEnv:
    __slots__ = ('remote_pwd_file', 'remote_account_address', 'remote_account_file', 'remote_key_file')


@pytest.fixture(scope='module', autouse=False)
def account_env(global_test_env) -> (Node, AccountEnv):
    log.info("module account begin.................................")

    env = global_test_env
    node = env.get_rand_node()
    log.info("Node::::::::::::::::::::::::::::::{}".format(node))

    remote_pwd_file = node.remote_node_path + "/password.txt"
    node.upload_file("./deploy/keystore/password.txt", remote_pwd_file)

    remote_account_file = node.remote_keystore_dir + "/UTC--2019-10-15T10-27-31.520865283Z--c198603d3793c11e5362c8564a65d3880bae341b"
    node.upload_file("./deploy/keystore/UTC--2019-10-15T10-27-31.520865283Z--c198603d3793c11e5362c8564a65d3880bae341b", remote_account_file)

    remote_key_file = node.remote_keystore_dir + "/key.pri"
    node.upload_file("./deploy/key.pri", remote_key_file)

    account_env = AccountEnv()
    account_env.remote_pwd_file = remote_pwd_file
    account_env.remote_account_file = remote_account_file
    account_env.remote_key_file = remote_key_file
    account_env.remote_account_address = "lax1785psd0qs0g8p79j54mnewh0ndwcvqq6g23h8h"

    yield node, account_env

    log.info("module account end.................................")
    # node.deleteRemoteFile(remote_pwd_file)
    # node.deleteRemoteFile(remote_pwd_file)
    # node.deleteRemoteFile(remote_pwd_file)


@allure.title("Specify the datadir and keystore paths and create a new account by entering a password.")
@pytest.mark.P1
@pytest.mark.SYNC
def test_CMD_002(account_env):
    node, env = account_env
    return_list = run_ssh_cmd(node.ssh, "{} account list --datadir {}".format(node.remote_bin_file, node.remote_data_dir))
    old_counts = len(return_list) - 1

    run_ssh_cmd(node.ssh, "{} account new --datadir {}  --keystore {}".format(node.remote_bin_file, node.remote_data_dir, node.remote_keystore_dir), "88888888", "88888888")
    time.sleep(0.2)
    return_list2 = run_ssh_cmd(node.ssh, "{} account list --datadir {}".format(node.remote_bin_file, node.remote_data_dir))
    new_counts = len(return_list2) - 1
    assert old_counts + 1 == new_counts


@allure.title("Specify datadir. In the default datadir/keystore, create a new account by entering a password.")
@pytest.mark.P1
@pytest.mark.SYNC
def test_CMD_003(account_env):
    node, env = account_env
    return_list = run_ssh_cmd(node.ssh, "{} account list --datadir {}".format(node.remote_bin_file, node.remote_data_dir))
    old_counts = len(return_list) - 1

    run_ssh_cmd(node.ssh, "{} account new --datadir {}".format(node.remote_bin_file, node.remote_data_dir), "88888888", "88888888")
    time.sleep(0.2)
    return_list2 = run_ssh_cmd(node.ssh, "{} account list --datadir {}".format(node.remote_bin_file, node.remote_data_dir))
    new_counts = len(return_list2) - 1
    assert old_counts + 1 == new_counts


@allure.title("Specify a keystore and create a new account by entering a password.")
@pytest.mark.P1
@pytest.mark.SYNC
def test_CMD_004(account_env):
    node, env = account_env
    return_list = run_ssh_cmd(node.ssh, "{} account list --keystore {}".format(node.remote_bin_file, node.remote_keystore_dir))
    old_counts = len(return_list) - 1

    run_ssh_cmd(node.ssh, "{} account new --keystore {}".format(node.remote_bin_file, node.remote_keystore_dir), "88888888", "88888888")
    time.sleep(0.2)
    return_list2 = run_ssh_cmd(node.ssh, "{} account list --keystore {}".format(node.remote_bin_file, node.remote_keystore_dir))
    new_counts = len(return_list2) - 1
    assert old_counts + 1 == new_counts


@allure.title("Specify the datadir and keystore paths to create a new account with a password file.")
@pytest.mark.P1
@pytest.mark.SYNC
def test_CMD_001(account_env):
    node, env = account_env
    return_list = run_ssh_cmd(node.ssh, "{} account list --datadir {}".format(node.remote_bin_file, node.remote_data_dir))
    old_counts = len(return_list) - 1

    run_ssh_cmd(node.ssh, "{} account new --datadir {} --keystore {} --password {}".format(node.remote_bin_file,
                                                                                           node.remote_data_dir,
                                                                                           node.remote_keystore_dir,
                                                                                           env.remote_pwd_file))
    time.sleep(0.2)
    return_list2 = run_ssh_cmd(node.ssh, "{} account list --datadir {}".format(node.remote_bin_file, node.remote_data_dir))
    new_counts = len(return_list2) - 1
    assert old_counts + 1 == new_counts


@allure.title("Specify datadir, create a new account with the password file in the default datadir/keystore")
@pytest.mark.P1
@pytest.mark.SYNC
def test_CMD_003_2(account_env):
    node, env = account_env
    return_list = run_ssh_cmd(node.ssh, "{} account list --datadir {}".format(node.remote_bin_file, node.remote_data_dir))
    old_counts = len(return_list) - 1

    run_ssh_cmd(node.ssh, "{} account new --datadir {} --password {}".format(node.remote_bin_file,
                                                                             node.remote_data_dir,
                                                                             env.remote_pwd_file))
    time.sleep(0.2)
    return_list2 = run_ssh_cmd(node.ssh, "{} account list --datadir {}".format(node.remote_bin_file, node.remote_data_dir))
    new_counts = len(return_list2) - 1
    assert old_counts + 1 == new_counts


@allure.title("Specify a keystore and create a new account by entering a password.")
@pytest.mark.P1
@pytest.mark.SYNC
def test_CMD_004_2(account_env):
    node, env = account_env
    return_list = run_ssh_cmd(node.ssh, "{} account list --keystore {}".format(node.remote_bin_file, node.remote_keystore_dir))
    old_counts = len(return_list) - 1

    run_ssh_cmd(node.ssh, "{} account new --keystore {} --password {}".format(node.remote_bin_file, node.remote_keystore_dir, env.remote_pwd_file))
    time.sleep(0.2)
    return_list2 = run_ssh_cmd(node.ssh, "{} account list --keystore {}".format(node.remote_bin_file, node.remote_keystore_dir))
    new_counts = len(return_list2) - 1
    assert old_counts + 1 == new_counts


@allure.title("Change account password, specify datadir")
@pytest.mark.P1
@pytest.mark.SYNC
def test_CMD_005(account_env):
    node, env = account_env
    returnList = run_ssh_cmd(node.ssh, "{} account update {} --datadir {}".format(node.remote_bin_file, env.remote_account_address, node.remote_data_dir), "88888888", "88888888", "88888888")

    assert len(returnList) == 6
    assert returnList[5].strip() == "Repeat passphrase:"


@allure.title("Change account password, specify keystore")
@pytest.mark.P1
@pytest.mark.SYNC
def test_CMD_006(account_env):
    node, env = account_env
    returnList = run_ssh_cmd(node.ssh, "{} account update {} --keystore {}".format(node.remote_bin_file, env.remote_account_address, node.remote_keystore_dir), "88888888", "88888888", "88888888")

    assert len(returnList) == 6
    assert returnList[5].strip() == "Repeat passphrase:"


@allure.title("Import account, do not specify password file, specify datadir")
@pytest.mark.P1
@pytest.mark.SYNC
def test_CMD_007(account_env):
    node, env = account_env
    log.info(node.node_mark)
    return_list = run_ssh_cmd(node.ssh, "{} account list --datadir {}".format(node.remote_bin_file, node.remote_data_dir))
    old_counts = len(return_list) - 1

    remote_key_file = node.remote_keystore_dir + "/key.pri"
    node.upload_file("./deploy/key.pri", remote_key_file)

    run_ssh_cmd(node.ssh, "{} account import {} --datadir {}".format(node.remote_bin_file, remote_key_file, node.remote_data_dir), "88888888", "88888888")
    time.sleep(2)
    return_list2 = run_ssh_cmd(node.ssh, "{} account list --datadir {}".format(node.remote_bin_file, node.remote_data_dir))

    new_counts = len(return_list2) - 1

    assert old_counts + 1 == new_counts


@allure.title("Import account, do not specify password file, specify keystore")
@pytest.mark.P1
@pytest.mark.SYNC
def test_CMD_010_CMD_034(account_env):
    node, env = account_env

    return_list = run_ssh_cmd(node.ssh, "{} account list --keystore {}".format(node.remote_bin_file, node.remote_keystore_dir))
    old_counts = len(return_list) - 1

    remote_key_file = node.remote_keystore_dir + "/key.pri_2"
    node.upload_file("./deploy/key.pri_2", remote_key_file)

    run_ssh_cmd(node.ssh, "{} account import {} --keystore {}".format(node.remote_bin_file, remote_key_file, node.remote_keystore_dir), "88888888", "88888888")
    time.sleep(0.2)
    return_list2 = run_ssh_cmd(node.ssh, "{} account list --keystore {}".format(node.remote_bin_file, node.remote_keystore_dir))

    new_counts = len(return_list2) - 1

    assert old_counts + 1 == new_counts


@allure.title("Import account, specify password file, specify datadir")
@pytest.mark.P1
@pytest.mark.SYNC
def test_CMD_009(account_env):
    node, env = account_env

    return_list = run_ssh_cmd(node.ssh, "{} account list --datadir {}".format(node.remote_bin_file, node.remote_data_dir))
    old_counts = len(return_list) - 1

    remote_key_file = node.remote_keystore_dir + "/key.pri_3"
    node.upload_file("./deploy/key.pri_3", remote_key_file)

    run_ssh_cmd(node.ssh, "{} account import {} --datadir {} --password {}".format(node.remote_bin_file, remote_key_file, node.remote_data_dir, env.remote_pwd_file))
    time.sleep(0.2)
    return_list2 = run_ssh_cmd(node.ssh, "{} account list --datadir {}".format(node.remote_bin_file, node.remote_data_dir))

    new_counts = len(return_list2) - 1

    assert old_counts + 1 == new_counts


@allure.title("Import account, specify password file, specify keystore")
@pytest.mark.P1
@pytest.mark.SYNC
def test_CMD_008(account_env):
    node, env = account_env

    return_list = run_ssh_cmd(node.ssh, "{} account list --keystore {}".format(node.remote_bin_file, node.remote_keystore_dir))
    old_counts = len(return_list) - 1

    remote_key_file = node.remote_keystore_dir + "/key.pri_4"
    node.upload_file("./deploy/key.pri_4", remote_key_file)

    run_ssh_cmd(node.ssh, "{} account import {} --keystore {}  --password {}".format(node.remote_bin_file, remote_key_file, node.remote_keystore_dir, env.remote_pwd_file))
    time.sleep(0.2)
    return_list2 = run_ssh_cmd(node.ssh, "{} account list --keystore {}".format(node.remote_bin_file, node.remote_keystore_dir))

    new_counts = len(return_list2) - 1

    assert old_counts + 1 == new_counts


@allure.title("List account")
@pytest.mark.P1
@pytest.mark.SYNC
def test_CMD_011(account_env):
    node, env = account_env

    return_list1 = run_ssh_cmd(node.ssh, "{} account list --datadir {}".format(node.remote_bin_file, node.remote_data_dir))

    counts1 = len(return_list1) - 1

    return_list2 = run_ssh_cmd(node.ssh, "{} account list --keystore {}".format(node.remote_bin_file, node.remote_keystore_dir))
    counts2 = len(return_list2) - 1

    assert counts1 == counts2


'''
platon attach http / ws
'''


@allure.title("Connect the node and open the js interactive console")
@pytest.mark.P3
def test_CMD_015(account_env):
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
@allure.title("Copy chain data")
@pytest.mark.P3
def test_CMD_016(global_test_env):
    global_env = global_test_env

    node = global_env.consensus_node_list[0]

    log.info("test copydb on host: {}".format(node.host))

    node.stop()

    # copy deploy data to bak
    bakremote_data_dir = node.remote_node_path + "/data_bak"

    run_ssh_cmd(node.ssh, "sudo -S -p '' cp -r {} {}".format(node.remote_data_dir, bakremote_data_dir), node.password)

    run_ssh_cmd(node.ssh, "sudo -S -p '' rm -rf {}/platon".format(node.remote_data_dir), node.password)
    # run_ssh_cmd(node.ssh, "sudo -S -p '' rm -rf {}/platon/chaindata".format(node.remote_data_dir), node.password)

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
        print("Serial number：{}".format(i), "result：{}".format(blockNumber[i]))

    bn = int(blockNumber[0])

    assert is_integer(bn)
    assert bn > 0
    # pass


@allure.title("Analyze a specific block")
@pytest.mark.P3
def test_CMD_017(global_test_env):
    global_env = global_test_env

    node = global_env.consensus_node_list[0]
    node.stop()

    # dump
    return_list = run_ssh_cmd(node.ssh, "sudo -S -p '' {} --datadir {} dump 0".format(node.remote_bin_file, node.remote_data_dir), node.password)

    node.start(False)

    assert len(return_list) > 0 and "root" in return_list[1]


@allure.title("Display configuration values(you can view the default configuration information of the node)")
@pytest.mark.P3
def test_CMD_018(global_test_env):
    global_env = global_test_env

    node = global_env.consensus_node_list[0]
    # dump
    returnList = run_ssh_cmd(node.ssh, "{} --nodekey {} --cbft.blskey {} dumpconfig".format(node.remote_bin_file, node.remote_nodekey_file, node.remote_blskey_file))
    assert returnList[0].strip() == '[Eth]'


@allure.title("Modify the exported value when exporting")
@pytest.mark.P3
def test_CMD_019(global_test_env):
    global_env = global_test_env

    node = global_env.consensus_node_list[0]
    # dump
    return_list = run_ssh_cmd(node.ssh, "{} --nodekey {} --cbft.blskey {} dumpconfig --networkid 1500".format(node.remote_bin_file, node.remote_nodekey_file, node.remote_blskey_file))

    assert return_list[1].strip() == 'NetworkId = 1500'


@allure.title("Import blocks from the hash image file")
@pytest.mark.P3
def test_CMD_025(global_test_env):
    global_env = global_test_env

    node = global_env.consensus_node_list[0]
    node.stop()

    # dump
    export_list = run_ssh_cmd(node.ssh, "sudo -S -p '' {} export-preimages exportPreImage --datadir {}".format(node.remote_bin_file, node.remote_data_dir), node.password)
    for i in range(len(export_list)):
        log.info("Serial number：{}   result：{}".format(i, export_list[i]))

    time.sleep(1)

    import_list = run_ssh_cmd(node.ssh, "sudo -S -p '' {} import-preimages exportPreImage --datadir {}".format(node.remote_bin_file, node.remote_data_dir), node.password)
    node.start(False)

    for i in range(len(import_list)):
        log.info("Serial number：{}   result：{}".format(i, import_list[i]))

    assert len(export_list) == 1
    assert len(import_list) == 1


@allure.title("Display version information")
@pytest.mark.P3
def test_CMD_026(global_test_env):
    global_env = global_test_env

    node = global_env.consensus_node_list[0]

    return_list = run_ssh_cmd(node.ssh, "{} license".format(node.remote_bin_file))
    # for i in range(len(returnList)):
    #     log.info("Serial number：{}   result：{}".format(i, returnList[i]))

    assert return_list[0].strip() == "platon is free software: you can redistribute it and/or modify"


@allure.title("Display chain version")
@pytest.mark.P3
def test_CMD_029(global_test_env):
    global_env = global_test_env

    node = global_env.consensus_node_list[0]

    return_list = run_ssh_cmd(node.ssh, "{} version".format(node.remote_bin_file))
    # for i in range(len(returnList)):
    #     log.info("Serial number：{}   Result: {}".format(i, returnList[i]))

    assert return_list[0].strip() == "PlatON"
    assert "Version:" in return_list[1]


@allure.title("Load configuration file")
@pytest.mark.P3
def test_CMD_033_CMD_034(global_test_env):
    global_env = global_test_env

    node = global_env.consensus_node_list[0]
    node.stop()

    run_ssh_cmd(node.ssh, "sed -i 's/\"NetworkId\": 1/\"NetworkId\": 111/g' {}".format(node.remote_config_file))

    node.start(False)

    time.sleep(2)

    ret = node.admin.nodeInfo
    # print(ret)
    assert ret["protocols"]["platon"]["network"] == 111

# Todo: use case without assertion
# def no_t1est_removedb(global_test_env):
#     globalEnv = global_test_env
#
#     node = globalEnv.consensus_node_list[0]
#     node.stop()
#
#     returnList = run_ssh_cmd(node.ssh, "{} removedb --datadir {}".format(node.remote_bin_file, node.remote_data_dir, "y", "y"))
#     for i in range(len(returnList)):
#         log.info("Serial number:{}   result".format(i, returnList[i]))
#
#     node.start(False)
    #assert returnList[0].strip()=="platon is free software: you can redistribute it and/or modify"
