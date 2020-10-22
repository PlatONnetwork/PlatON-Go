import os
import time

import rlp
from alaya import HTTPProvider, Web3, WebsocketProvider, ppos
from alaya.eth import Eth
from alaya.middleware import geth_poa_middleware
from alaya.ppos import Ppos
from hexbytes import HexBytes
from ruamel import yaml
from alaya.contract import ContractConstructor
import conf
from common.load_file import LoadFile

# from conf.settings import TMP_ADDRES, ACCOUNT_FILE, BASE_DIR
from tests.lib import Pip, json

accounts = {}


def connect_web3(url, chain_id=201030):
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


def increase_staking(url, type, node_id, amount, pri_key):
    web3 = connect_web3(url)
    ppos = Ppos(web3)
    result = ppos.increaseStaking(type, node_id, amount, pri_key)
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
    res = platon.waitForTransactionReceipt(result)
    print(res)


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


def get_listGovernParam(url, module=None, from_address=None):
    web3 = connect_web3(url)
    ppos = Ppos(web3)
    # ppos.web3.platon
    if module is None:
        module = ""
    data = rlp.encode([rlp.encode(int(2106)), rlp.encode(module)])
    raw_data = ppos.call_obj(ppos, from_address, Web3.pipAddress, data)
    data = str(raw_data, encoding="utf-8")
    if data == "":
        return ""
    print(json.loads(data))

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


def fff(url, from_address):
    web3 = connect_web3(url)
    platon = Eth(web3)
    # result = platon.getTransactionCount(from_address)
    result = platon.getBalance(from_address)
    print(result)


if __name__ == '__main__':
    # url = 'http://192.168.10.221:6789'
    # url = 'http://10.1.1.51:6789'
    # url = 'http://10.0.0.44:6789'
    url = 'http:// 47.241.4.217:6789'
    account = 'atx1v2sjx7whxggphhlyxumvxpxyw4mk05st5k92md'
    pri_key = 'f51ca759562e1daf9e5302d121f933a8152915d34fcbc27e542baf256b5e4b74'
    pri_key1 = '6558c64ea9b14a069b42148e373616c834ed75828b456ea001390996d7e206a2'
    from_address = 'atx1zkrxx6rf358jcvr7nruhyvr9hxpwv9unj58er9'
    epoch1 = 10
    epoch2 = 20
    amount = Web3.toWei(1, 'ether')
    list = ['atx1r8pvmt7hk6lk8uk7dtnfyrpcy9l8rfjry34uq9',
            'atx1nccsq48wery09qlma3rapree588cafwlpll8cr',
            'atx14w3m34dmx5xjr7wc5yhg6xtp59j3qy4m77t4vl',
            'atx10u83ynv0sdjrjg4na66pg3fqyx8uc89g7zvnyu',
            'atx13e2hlak7jytcdlv47vmqa2hkaeyv62ggv3hymk',
            'atx1z4e77tam2lv8vsrzkazz659avlvx3v8jyh0ywt',
            'atx1pvd37yxnmdtdujaf3dwc4xpcukdzsf5q2hyc9a',
            'atx17rf2ylcauvse5sxsvnn05c30zscey0zwzhs5wp',
            'atx1cux96yve5d335hqt3a7pdjygpj77pny3xqq7qu',
            'atx1s9qht6yzg6f53u26fe69vl4qrt6fs0sphth3ez',
            'atx1ckxg24sa4clv239y93talm79h7ac8r20t4dl8e',
            'atx19qtc92y2s9a6dyvuqxrqwpsaztz95mel4xuhkv']
    ac = 'atx1r8pvmt7hk6lk8uk7dtnfyrpcy9l8rfjry34uq9'
    # plan = [{'Epoch': epoch, 'Amount': amount}]
    # createRestrictingPlan(url, account, plan, pri_key)
    # delegate(url, 0, nodeid, amount, pri_key)
    plan = [{'Epoch': epoch1, 'Amount': amount}, {'Epoch': epoch2, 'Amount': amount}]
    address = 'atx1r8pvmt7hk6lk8uk7dtnfyrpcy9l8rfjry34uq9'
    node_id = '6cda52721a11a5034ae0dfc03ebe0a60a797e0240f9bba427957abeeb2e367c09ed099c6871bf17158c5d694c4d5ccad363b38055e345898ff02a88e17d66149'

    # node_id = '71bc24068d1f1f65331ad7573806bf58186375ef993dddf3ea51c8d0da162c801689aed5aa9e809396cd60273af1d2826d918e36ce4d003c578371a7b3b8b429'
    # pri_key1 = 'd357920de1df4ecb00cbce60ded2d73f3f51fd1e9fb79b08f366e301e849bd9d'
    # for i in list:
    #     print(i)
    # increase_staking(url, 0, node_id, amount, pri_key1)
    # createRestrictingPlan(url, ac, plan, pri_key)
    # get_RestrictingPlan(url, ac)
    # fff(url, i)
    # sendTransaction(url, from_address, pri_key, i, amount, 201030)

    # withdrewStaking(url, node_id, pri_key1)
    node_id = 'b53a0dde131a938ee6f98ee42f88b294897fbce32a028292feadda042ecf56d4b9f54f82407a45d64dbfa00ea450faad383dd7fd117c645c7dea0b16cf35c020'
    get_candinfo(url, node_id)
    # get_candidatelist(url)
    # addresss = 'lat13l39glde394a6kkrm5aenj4ty7m7565x8sgtrf'
    # get_RestrictingPlan(url, account)
    # fff()
    # get_listGovernParam(url)
