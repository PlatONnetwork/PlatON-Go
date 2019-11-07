from common.log import log
from dacite import from_dict
from tests.lib import Genesis
import pytest
from tests.lib.utils import get_pledge_list, upload_platon, wait_block_number, assert_code, get_governable_parameter_value
from tests.lib.client import get_client_obj
import time, math

def test_VP_SU_001(submit_version):
    pip_obj = submit_version
    proposalinfo = pip_obj.get_effect_proposal_info_of_vote()
    log.info('获取升级提案信息为{}'.format(proposalinfo))
    endvotingblock_count = math.ceil(proposalinfo.get('SubmitBlock') / pip_obj.economic.consensus_size + 5
                                     ) * pip_obj.economic.consensus_size - 20
    log.info('Calculated endvoting block{},interface returned endvoting block{}'.format(endvotingblock_count,
                                               proposalinfo.get('EndVotingBlock')))
    assert int(endvotingblock_count) == proposalinfo.get('EndVotingBlock')
    assert int(endvotingblock_count) + 21 == proposalinfo.get('ActiveBlock')

@pytest.mark.P0
def test_CP_SU_001_CP_UN_001(submit_cancel):
    pip_obj = submit_cancel
    proposalinfo = pip_obj.get_effect_proposal_info_of_vote(pip_obj.cfg.cancel_proposal)
    log.info('cancel proposalinfo : {}'.format(proposalinfo))
    endvotingblock_count = math.ceil(proposalinfo.get('SubmitBlock') / pip_obj.economic.consensus_size + 4
                                     ) * pip_obj.economic.consensus_size - 20
    log.info('Calculated endvoting block{},interface returned endvoting block{}'.format(endvotingblock_count,
                                               proposalinfo.get('EndVotingBlock')))
    assert int(endvotingblock_count) == proposalinfo.get('EndVotingBlock')
    result = pip_obj.submitCancel(pip_obj.node.node_id, str(time.time()), 1, proposalinfo.get('ProposalID'),
                         pip_obj.node.staking_address, transaction_cfg=pip_obj.cfg.transaction_cfg)
    log.info('Submit cancel proposal result : {}'.format(result))
    assert_code(result, 302014)

class TestsubmitCP():
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


    def test_CP_SU_002_CP_UN_002(self, submit_cancel_param):
        pip_obj = submit_cancel_param
        proposalinfo = pip_obj.get_effect_proposal_info_of_vote(pip_obj.cfg.cancel_proposal)
        log.info('cancel proposalinfo : {}'.format(proposalinfo))
        endvotingblock_count = math.ceil(proposalinfo.get('SubmitBlock') / pip_obj.economic.consensus_size + 3
                                         ) * pip_obj.economic.consensus_size - 20
        log.info('Calculated endvoting block{},interface returned endvoting block{}'.format(endvotingblock_count,
                                                   proposalinfo.get('EndVotingBlock')))
        assert int(endvotingblock_count) == proposalinfo.get('EndVotingBlock')
        result = pip_obj.submitCancel(pip_obj.node.node_id, str(time.time()), 1, proposalinfo.get('ProposalID'),
                             pip_obj.node.staking_address, transaction_cfg=pip_obj.cfg.transaction_cfg)
        log.info('Submit cancel proposal result : {}'.format(result))
        assert_code(result, 302014)

@pytest.mark.P0
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
    result = pip_obj.submitParam(pip_obj.node.node_id, str(time.time()), 'Slashing', 'SlashBlocksReward', '0',
                                 pip_obj.node.staking_address, transaction_cfg=pip_obj.cfg.transaction_cfg)
    log.info('There is a voting param proposal,submit param proposal result : {}'.format(result))
    assert_code(result, 302032)

    result = pip_obj.submitVersion(pip_obj.node.node_id, str(time.time()), pip_obj.cfg.version5, 1,
                                   pip_obj.node.staking_address, transaction_cfg=pip_obj.cfg.transaction_cfg)
    log.info('There is a voting param proposal,submit version proposal result : {}'.format(result))
    assert_code(result, 302032)

def test_VP_VE_001_to_VP_VE_004(no_vp_proposal):
    pip_obj_tmp = no_vp_proposal
    result = pip_obj_tmp.submitVersion(pip_obj_tmp.node.node_id, str(time.time()), pip_obj_tmp.cfg.version1, 1,
                                       pip_obj_tmp.node.staking_address, transaction_cfg=pip_obj_tmp.cfg.transaction_cfg)
    assert result.get("Code") == 302011

    result = pip_obj_tmp.submitVersion(pip_obj_tmp.node.node_id, str(time.time()), pip_obj_tmp.cfg.version2, 1,
                                       pip_obj_tmp.node.staking_address, transaction_cfg=pip_obj_tmp.cfg.transaction_cfg)
    assert result.get("Code") == 302011

    result = pip_obj_tmp.submitVersion(pip_obj_tmp.node.node_id, str(time.time()), pip_obj_tmp.cfg.version3, 1,
                                       pip_obj_tmp.node.staking_address, transaction_cfg=pip_obj_tmp.cfg.transaction_cfg)
    assert result.get("Code") == 302011

    result = pip_obj_tmp.submitVersion(pip_obj_tmp.node.node_id, str(time.time()), pip_obj_tmp.chain_version, 1,
                                   pip_obj_tmp.node.staking_address, transaction_cfg=pip_obj_tmp.cfg.transaction_cfg)
    assert result.get("Code") == 302011

def test_VP_WA_001(no_vp_proposal):
    pip_obj_tmp = no_vp_proposal
    address, _ = pip_obj_tmp.economic.account.generate_account(pip_obj_tmp.node.web3, 10**18 * 10000000)
    result = pip_obj_tmp.submitVersion(pip_obj_tmp.node.node_id, str(time.time()), pip_obj_tmp.cfg.version5, 1,
                                       address, transaction_cfg=pip_obj_tmp.cfg.transaction_cfg)
    log.info('发起升级提案结果为{}'.format(result))
    assert result.get('Code') == 302021

def test_TP_WA_001(client_verifier_obj):
    pip_obj = client_verifier_obj
    address, _ = pip_obj.economic.account.generate_account(pip_obj.node.web3, 10**18 * 10000000)
    result = pip_obj.submitVersion(pip_obj.node.node_id, str(time.time()), pip_obj.cfg.version5, 1,
                                       address, transaction_cfg=pip_obj.cfg.transaction_cfg)
    log.info('发起升级提案结果为{}'.format(result))
    assert result.get('Code') == 302021

def test_TP_UN_001(submit_text):
    pip_obj = submit_text
    result = pip_obj.submitText(pip_obj.node.node_id, str(time.time()), pip_obj.node.staking_address,
                                transaction_cfg=pip_obj.cfg.transaction_cfg)
    log.info('存在处于投票期文本提案，再次发起文本提案结果为{}'.format(result))
    assert result.get('Code') == 0

def test_VP_SU_001_VP_UN_001(submit_version):
    pip_obj = submit_version
    result = pip_obj.submitVersion(pip_obj.node.node_id, str(time.time()), pip_obj.cfg.version5, 1,
                                   pip_obj.node.staking_address, transaction_cfg=pip_obj.cfg.transaction_cfg)
    log.info('有处于投票期的升级提案，再次发起升级提案结果为{}'.format(result))
    assert result.get('Code') == 302012

def test_VP_UN_002_CP_ID_002_VP_UN_002(preactive_proposal_pipobj_list):
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

    result = pip_obj.submitParam(pip_obj.node.node_id, str(time.time()), 'Slashing', 'SlashBlocksReward',
                                 '84', pip_obj.node.staking_address, transaction_cfg=pip_obj.cfg.transaction_cfg)
    log.info('there is preactive version proposal, submit cancel param proposal result: {}'.format(result))
    assert_code(result, 302017)

