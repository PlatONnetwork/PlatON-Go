from common.log import log
from tests.lib import EconomicConfig
from tests.lib.utils import assert_code, get_pledge_list
from common.key import mock_duplicate_sign
from tests.lib.utils import wait_block_number, get_the_dynamic_parameter_gas_fee, get_getDelegateReward_gas_fee
import rlp
from typing import List
import time, math, pytest
from tests.govern.test_voting_statistics import createstaking
from copy import copy
from conf import settings as conf

delegate_amount = 10 ** 18 * 1000
init_amount = 10 ** 18 * 10000


def get_new_value(value):
    if value == 10000:
        return value - 1
    else:
        return value + 1


def assert_cumulativeIncome_delegateepoch(client, address, delegateepoch=None):
    stakingnum = client.staking.get_stakingblocknum()
    delegateinfo = client.ppos.getDelegateInfo(stakingnum, address, client.node.node_id)
    log.info('Address {} delegate information {}'.format(address, delegateinfo))
    assert delegateinfo.get('Ret').get('CumulativeIncome') == 0
    if delegateepoch:
        assert delegateinfo.get('Ret').get('DelegateEpoch') == delegateepoch


def staking_and_delegate(clients, address, amount=delegate_amount):
    createstaking(clients, reward_per=1000)
    if isinstance(clients, List):
        clients = clients[0]
    result = clients.delegate.delegate(0, address, amount=amount)
    assert_code(result, 0)
    result = clients.ppos.getCandidateInfo(clients.node.node_id)
    assert result.get('Ret').get('DelegateRewardTotal') == 0

class TestCreateStaking:
    def assert_rewardsper_staking(self, client, nextrewardsper):
        assert_code(client.staking.get_rewardper(client.node), nextrewardsper)
        assert_code(client.staking.get_rewardper(client.node, isnext=True), nextrewardsper)
        value, nextvalue = get_pledge_list(client.ppos.getCandidateList, client.node.node_id)
        assert_code(value, nextrewardsper)
        assert_code(nextvalue, nextrewardsper)

    @pytest.mark.P0
    def test_IV_032_IV_033_IV_037(self, client_new_node):
        staking = client_new_node.staking
        address, _ = staking.economic.account.generate_account(staking.node.web3,
                                                            3 * staking.economic.genesis.economicModel.staking.stakeThreshold)
        result = staking.create_staking(0, address, address, node_id=staking.node.node_id,
                                        amount=2 * staking.economic.genesis.economicModel.staking.stakeThreshold,
                                        reward_per=1)
        assert_code(result, 0)
        assert_code(staking.get_rewardper(staking.node), 1)
        assert_code(staking.get_rewardper(staking.node, isnext=True), 1)
        result = staking.create_staking(0, address, address, node_id=staking.node.node_id,
                                        amount=2 * staking.economic.genesis.economicModel.staking.stakeThreshold,
                                        reward_per=1)
        log.info('Repeat create staking result : {}'.format(result))
        assert_code(result, 301101)
        self.assert_rewardsper_staking(client_new_node, 1)


    def test_IV_043(self, client_new_node):
        staking = client_new_node.staking
        address, _ = staking.economic.account.generate_account(staking.node.web3, 10 ** 18 * 30000000)
        plan = [{'Epoch': 20, 'Amount': 10 ** 18 * 2000000}]
        result = client_new_node.restricting.createRestrictingPlan(address, plan, address)
        log.info('CreateRestrictingPlan result : {}'.format(result))
        assert_code(result, 0)
        result = staking.create_staking(1, address, address, node_id=staking.node.node_id,
                                        amount=2 * staking.economic.genesis.economicModel.staking.stakeThreshold,
                                        reward_per=1155)
        assert_code(result, 0)
        self.assert_rewardsper_staking(client_new_node, 1155)

    def test_IV_036(self, client_new_node):
        staking = client_new_node.staking
        address, _ = staking.economic.account.generate_account(staking.node.web3,
                                                               3 * staking.economic.genesis.economicModel.staking.stakeThreshold)
        result = staking.create_staking(0, address, address, node_id=staking.node.node_id,
                                        amount=2 * staking.economic.genesis.economicModel.staking.stakeThreshold,
                                        reward_per=0)
        assert_code(result, 0)
        self.assert_rewardsper_staking(client_new_node, 0)

    def test_IV_038_IV_039_IV_044_IV_040_IV_041(self, client_new_node):
        staking = client_new_node.staking
        address, _ = staking.economic.account.generate_account(staking.node.web3,
                                                               3 * staking.economic.genesis.economicModel.staking.stakeThreshold)
        try:
            staking.create_staking(0, address, address, node_id=staking.node.node_id,
                                            amount=2 * staking.economic.genesis.economicModel.staking.stakeThreshold,
                                            reward_per=100000)
        except ValueError as e:
            assert e.args[0].get('message') == "gas required exceeds allowance or always failing transaction"

        try:
            staking.create_staking(0, address, address, node_id=staking.node.node_id,
                                        amount=2 * staking.economic.genesis.economicModel.staking.stakeThreshold,
                                        reward_per=-1)
        except TypeError as e:
            assert str(e) == "Did not find sedes handling type int"
        try:
            staking.create_staking(0, address, address, node_id=staking.node.node_id,
                                        amount=2 * staking.economic.genesis.economicModel.staking.stakeThreshold,
                                        reward_per=1.1)
        except TypeError as e:
            assert str(e) == "Did not find sedes handling type float"

        result = staking.create_staking(0, address, address, node_id=staking.node.node_id,
                                        amount=2 * staking.economic.genesis.economicModel.staking.stakeThreshold,
                                        reward_per=10001)
        assert_code(result, 301007)
        try:

            staking.create_staking(0, address, address, node_id=staking.node.node_id,
                                        amount=2 * staking.economic.genesis.economicModel.staking.stakeThreshold,
                                        reward_per='1000')
        except ValueError as e:
            assert e.args[0].get('message') == "gas required exceeds allowance or always failing transaction"

        result = staking.create_staking(0, address, address, node_id=staking.node.node_id,
                                        amount=2 * staking.economic.genesis.economicModel.staking.stakeThreshold,
                                        reward_per=10000)
        assert_code(result, 0)
        self.assert_rewardsper_staking(client_new_node, 10000)

    @pytest.mark.P2
    def test_IV_034(self, new_genesis_env, client_verifier):
        client = client_verifier
        value = client.staking.get_rewardper(client.node)
        value = get_new_value(value)
        result = client.staking.withdrew_staking(client.node.staking_address)
        assert_code(result, 0)
        client.economic.wait_settlement_blocknum(client.node, 2)
        address, _ = client.economic.account.generate_account(client.node.web3, 10**18 * 10000000)
        result = client.staking.create_staking(0, address, address, reward_per=value)
        assert_code(result, 0)
        self.assert_rewardsper_staking(client, value)

    @pytest.mark.P2
    def test_IV_035(self, new_genesis_env, clients_consensus):
        new_genesis_env.deploy_all()
        client = clients_consensus[0]
        value = client.staking.get_rewardper(client.node)
        value = get_new_value(value)
        wait_block_number(client.node, 50)
        report_information = mock_duplicate_sign(1, clients_consensus[1].node.nodekey,
                                                 clients_consensus[1].node.blsprikey, 41)
        log.info("Report information: {}".format(report_information))
        address, _ = client.economic.account.generate_account(client.node.web3, 10 ** 18 * 1000)
        result = client.duplicatesign.reportDuplicateSign(1, report_information, address)
        log.info('Node duplicate block result : {}'.format(result))
        assert_code(result, 0)
        verifier_list = get_pledge_list(client.ppos.getVerifierList)
        assert clients_consensus[1].node.node_id not in verifier_list
        client.economic.wait_settlement_blocknum(client.node, 2)
        candidate_info = clients_consensus[1].ppos.getCandidateInfo(clients_consensus[1].node.node_id)
        log.info(candidate_info)
        address, _ = client.economic.account.generate_account(client.node.web3, 10**18 * 10000000)
        result = clients_consensus[1].staking.create_staking(0, address, address, reward_per=value)
        assert_code(result, 0)
        self.assert_rewardsper_staking(clients_consensus[1], value)


