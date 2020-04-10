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


def create_account_amount(client, amount1, amount2):
    # create account1
    lock_address, _ = client.economic.account.generate_account(client.node.web3, amount1)
    # create account2
    release_address, _ = client.economic.account.generate_account(client.node.web3, amount2)
    return lock_address, release_address


def restricting_plan_validation_release(client, economic, node):
    # create account
    address1, _ = economic.account.generate_account(node.web3, von_amount(economic.create_staking_limit, 2))
    # create Restricting Plan
    amount = economic.create_staking_limit
    plan = [{'Epoch': 1, 'Amount': amount}]
    result = client.restricting.createRestrictingPlan(address1, plan, address1)
    assert_code(result, 0)
    return address1


def restricting_plan_validation_staking(client, economic, node):
    # create restricting plan
    address1 = restricting_plan_validation_release(client, economic, node)
    # create staking
    staking_amount = economic.create_staking_limit
    result = client.staking.create_staking(1, address1, address1, amount=staking_amount)
    assert_code(result, 0)
    return address1


@pytest.mark.P2
@pytest.mark.compatibility
def test_UP_FV_001(client_new_node):
    """
    只有一个锁仓期，到达释放期返回解锁金额
    :param client_new_node:
    :return:
    """
    client = client_new_node
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
def test_UP_FV_002(client_new_node):
    """
    只有一个锁仓期，未达释放期返回解锁金额
    :param client_new_node:
    :return:
    """
    client = client_new_node
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
def test_UP_FV_003(client_new_node):
    """
    多个锁仓期，依次部分释放期返回解锁金额
    :param client_new_node:
    :return:
    """
    client = client_new_node
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
def test_UP_FV_004(client_new_node):
    """
    锁仓账户申请质押到释放期后释放锁定金额不足
    :param client_new_node:
    :return:
    """
    client = client_new_node
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


@pytest.mark.P1
def test_UP_FV_005(client_new_node):
    """
    到达释放期释放锁仓金额之后再申请退回质押金
    :param client_new_node:
    :return:
    """
    client = client_new_node
    economic = client.economic
    node = client.node
    # create restricting plan and staking
    address1 = restricting_plan_validation_staking(client, economic, node)
    # Waiting for the end of the settlement period
    economic.wait_settlement_blocknum(node)
    # Application for return of pledge
    result = client.staking.withdrew_staking(address1)
    assert_code(result, 0)
    # view restricting plan
    restricting_info = client.ppos.getRestrictingInfo(address1)
    log.info("restricting plan informtion: {}".format(restricting_info))
    info = restricting_info['Ret']
    assert info['debt'] == economic.create_staking_limit, 'ErrMsg: restricting debt amount {}'.format(
        info['debt'])
    # Waiting for the end of the settlement period
    economic.wait_settlement_blocknum(node, 2)
    # view restricting plan again
    restricting_info = client.ppos.getRestrictingInfo(address1)
    assert_code(restricting_info, 304005)


@pytest.mark.P1
def test_UP_FV_006(client_new_node):
    """
    多个锁仓期，质押一部分锁仓金额再依次释放
    :param client_new_node:
    :return:
    """
    client = client_new_node
    economic = client.economic
    node = client.node
    # create account1
    address1, _ = client.economic.account.generate_account(client.node.web3,
                                                           von_amount(economic.create_staking_limit, 2))
    # create Restricting Plan
    amount1 = economic.create_staking_limit
    amount2 = von_amount(economic.add_staking_limit, 10)
    plan = [{'Epoch': 1, 'Amount': amount1}, {'Epoch': 2, 'Amount': amount2}]
    result = client.restricting.createRestrictingPlan(address1, plan, address1)
    assert_code(result, 0)
    # create staking
    staking_amount = economic.create_staking_limit
    result = client.staking.create_staking(1, address1, address1, amount=staking_amount)
    assert_code(result, 0)
    # view restricting plan
    restricting_info = client.ppos.getRestrictingInfo(address1)
    log.info("restricting plan informtion: {}".format(restricting_info))
    # Waiting for the end of the settlement period
    economic.wait_settlement_blocknum(node)
    # view restricting plan
    restricting_info = client.ppos.getRestrictingInfo(address1)
    log.info("restricting plan informtion: {}".format(restricting_info))
    info = restricting_info['Ret']
    assert info['debt'] == amount1 - amount2, 'ErrMsg: restricting debt amount {}'.format(
        info['debt'])
    assert info['plans'][0]['amount'] == amount2, 'ErrMsg: restricting plans amount {}'.format(
        info['plans'][0]['amount'])
    # Waiting for the end of the settlement period
    economic.wait_settlement_blocknum(node)
    # view restricting plan again
    restricting_info = client.ppos.getRestrictingInfo(address1)
    log.info("restricting plan informtion: {}".format(restricting_info))
    info = restricting_info['Ret']
    assert info['debt'] == amount1, 'ErrMsg: restricting debt amount {}'.format(
        info['debt'])
    assert info['plans'] is None, 'ErrMsg: restricting plans'.format(info['plans'])