class TestEndVotingRounds():
    def test_VP_CR_001_VP_CR_002_VP_CR_007_TP_TE_002(self, pip_env, client_verifier_obj):
        '''
        投票周期配置成N个共识周期对应的秒数+1
        :param pip_env:
        :param pip_obj:
        :return:
        '''
        pip_obj = client_verifier_obj
        genesis = from_dict(data_class=Genesis, data=pip_env.genesis_config)
        genesis.EconomicModel.Gov.VersionProposalVote_DurationSeconds = 2 * pip_obj.cfg.consensus_block + 1
        genesis.EconomicModel.Gov.TextProposalVote_DurationSeconds = 5 * pip_obj.cfg.consensus_block + 1
        pip_env.set_genesis(genesis.to_dict())
        pip_env.deploy_all()
        result = pip_obj.submitVersion(pip_obj.node.node_id, str(time.time()), pip_obj.cfg.version5, 3,
                                       pip_obj.node.staking_address,
                                       transaction_cfg=pip_obj.cfg.transaction_cfg)
        log.info('投票共识轮数为3，发起升级提案结果为{}'.format(result))
        assert result.get('Code') == 302010

        result = pip_obj.submitVersion(pip_obj.node.node_id, str(time.time()), pip_obj.cfg.version5, 0,
                                       pip_obj.node.staking_address,
                                       transaction_cfg=pip_obj.cfg.transaction_cfg)
        log.info('投票共识轮数为0，发起升级提案结果为{}'.format(result))
        assert result.get('Code') == 302009

        result = pip_obj.submitVersion(pip_obj.node.node_id, str(time.time()), pip_obj.cfg.version5, 2,
                                       pip_obj.node.staking_address,
                                       transaction_cfg=pip_obj.cfg.transaction_cfg)
        log.info('投票共识轮数为2，发起升级提案结果为{}'.format(result))
        assert result.get('Code') == 0

        result = pip_obj.submitText(pip_obj.node.node_id, str(time.time()), pip_obj.node.staking_address,
                                    transaction_cfg=pip_obj.cfg.transaction_cfg)
        log.info('发起文本提案结果为{}'.format(result))
        assert result.get('Code') == 0
        proposalinfo = pip_obj.get_effect_proposal_info_of_vote(pip_obj.cfg.text_proposal)
        log.info('获取文本提案信息{}'.format(proposalinfo))
        endvotingblock_count = math.ceil(proposalinfo.get('SubmitBlock') / pip_obj.economic.consensus_size + 5
                                         ) * pip_obj.economic.consensus_size - 20
        log.info('计算投票截止块高为{},接口返回投票截止块高{}'.format(endvotingblock_count,
                                                   proposalinfo.get('EndVotingBlock')))
        assert int(endvotingblock_count) == proposalinfo.get('EndVotingBlock')

    def test_VP_CR_003_VP_CR_004_VP_CR_007_TP_TE_003(self, pip_env, client_verifier_obj):
        pip_obj = client_verifier_obj.pip
        genesis = from_dict(data_class=Genesis, data=pip_env.genesis_config)
        genesis.EconomicModel.Gov.VersionProposalVote_DurationSeconds = 3 * pip_obj.cfg.consensus_block - 1
        genesis.EconomicModel.Gov.TextProposalVote_DurationSeconds = 5 * pip_obj.cfg.consensus_block - 1
        pip_env.set_genesis(genesis.to_dict())
        pip_env.deploy_all()
        result = pip_obj.submitVersion(pip_obj.node.node_id, str(time.time()), pip_obj.cfg.version5, 3, pip_obj.node.staking_address,
                                       transaction_cfg=pip_obj.cfg.transaction_cfg)
        log.info('投票共识轮数为3，发起升级提案结果为{}'.format(result))
        assert result.get('Code') == 302010

        result = pip_obj.submitVersion(pip_obj.node.node_id, str(time.time()), pip_obj.cfg.version5, 0, pip_obj.node.staking_address,
                                       transaction_cfg=pip_obj.cfg.transaction_cfg)
        log.info('投票共识轮数为0，发起升级提案结果为{}'.format(result))
        assert result.get('Code') == 302009

        result = pip_obj.submitVersion(pip_obj.node.node_id, str(time.time()), pip_obj.cfg.version5, 2, pip_obj.node.staking_address,
                                       transaction_cfg=pip_obj.cfg.transaction_cfg)
        log.info('投票共识轮数为2，发起升级提案结果为{}'.format(result))
        assert result.get('Code') == 0

        result = pip_obj.submitText(pip_obj.node.node_id, str(time.time()), pip_obj.node.staking_address,
                                    transaction_cfg=pip_obj.cfg.transaction_cfg)
        log.info('发起文本提案结果为{}'.format(result))
        assert result.get('Code') == 0
        proposalinfo = pip_obj.get_effect_proposal_info_of_vote(pip_obj.cfg.text_proposal)
        log.info('获取文本提案信息{}'.format(proposalinfo))
        endvotingblock_count = math.ceil(proposalinfo.get('SubmitBlock') / pip_obj.economic.consensus_size + 4
                                         ) * pip_obj.economic.consensus_size - 20
        log.info('计算投票截止块高为{},接口返回投票截止块高{}'.format(endvotingblock_count,
                                                   proposalinfo.get('EndVotingBlock')))
        assert int(endvotingblock_count) == proposalinfo.get('EndVotingBlock')

    def test_VP_CR_005_VP_CR_006_TP_TE_001(self, pip_env, client_verifier_obj):
        pip_obj = client_verifier_obj.pip
        genesis = from_dict(data_class=Genesis, data=pip_env.genesis_config)
        genesis.EconomicModel.Gov.VersionProposalVote_DurationSeconds = 3 * pip_obj.economic.consensus_size
        genesis.EconomicModel.Gov.TextProposalVote_DurationSeconds = 5 * pip_obj.economic.consensus_size
        pip_env.set_genesis(genesis.to_dict())
        pip_env.deploy_all()
        result = pip_obj.submitVersion(pip_obj.node.node_id, str(time.time()), pip_obj.cfg.version5, 4,
                                       pip_obj.node.staking_address,
                                       transaction_cfg=pip_obj.cfg.transaction_cfg)
        log.info('投票共识轮数为3，发起升级提案结果为{}'.format(result))
        assert result.get('Code') == 302010

        result = pip_obj.submitVersion(pip_obj.node.node_id, str(time.time()), pip_obj.cfg.version5, 0,
                                       pip_obj.node.staking_address,
                                       transaction_cfg=pip_obj.cfg.transaction_cfg)
        log.info('投票共识轮数为0，发起升级提案结果为{}'.format(result))
        assert result.get('Code') == 302009

        result = pip_obj.submitVersion(pip_obj.node.node_id, str(time.time()), pip_obj.cfg.version5, 3,
                                       pip_obj.node.staking_address,
                                       transaction_cfg=pip_obj.cfg.transaction_cfg)
        log.info('投票共识轮数为2，发起升级提案结果为{}'.format(result))
        assert result.get('Code') == 0

        result = pip_obj.submitText(pip_obj.node.node_id, str(time.time()), pip_obj.node.staking_address,
                                    transaction_cfg=pip_obj.cfg.transaction_cfg)
        log.info('发起文本提案结果为{}'.format(result))
        assert result.get('Code') == 0
        proposalinfo = pip_obj.get_effect_proposal_info_of_vote(pip_obj.cfg.text_proposal)
        log.info('获取文本提案信息{}'.format(proposalinfo))
        endvotingblock_count = math.ceil(proposalinfo.get('SubmitBlock') / pip_obj.economic.consensus_size + 5
                                         ) * pip_obj.economic.consensus_size - 20
        log.info('计算投票截止块高为{},接口返回投票截止块高{}'.format(endvotingblock_count,
                                                   proposalinfo.get('EndVotingBlock')))
        assert int(endvotingblock_count) == proposalinfo.get('EndVotingBlock')

        proosalinfo = pip_obj.get_effect_proposal_info_of_vote(1)
        log.info('text proposalinfo: {}'.format(proosalinfo))
        result = pip_obj.submitCancel(pip_obj.node.node_id, str(time.time()), 1, proosalinfo.get('ProposalID'),
                             pip_obj.node.staking_address, transaction_cfg=pip_obj.cfg.transaction_cfg)
        log.info('submit cancel result: {}'.format(result))

