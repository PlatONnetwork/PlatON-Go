import time
import pytest
import allure
from dacite import from_dict
from common.key import get_pub_key, mock_duplicate_sign, generate_key
from common.log import log
from client_sdk_python import Web3
from decimal import Decimal
from tests.conftest import get_clients_noconsensus
from tests.lib import EconomicConfig, Genesis, StakingConfig, Staking, check_node_in_list, assert_code, von_amount, \
    get_governable_parameter_value, Client, update_param_by_dict, get_param_by_dict
from client_sdk_python.debug import Debug


def create_pledge_node_information(clients_noconsensus):
    first_client = clients_noconsensus[0]
    first_economic = first_client.economic
    first_node = first_client.node
    log.info("Current connection node：{}".format(first_client.node.node_mark))
    log.info("Start creating a pledge account Pledge_address")
    staking_address, _ = first_economic.account.generate_account(first_node.web3, von_amount(first_economic.create_staking_limit, 2))
    log.info("Created, account address: {} Amount: {}".format(staking_address, von_amount(first_economic.create_staking_limit, 2)))
    log.info("Start applying for a pledge node")
    result = first_client.staking.create_staking(0, staking_address, staking_address)
    assert_code(result, 0)
    first_economic.wait_settlement_blocknum(first_node)
    log.info("Current block height: {}".format(first_node.eth.blockNumber))
    result = first_node.ppos.getVerifierList()
    log.info("current Verifier List：{}".format(result))


def test_ZB_NP_01(new_genesis_env, clients_noconsensus):
    """
    节点未被选中验证人列表查询零出块记录表
    """
    # Change configuration parameters
    genesis = from_dict(data_class=Genesis, data=new_genesis_env.genesis_config)
    genesis.economicModel.staking.maxValidators = 4
    new_file = new_genesis_env.cfg.env_tmp + "/genesis.json"
    genesis.to_file(new_file)
    new_genesis_env.deploy_all(new_file)
    # start execution use case
    create_pledge_node_information(clients_noconsensus)
    # 查看零出块记录表
    result = Debug.getWaitSlashingNodeList()
    log.info("Slashing NodeList: {}".format(result))


def test_ZB_NP_02(clients_noconsensus):
    """
    节点未被选中共识验证人列表查询零出块记录表（不存在零出块记录）
    """
    # start execution use case
    create_pledge_node_information(clients_noconsensus)

