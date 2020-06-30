from common.log import log
from dacite import from_dict
from tests.lib import Genesis
import pytest, allure
from tests.lib.utils import wait_block_number, assert_code, get_governable_parameter_value
from tests.lib.client import get_client_by_nodeid
import time
import math
from tests.govern.test_voting_statistics import submitcppandvote, submitcvpandvote, submitppandvote


@pytest.mark.P0
@pytest.mark.compatibility
@allure.title('Submit version proposal function verification')
def test_VP_SU_001(submit_version):
    pip = submit_version
    proposalinfo = pip.get_effect_proposal_info_of_vote()
    log.info('Get version proposal information : {}'.format(proposalinfo))
    endvotingblock_count = math.ceil(proposalinfo.get('SubmitBlock') / pip.economic.consensus_size +
                                     proposalinfo.get('EndVotingRounds')
                                     ) * pip.economic.consensus_size - 20
    log.info('Calculated endvoting block{},interface returned endvoting block{}'.format(endvotingblock_count,
                                                                                        proposalinfo.get('EndVotingBlock')))
    assert int(endvotingblock_count) == proposalinfo.get('EndVotingBlock')
    assert int(endvotingblock_count) + 21 == proposalinfo.get('ActiveBlock')


@pytest.mark.P0
@pytest.mark.compatibility
@allure.title('Submit cancel proposal function verification')
def test_CP_SU_001_CP_UN_001(submit_cancel):
    pip = submit_cancel
    proposalinfo = pip.get_effect_proposal_info_of_vote(pip.cfg.cancel_proposal)
    log.info('cancel proposalinfo : {}'.format(proposalinfo))
    endvotingblock_count = math.ceil(proposalinfo.get('SubmitBlock') / pip.economic.consensus_size +
                                     proposalinfo.get('EndVotingRounds')
                                     ) * pip.economic.consensus_size - 20
    log.info('Calculated endvoting block{},interface returned endvoting block{}'.format(endvotingblock_count,
                                                                                        proposalinfo.get('EndVotingBlock')))
    assert int(endvotingblock_count) == proposalinfo.get('EndVotingBlock')
    result = pip.submitCancel(pip.node.node_id, str(time.time()), 1, proposalinfo.get('ProposalID'),
                                  pip.node.staking_address, transaction_cfg=pip.cfg.transaction_cfg)
    log.info('Submit cancel proposal result : {}'.format(result))
    assert_code(result, 302014)


class TestsubmitCP:
    @pytest.mark.P0
    @allure.title('Submit param proposal function verification')
    def test_CP_SU_002_CP_SU_003(self, submit_param):
        pip = submit_param
        proposalinfo = pip.get_effect_proposal_info_of_vote(pip.cfg.param_proposal)
        log.info('param proposalinfo : {}'.format(proposalinfo))
        endvotingrounds_count = (proposalinfo.get('EndVotingBlock') -
                                 math.ceil(pip.node.block_number/pip.economic.consensus_size) *
                                 pip.economic.consensus_size) / pip.economic.consensus_size
        log.info('caculated endvoting rounds is {}'.format(endvotingrounds_count))
        result = pip.submitCancel(pip.node.node_id, str(time.time()), endvotingrounds_count + 1, proposalinfo.get('ProposalID'),
                                      pip.node.staking_address, transaction_cfg=pip.cfg.transaction_cfg)
        log.info('Submit cancel proposal result : {}'.format(result))
        assert_code(result, 302010)
        result = pip.submitCancel(pip.node.node_id, str(time.time()), endvotingrounds_count, proposalinfo.get('ProposalID'),
                                      pip.node.staking_address, transaction_cfg=pip.cfg.transaction_cfg)
        log.info('Submit cancel proposal result : {}'.format(result))
        assert_code(result, 0)

    @pytest.mark.P0
    @allure.title('Submit cancel param proposal function verification')
    def test_CP_SU_002_CP_UN_002(self, submit_cancel_param):
        pip = submit_cancel_param
        proposalinfo = pip.get_effect_proposal_info_of_vote(pip.cfg.cancel_proposal)
        log.info('cancel proposalinfo : {}'.format(proposalinfo))
        endvotingblock_count = math.ceil(proposalinfo.get('SubmitBlock') / pip.economic.consensus_size +
                                         proposalinfo.get('EndVotingRounds')
                                         ) * pip.economic.consensus_size - 20
        log.info('Calculated endvoting block{},interface returned endvoting block{}'.format(endvotingblock_count,
                                                                                            proposalinfo.get('EndVotingBlock')))
        assert int(endvotingblock_count) == proposalinfo.get('EndVotingBlock')
        result = pip.submitCancel(pip.node.node_id, str(time.time()), 1, proposalinfo.get('ProposalID'),
                                      pip.node.staking_address, transaction_cfg=pip.cfg.transaction_cfg)
        log.info('Submit cancel proposal result : {}'.format(result))
        assert_code(result, 302014)


@pytest.mark.P0
@pytest.mark.compatibility
@allure.title('Submit param proposal function verification')
def test_PP_SU_001_PP_UN_001_VP_UN_003(submit_param):
    pip = submit_param
    log.info('test chain version : {}'.format(pip.chain_version))
    proposalinfo = pip.get_effect_proposal_info_of_vote(pip.cfg.param_proposal)
    log.info('param proposalinfo : {}'.format(proposalinfo))
    endvotingblock_count = math.ceil(proposalinfo.get('SubmitBlock') / pip.economic.settlement_size +
                                     pip.economic.pp_vote_settlement_wheel
                                     ) * pip.economic.settlement_size
    log.info('Calculated endvoting block {},interface returned endvoting block {}'.format(endvotingblock_count,
                                                                                          proposalinfo.get('EndVotingBlock')))
    assert int(endvotingblock_count) == proposalinfo.get('EndVotingBlock')
    result = pip.submitParam(pip.node.node_id, str(time.time()), 'slashing', 'slashBlocksReward', '22',
                                 pip.node.staking_address, transaction_cfg=pip.cfg.transaction_cfg)
    log.info('There is a voting param proposal,submit param proposal result : {}'.format(result))
    assert_code(result, 302032)

    result = pip.submitVersion(pip.node.node_id, str(time.time()), pip.cfg.version5, 1,
                                   pip.node.staking_address, transaction_cfg=pip.cfg.transaction_cfg)
    log.info('There is a voting param proposal,submit version proposal result : {}'.format(result))
    assert_code(result, 302032)

@pytest.mark.P2
def test_PP_SU_021(new_genesis_env, client_consensus):
    genesis = from_dict(data_class=Genesis, data=new_genesis_env.genesis_config)
    genesis.config.cbft.period = 1000 * 2 * genesis.config.cbft.amount
    genesis.economicModel.common.maxEpochMinutes = 6
    genesis.economicModel.gov.textProposalVoteDurationSeconds = 161
    new_genesis_env.set_genesis(genesis.to_dict())
    new_genesis_env.deploy_all()
    pip = client_consensus.pip
    result = pip.submitParam(pip.node.node_id, str(time.time()), 'staking', 'maxValidators', '99', pip.node.staking_address,
                             transaction_cfg=pip.cfg.transaction_cfg)
    assert_code(result, 0)
    proposalinfo = pip.get_effect_proposal_info_of_vote(pip.cfg.param_proposal)
    log.info('Param proposal information : {}'.format(proposalinfo))
    endvotingblock_count = math.ceil(proposalinfo.get('SubmitBlock') / pip.economic.settlement_size +
                                     pip.economic.pp_vote_settlement_wheel
                                     ) * pip.economic.settlement_size
    log.info('Calculated endvoting block {},interface returned endvoting block {}'.format(endvotingblock_count,
                                                                                          proposalinfo.get(
                                                                                              'EndVotingBlock')))
    assert int(endvotingblock_count) == proposalinfo.get('EndVotingBlock')
    result = pip.submitText(pip.node.node_id, str(time.time()), pip.node.staking_address,
                            transaction_cfg=pip.cfg.transaction_cfg)
    log.info('Submit text proposal result : {}'.format(result))
    assert_code(result, 0)
    proposalinfo = pip.get_effect_proposal_info_of_vote(pip.cfg.text_proposal)
    log.info('Get text proposal information :{}'.format(proposalinfo))
    endvotingblock_count = math.ceil(proposalinfo.get('SubmitBlock') / pip.economic.consensus_size + 4
                                     ) * pip.economic.consensus_size - 20
    log.info('calcuted endvoting block {},interface returns {}'.format(endvotingblock_count,
                                                                       proposalinfo.get('EndVotingBlock')))
    assert int(endvotingblock_count) == proposalinfo.get('EndVotingBlock')



@pytest.mark.P0
@allure.title('Submit version proposal function verification')
def test_VP_VE_001_to_VP_VE_005(no_vp_proposal):
    pip_tmp = no_vp_proposal
    result = pip_tmp.submitVersion(pip_tmp.node.node_id, str(time.time()), pip_tmp.cfg.version1, 1,
                                       pip_tmp.node.staking_address, transaction_cfg=pip_tmp.cfg.transaction_cfg)
    log.info('Submit version proposal result : {}'.format(result))
    assert_code(result, 302011)

    result = pip_tmp.submitVersion(pip_tmp.node.node_id, str(time.time()), pip_tmp.cfg.version2, 1,
                                       pip_tmp.node.staking_address, transaction_cfg=pip_tmp.cfg.transaction_cfg)
    log.info('Submit version proposal result : {}'.format(result))
    assert_code(result, 302011)

    result = pip_tmp.submitVersion(pip_tmp.node.node_id, str(time.time()), pip_tmp.cfg.version3, 1,
                                       pip_tmp.node.staking_address, transaction_cfg=pip_tmp.cfg.transaction_cfg)
    log.info('Submit version proposal result : {}'.format(result))
    assert_code(result, 302011)

    result = pip_tmp.submitVersion(pip_tmp.node.node_id, str(time.time()), pip_tmp.chain_version, 1,
                                       pip_tmp.node.staking_address, transaction_cfg=pip_tmp.cfg.transaction_cfg)
    log.info('Submit version proposal result : {}'.format(result))
    assert_code(result, 302011)

    result = pip_tmp.submitVersion(pip_tmp.node.node_id, str(time.time()), pip_tmp.cfg.version8, 1,
                                       pip_tmp.node.staking_address, transaction_cfg=pip_tmp.cfg.transaction_cfg)
    log.info('Submit version proposal result : {}'.format(result))
    assert_code(result, 0)


@pytest.mark.P2
@allure.title('Nostaking address, submit version proposal function verification')
def test_VP_WA_001(no_vp_proposal):
    pip_tmp = no_vp_proposal
    address, _ = pip_tmp.economic.account.generate_account(pip_tmp.node.web3, 10**18 * 10000000)
    result = pip_tmp.submitVersion(pip_tmp.node.node_id, str(time.time()), pip_tmp.cfg.version5, 1,
                                       address, transaction_cfg=pip_tmp.cfg.transaction_cfg)
    log.info('Submit version proposal reuslt : {}'.format(result))
    assert_code(result, 302021)


@pytest.mark.P2
@allure.title('Nostaking address, submit text proposal function verification')
def test_TP_WA_001(client_verifier):
    pip = client_verifier.pip
    address, _ = pip.economic.account.generate_account(pip.node.web3, 10**18 * 10000000)
    result = pip.submitText(pip.node.node_id, str(time.time()), address,
                                transaction_cfg=pip.cfg.transaction_cfg)
    log.info('Submit text proposal reuslt : {}'.format(result))
    assert_code(result, 302021)


@pytest.mark.P0
@allure.title('submit text proposal function verification')
def test_TP_UN_001(submit_text):
    pip = submit_text
    result = pip.submitText(pip.node.node_id, str(time.time()), pip.node.staking_address,
                                transaction_cfg=pip.cfg.transaction_cfg)
    log.info('There is voting text proposal, submit text proposal result : {}'.format(result))
    assert_code(result, 0)


@pytest.mark.P0
@pytest.mark.compatibility
@allure.title('Submit version proposal function verification')
def test_VP_SU_001_VP_UN_001(submit_version):
    pip = submit_version
    result = pip.submitVersion(pip.node.node_id, str(time.time()), pip.cfg.version5, 1,
                                   pip.node.staking_address, transaction_cfg=pip.cfg.transaction_cfg)
    log.info('There is voting version proposal, submit version proposal result : {}'.format(result))
    assert_code(result, 302012)


@pytest.mark.P0
@allure.title('There is preactive proposal, submit proposal function verification')
def test_VP_UN_002_CP_ID_002(preactive_proposal_pips, new_genesis_env):
    pip = preactive_proposal_pips[0]
    proposalinfo = pip.get_effect_proposal_info_of_preactive()
    log.info('Get preactive proposal info: {}'.format(proposalinfo))

    result = pip.submitVersion(pip.node.node_id, str(time.time()), pip.cfg.version5, 1,
                                   pip.node.staking_address,
                                   transaction_cfg=pip.cfg.transaction_cfg)
    log.info('There is preactive version proposal, submit version proposal result : {}'.format(result))
    assert_code(result, 302013)

    result = pip.submitCancel(pip.node.node_id, str(time.time()), 1, proposalinfo.get('ProposalID'),
                                  pip.node.staking_address, transaction_cfg=pip.cfg.transaction_cfg)
    log.info('there is preactive version proposal, submit cancel proposal result: {}'.format(result))
    assert_code(result, 302017)


@pytest.mark.P0
@allure.title('There is preactive proposal, submit cancel parammeter proposal function verification')
def test_PP_UN_003(preactive_proposal_pips, new_genesis_env):
    pip = preactive_proposal_pips[0]
    proposalinfo = pip.get_effect_proposal_info_of_preactive()
    log.info('Get preactive proposal info: {}'.format(proposalinfo))

    result = pip.submitParam(pip.node.node_id, str(time.time()), 'slashing', 'slashBlocksReward',
                                 '84', pip.node.staking_address, transaction_cfg=pip.cfg.transaction_cfg)
    log.info('there is preactive version proposal, submit cancel param proposal result: {}'.format(result))
    assert_code(result, 302013)


