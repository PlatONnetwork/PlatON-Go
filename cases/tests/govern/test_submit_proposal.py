from common.log import log
from dacite import from_dict
from tests.lib import Genesis
import pytest
from tests.lib.utils import wait_block_number, assert_code, get_governable_parameter_value
from tests.lib.client import get_client_obj
import time
import math
from tests.govern.test_voting_statistics import submitcppandvote, submitcvpandvote, submittpandvote


@pytest.mark.P0
@pytest.mark.compatibility
def test_VP_SU_001(submit_version):
    pip_obj = submit_version
    proposalinfo = pip_obj.get_effect_proposal_info_of_vote()
    log.info('Get version proposal information : {}'.format(proposalinfo))
    endvotingblock_count = math.ceil(proposalinfo.get('SubmitBlock') / pip_obj.economic.consensus_size +
                                     proposalinfo.get('EndVotingRounds')
                                     ) * pip_obj.economic.consensus_size - 20
    log.info('Calculated endvoting block{},interface returned endvoting block{}'.format(endvotingblock_count,
                                                                                        proposalinfo.get('EndVotingBlock')))
    assert int(endvotingblock_count) == proposalinfo.get('EndVotingBlock')
    assert int(endvotingblock_count) + 21 == proposalinfo.get('ActiveBlock')


@pytest.mark.P0
@pytest.mark.compatibility
def test_CP_SU_001_CP_UN_001(submit_cancel):
    pip_obj = submit_cancel
    proposalinfo = pip_obj.get_effect_proposal_info_of_vote(pip_obj.cfg.cancel_proposal)
    log.info('cancel proposalinfo : {}'.format(proposalinfo))
    endvotingblock_count = math.ceil(proposalinfo.get('SubmitBlock') / pip_obj.economic.consensus_size +
                                     proposalinfo.get('EndVotingRounds')
                                     ) * pip_obj.economic.consensus_size - 20
    log.info('Calculated endvoting block{},interface returned endvoting block{}'.format(endvotingblock_count,
                                                                                        proposalinfo.get('EndVotingBlock')))
    assert int(endvotingblock_count) == proposalinfo.get('EndVotingBlock')
    result = pip_obj.submitCancel(pip_obj.node.node_id, str(time.time()), 1, proposalinfo.get('ProposalID'),
                                  pip_obj.node.staking_address, transaction_cfg=pip_obj.cfg.transaction_cfg)
    log.info('Submit cancel proposal result : {}'.format(result))
    assert_code(result, 302014)


class TestsubmitCP():
    @pytest.mark.P0
    def test_CP_SU_002_CP_SU_003(self, submit_param):
        pip_obj = submit_param
        proposalinfo = pip_obj.get_effect_proposal_info_of_vote(pip_obj.cfg.param_proposal)
        log.info('param proposalinfo : {}'.format(proposalinfo))
        endvotingrounds_count = (proposalinfo.get('EndVotingBlock') -
                                 math.ceil(pip_obj.node.block_number/pip_obj.economic.consensus_size) *
                                 pip_obj.economic.consensus_size) / pip_obj.economic.consensus_size
        log.info('caculated endvoting rounds is {}'.format(endvotingrounds_count))
        result = pip_obj.submitCancel(pip_obj.node.node_id, str(time.time()), endvotingrounds_count + 1, proposalinfo.get('ProposalID'),
                                      pip_obj.node.staking_address, transaction_cfg=pip_obj.cfg.transaction_cfg)
        log.info('Submit cancel proposal result : {}'.format(result))
        assert_code(result, 302010)
        result = pip_obj.submitCancel(pip_obj.node.node_id, str(time.time()), endvotingrounds_count, proposalinfo.get('ProposalID'),
                                      pip_obj.node.staking_address, transaction_cfg=pip_obj.cfg.transaction_cfg)
        log.info('Submit cancel proposal result : {}'.format(result))
        assert_code(result, 0)

    @pytest.mark.P0
    def test_CP_SU_002_CP_UN_002(self, submit_cancel_param):
        pip_obj = submit_cancel_param
        proposalinfo = pip_obj.get_effect_proposal_info_of_vote(pip_obj.cfg.cancel_proposal)
        log.info('cancel proposalinfo : {}'.format(proposalinfo))
        endvotingblock_count = math.ceil(proposalinfo.get('SubmitBlock') / pip_obj.economic.consensus_size +
                                         proposalinfo.get('EndVotingRounds')
                                         ) * pip_obj.economic.consensus_size - 20
        log.info('Calculated endvoting block{},interface returned endvoting block{}'.format(endvotingblock_count,
                                                                                            proposalinfo.get('EndVotingBlock')))
        assert int(endvotingblock_count) == proposalinfo.get('EndVotingBlock')
        result = pip_obj.submitCancel(pip_obj.node.node_id, str(time.time()), 1, proposalinfo.get('ProposalID'),
                                      pip_obj.node.staking_address, transaction_cfg=pip_obj.cfg.transaction_cfg)
        log.info('Submit cancel proposal result : {}'.format(result))
        assert_code(result, 302014)


@pytest.mark.P0
@pytest.mark.compatibility
def test_PP_SU_001_PP_UN_001_VP_UN_003(submit_param):
    pip_obj = submit_param
    proposalinfo = pip_obj.get_effect_proposal_info_of_vote(pip_obj.cfg.param_proposal)
    log.info('param proposalinfo : {}'.format(proposalinfo))
    endvotingblock_count = math.ceil(proposalinfo.get('SubmitBlock') / pip_obj.economic.settlement_size +
                                     pip_obj.economic.pp_vote_settlement_wheel
                                     ) * pip_obj.economic.settlement_size
    log.info('Calculated endvoting block {},interface returned endvoting block {}'.format(endvotingblock_count,
                                                                                          proposalinfo.get('EndVotingBlock')))
    assert int(endvotingblock_count) == proposalinfo.get('EndVotingBlock')
    result = pip_obj.submitParam(pip_obj.node.node_id, str(time.time()), 'slashing', 'slashBlocksReward', '22',
                                 pip_obj.node.staking_address, transaction_cfg=pip_obj.cfg.transaction_cfg)
    log.info('There is a voting param proposal,submit param proposal result : {}'.format(result))
    assert_code(result, 302032)

    result = pip_obj.submitVersion(pip_obj.node.node_id, str(time.time()), pip_obj.cfg.version5, 1,
                                   pip_obj.node.staking_address, transaction_cfg=pip_obj.cfg.transaction_cfg)
    log.info('There is a voting param proposal,submit version proposal result : {}'.format(result))
    assert_code(result, 302032)


@pytest.mark.P0
def test_VP_VE_001_to_VP_VE_005(no_vp_proposal):
    pip_obj_tmp = no_vp_proposal
    result = pip_obj_tmp.submitVersion(pip_obj_tmp.node.node_id, str(time.time()), pip_obj_tmp.cfg.version1, 1,
                                       pip_obj_tmp.node.staking_address, transaction_cfg=pip_obj_tmp.cfg.transaction_cfg)
    log.info('Submit version proposal result : {}'.format(result))
    assert_code(result, 302011)

    result = pip_obj_tmp.submitVersion(pip_obj_tmp.node.node_id, str(time.time()), pip_obj_tmp.cfg.version2, 1,
                                       pip_obj_tmp.node.staking_address, transaction_cfg=pip_obj_tmp.cfg.transaction_cfg)
    log.info('Submit version proposal result : {}'.format(result))
    assert_code(result, 302011)

    result = pip_obj_tmp.submitVersion(pip_obj_tmp.node.node_id, str(time.time()), pip_obj_tmp.cfg.version3, 1,
                                       pip_obj_tmp.node.staking_address, transaction_cfg=pip_obj_tmp.cfg.transaction_cfg)
    log.info('Submit version proposal result : {}'.format(result))
    assert_code(result, 302011)

    result = pip_obj_tmp.submitVersion(pip_obj_tmp.node.node_id, str(time.time()), pip_obj_tmp.chain_version, 1,
                                       pip_obj_tmp.node.staking_address, transaction_cfg=pip_obj_tmp.cfg.transaction_cfg)
    log.info('Submit version proposal result : {}'.format(result))
    assert_code(result, 302011)

    result = pip_obj_tmp.submitVersion(pip_obj_tmp.node.node_id, str(time.time()), pip_obj_tmp.cfg.version8, 1,
                                       pip_obj_tmp.node.staking_address, transaction_cfg=pip_obj_tmp.cfg.transaction_cfg)
    log.info('Submit version proposal result : {}'.format(result))
    assert_code(result, 0)


@pytest.mark.P2
def test_VP_WA_001(no_vp_proposal):
    pip_obj_tmp = no_vp_proposal
    address, _ = pip_obj_tmp.economic.account.generate_account(pip_obj_tmp.node.web3, 10**18 * 10000000)
    result = pip_obj_tmp.submitVersion(pip_obj_tmp.node.node_id, str(time.time()), pip_obj_tmp.cfg.version5, 1,
                                       address, transaction_cfg=pip_obj_tmp.cfg.transaction_cfg)
    log.info('Submit version proposal reuslt : {}'.format(result))
    assert_code(result, 302021)


@pytest.mark.P2
def test_TP_WA_001(client_verifier_obj):
    pip_obj = client_verifier_obj.pip
    address, _ = pip_obj.economic.account.generate_account(pip_obj.node.web3, 10**18 * 10000000)
    result = pip_obj.submitVersion(pip_obj.node.node_id, str(time.time()), pip_obj.cfg.version5, 1,
                                   address, transaction_cfg=pip_obj.cfg.transaction_cfg)
    log.info('Submit version proposal reuslt : {}'.format(result))
    assert_code(result, 302021)


@pytest.mark.P0
def test_TP_UN_001(submit_text):
    pip_obj = submit_text
    result = pip_obj.submitText(pip_obj.node.node_id, str(time.time()), pip_obj.node.staking_address,
                                transaction_cfg=pip_obj.cfg.transaction_cfg)
    log.info('There is voting text proposal, submit text proposal result : {}'.format(result))
    assert_code(result, 0)


@pytest.mark.P0
@pytest.mark.compatibility
def test_VP_SU_001_VP_UN_001(submit_version):
    pip_obj = submit_version
    result = pip_obj.submitVersion(pip_obj.node.node_id, str(time.time()), pip_obj.cfg.version5, 1,
                                   pip_obj.node.staking_address, transaction_cfg=pip_obj.cfg.transaction_cfg)
    log.info('There is voting version proposal, submit version proposal result : {}'.format(result))
    assert_code(result, 302012)


@pytest.mark.P0
def test_VP_UN_002_CP_ID_002(preactive_proposal_pipobj_list, new_genesis_env):
    pip_obj = preactive_proposal_pipobj_list[0]
    proposalinfo = pip_obj.get_effect_proposal_info_of_preactive()
    log.info('Get preactive proposal info: {}'.format(proposalinfo))

    result = pip_obj.submitVersion(pip_obj.node.node_id, str(time.time()), pip_obj.cfg.version5, 1,
                                   pip_obj.node.staking_address,
                                   transaction_cfg=pip_obj.cfg.transaction_cfg)
    log.info('There is preactive version proposal, submit version proposal result : {}'.format(result))
    assert_code(result, 302013)

    result = pip_obj.submitCancel(pip_obj.node.node_id, str(time.time()), 1, proposalinfo.get('ProposalID'),
                                  pip_obj.node.staking_address, transaction_cfg=pip_obj.cfg.transaction_cfg)
    log.info('there is preactive version proposal, submit cancel proposal result: {}'.format(result))
    assert_code(result, 302017)


