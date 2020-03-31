from decimal import Decimal

from dacite import from_dict
from .utils import wait_block_number, get_pledge_list
from environment.node import Node
from .genesis import Genesis
from common.key import get_pub_key
import math
from .config import EconomicConfig
from environment.env import TestEnvironment


class Economic:
    cfg = EconomicConfig

    def __init__(self, env: TestEnvironment):
        self.env = env

        self.genesis = from_dict(data_class=Genesis, data=self.env.genesis_config)

        # Block rate parameter
        self.per_round_blocks = self.genesis.config.cbft.amount
        self.interval = int((self.genesis.config.cbft.period / self.per_round_blocks) / 1000)

        # Length of additional issuance cycle
        self.additional_cycle_time = self.genesis.economicModel.common.additionalCycleTime

        # Number of verification
        self.validator_count = self.genesis.economicModel.common.maxConsensusVals

        # Billing related
        # Billing cycle
        self.expected_minutes = self.genesis.economicModel.common.maxEpochMinutes
        # Consensus rounds
        # self.consensus_wheel = (self.expected_minutes * 60) // (
        #     self.interval * self.per_round_blocks * self.validator_count)
        # Number of settlement periods
        # self.settlement_size = self.consensus_wheel * (self.interval * self.per_round_blocks * self.validator_count)
        # Consensus round number
        # self.consensus_size = self.per_round_blocks * self.validator_count

        # Minimum amount limit
        self.create_staking_limit = self.genesis.economicModel.staking.stakeThreshold
        # Minimum deposit amount
        # Minimum holding amount
        self.add_staking_limit = self.genesis.economicModel.staking.operatingThreshold
        # Minimum commission amount
        self.delegate_limit = self.add_staking_limit
        # unstaking freeze duration
        self.unstaking_freeze_ratio = self.genesis.economicModel.staking.unStakeFreezeDuration
        # ParamProposalVote_DurationSeconds
        self.pp_vote_settlement_wheel = self.genesis.economicModel.gov.paramProposalVoteDurationSeconds // self.settlement_size
        # slash blocks reward
        self.slash_blocks_reward = self.genesis.economicModel.slashing.slashBlocksReward
        # text proposal vote duration senconds
        self.tp_vote_settlement_wheel = self.genesis.economicModel.gov.textProposalVoteDurationSeconds // (
            self.interval * self.per_round_blocks * self.validator_count)

    @property
    def consensus_wheel(self):
        return (self.expected_minutes * 60) // (
                self.interval * self.per_round_blocks * self.validator_count)

    @property
    def consensus_size(self):
        return self.per_round_blocks * self.validator_count

    @property
    def settlement_size(self):
        return self.consensus_wheel * (self.interval * self.per_round_blocks * self.validator_count)


    @property
    def account(self):
        return self.env.account

    def get_block_count_number(self, node: Node, roundnum=1):
        """
        Get the number of blocks out of the verification node
        """
        current_block = node.eth.blockNumber
        block_namber = self.consensus_size * roundnum
        count = 0
        for i in range(block_namber):
            if current_block > 0:
                node_id = node.eth.ecrecover(current_block)
                # node_id = get_pub_key(node.url, current_block)
                current_block = current_block - 1
                if node_id == node.node_id:
                    count = count + 1
            else:
                break
        print(current_block)
        return count

    def get_number_blocks_in_interval(self, node: Node, roundnum=1):
        """
        Get the number of blocks produced by the specified interval of the node
        """
        tmp_current_block = node.eth.blockNumber
        last_end_block = int(tmp_current_block / self.settlement_size) * self.settlement_size
        block_number = self.settlement_size * roundnum
        count = 0
        for i in range(block_number):
            node_id = node.eth.ecrecover(last_end_block)
            # node_id = get_pub_key(node.url, last_end_block)
            last_end_block = last_end_block - 1
            if node_id == node.node_id:
                count = count + 1
        return count

    def calculate_delegate_reward(self, node, block_reward, staking_reward, reward=None):
        block_number = self.get_number_blocks_in_interval(node)
        if reward is None:
            reward = node.ppos.getCandidateInfo(node.node_id)["Ret"]["RewardPer"]
        return int(Decimal(str(staking_reward))*Decimal(str(reward))/Decimal(str(10000)) + Decimal(str(int(Decimal(str(block_reward))*Decimal(str(reward))/Decimal(str(10000))))) * Decimal(str(block_number)))

    def delegate_cumulative_income(self, node, block_reward, staking_reward, delegate_total_amount, delegate_amount, reward=None):
        entrusted_income = self.calculate_delegate_reward(node, block_reward, staking_reward, reward)
        print("entrusted_income: ", entrusted_income)
        # unit_commission_award = math.floor(Decimal(str(entrusted_income)) / int((Decimal(str(delegate_total_amount)) / (10 ** 9))))
        # print("unit_commission_award: ", unit_commission_award)
        # current_commission_award = int((Decimal(str(delegate_amount)) / (10 ** 9)) * Decimal(str(unit_commission_award)))
        # print("current_commission_award: ", current_commission_award)
        current_commission_award = math.floor(Decimal(str(entrusted_income)) * Decimal(str(delegate_amount)) / Decimal(str(delegate_total_amount)))
        print("current_commission_award: ", current_commission_award)
        return current_commission_award

    def delegate_dividend_income(self, delegate_reward_total, delegate_total_amount, delegate_amount):
        # unit_commission_award = math.floor(Decimal(str(delegate_reward_total)) / int((Decimal(str(delegate_total_amount)) / (10 ** 9))))
        # print("unit_commission_award: ", unit_commission_award)
        # current_commission_award = int((Decimal(str(delegate_amount)) / (10 ** 9)) * Decimal(str(unit_commission_award)))
        # print("current_commission_award: ", current_commission_award)
        # return current_commission_award
        current_commission_award = math.floor(Decimal(str(delegate_reward_total)) * Decimal(str(delegate_amount)) / Decimal(str(delegate_total_amount)))
        print("current_commission_award: ", current_commission_award)
        return current_commission_award

    def get_current_year_reward(self, node: Node, verifier_num=None):
        """
        Get the first year of the block reward, pledge reward
        :return:
        """
        # if new_block_rate is None:
        #     new_block_rate = self.genesis.economicModel.reward.newBlockRate
        # # current_block = node.eth.blockNumber
        # annualcycle = (self.additional_cycle_time * 60) // self.settlement_size
        # annual_size = annualcycle * self.settlement_size
        # # starting_block_height = math.floor(current_block / annual_size) * annual_size
        if verifier_num is None:
            verifier_list = get_pledge_list(node.ppos.getVerifierList)
            verifier_num = len(verifier_list)
        # # amount = node.eth.getBalance(self.cfg.INCENTIVEPOOL_ADDRESS, starting_block_height)
        # if amount is None:
        #     amount = 262215742000000000000000000
        # block_proportion = str(new_block_rate / 100)
        # staking_proportion = str(1 - new_block_rate / 100)
        # block_reward = int(Decimal(str(amount)) * Decimal(str(block_proportion)) / Decimal(str(annual_size))) - node.web3.toWei(1 , 'ether')
        # staking_reward = int(
        #     Decimal(str(amount)) * Decimal(str(staking_proportion)) / Decimal(str(annualcycle)) / Decimal(
        #         str(verifier_num)))
        # # staking_reward = amount - block_reward
        result = node.ppos.getPackageReward()
        block_reward = result['Ret']
        result = node.ppos.getStakingReward()
        staking_reward = int(Decimal(str(result['Ret'])) / Decimal(str(verifier_num)))
        return block_reward, staking_reward

    def get_settlement_switchpoint(self, node: Node, number=0):
        """
        Get the last block of the current billing cycle
        :param node: node object
        :param number: number of billing cycles
        :return:
        """
        block_number = self.settlement_size * number
        tmp_current_block = node.eth.blockNumber
        current_end_block = math.ceil(tmp_current_block / self.settlement_size) * self.settlement_size + block_number
        return current_end_block

    def get_front_settlement_switchpoint(self, node: Node, number=0):
        """
        Get a block height before the current billing cycle
        :param node: node object
        :param number: number of billing cycles
        :return:
        """
        block_num = self.settlement_size * (number + 1)
        current_end_block = self.get_settlement_switchpoint(node)
        history_block = current_end_block - block_num + 1
        return history_block

    def wait_settlement_blocknum(self, node: Node, number=0):
        """
        Waiting for a billing cycle to settle
        :param node:
        :param number: number of billing cycles
        :return:
        """
        end_block = self.get_settlement_switchpoint(node, number)
        wait_block_number(node, end_block, self.interval)

    def get_annual_switchpoint(self, node: Node):
        """
        Get the number of annual settlement cycles
        """
        annual_cycle = (self.additional_cycle_time * 60) // self.settlement_size
        annualsize = annual_cycle * self.settlement_size
        current_block = node.eth.blockNumber
        current_end_block = math.ceil(current_block / annualsize) * annualsize
        return annual_cycle, annualsize, current_end_block

    def wait_annual_blocknum(self, node: Node):
        """
        Waiting for the end of the annual block high
        """
        annualcycle, annualsize, current_end_block = self.get_annual_switchpoint(node)
        current_block = node.eth.blockNumber
        differ_block = annualsize - (current_block % annualsize)
        annual_end_block = current_block + differ_block
        wait_block_number(node, annual_end_block, self.interval)

    def wait_consensus_blocknum(self, node: Node, number=0):
        """
        Waiting for a consensus round to end
        """
        end_block = self.get_consensus_switchpoint(node, number)
        wait_block_number(node, end_block, self.interval)

    def get_consensus_switchpoint(self, node: Node, number=0):
        """
        Get the specified consensus round high
        """
        block_number = self.consensus_size * number
        current_block = node.eth.blockNumber
        current_end_block = math.ceil(current_block / self.consensus_size) * self.consensus_size + block_number
        return current_end_block

    def get_report_reward(self, amount, penalty_ratio=None, proportion_ratio=None):
        """
        Gain income from double sign whistleblower and incentive pool
        :param node:
        :return:
        """
        if penalty_ratio is None:
            penalty_ratio = self.genesis.economicModel.slashing.slashFractionDuplicateSign
        if proportion_ratio is None:
            proportion_ratio = self.genesis.economicModel.slashing.duplicateSignReportReward
        penalty_reward = int(Decimal(str(amount)) * Decimal(str(penalty_ratio / 10000)))
        proportion_reward = int(Decimal(str(penalty_reward)) * Decimal(str(proportion_ratio / 100)))
        incentive_pool_reward = penalty_reward - proportion_reward
        return proportion_reward, incentive_pool_reward


if __name__ == '__main__':
    a = Economic()
    a.get_current_year_reward()