class TestEndVotingRounds:
    def update_setting(self, new_genesis_env, pip, value=0):
        genesis = from_dict(data_class=Genesis, data=new_genesis_env.genesis_config)
        genesis.economicModel.gov.versionProposalVoteDurationSeconds = 2 * pip.economic.consensus_size + value
        genesis.economicModel.gov.textProposalVoteDurationSeconds = 2 * pip.economic.consensus_size + value
        new_genesis_env.set_genesis(genesis.to_dict())
        new_genesis_env.deploy_all()

    @pytest.mark.P1
    @allure.title('Submit version and text proposal function verification--endvoting rounds')
    def test_VP_CR_001_VP_CR_002_VP_CR_007_TP_TE_002(self, new_genesis_env, client_verifier):
        '''
        Proposal vote duration set consensus size's accompanying number +1
        :param pip_env:
        :param pip:
        :return:
        '''
        pip = client_verifier.pip
        self.update_setting(new_genesis_env, pip, value=1)
        result = pip.submitVersion(pip.node.node_id, str(time.time()), pip.cfg.version5, 3,
                                       pip.node.staking_address,
                                       transaction_cfg=pip.cfg.transaction_cfg)
        log.info('endvoting rounds is three, subtmit version proposal result : {}'.format(result))
        assert_code(result, 302010)

        result = pip.submitVersion(pip.node.node_id, str(time.time()), pip.cfg.version5, 0,
                                       pip.node.staking_address,
                                       transaction_cfg=pip.cfg.transaction_cfg)
        log.info('endvoting rounds is zero, subtmit version proposal result : {}'.format(result))
        assert_code(result, 302009)

        result = pip.submitVersion(pip.node.node_id, str(time.time()), pip.cfg.version5, 2,
                                       pip.node.staking_address,
                                       transaction_cfg=pip.cfg.transaction_cfg)
        log.info('endvoting rounds is two, subtmit version proposal result : {}'.format(result))
        assert_code(result, 0)

        result = pip.submitText(pip.node.node_id, str(time.time()), pip.node.staking_address,
                                    transaction_cfg=pip.cfg.transaction_cfg)
        log.info('Submit text proposal result : {}'.format(result))
        assert_code(result, 0)
        proposalinfo = pip.get_effect_proposal_info_of_vote(pip.cfg.text_proposal)
        log.info('Get text proposal information :{}'.format(proposalinfo))
        endvotingblock_count = math.ceil(proposalinfo.get('SubmitBlock') / pip.economic.consensus_size + 2
                                         ) * pip.economic.consensus_size - 20
        log.info('calcuted endvoting block {},interface returns {}'.format(endvotingblock_count,
                                                                           proposalinfo.get('EndVotingBlock')))
        assert int(endvotingblock_count) == proposalinfo.get('EndVotingBlock')

    @pytest.mark.P1
    @allure.title('Submit version and text proposal function verification--endvoting rounds')
    def test_VP_CR_003_VP_CR_004_VP_CR_007_TP_TE_003(self, new_genesis_env, client_verifier):
        pip = client_verifier.pip
        self.update_setting(new_genesis_env, pip, value=-1)
        result = pip.submitVersion(pip.node.node_id, str(time.time()), pip.cfg.version5, 2, pip.node.staking_address,
                                       transaction_cfg=pip.cfg.transaction_cfg)
        log.info('endvoting rounds is two, subtmit version proposal result : {}'.format(result))
        assert_code(result, 302010)

        result = pip.submitVersion(pip.node.node_id, str(time.time()), pip.cfg.version5, 0, pip.node.staking_address,
                                       transaction_cfg=pip.cfg.transaction_cfg)
        log.info('Endvoting rounds is zero, subtmit version proposal result : {}'.format(result))
        assert_code(result, 302009)

        result = pip.submitVersion(pip.node.node_id, str(time.time()), pip.cfg.version5, 1, pip.node.staking_address,
                                       transaction_cfg=pip.cfg.transaction_cfg)
        log.info('Endvoting rounds is one, subtmit version proposal result : {}'.format(result))
        assert_code(result, 0)

        result = pip.submitText(pip.node.node_id, str(time.time()), pip.node.staking_address,
                                    transaction_cfg=pip.cfg.transaction_cfg)
        log.info('Submit text proposal result : {}'.format(result))
        assert_code(result, 0)
        proposalinfo = pip.get_effect_proposal_info_of_vote(pip.cfg.text_proposal)
        log.info('Get text proposal information : {}'.format(proposalinfo))
        endvotingblock_count = math.ceil(proposalinfo.get('SubmitBlock') / pip.economic.consensus_size + 1
                                         ) * pip.economic.consensus_size - 20
        log.info('Calcuted endvoting block {},interface return {}'.format(endvotingblock_count,
                                                                          proposalinfo.get('EndVotingBlock')))
        assert int(endvotingblock_count) == proposalinfo.get('EndVotingBlock')

    @pytest.mark.P1
    @pytest.mark.compatibility
    @allure.title('Submit version and text proposal function verification--endvoting rounds')
    def test_VP_CR_005_VP_CR_006_TP_TE_001(self, new_genesis_env, client_verifier):
        pip = client_verifier.pip
        self.update_setting(new_genesis_env, pip)
        result = pip.submitVersion(pip.node.node_id, str(time.time()), pip.cfg.version5, 3,
                                       pip.node.staking_address,
                                       transaction_cfg=pip.cfg.transaction_cfg)
        log.info('Endvoting rounds is three, subtmit version proposal result : {}'.format(result))
        assert_code(result, 302010)

        result = pip.submitVersion(pip.node.node_id, str(time.time()), pip.cfg.version5, 0,
                                       pip.node.staking_address,
                                       transaction_cfg=pip.cfg.transaction_cfg)
        log.info('Endvoting rounds is zero, subtmit version proposal result : {}'.format(result))
        assert_code(result, 302009)

        result = pip.submitVersion(pip.node.node_id, str(time.time()), pip.cfg.version5, 2,
                                       pip.node.staking_address,
                                       transaction_cfg=pip.cfg.transaction_cfg)
        log.info('Endvoting rounds is two, subtmit version proposal result : {}'.format(result))
        assert_code(result, 0)

        result = pip.submitText(pip.node.node_id, str(time.time()), pip.node.staking_address,
                                    transaction_cfg=pip.cfg.transaction_cfg)
        log.info('Submit text proposal result : {}'.format(result))
        assert_code(result, 0)
        proposalinfo = pip.get_effect_proposal_info_of_vote(pip.cfg.text_proposal)
        log.info('Get text proposal information {}'.format(proposalinfo))
        endvotingblock_count = math.ceil(proposalinfo.get('SubmitBlock') / pip.economic.consensus_size + 2
                                         ) * pip.economic.consensus_size - 20
        log.info('calcuted endvoting block {},interface return {}'.format(endvotingblock_count,
                                                                          proposalinfo.get('EndVotingBlock')))
        assert int(endvotingblock_count) == proposalinfo.get('EndVotingBlock')

        proosalinfo = pip.get_effect_proposal_info_of_vote(pip.cfg.text_proposal)
        log.info('text proposalinfo: {}'.format(proosalinfo))
        result = pip.submitCancel(pip.node.node_id, str(time.time()), 1, proosalinfo.get('ProposalID'),
                                      pip.node.staking_address, transaction_cfg=pip.cfg.transaction_cfg)
        log.info('submit cancel result: {}'.format(result))
        assert_code(result, 302016)


class TestNoVerifierSubmitProposal:
    @pytest.mark.P0
    @allure.title('New node submit version and text proposal function verification')
    def test_VP_PR_002_TP_PR_002(self, no_vp_proposal, client_new_node):
        pip = client_new_node.pip
        address, _ = pip.economic.account.generate_account(pip.node.web3, 10**18 * 10000000)
        result = pip.submitVersion(pip.node.node_id, str(time.time()), pip.cfg.version5, 1, address,
                                       transaction_cfg=pip.cfg.transaction_cfg)
        log.info('New node submit version proposal : {}'.format(result))
        assert_code(result, 302022)

        result = pip.submitText(pip.node.node_id, str(time.time()), address,
                                    transaction_cfg=pip.cfg.transaction_cfg)
        log.info('New node submit text proposal : {}'.format(result))
        assert_code(result, 302022)

    @pytest.mark.P0
    @allure.title('Candidate node submit version and text proposal function verification')
    def test_VP_PR_001_TP_PR_001(self, client_candidate):
        pip = client_candidate.pip
        result = pip.submitVersion(pip.node.node_id, str(time.time()),
                                       pip.cfg.version5, 1, pip.node.staking_address,
                                       transaction_cfg=pip.cfg.transaction_cfg)
        log.info('Candidate node {} submit version proposal result : {}'.format(pip.node.node_id, result))
        assert_code(result, 302022)

        result = pip.submitText(pip.node.node_id, str(time.time()), pip.node.staking_address,
                                    transaction_cfg=pip.cfg.transaction_cfg)
        log.info('candidate node {} submit text proposal result : {}'.format(pip.node.node_id, result))
        assert_code(result, 302022)

    @pytest.mark.P2
    @allure.title('Abnormal node submit proposal')
    def test_VP_PR_003_VP_PR_004_TP_PR_003_TP_PR_004(self, client_verifier):
        address = client_verifier.node.staking_address
        result = client_verifier.staking.withdrew_staking(address)
        log.info('Node {} withdrew staking result : {}'.format(client_verifier.node.node_id, result))
        assert_code(result, 0)
        result = client_verifier.pip.submitVersion(client_verifier.node.node_id, str(time.time()),
                                                   client_verifier.pip.cfg.version5, 1, address,
                                                   transaction_cfg=client_verifier.pip.cfg.transaction_cfg)
        log.info('Node exiting submit version proposal :{}'.format(result))
        assert_code(result, 302020)

        result = client_verifier.pip.submitText(client_verifier.node.node_id, str(time.time()), address,
                                                transaction_cfg=client_verifier.pip.cfg.transaction_cfg)
        log.info('Node exiting submit text proposal : {}'.format(result))
        assert_code(result, 302020)

        client_verifier.economic.wait_settlement_blocknum(client_verifier.node,
                                                          number=client_verifier.economic.unstaking_freeze_ratio)
        result = client_verifier.pip.submitVersion(client_verifier.node.node_id, str(time.time()),
                                                   client_verifier.pip.cfg.version5, 1, address,
                                                   transaction_cfg=client_verifier.pip.cfg.transaction_cfg)
        log.info('Node exited submit version proposal : {}'.format(result))
        assert_code(result, 302022)

        client_verifier.economic.wait_settlement_blocknum(client_verifier.node,
                                                          number=client_verifier.economic.unstaking_freeze_ratio)
        result = client_verifier.pip.submitText(client_verifier.node.node_id, str(time.time()), address,
                                                transaction_cfg=client_verifier.pip.cfg.transaction_cfg)
        log.info('Node exited submit text proposal : {}'.format(result))
        assert_code(result, 302022)


