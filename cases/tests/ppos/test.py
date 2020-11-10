from alaya import HTTPProvider, Web3, WebsocketProvider
from alaya.eth import Eth
from alaya.middleware import geth_poa_middleware
from alaya.ppos import Ppos
from hexbytes import HexBytes

# from conf.settings import TMP_ADDRES, ACCOUNT_FILE, BASE_DIR

accounts = {}


def connect_web3(url, chain_id=201018):
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


# def createstaking(url, typ, benifit_address, node_id, pri_key, amount):
#     web3 = connect_web3(url)
#     ppos = Ppos(web3)
#     external_id = None
#     node_name = None
#     website = None
#     details = None
#     result = ppos.createStaking(typ, benifit_address, node_id, external_id, node_name, website, details, amount,
#                                  program_version, program_version_sign, bls_pubkey, bls_proof, pri_key, reward_per, transaction_cfg=None)
#     print(result)


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


def get_VerifierList(url):
    web3 = connect_web3(url)
    ppos = Ppos(web3)
    result = ppos.getVerifierList()
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


# def get_listGovernParam(url, module=None, from_address=None):
#     web3 = connect_web3(url)
#     ppos = Ppos(web3)
#     # ppos.web3.platon
#     if module is None:
#         module = ""
#     data = rlp.encode([rlp.encode(int(2106)), rlp.encode(module)])
#     raw_data = ppos.call_obj(ppos, from_address, Web3.pipAddress, data)
#     data = str(raw_data, encoding="utf-8")
#     if data == "":
#         return ""
#     print(json.loads(data))

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


def fff(url, from_address=None):
    web3 = connect_web3(url)
    platon = Eth(web3)
    # print(platon.g)
    # result = platon.getTransactionCount(from_address)
    # result = platon.getBalance(from_address)
    # print(platon.blockNumber)
    # print(result)