@pytest.mark.P0
def test_PP_UN_003(preactive_proposal_pipobj_list, new_genesis_env):
    pip_obj = preactive_proposal_pipobj_list[0]
    proposalinfo = pip_obj.get_effect_proposal_info_of_preactive()
    log.info('Get preactive proposal info: {}'.format(proposalinfo))

    result = pip_obj.submitParam(pip_obj.node.node_id, str(time.time()), 'slashing', 'slashBlocksReward',
                                 '84', pip_obj.node.staking_address, transaction_cfg=pip_obj.cfg.transaction_cfg)
    log.info('there is preactive version proposal, submit cancel param proposal result: {}'.format(result))
    assert_code(result, 302013)


class TestEndVotingRounds():
    @pytest.mark.P1
    def test_VP_CR_001_VP_CR_002_VP_CR_007_TP_TE_002(self, new_genesis_env, client_verifier_obj):
        '''
        Proposal vote duration set consensus size's accompanying number +1
        :param pip_env:
        :param pip_obj:
        :return:
        '''
        pip_obj = client_verifier_obj.pip
        genesis = from_dict(data_class=Genesis, data=new_genesis_env.genesis_config)
        genesis.economicModel.gov.versionProposalVoteDurationSeconds = 2 * pip_obj.economic.consensus_size + 1
        genesis.economicModel.gov.textProposalVoteDurationSeconds = 5 * pip_obj.economic.consensus_size + 1
        new_genesis_env.set_genesis(genesis.to_dict())
        new_genesis_env.deploy_all()
        result = pip_obj.submitVersion(pip_obj.node.node_id, str(time.time()), pip_obj.cfg.version5, 3,
                                       pip_obj.node.staking_address,
                                       transaction_cfg=pip_obj.cfg.transaction_cfg)
        log.info('endvoting rounds is three, subtmit version proposal result : {}'.format(result))
        assert_code(result, 302010)

        result = pip_obj.submitVersion(pip_obj.node.node_id, str(time.time()), pip_obj.cfg.version5, 0,
                                       pip_obj.node.staking_address,
                                       transaction_cfg=pip_obj.cfg.transaction_cfg)
        log.info('endvoting rounds is zero, subtmit version proposal result : {}'.format(result))
        assert_code(result, 302009)

        result = pip_obj.submitVersion(pip_obj.node.node_id, str(time.time()), pip_obj.cfg.version5, 2,
                                       pip_obj.node.staking_address,
                                       transaction_cfg=pip_obj.cfg.transaction_cfg)
        log.info('endvoting rounds is two, subtmit version proposal result : {}'.format(result))
        assert_code(result, 0)

        result = pip_obj.submitText(pip_obj.node.node_id, str(time.time()), pip_obj.node.staking_address,
                                    transaction_cfg=pip_obj.cfg.transaction_cfg)
        log.info('Submit text proposal result : {}'.format(result))
        assert_code(result, 0)
        proposalinfo = pip_obj.get_effect_proposal_info_of_vote(pip_obj.cfg.text_proposal)
        log.info('Get text proposal information :{}'.format(proposalinfo))
        endvotingblock_count = math.ceil(proposalinfo.get('SubmitBlock') / pip_obj.economic.consensus_size + 5
                                         ) * pip_obj.economic.consensus_size - 20
        log.info('calcuted endvoting block {},interface returns {}'.format(endvotingblock_count,
                                                                           proposalinfo.get('EndVotingBlock')))
        assert int(endvotingblock_count) == proposalinfo.get('EndVotingBlock')

    @pytest.mark.P1
    def test_VP_CR_003_VP_CR_004_VP_CR_007_TP_TE_003(self, new_genesis_env, client_verifier_obj):
        pip_obj = client_verifier_obj.pip
        genesis = from_dict(data_class=Genesis, data=new_genesis_env.genesis_config)
        genesis.economicModel.gov.versionProposalVoteDurationSeconds = 3 * pip_obj.economic.consensus_size - 1
        genesis.economicModel.gov.textProposalVoteDurationSeconds = 5 * pip_obj.economic.consensus_size - 1
        new_genesis_env.set_genesis(genesis.to_dict())
        new_genesis_env.deploy_all()
        result = pip_obj.submitVersion(pip_obj.node.node_id, str(time.time()), pip_obj.cfg.version5, 3, pip_obj.node.staking_address,
                                       transaction_cfg=pip_obj.cfg.transaction_cfg)
        log.info('endvoting rounds is three, subtmit version proposal result : {}'.format(result))
        assert_code(result, 302010)

        result = pip_obj.submitVersion(pip_obj.node.node_id, str(time.time()), pip_obj.cfg.version5, 0, pip_obj.node.staking_address,
                                       transaction_cfg=pip_obj.cfg.transaction_cfg)
        log.info('Endvoting rounds is three, subtmit version proposal result : {}'.format(result))
        assert_code(result, 302009)

        result = pip_obj.submitVersion(pip_obj.node.node_id, str(time.time()), pip_obj.cfg.version5, 2, pip_obj.node.staking_address,
                                       transaction_cfg=pip_obj.cfg.transaction_cfg)
        log.info('Endvoting rounds is two, subtmit version proposal result : {}'.format(result))
        assert_code(result, 0)

        result = pip_obj.submitText(pip_obj.node.node_id, str(time.time()), pip_obj.node.staking_address,
                                    transaction_cfg=pip_obj.cfg.transaction_cfg)
        log.info('Submit text proposal result : {}'.format(result))
        assert_code(result, 0)
        proposalinfo = pip_obj.get_effect_proposal_info_of_vote(pip_obj.cfg.text_proposal)
        log.info('Get text proposal information : {}'.format(proposalinfo))
        endvotingblock_count = math.ceil(proposalinfo.get('SubmitBlock') / pip_obj.economic.consensus_size + 4
                                         ) * pip_obj.economic.consensus_size - 20
        log.info('Calcuted endvoting block {},interface return {}'.format(endvotingblock_count,
                                                                          proposalinfo.get('EndVotingBlock')))
        assert int(endvotingblock_count) == proposalinfo.get('EndVotingBlock')

    @pytest.mark.P1
    @pytest.mark.compatibility
    def test_VP_CR_005_VP_CR_006_TP_TE_001(self, new_genesis_env, client_verifier_obj):
        pip_obj = client_verifier_obj.pip
        genesis = from_dict(data_class=Genesis, data=new_genesis_env.genesis_config)
        genesis.economicModel.gov.versionProposalVoteDurationSeconds = 3 * pip_obj.economic.consensus_size
        genesis.economicModel.gov.textProposalVoteDurationSeconds = 5 * pip_obj.economic.consensus_size
        new_genesis_env.set_genesis(genesis.to_dict())
        new_genesis_env.deploy_all()
        result = pip_obj.submitVersion(pip_obj.node.node_id, str(time.time()), pip_obj.cfg.version5, 4,
                                       pip_obj.node.staking_address,
                                       transaction_cfg=pip_obj.cfg.transaction_cfg)
        log.info('Endvoting rounds is four, subtmit version proposal result : {}'.format(result))
        assert_code(result, 302010)

        result = pip_obj.submitVersion(pip_obj.node.node_id, str(time.time()), pip_obj.cfg.version5, 0,
                                       pip_obj.node.staking_address,
                                       transaction_cfg=pip_obj.cfg.transaction_cfg)
        log.info('Endvoting rounds is zero, subtmit version proposal result : {}'.format(result))
        assert_code(result, 302009)

        result = pip_obj.submitVersion(pip_obj.node.node_id, str(time.time()), pip_obj.cfg.version5, 3,
                                       pip_obj.node.staking_address,
                                       transaction_cfg=pip_obj.cfg.transaction_cfg)
        log.info('Endvoting rounds is zero, subtmit version proposal result : {}'.format(result))
        assert_code(result, 0)

        result = pip_obj.submitText(pip_obj.node.node_id, str(time.time()), pip_obj.node.staking_address,
                                    transaction_cfg=pip_obj.cfg.transaction_cfg)
        log.info('Submit text proposal result : {}'.format(result))
        assert_code(result, 0)
        proposalinfo = pip_obj.get_effect_proposal_info_of_vote(pip_obj.cfg.text_proposal)
        log.info('Get text proposal information {}'.format(proposalinfo))
        endvotingblock_count = math.ceil(proposalinfo.get('SubmitBlock') / pip_obj.economic.consensus_size + 5
                                         ) * pip_obj.economic.consensus_size - 20
        log.info('calcuted endvoting block {},interface return {}'.format(endvotingblock_count,
                                                                          proposalinfo.get('EndVotingBlock')))
        assert int(endvotingblock_count) == proposalinfo.get('EndVotingBlock')

        proosalinfo = pip_obj.get_effect_proposal_info_of_vote(1)
        log.info('text proposalinfo: {}'.format(proosalinfo))
        result = pip_obj.submitCancel(pip_obj.node.node_id, str(time.time()), 1, proosalinfo.get('ProposalID'),
                                      pip_obj.node.staking_address, transaction_cfg=pip_obj.cfg.transaction_cfg)
        log.info('submit cancel result: {}'.format(result))


class TestNoVerifierSubmitProposal():
    @pytest.mark.P0
    def test_VP_PR_002_TP_PR_002(self, no_vp_proposal, client_new_node_obj):
        pip_obj = client_new_node_obj.pip
        address, _ = pip_obj.economic.account.generate_account(pip_obj.node.web3, 10**18 * 10000000)
        result = pip_obj.submitVersion(pip_obj.node.node_id, str(time.time()), pip_obj.cfg.version5, 1, address,
                                       transaction_cfg=pip_obj.cfg.transaction_cfg)
        log.info('New node submit version proposal : {}'.format(result))
        assert_code(result, 302022)

        result = pip_obj.submitText(pip_obj.node.node_id, str(time.time()), address,
                                    transaction_cfg=pip_obj.cfg.transaction_cfg)
        log.info('New node submit text proposal : {}'.format(result))
        assert_code(result, 302022)

    @pytest.mark.P0
    def test_VP_PR_001_TP_PR_001(self, client_candidate_obj):
        pip_obj = client_candidate_obj.pip
        result = pip_obj.submitVersion(pip_obj.node.node_id, str(time.time()),
                                       pip_obj.cfg.version5, 1, pip_obj.node.staking_address,
                                       transaction_cfg=pip_obj.cfg.transaction_cfg)
        log.info('Candidate node {} submit version proposal result : {}'.format(pip_obj.node.node_id, result))
        assert_code(result, 302022)

        result = pip_obj.submitText(pip_obj.node.node_id, str(time.time()), pip_obj.node.staking_address,
                                    transaction_cfg=pip_obj.cfg.transaction_cfg)
        log.info('candidate node {} submit text proposal result : {}'.format(pip_obj.node.node_id, result))
        assert_code(result, 302022)

    @pytest.mark.P2
    def test_VP_PR_003_VP_PR_004_TP_PR_003_TP_PR_004(self, client_verifier_obj):
        address = client_verifier_obj.node.staking_address
        result = client_verifier_obj.staking.withdrew_staking(address)
        log.info('Node {} withdrew staking result : {}'.format(client_verifier_obj.node.node_id, result))
        assert_code(result, 0)
        log.info(client_verifier_obj.economic.account.find_pri_key(address))
        result = client_verifier_obj.pip.submitVersion(client_verifier_obj.node.node_id, str(time.time()),
                                                       client_verifier_obj.pip.cfg.version5, 1, address,
                                                       transaction_cfg=client_verifier_obj.pip.cfg.transaction_cfg)
        log.info('Node exiting submit version proposal :{}'.format(result))
        assert_code(result, 302020)

        result = client_verifier_obj.pip.submitText(client_verifier_obj.node.node_id, str(time.time()), address,
                                                    transaction_cfg=client_verifier_obj.pip.cfg.transaction_cfg)
        log.info('Node exiting submit text proposal : {}'.format(result))
        assert_code(result, 302020)

        client_verifier_obj.economic.wait_settlement_blocknum(client_verifier_obj.node,
                                                              number=client_verifier_obj.economic.unstaking_freeze_ratio)
        result = client_verifier_obj.pip.submitVersion(client_verifier_obj.node.node_id, str(time.time()),
                                                       client_verifier_obj.pip.cfg.version5, 1, address,
                                                       transaction_cfg=client_verifier_obj.pip.cfg.transaction_cfg)
        log.info('Node exited submit version proposal : {}'.format(result))
        assert_code(result, 302022)

        client_verifier_obj.economic.wait_settlement_blocknum(client_verifier_obj.node,
                                                              number=client_verifier_obj.economic.unstaking_freeze_ratio)
        result = client_verifier_obj.pip.submitText(client_verifier_obj.node.node_id, str(time.time()), address,
                                                    transaction_cfg=client_verifier_obj.pip.cfg.transaction_cfg)
        log.info('Node exited submit text proposal : {}'.format(result))
        assert_code(result, 302022)


