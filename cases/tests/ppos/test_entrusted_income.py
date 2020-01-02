import time
import pytest
import allure
import rlp
from client_sdk_python.utils.transactions import send_obj_transaction
from dacite import from_dict
from common.key import get_pub_key, mock_duplicate_sign
from common.log import log
from client_sdk_python import Web3
from decimal import Decimal
from tests.lib import EconomicConfig, Genesis, StakingConfig, Staking, check_node_in_list, assert_code, von_amount, \
    get_governable_parameter_value, get_the_dynamic_parameter_gas_fee

def create_free_staking_node(client):
    """
    创建一个自由质押节点
    :param client:
    :return:
    """
    economic = client.economic
    node = client.node
    staking_address, _ = economic.account.generate_account(node.web3, von_amount(economic.create_staking_limit, 2))
    benifit_address, _ = economic.account.generate_account(node.web3)
    result = client.staking.create_staking(benifit_address, staking_address, reward_per=1000)
    assert_code(result, 0)


def create_restricting_staking_node(client, epoch, amount):
    """
    创建一个锁仓质押节点
    :param client:
    :return:
    """
    economic = client.economic
    node = client.node
    # create restricting plan
    address, _ = economic.account.generate_account(node.web3, von_amount(economic.create_staking_limit, 2))
    restricting_address, _ = economic.account.generate_account(node.web3, node.web3.toWei(1000, 'ether'))
    plan = [{'Epoch': epoch, 'Amount': node.web3.toWei(amount, 'ether')}]
    result = client.restricting.createRestrictingPlan(restricting_address, plan, address)
    assert_code(result, 0)
    benifit_address, _ = economic.account.generate_account(node.web3)
    result = client.staking.create_staking(benifit_address, restricting_address, reward_per=1000)
    assert_code(result, 0)


def create_free_stakings_node(clients):
    """
    创建多个自由质押节点
    :param clients:
    :return:
    """
    first_client = clients[0]
    second_client = clients[1]
    first_economic = first_client.economic
    first_node = first_client.node
    second_economic = second_client.economic
    second_node = second_client.node
    first_staking_address, _ = first_economic.account.generate_account(first_node.web3, von_amount(first_economic.create_staking_limit, 2))
    first_benifit_address, _ = first_economic.account.generate_account(first_node.web3)
    result = first_client.staking.create_staking(first_benifit_address, first_staking_address, reward_per=1000)
    second_staking_address, _ = second_economic.account.generate_account(second_node.web3, von_amount(second_economic.create_staking_limit, 2))
    second_benifit_address, _ = second_economic.account.generate_account(second_node.web3)
    result = second_client.staking.create_staking(second_benifit_address, second_staking_address, reward_per=2000)
    assert_code(result, 0)


def get_dividend_information(client, address):
    """
    获取分红信息
    :param client:
    :return:
    """
    result = client.ppos.getCandidateInfo(client.node.node_id)
    blocknum = result['Ret']['StakingBlockNum']
    result = client.ppos.getDelegateInfo(blocknum, address, client.node.node_id)
    info = result['Ret']['Data']
    delegate_epoch = info['DelegateEpoch']
    cumulative_income = info['CumulativeIncome']
    return delegate_epoch, cumulative_income


def get_delegate_relevant_amount_and_epoch(client):
    result = client.ppos.getCandidateInfo(client.node.node_id)
    log.info('Current pledged node pledge information：{}'.format(result))
    last_delegate_epoch = result['Ret']['DelegateEpoch']
    delegate_total = result['Ret']['DelegateTotal']
    delegate_total_hes = result['Ret']['DelegateTotalHes']
    delegate_reward_total = result['Ret']['DelegateRewardTotal']
    return last_delegate_epoch, delegate_total, delegate_total_hes, delegate_reward_total


