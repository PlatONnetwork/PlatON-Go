import time

import pytest
import allure

from dacite import from_dict

from common.key import get_pub_key, mock_duplicate_sign
from common.log import log
from client_sdk_python import Web3
from decimal import Decimal
from tests.lib import EconomicConfig, Genesis, StakingConfig, Staking, check_node_in_list, assert_code, von_amount


@pytest.mark.P0
def test_LS_FV_001(client_consensus_obj):
    """
    查看锁仓账户计划
    :param client_consensus_obj:
    :return:
    """
    # Reset environment
    client_consensus_obj.economic.env.deploy_all()
    # view Lock in contract amount
    lock_up_amount = client_consensus_obj.node.eth.getBalance(EconomicConfig.FOUNDATION_LOCKUP_ADDRESS)
    log.info("Lock in contract amount: {}".format(lock_up_amount))
    # view Lockup plan
    result = client_consensus_obj.ppos.getRestrictingInfo(EconomicConfig.INCENTIVEPOOL_ADDRESS)
    release_plans_list = result['Data']['plans']
    assert_code(result, 0)
    log.info("Lockup plan information: {}".format(result))
    # assert louck up amount
    for i in release_plans_list:
        print("a", type(release_plans_list[i]))
        print("b", EconomicConfig.release_info[i])
        assert release_plans_list[i] == EconomicConfig.release_info[
            i], "Year {} Height of block to be released: {} Release amount: {}".format(i + 1, release_plans_list[i]['blockNumber'], release_plans_list[i]['amount'])
