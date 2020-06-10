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
    get_governable_parameter_value, get_the_dynamic_parameter_gas_fee, get_getDelegateReward_gas_fee


def create_staking_node(client):
    """
    创建一个自由质押节点
    :param client:
    :return:
    """
    economic = client.economic
    node = client.node
    staking_address, _ = economic.account.generate_account(node.web3, von_amount(economic.create_staking_limit, 3))
    benifit_address, _ = economic.account.generate_account(node.web3)
    result = client.staking.create_staking(0, staking_address, staking_address,
                                           amount=von_amount(economic.create_staking_limit, 2), reward_per=1000)
    assert_code(result, 0)
    print(staking_address)
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
    first_staking_address, _ = first_economic.account.generate_account(first_node.web3,
                                                                       von_amount(first_economic.create_staking_limit,
                                                                                  3))
    first_benifit_address, _ = first_economic.account.generate_account(first_node.web3)
    result = first_client.staking.create_staking(0, first_benifit_address, first_staking_address,
                                                 amount=von_amount(first_economic.create_staking_limit, 2),
                                                 reward_per=1000)
    assert_code(result, 0)
    second_staking_address, _ = second_economic.account.generate_account(second_node.web3, von_amount(
        second_economic.create_staking_limit, 3))
    second_benifit_address, _ = second_economic.account.generate_account(second_node.web3)
    result = second_client.staking.create_staking(0, second_benifit_address, second_staking_address,
                                                  amount=von_amount(second_economic.create_staking_limit, 2),
                                                  reward_per=2000)
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


def get_dividend_information(client, node_id, address):
    """
    获取分红信息
    :param client:
    :return:
    """
    result = client.ppos.getCandidateInfo(node_id)
    blocknum = result['Ret']['StakingBlockNum']
    result = client.ppos.getDelegateInfo(blocknum, address, node_id)
    log.info("Commission information：{}".format(result))
    info = result['Ret']
    delegate_epoch = info['DelegateEpoch']
    cumulative_income = info['CumulativeIncome']
    return delegate_epoch, cumulative_income


def get_delegate_relevant_amount_and_epoch(client, node_id):
    result = client.ppos.getCandidateInfo(node_id)
    log.info(result)
    log.info('Current pledged node pledge information：{}'.format(result))
    last_delegate_epoch = result['Ret']['DelegateEpoch']
    delegate_total = result['Ret']['DelegateTotal']
    delegate_total_hes = result['Ret']['DelegateTotalHes']
    delegate_reward_total = result['Ret']['DelegateRewardTotal']
    return last_delegate_epoch, delegate_total, delegate_total_hes, delegate_reward_total


@pytest.mark.P1
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
    delegate_epoch, cumulative_income = get_dividend_information(client, node.node_id, delegate_address)
    assert delegate_epoch == 1, "ErrMsg: Last time delegate epoch {}".format(delegate_epoch)
    assert cumulative_income == 0, "ErrMsg: Last time cumulative income {}".format(cumulative_income)
    economic.wait_consensus_blocknum(node)
    # initiate redemption
    result = client.ppos.getCandidateInfo(client.node.node_id)
    blocknum = result['Ret']['StakingBlockNum']
    result = client.delegate.withdrew_delegate(blocknum, delegate_address)
    assert_code(result, 0)
    delegate_epoch, cumulative_income = get_dividend_information(client, node.node_id, delegate_address)
    assert delegate_epoch == 1, "ErrMsg: Last time delegate epoch {}".format(delegate_epoch)
    assert cumulative_income == 0, "ErrMsg: Last time cumulative income {}".format(cumulative_income)
    time.sleep(3)

    # receive dividends
    result = client.ppos.getDelegateReward(delegate_address)
    log.info("result:{}".format(result))
    assert result['Ret'][0]['reward'] == 0, "ErrMsg: Withdraw commission award {}".format(result['Ret'][0]['reward'])
    result = client.delegate.withdraw_delegate_reward(delegate_address)
    assert_code(result, 0)
    delegate_epoch, cumulative_income = get_dividend_information(client, node.node_id, delegate_address)
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
    second_node = second_client.node
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
    first_delegate_epoch, first_cumulative_income = get_dividend_information(first_client, first_node.node_id, address)
    second_delegate_epoch, second_cumulative_income = get_dividend_information(first_client, second_node.node_id,
                                                                               address)
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
    first_delegate_epoch, first_cumulative_income = get_dividend_information(first_client, first_node.node_id, address)
    second_delegate_epoch, second_cumulative_income = get_dividend_information(first_client, second_node.node_id,
                                                                               address)
    assert first_delegate_epoch == 1, "ErrMsg: Last time delegate epoch {}".format(first_delegate_epoch)
    assert first_cumulative_income == 0, "ErrMsg: Last time cumulative income {}".format(first_cumulative_income)
    assert second_delegate_epoch == 1, "ErrMsg: Last time delegate epoch {}".format(second_delegate_epoch)
    assert second_cumulative_income == 0, "ErrMsg: Last time cumulative income {}".format(second_cumulative_income)

    # receive dividends
    result = first_client.ppos.getDelegateReward(address)
    assert result['Ret'][0]['reward'] == 0, "ErrMsg: Withdraw commission award {}".format(result['Ret'][0]['reward'])
    result = first_client.delegate.withdraw_delegate_reward(address)
    assert_code(result, 0)
    first_delegate_epoch, first_cumulative_income = get_dividend_information(first_client, first_node.node_id, address)
    second_delegate_epoch, second_cumulative_income = get_dividend_information(first_client, second_node.node_id,
                                                                               address)
    assert first_delegate_epoch == 1, "ErrMsg: Last time delegate epoch {}".format(first_delegate_epoch)
    assert first_cumulative_income == 0, "ErrMsg: Last time cumulative income {}".format(first_cumulative_income)
    assert second_delegate_epoch == 1, "ErrMsg: Last time delegate epoch {}".format(second_delegate_epoch)
    assert second_cumulative_income == 0, "ErrMsg: Last time cumulative income {}".format(second_cumulative_income)


@pytest.mark.P0
@pytest.mark.parametrize('delegate_type', [0, 1])
def test_EI_BC_003_007(client_new_node, delegate_type, reset_environment):
    """
    自由金额跨周期追加委托，验证待领取的委托收益（单节点）
    锁仓金额跨周期追加委托，验证待领取的委托收益（单节点）
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
    result = client.delegate.delegate(delegate_type, address, amount=delegate_amount)
    assert_code(result, 0)
    log.info("Commissioned successfully, commissioned amount：{}".format(delegate_amount))
    delegate_epoch, cumulative_income = get_dividend_information(client, node.node_id, address)
    assert delegate_epoch == 1, "ErrMsg: Last time delegate epoch {}".format(delegate_epoch)
    assert cumulative_income == 0, "ErrMsg: Last time cumulative income {}".format(cumulative_income)
    last_delegate_epoch, delegate_total, delegate_total_hes, delegate_reward_total = get_delegate_relevant_amount_and_epoch(
        client, node.node_id)
    assert last_delegate_epoch == 1, "ErrMsg: Last time delegate epoch {}".format(last_delegate_epoch)
    assert delegate_total == 0, "The total number of effective commissioned nodes: {}".format(delegate_total)
    assert delegate_total_hes == delegate_amount, "The total number of inactive nodes commissioned: {}".format(
        delegate_total_hes)
    assert delegate_reward_total == 0, "Total delegated rewards currently issued by the candidate: {}".format(
        delegate_reward_total)
    economic.wait_settlement_blocknum(node)
    log.info("Current settlement block height：{}".format(node.eth.blockNumber))
    result = client.delegate.delegate(delegate_type, address, amount=delegate_amount)
    assert_code(result, 0)
    log.info("Commissioned successfully, commissioned amount：{}".format(von_amount(delegate_amount, 2)))
    delegate_epoch, cumulative_income = get_dividend_information(client, node.node_id, address)
    assert delegate_epoch == 2, "ErrMsg: Last time delegate epoch {}".format(delegate_epoch)
    assert cumulative_income == 0, "ErrMsg: Last time cumulative income {}".format(cumulative_income)
    last_delegate_epoch, delegate_total, delegate_total_hes, delegate_reward_total = get_delegate_relevant_amount_and_epoch(
        client, node.node_id)
    assert last_delegate_epoch == 2, "ErrMsg: Last time delegate epoch {}".format(last_delegate_epoch)
    assert delegate_total == delegate_amount, "The total number of effective commissioned nodes: {}".format(
        delegate_total)
    assert delegate_total_hes == delegate_amount, "The total number of inactive nodes commissioned: {}".format(
        delegate_total_hes)
    assert delegate_reward_total == 0, "Total delegated rewards currently issued by the candidate: {}".format(
        delegate_reward_total)
    block_reward, staking_reward = economic.get_current_year_reward(node)
    log.info("The current node block reward: {} Pledge reward: {}".format(block_reward, staking_reward))
    economic.wait_settlement_blocknum(node)
    log.info("Current settlement block height：{}".format(node.eth.blockNumber))
    commission_award = economic.calculate_delegate_reward(node, block_reward, staking_reward)
    current_commission_award = economic.delegate_cumulative_income(node, block_reward, staking_reward, delegate_amount,
                                                                   delegate_amount)
    result = client.delegate.delegate(delegate_type, address, amount=delegate_amount)
    assert_code(result, 0)
    log.info("Commissioned successfully, commissioned amount：{}".format(von_amount(economic.delegate_limit, 2)))
    delegate_epoch, cumulative_income = get_dividend_information(client, node.node_id, address)
    assert delegate_epoch == 3, "ErrMsg: Last time delegate epoch {}".format(delegate_epoch)
    assert cumulative_income == current_commission_award, "ErrMsg: Last time cumulative income {}".format(
        cumulative_income)
    last_delegate_epoch, delegate_total, delegate_total_hes, delegate_reward_total = get_delegate_relevant_amount_and_epoch(
        client, node.node_id)
    assert last_delegate_epoch == 3, "ErrMsg: Last time delegate epoch {}".format(last_delegate_epoch)
    assert delegate_total == von_amount(delegate_amount,
                                        2), "The total number of effective commissioned nodes: {}".format(
        delegate_total)
    assert delegate_total_hes == delegate_amount, "The total number of inactive nodes commissioned: {}".format(
        delegate_total_hes)
    assert delegate_reward_total == commission_award, "Total delegated rewards currently issued by the candidate: {}".format(
        delegate_reward_total)


@pytest.mark.P1
@pytest.mark.parametrize('delegate_type', [0, 1])
def test_EI_BC_004_008(clients_new_node, delegate_type, reset_environment):
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
    first_economic.env.deploy_all()
    address = create_restricting_plan(first_client)
    log.info("Create delegate account：{}".format(address))
    create_stakings_node(clients_new_node)
    log.info("Create first pledge node id :{}".format(first_node.node_id))
    log.info("Create second pledge node id :{}".format(second_node.node_id))
    delegate_amount = von_amount(first_economic.delegate_limit, 10)
    result = first_client.delegate.delegate(delegate_type, address, amount=delegate_amount)
    assert_code(result, 0)
    log.info("Commissioned successfully, commissioned amount：{}".format(first_economic.delegate_limit))
    result = first_client.delegate.delegate(delegate_type, address, amount=delegate_amount, node_id=second_node.node_id)
    assert_code(result, 0)
    log.info("Commissioned successfully, commissioned amount：{}".format(second_economic.delegate_limit))
    first_delegate_epoch, first_cumulative_income = get_dividend_information(first_client, first_node.node_id, address)
    second_delegate_epoch, second_cumulative_income = get_dividend_information(first_client, second_node.node_id,
                                                                               address)
    assert first_delegate_epoch == 1, "ErrMsg: Last time first delegate epoch {}".format(first_delegate_epoch)
    assert first_cumulative_income == 0, "ErrMsg: Last time first cumulative income {}".format(first_cumulative_income)
    assert second_delegate_epoch == 1, "ErrMsg: Last time second delegate epoch {}".format(second_delegate_epoch)
    assert second_cumulative_income == 0, "ErrMsg: Last time second cumulative income {}".format(
        second_cumulative_income)
    first_last_delegate_epoch, first_delegate_total, first_delegate_total_hes, first_delegate_reward_total = get_delegate_relevant_amount_and_epoch(
        first_client, first_node.node_id)
    second_last_delegate_epoch, second_delegate_total, second_delegate_total_hes, second_delegate_reward_total = get_delegate_relevant_amount_and_epoch(
        first_client, second_node.node_id)
    assert first_last_delegate_epoch == 1, "ErrMsg: Last time delegate epoch {}".format(first_last_delegate_epoch)
    assert first_delegate_total == 0, "The total number of effective commissioned nodes: {}".format(
        first_delegate_total)
    assert first_delegate_total_hes == delegate_amount, "The total number of inactive nodes commissioned: {}".format(
        first_delegate_total_hes)
    assert first_delegate_reward_total == 0, "Total delegated rewards currently issued by the candidate: {}".format(
        first_delegate_reward_total)
    assert second_last_delegate_epoch == 1, "ErrMsg: Last time delegate epoch {}".format(second_last_delegate_epoch)
    assert second_delegate_total == 0, "The total number of effective commissioned nodes: {}".format(
        second_delegate_total)
    assert second_delegate_total_hes == delegate_amount, "The total number of inactive nodes commissioned: {}".format(
        second_delegate_total_hes)
    assert second_delegate_reward_total == 0, "Total delegated rewards currently issued by the candidate: {}".format(
        second_delegate_reward_total)
    first_economic.wait_settlement_blocknum(first_node)
    log.info("Current settlement block height：{}".format(first_node.eth.blockNumber))
    result = first_client.delegate.delegate(delegate_type, address, amount=delegate_amount)
    assert_code(result, 0)
    log.info("Commissioned successfully, commissioned amount：{}".format(first_economic.delegate_limit))
    result = first_client.delegate.delegate(delegate_type, address, amount=delegate_amount, node_id=second_node.node_id)
    assert_code(result, 0)
    log.info("Commissioned successfully, commissioned amount：{}".format(second_economic.delegate_limit))
    first_delegate_epoch, first_cumulative_income = get_dividend_information(first_client, first_node.node_id, address)
    second_delegate_epoch, second_cumulative_income = get_dividend_information(first_client, second_node.node_id,
                                                                               address)
    assert first_delegate_epoch == 2, "ErrMsg: Last time first delegate epoch {}".format(first_delegate_epoch)
    assert first_cumulative_income == 0, "ErrMsg: Last time first cumulative income {}".format(first_cumulative_income)
    assert second_delegate_epoch == 2, "ErrMsg: Last time second delegate epoch {}".format(second_delegate_epoch)
    assert second_cumulative_income == 0, "ErrMsg: Last time second cumulative income {}".format(
        second_cumulative_income)
    first_last_delegate_epoch, first_delegate_total, first_delegate_total_hes, first_delegate_reward_total = get_delegate_relevant_amount_and_epoch(
        first_client, first_node.node_id)
    second_last_delegate_epoch, second_delegate_total, second_delegate_total_hes, second_delegate_reward_total = get_delegate_relevant_amount_and_epoch(
        first_client, second_node.node_id)
    assert first_last_delegate_epoch == 2, "ErrMsg: Last time delegate epoch {}".format(first_last_delegate_epoch)
    assert first_delegate_total == delegate_amount, "The total number of effective commissioned nodes: {}".format(
        first_delegate_total)
    assert first_delegate_total_hes == delegate_amount, "The total number of inactive nodes commissioned: {}".format(
        first_delegate_total_hes)
    assert first_delegate_reward_total == 0, "Total delegated rewards currently issued by the candidate: {}".format(
        first_delegate_reward_total)
    assert second_last_delegate_epoch == 2, "ErrMsg: Last time delegate epoch {}".format(second_last_delegate_epoch)
    assert second_delegate_total == delegate_amount, "The total number of effective commissioned nodes: {}".format(
        second_delegate_total)
    assert second_delegate_total_hes == delegate_amount, "The total number of inactive nodes commissioned: {}".format(
        second_delegate_total_hes)
    assert second_delegate_reward_total == 0, "Total delegated rewards currently issued by the candidate: {}".format(
        second_delegate_reward_total)
    block_reward, staking_reward = first_economic.get_current_year_reward(first_node)
    first_economic.wait_settlement_blocknum(first_node)
    log.info("Current settlement block height：{}".format(first_node.eth.blockNumber))
    first_commission_award = first_economic.calculate_delegate_reward(first_node, block_reward, staking_reward)
    first_current_commission_award = first_economic.delegate_cumulative_income(first_node, block_reward, staking_reward,
                                                                               delegate_amount, delegate_amount)
    second_commission_award = second_economic.calculate_delegate_reward(second_node, block_reward, staking_reward)
    second_current_commission_award = first_economic.delegate_cumulative_income(second_node, block_reward,
                                                                                staking_reward, delegate_amount,
                                                                                delegate_amount)
    result = first_client.delegate.delegate(delegate_type, address, amount=delegate_amount)
    assert_code(result, 0)
    log.info("Commissioned successfully, commissioned amount：{}".format(first_economic.delegate_limit))
    result = first_client.delegate.delegate(delegate_type, address, amount=delegate_amount, node_id=second_node.node_id)
    assert_code(result, 0)
    log.info("Current settlement block height：{}".format(first_node.eth.blockNumber))
    first_delegate_epoch, first_cumulative_income = get_dividend_information(first_client, first_node.node_id, address)
    second_delegate_epoch, second_cumulative_income = get_dividend_information(first_client, second_node.node_id,
                                                                               address)
    assert first_delegate_epoch == 3, "ErrMsg: Last time first delegate epoch {}".format(first_delegate_epoch)
    assert first_cumulative_income == first_current_commission_award, "ErrMsg: Last time first cumulative income {}".format(
        first_cumulative_income)
    assert second_delegate_epoch == 3, "ErrMsg: Last time second delegate epoch {}".format(second_delegate_epoch)
    assert second_cumulative_income == second_current_commission_award, "ErrMsg: Last time second cumulative income {}".format(
        second_cumulative_income)
    first_last_delegate_epoch, first_delegate_total, first_delegate_total_hes, first_delegate_reward_total = get_delegate_relevant_amount_and_epoch(
        first_client, first_node.node_id)
    second_last_delegate_epoch, second_delegate_total, second_delegate_total_hes, second_delegate_reward_total = get_delegate_relevant_amount_and_epoch(
        first_client, second_node.node_id)
    assert first_last_delegate_epoch == 3, "ErrMsg: Last time delegate epoch {}".format(first_last_delegate_epoch)
    assert first_delegate_total == von_amount(delegate_amount,
                                              2), "The total number of effective commissioned nodes: {}".format(
        first_delegate_total)
    assert first_delegate_total_hes == delegate_amount, "The total number of inactive nodes commissioned: {}".format(
        first_delegate_total_hes)
    assert first_delegate_reward_total == first_commission_award, "Total delegated rewards currently issued by the candidate: {}".format(
        first_delegate_reward_total)
    assert second_last_delegate_epoch == 3, "ErrMsg: Last time delegate epoch {}".format(second_last_delegate_epoch)
    assert second_delegate_total == von_amount(delegate_amount,
                                               2), "The total number of effective commissioned nodes: {}".format(
        second_delegate_total)
    assert second_delegate_total_hes == delegate_amount, "The total number of inactive nodes commissioned: {}".format(
        second_delegate_total_hes)
    assert second_delegate_reward_total == second_commission_award, "Total delegated rewards currently issued by the candidate: {}".format(
        second_delegate_reward_total)


@pytest.mark.P0
@pytest.mark.parametrize('delegate_type', [0, 1])
def test_EI_BC_010_016(client_new_node, delegate_type, reset_environment):
    """
    自由金额首次部分赎回，验证待领取的委托收益（生效期N）
    锁仓金额首次部分赎回，验证待领取的委托收益（生效期N）
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
    delegate_epoch, cumulative_income = get_dividend_information(client, node.node_id, address)
    assert delegate_epoch == 2, "ErrMsg: Last time first delegate epoch {}".format(delegate_epoch)
    assert cumulative_income == 0, "ErrMsg: Last time first cumulative income {}".format(cumulative_income)
    last_delegate_epoch, delegate_total, delegate_total_hes, delegate_reward_total = get_delegate_relevant_amount_and_epoch(
        client, node.node_id)
    assert last_delegate_epoch == 1, "ErrMsg: Last time delegate epoch {}".format(last_delegate_epoch)
    assert delegate_total == delegate_amount - economic.delegate_limit, "The total number of effective commissioned nodes: {}".format(
        delegate_total)
    assert delegate_total_hes == 0, "The total number of inactive nodes commissioned: {}".format(delegate_total_hes)
    assert delegate_reward_total == 0, "Total delegated rewards currently issued by the candidate: {}".format(
        delegate_reward_total)


