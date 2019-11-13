import time
import pytest
import allure
from dacite import from_dict
from common.key import get_pub_key, mock_duplicate_sign, generate_key
from common.log import log
from client_sdk_python import Web3
from decimal import Decimal
from tests.lib import EconomicConfig, Genesis, StakingConfig, Staking, check_node_in_list, assert_code, von_amount, \
    get_governable_parameter_value, Client, update_param_by_dict, get_param_by_dict


def get_out_block_penalty_parameters(client_obj, node, amount_type):
    # view Consensus Amount of pledge
    candidate_info = client_obj.ppos.getCandidateInfo(node.node_id)
    log.info("Pledge node information: {}".format(candidate_info))
    pledge_amount1 = candidate_info['Ret'][amount_type]
    # view block_reward
    log.info("block: {}".format(node.eth.blockNumber))
    block_reward, staking_reward = client_obj.economic.get_current_year_reward(node)
    log.info("block_reward: {} staking_reward: {}".format(block_reward, staking_reward))
    # Get governable parameters
    slash_blocks = get_governable_parameter_value(client_obj, 'SlashBlocksReward')
    return pledge_amount1, block_reward, slash_blocks


@pytest.mark.P0
def test_VP_GPFV_003(client_new_node_obj_list, reset_environment):
    """
    低出快率最高处罚标准
    :param client_new_node_obj_list:
    :return:
    """
    client1 = client_new_node_obj_list[0]
    log.info("Current connection node1: {}".format(client1))
    client2 = client_new_node_obj_list[1]
    log.info("Current connection node2: {}".format(client2))
    economic = client1.economic
    node = client1.node
    # create account
    address, _ = economic.account.generate_account(node.web3, von_amount(economic.create_staking_limit, 2))
    # create staking
    result = client1.staking.create_staking(0, address, address)
    assert_code(result, 0)
    # Wait for the settlement round to end
    economic.wait_settlement_blocknum(node)
    # get pledge amount1 and block reward
    pledge_amount1, block_reward, slash_blocks = get_out_block_penalty_parameters(client1, node, 'Released')
    log.info("Current block height: {}".format(client1.node.eth.blockNumber))
    # stop node
    client1.node.stop()
    # Waiting for a settlement round
    client2.economic.wait_consensus_blocknum(client2.node, 2)
    log.info("Current block height: {}".format(client2.node.eth.blockNumber))
    # view verifier list
    verifier_list = client2.ppos.getVerifierList()
    log.info("verifier_list: {}".format(verifier_list))
    candidate_info = client2.ppos.getCandidateInfo(client1.node.node_id)
    log.info("Pledge node information： {}".format(candidate_info))
    pledge_amount2 = candidate_info['Ret']['Released']
    punishment_amonut = int(Decimal(str(block_reward)) * Decimal(str(slash_blocks)))
    if punishment_amonut < pledge_amount1:
        assert pledge_amount2 == pledge_amount1 - punishment_amonut, "ErrMsg:Consensus Amount of pledge {}".format(
            pledge_amount2)
    else:
        assert pledge_amount2 == 0, "ErrMsg:Consensus Amount of pledge {}".format(pledge_amount2)


@pytest.mark.P0
def test_VP_GPFV_004(client_new_node_obj_list, reset_environment):
    """
    锁仓质押被惩罚最高处罚标准
    :param client_new_node_obj_list:
    :return:
    """
    client1 = client_new_node_obj_list[0]
    log.info("Current connection node1: {}".format(client1))
    client2 = client_new_node_obj_list[1]
    log.info("Current connection node2: {}".format(client2))
    economic = client1.economic
    node = client1.node
    # create account
    address, _ = economic.account.generate_account(node.web3, von_amount(economic.create_staking_limit, 2))
    # create Restricting Plan
    amount = von_amount(economic.create_staking_limit, 1)
    plan = [{'Epoch': 2, 'Amount': amount}]
    result = client1.restricting.createRestrictingPlan(address, plan, address)
    assert_code(result, 0)
    # create staking
    result = client1.staking.create_staking(1, address, address)
    assert_code(result, 0)
    # Wait for the settlement round to end
    economic.wait_settlement_blocknum(node)
    # get pledge amount1 and block reward
    pledge_amount1, block_reward, slash_blocks = get_out_block_penalty_parameters(client1, node, 'RestrictingPlan')
    log.info("Current block height: {}".format(client1.node.eth.blockNumber))
    # stop node
    client1.node.stop()
    # Waiting for a settlement round
    client2.economic.wait_consensus_blocknum(client2.node, 2)
    log.info("Current block height: {}".format(client2.node.eth.blockNumber))
    # view verifier list
    verifier_list = client2.ppos.getVerifierList()
    log.info("verifier_list: {}".format(verifier_list))
    candidate_info = client2.ppos.getCandidateInfo(client1.node.node_id)
    log.info("Pledge node information： {}".format(candidate_info))
    pledge_amount2 = candidate_info['Ret']['Released']
    punishment_amonut = int(Decimal(str(block_reward)) * Decimal(str(slash_blocks)))
    if punishment_amonut < pledge_amount1:
        assert pledge_amount2 == pledge_amount1 - punishment_amonut, "ErrMsg:Consensus Amount of pledge {}".format(
            pledge_amount2)
    else:
        assert pledge_amount2 == 0, "ErrMsg:Consensus Amount of pledge {}".format(pledge_amount2)


