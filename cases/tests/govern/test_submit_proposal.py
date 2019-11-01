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