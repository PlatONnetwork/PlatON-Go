import time

import pytest
import allure

from dacite import from_dict
from common.log import log
from client_sdk_python import Web3
from decimal import Decimal
from tests.lib import EconomicConfig, Genesis, StakingConfig, Staking, check_node_in_list


@pytest.fixture(scope="function")
def staking_obj(global_test_env):
    cfg = StakingConfig("11111", "faker", "www.baidu.com", "how much")
    node = global_test_env.get_rand_node()
    return Staking(global_test_env, node, cfg)


@pytest.fixture(scope="class")
def staking_candidate(client_consensus_obj):
    address, _ = client_consensus_obj.economic.account.generate_account(client_consensus_obj.node.web3,
                                                                        client_consensus_obj.economic.create_staking_limit * 2)
    result = client_consensus_obj.staking.create_staking(0, address, address)
    assert result['Code'] == 0, "申请质押返回的状态：{}, {}用例失败".format(result['Code'], result['ErrMsg'])
    # 等待锁定期
    client_consensus_obj.economic.wait_settlement_blocknum(client_consensus_obj.node)
    return client_consensus_obj, address


@allure.title("链初始化各账号分配金额验证")
@pytest.mark.P0
def test_IT_IA_002_to_007(new_env):
    """
    IT_IA_002:链初始化-查看token发行总量账户初始值
    IT_IA_003:链初始化-查看platON基金会账户初始值
    IT_IA_004:链初始化-查看激励池账户
    IT_IA_005:链初始化-查看剩余总账户
    IT_IA_006:链初始化-查看锁仓账户余额
    IT_IA_007:链初始化-查看质押账户余额
    :return:验证链初始化后token各内置账户初始值
    """
    # 初始化genesis文件初始金额
    node_count = len(new_env.consensus_node_list)
    default_pledge_amount = Web3.toWei(node_count * 1500000, 'ether')
    node = new_env.get_rand_node()
    community_amount = default_pledge_amount + 259096239000000000000000000 + 62215742000000000000000000
    genesis = from_dict(data_class=Genesis, data=new_env.genesis_config)
    genesis.EconomicModel.InnerAcc.CDFBalance = community_amount
    surplus_amount = str(EconomicConfig.TOKEN_TOTAL - community_amount - 200000000000000000000000000)
    genesis.alloc = {
        "1000000000000000000000000000000000000003": {
            "balance": "200000000000000000000000000"
        },
        "0x2e95e3ce0a54951eb9a99152a6d5827872dfb4fd": {
            "balance": surplus_amount
        }
    }
    new_file = new_env.cfg.env_tmp + "/genesis.json"
    genesis.to_file(new_file)
    new_env.deploy_all(new_file)

    # 验证各个内置账户金额
    foundation_louckup = node.eth.getBalance(Web3.toChecksumAddress(EconomicConfig.FOUNDATION_LOCKUP_ADDRESS))
    log.info('初始锁仓合约地址： {} 金额：{}'.format(EconomicConfig.FOUNDATION_LOCKUP_ADDRESS, foundation_louckup))
    incentive_pool = node.eth.getBalance(Web3.toChecksumAddress(EconomicConfig.INCENTIVEPOOL_ADDRESS))
    log.info('激励池地址：{} 查询金额：{}'.format(EconomicConfig.INCENTIVEPOOL_ADDRESS, incentive_pool))
    staking = node.eth.getBalance(Web3.toChecksumAddress(EconomicConfig.STAKING_ADDRESS))
    log.info('STAKING地址：{} 查询金额：{}'.format(EconomicConfig.STAKING_ADDRESS, staking))
    foundation = node.eth.getBalance(Web3.toChecksumAddress(EconomicConfig.FOUNDATION_ADDRESS))
    log.info('FOUNDATION地址：{}查询金额：{}'.format(EconomicConfig.FOUNDATION_ADDRESS, foundation))
    remain = node.eth.getBalance(Web3.toChecksumAddress(EconomicConfig.REMAIN_ACCOUNT_ADDRESS))
    log.info('REMAINACCOUNT地址：{} 查询金额：{}'.format(EconomicConfig.REMAIN_ACCOUNT_ADDRESS, remain))
    develop = node.eth.getBalance(Web3.toChecksumAddress(EconomicConfig.DEVELOPER_FOUNDATAION_ADDRESS))
    log.info('COMMUNITYDEVELOPER地址：{} 查询金额：{}'.format(EconomicConfig.DEVELOPER_FOUNDATAION_ADDRESS, develop))
    reality_total = foundation_louckup + incentive_pool + staking + foundation + remain + develop
    log.info("创世区块发行总量：{}".format(reality_total))
    log.info("--------------分割线---------------")
    assert foundation == 0, "基金会初始金额:{}有误".format(foundation)
    assert foundation_louckup == 259096239000000000000000000, "基金会锁仓初始金额:{}有误".format(foundation_louckup)
    assert staking == default_pledge_amount, "初始质押账户金额:{}有误".format(staking)
    assert incentive_pool == 262215742000000000000000000, "奖励池初始金额:{}有误".format(incentive_pool)
    assert remain == int(surplus_amount), "剩余总账户初始金额:{}有误".format(remain)
    assert develop == 0, "社区开发者基金会账户金额：{} 有误".format(develop)
    assert reality_total == EconomicConfig.TOKEN_TOTAL, "初始化发行值{}有误".format(reality_total)

    
