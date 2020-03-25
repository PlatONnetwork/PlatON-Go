import time
import pytest
import allure
from dacite import from_dict
from common.key import get_pub_key, mock_duplicate_sign, generate_key
from common.log import log
from client_sdk_python import Web3
from decimal import Decimal
from tests.conftest import get_clients_noconsensus, param_governance_verify_before_endblock, param_governance_verify
from tests.lib import EconomicConfig, Genesis, StakingConfig, Staking, check_node_in_list, assert_code, von_amount, \
    get_governable_parameter_value, Client, update_param_by_dict, get_param_by_dict, get_pledge_list
from client_sdk_python.debug import Debug


def create_pledge_node_information(client):
    log.info("Current connection node：{}".format(client.node.node_mark))
    log.info("Start creating a pledge account Pledge_address")
    staking_address, pri_key = client.economic.account.generate_account(client.node.web3,
                                                                        von_amount(client.economic.create_staking_limit,
                                                                                   2))
    log.info("Created, account address: {} Amount: {}".format(staking_address,
                                                              von_amount(client.economic.create_staking_limit, 2)))
    log.info("Start applying for a pledge node")
    result = client.staking.create_staking(0, staking_address, staking_address)
    assert_code(result, 0)
    client.economic.wait_settlement_blocknum(client.node)
    log.info("Current block height: {}".format(client.node.eth.blockNumber))
    result = client.node.ppos.getVerifierList()
    log.info("current Verifier List：{}".format(result))
    return staking_address, pri_key


def update_zero_produce(new_genesis_env, cumulativetime=4, numberthreshold=3):
    genesis = from_dict(data_class=Genesis, data=new_genesis_env.genesis_config)
    genesis.economicModel.slashing.zeroProduceCumulativeTime = cumulativetime
    genesis.economicModel.slashing.zeroProduceNumberThreshold = numberthreshold
    new_file = new_genesis_env.cfg.env_tmp + "/genesis.json"
    genesis.to_file(new_file)
    new_genesis_env.deploy_all(new_file)


def get_slash_count(wait_slashing_node_list, node_id):
    for slashing in wait_slashing_node_list:
        if slashing["NodeId"] == node_id:
            return slashing["CountBit"]
    return 0


def to_int(value):
    return int(str(value), 2)


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
    first_node = first_client.node
    create_pledge_node_information(first_client)
    # view the zero out block record table
    result = Debug.getWaitSlashingNodeList(first_client.node)
    log.info("Slashing NodeList: {}".format(result))
    wait_slashing_list = get_slash_count(result, first_node.node_id)
    assert wait_slashing_list == 0


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
        wait_slashing_list = get_slash_count(result, first_node.node_id)
        assert wait_slashing_list == 0
    else:
        first_economic.wait_consensus_blocknum(first_node)


def test_ZB_NP_03(new_genesis_env, clients_noconsensus):
    """
    节点未被选中共识验证人列表查询零出块记录表（存在零出块记录）
    """
    update_zero_produce(new_genesis_env)
    first_client = clients_noconsensus[0]
    second_client = clients_noconsensus[1]
    first_node = first_client.node
    # start execution use case
    staking_address, pri_key = create_pledge_node_information(first_client)
    log.info("current validator: {}".format(get_pledge_list(first_node.ppos.getValidatorList)))
    verifier_nodeid_list = first_client.economic.env.consensus_node_id_list()
    verifier_nodeid_list.append(first_node.node_id)
    verifier_nodeid_list.pop(0)
    log.info("verifier nodeid list: {}".format(verifier_nodeid_list))
    result = second_client.node.debug.setValidatorList(verifier_nodeid_list, pri_key)
    assert_code(result, 0)
    # stop node
    first_client.node.stop()

    second_client.economic.wait_consensus_blocknum(second_client.node)
    log.info("current validator: {}".format(get_pledge_list(second_client.ppos.getValidatorList)))
    log.info("Current block height: {}".format(second_client.node.eth.blockNumber))
    result = second_client.node.debug.setValidatorList(second_client.economic.env.consensus_node_id_list(), pri_key)
    assert_code(result, 0)

    second_client.economic.wait_consensus_blocknum(second_client.node)
    log.info("current validator: {}".format(get_pledge_list(second_client.ppos.getValidatorList)))
    log.info("Current block height: {}".format(second_client.node.eth.blockNumber))
    result = Debug.getWaitSlashingNodeList(second_client.node)
    log.info("Slashing NodeList: {}".format(result))
    wait_slashing_list = get_slash_count(result, first_node.node_id)
    assert wait_slashing_list == 0

    second_client.economic.wait_consensus_blocknum(second_client.node)
    log.info("current validator: {}".format(get_pledge_list(second_client.ppos.getValidatorList)))
    log.info("Current block height: {}".format(second_client.node.eth.blockNumber))
    result = Debug.getWaitSlashingNodeList(second_client.node)
    log.info("Slashing NodeList: {}".format(result))
    wait_slashing_list = get_slash_count(result, first_node.node_id)
    assert wait_slashing_list == to_int(1)

    second_client.economic.wait_consensus_blocknum(second_client.node)
    log.info("current validator: {}".format(get_pledge_list(second_client.ppos.getValidatorList)))
    log.info("Current block height: {}".format(second_client.node.eth.blockNumber))
    result = Debug.getWaitSlashingNodeList(second_client.node)
    log.info("Slashing NodeList: {}".format(result))
    wait_slashing_list = get_slash_count(result, first_node.node_id)
    assert wait_slashing_list == to_int(1)


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
            result = first_node.debug.getWaitSlashingNodeList()
            log.info("Slashing NodeList: {}".format(result))
            wait_slashing_list = get_slash_count(result, first_node.node_id)
            assert wait_slashing_list == 0
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
    first_node = first_client.node
    staking_address, pri_key = create_pledge_node_information(first_client)
    verifier_nodeid_list = first_client.economic.env.consensus_node_id_list()
    verifier_nodeid_list.append(first_node.node_id)
    verifier_nodeid_list.pop(0)
    log.info("verifier nodeid list: {}".format(verifier_nodeid_list))
    result = second_client.node.debug.setValidatorList(verifier_nodeid_list, pri_key)
    assert_code(result, 0)
    first_node.stop()

    second_client.economic.wait_consensus_blocknum(second_client.node)
    result = second_client.node.debug.setValidatorList(verifier_nodeid_list, pri_key)
    assert_code(result, 0)

    second_client.economic.wait_consensus_blocknum(second_client.node)
    result = second_client.node.debug.getWaitSlashingNodeList()
    log.info("Slashing Node List:{}".format(result))
    wait_slashing_list = get_slash_count(result, first_node.node_id)
    assert wait_slashing_list == 0
    first_node.start(False)

    second_client.economic.wait_consensus_blocknum(second_client.node)
    result = second_client.node.debug.getWaitSlashingNodeList()
    log.info("Slashing Node List:{}".format(result))
    wait_slashing_list = get_slash_count(result, first_node.node_id)
    assert wait_slashing_list == to_int(1)

    second_client.economic.wait_consensus_blocknum(second_client.node)
    result = second_client.node.debug.getWaitSlashingNodeList()
    log.info("Slashing Node List:{}".format(result))
    wait_slashing_list = get_slash_count(result, first_node.node_id)
    assert wait_slashing_list == 0
    result = second_client.ppos.getCandidateInfo(first_node.node_id)
    log.info("Candidate Info:{}".format(result))
    #
    # second_client.economic.wait_consensus_blocknum(second_client.node)
    # result = second_client.node.debug.getWaitSlashingNodeList()
    # log.info("Slashing Node List:{}".format(result))
    # wait_slashing_list = get_slash_count(result, first_node.node_id)
    # assert wait_slashing_list == to_int(11)
    #
    # second_client.economic.wait_consensus_blocknum(second_client.node)
    # result = second_client.node.debug.getWaitSlashingNodeList()
    # log.info("Slashing Node List:{}".format(result))
    # wait_slashing_list = get_slash_count(result, first_node.node_id)
    # assert wait_slashing_list == to_int(111)


