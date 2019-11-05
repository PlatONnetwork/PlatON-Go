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
            i], "Year {} Height of block to be released: {} Release amount: {}".format(i + 1, release_plans_list[i][
            'blockNumber'], release_plans_list[i]['amount'])


def create_restrictingplan(client_new_node_obj, epoch, amount, multiple=2):
    # create restricting plan
    address, _ = client_new_node_obj.economic.account.generate_account(client_new_node_obj.node.web3,
                                                                       client_new_node_obj.economic.create_staking_limit * multiple)
    benifit_address, _ = client_new_node_obj.economic.account.generate_account(client_new_node_obj.node.web3,
                                                                               client_new_node_obj.node.web3.toWei(1000,
                                                                                                                   'ether'))
    plan = [{'Epoch': epoch, 'Amount': client_new_node_obj.node.web3.toWei(amount, 'ether')}]
    result = client_new_node_obj.restricting.createRestrictingPlan(benifit_address, plan, address)
    return result, address, benifit_address


@pytest.mark.P1
def test_LS_PV_001(client_new_node_obj):
    """
    锁仓参数的有效性验证:
                    None,
                    ""
    :param client_new_node_obj:
    :return:
    """
    # create restricting plan
    address, _ = client_new_node_obj.economic.account.generate_account(client_new_node_obj.node.web3,
                                                                       client_new_node_obj.economic.create_staking_limit)
    plan = [{'Epoch': 1, 'Amount': None}]
    try:
        result = client_new_node_obj.restricting.createRestrictingPlan(address, plan, address)
        assert_code(result, 304011)
    except Exception as e:
        log.info("Use case success, exception information：{} ".format(str(e)))

    # create restricting plan
    address, _ = client_new_node_obj.economic.account.generate_account(client_new_node_obj.node.web3,
                                                                       client_new_node_obj.economic.create_staking_limit)
    plan = [{'Epoch': 1, 'Amount': ""}]
    result = client_new_node_obj.restricting.createRestrictingPlan(address, plan, address)
    assert_code(result, 304011)


@pytest.mark.P1
def test_LS_PV_003(client_new_node_obj):
    """
    正常创建锁仓计划
    :param client_new_node_obj:
    :return:
    """
    result, address, benifit_address = create_restrictingplan(client_new_node_obj, 1, 1000)
    assert_code(result, 0)
    restricting_info = client_new_node_obj.ppos.getRestrictingInfo(benifit_address)
    assert_code(restricting_info, 0)
    assert restricting_info['Data']['balance'] == client_new_node_obj.node.web3.toWei(1000, 'ether')


@pytest.mark.P1
@pytest.mark.parametrize('epoch, amount', [(0.1, 10), (1, 0.1)])
def test_LS_PV_004_1(client_new_node_obj, epoch, amount):
    """
    锁仓参数的有效性验证:
                    number 0.1, amount 10
                    number 1, amount 0.1
    :param client_new_node_obj:
    :return:
    """
    try:
        result, address, benifit_address = create_restrictingplan(client_new_node_obj, epoch, amount)
        assert_code(result, 0)
    except Exception as e:
        log.info("Use case success, exception information：{} ".format(str(e)))


@pytest.mark.parametrize('epoch, amount', [(-1, 10), (1, -10)])
@pytest.mark.P1
def test_LS_PV_004_2(client_new_node_obj, epoch, amount):
    """
    锁仓参数的有效性验证:epoch -1, amount 10
                      epoch 1, amount -10
    :param client_new_node_obj:
    :return:
    """
    # create restricting plan
    address, _ = client_new_node_obj.economic.account.generate_account(client_new_node_obj.node.web3,
                                                                       client_new_node_obj.economic.create_staking_limit)
    plan = [{'Epoch': epoch, 'Amount': amount}]
    try:
        result = client_new_node_obj.restricting.createRestrictingPlan(address, plan, address)
        assert_code(result, 304011)
    except Exception as e:
        log.info("Use case success, exception information：{} ".format(str(e)))