@allure.title("二次分配：转账金额：{value}")
@pytest.mark.P0
@pytest.mark.parametrize('value', [1000, 0.000000000000000001, 100000000])
def test_IT_SD_004_to_006(client_consensus_obj, value):
    """
    IT_SD_006:二次分配：普通钱包转keyshard钱包
    IT_SD_004:二次分配：转账金额为1von
    IT_SD_005:二次分配：转账金额为1亿LAT
    :param client_consensus_obj:
    :param value:
    :return:
    """
    balance = client_consensus_obj.node.eth.getBalance(
        client_consensus_obj.node.web3.toChecksumAddress(0x493301712671Ada506ba6Ca7891F436D29185821))
    value = client_consensus_obj.node.web3.toWei(value, 'ether')
    address, _ = client_consensus_obj.economic.account.generate_account(client_consensus_obj.node.web3, value)
    balance = client_consensus_obj.node.eth.getBalance(client_consensus_obj.node.web3.toChecksumAddress(address))
    log.info("交易之后账户：{}的余额：{}".format(address, balance))
    assert balance == value, "转账金额:{}失败".format(balance)

@pytest.mark.P1
@pytest.mark.parametrize('code', [1, 2, 3])
def test_IT_SD_002_003_011(global_test_env, code):
    """
    IT_SD_002:二次分配：账户余额不足
    IT_SD_003:二次分配：转账手续费不足
    IT_SD_011:账户转账校验：转账gas费不足
    :param global_test_env:
    :param code:
    :return:
    """
    node = global_test_env.get_rand_node()
    value = node.web3.toWei(1000, 'ether')
    address, _ = global_test_env.account.generate_account(node.web3, value)
    if code == 1:
        # 账户余额不足转账
        try:
            address1, _ = global_test_env.account.generate_account(node.web3, 0)
            result = global_test_env.account.sendTransaction(node.web3, '', node.web3.toChecksumAddress(address),
                                                             node.web3.toChecksumAddress(address1),
                                                             node.web3.platon.gasPrice, 21000, 2000)
            return_info = node.eth.waitForTransactionReceipt(result)
            assert return_info is not None, "用例失败"
        except Exception as e:
            log.info("用例成功，异常信息：{} ".format(str(e)))
    elif code == 2:
        # 转账手续费不足
        try:
            address1, _ = global_test_env.account.generate_account(node.web3, 0)
            result = global_test_env.account.sendTransaction(node.web3, '', node.web3.toChecksumAddress(address),
                                                             node.web3.toChecksumAddress(address1),
                                                             node.web3.platon.gasPrice, 21000, 1000)
            return_info = node.eth.waitForTransactionReceipt(result)
            assert return_info is not None, "用例失败"
        except Exception as e:
            log.info("用例成功，异常信息：{} ".format(str(e)))
    elif code == 3:
        # 转账gas费不足
        try:
            address1, _ = global_test_env.account.generate_account(node.web3, 0)
            result = global_test_env.account.sendTransaction(node.web3, '', node.web3.toChecksumAddress(address),
                                                             node.web3.toChecksumAddress(address1),
                                                             node.web3.platon.gasPrice, 2100, 500)
            return_info = node.eth.waitForTransactionReceipt(result)
            assert return_info is not None, "用例失败"
        except Exception as e:
            log.info("用例成功，异常信息：{} ".format(str(e)))

