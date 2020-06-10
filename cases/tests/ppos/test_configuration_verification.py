import time
import json
from copy import copy

import pytest
from dacite import from_dict
from client_sdk_python import Web3
from tests.lib.genesis import Genesis
from environment.node import Node


def assert_error_deploy(env: Node, genesis_file, msg="Error config"):
    env.clean_db()
    is_success, err_msg = env.deploy_me(genesis_file)
    print(is_success)
    assert not is_success, "{}-{}".format(msg, err_msg)


@pytest.fixture(scope="module", autouse=True)
def stop(global_test_env):
    global_test_env.stop_all()
    yield
    # global_test_env.deploy_all()


@pytest.fixture()
def reset_cfg_env_node(global_test_env):
    new_env_node = copy(global_test_env.get_rand_node())
    cfg = global_test_env.cfg
    genesis_config = global_test_env.genesis_config
    backup_cfg = copy(cfg)
    genesis = from_dict(data_class=Genesis, data=genesis_config)
    id_cfg = id(cfg)
    setattr(new_env_node, "genesis", genesis)
    setattr(new_env_node, "genesis_path", global_test_env.cfg.env_tmp + "/genesis_0.13.0.json")
    yield new_env_node
    new_env_node.stop()
    if id_cfg != id(global_test_env.cfg) or id(genesis_config) != id(global_test_env.genesis_config):
        global_test_env.set_cfg(backup_cfg)


@pytest.mark.P2
@pytest.mark.parametrize('value', [2, 8, ""])
def test_IP_PR_001_to_012(value, reset_cfg_env_node):
    """
    IP_PR_001:校验结算周期是共识周期的倍数<4
    IP_PR_002:增发周期是结算周期的倍数<4
    """
    genesis = reset_cfg_env_node.genesis
    genesis.economicModel.common.maxEpochMinutes = value
    genesis.to_file(reset_cfg_env_node.genesis_path)
    assert_error_deploy(reset_cfg_env_node, reset_cfg_env_node.genesis_path, "Multiple of abnormal billing cycle")


@pytest.mark.P2
def test_IP_PR_003(reset_cfg_env_node):
    """
    IP_PR_003:备选验证人节点数小于验证节点数
    """
    genesis = reset_cfg_env_node.genesis
    genesis.economicModel.staking.maxValidators = 3
    genesis.to_file(reset_cfg_env_node.genesis_path)
    msg = "The number of alternate verifier nodes is less than the number of verified nodes"
    assert_error_deploy(reset_cfg_env_node, reset_cfg_env_node.genesis_path, msg)


@pytest.mark.P2
def test_IP_PR_004(reset_cfg_env_node):
    """
    正常范围内的质押Token数
    """
    value = Web3.toWei(1000000, "ether")
    genesis = reset_cfg_env_node.genesis
    genesis.economicModel.staking.stakeThreshold = value
    genesis.to_file(reset_cfg_env_node.genesis_path)
    reset_cfg_env_node.deploy_me(reset_cfg_env_node.genesis_path)
    config = reset_cfg_env_node.debug.economicConfig()
    assert value == config["staking"]["stakeThreshold"], "Inconsistent with the genesis file configuration amount"


@pytest.mark.P2
@pytest.mark.parametrize('value', [9, 0])
def test_IP_PR_004_005(value, reset_cfg_env_node):
    """
    Abnormal pledge Token number
    1、创建验证人最低的质押Token数<10 * 10^18 von
    2、创建验证人最低的质押Token数 = 0
    """
    value = Web3.toWei(value, "ether")
    genesis = reset_cfg_env_node.genesis
    genesis.economicModel.staking.stakeThreshold = value
    new_file = reset_cfg_env_node.genesis_path
    genesis.to_file(new_file)
    assert_error_deploy(reset_cfg_env_node, new_file, "Abnormal pledge Token number")


@pytest.mark.P2
def test_IP_PR_006_1(reset_cfg_env_node):
    value = Web3.toWei(100, 'ether')
    genesis = reset_cfg_env_node.genesis
    genesis.economicModel.staking.operatingThreshold = value
    new_file = reset_cfg_env_node.genesis_path
    genesis.to_file(new_file)
    reset_cfg_env_node.deploy_me(new_file)
    config = reset_cfg_env_node.debug.economicConfig()
    assert value == config["staking"]["operatingThreshold"]


