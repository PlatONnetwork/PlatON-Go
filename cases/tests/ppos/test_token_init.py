import time

import pytest
import allure

from dacite import from_dict
from common.log import log
from client_sdk_python import Web3
from decimal import Decimal
from tests.lib import EconomicConfig, Genesis, StakingConfig, Staking, check_node_in_list

@pytest.fixture(scope="function")
def staking_obj(global_test_env):
    cfg = StakingConfig("11111", "faker", "www.baidu.com", "how much")
    node = global_test_env.get_rand_node()
    return Staking(global_test_env, node, cfg)

@pytest.fixture(scope="class")
def staking_candidate(client_consensus_obj):
    address, _ = client_consensus_obj.economic.account.generate_account(client_consensus_obj.node.web3,
                                                                        client_consensus_obj.economic.create_staking_limit * 2)
    result = client_consensus_obj.staking.create_staking(0, address, address)
    assert result['Code'] == 0, "申请质押返回的状态：{}, {}用例失败".format(result['Code'], result['ErrMsg'])
    # 等待锁定期
    client_consensus_obj.economic.wait_settlement_blocknum(client_consensus_obj.node)
    return client_consensus_obj, address