@pytest.mark.P2
def test_IT_SD_007(global_test_env):
    """
    账户转账校验：本账户转本账户
    :return:
    """
    node = global_test_env.get_rand_node()
    value = node.web3.toWei(1000, 'ether')
    address, _ = global_test_env.account.generate_account(node.web3, value)
    balance = node.eth.getBalance(node.web3.toChecksumAddress(address))
    log.info("转账之前账户余额： {}".format(balance))
    result = global_test_env.account.sendTransaction(node.web3, '', address, address, node.eth.gasPrice, 21000, 100)
    assert result is not None, "用例失败"
    balance1 = node.eth.getBalance(node.web3.toChecksumAddress(address))
    log.info("转账之后账户余额： {}".format(balance1))
    log.info("手续费： {}".format(node.web3.platon.gasPrice * 21000))
    assert balance == balance1 + node.web3.platon.gasPrice * 21000, "转账之后账户余额： {} 有误".format(balance1)

    
@pytest.mark.P0
def test_IT_SD_008(global_test_env):
    """
    二次分配：普通账户转platON基金会账户
    :return:
    """
    node = global_test_env.get_rand_node()
    value = node.web3.toWei(1000, 'ether')
    address, _ = global_test_env.account.generate_account(node.web3, value)
    balance = node.eth.getBalance(EconomicConfig.INCENTIVEPOOL_ADDRESS)
    result = global_test_env.account.sendTransaction(node.web3, '', address, EconomicConfig.INCENTIVEPOOL_ADDRESS,
                                                     node.eth.gasPrice, 21000, node.web3.toWei(100, 'ether'))
    assert result is not None, "用例失败"
    balance1 = node.eth.getBalance(EconomicConfig.INCENTIVEPOOL_ADDRESS)
    log.info("转账之后账户余额： {}".format(balance1))
    log.info("手续费： {}".format(node.web3.platon.gasPrice * 21000))
    assert balance1 == balance + node.web3.toWei(100,'ether') + node.web3.platon.gasPrice * 21000, "转账之后账户余额： {} 有误".format(
        balance1)
    
def consensus_node_pledge_award_assertion(client_new_node_obj, address):
    """
    内置节点质押奖励断言
    :param client_new_node_obj:
    :param address:
    :return:
    """
    blockNumber = client_new_node_obj.node.eth.blockNumber
    log.info("当前块高：{}".format(blockNumber))
    incentive_pool_balance = client_new_node_obj.node.eth.getBalance(EconomicConfig.INCENTIVEPOOL_ADDRESS)
    log.info("激励池余额：{}".format(incentive_pool_balance))
    CandidateInfo = client_new_node_obj.ppos.getCandidateInfo(client_new_node_obj.node.node_id)
    log.info("质押人节点信息：{}".format(CandidateInfo))

    # 等待质押节点到锁定期
    client_new_node_obj.economic.wait_settlement_blocknum(client_new_node_obj.node)
    time.sleep(5)
    VerifierList = client_new_node_obj.ppos.getVerifierList()
    log.info("当前验证人列表：{}".format(VerifierList))
    ValidatorList = client_new_node_obj.ppos.getValidatorList()
    log.info("当前共识验证人列表：{}".format(ValidatorList))
    # 申请退回质押
    result = client_new_node_obj.staking.withdrew_staking(address)
    log.info("退回质押结果: {}".format(result))
    assert result['Code'] == 0, "申请退回质押返回的状态：{}, {}".format(result['Code'], result['ErrMsg'])
    # 等待当前结算结束
    client_new_node_obj.economic.wait_settlement_blocknum(client_new_node_obj.node)
    # 查看激励池金额
    incentive_pool_balance2 = client_new_node_obj.node.eth.getBalance(EconomicConfig.INCENTIVEPOOL_ADDRESS)
    log.info("地址：{}, 金额：{}".format(EconomicConfig.INCENTIVEPOOL_ADDRESS, incentive_pool_balance2))

    assert incentive_pool_balance2 - incentive_pool_balance < client_new_node_obj.node.web3.toWei(1, 'ether'), "激励池余额：{} 有误".format(
        incentive_pool_balance2)
 

