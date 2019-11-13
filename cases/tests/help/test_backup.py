from decimal import Decimal

import allure
import pytest
import time
from common.log import log
from client_sdk_python import Web3
from tests.lib.utils import get_pledge_list, get_block_count_number, assert_code
from common.key import generate_key


def calculate(big_int, mul):
    return int(Decimal(str(big_int))*Decimal(mul))


@pytest.fixture()
def staking_client(client_new_node_obj):
    amount = calculate(client_new_node_obj.economic.create_staking_limit, 5)
    staking_amount = calculate(client_new_node_obj.economic.create_staking_limit, 2)
    staking_address, _ = client_new_node_obj.economic.account.generate_account(client_new_node_obj.node.web3, amount)
    delegate_address, _ = client_new_node_obj.economic.account.generate_account(client_new_node_obj.node.web3,
                                                                                client_new_node_obj.economic.add_staking_limit * 2)
    client_new_node_obj.staking.create_staking(0, staking_address, staking_address, amount=staking_amount)
    setattr(client_new_node_obj, "staking_address", staking_address)
    setattr(client_new_node_obj, "delegate_address", delegate_address)
    setattr(client_new_node_obj, "amount", amount)
    setattr(client_new_node_obj, "staking_amount", staking_amount)
    return client_new_node_obj


@allure.title("验证人退回质押金（锁定期）")
@pytest.mark.P1
def test_back_unstaking(staking_client):
    """
    验证人退回质押金（未到达可解锁期）
    质押成为下个周期验证人，退出后，下一个结算周期退出
    """
    client = staking_client
    staking_address = client.staking_address
    node = client.node
    economic = client.economic
    staking_address_balance = node.eth.getBalance(staking_address)
    log.info(staking_address_balance)
    economic.wait_settlement_blocknum(node)
    log.info("查询第2个结算周期的验证人")
    node_list = get_pledge_list(client.ppos.getVerifierList)
    log.info(node_list)
    assert node.node_id in node_list
    log.info("节点1在第2结算周期，锁定期申请退回")
    client.staking.withdrew_staking(staking_address)
    """发起退回消耗一定gas"""
    staking_address_balance_1 = node.eth.getBalance(staking_address)
    log.info(staking_address_balance_1)
    log.info("进入第3个结算周期")
    economic.wait_settlement_blocknum(node)
    staking_address_balance_2 = node.eth.getBalance(staking_address)
    log.info(staking_address_balance_2)
    node_list = get_pledge_list(client.ppos.getVerifierList)
    log.info(node_list)
    assert node.node_id not in node_list
    log.info("进入第4个结算周期")
    economic.wait_settlement_blocknum(node)
    msg = client.ppos.getCandidateInfo(node.node_id)
    log.info(msg)
    staking_address_balance_3 = node.eth.getBalance(staking_address)
    log.info(staking_address_balance_3)
    """ 第3个结算周期结束后质押金额已退 """
    log.info(staking_address_balance_3 - staking_address_balance_1)
    assert staking_address_balance_3 - staking_address_balance_1 > client.staking_amount, "退回的交易金额应大于返回质押金额"


@allure.title("退出质押后锁定期不能增持和委托")
@pytest.mark.P1
def test_locked_quit_addstaking_delegate(staking_client):
    """
    退出质押后不能增持和委托
    """
    client = staking_client
    node = client.node
    staking_address = client.staking_address
    economic = client.economic
    log.info("进入锁定期")
    economic.wait_settlement_blocknum(node)
    log.info("节点1退出质押")
    client.staking.withdrew_staking(staking_address)
    log.info("节点1做增持")
    msg = client.staking.increase_staking(0, staking_address, amount=economic.add_staking_limit)
    assert_code(msg, 301103)
    log.info("节点1做委托")
    msg = client.delegate.delegate(0, client.delegate_address)
    assert_code(msg, 301103)