class TestEditCandidate:
    def assert_rewardsper(self, client, rewardsper, nextrewardsper):
        assert_code(client.staking.get_rewardper(client.node), rewardsper)
        assert_code(client.staking.get_rewardper(client.node, isnext=True), nextrewardsper)
        value, nextvalue = get_pledge_list(client.ppos.getCandidateList, client.node.node_id)
        assert_code(value, rewardsper)
        assert_code(nextvalue, nextrewardsper)
        value, nextvalue = get_pledge_list(client.ppos.getVerifierList, client.node.node_id)
        assert_code(value, rewardsper)
        assert_code(nextvalue, nextrewardsper)
        value, nextvalue = get_pledge_list(client.ppos.getValidatorList, client.node.node_id)
        assert_code(value, rewardsper)
        assert_code(nextvalue, nextrewardsper)

    @pytest.mark.P0
    def test_MPI_018_to_027(self, client_verifier):
        client = client_verifier
        value = client.staking.get_rewardper(client.node)
        newvalue = get_new_value(value)
        print(value, newvalue)
        # client.economic.wait_settlement_blocknum(client.node, 1)
        # result = client.staking.edit_candidate(client.node.staking_address, client.node.staking_address,
        #                                        reward_per=newvalue)
        # assert_code(result, 0)
        # self.assert_rewardsper(client, value, newvalue)
        # client.economic.wait_settlement_blocknum(client.node, 1)
        # result = client.staking.edit_candidate(client.node.staking_address, client.node.staking_address,
        #                                        reward_per=0)
        # print(result)
        # assert_code(result, 0)
        # self.assert_rewardsper(client, 1, 0)
        #
        # client.economic.wait_settlement_blocknum(client.node, 1)
        # result = client.staking.edit_candidate(client.node.staking_address, client.node.staking_address,
        #                                        reward_per=10000)
        # assert_code(result, 301009)
        #
        # client.economic.wait_settlement_blocknum(client.node, 1)
        # result = client.staking.edit_candidate(client.node.staking_address, client.node.staking_address,
        #                                        reward_per=10001)
        # log.info('Edit candidate information reward percent is 10001, result : {}'.format(result))
        # assert_code(result, 301007)
        #
        # try:
        #     client.staking.edit_candidate(client.node.staking_address, client.node.staking_address,
        #                                        reward_per=-1)
        # except TypeError as e:
        #     str(e) == 'Did not find sedes handling type int'
        #
        # result = client.staking.edit_candidate(client.node.staking_address, client.node.staking_address,
        #                                        reward_per='2')
        # log.info('Edit candidate information reward percent is string, result : {}'.format(result))
        # assert_code(result, 0)

        # client.economic.wait_settlement_blocknum(client.node, 1)
        # result = client.staking.edit_candidate(client.node.staking_address, client.node.staking_address,
        #                                    reward_per=10000000000000000000000000000000000000000000000)
        # assert_code(result, 301007)

    @pytest.mark.P2
    def test_MPI_028_MPI_030_MPI_032(self, client_new_node):
        client = client_new_node
        address, _ = client.economic.account.generate_account(client.node.web3, 10**18 * 10000000)
        result = client.staking.create_staking(0, address, address,
                                               amount=2*client.economic.genesis.economicModel.staking.stakeThreshold,
                                               reward_per=100)
        assert_code(result, 0)
        client.economic.wait_settlement_blocknum(client.node, 1)
        verifier_list = get_pledge_list(client.ppos.getVerifierList)
        assert client.node.node_id in verifier_list
        address_delegate, _ = client.economic.account.generate_account(client.node.web3, init_amount)
        result = client.delegate.delegate(0, address_delegate, amount=10**18 * 1000)
        assert_code(result, 0)
        result = client.staking.edit_candidate(address, address, reward_per=101)

        assert_code(result, 0)
        assert_code(100, client.staking.get_rewardper())
        assert_code(101, client.staking.get_rewardper(isnext=True))
        vaule, newvalue = get_pledge_list(client.ppos.getCandidateList, nodeid=client.node.node_id)
        assert_code(vaule, 100)
        assert_code(newvalue, 101)

        client.economic.wait_settlement_blocknum(client.node, 1)
        result = client.staking.edit_candidate(address, address, reward_per=99)

        assert_code(result, 0)
        assert_code(101, client.staking.get_rewardper())
        assert_code(99, client.staking.get_rewardper(isnext=True))
        vaule, newvalue = get_pledge_list(client.ppos.getCandidateList, nodeid=client.node.node_id)
        assert_code(vaule, 101)
        assert_code(newvalue, 99)

        client.economic.wait_settlement_blocknum(client.node, 1)
        result = client.staking.edit_candidate(address, address, reward_per=9999)
        assert_code(result, 301009)

    @pytest.mark.P2
    def test_MPI_033(self,new_genesis_env, clients_verifier):
        client0 = clients_verifier[0]
        client1 = clients_verifier[1]
        wait_block_number(client1.node, client0.node.block_number + 40)
        report_information = mock_duplicate_sign(1, client0.node.nodekey,
                                                 client0.node.blsprikey, client0.node.block_number)
        log.info("Report information: {}".format(report_information))
        address1, _ = client1.economic.account.generate_account(client1.node.web3, init_amount)
        result = client1.duplicatesign.reportDuplicateSign(1, report_information, address1)
        log.info('Node duplicate block result : {}'.format(result))
        assert_code(result, 0)
        value = client1.staking.get_rewardper(client1.node)
        newvalue = get_new_value(value)
        result = client0.staking.edit_candidate(client0.node.staking_address, client0.node.benifit_address,
                                                reward_per=newvalue)
        log.info('result {}'.format(result))
        assert_code(result, 301103)
        assert_code(value, client1.staking.get_rewardper())
        assert_code(value, client1.staking.get_rewardper(isnext=True))
        vaule_per, newvalue_per = get_pledge_list(client1.ppos.getCandidateList, nodeid=client1.node.node_id)
        assert_code(vaule_per, value)
        assert_code(newvalue_per, value)

    @pytest.mark.P2
    def test_MPI_034(self, client_new_node):
        client = client_new_node
        address, _ = client.economic.account.generate_account(client.node.web3, 10**18 * 10000000)
        result = client.staking.create_staking(0, address, address, reward_per=100)
        assert_code(result, 0)

        client.economic.wait_settlement_blocknum(client.node, 1)
        result = client.staking.edit_candidate(address, address, reward_per=101)
        assert_code(result, 0)
        assert_code(100, client.staking.get_rewardper())
        assert_code(101, client.staking.get_rewardper(isnext=True))
        vaule, newvalue = get_pledge_list(client.ppos.getCandidateList, nodeid=client.node.node_id)
        assert_code(vaule, 100)
        assert_code(newvalue, 101)

    @pytest.mark.P2
    def test_MPI_035(self, new_genesis_env, client_verifier):
        client = client_verifier
        value, nextvalue = get_pledge_list(client.ppos.getVerifierList, nodeid=client.node.node_id)
        newvalue = get_new_value(value)
        result = client.staking.withdrew_staking(client.node.staking_address)
        assert_code(result, 0)
        result = client.staking.edit_candidate(client.node.staking_address, client.node.benifit_address,
                                               reward_per=newvalue)
        log.info('Edit exiting candidate information result : {}'.format(result))
        assert_code(result, 301103)
        assert_code(client.staking.get_rewardper(client.node), value)
        assert_code(client.staking.get_rewardper(client.node, isnext=True), nextvalue)
        candidates = get_pledge_list(client.ppos.getCandidateList)
        verifiers = get_pledge_list(client.ppos.getVerifierList)
        validators = get_pledge_list(client.ppos.getVerifierList)
        assert client.node.node_id not in candidates
        assert client.node.node_id in verifiers
        assert client.node.node_id in validators
        value_exiting, nextvalue_exiting = get_pledge_list(client.ppos.getVerifierList, client.node.node_id)
        assert_code(value, value_exiting)
        assert_code(nextvalue, nextvalue_exiting)

    @pytest.mark.P2
    def test_MPI_036(self, client_verifier):
        client = client_verifier
        value, nextvalue = get_pledge_list(client.ppos.getVerifierList, client.node.node_id)
        newvalue = get_new_value(value)
        client.economic.wait_settlement_blocknum(client.node, 1)
        result = client.staking.edit_candidate(client.node.staking_address, client.node.staking_address,
                                               reward_per=newvalue)
        assert_code(result, 0)
        self.assert_rewardsper(client, value, newvalue)
        verifier_info = client.ppos.getVerifierList()
        log.info('nodeid {} verifier_info : {}'.format(client.node.node_id, verifier_info))
        result = client.staking.withdrew_staking(client.node.staking_address)
        assert_code(result, 0)
        assert_code(client.staking.get_rewardper(client.node), value)
        assert_code(client.staking.get_rewardper(client.node, isnext=True), newvalue)
        verifiers = get_pledge_list(client.ppos.getVerifierList)
        candidates = get_pledge_list(client.ppos.getCandidateList)
        validators = get_pledge_list(client.ppos.getVerifierList)
        assert client.node.node_id not in candidates
        assert client.node.node_id in verifiers
        assert client.node.node_id in validators
        verifier_info = client.ppos.getVerifierList()
        log.info('nodeid {} verifier_info : {}'.format(client.node.node_id, verifier_info))
        value_exiting, nextvalue_exiting = get_pledge_list(client.ppos.getVerifierList, client.node.node_id)
        assert_code(value, value_exiting)
        assert_code(newvalue, nextvalue_exiting)