def no_consensus_node_pledge_award_assertion(client_new_node_obj, benifit_address, from_address):
    """
    非内置节点质押奖励断言
    :param client_new_node_obj:
    :param benifit_address: 节点收益地址
    :param from_address: 节点质押地址
    :return:
    """
    CandidateInfo = client_new_node_obj.ppos.getCandidateInfo(client_new_node_obj.node.node_id)
    log.info("质押人节点信息：{}".format(CandidateInfo))
    balance = client_new_node_obj.node.eth.getBalance(client_new_node_obj.node.web3.toChecksumAddress(benifit_address))
    log.info("地址：{} 的金额： {}".format(benifit_address, balance))

    # 等待质押节点到锁定期
    client_new_node_obj.economic.wait_settlement_blocknum(client_new_node_obj.node)
    time.sleep(5)
    VerifierList = client_new_node_obj.ppos.getVerifierList()
    log.info("当前验证人列表：{}".format(VerifierList))
    ValidatorList = client_new_node_obj.ppos.getValidatorList()
    log.info("当前共识验证人列表：{}".format(ValidatorList))
    block_reward, staking_reward = client_new_node_obj.economic.get_current_year_reward(client_new_node_obj.node)
    for i in range(4):
        result = check_node_in_list(client_new_node_obj.node.node_id, client_new_node_obj.ppos.getValidatorList)
        log.info("当前节点是否在共识列表：{}".format(result))
        if result:
            # 等待一个共识轮
            client_new_node_obj.economic.wait_consensus_blocknum(client_new_node_obj.node)
            # 申请退回质押
            result = client_new_node_obj.staking.withdrew_staking(from_address)
            log.info("退回质押结果: {}".format(result))
            assert result['Code'] == 0, "申请退回质押返回的状态：{}, {}".format(result['Code'], result['ErrMsg'])
            incentive_pool_balance1 = client_new_node_obj.node.eth.getBalance(EconomicConfig.INCENTIVEPOOL_ADDRESS)
            log.info("激励池余额：{}".format(incentive_pool_balance1))
            # 等待当前结算结束
            client_new_node_obj.economic.wait_settlement_blocknum(client_new_node_obj.node)
            # 统计质押节点出块数
            blocknumber = client_new_node_obj.economic.get_block_count_number(client_new_node_obj.node, 5)
            balance1 = client_new_node_obj.node.eth.getBalance(
                client_new_node_obj.node.web3.toChecksumAddress(benifit_address))
            log.info("地址：{} 的金额： {}".format(benifit_address, balance1))

            # 验证出块奖励
            log.info("预计出块奖励：{}".format(Decimal(str(block_reward)) * Decimal(blocknumber)))
            assert balance + (Decimal(str(block_reward)) * Decimal(
                blocknumber)) - balance1 < client_new_node_obj.node.web3.toWei(1, 'ether'), "地址: {} 的金额: {} 有误".format(
                benifit_address, balance1)
            break
        else:
            # 等一个共识轮切换共识验证人
            client_new_node_obj.economic.wait_consensus_blocknum(client_new_node_obj.node)



@pytest.mark.P1
def test_AL_IE_003(client_new_node_obj_list):
    """
    自由账户创建质押节点且收益地址为激励池
    :param client_new_node_obj_list:
    :return:
    """
    log.info("节点id：{}".format(client_new_node_obj_list[0].node.node_id))
    address, _ = client_new_node_obj_list[0].economic.account.generate_account(client_new_node_obj_list[0].node.web3,
                                                                               client_new_node_obj_list[
                                                                                   0].economic.create_staking_limit * 2)
    log.info("质押账户地址: {}".format(address))
    # 自由金额申请质押节点
    result = client_new_node_obj_list[0].staking.create_staking(0, EconomicConfig.INCENTIVEPOOL_ADDRESS, address)
    log.info("质押结果: {}".format(result))
    assert result['Code'] == 0, "申请质押返回的状态：{}, {}".format(result['Code'], result['ErrMsg'])
    consensus_node_pledge_award_assertion(client_new_node_obj_list[0], address)