class TestSubmitCancel:
    @pytest.mark.P0
    @allure.title('Nostaking address submit cancel proposal')
    def test_CP_WA_001(self, submit_version):
        pip = submit_version
        address, _ = pip.economic.account.generate_account(pip.node.web3, 10**18 * 10000)
        proposalinfo = pip.get_effect_proposal_info_of_vote()
        log.info('Get version proposal information : {}'.format(proposalinfo))
        result = pip.submitCancel(pip.node.node_id, str(time.time()), 1, proposalinfo.get('ProposalID'),
                                      address, transaction_cfg=pip.cfg.transaction_cfg)
        log.info('Submit cancel proposal result : {}'.format(result))
        assert_code(result, 302021)

    @pytest.mark.P0
    @allure.title('New node submit cancel proposal')
    def test_CP_PR_001(self, new_node_has_proposal):
        pip = new_node_has_proposal
        address, _ = pip.economic.account.generate_account(pip.node.web3, 10**18 * 10000)
        proposalinfo = pip.get_effect_proposal_info_of_vote()
        log.info('Get version proposal information : {}'.format(proposalinfo))
        result = pip.submitCancel(pip.node.node_id, str(time.time()), 1,
                                      proposalinfo.get('ProposalID'),
                                      address, transaction_cfg=pip.cfg.transaction_cfg)
        log.info('Submit cancel proposal resullt : {}'.format(result))
        assert_code(result, 302022)

    @pytest.mark.P0
    @allure.title('Candidate node submit cancel proposal')
    def test_CP_PR_002(self, candidate_has_proposal):
        pip = candidate_has_proposal
        proposalinfo = pip.get_effect_proposal_info_of_vote()
        if not proposalinfo:
            time.sleep(10)
            proposalinfo = pip.get_effect_proposal_info_of_vote()
        log.info('Get proposal information {}'.format(proposalinfo))
        result = pip.submitCancel(pip.node.node_id, str(time.time()), 1,
                                      proposalinfo.get('ProposalID'),
                                      pip.node.staking_address, transaction_cfg=pip.cfg.transaction_cfg)
        log.info('Candidate submit cancel proposal result : {}'.format(result))
        assert_code(result, 302022)

    @pytest.mark.P2
    @allure.title('Abnormal node submit cancel proposal')
    def test_CP_PR_003_CP_PR_004(self, new_genesis_env, client_consensus):
        genesis = from_dict(data_class=Genesis, data=new_genesis_env.genesis_config)
        genesis.economicModel.gov.versionProposalVoteDurationSeconds = 10000
        new_genesis_env.set_genesis(genesis.to_dict())
        new_genesis_env.deploy_all()
        pip = client_consensus.pip
        result = pip.submitVersion(pip.node.node_id, str(time.time()), pip.cfg.version5, 20,
                                       pip.node.staking_address, transaction_cfg=pip.cfg.transaction_cfg)
        log.info('Submit version proposal result : {}'.format(result))
        assert_code(result, 0)
        proposalinfo = pip.get_effect_proposal_info_of_vote()
        log.info('proposalinfo: {}'.format(proposalinfo))
        client = client_consensus
        address = pip.node.staking_address
        result = client.staking.withdrew_staking(address)
        log.info('nodeid: {} withdrewstaking result: {}'.format(client.node.node_id, result))
        assert_code(result, 0)
        result = client.pip.submitCancel(client.node.node_id, str(time.time()), 1,
                                             proposalinfo.get('ProposalID'), address,
                                             transaction_cfg=client.pip.cfg.transaction_cfg)
        log.info('node exiting，cancel proposal result: {}'.format(result))
        assert_code(result, 302020)

        client.economic.wait_settlement_blocknum(client.node, number=client.economic.unstaking_freeze_ratio)
        result = client.pip.submitCancel(client.node.node_id, str(time.time()), 1,
                                             proposalinfo.get('ProposalID'), address,
                                             transaction_cfg=client.pip.cfg.transaction_cfg)
        log.info('exited node，cancel proposal result: {}'.format(result))
        assert_code(result, 302022)

    @pytest.mark.P0
    @allure.title('Submit version  proposal function verification--endvoting rounds')
    def test_CP_CR_001_CP_CR_002(self, submit_version):
        pip = submit_version
        proposalinfo = pip.get_effect_proposal_info_of_vote()
        log.info('proposalinfo: {}'.format(proposalinfo))
        endvoting_rounds = (proposalinfo.get('EndVotingBlock') + 20 - math.ceil(
            pip.node.block_number / pip.economic.consensus_size) * pip.economic.consensus_size
            ) / pip.economic.consensus_size
        result = pip.submitCancel(pip.node.node_id, str(time.time()), endvoting_rounds,
                                      proposalinfo.get('ProposalID'), pip.node.staking_address,
                                      transaction_cfg=pip.cfg.transaction_cfg)
        log.info('endvoting_rounds:{}， cancel proposal result:{}'.format(endvoting_rounds, result))
        assert_code(result, 302010)

        result = pip.submitCancel(pip.node.node_id, str(time.time()), endvoting_rounds + 1,
                                      proposalinfo.get('ProposalID'), pip.node.staking_address,
                                      transaction_cfg=pip.cfg.transaction_cfg)
        log.info('endvoting_rounds:{}， cancel proposal result:{}'.format(endvoting_rounds + 1, result))
        assert_code(result, 302010)

    @pytest.mark.P0
    @allure.title('Submit cancel  proposal function verification--ineffective proposal id')
    def test_CP_ID_001(self, no_vp_proposal):
        pip = no_vp_proposal
        result = pip.submitCancel(pip.node.node_id, str(time.time()), 1,
                                      '0x49b83cfc4b99462f7131d14d80c73b6657237753cd1e878e8d62dc2e9f574123',
                                      pip.node.staking_address, transaction_cfg=pip.cfg.transaction_cfg)
        log.info('cancel proposal result: {}'.format(result))
        assert_code(result, 302015)

    @pytest.mark.P0
    @allure.title('Submit parammeter  proposal function verification')
    def test_CP_ID_004_CP_ID_003(self, new_genesis_env, client_consensus):
        genesis = from_dict(data_class=Genesis, data=new_genesis_env.genesis_config)
        genesis.economicModel.gov.paramProposalVoteDurationSeconds = 0
        new_genesis_env.set_genesis(genesis.to_dict())
        new_genesis_env.deploy_all()
        pip = client_consensus.pip
        result = pip.submitParam(pip.node.node_id, str(time.time()), 'slashing', 'slashBlocksReward', '123',
                                     pip.node.staking_address, transaction_cfg=pip.cfg.transaction_cfg)
        log.info('Submit parameter proposal result : {}'.format(result))
        assert_code(result, 0)
        result = pip.submitText(pip.node.node_id, str(time.time()), pip.node.staking_address,
                                    transaction_cfg=pip.cfg.transaction_cfg)
        log.info('Submit text proposal result : {}'.format(result))
        assert_code(result, 0)
        proposalinfo_param = pip.get_effect_proposal_info_of_vote(pip.cfg.param_proposal)
        log.info('Get parameter proposal information : {}'.format(proposalinfo_param))
        proposalinfo_text = pip.get_effect_proposal_info_of_vote(pip.cfg.text_proposal)
        log.info('Get text proposal information : {}'.format(proposalinfo_text))
        result = pip.submitCancel(pip.node.node_id, str(time.time()), 1, proposalinfo_text.get('ProposalID'),
                                      pip.node.staking_address, transaction_cfg=pip.cfg.transaction_cfg)
        log.info('Submit cancel proposal result : {}'.format(result))
        wait_block_number(pip.node, proposalinfo_param.get('EndVotingBlock'))
        result = pip.submitCancel(pip.node.node_id, str(time.time()), 1, proposalinfo_param.get('ProposalID'),
                                      pip.node.staking_address, transaction_cfg=pip.cfg.transaction_cfg)
        log.info('Submit cancel proposal result : {}'.format(result))
        assert_code(result, 302017)