@pytest.mark.P2
@pytest.mark.parametrize('value', [9, 0])
def test_IP_PR_006_2(value, reset_cfg_env_node):
    """
    修改每次委托及赎回的最低Token数
    1、委托人每次委托及赎回的最低Token数<10 * 10^18 von
    2、委托人每次委托及赎回的最低Token数 = 0
    """
    value = Web3.toWei(value, 'ether')
    genesis = reset_cfg_env_node.genesis
    genesis.economicModel.staking.operatingThreshold = value
    new_file = reset_cfg_env_node.genesis_path
    genesis.to_file(new_file)
    assert_error_deploy(reset_cfg_env_node, new_file, "The abnormal redemption amount")


# @pytest.mark.P2
# def test_IP_PR_007_1(reset_cfg_env_node):
#     """
#     正常范围内的犹豫期(多少个结算周期)
#     """
#     value = 3
#     genesis = reset_cfg_env_node.genesis
#     genesis.economicModel.staking.hesitateRatio = value
#     new_file = reset_cfg_env_node.genesis_path
#     genesis.to_file(new_file)
#     reset_cfg_env_node.deploy_me(new_file)
#     config = reset_cfg_env_node.debug.economicConfig()
#     assert value == config["staking"]["hesitateRatio"]


# @pytest.mark.P2
# @pytest.mark.parametrize('value', [-1, 0, ""])
# def test_IP_PR_007(value, reset_cfg_env_node):
#     """
#     修改犹豫期(多少个结算周期)
#     1、犹豫期(多少个结算周期)<=0
#     2、犹豫期(多少个结算周期)=0
#     """
#     genesis = reset_cfg_env_node.genesis
#     genesis.economicModel.staking.hesitateRatio = value
#     new_file = reset_cfg_env_node.genesis_path
#     genesis.to_file(new_file)
#     assert_error_deploy(reset_cfg_env_node, new_file, "An abnormal billing cycle")


@pytest.mark.P2
def test_IP_PR_008_1(reset_cfg_env_node):
    """
    正常范围内的节点质押退回锁定周期
    """
    value = 3
    genesis = reset_cfg_env_node.genesis
    genesis.economicModel.staking.unStakeFreezeDuration = value
    new_file = reset_cfg_env_node.genesis_path
    genesis.to_file(new_file)
    reset_cfg_env_node.deploy_me(new_file)
    config = reset_cfg_env_node.debug.economicConfig()
    assert value == config["staking"]["unStakeFreezeDuration"]


@pytest.mark.P2
@pytest.mark.parametrize('value', [-1, 0, ""])
def test_IP_PR_008_2(value, reset_cfg_env_node):
    """
    修改点质押退回锁定周期
    1、节点质押退回锁定周期<0
    2、节点质押退回锁定周期=0
    """
    genesis = reset_cfg_env_node.genesis
    genesis.economicModel.staking.unStakeFreezeDuration = value
    new_file = reset_cfg_env_node.genesis_path
    genesis.to_file(new_file)
    assert_error_deploy(reset_cfg_env_node, new_file, "Abnormal pledge return lock cycle")


@pytest.mark.P2
def test_IP_PR_009_1(reset_cfg_env_node):
    """
    正常范围内的基金会分配年
    """
    value = 3
    genesis = reset_cfg_env_node.genesis
    genesis.economicModel.reward.platONFoundationYear = value
    new_file = reset_cfg_env_node.genesis_path
    genesis.to_file(new_file)
    reset_cfg_env_node.deploy_me(new_file)
    config = reset_cfg_env_node.debug.economicConfig()
    assert value == config["reward"]["platonFoundationYear"]


@pytest.mark.P2
@pytest.mark.parametrize('value', [-1, 0, ""])
def test_IP_PR_009(value, reset_cfg_env_node):
    """
    修改基金会分配年
    1、基金会分配年<0
    2、基金会分配年=0
    """
    genesis = reset_cfg_env_node.genesis
    genesis.economicModel.reward.platONFoundationYear = value
    new_file = reset_cfg_env_node.genesis_path
    genesis.to_file(new_file)
    assert_error_deploy(reset_cfg_env_node, new_file, "Abnormal fund allocation year")


@pytest.mark.P2
@pytest.mark.parametrize('value', [40, 0, 100])
def test_IP_PR_011(reset_cfg_env_node, value):
    """
    正常范围内的奖励池分配给出块奖励的比例
    IP_PR_011:奖励池分配给出块奖励的比例=0
    IP_PR_013：奖励池分配给出块奖励的比例=100
    """
    genesis = reset_cfg_env_node.genesis
    genesis.economicModel.reward.newBlockRate = value
    new_file = reset_cfg_env_node.genesis_path
    genesis.to_file(new_file)
    reset_cfg_env_node.deploy_me(new_file)
    config = reset_cfg_env_node.debug.economicConfig()
    assert value == config["reward"]["newBlockRate"]


