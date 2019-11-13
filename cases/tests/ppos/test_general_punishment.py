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


def information_before_slash_blocks(client_obj, node):
    # view Consensus Amount of pledge
    candidate_info1 = client_obj.ppos.getCandidateInfo(node.node_id)
    pledge_amount1 = candidate_info1['Ret']['Released']
    # view block_reward
    log.info("block: {}".format(node.eth.blockNumber))
    block_reward, staking_reward = client_obj.economic.get_current_year_reward(node)
    log.info("block_reward: {} staking_reward: {}".format(block_reward, staking_reward))
    # Get governable parameters
    slash_blocks1 = get_governable_parameter_value(client_obj, 'SlashBlocksReward')
    return pledge_amount1, block_reward, slash_blocks1


@pytest.mark.P0
def test_VP_GPFV_003(client_new_node_obj_list):
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
    pledge_amount1, block_reward, slash_blocks = information_before_slash_blocks(client1, node)
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
def test_VP_GPFV_004(client_new_node_obj_list):
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
    pledge_amount1, block_reward, slash_blocks = information_before_slash_blocks(client1, node)
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
def test_VP_GPFV_005(client_new_node_obj_list):
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
    amount = von_amount(economic.create_staking_limit, 2)
    plan = [{'Epoch': 1, 'Amount': amount}]
    result = client1.restricting.createRestrictingPlan(address1, plan, address1)
    assert_code(result, 0)
    # create Restricting Plan
    amount = von_amount(economic.create_staking_limit, 1)
    plan = [{'Epoch': 1, 'Amount': amount}]
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
    pledge_amount1, block_reward, slash_blocks = information_before_slash_blocks(client1, node)
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