class TestPP:
    @pytest.mark.P0
    @allure.title('Submit parammeter  proposal function verification')
    def test_PP_SU_001_PP_SU_002(self, no_vp_proposal, all_clients):
        pip = no_vp_proposal
        client = get_client_by_nodeid(pip.node.node_id, all_clients)
        result = pip.submitParam(pip.node.node_id, str(time.time()), 'slashing', 'slashBlocksReward', '',
                                     pip.node.staking_address, transaction_cfg=pip.cfg.transaction_cfg)
        log.info('Submit param proposal result : {}'.format(result))
        assert_code(result, 3)

        result = pip.submitParam(pip.node.node_id, str(time.time()), 'slashing', 'slashBlocksReward', '1.1',
                                     pip.node.staking_address, transaction_cfg=pip.cfg.transaction_cfg)
        log.info('Submit param proposal result : {}'.format(result))
        assert_code(result, 3)

        result = pip.submitParam(pip.node.node_id, str(time.time()), 'slashing', 'slashBlocksReward', '-1',
                                     pip.node.staking_address, transaction_cfg=pip.cfg.transaction_cfg)
        log.info('Submit param proposal result : {}'.format(result))
        assert_code(result, 3)

        result = pip.submitParam(pip.node.node_id, str(time.time()), 'slashing', 'slashBlocksReward', '50000',
                                     pip.node.staking_address, transaction_cfg=pip.cfg.transaction_cfg)
        log.info('Submit param proposal result : {}'.format(result))
        assert_code(result, 3)

        result = pip.submitParam(pip.node.node_id, str(time.time()), 'slashing', 'slashBlocksReward', 1,
                                     pip.node.staking_address, transaction_cfg=pip.cfg.transaction_cfg)
        log.info('Submit param proposal result : {}'.format(result))
        assert_code(result, 3)

        result = pip.submitParam(pip.node.node_id, str(time.time()), 'slashing', 'slashBlocksReward', '1000000',
                                     pip.node.staking_address, transaction_cfg=pip.cfg.transaction_cfg)
        log.info('Submit param proposal result : {}'.format(result))
        assert_code(result, 3)

        result = pip.submitParam(pip.node.node_id, str(time.time()), 'slashing', 'slashBlocksReward',
                                     str(get_governable_parameter_value(client, 'slashBlocksReward')),
                                     pip.node.staking_address, transaction_cfg=pip.cfg.transaction_cfg)
        log.info('Submit param proposal result : {}'.format(result))
        assert_code(result, 302034)

        if pip.economic.slash_blocks_reward != 0:
            result = pip.submitParam(pip.node.node_id, str(time.time()), 'slashing', 'slashBlocksReward', '0',
                                         pip.node.staking_address, transaction_cfg=pip.cfg.transaction_cfg)
            log.info('Submit param proposal result : {}'.format(result))
            assert_code(result, 0)

    @pytest.mark.P0
    @allure.title('Submit parammeter  proposal function verification')
    def test_PP_SU_002(self, no_vp_proposal, all_clients):
        pip = no_vp_proposal
        if str(get_governable_parameter_value(all_clients[0], 'slashBlocksReward')) != 49999:
            result = pip.submitParam(pip.node.node_id, str(time.time()), 'slashing', 'slashBlocksReward', '49999',
                                         pip.node.staking_address, transaction_cfg=pip.cfg.transaction_cfg)
            log.info('Submit param proposal result : {}'.format(result))
            assert_code(result, 0)

    @pytest.mark.P0
    @allure.title('Submit parammeter  proposal function verification')
    def test_PP_SU_003_PP_SU_004(self, no_vp_proposal, all_clients):
        pip = no_vp_proposal
        client = get_client_by_nodeid(pip.node.node_id, all_clients)
        result = pip.submitParam(pip.node.node_id, str(time.time()), 'slashing', 'maxEvidenceAge', '',
                                     pip.node.staking_address, transaction_cfg=pip.cfg.transaction_cfg)
        log.info('Submit param proposal result : {}'.format(result))
        assert_code(result, 3)

        result = pip.submitParam(pip.node.node_id, str(time.time()), 'slashing', 'maxEvidenceAge', '1.1',
                                     pip.node.staking_address, transaction_cfg=pip.cfg.transaction_cfg)
        log.info('Submit param proposal result : {}'.format(result))
        assert_code(result, 3)

        result = pip.submitParam(pip.node.node_id, str(time.time()), 'slashing', 'maxEvidenceAge', '-1',
                                     pip.node.staking_address, transaction_cfg=pip.cfg.transaction_cfg)
        log.info('Submit param proposal result : {}'.format(result))
        assert_code(result, 3)

        result = pip.submitParam(pip.node.node_id, str(time.time()), 'slashing', 'maxEvidenceAge', '0',
                                     pip.node.staking_address, transaction_cfg=pip.cfg.transaction_cfg)
        log.info('Submit param proposal result : {}'.format(result))
        assert_code(result, 3)

        result = pip.submitParam(pip.node.node_id, str(time.time()), 'slashing', 'maxEvidenceAge', 1,
                                     pip.node.staking_address, transaction_cfg=pip.cfg.transaction_cfg)
        log.info('Submit param proposal result : {}'.format(result))
        assert_code(result, 3)

        result = pip.submitParam(pip.node.node_id, str(time.time()), 'slashing', 'maxEvidenceAge',
                                     str(get_governable_parameter_value(client, 'maxEvidenceAge')),
                                     pip.node.staking_address, transaction_cfg=pip.cfg.transaction_cfg)
        log.info('Submit param proposal result : {}'.format(result))
        assert_code(result, 302034)

        result = pip.submitParam(pip.node.node_id, str(time.time()), 'slashing', 'maxEvidenceAge',
                                     str(get_governable_parameter_value(client, 'unStakeFreezeDuration')),
                                     pip.node.staking_address, transaction_cfg=pip.cfg.transaction_cfg)
        log.info('Submit param proposal result : {}'.format(result))
        assert_code(result, 3)

        if int(get_governable_parameter_value(client, 'maxEvidenceAge')) != 1:
            result = pip.submitParam(pip.node.node_id, str(time.time()), 'slashing', 'maxEvidenceAge', '1',
                                         pip.node.staking_address, transaction_cfg=pip.cfg.transaction_cfg)
            log.info('Submit param proposal result : {}'.format(result))
            assert_code(result, 0)

    @pytest.mark.P0
    @allure.title('Submit parammeter  proposal function verification')
    def test_PP_SU_004(self, no_vp_proposal, all_clients):
        pip = no_vp_proposal
        if str(get_governable_parameter_value(all_clients[0], 'maxEvidenceAge')) != str(
            pip.economic.unstaking_freeze_ratio - 1) and str(pip.economic.unstaking_freeze_ratio - 1) <= str(
                get_governable_parameter_value(all_clients[0], 'unStakeFreezeDuration')):
            result = pip.submitParam(pip.node.node_id, str(time.time()), 'slashing', 'maxEvidenceAge',
                                         str(pip.economic.unstaking_freeze_ratio - 1),
                                         pip.node.staking_address, transaction_cfg=pip.cfg.transaction_cfg)
            log.info('Submit param proposal result : {}'.format(result))
            assert_code(result, 0)

    @pytest.mark.P0
    @allure.title('Submit parammeter  proposal function verification')
    def test_PP_SU_005_PP_SU_006(self, no_vp_proposal, all_clients):
        pip = no_vp_proposal
        client = get_client_by_nodeid(pip.node.node_id, all_clients)
        result = pip.submitParam(pip.node.node_id, str(time.time()), 'slashing', 'slashFractionDuplicateSign', '',
                                     pip.node.staking_address, transaction_cfg=pip.cfg.transaction_cfg)
        log.info('Submit param proposal result : {}'.format(result))
        assert_code(result, 3)

        result = pip.submitParam(pip.node.node_id, str(time.time()), 'slashing', 'slashFractionDuplicateSign', '1.1',
                                     pip.node.staking_address, transaction_cfg=pip.cfg.transaction_cfg)
        log.info('Submit param proposal result : {}'.format(result))
        assert_code(result, 3)

        result = pip.submitParam(pip.node.node_id, str(time.time()), 'slashing', 'slashFractionDuplicateSign', '-1',
                                     pip.node.staking_address, transaction_cfg=pip.cfg.transaction_cfg)
        log.info('Submit param proposal result : {}'.format(result))
        assert_code(result, 3)

        result = pip.submitParam(pip.node.node_id, str(time.time()), 'slashing', 'slashFractionDuplicateSign', '0',
                                     pip.node.staking_address, transaction_cfg=pip.cfg.transaction_cfg)
        log.info('Submit param proposal result : {}'.format(result))
        assert_code(result, 3)

        result = pip.submitParam(pip.node.node_id, str(time.time()), 'slashing', 'slashFractionDuplicateSign', 1,
                                     pip.node.staking_address, transaction_cfg=pip.cfg.transaction_cfg)
        log.info('Submit param proposal result : {}'.format(result))
        assert_code(result, 3)

        result = pip.submitParam(pip.node.node_id, str(time.time()), 'slashing', 'slashFractionDuplicateSign',
                                     '10001', pip.node.staking_address, transaction_cfg=pip.cfg.transaction_cfg)
        log.info('Submit param proposal result : {}'.format(result))
        assert_code(result, 3)

        result = pip.submitParam(pip.node.node_id, str(time.time()), 'slashing', 'slashFractionDuplicateSign',
                                     str(get_governable_parameter_value(client, 'slashFractionDuplicateSign')),
                                     pip.node.staking_address, transaction_cfg=pip.cfg.transaction_cfg)
        log.info('Submit param proposal result : {}'.format(result))
        assert_code(result, 302034)

        if int(get_governable_parameter_value(client, 'maxEvidenceAge')) != 1:
            result = pip.submitParam(pip.node.node_id, str(time.time()), 'slashing', 'slashFractionDuplicateSign', '1',
                                         pip.node.staking_address, transaction_cfg=pip.cfg.transaction_cfg)
            log.info('Submit param proposal result : {}'.format(result))
            assert_code(result, 0)

    @pytest.mark.P0
    @allure.title('Submit parammeter  proposal function verification')
    def test_PP_SU_006(self, no_vp_proposal):
        pip = no_vp_proposal
        result = pip.submitParam(pip.node.node_id, str(time.time()), 'slashing', 'slashFractionDuplicateSign', '10000',
                                     pip.node.staking_address, transaction_cfg=pip.cfg.transaction_cfg)
        log.info('Submit param proposal result : {}'.format(result))
        assert_code(result, 0)

    @pytest.mark.P0
    @allure.title('Submit parammeter  proposal function verification')
    def test_PP_SU_007_PP_SU_008(self, no_vp_proposal, all_clients):
        pip = no_vp_proposal
        client = get_client_by_nodeid(pip.node.node_id, all_clients)
        result = pip.submitParam(pip.node.node_id, str(time.time()), 'slashing', 'duplicateSignReportReward', '',
                                     pip.node.staking_address, transaction_cfg=pip.cfg.transaction_cfg)
        log.info('Submit param proposal result : {}'.format(result))
        assert_code(result, 3)

        result = pip.submitParam(pip.node.node_id, str(time.time()), 'slashing', 'duplicateSignReportReward', '1.1',
                                     pip.node.staking_address, transaction_cfg=pip.cfg.transaction_cfg)
        log.info('Submit param proposal result : {}'.format(result))
        assert_code(result, 3)

        result = pip.submitParam(pip.node.node_id, str(time.time()), 'slashing', 'duplicateSignReportReward', '-1',
                                     pip.node.staking_address, transaction_cfg=pip.cfg.transaction_cfg)
        log.info('Submit param proposal result : {}'.format(result))
        assert_code(result, 3)

        result = pip.submitParam(pip.node.node_id, str(time.time()), 'slashing', 'duplicateSignReportReward', '0',
                                     pip.node.staking_address, transaction_cfg=pip.cfg.transaction_cfg)
        log.info('Submit param proposal result : {}'.format(result))
        assert_code(result, 3)

        result = pip.submitParam(pip.node.node_id, str(time.time()), 'slashing', 'duplicateSignReportReward', 1,
                                     pip.node.staking_address, transaction_cfg=pip.cfg.transaction_cfg)
        log.info('Submit param proposal result : {}'.format(result))
        assert_code(result, 3)

        result = pip.submitParam(pip.node.node_id, str(time.time()), 'slashing', 'duplicateSignReportReward',
                                     '81', pip.node.staking_address, transaction_cfg=pip.cfg.transaction_cfg)
        log.info('Submit param proposal result : {}'.format(result))
        assert_code(result, 3)

        result = pip.submitParam(pip.node.node_id, str(time.time()), 'slashing', 'duplicateSignReportReward',
                                     str(get_governable_parameter_value(client, 'duplicateSignReportReward')),
                                     pip.node.staking_address, transaction_cfg=pip.cfg.transaction_cfg)
        log.info('Submit param proposal result : {}'.format(result))
        assert_code(result, 302034)

        if int(get_governable_parameter_value(client, 'duplicateSignReportReward')) != 1:
            result = pip.submitParam(pip.node.node_id, str(time.time()), 'slashing', 'duplicateSignReportReward', '1',
                                         pip.node.staking_address, transaction_cfg=pip.cfg.transaction_cfg)
            log.info('Submit param proposal result : {}'.format(result))
            assert_code(result, 0)

    @pytest.mark.P0
    @allure.title('Submit parammeter  proposal function verification')
    def test_PP_SU_008(self, no_vp_proposal, all_clients):
        pip = no_vp_proposal
        if int(get_governable_parameter_value(all_clients[0], 'duplicateSignReportReward')) != 80:
            result = pip.submitParam(pip.node.node_id, str(time.time()), 'slashing', 'duplicateSignReportReward',
                                         '80', pip.node.staking_address,
                                         transaction_cfg=pip.cfg.transaction_cfg)
            log.info('Submit param proposal result : {}'.format(result))
            assert_code(result, 0)

    @pytest.mark.P0
    @allure.title('Submit parammeter  proposal function verification')
    def test_PP_SU_009_PP_SU_010(self, no_vp_proposal, all_clients):
        pip = no_vp_proposal
        client = get_client_by_nodeid(pip.node.node_id, all_clients)
        result = pip.submitParam(pip.node.node_id, str(time.time()), 'staking', 'stakeThreshold', '',
                                     pip.node.staking_address, transaction_cfg=pip.cfg.transaction_cfg)
        log.info('Submit param proposal result : {}'.format(result))
        assert_code(result, 3)

        value = 10**18 * 1000000 + 0.1
        result = pip.submitParam(pip.node.node_id, str(time.time()), 'staking', 'stakeThreshold', str(value),
                                     pip.node.staking_address, transaction_cfg=pip.cfg.transaction_cfg)
        log.info('Submit param proposal result : {}'.format(result))
        assert_code(result, 3)

        value = -10**18 * 1000000
        result = pip.submitParam(pip.node.node_id, str(time.time()), 'staking', 'stakeThreshold', str(value),
                                     pip.node.staking_address, transaction_cfg=pip.cfg.transaction_cfg)
        log.info('Submit param proposal result : {}'.format(result))
        assert_code(result, 3)

        result = pip.submitParam(pip.node.node_id, str(time.time()), 'staking', 'stakeThreshold', '0',
                                     pip.node.staking_address, transaction_cfg=pip.cfg.transaction_cfg)
        log.info('Submit param proposal result : {}'.format(result))
        assert_code(result, 3)

        result = pip.submitParam(pip.node.node_id, str(time.time()), 'staking', 'stakeThreshold', 10**18 * 1000000,
                                     pip.node.staking_address, transaction_cfg=pip.cfg.transaction_cfg)
        log.info('Submit param proposal result : {}'.format(result))
        assert_code(result, 3)

        value = 10**18 * 1000000 - 1
        result = pip.submitParam(pip.node.node_id, str(time.time()), 'staking', 'stakeThreshold',
                                     str(value), pip.node.staking_address, transaction_cfg=pip.cfg.transaction_cfg)
        log.info('Submit param proposal result : {}'.format(result))
        assert_code(result, 3)

        value = 10**18 * 10000000 + 1
        result = pip.submitParam(pip.node.node_id, str(time.time()), 'staking', 'stakeThreshold',
                                     str(value), pip.node.staking_address, transaction_cfg=pip.cfg.transaction_cfg)
        log.info('Submit param proposal result : {}'.format(result))
        assert_code(result, 3)

        result = pip.submitParam(pip.node.node_id, str(time.time()), 'staking', 'stakeThreshold',
                                     str(get_governable_parameter_value(client, 'stakeThreshold')),
                                     pip.node.staking_address, transaction_cfg=pip.cfg.transaction_cfg)
        log.info('Submit param proposal result : {}'.format(result))
        assert_code(result, 302034)

        if int(get_governable_parameter_value(client, 'stakeThreshold')) != 10**18 * 1000000:
            value = 10**18 * 1000000
            result = pip.submitParam(pip.node.node_id, str(time.time()), 'staking', 'stakeThreshold', str(value),
                                         pip.node.staking_address, transaction_cfg=pip.cfg.transaction_cfg)
            log.info('Submit param proposal result : {}'.format(result))
            assert_code(result, 0)

    @pytest.mark.P0
    @allure.title('Submit parammeter  proposal function verification')
    def test_PP_SU_010(self, no_vp_proposal):
        pip = no_vp_proposal
        value = 10**18 * 10000000
        result = pip.submitParam(pip.node.node_id, str(time.time()), 'staking', 'stakeThreshold',
                                     str(value),
                                     pip.node.staking_address, transaction_cfg=pip.cfg.transaction_cfg)
        log.info('Submit param proposal result : {}'.format(result))
        assert_code(result, 0)

    @pytest.mark.P0
    @allure.title('Submit parammeter  proposal function verification')
    def test_PP_SU_011_PP_SU_012(self, no_vp_proposal, all_clients):
        pip = no_vp_proposal
        client = get_client_by_nodeid(pip.node.node_id, all_clients)
        result = pip.submitParam(pip.node.node_id, str(time.time()), 'staking', 'operatingThreshold', '',
                                     pip.node.staking_address, transaction_cfg=pip.cfg.transaction_cfg)
        log.info('Submit param proposal result : {}'.format(result))
        assert_code(result, 3)

        value = 10**18 * 10 + 0.5
        result = pip.submitParam(pip.node.node_id, str(time.time()), 'staking', 'operatingThreshold', str(value),
                                     pip.node.staking_address, transaction_cfg=pip.cfg.transaction_cfg)
        log.info('Submit param proposal result : {}'.format(result))
        assert_code(result, 3)

        value = -10**18 * 10
        result = pip.submitParam(pip.node.node_id, str(time.time()), 'staking', 'operatingThreshold', str(value),
                                     pip.node.staking_address, transaction_cfg=pip.cfg.transaction_cfg)
        log.info('Submit param proposal result : {}'.format(result))
        assert_code(result, 3)

        result = pip.submitParam(pip.node.node_id, str(time.time()), 'staking', 'operatingThreshold', '0',
                                     pip.node.staking_address, transaction_cfg=pip.cfg.transaction_cfg)
        log.info('Submit param proposal result : {}'.format(result))
        assert_code(result, 3)

        result = pip.submitParam(pip.node.node_id, str(time.time()), 'staking', 'operatingThreshold', 1,
                                     pip.node.staking_address, transaction_cfg=pip.cfg.transaction_cfg)
        log.info('Submit param proposal result : {}'.format(result))
        assert_code(result, 3)

        value = 10**18 * 10 - 1
        result = pip.submitParam(pip.node.node_id, str(time.time()), 'staking', 'operatingThreshold',
                                     str(value), pip.node.staking_address, transaction_cfg=pip.cfg.transaction_cfg)
        log.info('Submit param proposal result : {}'.format(result))
        assert_code(result, 3)

        value = 10**18 * 10000 + 1
        result = pip.submitParam(pip.node.node_id, str(time.time()), 'staking', 'operatingThreshold',
                                     str(value), pip.node.staking_address, transaction_cfg=pip.cfg.transaction_cfg)
        log.info('Submit param proposal result : {}'.format(result))
        assert_code(result, 3)

        result = pip.submitParam(pip.node.node_id, str(time.time()), 'staking', 'operatingThreshold',
                                     str(get_governable_parameter_value(client, 'operatingThreshold')),
                                     pip.node.staking_address, transaction_cfg=pip.cfg.transaction_cfg)
        log.info('Submit param proposal result : {}'.format(result))
        assert_code(result, 302034)

        if int(get_governable_parameter_value(client, 'operatingThreshold')) != 10**18 * 10:
            value = 10**18 * 10
            result = pip.submitParam(pip.node.node_id, str(time.time()), 'staking', 'operatingThreshold', str(value),
                                         pip.node.staking_address, transaction_cfg=pip.cfg.transaction_cfg)
            log.info('Submit param proposal result : {}'.format(result))
            assert_code(result, 0)

    @pytest.mark.P0
    @allure.title('Submit parammeter  proposal function verification')
    def test_PP_SU_012(self, no_vp_proposal):
        pip = no_vp_proposal
        value = 10**18 * 10000
        result = pip.submitParam(pip.node.node_id, str(time.time()), 'staking', 'operatingThreshold',
                                     str(value), pip.node.staking_address,
                                     transaction_cfg=pip.cfg.transaction_cfg)
        log.info('Submit param proposal result : {}'.format(result))
        assert_code(result, 0)

    @pytest.mark.P0
    @allure.title('Submit parammeter  proposal function verification')
    def test_PP_SU_013_PP_SU_014(self, no_vp_proposal, all_clients):
        pip = no_vp_proposal
        client = get_client_by_nodeid(pip.node.node_id, all_clients)
        result = pip.submitParam(pip.node.node_id, str(time.time()), 'staking', 'unStakeFreezeDuration', '',
                                     pip.node.staking_address, transaction_cfg=pip.cfg.transaction_cfg)
        log.info('Submit param proposal result : {}'.format(result))
        assert_code(result, 3)

        result = pip.submitParam(pip.node.node_id, str(time.time()), 'staking', 'unStakeFreezeDuration', '10.5',
                                     pip.node.staking_address, transaction_cfg=pip.cfg.transaction_cfg)
        log.info('Submit param proposal result : {}'.format(result))
        assert_code(result, 3)

        result = pip.submitParam(pip.node.node_id, str(time.time()), 'staking', 'unStakeFreezeDuration', '-100',
                                     pip.node.staking_address, transaction_cfg=pip.cfg.transaction_cfg)
        log.info('Submit param proposal result : {}'.format(result))
        assert_code(result, 3)

        result = pip.submitParam(pip.node.node_id, str(time.time()), 'staking', 'unStakeFreezeDuration', '0',
                                     pip.node.staking_address, transaction_cfg=pip.cfg.transaction_cfg)
        log.info('Submit param proposal result : {}'.format(result))
        assert_code(result, 3)

        result = pip.submitParam(pip.node.node_id, str(time.time()), 'staking', 'unStakeFreezeDuration', 11,
                                     pip.node.staking_address, transaction_cfg=pip.cfg.transaction_cfg)
        log.info('Submit param proposal result : {}'.format(result))
        assert_code(result, 3)

        result = pip.submitParam(pip.node.node_id, str(time.time()), 'staking', 'unStakeFreezeDuration',
                                     '113', pip.node.staking_address, transaction_cfg=pip.cfg.transaction_cfg)
        log.info('Submit param proposal result : {}'.format(result))
        assert_code(result, 3)

        result = pip.submitParam(pip.node.node_id, str(time.time()), 'staking', 'unStakeFreezeDuration',
                                     str(get_governable_parameter_value(client, 'unStakeFreezeDuration')),
                                     pip.node.staking_address, transaction_cfg=pip.cfg.transaction_cfg)
        log.info('Submit param proposal result : {}'.format(result))
        assert_code(result, 302034)

        result = pip.submitParam(pip.node.node_id, str(time.time()), 'staking', 'unStakeFreezeDuration',
                                     str(get_governable_parameter_value(client, 'maxEvidenceAge')),
                                     pip.node.staking_address, transaction_cfg=pip.cfg.transaction_cfg)
        log.info('Submit param proposal result : {}'.format(result))
        assert_code(result, 3)

        if int(get_governable_parameter_value(client, 'unStakeFreezeDuration')) != 112:
            result = pip.submitParam(pip.node.node_id, str(time.time()), 'staking', 'unStakeFreezeDuration', '112',
                                         pip.node.staking_address, transaction_cfg=pip.cfg.transaction_cfg)
            log.info('Submit param proposal result : {}'.format(result))
            assert_code(result, 0)

    @pytest.mark.P0
    @allure.title('Submit parammeter  proposal function verification')
    def test__PP_SU_014(self, no_vp_proposal, all_clients):
        pip = no_vp_proposal
        if int(get_governable_parameter_value(all_clients[0], 'unStakeFreezeDuration')) != str(
                int(get_governable_parameter_value(all_clients[0], 'maxEvidenceAge')) - 1):
            result = pip.submitParam(pip.node.node_id, str(time.time()), 'staking', 'unStakeFreezeDuration',
                                     str(int(get_governable_parameter_value(all_clients[0], 'maxEvidenceAge')) + 5),
                                     pip.node.staking_address, transaction_cfg=pip.cfg.transaction_cfg)
            log.info('Submit param proposal result : {}'.format(result))
            assert_code(result, 0)

    @pytest.mark.P0
    @allure.title('Submit parammeter  proposal function verification')
    def test_PP_SU_015_PP_SU_016(self, no_vp_proposal, all_clients):
        pip = no_vp_proposal
        client = get_client_by_nodeid(pip.node.node_id, all_clients)
        result = pip.submitParam(pip.node.node_id, str(time.time()), 'staking', 'maxValidators', '',
                                     pip.node.staking_address, transaction_cfg=pip.cfg.transaction_cfg)
        log.info('Submit param proposal result : {}'.format(result))
        assert_code(result, 3)

        result = pip.submitParam(pip.node.node_id, str(time.time()), 'staking', 'maxValidators', '30.5',
                                     pip.node.staking_address, transaction_cfg=pip.cfg.transaction_cfg)
        log.info('Submit param proposal result : {}'.format(result))
        assert_code(result, 3)

        result = pip.submitParam(pip.node.node_id, str(time.time()), 'staking', 'maxValidators', '-100',
                                     pip.node.staking_address, transaction_cfg=pip.cfg.transaction_cfg)
        log.info('Submit param proposal result : {}'.format(result))
        assert_code(result, 3)

        result = pip.submitParam(pip.node.node_id, str(time.time()), 'staking', 'maxValidators', '0',
                                     pip.node.staking_address, transaction_cfg=pip.cfg.transaction_cfg)
        log.info('Submit param proposal result : {}'.format(result))
        assert_code(result, 3)

        result = pip.submitParam(pip.node.node_id, str(time.time()), 'staking', 'maxValidators', 25,
                                     pip.node.staking_address, transaction_cfg=pip.cfg.transaction_cfg)
        log.info('Submit param proposal result : {}'.format(result))
        assert_code(result, 3)

        result = pip.submitParam(pip.node.node_id, str(time.time()), 'staking', 'maxValidators',
                                     '3', pip.node.staking_address, transaction_cfg=pip.cfg.transaction_cfg)
        log.info('Submit param proposal result : {}'.format(result))
        assert_code(result, 3)

        result = pip.submitParam(pip.node.node_id, str(time.time()), 'staking', 'maxValidators',
                                     '202', pip.node.staking_address, transaction_cfg=pip.cfg.transaction_cfg)
        log.info('Submit param proposal result : {}'.format(result))
        assert_code(result, 3)

        result = pip.submitParam(pip.node.node_id, str(time.time()), 'staking', 'maxValidators',
                                     str(get_governable_parameter_value(client, 'maxValidators')),
                                     pip.node.staking_address, transaction_cfg=pip.cfg.transaction_cfg)
        log.info('Submit param proposal result : {}'.format(result))
        assert_code(result, 302034)

        if int(get_governable_parameter_value(client, 'maxValidators')) != 4:
            result = pip.submitParam(pip.node.node_id, str(time.time()), 'staking', 'maxValidators', '4',
                                         pip.node.staking_address, transaction_cfg=pip.cfg.transaction_cfg)
            log.info('Submit param proposal result : {}'.format(result))
            assert_code(result, 0)

    @pytest.mark.P0
    @allure.title('Submit parammeter  proposal function verification')
    def test_PP_SU_016(self, no_vp_proposal, all_clients):
        pip = no_vp_proposal
        if int(get_governable_parameter_value(all_clients[0], 'maxValidators')) != 201:
            result = pip.submitParam(pip.node.node_id, str(time.time()), 'staking', 'maxValidators', '201',
                                         pip.node.staking_address, transaction_cfg=pip.cfg.transaction_cfg)
            log.info('Submit param proposal result : {}'.format(result))
            assert_code(result, 0)

    @pytest.mark.P0
    @allure.title('Submit parammeter  proposal function verification')
    def test_PP_SU_016(self, no_vp_proposal, all_clients):
        pip = no_vp_proposal
        client = get_client_by_nodeid(pip.node.node_id, all_clients)
        if int(get_governable_parameter_value(client, 'maxValidators')) != 201:
            result = pip.submitParam(pip.node.node_id, str(time.time()), 'staking', 'maxValidators', '201',
                                         pip.node.staking_address, transaction_cfg=pip.cfg.transaction_cfg)
            log.info('Submit param proposal result : {}'.format(result))
            assert_code(result, 0)

    @pytest.mark.P0
    @allure.title('Submit parammeter  proposal function verification')
    def test_PP_SU_017_PP_SU_018(self, no_vp_proposal, all_clients):
        pip = no_vp_proposal
        client = get_client_by_nodeid(pip.node.node_id, all_clients)
        result = pip.submitParam(pip.node.node_id, str(time.time()), 'block', 'maxBlockGasLimit', '',
                                     pip.node.staking_address, transaction_cfg=pip.cfg.transaction_cfg)
        log.info('Submit param proposal result : {}'.format(result))
        assert_code(result, 3)

        result = pip.submitParam(pip.node.node_id, str(time.time()), 'block', 'maxBlockGasLimit', '4712388.5',
                                     pip.node.staking_address, transaction_cfg=pip.cfg.transaction_cfg)
        log.info('Submit param proposal result : {}'.format(result))
        assert_code(result, 3)

        result = pip.submitParam(pip.node.node_id, str(time.time()), 'block', 'maxBlockGasLimit', '-4712388',
                                     pip.node.staking_address, transaction_cfg=pip.cfg.transaction_cfg)
        log.info('Submit param proposal result : {}'.format(result))
        assert_code(result, 3)

        result = pip.submitParam(pip.node.node_id, str(time.time()), 'block', 'maxBlockGasLimit', '0',
                                     pip.node.staking_address, transaction_cfg=pip.cfg.transaction_cfg)
        log.info('Submit param proposal result : {}'.format(result))
        assert_code(result, 3)

        result = pip.submitParam(pip.node.node_id, str(time.time()), 'block', 'maxBlockGasLimit', 4712388,
                                     pip.node.staking_address, transaction_cfg=pip.cfg.transaction_cfg)
        log.info('Submit param proposal result : {}'.format(result))
        assert_code(result, 3)

        result = pip.submitParam(pip.node.node_id, str(time.time()), 'block', 'maxBlockGasLimit',
                                     '4712387', pip.node.staking_address, transaction_cfg=pip.cfg.transaction_cfg)
        log.info('Submit param proposal result : {}'.format(result))
        assert_code(result, 3)

        result = pip.submitParam(pip.node.node_id, str(time.time()), 'block', 'maxBlockGasLimit',
                                     '210000001', pip.node.staking_address, transaction_cfg=pip.cfg.transaction_cfg)
        log.info('Submit param proposal result : {}'.format(result))
        assert_code(result, 3)

        result = pip.submitParam(pip.node.node_id, str(time.time()), 'block', 'maxBlockGasLimit',
                                     str(get_governable_parameter_value(client, 'maxBlockGasLimit')),
                                     pip.node.staking_address, transaction_cfg=pip.cfg.transaction_cfg)
        log.info('Submit param proposal result : {}'.format(result))
        assert_code(result, 302034)

        if int(get_governable_parameter_value(client, 'maxBlockGasLimit')) != 4712388:
            result = pip.submitParam(pip.node.node_id, str(time.time()), 'block', 'maxBlockGasLimit', '4712388',
                                         pip.node.staking_address, transaction_cfg=pip.cfg.transaction_cfg)
            log.info('Submit param proposal result : {}'.format(result))
            assert_code(result, 0)



    @pytest.mark.P0
    @allure.title('Submit parammeter  proposal function verification')
    def test_PP_SU_018(self, no_vp_proposal):
        pip = no_vp_proposal
        result = pip.submitParam(pip.node.node_id, str(time.time()), 'block', 'maxBlockGasLimit', '210000000',
                                     pip.node.staking_address, transaction_cfg=pip.cfg.transaction_cfg)
        log.info('Submit param proposal result : {}'.format(result))
        assert_code(result, 0)

    @pytest.mark.P0
    @allure.title('Submit parammeter  proposal function verification')
    def test_PP_SU_022(self, new_genesis_env, clients_consensus):
        genesis = from_dict(data_class=Genesis, data=new_genesis_env.genesis_config)
        genesis.economicModel.slashing.zeroProduceNumberThreshold = 2
        genesis.economicModel.slashing.zeroProduceCumulativeTime = 3
        new_genesis_env.set_genesis(genesis.to_dict())
        new_genesis_env.deploy_all()
        client = clients_consensus[0]
        pip = client.pip
        result = pip.submitParam(pip.node.node_id, str(time.time()), 'slashing', 'zeroProduceCumulativeTime', '',
                                 pip.node.staking_address, transaction_cfg=pip.cfg.transaction_cfg)
        log.info('Submit param proposal result : {}'.format(result))
        assert_code(result, 3)

        result = pip.submitParam(pip.node.node_id, str(time.time()), 'slashing', 'zeroProduceCumulativeTime', '1.1',
                                 pip.node.staking_address, transaction_cfg=pip.cfg.transaction_cfg)
        log.info('Submit param proposal result : {}'.format(result))
        assert_code(result, 3)

        result = pip.submitParam(pip.node.node_id, str(time.time()), 'slashing', 'zeroProduceCumulativeTime', '-1',
                                 pip.node.staking_address, transaction_cfg=pip.cfg.transaction_cfg)
        log.info('Submit param proposal result : {}'.format(result))
        assert_code(result, 3)

        result = pip.submitParam(pip.node.node_id, str(time.time()), 'slashing', 'zeroProduceCumulativeTime', '0',
                                 pip.node.staking_address, transaction_cfg=pip.cfg.transaction_cfg)
        log.info('Submit param proposal result : {}'.format(result))
        assert_code(result, 3)

        result = pip.submitParam(pip.node.node_id, str(time.time()), 'slashing', 'zeroProduceCumulativeTime', 4,
                                 pip.node.staking_address, transaction_cfg=pip.cfg.transaction_cfg)
        log.info('Submit param proposal result : {}'.format(result))
        assert_code(result, 3)

        result = pip.submitParam(pip.node.node_id, str(time.time()), 'slashing', 'zeroProduceCumulativeTime',
                                 str(pip.economic.consensus_wheel + 1), pip.node.staking_address,
                                 transaction_cfg=pip.cfg.transaction_cfg)
        log.info('Submit param proposal result : {}'.format(result))
        assert_code(result, 3)


        result = pip.submitParam(pip.node.node_id, str(time.time()), 'slashing', 'zeroProduceCumulativeTime',
                                 str(int(get_governable_parameter_value(client, 'zeroProduceNumberThreshold')) - 1),
                                 pip.node.staking_address, transaction_cfg=pip.cfg.transaction_cfg)
        log.info('Submit param proposal result : {}'.format(result))
        assert_code(result, 3)


        result = pip.submitParam(pip.node.node_id, str(time.time()), 'slashing', 'zeroProduceCumulativeTime',
                                 str(pip.economic.consensus_wheel), pip.node.staking_address,
                                 transaction_cfg=pip.cfg.transaction_cfg)
        log.info('Submit param proposal result : {}'.format(result))
        assert_code(result, 0)

    @pytest.mark.P0
    @allure.title('Submit parammeter  proposal function verification')
    def test_PP_SU_023(self, new_genesis_env, clients_consensus):
        genesis = from_dict(data_class=Genesis, data=new_genesis_env.genesis_config)
        genesis.economicModel.slashing.zeroProduceNumberThreshold = 2
        genesis.economicModel.slashing.zeroProduceCumulativeTime = 3
        new_genesis_env.set_genesis(genesis.to_dict())
        new_genesis_env.deploy_all()
        client = clients_consensus[0]
        pip = client.pip
        result = pip.submitParam(pip.node.node_id, str(time.time()), 'slashing', 'zeroProduceNumberThreshold', '',
                                 pip.node.staking_address, transaction_cfg=pip.cfg.transaction_cfg)
        log.info('Submit param proposal result : {}'.format(result))
        assert_code(result, 3)

        result = pip.submitParam(pip.node.node_id, str(time.time()), 'slashing', 'zeroProduceNumberThreshold', '1.1',
                                 pip.node.staking_address, transaction_cfg=pip.cfg.transaction_cfg)
        log.info('Submit param proposal result : {}'.format(result))
        assert_code(result, 3)

        result = pip.submitParam(pip.node.node_id, str(time.time()), 'slashing', 'zeroProduceNumberThreshold', '-2',
                                 pip.node.staking_address, transaction_cfg=pip.cfg.transaction_cfg)
        log.info('Submit param proposal result : {}'.format(result))
        assert_code(result, 3)

        result = pip.submitParam(pip.node.node_id, str(time.time()), 'slashing', 'zeroProduceNumberThreshold', '0',
                                 pip.node.staking_address, transaction_cfg=pip.cfg.transaction_cfg)
        log.info('Submit param proposal result : {}'.format(result))
        assert_code(result, 3)

        result = pip.submitParam(pip.node.node_id, str(time.time()), 'slashing', 'zeroProduceNumberThreshold', 1,
                                 pip.node.staking_address, transaction_cfg=pip.cfg.transaction_cfg)
        log.info('Submit param proposal result : {}'.format(result))
        assert_code(result, 3)

        result = pip.submitParam(pip.node.node_id, str(time.time()), 'slashing', 'zeroProduceNumberThreshold',
                                 str(int(get_governable_parameter_value(client, 'zeroProduceCumulativeTime')) + 1),
                                 pip.node.staking_address, transaction_cfg=pip.cfg.transaction_cfg)
        log.info('Submit param proposal result : {}'.format(result))
        assert_code(result, 3)

        result = pip.submitParam(pip.node.node_id, str(time.time()), 'slashing', 'zeroProduceNumberThreshold',
                                 str(get_governable_parameter_value(client, 'zeroProduceCumulativeTime')),
                                     pip.node.staking_address,
                                 transaction_cfg=pip.cfg.transaction_cfg)
        log.info('Submit param proposal result : {}'.format(result))
        assert_code(result, 0)

    @pytest.mark.P0
    @allure.title('Submit parammeter  proposal function verification')
    def test_PP_SU_023_2(self, new_genesis_env, clients_consensus):
        genesis = from_dict(data_class=Genesis, data=new_genesis_env.genesis_config)
        genesis.economicModel.slashing.zeroProduceNumberThreshold = 2
        genesis.economicModel.slashing.zeroProduceCumulativeTime = 3
        new_genesis_env.set_genesis(genesis.to_dict())
        new_genesis_env.deploy_all()
        client = clients_consensus[0]
        pip = client.pip
        result = pip.submitParam(pip.node.node_id, str(time.time()), 'slashing', 'zeroProduceNumberThreshold',
                                 '1', pip.node.staking_address,
                                 transaction_cfg=pip.cfg.transaction_cfg)
        assert_code(result, 0)

    @pytest.mark.P0
    @allure.title('Submit parammeter  proposal function verification')
    def test_PP_SU_026(self, no_vp_proposal, client_consensus):
        pip = no_vp_proposal
        result = pip.submitParam(pip.node.node_id, str(time.time()), 'staking', 'rewardPerMaxChangeRange', '',
                                 pip.node.staking_address, transaction_cfg=pip.cfg.transaction_cfg)
        log.info('Submit param proposal result : {}'.format(result))
        assert_code(result, 3)

        result = pip.submitParam(pip.node.node_id, str(time.time()), 'staking', 'rewardPerMaxChangeRange', '1.1',
                                 pip.node.staking_address, transaction_cfg=pip.cfg.transaction_cfg)
        log.info('Submit param proposal result : {}'.format(result))
        assert_code(result, 3)

        result = pip.submitParam(pip.node.node_id, str(time.time()), 'staking', 'rewardPerMaxChangeRange', '-2',
                                 pip.node.staking_address, transaction_cfg=pip.cfg.transaction_cfg)
        log.info('Submit param proposal result : {}'.format(result))
        assert_code(result, 3)

        result = pip.submitParam(pip.node.node_id, str(time.time()), 'staking', 'rewardPerMaxChangeRange', '0',
                                 pip.node.staking_address, transaction_cfg=pip.cfg.transaction_cfg)
        log.info('Submit param proposal result : {}'.format(result))
        assert_code(result, 3)

        result = pip.submitParam(pip.node.node_id, str(time.time()), 'staking', 'rewardPerMaxChangeRange', 6,
                                 pip.node.staking_address, transaction_cfg=pip.cfg.transaction_cfg)
        log.info('Submit param proposal result : {}'.format(result))
        assert_code(result, 3)

        result = pip.submitParam(pip.node.node_id, str(time.time()), 'staking', 'rewardPerMaxChangeRange',
                                 '2001',
                                 pip.node.staking_address, transaction_cfg=pip.cfg.transaction_cfg)
        log.info('Submit param proposal result : {}'.format(result))
        assert_code(result, 3)

        result = pip.submitParam(pip.node.node_id, str(time.time()), 'staking', 'rewardPerMaxChangeRange',
                                 '2000',
                                 pip.node.staking_address,
                                 transaction_cfg=pip.cfg.transaction_cfg)
        log.info('Submit param proposal result : {}'.format(result))
        assert_code(result, 0)

    @pytest.mark.P0
    @allure.title('Submit parammeter  proposal function verification')
    def test_PP_SU_026_2(self, new_genesis_env, clients_consensus):
        genesis = from_dict(data_class=Genesis, data=new_genesis_env.genesis_config)
        genesis.economicModel.staking.rewardPerMaxChangeRange = 2
        new_genesis_env.set_genesis(genesis.to_dict())
        new_genesis_env.deploy_all()
        client = clients_consensus[0]
        pip = client.pip

        result = pip.submitParam(pip.node.node_id, str(time.time()), 'staking', 'rewardPerMaxChangeRange',
                                 '2', pip.node.staking_address,
                                 transaction_cfg=pip.cfg.transaction_cfg)
        assert_code(result, 302034)
        result = pip.submitParam(pip.node.node_id, str(time.time()), 'staking', 'rewardPerMaxChangeRange',
                                 '1', pip.node.staking_address,
                                 transaction_cfg=pip.cfg.transaction_cfg)
        assert_code(result, 0)

    @pytest.mark.P0
    @allure.title('Submit parammeter  proposal function verification')
    def test_PP_SU_027(self, no_vp_proposal, client_consensus):
        pip = no_vp_proposal
        result = pip.submitParam(pip.node.node_id, str(time.time()), 'staking', 'rewardPerChangeInterval', '',
                                 pip.node.staking_address, transaction_cfg=pip.cfg.transaction_cfg)
        log.info('Submit param proposal result : {}'.format(result))
        assert_code(result, 3)

        result = pip.submitParam(pip.node.node_id, str(time.time()), 'staking', 'rewardPerChangeInterval', '1.1',
                                 pip.node.staking_address, transaction_cfg=pip.cfg.transaction_cfg)
        log.info('Submit param proposal result : {}'.format(result))
        assert_code(result, 3)

        result = pip.submitParam(pip.node.node_id, str(time.time()), 'staking', 'rewardPerChangeInterval', '-2',
                                 pip.node.staking_address, transaction_cfg=pip.cfg.transaction_cfg)
        log.info('Submit param proposal result : {}'.format(result))
        assert_code(result, 3)

        result = pip.submitParam(pip.node.node_id, str(time.time()), 'staking', 'rewardPerChangeInterval', '1',
                                 pip.node.staking_address, transaction_cfg=pip.cfg.transaction_cfg)
        log.info('Submit param proposal result : {}'.format(result))
        assert_code(result, 3)

        result = pip.submitParam(pip.node.node_id, str(time.time()), 'staking', 'rewardPerChangeInterval', '0',
                                 pip.node.staking_address, transaction_cfg=pip.cfg.transaction_cfg)
        log.info('Submit param proposal result : {}'.format(result))
        assert_code(result, 3)

        result = pip.submitParam(pip.node.node_id, str(time.time()), 'staking', 'rewardPerChangeInterval', 6,
                                 pip.node.staking_address, transaction_cfg=pip.cfg.transaction_cfg)
        log.info('Submit param proposal result : {}'.format(result))
        assert_code(result, 3)

        result = pip.submitParam(pip.node.node_id, str(time.time()), 'staking', 'rewardPerChangeInterval',
                                 str(pip.economic.additional_cycle_time * 60 //(
                                         pip.economic.settlement_size * pip.economic.interval) + 1),
                                 pip.node.staking_address, transaction_cfg=pip.cfg.transaction_cfg)
        log.info('Submit param proposal result : {}'.format(result))
        assert_code(result, 0)

        result = pip.submitParam(pip.node.node_id, str(time.time()), 'staking', 'rewardPerChangeInterval',
                                 str(pip.economic.additional_cycle_time * 60 //(
                                         pip.economic.settlement_size * pip.economic.interval)),
                                 pip.node.staking_address,
                                 transaction_cfg=pip.cfg.transaction_cfg)
        log.info('Submit param proposal result : {}'.format(result))
        assert_code(result, 302032)

    @pytest.mark.P0
    @allure.title('Submit parammeter  proposal function verification')
    def test_PP_SU_027_2(self, new_genesis_env, clients_consensus):
        genesis = from_dict(data_class=Genesis, data=new_genesis_env.genesis_config)
        genesis.economicModel.staking.rewardPerChangeInterval = 3
        new_genesis_env.set_genesis(genesis.to_dict())
        new_genesis_env.deploy_all()
        client = clients_consensus[0]
        pip = client.pip
        result = pip.submitParam(pip.node.node_id, str(time.time()), 'staking', 'rewardPerChangeInterval',
                                 '3', pip.node.staking_address,
                                 transaction_cfg=pip.cfg.transaction_cfg)
        assert_code(result, 302034)
        result = pip.submitParam(pip.node.node_id, str(time.time()), 'staking', 'rewardPerChangeInterval',
                                 '2', pip.node.staking_address,
                                 transaction_cfg=pip.cfg.transaction_cfg)
        assert_code(result, 0)

    @pytest.mark.P0
    @allure.title('Submit parammeter  proposal function verification')
    def test_PP_SU_028(self, no_vp_proposal, client_consensus):
        pip = no_vp_proposal
        result = pip.submitParam(pip.node.node_id, str(time.time()), 'reward', 'increaseIssuanceRatio', '',
                                 pip.node.staking_address, transaction_cfg=pip.cfg.transaction_cfg)
        log.info('Submit param proposal result : {}'.format(result))
        assert_code(result, 3)
        result = pip.submitParam(pip.node.node_id, str(time.time()), 'reward', 'increaseIssuanceRatio', '1.1',
                                 pip.node.staking_address, transaction_cfg=pip.cfg.transaction_cfg)
        log.info('Submit param proposal result : {}'.format(result))
        assert_code(result, 3)
        result = pip.submitParam(pip.node.node_id, str(time.time()), 'reward', 'increaseIssuanceRatio', '-2',
                                 pip.node.staking_address, transaction_cfg=pip.cfg.transaction_cfg)
        log.info('Submit param proposal result : {}'.format(result))
        assert_code(result, 3)
        result = pip.submitParam(pip.node.node_id, str(time.time()), 'reward', 'increaseIssuanceRatio', '2001',
                                 pip.node.staking_address, transaction_cfg=pip.cfg.transaction_cfg)
        log.info('Submit param proposal result : {}'.format(result))
        assert_code(result, 3)
        result = pip.submitParam(pip.node.node_id, str(time.time()), 'reward', 'increaseIssuanceRatio', 6,
                                 pip.node.staking_address, transaction_cfg=pip.cfg.transaction_cfg)
        log.info('Submit param proposal result : {}'.format(result))
        assert_code(result, 3)
        result = pip.submitParam(pip.node.node_id, str(time.time()), 'reward', 'increaseIssuanceRatio', '2000',
                                 pip.node.staking_address,
                                 transaction_cfg=pip.cfg.transaction_cfg)
        log.info('Submit param proposal result : {}'.format(result))
        assert_code(result, 0)

    @pytest.mark.P0
    @allure.title('Submit parammeter  proposal function verification')
    def test_PP_SU_028_2(self, new_genesis_env, clients_consensus):
        genesis = from_dict(data_class=Genesis, data=new_genesis_env.genesis_config)
        genesis.economicModel.reward.increaseIssuanceRatio = 3
        new_genesis_env.set_genesis(genesis.to_dict())
        new_genesis_env.deploy_all()
        client = clients_consensus[0]
        pip = client.pip

        result = pip.submitParam(pip.node.node_id, str(time.time()), 'reward', 'increaseIssuanceRatio',
                                 '3', pip.node.staking_address,
                                 transaction_cfg=pip.cfg.transaction_cfg)
        assert_code(result, 302034)
        result = pip.submitParam(pip.node.node_id, str(time.time()), 'reward', 'increaseIssuanceRatio',
                                 '0', pip.node.staking_address,
                                 transaction_cfg=pip.cfg.transaction_cfg)
        assert_code(result, 0)

    @pytest.mark.P2
    @allure.title('Submit parammeter  proposal function verification')
    def test_PP_SU_024_UP_PA_008_PP_VO_004(self, new_genesis_env, clients_consensus):
        genesis = from_dict(data_class=Genesis, data=new_genesis_env.genesis_config)
        genesis.economicModel.gov.paramProposalVoteDurationSeconds = 0
        new_genesis_env.set_genesis(genesis.to_dict())
        new_genesis_env.deploy_all()
        client = clients_consensus[0]
        pip = client.pip

        result = pip.submitParam(pip.node.node_id, str(time.time()), 'slashing', 'zeroProduceNumberThreshold',
                                 '4',pip.node.staking_address, transaction_cfg=pip.cfg.transaction_cfg)
        log.info('zeroProduceNumberThreshold {} submit param proposal result :{}'.format(4, result))
        assert_code(result, 3)

        result = pip.submitParam(pip.node.node_id, str(time.time()), 'slashing', 'zeroProduceCumulativeTime',
                                 '4', pip.node.staking_address, transaction_cfg=pip.cfg.transaction_cfg)
        assert_code(result, 0)
        proposal_info = pip.get_effect_proposal_info_of_vote(pip.cfg.param_proposal)
        log.info('Param praposal info : {}'.format(proposal_info))
        for client in clients_consensus:
            pip = client.pip
            result = pip.vote(pip.node.node_id, proposal_info.get('ProposalID'), pip.cfg.vote_option_yeas,
                              pip.node.staking_address)
            log.info('node {} vote result {}'.format(pip.node.node_id, result))
            assert_code(result, 0)
        wait_block_number(pip.node, proposal_info.get('EndVotingBlock'))
        value = client.pip.pip.getGovernParamValue('slashing', 'zeroProduceCumulativeTime').get('Ret')
        log.info('zeroProduceCumulativeTime new value : {}'.format(value))
        assert int(value) == 4

        value, oldvalue = get_governable_parameter_value(client, 'zeroProduceCumulativeTime', flag=1)
        assert value == 4
        assert oldvalue == 1

        result = pip.submitParam(pip.node.node_id, str(time.time()), 'slashing', 'zeroProduceNumberThreshold',
                                 '4', pip.node.staking_address, transaction_cfg=pip.cfg.transaction_cfg)
        assert_code(result, 0)

    @pytest.mark.P2
    @allure.title('Submit parammeter  proposal function verification')
    def test_PP_SU_025_UP_PA_009_PP_VO_005(self, new_genesis_env, clients_consensus):
        genesis = from_dict(data_class=Genesis, data=new_genesis_env.genesis_config)
        genesis.economicModel.gov.paramProposalVoteDurationSeconds = 0
        genesis.economicModel.slashing.zeroProduceNumberThreshold = 2
        genesis.economicModel.slashing.zeroProduceCumulativeTime = 3
        new_genesis_env.set_genesis(genesis.to_dict())
        new_genesis_env.deploy_all()
        client = clients_consensus[0]
        pip = client.pip
        result = pip.submitParam(pip.node.node_id, str(time.time()), 'slashing', 'zeroProduceCumulativeTime',
                                 '1',pip.node.staking_address, transaction_cfg=pip.cfg.transaction_cfg)
        log.info('zeroProduceCumulativeTime {} submit param proposal result :{}'.format(1, result))
        assert_code(result, 3)
        result = pip.submitParam(pip.node.node_id, str(time.time()), 'slashing', 'zeroProduceNumberThreshold',
                                 '1', pip.node.staking_address, transaction_cfg=pip.cfg.transaction_cfg)
        assert_code(result, 0)
        proposal_info = pip.get_effect_proposal_info_of_vote(pip.cfg.param_proposal)
        log.info('Param praposal info : {}'.format(proposal_info))
        for client in clients_consensus:
            pip = client.pip
            result = pip.vote(pip.node.node_id, proposal_info.get('ProposalID'), pip.cfg.vote_option_yeas,
                              pip.node.staking_address)
            log.info('node {} vote result {}'.format(pip.node.node_id, result))
            assert_code(result, 0)
        wait_block_number(pip.node, proposal_info.get('EndVotingBlock'))
        value = client.pip.pip.getGovernParamValue('slashing', 'zeroProduceNumberThreshold').get('Ret')
        log.info('zeroProduceNumberThreshold new value : {}'.format(value))
        assert int(value) == 1
        value, oldvalue = get_governable_parameter_value(client, 'zeroProduceNumberThreshold', flag=1)
        assert value == 1
        assert oldvalue == 2
        result = pip.submitParam(pip.node.node_id, str(time.time()), 'slashing', 'zeroProduceCumulativeTime',
                                 '1', pip.node.staking_address, transaction_cfg=pip.cfg.transaction_cfg)
        assert_code(result, 0)


    @pytest.mark.P0
    @allure.title('Submit parammeter  proposal function verification')
    def test_PP_SU_019(self, no_vp_proposal):
        pip = no_vp_proposal
        result = pip.submitParam(pip.node.node_id, str(time.time()), 'block', '', '1',
                                     pip.node.staking_address,
                                     transaction_cfg=pip.cfg.transaction_cfg)
        log.info('Submit param proposal result : {}'.format(result))
        assert_code(result, 302031)
        result = pip.submitParam(pip.node.node_id, str(time.time()), 'block', 'unStakeFreezeDuration', '100',
                                     pip.node.staking_address,
                                     transaction_cfg=pip.cfg.transaction_cfg)
        log.info('Submit param proposal result : {}'.format(result))
        assert_code(result, 302031)
        result = pip.submitParam(pip.node.node_id, str(time.time()), 'slashing', 'SlashBlocksReward', '100',
                                     pip.node.staking_address,
                                     transaction_cfg=pip.cfg.transaction_cfg)
        log.info('Submit param proposal result : {}'.format(result))
        assert_code(result, 302031)
        result = pip.submitParam(pip.node.node_id, str(time.time()), 'slashing', 'slash BlocksReward', '100',
                                     pip.node.staking_address,
                                     transaction_cfg=pip.cfg.transaction_cfg)
        log.info('Submit param proposal result : {}'.format(result))
        assert_code(result, 302031)
        result = pip.submitParam(pip.node.node_id, str(time.time()), 'slashing', 'slashBlocks./,.Reward', '100',
                                     pip.node.staking_address,
                                     transaction_cfg=pip.cfg.transaction_cfg)
        log.info('Submit param proposal result : {}'.format(result))
        assert_code(result, 302031)

    @pytest.mark.P0
    @allure.title('Submit parammeter  proposal function verification')
    def test_PP_SU_020(self, no_vp_proposal):
        pip = no_vp_proposal
        result = pip.submitParam(pip.node.node_id, str(time.time()), '', 'slashBlocksReward', '100',
                                     pip.node.staking_address,
                                     transaction_cfg=pip.cfg.transaction_cfg)
        log.info('Submit param proposal result : {}'.format(result))
        assert_code(result, 302031)

        result = pip.submitParam(pip.node.node_id, str(time.time()), 'Slashing', 'slashBlocksReward', '100',
                                     pip.node.staking_address,
                                     transaction_cfg=pip.cfg.transaction_cfg)
        log.info('Submit param proposal result : {}'.format(result))
        assert_code(result, 302031)

        result = pip.submitParam(pip.node.node_id, str(time.time()), 'sLashing', 'SlashBlocksReward', '100',
                                     pip.node.staking_address,
                                     transaction_cfg=pip.cfg.transaction_cfg)
        log.info('Submit param proposal result : {}'.format(result))
        assert_code(result, 302031)

        result = pip.submitParam(pip.node.node_id, str(time.time()), 'slashing123', 'slash BlocksReward', '100',
                                     pip.node.staking_address,
                                     transaction_cfg=pip.cfg.transaction_cfg)
        log.info('Submit param proposal result : {}'.format(result))
        assert_code(result, 302031)

        result = pip.submitParam(pip.node.node_id, str(time.time()), 's lashing', 'slashBlocks./,.Reward', '100',
                                     pip.node.staking_address,
                                     transaction_cfg=pip.cfg.transaction_cfg)
        log.info('Submit param proposal result : {}'.format(result))
        assert_code(result, 302031)

        result = pip.submitParam(pip.node.node_id, str(time.time()), 's.,.lashing', 'slashBlocks./,.Reward', '100',
                                     pip.node.staking_address,
                                     transaction_cfg=pip.cfg.transaction_cfg)
        log.info('Submit param proposal result : {}'.format(result))
        assert_code(result, 302031)


