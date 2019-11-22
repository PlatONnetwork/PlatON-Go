from tests.lib.utils import upload_platon, assert_code, get_pledge_list, wait_block_number
from common.log import log
from tests.lib.client import Client, get_client_obj, get_client_obj_list, StakingConfig
import pytest, allure
import time
import math
from tests.govern.conftest import version_proposal_vote, proposal_vote
from tests.lib import Genesis, PipConfig
from dacite import from_dict
from tests.govern.test_voting_statistics import submitcvpandvote, submitcppandvote, submittpandvote, \
    submitvpandvote, submitppandvote


def create_lockup_plan(client_obj):
    address, _ = client_obj.pip.economic.account.generate_account(client_obj.node.web3,
                                                                  3*client_obj.economic.genesis.economicModel.staking.stakeThreshold)
    plan = [{'Epoch': 20, 'Amount': 2*client_obj.economic.genesis.economicModel.staking.stakeThreshold}]
    result = client_obj.restricting.createRestrictingPlan(address, plan, address,
                                                          transaction_cfg=client_obj.pip.cfg.transaction_cfg)
    log.info('CreateRestrictingPlan result : {}'.format(result))
    assert_code(result, 0)
    result = client_obj.staking.create_staking(1, address, address,
                                               amount=int(1.8*client_obj.economic.genesis.economicModel.staking.stakeThreshold),
                                               transaction_cfg=client_obj.pip.cfg.transaction_cfg)
    log.info('Create staking result : {}'.format(result))
    assert_code(result, 0)
    client_obj.economic.wait_settlement_blocknum(client_obj.node)


@pytest.fixture()
def new_node_no_proposal(no_vp_proposal, client_noc_list_obj, client_list_obj):
    pip = no_vp_proposal
    client_obj = get_client_obj(pip.node.node_id, client_list_obj)
    candidate_list = get_pledge_list(client_obj.ppos.getCandidateList)
    log.info('candidate_list: {}'.format(candidate_list))
    for client_obj in client_noc_list_obj:
        if client_obj.node.node_id not in candidate_list:
            return client_obj.pip
    log.info('All nodes are staked, restart the chain')
    pip.economic.env.deploy_all()
    return client_noc_list_obj[0].pip


