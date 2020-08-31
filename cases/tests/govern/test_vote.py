import pytest
import allure
from common.log import log
import time
from tests.lib.utils import assert_code, wait_block_number, upload_platon
from tests.lib.client import get_client_by_nodeid
from tests.govern.conftest import version_proposal_vote, get_refund_to_account_block, proposal_vote, verifier_node_version
from dacite import from_dict
from tests.lib.genesis import Genesis
from tests.govern.test_voting_statistics import submitvpandvote, submittpandvote, submitcppandvote


def replace_platon_vote(pip, bin=None, program_version=None, version_sign=None):
    '''
    Replace the bin package of the node, restart the node
    :param pip:
    :param bin:
    :return:
    '''
    if bin:
        upload_platon(pip.node, bin)
        pip.node.restart()
    if program_version is None:
        program_version = pip.node.program_version
    if version_sign is None:
        version_sign = pip.node.program_version_sign
    proposalinfo = pip.get_effect_proposal_info_of_vote()
    log.info('Get version proposal information {}'.format(proposalinfo))
    result = pip.vote(pip.node.node_id, proposalinfo.get('ProposalID'), pip.cfg.vote_option_yeas,
                      pip.node.staking_address, program_version=program_version, version_sign=version_sign,
                      transaction_cfg=pip.cfg.transaction_cfg)
    log.info('Node {} vote result : {}'.format(pip.node.node_id, result))
    return result


@pytest.fixture()
def voting_version_proposal_verifier_pip(client_verifier):
    pip = client_verifier.pip
    if pip.chain_version != pip.cfg.version0:
        log.info('The chain has been upgraded,restart!')
        client_verifier.economic.env.deploy_all()
    if pip.is_exist_effective_proposal:
        if pip.is_exist_effective_proposal_for_vote():
            proposalinfo = pip.get_effect_proposal_info_of_vote()
            log.info('get version proposalinfo : {}'.format(proposalinfo))
            if proposalinfo.get('EndVotingBlock') - pip.node.block_number > pip.economic.consensus_size * 2:
                return pip
        client_verifier.economic.env.deploy_all()
    result = pip.submitVersion(pip.node.node_id, str(time.time()), pip.cfg.version5, 10, pip.node.staking_address,
                               transaction_cfg=pip.cfg.transaction_cfg)
    log.info('node {} submit version proposal {}'.format(pip.node.node_id, result))
    assert_code(result, 0)
    return pip


@pytest.fixture()
def voting_text_proposal_verifier_pip(client_verifier):
    pip = client_verifier.pip
    if pip.is_exist_effective_proposal_for_vote(pip.cfg.text_proposal):
        proposalinfo = pip.get_effect_proposal_info_of_vote()
        if proposalinfo.get('EndVotingBlock') - pip.node.eth.blockNumber > 2 * pip.economic.consensus_size:
            return pip
        else:
            client_verifier.economic.env.deploy_all()
    result = pip.submitText(pip.node.node_id, str(time.time()), pip.node.staking_address,
                            transaction_cfg=pip.cfg.transaction_cfg)
    log.info('Submit text proposal result {}'.format(result))
    assert_code(result, 0)
    return pip


class TestVoteVP:
    @pytest.mark.P1
    @allure.title('Version proposal voting function verification--voting stage')
    def test_V_STA_2_to_5(self, no_vp_proposal, clients_verifier):
        pip = no_vp_proposal
        value = len(clients_verifier) - 2
        submitvpandvote(clients_verifier[:value], votingrounds=4)
        proposalinfo = pip.get_effect_proposal_info_of_vote()
        log.info('Get version proposal information : {}'.format(proposalinfo))
        result = version_proposal_vote(clients_verifier[-2].pip)
        log.info('Node {} vote proposal result : {}'.format(clients_verifier[-1].node.node_id, result))
        assert_code(result, 0)
        wait_block_number(pip.node, proposalinfo.get('EndVotingBlock'))
        pip = clients_verifier[-1].pip
        upload_platon(pip.node, pip.cfg.PLATON_NEW_BIN)
        pip.node.restart()
        log.info('Replace the platon bin and restart the node {}'.format(pip.node.node_id))
        result = pip.vote(pip.node.node_id, proposalinfo.get('ProposalID'), pip.cfg.vote_option_yeas,
                          pip.node.staking_address, transaction_cfg=pip.cfg.transaction_cfg)
        log.info('Node {} vote proposal result : {}'.format(clients_verifier[-1].node.node_id, result))
        assert_code(result, 302026)
        log.info('{}'.format(pip.pip.getTallyResult(proposalinfo.get('ProposalID'))))
        assert_code(pip.get_status_of_proposal(proposalinfo.get('ProposalID')), 4)
        result = pip.vote(pip.node.node_id, proposalinfo.get('ProposalID'), pip.cfg.vote_option_yeas,
                          pip.node.staking_address, transaction_cfg=pip.cfg.transaction_cfg)
        log.info('Node {} vote proposal result : {}'.format(clients_verifier[-1].node.node_id, result))
        assert_code(result, 302026)
        wait_block_number(pip.node, proposalinfo.get('ActiveBlock'))
        assert_code(pip.get_status_of_proposal(proposalinfo.get('ProposalID')), 5)