def test_ZB_NP_06_07(new_genesis_env, clients_noconsensus):
    """
    ZB_NP_06:节点被选中共识验证人未出块查询零出块记录表（不存在零出块记录）
    ZB_NP_07:节点被选中共识验证人未出块查询零出块记录表（存在零出块记录）
    """
    update_zero_produce(new_genesis_env)
    first_client = clients_noconsensus[0]
    second_client = clients_noconsensus[1]
    first_node = first_client.node
    staking_address, pri_key = create_pledge_node_information(first_client)

    verifier_nodeid_list = first_client.economic.env.consensus_node_id_list()
    verifier_nodeid_list.append(first_node.node_id)
    verifier_nodeid_list.pop(0)
    log.info("verifier nodeid list: {}".format(verifier_nodeid_list))
    result = second_client.node.debug.setValidatorList(verifier_nodeid_list, pri_key)
    assert_code(result, 0)
    first_node.stop()

    second_client.economic.wait_consensus_blocknum(second_client.node)
    result = second_client.node.debug.setValidatorList(verifier_nodeid_list, pri_key)
    assert_code(result, 0)

    second_client.economic.wait_consensus_blocknum(second_client.node)
    result = second_client.node.debug.getWaitSlashingNodeList()
    log.info("Slashing Node List:{}".format(result))
    wait_slashing_list = get_slash_count(result, first_node.node_id)
    assert wait_slashing_list == 0

    second_client.economic.wait_consensus_blocknum(second_client.node)
    result = second_client.node.debug.getWaitSlashingNodeList()
    log.info("Slashing Node List:{}".format(result))
    wait_slashing_list = get_slash_count(result, first_node.node_id)
    assert wait_slashing_list == to_int(1)

    second_client.economic.wait_consensus_blocknum(second_client.node)
    result = second_client.node.debug.getWaitSlashingNodeList()
    log.info("Slashing Node List:{}".format(result))
    wait_slashing_list = get_slash_count(result, first_node.node_id)
    assert wait_slashing_list == to_int(11)


def test_ZB_NP_08(new_genesis_env, clients_noconsensus):
    """
    节点达到共识轮阈值，未达到累计次数被选中共识验证人未出块
    """
    update_zero_produce(new_genesis_env)
    first_client = clients_noconsensus[0]
    second_client = clients_noconsensus[1]
    first_node = first_client.node
    staking_address, pri_key = create_pledge_node_information(first_client)

    log.info("current validator: {}".format(get_pledge_list(second_client.ppos.getValidatorList)))
    verifier_nodeid_list = first_client.economic.env.consensus_node_id_list()
    verifier_nodeid_list.append(first_node.node_id)
    verifier_nodeid_list.pop(0)
    log.info("verifier nodeid list: {}".format(verifier_nodeid_list))
    result = second_client.node.debug.setValidatorList(verifier_nodeid_list, pri_key)
    assert_code(result, 0)
    first_node.stop()

    second_client.economic.wait_consensus_blocknum(second_client.node)
    log.info("current validator: {}".format(get_pledge_list(second_client.ppos.getValidatorList)))
    result = second_client.node.debug.setValidatorList(second_client.economic.env.consensus_node_id_list(), pri_key)
    assert_code(result, 0)

    second_client.economic.wait_consensus_blocknum(second_client.node)
    log.info("current validator: {}".format(get_pledge_list(second_client.ppos.getValidatorList)))
    result = second_client.node.debug.setValidatorList(second_client.economic.env.consensus_node_id_list(), pri_key)
    assert_code(result, 0)
    result = second_client.node.debug.getWaitSlashingNodeList()
    log.info("Slashing Node List:{}".format(result))
    wait_slashing_list = get_slash_count(result, first_node.node_id)
    assert wait_slashing_list == 0

    second_client.economic.wait_consensus_blocknum(second_client.node)
    log.info("current validator: {}".format(get_pledge_list(second_client.ppos.getValidatorList)))
    result = second_client.node.debug.setValidatorList(verifier_nodeid_list, pri_key)
    assert_code(result, 0)
    result = second_client.node.debug.getWaitSlashingNodeList()
    log.info("Slashing Node List:{}".format(result))
    wait_slashing_list = get_slash_count(result, first_node.node_id)
    assert wait_slashing_list == to_int(1)

    second_client.economic.wait_consensus_blocknum(second_client.node)
    log.info("current validator: {}".format(get_pledge_list(second_client.ppos.getValidatorList)))
    result = second_client.node.debug.setValidatorList(verifier_nodeid_list, pri_key)
    assert_code(result, 0)
    result = second_client.node.debug.getWaitSlashingNodeList()
    log.info("Slashing Node List:{}".format(result))
    wait_slashing_list = get_slash_count(result, first_node.node_id)
    assert wait_slashing_list == to_int(1)

    second_client.economic.wait_consensus_blocknum(second_client.node)
    log.info("current validator: {}".format(get_pledge_list(second_client.ppos.getValidatorList)))
    result = second_client.node.debug.getWaitSlashingNodeList()
    log.info("Slashing Node List:{}".format(result))
    wait_slashing_list = get_slash_count(result, first_node.node_id)
    assert wait_slashing_list == to_int(1)

    second_client.economic.wait_consensus_blocknum(second_client.node)
    log.info("current validator: {}".format(get_pledge_list(second_client.ppos.getValidatorList)))
    result = second_client.node.debug.getWaitSlashingNodeList()
    log.info("Slashing Node List:{}".format(result))
    wait_slashing_list = get_slash_count(result, first_node.node_id)
    assert wait_slashing_list == to_int(1001)

    second_client.economic.wait_consensus_blocknum(second_client.node)
    log.info("current validator: {}".format(get_pledge_list(second_client.ppos.getValidatorList)))
    result = second_client.node.debug.getWaitSlashingNodeList()
    log.info("Slashing Node List:{}".format(result))
    wait_slashing_list = get_slash_count(result, first_node.node_id)
    assert wait_slashing_list == to_int(11)


