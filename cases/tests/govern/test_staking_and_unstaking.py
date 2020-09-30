from tests.lib.utils import upload_platon, assert_code, get_pledge_list, wait_block_number
from common.log import log
from tests.lib.client import Client, get_client_by_nodeid, get_clients_by_nodeid, StakingConfig
import pytest
import allure
import time
import math
from tests.govern.conftest import version_proposal_vote, proposal_vote
from tests.lib import Genesis, PipConfig
from dacite import from_dict
from tests.govern.test_voting_statistics import submitcvpandvote, submitcppandvote, submittpandvote, \
    submitvpandvote, submitppandvote
from common.key import mock_duplicate_sign


def create_lockup_plan(client):
    address, _ = client.pip.economic.account.generate_account(client.node.web3,
                                                              3 * client.economic.genesis.economicModel.staking.stakeThreshold)
    plan = [{'Epoch': 20, 'Amount': 2 * client.economic.genesis.economicModel.staking.stakeThreshold}]
    result = client.restricting.createRestrictingPlan(address, plan, address,
                                                      transaction_cfg=client.pip.cfg.transaction_cfg)
    log.info('CreateRestrictingPlan result : {}'.format(result))
    assert_code(result, 0)
    result = client.staking.create_staking(1, address, address,
                                           amount=int(1.8 * client.economic.genesis.economicModel.staking.stakeThreshold),
                                           transaction_cfg=client.pip.cfg.transaction_cfg)
    log.info('Create staking result : {}'.format(result))
    assert_code(result, 0)
    client.economic.wait_settlement_blocknum(client.node)


@pytest.fixture()
def new_node_no_proposal(no_vp_proposal, clients_noconsensus, all_clients):
    pip = no_vp_proposal
    client = get_client_by_nodeid(pip.node.node_id, all_clients)
    candidate_list = get_pledge_list(client.ppos.getCandidateList)
    log.info('candidate_list: {}'.format(candidate_list))
    for client in clients_noconsensus:
        if client.node.node_id not in candidate_list:
            return client.pip
    log.info('All nodes are staked, restart the chain')
    pip.economic.env.deploy_all()
    return clients_noconsensus[0].pip


def replace_platon_and_staking(pip, platon_bin):
    all_nodes = pip.economic.env.get_all_nodes()
    all_clients = []
    for node in all_nodes:
        all_clients.append(Client(pip.economic.env, node, StakingConfig("externalId", "nodeName", "website",
                                                                        "details")))
    client = get_client_by_nodeid(pip.node.node_id, all_clients)
    upload_platon(pip.node, platon_bin)
    log.info('Replace the platon of the node {}'.format(pip.node.node_id))
    pip.node.restart()
    log.info('Restart the node {}'.format(pip.node.node_id))
    address, _ = pip.economic.account.generate_account(pip.node.web3,
                                                       10 * pip.economic.genesis.economicModel.staking.stakeThreshold)
    result = client.staking.create_staking(0, address, address, transaction_cfg=pip.cfg.transaction_cfg)
    log.info('Node {} staking result {}'.format(pip.node.node_id, result))
    return result


class TestVotingProposalStaking:
    @pytest.mark.P1
    @allure.title('Verify stake function')
    @pytest.mark.parametrize('platon_bin', [getattr(PipConfig(), 'PLATON_NEW_BIN2'),
                                            getattr(PipConfig(), 'PLATON_NEW_BIN1'),
                                            getattr(PipConfig(), 'PLATON_NEW_BIN0'),
                                            getattr(PipConfig(), 'PLATON_NEW_BIN3'),
                                            getattr(PipConfig(), 'PLATON_NEW_BIN')])
    def test_ST_VS_001_to_005(self, new_node_has_proposal, platon_bin):
        pip = new_node_has_proposal
        result = replace_platon_and_staking(pip, platon_bin)
        if platon_bin != pip.cfg.PLATON_NEW_BIN1:
            assert_code(result, 0)
        else:
            assert_code(result, 301004)


class TestNoProposalStaking:
    @pytest.mark.P1
    @allure.title('No proposal, verify stake function')
    def test_ST_NO_001(self, new_node_no_proposal):
        pip = new_node_no_proposal
        result = replace_platon_and_staking(pip, pip.cfg.PLATON_NEW_BIN2)
        assert_code(result, 0)

    @pytest.mark.P1
    @allure.title('No proposal, verify stake function')
    def test_ST_NO_002(self, new_node_no_proposal):
        pip = new_node_no_proposal
        result = replace_platon_and_staking(pip, pip.cfg.PLATON_NEW_BIN1)
        assert_code(result, 301004)

    @pytest.mark.P1
    @allure.title('No proposal, verify stake function')
    def test_ST_NO_003(self, new_node_no_proposal):
        pip = new_node_no_proposal
        result = replace_platon_and_staking(pip, pip.cfg.PLATON_NEW_BIN0)
        assert_code(result, 0)

    @pytest.mark.P1
    def test_ST_NO_004(self, new_node_no_proposal):
        pip = new_node_no_proposal
        result = replace_platon_and_staking(pip, pip.cfg.PLATON_NEW_BIN3)
        assert_code(result, 0)

    @pytest.mark.P1
    @allure.title('No proposal, verify stake function')
    def test_ST_NO_005(self, new_node_no_proposal):
        pip = new_node_no_proposal
        result = replace_platon_and_staking(pip, pip.cfg.PLATON_NEW_BIN)
        assert_code(result, 301005)


class TestPreactiveProposalStaking:
    def preactive_proposal(self, all_clients):
        verifier_list = get_pledge_list(all_clients[0].ppos.getVerifierList)
        log.info('verifierlist :{}'.format(verifier_list))
        client_verifiers = get_clients_by_nodeid(verifier_list, all_clients)
        pips = [client.pip for client in client_verifiers]
        result = pips[0].submitVersion(pips[0].node.node_id, str(time.time()),
                                       pips[0].cfg.version5, 2, pips[0].node.staking_address,
                                       transaction_cfg=pips[0].cfg.transaction_cfg)
        log.info('submit version proposal, result : {}'.format(result))
        proposalinfo = pips[0].get_effect_proposal_info_of_vote()
        log.info('Version proposalinfo: {}'.format(proposalinfo))
        for pip in pips:
            result = version_proposal_vote(pip)
            assert_code(result, 0)
        wait_block_number(pip.node, proposalinfo.get('EndVotingBlock'))
        assert pip.get_status_of_proposal(proposalinfo.get('ProposalID')) == 4

    @pytest.mark.P1
    @allure.title('There is preactive proposal, verify stake function')
    def test_ST_PR_001(self, new_genesis_env, new_node_no_proposal, all_clients):
        pip = new_node_no_proposal
        self.preactive_proposal(all_clients)
        result = replace_platon_and_staking(pip, pip.cfg.PLATON_NEW_BIN2)
        assert_code(result, 301004)

    @pytest.mark.P1
    @allure.title('There is preactive proposal, verify stake function')
    def test_ST_PR_002(self, new_genesis_env, new_node_no_proposal, all_clients):
        pip = new_node_no_proposal
        self.preactive_proposal(all_clients)
        result = replace_platon_and_staking(pip, pip.cfg.PLATON_NEW_BIN1)
        assert_code(result, 301004)

    @pytest.mark.P1
    @allure.title('There is preactive proposal, verify stake function')
    def test_ST_PR_003(self, new_genesis_env, new_node_no_proposal, all_clients):
        pip = new_node_no_proposal
        self.preactive_proposal(all_clients)
        result = replace_platon_and_staking(pip, pip.cfg.PLATON_NEW_BIN0)
        assert_code(result, 301004)

    @pytest.mark.P1
    @allure.title('There is preactive proposal, verify stake function')
    def test_ST_PR_004(self, new_genesis_env, new_node_no_proposal, all_clients):
        pip = new_node_no_proposal
        self.preactive_proposal(all_clients)
        result = replace_platon_and_staking(pip, pip.cfg.PLATON_NEW_BIN3)
        assert_code(result, 301004)

    @pytest.mark.P1
    @allure.title('There is preactive proposal, verify stake function')
    def test_ST_PR_005(self, new_genesis_env, new_node_no_proposal, all_clients):
        pip = new_node_no_proposal
        self.preactive_proposal(all_clients)
        result = replace_platon_and_staking(pip, pip.cfg.PLATON_NEW_BIN)
        assert_code(result, 0)

    @pytest.mark.P1
    @allure.title('There is preactive proposal, verify stake function')
    def test_ST_PR_006(self, new_genesis_env, new_node_no_proposal, all_clients):
        pip = new_node_no_proposal
        self.preactive_proposal(all_clients)
        result = replace_platon_and_staking(pip, pip.cfg.PLATON_NEW_BIN8)
        assert_code(result, 301005)


