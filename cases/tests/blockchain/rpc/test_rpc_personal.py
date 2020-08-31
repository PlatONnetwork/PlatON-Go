# -*- coding: utf-8 -*-
'''
@Description: rpc cases
'''
import time

import allure
import pytest
from client_sdk_python import Web3
from client_sdk_python.eth import Eth
from hexbytes import HexBytes
password = "88888888"
to_address = "lax1zy3rg32u8ge06yfrp3pw00xdf2zwqgqshrlnax"
g_txHash = None


@allure.title("List all account addresses")
@pytest.mark.P1
@pytest.mark.compatibility
def test_personal_listAccounts(global_running_env):
    node = global_running_env.get_rand_node()
    assert len(node.personal.listAccounts) >= 0


@allure.title("List all wallet information")
@pytest.mark.P1
@pytest.mark.compatibility
def test_personal_listWallets(global_running_env):
    node = global_running_env.get_rand_node()
    listWallets = node.personal.listWallets
    assert len(listWallets) >= 0


@allure.title("Create a new account")
@pytest.mark.P1
@pytest.mark.compatibility
@pytest.fixture()
def test_personal_newAccount(global_running_env):
    node = global_running_env.get_rand_node()
    before = len(node.eth.accounts)
    new_account = node.personal.newAccount(password)
    time.sleep(2)
    to_account = Web3.toChecksumAddress(new_account)
    after = len(node.eth.accounts)
    assert len(to_account) == 42
    assert after == (before + 1)
    yield node


@allure.title("Open a wallet")
@pytest.mark.P1
@pytest.mark.compatibility
def test_personal_openWallet(global_running_env):
    node = global_running_env.get_rand_node()
    listWallet = node.personal.listWallets
    if len(listWallet) > 0:
        assert None == node.personal.openWallet(listWallet[0]["url"], password)


@allure.title("Unlock account")
@pytest.mark.P1
@pytest.mark.compatibility
def test_personal_unlockAccount(test_personal_newAccount):
    listWallet = test_personal_newAccount.personal.listWallets
    if len(listWallet) > 0:
        addr1 = Web3.toChecksumAddress(listWallet[0]["accounts"][0]["address"])
        assert True == test_personal_newAccount.personal.unlockAccount(addr1, password)
        listWallet = test_personal_newAccount.personal.listWallets
        assert "Unlocked" == listWallet[0]["status"]


@allure.title("Lock account")
@pytest.mark.P1
@pytest.mark.compatibility
def test_personal_lockAccount(test_personal_newAccount):
    listWallet = test_personal_newAccount.personal.listWallets
    addr1 = Web3.toChecksumAddress(listWallet[0]["accounts"][0]["address"])
    assert True == test_personal_newAccount.personal.lockAccount(addr1)
    listWallet = test_personal_newAccount.personal.listWallets
    assert "Locked" == listWallet[0]["status"]


@allure.title("Import raw key")
@pytest.mark.P1
@pytest.mark.compatibility
@pytest.fixture()
def test_personal_importRawKey(global_running_env):
    node = global_running_env.get_rand_node()
    account = global_running_env.account.get_rand_account()
    prikey = account["prikey"]
    addr = node.personal.importRawKey(prikey, password)
    assert 42 == len(addr)
    yield node


@allure.title("data sign and ecRecover")
@pytest.mark.P1
@pytest.mark.compatibility
def test_personal_sign_ecRecover(global_running_env):
    node = global_running_env.get_rand_node()
    test_data = "0x11"
    signer = Web3.toChecksumAddress(node.eth.accounts[0])
    sign_data = node.personal.sign(test_data, signer, password)
    assert len(sign_data) == 132
    assert Web3.toChecksumAddress(node.personal.ecRecover(test_data, sign_data)) == signer


@allure.title("Sign transaction")
@pytest.mark.P1
@pytest.mark.compatibility
def test_personal_signTransaction(global_running_env):
    node = global_running_env.get_rand_node()
    account = global_running_env.account
    platon = Eth(node.web3)
    addr = account.account_with_money["address"]

    nonce = hex(platon.getTransactionCount(Web3.toChecksumAddress(addr)))
    transaction_dict = {
        "from": Web3.toChecksumAddress(addr),
        "to": Web3.toChecksumAddress(to_address),
        "value": "0x10000000000000",
        "data": "0x11",
        "gasPrice": "0x8250de00",
        "gas": "0x6fffffff",
        "nonce": nonce,
    }
    ret = node.personal.signTransaction(transaction_dict, password)
    assert ret is not None


def transaction_func(node, from_addr="", to_addr=to_address, value=1000, data='', gasPrice='100000000',
                     gas='21068', nonce=0, password=password):
    transaction_dict = {
        "from": Web3.toChecksumAddress(from_addr),
        "to": Web3.toChecksumAddress(to_addr),
        "value": value,
        "data": data,
        "gasPrice": gasPrice,
        "nonce": nonce,
    }
    if gas == '':
        gas = node.eth.estimateGas(transaction_dict)
        return gas
    transaction_dict["gas"] = gas
    global g_txHash
    g_txHash = node.personal.sendTransaction(transaction_dict, password)
    return g_txHash