@allure.title("最高惩罚后,返回金额&重新质押、委托、赎回")
@pytest.mark.P1
def test_punishment_refund(staking_client, global_test_env):
    """
    最高惩罚后,返回金额
    """
    other_node = global_test_env.get_rand_node()
    client = staking_client
    staking_address = client.staking_address
    node = client.node
    economic = client.economic
    balance = node.eth.getBalance(staking_address)
    log.info(balance)
    msg = client.ppos.getCandidateInfo(node.node_id)
    log.info(msg)
    log.info("把新的验证人节点停掉")
    node.stop()
    log.info("进入到下个结算周期")
    economic.wait_settlement_blocknum(other_node)
    msg = get_pledge_list(other_node.ppos.getCandidateList)
    log.info("实时验证人列表{}".format(msg))
    msg = get_pledge_list(other_node.ppos.getVerifierList)
    log.info("当前结算周期验证人{}".format(msg))
    msg = get_pledge_list(other_node.ppos.getValidatorList)
    log.info("当前共识轮验证人{}".format(msg))
    log.info("进入到下个结算周期")
    economic.wait_settlement_blocknum(other_node)
    msg = get_pledge_list(other_node.ppos.getCandidateList)
    log.info("实时验证人列表{}".format(msg))
    verifier_list = get_pledge_list(other_node.ppos.getVerifierList)
    log.info("当前结算周期验证人{}".format(verifier_list))
    assert node.node_id not in verifier_list, "预期退出验证人列表"
    candidate_info = other_node.ppos.getCandidateInfo(node.node_id)
    balance_before = other_node.eth.getBalance(staking_address)
    log.info("查询被惩罚后的账户余额:{}".format(balance_before))
    log.info("进入到下个结算周期")
    economic.wait_settlement_blocknum(other_node, 1)
    time.sleep(10)
    balance_after = other_node.eth.getBalance(staking_address)
    log.info("被罚后剩余的金额退回到账户后的余额:{}".format(balance_after))
    assert balance_before + candidate_info["Ret"]["Released"] == balance_after, "被罚出移出验证人后，金额退还异常"
    msg = other_node.ppos.getCandidateInfo(node.node_id)
    log.info(msg)
    node.start()
    time.sleep(10)
    staking_result = client.staking.create_staking(0, staking_address, staking_address)
    assert_code(staking_result, 0)
    candidate_info = node.ppos.getCandidateInfo(node.node_id)
    log.info(candidate_info)
    staking_blocknum = candidate_info["Ret"]["StakingBlockNum"]
    log.info("钱包3委托给节点1")
    msg = client.delegate.delegate(0, client.delegate_address, node.node_id)
    assert_code(msg, 0)
    msg = client.delegate.withdrew_delegate(staking_blocknum, client.delegate_address, node.node_id)
    assert_code(msg, 0)


@allure.title("退出中修改质押信息")
@pytest.mark.P2
def test_quiting_updateStakingInfo(staking_client):
    """
    退出中修改质押信息
    """
    node_name = "wuyiqin"
    client = staking_client
    staking_address = client.staking_address
    node = client.node
    economic = client.economic
    economic.wait_settlement_blocknum(node)
    msg = client.staking.withdrew_staking(staking_address)
    log.info(msg)
    msg = node.ppos.getCandidateInfo(node.node_id)
    log.info(msg)
    log.info("节点2修改节点信息")
    client.staking.cfg.node_name = node_name
    msg = client.staking.edit_candidate(staking_address, staking_address)
    assert_code(msg, 301103)


@allure.title("已退出修改质押信息")
@pytest.mark.P2
def test_quited_updateStakingInfo(staking_client):
    """
    已退出修改质押信息
    """
    node_name = "wuyiqin"
    client = staking_client
    staking_address = client.staking_address
    node = client.node
    economic = client.economic
    economic.wait_settlement_blocknum(node)
    msg = client.staking.withdrew_staking(staking_address)
    log.info(msg)
    msg = node.ppos.getCandidateInfo(node.node_id)
    log.info(msg)
    economic.wait_settlement_blocknum(node, 2)
    log.info("节点2修改节点信息")
    client.staking.cfg.node_name = node_name
    msg = client.staking.edit_candidate(staking_address, staking_address)
    assert_code(msg, 301102)


@allure.title("退出验证人后，返回质押金+出块奖励+质押奖励")
@pytest.mark.P1
def test_into_quit_block_reward(staking_client):
    """
    成为验证人后，有质押奖励和出块奖励
    """
    client = staking_client
    staking_address = client.staking_address
    node = client.node
    economic = client.economic
    economic.wait_settlement_blocknum(node)
    log.info("进入下个周期")
    block_reward, staking_reward = economic.get_current_year_reward(node)
    msg = client.staking.withdrew_staking(staking_address)
    log.info(msg)
    balance_1 = node.eth.getBalance(staking_address)
    log.info(balance_1)
    log.info("进入下个周期")
    economic.wait_settlement_blocknum(node, 2)
    balance_2 = node.eth.getBalance(staking_address)
    log.info(balance_2)
    verifier_list = get_pledge_list(node.ppos.getVerifierList)
    log.info("当前验证人列表：{}".format(verifier_list))
    validator_list = get_pledge_list(node.ppos.getValidatorList)
    log.info("当前共识验证人列表：{}".format(validator_list))
    block_number = get_block_count_number(node, economic.settlement_size)
    sum_block_reward = calculate(block_reward, block_number)
    reward_sum = sum_block_reward + staking_reward
    log.info("奖励的总金额{}".format(reward_sum))
    assert balance_1 + reward_sum + client.staking_amount == balance_2, "奖励金额异常"


