import pytest
from common.log import log
import time
from tests.lib.utils import assert_code, wait_block_number, upload_platon
from tests.lib.client import get_client_obj
from tests.govern.conftest import version_proposal_vote, get_refund_to_account_block, proposal_vote
from dacite import from_dict
from tests.lib.genesis import Genesis
from tests.govern.test_voting_statistics import submitvpandvote, submittpandvote, submitcppandvote


def replace_platon_vote(pip_obj, bin=None, program_version=None, version_sign=None):
    '''
    Replace the bin package of the node, restart the node
    :param pip_obj:
    :param bin:
    :return:
    '''
    if bin:
        upload_platon(pip_obj.node, bin)
        pip_obj.node.restart()
    if program_version is None:
        program_version = pip_obj.node.program_version
    if version_sign is None:
        version_sign = pip_obj.node.program_version_sign
    proposalinfo = pip_obj.get_effect_proposal_info_of_vote()
    log.info('Get version proposal information {}'.format(proposalinfo))
    result = pip_obj.vote(pip_obj.node.node_id, proposalinfo.get('ProposalID'), pip_obj.cfg.vote_option_yeas,
                          pip_obj.node.staking_address, program_version=program_version, version_sign=version_sign,
                          transaction_cfg=pip_obj.cfg.transaction_cfg)
    log.info('Node {} vote result : {}'.format(pip_obj.node.node_id, result))
    return result


def text_proposal_vote(pip_obj):
    proposalinfo = pip_obj.get_effect_proposal_info_of_vote(pip_obj.cfg.text_proposal)
    log.info('proposalinfo: {}'.format(proposalinfo))
    result = pip_obj.vote(pip_obj.node.node_id, proposalinfo.get('ProposalID'), pip_obj.cfg.vote_option_yeas,
                          pip_obj.node.staking_address, transaction_cfg=pip_obj.cfg.transaction_cfg)
    log.info('Node {} vote text proposal result {}'.format(pip_obj.node.node_id, result))
    return result


def cancel_proposal_vote(pip_obj):
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


class TestVoteVP():
    @pytest.mark.P1
    def test_V_STA_2_to_5(self, no_vp_proposal, client_verifier_obj_list):
        pip_obj = no_vp_proposal
        result = pip_obj.submitVersion(pip_obj.node.node_id, str(time.time()), pip_obj.cfg.version5, 2,
                                       pip_obj.node.staking_address, transaction_cfg=pip_obj.cfg.transaction_cfg)
        log.info('Submit Version proposal result : {}'.format(result))
        assert_code(result, 0)
        proposalinfo = pip_obj.get_effect_proposal_info_of_vote(pip_obj.cfg.version_proposal)
        log.info('Get version proposal information : {}'.format(proposalinfo))
        for client_obj in client_verifier_obj_list[:2]:
            result = version_proposal_vote(client_obj.pip)
            log.info('Node {} vote proposal result : {}'.format(client_obj.node.node_id, result))
            assert_code(result, 0)
        upload_platon(client_verifier_obj_list[-2].node, pip_obj.cfg.PLATON_NEW_BIN)
        client_verifier_obj_list[-2].node.restart()
        log.info('Replace the platon of the Node {}, restart the node'.format(client_verifier_obj_list[-2].node.node_id))
        wait_block_number(pip_obj.node, proposalinfo.get('EndVotingBlock') - 10)
        result = client_verifier_obj_list[-2].pip.vote(client_verifier_obj_list[-2].pip.node.node_id, proposalinfo.get('ProposalID'),
                                                       client_verifier_obj_list[-2].pip.cfg.vote_option_yeas,
                                                       client_verifier_obj_list[-2].pip.node.staking_address, transaction_cfg=pip_obj.cfg.transaction_cfg)
        log.info('Node {} vote proposal result : {}'.format(client_verifier_obj_list[-1].node.node_id, result))
        assert_code(result, 0)
        wait_block_number(pip_obj.node, proposalinfo.get('EndVotingBlock'))
        pip_obj = client_verifier_obj_list[-1].pip
        result = pip_obj.vote(pip_obj.node.node_id, proposalinfo.get('ProposalID'), pip_obj.cfg.vote_option_yeas,
                              pip_obj.node.staking_address, transaction_cfg=pip_obj.cfg.transaction_cfg)
        log.info('Node {} vote proposal result : {}'.format(client_verifier_obj_list[-1].node.node_id, result))
        assert_code(pip_obj.get_status_of_proposal(proposalinfo.get('ProposalID')), 4)
        wait_block_number(pip_obj.node, proposalinfo.get('ActiveBlock'))
        assert_code(pip_obj.get_status_of_proposal(proposalinfo.get('ProposalID')), 5)
        result = pip_obj.vote(pip_obj.node.node_id, proposalinfo.get('ProposalID'), pip_obj.cfg.vote_option_yeas,
                              pip_obj.node.staking_address, transaction_cfg=pip_obj.cfg.transaction_cfg)
        log.info('Node {} vote proposal result : {}'.format(client_verifier_obj_list[-1].node.node_id, result))