@pytest.mark.P1
def test_UP_FV_007(client_new_node):
    """
    锁仓账户申请委托到释放期后释放锁定金额不足
    :param client_new_node:
    :return:
    """
    client = client_new_node
    economic = client.economic
    node = client.node
    # create account
    amount1 = von_amount(economic.create_staking_limit, 2)
    amount2 = von_amount(economic.create_staking_limit, 1)
    address1, address2 = create_account_amount(client, amount1, amount2)
    # create Restricting Plan
    delegate_amount = von_amount(economic.delegate_limit, 10)
    plan = [{'Epoch': 1, 'Amount': delegate_amount}]
    result = client.restricting.createRestrictingPlan(address2, plan, address2)
    assert_code(result, 0)
    # view restricting plan
    restricting_info = client.ppos.getRestrictingInfo(address2)
    log.info("restricting plan informtion: {}".format(restricting_info))
    # create staking
    result = client.staking.create_staking(0, address1, address1)
    assert_code(result, 0)
    # Application for Commission
    result = client.delegate.delegate(1, address2, amount=delegate_amount)
    assert_code(result, 0)
    # view restricting plan
    restricting_info = client.ppos.getRestrictingInfo(address2)
    log.info("restricting plan informtion: {}".format(restricting_info))
    # Waiting for the end of the settlement period
    economic.wait_settlement_blocknum(node)
    # view restricting plan
    restricting_info = client.ppos.getRestrictingInfo(address2)
    log.info("restricting plan informtion: {}".format(restricting_info))
    info = restricting_info['Ret']
    assert info['debt'] == delegate_amount, 'ErrMsg: restricting debt amount {}'.format(
        info['debt'])


@pytest.mark.P1
def test_UP_FV_008(client_new_node):
    """
    到达释放期释放锁仓金额之后再申请赎回委托
    :param client_new_node:
    :return:
    """
    client = client_new_node
    economic = client.economic
    node = client.node
    # create account
    amount1 = von_amount(economic.create_staking_limit, 2)
    amount2 = von_amount(economic.create_staking_limit, 1)
    address1, address2 = create_account_amount(client, amount1, amount2)
    # create Restricting Plan
    delegate_amount = von_amount(economic.delegate_limit, 10)
    plan = [{'Epoch': 1, 'Amount': delegate_amount}]
    result = client.restricting.createRestrictingPlan(address2, plan, address2)
    assert_code(result, 0)
    # create staking
    result = client.staking.create_staking(0, address1, address1)
    assert_code(result, 0)
    # Application for Commission
    result = client.delegate.delegate(1, address2, amount=delegate_amount)
    assert_code(result, 0)
    # Waiting for the end of the settlement period
    economic.wait_settlement_blocknum(node)
    # view restricting plan
    restricting_info = client.ppos.getRestrictingInfo(address2)
    log.info("restricting plan informtion: {}".format(restricting_info))
    info = restricting_info['Ret']
    assert info['debt'] == delegate_amount, 'ErrMsg: restricting debt amount {}'.format(
        info['debt'])
    # Access to pledge information
    candidate_info = client.ppos.getCandidateInfo(node.node_id)
    info = candidate_info['Ret']
    staking_blocknum = info['StakingBlockNum']
    # withdrew delegate
    result = client.delegate.withdrew_delegate(staking_blocknum, address2, amount=delegate_amount)
    assert_code(result, 0)
    # view restricting plan
    restricting_info = client.ppos.getRestrictingInfo(address1)
    assert_code(restricting_info, 304005)


