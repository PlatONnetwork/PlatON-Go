import time
import pytest
import allure
import rlp
from client_sdk_python.utils.transactions import send_obj_transaction
from dacite import from_dict
from platon_account.internal.transactions import bech32_address_bytes

from common.key import get_pub_key, mock_duplicate_sign
from common.log import log
from client_sdk_python import Web3
from decimal import Decimal
from tests.lib import EconomicConfig, Genesis, StakingConfig, Staking, check_node_in_list, assert_code, von_amount, \
    get_governable_parameter_value, get_the_dynamic_parameter_gas_fee


@pytest.mark.P0
@pytest.mark.compatibility
def test_LS_FV_001(client_consensus):
    """
    查看锁仓账户计划
    :param client_consensus:
    :return:
    """
    # Reset environment
    client_consensus.economic.env.deploy_all()
    # view Lock in contract amount
    lock_up_amount = client_consensus.node.eth.getBalance(EconomicConfig.FOUNDATION_LOCKUP_ADDRESS)
    log.info("Lock in contract amount: {}".format(lock_up_amount))
    # view Lockup plan
    result = client_consensus.ppos.getRestrictingInfo(EconomicConfig.INCENTIVEPOOL_ADDRESS)
    assert_code(result, 304005)
    # release_plans_list = result['Ret']['plans']
    # assert_code(result, 0)
    # log.info("Lockup plan information: {}".format(result))
    # # assert louck up amount
    # for i in range(len(release_plans_list)):
    #     assert release_plans_list[i] == EconomicConfig.release_info[
    #         i], "Year {} Height of block to be released: {} Release amount: {}".format(i + 1, release_plans_list[i][
    #             'blockNumber'], release_plans_list[i]['amount'])


def create_restrictingplan(client, epoch, amount, multiple=2):
    # create restricting plan
    address, _ = client.economic.account.generate_account(client.node.web3,
                                                          client.economic.create_staking_limit * multiple)
    benifit_address, _ = client.economic.account.generate_account(client.node.web3,
                                                                  client.node.web3.toWei(1000,
                                                                                         'ether'))
    plan = [{'Epoch': epoch, 'Amount': client.node.web3.toWei(amount, 'ether')}]
    result = client.restricting.createRestrictingPlan(benifit_address, plan, address)
    return result, address, benifit_address


@pytest.mark.P1
@pytest.mark.compatibility
def test_LS_UPV_001(client_new_node):
    """
    锁仓参数的有效性验证:
                    None,
                    ""
    :param client_new_node:
    :return:
    """
    client = client_new_node
    economic = client.economic
    node = client.node
    status = True
    # create restricting plan
    address, _ = economic.account.generate_account(node.web3, economic.create_staking_limit)
    plan = [{'Epoch': 1, 'Amount': ""}]
    result = client.restricting.createRestrictingPlan(address, plan, address)
    assert_code(result, 304011)

    # create restricting plan
    address, _ = economic.account.generate_account(node.web3, economic.create_staking_limit)
    plan = [{'Epoch': 1, 'Amount': None}]
    try:
        client.restricting.createRestrictingPlan(address, plan, address)
        status = False
    except Exception as e:
        log.info("Use case success, exception information：{} ".format(str(e)))
    assert status, "ErrMsg: create restricting result {}".format(status)


@pytest.mark.P2
def test_LS_UPV_002_1(client_new_node):
    """
    创建锁仓计划Gas费- 单个解锁次数
    :param client_new_node:
    :return:
    """
    client = client_new_node
    economic = client.economic
    node = client.node
    # create restricting plan
    address, _ = economic.account.generate_account(node.web3, economic.create_staking_limit)
    lock_amount = node.web3.toWei(1000, 'ether')
    plan = [{'Epoch': 1, 'Amount': lock_amount}]
    account = bech32_address_bytes(address)
    plan_list = []
    for dict_ in plan:
        v = [dict_[k] for k in dict_]
        plan_list.append(v)
    rlp_list = rlp.encode(plan_list)
    data = rlp.encode([rlp.encode(int(4000)), rlp.encode(account), rlp_list])
    transaction_data = {"to": address, "data": data}
    aa = node.eth.estimateGas(transaction_data)
    dynamic_gas = get_the_dynamic_parameter_gas_fee(data)
    gas_total = 21000 + 18000 + 8000 + 21000 + dynamic_gas
    log.info("gas_total: {}".format(gas_total))
    balance = node.eth.getBalance(address)
    # Create a lockout plan
    result = client.restricting.createRestrictingPlan(address, plan, address)
    assert_code(result, 0)
    balance1 = node.eth.getBalance(address)
    transaction_fees = gas_total * node.eth.gasPrice
    assert balance - balance1 - lock_amount == transaction_fees, "ErrMsg: transaction fees {}".format(transaction_fees)


@pytest.mark.P2
def test_LS_UPV_002_2(client_new_node):
    """
    创建锁仓计划Gas费 - 多个解锁次数
    :param client_new_node:
    :return:
    """
    client = client_new_node
    economic = client.economic
    node = client.node
    # create restricting plan
    address, _ = economic.account.generate_account(node.web3, economic.create_staking_limit)
    lock_amount = node.web3.toWei(1000, 'ether')
    plan = [{'Epoch': 1, 'Amount': lock_amount}, {'Epoch': 2, 'Amount': lock_amount}]
    account = bech32_address_bytes(address)
    plan_list = []
    for dict_ in plan:
        v = [dict_[k] for k in dict_]
        plan_list.append(v)
    rlp_list = rlp.encode(plan_list)
    data = rlp.encode([rlp.encode(int(4000)), rlp.encode(account), rlp_list])
    dynamic_gas = get_the_dynamic_parameter_gas_fee(data)
    gas_total = 21000 + 18000 + 8000 + 21000 * 2 + dynamic_gas
    log.info("gas_total: {}".format(gas_total))
    balance = node.eth.getBalance(address)
    # Create a lockout plan
    result = client.restricting.createRestrictingPlan(address, plan, address)
    assert_code(result, 0)
    balance1 = node.eth.getBalance(address)
    transaction_fees = gas_total * node.eth.gasPrice
    assert balance - balance1 - lock_amount * 2 == transaction_fees, "ErrMsg: transaction fees {}".format(
        transaction_fees)


@pytest.mark.P1
@pytest.mark.compatibility
def test_LS_UPV_003(client_new_node):
    """
    正常创建锁仓计划
    :param client_new_node:
    :return:
    """
    result, address, benifit_address = create_restrictingplan(client_new_node, 1, 1000)
    assert_code(result, 0)
    restricting_info = client_new_node.ppos.getRestrictingInfo(benifit_address)
    assert_code(restricting_info, 0)
    assert restricting_info['Ret']['balance'] == client_new_node.node.web3.toWei(1000, 'ether')


@pytest.mark.P1
def test_LS_UPV_004_1(client_new_node):
    """
    锁仓参数的有效性验证:number 1, amount 0.1
                      number 0.1, amount 10
    :param client_new_node:
    :return:
    """
    status = True
    # number 1, amount 0.1
    result, address, benifit_address = create_restrictingplan(client_new_node, 1, 0.1)
    log.info('result: {}'.format(result))
    assert_code(result, 304003)
    # number 0.1, amount 10
    try:
        result = create_restrictingplan(client_new_node, 0.1, 10)
        log.info('result: {}'.format(result))
        status = False
    except Exception as e:
        log.info("Use case success, exception information：{} ".format(str(e)))
    assert status, "ErrMsg: create restricting result {}".format(status)


@pytest.mark.parametrize('epoch, amount', [(-1, 10), (1, -10)])
@pytest.mark.P1
def test_LS_UPV_004_2(client_new_node, epoch, amount):
    """
    锁仓参数的有效性验证:epoch -1, amount 10
                      epoch 1, amount -10
    :param client_new_node:
    :return:
    """
    status = True
    # create restricting plan
    address, _ = client_new_node.economic.account.generate_account(client_new_node.node.web3,
                                                                   client_new_node.economic.create_staking_limit)
    plan = [{'Epoch': epoch, 'Amount': amount}]
    try:
        client_new_node.restricting.createRestrictingPlan(address, plan, address)
        status = False
    except Exception as e:
        log.info("Use case success, exception information：{} ".format(str(e)))
    assert status, "ErrMsg: create restricting result {}".format(status)


@pytest.mark.P1
def test_LS_UPV_005(client_new_node):
    """
    锁仓参数的有效性验证:epoch 0, amount 10
    :param client_new_node:
    :return:
    """
    result, address, benifit_address = create_restrictingplan(client_new_node, 0, 10)
    assert_code(result, 304001)


@pytest.mark.P1
@pytest.mark.parametrize('number', [1, 5, 36])
def test_LS_UPV_006(client_new_node, number):
    """
    创建锁仓计划1<= 释放计划个数N <=36
    :param client_new_node:
    :return:
    """
    # create restricting plan
    address, _ = client_new_node.economic.account.generate_account(client_new_node.node.web3,
                                                                   client_new_node.economic.create_staking_limit)
    plan = []
    for i in range(number):
        plan.append({'Epoch': i + 1, 'Amount': client_new_node.node.web3.toWei(10, 'ether')})
    log.info("Create lock plan parameters：{}".format(plan))
    result = client_new_node.restricting.createRestrictingPlan(address, plan, address)
    assert_code(result, 0)


@pytest.mark.P1
def test_LS_UPV_007(client_new_node):
    """
    创建锁仓计划-释放计划的锁定期个数 > 36
    :param client_new_node:
    :return:
    """
    # create restricting plan
    address, _ = client_new_node.economic.account.generate_account(client_new_node.node.web3,
                                                                   client_new_node.economic.create_staking_limit)
    plan = []
    for i in range(37):
        plan.append({'Epoch': i + 1, 'Amount': client_new_node.node.web3.toWei(10, 'ether')})
    log.info("Create lock plan parameters：{}".format(plan))
    result = client_new_node.restricting.createRestrictingPlan(address, plan, address)
    assert_code(result, 304002)


@pytest.mark.P1
def test_LS_UPV_008(client_new_node):
    """
    锁仓参数的有效性验证:epoch 1, amount 0
    :param client_new_node:
    :return:
    """
    # create restricting plan
    result, address, benifit_address = create_restrictingplan(client_new_node, 1, 0)
    assert_code(result, 304011)


@pytest.mark.P2
def test_LS_UPV_009(client_new_node):
    """
    创建锁仓计划-锁仓金额中文、特殊字符字符测试
    :param client_new_node:
    :return:
    """
    # create restricting plan
    address, _ = client_new_node.economic.account.generate_account(client_new_node.node.web3,
                                                                   client_new_node.economic.create_staking_limit)
    plan = [{'Epoch': 1, 'Amount': '测试 @！'}]
    result = client_new_node.restricting.createRestrictingPlan(address, plan, address)
    assert_code(result, 304004)


@pytest.mark.P1
def test_LS_RV_001(client_new_node):
    """
    创建锁仓计划-单个释放锁定期金额大于账户金额
    :param client_new_node:
    :return:
    """
    client = client_new_node
    economic = client.economic
    node = client.node
    # create restricting plan
    address, _ = economic.account.generate_account(node.web3, node.web3.toWei(1000, 'ether'))
    plan = [{'Epoch': 1, 'Amount': node.web3.toWei(1001, 'ether')}]
    result = client_new_node.restricting.createRestrictingPlan(address, plan, address)
    assert_code(result, 304004)