@pytest.mark.P0
@pytest.mark.parametrize('delegate_type', [0, 1])
def test_EI_BC_011_074(client_new_node, delegate_type, reset_environment):
    """
    自由金额跨结算周期首次部分赎回，验证待领取的委托收益（生效期N）
    锁仓金额跨结算周期首次部分赎回，验证待领取的委托收益（生效期N）
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
    result = client.delegate.delegate(delegate_type, address, amount=delegate_amount)
    assert_code(result, 0)
    log.info("Commissioned successfully, commissioned amount：{}".format(economic.delegate_limit))
    economic.wait_settlement_blocknum(node)
    block_reward, staking_reward = economic.get_current_year_reward(node)
    economic.wait_settlement_blocknum(node)
    log.info("Current block height：{}".format(node.eth.blockNumber))
    commission_award = economic.calculate_delegate_reward(node, block_reward, staking_reward)
    current_commission_award = economic.delegate_cumulative_income(node, block_reward, staking_reward, delegate_amount,
                                                                   delegate_amount)
    result = client.ppos.getCandidateInfo(client.node.node_id)
    blocknum = result['Ret']['StakingBlockNum']
    result = client.delegate.withdrew_delegate(blocknum, address)
    assert_code(result, 0)
    delegate_epoch, cumulative_income = get_dividend_information(client, node.node_id, address)
    assert delegate_epoch == 3, "ErrMsg: Last time first delegate epoch {}".format(delegate_epoch)
    assert cumulative_income == current_commission_award, "ErrMsg: Last time first cumulative income {}".format(
        cumulative_income)
    last_delegate_epoch, delegate_total, delegate_total_hes, delegate_reward_total = get_delegate_relevant_amount_and_epoch(
        client, node.node_id)
    assert last_delegate_epoch == 1, "ErrMsg: Last time delegate epoch {}".format(last_delegate_epoch)
    assert delegate_total == delegate_amount - economic.delegate_limit, "The total number of effective commissioned nodes: {}".format(
        delegate_total)
    assert delegate_total_hes == 0, "The total number of inactive nodes commissioned: {}".format(delegate_total_hes)
    assert delegate_reward_total == commission_award, "Total delegated rewards currently issued by the candidate: {}".format(
        delegate_reward_total)
    result = client.ppos.getDelegateReward(address)
    assert result['Ret'][0]['reward'] == current_commission_award


@pytest.mark.P0
@pytest.mark.parametrize('delegate_type', [0, 1])
def test_EI_BC_012_017(client_new_node, delegate_type, reset_environment):
    """
    自由金额多次部分赎回，验证待领取的委托收益（单节点）
    锁仓金额多次部分赎回，验证待领取的委托收益（单节点）
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
    result = client.delegate.delegate(delegate_type, address, amount=delegate_amount)
    assert_code(result, 0)
    log.info("Commissioned successfully, commissioned amount：{}".format(economic.delegate_limit))
    result = client.ppos.getCandidateInfo(node.node_id)
    blocknum = result['Ret']['StakingBlockNum']
    result = client.delegate.withdrew_delegate(blocknum, address)
    assert_code(result, 0)
    delegate_epoch, cumulative_income = get_dividend_information(client, node.node_id, address)
    assert delegate_epoch == 1, "ErrMsg: Last time first delegate epoch {}".format(delegate_epoch)
    assert cumulative_income == 0, "ErrMsg: Last time first cumulative income {}".format(
        cumulative_income)
    last_delegate_epoch, delegate_total, delegate_total_hes, delegate_reward_total = get_delegate_relevant_amount_and_epoch(
        client, node.node_id)
    assert last_delegate_epoch == 1, "ErrMsg: Last time delegate epoch {}".format(last_delegate_epoch)
    assert delegate_total == 0, "The total number of effective commissioned nodes: {}".format(
        delegate_total)
    assert delegate_total_hes == delegate_amount - economic.delegate_limit, "The total number of inactive nodes commissioned: {}".format(
        delegate_total_hes)
    assert delegate_reward_total == 0, "Total delegated rewards currently issued by the candidate: {}".format(
        delegate_reward_total)
    economic.wait_settlement_blocknum(node)
    result = client.delegate.withdrew_delegate(blocknum, address)
    assert_code(result, 0)
    delegate_epoch, cumulative_income = get_dividend_information(client, node.node_id, address)
    assert delegate_epoch == 2, "ErrMsg: Last time first delegate epoch {}".format(delegate_epoch)
    assert cumulative_income == 0, "ErrMsg: Last time first cumulative income {}".format(
        cumulative_income)
    last_delegate_epoch, delegate_total, delegate_total_hes, delegate_reward_total = get_delegate_relevant_amount_and_epoch(
        client, node.node_id)
    assert last_delegate_epoch == 1, "ErrMsg: Last time delegate epoch {}".format(last_delegate_epoch)
    assert delegate_total == delegate_amount - von_amount(economic.delegate_limit,
                                                          2), "The total number of effective commissioned nodes: {}".format(
        delegate_total)
    assert delegate_total_hes == 0, "The total number of inactive nodes commissioned: {}".format(
        delegate_total_hes)
    assert delegate_reward_total == 0, "Total delegated rewards currently issued by the candidate: {}".format(
        delegate_reward_total)
    block_reward, staking_reward = economic.get_current_year_reward(node)
    economic.wait_settlement_blocknum(node)
    result = client.delegate.withdrew_delegate(blocknum, address)
    assert_code(result, 0)
    delegate_amount_total = delegate_amount - von_amount(economic.delegate_limit, 2)
    commission_total_reward = economic.calculate_delegate_reward(node, block_reward, staking_reward)
    current_commission_reward = economic.delegate_cumulative_income(node, block_reward, staking_reward,
                                                                    delegate_amount_total, delegate_amount_total)
    delegate_epoch, cumulative_income = get_dividend_information(client, node.node_id, address)
    assert delegate_epoch == 3, "ErrMsg: Last time first delegate epoch {}".format(delegate_epoch)
    assert cumulative_income == current_commission_reward, "ErrMsg: Last time first cumulative income {}".format(
        cumulative_income)
    last_delegate_epoch, delegate_total, delegate_total_hes, delegate_reward_total = get_delegate_relevant_amount_and_epoch(
        client, node.node_id)
    assert last_delegate_epoch == 1, "ErrMsg: Last time delegate epoch {}".format(last_delegate_epoch)
    assert delegate_total == delegate_amount - von_amount(economic.delegate_limit,
                                                          3), "The total number of effective commissioned nodes: {}".format(
        delegate_total)
    assert delegate_total_hes == 0, "The total number of inactive nodes commissioned: {}".format(
        delegate_total_hes)
    assert delegate_reward_total == commission_total_reward, "Total delegated rewards currently issued by the candidate: {}".format(
        delegate_reward_total)
    result = client.ppos.getDelegateReward(address)
    log.info("Dividend information currently available：{}".format(result))
    assert result['Ret'][0][
               'reward'] == current_commission_reward, "ErrMsg:Delegate rewards currently available {}".format(
        result['Ret'][0]['reward'])


@pytest.mark.P0
def test_EI_BC_013(client_new_node, reset_environment):
    """
    节点被多账户委托，跨结算周期部分赎回，验证待领取的委托收益（生效期N）
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
    commission_total_reward = economic.calculate_delegate_reward(node, block_reward, staking_reward)
    current_commission_award = economic.delegate_cumulative_income(node, block_reward, staking_reward,
                                                                   von_amount(delegate_amount, 2), delegate_amount)
    result = client.ppos.getCandidateInfo(client.node.node_id)
    blocknum = result['Ret']['StakingBlockNum']
    result = client.delegate.withdrew_delegate(blocknum, first_address)
    assert_code(result, 0)
    first_delegate_epoch, first_cumulative_income = get_dividend_information(client, node.node_id, first_address)
    assert first_delegate_epoch == 3, "ErrMsg: Last time first delegate epoch {}".format(first_delegate_epoch)
    assert first_cumulative_income == current_commission_award, "ErrMsg: Last time first cumulative income {}".format(
        first_cumulative_income)
    last_delegate_epoch, delegate_total, delegate_total_hes, delegate_reward_total = get_delegate_relevant_amount_and_epoch(
        client, node.node_id)
    assert last_delegate_epoch == 1, "ErrMsg: Last time delegate epoch {}".format(last_delegate_epoch)
    assert delegate_total == von_amount(delegate_amount,
                                        2) - economic.delegate_limit, "The total number of effective commissioned nodes: {}".format(
        delegate_total)
    assert delegate_total_hes == 0, "The total number of inactive nodes commissioned: {}".format(delegate_total_hes)
    assert delegate_reward_total == commission_total_reward, "Total delegated rewards currently issued by the candidate: {}".format(
        delegate_reward_total)
    result = client.ppos.getDelegateReward(first_address)
    withdrawal_commission = result['Ret'][0]['reward']
    log.info("{} Dividends can be collected in the current settlement period： {}".format(first_address,
                                                                                         withdrawal_commission))
    assert withdrawal_commission == current_commission_award, "ErrMsg: Dividends currently available {}".format(
        withdrawal_commission)


@pytest.mark.P1
@pytest.mark.parametrize('delegate_type', [0, 1])
def test_EI_BC_014_018(client_new_node, delegate_type, reset_environment):
    """
    自由金额赎回全部委托，验证待领取的委托收益（未生效期）
    锁仓金额赎回全部委托，验证待领取的委托收益（未生效期）
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
    result = client.delegate.delegate(delegate_type, address, amount=delegate_amount)
    assert_code(result, 0)
    log.info("Commissioned successfully, commissioned amount：{}".format(economic.delegate_limit))
    result = client.ppos.getCandidateInfo(node.node_id)
    blocknum = result['Ret']['StakingBlockNum']
    result = client.delegate.withdrew_delegate(blocknum, address, amount=delegate_amount)
    assert_code(result, 0)
    result = client.ppos.getDelegateInfo(blocknum, address, node.node_id)
    assert_code(result, 301205)
    last_delegate_epoch, delegate_total, delegate_total_hes, delegate_reward_total = get_delegate_relevant_amount_and_epoch(
        client, node.node_id)
    assert last_delegate_epoch == 1, "ErrMsg: Last time delegate epoch {}".format(last_delegate_epoch)
    assert delegate_total == 0, "The total number of effective commissioned nodes: {}".format(
        delegate_total)
    assert delegate_total_hes == 0, "The total number of inactive nodes commissioned: {}".format(
        delegate_total_hes)
    assert delegate_reward_total == 0, "Total delegated rewards currently issued by the candidate: {}".format(
        delegate_reward_total)


@pytest.mark.P1
@pytest.mark.parametrize('delegate_type', [0, 1])
def test_EI_BC_019_023(clients_new_node, delegate_type, reset_environment):
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
    first_delegate_epoch, first_cumulative_income = get_dividend_information(first_client, first_node.node_id, address)
    assert first_delegate_epoch == 1, "ErrMsg: Last time first delegate epoch {}".format(first_delegate_epoch)
    assert first_cumulative_income == 0, "ErrMsg: Last time first cumulative income {}".format(
        first_cumulative_income)
    second_delegate_epoch, second_cumulative_income = get_dividend_information(first_client, second_node.node_id,
                                                                               address)
    assert second_delegate_epoch == 1, "ErrMsg: Last time first delegate epoch {}".format(second_delegate_epoch)
    assert second_cumulative_income == 0, "ErrMsg: Last time first cumulative income {}".format(
        second_cumulative_income)
    first_last_delegate_epoch, first_delegate_total, first_delegate_total_hes, first_delegate_reward_total = get_delegate_relevant_amount_and_epoch(
        first_client, first_node.node_id)
    assert first_last_delegate_epoch == 1, "ErrMsg: Last time delegate epoch {}".format(first_last_delegate_epoch)
    assert first_delegate_total == 0, "The total number of effective commissioned nodes: {}".format(
        first_delegate_total)
    assert first_delegate_total_hes == delegate_amount - first_economic.delegate_limit, "The total number of inactive nodes commissioned: {}".format(
        first_delegate_total_hes)
    assert first_delegate_reward_total == 0, "Total delegated rewards currently issued by the candidate: {}".format(
        first_delegate_reward_total)
    second_last_delegate_epoch, second_delegate_total, second_delegate_total_hes, second_delegate_reward_total = get_delegate_relevant_amount_and_epoch(
        first_client, second_node.node_id)
    assert second_last_delegate_epoch == 1, "ErrMsg: Last time delegate epoch {}".format(second_last_delegate_epoch)
    assert second_delegate_total == 0, "The total number of effective commissioned nodes: {}".format(
        second_delegate_total)
    assert second_delegate_total_hes == delegate_amount - second_economic.delegate_limit, "The total number of inactive nodes commissioned: {}".format(
        second_delegate_total_hes)
    assert second_delegate_reward_total == 0, "Total delegated rewards currently issued by the candidate: {}".format(
        second_delegate_reward_total)


@pytest.mark.P1
@pytest.mark.parametrize('delegate_type', [0, 1])
def test_EI_BC_020_24(clients_new_node, delegate_type, reset_environment):
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
    first_economic.env.deploy_all()
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
    first_delegate_epoch, first_cumulative_income = get_dividend_information(first_client, first_node.node_id, address)
    assert first_delegate_epoch == 1, "ErrMsg: Last time first delegate epoch {}".format(first_delegate_epoch)
    assert first_cumulative_income == 0, "ErrMsg: Last time first cumulative income {}".format(
        first_cumulative_income)
    second_delegate_epoch, second_cumulative_income = get_dividend_information(first_client, second_node.node_id,
                                                                               address)
    assert second_delegate_epoch == 1, "ErrMsg: Last time first delegate epoch {}".format(second_delegate_epoch)
    assert second_cumulative_income == 0, "ErrMsg: Last time first cumulative income {}".format(
        second_cumulative_income)
    first_last_delegate_epoch, first_delegate_total, first_delegate_total_hes, first_delegate_reward_total = get_delegate_relevant_amount_and_epoch(
        first_client, first_node.node_id)
    assert first_last_delegate_epoch == 1, "ErrMsg: Last time delegate epoch {}".format(first_last_delegate_epoch)
    assert first_delegate_total == 0, "The total number of effective commissioned nodes: {}".format(
        first_delegate_total)
    assert first_delegate_total_hes == delegate_amount - first_economic.delegate_limit, "The total number of inactive nodes commissioned: {}".format(
        first_delegate_total_hes)
    assert first_delegate_reward_total == 0, "Total delegated rewards currently issued by the candidate: {}".format(
        first_delegate_reward_total)
    second_last_delegate_epoch, second_delegate_total, second_delegate_total_hes, second_delegate_reward_total = get_delegate_relevant_amount_and_epoch(
        first_client, second_node.node_id)
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
    first_delegate_epoch, first_cumulative_income = get_dividend_information(first_client, first_node.node_id, address)
    assert first_delegate_epoch == 2, "ErrMsg: Last time first delegate epoch {}".format(first_delegate_epoch)
    assert first_cumulative_income == 0, "ErrMsg: Last time first cumulative income {}".format(
        first_cumulative_income)
    second_delegate_epoch, second_cumulative_income = get_dividend_information(first_client, second_node.node_id,
                                                                               address)
    assert second_delegate_epoch == 2, "ErrMsg: Last time first delegate epoch {}".format(second_delegate_epoch)
    assert second_cumulative_income == 0, "ErrMsg: Last time first cumulative income {}".format(
        second_cumulative_income)
    first_last_delegate_epoch, first_delegate_total, first_delegate_total_hes, first_delegate_reward_total = get_delegate_relevant_amount_and_epoch(
        first_client, first_node.node_id)
    assert first_last_delegate_epoch == 1, "ErrMsg: Last time delegate epoch {}".format(first_last_delegate_epoch)
    assert first_delegate_total == delegate_amount - von_amount(first_economic.delegate_limit,
                                                                2), "The total number of effective commissioned nodes: {}".format(
        first_delegate_total)
    assert first_delegate_total_hes == 0, "The total number of inactive nodes commissioned: {}".format(
        first_delegate_total_hes)
    assert first_delegate_reward_total == 0, "Total delegated rewards currently issued by the candidate: {}".format(
        first_delegate_reward_total)
    second_last_delegate_epoch, second_delegate_total, second_delegate_total_hes, second_delegate_reward_total = get_delegate_relevant_amount_and_epoch(
        first_client, second_node.node_id)
    assert second_last_delegate_epoch == 1, "ErrMsg: Last time delegate epoch {}".format(second_last_delegate_epoch)
    assert second_delegate_total == delegate_amount - von_amount(second_economic.delegate_limit,
                                                                 2), "The total number of effective commissioned nodes: {}".format(
        second_delegate_total)
    assert second_delegate_total_hes == 0, "The total number of inactive nodes commissioned: {}".format(
        second_delegate_total_hes)
    assert second_delegate_reward_total == 0, "Total delegated rewards currently issued by the candidate: {}".format(
        second_delegate_reward_total)
    block_reward, staking_reward = first_economic.get_current_year_reward(first_node)
    first_economic.wait_settlement_blocknum(first_node)
    log.info("Current settlement block height：{}".format(first_node.eth.blockNumber))
    delegate_amount_total = delegate_amount - von_amount(first_economic.delegate_limit, 2)
    first_commission_award = first_economic.calculate_delegate_reward(first_node, block_reward, staking_reward)
    first_current_commission_award = first_economic.delegate_cumulative_income(first_node, block_reward, staking_reward,
                                                                               delegate_amount_total,
                                                                               delegate_amount_total)
    second_commission_award = second_economic.calculate_delegate_reward(second_node, block_reward, staking_reward)
    second_current_commission_award = second_economic.delegate_cumulative_income(second_node, block_reward,
                                                                                 staking_reward, delegate_amount_total,
                                                                                 delegate_amount_total)
    result = first_client.delegate.withdrew_delegate(first_blocknum, address)
    assert_code(result, 0)
    result = first_client.delegate.withdrew_delegate(second_blocknum, address, node_id=second_node.node_id)
    assert_code(result, 0)
    log.info("Current settlement block height：{}".format(first_node.eth.blockNumber))
    first_delegate_epoch, first_cumulative_income = get_dividend_information(first_client, first_node.node_id, address)
    assert first_delegate_epoch == 3, "ErrMsg: Last time first delegate epoch {}".format(first_delegate_epoch)
    assert first_cumulative_income == first_current_commission_award, "ErrMsg: Last time first cumulative income {}".format(
        first_cumulative_income)
    second_delegate_epoch, second_cumulative_income = get_dividend_information(first_client, second_node.node_id,
                                                                               address)
    assert second_delegate_epoch == 3, "ErrMsg: Last time first delegate epoch {}".format(second_delegate_epoch)
    assert second_cumulative_income == second_current_commission_award, "ErrMsg: Last time first cumulative income {}".format(
        second_cumulative_income)
    first_last_delegate_epoch, first_delegate_total, first_delegate_total_hes, first_delegate_reward_total = get_delegate_relevant_amount_and_epoch(
        first_client, first_node.node_id)
    assert first_last_delegate_epoch == 1, "ErrMsg: Last time delegate epoch {}".format(first_last_delegate_epoch)
    assert first_delegate_total == delegate_amount - von_amount(first_economic.delegate_limit,
                                                                3), "The total number of effective commissioned nodes: {}".format(
        first_delegate_total)
    assert first_delegate_total_hes == 0, "The total number of inactive nodes commissioned: {}".format(
        first_delegate_total_hes)
    assert first_delegate_reward_total == first_commission_award, "Total delegated rewards currently issued by the candidate: {}".format(
        first_delegate_reward_total)
    second_last_delegate_epoch, second_delegate_total, second_delegate_total_hes, second_delegate_reward_total = get_delegate_relevant_amount_and_epoch(
        first_client, second_node.node_id)
    assert second_last_delegate_epoch == 1, "ErrMsg: Last time delegate epoch {}".format(second_last_delegate_epoch)
    assert second_delegate_total == delegate_amount - von_amount(second_economic.delegate_limit,
                                                                 3), "The total number of effective commissioned nodes: {}".format(
        second_delegate_total)
    assert second_delegate_total_hes == 0, "The total number of inactive nodes commissioned: {}".format(
        second_delegate_total_hes)
    assert second_delegate_reward_total == second_commission_award, "Total delegated rewards currently issued by the candidate: {}".format(
        second_delegate_reward_total)
    result = first_client.ppos.getDelegateReward(address, node_ids=[first_node.node_id])
    withdrawal_commission = result['Ret'][0]['reward']
    log.info(
        "{} Dividends can be collected in the current settlement period： {}".format(address, withdrawal_commission))
    assert withdrawal_commission == first_current_commission_award, "ErrMsg: Dividends currently available {}".format(
        withdrawal_commission)
    result = first_client.ppos.getDelegateReward(address, node_ids=[second_node.node_id])
    withdrawal_commission = result['Ret'][0]['reward']
    log.info(
        "{} Dividends can be collected in the current settlement period： {}".format(address, withdrawal_commission))
    assert withdrawal_commission == second_current_commission_award, "ErrMsg: Dividends currently available {}".format(
        withdrawal_commission)


