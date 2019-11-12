from dacite import from_dict

from common.log import log
from tests.lib.genesis import Genesis
from environment.env import TestEnvironment


def assert_error_deploy(env:TestEnvironment, genesis_file, msg="Error config"):
    try:
        env.deploy_all(genesis_file)
        assert False, "{},but deploy success".format(msg)
    except Exception as e:
        log.info("Deploy error info:{}".format(e))


def update_staking(reset_cfg_env, value):
    genesis = from_dict(data_class=Genesis, data=reset_cfg_env.genesis_cfg)
    genesis.economicModel.staking.stakeThreshold = value
    new_file = reset_cfg_env.cfg.env_tmp + "/genesis.json"
    genesis.to_file(new_file)
    return new_file


def update_unstaking(reset_cfg_env, value):
    genesis = from_dict(data_class=Genesis, data=reset_cfg_env.genesis_cfg)
    genesis.economicModel.staking.operatingThreshold = value
    new_file = reset_cfg_env.cfg.env_tmp + "/genesis.json"
    genesis.to_file(new_file)
    return new_file


def update_hesitate_ratio(reset_cfg_env, value):
    genesis = from_dict(data_class=Genesis, data=reset_cfg_env.genesis_cfg)
    genesis.economicModel.staking.hesitateRatio = value
    new_file = reset_cfg_env.cfg.env_tmp + "/genesis.json"
    genesis.to_file(new_file)
    return new_file


def update_unstake_freeze_duration(reset_cfg_env, value):
    genesis = from_dict(data_class=Genesis, data=reset_cfg_env.genesis_cfg)
    genesis.economicModel.staking.unStakeFreezeDuration = value
    new_file = reset_cfg_env.cfg.env_tmp + "/genesis.json"
    genesis.to_file(new_file)
    return new_file