class TestSubmitCancel():
    @pytest.mark.P0
    def test_CP_WA_001(self, submit_version):
        pip_obj = submit_version
        address, _ = pip_obj.economic.account.generate_account(pip_obj.node.web3, 10**18 * 10000)
        proposalinfo = pip_obj.get_effect_proposal_info_of_vote()
        log.info('Get version proposal information : {}'.format(proposalinfo))
        result = pip_obj.submitCancel(pip_obj.node.node_id, str(time.time()), 1, proposalinfo.get('ProposalID'), address,
                                      transaction_cfg=pip_obj.cfg.transaction_cfg)
        log.info('Submit cancel proposal result : {}'.format(result))
        assert_code(result, 302021)

    @pytest.mark.P0
    def test_CP_PR_001(self, new_node_has_proposal):
        pip_obj = new_node_has_proposal
        address, _ = pip_obj.economic.account.generate_account(pip_obj.node.web3, 10**18 * 10000)
        proposalinfo = pip_obj.get_effect_proposal_info_of_vote()
        log.info('Get version proposal information : {}'.format(proposalinfo))
        result = pip_obj.submitCancel(pip_obj.node.node_id, str(time.time()), 1,
                                      proposalinfo.get('ProposalID'),
                                      address, transaction_cfg=pip_obj.cfg.transaction_cfg)
        log.info('Submit cancel proposal resullt : {}'.format(result))
        assert_code(result, 302022)

    @pytest.mark.P0
    def test_CP_PR_002(self, candidate_has_proposal):
        pip_obj = candidate_has_proposal
        proposalinfo = pip_obj.get_effect_proposal_info_of_vote()
        if not proposalinfo:
            time.sleep(10)
            proposalinfo = pip_obj.get_effect_proposal_info_of_vote()
        log.info('Get proposal information {}'.format(proposalinfo))
        result = pip_obj.submitCancel(pip_obj.node.node_id, str(time.time()), 1,
                                      proposalinfo.get('ProposalID'),
                                      pip_obj.node.staking_address, transaction_cfg=pip_obj.cfg.transaction_cfg)
        log.info('Candidate submit cancel proposal result : {}'.format(result))
        assert_code(result, 302022)

    @pytest.mark.P2
    def test_CP_PR_003_CP_PR_004(self, new_genesis_env, client_consensus_obj):
        genesis = from_dict(data_class=Genesis, data=new_genesis_env.genesis_config)
        genesis.economicModel.gov.versionProposalVoteDurationSeconds = 10000
        new_genesis_env.set_genesis(genesis.to_dict())
        new_genesis_env.deploy_all()
        pip_obj = client_consensus_obj.pip
        result = pip_obj.submitVersion(pip_obj.node.node_id, str(time.time()), pip_obj.cfg.version5, 20,
                                       pip_obj.node.staking_address, transaction_cfg=pip_obj.cfg.transaction_cfg)
        log.info('Submit version proposal result : {}'.format(result))
        assert_code(result, 0)
        proposalinfo = pip_obj.get_effect_proposal_info_of_vote()
        log.info('proposalinfo: {}'.format(proposalinfo))
        client_obj = client_consensus_obj
        address = pip_obj.node.staking_address
        result = client_obj.staking.withdrew_staking(address)
        log.info('nodeid: {} withdrewstaking result: {}'.format(client_obj.node.node_id, result))
        assert_code(result, 0)
        result = client_obj.pip.submitCancel(client_obj.node.node_id, str(time.time()), 1,
                                             proposalinfo.get('ProposalID'), address,
                                             transaction_cfg=client_obj.pip.cfg.transaction_cfg)
        log.info('node exiting，cancel proposal result: {}'.format(result))
        assert_code(result, 302020)

        client_obj.economic.wait_settlement_blocknum(client_obj.node,
                                                     number=client_obj.economic.unstaking_freeze_ratio)
        result = client_obj.pip.submitCancel(client_obj.node.node_id, str(time.time()), 1,
                                             proposalinfo.get('ProposalID'), address,
                                             transaction_cfg=client_obj.pip.cfg.transaction_cfg)
        log.info('exited node，cancel proposal result: {}'.format(result))
        assert_code(result, 302022)

    @pytest.mark.P0
    def test_CP_CR_001_CP_CR_002(self, submit_version):
        pip_obj = submit_version
        proposalinfo = pip_obj.get_effect_proposal_info_of_vote()
        log.info('proposalinfo: {}'.format(proposalinfo))
        endvoting_rounds = (proposalinfo.get('EndVotingBlock') + 20 - math.ceil(
            pip_obj.node.block_number / pip_obj.economic.consensus_size) * pip_obj.economic.consensus_size
            ) / pip_obj.economic.consensus_size
        result = pip_obj.submitCancel(pip_obj.node.node_id, str(time.time()), endvoting_rounds,
                                      proposalinfo.get('ProposalID'), pip_obj.node.staking_address,
                                      transaction_cfg=pip_obj.cfg.transaction_cfg)
        log.info('endvoting_rounds:{}， cancel proposal result:{}'.format(endvoting_rounds, result))
        assert_code(result, 302010)

        result = pip_obj.submitCancel(pip_obj.node.node_id, str(time.time()), endvoting_rounds + 1,
                                      proposalinfo.get('ProposalID'), pip_obj.node.staking_address,
                                      transaction_cfg=pip_obj.cfg.transaction_cfg)
        log.info('endvoting_rounds:{}， cancel proposal result:{}'.format(endvoting_rounds + 1, result))
        assert_code(result, 302010)

    @pytest.mark.P0
    def test_CP_ID_001(self, no_vp_proposal):
        pip_obj = no_vp_proposal
        result = pip_obj.submitCancel(pip_obj.node.node_id, str(time.time()), 1,
                                      '0x49b83cfc4b99462f7131d14d80c73b6657237753cd1e878e8d62dc2e9f574123',
                                      pip_obj.node.staking_address, transaction_cfg=pip_obj.cfg.transaction_cfg)
        log.info('cancel proposal result: {}'.format(result))
        assert_code(result, 302015)

    @pytest.mark.P0
    def test_CP_ID_004_CP_ID_003(self, new_genesis_env, client_consensus_obj):
        genesis = from_dict(data_class=Genesis, data=new_genesis_env.genesis_config)
        genesis.economicModel.gov.paramProposalVoteDurationSeconds = 0
        new_genesis_env.set_genesis(genesis.to_dict())
        new_genesis_env.deploy_all()
        pip_obj = client_consensus_obj.pip
        result = pip_obj.submitParam(pip_obj.node.node_id, str(time.time()), 'slashing', 'slashBlocksReward', '123',
                                     pip_obj.node.staking_address, transaction_cfg=pip_obj.cfg.transaction_cfg)
        log.info('Submit parameter proposal result : {}'.format(result))
        assert_code(result, 0)
        result = pip_obj.submitText(pip_obj.node.node_id, str(time.time()), pip_obj.node.staking_address,
                                    transaction_cfg=pip_obj.cfg.transaction_cfg)
        log.info('Submit text proposal result : {}'.format(result))
        assert_code(result, 0)
        proposalinfo_param = pip_obj.get_effect_proposal_info_of_vote(pip_obj.cfg.param_proposal)
        log.info('Get parameter proposal information : {}'.format(proposalinfo_param))
        proposalinfo_text = pip_obj.get_effect_proposal_info_of_vote(pip_obj.cfg.text_proposal)
        log.info('Get text proposal information : {}'.format(proposalinfo_text))
        result = pip_obj.submitCancel(pip_obj.node.node_id, str(time.time()), 1, proposalinfo_text.get('ProposalID'),
                                      pip_obj.node.staking_address, transaction_cfg=pip_obj.cfg.transaction_cfg)
        log.info('Submit cancel proposal result : {}'.format(result))
        wait_block_number(pip_obj.node, proposalinfo_param.get('EndVotingBlock'))
        result = pip_obj.submitCancel(pip_obj.node.node_id, str(time.time()), 1, proposalinfo_param.get('ProposalID'),
                                      pip_obj.node.staking_address, transaction_cfg=pip_obj.cfg.transaction_cfg)
        log.info('Submit cancel proposal result : {}'.format(result))