class TestUpgradedProposalStaking:
    def upgraded_proposal(self, all_clients):
        verifier_list = get_pledge_list(all_clients[0].ppos.getVerifierList)
        log.info('verifierlist :{}'.format(verifier_list))
        client_verifiers = get_clients_by_nodeid(verifier_list, all_clients)
        pips = [client.pip for client in client_verifiers]
        result = pips[0].submitVersion(pips[0].node.node_id, str(time.time()),
                                       pips[0].cfg.version5, 2, pips[0].node.staking_address,
                                       transaction_cfg=pips[0].cfg.transaction_cfg)
        log.info('submit version proposal, result : {}'.format(result))
        proposalinfo = pips[0].get_effect_proposal_info_of_vote()
        log.info('Version proposalinfo: {}'.format(proposalinfo))
        for pip in pips:
            result = version_proposal_vote(pip)
            assert_code(result, 0)
        wait_block_number(pip.node, proposalinfo.get('EndVotingBlock'))
        assert pip.get_status_of_proposal(proposalinfo.get('ProposalID')) == 4
        wait_block_number(pip.node, proposalinfo.get('ActiveBlock'))
        assert pip.get_status_of_proposal(proposalinfo.get('ProposalID')) == 5

    @pytest.mark.P2
    @allure.title('Chain upgrade completed, verify stake function')
    def test_ST_UPG_001(self, new_genesis_env, new_node_no_proposal, all_clients):
        pip = new_node_no_proposal
        self.upgraded_proposal(all_clients)
        result = replace_platon_and_staking(pip, pip.cfg.PLATON_NEW_BIN4)
        assert_code(result, 0)

    @pytest.mark.P2
    @allure.title('Chain upgrade completed, verify stake function')
    def test_ST_UPG_002(self, new_genesis_env, new_node_no_proposal, all_clients):
        pip = new_node_no_proposal
        self.upgraded_proposal(all_clients)
        result = replace_platon_and_staking(pip, pip.cfg.PLATON_NEW_BIN0)
        assert_code(result, 301004)

    @pytest.mark.P2
    @allure.title('Chain upgrade completed, verify stake function')
    def test_ST_UPG_003(self, new_genesis_env, new_node_no_proposal, all_clients):
        pip = new_node_no_proposal
        self.upgraded_proposal(all_clients)
        result = replace_platon_and_staking(pip, pip.cfg.PLATON_NEW_BIN)
        assert_code(result, 0)

    @pytest.mark.P2
    @allure.title('Chain upgrade completed, verify stake function')
    def test_ST_UPG_004(self, new_genesis_env, new_node_no_proposal, all_clients):
        pip = new_node_no_proposal
        self.upgraded_proposal(all_clients)
        result = replace_platon_and_staking(pip, pip.cfg.PLATON_NEW_BIN6)
        assert_code(result, 0)

    @pytest.mark.P2
    @allure.title('Chain upgrade completed, verify stake function')
    def test_ST_UPG_005(self, new_genesis_env, new_node_no_proposal, all_clients):
        pip = new_node_no_proposal
        self.upgraded_proposal(all_clients)
        result = replace_platon_and_staking(pip, pip.cfg.PLATON_NEW_BIN7)
        assert_code(result, 301005)