class TestNoVerifierSubmitProposal():
    def test_VP_PR_002_TP_PR_002(self, client_new_node_obj):
        pip_obj = client_new_node_obj
        address, _ = pip_obj.economic.account.generate_account(pip_obj.node.web3, 10**18 * 10000000)
        result = pip_obj.submitVersion(pip_obj.node.node_id, str(time.time()), pip_obj.cfg.version5, 1, address,
                                       transaction_cfg=pip_obj.cfg.transaction_cfg)
        log.info('新节点发起版本声明，结果为{}'.format(result))
        assert result.get('Code') == 302022

        result = pip_obj.submitText(pip_obj.node.node_id, str(time.time()), address,
                                       transaction_cfg=pip_obj.cfg.transaction_cfg)
        log.info('新节点发起版本声明，结果为{}'.format(result))
        assert result.get('Code') == 302022

    def test_VP_PR_001_TP_PR_001(self, client_candidate_obj):
        pip_obj = client_candidate_obj
        result = pip_obj.submitVersion(pip_obj.node.node_id, str(time.time()),
                                                pip_obj.cfg.version5, 1, pip_obj.node.staking_address,
                                                transaction_cfg=pip_obj.cfg.transaction_cfg)
        log.info('候选节点节点{}发起升级提案，结果为{}'.format(pip_obj.node.node_id, result))
        assert result.get('Code') == 302022

        result = pip_obj.submitText(pip_obj.node.node_id, str(time.time()), pip_obj.node.staking_address,
                                                transaction_cfg=pip_obj.cfg.transaction_cfg)
        log.info('候选节点节点{}发起文本提案，结果为{}'.format(pip_obj.node.node_id, result))
        assert result.get('Code') == 302022

    def test_VP_PR_003_VP_PR_004_TP_PR_003_TP_PR_004(self, client_verifier_obj):
        address = client_verifier_obj.node.staking_address
        result = client_verifier_obj.staking.withdrew_staking(address)
        log.info('节点{}发起退质押结果为{}'.format(client_verifier_obj.node.node_id, result))
        assert result.get("Code") == 0
        log.info(client_verifier_obj.economic.account.find_pri_key(address))
        result = client_verifier_obj.pip.submitVersion(client_verifier_obj.node.node_id, str(time.time()),
                                                       client_verifier_obj.pip.cfg.version5, 1, address,
                                              transaction_cfg=client_verifier_obj.pip.cfg.transaction_cfg)
        log.info('节点退出中，发起升级提案结果为{}'.format(result))
        assert result.get('Code') == 302020

        result = client_verifier_obj.pip.submitText(client_verifier_obj.node.node_id, str(time.time()), address,
                                           transaction_cfg=client_verifier_obj.pip.cfg.transaction_cfg)
        log.info('节点退出中，发起文本提案结果为{}'.format(result))
        assert result.get('Code') == 302020

        client_verifier_obj.economic.wait_settlement_blocknum(client_verifier_obj.node,
                                                              number=client_verifier_obj.economic.unstaking_freeze_ratio)
        result = client_verifier_obj.pip.submitVersion(client_verifier_obj.node.node_id, str(time.time()),
                                                       client_verifier_obj.pip.cfg.version5, 1, address,
                                              transaction_cfg=client_verifier_obj.pip.cfg.transaction_cfg)
        log.info('节点已退出，发起升级提案结果为{}'.format(result))
        assert result.get('Code') == 302022

        client_verifier_obj.economic.wait_settlement_blocknum(client_verifier_obj.node,
                                                              number=client_verifier_obj.economic.unstaking_freeze_ratio)
        result = client_verifier_obj.pip.submitText(client_verifier_obj.node.node_id, str(time.time()), address,
                                           transaction_cfg=client_verifier_obj.pip.cfg.transaction_cfg)
        log.info('节点已退出，发起文本提案结果为{}'.format(result))
        assert result.get('Code') == 302022

class TestSubmitCancel():
    @pytest.mark.P0
    def test_CP_WA_001(self, submit_version):
        pip_obj = submit_version
        address, _ = pip_obj.economic.account.generate_account(pip_obj.node.web3, 10**18 * 10000)
        proposalinfo = pip_obj.get_effect_proposal_info_of_vote()
        log.info('获取升级提案信息为{}'.format(proposalinfo))
        result = pip_obj.submitCancel(pip_obj.node.node_id, str(time.time()), 1, proposalinfo.get('ProposalID'), address,
                             transaction_cfg=pip_obj.cfg.transaction_cfg)
        assert result.get('Code') == 302021

    @pytest.mark.P0
    def test_CP_PR_001(self, new_node_has_proposal):
        pip_obj = new_node_has_proposal
        address, _ = pip_obj.economic.account.generate_account(pip_obj.node.web3, 10**18 * 10000)
        proposalinfo = pip_obj.get_effect_proposal_info_of_vote()
        log.info('升级提案信息为{}'.format(proposalinfo))
        result = pip_obj.submitCancel(pip_obj.node.node_id, str(time.time()), 1,
                                                    proposalinfo.get('ProposalID'),
                                               address, transaction_cfg=pip_obj.cfg.transaction_cfg)
        assert result.get('Code') == 302022

    @pytest.mark.P0
    def test_CP_PR_002(self, candidate_has_proposal):
        pip_obj = candidate_has_proposal
        proposalinfo = pip_obj.get_effect_proposal_info_of_vote()
        log.info('升级提案信息为{}'.format(proposalinfo))
        result = pip_obj.submitCancel(pip_obj.node.node_id, str(time.time()), 1,
                                                    proposalinfo.get('ProposalID'),
                                               pip_obj.node.staking_address, transaction_cfg=pip_obj.cfg.transaction_cfg)
        log.info('候选人发起升级提案结果为{}'.format(result))
        assert result.get('Code') == 302022

    @pytest.mark.P2
    def test_CP_PR_003_CP_PR_004(self, submit_version, client_list_obj):
        pip_obj = submit_version
        proposalinfo = pip_obj.get_effect_proposal_info_of_vote()
        client_obj = get_client_obj(pip_obj.node.node_id, client_list_obj)
        address = client_obj.node.staking_address
        log.info('proposalinfo: {}'.format(proposalinfo))
        result = client_obj.staking.withdrew_staking(address)
        log.info('nodeid: {} withdrewstaking result: {}'.format(client_obj.node.node_id, result))
        assert result.get("Code") == 0
        result = client_obj.pip.submitCancel(client_obj.node.node_id, str(time.time()), 1,
                                                      proposalinfo.get('ProposalID'), address,
                                              transaction_cfg=client_obj.pip.cfg.transaction_cfg)
        log.info('node exiting，cancel proposal result: {}'.format(result))
        assert result.get('Code') == 302020

        client_obj.economic.wait_settlement_blocknum(client_obj.node,
                                                     number=client_obj.economic.unstaking_freeze_ratio)
        result = client_obj.pip.submitCancel(client_obj.node.node_id, str(time.time()), 1,
                                                      proposalinfo.get('ProposalID'), address,
                                              transaction_cfg=client_obj.pip.cfg.transaction_cfg)
        log.info('exited node，cancel proposal result: {}'.format(result))
        assert result.get('Code') == 302022

    def test_CP_CR_001(self, submit_version):
        pip_obj = submit_version
        proposalinfo = pip_obj.get_effect_proposal_info_of_vote()
        log.info('proposalinfo: {}'.format(proposalinfo))
        endvoting_rounds = (math.ceil(proposalinfo.get('EndVotingBlock')/pip_obj.economic.consensus_size) - math.ceil(
            pip_obj.node.block_number/pip_obj.economic.consensus_size)) / pip_obj.economic.consensus_size
        result = pip_obj.submitCancel(pip_obj.node.node_id, str(time.time()), endvoting_rounds,
                                                      proposalinfo.get('ProposalID'), pip_obj.node.staking_address,
                                              transaction_cfg=pip_obj.cfg.transaction_cfg)
        log.info('endvoting_rounds:{}， cancel proposal result:{}'.format(endvoting_rounds, result))
        assert result.get('Code') == 302009

        result = pip_obj.submitCancel(pip_obj.node.node_id, str(time.time()), endvoting_rounds+1,
                                                      proposalinfo.get('ProposalID'), pip_obj.node.staking_address,
                                              transaction_cfg=pip_obj.cfg.transaction_cfg)
        log.info('endvoting_rounds:{}， cancel proposal result:{}'.format(endvoting_rounds+1, result))
        assert result.get('Code') == 302009

    def test_CP_ID_001(self, no_vp_proposal):
        pip_obj = no_vp_proposal
        result = pip_obj.submitCancel(pip_obj.node.node_id, str(time.time()), 1,
                                                      '0x49b83cfc4b99462f7131d14d80c73b6657237753cd1e878e8d62dc2e9f574123',
                             pip_obj.node.staking_address, transaction_cfg=pip_obj.cfg.transaction_cfg)
        log.info('cancel proposal result: {}'.format(result))
        assert result.get('Code') == 302015