@pytest.mark.compatibility
@pytest.mark.P0
def test_VO_VO_001_V0_RE_001_V0_WA_001_V_STA_1_VO_OP_001_VO_OP_002(no_vp_proposal):
    pip_obj = no_vp_proposal
    result = pip_obj.submitVersion(pip_obj.node.node_id, str(time.time()), pip_obj.cfg.version8, 2,
                                   pip_obj.node.staking_address, transaction_cfg=pip_obj.cfg.transaction_cfg)
    log.info('Submit version proposal result : {}'.format(result))
    assert_code(result, 0)
    proposalinfo = pip_obj.get_effect_proposal_info_of_vote()
    log.info('Get version proposalinfo {}'.format(proposalinfo))

    upload_platon(pip_obj.node, pip_obj.cfg.PLATON_NEW_BIN)
    log.info('Replace the node platon package to {}'.format(pip_obj.cfg.version5))
    pip_obj.node.restart()
    log.info('Restart the node {}'.format(pip_obj.node.node_id))
    address, _ = pip_obj.economic.account.generate_account(pip_obj.node.web3, 10**18 * 10000)
    result = pip_obj.vote(pip_obj.node.node_id, proposalinfo.get('ProposalID'), pip_obj.cfg.vote_option_yeas,
                          address, transaction_cfg=pip_obj.cfg.transaction_cfg)
    log.info('Not staking address vote result : {}'.format(result))
    assert_code(result, 302021)

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


@pytest.mark.P0
def test_VO_VO_003_V_STA_9_V_STA_10_V_STA_11_V0_WA_003_V0_RE_003(voting_proposal_te_pipobj, client_verifier_obj_list):
    pip_obj = voting_proposal_te_pipobj
    proposalinfo = pip_obj.get_effect_proposal_info_of_vote(pip_obj.cfg.text_proposal)
    address, _ = pip_obj.economic.account.generate_account(pip_obj.node.web3, 10**18 * 10000)
    result = pip_obj.vote(pip_obj.node.node_id, proposalinfo.get('ProposalID'), pip_obj.cfg.vote_option_yeas, address,
                          transaction_cfg=pip_obj.cfg.transaction_cfg)
    log.info('Not staking address vote result {}'.format(result))
    assert_code(result, 302021)

    result = text_proposal_vote(pip_obj)
    log.info('vote result {}'.format(result))
    assert_code(result, 0)

    result = text_proposal_vote(pip_obj)
    log.info('Repeat vote  result {}'.format(result))
    assert_code(result, 302027)
    for client_obj in client_verifier_obj_list:
        if client_obj.node.node_id != pip_obj.node.node_id:
            pip_obj_test = client_obj.pip
            break

    wait_block_number(pip_obj.node, proposalinfo.get('EndVotingBlock') - 10)
    result = text_proposal_vote(pip_obj_test)
    log.info('Node {} vote result {}'.format(pip_obj_test.node.node_id, result))
    assert_code(result, 0)

    wait_block_number(pip_obj.node, proposalinfo.get('EndVotingBlock'))
    result = pip_obj.vote(pip_obj.node.node_id, proposalinfo.get('ProposalID'), pip_obj.cfg.vote_option_nays,
                          pip_obj.node.staking_address, transaction_cfg=pip_obj.cfg.transaction_cfg)
    log.info('Endvoting block vote result {}'.format(result))
    assert_code(result, 302026)