@pytest.mark.P1
@pytest.mark.parametrize('balance1, balance2', [(0, 0), (300, 300), (500, 500), (500, 600)])
def test_LS_RV_002(client_new_node, balance1, balance2):
    """
    创建锁仓计划-多个释放锁定期合计金额大于账户金额
    :param client_new_node:
    :return:
    """
    # create restricting plan
    address, _ = client_new_node.economic.account.generate_account(client_new_node.node.web3,
                                                                   client_new_node.node.web3.toWei(1000,
                                                                                                   'ether'))
    lock_up_balance1 = client_new_node.node.web3.toWei(balance1, 'ether')
    lock_up_balance2 = client_new_node.node.web3.toWei(balance2, 'ether')
    plan = [{'Epoch': 1, 'Amount': lock_up_balance1}, {'Epoch': 2, 'Amount': lock_up_balance2}]
    result = client_new_node.restricting.createRestrictingPlan(address, plan, address)
    if 0 < balance1 + balance2 < 1000:
        assert_code(result, 0)
    elif 1000 <= balance1 + balance2:
        assert_code(result, 304004)
    else:
        assert_code(result, 304011)


def create_restricting_plan(client, plan, benifit_address, address):
    """
    create restricting plan
    :param client:
    :param plan:
    :param benifit_address:
    :param address:
    :return:
    """
    # create restricting
    result = client.restricting.createRestrictingPlan(benifit_address, plan, address)
    assert_code(result, 0)
    # view restricting plan
    restricting_info = client.ppos.getRestrictingInfo(benifit_address)
    log.info("Restricting information: {}".format(restricting_info))
    assert_code(restricting_info, 0)
    return restricting_info


@pytest.mark.P1
def test_LS_RV_003(client_new_node):
    """
    创建锁仓计划-锁仓计划里两个锁仓计划的解锁期相同
    :param client_new_node:
    :return:
    """
    # create account
    address, _ = client_new_node.economic.account.generate_account(client_new_node.node.web3,
                                                                   client_new_node.node.web3.toWei(1000,
                                                                                                   'ether'))
    louk_up_balace = client_new_node.node.web3.toWei(100, 'ether')
    plan = [{'Epoch': 1, 'Amount': louk_up_balace}, {'Epoch': 1, 'Amount': louk_up_balace}]
    # create restricting plan
    restricting_info = create_restricting_plan(client_new_node, plan, address, address)
    # assert restricting plan
    assert restricting_info['Ret']['balance'] == louk_up_balace * 2, "ErrMsg:Restricting balance：{}".format(
        restricting_info['Ret']['balance'])
    assert restricting_info['Ret']['plans'][0][
        'blockNumber'] == client_new_node.economic.get_settlement_switchpoint(
        client_new_node.node), "ErrMsg:Restricting blockNumber {}".format(
        restricting_info['Ret']['plans'][0]['blockNumber'])
    assert restricting_info['Ret']['plans'][0][
        'amount'] == louk_up_balace * 2, "ErrMsg:Restricting amount {}".format(
        restricting_info['Ret']['plans'][0]['amount'])


@pytest.mark.P1
def test_LS_RV_004(client_new_node):
    """
    创建锁仓计划-新建锁仓计划里两个锁仓计划的解锁期不同
    :param client_new_node:
    :return:
    """
    address, _ = client_new_node.economic.account.generate_account(client_new_node.node.web3,
                                                                   client_new_node.node.web3.toWei(1000,
                                                                                                   'ether'))

    louk_up_balace = client_new_node.node.web3.toWei(100, 'ether')
    plan = [{'Epoch': 1, 'Amount': louk_up_balace}, {'Epoch': 2, 'Amount': louk_up_balace}]
    # create restricting plan
    restricting_info = create_restricting_plan(client_new_node, plan, address, address)
    # assert restricting plan
    assert restricting_info['Ret']['balance'] == louk_up_balace * 2, "ErrMsg:Restricting balance：{}".format(
        restricting_info['Ret']['balance'])
    assert restricting_info['Ret']['plans'][0][
        'blockNumber'] == client_new_node.economic.get_settlement_switchpoint(
        client_new_node.node), "ErrMsg:Restricting blockNumber {}".format(
        restricting_info['Ret']['plans'][0]['blockNumber'])
    assert restricting_info['Ret']['plans'][0][
        'amount'] == louk_up_balace, "ErrMsg:Restricting amount {}".format(
        restricting_info['Ret']['plans'][0]['amount'])
    assert restricting_info['Ret']['plans'][1][
        'amount'] == louk_up_balace, "ErrMsg:Restricting amount {}".format(
        restricting_info['Ret']['plans'][1]['amount'])


@pytest.mark.P1
def test_LS_RV_005(client_new_node):
    """
    创建锁仓计划-创建不同锁仓计划里2个相同解锁期
    :param client_new_node:
    :return:
    """
    # create account
    address, _ = client_new_node.economic.account.generate_account(client_new_node.node.web3,
                                                                   client_new_node.node.web3.toWei(1000,
                                                                                                   'ether'))

    louk_up_balace = client_new_node.node.web3.toWei(100, 'ether')
    plan = [{'Epoch': 1, 'Amount': louk_up_balace}]
    # create restricting plan
    restricting_info = create_restricting_plan(client_new_node, plan, address, address)
    # create restricting plan
    restricting_info = create_restricting_plan(client_new_node, plan, address, address)
    # assert restricting plan
    assert restricting_info['Ret']['balance'] == louk_up_balace * 2, "ErrMsg:Restricting balance：{}".format(
        restricting_info['Ret']['balance'])
    assert restricting_info['Ret']['plans'][0][
        'blockNumber'] == client_new_node.economic.get_settlement_switchpoint(
        client_new_node.node), "ErrMsg:Restricting blockNumber {}".format(
        restricting_info['Ret']['plans'][0]['blockNumber'])
    assert restricting_info['Ret']['plans'][0][
        'amount'] == louk_up_balace * 2, "ErrMsg:Restricting amount {}".format(
        restricting_info['Ret']['plans'][0]['amount'])


def create_lock_release_amount(client, first_amount, second_amount):
    # create first_address
    first_address, _ = client.economic.account.generate_account(client.node.web3, first_amount)
    # create second_address
    second_address, _ = client.economic.account.generate_account(client.node.web3, second_amount)
    return first_address, second_address


@pytest.mark.P1
def test_LS_RV_006(client_new_node):
    """
    创建锁仓计划-不同个账户创建不同锁仓计划里有相同解锁期
    :param client_new_node:
    :return:
    """
    # create account
    amount1 = client_new_node.node.web3.toWei(1000, 'ether')
    amount2 = client_new_node.node.web3.toWei(1000, 'ether')
    address1, address2 = create_lock_release_amount(client_new_node, amount1, amount2)
    louk_up_balace = client_new_node.node.web3.toWei(100, 'ether')
    plan = [{'Epoch': 1, 'Amount': louk_up_balace}, {'Epoch': 2, 'Amount': louk_up_balace}]
    # create restricting plan1
    restricting_info = create_restricting_plan(client_new_node, plan, address1, address1)
    # create restricting plan2
    restricting_info = create_restricting_plan(client_new_node, plan, address1, address2)
    # assert restricting plan1
    assert restricting_info['Ret']['balance'] == louk_up_balace * 4, "ErrMsg:Restricting balance：{}".format(
        restricting_info['Ret']['balance'])
    assert restricting_info['Ret']['plans'][0][
        'blockNumber'] == client_new_node.economic.get_settlement_switchpoint(
        client_new_node.node), "ErrMsg:Restricting blockNumber {}".format(
        restricting_info['Ret']['plans'][0]['blockNumber'])
    assert restricting_info['Ret']['plans'][0][
        'amount'] == louk_up_balace * 2, "ErrMsg:Restricting amount {}".format(
        restricting_info['Ret']['plans'][0]['amount'])
    assert restricting_info['Ret']['plans'][1][
        'amount'] == louk_up_balace * 2, "ErrMsg:Restricting amount {}".format(
        restricting_info['Ret']['plans'][1]['amount'])


@pytest.mark.P1
def test_LS_RV_007(client_new_node):
    """
    创建锁仓计划-不同账户创建不同锁仓计划里有不相同解锁期
    :param client_new_node:
    :return:
    """
    # create account
    amount1 = client_new_node.node.web3.toWei(1000, 'ether')
    amount2 = client_new_node.node.web3.toWei(1000, 'ether')
    address1, address2 = create_lock_release_amount(client_new_node, amount1, amount2)
    louk_up_balace = client_new_node.node.web3.toWei(100, 'ether')
    plan1 = [{'Epoch': 1, 'Amount': louk_up_balace}, {'Epoch': 2, 'Amount': louk_up_balace}]
    plan2 = [{'Epoch': 1, 'Amount': louk_up_balace}, {'Epoch': 3, 'Amount': louk_up_balace}]
    # create restricting plan1
    restricting_info = create_restricting_plan(client_new_node, plan1, address1, address1)
    # create restricting plan2
    restricting_info = create_restricting_plan(client_new_node, plan2, address1, address2)
    # assert restricting plan1
    assert restricting_info['Ret']['balance'] == louk_up_balace * 4, "ErrMsg:Restricting balance：{}".format(
        restricting_info['Ret']['balance'])
    assert restricting_info['Ret']['plans'][0][
        'blockNumber'] == client_new_node.economic.get_settlement_switchpoint(
        client_new_node.node), "ErrMsg:Restricting blockNumber {}".format(
        restricting_info['Ret']['plans'][0]['blockNumber'])
    assert restricting_info['Ret']['plans'][0][
        'amount'] == louk_up_balace * 2, "ErrMsg:Restricting amount {}".format(
        restricting_info['Ret']['plans'][0]['amount'])
    assert restricting_info['Ret']['plans'][1][
        'amount'] == louk_up_balace, "ErrMsg:Restricting amount {}".format(
        restricting_info['Ret']['plans'][1]['amount'])
    assert restricting_info['Ret']['plans'][2][
        'amount'] == louk_up_balace, "ErrMsg:Restricting amount {}".format(
        restricting_info['Ret']['plans'][2]['amount'])


def create_restricting_plan_and_staking(client, economic, node):
    # create account
    amount1 = von_amount(economic.create_staking_limit, 4)
    amount2 = client.node.web3.toWei(1000, 'ether')
    address1, address2 = create_lock_release_amount(client, amount1, amount2)
    # create Restricting Plan
    plan = [{'Epoch': 1, 'Amount': economic.create_staking_limit}]
    result = client.restricting.createRestrictingPlan(address2, plan, address1)
    assert_code(result, 0)
    # create staking
    result = client.staking.create_staking(1, address2, address2)
    assert_code(result, 0)
    # view Restricting Plan
    restricting_info1 = client.ppos.getRestrictingInfo(address2)
    log.info("restricting info: {}".format(restricting_info1))
    assert_code(restricting_info1, 0)
    info = restricting_info1['Ret']
    assert info['Pledge'] == economic.create_staking_limit, 'ErrMsg: restricting Pledge amount {}'.format(
        info['Pledge'])
    # wait settlement block
    economic.wait_settlement_blocknum(node)
    restricting_info2 = client.ppos.getRestrictingInfo(address2)
    log.info("current block: {}".format(node.block_number))
    log.info("restricting info: {}".format(restricting_info2))
    info = restricting_info2['Ret']
    assert info['debt'] == economic.create_staking_limit, 'ErrMsg: restricting debt amount {}'.format(
        info['debt'])
    return address1, address2