class TestPP():
    @pytest.mark.P0
    def test_PP_SU_001_PP_SU_002(self, no_vp_proposal, client_list_obj):
        pip_obj = no_vp_proposal
        client_obj = get_client_obj(pip_obj.node.node_id, client_list_obj)
        result = pip_obj.submitParam(pip_obj.node.node_id, str(time.time()), 'slashing', 'slashBlocksReward', '',
                                     pip_obj.node.staking_address, transaction_cfg=pip_obj.cfg.transaction_cfg)
        log.info('Submit param proposal result : {}'.format(result))
        assert_code(result, 3)

        result = pip_obj.submitParam(pip_obj.node.node_id, str(time.time()), 'slashing', 'slashBlocksReward', '1.1',
                                     pip_obj.node.staking_address, transaction_cfg=pip_obj.cfg.transaction_cfg)
        log.info('Submit param proposal result : {}'.format(result))
        assert_code(result, 3)

        result = pip_obj.submitParam(pip_obj.node.node_id, str(time.time()), 'slashing', 'slashBlocksReward', '-1',
                                     pip_obj.node.staking_address, transaction_cfg=pip_obj.cfg.transaction_cfg)
        log.info('Submit param proposal result : {}'.format(result))
        assert_code(result, 3)

        result = pip_obj.submitParam(pip_obj.node.node_id, str(time.time()), 'slashing', 'slashBlocksReward', '60101',
                                     pip_obj.node.staking_address, transaction_cfg=pip_obj.cfg.transaction_cfg)
        log.info('Submit param proposal result : {}'.format(result))
        assert_code(result, 3)

        result = pip_obj.submitParam(pip_obj.node.node_id, str(time.time()), 'slashing', 'slashBlocksReward', 1,
                                     pip_obj.node.staking_address, transaction_cfg=pip_obj.cfg.transaction_cfg)
        log.info('Submit param proposal result : {}'.format(result))
        assert_code(result, 3)

        result = pip_obj.submitParam(pip_obj.node.node_id, str(time.time()), 'slashing', 'slashBlocksReward', '1000000',
                                     pip_obj.node.staking_address, transaction_cfg=pip_obj.cfg.transaction_cfg)
        log.info('Submit param proposal result : {}'.format(result))
        assert_code(result, 3)

        result = pip_obj.submitParam(pip_obj.node.node_id, str(time.time()), 'slashing', 'slashBlocksReward',
                                     str(get_governable_parameter_value(client_obj, 'slashBlocksReward')),
                                     pip_obj.node.staking_address, transaction_cfg=pip_obj.cfg.transaction_cfg)
        log.info('Submit param proposal result : {}'.format(result))
        assert_code(result, 302034)

        if pip_obj.economic.slash_blocks_reward != 0:
            result = pip_obj.submitParam(pip_obj.node.node_id, str(time.time()), 'slashing', 'slashBlocksReward', '0',
                                         pip_obj.node.staking_address, transaction_cfg=pip_obj.cfg.transaction_cfg)
            log.info('Submit param proposal result : {}'.format(result))
            assert_code(result, 0)

    @pytest.mark.P0
    def test_PP_SU_002(self, no_vp_proposal, client_list_obj):
        pip_obj = no_vp_proposal
        if str(get_governable_parameter_value(client_list_obj[0], 'slashBlocksReward')) != 60100:
            result = pip_obj.submitParam(pip_obj.node.node_id, str(time.time()), 'slashing', 'slashBlocksReward', '60100',
                                         pip_obj.node.staking_address, transaction_cfg=pip_obj.cfg.transaction_cfg)
            log.info('Submit param proposal result : {}'.format(result))
            assert_code(result, 0)

    @pytest.mark.P0
    def test_PP_SU_003_PP_SU_004(self, no_vp_proposal, client_list_obj):
        pip_obj = no_vp_proposal
        client_obj = get_client_obj(pip_obj.node.node_id, client_list_obj)
        result = pip_obj.submitParam(pip_obj.node.node_id, str(time.time()), 'slashing', 'maxEvidenceAge', '',
                                     pip_obj.node.staking_address, transaction_cfg=pip_obj.cfg.transaction_cfg)
        log.info('Submit param proposal result : {}'.format(result))
        assert_code(result, 3)

        result = pip_obj.submitParam(pip_obj.node.node_id, str(time.time()), 'slashing', 'maxEvidenceAge', '1.1',
                                     pip_obj.node.staking_address, transaction_cfg=pip_obj.cfg.transaction_cfg)
        log.info('Submit param proposal result : {}'.format(result))
        assert_code(result, 3)

        result = pip_obj.submitParam(pip_obj.node.node_id, str(time.time()), 'slashing', 'maxEvidenceAge', '-1',
                                     pip_obj.node.staking_address, transaction_cfg=pip_obj.cfg.transaction_cfg)
        log.info('Submit param proposal result : {}'.format(result))
        assert_code(result, 3)

        result = pip_obj.submitParam(pip_obj.node.node_id, str(time.time()), 'slashing', 'maxEvidenceAge', '0',
                                     pip_obj.node.staking_address, transaction_cfg=pip_obj.cfg.transaction_cfg)
        log.info('Submit param proposal result : {}'.format(result))
        assert_code(result, 3)

        result = pip_obj.submitParam(pip_obj.node.node_id, str(time.time()), 'slashing', 'maxEvidenceAge', 1,
                                     pip_obj.node.staking_address, transaction_cfg=pip_obj.cfg.transaction_cfg)
        log.info('Submit param proposal result : {}'.format(result))
        assert_code(result, 3)

        result = pip_obj.submitParam(pip_obj.node.node_id, str(time.time()), 'slashing', 'maxEvidenceAge',
                                     str(get_governable_parameter_value(client_obj, 'maxEvidenceAge')),
                                     pip_obj.node.staking_address, transaction_cfg=pip_obj.cfg.transaction_cfg)
        log.info('Submit param proposal result : {}'.format(result))
        assert_code(result, 302034)

        result = pip_obj.submitParam(pip_obj.node.node_id, str(time.time()), 'slashing', 'maxEvidenceAge',
                                     str(get_governable_parameter_value(client_obj, 'unStakeFreezeDuration')),
                                     pip_obj.node.staking_address, transaction_cfg=pip_obj.cfg.transaction_cfg)
        log.info('Submit param proposal result : {}'.format(result))
        assert_code(result, 3)

        if int(get_governable_parameter_value(client_obj, 'maxEvidenceAge')) != 1:
            result = pip_obj.submitParam(pip_obj.node.node_id, str(time.time()), 'slashing', 'maxEvidenceAge', '1',
                                         pip_obj.node.staking_address, transaction_cfg=pip_obj.cfg.transaction_cfg)
            log.info('Submit param proposal result : {}'.format(result))
            assert_code(result, 0)

    @pytest.mark.P0
    def test_PP_SU_004(self, no_vp_proposal, client_list_obj):
        pip_obj = no_vp_proposal
        if str(get_governable_parameter_value(client_list_obj[0], 'maxEvidenceAge')) != str(
            pip_obj.economic.unstaking_freeze_ratio - 1) and str(pip_obj.economic.unstaking_freeze_ratio - 1) <= str(
                get_governable_parameter_value(client_list_obj[0], 'unStakeFreezeDuration')):
            result = pip_obj.submitParam(pip_obj.node.node_id, str(time.time()), 'slashing', 'maxEvidenceAge',
                                         str(pip_obj.economic.unstaking_freeze_ratio - 1),
                                         pip_obj.node.staking_address, transaction_cfg=pip_obj.cfg.transaction_cfg)
            log.info('Submit param proposal result : {}'.format(result))
            assert_code(result, 0)

    @pytest.mark.P0
    def test_PP_SU_005_PP_SU_006(self, no_vp_proposal, client_list_obj):
        pip_obj = no_vp_proposal
        client_obj = get_client_obj(pip_obj.node.node_id, client_list_obj)
        result = pip_obj.submitParam(pip_obj.node.node_id, str(time.time()), 'slashing', 'slashFractionDuplicateSign', '',
                                     pip_obj.node.staking_address, transaction_cfg=pip_obj.cfg.transaction_cfg)
        log.info('Submit param proposal result : {}'.format(result))
        assert_code(result, 3)

        result = pip_obj.submitParam(pip_obj.node.node_id, str(time.time()), 'slashing', 'slashFractionDuplicateSign', '1.1',
                                     pip_obj.node.staking_address, transaction_cfg=pip_obj.cfg.transaction_cfg)
        log.info('Submit param proposal result : {}'.format(result))
        assert_code(result, 3)

        result = pip_obj.submitParam(pip_obj.node.node_id, str(time.time()), 'slashing', 'slashFractionDuplicateSign', '-1',
                                     pip_obj.node.staking_address, transaction_cfg=pip_obj.cfg.transaction_cfg)
        log.info('Submit param proposal result : {}'.format(result))
        assert_code(result, 3)

        result = pip_obj.submitParam(pip_obj.node.node_id, str(time.time()), 'slashing', 'slashFractionDuplicateSign', '0',
                                     pip_obj.node.staking_address, transaction_cfg=pip_obj.cfg.transaction_cfg)
        log.info('Submit param proposal result : {}'.format(result))
        assert_code(result, 3)

        result = pip_obj.submitParam(pip_obj.node.node_id, str(time.time()), 'slashing', 'slashFractionDuplicateSign', 1,
                                     pip_obj.node.staking_address, transaction_cfg=pip_obj.cfg.transaction_cfg)
        log.info('Submit param proposal result : {}'.format(result))
        assert_code(result, 3)

        result = pip_obj.submitParam(pip_obj.node.node_id, str(time.time()), 'slashing', 'slashFractionDuplicateSign',
                                     '10001', pip_obj.node.staking_address, transaction_cfg=pip_obj.cfg.transaction_cfg)
        log.info('Submit param proposal result : {}'.format(result))
        assert_code(result, 3)

        result = pip_obj.submitParam(pip_obj.node.node_id, str(time.time()), 'slashing', 'slashFractionDuplicateSign',
                                     str(get_governable_parameter_value(client_obj, 'slashFractionDuplicateSign')),
                                     pip_obj.node.staking_address, transaction_cfg=pip_obj.cfg.transaction_cfg)
        log.info('Submit param proposal result : {}'.format(result))
        assert_code(result, 302034)

        if int(get_governable_parameter_value(client_obj, 'maxEvidenceAge')) != 1:
            result = pip_obj.submitParam(pip_obj.node.node_id, str(time.time()), 'slashing', 'slashFractionDuplicateSign', '1',
                                         pip_obj.node.staking_address, transaction_cfg=pip_obj.cfg.transaction_cfg)
            log.info('Submit param proposal result : {}'.format(result))
            assert_code(result, 0)

        if int(get_governable_parameter_value(client_obj, 'maxEvidenceAge')) != 10000:
            result = pip_obj.submitParam(pip_obj.node.node_id, str(time.time()), 'slashing', 'slashFractionDuplicateSign', '10000',
                                         pip_obj.node.staking_address, transaction_cfg=pip_obj.cfg.transaction_cfg)
            log.info('Submit param proposal result : {}'.format(result))
            assert_code(result, 0)

    @pytest.mark.P0
    def test_PP_SU_007_PP_SU_008(self, no_vp_proposal, client_list_obj):
        pip_obj = no_vp_proposal
        client_obj = get_client_obj(pip_obj.node.node_id, client_list_obj)
        result = pip_obj.submitParam(pip_obj.node.node_id, str(time.time()), 'slashing', 'duplicateSignReportReward', '',
                                     pip_obj.node.staking_address, transaction_cfg=pip_obj.cfg.transaction_cfg)
        log.info('Submit param proposal result : {}'.format(result))
        assert_code(result, 3)

        result = pip_obj.submitParam(pip_obj.node.node_id, str(time.time()), 'slashing', 'duplicateSignReportReward', '1.1',
                                     pip_obj.node.staking_address, transaction_cfg=pip_obj.cfg.transaction_cfg)
        log.info('Submit param proposal result : {}'.format(result))
        assert_code(result, 3)

        result = pip_obj.submitParam(pip_obj.node.node_id, str(time.time()), 'slashing', 'duplicateSignReportReward', '-1',
                                     pip_obj.node.staking_address, transaction_cfg=pip_obj.cfg.transaction_cfg)
        log.info('Submit param proposal result : {}'.format(result))
        assert_code(result, 3)

        result = pip_obj.submitParam(pip_obj.node.node_id, str(time.time()), 'slashing', 'duplicateSignReportReward', '0',
                                     pip_obj.node.staking_address, transaction_cfg=pip_obj.cfg.transaction_cfg)
        log.info('Submit param proposal result : {}'.format(result))
        assert_code(result, 3)

        result = pip_obj.submitParam(pip_obj.node.node_id, str(time.time()), 'slashing', 'duplicateSignReportReward', 1,
                                     pip_obj.node.staking_address, transaction_cfg=pip_obj.cfg.transaction_cfg)
        log.info('Submit param proposal result : {}'.format(result))
        assert_code(result, 3)

        result = pip_obj.submitParam(pip_obj.node.node_id, str(time.time()), 'slashing', 'duplicateSignReportReward',
                                     '81', pip_obj.node.staking_address, transaction_cfg=pip_obj.cfg.transaction_cfg)
        log.info('Submit param proposal result : {}'.format(result))
        assert_code(result, 3)

        result = pip_obj.submitParam(pip_obj.node.node_id, str(time.time()), 'slashing', 'duplicateSignReportReward',
                                     str(get_governable_parameter_value(client_obj, 'duplicateSignReportReward')),
                                     pip_obj.node.staking_address, transaction_cfg=pip_obj.cfg.transaction_cfg)
        log.info('Submit param proposal result : {}'.format(result))
        assert_code(result, 302034)

        if int(get_governable_parameter_value(client_obj, 'duplicateSignReportReward')) != 1:
            result = pip_obj.submitParam(pip_obj.node.node_id, str(time.time()), 'slashing', 'duplicateSignReportReward', '1',
                                         pip_obj.node.staking_address, transaction_cfg=pip_obj.cfg.transaction_cfg)
            log.info('Submit param proposal result : {}'.format(result))
            assert_code(result, 0)

    @pytest.mark.P0
    def test_PP_SU_008(self, no_vp_proposal, client_list_obj):
        pip_obj = no_vp_proposal
        if int(get_governable_parameter_value(client_list_obj[0], 'duplicateSignReportReward')) != 80:
            result = pip_obj.submitParam(pip_obj.node.node_id, str(time.time()), 'slashing', 'duplicateSignReportReward',
                                         '80', pip_obj.node.staking_address,
                                         transaction_cfg=pip_obj.cfg.transaction_cfg)
            log.info('Submit param proposal result : {}'.format(result))
            assert_code(result, 0)

    @pytest.mark.P0
    def test_PP_SU_009_PP_SU_010(self, no_vp_proposal, client_list_obj):
        pip_obj = no_vp_proposal
        client_obj = get_client_obj(pip_obj.node.node_id, client_list_obj)
        result = pip_obj.submitParam(pip_obj.node.node_id, str(time.time()), 'staking', 'stakeThreshold', '',
                                     pip_obj.node.staking_address, transaction_cfg=pip_obj.cfg.transaction_cfg)
        log.info('Submit param proposal result : {}'.format(result))
        assert_code(result, 3)

        value = 10**18 * 1000000 + 0.1
        result = pip_obj.submitParam(pip_obj.node.node_id, str(time.time()), 'staking', 'stakeThreshold', str(value),
                                     pip_obj.node.staking_address, transaction_cfg=pip_obj.cfg.transaction_cfg)
        log.info('Submit param proposal result : {}'.format(result))
        assert_code(result, 3)

        value = -10**18 * 1000000
        result = pip_obj.submitParam(pip_obj.node.node_id, str(time.time()), 'staking', 'stakeThreshold', str(value),
                                     pip_obj.node.staking_address, transaction_cfg=pip_obj.cfg.transaction_cfg)
        log.info('Submit param proposal result : {}'.format(result))
        assert_code(result, 3)

        result = pip_obj.submitParam(pip_obj.node.node_id, str(time.time()), 'staking', 'stakeThreshold', '0',
                                     pip_obj.node.staking_address, transaction_cfg=pip_obj.cfg.transaction_cfg)
        log.info('Submit param proposal result : {}'.format(result))
        assert_code(result, 3)

        result = pip_obj.submitParam(pip_obj.node.node_id, str(time.time()), 'staking', 'stakeThreshold', 10**18 * 1000000,
                                     pip_obj.node.staking_address, transaction_cfg=pip_obj.cfg.transaction_cfg)
        log.info('Submit param proposal result : {}'.format(result))
        assert_code(result, 3)

        value = 10**18 * 1000000 - 1
        result = pip_obj.submitParam(pip_obj.node.node_id, str(time.time()), 'staking', 'stakeThreshold',
                                     str(value), pip_obj.node.staking_address, transaction_cfg=pip_obj.cfg.transaction_cfg)
        log.info('Submit param proposal result : {}'.format(result))
        assert_code(result, 3)

        value = 10**18 * 10000000
        result = pip_obj.submitParam(pip_obj.node.node_id, str(time.time()), 'staking', 'stakeThreshold',
                                     str(value), pip_obj.node.staking_address, transaction_cfg=pip_obj.cfg.transaction_cfg)
        log.info('Submit param proposal result : {}'.format(result))
        assert_code(result, 3)

        result = pip_obj.submitParam(pip_obj.node.node_id, str(time.time()), 'staking', 'stakeThreshold',
                                     str(get_governable_parameter_value(client_obj, 'stakeThreshold')),
                                     pip_obj.node.staking_address, transaction_cfg=pip_obj.cfg.transaction_cfg)
        log.info('Submit param proposal result : {}'.format(result))
        assert_code(result, 302034)

        if int(get_governable_parameter_value(client_obj, 'stakeThreshold')) != 10**18 * 1000000:
            value = 10**18 * 1000000
            result = pip_obj.submitParam(pip_obj.node.node_id, str(time.time()), 'staking', 'stakeThreshold', str(value),
                                         pip_obj.node.staking_address, transaction_cfg=pip_obj.cfg.transaction_cfg)
            log.info('Submit param proposal result : {}'.format(result))
            assert_code(result, 0)

    @pytest.mark.P0
    def test_PP_SU_010(self, no_vp_proposal):
        pip_obj = no_vp_proposal
        value = 10**18 * 10000000 - 1
        result = pip_obj.submitParam(pip_obj.node.node_id, str(time.time()), 'staking', 'stakeThreshold',
                                     str(value),
                                     pip_obj.node.staking_address, transaction_cfg=pip_obj.cfg.transaction_cfg)
        log.info('Submit param proposal result : {}'.format(result))
        assert_code(result, 0)

    @pytest.mark.P0
    def test_PP_SU_011_PP_SU_012(self, no_vp_proposal, client_list_obj):
        pip_obj = no_vp_proposal
        client_obj = get_client_obj(pip_obj.node.node_id, client_list_obj)
        result = pip_obj.submitParam(pip_obj.node.node_id, str(time.time()), 'staking', 'operatingThreshold', '',
                                     pip_obj.node.staking_address, transaction_cfg=pip_obj.cfg.transaction_cfg)
        log.info('Submit param proposal result : {}'.format(result))
        assert_code(result, 3)

        value = 10**18 * 10 + 0.5
        result = pip_obj.submitParam(pip_obj.node.node_id, str(time.time()), 'staking', 'operatingThreshold', str(value),
                                     pip_obj.node.staking_address, transaction_cfg=pip_obj.cfg.transaction_cfg)
        log.info('Submit param proposal result : {}'.format(result))
        assert_code(result, 3)

        value = -10**18 * 10
        result = pip_obj.submitParam(pip_obj.node.node_id, str(time.time()), 'staking', 'operatingThreshold', str(value),
                                     pip_obj.node.staking_address, transaction_cfg=pip_obj.cfg.transaction_cfg)
        log.info('Submit param proposal result : {}'.format(result))
        assert_code(result, 3)

        result = pip_obj.submitParam(pip_obj.node.node_id, str(time.time()), 'staking', 'operatingThreshold', '0',
                                     pip_obj.node.staking_address, transaction_cfg=pip_obj.cfg.transaction_cfg)
        log.info('Submit param proposal result : {}'.format(result))
        assert_code(result, 3)

        result = pip_obj.submitParam(pip_obj.node.node_id, str(time.time()), 'staking', 'operatingThreshold', 1,
                                     pip_obj.node.staking_address, transaction_cfg=pip_obj.cfg.transaction_cfg)
        log.info('Submit param proposal result : {}'.format(result))
        assert_code(result, 3)

        value = 10**18 * 10 - 1
        result = pip_obj.submitParam(pip_obj.node.node_id, str(time.time()), 'staking', 'operatingThreshold',
                                     str(value), pip_obj.node.staking_address, transaction_cfg=pip_obj.cfg.transaction_cfg)
        log.info('Submit param proposal result : {}'.format(result))
        assert_code(result, 3)

        value = 10**18 * 10000
        result = pip_obj.submitParam(pip_obj.node.node_id, str(time.time()), 'staking', 'operatingThreshold',
                                     str(value), pip_obj.node.staking_address, transaction_cfg=pip_obj.cfg.transaction_cfg)
        log.info('Submit param proposal result : {}'.format(result))
        assert_code(result, 3)

        result = pip_obj.submitParam(pip_obj.node.node_id, str(time.time()), 'staking', 'operatingThreshold',
                                     str(get_governable_parameter_value(client_obj, 'operatingThreshold')),
                                     pip_obj.node.staking_address, transaction_cfg=pip_obj.cfg.transaction_cfg)
        log.info('Submit param proposal result : {}'.format(result))
        assert_code(result, 302034)

        if int(get_governable_parameter_value(client_obj, 'operatingThreshold')) != 10**18 * 10:
            value = 10**18 * 10
            result = pip_obj.submitParam(pip_obj.node.node_id, str(time.time()), 'staking', 'operatingThreshold', str(value),
                                         pip_obj.node.staking_address, transaction_cfg=pip_obj.cfg.transaction_cfg)
            log.info('Submit param proposal result : {}'.format(result))
            assert_code(result, 0)

    @pytest.mark.P0
    def test_PP_SU_012(self, no_vp_proposal):
        pip_obj = no_vp_proposal
        value = 10**18 * 10000 - 1
        result = pip_obj.submitParam(pip_obj.node.node_id, str(time.time()), 'staking', 'operatingThreshold',
                                     str(value), pip_obj.node.staking_address,
                                     transaction_cfg=pip_obj.cfg.transaction_cfg)
        log.info('Submit param proposal result : {}'.format(result))
        assert_code(result, 0)

    @pytest.mark.P0
    def test_PP_SU_013_PP_SU_014(self, no_vp_proposal, client_list_obj):
        pip_obj = no_vp_proposal
        client_obj = get_client_obj(pip_obj.node.node_id, client_list_obj)
        result = pip_obj.submitParam(pip_obj.node.node_id, str(time.time()), 'staking', 'unStakeFreezeDuration', '',
                                     pip_obj.node.staking_address, transaction_cfg=pip_obj.cfg.transaction_cfg)
        log.info('Submit param proposal result : {}'.format(result))
        assert_code(result, 3)

        result = pip_obj.submitParam(pip_obj.node.node_id, str(time.time()), 'staking', 'unStakeFreezeDuration', '10.5',
                                     pip_obj.node.staking_address, transaction_cfg=pip_obj.cfg.transaction_cfg)
        log.info('Submit param proposal result : {}'.format(result))
        assert_code(result, 3)

        result = pip_obj.submitParam(pip_obj.node.node_id, str(time.time()), 'staking', 'unStakeFreezeDuration', '-100',
                                     pip_obj.node.staking_address, transaction_cfg=pip_obj.cfg.transaction_cfg)
        log.info('Submit param proposal result : {}'.format(result))
        assert_code(result, 3)

        result = pip_obj.submitParam(pip_obj.node.node_id, str(time.time()), 'staking', 'unStakeFreezeDuration', '0',
                                     pip_obj.node.staking_address, transaction_cfg=pip_obj.cfg.transaction_cfg)
        log.info('Submit param proposal result : {}'.format(result))
        assert_code(result, 3)

        result = pip_obj.submitParam(pip_obj.node.node_id, str(time.time()), 'staking', 'unStakeFreezeDuration', 11,
                                     pip_obj.node.staking_address, transaction_cfg=pip_obj.cfg.transaction_cfg)
        log.info('Submit param proposal result : {}'.format(result))
        assert_code(result, 3)

        result = pip_obj.submitParam(pip_obj.node.node_id, str(time.time()), 'staking', 'unStakeFreezeDuration',
                                     '113', pip_obj.node.staking_address, transaction_cfg=pip_obj.cfg.transaction_cfg)
        log.info('Submit param proposal result : {}'.format(result))
        assert_code(result, 3)

        result = pip_obj.submitParam(pip_obj.node.node_id, str(time.time()), 'staking', 'unStakeFreezeDuration',
                                     str(get_governable_parameter_value(client_obj, 'unStakeFreezeDuration')),
                                     pip_obj.node.staking_address, transaction_cfg=pip_obj.cfg.transaction_cfg)
        log.info('Submit param proposal result : {}'.format(result))
        assert_code(result, 302034)

        result = pip_obj.submitParam(pip_obj.node.node_id, str(time.time()), 'staking', 'unStakeFreezeDuration',
                                     str(get_governable_parameter_value(client_obj, 'maxEvidenceAge')),
                                     pip_obj.node.staking_address, transaction_cfg=pip_obj.cfg.transaction_cfg)
        log.info('Submit param proposal result : {}'.format(result))
        assert_code(result, 3)

        if int(get_governable_parameter_value(client_obj, 'unStakeFreezeDuration')) != 112:
            result = pip_obj.submitParam(pip_obj.node.node_id, str(time.time()), 'staking', 'unStakeFreezeDuration', '112',
                                         pip_obj.node.staking_address, transaction_cfg=pip_obj.cfg.transaction_cfg)
            log.info('Submit param proposal result : {}'.format(result))
            assert_code(result, 0)

    @pytest.mark.P0
    def test__PP_SU_014(self, no_vp_proposal, client_list_obj):
        pip_obj = no_vp_proposal
        if int(get_governable_parameter_value(client_list_obj[0], 'unStakeFreezeDuration')) != str(
                int(get_governable_parameter_value(client_list_obj[0], 'maxEvidenceAge')) - 1):
            result = pip_obj.submitParam(pip_obj.node.node_id, str(time.time()), 'staking', 'unStakeFreezeDuration',
                                         str(int(get_governable_parameter_value(client_list_obj[0], 'maxEvidenceAge')) + 5),
                                         pip_obj.node.staking_address, transaction_cfg=pip_obj.cfg.transaction_cfg)
            log.info('Submit param proposal result : {}'.format(result))
            assert_code(result, 0)

    @pytest.mark.P0
    def test_PP_SU_015_PP_SU_016(self, no_vp_proposal, client_list_obj):
        pip_obj = no_vp_proposal
        client_obj = get_client_obj(pip_obj.node.node_id, client_list_obj)
        result = pip_obj.submitParam(pip_obj.node.node_id, str(time.time()), 'staking', 'maxValidators', '',
                                     pip_obj.node.staking_address, transaction_cfg=pip_obj.cfg.transaction_cfg)
        log.info('Submit param proposal result : {}'.format(result))
        assert_code(result, 3)

        result = pip_obj.submitParam(pip_obj.node.node_id, str(time.time()), 'staking', 'maxValidators', '30.5',
                                     pip_obj.node.staking_address, transaction_cfg=pip_obj.cfg.transaction_cfg)
        log.info('Submit param proposal result : {}'.format(result))
        assert_code(result, 3)

        result = pip_obj.submitParam(pip_obj.node.node_id, str(time.time()), 'staking', 'maxValidators', '-100',
                                     pip_obj.node.staking_address, transaction_cfg=pip_obj.cfg.transaction_cfg)
        log.info('Submit param proposal result : {}'.format(result))
        assert_code(result, 3)

        result = pip_obj.submitParam(pip_obj.node.node_id, str(time.time()), 'staking', 'maxValidators', '0',
                                     pip_obj.node.staking_address, transaction_cfg=pip_obj.cfg.transaction_cfg)
        log.info('Submit param proposal result : {}'.format(result))
        assert_code(result, 3)

        result = pip_obj.submitParam(pip_obj.node.node_id, str(time.time()), 'staking', 'maxValidators', 25,
                                     pip_obj.node.staking_address, transaction_cfg=pip_obj.cfg.transaction_cfg)
        log.info('Submit param proposal result : {}'.format(result))
        assert_code(result, 3)

        result = pip_obj.submitParam(pip_obj.node.node_id, str(time.time()), 'staking', 'maxValidators',
                                     '3', pip_obj.node.staking_address, transaction_cfg=pip_obj.cfg.transaction_cfg)
        log.info('Submit param proposal result : {}'.format(result))
        assert_code(result, 3)

        result = pip_obj.submitParam(pip_obj.node.node_id, str(time.time()), 'staking', 'maxValidators',
                                     '202', pip_obj.node.staking_address, transaction_cfg=pip_obj.cfg.transaction_cfg)
        log.info('Submit param proposal result : {}'.format(result))
        assert_code(result, 3)

        result = pip_obj.submitParam(pip_obj.node.node_id, str(time.time()), 'staking', 'maxValidators',
                                     str(get_governable_parameter_value(client_obj, 'maxValidators')),
                                     pip_obj.node.staking_address, transaction_cfg=pip_obj.cfg.transaction_cfg)
        log.info('Submit param proposal result : {}'.format(result))
        assert_code(result, 302034)

        if int(get_governable_parameter_value(client_obj, 'maxValidators')) != 4:
            result = pip_obj.submitParam(pip_obj.node.node_id, str(time.time()), 'staking', 'maxValidators', '4',
                                         pip_obj.node.staking_address, transaction_cfg=pip_obj.cfg.transaction_cfg)
            log.info('Submit param proposal result : {}'.format(result))
            assert_code(result, 0)

    def test_PP_SU_016(self, no_vp_proposal, client_list_obj):
        pip_obj = no_vp_proposal
        if int(get_governable_parameter_value(client_list_obj[0], 'maxValidators')) != 201:
            result = pip_obj.submitParam(pip_obj.node.node_id, str(time.time()), 'staking', 'maxValidators', '201',
                                         pip_obj.node.staking_address, transaction_cfg=pip_obj.cfg.transaction_cfg)
            log.info('Submit param proposal result : {}'.format(result))
            assert_code(result, 0)

    @pytest.mark.P0
    def test_PP_SU_016(self, no_vp_proposal, client_list_obj):
        pip_obj = no_vp_proposal
        client_obj = get_client_obj(pip_obj.node.node_id, client_list_obj)
        if int(get_governable_parameter_value(client_obj, 'maxValidators')) != 201:
            result = pip_obj.submitParam(pip_obj.node.node_id, str(time.time()), 'staking', 'maxValidators', '201',
                                         pip_obj.node.staking_address, transaction_cfg=pip_obj.cfg.transaction_cfg)
            log.info('Submit param proposal result : {}'.format(result))
            assert_code(result, 0)

    @pytest.mark.P0
    def test_PP_SU_017_PP_SU_018(self, no_vp_proposal, client_list_obj):
        pip_obj = no_vp_proposal
        client_obj = get_client_obj(pip_obj.node.node_id, client_list_obj)
        result = pip_obj.submitParam(pip_obj.node.node_id, str(time.time()), 'block', 'maxBlockGasLimit', '',
                                     pip_obj.node.staking_address, transaction_cfg=pip_obj.cfg.transaction_cfg)
        log.info('Submit param proposal result : {}'.format(result))
        assert_code(result, 3)

        result = pip_obj.submitParam(pip_obj.node.node_id, str(time.time()), 'block', 'maxBlockGasLimit', '4712388.5',
                                     pip_obj.node.staking_address, transaction_cfg=pip_obj.cfg.transaction_cfg)
        log.info('Submit param proposal result : {}'.format(result))
        assert_code(result, 3)

        result = pip_obj.submitParam(pip_obj.node.node_id, str(time.time()), 'block', 'maxBlockGasLimit', '-4712388',
                                     pip_obj.node.staking_address, transaction_cfg=pip_obj.cfg.transaction_cfg)
        log.info('Submit param proposal result : {}'.format(result))
        assert_code(result, 3)

        result = pip_obj.submitParam(pip_obj.node.node_id, str(time.time()), 'block', 'maxBlockGasLimit', '0',
                                     pip_obj.node.staking_address, transaction_cfg=pip_obj.cfg.transaction_cfg)
        log.info('Submit param proposal result : {}'.format(result))
        assert_code(result, 3)

        result = pip_obj.submitParam(pip_obj.node.node_id, str(time.time()), 'block', 'maxBlockGasLimit', 4712388,
                                     pip_obj.node.staking_address, transaction_cfg=pip_obj.cfg.transaction_cfg)
        log.info('Submit param proposal result : {}'.format(result))
        assert_code(result, 3)

        result = pip_obj.submitParam(pip_obj.node.node_id, str(time.time()), 'block', 'maxBlockGasLimit',
                                     '4712387', pip_obj.node.staking_address, transaction_cfg=pip_obj.cfg.transaction_cfg)
        log.info('Submit param proposal result : {}'.format(result))
        assert_code(result, 3)

        result = pip_obj.submitParam(pip_obj.node.node_id, str(time.time()), 'block', 'maxBlockGasLimit',
                                     '210000001', pip_obj.node.staking_address, transaction_cfg=pip_obj.cfg.transaction_cfg)
        log.info('Submit param proposal result : {}'.format(result))
        assert_code(result, 3)

        result = pip_obj.submitParam(pip_obj.node.node_id, str(time.time()), 'block', 'maxBlockGasLimit',
                                     str(get_governable_parameter_value(client_obj, 'maxBlockGasLimit')),
                                     pip_obj.node.staking_address, transaction_cfg=pip_obj.cfg.transaction_cfg)
        log.info('Submit param proposal result : {}'.format(result))
        assert_code(result, 302034)

        if int(get_governable_parameter_value(client_obj, 'maxBlockGasLimit')) != 4712388:
            result = pip_obj.submitParam(pip_obj.node.node_id, str(time.time()), 'block', 'maxBlockGasLimit', '4712388',
                                         pip_obj.node.staking_address, transaction_cfg=pip_obj.cfg.transaction_cfg)
            log.info('Submit param proposal result : {}'.format(result))
            assert_code(result, 0)

    @pytest.mark.P0
    def test_PP_SU_018(self, no_vp_proposal):
        pip_obj = no_vp_proposal
        result = pip_obj.submitParam(pip_obj.node.node_id, str(time.time()), 'block', 'maxBlockGasLimit', '210000000',
                                     pip_obj.node.staking_address, transaction_cfg=pip_obj.cfg.transaction_cfg)
        log.info('Submit param proposal result : {}'.format(result))
        assert_code(result, 0)

    @pytest.mark.P0
    def test_PP_SU_019(self, no_vp_proposal):
        pip_obj = no_vp_proposal
        result = pip_obj.submitParam(pip_obj.node.node_id, str(time.time()), 'block', '', '1',
                                     pip_obj.node.staking_address,
                                     transaction_cfg=pip_obj.cfg.transaction_cfg)
        log.info('Submit param proposal result : {}'.format(result))
        assert_code(result, 302031)

        result = pip_obj.submitParam(pip_obj.node.node_id, str(time.time()), 'block', 'unStakeFreezeDuration', '100',
                                     pip_obj.node.staking_address,
                                     transaction_cfg=pip_obj.cfg.transaction_cfg)
        log.info('Submit param proposal result : {}'.format(result))
        assert_code(result, 302031)

        result = pip_obj.submitParam(pip_obj.node.node_id, str(time.time()), 'slashing', 'SlashBlocksReward', '100',
                                     pip_obj.node.staking_address,
                                     transaction_cfg=pip_obj.cfg.transaction_cfg)
        log.info('Submit param proposal result : {}'.format(result))
        assert_code(result, 302031)

        result = pip_obj.submitParam(pip_obj.node.node_id, str(time.time()), 'slashing', 'slash BlocksReward', '100',
                                     pip_obj.node.staking_address,
                                     transaction_cfg=pip_obj.cfg.transaction_cfg)
        log.info('Submit param proposal result : {}'.format(result))
        assert_code(result, 302031)

        result = pip_obj.submitParam(pip_obj.node.node_id, str(time.time()), 'slashing', 'slashBlocks./,.Reward', '100',
                                     pip_obj.node.staking_address,
                                     transaction_cfg=pip_obj.cfg.transaction_cfg)
        log.info('Submit param proposal result : {}'.format(result))
        assert_code(result, 302031)

    @pytest.mark.P0
    def test_PP_SU_020(self, no_vp_proposal):
        pip_obj = no_vp_proposal
        result = pip_obj.submitParam(pip_obj.node.node_id, str(time.time()), '', 'slashBlocksReward', '100',
                                     pip_obj.node.staking_address,
                                     transaction_cfg=pip_obj.cfg.transaction_cfg)
        log.info('Submit param proposal result : {}'.format(result))
        assert_code(result, 302031)

        result = pip_obj.submitParam(pip_obj.node.node_id, str(time.time()), 'Slashing', 'slashBlocksReward', '100',
                                     pip_obj.node.staking_address,
                                     transaction_cfg=pip_obj.cfg.transaction_cfg)
        log.info('Submit param proposal result : {}'.format(result))
        assert_code(result, 302031)

        result = pip_obj.submitParam(pip_obj.node.node_id, str(time.time()), 'sLashing', 'SlashBlocksReward', '100',
                                     pip_obj.node.staking_address,
                                     transaction_cfg=pip_obj.cfg.transaction_cfg)
        log.info('Submit param proposal result : {}'.format(result))
        assert_code(result, 302031)

        result = pip_obj.submitParam(pip_obj.node.node_id, str(time.time()), 'slashing123', 'slash BlocksReward', '100',
                                     pip_obj.node.staking_address,
                                     transaction_cfg=pip_obj.cfg.transaction_cfg)
        log.info('Submit param proposal result : {}'.format(result))
        assert_code(result, 302031)

        result = pip_obj.submitParam(pip_obj.node.node_id, str(time.time()), 's lashing', 'slashBlocks./,.Reward', '100',
                                     pip_obj.node.staking_address,
                                     transaction_cfg=pip_obj.cfg.transaction_cfg)
        log.info('Submit param proposal result : {}'.format(result))
        assert_code(result, 302031)

        result = pip_obj.submitParam(pip_obj.node.node_id, str(time.time()), 's.,.lashing', 'slashBlocks./,.Reward', '100',
                                     pip_obj.node.staking_address,
                                     transaction_cfg=pip_obj.cfg.transaction_cfg)
        log.info('Submit param proposal result : {}'.format(result))
        assert_code(result, 302031)


