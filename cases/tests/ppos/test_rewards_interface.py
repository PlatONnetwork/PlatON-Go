from common.log import log
from tests.lib.utils import assert_code, get_pledge_list
from common.key import mock_duplicate_sign
from tests.lib.utils import wait_block_number, get_the_dynamic_parameter_gas_fee, get_getDelegateReward_gas_fee
import rlp
from typing import List
import time
from tests.govern.test_voting_statistics import createstaking

def get_new_value(value):
    if value == 10000:
        return value - 1
    else:
        return value + 1

def staking_and_delegate(clients, address, amount=10**18 * 1000):
    createstaking(clients, reward_per=1000)
    if isinstance(clients, List):
        clients = clients[0]
    result = clients.delegate.delegate(0, address, amount=amount)
    assert_code(result, 0)

class TestCreateStaking:
    def assert_rewardsper(self, client, nextrewardsper):
        assert_code(client.staking.get_rewardper(client.node), nextrewardsper)
        assert_code(client.staking.get_rewardper(client.node, isnext=True), nextrewardsper)
        value, nextvalue = get_pledge_list(client.ppos.getCandidateList, client.node.node_id)
        assert_code(value, nextrewardsper)
        assert_code(nextvalue, nextrewardsper)

    def test_IV_032_IV_033_IV_037(self, client_new_node):
        staking = client_new_node.staking
        address, _ = staking.economic.account.generate_account(staking.node.web3,
                                                            3 * staking.economic.genesis.economicModel.staking.stakeThreshold)
        result = staking.create_staking(0, address, address, node_id=staking.node.node_id,
                                        amount=2 * staking.economic.genesis.economicModel.staking.stakeThreshold,
                                        reward_per=1155)
        assert_code(result, 0)
        assert_code(staking.get_rewardper(staking.node), 1155)
        assert_code(staking.get_rewardper(staking.node, isnext=True), 1155)
        result = staking.create_staking(0, address, address, node_id=staking.node.node_id,
                                        amount=2 * staking.economic.genesis.economicModel.staking.stakeThreshold,
                                        reward_per=1)
        log.info('Repeat create staking result : {}'.format(result))
        assert_code(result, 301101)
        self.assert_rewardsper(client_new_node, 0, 1155)



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
        self.assert_rewardsper(client_new_node, 0, 1155)

    def test_IV_036(self, client_new_node):
        staking = client_new_node.staking
        address, _ = staking.economic.account.generate_account(staking.node.web3,
                                                               3 * staking.economic.genesis.economicModel.staking.stakeThreshold)
        result = staking.create_staking(0, address, address, node_id=staking.node.node_id,
                                        amount=2 * staking.economic.genesis.economicModel.staking.stakeThreshold,
                                        reward_per=0)
        assert_code(result, 0)
        self.assert_rewardsper(client_new_node, 0, 0)

    def test_IV_038_IV_039_IV_044_IV_040(self, client_new_node):
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

        result = staking.create_staking(0, address, address, node_id=staking.node.node_id,
                                        amount=2 * staking.economic.genesis.economicModel.staking.stakeThreshold,
                                        reward_per=10000)
        assert_code(result, 0)
        self.assert_rewardsper(client_new_node, 0, 10000)

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
        self.assert_rewardsper(client, 0, value)

    def test_IV_035(self, new_genesis_env, clients_consensus):
        new_genesis_env.deploy_all()
        client = clients_consensus[-1]
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
        client.economic.wait_settlement_blocknum(client.node, 2)
        address, _ = client.economic.account.generate_account(client.node.web3, 10**18 * 10000000)
        result = client.staking.create_staking(0, address, address, reward_per=value)
        assert_code(result, 0)
        self.assert_rewardsper(client, 0, value)

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

    def test_MPI_018_to_027(self, client_verifier):
        client = client_verifier
        value = client.staking.get_rewardper(client.node)
        newvalue = get_new_value(value)
        result = client.staking.edit_candidate(client.node.staking_address, client.node.staking_address,
                                               reward_per=newvalue)
        assert_code(result, 0)
        self.assert_rewardsper(client, value, newvalue)

        result = client.staking.edit_candidate(client.node.staking_address, client.node.staking_address,
                                               reward_per=0)
        assert_code(result, 0)
        self.assert_rewardsper(client, value, 0)

        result = client.staking.edit_candidate(client.node.staking_address, client.node.staking_address,
                                               reward_per=1)
        assert_code(result, 0)
        self.assert_rewardsper(client, value, 1)

        result = client.staking.edit_candidate(client.node.staking_address, client.node.staking_address,
                                               reward_per=10000)
        assert_code(result, 0)
        self.assert_rewardsper(client, value, 10000)
        result = client.staking.edit_candidate(client.node.staking_address, client.node.staking_address,
                                               reward_per=10001)
        log.info('Edit candidate information reward percent is 10001, result : {}'.format(result))
        assert_code(result, 301007)

        try:
            client.staking.edit_candidate(client.node.staking_address, client.node.staking_address,
                                               reward_per=-1)
        except TypeError as e:
            str(e) == 'Did not find sedes handling type int'
        result = client.staking.edit_candidate(client.node.staking_address, client.node.staking_address,
                                               reward_per='1')
        log.info('Edit candidate information reward percent is string, result : {}'.format(result))
        assert_code(result, 0)

        try:
            client.staking.edit_candidate(client.node.staking_address, client.node.staking_address,
                                               reward_per=10000000000000000000000000000000000000000000000)
        except ValueError as e:
            str(e) == "gas required exceeds allowance or always failing transaction"


    def test_MPI_034(self, client_new_node):
        client = client_new_node
        address, _ = client.economic.account.generate_account(client.node.web3, 10**18 * 10000000)
        result = client.staking.create_staking(0, address, address, reward_per=100)
        assert_code(result, 0)
        result = client.staking.edit_candidate(address, address, reward_per=101)
        assert_code(result, 0)
        assert_code(100, client.staking.get_rewardper())
        assert_code(101, client.staking.get_rewardper(isnext=True))
        vaule, newvalue = get_pledge_list(client.ppos.getCandidateList, nodeid=client.node.node_id)
        assert_code(vaule, 100)
        assert_code(newvalue, 101)

    def test_MPI_035(self, client_verifier):
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
        value_exiting, nextvalue_exiting = get_pledge_list(client.ppos.getCandidateList, client.node.node_id)
        assert_code(value, value_exiting)
        assert_code(nextvalue, nextvalue_exiting)

    def test_MPI_036(self, client_verifier):
        client = client_verifier
        value, nextvalue = get_pledge_list(client.ppos.getVerifierList, client.node.node_id)
        newvalue = get_new_value(value)
        result = client.staking.edit_candidate(client.node.staking_address, client.node.staking_address,
                                               reward_per=newvalue)
        assert_code(result, 0)
        self.assert_rewardsper(client, value, newvalue)
        result = client.staking.withdrew_staking(client.node.staking_address)
        assert_code(result, 0)
        assert_code(client.staking.get_rewardper(client.node), value)
        assert_code(client.staking.get_rewardper(client.node, isnext=True), newvalue)
        value_exiting, nextvalue_exiting = get_pledge_list(client.ppos.getCandidateList, client.node.node_id)
        assert_code(value, value_exiting)
        assert_code(newvalue, nextvalue_exiting)