@pytest.mark.P1
def test_LS_RV_008(client_new_node):
    """
    创建锁仓计划-锁仓欠释放金额<新增锁仓计划总金额
    :param client_new_node:
    :return:
    """
    client = client_new_node
    economic = client.economic
    node = client.node
    address1, address2 = create_restricting_plan_and_staking(client, economic, node)
    # create Restricting Plan again
    plan = [{'Epoch': 1, 'Amount': von_amount(economic.create_staking_limit, 2)}]
    result = client.restricting.createRestrictingPlan(address2, plan, address1)
    assert_code(result, 0)
    # view Restricting Plan
    restricting_info = client.ppos.getRestrictingInfo(address2)
    log.info("restricting info: {}".format(restricting_info))
    assert_code(restricting_info, 0)
    info = restricting_info['Ret']
    assert info['debt'] == 0, "rrMsg: restricting debt amount {}".format(info['debt'])


@pytest.mark.P1
def test_LS_RV_009(client_new_node):
    """
    创建锁仓计划-锁仓欠释放金额>新增锁仓计划总金额
    :param client_new_node:
    :return:
    """
    client = client_new_node
    economic = client.economic
    node = client.node
    address1, address2 = create_restricting_plan_and_staking(client, economic, node)
    # create Restricting Plan again
    plan = [{'Epoch': 1, 'Amount': von_amount(economic.create_staking_limit, 0.8)}]
    result = client.restricting.createRestrictingPlan(address2, plan, address1)
    assert_code(result, 0)
    # view Restricting Plan
    restricting_info = client.ppos.getRestrictingInfo(address2)
    log.info("restricting info: {}".format(restricting_info))
    assert_code(restricting_info, 0)
    info = restricting_info['Ret']
    assert info['debt'] == economic.create_staking_limit - von_amount(economic.create_staking_limit,
                                                                      0.8), "rrMsg: restricting debt amount {}".format(
        info['debt'])


@pytest.mark.P1
def test_LS_RV_010(client_new_node):
    """
    创建锁仓计划-锁仓欠释放金额=新增锁仓计划总金额
    :param client_new_node:
    :return:
    """
    client = client_new_node
    economic = client.economic
    node = client.node
    address1, address2 = create_restricting_plan_and_staking(client, economic, node)
    # create Restricting Plan again
    plan = [{'Epoch': 1, 'Amount': von_amount(economic.create_staking_limit, 1)}]
    result = client.restricting.createRestrictingPlan(address2, plan, address1)
    assert_code(result, 0)
    # view Restricting Plan
    restricting_info = client.ppos.getRestrictingInfo(address2)
    log.info("restricting info: {}".format(restricting_info))
    assert_code(restricting_info, 0)
    info = restricting_info['Ret']
    assert info['debt'] == 0, "rrMsg: restricting debt amount {}".format(info['debt'])


def create_restricting_plan_and_entrust(client, node, economic):
    # create account
    amount1 = von_amount(economic.create_staking_limit, 2)
    amount2 = client.node.web3.toWei(1000, 'ether')
    address1, address2 = create_lock_release_amount(client, amount1, amount2)
    # create Restricting Plan
    plan = [{'Epoch': 1, 'Amount': von_amount(economic.delegate_limit, 1)}]
    result = client.restricting.createRestrictingPlan(address2, plan, address1)
    assert_code(result, 0)
    # create staking
    result = client.staking.create_staking(0, address1, address1)
    assert_code(result, 0)
    # Application for Commission
    result = client.delegate.delegate(1, address2)
    assert_code(result, 0)
    # view Restricting Plan
    restricting_info = client.ppos.getRestrictingInfo(address2)
    log.info("restricting info: {}".format(restricting_info))
    assert_code(restricting_info, 0)
    info = restricting_info['Ret']
    assert info['Pledge'] == economic.delegate_limit, 'ErrMsg: restricting Pledge amount {}'.format(
        info['Pledge'])
    # wait settlement block
    economic.wait_settlement_blocknum(node)
    log.info("current block: {}".format(node.block_number))
    # view Restricting Plan
    restricting_info = client.ppos.getRestrictingInfo(address2)
    log.info("restricting info: {}".format(restricting_info))
    assert_code(restricting_info, 0)
    info = restricting_info['Ret']
    assert info['debt'] == economic.delegate_limit, 'ErrMsg: restricting debt amount {}'.format(
        info['debt'])
    return address1, address2


@pytest.mark.P1
def test_LS_RV_011(client_new_node):
    """
    创建锁仓计划-锁仓委托释放后再次创建锁仓计划
    :param client_new_node:
    :return:
    """
    client = client_new_node
    economic = client.economic
    node = client.node
    address1, address2 = create_restricting_plan_and_entrust(client, node, economic)
    # create Restricting Plan again
    plan = [{'Epoch': 1, 'Amount': von_amount(economic.delegate_limit, 2)}]
    result = client.restricting.createRestrictingPlan(address2, plan, address1)
    assert_code(result, 0)
    # view Restricting Plan
    restricting_info = client.ppos.getRestrictingInfo(address2)
    log.info("restricting info: {}".format(restricting_info))
    assert_code(restricting_info, 0)
    info = restricting_info['Ret']
    assert info['debt'] == 0, "rrMsg: restricting debt amount {}".format(info['debt'])


@pytest.mark.P1
def test_LS_RV_012(new_genesis_env, clients_new_node):
    """
    创建锁仓计划-锁仓质押释放后被处罚再次创建锁仓计划(处罚金额大于质押金额)
    :param:client_new_node:
    :return:
    """
    # Change configuration parameters
    genesis = from_dict(data_class=Genesis, data=new_genesis_env.genesis_config)
    genesis.economicModel.slashing.slashBlocksReward = 30
    new_file = new_genesis_env.cfg.env_tmp + "/genesis_0.13.0.json"
    genesis.to_file(new_file)
    new_genesis_env.deploy_all(new_file)

    client1 = clients_new_node[0]
    log.info("Current linked client1: {}".format(client1.node.node_mark))
    client2 = clients_new_node[1]
    log.info("Current linked client2: {}".format(client2.node.node_mark))
    economic = client1.economic
    node = client1.node
    # create restricting plan and staking
    address1, address2 = create_restricting_plan_and_staking(client1, economic, node)
    # view
    candidate_info = client1.ppos.getCandidateInfo(node.node_id)
    pledge_amount = candidate_info['Ret']['Shares']
    log.info("pledge_amount: {}".format(pledge_amount))
    # Obtain pledge reward and block out reward
    block_reward, staking_reward = client1.economic.get_current_year_reward(node)
    log.info("block_reward: {} staking_reward: {}".format(block_reward, staking_reward))
    # Get 0 block rate penalties
    slash_blocks = get_governable_parameter_value(client1, 'slashBlocksReward')
    log.info("Current block height: {}".format(client2.node.eth.blockNumber))
    # stop node
    node.stop()
    # Waiting 2 consensus block
    client2.economic.wait_consensus_blocknum(client2.node, 3)
    log.info("Current block height: {}".format(client2.node.eth.blockNumber))
    # view verifier list
    verifier_list = client2.ppos.getVerifierList()
    log.info("Current settlement cycle verifier list: {}".format(verifier_list))
    # Amount of penalty
    punishment_amonut = int(Decimal(str(block_reward)) * Decimal(str(slash_blocks)))
    log.info("punishment_amonut: {}".format(punishment_amonut))
    # view Restricting Plan
    restricting_info = client2.ppos.getRestrictingInfo(address2)
    log.info("restricting info: {}".format(restricting_info))
    assert_code(restricting_info, 304005)


@pytest.mark.P1
def test_LS_RV_019(new_genesis_env, clients_noconsensus):
    """
    创建锁仓计划-锁仓质押释放后被处罚再次创建锁仓计划
    :param new_genesis_env:
    :return:
    """
    # Change configuration parameters
    genesis = from_dict(data_class=Genesis, data=new_genesis_env.genesis_config)
    genesis.economicModel.slashing.slashBlocksReward = 5
    new_file = new_genesis_env.cfg.env_tmp + "/genesis_0.13.0.json"
    genesis.to_file(new_file)
    new_genesis_env.deploy_all(new_file)

    client1 = clients_noconsensus[0]
    log.info("Current linked client1: {}".format(client1.node.node_mark))
    client2 = clients_noconsensus[1]
    log.info("Current linked client2: {}".format(client2.node.node_mark))
    economic = client1.economic
    node = client1.node
    # create restricting plan and staking
    address1, address2 = create_restricting_plan_and_staking(client1, economic, node)
    # view
    candidate_info = client1.ppos.getCandidateInfo(node.node_id)
    pledge_amount = candidate_info['Ret']['Shares']
    log.info("pledge_amount: {}".format(pledge_amount))
    # Obtain pledge reward and block out reward
    block_reward, staking_reward = client1.economic.get_current_year_reward(node)
    log.info("block_reward: {} staking_reward: {}".format(block_reward, staking_reward))
    # Get 0 block rate penalties
    slash_blocks = get_governable_parameter_value(client1, 'slashBlocksReward')
    log.info("Current block height: {}".format(client2.node.eth.blockNumber))
    # stop node
    node.stop()
    # Waiting 2 consensus block
    client2.economic.wait_consensus_blocknum(client2.node, 3)
    log.info("Current block height: {}".format(client2.node.eth.blockNumber))
    # view verifier list
    verifier_list = client2.ppos.getVerifierList()
    log.info("Current settlement cycle verifier list: {}".format(verifier_list))
    # Amount of penalty
    punishment_amonut = int(Decimal(str(block_reward)) * Decimal(str(slash_blocks)))
    log.info("punishment_amonut: {}".format(punishment_amonut))
    # view Restricting Plan
    restricting_info = client2.ppos.getRestrictingInfo(address2)
    log.info("restricting info: {}".format(restricting_info))
    info = restricting_info['Ret']
    assert (info['Pledge'] == pledge_amount - punishment_amonut * 2) or (info['Pledge'] == pledge_amount - punishment_amonut), 'ErrMsg: restricting Pledge amount {}'.format(
        info['Pledge'])
    assert (info['balance'] == pledge_amount - punishment_amonut * 2) or (info['balance'] == pledge_amount - punishment_amonut), 'ErrMsg: restricting balance amount {}'.format(
        info['balance'])
    # create Restricting Plan again
    staking_amount = von_amount(economic.create_staking_limit, 2)
    plan = [{'Epoch': 1, 'Amount': staking_amount}]
    result = client2.restricting.createRestrictingPlan(address2, plan, address1)
    assert_code(result, 0)
    # view Restricting Plan
    restricting_info3 = client2.ppos.getRestrictingInfo(address2)
    log.info("restricting info: {}".format(restricting_info3))
    assert_code(restricting_info3, 0)
    info2 = restricting_info3['Ret']
    assert info2['balance'] == staking_amount + info['balance'] - info[
        'debt'], "rrMsg: restricting balance amount {}".format(info2['balance'])
    assert info2['debt'] == 0, "rrMsg: restricting debt amount {}".format(info2['debt'])
    assert info2['plans'][0]['amount'] == staking_amount, "rrMsg: restricting plans amount {}".format(
        info2['plans'][0]['amount'])
    assert info2['Pledge'] == info['Pledge'], "rrMsg: restricting Pledge amount {}".format(info['Pledge'])
    # Waiting for the end of the 2 settlement
    client2.economic.wait_settlement_blocknum(client2.node, 2)
    # view Restricting Plan
    restricting_info3 = client2.ppos.getRestrictingInfo(address2)
    log.info("restricting info: {}".format(restricting_info3))
    assert_code(restricting_info3, 304005)


