# -*- coding: utf-8 -*-
'''
@Description: rpc用例
'''
import time

import allure
import pytest
from client_sdk_python import Web3
from hexbytes import HexBytes
password = "88888888"
address = "0x493301712671Ada506ba6Ca7891F436D29185821"
priKey = "a11859ce23effc663a9460e332ca09bd812acc390497f8dc7542b6938e13f8d7"
dataDir = "/home/platon/"
startApi = "eth,web3,net,txpool,platon,admin,personal"

w3 = None
id = ""
host = ""
http_url = ""
http_port = "6789"
p2p_prot = "16789"

test_data = "0x11"
sign_data = ""
signer = ""

addressCount = 1

@pytest.fixture()
def setNodeInfo(global_test_env):
    collusion_list = global_test_env.consensus_node_list
    if len(collusion_list) > 0:
        try:
            global w3
            global ws
            global id
            global host
            global http_url
            global http_port
            global p2p_prot
            global addressCount
            global priKey

            test_node = collusion_list[0]
            id = test_node.node_id
            host = test_node.host
            http_url = test_node.url
            http_port = test_node.rpc_port
            p2p_prot = test_node.p2p_port

            # rpc连接
            w3 = test_node.web3

            addressCount = len(w3.personal.listAccounts)
            addressCountTmp = len(w3.platon.accounts)
            assert addressCount == addressCountTmp
        except Exception as e:
            print("setNodeInfo error:{}>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>".format(e))
            w3 = None
            ws = None

@allure.title("列出所有的账户地址:personal.listAccounts")
@pytest.mark.P1
def test_personal_listAccounts(setNodeInfo):
    if w3 != None:
        listAccounts = w3.personal.listAccounts
        assert addressCount == len(listAccounts)
        print("\n所有账户地址:{}".format(listAccounts))

@allure.title("列出所有的钱包信息:personal.listWallets")
@pytest.mark.P1
def test_personal_listWallets(setNodeInfo):
    if w3 != None:
        listWallets = w3.personal.listWallets
        assert addressCount == len(listWallets)
        print("\n所有钱包信息:{}".format(listWallets))

@allure.title("创建一个新账户并产生一个新钱包:personal.newAccount")
@pytest.mark.P1
def test_personal_newAccount(setNodeInfo):
    if w3 != None:
        print("\n 本地钱包个数:{}".format(len(w3.platon.accounts)))
        new_account = w3.personal.newAccount(password)
        time.sleep(2)
        to_account = w3.toChecksumAddress(new_account)
        assert len(to_account) == 42
        print("\n创建账户成功:{}, 本地钱包个数:{}".format(to_account, len(w3.platon.accounts)))

@allure.title("打开一个钱包:personal.openWallet")
@pytest.mark.P1
def test_personal_openWallet(setNodeInfo):
    if w3 != None:
        listWallet = w3.personal.listWallets
        if len(listWallet) > 0 :
            assert None == w3.personal.openWallet(listWallet[0]["url"], password)
            print("\n打开钱包成功,钱包路径:{}".format(listWallet[0]["url"]))

@allure.title("解锁钱包:personal.unlockAccount")
@pytest.mark.P1
def test_personal_unlockAccount(setNodeInfo):
    if w3 != None:
        listWallet = w3.personal.listWallets
        if len(listWallet) > 0:
            addr1 = w3.toChecksumAddress(listWallet[0]["accounts"][0]["address"])
            try:
                assert True == w3.personal.unlockAccount(addr1, password)
                listWallet = w3.personal.listWallets
                assert "Unlocked" == listWallet[0]["status"]
                print("\n解锁钱包成功,钱包地址:{},状态:{}".format(addr1, listWallet[0]["status"]))
            except Exception as e:
                print("\n解锁失败, 钱包地址:{}, error message:{}".format(addr1, e))