@allure.title("验证人申请退回质押金（犹豫期）")
@pytest.mark.P0
def test_back_unStaking(staking_client):
    """
    用例id 81 验证人申请退回质押金（犹豫期）
    """
    client = staking_client
    staking_address = client.staking_address
    node = client.node
    balance_before = node.eth.getBalance(staking_address)
    log.info("节点4对应的钱包余额{}".format(balance_before))
    client.staking.withdrew_staking(staking_address)
    balance_after = node.eth.getBalance(staking_address)
    log.info("节点4退出质押后钱包余额{}".format(balance_after))
    assert balance_after > balance_before, "退出质押后，钱包余额未增加"
    log.info("因为质押消耗的gas值大于撤销质押的gas值")
    assert balance_after > client.amount - 10**18
    node_list = get_pledge_list(node.ppos.getCandidateList)
    assert node.node_id not in node_list, "验证节点退出异常"


@allure.title("发起撤销质押（质押金+增持金额））")
@pytest.mark.P1
def test_unstaking_all(staking_client):
    client = staking_client
    node = client.node
    staking_address = client.staking_address
    economic = client.economic
    value_before = client.amount
    log.info("发起质押前的余额{}".format(value_before))

    log.info("进入第2个结算周期,节点1增持金额")
    economic.wait_settlement_blocknum(node)
    client.staking.increase_staking(0, staking_address)
    value2 = node.eth.getBalance(staking_address)
    log.info("做了质押+增持后的余额{}".format(value2))
    log.info("进入第3个结算周期,节点6发起退回")
    economic.wait_settlement_blocknum(node)
    value3 = node.eth.getBalance(staking_address)
    log.info("第3个周期的余额{}".format(value3))
    client.staking.withdrew_staking(staking_address)
    log.info("进入第4个结算周期")
    economic.wait_settlement_blocknum(node)
    value4 = node.eth.getBalance(staking_address)
    log.info("第4个结算周期的余额(包括第3周期的奖励){}".format(value4))
    log.info("进入第5个结算周期")
    economic.wait_settlement_blocknum(node)
    value5 = node.eth.getBalance(staking_address)
    log.info("到了解锁期退回质押+增持后的余额:{}".format(value5))
    log.info(value5 - value_before)
    amount_sum = client.staking_amount + economic.add_staking_limit
    assert value5 > value_before, "出块奖励异常"
    assert value5 > amount_sum, "解锁期的余额大于锁定期的余额+质押+增持金额，但是发生异常"


@allure.title("验证人申请退回质押金（犹豫期+锁定期）")
@pytest.mark.P1
def test_520(staking_client):
    """
    验证人申请退回质押金（犹豫期+锁定期）
    """
    client = staking_client
    staking_address = client.staking_address
    node = client.node
    economic = client.economic
    log.info("进入下个周期")
    economic.wait_settlement_blocknum(node)
    msg = client.staking.increase_staking(0, staking_address)
    assert_code(msg, 0)
    msg = node.ppos.getCandidateInfo(node.node_id)
    log.info("质押信息{}".format(msg))
    assert msg["Ret"]["Shares"] == client.staking_amount + economic.add_staking_limit, "预期显示质押金额+增持金额"
    assert msg["Ret"]["Released"] == client.staking_amount, "预期显示质押金额"
    assert msg["Ret"]["ReleasedHes"] == economic.add_staking_limit, "预期增持金额显示在犹豫期"
    block_reward, staking_reward = economic.get_current_year_reward(node)

    balance = node.eth.getBalance(staking_address)
    log.info("发起退质押前的余额{}".format(balance))

    log.info("节点1在第2周期发起退回质押")
    msg = client.staking.withdrew_staking(staking_address)
    assert_code(msg, 0)
    msg = node.ppos.getCandidateInfo(node.node_id)
    log.info("发起退回后质押信息{}".format(msg))
    assert msg["Ret"]["ReleasedHes"] == 0, "预期增持的金额已退回，显示0"
    balance1 = node.eth.getBalance(client.staking_address)
    log.info(balance1)
    log.info("进入第3个周期")
    economic.wait_settlement_blocknum(node, 2)

    balance2 = node.eth.getBalance(staking_address)
    log.info(balance2)

    block_number = get_block_count_number(node, economic.settlement_size * 2)
    sum_block_reward = calculate(block_reward, block_number)
    reward_sum = sum_block_reward + staking_reward
    log.info("奖励的总金额{}".format(reward_sum))
    assert balance1 + reward_sum + client.staking_amount == balance2, "奖励金额异常"