class TestgetDelegateReward:
    @pytest.mark.P0
    def test_IN_GR_001_IN_GR_002_IN_GR_004_IN_GR_005_IN_GR_006_IN_DR_018_IN_DR_019(self, clients_new_node):
        client0 = clients_new_node[0]
        client1 = clients_new_node[1]
        delegate_address, _ = client0.economic.account.generate_account(client0.node.web3, init_amount)
        staking_and_delegate(client0, delegate_address)
        staking_and_delegate(client1, delegate_address)
        staking_num0 = client0.staking.get_stakingblocknum(client0.node)
        staking_num1 = client1.staking.get_stakingblocknum(client1.node)
        result = client1.ppos.getDelegateInfo(staking_num1, delegate_address, client1.node.node_id)
        log.info(result)
        assert_code(result, 0)
        result = client0.ppos.getDelegateInfo(staking_num0, delegate_address, client0.node.node_id)
        log.info(result)
        assert_code(result, 0)
        client0.economic.wait_settlement_blocknum(client0.node)
        verifier_list = get_pledge_list(client0.ppos.getVerifierList)
        assert client0.node.node_id in verifier_list
        assert client1.node.node_id in verifier_list

        nodeid = 'a128496e6d88a3c231ceda1c095b3cc51a9b25ea14b8d1c8839e5cce65888' \
                 'f68e41231e7ebec710e41bcd5c53a704da4ded75a1398762d4c01a165897c911111'
        result = client0.ppos.getDelegateReward(delegate_address, node_ids=[nodeid])
        assert_code(result, 305001)

        result = client0.ppos.getDelegateReward(delegate_address, node_ids=[nodeid, client0.node.node_id])
        assert_code(result, 0)
        assert len(result.get('Ret')) == 1
        assert result.get('Ret')[0].get('nodeID') == client0.node.node_id
        assert result.get('Ret')[0].get('stakingNum') == staking_num0

        result = client0.ppos.getDelegateReward(delegate_address)
        log.info('Do not given nodeid, get address delegate reward : {}'.format(result))
        assert_code(result, 0)
        assert len(result.get('Ret')) == 2
        for rewardinfo in result.get('Ret'):
            if client0.node.node_id == rewardinfo.get('nodeID'):
                rewardinfo.get('reward') == 0
                rewardinfo.get('stakingNum') == staking_num0
            else:
                rewardinfo.get('reward') == 0
                rewardinfo.get('stakingNum') == staking_num1
        result = client0.ppos.getDelegateReward(delegate_address, node_ids=[client0.node.node_id])
        log.info('Get address delegate nodeid {} reward : {}'.format(client0.node.node_id, result))
        assert_code(result, 0)
        assert result.get('Ret')[0].get('nodeID') == client0.node.node_id
        assert result.get('Ret')[0].get('stakingNum') == staking_num0

        result = client0.ppos.getDelegateReward(delegate_address, node_ids=[client1.node.node_id])
        log.info('Get address delegate nodeid {} reward : {}'.format(client1.node.node_id, result))
        assert_code(result, 0)
        assert result.get('Ret')[0].get('nodeID') == client1.node.node_id
        assert result.get('Ret')[0].get('stakingNum') == staking_num1

        result = client0.ppos.getDelegateReward(delegate_address, node_ids=[client0.node.node_id, client1.node.node_id])
        log.info('Get address delegate nodeid {},{} reward : {}'.format(client1.node.node_id,
                                                                        client1.node.node_id, result))
        assert_code(result, 0)
        assert len(result.get('Ret')) == 2
        for rewardinfo in result.get('Ret'):
            if client0.node.node_id == rewardinfo.get('nodeID'):
                rewardinfo.get('reward') == 0
                rewardinfo.get('stakingNum') == staking_num0
            else:
                rewardinfo.get('reward') == 0
                rewardinfo.get('stakingNum') == staking_num1
        client0.economic.wait_settlement_blocknum(client0.node)

        reward = client0.delegate.get_delegate_reward_by_nodeid(delegate_address)
        log.info('Address {} reward : {}'.format(delegate_address, reward))
        assert reward > 0
        balance_before = client0.node.eth.getBalance(delegate_address)
        log.info('Before delegate address balance : {}'.format(balance_before))

        gas = get_getDelegateReward_gas_fee(client0, 2, 2)

        log.info('Calculated gas : {}'.format(gas))
        result = client0.delegate.withdraw_delegate_reward(delegate_address)
        assert_code(result, 0)
        balance_after = client0.node.eth.getBalance(delegate_address)
        log.info('After withdraw reward, address balance : {}'.format(balance_after))
        assert balance_before + reward - gas == balance_after
        result = client0.delegate.withdraw_delegate_reward(delegate_address)
        assert_code(result, 0)
        assert_cumulativeIncome_delegateepoch(client0, delegate_address, delegateepoch=math.ceil(
            client0.node.block_number/client0.economic.settlement_size))
        assert_cumulativeIncome_delegateepoch(client1, delegate_address, delegateepoch=math.ceil(
            client0.node.block_number/client0.economic.settlement_size))
        balance_after_1 = client0.node.eth.getBalance(delegate_address)
        log.info('After repeat withdraw reward address balance : {}'.format(balance_after_1))
        gas = get_getDelegateReward_gas_fee(client0, 2, 0)
        assert balance_after - gas == balance_after_1
        client0.economic.wait_settlement_blocknum(client0.node)
        reward = client0.delegate.get_delegate_reward_by_nodeid(delegate_address)
        assert reward > 0
        result = client0.delegate.withdraw_delegate_reward(delegate_address)
        assert_code(result, 0)
        balance_after_2 = client0.node.eth.getBalance(delegate_address)
        log.info('After {} withdraw reward address balance : {}'.format(delegate_address, balance_after_2))
        assert_cumulativeIncome_delegateepoch(client0, delegate_address, delegateepoch=math.ceil(
            client0.node.block_number/client0.economic.settlement_size))
        assert_cumulativeIncome_delegateepoch(client1, delegate_address, delegateepoch=math.ceil(
            client0.node.block_number/client0.economic.settlement_size))
        gas = get_getDelegateReward_gas_fee(client1, 2, 2)
        assert balance_after_1 - gas + reward == balance_after_2

    @pytest.mark.P2
    def test_IN_GR_008(self, client_verifier):
        client = client_verifier
        address, _ = client.economic.account.generate_account(client.node.web3, init_amount)
        try:
            client.ppos.getDelegateReward(address, node_ids=client.node.node_id)
        except ValueError as e:
            str(e) == 'non-hexadecimal number found in fromhex() arg at position 1'

    @pytest.mark.P2
    def test_IN_GR_010_IN_GR_012(self, client_new_node):
        client = client_new_node
        address, _ = client.economic.account.generate_account(client.node.web3, init_amount)
        epoch = math.ceil(client.node.block_number/client.economic.settlement_size)
        staking_and_delegate(client, address)
        staking_blocknum = client.staking.get_stakingblocknum()
        client.economic.wait_settlement_blocknum(client.node)
        verifier_list = get_pledge_list(client.ppos.getVerifierList)
        assert client.node.node_id in verifier_list
        rewardinfo = client.ppos.getDelegateReward(address)
        log.info('Get address {} reward information : {}'.format(address, rewardinfo))
        assert len(rewardinfo.get('Ret')) == 1
        assert rewardinfo.get('Ret')[0].get('stakingNum') == staking_blocknum
        assert rewardinfo.get('Ret')[0].get('reward') == 0
        client.economic.wait_settlement_blocknum(client.node)
        delegate_info = client.ppos.getDelegateInfo(staking_blocknum, address, client.node.node_id)
        log.info('Address {} delegate information : {}'.format(address, delegate_info))
        assert_code(delegate_info, 0)
        assert delegate_info.get('Ret').get('DelegateEpoch') == epoch
        rewardinfo = client.ppos.getDelegateReward(address)
        log.info('Get address {} reward information : {}'.format(address, rewardinfo))
        reward = rewardinfo.get('Ret')[0].get('reward')
        assert reward > 0
        client.staking.withdrew_staking(client.node.staking_address)
        reward1 = client.delegate.get_delegate_reward_by_nodeid(address)
        assert reward == reward1

        before_withdraw_balance = client.node.eth.getBalance(address)
        log.info('Address {} before withdraw {}'.format(address, before_withdraw_balance))
        result = client.delegate.withdraw_delegate_reward(address)
        assert_code(result, 0)
        after_withdraw_balance = client.node.eth.getBalance(address)
        log.info('Address {} before withdraw {}'.format(address, after_withdraw_balance))
        gas = get_getDelegateReward_gas_fee(client, 1, 1)
        assert before_withdraw_balance - gas + reward1 == after_withdraw_balance

    @pytest.mark.P2
    def test_IN_GR_011_IN_GR_012(self, client_new_node, client_verifier):
        client = client_new_node
        address, _ = client.economic.account.generate_account(client.node.web3, init_amount)
        staking_and_delegate(client_new_node, address)
        staking_blocknum = client.staking.get_stakingblocknum()
        client.economic.wait_settlement_blocknum(client.node)
        verifier_list = get_pledge_list(client.ppos.getVerifierList)
        assert client.node.node_id in verifier_list
        rewardinfo = client.ppos.getDelegateReward(address)
        log.info('Get address {} reward information : {}'.format(address, rewardinfo))
        assert rewardinfo.get('Ret')[0].get('stakingNum') == staking_blocknum
        client.economic.wait_settlement_blocknum(client.node)
        staking_info = client.ppos.getCandidateInfo(client.node.node_id)
        log.info('Node {} candidate information : {}'.format(client.node.node_id, staking_info))
        delegate_info = client.ppos.getDelegateInfo(staking_blocknum, address, client.node.node_id)
        log.info('Address {} delegate information : {}'.format(address, delegate_info))
        rewardinfo = client.ppos.getDelegateReward(address)
        log.info('Get address {} reward information : {}'.format(address, rewardinfo))
        reward = rewardinfo.get('Ret')[0].get('reward')
        assert reward > 0
        report_information = mock_duplicate_sign(1, client.node.nodekey,
                                                 client.node.blsprikey, client.node.block_number - 10)
        log.info("Report information: {}".format(report_information))
        address1, _ = client_verifier.economic.account.generate_account(client_verifier.node.web3, init_amount)
        result = client_verifier.duplicatesign.reportDuplicateSign(1, report_information, address1)
        log.info('Node duplicate block result : {}'.format(result))
        assert_code(result, 0)
        verifier_list = get_pledge_list(client_verifier.ppos.getVerifierList)
        assert client.node.node_id not in verifier_list
        rewardinfo = client.ppos.getDelegateReward(address)
        assert reward == rewardinfo.get('Ret')[0].get('reward')
        client.economic.wait_settlement_blocknum(client.node, client.economic.genesis.economicModel.staking.unStakeFreezeDuration)

        rewardinfo = client.ppos.getDelegateReward(address)
        assert reward == rewardinfo.get('Ret')[0].get('reward')
        before_balance = client.node.eth.getBalance(address)
        log.info('Before withdraw address {} balance {}'.format(address, before_balance))
        reward = client.delegate.get_delegate_reward_by_nodeid(address)
        log.info('Address {} reward {}'.format(address, reward))
        gas = get_getDelegateReward_gas_fee(client, 1, 1)
        result = client.delegate.withdraw_delegate_reward(address)
        assert_code(result, 0)
        after_balance = client.node.eth.getBalance(address)
        log.info('After withdraw address {} balance {}'.format(address, after_balance))
        assert before_balance + reward - gas == after_balance

    @pytest.mark.P2
    def test_IN_GR_013(self, client_new_node):
        client = client_new_node
        address, _ = client.economic.account.generate_account(client.node.web3, init_amount)
        result = client.ppos.getDelegateReward(address, [client.node.node_id])
        assert_code(result, 305001)

    @pytest.mark.P2
    def test_IN_GR_014(self, client_new_node, client_consensus):
        client = client_new_node
        address, _ = client.economic.account.generate_account(client.node.web3, init_amount)
        result = client.ppos.getDelegateReward(address, [client_consensus.node.node_id])
        assert_code(result, 305001)

    @pytest.mark.P2
    def test_IN_GR_015_to_IN_GR_018(self, clients_new_node):
        client1 = clients_new_node[0]
        client2 = clients_new_node[1]
        address1, _ = client1.economic.account.generate_account(client1.node.web3, init_amount)
        address2, _ = client1.economic.account.generate_account(client1.node.web3, init_amount)
        address3, _ = client1.economic.account.generate_account(client1.node.web3, init_amount)
        staking_and_delegate(client1, address1)
        staking_and_delegate(client2, address2)
        staking_num1 = client1.staking.get_stakingblocknum()
        log.info('Node {} staking block num : {}'.format(client1.node.node_id, staking_num1))
        staking_num2 = client2.staking.get_stakingblocknum()
        log.info('Node {} staking block num : {}'.format(client2.node.node_id, staking_num2))
        client2.delegate.delegate(0, address1, amount=10**18 * 1000)
        client1.economic.wait_settlement_blocknum(client1.node, number=1)
        reward_info_address1_1 = client1.ppos.getDelegateReward(address1)
        log.info('Not incoming nodeid, Address {} reward information : {}'.format(address1, reward_info_address1_1))
        assert_code(reward_info_address1_1, 0)
        assert_code(len(reward_info_address1_1.get('Ret')), 2)
        assert client1.delegate.get_staking_num_by_nodeid(address1, client1.node.node_id) == staking_num1
        assert client2.delegate.get_staking_num_by_nodeid(address1, client2.node.node_id) == staking_num2
        assert reward_info_address1_1.get('Ret')[0].get('reward') > 0
        assert reward_info_address1_1.get('Ret')[0].get('reward') != reward_info_address1_1.get('Ret')[1].get('reward')

        reward_info_address1_2 = client1.ppos.getDelegateReward(address1, node_ids=[client1.node.node_id])
        log.info('incoming nodeid {}, Address {} reward information : {}'.format(client1.node.node_id,
                                                                              address1, reward_info_address1_2))
        assert_code(reward_info_address1_2, 0)
        assert_code(len(reward_info_address1_2.get('Ret')), 1)
        assert client2.delegate.get_staking_num_by_nodeid(address1, client1.node.node_id) == staking_num1
        assert reward_info_address1_2.get('Ret')[0].get('reward') != 0

        reward_info_address1_3 = client1.ppos.getDelegateReward(address1, node_ids=[client2.node.node_id])
        log.info('incoming nodeid {}, Address {} reward information : {}'.format(client2.node.node_id,
                                                                              address1, reward_info_address1_3))
        assert_code(reward_info_address1_3, 0)
        assert_code(len(reward_info_address1_3.get('Ret')), 1)
        assert client2.delegate.get_staking_num_by_nodeid(address1, client2.node.node_id) == staking_num2
        assert reward_info_address1_3.get('Ret')[0].get('reward') != 0
        reward_info_address2_1 = client2.ppos.getDelegateReward(address2, node_ids=[client2.node.node_id])
        log.info('incoming nodeid {}, address {}, reward information : {}'.format(
            client2.node.node_id, address2, reward_info_address2_1
        ))
        assert_code(reward_info_address2_1, 0)
        assert_code(len(reward_info_address2_1.get('Ret')), 1)
        assert client2.delegate.get_staking_num_by_nodeid(address2, client2.node.node_id) == staking_num2
        assert reward_info_address2_1.get('Ret')[0].get('reward') != 0
        reward_info_address2_2 = client2.ppos.getDelegateReward(address2)
        log.info('Not incoming nodeid {}, address {}, reward information : {}'.format(
            client2.node.node_id, address2, reward_info_address2_2
        ))
        assert_code(reward_info_address2_2, 0)
        assert_code(len(reward_info_address2_2.get('Ret')), 1)
        assert client2.delegate.get_staking_num_by_nodeid(address2, client2.node.node_id) == staking_num2
        assert reward_info_address2_2.get('Ret')[0].get('reward') != 0
        reward_info_address2_3 = client2.ppos.getDelegateReward(address2, [client1.node.node_id])
        log.info('Not incoming nodeid {}, address {}, reward information : {}'.format(
            client2.node.node_id, address2, reward_info_address2_3
        ))
        assert_code(reward_info_address2_3, 305001)
        reward_info_address2_4 = client2.ppos.getDelegateReward(address2, [client1.node.node_id, client2.node.node_id])
        log.info('Incoming nodeid {}，{}, address {}, reward information : {}'.format(client1.node.node_id,
            client2.node.node_id, address2, reward_info_address2_4
        ))
        assert_code(reward_info_address2_2, 0)
        assert_code(len(reward_info_address2_2.get('Ret')), 1)
        assert client2.delegate.get_staking_num_by_nodeid(address2, client2.node.node_id) == staking_num2
        assert reward_info_address2_2.get('Ret')[0].get('reward') != 0

        reward_info_address3_1 = client1.ppos.getDelegateReward(address3)
        assert_code(reward_info_address3_1, 305001)
        reward_info_address3_2 = client1.ppos.getDelegateReward(address3, node_ids=[client1.node.node_id])
        assert_code(reward_info_address3_2, 305001)
        reward_info_address3_3 = client1.ppos.getDelegateReward(address3, node_ids=[client2.node.node_id])
        assert_code(reward_info_address3_3, 305001)
        address1_balance_before = client1.node.eth.getBalance(address1)
        log.info('Before getDelegateReward, the address {} balance: {}'.format(address1, address1_balance_before))
        reward_address1 = client1.delegate.get_delegate_reward_by_nodeid(address1, node_ids=[client1.node.node_id])
        log.info('Address {} delegate node {} reward : {}'.format(address1, client1.node.node_id, reward_address1))
        result = client1.delegate.withdrew_delegate(staking_num1, address1, node_id=client1.node.node_id,
                                                    amount=10**18 * 1000, transaction_cfg=client1.pip.cfg.transaction_cfg)
        assert_code(result, 0)
        reward_info_address1 = client1.ppos.getDelegateReward(address1, node_ids=[client1.node.node_id])
        assert_code(reward_info_address1, 305001)
        address1_balance_after = client1.node.eth.getBalance(address1)
        log.info('Before getDelegateReward, the address {} balance: {}'.format(address1, address1_balance_after))
        data = rlp.encode([rlp.encode(int(1005)), rlp.encode(staking_num1), rlp.encode(bytes.fromhex(client1.node.node_id)),
                           rlp.encode(10**18 * 1000)])
        connt = get_the_dynamic_parameter_gas_fee(data)
        gas = (21000 + 6000 + 8000 + 100 + connt) * client1.pip.cfg.transaction_cfg.get('gasPrice')
        assert address1_balance_before + reward_address1 - gas + 10**18 * 1000 == address1_balance_after

    @pytest.mark.P2
    def test_IN_GR_019(self, clients_new_node):
        client1 = clients_new_node[0]
        client2 = clients_new_node[1]
        address, _ = client1.economic.account.generate_account(client1.node.web3, init_amount)
        staking_and_delegate(client1, address)
        stakingnum = client1.staking.get_stakingblocknum()
        result = client1.ppos.getDelegateReward(address, [client1.node.node_id, client2.node.node_id])
        assert len(result.get('Ret')) == 1
        assert result.get('Ret')[0].get('nodeID') == client1.node.node_id
        assert result.get('Ret')[0].get('reward') == 0
        assert result.get('Ret')[0].get('stakingNum') == stakingnum
        nodeid = 'ad5ef1acceadfb7f2c434e4451dbb8dd02f866d4344495f01908093d40cbe5a61a03' \
                 '93eb09ec53abd7d270111f43f4a74c40e8c8680d3c47e1f190b5af991111'
        result = client1.ppos.getDelegateReward(address, [client1.node.node_id, nodeid])
        assert len(result.get('Ret')) == 1
        assert result.get('Ret')[0].get('nodeID') == client1.node.node_id
        assert result.get('Ret')[0].get('reward') == 0
        assert result.get('Ret')[0].get('stakingNum') == stakingnum

        client1.economic.wait_settlement_blocknum(client1.node)
        wait_block_number(client1.node, client1.node.block_number + 80)
        verifier_list = get_pledge_list(client1.ppos.getVerifierList)
        assert client1.node.node_id in verifier_list
        balance_before_withdraw = client1.node.eth.getBalance(address)
        log.info('Address {} before withdraw balance : {}'.format(address, balance_before_withdraw))
        result = client1.delegate.withdrew_delegate(stakingnum, address, amount=delegate_amount)
        assert_code(result, 0)
        data = rlp.encode([rlp.encode(int(1005)), rlp.encode(stakingnum),
                           rlp.encode(bytes.fromhex(client1.node.node_id)), rlp.encode(delegate_amount)])
        gas = (21000 + 6000 + 8000 + get_the_dynamic_parameter_gas_fee(data)) * client1.node.eth.gasPrice
        balance_after_withdraw = client1.node.eth.getBalance(address)
        log.info('Address {} after withdraw balance : {}'.format(address, balance_after_withdraw))
        assert balance_before_withdraw - gas + delegate_amount == balance_after_withdraw