@pytest.mark.P1
def test_LS_PV_005(client_new_node_obj):
    """
    锁仓参数的有效性验证:epoch 0, amount 10
    :param client_new_node_obj:
    :return:
    """
    result, address, benifit_address = create_restrictingplan(client_new_node_obj, 0, 10)
    assert_code(result, 304001)


@pytest.mark.P1
@pytest.mark.parametrize('number', [1, 5, 36])
def test_LS_PV_006(client_new_node_obj, number):
    """
    创建锁仓计划1<= 释放计划个数N <=36
    :param client_new_node_obj:
    :return:
    """
    # create restricting plan
    address, _ = client_new_node_obj.economic.account.generate_account(client_new_node_obj.node.web3,
                                                                       client_new_node_obj.economic.create_staking_limit)
    plan = []
    for i in range(number):
        plan.append({'Epoch': i + 1, 'Amount': client_new_node_obj.node.web3.toWei(10, 'ether')})
    log.info("Create lock plan parameters：{}".format(plan))
    result = client_new_node_obj.restricting.createRestrictingPlan(address, plan, address)
    assert_code(result, 0)


@pytest.mark.P1
def test_LS_PV_007(client_new_node_obj):
    """
    创建锁仓计划-释放计划的锁定期个数 > 36
    :param client_new_node_obj:
    :return:
    """
    # create restricting plan
    address, _ = client_new_node_obj.economic.account.generate_account(client_new_node_obj.node.web3,
                                                                       client_new_node_obj.economic.create_staking_limit)
    plan = []
    for i in range(37):
        plan.append({'Epoch': i + 1, 'Amount': client_new_node_obj.node.web3.toWei(10, 'ether')})
    log.info("Create lock plan parameters：{}".format(plan))
    result = client_new_node_obj.restricting.createRestrictingPlan(address, plan, address)
    assert_code(result, 304002)


@pytest.mark.P1
def test_LS_PV_008(client_new_node_obj):
    """
    锁仓参数的有效性验证:epoch 1, amount 0
    :param client_new_node_obj:
    :return:
    """
    # create restricting plan
    result, address, benifit_address = create_restrictingplan(client_new_node_obj, 1, 0)
    assert_code(result, 304011)


@pytest.mark.P2
def test_LS_PV_009(client_new_node_obj):
    """
    创建锁仓计划-锁仓金额中文、特殊字符字符测试
    :param client_new_node_obj:
    :return:
    """
    # create restricting plan
    address, _ = client_new_node_obj.economic.account.generate_account(client_new_node_obj.node.web3,
                                                                       client_new_node_obj.economic.create_staking_limit)
    plan = [{'Epoch': 1, 'Amount': '测试 @！'}]
    result = client_new_node_obj.restricting.createRestrictingPlan(address, plan, address)
    assert_code(result, 304004)


@pytest.mark.P1
def test_LS_RV_001(client_new_node_obj):
    """
    创建锁仓计划-单个释放锁定期金额大于账户金额
    :param client_new_node_obj:
    :return:
    """
    # create restricting plan
    account_balance = client_new_node_obj.node.web3.toWei(1000, 'ether')
    Lock_in_amount = client_new_node_obj.node.web3.toWei(1001, 'ether')
    address, _ = client_new_node_obj.economic.account.generate_account(client_new_node_obj.node.web3,
                                                                       client_new_node_obj.node.web3.toWei(
                                                                           account_balance, 'ether'))
    plan = [{'Epoch': 1, 'Amount': client_new_node_obj.node.web3.toWei(Lock_in_amount, 'ether')}]
    result = client_new_node_obj.restricting.createRestrictingPlan(address, plan, address)
    assert_code(result, 304004)