@allure.title("上锁钱包:personal.lockAccount")
@pytest.mark.P1
def test_personal_lockAccount(setNodeInfo):
    if w3 != None:
        listWallet = w3.personal.listWallets
        if len(listWallet) > 0:
            addr1 = w3.toChecksumAddress(listWallet[0]["accounts"][0]["address"])
            assert True == w3.personal.lockAccount(addr1)
            listWallet = w3.personal.listWallets
            assert "Locked" == listWallet[0]["status"]
            print("\n锁钱包成功,钱包地址:{}，状态:{}".format(addr1, listWallet[0]["status"]))

@allure.title("导入钱包私钥:personal.importRawKey")
@pytest.mark.P1
def test_personal_importRawKey(setNodeInfo):
    if w3 != None:
        try:
            addr = w3.personal.importRawKey(priKey, password)
            assert 42 == len(addr)
            print("\n导入私钥成功,钱包地址:{}".format(addr))
        except Exception as e:
            print("\n导入私钥失败:{}".format(e))

@allure.title("签名数据和解签:personal.sign()/personal.ecRecover()")
@pytest.mark.P1
def test_personal_sign_ecRecover(setNodeInfo):
    if w3 != None:
        signer = Web3.toChecksumAddress(w3.platon.accounts[0])
        sign_data = w3.personal.sign(test_data, signer, password)
        assert len(sign_data) == 132
        print("\n签名数据成功,测试数据:{},签名数据sign_data:{}".format(test_data, sign_data))
        assert w3.toChecksumAddress(w3.personal.ecRecover(test_data, sign_data)) == signer
        print("\n解签成功,签名数据sign_data:{}, 解签数据:{}".format(sign_data, test_data))

# 签名交易
@allure.title("签名交易:personal.signTransaction")
@pytest.mark.P1
def test_personal_signTransaction(setNodeInfo):
    if w3 != None:
        nonce = hex(w3.platon.getTransactionCount(Web3.toChecksumAddress(address)))
        to = w3.platon.accounts[0]
        transaction_dict = {
            "from": Web3.toChecksumAddress(address),
            "to": Web3.toChecksumAddress(to),
            "value": "0x10000000000000",
            "data": "0x11",
            "gasPrice": "0x8250de00",
            "gas": "0x6fffffff",
            "nonce": nonce,
        }
        try:
            ret = w3.personal.signTransaction(transaction_dict, password)
            assert ret != None
            print("\n签名交易成功, 签名钱包地址:{}, 签名数据:{}".format(address, ret))
        except Exception as e:
            print("\n签名交易失败,error message:{}".format(e))


def transaction_func(from_addr=address, to_addr='', value='0x10000000000000', data='0x11', gasPrice='100000000',
                     gas='21068', nonce='0x11', password=password) :
    if w3 != None:
        transaction_dict = {
            "from": w3.toChecksumAddress(from_addr),
            "to": w3.toChecksumAddress(to_addr),
            "value": value,
            "data": data,
            "gasPrice": gasPrice,
            "nonce": nonce,
        }
        # 获取交易的预估gas
        if gas == '' :
            gas = w3.platon.estimateGas(transaction_dict)
            return gas
        transaction_dict["gas"] = gas
        tx_hash = w3.personal.sendTransaction(transaction_dict, password)
        return tx_hash

# 发送交易
@allure.title("签名交易:personal.sendTransaction")
@pytest.fixture()
def personal_sendTransaction(setNodeInfo):
    if w3 != None:
        nonce = hex(w3.platon.getTransactionCount(Web3.toChecksumAddress(address)))
        to = w3.platon.accounts[0]
        transaction_dict = {
            "from": Web3.toChecksumAddress(address),
            "to": Web3.toChecksumAddress(to),
            "value": "0x10000000000000",
            "data": "0x11",
            "gasPrice": "0x8250de00",
            "gas": "0x6fffffff",
            "nonce": nonce,
        }
        try:
            tx_hash = w3.personal.sendTransaction(transaction_dict, password)
            assert len(tx_hash) == 32
            print("\n发送交易成功,交易hash:{}".format(HexBytes(tx_hash).hex()))
            result = w3.platon.waitForTransactionReceipt(HexBytes(tx_hash).hex())
            assert None != result
            print("\n交易上链成功, 交易回执:{}".format(result))
        except Exception as e:
            print("\n发送交易失败，error message:{}".format(e))
            return ""
        return tx_hash