class TestSubmitPPAbnormal:
    @pytest.mark.P0
    @allure.title('New node submit parammeter  proposal function verification')
    def test_PP_PR_002(self, no_vp_proposal, client_new_node):
        pip = client_new_node.pip
        address, _ = pip.economic.account.generate_account(pip.node.web3, 10**18 * 10000)
        result = pip.submitParam(pip.node.node_id, str(time.time()), 'slashing', 'slashBlocksReward', '88',
                                     address, transaction_cfg=pip.cfg.transaction_cfg)
        log.info('new node submit param proposal result : {}'.format(result))
        assert_code(result, 302022)

    @pytest.mark.P0
    @allure.title('Candidate submit parammeter  proposal function verification')
    def test_PP_PR_001(self, no_vp_proposal, client_candidate):
        pip = client_candidate.pip
        result = pip.submitParam(pip.node.node_id, str(time.time()), 'slashing', 'slashBlocksReward', '87',
                                     pip.node.staking_address, transaction_cfg=pip.cfg.transaction_cfg)
        log.info('candidate submit param proposal result :{}'.format(result))
        assert_code(result, 302022)

    @pytest.mark.P0
    @allure.title('There is a voting version proposal, submit parammeter  proposal function verification')
    def test_PP_UN_002(self, submit_version):
        pip = submit_version
        result = pip.submitParam(pip.node.node_id, str(time.time()), 'slashing', 'slashBlocksReward', '99',
                                     pip.node.staking_address, transaction_cfg=pip.cfg.transaction_cfg)
        log.info('There is voting version proposal, submit a param proposal : {}'.format(result))
        assert_code(result, 302012)

    @pytest.mark.P2
    @allure.title('Abnormal submit parammeter  proposal function verification')
    def test_PP_PR_003_PP_PR_004(self, no_vp_proposal, all_clients):
        pip = no_vp_proposal
        client = get_client_by_nodeid(pip.node.node_id, all_clients)
        address = pip.node.staking_address
        result = client.staking.withdrew_staking(address)
        log.info('nodeid: {} withdrewstaking result: {}'.format(client.node.node_id, result))
        assert_code(result, 0)
        result = pip.submitParam(pip.node.node_id, str(time.time()), 'slashing', 'slashBlocksReward',
                                     '86', address, transaction_cfg=pip.cfg.transaction_cfg)
        log.info('node exiting，param proposal result: {}'.format(result))
        assert_code(result, 302020)

        client.economic.wait_settlement_blocknum(client.node,
                                                     number=client.economic.unstaking_freeze_ratio)
        result = pip.submitParam(pip.node.node_id, str(time.time()), 'slashing', 'slashBlocksReward',
                                     '86', address, transaction_cfg=pip.cfg.transaction_cfg)
        log.info('exited node，cancel proposal result: {}'.format(result))
        assert_code(result, 302022)

    @pytest.mark.P2
    @allure.title('Not staking address submit parammeter  proposal function verification')
    def test_PP_WA_001(self, no_vp_proposal):
        pip = no_vp_proposal
        address, _ = pip.economic.account.generate_account(pip.node.web3, 10**18 * 10000)
        result = pip.submitParam(pip.node.node_id, str(time.time()), 'slashing', 'slashBlocksReward', '87',
                                     address, transaction_cfg=pip.cfg.transaction_cfg)
        log.info('candidate submit param proposal result :{}'.format(result))
        assert_code(result, 302021)

