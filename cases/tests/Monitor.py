import time

import rlp
import json
import os
from client_sdk_python import HTTPProvider, Web3, WebsocketProvider
from client_sdk_python.eth import Eth
from client_sdk_python.middleware import geth_poa_middleware
from hexbytes import HexBytes
from smtplib import SMTP_SSL
from email.mime.text import MIMEText

config = {
    "url": "http://149.129.129.122:6789",
    "stak_address": "0x48c867ddBF22062704D6c81d3FA256bc6fc8b6bC",
    "stak_pri": "96f1f76c45bc2dd9c0f84a11da4ec104ae95661871284a29413a454a70b15307",
    "benifit_address": "0xCeCA295e1471B3008D20b017c7Df7d4F338A7FbA",
    "benifit_pri": "9614c2b32f2d5d3421591ab3ffc03ac66c831fb6807b532f6e3a8e7aac31f1d9",
    "node_pri": "71f7bc4797db4b7eb65c7e8648a11c0f45a1b2abbd0eb2ef309374aaa4f80297",
    "node_pub": "c7d843e4317a749071433ac021c7ec407bae2c61291b55986b0bc02e779ad154171db11d4236ae23b3cf4cdf064b3879c0ec0fed629234e88fad60fa816c6957",
    "bls_pri": "826e0495b92229f215064915ede29209253bd25135b96646d5d6f976e3b0584b",
    "bls_pub": "6cd1875fec5662384313702b7dc174234042417ea565993d29b4aa78e94629e99965b7332d684bdbecfb0201a1bf6013984048a223514adecc4d565ce0934c4f0494daf450c018761ded3602ab3a39238bcf988a5f5e4bd1635bdb1977954288",
    "ip": "13.54.77.34",
    "p2pport": "16798",
    "node_id": "c7d843e4317a749071433ac021c7ec407bae2c61291b55986b0bc02e779ad154171db11d4236ae23b3cf4cdf064b3879c0ec0fed629234e88fad60fa816c6957",
    "externalId": "111111",
    "nodeName": "starry_sky",
    "website": "https://www.test.network",
    "details": "supper node",
    "email_user": "230469827@qq.com",
    "email_authorization_code": "skeshksdjmsgcbbb"
}


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
    def __init__(self, url, address, privatekey, chainid=101):
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
        print('nonce: ', nonce)
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

    def apply_delegate(self, typ, nodeId, amount, privatekey=None, from_address=None, gasPrice=None, gas=None):
        to_address = "0x1000000000000000000000000000000000000002"
        data = rlp.encode([rlp.encode(int(1004)),
                           rlp.encode(int(typ)),
                           rlp.encode(bytes.fromhex(nodeId)),
                           rlp.encode(self.web3.toWei(amount, 'ether'))])
        if not privatekey:
            privatekey = self.privatekey
        if not from_address:
            from_address = self.address
        if not gasPrice:
            gasPrice = self.gasPrice
        if not gas:
            transactiondict = {"to": to_address, "data": data}
            gas = self.eth.estimateGas(transactiondict)
        result = self.send_raw_transaction(data, from_address, to_address, gasPrice, gas, 0, privatekey)
        return self.get_result(result)

    def getCandidateInfo(self, nodeId):
        data = rlp.encode([rlp.encode(int(1105)), rlp.encode(bytes.fromhex(nodeId))])
        to_address = "0x1000000000000000000000000000000000000002"
        raw_data = self.eth.call({
            "from": self.address,
            "to": to_address,
            "data": data
        })
        parse = str(raw_data, encoding="utf8").replace('\\', '').replace('"{', '{').replace('}"', '}')
        receive = json.loads(parse)
        try:
            raw_data_dict = receive["Ret"]
            raw_data_dict["Shares"] = int(raw_data_dict["Shares"], 16)
            raw_data_dict["Released"] = int(raw_data_dict["Released"], 16)
            raw_data_dict["ReleasedHes"] = int(raw_data_dict["ReleasedHes"], 16)
            raw_data_dict["RestrictingPlan"] = int(raw_data_dict["RestrictingPlan"], 16)
            raw_data_dict["RestrictingPlanHes"] = int(raw_data_dict["RestrictingPlanHes"], 16)
            raw_data_dict["DelegateRewardTotal"] = int(raw_data_dict["DelegateRewardTotal"], 16)
            raw_data_dict["DelegateTotal"] = int(raw_data_dict["DelegateTotal"], 16)
            raw_data_dict["DelegateTotalHes"] = int(raw_data_dict["DelegateTotalHes"], 16)
            receive["Ret"] = raw_data_dict
        except:
            ...
        return receive

    def getCandidateList(self):
        """
        Query all real-time candidate lists
        :param from_address: Used to call the rpc call method
        :return:
        todo fill
        """
        data = rlp.encode([rlp.encode(int(1102))])
        raw_data = self.eth.call({
            "from": self.address,
            "to": '0x1000000000000000000000000000000000000002',
            "data": data
        })
        parse = str(raw_data, encoding="utf8").replace('\\', '').replace('"{', '{').replace('}"', '}')
        try:
            raw_data = parse["Ret"]
            for i in raw_data:
                i["Shares"] = int(i["Shares"], 16)
                i["Released"] = int(i["Released"], 16)
                i["ReleasedHes"] = int(i["ReleasedHes"], 16)
                i["RestrictingPlan"] = int(i["RestrictingPlan"], 16)
                i["RestrictingPlanHes"] = int(i["RestrictingPlanHes"], 16)
        except:...
        return parse

    def getVerifierList(self):
        """
        Query the certified queue for the current billing cycle
        :param from_address: Used to call the rpc call method
        :return:
        todo fill
        """
        data = rlp.encode([rlp.encode(int(1100))])
        raw_data = self.eth.call({
            "from": self.address,
            "to": '0x1000000000000000000000000000000000000002',
            "data": data
        })
        parse = str(raw_data, encoding="utf8").replace('\\', '').replace('"{', '{').replace('}"', '}')
        try:
            raw_data = parse["Ret"]
            for i in raw_data:
                i["Shares"] = int(i["Shares"], 16)
        except:...
        return parse