class TestUnstaking:
    @pytest.mark.P1
    @allure.title('Verify unstake function')
    def test_UNS_AM_003_007(self, new_genesis_env, client_verifier):
        genesis = from_dict(data_class=Genesis, data=new_genesis_env.genesis_config)
        genesis.economicModel.staking.unStakeFreezeDuration = 2
        genesis.economicModel.slashing.maxEvidenceAge = 1
        new_genesis_env.set_genesis(genesis.to_dict())
        new_genesis_env.deploy_all()
        pip = client_verifier.pip
        address = pip.node.staking_address
        submitcvpandvote([client_verifier], 1)
        result = version_proposal_vote(pip)
        assert_code(result, 0)
        shares = client_verifier.staking.get_staking_amount(pip.node)
        result = client_verifier.staking.withdrew_staking(address)
        log.info('Node withdrew staking result : {}'.format(result))
        assert_code(result, 0)
        calculated_block = 480
        wait_block_number(pip.node, calculated_block)
        balance_before = pip.node.eth.getBalance(address, calculated_block - 1)
        log.info('Block bumber {} staking address balance {}'.format(calculated_block - 1, balance_before))
        balance_after = pip.node.eth.getBalance(address, calculated_block)
        log.info('Block bumber {} staking address balance {}'.format(calculated_block, balance_after))
        log.info('{}'.format(pip.economic.get_current_year_reward(pip.node)))
        assert balance_after - balance_before == shares

    @pytest.mark.P1
    @allure.title('Verify unstake function')
    def test_UNS_AM_005(self, new_genesis_env, client_verifier):
        genesis = from_dict(data_class=Genesis, data=new_genesis_env.genesis_config)
        genesis.economicModel.staking.unStakeFreezeDuration = 2
        genesis.economicModel.slashing.maxEvidenceAge = 1
        new_genesis_env.set_genesis(genesis.to_dict())
        new_genesis_env.deploy_all()
        pip = client_verifier.pip
        address = pip.node.staking_address
        result = pip.submitText(pip.node.node_id, str(time.time()), address, transaction_cfg=pip.cfg.transaction_cfg)
        assert_code(result, 0)

        submitcppandvote([client_verifier], [1])
        result = proposal_vote(pip, proposaltype=pip.cfg.text_proposal)
        assert_code(result, 0)
        shares = client_verifier.staking.get_staking_amount(pip.node)
        result = client_verifier.staking.withdrew_staking(address)
        log.info('Node withdrew staking result : {}'.format(result))
        assert_code(result, 0)
        calculated_block = 480
        wait_block_number(pip.node, calculated_block)
        balance_before = pip.node.eth.getBalance(address, calculated_block - 1)
        log.info('Block bumber {} staking address balance {}'.format(calculated_block - 1, balance_before))
        balance_after = pip.node.eth.getBalance(address, calculated_block)
        log.info('Block bumber {} staking address balance {}'.format(calculated_block, balance_after))
        assert balance_after - balance_before == shares

    @pytest.mark.P1
    @allure.title('Verify unstake function')
    def test_UNS_AM_004_006_008(self, new_genesis_env, clients_verifier):
        genesis = from_dict(data_class=Genesis, data=new_genesis_env.genesis_config)
        genesis.economicModel.staking.unStakeFreezeDuration = 2
        genesis.economicModel.slashing.maxEvidenceAge = 1
        genesis.economicModel.gov.versionProposalVoteDurationSeconds = 1000
        genesis.economicModel.gov.textProposalVoteDurationSeconds = 840
        new_genesis_env.set_genesis(genesis.to_dict())
        new_genesis_env.deploy_all()
        pip_one = clients_verifier[0].pip
        pip_two = clients_verifier[1].pip
        pip_three = clients_verifier[2].pip
        address = pip_one.node.staking_address
        result = pip_one.submitVersion(pip_one.node.node_id, str(time.time()), pip_one.cfg.version5, 17, address,
                                       transaction_cfg=pip_one.cfg.transaction_cfg)
        log.info('Submit version proposal result : {}'.format(result))
        proposalinfo_version = pip_one.get_effect_proposal_info_of_vote(pip_one.cfg.version_proposal)
        log.info('Version proposal information : {}'.format(proposalinfo_version))

        result = pip_one.submitCancel(pip_one.node.node_id, str(time.time()), 13, proposalinfo_version.get('ProposalID'),
                                      address, transaction_cfg=pip_one.cfg.transaction_cfg)
        log.info('Submit cancel proposal result : {}'.format(result))
        assert_code(result, 0)
        result_text = pip_one.submitText(pip_one.node.node_id, str(time.time()), address,
                                         transaction_cfg=pip_one.cfg.transaction_cfg)
        log.info('Submit text proposal result : {}'.format(result_text))
        result = proposal_vote(pip_one, proposaltype=pip_one.cfg.cancel_proposal)
        assert_code(result, 0)
        result = version_proposal_vote(pip_two)
        assert_code(result, 0)
        result = proposal_vote(pip_three, proposaltype=pip_three.cfg.text_proposal)
        assert_code(result, 0)
        proposalinfo_cancel = pip_one.get_effect_proposal_info_of_vote(pip_one.cfg.cancel_proposal)
        log.info('Get cancel proposal information {}'.format(proposalinfo_cancel))
        proposalinfo_text = pip_one.get_effect_proposal_info_of_vote(pip_one.cfg.text_proposal)
        log.info('Get cancel proposal information {}'.format(proposalinfo_text))
        shares_one = clients_verifier[0].staking.get_staking_amount(pip_one.node)
        shares_two = clients_verifier[1].staking.get_staking_amount(pip_two.node)
        shares_three = clients_verifier[2].staking.get_staking_amount(pip_three.node)
        result = clients_verifier[0].staking.withdrew_staking(address)
        log.info('Node {} withdrew staking result : {}'.format(pip_one.node.node_id, result))
        assert_code(result, 0)

        result = clients_verifier[1].staking.withdrew_staking(address)
        log.info('Node {} withdrew staking result : {}'.format(pip_two.node.node_id, result))
        assert_code(result, 0)

        result = clients_verifier[2].staking.withdrew_staking(address)
        log.info('Node {} withdrew staking result : {}'.format(pip_three.node.node_id, result))
        assert_code(result, 0)
        calculated_block = 480
        wait_block_number(pip_one.node, calculated_block)
        balance_before = pip_one.node.eth.getBalance(address, calculated_block - 1)
        log.info('Block bumber {} staking address balance {}'.format(calculated_block - 1, balance_before))
        balance_after = pip_one.node.eth.getBalance(address, calculated_block)
        log.info('Block bumber {} staking address balance {}'.format(calculated_block, balance_after))
        assert balance_after == balance_before

        blocknumber = math.ceil(proposalinfo_cancel.get('EndVotingBlock') / pip_one.economic.settlement_size
                                ) * pip_one.economic.settlement_size
        wait_block_number(pip_one.node, blocknumber)
        balance_before = pip_one.node.eth.getBalance(address, blocknumber - 1)
        log.info('Block bumber {} staking address balance {}'.format(blocknumber - 1, balance_before))
        balance_after = pip_one.node.eth.getBalance(address, blocknumber)
        log.info('Block bumber {} staking address balance {}'.format(blocknumber, balance_after))
        assert balance_after - balance_before == shares_one

        blocknumber = math.ceil(proposalinfo_version.get('EndVotingBlock') / pip_one.economic.settlement_size
                                ) * pip_one.economic.settlement_size
        wait_block_number(pip_one.node, blocknumber)
        balance_before = pip_one.node.eth.getBalance(address, blocknumber - 1)
        log.info('Block bumber {} staking address balance {}'.format(blocknumber - 1, balance_before))
        balance_after = pip_one.node.eth.getBalance(address, blocknumber)
        log.info('Block bumber {} staking address balance {}'.format(blocknumber, balance_after))
        assert balance_after - balance_before == shares_two

        blocknumber = math.ceil(proposalinfo_text.get('EndVotingBlock') / pip_one.economic.settlement_size
                                ) * pip_one.economic.settlement_size
        wait_block_number(pip_one.node, blocknumber)
        balance_before = pip_one.node.eth.getBalance(address, blocknumber - 1)
        log.info('Block bumber {} staking address balance {}'.format(blocknumber - 1, balance_before))
        balance_after = pip_one.node.eth.getBalance(address, blocknumber)
        log.info('Block bumber {} staking address balance {}'.format(blocknumber, balance_after))
        assert balance_after - balance_before == shares_three

    @pytest.mark.P2
    @allure.title('Verify unstake function')
    def test_UNS_AM_009_011_013(self, new_genesis_env, clients_noconsensus):
        genesis = from_dict(data_class=Genesis, data=new_genesis_env.genesis_config)
        genesis.economicModel.staking.unStakeFreezeDuration = 2
        genesis.economicModel.slashing.maxEvidenceAge = 1
        genesis.economicModel.gov.versionProposalVoteDurationSeconds = 1000
        genesis.economicModel.gov.textProposalVoteDurationSeconds = 200
        new_genesis_env.set_genesis(genesis.to_dict())
        new_genesis_env.deploy_all()
        pip_test = clients_noconsensus[0].pip
        address, _ = pip_test.economic.account.generate_account(pip_test.node.web3, 10**18 * 20000000)
        plan = [{'Epoch': 20, 'Amount': 10**18 * 2000000}]
        result = clients_noconsensus[0].restricting.createRestrictingPlan(address, plan, address,
                                                                          transaction_cfg=pip_test.cfg.transaction_cfg)
        log.info('CreateRestrictingPlan result : {}'.format(result))
        assert_code(result, 0)
        result = clients_noconsensus[0].staking.create_staking(1, address, address,
                                                               transaction_cfg=pip_test.cfg.transaction_cfg)
        log.info('Create staking result : {}'.format(result))
        assert_code(result, 0)
        pip_test.economic.wait_settlement_blocknum(pip_test.node)
        result = pip_test.submitVersion(pip_test.node.node_id, str(time.time()), pip_test.cfg.version5,
                                        4, pip_test.node.staking_address, transaction_cfg=pip_test.cfg.transaction_cfg)
        log.info('Submit version proposal result : {}'.format(result))
        assert_code(result, 0)
        result = pip_test.submitText(pip_test.node.node_id, str(time.time()), pip_test.node.staking_address,
                                     transaction_cfg=pip_test.cfg.transaction_cfg)
        log.info('Submit text proposal result : {}'.format(result))
        assert_code(result, 0)

        proposalinfo_version = pip_test.get_effect_proposal_info_of_vote()
        log.info('Get version proposal information : {}'.format(proposalinfo_version))
        result = pip_test.submitCancel(pip_test.node.node_id, str(time.time()), 2, proposalinfo_version.get('ProposalID'),
                                       pip_test.node.staking_address, transaction_cfg=pip_test.cfg.transaction_cfg)
        log.info('Submit cancel result : {}'.format(result))
        assert_code(result, 0)
        result = version_proposal_vote(pip_test)
        assert_code(result, 0)
        result = proposal_vote(pip_test, proposaltype=pip_test.cfg.text_proposal)
        assert_code(result, 0)
        result = proposal_vote(pip_test, proposaltype=pip_test.cfg.cancel_proposal)
        assert_code(result, 0)
        proposalinfo_cancel = pip_test.get_effect_proposal_info_of_vote(pip_test.cfg.cancel_proposal)
        log.info('Get cancel proposal information : {}'.format(proposalinfo_cancel))
        proposalinfo_text = pip_test.get_effect_proposal_info_of_vote(pip_test.cfg.text_proposal)
        log.info('Get text proposal information : {}'.format(proposalinfo_text))

        shares = clients_noconsensus[0].staking.get_staking_amount(pip_test.node)
        result = clients_noconsensus[0].staking.withdrew_staking(address)
        log.info('Node withdrew result : {}'.format(result))
        assert_code(result, 0)

        wait_block_number(pip_test.node, 4 * pip_test.economic.settlement_size)
        balance_before = pip_test.node.eth.getBalance(address, 4 * pip_test.economic.settlement_size - 1)
        log.info('Block bumber {} staking address balance {}'.format(4 * pip_test.economic.settlement_size - 1, balance_before))
        balance_before_lockup = pip_test.node.eth.getBalance(pip_test.cfg.FOUNDATION_LOCKUP_ADDRESS,
                                                             4 * pip_test.economic.settlement_size - 1)
        log.info('Block bumber {} FOUNDATION_LOCKUP_ADDRESS balance {}'.format(4 * pip_test.economic.settlement_size - 1,
                                                                               balance_before_lockup))
        balance_after = pip_test.node.eth.getBalance(address, 4 * pip_test.economic.settlement_size)
        balance_after_lockup = pip_test.node.eth.getBalance(pip_test.cfg.FOUNDATION_LOCKUP_ADDRESS,
                                                            4 * pip_test.economic.settlement_size)
        log.info('Block bumber {} staking address balance {}'.format(4 * pip_test.economic.settlement_size, balance_after))
        log.info('Block bumber {} FOUNDATION_LOCKUP_ADDRESS balance {}'.format(4 * pip_test.economic.settlement_size,
                                                                               balance_after_lockup))
        assert balance_after_lockup - balance_before_lockup == shares

    @pytest.mark.P2
    @allure.title('Verify unstake function')
    def test_UNS_AM_010_012_014(self, new_genesis_env, clients_noconsensus):
        genesis = from_dict(data_class=Genesis, data=new_genesis_env.genesis_config)
        genesis.economicModel.staking.unStakeFreezeDuration = 2
        genesis.economicModel.slashing.maxEvidenceAge = 1
        genesis.economicModel.gov.paramProposalVoteDurationSeconds = 640
        genesis.economicModel.gov.textProposalVoteDurationSeconds = 840
        new_genesis_env.set_genesis(genesis.to_dict())
        new_genesis_env.deploy_all()
        for client in clients_noconsensus:
            pip = client.pip
            address, _ = pip.economic.account.generate_account(pip.node.web3, 10**18 * 20000000)
            plan = [{'Epoch': 20, 'Amount': 10**18 * 2000000}]
            result = client.restricting.createRestrictingPlan(address, plan, address,
                                                              transaction_cfg=pip.cfg.transaction_cfg)
            log.info('CreateRestrictingPlan result : {}'.format(result))
            assert_code(result, 0)
            result = client.staking.create_staking(1, address, address, amount=10**18 * 1800000,
                                                   transaction_cfg=pip.cfg.transaction_cfg)
            log.info('Create staking result : {}'.format(result))
            assert_code(result, 0)
        pip.economic.wait_settlement_blocknum(pip.node)
        result = pip.submitParam(pip.node.node_id, str(time.time()), 'slashing', 'slashBlocksReward',
                                 '1116', address, transaction_cfg=pip.cfg.transaction_cfg)
        log.info('Submit param proposal result : {}'.format(result))
        assert_code(result, 0)
        result = pip.submitText(pip.node.node_id, str(time.time()), pip.node.staking_address,
                                transaction_cfg=pip.cfg.transaction_cfg)
        log.info('Submit text proposal result : {}'.format(result))
        assert_code(result, 0)

        proposalinfo_param = pip.get_effect_proposal_info_of_vote(pip.cfg.param_proposal)
        log.info('Get param proposal information : {}'.format(proposalinfo_param))
        result = pip.submitCancel(pip.node.node_id, str(time.time()), 14, proposalinfo_param.get('ProposalID'),
                                  pip.node.staking_address, transaction_cfg=pip.cfg.transaction_cfg)
        log.info('Submit cancel result : {}'.format(result))
        assert_code(result, 0)
        result = proposal_vote(clients_noconsensus[0].pip, proposaltype=pip.cfg.param_proposal)
        assert_code(result, 0)
        result = proposal_vote(clients_noconsensus[1].pip, proposaltype=pip.cfg.text_proposal)
        assert_code(result, 0)
        result = proposal_vote(clients_noconsensus[2].pip, proposaltype=pip.cfg.cancel_proposal)
        assert_code(result, 0)
        proposalinfo_cancel = pip.get_effect_proposal_info_of_vote(pip.cfg.cancel_proposal)
        log.info('Get cancel proposal information : {}'.format(proposalinfo_cancel))
        proposalinfo_text = pip.get_effect_proposal_info_of_vote(pip.cfg.text_proposal)
        log.info('Get text proposal information : {}'.format(proposalinfo_text))

        shares0 = clients_noconsensus[0].staking.get_staking_amount(clients_noconsensus[0].node)
        shares1 = clients_noconsensus[1].staking.get_staking_amount(clients_noconsensus[1].node)
        shares2 = clients_noconsensus[2].staking.get_staking_amount(clients_noconsensus[2].node)
        address0 = clients_noconsensus[0].node.staking_address
        address1 = clients_noconsensus[1].node.staking_address
        address2 = clients_noconsensus[2].node.staking_address
        result = clients_noconsensus[0].staking.withdrew_staking(address0)
        log.info('Node {} withdrew result : {}'.format(clients_noconsensus[0].node.node_id, result))
        assert_code(result, 0)
        result = clients_noconsensus[1].staking.withdrew_staking(address1)
        log.info('Node {} withdrew result : {}'.format(clients_noconsensus[0].node.node_id, result))
        assert_code(result, 0)
        result = clients_noconsensus[2].staking.withdrew_staking(address2)
        log.info('Node {} withdrew result : {}'.format(clients_noconsensus[0].node.node_id, result))
        assert_code(result, 0)
        wait_block_number(pip.node, 4 * pip.economic.settlement_size)
        balance_before_lockup = pip.node.eth.getBalance(pip.cfg.FOUNDATION_LOCKUP_ADDRESS,
                                                        4 * pip.economic.settlement_size - 1)
        log.info('Block bumber {} FOUNDATION_LOCKUP_ADDRESS balance {}'.format(4 * pip.economic.settlement_size - 1,
                                                                               balance_before_lockup))
        balance_after_lockup = pip.node.eth.getBalance(pip.cfg.FOUNDATION_LOCKUP_ADDRESS,
                                                       4 * pip.economic.settlement_size)
        log.info('Block bumber {} FOUNDATION_LOCKUP_ADDRESS balance {}'.format(4 * pip.economic.settlement_size,
                                                                               balance_after_lockup))
        assert balance_after_lockup == balance_before_lockup

        wait_block_number(pip.node, 5 * pip.economic.settlement_size)
        balance_before = pip.node.eth.getBalance(address2, 5 * pip.economic.settlement_size - 1)
        balance_before_lockup = pip.node.eth.getBalance(pip.cfg.FOUNDATION_LOCKUP_ADDRESS,
                                                        5 * pip.economic.settlement_size - 1)
        log.info('Block bumber {} staking address balance {}'.format(5 * pip.economic.settlement_size - 1, balance_before))
        log.info('Block bumber {} FOUNDATION_LOCKUP_ADDRESS balance {}'.format(5 * pip.economic.settlement_size - 1,
                                                                               balance_before_lockup))
        balance_after = pip.node.eth.getBalance(address2, 5 * pip.economic.settlement_size)
        balance_after_lockup = pip.node.eth.getBalance(pip.cfg.FOUNDATION_LOCKUP_ADDRESS,
                                                       5 * pip.economic.settlement_size)
        log.info('Block bumber {} staking address balance {}'.format(5 * pip.economic.settlement_size, balance_after))
        log.info('Block bumber {} FOUNDATION_LOCKUP_ADDRESS balance {}'.format(5 * pip.economic.settlement_size,
                                                                               balance_after_lockup))
        assert balance_after == balance_before
        assert balance_after_lockup - balance_before_lockup == shares2

        wait_block_number(pip.node, 6 * pip.economic.settlement_size)
        balance_before = pip.node.eth.getBalance(address0, 6 * pip.economic.settlement_size - 1)
        balance_before_lockup = pip.node.eth.getBalance(pip.cfg.FOUNDATION_LOCKUP_ADDRESS,
                                                        6 * pip.economic.settlement_size - 1)
        log.info('Block bumber {} staking address balance {}'.format(6 * pip.economic.settlement_size - 1, balance_before))
        log.info('Block bumber {} FOUNDATION_LOCKUP_ADDRESS balance {}'.format(6 * pip.economic.settlement_size - 1,
                                                                               balance_before_lockup))
        balance_after = pip.node.eth.getBalance(address0, 6 * pip.economic.settlement_size)
        balance_after_lockup = pip.node.eth.getBalance(pip.cfg.FOUNDATION_LOCKUP_ADDRESS,
                                                       6 * pip.economic.settlement_size)
        log.info('Block bumber {} staking address balance {}'.format(6 * pip.economic.settlement_size, balance_after))
        log.info('Block bumber {} FOUNDATION_LOCKUP_ADDRESS balance {}'.format(6 * pip.economic.settlement_size, balance_after_lockup))

        assert balance_after == balance_before
        assert balance_after_lockup - balance_before_lockup == shares0

        wait_block_number(pip.node, 7 * pip.economic.settlement_size)
        balance_before = pip.node.eth.getBalance(address1, 7 * pip.economic.settlement_size - 1)
        balance_before_lockup = pip.node.eth.getBalance(pip.cfg.FOUNDATION_LOCKUP_ADDRESS,
                                                        7 * pip.economic.settlement_size - 1)
        log.info('Block bumber {} staking address balance {}'.format(7 * pip.economic.settlement_size - 1,
                                                                     balance_before))
        log.info('Block bumber {} FOUNDATION_LOCKUP_ADDRESS balance {}'.format(7 * pip.economic.settlement_size - 1,
                                                                               balance_before_lockup))
        balance_after = pip.node.eth.getBalance(address1, 7 * pip.economic.settlement_size)
        balance_after_lockup = pip.node.eth.getBalance(pip.cfg.FOUNDATION_LOCKUP_ADDRESS,
                                                       7 * pip.economic.settlement_size)
        log.info('Block bumber {} staking address balance {}'.format(7 * pip.economic.settlement_size, balance_after))
        log.info('Block bumber {} FOUNDATION_LOCKUP_ADDRESS balance {}'.format(7 * pip.economic.settlement_size,
                                                                               balance_after_lockup))

        assert balance_after == balance_before
        assert balance_after_lockup - balance_before_lockup == shares1


