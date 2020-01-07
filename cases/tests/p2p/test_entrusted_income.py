import math
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


def create_staking_node(client):
    """
    创建一个自由质押节点
    :param client:
    :return:
    """
    economic = client.economic
    node = client.node
    staking_address, _ = economic.account.generate_account(node.web3, von_amount(economic.create_staking_limit, 2))
    benifit_address, _ = economic.account.generate_account(node.web3)
    result = client.staking.create_staking(0, benifit_address, staking_address, reward_per=1000)
    assert_code(result, 0)
    return staking_address


def create_stakings_node(clients):
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
    result = first_client.staking.create_staking(0, first_benifit_address, first_staking_address, reward_per=1000)
    assert_code(result, 0)
    second_staking_address, _ = second_economic.account.generate_account(second_node.web3, von_amount(second_economic.create_staking_limit, 2))
    second_benifit_address, _ = second_economic.account.generate_account(second_node.web3)
    result = second_client.staking.create_staking(0, second_benifit_address, second_staking_address, reward_per=2000)
    assert_code(result, 0)


def create_restricting_plan(client):
    economic = client.economic
    node = client.node
    # create restricting plan
    address, _ = economic.account.generate_account(node.web3, economic.create_staking_limit)
    benifit_address, _ = economic.account.generate_account(node.web3, economic.create_staking_limit)
    plan = [{'Epoch': 5, 'Amount': client.node.web3.toWei(1000, 'ether')}]
    result = client.restricting.createRestrictingPlan(benifit_address, plan, address)
    assert_code(result, 0)
    return benifit_address


def get_dividend_information(client, address):
    """
    获取分红信息
    :param client:
    :return:
    """
    result = client.ppos.getCandidateInfo(client.node.node_id)
    blocknum = result['Ret']['StakingBlockNum']
    result = client.ppos.getDelegateInfo(blocknum, address, client.node.node_id)
    log.info("Commission information：{}".format(result))
    info = result['Ret']
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
def test_EI_BC_001_005_009_015_051_057(client_new_node):
    """
    自由金额首次委托，验证待领取的委托收益（未生效期N）
    锁仓金额首次委托，验证待领取的委托收益（未生效期N）
    自由金额首次部分赎回，验证待领取的委托收益（未生效期N）
    锁仓金额首次部分赎回，验证待领取的委托收益（未生效期N）
    自由金额委托首次领取分红,验证待领取的委托收益（未生效期N）
    锁仓金额委托首次领取分红,验证待领取的委托收益（未生效期N）
    :param client_new_node:
    :return:
    """
    client = client_new_node
    economic = client.economic
    node = client.node
    delegate_address, _ = economic.account.generate_account(node.web3, von_amount(economic.delegate_limit, 10))
    log.info("Create delegate account：{}".format(delegate_address))
    create_staking_node(client)
    log.info("Create pledge node id :{}".format(node.node_id))

    # initiate a commission
    result = client.delegate.delegate(0, delegate_address, amount=von_amount(economic.delegate_limit, 5))
    assert_code(result, 0)
    delegate_epoch, cumulative_income = get_dividend_information(client, delegate_address)
    assert delegate_epoch == 1, "ErrMsg: Last time delegate epoch {}".format(delegate_epoch)
    assert cumulative_income == 0, "ErrMsg: Last time cumulative income {}".format(cumulative_income)
    economic.wait_consensus_blocknum(node)
    # initiate redemption
    result = client.ppos.getCandidateInfo(client.node.node_id)
    blocknum = result['Ret']['StakingBlockNum']
    result = client.delegate.withdrew_delegate(blocknum, delegate_address)
    assert_code(result, 0)
    delegate_epoch, cumulative_income = get_dividend_information(client, delegate_address)
    assert delegate_epoch == 1, "ErrMsg: Last time delegate epoch {}".format(delegate_epoch)
    assert cumulative_income == 0, "ErrMsg: Last time cumulative income {}".format(cumulative_income)
    time.sleep(3)

    # receive dividends
    result = client.ppos.getDelegateReward(delegate_address)
    log.info("result:{}".format(result))
    assert result['Ret'][0]['reward'] == 0, "ErrMsg: Withdraw commission award {}".format(result['Ret'][0]['reward'])
    result = client.delegate.withdraw_delegate_reward(delegate_address)
    assert_code(result, 0)
    delegate_epoch, cumulative_income = get_dividend_information(client, delegate_address)
    assert delegate_epoch == 1, "ErrMsg: Last time delegate epoch {}".format(delegate_epoch)
    assert cumulative_income == 0, "ErrMsg: Last time cumulative income {}".format(cumulative_income)


@pytest.mark.P1
@pytest.mark.parametrize('delegate_type', [0, 1])
def test_EI_BC_002_006(clients_new_node, delegate_type):
    """
    自由金额首次委托，验证待领取的委托收益（多节点）
    锁仓金额首次委托，验证待领取的委托收益（多节点）
    :param clients_new_node:
    :return:
    """
    first_client = clients_new_node[0]
    second_client = clients_new_node[1]
    first_economic = first_client.economic
    first_node = first_client.node
    second_economic = second_client.economic
    second_node = second_client.node
    if delegate_type == 0:
        address, _ = first_economic.account.generate_account(first_node.web3, von_amount(first_economic.delegate_limit, 100))
    else:
        address = create_restricting_plan(first_client)
    log.info("Create delegate account：{}".format(address))
    create_stakings_node(clients_new_node)
    log.info("Create first pledge node id :{}".format(first_node.node_id))
    log.info("Create second pledge node id :{}".format(second_node.node_id))
    first_economic.wait_consensus_blocknum(first_node)
    # initiate a commission
    delegate_amount = von_amount(first_economic.delegate_limit, 10)
    result = first_client.delegate.delegate(delegate_type, address, amount=delegate_amount)
    assert_code(result, 0)
    time.sleep(5)
    result = first_client.delegate.delegate(delegate_type, address, node_id=second_node.node_id, amount=delegate_amount)
    assert_code(result, 0)
    first_delegate_epoch, first_cumulative_income = get_dividend_information(first_client, address)
    second_delegate_epoch, second_cumulative_income = get_dividend_information(second_client, address)
    assert first_delegate_epoch == 1, "ErrMsg: Last time delegate epoch {}".format(first_delegate_epoch)
    assert first_cumulative_income == 0, "ErrMsg: Last time cumulative income {}".format(first_cumulative_income)
    assert second_delegate_epoch == 1, "ErrMsg: Last time delegate epoch {}".format(second_delegate_epoch)
    assert second_cumulative_income == 0, "ErrMsg: Last time cumulative income {}".format(second_cumulative_income)

    # initiate redemption
    result = first_client.ppos.getCandidateInfo(first_node.node_id)
    first_blocknum = result['Ret']['StakingBlockNum']
    result = first_client.delegate.withdrew_delegate(first_blocknum, address)
    assert_code(result, 0)
    result = second_client.ppos.getCandidateInfo(second_node.node_id)
    second_blocknum = result['Ret']['StakingBlockNum']
    result = first_client.delegate.withdrew_delegate(second_blocknum, address, node_id=second_node.node_id)
    assert_code(result, 0)
    first_delegate_epoch, first_cumulative_income = get_dividend_information(first_client, address)
    second_delegate_epoch, second_cumulative_income = get_dividend_information(second_client, address)
    assert first_delegate_epoch == 1, "ErrMsg: Last time delegate epoch {}".format(first_delegate_epoch)
    assert first_cumulative_income == 0, "ErrMsg: Last time cumulative income {}".format(first_cumulative_income)
    assert second_delegate_epoch == 1, "ErrMsg: Last time delegate epoch {}".format(second_delegate_epoch)
    assert second_cumulative_income == 0, "ErrMsg: Last time cumulative income {}".format(second_cumulative_income)

    # receive dividends
    result = first_client.ppos.getDelegateReward(address)
    assert result['Ret'][0]['reward'] == 0, "ErrMsg: Withdraw commission award {}".format(result['Ret'][0]['reward'])
    result = first_client.delegate.withdraw_delegate_reward(address)
    assert_code(result, 0)
    first_delegate_epoch, first_cumulative_income = get_dividend_information(first_client, address)
    second_delegate_epoch, second_cumulative_income = get_dividend_information(second_client, address)
    assert first_delegate_epoch == 1, "ErrMsg: Last time delegate epoch {}".format(first_delegate_epoch)
    assert first_cumulative_income == 0, "ErrMsg: Last time cumulative income {}".format(first_cumulative_income)
    assert second_delegate_epoch == 1, "ErrMsg: Last time delegate epoch {}".format(second_delegate_epoch)
    assert second_cumulative_income == 0, "ErrMsg: Last time cumulative income {}".format(second_cumulative_income)


@pytest.mark.P0
@pytest.mark.parametrize('delegate_type', [0, 1])
def test_EI_BC_003_007(client_new_node, delegate_type):
    """
    自由金额跨周期追加委托，验证待领取的委托收益（单节点）
    锁仓金额跨周期追加委托，验证待领取的委托收益（单节点）
    :param client_new_node:
    :return:
    """
    client = client_new_node
    economic = client.economic
    node = client.node
    if delegate_type == 0:
        address, _ = economic.account.generate_account(node.web3, von_amount(economic.delegate_limit, 100))
    else:
        address = create_restricting_plan(client)
    log.info("Create delegate account：{}".format(address))
    create_staking_node(client)
    log.info("Create pledge node id :{}".format(node.node_id))
    result = client.delegate.delegate(delegate_type, address)
    assert_code(result, 0)
    log.info("Commissioned successfully, commissioned amount：{}".format(economic.delegate_limit))
    delegate_epoch, cumulative_income = get_dividend_information(client, address)
    assert delegate_epoch == 1, "ErrMsg: Last time delegate epoch {}".format(delegate_epoch)
    assert cumulative_income == 0, "ErrMsg: Last time cumulative income {}".format(cumulative_income)
    last_delegate_epoch, delegate_total, delegate_total_hes, delegate_reward_total = get_delegate_relevant_amount_and_epoch(client)
    assert last_delegate_epoch == 1, "ErrMsg: Last time delegate epoch {}".format(last_delegate_epoch)
    assert delegate_total == 0, "The total number of effective commissioned nodes: {}".format(delegate_total)
    assert delegate_total_hes == economic.delegate_limit, "The total number of inactive nodes commissioned: {}".format(delegate_total_hes)
    assert delegate_reward_total == 0, "Total delegated rewards currently issued by the candidate: {}".format(delegate_reward_total)
    economic.wait_settlement_blocknum(node)
    log.info("Current settlement block height：{}".format(node.eth.blockNumber))
    result = client.delegate.delegate(delegate_type, address)
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
    log.info("The current node block reward: {} Pledge reward: {}".format(block_reward, staking_reward))
    economic.wait_settlement_blocknum(node)
    log.info("Current settlement block height：{}".format(node.eth.blockNumber))
    commission_award = economic.calculate_delegate_reward(node, block_reward, staking_reward)
    result = client.delegate.delegate(delegate_type, address)
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