@pytest.mark.P1
def test_LS_RV_013(client_new_node):
    """
    同个账号锁仓给多个人
    :param client_new_node:
    :return:
    """
    client = client_new_node
    economic = client.economic
    node = client.node
    # create account
    address1, _ = economic.account.generate_account(node.web3, economic.create_staking_limit)
    address2, _ = economic.account.generate_account(node.web3, 0)
    address3, _ = economic.account.generate_account(node.web3, 0)
    # create Restricting Plan1
    plan = [{'Epoch': 1, 'Amount': economic.delegate_limit}]
    result = client.restricting.createRestrictingPlan(address2, plan, address1)
    assert_code(result, 0)
    restricting_info = client.ppos.getRestrictingInfo(address2)
    assert_code(restricting_info, 0)
    # create Restricting Plan1
    plan = [{'Epoch': 1, 'Amount': economic.delegate_limit}]
    result = client.restricting.createRestrictingPlan(address3, plan, address1)
    assert_code(result, 0)
    restricting_info = client.ppos.getRestrictingInfo(address3)
    assert_code(restricting_info, 0)


def create_a_multiplayer_lockout_plan(client):
    economic = client.economic
    node = client.node
    # create account
    first_address, _ = economic.account.generate_account(node.web3, von_amount(economic.create_staking_limit, 2))
    second_address, _ = economic.account.generate_account(node.web3, von_amount(economic.create_staking_limit, 2))
    locked_address, _ = economic.account.generate_account(node.web3, node.web3.toWei(1000, 'ether'))
    # create Restricting Plan1
    plan = [{'Epoch': 1, 'Amount': economic.create_staking_limit}]
    result = client.restricting.createRestrictingPlan(locked_address, plan, first_address)
    assert_code(result, 0)
    restricting_info = client.ppos.getRestrictingInfo(locked_address)
    assert_code(restricting_info, 0)
    # create Restricting Plan1
    plan = [{'Epoch': 1, 'Amount': economic.create_staking_limit}]
    result = client.restricting.createRestrictingPlan(locked_address, plan, second_address)
    assert_code(result, 0)
    restricting_info = client.ppos.getRestrictingInfo(locked_address)
    return locked_address, restricting_info


@pytest.mark.P1
def test_LS_RV_014(client_new_node):
    """
    同个账号被多个人锁仓
    :param client_new_node:
    :return:
    """
    client = client_new_node
    locked_address, restricting_info = create_a_multiplayer_lockout_plan(client)
    assert_code(restricting_info, 0)


@pytest.mark.P1
def test_LS_RV_015(client_new_node):
    """
    使用多人锁仓金额质押
    :param client_new_node:
    :return:
    """
    client = client_new_node
    locked_address, restricting_info = create_a_multiplayer_lockout_plan(client)
    # create staking
    result = client.staking.create_staking(1, locked_address, locked_address)
    assert_code(result, 0)


@pytest.mark.P1
def test_LS_RV_016(client_new_node):
    """
    使用多人锁仓金额委托
    :param client_new_node:
    :return:
    """
    client = client_new_node
    economic = client.economic
    node = client.node
    locked_address, restricting_info = create_a_multiplayer_lockout_plan(client)
    # create account
    pledge_address, _ = economic.account.generate_account(node.web3, von_amount(economic.create_staking_limit, 2))
    # create staking
    result = client.staking.create_staking(0, pledge_address, pledge_address)
    assert_code(result, 0)
    # Application for Commission
    result = client.delegate.delegate(1, locked_address, amount=economic.create_staking_limit)
    assert_code(result, 0)


@pytest.mark.P1
def test_LS_RV_017(client_new_node):
    """
    使用多人锁仓金额增持
    :param client_new_node:
    :return:
    """
    client = client_new_node
    locked_address, restricting_info = create_a_multiplayer_lockout_plan(client)
    # create staking
    result = client.staking.create_staking(1, locked_address, locked_address)
    assert_code(result, 0)
    # Apply for additional pledge
    result = client.staking.increase_staking(1, locked_address)
    assert_code(result, 0)


@pytest.mark.P2
def test_LS_RV_018(clients_new_node, reset_environment):
    """
    验证人非正常状态下创建锁仓计划（节点退出创建锁仓）
    :param clients_new_node:
    :return:
    """
    client1 = clients_new_node[0]
    log.info("Current linked client1: {}".format(client1.node.node_mark))
    client2 = clients_new_node[1]
    log.info("Current linked client2: {}".format(client2.node.node_mark))
    economic = client1.economic
    node = client1.node
    # create account
    address1, _ = economic.account.generate_account(node.web3, von_amount(economic.create_staking_limit, 2))
    # create staking
    result = client1.staking.create_staking(0, address1, address1)
    assert_code(result, 0)
    # Waiting settlement block
    client1.economic.wait_settlement_blocknum(client1.node)
    # stop node
    client1.node.stop()
    # Waiting 2 consensus block
    client2.economic.wait_consensus_blocknum(client2.node, 2)
    log.info("Current block height: {}".format(client2.node.eth.blockNumber))
    # create Restricting Plan1
    plan = [{'Epoch': 1, 'Amount': economic.delegate_limit}]
    result = client2.restricting.createRestrictingPlan(address1, plan, address1)
    assert_code(result, 0)


def create_account_restricting_plan(client, economic, node):
    # create account
    address1, _ = economic.account.generate_account(node.web3, von_amount(economic.create_staking_limit, 2))
    address2, _ = economic.account.generate_account(node.web3, node.web3.toWei(1000, 'ether'))
    # create Restricting Plan
    amount = economic.create_staking_limit
    plan = [{'Epoch': 1, 'Amount': amount}]
    result = client.restricting.createRestrictingPlan(address2, plan, address1)
    assert_code(result, 0)
    # view restricting info
    restricting_info = client.ppos.getRestrictingInfo(address2)
    log.info("restricting info: {}".format(restricting_info))
    assert_code(restricting_info, 0)
    info = restricting_info['Ret']
    assert info['balance'] == amount, 'ErrMsg: restricting balance amount {}'.format(info['balance'])
    assert info['Pledge'] == 0, 'ErrMsg: restricting Pledge amount {}'.format(info['Pledge'])
    return address2


@pytest.mark.P1
@pytest.mark.compatibility
def test_LS_PV_001(client_new_node):
    """
    锁仓账户质押正常节点
    :param client_new_node:
    :return:
    """
    client = client_new_node
    economic = client.economic
    node = client.node
    # create account restricting plan
    address2 = create_account_restricting_plan(client, economic, node)
    # create staking
    result = client.staking.create_staking(1, address2, address2)
    assert_code(result, 0)
    # view restricting info
    restricting_info = client.ppos.getRestrictingInfo(address2)
    log.info("restricting info: {}".format(restricting_info))
    assert_code(restricting_info, 0)
    info = restricting_info['Ret']
    assert info['Pledge'] == economic.create_staking_limit, 'ErrMsg: restricting Pledge amount {}'.format(
        info['Pledge'])


@pytest.mark.P1
def test_LS_PV_002(client_new_node):
    """
    创建计划质押-未找到锁仓信息
    :param client_new_node:
    :return:
    """
    client = client_new_node
    economic = client.economic
    node = client.node
    # create account
    address1, _ = economic.account.generate_account(node.web3, von_amount(economic.create_staking_limit, 2))
    address2, _ = economic.account.generate_account(node.web3, node.web3.toWei(1000, 'ether'))
    # create staking
    result = client.staking.create_staking(1, address2, address2)
    assert_code(result, 304005)


@pytest.mark.P1
@pytest.mark.compatibility
def test_LS_PV_003(client_new_node):
    """
    创建计划质押-锁仓计划质押金额<0
    :param client_new_node:
    :return:
    """
    client = client_new_node
    economic = client.economic
    node = client.node
    status = True
    # create account restricting plan
    address2 = create_account_restricting_plan(client, economic, node)
    try:
        # create staking
        client.staking.create_staking(1, address2, address2, amount=-1)
        status = False
    except Exception as e:
        log.info("Use case success, exception information：{} ".format(str(e)))
    assert status, "ErrMsg: create restricting result {}".format(status)


@pytest.mark.P1
def test_LS_PV_004(client_new_node):
    """
    创建计划质押-锁仓计划质押金额=0
    :param client_new_node:
    :return:
    """
    client = client_new_node
    economic = client.economic
    node = client.node
    # create account restricting plan
    address2 = create_account_restricting_plan(client, economic, node)
    # create staking
    result = client.staking.create_staking(1, address2, address2, amount=0)
    assert_code(result, 301100)


@pytest.mark.P1
def test_LS_PV_005(client_new_node):
    """
    创建计划质押-锁仓计划质押金额小于最低门槛
    :param client_new_node:
    :return:
    """
    client = client_new_node
    economic = client.economic
    node = client.node
    # create account restricting plan
    address2 = create_account_restricting_plan(client, economic, node)
    # create staking
    staking_amount = von_amount(economic.create_staking_limit, 0.8)
    result = client.staking.create_staking(1, address2, address2, amount=staking_amount)
    assert_code(result, 301100)


@pytest.mark.P2
def test_LS_PV_006(client_new_node):
    """
    创建计划质押-锁仓账户余额为0的情况下申请质押
    :param client_new_node:
    :return:
    """
    client = client_new_node
    economic = client.economic
    node = client.node
    status = True
    # create account
    address1, _ = economic.account.generate_account(node.web3, von_amount(economic.create_staking_limit, 2))
    address2, _ = economic.account.generate_account(node.web3, 0)
    # create Restricting Plan
    amount = economic.create_staking_limit
    plan = [{'Epoch': 1, 'Amount': amount}]
    result = client.restricting.createRestrictingPlan(address2, plan, address1)
    assert_code(result, 0)
    try:
        # create staking
        client.staking.create_staking(1, address2, address2)
        status = False
    except Exception as e:
        log.info("Use case success, exception information：{} ".format(str(e)))
    assert status, "ErrMsg: create restricting result {}".format(status)


