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


def create_pledge_node_information(client):
    log.info("Current connection node：{}".format(client.node.node_mark))
    log.info("Start creating a pledge account Pledge_address")
    staking_address, _ = client.economic.account.generate_account(client.node.web3, von_amount(client.economic.create_staking_limit, 2))
    log.info("Created, account address: {} Amount: {}".format(staking_address, von_amount(client.economic.create_staking_limit, 2)))
    log.info("Start applying for a pledge node")
    result = client.staking.create_staking(0, staking_address, staking_address)
    assert_code(result, 0)
    client.economic.wait_settlement_blocknum(client.node)
    log.info("Current block height: {}".format(client.node.eth.blockNumber))
    result = client.node.ppos.getVerifierList()
    log.info("current Verifier List：{}".format(result))


def update_zero_produce(new_genesis_env, cumulativetime=4, numberthreshold=3):
    genesis = from_dict(data_class=Genesis, data=new_genesis_env.genesis_config)
    genesis.economicModel.slashing.zeroProduceCumulativeTime = cumulativetime
    genesis.economicModel.slashing.zeroProduceNumberThreshold = numberthreshold
    new_file = new_genesis_env.cfg.env_tmp + "/genesis.json"
    genesis.to_file(new_file)
    new_genesis_env.deploy_all(new_file)


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
    first_client = clients_noconsensus[0]
    create_pledge_node_information(first_client)
    # view the zero out block record table
    result = Debug.getWaitSlashingNodeList(first_client.node)
    log.info("Slashing NodeList: {}".format(result))
    assert result == []


def test_ZB_NP_02(clients_noconsensus):
    """
    节点未被选中共识验证人列表查询零出块记录表（不存在零出块记录）
    """
    # start execution use case
    first_client = clients_noconsensus[0]
    first_economic = first_client.economic
    first_node = first_client.node
    create_pledge_node_information(first_client)
    result = check_node_in_list(first_node.node_id, first_client.ppos.getValidatorList)
    log.info("Current node in consensus list status：{}".format(result))
    if result is False:
        # view the zero out block record table
        result = Debug.getWaitSlashingNodeList(first_client.node)
        log.info("Slashing NodeList: {}".format(result))
        assert result == []
    else:
        first_economic.wait_consensus_blocknum(first_node)


def test_ZB_NP_03(new_genesis_env, clients_noconsensus):
    """
    节点未被选中共识验证人列表查询零出块记录表（存在零出块记录）
    """
    update_zero_produce(new_genesis_env)
    first_client = clients_noconsensus[0]
    second_client = clients_noconsensus[1]
    first_economic = first_client.economic
    first_node = first_client.node

    create_pledge_node_information(first_client)
    verifier_nodeid_list = first_client.economic.env.consensus_node_id_list()
    verifier_nodeid_list.append(first_node.node_id)
    verifier_nodeid_list.pop(4)
    log.info("verifier nodeid list: {}".format(verifier_nodeid_list))
    for i in range(4):
        result = check_node_in_list(first_node.node_id, first_client.ppos.getValidatorList)
        log.info("Current node in consensus list status：{}".format(result))
        if result:
            result = second_client.node.debug.setValidatorList(verifier_nodeid_list)
            assert result is None
            # stop node
            first_client.node.stop()
            second_client.economic.wait_consensus_blocknum(second_client.node)
            log.info("Current block height: {}".format(second_client.node.eth.blockNumber))
            result = check_node_in_list(first_node.node_id, second_client.ppos.getValidatorList)
            log.info("node in consensus list status：{}".format(result))
            result = second_client.ppos.getValidatorList()
            log.info("Validator List:{}".format(result))
            first_node.start()
            second_client.economic.wait_consensus_blocknum(second_client.node)
            log.info("Current block height: {}".format(second_client.node.eth.blockNumber))
            result = second_client.node.debug.getWaitSlashingNodeList()
            log.info("Slashing NodeList: {}".format(result))
            second_client.economic.wait_consensus_blocknum(second_client.node)
            result = second_client.node.debug.getWaitSlashingNodeList()
            log.info("Slashing NodeList: {}".format(result))
            break
        else:
            first_economic.wait_consensus_blocknum(first_node)


def test_ZB_NP_04(clients_noconsensus):
    """
    节点被选中共识验证人查询零出块记录表（不存在零出块记录）
    """
    first_client = clients_noconsensus[0]
    first_economic = first_client.economic
    first_node = first_client.node
    create_pledge_node_information(first_client)
    for i in range(4):
        result = check_node_in_list(first_node.node_id, first_client.ppos.getValidatorList)
        log.info("Current node in consensus list status：{}".format(result))
        if result:
            # view the zero out block record table
            result = Debug.getWaitSlashingNodeList(first_client.node)
            log.info("Slashing NodeList: {}".format(result))
            assert result == []
            break
        else:
            first_economic.wait_consensus_blocknum(first_node)


def test_ZB_NP_05(new_genesis_env, clients_noconsensus):
    """
    节点被选中共识验证人查询零出块记录表（存在零出块记录）
    """
    update_zero_produce(new_genesis_env)
    first_client = clients_noconsensus[0]
    second_client = clients_noconsensus[1]
    first_economic = first_client.economic
    first_node = first_client.node
    create_pledge_node_information(first_client)
    verifier_nodeid_list = first_client.economic.env.consensus_node_id_list()
    verifier_nodeid_list.append(first_node.node_id)
    verifier_nodeid_list.pop(0)
    log.info("verifier nodeid list: {}".format(verifier_nodeid_list))
    second_client.node.debug.setValidatorList(verifier_nodeid_list)
    first_node.stop()
    second_client.economic.wait_consensus_blocknum(second_client.node, 1)
    first_node.start()
    second_client.economic.wait_consensus_blocknum(second_client.node)
    result = second_client.node.debug.getWaitSlashingNodeList()
    log.info("Slashing Node List:{}".format(result))
    second_client.economic.wait_consensus_blocknum(second_client.node)
    result = second_client.node.debug.getWaitSlashingNodeList()
    log.info("Slashing Node List:{}".format(result))


def test_ZB_NP_06(new_genesis_env, clients_noconsensus):
    """
    节点被选中共识验证人未出块查询零出块记录表（不存在零出块记录）
    """
    update_zero_produce(new_genesis_env)
    first_client = clients_noconsensus[0]
    second_client = clients_noconsensus[1]
    first_economic = first_client.economic
    first_node = first_client.node
    create_pledge_node_information(first_client)
    verifier_nodeid_list = first_client.economic.env.consensus_node_id_list()
    verifier_nodeid_list.append(first_node.node_id)
    verifier_nodeid_list.pop(0)
    log.info("verifier nodeid list: {}".format(verifier_nodeid_list))
    second_client.node.debug.setValidatorList(verifier_nodeid_list)
    first_node.stop()
    second_client.economic.wait_consensus_blocknum(second_client.node, 1)
    result = second_client.node.debug.getWaitSlashingNodeList()
    log.info("Slashing Node List:{}".format(result))