class TestVoteNodeException():
    @pytest.mark.P0
    def test_VO_TE_001_002_PP_VO_009_010_PP_VO_011_012_PP_VO_014_VO_TER_008_VO_TER_006(self, new_genesis_env,
                                                                                       client_con_list_obj, client_noconsensus_obj):
        genesis = from_dict(data_class=Genesis, data=new_genesis_env.genesis_config)
        genesis.economicModel.gov.paramProposalVoteDurationSeconds = 640
        genesis.economicModel.staking.unStakeFreezeDuration = 2
        new_genesis_env.set_genesis(genesis.to_dict())
        new_genesis_env.deploy_all()
        pip_obj = client_con_list_obj[0].pip
        result = pip_obj.submitText(pip_obj.node.node_id, str(time.time()), pip_obj.node.staking_address,
                                    transaction_cfg=pip_obj.cfg.transaction_cfg)
        log.info('Submit text proposal information : {}'.format(result))
        assert_code(result, 0)
        proposalinfo_text = pip_obj.get_effect_proposal_info_of_vote(pip_obj.cfg.text_proposal)
        log.info('Get text proposal information : {}'.format(proposalinfo_text))

        result = pip_obj.submitParam(pip_obj.node.node_id, str(time.time()), 'slashing', 'slashBlocksReward',
                                     '123', pip_obj.node.staking_address, transaction_cfg=pip_obj.cfg.transaction_cfg)
        log.info('Submit param proposal information : {}'.format(result))
        assert_code(result, 0)
        proposalinfo_param = pip_obj.get_effect_proposal_info_of_vote(pip_obj.cfg.param_proposal)
        log.info('Get param proposal information : {}'.format(proposalinfo_param))

        result = pip_obj.submitCancel(pip_obj.node.node_id, str(time.time()), 10, proposalinfo_param.get('ProposalID'),
                                      pip_obj.node.staking_address, transaction_cfg=pip_obj.cfg.transaction_cfg)
        log.info('Submit cancel proposal information : {}'.format(result))
        assert_code(result, 0)
        proposalinfo_cancel = pip_obj.get_effect_proposal_info_of_vote(pip_obj.cfg.cancel_proposal)
        log.info('Get cancel proposal information : {}'.format(proposalinfo_cancel))

        result = client_con_list_obj[0].staking.withdrew_staking(client_con_list_obj[0].node.staking_address)
        log.info('Node {} withdrew staking result {}'.format(client_con_list_obj[0].node.node_id, result))
        assert_code(result, 0)
        address = client_con_list_obj[0].node.staking_address

        result = text_proposal_vote(pip_obj)
        log.info('node vote text proposal result {}'.format(result))
        assert_code(result, 302020)

        result = proposal_vote(pip_obj)
        log.info('node vote param proposal result {}'.format(result))
        assert_code(result, 302020)

        result = cancel_proposal_vote(pip_obj)
        log.info('node vote cancel proposal result {}'.format(result))
        assert_code(result, 302020)

        address_test, _ = pip_obj.economic.account.generate_account(pip_obj.node.web3, 10**18 * 10000)
        result = client_noconsensus_obj.pip.vote(client_noconsensus_obj.node.node_id, proposalinfo_text.get('ProposalID'),
                                                 pip_obj.cfg.vote_option_yeas, address_test,
                                                 transaction_cfg=pip_obj.cfg.transaction_cfg)
        log.info('node {} vote text proposal result {}'.format(client_noconsensus_obj.node.node_id, result))
        assert_code(result, 302022)

        result = client_noconsensus_obj.pip.vote(client_noconsensus_obj.node.node_id, proposalinfo_param.get('ProposalID'),
                                                 pip_obj.cfg.vote_option_yeas, address_test,
                                                 transaction_cfg=pip_obj.cfg.transaction_cfg)
        log.info('node {} vote param proposal result {}'.format(client_noconsensus_obj.node.node_id, result))
        assert_code(result, 302022)

        result = client_noconsensus_obj.pip.vote(client_noconsensus_obj.node.node_id, proposalinfo_cancel.get('ProposalID'),
                                                 pip_obj.cfg.vote_option_yeas, address_test,
                                                 transaction_cfg=pip_obj.cfg.transaction_cfg)
        log.info('node {} vote cancel proposal result {}'.format(client_noconsensus_obj.node.node_id, result))
        assert_code(result, 302022)

        pip_obj.economic.wait_settlement_blocknum(pip_obj.node, pip_obj.economic.unstaking_freeze_ratio)
        result = pip_obj.vote(pip_obj.node.node_id, proposalinfo_text.get('ProposalID'), pip_obj.cfg.vote_option_nays,
                              address, transaction_cfg=pip_obj.cfg.transaction_cfg)
        log.info('Exited node vote text proposal result {}'.format(result))
        assert_code(result, 302022)

        result = pip_obj.vote(pip_obj.node.node_id, proposalinfo_param.get('ProposalID'), pip_obj.cfg.vote_option_nays,
                              address, transaction_cfg=pip_obj.cfg.transaction_cfg)
        log.info('Exited node vote param proposal result {}'.format(result))
        assert_code(result, 302022)

        result = pip_obj.vote(pip_obj.node.node_id, proposalinfo_cancel.get('ProposalID'), pip_obj.cfg.vote_option_nays,
                              address, transaction_cfg=pip_obj.cfg.transaction_cfg)
        log.info('Exited node vote cancel proposal result {}'.format(result))
        assert_code(result, 302022)

    @pytest.mark.P0
    def test_VO_VE_001_002_VO_CA_001_002_VO_TER_002_VO_TER_004(self, new_genesis_env, client_con_list_obj, client_noconsensus_obj):
        pip_obj = client_con_list_obj[0].pip
        genesis = from_dict(data_class=Genesis, data=new_genesis_env.genesis_config)
        genesis.economicModel.gov.versionProposalVoteDurationSeconds = 3200
        genesis.economicModel.staking.unStakeFreezeDuration = 2
        new_genesis_env.set_genesis(genesis.to_dict())
        new_genesis_env.deploy_all()
        result = pip_obj.submitVersion(pip_obj.node.node_id, str(time.time()), pip_obj.cfg.version5, 20,
                                       pip_obj.node.staking_address, transaction_cfg=pip_obj.cfg.transaction_cfg)
        log.info('Submit version proposal information : {}'.format(result))
        assert_code(result, 0)
        proposalinfo_version = pip_obj.get_effect_proposal_info_of_vote(pip_obj.cfg.version_proposal)
        log.info('Get param proposal information : {}'.format(proposalinfo_version))

        result = pip_obj.submitCancel(pip_obj.node.node_id, str(time.time()), 10, proposalinfo_version.get('ProposalID'),
                                      pip_obj.node.staking_address, transaction_cfg=pip_obj.cfg.transaction_cfg)
        log.info('Submit cancel proposal information : {}'.format(result))
        assert_code(result, 0)
        proposalinfo_cancel = pip_obj.get_effect_proposal_info_of_vote(pip_obj.cfg.cancel_proposal)
        log.info('Get cancel proposal information : {}'.format(proposalinfo_cancel))

        result = client_con_list_obj[0].staking.withdrew_staking(client_con_list_obj[0].node.staking_address)
        log.info('Node {} withdrew staking result {}'.format(client_con_list_obj[0].node.node_id, result))
        assert_code(result, 0)
        address = client_con_list_obj[0].node.staking_address

        result = version_proposal_vote(pip_obj)
        log.info('node vote version proposal result {}'.format(result))
        assert_code(result, 302020)

        result = cancel_proposal_vote(pip_obj)
        log.info('node vote cancel proposal result {}'.format(result))
        assert_code(result, 302020)

        address_test, _ = pip_obj.economic.account.generate_account(pip_obj.node.web3, 10**18 * 100000)
        result = client_noconsensus_obj.pip.vote(client_noconsensus_obj.node.node_id, proposalinfo_version.get('ProposalID'),
                                                 pip_obj.cfg.vote_option_yeas, address_test,
                                                 transaction_cfg=pip_obj.cfg.transaction_cfg)
        log.info('node {} vote param proposal result {}'.format(client_noconsensus_obj.node.node_id, result))
        assert_code(result, 302022)

        result = client_noconsensus_obj.pip.vote(client_noconsensus_obj.node.node_id, proposalinfo_cancel.get('ProposalID'),
                                                 pip_obj.cfg.vote_option_yeas, address_test,
                                                 transaction_cfg=pip_obj.cfg.transaction_cfg)
        log.info('node {} vote cancel proposal result {}'.format(client_noconsensus_obj.node.node_id, result))
        assert_code(result, 302022)

        pip_obj.economic.wait_settlement_blocknum(pip_obj.node, pip_obj.economic.unstaking_freeze_ratio)
        result = pip_obj.vote(pip_obj.node.node_id, proposalinfo_version.get('ProposalID'), pip_obj.cfg.vote_option_yeas,
                              address, transaction_cfg=pip_obj.cfg.transaction_cfg)
        log.info('Exited node vote version proposal result {}'.format(result))
        assert_code(result, 302022)

        result = pip_obj.vote(pip_obj.node.node_id, proposalinfo_cancel.get('ProposalID'), pip_obj.cfg.vote_option_nays,
                              address, transaction_cfg=pip_obj.cfg.transaction_cfg)
        log.info('Exited node vote cancel proposal result {}'.format(result))
        assert_code(result, 302022)

    @pytest.mark.P1
    def test_VO_TER_002_004(self, no_vp_proposal, client_candidate_obj, client_verifier_obj_list):
        pip_obj = client_candidate_obj.pip
        ver_pip_obj = client_verifier_obj_list[0].pip
        result = ver_pip_obj.submitParam(ver_pip_obj.node.node_id, str(time.time()), 'slashing', 'slashBlocksReward',
                                         '111', ver_pip_obj.node.staking_address,
                                         transaction_cfg=ver_pip_obj.cfg.transaction_cfg)
        log.info('Submit version proposal result : {}'.format(result))
        assert_code(result, 0)
        proposalinfo = ver_pip_obj.get_effect_proposal_info_of_vote(pip_obj.cfg.param_proposal)
        log.info('Get param proposal information : {}'.format(proposalinfo))
        result = ver_pip_obj.submitCancel(ver_pip_obj.node.node_id, str(time.time()), 2, proposalinfo.get('ProposalID'),
                                          ver_pip_obj.node.staking_address, transaction_cfg=pip_obj.cfg.transaction_cfg)
        log.info('Submit cancel proposal result : {}'.format(result))
        assert_code(result, 0)
        result = proposal_vote(pip_obj)
        log.info('Candidate node {} vote param proposal result : {}'.format(pip_obj.node.node_id, result))
        assert_code(result, 302022)

        result = cancel_proposal_vote(pip_obj)
        log.info('Candidate node {} vote cancel proposal result : {}'.format(pip_obj.node.node_id, result))
        assert_code(result, 302022)

    @pytest.mark.P1
    def test_VO_TER_001_003_005(self, proposal_voted_ca_pipobj_list, client_verifier_obj_list):
        pip_obj = proposal_voted_ca_pipobj_list[0]
        proposalinfo_version = client_verifier_obj_list[0].pip.get_effect_proposal_info_of_vote()
        log.info('Get proposal information :{}'.format(proposalinfo_version))
        pip_obj_test = client_verifier_obj_list[0].pip
        result = pip_obj_test.submitCancel(pip_obj_test.node.node_id, str(time.time()), 1, proposalinfo_version.get('ProposalID'),
                                           pip_obj_test.node.staking_address, transaction_cfg=pip_obj_test.cfg.transaction_cfg)
        log.info('Submit cancel proposal result : {}'.format(result))
        assert_code(result, 0)
        result = pip_obj_test.submitText(pip_obj_test.node.node_id, str(time.time()), pip_obj_test.node.staking_address,
                                         transaction_cfg=pip_obj_test.cfg.transaction_cfg)
        log.info('Submit text proposal result : {}'.format(result))
        assert_code(result, 0)
        result = version_proposal_vote(pip_obj)
        log.info('Candidate node {} vote version proposal result : {}'.format(pip_obj.node.node_id, result))
        assert_code(result, 302022)

        result = cancel_proposal_vote(pip_obj)
        log.info('Candidate node {} vote cancel proposal result : {}'.format(pip_obj.node.node_id, result))
        assert_code(result, 302022)

        result = text_proposal_vote(pip_obj)
        log.info('Candidate node {} vote text proposal result : {}'.format(pip_obj.node.node_id, result))
        assert_code(result, 302022)


