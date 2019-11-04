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

@pytest.fixture()
def candidate_has_proposal(global_test_env, client_candidate_obj, client_verifier_obj, client_list_obj):
    pip_obj = client_verifier_obj.pip
    if pip_obj.is_exist_effective_proposal():
        proposalinfo = pip_obj.get_effect_proposal_info_of_vote()
        log.info('升级提案信息为{}'.format(proposalinfo))
        if proposalinfo.get('EndVotingBlock') - pip_obj.node.block_number < 2 * pip_obj.economic.consensus_size:
            global_test_env.deploy_all()
            normal_node_obj_list = global_test_env.normal_node_list
            for normal_node_obj in normal_node_obj_list:
                client_obj = get_client_obj(normal_node_obj.node_id, client_list_obj)
                address, _ = pip_obj.economic.account.generate_account(pip_obj.node.web3, 10 ** 18 * 10000000)
                log.info('对节点{}进行质押操作'.format(normal_node_obj.node_id))
                result = client_obj.staking.create_staking(0, address, address)
                log.info('节点{}质押结果为{}'.format(normal_node_obj.node_id, result))
                assert result.get('Code') == 0
            pip_obj.economic.wait_settlement_blocknum(pip_obj.node)
            node_id_list = pip_obj.get_candidate_list_not_verifier()
            if not node_id_list:
                raise Exception('获取候选人失败')
            client_candidate_obj = get_client_obj(node_id_list[0], client_list_obj)
        else:
            return client_candidate_obj.pip
    result = pip_obj.submitVersion(pip_obj.node.node_id, str(time.time()), pip_obj.cfg.version5, 5,
                          pip_obj.node.staking_address, transaction_cfg=pip_obj.cfg.transaction_cfg)
    log.info('发起升级提案结果为{}'.format(result))
    assert result.get('Code') == 0
    return client_candidate_obj.pip