import json
import time

import pytest
import allure

from dacite import from_dict

from common.key import get_pub_key, mock_duplicate_sign
from common.log import log
from client_sdk_python import Web3
from decimal import Decimal

from tests.conftest import param_governance_verify, param_governance_verify_before_endblock
from tests.lib import EconomicConfig, Genesis, StakingConfig, Staking, check_node_in_list, assert_code, von_amount, \
    get_governable_parameter_value


def pledge_punishment(client_con_list_obj):
    """
    :return:
    """
    log.info("Current block height: {}".format(client_con_list_obj[1].node.eth.blockNumber))
    # stop node
    client_con_list_obj[0].node.stop()
    # Waiting for a settlement round
    client_con_list_obj[1].economic.wait_consensus_blocknum(client_con_list_obj[1].node, 2)
    log.info("Current block height: {}".format(client_con_list_obj[1].node.eth.blockNumber))
    # view verifier list
    verifier_list = client_con_list_obj[1].ppos.getVerifierList()
    log.info("verifier_list: {}".format(verifier_list))
    candidate_info = client_con_list_obj[1].ppos.getCandidateInfo(client_con_list_obj[0].node.node_id)
    log.info("Pledge node information： {}".format(candidate_info))
    return candidate_info


def Information_before_governance(client_obj):
    # view Consensus Amount of pledge
    candidate_info1 = client_obj.ppos.getCandidateInfo(client_obj.node.node_id)
    pledge_amount1 = candidate_info1['Ret']['Released']
    # view block_reward
    log.info("block: {}".format(client_obj.node.eth.blockNumber))
    block_reward, staking_reward = client_obj.economic.get_current_year_reward(
        client_obj.node)
    log.info("block_reward: {} staking_reward: {}".format(block_reward, staking_reward))
    slash_blocks = get_governable_parameter_value(client_obj, 'SlashBlocksReward')


@pytest.mark.P1
def test_PIP_PVF_001(client_con_list_obj, reset_environment):
    """
    治理修改低出块率扣除验证人自有质押金比例投票失败
    :param client_con_list_obj:
    :return:
    """
    # Initialize environment
    client_con_list_obj[0].economic.env.deploy_all()
    time.sleep(3)
    # view Consensus Amount of pledge
    candidate_info1 = client_con_list_obj[0].ppos.getCandidateInfo(client_con_list_obj[0].node.node_id)
    pledge_amount1 = candidate_info1['Ret']['Released']
    # view block_reward
    log.info("block: {}".format(client_con_list_obj[0].node.eth.blockNumber))
    block_reward, staking_reward = client_con_list_obj[0].economic.get_current_year_reward(
        client_con_list_obj[0].node)
    log.info("block_reward: {} staking_reward: {}".format(block_reward, staking_reward))
    slash_blocks = get_governable_parameter_value(client_con_list_obj[0], 'Slashing', 'SlashBlocksReward')
    # create Parametric proposal
    param_governance_verify(client_con_list_obj[0], 'Slashing', 'SlashBlocksReward', '0', False)
    # Verify changed parameters
    candidate_info2 = pledge_punishment(client_con_list_obj)
    pledge_amount2 = candidate_info2['Ret']['Released']
    punishment_amonut = int(Decimal(str(block_reward)) * Decimal(str(slash_blocks)))
    if punishment_amonut < pledge_amount1:
        assert pledge_amount2 == pledge_amount1 - punishment_amonut, "ErrMsg:Consensus Amount of pledge {}".format(
            pledge_amount2)
    else:
        assert pledge_amount2 == 0, "ErrMsg:Consensus Amount of pledge {}".format(pledge_amount2)


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
    slash_blocks = get_governable_parameter_value(client_con_list_obj[0], 'Slashing', 'SlashBlocksReward')
    # create Parametric proposal
    End_voting_block = param_governance_verify_before_endblock(client_con_list_obj[0], 'Slashing', 'SlashBlocksReward',
                                                               '0')
    # Verify changed parameters
    candidate_info2 = pledge_punishment(client_con_list_obj)
    pledge_amount2 = candidate_info2['Ret']['Released']
    punishment_amonut = int(Decimal(str(block_reward)) * Decimal(str(slash_blocks)))
    if punishment_amonut < pledge_amount1:
        assert pledge_amount2 == pledge_amount1 - punishment_amonut, "ErrMsg:Consensus Amount of pledge {}".format(
            pledge_amount2)
    else:
        assert pledge_amount2 == 0, "ErrMsg:Consensus Amount of pledge {}".format(pledge_amount2)


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
    slash_blocks1 = get_governable_parameter_value(client_con_list_obj[0], 'Slashing', 'SlashBlocksReward')
    # create Parametric proposal
    param_governance_verify(client_con_list_obj[0], 'Slashing', 'SlashBlocksReward', '0')
    log.info("Current block height: {}".format(client_con_list_obj[0].node.eth.blockNumber))
    # Get governable parameters
    slash_blocks2 = get_governable_parameter_value(client_con_list_obj[0], 'Slashing', 'SlashBlocksReward')
    assert slash_blocks2 == '0', "ErrMsg:Change parameters {}".format(slash_blocks2)
    # Verify changed parameters
    candidate_info2 = pledge_punishment(client_con_list_obj)
    pledge_amount2 = candidate_info2['Ret']['Released']
    punishment_amonut = int(Decimal(str(block_reward)) * Decimal(str(slash_blocks2)))
    if punishment_amonut < pledge_amount1:
        assert pledge_amount2 == pledge_amount1 - punishment_amonut, "ErrMsg:Consensus Amount of pledge {}".format(
            pledge_amount2)
    else:
        assert pledge_amount2 == 0, "ErrMsg:Consensus Amount of pledge {}".format(pledge_amount2)