class TestVoteCancelVersion():
    @pytest.mark.compatibility
    @pytest.mark.P0
    def test_VO_VO_002_V0_WA_002_V0_RE_002_V_STA_8(self, submit_cancel):
        pip_obj = submit_cancel
        address, _ = pip_obj.economic.account.generate_account(pip_obj.node.web3, 10 ** 18 * 10000)
        proposalinfo = pip_obj.get_effect_proposal_info_of_vote(pip_obj.cfg.cancel_proposal)
        log.info('Cancel proposal info : {}'.format(proposalinfo))
        result = pip_obj.vote(pip_obj.node.node_id, proposalinfo.get('ProposalID'), pip_obj.cfg.vote_option_Abstentions,
                              address, transaction_cfg=pip_obj.cfg.transaction_cfg)
        log.info('Not staking address , node {}, vote cancel proposal result {}'.format(pip_obj.node.node_id, result))
        assert_code(result, 302021)
        result = cancel_proposal_vote(pip_obj)
        assert_code(result, 0)
        result = cancel_proposal_vote(pip_obj)
        assert_code(result, 302027)

    @pytest.mark.P1
    def test_V_STA_6_7(self, submit_cancel, client_verifier_obj_list):
        pip_obj = submit_cancel
        proposalinfo = pip_obj.get_effect_proposal_info_of_vote(pip_obj.cfg.cancel_proposal)
        log.info('Cancel proposal info : {}'.format(proposalinfo))
        wait_block_number(pip_obj.node, proposalinfo.get('EndVotingBlock') - 10)
        result = cancel_proposal_vote(pip_obj)
        assert_code(result, 0)
        wait_block_number(pip_obj.node, proposalinfo.get('EndVotingBlock'))
        for client_obj in client_verifier_obj_list:
            if client_obj.node.node_id != pip_obj.node.node_id:
                pip_obj_test = client_obj.pip
                break
        result = pip_obj_test.vote(pip_obj_test.node.node_id, proposalinfo.get('ProposalID'), pip_obj_test.cfg.vote_option_Abstentions,
                                   pip_obj_test.node.staking_address, transaction_cfg=pip_obj_test.cfg.transaction_cfg)
        log.info('Node {} vote result {}'.format(pip_obj_test.node.node_id, result))
        assert_code(result, 302026)


