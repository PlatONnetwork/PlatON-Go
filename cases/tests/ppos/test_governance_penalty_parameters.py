import time

import pytest
import allure

from dacite import from_dict

from common.key import get_pub_key, mock_duplicate_sign
from common.log import log
from client_sdk_python import Web3
from decimal import Decimal
from tests.lib import EconomicConfig, Genesis, StakingConfig, Staking, check_node_in_list, assert_code, von_amount


def pledge_punishment(client_con_list_obj, punish_reward):
    """
    :return:
    """
    # view block_reward
    block_reward, staking_reward = client_con_list_obj[0].economic.get_current_year_reward(
        client_con_list_obj[0].node)
    log.info("block_reward: {} staking_reward: {}".format(block_reward, staking_reward))
    # stop node
    client_con_list_obj[0].node.stop()
    # Waiting for a settlement round
    client_con_list_obj[1].economic.wait_consensus_blocknum(client_con_list_obj[1].node)
    # view verifier list
    verifier_list = client_con_list_obj[1].ppos.getVerifierList()
    log.info("verifier_list: {}".format(verifier_list))
    punishment_amonut = block_reward * punish_reward
    return punishment_amonut

#
# @pytest.mark.P1
# def test_PIP_PVF_001(client_new_node_obj):