@pytest.mark.P1
@pytest.mark.parametrize('delegate_type', [0, 1])
def test_EI_BC_004_008(clients_new_node, delegate_type):
    """
    自由金额跨周期追加委托，验证待领取的委托收益（多节点）
    锁仓金额跨周期追加委托，验证待领取的委托收益（多节点）
    :param clients_new_node:
    :return:
    """
    first_client = clients_new_node[0]
    second_client = clients_new_node[1]
    first_economic = first_client.economic
    first_node = first_client.node
    second_economic = second_client.economic
    second_node = second_client.node
    if delegate_type == 0:
        address, _ = first_economic.account.generate_account(first_node.web3, von_amount(first_economic.delegate_limit, 100))
    else:
        address = create_restricting_plan(first_client)
    log.info("Create delegate account：{}".format(address))
    create_stakings_node(clients_new_node)
    log.info("Create first pledge node id :{}".format(first_node.node_id))
    log.info("Create second pledge node id :{}".format(second_node.node_id))
    result = first_client.delegate.delegate(delegate_type, address)
    assert_code(result, 0)
    log.info("Commissioned successfully, commissioned amount：{}".format(first_economic.delegate_limit))
    result = first_client.delegate.delegate(delegate_type, address, node_id=second_node.node_id)
    assert_code(result, 0)
    log.info("Commissioned successfully, commissioned amount：{}".format(second_economic.delegate_limit))
    first_delegate_epoch, first_cumulative_income = get_dividend_information(first_client, address)
    second_delegate_epoch, second_cumulative_income = get_dividend_information(second_client, address)
    assert first_delegate_epoch == 1, "ErrMsg: Last time first delegate epoch {}".format(first_delegate_epoch)
    assert first_cumulative_income == 0, "ErrMsg: Last time first cumulative income {}".format(first_cumulative_income)
    assert second_delegate_epoch == 1, "ErrMsg: Last time second delegate epoch {}".format(second_delegate_epoch)
    assert second_cumulative_income == 0, "ErrMsg: Last time second cumulative income {}".format(
        second_cumulative_income)
    first_last_delegate_epoch, first_delegate_total, first_delegate_total_hes, first_delegate_reward_total = get_delegate_relevant_amount_and_epoch(
        first_client)
    second_last_delegate_epoch, second_delegate_total, second_delegate_total_hes, second_delegate_reward_total = get_delegate_relevant_amount_and_epoch(
        second_client)
    assert first_last_delegate_epoch == 1, "ErrMsg: Last time delegate epoch {}".format(first_last_delegate_epoch)
    assert first_delegate_total == 0, "The total number of effective commissioned nodes: {}".format(
        first_delegate_total)
    assert first_delegate_total_hes == first_economic.delegate_limit, "The total number of inactive nodes commissioned: {}".format(
        first_delegate_total_hes)
    assert first_delegate_reward_total == 0, "Total delegated rewards currently issued by the candidate: {}".format(
        first_delegate_reward_total)
    assert second_last_delegate_epoch == 1, "ErrMsg: Last time delegate epoch {}".format(second_last_delegate_epoch)
    assert second_delegate_total == 0, "The total number of effective commissioned nodes: {}".format(
        second_delegate_total)
    assert second_delegate_total_hes == second_economic.delegate_limit, "The total number of inactive nodes commissioned: {}".format(
        second_delegate_total_hes)
    assert second_delegate_reward_total == 0, "Total delegated rewards currently issued by the candidate: {}".format(
        second_delegate_reward_total)
    first_economic.wait_settlement_blocknum(first_node)
    log.info("Current settlement block height：{}".format(first_node.eth.blockNumber))
    result = first_client.delegate.delegate(delegate_type, address)
    assert_code(result, 0)
    log.info("Commissioned successfully, commissioned amount：{}".format(first_economic.delegate_limit))
    result = first_client.delegate.delegate(delegate_type, address, node_id=second_node.node_id)
    assert_code(result, 0)
    log.info("Commissioned successfully, commissioned amount：{}".format(second_economic.delegate_limit))
    first_delegate_epoch, first_cumulative_income = get_dividend_information(first_client, address)
    second_delegate_epoch, second_cumulative_income = get_dividend_information(second_client, address)
    assert first_delegate_epoch == 2, "ErrMsg: Last time first delegate epoch {}".format(first_delegate_epoch)
    assert first_cumulative_income == 0, "ErrMsg: Last time first cumulative income {}".format(first_cumulative_income)
    assert second_delegate_epoch == 2, "ErrMsg: Last time second delegate epoch {}".format(second_delegate_epoch)
    assert second_cumulative_income == 0, "ErrMsg: Last time second cumulative income {}".format(
        second_cumulative_income)
    first_last_delegate_epoch, first_delegate_total, first_delegate_total_hes, first_delegate_reward_total = get_delegate_relevant_amount_and_epoch(
        first_client)
    second_last_delegate_epoch, second_delegate_total, second_delegate_total_hes, second_delegate_reward_total = get_delegate_relevant_amount_and_epoch(
        second_client)
    assert first_last_delegate_epoch == 2, "ErrMsg: Last time delegate epoch {}".format(first_last_delegate_epoch)
    assert first_delegate_total == first_economic.delegate_limit, "The total number of effective commissioned nodes: {}".format(
        first_delegate_total)
    assert first_delegate_total_hes == first_economic.delegate_limit, "The total number of inactive nodes commissioned: {}".format(
        first_delegate_total_hes)
    assert first_delegate_reward_total == 0, "Total delegated rewards currently issued by the candidate: {}".format(
        first_delegate_reward_total)
    assert second_last_delegate_epoch == 2, "ErrMsg: Last time delegate epoch {}".format(second_last_delegate_epoch)
    assert second_delegate_total == second_economic.delegate_limit, "The total number of effective commissioned nodes: {}".format(
        second_delegate_total)
    assert second_delegate_total_hes == second_economic.delegate_limit, "The total number of inactive nodes commissioned: {}".format(
        second_delegate_total_hes)
    assert second_delegate_reward_total == 0, "Total delegated rewards currently issued by the candidate: {}".format(
        second_delegate_reward_total)
    block_reward, staking_reward = first_economic.get_current_year_reward(first_node)
    first_economic.wait_settlement_blocknum(first_node)
    log.info("Current settlement block height：{}".format(first_node.eth.blockNumber))
    first_commission_award = first_economic.calculate_delegate_reward(first_node, block_reward, staking_reward)
    second_commission_award = second_economic.calculate_delegate_reward(second_node, block_reward, staking_reward)
    result = first_client.delegate.delegate(delegate_type, address)
    assert_code(result, 0)
    log.info("Commissioned successfully, commissioned amount：{}".format(first_economic.delegate_limit))
    result = first_client.delegate.delegate(delegate_type, address, node_id=second_node.node_id)
    assert_code(result, 0)
    log.info("Commissioned successfully, commissioned amount：{}".format(second_economic.delegate_limit))
    first_delegate_epoch, first_cumulative_income = get_dividend_information(first_client, address)
    second_delegate_epoch, second_cumulative_income = get_dividend_information(second_client, address)
    assert first_delegate_epoch == 3, "ErrMsg: Last time first delegate epoch {}".format(first_delegate_epoch)
    assert first_cumulative_income == first_commission_award, "ErrMsg: Last time first cumulative income {}".format(
        first_cumulative_income)
    assert second_delegate_epoch == 3, "ErrMsg: Last time second delegate epoch {}".format(second_delegate_epoch)
    assert second_cumulative_income == second_commission_award, "ErrMsg: Last time second cumulative income {}".format(
        second_cumulative_income)
    first_last_delegate_epoch, first_delegate_total, first_delegate_total_hes, first_delegate_reward_total = get_delegate_relevant_amount_and_epoch(
        first_client)
    second_last_delegate_epoch, second_delegate_total, second_delegate_total_hes, second_delegate_reward_total = get_delegate_relevant_amount_and_epoch(
        second_client)
    assert first_last_delegate_epoch == 3, "ErrMsg: Last time delegate epoch {}".format(first_last_delegate_epoch)
    assert first_delegate_total == von_amount(first_economic.delegate_limit,
                                              2), "The total number of effective commissioned nodes: {}".format(
        first_delegate_total)
    assert first_delegate_total_hes == first_economic.delegate_limit, "The total number of inactive nodes commissioned: {}".format(
        first_delegate_total_hes)
    assert first_delegate_reward_total == first_commission_award, "Total delegated rewards currently issued by the candidate: {}".format(
        first_delegate_reward_total)
    assert second_last_delegate_epoch == 3, "ErrMsg: Last time delegate epoch {}".format(second_last_delegate_epoch)
    assert second_delegate_total == von_amount(second_economic.delegate_limit,
                                               2), "The total number of effective commissioned nodes: {}".format(
        second_delegate_total)
    assert second_delegate_total_hes == second_economic.delegate_limit, "The total number of inactive nodes commissioned: {}".format(
        second_delegate_total_hes)
    assert second_delegate_reward_total == second_commission_award, "Total delegated rewards currently issued by the candidate: {}".format(
        second_delegate_reward_total)


@pytest.mark.P0
@pytest.mark.parametrize('delegate_type', [0, 1])
def test_EI_BC_010_016(client_new_node, delegate_type):
    """
    自由金额首次部分赎回，验证待领取的委托收益（生效期N）
    锁仓金额首次部分赎回，验证待领取的委托收益（生效期N）
    :param client_new_node:
    :return:
    """
    client = client_new_node
    economic = client.economic
    node = client.node
    if delegate_type == 0:
        address, _ = economic.account.generate_account(node.web3, von_amount(economic.delegate_limit, 100))
    else:
        address = create_restricting_plan(client)
    log.info("Create delegate account：{}".format(address))
    create_staking_node(client)
    log.info("Create pledge node id :{}".format(node.node_id))
    delegate_amount = von_amount(economic.delegate_limit, 10)
    result = client.delegate.delegate(delegate_type, address, amount=delegate_amount)
    assert_code(result, 0)
    log.info("Commissioned successfully, commissioned amount：{}".format(economic.delegate_limit))
    economic.wait_settlement_blocknum(node)
    log.info("Current block height：{}".format(node.eth.blockNumber))
    result = client.ppos.getCandidateInfo(client.node.node_id)
    blocknum = result['Ret']['StakingBlockNum']
    result = client.ppos.getDelegateInfo(blocknum, address, node.node_id)
    log.info("commission information：{}".format(result))
    result = client.delegate.withdrew_delegate(blocknum, address)
    assert_code(result, 0)
    delegate_epoch, cumulative_income = get_dividend_information(client, address)
    assert delegate_epoch == 2, "ErrMsg: Last time first delegate epoch {}".format(delegate_epoch)
    assert cumulative_income == 0, "ErrMsg: Last time first cumulative income {}".format(cumulative_income)
    last_delegate_epoch, delegate_total, delegate_total_hes, delegate_reward_total = get_delegate_relevant_amount_and_epoch(client)
    assert last_delegate_epoch == 1, "ErrMsg: Last time delegate epoch {}".format(last_delegate_epoch)
    assert delegate_total == delegate_amount - economic.delegate_limit, "The total number of effective commissioned nodes: {}".format(delegate_total)
    assert delegate_total_hes == 0, "The total number of inactive nodes commissioned: {}".format(delegate_total_hes)
    assert delegate_reward_total == 0, "Total delegated rewards currently issued by the candidate: {}".format(delegate_reward_total)


@pytest.mark.P0
@pytest.mark.parametrize('delegate_type', [0, 1])
def test_EI_BC_011_072(client_new_node, delegate_type):
    """
    自由金额跨结算周期首次部分赎回，验证待领取的委托收益（生效期N）
    锁仓金额跨结算周期首次部分赎回，验证待领取的委托收益（生效期N）
    :param client_new_node:
    :return:
    """
    client = client_new_node
    economic = client.economic
    node = client.node
    if delegate_type == 0:
        address, _ = economic.account.generate_account(node.web3, von_amount(economic.delegate_limit, 100))
    else:
        address = create_restricting_plan(client)
    log.info("Create delegate account：{}".format(address))
    create_staking_node(client)
    log.info("Create pledge node id :{}".format(node.node_id))
    delegate_amount = von_amount(economic.delegate_limit, 10)
    result = client.delegate.delegate(delegate_type, address, amount=delegate_amount)
    assert_code(result, 0)
    log.info("Commissioned successfully, commissioned amount：{}".format(economic.delegate_limit))
    economic.wait_settlement_blocknum(node)
    block_reward, staking_reward = economic.get_current_year_reward(node)
    economic.wait_settlement_blocknum(node)
    log.info("Current block height：{}".format(node.eth.blockNumber))
    current_entrusted_income = economic.calculate_delegate_reward(node, block_reward, staking_reward)
    result = client.ppos.getCandidateInfo(client.node.node_id)
    blocknum = result['Ret']['StakingBlockNum']
    result = client.delegate.withdrew_delegate(blocknum, address)
    assert_code(result, 0)
    delegate_epoch, cumulative_income = get_dividend_information(client, address)
    assert delegate_epoch == 3, "ErrMsg: Last time first delegate epoch {}".format(delegate_epoch)
    assert cumulative_income == current_entrusted_income, "ErrMsg: Last time first cumulative income {}".format(cumulative_income)
    last_delegate_epoch, delegate_total, delegate_total_hes, delegate_reward_total = get_delegate_relevant_amount_and_epoch(
        client)
    assert last_delegate_epoch == 3, "ErrMsg: Last time delegate epoch {}".format(last_delegate_epoch)
    assert delegate_total == delegate_amount - economic.delegate_limit, "The total number of effective commissioned nodes: {}".format(
        delegate_total)
    assert delegate_total_hes == 0, "The total number of inactive nodes commissioned: {}".format(delegate_total_hes)
    assert delegate_reward_total == current_entrusted_income, "Total delegated rewards currently issued by the candidate: {}".format(
        delegate_reward_total)