def replace_platon_and_staking(pip, platon_bin):
    node_obj_list = pip.economic.env.get_all_nodes()
    client_list_obj = []
    for node_obj in node_obj_list:
        client_list_obj.append(Client(pip.economic.env, node_obj, StakingConfig("externalId", "nodeName", "website",
                                                                                    "details")))
    client_obj = get_client_obj(pip.node.node_id, client_list_obj)
    upload_platon(pip.node, platon_bin)
    log.info('Replace the platon of the node {}'.format(pip.node.node_id))
    pip.node.restart()
    log.info('Restart the node {}'.format(pip.node.node_id))
    address, _ = pip.economic.account.generate_account(pip.node.web3,
                                                           10*pip.economic.genesis.economicModel.staking.stakeThreshold)
    result = client_obj.staking.create_staking(0, address, address, transaction_cfg=pip.cfg.transaction_cfg)
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
    def preactive_proposal(self, client_list_obj):
        verifier_list = get_pledge_list(client_list_obj[0].ppos.getVerifierList)
        log.info('verifierlist :{}'.format(verifier_list))
        client_verifier_list_obj = get_client_obj_list(verifier_list, client_list_obj)
        pip_list_obj = [client_obj.pip for client_obj in client_verifier_list_obj]
        result = pip_list_obj[0].submitVersion(pip_list_obj[0].node.node_id, str(time.time()),
                                               pip_list_obj[0].cfg.version5, 2, pip_list_obj[0].node.staking_address,
                                               transaction_cfg=pip_list_obj[0].cfg.transaction_cfg)
        log.info('submit version proposal, result : {}'.format(result))
        proposalinfo = pip_list_obj[0].get_effect_proposal_info_of_vote()
        log.info('Version proposalinfo: {}'.format(proposalinfo))
        for pip in pip_list_obj:
            result = version_proposal_vote(pip)
            assert_code(result, 0)
        wait_block_number(pip.node, proposalinfo.get('EndVotingBlock'))
        assert pip.get_status_of_proposal(proposalinfo.get('ProposalID')) == 4

    @pytest.mark.P1
    @allure.title('There is preactive proposal, verify stake function')
    def test_ST_PR_001(self, new_genesis_env, new_node_no_proposal, client_list_obj):
        pip = new_node_no_proposal
        self.preactive_proposal(client_list_obj)
        result = replace_platon_and_staking(pip, pip.cfg.PLATON_NEW_BIN2)
        assert_code(result, 301004)

    @pytest.mark.P1
    @allure.title('There is preactive proposal, verify stake function')
    def test_ST_PR_002(self, new_genesis_env, new_node_no_proposal, client_list_obj):
        pip = new_node_no_proposal
        self.preactive_proposal(client_list_obj)
        result = replace_platon_and_staking(pip, pip.cfg.PLATON_NEW_BIN1)
        assert_code(result, 301004)

    @pytest.mark.P1
    @allure.title('There is preactive proposal, verify stake function')
    def test_ST_PR_003(self, new_genesis_env, new_node_no_proposal, client_list_obj):
        pip = new_node_no_proposal
        self.preactive_proposal(client_list_obj)
        result = replace_platon_and_staking(pip, pip.cfg.PLATON_NEW_BIN0)
        assert_code(result, 301004)

    @pytest.mark.P1
    @allure.title('There is preactive proposal, verify stake function')
    def test_ST_PR_004(self, new_genesis_env, new_node_no_proposal, client_list_obj):
        pip = new_node_no_proposal
        self.preactive_proposal(client_list_obj)
        result = replace_platon_and_staking(pip, pip.cfg.PLATON_NEW_BIN3)
        assert_code(result, 301004)

    @pytest.mark.P1
    @allure.title('There is preactive proposal, verify stake function')
    def test_ST_PR_005(self, new_genesis_env, new_node_no_proposal, client_list_obj):
        pip = new_node_no_proposal
        self.preactive_proposal(client_list_obj)
        result = replace_platon_and_staking(pip, pip.cfg.PLATON_NEW_BIN)
        assert_code(result, 0)

    @pytest.mark.P1
    @allure.title('There is preactive proposal, verify stake function')
    def test_ST_PR_006(self, new_genesis_env, new_node_no_proposal, client_list_obj):
        pip = new_node_no_proposal
        self.preactive_proposal(client_list_obj)
        result = replace_platon_and_staking(pip, pip.cfg.PLATON_NEW_BIN8)
        assert_code(result, 301005)


class TestUpgradedProposalStaking:
    def upgraded_proposal(self, client_list_obj):
        verifier_list = get_pledge_list(client_list_obj[0].ppos.getVerifierList)
        log.info('verifierlist :{}'.format(verifier_list))
        client_verifier_list_obj = get_client_obj_list(verifier_list, client_list_obj)
        pip_list_obj = [client_obj.pip for client_obj in client_verifier_list_obj]
        result = pip_list_obj[0].submitVersion(pip_list_obj[0].node.node_id, str(time.time()),
                                               pip_list_obj[0].cfg.version5, 2, pip_list_obj[0].node.staking_address,
                                               transaction_cfg=pip_list_obj[0].cfg.transaction_cfg)
        log.info('submit version proposal, result : {}'.format(result))
        proposalinfo = pip_list_obj[0].get_effect_proposal_info_of_vote()
        log.info('Version proposalinfo: {}'.format(proposalinfo))
        for pip in pip_list_obj:
            result = version_proposal_vote(pip)
            assert_code(result, 0)
        wait_block_number(pip.node, proposalinfo.get('EndVotingBlock'))
        assert pip.get_status_of_proposal(proposalinfo.get('ProposalID')) == 4
        wait_block_number(pip.node, proposalinfo.get('ActiveBlock'))
        assert pip.get_status_of_proposal(proposalinfo.get('ProposalID')) == 5

    @pytest.mark.P2
    @allure.title('Chain upgrade completed, verify stake function')
    def test_ST_UPG_001(self, new_genesis_env, new_node_no_proposal, client_list_obj):
        pip = new_node_no_proposal
        self.upgraded_proposal(client_list_obj)
        result = replace_platon_and_staking(pip, pip.cfg.PLATON_NEW_BIN4)
        assert_code(result, 0)

    @pytest.mark.P2
    @allure.title('Chain upgrade completed, verify stake function')
    def test_ST_UPG_002(self, new_genesis_env, new_node_no_proposal, client_list_obj):
        pip = new_node_no_proposal
        self.upgraded_proposal(client_list_obj)
        result = replace_platon_and_staking(pip, pip.cfg.PLATON_NEW_BIN0)
        assert_code(result, 301004)

    @pytest.mark.P2
    @allure.title('Chain upgrade completed, verify stake function')
    def test_ST_UPG_003(self, new_genesis_env, new_node_no_proposal, client_list_obj):
        pip = new_node_no_proposal
        self.upgraded_proposal(client_list_obj)
        result = replace_platon_and_staking(pip, pip.cfg.PLATON_NEW_BIN)
        assert_code(result, 0)

    @pytest.mark.P2
    @allure.title('Chain upgrade completed, verify stake function')
    def test_ST_UPG_004(self, new_genesis_env, new_node_no_proposal, client_list_obj):
        pip = new_node_no_proposal
        self.upgraded_proposal(client_list_obj)
        result = replace_platon_and_staking(pip, pip.cfg.PLATON_NEW_BIN6)
        assert_code(result, 0)

    @pytest.mark.P2
    @allure.title('Chain upgrade completed, verify stake function')
    def test_ST_UPG_005(self, new_genesis_env, new_node_no_proposal, client_list_obj):
        pip = new_node_no_proposal
        self.upgraded_proposal(client_list_obj)
        result = replace_platon_and_staking(pip, pip.cfg.PLATON_NEW_BIN7)
        assert_code(result, 301005)


