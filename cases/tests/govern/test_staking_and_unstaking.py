from tests.lib.utils import upload_platon, assert_code, get_pledge_list, wait_block_number
from common.log import log
from tests.lib.client import Client, get_client_obj, get_client_obj_list, StakingConfig
import pytest, time, math
from tests.govern.conftest import version_proposal_vote, proposal_vote
from tests.lib import Genesis
from dacite import from_dict
from tests.govern.test_voting_statistics import submitcvpandvote, submitcppandvote


@pytest.fixture()
def new_node_no_proposal(no_vp_proposal, client_noc_list_obj, client_list_obj):
    pip_obj = no_vp_proposal
    client_obj = get_client_obj(pip_obj.node.node_id, client_list_obj)
    candidate_list = get_pledge_list(client_obj.ppos.getCandidateList)
    log.info('candidate_list: {}'.format(candidate_list))
    for client_obj in client_noc_list_obj:
        if client_obj.node.node_id not in candidate_list:
            return client_obj.pip
    log.info('All nodes are staked, restart the chain')
    pip_obj.economic.env.deploy_all()
    return client_noc_list_obj[0].pip


def replace_platon_and_staking(pip_obj, bin):
    node_obj_list = pip_obj.economic.env.get_all_nodes()
    client_list_obj = []
    for node_obj in node_obj_list:
        client_list_obj.append(Client(pip_obj.economic.env, node_obj, StakingConfig("externalId", "nodeName", "website",
                                                                                    "details")))
    client_obj = get_client_obj(pip_obj.node.node_id, client_list_obj)
    upload_platon(pip_obj.node, bin)
    log.info('Replace the platon of the node {}'.format(pip_obj.node.node_id))
    pip_obj.node.restart()
    log.info('Restart the node {}'.format(pip_obj.node.node_id))
    address, _ = pip_obj.economic.account.generate_account(pip_obj.node.web3, 10**18 * 10000000)
    result = client_obj.staking.create_staking(0, address, address, transaction_cfg=pip_obj.cfg.transaction_cfg)
    log.info('Node {} staking result {}'.format(pip_obj.node.node_id, result))
    return result


class TestVotingProposalStaking():
    def test_ST_VS_001(self, new_node_has_proposal):
        pip_obj = new_node_has_proposal
        result = replace_platon_and_staking(pip_obj, pip_obj.cfg.PLATON_NEW_BIN2)
        assert_code(result, 0)

    def test_ST_VS_002(self, new_node_has_proposal):
        pip_obj = new_node_has_proposal
        result = replace_platon_and_staking(pip_obj, pip_obj.cfg.PLATON_NEW_BIN1)
        assert_code(result, 301004)

    def test_ST_VS_003(self, new_node_has_proposal):
        pip_obj = new_node_has_proposal
        result = replace_platon_and_staking(pip_obj, pip_obj.cfg.PLATON_NEW_BIN0)
        assert_code(result, 0)

    def test_ST_VS_004(self, new_node_has_proposal):
        pip_obj = new_node_has_proposal
        result = replace_platon_and_staking(pip_obj, pip_obj.cfg.PLATON_NEW_BIN3)
        assert_code(result, 0)

    def test_ST_VS_005(self, new_node_has_proposal):
        pip_obj = new_node_has_proposal
        result = replace_platon_and_staking(pip_obj, pip_obj.cfg.PLATON_NEW_BIN)
        assert_code(result, 0)