@pytest.mark.P1
def test_LS_PV_007(clients_new_node):
    """
    创建计划退回质押-退回质押金额>锁仓质押金额
    :param clients_new_node:
    :return:
    """
    client1 = clients_new_node[0]
    log.info("Current linked client1: {}".format(client1.node.node_mark))
    client2 = clients_new_node[1]
    log.info("Current linked client2: {}".format(client2.node.node_mark))
    economic = client1.economic
    node = client1.node
    # create account
    amount1 = von_amount(economic.create_staking_limit, 2)
    amount2 = von_amount(economic.create_staking_limit, 2)
    address1, address2 = create_lock_release_amount(client1, amount1, amount2)
    # create Restricting Plan
    plan = [{'Epoch': 1, 'Amount': economic.create_staking_limit}]
    result = client1.restricting.createRestrictingPlan(address2, plan, address1)
    assert_code(result, 0)
    # create Restricting amount staking
    result = client1.staking.create_staking(1, address2, address2)
    assert_code(result, 0)
    time.sleep(3)
    # create Free amount staking
    result = client2.staking.create_staking(0, address2, address2)
    assert_code(result, 0)
    # withdrew staking
    result = client2.staking.withdrew_staking(address2)
    assert_code(result, 0)


@pytest.mark.P1
def test_LS_PV_008(client_new_node):
    """
    创建计划退回质押-欠释放金额=回退金额
    :param client_new_node:
    :return:
    """
    client = client_new_node
    economic = client.economic
    node = client.node
    # create restricting plan and staking
    address1, address2 = create_restricting_plan_and_staking(client, economic, node)
    # withdrew staking
    result = client.staking.withdrew_staking(address2)
    assert_code(result, 0)


@pytest.mark.P1
def test_LS_PV_009(client_new_node):
    """
    创建计划退回质押-欠释放金额<回退金额
    :param client_new_node:
    :return:
    """
    client1 = client_new_node
    log.info("Current linked client1: {}".format(client1.node.node_mark))
    economic = client1.economic
    node = client1.node
    # create account
    amount1 = von_amount(economic.create_staking_limit, 2)
    amount2 = von_amount(economic.create_staking_limit, 2)
    address1, address2 = create_lock_release_amount(client1, amount1, amount2)
    # create Restricting Plan
    plan = [{'Epoch': 1, 'Amount': economic.create_staking_limit}]
    result = client1.restricting.createRestrictingPlan(address2, plan, address1)
    assert_code(result, 0)
    # create Restricting amount staking
    result = client1.staking.create_staking(1, address2, address2)
    assert_code(result, 0)
    # wait settlement block
    economic.wait_settlement_blocknum(node)
    # view restricting info
    restricting_info = client1.ppos.getRestrictingInfo(address2)
    info = restricting_info['Ret']
    assert info['debt'] == economic.create_staking_limit, "rrMsg: restricting debt amount {}".format(info['debt'])
    # create Free amount staking
    result = client1.staking.increase_staking(0, address2)
    assert_code(result, 0)
    # withdrew staking
    result = client1.staking.withdrew_staking(address2)
    assert_code(result, 0)
    # view Restricting plan
    restricting_info = client1.ppos.getRestrictingInfo(address2)
    assert_code(restricting_info, 0)
    info = restricting_info['Ret']
    assert info['debt'] == economic.create_staking_limit, "errMsg: restricting debt amount {}".format(info['debt'])


@pytest.mark.P2
def test_LS_PV_010(client_new_node):
    """
    创建计划退回质押-锁仓账户余额不足的情况下申请退回质押
    :param client_new_node:
    :return:
    """
    client = client_new_node
    economic = client.economic
    node = client.node
    status = True
    # create account
    amount1 = von_amount(economic.create_staking_limit, 2)
    amount2 = node.web3.toWei(0.000009, 'ether')
    address1, address2 = create_lock_release_amount(client, amount1, amount2)
    # create Restricting Plan
    plan = [{'Epoch': 1, 'Amount': economic.create_staking_limit}]
    result = client.restricting.createRestrictingPlan(address2, plan, address1)
    assert_code(result, 0)
    # create Restricting amount staking
    result = client.staking.create_staking(1, address2, address2)
    assert_code(result, 0)
    log.info("address amount: {}".format(node.eth.getBalance(address2)))
    try:
        # withdrew staking
        client.staking.withdrew_staking(address2)
        status = False
    except Exception as e:
        log.info("Use case success, exception information：{} ".format(str(e)))
    assert status, "ErrMsg: create restricting result {}".format(status)


@pytest.mark.P2
def test_LS_PV_011(client_new_node):
    """
    锁仓账户退回质押金中，申请质押节点
    :param client_new_node:
    :return:
    """
    client = client_new_node
    economic = client.economic
    node = client.node
    # create restricting plan and staking
    address1, address2 = create_restricting_plan_and_staking(client, economic, node)
    # withdrew staking
    result = client.staking.withdrew_staking(address2)
    assert_code(result, 0)
    # create Restricting amount staking
    result = client.staking.create_staking(1, address2, address2)
    assert_code(result, 301101)


@pytest.mark.P2
def test_LS_PV_012(client_new_node):
    """
    锁仓账户申请完质押后又退回质押金（犹豫期）
    :param client_new_node:
    :return:
    """
    client = client_new_node
    economic = client.economic
    node = client.node
    # create account restricting plan
    address2 = create_account_restricting_plan(client, economic, node)
    # create staking
    result = client.staking.create_staking(1, address2, address2)
    assert_code(result, 0)
    # withdrew staking
    result = client.staking.withdrew_staking(address2)
    assert_code(result, 0)


@pytest.mark.P1
def test_LS_PV_013(client_new_node):
    """
    锁仓账户申请完质押后又退回质押金（锁定期）
    :param client_new_node:
    :return:
    """
    client = client_new_node
    economic = client.economic
    node = client.node
    # create account restricting plan
    address2 = create_account_restricting_plan(client, economic, node)
    # create staking
    result = client.staking.create_staking(1, address2, address2)
    assert_code(result, 0)
    # wait settlement block
    economic.wait_settlement_blocknum(node)
    # withdrew staking
    result = client.staking.withdrew_staking(address2)
    assert_code(result, 0)
    # wait settlement block
    economic.wait_settlement_blocknum(node, 2)
    # view restricting info
    restricting_info = client.ppos.getRestrictingInfo(address2)
    log.info("restricting info: {}".format(restricting_info))
    assert_code(restricting_info, 304005)


def create_free_pledge(client, economic):
    # create account
    amount1 = von_amount(economic.create_staking_limit, 2)
    amount2 = client.node.web3.toWei(1000, 'ether')
    address1, address2 = create_lock_release_amount(client, amount1, amount2)
    # create Restricting Plan
    plan = [{'Epoch': 1, 'Amount': von_amount(economic.delegate_limit, 10)}]
    result = client.restricting.createRestrictingPlan(address2, plan, address1)
    assert_code(result, 0)
    # view Restricting Plan informtion
    restricting_info = client.ppos.getRestrictingInfo(address2)
    log.info("Restricting Plan informtion: {}".format(restricting_info))
    # create staking
    result = client.staking.create_staking(0, address1, address1)
    assert_code(result, 0)
    return address2


@pytest.mark.P1
def test_LS_EV_001(client_new_node):
    """
    创建计划委托-委托正常节点
    :param client_new_node:
    :return:
    """
    client = client_new_node
    economic = client.economic
    address2 = create_free_pledge(client, economic)
    # Application for Commission
    result = client.delegate.delegate(1, address2)
    assert_code(result, 0)
    # view Restricting Plan
    restricting_info = client.ppos.getRestrictingInfo(address2)
    log.info("restricting info: {}".format(restricting_info))
    assert_code(restricting_info, 0)
    info = restricting_info['Ret']
    assert info['Pledge'] == economic.delegate_limit, 'ErrMsg: restricting Pledge amount {}'.format(
        info['Pledge'])


@pytest.mark.P1
def test_LS_EV_002(client_new_node):
    """
    创建计划委托-未找到锁仓信息
    :param client_new_node:
    :return:
    """
    client = client_new_node
    economic = client.economic
    node = client.node
    # create account
    amount1 = von_amount(economic.create_staking_limit, 2)
    amount2 = client.node.web3.toWei(1000, 'ether')
    address1, address2 = create_lock_release_amount(client, amount1, amount2)
    # create staking
    result = client.staking.create_staking(0, address1, address1)
    assert_code(result, 0)
    # Application for Commission
    result = client.delegate.delegate(1, address2)
    assert_code(result, 304005)


@pytest.mark.P1
def test_LS_EV_003(client_new_node, client_consensus):
    """
    锁仓账户委托基金会节点
    :param client_new_node:
    :param client_consensus:
    :return:
    """
    client = client_new_node
    economic = client.economic
    address2 = create_free_pledge(client, economic)
    # Application for Commission
    result = client_consensus.delegate.delegate(1, address2)
    assert_code(result, 301107)


@pytest.mark.P1
def test_LS_EV_004(client_new_node):
    """
    锁仓账户委托非候选节点(锁定期)
    :param client_new_node:
    :return:
    """
    client = client_new_node
    economic = client.economic
    node = client.node
    # create account
    pledge_acount = von_amount(economic.create_staking_limit, 2)
    lock_acount = node.web3.toWei(1000, 'ether')
    pledge_address, lock_address = create_lock_release_amount(client, pledge_acount, lock_acount)
    # create staking
    result = client.staking.create_staking(0, pledge_address, pledge_address)
    assert_code(result, 0)
    # Waiting for the end of the settlement
    economic.wait_settlement_blocknum(node)
    # create Restricting Plan
    plan = [{'Epoch': 1, 'Amount': von_amount(economic.delegate_limit, 10)}]
    result = client.restricting.createRestrictingPlan(lock_address, plan, pledge_address)
    assert_code(result, 0)
    # Application for Commission
    result = client.delegate.delegate(1, lock_address)
    assert_code(result, 0)
    # view Restricting Plan
    restricting_info = client.ppos.getRestrictingInfo(lock_address)
    log.info("restricting info: {}".format(restricting_info))
    assert_code(restricting_info, 0)
    info = restricting_info['Ret']
    assert info['Pledge'] == economic.delegate_limit, 'ErrMsg: restricting Pledge amount {}'.format(
        info['Pledge'])


@pytest.mark.P1
def test_LS_EV_005(client_new_node):
    """
    锁仓账户委托金额小于最低委托金
    :param client_new_node:
    :return:
    """
    client = client_new_node
    economic = client.economic
    address2 = create_free_pledge(client, economic)
    # Application for Commission
    delegate_amount = von_amount(economic.delegate_limit, 0.8)
    result = client.delegate.delegate(1, address2, amount=delegate_amount)
    assert_code(result, 301105)


@pytest.mark.P1
def test_LS_EV_006(client_new_node):
    """
    有锁仓可用金额，但是账户余额为0的情况下申请委托
    :param client_new_node:
    :return:
    """
    client = client_new_node
    economic = client.economic
    status = True
    # create account
    amount1 = von_amount(economic.create_staking_limit, 2)
    address1, address2 = create_lock_release_amount(client, amount1, 0)
    # create Restricting Plan
    plan = [{'Epoch': 1, 'Amount': economic.delegate_limit}]
    result = client.restricting.createRestrictingPlan(address2, plan, address1)
    assert_code(result, 0)
    # create staking
    result = client.staking.create_staking(0, address1, address1)
    assert_code(result, 0)
    try:
        # Application for Commission
        client.delegate.delegate(1, address2)
        status = False
    except Exception as e:
        log.info("Use case success, exception information：{} ".format(str(e)))
    assert status, "ErrMsg:Transfer result {}".format(status)


