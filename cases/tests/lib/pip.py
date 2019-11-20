from environment.env import TestEnvironment
from environment.node import Node
from .config import PipConfig
from .economic import Economic
from .utils import int_to_bytes, get_blockhash, proposal_list_effective, proposal_effective, find_proposal
import json
from typing import List
import time


class Pip:
    """
    Used to initiate a pip transaction,
    if you need to use the call method, please call pip
    example:
    >>>pip=Pip(env, node)
    >>>pip.pip.getActiveVersion(...)
    """
    cfg = PipConfig

    def __init__(self, env: TestEnvironment, node: Node):
        self.node = node
        self.economic = Economic(env)

    def submitText(self, verifier, pip_id, from_address, transaction_cfg=None):
        """
        Submit a text proposal
        :param verifier: The certified submitting the proposal
        :param pip_id: PIPID
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
        pri_key = self.economic.account.find_pri_key(from_address)
        return self.pip.submitText(verifier, pip_id, pri_key, transaction_cfg)

    def submitVersion(self, verifier, pip_id, new_version, end_voting_rounds, from_address, transaction_cfg=None):
        """
        Submit an upgrade proposal
        :param verifier:  The certified submitting the proposal
        :param pip_id:  PIPID
        :param new_version: upgraded version
        :param end_voting_rounds: The number of voting consensus rounds.
            Explanation: Assume that the transaction submitted by the proposal is rounded when the consensus round
            number of the package is packed into the block, then the proposal voting block is high,
            which is the 230th block height of the round of the round1 + endVotingRounds
            (assuming a consensus round out of block 250, ppos The list is 20 blocks high in advance,
             250, 20 are configurable), where 0 < endVotingRounds <= 4840 (about 2 weeks, the actual discussion
             can be calculated according to the configuration), and is an integer)
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
        pri_key = self.economic.account.find_pri_key(from_address)
        return self.pip.submitVersion(verifier, pip_id, new_version, end_voting_rounds, pri_key, transaction_cfg)

    def submitParam(self, verifier, pip_id, module, name, new_value, from_address, transaction_cfg=None):
        """
        Submit an param proposal
        :param verifier: The certified submitting the proposal
        :param pip_id: PIPID
        :param module: parameter module
        :param name: parameter name
        :param new_value: New parameter value
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
        pri_key = self.economic.account.find_pri_key(from_address)
        return self.pip.submitParam(verifier, pip_id, module, name, new_value, pri_key, transaction_cfg)

    def submitCancel(self, verifier, pip_id, end_voting_rounds, tobe_canceled_proposal_id, from_address, transaction_cfg=None):
        """
        Submit cancellation proposal
        :param verifier: The certified submitting the proposal
        :param pip_id: PIPID
        :param end_voting_rounds:
           The number of voting consensus rounds. Refer to the instructions for submitting the upgrade proposal.
           At the same time, the value of this parameter in this interface
           cannot be greater than the value in the corresponding upgrade proposal.
        :param tobe_canceled_proposal_id: Upgrade proposal ID to be cancelled
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
        pri_key = self.economic.account.find_pri_key(from_address)
        return self.pip.submitCancel(verifier, pip_id, end_voting_rounds, tobe_canceled_proposal_id, pri_key, transaction_cfg)

    def vote(self, verifier, proposal_id, option, from_address, program_version=None, version_sign=None, transaction_cfg=None):
        """
        Vote for proposal
        :param verifier:  The certified submitting the proposal
        :param proposal_id: Proposal ID
        :param option: Voting option
        :param program_version: Node code version, obtained by rpc getProgramVersion interface
        :param version_sign: Code version signature, obtained by rpc getProgramVersion interface
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
        pri_key = self.economic.account.find_pri_key(from_address)
        if program_version is None:
            program_version = self.node.program_version
        if version_sign is None:
            version_sign = self.node.program_version_sign
        return self.pip.vote(verifier, proposal_id, option, program_version, version_sign, pri_key, transaction_cfg)

    def declareVersion(self, active_node, from_address, program_version=None, version_sign=None, transaction_cfg=None):
        """
        Version statement
        :param active_node: The declared node can only be a verifier/candidate
        :param program_version: The declared version, obtained by rpc's getProgramVersion interface
        :param version_sign: The signed version signature, obtained by rpc's getProgramVersion interface
        :param from_address: address transaction
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
        pri_key = self.economic.account.find_pri_key(from_address)
        if program_version is None:
            program_version = self.node.program_version
        if version_sign is None:
            version_sign = self.node.program_version_sign
        return self.pip.declareVersion(active_node, program_version, version_sign, pri_key, transaction_cfg)

    @property
    def pip(self):
        """
        use sdk pip
        :return:
        """
        return self.node.pip

    def get_status_of_proposal(self, proposal_id):
        """
        Get proposal voting results
        :param proposal_id:
        :return:
        """
        result = self.pip.getTallyResult(proposal_id)
        data = result.get('Ret')
        # data = json.loads(data)
        if not data:
            raise Exception('Failed to query proposal result based on given proposal id')
        return data.get('status')

    def get_accu_verifiers_of_proposal(self, proposal_id):
        """
        Get the total number of certifiers who have voted for the entire voting period
        :param proposal_id:
        :return:
        """
        result = self.pip.getTallyResult(proposal_id)
        resultinfo = result.get('Ret')
        # resultinfo = json.loads(resultinfo)
        if not resultinfo:
            raise Exception('Failed to query proposal result based on given proposal id')
        return resultinfo.get('accuVerifiers')

    def get_yeas_of_proposal(self, proposal_id):
        """
        Get the number of people who voted for the entire voting period
        :param proposal_id:
        :return:
        """
        result = self.pip.getTallyResult(proposal_id)
        data = result.get('Ret')
        # data = json.loads(data)
        if not data:
            raise Exception('Failed to query proposal result based on given proposal id')
        return data.get('yeas')

    def get_nays_of_proposal(self, proposal_id):
        """
        Get the number of votes against the entire voting period
        :param proposal_id:
        :return:
        """
        result = self.pip.getTallyResult(proposal_id)
        data = result.get('Ret')
        # data = json.loads(data)
        if not data:
            raise Exception('Failed to query proposal result based on given proposal id')
        return data.get('nays')

    def get_abstentions_of_proposal(self, proposal_id):
        """
        Obtain the number of abstentions during the entire voting period
        :param proposal_id:
        :return:
        """
        result = self.pip.getTallyResult(proposal_id)
        data = result.get('Ret')
        # data = json.loads(data)
        if not data:
            raise Exception('Failed to query proposal result based on given proposal id')
        return data.get('abstentions')

    def get_canceledby_of_proposal(self, proposal_id):
        """
        Obtain the number of abstentions during the entire voting period
        :param proposal_id:
        :return:
        """
        result = self.pip.getTallyResult(proposal_id)
        data = result.get('Ret')
        # data = json.loads(data)
        if not data:
            raise Exception('Failed to query proposal result based on given proposal id')
        return data.get('canceledBy')

    @property
    def chain_version(self):
        """
        Get the version number on the chain
        :return:
        """
        result = self.pip.getActiveVersion()
        return int(result.get('Ret'))

    def get_version_small_version(self, flag=3):
        """
        Determine if the minor version of the incoming version number is 0
        :param flag:
        :return:
        """
        flag = int(flag)
        if flag > 3 or flag < 1:
            raise Exception("Incorrect parameters passed in")
        version = int(self.chain_version)
        version_byte = int_to_bytes(version)
        return version_byte[flag]

    def get_accuverifiers_count(self, proposal_id, blocknumber=None):
        """
        Get proposal real-time votes
        :param proposal_id:
        :param blocknumber:
        :return:
        """
        if blocknumber is None:
            blocknumber = self.node.block_number
        blockhash = get_blockhash(self.node, blocknumber)
        result = self.pip.getAccuVerifiersCount(proposal_id, blockhash)
        voteinfo = result.get('Ret')
        # vote_result = eval(voteinfo)
        return voteinfo

    def get_rate_of_voting(self, proposal_id):
        """
        Calculate the voting rate of the upgrade proposal
        :param proposal_id:
        :return:
        """
        result = self.pip.getTallyResult(proposal_id).get('Ret')
        # result = json.loads(result)
        if not result:
            raise Exception('Failed to query proposal result based on given proposal id')
        yeas = result.get('yeas')
        accu_verifiers = result.get('accuVerifiers')
        return yeas / accu_verifiers

    def get_effect_proposal_info_of_preactive(self):
        """
        Get pre-valid proposal information on the chain
        :return:
        """
        result = self.pip.listProposal().get('Ret')
        # result = json.loads(result)
        for pid_list in result:
            if pid_list.get('ProposalType') == 2:
                if self.get_status_of_proposal(pid_list.get('ProposalID')) == 4:
                    return pid_list
        raise Exception('There is no pre-validation upgrade proposal')

    def get_effect_proposal_info_of_vote(self, proposaltype=cfg.version_proposal):
        """
        Get pre-valid proposal information on the chain
        :return:
        """
        if not self.is_exist_effective_proposal_for_vote(self.cfg.text_proposal) and proposaltype == self.cfg.text_proposal:
            return None

        if not self.is_exist_effective_proposal_for_vote(self.cfg.version_proposal) and proposaltype == self.cfg.version_proposal:
            return None

        if not self.is_exist_effective_proposal_for_vote(self.cfg.cancel_proposal) and proposaltype == self.cfg.cancel_proposal:
            return None

        if not self.is_exist_effective_proposal_for_vote(self.cfg.param_proposal) and proposaltype == self.cfg.param_proposal:
            return None

        proposal_info = self.pip.listProposal().get('Ret')
        # proposal_info = json.loads(proposal_info)
        proposal_list_text = []
        proposal_list_version = []
        proposal_list_param = []
        proposal_list_cancel = []
        for pid_list in proposal_info:
            if pid_list.get('ProposalType') == self.cfg.text_proposal:
                proposal_list_text.append(pid_list)

            elif pid_list.get('ProposalType') == self.cfg.version_proposal:
                proposal_list_version.append(pid_list)

            elif pid_list.get('ProposalType') == self.cfg.cancel_proposal:
                proposal_list_cancel.append(pid_list)

            elif pid_list.get('ProposalType') == self.cfg.param_proposal:
                proposal_list_param.append(pid_list)
            else:
                raise Exception("Unknown proposal type")
        # Current block height
        block_number = self.node.eth.blockNumber
        if proposaltype == self.cfg.text_proposal:
            return find_proposal(proposal_list_text, block_number)

        elif proposaltype == self.cfg.version_proposal:
            return find_proposal(proposal_list_version, block_number)

        elif proposaltype == self.cfg.cancel_proposal:
            return find_proposal(proposal_list_cancel, block_number)

        elif proposaltype == self.cfg.param_proposal:
            return find_proposal(proposal_list_param, block_number)
        else:
            raise Exception("listProposal interface gets the wrong proposal type")

    def get_proposal_info_list(self):
        """
        Get a list of proposals
        :return:
        """
        proposal_info_list = self.pip.listProposal().get('Ret')
        version_proposal_list, text_proposal_list, cancel_proposal_list, param_proposal_list = [], [], [], []
        if proposal_info_list != 'Object not found':
            # proposal_info_list = json.loads(proposal_info_list)
            for proposal_info in proposal_info_list:
                if proposal_info.get('ProposalType') == self.cfg.version_proposal:
                    version_proposal_list.append(proposal_info)
                elif proposal_info.get('ProposalType') == self.cfg.text_proposal:
                    text_proposal_list.append(proposal_info)
                elif proposal_info.get('ProposalType') == self.cfg.cancel_proposal:
                    cancel_proposal_list.append(proposal_info)
                elif proposal_info.get('ProposalType') == self.cfg.param_proposal:
                    param_proposal_list.append(proposal_info)
                else:
                    raise Exception('listProposal interface gets the wrong proposal type')
        return version_proposal_list, text_proposal_list, cancel_proposal_list, param_proposal_list

    def is_exist_effective_proposal(self, proposal_type=None):
        """
        Determine if there is a valid upgrade proposal on the chain - to determine if a proposal can be initiated
                                :param proposal_type: 2 is the upgrade proposal 1. Text proposal 4. Cancel the proposal
                                :return:
        """
        if proposal_type is None:
            proposal_type = self.cfg.version_proposal
        version_proposal_list, text_proposal_list, cancel_proposal_list, param_proposal_list = self.get_proposal_info_list()
        block_number = self.node.eth.blockNumber
        if proposal_type == self.cfg.version_proposal:
            for version_proposal in version_proposal_list:
                if proposal_effective(version_proposal, block_number):
                    return True
                else:
                    status = self.get_status_of_proposal(version_proposal["ProposalID"].strip())
                    if status == 4:
                        return True

        elif proposal_type == self.cfg.text_proposal:
            return proposal_list_effective(text_proposal_list, block_number)

        elif proposal_type == self.cfg.cancel_proposal:
            return proposal_list_effective(cancel_proposal_list, block_number)

        elif proposal_type == self.cfg.param_proposal:
            return proposal_list_effective(param_proposal_list, block_number)
        else:
            raise Exception("Incoming type error")
        return False

    def is_exist_effective_proposal_for_vote(self, proposal_type=None):
        """
        Is there a valid proposal for voting?
        :param proposal_type:
        :return:
        """
        if proposal_type is None:
            proposal_type = self.cfg.version_proposal
        version_proposal_list, text_proposal_list, cancel_proposal_list, param_proposal_list = self.get_proposal_info_list()
        block_number = self.node.eth.blockNumber
        if proposal_type == self.cfg.version_proposal:
            return proposal_list_effective(version_proposal_list, block_number)

        elif proposal_type == self.cfg.text_proposal:
            return proposal_list_effective(text_proposal_list, block_number)

        elif proposal_type == self.cfg.cancel_proposal:
            return proposal_list_effective(cancel_proposal_list, block_number)

        elif proposal_type == self.cfg.param_proposal:
            return proposal_list_effective(param_proposal_list, block_number)
        else:
            raise Exception("Incoming type error")

    def get_candidate_list_not_verifier(self):
        """
        获取当前结算周期非验证人的候选人列表
        :return:
        """
        candidate_list = self.node.ppos.getCandidateList().get('Ret')
        verifier_list = self.node.ppos.getVerifierList().get('Ret')
        if verifier_list == "Getting verifierList is failed:The validator is not exist":
            time.sleep(10)
            verifier_list = self.node.ppos.getVerifierList().get('Ret')
        candidate_no_verify_list = []
        verifier_node_list = [node_info.get("NodeId") for node_info in verifier_list]
        for node_info in candidate_list:
            node_id = node_info.get("NodeId")
            if node_id not in verifier_node_list:
                candidate_no_verify_list.append(node_id)
        return candidate_no_verify_list

    def get_version(self, version=None):
        # todo implement
        pass


def get_pip_obj(nodeid, pip_obj_list: List[Pip]) -> Pip:
    """
    Get the pip object according to the node id
    :param nodeid:
    :param pip_obj_list:
    :return:
    """
    for pip_obj in pip_obj_list:
        if nodeid == pip_obj.node.node_id:
            return pip_obj


def get_pip_obj_list(nodeid_list, pip_obj_list: List[Pip]) -> List[Pip]:
    """
    Get a list of pip objects based on the node id list
    :param nodeid_list:
    :param pip_obj_list:
    :return:
    """
    new_pip_obj_list = []
    for nodeid in nodeid_list:
        new_pip_obj_list.append(get_pip_obj(nodeid, pip_obj_list))
    return new_pip_obj_list