def test_ZB_NP_09(new_genesis_env, clients_noconsensus):
    """
    节点先双签被举报再触发零出块
    """
    update_zero_produce(new_genesis_env)
    first_client = clients_noconsensus[0]
    second_client = clients_noconsensus[1]
    first_economic = first_client.economic
    first_node = first_client.node
    report_address, _ = first_economic.account.generate_account(first_node.web3, first_node.web3.toWei(1000, 'ether'))
    create_pledge_node_information(first_client)
    # verifier_nodeid_list = first_client.economic.env.consensus_node_id_list()
    # verifier_nodeid_list.append(first_node.node_id)
    # new_verifier_nodeid_list = verifier_nodeid_list.pop(0)
    # log.info("verifier nodeid list: {}".format(new_verifier_nodeid_list))
    # second_client.node.debug.setValidatorList(new_verifier_nodeid_list, pri_key)
    first_economic.wait_consensus_blocknum(first_node)
    log.info("Current block height: {}".format(first_node.eth.blockNumber))
    report_information = mock_duplicate_sign(1, first_node.nodekey, first_node.blsprikey, first_node.eth.blockNumber)
    log.info("Report information: {}".format(report_information))
    result = first_client.duplicatesign.reportDuplicateSign(1, report_information, report_address)
    assert_code(result, 0)
    first_node.stop()
    second_client.economic.wait_consensus_blocknum(second_client.node)
    result = second_client.node.debug.getWaitSlashingNodeList()
    log.info("Slashing Node List:{}".format(result))
    wait_slashing_list = get_slash_count(result, first_node.node_id)
    assert wait_slashing_list == 0
    second_client.economic.wait_consensus_blocknum(second_client.node)
    result = second_client.node.debug.getWaitSlashingNodeList()
    log.info("Slashing Node List:{}".format(result))
    wait_slashing_list = get_slash_count(result, first_node.node_id)
    assert wait_slashing_list == to_int(1)
    second_client.economic.wait_consensus_blocknum(second_client.node)
    result = second_client.node.debug.getWaitSlashingNodeList()
    log.info("Slashing Node List:{}".format(result))
    wait_slashing_list = get_slash_count(result, first_node.node_id)
    assert wait_slashing_list == to_int(1)


def test_ZB_NP_10(new_genesis_env, clients_noconsensus):
    """
    节点触发零出块再双签被举报
    """
    update_zero_produce(new_genesis_env)
    first_client = clients_noconsensus[0]
    second_client = clients_noconsensus[1]
    first_economic = first_client.economic
    first_node = first_client.node
    report_address, _ = first_economic.account.generate_account(first_node.web3, first_node.web3.toWei(1000, 'ether'))
    staking_address, pri_key = create_pledge_node_information(first_client)

    verifier_nodeid_list = first_client.economic.env.consensus_node_id_list()
    verifier_nodeid_list.append(first_node.node_id)
    verifier_nodeid_list.pop(0)
    log.info("verifier nodeid list: {}".format(verifier_nodeid_list))
    result = first_node.debug.setValidatorList(verifier_nodeid_list, pri_key)
    assert_code(result, 0)
    first_node.stop()

    second_client.economic.wait_consensus_blocknum(second_client.node)
    current_block = second_client.node.eth.blockNumber
    log.info("Current block height: {}".format(current_block))
    result = second_client.node.debug.setValidatorList(verifier_nodeid_list, pri_key)
    assert_code(result, 0)

    second_client.economic.wait_consensus_blocknum(second_client.node)
    log.info("current validator: {}".format(get_pledge_list(second_client.ppos.getValidatorList)))
    result = second_client.node.debug.getWaitSlashingNodeList()
    log.info("Slashing Node List:{}".format(result))
    wait_slashing_list = get_slash_count(result, first_node.node_id)
    assert wait_slashing_list == 0

    second_client.economic.wait_consensus_blocknum(second_client.node)
    log.info("current validator: {}".format(get_pledge_list(second_client.ppos.getValidatorList)))
    result = second_client.node.debug.getWaitSlashingNodeList()
    log.info("Slashing Node List:{}".format(result))
    wait_slashing_list = get_slash_count(result, first_node.node_id)
    assert wait_slashing_list == to_int(1)
    report_information = mock_duplicate_sign(1, first_node.nodekey, first_node.blsprikey, current_block)
    log.info("Report information: {}".format(report_information))
    result = second_client.duplicatesign.reportDuplicateSign(1, report_information, report_address)
    assert_code(result, 0)

    second_client.economic.wait_consensus_blocknum(second_client.node)
    log.info("current validator: {}".format(get_pledge_list(second_client.ppos.getValidatorList)))
    result = check_node_in_list(first_node.node_id, second_client.ppos.getVerifierList)
    log.info("Current node in consensus list status：{}".format(result))
    assert result is False
    result = second_client.ppos.getCandidateInfo(first_node.node_id)
    log.info("Candidate Info :{}".format(result))
    # result['Ret'][]
    result = second_client.node.debug.getWaitSlashingNodeList()
    log.info("Slashing Node List:{}".format(result))
    wait_slashing_list = get_slash_count(result, first_node.node_id)
    assert wait_slashing_list == to_int(11)


