from tests.lib.utils import upload_platon, assert_code, get_pledge_list, wait_block_number
from common.log import log
from tests.lib.client import Client, get_client_obj, get_client_obj_list, StakingConfig
import pytest, time
from tests.govern.conftest import version_proposal_vote

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
    return client_noc_list_obj[0]


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
        assert_code(result, 301004)

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
        self.preactiveproposal(client_list_obj)
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
        assert_code(result, 301004)

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
        result = replace_platon_and_staking(pip_obj, pip_obj.cfg.PLATON_NEW_BIN2)
        assert_code(result, 0)

    def test_ST_UPG_002(self, new_genesis_env, new_node_no_proposal, client_list_obj):
        pip_obj = new_node_no_proposal
        self.upgraded_proposal(client_list_obj)
        result = replace_platon_and_staking(pip_obj, pip_obj.cfg.PLATON_NEW_BIN2)
        assert_code(result, 301004)

    def test_ST_UPG_003(self, new_genesis_env, new_node_no_proposal, client_list_obj):
        pip_obj = new_node_no_proposal
        self.upgraded_proposal(client_list_obj)
        result = replace_platon_and_staking(pip_obj, pip_obj.cfg.PLATON_NEW_BIN0)
        assert_code(result, 0)

    def test_ST_UPG_004(self, new_genesis_env, new_node_no_proposal, client_list_obj):
        pip_obj = new_node_no_proposal
        self.upgraded_proposal(client_list_obj)
        result = replace_platon_and_staking(pip_obj, pip_obj.cfg.PLATON_NEW_BIN3)
        assert_code(result, 0)

    def test_ST_UPG_005(self, new_genesis_env, new_node_no_proposal, client_list_obj):
        pip_obj = new_node_no_proposal
        self.upgraded_proposal(client_list_obj)
        result = replace_platon_and_staking(pip_obj, pip_obj.cfg.PLATON_NEW_BIN)
        assert_code(result, 301004)




