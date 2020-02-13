# -*- coding: utf-8 -*-
from tests.lib.utils import *
import pytest
import allure


@allure.title("Normal overweight")
@pytest.mark.P0
def test_AS_001_002_009(client_new_node):
    """
    001:Normal overweight
    002:The verifier initiates the overweight with the amount of free account, meeting the minimum threshold
    009:Hesitation period add pledge, inquire pledge information
    """
    StakeThreshold = get_governable_parameter_value(client_new_node, "StakeThreshold")
    log.info(StakeThreshold)
    address, _ = client_new_node.economic.account.generate_account(client_new_node.node.web3,
                                                                   10 ** 18 * 10000000)
    result = client_new_node.staking.create_staking(0, address, address)
    assert_code(result, 0)
    result = client_new_node.staking.increase_staking(0, address)
    assert_code(result, 0)
    result = client_new_node.ppos.getCandidateInfo(client_new_node.node.node_id)
    staking_amount = client_new_node.economic.create_staking_limit
    add_staking_amount = client_new_node.economic.add_staking_limit
    assert result["Ret"]["Shares"] == staking_amount + add_staking_amount


@allure.title("The verifier is not on the verifier and candidate list")
@pytest.mark.P2
def test_AS_003(client_new_node):
    """
    :param client_new_node_obj:
    :param get_generate_account:
    :return:
    """
    address, _ = client_new_node.economic.account.generate_account(client_new_node.node.web3,
                                                                   10 ** 18 * 10000000)
    result = client_new_node.staking.increase_staking(0, address)
    log.info(result)
    assert_code(result, 301102)


@allure.title("Undersupply of gas")
@pytest.mark.P3
def test_AS_004(client_new_node):
    """
    :param client_new_node_obj:
    :param get_generate_account:
    :return:
    """
    address, _ = client_new_node.economic.account.generate_account(client_new_node.node.web3,
                                                                   10 ** 18 * 10000000)
    client_new_node.staking.create_staking(0, address, address)
    fig = {"gas": 1}
    status = 0
    try:
        result = client_new_node.staking.increase_staking(0, address, transaction_cfg=fig)
        log.info(result)
    except BaseException:
        status = 1
    assert status == 1


@allure.title("Insufficient balance initiated overweight")
@pytest.mark.P3
def test_AS_005(client_new_node):
    """
    :param client_new_node_obj:
    :return:
    """
    account = client_new_node.economic.account
    node = client_new_node.node
    address, _ = account.generate_account(node.web3, 10)
    status = 0
    try:
        result = client_new_node.staking.increase_staking(0, address)
        log.info(result)
    except BaseException:
        status = 1
    assert status == 1


@allure.title("(hesitation period) holdings less than the minimum threshold")
@pytest.mark.P1
def test_AS_007(client_new_node):
    """
    :param client_new_node_obj:
    :return:
    """
    address, _ = client_new_node.economic.account.generate_account(client_new_node.node.web3,
                                                                   10 ** 18 * 10000000)
    client_new_node.staking.create_staking(0, address, address)
    add_staking_amount = client_new_node.economic.add_staking_limit
    result = client_new_node.staking.increase_staking(0, address, amount=add_staking_amount - 1)
    log.info(result)
    assert_code(result, 301104)


@allure.title("(hesitation period) when the verifier revoks the pledge, he/she shall apply for adding the pledge")
@pytest.mark.P1
def test_AS_008(client_new_node):
    """
    :param client_new_node_obj:
    :return:
    """
    address, _ = client_new_node.economic.account.generate_account(client_new_node.node.web3,
                                                                   10 ** 18 * 10000000)
    client_new_node.staking.create_staking(0, address, address)
    log.info("进入下个周期")
    client_new_node.economic.wait_settlement_blocknum(client_new_node.node)
    client_new_node.staking.withdrew_staking(address)
    result = client_new_node.staking.increase_staking(0, address)
    log.info(result)
    assert_code(result, 301103)


