from common.log import log
from dacite import from_dict
from tests.lib import Genesis
import pytest
from tests.lib.utils import get_client_obj, get_pledge_list, upload_platon, wait_block_number
import time, math

def test_VP_SU_001(submit_version):
    pip_obj = submit_version
    proposalinfo = pip_obj.get_effect_proposal_info_of_vote()
    log.info('获取升级提案信息为{}'.format(proposalinfo))
    endvotingblock_count = math.ceil(proposalinfo.get('SubmitBlock') / pip_obj.economic.consensus_size + 5
                                     ) * pip_obj.economic.consensus_size - 20
    log.info('计算投票截止块高为{},接口返回投票截止块高{}'.format(endvotingblock_count,
                                               proposalinfo.get('EndVotingBlock')))
    assert int(endvotingblock_count) == proposalinfo.get('EndVotingBlock')
    assert int(endvotingblock_count) + 21 == proposalinfo.get('ActiveBlock')

@pytest.mark.P0
def test_CP_SU_001_CP_UN_001(submit_cancel):
    pip_obj = submit_cancel
    proposalinfo = pip_obj.get_effect_proposal_info_of_vote(4)
    log.info('获取取消提案信息为{}'.format(proposalinfo))
    endvotingblock_count = math.ceil(proposalinfo.get('SubmitBlock') / pip_obj.economic.consensus_size + 4
                                     ) * pip_obj.economic.consensus_size - 20
    log.info('计算投票截止块高为{},接口返回投票截止块高{}'.format(endvotingblock_count,
                                               proposalinfo.get('EndVotingBlock')))
    assert int(endvotingblock_count) == proposalinfo.get('EndVotingBlock')
    pip_obj.submitCancel(pip_obj.node.node_id, str(time.time()), 1, proposalinfo.get('ProposalID'),
                         pip_obj.node.staking_address, transaction_cfg=pip_obj.cfg.transaction_cfg)

def test_VP_VE_001_to_VP_VE_004(no_version_proposal):
    pip_obj_tmp = no_version_proposal
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

def test_VP_WA_001(no_version_proposal):
    pip_obj_tmp = no_version_proposal
    address, _ = pip_obj_tmp.economic.account.generate_account(pip_obj_tmp.node.web3, 10**18 * 10000000)
    result = pip_obj_tmp.submitVersion(pip_obj_tmp.node.node_id, str(time.time()), pip_obj_tmp.cfg.version5, 1,
                                       address, transaction_cfg=pip_obj_tmp.cfg.transaction_cfg)
    log.info('发起升级提案结果为{}'.format(result))
    assert result.get('Code') == 302021