@pytest.mark.P1
@pytest.mark.parametrize('delegate_type', [0, 1])
def test_EI_BC_021_25(clients_new_node, delegate_type, reset_environment):
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
    second_node = second_client.node
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
    result = first_client.delegate.withdrew_delegate(first_blocknum, address, amount=delegate_amount)
    assert_code(result, 0)
    result = first_client.delegate.withdrew_delegate(second_blocknum, address, node_id=second_node.node_id,
                                                     amount=delegate_amount)
    assert_code(result, 0)
    result = first_client.ppos.getDelegateInfo(first_blocknum, address, first_node.node_id)
    assert_code(result, 301205)
    result = second_client.ppos.getDelegateInfo(first_blocknum, address, second_node.node_id)
    assert_code(result, 301205)
    first_last_delegate_epoch, first_delegate_total, first_delegate_total_hes, first_delegate_reward_total = get_delegate_relevant_amount_and_epoch(
        first_client, first_node.node_id)
    assert first_last_delegate_epoch == 1, "ErrMsg: Last time delegate epoch {}".format(first_last_delegate_epoch)
    assert first_delegate_total == 0, "The total number of effective commissioned nodes: {}".format(
        first_delegate_total)
    assert first_delegate_total_hes == 0, "The total number of inactive nodes commissioned: {}".format(
        first_delegate_total_hes)
    assert first_delegate_reward_total == 0, "Total delegated rewards currently issued by the candidate: {}".format(
        first_delegate_reward_total)
    second_last_delegate_epoch, second_delegate_total, second_delegate_total_hes, second_delegate_reward_total = get_delegate_relevant_amount_and_epoch(
        first_client, second_node.node_id)
    assert second_last_delegate_epoch == 1, "ErrMsg: Last time delegate epoch {}".format(second_last_delegate_epoch)
    assert second_delegate_total == 0, "The total number of effective commissioned nodes: {}".format(
        second_delegate_total)
    assert second_delegate_total_hes == 0, "The total number of inactive nodes commissioned: {}".format(
        second_delegate_total_hes)
    assert second_delegate_reward_total == 0, "Total delegated rewards currently issued by the candidate: {}".format(
        second_delegate_reward_total)


@pytest.mark.P1
@pytest.mark.parametrize('delegate_type', [0, 1])
def test_EI_BC_022_026(clients_new_node, delegate_type, reset_environment):
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
    first_economic.env.deploy_all()
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
    first_delegate_balance = first_node.eth.getBalance(address)
    log.info("Entrusted account balance：{}".format(first_delegate_balance))
    first_commission_award = first_economic.calculate_delegate_reward(first_node, block_reward, staking_reward)
    first_current_commission_award = first_economic.delegate_cumulative_income(first_node, block_reward, staking_reward,
                                                                               delegate_amount, delegate_amount)
    second_commission_award = second_economic.calculate_delegate_reward(second_node, block_reward, staking_reward)
    second_current_commission_award = first_economic.delegate_cumulative_income(first_node, block_reward,
                                                                                staking_reward, delegate_amount,
                                                                                delegate_amount)
    first_delegate_balance = first_node.eth.getBalance(address)
    log.info("Entrusted account balance：{}".format(first_delegate_balance))
    result = first_client.delegate.withdrew_delegate(first_blocknum, address, amount=delegate_amount)
    assert_code(result, 0)
    result = first_client.delegate.withdrew_delegate(second_blocknum, address, node_id=second_node.node_id,
                                                     amount=delegate_amount)
    assert_code(result, 0)
    second_delegate_balance = first_node.eth.getBalance(address)
    log.info("Entrusted account balance：{}".format(second_delegate_balance))
    result = first_client.ppos.getDelegateInfo(first_blocknum, address, first_node.node_id)
    assert_code(result, 301205)
    result = second_client.ppos.getDelegateInfo(first_blocknum, address, second_node.node_id)
    assert_code(result, 301205)
    first_last_delegate_epoch, first_delegate_total, first_delegate_total_hes, first_delegate_reward_total = get_delegate_relevant_amount_and_epoch(
        first_client, first_node.node_id)
    assert first_last_delegate_epoch == 1, "ErrMsg: Last time delegate epoch {}".format(first_last_delegate_epoch)
    assert first_delegate_total == 0, "The total number of effective commissioned nodes: {}".format(
        first_delegate_total)
    assert first_delegate_total_hes == 0, "The total number of inactive nodes commissioned: {}".format(
        first_delegate_total_hes)
    assert first_delegate_reward_total == first_commission_award, "Total delegated rewards currently issued by the candidate: {}".format(
        first_delegate_reward_total)
    second_last_delegate_epoch, second_delegate_total, second_delegate_total_hes, second_delegate_reward_total = get_delegate_relevant_amount_and_epoch(
        first_client, second_node.node_id)
    assert second_last_delegate_epoch == 1, "ErrMsg: Last time delegate epoch {}".format(second_last_delegate_epoch)
    assert second_delegate_total == 0, "The total number of effective commissioned nodes: {}".format(
        second_delegate_total)
    assert second_delegate_total_hes == 0, "The total number of inactive nodes commissioned: {}".format(
        second_delegate_total_hes)
    assert second_delegate_reward_total == second_commission_award, "Total delegated rewards currently issued by the candidate: {}".format(
        second_delegate_reward_total)
    result = first_client.ppos.getDelegateReward(address, node_ids=[first_node.node_id])
    log.info("Dividend reward information currently available：{}".format(result))
    assert_code(result, 305001)
    result = first_client.ppos.getDelegateReward(address, node_ids=[second_node.node_id])
    log.info("Dividend reward information currently available：{}".format(result))
    assert_code(result, 305001)
    assert first_delegate_balance + first_current_commission_award + second_current_commission_award - second_delegate_balance < first_node.web3.toWei(
        1, 'ether'), "ErrMsg: 账户余额 {}".format(second_delegate_balance)


@pytest.mark.P1
@pytest.mark.parametrize('delegate_type', [0, 1])
def test_EI_BC_027_029(client_new_node, delegate_type, reset_environment):
    """
    跨结算期赎回全部自由委托（生效期N赎回）
    跨结算期赎回全部锁仓委托（生效期N赎回）
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
    result = client.delegate.delegate(delegate_type, address, amount=delegate_amount)
    assert_code(result, 0)
    log.info("Commissioned successfully, commissioned amount：{}".format(economic.delegate_limit))
    economic.wait_settlement_blocknum(node)
    log.info("Current settlement block height：{}".format(node.eth.blockNumber))
    result = client.ppos.getCandidateInfo(node.node_id)
    blocknum = result['Ret']['StakingBlockNum']
    result = client.delegate.withdrew_delegate(blocknum, address, amount=delegate_amount)
    assert_code(result, 0)
    result = client.ppos.getDelegateInfo(blocknum, address, node.node_id)
    assert_code(result, 301205)
    last_delegate_epoch, delegate_total, delegate_total_hes, delegate_reward_total = get_delegate_relevant_amount_and_epoch(
        client, node.node_id)
    assert last_delegate_epoch == 1, "ErrMsg: Last time delegate epoch {}".format(last_delegate_epoch)
    assert delegate_total == 0, "The total number of effective commissioned nodes: {}".format(
        delegate_total)
    assert delegate_total_hes == 0, "The total number of inactive nodes commissioned: {}".format(
        delegate_total_hes)
    assert delegate_reward_total == 0, "Total delegated rewards currently issued by the candidate: {}".format(
        delegate_reward_total)
    result = client.ppos.getDelegateReward(address)
    assert_code(result, 305001)


@pytest.mark.P1
@pytest.mark.parametrize('delegate_type', [0, 1])
def test_EI_BC_028_030(client_new_node, delegate_type, reset_environment):
    """
    跨结算期赎回全部自由委托（生效期N+1赎回）
    跨结算期赎回全部锁仓委托（生效期N+1赎回）
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
    result = client.delegate.delegate(delegate_type, address, amount=delegate_amount)
    assert_code(result, 0)
    log.info("Commissioned successfully, commissioned amount：{}".format(economic.delegate_limit))
    economic.wait_settlement_blocknum(node)
    log.info("Current settlement block height：{}".format(node.eth.blockNumber))
    result = client.ppos.getCandidateInfo(node.node_id)
    blocknum = result['Ret']['StakingBlockNum']
    block_reward, staking_reward = economic.get_current_year_reward(node)
    economic.wait_settlement_blocknum(node)
    log.info("Current settlement block height：{}".format(node.eth.blockNumber))
    first_delegate_balance = node.eth.getBalance(address)
    log.info("Entrusted account balance： {}".format(first_delegate_balance))
    commission_award = economic.calculate_delegate_reward(node, block_reward, staking_reward)
    current_commission_award = economic.delegate_cumulative_income(node, block_reward, staking_reward, delegate_amount,
                                                                   delegate_amount)
    result = client.delegate.withdrew_delegate(blocknum, address, amount=delegate_amount)
    assert_code(result, 0)
    result = client.ppos.getDelegateInfo(blocknum, address, node.node_id)
    assert_code(result, 301205)
    last_delegate_epoch, delegate_total, delegate_total_hes, delegate_reward_total = get_delegate_relevant_amount_and_epoch(
        client, node.node_id)
    assert last_delegate_epoch == 1, "ErrMsg: Last time delegate epoch {}".format(last_delegate_epoch)
    assert delegate_total == 0, "The total number of effective commissioned nodes: {}".format(
        delegate_total)
    assert delegate_total_hes == 0, "The total number of inactive nodes commissioned: {}".format(
        delegate_total_hes)
    assert delegate_reward_total == commission_award, "Total delegated rewards currently issued by the candidate: {}".format(
        delegate_reward_total)
    result = client.ppos.getDelegateReward(address)
    assert_code(result, 305001)
    second_delegate_balance = node.eth.getBalance(address)
    log.info("Entrusted account balance： {}".format(second_delegate_balance))
    assert first_delegate_balance + current_commission_award - second_delegate_balance < node.web3.toWei(1,
                                                                                                         'ether'), "ErrMsg:账户余额 {}".format(
        second_delegate_balance)


@pytest.mark.P1
@pytest.mark.parametrize('amount', [10, 100, 150])
def test_EI_BC_031_032_033(client_new_node, amount, reset_environment):
    """
    生效期N自由金额再委托，赎回部分委托，赎回委托金额<生效期N自由金额（自由首次委托）
    生效期N自由金额再委托，赎回部分委托，赎回委托金额=生效期N自由金额（自由首次委托)
    生效期N自由金额再委托，赎回部分委托，赎回委托金额大于生效期N自由金额（自由首次委托）
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
    log.info("Commissioned successfully, commissioned amount：{}".format(economic.delegate_limit))
    economic.wait_settlement_blocknum(node)
    log.info("Current settlement block height：{}".format(node.eth.blockNumber))
    result = client.delegate.delegate(0, address, amount=delegate_amount)
    assert_code(result, 0)
    result = client.ppos.getCandidateInfo(node.node_id)
    blocknum = result['Ret']['StakingBlockNum']
    redemption_amount = node.web3.toWei(amount, 'ether')
    result = client.delegate.withdrew_delegate(blocknum, address, amount=redemption_amount)
    assert_code(result, 0)
    block_reward, staking_reward = economic.get_current_year_reward(node)
    economic.wait_settlement_blocknum(node)
    log.info("Current block height：{}".format(node.eth.blockNumber))
    current_delegate_amount = von_amount(delegate_amount, 2) - redemption_amount
    if current_delegate_amount < delegate_amount:
        current_delegate_amount = current_delegate_amount
    else:
        current_delegate_amount = delegate_amount
    commission_award_total = economic.calculate_delegate_reward(node, block_reward, staking_reward)
    current_commission_award = economic.delegate_cumulative_income(node, block_reward, staking_reward,
                                                                   current_delegate_amount, current_delegate_amount)
    result = client.ppos.getCandidateInfo(client.node.node_id)
    blocknum = result['Ret']['StakingBlockNum']
    result = client.delegate.withdrew_delegate(blocknum, address)
    assert_code(result, 0)
    delegate_epoch, cumulative_income = get_dividend_information(client, node.node_id, address)
    assert delegate_epoch == 3, "ErrMsg: Last time first delegate epoch {}".format(delegate_epoch)
    assert cumulative_income == current_commission_award, "ErrMsg: Last time first cumulative income {}".format(
        cumulative_income)
    last_delegate_epoch, delegate_total, delegate_total_hes, delegate_reward_total = get_delegate_relevant_amount_and_epoch(
        client, node.node_id)
    assert last_delegate_epoch == 2, "ErrMsg: Last time delegate epoch {}".format(last_delegate_epoch)
    # assert delegate_total == delegate_amount, "The total number of effective commissioned nodes: {}".format(delegate_total)
    assert delegate_total_hes == 0, "The total number of inactive nodes commissioned: {}".format(delegate_total_hes)
    assert delegate_reward_total == commission_award_total, "Total delegated rewards currently issued by the candidate: {}".format(
        delegate_reward_total)
    result = client.ppos.getDelegateReward(address)
    assert result['Ret'][0]['reward'] == current_commission_award, "ErrMsg:Dividends are currently available {}".format(
        result['Ret'][0]['reward'])


@pytest.mark.P1
@pytest.mark.parametrize('first_type,second_type', [(0, 1), (1, 0)])
def test_EI_BC_036_037(client_new_node, first_type, second_type, reset_environment):
    """
    生效期N自由金额再委托，赎回部分委托（自由首次委托）
    生效期N自由金额再委托，赎回部分委托（锁仓首次委托）
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
    result = client.delegate.delegate(first_type, address, amount=delegate_amount)
    assert_code(result, 0)
    log.info("Commissioned successfully, commissioned amount：{}".format(economic.delegate_limit))
    economic.wait_settlement_blocknum(node)
    log.info("Current settlement block height：{}".format(node.eth.blockNumber))
    result = client.delegate.delegate(second_type, address, amount=delegate_amount)
    assert_code(result, 0)
    result = client.ppos.getCandidateInfo(node.node_id)
    blocknum = result['Ret']['StakingBlockNum']
    redemption_amount = node.web3.toWei(150, 'ether')
    result = client.delegate.withdrew_delegate(blocknum, address, amount=redemption_amount)
    assert_code(result, 0)
    delegate_epoch, cumulative_income = get_dividend_information(client, node.node_id, address)
    assert delegate_epoch == 2, "ErrMsg: Last time first delegate epoch {}".format(delegate_epoch)
    assert cumulative_income == 0, "ErrMsg: Last time first cumulative income {}".format(
        cumulative_income)
    last_delegate_epoch, delegate_total, delegate_total_hes, delegate_reward_total = get_delegate_relevant_amount_and_epoch(
        client, node.node_id)
    assert last_delegate_epoch == 2, "ErrMsg: Last time delegate epoch {}".format(last_delegate_epoch)
    assert delegate_total == von_amount(delegate_amount,
                                        2) - redemption_amount, "The total number of effective commissioned nodes: {}".format(
        delegate_total)
    assert delegate_total_hes == 0, "The total number of inactive nodes commissioned: {}".format(delegate_total_hes)
    assert delegate_reward_total == 0, "Total delegated rewards currently issued by the candidate: {}".format(
        delegate_reward_total)


@pytest.mark.P1
def test_EI_BC_038(client_new_node, reset_environment):
    """
    自由、锁仓金额同时委托，赎回部分委托
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
    log.info("Commissioned successfully, commissioned amount：{}".format(economic.delegate_limit))
    result = client.delegate.delegate(1, address, amount=delegate_amount)
    assert_code(result, 0)
    log.info("Commissioned successfully, commissioned amount：{}".format(economic.delegate_limit))
    economic.wait_settlement_blocknum(node)
    log.info("Current settlement block height：{}".format(node.eth.blockNumber))
    result = client.ppos.getCandidateInfo(node.node_id)
    blocknum = result['Ret']['StakingBlockNum']
    redemption_amount = node.web3.toWei(150, 'ether')
    result = client.delegate.withdrew_delegate(blocknum, address, amount=redemption_amount)
    assert_code(result, 0)
    delegate_epoch, cumulative_income = get_dividend_information(client, node.node_id, address)
    assert delegate_epoch == 2, "ErrMsg: Last time first delegate epoch {}".format(delegate_epoch)
    assert cumulative_income == 0, "ErrMsg: Last time first cumulative income {}".format(cumulative_income)
    last_delegate_epoch, delegate_total, delegate_total_hes, delegate_reward_total = get_delegate_relevant_amount_and_epoch(
        client, node.node_id)
    assert last_delegate_epoch == 1, "ErrMsg: Last time delegate epoch {}".format(last_delegate_epoch)
    assert delegate_total == von_amount(delegate_amount,
                                        2) - redemption_amount, "The total number of effective commissioned nodes: {}".format(
        delegate_total)
    assert delegate_total_hes == 0, "The total number of inactive nodes commissioned: {}".format(delegate_total_hes)
    assert delegate_reward_total == 0, "Total delegated rewards currently issued by the candidate: {}".format(
        delegate_reward_total)


@pytest.mark.P1
def test_EI_BC_039(client_new_node, reset_environment):
    """
    自由、锁仓金额同时委托，生效期再委托，赎回部分委托
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
    log.info("Commissioned successfully, commissioned amount：{}".format(economic.delegate_limit))
    result = client.delegate.delegate(1, address, amount=delegate_amount)
    assert_code(result, 0)
    log.info("Commissioned successfully, commissioned amount：{}".format(economic.delegate_limit))
    economic.wait_settlement_blocknum(node)
    log.info("Current settlement block height：{}".format(node.eth.blockNumber))
    result = client.delegate.delegate(0, address, amount=delegate_amount)
    assert_code(result, 0)
    result = client.delegate.delegate(1, address, amount=delegate_amount)
    assert_code(result, 0)
    delegate_epoch, cumulative_income = get_dividend_information(client, node.node_id, address)
    assert delegate_epoch == 2, "ErrMsg: Last time first delegate epoch {}".format(delegate_epoch)
    assert cumulative_income == 0, "ErrMsg: Last time first cumulative income {}".format(
        cumulative_income)
    last_delegate_epoch, delegate_total, delegate_total_hes, delegate_reward_total = get_delegate_relevant_amount_and_epoch(
        client, node.node_id)
    assert last_delegate_epoch == 2, "ErrMsg: Last time delegate epoch {}".format(last_delegate_epoch)
    assert delegate_total == von_amount(delegate_amount,
                                        2), "The total number of effective commissioned nodes: {}".format(
        delegate_total)
    assert delegate_total_hes == von_amount(delegate_amount,
                                            2), "The total number of inactive nodes commissioned: {}".format(
        delegate_total_hes)
    assert delegate_reward_total == 0, "Total delegated rewards currently issued by the candidate: {}".format(
        delegate_reward_total)
    result = client.ppos.getCandidateInfo(node.node_id)
    blocknum = result['Ret']['StakingBlockNum']
    block_reward, staking_reward = economic.get_current_year_reward(node)
    economic.wait_settlement_blocknum(node)
    log.info("Current settlement block height：{}".format(node.eth.blockNumber))
    delegate_amount_total = von_amount(delegate_amount, 2)
    commission_award_total = economic.calculate_delegate_reward(node, block_reward, staking_reward)
    current_commission_award = economic.delegate_cumulative_income(node, block_reward, staking_reward,
                                                                   delegate_amount_total, delegate_amount_total)
    redemption_amount = node.web3.toWei(350, 'ether')
    result = client.delegate.withdrew_delegate(blocknum, address, amount=redemption_amount)
    assert_code(result, 0)
    delegate_epoch, cumulative_income = get_dividend_information(client, node.node_id, address)
    assert delegate_epoch == 3, "ErrMsg: Last time first delegate epoch {}".format(delegate_epoch)
    assert cumulative_income == current_commission_award, "ErrMsg: Last time first cumulative income {}".format(
        cumulative_income)
    last_delegate_epoch, delegate_total, delegate_total_hes, delegate_reward_total = get_delegate_relevant_amount_and_epoch(
        client, node.node_id)
    assert last_delegate_epoch == 2, "ErrMsg: Last time delegate epoch {}".format(last_delegate_epoch)
    assert delegate_total == von_amount(delegate_amount,
                                        4) - redemption_amount, "The total number of effective commissioned nodes: {}".format(
        delegate_total)
    assert delegate_total_hes == 0, "The total number of inactive nodes commissioned: {}".format(delegate_total_hes)
    assert delegate_reward_total == commission_award_total, "Total delegated rewards currently issued by the candidate: {}".format(
        delegate_reward_total)