class TestVoteCancelParam():
    @pytest.mark.compatibility
    @pytest.mark.P0
    def test_PP_VO_001_PP_VO_005_PP_VO_015_PP_VO_017(self, submit_cancel_param):
        pip_obj = submit_cancel_param
        address, _ = pip_obj.economic.account.generate_account(pip_obj.node.web3, 10**18 * 10000)
        proposalinfo = pip_obj.get_effect_proposal_info_of_vote(pip_obj.cfg.cancel_proposal)
        log.info('Cancel proposal info : {}'.format(proposalinfo))
        result = pip_obj.vote(pip_obj.node.node_id, proposalinfo.get('ProposalID'), pip_obj.cfg.vote_option_Abstentions,
                              address, transaction_cfg=pip_obj.cfg.transaction_cfg)
        log.info('Not staking address , node {}, vote cancel proposal result {}'.format(pip_obj.node.node_id, result))
        assert_code(result, 302021)
        result = cancel_proposal_vote(pip_obj)
        assert_code(result, 0)
        result = cancel_proposal_vote(pip_obj)
        assert_code(result, 302027)


class TestVoteParam():
    @pytest.mark.compatibility
    @pytest.mark.P0
    def test_PP_VO_002_PP_VO_008_PP_VO_018_PP_VO_016(self, submit_param):
        pip_obj = submit_param
        address, _ = pip_obj.economic.account.generate_account(pip_obj.node.web3, 10**18 * 10000)
        proposalinfo = pip_obj.get_effect_proposal_info_of_vote(pip_obj.cfg.param_proposal)
        log.info('param proposal info : {}'.format(proposalinfo))
        result = pip_obj.vote(pip_obj.node.node_id, proposalinfo.get('ProposalID'), pip_obj.cfg.vote_option_Abstentions,
                              address, transaction_cfg=pip_obj.cfg.transaction_cfg)
        log.info('Not staking address , node {}, vote param proposal result {}'.format(pip_obj.node.node_id, result))
        assert_code(result, 302021)
        result = proposal_vote(pip_obj)
        assert_code(result, 0)
        result = proposal_vote(pip_obj)
        assert_code(result, 302027)

    @pytest.mark.P2
    def test_PP_VO_009_PP_VO_010_V0_TE_001_V0_TE_002(self, submit_param, client_list_obj):
        pip_obj = submit_param
        result = pip_obj.submitText(pip_obj.node.node_id, str(time.time()), pip_obj.node.staking_address,
                                    transaction_cfg=pip_obj.cfg.transaction_cfg)
        log.info('Submit text proposal result : {}'.format(result))
        assert_code(result, 0)
        address = pip_obj.node.staking_address
        client_obj = get_client_obj(pip_obj.node.node_id, client_list_obj)
        result = client_obj.staking.withdrew_staking(pip_obj.node.staking_address)
        endblock = get_refund_to_account_block(pip_obj)
        log.info('Node {} withdrew staking result {}'.format(pip_obj.node.node_id, result))
        assert_code(result, 0)
        proposalinfo = pip_obj.get_effect_proposal_info_of_vote(pip_obj.cfg.param_proposal)
        log.info('Cancel proposal info : {}'.format(proposalinfo))
        proposalinfo_text = pip_obj.get_effect_proposal_info_of_vote(pip_obj.cfg.param_proposal)
        log.info('Text proposal info : {}'.format(proposalinfo_text))
        result = proposal_vote(pip_obj)
        assert_code(result, 302020)
        result = proposal_vote(pip_obj, proposaltype=pip_obj.cfg.text_proposal)
        assert_code(result, 302020)
        wait_block_number(pip_obj.node, endblock)
        result = pip_obj.vote(pip_obj.node.node_id, proposalinfo.get('ProposalID'), pip_obj.cfg.vote_option_nays, address,
                              transaction_cfg=pip_obj.cfg.transaction_cfg)
        assert_code(result, 302022)
        result = pip_obj.vote(pip_obj.node.node_id, proposalinfo_text.get('ProposalID'), pip_obj.cfg.vote_option_yeas,
                              address, transaction_cfg=pip_obj.cfg.transaction_cfg)
        assert_code(result, 302022)


