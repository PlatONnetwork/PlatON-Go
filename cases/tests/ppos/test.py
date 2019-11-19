import rlp
import json
from client_sdk_python import HTTPProvider, Web3, WebsocketProvider
from client_sdk_python.eth import Eth
from client_sdk_python.middleware import geth_poa_middleware
from client_sdk_python.utils.transactions import send_obj_transaction
from hexbytes import HexBytes


def connect_web3(url):
    '''
    连接web3服务,增加区块查询中间件,用于实现eth_getBlockByHash,eth_getBlockByNumber等方法
    '''
    if "ws" in url:
        w3 = Web3(WebsocketProvider(url))
    else:
        w3 = Web3(HTTPProvider(url))
    w3.middleware_stack.inject(geth_poa_middleware, layer=0)

    return w3


class Ppos:
    def __init__(self, url, address, privatekey, chainid=120):
        self.web3 = connect_web3(url)
        if not self.web3.isConnected():
            raise Exception("node connection failed")
        self.eth = Eth(self.web3)
        self.address = Web3.toChecksumAddress(address)
        self.privatekey = privatekey
        self.gasPrice = 1000000000
        self.gas = "0x6fffff"
        self.chainid = chainid

    def get_result(self, tx_hash):
        result = self.eth.waitForTransactionReceipt(tx_hash)
        """查看eventData"""
        data = result['logs'][0]['data']
        if data[:2] == '0x':
            data = data[2:]
        data_bytes = rlp.decode(bytes.fromhex(data))[0]
        event_data = bytes.decode(data_bytes)
        event_data = json.loads(event_data)
        print(event_data)
        return event_data

    def send_raw_transaction(self, data, from_address, to_address, gasPrice, gas, value, privatekey=None):
        nonce = self.eth.getTransactionCount(from_address)
        if not privatekey:
            privatekey = self.privatekey
        if value > 0:
            transaction_dict = {
                "to": to_address,
                "gasPrice": gasPrice,
                "gas": gas,
                "nonce": nonce,
                "data": data,
                "chainId": self.chainid,
                "value": self.web3.toWei(value, "ether")
            }
        else:
            transaction_dict = {
                "to": to_address,
                "gasPrice": gasPrice,
                "gas": gas,
                "nonce": nonce,
                "data": data,
                "chainId": self.chainid
            }
        signedTransactionDict = self.eth.account.signTransaction(
            transaction_dict, privatekey
        )
        data = signedTransactionDict.rawTransaction
        result = HexBytes(self.eth.sendRawTransaction(data)).hex()
        return result

    def createStaking(self, typ, benifitAddress, nodeId, externalId, nodeName, website, details, amount, programVersion,
                      programVersionSign, blsPubKey, blsProof, privatekey=None, from_address=None, gasPrice=None,
                      gas=None):
        to_address = "0x1000000000000000000000000000000000000002"
        if benifitAddress[:2] == '0x':
            benifitAddress = benifitAddress[2:]
        if programVersionSign[:2] == '0x':
            programVersionSign = programVersionSign[2:]
            data = HexBytes(rlp.encode([rlp.encode(int(1000)),
                                        rlp.encode(int(typ)),
                                        rlp.encode(bytes.fromhex(benifitAddress)),
                                        rlp.encode(bytes.fromhex(nodeId)),
                                        rlp.encode(externalId), rlp.encode(nodeName), rlp.encode(website),
                                        rlp.encode(details), rlp.encode(self.web3.toWei(amount, 'ether')),
                                        rlp.encode(programVersion),
                                        rlp.encode(bytes.fromhex(programVersionSign)),
                                        rlp.encode(bytes.fromhex(blsPubKey)),
                                        rlp.encode(bytes.fromhex(blsProof))])).hex()
            if not privatekey:
                privatekey = self.privatekey
            if not from_address:
                from_address = self.address
            if not gasPrice:
                gasPrice = self.gasPrice
            if not gas:
                # transactiondict = {"to": to_address, "data": data}
                gas = self.gas
            result = self.send_raw_transaction(data, from_address, to_address, gasPrice, gas, 0, privatekey)
            return self.get_result(result)

    def getCandidateInfo(self, nodeId):
        data = rlp.encode([rlp.encode(int(1105)), rlp.encode(bytes.fromhex(nodeId))])
        to_address = "0x1000000000000000000000000000000000000002"
        recive = self.eth.call({
            "from": self.address,
            "to": to_address,
            "data": data
        })
        recive = str(recive, encoding="utf8")
        recive = recive.replace('\\', '').replace('"{', '{').replace('}"', '}')
        recive = json.loads(recive)
        # print(recive)
        if recive["Data"] != "":
            recive["Data"]["Shares"] = int(recive["Data"]["Shares"], 16)
            recive["Data"]["Released"] = int(recive["Data"]["Released"], 16)
            recive["Data"]["ReleasedHes"] = int(recive["Data"]["ReleasedHes"], 16)
            recive["Data"]["RestrictingPlan"] = int(recive["Data"]["RestrictingPlan"], 16)
            recive["Data"]["RestrictingPlanHes"] = int(recive["Data"]["RestrictingPlanHes"], 16)
        return recive

    def createRestrictingPlan(self, account, plan, pri_key, transaction_cfg=None):
        """
        Create a lockout plan
        :param account: Locked account release account
        :param plan:
        """
        if account[:2] == '0x':
            account = account[2:]
        plan_list = []
        for dict_ in plan:
            v = [dict_[k] for k in dict_]
            plan_list.append(v)
        rlp_list = rlp.encode(plan_list)
        data = rlp.encode([rlp.encode(int(4000)), rlp.encode(bytes.fromhex(account)), rlp_list])
        return send_obj_transaction(self, data, self.web3.restrictingAddress, pri_key, transaction_cfg)