class TestwithdrawDelegateReward():
    @pytest.mark.P0
    def test_IN_GR_020_IN_GR_021_IN_GR_022_IN_GR_023(self, client_new_node):
        client = client_new_node
        address, _ = client.economic.account.generate_account(client.node.web3, init_amount)
        staking_and_delegate(client, address)
        assert client.staking.get_rewardper() == 1000
        result = client.ppos.getDelegateReward(address)
        log.info('Address {} delegate reward information : {}'.format(address, result))
        assert_code(result, 0)
        balance_before_withdraw = client.node.eth.getBalance(address)
        log.info('Address {} before withdraw reward balance : {}'.format(address, balance_before_withdraw))
        result = client.delegate.withdraw_delegate_reward(address)
        log.info('Address {} withdraw delegate reward result : {}'.format(address, result))
        assert_code(result, 0)
        balance_after_withdraw = client.node.eth.getBalance(address)
        log.info('Address {} after withdraw reward balance : {}'.format(address, balance_after_withdraw))
        gas = get_getDelegateReward_gas_fee(client, 1, 0)
        assert balance_before_withdraw - gas == balance_after_withdraw
        client.economic.wait_settlement_blocknum(client.node)
        blocknum = client.node.block_number
        log.info('Block number : {}'.format(blocknum))
        verifier_list = get_pledge_list(client.ppos.getVerifierList)
        assert client.node.node_id in verifier_list

        result = client.ppos.getDelegateReward(address)
        log.info('Address {} delegate reward information : {}'.format(address, result))
        assert_code(result, 0)
        assert result.get('Ret')[0].get('reward') == 0
        result = client.delegate.withdraw_delegate_reward(address)
        log.info('Address {} withdraw delegate reward result : {}'.format(address, result))
        assert_code(result, 0)
        balance_after_withdraw2 = client.node.eth.getBalance(address)
        log.info('Address {} after withdraw reward balance : {}'.format(address, balance_after_withdraw2))
        gas = get_getDelegateReward_gas_fee(client, 1, 0)
        assert balance_after_withdraw - gas == balance_after_withdraw2
        wait_block_number(client.node, math.ceil(blocknum/client.economic.settlement_size
                                                 ) * client.economic.settlement_size - 10)
        result = client.ppos.getDelegateReward(address)
        log.info('Address {} delegate reward information : {}'.format(address, result))
        assert_code(result, 0)
        assert result.get('Ret')[0].get('reward') == 0
        wait_block_number(client.node, math.ceil(blocknum/client.economic.settlement_size
                                                 ) * client.economic.settlement_size)
        log.info('blocknumber : {}'.format(client.node.block_number))
        result = client.ppos.getDelegateReward(address)
        log.info('Address {} delegate reward information : {}'.format(address, result))
        assert_code(result, 0)
        delegate_info = client.ppos.getDelegateInfo(client.staking.get_stakingblocknum(), address,
                                                       client.node.node_id)
        log.info('Address delegate information : {}'.format(delegate_info))
        candidate_info = client.ppos.getCandidateInfo(client.node.node_id)
        log.info('nodeid {} candidate information : {}'.format(client.node.node_id, candidate_info))
        reward = client.delegate.get_delegate_reward_by_nodeid(address, [client.node.node_id])
        log.info('Address {} reward : {}'.format(address, reward))
        assert reward != 0
        result = client.delegate.withdraw_delegate_reward(address)
        log.info('Address {} withdraw delegate reward result : {}'.format(address, result))
        assert_code(result, 0)
        balance_after_withdraw3 = client.node.eth.getBalance(address)
        log.info('Address {} after withdraw reward balance : {}'.format(address, balance_after_withdraw3))
        gas = get_getDelegateReward_gas_fee(client, 1, 1)
        assert balance_after_withdraw2 + reward - gas == balance_after_withdraw3

    @pytest.mark.P2
    def test_IN_GR_024_IN_DR_002_IN_DR_003(self, clients_new_node):
        client0 = clients_new_node[0]
        client1 = clients_new_node[1]
        address0, _ = client0.economic.account.generate_account(client0.node.web3, init_amount)
        address1, _ = client0.economic.account.generate_account(client0.node.web3, init_amount)
        address2, _ = client0.economic.account.generate_account(client0.node.web3, init_amount)
        staking_and_delegate(client0, address0)
        staking_and_delegate(client1, address1)

        stakingnum0 = client0.staking.get_stakingblocknum(client0.node)
        stakingnum1 = client1.staking.get_stakingblocknum(client1.node)

        client0.economic.wait_settlement_blocknum(client0.node)
        result = client1.delegate.delegate(0, address0, amount=delegate_amount)
        assert_code(result, 0)
        client0.economic.wait_settlement_blocknum(client0.node, 1)
        balance_before_withdraw0 = client0.node.eth.getBalance(address0)
        log.info('Address {} before withdraw balance : {}'.format(address0, balance_before_withdraw0))
        balance_before_withdraw1 = client0.node.eth.getBalance(address1)
        log.info('Address {} before withdraw balance : {}'.format(address0, balance_before_withdraw1))
        rewards0 = client0.delegate.get_delegate_reward_by_nodeid(address0)
        log.info('Address rewards : {}'.format(rewards0))
        rewards1 = client1.delegate.get_delegate_reward_by_nodeid(address1)
        log.info('Address rewards : {}'.format(rewards1))
        result = client0.ppos.getDelegateReward(address2)
        assert_code(result, 305001)
        delegate_info0 = client0.ppos.getDelegateInfo(stakingnum0, address0, client0.node.node_id)
        log.info('Address delegate nodeid {} information : {}'.format(client0.node.node_id, delegate_info0))
        delegate_info1 = client1.ppos.getDelegateInfo(stakingnum1, address0, client1.node.node_id)
        log.info('Address delegate nodeid {} information : {}'.format(client1.node.node_id, delegate_info1))
        gas0 = get_getDelegateReward_gas_fee(client0, 2, 3)
        gas1 = get_getDelegateReward_gas_fee(client1, 1, 2)
        result = client0.delegate.withdraw_delegate_reward(address0)
        assert_code(result, 0)
        result = client1.delegate.withdraw_delegate_reward(address1)
        assert_code(result, 0)
        epoch0 = math.ceil(client0.node.block_number/client0.economic.settlement_size)
        assert_cumulativeIncome_delegateepoch(client0, address0, delegateepoch=epoch0)
        assert_cumulativeIncome_delegateepoch(client1, address0, delegateepoch=epoch0)
        balance_after_withdraw0 = client0.node.eth.getBalance(address0)
        log.info('After withdraw address {} balance : {}'.format(address0, balance_after_withdraw0))
        balance_after_withdraw1 = client0.node.eth.getBalance(address1)
        log.info('After withdraw address {} balance : {}'.format(address1, balance_after_withdraw1))
        assert balance_before_withdraw0 - gas0 + rewards0 == balance_after_withdraw0
        assert balance_before_withdraw1 - gas1 + rewards1 == balance_after_withdraw1

    @pytest.mark.P2
    def test_IN_DR_004_IN_DR_014(self, clients_new_node, reset_environment):
        client0 = clients_new_node[0]
        client1 = clients_new_node[1]
        address0, _ = client0.economic.account.generate_account(client0.node.web3, init_amount)
        address1, _ = client0.economic.account.generate_account(client0.node.web3, init_amount)
        staking_and_delegate(client0, address0)
        staking_and_delegate(client1, address1)
        staking_num0 = client0.staking.get_stakingblocknum()
        staking_num1 = client1.staking.get_stakingblocknum()
        result = client1.delegate.delegate(0, address1, node_id=client0.node.node_id, amount=delegate_amount)
        assert_code(result, 0)
        epoch = math.ceil(client0.node.block_number/client0.economic.settlement_size)
        client0.economic.wait_settlement_blocknum(client0.node, 1)
        rewards0 = client0.delegate.get_delegate_reward_by_nodeid(address0)
        log.info('Address {} delegate rewards : {}'.format(address0, rewards0))
        rewards1 = client1.delegate.get_delegate_reward_by_nodeid(address1)
        log.info('Address {} delegate rewards : {}'.format(address1, rewards1))
        rewards1_node0 = client1.delegate.get_delegate_reward_by_nodeid(address1, node_ids=[client0.node.node_id])
        log.info('Address {} delegate nodeid {} rewards : {}'.format(address1, client0.node.node_id, rewards1_node0))

        balance_before0 = client0.node.eth.getBalance(address0)
        log.info('Address {} before withdraw delegate balance : {}'.format(address0, balance_before0))
        balance_before1 = client1.node.eth.getBalance(address1)
        log.info('Address {} before withdraw delegate balance : {}'.format(address0, balance_before1))
        result = client0.delegate.withdrew_delegate(staking_num0, address0, amount=delegate_amount)
        assert_code(result, 0)
        result = client0.delegate.withdrew_delegate(staking_num0, address1, amount=delegate_amount)
        assert_code(result, 0)

        assert_cumulativeIncome_delegateepoch(client1, address1, delegateepoch=epoch)
        data0 = rlp.encode([rlp.encode(int(1005)), rlp.encode(staking_num0), rlp.encode(bytes.fromhex(
            client0.node.node_id)), rlp.encode(delegate_amount)])
        log.info('data {}'.format(data0))
        data1 = rlp.encode([rlp.encode(int(1005)), rlp.encode(staking_num1), rlp.encode(bytes.fromhex(
            client1.node.node_id)), rlp.encode(delegate_amount)])
        gas0 = (21000 + 6000 + 8000 + 100 + get_the_dynamic_parameter_gas_fee(data0)) * client0.node.eth.gasPrice
        gas1 = (21000 + 6000 + 8000 + 100 + get_the_dynamic_parameter_gas_fee(data1)) * client0.node.eth.gasPrice
        print(gas0)
        print(get_the_dynamic_parameter_gas_fee(data0))
        balance_after0 = client0.node.eth.getBalance(address0)
        log.info('Address {} before withdraw delegate balance : {}'.format(address0, balance_after0))
        balance_after1 = client0.node.eth.getBalance(address1)
        log.info('Address {} before withdraw delegate balance : {}'.format(address0, balance_after1))
        assert balance_before0 - gas0 + rewards0 + delegate_amount == balance_after0
        assert balance_before1 - gas1 + rewards1_node0 + delegate_amount == balance_after1

    @pytest.mark.P2
    def test_withdraw_rewards_gas(self, clients_new_node, reset_environment):
        client0 = clients_new_node[0]
        # client1 = clients_new_node[1]
        address0, _ = client0.economic.account.generate_account(client0.node.web3, init_amount)
        address1, _ = client0.economic.account.generate_account(client0.node.web3, init_amount)
        staking_and_delegate(client0, address0)
        # staking_and_delegate(client1, address1)
        staking_num0 = client0.staking.get_stakingblocknum()
        # staking_num1 = client1.staking.get_stakingblocknum()
        client0.economic.wait_settlement_blocknum(client0.node)
        result = client0.delegate.delegate(0, address1, node_id=client0.node.node_id, amount=delegate_amount)
        assert_code(result, 0)
        # epoch = math.ceil(client0.node.block_number/client0.economic.settlement_size)
        client0.economic.wait_settlement_blocknum(client0.node, 1)
        rewards0 = client0.delegate.get_delegate_reward_by_nodeid(address0)
        log.info('Address {} delegate rewards : {}'.format(address0, rewards0))
        rewards1 = client0.delegate.get_delegate_reward_by_nodeid(address1)
        log.info('Address {} delegate rewards : {}'.format(address1, rewards1))

        balance_before0 = client0.node.eth.getBalance(address0)
        log.info('Address {} before withdraw delegate balance : {}'.format(address0, balance_before0))
        balance_before1 = client0.node.eth.getBalance(address1)
        log.info('Address {} before withdraw delegate balance : {}'.format(address0, balance_before1))
        result = client0.delegate.withdraw_delegate_reward(address1)
        assert_code(result, 0)
        result = client0.delegate.withdraw_delegate_reward(address0)
        assert_code(result, 0)

        balance_after0 = client0.node.eth.getBalance(address0)
        log.info('Address {} before withdraw delegate balance : {}'.format(address0, balance_after0))
        balance_after1 = client0.node.eth.getBalance(address1)
        log.info('Address {} before withdraw delegate balance : {}'.format(address0, balance_after1))
        gas0 = get_getDelegateReward_gas_fee(client0, 1, 2)
        gas1 = get_getDelegateReward_gas_fee(client0, 1, 1)
        assert balance_before0 - gas0 + rewards0 == balance_after0
        assert balance_before1 - gas1 + rewards1 == balance_after1

    @pytest.mark.P2
    def test_IN_DR_005_IN_DR_015(self, clients_new_node, reset_environment):
        client0 = clients_new_node[0]
        client1 = clients_new_node[1]
        address0, _ = client0.economic.account.generate_account(client0.node.web3, init_amount)
        address1, _ = client0.economic.account.generate_account(client0.node.web3, init_amount)
        staking_and_delegate(client0, address0)
        staking_and_delegate(client1, address1)
        staking_num0 = client0.staking.get_stakingblocknum()
        staking_num1 = client1.staking.get_stakingblocknum()
        result = client1.delegate.delegate(0, address1, node_id=client0.node.node_id, amount=delegate_amount)
        assert_code(result, 0)
        client0.economic.wait_settlement_blocknum(client0.node, 1)
        rewards0 = client0.delegate.get_delegate_reward_by_nodeid(address0)
        log.info('Address {} delegate rewards : {}'.format(address0, rewards0))
        assert rewards0 > 0
        rewards1 = client1.delegate.get_delegate_reward_by_nodeid(address1)
        log.info('Address {} delegate rewards : {}'.format(address1, rewards1))
        assert rewards1 > 0
        balance_before0 = client0.node.eth.getBalance(address0)
        log.info('Address {} before withdraw delegate balance : {}'.format(address0, balance_before0))
        balance_before1 = client1.node.eth.getBalance(address1)
        log.info('Address {} before withdraw delegate balance : {}'.format(address0, balance_before1))
        result = client0.delegate.withdrew_delegate(staking_num0, address0, amount=int(delegate_amount/2))
        assert_code(result, 0)
        result = client0.delegate.withdrew_delegate(staking_num0, address1, amount=int(delegate_amount/2))
        assert_code(result, 0)

        data0 = rlp.encode([rlp.encode(int(1005)), rlp.encode(staking_num0), rlp.encode(bytes.fromhex(
            client0.node.node_id)), rlp.encode(int(delegate_amount/2))])
        data1 = rlp.encode([rlp.encode(int(1005)), rlp.encode(staking_num1), rlp.encode(bytes.fromhex(
            client1.node.node_id)), rlp.encode(int(delegate_amount/2))])
        gas0 = (21000 + 6000 + 8000 + 100 + get_the_dynamic_parameter_gas_fee(data0)) * client0.node.eth.gasPrice
        gas1 = (21000 + 6000 + 8000 + 100 + get_the_dynamic_parameter_gas_fee(data1)) * client0.node.eth.gasPrice
        balance_after0 = client0.node.eth.getBalance(address0)
        log.info('Address {} before withdraw delegate balance : {}'.format(address0, balance_after0))
        balance_after1 = client0.node.eth.getBalance(address1)
        log.info('Address {} before withdraw delegate balance : {}'.format(address0, balance_after1))

        rewards0_after = client0.delegate.get_delegate_reward_by_nodeid(address0)
        log.info('Address {} delegate rewards : {}'.format(address0, rewards0))
        rewards1_after = client0.delegate.get_delegate_reward_by_nodeid(address1)
        log.info('Address {} delegate rewards : {}'.format(address1, rewards1))
        assert rewards0 == rewards0_after
        assert rewards1 == rewards1_after
        assert balance_before0 - gas0 + int(delegate_amount/2) == balance_after0
        assert balance_before1 - gas1 + int(delegate_amount/2) == balance_after1

        result = client0.delegate.withdraw_delegate_reward(address0)
        assert_code(result, 0)
        result = client0.delegate.withdraw_delegate_reward(address1)
        assert_code(result, 0)
        balance_after0_withdraw = client0.node.eth.getBalance(address0)
        log.info('Address after withdraw rewards balance : {}'.format(balance_after0_withdraw))
        balance_after1_withdraw = client0.node.eth.getBalance(address1)
        log.info('Address after withdraw rewards balance : {}'.format(balance_after1_withdraw))
        gas0 = get_getDelegateReward_gas_fee(client0, 1, 0)
        gas1 = get_getDelegateReward_gas_fee(client0, 2, 1)
        assert balance_after0 - gas0 + rewards0 == balance_after0_withdraw
        assert balance_after1 - gas1 + rewards1 == balance_after1_withdraw

    @pytest.mark.P0
    def test_IN_DR_001_IN_DR_002_IN_DR_016_IN_DR_017(self, clients_new_node):
        client1 = clients_new_node[0]
        client2 = clients_new_node[1]
        address1, _ = client1.economic.account.generate_account(client1.node.web3, init_amount)
        address2, _ = client1.economic.account.generate_account(client1.node.web3, init_amount)
        result = client1.delegate.withdraw_delegate_reward(address1)
        assert_code(result, 305001)
        gas = get_getDelegateReward_gas_fee(client1, 0, 0)
        balance_after_withdraw_reward = client1.node.eth.getBalance(address1)
        log.info('Address {} after withdraw reward balance {}'.format(address1, balance_after_withdraw_reward))
        assert init_amount - gas == balance_after_withdraw_reward
        staking_and_delegate(client1, address1)
        staking_and_delegate(client2, address1)
        result = client2.delegate.delegate(0, address2, amount=10**18 * 1000)
        assert_code(result, 0)

        client1.economic.wait_settlement_blocknum(client1.node, 1)
        reward_address1 = client1.delegate.get_delegate_reward_by_nodeid(address1)
        log.info('Address {} reward : {}'.format(address1, reward_address1))
        assert reward_address1 > 0
        reward_address2 = client1.delegate.get_delegate_reward_by_nodeid(address2)
        log.info('Address {} reward : {}'.format(address2, reward_address2))
        assert reward_address2 > 0
        balance_before_address1 = client1.node.eth.getBalance(address1)
        log.info('Before withdraw reward Address {} balance : {}'.format(address1, balance_before_address1))
        balance_before_address2 = client1.node.eth.getBalance(address2)
        log.info('Before withdraw reward Address {} balance : {}'.format(address2, balance_before_address2))
        result = client1.delegate.withdraw_delegate_reward(address1)
        assert_code(result, 0)
        result = client1.delegate.withdraw_delegate_reward(address2)
        assert_code(result, 0)
        gas_address1 = get_getDelegateReward_gas_fee(client1, 2, 2)
        log.info('Withdraw reward address gas : {}'.format(gas_address1))
        gas_address2 = get_getDelegateReward_gas_fee(client1, 1, 1)
        log.info('Withdraw reward address gas : {}'.format(gas_address2))

        balance_after_address1 = client1.node.eth.getBalance(address1)
        log.info('After withdraw reward Address {} balance : {}'.format(address1, balance_after_address1))
        balance_after_address2 = client1.node.eth.getBalance(address2)
        log.info('After withdraw reward Address {} balance : {}'.format(address2, balance_after_address2))

        assert balance_before_address1 + reward_address1 - gas_address1 == balance_after_address1
        assert balance_before_address2 + reward_address2 - gas_address2 == balance_after_address2

        client1.economic.wait_settlement_blocknum(client1.node)
        reward_address1 = client1.delegate.get_delegate_reward_by_nodeid(address1)
        log.info('Address {} reward : {}'.format(address1, reward_address1))
        assert reward_address1 != 0
        reward_address2 = client1.delegate.get_delegate_reward_by_nodeid(address2)
        log.info('Address {} reward : {}'.format(address2, reward_address2))
        assert reward_address2 != 0

        result = client1.delegate.withdraw_delegate_reward(address1)
        assert_code(result, 0)
        result = client1.delegate.withdraw_delegate_reward(address2)
        assert_code(result, 0)
        gas_address1 = get_getDelegateReward_gas_fee(client1, 2, 2)
        log.info('Withdraw reward address gas : {}'.format(gas_address1))
        gas_address2 = get_getDelegateReward_gas_fee(client1, 1, 1)
        log.info('Withdraw reward address gas : {}'.format(gas_address2))

        balance_after_address1_1 = client1.node.eth.getBalance(address1)
        log.info('After withdraw reward Address {} balance : {}'.format(address1, balance_after_address1_1))
        balance_after_address2_1 = client1.node.eth.getBalance(address2)
        log.info('After withdraw reward Address {} balance : {}'.format(address2, balance_after_address2_1))

        assert balance_after_address1 + reward_address1 - gas_address1 == balance_after_address1_1
        assert balance_after_address2 + reward_address2 - gas_address2 == balance_after_address2_1

    @pytest.mark.P2
    def test_IN_DR_006(self, clients_new_node):
        client1 = clients_new_node[0]
        client2 = clients_new_node[1]
        address1, _ = client1.economic.account.generate_account(client1.node.web3, init_amount)
        address2, _ = client1.economic.account.generate_account(client1.node.web3, init_amount)
        staking_and_delegate(client1, address1)
        staking_and_delegate(client2, address1)
        result = client2.delegate.delegate(0, address2, amount=10 ** 18 * 1000)
        assert_code(result, 0)
        wait_block_number(client1.node, 2*client1.economic.settlement_size - 40)
        result = client1.staking.withdrew_staking(client1.node.staking_address, client1.node.node_id)
        assert_code(result, 0)
        candidate_info = client1.ppos.getCandidateInfo(client1.node.node_id)
        log.info('Nodeid {} candiate information : {}'.format(client1.node.node_id, candidate_info))
        wait_block_number(client1.node, 2*client1.economic.settlement_size)
        candidate_info = client1.ppos.getCandidateInfo(client1.node.node_id)
        log.info('Nodeid {} candiate information : {}'.format(client1.node.node_id, candidate_info))
        reward_address1_node1 = client1.delegate.get_delegate_reward_by_nodeid(address1, [client1.node.node_id])
        log.info('Address {} delegate node {} reward {}'.format(address1, client1.node.node_id, reward_address1_node1))
        assert reward_address1_node1 != 0
        reward_address1_node2 = client1.delegate.get_delegate_reward_by_nodeid(address1, [client2.node.node_id])
        log.info('Address {} delegate node {} reward {}'.format(address1, client2.node.node_id, reward_address1_node2))
        assert reward_address1_node2 != 0
        reward_address2 = client1.delegate.get_delegate_reward_by_nodeid(address2)
        log.info('Address {} delegate reward : {}'.format(address2, reward_address2))
        assert reward_address2 != 0
        balance_before_address1 = client1.node.eth.getBalance(address1)
        log.info('Before withdraw reward, address {} balance {}'.format(address1, balance_before_address1))
        balance_before_address2 = client1.node.eth.getBalance(address2)
        log.info('Before withdraw reward, address {} balance {}'.format(address2, balance_before_address2))
        result = client1.delegate.withdraw_delegate_reward(address1)
        assert_code(result, 0)
        result = client1.delegate.withdraw_delegate_reward(address2)
        assert_code(result, 0)
        gas_address1 = get_getDelegateReward_gas_fee(client1, 2, 2)
        gas_address2 = get_getDelegateReward_gas_fee(client1, 1, 1)
        balance_after_address1 = client1.node.eth.getBalance(address1)
        log.info('Address {} balance : {}'.format(address1, balance_after_address1))
        balance_after_address2 = client1.node.eth.getBalance(address2)
        log.info('Address {} balance : {}'.format(address2, balance_after_address2))
        assert balance_before_address1 + reward_address1_node1 + reward_address1_node2 - gas_address1 == balance_after_address1
        assert balance_before_address2 + reward_address2 - gas_address2 == balance_after_address2

    @pytest.mark.P2
    def test_IN_DR_007(self, clients_new_node):
        client1 = clients_new_node[0]
        client2 = clients_new_node[1]
        address1, _ = client1.economic.account.generate_account(client1.node.web3, init_amount)
        address2, _ = client1.economic.account.generate_account(client1.node.web3, init_amount)
        staking_and_delegate(client1, address1)
        staking_and_delegate(client2, address1)
        result = client2.delegate.delegate(0, address2, amount=10 ** 18 * 1000)
        assert_code(result, 0)
        wait_block_number(client1.node, 2*client1.economic.settlement_size - 40)
        result = client1.staking.withdrew_staking(client1.node.staking_address, client1.node.node_id)
        assert_code(result, 0)
        wait_block_number(client1.node, 2*client1.economic.settlement_size)
        reward_address1_node1 = client1.delegate.get_delegate_reward_by_nodeid(address1, [client1.node.node_id])
        log.info('Address {} delegate node {} reward {}'.format(address1, client1.node.node_id, reward_address1_node1))
        assert reward_address1_node1 != 0
        reward_address1_node2 = client1.delegate.get_delegate_reward_by_nodeid(address1, [client2.node.node_id])
        log.info('Address {} delegate node {} reward {}'.format(address1, client2.node.node_id, reward_address1_node2))
        assert reward_address1_node2 != 0
        reward_address2 = client1.delegate.get_delegate_reward_by_nodeid(address2)
        log.info('Address {} delegate reward : {}'.format(address2, reward_address2))
        assert reward_address2 != 0
        wait_block_number(client1.node, (
                2 + client1.economic.genesis.economicModel.staking.unStakeFreezeDuration)*client1.economic.settlement_size)
        balance_before_address1 = client1.node.eth.getBalance(address1)
        log.info('Before withdraw reward, address {} balance {}'.format(address1, balance_before_address1))
        balance_before_address2 = client1.node.eth.getBalance(address2)
        log.info('Before withdraw reward, address {} balance {}'.format(address2, balance_before_address2))

        reward_address1_node1_2 = client1.delegate.get_delegate_reward_by_nodeid(address1, [client1.node.node_id])
        log.info('Address {} delegate node {} reward {}'.format(address1, client1.node.node_id, reward_address1_node1_2))
        assert reward_address1_node1 == reward_address1_node1_2
        reward_address1_node2_2 = client1.delegate.get_delegate_reward_by_nodeid(address1, [client2.node.node_id])
        log.info('Address {} delegate node {} reward {}'.format(address1, client2.node.node_id, reward_address1_node2_2))
        assert reward_address1_node2 < reward_address1_node2_2
        reward_address2_2 = client1.delegate.get_delegate_reward_by_nodeid(address2)
        log.info('Address {} delegate node {} reward {}'.format(address2, client2.node.node_id, reward_address2_2))
        assert reward_address2 < reward_address2_2

        result = client1.delegate.withdraw_delegate_reward(address1)
        assert_code(result, 0)
        result = client1.delegate.withdraw_delegate_reward(address2)
        assert_code(result, 0)
        gas_address1 = get_getDelegateReward_gas_fee(client1, 2, 4)
        gas_address2 = get_getDelegateReward_gas_fee(client1, 1, 3)
        balance_after_address1 = client1.node.eth.getBalance(address1)
        log.info('Address {} balance : {}'.format(address1, balance_after_address1))
        balance_after_address2 = client1.node.eth.getBalance(address2)
        log.info('Address {} balance : {}'.format(address2, balance_after_address2))

        assert balance_before_address1 + reward_address1_node1_2 + reward_address1_node2_2 - gas_address1 == balance_after_address1
        assert balance_before_address2 + reward_address2_2 - gas_address2 == balance_after_address2

    @pytest.mark.P2
    def test_IN_DR_008(self, clients_new_node):
        client1 = clients_new_node[0]
        client2 = clients_new_node[1]
        address1, _ = client1.economic.account.generate_account(client1.node.web3, init_amount)
        address2, _ = client1.economic.account.generate_account(client1.node.web3, init_amount)
        staking_and_delegate(client1, address1)
        staking_and_delegate(client2, address1)
        result = client1.delegate.delegate(0, address2, amount=10 ** 18 * 1000)
        assert_code(result, 0)
        wait_block_number(client1.node, 2 * client1.economic.settlement_size - 40)
        report_information = mock_duplicate_sign(1, client1.node.nodekey, client1.node.blsprikey,
                                                 2 * client1.economic.settlement_size - 50)
        log.info("Report information: {}".format(report_information))
        address, _ = client1.economic.account.generate_account(client1.node.web3, 10 ** 18 * 1000)
        result = client2.duplicatesign.reportDuplicateSign(1, report_information, address)
        log.info('Node duplicate block result : {}'.format(result))
        assert_code(result, 0)
        wait_block_number(client2.node, 2 * client1.economic.settlement_size)
        reward_address1_node1 = client2.delegate.get_delegate_reward_by_nodeid(address1, [client1.node.node_id])
        log.info('Address {} delegate node {} reward is {}'.format(address1, client1.node.node_id, reward_address1_node1))
        assert reward_address1_node1 == 0
        reward_address1_node2 = client2.delegate.get_delegate_reward_by_nodeid(address1, [client2.node.node_id])
        log.info(
            'Address {} delegate node {} reward is {}'.format(address1, client2.node.node_id, reward_address1_node2))
        assert reward_address1_node2 > 0
        reward_address2 = client2.delegate.get_delegate_reward_by_nodeid(address2)
        log.info('Address {} reward is {}'.format(address2, reward_address2))
        assert reward_address2 == 0
        balance_address1_before_withdraw = client1.node.eth.getBalance(address1)
        log.info('Address {} before withdraw balance {}'.format(address1, balance_address1_before_withdraw))
        balance_address2_before_withdraw = client1.node.eth.getBalance(address2)
        log.info('Address {} before withdraw balance {}'.format(address2, balance_address2_before_withdraw))
        result = client1.delegate.withdraw_delegate_reward(address1)
        assert_code(result, 0)
        result = client1.delegate.withdraw_delegate_reward(address2)
        assert_code(result, 0)
        gas1 = get_getDelegateReward_gas_fee(client1, 2, 1)
        gas2 = get_getDelegateReward_gas_fee(client2, 1, 0)
        balance_address1_after_withdraw = client1.node.eth.getBalance(address1)
        log.info('Address {} before withdraw balance {}'.format(address1, balance_address1_after_withdraw))
        balance_address2_after_withdraw = client1.node.eth.getBalance(address2)
        log.info('Address {} before withdraw balance {}'.format(address2, balance_address2_after_withdraw))
        assert balance_address1_before_withdraw - gas1 + reward_address1_node2 == balance_address1_after_withdraw
        assert balance_address2_before_withdraw - gas2 == balance_address2_after_withdraw

    @pytest.mark.P2
    def test_IN_DR_010(self, clients_new_node):
        client1 = clients_new_node[0]
        client2 = clients_new_node[1]
        address1, _ = client1.economic.account.generate_account(client1.node.web3, init_amount)
        address2, _ = client1.economic.account.generate_account(client1.node.web3, init_amount)
        staking_and_delegate(client1, address1)
        staking_and_delegate(client2, address1)
        stakingnum1 = client1.staking.get_stakingblocknum(client1.node)
        result = client1.delegate.delegate(0, address2, amount=10 ** 18 * 1000)
        assert_code(result, 0)
        client1.economic.wait_settlement_blocknum(client1.node)
        verifier_list = get_pledge_list(client1.ppos.getVerifierList)
        assert client1.node.node_id in verifier_list
        assert client2.node.node_id in verifier_list
        wait_block_number(client1.node, 2*client1.economic.settlement_size - 40)

        result = client1.delegate.withdrew_delegate(stakingnum1, address1, amount=10**18 * 1000)
        assert_code(result, 0)
        result = client1.delegate.withdrew_delegate(stakingnum1, address2, amount=10**18 * 500)
        assert_code(result, 0)
        result = client1.ppos.getDelegateReward(address1, [client1.node.node_id])
        log.info('Address {} delegate nodeid {} reward {}'.format(address1, client1.node.node_id, result))
        assert_code(result, 305001)
        result = client1.ppos.getDelegateReward(address1, [client2.node.node_id])
        log.info(
            'Address {} delegate nodeid {} reward {}'.format(address1, client1.node.node_id, result))
        assert_code(result, 0)
        assert result.get('Ret')[0].get('reward') == 0
        delegateinfo_address2 = client1.ppos.getDelegateInfo(stakingnum1, address2, client1.node.node_id)
        log.info('Address {} delegate information : {}'.format(address2, delegateinfo_address2))
        result = client1.ppos.getDelegateReward(address2)
        assert_code(result, 0)
        assert result.get('Ret')[0].get('reward') == 0

        client1.economic.wait_settlement_blocknum(client1.node)

        result = client1.ppos.getDelegateReward(address1, [client1.node.node_id])
        log.info('Address {} delegate nodeid {} reward {}'.format(address1, client1.node.node_id, result))
        assert_code(result, 305001)
        result = client1.ppos.getDelegateReward(address1, [client2.node.node_id])
        log.info(
            'Address {} delegate nodeid {} reward {}'.format(address1, client1.node.node_id, result))
        assert_code(result, 0)
        assert result.get('Ret')[0].get('reward') > 0
        result = client1.ppos.getDelegateReward(address2)
        log.info('Address {} delegate nodeid {} reward {}'.format(address2, client1.node.node_id, result))
        assert_code(result, 0)
        assert result.get('Ret')[0].get('reward') > 0

    @pytest.mark.P2
    def test_IN_DR_012(self, clients_new_node):
        client1 = clients_new_node[0]
        client2 = clients_new_node[1]
        address1, _ = client1.economic.account.generate_account(client1.node.web3, init_amount)
        address2, _ = client1.economic.account.generate_account(client1.node.web3, init_amount)
        staking_and_delegate(client1, address1)
        staking_and_delegate(client2, address1)
        result = client1.delegate.delegate(0, address2, amount=10 ** 18 * 1000)
        assert_code(result, 0)
        client1.economic.wait_settlement_blocknum(client1.node)
        verifier_list = get_pledge_list(client1.ppos.getVerifierList)
        assert client1.node.node_id in verifier_list
        assert client2.node.node_id in verifier_list
        wait_block_number(client1.node, 2*client1.economic.settlement_size - 40)
        result = client1.delegate.delegate(0, address1, amount=10 ** 18 * 1000)
        assert_code(result, 0)

        result = client1.delegate.delegate(0, address2, amount=10 ** 18 * 1000)
        assert_code(result, 0)
        client1.economic.wait_settlement_blocknum(client1.node)

        reward_address1_node1 = client1.delegate.get_delegate_reward_by_nodeid(address1, [client1.node.node_id])
        log.info('Address {} delegate nodeid {} reward {}'.format(address1, client1.node.node_id, reward_address1_node1))
        assert reward_address1_node1 > 0
        reward_address1_node2 = client1.delegate.get_delegate_reward_by_nodeid(address1, [client2.node.node_id])
        log.info(
            'Address {} delegate nodeid {} reward {}'.format(address1, client2.node.node_id, reward_address1_node2))
        assert reward_address1_node2 > 0
        reward_address2 = client1.delegate.get_delegate_reward_by_nodeid(address2, [client1.node.node_id])
        log.info(
            'Address {} delegate nodeid {} reward {}'.format(address2, client1.node.node_id, reward_address2))
        assert reward_address2 > 0

        balance_address1_before_withdraw = client1.node.eth.getBalance(address1)
        log.info('Address {} before withdraw  balance {}'.format(address1, balance_address1_before_withdraw))
        result = client1.delegate.withdraw_delegate_reward(address1)
        assert_code(result, 0)
        balance_address1 = client1.node.eth.getBalance(address1)
        log.info('Address {} after withdraw balance {}'.format(address1, balance_address1))
        gas = get_getDelegateReward_gas_fee(client1, 2, 2)
        assert balance_address1_before_withdraw - gas + reward_address1_node1 + reward_address1_node2 == balance_address1

        balance_address2_before_withdraw = client1.node.eth.getBalance(address2)
        log.info('Address {} before withdraw  balance {}'.format(address2, balance_address2_before_withdraw))
        result = client1.delegate.withdraw_delegate_reward(address2)
        assert_code(result, 0)
        balance_address2 = client1.node.eth.getBalance(address2)
        log.info('Address {} after withdraw balance {}'.format(address2, balance_address2))
        gas = get_getDelegateReward_gas_fee(client1, 1, 1)
        assert balance_address2_before_withdraw - gas + reward_address2 == balance_address2

    @pytest.mark.P2
    def test_IN_DR_013(self, clients_new_node):
        client1 = clients_new_node[0]
        client2 = clients_new_node[1]
        address1, _ = client1.economic.account.generate_account(client1.node.web3, init_amount)
        address2, _ = client1.economic.account.generate_account(client1.node.web3, init_amount)
        staking_and_delegate(client1, address1)
        staking_and_delegate(client2, address1)
        result = client1.delegate.delegate(0, address2, amount=10 ** 18 * 1000)
        assert_code(result, 0)
        client1.economic.wait_settlement_blocknum(client1.node)
        verifier_list = get_pledge_list(client1.ppos.getVerifierList)
        assert client1.node.node_id in verifier_list
        assert client2.node.node_id in verifier_list
        wait_block_number(client1.node, 2*client1.economic.settlement_size - 40)

        client1.economic.wait_settlement_blocknum(client1.node)

        reward_address1_node1 = client1.delegate.get_delegate_reward_by_nodeid(address1, [client1.node.node_id])
        log.info('Address {} delegate nodeid {} reward {}'.format(address1, client1.node.node_id, reward_address1_node1))
        assert reward_address1_node1 > 0
        reward_address1_node2 = client1.delegate.get_delegate_reward_by_nodeid(address1, [client2.node.node_id])
        log.info(
            'Address {} delegate nodeid {} reward {}'.format(address1, client2.node.node_id, reward_address1_node2))
        assert reward_address1_node2 > 0
        reward_address2 = client1.delegate.get_delegate_reward_by_nodeid(address2, [client1.node.node_id])
        log.info(
            'Address {} delegate nodeid {} reward {}'.format(address2, client1.node.node_id, reward_address2))
        assert reward_address2 > 0

        balance_address1_before_withdraw = client1.node.eth.getBalance(address1)
        log.info('Address {} before withdraw  balance {}'.format(address1, balance_address1_before_withdraw))
        result = client1.delegate.withdraw_delegate_reward(address1)
        assert_code(result, 0)
        balance_address1 = client1.node.eth.getBalance(address1)
        log.info('Address {} after withdraw balance {}'.format(address1, balance_address1))
        gas = get_getDelegateReward_gas_fee(client1, 2, 2)
        assert balance_address1_before_withdraw - gas + reward_address1_node1 + reward_address1_node2 == balance_address1

        balance_address2_before_withdraw = client1.node.eth.getBalance(address2)
        log.info('Address {} before withdraw  balance {}'.format(address2, balance_address2_before_withdraw))
        result = client1.delegate.withdraw_delegate_reward(address2)
        assert_code(result, 0)
        balance_address2 = client1.node.eth.getBalance(address2)
        log.info('Address {} after withdraw balance {}'.format(address2, balance_address2))
        gas = get_getDelegateReward_gas_fee(client1, 1, 1)
        assert balance_address2_before_withdraw - gas + reward_address2 == balance_address2