@pytest.mark.compatibility
@pytest.mark.P0
def test_PP_VO_003_PP_VO_004_VS_EP_002_VS_EP_003(new_genesis_env, client_con_list_obj):
    genesis = from_dict(data_class=Genesis, data=new_genesis_env.genesis_config)
    genesis.economicModel.gov.paramProposalVoteDurationSeconds = 0
    new_genesis_env.set_genesis(genesis.to_dict())
    new_genesis_env.deploy_all()
    pip_obj = client_con_list_obj[0].pip
    result = pip_obj.submitParam(pip_obj.node.node_id, str(time.time()), 'slashing', 'slashBlocksReward', '123',
                                 pip_obj.node.staking_address, transaction_cfg=pip_obj.cfg.transaction_cfg)
    log.info('Submit param proposal result : {}'.format(result))
    assert_code(result, 0)
    proposalinfo = pip_obj.get_effect_proposal_info_of_vote(pip_obj.cfg.param_proposal)
    log.info('param proposal info {}'.format(proposalinfo))
    wait_block_number(pip_obj.node, proposalinfo.get('EndVotingBlock') - 10)
    result = proposal_vote(pip_obj)
    assert_code(result, 0)
    result = pip_obj.pip.getTallyResult(proposalinfo.get('ProposalID'))
    log.info('Interface getTallyResult result is {}'.format(result))
    assert_code(result, 302030)
    wait_block_number(pip_obj.node, proposalinfo.get('EndVotingBlock'))
    result = pip_obj.vote(pip_obj.node.node_id, proposalinfo.get('ProposalID'), pip_obj.cfg.vote_option_yeas,
                          pip_obj.node.staking_address, transaction_cfg=pip_obj.cfg.transaction_cfg)
    log.info('Node {} vote param proposal result : {}'.format(pip_obj.node.node_id, result))
    result = pip_obj.pip.getTallyResult(proposalinfo.get('ProposalID'))
    log.info('Interface getTallyResult result is {}'.format(result))
    assert_code(result, 0)