class TestSubmitPPAbnormal():
    @pytest.mark.P0
    def test_PP_PR_002(self, no_vp_proposal, client_new_node_obj):
        pip_obj = client_new_node_obj.pip
        address, _ = pip_obj.economic.account.generate_account(pip_obj.node.web3, 10**18 * 10000)
        result = pip_obj.submitParam(pip_obj.node.node_id, str(time.time()), 'slashing', 'slashBlocksReward', '88',
                                     address, transaction_cfg=pip_obj.cfg.transaction_cfg)
        log.info('new node submit param proposal result : {}'.format(result))
        assert_code(result, 302022)

    @pytest.mark.P0
    def test_PP_PR_001(self, no_vp_proposal, client_candidate_obj):
        pip_obj = client_candidate_obj.pip
        result = pip_obj.submitParam(pip_obj.node.node_id, str(time.time()), 'slashing', 'slashBlocksReward', '87',
                                     pip_obj.node.staking_address, transaction_cfg=pip_obj.cfg.transaction_cfg)
        log.info('candidate submit param proposal result :{}'.format(result))
        assert_code(result, 302022)

    @pytest.mark.P0
    def test_PP_UN_002(self, submit_version):
        pip_obj = submit_version
        result = pip_obj.submitParam(pip_obj.node.node_id, str(time.time()), 'slashing', 'slashBlocksReward', '99',
                                     pip_obj.node.staking_address, transaction_cfg=pip_obj.cfg.transaction_cfg)
        log.info('There is voting version proposal, submit a param proposal : {}'.format(result))
        assert_code(result, 302012)

    @pytest.mark.P2
    def test_PP_PR_003_PP_PR_004(self, no_vp_proposal, client_list_obj):
        pip_obj = no_vp_proposal
        client_obj = get_client_obj(pip_obj.node.node_id, client_list_obj)
        address = pip_obj.node.staking_address
        result = client_obj.staking.withdrew_staking(address)
        log.info('nodeid: {} withdrewstaking result: {}'.format(client_obj.node.node_id, result))
        assert_code(result, 0)
        result = pip_obj.submitParam(pip_obj.node.node_id, str(time.time()), 'slashing', 'slashBlocksReward',
                                     '86', address, transaction_cfg=pip_obj.cfg.transaction_cfg)
        log.info('node exiting，param proposal result: {}'.format(result))
        assert_code(result, 302020)

        client_obj.economic.wait_settlement_blocknum(client_obj.node,
                                                     number=client_obj.economic.unstaking_freeze_ratio)
        result = pip_obj.submitParam(pip_obj.node.node_id, str(time.time()), 'slashing', 'slashBlocksReward',
                                     '86', address, transaction_cfg=pip_obj.cfg.transaction_cfg)
        log.info('exited node，cancel proposal result: {}'.format(result))
        assert_code(result, 302022)

    @pytest.mark.P2
    def test_PP_WA_001(self, no_vp_proposal):
        pip_obj = no_vp_proposal
        address, _ = pip_obj.economic.account.generate_account(pip_obj.node.web3, 10**18 * 10000)
        result = pip_obj.submitParam(pip_obj.node.node_id, str(time.time()), 'slashing', 'slashBlocksReward', '87',
                                     address, transaction_cfg=pip_obj.cfg.transaction_cfg)
        log.info('candidate submit param proposal result :{}'.format(result))
        assert_code(result, 302021)