# GetTransactionByHash
# GetRawTransactionByHash
# GetTransactionReceipt
@allure.title("Get raw transaction")
@pytest.mark.P1
def test_platon_getTransaction(global_running_env):
    node = global_running_env.get_rand_node()
    ret = node.eth.getTransaction("0x1111111111111111111111111111111111111111111111111111111111111111")
    assert ret is None
    print("check succeed: getTransaction by not exist hash!,ret:{}".format(ret))
    ret = node.eth.getRawTransaction(
        HexBytes("0x1111111111111111111111111111111111111111111111111111111111111111").hex())
    assert ret == "0x"


@allure.title("Send transaction based on gasprice's recommendations")
@pytest.mark.P1
def test_platon_gasPrice(global_running_env):
    node = global_running_env.get_rand_node()
    platon = Eth(node.web3)
    account = global_running_env.account
    from_address = account.account_with_money["address"]

    nCount = platon.getTransactionCount(Web3.toChecksumAddress(from_address))
    nonce = hex(nCount)
    gasprice = node.eth.gasPrice
    tx_hash = transaction_func(node=node, from_addr=from_address, to_addr=to_address, nonce=nonce, gasPrice=gasprice)
    assert len(tx_hash) == 32

    gasprice = node.eth.gasPrice * 2
    nCount = nCount + 1
    nonce = hex(nCount)
    tx_hash = transaction_func(node=node, from_addr=from_address, to_addr=to_address, nonce=nonce, gasPrice=gasprice)
    assert len(tx_hash) == 32

    gasprice = int(node.eth.gasPrice / 2)
    nCount = nCount + 1
    nonce = hex(nCount)

    status = 0
    try:
        transaction_func(node=node, from_addr=from_address, to_addr=to_address, nonce=nonce, gasPrice=gasprice)
        status = 1
    except Exception as e:
        print("\nUse less than the recommended gasprice:{}, nonce:{}, Send transaction failed,error message:{}".format(gasprice, nonce, e))
    assert status == 0


@allure.title("Get the block based on block number and block hash")
@pytest.mark.P1
def test_platon_GetBlock(global_running_env):
    node = global_running_env.get_rand_node()
    account = global_running_env.account
    platon = Eth(node.web3)
    address = account.account_with_money["address"]
    nCount = platon.getTransactionCount(Web3.toChecksumAddress(address))
    nonce = hex(nCount)

    gasprice = node.eth.gasPrice
    # send transaction
    if g_txHash is None:
        tx_hash = transaction_func(node=node, from_addr=address, to_addr=to_address, nonce=nonce, gasPrice=gasprice)
    else:
        tx_hash = g_txHash

    assert len(tx_hash) == 32
    tx_hash = HexBytes(tx_hash).hex()
    print("\ntransaction hash：{}".format(tx_hash))
    # Waiting for transaction on the chain
    result = node.eth.waitForTransactionReceipt(tx_hash)
    assert None != result
    # get block info by transaction receipt
    blockHash = result["blockHash"]
    blockNumber = result["blockNumber"]
    assert len(blockHash) == 32
    blockHash = HexBytes(blockHash).hex()
    # get block by blockHash
    # fullTx:Flase
    blockInfo = node.eth.getBlock(blockHash, False)
    blockHash = blockInfo['hash']
    blockNumber = blockInfo['number']
    assert len(blockHash) == 32
    assert blockNumber > 0
    fullTransaction = blockInfo["transactions"]
    assert len(fullTransaction) > 0
    # fullTx:True
    blockInfo = node.eth.getBlock(blockHash, True)
    blockHash = blockInfo['hash']
    assert len(blockHash) == 32
    fullTransaction = blockInfo["transactions"]
    assert len(fullTransaction) > 0
    # get block by blockNumber
    # fullTx:Flase
    blockInfo = node.eth.getBlock(blockNumber, False)
    blockHash = blockInfo['hash']
    assert len(blockHash) == 32
    fullTransaction = blockInfo["transactions"]
    assert len(fullTransaction) > 0
    # fullTx:True
    blockInfo = node.eth.getBlock(blockNumber, True)
    blockHash = blockInfo['hash']
    assert len(blockHash) == 32
    fullTransaction = blockInfo["transactions"]
    assert len(fullTransaction) > 0


@allure.title("Send the trade based on the gas estimate")
@pytest.mark.P1
def test_platon_estimateGas(global_running_env):
    node = global_running_env.get_rand_node()
    account = global_running_env.account
    platon = Eth(node.web3)
    address = account.account_with_money["address"]

    nCount = platon.getTransactionCount(Web3.toChecksumAddress(address))
    nonce = hex(nCount)
    estimateGas = transaction_func(node=node, from_addr=address, to_addr=to_address, nonce=nonce, gas='')
    gas = estimateGas

    tx_hash = transaction_func(node=node, from_addr=address, to_addr=to_address, nonce=nonce, gas=gas)
    assert len(tx_hash) == 32
    nCount = nCount + 1
    nonce = hex(nCount)

    gas = int(estimateGas * 2)
    tx_hash = transaction_func(node=node, from_addr=address, to_addr=to_address, nonce=nonce, gas=gas)
    assert len(tx_hash) == 32
    nCount = nCount + 1
    nonce = hex(nCount)

    gas = int(estimateGas / 2)
    status = 0
    try:
        transaction_func(node=node, from_addr=address, to_addr=to_address, nonce=nonce, gas=gas)
        status = 1
    except Exception as e:
        print("\nUse less gas than expected:【{}】,Send transaction failed,error message:{}".format(gas, e))
    assert status == 0


if __name__ == '__main__':
    pytest.main(['-v', 'test_rpc_personal.py'])
