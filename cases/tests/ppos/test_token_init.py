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