class TestSubmitAgain:
    @pytest.mark.P2
    @allure.title('Submit parammeter  proposal function verification')
    def test_PP_TI_001_002(self, new_genesis_env, clients_consensus):
        genesis = from_dict(data_class=Genesis, data=new_genesis_env.genesis_config)
        genesis.economicModel.gov.paramProposalVoteDurationSeconds = 0
        new_genesis_env.set_genesis(genesis.to_dict())
        new_genesis_env.deploy_all()
        submitcppandvote(clients_consensus[:3], [1, 1, 1])
        pip = clients_consensus[0].pip
        proposalinfo_param = pip.get_effect_proposal_info_of_vote(pip.cfg.param_proposal)
        log.info('Param proposal information : {}'.format(proposalinfo_param))
        proposalinfo_cancel = pip.get_effect_proposal_info_of_vote(pip.cfg.cancel_proposal)
        log.info('Cancel proposal information : {}'.format(proposalinfo_cancel))
        wait_block_number(pip.node, proposalinfo_cancel.get('EndVotingBlock'))
        assert_code(pip.get_status_of_proposal(proposalinfo_cancel.get('ProposalID')), 2)
        assert_code(pip.get_status_of_proposal(proposalinfo_param.get('ProposalID')), 6)

        result = pip.submitParam(pip.node.node_id, str(time.time()), 'slashing', 'slashBlocksReward', '998',
                                     pip.node.staking_address, transaction_cfg=pip.cfg.transaction_cfg)
        log.info('Submit param proposal result : {}'.format(result))
        assert_code(result, 0)
        proposalinfo_param = pip.get_effect_proposal_info_of_vote(pip.cfg.param_proposal)
        log.info('Param proposal information : {}'.format(proposalinfo_param))
        wait_block_number(pip.node, proposalinfo_param.get('EndVotingBlock'))

        result = pip.submitParam(pip.node.node_id, str(time.time()), 'slashing', 'slashBlocksReward', '998',
                                     pip.node.staking_address, transaction_cfg=pip.cfg.transaction_cfg)
        log.info('Submit param proposal result : {}'.format(result))
        assert_code(result, 0)

    @pytest.mark.P2
    @allure.title('Submit parammeter  proposal function verification')
    def test_VP_TI_001_002(self, no_vp_proposal, clients_verifier):
        submitcvpandvote(clients_verifier[:3], 1, 1, 1)
        pip = clients_verifier[0].pip
        proposalinfo_version = pip.get_effect_proposal_info_of_vote(pip.cfg.version_proposal)
        log.info('Version proposal information : {}'.format(proposalinfo_version))
        proposalinfo_cancel = pip.get_effect_proposal_info_of_vote(pip.cfg.cancel_proposal)
        log.info('Cancel proposal information : {}'.format(proposalinfo_cancel))
        wait_block_number(pip.node, proposalinfo_cancel.get('EndVotingBlock'))
        assert_code(pip.get_status_of_proposal(proposalinfo_cancel.get('ProposalID')), 2)
        assert_code(pip.get_status_of_proposal(proposalinfo_version.get('ProposalID')), 6)

        result = pip.submitVersion(pip.node.node_id, str(time.time()), pip.cfg.version5, 1,
                                       pip.node.staking_address, transaction_cfg=pip.cfg.transaction_cfg)
        log.info('Submit version proposal result : {}'.format(result))
        assert_code(result, 0)
        proposalinfo_version = pip.get_effect_proposal_info_of_vote(pip.cfg.version_proposal)
        log.info('Version proposal information : {}'.format(proposalinfo_version))
        wait_block_number(pip.node, proposalinfo_version.get('EndVotingBlock'))

        result = pip.submitVersion(pip.node.node_id, str(time.time()), pip.cfg.version5, 1,
                                       pip.node.staking_address, transaction_cfg=pip.cfg.transaction_cfg)
        log.info('Submit version proposal result : {}'.format(result))
        assert_code(result, 0)