@pytest.mark.P2
@pytest.mark.parametrize('value', [-1, 110, ""])
def test_IP_PR_010_012(reset_cfg_env_node, value):
    """
    IP_PR_010:奖励池分配给出块奖励的比例<0
    IP_PR_012:奖励池分配给出块奖励的比例>100
    """
    genesis = reset_cfg_env_node.genesis
    genesis.economicModel.reward.newBlockRate = value
    new_file = reset_cfg_env_node.genesis_path
    genesis.to_file(new_file)
    assert_error_deploy(reset_cfg_env_node, new_file, "Abnormal reward pool allocation")


@pytest.mark.P2
@pytest.mark.parametrize('value', ["ss", ""])
def test_IP_CP_001(value, reset_cfg_env_node):
    """
    创世文件链参数验证
    修改eip155Block参数字符串
    修改eip155Block参数空值
    """
    genesis = reset_cfg_env_node.genesis
    genesis.config.eip155Block = value
    new_file = reset_cfg_env_node.genesis_path
    genesis.to_file(new_file)
    assert_error_deploy(reset_cfg_env_node, new_file, "Abnormal eip155Block")


@pytest.mark.P2
@pytest.mark.parametrize('value', ["ss", ""])
def test_IP_CP_002_amount(reset_cfg_env_node, value):
    """
    创世文件共识参数验证
    :return:
    """
    genesis = reset_cfg_env_node.genesis
    genesis.config.cbft.amount = value
    new_file = reset_cfg_env_node.genesis_path
    genesis.to_file(new_file)
    assert_error_deploy(reset_cfg_env_node, new_file, "Abnormal amount")


@pytest.mark.P2
@pytest.mark.parametrize('value', [0.1, "ss", ""])
def test_IP_CP_002_validator_mode(reset_cfg_env_node, value):
    """
    创世文件共识参数验证
    :return:
    """
    if value == "":
        genesis = reset_cfg_env_node.genesis
        genesis.config.cbft.validatorMode = value
        new_file = reset_cfg_env_node.genesis_path
        genesis.to_file(new_file)
        reset_cfg_env_node.deploy_me(new_file)
    else:
        genesis = reset_cfg_env_node.genesis
        genesis.config.cbft.validatorMode = value
        new_file = reset_cfg_env_node.genesis_path
        genesis.to_file(new_file)
        assert_error_deploy(reset_cfg_env_node, new_file, "Abnormal validatorMode")


@pytest.mark.P2
@pytest.mark.parametrize('value', [0.1, "ss", ""])
def test_IP_CP_002_period(reset_cfg_env_node, value):
    """
    创世文件共识参数验证
    :return:
    """
    genesis = reset_cfg_env_node.genesis
    genesis.config.cbft.period = value
    new_file = reset_cfg_env_node.genesis_path
    genesis.to_file(new_file)
    assert_error_deploy(reset_cfg_env_node, new_file, "Abnormal validatorMode")


def test_IP_CP_003_1(reset_cfg_env_node):
    """
    创世文件经济参数验证
    """
    genesis = reset_cfg_env_node.genesis
    genesis.economicModel.common.additionalCycleTime = ""
    new_file = reset_cfg_env_node.genesis_path
    genesis.to_file(new_file)
    assert_error_deploy(reset_cfg_env_node, new_file, "Abnormal economicModel")


@pytest.mark.P2
def test_IP_CP_003_2(reset_cfg_env_node):
    genesis = reset_cfg_env_node.genesis
    genesis.economicModel.innerAcc.cdfAccount = ""
    new_file = reset_cfg_env_node.genesis_path
    genesis.to_file(new_file)
    assert_error_deploy(reset_cfg_env_node, new_file, "Abnormal economicModel")


@pytest.mark.P2
@pytest.mark.parametrize("key", ["slashFractionDuplicateSign", "duplicateSignReportReward", "slashBlocksReward", "maxEvidenceAge"])
def test_IP_CP_004(reset_cfg_env_node, key):
    """
    创世文件处罚参数验证
    :return:
    """
    genesis = reset_cfg_env_node.genesis
    setattr(genesis.economicModel.slashing, key, "")
    new_file = reset_cfg_env_node.genesis_path
    genesis.to_file(new_file)
    assert_error_deploy(reset_cfg_env_node, new_file, "Abnormal Slashing")


@pytest.mark.P2
@pytest.mark.parametrize('value', [0.1, "ss", ""])
def test_IP_CP_005(reset_cfg_env_node, value):
    """
    创世文件共识参数验证
    :return:
    """
    genesis = reset_cfg_env_node.genesis
    genesis.config.chainId = value
    new_file = reset_cfg_env_node.genesis_path
    genesis.to_file(new_file)
    assert_error_deploy(reset_cfg_env_node, new_file, "Abnormal chain id")
