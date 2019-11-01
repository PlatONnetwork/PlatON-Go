import pytest
from tests.lib import Pip, StakingConfig
from common.log import log
from tests.lib.client import Client
from copy import copy
from tests.lib.utils import get_pledge_list, get_client_obj, get_staking_address
import time

@pytest.fixture(scope="class")
def pip_env(global_test_env):
    cfg_copy = copy(global_test_env.cfg)
    yield global_test_env
    # global_test_env.set_cfg(cfg_copy)
    # global_test_env.deploy_all()

@pytest.fixture()
def no_version_proposal(global_test_env, pip_obj):
    if pip_obj.is_exist_effective_proposal():
        log.info('存在有效升级提案，重新启链')
        global_test_env.deploy_all()
    return pip_obj

@pytest.fixture()
def submit_version(no_version_proposal):
    pip_obj = no_version_proposal
    address = get_staking_address(pip_obj.node)
    result = pip_obj.submitVersion(pip_obj.node.node_id, str(time.time()), pip_obj.cfg.version5, 5, address,
                                   transaction_cfg=pip_obj.cfg.transaction_cfg)
    log.info('submit version result:'.format(result))
    assert result.get('Code') == 0
    return pip_obj

@pytest.fixture()
def submit_cancel(submit_version):
    pip_obj = submit_version
    address = get_staking_address(pip_obj.node)
    propolsalinfo = pip_obj.get_effect_proposal_info_of_vote()
    log.info('获取处于投票期的升级提案信息{}'.format(propolsalinfo))
    result = pip_obj.submitCancel(pip_obj.node.node_id, str(time.time()), 4, propolsalinfo.get('ProposalID'),
                                  address, transaction_cfg=pip_obj.cfg.transaction_cfg)
    log.info('发起取消提案结果为{}'.format(result))
    assert result.get('Code') == 0
    return pip_obj

@pytest.fixture()
def submit_text(pip_obj):
    address = get_staking_address(pip_obj.node)
    result = pip_obj.submitText(pip_obj.node.node_id, str(time.time()), address, transaction_cfg=pip_obj.cfg.transaction_cfg)
    log.info('submit text result:'.format(result))
    assert result.get('Code') == 0
    return pip_obj