@pytest.mark.P0
def test_EI_BC_001(client_new_node):
    """
    自由金额首次委托，验证待领取的委托收益（单节点）
    :param client_new_node:
    :return:
    """
    client = client_new_node
    economic = client.economic
    node = client.node
    address, _ = economic.account.generate_account(node.web3, 1000)
    log.info("Create delegate account：{}".format(address))
    create_free_staking_node(client)
    log.info("Create pledge node id :{}".format(node.node_id))

    # initiate a commission
    result = client.delegate.delegate(0, address)
    assert_code(result, 0)
    delegate_epoch, cumulative_income = get_dividend_information(client, address)
    assert delegate_epoch == 1, "ErrMsg: Last time delegate epoch {}".format(delegate_epoch)
    assert cumulative_income == 0, "ErrMsg: Last time cumulative income {}".format(cumulative_income)

    # initiate redemption
    result = client.ppos.getCandidateInfo(client.node.node_id)
    blocknum = result['Ret']['StakingBlockNum']
    result = client.delegate.withdrew_delegate(blocknum, address)
    assert_code(result, 0)
    delegate_epoch, cumulative_income = get_dividend_information(client, address)
    assert delegate_epoch == 1, "ErrMsg: Last time delegate epoch {}".format(delegate_epoch)
    assert cumulative_income == 0, "ErrMsg: Last time cumulative income {}".format(cumulative_income)

    # receive dividends
    result = client.ppos.getDelegateReward(address)
    assert result['Ret'][0]['reward'] == 0, "ErrMsg: Withdraw commission award {}".format(result['Ret'][0]['reward'])
    result = client.delegate.withdraw_delegate_reward(address)
    assert_code(result, 0)
    delegate_epoch, cumulative_income = get_dividend_information(client, address)
    assert delegate_epoch == 1, "ErrMsg: Last time delegate epoch {}".format(delegate_epoch)
    assert cumulative_income == 0, "ErrMsg: Last time cumulative income {}".format(cumulative_income)


@pytest.mark.P1
def test_EI_BC_002(clients_new_node):
    """
    自由金额首次委托，验证待领取的委托收益（多节点）
    :param clients_new_node:
    :return:
    """
    first_client = clients_new_node[0]
    second_client = clients_new_node[1]
    first_economic = first_client.economic
    first_node = first_client.node
    second_economic = second_client.economic
    second_node = second_client.node
    address, _ = first_economic.account.generate_account(first_node.web3, 1000)
    log.info("Create delegate account：{}".format(address))
    create_free_stakings_node(clients_new_node)
    log.info("Create first pledge node id :{}".format(first_node.node_id))
    log.info("Create second pledge node id :{}".format(second_node.node_id))

    # initiate a commission
    result = first_client.delegate.delegate(0, address)
    assert_code(result, 0)
    delegate_epoch, cumulative_income = get_dividend_information(first_client, address)
    assert delegate_epoch == 1, "ErrMsg: Last time delegate epoch {}".format(delegate_epoch)
    assert cumulative_income == 0, "ErrMsg: Last time cumulative income {}".format(cumulative_income)

    # initiate redemption
    result = first_client.ppos.getCandidateInfo(first_client.node.node_id)
    blocknum = result['Ret']['StakingBlockNum']
    result = first_client.delegate.withdrew_delegate(blocknum, address)
    assert_code(result, 0)
    delegate_epoch, cumulative_income = get_dividend_information(first_client, address)
    assert delegate_epoch == 1, "ErrMsg: Last time delegate epoch {}".format(delegate_epoch)
    assert cumulative_income == 0, "ErrMsg: Last time cumulative income {}".format(cumulative_income)

    # receive dividends
    result = first_client.ppos.getDelegateReward(address)
    assert result['Ret'][0]['reward'] == 0, "ErrMsg: Withdraw commission award {}".format(result['Ret'][0]['reward'])
    result = first_client.delegate.withdraw_delegate_reward(address)
    assert_code(result, 0)
    delegate_epoch, cumulative_income = get_dividend_information(first_client, address)
    assert delegate_epoch == 1, "ErrMsg: Last time delegate epoch {}".format(delegate_epoch)
    assert cumulative_income == 0, "ErrMsg: Last time cumulative income {}".format(cumulative_income)