def test_ZB_NP_31(new_genesis_env, clients_noconsensus):
    """
    治理提案投票未通过，节点零出块率处罚（不存在零出块记录）
    """
    update_zero_produce(new_genesis_env, 3)
    first_client = clients_noconsensus[0]
    second_client = clients_noconsensus[1]
    first_economic = first_client.economic
    first_node = first_client.node

    # Voting failed and expired
    param_governance_verify(first_client, 'slashing', 'zeroProduceCumulativeTime', '4', False)
    # Start execution use case
    staking_address, pri_key = create_pledge_node_information(first_client)

    verifier_nodeid_list = first_client.economic.env.consensus_node_id_list()
    verifier_nodeid_list.append(first_node.node_id)
    verifier_nodeid_list.pop(0)
    log.info("verifier nodeid list: {}".format(verifier_nodeid_list))
    result = second_client.node.debug.setValidatorList(verifier_nodeid_list, pri_key)
    assert_code(result, 0)
    first_node.stop()

    second_client.economic.wait_consensus_blocknum(second_client.node)
    result = second_client.node.debug.setValidatorList(verifier_nodeid_list, pri_key)
    assert_code(result, 0)

    second_client.economic.wait_consensus_blocknum(second_client.node)
    result = second_client.node.debug.setValidatorList(verifier_nodeid_list, pri_key)
    assert_code(result, 0)
    result = second_client.node.debug.getWaitSlashingNodeList()
    log.info("Slashing Node List:{}".format(result))
    wait_slashing_list = get_slash_count(result, first_node.node_id)
    assert wait_slashing_list == 0

    second_client.economic.wait_consensus_blocknum(second_client.node)
    result = second_client.node.debug.getWaitSlashingNodeList()
    log.info("Slashing Node List:{}".format(result))
    wait_slashing_list = get_slash_count(result, first_node.node_id)
    assert wait_slashing_list == to_int(1)

    second_client.economic.wait_consensus_blocknum(second_client.node)
    result = second_client.node.debug.getWaitSlashingNodeList()
    log.info("Slashing Node List:{}".format(result))
    wait_slashing_list = get_slash_count(result, first_node.node_id)
    assert wait_slashing_list == to_int(11)

    second_client.economic.wait_consensus_blocknum(second_client.node)
    result = second_client.node.debug.getWaitSlashingNodeList()
    log.info("Slashing Node List:{}".format(result))
    wait_slashing_list = get_slash_count(result, first_node.node_id)
    assert wait_slashing_list == 0
    result = second_client.ppos.getCandidateInfo(first_node.node_id)
    log.info("Candidate Info:{}".format(result))
    assert result['Ret']['Shares'] != 1000000000000000000000000


def test_ZB_NP_32(new_genesis_env, clients_noconsensus):
    """
    治理提案投票未通过，节点零出块率处罚（存在零出块记录）
    """
    update_zero_produce(new_genesis_env, 3)
    first_client = clients_noconsensus[0]
    second_client = clients_noconsensus[1]
    first_economic = first_client.economic
    first_node = first_client.node

    # Voting failed and expired
    param_governance_verify_before_endblock(first_client, 'slashing', 'zeroProduceCumulativeTime', '4', False)
    # Start execution use case
    staking_address, pri_key = create_pledge_node_information(first_client)

    log.info("Current block height: {}".format(first_node.eth.blockNumber))
    first_economic.wait_consensus_blocknum()
    verifier_nodeid_list = first_client.economic.env.consensus_node_id_list()
    verifier_nodeid_list.append(first_node.node_id)
    verifier_nodeid_list.pop(0)
    log.info("verifier nodeid list: {}".format(verifier_nodeid_list))
    result = second_client.node.debug.setValidatorList(first_economic.env.consensus_node_id_list(), pri_key)
    assert_code(result, 0)
    first_node.stop()

    second_client.economic.wait_consensus_blocknum(second_client.node)
    log.info("Current block height: {}".format(second_client.node.eth.blockNumber))
    result = second_client.node.debug.setValidatorList(verifier_nodeid_list, pri_key)
    assert_code(result, 0)

    second_client.economic.wait_consensus_blocknum(second_client.node)
    log.info("Current block height: {}".format(second_client.node.eth.blockNumber))
    result = second_client.node.debug.setValidatorList(verifier_nodeid_list, pri_key)
    assert_code(result, 0)
    result = second_client.node.debug.getWaitSlashingNodeList()
    log.info("Slashing Node List:{}".format(result))
    wait_slashing_list = get_slash_count(result, first_node.node_id)
    assert wait_slashing_list == 0

    second_client.economic.wait_consensus_blocknum(second_client.node)
    log.info("Current block height: {}".format(second_client.node.eth.blockNumber))
    result = second_client.node.debug.setValidatorList(verifier_nodeid_list, pri_key)
    assert_code(result, 0)
    result = second_client.node.debug.getWaitSlashingNodeList()
    log.info("Slashing Node List:{}".format(result))
    wait_slashing_list = get_slash_count(result, first_node.node_id)
    assert wait_slashing_list == 0

    second_client.economic.wait_consensus_blocknum(second_client.node)
    log.info("Current block height: {}".format(second_client.node.eth.blockNumber))
    result = second_client.node.debug.setValidatorList(verifier_nodeid_list, pri_key)
    assert_code(result, 0)
    result = second_client.node.debug.getWaitSlashingNodeList()
    log.info("Slashing Node List:{}".format(result))
    wait_slashing_list = get_slash_count(result, first_node.node_id)
    assert wait_slashing_list == to_int(1)

    second_client.economic.wait_consensus_blocknum(second_client.node)
    log.info("Current block height: {}".format(second_client.node.eth.blockNumber))
    result = second_client.node.debug.setValidatorList(verifier_nodeid_list, pri_key)
    assert_code(result, 0)
    result = second_client.node.debug.getWaitSlashingNodeList()
    log.info("Slashing Node List:{}".format(result))
    wait_slashing_list = get_slash_count(result, first_node.node_id)
    assert wait_slashing_list == to_int(11)

    second_client.economic.wait_consensus_blocknum(second_client.node)
    result = second_client.node.debug.getWaitSlashingNodeList()
    log.info("Slashing Node List:{}".format(result))
    wait_slashing_list = get_slash_count(result, first_node.node_id)
    assert wait_slashing_list == 0
    result = second_client.ppos.getCandidateInfo(first_node.node_id)
    log.info("Candidate Info:{}".format(result))
    assert result['Ret']['Shares'] != 1000000000000000000000000