@allure.title("撤销5种身份候选人，验证人，共识验证人，不存在的候选人，已失效的候选人")
@pytest.mark.P1
def test_backup_identity(client_new_node_obj_list):
    """
    由于其他用例有验证过退质押的金额，这里不做断言
    :param status:
    0: 候选人
    1:验证人
    2:共识验证人
    3:不存在的候选人
    4：已失效的候选人
    :return:
    """
    client_a = client_new_node_obj_list[0]
    node_a = client_a.node
    client_b = client_new_node_obj_list[1]
    node_b = client_b.node
    amount_a = client_b.economic.create_staking_limit * 6
    amount_b = client_b.economic.create_staking_limit * 7
    amount = Web3.toWei(amount_a + amount_b + 10, "ether")
    address, _ = client_a.economic.account.generate_account(node_a.web3, amount)
    msg = client_a.staking.create_staking(0, address, address, amount=amount_a)

    assert_code(msg, 0)
    msg = client_b.staking.create_staking(0, address, address, amount=amount_b)

    assert_code(msg, 0)
    log.info("进入下个周期")
    client_b.economic.wait_settlement_blocknum(node_b)
    msg = client_b.staking.withdrew_staking(address)
    assert_code(msg, 0)


def test_withdrew_staking_001(staking_client):

    client = staking_client
    staking_address = client.staking_address
    node = client.node
    economic = client.economic
    log.info("进入下个周期")
    economic.wait_settlement_blocknum(node, 1)
    verifier_list = get_pledge_list(node.ppos.getVerifierList)
    log.info(log.info("当前结算周期验证人{}".format(verifier_list)))
    assert node.node_id in verifier_list
    msg = client.staking.withdrew_staking(staking_address)
    assert_code(msg, 0)


def test_withdrew_staking_002(staking_client):
    client = staking_client
    node = client.node
    economic = client.economic
    staking_address = client.staking_address
    log.info("进入下个周期")
    economic.wait_settlement_blocknum(node)
    log.info("进入下一个共识轮")
    economic.wait_consensus_blocknum(node)

    validator_list = get_pledge_list(node.ppos.getValidatorList)
    log.info("共识验证人列表:{}".format(validator_list))
    assert node.node_id in validator_list
    msg = client.staking.withdrew_staking(staking_address)
    assert_code(msg, 0)


def test_withdrew_staking_003(staking_client):
    _, node_id = generate_key()
    msg = staking_client.staking.withdrew_staking(staking_client.staking_address, node_id=node_id)
    log.info(msg)
    assert_code(msg, 301102)


def test_withdrew_staking_004(staking_client):
    client = staking_client
    node = client.node
    staking_address = client.staking_address
    msg = client.staking.withdrew_staking(staking_address)
    assert_code(msg, 0)
    msg = node.ppos.getCandidateInfo(node.node_id)
    assert msg["Ret"] == "Query candidate info failed:Candidate info is not found", "预期退质押成功；质押信息被删除"
    msg = client.staking.withdrew_staking(staking_address)
    assert_code(msg, 301102)