class TestSubmitAgain():
    @pytest.mark.P2
    def test_PP_TI_001_002(self, new_genesis_env, client_con_list_obj):
        genesis = from_dict(data_class=Genesis, data=new_genesis_env.genesis_config)
        genesis.economicModel.gov.paramProposalVoteDurationSeconds = 0
        new_genesis_env.set_genesis(genesis.to_dict())
        new_genesis_env.deploy_all()
        submitcppandvote(client_con_list_obj[:3], [1, 1, 1])
        pip_obj = client_con_list_obj[0].pip
        proposalinfo_param = pip_obj.get_effect_proposal_info_of_vote(pip_obj.cfg.param_proposal)
        log.info('Param proposal information : {}'.format(proposalinfo_param))
        proposalinfo_cancel = pip_obj.get_effect_proposal_info_of_vote(pip_obj.cfg.cancel_proposal)
        log.info('Cancel proposal information : {}'.format(proposalinfo_cancel))
        wait_block_number(pip_obj.node, proposalinfo_cancel.get('EndVotingBlock'))
        assert_code(pip_obj.get_status_of_proposal(proposalinfo_cancel.get('ProposalID')), 2)
        assert_code(pip_obj.get_status_of_proposal(proposalinfo_param.get('ProposalID')), 6)

        result = pip_obj.submitParam(pip_obj.node.node_id, str(time.time()), 'slashing', 'slashBlocksReward', '998',
                                     pip_obj.node.staking_address, transaction_cfg=pip_obj.cfg.transaction_cfg)
        log.info('Submit param proposal result : {}'.format(result))
        assert_code(result, 0)
        proposalinfo_param = pip_obj.get_effect_proposal_info_of_vote(pip_obj.cfg.param_proposal)
        log.info('Param proposal information : {}'.format(proposalinfo_param))
        wait_block_number(pip_obj.node, proposalinfo_param.get('EndVotingBlock'))

        result = pip_obj.submitParam(pip_obj.node.node_id, str(time.time()), 'slashing', 'slashBlocksReward', '998',
                                     pip_obj.node.staking_address, transaction_cfg=pip_obj.cfg.transaction_cfg)
        log.info('Submit param proposal result : {}'.format(result))
        assert_code(result, 0)

    @pytest.mark.P2
    def test_VP_TI_001_002(self, no_vp_proposal, client_verifier_obj_list):
        submitcvpandvote(client_verifier_obj_list[:3], 1, 1, 1)
        pip_obj = client_verifier_obj_list[0].pip
        proposalinfo_version = pip_obj.get_effect_proposal_info_of_vote(pip_obj.cfg.version_proposal)
        log.info('Version proposal information : {}'.format(proposalinfo_version))
        proposalinfo_cancel = pip_obj.get_effect_proposal_info_of_vote(pip_obj.cfg.cancel_proposal)
        log.info('Cancel proposal information : {}'.format(proposalinfo_cancel))
        wait_block_number(pip_obj.node, proposalinfo_cancel.get('EndVotingBlock'))
        assert_code(pip_obj.get_status_of_proposal(proposalinfo_cancel.get('ProposalID')), 2)
        assert_code(pip_obj.get_status_of_proposal(proposalinfo_version.get('ProposalID')), 6)

        result = pip_obj.submitVersion(pip_obj.node.node_id, str(time.time()), pip_obj.cfg.version5, 1,
                                       pip_obj.node.staking_address, transaction_cfg=pip_obj.cfg.transaction_cfg)
        log.info('Submit version proposal result : {}'.format(result))
        assert_code(result, 0)
        proposalinfo_version = pip_obj.get_effect_proposal_info_of_vote(pip_obj.cfg.version_proposal)
        log.info('Version proposal information : {}'.format(proposalinfo_version))
        wait_block_number(pip_obj.node, proposalinfo_version.get('EndVotingBlock'))

        result = pip_obj.submitVersion(pip_obj.node.node_id, str(time.time()), pip_obj.cfg.version5, 1,
                                       pip_obj.node.staking_address, transaction_cfg=pip_obj.cfg.transaction_cfg)
        log.info('Submit version proposal result : {}'.format(result))
        assert_code(result, 0)


