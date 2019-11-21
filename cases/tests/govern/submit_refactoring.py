import pytest

from common.log import log
from dacite import from_dict
from tests.lib import Genesis
import time
import math


def assert_code(result, code):
    assert result.get('Code') == code, "状态码错误，预期状态码：{}，实际状态码:{}".format(code, result.get("Code"))


def voting_round_deploy(pip_obj, pip_env, param):
    genesis = from_dict(data_class=Genesis, data=pip_env.genesis_config)
    genesis.economicModel.gov.versionProposalVote_DurationSeconds = 2 * pip_obj.economic.consensus_size + param
    genesis.economicModel.gov.textProposalVote_DurationSeconds = 5 * pip_obj.economic.consensus_size + param
    pip_env.set_genesis(genesis.to_dict())
    pip_env.deploy_all()


def voting_round_assert(pip_obj, param):
    result = pip_obj.submitVersion(pip_obj.node.node_id, str(time.time()), pip_obj.cfg.version5, 3,
                                   pip_obj.node.staking_address,
                                   transaction_cfg=pip_obj.cfg.transaction_cfg)
    log.info('投票共识轮数为3，发起升级提案结果为{}'.format(result))
    assert_code(result, 302010)

    result = pip_obj.submitVersion(pip_obj.node.node_id, str(time.time()), pip_obj.cfg.version5, 0,
                                   pip_obj.node.staking_address,
                                   transaction_cfg=pip_obj.cfg.transaction_cfg)
    log.info('投票共识轮数为0，发起升级提案结果为{}'.format(result))
    assert_code(result, 302009)
    result = pip_obj.submitVersion(pip_obj.node.node_id, str(time.time()), pip_obj.cfg.version5, 2,
                                   pip_obj.node.staking_address,
                                   transaction_cfg=pip_obj.cfg.transaction_cfg)
    log.info('投票共识轮数为2，发起升级提案结果为{}'.format(result))
    assert_code(result, 0)
    result = pip_obj.submitText(pip_obj.node.node_id, str(time.time()), pip_obj.node.staking_address,
                                transaction_cfg=pip_obj.cfg.transaction_cfg)
    log.info('发起文本提案结果为{}'.format(result))
    assert_code(result, 0)
    proposalinfo = pip_obj.get_effect_proposal_info_of_vote(pip_obj.cfg.text_proposal)
    log.info('获取文本提案信息{}'.format(proposalinfo))
    endvotingblock_count = math.ceil(proposalinfo.get('SubmitBlock') / pip_obj.economic.consensus_size + param
                                     ) * pip_obj.economic.consensus_size - 20
    log.info('计算投票截止块高为{},接口返回投票截止块高{}'.format(endvotingblock_count, proposalinfo.get('EndVotingBlock')))
    assert int(endvotingblock_count) == proposalinfo.get('EndVotingBlock')


class TestEndVotingRounds:
    def test_VP_CR_001_VP_CR_002_VP_CR_007_TP_TE_002(self, pip_env, client_verifier_obj):
        '''
        投票周期配置成N个共识周期对应的秒数+1
        :param pip_env:
        :param pip_obj:
        :return:
        '''
        pip_obj = client_verifier_obj.pip
        voting_round_deploy(pip_obj, pip_env, 1)
        voting_round_assert(pip_obj, 5)

    def test_VP_CR_003_VP_CR_004_VP_CR_007_TP_TE_003(self, pip_env, client_verifier_obj):
        pip_obj = client_verifier_obj.pip
        voting_round_deploy(pip_obj, pip_env, -1)
        voting_round_assert(pip_obj, 4)

    @pytest.mark.compatibility
    def test_VP_CR_005_VP_CR_006_TP_TE_001(self, pip_env, client_verifier_obj):
        pip_obj = client_verifier_obj.pip
        voting_round_deploy(pip_obj, pip_env, 0)
        voting_round_assert(pip_obj, 5)
        prosalinfo = pip_obj.get_effect_proposal_info_of_vote(1)
        log.info('text proposal info: {}'.format(prosalinfo))
        result = pip_obj.submitCancel(pip_obj.node.node_id, str(time.time()), 1, prosalinfo.get('ProposalID'),
                                      pip_obj.node.staking_address, transaction_cfg=pip_obj.cfg.transaction_cfg)
        log.info('submit cancel result: {}'.format(result))  # why not assert