@allure.title("自由账户质押+锁仓账户增持(犹豫期退质押)")
@pytest.mark.P1
def test_006(staking_client):
    client = staking_client
    staking_address = client.staking_address
    node = client.node
    economic = client.economic
    log.info("创建锁仓计划")
    lockup_amount = economic.add_staking_limit * 2
    plan = [{'Epoch': 1, 'Amount': lockup_amount}]
    msg = client.restricting.createRestrictingPlan(staking_address, plan, economic.account.account_with_money["address"])
    assert_code(msg, 0)
    locked_info = client.ppos.getRestrictingInfo(staking_address)
    log.info(locked_info)
    before_create_balance = client.amount
    log.info("发起质押前的余额{}".format(before_create_balance))

    msg = client.staking.increase_staking(1, staking_address)
    assert_code(msg, 0)
    msg = client.ppos.getCandidateInfo(node.node_id)
    log.info("查询质押情况{}".format(msg))
    log.info("发起撤销质押")
    msg = client.staking.withdrew_staking(staking_address)
    assert_code(msg, 0)

    after_balance_1 = node.eth.getBalance(staking_address)
    log.info("犹豫期发起退回后的余额{}".format(after_balance_1))
    """退回后的余额肯定小于质押前的余额，消耗小于1 eth"""
    assert before_create_balance - after_balance_1 < Web3.toWei(1, "ether"), "退回金额异常"
    locked_info = client.ppos.getRestrictingInfo(staking_address)
    log.info(locked_info)

    msg = client.ppos.getCandidateInfo(node.node_id)
    assert_code(msg, 301204)
    log.info("进入下个周期")
    economic.wait_settlement_blocknum(node)
    locked_info = client.ppos.getRestrictingInfo(staking_address)
    log.info(locked_info)
    after_account = node.eth.getBalance(staking_address)
    log.info("锁仓释放后的账户余额{}".format(after_account))
    assert after_account - after_balance_1 == lockup_amount, "锁仓退回金额异常"


@allure.title("自由账户质押+锁仓账户增持(锁定期退质押)")
@pytest.mark.P1
def test_007(staking_client):
    client = staking_client
    node = client.node
    staking_address = client.staking_address
    economic = client.economic
    log.info("创建锁仓计划")
    lockup_amount = economic.add_staking_limit * 2
    plan = [{'Epoch': 1, 'Amount': lockup_amount}]
    msg = client.restricting.createRestrictingPlan(staking_address, plan, economic.account.account_with_money["address"])
    assert_code(msg, 0)
    locked_info = client.ppos.getRestrictingInfo(staking_address)
    log.info(locked_info)
    before_create_balance = client.amount
    log.info("发起质押前的余额{}".format(before_create_balance))

    msg = client.staking.increase_staking(1, staking_address)
    assert_code(msg, 0)
    economic.wait_settlement_blocknum(node)

    msg = client.ppos.getCandidateInfo(node.node_id)
    log.info("查询质押情况{}".format(msg))
    assert msg["Ret"]["Shares"] == client.staking_amount + economic.add_staking_limit
    assert msg["Ret"]["Released"] == client.staking_amount
    assert msg["Ret"]["RestrictingPlan"] == economic.add_staking_limit

    block_reward, staking_reward = economic.get_current_year_reward(node)
    msg = client.staking.withdrew_staking(staking_address)
    assert_code(msg, 0)
    balance_withdrew = node.eth.getBalance(staking_address)
    log.info("第2个周期发起撤销后的余额{}".format(balance_withdrew))
    log.info("进入第3个周期")
    economic.wait_settlement_blocknum(node)

    balance_settlement = node.eth.getBalance(staking_address)
    log.info("第3个周期发起撤销后的余额{}".format(balance_settlement))

    log.info("进入第4个周期")
    economic.wait_settlement_blocknum(node, 1)

    balance_settlement_2 = node.eth.getBalance(staking_address)
    log.info("第4个周期发起撤销后的余额{}".format(balance_settlement_2))

    """算出块奖励+质押奖励"""
    log.info("以下为获取节点2出的块数")
    block_number = get_block_count_number(node, economic.settlement_size * 2)
    sum_block_reward = calculate(block_reward, block_number)
    reward_sum = sum_block_reward + staking_reward
    log.info("奖励的总金额{}".format(reward_sum))
    assert before_create_balance + reward_sum + lockup_amount - balance_settlement_2 < Web3.toWei(1, "ether"), "预期结果解锁期后，钱已退+出块奖励+质押奖励"


