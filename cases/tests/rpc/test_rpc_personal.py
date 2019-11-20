# -*- coding: utf-8 -*-
'''
@Description: rpc用例
'''
import time

import allure
import pytest
from client_sdk_python import Web3
from client_sdk_python.eth import Eth
from hexbytes import HexBytes
password = "88888888"
to_address = "0xdfdbb962a03bd270e1e8235a3d11b5775334c7d7"
g_txHash = None


@allure.title("列出所有的账户地址:personal.listAccounts")
@pytest.mark.P1
@pytest.mark.compatibility
def test_personal_listAccounts(global_running_env):
    node = global_running_env.get_rand_node()
    assert len(node.personal.listAccounts) >= 0


@allure.title("列出所有的钱包信息:personal.listWallets")
@pytest.mark.P1
@pytest.mark.compatibility
def test_personal_listWallets(global_running_env):
    node = global_running_env.get_rand_node()
    listWallets = node.personal.listWallets
    assert len(listWallets) >= 0
    print("\n所有钱包信息:{}".format(listWallets))


@allure.title("创建一个新账户并产生一个新钱包:personal.newAccount")
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
#    print("\n创建账户成功:{}, 本地钱包个数:{}".format(to_account, after))
    yield node


@allure.title("打开一个钱包:personal.openWallet")
@pytest.mark.P1
@pytest.mark.compatibility
def test_personal_openWallet(global_running_env):
    node = global_running_env.get_rand_node()
    listWallet = node.personal.listWallets
    if len(listWallet) > 0:
        assert None == node.personal.openWallet(listWallet[0]["url"], password)
        print("\n打开钱包成功,钱包路径:{}".format(listWallet[0]["url"]))


@allure.title("解锁钱包:personal.unlockAccount")
@pytest.mark.P1
@pytest.mark.compatibility
def test_personal_unlockAccount(test_personal_newAccount):
    listWallet = test_personal_newAccount.personal.listWallets
    if len(listWallet) > 0:
        addr1 = Web3.toChecksumAddress(listWallet[0]["accounts"][0]["address"])
        assert True == test_personal_newAccount.personal.unlockAccount(addr1, password)
        listWallet = test_personal_newAccount.personal.listWallets
        assert "Unlocked" == listWallet[0]["status"]
        print("\n解锁钱包成功,钱包地址:{},状态:{}".format(addr1, listWallet[0]["status"]))


@allure.title("上锁钱包:personal.lockAccount")
@pytest.mark.P1
@pytest.mark.compatibility
def test_personal_lockAccount(test_personal_newAccount):
    listWallet = test_personal_newAccount.personal.listWallets
    addr1 = Web3.toChecksumAddress(listWallet[0]["accounts"][0]["address"])
    assert True == test_personal_newAccount.personal.lockAccount(addr1)
    listWallet = test_personal_newAccount.personal.listWallets
    assert "Locked" == listWallet[0]["status"]
    print("\n锁钱包成功,钱包地址:{}，状态:{}".format(addr1, listWallet[0]["status"]))


@allure.title("导入钱包私钥:personal.importRawKey")
@pytest.mark.P1
@pytest.mark.compatibility
@pytest.fixture()
def test_personal_importRawKey(global_running_env):
    node = global_running_env.get_rand_node()
    account = global_running_env.account.get_rand_account()
    prikey = account["prikey"]
    addr = node.personal.importRawKey(prikey, password)
    assert 42 == len(addr)
    print("\n导入私钥成功,钱包地址:{}".format(addr))
    yield node


@allure.title("签名数据和解签:personal.sign()/personal.ecRecover()")
@pytest.mark.P1
@pytest.mark.compatibility
def test_personal_sign_ecRecover(global_running_env):
    node = global_running_env.get_rand_node()
    test_data = "0x11"
    signer = Web3.toChecksumAddress(node.eth.accounts[0])
    sign_data = node.personal.sign(test_data, signer, password)
    assert len(sign_data) == 132
    print("\n签名数据成功,测试数据:{},签名数据sign_data:{}".format(test_data, sign_data))
    assert Web3.toChecksumAddress(node.personal.ecRecover(test_data, sign_data)) == signer
    print("\n解签成功,签名数据sign_data:{}, 解签数据:{}".format(sign_data, test_data))

# 签名交易
@allure.title("签名交易:personal.signTransaction")
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
    print("\n签名交易成功, 签名钱包地址:{}, 签名数据:{}".format(addr, ret))


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
    # 获取交易的预估gas
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
@allure.title("获取交易信息:personal.getTransaction/personal.getRawTransaction")
@pytest.mark.P1
def test_platon_getTransaction(global_running_env):
    node = global_running_env.get_rand_node()
    ret = node.eth.getTransaction("0x1111111111111111111111111111111111111111111111111111111111111111")
    assert ret is None
    print("check succeed: getTransaction by not exist hash!,ret:{}".format(ret))
    ret = node.eth.getRawTransaction(
        HexBytes("0x1111111111111111111111111111111111111111111111111111111111111111").hex())
    assert ret == "0x"
    print("check succeed: getRawTransaction by not exist hash!,ret:{}".format(ret))