@pytest.mark.P0
@pytest.mark.parametrize('delegate_type', [0, 1])
def test_EI_BC_012_017(client_new_node, delegate_type):
    """
    自由金额多次部分赎回，验证待领取的委托收益（单节点）
    锁仓金额多次部分赎回，验证待领取的委托收益（单节点）
    :param client_new_node:
    :return:
    """
    client = client_new_node
    economic = client.economic
    node = client.node
    if delegate_type == 0:
        address, _ = economic.account.generate_account(node.web3, von_amount(economic.delegate_limit, 100))
    else:
        address = create_restricting_plan(client)
    log.info("Create delegate account：{}".format(address))
    create_staking_node(client)
    log.info("Create pledge node id :{}".format(node.node_id))
    delegate_amount = von_amount(economic.delegate_limit, 10)
    result = client.delegate.delegate(0, address, amount=delegate_amount)
    assert_code(result, 0)
    log.info("Commissioned successfully, commissioned amount：{}".format(economic.delegate_limit))
    result = client.ppos.getCandidateInfo(node.node_id)
    blocknum = result['Ret']['StakingBlockNum']
    result = client.delegate.withdrew_delegate(blocknum, address)
    assert_code(result, 0)
    delegate_epoch, cumulative_income = get_dividend_information(client, address)
    assert delegate_epoch == 1, "ErrMsg: Last time first delegate epoch {}".format(delegate_epoch)
    assert cumulative_income == 0, "ErrMsg: Last time first cumulative income {}".format(
        cumulative_income)
    last_delegate_epoch, delegate_total, delegate_total_hes, delegate_reward_total = get_delegate_relevant_amount_and_epoch(
        client)
    assert last_delegate_epoch == 1, "ErrMsg: Last time delegate epoch {}".format(last_delegate_epoch)
    assert delegate_total == 0, "The total number of effective commissioned nodes: {}".format(
        delegate_total)
    assert delegate_total_hes == delegate_amount, "The total number of inactive nodes commissioned: {}".format(delegate_total_hes)
    assert delegate_reward_total == 0, "Total delegated rewards currently issued by the candidate: {}".format(
        delegate_reward_total)
    economic.wait_settlement_blocknum(node)
    result = client.delegate.withdrew_delegate(blocknum, address)
    assert_code(result, 0)
    delegate_epoch, cumulative_income = get_dividend_information(client, address)
    assert delegate_epoch == 2, "ErrMsg: Last time first delegate epoch {}".format(delegate_epoch)
    assert cumulative_income == 0, "ErrMsg: Last time first cumulative income {}".format(
        cumulative_income)
    last_delegate_epoch, delegate_total, delegate_total_hes, delegate_reward_total = get_delegate_relevant_amount_and_epoch(
        client)
    assert last_delegate_epoch == 2, "ErrMsg: Last time delegate epoch {}".format(last_delegate_epoch)
    assert delegate_total == delegate_amount - economic.delegate_limit, "The total number of effective commissioned nodes: {}".format(
        delegate_total)
    assert delegate_total_hes == 0, "The total number of inactive nodes commissioned: {}".format(
        delegate_total_hes)
    assert delegate_reward_total == 0, "Total delegated rewards currently issued by the candidate: {}".format(
        delegate_reward_total)
    block_reward, staking_reward = economic.get_current_year_reward(node)
    economic.wait_settlement_blocknum(node)
    result = client.delegate.withdrew_delegate(blocknum, address)
    assert_code(result, 0)
    current_entrusted_income = economic.calculate_delegate_reward(node, block_reward, staking_reward)
    delegate_epoch, cumulative_income = get_dividend_information(client, address)
    assert delegate_epoch == 3, "ErrMsg: Last time first delegate epoch {}".format(delegate_epoch)
    assert cumulative_income == current_entrusted_income, "ErrMsg: Last time first cumulative income {}".format(
        cumulative_income)
    last_delegate_epoch, delegate_total, delegate_total_hes, delegate_reward_total = get_delegate_relevant_amount_and_epoch(
        client)
    assert last_delegate_epoch == 3, "ErrMsg: Last time delegate epoch {}".format(last_delegate_epoch)
    assert delegate_total == delegate_amount - von_amount(economic.delegate_limit, 2), "The total number of effective commissioned nodes: {}".format(
        delegate_total)
    assert delegate_total_hes == 0, "The total number of inactive nodes commissioned: {}".format(
        delegate_total_hes)
    assert delegate_reward_total == current_entrusted_income, "Total delegated rewards currently issued by the candidate: {}".format(
        delegate_reward_total)


@pytest.mark.P0
def test_EI_BC_013(client_new_node):
    """
    节点被多账户自委托，跨结算周期部分赎回，验证待领取的委托收益（生效期N）
    :param client_new_node:
    :return:
    """
    client = client_new_node
    economic = client.economic
    node = client.node
    first_address, _ = economic.account.generate_account(node.web3, von_amount(economic.delegate_limit, 100))
    second_address, _ = economic.account.generate_account(node.web3, von_amount(economic.delegate_limit, 1000))
    log.info("Create delegate account：{}".format(first_address))
    create_staking_node(client)
    log.info("Create pledge node id :{}".format(node.node_id))
    delegate_amount = von_amount(economic.delegate_limit, 10)
    result = client.delegate.delegate(0, first_address, amount=delegate_amount)
    assert_code(result, 0)
    result = client.delegate.delegate(0, second_address, amount=delegate_amount)
    assert_code(result, 0)
    log.info("Commissioned successfully, commissioned amount：{}".format(economic.delegate_limit))
    economic.wait_settlement_blocknum(node)
    block_reward, staking_reward = economic.get_current_year_reward(node)
    economic.wait_settlement_blocknum(node)
    log.info("Current block height：{}".format(node.eth.blockNumber))
    current_entrusted_income = economic.calculate_delegate_reward(node, block_reward, staking_reward)
    unit_commission_award = math.floor(current_entrusted_income / int((von_amount(economic.delegate_limit, 20) / (10 ** 9))))
    first_commission_award = int((von_amount(economic.delegate_limit, 10) / (10 ** 9)) * unit_commission_award)
    result = client.ppos.getCandidateInfo(client.node.node_id)
    blocknum = result['Ret']['StakingBlockNum']
    result = client.delegate.withdrew_delegate(blocknum, first_address)
    assert_code(result, 0)
    first_delegate_epoch, first_cumulative_income = get_dividend_information(client, first_address)
    assert first_delegate_epoch == 3, "ErrMsg: Last time first delegate epoch {}".format(first_delegate_epoch)
    assert first_cumulative_income == first_commission_award, "ErrMsg: Last time first cumulative income {}".format(
        first_cumulative_income)
    last_delegate_epoch, delegate_total, delegate_total_hes, delegate_reward_total = get_delegate_relevant_amount_and_epoch(
        client)
    assert last_delegate_epoch == 3, "ErrMsg: Last time delegate epoch {}".format(last_delegate_epoch)
    assert delegate_total == delegate_amount - economic.delegate_limit, "The total number of effective commissioned nodes: {}".format(
        delegate_total)
    assert delegate_total_hes == 0, "The total number of inactive nodes commissioned: {}".format(delegate_total_hes)
    assert delegate_reward_total == current_entrusted_income, "Total delegated rewards currently issued by the candidate: {}".format(
        delegate_reward_total)
    result = client.ppos.getDelegateReward(first_address)
    withdrawal_commission = result['Ret'][0]['reward']
    log.info("{} Dividends can be collected in the current settlement period： {}".format(first_address, withdrawal_commission))
    assert withdrawal_commission == first_commission_award, "ErrMsg: Dividends currently available {}".format(withdrawal_commission)


@pytest.mark.P1
@pytest.mark.parametrize('delegate_type', [0, 1])
def test_EI_BC_014_018(client_new_node, delegate_type):
    """
    自由金额赎回全部委托，验证待领取的委托收益（未生效期）
    锁仓金额赎回全部委托，验证待领取的委托收益（未生效期）
    :param client_new_node:
    :return:
    """
    client = client_new_node
    economic = client.economic
    node = client.node
    if delegate_type == 0:
        address, _ = economic.account.generate_account(node.web3, von_amount(economic.delegate_limit, 100))
    else:
        address = create_restricting_plan(client)
    log.info("Create delegate account：{}".format(address))
    create_staking_node(client)
    log.info("Create pledge node id :{}".format(node.node_id))
    delegate_amount = von_amount(economic.delegate_limit, 10)
    result = client.delegate.delegate(0, address, amount=delegate_amount)
    assert_code(result, 0)
    log.info("Commissioned successfully, commissioned amount：{}".format(economic.delegate_limit))
    result = client.ppos.getCandidateInfo(node.node_id)
    blocknum = result['Ret']['StakingBlockNum']
    result = client.delegate.withdrew_delegate(blocknum, address, amount=delegate_amount)
    assert_code(result, 0)
    result = client.ppos.getDelegateInfo(blocknum, address, node.node_id)
    assert_code(result, 301205)
    last_delegate_epoch, delegate_total, delegate_total_hes, delegate_reward_total = get_delegate_relevant_amount_and_epoch(
        client)
    assert last_delegate_epoch == 1, "ErrMsg: Last time delegate epoch {}".format(last_delegate_epoch)
    assert delegate_total == 0, "The total number of effective commissioned nodes: {}".format(
        delegate_total)
    assert delegate_total_hes == 0, "The total number of inactive nodes commissioned: {}".format(
        delegate_total_hes)
    assert delegate_reward_total == 0, "Total delegated rewards currently issued by the candidate: {}".format(
        delegate_reward_total)


@pytest.mark.P1
@pytest.mark.parametrize('delegate_type', [0, 1])
def test_EI_BC_019_023(clients_new_node, delegate_type):
    """
    自由金额首次部分赎回，验证待领取的委托收益（多节点）
    锁仓金额首次部分赎回，验证待领取的委托收益（多节点）
    :param clients_new_node:
    :return:
    """
    first_client = clients_new_node[0]
    second_client = clients_new_node[1]
    first_economic = first_client.economic
    first_node = first_client.node
    second_economic = second_client.economic
    second_node = second_client.node
    if delegate_type == 0:
        address, _ = first_economic.account.generate_account(first_node.web3, von_amount(first_economic.delegate_limit, 100))
    else:
        address = create_restricting_plan(first_client)
    log.info("Create delegate account：{}".format(address))
    create_stakings_node(clients_new_node)
    log.info("Create first pledge node id :{}".format(first_node.node_id))
    log.info("Create second pledge node id :{}".format(second_node.node_id))
    delegate_amount = von_amount(first_economic.delegate_limit, 10)
    result = first_client.delegate.delegate(delegate_type, address, amount=delegate_amount)
    assert_code(result, 0)
    log.info("Commissioned successfully, commissioned amount：{}".format(delegate_amount))
    result = first_client.delegate.delegate(delegate_type, address, node_id=second_node.node_id, amount=delegate_amount)
    assert_code(result, 0)
    log.info("Commissioned successfully, commissioned amount：{}".format(delegate_amount))
    result = first_client.ppos.getCandidateInfo(first_node.node_id)
    first_blocknum = result['Ret']['StakingBlockNum']
    result = second_client.ppos.getCandidateInfo(second_node.node_id)
    second_blocknum = result['Ret']['StakingBlockNum']
    result = first_client.delegate.withdrew_delegate(first_blocknum, address)
    assert_code(result, 0)
    result = first_client.delegate.withdrew_delegate(second_blocknum, address, node_id=second_node.node_id)
    assert_code(result, 0)
    first_delegate_epoch, first_cumulative_income = get_dividend_information(first_client, address)
    assert first_delegate_epoch == 1, "ErrMsg: Last time first delegate epoch {}".format(first_delegate_epoch)
    assert first_cumulative_income == 0, "ErrMsg: Last time first cumulative income {}".format(
        first_cumulative_income)
    second_delegate_epoch, second_cumulative_income = get_dividend_information(second_client, address)
    assert second_delegate_epoch == 1, "ErrMsg: Last time first delegate epoch {}".format(second_delegate_epoch)
    assert second_cumulative_income == 0, "ErrMsg: Last time first cumulative income {}".format(
        second_cumulative_income)
    first_last_delegate_epoch, first_delegate_total, first_delegate_total_hes, first_delegate_reward_total = get_delegate_relevant_amount_and_epoch(
        first_client)
    assert first_last_delegate_epoch == 1, "ErrMsg: Last time delegate epoch {}".format(first_last_delegate_epoch)
    assert first_delegate_total == 0, "The total number of effective commissioned nodes: {}".format(
        first_delegate_total)
    assert first_delegate_total_hes == delegate_amount - first_economic.delegate_limit, "The total number of inactive nodes commissioned: {}".format(
        first_delegate_total_hes)
    assert first_delegate_reward_total == 0, "Total delegated rewards currently issued by the candidate: {}".format(
        first_delegate_reward_total)
    second_last_delegate_epoch, second_delegate_total, second_delegate_total_hes, second_delegate_reward_total = get_delegate_relevant_amount_and_epoch(
        second_client)
    assert second_last_delegate_epoch == 1, "ErrMsg: Last time delegate epoch {}".format(second_last_delegate_epoch)
    assert second_delegate_total == 0, "The total number of effective commissioned nodes: {}".format(
        second_delegate_total)
    assert second_delegate_total_hes == delegate_amount - second_economic.delegate_limit, "The total number of inactive nodes commissioned: {}".format(
        second_delegate_total_hes)
    assert second_delegate_reward_total == 0, "Total delegated rewards currently issued by the candidate: {}".format(
        second_delegate_reward_total)