@pytest.mark.P1
@pytest.mark.parametrize('balace1, balace2', [(0, 0), (300, 300), (500, 500), (500, 600)])
def test_LS_RV_002(client_new_node_obj, balace1, balace2):
    """
    创建锁仓计划-多个释放锁定期合计金额大于账户金额
    :param client_new_node_obj:
    :return:
    """
    # create restricting plan
    address, _ = client_new_node_obj.economic.account.generate_account(client_new_node_obj.node.web3,
                                                                       client_new_node_obj.node.web3.toWei(1000,
                                                                                                           'ether'))
    louk_up_balace1 = client_new_node_obj.node.web3.toWei(balace1, 'ether')
    louk_up_balace2 = client_new_node_obj.node.web3.toWei(balace2, 'ether')
    plan = [{'Epoch': 1, 'Amount': louk_up_balace1}, {'Epoch': 2, 'Amount': louk_up_balace2}]
    result = client_new_node_obj.restricting.createRestrictingPlan(address, plan, address)
    if 0 < balace1 + balace2 < 1000:
        assert_code(result, 0)
    elif 1000 <= balace1 + balace2:
        assert_code(result, 304004)
    else:
        assert_code(result, 304011)


@pytest.mark.P1
def test_LS_RV_003(client_new_node_obj):
    """
    创建锁仓计划-锁仓计划里两个锁仓计划的解锁期相同
    :param client_new_node_obj:
    :return:
    """
    # create restricting plan
    louk_up_balace = client_new_node_obj.node.web3.toWei(100, 'ether')
    address, _ = client_new_node_obj.economic.account.generate_account(client_new_node_obj.node.web3,
                                                                       client_new_node_obj.node.web3.toWei(1000,
                                                                                                           'ether'))
    plan = [{'Epoch': 1, 'Amount': louk_up_balace}, {'Epoch': 1, 'Amount': louk_up_balace}]
    result = client_new_node_obj.restricting.createRestrictingPlan(address, plan, address)
    assert_code(result, 0)
    restricting_info = client_new_node_obj.ppos.getRestrictingInfo(address)
    assert_code(restricting_info, 0)
    # assert restricting plan
    assert restricting_info['Data']['balance'] == louk_up_balace * 2, "ErrMsg:Restricting balance：{}".format(
        restricting_info['Data']['balance'])
    assert restricting_info['Data']['plans'][0][
               'blockNumber'] == client_new_node_obj.economic.get_settlement_switchpoint(
        client_new_node_obj.node), "ErrMsg:Restricting blockNumber {}".format(
        restricting_info['Data']['plans'][0]['blockNumber'])
    assert restricting_info['Data']['plans'][0][
               'amount'] == louk_up_balace * 2, "ErrMsg:Restricting amount {}".format(
        restricting_info['Data']['plans'][0]['amount'])


@pytest.mark.P1
def test_LS_RV_004(client_new_node_obj):
    """
    创建锁仓计划-新建锁仓计划里两个锁仓计划的解锁期不同
    :param client_new_node_obj:
    :return:
    """
    # create restricting plan
    louk_up_balace = client_new_node_obj.node.web3.toWei(100, 'ether')
    address, _ = client_new_node_obj.economic.account.generate_account(client_new_node_obj.node.web3,
                                                                       client_new_node_obj.node.web3.toWei(1000,
                                                                                                           'ether'))
    plan = [{'Epoch': 1, 'Amount': louk_up_balace}, {'Epoch': 2, 'Amount': louk_up_balace}]
    # create restricting
    result = client_new_node_obj.restricting.createRestrictingPlan(address, plan, address)
    assert_code(result, 0)
    # view restricting plan
    restricting_info = client_new_node_obj.ppos.getRestrictingInfo(address)
    log.info("Restricting information: {}".format(restricting_info))
    assert_code(restricting_info, 0)
    # assert restricting plan
    assert restricting_info['Data']['balance'] == louk_up_balace * 2, "ErrMsg:Restricting balance：{}".format(
        restricting_info['Data']['balance'])
    assert restricting_info['Data']['plans'][0][
               'blockNumber'] == client_new_node_obj.economic.get_settlement_switchpoint(
        client_new_node_obj.node), "ErrMsg:Restricting blockNumber {}".format(
        restricting_info['Data']['plans'][0]['blockNumber'])
    assert restricting_info['Data']['plans'][0][
               'amount'] == louk_up_balace, "ErrMsg:Restricting amount {}".format(
        restricting_info['Data']['plans'][0]['amount'])
    assert restricting_info['Data']['plans'][1][
               'amount'] == louk_up_balace, "ErrMsg:Restricting amount {}".format(
        restricting_info['Data']['plans'][1]['amount'])
