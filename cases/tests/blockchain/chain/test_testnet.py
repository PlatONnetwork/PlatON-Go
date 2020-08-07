import time
from copy import copy

import pytest
from client_sdk_python.eth import PlatON
from hexbytes import HexBytes

from common.log import log

@pytest.mark.skip("Test case process is random and needs to be executed multiple times manually")
def test_testnet_fast(global_test_env):
    test_node = copy(global_test_env.get_a_normal_node())
    test_node.clean()
    new_cfg = copy(global_test_env.cfg)
    new_cfg.init_chain = False
    new_cfg.append_cmd = "--testnet"
    new_cfg.syncmode = "fast"
    test_node.cfg = new_cfg
    log.info("start deploy {}".format(test_node.node_mark))
    log.info("is init:{}".format(test_node.cfg.init_chain))
    test_node.deploy_me(genesis_file=None)
    log.info("deploy end")
    time.sleep(5)
    assert test_node.web3.net.peerCount > 1
    time.sleep(10)
    t = 2000
    while t:
        print(test_node.block_number)
        time.sleep(10)
        t -= 10
    assert test_node.block_number >= 1000000

@pytest.mark.skip("Test case process is random and needs to be executed multiple times manually")
def test_testnet_full(global_test_env):
    test_node = copy(global_test_env.get_a_normal_node())
    test_node.clean()
    new_cfg = copy(global_test_env.cfg)
    new_cfg.init_chain = False
    new_cfg.append_cmd = "--testnet"
    test_node.cfg = new_cfg
    log.info("start deploy {}".format(test_node.node_mark))
    log.info("is init:{}".format(test_node.cfg.init_chain))
    test_node.deploy_me(genesis_file=None)
    log.info("deploy end")
    time.sleep(5)
    assert test_node.web3.net.peerCount > 1
    time.sleep(10)
    t = 18000
    while t:
        print(test_node.block_number)
        time.sleep(10)
        t -= 10
    assert test_node.block_number >= 1000000


def sendTransaction(w3, nonce, from_address):
    # print(nonce)
    to_address = w3.toChecksumAddress(
        "0x54a7a3c6822eb222c53F76443772a60b0f9A8bab")
    tmp_from_address = w3.toChecksumAddress(from_address)
    platon = PlatON(w3)
    transaction_dict = {
        "to": to_address,
        "gasPrice": platon.gasPrice,
        "gas": 21000,
        "nonce": nonce,
        "data": "",
        "chainId": 100,
        "value": 1000,
        'from': tmp_from_address,
    }

    signedTransactionDict = platon.account.signTransaction(
        transaction_dict, "a689f0879f53710e9e0c1025af410a530d6381eebb5916773195326e123b822b"
    )

    # log.debug("signedTransactionDict:::::::{}，nonce::::::::::{}".format(signedTransactionDict, nonce))

    data = signedTransactionDict.rawTransaction
    return HexBytes(platon.sendRawTransaction(data)).hex()