@pytest.mark.P2
def test_VP_GPFV_005(client_new_node_obj_list, reset_environment):
    """
    锁仓增持/委托后被惩罚
    :param client_new_node_obj_list:
    :return:
    """
    client1 = client_new_node_obj_list[0]
    log.info("Current connection node1: {}".format(client1))
    client2 = client_new_node_obj_list[1]
    log.info("Current connection node2: {}".format(client2))
    economic = client1.economic
    node = client1.node
    # create account
    address1, _ = economic.account.generate_account(node.web3, von_amount(economic.create_staking_limit, 3))
    # create account
    address2, _ = economic.account.generate_account(node.web3, von_amount(economic.create_staking_limit, 1))
    # create Restricting Plan
    amount1 = von_amount(economic.create_staking_limit, 2)
    plan = [{'Epoch': 1, 'Amount': amount1}]
    result = client1.restricting.createRestrictingPlan(address1, plan, address1)
    assert_code(result, 0)
    # create Restricting Plan
    amount2 = von_amount(economic.delegate_limit, 100)
    plan = [{'Epoch': 1, 'Amount': amount2}]
    result = client1.restricting.createRestrictingPlan(address2, plan, address2)
    assert_code(result, 0)
    # create staking
    result = client1.staking.create_staking(1, address1, address1)
    assert_code(result, 0)
    # increase staking
    result = client1.staking.increase_staking(1, address1)
    assert_code(result, 0)
    # Additional pledge
    result = client1.delegate.delegate(1, address2, amount=von_amount(economic.delegate_limit, 100))
    assert_code(result, 0)
    # Wait for the settlement round to end
    economic.wait_settlement_blocknum(node)
    # get pledge amount1 and block reward
    pledge_amount1, block_reward, slash_blocks = get_out_block_penalty_parameters(client1, node, 'RestrictingPlan')
    log.info("Current block height: {}".format(client1.node.eth.blockNumber))
    # stop node
    client1.node.stop()
    # Waiting for a settlement round
    client2.economic.wait_consensus_blocknum(client2.node, 2)
    log.info("Current block height: {}".format(client2.node.eth.blockNumber))
    # view verifier list
    verifier_list = client2.ppos.getVerifierList()
    log.info("verifier_list: {}".format(verifier_list))
    candidate_info = client2.ppos.getCandidateInfo(client1.node.node_id)
    log.info("Pledge node information： {}".format(candidate_info))
    pledge_amount2 = candidate_info['Ret']['Released']
    punishment_amonut = int(Decimal(str(block_reward)) * Decimal(str(slash_blocks)))
    if punishment_amonut < pledge_amount1:
        assert pledge_amount2 == pledge_amount1 - punishment_amonut, "ErrMsg:Consensus Amount of pledge {}".format(
            pledge_amount2)
    else:
        assert pledge_amount2 == 0, "ErrMsg:Consensus Amount of pledge {}".format(pledge_amount2)


@pytest.mark.P2
def test_VP_GPFV_006(client_new_node_obj_list, reset_environment):
    """
    自由金额增持/委托后被惩罚
    :param client_new_node_obj_list:
    :return:
    """
    client1 = client_new_node_obj_list[0]
    log.info("Current connection node1: {}".format(client1))
    client2 = client_new_node_obj_list[1]
    log.info("Current connection node2: {}".format(client2))
    economic = client1.economic
    node = client1.node
    # create account
    address1, _ = economic.account.generate_account(node.web3, von_amount(economic.create_staking_limit, 3))
    # create account
    address2, _ = economic.account.generate_account(node.web3, von_amount(economic.create_staking_limit, 1))
    # create staking
    result = client1.staking.create_staking(0, address1, address1)
    assert_code(result, 0)
    # increase staking
    result = client1.staking.increase_staking(0, address1)
    assert_code(result, 0)
    # Additional pledge
    result = client1.delegate.delegate(0, address2, amount=von_amount(economic.delegate_limit, 100))
    assert_code(result, 0)
    # Wait for the settlement round to end
    economic.wait_settlement_blocknum(node)
    # get pledge amount1 and block reward
    pledge_amount1, block_reward, slash_blocks = get_out_block_penalty_parameters(client1, node, 'Released')
    log.info("Current block height: {}".format(client1.node.eth.blockNumber))
    # stop node
    client1.node.stop()
    # Waiting for a settlement round
    client2.economic.wait_consensus_blocknum(client2.node, 2)
    log.info("Current block height: {}".format(client2.node.eth.blockNumber))
    # view verifier list
    verifier_list = client2.ppos.getVerifierList()
    log.info("verifier_list: {}".format(verifier_list))
    candidate_info = client2.ppos.getCandidateInfo(client1.node.node_id)
    log.info("Pledge node information： {}".format(candidate_info))
    pledge_amount2 = candidate_info['Ret']['Released']
    punishment_amonut = int(Decimal(str(block_reward)) * Decimal(str(slash_blocks)))
    if punishment_amonut < pledge_amount1:
        assert pledge_amount2 == pledge_amount1 - punishment_amonut, "ErrMsg:Consensus Amount of pledge {}".format(
            pledge_amount2)
    else:
        assert pledge_amount2 == 0, "ErrMsg:Consensus Amount of pledge {}".format(pledge_amount2)