@pytest.mark.P1
def test_UP_FV_009(clients_new_node):
    """
    锁仓账号申请质押，验证人违规被扣除节点自有质押金k
    :param clients_new_node:
    :return:
    """
    client1 = clients_new_node[0]
    client2 = clients_new_node[1]
    economic = client1.economic
    node = client1.node
    # create restricting plan and staking
    address1 = restricting_plan_validation_staking(client1, economic, node)
    # Waiting for the end of the settlement period
    economic.wait_settlement_blocknum(node)
    # Obtain block bonus and pledge bonus
    block_reward, staking_reward = client1.economic.get_current_year_reward(node)
    # Get penalty blocks
    slash_blocks = get_governable_parameter_value(client1, 'slashBlocksReward')
    # view restricting plan
    restricting_info = client1.ppos.getRestrictingInfo(address1)
    log.info("restricting plan informtion: {}".format(restricting_info))
    info = restricting_info['Ret']
    assert info['debt'] == economic.create_staking_limit, 'ErrMsg: restricting debt amount {}'.format(
        info['debt'])
    log.info("Current block height: {}".format(client2.node.eth.blockNumber))
    # stop node
    client1.node.stop()
    # Waiting for a 3 consensus round
    client2.economic.wait_consensus_blocknum(client2.node, 3)
    log.info("Current block height: {}".format(client2.node.eth.blockNumber))
    # view verifier list
    verifier_list = client2.ppos.getVerifierList()
    log.info("verifier_list: {}".format(verifier_list))
    punishment_amonut = int(Decimal(str(block_reward)) * Decimal(str(slash_blocks)))
    log.info("punishment_amonut: {}".format(punishment_amonut))
    # view restricting plan again
    restricting_info = client2.ppos.getRestrictingInfo(address1)
    log.info("restricting plan informtion: {}".format(restricting_info))
    info = restricting_info['Ret']
    if punishment_amonut < economic.create_staking_limit:
        assert (info['balance'] == economic.create_staking_limit - punishment_amonut) or (info['balance'] == economic.create_staking_limit - punishment_amonut * 2), 'ErrMsg: restricting balance amount {}'.format(
            info['balance'])
    else:
        assert_code(restricting_info, 304005)


@pytest.mark.P2
def test_UP_FV_010(client_new_node, reset_environment):
    """
    锁仓验证人违规被剔除验证人列表，申请质押节点
    :param client_new_node:
    :return:
    """
    client1 = client_new_node
    economic = client1.economic
    node = client1.node
    # create account
    amount1 = von_amount(economic.create_staking_limit, 3)
    amount2 = von_amount(economic.create_staking_limit, 1)
    address1, report_address = create_account_amount(client1, amount1, amount2)
    # create Restricting Plan
    delegate_amount = von_amount(economic.create_staking_limit, 2)
    plan = [{'Epoch': 3, 'Amount': delegate_amount}]
    result = client1.restricting.createRestrictingPlan(address1, plan, address1)
    assert_code(result, 0)
    # view restricting plan again
    restricting_info = client1.ppos.getRestrictingInfo(address1)
    log.info("restricting plan informtion: {}".format(restricting_info))
    # create staking
    result = client1.staking.create_staking(1, address1, address1)
    assert_code(result, 0)
    # view restricting plan again
    restricting_info = client1.ppos.getRestrictingInfo(address1)
    log.info("restricting plan informtion: {}".format(restricting_info))
    # Waiting for the end of the settlement period
    economic.wait_settlement_blocknum(node)
    #
    for i in range(4):
        result = check_node_in_list(node.node_id, client1.ppos.getValidatorList)
        log.info("Current node in consensus list status：{}".format(result))
        if result:
            # view current block
            current_block = node.eth.blockNumber
            log.info("Current block: {}".format(current_block))
            # Report prepareblock signature
            report_information = mock_duplicate_sign(1, node.nodekey, node.blsprikey, current_block)
            log.info("Report information: {}".format(report_information))
            result = client1.duplicatesign.reportDuplicateSign(1, report_information, report_address)
            assert_code(result, 0)
            time.sleep(3)
            # create staking
            result = client1.staking.create_staking(1, address1, address1)
            assert_code(result, 301101)
            break
        else:
            # wait consensus block
            client1.economic.wait_consensus_blocknum(node)