@pytest.mark.P1
@pytest.mark.parametrize('delegate_type', [0, 1])
def test_EI_BC_040_041(client_new_node, delegate_type, reset_environment):
    """
    生效期N自由金额再委托，赎回全部委托（自由首次委托）
    生效期N锁仓金额再委托，赎回全部委托（锁仓首次委托）
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
    result = client.delegate.delegate(delegate_type, address, amount=delegate_amount)
    assert_code(result, 0)
    log.info("Commissioned successfully, commissioned amount：{}".format(economic.delegate_limit))
    economic.wait_settlement_blocknum(node)
    log.info("Current settlement block height：{}".format(node.eth.blockNumber))
    result = client.delegate.delegate(delegate_type, address, amount=delegate_amount)
    assert_code(result, 0)
    result = client.ppos.getCandidateInfo(node.node_id)
    blocknum = result['Ret']['StakingBlockNum']
    result = client.delegate.withdrew_delegate(blocknum, address, amount=von_amount(delegate_amount, 2))
    assert_code(result, 0)
    result = client.ppos.getDelegateInfo(blocknum, address, node.node_id)
    assert_code(result, 301205)
    last_delegate_epoch, delegate_total, delegate_total_hes, delegate_reward_total = get_delegate_relevant_amount_and_epoch(
        client, node.node_id)
    assert last_delegate_epoch == 2, "ErrMsg: Last time delegate epoch {}".format(last_delegate_epoch)
    assert delegate_total == 0, "The total number of effective commissioned nodes: {}".format(
        delegate_total)
    assert delegate_total_hes == 0, "The total number of inactive nodes commissioned: {}".format(delegate_total_hes)
    assert delegate_reward_total == 0, "Total delegated rewards currently issued by the candidate: {}".format(
        delegate_reward_total)
    economic.wait_settlement_blocknum(node)
    log.info("Current settlement block height：{}".format(node.eth.blockNumber))
    result = client.ppos.getDelegateInfo(blocknum, address, node.node_id)
    assert_code(result, 301205)
    last_delegate_epoch, delegate_total, delegate_total_hes, delegate_reward_total = get_delegate_relevant_amount_and_epoch(
        client, node.node_id)
    assert last_delegate_epoch == 2, "ErrMsg: Last time delegate epoch {}".format(last_delegate_epoch)
    assert delegate_total == 0, "The total number of effective commissioned nodes: {}".format(
        delegate_total)
    assert delegate_total_hes == 0, "The total number of inactive nodes commissioned: {}".format(delegate_total_hes)
    assert delegate_reward_total == 0, "Total delegated rewards currently issued by the candidate: {}".format(
        delegate_reward_total)


@pytest.mark.P1
@pytest.mark.parametrize('first_type,second_type', [(0, 1), (1, 0)])
def test_EI_BC_042_043(client_new_node, first_type, second_type, reset_environment):
    """
    生效期N锁仓锁仓委托，赎回全部委托（自由首次委托）
    生效期N自由锁仓委托，赎回全部委托（锁仓首次委托）
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
    result = client.delegate.delegate(first_type, address, amount=delegate_amount)
    assert_code(result, 0)
    log.info("Commissioned successfully, commissioned amount：{}".format(economic.delegate_limit))
    economic.wait_settlement_blocknum(node)
    log.info("Current settlement block height：{}".format(node.eth.blockNumber))
    result = client.delegate.delegate(second_type, address, amount=delegate_amount)
    assert_code(result, 0)
    result = client.ppos.getCandidateInfo(node.node_id)
    blocknum = result['Ret']['StakingBlockNum']
    result = client.delegate.withdrew_delegate(blocknum, address, amount=von_amount(delegate_amount, 2))
    assert_code(result, 0)
    result = client.ppos.getDelegateInfo(blocknum, address, node.node_id)
    assert_code(result, 301205)
    last_delegate_epoch, delegate_total, delegate_total_hes, delegate_reward_total = get_delegate_relevant_amount_and_epoch(
        client, node.node_id)
    assert last_delegate_epoch == 2, "ErrMsg: Last time delegate epoch {}".format(last_delegate_epoch)
    assert delegate_total == 0, "The total number of effective commissioned nodes: {}".format(
        delegate_total)
    assert delegate_total_hes == 0, "The total number of inactive nodes commissioned: {}".format(delegate_total_hes)
    assert delegate_reward_total == 0, "Total delegated rewards currently issued by the candidate: {}".format(
        delegate_reward_total)
    economic.wait_settlement_blocknum(node)
    log.info("Current settlement block height：{}".format(node.eth.blockNumber))
    result = client.ppos.getDelegateInfo(blocknum, address, node.node_id)
    assert_code(result, 301205)
    last_delegate_epoch, delegate_total, delegate_total_hes, delegate_reward_total = get_delegate_relevant_amount_and_epoch(
        client, node.node_id)
    assert last_delegate_epoch == 2, "ErrMsg: Last time delegate epoch {}".format(last_delegate_epoch)
    assert delegate_total == 0, "The total number of effective commissioned nodes: {}".format(
        delegate_total)
    assert delegate_total_hes == 0, "The total number of inactive nodes commissioned: {}".format(delegate_total_hes)
    assert delegate_reward_total == 0, "Total delegated rewards currently issued by the candidate: {}".format(
        delegate_reward_total)


@pytest.mark.P1
def test_EI_BC_044(clients_new_node, reset_environment):
    """
    连续赎回不同节点委托，验证待领取的委托收益
    :param clients_new_node:
    :return:
    """
    first_client = clients_new_node[0]
    second_client = clients_new_node[1]
    first_economic = first_client.economic
    first_node = first_client.node
    second_node = second_client.node
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
    first_economic.wait_settlement_blocknum(first_node)
    log.info("Current settlement block height：{}".format(first_node.eth.blockNumber))
    result = first_client.ppos.getCandidateInfo(first_node.node_id)
    first_blocknum = result['Ret']['StakingBlockNum']
    result = first_client.delegate.withdrew_delegate(first_blocknum, address)
    assert_code(result, 0)
    result = second_client.ppos.getCandidateInfo(second_node.node_id)
    second_blocknum = result['Ret']['StakingBlockNum']
    result = first_client.delegate.withdrew_delegate(second_blocknum, address, node_id=second_node.node_id)
    assert_code(result, 0)
    first_delegate_epoch, first_cumulative_income = get_dividend_information(first_client, first_node.node_id, address)
    second_delegate_epoch, second_cumulative_income = get_dividend_information(first_client, second_node.node_id,
                                                                               address)
    assert first_delegate_epoch == 2, "ErrMsg: Last time first delegate epoch {}".format(first_delegate_epoch)
    assert first_cumulative_income == 0, "ErrMsg: Last time first cumulative income {}".format(first_cumulative_income)
    assert second_delegate_epoch == 2, "ErrMsg: Last time first delegate epoch {}".format(second_delegate_epoch)
    assert second_cumulative_income == 0, "ErrMsg: Last time first cumulative income {}".format(
        second_cumulative_income)
    first_last_delegate_epoch, first_delegate_total, first_delegate_total_hes, first_delegate_reward_total = get_delegate_relevant_amount_and_epoch(
        first_client, first_node.node_id)
    second_last_delegate_epoch, second_delegate_total, second_delegate_total_hes, second_delegate_reward_total = get_delegate_relevant_amount_and_epoch(
        first_client, second_node.node_id)
    assert first_last_delegate_epoch == 1, "ErrMsg: Last time delegate epoch {}".format(first_last_delegate_epoch)
    assert first_delegate_total == delegate_amount - first_economic.delegate_limit, "The total number of effective commissioned nodes: {}".format(
        first_delegate_total)
    assert first_delegate_total_hes == 0, "The total number of inactive nodes commissioned: {}".format(
        first_delegate_total_hes)
    assert first_delegate_reward_total == 0, "Total delegated rewards currently issued by the candidate: {}".format(
        first_delegate_reward_total)
    result = first_client.ppos.getDelegateReward(address)
    assert result['Ret'][0]['reward'] == 0, "ErrMsg:Receive dividends {}".format(result['Ret'][0]['reward'])
    assert second_last_delegate_epoch == 1, "ErrMsg: Last time delegate epoch {}".format(second_last_delegate_epoch)
    assert second_delegate_total == delegate_amount - first_economic.delegate_limit, "The total number of effective commissioned nodes: {}".format(
        second_delegate_total)
    assert second_delegate_total_hes == 0, "The total number of inactive nodes commissioned: {}".format(
        second_delegate_total_hes)
    assert second_delegate_reward_total == 0, "Total delegated rewards currently issued by the candidate: {}".format(
        second_delegate_reward_total)
    result = second_client.ppos.getDelegateReward(address)
    assert result['Ret'][0]['reward'] == 0, "ErrMsg:Receive dividends {}".format(result['Ret'][0]['reward'])


@pytest.mark.P1
def test_EI_BC_045(client_new_node, reset_environment):
    """
    跨结算周期赎回委托，验证待领取的委托收益
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
    log.info("Commissioned successfully, commissioned amount：{}".format(economic.delegate_limit))
    result = client.ppos.getCandidateInfo(node.node_id)
    blocknum = result['Ret']['StakingBlockNum']
    result = client.delegate.withdrew_delegate(blocknum, address)
    assert_code(result, 0)
    delegate_epoch, cumulative_income = get_dividend_information(client, node.node_id, address)
    assert delegate_epoch == 1, "ErrMsg: Last time first delegate epoch {}".format(delegate_epoch)
    assert cumulative_income == 0, "ErrMsg: Last time first cumulative income {}".format(
        cumulative_income)
    last_delegate_epoch, delegate_total, delegate_total_hes, delegate_reward_total = get_delegate_relevant_amount_and_epoch(
        client, node.node_id)
    assert last_delegate_epoch == 1, "ErrMsg: Last time delegate epoch {}".format(last_delegate_epoch)
    assert delegate_total == 0, "The total number of effective commissioned nodes: {}".format(delegate_total)
    assert delegate_total_hes == delegate_amount - economic.delegate_limit, "The total number of inactive nodes commissioned: {}".format(
        delegate_total_hes)
    assert delegate_reward_total == 0, "Total delegated rewards currently issued by the candidate: {}".format(
        delegate_reward_total)
    economic.wait_settlement_blocknum(node)
    log.info("Current settlement block height：{}".format(node.eth.blockNumber))
    block_reward, staking_reward = economic.get_current_year_reward(node)
    economic.wait_settlement_blocknum(node)
    log.info("Current settlement block height：{}".format(node.eth.blockNumber))
    current_delegate_amount = delegate_amount - economic.delegate_limit
    commission_award = economic.calculate_delegate_reward(node, block_reward, staking_reward)
    current_commission_award = economic.delegate_cumulative_income(node, block_reward, staking_reward,
                                                                   current_delegate_amount, current_delegate_amount)
    result = client.delegate.withdrew_delegate(blocknum, address)
    assert_code(result, 0)
    delegate_epoch, cumulative_income = get_dividend_information(client, node.node_id, address)
    assert delegate_epoch == 3, "ErrMsg: Last time first delegate epoch {}".format(delegate_epoch)
    assert cumulative_income == current_commission_award, "ErrMsg: Last time first cumulative income {}".format(
        cumulative_income)
    last_delegate_epoch, delegate_total, delegate_total_hes, delegate_reward_total = get_delegate_relevant_amount_and_epoch(
        client, node.node_id)
    assert last_delegate_epoch == 1, "ErrMsg: Last time delegate epoch {}".format(last_delegate_epoch)
    assert delegate_total == delegate_amount - von_amount(economic.delegate_limit,
                                                          2), "The total number of effective commissioned nodes: {}".format(
        delegate_total)
    assert delegate_total_hes == 0, "The total number of inactive nodes commissioned: {}".format(delegate_total_hes)
    assert delegate_reward_total == commission_award, "Total delegated rewards currently issued by the candidate: {}".format(
        delegate_reward_total)


@pytest.mark.P1
def test_EI_BC_046(client_new_node, reset_environment):
    """
    节点退出中赎回部分委托，验证待领取的委托收益
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
    log.info("Commissioned successfully, commissioned amount：{}".format(economic.delegate_limit))

    economic.wait_settlement_blocknum(node)
    log.info("Current settlement block height：{}".format(node.eth.blockNumber))
    result = client.staking.withdrew_staking(staking_address)
    assert_code(result, 0)
    result = client.ppos.getCandidateInfo(node.node_id)
    blocknum = result['Ret']['StakingBlockNum']
    result = client.delegate.withdrew_delegate(blocknum, address)
    assert_code(result, 0)
    delegate_epoch, cumulative_income = get_dividend_information(client, node.node_id, address)
    assert delegate_epoch == 2, "ErrMsg: Last time first delegate epoch {}".format(delegate_epoch)
    assert cumulative_income == 0, "ErrMsg: Last time first cumulative income {}".format(
        cumulative_income)
    last_delegate_epoch, delegate_total, delegate_total_hes, delegate_reward_total = get_delegate_relevant_amount_and_epoch(
        client, node.node_id)
    assert last_delegate_epoch == 1, "ErrMsg: Last time delegate epoch {}".format(last_delegate_epoch)
    assert delegate_total == delegate_amount - economic.delegate_limit, "The total number of effective commissioned nodes: {}".format(
        delegate_total)
    assert delegate_total_hes == 0, "The total number of inactive nodes commissioned: {}".format(
        delegate_total_hes)
    assert delegate_reward_total == 0, "Total delegated rewards currently issued by the candidate: {}".format(
        delegate_reward_total)
    block_reward, staking_reward = economic.get_current_year_reward(node)

    economic.wait_settlement_blocknum(node)
    log.info("Current settlement block height：{}".format(node.eth.blockNumber))
    currrent_delegate_amount = delegate_amount - economic.delegate_limit
    current_commission_award = economic.delegate_cumulative_income(node, block_reward, staking_reward,
                                                                   currrent_delegate_amount, currrent_delegate_amount)
    result = client.delegate.withdrew_delegate(blocknum, address)
    assert_code(result, 0)
    delegate_epoch, cumulative_income = get_dividend_information(client, node.node_id, address)
    delegate_reward = economic.calculate_delegate_reward(node, block_reward, staking_reward)
    assert delegate_epoch == 3, "ErrMsg: Last time first delegate epoch {}".format(delegate_epoch)
    assert cumulative_income == current_commission_award, "ErrMsg: Last time first cumulative income {}".format(
        cumulative_income)
    last_delegate_epoch, delegate_total, delegate_total_hes, delegate_reward_total = get_delegate_relevant_amount_and_epoch(
        client, node.node_id)
    assert last_delegate_epoch == 1, "ErrMsg: Last time delegate epoch {}".format(last_delegate_epoch)
    assert delegate_total == delegate_amount - economic.delegate_limit - economic.delegate_limit, "The total number of effective commissioned nodes: {}".format(
        delegate_total)
    assert delegate_total_hes == 0, "The total number of inactive nodes commissioned: {}".format(delegate_total_hes)
    assert delegate_reward_total == delegate_reward, "Total delegated rewards currently issued by the candidate: {}".format(
        delegate_reward_total)
    result = client.ppos.getDelegateReward(address)
    assert result['Ret'][0]['reward'] == current_commission_award


@pytest.mark.P1
def test_EI_BC_047(client_new_node, reset_environment):
    """
    节点退出中赎回全部委托，验证待领取的委托收益
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
    log.info("Commissioned successfully, commissioned amount：{}".format(economic.delegate_limit))
    economic.wait_settlement_blocknum(node)
    log.info("Current settlement block height：{}".format(node.eth.blockNumber))

    result = client.staking.withdrew_staking(staking_address)
    assert_code(result, 0)
    first_balance = node.eth.getBalance(address)
    log.info("Entrusted account balance：{}".format(first_balance))
    last_delegate_epoch, delegate_total, delegate_total_hes, delegate_reward_total = get_delegate_relevant_amount_and_epoch(
        client, node.node_id)
    assert last_delegate_epoch == 1, "ErrMsg: Last time delegate epoch {}".format(last_delegate_epoch)
    assert delegate_total == delegate_amount, "The total number of effective commissioned nodes: {}".format(
        delegate_total)
    assert delegate_total_hes == 0, "The total number of inactive nodes commissioned: {}".format(
        delegate_total_hes)
    assert delegate_reward_total == 0, "Total delegated rewards currently issued by the candidate: {}".format(
        delegate_reward_total)
    block_reward, staking_reward = economic.get_current_year_reward(node)
    economic.wait_settlement_blocknum(node)
    log.info("Current settlement block height：{}".format(node.eth.blockNumber))
    result = client.ppos.getCandidateInfo(node.node_id)
    blocknum = result['Ret']['StakingBlockNum']
    result = client.delegate.withdrew_delegate(blocknum, address, amount=delegate_amount)
    assert_code(result, 0)
    commission_award = economic.calculate_delegate_reward(node, block_reward, staking_reward)
    current_commission_award = economic.delegate_cumulative_income(node, block_reward, staking_reward, delegate_amount,
                                                                   delegate_amount)
    result = client.ppos.getDelegateInfo(blocknum, address, node.node_id)
    assert_code(result, 301205)
    last_delegate_epoch, delegate_total, delegate_total_hes, delegate_reward_total = get_delegate_relevant_amount_and_epoch(
        client, node.node_id)
    assert last_delegate_epoch == 1, "ErrMsg: Last time delegate epoch {}".format(last_delegate_epoch)
    assert delegate_total == 0, "The total number of effective commissioned nodes: {}".format(
        delegate_total)
    assert delegate_total_hes == 0, "The total number of inactive nodes commissioned: {}".format(delegate_total_hes)
    assert delegate_reward_total == commission_award, "Total delegated rewards currently issued by the candidate: {}".format(
        delegate_reward_total)
    result = client.ppos.getDelegateReward(address)
    assert_code(result, 305001)
    second_balance = node.eth.getBalance(address)
    log.info("Entrusted account balance：{}".format(second_balance))
    assert (first_balance + current_commission_award) - second_balance < node.web3.toWei(1,
                                                                                         'ether'), "ErrMsg: Account Balance {}".format(
        second_balance)


@pytest.mark.P1
def test_EI_BC_048(client_new_node, reset_environment):
    """
    节点已退出质押赎回部分委托，验证待领取的委托收益
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
    log.info("Commissioned successfully, commissioned amount：{}".format(economic.delegate_limit))

    economic.wait_settlement_blocknum(node)
    log.info("Current settlement block height：{}".format(node.eth.blockNumber))
    result = client.staking.withdrew_staking(staking_address)
    assert_code(result, 0)
    result = client.ppos.getCandidateInfo(node.node_id)
    blocknum = result['Ret']['StakingBlockNum']
    block_reward, staking_reward = economic.get_current_year_reward(node)
    economic.wait_settlement_blocknum(node)
    log.info("Current settlement block height：{}".format(node.eth.blockNumber))
    current_commission_award = economic.delegate_cumulative_income(node, block_reward, staking_reward, delegate_amount,
                                                                   delegate_amount)

    economic.wait_settlement_blocknum(node, 2)
    log.info("Current settlement block height：{}".format(node.eth.blockNumber))
    result = client.delegate.withdrew_delegate(blocknum, address)
    assert_code(result, 0)
    first_commission_balance = node.eth.getBalance(address)
    log.info("Account balance before receiving dividends： {}".format(first_commission_balance))

    result = client.ppos.getDelegateInfo(blocknum, address, node.node_id)
    log.info("Commission information：{}".format(result))
    info = result['Ret']
    delegate_epoch = info['DelegateEpoch']
    cumulative_income = info['CumulativeIncome']

    assert delegate_epoch == 6, "ErrMsg: Last time first delegate epoch {}".format(delegate_epoch)
    assert cumulative_income == current_commission_award, "ErrMsg: Last time first cumulative income {}".format(
        cumulative_income)
    result = client.ppos.getCandidateInfo(node.node_id)
    assert_code(result, 301204)
    second_commission_balance = node.eth.getBalance(address)
    log.info("Account balance before receiving dividends： {}".format(second_commission_balance))
    assert second_commission_balance == first_commission_balance, "ErrMsg: Account Balance {}".format(
        second_commission_balance)