@pytest.mark.P1
@pytest.mark.parametrize('delegate_type', [0, 1])
def test_EI_BC_020_24(clients_new_node, delegate_type):
    """
    自由金额多次部分赎回，验证待领取的委托收益（多节点）
    锁仓金额多次部分赎回，验证待领取的委托收益（多节点）
    :param clients_new_node:
    :return:
    """
    first_client = clients_new_node[0]
    second_client = clients_new_node[1]
    first_economic = first_client.economic
    first_node = first_client.node
    second_economic = second_client.economic
    second_node = second_client.node
    if delegate_type == 0:
        address, _ = first_economic.account.generate_account(first_node.web3, von_amount(first_economic.delegate_limit, 100))
    else:
        address = create_restricting_plan(first_client)
    log.info("Create delegate account：{}".format(address))
    create_stakings_node(clients_new_node)
    log.info("Create first pledge node id :{}".format(first_node.node_id))
    log.info("Create second pledge node id :{}".format(second_node.node_id))
    delegate_amount = von_amount(first_economic.delegate_limit, 10)
    result = first_client.delegate.delegate(0, address, amount=delegate_amount)
    assert_code(result, 0)
    log.info("Commissioned successfully, commissioned amount：{}".format(delegate_amount))
    result = first_client.delegate.delegate(0, address, node_id=second_node.node_id, amount=delegate_amount)
    assert_code(result, 0)
    log.info("Commissioned successfully, commissioned amount：{}".format(delegate_amount))
    result = first_client.ppos.getCandidateInfo(first_node.node_id)
    first_blocknum = result['Ret']['StakingBlockNum']
    result = second_client.ppos.getCandidateInfo(second_node.node_id)
    second_blocknum = result['Ret']['StakingBlockNum']
    result = first_client.delegate.withdrew_delegate(first_blocknum, address)
    assert_code(result, 0)
    result = first_client.delegate.withdrew_delegate(second_blocknum, address, node_id=second_node.node_id)
    assert_code(result, 0)
    first_delegate_epoch, first_cumulative_income = get_dividend_information(first_client, address)
    assert first_delegate_epoch == 1, "ErrMsg: Last time first delegate epoch {}".format(first_delegate_epoch)
    assert first_cumulative_income == 0, "ErrMsg: Last time first cumulative income {}".format(
        first_cumulative_income)
    second_delegate_epoch, second_cumulative_income = get_dividend_information(second_client, address)
    assert second_delegate_epoch == 1, "ErrMsg: Last time first delegate epoch {}".format(second_delegate_epoch)
    assert second_cumulative_income == 0, "ErrMsg: Last time first cumulative income {}".format(
        second_cumulative_income)
    first_last_delegate_epoch, first_delegate_total, first_delegate_total_hes, first_delegate_reward_total = get_delegate_relevant_amount_and_epoch(
        first_client)
    assert first_last_delegate_epoch == 1, "ErrMsg: Last time delegate epoch {}".format(first_last_delegate_epoch)
    assert first_delegate_total == 0, "The total number of effective commissioned nodes: {}".format(
        first_delegate_total)
    assert first_delegate_total_hes == delegate_amount - first_economic.delegate_limit, "The total number of inactive nodes commissioned: {}".format(
        first_delegate_total_hes)
    assert first_delegate_reward_total == 0, "Total delegated rewards currently issued by the candidate: {}".format(
        first_delegate_reward_total)
    second_last_delegate_epoch, second_delegate_total, second_delegate_total_hes, second_delegate_reward_total = get_delegate_relevant_amount_and_epoch(
        second_client)
    assert second_last_delegate_epoch == 1, "ErrMsg: Last time delegate epoch {}".format(second_last_delegate_epoch)
    assert second_delegate_total == 0, "The total number of effective commissioned nodes: {}".format(
        second_delegate_total)
    assert second_delegate_total_hes == delegate_amount - second_economic.delegate_limit, "The total number of inactive nodes commissioned: {}".format(
        second_delegate_total_hes)
    assert second_delegate_reward_total == 0, "Total delegated rewards currently issued by the candidate: {}".format(
        second_delegate_reward_total)
    first_economic.wait_settlement_blocknum(first_node)
    log.info("Current settlement block height：{}".format(first_node.eth.blockNumber))
    result = first_client.delegate.withdrew_delegate(first_blocknum, address)
    assert_code(result, 0)
    result = first_client.delegate.withdrew_delegate(second_blocknum, address, node_id=second_node.node_id)
    assert_code(result, 0)
    first_delegate_epoch, first_cumulative_income = get_dividend_information(first_client, address)
    assert first_delegate_epoch == 2, "ErrMsg: Last time first delegate epoch {}".format(first_delegate_epoch)
    assert first_cumulative_income == 0, "ErrMsg: Last time first cumulative income {}".format(
        first_cumulative_income)
    second_delegate_epoch, second_cumulative_income = get_dividend_information(second_client, address)
    assert second_delegate_epoch == 2, "ErrMsg: Last time first delegate epoch {}".format(second_delegate_epoch)
    assert second_cumulative_income == 0, "ErrMsg: Last time first cumulative income {}".format(
        second_cumulative_income)
    first_last_delegate_epoch, first_delegate_total, first_delegate_total_hes, first_delegate_reward_total = get_delegate_relevant_amount_and_epoch(
        first_client)
    assert first_last_delegate_epoch == 1, "ErrMsg: Last time delegate epoch {}".format(first_last_delegate_epoch)
    assert first_delegate_total == 0, "The total number of effective commissioned nodes: {}".format(
        first_delegate_total)
    assert first_delegate_total_hes == delegate_amount - von_amount(first_economic.delegate_limit, 2), "The total number of inactive nodes commissioned: {}".format(
        first_delegate_total_hes)
    assert first_delegate_reward_total == 0, "Total delegated rewards currently issued by the candidate: {}".format(
        first_delegate_reward_total)
    second_last_delegate_epoch, second_delegate_total, second_delegate_total_hes, second_delegate_reward_total = get_delegate_relevant_amount_and_epoch(
        second_client)
    assert second_last_delegate_epoch == 1, "ErrMsg: Last time delegate epoch {}".format(second_last_delegate_epoch)
    assert second_delegate_total == 0, "The total number of effective commissioned nodes: {}".format(
        second_delegate_total)
    assert second_delegate_total_hes == delegate_amount - von_amount(second_economic.delegate_limit, 2), "The total number of inactive nodes commissioned: {}".format(
        second_delegate_total_hes)
    assert second_delegate_reward_total == 0, "Total delegated rewards currently issued by the candidate: {}".format(
        second_delegate_reward_total)
    block_reward, staking_reward = first_economic.get_current_year_reward(first_node)
    first_economic.wait_settlement_blocknum(first_node)
    log.info("Current settlement block height：{}".format(first_node.eth.blockNumber))
    first_commission_award = first_economic.calculate_delegate_reward(first_node, block_reward, staking_reward)
    second_commission_award = second_economic.calculate_delegate_reward(second_node, block_reward, staking_reward)

    result = first_client.delegate.withdrew_delegate(first_blocknum, address)
    assert_code(result, 0)
    result = first_client.delegate.withdrew_delegate(second_blocknum, address, node_id=second_node.node_id)
    assert_code(result, 0)
    first_delegate_epoch, first_cumulative_income = get_dividend_information(first_client, address)
    assert first_delegate_epoch == 3, "ErrMsg: Last time first delegate epoch {}".format(first_delegate_epoch)
    assert first_cumulative_income == first_commission_award, "ErrMsg: Last time first cumulative income {}".format(
        first_cumulative_income)
    second_delegate_epoch, second_cumulative_income = get_dividend_information(second_client, address)
    assert second_delegate_epoch == 3, "ErrMsg: Last time first delegate epoch {}".format(second_delegate_epoch)
    assert second_cumulative_income == second_commission_award, "ErrMsg: Last time first cumulative income {}".format(
        second_cumulative_income)
    first_last_delegate_epoch, first_delegate_total, first_delegate_total_hes, first_delegate_reward_total = get_delegate_relevant_amount_and_epoch(
        first_client)
    assert first_last_delegate_epoch == 1, "ErrMsg: Last time delegate epoch {}".format(first_last_delegate_epoch)
    assert first_delegate_total == 0, "The total number of effective commissioned nodes: {}".format(
        first_delegate_total)
    assert first_delegate_total_hes == delegate_amount - von_amount(first_economic.delegate_limit,
                                                                    3), "The total number of inactive nodes commissioned: {}".format(
        first_delegate_total_hes)
    assert first_delegate_reward_total == first_commission_award, "Total delegated rewards currently issued by the candidate: {}".format(
        first_delegate_reward_total)
    second_last_delegate_epoch, second_delegate_total, second_delegate_total_hes, second_delegate_reward_total = get_delegate_relevant_amount_and_epoch(
        second_client)
    assert second_last_delegate_epoch == 1, "ErrMsg: Last time delegate epoch {}".format(second_last_delegate_epoch)
    assert second_delegate_total == 0, "The total number of effective commissioned nodes: {}".format(
        second_delegate_total)
    assert second_delegate_total_hes == delegate_amount - von_amount(second_economic.delegate_limit,
                                                                     3), "The total number of inactive nodes commissioned: {}".format(
        second_delegate_total_hes)
    assert second_delegate_reward_total == second_commission_award, "Total delegated rewards currently issued by the candidate: {}".format(
        second_delegate_reward_total)


@pytest.mark.P1
@pytest.mark.parametrize('delegate_type', [0, 1])
def test_EI_BC_021_25(clients_new_node, delegate_type):
    """
    未生效期N自由金额赎回全部委托，验证待领取的委托收益（多节点）
    未生效期N锁仓金额赎回全部委托，验证待领取的委托收益（多节点）
    :param clients_new_node:
    :return:
    """
    first_client = clients_new_node[0]
    second_client = clients_new_node[1]
    first_economic = first_client.economic
    first_node = first_client.node
    second_economic = second_client.economic
    second_node = second_client.node
    if delegate_type == 0:
        address, _ = first_economic.account.generate_account(first_node.web3, von_amount(first_economic.delegate_limit, 100))
    else:
        address = create_restricting_plan(first_client)
    log.info("Create delegate account：{}".format(address))
    create_stakings_node(clients_new_node)
    log.info("Create first pledge node id :{}".format(first_node.node_id))
    log.info("Create second pledge node id :{}".format(second_node.node_id))
    delegate_amount = von_amount(first_economic.delegate_limit, 10)
    result = first_client.delegate.delegate(0, address, amount=delegate_amount)
    assert_code(result, 0)
    log.info("Commissioned successfully, commissioned amount：{}".format(delegate_amount))
    result = first_client.delegate.delegate(0, address, node_id=second_node.node_id, amount=delegate_amount)
    assert_code(result, 0)
    log.info("Commissioned successfully, commissioned amount：{}".format(delegate_amount))
    result = first_client.ppos.getCandidateInfo(first_node.node_id)
    first_blocknum = result['Ret']['StakingBlockNum']
    result = second_client.ppos.getCandidateInfo(second_node.node_id)
    second_blocknum = result['Ret']['StakingBlockNum']
    result = first_client.delegate.withdrew_delegate(first_blocknum, address, amount=delegate_amount)
    assert_code(result, 0)
    result = first_client.delegate.withdrew_delegate(second_blocknum, address, node_id=second_node.node_id,amount=delegate_amount)
    assert_code(result, 0)
    result = first_client.ppos.getDelegateInfo(first_blocknum, address, first_node.node_id)
    assert_code(result, 301205)
    result = second_client.ppos.getDelegateInfo(first_blocknum, address, second_node.node_id)
    assert_code(result, 301205)
    first_last_delegate_epoch, first_delegate_total, first_delegate_total_hes, first_delegate_reward_total = get_delegate_relevant_amount_and_epoch(
        first_client)
    assert first_last_delegate_epoch == 1, "ErrMsg: Last time delegate epoch {}".format(first_last_delegate_epoch)
    assert first_delegate_total == 0, "The total number of effective commissioned nodes: {}".format(
        first_delegate_total)
    assert first_delegate_total_hes == 0, "The total number of inactive nodes commissioned: {}".format(
        first_delegate_total_hes)
    assert first_delegate_reward_total == 0, "Total delegated rewards currently issued by the candidate: {}".format(
        first_delegate_reward_total)
    second_last_delegate_epoch, second_delegate_total, second_delegate_total_hes, second_delegate_reward_total = get_delegate_relevant_amount_and_epoch(
        second_client)
    assert second_last_delegate_epoch == 1, "ErrMsg: Last time delegate epoch {}".format(second_last_delegate_epoch)
    assert second_delegate_total == 0, "The total number of effective commissioned nodes: {}".format(
        second_delegate_total)
    assert second_delegate_total_hes == 0, "The total number of inactive nodes commissioned: {}".format(
        second_delegate_total_hes)
    assert second_delegate_reward_total == 0, "Total delegated rewards currently issued by the candidate: {}".format(
        second_delegate_reward_total)