@pytest.mark.P2
def test_UP_FV_011(client_new_node, reset_environment):
    """
    锁仓验证人违规被剔除验证人列表，申请委托节点
    :param client_new_node:
    :return:
    """
    client1 = client_new_node
    economic = client1.economic
    node = client1.node
    # create account
    amount1 = von_amount(economic.create_staking_limit, 2)
    amount2 = von_amount(economic.create_staking_limit, 1)
    address1, report_address = create_account_amount(client1, amount1, amount2)
    # create Restricting Plan
    delegate_amount = von_amount(economic.delegate_limit, 10)
    plan = [{'Epoch': 3, 'Amount': delegate_amount}]
    result = client1.restricting.createRestrictingPlan(report_address, plan, report_address)
    assert_code(result, 0)
    # create staking
    result = client1.staking.create_staking(0, address1, address1)
    assert_code(result, 0)
    # Application for Commission
    result = client1.delegate.delegate(1, report_address)
    assert_code(result, 0)
    # Waiting for the end of the settlement period
    economic.wait_settlement_blocknum(node)
    #
    for i in range(4):
        result = check_node_in_list(node.node_id, client1.ppos.getValidatorList)
        log.info("Current node in consensus list status：{}".format(result))
        if result:
            # view current block
            current_block = node.eth.blockNumber
            log.info("Current block: {}".format(current_block))
            # Report prepareblock signature
            report_information = mock_duplicate_sign(1, node.nodekey, node.blsprikey, current_block)
            log.info("Report information: {}".format(report_information))
            result = client1.duplicatesign.reportDuplicateSign(1, report_information, report_address)
            assert_code(result, 0)
            time.sleep(3)
            # Application for Commission
            result = client1.delegate.delegate(1, report_address)
            assert_code(result, 301103)
            break
        else:
            # wait consensus block
            client1.economic.wait_consensus_blocknum(node)


@pytest.mark.P2
def test_UP_FV_012(client_new_node, reset_environment):
    """
    锁仓验证人违规被剔除验证人列表，申请增持质押
    :param client_new_node:
    :return:
    """
    client1 = client_new_node
    economic = client1.economic
    node = client1.node
    # create account
    amount1 = von_amount(economic.create_staking_limit, 3)
    amount2 = von_amount(economic.create_staking_limit, 1)
    address1, report_address = create_account_amount(client1, amount1, amount2)
    # create Restricting Plan
    staking_amount = von_amount(economic.create_staking_limit, 2)
    plan = [{'Epoch': 3, 'Amount': staking_amount}]
    result = client1.restricting.createRestrictingPlan(address1, plan, address1)
    assert_code(result, 0)
    # create staking
    result = client1.staking.create_staking(1, address1, address1)
    assert_code(result, 0)
    # Apply for additional pledge
    result = client1.staking.increase_staking(1, address1)
    assert_code(result, 0)
    # Waiting for the end of the settlement period
    economic.wait_settlement_blocknum(node)
    #
    for i in range(4):
        result = check_node_in_list(node.node_id, client1.ppos.getValidatorList)
        log.info("Current node in consensus list status：{}".format(result))
        if result:
            # view current block
            current_block = node.eth.blockNumber
            log.info("Current block: {}".format(current_block))
            # Report prepareblock signature
            report_information = mock_duplicate_sign(1, node.nodekey, node.blsprikey, current_block)
            log.info("Report information: {}".format(report_information))
            result = client1.duplicatesign.reportDuplicateSign(1, report_information, report_address)
            assert_code(result, 0)
            time.sleep(3)
            # Apply for additional pledge
            result = client1.staking.increase_staking(1, address1)
            assert_code(result, 301103)
            break
        else:
            # wait consensus block
            client1.economic.wait_consensus_blocknum(node)