if __name__ == '__main__':
    # url = 'http://192.168.10.221:6789'
    # url = 'http://10.1.1.51:6789'
    # url = 'http://10.0.0.44:6789'
    # url = 'http:// 47.241.4.217:6789'
    url = 'http://154.85.35.163:80'
    account = 'atx1nqee6ze7flv6hunuajqt7el70u7dp49758fxch'
    pri_key = 'f51ca759562e1daf9e5302d121f933a8152915d34fcbc27e542baf256b5e4b74'
    pri_key1 = '6558c64ea9b14a069b42148e373616c834ed75828b456ea001390996d7e206a2'
    from_address = 'atx1zkrxx6rf358jcvr7nruhyvr9hxpwv9unj58er9'
    epoch1 = 10
    epoch2 = 20
    # amount = Web3.toWei(1, 'ether')
    amount1 = Web3.toWei(833, 'ether')
    amount2 = Web3.toWei(837, 'ether')
    # list = ['atx1r8pvmt7hk6lk8uk7dtnfyrpcy9l8rfjry34uq9',
    #         'atx1nccsq48wery09qlma3rapree588cafwlpll8cr',
    #         'atx14w3m34dmx5xjr7wc5yhg6xtp59j3qy4m77t4vl',
    #         'atx10u83ynv0sdjrjg4na66pg3fqyx8uc89g7zvnyu',
    #         'atx13e2hlak7jytcdlv47vmqa2hkaeyv62ggv3hymk',
    #         'atx1z4e77tam2lv8vsrzkazz659avlvx3v8jyh0ywt',
    #         'atx1pvd37yxnmdtdujaf3dwc4xpcukdzsf5q2hyc9a',
    #         'atx17rf2ylcauvse5sxsvnn05c30zscey0zwzhs5wp',
    #         'atx1cux96yve5d335hqt3a7pdjygpj77pny3xqq7qu',
    #         'atx1s9qht6yzg6f53u26fe69vl4qrt6fs0sphth3ez',
    #         'atx1ckxg24sa4clv239y93talm79h7ac8r20t4dl8e',
    #         'atx19qtc92y2s9a6dyvuqxrqwpsaztz95mel4xuhkv']
    # ac = 'atx1r8pvmt7hk6lk8uk7dtnfyrpcy9l8rfjry34uq9'
    list = ['atx1y2pejagyaj2s8jzrcrtrc2xl3x5l5r3hxfu8z3',
            'atx1n39fsmhhleyrru6smcy9m9q4w6t4mxz5439tc8',
            'atx124u7787cqrsclat26jv2vtph29slanruhsz95j',
            'atx1jnh4qfvuvam86hfhy08daw5ahsnmryjcuyzlx8',
            'atx1quxk2dju0774w65lu4dzuhac28q8c5ak8lk46e',
            'atx1j89t4zw9znq7582yf8ut99a8ynz8tp0m75ksj6',
            'atx1lr0tn4c7q37y2yjfe5a6vks7trjgznr3t7yy3h',
            'atx1we7rgld9n8att6l9tw5m5lpkqfa5y3dedk7yuk',
            'atx12mevpl56jz07wu0lk4hmlmqwmcfwhukxqxejmy',
            'atx1wjk3qwhmvg3e5chnk5k8kxp23cyjh3e4np4xwx',
            'atx1740vqwuaxetpeym6j7kvp6w0c28r4vyz4m5znr',
            'atx1wuth2jvd98dfaw40dg294xzj8hlhq3s3r69v39',
            'atx1sg03frnmah7rxahaxphklfxl5aug37ytgayfu7',
            'atx1nyzkjzy8lgmmswpyvpw6dfe9hw9krkhm8gzdue',
            'atx10sccfk6v8rmv35xh9fhaaxchw48u2xggt3vamv']
    # plan = [{'Epoch': epoch, 'Amount': amount}]
    # createRestrictingPlan(url, account, plan, pri_key)
    # delegate(url, 0, nodeid, amount, pri_key)
    plan = [{'Epoch': 10, 'Amount': amount1},
            {'Epoch': 20, 'Amount': amount1},
            {'Epoch': 30, 'Amount': amount1},
            {'Epoch': 40, 'Amount': amount1},
            {'Epoch': 50, 'Amount': amount1},
            {'Epoch': 60, 'Amount': amount1},
            {'Epoch': 70, 'Amount': amount1},
            {'Epoch': 80, 'Amount': amount1},
            {'Epoch': 90, 'Amount': amount1},
            {'Epoch': 100, 'Amount': amount1},
            {'Epoch': 110, 'Amount': amount1},
            {'Epoch': 120, 'Amount': amount2}]
    # address = 'atx1r8pvmt7hk6lk8uk7dtnfyrpcy9l8rfjry34uq9'
    node_id = 'eac3b0c6569786ca0b6acaa5c80ee71854a00b6e3a53852dfdefe2ea1c6f1ada7eca7208ce942bbb5357adb16b1efe499b91e90655b3a05cc23603c6c421d4e9'

    # node_id = '71bc24068d1f1f65331ad7573806bf58186375ef993dddf3ea51c8d0da162c801689aed5aa9e809396cd60273af1d2826d918e36ce4d003c578371a7b3b8b429'
    # pri_key1 = 'd357920de1df4ecb00cbce60ded2d73f3f51fd1e9fb79b08f366e301e849bd9d'
    # for i in list:
    #     print(i)
    # # increase_staking(url, 0, node_id, amount, pri_key1)
    #     createRestrictingPlan(url, i, plan, pri_key)
    #     get_RestrictingPlan(url, account)
    # fff(url)
    #     sendTransaction(url, from_address, pri_key, i, amount, 201030)

    # withdrewStaking(url, node_id, pri_key1)
    # node_id = 'b53a0dde131a938ee6f98ee42f88b294897fbce32a028292feadda042ecf56d4b9f54f82407a45d64dbfa00ea450faad383dd7fd117c645c7dea0b16cf35c020'
    # get_VerifierList(url)
    get_candidatelist(url)
    # get_candinfo(url, node_id)
    # get_candidatelist(url)
    # addresss = 'lat13l39glde394a6kkrm5aenj4ty7m7565x8sgtrf'
    # get_RestrictingPlan(url, account)
    # fff()
    # get_listGovernParam(url)