def test_ZB_NP_33(new_genesis_env, clients_noconsensus):
    """
    共识轮阈值小值变大值提案未生效时，节点零出块率处罚（不存在零出块记录）
    """
    update_zero_produce(new_genesis_env, 3)
    first_client = clients_noconsensus[0]
    second_client = clients_noconsensus[1]
    first_node = first_client.node
    staking_address, pri_key = create_pledge_node_information(first_client)

    # Voting failed and expired
    param_governance_verify_before_endblock(first_client, 'slashing', 'zeroProduceCumulativeTime', '4', True)

    # Start execution use case
    verifier_nodeid_list = first_client.economic.env.consensus_node_id_list()
    verifier_nodeid_list.append(first_node.node_id)
    verifier_nodeid_list.pop(0)
    log.info("verifier nodeid list: {}".format(verifier_nodeid_list))
    result = second_client.node.debug.setValidatorList(verifier_nodeid_list, pri_key)
    assert_code(result, 0)
    first_node.stop()

    second_client.economic.wait_consensus_blocknum(second_client.node)
    result = second_client.node.debug.setValidatorList(verifier_nodeid_list, pri_key)
    assert_code(result, 0)

    second_client.economic.wait_consensus_blocknum(second_client.node)
    result = second_client.node.debug.setValidatorList(verifier_nodeid_list, pri_key)
    assert_code(result, 0)
    result = second_client.node.debug.getWaitSlashingNodeList()
    log.info("Slashing Node List:{}".format(result))
    wait_slashing_list = get_slash_count(result, first_node.node_id)
    assert wait_slashing_list == 0

    second_client.economic.wait_consensus_blocknum(second_client.node)
    result = second_client.node.debug.getWaitSlashingNodeList()
    log.info("Slashing Node List:{}".format(result))
    wait_slashing_list = get_slash_count(result, first_node.node_id)
    assert wait_slashing_list == to_int(1)

    second_client.economic.wait_consensus_blocknum(second_client.node)
    result = second_client.node.debug.getWaitSlashingNodeList()
    log.info("Slashing Node List:{}".format(result))
    wait_slashing_list = get_slash_count(result, first_node.node_id)
    assert wait_slashing_list == to_int(11)

    second_client.economic.wait_consensus_blocknum(second_client.node)
    result = second_client.node.debug.getWaitSlashingNodeList()
    log.info("Slashing Node List:{}".format(result))
    wait_slashing_list = get_slash_count(result, first_node.node_id)
    assert wait_slashing_list == 0
    result = second_client.ppos.getCandidateInfo(first_node.node_id)
    log.info("Candidate Info:{}".format(result))
    assert result['Ret']['Shares'] != 1000000000000000000000000


def test_ZB_NP_34(new_genesis_env, clients_noconsensus):
    """
    共识轮阈值小值变大值提案未生效时，节点零出块率处罚（存在零出块记录）
    """
    update_zero_produce(new_genesis_env, 3)
    first_client = clients_noconsensus[0]
    second_client = clients_noconsensus[1]
    first_node = first_client.node
    staking_address, pri_key = create_pledge_node_information(first_client)

    # Start execution use case
    verifier_nodeid_list = first_client.economic.env.consensus_node_id_list()
    verifier_nodeid_list.append(first_node.node_id)
    verifier_nodeid_list.pop(0)
    log.info("verifier nodeid list: {}".format(verifier_nodeid_list))
    result = second_client.node.debug.setValidatorList(verifier_nodeid_list, pri_key)
    assert_code(result, 0)
    first_node.stop()

    second_client.economic.wait_consensus_blocknum(second_client.node)
    result = second_client.node.debug.setValidatorList(verifier_nodeid_list, pri_key)
    assert_code(result, 0)

    second_client.economic.wait_consensus_blocknum(second_client.node)
    result = second_client.node.debug.setValidatorList(verifier_nodeid_list, pri_key)
    assert_code(result, 0)
    result = second_client.node.debug.getWaitSlashingNodeList()
    log.info("Slashing Node List:{}".format(result))
    wait_slashing_list = get_slash_count(result, first_node.node_id)
    assert wait_slashing_list == 0

    second_client.economic.wait_consensus_blocknum(second_client.node)
    result = second_client.node.debug.getWaitSlashingNodeList()
    log.info("Slashing Node List:{}".format(result))
    wait_slashing_list = get_slash_count(result, first_node.node_id)
    assert wait_slashing_list == to_int(1)
    # Voting failed and expired
    param_governance_verify_before_endblock(second_client, 'slashing', 'zeroProduceCumulativeTime', '4', True)

    second_client.economic.wait_consensus_blocknum(second_client.node)
    result = second_client.node.debug.getWaitSlashingNodeList()
    log.info("Slashing Node List:{}".format(result))
    wait_slashing_list = get_slash_count(result, first_node.node_id)
    assert wait_slashing_list == to_int(11)

    second_client.economic.wait_consensus_blocknum(second_client.node)
    result = second_client.node.debug.getWaitSlashingNodeList()
    log.info("Slashing Node List:{}".format(result))
    wait_slashing_list = get_slash_count(result, first_node.node_id)
    assert wait_slashing_list == 0
    result = second_client.ppos.getCandidateInfo(first_node.node_id)
    log.info("Candidate Info:{}".format(result))
    assert result['Ret']['Shares'] != 1000000000000000000000000


def test_ZB_NP_35(new_genesis_env, clients_noconsensus):
    """
    共识轮阈值大值变小值提案未生效时，节点零出块率处罚（不存在零出块记录）
    """
    update_zero_produce(new_genesis_env)
    first_client = clients_noconsensus[0]
    second_client = clients_noconsensus[1]
    first_node = first_client.node
    staking_address, pri_key = create_pledge_node_information(first_client)

    # Voting failed and expired
    param_governance_verify_before_endblock(first_client, 'slashing', 'zeroProduceCumulativeTime', '3', True)

    # Start execution use case
    verifier_nodeid_list = first_client.economic.env.consensus_node_id_list()
    verifier_nodeid_list.append(first_node.node_id)
    verifier_nodeid_list.pop(0)
    log.info("verifier nodeid list: {}".format(verifier_nodeid_list))
    result = second_client.node.debug.setValidatorList(verifier_nodeid_list, pri_key)
    assert_code(result, 0)
    first_node.stop()

    second_client.economic.wait_consensus_blocknum(second_client.node)
    result = second_client.node.debug.setValidatorList(verifier_nodeid_list, pri_key)
    assert_code(result, 0)

    second_client.economic.wait_consensus_blocknum(second_client.node)
    result = second_client.node.debug.setValidatorList(verifier_nodeid_list, pri_key)
    assert_code(result, 0)
    result = second_client.node.debug.getWaitSlashingNodeList()
    log.info("Slashing Node List:{}".format(result))
    wait_slashing_list = get_slash_count(result, first_node.node_id)
    assert wait_slashing_list == 0

    second_client.economic.wait_consensus_blocknum(second_client.node)
    result = second_client.node.debug.setValidatorList(verifier_nodeid_list, pri_key)
    assert_code(result, 0)
    result = second_client.node.debug.getWaitSlashingNodeList()
    log.info("Slashing Node List:{}".format(result))
    wait_slashing_list = get_slash_count(result, first_node.node_id)
    assert wait_slashing_list == to_int(1)

    second_client.economic.wait_consensus_blocknum(second_client.node)
    result = second_client.node.debug.getWaitSlashingNodeList()
    log.info("Slashing Node List:{}".format(result))
    wait_slashing_list = get_slash_count(result, first_node.node_id)
    assert wait_slashing_list == to_int(11)

    second_client.economic.wait_consensus_blocknum(second_client.node)
    result = second_client.node.debug.getWaitSlashingNodeList()
    log.info("Slashing Node List:{}".format(result))
    wait_slashing_list = get_slash_count(result, first_node.node_id)
    assert wait_slashing_list == to_int(111)

    second_client.economic.wait_consensus_blocknum(second_client.node)
    result = second_client.node.debug.getWaitSlashingNodeList()
    log.info("Slashing Node List:{}".format(result))
    wait_slashing_list = get_slash_count(result, first_node.node_id)
    assert wait_slashing_list == 0
    result = second_client.ppos.getCandidateInfo(first_node.node_id)
    log.info("Candidate Info:{}".format(result))
    assert result['Ret']['Shares'] != 1000000000000000000000000