@pytest.mark.compatibility
@pytest.mark.P0
@allure.title('Version proposal voting function verification')
def test_VO_VO_001_V0_RE_001_V0_WA_001_V_STA_1_VO_OP_001_VO_OP_002(no_vp_proposal):
    pip = no_vp_proposal
    result = pip.submitVersion(pip.node.node_id, str(time.time()), pip.cfg.version8, 2,
                               pip.node.staking_address, transaction_cfg=pip.cfg.transaction_cfg)
    log.info('Submit version proposal result : {}'.format(result))
    assert_code(result, 0)
    proposalinfo = pip.get_effect_proposal_info_of_vote()
    log.info('Get version proposalinfo {}'.format(proposalinfo))

    upload_platon(pip.node, pip.cfg.PLATON_NEW_BIN)
    log.info('Replace the node platon package to {}'.format(pip.cfg.version5))
    pip.node.restart()
    log.info('Restart the node {}'.format(pip.node.node_id))
    address, _ = pip.economic.account.generate_account(pip.node.web3, 10**18 * 10000)
    result = pip.vote(pip.node.node_id, proposalinfo.get('ProposalID'), pip.cfg.vote_option_yeas,
                      address, transaction_cfg=pip.cfg.transaction_cfg)
    log.info('Not staking address vote result : {}'.format(result))
    assert_code(result, 302021)

    result = pip.vote(pip.node.node_id, proposalinfo.get('ProposalID'), pip.cfg.vote_option_nays,
                      pip.node.staking_address, transaction_cfg=pip.cfg.transaction_cfg)
    log.info('vote option {} result {}'.format(pip.cfg.vote_option_nays, result))
    assert_code(result, 302002)

    result = pip.vote(pip.node.node_id, proposalinfo.get('ProposalID'), pip.cfg.vote_option_Abstentions,
                      pip.node.staking_address, transaction_cfg=pip.cfg.transaction_cfg)
    log.info('vote option {} result {}'.format(pip.cfg.vote_option_Abstentions, result))
    assert_code(result, 302002)

    result = pip.vote(pip.node.node_id, proposalinfo.get('ProposalID'), 0,
                      pip.node.staking_address, transaction_cfg=pip.cfg.transaction_cfg)
    log.info('vote option {} result {}'.format(0, result))
    assert_code(result, 302002)

    result = pip.vote(pip.node.node_id, proposalinfo.get('ProposalID'), 'a',
                      pip.node.staking_address, transaction_cfg=pip.cfg.transaction_cfg)
    log.info('vote option {} result {}'.format(pip.cfg.vote_option_nays, result))
    assert_code(result, 302002)

    address, _ = pip.economic.account.generate_account(pip.node.web3, 10**18 * 10000)
    result = pip.vote(pip.node.node_id, proposalinfo.get('ProposalID'), pip.cfg.vote_option_yeas,
                      address, transaction_cfg=pip.cfg.transaction_cfg)
    assert_code(result, 302021)

    node_version = verifier_node_version(pip)
    result = version_proposal_vote(pip)
    assert_code(result, 0)
    verifier_node_version(pip, node_version)

    result = pip.vote(pip.node.node_id, proposalinfo.get('ProposalID'), pip.cfg.vote_option_yeas,
                      pip.node.staking_address, transaction_cfg=pip.cfg.transaction_cfg)
    log.info('vote duplicated result: {}'.format(result))
    assert_code(result, 302027)

    wait_block_number(pip.node, proposalinfo.get('EndVotingBlock'))
    result = pip.vote(pip.node.node_id, proposalinfo.get('ProposalID'), pip.cfg.vote_option_yeas,
                      pip.node.staking_address, transaction_cfg=pip.cfg.transaction_cfg)
    log.info('endblock vote result: {}'.format(result))
    assert_code(result, 302026)


@pytest.mark.P0
@allure.title('Text proposal voting function verification')
def test_VO_VO_003_V_STA_9_V_STA_10_V_STA_11_V0_WA_003_V0_RE_003(voting_text_proposal_verifier_pip, clients_verifier):
    pip = voting_text_proposal_verifier_pip
    proposalinfo = pip.get_effect_proposal_info_of_vote(pip.cfg.text_proposal)
    address, _ = pip.economic.account.generate_account(pip.node.web3, 10**18 * 10000)
    result = pip.vote(pip.node.node_id, proposalinfo.get('ProposalID'), pip.cfg.vote_option_yeas, address,
                      transaction_cfg=pip.cfg.transaction_cfg)
    log.info('Not staking address vote result {}'.format(result))
    assert_code(result, 302021)

    result = proposal_vote(pip, proposaltype=pip.cfg.text_proposal)
    log.info('vote result {}'.format(result))
    assert_code(result, 0)

    result = proposal_vote(pip, proposaltype=pip.cfg.text_proposal)
    log.info('Repeat vote  result {}'.format(result))
    assert_code(result, 302027)
    for client in clients_verifier:
        if client.node.node_id != pip.node.node_id:
            pip_test = client.pip
            break

    wait_block_number(pip.node, proposalinfo.get('EndVotingBlock') - 10)
    result = proposal_vote(pip_test, proposaltype=pip.cfg.text_proposal)
    log.info('Node {} vote result {}'.format(pip_test.node.node_id, result))
    assert_code(result, 0)

    wait_block_number(pip.node, proposalinfo.get('EndVotingBlock'))
    result = pip.vote(pip.node.node_id, proposalinfo.get('ProposalID'), pip.cfg.vote_option_nays,
                      pip.node.staking_address, transaction_cfg=pip.cfg.transaction_cfg)
    log.info('Endvoting block vote result {}'.format(result))
    assert_code(result, 302026)