@pytest.mark.P1
@pytest.mark.parametrize('delegate_type', [0, 1])
def test_EI_BC_022_026(clients_new_node, delegate_type):
    """
    生效期N自由金额赎回全部委托，验证待领取的委托收益（多节点）
    生效期N锁仓金额赎回全部委托，验证待领取的委托收益（多节点）
    :param clients_new_node:
    :return:
    """
    first_client = clients_new_node[0]
    second_client = clients_new_node[1]
    first_economic = first_client.economic
    first_node = first_client.node
    second_economic = second_client.economic
    second_node = second_client.node
    if delegate_type == 0:
        address, _ = first_economic.account.generate_account(first_node.web3,von_amount(first_economic.delegate_limit, 100))
    else:
        address = create_restricting_plan(first_client)
    log.info("Create delegate account：{}".format(address))
    create_stakings_node(clients_new_node)
    log.info("Create first pledge node id :{}".format(first_node.node_id))
    log.info("Create second pledge node id :{}".format(second_node.node_id))
    delegate_amount = von_amount(first_economic.delegate_limit, 10)
    result = first_client.delegate.delegate(delegate_type, address, amount=delegate_amount)
    assert_code(result, 0)
    log.info("Commissioned successfully, commissioned amount：{}".format(delegate_amount))
    result = first_client.delegate.delegate(delegate_type, address, node_id=second_node.node_id, amount=delegate_amount)
    assert_code(result, 0)
    log.info("Commissioned successfully, commissioned amount：{}".format(delegate_amount))
    first_economic.wait_settlement_blocknum(first_node)
    log.info("Current settlement block height：{}".format(first_node.eth.blockNumber))
    block_reward, staking_reward = first_economic.get_current_year_reward(first_node)
    first_economic.wait_settlement_blocknum(first_node)
    log.info("Current settlement block height：{}".format(first_node.eth.blockNumber))
    result = first_client.ppos.getCandidateInfo(first_node.node_id)
    first_blocknum = result['Ret']['StakingBlockNum']
    result = second_client.ppos.getCandidateInfo(second_node.node_id)
    second_blocknum = result['Ret']['StakingBlockNum']
    first_commission_award = first_economic.calculate_delegate_reward(first_node, block_reward, staking_reward)
    second_commission_award = second_economic.calculate_delegate_reward(second_node, block_reward, staking_reward)
    result = first_client.delegate.withdrew_delegate(first_blocknum, address, amount=delegate_amount)
    assert_code(result, 0)
    result = first_client.delegate.withdrew_delegate(second_blocknum, address, node_id=second_node.node_id, amount=delegate_amount)
    assert_code(result, 0)
    result = first_client.ppos.getDelegateInfo(first_blocknum, address, first_node.node_id)
    assert_code(result, 301205)
    result = second_client.ppos.getDelegateInfo(first_blocknum, address, second_node.node_id)
    assert_code(result, 301205)
    first_last_delegate_epoch, first_delegate_total, first_delegate_total_hes, first_delegate_reward_total = get_delegate_relevant_amount_and_epoch(
        first_client)
    assert first_last_delegate_epoch == 1, "ErrMsg: Last time delegate epoch {}".format(first_last_delegate_epoch)
    assert first_delegate_total == 0, "The total number of effective commissioned nodes: {}".format(first_delegate_total)
    assert first_delegate_total_hes == 0, "The total number of inactive nodes commissioned: {}".format(
        first_delegate_total_hes)
    assert first_delegate_reward_total == first_commission_award, "Total delegated rewards currently issued by the candidate: {}".format(
        first_delegate_reward_total)
    second_last_delegate_epoch, second_delegate_total, second_delegate_total_hes, second_delegate_reward_total = get_delegate_relevant_amount_and_epoch(
        second_client)
    assert second_last_delegate_epoch == 1, "ErrMsg: Last time delegate epoch {}".format(second_last_delegate_epoch)
    assert second_delegate_total == 0, "The total number of effective commissioned nodes: {}".format(second_delegate_total)
    assert second_delegate_total_hes == 0, "The total number of inactive nodes commissioned: {}".format(
        second_delegate_total_hes)
    assert second_delegate_reward_total == second_commission_award, "Total delegated rewards currently issued by the candidate: {}".format(
        second_delegate_reward_total)
    result = first_client.ppos.getDelegateReward(address, node_ids=first_node.node_id)
    assert result['Ret'][0]['reward'] == 0, "ErrMsg: Receive dividends {}".format(result['Ret'][0]['reward'])
    result = first_client.ppos.getDelegateReward(address, node_ids=second_node.node_id)
    assert result['Ret'][0]['reward'] == 0, "ErrMsg: Receive dividends {}".format(result['Ret'][0]['reward'])


@pytest.mark.P1
@pytest.mark.parametrize('delegate_type', [0, 1])
def test_EI_BC_027_029(client_new_node, delegate_type):
    """
    跨结算期赎回全部自由委托（生效期N赎回）
    跨结算期赎回全部锁仓委托（生效期N赎回）
    :param client_new_node:
    :return:
    """
    client = client_new_node
    economic = client.economic
    node = client.node
    if delegate_type == 0:
        address, _ = economic.account.generate_account(node.web3, von_amount(economic.delegate_limit, 100))
    else:
        address = create_restricting_plan(client)
    log.info("Create delegate account：{}".format(address))
    create_staking_node(client)
    log.info("Create pledge node id :{}".format(node.node_id))
    delegate_amount = von_amount(economic.delegate_limit, 10)
    result = client.delegate.delegate(0, address, amount=delegate_amount)
    assert_code(result, 0)
    log.info("Commissioned successfully, commissioned amount：{}".format(economic.delegate_limit))
    result = client.ppos.getCandidateInfo(node.node_id)
    blocknum = result['Ret']['StakingBlockNum']
    result = client.delegate.withdrew_delegate(blocknum, address, amount=delegate_amount)
    assert_code(result, 0)
    economic.wait_settlement_blocknum(node)
    log.info("Current settlement block height：{}".format(node.eth.blockNumber))
    result = client.ppos.getCandidateInfo(node.node_id)
    blocknum = result['Ret']['StakingBlockNum']
    result = client.delegate.withdrew_delegate(blocknum, address, amount=delegate_amount)
    assert_code(result, 0)
    result = client.ppos.getDelegateInfo(blocknum, address, node.node_id)
    assert_code(result, 301205)
    first_last_delegate_epoch, first_delegate_total, first_delegate_total_hes, first_delegate_reward_total = get_delegate_relevant_amount_and_epoch(
        client)
    assert first_last_delegate_epoch == 1, "ErrMsg: Last time delegate epoch {}".format(first_last_delegate_epoch)
    assert first_delegate_total == 0, "The total number of effective commissioned nodes: {}".format(
        first_delegate_total)
    assert first_delegate_total_hes == 0, "The total number of inactive nodes commissioned: {}".format(
        first_delegate_total_hes)
    assert first_delegate_reward_total == 0, "Total delegated rewards currently issued by the candidate: {}".format(
        first_delegate_reward_total)
    result = client.ppos.getDelegateReward(address)
    assert_code(result, 00)


@pytest.mark.P1
@pytest.mark.parametrize('delegate_type', [0, 1])
def test_EI_BC_028_030(client_new_node, delegate_type):
    """
    跨结算期赎回全部自由委托（生效期N+1赎回）
    跨结算期赎回全部锁仓委托（生效期N+1赎回）
    :param client_new_node:
    :return:
    """
    client = client_new_node
    economic = client.economic
    node = client.node
    if delegate_type == 0:
        address, _ = economic.account.generate_account(node.web3, von_amount(economic.delegate_limit, 100))
    else:
        address = create_restricting_plan(client)
    log.info("Create delegate account：{}".format(address))
    create_staking_node(client)
    log.info("Create pledge node id :{}".format(node.node_id))
    delegate_amount = von_amount(economic.delegate_limit, 10)
    result = client.delegate.delegate(0, address, amount=delegate_amount)
    assert_code(result, 0)
    log.info("Commissioned successfully, commissioned amount：{}".format(economic.delegate_limit))
    economic.wait_settlement_blocknum(node)
    log.info("Current settlement block height：{}".format(node.eth.blockNumber))
    result = client.ppos.getCandidateInfo(node.node_id)
    blocknum = result['Ret']['StakingBlockNum']
    block_reward, staking_reward = economic.get_current_year_reward(node)
    economic.wait_settlement_blocknum(node)
    log.info("Current settlement block height：{}".format(node.eth.blockNumber))
    commission_award = economic.calculate_delegate_reward(node, block_reward, staking_reward)
    result = client.delegate.withdrew_delegate(blocknum, address, amount=delegate_amount)
    assert_code(result, 0)
    result = client.ppos.getDelegateInfo(blocknum, address, node.node_id)
    assert_code(result, 301205)
    first_last_delegate_epoch, first_delegate_total, first_delegate_total_hes, first_delegate_reward_total = get_delegate_relevant_amount_and_epoch(
        client)
    assert first_last_delegate_epoch == 1, "ErrMsg: Last time delegate epoch {}".format(first_last_delegate_epoch)
    assert first_delegate_total == 0, "The total number of effective commissioned nodes: {}".format(
        first_delegate_total)
    assert first_delegate_total_hes == 0, "The total number of inactive nodes commissioned: {}".format(
        first_delegate_total_hes)
    assert first_delegate_reward_total == commission_award, "Total delegated rewards currently issued by the candidate: {}".format(
        first_delegate_reward_total)
    result = client.ppos.getDelegateReward(address)
    assert_code(result, 00)


@pytest.mark.P1
@pytest.mark.parametrize('delegate_type', [0, 1])
def test_EI_BC_055_058(client_new_node, delegate_type):
    """
    自由金额委托首次领取分红（生效期N）
    锁仓金额委托首次领取分红（生效期N）
    :param client_new_node:
    :return:
    """
    client = client_new_node
    economic = client.economic
    node = client.node
    if delegate_type == 0:
        address, _ = economic.account.generate_account(node.web3, von_amount(economic.delegate_limit, 100))
    else:
        address = create_restricting_plan(client)
    log.info("Create delegate account：{}".format(address))
    create_staking_node(client)
    log.info("Create pledge node id :{}".format(node.node_id))
    delegate_amount = von_amount(economic.delegate_limit, 10)
    result = client.delegate.delegate(0, address, amount=delegate_amount)
    assert_code(result, 0)
    log.info("Commissioned successfully, commissioned amount：{}".format(economic.delegate_limit))
    economic.wait_settlement_blocknum(node)
    log.info("Current settlement block height：{}".format(node.eth.blockNumber))
    result = client.delegate.withdraw_delegate_reward(address)
    assert_code(result, 0)
    delegate_epoch, cumulative_income = get_dividend_information(client, address)
    assert delegate_epoch == 2, "ErrMsg: Last time first delegate epoch {}".format(delegate_epoch)
    assert cumulative_income == 0, "ErrMsg: Last time first cumulative income {}".format(cumulative_income)
    last_delegate_epoch, delegate_total, delegate_total_hes, delegate_reward_total = get_delegate_relevant_amount_and_epoch(
        client)
    assert last_delegate_epoch == 1, "ErrMsg: Last time delegate epoch {}".format(last_delegate_epoch)
    assert delegate_total == delegate_amount, "The total number of effective commissioned nodes: {}".format(
        delegate_total)
    assert delegate_total_hes == 0, "The total number of inactive nodes commissioned: {}".format(delegate_total_hes)
    assert delegate_reward_total == 0, "Total delegated rewards currently issued by the candidate: {}".format(
        delegate_reward_total)
    result = client.ppos.getDelegateReward(address)
    assert result['Ret'][0]['reward'] == 0, "ErrMsg:Receive dividends {}".format(result['Ret'][0]['reward'])


