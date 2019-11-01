import pytest
from common.log import log
from copy import copy
import time
from tests.lib import get_client_obj

@pytest.fixture(scope="class")
def pip_env(global_test_env):
    cfg_copy = copy(global_test_env.cfg)
    yield global_test_env
    # global_test_env.set_cfg(cfg_copy)
    # global_test_env.deploy_all()

@pytest.fixture()
def no_version_proposal(global_test_env, client_verifier_obj):
    pip_obj = client_verifier_obj.pip
    if pip_obj.is_exist_effective_proposal():
        log.info('存在有效升级提案，重新启链')
        global_test_env.deploy_all()
    return pip_obj

@pytest.fixture()
def submit_version(no_version_proposal):
    pip_obj = no_version_proposal
    result = pip_obj.submitVersion(pip_obj.node.node_id, str(time.time()), pip_obj.cfg.version5, 5,
                                   pip_obj.node.staking_address,
                                   transaction_cfg=pip_obj.cfg.transaction_cfg)
    log.info('submit version result:'.format(result))
    assert result.get('Code') == 0
    return pip_obj

@pytest.fixture()
def submit_cancel(submit_version):
    pip_obj = submit_version
    propolsalinfo = pip_obj.get_effect_proposal_info_of_vote()
    log.info('获取处于投票期的升级提案信息{}'.format(propolsalinfo))
    result = pip_obj.submitCancel(pip_obj.node.node_id, str(time.time()), 4, propolsalinfo.get('ProposalID'),
                                  pip_obj.node.staking_address, transaction_cfg=pip_obj.cfg.transaction_cfg)
    log.info('发起取消提案结果为{}'.format(result))
    assert result.get('Code') == 0
    return pip_obj

@pytest.fixture()
def submit_text(client_verifier_obj):
    pip_obj = client_verifier_obj.pip
    result = pip_obj.submitText(pip_obj.node.node_id, str(time.time()), pip_obj.node.staking_address,
                                transaction_cfg=pip_obj.cfg.transaction_cfg)
    log.info('submit text result:'.format(result))
    assert result.get('Code') == 0
    return pip_obj

@pytest.fixture()
def new_node_has_proposal(global_test_env, client_new_node_obj, client_verifier_obj, client_noconsensus_obj):
    pip_obj = client_verifier_obj.pip
    if pip_obj.is_exist_effective_proposal():
        proposalinfo = pip_obj.get_effect_proposal_info_of_vote()
        log.info('升级提案信息为{}'.format(proposalinfo))
        if proposalinfo.get('EndVotingBlock') - pip_obj.node.block_number < 2 * pip_obj.economic.consensus_size:
            global_test_env.deploy_all()
            result = pip_obj.submitVersion(pip_obj.node.node_id, str(time.time()), pip_obj.cfg.version5, 5,
                                           pip_obj.node.staking_address,
                                           transaction_cfg=pip_obj.cfg.transaction_cfg)
            assert result.get('Code') == 0
            return client_noconsensus_obj.pip
        else:
            return client_new_node_obj.pip
    result = pip_obj.submitVersion(pip_obj.node.node_id, str(time.time()), pip_obj.cfg.version5, 5,
                          pip_obj.node.staking_address, transaction_cfg=pip_obj.cfg.transaction_cfg)
    assert result.get('Code') == 0
    return client_new_node_obj.pip