def test_ZB_NP_36(new_genesis_env, clients_noconsensus):
    """
    共识轮阈值大值变小值提案未生效时，节点零出块率处罚（存在零出块记录）
    """
    update_zero_produce(new_genesis_env)
    first_client = clients_noconsensus[0]
    second_client = clients_noconsensus[1]
    first_node = first_client.node
    staking_address, pri_key = create_pledge_node_information(first_client)

    # Start execution use case
    verifier_nodeid_list = first_client.economic.env.consensus_node_id_list()
    verifier_nodeid_list.append(first_node.node_id)
    verifier_nodeid_list.pop(0)
    log.info("verifier nodeid list: {}".format(verifier_nodeid_list))
    result = second_client.node.debug.setValidatorList(verifier_nodeid_list, pri_key)
    assert_code(result, 0)
    first_node.stop()

    second_client.economic.wait_consensus_blocknum(second_client.node)
    result = second_client.node.debug.setValidatorList(verifier_nodeid_list, pri_key)
    assert_code(result, 0)

    second_client.economic.wait_consensus_blocknum(second_client.node)
    result = second_client.node.debug.setValidatorList(verifier_nodeid_list, pri_key)
    assert_code(result, 0)
    result = second_client.node.debug.getWaitSlashingNodeList()
    log.info("Slashing Node List:{}".format(result))
    wait_slashing_list = get_slash_count(result, first_node.node_id)
    assert wait_slashing_list == 0

    second_client.economic.wait_consensus_blocknum(second_client.node)
    result = second_client.node.debug.setValidatorList(verifier_nodeid_list, pri_key)
    assert_code(result, 0)
    result = second_client.node.debug.getWaitSlashingNodeList()
    log.info("Slashing Node List:{}".format(result))
    wait_slashing_list = get_slash_count(result, first_node.node_id)
    assert wait_slashing_list == to_int(1)

    second_client.economic.wait_consensus_blocknum(second_client.node)
    result = second_client.node.debug.getWaitSlashingNodeList()
    log.info("Slashing Node List:{}".format(result))
    wait_slashing_list = get_slash_count(result, first_node.node_id)
    assert wait_slashing_list == to_int(11)
    # Voting failed and expired
    param_governance_verify_before_endblock(second_client, 'slashing', 'zeroProduceCumulativeTime', '3', True)

    second_client.economic.wait_consensus_blocknum(second_client.node)
    result = second_client.node.debug.getWaitSlashingNodeList()
    log.info("Slashing Node List:{}".format(result))
    wait_slashing_list = get_slash_count(result, first_node.node_id)
    assert wait_slashing_list == to_int(111)

    second_client.economic.wait_consensus_blocknum(second_client.node)
    result = second_client.node.debug.getWaitSlashingNodeList()
    log.info("Slashing Node List:{}".format(result))
    wait_slashing_list = get_slash_count(result, first_node.node_id)
    assert wait_slashing_list == 0
    result = second_client.ppos.getCandidateInfo(first_node.node_id)
    log.info("Candidate Info:{}".format(result))
    assert result['Ret']['Shares'] != 1000000000000000000000000


def test_ZB_NP_37(new_genesis_env, clients_noconsensus):
    """
    零出块率次数大值变小值提案未生效时，节点零出块率处罚（不存在零出块记录）
    """
    update_zero_produce(new_genesis_env)
    first_client = clients_noconsensus[0]
    second_client = clients_noconsensus[1]
    first_node = first_client.node
    staking_address, pri_key = create_pledge_node_information(first_client)

    # Voting failed and expired
    param_governance_verify_before_endblock(first_client, 'slashing', 'zeroProduceNumberThreshold', '2', True)

    # Start execution use case
    log.info("current validator: {}".format(get_pledge_list(second_client.ppos.getValidatorList)))
    verifier_nodeid_list = first_client.economic.env.consensus_node_id_list()
    verifier_nodeid_list.append(first_node.node_id)
    verifier_nodeid_list.pop(0)
    log.info("verifier nodeid list: {}".format(verifier_nodeid_list))
    result = second_client.node.debug.setValidatorList(verifier_nodeid_list, pri_key)
    assert_code(result, 0)
    first_node.stop()

    second_client.economic.wait_consensus_blocknum(second_client.node)
    log.info("current validator: {}".format(get_pledge_list(second_client.ppos.getValidatorList)))
    result = second_client.node.debug.setValidatorList(second_client.economic.env.consensus_node_id_list(), pri_key)
    assert_code(result, 0)

    second_client.economic.wait_consensus_blocknum(second_client.node)
    log.info("current validator: {}".format(get_pledge_list(second_client.ppos.getValidatorList)))
    result = second_client.node.debug.setValidatorList(second_client.economic.env.consensus_node_id_list(), pri_key)
    assert_code(result, 0)
    result = second_client.node.debug.getWaitSlashingNodeList()
    log.info("Slashing Node List:{}".format(result))
    wait_slashing_list = get_slash_count(result, first_node.node_id)
    assert wait_slashing_list == 0

    second_client.economic.wait_consensus_blocknum(second_client.node)
    log.info("current validator: {}".format(get_pledge_list(second_client.ppos.getValidatorList)))
    result = second_client.node.debug.setValidatorList(verifier_nodeid_list, pri_key)
    assert_code(result, 0)
    result = second_client.node.debug.getWaitSlashingNodeList()
    log.info("Slashing Node List:{}".format(result))
    wait_slashing_list = get_slash_count(result, first_node.node_id)
    assert wait_slashing_list == to_int(1)
    # Voting failed and expired
    param_governance_verify_before_endblock(first_client, 'slashing', 'zeroProduceNumberThreshold', '2', True)

    second_client.economic.wait_consensus_blocknum(second_client.node)
    log.info("current validator: {}".format(get_pledge_list(second_client.ppos.getValidatorList)))
    result = second_client.node.debug.setValidatorList(verifier_nodeid_list, pri_key)
    assert_code(result, 0)
    result = second_client.node.debug.getWaitSlashingNodeList()
    log.info("Slashing Node List:{}".format(result))
    wait_slashing_list = get_slash_count(result, first_node.node_id)
    assert wait_slashing_list == to_int(1)

    second_client.economic.wait_consensus_blocknum(second_client.node)
    log.info("current validator: {}".format(get_pledge_list(second_client.ppos.getValidatorList)))
    result = second_client.node.debug.getWaitSlashingNodeList()
    log.info("Slashing Node List:{}".format(result))
    wait_slashing_list = get_slash_count(result, first_node.node_id)
    assert wait_slashing_list == to_int(1)

    second_client.economic.wait_consensus_blocknum(second_client.node)
    log.info("current validator: {}".format(get_pledge_list(second_client.ppos.getValidatorList)))
    result = second_client.node.debug.getWaitSlashingNodeList()
    log.info("Slashing Node List:{}".format(result))
    wait_slashing_list = get_slash_count(result, first_node.node_id)
    assert wait_slashing_list == to_int(1001)