class TestSlashing:
    @pytest.mark.P1
    @allure.title('Node be slashed, verify unstake function')
    def test_UNS_PU_003_005_007_017(self, new_genesis_env, clients_consensus):
        genesis = from_dict(data_class=Genesis, data=new_genesis_env.genesis_config)
        genesis.economicModel.staking.unStakeFreezeDuration = 2
        genesis.economicModel.slashing.maxEvidenceAge = 1
        genesis.economicModel.gov.paramProposalVoteDurationSeconds = 200
        genesis.economicModel.gov.textProposalVoteDurationSeconds = 200
        genesis.economicModel.slashing.slashBlocksReward = 0
        new_genesis_env.set_genesis(genesis.to_dict())
        new_genesis_env.deploy_all()
        pip = clients_consensus[0].pip
        pip_test = clients_consensus[1].pip
        address = pip.node.staking_address
        result = pip.submitVersion(pip.node.node_id, str(time.time()), pip.cfg.version5, 1,
                                   address, transaction_cfg=pip.cfg.transaction_cfg)
        log.info('Submit version proposal result : {}'.format(result))
        assert_code(result, 0)
        proposalinfo_version = pip.get_effect_proposal_info_of_vote()
        log.info('Get version proposal information : {}'.format(proposalinfo_version))
        result = version_proposal_vote(pip)
        assert_code(result, 0)
        wait_block_number(pip.node, proposalinfo_version.get('EndVotingBlock'))
        submittpandvote([clients_consensus[0]], 3)
        submitcppandvote([clients_consensus[0]], [2])
        result = proposal_vote(pip, proposaltype=pip.cfg.param_proposal)
        assert_code(result, 0)
        log.info('Stop the node {}'.format(pip.node.node_id))
        shares = clients_consensus[1].staking.get_staking_amount(pip_test.node)
        pip.node.stop()
        wait_block_number(pip_test.node, 3 * pip_test.economic.settlement_size)
        balance_before = pip_test.node.eth.getBalance(address, 3 * pip_test.economic.settlement_size - 1)
        log.info('Block bumber {} staking address balance {}'.format(3 * pip_test.economic.settlement_size - 1,
                                                                     balance_before))
        balance_after = pip_test.node.eth.getBalance(address, 3 * pip_test.economic.settlement_size)

        log.info('Block bumber {} staking address balance {}'.format(3 * pip_test.economic.settlement_size,
                                                                     balance_after))
        assert balance_after - balance_before == shares

    @pytest.mark.P2
    @allure.title('Node be slashed, verify unstake function')
    def test_UNS_PU_016(self, new_genesis_env, clients_consensus):
        genesis = from_dict(data_class=Genesis, data=new_genesis_env.genesis_config)
        genesis.economicModel.staking.unStakeFreezeDuration = 2
        genesis.economicModel.slashing.maxEvidenceAge = 1
        genesis.economicModel.gov.paramProposalVoteDurationSeconds = 640
        genesis.economicModel.gov.textProposalVoteDurationSeconds = 200
        genesis.economicModel.slashing.slashBlocksReward = 0
        new_genesis_env.set_genesis(genesis.to_dict())
        new_genesis_env.deploy_all()
        pip = clients_consensus[0].pip
        pip_test = clients_consensus[1].pip
        address = pip.node.staking_address
        result = pip.submitParam(pip.node.node_id, str(time.time()), 'slashing', 'slashBlocksReward',
                                 '1116', address, transaction_cfg=pip.cfg.transaction_cfg)
        log.info('Submit param proposal result : {}'.format(result))
        assert_code(result, 0)
        result = pip.submitText(pip.node.node_id, str(time.time()), pip.node.staking_address,
                                transaction_cfg=pip.cfg.transaction_cfg)
        log.info('Submit text proposal result : {}'.format(result))
        assert_code(result, 0)

        proposalinfo_param = pip.get_effect_proposal_info_of_vote(pip.cfg.param_proposal)
        log.info('Get param proposal information : {}'.format(proposalinfo_param))
        result = pip.submitCancel(pip.node.node_id, str(time.time()), 14, proposalinfo_param.get('ProposalID'),
                                  pip.node.staking_address, transaction_cfg=pip.cfg.transaction_cfg)
        log.info('Submit cancel result : {}'.format(result))
        assert_code(result, 0)
        result = proposal_vote(pip, proposaltype=pip.cfg.param_proposal)
        assert_code(result, 0)
        result = proposal_vote(pip, proposaltype=pip.cfg.text_proposal)
        assert_code(result, 0)
        result = proposal_vote(pip, proposaltype=pip.cfg.cancel_proposal)
        assert_code(result, 0)
        proposalinfo_cancel = pip.get_effect_proposal_info_of_vote(pip.cfg.cancel_proposal)
        log.info('Get cancel proposal information : {}'.format(proposalinfo_cancel))
        proposalinfo_text = pip.get_effect_proposal_info_of_vote(pip.cfg.text_proposal)
        log.info('Get text proposal information : {}'.format(proposalinfo_text))

        shares0 = clients_consensus[0].staking.get_staking_amount(clients_consensus[0].node)
        log.info('Stop node {}'.format(pip.node.node_id))
        pip.node.stop()
        wait_block_number(pip_test.node, 3 * pip_test.economic.settlement_size)
        balance_before = pip_test.node.eth.getBalance(address, 3 * pip_test.economic.settlement_size - 1)
        log.info('Block bumber {} staking address balance {}'.format(3 * pip_test.economic.settlement_size - 1,
                                                                     balance_before))
        balance_after = pip_test.node.eth.getBalance(address, 3 * pip_test.economic.settlement_size)

        log.info('Block bumber {} staking address balance {}'.format(3 * pip_test.economic.settlement_size,
                                                                     balance_after))
        assert balance_after == balance_before

        wait_block_number(pip_test.node, 5 * pip_test.economic.settlement_size)
        balance_before = pip_test.node.eth.getBalance(address, 5 * pip_test.economic.settlement_size - 1)
        log.info('Block bumber {} staking address balance {}'.format(5 * pip_test.economic.settlement_size - 1,
                                                                     balance_before))
        balance_after = pip_test.node.eth.getBalance(address, 5 * pip_test.economic.settlement_size)

        log.info('Block bumber {} staking address balance {}'.format(5 * pip_test.economic.settlement_size,
                                                                     balance_after))
        assert balance_after - balance_before == shares0

    @pytest.mark.P1
    @allure.title('Node be slashed, verify unstake function')
    def test_UNS_PU_004(self, new_genesis_env, clients_consensus):
        genesis = from_dict(data_class=Genesis, data=new_genesis_env.genesis_config)
        genesis.economicModel.staking.unStakeFreezeDuration = 2
        genesis.economicModel.slashing.maxEvidenceAge = 1
        genesis.economicModel.gov.versionProposalVoteDurationSeconds = 1000
        genesis.economicModel.gov.textProposalVoteDurationSeconds = 200
        genesis.economicModel.slashing.slashBlocksReward = 0
        new_genesis_env.set_genesis(genesis.to_dict())
        new_genesis_env.deploy_all()
        pip = clients_consensus[0].pip
        pip_test = clients_consensus[1].pip
        address = pip.node.staking_address
        result = pip.submitVersion(pip.node.node_id, str(time.time()), pip.cfg.version5, 13, address,
                                   transaction_cfg=pip.cfg.transaction_cfg)
        log.info('Submit version proposal result : {}'.format(result))
        assert_code(result, 0)
        result = pip.submitText(pip.node.node_id, str(time.time()), pip.node.staking_address,
                                transaction_cfg=pip.cfg.transaction_cfg)
        log.info('Submit text proposal result : {}'.format(result))
        assert_code(result, 0)

        proposalinfo_version = pip.get_effect_proposal_info_of_vote(pip.cfg.version_proposal)
        log.info('Get version proposal information : {}'.format(proposalinfo_version))
        result = version_proposal_vote(pip)
        assert_code(result, 0)
        result = proposal_vote(pip, proposaltype=pip.cfg.text_proposal)
        assert_code(result, 0)

        proposalinfo_text = pip.get_effect_proposal_info_of_vote(pip.cfg.text_proposal)
        log.info('Get text proposal information : {}'.format(proposalinfo_text))

        shares0 = clients_consensus[0].staking.get_staking_amount(clients_consensus[0].node)
        log.info('Stop node {}'.format(pip.node.node_id))
        pip.node.stop()
        wait_block_number(pip_test.node, 3 * pip_test.economic.settlement_size)
        balance_before = pip_test.node.eth.getBalance(address, 3 * pip_test.economic.settlement_size - 1)
        log.info('Block bumber {} staking address balance {}'.format(3 * pip_test.economic.settlement_size - 1,
                                                                     balance_before))
        balance_after = pip_test.node.eth.getBalance(address, 3 * pip_test.economic.settlement_size)

        log.info('Block bumber {} staking address balance {}'.format(3 * pip_test.economic.settlement_size,
                                                                     balance_after))
        assert balance_after == balance_before

        wait_block_number(pip_test.node, 4 * pip_test.economic.settlement_size)
        balance_before = pip_test.node.eth.getBalance(address, 4 * pip_test.economic.settlement_size - 1)
        log.info('Block bumber {} staking address balance {}'.format(4 * pip_test.economic.settlement_size - 1,
                                                                     balance_before))
        balance_after = pip_test.node.eth.getBalance(address, 4 * pip_test.economic.settlement_size)

        log.info('Block bumber {} staking address balance {}'.format(4 * pip_test.economic.settlement_size,
                                                                     balance_after))
        assert balance_after - balance_before == shares0

    @pytest.mark.P1
    @allure.title('Node be slashed, verify unstake function')
    def test_UNS_PU_006(self, new_genesis_env, clients_consensus):
        genesis = from_dict(data_class=Genesis, data=new_genesis_env.genesis_config)
        genesis.economicModel.staking.unStakeFreezeDuration = 2
        genesis.economicModel.slashing.maxEvidenceAge = 1
        genesis.economicModel.gov.textProposalVoteDurationSeconds = 520
        genesis.economicModel.slashing.slashBlocksReward = 0
        new_genesis_env.set_genesis(genesis.to_dict())
        new_genesis_env.deploy_all()
        pip = clients_consensus[0].pip
        pip_test = clients_consensus[1].pip
        address = pip.node.staking_address

        result = pip.submitText(pip.node.node_id, str(time.time()), pip.node.staking_address,
                                transaction_cfg=pip.cfg.transaction_cfg)
        log.info('Submit text proposal result : {}'.format(result))
        assert_code(result, 0)

        result = proposal_vote(pip, proposaltype=pip.cfg.text_proposal)
        assert_code(result, 0)
        proposalinfo_text = pip.get_effect_proposal_info_of_vote(pip.cfg.text_proposal)
        log.info('Get text proposal information : {}'.format(proposalinfo_text))

        shares0 = clients_consensus[0].staking.get_staking_amount(clients_consensus[0].node)
        log.info('Stop node {}'.format(pip.node.node_id))
        pip.node.stop()
        wait_block_number(pip_test.node, 3 * pip_test.economic.settlement_size)
        balance_before = pip_test.node.eth.getBalance(address, 3 * pip_test.economic.settlement_size - 1)
        log.info('Block bumber {} staking address balance {}'.format(3 * pip_test.economic.settlement_size - 1,
                                                                     balance_before))
        balance_after = pip_test.node.eth.getBalance(address, 3 * pip_test.economic.settlement_size)

        log.info('Block bumber {} staking address balance {}'.format(3 * pip_test.economic.settlement_size,
                                                                     balance_after))
        assert balance_after == balance_before

        wait_block_number(pip_test.node, 4 * pip_test.economic.settlement_size)
        balance_before = pip_test.node.eth.getBalance(address, 4 * pip_test.economic.settlement_size - 1)
        log.info('Block bumber {} staking address balance {}'.format(4 * pip_test.economic.settlement_size - 1,
                                                                     balance_before))
        balance_after = pip_test.node.eth.getBalance(address, 4 * pip_test.economic.settlement_size)

        log.info('Block bumber {} staking address balance {}'.format(4 * pip_test.economic.settlement_size,
                                                                     balance_after))
        assert balance_after - balance_before == shares0

    @pytest.mark.P2
    @allure.title('Node be slashed, verify unstake function')
    def test_UNS_PU_008(self, new_genesis_env, clients_consensus):
        genesis = from_dict(data_class=Genesis, data=new_genesis_env.genesis_config)
        genesis.economicModel.staking.unStakeFreezeDuration = 2
        genesis.economicModel.slashing.maxEvidenceAge = 1
        genesis.economicModel.gov.paramProposalVoteDurationSeconds = 700
        genesis.economicModel.slashing.slashBlocksReward = 0
        new_genesis_env.set_genesis(genesis.to_dict())
        new_genesis_env.deploy_all()
        pip = clients_consensus[0].pip
        pip_test = clients_consensus[1].pip
        address = pip.node.staking_address
        balance_before = pip.node.eth.getBalance(address)
        result = pip.submitParam(pip.node.node_id, str(time.time()), 'slashing', 'slashBlocksReward',
                                 '1', address, transaction_cfg=pip.cfg.transaction_cfg)
        log.info('Submit param proposal result : {}'.format(result))
        assert_code(result, 0)
        result = proposal_vote(pip, proposaltype=pip.cfg.param_proposal)
        assert_code(result, 0)
        proposalinfo_param = pip.get_effect_proposal_info_of_vote(pip.cfg.param_proposal)
        log.info('Get param proposal information : {}'.format(proposalinfo_param))

        result = pip.submitCancel(pip.node.node_id, str(time.time()), 14, proposalinfo_param.get('ProposalID'),
                                  pip.node.staking_address, transaction_cfg=pip.cfg.transaction_cfg)
        log.info('Submit cancel result : {}'.format(result))
        assert_code(result, 0)
        result = proposal_vote(pip, proposaltype=pip.cfg.cancel_proposal)
        assert_code(result, 0)
        proposalinfo_cancel = pip.get_effect_proposal_info_of_vote(pip.cfg.cancel_proposal)
        log.info('Get cancel proposal information : {}'.format(proposalinfo_cancel))

        shares0 = clients_consensus[0].staking.get_staking_amount(clients_consensus[0].node)
        log.info('get staking amount {}'.format(shares0))
        log.info('Stop node {}'.format(pip.node.node_id))
        pip.node.stop()

        # wait_block_number(pip_test.node, 3 * pip_test.economic.settlement_size)
        # balance_before = pip_test.node.eth.getBalance(address, 3 * pip_test.economic.settlement_size - 1)
        # log.info('Block bumber {} staking address balance {}'.format(3 * pip_test.economic.settlement_size - 1,
        #                                                              balance_before))
        # balance_after = pip_test.node.eth.getBalance(address, 3 * pip_test.economic.settlement_size)
        # log.info('Block bumber {} staking address balance {}'.format(3 * pip_test.economic.settlement_size,
        #                                                              balance_after))
        # assert balance_after == balance_before

        wait_block_number(pip_test.node, 4 * pip_test.economic.settlement_size)
        balance_before = pip_test.node.eth.getBalance(address, 4 * pip_test.economic.settlement_size - 1)
        log.info('Block bumber {} staking address balance {}'.format(4 * pip_test.economic.settlement_size - 1,
                                                                     balance_before))
        wait_block_number(pip_test.node, 8 * pip_test.economic.settlement_size)
        balance_after = pip_test.node.eth.getBalance(address, 8 * pip_test.economic.settlement_size + 1)
        log.info('Block bumber {} staking address balance {}'.format(8 * pip_test.economic.settlement_size,
                                                                     balance_after))
        assert balance_after - balance_before == shares0

    @pytest.mark.P2
    @allure.title('Node be slashed, verify unstake function')
    def test_UNS_PU_009_011_013_019(self, new_genesis_env, clients_noconsensus):
        genesis = from_dict(data_class=Genesis, data=new_genesis_env.genesis_config)
        genesis.economicModel.staking.unStakeFreezeDuration = 2
        genesis.economicModel.slashing.maxEvidenceAge = 1
        genesis.economicModel.gov.paramProposalVoteDurationSeconds = 300
        genesis.economicModel.slashing.slashBlocksReward = 0
        new_genesis_env.set_genesis(genesis.to_dict())
        new_genesis_env.deploy_all()
        pip = clients_noconsensus[0].pip
        pip_test = clients_noconsensus[1].pip

        create_lockup_plan(clients_noconsensus[0])
        address = pip.node.staking_address
        submitvpandvote([clients_noconsensus[0]], votingrounds=1)
        proposalinfo_version = pip.get_effect_proposal_info_of_vote()
        log.info('Get version proposal information : {}'.format(proposalinfo_version))
        wait_block_number(pip.node, proposalinfo_version.get('EndVotingBlock'))
        submitcppandvote([clients_noconsensus[0]], [1])
        submittpandvote([clients_noconsensus[0]], 1)
        proposalinfo_param = pip.get_effect_proposal_info_of_vote(pip.cfg.param_proposal)
        log.info('Get param proposal information : {}'.format(proposalinfo_param))
        proposal_vote(pip, proposaltype=pip.cfg.param_proposal)
        log.info('Stop the node {}'.format(pip.node.node_id))
        shares = clients_noconsensus[0].staking.get_staking_amount()
        log.info('Node staking amount : {}'.format(shares))
        pip.node.stop()

        self.verify_amount_block(pip_test, address, shares, value=4)

    @pytest.mark.P2
    @allure.title('Node be slashed, verify unstake function')
    def test_UNS_PU_010(self, new_genesis_env, clients_noconsensus):
        genesis = from_dict(data_class=Genesis, data=new_genesis_env.genesis_config)
        genesis.economicModel.staking.unStakeFreezeDuration = 2
        genesis.economicModel.slashing.maxEvidenceAge = 1
        genesis.economicModel.gov.versionProposalVoteDurationSeconds = 2000
        genesis.economicModel.slashing.slashBlocksReward = 0
        new_genesis_env.set_genesis(genesis.to_dict())
        new_genesis_env.deploy_all()
        pip = clients_noconsensus[0].pip
        pip_test = clients_noconsensus[1].pip

        create_lockup_plan(clients_noconsensus[0])
        address = pip.node.staking_address
        submitvpandvote([clients_noconsensus[0]], votingrounds=14)
        shares = clients_noconsensus[0].staking.get_staking_amount()
        log.info('Node staking amount : {}'.format(shares))
        pip.node.stop()
        self.verify_amount(pip_test, address, shares)

    @pytest.mark.P2
    @allure.title('Node be slashed, verify unstake function')
    def test_UNS_PU_012(self, new_genesis_env, clients_noconsensus):
        genesis = from_dict(data_class=Genesis, data=new_genesis_env.genesis_config)
        genesis.economicModel.staking.unStakeFreezeDuration = 2
        genesis.economicModel.slashing.maxEvidenceAge = 1
        genesis.economicModel.gov.textProposalVoteDurationSeconds = 480
        genesis.economicModel.slashing.slashBlocksReward = 0
        new_genesis_env.set_genesis(genesis.to_dict())
        new_genesis_env.deploy_all()
        pip = clients_noconsensus[0].pip
        pip_test = clients_noconsensus[1].pip

        create_lockup_plan(clients_noconsensus[0])
        address = pip.node.staking_address
        submittpandvote([clients_noconsensus[0]], 1)
        shares = clients_noconsensus[0].staking.get_staking_amount()
        log.info('Node staking amount : {}'.format(shares))
        pip.node.stop()
        self.verify_amount(pip_test, address, shares)

    @pytest.mark.P2
    @allure.title('Node be slashed, verify unstake function')
    def test_UNS_PU_014(self, new_genesis_env, clients_noconsensus):
        genesis = from_dict(data_class=Genesis, data=new_genesis_env.genesis_config)
        genesis.economicModel.staking.unStakeFreezeDuration = 2
        genesis.economicModel.slashing.maxEvidenceAge = 1
        genesis.economicModel.gov.paramProposalVoteDurationSeconds = 640
        genesis.economicModel.slashing.slashBlocksReward = 0
        new_genesis_env.set_genesis(genesis.to_dict())
        new_genesis_env.deploy_all()
        pip = clients_noconsensus[0].pip
        pip_test = clients_noconsensus[1].pip

        create_lockup_plan(clients_noconsensus[0])
        address = pip.node.staking_address
        result = pip.submitParam(pip.node.node_id, str(time.time()), 'slashing', 'slashBlocksReward', '1116',
                                 address, transaction_cfg=pip.cfg.transaction_cfg)
        log.info('Submit param proposal result : {}'.format(result))
        proposalinfo_param = pip.get_effect_proposal_info_of_vote(pip.cfg.param_proposal)
        log.info('Get param proposal information : {}'.format(proposalinfo_param))
        result = pip.submitCancel(pip.node.node_id, str(time.time()), 13, proposalinfo_param.get('ProposalID'),
                                  address, transaction_cfg=pip.cfg.transaction_cfg)
        log.info('Submit cancel proposal result : {}'.format(result))
        result = proposal_vote(pip, proposaltype=pip.cfg.cancel_proposal)
        assert_code(result, 0)
        shares = clients_noconsensus[0].staking.get_staking_amount()
        log.info('Node staking amount : {}'.format(shares))
        pip.node.stop()
        self.verify_amount(pip_test, address, shares)

    @pytest.mark.P2
    @allure.title('Node be slashed, verify unstake function')
    def test_UNS_PU_020(self, new_genesis_env, clients_noconsensus):
        genesis = from_dict(data_class=Genesis, data=new_genesis_env.genesis_config)
        genesis.economicModel.staking.unStakeFreezeDuration = 2
        genesis.economicModel.slashing.maxEvidenceAge = 1
        genesis.economicModel.gov.paramProposalVoteDurationSeconds = 480
        genesis.economicModel.slashing.slashBlocksReward = 0
        new_genesis_env.set_genesis(genesis.to_dict())
        new_genesis_env.deploy_all()
        pip = clients_noconsensus[0].pip
        pip_test = clients_noconsensus[1].pip

        create_lockup_plan(clients_noconsensus[0])
        address = pip.node.staking_address
        submitppandvote([clients_noconsensus[0]], 3)
        shares = clients_noconsensus[0].staking.get_staking_amount()
        log.info('Node staking amount : {}'.format(shares))
        pip.node.stop()
        self.verify_amount(pip_test, address, shares)

    def verify_amount(self, pip, address, shares):
        self.verify_amount_block(pip, address, shares, value=4, tag=False)
        self.verify_amount_block(pip, address, shares, value=5)

    def verify_amount_block(self, pip, address, shares, value, tag=True):
        wait_block_number(pip.node, value * pip.economic.settlement_size)
        balance_before = pip.node.eth.getBalance(address, value * pip.economic.settlement_size - 1)
        balance_before_lockup = pip.node.eth.getBalance(pip.cfg.FOUNDATION_LOCKUP_ADDRESS,
                                                        value * pip.economic.settlement_size - 1)
        log.info('Block bumber {} staking address balance {}'.format(value * pip.economic.settlement_size - 1,
                                                                     balance_before))
        log.info('Block bumber {} FOUNDATION_LOCKUP_ADDRESS balance {}'.format(value * pip.economic.settlement_size - 1,
                                                                               balance_before_lockup))
        balance_after = pip.node.eth.getBalance(address, value * pip.economic.settlement_size)
        balance_after_lockup = pip.node.eth.getBalance(pip.cfg.FOUNDATION_LOCKUP_ADDRESS,
                                                       value * pip.economic.settlement_size)
        log.info('Block bumber {} staking address balance {}'.format(value * pip.economic.settlement_size,
                                                                     balance_after))
        log.info('Block bumber {} FOUNDATION_LOCKUP_ADDRESS balance {}'.format(value * pip.economic.settlement_size,
                                                                               balance_after_lockup))
        assert balance_after == balance_before
        if tag:
            assert balance_after_lockup - balance_before_lockup == shares
        else:
            assert balance_after_lockup == balance_before_lockup