@pytest.mark.P0
def test_PP_VO_001_PP_VO_006_PP_VO_007_VS_EP_001(submit_cancel_param):
    pip_obj = submit_cancel_param
    proposalinfo = pip_obj.get_effect_proposal_info_of_vote(pip_obj.cfg.cancel_proposal)
    log.info('cancel proposal info {}'.format(proposalinfo))
    wait_block_number(pip_obj.node, proposalinfo.get('EndVotingBlock') - 8)
    result = cancel_proposal_vote(pip_obj)
    assert_code(result, 0)
    result = pip_obj.pip.getTallyResult(proposalinfo.get('ProposalID'))
    log.info('Interface getTallyResult result is {}'.format(result))
    assert_code(result, 302030)
    wait_block_number(pip_obj.node, proposalinfo.get('EndVotingBlock'))
    result = pip_obj.vote(pip_obj.node.node_id, proposalinfo.get('ProposalID'), pip_obj.cfg.vote_option_yeas,
                          pip_obj.node.staking_address, transaction_cfg=pip_obj.cfg.transaction_cfg)
    log.info('Node {} vote cancel proposal result : {}'.format(pip_obj.node.node_id, result))
    result = pip_obj.pip.getTallyResult(proposalinfo.get('ProposalID'))
    log.info('Interface getTallyResult result is {}'.format(result))
    assert_code(result, 0)


class TestVoteVPVerify():
    def vote_wrong_version(self, pip_obj, proposaltype):
        proposalinfo = pip_obj.get_effect_proposal_info_of_vote(proposaltype)
        log.info('Get proposal information : {}'.format(proposalinfo))
        program_version = pip_obj.cfg.version1
        if pip_obj.node.program_version == pip_obj.cfg.version1:
            program_version = pip_obj.cfg.version2
        result = pip_obj.vote(pip_obj.node.node_id, proposalinfo.get('ProposalID'), pip_obj.cfg.vote_option_yeas,
                              pip_obj.node.staking_address, program_version=program_version,
                              transaction_cfg=pip_obj.cfg.transaction_cfg)
        log.info('Wrong  program version vote result : {}'.format(result))
        return result

    def vote_wrong_versionsign(self, pip_obj, proposaltype):
        proposalinfo = pip_obj.get_effect_proposal_info_of_vote(proposaltype)
        log.info('Get proposal information : {}'.format(proposalinfo))
        version_sign = pip_obj.node.program_version_sign
        version_sign = version_sign.replace(version_sign[2:4], '11')
        result = pip_obj.vote(pip_obj.node.node_id, proposalinfo.get('ProposalID'), pip_obj.cfg.vote_option_yeas,
                              pip_obj.node.staking_address, version_sign=version_sign,
                              transaction_cfg=pip_obj.cfg.transaction_cfg)
        log.info('Wrong version sign vote result : {}'.format(result))
        return result

    @pytest.mark.P1
    def test_VO_VER_001_003_VO_SI_001_V_UP_1(self, submit_version):
        pip_obj = submit_version
        result = replace_platon_vote(pip_obj, bin=pip_obj.cfg.PLATON_NEW_BIN1)
        assert_code(result, 302025)
        result = replace_platon_vote(pip_obj, bin=pip_obj.cfg.PLATON_NEW_BIN2)
        assert_code(result, 302025)
        result = replace_platon_vote(pip_obj, bin=pip_obj.cfg.PLATON_NEW_BIN0)
        assert_code(result, 302025)
        result = replace_platon_vote(pip_obj, bin=pip_obj.cfg.PLATON_NEW_BIN3)
        assert_code(result, 302025)
        result = replace_platon_vote(pip_obj, bin=pip_obj.cfg.PLATON_NEW_BIN4)
        assert_code(result, 302025)
        version_sign = pip_obj.node.program_version_sign
        result = replace_platon_vote(pip_obj, bin=pip_obj.cfg.PLATON_NEW_BIN, version_sign=version_sign)
        assert_code(result, 302024)
        result = replace_platon_vote(pip_obj, bin=pip_obj.cfg.PLATON_NEW_BIN6)
        assert_code(result, 302025)
        result = replace_platon_vote(pip_obj, bin=pip_obj.cfg.PLATON_NEW_BIN7)
        assert_code(result, 302025)
        result = replace_platon_vote(pip_obj, bin=pip_obj.cfg.PLATON_NEW_BIN8)
        assert_code(result, 302025)

    @pytest.mark.P1
    def test_VO_VER_002_004_VO_SI_002(self, no_vp_proposal):
        pip_obj = no_vp_proposal
        result = pip_obj.submitVersion(pip_obj.node.node_id, str(time.time()), pip_obj.cfg.version9, 4,
                                       pip_obj.node.staking_address, transaction_cfg=pip_obj.cfg.transaction_cfg)
        log.info('Node {} submit version proposal result : {}'.format(pip_obj.node.node_id, result))
        assert_code(result, 0)
        result = replace_platon_vote(pip_obj, bin=pip_obj.cfg.PLATON_NEW_BIN0)
        assert_code(result, 302025)
        result = replace_platon_vote(pip_obj, bin=pip_obj.cfg.PLATON_NEW_BIN1)
        assert_code(result, 302025)
        result = replace_platon_vote(pip_obj, bin=pip_obj.cfg.PLATON_NEW_BIN2)
        assert_code(result, 302025)
        result = replace_platon_vote(pip_obj, bin=pip_obj.cfg.PLATON_NEW_BIN3)
        assert_code(result, 302025)
        result = replace_platon_vote(pip_obj, bin=pip_obj.cfg.PLATON_NEW_BIN4)
        assert_code(result, 302025)
        result = replace_platon_vote(pip_obj, bin=pip_obj.cfg.PLATON_NEW_BIN)
        assert_code(result, 302025)
        result = replace_platon_vote(pip_obj, bin=pip_obj.cfg.PLATON_NEW_BIN6)
        assert_code(result, 302025)
        result = replace_platon_vote(pip_obj, bin=pip_obj.cfg.PLATON_NEW_BIN7)
        assert_code(result, 302025)
        result = replace_platon_vote(pip_obj, bin=pip_obj.cfg.PLATON_NEW_BIN8)
        assert_code(result, 302025)
        version_sign = pip_obj.node.program_version_sign
        result = replace_platon_vote(pip_obj, bin=pip_obj.cfg.PLATON_NEW_BIN9, version_sign=version_sign)
        assert_code(result, 302024)

    @pytest.mark.P2
    def test_VO_SI_011_012(self, client_verifier_obj_list):
        pip_obj = client_verifier_obj_list[0].pip
        pip_obj_two = client_verifier_obj_list[1].pip
        result = pip_obj.submitText(pip_obj.node.node_id, str(time.time()), pip_obj.node.staking_address,
                                    transaction_cfg=pip_obj.cfg.transaction_cfg)
        log.info('Submit text proposal result : {}'.format(result))
        assert_code(result, 0)
        result = self.vote_wrong_version(pip_obj, pip_obj.cfg.text_proposal)
        assert_code(result, 0)
        result = self.vote_wrong_versionsign(pip_obj_two, pip_obj.cfg.text_proposal)
        assert_code(result, 0)

    @pytest.mark.P2
    def test_VO_SI_013_VO_SI_014_VO_SI_015_VO_SI_016(self, submit_cancel_param, client_verifier_obj_list):
        pip_obj = submit_cancel_param
        for client_obj in client_verifier_obj_list:
            if pip_obj.node.node_id != client_obj.node.node_id:
                pip_obj_two = client_obj.pip
                break
        result = self.vote_wrong_version(pip_obj, pip_obj.cfg.param_proposal)
        assert_code(result, 0)
        result = self.vote_wrong_versionsign(pip_obj_two, pip_obj.cfg.param_proposal)
        assert_code(result, 0)

        result = self.vote_wrong_version(pip_obj, pip_obj.cfg.cancel_proposal)
        assert_code(result, 0)
        result = self.vote_wrong_versionsign(pip_obj_two, pip_obj.cfg.cancel_proposal)
        assert_code(result, 0)

    @pytest.mark.P2
    def test_V0_POI_001(self, client_verifier_obj):
        pip_obj = client_verifier_obj.pip
        result = pip_obj.vote(pip_obj.node.node_id, '0x29b553fb979855751890aecf3e105948a11a21f121cad11f9e455c1f01b12345',
                              pip_obj.cfg.vote_option_yeas, pip_obj.node.staking_address,
                              transaction_cfg=pip_obj.cfg.transaction_cfg)
        log.info('Ineffective proposalid, vote result : {}'.format(result))
        assert_code(result, 302006)