class TestPIPVerify():
    @pytest.mark.P0
    def test_VP_PIP_001_003_TP_PI_001_003_CP_PI_001_003_CP_PI_001_003(self, no_vp_proposal):
        pip_obj = no_vp_proposal
        pip_id_text = str(time.time())
        result = pip_obj.submitText(pip_obj.node.node_id, pip_id_text, pip_obj.node.staking_address,
                                    transaction_cfg=pip_obj.cfg.transaction_cfg)
        log.info('Submit text proposal result : {}'.format(result))
        assert_code(result, 0)

        result = pip_obj.submitVersion(pip_obj.node.node_id, pip_id_text, pip_obj.cfg.version5, 1, pip_obj.node.staking_address,
                                       transaction_cfg=pip_obj.cfg.transaction_cfg)
        log.info('Same PIPID, submit version proposal result : {}'.format(result))
        assert_code(result, 302008)

        result = pip_obj.submitText(pip_obj.node.node_id, pip_id_text, pip_obj.node.staking_address,
                                    transaction_cfg=pip_obj.cfg.transaction_cfg)
        log.info('Same PIPID, submit text proposal result : {}'.format(result))
        assert_code(result, 302008)

        result = pip_obj.submitParam(pip_obj.node.node_id, pip_id_text, 'slashing', 'slashBlocksReward', '889',
                                     pip_obj.node.staking_address, transaction_cfg=pip_obj.cfg.transaction_cfg)
        log.info('Same PIPID, submit param proposal result : {}'.format(result))
        assert_code(result, 302008)

        result = pip_obj.submitVersion(pip_obj.node.node_id, str(time.time()), pip_obj.cfg.version5, 3,
                                       pip_obj.node.staking_address, transaction_cfg=pip_obj.cfg.transaction_cfg)
        log.info('Differ PIPID, submit version proposal result : {}'.format(result))
        assert_code(result, 0)

        proposalinfo_version = pip_obj.get_effect_proposal_info_of_vote(pip_obj.cfg.version_proposal)
        log.info('Version proposal information : {}'.format(proposalinfo_version))
        pip_id_version = proposalinfo_version.get('PIPID')

        result = pip_obj.submitCancel(pip_obj.node.node_id, pip_id_text, 1, proposalinfo_version.get('ProposalID'),
                                      pip_obj.node.staking_address, transaction_cfg=pip_obj.cfg.transaction_cfg)
        log.info('Same PIPID, submit cancel proposal result : {}'.format(result))
        assert_code(result, 302008)

        result = pip_obj.submitCancel(pip_obj.node.node_id, pip_id_version, 1, proposalinfo_version.get('ProposalID'),
                                      pip_obj.node.staking_address, transaction_cfg=pip_obj.cfg.transaction_cfg)
        log.info('Same PIPID, submit cancel proposal result : {}'.format(result))
        assert_code(result, 302008)

        result = pip_obj.submitCancel(pip_obj.node.node_id, str(time.time()), 1, proposalinfo_version.get('ProposalID'),
                                      pip_obj.node.staking_address, transaction_cfg=pip_obj.cfg.transaction_cfg)
        log.info('Differ PIPID, submit cancel proposal result : {}'.format(result))
        assert_code(result, 0)

        proposalinfo_cancel = pip_obj.get_effect_proposal_info_of_vote(pip_obj.cfg.cancel_proposal)
        log.info('Version proposal information : {}'.format(proposalinfo_cancel))
        pip_id_cancel = proposalinfo_cancel.get('PIPID')

        result = pip_obj.submitText(pip_obj.node.node_id, pip_id_version, pip_obj.node.staking_address,
                                    transaction_cfg=pip_obj.cfg.transaction_cfg)
        log.info('Same PIPID, submit text proposal result : {}'.format(result))
        assert_code(result, 302008)

        result = pip_obj.submitText(pip_obj.node.node_id, pip_id_cancel, pip_obj.node.staking_address,
                                    transaction_cfg=pip_obj.cfg.transaction_cfg)
        log.info('Same PIPID, submit text proposal result : {}'.format(result))
        assert_code(result, 302008)

        wait_block_number(pip_obj.node, proposalinfo_version.get('EndVotingBlock'))
        result = pip_obj.submitVersion(pip_obj.node.node_id, pip_id_version, pip_obj.cfg.version5, 1,
                                       pip_obj.node.staking_address, transaction_cfg=pip_obj.cfg.transaction_cfg)
        log.info('Same PIPID, submit version proposal result : {}'.format(result))
        assert_code(result, 0)

        proposalinfo_version = pip_obj.get_effect_proposal_info_of_vote(pip_obj.cfg.version_proposal)
        log.info('Version proposal information : {}'.format(proposalinfo_version))
        wait_block_number(pip_obj.node, proposalinfo_version.get('EndVotingBlock'))
        result = pip_obj.submitParam(pip_obj.node.node_id, pip_id_version, 'slashing', 'slashBlocksReward', '889',
                                     pip_obj.node.staking_address, transaction_cfg=pip_obj.cfg.transaction_cfg)
        log.info('Differ PIPID, submit param proposal result : {}'.format(result))
        assert_code(result, 0)

        proposalinfo_param = pip_obj.get_effect_proposal_info_of_vote(pip_obj.cfg.param_proposal)
        log.info('Param proposal information : {}'.format(proposalinfo_param))

        result = pip_obj.submitCancel(pip_obj.node.node_id, pip_id_cancel, 1, proposalinfo_param.get('ProposalID'),
                                      pip_obj.node.staking_address, transaction_cfg=pip_obj.cfg.transaction_cfg)
        log.info('Same PIPID, submit cancel proposal result : {}'.format(result))
        assert_code(result, 0)

        proposalinfo_cancel = pip_obj.get_effect_proposal_info_of_vote(pip_obj.cfg.cancel_proposal)
        log.info('Cancel proposal information : {}'.format(proposalinfo_cancel))

        wait_block_number(pip_obj.node, proposalinfo_cancel.get('EndVotingBlock'))
        result = pip_obj.submitText(pip_obj.node.node_id, pip_id_cancel, pip_obj.node.staking_address,
                                    transaction_cfg=pip_obj.cfg.transaction_cfg)
        log.info('Same PIPID, submit text proposal result : {}'.format(result))
        assert_code(result, 0)

    @pytest.mark.P0
    def test_VP_PIP_002_TP_PI_002_CP_PI_002_CP_PI_002(self, no_vp_proposal, client_verifier_obj_list):
        pip_obj = client_verifier_obj_list[0].pip
        submitcvpandvote(client_verifier_obj_list, 1, 1, 1, 1)
        proposalinfo = pip_obj.get_effect_proposal_info_of_vote(pip_obj.cfg.cancel_proposal)
        pip_id = proposalinfo.get('PIPID')
        wait_block_number(pip_obj.node, proposalinfo.get('EndVotingBlock'))
        assert_code(pip_obj.get_status_of_proposal(proposalinfo.get('ProposalID')), 2)
        result = pip_obj.submitText(pip_obj.node.node_id, pip_id, pip_obj.node.staking_address,
                                    transaction_cfg=pip_obj.cfg.transaction_cfg)
        log.info('Same pipid, submit text proposal result : {}'.format(result))
        assert_code(result, 302008)
        result = pip_obj.submitVersion(pip_obj.node.node_id, pip_id, pip_obj.cfg.version5, 1, pip_obj.node.staking_address,
                                       transaction_cfg=pip_obj.cfg.transaction_cfg)
        log.info('Same pipid, submit text proposal result : {}'.format(result))
        assert_code(result, 302008)

        result = pip_obj.submitVersion(pip_obj.node.node_id, str(time.time()), pip_obj.cfg.version5, 2,
                                       pip_obj.node.staking_address,
                                       transaction_cfg=pip_obj.cfg.transaction_cfg)
        log.info('Differ pipid, submit version proposal result : {}'.format(result))
        assert_code(result, 0)
        proposalinfo_version = pip_obj.get_effect_proposal_info_of_vote()
        log.info('Proposal information : {}'.format(proposalinfo_version))
        result = pip_obj.submitCancel(pip_obj.node.node_id, pip_id, 1, proposalinfo_version.get('ProposalID'),
                                      pip_obj.node.staking_address, transaction_cfg=pip_obj.cfg.transaction_cfg)
        log.info('Same pipid, submit cancel proposal result : {}'.format(result))
        assert_code(result, 302008)

        wait_block_number(pip_obj.node, proposalinfo_version.get('EndVotingBlock'))
        result = pip_obj.submitParam(pip_obj.node.node_id, pip_id, 'slashing', 'slashBlocksReward', '19',
                                     pip_obj.node.staking_address, transaction_cfg=pip_obj.cfg.transaction_cfg)
        log.info('Same pipid, submit param proposal result : {}'.format(result))
        assert_code(result, 302008)

        result = pip_obj.submitParam(pip_obj.node.node_id, str(time.time()), 'slashing', 'slashBlocksReward', '19',
                                     pip_obj.node.staking_address, transaction_cfg=pip_obj.cfg.transaction_cfg)
        log.info('Differ pipid, submit param proposal result : {}'.format(result))
        assert_code(result, 0)

        proposalinfo_param = pip_obj.get_effect_proposal_info_of_vote(pip_obj.cfg.param_proposal)
        log.info('Get param proposal information : {}'.format(proposalinfo_param))

        result = pip_obj.submitCancel(pip_obj.node.node_id, pip_id, 1, proposalinfo_param.get('ProposalID'),
                                      pip_obj.node.staking_address, transaction_cfg=pip_obj.cfg.transaction_cfg)
        log.info('Same pipid, submit cancel proposal result : {}'.format(result))
        assert_code(result, 302008)