def get_staking_program(w3):
    """
    根据rpc获取发质押的4个参数
    """
    msg = w3.admin.getProgramVersion()
    proof = w3.admin.getSchnorrNIZKProve()
    ProgramVersionSign = msg["Sign"]
    ProgramVersion = msg["Version"]
    return ProgramVersion, ProgramVersionSign, proof


def create():
    url = "http://192.168.16.189:6789"
    stak_address = "0x48c867ddBF22062704D6c81d3FA256bc6fc8b6bC"
    stak_pri = "96f1f76c45bc2dd9c0f84a11da4ec104ae95661871284a29413a454a70b15307"
    ppos = Ppos(url, stak_address, stak_pri)
    benifitAddress = "0xCeCA295e1471B3008D20b017c7Df7d4F338A7FbA"
    nodeId = "2d25f7686573602334589ac2e606a3743d34fcae0c7d34c6eadc01dbecd21f349d93ec227b2c43a5f61eab7fff1e0382e8a9f61a2cce9cf8eb0730a697a98159"
    externalId = "111111"
    nodeName = "starry_sky"
    website = "https://www.test.network"
    details = "supper node"
    amount = 9999999
    programVersion, programVersionSign, blsProof = get_staking_program(ppos.web3)
    blsPubKey = "3f686e42718be2f7244c7c1ed6fce0aeb084f09d1668e6ad75aab672fb49d372dec42bede760a9c4e487fde485cd0f11994dc3223b355dd822ec499f00e67a4e0fe810dfe5aebac0a7ad3c8ad446e048ad83842719e227f183693a213680d098"
    result = ppos.createStaking(typ=0, benifitAddress=benifitAddress,
                                nodeId=nodeId, externalId=externalId,
                                nodeName=nodeName, website=website, details=details, amount=amount,
                                programVersion=programVersion, programVersionSign=programVersionSign, blsPubKey=blsPubKey, blsProof=blsProof,
                                )
    print(result)
    info = ppos.getCandidateInfo(nodeId)
    print(info)


def create_Restricting_Plan():
    url = "http://192.168.16.189:6789"
    stak_address = "0x2e95E3ce0a54951eB9A99152A6d5827872dFB4FD"
    stak_pri = "a689f0879f53710e9e0c1025af410a530d6381eebb5916773195326e123b822b"
    ppos = Ppos(url, stak_address, stak_pri)
    benifitAddress = "0xCeCA295e1471B3008D20b017c7Df7d4F338A7FbA"
    amount = ppos.web3.toWei(1000, 'ether')
    plan = [{'Epoch': 1, 'Amount': amount}]
    result = ppos.createRestrictingPlan(benifitAddress, plan, stak_pri)
    info = ppos.get_result(result)


if __name__ == "__main__":
    #create()
    create_Restricting_Plan()