class TestVoteNodeException:
    @pytest.mark.P0
    @allure.title('Voting function verification---Abnormal node')
    def test_VO_TE_001_002_PP_VO_009_010_PP_VO_011_012_PP_VO_014_VO_TER_008_VO_TER_006(self, new_genesis_env,
                                                                                       clients_consensus,
                                                                                       client_noconsensus):
        genesis = from_dict(data_class=Genesis, data=new_genesis_env.genesis_config)
        genesis.economicModel.gov.paramProposalVoteDurationSeconds = 640
        genesis.economicModel.staking.unStakeFreezeDuration = 2
        new_genesis_env.set_genesis(genesis.to_dict())
        new_genesis_env.deploy_all()
        pip = clients_consensus[0].pip
        result = pip.submitText(pip.node.node_id, str(time.time()), pip.node.staking_address,
                                transaction_cfg=pip.cfg.transaction_cfg)
        log.info('Submit text proposal information : {}'.format(result))
        assert_code(result, 0)
        proposalinfo_text = pip.get_effect_proposal_info_of_vote(pip.cfg.text_proposal)
        log.info('Get text proposal information : {}'.format(proposalinfo_text))

        result = pip.submitParam(pip.node.node_id, str(time.time()), 'slashing', 'slashBlocksReward',
                                 '123', pip.node.staking_address, transaction_cfg=pip.cfg.transaction_cfg)
        log.info('Submit param proposal information : {}'.format(result))
        assert_code(result, 0)
        proposalinfo_param = pip.get_effect_proposal_info_of_vote(pip.cfg.param_proposal)
        log.info('Get param proposal information : {}'.format(proposalinfo_param))

        result = pip.submitCancel(pip.node.node_id, str(time.time()), 10, proposalinfo_param.get('ProposalID'),
                                  pip.node.staking_address, transaction_cfg=pip.cfg.transaction_cfg)
        log.info('Submit cancel proposal information : {}'.format(result))
        assert_code(result, 0)
        proposalinfo_cancel = pip.get_effect_proposal_info_of_vote(pip.cfg.cancel_proposal)
        log.info('Get cancel proposal information : {}'.format(proposalinfo_cancel))

        result = clients_consensus[0].staking.withdrew_staking(clients_consensus[0].node.staking_address)
        log.info('Node {} withdrew staking result {}'.format(clients_consensus[0].node.node_id, result))
        assert_code(result, 0)
        address = clients_consensus[0].node.staking_address

        result = proposal_vote(pip, proposaltype=pip.cfg.text_proposal)
        log.info('node vote text proposal result {}'.format(result))
        assert_code(result, 302020)

        result = proposal_vote(pip)
        log.info('node vote param proposal result {}'.format(result))
        assert_code(result, 302020)

        result = proposal_vote(pip, proposaltype=pip.cfg.cancel_proposal)
        log.info('node vote cancel proposal result {}'.format(result))
        assert_code(result, 302020)

        address_test, _ = pip.economic.account.generate_account(pip.node.web3, 10**18 * 10000)
        result = client_noconsensus.pip.vote(client_noconsensus.node.node_id, proposalinfo_text.get('ProposalID'),
                                             pip.cfg.vote_option_yeas, address_test,
                                             transaction_cfg=pip.cfg.transaction_cfg)
        log.info('node {} vote text proposal result {}'.format(client_noconsensus.node.node_id, result))
        assert_code(result, 302022)

        result = client_noconsensus.pip.vote(client_noconsensus.node.node_id, proposalinfo_param.get('ProposalID'),
                                             pip.cfg.vote_option_yeas, address_test,
                                             transaction_cfg=pip.cfg.transaction_cfg)
        log.info('node {} vote param proposal result {}'.format(client_noconsensus.node.node_id, result))
        assert_code(result, 302022)

        result = client_noconsensus.pip.vote(client_noconsensus.node.node_id, proposalinfo_cancel.get('ProposalID'),
                                             pip.cfg.vote_option_yeas, address_test,
                                             transaction_cfg=pip.cfg.transaction_cfg)
        log.info('node {} vote cancel proposal result {}'.format(client_noconsensus.node.node_id, result))
        assert_code(result, 302022)

        pip.economic.wait_settlement_blocknum(pip.node, pip.economic.unstaking_freeze_ratio)
        result = pip.vote(pip.node.node_id, proposalinfo_text.get('ProposalID'), pip.cfg.vote_option_nays,
                          address, transaction_cfg=pip.cfg.transaction_cfg)
        log.info('Exited node vote text proposal result {}'.format(result))
        assert_code(result, 302022)

        result = pip.vote(pip.node.node_id, proposalinfo_param.get('ProposalID'), pip.cfg.vote_option_nays,
                          address, transaction_cfg=pip.cfg.transaction_cfg)
        log.info('Exited node vote param proposal result {}'.format(result))
        assert_code(result, 302022)

        result = pip.vote(pip.node.node_id, proposalinfo_cancel.get('ProposalID'), pip.cfg.vote_option_nays,
                          address, transaction_cfg=pip.cfg.transaction_cfg)
        log.info('Exited node vote cancel proposal result {}'.format(result))
        assert_code(result, 302022)

    @pytest.mark.P0
    @allure.title('Voting function verification---Abnormal node')
    def test_VO_VE_001_002_VO_CA_001_002_VO_TER_002_VO_TER_004(self, new_genesis_env, clients_consensus,
                                                               client_noconsensus):
        pip = clients_consensus[0].pip
        genesis = from_dict(data_class=Genesis, data=new_genesis_env.genesis_config)
        genesis.economicModel.gov.versionProposalVoteDurationSeconds = 3200
        genesis.economicModel.staking.unStakeFreezeDuration = 2
        new_genesis_env.set_genesis(genesis.to_dict())
        new_genesis_env.deploy_all()
        result = pip.submitVersion(pip.node.node_id, str(time.time()), pip.cfg.version5, 20,
                                   pip.node.staking_address, transaction_cfg=pip.cfg.transaction_cfg)
        log.info('Submit version proposal information : {}'.format(result))
        assert_code(result, 0)
        proposalinfo_version = pip.get_effect_proposal_info_of_vote(pip.cfg.version_proposal)
        log.info('Get param proposal information : {}'.format(proposalinfo_version))

        result = pip.submitCancel(pip.node.node_id, str(time.time()), 10, proposalinfo_version.get('ProposalID'),
                                  pip.node.staking_address, transaction_cfg=pip.cfg.transaction_cfg)
        log.info('Submit cancel proposal information : {}'.format(result))
        assert_code(result, 0)
        proposalinfo_cancel = pip.get_effect_proposal_info_of_vote(pip.cfg.cancel_proposal)
        log.info('Get cancel proposal information : {}'.format(proposalinfo_cancel))

        result = clients_consensus[0].staking.withdrew_staking(clients_consensus[0].node.staking_address)
        log.info('Node {} withdrew staking result {}'.format(clients_consensus[0].node.node_id, result))
        assert_code(result, 0)
        address = clients_consensus[0].node.staking_address

        result = version_proposal_vote(pip)
        log.info('node vote version proposal result {}'.format(result))
        assert_code(result, 302020)

        result = proposal_vote(pip, proposaltype=pip.cfg.cancel_proposal)
        log.info('node vote cancel proposal result {}'.format(result))
        assert_code(result, 302020)

        address_test, _ = pip.economic.account.generate_account(pip.node.web3, 10**18 * 100000)
        result = client_noconsensus.pip.vote(client_noconsensus.node.node_id, proposalinfo_version.get('ProposalID'),
                                             pip.cfg.vote_option_yeas, address_test,
                                             transaction_cfg=pip.cfg.transaction_cfg)
        log.info('node {} vote param proposal result {}'.format(client_noconsensus.node.node_id, result))
        assert_code(result, 302022)

        result = client_noconsensus.pip.vote(client_noconsensus.node.node_id, proposalinfo_cancel.get('ProposalID'),
                                             pip.cfg.vote_option_yeas, address_test,
                                             transaction_cfg=pip.cfg.transaction_cfg)
        log.info('node {} vote cancel proposal result {}'.format(client_noconsensus.node.node_id, result))
        assert_code(result, 302022)

        pip.economic.wait_settlement_blocknum(pip.node, pip.economic.unstaking_freeze_ratio)
        result = pip.vote(pip.node.node_id, proposalinfo_version.get('ProposalID'), pip.cfg.vote_option_yeas,
                          address, transaction_cfg=pip.cfg.transaction_cfg)
        log.info('Exited node vote version proposal result {}'.format(result))
        assert_code(result, 302022)

        result = pip.vote(pip.node.node_id, proposalinfo_cancel.get('ProposalID'), pip.cfg.vote_option_nays,
                          address, transaction_cfg=pip.cfg.transaction_cfg)
        log.info('Exited node vote cancel proposal result {}'.format(result))
        assert_code(result, 302022)

    @pytest.mark.P1
    @allure.title('Voting function verification---Abnormal node')
    def test_VO_TER_002_004(self, no_vp_proposal, client_candidate, clients_verifier):
        pip = client_candidate.pip
        ver_pip = clients_verifier[0].pip
        result = ver_pip.submitParam(ver_pip.node.node_id, str(time.time()), 'slashing', 'slashBlocksReward',
                                     '111', ver_pip.node.staking_address,
                                     transaction_cfg=ver_pip.cfg.transaction_cfg)
        log.info('Submit version proposal result : {}'.format(result))
        assert_code(result, 0)
        proposalinfo = ver_pip.get_effect_proposal_info_of_vote(pip.cfg.param_proposal)
        log.info('Get param proposal information : {}'.format(proposalinfo))
        result = ver_pip.submitCancel(ver_pip.node.node_id, str(time.time()), 2, proposalinfo.get('ProposalID'),
                                      ver_pip.node.staking_address, transaction_cfg=pip.cfg.transaction_cfg)
        log.info('Submit cancel proposal result : {}'.format(result))
        assert_code(result, 0)
        result = proposal_vote(pip)
        log.info('Candidate node {} vote param proposal result : {}'.format(pip.node.node_id, result))
        assert_code(result, 302022)

        result = proposal_vote(pip, proposaltype=pip.cfg.cancel_proposal)
        log.info('Candidate node {} vote cancel proposal result : {}'.format(pip.node.node_id, result))
        assert_code(result, 302022)

    @pytest.mark.P1
    @allure.title('Voting function verification')
    def test_VO_TER_001_003_005(self, candidate_has_proposal, client_verifier):
        pip = candidate_has_proposal
        pip_test = client_verifier.pip
        proposalinfo_version = pip_test.get_effect_proposal_info_of_vote()
        log.info('Get proposal information :{}'.format(proposalinfo_version))
        result = pip_test.submitCancel(pip_test.node.node_id, str(time.time()), 1, proposalinfo_version.get('ProposalID'),
                                       pip_test.node.staking_address, transaction_cfg=pip_test.cfg.transaction_cfg)
        log.info('Submit cancel proposal result : {}'.format(result))
        assert_code(result, 0)
        result = pip_test.submitText(pip_test.node.node_id, str(time.time()), pip_test.node.staking_address,
                                     transaction_cfg=pip_test.cfg.transaction_cfg)
        log.info('Submit text proposal result : {}'.format(result))
        assert_code(result, 0)
        result = version_proposal_vote(pip)
        log.info('Candidate node {} vote version proposal result : {}'.format(pip.node.node_id, result))
        assert_code(result, 302022)

        result = proposal_vote(pip, proposaltype=pip.cfg.cancel_proposal)
        log.info('Candidate node {} vote cancel proposal result : {}'.format(pip.node.node_id, result))
        assert_code(result, 302022)

        result = proposal_vote(pip, proposaltype=pip.cfg.text_proposal)
        log.info('Candidate node {} vote text proposal result : {}'.format(pip.node.node_id, result))
        assert_code(result, 302022)