@allure.title("(lockup period) normal increase")
@pytest.mark.P0
@pytest.mark.compatibility
def test_AS_011_012_013_014(client_new_node):
    """
    011:(lockup period) normal increase
    012:(lockup period) overweight meets the minimum threshold
    013:(lock-up period) gas underinitiation overweight
    014:(lockup period) insufficient balance to initiate overweight pledge
    :param client_new_node_obj:
    :return:
    """
    address, _ = client_new_node.economic.account.generate_account(client_new_node.node.web3,
                                                                   10 ** 18 * 10000000)
    client_new_node.staking.create_staking(0, address, address)
    log.info("进入下个周期")
    client_new_node.economic.wait_settlement_blocknum(client_new_node.node)
    result = client_new_node.staking.increase_staking(0, address)
    log.info(result)
    assert_code(result, 0)
    result = client_new_node.ppos.getCandidateInfo(client_new_node.node.node_id)
    log.info(result)
    staking_amount = client_new_node.economic.create_staking_limit
    add_staking_amount = client_new_node.economic.add_staking_limit
    assert result["Ret"]["Shares"] == staking_amount + add_staking_amount
    assert result["Ret"]["Released"] == staking_amount
    assert result["Ret"]["ReleasedHes"] == add_staking_amount
    fig = {"gas": 1}
    status = 0
    try:
        result = client_new_node.staking.increase_staking(0, address, transaction_cfg=fig)
        log.info(result)
    except BaseException:
        status = 1
    assert status == 1
    account = client_new_node.economic.account
    node = client_new_node.node
    address, _ = account.generate_account(node.web3, 10)
    status = 0
    try:
        result = client_new_node.staking.increase_staking(0, address, address)
        log.info(result)
    except BaseException:
        status = 1
    assert status == 1


@allure.title("The free amount is insufficient, the lock position is sufficient, and the free amount is added")
@pytest.mark.P1
def test_AS_015(client_new_node):
    """
    :param client_new_node_obj:
    :param get_generate_account:
    :return:
    """
    address, _ = client_new_node.economic.account.generate_account(client_new_node.node.web3,
                                                                   10 ** 18 * 10000000)
    result = client_new_node.staking.create_staking(0, address, address)
    log.info(result)
    node = client_new_node.node
    amount = node.eth.getBalance(address)
    log.info("Wallet balance{}".format(amount))
    locked_amount = amount - node.web3.toWei(1, "ether")

    plan = [{'Epoch': 1, 'Amount': locked_amount}]
    log.info("The balance of the wallet is used as a lock")
    result = client_new_node.restricting.createRestrictingPlan(address, plan, address)
    log.info(result)
    amount = node.eth.getBalance(address)
    log.info("Check your wallet balance{}".format(amount))
    result = client_new_node.staking.increase_staking(0, address)
    assert_code(result, 301111)


@allure.title("The free amount is insufficient, the lock position is sufficient, and the free amount is added")
@pytest.mark.P1
def test_AS_016(client_new_node):
    """
    :param client_new_node_obj:
    :return:
    """
    address, _ = client_new_node.economic.account.generate_account(client_new_node.node.web3,
                                                                   10 ** 18 * 10000000)
    result = client_new_node.staking.create_staking(0, address, address)
    log.info(result)
    node = client_new_node.node
    amount = node.eth.getBalance(address)
    log.info("Wallet balance{}".format(amount))
    locked_amount = 100000000000000000000
    plan = [{'Epoch': 1, 'Amount': locked_amount}]
    result = client_new_node.restricting.createRestrictingPlan(address, plan, address)
    log.info(result)
    assert_code(result, 0)
    value = 101000000000000000000
    result = client_new_node.staking.increase_staking(1, address, amount=value)
    log.info(result)
    assert_code(result, 304013)