@pytest.mark.P1
@pytest.mark.parametrize('delegate_type', [0, 1])
def test_EI_BC_056_059(client_new_node, delegate_type):
    """
    跨结算期自由金额委托首次领取分红
    跨结算期锁仓金额委托首次领取分红
    :param client_new_node:
    :return:
    """
    client = client_new_node
    economic = client.economic
    node = client.node
    if delegate_type == 0:
        address, _ = economic.account.generate_account(node.web3, von_amount(economic.delegate_limit, 100))
    else:
        address = create_restricting_plan(client)
    log.info("Create delegate account：{}".format(address))
    create_staking_node(client)
    log.info("Create pledge node id :{}".format(node.node_id))
    delegate_amount = von_amount(economic.delegate_limit, 10)
    result = client.delegate.delegate(delegate_type, address, amount=delegate_amount)
    assert_code(result, 0)
    log.info("Commissioned successfully, commissioned amount：{}".format(economic.delegate_limit))
    result = client.delegate.withdraw_delegate_reward(address)
    assert_code(result, 0)
    delegate_epoch, cumulative_income = get_dividend_information(client, address)
    assert delegate_epoch == 1, "ErrMsg: Last time first delegate epoch {}".format(delegate_epoch)
    assert cumulative_income == 0, "ErrMsg: Last time first cumulative income {}".format(cumulative_income)
    last_delegate_epoch, delegate_total, delegate_total_hes, delegate_reward_total = get_delegate_relevant_amount_and_epoch(
        client)
    assert last_delegate_epoch == 1, "ErrMsg: Last time delegate epoch {}".format(last_delegate_epoch)
    assert delegate_total == 0, "The total number of effective commissioned nodes: {}".format(
        delegate_total)
    assert delegate_total_hes == delegate_amount, "The total number of inactive nodes commissioned: {}".format(delegate_total_hes)
    assert delegate_reward_total == 0, "Total delegated rewards currently issued by the candidate: {}".format(
        delegate_reward_total)
    result = client.ppos.getDelegateReward(address)
    assert result['Ret'][0]['reward'] == 0, "ErrMsg:Receive dividends {}".format(result['Ret'][0]['reward'])
    economic.wait_settlement_blocknum(node)
    log.info("Current settlement block height：{}".format(node.eth.blockNumber))
    result = client.delegate.withdraw_delegate_reward(address)
    assert_code(result, 0)
    delegate_epoch, cumulative_income = get_dividend_information(client, address)
    assert delegate_epoch == 2, "ErrMsg: Last time first delegate epoch {}".format(delegate_epoch)
    assert cumulative_income == 0, "ErrMsg: Last time first cumulative income {}".format(cumulative_income)
    last_delegate_epoch, delegate_total, delegate_total_hes, delegate_reward_total = get_delegate_relevant_amount_and_epoch(
        client)
    assert last_delegate_epoch == 1, "ErrMsg: Last time delegate epoch {}".format(last_delegate_epoch)
    assert delegate_total == delegate_amount, "The total number of effective commissioned nodes: {}".format(
        delegate_total)
    assert delegate_total_hes == 0, "The total number of inactive nodes commissioned: {}".format(
        delegate_total_hes)
    assert delegate_reward_total == 0, "Total delegated rewards currently issued by the candidate: {}".format(
        delegate_reward_total)
    result = client.ppos.getDelegateReward(address)
    assert result['Ret'][0]['reward'] == 0, "ErrMsg:Receive dividends {}".format(result['Ret'][0]['reward'])
    block_reward, staking_reward = economic.get_current_year_reward(node)
    economic.wait_settlement_blocknum(node)
    log.info("Current settlement block height：{}".format(node.eth.blockNumber))
    commission_award = economic.calculate_delegate_reward(node, block_reward, staking_reward)
    result = client.delegate.withdraw_delegate_reward(address)
    assert_code(result, 0)
    delegate_epoch, cumulative_income = get_dividend_information(client, address)
    assert delegate_epoch == 3, "ErrMsg: Last time first delegate epoch {}".format(delegate_epoch)
    assert cumulative_income == 0, "ErrMsg: Last time first cumulative income {}".format(cumulative_income)
    last_delegate_epoch, delegate_total, delegate_total_hes, delegate_reward_total = get_delegate_relevant_amount_and_epoch(
        client)
    assert last_delegate_epoch == 1, "ErrMsg: Last time delegate epoch {}".format(last_delegate_epoch)
    assert delegate_total == delegate_amount, "The total number of effective commissioned nodes: {}".format(
        delegate_total)
    assert delegate_total_hes == 0, "The total number of inactive nodes commissioned: {}".format(
        delegate_total_hes)
    assert delegate_reward_total == commission_award, "Total delegated rewards currently issued by the candidate: {}".format(
        delegate_reward_total)
    result = client.ppos.getDelegateReward(address)
    assert result['Ret'][0]['reward'] == 0, "ErrMsg:Receive dividends {}".format(result['Ret'][0]['reward'])


@pytest.mark.P1
def test_EI_BC_060(client_new_node):
    """
    自由、锁仓金额同时委托首次领取分红（未生效期N）
    :param client_new_node:
    :return:
    """
    client = client_new_node
    economic = client.economic
    node = client.node
    address = create_restricting_plan(client)
    log.info("Create delegate account：{}".format(address))
    create_staking_node(client)
    log.info("Create pledge node id :{}".format(node.node_id))
    delegate_amount = von_amount(economic.delegate_limit, 10)
    result = client.delegate.delegate(0, address, amount=delegate_amount)
    assert_code(result, 0)
    result = client.delegate.delegate(1, address, amount=delegate_amount)
    assert_code(result, 0)
    result = client.delegate.withdraw_delegate_reward(address)
    assert_code(result, 0)
    delegate_epoch, cumulative_income = get_dividend_information(client, address)
    assert delegate_epoch == 1, "ErrMsg: Last time first delegate epoch {}".format(delegate_epoch)
    assert cumulative_income == 0, "ErrMsg: Last time first cumulative income {}".format(cumulative_income)
    last_delegate_epoch, delegate_total, delegate_total_hes, delegate_reward_total = get_delegate_relevant_amount_and_epoch(
        client)
    assert last_delegate_epoch == 1, "ErrMsg: Last time delegate epoch {}".format(last_delegate_epoch)
    assert delegate_total == 0, "The total number of effective commissioned nodes: {}".format(
        delegate_total)
    assert delegate_total_hes == von_amount(delegate_amount, 2), "The total number of inactive nodes commissioned: {}".format(
        delegate_total_hes)
    assert delegate_reward_total == 0, "Total delegated rewards currently issued by the candidate: {}".format(
        delegate_reward_total)
    result = client.ppos.getDelegateReward(address)
    assert result['Ret'][0]['reward'] == 0, "ErrMsg:Receive dividends {}".format(result['Ret'][0]['reward'])


@pytest.mark.P1
def test_EI_BC_061(client_new_node):
    """
    自由、锁仓金额同时委托首次领取分红（生效期N）
    :param client_new_node:
    :return:
    """
    client = client_new_node
    economic = client.economic
    node = client.node
    address = create_restricting_plan(client)
    log.info("Create delegate account：{}".format(address))
    create_staking_node(client)
    log.info("Create pledge node id :{}".format(node.node_id))
    delegate_amount = von_amount(economic.delegate_limit, 10)
    result = client.delegate.delegate(0, address, amount=delegate_amount)
    assert_code(result, 0)
    result = client.delegate.delegate(1, address, amount=delegate_amount)
    assert_code(result, 0)
    economic.wait_settlement_blocknum(node)
    log.info("Current settlement block height：{}".format(node.eth.blockNumber))
    result = client.delegate.withdraw_delegate_reward(address)
    assert_code(result, 0)
    delegate_epoch, cumulative_income = get_dividend_information(client, address)
    assert delegate_epoch == 2, "ErrMsg: Last time first delegate epoch {}".format(delegate_epoch)
    assert cumulative_income == 0, "ErrMsg: Last time first cumulative income {}".format(cumulative_income)
    last_delegate_epoch, delegate_total, delegate_total_hes, delegate_reward_total = get_delegate_relevant_amount_and_epoch(
        client)
    assert last_delegate_epoch == 1, "ErrMsg: Last time delegate epoch {}".format(last_delegate_epoch)
    assert delegate_total == von_amount(delegate_amount, 2), "The total number of effective commissioned nodes: {}".format(
        delegate_total)
    assert delegate_total_hes == 0, "The total number of inactive nodes commissioned: {}".format(
        delegate_total_hes)
    assert delegate_reward_total == 0, "Total delegated rewards currently issued by the candidate: {}".format(
        delegate_reward_total)
    result = client.ppos.getDelegateReward(address)
    assert result['Ret'][0]['reward'] == 0, "ErrMsg:Receive dividends {}".format(result['Ret'][0]['reward'])


@pytest.mark.P1
def test_EI_BC_062(client_new_node):
    """
    跨结算期自由、锁仓金额同时委托首次领取分红
    :param client_new_node:
    :return:
    """
    client = client_new_node
    economic = client.economic
    node = client.node
    address = create_restricting_plan(client)
    log.info("Create delegate account：{}".format(address))
    create_staking_node(client)
    log.info("Create pledge node id :{}".format(node.node_id))
    delegate_amount = von_amount(economic.delegate_limit, 10)
    result = client.delegate.delegate(0, address, amount=delegate_amount)
    assert_code(result, 0)
    result = client.delegate.delegate(1, address, amount=delegate_amount)
    assert_code(result, 0)

    result = client.delegate.withdraw_delegate_reward(address)
    assert_code(result, 0)
    delegate_epoch, cumulative_income = get_dividend_information(client, address)
    assert delegate_epoch == 1, "ErrMsg: Last time first delegate epoch {}".format(delegate_epoch)
    assert cumulative_income == 0, "ErrMsg: Last time first cumulative income {}".format(cumulative_income)
    last_delegate_epoch, delegate_total, delegate_total_hes, delegate_reward_total = get_delegate_relevant_amount_and_epoch(
        client)
    assert last_delegate_epoch == 1, "ErrMsg: Last time delegate epoch {}".format(last_delegate_epoch)
    assert delegate_total == 0, "The total number of effective commissioned nodes: {}".format(
        delegate_total)
    assert delegate_total_hes == von_amount(delegate_amount, 2), "The total number of inactive nodes commissioned: {}".format(
        delegate_total_hes)
    assert delegate_reward_total == 0, "Total delegated rewards currently issued by the candidate: {}".format(
        delegate_reward_total)
    result = client.ppos.getDelegateReward(address)
    assert result['Ret'][0]['reward'] == 0, "ErrMsg:Receive dividends {}".format(result['Ret'][0]['reward'])
    economic.wait_settlement_blocknum(node)
    log.info("Current settlement block height：{}".format(node.eth.blockNumber))
    result = client.delegate.withdraw_delegate_reward(address)
    assert_code(result, 0)
    delegate_epoch, cumulative_income = get_dividend_information(client, address)
    assert delegate_epoch == 2, "ErrMsg: Last time first delegate epoch {}".format(delegate_epoch)
    assert cumulative_income == 0, "ErrMsg: Last time first cumulative income {}".format(cumulative_income)
    last_delegate_epoch, delegate_total, delegate_total_hes, delegate_reward_total = get_delegate_relevant_amount_and_epoch(
        client)
    assert last_delegate_epoch == 1, "ErrMsg: Last time delegate epoch {}".format(last_delegate_epoch)
    assert delegate_total == von_amount(delegate_amount, 2), "The total number of effective commissioned nodes: {}".format(
        delegate_total)
    assert delegate_total_hes == 0, "The total number of inactive nodes commissioned: {}".format(
        delegate_total_hes)
    assert delegate_reward_total == 0, "Total delegated rewards currently issued by the candidate: {}".format(
        delegate_reward_total)
    result = client.ppos.getDelegateReward(address)
    assert result['Ret'][0]['reward'] == 0, "ErrMsg:Receive dividends {}".format(result['Ret'][0]['reward'])
    block_reward, staking_reward = economic.get_current_year_reward(node)
    economic.wait_settlement_blocknum(node)
    log.info("Current settlement block height：{}".format(node.eth.blockNumber))
    commission_award = economic.calculate_delegate_reward(node, block_reward, staking_reward)
    result = client.delegate.withdraw_delegate_reward(address)
    assert_code(result, 0)
    delegate_epoch, cumulative_income = get_dividend_information(client, address)
    assert delegate_epoch == 3, "ErrMsg: Last time first delegate epoch {}".format(delegate_epoch)
    assert cumulative_income == 0, "ErrMsg: Last time first cumulative income {}".format(cumulative_income)
    last_delegate_epoch, delegate_total, delegate_total_hes, delegate_reward_total = get_delegate_relevant_amount_and_epoch(
        client)
    assert last_delegate_epoch == 1, "ErrMsg: Last time delegate epoch {}".format(last_delegate_epoch)
    assert delegate_total == delegate_amount, "The total number of effective commissioned nodes: {}".format(
        delegate_total)
    assert delegate_total_hes == 0, "The total number of inactive nodes commissioned: {}".format(
        delegate_total_hes)
    assert delegate_reward_total == commission_award, "Total delegated rewards currently issued by the candidate: {}".format(
        delegate_reward_total)
    result = client.ppos.getDelegateReward(address)
    assert result['Ret'][0]['reward'] == 0, "ErrMsg:Receive dividends {}".format(result['Ret'][0]['reward'])


