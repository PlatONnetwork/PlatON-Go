import time
import pytest
import allure
from dacite import from_dict
from common.key import get_pub_key, mock_duplicate_sign
from common.log import log
from client_sdk_python import Web3
from decimal import Decimal
from tests.lib import EconomicConfig, Genesis, StakingConfig, Staking, check_node_in_list, assert_code, von_amount, \
    get_governable_parameter_value


def restricting_plan_validation_release(client, economic, node):
    # create account
    address1, _ = economic.account.generate_account(node.web3, von_amount(economic.create_staking_limit, 2))
    # create Restricting Plan
    amount = economic.create_staking_limit
    plan = [{'Epoch': 1, 'Amount': amount}]
    result = client.restricting.createRestrictingPlan(address1, plan, address1)
    assert_code(result, 0)
    return address1


@pytest.mark.P2
def test_UP_FV_001(client_new_node_obj):
    """
    只有一个锁仓期，到达释放期返回解锁金额
    :param client_new_node_obj:
    :return:
    """
    client = client_new_node_obj
    economic = client.economic
    node = client.node
    # create restricting plan
    address1 = restricting_plan_validation_release(client, economic, node)
    # view Account balance
    balance = node.eth.getBalance(address1)
    log.info("Account balance: {}".format(balance))
    # Waiting for the end of the settlement period
    economic.wait_settlement_blocknum(node)
    # view Account balance again
    balance1 = node.eth.getBalance(address1)
    log.info("Account balance: {}".format(balance1))
    assert balance1 == balance + economic.create_staking_limit, "ErrMsg:Account balance: {}".format(balance1)


@pytest.mark.P1
def test_UP_FV_002(client_new_node_obj):
    """
    只有一个锁仓期，未达释放期返回解锁金额
    :param client_new_node_obj:
    :return:
    """
    client = client_new_node_obj
    economic = client.economic
    node = client.node
    # create restricting plan
    address1 = restricting_plan_validation_release(client, economic, node)
    # view restricting plan index 0 amount
    restricting_info = client.ppos.getRestrictingInfo(address1)
    log.info("restricting plan information: {}".format(restricting_info))
    amount = restricting_info['Ret']['plans'][0]['amount']
    # view Account balance
    balance = node.eth.getBalance(address1)
    log.info("Account balance: {}".format(balance))
    # Waiting for the end of the settlement period
    economic.wait_consensus_blocknum(node)
    # view restricting plan index 0 amount again
    restricting_info = client.ppos.getRestrictingInfo(address1)
    amount1 = restricting_info['Ret']['plans'][0]['amount']
    # view Account balance again
    balance1 = node.eth.getBalance(address1)
    log.info("Account balance: {}".format(balance1))
    assert amount1 == amount, "ErrMsg:restricting index 0 amount: {}".format(amount1)
    assert balance1 == balance, "ErrMsg:Account balance: {}".format(balance1)


@pytest.mark.P1
def test_UP_FV_003(client_new_node_obj):
    """
    多个锁仓期，依次部分释放期返回解锁金额
    :param client_new_node_obj:
    :return:
    """
    client = client_new_node_obj
    economic = client.economic
    node = client.node
    # create account
    address1, _ = economic.account.generate_account(node.web3, von_amount(economic.create_staking_limit, 2))
    # create Restricting Plan
    amount = von_amount(economic.delegate_limit, 10)
    plan = [{'Epoch': 1, 'Amount': amount}, {'Epoch': 2, 'Amount': amount}]
    result = client.restricting.createRestrictingPlan(address1, plan, address1)
    assert_code(result, 0)
    # view Restricting Plan again
    restricting_info = client.ppos.getRestrictingInfo(address1)
    log.info("restricting plan information: {}".format(restricting_info))
    assert len(restricting_info['Ret']['plans']) == 2, "ErrMsg:Planned releases: {}".format(
        len(restricting_info['Ret']['plans']))
    # Waiting for the end of the settlement period
    economic.wait_settlement_blocknum(node)
    # view Restricting Plan
    restricting_info = client.ppos.getRestrictingInfo(address1)
    log.info("restricting plan information: {}".format(restricting_info))
    assert len(restricting_info['Ret']['plans']) == 1, "ErrMsg:Planned releases: {}".format(
        len(restricting_info['Ret']['plans']))


@pytest.mark.P1
def test_UP_FV_004(client_new_node_obj):
    """
    锁仓账户申请质押到释放期后释放锁定金额不足
    :param client_new_node_obj:
    :return:
    """
    client = client_new_node_obj
    economic = client.economic
    node = client.node
    # create restricting plan
    address1 = restricting_plan_validation_release(client, economic, node)
    # create staking
    staking_amount = economic.create_staking_limit
    result = client.staking.create_staking(1, address1, address1, amount=staking_amount)
    assert_code(result, 0)
    # Waiting for the end of the settlement period
    economic.wait_settlement_blocknum(node)
    # view restricting plan
    restricting_info = client.ppos.getRestrictingInfo(address1)
    log.info("restricting plan informtion: {}".format(restricting_info))
    info = restricting_info['Ret']
    assert info['debt'] == economic.create_staking_limit, 'ErrMsg: restricting debt amount {}'.format(
        info['debt'])

