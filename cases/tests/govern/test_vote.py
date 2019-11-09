import pytest
from common.log import log
import time
from tests.lib.utils import assert_code, wait_block_number, upload_platon
from tests.lib.client import get_client_obj
from tests.govern.conftest import version_proposal_vote, get_refund_to_account_block, param_proposal_vote


def text_proposal_vote(pip_obj):
    proposalinfo = pip_obj.get_effect_proposal_info_of_vote(pip_obj.cfg.text_proposal)
    log.info('proposalinfo: {}'.format(proposalinfo))
    result = pip_obj.vote(pip_obj.node.node_id, proposalinfo.get('ProposalID'), pip_obj.cfg.vote_option_yeas,
                          pip_obj.node.staking_address, transaction_cfg=pip_obj.cfg.transaction_cfg)
    log.info('Node {} vote text proposal result {}'.format(pip_obj.node.node_id, result))
    return result

def cancel_param_proposal_vote(pip_obj):
    proposalinfo = pip_obj.get_effect_proposal_info_of_vote(pip_obj.cfg.cancel_proposal)
    log.info('proposalinfo: {}'.format(proposalinfo))
    result = pip_obj.vote(pip_obj.node.node_id, proposalinfo.get('ProposalID'), pip_obj.cfg.vote_option_yeas,
                          pip_obj.node.staking_address, transaction_cfg=pip_obj.cfg.transaction_cfg)
    log.info('Node {} vote cancel proposal result {}'.format(pip_obj.node.node_id, result))
    return result

@pytest.fixture()
def voting_proposal_ve_pipobj(global_test_env, client_verifier_obj):
    pip_obj = client_verifier_obj.pip
    if pip_obj.chain_version != pip_obj.cfg.version0:
        log.info('The chain has been upgraded,restart!')
        global_test_env.deploy_all()
    if pip_obj.is_exist_effective_proposal:
        if pip_obj.is_exist_effective_proposal_for_vote():
            proposalinfo = pip_obj.get_effect_proposal_info_of_vote()
            log.info('get version proposalinfo : {}'.format(proposalinfo))
            if proposalinfo.get('EndVotingBlock') - pip_obj.node.block_number > pip_obj.economic.consensus_size * 2:
                return pip_obj
        global_test_env.deploy_all()
    result = pip_obj.submitVersion(pip_obj.node.node_id, str(time.time()), pip_obj.cfg.version5, 10, pip_obj.node.staking_address,
                          transaction_cfg=pip_obj.cfg.transaction_cfg)
    log.info('node {} submit version proposal {}'.format(pip_obj.node.node_id, result))
    assert_code(result, 0)
    return pip_obj

@pytest.fixture()
def voting_proposal_te_pipobj(global_test_env, client_verifier_obj):
    pip_obj = client_verifier_obj.pip
    if pip_obj.is_exist_effective_proposal_for_vote(pip_obj.cfg.text_proposal):
        proposalinfo = pip_obj.get_effect_proposal_info_of_vote()
        if proposalinfo.get('EndVotingBlock') - pip_obj.node.eth.blockNumber > 2 * pip_obj.economic.consensus_size:
            return pip_obj
        else:
            global_test_env.deploy_all()
    result = pip_obj.submitText(pip_obj.node.node_id, str(time.time()), pip_obj.node.staking_address,
                       transaction_cfg=pip_obj.cfg.transaction_cfg)
    log.info('Submit text proposal result {}'.format(result))
    assert_code(result, 0)
    return pip_obj