class TestCadidateVote():
    @pytest.mark.P1
    def test_VO_TER_003_VO_TER_007_VO_TER_005_PP_VO_013(self, no_vp_proposal, client_candidate_obj, client_verifier_obj):
        ca_pip_obj = client_candidate_obj.pip
        ve_pip_obj = client_verifier_obj.pip
        submittpandvote([client_verifier_obj], 2)
        submitvpandvote([client_verifier_obj], votingrounds=1)
        proposalinfo_version = ve_pip_obj.get_effect_proposal_info_of_vote()
        log.info('Version proposal information {}'.format(proposalinfo_version))
        result = version_proposal_vote(ca_pip_obj)
        assert_code(result, 302022)
        result = proposal_vote(ca_pip_obj, proposaltype=ca_pip_obj.cfg.text_proposal)
        assert_code(result, 302022)
        wait_block_number(ca_pip_obj.node, proposalinfo_version.get('EndVotingBlock'))
        submitcppandvote([client_verifier_obj], [2])
        result = proposal_vote(ca_pip_obj, proposaltype=ca_pip_obj.cfg.param_proposal)
        assert_code(result, 302022)
        result = proposal_vote(ca_pip_obj, proposaltype=ca_pip_obj.cfg.cancel_proposal)
        assert_code(result, 302022)