@pytest.mark.P1
def test_EI_BC_049(client_new_node, reset_environment):
    """
    节点已退出质押赎回全部委托，验证待领取的委托收益
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
    log.info("Commissioned successfully, commissioned amount：{}".format(economic.delegate_limit))

    economic.wait_settlement_blocknum(node)
    log.info("Current settlement block height：{}".format(node.eth.blockNumber))
    result = client.staking.withdrew_staking(staking_address)
    assert_code(result, 0)
    result = client.ppos.getCandidateInfo(node.node_id)
    blocknum = result['Ret']['StakingBlockNum']
    block_reward, staking_reward = economic.get_current_year_reward(node)

    economic.wait_settlement_blocknum(node)
    log.info("Current settlement block height：{}".format(node.eth.blockNumber))
    current_commission_award = economic.delegate_cumulative_income(node, block_reward, staking_reward, delegate_amount,
                                                                   delegate_amount)

    first_commission_balance = node.eth.getBalance(address)
    log.info("Account balance before receiving dividends： {}".format(first_commission_balance))

    economic.wait_settlement_blocknum(node, 2)
    log.info("Current settlement block height：{}".format(node.eth.blockNumber))
    result = client.delegate.withdrew_delegate(blocknum, address, amount=delegate_amount)
    assert_code(result, 0)
    result = client.ppos.getDelegateInfo(blocknum, address, node.node_id)
    assert_code(result, 301205)
    result = client.ppos.getCandidateInfo(node.node_id)
    assert_code(result, 301204)
    result = client.ppos.getDelegateReward(address)
    assert_code(result, 305001)
    second_commission_balance = node.eth.getBalance(address)
    log.info("Account balance before receiving dividends： {}".format(second_commission_balance))
    assert (first_commission_balance + current_commission_award) - second_commission_balance < node.web3.toWei(1,
                                                                                                               'ether'), "ErrMsg: Account Balance {}".format(
        second_commission_balance)


@pytest.mark.P1
def test_EI_BC_050(client_new_node, reset_environment):
    """
    节点被惩罚赎回委托，验证待领取的委托收益
    :param client_new_node:
    :return:
    """
    client = client_new_node
    economic = client.economic
    node = client.node
    address = create_restricting_plan(client)
    report_address, _ = economic.account.generate_account(node.web3, node.web3.toWei(1000, 'ether'))
    log.info("Create delegate account：{}".format(address))
    create_staking_node(client)
    log.info("Create pledge node id :{}".format(node.node_id))
    delegate_amount = von_amount(economic.delegate_limit, 10)
    result = client.delegate.delegate(0, address, amount=delegate_amount)
    assert_code(result, 0)
    result = client.ppos.getCandidateInfo(node.node_id)
    blocknum = result['Ret']['StakingBlockNum']

    economic.wait_settlement_blocknum(node)
    log.info("Current settlement block height：{}".format(node.eth.blockNumber))
    for i in range(4):
        result = check_node_in_list(node.node_id, client.ppos.getValidatorList)
        log.info("Current node in consensus list status：{}".format(result))
        if result:
            # view Current block
            current_block = client_new_node.node.eth.blockNumber
            log.info("Current block: {}".format(current_block))
            # Report prepareblock signature
            report_information = mock_duplicate_sign(1, node.nodekey, node.blsprikey, current_block)
            log.info("Report information: {}".format(report_information))
            result = client_new_node.duplicatesign.reportDuplicateSign(1, report_information, report_address)
            assert_code(result, 0)
            result = client.delegate.withdrew_delegate(blocknum, address)
            assert_code(result, 0)
            delegate_epoch, cumulative_income = get_dividend_information(client, node.node_id, address)
            assert delegate_epoch == 2, "ErrMsg: Last time first delegate epoch {}".format(delegate_epoch)
            assert cumulative_income == 0, "ErrMsg: Last time first cumulative income {}".format(cumulative_income)
            last_delegate_epoch, delegate_total, delegate_total_hes, delegate_reward_total = get_delegate_relevant_amount_and_epoch(
                client, node.node_id)
            assert last_delegate_epoch == 1, "ErrMsg: Last time delegate epoch {}".format(last_delegate_epoch)
            assert delegate_total == delegate_amount - economic.delegate_limit, "The total number of effective commissioned nodes: {}".format(
                delegate_total)
            assert delegate_total_hes == 0, "The total number of inactive nodes commissioned: {}".format(
                delegate_total_hes)
            assert delegate_reward_total == 0, "Total delegated rewards currently issued by the candidate: {}".format(
                delegate_reward_total)
            result = client.ppos.getDelegateReward(address)
            assert result['Ret'][0]['reward'] == 0, "ErrMsg:Receive dividends {}".format(result['Ret'][0]['reward'])

            economic.wait_settlement_blocknum(node)
            log.info("Current settlement block height：{}".format(node.eth.blockNumber))
            result = client.delegate.withdrew_delegate(blocknum, address)
            assert_code(result, 0)
            delegate_epoch, cumulative_income = get_dividend_information(client, node.node_id, address)
            assert delegate_epoch == 3, "ErrMsg: Last time first delegate epoch {}".format(delegate_epoch)
            assert cumulative_income == 0, "ErrMsg: Last time first cumulative income {}".format(cumulative_income)
            last_delegate_epoch, delegate_total, delegate_total_hes, delegate_reward_total = get_delegate_relevant_amount_and_epoch(
                client, node.node_id)
            assert last_delegate_epoch == 1, "ErrMsg: Last time delegate epoch {}".format(last_delegate_epoch)
            assert delegate_total == delegate_amount - von_amount(economic.delegate_limit,
                                                                  2), "The total number of effective commissioned nodes: {}".format(
                delegate_total)
            assert delegate_total_hes == 0, "The total number of inactive nodes commissioned: {}".format(
                delegate_total_hes)
            assert delegate_reward_total == 0, "Total delegated rewards currently issued by the candidate: {}".format(
                delegate_reward_total)
            result = client.ppos.getDelegateReward(address)
            assert result['Ret'][0]['reward'] == 0, "ErrMsg:Receive dividends {}".format(result['Ret'][0]['reward'])
            break
        else:
            # wait consensus block
            economic.wait_consensus_blocknum(node)


@pytest.mark.P1
@pytest.mark.parametrize('delegate_type', [0, 1])
def test_EI_BC_055_058(client_new_node, delegate_type, reset_environment):
    """
    自由金额委托首次领取分红（生效期N）
    锁仓金额委托首次领取分红（生效期N）
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
    result = client.delegate.delegate(delegate_type, address, amount=delegate_amount)
    assert_code(result, 0)
    log.info("Commissioned successfully, commissioned amount：{}".format(economic.delegate_limit))
    economic.wait_settlement_blocknum(node)
    log.info("Current settlement block height：{}".format(node.eth.blockNumber))
    result = client.delegate.withdraw_delegate_reward(address)
    assert_code(result, 0)
    delegate_epoch, cumulative_income = get_dividend_information(client, node.node_id, address)
    assert delegate_epoch == 1, "ErrMsg: Last time first delegate epoch {}".format(delegate_epoch)
    assert cumulative_income == 0, "ErrMsg: Last time first cumulative income {}".format(cumulative_income)
    last_delegate_epoch, delegate_total, delegate_total_hes, delegate_reward_total = get_delegate_relevant_amount_and_epoch(
        client, node.node_id)
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
def test_EI_BC_056_059(client_new_node, delegate_type, reset_environment):
    """
    跨结算期自由金额委托首次领取分红
    跨结算期锁仓金额委托首次领取分红
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
    result = client.delegate.delegate(delegate_type, address, amount=delegate_amount)
    assert_code(result, 0)
    log.info("Commissioned successfully, commissioned amount：{}".format(economic.delegate_limit))
    result = client.delegate.withdraw_delegate_reward(address)
    assert_code(result, 0)
    delegate_epoch, cumulative_income = get_dividend_information(client, node.node_id, address)
    assert delegate_epoch == 1, "ErrMsg: Last time first delegate epoch {}".format(delegate_epoch)
    assert cumulative_income == 0, "ErrMsg: Last time first cumulative income {}".format(cumulative_income)
    last_delegate_epoch, delegate_total, delegate_total_hes, delegate_reward_total = get_delegate_relevant_amount_and_epoch(
        client, node.node_id)
    assert last_delegate_epoch == 1, "ErrMsg: Last time delegate epoch {}".format(last_delegate_epoch)
    assert delegate_total == 0, "The total number of effective commissioned nodes: {}".format(
        delegate_total)
    assert delegate_total_hes == delegate_amount, "The total number of inactive nodes commissioned: {}".format(
        delegate_total_hes)
    assert delegate_reward_total == 0, "Total delegated rewards currently issued by the candidate: {}".format(
        delegate_reward_total)
    result = client.ppos.getDelegateReward(address)
    assert result['Ret'][0]['reward'] == 0, "ErrMsg:Receive dividends {}".format(result['Ret'][0]['reward'])
    economic.wait_settlement_blocknum(node)
    log.info("Current settlement block height：{}".format(node.eth.blockNumber))
    result = client.delegate.withdraw_delegate_reward(address)
    assert_code(result, 0)
    delegate_epoch, cumulative_income = get_dividend_information(client, node.node_id, address)
    assert delegate_epoch == 1, "ErrMsg: Last time first delegate epoch {}".format(delegate_epoch)
    assert cumulative_income == 0, "ErrMsg: Last time first cumulative income {}".format(cumulative_income)
    last_delegate_epoch, delegate_total, delegate_total_hes, delegate_reward_total = get_delegate_relevant_amount_and_epoch(
        client, node.node_id)
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
    delegate_epoch, cumulative_income = get_dividend_information(client, node.node_id, address)
    assert delegate_epoch == 3, "ErrMsg: Last time first delegate epoch {}".format(delegate_epoch)
    assert cumulative_income == 0, "ErrMsg: Last time first cumulative income {}".format(cumulative_income)
    last_delegate_epoch, delegate_total, delegate_total_hes, delegate_reward_total = get_delegate_relevant_amount_and_epoch(
        client, node.node_id)
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
def test_EI_BC_060(client_new_node, reset_environment):
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
    delegate_epoch, cumulative_income = get_dividend_information(client, node.node_id, address)
    assert delegate_epoch == 1, "ErrMsg: Last time first delegate epoch {}".format(delegate_epoch)
    assert cumulative_income == 0, "ErrMsg: Last time first cumulative income {}".format(cumulative_income)
    last_delegate_epoch, delegate_total, delegate_total_hes, delegate_reward_total = get_delegate_relevant_amount_and_epoch(
        client, node.node_id)
    assert last_delegate_epoch == 1, "ErrMsg: Last time delegate epoch {}".format(last_delegate_epoch)
    assert delegate_total == 0, "The total number of effective commissioned nodes: {}".format(
        delegate_total)
    assert delegate_total_hes == von_amount(delegate_amount,
                                            2), "The total number of inactive nodes commissioned: {}".format(
        delegate_total_hes)
    assert delegate_reward_total == 0, "Total delegated rewards currently issued by the candidate: {}".format(
        delegate_reward_total)
    result = client.ppos.getDelegateReward(address)
    assert result['Ret'][0]['reward'] == 0, "ErrMsg:Receive dividends {}".format(result['Ret'][0]['reward'])


@pytest.mark.P1
def test_EI_BC_061(client_new_node, reset_environment):
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
    delegate_epoch, cumulative_income = get_dividend_information(client, node.node_id, address)
    assert delegate_epoch == 1, "ErrMsg: Last time first delegate epoch {}".format(delegate_epoch)
    assert cumulative_income == 0, "ErrMsg: Last time first cumulative income {}".format(cumulative_income)
    last_delegate_epoch, delegate_total, delegate_total_hes, delegate_reward_total = get_delegate_relevant_amount_and_epoch(
        client, node.node_id)
    assert last_delegate_epoch == 1, "ErrMsg: Last time delegate epoch {}".format(last_delegate_epoch)
    assert delegate_total == von_amount(delegate_amount,
                                        2), "The total number of effective commissioned nodes: {}".format(
        delegate_total)
    assert delegate_total_hes == 0, "The total number of inactive nodes commissioned: {}".format(
        delegate_total_hes)
    assert delegate_reward_total == 0, "Total delegated rewards currently issued by the candidate: {}".format(
        delegate_reward_total)
    result = client.ppos.getDelegateReward(address)
    assert result['Ret'][0]['reward'] == 0, "ErrMsg:Receive dividends {}".format(result['Ret'][0]['reward'])


@pytest.mark.P1
def test_EI_BC_062(client_new_node, reset_environment):
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
    delegate_epoch, cumulative_income = get_dividend_information(client, node.node_id, address)
    assert delegate_epoch == 1, "ErrMsg: Last time first delegate epoch {}".format(delegate_epoch)
    assert cumulative_income == 0, "ErrMsg: Last time first cumulative income {}".format(cumulative_income)
    last_delegate_epoch, delegate_total, delegate_total_hes, delegate_reward_total = get_delegate_relevant_amount_and_epoch(
        client, node.node_id)
    assert last_delegate_epoch == 1, "ErrMsg: Last time delegate epoch {}".format(last_delegate_epoch)
    assert delegate_total == 0, "The total number of effective commissioned nodes: {}".format(
        delegate_total)
    assert delegate_total_hes == von_amount(delegate_amount,
                                            2), "The total number of inactive nodes commissioned: {}".format(
        delegate_total_hes)
    assert delegate_reward_total == 0, "Total delegated rewards currently issued by the candidate: {}".format(
        delegate_reward_total)
    result = client.ppos.getDelegateReward(address)
    assert result['Ret'][0]['reward'] == 0, "ErrMsg:Receive dividends {}".format(result['Ret'][0]['reward'])
    economic.wait_settlement_blocknum(node)
    log.info("Current settlement block height：{}".format(node.eth.blockNumber))
    result = client.delegate.withdraw_delegate_reward(address)
    assert_code(result, 0)
    delegate_epoch, cumulative_income = get_dividend_information(client, node.node_id, address)
    assert delegate_epoch == 1, "ErrMsg: Last time first delegate epoch {}".format(delegate_epoch)
    assert cumulative_income == 0, "ErrMsg: Last time first cumulative income {}".format(cumulative_income)
    last_delegate_epoch, delegate_total, delegate_total_hes, delegate_reward_total = get_delegate_relevant_amount_and_epoch(
        client, node.node_id)
    assert last_delegate_epoch == 1, "ErrMsg: Last time delegate epoch {}".format(last_delegate_epoch)
    assert delegate_total == von_amount(delegate_amount,
                                        2), "The total number of effective commissioned nodes: {}".format(
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
    delegate_epoch, cumulative_income = get_dividend_information(client, node.node_id, address)
    assert delegate_epoch == 3, "ErrMsg: Last time first delegate epoch {}".format(delegate_epoch)
    assert cumulative_income == 0, "ErrMsg: Last time first cumulative income {}".format(cumulative_income)
    last_delegate_epoch, delegate_total, delegate_total_hes, delegate_reward_total = get_delegate_relevant_amount_and_epoch(
        client, node.node_id)
    assert last_delegate_epoch == 1, "ErrMsg: Last time delegate epoch {}".format(last_delegate_epoch)
    assert delegate_total == von_amount(delegate_amount,
                                        2), "The total number of effective commissioned nodes: {}".format(
        delegate_total)
    assert delegate_total_hes == 0, "The total number of inactive nodes commissioned: {}".format(
        delegate_total_hes)
    assert delegate_reward_total == commission_award, "Total delegated rewards currently issued by the candidate: {}".format(
        delegate_reward_total)
    result = client.ppos.getDelegateReward(address)
    assert result['Ret'][0]['reward'] == 0, "ErrMsg:Receive dividends {}".format(result['Ret'][0]['reward'])


@pytest.mark.P1
@pytest.mark.parametrize('delegate_type', [0, 1])
def test_EI_BC_063_064(clients_new_node, delegate_type, reset_environment):
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
    second_node = second_client.node
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
    first_delegate_epoch, first_cumulative_income = get_dividend_information(first_client, first_node.node_id, address)
    second_delegate_epoch, second_cumulative_income = get_dividend_information(first_client, second_node.node_id,
                                                                               address)
    assert first_delegate_epoch == 1, "ErrMsg: Last time first delegate epoch {}".format(first_delegate_epoch)
    assert first_cumulative_income == 0, "ErrMsg: Last time first cumulative income {}".format(first_cumulative_income)
    assert second_delegate_epoch == 1, "ErrMsg: Last time first delegate epoch {}".format(second_delegate_epoch)
    assert second_cumulative_income == 0, "ErrMsg: Last time first cumulative income {}".format(
        second_cumulative_income)
    first_last_delegate_epoch, first_delegate_total, first_delegate_total_hes, first_delegate_reward_total = get_delegate_relevant_amount_and_epoch(
        first_client, first_node.node_id)
    second_last_delegate_epoch, second_delegate_total, second_delegate_total_hes, second_delegate_reward_total = get_delegate_relevant_amount_and_epoch(
        first_client, second_node.node_id)
    assert first_last_delegate_epoch == 1, "ErrMsg: Last time delegate epoch {}".format(first_last_delegate_epoch)
    assert first_delegate_total == 0, "The total number of effective commissioned nodes: {}".format(
        first_delegate_total)
    assert first_delegate_total_hes == delegate_amount, "The total number of inactive nodes commissioned: {}".format(
        first_delegate_total_hes)
    assert first_delegate_reward_total == 0, "Total delegated rewards currently issued by the candidate: {}".format(
        first_delegate_reward_total)
    result = first_client.ppos.getDelegateReward(address)
    assert result['Ret'][0]['reward'] == 0, "ErrMsg:Receive dividends {}".format(result['Ret'][0]['reward'])
    assert second_last_delegate_epoch == 1, "ErrMsg: Last time delegate epoch {}".format(second_last_delegate_epoch)
    assert second_delegate_total == 0, "The total number of effective commissioned nodes: {}".format(
        second_delegate_total)
    assert second_delegate_total_hes == delegate_amount, "The total number of inactive nodes commissioned: {}".format(
        second_delegate_total_hes)
    assert second_delegate_reward_total == 0, "Total delegated rewards currently issued by the candidate: {}".format(
        second_delegate_reward_total)
    result = second_client.ppos.getDelegateReward(address)
    assert result['Ret'][0]['reward'] == 0, "ErrMsg:Receive dividends {}".format(result['Ret'][0]['reward'])


@pytest.mark.P1
@pytest.mark.parametrize('delegate_type', [0, 1])
def test_EI_BC_065_066(clients_new_node, delegate_type, reset_environment):
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
    second_node = second_client.node
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
    first_delegate_epoch, first_cumulative_income = get_dividend_information(first_client, first_node.node_id, address)
    second_delegate_epoch, second_cumulative_income = get_dividend_information(first_client, second_node.node_id,
                                                                               address)
    assert first_delegate_epoch == 1, "ErrMsg: Last time first delegate epoch {}".format(first_delegate_epoch)
    assert first_cumulative_income == 0, "ErrMsg: Last time first cumulative income {}".format(first_cumulative_income)
    assert second_delegate_epoch == 1, "ErrMsg: Last time first delegate epoch {}".format(second_delegate_epoch)
    assert second_cumulative_income == 0, "ErrMsg: Last time first cumulative income {}".format(
        second_cumulative_income)
    first_last_delegate_epoch, first_delegate_total, first_delegate_total_hes, first_delegate_reward_total = get_delegate_relevant_amount_and_epoch(
        first_client, first_node.node_id)
    second_last_delegate_epoch, second_delegate_total, second_delegate_total_hes, second_delegate_reward_total = get_delegate_relevant_amount_and_epoch(
        first_client, second_node.node_id)
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
def test_EI_BC_067_068(clients_new_node, delegate_type, reset_environment):
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
    first_delegate_epoch, first_cumulative_income = get_dividend_information(first_client, first_node.node_id, address)
    second_delegate_epoch, second_cumulative_income = get_dividend_information(first_client, second_node.node_id,
                                                                               address)
    assert first_delegate_epoch == 3, "ErrMsg: Last time first delegate epoch {}".format(first_delegate_epoch)
    assert first_cumulative_income == 0, "ErrMsg: Last time first cumulative income {}".format(first_cumulative_income)
    assert second_delegate_epoch == 3, "ErrMsg: Last time first delegate epoch {}".format(second_delegate_epoch)
    assert second_cumulative_income == 0, "ErrMsg: Last time first cumulative income {}".format(
        second_cumulative_income)
    first_last_delegate_epoch, first_delegate_total, first_delegate_total_hes, first_delegate_reward_total = get_delegate_relevant_amount_and_epoch(
        first_client, first_node.node_id)
    second_last_delegate_epoch, second_delegate_total, second_delegate_total_hes, second_delegate_reward_total = get_delegate_relevant_amount_and_epoch(
        first_client, second_node.node_id)
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
def test_EI_BC_069(client_new_node, reset_environment):
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
    delegate_epoch, cumulative_income = get_dividend_information(client, node.node_id, address)
    assert delegate_epoch == 1, "ErrMsg: Last time first delegate epoch {}".format(delegate_epoch)
    assert cumulative_income == 0, "ErrMsg: Last time first cumulative income {}".format(cumulative_income)
    last_delegate_epoch, delegate_total, delegate_total_hes, delegate_reward_total = get_delegate_relevant_amount_and_epoch(
        client, node.node_id)
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
    delegate_epoch, cumulative_income = get_dividend_information(client, node.node_id, address)
    assert delegate_epoch == 3, "ErrMsg: Last time first delegate epoch {}".format(delegate_epoch)
    assert cumulative_income == 0, "ErrMsg: Last time first cumulative income {}".format(cumulative_income)
    last_delegate_epoch, delegate_total, delegate_total_hes, delegate_reward_total = get_delegate_relevant_amount_and_epoch(
        client, node.node_id)
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
def test_EI_BC_070(client_new_node, reset_environment):
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
    candadite_info = node.ppos.getCandidateInfo(node.node_id)
    staking_num = candadite_info['Ret']['StakingBlockNum']

    economic.wait_settlement_blocknum(node)
    log.info("Current settlement block height：{}".format(node.eth.blockNumber))
    result = client.staking.withdrew_staking(staking_address)
    assert_code(result, 0)
    result = client.delegate.withdraw_delegate_reward(address)
    assert_code(result, 0)
    delegate_epoch, cumulative_income = get_dividend_information(client, node.node_id, address)
    assert delegate_epoch == 1, "ErrMsg: Last time first delegate epoch {}".format(delegate_epoch)
    assert cumulative_income == 0, "ErrMsg: Last time first cumulative income {}".format(cumulative_income)
    last_delegate_epoch, delegate_total, delegate_total_hes, delegate_reward_total = get_delegate_relevant_amount_and_epoch(
        client, node.node_id)
    assert last_delegate_epoch == 1, "ErrMsg: Last time delegate epoch {}".format(last_delegate_epoch)
    assert delegate_total == delegate_amount, "The total number of effective commissioned nodes: {}".format(
        delegate_total)
    assert delegate_total_hes == 0, "The total number of inactive nodes commissioned: {}".format(
        delegate_total_hes)
    assert delegate_reward_total == 0, "Total delegated rewards currently issued by the candidate: {}".format(
        delegate_reward_total)
    result = client.ppos.getDelegateReward(address)
    assert result['Ret'][0]['reward'] == 0, "ErrMsg:Receive dividends {}".format(result['Ret'][0]['reward'])
    # block_reward, staking_reward = economic.get_current_year_reward(node)

    # economic.wait_settlement_blocknum(node)
    # log.info("Current settlement block height：{}".format(node.eth.blockNumber))

    economic.wait_settlement_blocknum(node, 2)
    log.info("Current settlement block height：{}".format(node.eth.blockNumber))
    delegate_balance = node.eth.getBalance(address)
    log.info("Entrusted account balance： {}".format(delegate_balance))
    result = client.delegate.withdraw_delegate_reward(address)
    assert_code(result, 0)

    result = client.ppos.getDelegateInfo(staking_num, address, node.node_id)
    delegate_epoch = result['Ret']['DelegateEpoch']
    cumulative_income = result['Ret']['CumulativeIncome']

    assert delegate_epoch == 5, "ErrMsg: Last time first delegate epoch {}".format(delegate_epoch)
    assert cumulative_income == 0, "ErrMsg: Last time first cumulative income {}".format(cumulative_income)
    result = client.ppos.getCandidateInfo(node.node_id)
    assert_code(result, 301204)
    result = client.ppos.getDelegateReward(address)
    assert result['Ret'][0]['reward'] == 0, "ErrMsg:Receive dividends {}".format(result['Ret'][0]['reward'])


@pytest.mark.P1
def test_EI_BC_071(client_new_node, reset_environment):
    """
    节点被惩罚领取分红，验证待领取的委托收益
    :param client_new_node:
    :return:
    """
    client = client_new_node
    economic = client.economic
    node = client.node
    address = create_restricting_plan(client)
    report_address, _ = economic.account.generate_account(node.web3, node.web3.toWei(1000, 'ether'))
    log.info("Create delegate account：{}".format(address))
    create_staking_node(client)
    log.info("Create pledge node id :{}".format(node.node_id))
    delegate_amount = von_amount(economic.delegate_limit, 10)
    result = client.delegate.delegate(0, address, amount=delegate_amount)
    assert_code(result, 0)
    economic.wait_settlement_blocknum(node)
    log.info("Current settlement block height：{}".format(node.eth.blockNumber))
    for i in range(4):
        result = check_node_in_list(node.node_id, client.ppos.getValidatorList)
        log.info("Current node in consensus list status：{}".format(result))
        if result:
            # view Current block
            current_block = client_new_node.node.eth.blockNumber
            log.info("Current block: {}".format(current_block))
            # Report prepareblock signature
            report_information = mock_duplicate_sign(1, node.nodekey, node.blsprikey, current_block)
            log.info("Report information: {}".format(report_information))
            result = client_new_node.duplicatesign.reportDuplicateSign(1, report_information, report_address)
            assert_code(result, 0)
            result = client.delegate.withdraw_delegate_reward(address)
            assert_code(result, 0)
            delegate_epoch, cumulative_income = get_dividend_information(client, node.node_id, address)
            assert delegate_epoch == 1, "ErrMsg: Last time first delegate epoch {}".format(delegate_epoch)
            assert cumulative_income == 0, "ErrMsg: Last time first cumulative income {}".format(cumulative_income)
            last_delegate_epoch, delegate_total, delegate_total_hes, delegate_reward_total = get_delegate_relevant_amount_and_epoch(
                client, node.node_id)
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
            delegate_epoch, cumulative_income = get_dividend_information(client, node.node_id, address)
            assert delegate_epoch == 1, "ErrMsg: Last time first delegate epoch {}".format(delegate_epoch)
            assert cumulative_income == 0, "ErrMsg: Last time first cumulative income {}".format(cumulative_income)
            last_delegate_epoch, delegate_total, delegate_total_hes, delegate_reward_total = get_delegate_relevant_amount_and_epoch(
                client, node.node_id)
            assert last_delegate_epoch == 1, "ErrMsg: Last time delegate epoch {}".format(last_delegate_epoch)
            assert delegate_total == delegate_amount, "The total number of effective commissioned nodes: {}".format(
                delegate_total)
            assert delegate_total_hes == 0, "The total number of inactive nodes commissioned: {}".format(
                delegate_total_hes)
            assert delegate_reward_total == 0, "Total delegated rewards currently issued by the candidate: {}".format(
                delegate_reward_total)
            result = client.ppos.getDelegateReward(address)
            assert result['Ret'][0]['reward'] == 0, "ErrMsg:Receive dividends {}".format(result['Ret'][0]['reward'])
            break
        else:
            # wait consensus block
            economic.wait_consensus_blocknum(node)


@pytest.mark.P1
def test_EI_BC_072(client_new_node, reset_environment):
    """
    同个生效期N二次委托、赎回部分委托、领取分红
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
    log.info("Commissioned successfully, commissioned amount：{}".format(economic.delegate_limit))

    economic.wait_settlement_blocknum(node)
    log.info("Current settlement block height：{}".format(node.eth.blockNumber))
    result = client.delegate.delegate(0, address, amount=delegate_amount)
    assert_code(result, 0)
    delegate_epoch, cumulative_income = get_dividend_information(client, node.node_id, address)
    assert delegate_epoch == 2, "ErrMsg: Last time delegate epoch {}".format(delegate_epoch)
    assert cumulative_income == 0, "ErrMsg: Last time cumulative income {}".format(cumulative_income)
    last_delegate_epoch, delegate_total, delegate_total_hes, delegate_reward_total = get_delegate_relevant_amount_and_epoch(
        client, node.node_id)
    assert last_delegate_epoch == 2, "ErrMsg: Last time delegate epoch {}".format(last_delegate_epoch)
    assert delegate_total == delegate_amount, "The total number of effective commissioned nodes: {}".format(
        delegate_total)
    assert delegate_total_hes == delegate_amount, "The total number of inactive nodes commissioned: {}".format(
        delegate_total_hes)
    assert delegate_reward_total == 0, "Total delegated rewards currently issued by the candidate: {}".format(
        delegate_reward_total)
    result = client.ppos.getCandidateInfo(node.node_id)
    blocknum = result['Ret']['StakingBlockNum']
    result = client.delegate.withdrew_delegate(blocknum, address)
    assert_code(result, 0)
    delegate_epoch, cumulative_income = get_dividend_information(client, node.node_id, address)
    assert delegate_epoch == 2, "ErrMsg: Last time delegate epoch {}".format(delegate_epoch)
    assert cumulative_income == 0, "ErrMsg: Last time cumulative income {}".format(cumulative_income)
    last_delegate_epoch, delegate_total, delegate_total_hes, delegate_reward_total = get_delegate_relevant_amount_and_epoch(
        client, node.node_id)
    assert last_delegate_epoch == 2, "ErrMsg: Last time delegate epoch {}".format(last_delegate_epoch)
    assert delegate_total == delegate_amount, "The total number of effective commissioned nodes: {}".format(
        delegate_total)
    assert delegate_total_hes == delegate_amount - economic.delegate_limit, "The total number of inactive nodes commissioned: {}".format(
        delegate_total_hes)
    assert delegate_reward_total == 0, "Total delegated rewards currently issued by the candidate: {}".format(
        delegate_reward_total)
    result = client.delegate.withdraw_delegate_reward(address)
    assert_code(result, 0)
    delegate_epoch, cumulative_income = get_dividend_information(client, node.node_id, address)
    assert delegate_epoch == 2, "ErrMsg: Last time delegate epoch {}".format(delegate_epoch)
    assert cumulative_income == 0, "ErrMsg: Last time cumulative income {}".format(cumulative_income)
    last_delegate_epoch, delegate_total, delegate_total_hes, delegate_reward_total = get_delegate_relevant_amount_and_epoch(
        client, node.node_id)
    assert last_delegate_epoch == 2, "ErrMsg: Last time delegate epoch {}".format(last_delegate_epoch)
    assert delegate_total == delegate_amount, "The total number of effective commissioned nodes: {}".format(
        delegate_total)
    assert delegate_total_hes == delegate_amount - economic.delegate_limit, "The total number of inactive nodes commissioned: {}".format(
        delegate_total_hes)
    assert delegate_reward_total == 0, "Total delegated rewards currently issued by the candidate: {}".format(
        delegate_reward_total)