class TestPP():
    def test_PP_SU_001_PP_SU_002(self, no_vp_proposal, client_list_obj):
        pip_obj = no_vp_proposal
        client_obj = get_client_obj(pip_obj.node.node_id, client_list_obj)
        result = pip_obj.submitParam(pip_obj.node.node_id, str(time.time()), 'Slashing', 'SlashBlocksReward', '',
                            pip_obj.node.staking_address, transaction_cfg=pip_obj.cfg.transaction_cfg)
        log.info('Submit param proposal result : {}'.format(result))
        assert_code(result, 3)

        result = pip_obj.submitParam(pip_obj.node.node_id, str(time.time()), 'Slashing', 'SlashBlocksReward', '1.1',
                            pip_obj.node.staking_address, transaction_cfg=pip_obj.cfg.transaction_cfg)
        log.info('Submit param proposal result : {}'.format(result))
        assert_code(result, 3)

        result = pip_obj.submitParam(pip_obj.node.node_id, str(time.time()), 'Slashing', 'SlashBlocksReward', '-1',
                            pip_obj.node.staking_address, transaction_cfg=pip_obj.cfg.transaction_cfg)
        log.info('Submit param proposal result : {}'.format(result))
        assert_code(result, 3)

        result = pip_obj.submitParam(pip_obj.node.node_id, str(time.time()), 'Slashing', 'SlashBlocksReward', '60101',
                            pip_obj.node.staking_address, transaction_cfg=pip_obj.cfg.transaction_cfg)
        log.info('Submit param proposal result : {}'.format(result))
        assert_code(result, 3)

        result = pip_obj.submitParam(pip_obj.node.node_id, str(time.time()), 'Slashing', 'SlashBlocksReward', 1,
                            pip_obj.node.staking_address, transaction_cfg=pip_obj.cfg.transaction_cfg)
        log.info('Submit param proposal result : {}'.format(result))
        assert_code(result, 3)

        result = pip_obj.submitParam(pip_obj.node.node_id, str(time.time()), 'Slashing', 'SlashBlocksReward', '1000000',
                            pip_obj.node.staking_address, transaction_cfg=pip_obj.cfg.transaction_cfg)
        log.info('Submit param proposal result : {}'.format(result))
        assert_code(result, 3)

        result = pip_obj.submitParam(pip_obj.node.node_id, str(time.time()), 'Slashing', 'SlashBlocksReward',
                                     str(get_governable_parameter_value(client_obj, 'SlashBlocksReward')),
                            pip_obj.node.staking_address, transaction_cfg=pip_obj.cfg.transaction_cfg)
        log.info('Submit param proposal result : {}'.format(result))
        assert_code(result, 302033)

        if pip_obj.economic.slash_blocks_reward != 0:
            result = pip_obj.submitParam(pip_obj.node.node_id, str(time.time()), 'Slashing', 'SlashBlocksReward', '0',
                                pip_obj.node.staking_address, transaction_cfg=pip_obj.cfg.transaction_cfg)
            log.info('Submit param proposal result : {}'.format(result))
            assert_code(result, 0)

        if pip_obj.economic.slash_blocks_reward != 60100:
            result = pip_obj.submitParam(pip_obj.node.node_id, str(time.time()), 'Slashing', 'SlashBlocksReward', '60100',
                                         pip_obj.node.staking_address, transaction_cfg=pip_obj.cfg.transaction_cfg)
            log.info('Submit param proposal result : {}'.format(result))
            assert_code(result, 0)

    def test_PP_SU_003_PP_SU_004(self, no_vp_proposal, client_list_obj):
        pip_obj = no_vp_proposal
        client_obj = get_client_obj(pip_obj.node.node_id, client_list_obj)
        result = pip_obj.submitParam(pip_obj.node.node_id, str(time.time()), 'Slashing', 'MaxEvidenceAge', '',
                            pip_obj.node.staking_address, transaction_cfg=pip_obj.cfg.transaction_cfg)
        log.info('Submit param proposal result : {}'.format(result))
        assert_code(result, 3)

        result = pip_obj.submitParam(pip_obj.node.node_id, str(time.time()), 'Slashing', 'MaxEvidenceAge', '1.1',
                            pip_obj.node.staking_address, transaction_cfg=pip_obj.cfg.transaction_cfg)
        log.info('Submit param proposal result : {}'.format(result))
        assert_code(result, 3)

        result = pip_obj.submitParam(pip_obj.node.node_id, str(time.time()), 'Slashing', 'MaxEvidenceAge', '-1',
                            pip_obj.node.staking_address, transaction_cfg=pip_obj.cfg.transaction_cfg)
        log.info('Submit param proposal result : {}'.format(result))
        assert_code(result, 3)

        result = pip_obj.submitParam(pip_obj.node.node_id, str(time.time()), 'Slashing', 'MaxEvidenceAge', '0',
                            pip_obj.node.staking_address, transaction_cfg=pip_obj.cfg.transaction_cfg)
        log.info('Submit param proposal result : {}'.format(result))
        assert_code(result, 3)

        result = pip_obj.submitParam(pip_obj.node.node_id, str(time.time()), 'Slashing', 'MaxEvidenceAge', 1,
                            pip_obj.node.staking_address, transaction_cfg=pip_obj.cfg.transaction_cfg)
        log.info('Submit param proposal result : {}'.format(result))
        assert_code(result, 3)

        result = pip_obj.submitParam(pip_obj.node.node_id, str(time.time()), 'Slashing', 'MaxEvidenceAge',
                                     str(get_governable_parameter_value(client_obj, 'MaxEvidenceAge')),
                            pip_obj.node.staking_address, transaction_cfg=pip_obj.cfg.transaction_cfg)
        log.info('Submit param proposal result : {}'.format(result))
        assert_code(result, 302033)

        result = pip_obj.submitParam(pip_obj.node.node_id, str(time.time()), 'Slashing', 'MaxEvidenceAge',
                                     str(get_governable_parameter_value(client_obj, 'UnStakeFreezeDuration')),
                            pip_obj.node.staking_address, transaction_cfg=pip_obj.cfg.transaction_cfg)
        log.info('Submit param proposal result : {}'.format(result))
        assert_code(result, 3)

        if int(get_governable_parameter_value(client_obj, 'MaxEvidenceAge')) != 1:
            result = pip_obj.submitParam(pip_obj.node.node_id, str(time.time()), 'Slashing', 'MaxEvidenceAge', '1',
                                pip_obj.node.staking_address, transaction_cfg=pip_obj.cfg.transaction_cfg)
            log.info('Submit param proposal result : {}'.format(result))
            assert_code(result, 0)

    def test_PP_SU_004(self, no_vp_proposal):
        pip_obj = no_vp_proposal
        result = pip_obj.submitParam(pip_obj.node.node_id, str(time.time()), 'Slashing', 'MaxEvidenceAge',
                                     str(pip_obj.economic.unstaking_freeze_ratio - 1),
                                     pip_obj.node.staking_address, transaction_cfg=pip_obj.cfg.transaction_cfg)
        log.info('Submit param proposal result : {}'.format(result))
        assert_code(result, 0)

    def test_PP_SU_005_PP_SU_006(self, no_vp_proposal, client_list_obj):
        pip_obj = no_vp_proposal
        client_obj = get_client_obj(pip_obj.node.node_id, client_list_obj)
        result = pip_obj.submitParam(pip_obj.node.node_id, str(time.time()), 'Slashing', 'SlashFractionDuplicateSign', '',
                            pip_obj.node.staking_address, transaction_cfg=pip_obj.cfg.transaction_cfg)
        log.info('Submit param proposal result : {}'.format(result))
        assert_code(result, 3)

        result = pip_obj.submitParam(pip_obj.node.node_id, str(time.time()), 'Slashing', 'SlashFractionDuplicateSign', '1.1',
                            pip_obj.node.staking_address, transaction_cfg=pip_obj.cfg.transaction_cfg)
        log.info('Submit param proposal result : {}'.format(result))
        assert_code(result, 3)

        result = pip_obj.submitParam(pip_obj.node.node_id, str(time.time()), 'Slashing', 'SlashFractionDuplicateSign', '-1',
                            pip_obj.node.staking_address, transaction_cfg=pip_obj.cfg.transaction_cfg)
        log.info('Submit param proposal result : {}'.format(result))
        assert_code(result, 3)

        result = pip_obj.submitParam(pip_obj.node.node_id, str(time.time()), 'Slashing', 'SlashFractionDuplicateSign', '0',
                            pip_obj.node.staking_address, transaction_cfg=pip_obj.cfg.transaction_cfg)
        log.info('Submit param proposal result : {}'.format(result))
        assert_code(result, 3)

        result = pip_obj.submitParam(pip_obj.node.node_id, str(time.time()), 'Slashing', 'SlashFractionDuplicateSign', 1,
                            pip_obj.node.staking_address, transaction_cfg=pip_obj.cfg.transaction_cfg)
        log.info('Submit param proposal result : {}'.format(result))
        assert_code(result, 3)

        result = pip_obj.submitParam(pip_obj.node.node_id, str(time.time()), 'Slashing', 'SlashFractionDuplicateSign',
                                     '10001', pip_obj.node.staking_address, transaction_cfg=pip_obj.cfg.transaction_cfg)
        log.info('Submit param proposal result : {}'.format(result))
        assert_code(result, 3)

        result = pip_obj.submitParam(pip_obj.node.node_id, str(time.time()), 'Slashing', 'SlashFractionDuplicateSign',
                                     str(get_governable_parameter_value(client_obj, 'SlashFractionDuplicateSign')),
                            pip_obj.node.staking_address, transaction_cfg=pip_obj.cfg.transaction_cfg)
        log.info('Submit param proposal result : {}'.format(result))
        assert_code(result, 302033)

        if int(get_governable_parameter_value(client_obj, 'MaxEvidenceAge')) != 1:
            result = pip_obj.submitParam(pip_obj.node.node_id, str(time.time()), 'Slashing', 'SlashFractionDuplicateSign', '1',
                                pip_obj.node.staking_address, transaction_cfg=pip_obj.cfg.transaction_cfg)
            log.info('Submit param proposal result : {}'.format(result))
            assert_code(result, 0)

        if int(get_governable_parameter_value(client_obj, 'MaxEvidenceAge')) != 10000:
            result = pip_obj.submitParam(pip_obj.node.node_id, str(time.time()), 'Slashing', 'SlashFractionDuplicateSign', '10000',
                                pip_obj.node.staking_address, transaction_cfg=pip_obj.cfg.transaction_cfg)
            log.info('Submit param proposal result : {}'.format(result))
            assert_code(result, 0)

    def test_PP_SU_007_PP_SU_008(self, no_vp_proposal, client_list_obj):
        pip_obj = no_vp_proposal
        client_obj = get_client_obj(pip_obj.node.node_id, client_list_obj)
        result = pip_obj.submitParam(pip_obj.node.node_id, str(time.time()), 'Slashing', 'DuplicateSignReportReward', '',
                            pip_obj.node.staking_address, transaction_cfg=pip_obj.cfg.transaction_cfg)
        log.info('Submit param proposal result : {}'.format(result))
        assert_code(result, 3)

        result = pip_obj.submitParam(pip_obj.node.node_id, str(time.time()), 'Slashing', 'DuplicateSignReportReward', '1.1',
                            pip_obj.node.staking_address, transaction_cfg=pip_obj.cfg.transaction_cfg)
        log.info('Submit param proposal result : {}'.format(result))
        assert_code(result, 3)

        result = pip_obj.submitParam(pip_obj.node.node_id, str(time.time()), 'Slashing', 'DuplicateSignReportReward', '-1',
                            pip_obj.node.staking_address, transaction_cfg=pip_obj.cfg.transaction_cfg)
        log.info('Submit param proposal result : {}'.format(result))
        assert_code(result, 3)

        result = pip_obj.submitParam(pip_obj.node.node_id, str(time.time()), 'Slashing', 'DuplicateSignReportReward', '0',
                            pip_obj.node.staking_address, transaction_cfg=pip_obj.cfg.transaction_cfg)
        log.info('Submit param proposal result : {}'.format(result))
        assert_code(result, 3)

        result = pip_obj.submitParam(pip_obj.node.node_id, str(time.time()), 'Slashing', 'DuplicateSignReportReward', 1,
                            pip_obj.node.staking_address, transaction_cfg=pip_obj.cfg.transaction_cfg)
        log.info('Submit param proposal result : {}'.format(result))
        assert_code(result, 3)

        result = pip_obj.submitParam(pip_obj.node.node_id, str(time.time()), 'Slashing', 'DuplicateSignReportReward',
                                     '81', pip_obj.node.staking_address, transaction_cfg=pip_obj.cfg.transaction_cfg)
        log.info('Submit param proposal result : {}'.format(result))
        assert_code(result, 3)

        result = pip_obj.submitParam(pip_obj.node.node_id, str(time.time()), 'Slashing', 'DuplicateSignReportReward',
                                     str(get_governable_parameter_value(client_obj, 'DuplicateSignReportReward')),
                            pip_obj.node.staking_address, transaction_cfg=pip_obj.cfg.transaction_cfg)
        log.info('Submit param proposal result : {}'.format(result))
        assert_code(result, 302033)

        if int(get_governable_parameter_value(client_obj, 'DuplicateSignReportReward')) != 1:
            result = pip_obj.submitParam(pip_obj.node.node_id, str(time.time()), 'Slashing', 'DuplicateSignReportReward', '1',
                                pip_obj.node.staking_address, transaction_cfg=pip_obj.cfg.transaction_cfg)
            log.info('Submit param proposal result : {}'.format(result))
            assert_code(result, 0)

        if int(get_governable_parameter_value(client_obj, 'DuplicateSignReportReward')) != 80:
            result = pip_obj.submitParam(pip_obj.node.node_id, str(time.time()), 'Slashing', 'DuplicateSignReportReward',
                                         '80', pip_obj.node.staking_address,
                                         transaction_cfg=pip_obj.cfg.transaction_cfg)
            log.info('Submit param proposal result : {}'.format(result))
            assert_code(result, 0)

    def test_PP_SU_009_PP_SU_010(self, no_vp_proposal, client_list_obj):
        pip_obj = no_vp_proposal
        client_obj = get_client_obj(pip_obj.node.node_id, client_list_obj)
        result = pip_obj.submitParam(pip_obj.node.node_id, str(time.time()), 'Staking', 'StakeThreshold', '',
                            pip_obj.node.staking_address, transaction_cfg=pip_obj.cfg.transaction_cfg)
        log.info('Submit param proposal result : {}'.format(result))
        assert_code(result, 3)

        result = pip_obj.submitParam(pip_obj.node.node_id, str(time.time()), 'Staking', 'StakeThreshold', '100000.1',
                            pip_obj.node.staking_address, transaction_cfg=pip_obj.cfg.transaction_cfg)
        log.info('Submit param proposal result : {}'.format(result))
        assert_code(result, 3)

        result = pip_obj.submitParam(pip_obj.node.node_id, str(time.time()), 'Staking', 'StakeThreshold', '-10**18 * 1000000',
                            pip_obj.node.staking_address, transaction_cfg=pip_obj.cfg.transaction_cfg)
        log.info('Submit param proposal result : {}'.format(result))
        assert_code(result, 3)

        result = pip_obj.submitParam(pip_obj.node.node_id, str(time.time()), 'Staking', 'StakeThreshold', '0',
                            pip_obj.node.staking_address, transaction_cfg=pip_obj.cfg.transaction_cfg)
        log.info('Submit param proposal result : {}'.format(result))
        assert_code(result, 3)

        result = pip_obj.submitParam(pip_obj.node.node_id, str(time.time()), 'Staking', 'StakeThreshold', 10**18 * 1000000,
                            pip_obj.node.staking_address, transaction_cfg=pip_obj.cfg.transaction_cfg)
        log.info('Submit param proposal result : {}'.format(result))
        assert_code(result, 3)

        result = pip_obj.submitParam(pip_obj.node.node_id, str(time.time()), 'Staking', 'StakeThreshold',
                                     '10**18 * 1000000 - 1', pip_obj.node.staking_address, transaction_cfg=pip_obj.cfg.transaction_cfg)
        log.info('Submit param proposal result : {}'.format(result))
        assert_code(result, 3)
        result = pip_obj.submitParam(pip_obj.node.node_id, str(time.time()), 'Staking', 'StakeThreshold',
                                     str(get_governable_parameter_value(client_obj, 'StakeThreshold')),
                            pip_obj.node.staking_address, transaction_cfg=pip_obj.cfg.transaction_cfg)
        log.info('Submit param proposal result : {}'.format(result))
        assert_code(result, 302033)

        if int(get_governable_parameter_value(client_obj, 'StakeThreshold')) != 10**18 * 1000000:
            result = pip_obj.submitParam(pip_obj.node.node_id, str(time.time()), 'Staking', 'StakeThreshold', '10**18 * 1000000',
                                pip_obj.node.staking_address, transaction_cfg=pip_obj.cfg.transaction_cfg)
            log.info('Submit param proposal result : {}'.format(result))
            assert_code(result, 0)

    def test_PP_SU_010(self, no_vp_proposal):
        pip_obj = no_vp_proposal
        result = pip_obj.submitParam(pip_obj.node.node_id, str(time.time()), 'Staking', 'StakeThreshold',
                                     '10000000000000000000000000000000000000000000000000000000000000000000000000000000000000',
                                     pip_obj.node.staking_address,
                                     transaction_cfg=pip_obj.cfg.transaction_cfg)
        log.info('Submit param proposal result : {}'.format(result))
        assert_code(result, 0)

    def test_PP_SU_011_PP_SU_012(self, no_vp_proposal, client_list_obj):
        pip_obj = no_vp_proposal
        client_obj = get_client_obj(pip_obj.node.node_id, client_list_obj)
        result = pip_obj.submitParam(pip_obj.node.node_id, str(time.time()), 'Staking', 'OperatingThreshold', '',
                            pip_obj.node.staking_address, transaction_cfg=pip_obj.cfg.transaction_cfg)
        log.info('Submit param proposal result : {}'.format(result))
        assert_code(result, 3)

        result = pip_obj.submitParam(pip_obj.node.node_id, str(time.time()), 'Staking', 'OperatingThreshold', '10**18 * 10 + 0.5',
                            pip_obj.node.staking_address, transaction_cfg=pip_obj.cfg.transaction_cfg)
        log.info('Submit param proposal result : {}'.format(result))
        assert_code(result, 3)

        result = pip_obj.submitParam(pip_obj.node.node_id, str(time.time()), 'Staking', 'OperatingThreshold', '-10**18 * 10',
                            pip_obj.node.staking_address, transaction_cfg=pip_obj.cfg.transaction_cfg)
        log.info('Submit param proposal result : {}'.format(result))
        assert_code(result, 3)

        result = pip_obj.submitParam(pip_obj.node.node_id, str(time.time()), 'Staking', 'OperatingThreshold', '0',
                            pip_obj.node.staking_address, transaction_cfg=pip_obj.cfg.transaction_cfg)
        log.info('Submit param proposal result : {}'.format(result))
        assert_code(result, 3)

        result = pip_obj.submitParam(pip_obj.node.node_id, str(time.time()), 'Staking', 'OperatingThreshold', 1,
                            pip_obj.node.staking_address, transaction_cfg=pip_obj.cfg.transaction_cfg)
        log.info('Submit param proposal result : {}'.format(result))
        assert_code(result, 3)

        result = pip_obj.submitParam(pip_obj.node.node_id, str(time.time()), 'Staking', 'OperatingThreshold',
                                     '9', pip_obj.node.staking_address, transaction_cfg=pip_obj.cfg.transaction_cfg)
        log.info('Submit param proposal result : {}'.format(result))
        assert_code(result, 3)

        result = pip_obj.submitParam(pip_obj.node.node_id, str(time.time()), 'Staking', 'OperatingThreshold',
                                     str(get_governable_parameter_value(client_obj, 'OperatingThreshold')),
                            pip_obj.node.staking_address, transaction_cfg=pip_obj.cfg.transaction_cfg)
        log.info('Submit param proposal result : {}'.format(result))
        assert_code(result, 302033)

        if int(get_governable_parameter_value(client_obj, 'OperatingThreshold')) != 10**18 * 10:
            result = pip_obj.submitParam(pip_obj.node.node_id, str(time.time()), 'Staking', 'OperatingThreshold', '10**18 * 10',
                                pip_obj.node.staking_address, transaction_cfg=pip_obj.cfg.transaction_cfg)
            log.info('Submit param proposal result : {}'.format(result))
            assert_code(result, 0)

    def test_PP_SU_012(self, no_vp_proposal):
        pip_obj = no_vp_proposal
        result = pip_obj.submitParam(pip_obj.node.node_id, str(time.time()), 'Staking', 'OperatingThreshold',
                                     '10000000000000000000000000000000000000000000000000000000000000000000000000000000000000',
                                     pip_obj.node.staking_address,
                                     transaction_cfg=pip_obj.cfg.transaction_cfg)
        log.info('Submit param proposal result : {}'.format(result))
        assert_code(result, 0)

    def test_PP_SU_013_PP_SU_014(self, no_vp_proposal, client_list_obj):
        pip_obj = no_vp_proposal
        client_obj = get_client_obj(pip_obj.node.node_id, client_list_obj)
        result = pip_obj.submitParam(pip_obj.node.node_id, str(time.time()), 'Staking', 'UnStakeFreezeDuration', '',
                            pip_obj.node.staking_address, transaction_cfg=pip_obj.cfg.transaction_cfg)
        log.info('Submit param proposal result : {}'.format(result))
        assert_code(result, 3)

        result = pip_obj.submitParam(pip_obj.node.node_id, str(time.time()), 'Staking', 'UnStakeFreezeDuration', '10.5',
                            pip_obj.node.staking_address, transaction_cfg=pip_obj.cfg.transaction_cfg)
        log.info('Submit param proposal result : {}'.format(result))
        assert_code(result, 3)

        result = pip_obj.submitParam(pip_obj.node.node_id, str(time.time()), 'Staking', 'UnStakeFreezeDuration', '-100',
                            pip_obj.node.staking_address, transaction_cfg=pip_obj.cfg.transaction_cfg)
        log.info('Submit param proposal result : {}'.format(result))
        assert_code(result, 3)

        result = pip_obj.submitParam(pip_obj.node.node_id, str(time.time()), 'Staking', 'UnStakeFreezeDuration', '0',
                            pip_obj.node.staking_address, transaction_cfg=pip_obj.cfg.transaction_cfg)
        log.info('Submit param proposal result : {}'.format(result))
        assert_code(result, 3)

        result = pip_obj.submitParam(pip_obj.node.node_id, str(time.time()), 'Staking', 'UnStakeFreezeDuration', 11,
                            pip_obj.node.staking_address, transaction_cfg=pip_obj.cfg.transaction_cfg)
        log.info('Submit param proposal result : {}'.format(result))
        assert_code(result, 3)

        result = pip_obj.submitParam(pip_obj.node.node_id, str(time.time()), 'Staking', 'UnStakeFreezeDuration',
                                     '113', pip_obj.node.staking_address, transaction_cfg=pip_obj.cfg.transaction_cfg)
        log.info('Submit param proposal result : {}'.format(result))
        assert_code(result, 3)

        result = pip_obj.submitParam(pip_obj.node.node_id, str(time.time()), 'Staking', 'UnStakeFreezeDuration',
                                     str(get_governable_parameter_value(client_obj, 'UnStakeFreezeDuration')),
                            pip_obj.node.staking_address, transaction_cfg=pip_obj.cfg.transaction_cfg)
        log.info('Submit param proposal result : {}'.format(result))
        assert_code(result, 302033)

        result = pip_obj.submitParam(pip_obj.node.node_id, str(time.time()), 'Staking', 'UnStakeFreezeDuration',
                                     str(get_governable_parameter_value(client_obj, 'MaxEvidenceAge')),
                            pip_obj.node.staking_address, transaction_cfg=pip_obj.cfg.transaction_cfg)
        log.info('Submit param proposal result : {}'.format(result))
        assert_code(result, 3)

        if int(get_governable_parameter_value(client_obj, 'UnStakeFreezeDuration')) != 112:
            result = pip_obj.submitParam(pip_obj.node.node_id, str(time.time()), 'Staking', 'UnStakeFreezeDuration', '112',
                                pip_obj.node.staking_address, transaction_cfg=pip_obj.cfg.transaction_cfg)
            log.info('Submit param proposal result : {}'.format(result))
            assert_code(result, 0)

        if int(get_governable_parameter_value(client_obj, 'UnStakeFreezeDuration')) != str(
                int(get_governable_parameter_value(client_obj, 'MaxEvidenceAge'))-1):
            result = pip_obj.submitParam(pip_obj.node.node_id, str(time.time()), 'Staking', 'UnStakeFreezeDuration',
                                         str(int(get_governable_parameter_value(client_obj, 'MaxEvidenceAge')) - 1),
                                pip_obj.node.staking_address, transaction_cfg=pip_obj.cfg.transaction_cfg)
            log.info('Submit param proposal result : {}'.format(result))
            assert_code(result, 0)


    def test_PP_SU_015_PP_SU_016(self, no_vp_proposal, client_list_obj):
        pip_obj = no_vp_proposal
        client_obj = get_client_obj(pip_obj.node.node_id, client_list_obj)
        result = pip_obj.submitParam(pip_obj.node.node_id, str(time.time()), 'Staking', 'MaxValidators', '',
                            pip_obj.node.staking_address, transaction_cfg=pip_obj.cfg.transaction_cfg)
        log.info('Submit param proposal result : {}'.format(result))
        assert_code(result, 3)

        result = pip_obj.submitParam(pip_obj.node.node_id, str(time.time()), 'Staking', 'MaxValidators', '30.5',
                            pip_obj.node.staking_address, transaction_cfg=pip_obj.cfg.transaction_cfg)
        log.info('Submit param proposal result : {}'.format(result))
        assert_code(result, 3)

        result = pip_obj.submitParam(pip_obj.node.node_id, str(time.time()), 'Staking', 'MaxValidators', '-100',
                            pip_obj.node.staking_address, transaction_cfg=pip_obj.cfg.transaction_cfg)
        log.info('Submit param proposal result : {}'.format(result))
        assert_code(result, 3)

        result = pip_obj.submitParam(pip_obj.node.node_id, str(time.time()), 'Staking', 'MaxValidators', '0',
                            pip_obj.node.staking_address, transaction_cfg=pip_obj.cfg.transaction_cfg)
        log.info('Submit param proposal result : {}'.format(result))
        assert_code(result, 3)

        result = pip_obj.submitParam(pip_obj.node.node_id, str(time.time()), 'Staking', 'MaxValidators', 25,
                            pip_obj.node.staking_address, transaction_cfg=pip_obj.cfg.transaction_cfg)
        log.info('Submit param proposal result : {}'.format(result))
        assert_code(result, 3)

        result = pip_obj.submitParam(pip_obj.node.node_id, str(time.time()), 'Staking', 'MaxValidators',
                                     '24', pip_obj.node.staking_address, transaction_cfg=pip_obj.cfg.transaction_cfg)
        log.info('Submit param proposal result : {}'.format(result))
        assert_code(result, 3)

        result = pip_obj.submitParam(pip_obj.node.node_id, str(time.time()), 'Staking', 'MaxValidators',
                                     '202', pip_obj.node.staking_address, transaction_cfg=pip_obj.cfg.transaction_cfg)
        log.info('Submit param proposal result : {}'.format(result))
        assert_code(result, 3)

        result = pip_obj.submitParam(pip_obj.node.node_id, str(time.time()), 'Staking', 'MaxValidators',
                                     str(get_governable_parameter_value(client_obj, 'MaxValidators')),
                            pip_obj.node.staking_address, transaction_cfg=pip_obj.cfg.transaction_cfg)
        log.info('Submit param proposal result : {}'.format(result))
        assert_code(result, 302033)

        if int(get_governable_parameter_value(client_obj, 'MaxValidators')) != 25:
            result = pip_obj.submitParam(pip_obj.node.node_id, str(time.time()), 'Staking', 'MaxValidators', '25',
                                pip_obj.node.staking_address, transaction_cfg=pip_obj.cfg.transaction_cfg)
            log.info('Submit param proposal result : {}'.format(result))
            assert_code(result, 0)

        if int(get_governable_parameter_value(client_obj, 'MaxValidators')) != 201:
            result = pip_obj.submitParam(pip_obj.node.node_id, str(time.time()), 'Staking', 'MaxValidators', '201',
                                pip_obj.node.staking_address, transaction_cfg=pip_obj.cfg.transaction_cfg)
            log.info('Submit param proposal result : {}'.format(result))
            assert_code(result, 0)

    def test_PP_SU_016(self, no_vp_proposal, client_list_obj):
        pip_obj = no_vp_proposal
        client_obj = get_client_obj(pip_obj.node.node_id, client_list_obj)
        if int(get_governable_parameter_value(client_obj, 'MaxValidators')) != 201:
            result = pip_obj.submitParam(pip_obj.node.node_id, str(time.time()), 'Staking', 'MaxValidators', '201',
                                pip_obj.node.staking_address, transaction_cfg=pip_obj.cfg.transaction_cfg)
            log.info('Submit param proposal result : {}'.format(result))
            assert_code(result, 0)

    def test_PP_SU_017_PP_SU_018(self, no_vp_proposal, client_list_obj):
        pip_obj = no_vp_proposal
        client_obj = get_client_obj(pip_obj.node.node_id, client_list_obj)
        result = pip_obj.submitParam(pip_obj.node.node_id, str(time.time()), 'Block', 'MaxBlockGasLimit', '',
                            pip_obj.node.staking_address, transaction_cfg=pip_obj.cfg.transaction_cfg)
        log.info('Submit param proposal result : {}'.format(result))
        assert_code(result, 3)

        result = pip_obj.submitParam(pip_obj.node.node_id, str(time.time()), 'Block', 'MaxBlockGasLimit', '21000*200 + 0.5',
                            pip_obj.node.staking_address, transaction_cfg=pip_obj.cfg.transaction_cfg)
        log.info('Submit param proposal result : {}'.format(result))
        assert_code(result, 3)

        result = pip_obj.submitParam(pip_obj.node.node_id, str(time.time()), 'Block', 'MaxBlockGasLimit', '-21000*200',
                            pip_obj.node.staking_address, transaction_cfg=pip_obj.cfg.transaction_cfg)
        log.info('Submit param proposal result : {}'.format(result))
        assert_code(result, 3)

        result = pip_obj.submitParam(pip_obj.node.node_id, str(time.time()), 'Block', 'MaxBlockGasLimit', '0',
                            pip_obj.node.staking_address, transaction_cfg=pip_obj.cfg.transaction_cfg)
        log.info('Submit param proposal result : {}'.format(result))
        assert_code(result, 3)

        result = pip_obj.submitParam(pip_obj.node.node_id, str(time.time()), 'Block', 'MaxBlockGasLimit', 21000*200,
                            pip_obj.node.staking_address, transaction_cfg=pip_obj.cfg.transaction_cfg)
        log.info('Submit param proposal result : {}'.format(result))
        assert_code(result, 3)

        result = pip_obj.submitParam(pip_obj.node.node_id, str(time.time()), 'Block', 'MaxBlockGasLimit',
                                     '21000*200 - 1', pip_obj.node.staking_address, transaction_cfg=pip_obj.cfg.transaction_cfg)
        log.info('Submit param proposal result : {}'.format(result))
        assert_code(result, 3)

        result = pip_obj.submitParam(pip_obj.node.node_id, str(time.time()), 'Block', 'MaxBlockGasLimit',
                                     str(get_governable_parameter_value(client_obj, 'MaxBlockGasLimit')),
                            pip_obj.node.staking_address, transaction_cfg=pip_obj.cfg.transaction_cfg)
        log.info('Submit param proposal result : {}'.format(result))
        assert_code(result, 302033)

        if int(get_governable_parameter_value(client_obj, 'MaxBlockGasLimit')) != 21000*200:
            result = pip_obj.submitParam(pip_obj.node.node_id, str(time.time()), 'Block', 'MaxBlockGasLimit', '21000*200',
                                pip_obj.node.staking_address, transaction_cfg=pip_obj.cfg.transaction_cfg)
            log.info('Submit param proposal result : {}'.format(result))
            assert_code(result, 0)

    def test_PP_SU_018(self, no_vp_proposal):
        pip_obj = no_vp_proposal
        result = pip_obj.submitParam(pip_obj.node.node_id, str(time.time()), 'Block', 'MaxBlockGasLimit',
                                     '10000000000000000000000000000000000000000000000000000000000000000000000000000000000000',
                                     pip_obj.node.staking_address,
                                     transaction_cfg=pip_obj.cfg.transaction_cfg)
        log.info('Submit param proposal result : {}'.format(result))
        assert_code(result, 0)

    def test_PP_SU_019(self, no_vp_proposal):
        pip_obj = no_vp_proposal
        result = pip_obj.submitParam(pip_obj.node.node_id, str(time.time()), 'Block', '', '1',
                                     pip_obj.node.staking_address,
                                     transaction_cfg=pip_obj.cfg.transaction_cfg)
        log.info('Submit param proposal result : {}'.format(result))
        assert_code(result, 302033)

        result = pip_obj.submitParam(pip_obj.node.node_id, str(time.time()), 'Block', 'UnStakeFreezeDuration', '100',
                                     pip_obj.node.staking_address,
                                     transaction_cfg=pip_obj.cfg.transaction_cfg)
        log.info('Submit param proposal result : {}'.format(result))
        assert_code(result, 302033)

        result = pip_obj.submitParam(pip_obj.node.node_id, str(time.time()), 'Slashing', 'slashBlocksReward', '100',
                                     pip_obj.node.staking_address,
                                     transaction_cfg=pip_obj.cfg.transaction_cfg)
        log.info('Submit param proposal result : {}'.format(result))
        assert_code(result, 302033)

        result = pip_obj.submitParam(pip_obj.node.node_id, str(time.time()), 'Slashing', 'Slash BlocksReward', '100',
                                     pip_obj.node.staking_address,
                                     transaction_cfg=pip_obj.cfg.transaction_cfg)
        log.info('Submit param proposal result : {}'.format(result))
        assert_code(result, 302033)

        result = pip_obj.submitParam(pip_obj.node.node_id, str(time.time()), 'Slashing', 'SlashBlocks./,.Reward', '100',
                                     pip_obj.node.staking_address,
                                     transaction_cfg=pip_obj.cfg.transaction_cfg)
        log.info('Submit param proposal result : {}'.format(result))
        assert_code(result, 302033)


    def test_PP_SU_020(self, no_vp_proposal):
        pip_obj = no_vp_proposal
        result = pip_obj.submitParam(pip_obj.node.node_id, str(time.time()), '', 'SlashBlocksReward', '100',
                                     pip_obj.node.staking_address,
                                     transaction_cfg=pip_obj.cfg.transaction_cfg)
        log.info('Submit param proposal result : {}'.format(result))
        assert_code(result, 302033)

        result = pip_obj.submitParam(pip_obj.node.node_id, str(time.time()), 'slashing', 'SlashBlocksReward', '100',
                                     pip_obj.node.staking_address,
                                     transaction_cfg=pip_obj.cfg.transaction_cfg)
        log.info('Submit param proposal result : {}'.format(result))
        assert_code(result, 302033)

        result = pip_obj.submitParam(pip_obj.node.node_id, str(time.time()), 'SLashing', 'slashBlocksReward', '100',
                                     pip_obj.node.staking_address,
                                     transaction_cfg=pip_obj.cfg.transaction_cfg)
        log.info('Submit param proposal result : {}'.format(result))
        assert_code(result, 302033)

        result = pip_obj.submitParam(pip_obj.node.node_id, str(time.time()), 'Slashing123', 'Slash BlocksReward', '100',
                                     pip_obj.node.staking_address,
                                     transaction_cfg=pip_obj.cfg.transaction_cfg)
        log.info('Submit param proposal result : {}'.format(result))
        assert_code(result, 302033)

        result = pip_obj.submitParam(pip_obj.node.node_id, str(time.time()), 'S lashing', 'SlashBlocks./,.Reward', '100',
                                     pip_obj.node.staking_address,
                                     transaction_cfg=pip_obj.cfg.transaction_cfg)
        log.info('Submit param proposal result : {}'.format(result))
        assert_code(result, 302033)

        result = pip_obj.submitParam(pip_obj.node.node_id, str(time.time()), 'S.,.lashing', 'SlashBlocks./,.Reward', '100',
                                     pip_obj.node.staking_address,
                                     transaction_cfg=pip_obj.cfg.transaction_cfg)
        log.info('Submit param proposal result : {}'.format(result))
        assert_code(result, 302033)