class TestVoteCancelVersion:
    @pytest.mark.compatibility
    @pytest.mark.P0
    @allure.title('Cancel proposal voting function verification')
    def test_VO_VO_002_V0_WA_002_V0_RE_002_V_STA_8(self, submit_cancel):
        pip = submit_cancel
        address, _ = pip.economic.account.generate_account(pip.node.web3, 10 ** 18 * 10000)
        proposalinfo = pip.get_effect_proposal_info_of_vote(pip.cfg.cancel_proposal)
        log.info('Cancel proposal info : {}'.format(proposalinfo))
        result = pip.vote(pip.node.node_id, proposalinfo.get('ProposalID'), pip.cfg.vote_option_Abstentions,
                          address, transaction_cfg=pip.cfg.transaction_cfg)
        log.info('Not staking address , node {}, vote cancel proposal result {}'.format(pip.node.node_id, result))
        assert_code(result, 302021)
        result = proposal_vote(pip, proposaltype=pip.cfg.cancel_proposal)
        assert_code(result, 0)
        result = proposal_vote(pip, proposaltype=pip.cfg.cancel_proposal)
        assert_code(result, 302027)

    @pytest.mark.P1
    @allure.title('Cancel proposal voting function verification--candidate')
    def test_V_STA_6_7(self, submit_cancel, clients_verifier):
        pip = submit_cancel
        proposalinfo = pip.get_effect_proposal_info_of_vote(pip.cfg.cancel_proposal)
        log.info('Cancel proposal info : {}'.format(proposalinfo))
        wait_block_number(pip.node, proposalinfo.get('EndVotingBlock') - 10)
        result = proposal_vote(pip, proposaltype=pip.cfg.cancel_proposal)
        assert_code(result, 0)
        wait_block_number(pip.node, proposalinfo.get('EndVotingBlock'))
        for client in clients_verifier:
            if client.node.node_id != pip.node.node_id:
                pip_test = client.pip
                break
        result = pip_test.vote(pip_test.node.node_id, proposalinfo.get('ProposalID'), pip_test.cfg.vote_option_Abstentions,
                               pip_test.node.staking_address, transaction_cfg=pip_test.cfg.transaction_cfg)
        log.info('Node {} vote result {}'.format(pip_test.node.node_id, result))
        assert_code(result, 302026)