@pytest.mark.P2
def test_UP_FV_013(client_new_node, reset_environment):
    """
    锁仓验证人违规被剔除验证人列表，申请退回质押金
    :param client_new_node:
    :return:
    """
    client1 = client_new_node
    economic = client1.economic
    node = client1.node
    # create account
    amount1 = von_amount(economic.create_staking_limit, 3)
    amount2 = von_amount(economic.create_staking_limit, 1)
    address1, report_address = create_account_amount(client1, amount1, amount2)
    # create Restricting Plan
    delegate_amount = von_amount(economic.create_staking_limit, 2)
    plan = [{'Epoch': 3, 'Amount': delegate_amount}]
    result = client1.restricting.createRestrictingPlan(address1, plan, address1)
    assert_code(result, 0)
    # view restricting plan again
    restricting_info = client1.ppos.getRestrictingInfo(address1)
    log.info("restricting plan informtion: {}".format(restricting_info))
    # create staking
    result = client1.staking.create_staking(1, address1, address1)
    assert_code(result, 0)
    # view restricting plan again
    restricting_info = client1.ppos.getRestrictingInfo(address1)
    log.info("restricting plan informtion: {}".format(restricting_info))
    # Waiting for the end of the settlement period
    economic.wait_settlement_blocknum(node)
    #
    for i in range(4):
        result = check_node_in_list(node.node_id, client1.ppos.getValidatorList)
        log.info("Current node in consensus list status：{}".format(result))
        if result:
            # view current block
            current_block = node.eth.blockNumber
            log.info("Current block: {}".format(current_block))
            # Report prepareblock signature
            report_information = mock_duplicate_sign(1, node.nodekey, node.blsprikey, current_block)
            log.info("Report information: {}".format(report_information))
            result = client1.duplicatesign.reportDuplicateSign(1, report_information, report_address)
            assert_code(result, 0)
            time.sleep(3)
            # withdrew staking
            result = client1.staking.withdrew_staking(address1)
            assert_code(result, 301103)
            break
        else:
            # wait consensus block
            client1.economic.wait_consensus_blocknum(node)


@pytest.mark.P2
def test_UP_FV_014(client_new_node, reset_environment):
    """
    锁仓验证人违规被剔除验证人列表，申请赎回委托金
    :param client_new_node:
    :param reset_environment:
    :return:
    """
    client = client_new_node
    economic = client.economic
    node = client.node
    # create account
    amount1 = von_amount(economic.create_staking_limit, 2)
    amount2 = von_amount(economic.create_staking_limit, 1)
    address1, report_address = create_account_amount(client, amount1, amount2)
    # create Restricting Plan
    delegate_amount = von_amount(economic.delegate_limit, 10)
    plan = [{'Epoch': 3, 'Amount': delegate_amount}]
    result = client.restricting.createRestrictingPlan(report_address, plan, report_address)
    assert_code(result, 0)
    # create staking
    result = client.staking.create_staking(0, address1, address1)
    assert_code(result, 0)
    # Application for Commission
    result = client.delegate.delegate(1, report_address)
    assert_code(result, 0)
    # Waiting for the end of the settlement period
    economic.wait_settlement_blocknum(node)
    #
    for i in range(4):
        result = check_node_in_list(node.node_id, client.ppos.getValidatorList)
        log.info("Current node in consensus list status：{}".format(result))
        if result:
            # view current block
            current_block = node.eth.blockNumber
            log.info("Current block: {}".format(current_block))
            # Report prepareblock signature
            report_information = mock_duplicate_sign(1, node.nodekey, node.blsprikey, current_block)
            log.info("Report information: {}".format(report_information))
            result = client.duplicatesign.reportDuplicateSign(1, report_information, report_address)
            assert_code(result, 0)
            time.sleep(3)
            # Access to pledge information
            candidate_info = client.ppos.getCandidateInfo(node.node_id)
            info = candidate_info['Ret']
            staking_blocknum = info['StakingBlockNum']
            # withdrew delegate
            result = client.delegate.withdrew_delegate(staking_blocknum, report_address)
            assert_code(result, 0)
            break
        else:
            # wait consensus block
            client.economic.wait_consensus_blocknum(node)


