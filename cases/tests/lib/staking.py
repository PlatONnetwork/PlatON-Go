from environment.env import TestEnvironment
from environment.node import Node
from .config import StakingConfig
from .economic import Economic
import time


class Staking:
    """
    Used to initiate a Staking transaction,
    if you need to use the call method, please call ppos
    example:
    >>>staking=Staking(env, node)
    >>>staking.ppos.getValidatorList()
    """

    def __init__(self, env: TestEnvironment, node: Node, cfg: StakingConfig):
        self.cfg = cfg
        self.node = node
        self.economic = Economic(env)

    @property
    def ppos(self):
        return self.node.ppos

    def create_staking(self, typ, benifit_address, from_address, node_id=None, amount=None, program_version=None,
                       program_version_sign=None, bls_pubkey=None, bls_proof=None, transaction_cfg=None, reward_per=0):
        """
        Initiate Staking
        :param typ: Indicates whether the account free amount or the account's lock amount is used for staking, 0: free amount; 1: lock amount
        :param benifit_address: Income account for accepting block rewards and staking rewards
        :param node_id: The idled node Id (also called the candidate's node Id)
        :param amount: staking von (unit:von, 1LAT = 10**18 von)
        :param program_version: The real version of the program, admin_getProgramVersion
        :param program_version_sign: The real version of the program is signed, admin_getProgramVersion
        :param bls_pubkey: Bls public key
        :param bls_proof: Proof of bls, obtained by pulling the proof interface, admin_getSchnorrNIZKProve
        :param from_address: address for transaction
        :param transaction_cfg: Transaction basic configuration
              type: dict
              example:cfg = {
                  "gas":100000000,
                  "gasPrice":2000000000000,
                  "nonce":1,
              }
        :param reward_per: Proportion of the reward share obtained from the commission, using BasePoint 1BP = 0.01%
        :return: if is need analyze return transaction result dict
                if is not need analyze return transaction hash
        """
        if node_id is None:
            node_id = self.node.node_id
        if amount is None:
            amount = self.economic.create_staking_limit
        if program_version is None:
            program_version = self.node.program_version
        if program_version_sign is None:
            program_version_sign = self.node.program_version_sign
        if bls_pubkey is None:
            bls_pubkey = self.node.blspubkey
        if bls_proof is None:
            bls_proof = self.node.schnorr_NIZK_prove
        pri_key = self.economic.account.find_pri_key(from_address)
        return self.ppos.createStaking(typ, benifit_address, node_id, self.cfg.external_id, self.cfg.node_name,
                                       self.cfg.website, self.cfg.details, amount, program_version, program_version_sign,
                                       bls_pubkey, bls_proof, pri_key, reward_per, transaction_cfg=transaction_cfg)

    def edit_candidate(self, from_address, benifit_address, node_id=None, transaction_cfg=None, reward_per=0):
        """
        Modify staking information
        :param benifit_address: Income account for accepting block rewards and staking rewards
        :param node_id: The idled node Id (also called the candidate's node Id)
        :param from_address: address for transaction
        :param transaction_cfg: Transaction basic configuration
              type: dict
              example:cfg = {
                  "gas":100000000,
                  "gasPrice":2000000000000,
                  "nonce":1,
              }
        :param reward_per: Proportion of the reward share obtained from the commission, using BasePoint 1BP = 0.01%
        :return: if is need analyze return transaction result dict
                if is not need analyze return transaction hash
        """
        if node_id is None:
            node_id = self.node.node_id
        pri_key = self.economic.account.find_pri_key(from_address)
        return self.ppos.editCandidate(benifit_address, node_id, self.cfg.external_id, self.cfg.node_name, self.cfg.website, self.cfg.details,
                                       pri_key, reward_per, transaction_cfg=transaction_cfg)

    def increase_staking(self, typ, from_address, node_id=None, amount=None, transaction_cfg=None):
        """
        Increase staking
        :param typ: Indicates whether the account free amount or the account's lock amount is used for staking, 0: free amount; 1: lock amount
        :param node_id: The idled node Id (also called the candidate's node Id)
        :param amount: staking von (unit:von, 1LAT = 10**18 von)
        :param from_address: address for transaction
        :param transaction_cfg: Transaction basic configuration
              type: dict
              example:cfg = {
                  "gas":100000000,
                  "gasPrice":2000000000000,
                  "nonce":1,
              }
        :return: if is need analyze return transaction result dict
                if is not need analyze return transaction hash
        """
        if node_id is None:
            node_id = self.node.node_id
        if amount is None:
            amount = self.economic.add_staking_limit
        pri_key = self.economic.account.find_pri_key(from_address)
        return self.ppos.increaseStaking(typ, node_id, amount, pri_key, transaction_cfg=transaction_cfg)

    def withdrew_staking(self, from_address, node_id=None, transaction_cfg=None):
        """
        Withdrawal of staking (one-time initiation of all cancellations, multiple arrivals)
        :param node_id: The idled node Id (also called the candidate's node Id)
        :param from_address: address for transaction
        :param transaction_cfg: Transaction basic configuration
              type: dict
              example:cfg = {
                  "gas":100000000,
                  "gasPrice":2000000000000,
                  "nonce":1,
              }
        :return: if is need analyze return transaction result dict
                if is not need analyze return transaction hash
        """
        if node_id is None:
            node_id = self.node.node_id
        pri_key = self.economic.account.find_pri_key(from_address)
        return self.ppos.withdrewStaking(node_id, pri_key, transaction_cfg=transaction_cfg)

    def get_staking_address(self):
        """
        Get the pledge wallet address
        """
        result = self.ppos.getCandidateInfo(self.node.node_id)
        candidate_info = result.get('Ret')
        address = candidate_info.get('StakingAddress')
        return self.node.web3.toChecksumAddress(address)

    def get_candidate_list_not_verifier(self):
        """
        Get a list of candidates for non-verifiers for the current billing cycle
        """
        candidate_list = self.ppos.getCandidateList().get('Ret')
        verifier_list = self.ppos.getVerifierList().get('Ret')
        if verifier_list == "Getting verifierList is failed:The validator is not exist":
            time.sleep(10)
            verifier_list = self.ppos.getVerifierList().get('Ret')
        candidate_no_verify_list = []
        verifier_node_list = [node_info.get("NodeId") for node_info in verifier_list]
        for node_info in candidate_list:
            node_id = node_info.get("NodeId")
            if node_id not in verifier_node_list:
                candidate_no_verify_list.append(node_id)
        return candidate_no_verify_list

    def get_staking_amount(self, node=None, flag=0):
        """
        According to the node to obtain the amount of the deposit
        """
        if node is None:
            node = self.node
        flag = int(flag)
        stakinginfo = node.ppos.getCandidateInfo(node.node_id)
        staking_data = stakinginfo.get('Ret')
        shares = int(staking_data.get('Shares'))
        released = int(staking_data.get('Released'))
        restrictingplan = int(staking_data.get('RestrictingPlan'))
        return [shares, released, restrictingplan][flag]

    def get_rewardper(self, node=None, isnext=False):
        '''
        According to the node to obtain the reward percent
        :param node:
        :param isnext:
        :return:
        '''
        if node is None:
            node = self.node
        stakinginfo = node.ppos.getCandidateInfo(node.node_id)
        print(stakinginfo)
        staking_data = stakinginfo.get('Ret')
        rewardper = int(staking_data.get('RewardPer'))
        nextrewardper = int(staking_data.get('NextRewardPer'))
        if isnext:
            return nextrewardper
        else:
            return rewardper


    def get_version(self, node=None):
        """
        According to the node to obtain the amount of the deposit
        """
        if node is None:
            node = self.node
        stakinginfo = node.ppos.getCandidateInfo(node.node_id)
        staking_data = stakinginfo.get('Ret')
        programversion = staking_data.get('ProgramVersion')
        return programversion

    def get_stakingblocknum(self, node=None):
        """
        According to the node to obtain the amount of the deposit
        """
        if node is None:
            node = self.node
        stakinginfo = node.ppos.getCandidateInfo(node.node_id)
        staking_data = stakinginfo.get('Ret')
        stakingblocknum = staking_data.get('StakingBlockNum')
        return int(stakingblocknum)