def test_ZB_NP_38(new_genesis_env, clients_noconsensus):
    """
    零出块率次数大值变小值提案未生效时，节点零出块率处罚（存在零出块记录）
    """
    update_zero_produce(new_genesis_env)
    first_client = clients_noconsensus[0]
    second_client = clients_noconsensus[1]
    first_node = first_client.node
    staking_address, pri_key = create_pledge_node_information(first_client)

    log.info("current validator: {}".format(get_pledge_list(second_client.ppos.getValidatorList)))
    verifier_nodeid_list = first_client.economic.env.consensus_node_id_list()
    verifier_nodeid_list.append(first_node.node_id)
    verifier_nodeid_list.pop(0)
    log.info("verifier nodeid list: {}".format(verifier_nodeid_list))
    result = second_client.node.debug.setValidatorList(verifier_nodeid_list, pri_key)
    assert_code(result, 0)
    first_node.stop()

    second_client.economic.wait_consensus_blocknum(second_client.node)
    log.info("current validator: {}".format(get_pledge_list(second_client.ppos.getValidatorList)))
    result = second_client.node.debug.setValidatorList(second_client.economic.env.consensus_node_id_list(), pri_key)
    assert_code(result, 0)

    second_client.economic.wait_consensus_blocknum(second_client.node)
    log.info("current validator: {}".format(get_pledge_list(second_client.ppos.getValidatorList)))
    result = second_client.node.debug.setValidatorList(second_client.economic.env.consensus_node_id_list(), pri_key)
    assert_code(result, 0)
    result = second_client.node.debug.getWaitSlashingNodeList()
    log.info("Slashing Node List:{}".format(result))
    wait_slashing_list = get_slash_count(result, first_node.node_id)
    assert wait_slashing_list == 0

    second_client.economic.wait_consensus_blocknum(second_client.node)
    log.info("current validator: {}".format(get_pledge_list(second_client.ppos.getValidatorList)))
    result = second_client.node.debug.setValidatorList(verifier_nodeid_list, pri_key)
    assert_code(result, 0)
    result = second_client.node.debug.getWaitSlashingNodeList()
    log.info("Slashing Node List:{}".format(result))
    wait_slashing_list = get_slash_count(result, first_node.node_id)
    assert wait_slashing_list == to_int(1)
    # Voting failed and expired
    param_governance_verify_before_endblock(first_client, 'slashing', 'zeroProduceNumberThreshold', '2', True)

    second_client.economic.wait_consensus_blocknum(second_client.node)
    log.info("current validator: {}".format(get_pledge_list(second_client.ppos.getValidatorList)))
    result = second_client.node.debug.setValidatorList(verifier_nodeid_list, pri_key)
    assert_code(result, 0)
    result = second_client.node.debug.getWaitSlashingNodeList()
    log.info("Slashing Node List:{}".format(result))
    wait_slashing_list = get_slash_count(result, first_node.node_id)
    assert wait_slashing_list == to_int(1)

    second_client.economic.wait_consensus_blocknum(second_client.node)
    log.info("current validator: {}".format(get_pledge_list(second_client.ppos.getValidatorList)))
    result = second_client.node.debug.getWaitSlashingNodeList()
    log.info("Slashing Node List:{}".format(result))
    wait_slashing_list = get_slash_count(result, first_node.node_id)
    assert wait_slashing_list == to_int(1)

    second_client.economic.wait_consensus_blocknum(second_client.node)
    log.info("current validator: {}".format(get_pledge_list(second_client.ppos.getValidatorList)))
    result = second_client.node.debug.getWaitSlashingNodeList()
    log.info("Slashing Node List:{}".format(result))
    wait_slashing_list = get_slash_count(result, first_node.node_id)
    assert wait_slashing_list == to_int(1001)