@allure.title("自由账户质押+锁仓账户增持(都存在犹豫期+锁定期)")
@pytest.mark.P1
def test_009(staking_client):
    client = staking_client
    staking_address = client.staking_address
    node = client.node
    economic = client.economic
    log.info("创建锁仓计划")
    lockup_amount = economic.add_staking_limit * 5
    plan = [{'Epoch': 3, 'Amount': lockup_amount}]
    msg = client.restricting.createRestrictingPlan(staking_address, plan, economic.account.account_with_money["address"])
    assert_code(msg, 0), "创建锁仓计划失败"
    locked_info = client.ppos.getRestrictingInfo(staking_address)
    log.info(locked_info)

    msg = client.staking.increase_staking(1, staking_address)
    assert_code(msg, 0)
    log.info("进入第2个周期")
    economic.wait_settlement_blocknum(node)

    msg = client.staking.increase_staking(1, staking_address)
    assert_code(msg, 0)
    msg = client.staking.increase_staking(0, staking_address)
    assert_code(msg, 0)
    msg = client.ppos.getCandidateInfo(node.node_id)
    log.info("查询节点的质押情况{}".format(msg))

    assert msg["Ret"]["Shares"] == client.staking_amount + economic.add_staking_limit * 3
    assert msg["Ret"]["Released"] == client.staking_amount
    assert msg["Ret"]["RestrictingPlan"] == economic.add_staking_limit
    assert msg["Ret"]["RestrictingPlanHes"] == economic.add_staking_limit
    block_reward, staking_reward = economic.get_current_year_reward(node)

    log.info("节点2发起撤销质押")
    msg = client.staking.withdrew_staking(staking_address)
    assert_code(msg, 0)
    balance2 = node.eth.getBalance(staking_address)
    log.info("第2个周期发起撤销后的余额{}".format(balance2))
    """当前自由资金的增持已退,以下相减为手续费"""
    assert client.amount - balance2 < Web3.toWei(1, "ether")
    locked_info = client.ppos.getRestrictingInfo(staking_address)
    log.info("第2个周期发起撤销后查询锁仓计划{}".format(locked_info))
    assert_code(locked_info, 0)
    assert locked_info["Ret"]["Pledge"] == lockup_amount, "预期锁仓计划里的金额为锁定期金额"

    msg = client.ppos.getCandidateInfo(node.node_id)
    log.info("查询节点2的质押情况{}".format(msg))

    assert msg["Ret"]["ReleasedHes"] == 0, "预期犹豫期的金额已退"
    assert msg["Ret"]["RestrictingPlanHes"] == 0, "预期犹豫期的锁仓金额已退"

    log.info("进入第3个周期")
    economic.wait_settlement_blocknum(node)
    balance3 = node.eth.getBalance(staking_address)
    log.info("第3个周期发起撤销后的余额{}".format(balance3))

    log.info("进入第4个周期")
    economic.wait_settlement_blocknum(node, 1)
    balance4 = node.eth.getBalance(staking_address)
    log.info("第4个周期发起撤销后的余额{}".format(balance4))

    locked_info = client.ppos.getRestrictingInfo(staking_address)
    assert_code(locked_info, 304005)

    msg = client.ppos.getCandidateInfo(node.node_id)
    log.info("查询节点2的质押情况{}".format(msg))
    assert_code(msg, 301204)

    """算出块奖励+质押奖励"""
    log.info("以下为获取节点2出的块数")
    block_number = get_block_count_number(node, economic.settlement_size * 2)
    sum_block_reward = calculate(block_reward, block_number)
    reward_sum = sum_block_reward + staking_reward
    log.info("奖励的总金额{}".format(reward_sum))

    assert client.amount + reward_sum - balance4 < Web3.toWei(1, "ether"), "预期结果解锁期后，钱已退+出块奖励+质押奖励"


@allure.title("修改节点收益地址，再做退回：验证质押奖励+出块奖励")
@pytest.mark.P0
def test_alter_address_backup(staking_client):
    """
    修改钱包地址，更改后的地址收益正常
    """
    client = staking_client
    node = client.node
    staking_address = client.staking_address
    economic = client.economic
    ben_address, _ = economic.account.generate_account(node.web3)
    log.info("节点2修改节点信息")
    msg = client.staking.edit_candidate(staking_address, ben_address)
    assert_code(msg, 0)

    log.info("进入第2个结算周期")
    economic.wait_settlement_blocknum(node)

    block_reward, staking_reward = economic.get_current_year_reward(node)
    msg = client.staking.withdrew_staking(staking_address)
    assert_code(msg, 0)
    balance_before = node.eth.getBalance(ben_address)
    log.info("退出质押后新钱包余额：{}".format(balance_before))
    log.info("进入第3个结算周期")
    economic.wait_settlement_blocknum(node, 2)

    balance_after = node.eth.getBalance(ben_address)
    log.info("新钱包解锁期后的余额{}".format(balance_after))

    """算出块奖励+质押奖励"""
    log.info("以下为获取节点2出的块数")
    block_number = get_block_count_number(node, economic.settlement_size * 2)
    sum_block_reward = calculate(block_reward, block_number)
    reward_sum = sum_block_reward + staking_reward
    log.info("奖励的总金额{}".format(reward_sum))
    assert balance_after == reward_sum, "预期新钱包余额==收益奖励"
