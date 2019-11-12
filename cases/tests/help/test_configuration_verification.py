import time
from copy import copy

import pytest
from dacite import from_dict

from common.log import log
from tests.lib.genesis import Genesis
from .utils import (
    assert_error_deploy,
    update_hesitate_ratio,
    update_staking,
    update_unstaking,
    update_unstake_freeze_duration,
)


@pytest.fixture()
def reset_cfg_env(global_test_env):
    cfg = global_test_env.cfg
    genesis = global_test_env.genesis_config
    backup_cfg = copy(cfg)
    id_cfg = id(cfg)
    yield global_test_env
    if id_cfg != id(global_test_env.cfg) or id(genesis) != id(global_test_env.genesis_config):
        global_test_env.set_cfg(backup_cfg)


@pytest.mark.P2
@pytest.mark.parametrize('value', [2, 8])
def test_IP_PR_001_to_012(value, reset_cfg_env):
    """
    IP_PR_001:校验结算周期是共识周期的倍数<4
    IP_PR_002:增发周期是结算周期的倍数<4
    """
    genesis = from_dict(data_class=Genesis, data=reset_cfg_env)
    genesis.economicModel.common.maxEpochMinutes = value
    new_file = reset_cfg_env.cfg.env_tmp + "/genesis.json"
    genesis.to_file(new_file)
    assert_error_deploy(reset_cfg_env, new_file, "Multiple of abnormal billing cycle")


@pytest.mark.P2
def test_IP_PR_003(reset_cfg_env):
    """
    IP_PR_003:备选验证人节点数小于验证节点数
    """
    genesis = from_dict(data_class=Genesis, data=reset_cfg_env)
    genesis.economicModel.staking.maxValidators = 3
    new_file = reset_cfg_env.cfg.env_tmp + "/genesis.json"
    genesis.to_file(new_file)
    msg = "The number of alternate verifier nodes is less than the number of verified nodes"
    assert_error_deploy(reset_cfg_env,new_file,msg)


@pytest.mark.P2
def test_IP_PR_004(reset_cfg_env):
    """
    正常范围内的质押Token数
    """
    node = reset_cfg_env.get_rand_node()
    value = node.web3.toWei(100,"ether")
    new_file = update_staking(reset_cfg_env, value)
    reset_cfg_env.deploy_all(new_file)
    config = node.debug.economicConfig()
    assert value == config["economicModel"]["staking"]["stakeThreshold"], "Inconsistent with the genesis file configuration amount"


@pytest.mark.P2
@pytest.mark.parametrize('value', [10000001, 9, 0])
def test_IP_PR_004_005(value, reset_cfg_env):
    """
    Abnormal pledge Token number
    1、创建验证人最低的质押Token数>10000000 * 10^18 von
    2、创建验证人最低的质押Token数<10 * 10^18 von
    3、创建验证人最低的质押Token数 = 0
    """
    node = reset_cfg_env.get_rand_node()
    value = node.web3.toWei(value, "ether")
    new_file = update_staking(reset_cfg_env, value)
    assert_error_deploy(reset_cfg_env, new_file, "Abnormal pledge Token number")


@pytest.mark.P2
def test_IP_PR_006_1(reset_cfg_env):
    node = reset_cfg_env.get_rand_node()
    value = node.web3.toWei(100, 'ether')
    new_file = update_unstaking(reset_cfg_env, value)
    reset_cfg_env.deploy_all(new_file)
    config = node.debug.economicConfig()
    assert value == config["economicModel"]["staking"]["operatingThreshold"]


@pytest.mark.P2
@pytest.mark.parametrize('value', [9, 0])
def test_IP_PR_006_2(value, reset_cfg_env):
    """
    修改每次委托及赎回的最低Token数
    1、委托人每次委托及赎回的最低Token数<10 * 10^18 von
    2、委托人每次委托及赎回的最低Token数 = 0
    """
    node = reset_cfg_env.get_rand_node()
    value = node.web3.toWei(value, 'ether')
    new_file = update_unstaking(reset_cfg_env, value)
    assert_error_deploy(reset_cfg_env, new_file, "The abnormal redemption amount")


@pytest.mark.P2
def test_IP_PR_007_1(reset_cfg_env):
    """
    正常范围内的犹豫期(多少个结算周期)
    """
    value = 3
    new_file = update_hesitate_ratio(reset_cfg_env, value)
    reset_cfg_env.deploy_all(new_file)
    config = reset_cfg_env.get_rand_node().debug.economicConfig()
    assert value == config["economicModel"]["staking"]["hesitateRatio"]


@pytest.mark.P2
@pytest.mark.parametrize('value', [0.1, 0])
def test_IP_PR_007(value, reset_cfg_env):
    """
    修改犹豫期(多少个结算周期)
    1、犹豫期(多少个结算周期)<=0
    2、犹豫期(多少个结算周期)=0
    """
    new_file = update_hesitate_ratio(reset_cfg_env, value)
    assert_error_deploy(reset_cfg_env,new_file,"An abnormal billing cycle")