class TestGas:
    @pytest.mark.P1
    def test_IN_GA_001_IN_GA_002(self, clients_new_node, reset_environment):
        client1 = clients_new_node[0]
        client2 = clients_new_node[1]
        address1, _ = client1.economic.account.generate_account(client1.node.web3, init_amount)
        address2, _ = client1.economic.account.generate_account(client1.node.web3, init_amount)
        balance = client1.node.eth.getBalance(address1)
        print(balance)
        staking_and_delegate(client1, address1)
        stakingnum = client1.staking.get_stakingblocknum()
        balance_address1 = client1.node.eth.getBalance(address1)
        log.info('Address {} balance : {}'.format(address1, balance_address1))
        data = rlp.encode([rlp.encode(int(1004)), rlp.encode(0), rlp.encode(bytes.fromhex(client1.node.node_id)),
                           rlp.encode(10 ** 18 * 1000)])
        gas = (21000 + 6000 + 16000 + get_the_dynamic_parameter_gas_fee(data)) * client1.node.eth.gasPrice
        assert balance - delegate_amount - gas == balance_address1
        time.sleep(10)

        balance_address1 = client2.node.eth.getBalance(address1)
        print('balance_address1: {}'.format(balance_address1))
        staking_and_delegate(client2, address1)
        balance_address1_1 = client2.node.eth.getBalance(address1)
        log.info('Address {} balance : {}'.format(address1, balance_address1_1))
        data = rlp.encode([rlp.encode(int(1004)), rlp.encode(0), rlp.encode(bytes.fromhex(client1.node.node_id)),
                           rlp.encode(10 ** 18 * 1000)])
        gas = (21000 + 6000 + 16000 + get_the_dynamic_parameter_gas_fee(data)) * client1.node.eth.gasPrice
        transaction_data = {"to": EconomicConfig.STAKING_ADDRESS, "data": data, "from": address1}
        estimated_gas = client1.node.eth.estimateGas(transaction_data)
        print(gas, estimated_gas)
        assert balance_address1 - gas - delegate_amount == balance_address1_1

        client1.delegate.withdrew_delegate(stakingnum, address1, amount=10**18 * 100)
        banlance_address1_after_withdraw = client1.node.eth.getBalance(address1)
        log.info('Address {} after withdraw balance : {}'.format(address1, banlance_address1_after_withdraw))
        data = rlp.encode([rlp.encode(int(1005)), rlp.encode(stakingnum), rlp.encode(bytes.fromhex(client1.node.node_id)),
                           rlp.encode(10**18 * 100)])
        gas = (21000 + 6000 + 8000 + get_the_dynamic_parameter_gas_fee(data)) * client1.node.eth.gasPrice
        assert balance_address1_1 - gas + 10**18 * 100 == banlance_address1_after_withdraw

        result = client1.delegate.delegate(0, address2, amount=10 ** 18 * 1000)
        assert_code(result, 0)
        balance_address2 = client1.node.eth.getBalance(address2)
        log.info('Address {} balance : {}'.format(address2, balance_address2))
        data = rlp.encode([rlp.encode(int(1004)), rlp.encode(0), rlp.encode(bytes.fromhex(client1.node.node_id)),
                           rlp.encode(10 ** 18 * 1000)])
        gas = (21000 + 6000 + 16000 + get_the_dynamic_parameter_gas_fee(data)) * client1.node.eth.gasPrice
        assert 10 ** 18 * 9000 - gas == balance_address1
        client1.economic.wait_settlement_blocknum(client1.node)
        verifier_list = get_pledge_list(client1.ppos.getVerifierList)
        assert client1.node.node_id in verifier_list
        assert client2.node.node_id in verifier_list
        wait_block_number(client1.node, 2 * client1.economic.settlement_size - 40)
        result = client1.delegate.delegate(0, address1, amount=10 ** 18 * 1000)
        assert_code(result, 0)
        balance_address1_2 = client1.node.eth.getBalance(address1)
        log.info('Address {} balance {}'.format(address1, balance_address2))
        assert banlance_address1_after_withdraw - gas - 10 ** 18 * 1000 == balance_address1_2


        result = client1.delegate.withdrew_delegate(stakingnum, address1, amount=10**18 * 100)
        assert_code(result, 0)
        banlance_address1_after_withdraw_1 = client1.node.eth.getBalance(address1)
        log.info('Address {} after withdraw balance : {}'.format(address1, banlance_address1_after_withdraw))
        data = rlp.encode([rlp.encode(int(1005)), rlp.encode(stakingnum), rlp.encode(bytes.fromhex(client1.node.node_id)),
                           rlp.encode(10**18 * 100)])
        gas = (21000 + 6000 + 8000 + get_the_dynamic_parameter_gas_fee(data)) * client1.node.eth.gasPrice
        assert balance_address1_2 + 10**18 * 100 - gas == banlance_address1_after_withdraw_1
        client1.economic.wait_settlement_blocknum(client1.node)

        result = client1.delegate.delegate(0, address1, amount=10**18*1000)
        assert_code(result, 0)
        balance_address1_3 = client1.node.eth.getBalance(address1)
        log.info('Address {} balance : {}'.format(address1, balance_address1_3))
        data = rlp.encode([rlp.encode(int(1004)), rlp.encode(0), rlp.encode(bytes.fromhex(client1.node.node_id)),
                           rlp.encode(10 ** 18 * 1000)])
        gas = (21000 + 6000 + 16000 + get_the_dynamic_parameter_gas_fee(data) + 100) * client1.node.eth.gasPrice
        assert banlance_address1_after_withdraw_1 - 10**18*1000 - gas == balance_address1_3

        result = client1.delegate.withdrew_delegate(stakingnum, address1, amount=10**18*100)
        assert_code(result, 0)
        balance_address1_4 = client1.node.eth.getBalance(address1)
        data = rlp.encode(
            [rlp.encode(int(1005)), rlp.encode(stakingnum), rlp.encode(bytes.fromhex(client1.node.node_id)),
             rlp.encode(10 ** 18 * 100)])
        gas = (21000 + 6000 + 8000 + get_the_dynamic_parameter_gas_fee(data)) * client1.node.eth.gasPrice
        assert balance_address1_3 + 10**18*100 - gas == balance_address1_4

        client1.economic.wait_settlement_blocknum(client1.node)

        result = client1.delegate.delegate(0, address1, amount=10**18 * 100)
        assert_code(result, 0)
        balance_address1_5 = client1.node.eth.getBalance(address1)
        data = rlp.encode([rlp.encode(int(1004)), rlp.encode(0), rlp.encode(bytes.fromhex(client1.node.node_id)),
                           rlp.encode(10 ** 18 * 100)])
        gas = (21000 + 6000 + 16000 + get_the_dynamic_parameter_gas_fee(data) + 100) * client1.node.eth.gasPrice
        assert balance_address1_4 - 10**18 * 100 - gas == balance_address1_5

        result = client1.delegate.withdrew_delegate(stakingnum, address1, amount=10**18 * 100)
        assert_code(result, 0)
        balance_address1_6 = client1.node.eth.getBalance(address1)
        data = rlp.encode(
            [rlp.encode(int(1005)), rlp.encode(stakingnum), rlp.encode(bytes.fromhex(client1.node.node_id)),
             rlp.encode(10 ** 18 * 100)])
        gas = (21000 + 6000 + 8000 + get_the_dynamic_parameter_gas_fee(data)) * client1.node.eth.gasPrice
        assert balance_address1_5 + 10**18 * 100 - gas == balance_address1_6

    @pytest.mark.P2
    def test_IN_GA_003_IN_GA_004(self, client_new_node):
        client = client_new_node
        address, _ = client.economic.account.generate_account(client.node.web3, init_amount)
        staking_and_delegate(client, address)
        stakingnum = client.staking.get_stakingblocknum()
        balance_address = client.node.eth.getBalance(address)
        log.info('Address {} balance : {}'.format(address, balance_address))
        data = rlp.encode([rlp.encode(int(1004)), rlp.encode(0), rlp.encode(bytes.fromhex(client.node.node_id)),
                           rlp.encode(10 ** 18 * 1000)])
        gas = (21000 + 6000 + 16000 + get_the_dynamic_parameter_gas_fee(data)) * client.node.eth.gasPrice
        assert 10 ** 18 * 9000 - gas == balance_address

        client.economic.wait_settlement_blocknum(client.node, 2)

        result = client.delegate.withdrew_delegate(stakingnum, address, amount=10 ** 18 * 900)
        assert_code(result, 0)
        balance_address_1 = client.node.eth.getBalance(address)
        log.info('Address {} balance : {}'.format(address, balance_address_1))
        data = rlp.encode(
            [rlp.encode(int(1005)), rlp.encode(stakingnum), rlp.encode(bytes.fromhex(client.node.node_id)),
             rlp.encode(10 ** 18 * 900)])
        gas = (21000 + 6000 + 8000 + get_the_dynamic_parameter_gas_fee(data) + 200) * client.node.eth.gasPrice
        assert balance_address - gas + 10 ** 18 * 900 == balance_address_1

    @pytest.mark.P2
    def test_IN_GA_005_IN_GA_006(self, client_new_node):
        client = client_new_node
        address, _ = client.economic.account.generate_account(client.node.web3, init_amount)
        staking_and_delegate(client, address)
        stakingnum = client.staking.get_stakingblocknum()
        balance_address = client.node.eth.getBalance(address)
        log.info('Address {} balance : {}'.format(address, balance_address))
        data = rlp.encode([rlp.encode(int(1004)), rlp.encode(0), rlp.encode(bytes.fromhex(client.node.node_id)),
                           rlp.encode(10 ** 18 * 1000)])
        gas = (21000 + 6000 + 16000 + get_the_dynamic_parameter_gas_fee(data)) * client.node.eth.gasPrice
        assert 10 ** 18 * 9000 - gas == balance_address

        client.economic.wait_settlement_blocknum(client.node, 3)
        reward = client.delegate.get_delegate_reward_by_nodeid(address)
        log.info('Address {} delegate reward {}'.format(address, reward))

        result = client.delegate.withdrew_delegate(stakingnum, address, amount=10 ** 18 * 1000)
        assert_code(result, 0)
        balance_address_1 = client.node.eth.getBalance(address)
        log.info('Address {} balance : {}'.format(address, balance_address_1))
        data = rlp.encode(
            [rlp.encode(int(1005)), rlp.encode(stakingnum), rlp.encode(bytes.fromhex(client.node.node_id)),
             rlp.encode(10 ** 18 * 1000)])
        gas = (21000 + 6000 + 8000 + get_the_dynamic_parameter_gas_fee(data) + 300) * client.node.eth.gasPrice
        assert balance_address - gas + 10 ** 18 * 1000 + reward == balance_address_1

    @pytest.mark.P2
    def test_IN_GA_007_IN_GA_008(self, client_new_node, reset_environment):
        client = client_new_node
        address, _ = client.economic.account.generate_account(client.node.web3, init_amount)
        staking_and_delegate(client, address)
        stakingnum = client.staking.get_stakingblocknum()
        balance_address = client.node.eth.getBalance(address)
        log.info('Address {} balance : {}'.format(address, balance_address))
        data = rlp.encode([rlp.encode(int(1004)), rlp.encode(0), rlp.encode(bytes.fromhex(client.node.node_id)),
                           rlp.encode(10 ** 18 * 1000)])
        gas = (21000 + 6000 + 16000 + get_the_dynamic_parameter_gas_fee(data)) * client.node.eth.gasPrice
        assert 10 ** 18 * 9000 - gas == balance_address

        client.economic.wait_settlement_blocknum(client.node, 2)

        data = rlp.encode(
            [rlp.encode(int(1005)), rlp.encode(stakingnum), rlp.encode(bytes.fromhex(client.node.node_id)),
             rlp.encode(10 ** 18 * 900)])
        gas = 21000 + 6000 + 8000 + get_the_dynamic_parameter_gas_fee(data) + 99

        try:
            client.delegate.withdrew_delegate(stakingnum, address, amount=10 ** 18 * 900,
                                                   transaction_cfg={'gas': gas})
        except IndexError as e:
            assert str(e) == "list index out of range"

        gas = 21000 + get_the_dynamic_parameter_gas_fee(data)
        try:
            client.delegate.withdrew_delegate(stakingnum, address, amount=10 ** 18 * 900,
                                              transaction_cfg={'gas': gas})
        except IndexError as e:
            assert str(e) == "list index out of range"

        gas = 21000 + get_the_dynamic_parameter_gas_fee(data) - 1
        try:
            client.delegate.withdrew_delegate(stakingnum, address, amount=10 ** 18 * 900,
                                              transaction_cfg={'gas': gas})
        except ValueError as e:
            assert e.args[0].get('message') == "intrinsic gas too low"

    @pytest.mark.P2
    def test_IN_GA_009_IN_GA_010(self, client_new_node, reset_environment):
        client = client_new_node
        a = get_getDelegateReward_gas_fee(client, 1, 1)
        log.info(a)
        log.info(client.node.eth.gasPrice)
        address, _ = client.economic.account.generate_account(client.node.web3, init_amount)
        data = rlp.encode([rlp.encode(int(5000))])

        gas = 21000 + get_the_dynamic_parameter_gas_fee(data) - 1
        try:
            client.delegate.withdraw_delegate_reward(address, transaction_cfg={'gas': gas})
        except ValueError as e:
            assert e.args[0].get('message') == "intrinsic gas too low"
        gas = 21000 + get_the_dynamic_parameter_gas_fee(data)
        try:
            client.delegate.withdraw_delegate_reward(address, transaction_cfg={'gas': gas})
        except IndexError as e:
            assert str(e) == "list index out of range"

        gas = 21000 + 3000 + 8000 + get_the_dynamic_parameter_gas_fee(data) - 1
        try:
            client.delegate.withdraw_delegate_reward(address, transaction_cfg={'gas': gas})
        except IndexError as e:
            assert str(e) == "list index out of range"

        result = client.delegate.withdraw_delegate_reward(address, transaction_cfg={'gas': gas + 1})
        assert_code(result, 305001)

        staking_and_delegate(client, address)
        gas = 21000 + 3000 + 8000 + get_the_dynamic_parameter_gas_fee(data) + 1000
        balance_before_fail = client.node.eth.getBalance(address)
        log.info('Before withdraw reward--fail address balance : {}'.format(balance_before_fail))
        try:
            client.delegate.withdraw_delegate_reward(address, transaction_cfg={'gas': gas - 1})
        except IndexError as e:
            assert str(e) == "list index out of range"
        balance_before = client.node.eth.getBalance(address)
        log.info('Before withdraw reward address balance : {}'.format(balance_before))
        assert balance_before_fail - (gas - 1) * client.node.eth.gasPrice == balance_before
        result = client.delegate.withdraw_delegate_reward(address, transaction_cfg={'gas': gas})
        assert_code(result, 0)

        balance_after = client.node.eth.getBalance(address)
        log.info('After withdraw reward address balance : {}'.format(balance_after))
        assert balance_before - gas*client.node.eth.gasPrice == balance_after

    @pytest.mark.P2
    def test_IN_GA_011_IN_GA_012(self, clients_new_node):
        client0 = clients_new_node[0]
        client1 = clients_new_node[1]
        address0, _ = client0.economic.account.generate_account(client0.node.web3, init_amount)
        address1, _ = client0.economic.account.generate_account(client0.node.web3, init_amount)
        staking_and_delegate(client0, address0)
        staking_and_delegate(client1, address1)
        staking_num0 = client0.staking.get_stakingblocknum(client0.node)
        staking_num1 = client1.staking.get_stakingblocknum(client1.node)
        client0.economic.wait_settlement_blocknum(client0.node, 2)
        balance_before_withdraw0 = client0.node.eth.getBalance(address0)
        log.info('Before withdraw address {} balance {}'.format(address0, balance_before_withdraw0))
        balance_before_withdraw1 = client0.node.eth.getBalance(address1)
        log.info('Before withdraw address {} balance {}'.format(address1, balance_before_withdraw1))
        result = client0.delegate.withdrew_delegate(staking_num0, address0, amount=int(delegate_amount/2))
        assert_code(result, 0)
        rewards = client1.delegate.get_delegate_reward_by_nodeid(address1)
        log.info('Address {} rewards {}'.format(address1, rewards))
        assert rewards > 0
        result = client1.delegate.withdrew_delegate(staking_num1, address1, amount=delegate_amount)
        assert_code(result, 0)
        data0 = rlp.encode([rlp.encode(int(1005)), rlp.encode(staking_num0), rlp.encode(bytes.fromhex(
            client0.node.node_id)), rlp.encode(int(delegate_amount/2))])
        data1 = rlp.encode([rlp.encode(int(1005)), rlp.encode(staking_num1), rlp.encode(bytes.fromhex(
            client1.node.node_id)), rlp.encode(delegate_amount)])
        gas0 = (21000 + 8000 + 6000 + get_the_dynamic_parameter_gas_fee(data0) + 200) * client0.node.eth.gasPrice
        gas1 = (21000 + 8000 + 6000 + get_the_dynamic_parameter_gas_fee(data1) + 200) * client0.node.eth.gasPrice
        balance_after_withdraw0 = client0.node.eth.getBalance(address0)
        log.info('After withdraw address {} balance {}'.format(address0, balance_after_withdraw0))
        balance_after_withdraw1 = client0.node.eth.getBalance(address1)
        log.info('After withdraw address {} balance {}'.format(address1, balance_after_withdraw1))
        assert balance_before_withdraw0 - gas0 + int(delegate_amount/2) == balance_after_withdraw0
        assert balance_before_withdraw1 - gas1 + delegate_amount + rewards == balance_after_withdraw1