class TestVoteCancelParam:
    @pytest.mark.compatibility
    @pytest.mark.P0
    @allure.title('Cancel proposal voting function verification')
    def test_PP_VO_001_PP_VO_005_PP_VO_015_PP_VO_017(self, submit_cancel_param):
        pip = submit_cancel_param
        address, _ = pip.economic.account.generate_account(pip.node.web3, 10**18 * 10000)
        proposalinfo = pip.get_effect_proposal_info_of_vote(pip.cfg.cancel_proposal)
        log.info('Cancel proposal info : {}'.format(proposalinfo))
        result = pip.vote(pip.node.node_id, proposalinfo.get('ProposalID'), pip.cfg.vote_option_Abstentions,
                          address, transaction_cfg=pip.cfg.transaction_cfg)
        log.info('Not staking address , node {}, vote cancel proposal result {}'.format(pip.node.node_id, result))
        assert_code(result, 302021)
        result = proposal_vote(pip, proposaltype=pip.cfg.cancel_proposal)
        assert_code(result, 0)
        result = proposal_vote(pip, proposaltype=pip.cfg.cancel_proposal)
        assert_code(result, 302027)


class TestVoteParam:
    @pytest.mark.compatibility
    @pytest.mark.P0
    @allure.title('Param proposal voting function verification')
    def test_PP_VO_002_PP_VO_008_PP_VO_018_PP_VO_016(self, submit_param):
        pip = submit_param
        address, _ = pip.economic.account.generate_account(pip.node.web3, 10**18 * 10000)
        proposalinfo = pip.get_effect_proposal_info_of_vote(pip.cfg.param_proposal)
        log.info('param proposal info : {}'.format(proposalinfo))
        result = pip.vote(pip.node.node_id, proposalinfo.get('ProposalID'), pip.cfg.vote_option_Abstentions,
                          address, transaction_cfg=pip.cfg.transaction_cfg)
        log.info('Not staking address , node {}, vote param proposal result {}'.format(pip.node.node_id, result))
        assert_code(result, 302021)
        result = proposal_vote(pip)
        assert_code(result, 0)
        result = proposal_vote(pip)
        assert_code(result, 302027)

    @pytest.mark.P2
    @allure.title('voting function verification')
    def test_PP_VO_009_PP_VO_010_V0_TE_001_V0_TE_002(self, submit_param, all_clients):
        pip = submit_param
        result = pip.submitText(pip.node.node_id, str(time.time()), pip.node.staking_address,
                                transaction_cfg=pip.cfg.transaction_cfg)
        log.info('Submit text proposal result : {}'.format(result))
        assert_code(result, 0)
        address = pip.node.staking_address
        client = get_client_by_nodeid(pip.node.node_id, all_clients)
        result = client.staking.withdrew_staking(pip.node.staking_address)
        endblock = get_refund_to_account_block(pip)
        log.info('Node {} withdrew staking result {}'.format(pip.node.node_id, result))
        assert_code(result, 0)
        proposalinfo = pip.get_effect_proposal_info_of_vote(pip.cfg.param_proposal)
        log.info('Cancel proposal info : {}'.format(proposalinfo))
        proposalinfo_text = pip.get_effect_proposal_info_of_vote(pip.cfg.param_proposal)
        log.info('Text proposal info : {}'.format(proposalinfo_text))
        result = proposal_vote(pip)
        assert_code(result, 302020)
        result = proposal_vote(pip, proposaltype=pip.cfg.text_proposal)
        assert_code(result, 302020)
        wait_block_number(pip.node, endblock)
        result = pip.vote(pip.node.node_id, proposalinfo.get('ProposalID'), pip.cfg.vote_option_nays, address,
                          transaction_cfg=pip.cfg.transaction_cfg)
        assert_code(result, 302022)
        result = pip.vote(pip.node.node_id, proposalinfo_text.get('ProposalID'), pip.cfg.vote_option_yeas,
                          address, transaction_cfg=pip.cfg.transaction_cfg)
        assert_code(result, 302022)