@pytest.mark.P0
def test_CP_CR_003_CP_CR_004(submit_param):
    pip_obj = submit_param
    proposalinfo = pip_obj.get_effect_proposal_info_of_vote(pip_obj.cfg.param_proposal)
    log.info('proposalinfo: {}'.format(proposalinfo))
    endvoting_rounds = (proposalinfo.get('EndVotingBlock') - math.ceil(
        pip_obj.node.block_number / pip_obj.economic.consensus_size) * pip_obj.economic.consensus_size
        ) / pip_obj.economic.consensus_size
    result = pip_obj.submitCancel(pip_obj.node.node_id, str(time.time()), endvoting_rounds + 1,
                                  proposalinfo.get('ProposalID'), pip_obj.node.staking_address,
                                  transaction_cfg=pip_obj.cfg.transaction_cfg)
    log.info('endvoting_rounds:{}， cancel proposal result:{}'.format(endvoting_rounds, result))
    assert_code(result, 302010)

    result = pip_obj.submitCancel(pip_obj.node.node_id, str(time.time()), endvoting_rounds,
                                  proposalinfo.get('ProposalID'), pip_obj.node.staking_address,
                                  transaction_cfg=pip_obj.cfg.transaction_cfg)
    log.info('endvoting_rounds:{}， cancel proposal result:{}'.format(endvoting_rounds + 1, result))
    assert_code(result, 0)


class TestGas():
    def test_VP_GA_001(self, no_vp_proposal):
        pip_obj = no_vp_proposal
        transaction_cfg = {"gasPrice": 2100000000000000 - 1}
        try:
            pip_obj.submitVersion(pip_obj.node.node_id, str(time.time()), pip_obj.cfg.version5, 1,
                                  pip_obj.node.staking_address, transaction_cfg=transaction_cfg)
        except ValueError as e:
            assert e.args[0].get('message') == "the tx data is invalid: Invalid parameter:Gas price under the min gas price."

        transaction_cfg = {"gasPrice": 2100000000000000}
        result = pip_obj.submitVersion(pip_obj.node.node_id, str(time.time()), pip_obj.cfg.version5, 1,
                                       pip_obj.node.staking_address, transaction_cfg=transaction_cfg)
        log.info('Submit version proposal result : {}'.format(result))
        assert_code(result, 0)

    def test_CP_GA_001(self, no_vp_proposal):
        pip_obj = no_vp_proposal
        transaction_cfg = {"gasPrice": 2000000000000000 - 1}
        try:
            pip_obj.submitParam(pip_obj.node.node_id, str(time.time()), 'slashing', 'slashBlocksReward', '123',
                                pip_obj.node.staking_address, transaction_cfg=transaction_cfg)
        except ValueError as e:
            assert e.args[0].get('message') == "the tx data is invalid: Invalid parameter:Gas price under the min gas price."

        transaction_cfg = {"gasPrice": 2000000000000000}
        result = pip_obj.submitParam(pip_obj.node.node_id, str(time.time()), 'slashing', 'slashBlocksReward', '123',
                                     pip_obj.node.staking_address, transaction_cfg=transaction_cfg)
        log.info('Submit param proposal result : {}'.format(result))
        assert_code(result, 0)

    def test_TP_GA_001(self, client_verifier_obj):
        pip_obj = client_verifier_obj.pip
        transaction_cfg = {"gasPrice": 1500000000000000 - 1}
        try:
            pip_obj.submitText(pip_obj.node.node_id, str(time.time()), pip_obj.node.staking_address,
                               transaction_cfg=transaction_cfg)
        except ValueError as e:
            assert e.args[0].get('message') == "the tx data is invalid: Invalid parameter:Gas price under the min gas price."

        transaction_cfg = {"gasPrice": 1500000000000000}
        result = pip_obj.submitText(pip_obj.node.node_id, str(time.time()), pip_obj.node.staking_address,
                                    transaction_cfg=transaction_cfg)
        log.info('Submit text proposal result : {}'.format(result))
        assert_code(result, 0)

    def test_PP_GA_001(self, no_vp_proposal):
        pip_obj = no_vp_proposal
        transaction_cfg = {"gasPrice": 3000000000000000 - 1, "gas": 100000}
        try:
            result = pip_obj.submitVersion(pip_obj.node.node_id, str(time.time()), pip_obj.cfg.version5, 3,
                                           pip_obj.node.staking_address, transaction_cfg=pip_obj.cfg.transaction_cfg)
            assert_code(result, 0)
            proposalinfo = pip_obj.get_effect_proposal_info_of_vote()
            log.info('Get proposal information {}'.format(proposalinfo))
            pip_obj.submitCancel(pip_obj.node.node_id, str(time.time()), 1, proposalinfo.get('ProposalID'),
                                 pip_obj.node.staking_address, transaction_cfg=transaction_cfg)
        except ValueError as e:
            assert e.args[0].get('message') == "the tx data is invalid: Invalid parameter:Gas price under the min gas price."
        transaction_cfg = {"gasPrice": 3000000000000000}
        result = pip_obj.submitCancel(pip_obj.node.node_id, str(time.time()), 1, proposalinfo.get('ProposalID'),
                                      pip_obj.node.staking_address, transaction_cfg=transaction_cfg)
        log.info('Submit cancel proposal result : {}'.format(result))
        assert_code(result, 0)


if __name__ == '__main__':
    pytest.main(['./tests/govern/', '-s', '-q', '--alluredir', './report/report'])