@pytest.mark.P1
def test_LS_EV_007(client_new_node):
    """
    创建计划委托-锁仓计划委托金额<0
    :param client_new_node:
    :return:
    """
    client = client_new_node
    economic = client.economic
    address2 = create_free_pledge(client, economic)
    status = True
    try:
        # Application for Commission
        # delegate_amount = von_amount(economic.delegate_limit, -8)
        result = client.delegate.delegate(1, address2, amount=-8)
        log.info("result: {}".format(result))
        status = False
    except Exception as e:
        log.info("Use case success, exception information：{} ".format(str(e)))
    assert status, "ErrMsg:Transfer result {}".format(status)


@pytest.mark.P1
def test_LS_EV_008(client_new_node):
    """
    创建计划委托-锁仓计划委托金额=0
    :param client_new_node:
    :return:
    """
    client = client_new_node
    economic = client.economic
    address2 = create_free_pledge(client, economic)
    # Application for Commission
    result = client.delegate.delegate(1, address2, amount=0)
    assert_code(result, 301105)


def create_delegation_information(client, economic, node, base):
    address2 = create_free_pledge(client, economic)
    # Application for Commission
    delegate_amount = von_amount(economic.delegate_limit, base)
    result = client.delegate.delegate(1, address2, amount=delegate_amount)
    assert_code(result, 0)
    # view restricting info
    restricting_info = client.ppos.getRestrictingInfo(address2)
    assert_code(restricting_info, 0)
    info = restricting_info['Ret']
    assert info['Pledge'] == delegate_amount, 'ErrMsg: restricting Pledge amount {}'.format(info['Pledge'])
    # get Pledge node information
    candidate_info = client.ppos.getCandidateInfo(node.node_id)
    info = candidate_info['Ret']
    staking_blocknum = info['StakingBlockNum']
    return address2, delegate_amount, staking_blocknum


@pytest.mark.P2
def test_LS_EV_009(client_new_node):
    """
    锁仓账户发起委托之后赎回部分委托验证（犹豫期）
    :param client_new_node:
    :return:
    """
    client = client_new_node
    economic = client.economic
    node = client.node
    # create delegation information
    address2, delegate_amount, staking_blocknum = create_delegation_information(client, economic, node, 10)
    # withdrew delegate
    redemption_amount = von_amount(economic.delegate_limit, 5)
    result = client.delegate.withdrew_delegate(staking_blocknum, address2, amount=redemption_amount)
    assert_code(result, 0)
    # view restricting info again
    restricting_info = client.ppos.getRestrictingInfo(address2)
    assert_code(restricting_info, 0)
    info = restricting_info['Ret']
    assert info['Pledge'] == delegate_amount - redemption_amount, 'ErrMsg: restricting Pledge amount {}'.format(
        info['Pledge'])


@pytest.mark.P2
def test_LS_EV_010(client_new_node):
    """
    锁仓账户发起委托之后赎回全部委托验证（犹豫期）
    :param client_new_node:
    :return:
    """
    client = client_new_node
    economic = client.economic
    node = client.node
    # create delegation information
    address2, delegate_amount, staking_blocknum = create_delegation_information(client, economic, node, 10)
    # withdrew delegate
    redemption_amount = von_amount(economic.delegate_limit, 10)
    client.delegate.withdrew_delegate(staking_blocknum, address2, amount=redemption_amount)
    # view restricting info again
    restricting_info = client.ppos.getRestrictingInfo(address2)
    assert_code(restricting_info, 0)
    info = restricting_info['Ret']
    assert info['Pledge'] == delegate_amount - redemption_amount, 'ErrMsg: restricting Pledge amount {}'.format(
        info['Pledge'])


@pytest.mark.P2
def test_LS_EV_011(client_new_node):
    """
    锁仓账户发起委托之后赎回委托验证（锁定期）
    :param client_new_node:
    :return:
    """
    client = client_new_node
    economic = client.economic
    node = client.node
    # create delegation information
    address2, delegate_amount, staking_blocknum = create_delegation_information(client, economic, node, 10)
    # Waiting for the end of the settlement cycle
    economic.wait_settlement_blocknum(node)
    # withdrew delegate
    redemption_amount = von_amount(economic.delegate_limit, 10)
    client.delegate.withdrew_delegate(staking_blocknum, address2, amount=redemption_amount)
    # view restricting info again
    restricting_info = client.ppos.getRestrictingInfo(address2)
    assert_code(restricting_info, 304005)


@pytest.mark.P1
def test_LS_EV_012(client_new_node):
    """
    锁仓赎回委托金额小于委托最低门槛
    :param client_new_node:
    :return:
    """
    client = client_new_node
    economic = client.economic
    node = client.node
    # create delegation information
    address2, delegate_amount, staking_blocknum = create_delegation_information(client, economic, node, 10)
    # withdrew delegate
    redemption_amount = von_amount(economic.delegate_limit, 0.8)
    result = client.delegate.withdrew_delegate(staking_blocknum, address2, amount=redemption_amount)
    assert_code(result, 301108)


@pytest.mark.P1
def test_LS_EV_013(client_new_node):
    """
    锁仓赎回委托后剩余委托小于委托最低门槛
    :param client_new_node:
    :return:
    """
    client = client_new_node
    economic = client.economic
    node = client.node
    # create delegation information
    address2, delegate_amount, staking_blocknum = create_delegation_information(client, economic, node, 2)
    # view restricting plan information
    restricting_info = client.ppos.getRestrictingInfo(address2)
    log.info("restricting plan information: {}".format(restricting_info))
    info = restricting_info['Ret']
    assert info['Pledge'] == von_amount(economic.delegate_limit, 2), 'ErrMsg: restricting Pledge amount {}'.format(
        info['Pledge'])
    # withdrew delegate
    redemption_amount = von_amount(economic.delegate_limit, 1.5)
    result = client.delegate.withdrew_delegate(staking_blocknum, address2, amount=redemption_amount)
    assert_code(result, 0)
    # view restricting plan information again
    restricting_info = client.ppos.getRestrictingInfo(address2)
    log.info("restricting plan information: {}".format(restricting_info))
    info = restricting_info['Ret']
    assert info['Pledge'] == 0, 'ErrMsg: restricting Pledge amount {}'.format(info['Pledge'])


@pytest.mark.P2
def test_LS_EV_014(clients_new_node, reset_environment):
    """
    锁仓账户委托节点状态异常验证人（节点已挂）
    :param clients_new_node:
    :param reset_environment:
    :return:
    """
    client1 = clients_new_node[0]
    log.info("Current linked client1: {}".format(client1.node.node_mark))
    client2 = clients_new_node[1]
    log.info("Current linked client2: {}".format(client2.node.node_mark))
    economic = client1.economic
    node = client1.node
    # create free pledge
    address2 = create_free_pledge(client1, economic)
    # stop pledge node
    node.stop()
    # Wait for the consensus round to end
    client2.economic.wait_consensus_blocknum(client2.node)
    # Application for Commission
    result = client2.delegate.delegate(1, address2, node_id=node.node_id)
    assert_code(result, 0)


@pytest.mark.P2
def test_LS_EV_015(client_new_node):
    """
    创建计划委托-锁仓账户余额为0的情况下申请委托
    :param client_new_node:
    :return:
    """
    client = client_new_node
    economic = client.economic
    status = True
    # create account
    amount1 = von_amount(economic.create_staking_limit, 2)
    address1, address2 = create_lock_release_amount(client, amount1, 0)
    # create Restricting Plan
    plan = [{'Epoch': 1, 'Amount': von_amount(economic.delegate_limit, 10)}]
    result = client.restricting.createRestrictingPlan(address2, plan, address1)
    assert_code(result, 0)
    # create staking
    result = client.staking.create_staking(0, address1, address1)
    assert_code(result, 0)
    try:
        # Application for Commission
        client.delegate.delegate(1, address2)
        status = False
    except Exception as e:
        log.info("Use case success, exception information：{} ".format(str(e)))
    assert status, "ErrMsg:Transfer result {}".format(status)


@pytest.mark.P1
def test_LS_EV_016(client_new_node):
    """
    创建计划退回委托-锁仓计划退回委托金额<0
    :param client_new_node:
    :return:
    """
    client = client_new_node
    economic = client.economic
    node = client.node
    status = True
    # create delegation information
    address2, delegate_amount, staking_blocknum = create_delegation_information(client, economic, node, 10)
    try:
        # withdrew delegate
        client.delegate.withdrew_delegate(staking_blocknum, address2, amount=-100)
        status = False
    except Exception as e:
        log.info("Use case success, exception information：{} ".format(str(e)))
    assert status, "ErrMsg:Transfer result {}".format(status)


@pytest.mark.P1
def test_LS_EV_017(client_new_node):
    """
    创建计划退回委托-锁仓计划退回委托金额=0
    :param client_new_node:
    :return:
    """
    client = client_new_node
    economic = client.economic
    node = client.node
    # create delegation information
    address2, delegate_amount, staking_blocknum = create_delegation_information(client, economic, node, 10)
    # withdrew delegate
    result = client.delegate.withdrew_delegate(staking_blocknum, address2, amount=0)
    assert_code(result, 301108)


@pytest.mark.P1
def test_LS_EV_018(client_new_node):
    """
    创建计划退回委托-锁仓计划退回委托金额>锁仓委托金额
    :param client_new_node:
    :return:
    """
    client = client_new_node
    economic = client.economic
    node = client.node
    # create delegation information
    address2, delegate_amount, staking_blocknum = create_delegation_information(client, economic, node, 10)
    # withdrew delegate
    redemption_amount = von_amount(economic.delegate_limit, 11)
    result = client.delegate.withdrew_delegate(staking_blocknum, address2, amount=redemption_amount)
    assert_code(result, 301113)


@pytest.mark.P1
def test_LS_EV_019(client_new_node):
    """
    创建计划退回委托-欠释放金额>赎回委托金额
    :param client_new_node:
    :return:
    """
    client = client_new_node
    economic = client.economic
    node = client.node
    address2, delegate_amount, staking_blocknum = create_delegation_information(client, economic, node, 10)
    # Waiting for the end of the settlement cycle
    economic.wait_settlement_blocknum(node)
    # view restricting plan informtion
    restricting_info = client.ppos.getRestrictingInfo(address2)
    log.info("restricting plan informtion: {}".format(restricting_info))
    info = restricting_info['Ret']
    assert info['debt'] == von_amount(economic.delegate_limit, 10), "rrMsg: restricting debt amount {}".format(
        info['debt'])
    # withdrew delegate
    redemption_amount = von_amount(economic.delegate_limit, 5)
    result = client.delegate.withdrew_delegate(staking_blocknum, address2, amount=redemption_amount)
    assert_code(result, 0)
    # view restricting plan informtion again
    restricting_info = client.ppos.getRestrictingInfo(address2)
    log.info("restricting plan informtion: {}".format(restricting_info))
    info = restricting_info['Ret']
    assert info['debt'] == von_amount(economic.delegate_limit,
                                      10) - redemption_amount, "rrMsg: restricting debt amount {}".format(info['debt'])