@pytest.mark.P0
def test_EI_BC_073(client_new_node, reset_environment):
    """
    不同生效期N二次委托、赎回部分委托、领取分红
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
    log.info("Commissioned successfully, commissioned amount：{}".format(economic.delegate_limit))
    economic.wait_settlement_blocknum(node)
    log.info("Current settlement block height：{}".format(node.eth.blockNumber))
    result = client.delegate.delegate(0, address, amount=delegate_amount)
    assert_code(result, 0)
    delegate_epoch, cumulative_income = get_dividend_information(client, node.node_id, address)
    assert delegate_epoch == 2, "ErrMsg: Last time delegate epoch {}".format(delegate_epoch)
    assert cumulative_income == 0, "ErrMsg: Last time cumulative income {}".format(cumulative_income)
    last_delegate_epoch, delegate_total, delegate_total_hes, delegate_reward_total = get_delegate_relevant_amount_and_epoch(
        client, node.node_id)
    assert last_delegate_epoch == 2, "ErrMsg: Last time delegate epoch {}".format(last_delegate_epoch)
    assert delegate_total == delegate_amount, "The total number of effective commissioned nodes: {}".format(
        delegate_total)
    assert delegate_total_hes == delegate_amount, "The total number of inactive nodes commissioned: {}".format(
        delegate_total_hes)
    assert delegate_reward_total == 0, "Total delegated rewards currently issued by the candidate: {}".format(
        delegate_reward_total)
    block_reward, staking_reward = economic.get_current_year_reward(node)
    economic.wait_settlement_blocknum(node)
    log.info("Current settlement block height：{}".format(node.eth.blockNumber))
    commission_award_total = economic.calculate_delegate_reward(node, block_reward, staking_reward)
    third_current_commission_award = economic.delegate_cumulative_income(node, block_reward, staking_reward,
                                                                         delegate_amount, delegate_amount)
    result = client.ppos.getCandidateInfo(node.node_id)
    blocknum = result['Ret']['StakingBlockNum']
    result = client.delegate.withdrew_delegate(blocknum, address)
    assert_code(result, 0)
    delegate_epoch, cumulative_income = get_dividend_information(client, node.node_id, address)
    assert delegate_epoch == 3, "ErrMsg: Last time delegate epoch {}".format(delegate_epoch)
    assert cumulative_income == third_current_commission_award, "ErrMsg: Last time cumulative income {}".format(
        cumulative_income)
    last_delegate_epoch, delegate_total, delegate_total_hes, delegate_reward_total = get_delegate_relevant_amount_and_epoch(
        client, node.node_id)
    assert last_delegate_epoch == 2, "ErrMsg: Last time delegate epoch {}".format(last_delegate_epoch)
    assert delegate_total == von_amount(delegate_amount,
                                        2) - economic.delegate_limit, "The total number of effective commissioned nodes: {}".format(
        delegate_total)
    assert delegate_total_hes == 0, "The total number of inactive nodes commissioned: {}".format(
        delegate_total_hes)
    assert delegate_reward_total == commission_award_total, "Total delegated rewards currently issued by the candidate: {}".format(
        delegate_reward_total)
    block_reward, staking_reward = economic.get_current_year_reward(node)
    first_balance = node.eth.getBalance(address)
    log.info("Entrusted account balance： {}".format(first_balance))
    economic.wait_settlement_blocknum(node)
    log.info("Current settlement block height：{}".format(node.eth.blockNumber))
    commission_award_total = economic.calculate_delegate_reward(node, block_reward, staking_reward)
    current_delegate_amount = von_amount(delegate_amount, 2) - economic.delegate_limit
    current_commission_award = economic.delegate_cumulative_income(node, block_reward, staking_reward,
                                                                   current_delegate_amount, current_delegate_amount)
    result = client.ppos.getDelegateReward(address)
    log.info("Receive dividends：{}".format(result))
    result = client.delegate.withdraw_delegate_reward(address)
    assert_code(result, 0)
    second_balance = node.eth.getBalance(address)
    log.info("Entrusted account balance： {}".format(second_balance))
    gas = get_getDelegateReward_gas_fee(client, 1, 1)
    delegate_epoch, cumulative_income = get_dividend_information(client, node.node_id, address)
    assert delegate_epoch == 4, "ErrMsg: Last time delegate epoch {}".format(delegate_epoch)
    assert cumulative_income == 0, "ErrMsg: Last time cumulative income {}".format(cumulative_income)
    last_delegate_epoch, delegate_total, delegate_total_hes, delegate_reward_total = get_delegate_relevant_amount_and_epoch(
        client, node.node_id)
    assert last_delegate_epoch == 2, "ErrMsg: Last time delegate epoch {}".format(last_delegate_epoch)
    assert delegate_total == von_amount(delegate_amount,
                                        2) - economic.delegate_limit, "The total number of effective commissioned nodes: {}".format(
        delegate_total)
    assert delegate_total_hes == 0, "The total number of inactive nodes commissioned: {}".format(
        delegate_total_hes)
    assert delegate_reward_total == von_amount(commission_award_total,
                                               2), "Total delegated rewards currently issued by the candidate: {}".format(
        delegate_reward_total)
    result = client.ppos.getDelegateReward(address)
    assert result['Ret'][0]['reward'] == 0, "ErrMsg: Dividends currently available {}".format(
        result['Ret'][0]['reward'])
    assert first_balance + third_current_commission_award + current_commission_award - gas == second_balance


@pytest.mark.P1
def test_EI_BC_075(clients_new_node, reset_environment):
    """
    委托多节点，生效期+1赎回单个节点全部委托
    :param clients_new_node:
    :return:
    """
    first_client = clients_new_node[0]
    second_client = clients_new_node[1]
    first_economic = first_client.economic
    first_node = first_client.node
    second_economic = second_client.economic
    second_node = second_client.node
    first_economic.env.deploy_all()
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
    first_economic.wait_settlement_blocknum(first_node)
    log.info("Current settlement block height：{}".format(first_node.eth.blockNumber))
    result = first_client.ppos.getCandidateInfo(first_node.node_id)
    blocknum = result['Ret']['StakingBlockNum']
    block_reward, staking_reward = first_economic.get_current_year_reward(first_node)
    first_balance = first_node.eth.getBalance(address)
    log.info("Entrusted account balance: {}".format(first_balance))
    first_economic.wait_settlement_blocknum(first_node)
    log.info("Current settlement block height：{}".format(first_node.eth.blockNumber))
    first_award_total = first_economic.calculate_delegate_reward(first_node, block_reward, staking_reward)
    first_current_commission_award = first_economic.delegate_cumulative_income(first_node, block_reward, staking_reward,
                                                                               delegate_amount, delegate_amount)
    second_award_total = second_economic.calculate_delegate_reward(second_node, block_reward, staking_reward)
    second_current_commission_award = first_economic.delegate_cumulative_income(second_node, block_reward,
                                                                                staking_reward, delegate_amount,
                                                                                delegate_amount)
    result = first_client.delegate.withdrew_delegate(blocknum, address, amount=delegate_amount,
                                                     node_id=first_node.node_id)
    assert_code(result, 0)
    second_balance = first_node.eth.getBalance(address)
    log.info("Entrusted account balance: {}".format(second_balance))
    result = first_client.ppos.getDelegateInfo(blocknum, address, node_id=first_node.node_id)
    assert_code(result, 301205)
    second_delegate_epoch, second_cumulative_income = get_dividend_information(first_client, second_node.node_id,
                                                                               address)
    assert second_delegate_epoch == 1, "ErrMsg: Last time first delegate epoch {}".format(second_delegate_epoch)
    assert second_cumulative_income == 0, "ErrMsg: Last time first cumulative income {}".format(
        second_cumulative_income)
    first_last_delegate_epoch, first_delegate_total, first_delegate_total_hes, first_delegate_reward_total = get_delegate_relevant_amount_and_epoch(
        first_client, first_node.node_id)
    second_last_delegate_epoch, second_delegate_total, second_delegate_total_hes, second_delegate_reward_total = get_delegate_relevant_amount_and_epoch(
        first_client, second_node.node_id)
    assert first_last_delegate_epoch == 1, "ErrMsg: Last time delegate epoch {}".format(first_last_delegate_epoch)
    assert first_delegate_total == 0, "The total number of effective commissioned nodes: {}".format(
        first_delegate_total)
    assert first_delegate_total_hes == 0, "The total number of inactive nodes commissioned: {}".format(
        first_delegate_total_hes)
    assert first_delegate_reward_total == first_award_total, "Total delegated rewards currently issued by the candidate: {}".format(
        first_delegate_reward_total)
    result = first_client.ppos.getDelegateReward(address, node_ids=[first_node.node_id])
    assert_code(result, 305001)
    assert first_balance + first_current_commission_award - second_balance < first_node.web3.toWei(1, 'ether')
    assert second_last_delegate_epoch == 1, "ErrMsg: Last time delegate epoch {}".format(second_last_delegate_epoch)
    assert second_delegate_total == delegate_amount, "The total number of effective commissioned nodes: {}".format(
        second_delegate_total)
    assert second_delegate_total_hes == 0, "The total number of inactive nodes commissioned: {}".format(
        second_delegate_total_hes)
    assert second_delegate_reward_total == second_award_total, "Total delegated rewards currently issued by the candidate: {}".format(
        second_delegate_reward_total)
    result = second_client.ppos.getDelegateReward(address, node_ids=[second_node.node_id])
    assert result['Ret'][0]['reward'] == second_current_commission_award, "ErrMsg:Receive dividends {}".format(
        result['Ret'][0]['reward'])


@pytest.mark.P1
def test_EI_BC_076(client_new_node, reset_environment):
    """
    调整分红比例，查询账户在各节点未提取委托奖励
    :param client_new_node:
    :param reset_environment:
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
    log.info("Commissioned successfully, commissioned amount：{}".format(economic.delegate_limit))
    economic.wait_settlement_blocknum(node, 1)
    log.info("Current settlement block height：{}".format(node.eth.blockNumber))
    log.info("staking_address {} balance: {}".format(staking_address, node.eth.getBalance(staking_address)))
    result = client.staking.edit_candidate(staking_address, staking_address, reward_per=2000)
    assert_code(result, 0)
    block_reward, staking_reward = economic.get_current_year_reward(node)
    log.info("block_reward: {}  staking_reward: {}".format(block_reward, staking_reward))
    economic.wait_settlement_blocknum(node)
    log.info("Current settlement block height：{}".format(node.eth.blockNumber))
    commission_award_total = economic.calculate_delegate_reward(node, block_reward, staking_reward, reward=1000)
    current_commission_award = economic.delegate_cumulative_income(node, block_reward, staking_reward, delegate_amount,
                                                                   delegate_amount, reward=1000)
    result = client.delegate.delegate(0, address, amount=delegate_amount)
    assert_code(result, 0)
    delegate_epoch, cumulative_income = get_dividend_information(client, node.node_id, address)
    assert delegate_epoch == 4, "ErrMsg: Last time delegate epoch {}".format(delegate_epoch)
    assert cumulative_income == von_amount(current_commission_award, 2), "ErrMsg: Last time cumulative income {}".format(
        cumulative_income)
    last_delegate_epoch, delegate_total, delegate_total_hes, delegate_reward_total = get_delegate_relevant_amount_and_epoch(
        client, node.node_id)
    assert last_delegate_epoch == 4, "ErrMsg: Last time delegate epoch {}".format(last_delegate_epoch)
    assert delegate_total == delegate_amount, "The total number of effective commissioned nodes: {}".format(
        delegate_total)
    assert delegate_total_hes == delegate_amount, "The total number of inactive nodes commissioned: {}".format(
        delegate_total_hes)
    assert delegate_reward_total == von_amount(commission_award_total, 2), "Total delegated rewards currently issued by the candidate: {}".format(
        delegate_reward_total)
    result = client.ppos.getDelegateReward(address)
    assert result['Ret'][0]['reward'] == von_amount(commission_award_total, 2), "ErrMsg: Dividends currently available {}".format(
        result['Ret'][0]['reward'])
    block_reward, staking_reward = economic.get_current_year_reward(node)
    economic.wait_settlement_blocknum(node)
    log.info("Current settlement block height：{}".format(node.eth.blockNumber))
    second_commission_award_total = economic.calculate_delegate_reward(node, block_reward, staking_reward)
    second_current_commission_award = economic.delegate_cumulative_income(node, block_reward, staking_reward,
                                                                          delegate_amount,
                                                                          delegate_amount)
    result = client.delegate.delegate(0, address, amount=delegate_amount)
    assert_code(result, 0)
    delegate_epoch, cumulative_income = get_dividend_information(client, node.node_id, address)
    assert delegate_epoch == 5, "ErrMsg: Last time delegate epoch {}".format(delegate_epoch)
    assert cumulative_income == von_amount(current_commission_award, 2) + second_current_commission_award, "ErrMsg: Last time cumulative income {}".format(
        cumulative_income)
    last_delegate_epoch, delegate_total, delegate_total_hes, delegate_reward_total = get_delegate_relevant_amount_and_epoch(
        client, node.node_id)
    assert last_delegate_epoch == 5, "ErrMsg: Last time delegate epoch {}".format(last_delegate_epoch)
    assert delegate_total == von_amount(delegate_amount,
                                        2), "The total number of effective commissioned nodes: {}".format(
        delegate_total)
    assert delegate_total_hes == delegate_amount, "The total number of inactive nodes commissioned: {}".format(
        delegate_total_hes)
    assert delegate_reward_total == von_amount(commission_award_total, 2) + second_commission_award_total, "Total delegated rewards currently issued by the candidate: {}".format(
        delegate_reward_total)
    result = client.ppos.getDelegateReward(address)
    assert result['Ret'][0][
               'reward'] == von_amount(current_commission_award, 2) + second_current_commission_award, "ErrMsg: Dividends currently available {}".format(
        result['Ret'][0]['reward'])