@pytest.mark.P2
def test_UP_FV_015(client_new_node):
    """
    锁仓账户申请之后到达释放期，账户锁仓不足再新增新的锁仓计划
    :param client_new_node:
    :return:
    """
    client = client_new_node
    economic = client.economic
    node = client.node
    # create restricting plan and staking
    address1 = restricting_plan_validation_staking(client, economic, node)
    # create account2
    address2, _ = economic.account.generate_account(node.web3, von_amount(economic.create_staking_limit, 2))
    # Waiting for the end of the settlement period
    economic.wait_settlement_blocknum(node)
    # view restricting plan
    restricting_info = client.ppos.getRestrictingInfo(address1)
    log.info("restricting plan informtion: {}".format(restricting_info))
    info = restricting_info['Ret']
    assert info['debt'] == economic.create_staking_limit, 'ErrMsg: restricting debt amount {}'.format(
        info['debt'])
    # create Restricting Plan
    delegate_amount = von_amount(economic.create_staking_limit, 1)
    plan = [{'Epoch': 1, 'Amount': delegate_amount}]
    result = client.restricting.createRestrictingPlan(address1, plan, address2)
    assert_code(result, 0)
    # view restricting plan
    restricting_info = client.ppos.getRestrictingInfo(address1)
    log.info("restricting plan informtion: {}".format(restricting_info))
    info = restricting_info['Ret']
    assert info['debt'] == 0, 'ErrMsg: restricting debt amount {}'.format(info['debt'])


@pytest.mark.P1
def test_UP_FV_016(client_new_node):
    """
    自由资金质押，锁仓再增持
    :param client_new_node:
    :return:
    """
    client = client_new_node
    economic = client.economic
    node = client.node
    # create account
    amount1 = von_amount(economic.create_staking_limit, 2)
    amount2 = von_amount(economic.create_staking_limit, 1)
    address1, address2 = create_account_amount(client, amount1, amount2)
    # create Restricting Plan
    delegate_amount = von_amount(economic.add_staking_limit, 10)
    plan = [{'Epoch': 3, 'Amount': delegate_amount}]
    result = client.restricting.createRestrictingPlan(address1, plan, address1)
    assert_code(result, 0)
    # create staking
    result = client.staking.create_staking(0, address1, address1)
    assert_code(result, 0)
    # Apply for additional pledge
    result = client.staking.increase_staking(1, address1)
    assert_code(result, 0)
    # Waiting for the end of the settlement period
    economic.wait_settlement_blocknum(node)
    # Apply for additional pledge
    result = client.staking.increase_staking(1, address1)
    assert_code(result, 0)


@pytest.mark.P1
def test_UP_FV_017(client_new_node):
    """
    锁仓账户质押，自由资金再增持
    :param client_new_node:
    :return:
    """
    client = client_new_node
    economic = client.economic
    node = client.node
    # create account
    amount1 = von_amount(economic.create_staking_limit, 2)
    amount2 = von_amount(economic.create_staking_limit, 1)
    address1, address2 = create_account_amount(client, amount1, amount2)
    # create Restricting Plan
    delegate_amount = von_amount(economic.create_staking_limit, 1)
    plan = [{'Epoch': 3, 'Amount': delegate_amount}]
    result = client.restricting.createRestrictingPlan(address1, plan, address1)
    assert_code(result, 0)
    # create staking
    result = client.staking.create_staking(1, address1, address1)
    assert_code(result, 0)
    # Apply for additional pledge
    result = client.staking.increase_staking(0, address1)
    assert_code(result, 0)
    # Waiting for the end of the settlement period
    economic.wait_settlement_blocknum(node)
    # Apply for additional pledge
    result = client.staking.increase_staking(0, address1)
    assert_code(result, 0)