def test_V0_VO_001_V0_RE_001_V0_WA_001_V_STA_1_V_OP_1_V_OP_2(voting_proposal_ve_pipobj):
    pip_obj = voting_proposal_ve_pipobj
    proposalinfo = pip_obj.get_effect_proposal_info_of_vote()
    log.info('Get version proposalinfo {}'.format(proposalinfo))

    upload_platon(pip_obj.node, pip_obj.cfg.PLATON_NEW_BIN)
    log.info('Replace the node platon package to {}'.format(pip_obj.cfg.version5))
    pip_obj.node.restart()
    log.info('Restart the node {}'.format(pip_obj.node.node_id))
    result = pip_obj.vote(pip_obj.node.node_id, proposalinfo.get('ProposalID'), pip_obj.cfg.vote_option_nays,
                          pip_obj.node.staking_address, transaction_cfg=pip_obj.cfg.transaction_cfg)
    log.info('vote option {} result {}'.format(pip_obj.cfg.vote_option_nays, result))
    assert_code(result, 302002)

    result = pip_obj.vote(pip_obj.node.node_id, proposalinfo.get('ProposalID'), pip_obj.cfg.vote_option_Abstentions,
                          pip_obj.node.staking_address, transaction_cfg=pip_obj.cfg.transaction_cfg)
    log.info('vote option {} result {}'.format(pip_obj.cfg.vote_option_Abstentions, result))
    assert_code(result, 302002)

    result = pip_obj.vote(pip_obj.node.node_id, proposalinfo.get('ProposalID'), 0,
                          pip_obj.node.staking_address, transaction_cfg=pip_obj.cfg.transaction_cfg)
    log.info('vote option {} result {}'.format(0, result))
    assert_code(result, 302002)

    result = pip_obj.vote(pip_obj.node.node_id, proposalinfo.get('ProposalID'), 'a',
                          pip_obj.node.staking_address, transaction_cfg=pip_obj.cfg.transaction_cfg)
    log.info('vote option {} result {}'.format(pip_obj.cfg.vote_option_nays, result))
    assert_code(result, 302002)

    address, _ = pip_obj.economic.account.generate_account(pip_obj.node.web3, 10**18 * 10000)
    result = pip_obj.vote(pip_obj.node.node_id, proposalinfo.get('ProposalID'), pip_obj.cfg.vote_option_yeas,
                          address, transaction_cfg=pip_obj.cfg.transaction_cfg)
    assert_code(result, 302021)

    result = version_proposal_vote(pip_obj)
    assert_code(result, 0)

    result = pip_obj.vote(pip_obj.node.node_id, proposalinfo.get('ProposalID'), pip_obj.cfg.vote_option_yeas,
                 pip_obj.node.staking_address, transaction_cfg=pip_obj.cfg.transaction_cfg)
    log.info('vote duplicated result: {}'.format(result))
    assert_code(result, 302027)

    wait_block_number(pip_obj.node, proposalinfo.get('EndVotingBlock'))
    result = pip_obj.vote(pip_obj.node.node_id, proposalinfo.get('ProposalID'), pip_obj.cfg.vote_option_yeas,
                 pip_obj.node.staking_address, transaction_cfg=pip_obj.cfg.transaction_cfg)
    log.info('endblock vote result: {}'.format(result))
    assert_code(result, 302026)

def test_V0_WA_001_V0_VO_003_V_STA_9_V_STA_10(voting_proposal_te_pipobj):
    pip_obj = voting_proposal_te_pipobj
    proposalinfo = pip_obj.get_effect_proposal_info_of_vote(pip_obj.cfg.text_proposal)
    address, _ = pip_obj.economic.account.generate_account(pip_obj.node.web3, 10**18 * 10000)
    result = pip_obj.vote(pip_obj.node.node_id, proposalinfo.get('ProposalID'), pip_obj.cfg.vote_option_yeas, address,
                 transaction_cfg=pip_obj.cfg.transaction_cfg)
    log.info('Not staking address vote result {}'.format(result))
    assert_code(result, 302021)

    wait_block_number(pip_obj.node, proposalinfo.get('EndVotingBlock') - 10)
    result = text_proposal_vote(pip_obj)
    log.info('vote result {}'.format(result))
    assert_code(result, 0)

    wait_block_number(pip_obj.node, proposalinfo.get('EndVotingBlock'))
    result = pip_obj.vote(pip_obj.node.node_id, proposalinfo.get('ProposalID'), pip_obj.cfg.vote_option_nays,
                          address, transaction_cfg=pip_obj.cfg.transaction_cfg)
    log.info('Endvoting block vote result {}'.format(result))
    assert_code(result, 302026)

class TestVoteNodeExceptionVP():
    def test_V_VE_1_V_VE_2(self, voting_proposal_ve_pipobj, client_list_obj):
        pip_obj = voting_proposal_ve_pipobj
        client_obj = get_client_obj(pip_obj.node.node_id, client_list_obj)
        address = client_obj.node.staking_address
        proposalinfo = pip_obj.get_effect_proposal_info_of_vote()
        log.info('proposalinfo {}'.format(proposalinfo))
        result = client_obj.staking.withdrew_staking(address)
        log.info('Node {} withdrew staking result : {}'.format(client_obj.node.node_id, result))
        assert_code(result, 0)

        result = version_proposal_vote(pip_obj)
        log.info('node vote result {}'.format(result))
        assert_code(result, 302020)

        pip_obj.economic.wait_settlement_blocknum(pip_obj.node, pip_obj.economic.unstaking_freeze_ratio)
        result = pip_obj.vote(pip_obj.node.node_id, proposalinfo.get('ProposalID'), pip_obj.cfg.vote_option_nays,
                     address, program_version=pip_obj.node.program_version, version_sign=pip_obj.node.program_version_sign,
                     transaction_cfg=pip_obj.cfg.transaction_cfg)
        log.info('Exited node vote result {}'.format(result))
        assert_code(result, 302022)

