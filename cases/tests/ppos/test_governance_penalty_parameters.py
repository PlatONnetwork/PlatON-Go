import json
import time

import pytest
import allure

from dacite import from_dict

from common.key import get_pub_key, mock_duplicate_sign
from common.log import log
from client_sdk_python import Web3
from decimal import Decimal

from tests.conftest import param_governance_verify
from tests.lib import EconomicConfig, Genesis, StakingConfig, Staking, check_node_in_list, assert_code, von_amount, get_governable_parameter_value


def pledge_punishment(client_con_list_obj):
    """
    :return:
    """
    # stop node
    client_con_list_obj[0].node.stop()
    # Waiting for a settlement round
    client_con_list_obj[1].economic.wait_consensus_blocknum(client_con_list_obj[1].node)
    # view verifier list
    verifier_list = client_con_list_obj[1].ppos.getVerifierList()
    log.info("verifier_list: {}".format(verifier_list))
    candidate_info = client_con_list_obj[1].ppos.getCandidateInfo(client_con_list_obj[0].node.node_id)
    assert_code(candidate_info, 0)
    log.info("Pledge node information： {}".format(candidate_info))
    return candidate_info


@pytest.mark.P1
def test_PIP_PVF_001(client_con_list_obj, reset_environment):
    """
    治理修改低出块率扣除验证人自有质押金比例投票失败
    :param client_con_list_obj:
    :return:
    """
    # Initialize environment
    client_con_list_obj[0].economic.env.deploy_all()
    # view Consensus Amount of pledge
    candidate_info1 = client_con_list_obj[0].ppos.getCandidateInfo(client_con_list_obj[0].node.node_id)
    pledge_amount1 = candidate_info1['Ret']['Released']
    # view block_reward
    block_reward, staking_reward = client_con_list_obj[0].economic.get_current_year_reward(
        client_con_list_obj[0].node)
    log.info("block_reward: {} staking_reward: {}".format(block_reward, staking_reward))
    slash_blocks = get_governable_parameter_value(client_con_list_obj[0], 'SlashBlocksReward')
    # create Parametric proposal
    param_governance_verify(client_con_list_obj[0], 'Slashing', 'SlashBlocksReward', 0, 1)
    # Verify changed parameters
    candidate_info2 = pledge_punishment(client_con_list_obj)
    pledge_amount2 = candidate_info2['Ret']['Released']
    punishment_amonut = int(Decimal(str(block_reward)) * Decimal(str(slash_blocks)))
    assert pledge_amount2 == pledge_amount1 - punishment_amonut, "ErrMsg:Consensus Amount of pledge {}".format(pledge_amount2)


@pytest.mark.P1
def test_PIP_PVF_002(client_con_list_obj, reset_environment):
    """
    治理修改低出块率扣除验证人自有质押金比例处于未生效期
    :param client_con_list_obj:
    :return:
    """
    # view Consensus Amount of pledge
    candidate_info1 = client_con_list_obj[0].ppos.getCandidateInfo(client_con_list_obj[0].node.node_id)
    pledge_amount1 = candidate_info1['Ret']['Released']
    # view block_reward
    block_reward, staking_reward = client_con_list_obj[0].economic.get_current_year_reward(
        client_con_list_obj[0].node)
    log.info("block_reward: {} staking_reward: {}".format(block_reward, staking_reward))
    slash_blocks = get_governable_parameter_value(client_con_list_obj[0], 'SlashBlocksReward')
    # create Parametric proposal
    param_governance_verify(client_con_list_obj[0], 'Slashing', 'SlashBlocksReward', 0)
    # Verify changed parameters
    candidate_info2 = pledge_punishment(client_con_list_obj)
    pledge_amount2 = candidate_info2['Ret']['Released']
    punishment_amonut = int(Decimal(str(block_reward)) * Decimal(str(slash_blocks)))
    assert pledge_amount2 == pledge_amount1 - punishment_amonut, "ErrMsg:Consensus Amount of pledge {}".format(
        pledge_amount2)


@pytest.mark.P1
def test_PIP_PVF_003(client_con_list_obj, reset_environment):
    """
    治理修改低出块率扣除验证人自有质押金比例处于已生效期
    :param client_con_list_obj:
    :param reset_environment:
    :return:
    """
    # view Consensus Amount of pledge
    candidate_info1 = client_con_list_obj[0].ppos.getCandidateInfo(client_con_list_obj[0].node.node_id)
    pledge_amount1 = candidate_info1['Ret']['Released']
    # view block_reward
    block_reward, staking_reward = client_con_list_obj[0].economic.get_current_year_reward(
        client_con_list_obj[0].node)
    log.info("block_reward: {} staking_reward: {}".format(block_reward, staking_reward))
    # Get governable parameters
    slash_blocks1 = get_governable_parameter_value(client_con_list_obj[0], 'SlashBlocksReward')
    # create Parametric proposal
    param_governance_verify(client_con_list_obj[0], 'Slashing', 'SlashBlocksReward', 0)
    # wait settlement block
    client_con_list_obj[1].economic.get_settlement_switchpoint(1)
    # Get governable parameters
    slash_blocks2 = get_governable_parameter_value(client_con_list_obj[0], 'SlashBlocksReward')
    assert slash_blocks1 != slash_blocks2,"ErrMsg:Change parameters {}".format(slash_blocks2)
    # Verify changed parameters
    candidate_info2 = pledge_punishment(client_con_list_obj)
    pledge_amount2 = candidate_info2['Ret']['Released']
    punishment_amonut = int(Decimal(str(block_reward)) * Decimal(str(slash_blocks2)))
    assert pledge_amount2 == pledge_amount1 - punishment_amonut, "ErrMsg:Consensus Amount of pledge {}".format(
        pledge_amount2)