def test_fixbug(new_genesis_env, clients_consensus):
    genesis = from_dict(data_class=Genesis, data=new_genesis_env.genesis_config)
    genesis.economicModel.gov.versionProposalVoteDurationSeconds = 1000
    new_genesis_env.set_genesis(genesis.to_dict())
    new_genesis_env.deploy_all()
    pip_stop = clients_consensus[0].pip
    pip = clients_consensus[1].pip
    submitvpandvote(clients_consensus, votingrounds=15)
    proprosalinfo = pip.get_effect_proposal_info_of_vote()
    log.info('Proposalinfo : {}'.format(proprosalinfo))
    log.info('Stop node {}'.format(pip_stop.node.node_id))
    log.info('stop node nodeid {}'.format(pip_stop.node.node_id))
    pip_stop.node.stop()
    pip.economic.wait_settlement_blocknum(pip.node)
    pip.economic.wait_consensus_blocknum(pip.node, 1)
    verifier_list = get_pledge_list(clients_consensus[1].ppos.getVerifierList)
    log.info('Verifier list : {}'.format(verifier_list))
    validator_list = get_pledge_list(clients_consensus[1].ppos.getValidatorList)
    log.info('Validator list : {}'.format(validator_list))
    assert pip_stop.node.node_id not in verifier_list
    assert pip_stop.node.node_id not in validator_list
    wait_block_number(pip.node, proprosalinfo.get('ActiveBlock'))
    assert pip.chain_version == proprosalinfo.get('NewVersion')
    verifier_list = get_pledge_list(clients_consensus[1].ppos.getVerifierList)
    log.info('Verifier list : {}'.format(verifier_list))
    validator_list = get_pledge_list(clients_consensus[1].ppos.getValidatorList)
    log.info('Validator list : {}'.format(validator_list))
    assert pip_stop.node.node_id not in verifier_list
    assert pip_stop.node.node_id not in validator_list
    result = clients_consensus[1].ppos.getCandidateInfo(pip_stop.node.node_id)
    log.info('Get nodeid {} candidate infor {}'.format(pip_stop.node.node_id, result))
    assert_code(result, 301204)
    assert result.get('Ret') == 'Query candidate info failed:Candidate info is not found'