@allure.title("The amount of the increase is less than the threshold")
@pytest.mark.P1
def test_AS_017(client_new_node):
    """
    :param client_new_node_obj:
    :return:
    """
    address, _ = client_new_node.economic.account.generate_account(client_new_node.node.web3,
                                                                   10 ** 18 * 10000000)
    result = client_new_node.staking.create_staking(0, address, address)
    log.info(result)
    add_staking_amount = client_new_node.economic.add_staking_limit
    result = client_new_node.staking.increase_staking(0, address, amount=add_staking_amount - 1)
    log.info(result)
    assert_code(result, 301104)


@allure.title("Increase the number of active withdrawal but still in the freeze period of the candidate")
@pytest.mark.P0
def test_AS_018_019(client_new_node):
    """
    018:Increase the number of active withdrawal but still in the freeze period of the candidate
    019:Candidates whose holdings have been actively withdrawn and who have passed the freeze period
    :param client_new_node_obj:
    :return:
    """
    address, _ = client_new_node.economic.account.generate_account(client_new_node.node.web3,
                                                                   10 ** 18 * 10000000)
    result = client_new_node.staking.create_staking(0, address, address)
    assert_code(result, 0)
    log.info("Next settlement period")
    client_new_node.economic.wait_settlement_blocknum(client_new_node.node)
    result = client_new_node.staking.withdrew_staking(address)
    assert_code(result, 0)
    result = client_new_node.staking.increase_staking(0, address)
    log.info(result)
    assert_code(result, 301103)
    log.info("Next settlement period")
    client_new_node.economic.wait_settlement_blocknum(client_new_node.node, number=2)
    result = client_new_node.staking.increase_staking(0, address)
    log.info(result)
    assert_code(result, 301102)


@allure.title("Add to the list of candidates who have been penalized and are still in the freeze period")
@pytest.mark.P0
def test_AS_020_021(clients_new_node, client_consensus):
    """
    020:Add to the list of candidates who have been penalized and are still in the freeze period
    021:A candidate whose holdings have been penalized has passed the freeze period
    :param client_new_node_obj:
    :return:
    """
    client = clients_new_node[0]
    node = client.node
    other_node = client_consensus.node
    economic = client.economic
    address, pri_key = economic.account.generate_account(node.web3, 10 ** 18 * 10000000)

    value = economic.create_staking_limit * 2
    result = client.staking.create_staking(0, address, address, amount=value)
    assert_code(result, 0)
    economic.wait_consensus_blocknum(other_node, number=4)
    validator_list = get_pledge_list(other_node.ppos.getValidatorList)
    assert node.node_id in validator_list
    candidate_info = other_node.ppos.getCandidateInfo(node.node_id)
    log.info(candidate_info)
    log.info("Close one node")
    node.stop()
    for i in range(4):
        economic.wait_consensus_blocknum(other_node, number=i)
        candidate_info = other_node.ppos.getCandidateInfo(node.node_id)
        log.info(candidate_info)
        if candidate_info["Ret"]["Released"] < value:
            break
    log.info("Restart the node")
    node.start()
    time.sleep(10)
    result = client.staking.increase_staking(0, address)
    log.info(result)
    assert_code(result, 301103)
    log.info("Next settlement period")
    economic.wait_settlement_blocknum(node, number=2)
    result = client.staking.increase_staking(0, address)
    assert_code(result, 301102)


@allure.title("Increase your holdings with a new wallet")
@pytest.mark.P3
def test_AS_022(client_new_node):
    """
    :param client_new_node_obj:
    :return:
    """
    address, _ = client_new_node.economic.account.generate_account(client_new_node.node.web3,
                                                                   10 ** 18 * 10000000)
    address1, _ = client_new_node.economic.account.generate_account(client_new_node.node.web3,
                                                                    10 ** 18 * 10000000)
    result = client_new_node.staking.create_staking(0, address, address)
    assert_code(result, 0)
    result = client_new_node.staking.increase_staking(0, address1)
    log.info(result)
    assert_code(result, 301006)