@pytest.mark.P1
def test_PIP_PVF_004(client_con_list_obj, client_new_node_obj_list, reset_environment):
    """

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
    slash_blocks1 = get_governable_parameter_value(client_con_list_obj[0], 'Slashing', 'SlashBlocksReward')
    # create Parametric proposal
    param_governance_verify(client_con_list_obj[0], 'Slashing', 'SlashBlocksReward', '60100')
    log.info("Current block height: {}".format(client_con_list_obj[0].node.eth.blockNumber))
    # Get governable parameters
    slash_blocks2 = get_governable_parameter_value(client_con_list_obj[0], 'Slashing', 'SlashBlocksReward')
    assert slash_blocks2 == '60100', "ErrMsg:Change parameters {}".format(slash_blocks2)
    # create account
    address, _ = client_con_list_obj[0].economic.account.generate_account(client_con_list_obj[0].node.web3,
                                                                          client_con_list_obj[
                                                                              0].economic.create_staking_limit * 2)
    # create staking
    result = client_new_node_obj_list[0].staking.create_staking(0, address, address)
    assert_code(result, 0)
    # wait settlement block
    log.info(client_new_node_obj_list[0].node)
    client_new_node_obj_list[0].economic.wait_settlement_blocknum(client_new_node_obj_list[0].node)
    for i in range(4):
        result = check_node_in_list(client_con_list_obj[0].node.node_id, client_con_list_obj[0].ppos.getValidatorList)
        log.info("Current node in consensus list status：{}".format(result))
        if not result:
            # Verify changed parameters
            candidate_info2 = pledge_punishment(client_new_node_obj_list)
            pledge_amount2 = candidate_info2['Ret']['Released']
            punishment_amonut = int(Decimal(str(block_reward)) * Decimal(str(slash_blocks2)))
            if punishment_amonut < pledge_amount1:
                assert pledge_amount2 == pledge_amount1 - punishment_amonut, "ErrMsg:Consensus Amount of pledge {}".format(
                    pledge_amount2)
            else:
                assert pledge_amount2 == 0, "ErrMsg:Consensus Amount of pledge {}".format(pledge_amount2)
        else:
            # wait consensus block
            client_con_list_obj[0].economic.wait_consensus_blocknum(client_con_list_obj[0].node)


@pytest.mark.P1
def test_PIP_PVF_005(client_con_list_obj, client_noc_list_obj, reset_environment):
    """
    治理修改低出块率扣除验证人自有质押金比例扣除区块奖励块数60100-锁仓金额质押
    :param client_con_list_obj:
    :param client_noc_list_obj:
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
    slash_blocks1 = get_governable_parameter_value(client_con_list_obj[0], 'Slashing', 'SlashBlocksReward')
    # create Parametric proposal
    param_governance_verify(client_con_list_obj[0], 'Slashing', 'SlashBlocksReward', '60100')
    # wait settlement block
    client_con_list_obj[1].economic.get_settlement_switchpoint(1)
    # Get governable parameters
    slash_blocks2 = get_governable_parameter_value(client_con_list_obj[0], 'Slashing', 'SlashBlocksReward')
    assert slash_blocks2 == '60100', "ErrMsg:Change parameters {}".format(slash_blocks2)
    # create account
    address, _ = client_con_list_obj[0].economic.account.generate_account(client_con_list_obj[0].node.web3,
                                                                          client_con_list_obj[
                                                                              0].economic.create_staking_limit * 2)
    address, _ = client_con_list_obj[0].economic.account.generate_account(client_con_list_obj[0].node.web3, client_con_list_obj[0].node.web3.toWei(1000, 'ether'))
    # Create restricting plan
    plan = [{'Epoch': 1, 'Amount': client_noc_list_obj[0].economic.create_staking_limit}]
    result = client_noc_list_obj[0].restricting.createRestrictingPlan(address, plan, address)
    assert_code(result, 0)
    # create staking
    result = client_noc_list_obj[0].staking.create_staking(0, address, address)
    assert_code(result, 0)
    # wait settlement block
    client_noc_list_obj[0].economic.wait_settlement_blocknum(client_noc_list_obj[0].node)
    for i in range(4):
        result = check_node_in_list(client_con_list_obj[0].node.node_id, client_con_list_obj[0].ppos.getValidatorList)
        log.info("Current node in consensus list status：{}".format(result))
        if not result:
            # Verify changed parameters
            candidate_info2 = pledge_punishment(client_noc_list_obj)
            pledge_amount2 = candidate_info2['Ret']['Released']
            punishment_amonut = int(Decimal(str(block_reward)) * Decimal(str(slash_blocks2)))
            if punishment_amonut < pledge_amount1:
                assert pledge_amount2 == pledge_amount1 - punishment_amonut, "ErrMsg:Consensus Amount of pledge {}".format(
                    pledge_amount2)
            else:
                assert pledge_amount2 == 0, "ErrMsg:Consensus Amount of pledge {}".format(pledge_amount2)
        else:
            # wait consensus block
            client_con_list_obj[0].economic.wait_consensus_blocknum(client_con_list_obj[0].node)