# GetTransactionByHash
# GetRawTransactionByHash
# GetTransactionReceipt
@allure.title("获取交易信息:personal.getTransaction/personal.getRawTransaction")
@pytest.mark.P1
def test_platon_getTransaction(personal_sendTransaction):
    if w3 != None:
        tx_hash = HexBytes(personal_sendTransaction).hex()
        # await
        try:
            ret = w3.platon.getTransaction(tx_hash)
            assert ret != None
            print("\ncheck succeed: getTransaction by exist hash!,ret:{}".format(ret))
            ret = w3.platon.getRawTransaction(tx_hash)
            assert ret != None
            print("check succeed: getRawTransaction by exist hash!,ret:{}".format(ret))
            ret = w3.platon.getTransaction("0x1111111111111111111111111111111111111111111111111111111111111111")
            assert ret == None
            print("check succeed: getTransaction by not exist hash!,ret:{}".format(ret))
            ret = w3.platon.getRawTransaction(HexBytes("0x1111111111111111111111111111111111111111111111111111111111111111").hex())
            assert ret == "0x"
            print("check succeed: getRawTransaction by not exist hash!,ret:{}".format(ret))
        except Exception as e:
            print("\nGetTransaction error:{}".format(e))

@allure.title("根据gasprice的建议值,发送交易")
@pytest.mark.P1
def test_platon_gasPrice(setNodeInfo):
    if w3 != None:
        nCount = w3.platon.getTransactionCount(Web3.toChecksumAddress(address))
        nonce = hex(nCount)
        to = w3.platon.accounts[0]
        gasprice = w3.platon.gasPrice
        try:
            tx_hash = transaction_func(to_addr=to, nonce=nonce, gasPrice=gasprice)
            assert len(tx_hash) == 32
            print("\n使用建议的交易gasprice:{},发送交易成功,交易hash：{}".format(gasprice, HexBytes(tx_hash).hex()))
        except Exception as e:
            print("\n使用建议的交易gasprice:{}, nonce:{}, 发送交易失败,error message:{}".format(gasprice, nonce, e))
        gasprice = w3.platon.gasPrice * 2
        nCount = nCount + 1
        nonce = hex(nCount)
        try:
            tx_hash = transaction_func(to_addr=to, nonce=nonce, gasPrice=gasprice)
            assert len(tx_hash) == 32
            print("\n使用大于建议的交易gasprice,发送交易成功,交易hash：{},  gasprice:{}".format(HexBytes(tx_hash).hex(), gasprice))
        except Exception as e:
            print("\n使用大于建议的交易gasprice:{}, nonce:{}, 发送交易失败,error message:{}".format(gasprice, nonce, e))
        gasprice = w3.platon.gasPrice - 100000
        nCount = nCount + 1
        nonce = hex(nCount)
        try:
            tx_hash = transaction_func(to_addr=to, nonce=nonce, gasPrice=gasprice)
            assert len(tx_hash) == 32
            print("\n使用小于建议的交易gasprice,发送交易成功,交易hash：{},  gasprice:{}".format(HexBytes(tx_hash).hex(), gasprice))
        except Exception as e:
            print("\n使用小于建议的交易gasprice;{}, nonce:{}, 发送交易失败,error message:{}".format(gasprice, nonce, e))