class TestVoteNodeExceptionTP():
    def test_V0_TE_001_V0_TE_002(self, voting_proposal_te_pipobj, client_list_obj):
        pip_obj = voting_proposal_te_pipobj
        client_obj = get_client_obj(pip_obj.node.node_id, client_list_obj)
        result = client_obj.staking.withdrew_staking(client_obj.node.staking_address)
        log.info('Node {} withdrew staking result {}'.format(pip_obj.node.node_id, result))
        assert_code(result, 0)
        address = client_obj.node.staking_address
        proposalinfo = pip_obj.get_effect_proposal_info_of_vote()
        log.info('proposalinfo {}'.format(proposalinfo))

        result = text_proposal_vote(pip_obj)
        log.info('node vote result {}'.format(result))
        assert_code(result, 302020)

        pip_obj.economic.wait_settlement_blocknum(pip_obj.node, pip_obj.economic.unstaking_freeze_ratio)
        result = pip_obj.vote(pip_obj.node.node_id, proposalinfo.get('ProposalID'), pip_obj.cfg.vote_option_nays,
                     address, program_version=pip_obj.node.program_version, version_sign=pip_obj.node.program_version_sign,
                     transaction_cfg=pip_obj.cfg.transaction_cfg)
        log.info('Exited node vote result {}'.format(result))
        assert_code(result, 302022)

class TestVoteCancelParam():
    def test_PP_VO_001_PP_VO_005_PP_VO_015_PP_VO_017(self, submit_cancel_param):
        pip_obj = submit_cancel_param
        address, _ = pip_obj.economic.account.generate_account(pip_obj.node.web3, 10**18 * 10000)
        proposalinfo = pip_obj.get_effect_proposal_info_of_vote(pip_obj.cfg.cancel_proposal)
        log.info('Cancel proposal info : {}'.format(proposalinfo))
        result = pip_obj.vote(pip_obj.node.node_id, proposalinfo.get('ProposalID'), pip_obj.cfg.vote_option_Abstentions,
                              address, transaction_cfg=pip_obj.cfg.transaction_cfg)
        log.info('Not staking address , node {}, vote cancel proposal result {}'.format(pip_obj.node.node_id, result))
        assert_code(result, 302021)
        result = cancel_param_proposal_vote(pip_obj)
        assert_code(result, 0)
        result = cancel_param_proposal_vote(pip_obj)
        assert_code(result, 302027)

    def test_PP_VO_009_PP_VO_010(self, submit_cancel_param, client_list_obj):
        pip_obj = submit_cancel_param
        address = pip_obj.node.staking_address
        client_obj = get_client_obj(pip_obj.node.node_id, client_list_obj)
        result = client_obj.staking.withdrew_staking(pip_obj.node.staking_address)
        endblock = get_refund_to_account_block(pip_obj)
        log.info('Node {} withdrew staking result {}'.format(pip_obj.node.node_id, result))
        assert_code(result, 0)
        proposalinfo = pip_obj.get_effect_proposal_info_of_vote(pip_obj.cfg.cancel_proposal)
        log.info('Cancel proposal info : {}'.format(proposalinfo))
        result = cancel_param_proposal_vote(pip_obj)
        assert_code(result, 302020)
        wait_block_number(pip_obj.node, endblock)
        result = pip_obj.vote(pip_obj.node.node_id, proposalinfo.get('ProposalID'), pip_obj.cfg.vote_option_nays, address,
                              transaction_cfg=pip_obj.cfg.transaction_cfg)
        assert_code(result, 302022)