class TestNoProposalStaking():
    def test_ST_NO_001(self, new_node_no_proposal):
        pip_obj = new_node_no_proposal
        result = replace_platon_and_staking(pip_obj, pip_obj.cfg.PLATON_NEW_BIN2)
        assert_code(result, 0)

    def test_ST_NO_002(self, new_node_no_proposal):
        pip_obj = new_node_no_proposal
        result = replace_platon_and_staking(pip_obj, pip_obj.cfg.PLATON_NEW_BIN1)
        assert_code(result, 301004)

    def test_ST_NO_003(self, new_node_no_proposal):
        pip_obj = new_node_no_proposal
        result = replace_platon_and_staking(pip_obj, pip_obj.cfg.PLATON_NEW_BIN0)
        assert_code(result, 0)

    def test_ST_NO_004(self, new_node_no_proposal):
        pip_obj = new_node_no_proposal
        result = replace_platon_and_staking(pip_obj, pip_obj.cfg.PLATON_NEW_BIN3)
        assert_code(result, 0)

    def test_ST_NO_005(self, new_node_no_proposal):
        pip_obj = new_node_no_proposal
        result = replace_platon_and_staking(pip_obj, pip_obj.cfg.PLATON_NEW_BIN)
        assert_code(result, 301005)

class TestPreactiveProposalStaking():
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
        for pip_obj in pip_list_obj:
            result = version_proposal_vote(pip_obj)
            assert_code(result, 0)
        wait_block_number(pip_obj.node, proposalinfo.get('EndVotingBlock'))
        assert pip_obj.get_status_of_proposal(proposalinfo.get('ProposalID')) == 4

    def test_ST_PR_001(self, new_genesis_env, new_node_no_proposal, client_list_obj):
        pip_obj = new_node_no_proposal
        self.preactive_proposal(client_list_obj)
        result = replace_platon_and_staking(pip_obj, pip_obj.cfg.PLATON_NEW_BIN2)
        assert_code(result, 301004)

    def test_ST_PR_002(self, new_genesis_env, new_node_no_proposal, client_list_obj):
        pip_obj = new_node_no_proposal
        self.preactive_proposal(client_list_obj)
        result = replace_platon_and_staking(pip_obj, pip_obj.cfg.PLATON_NEW_BIN1)
        assert_code(result, 301004)

    def test_ST_PR_003(self, new_genesis_env, new_node_no_proposal, client_list_obj):
        pip_obj = new_node_no_proposal
        self.preactive_proposal(client_list_obj)
        result = replace_platon_and_staking(pip_obj, pip_obj.cfg.PLATON_NEW_BIN0)
        assert_code(result, 301004)

    def test_ST_PR_004(self, new_genesis_env, new_node_no_proposal, client_list_obj):
        pip_obj = new_node_no_proposal
        self.preactive_proposal(client_list_obj)
        result = replace_platon_and_staking(pip_obj, pip_obj.cfg.PLATON_NEW_BIN3)
        assert_code(result, 301004)

    def test_ST_PR_005(self, new_genesis_env, new_node_no_proposal, client_list_obj):
        pip_obj = new_node_no_proposal
        self.preactive_proposal(client_list_obj)
        result = replace_platon_and_staking(pip_obj, pip_obj.cfg.PLATON_NEW_BIN)
        assert_code(result, 0)

    def test_ST_PR_006(self, new_genesis_env, new_node_no_proposal, client_list_obj):
        pip_obj = new_node_no_proposal
        self.preactive_proposal(client_list_obj)
        result = replace_platon_and_staking(pip_obj, pip_obj.cfg.PLATON_NEW_BIN8)
        assert_code(result, 301005)