class TestgetDelegateReward:



    def test(self, client_verifier, client_noconsensus):
        address, _ = client_verifier.economic.account.generate_account(client_verifier.node.web3, 10**18*10000)
        result = client_verifier.ppos.getDelegateReward(address, node_ids=[client_verifier.node.node_id])
        log.info(result)
        assert_code(result, 2)
        result = client_verifier.delegate.withdraw_delegate_reward(address)
        log.info(result)

    def test_IN_GR_001_IN_GR_002_IN_GR_003(self, clients_new_node):
        client0 = clients_new_node[0]
        client1 = clients_new_node[1]
        address0, _ = client0.economic.account.generate_account(client0.node.web3, 10**18 * 10000000)
        address1, _ = client0.economic.account.generate_account(client0.node.web3, 10**18 * 10000000)
        result = client0.staking.create_staking(0, address0, address0,
                                                amount=2 * client0.economic.genesis.economicModel.staking.stakeThreshold,
                                                reward_per=1111)
        assert_code(result, 0)
        staking_block0 = client0.staking.get_stakingblocknum(client0.node)
        result = client1.staking.create_staking(0, address1, address1,
                                                amount=2 * client1.economic.genesis.economicModel.staking.stakeThreshold,
                                                reward_per=1111)
        assert_code(result, 0)
        staking_block1 = client1.staking.get_stakingblocknum(client1.node)
        delegate_address, _ = client1.economic.account.generate_account(client1.node.web3, 10**18 * 100000)
        result = client0.delegate.delegate(0, delegate_address, amount=10**18 * 1000)
        assert_code(result, 0)
        time.sleep(3)
        result = client1.delegate.delegate(0, delegate_address, amount=10**18 * 1000)
        assert_code(result, 0)
        result = client1.ppos.getDelegateInfo(staking_block1, delegate_address, client1.node.node_id)
        log.info(result)
        result = client0.ppos.getDelegateInfo(staking_block0, delegate_address, client0.node.node_id)
        log.info(result)
        client0.economic.wait_settlement_blocknum(client0.node)
        verifier_list = get_pledge_list(client0.ppos.getVerifierList)
        assert client0.node.node_id in verifier_list
        assert client1.node.node_id in verifier_list
        result = client0.ppos.getDelegateReward(delegate_address)
        log.info('Do not given nodeid, get address delegate reward : {}'.format(result))
        result = client0.ppos.getDelegateReward(delegate_address, node_ids=[client0.node.node_id])
        log.info('Get address delegate nodeid {} reward : {}'.format(client0.node.node_id, result))

        result = client0.ppos.getDelegateReward(delegate_address, node_ids=[client1.node.node_id])
        log.info('Get address delegate nodeid {} reward : {}'.format(client1.node.node_id, result))

        result = client0.ppos.getDelegateReward(delegate_address, node_ids=[client0.node.node_id, client1.node.node_id])
        log.info('Get address delegate nodeid {},{} reward : {}'.format(client1.node.node_id,
                                                                        client1.node.node_id, result))

        client0.economic.wait_settlement_blocknum(client0.node)
        result = client0.ppos.getCandidateInfo(client0.node.node_id)
        log.info('nodeid {} candidate information : {}'.format(client0.node.node_id, result))
        result = client1.ppos.getCandidateInfo(client1.node.node_id)
        log.info('nodeid {} candidate information : {}'.format(client1.node.node_id, result))

        result = client1.ppos.getDelegateInfo(staking_block1, delegate_address, client1.node.node_id)
        log.info('nodeid {} delegate information : {}'.format(client1.node.node_id, result))
        result = client0.ppos.getDelegateInfo(staking_block0, delegate_address, client0.node.node_id)
        log.info('nodeid {} delegate information : {}'.format(client0.node.node_id, result))

        reward = client0.delegate.get_delegate_reward_by_nodeid(delegate_address)
        log.info('Address {} reward : {}'.format(delegate_address, reward))
        balance_before = client0.node.eth.getBalance(delegate_address)
        log.info('Before delegate address balance : {}'.format(balance_before))

        data = rlp.encode([rlp.encode(int(5000))])
        gas = get_the_dynamic_parameter_gas_fee(data) + 21000 + 3000 + 8000 + 2 * 1000 + 2 * 100
        log.info('Calculated gas : {}'.format(gas))
        result = client0.delegate.withdraw_delegate_reward(delegate_address, transaction_cfg=client0.pip.cfg.transaction_cfg)
        log.info(result)
        assert_code(result, 0)
        balance_after = client0.node.eth.getBalance(delegate_address)
        log.info('After delegate address balance : {}'.format(balance_after))
        assert balance_before + reward - gas * client0.pip.cfg.transaction_cfg.get('gasPrice') == balance_after


    def test_IN_GR_008(self, client_verifier):
        client = client_verifier
        address, _ = client.economic.account.generate_account(client.node.web3, 10**18 * 10000)
        try:
            client.ppos.getDelegateReward(address, node_ids=client.node.node_id)
        except ValueError as e:
            str(e) == 'non-hexadecimal number found in fromhex() arg at position 1'

    def test_IN_GR_009_IN_GR_010_IN_GR_012(self, client_new_node):
        client = client_new_node
        address, _ = client.economic.account.generate_account(client.node.web3, 10**18 * 100000)
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
        assert reward != 0
        client.staking.withdrew_staking(client.node.staking_address)
        rewardinfo = client.ppos.getDelegateReward(address)
        assert reward == rewardinfo.get('Ret')[0].get('reward')
        client.economic.wait_settlement_blocknum(client.node, client.economic.genesis.economicModel.staking.unStakeFreezeDuration)
        rewardinfo = client.ppos.getDelegateReward(address)
        assert reward == rewardinfo.get('Ret')[0].get('reward')

    def test_IN_GR_011(self, client_new_node, client_verifier):
        client = client_new_node
        address, _ = client.economic.account.generate_account(client.node.web3, 10**18 * 100000)
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
        assert reward != 0
        report_information = mock_duplicate_sign(1, client.node.nodekey,
                                                 client.node.blsprikey, client.node.block_number - 10)
        log.info("Report information: {}".format(report_information))
        address, _ = client_verifier.economic.account.generate_account(client_verifier.node.web3, 10 ** 18 * 1000)
        result = client_verifier.duplicatesign.reportDuplicateSign(1, report_information, address)
        log.info('Node duplicate block result : {}'.format(result))
        assert_code(result, 0)
        verifier_list = get_pledge_list(client_verifier.ppos.getVerifierList)
        assert client.node.node_id not in verifier_list
        rewardinfo = client.ppos.getDelegateReward(address)
        assert reward == rewardinfo.get('Ret')[0].get('reward')
        client.economic.wait_settlement_blocknum(client.node, client.economic.genesis.economicModel.staking.unStakeFreezeDuration)
        client.staking.withdrew_staking(client.node.staking_address)
        rewardinfo = client.ppos.getDelegateReward(address)
        assert reward == rewardinfo.get('Ret')[0].get('reward')

    def test_IN_GR_013(self, client_new_node):
        client = client_new_node
        address, _ = client.economic.account.generate_account(client.node.web3, 10**18 * 10000)
        result = client.ppos.getDelegateReward(address, [client.node.node_id])
        assert_code(result, 2)

    def test_IN_GR_013(self, client_new_node, client_consensus):
        client = client_new_node
        address, _ = client.economic.account.generate_account(client.node.web3, 10**18 * 10000)
        result = client.ppos.getDelegateReward(address, [client_consensus.node.node_id])
        assert_code(result, 2)

    def test_IN_GR_015_to_IN_GR_018(self, clients_new_node):
        client1 = clients_new_node[0]
        client2 = clients_new_node[1]
        address1, _ = client1.economic.account.generate_account(client1.node.web3, 10**18 * 10000)
        address2, _ = client1.economic.account.generate_account(client1.node.web3, 10**18 * 10000)
        address3, _ = client1.economic.account.generate_account(client1.node.web3, 10**18 * 10000)
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
        assert reward_info_address1_1.get('Ret')[0].get('reward') != 0
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
        assert_code(reward_info_address2_3, 2)
        reward_info_address2_4 = client2.ppos.getDelegateReward(address2, [client1.node.node_id, client2.node.node_id])
        log.info('Incoming nodeid {}，{}, address {}, reward information : {}'.format(client1.node.node_id,
            client2.node.node_id, address2, reward_info_address2_4
        ))
        assert_code(reward_info_address2_2, 0)
        assert_code(len(reward_info_address2_2.get('Ret')), 1)
        assert client2.delegate.get_staking_num_by_nodeid(address2, client2.node.node_id) == staking_num2
        assert reward_info_address2_2.get('Ret')[0].get('reward') != 0

        reward_info_address3_1 = client1.ppos.getDelegateReward(address3)
        assert_code(reward_info_address3_1, 2)
        reward_info_address3_2 = client1.ppos.getDelegateReward(address3, node_ids=[client1.node.node_id])
        assert_code(reward_info_address3_2, 2)
        reward_info_address3_3 = client1.ppos.getDelegateReward(address3, node_ids=[client2.node.node_id])
        assert_code(reward_info_address3_3, 2)
        address1_balance_before = client1.node.eth.getBalance(address1)
        log.info('Before getDelegateReward, the address {} balance: {}'.format(address1, address1_balance_before))
        reward_address1 = client1.delegate.get_delegate_reward_by_nodeid(address1, node_ids=[client1.node.node_id])
        log.info('Address {} delegate node {} reward : {}'.format(address1, client1.node.node_id, reward_address1))
        result = client1.delegate.withdrew_delegate(staking_num1, address1, node_id=client1.node.node_id,
                                                    amount=10**18 * 1000, transaction_cfg=client1.pip.cfg.transaction_cfg)
        assert_code(result, 0)
        reward_info_address1 = client1.ppos.getDelegateReward(address1, node_ids=[client1.node.node_id])
        assert_code(reward_info_address1, 2)
        address1_balance_after = client1.node.eth.getBalance(address1)
        log.info('Before getDelegateReward, the address {} balance: {}'.format(address1, address1_balance_after))
        data = rlp.encode([rlp.encode(int(1005)), rlp.encode(staking_num1), rlp.encode(bytes.fromhex(client1.node.node_id)),
                           rlp.encode(10**18 * 1000)])
        connt = get_the_dynamic_parameter_gas_fee(data)
        gas = (21000 + 6000 + 8000 + connt) * client1.pip.cfg.transaction_cfg.get('gasPrice')
        assert address1_balance_before + reward_address1 - gas + 10**18 * 1000 == address1_balance_after