class TestVoteParam():
    def test_PP_VO_002_PP_VO_008_PP_VO_015_PP_VO_017(self, submit_param):
        pip_obj = submit_param
        address, _ = pip_obj.economic.account.generate_account(pip_obj.node.web3, 10**18 * 10000)
        proposalinfo = pip_obj.get_effect_proposal_info_of_vote(pip_obj.cfg.param_proposal)
        log.info('param proposal info : {}'.format(proposalinfo))
        result = pip_obj.vote(pip_obj.node.node_id, proposalinfo.get('ProposalID'), pip_obj.cfg.vote_option_Abstentions,
                              address, transaction_cfg=pip_obj.cfg.transaction_cfg)
        log.info('Not staking address , node {}, vote param proposal result {}'.format(pip_obj.node.node_id, result))
        assert_code(result, 302021)
        result = param_proposal_vote(pip_obj)
        assert_code(result, 0)
        result = param_proposal_vote(pip_obj)
        assert_code(result, 302027)

    def test_PP_VO_009_PP_VO_010(self, submit_param, client_list_obj):
        pip_obj = submit_param
        address = pip_obj.node.staking_address
        client_obj = get_client_obj(pip_obj.node.node_id, client_list_obj)
        result = client_obj.staking.withdrew_staking(pip_obj.node.staking_address)
        endblock = get_refund_to_account_block(pip_obj)
        log.info('Node {} withdrew staking result {}'.format(pip_obj.node.node_id, result))
        assert_code(result, 0)
        proposalinfo = pip_obj.get_effect_proposal_info_of_vote(pip_obj.cfg.param_proposal)
        log.info('Cancel proposal info : {}'.format(proposalinfo))
        result = param_proposal_vote(pip_obj)
        assert_code(result, 302020)
        wait_block_number(pip_obj.node, endblock)
        result = pip_obj.vote(pip_obj.node.node_id, proposalinfo.get('ProposalID'), pip_obj.cfg.vote_option_nays, address,
                              transaction_cfg=pip_obj.cfg.transaction_cfg)
        assert_code(result, 302022)

def test_PP_VO_003_PP_VO_004_VS_EP_002_VS_EP_003(submit_param):
    pip_obj = submit_param
    proposalinfo = pip_obj.get_effect_proposal_info_of_vote(pip_obj.cfg.param_proposal)
    log.info('param proposal info {}'.format(proposalinfo))
    wait_block_number(pip_obj.node, proposalinfo.get('EndVotingBlock') - 5)
    result = param_proposal_vote(pip_obj)
    assert_code(result, 0)
    result = pip_obj.pip.getTallyResult(proposalinfo.get('ProposalID'))
    log.info('Interface getTallyResult result is {}'.format(result))
    assert_code(result, 302030)
    wait_block_number(pip_obj.node, proposalinfo.get('EndVotingBlock'))
    result = pip_obj.vote(pip_obj.node.node_id, proposalinfo.get('ProposalID'), pip_obj.cfg.vote_option_yeas,
                          pip_obj.node.staking_address, transaction_cfg=pip_obj.cfg.transaction_cfg)
    log.info('Node {} vote param proposal result : {]'.format(pip_obj.node.node_id, result))
    result = pip_obj.pip.getTallyResult(proposalinfo.get('ProposalID'))
    log.info('Interface getTallyResult result is {}'.format(result))
    assert_code(result, 0)

def test_PP_VO_001_PP_VO_006_PP_VO_007_VS_EP_001(submit_cancel_param):
    pip_obj = submit_cancel_param
    proposalinfo = pip_obj.get_effect_proposal_info_of_vote(pip_obj.cfg.cancel_proposal)
    log.info('cancel proposal info {}'.format(proposalinfo))
    wait_block_number(pip_obj.node, proposalinfo.get('EndVotingBlock') - 8)
    result = cancel_param_proposal_vote(pip_obj)
    assert_code(result, 0)
    result = pip_obj.pip.getTallyResult(proposalinfo.get('ProposalID'))
    log.info('Interface getTallyResult result is {}'.format(result))
    assert_code(result, 302030)
    wait_block_number(pip_obj.node, proposalinfo.get('EndVotingBlock'))
    result = pip_obj.vote(pip_obj.node.node_id, proposalinfo.get('ProposalID'), pip_obj.cfg.vote_option_yeas,
                          pip_obj.node.staking_address, transaction_cfg=pip_obj.cfg.transaction_cfg)
    log.info('Node {} vote cancel proposal result : {]'.format(pip_obj.node.node_id, result))
    result = pip_obj.pip.getTallyResult(proposalinfo.get('ProposalID'))
    log.info('Interface getTallyResult result is {}'.format(result))
    assert_code(result, 0)