class TestPIPVerify:
    @pytest.mark.P0
    @allure.title('Submit  proposal function verification---PIPID')
    def test_VP_PIP_001_003_TP_PI_001_003_CP_PI_001_003_CP_PI_001_003(self, no_vp_proposal):
        pip = no_vp_proposal
        pip_id_text = str(time.time())
        result = pip.submitText(pip.node.node_id, pip_id_text, pip.node.staking_address,
                                    transaction_cfg=pip.cfg.transaction_cfg)
        log.info('Submit text proposal result : {}'.format(result))
        assert_code(result, 0)

        result = pip.submitVersion(pip.node.node_id, pip_id_text, pip.cfg.version5, 1, pip.node.staking_address,
                                       transaction_cfg=pip.cfg.transaction_cfg)
        log.info('Same PIPID, submit version proposal result : {}'.format(result))
        assert_code(result, 302008)

        result = pip.submitText(pip.node.node_id, pip_id_text, pip.node.staking_address,
                                    transaction_cfg=pip.cfg.transaction_cfg)
        log.info('Same PIPID, submit text proposal result : {}'.format(result))
        assert_code(result, 302008)

        result = pip.submitParam(pip.node.node_id, pip_id_text, 'slashing', 'slashBlocksReward', '889',
                                     pip.node.staking_address, transaction_cfg=pip.cfg.transaction_cfg)
        log.info('Same PIPID, submit param proposal result : {}'.format(result))
        assert_code(result, 302008)

        result = pip.submitVersion(pip.node.node_id, str(time.time()), pip.cfg.version5, 3,
                                       pip.node.staking_address, transaction_cfg=pip.cfg.transaction_cfg)
        log.info('Differ PIPID, submit version proposal result : {}'.format(result))
        assert_code(result, 0)

        proposalinfo_version = pip.get_effect_proposal_info_of_vote(pip.cfg.version_proposal)
        log.info('Version proposal information : {}'.format(proposalinfo_version))
        pip_id_version = proposalinfo_version.get('PIPID')

        result = pip.submitCancel(pip.node.node_id, pip_id_text, 1, proposalinfo_version.get('ProposalID'),
                                      pip.node.staking_address, transaction_cfg=pip.cfg.transaction_cfg)
        log.info('Same PIPID, submit cancel proposal result : {}'.format(result))
        assert_code(result, 302008)

        result = pip.submitCancel(pip.node.node_id, pip_id_version, 1, proposalinfo_version.get('ProposalID'),
                                      pip.node.staking_address, transaction_cfg=pip.cfg.transaction_cfg)
        log.info('Same PIPID, submit cancel proposal result : {}'.format(result))
        assert_code(result, 302008)

        result = pip.submitCancel(pip.node.node_id, str(time.time()), 1, proposalinfo_version.get('ProposalID'),
                                      pip.node.staking_address, transaction_cfg=pip.cfg.transaction_cfg)
        log.info('Differ PIPID, submit cancel proposal result : {}'.format(result))
        assert_code(result, 0)

        proposalinfo_cancel = pip.get_effect_proposal_info_of_vote(pip.cfg.cancel_proposal)
        log.info('Version proposal information : {}'.format(proposalinfo_cancel))
        pip_id_cancel = proposalinfo_cancel.get('PIPID')

        result = pip.submitText(pip.node.node_id, pip_id_version, pip.node.staking_address,
                                    transaction_cfg=pip.cfg.transaction_cfg)
        log.info('Same PIPID, submit text proposal result : {}'.format(result))
        assert_code(result, 302008)

        result = pip.submitText(pip.node.node_id, pip_id_cancel, pip.node.staking_address,
                                    transaction_cfg=pip.cfg.transaction_cfg)
        log.info('Same PIPID, submit text proposal result : {}'.format(result))
        assert_code(result, 302008)

        wait_block_number(pip.node, proposalinfo_version.get('EndVotingBlock'))
        result = pip.submitVersion(pip.node.node_id, pip_id_version, pip.cfg.version5, 1,
                                       pip.node.staking_address, transaction_cfg=pip.cfg.transaction_cfg)
        log.info('Same PIPID, submit version proposal result : {}'.format(result))
        assert_code(result, 0)

        proposalinfo_version = pip.get_effect_proposal_info_of_vote(pip.cfg.version_proposal)
        log.info('Version proposal information : {}'.format(proposalinfo_version))
        wait_block_number(pip.node, proposalinfo_version.get('EndVotingBlock'))
        result = pip.submitParam(pip.node.node_id, pip_id_version, 'slashing', 'slashBlocksReward', '889',
                                     pip.node.staking_address, transaction_cfg=pip.cfg.transaction_cfg)
        log.info('Differ PIPID, submit param proposal result : {}'.format(result))
        assert_code(result, 0)

        proposalinfo_param = pip.get_effect_proposal_info_of_vote(pip.cfg.param_proposal)
        log.info('Param proposal information : {}'.format(proposalinfo_param))

        result = pip.submitCancel(pip.node.node_id, pip_id_cancel, 1, proposalinfo_param.get('ProposalID'),
                                      pip.node.staking_address, transaction_cfg=pip.cfg.transaction_cfg)
        log.info('Same PIPID, submit cancel proposal result : {}'.format(result))
        assert_code(result, 0)

        proposalinfo_cancel = pip.get_effect_proposal_info_of_vote(pip.cfg.cancel_proposal)
        log.info('Cancel proposal information : {}'.format(proposalinfo_cancel))

        wait_block_number(pip.node, proposalinfo_cancel.get('EndVotingBlock'))
        result = pip.submitText(pip.node.node_id, pip_id_cancel, pip.node.staking_address,
                                    transaction_cfg=pip.cfg.transaction_cfg)
        log.info('Same PIPID, submit text proposal result : {}'.format(result))
        assert_code(result, 0)

    @pytest.mark.P0
    @allure.title('Submit  proposal function verification---PIPID')
    def test_VP_PIP_002_TP_PI_002_CP_PI_002_CP_PI_002(self, no_vp_proposal, clients_verifier):
        pip = clients_verifier[0].pip
        submitcvpandvote(clients_verifier, 1, 1, 1, 1)
        proposalinfo = pip.get_effect_proposal_info_of_vote(pip.cfg.cancel_proposal)
        pip_id = proposalinfo.get('PIPID')
        wait_block_number(pip.node, proposalinfo.get('EndVotingBlock'))
        assert_code(pip.get_status_of_proposal(proposalinfo.get('ProposalID')), 2)
        result = pip.submitText(pip.node.node_id, pip_id, pip.node.staking_address,
                                    transaction_cfg=pip.cfg.transaction_cfg)
        log.info('Same pipid, submit text proposal result : {}'.format(result))
        assert_code(result, 302008)
        result = pip.submitVersion(pip.node.node_id, pip_id, pip.cfg.version5, 1, pip.node.staking_address,
                                       transaction_cfg=pip.cfg.transaction_cfg)
        log.info('Same pipid, submit text proposal result : {}'.format(result))
        assert_code(result, 302008)

        result = pip.submitVersion(pip.node.node_id, str(time.time()), pip.cfg.version5, 2,
                                       pip.node.staking_address,
                                       transaction_cfg=pip.cfg.transaction_cfg)
        log.info('Differ pipid, submit version proposal result : {}'.format(result))
        assert_code(result, 0)
        proposalinfo_version = pip.get_effect_proposal_info_of_vote()
        log.info('Proposal information : {}'.format(proposalinfo_version))
        result = pip.submitCancel(pip.node.node_id, pip_id, 1, proposalinfo_version.get('ProposalID'),
                                      pip.node.staking_address, transaction_cfg=pip.cfg.transaction_cfg)
        log.info('Same pipid, submit cancel proposal result : {}'.format(result))
        assert_code(result, 302008)

        wait_block_number(pip.node, proposalinfo_version.get('EndVotingBlock'))
        result = pip.submitParam(pip.node.node_id, pip_id, 'slashing', 'slashBlocksReward', '19',
                                     pip.node.staking_address, transaction_cfg=pip.cfg.transaction_cfg)
        log.info('Same pipid, submit param proposal result : {}'.format(result))
        assert_code(result, 302008)

        result = pip.submitParam(pip.node.node_id, str(time.time()), 'slashing', 'slashBlocksReward', '19',
                                     pip.node.staking_address, transaction_cfg=pip.cfg.transaction_cfg)
        log.info('Differ pipid, submit param proposal result : {}'.format(result))
        assert_code(result, 0)

        proposalinfo_param = pip.get_effect_proposal_info_of_vote(pip.cfg.param_proposal)
        log.info('Get param proposal information : {}'.format(proposalinfo_param))

        result = pip.submitCancel(pip.node.node_id, pip_id, 1, proposalinfo_param.get('ProposalID'),
                                      pip.node.staking_address, transaction_cfg=pip.cfg.transaction_cfg)
        log.info('Same pipid, submit cancel proposal result : {}'.format(result))
        assert_code(result, 302008)


