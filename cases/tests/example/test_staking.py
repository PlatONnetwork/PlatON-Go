import pytest
from tests.lib import Staking
from tests.lib import StakingConfig
from common.log import log


@pytest.fixture(scope="module", autouse=True)
def staking_obj(global_test_env):
    cfg = StakingConfig("11111", "faker", "www.baidu.com", "how much")
    node = global_test_env.get_a_normal_node()
    return Staking(global_test_env, node, cfg)


def test_staking(staking_obj):
    address, _ = staking_obj.economic.account.generate_account(staking_obj.node.web3, 10**18 * 10000000)
    log.info("Generate address:{}".format(address))
    result = staking_obj.create_staking(0, address, address)
    log.info("Staking result:{}".format(result))
    assert result["Code"] == 0
    assert result["ErrMsg"] == "ok"