class TestUpgradedProposalStaking():
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
        for pip_obj in pip_list_obj:
            result = version_proposal_vote(pip_obj)
            assert_code(result, 0)
        wait_block_number(pip_obj.node, proposalinfo.get('EndVotingBlock'))
        assert pip_obj.get_status_of_proposal(proposalinfo.get('ProposalID')) == 4
        wait_block_number(pip_obj.node, proposalinfo.get('ActiveBlock'))
        assert pip_obj.get_status_of_proposal(proposalinfo.get('ProposalID')) == 5

    def test_ST_UPG_001(self, new_genesis_env, new_node_no_proposal, client_list_obj):
        pip_obj = new_node_no_proposal
        self.upgraded_proposal(client_list_obj)
        result = replace_platon_and_staking(pip_obj, pip_obj.cfg.PLATON_NEW_BIN4)
        assert_code(result, 0)

    def test_ST_UPG_002(self, new_genesis_env, new_node_no_proposal, client_list_obj):
        pip_obj = new_node_no_proposal
        self.upgraded_proposal(client_list_obj)
        result = replace_platon_and_staking(pip_obj, pip_obj.cfg.PLATON_NEW_BIN0)
        assert_code(result, 301004)

    def test_ST_UPG_003(self, new_genesis_env, new_node_no_proposal, client_list_obj):
        pip_obj = new_node_no_proposal
        self.upgraded_proposal(client_list_obj)
        result = replace_platon_and_staking(pip_obj, pip_obj.cfg.PLATON_NEW_BIN)
        assert_code(result, 0)

    def test_ST_UPG_004(self, new_genesis_env, new_node_no_proposal, client_list_obj):
        pip_obj = new_node_no_proposal
        self.upgraded_proposal(client_list_obj)
        result = replace_platon_and_staking(pip_obj, pip_obj.cfg.PLATON_NEW_BIN6)
        assert_code(result, 0)

    def test_ST_UPG_005(self, new_genesis_env, new_node_no_proposal, client_list_obj):
        pip_obj = new_node_no_proposal
        self.upgraded_proposal(client_list_obj)
        result = replace_platon_and_staking(pip_obj, pip_obj.cfg.PLATON_NEW_BIN7)
        assert_code(result, 301005)