@pytest.mark.P1
@pytest.mark.parametrize('delegate_type', [0, 1])
def test_EI_BC_063_064(clients_new_node, delegate_type):
    """
    自由金额多节点被委托，领取单个节点分红（未生效期N）
    锁仓金额多节点被委托，领取单个节点分红（未生效期N）
    :param clients_new_node:
    :return:
    """
    first_client = clients_new_node[0]
    second_client = clients_new_node[1]
    first_economic = first_client.economic
    first_node = first_client.node
    second_economic = second_client.economic
    second_node = second_client.node
    if delegate_type == 0:
        address, _ = first_economic.account.generate_account(first_node.web3,
                                                             von_amount(first_economic.delegate_limit, 100))
    else:
        address = create_restricting_plan(first_client)
    log.info("Create delegate account：{}".format(address))
    create_stakings_node(clients_new_node)
    log.info("Create first pledge node id :{}".format(first_node.node_id))
    log.info("Create second pledge node id :{}".format(second_node.node_id))
    delegate_amount = von_amount(first_economic.delegate_limit, 10)
    result = first_client.delegate.delegate(delegate_type, address, amount=delegate_amount)
    assert_code(result, 0)
    log.info("Commissioned successfully, commissioned amount：{}".format(delegate_amount))
    result = first_client.delegate.delegate(delegate_type, address, node_id=second_node.node_id, amount=delegate_amount)
    assert_code(result, 0)
    log.info("Commissioned successfully, commissioned amount：{}".format(delegate_amount))
    result = first_client.delegate.withdraw_delegate_reward(address)
    assert_code(result, 0)
    first_delegate_epoch, first_cumulative_income = get_dividend_information(first_client, address)
    second_delegate_epoch, second_cumulative_income = get_dividend_information(second_client, address)
    assert first_delegate_epoch == 1, "ErrMsg: Last time first delegate epoch {}".format(first_delegate_epoch)
    assert first_cumulative_income == 0, "ErrMsg: Last time first cumulative income {}".format(first_cumulative_income)
    assert second_delegate_epoch == 1, "ErrMsg: Last time first delegate epoch {}".format(second_delegate_epoch)
    assert second_cumulative_income == 0, "ErrMsg: Last time first cumulative income {}".format(second_cumulative_income)
    first_last_delegate_epoch, first_delegate_total, first_delegate_total_hes, first_delegate_reward_total = get_delegate_relevant_amount_and_epoch(first_client)
    second_last_delegate_epoch, second_delegate_total, second_delegate_total_hes, second_delegate_reward_total = get_delegate_relevant_amount_and_epoch(second_client)
    assert first_last_delegate_epoch == 1, "ErrMsg: Last time delegate epoch {}".format(first_last_delegate_epoch)
    assert first_delegate_total == 0, "The total number of effective commissioned nodes: {}".format(first_delegate_total)
    assert first_delegate_total_hes == delegate_amount, "The total number of inactive nodes commissioned: {}".format(
        first_delegate_total_hes)
    assert first_delegate_reward_total == 0, "Total delegated rewards currently issued by the candidate: {}".format(
        first_delegate_reward_total)
    result = first_client.ppos.getDelegateReward(address)
    assert result['Ret'][0]['reward'] == 0, "ErrMsg:Receive dividends {}".format(result['Ret'][0]['reward'])
    assert second_last_delegate_epoch == 1, "ErrMsg: Last time delegate epoch {}".format(second_last_delegate_epoch)
    assert second_delegate_total == 0, "The total number of effective commissioned nodes: {}".format(second_delegate_total)
    assert second_delegate_total_hes == delegate_amount, "The total number of inactive nodes commissioned: {}".format(
        second_delegate_total_hes)
    assert second_delegate_reward_total == 0, "Total delegated rewards currently issued by the candidate: {}".format(
        second_delegate_reward_total)
    result = second_client.ppos.getDelegateReward(address)
    assert result['Ret'][0]['reward'] == 0, "ErrMsg:Receive dividends {}".format(result['Ret'][0]['reward'])


@pytest.mark.P1
@pytest.mark.parametrize('delegate_type', [0, 1])
def test_EI_BC_065_066(clients_new_node, delegate_type):
    """
    自由金额多节点被委托，领取分红（生效期N）
    锁仓金额多节点被委托，领取分红（生效期N）
    :param clients_new_node:
    :return:
    """
    first_client = clients_new_node[0]
    second_client = clients_new_node[1]
    first_economic = first_client.economic
    first_node = first_client.node
    second_economic = second_client.economic
    second_node = second_client.node
    if delegate_type == 0:
        address, _ = first_economic.account.generate_account(first_node.web3,von_amount(first_economic.delegate_limit, 100))
    else:
        address = create_restricting_plan(first_client)
    log.info("Create delegate account：{}".format(address))
    create_stakings_node(clients_new_node)
    log.info("Create first pledge node id :{}".format(first_node.node_id))
    log.info("Create second pledge node id :{}".format(second_node.node_id))
    delegate_amount = von_amount(first_economic.delegate_limit, 10)
    result = first_client.delegate.delegate(delegate_type, address, amount=delegate_amount)
    assert_code(result, 0)
    log.info("Commissioned successfully, commissioned amount：{}".format(delegate_amount))
    result = first_client.delegate.delegate(delegate_type, address, node_id=second_node.node_id, amount=delegate_amount)
    assert_code(result, 0)
    log.info("Commissioned successfully, commissioned amount：{}".format(delegate_amount))
    first_economic.wait_settlement_blocknum(first_node)
    log.info("Current settlement block height：{}".format(first_node.eth.blockNumber))
    result = first_client.delegate.withdraw_delegate_reward(address)
    assert_code(result, 0)
    first_delegate_epoch, first_cumulative_income = get_dividend_information(first_client, address)
    second_delegate_epoch, second_cumulative_income = get_dividend_information(second_client, address)
    assert first_delegate_epoch == 2, "ErrMsg: Last time first delegate epoch {}".format(first_delegate_epoch)
    assert first_cumulative_income == 0, "ErrMsg: Last time first cumulative income {}".format(first_cumulative_income)
    assert second_delegate_epoch == 2, "ErrMsg: Last time first delegate epoch {}".format(second_delegate_epoch)
    assert second_cumulative_income == 0, "ErrMsg: Last time first cumulative income {}".format(
        second_cumulative_income)
    first_last_delegate_epoch, first_delegate_total, first_delegate_total_hes, first_delegate_reward_total = get_delegate_relevant_amount_and_epoch(
        first_client)
    second_last_delegate_epoch, second_delegate_total, second_delegate_total_hes, second_delegate_reward_total = get_delegate_relevant_amount_and_epoch(
        second_client)
    assert first_last_delegate_epoch == 1, "ErrMsg: Last time delegate epoch {}".format(first_last_delegate_epoch)
    assert first_delegate_total == delegate_amount, "The total number of effective commissioned nodes: {}".format(
        first_delegate_total)
    assert first_delegate_total_hes == 0, "The total number of inactive nodes commissioned: {}".format(
        first_delegate_total_hes)
    assert first_delegate_reward_total == 0, "Total delegated rewards currently issued by the candidate: {}".format(
        first_delegate_reward_total)
    result = first_client.ppos.getDelegateReward(address)
    assert result['Ret'][0]['reward'] == 0, "ErrMsg:Receive dividends {}".format(result['Ret'][0]['reward'])
    assert second_last_delegate_epoch == 1, "ErrMsg: Last time delegate epoch {}".format(second_last_delegate_epoch)
    assert second_delegate_total == delegate_amount, "The total number of effective commissioned nodes: {}".format(
        second_delegate_total)
    assert second_delegate_total_hes == 0, "The total number of inactive nodes commissioned: {}".format(
        second_delegate_total_hes)
    assert second_delegate_reward_total == 0, "Total delegated rewards currently issued by the candidate: {}".format(
        second_delegate_reward_total)
    result = second_client.ppos.getDelegateReward(address)
    assert result['Ret'][0]['reward'] == 0, "ErrMsg:Receive dividends {}".format(result['Ret'][0]['reward'])


@pytest.mark.P1
@pytest.mark.parametrize('delegate_type', [0, 1])
def test_EI_BC_067_068(clients_new_node, delegate_type):
    """
    自由金额跨结算期多节点被委托，领取分红
    锁仓金额跨结算期多节点被委托，领取分红
    :param clients_new_node:
    :return:
    """
    first_client = clients_new_node[0]
    second_client = clients_new_node[1]
    first_economic = first_client.economic
    first_node = first_client.node
    second_economic = second_client.economic
    second_node = second_client.node
    if delegate_type == 0:
        address, _ = first_economic.account.generate_account(first_node.web3,
                                                             von_amount(first_economic.delegate_limit, 100))
    else:
        address = create_restricting_plan(first_client)
    log.info("Create delegate account：{}".format(address))
    create_stakings_node(clients_new_node)
    log.info("Create first pledge node id :{}".format(first_node.node_id))
    log.info("Create second pledge node id :{}".format(second_node.node_id))
    delegate_amount = von_amount(first_economic.delegate_limit, 10)
    result = first_client.delegate.delegate(delegate_type, address, amount=delegate_amount)
    assert_code(result, 0)
    log.info("Commissioned successfully, commissioned amount：{}".format(delegate_amount))
    result = first_client.delegate.delegate(delegate_type, address, node_id=second_node.node_id, amount=delegate_amount)
    assert_code(result, 0)
    log.info("Commissioned successfully, commissioned amount：{}".format(delegate_amount))
    first_economic.wait_settlement_blocknum(first_node)
    log.info("Current settlement block height：{}".format(first_node.eth.blockNumber))
    block_reward, staking_reward = first_economic.get_current_year_reward(first_node)
    first_economic.wait_settlement_blocknum(first_node)
    log.info("Current settlement block height：{}".format(first_node.eth.blockNumber))
    first_commission_award = first_economic.calculate_delegate_reward(first_node, block_reward, staking_reward)
    second_commission_award = second_economic.calculate_delegate_reward(second_node, block_reward, staking_reward)
    result = first_client.delegate.withdraw_delegate_reward(address)
    assert_code(result, 0)
    first_delegate_epoch, first_cumulative_income = get_dividend_information(first_client, address)
    second_delegate_epoch, second_cumulative_income = get_dividend_information(second_client, address)
    assert first_delegate_epoch == 3, "ErrMsg: Last time first delegate epoch {}".format(first_delegate_epoch)
    assert first_cumulative_income == 0, "ErrMsg: Last time first cumulative income {}".format(first_cumulative_income)
    assert second_delegate_epoch == 3, "ErrMsg: Last time first delegate epoch {}".format(second_delegate_epoch)
    assert second_cumulative_income == 0, "ErrMsg: Last time first cumulative income {}".format(
        second_cumulative_income)
    first_last_delegate_epoch, first_delegate_total, first_delegate_total_hes, first_delegate_reward_total = get_delegate_relevant_amount_and_epoch(
        first_client)
    second_last_delegate_epoch, second_delegate_total, second_delegate_total_hes, second_delegate_reward_total = get_delegate_relevant_amount_and_epoch(
        second_client)
    assert first_last_delegate_epoch == 1, "ErrMsg: Last time delegate epoch {}".format(first_last_delegate_epoch)
    assert first_delegate_total == delegate_amount, "The total number of effective commissioned nodes: {}".format(
        first_delegate_total)
    assert first_delegate_total_hes == 0, "The total number of inactive nodes commissioned: {}".format(
        first_delegate_total_hes)
    assert first_delegate_reward_total == first_commission_award, "Total delegated rewards currently issued by the candidate: {}".format(
        first_delegate_reward_total)
    result = first_client.ppos.getDelegateReward(address)
    assert result['Ret'][0]['reward'] == 0, "ErrMsg:Receive dividends {}".format(result['Ret'][0]['reward'])
    assert second_last_delegate_epoch == 1, "ErrMsg: Last time delegate epoch {}".format(second_last_delegate_epoch)
    assert second_delegate_total == delegate_amount, "The total number of effective commissioned nodes: {}".format(
        second_delegate_total)
    assert second_delegate_total_hes == 0, "The total number of inactive nodes commissioned: {}".format(
        second_delegate_total_hes)
    assert second_delegate_reward_total == second_commission_award, "Total delegated rewards currently issued by the candidate: {}".format(
        second_delegate_reward_total)
    result = second_client.ppos.getDelegateReward(address)
    assert result['Ret'][0]['reward'] == 0, "ErrMsg:Receive dividends {}".format(result['Ret'][0]['reward'])


@pytest.mark.P1
def test_EI_BC_069(client_new_node):
    """
    节点退出中领取分红，验证待领取的委托收益
    :param client_new_node:
    :return:
    """
    client = client_new_node
    economic = client.economic
    node = client.node
    address = create_restricting_plan(client)
    log.info("Create delegate account：{}".format(address))
    staking_address = create_staking_node(client)
    log.info("Create pledge node id :{}".format(node.node_id))
    delegate_amount = von_amount(economic.delegate_limit, 10)
    result = client.delegate.delegate(0, address, amount=delegate_amount)
    assert_code(result, 0)
    economic.wait_settlement_blocknum(node)
    log.info("Current settlement block height：{}".format(node.eth.blockNumber))
    result = client.staking.withdrew_staking(staking_address)
    assert_code(result, 0)
    result = client.delegate.withdraw_delegate_reward(address)
    assert_code(result, 0)
    delegate_epoch, cumulative_income = get_dividend_information(client, address)
    assert delegate_epoch == 2, "ErrMsg: Last time first delegate epoch {}".format(delegate_epoch)
    assert cumulative_income == 0, "ErrMsg: Last time first cumulative income {}".format(cumulative_income)
    last_delegate_epoch, delegate_total, delegate_total_hes, delegate_reward_total = get_delegate_relevant_amount_and_epoch(
        client)
    assert last_delegate_epoch == 1, "ErrMsg: Last time delegate epoch {}".format(last_delegate_epoch)
    assert delegate_total == delegate_amount, "The total number of effective commissioned nodes: {}".format(
        delegate_total)
    assert delegate_total_hes == 0, "The total number of inactive nodes commissioned: {}".format(
        delegate_total_hes)
    assert delegate_reward_total == 0, "Total delegated rewards currently issued by the candidate: {}".format(
        delegate_reward_total)
    result = client.ppos.getDelegateReward(address)
    assert result['Ret'][0]['reward'] == 0, "ErrMsg:Receive dividends {}".format(result['Ret'][0]['reward'])
    block_reward, staking_reward = economic.get_current_year_reward(node)
    economic.wait_settlement_blocknum(node)
    log.info("Current settlement block height：{}".format(node.eth.blockNumber))
    commission_award = economic.calculate_delegate_reward(node, block_reward, staking_reward)
    result = client.delegate.withdraw_delegate_reward(address)
    assert_code(result, 0)
    delegate_epoch, cumulative_income = get_dividend_information(client, address)
    assert delegate_epoch == 3, "ErrMsg: Last time first delegate epoch {}".format(delegate_epoch)
    assert cumulative_income == 0, "ErrMsg: Last time first cumulative income {}".format(cumulative_income)
    last_delegate_epoch, delegate_total, delegate_total_hes, delegate_reward_total = get_delegate_relevant_amount_and_epoch(
        client)
    assert last_delegate_epoch == 1, "ErrMsg: Last time delegate epoch {}".format(last_delegate_epoch)
    assert delegate_total == delegate_amount, "The total number of effective commissioned nodes: {}".format(
        delegate_total)
    assert delegate_total_hes == 0, "The total number of inactive nodes commissioned: {}".format(
        delegate_total_hes)
    assert delegate_reward_total == commission_award, "Total delegated rewards currently issued by the candidate: {}".format(
        delegate_reward_total)
    result = client.ppos.getDelegateReward(address)
    assert result['Ret'][0]['reward'] == 0, "ErrMsg:Receive dividends {}".format(result['Ret'][0]['reward'])