def get_staking_program(w3):
    """
    根据rpc获取发质押的4个参数
    """
    msg = w3.admin.getProgramVersion()
    proof = w3.admin.getSchnorrNIZKProve()
    ProgramVersionSign = msg["Sign"]
    ProgramVersion = msg["Version"]
    return ProgramVersion, ProgramVersionSign, proof


def automatically_send_alert_emails(user, authorization_code):
    with SMTP_SSL(host="smtp.qq.com") as smtp:
        smtp.login(user=user, password=authorization_code)
        msg = MIMEText("The node has stopped exporting, please inform the relevant personnel to check", _charset="utf8")
        msg["Subject"] = "Monitoring node status is abnormal - alarm!!!"
        msg["from"] = user
        msg["to"] = user

        smtp.sendmail(from_addr=user, to_addrs=user, msg=msg.as_string())


def check_if_you_need_to_restart_node():
    """
    Check node status and whether the height of the block is normal
    :return:
    """

    last_query_block_height = 0
    ppos = Ppos(config['url'], config['stak_address'], config['stak_pri'])
    while True:
        pid = os.popen("ps -ef|grep platon|grep port|grep %s|grep -v grep|awk {'print $2'}")
        print('pid: ', pid.read())
        if pid is None:
            os.popen("nohup platon --identity 'platon' --datadir ./data --port 16798 --rpcport 6789 --verbosity 1 --rpcapi 'platon,web3,admin' --rpc  --nodekey ./data/nodekey --cbft.blskey ./data/nodeblskey >> platon.log 2>&1 &")
        else:
            print("Current process: running normally PID: ", pid.read())
        current_query_block_height = ppos.eth.blockNumber
        print('current query block height: ', current_query_block_height)
        if last_query_block_height == current_query_block_height:
            automatically_send_alert_emails(config['email_user'], config['email_authorization_code'])
        last_query_block_height = current_query_block_height
        # Automatic delegation
        node_revenue_amount = ppos.eth.getBalance(ppos.web3.toChecksumAddress(config['benifit_address']))
        print('Current income address balance： ', node_revenue_amount)
        if node_revenue_amount > ppos.web3.toWei(100, 'ether'):
            amount = node_revenue_amount - ppos.web3.toWei(1, 'ether')
            print('Application entrustment amount： ', amount)
            ppos.apply_delegate(0, config['node_id'], amount, config['benifit_pri'], config['benifit_pri'])
        time.sleep(5)


def automatic_account_entrustment(amount=None):
    ppos = Ppos(config['url'], config['stak_address'], config['stak_pri'])
    if amount is None:
        node_revenue_amount = ppos.eth.getBalance(ppos.web3.toChecksumAddress(config['benifit_address']))
        print('Current income address balance： ', node_revenue_amount)
        if node_revenue_amount > ppos.web3.toWei(100, 'ether'):
            amount = node_revenue_amount - ppos.web3.toWei(1, 'ether')
    print('Application entrustment amount： ', amount)
    ppos.apply_delegate(0, config['node_id'], amount, config['benifit_pri'], config['benifit_pri'])

def gettttt(nodeid):
    ppos = Ppos(config['url'], config['stak_address'], config['stak_pri'])
    result = ppos.getCandidateInfo(nodeid)
    print(result)

def getsssss():
    ppos = Ppos(config['url'], config['stak_address'], config['stak_pri'])
    result = ppos.getCandidateList()
    print(result)

def getaaaa():
    ppos = Ppos(config['url'], config['stak_address'], config['stak_pri'])
    result = ppos.getVerifierList()
    print(result)

if __name__ == '__main__':
    # automatically_send_alert_emails('230469827@qq.com','skeshksdjmsgcbbb')
    # check_if_you_need_to_restart_node()
    # automatic_account_entrustment()
    gettttt('df672cbf413dc740036ef5bf2545180fe9309561189b22f220b77beae7fce61cba7faee821433eab0fce29b3a564288629cff5dee314195c33f601c165d1949a')
    # getsssss()
    # getaaaa()