@pytest.mark.P1
def test_UP_FV_018(client_new_node):
    """
    账户自由金额和锁仓金额申请委托同一个验证人，再申请赎回
    :param client_new_node:
    :return:
    """
    client = client_new_node
    economic = client.economic
    node = client.node
    # create account
    amount1 = von_amount(economic.create_staking_limit, 2)
    amount2 = von_amount(economic.create_staking_limit, 1)
    address1, address2 = create_account_amount(client, amount1, amount2)
    # create Restricting Plan
    delegate_amount = von_amount(economic.delegate_limit, 10)
    plan = [{'Epoch': 1, 'Amount': delegate_amount}]
    result = client.restricting.createRestrictingPlan(address2, plan, address2)
    assert_code(result, 0)
    # create staking
    result = client.staking.create_staking(0, address1, address1)
    assert_code(result, 0)
    # Lock in amount Application for delegate
    result = client.delegate.delegate(1, address2, amount=delegate_amount)
    assert_code(result, 0)
    # Free amount Application for delegate
    result = client.delegate.delegate(0, address2, amount=delegate_amount)
    assert_code(result, 0)
    # Waiting for the end of the settlement period
    economic.wait_settlement_blocknum(node)
    # Access to pledge information
    candidate_info = client.ppos.getCandidateInfo(node.node_id)
    info = candidate_info['Ret']
    staking_blocknum = info['StakingBlockNum']
    # withdrew delegate
    withdrew_amount = von_amount(economic.delegate_limit, 15)
    result = client.delegate.withdrew_delegate(staking_blocknum, address2, amount=withdrew_amount)
    assert_code(result, 0)
    # view restricting plan
    restricting_info = client.ppos.getRestrictingInfo(address2)
    log.info("restricting plan informtion: {}".format(restricting_info))
    info = restricting_info['Ret']
    assert info['debt'] == von_amount(economic.delegate_limit, 5), 'ErrMsg: restricting debt amount {}'.format(info['debt'])


@pytest.mark.P1
def test_UP_FV_019(client_new_node):
    """
    账户自由金额和锁仓金额申请委托同一个验证人，再申请赎回
    :param client_new_node:
    :return:
    """
    client = client_new_node
    economic = client.economic
    node = client.node
    # create account
    amount = economic.create_staking_limit
    first_address, second_address = create_account_amount(client, amount, amount)
    delegate_amount = von_amount(economic.delegate_limit, 10)
    plan = [{'Epoch': 2, 'Amount': delegate_amount}]

    # create Restricting Plan1
    result = client.restricting.createRestrictingPlan(first_address, plan, first_address)
    assert_code(result, 0)
    economic.wait_settlement_blocknum(node)
    log.info("Current block height：{}".format(node.eth.blockNumber))

    # create Restricting Plan2
    first_balance1 = node.eth.getBalance(first_address)
    log.info("first_balance1: {}".format(first_balance1))
    result = client.restricting.createRestrictingPlan(second_address, plan, second_address)
    assert_code(result, 0)

    economic.wait_settlement_blocknum(node)
    log.info("Current block height：{}".format(node.eth.blockNumber))
    second_balance1 = node.eth.getBalance(second_address)
    log.info("second_balance1: {}".format(second_balance1))
    restricting_info = client.ppos.getRestrictingInfo(first_address)
    log.info("restricting plan1 informtion: {}".format(restricting_info))
    first_balance2 = node.eth.getBalance(first_address)
    log.info("first_balance2: {}".format(first_balance2))
    assert first_balance2 == first_balance1 + delegate_amount

    economic.wait_settlement_blocknum(node)
    log.info("Current block height：{}".format(node.eth.blockNumber))
    restricting_info = client.ppos.getRestrictingInfo(second_address)
    log.info("restricting plan2 informtion: {}".format(restricting_info))
    second_balance2 = node.eth.getBalance(second_address)
    log.info("second_balance2: {}".format(second_balance2))
    assert second_balance2 == second_balance1 + delegate_amount