class TestNet:
    def deploy_me(self, global_test_env, net):
        global_test_env.deploy_all()
        node = global_test_env.get_a_normal_node()
        log.info('Node id {}'.format(node.node_id))
        test_node = copy(node)
        test_node.clean()
        new_cfg = copy(global_test_env.cfg)
        new_cfg.init_chain = False
        if net == 'main':
            new_cfg.append_cmd = "--main"
        elif net == 'testnet':
            new_cfg.append_cmd = "--testnet"
        elif net == 'rallynet':
            new_cfg.append_cmd = "--rallynet"
        elif net == 'uatnet':
            new_cfg.append_cmd = "--uatnet"
        elif net == 'demonet':
            new_cfg.append_cmd = "--demonet"
        else:
            raise Exception('No {} net'.format(net))
        test_node.cfg = new_cfg
        log.info("start deploy {}".format(test_node.node_mark))
        log.info("is init:{}".format(test_node.cfg.init_chain))
        test_node.deploy_me(genesis_file=None)
        log.info("deploy end")
        return test_node

    # @pytest.mark.P2
    # def test_DD_NE_001(self, global_test_env):
    #     test_node = self.deploy_me(global_test_env, 'main')
    #     try:
    #         log.info(test_node.admin.nodeInfo)
    #     except Exception as e:
    #         e == "20 seconds"
    #     log_file = "/home/platon/trantor_test/node-{}/log/platon.log".format(test_node.p2p_port)
    #     test_node.sftp.get(log_file, conf.LOG_FILE)
    #     with open(conf.LOG_FILE, 'r') as f:
    #         info = f.readlines()
    #         log.info(info[-1])
    #         assert info[-1].strip() == "Fatal: Error starting protocol stack: EconomicModel config is nil"

    @pytest.mark.P2
    def test_DD_NE_002(self, global_test_env):
        test_node = self.deploy_me(global_test_env, 'testnet')
        assert test_node.admin.nodeInfo.get('protocols').get('platon').get('config').get('chainId') == 103
        assert test_node.admin.nodeInfo.get('protocols').get('platon').get('network') == 2000

    @pytest.mark.P2
    def test_DD_NE_003(self, global_test_env):
        pass
        # test_node = self.deploy_me(global_test_env, 'rallynet')
        # assert test_node.admin.nodeInfo.get('protocols').get('platon').get('config').get('chainId') == 94
        # assert test_node.admin.nodeInfo.get('protocols').get('platon').get('network') == 3000

    @pytest.mark.P2
    def test_DD_NE_004(self, global_test_env):
        pass
        # delete for 0.8.0
        # test_node = self.deploy_me(global_test_env, 'uatnet')
        # assert test_node.admin.nodeInfo.get('protocols').get('platon').get('config').get('chainId') == 299
        # assert test_node.admin.nodeInfo.get('protocols').get('platon').get('network') == 4000

    @pytest.mark.P2
    def test_DD_NE_005(self, global_test_env):
        test_node = self.deploy_me(global_test_env, 'demonet')
        try:
            log.info(test_node.admin.nodeInfo)
        except Exception as e:
            e == "20 seconds"
        log_file = "/home/platon/trantor_test/node-{}/log/platon.log".format(test_node.p2p_port)
        test_node.sftp.get(log_file, conf.LOG_FILE)
        with open(conf.LOG_FILE, 'r') as f:
            info = f.readlines()
            log.info(info[-1])
            assert info[-1].strip() == "Fatal: Error starting protocol stack: genesis has no chain configuration"