@pytest.mark.P2
def test_IP_PR_008_1(reset_cfg_env):
    value = 3
    new_file = update_unstake_freeze_duration(reset_cfg_env, value)
    reset_cfg_env.deploy_all(new_file)
    config = reset_cfg_env.get_rand_node().debug.economicConfig()
    assert value == config["economicModel"]["staking"]["unStakeFreezeDuration"]

@pytest.mark.P2
@pytest.mark.parametrize('code', [1, 2, 3])
def test_IP_PR_008(self, code):
    """
    修改点质押退回锁定周期
    :param code:1、正常范围内的节点质押退回锁定周期
                2、节点质押退回锁定周期<0
                3、节点质押退回锁定周期=0
    :return:
    """
    if code == 1:
        start_init()
        # 修改ppos参数
        update_config('EconomicModel', 'Staking', 'UnStakeFreezeRatio', 3)
        try:
            # 启动节点
            auto = AutoDeployPlaton()
            auto.start_all_node(node_yml_path)
            time.sleep(3)
            result = get_economicConfig()
            result = get_json_parameter(result, 'EconomicModel', 'Staking', 'UnStakeFreezeRatio')
            assert result == 3, "创建验证人最低的质押：{} 有误".format(result)
        except Exception as e:
            log.info("异常信息：{} ".format(str(e)))

    elif code == 2:
        # 修改ppos参数
        update_config('EconomicModel', 'Staking', 'UnStakeFreezeRatio', 0.1)

        # 初始化链
        abnormal_start_up()

    elif code == 3:
        # 修改ppos参数
        update_config('EconomicModel', 'Staking', 'UnStakeFreezeRatio', 0)

        # 初始化链
        abnormal_start_up()

@pytest.mark.P2
@pytest.mark.parametrize('code', [1, 2, 3])
def test_IP_PR_009(self, code):
    """
    修改基金会分配年<=0
    :param code:1、正常范围内的基金会分配年
                2、基金会分配年<0
                3、基金会分配年=0
    :return:
    """
    if code == 1:
        start_init()
        # 修改ppos参数
        update_config('EconomicModel', 'Reward', 'PlatONFoundationYear', 3)

        try:
            # 启动节点
            auto = AutoDeployPlaton()
            auto.start_all_node(node_yml_path)
            time.sleep(3)
            result = get_economicConfig()
            result = get_json_parameter(result, 'EconomicModel', 'Reward', 'PlatONFoundationYear')
            assert result == 3, "创建验证人最低的质押：{} 有误".format(result)
        except Exception as e:
            log.info("异常信息：{} ".format(str(e)))

    elif code == 2:
        # 修改ppos参数
        update_config('EconomicModel', 'Reward', 'PlatONFoundationYear', 0.1)

        # 初始化链
        abnormal_start_up()

    elif code == 3:
        # 修改ppos参数
        update_config('EconomicModel', 'Reward', 'PlatONFoundationYear', 0)

        # 初始化链
        abnormal_start_up()

@pytest.mark.P2
@pytest.mark.parametrize('code', [1, 2, 3])
def test_IP_PR_011_013(self, code):
    """
    正常范围内的奖励池分配给出块奖励的比例
    IP_PR_011:奖励池分配给出块奖励的比例=0
    IP_PR_013：奖励池分配给出块奖励的比例=100
    :param code:
    :return:
    """
    start_init()

    if code == 1:
        # 修改ppos参数
        update_config('EconomicModel', 'Reward', 'NewBlockRate', 40)

        # 启动节点
        auto = AutoDeployPlaton()
        auto.start_all_node(node_yml_path)
        time.sleep(3)
        result = get_economicConfig()
        result = get_json_parameter(result, 'Reward', 'NewBlockRate',None)
        log.info("奖励池分配给出块奖励: {}".format(result))
        assert result == 40, "奖励池分配给出块奖励：{} 有误".format(result)

    elif code == 2:
        # 修改ppos参数
        update_config('EconomicModel', 'Reward', 'NewBlockRate', 0)

        # 启动节点
        auto = AutoDeployPlaton()
        auto.start_all_node(node_yml_path)
        time.sleep(3)
        result = get_economicConfig()
        result = get_json_parameter(result, 'Reward', 'NewBlockRate',None)
        log.info("奖励池分配给出块奖励: {}".format(result))
        assert result == 0, "奖励池分配给出块奖励：{} 有误".format(result)

    elif code == 3:
        # 修改ppos参数
        update_config('EconomicModel', 'Reward', 'NewBlockRate', 100)

        # 启动节点
        auto = AutoDeployPlaton()
        auto.start_all_node(node_yml_path)
        time.sleep(3)
        result = get_economicConfig()
        result = get_json_parameter(result, 'Reward', 'NewBlockRate',None)
        log.info("奖励池分配给出块奖励: {}".format(result))
        assert result == 100, "奖励池分配给出块奖励：{} 有误".format(result)