def test_ZB_NP_39(new_genesis_env, clients_noconsensus):
    """
    零出块率次数小值变大值提案未生效时，节点零出块率处罚（不存在零出块记录）
    """
    update_zero_produce(new_genesis_env)
    first_client = clients_noconsensus[0]
    second_client = clients_noconsensus[1]
    first_node = first_client.node
    staking_address, pri_key = create_pledge_node_information(first_client)

    # Voting failed and expired
    param_governance_verify_before_endblock(first_client, 'slashing', 'zeroProduceNumberThreshold', '4', True)

    # Start execution use case
    log.info("current validator: {}".format(get_pledge_list(second_client.ppos.getValidatorList)))
    verifier_nodeid_list = first_client.economic.env.consensus_node_id_list()
    verifier_nodeid_list.append(first_node.node_id)
    verifier_nodeid_list.pop(0)
    log.info("verifier nodeid list: {}".format(verifier_nodeid_list))
    result = second_client.node.debug.setValidatorList(verifier_nodeid_list, pri_key)
    assert_code(result, 0)
    first_node.stop()

    second_client.economic.wait_consensus_blocknum(second_client.node)
    log.info("current validator: {}".format(get_pledge_list(second_client.ppos.getValidatorList)))
    result = second_client.node.debug.setValidatorList(second_client.economic.env.consensus_node_id_list(), pri_key)
    assert_code(result, 0)

    second_client.economic.wait_consensus_blocknum(second_client.node)
    log.info("current validator: {}".format(get_pledge_list(second_client.ppos.getValidatorList)))
    result = second_client.node.debug.setValidatorList(verifier_nodeid_list, pri_key)
    assert_code(result, 0)
    result = second_client.node.debug.getWaitSlashingNodeList()
    log.info("Slashing Node List:{}".format(result))
    wait_slashing_list = get_slash_count(result, first_node.node_id)
    assert wait_slashing_list == 0

    second_client.economic.wait_consensus_blocknum(second_client.node)
    log.info("current validator: {}".format(get_pledge_list(second_client.ppos.getValidatorList)))
    result = second_client.node.debug.setValidatorList(verifier_nodeid_list, pri_key)
    assert_code(result, 0)
    result = second_client.node.debug.getWaitSlashingNodeList()
    log.info("Slashing Node List:{}".format(result))
    wait_slashing_list = get_slash_count(result, first_node.node_id)
    assert wait_slashing_list == to_int(1)

    second_client.economic.wait_consensus_blocknum(second_client.node)
    log.info("current validator: {}".format(get_pledge_list(second_client.ppos.getValidatorList)))
    result = second_client.node.debug.setValidatorList(verifier_nodeid_list, pri_key)
    assert_code(result, 0)
    result = second_client.node.debug.getWaitSlashingNodeList()
    log.info("Slashing Node List:{}".format(result))
    wait_slashing_list = get_slash_count(result, first_node.node_id)
    assert wait_slashing_list == to_int(1)

    second_client.economic.wait_consensus_blocknum(second_client.node)
    log.info("current validator: {}".format(get_pledge_list(second_client.ppos.getValidatorList)))
    result = second_client.node.debug.getWaitSlashingNodeList()
    log.info("Slashing Node List:{}".format(result))
    wait_slashing_list = get_slash_count(result, first_node.node_id)
    assert wait_slashing_list == to_int(101)

    second_client.economic.wait_consensus_blocknum(second_client.node)
    result = second_client.node.debug.getWaitSlashingNodeList()
    log.info("Slashing Node List:{}".format(result))
    wait_slashing_list = get_slash_count(result, first_node.node_id)
    assert wait_slashing_list == 0
    result = second_client.ppos.getCandidateInfo(first_node.node_id)
    log.info("Candidate Info:{}".format(result))
    assert result['Ret']['Shares'] != 1000000000000000000000000


def test_ZB_NP_40(new_genesis_env, clients_noconsensus):
    """
    零出块率次数小值变大值提案未生效时，节点零出块率处罚（存在零出块记录）
    """
    update_zero_produce(new_genesis_env)
    first_client = clients_noconsensus[0]
    second_client = clients_noconsensus[1]
    first_node = first_client.node
    staking_address, pri_key = create_pledge_node_information(first_client)

    # Start execution use case
    log.info("current validator: {}".format(get_pledge_list(second_client.ppos.getValidatorList)))
    verifier_nodeid_list = first_client.economic.env.consensus_node_id_list()
    verifier_nodeid_list.append(first_node.node_id)
    verifier_nodeid_list.pop(0)
    log.info("verifier nodeid list: {}".format(verifier_nodeid_list))
    result = second_client.node.debug.setValidatorList(verifier_nodeid_list, pri_key)
    assert_code(result, 0)
    first_node.stop()

    second_client.economic.wait_consensus_blocknum(second_client.node)
    log.info("current validator: {}".format(get_pledge_list(second_client.ppos.getValidatorList)))
    result = second_client.node.debug.setValidatorList(verifier_nodeid_list, pri_key)
    assert_code(result, 0)

    second_client.economic.wait_consensus_blocknum(second_client.node)
    log.info("current validator: {}".format(get_pledge_list(second_client.ppos.getValidatorList)))
    result = second_client.node.debug.setValidatorList(verifier_nodeid_list, pri_key)
    assert_code(result, 0)
    result = second_client.node.debug.getWaitSlashingNodeList()
    log.info("Slashing Node List:{}".format(result))
    wait_slashing_list = get_slash_count(result, first_node.node_id)
    assert wait_slashing_list == 0
    # Voting failed and expired
    param_governance_verify_before_endblock(first_client, 'slashing', 'zeroProduceNumberThreshold', '4', True)

    second_client.economic.wait_consensus_blocknum(second_client.node)
    log.info("current validator: {}".format(get_pledge_list(second_client.ppos.getValidatorList)))
    result = second_client.node.debug.setValidatorList(second_client.economic.env.consensus_node_id_list(), pri_key)
    assert_code(result, 0)
    result = second_client.node.debug.getWaitSlashingNodeList()
    log.info("Slashing Node List:{}".format(result))
    wait_slashing_list = get_slash_count(result, first_node.node_id)
    assert wait_slashing_list == to_int(1)

    second_client.economic.wait_consensus_blocknum(second_client.node)
    log.info("current validator: {}".format(get_pledge_list(second_client.ppos.getValidatorList)))
    assert_code(result, 0)
    result = second_client.node.debug.getWaitSlashingNodeList()
    log.info("Slashing Node List:{}".format(result))
    wait_slashing_list = get_slash_count(result, first_node.node_id)
    assert wait_slashing_list == to_int(11)

    second_client.economic.wait_consensus_blocknum(second_client.node)
    log.info("current validator: {}".format(get_pledge_list(second_client.ppos.getValidatorList)))
    assert_code(result, 0)
    result = second_client.node.debug.getWaitSlashingNodeList()
    log.info("Slashing Node List:{}".format(result))
    wait_slashing_list = get_slash_count(result, first_node.node_id)
    assert wait_slashing_list == to_int(111)

    second_client.economic.wait_consensus_blocknum(second_client.node)
    result = second_client.node.debug.getWaitSlashingNodeList()
    log.info("Slashing Node List:{}".format(result))
    wait_slashing_list = get_slash_count(result, first_node.node_id)
    assert wait_slashing_list == 0
    result = second_client.ppos.getCandidateInfo(first_node.node_id)
    log.info("Candidate Info:{}".format(result))
    assert result['Ret']['Shares'] != 1000000000000000000000000