@pytest.mark.P1
def test_AL_IE_004(client_new_node_obj_list):
    """
    锁仓账户创建质押节点且收益地址为激励池
    :param client_new_node_obj_list:
    :return:
    """
    log.info("节点id：{}".format(client_new_node_obj_list[1].node.node_id))
    address, _ = client_new_node_obj_list[1].economic.account.generate_account(client_new_node_obj_list[1].node.web3,
                                                                               client_new_node_obj_list[
                                                                                   1].economic.create_staking_limit * 2)
    log.info("质押账户地址: {}".format(address))
    # 创建锁仓计划
    staking_amount = client_new_node_obj_list[1].economic.create_staking_limit
    log.info("质押金额：{}".format(staking_amount))
    plan = [{'Epoch': 1, 'Amount': staking_amount}]
    result = client_new_node_obj_list[1].restricting.createRestrictingPlan(address, plan, address)
    assert result['Code'] == 0, "创建锁仓计划返回的状态：{}, {}".format(result['Code'], result['ErrMsg'])
    # 锁仓金额申请质押节点
    result = client_new_node_obj_list[1].staking.create_staking(1, EconomicConfig.INCENTIVEPOOL_ADDRESS, address)
    log.info("质押结果: {}".format(result))
    assert result['Code'] == 0, "申请质押返回的状态：{}, {}".format(result['Code'], result['ErrMsg'])
    consensus_node_pledge_award_assertion(client_new_node_obj_list[1], address)


@pytest.mark.P1
def test_AL_BI_004(client_consensus_obj):
    """
    初始验证人退出后重新质押进来
    :param client_consensus_obj:
    :return:
    """
    developer_foundation_balance = client_consensus_obj.node.eth.getBalance(
        client_consensus_obj.node.web3.toChecksumAddress(EconomicConfig.DEVELOPER_FOUNDATAION_ADDRESS))
    log.info("incentive_pool_balance: {}".format(developer_foundation_balance))
    staking_balance = client_consensus_obj.node.eth.getBalance(EconomicConfig.STAKING_ADDRESS)
    log.info("staking_balance: {}".format(staking_balance))

    # 内置节点退回质押
    result = client_consensus_obj.staking.withdrew_staking(
        client_consensus_obj.node.web3.toChecksumAddress(EconomicConfig.DEVELOPER_FOUNDATAION_ADDRESS))
    assert result['Code'] == 0, "申请退回质押返回的状态：{}, {}".format(result['Code'], result['ErrMsg'])
    # create account
    address, _ = client_consensus_obj.economic.account.generate_account(client_consensus_obj.node.web3,
                                                                        client_consensus_obj.economic.create_staking_limit * 2)

    # 等待退出质押
    client_consensus_obj.economic.wait_settlement_blocknum(client_consensus_obj.node, 1)
    # 查看账户金额变化
    developer_foundation_balance1 = client_consensus_obj.node.eth.getBalance(
        client_consensus_obj.node.web3.toChecksumAddress(EconomicConfig.DEVELOPER_FOUNDATAION_ADDRESS))
    log.info("incentive_pool_balance: {}".format(developer_foundation_balance1))
    staking_balance1 = client_consensus_obj.node.eth.getBalance(EconomicConfig.STAKING_ADDRESS)
    log.info("staking_balance: {}".format(staking_balance1))
    # 断言账户金额变化
    assert developer_foundation_balance + EconomicConfig.DEVELOPER_STAKING_AMOUNT - developer_foundation_balance1 < client_consensus_obj.node.web3.toWei(
        1, 'ether'), "error: developer_foundation_balance1: {}".format(developer_foundation_balance1)
    assert staking_balance1 == staking_balance - EconomicConfig.DEVELOPER_STAKING_AMOUNT, "error: staking_balance1: {}".format(
        staking_balance1)
    # 节点重新质押
    result = client_consensus_obj.staking.create_staking(0, address, address)
    assert result['Code'] == 0, "申请质押返回的状态：{}, {}".format(result['Code'], result['ErrMsg'])
    # 质押奖励和出块奖励断言
    no_consensus_node_pledge_award_assertion(client_consensus_obj, address, address)