@pytest.mark.P1
def test_LS_EV_020(client_new_node):
    """
    创建计划退回委托-欠释放金额=撤销委托金额
    :param client_new_node:
    :return:
    """
    client = client_new_node
    economic = client.economic
    node = client.node
    address2, delegate_amount, staking_blocknum = create_delegation_information(client, economic, node, 10)
    # Waiting for the end of the settlement cycle
    economic.wait_settlement_blocknum(node)
    # view restricting plan informtion
    restricting_info = client.ppos.getRestrictingInfo(address2)
    log.info("restricting plan informtion: {}".format(restricting_info))
    info = restricting_info['Ret']
    assert info['debt'] == von_amount(economic.delegate_limit, 10), "rrMsg: restricting debt amount {}".format(
        info['debt'])
    # withdrew delegate
    redemption_amount = von_amount(economic.delegate_limit, 10)
    result = client.delegate.withdrew_delegate(staking_blocknum, address2, amount=redemption_amount)
    assert_code(result, 0)
    # view restricting plan informtion again
    restricting_info = client.ppos.getRestrictingInfo(address2)
    log.info("restricting plan informtion: {}".format(restricting_info))
    assert_code(restricting_info, 304005)


@pytest.mark.P1
def test_LS_EV_021(client_new_node):
    """
    创建计划退回委托-欠释放金额<撤销委托金额
    :param client_new_node:
    :return:
    """
    client = client_new_node
    economic = client.economic
    node = client.node
    address2, delegate_amount, staking_blocknum = create_delegation_information(client, economic, node, 5)
    # Waiting for the end of the settlement cycle
    economic.wait_settlement_blocknum(node)
    # view restricting plan informtion
    restricting_info = client.ppos.getRestrictingInfo(address2)
    log.info("restricting plan informtion: {}".format(restricting_info))
    info = restricting_info['Ret']
    assert info['debt'] == von_amount(economic.delegate_limit, 5), "rrMsg: restricting debt amount {}".format(
        info['debt'])
    # Application for Commission
    delegate_amount2 = von_amount(economic.delegate_limit, 5)
    result = client.delegate.delegate(0, address2, amount=delegate_amount2)
    assert_code(result, 0)
    # withdrew delegate
    redemption_amount = von_amount(economic.delegate_limit, 10)
    result = client.delegate.withdrew_delegate(staking_blocknum, address2, amount=redemption_amount)
    assert_code(result, 0)
    # view restricting plan informtion again
    restricting_info = client.ppos.getRestrictingInfo(address2)
    log.info("restricting plan informtion: {}".format(restricting_info))
    assert_code(restricting_info, 304005)


@pytest.mark.P1
def test_LS_EV_022(client_new_node):
    """
    创建计划退回委托-锁仓账户余额不足的情况下申请退回委托
    :param client_new_node:
    :return:
    """
    client = client_new_node
    economic = client.economic
    node = client.node
    status = True
    # create account
    amount1 = von_amount(economic.create_staking_limit, 2)
    amount2 = node.web3.toWei(0.000006, 'ether')
    address1, address2 = create_lock_release_amount(client, amount1, amount2)
    # create Restricting Plan
    plan = [{'Epoch': 1, 'Amount': von_amount(economic.delegate_limit, 10)}]
    result = client.restricting.createRestrictingPlan(address2, plan, address1)
    assert_code(result, 0)
    # create staking
    result = client.staking.create_staking(0, address1, address1)
    assert_code(result, 0)
    # Application for Commission
    result = client.delegate.delegate(1, address2)
    assert_code(result, 0)
    try:
        # get Pledge node information
        candidate_info = client.ppos.getCandidateInfo(node.node_id)
        info = candidate_info['Ret']
        staking_blocknum = info['StakingBlockNum']
        # withdrew delegate
        client.delegate.withdrew_delegate(staking_blocknum, address2)
        status = False
    except Exception as e:
        log.info("Use case success, exception information：{} ".format(str(e)))
    assert status, "ErrMsg:Transfer result {}".format(status)


def create_restricting_increase_staking(client, economic, node):
    # create account
    address1, _ = economic.account.generate_account(node.web3, economic.create_staking_limit)
    address2, _ = economic.account.generate_account(node.web3, von_amount(economic.create_staking_limit, 2))
    # create Restricting Plan1
    add_staking_amount = von_amount(economic.add_staking_limit, 10)
    plan = [{'Epoch': 1, 'Amount': add_staking_amount}]
    result = client.restricting.createRestrictingPlan(address2, plan, address1)
    assert_code(result, 0)
    restricting_info = client.ppos.getRestrictingInfo(address2)
    log.info("restricting plan informtion: {}".format(restricting_info))
    # create staking
    result = client.staking.create_staking(0, address2, address2)
    assert_code(result, 0)
    return address2


@pytest.mark.P1
def test_LS_IV_001(client_new_node):
    """
    锁仓账户申请质押后用锁仓余额进行增持质押
    :param client_new_node:
    :return:
    """
    client = client_new_node
    economic = client.economic
    node = client.node
    address2 = create_restricting_increase_staking(client, economic, node)
    # Create pledge of increasing holding
    result = client.staking.increase_staking(1, address2)
    assert_code(result, 0)


@pytest.mark.P1
def test_LS_IV_002(client_new_node):
    """
    有锁仓可用金额，但是账户gas不足的情况下申请增持
    :param client_new_node:
    :return:
    """
    client = client_new_node
    economic = client.economic
    node = client.node
    status = True
    # create account
    address1, _ = economic.account.generate_account(node.web3, von_amount(economic.create_staking_limit, 2))
    address2, _ = economic.account.generate_account(node.web3,
                                                    economic.create_staking_limit + node.web3.toWei(0.000009, 'ether'))
    # create Restricting Plan
    add_staking_amount = von_amount(economic.add_staking_limit, 10)
    plan = [{'Epoch': 1, 'Amount': add_staking_amount}]
    result = client.restricting.createRestrictingPlan(address2, plan, address1)
    assert_code(result, 0)
    restricting_info = client.ppos.getRestrictingInfo(address2)
    log.info("restricting plan informtion: {}".format(restricting_info))
    # create staking
    result = client.staking.create_staking(0, address2, address2)
    assert_code(result, 0)
    try:
        # Create pledge of increasing holding
        client.staking.increase_staking(1, address2)
        status = False
    except Exception as e:
        log.info("Use case success, exception information：{} ".format(str(e)))
    assert status, "ErrMsg:Transfer result {}".format(status)


@pytest.mark.P1
def test_LS_IV_003(clients_new_node, reset_environment):
    """
    锁仓账户增持状态异常验证人（节点已挂）
    :param clients_new_node:
    :param reset_environment:
    :return:
    """
    client1 = clients_new_node[0]
    log.info("Current linked client1: {}".format(client1.node.node_mark))
    client2 = clients_new_node[1]
    log.info("Current linked client2: {}".format(client2.node.node_mark))
    economic = client1.economic
    node = client1.node
    # Create restricting plan and free pledge
    address2 = create_restricting_increase_staking(client1, economic, node)
    # stop pledge node
    node.stop()
    # Wait for the consensus round to end
    client2.economic.wait_consensus_blocknum(client2.node)
    # Create pledge of increasing holding
    result = client2.staking.increase_staking(1, address2, node_id=node.node_id)
    assert_code(result, 0)


def restricting_plan_verification_pledge(client, economic, node):
    # create account
    address1, _ = economic.account.generate_account(node.web3, von_amount(economic.create_staking_limit, 2))
    # create Restricting Plan
    amount = economic.create_staking_limit
    plan = [{'Epoch': 1, 'Amount': amount}]
    result = client.restricting.createRestrictingPlan(address1, plan, address1)
    assert_code(result, 0)
    return address1


@pytest.mark.P2
def test_LS_CSV_001(client_new_node):
    """
    创建计划质押-锁仓账户和释放账户是同一个账户账户进行质押（质押金额小于锁仓金额）
    :param client_new_node:
    :return:
    """
    client = client_new_node
    economic = client.economic
    node = client.node
    # Create restricting plan
    address1 = restricting_plan_verification_pledge(client, economic, node)
    # create staking
    staking_amount = economic.create_staking_limit
    result = client.staking.create_staking(1, address1, address1, amount=staking_amount)
    assert_code(result, 0)


@pytest.mark.P2
def test_LS_CSV_002(client_new_node):
    """
    创建计划质押-锁仓账户和释放账户是同一个账户账户进行质押（质押金额大于锁仓金额）
    :param client_new_node:
    :return:
    """
    client = client_new_node
    economic = client.economic
    node = client.node
    # Create restricting plan
    address1 = restricting_plan_verification_pledge(client, economic, node)
    # create staking
    staking_amount = von_amount(economic.create_staking_limit, 2)
    result = client.staking.create_staking(0, address1, address1, amount=staking_amount)
    assert_code(result, 301111)


def restricting_plan_verification_pledge2(client, economic, node):
    # create account
    address1, _ = economic.account.generate_account(node.web3, von_amount(economic.create_staking_limit, 2))
    address2, _ = economic.account.generate_account(node.web3, node.web3.toWei(1000, 'ether'))
    # create Restricting Plan
    amount = economic.create_staking_limit
    plan = [{'Epoch': 1, 'Amount': amount}]
    result = client.restricting.createRestrictingPlan(address2, plan, address1)
    assert_code(result, 0)
    return address2


@pytest.mark.P2
def test_LS_CSV_003(client_new_node):
    """
    创建计划质押-锁仓账户和释放账户不同时进行质押（质押金额小于等于锁仓金额）
    :param client_new_node:
    :return:
    """
    client = client_new_node
    economic = client.economic
    node = client.node
    # Create restricting plan
    address2 = restricting_plan_verification_pledge2(client, economic, node)
    # create staking
    staking_amount = economic.create_staking_limit
    result = client.staking.create_staking(1, address2, address2, amount=staking_amount)
    assert_code(result, 0)


@pytest.mark.P2
def test_LS_CSV_004(client_new_node):
    """
    创建计划质押-锁仓账户和释放账户不同时进行质押（质押金额大于锁仓金额）
    :param client_new_node:
    :return:
    """
    client = client_new_node
    economic = client.economic
    node = client.node
    # Create restricting plan
    address2 = restricting_plan_verification_pledge2(client, economic, node)
    # create staking
    staking_amount = von_amount(economic.create_staking_limit, 2)
    result = client.staking.create_staking(1, address2, address2, amount=staking_amount)
    assert_code(result, 304013)


def restricting_plan_verification_add_staking(client, economic, node):
    # create account
    address1, _ = economic.account.generate_account(node.web3, von_amount(economic.create_staking_limit, 3))
    # create Restricting Plan
    amount = von_amount(economic.create_staking_limit, 2)
    plan = [{'Epoch': 1, 'Amount': amount}]
    result = client.restricting.createRestrictingPlan(address1, plan, address1)
    assert_code(result, 0)
    # create staking
    result = client.staking.create_staking(1, address1, address1)
    assert_code(result, 0)
    return address1


@pytest.mark.P2
def test_LS_CSV_005(client_new_node):
    """
    锁仓账户和释放账户是同一个账户账户进行增持质押（质押金额小于锁仓金额）
    :param client_new_node:
    :return:
    """
    client = client_new_node
    economic = client.economic
    node = client.node
    # create restricting plan staking
    address1 = restricting_plan_verification_add_staking(client, economic, node)
    # Additional pledge
    increase_amount = von_amount(economic.delegate_limit, 5)
    result = client.staking.increase_staking(1, address1, amount=increase_amount)
    assert_code(result, 0)