@allure.title("验证获取区块:根据区块高度和区块hash, fullTx:Flase/True")
@pytest.mark.P1
def test_platon_GetBlock(setNodeInfo):
    if w3 != None:
        nCount = w3.platon.getTransactionCount(Web3.toChecksumAddress(address))
        nonce = hex(nCount)
        to = w3.platon.accounts[0]
        gasprice = w3.platon.gasPrice
        # send transaction
        try:
            tx_hash = transaction_func(to_addr=to, nonce=nonce, gasPrice=gasprice)
            assert len(tx_hash) == 32
            tx_hash = HexBytes(tx_hash).hex()
            print("\n交易hash：{}".format(tx_hash))
        # Waiting for transaction on the chain
            result = w3.platon.waitForTransactionReceipt(tx_hash)
            assert None != result
        # get block info by transaction receipt
            blockHash = result["blockHash"]
            blockNumber = result["blockNumber"]
            assert len(blockHash) == 32
            blockHash = HexBytes(blockHash).hex()
            print("\n交易：【{}】已经上链, 所在区块高度:【{}]，区块hash:【{}】".format(tx_hash,blockNumber, blockHash))
        # get block by blockHash
        # fullTx:Flase
            blockInfo = w3.platon.getBlock(blockHash, False)
            blockHash = blockInfo['hash']
            blockNumber = blockInfo['number']
            assert len(blockHash) == 32
            assert blockNumber > 0
            fullTransaction = blockInfo["transactions"]
            assert len(fullTransaction) > 0
            print("\ngetBlock by blockHash fullTx为Flase,返回transactions信息为:{}".format(fullTransaction))
        # fullTx:True
            blockInfo = w3.platon.getBlock(blockHash, True)
            blockHash = blockInfo['hash']
            assert len(blockHash) == 32
            fullTransaction = blockInfo["transactions"]
            assert len(fullTransaction) > 0
            print("\ngetBlock by blockHash fullTx为True,返回transactions信息为:{}".format(fullTransaction))
            # get block by blockNumber
            # fullTx:Flase
            blockInfo = w3.platon.getBlock(blockNumber, False)
            blockHash = blockInfo['hash']
            assert len(blockHash) == 32
            fullTransaction = blockInfo["transactions"]
            assert len(fullTransaction) > 0
            print("\ngetBlock by blockNumber fullTx为Flase,返回transactions信息为:{}".format(fullTransaction))
            # fullTx:True
            blockInfo = w3.platon.getBlock(blockNumber, True)
            blockHash = blockInfo['hash']
            assert len(blockHash) == 32
            fullTransaction = blockInfo["transactions"]
            assert len(fullTransaction) > 0
            print("\ngetBlock by blockNumber fullTx为True,返回transactions信息为:{}".format(fullTransaction))
        except Exception as e:
            print("\n error message:{}".format(e))

@allure.title("根据交易的gas预估值,发送交易")
@pytest.mark.P1
def test_platon_estimateGas(setNodeInfo):
    if w3 != None:
        nCount = w3.platon.getTransactionCount(Web3.toChecksumAddress(address))
        nonce = hex(nCount)
        to = w3.platon.accounts[0]
        # 获取交易的预估值
        estimateGas = transaction_func(to_addr=to,nonce=nonce,gas='')
        gas = estimateGas
        print("\n交易预估gas:{}".format(gas))
        # 发送交易
        try:
            tx_hash = transaction_func(to_addr=to, nonce=nonce, gas=gas)
            assert len(tx_hash) == 32
            print("\n使用建议的交易gas:【{}】,发送交易成功,交易hash:【{}】".format(gas, HexBytes(tx_hash).hex()))
            nCount = nCount + 1
            nonce = hex(nCount)
        except Exception as e:
            print("\n使用建议的交易gas:【{}】,发送交易失败,error message:{}".format(gas, e))
        gas = estimateGas * 2
        try:
            tx_hash = transaction_func(to_addr=to, nonce=nonce, gas=gas)
            assert len(tx_hash) == 32
            print("\n使用大于建议的交易gas:【{}】,发送交易成功,交易hash:【{}】".format(gas, HexBytes(tx_hash).hex()))
            nCount = nCount + 1
            nonce = hex(nCount)
        except Exception as e:
            print("\n使用大于建议的交易gas:【{}】,发送交易失败,error message:{}".format(gas, e))
        gas = int(estimateGas/2)
        try:
            tx_hash = transaction_func(to_addr=to, nonce=nonce, gas=gas)
            assert len(tx_hash) == 32
            print("\n使用小于建议的交易gas:【{}】,发送交易成功,交易hash:【{}】".format(gas, HexBytes(tx_hash).hex()))
        except Exception as e:
            print("\n使用小于建议的交易gas:【{}】,发送交易失败,error message:{}".format(gas, e))

if __name__ == '__main__':
    pytest.main(['-v','test_rpc_personal.py'])