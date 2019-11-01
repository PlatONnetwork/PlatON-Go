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