class TestSubmitPPAbnormal():
    @pytest.mark.P0
    def test_PP_UN_002(self, submit_version):
        pip_obj = submit_version
        result = pip_obj.submitParam(pip_obj.node.node_id, str(time.time()), 'Slashing', 'SlashBlocksReward', '99',
                                     pip_obj.node.staking_address, transaction_cfg=pip_obj.cfg.transaction_cfg)
        log.info('There is voting version proposal, submit a param proposal : {}'.format(result))
        assert_code(result, 302012)

    @pytest.mark.P0
    def test_PP_PR_002(self, client_new_node_obj):
        pip_obj = client_new_node_obj.pip
        address, _ = pip_obj.economic.account.generate_account(pip_obj.node.web3, 10**18 * 10000)
        result = pip_obj.submitParam(pip_obj.node.node_id, str(time.time()), 'Slashing', 'SlashBlocksReward', '88',
                                     address, transaction_cfg=pip_obj.cfg.transaction_cfg)
        log.info('new node submit param proposal result : {}'.format(result))
        assert result.get('Code') == 302022

    @pytest.mark.P0
    def test_PP_PR_001(self, client_candidate_obj):
        pip_obj = client_candidate_obj
        result = pip_obj.submitParam(pip_obj.node.node_id, str(time.time()), 'Slashing', 'SlashBlocksReward', '87',
                                     pip_obj.node.staking_address, transaction_cfg=pip_obj.cfg.transaction_cfg)
        log.info('candidate submit param proposal result :{}'.format(result))
        assert result.get('Code') == 302022

    @pytest.mark.P2
    def test_PP_PR_003_PP_PR_004(self, no_vp_proposal, client_list_obj):
        pip_obj = no_vp_proposal
        client_obj = get_client_obj(pip_obj.node.node_id, client_list_obj)
        address = pip_obj.node.staking_address
        result = client_obj.staking.withdrew_staking(address)
        log.info('nodeid: {} withdrewstaking result: {}'.format(client_obj.node.node_id, result))
        assert_code(result, 0)
        result = pip_obj.submitParam(pip_obj.node.node_id, str(time.time()), 'Slashing', 'SlashBlocksReward',
                                     '86', address, transaction_cfg=pip_obj.cfg.transaction_cfg)
        log.info('node exiting，param proposal result: {}'.format(result))
        assert_code(result, 302020)

        client_obj.economic.wait_settlement_blocknum(client_obj.node,
                                                     number=client_obj.economic.unstaking_freeze_ratio)
        result = pip_obj.submitParam(pip_obj.node.node_id, str(time.time()), 'Slashing', 'SlashBlocksReward',
                                     '86', address, transaction_cfg=pip_obj.cfg.transaction_cfg)
        log.info('exited node，cancel proposal result: {}'.format(result))
        assert_code(result, 302022)

    def test_PP_WA_001(self, no_vp_proposal):
        pip_obj = no_vp_proposal
        address, _ = pip_obj.economic.account.generate_account(pip_obj.node.web3, 10**18 * 10000)
        result = pip_obj.submitParam(pip_obj.node.node_id, str(time.time()), 'Slashing', 'SlashBlocksReward', '87',
                                     address, transaction_cfg=pip_obj.cfg.transaction_cfg)
        log.info('candidate submit param proposal result :{}'.format(result))
        assert result.get('Code') == 302021