@pytest.mark.P2
def test_EI_BC_070(client_new_node):
    """
    节点已退出质押领取分红，验证待领取的委托收益
    :param client_new_node:
    :return:
    """
    client = client_new_node
    economic = client.economic
    node = client.node
    address = create_restricting_plan(client)
    log.info("Create delegate account：{}".format(address))
    staking_address = create_staking_node(client)
    log.info("Create pledge node id :{}".format(node.node_id))
    delegate_amount = von_amount(economic.delegate_limit, 10)
    result = client.delegate.delegate(0, address, amount=delegate_amount)
    assert_code(result, 0)
    economic.wait_settlement_blocknum(node)
    log.info("Current settlement block height：{}".format(node.eth.blockNumber))
    result = client.staking.withdrew_staking(staking_address)
    assert_code(result, 0)
    result = client.delegate.withdraw_delegate_reward(address)
    assert_code(result, 0)
    delegate_epoch, cumulative_income = get_dividend_information(client, address)
    assert delegate_epoch == 2, "ErrMsg: Last time first delegate epoch {}".format(delegate_epoch)
    assert cumulative_income == 0, "ErrMsg: Last time first cumulative income {}".format(cumulative_income)
    last_delegate_epoch, delegate_total, delegate_total_hes, delegate_reward_total = get_delegate_relevant_amount_and_epoch(
        client)
    assert last_delegate_epoch == 1, "ErrMsg: Last time delegate epoch {}".format(last_delegate_epoch)
    assert delegate_total == delegate_amount, "The total number of effective commissioned nodes: {}".format(
        delegate_total)
    assert delegate_total_hes == 0, "The total number of inactive nodes commissioned: {}".format(
        delegate_total_hes)
    assert delegate_reward_total == 0, "Total delegated rewards currently issued by the candidate: {}".format(
        delegate_reward_total)
    result = client.ppos.getDelegateReward(address)
    assert result['Ret'][0]['reward'] == 0, "ErrMsg:Receive dividends {}".format(result['Ret'][0]['reward'])
    block_reward, staking_reward = economic.get_current_year_reward(node)
    economic.wait_settlement_blocknum(node)
    log.info("Current settlement block height：{}".format(node.eth.blockNumber))
    commission_award = economic.calculate_delegate_reward(node, block_reward, staking_reward)
    economic.wait_settlement_blocknum(node, 2)
    log.info("Current settlement block height：{}".format(node.eth.blockNumber))
    delegate_balance = node.eth.getBalance(address)
    log.info("Entrusted account balance： {}".format(delegate_balance))
    result = client.delegate.withdraw_delegate_reward(address)
    assert_code(result, 0)
    delegate_epoch, cumulative_income = get_dividend_information(client, address)
    assert delegate_epoch == 3, "ErrMsg: Last time first delegate epoch {}".format(delegate_epoch)
    assert cumulative_income == 0, "ErrMsg: Last time first cumulative income {}".format(cumulative_income)
    result = client.ppos.getCandidateInfo(node.node_id)
    assert_code(result, 301102)
    result = client.ppos.getDelegateReward(address)
    assert result['Ret'][0]['reward'] == 0, "ErrMsg:Receive dividends {}".format(result['Ret'][0]['reward'])


@pytest.mark.P1
def test_EI_BC_071(client_new_node):
    """
    节点被惩罚领取分红，验证待领取的委托收益
    :param client_new_node:
    :return:
    """
    client = client_new_node
    economic = client.economic
    node = client.node
    address = create_restricting_plan(client)
    log.info("Create delegate account：{}".format(address))
    staking_address = create_staking_node(client)
    log.info("Create pledge node id :{}".format(node.node_id))
    delegate_amount = von_amount(economic.delegate_limit, 10)
    result = client.delegate.delegate(0, address, amount=delegate_amount)
    assert_code(result, 0)
    economic.wait_settlement_blocknum(node)
    log.info("Current settlement block height：{}".format(node.eth.blockNumber))
    result = client.staking.withdrew_staking(staking_address)
    assert_code(result, 0)
    result = client.delegate.withdraw_delegate_reward(address)
    assert_code(result, 0)
    delegate_epoch, cumulative_income = get_dividend_information(client, address)
    assert delegate_epoch == 2, "ErrMsg: Last time first delegate epoch {}".format(delegate_epoch)
    assert cumulative_income == 0, "ErrMsg: Last time first cumulative income {}".format(cumulative_income)
    last_delegate_epoch, delegate_total, delegate_total_hes, delegate_reward_total = get_delegate_relevant_amount_and_epoch(
        client)
    assert last_delegate_epoch == 1, "ErrMsg: Last time delegate epoch {}".format(last_delegate_epoch)
    assert delegate_total == delegate_amount, "The total number of effective commissioned nodes: {}".format(
        delegate_total)
    assert delegate_total_hes == 0, "The total number of inactive nodes commissioned: {}".format(
        delegate_total_hes)
    assert delegate_reward_total == 0, "Total delegated rewards currently issued by the candidate: {}".format(
        delegate_reward_total)
    result = client.ppos.getDelegateReward(address)
    assert result['Ret'][0]['reward'] == 0, "ErrMsg:Receive dividends {}".format(result['Ret'][0]['reward'])
    block_reward, staking_reward = economic.get_current_year_reward(node)
    economic.wait_settlement_blocknum(node)
    log.info("Current settlement block height：{}".format(node.eth.blockNumber))
    commission_award = economic.calculate_delegate_reward(node, block_reward, staking_reward)
    result = client.delegate.withdraw_delegate_reward(address)
    assert_code(result, 0)
    delegate_epoch, cumulative_income = get_dividend_information(client, address)
    assert delegate_epoch == 3, "ErrMsg: Last time first delegate epoch {}".format(delegate_epoch)
    assert cumulative_income == 0, "ErrMsg: Last time first cumulative income {}".format(cumulative_income)
    last_delegate_epoch, delegate_total, delegate_total_hes, delegate_reward_total = get_delegate_relevant_amount_and_epoch(
        client)
    assert last_delegate_epoch == 1, "ErrMsg: Last time delegate epoch {}".format(last_delegate_epoch)
    assert delegate_total == delegate_amount, "The total number of effective commissioned nodes: {}".format(
        delegate_total)
    assert delegate_total_hes == 0, "The total number of inactive nodes commissioned: {}".format(
        delegate_total_hes)
    assert delegate_reward_total == commission_award, "Total delegated rewards currently issued by the candidate: {}".format(
        delegate_reward_total)
    result = client.ppos.getDelegateReward(address)
    assert result['Ret'][0]['reward'] == 0, "ErrMsg:Receive dividends {}".format(result['Ret'][0]['reward'])



@pytest.mark.P0
def test_EI_BC_075(client_new_node):
    """
    不同生效期N二次委托、赎回部分委托、领取分红
    :param client_new_node:
    :return:
    """
    client = client_new_node
    economic = client.economic
    node = client.node
    address, _ = economic.account.generate_account(node.web3, von_amount(economic.delegate_limit, 100))
    log.info("Create delegate account：{}".format(address))
    create_staking_node(client)
    log.info("Create pledge node id :{}".format(node.node_id))
    delegate_amount = von_amount(economic.delegate_limit, 10)
    result = client.delegate.delegate(0, address, amount=delegate_amount)
    assert_code(result, 0)
    log.info("Commissioned successfully, commissioned amount：{}".format(economic.delegate_limit))
    economic.wait_settlement_blocknum(node)
    log.info("Current settlement block height：{}".format(node.eth.blockNumber))
    result = client.delegate.delegate(0, address, amount=delegate_amount)
    assert_code(result, 0)
    delegate_epoch, cumulative_income = get_dividend_information(client, address)
    assert delegate_epoch == 2, "ErrMsg: Last time delegate epoch {}".format(delegate_epoch)
    assert cumulative_income == 0, "ErrMsg: Last time cumulative income {}".format(cumulative_income)
    last_delegate_epoch, delegate_total, delegate_total_hes, delegate_reward_total = get_delegate_relevant_amount_and_epoch(
        client)
    assert last_delegate_epoch == 2, "ErrMsg: Last time delegate epoch {}".format(last_delegate_epoch)
    assert delegate_total == delegate_amount, "The total number of effective commissioned nodes: {}".format(delegate_total)
    assert delegate_total_hes == delegate_amount, "The total number of inactive nodes commissioned: {}".format(
        delegate_total_hes)
    assert delegate_reward_total == 0, "Total delegated rewards currently issued by the candidate: {}".format(
        delegate_reward_total)
    block_reward, staking_reward = economic.get_current_year_reward(node)
    economic.wait_settlement_blocknum(node)
    log.info("Current settlement block height：{}".format(node.eth.blockNumber))
    commission_award = economic.calculate_delegate_reward(node, block_reward, staking_reward)
    result = client.ppos.getCandidateInfo(node.node_id)
    blocknum = result['Ret']['StakingBlockNum']
    client.delegate.withdrew_delegate(blocknum, address)
    delegate_epoch, cumulative_income = get_dividend_information(client, address)
    assert delegate_epoch == 3, "ErrMsg: Last time delegate epoch {}".format(delegate_epoch)
    assert cumulative_income == commission_award, "ErrMsg: Last time cumulative income {}".format(cumulative_income)
    last_delegate_epoch, delegate_total, delegate_total_hes, delegate_reward_total = get_delegate_relevant_amount_and_epoch(
        client)
    assert last_delegate_epoch == 2, "ErrMsg: Last time delegate epoch {}".format(last_delegate_epoch)
    assert delegate_total == von_amount(delegate_amount, 2) - economic.delegate_limit, "The total number of effective commissioned nodes: {}".format(
        delegate_total)
    assert delegate_total_hes == 0, "The total number of inactive nodes commissioned: {}".format(
        delegate_total_hes)
    assert delegate_reward_total == commission_award, "Total delegated rewards currently issued by the candidate: {}".format(
        delegate_reward_total)
    block_reward, staking_reward = economic.get_current_year_reward(node)
    economic.wait_settlement_blocknum(node)
    log.info("Current settlement block height：{}".format(node.eth.blockNumber))
    commission_award = economic.calculate_delegate_reward(node, block_reward, staking_reward)
    result = client.delegate.withdraw_delegate_reward(address)
    assert_code(result, 0)
    delegate_epoch, cumulative_income = get_dividend_information(client, address)
    assert delegate_epoch == 4, "ErrMsg: Last time delegate epoch {}".format(delegate_epoch)
    assert cumulative_income == 0, "ErrMsg: Last time cumulative income {}".format(cumulative_income)
    last_delegate_epoch, delegate_total, delegate_total_hes, delegate_reward_total = get_delegate_relevant_amount_and_epoch(
        client)
    assert last_delegate_epoch == 2, "ErrMsg: Last time delegate epoch {}".format(last_delegate_epoch)
    assert delegate_total == von_amount(delegate_amount, 2) - economic.delegate_limit, "The total number of effective commissioned nodes: {}".format(
        delegate_total)
    assert delegate_total_hes == 0, "The total number of inactive nodes commissioned: {}".format(
        delegate_total_hes)
    assert delegate_reward_total == commission_award, "Total delegated rewards currently issued by the candidate: {}".format(
        delegate_reward_total)
    result = client.ppos.getDelegateReward(address)
    assert result['Ret'][0]['reward'] == 0, "ErrMsg: Dividends currently available {}".format(
        result['Ret'][0]['reward'])































