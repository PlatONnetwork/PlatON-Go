import pytest
from common.log import log
import time
from tests.lib.utils import assert_code

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


def test_V0_VO_001_V0_RE_001(voting_proposal_ve_pipobj):
    pip_obj = voting_proposal_ve_pipobj
    proposalinfo = pip_obj.get_effect_proposal_info_of_vote()
    log.info('Get version proposalinfo {}'.format(proposalinfo))
    result = pip_obj.vote(pip_obj.node.node_id, proposalinfo.get('ProposalID'), pip_obj.cfg.vote_option_yeas,
                 pip_obj.node.staking_address, transaction_cfg=pip_obj.cfg.transaction_cfg)
    log.info('vote result: {}'.format(result))
    assert_code(result, 0)

    result = pip_obj.vote(pip_obj.node.node_id, proposalinfo.get('ProposalID'), pip_obj.cfg.vote_option_yeas,
                 pip_obj.node.staking_address, transaction_cfg=pip_obj.cfg.transaction_cfg)
    log.info('vote duplicated result: {}'.format(result))
    assert_code(result, 302027)