@pytest.mark.compatibility
@pytest.mark.P0
@allure.title('Param proposal voting function verification')
def test_PP_VO_003_PP_VO_004_VS_EP_002_VS_EP_003(new_genesis_env, clients_consensus):
    genesis = from_dict(data_class=Genesis, data=new_genesis_env.genesis_config)
    genesis.economicModel.gov.paramProposalVoteDurationSeconds = 0
    new_genesis_env.set_genesis(genesis.to_dict())
    new_genesis_env.deploy_all()
    pip = clients_consensus[0].pip
    result = pip.submitParam(pip.node.node_id, str(time.time()), 'slashing', 'slashBlocksReward', '123',
                             pip.node.staking_address, transaction_cfg=pip.cfg.transaction_cfg)
    log.info('Submit param proposal result : {}'.format(result))
    assert_code(result, 0)
    proposalinfo = pip.get_effect_proposal_info_of_vote(pip.cfg.param_proposal)
    log.info('param proposal info {}'.format(proposalinfo))
    wait_block_number(pip.node, proposalinfo.get('EndVotingBlock') - 10)
    result = proposal_vote(pip)
    assert_code(result, 0)
    result = pip.pip.getTallyResult(proposalinfo.get('ProposalID'))
    log.info('Interface getTallyResult result is {}'.format(result))
    assert_code(result, 302030)
    wait_block_number(pip.node, proposalinfo.get('EndVotingBlock'))
    result = pip.vote(pip.node.node_id, proposalinfo.get('ProposalID'), pip.cfg.vote_option_yeas,
                      pip.node.staking_address, transaction_cfg=pip.cfg.transaction_cfg)
    log.info('Node {} vote param proposal result : {}'.format(pip.node.node_id, result))
    result = pip.pip.getTallyResult(proposalinfo.get('ProposalID'))
    log.info('Interface getTallyResult result is {}'.format(result))
    assert_code(result, 0)


