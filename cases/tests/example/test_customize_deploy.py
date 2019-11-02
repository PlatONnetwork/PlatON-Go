import pytest
import allure
from dacite import from_dict
from copy import copy
from common.log import log
from tests.lib import Genesis
import json


@pytest.fixture
def new_env(global_test_env):
    cfg = copy(global_test_env.cfg)
    yield global_test_env
    log.info("reset deploy.................")
    global_test_env.set_cfg(cfg)
    print(global_test_env.genesis_config)
    print(global_test_env.get_rand_node().genesis_config)
    print(global_test_env.chain_id)


def test_copy_cfg(new_env):
    # Do not copy the environment, retain the original environment configuration through the fixture
    genesis = from_dict(data_class=Genesis, data=new_env.genesis_config)
    genesis.EconomicModel.Staking.StakeThreshold = 500000000000000000000000
    new_env.set_genesis(genesis.to_dict())
    print("test copy cfg")
    print(new_env.get_rand_node().genesis_config)
    new_env.deploy_all()
    log.info("new data:{}".format(json.loads(new_env.get_rand_node().debug.economicConfig())["Staking"]["StakeThreshold"]))


def test_copy_env(global_test_env):
    # Copy the environment, copy the environment, do not configure the environment, account and node copy
    new_env = global_test_env.copy_env()
    genesis = from_dict(data_class=Genesis, data=new_env.genesis_config)
    genesis.EconomicModel.Staking.StakeThreshold = 500000000000000000000000
    new_env.genesis_config = genesis.to_dict()
    new_env.deploy_all()
    log.info("new data:{}".format(json.loads(new_env.get_rand_node().debug.economicConfig())["Staking"]["StakeThreshold"]))
    log.info("start genesis:{}".format(global_test_env.genesis_config))
    global_test_env.deploy_all()
    print("test_copy_env")
    print(global_test_env.get_rand_node().debug.economicConfig())


def test_use_genesis(global_test_env):
    genesis = from_dict(data_class=Genesis, data=global_test_env.genesis_config)
    genesis.EconomicModel.Staking.StakeThreshold = 500000000000000000000000
    new_env.genesis_config = genesis.to_dict()
    new_file = global_test_env.cfg.env_tmp + "/genesis.json"
    genesis.to_file(new_file)
    global_test_env.deploy_all(new_file)