class TestwithdrawDelegateReward():
    def test_IN_GR_020_IN_GR_021(self, client_new_node):
        client = client_new_node
        address, _ = client.economic.account.generate_account(client.node.web3, 10**18 * 100000)
        staking_and_delegate(client, address)
        result = client.ppos.getDelegateReward(address)
        log.info('Address {} delegate reward information : {}'.format(address, result))
        assert_code(result, 0)
        balance_before_withdraw = client.node.eth.getBalance(address)
        log.info('Address {} before withdraw reward balance : {}'.format(address, balance_before_withdraw))
        result = client.delegate.withdraw_delegate_reward(address, transaction_cfg=client.pip.cfg.transaction_cfg)
        log.info('Address {} withdraw delegate reward result : {}'.format(address, result))
        assert_code(result, 0)
        balance_after_withdraw = client.node.eth.getBalance(address)
        log.info('Address {} after withdraw reward balance : {}'.format(address, balance_after_withdraw))
        data = rlp.encode([rlp.encode(int(5000))])
        gas = (get_the_dynamic_parameter_gas_fee(data) + 8000 + 3000 + 21000 + 1000
               ) * client.pip.cfg.transaction_cfg.get('gasPrice')
        assert balance_before_withdraw - gas == balance_after_withdraw
        client.economic.wait_settlement_blocknum(client.node)
        blocknum = client.node.block_number
        log.info('Block number : {}'.format(blocknum))
        result = client.ppos.getDelegateReward(address)
        log.info('Address {} delegate reward information : {}'.format(address, result))
        assert_code(result, 0)
        assert result.get('Ret')[0].get('reward') == 0
        result = client.delegate.withdraw_delegate_reward(address, transaction_cfg=client.pip.cfg.transaction_cfg)
        log.info('Address {} withdraw delegate reward result : {}'.format(address, result))
        assert_code(result, 0)
        balance_after_withdraw2 = client.node.eth.getBalance(address)
        log.info('Address {} after withdraw reward balance : {}'.format(address, balance_after_withdraw2))
        gas = (get_the_dynamic_parameter_gas_fee(data) + 8000 + 3000 + 21000 + 1000
               ) * client.pip.cfg.transaction_cfg.get('gasPrice')
        assert balance_before_withdraw - gas == balance_after_withdraw
        wait_block_number()
        client.node.eth.generateGasPrice()

    def test(self, client_verifier):
        a = get_getDelegateReward_gas_fee(client_verifier, 1, 1)
        log.info(a)























