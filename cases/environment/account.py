from common.load_file import LoadFile
from client_sdk_python.eth import Eth
from client_sdk_python.personal import Personal
from hexbytes import HexBytes
import random
import rlp
import os


class Account:
    def __init__(self, accountFile, chainId):
        """
        accounts 包含的属性: address,prikey,nonce,balance,node_id,passwd
        :param accountFile:
        :param chainId:
        """
        self.accounts = {}
        self.raw_accounts = LoadFile(accountFile).get_data()
        # log.info(self.raw_accounts)
        self.chain_id = chainId
        self.account_with_money = self.raw_accounts[0]
        self.reset()
        # for account in self.raw_accounts:
        #     self.accounts[account['address']] = account

    def reset(self):
        self.accounts = {}
        for account in self.raw_accounts:
            self.accounts[account['address']] = account

    def get_all_accounts(self):
        accounts = []
        for account in self.accounts.values():
            accounts.append(account)
        return accounts

    def get_rand_account(self):
        # todo 实现随机
        # for account in self.accounts.values():
        #     return account
        return random.choice(list(self.accounts.values()))

    def sendTransaction(self, connect, data, from_address, to_address, gasPrice, gas, value, check_address=True):
        platon = Eth(connect)
        account = self.accounts[from_address]
        nonce = platon.getTransactionCount(from_address)
        transaction_dict = {
            "to": to_address,
            "gasPrice": gasPrice,
            "gas": gas,
            "nonce": nonce,
            "data": data,
            "chainId": self.chain_id,
            "value": value,
        }

        signedTransactionDict = platon.account.signTransaction(
            transaction_dict, account['prikey']
        )

        data = signedTransactionDict.rawTransaction
        result = HexBytes(platon.sendRawTransaction(data)).hex()
        # log.debug("result:::::::{}".format(result))
        res = platon.waitForTransactionReceipt(result)
        account['nonce'] = nonce + 1
        self.accounts[from_address] = account
        return res

    def generate_account_in_node(self, node, passwd, balance=0):
        personal = Personal(node.web3)
        address = personal.newAccount(passwd)
        # log.debug(address)
        if balance > 0:
            self.sendTransaction(node.web3, '', self.account_with_money['address'], address, node.eth.gasPrice, 40000, balance)
        account = {
            "node_id": node.node_id,
            "address": address,
            "nonce": 0,
            "balance": balance,
            "prikey": '',
            'passwd': passwd
        }
        self.accounts[address] = account
        return address

    def unlock_account(self, node, address):
        account = self.accounts[address]
        personal = Personal(node.web3)
        personal.unlockAccount(account['address'], account['passwd'])

    def get_rand_account_in_node(self, node):
        for account in self.accounts.values():
            if account['node_id'] == node.id:
                return account
        self.generate_account_in_node(node, '123456')

    def create_restricting_plan(self, connect, receive_address, plan, from_address, gasPrice=None, gas=None):
        '''
        创建锁仓计划
        :param account: 20bytes
        :param plan: []RestrictingPlan
        :param from_address:
        :param gasPrice:
        :param gas:
        :return:
        '''
        to_address = "0x1000000000000000000000000000000000000001"
        if receive_address[:2] == '0x':
            receive_address = receive_address[2:]
        plan_list = []
        for dict_ in plan:
            v = [dict_[k] for k in dict_]
            plan_list.append(v)
        rlp_list = rlp.encode(plan_list)
        data = rlp.encode([rlp.encode(int(4000)),
                           rlp.encode(bytes.fromhex(receive_address)),
                           rlp_list])
        # print ("len:", len (data))
        # l = [hex (int (i)) for i in data]
        # print (" ".join (l))
        result = self.sendTransaction(connect, data, from_address, to_address, gasPrice, gas, 0)
        return result

    def generate_account(self, web3, balance=0):
        platon = Eth(web3)
        account = platon.account.create(net_type=web3.net_type)
        address = account.address

        prikey = account.privateKey.hex()[2:]
        if balance != 0:
            self.sendTransaction(web3, '', self.account_with_money['address'], address, web3.platon.gasPrice, 21000, balance)
        account = {
            "address": address,
            "nonce": 0,
            "balance": balance,
            "prikey": prikey,
        }
        self.accounts[address] = account
        # todo delete debug

        def debug():
            from conf.settings import BASE_DIR
            from ruamel import yaml
            accounts = list(self.accounts.values())
            with open(os.path.join(BASE_DIR, "deploy/tmp/accounts.yml"), mode="w", encoding="UTF-8") as f:
                yaml.dump(accounts, f, Dumper=yaml.RoundTripDumper)
        debug()
        return address, prikey

    def find_pri_key(self, address):
        return self.accounts[address]["prikey"]
