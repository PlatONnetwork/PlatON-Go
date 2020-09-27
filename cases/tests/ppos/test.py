import os
import time

import rlp
from alaya import HTTPProvider, Web3, WebsocketProvider, ppos
from alaya.eth import Eth
from alaya.middleware import geth_poa_middleware
from alaya.ppos import Ppos
from hexbytes import HexBytes
from ruamel import yaml

import conf
from common.load_file import LoadFile
# from conf.settings import TMP_ADDRES, ACCOUNT_FILE, BASE_DIR

accounts = {}


def connect_web3(url, chain_id=101):
    if "ws" in url:
        w3 = Web3(WebsocketProvider(url), chain_id=chain_id)
    else:
        w3 = Web3(HTTPProvider(url), chain_id=chain_id)
    w3.middleware_stack.inject(geth_poa_middleware, layer=0)
    return w3


def createRestrictingPlan(url, account, plan, pri_key):
    web3 = connect_web3(url)
    ppos = Ppos(web3)
    result = ppos.createRestrictingPlan(account, plan, pri_key)
    print(result)


def withdrewStaking(url, node_id, pri_key):
    web3 = connect_web3(url)
    ppos = Ppos(web3)
    result = ppos.withdrewStaking(node_id, pri_key)
    print(result)


def delegate(url, typ, node_id, amount, pri_key):
    web3 = connect_web3(url)
    ppos = Ppos(web3)
    result = ppos.delegate(typ, node_id, amount, pri_key)
    print(result)


def sendTransaction(url, from_address, prikey, to_address, value, chain_id):
    web3 = connect_web3(url)
    platon = Eth(web3)
    nonce = platon.getTransactionCount(from_address)
    print(nonce)
    gasPrice = platon.gasPrice
    transaction_dict = {
        "to": to_address,
        "gasPrice": gasPrice,
        "gas": 21000,
        "nonce": nonce,
        "data": '',
        "chainId": chain_id,
        "value": value,
    }

    signedTransactionDict = platon.account.signTransaction(
        transaction_dict, prikey
    )

    data = signedTransactionDict.rawTransaction
    result = HexBytes(platon.sendRawTransaction(data)).hex()
    # print(result)
    # log.debug("result:::::::{}".format(result))
    # res = platon.waitForTransactionReceipt(result)
    # print(res)


def get_candidatelist(url):
    web3 = connect_web3(url)
    ppos = Ppos(web3)
    result = ppos.getCandidateList()
    print(result)


def get_candinfo(url, node_id):
    web3 = connect_web3(url)
    ppos = Ppos(web3)
    result = ppos.getCandidateInfo(node_id)
    print(result)

def get_RestrictingPlan(url, address):
    web3 = connect_web3(url)
    ppos = Ppos(web3)
    result = ppos.getRestrictingInfo(address)
    print(result)
#
# def create_address(url):
#     """
#     创建新钱包地址
#     """
#     web3 = connect_web3(url)
#     platon = Eth(web3)
#     account = platon.account.create(net_type=web3.net_type)
#     address = account.address
#     prikey = account.privateKey.hex()[2:]
#     account = {
#         "address": address,
#         "nonce": 0,
#         "balance": 0,
#         "prikey": prikey,
#     }
#     accounts = {}
#     raw_accounts = LoadFile(ACCOUNT_FILE).get_data()
#     print(raw_accounts)
#     for account1 in raw_accounts:
#         accounts[account1['address']] = account1
#     print(accounts)
#     accounts[address] = account
#     # todo delete debug
#     accounts = list(accounts.values())
#     with open(os.path.join(BASE_DIR, "deploy/tmp/accounts.yml"), mode="w", encoding="UTF-8") as f:
#         yaml.dump(accounts, f, Dumper=yaml.RoundTripDumper)
#
#
# def cycle_sendTransaction(url):
#     """
#
#     """
#
    # with open(TMP_ADDRES, 'a', encoding='utf-8') as f:
    #     f.write("2")


if __name__ == '__main__':
    # url = 'http://192.168.9.222:6789'
    url = 'http://10.1.1.58:6789'
    # url = 'http://10.0.0.44:6789'
    # url = 'http://149.129.180.78:6789'
    # epoch = 100
    amount = Web3.toWei(1, 'ether')
    # plan = [{'Epoch': epoch, 'Amount': amount}]
    # createRestrictingPlan(url, account, plan, pri_key)
    # delegate(url, 0, nodeid, amount, pri_key)
    to_address = 'lat13l39glde394a6kkrm5aenj4ty7m7565x8sgtrf'
    pri_key = '91751513fa39f02ada9a7110bef0a20e03375e9b05d78036e84e91366276e5d8'
    from_address = 'lax184zj2xdms82dvg5ypacsk48qw3ch0q9rtfrmp3'
    # node_id = '71bc24068d1f1f65331ad7573806bf58186375ef993dddf3ea51c8d0da162c801689aed5aa9e809396cd60273af1d2826d918e36ce4d003c578371a7b3b8b429'
    # pri_key1 = 'd357920de1df4ecb00cbce60ded2d73f3f51fd1e9fb79b08f366e301e849bd9d'
    while 1:
        sendTransaction(url, from_address, pri_key, to_address, amount, 298)
    # withdrewStaking(url, node_id, pri_key1)
    # node_id = 'e2181d8dc731b14117ba6d982ce163fc7b9b14bbbaf9cb3c343ef72c24cf3ed568cac6ecbc30fddf9012320fab99f6be6ab37132d083cb514100bdb4b90fff5e'
    # get_candinfo(url, node_id)
    # get_candidatelist(url)
    # addresss = 'lat13l39glde394a6kkrm5aenj4ty7m7565x8sgtrf'
    # get_RestrictingPlan(url, 'lax17ax3lr6qy405sf03ncema3nmusdqmpzq7ujvpz')