@pytest.mark.P2
def test_VP_GPFV_007(client_new_node_obj_list, reset_environment):
    """
    在被惩罚前退出质押
    :param client_new_node_obj_list:
    :param reset_environment:
    :return:
    """
    client1 = client_new_node_obj_list[0]
    log.info("Current connection node1: {}".format(client1.node.node_mark))
    client2 = client_new_node_obj_list[1]
    log.info("Current connection node2: {}".format(client2.node.node_mark))
    economic = client1.economic
    node = client1.node
    # create account
    address, _ = economic.account.generate_account(node.web3, von_amount(economic.create_staking_limit, 2))
    # create staking
    result = client1.staking.create_staking(0, address, address)
    assert_code(result, 0)
    # Wait for the settlement round to end
    economic.wait_settlement_blocknum(node)
    # get pledge amount1 and block reward
    pledge_amount1, block_reward, slash_blocks = get_out_block_penalty_parameters(client1, node, 'Released')
    log.info("Current block height: {}".format(client1.node.eth.blockNumber))
    # stop node
    client1.node.stop()
    # Application for return of pledge
    result = client2.staking.withdrew_staking(address)
    assert_code(result, 0)
    # Waiting for a settlement round
    client2.economic.wait_consensus_blocknum(client2.node, 2)
    log.info("Current block height: {}".format(client2.node.eth.blockNumber))
    # view verifier list
    verifier_list = client2.ppos.getVerifierList()
    log.info("verifier_list: {}".format(verifier_list))
    candidate_info = client2.ppos.getCandidateInfo(client1.node.node_id)
    log.info("Pledge node information： {}".format(candidate_info))
    pledge_amount2 = candidate_info['Ret']['Released']
    punishment_amonut = int(Decimal(str(block_reward)) * Decimal(str(slash_blocks)))
    if punishment_amonut < pledge_amount1:
        assert pledge_amount2 == pledge_amount1 - punishment_amonut, "ErrMsg:Consensus Amount of pledge {}".format(
            pledge_amount2)
    else:
        assert pledge_amount2 == 0, "ErrMsg:Consensus Amount of pledge {}".format(pledge_amount2)


@pytest.mark.P2
def test_VP_GPFV_008(client_new_node_obj_list, reset_environment):
    """
    被处罚前进行增持
    :param client_new_node_obj_list:
    :param reset_environment:
    :return:
    """
    client1 = client_new_node_obj_list[0]
    log.info("Current connection node1: {}".format(client1.node.node_mark))
    client2 = client_new_node_obj_list[1]
    log.info("Current connection node2: {}".format(client2.node.node_mark))
    economic = client1.economic
    node = client1.node
    # create account
    address, _ = economic.account.generate_account(node.web3, von_amount(economic.create_staking_limit, 3))
    # create staking
    result = client1.staking.create_staking(0, address, address)
    assert_code(result, 0)
    # Wait for the settlement round to end
    economic.wait_settlement_blocknum(node)
    # get pledge amount1 and block reward
    pledge_amount1, block_reward, slash_blocks = get_out_block_penalty_parameters(client1, node, 'Released')
    log.info("Current block height: {}".format(client1.node.eth.blockNumber))
    # stop node
    client1.node.stop()
    # Additional pledge
    result = client2.staking.increase_staking(1, address, node_id=node.node_id, amount=economic.create_staking_limit)
    assert_code(result, 0)
    # Waiting for a settlement round
    client2.economic.wait_consensus_blocknum(client2.node, 2)
    log.info("Current block height: {}".format(client2.node.eth.blockNumber))
    # view verifier list
    verifier_list = client2.ppos.getVerifierList()
    log.info("verifier_list: {}".format(verifier_list))
    candidate_info = client2.ppos.getCandidateInfo(client1.node.node_id)
    log.info("Pledge node information： {}".format(candidate_info))
    pledge_amount2 = candidate_info['Ret']['Released']
    punishment_amonut = int(Decimal(str(block_reward)) * Decimal(str(slash_blocks)))
    if punishment_amonut < pledge_amount1:
        assert pledge_amount2 == pledge_amount1 - punishment_amonut, "ErrMsg:Consensus Amount of pledge {}".format(
            pledge_amount2)
    else:
        assert pledge_amount2 == 0, "ErrMsg:Consensus Amount of pledge {}".format(pledge_amount2)