class TestUnstaking:
    @pytest.mark.P1
    @allure.title('Verify unstake function')
    def test_UNS_AM_003_007(self, new_genesis_env, client_verifier_obj):
        genesis = from_dict(data_class=Genesis, data=new_genesis_env.genesis_config)
        genesis.economicModel.staking.unStakeFreezeDuration == 2
        genesis.economicModel.slashing.maxEvidenceAge == 1
        new_genesis_env.set_genesis(genesis.to_dict())
        new_genesis_env.deploy_all()
        pip = client_verifier_obj.pip
        address = pip.node.staking_address
        list_obj = [client_verifier_obj]
        submitcvpandvote(list_obj, 1)
        result = version_proposal_vote(pip)
        assert_code(result, 0)
        shares = client_verifier_obj.staking.get_staking_amount(pip.node)
        result = client_verifier_obj.staking.withdrew_staking(address)
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
    def test_UNS_AM_005(self, new_genesis_env, client_verifier_obj):
        genesis = from_dict(data_class=Genesis, data=new_genesis_env.genesis_config)
        genesis.economicModel.staking.unStakeFreezeDuration == 2
        genesis.economicModel.slashing.maxEvidenceAge == 1
        new_genesis_env.set_genesis(genesis.to_dict())
        new_genesis_env.deploy_all()
        pip = client_verifier_obj.pip
        address = pip.node.staking_address
        result = pip.submitText(pip.node.node_id, str(time.time()), address, transaction_cfg=pip.cfg.transaction_cfg)
        assert_code(result, 0)

        submitcppandvote([client_verifier_obj], [1])
        result = proposal_vote(pip, proposaltype=pip.cfg.text_proposal)
        assert_code(result, 0)
        shares = client_verifier_obj.staking.get_staking_amount(pip.node)
        result = client_verifier_obj.staking.withdrew_staking(address)
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
    def test_UNS_AM_004_006_008(self, new_genesis_env, client_verifier_obj_list):
        genesis = from_dict(data_class=Genesis, data=new_genesis_env.genesis_config)
        genesis.economicModel.staking.unStakeFreezeDuration == 2
        genesis.economicModel.slashing.maxEvidenceAge == 1
        genesis.economicModel.gov.versionProposalVoteDurationSeconds = 1000
        genesis.economicModel.gov.textProposalVoteDurationSeconds = 840
        new_genesis_env.set_genesis(genesis.to_dict())
        new_genesis_env.deploy_all()
        pip_one = client_verifier_obj_list[0].pip
        pip_two = client_verifier_obj_list[1].pip
        pip_three = client_verifier_obj_list[2].pip
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
        shares_one = client_verifier_obj_list[0].staking.get_staking_amount(pip_one.node)
        shares_two = client_verifier_obj_list[1].staking.get_staking_amount(pip_two.node)
        shares_three = client_verifier_obj_list[2].staking.get_staking_amount(pip_three.node)
        result = client_verifier_obj_list[0].staking.withdrew_staking(address)
        log.info('Node {} withdrew staking result : {}'.format(pip_one.node.node_id, result))
        assert_code(result, 0)

        result = client_verifier_obj_list[1].staking.withdrew_staking(address)
        log.info('Node {} withdrew staking result : {}'.format(pip_two.node.node_id, result))
        assert_code(result, 0)

        result = client_verifier_obj_list[2].staking.withdrew_staking(address)
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
    def test_UNS_AM_009_011_013(self, new_genesis_env, client_noc_list_obj):
        genesis = from_dict(data_class=Genesis, data=new_genesis_env.genesis_config)
        genesis.economicModel.staking.unStakeFreezeDuration == 2
        genesis.economicModel.slashing.maxEvidenceAge == 1
        genesis.economicModel.gov.versionProposalVoteDurationSeconds = 1000
        genesis.economicModel.gov.textProposalVoteDurationSeconds = 200
        new_genesis_env.set_genesis(genesis.to_dict())
        new_genesis_env.deploy_all()
        pip_test = client_noc_list_obj[0].pip
        address, _ = pip_test.economic.account.generate_account(pip_test.node.web3, 10**18 * 20000000)
        plan = [{'Epoch': 20, 'Amount': 10**18 * 2000000}]
        result = client_noc_list_obj[0].restricting.createRestrictingPlan(address, plan, address,
                                                                          transaction_cfg=pip_test.cfg.transaction_cfg)
        log.info('CreateRestrictingPlan result : {}'.format(result))
        assert_code(result, 0)
        result = client_noc_list_obj[0].staking.create_staking(1, address, address,
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

        shares = client_noc_list_obj[0].staking.get_staking_amount(pip_test.node)
        result = client_noc_list_obj[0].staking.withdrew_staking(address)
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
    def test_UNS_AM_010_012_014(self, new_genesis_env, client_noc_list_obj):
        genesis = from_dict(data_class=Genesis, data=new_genesis_env.genesis_config)
        genesis.economicModel.staking.unStakeFreezeDuration == 2
        genesis.economicModel.slashing.maxEvidenceAge == 1
        genesis.economicModel.gov.paramProposalVoteDurationSeconds = 640
        genesis.economicModel.gov.textProposalVoteDurationSeconds = 840
        new_genesis_env.set_genesis(genesis.to_dict())
        new_genesis_env.deploy_all()
        for client_obj in client_noc_list_obj:
            pip = client_obj.pip
            address, _ = pip.economic.account.generate_account(pip.node.web3, 10**18 * 20000000)
            plan = [{'Epoch': 20, 'Amount': 10**18 * 2000000}]
            result = client_obj.restricting.createRestrictingPlan(address, plan, address,
                                                                  transaction_cfg=pip.cfg.transaction_cfg)
            log.info('CreateRestrictingPlan result : {}'.format(result))
            assert_code(result, 0)
            result = client_obj.staking.create_staking(1, address, address, amount=10**18 * 1800000,
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
        result = proposal_vote(client_noc_list_obj[0].pip, proposaltype=pip.cfg.param_proposal)
        assert_code(result, 0)
        result = proposal_vote(client_noc_list_obj[1].pip, proposaltype=pip.cfg.text_proposal)
        assert_code(result, 0)
        result = proposal_vote(client_noc_list_obj[2].pip, proposaltype=pip.cfg.cancel_proposal)
        assert_code(result, 0)
        proposalinfo_cancel = pip.get_effect_proposal_info_of_vote(pip.cfg.cancel_proposal)
        log.info('Get cancel proposal information : {}'.format(proposalinfo_cancel))
        proposalinfo_text = pip.get_effect_proposal_info_of_vote(pip.cfg.text_proposal)
        log.info('Get text proposal information : {}'.format(proposalinfo_text))

        shares0 = client_noc_list_obj[0].staking.get_staking_amount(client_noc_list_obj[0].node)
        shares1 = client_noc_list_obj[1].staking.get_staking_amount(client_noc_list_obj[1].node)
        shares2 = client_noc_list_obj[2].staking.get_staking_amount(client_noc_list_obj[2].node)
        address0 = client_noc_list_obj[0].node.staking_address
        address1 = client_noc_list_obj[1].node.staking_address
        address2 = client_noc_list_obj[2].node.staking_address
        result = client_noc_list_obj[0].staking.withdrew_staking(address0)
        log.info('Node {} withdrew result : {}'.format(client_noc_list_obj[0].node.node_id, result))
        assert_code(result, 0)
        result = client_noc_list_obj[1].staking.withdrew_staking(address1)
        log.info('Node {} withdrew result : {}'.format(client_noc_list_obj[0].node.node_id, result))
        assert_code(result, 0)
        result = client_noc_list_obj[2].staking.withdrew_staking(address2)
        log.info('Node {} withdrew result : {}'.format(client_noc_list_obj[0].node.node_id, result))
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
    def test_UNS_PU_003_005_007_017(self, new_genesis_env, client_con_list_obj):
        genesis = from_dict(data_class=Genesis, data=new_genesis_env.genesis_config)
        genesis.economicModel.staking.unStakeFreezeDuration = 2
        genesis.economicModel.slashing.maxEvidenceAge = 1
        genesis.economicModel.gov.paramProposalVoteDurationSeconds = 200
        genesis.economicModel.gov.textProposalVoteDurationSeconds = 200
        genesis.economicModel.slashing.slashBlocksReward = 0
        new_genesis_env.set_genesis(genesis.to_dict())
        new_genesis_env.deploy_all()
        pip = client_con_list_obj[0].pip
        pip_test = client_con_list_obj[1].pip
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
        submittpandvote([client_con_list_obj[0]], 3)
        submitcppandvote([client_con_list_obj[0]], [2])
        result = proposal_vote(pip, proposaltype=pip.cfg.param_proposal)
        assert_code(result, 0)
        log.info('Stop the node {}'.format(pip.node.node_id))
        shares = client_con_list_obj[1].staking.get_staking_amount(pip_test.node)
        pip.node.stop()
        wait_block_number(pip_test.node, 4 * pip_test.economic.settlement_size)
        balance_before = pip_test.node.eth.getBalance(address, 4 * pip_test.economic.settlement_size - 1)
        log.info('Block bumber {} staking address balance {}'.format(4 * pip_test.economic.settlement_size - 1,
                                                                     balance_before))
        balance_after = pip_test.node.eth.getBalance(address, 4 * pip_test.economic.settlement_size)

        log.info('Block bumber {} staking address balance {}'.format(4 * pip_test.economic.settlement_size,
                                                                     balance_after))
        assert balance_after - balance_before == shares

    @pytest.mark.P2
    @allure.title('Node be slashed, verify unstake function')
    def test_UNS_PU_016(self, new_genesis_env, client_con_list_obj):
        genesis = from_dict(data_class=Genesis, data=new_genesis_env.genesis_config)
        genesis.economicModel.staking.unStakeFreezeDuration = 2
        genesis.economicModel.slashing.maxEvidenceAge = 1
        genesis.economicModel.gov.paramProposalVoteDurationSeconds = 640
        genesis.economicModel.gov.textProposalVoteDurationSeconds = 200
        genesis.economicModel.slashing.slashBlocksReward = 0
        new_genesis_env.set_genesis(genesis.to_dict())
        new_genesis_env.deploy_all()
        pip = client_con_list_obj[0].pip
        pip_test = client_con_list_obj[1].pip
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

        shares0 = client_con_list_obj[0].staking.get_staking_amount(client_con_list_obj[0].node)
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
    def test_UNS_PU_004(self, new_genesis_env, client_con_list_obj):
        genesis = from_dict(data_class=Genesis, data=new_genesis_env.genesis_config)
        genesis.economicModel.staking.unStakeFreezeDuration = 2
        genesis.economicModel.slashing.maxEvidenceAge = 1
        genesis.economicModel.gov.versionProposalVoteDurationSeconds = 1000
        genesis.economicModel.gov.textProposalVoteDurationSeconds = 200
        genesis.economicModel.slashing.slashBlocksReward = 0
        new_genesis_env.set_genesis(genesis.to_dict())
        new_genesis_env.deploy_all()
        pip = client_con_list_obj[0].pip
        pip_test = client_con_list_obj[1].pip
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

        shares0 = client_con_list_obj[0].staking.get_staking_amount(client_con_list_obj[0].node)
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
    def test_UNS_PU_006(self, new_genesis_env, client_con_list_obj):
        genesis = from_dict(data_class=Genesis, data=new_genesis_env.genesis_config)
        genesis.economicModel.staking.unStakeFreezeDuration = 2
        genesis.economicModel.slashing.maxEvidenceAge = 1
        genesis.economicModel.gov.textProposalVoteDurationSeconds = 520
        genesis.economicModel.slashing.slashBlocksReward = 0
        new_genesis_env.set_genesis(genesis.to_dict())
        new_genesis_env.deploy_all()
        pip = client_con_list_obj[0].pip
        pip_test = client_con_list_obj[1].pip
        address = pip.node.staking_address

        result = pip.submitText(pip.node.node_id, str(time.time()), pip.node.staking_address,
                                    transaction_cfg=pip.cfg.transaction_cfg)
        log.info('Submit text proposal result : {}'.format(result))
        assert_code(result, 0)

        result = proposal_vote(pip, proposaltype=pip.cfg.text_proposal)
        assert_code(result, 0)
        proposalinfo_text = pip.get_effect_proposal_info_of_vote(pip.cfg.text_proposal)
        log.info('Get text proposal information : {}'.format(proposalinfo_text))

        shares0 = client_con_list_obj[0].staking.get_staking_amount(client_con_list_obj[0].node)
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
    def test_UNS_PU_008(self, new_genesis_env, client_con_list_obj):
        genesis = from_dict(data_class=Genesis, data=new_genesis_env.genesis_config)
        genesis.economicModel.staking.unStakeFreezeDuration = 2
        genesis.economicModel.slashing.maxEvidenceAge = 1
        genesis.economicModel.gov.paramProposalVoteDurationSeconds = 700
        genesis.economicModel.slashing.slashBlocksReward = 0
        new_genesis_env.set_genesis(genesis.to_dict())
        new_genesis_env.deploy_all()
        pip = client_con_list_obj[0].pip
        pip_test = client_con_list_obj[1].pip
        address = pip.node.staking_address
        result = pip.submitParam(pip.node.node_id, str(time.time()), 'slashing', 'slashBlocksReward',
                                     '1116', address, transaction_cfg=pip.cfg.transaction_cfg)
        log.info('Submit param proposal result : {}'.format(result))
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

        shares0 = client_con_list_obj[0].staking.get_staking_amount(client_con_list_obj[0].node)
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
    def test_UNS_PU_009_011_013_019(self, new_genesis_env, client_noc_list_obj):
        genesis = from_dict(data_class=Genesis, data=new_genesis_env.genesis_config)
        genesis.economicModel.staking.unStakeFreezeDuration = 2
        genesis.economicModel.slashing.maxEvidenceAge = 1
        genesis.economicModel.gov.paramProposalVoteDurationSeconds = 300
        genesis.economicModel.slashing.slashBlocksReward = 0
        new_genesis_env.set_genesis(genesis.to_dict())
        new_genesis_env.deploy_all()
        pip = client_noc_list_obj[0].pip
        pip_test = client_noc_list_obj[1].pip

        create_lockup_plan(client_noc_list_obj[0])
        address = pip.node.staking_address
        submitvpandvote([client_noc_list_obj[0]], votingrounds=1)
        proposalinfo_version = pip.get_effect_proposal_info_of_vote()
        log.info('Get version proposal information : {}'.format(proposalinfo_version))
        wait_block_number(pip.node, proposalinfo_version.get('EndVotingBlock'))
        submitcppandvote([client_noc_list_obj[0]], [1])
        submittpandvote([client_noc_list_obj[0]], 1)
        proposalinfo_param = pip.get_effect_proposal_info_of_vote(pip.cfg.param_proposal)
        log.info('Get param proposal information : {}'.format(proposalinfo_param))
        proposal_vote(pip, proposaltype=pip.cfg.param_proposal)
        log.info('Stop the node {}'.format(pip.node.node_id))
        shares = client_noc_list_obj[0].staking.get_staking_amount()
        log.info('Node staking amount : {}'.format(shares))
        pip.node.stop()

        self.verify_amount_block(pip_test, address, shares, value=4)

    @pytest.mark.P2
    @allure.title('Node be slashed, verify unstake function')
    def test_UNS_PU_010(self, new_genesis_env, client_noc_list_obj):
        genesis = from_dict(data_class=Genesis, data=new_genesis_env.genesis_config)
        genesis.economicModel.staking.unStakeFreezeDuration = 2
        genesis.economicModel.slashing.maxEvidenceAge = 1
        genesis.economicModel.gov.versionProposalVoteDurationSeconds = 2000
        genesis.economicModel.slashing.slashBlocksReward = 0
        new_genesis_env.set_genesis(genesis.to_dict())
        new_genesis_env.deploy_all()
        pip = client_noc_list_obj[0].pip
        pip_test = client_noc_list_obj[1].pip

        create_lockup_plan(client_noc_list_obj[0])
        address = pip.node.staking_address
        submitvpandvote([client_noc_list_obj[0]], votingrounds=14)
        shares = client_noc_list_obj[0].staking.get_staking_amount()
        log.info('Node staking amount : {}'.format(shares))
        pip.node.stop()
        self.verify_amount(pip_test, address, shares)

    @pytest.mark.P2
    @allure.title('Node be slashed, verify unstake function')
    def test_UNS_PU_012(self, new_genesis_env, client_noc_list_obj):
        genesis = from_dict(data_class=Genesis, data=new_genesis_env.genesis_config)
        genesis.economicModel.staking.unStakeFreezeDuration = 2
        genesis.economicModel.slashing.maxEvidenceAge = 1
        genesis.economicModel.gov.textProposalVoteDurationSeconds = 480
        genesis.economicModel.slashing.slashBlocksReward = 0
        new_genesis_env.set_genesis(genesis.to_dict())
        new_genesis_env.deploy_all()
        pip = client_noc_list_obj[0].pip
        pip_test = client_noc_list_obj[1].pip

        create_lockup_plan(client_noc_list_obj[0])
        address = pip.node.staking_address
        submittpandvote([client_noc_list_obj[0]], 1)
        shares = client_noc_list_obj[0].staking.get_staking_amount()
        log.info('Node staking amount : {}'.format(shares))
        pip.node.stop()
        self.verify_amount(pip_test, address, shares)

    @pytest.mark.P2
    @allure.title('Node be slashed, verify unstake function')
    def test_UNS_PU_014(self, new_genesis_env, client_noc_list_obj):
        genesis = from_dict(data_class=Genesis, data=new_genesis_env.genesis_config)
        genesis.economicModel.staking.unStakeFreezeDuration = 2
        genesis.economicModel.slashing.maxEvidenceAge = 1
        genesis.economicModel.gov.paramProposalVoteDurationSeconds = 640
        genesis.economicModel.slashing.slashBlocksReward = 0
        new_genesis_env.set_genesis(genesis.to_dict())
        new_genesis_env.deploy_all()
        pip = client_noc_list_obj[0].pip
        pip_test = client_noc_list_obj[1].pip

        create_lockup_plan(client_noc_list_obj[0])
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
        shares = client_noc_list_obj[0].staking.get_staking_amount()
        log.info('Node staking amount : {}'.format(shares))
        pip.node.stop()
        self.verify_amount(pip_test, address, shares)

    @pytest.mark.P2
    @allure.title('Node be slashed, verify unstake function')
    def test_UNS_PU_020(self, new_genesis_env, client_noc_list_obj):
        genesis = from_dict(data_class=Genesis, data=new_genesis_env.genesis_config)
        genesis.economicModel.staking.unStakeFreezeDuration = 2
        genesis.economicModel.slashing.maxEvidenceAge = 1
        genesis.economicModel.gov.paramProposalVoteDurationSeconds = 480
        genesis.economicModel.slashing.slashBlocksReward = 0
        new_genesis_env.set_genesis(genesis.to_dict())
        new_genesis_env.deploy_all()
        pip = client_noc_list_obj[0].pip
        pip_test = client_noc_list_obj[1].pip

        create_lockup_plan(client_noc_list_obj[0])
        address = pip.node.staking_address
        submitppandvote([client_noc_list_obj[0]], 3)
        shares = client_noc_list_obj[0].staking.get_staking_amount()
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