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