@pytest.mark.P2
def test_LS_CSV_006(client_new_node):
    """
    锁仓账户和释放账户是同一个账户账户进行增持质押（质押金额大于锁仓金额）
    :param client_new_node:
    :return:
    """
    client = client_new_node
    economic = client.economic
    node = client.node
    # create restricting plan staking
    address1 = restricting_plan_verification_add_staking(client, economic, node)
    # Additional pledge
    increase_amount = von_amount(economic.create_staking_limit, 2)
    result = client.staking.increase_staking(1, address1, amount=increase_amount)
    assert_code(result, 304013)


def restricting_plan_verification_add_staking2(client, economic, node):
    # create account
    address1, _ = economic.account.generate_account(node.web3, von_amount(economic.create_staking_limit, 2))
    address2, _ = economic.account.generate_account(node.web3, von_amount(economic.create_staking_limit, 2))
    # create Restricting Plan
    amount = von_amount(economic.add_staking_limit, 10)
    plan = [{'Epoch': 1, 'Amount': amount}]
    result = client.restricting.createRestrictingPlan(address2, plan, address1)
    assert_code(result, 0)
    # create staking
    result = client.staking.create_staking(0, address2, address2)
    assert_code(result, 0)
    return address2


@pytest.mark.P2
def test_LS_CSV_007(client_new_node):
    """
    锁仓账户和释放账户不同时进行质押（增持质押金额小于锁仓金额）
    :param client_new_node:
    :return:
    """
    client = client_new_node
    economic = client.economic
    node = client.node
    # create restricting plan staking
    address2 = restricting_plan_verification_add_staking2(client, economic, node)
    # Additional pledge
    increase_amount = von_amount(economic.add_staking_limit, 5)
    result = client.staking.increase_staking(1, address2, amount=increase_amount)
    assert_code(result, 0)


@pytest.mark.P2
def test_LS_CSV_008(client_new_node):
    """
    锁仓账户和释放账户不同时进行质押（增持质押金额大于锁仓金额）
    :param client_new_node:
    :return:
    """
    client = client_new_node
    economic = client.economic
    node = client.node
    # create restricting plan staking
    address2 = restricting_plan_verification_add_staking2(client, economic, node)
    # Additional pledge
    increase_amount = von_amount(economic.add_staking_limit, 15)
    result = client.staking.increase_staking(1, address2, amount=increase_amount)
    assert_code(result, 304013)


def restricting_plan_verification_delegate(client, economic, node):
    # create account
    address1, _ = economic.account.generate_account(node.web3, von_amount(economic.create_staking_limit, 2))
    address2, _ = economic.account.generate_account(node.web3, node.web3.toWei(1000, 'ether'))
    # create Restricting Plan
    amount = von_amount(economic.delegate_limit, 10)
    plan = [{'Epoch': 1, 'Amount': amount}]
    result = client.restricting.createRestrictingPlan(address2, plan, address1)
    assert_code(result, 0)
    # create staking
    result = client.staking.create_staking(0, address1, address1)
    assert_code(result, 0)
    return address2


@pytest.mark.P2
def test_LS_CSV_009(client_new_node):
    """
    锁仓账户和释放账户不是同一个账号进行委托（委托金额小于锁仓金额）
    :param client_new_node:
    :return:
    """
    client = client_new_node
    economic = client.economic
    node = client.node
    # create restricting plan staking
    address2 = restricting_plan_verification_delegate(client, economic, node)
    # Additional pledge
    delegate_amount = von_amount(economic.delegate_limit, 5)
    result = client.delegate.delegate(1, address2, amount=delegate_amount)
    assert_code(result, 0)


@pytest.mark.P2
def test_LS_CSV_010(client_new_node):
    """
    锁仓账户和释放账户不是同一个账号进行委托（委托金额大于锁仓金额）
    :param client_new_node:
    :return:
    """
    client = client_new_node
    economic = client.economic
    node = client.node
    # create restricting plan staking
    address2 = restricting_plan_verification_delegate(client, economic, node)
    # Additional pledge
    delegate_amount = von_amount(economic.delegate_limit, 15)
    result = client.delegate.delegate(1, address2, amount=delegate_amount)
    assert_code(result, 304013)


@pytest.mark.P2
def test_LS_CSV_011(client_new_node):
    """
    锁仓账号在犹豫期申请质押后，在锁定期申请增持后，在申请退回质押金
    :param client_new_node:
    :return:
    """
    client = client_new_node
    economic = client.economic
    node = client.node
    # create account
    address1, _ = economic.account.generate_account(node.web3, von_amount(economic.create_staking_limit, 3))
    # create Restricting Plan
    amount = von_amount(economic.create_staking_limit, 2)
    plan = [{'Epoch': 2, 'Amount': amount}]
    result = client.restricting.createRestrictingPlan(address1, plan, address1)
    assert_code(result, 0)
    # create staking
    result = client.staking.create_staking(1, address1, address1)
    assert_code(result, 0)
    # Waiting for the end of the settlement period
    economic.wait_settlement_blocknum(node)
    # Additional pledge
    result = client.staking.increase_staking(1, address1)
    assert_code(result, 0)
    balance = node.eth.getBalance(address1)
    log.info("Account address: {} balance: {}".format(address1, balance))
    # Application for return of pledge
    result = client.staking.withdrew_staking(address1)
    assert_code(result, 0)
    balance1 = node.eth.getBalance(address1)
    log.info("Account address: {} balance: {}".format(address1, balance1))
    # Waiting for the end of the 2 settlement period
    economic.wait_settlement_blocknum(node, 2)
    balance2 = node.eth.getBalance(address1)
    log.info("Account address: {} balance: {}".format(address1, balance2))
    assert balance2 - balance1 > economic.create_staking_limit, "errMsg: Account address: {} balance: {}".format(
        address1, balance2)


@pytest.mark.P2
def test_LS_CSV_012(client_new_node):
    """
    锁仓账户退回质押金中，申请委托节点
    :param client_new_node:
    :return:
    """
    client = client_new_node
    economic = client.economic
    node = client.node
    # create restricting plan staking
    address1 = restricting_plan_verification_pledge(client, economic, node)
    # create staking
    result = client.staking.create_staking(1, address1, address1)
    assert_code(result, 0)
    # create account
    address2, _ = economic.account.generate_account(node.web3, node.web3.toWei(1000, 'ether'))
    # Waiting for the end of the settlement period
    economic.wait_settlement_blocknum(node)
    # Application for return of pledge
    result = client.staking.withdrew_staking(address1)
    assert_code(result, 0)
    # create Restricting Plan
    amount = von_amount(economic.delegate_limit, 10)
    plan = [{'Epoch': 1, 'Amount': amount}]
    result = client.restricting.createRestrictingPlan(address2, plan, address2)
    assert_code(result, 0)
    # Free amount Entrust node
    result = client.delegate.delegate(0, address2)
    assert_code(result, 301103)
    # Restricting amount Entrust node
    result = client.delegate.delegate(1, address2)
    assert_code(result, 301103)


@pytest.mark.P2
def test_LS_CSV_013(client_new_node):
    """
    锁仓账户退回质押金中，申请增持质押
    :param client_new_node:
    :return:
    """
    client = client_new_node
    economic = client.economic
    node = client.node
    # create restricting plan staking
    address1 = restricting_plan_verification_add_staking(client, economic, node)
    # Waiting for the end of the settlement period
    economic.wait_settlement_blocknum(node)
    # Application for return of pledge
    result = client.staking.withdrew_staking(address1)
    assert_code(result, 0)
    # create Restricting Plan
    amount = von_amount(economic.add_staking_limit, 5)
    plan = [{'Epoch': 1, 'Amount': amount}]
    result = client.restricting.createRestrictingPlan(address1, plan, address1)
    assert_code(result, 0)
    # Restricting amount Additional pledge
    result = client.staking.increase_staking(1, address1)
    assert_code(result, 301103)
    # Free amount Additional pledge
    result = client.staking.increase_staking(0, address1)
    assert_code(result, 301103)


def steps_of_returning_pledge(client, economic, node):
    # create restricting plan staking
    address1 = restricting_plan_verification_pledge(client, economic, node)
    # create staking
    result = client.staking.create_staking(1, address1, address1)
    assert_code(result, 0)
    # create account
    address2, _ = economic.account.generate_account(node.web3, node.web3.toWei(1000, 'ether'))
    # Waiting for the end of the settlement period
    economic.wait_settlement_blocknum(node)
    # Application for return of pledge
    result = client.staking.withdrew_staking(address1)
    assert_code(result, 0)
    # Waiting for the end of the 2 settlement period
    economic.wait_settlement_blocknum(node, 2)
    log.info("Pledge node information: {}".format(client.ppos.getCandidateInfo(node.node_id)))
    return address1, address2


@pytest.mark.P2
def test_LS_CSV_014(client_new_node):
    """
    锁仓账户退回质押金后，重新申请质押节点
    :param client_new_node:
    :return:
    """
    client = client_new_node
    economic = client.economic
    node = client.node
    # After returning the deposit
    address1, address2 = steps_of_returning_pledge(client, economic, node)
    # create Restricting Plan
    amount = economic.create_staking_limit
    plan = [{'Epoch': 1, 'Amount': amount}]
    result = client.restricting.createRestrictingPlan(address1, plan, address1)
    assert_code(result, 0)
    restricting_info = client.ppos.getRestrictingInfo(address1)
    log.info("restricting plan informtion: {}".format(restricting_info))
    # create staking
    result = client.staking.create_staking(1, address1, address1)
    assert_code(result, 0)


@pytest.mark.P2
def test_LS_CSV_015(client_new_node):
    """
    锁仓账户退回质押金后，重新申请委托节点
    :param client_new_node:
    :return:
    """
    client = client_new_node
    economic = client.economic
    node = client.node
    # After returning the deposit
    address1, address2 = steps_of_returning_pledge(client, economic, node)
    # create Restricting Plan
    amount = von_amount(economic.delegate_limit, 10)
    plan = [{'Epoch': 1, 'Amount': amount}]
    result = client.restricting.createRestrictingPlan(address2, plan, address2)
    assert_code(result, 0)
    # Free amount Entrust node
    result = client.delegate.delegate(0, address2)
    assert_code(result, 301102)
    # Restricting amount Entrust node
    result = client.delegate.delegate(1, address2)
    assert_code(result, 301102)


@pytest.mark.P2
def test_LS_CSV_016(client_new_node):
    """
    锁仓账户退回质押金后，重新申请增持质押
    :param client_new_node:
    :return:
    """
    client = client_new_node
    economic = client.economic
    node = client.node
    # After returning the deposit
    address1, address2 = steps_of_returning_pledge(client, economic, node)
    # create Restricting Plan
    amount = von_amount(economic.add_staking_limit, 5)
    plan = [{'Epoch': 1, 'Amount': amount}]
    result = client.restricting.createRestrictingPlan(address1, plan, address1)
    assert_code(result, 0)
    # Restricting amount Additional pledge
    result = client.staking.increase_staking(1, address1)
    assert_code(result, 301102)
    # Free amount Additional pledge
    result = client.staking.increase_staking(0, address1)
    assert_code(result, 301102)