@pytest.mark.P2
@pytest.mark.parametrize('code', [1, 2])
def test_IP_PR_010_012(self, code):
    """
    IP_PR_010:奖励池分配给出块奖励的比例<0
    IP_PR_012:奖励池分配给出块奖励的比例>100
    :param code:
    :return:
    """
    if code == 1:
        # 修改ppos参数
        update_config('EconomicModel', 'Reward', 'NewBlockRate', '-1')

        # 初始化链
        abnormal_start_up()

    elif code == 2:
        # 修改ppos参数
        update_config('EconomicModel', 'Reward', 'NewBlockRate', 110)

        # 初始化链
        abnormal_start_up()

@pytest.mark.P2
def test_IP_CP_001(self):
    """
    创世文件链参数验证
    :return:
    """
    # 修改eip155Block参数字符串
    update_config('config', 'eip155Block', None, 'ss')
    # 初始化链
    abnormal_start_up()

    # 修改eip155Block参数空值
    update_config('config', 'eip155Block', None, None)
    # 初始化链
    abnormal_start_up()


@pytest.mark.P2
def test_IP_CP_002(self):
    """
    创世文件共识参数验证
    :return:
    """
    # 修改cbft-amount参数字符串
    update_config('config', 'cbft', 'amount', 'ss')
    # 初始化链
    abnormal_start_up()

    # 修改cbft-amount参数空值
    update_config('config', 'cbft', 'amount', None)
    # 初始化链
    abnormal_start_up()

    # 修改cbft-epoch参数非正整数
    update_config('config', 'cbft', 'epoch', 0.1)
    # 初始化链
    abnormal_start_up()

    # 修改cbft-epoch参数字符串
    update_config('config', 'cbft', 'epoch', 'ss')
    # 初始化链
    abnormal_start_up()

    # 修改cbft-epoch参数空值
    update_config('config', 'cbft', 'epoch', None)
    # 初始化链
    abnormal_start_up()

    # 修改cbft-validatorMode参数字符串
    update_config('config', 'cbft', 'validatorMode', 'ss')
    # 初始化链
    abnormal_start_up()

    # 修改cbft-validatorMode参数空值
    update_config('config', 'cbft', 'validatorMode', None)
    # 初始化链
    abnormal_start_up()

    # 修改cbft-period参数空值
    update_config('config', 'cbft', 'period', None)
    # 初始化链
    abnormal_start_up()

@pytest.mark.P2
def test_IP_CP_003(self):
    """
    创世文件经济参数验证
    :return:
    """
    # 修改Common-AdditionalCycleTime参数空值
    update_config('EconomicModel', 'Common', 'AdditionalCycleTime', None)
    # 初始化链
    abnormal_start_up()

    # 修改Staking-PassiveUnDelegateFreezeRatio参数空值
    update_config('EconomicModel', 'Staking', 'PassiveUnDelegateFreezeRatio', None)
    # 初始化链
    abnormal_start_up()

    # 修改Gov-VersionProposalActive_ConsensusRounds参数空值
    update_config('EconomicModel', 'Gov', 'VersionProposalActive_ConsensusRounds', None)
    # 初始化链
    abnormal_start_up()

    # 修改Reward-NewBlockRate参数空值
    update_config('EconomicModel', 'Reward', 'NewBlockRate', None)
    # 初始化链
    abnormal_start_up()

    # 修改InnerAcc-CDFAccount参数空值
    update_config('EconomicModel', 'InnerAcc', 'CDFAccount', None)
    # 初始化链
    abnormal_start_up()

@pytest.mark.P2
def test_IP_CP_004(self):
    """
    创世文件处罚参数验证
    :return:
    """
    # 修改Slashing-PackAmountAbnormal参数空值
    update_config('EconomicModel', 'Slashing', 'PackAmountAbnormal', None)
    # 初始化链
    abnormal_start_up()

    # 修改Slashing-DuplicateSignHighSlashing参数空值
    update_config('EconomicModel', 'Slashing', 'DuplicateSignHighSlashing', None)
    # 初始化链
    abnormal_start_up()

    # 修改Slashing-NumberOfBlockRewardForSlashing参数空值
    update_config('EconomicModel', 'Slashing', 'NumberOfBlockRewardForSlashing', None)
    # 初始化链
    abnormal_start_up()

    # 修改Slashing-EvidenceValidEpoch参数空值
    update_config('EconomicModel', 'Slashing', 'EvidenceValidEpoch', None)
    # 初始化链
    abnormal_start_up()

    # 修改Slashing-PackAmountHighAbnormal参数空值
    update_config('EconomicModel', 'Slashing', 'PackAmountHighAbnormal', None)
    # 初始化链
    abnormal_start_up()

@pytest.mark.P2
def test_IP_CP_005(self):
    """
    创世文件参数chainid验证
    :return:
    """

    # 修改Slashing-PackAmountHighAbnormal参数非正整数
    update_config('config', 'chainId', None, 0.1)
    # 初始化链
    abnormal_start_up()

    # 修改Slashing-PackAmountHighAbnormal参数字符串
    update_config('config', 'chainId', None, 'ss')
    # 初始化链
    abnormal_start_up()

    # 修改Slashing-PackAmountHighAbnormal参数空值
    update_config('config', 'chainId', None, None)
    # 初始化链
    abnormal_start_up()


if __name__ == '__main__':
    a = TestConfiguration()
    a.test_IP_PR_011_013(3)