@pytest.mark.P1
def test_EI_BC_003(client_new_node):
    """
    自由金额跨周期追加委托，验证待领取的委托收益（单节点）
    :param client_new_node:
    :return:
    """
    client = client_new_node
    economic = client.economic
    node = client.node
    address, _ = economic.account.generate_account(node.web3, 1000)
    log.info("Create delegate account：{}".format(address))
    create_free_staking_node(client)
    log.info("Create pledge node id :{}".format(node.node_id))
    result = client.delegate.delegate(0, address)
    assert_code(result, 0)
    log.info("Commissioned successfully, commissioned amount：{}".format(economic.delegate_limit))
    delegate_epoch, cumulative_income = get_dividend_information(client, address)
    assert delegate_epoch == 1, "ErrMsg: Last time delegate epoch {}".format(delegate_epoch)
    assert cumulative_income == 0, "ErrMsg: Last time cumulative income {}".format(cumulative_income)
    last_delegate_epoch, delegate_total, delegate_total_hes, delegate_reward_total = get_delegate_relevant_amount_and_epoch(client)
    assert last_delegate_epoch == 1, "ErrMsg: Last time delegate epoch {}".format(last_delegate_epoch)
    assert delegate_total == 0, "The total number of effective commissioned nodes: {}".format(delegate_total)
    assert delegate_total_hes == 0, "The total number of inactive nodes commissioned: {}".format(delegate_total_hes)
    assert delegate_reward_total == 0, "Total delegated rewards currently issued by the candidate: {}".format(delegate_reward_total)
    economic.wait_settlement_blocknum(node)
    log.info("Current settlement block height：{}".format(node.eth.blockNumber))
    result = client.delegate.delegate(0, address)
    assert_code(result, 0)
    log.info("Commissioned successfully, commissioned amount：{}".format(von_amount(economic.delegate_limit,2)))
    delegate_epoch, cumulative_income = get_dividend_information(client, address)
    assert delegate_epoch == 2, "ErrMsg: Last time delegate epoch {}".format(delegate_epoch)
    assert cumulative_income == 0, "ErrMsg: Last time cumulative income {}".format(cumulative_income)
    last_delegate_epoch, delegate_total, delegate_total_hes, delegate_reward_total = get_delegate_relevant_amount_and_epoch(client)
    assert last_delegate_epoch == 2, "ErrMsg: Last time delegate epoch {}".format(last_delegate_epoch)
    assert delegate_total == economic.delegate_limit, "The total number of effective commissioned nodes: {}".format(delegate_total)
    assert delegate_total_hes == economic.delegate_limit, "The total number of inactive nodes commissioned: {}".format(delegate_total_hes)
    assert delegate_reward_total == 0, "Total delegated rewards currently issued by the candidate: {}".format(
        delegate_reward_total)
    block_reward, staking_reward = economic.get_current_year_reward(node)
    economic.wait_settlement_blocknum(node)
    log.info("Current settlement block height：{}".format(node.eth.blockNumber))
    block_number = economic.get_number_blocks_in_interval(node)
    commission_award = int(Decimal(str(block_reward)) * Decimal(str(block_number)) * Decimal(str(0.1)) + Decimal(str(staking_reward) * Decimal(str(0.1))))
    result = client.delegate.delegate(0, address)
    assert_code(result, 0)
    log.info("Commissioned successfully, commissioned amount：{}".format(von_amount(economic.delegate_limit, 2)))
    delegate_epoch, cumulative_income = get_dividend_information(client, address)
    assert delegate_epoch == 3, "ErrMsg: Last time delegate epoch {}".format(delegate_epoch)
    assert cumulative_income == commission_award, "ErrMsg: Last time cumulative income {}".format(cumulative_income)
    last_delegate_epoch, delegate_total, delegate_total_hes, delegate_reward_total = get_delegate_relevant_amount_and_epoch(
        client)
    assert last_delegate_epoch == 3, "ErrMsg: Last time delegate epoch {}".format(last_delegate_epoch)
    assert delegate_total == von_amount(economic.delegate_limit, 2), "The total number of effective commissioned nodes: {}".format(
        delegate_total)
    assert delegate_total_hes == economic.delegate_limit, "The total number of inactive nodes commissioned: {}".format(
        delegate_total_hes)
    assert delegate_reward_total == commission_award, "Total delegated rewards currently issued by the candidate: {}".format(
        delegate_reward_total)