@allure.title("根据gasprice的建议值,发送交易")
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
    print("\n使用建议的交易gasprice:{},发送交易成功,交易hash：{}".format(gasprice, HexBytes(tx_hash).hex()))

    gasprice = node.eth.gasPrice * 2
    nCount = nCount + 1
    nonce = hex(nCount)
    tx_hash = transaction_func(node=node, from_addr=from_address, to_addr=to_address, nonce=nonce, gasPrice=gasprice)
    assert len(tx_hash) == 32
    print("\n使用大于建议的交易gasprice,发送交易成功,交易hash：{},  gasprice:{}".format(HexBytes(tx_hash).hex(), gasprice))

    gasprice = int(node.eth.gasPrice / 2)
    nCount = nCount + 1
    nonce = hex(nCount)

    # 异常测试场景
    status = 0
    try:
        transaction_func(node=node, from_addr=from_address, to_addr=to_address, nonce=nonce, gasPrice=gasprice)
        status = 1
    except Exception as e:
        print("\n使用小于建议的交易gasprice:{}, nonce:{}, 发送交易失败,error message:{}".format(gasprice, nonce, e))
    assert status == 0


@allure.title("验证获取区块:根据区块高度和区块hash, fullTx:Flase/True")
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
    print("\n交易hash：{}".format(tx_hash))
    # Waiting for transaction on the chain
    result = node.eth.waitForTransactionReceipt(tx_hash)
    assert None != result
    # get block info by transaction receipt
    blockHash = result["blockHash"]
    blockNumber = result["blockNumber"]
    assert len(blockHash) == 32
    blockHash = HexBytes(blockHash).hex()
    print("\n交易：【{}】已经上链, 所在区块高度:【{}]，区块hash:【{}】".format(tx_hash, blockNumber, blockHash))
    # get block by blockHash
    # fullTx:Flase
    blockInfo = node.eth.getBlock(blockHash, False)
    blockHash = blockInfo['hash']
    blockNumber = blockInfo['number']
    assert len(blockHash) == 32
    assert blockNumber > 0
    fullTransaction = blockInfo["transactions"]
    assert len(fullTransaction) > 0
    print("\ngetBlock by blockHash fullTx为Flase,返回transactions信息为:{}".format(fullTransaction))
    # fullTx:True
    blockInfo = node.eth.getBlock(blockHash, True)
    blockHash = blockInfo['hash']
    assert len(blockHash) == 32
    fullTransaction = blockInfo["transactions"]
    assert len(fullTransaction) > 0
    print("\ngetBlock by blockHash fullTx为True,返回transactions信息为:{}".format(fullTransaction))
    # get block by blockNumber
    # fullTx:Flase
    blockInfo = node.eth.getBlock(blockNumber, False)
    blockHash = blockInfo['hash']
    assert len(blockHash) == 32
    fullTransaction = blockInfo["transactions"]
    assert len(fullTransaction) > 0
    print("\ngetBlock by blockNumber fullTx为Flase,返回transactions信息为:{}".format(fullTransaction))
    # fullTx:True
    blockInfo = node.eth.getBlock(blockNumber, True)
    blockHash = blockInfo['hash']
    assert len(blockHash) == 32
    fullTransaction = blockInfo["transactions"]
    assert len(fullTransaction) > 0
    print("\ngetBlock by blockNumber fullTx为True,返回transactions信息为:{}".format(fullTransaction))


@allure.title("根据交易的gas预估值,发送交易")
@pytest.mark.P1
def test_platon_estimateGas(global_running_env):
    node = global_running_env.get_rand_node()
    account = global_running_env.account
    platon = Eth(node.web3)
    address = account.account_with_money["address"]

    nCount = platon.getTransactionCount(Web3.toChecksumAddress(address))
    nonce = hex(nCount)
    # 获取交易的预估值
    estimateGas = transaction_func(node=node, from_addr=address, to_addr=to_address, nonce=nonce, gas='')
    gas = estimateGas
    print("\n交易预估gas:{}".format(gas))

    # 发送交易
    tx_hash = transaction_func(node=node, from_addr=address, to_addr=to_address, nonce=nonce, gas=gas)
    assert len(tx_hash) == 32
    print("\n使用预估的交易gas:【{}】,发送交易成功,交易hash:【{}】".format(gas, HexBytes(tx_hash).hex()))
    nCount = nCount + 1
    nonce = hex(nCount)

    gas = int(estimateGas * 2)
    tx_hash = transaction_func(node=node, from_addr=address, to_addr=to_address, nonce=nonce, gas=gas)
    assert len(tx_hash) == 32
    print("\n使用大于预估的交易gas:【{}】,发送交易成功,交易hash:【{}】".format(gas, HexBytes(tx_hash).hex()))
    nCount = nCount + 1
    nonce = hex(nCount)

    gas = int(estimateGas / 2)
    # 异常测试场景
    status = 0
    try:
        transaction_func(node=node, from_addr=address, to_addr=to_address, nonce=nonce, gas=gas)
        status = 1
    except Exception as e:
        print("\n使用小于预估的交易gas:【{}】,发送交易失败,error message:{}".format(gas, e))
    assert status == 0


if __name__ == '__main__':
    pytest.main(['-v', 'test_rpc_personal.py'])