@pytest.mark.P0
@allure.title('Cancel proposal voting function verification')
def test_PP_VO_001_PP_VO_006_PP_VO_007_VS_EP_001(submit_cancel_param):
    pip = submit_cancel_param
    proposalinfo = pip.get_effect_proposal_info_of_vote(pip.cfg.cancel_proposal)
    log.info('cancel proposal info {}'.format(proposalinfo))
    wait_block_number(pip.node, proposalinfo.get('EndVotingBlock') - 8)
    result = proposal_vote(pip, proposaltype=pip.cfg.cancel_proposal)
    assert_code(result, 0)
    result = pip.pip.getTallyResult(proposalinfo.get('ProposalID'))
    log.info('Interface getTallyResult result is {}'.format(result))
    assert_code(result, 302030)
    wait_block_number(pip.node, proposalinfo.get('EndVotingBlock'))
    result = pip.vote(pip.node.node_id, proposalinfo.get('ProposalID'), pip.cfg.vote_option_yeas,
                      pip.node.staking_address, transaction_cfg=pip.cfg.transaction_cfg)
    log.info('Node {} vote cancel proposal result : {}'.format(pip.node.node_id, result))
    result = pip.pip.getTallyResult(proposalinfo.get('ProposalID'))
    log.info('Interface getTallyResult result is {}'.format(result))
    assert_code(result, 0)