@pytest.mark.P2
def test_EI_BC_077(client_new_node, reset_environment):
    """
    调整分红比例，查询账户在各节点未提取委托奖励(多个委托)
    :param client_new_node:
    :param reset_environment:
    :return:
    """
    client = client_new_node
    economic = client.economic
    node = client.node
    first_address = create_restricting_plan(client)
    second_address = create_restricting_plan(client)
    log.info("Create delegate account：{}".format(first_address))
    log.info("Create delegate account：{}".format(second_address))
    staking_address = create_staking_node(client)
    log.info("Create pledge node id :{}".format(node.node_id))
    delegate_amount = von_amount(economic.delegate_limit, 10)
    result = client.delegate.delegate(0, first_address, amount=delegate_amount)
    assert_code(result, 0)
    result = client.delegate.delegate(0, second_address, amount=delegate_amount)
    assert_code(result, 0)
    log.info("Commissioned successfully, commissioned amount：{}".format(economic.delegate_limit))
    economic.wait_settlement_blocknum(node, 1)
    log.info("Current settlement block height：{}".format(node.eth.blockNumber))
    result = client.staking.edit_candidate(staking_address, staking_address, reward_per=2000)
    assert_code(result, 0)
    block_reward, staking_reward = economic.get_current_year_reward(node)
    economic.wait_settlement_blocknum(node)
    log.info("Current settlement block height：{}".format(node.eth.blockNumber))
    commission_award_total = economic.calculate_delegate_reward(node, block_reward, staking_reward, reward=1000)
    current_commission_award = economic.delegate_cumulative_income(node, block_reward, staking_reward, von_amount(delegate_amount, 2),delegate_amount, reward=1000)
    result = client.delegate.delegate(0, first_address, amount=delegate_amount)
    assert_code(result, 0)
    log.info("Current settlement block height：{}".format(node.eth.blockNumber))
    result = client.delegate.delegate(0, second_address, amount=delegate_amount)
    assert_code(result, 0)
    log.info("Current settlement block height：{}".format(node.eth.blockNumber))
    delegate_epoch, cumulative_income = get_dividend_information(client, node.node_id, first_address)
    assert delegate_epoch == 4, "ErrMsg: Last time delegate epoch {}".format(delegate_epoch)
    assert cumulative_income == von_amount(current_commission_award, 2), "ErrMsg: Last time cumulative income {}".format(
        cumulative_income)
    last_delegate_epoch, delegate_total, delegate_total_hes, delegate_reward_total = get_delegate_relevant_amount_and_epoch(
        client, node.node_id)
    assert last_delegate_epoch == 4, "ErrMsg: Last time delegate epoch {}".format(last_delegate_epoch)
    assert delegate_total == von_amount(delegate_amount,
                                        2), "The total number of effective commissioned nodes: {}".format(
        delegate_total)
    assert delegate_total_hes == von_amount(delegate_amount,
                                            2), "The total number of inactive nodes commissioned: {}".format(
        delegate_total_hes)
    assert delegate_reward_total == von_amount(commission_award_total, 2), "Total delegated rewards currently issued by the candidate: {}".format(
        delegate_reward_total)
    result = client.ppos.getDelegateReward(first_address)
    assert result['Ret'][0]['reward'] == von_amount(current_commission_award, 2), "ErrMsg: Dividends currently available {}".format(
        result['Ret'][0]['reward'])
    block_reward, staking_reward = economic.get_current_year_reward(node)
    economic.wait_settlement_blocknum(node)
    log.info("Current settlement block height：{}".format(node.eth.blockNumber))
    second_commission_award_total = economic.calculate_delegate_reward(node, block_reward, staking_reward)
    second_current_commission_award = economic.delegate_dividend_income(second_commission_award_total,
                                                                        von_amount(delegate_amount, 2),
                                                                        delegate_amount)
    result = client.delegate.delegate(0, first_address, amount=delegate_amount)
    assert_code(result, 0)
    result = client.delegate.delegate(0, second_address, amount=delegate_amount)
    assert_code(result, 0)
    delegate_epoch, cumulative_income = get_dividend_information(client, node.node_id, first_address)
    assert delegate_epoch == 5, "ErrMsg: Last time delegate epoch {}".format(delegate_epoch)
    assert cumulative_income == von_amount(current_commission_award, 2) + second_current_commission_award, "ErrMsg: Last time cumulative income {}".format(
        cumulative_income)
    last_delegate_epoch, delegate_total, delegate_total_hes, delegate_reward_total = get_delegate_relevant_amount_and_epoch(
        client, node.node_id)
    assert last_delegate_epoch == 5, "ErrMsg: Last time delegate epoch {}".format(last_delegate_epoch)
    assert delegate_total == von_amount(delegate_amount,
                                        4), "The total number of effective commissioned nodes: {}".format(
        delegate_total)
    assert delegate_total_hes == von_amount(delegate_amount,
                                            2), "The total number of inactive nodes commissioned: {}".format(
        delegate_total_hes)
    assert delegate_reward_total == von_amount(commission_award_total, 2) + second_commission_award_total, "Total delegated rewards currently issued by the candidate: {}".format(
        delegate_reward_total)
    result = client.ppos.getDelegateReward(first_address)
    assert result['Ret'][0][
               'reward'] == von_amount(current_commission_award, 2) + second_current_commission_award, "ErrMsg: Dividends currently available {}".format(
        result['Ret'][0]['reward'])


@pytest.mark.P1
def test_EI_BC_078(client_new_node, reset_environment):
    """
    多账户委托，其中一个节点赎回
    :param client_new_node:
    :param reset_environment:
    :return:
    """
    client = client_new_node
    economic = client.economic
    node = client.node
    first_address = create_restricting_plan(client)
    second_address = create_restricting_plan(client)
    log.info("Create delegate account：{}".format(first_address))
    log.info("Create delegate account：{}".format(second_address))
    create_staking_node(client)
    log.info("Create pledge node id :{}".format(node.node_id))
    delegate_amount = von_amount(economic.delegate_limit, 10)
    result = client.delegate.delegate(0, first_address, amount=delegate_amount)
    assert_code(result, 0)
    result = client.delegate.delegate(0, second_address, amount=delegate_amount)
    assert_code(result, 0)
    log.info("Commissioned successfully, commissioned amount：{}".format(economic.delegate_limit))
    economic.wait_settlement_blocknum(node)
    log.info("Current settlement block height：{}".format(node.eth.blockNumber))
    result = client.ppos.getCandidateInfo(node.node_id)
    blocknum = result['Ret']['StakingBlockNum']
    result = client.delegate.withdrew_delegate(blocknum, first_address)
    assert_code(result, 0)
    block_reward, staking_reward = economic.get_current_year_reward(node)
    economic.wait_settlement_blocknum(node)
    log.info("Current settlement block height：{}".format(node.eth.blockNumber))
    delegate_amount_total = von_amount(delegate_amount, 2) - economic.delegate_limit
    commission_award_total = economic.calculate_delegate_reward(node, block_reward, staking_reward)
    first_current_commission_award = economic.delegate_cumulative_income(node, block_reward, staking_reward,
                                                                         delegate_amount_total,
                                                                         delegate_amount - economic.delegate_limit)
    second_current_commission_award = economic.delegate_cumulative_income(node, block_reward, staking_reward,
                                                                          delegate_amount_total,
                                                                          delegate_amount)

    result = client.delegate.delegate(0, first_address, amount=delegate_amount)
    assert_code(result, 0)

    result = client.delegate.withdrew_delegate(blocknum, second_address)
    assert_code(result, 0)

    delegate_epoch, cumulative_income = get_dividend_information(client, node.node_id, first_address)
    assert delegate_epoch == 3, "ErrMsg: Last time delegate epoch {}".format(delegate_epoch)
    assert cumulative_income == first_current_commission_award, "ErrMsg: Last time cumulative income {}".format(
        cumulative_income)
    delegate_epoch, cumulative_income = get_dividend_information(client, node.node_id, second_address)
    assert delegate_epoch == 3, "ErrMsg: Last time delegate epoch {}".format(delegate_epoch)
    assert cumulative_income == second_current_commission_award, "ErrMsg: Last time cumulative income {}".format(
        cumulative_income)
    last_delegate_epoch, delegate_total, delegate_total_hes, delegate_reward_total = get_delegate_relevant_amount_and_epoch(
        client, node.node_id)
    assert last_delegate_epoch == 3, "ErrMsg: Last time delegate epoch {}".format(last_delegate_epoch)
    assert delegate_total == delegate_amount_total - economic.delegate_limit, "The total number of effective commissioned nodes: {}".format(
        delegate_total)
    assert delegate_total_hes == delegate_amount, "The total number of inactive nodes commissioned: {}".format(
        delegate_total_hes)
    assert delegate_reward_total == commission_award_total, "Total delegated rewards currently issued by the candidate: {}".format(
        delegate_reward_total)
    result = client.ppos.getDelegateReward(first_address)
    assert result['Ret'][0][
               'reward'] == first_current_commission_award, "ErrMsg: Dividends currently available {}".format(
        result['Ret'][0]['reward'])
    result = client.ppos.getDelegateReward(second_address)
    assert result['Ret'][0][
               'reward'] == second_current_commission_award, "ErrMsg: Dividends currently available {}".format(
        result['Ret'][0]['reward'])


@pytest.mark.P1
def test_EI_BC_079(client_new_node, reset_environment):
    """
    多账户委托，不同结算期操作
    :param client_new_node:
    :param reset_environment:
    :return:
    """
    client = client_new_node
    economic = client.economic
    node = client.node
    first_address = create_restricting_plan(client)
    second_address = create_restricting_plan(client)
    log.info("Create delegate account：{}".format(first_address))
    log.info("Create delegate account：{}".format(second_address))
    create_staking_node(client)
    log.info("Create pledge node id :{}".format(node.url))
    delegate_amount = von_amount(economic.delegate_limit, 10)
    result = client.delegate.delegate(0, first_address, amount=delegate_amount)
    assert_code(result, 0)
    log.info("Commissioned successfully, commissioned amount：{}".format(economic.delegate_limit))
    economic.wait_settlement_blocknum(node)
    log.info("Current settlement block height：{}".format(node.eth.blockNumber))
    result = client.delegate.delegate(0, second_address, amount=delegate_amount)
    assert_code(result, 0)
    result = client.ppos.getCandidateInfo(node.node_id)
    blocknum = result['Ret']['StakingBlockNum']
    result = client.delegate.withdrew_delegate(blocknum, first_address)
    assert_code(result, 0)
    block_reward, staking_reward = economic.get_current_year_reward(node)
    economic.wait_settlement_blocknum(node)
    log.info("Current settlement block height：{}".format(node.eth.blockNumber))
    delegate_amount_total = delegate_amount - economic.delegate_limit
    first_commission_award_total = economic.calculate_delegate_reward(node, block_reward, staking_reward)
    first_current_commission_award = economic.delegate_cumulative_income(node, block_reward, staking_reward,
                                                                         delegate_amount_total,
                                                                         delegate_amount_total)
    result = client.delegate.delegate(0, first_address, amount=delegate_amount)
    assert_code(result, 0)

    result = client.delegate.withdrew_delegate(blocknum, second_address)
    assert_code(result, 0)
    delegate_epoch, cumulative_income = get_dividend_information(client, node.node_id, first_address)
    assert delegate_epoch == 3, "ErrMsg: Last time delegate epoch {}".format(delegate_epoch)
    assert cumulative_income == first_current_commission_award, "ErrMsg: Last time cumulative income {}".format(
        cumulative_income)
    delegate_epoch, cumulative_income = get_dividend_information(client, node.node_id, second_address)
    assert delegate_epoch == 3, "ErrMsg: Last time delegate epoch {}".format(delegate_epoch)
    assert cumulative_income == 0, "ErrMsg: Last time cumulative income {}".format(
        cumulative_income)
    last_delegate_epoch, delegate_total, delegate_total_hes, delegate_reward_total = get_delegate_relevant_amount_and_epoch(
        client, node.node_id)
    assert last_delegate_epoch == 3, "ErrMsg: Last time delegate epoch {}".format(last_delegate_epoch)
    assert delegate_total == von_amount(delegate_amount, 2) - von_amount(economic.delegate_limit,
                                                                         2), "The total number of effective commissioned nodes: {}".format(
        delegate_total)
    assert delegate_total_hes == delegate_amount, "The total number of inactive nodes commissioned: {}".format(
        delegate_total_hes)
    assert delegate_reward_total == first_commission_award_total, "Total delegated rewards currently issued by the candidate: {}".format(
        delegate_reward_total)
    result = client.ppos.getDelegateReward(first_address)
    assert result['Ret'][0][
               'reward'] == first_current_commission_award, "ErrMsg: Dividends currently available {}".format(
        result['Ret'][0]['reward'])
    result = client.ppos.getDelegateReward(second_address)
    assert result['Ret'][0][
               'reward'] == 0, "ErrMsg: Dividends currently available {}".format(
        result['Ret'][0]['reward'])
    block_reward, staking_reward = economic.get_current_year_reward(node)
    economic.wait_settlement_blocknum(node)
    log.info("Current settlement block height：{}".format(node.eth.blockNumber))
    delegate_amount_total = von_amount(delegate_amount, 2) - von_amount(economic.delegate_limit, 2)
    second_commission_award_total = economic.calculate_delegate_reward(node, block_reward, staking_reward)
    second_current_commission_award = economic.delegate_cumulative_income(node, block_reward, staking_reward,
                                                                          delegate_amount_total,
                                                                          delegate_amount - economic.delegate_limit)
    result = client.delegate.withdraw_delegate_reward(first_address)
    assert_code(result, 0)
    result = client.delegate.delegate(0, second_address, amount=delegate_amount)
    assert_code(result, 0)
    delegate_epoch, cumulative_income = get_dividend_information(client, node.node_id, first_address)
    assert delegate_epoch == 4, "ErrMsg: Last time delegate epoch {}".format(delegate_epoch)
    assert cumulative_income == 0, "ErrMsg: Last time cumulative income {}".format(
        cumulative_income)
    delegate_epoch, cumulative_income = get_dividend_information(client, node.node_id, second_address)
    assert delegate_epoch == 4, "ErrMsg: Last time delegate epoch {}".format(delegate_epoch)
    assert cumulative_income == second_current_commission_award, "ErrMsg: Last time cumulative income {}".format(
        cumulative_income)
    last_delegate_epoch, delegate_total, delegate_total_hes, delegate_reward_total = get_delegate_relevant_amount_and_epoch(
        client, node.node_id)
    assert last_delegate_epoch == 4, "ErrMsg: Last time delegate epoch {}".format(last_delegate_epoch)
    assert delegate_total == von_amount(delegate_amount, 3) - von_amount(economic.delegate_limit, 2), "The total number of effective commissioned nodes: {}".format(
        delegate_total)
    assert delegate_total_hes == delegate_amount, "The total number of inactive nodes commissioned: {}".format(
        delegate_total_hes)
    assert delegate_reward_total == first_commission_award_total + second_commission_award_total, "Total delegated rewards currently issued by the candidate: {}".format(
        delegate_reward_total)
    result = client.ppos.getDelegateReward(first_address)
    assert result['Ret'][0]['reward'] == 0, "ErrMsg: Dividends currently available {}".format(
        result['Ret'][0]['reward'])
    result = client.ppos.getDelegateReward(second_address)
    assert result['Ret'][0]['reward'] == second_current_commission_award, "ErrMsg: Dividends currently available {}".format(
        result['Ret'][0]['reward'])