class TestUnstaking():
    def test_UNS_AM_003_007(self, new_genesis_env, client_verifier_obj):
        genesis = from_dict(data_class=Genesis, data=new_genesis_env.genesis_config)
        genesis.economicModel.staking.unStakeFreezeDuration == 2
        genesis.economicModel.slashing.maxEvidenceAge == 1
        new_genesis_env.set_genesis(genesis.to_dict())
        new_genesis_env.deploy_all()
        pip_obj = client_verifier_obj.pip
        address = pip_obj.node.staking_address
        list_obj = [client_verifier_obj]
        submitcvpandvote(list_obj, 1)
        result = version_proposal_vote(pip_obj)
        assert_code(result, 0)
        shares = client_verifier_obj.staking.get_staking_amount(pip_obj.node)
        result = client_verifier_obj.staking.withdrew_staking(address)
        log.info('Node withdrew staking result : {}'.format(result))
        assert_code(result, 0)
        calculated_block = 480
        wait_block_number(pip_obj.node, calculated_block)
        balance_before = pip_obj.node.eth.getBalance(address, calculated_block-1)
        log.info('Block bumber {} staking address balance {}'.format(calculated_block-1, balance_before))
        balance_after = pip_obj.node.eth.getBalance(address, calculated_block)
        log.info('Block bumber {} staking address balance {}'.format(calculated_block, balance_after))
        log.info('{}'.format(pip_obj.economic.get_current_year_reward(pip_obj.node)))
        assert balance_after - balance_before == shares

    def test_UNS_AM_005(self, new_genesis_env, client_verifier_obj):
        genesis = from_dict(data_class=Genesis, data=new_genesis_env.genesis_config)
        genesis.economicModel.staking.unStakeFreezeDuration == 2
        genesis.economicModel.slashing.maxEvidenceAge == 1
        new_genesis_env.set_genesis(genesis.to_dict())
        new_genesis_env.deploy_all()
        pip_obj = client_verifier_obj.pip
        address = pip_obj.node.staking_address
        list_obj = [client_verifier_obj]
        result = pip_obj.submitText(pip_obj.node.node_id, str(time.time()), address, transaction_cfg=pip_obj.cfg.transaction_cfg)
        assert_code(result, 0)

        submitcppandvote(list_obj, 1)
        result = proposal_vote(pip_obj, proposaltype=pip_obj.cfg.text_proposal)
        assert_code(result, 0)
        shares = client_verifier_obj.staking.get_staking_amount(pip_obj.node)
        result = client_verifier_obj.staking.withdrew_staking(address)
        log.info('Node withdrew staking result : {}'.format(result))
        assert_code(result, 0)
        calculated_block = 480
        wait_block_number(pip_obj.node, calculated_block)
        balance_before = pip_obj.node.eth.getBalance(address, calculated_block-1)
        log.info('Block bumber {} staking address balance {}'.format(calculated_block-1, balance_before))
        balance_after = pip_obj.node.eth.getBalance(address, calculated_block)
        log.info('Block bumber {} staking address balance {}'.format(calculated_block, balance_after))
        assert balance_after - balance_before == shares

    def test_UNS_AM_004_006_008(self, new_genesis_env, client_verifier_obj_list):
        genesis = from_dict(data_class=Genesis, data=new_genesis_env.genesis_config)
        genesis.economicModel.staking.unStakeFreezeDuration == 2
        genesis.economicModel.slashing.maxEvidenceAge == 1
        genesis.economicModel.gov.versionProposalVoteDurationSeconds = 1000
        genesis.economicModel.gov.textProposalVoteDurationSeconds = 840
        new_genesis_env.set_genesis(genesis.to_dict())
        new_genesis_env.deploy_all()
        pip_obj_one = client_verifier_obj_list[0].pip
        pip_obj_two = client_verifier_obj_list[1].pip
        pip_obj_three = client_verifier_obj_list[2].pip
        address = pip_obj_one.node.staking_address
        result = pip_obj_one.submitVersion(pip_obj_one.node.node_id, str(time.time()), pip_obj_one.cfg.version5, 17, address,
                                       transaction_cfg=pip_obj_one.cfg.transaction_cfg)
        log.info('Submit version proposal result : {}'.format(result))
        proposalinfo_version = pip_obj_one.get_effect_proposal_info_of_vote(pip_obj_one.cfg.version_proposal)
        log.info('Version proposal information : {}'.format(proposalinfo_version))

        result = pip_obj_one.submitCancel(pip_obj_one.node.node_id, str(time.time()), 13, proposalinfo_version.get('ProposalID'),
                                      address, transaction_cfg=pip_obj_one.cfg.transaction_cfg)
        log.info('Submit cancel proposal result : {}'.format(result))
        assert_code(result, 0)
        result_text = pip_obj_one.submitText(pip_obj_one.node.node_id, str(time.time()), address,
                                        transaction_cfg=pip_obj_one.cfg.transaction_cfg)
        log.info('Submit text proposal result : {}'.format(result_text))
        result = proposal_vote(pip_obj_one, proposaltype=pip_obj_one.cfg.cancel_proposal)
        assert_code(result, 0)
        result = version_proposal_vote(pip_obj_two)
        assert_code(result, 0)
        result = proposal_vote(pip_obj_three, proposaltype=pip_obj_three.cfg.text_proposal)
        assert_code(result, 0)
        proposalinfo_cancel = pip_obj_one.get_effect_proposal_info_of_vote(pip_obj_one.cfg.cancel_proposal)
        log.info('Get cancel proposal information {}'.format(proposalinfo_cancel))
        proposalinfo_text = pip_obj_one.get_effect_proposal_info_of_vote(pip_obj_one.cfg.text_proposal)
        log.info('Get cancel proposal information {}'.format(proposalinfo_text))
        shares_one = client_verifier_obj_list[0].staking.get_staking_amount(pip_obj_one.node)
        shares_two = client_verifier_obj_list[1].staking.get_staking_amount(pip_obj_two.node)
        shares_three = client_verifier_obj_list[2].staking.get_staking_amount(pip_obj_three.node)
        result = client_verifier_obj_list[0].staking.withdrew_staking(address)
        log.info('Node {} withdrew staking result : {}'.format(pip_obj_one.node.node_id, result))
        assert_code(result, 0)

        result = client_verifier_obj_list[1].staking.withdrew_staking(address)
        log.info('Node {} withdrew staking result : {}'.format(pip_obj_two.node.node_id, result))
        assert_code(result, 0)

        result = client_verifier_obj_list[2].staking.withdrew_staking(address)
        log.info('Node {} withdrew staking result : {}'.format(pip_obj_three.node.node_id, result))
        assert_code(result, 0)
        calculated_block = 480
        wait_block_number(pip_obj_one.node, calculated_block)
        balance_before = pip_obj_one.node.eth.getBalance(address, calculated_block-1)
        log.info('Block bumber {} staking address balance {}'.format(calculated_block-1, balance_before))
        balance_after = pip_obj_one.node.eth.getBalance(address, calculated_block)
        log.info('Block bumber {} staking address balance {}'.format(calculated_block, balance_after))
        assert balance_after == balance_before

        blocknumber = math.ceil(proposalinfo_cancel.get('EndVotingBlock')/pip_obj_one.economic.settlement_size
                                ) * pip_obj_one.economic.settlement_size
        wait_block_number(pip_obj_one.node, blocknumber)
        balance_before = pip_obj_one.node.eth.getBalance(address, blocknumber-1)
        log.info('Block bumber {} staking address balance {}'.format(blocknumber-1, balance_before))
        balance_after = pip_obj_one.node.eth.getBalance(address, blocknumber)
        log.info('Block bumber {} staking address balance {}'.format(blocknumber, balance_after))
        assert balance_after - balance_before == shares_one

        blocknumber = math.ceil(proposalinfo_version.get('EndVotingBlock')/pip_obj_one.economic.settlement_size
                                ) * pip_obj_one.economic.settlement_size
        wait_block_number(pip_obj_one.node, blocknumber)
        balance_before = pip_obj_one.node.eth.getBalance(address, blocknumber-1)
        log.info('Block bumber {} staking address balance {}'.format(blocknumber-1, balance_before))
        balance_after = pip_obj_one.node.eth.getBalance(address, blocknumber)
        log.info('Block bumber {} staking address balance {}'.format(blocknumber, balance_after))
        assert balance_after - balance_before == shares_two

        blocknumber = math.ceil(proposalinfo_text.get('EndVotingBlock')/pip_obj_one.economic.settlement_size
                                ) * pip_obj_one.economic.settlement_size
        wait_block_number(pip_obj_one.node, blocknumber)
        balance_before = pip_obj_one.node.eth.getBalance(address, blocknumber-1)
        log.info('Block bumber {} staking address balance {}'.format(blocknumber-1, balance_before))
        balance_after = pip_obj_one.node.eth.getBalance(address, blocknumber)
        log.info('Block bumber {} staking address balance {}'.format(blocknumber, balance_after))
        assert balance_after - balance_before == shares_three

    def test_UNS_AM_009(self, new_genesis_env, client_con_list_obj, client_noc_list_obj):
        genesis = from_dict(data_class=Genesis, data=new_genesis_env.genesis_config)
        genesis.economicModel.staking.unStakeFreezeDuration == 2
        genesis.economicModel.slashing.maxEvidenceAge == 1
        genesis.economicModel.gov.versionProposalVoteDurationSeconds = 1000
        genesis.economicModel.gov.textProposalVoteDurationSeconds = 840
        new_genesis_env.set_genesis(genesis.to_dict())
        new_genesis_env.deploy_all()
        pip_obj_test = client_noc_list_obj[0].pip
        address, _ = pip_obj_test.economic.account.generate_account(pip_obj_test.node.web3, 10**18 * 2000000)
        plan = [{'Epoch': 20, 'Amount': 10**18 * 2000000}]
        result = client_noc_list_obj[0].restricting.createRestrictingPlan(address, plan, address,
                                                                          transaction_cfg=pip_obj_test.cfg.transaction_cfg)
        log.info('CreateRestrictingPlan result : {}'.format(result))
        assert_code(result, 0)
        result = client_noc_list_obj[0].staking.create_staking(1, address, address,
                                                               transaction_cfg=pip_obj_test.cfg.transaction_cfg)
        log.info('Create staking result : {}'.format(result))
        assert_code(result, 0)