class TestVoteVPVerify:
    def vote_wrong_version(self, pip, proposaltype):
        proposalinfo = pip.get_effect_proposal_info_of_vote(proposaltype)
        log.info('Get proposal information : {}'.format(proposalinfo))
        program_version = pip.cfg.version1
        if pip.node.program_version == pip.cfg.version1:
            program_version = pip.cfg.version2
        result = pip.vote(pip.node.node_id, proposalinfo.get('ProposalID'), pip.cfg.vote_option_yeas,
                          pip.node.staking_address, program_version=program_version,
                          transaction_cfg=pip.cfg.transaction_cfg)
        log.info('Wrong  program version vote result : {}'.format(result))
        return result

    def vote_wrong_versionsign(self, pip, proposaltype):
        proposalinfo = pip.get_effect_proposal_info_of_vote(proposaltype)
        log.info('Get proposal information : {}'.format(proposalinfo))
        version_sign = pip.node.program_version_sign
        version_sign = version_sign.replace(version_sign[2:4], '11')
        result = pip.vote(pip.node.node_id, proposalinfo.get('ProposalID'), pip.cfg.vote_option_yeas,
                          pip.node.staking_address, version_sign=version_sign,
                          transaction_cfg=pip.cfg.transaction_cfg)
        log.info('Wrong version sign vote result : {}'.format(result))
        return result

    @pytest.mark.P1
    @allure.title('Version proposal voting function verification--platon version')
    def test_VO_VER_001_003_VO_SI_001_V_UP_1(self, submit_version):
        pip = submit_version
        result = replace_platon_vote(pip, bin=pip.cfg.PLATON_NEW_BIN1)
        assert_code(result, 302025)
        result = replace_platon_vote(pip, bin=pip.cfg.PLATON_NEW_BIN2)
        assert_code(result, 302025)
        result = replace_platon_vote(pip, bin=pip.cfg.PLATON_NEW_BIN0)
        assert_code(result, 302025)
        result = replace_platon_vote(pip, bin=pip.cfg.PLATON_NEW_BIN3)
        assert_code(result, 302025)
        result = replace_platon_vote(pip, bin=pip.cfg.PLATON_NEW_BIN4)
        assert_code(result, 302025)
        version_sign = pip.node.program_version_sign
        result = replace_platon_vote(pip, bin=pip.cfg.PLATON_NEW_BIN, version_sign=version_sign)
        assert_code(result, 302024)
        result = replace_platon_vote(pip, bin=pip.cfg.PLATON_NEW_BIN6)
        assert_code(result, 302025)
        result = replace_platon_vote(pip, bin=pip.cfg.PLATON_NEW_BIN7)
        assert_code(result, 302025)
        result = replace_platon_vote(pip, bin=pip.cfg.PLATON_NEW_BIN8)
        assert_code(result, 302025)

    @pytest.mark.P1
    @allure.title('Version proposal voting function verification--platon version')
    def test_VO_VER_002_004_VO_SI_002(self, no_vp_proposal):
        pip = no_vp_proposal
        result = pip.submitVersion(pip.node.node_id, str(time.time()), pip.cfg.version9, 4,
                                   pip.node.staking_address, transaction_cfg=pip.cfg.transaction_cfg)
        log.info('Node {} submit version proposal result : {}'.format(pip.node.node_id, result))
        assert_code(result, 0)
        result = replace_platon_vote(pip, bin=pip.cfg.PLATON_NEW_BIN0)
        assert_code(result, 302025)
        result = replace_platon_vote(pip, bin=pip.cfg.PLATON_NEW_BIN1)
        assert_code(result, 302025)
        result = replace_platon_vote(pip, bin=pip.cfg.PLATON_NEW_BIN2)
        assert_code(result, 302025)
        result = replace_platon_vote(pip, bin=pip.cfg.PLATON_NEW_BIN3)
        assert_code(result, 302025)
        result = replace_platon_vote(pip, bin=pip.cfg.PLATON_NEW_BIN4)
        assert_code(result, 302025)
        result = replace_platon_vote(pip, bin=pip.cfg.PLATON_NEW_BIN)
        assert_code(result, 302025)
        result = replace_platon_vote(pip, bin=pip.cfg.PLATON_NEW_BIN6)
        assert_code(result, 302025)
        result = replace_platon_vote(pip, bin=pip.cfg.PLATON_NEW_BIN7)
        assert_code(result, 302025)
        result = replace_platon_vote(pip, bin=pip.cfg.PLATON_NEW_BIN8)
        assert_code(result, 302025)
        version_sign = pip.node.program_version_sign
        result = replace_platon_vote(pip, bin=pip.cfg.PLATON_NEW_BIN9, version_sign=version_sign)
        assert_code(result, 302024)

    @pytest.mark.P2
    @allure.title('Text proposal voting function verification')
    def test_VO_SI_011_012(self, clients_verifier):
        pip = clients_verifier[0].pip
        pip_two = clients_verifier[1].pip
        result = pip.submitText(pip.node.node_id, str(time.time()), pip.node.staking_address,
                                transaction_cfg=pip.cfg.transaction_cfg)
        log.info('Submit text proposal result : {}'.format(result))
        assert_code(result, 0)
        result = self.vote_wrong_version(pip, pip.cfg.text_proposal)
        assert_code(result, 0)
        result = self.vote_wrong_versionsign(pip_two, pip.cfg.text_proposal)
        assert_code(result, 0)

    @pytest.mark.P2
    @allure.title('Cancel proposal voting function verification')
    def test_VO_SI_013_VO_SI_014_VO_SI_015_VO_SI_016(self, submit_cancel_param, clients_verifier):
        pip = submit_cancel_param
        for client in clients_verifier:
            if pip.node.node_id != client.node.node_id:
                pip_two = client.pip
                break
        result = self.vote_wrong_version(pip, pip.cfg.param_proposal)
        assert_code(result, 0)
        result = self.vote_wrong_versionsign(pip_two, pip.cfg.param_proposal)
        assert_code(result, 0)

        result = self.vote_wrong_version(pip, pip.cfg.cancel_proposal)
        assert_code(result, 0)
        result = self.vote_wrong_versionsign(pip_two, pip.cfg.cancel_proposal)
        assert_code(result, 0)

    @pytest.mark.P2
    @allure.title('Voting function verification--effective proposal id')
    def test_V0_POI_001(self, client_verifier):
        pip = client_verifier.pip
        result = pip.vote(pip.node.node_id, '0x29b553fb979855751890aecf3e105948a11a21f121cad11f9e455c1f01b12345',
                          pip.cfg.vote_option_yeas, pip.node.staking_address,
                          transaction_cfg=pip.cfg.transaction_cfg)
        log.info('Ineffective proposalid, vote result : {}'.format(result))
        assert_code(result, 302006)


class TestCadidateVote:
    @pytest.mark.P1
    @allure.title('Voting function verification--candidate')
    def test_VO_TER_003_VO_TER_007_VO_TER_005_PP_VO_013(self, no_vp_proposal, client_candidate, client_verifier):
        ca_pip = client_candidate.pip
        ve_pip = client_verifier.pip
        submittpandvote([client_verifier], 2)
        submitvpandvote([client_verifier], votingrounds=1)
        proposalinfo_version = ve_pip.get_effect_proposal_info_of_vote()
        log.info('Version proposal information {}'.format(proposalinfo_version))
        result = version_proposal_vote(ca_pip)
        assert_code(result, 302022)
        result = proposal_vote(ca_pip, proposaltype=ca_pip.cfg.text_proposal)
        assert_code(result, 302022)
        wait_block_number(ca_pip.node, proposalinfo_version.get('EndVotingBlock'))
        submitcppandvote([client_verifier], [2])
        result = proposal_vote(ca_pip, proposaltype=ca_pip.cfg.param_proposal)
        assert_code(result, 302022)
        result = proposal_vote(ca_pip, proposaltype=ca_pip.cfg.cancel_proposal)
        assert_code(result, 302022)