#
# @pytest.mark.P2
# def test_EI_BC_080(client_new_node, reset_environment):
#     client = client_new_node
#     economic = client.economic
#     node = client.node
#     address = create_restricting_plan(client)
#     log.info("Create delegate account：{}".format(address))
#     staking_address = create_staking_node(client)
#     log.info("Create pledge node id :{}".format(node.node_id))
#     delegate_amount = von_amount(economic.delegate_limit, 10)
#     result = client.delegate.delegate(0, address, amount=delegate_amount)
#     assert_code(result, 0)
#     log.info("Commissioned successfully, commissioned amount：{}".format(economic.delegate_limit))
#     economic.wait_settlement_blocknum(node)
#     log.info("Current settlement block height：{}".format(node.eth.blockNumber))
#     block_reward, staking_reward = economic.get_current_year_reward(node)
#     log.info("block_reward: {}  staking_reward: {}".format(block_reward, staking_reward))
#     economic.wait_settlement_blocknum(node)
#     log.info("Current settlement block height：{}".format(node.eth.blockNumber))
#     commission_award_total = economic.calculate_delegate_reward(node, block_reward, staking_reward, reward=1000)
#     current_commission_award = economic.delegate_cumulative_income(node, block_reward, staking_reward, delegate_amount,
#                                                                    delegate_amount, reward=1000)
#     result = client.delegate.delegate(0, address, amount=delegate_amount)
#     assert_code(result, 0)
#     delegate_epoch, cumulative_income = get_dividend_information(client, node.node_id, address)
#     assert delegate_epoch == 3, "ErrMsg: Last time delegate epoch {}".format(delegate_epoch)
#     assert cumulative_income == current_commission_award, "ErrMsg: Last time cumulative income {}".format(
#         cumulative_income)
#     last_delegate_epoch, delegate_total, delegate_total_hes, delegate_reward_total = get_delegate_relevant_amount_and_epoch(
#         client, node.node_id)
#     assert last_delegate_epoch == 3, "ErrMsg: Last time delegate epoch {}".format(last_delegate_epoch)
#     assert delegate_total == delegate_amount, "The total number of effective commissioned nodes: {}".format(
#         delegate_total)
#     assert delegate_total_hes == delegate_amount, "The total number of inactive nodes commissioned: {}".format(
#         delegate_total_hes)
#     assert delegate_reward_total == commission_award_total, "Total delegated rewards currently issued by the candidate: {}".format(
#         delegate_reward_total)
#     result = client.ppos.getDelegateReward(address)
#     assert result['Ret'][0]['reward'] == current_commission_award, "ErrMsg: Dividends currently available {}".format(
#         result['Ret'][0]['reward'])
#     block_reward, staking_reward = economic.get_current_year_reward(node)
#     economic.wait_settlement_blocknum(node)
#     log.info("Current settlement block height：{}".format(node.eth.blockNumber))
#     second_commission_award_total = economic.calculate_delegate_reward(node, block_reward, staking_reward)
#     second_current_commission_award = economic.delegate_cumulative_income(node, block_reward, staking_reward,
#                                                                           delegate_amount,
#                                                                           delegate_amount)
#     result = client.delegate.delegate(0, address, amount=delegate_amount)
#     assert_code(result, 0)
#     delegate_epoch, cumulative_income = get_dividend_information(client, node.node_id, address)
#     assert delegate_epoch == 4, "ErrMsg: Last time delegate epoch {}".format(delegate_epoch)
#     assert cumulative_income == current_commission_award + second_current_commission_award, "ErrMsg: Last time cumulative income {}".format(
#         cumulative_income)
#     last_delegate_epoch, delegate_total, delegate_total_hes, delegate_reward_total = get_delegate_relevant_amount_and_epoch(
#         client, node.node_id)
#     assert last_delegate_epoch == 4, "ErrMsg: Last time delegate epoch {}".format(last_delegate_epoch)
#     assert delegate_total == von_amount(delegate_amount,
#                                         2), "The total number of effective commissioned nodes: {}".format(
#         delegate_total)
#     assert delegate_total_hes == delegate_amount, "The total number of inactive nodes commissioned: {}".format(
#         delegate_total_hes)
#     assert delegate_reward_total == commission_award_total + second_commission_award_total, "Total delegated rewards currently issued by the candidate: {}".format(
#         delegate_reward_total)
#     result = client.ppos.getDelegateReward(address)
#     assert result['Ret'][0][
#                'reward'] == current_commission_award + second_current_commission_award, "ErrMsg: Dividends currently available {}".format(
#         result['Ret'][0]['reward'])


@pytest.mark.P1
def test_EI_BC_081(clients_new_node, reset_environment):
    """
    委托多节点，委托犹豫期多次领取
    :param clients_new_node:
    :return:
    """
    first_client = clients_new_node[0]
    second_client = clients_new_node[1]
    first_economic = first_client.economic
    first_node = first_client.node
    second_economic = second_client.economic
    second_node = second_client.node
    first_economic.env.deploy_all()
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
    first_balance = first_client.node.eth.getBalance(address)
    log.info("delegate account balance:{}".format(first_balance))
    first_client.economic.wait_consensus_blocknum(first_node)
    log.info("Current block height：{}".format(first_node.eth.blockNumber))
    result = first_client.delegate.withdraw_delegate_reward(address)
    assert_code(result, 0)
    second_balance = first_client.node.eth.getBalance(address)
    log.info("delegate account balance:{}".format(second_balance))
    gas = get_getDelegateReward_gas_fee(first_client, 2, 0)
    assert first_balance - gas == second_balance, "ErrMsg:delegate account balance {}".format(second_balance)
    result = first_client.delegate.withdraw_delegate_reward(address)
    assert_code(result, 0)
    third_balance = first_client.node.eth.getBalance(address)
    gas = get_getDelegateReward_gas_fee(first_client, 2, 0)
    assert second_balance - gas == third_balance, "ErrMsg:delegate account balance {}".format(third_balance)


@pytest.mark.P1
def test_EI_BC_082(clients_new_node, reset_environment):
    """
    委托多节点，委托锁定期多次领取
    :param clients_new_node:
    :return:
    """
    first_client = clients_new_node[0]
    second_client = clients_new_node[1]
    first_economic = first_client.economic
    first_node = first_client.node
    second_economic = second_client.economic
    second_node = second_client.node
    first_economic.env.deploy_all()
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
    first_balance = first_client.node.eth.getBalance(address)
    log.info("delegate account balance:{}".format(first_balance))
    first_client.economic.wait_settlement_blocknum(first_node)
    log.info("Current block height：{}".format(first_node.eth.blockNumber))
    result = first_client.delegate.withdraw_delegate_reward(address)
    assert_code(result, 0)
    result = first_client.ppos.getCandidateInfo(first_node.node_id)
    blocknum = result['Ret']['StakingBlockNum']
    result = first_client.ppos.getDelegateInfo(blocknum, address, first_node.node_id)
    log.info("Commission information：{}".format(result))
    delegate_epoch, cumulative_income = get_dividend_information(first_client, first_node.node_id, address)
    assert delegate_epoch == 1, "ErrMsg: Last time delegate epoch {}".format(delegate_epoch)
    assert cumulative_income == 0, "ErrMsg: Last time cumulative income {}".format(cumulative_income)
    delegate_epoch, cumulative_income = get_dividend_information(second_client, second_node.node_id, address)
    assert delegate_epoch == 1, "ErrMsg: Last time delegate epoch {}".format(delegate_epoch)
    assert cumulative_income == 0, "ErrMsg: Last time cumulative income {}".format(cumulative_income)
    second_balance = first_client.node.eth.getBalance(address)
    log.info("delegate account balance:{}".format(second_balance))
    gas = get_getDelegateReward_gas_fee(first_client, 2, 0)
    assert first_balance - gas == second_balance, "ErrMsg:delegate account balance {}".format(second_balance)
    result = first_client.delegate.withdraw_delegate_reward(address)
    assert_code(result, 0)
    third_balance = first_client.node.eth.getBalance(address)
    gas = get_getDelegateReward_gas_fee(first_client, 2, 0)
    assert second_balance - gas == third_balance, "ErrMsg:delegate account balance {}".format(third_balance)
    first_current_block = first_node.eth.blockNumber
    log.info("Current block height：{}".format(first_current_block))
    first_client.economic.wait_settlement_blocknum(first_node)
    second_current_block = first_node.eth.blockNumber
    log.info("Current block height：{}".format(second_current_block))
    assert second_current_block > first_current_block


@pytest.mark.P1
def test_EI_BC_083(clients_new_node, reset_environment):
    """
    委托多节点，委托跨结算期多次领取
    :param clients_new_node:
    :return:
    """
    first_client = clients_new_node[0]
    second_client = clients_new_node[1]
    first_economic = first_client.economic
    first_node = first_client.node
    second_economic = second_client.economic
    second_node = second_client.node
    first_economic.env.deploy_all()
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
    first_client.economic.wait_settlement_blocknum(first_node)
    log.info("Current block height：{}".format(first_node.eth.blockNumber))
    block_reward, staking_reward = first_economic.get_current_year_reward(first_node)
    first_balance = first_node.eth.getBalance(address)
    log.info("Entrusted account balance: {}".format(first_balance))
    first_economic.wait_settlement_blocknum(first_node)
    log.info("Current settlement block height：{}".format(first_node.eth.blockNumber))
    first_award_total = first_economic.calculate_delegate_reward(first_node, block_reward, staking_reward)
    first_current_commission_award = first_economic.delegate_cumulative_income(first_node, block_reward, staking_reward,
                                                                               delegate_amount, delegate_amount)
    second_award_total = first_economic.calculate_delegate_reward(second_node, block_reward, staking_reward)
    second_current_commission_award = second_economic.delegate_cumulative_income(second_node, block_reward, staking_reward,
                                                                               delegate_amount, delegate_amount)
    result = first_client.delegate.withdraw_delegate_reward(address)
    assert_code(result, 0)
    delegate_epoch, cumulative_income = get_dividend_information(first_client, first_node.node_id, address)
    assert delegate_epoch == 3, "ErrMsg: Last time delegate epoch {}".format(delegate_epoch)
    assert cumulative_income == 0, "ErrMsg: Last time cumulative income {}".format(cumulative_income)
    delegate_epoch, cumulative_income = get_dividend_information(second_client, second_node.node_id, address)
    assert delegate_epoch == 3, "ErrMsg: Last time delegate epoch {}".format(delegate_epoch)
    assert cumulative_income == 0, "ErrMsg: Last time cumulative income {}".format(cumulative_income)
    second_balance = first_client.node.eth.getBalance(address)
    log.info("delegate1 account balance:{}".format(second_balance))
    gas = get_getDelegateReward_gas_fee(first_client, 2, 2)
    log.info("gas:{}".format(gas))
    assert first_balance + first_current_commission_award + second_current_commission_award - gas == second_balance, "ErrMsg:delegate account balance {}".format(second_balance)
    result = first_client.delegate.withdraw_delegate_reward(address)
    assert_code(result, 0)
    third_balance = first_client.node.eth.getBalance(address)
    gas = get_getDelegateReward_gas_fee(first_client, 2, 0)
    log.info("gas:{}".format(gas))
    assert second_balance - gas == third_balance, "ErrMsg:delegate account balance {}".format(third_balance)


@pytest.mark.P1
def test_EI_BC_084(clients_new_node, reset_environment):
    """
    多账户委托多节点，委托跨结算期多次领取
    :param clients_new_node:
    :return:
    """
    first_client = clients_new_node[0]
    second_client = clients_new_node[1]
    first_economic = first_client.economic
    first_node = first_client.node
    second_economic = second_client.economic
    second_node = second_client.node
    first_economic.env.deploy_all()
    address = create_restricting_plan(first_client)
    address2 = create_restricting_plan(first_client)
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
    first_client.economic.wait_settlement_blocknum(first_node)
    log.info("Current block height：{}".format(first_node.eth.blockNumber))
    result = first_client.delegate.delegate(0, address2, amount=delegate_amount)
    assert_code(result, 0)
    block_reward, staking_reward = first_economic.get_current_year_reward(first_node)
    first_balance = first_node.eth.getBalance(address2)
    log.info("Entrusted account balance: {}".format(first_balance))
    first_economic.wait_settlement_blocknum(first_node)
    log.info("Current settlement block height：{}".format(first_node.eth.blockNumber))
    first_award_total = first_economic.calculate_delegate_reward(first_node, block_reward, staking_reward)
    first_current_commission_award = first_economic.delegate_cumulative_income(first_node, block_reward, staking_reward,
                                                                               delegate_amount, delegate_amount)
    result = first_client.delegate.withdraw_delegate_reward(address2)
    assert_code(result, 0)
    delegate_epoch, cumulative_income = get_dividend_information(first_client, first_node.node_id, address)
    assert delegate_epoch == 1, "ErrMsg: Last time delegate epoch {}".format(delegate_epoch)
    assert cumulative_income == 0, "ErrMsg: Last time cumulative income {}".format(cumulative_income)
    delegate_epoch, cumulative_income = get_dividend_information(first_client, first_node.node_id, address2)
    assert delegate_epoch == 2, "ErrMsg: Last time delegate epoch {}".format(delegate_epoch)
    assert cumulative_income == 0, "ErrMsg: Last time cumulative income {}".format(cumulative_income)
    second_balance = first_client.node.eth.getBalance(address2)
    log.info("delegate1 account balance:{}".format(second_balance))
    gas = get_getDelegateReward_gas_fee(first_client, 1, 0)
    assert first_balance - gas == second_balance, "ErrMsg:delegate account balance {}".format(second_balance)
    result = first_client.delegate.withdraw_delegate_reward(address2)
    assert_code(result, 0)
    third_balance = first_client.node.eth.getBalance(address2)
    gas = get_getDelegateReward_gas_fee(first_client, 1, 0)
    assert second_balance - gas == third_balance, "ErrMsg:delegate account balance {}".format(third_balance)


@pytest.mark.P1
def test_EI_BC_085(client_new_node, reset_environment):
    """
    不同结算期委托，多次领取
    :param client_new_node:
    :param reset_environment:
    :return:
    """
    client = client_new_node
    economic = client.economic
    node = client.node
    address = create_restricting_plan(client)
    log.info("Create delegate account：{}".format(address))
    create_staking_node(client)
    log.info("Create pledge node id :{}".format(node.url))
    delegate_amount = von_amount(economic.delegate_limit, 10)
    result = client.delegate.delegate(0, address, amount=delegate_amount)
    assert_code(result, 0)
    log.info("Commissioned successfully, commissioned amount：{}".format(economic.delegate_limit))
    economic.wait_settlement_blocknum(node)
    log.info("Current settlement block height：{}".format(node.eth.blockNumber))
    result = client.delegate.delegate(0, address, amount=delegate_amount)
    assert_code(result, 0)
    result = client.ppos.getCandidateInfo(node.node_id)
    blocknum = result['Ret']['StakingBlockNum']
    result = client.ppos.getDelegateInfo(blocknum, address, node.node_id)
    log.info("Commission information：{}".format(result))
    info = result['Ret']
    assert info['ReleasedHes'] == delegate_amount
    block_reward, staking_reward = economic.get_current_year_reward(node)
    economic.wait_settlement_blocknum(node)
    log.info("Current settlement block height：{}".format(node.eth.blockNumber))
    first_commission_award_total = economic.calculate_delegate_reward(node, block_reward, staking_reward)
    first_current_commission_award = economic.delegate_cumulative_income(node, block_reward, staking_reward,
                                                                         delegate_amount,
                                                                         delegate_amount)
    result = client.delegate.delegate(0, address, amount=delegate_amount)
    assert_code(result, 0)
    result = client.ppos.getDelegateInfo(blocknum, address, node.node_id)
    log.info("Commission information：{}".format(result))
    first_balance = node.eth.getBalance(address)
    log.info("delegate amount balance: {}".format(first_balance))
    result = client.delegate.withdraw_delegate_reward(address)
    assert_code(result, 0)
    delegate_epoch, cumulative_income = get_dividend_information(client, node.node_id, address)
    assert delegate_epoch == 3, "ErrMsg: Last time delegate epoch {}".format(delegate_epoch)
    assert cumulative_income == 0, "ErrMsg: Last time cumulative income {}".format(
        cumulative_income)
    second_balance = node.eth.getBalance(address)
    log.info("delegate amount balance: {}".format(second_balance))
    gas = get_getDelegateReward_gas_fee(client, 1, 0)
    assert first_balance + first_current_commission_award - gas == second_balance
    result = client.delegate.withdraw_delegate_reward(address)
    assert_code(result, 0)
    third_balance = node.eth.getBalance(address)
    log.info("delegate amount balance: {}".format(third_balance))
    assert second_balance - gas == third_balance


@pytest.mark.P1
def test_EI_BC_086(client_new_node, reset_environment):
    """
    账号全部赎回后再次领取
    :param client_new_node:
    :param reset_environment:
    :return:
    """
    client = client_new_node
    economic = client.economic
    node = client.node
    address = create_restricting_plan(client)
    log.info("Create delegate account：{}".format(address))
    create_staking_node(client)
    log.info("Create pledge node id :{}".format(node.url))
    delegate_amount = von_amount(economic.delegate_limit, 10)
    result = client.delegate.delegate(0, address, amount=delegate_amount)
    assert_code(result, 0)
    log.info("Commissioned successfully, commissioned amount：{}".format(economic.delegate_limit))
    economic.wait_settlement_blocknum(node)
    log.info("Current settlement block height：{}".format(node.eth.blockNumber))
    result = client.ppos.getCandidateInfo(node.node_id)
    blocknum = result['Ret']['StakingBlockNum']
    result = client.ppos.getDelegateInfo(blocknum, address, node.node_id)
    log.info("Commission information：{}".format(result))
    info = result['Ret']
    # assert info['ReleasedHes'] == delegate_amount
    block_reward, staking_reward = economic.get_current_year_reward(node)
    economic.wait_settlement_blocknum(node)
    log.info("Current settlement block height：{}".format(node.eth.blockNumber))
    first_commission_award_total = economic.calculate_delegate_reward(node, block_reward, staking_reward)
    first_current_commission_award = economic.delegate_cumulative_income(node, block_reward, staking_reward,
                                                                         delegate_amount,
                                                                         delegate_amount)
    result = client.delegate.withdrew_delegate(blocknum, address, amount=delegate_amount)
    assert_code(result, 0)
    first_balance = node.eth.getBalance(address)
    log.info("delegate amount balance: {}".format(first_balance))
    # delegate_epoch, cumulative_income = get_dividend_information(client, node.node_id, address)
    # assert delegate_epoch == 3, "ErrMsg: Last time delegate epoch {}".format(delegate_epoch)
    # assert cumulative_income == 0, "ErrMsg: Last time cumulative income {}".format(
    #     cumulative_income)
    second_balance = node.eth.getBalance(address)
    log.info("delegate amount balance: {}".format(second_balance))
    # assert first_balance + first_current_commission_award - second_balance < client.node.web3.toWei(1, 'ether')
    result = client.delegate.withdraw_delegate_reward(address)
    assert_code(result, 305001)