@pytest.mark.P0
@allure.title('Submit cancel  proposal function verification---endvoting rounds')
def test_CP_CR_003_CP_CR_004(submit_param):
    pip = submit_param
    proposalinfo = pip.get_effect_proposal_info_of_vote(pip.cfg.param_proposal)
    log.info('proposalinfo: {}'.format(proposalinfo))
    endvoting_rounds = (proposalinfo.get('EndVotingBlock') - math.ceil(
        pip.node.block_number / pip.economic.consensus_size) * pip.economic.consensus_size
        ) / pip.economic.consensus_size
    result = pip.submitCancel(pip.node.node_id, str(time.time()), endvoting_rounds + 1,
                                  proposalinfo.get('ProposalID'), pip.node.staking_address,
                                  transaction_cfg=pip.cfg.transaction_cfg)
    log.info('endvoting_rounds:{}， cancel proposal result:{}'.format(endvoting_rounds, result))
    assert_code(result, 302010)

    result = pip.submitCancel(pip.node.node_id, str(time.time()), endvoting_rounds,
                                  proposalinfo.get('ProposalID'), pip.node.staking_address,
                                  transaction_cfg=pip.cfg.transaction_cfg)
    log.info('endvoting_rounds:{}， cancel proposal result:{}'.format(endvoting_rounds + 1, result))
    assert_code(result, 0)


class TestGas:
    @pytest.mark.P2
    @allure.title('Submit version proposal function verification---gasprice')
    def test_VP_GP_001_VP_GP_002(self, no_vp_proposal):
        pip = no_vp_proposal
        transaction_cfg = {"gasPrice": 2100000000000000 - 1}
        try:
            pip.submitVersion(pip.node.node_id, str(time.time()), pip.cfg.version5, 1,
                                  pip.node.staking_address, transaction_cfg=transaction_cfg)
        except ValueError as e:
            assert e.args[0].get('message') == "the tx data is invalid: Invalid parameter:gas price is lower than minimum"

        transaction_cfg = {"gasPrice": 2100000000000000}
        result = pip.submitVersion(pip.node.node_id, str(time.time()), pip.cfg.version5, 1,
                                       pip.node.staking_address, transaction_cfg=transaction_cfg)
        log.info('Submit version proposal result : {}'.format(result))
        assert_code(result, 0)

    @pytest.mark.P2
    @allure.title('Submit param proposal function verification---gasprice')
    def test_PP_GP_001_PP_GP_002(self, no_vp_proposal):
        pip = no_vp_proposal
        transaction_cfg = {"gasPrice": 2000000000000000 - 1}
        try:
            pip.submitParam(pip.node.node_id, str(time.time()), 'slashing', 'slashBlocksReward', '123',
                                pip.node.staking_address, transaction_cfg=transaction_cfg)
        except ValueError as e:
            assert e.args[0].get('message') == "the tx data is invalid: Invalid parameter:gas price is lower than minimum"

        transaction_cfg = {"gasPrice": 2000000000000000}
        result = pip.submitParam(pip.node.node_id, str(time.time()), 'slashing', 'slashBlocksReward', '123',
                                     pip.node.staking_address, transaction_cfg=transaction_cfg)
        log.info('Submit param proposal result : {}'.format(result))
        assert_code(result, 0)

    @pytest.mark.P2
    @allure.title('Submit text proposal function verification---gasprice')
    def test_TP_GP_001_TP_GP_002(self, client_verifier):
        pip = client_verifier.pip
        transaction_cfg = {"gasPrice": 1500000000000000 - 1}
        try:
            pip.submitText(pip.node.node_id, str(time.time()), pip.node.staking_address,
                               transaction_cfg=transaction_cfg)
        except ValueError as e:
            assert e.args[0].get('message') == "the tx data is invalid: Invalid parameter:gas price is lower than minimum"

        transaction_cfg = {"gasPrice": 1500000000000000}
        result = pip.submitText(pip.node.node_id, str(time.time()), pip.node.staking_address,
                                    transaction_cfg=transaction_cfg)
        log.info('Submit text proposal result : {}'.format(result))
        assert_code(result, 0)

    @pytest.mark.P2
    @allure.title('Submit cancel proposal function verification---gas')
    def test_CP_GP_001_CP_GP_002(self, no_vp_proposal):
        pip = no_vp_proposal
        transaction_cfg = {"gasPrice": 3000000000000000 - 1, "gas": 100000}
        try:
            result = pip.submitVersion(pip.node.node_id, str(time.time()), pip.cfg.version5, 3,
                                           pip.node.staking_address, transaction_cfg=pip.cfg.transaction_cfg)
            assert_code(result, 0)
            proposalinfo = pip.get_effect_proposal_info_of_vote()
            log.info('Get proposal information {}'.format(proposalinfo))
            pip.submitCancel(pip.node.node_id, str(time.time()), 1, proposalinfo.get('ProposalID'),
                                 pip.node.staking_address, transaction_cfg=transaction_cfg)
        except ValueError as e:
            assert e.args[0].get('message') == "the tx data is invalid: Invalid parameter:gas price is lower than minimum"
        transaction_cfg = {"gasPrice": 3000000000000000}
        result = pip.submitCancel(pip.node.node_id, str(time.time()), 1, proposalinfo.get('ProposalID'),
                                      pip.node.staking_address, transaction_cfg=transaction_cfg)
        log.info('Submit cancel proposal result : {}'.format(result))
        assert_code(result, 0)


def TP_TE_004(new_genesis_env, client_consensus):
    genesis = from_dict(data_class=Genesis, data=new_genesis_env.genesis_config)
    genesis.economicModel.gov.textProposalVoteDurationSeconds = 0
    new_genesis_env.set_genesis(genesis.to_dict())
    new_genesis_env.deploy_all()
    client = client_consensus
    consensus_size = client.economic.consensus_size
    log.info(consensus_size)
    while True:
        if client.node.block_number % consensus_size > consensus_size - 20:
            log.info(client.node.block_number)
            result = client.pip.submitText(client.node.node_id, str(time.time()), client.node.staking_address,
                                           transaction_cfg=client.pip.cfg.transaction_cfg)
            log.info('Submit text proposal result : {}'.format(result))
            log.info(client.pip.pip.listProposal())
            assert_code(result, 1)
            break

if __name__ == '__main__':
    pytest.main(['./tests/govern/', '-s', '-q', '--alluredir', './report/report'])
