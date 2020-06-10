import time
from decimal import Decimal

from dacite import from_dict
from dacite import from_dict
import copy
from common.key import mock_duplicate_sign
from common.log import log
from tests.conftest import param_governance_verify_before_endblock, param_governance_verify
from tests.lib import Genesis, check_node_in_list, assert_code, von_amount, get_pledge_list, wait_block_number


def update_zero_produce(new_genesis_env, cumulativetime=4, numberthreshold=3):
    global pri_key
    pri_key = new_genesis_env.account.account_with_money["prikey"]
    genesis = from_dict(data_class=Genesis, data=new_genesis_env.genesis_config)
    genesis.economicModel.slashing.zeroProduceCumulativeTime = cumulativetime
    genesis.economicModel.slashing.zeroProduceNumberThreshold = numberthreshold
    new_file = new_genesis_env.cfg.env_tmp + "/genesis_0.13.0.json"
    genesis.to_file(new_file)
    new_genesis_env.deploy_all(new_file)


def assert_slashing(candidate_info, staking_amount):
    ret = candidate_info["Ret"]
    assert ret["Status"] != 0
    assert ret["Released"] < staking_amount


def assert_not_slashing(candidate_info, staking_amount):
    ret = candidate_info["Ret"]
    assert ret["Status"] == 0
    assert ret["Released"] == staking_amount


def assert_set_validator_list(node, validator_list):
    global pri_key
    result = node.debug.setValidatorList(validator_list, pri_key)
    if result:
        raise Exception("setValidatorList:{}".format(result))


def get_slash_count(wait_slashing_node_list, node_id):
    for slashing in wait_slashing_node_list:
        if slashing["NodeId"] == node_id:
            return slashing["CountBit"]
    return 0


def gen_validator_list(initial_validator, slashing_node_id):
    slashing_node_list = [node_id for node_id in initial_validator]
    slashing_node_list.append(slashing_node_id)
    slashing_node_list.pop(0)
    return initial_validator, slashing_node_list


def to_int(value):
    return int(str(value), 2)


def set_slashing(initial_validator_list, contain_slashing_list, node, economic, status):
    do = {
        "1": contain_slashing_list,
        "0": initial_validator_list
    }
    status = str(status)
    for state in status:
        assert_set_validator_list(node, do[state])
        economic.wait_consensus_blocknum(node)


def create_pledge_node_information(client):
    log.info("Current connection node：{}".format(client.node.node_mark))
    log.info("Start creating a pledge account Pledge_address")
    staking_amount = von_amount(client.economic.create_staking_limit, 2)
    staking_address, pri_key = client.economic.account.generate_account(client.node.web3, staking_amount)
    log.info("Created, account address: {} Amount: {}".format(staking_address, staking_amount))
    log.info("Start applying for a pledge node")
    result = client.staking.create_staking(0, staking_address, staking_address)
    assert_code(result, 0)
    client.economic.wait_settlement_blocknum(client.node)
    log.info("Current block height: {}".format(client.node.eth.blockNumber))
    result = client.node.ppos.getVerifierList()
    log.info("current Verifier List：{}".format(result))
    return staking_address, pri_key


def test_ZB_NP_01(new_genesis_env, client_noconsensus):
    """
    节点未被选中验证人列表查询零出块记录表
    """
    # Change configuration parameters
    genesis = from_dict(data_class=Genesis, data=new_genesis_env.genesis_config)
    genesis.economicModel.staking.maxValidators = 4
    new_file = new_genesis_env.cfg.env_tmp + "/genesis_0.13.0.json"
    genesis.to_file(new_file)
    new_genesis_env.deploy_all(new_file)
    # start execution use case
    create_pledge_node_information(client_noconsensus)
    # view the zero out block record table
    wait_slashing_list = client_noconsensus.node.debug.getWaitSlashingNodeList()
    assert get_slash_count(wait_slashing_list, client_noconsensus.node.node_id) == 0


def test_ZB_NP_02(client_noconsensus):
    """
    节点未被选中共识验证人列表查询零出块记录表（不存在零出块记录）
    """
    # start execution use case
    create_pledge_node_information(client_noconsensus)
    result = check_node_in_list(client_noconsensus.node.node_id, client_noconsensus.ppos.getValidatorList)
    log.info("Current node in consensus list status：{}".format(result))
    if result is False:
        # view the zero out block record table
        wait_slashing_list = client_noconsensus.node.debug.getWaitSlashingNodeList()
        assert get_slash_count(wait_slashing_list, client_noconsensus.node.node_id) == 0
    else:
        client_noconsensus.economic.wait_consensus_blocknum(client_noconsensus.node)


def test_ZB_NP_03(new_genesis_env, client_noconsensus):
    """
    节点未被选中共识验证人列表查询零出块记录表（存在零出块记录）
    """
    update_zero_produce(new_genesis_env)
    economic = client_noconsensus.economic
    node = economic.env.get_consensus_node_by_index(0)
    # start execution use case
    create_pledge_node_information(client_noconsensus)
    current_validator = get_pledge_list(node.ppos.getValidatorList)
    log.info("current validator: {}".format(current_validator))
    initial_validator, slashing_node_list = gen_validator_list(economic.env.consensus_node_id_list(), client_noconsensus.node.node_id)
    log.info("verifier nodeid list: {}".format(slashing_node_list))
    assert_set_validator_list(node, slashing_node_list)
    # stop node
    client_noconsensus.node.stop()

    economic.wait_consensus_blocknum(node)
    log.info("Current block height: {}".format(node.eth.blockNumber))
    current_validator = get_pledge_list(node.ppos.getValidatorList)
    log.info("current validator: {}".format(current_validator))
    assert_set_validator_list(node, initial_validator)

    economic.wait_consensus_blocknum(node)
    log.info("Current block height: {}".format(node.eth.blockNumber))
    current_validator = get_pledge_list(node.ppos.getValidatorList)
    log.info("current validator: {}".format(current_validator))
    wait_slashing_list = node.debug.getWaitSlashingNodeList()
    log.info("Zero block record table：{}".format(wait_slashing_list))
    assert_set_validator_list(node, initial_validator)
    assert get_slash_count(wait_slashing_list, client_noconsensus.node.node_id) == 0

    economic.wait_consensus_blocknum(node)
    log.info("Current block height: {}".format(node.eth.blockNumber))
    current_validator = get_pledge_list(node.ppos.getValidatorList)
    log.info("current validator: {}".format(current_validator))
    wait_slashing_list = node.debug.getWaitSlashingNodeList()
    log.info("Zero block record table：{}".format(wait_slashing_list))
    assert_set_validator_list(node, initial_validator)
    assert get_slash_count(wait_slashing_list, client_noconsensus.node.node_id) == to_int(1)

    economic.wait_consensus_blocknum(node)
    log.info("Current block height: {}".format(node.eth.blockNumber))
    current_validator = get_pledge_list(node.ppos.getValidatorList)
    log.info("current validator: {}".format(current_validator))
    wait_slashing_list = node.debug.getWaitSlashingNodeList()
    log.info("Zero block record table：{}".format(wait_slashing_list))
    assert_set_validator_list(node, initial_validator)
    assert get_slash_count(wait_slashing_list, client_noconsensus.node.node_id) == to_int(1)


def test_ZB_NP_04(client_noconsensus):
    """
    节点被选中共识验证人查询零出块记录表（不存在零出块记录）
    """
    economic = client_noconsensus.economic
    node = economic.env.get_consensus_node_by_index(0)
    create_pledge_node_information(client_noconsensus)
    for i in range(4):
        result = check_node_in_list(client_noconsensus.node.node_id, node.ppos.getValidatorList)
        log.info("Current node in consensus list status：{}".format(result))
        if result:
            # view the zero out block record table
            wait_slashing_list = node.debug.getWaitSlashingNodeList()
            log.info("Zero block record table：{}".format(wait_slashing_list))
            assert get_slash_count(wait_slashing_list, client_noconsensus.node.node_id) == 0
            break
        else:
            economic.wait_consensus_blocknum(node)


def test_ZB_NP_05(new_genesis_env, client_noconsensus):
    """
    节点被选中共识验证人查询零出块记录表（存在零出块记录）
    """
    update_zero_produce(new_genesis_env)
    economic = client_noconsensus.economic
    node = economic.env.get_consensus_node_by_index(0)
    create_pledge_node_information(client_noconsensus)
    initial_validator, slashing_node_list = gen_validator_list(economic.env.consensus_node_id_list(), client_noconsensus.node.node_id)
    log.info("verifier nodeid list: {}".format(slashing_node_list))
    assert_set_validator_list(node, slashing_node_list)
    client_noconsensus.node.stop()

    economic.wait_consensus_blocknum(node)
    log.info("Current block height: {}".format(node.eth.blockNumber))
    current_validator = get_pledge_list(node.ppos.getValidatorList)
    log.info("current validator: {}".format(current_validator))
    assert_set_validator_list(node, slashing_node_list)

    economic.wait_consensus_blocknum(node)
    log.info("Current block height: {}".format(node.eth.blockNumber))
    current_validator = get_pledge_list(node.ppos.getValidatorList)
    log.info("current validator: {}".format(current_validator))
    assert_set_validator_list(node, slashing_node_list)
    wait_slashing_list = node.debug.getWaitSlashingNodeList()
    log.info("Zero block record table：{}".format(wait_slashing_list))
    assert get_slash_count(wait_slashing_list, client_noconsensus.node.node_id) == 0
    client_noconsensus.node.start(False)

    economic.wait_consensus_blocknum(node)
    log.info("Current block height: {}".format(node.eth.blockNumber))
    current_validator = get_pledge_list(node.ppos.getValidatorList)
    log.info("current validator: {}".format(current_validator))
    wait_slashing_list = node.debug.getWaitSlashingNodeList()
    log.info("Zero block record table：{}".format(wait_slashing_list))
    assert get_slash_count(wait_slashing_list, client_noconsensus.node.node_id) == to_int(1)

    economic.wait_consensus_blocknum(node)
    log.info("Current block height: {}".format(node.eth.blockNumber))
    current_validator = get_pledge_list(node.ppos.getValidatorList)
    log.info("current validator: {}".format(current_validator))
    wait_slashing_list = node.debug.getWaitSlashingNodeList()
    log.info("Zero block record table：{}".format(wait_slashing_list))
    assert get_slash_count(wait_slashing_list, client_noconsensus.node.node_id) == 0
    assert_not_slashing(node.ppos.getCandidateInfo(client_noconsensus.node.node_id), economic.create_staking_limit)


def test_ZB_NP_06_07(new_genesis_env, client_noconsensus):
    """
    ZB_NP_06:节点被选中共识验证人未出块查询零出块记录表（不存在零出块记录）
    ZB_NP_07:节点被选中共识验证人未出块查询零出块记录表（存在零出块记录）
    """
    update_zero_produce(new_genesis_env)
    economic = client_noconsensus.economic
    node = economic.env.get_consensus_node_by_index(0)
    create_pledge_node_information(client_noconsensus)
    initial_validator, slashing_node_list = gen_validator_list(economic.env.consensus_node_id_list(), client_noconsensus.node.node_id)
    log.info("verifier nodeid list: {}".format(slashing_node_list))
    assert_set_validator_list(node, slashing_node_list)
    client_noconsensus.node.stop()

    economic.wait_consensus_blocknum(node)
    log.info("Current block height: {}".format(node.eth.blockNumber))
    current_validator = get_pledge_list(node.ppos.getValidatorList)
    log.info("current validator: {}".format(current_validator))
    assert_set_validator_list(node, slashing_node_list)

    economic.wait_consensus_blocknum(node)
    log.info("Current block height: {}".format(node.eth.blockNumber))
    current_validator = get_pledge_list(node.ppos.getValidatorList)
    log.info("current validator: {}".format(current_validator))
    wait_slashing_list = node.debug.getWaitSlashingNodeList()
    log.info("Zero block record table：{}".format(wait_slashing_list))
    assert get_slash_count(wait_slashing_list, client_noconsensus.node.node_id) == 0

    economic.wait_consensus_blocknum(node)
    log.info("Current block height: {}".format(node.eth.blockNumber))
    current_validator = get_pledge_list(node.ppos.getValidatorList)
    log.info("current validator: {}".format(current_validator))
    wait_slashing_list = node.debug.getWaitSlashingNodeList()
    log.info("Zero block record table：{}".format(wait_slashing_list))
    assert get_slash_count(wait_slashing_list, client_noconsensus.node.node_id) == to_int(1)

    economic.wait_consensus_blocknum(node)
    log.info("Current block height: {}".format(node.eth.blockNumber))
    current_validator = get_pledge_list(node.ppos.getValidatorList)
    log.info("current validator: {}".format(current_validator))
    wait_slashing_list = node.debug.getWaitSlashingNodeList()
    log.info("Zero block record table：{}".format(wait_slashing_list))
    assert get_slash_count(wait_slashing_list, client_noconsensus.node.node_id) == to_int(11)


def test_ZB_NP_08(new_genesis_env, client_noconsensus):
    """
    节点达到共识轮阈值，未达到累计次数被选中共识验证人未出块
    """
    update_zero_produce(new_genesis_env)
    economic = client_noconsensus.economic
    node = economic.env.get_consensus_node_by_index(0)
    create_pledge_node_information(client_noconsensus)
    initial_validator, slashing_node_list = gen_validator_list(economic.env.consensus_node_id_list(), client_noconsensus.node.node_id)
    log.info("verifier nodeid list: {}".format(slashing_node_list))
    assert_set_validator_list(node, slashing_node_list)
    client_noconsensus.node.stop()

    set_slashing(initial_validator, slashing_node_list, node, economic, "1001")
    assert_set_validator_list(node, slashing_node_list)

    economic.wait_consensus_blocknum(node, 1)
    log.info("Current block height: {}".format(node.eth.blockNumber))
    current_validator = get_pledge_list(node.ppos.getValidatorList)
    log.info("current validator: {}".format(current_validator))
    wait_slashing_list = node.debug.getWaitSlashingNodeList()
    log.info("Zero block record table：{}".format(wait_slashing_list))
    assert get_slash_count(wait_slashing_list, client_noconsensus.node.node_id) == to_int(1001)

    economic.wait_consensus_blocknum(node)
    log.info("Current block height: {}".format(node.eth.blockNumber))
    current_validator = get_pledge_list(node.ppos.getValidatorList)
    log.info("current validator: {}".format(current_validator))
    wait_slashing_list = node.debug.getWaitSlashingNodeList()
    log.info("Zero block record table：{}".format(wait_slashing_list))
    assert get_slash_count(wait_slashing_list, client_noconsensus.node.node_id) == to_int(11)
    assert_not_slashing(node.ppos.getCandidateInfo(client_noconsensus.node.node_id), economic.create_staking_limit)


def test_ZB_NP_09(new_genesis_env, client_noconsensus):
    """
    节点先双签被举报再触发零出块
    """
    # update_zero_produce(new_genesis_env, 1, 1)
    economic = client_noconsensus.economic
    node = economic.env.get_consensus_node_by_index(0)
    report_address, _ = economic.account.generate_account(node.web3, node.web3.toWei(1000, 'ether'))
    create_pledge_node_information(client_noconsensus)
    result = node.ppos.getCandidateInfo(client_noconsensus.node.node_id)
    log.info("Candidate Info:{}".format(result))
    amount_of_pledge = result['Ret']['Released']
    block_reward, staking_reward = economic.get_current_year_reward(node)
    log.info("block_reward: {} staking_reward: {}".format(block_reward, staking_reward))
    economic.wait_consensus_blocknum(node)
    log.info("Current block height: {}".format(node.eth.blockNumber))
    report_information = mock_duplicate_sign(1, client_noconsensus.node.nodekey, client_noconsensus.node.blsprikey, node.eth.blockNumber)
    log.info("Report information: {}".format(report_information))
    result = client_noconsensus.duplicatesign.reportDuplicateSign(1, report_information, report_address)
    assert_code(result, 0)
    client_noconsensus.node.stop()

    economic.wait_consensus_blocknum(node)
    log.info("Current block height: {}".format(node.eth.blockNumber))
    current_validator = get_pledge_list(node.ppos.getValidatorList)
    log.info("current validator: {}".format(current_validator))
    wait_slashing_list = node.debug.getWaitSlashingNodeList()
    log.info("Zero block record table：{}".format(wait_slashing_list))
    assert get_slash_count(wait_slashing_list, client_noconsensus.node.node_id) == 0
    penalty_for_double_signing = int(Decimal(str(amount_of_pledge)) * Decimal(str(economic.genesis.economicModel.slashing.slashFractionDuplicateSign / 10000)))
    log.info("penalty_for_double_signing: {}".format(penalty_for_double_signing))
    economic.wait_consensus_blocknum(node)
    log.info("Current block height: {}".format(node.eth.blockNumber))
    current_validator = get_pledge_list(node.ppos.getValidatorList)
    log.info("current validator: {}".format(current_validator))
    wait_slashing_list = node.debug.getWaitSlashingNodeList()
    log.info("Zero block record table：{}".format(wait_slashing_list))
    assert get_slash_count(wait_slashing_list, client_noconsensus.node.node_id) == 0
    assert_slashing(node.ppos.getCandidateInfo(client_noconsensus.node.node_id), economic.create_staking_limit)
    result = node.ppos.getCandidateInfo(client_noconsensus.node.node_id)
    log.info("Candidate Info:{}".format(result))
    punishment_amonut = int(Decimal(str(block_reward)) * int(Decimal(str(economic.genesis.economicModel.slashing.slashBlocksReward))))
    log.info("punishment_amonut:{}".format(punishment_amonut))
    assert result['Ret']['Released'] == amount_of_pledge - penalty_for_double_signing - punishment_amonut



def test_ZB_NP_10(new_genesis_env, clients_noconsensus):
    """
    节点触发零出块再双签被举报
    """
    update_zero_produce(new_genesis_env)
    economic = clients_noconsensus[0].economic
    second_client = clients_noconsensus[1]
    node = economic.env.get_consensus_node_by_index(0)
    report_address, _ = economic.account.generate_account(node.web3, node.web3.toWei(1000, 'ether'))
    create_pledge_node_information(clients_noconsensus[0])
    initial_validator, slashing_node_list = gen_validator_list(economic.env.consensus_node_id_list(), clients_noconsensus[0].node.node_id)
    log.info("verifier nodeid list: {}".format(slashing_node_list))
    assert_set_validator_list(node, slashing_node_list)
    clients_noconsensus[0].node.stop()

    economic.wait_consensus_blocknum(node)
    current_block = node.eth.blockNumber
    log.info("Current block height: {}".format(current_block))
    assert_set_validator_list(node, slashing_node_list)

    economic.wait_consensus_blocknum(node)
    log.info("Current block height: {}".format(node.eth.blockNumber))
    current_validator = get_pledge_list(node.ppos.getValidatorList)
    log.info("current validator: {}".format(current_validator))
    wait_slashing_list = node.debug.getWaitSlashingNodeList()
    log.info("Zero block record table：{}".format(wait_slashing_list))
    assert get_slash_count(wait_slashing_list, clients_noconsensus[0].node.node_id) == 0

    economic.wait_consensus_blocknum(node)
    log.info("Current block height: {}".format(node.eth.blockNumber))
    current_validator = get_pledge_list(node.ppos.getValidatorList)
    log.info("current validator: {}".format(current_validator))
    wait_slashing_list = node.debug.getWaitSlashingNodeList()
    log.info("Zero block record table：{}".format(wait_slashing_list))
    assert get_slash_count(wait_slashing_list, clients_noconsensus[0].node.node_id) == to_int(1)
    report_information = mock_duplicate_sign(1, clients_noconsensus[0].node.nodekey, clients_noconsensus[0].node.blsprikey, current_block)
    log.info("Report information: {}".format(report_information))
    result = second_client.duplicatesign.reportDuplicateSign(1, report_information, report_address)
    assert_code(result, 0)

    economic.wait_consensus_blocknum(node)
    log.info("Current block height: {}".format(node.eth.blockNumber))
    current_validator = get_pledge_list(node.ppos.getValidatorList)
    log.info("current validator: {}".format(current_validator))
    result = check_node_in_list(clients_noconsensus[0].node.node_id, node.ppos.getVerifierList)
    log.info("Current node in consensus list status：{}".format(result))
    assert result is False
    wait_slashing_list = node.debug.getWaitSlashingNodeList()
    log.info("Zero block record table：{}".format(wait_slashing_list))
    assert get_slash_count(wait_slashing_list, clients_noconsensus[0].node.node_id) == to_int(11)
    assert_slashing(node.ppos.getCandidateInfo(clients_noconsensus[0].node.node_id), economic.create_staking_limit)


def test_ZB_NP_31(new_genesis_env, client_consensus, client_noconsensus):
    """
    治理提案投票未通过，节点零出块率处罚（不存在零出块记录）
    """
    update_zero_produce(new_genesis_env, 3)
    economic = client_noconsensus.economic
    node = economic.env.get_consensus_node_by_index(0)
    # Voting failed and expired
    param_governance_verify(client_consensus, 'slashing', 'zeroProduceCumulativeTime', '4', False)
    # Start execution use case
    create_pledge_node_information(client_noconsensus)
    initial_validator, slashing_node_list = gen_validator_list(economic.env.consensus_node_id_list(), client_noconsensus.node.node_id)
    log.info("verifier nodeid list: {}".format(slashing_node_list))
    # assert_set_validator_list(node, slashing_node_list)
    client_noconsensus.node.stop()

    set_slashing(initial_validator, slashing_node_list, node, economic, "111")

    economic.wait_consensus_blocknum(node, 1)
    log.info("Current block height: {}".format(node.eth.blockNumber))
    current_validator = get_pledge_list(node.ppos.getValidatorList)
    log.info("current validator: {}".format(current_validator))
    assert_set_validator_list(node, slashing_node_list)
    wait_slashing_list = node.debug.getWaitSlashingNodeList()
    log.info("Zero block record table：{}".format(wait_slashing_list))
    assert get_slash_count(wait_slashing_list, client_noconsensus.node.node_id) == 0
    assert_slashing(node.ppos.getCandidateInfo(client_noconsensus.node.node_id), economic.create_staking_limit)


def test_ZB_NP_32(new_genesis_env, client_consensus, client_noconsensus):
    """
    治理提案投票未通过，节点零出块率处罚（存在零出块记录）
    """
    update_zero_produce(new_genesis_env, 3)
    economic = client_noconsensus.economic
    node = economic.env.get_consensus_node_by_index(0)
    # Voting failed and expired
    param_governance_verify_before_endblock(client_consensus, 'slashing', 'zeroProduceCumulativeTime', '4', False)
    # Start execution use case
    create_pledge_node_information(client_noconsensus)
    log.info("Current block height: {}".format(node.eth.blockNumber))
    initial_validator, slashing_node_list = gen_validator_list(economic.env.consensus_node_id_list(), client_noconsensus.node.node_id)
    log.info("verifier nodeid list: {}".format(initial_validator))
    client_noconsensus.node.stop()

    current_validator = get_pledge_list(node.ppos.getValidatorList)
    log.info("current validator: {}".format(current_validator))
    set_slashing(initial_validator, slashing_node_list, node, economic, "0111")

    economic.wait_consensus_blocknum(node, 1)
    log.info("Current block height: {}".format(node.eth.blockNumber))
    current_validator = get_pledge_list(node.ppos.getValidatorList)
    log.info("current validator: {}".format(current_validator))
    assert_set_validator_list(node, slashing_node_list)
    wait_slashing_list = node.debug.getWaitSlashingNodeList()
    log.info("Zero block record table：{}".format(wait_slashing_list))
    assert get_slash_count(wait_slashing_list, client_noconsensus.node.node_id) == 0
    assert_slashing(node.ppos.getCandidateInfo(client_noconsensus.node.node_id), economic.create_staking_limit)


def test_ZB_NP_33(new_genesis_env,client_consensus, client_noconsensus):
    """
    共识轮阈值小值变大值提案未生效时，节点零出块率处罚（不存在零出块记录）
    """
    update_zero_produce(new_genesis_env, 3)
    economic = client_noconsensus.economic
    node = economic.env.get_consensus_node_by_index(0)
    create_pledge_node_information(client_noconsensus)

    # Voting failed and expired
    param_governance_verify_before_endblock(client_consensus, 'slashing', 'zeroProduceCumulativeTime', '4', True)

    # Start execution use case
    initial_validator, slashing_node_list = gen_validator_list(economic.env.consensus_node_id_list(), client_noconsensus.node.node_id)
    log.info("verifier nodeid list: {}".format(initial_validator))
    client_noconsensus.node.stop()

    set_slashing(initial_validator, slashing_node_list, node, economic, "111")

    economic.wait_consensus_blocknum(node, 1)
    log.info("Current block height: {}".format(node.eth.blockNumber))
    current_validator = get_pledge_list(node.ppos.getValidatorList)
    log.info("current validator: {}".format(current_validator))
    assert_set_validator_list(node, slashing_node_list)
    wait_slashing_list = node.debug.getWaitSlashingNodeList()
    log.info("Zero block record table：{}".format(wait_slashing_list))
    assert get_slash_count(wait_slashing_list, client_noconsensus.node.node_id) == 0
    assert_slashing(node.ppos.getCandidateInfo(client_noconsensus.node.node_id), economic.create_staking_limit)


def test_ZB_NP_34(new_genesis_env, client_consensus, client_noconsensus):
    """
    共识轮阈值小值变大值提案未生效时，节点零出块率处罚（存在零出块记录）
    """
    update_zero_produce(new_genesis_env, 3)
    economic = client_noconsensus.economic
    node = economic.env.get_consensus_node_by_index(0)
    create_pledge_node_information(client_noconsensus)

    # Start execution use case
    initial_validator, slashing_node_list = gen_validator_list(economic.env.consensus_node_id_list(), client_noconsensus.node.node_id)
    log.info("verifier nodeid list: {}".format(initial_validator))
    assert_set_validator_list(node, slashing_node_list)
    client_noconsensus.node.stop()

    economic.wait_consensus_blocknum(node)
    assert_set_validator_list(node, slashing_node_list)

    economic.wait_consensus_blocknum(node)
    log.info("Current block height: {}".format(node.eth.blockNumber))
    current_validator = get_pledge_list(node.ppos.getValidatorList)
    log.info("current validator: {}".format(current_validator))
    assert_set_validator_list(node, slashing_node_list)
    wait_slashing_list = node.debug.getWaitSlashingNodeList()
    log.info("Zero block record table：{}".format(wait_slashing_list))
    assert get_slash_count(wait_slashing_list, client_noconsensus.node.node_id) == 0

    economic.wait_consensus_blocknum(node)
    log.info("Current block height: {}".format(node.eth.blockNumber))
    current_validator = get_pledge_list(node.ppos.getValidatorList)
    log.info("current validator: {}".format(current_validator))
    wait_slashing_list = node.debug.getWaitSlashingNodeList()
    log.info("Zero block record table：{}".format(wait_slashing_list))
    assert get_slash_count(wait_slashing_list, client_noconsensus.node.node_id) == to_int(1)
    # Voting failed and expired
    print(client_consensus.node.node_mark)
    param_governance_verify_before_endblock(client_consensus, 'slashing', 'zeroProduceCumulativeTime', '4', True)

    economic.wait_consensus_blocknum(node)
    log.info("Current block height: {}".format(node.eth.blockNumber))
    current_validator = get_pledge_list(node.ppos.getValidatorList)
    log.info("current validator: {}".format(current_validator))
    wait_slashing_list = node.debug.getWaitSlashingNodeList()
    log.info("Zero block record table：{}".format(wait_slashing_list))
    assert get_slash_count(wait_slashing_list, client_noconsensus.node.node_id) == to_int(11)

    economic.wait_consensus_blocknum(node)
    log.info("Current block height: {}".format(node.eth.blockNumber))
    current_validator = get_pledge_list(node.ppos.getValidatorList)
    log.info("current validator: {}".format(current_validator))
    wait_slashing_list = node.debug.getWaitSlashingNodeList()
    log.info("Zero block record table：{}".format(wait_slashing_list))
    assert get_slash_count(wait_slashing_list, client_noconsensus.node.node_id) == 0
    assert_slashing(node.ppos.getCandidateInfo(client_noconsensus.node.node_id), economic.create_staking_limit)


def test_ZB_NP_35(new_genesis_env, client_consensus,  client_noconsensus):
    """
    共识轮阈值大值变小值提案未生效时，节点零出块率处罚（不存在零出块记录）
    """
    update_zero_produce(new_genesis_env)
    economic = client_noconsensus.economic
    node = economic.env.get_consensus_node_by_index(0)
    create_pledge_node_information(client_noconsensus)

    # Voting failed and expired
    param_governance_verify_before_endblock(client_consensus, 'slashing', 'zeroProduceCumulativeTime', '3', True)

    # Start execution use case
    initial_validator, slashing_node_list = gen_validator_list(economic.env.consensus_node_id_list(), client_noconsensus.node.node_id)
    log.info("verifier nodeid list: {}".format(initial_validator))
    client_noconsensus.node.stop()

    set_slashing(initial_validator, slashing_node_list, node, economic, "1101")

    economic.wait_consensus_blocknum(node, 1)
    log.info("Current block height: {}".format(node.eth.blockNumber))
    current_validator = get_pledge_list(node.ppos.getValidatorList)
    log.info("current validator: {}".format(current_validator))
    wait_slashing_list = node.debug.getWaitSlashingNodeList()
    log.info("Zero block record table：{}".format(wait_slashing_list))
    assert get_slash_count(wait_slashing_list, client_noconsensus.node.node_id) == 0
    assert_slashing(node.ppos.getCandidateInfo(client_noconsensus.node.node_id), economic.create_staking_limit)


def test_ZB_NP_36(new_genesis_env, client_consensus, client_noconsensus):
    """
    共识轮阈值大值变小值提案未生效时，节点零出块率处罚（存在零出块记录）
    """
    update_zero_produce(new_genesis_env)
    economic = client_noconsensus.economic
    node = economic.env.get_consensus_node_by_index(0)
    create_pledge_node_information(client_noconsensus)

    # Start execution use case
    initial_validator, slashing_node_list = gen_validator_list(economic.env.consensus_node_id_list(), client_noconsensus.node.node_id)
    log.info("verifier nodeid list: {}".format(initial_validator))
    assert_set_validator_list(node, slashing_node_list)
    client_noconsensus.node.stop()

    economic.wait_consensus_blocknum(node)
    log.info("Current block height: {}".format(node.eth.blockNumber))
    current_validator = get_pledge_list(node.ppos.getValidatorList)
    log.info("current validator: {}".format(current_validator))
    assert_set_validator_list(node, slashing_node_list)
    wait_slashing_list = node.debug.getWaitSlashingNodeList()
    log.info("Zero block record table：{}".format(wait_slashing_list))
    assert get_slash_count(wait_slashing_list, client_noconsensus.node.node_id) == 0

    economic.wait_consensus_blocknum(node)
    log.info("Current block height: {}".format(node.eth.blockNumber))
    current_validator = get_pledge_list(node.ppos.getValidatorList)
    log.info("current validator: {}".format(current_validator))
    assert_set_validator_list(node, initial_validator)
    wait_slashing_list = node.debug.getWaitSlashingNodeList()
    log.info("Zero block record table：{}".format(wait_slashing_list))
    assert get_slash_count(wait_slashing_list, client_noconsensus.node.node_id) == 0

    economic.wait_consensus_blocknum(node)
    log.info("Current block height: {}".format(node.eth.blockNumber))
    current_validator = get_pledge_list(node.ppos.getValidatorList)
    log.info("current validator: {}".format(current_validator))
    assert_set_validator_list(node, slashing_node_list)
    wait_slashing_list = node.debug.getWaitSlashingNodeList()
    log.info("Zero block record table：{}".format(wait_slashing_list))
    assert get_slash_count(wait_slashing_list, client_noconsensus.node.node_id) == to_int(1)
    param_governance_verify_before_endblock(client_consensus, 'slashing', 'zeroProduceCumulativeTime', '3', True)

    economic.wait_consensus_blocknum(node)
    log.info("Current block height: {}".format(node.eth.blockNumber))
    current_validator = get_pledge_list(node.ppos.getValidatorList)
    log.info("current validator: {}".format(current_validator))
    wait_slashing_list = node.debug.getWaitSlashingNodeList()
    log.info("Zero block record table：{}".format(wait_slashing_list))
    assert get_slash_count(wait_slashing_list, client_noconsensus.node.node_id) == to_int(11)

    economic.wait_consensus_blocknum(node)
    log.info("Current block height: {}".format(node.eth.blockNumber))
    current_validator = get_pledge_list(node.ppos.getValidatorList)
    log.info("current validator: {}".format(current_validator))
    wait_slashing_list = node.debug.getWaitSlashingNodeList()
    log.info("Zero block record table：{}".format(wait_slashing_list))
    assert get_slash_count(wait_slashing_list, client_noconsensus.node.node_id) == to_int(11)

    economic.wait_consensus_blocknum(node)
    log.info("Current block height: {}".format(node.eth.blockNumber))
    current_validator = get_pledge_list(node.ppos.getValidatorList)
    log.info("current validator: {}".format(current_validator))
    wait_slashing_list = node.debug.getWaitSlashingNodeList()
    log.info("Zero block record table：{}".format(wait_slashing_list))
    assert get_slash_count(wait_slashing_list, client_noconsensus.node.node_id) == 0
    assert_slashing(node.ppos.getCandidateInfo(client_noconsensus.node.node_id), economic.create_staking_limit)


def test_ZB_NP_37(new_genesis_env, client_consensus, client_noconsensus):
    """
    零出块率次数大值变小值提案未生效时，节点零出块率处罚（不存在零出块记录）
    """
    update_zero_produce(new_genesis_env)
    economic = client_noconsensus.economic
    node = economic.env.get_consensus_node_by_index(0)
    create_pledge_node_information(client_noconsensus)

    # Voting failed and expired
    param_governance_verify_before_endblock(client_consensus, 'slashing', 'zeroProduceNumberThreshold', '2', True)

    # Start execution use case
    initial_validator, slashing_node_list = gen_validator_list(economic.env.consensus_node_id_list(), client_noconsensus.node.node_id)
    log.info("verifier nodeid list: {}".format(initial_validator))
    client_noconsensus.node.stop()

    set_slashing(initial_validator, slashing_node_list, node, economic, "1011")

    economic.wait_consensus_blocknum(node, 1)
    log.info("Current block height: {}".format(node.eth.blockNumber))
    current_validator = get_pledge_list(node.ppos.getValidatorList)
    log.info("current validator: {}".format(current_validator))
    wait_slashing_list = node.debug.getWaitSlashingNodeList()
    log.info("Zero block record table：{}".format(wait_slashing_list))
    assert get_slash_count(wait_slashing_list, client_noconsensus.node.node_id) == 0
    assert_slashing(node.ppos.getCandidateInfo(client_noconsensus.node.node_id), economic.create_staking_limit)


def test_ZB_NP_38(new_genesis_env, client_consensus, client_noconsensus):
    """
    零出块率次数大值变小值提案未生效时，节点零出块率处罚（存在零出块记录）
    """
    update_zero_produce(new_genesis_env)
    economic = client_noconsensus.economic
    node = economic.env.get_consensus_node_by_index(0)
    create_pledge_node_information(client_noconsensus)

    initial_validator, slashing_node_list = gen_validator_list(economic.env.consensus_node_id_list(), client_noconsensus.node.node_id)
    log.info("verifier nodeid list: {}".format(initial_validator))
    client_noconsensus.node.stop()

    set_slashing(initial_validator, slashing_node_list, node, economic, "11")
    # Voting failed and expired
    assert_set_validator_list(node, initial_validator)
    param_governance_verify_before_endblock(client_consensus, 'slashing', 'zeroProduceNumberThreshold', '2', True)

    economic.wait_consensus_blocknum(node)
    log.info("Current block height: {}".format(node.eth.blockNumber))
    current_validator = get_pledge_list(node.ppos.getValidatorList)
    log.info("current validator: {}".format(current_validator))
    assert_set_validator_list(node, slashing_node_list)

    economic.wait_consensus_blocknum(node)
    log.info("Current block height: {}".format(node.eth.blockNumber))
    current_validator = get_pledge_list(node.ppos.getValidatorList)
    log.info("current validator: {}".format(current_validator))
    wait_slashing_list = node.debug.getWaitSlashingNodeList()
    log.info("Zero block record table：{}".format(wait_slashing_list))
    assert get_slash_count(wait_slashing_list, client_noconsensus.node.node_id) == to_int(11)

    economic.wait_consensus_blocknum(node)
    log.info("Current block height: {}".format(node.eth.blockNumber))
    current_validator = get_pledge_list(node.ppos.getValidatorList)
    log.info("current validator: {}".format(current_validator))
    wait_slashing_list = node.debug.getWaitSlashingNodeList()
    log.info("Zero block record table：{}".format(wait_slashing_list))
    assert get_slash_count(wait_slashing_list, client_noconsensus.node.node_id) == to_int(11)

    economic.wait_consensus_blocknum(node)
    log.info("Current block height: {}".format(node.eth.blockNumber))
    current_validator = get_pledge_list(node.ppos.getValidatorList)
    log.info("current validator: {}".format(current_validator))
    wait_slashing_list = node.debug.getWaitSlashingNodeList()
    log.info("Zero block record table：{}".format(wait_slashing_list))
    assert get_slash_count(wait_slashing_list, client_noconsensus.node.node_id) == 0
    assert_slashing(node.ppos.getCandidateInfo(client_noconsensus.node.node_id), economic.create_staking_limit)


def test_ZB_NP_39(new_genesis_env, client_consensus, client_noconsensus):
    """
    零出块率次数小值变大值提案未生效时，节点零出块率处罚（不存在零出块记录）
    """
    update_zero_produce(new_genesis_env)
    economic = client_noconsensus.economic
    node = economic.env.get_consensus_node_by_index(0)
    create_pledge_node_information(client_noconsensus)

    # Voting failed and expired
    param_governance_verify_before_endblock(client_consensus, 'slashing', 'zeroProduceNumberThreshold', '4', True)

    # Start execution use case
    initial_validator, slashing_node_list = gen_validator_list(economic.env.consensus_node_id_list(), client_noconsensus.node.node_id)
    log.info("verifier nodeid list: {}".format(initial_validator))
    client_noconsensus.node.stop()

    set_slashing(initial_validator, slashing_node_list, node, economic, "1110")

    economic.wait_consensus_blocknum(node, 1)
    log.info("Current block height: {}".format(node.eth.blockNumber))
    current_validator = get_pledge_list(node.ppos.getValidatorList)
    log.info("current validator: {}".format(current_validator))
    wait_slashing_list = node.debug.getWaitSlashingNodeList()
    log.info("Zero block record table：{}".format(wait_slashing_list))
    assert get_slash_count(wait_slashing_list, client_noconsensus.node.node_id) == 0
    assert_slashing(node.ppos.getCandidateInfo(client_noconsensus.node.node_id), economic.create_staking_limit)


def test_ZB_NP_40(new_genesis_env, client_consensus, client_noconsensus):
    """
    零出块率次数小值变大值提案未生效时，节点零出块率处罚（存在零出块记录）
    """
    update_zero_produce(new_genesis_env)
    economic = client_noconsensus.economic
    node = economic.env.get_consensus_node_by_index(0)
    create_pledge_node_information(client_noconsensus)

    initial_validator, slashing_node_list = gen_validator_list(economic.env.consensus_node_id_list(), client_noconsensus.node.node_id)
    log.info("verifier nodeid list: {}".format(initial_validator))
    client_noconsensus.node.stop()

    set_slashing(initial_validator, slashing_node_list, node, economic, "11")
    # Voting failed and expired
    assert_set_validator_list(node, initial_validator)
    param_governance_verify_before_endblock(client_consensus, 'slashing', 'zeroProduceNumberThreshold', '4', True)

    economic.wait_consensus_blocknum(node)
    log.info("Current block height: {}".format(node.eth.blockNumber))
    current_validator = get_pledge_list(node.ppos.getValidatorList)
    log.info("current validator: {}".format(current_validator))
    assert_set_validator_list(node, slashing_node_list)

    economic.wait_consensus_blocknum(node)
    log.info("Current block height: {}".format(node.eth.blockNumber))
    current_validator = get_pledge_list(node.ppos.getValidatorList)
    log.info("current validator: {}".format(current_validator))
    wait_slashing_list = node.debug.getWaitSlashingNodeList()
    log.info("Zero block record table：{}".format(wait_slashing_list))
    assert get_slash_count(wait_slashing_list, client_noconsensus.node.node_id) == to_int(11)

    economic.wait_consensus_blocknum(node)
    log.info("Current block height: {}".format(node.eth.blockNumber))
    current_validator = get_pledge_list(node.ppos.getValidatorList)
    log.info("current validator: {}".format(current_validator))
    wait_slashing_list = node.debug.getWaitSlashingNodeList()
    log.info("Zero block record table：{}".format(wait_slashing_list))
    assert get_slash_count(wait_slashing_list, client_noconsensus.node.node_id) == to_int(11)

    economic.wait_consensus_blocknum(node)
    log.info("Current block height: {}".format(node.eth.blockNumber))
    current_validator = get_pledge_list(node.ppos.getValidatorList)
    log.info("current validator: {}".format(current_validator))
    wait_slashing_list = node.debug.getWaitSlashingNodeList()
    log.info("Zero block record table：{}".format(wait_slashing_list))
    assert get_slash_count(wait_slashing_list, client_noconsensus.node.node_id) == 0
    assert_slashing(node.ppos.getCandidateInfo(client_noconsensus.node.node_id), economic.create_staking_limit)


def test_ZB_NP_41_42(new_genesis_env, client_consensus, client_noconsensus):
    """
    ZB_NP_41:共识轮阈值小值变大值提案生效，未到新共识轮位数（不存在零出块记录）
    ZB_NP_42:共识轮阈值小值变大值提案生效，等待共识轮且不当选共识验证人（不存在零出块记录）
    """
    update_zero_produce(new_genesis_env, 3)
    economic = client_noconsensus.economic
    node = economic.env.get_consensus_node_by_index(0)

    # Voting failed and expired
    param_governance_verify(client_consensus, 'slashing', 'zeroProduceCumulativeTime', '4', True)

    # Start execution use case
    create_pledge_node_information(client_noconsensus)
    initial_validator, slashing_node_list = gen_validator_list(economic.env.consensus_node_id_list(), client_noconsensus.node.node_id)
    log.info("verifier nodeid list: {}".format(initial_validator))
    client_noconsensus.node.stop()

    set_slashing(initial_validator, slashing_node_list, node, economic, "111")
    assert_set_validator_list(node, initial_validator)

    economic.wait_consensus_blocknum(node, 1)
    log.info("Current block height: {}".format(node.eth.blockNumber))
    current_validator = get_pledge_list(node.ppos.getValidatorList)
    log.info("current validator: {}".format(current_validator))
    wait_slashing_list = node.debug.getWaitSlashingNodeList()
    log.info("Zero block record table：{}".format(wait_slashing_list))
    assert get_slash_count(wait_slashing_list, client_noconsensus.node.node_id) == to_int(111)
    assert_not_slashing(node.ppos.getCandidateInfo(client_noconsensus.node.node_id), economic.create_staking_limit)

    economic.wait_consensus_blocknum(node)
    log.info("Current block height: {}".format(node.eth.blockNumber))
    current_validator = get_pledge_list(node.ppos.getValidatorList)
    log.info("current validator: {}".format(current_validator))
    wait_slashing_list = node.debug.getWaitSlashingNodeList()
    log.info("Zero block record table：{}".format(wait_slashing_list))
    assert get_slash_count(wait_slashing_list, client_noconsensus.node.node_id) == 0
    assert_slashing(node.ppos.getCandidateInfo(client_noconsensus.node.node_id), economic.create_staking_limit)


def test_ZB_NP_43_44(new_genesis_env, client_consensus, client_noconsensus):
    """
    ZB_NP_43:共识轮阈值小值变大值提案生效，未到新共识轮位数（存在零出块记录）
    ZB_NP_44:共识轮阈值小值变大值提案生效，等待共识轮且不当选共识验证人（存在零出块记录）
    """
    update_zero_produce(new_genesis_env, 3)
    economic = client_noconsensus.economic
    node = economic.env.get_consensus_node_by_index(0)

    # Voting failed and expired
    param_governance_verify_before_endblock(client_consensus, 'slashing', 'zeroProduceCumulativeTime', '4', True)

    # Start execution use case
    create_pledge_node_information(client_noconsensus)
    initial_validator, slashing_node_list = gen_validator_list(economic.env.consensus_node_id_list(), client_noconsensus.node.node_id)
    log.info("verifier nodeid list: {}".format(initial_validator))
    client_noconsensus.node.stop()

    set_slashing(initial_validator, slashing_node_list, node, economic, "0111")
    assert_set_validator_list(node, initial_validator)

    economic.wait_consensus_blocknum(node, 1)
    log.info("Current block height: {}".format(node.eth.blockNumber))
    current_validator = get_pledge_list(node.ppos.getValidatorList)
    log.info("current validator: {}".format(current_validator))
    wait_slashing_list = node.debug.getWaitSlashingNodeList()
    log.info("Zero block record table：{}".format(wait_slashing_list))
    assert get_slash_count(wait_slashing_list, client_noconsensus.node.node_id) == to_int(111)
    assert_not_slashing(node.ppos.getCandidateInfo(client_noconsensus.node.node_id), economic.create_staking_limit)

    economic.wait_consensus_blocknum(node)
    log.info("Current block height: {}".format(node.eth.blockNumber))
    current_validator = get_pledge_list(node.ppos.getValidatorList)
    log.info("current validator: {}".format(current_validator))
    wait_slashing_list = node.debug.getWaitSlashingNodeList()
    log.info("Zero block record table：{}".format(wait_slashing_list))
    assert get_slash_count(wait_slashing_list, client_noconsensus.node.node_id) == 0
    assert_slashing(node.ppos.getCandidateInfo(client_noconsensus.node.node_id), economic.create_staking_limit)


def test_ZB_NP_45(new_genesis_env, client_consensus, client_noconsensus):
    """
    共识轮阈值大值变小值提案生效，节点达到零出块率次数，（不存在零出块记录）
    """
    update_zero_produce(new_genesis_env)
    economic = client_noconsensus.economic
    node = economic.env.get_consensus_node_by_index(0)

    # Voting failed and expired
    param_governance_verify(client_consensus, 'slashing', 'zeroProduceCumulativeTime', '3', True)

    # Start execution use case
    create_pledge_node_information(client_noconsensus)
    initial_validator, slashing_node_list = gen_validator_list(economic.env.consensus_node_id_list(), client_noconsensus.node.node_id)
    log.info("verifier nodeid list: {}".format(initial_validator))
    client_noconsensus.node.stop()

    set_slashing(initial_validator, slashing_node_list, node, economic, "111")
    assert_set_validator_list(node, initial_validator)

    economic.wait_consensus_blocknum(node, 1)
    log.info("Current block height: {}".format(node.eth.blockNumber))
    current_validator = get_pledge_list(node.ppos.getValidatorList)
    log.info("current validator: {}".format(current_validator))
    wait_slashing_list = node.debug.getWaitSlashingNodeList()
    log.info("Zero block record table：{}".format(wait_slashing_list))
    assert get_slash_count(wait_slashing_list, client_noconsensus.node.node_id) == 0
    assert_slashing(node.ppos.getCandidateInfo(client_noconsensus.node.node_id), economic.create_staking_limit)


def test_ZB_NP_47(new_genesis_env, client_consensus, client_noconsensus):
    """
    共识轮阈值大值变小值提案生效，节点达到零出块率次数（存在零出块记录）（一）
    """
    update_zero_produce(new_genesis_env, 4, 3)
    economic = client_noconsensus.economic
    node = economic.env.get_consensus_node_by_index(0)

    # Voting failed and expired
    param_governance_verify_before_endblock(client_consensus, 'slashing', 'zeroProduceCumulativeTime', '3', True)

    # Start execution use case
    create_pledge_node_information(client_noconsensus)
    initial_validator, slashing_node_list = gen_validator_list(economic.env.consensus_node_id_list(), client_noconsensus.node.node_id)
    log.info("verifier nodeid list: {}".format(initial_validator))
    client_noconsensus.node.stop()

    set_slashing(initial_validator, slashing_node_list, node, economic, "011")
    assert_set_validator_list(node, slashing_node_list)

    economic.wait_consensus_blocknum(node, 1)
    log.info("Current block height: {}".format(node.eth.blockNumber))
    current_validator = get_pledge_list(node.ppos.getValidatorList)
    log.info("current validator: {}".format(current_validator))
    wait_slashing_list = node.debug.getWaitSlashingNodeList()
    log.info("Zero block record table：{}".format(wait_slashing_list))
    assert get_slash_count(wait_slashing_list, client_noconsensus.node.node_id) == to_int(11)
    assert_not_slashing(node.ppos.getCandidateInfo(client_noconsensus.node.node_id), economic.create_staking_limit)

    economic.wait_consensus_blocknum(node)
    log.info("Current block height: {}".format(node.eth.blockNumber))
    current_validator = get_pledge_list(node.ppos.getValidatorList)
    log.info("current validator: {}".format(current_validator))
    wait_slashing_list = node.debug.getWaitSlashingNodeList()
    log.info("Zero block record table：{}".format(wait_slashing_list))
    assert get_slash_count(wait_slashing_list, client_noconsensus.node.node_id) == 0
    assert_slashing(node.ppos.getCandidateInfo(client_noconsensus.node.node_id), economic.create_staking_limit)


def test_ZB_NP_48(new_genesis_env, client_consensus, client_noconsensus):
    """
    共识轮阈值大值变小值提案生效，节点达到零出块率次数（存在零出块记录）（二）
    """
    update_zero_produce(new_genesis_env, 4, 2)
    economic = client_noconsensus.economic
    node = economic.env.get_consensus_node_by_index(0)

    # Voting failed and expired
    param_governance_verify_before_endblock(client_consensus, 'slashing', 'zeroProduceCumulativeTime', '3', True)

    # Start execution use case
    create_pledge_node_information(client_noconsensus)
    initial_validator, slashing_node_list = gen_validator_list(economic.env.consensus_node_id_list(), client_noconsensus.node.node_id)
    log.info("verifier nodeid list: {}".format(initial_validator))
    client_noconsensus.node.stop()

    set_slashing(initial_validator, slashing_node_list, node, economic, "011")
    assert_set_validator_list(node, initial_validator)

    economic.wait_consensus_blocknum(node, 1)
    log.info("Current block height: {}".format(node.eth.blockNumber))
    current_validator = get_pledge_list(node.ppos.getValidatorList)
    log.info("current validator: {}".format(current_validator))
    wait_slashing_list = node.debug.getWaitSlashingNodeList()
    log.info("Zero block record table：{}".format(wait_slashing_list))
    assert get_slash_count(wait_slashing_list, client_noconsensus.node.node_id) == to_int(11)
    assert_not_slashing(node.ppos.getCandidateInfo(client_noconsensus.node.node_id), economic.create_staking_limit)

    economic.wait_consensus_blocknum(node)
    log.info("Current block height: {}".format(node.eth.blockNumber))
    current_validator = get_pledge_list(node.ppos.getValidatorList)
    log.info("current validator: {}".format(current_validator))
    wait_slashing_list = node.debug.getWaitSlashingNodeList()
    log.info("Zero block record table：{}".format(wait_slashing_list))
    assert get_slash_count(wait_slashing_list, client_noconsensus.node.node_id) == 0
    assert_slashing(node.ppos.getCandidateInfo(client_noconsensus.node.node_id), economic.create_staking_limit)


def test_ZB_NP_49(new_genesis_env, client_consensus, client_noconsensus):
    """
    共识轮阈值大值变小值提案生效，节点未达到零出块率次数（存在零出块记录）
    """
    update_zero_produce(new_genesis_env)
    economic = client_noconsensus.economic
    node = economic.env.get_consensus_node_by_index(0)

    # Voting failed and expired
    param_governance_verify_before_endblock(client_consensus, 'slashing', 'zeroProduceCumulativeTime', '3', True)

    # Start execution use case
    create_pledge_node_information(client_noconsensus)
    initial_validator, slashing_node_list = gen_validator_list(economic.env.consensus_node_id_list(), client_noconsensus.node.node_id)
    log.info("verifier nodeid list: {}".format(initial_validator))
    client_noconsensus.node.stop()

    set_slashing(initial_validator, slashing_node_list, node, economic, "010")
    assert_set_validator_list(node, slashing_node_list)

    economic.wait_consensus_blocknum(node, 1)
    log.info("Current block height: {}".format(node.eth.blockNumber))
    current_validator = get_pledge_list(node.ppos.getValidatorList)
    log.info("current validator: {}".format(current_validator))
    wait_slashing_list = node.debug.getWaitSlashingNodeList()
    log.info("Zero block record table：{}".format(wait_slashing_list))
    assert get_slash_count(wait_slashing_list, client_noconsensus.node.node_id) == to_int(1)
    assert_not_slashing(node.ppos.getCandidateInfo(client_noconsensus.node.node_id), economic.create_staking_limit)

    economic.wait_consensus_blocknum(node)
    log.info("Current block height: {}".format(node.eth.blockNumber))
    current_validator = get_pledge_list(node.ppos.getValidatorList)
    log.info("current validator: {}".format(current_validator))
    wait_slashing_list = node.debug.getWaitSlashingNodeList()
    log.info("Zero block record table：{}".format(wait_slashing_list))
    assert get_slash_count(wait_slashing_list, client_noconsensus.node.node_id) == to_int(101)
    assert_not_slashing(node.ppos.getCandidateInfo(client_noconsensus.node.node_id), economic.create_staking_limit)


def test_ZB_NP_50(new_genesis_env, client_consensus, client_noconsensus):
    """
    零出块率次数大值变小值提案生效，未到达共识轮位数（不存在零出块记录）
    """
    update_zero_produce(new_genesis_env)
    economic = client_noconsensus.economic
    node = economic.env.get_consensus_node_by_index(0)

    # Voting failed and expired
    param_governance_verify(client_consensus, 'slashing', 'zeroProduceNumberThreshold', '2', True)

    # Start execution use case
    create_pledge_node_information(client_noconsensus)
    initial_validator, slashing_node_list = gen_validator_list(economic.env.consensus_node_id_list(), client_noconsensus.node.node_id)
    log.info("verifier nodeid list: {}".format(initial_validator))
    client_noconsensus.node.stop()

    set_slashing(initial_validator, slashing_node_list, node, economic, "101")

    economic.wait_consensus_blocknum(node, 1)
    log.info("Current block height: {}".format(node.eth.blockNumber))
    current_validator = get_pledge_list(node.ppos.getValidatorList)
    log.info("current validator: {}".format(current_validator))
    wait_slashing_list = node.debug.getWaitSlashingNodeList()
    log.info("Zero block record table：{}".format(wait_slashing_list))
    assert get_slash_count(wait_slashing_list, client_noconsensus.node.node_id) == to_int(101)
    assert_not_slashing(node.ppos.getCandidateInfo(client_noconsensus.node.node_id), economic.create_staking_limit)


def test_ZB_NP_51(new_genesis_env, client_consensus, client_noconsensus):
    """
    零出块率次数大值变小值提案生效，到达共识轮位数，未达到零出块率次数不处罚（不存在零出块记录）
    """
    update_zero_produce(new_genesis_env)
    economic = client_noconsensus.economic
    node = economic.env.get_consensus_node_by_index(0)

    # Voting failed and expired
    param_governance_verify(client_consensus, 'slashing', 'zeroProduceNumberThreshold', '2', True)

    # Start execution use case
    create_pledge_node_information(client_noconsensus)
    initial_validator, slashing_node_list = gen_validator_list(economic.env.consensus_node_id_list(), client_noconsensus.node.node_id)
    log.info("verifier nodeid list: {}".format(initial_validator))
    client_noconsensus.node.stop()

    set_slashing(initial_validator, slashing_node_list, node, economic, "1000")

    economic.wait_consensus_blocknum(node, 1)
    log.info("Current block height: {}".format(node.eth.blockNumber))
    current_validator = get_pledge_list(node.ppos.getValidatorList)
    log.info("current validator: {}".format(current_validator))
    wait_slashing_list = node.debug.getWaitSlashingNodeList()
    log.info("Zero block record table：{}".format(wait_slashing_list))
    assert get_slash_count(wait_slashing_list, client_noconsensus.node.node_id) == to_int(1)
    assert_not_slashing(node.ppos.getCandidateInfo(client_noconsensus.node.node_id), economic.create_staking_limit)

    economic.wait_consensus_blocknum(node)
    log.info("Current block height: {}".format(node.eth.blockNumber))
    current_validator = get_pledge_list(node.ppos.getValidatorList)
    log.info("current validator: {}".format(current_validator))
    wait_slashing_list = node.debug.getWaitSlashingNodeList()
    log.info("Zero block record table：{}".format(wait_slashing_list))
    assert get_slash_count(wait_slashing_list, client_noconsensus.node.node_id) == 0
    assert_not_slashing(node.ppos.getCandidateInfo(client_noconsensus.node.node_id), economic.create_staking_limit)


def test_ZB_NP_52(new_genesis_env, client_consensus, client_noconsensus):
    """
    零出块率次数大值变小值提案生效，到达共识轮位数，达到零出块率次数处罚（不存在零出块记录）
    """
    update_zero_produce(new_genesis_env)
    economic = client_noconsensus.economic
    node = economic.env.get_consensus_node_by_index(0)

    # Voting failed and expired
    param_governance_verify(client_consensus, 'slashing', 'zeroProduceNumberThreshold', '2', True)

    # Start execution use case
    create_pledge_node_information(client_noconsensus)
    initial_validator, slashing_node_list = gen_validator_list(economic.env.consensus_node_id_list(), client_noconsensus.node.node_id)
    log.info("verifier nodeid list: {}".format(initial_validator))
    client_noconsensus.node.stop()

    set_slashing(initial_validator, slashing_node_list, node, economic, "1010")

    economic.wait_consensus_blocknum(node, 1)
    log.info("Current block height: {}".format(node.eth.blockNumber))
    current_validator = get_pledge_list(node.ppos.getValidatorList)
    log.info("current validator: {}".format(current_validator))
    wait_slashing_list = node.debug.getWaitSlashingNodeList()
    log.info("Zero block record table：{}".format(wait_slashing_list))
    assert get_slash_count(wait_slashing_list, client_noconsensus.node.node_id) == 0
    assert_slashing(node.ppos.getCandidateInfo(client_noconsensus.node.node_id), economic.create_staking_limit)


def test_ZB_NP_53(new_genesis_env, client_consensus, client_noconsensus):
    """
    零出块率次数大值变小值提案生效，未到达共识轮位数（存在零出块记录）
    """
    update_zero_produce(new_genesis_env)
    economic = client_noconsensus.economic
    node = economic.env.get_consensus_node_by_index(0)

    # Voting failed and expired
    param_governance_verify_before_endblock(client_consensus, 'slashing', 'zeroProduceNumberThreshold', '2', True)

    # Start execution use case
    create_pledge_node_information(client_noconsensus)
    initial_validator, slashing_node_list = gen_validator_list(economic.env.consensus_node_id_list(), client_noconsensus.node.node_id)
    log.info("verifier nodeid list: {}".format(initial_validator))
    client_noconsensus.node.stop()

    set_slashing(initial_validator, slashing_node_list, node, economic, "010")
    assert_set_validator_list(node, slashing_node_list)

    economic.wait_consensus_blocknum(node, 1)
    log.info("Current block height: {}".format(node.eth.blockNumber))
    current_validator = get_pledge_list(node.ppos.getValidatorList)
    log.info("current validator: {}".format(current_validator))
    wait_slashing_list = node.debug.getWaitSlashingNodeList()
    log.info("Zero block record table：{}".format(wait_slashing_list))
    assert get_slash_count(wait_slashing_list, client_noconsensus.node.node_id) == to_int(1)
    assert_not_slashing(node.ppos.getCandidateInfo(client_noconsensus.node.node_id), economic.create_staking_limit)

    economic.wait_consensus_blocknum(node)
    log.info("Current block height: {}".format(node.eth.blockNumber))
    current_validator = get_pledge_list(node.ppos.getValidatorList)
    log.info("current validator: {}".format(current_validator))
    wait_slashing_list = node.debug.getWaitSlashingNodeList()
    log.info("Zero block record table：{}".format(wait_slashing_list))
    assert get_slash_count(wait_slashing_list, client_noconsensus.node.node_id) == to_int(101)
    assert_not_slashing(node.ppos.getCandidateInfo(client_noconsensus.node.node_id), economic.create_staking_limit)


def test_ZB_NP_54(new_genesis_env, client_consensus, client_noconsensus):
    """
    零出块率次数大值变小值提案生效，到达共识轮位数，未达到零出块率次数不处罚（存在零出块记录）
    """

    update_zero_produce(new_genesis_env)
    economic = client_noconsensus.economic
    node = economic.env.get_consensus_node_by_index(0)

    # Voting failed and expired
    param_governance_verify_before_endblock(client_consensus, 'slashing', 'zeroProduceNumberThreshold', '2', True)

    # Start execution use case
    create_pledge_node_information(client_noconsensus)
    initial_validator, slashing_node_list = gen_validator_list(economic.env.consensus_node_id_list(), client_noconsensus.node.node_id)
    log.info("verifier nodeid list: {}".format(initial_validator))
    client_noconsensus.node.stop()

    set_slashing(initial_validator, slashing_node_list, node, economic, "100")
    assert_set_validator_list(node, initial_validator)

    economic.wait_consensus_blocknum(node, 1)
    log.info("Current block height: {}".format(node.eth.blockNumber))
    current_validator = get_pledge_list(node.ppos.getValidatorList)
    log.info("current validator: {}".format(current_validator))
    wait_slashing_list = node.debug.getWaitSlashingNodeList()
    log.info("Zero block record table：{}".format(wait_slashing_list))
    assert get_slash_count(wait_slashing_list, client_noconsensus.node.node_id) == to_int(1)
    assert_not_slashing(node.ppos.getCandidateInfo(client_noconsensus.node.node_id), economic.create_staking_limit)

    economic.wait_consensus_blocknum(node)
    log.info("Current block height: {}".format(node.eth.blockNumber))
    current_validator = get_pledge_list(node.ppos.getValidatorList)
    log.info("current validator: {}".format(current_validator))
    wait_slashing_list = node.debug.getWaitSlashingNodeList()
    log.info("Zero block record table：{}".format(wait_slashing_list))
    assert get_slash_count(wait_slashing_list, client_noconsensus.node.node_id) == to_int(1)
    assert_not_slashing(node.ppos.getCandidateInfo(client_noconsensus.node.node_id), economic.create_staking_limit)


def test_ZB_NP_55(new_genesis_env, client_consensus, client_noconsensus):
    """
    零出块率次数大值变小值提案生效，到达共识轮位数，达到零出块率次数（存在零出块记录）
    """
    update_zero_produce(new_genesis_env)
    economic = client_noconsensus.economic
    node = economic.env.get_consensus_node_by_index(0)

    # Voting failed and expired
    param_governance_verify_before_endblock(client_consensus, 'slashing', 'zeroProduceNumberThreshold', '2', True)

    # Start execution use case
    create_pledge_node_information(client_noconsensus)
    initial_validator, slashing_node_list = gen_validator_list(economic.env.consensus_node_id_list(), client_noconsensus.node.node_id)
    log.info("verifier nodeid list: {}".format(initial_validator))
    client_noconsensus.node.stop()

    set_slashing(initial_validator, slashing_node_list, node, economic, "100")
    assert_set_validator_list(node, slashing_node_list)

    economic.wait_consensus_blocknum(node, 1)
    log.info("Current block height: {}".format(node.eth.blockNumber))
    current_validator = get_pledge_list(node.ppos.getValidatorList)
    log.info("current validator: {}".format(current_validator))
    wait_slashing_list = node.debug.getWaitSlashingNodeList()
    log.info("Zero block record table：{}".format(wait_slashing_list))
    assert get_slash_count(wait_slashing_list, client_noconsensus.node.node_id) == to_int(1)
    assert_not_slashing(node.ppos.getCandidateInfo(client_noconsensus.node.node_id), economic.create_staking_limit)

    economic.wait_consensus_blocknum(node)
    log.info("Current block height: {}".format(node.eth.blockNumber))
    current_validator = get_pledge_list(node.ppos.getValidatorList)
    log.info("current validator: {}".format(current_validator))
    wait_slashing_list = node.debug.getWaitSlashingNodeList()
    log.info("Zero block record table：{}".format(wait_slashing_list))
    assert get_slash_count(wait_slashing_list, client_noconsensus.node.node_id) == 0
    assert_slashing(node.ppos.getCandidateInfo(client_noconsensus.node.node_id), economic.create_staking_limit)


def test_ZB_NP_56(new_genesis_env, client_consensus, client_noconsensus):
    """
    ZB_NP_56:零出块率次数小值变大值提案生效，未到达共识轮位数，达到零出块率次数（不存在零出块记录）
    ZB_NP_57:零出块率次数小值变大值提案生效，未到达共识轮位数，达到零出块率次数，等待下共识轮且不当选共识验证人（不存在零出块记录）
    """
    update_zero_produce(new_genesis_env, 4, 2)
    economic = client_noconsensus.economic
    node = economic.env.get_consensus_node_by_index(0)

    # Voting failed and expired
    param_governance_verify(client_consensus, 'slashing', 'zeroProduceNumberThreshold', '3', True)

    # Start execution use case
    create_pledge_node_information(client_noconsensus)
    initial_validator, slashing_node_list = gen_validator_list(economic.env.consensus_node_id_list(), client_noconsensus.node.node_id)
    log.info("verifier nodeid list: {}".format(initial_validator))
    client_noconsensus.node.stop()

    set_slashing(initial_validator, slashing_node_list, node, economic, "111")
    assert_set_validator_list(node, slashing_node_list)

    economic.wait_consensus_blocknum(node, 1)
    log.info("Current block height: {}".format(node.eth.blockNumber))
    current_validator = get_pledge_list(node.ppos.getValidatorList)
    log.info("current validator: {}".format(current_validator))
    wait_slashing_list = node.debug.getWaitSlashingNodeList()
    log.info("Zero block record table：{}".format(wait_slashing_list))
    assert get_slash_count(wait_slashing_list, client_noconsensus.node.node_id) == to_int(111)
    assert_not_slashing(node.ppos.getCandidateInfo(client_noconsensus.node.node_id), economic.create_staking_limit)

    economic.wait_consensus_blocknum(node)
    log.info("Current block height: {}".format(node.eth.blockNumber))
    current_validator = get_pledge_list(node.ppos.getValidatorList)
    log.info("current validator: {}".format(current_validator))
    wait_slashing_list = node.debug.getWaitSlashingNodeList()
    log.info("Zero block record table：{}".format(wait_slashing_list))
    assert get_slash_count(wait_slashing_list, client_noconsensus.node.node_id) == 0
    assert_slashing(node.ppos.getCandidateInfo(client_noconsensus.node.node_id), economic.create_staking_limit)


def test_ZB_NP_58(new_genesis_env, client_consensus, client_noconsensus):
    """
    零出块率次数小值变大值提案生效，未到达共识轮位数，未达到零出块率次数，等待下共识轮且不当选共识验证人（不存在零出块记录）
    """
    update_zero_produce(new_genesis_env, 4, 2)
    economic = client_noconsensus.economic
    node = economic.env.get_consensus_node_by_index(0)

    # Voting failed and expired
    param_governance_verify(client_consensus, 'slashing', 'zeroProduceNumberThreshold', '3', True)

    # Start execution use case
    create_pledge_node_information(client_noconsensus)
    initial_validator, slashing_node_list = gen_validator_list(economic.env.consensus_node_id_list(), client_noconsensus.node.node_id)
    log.info("verifier nodeid list: {}".format(initial_validator))
    client_noconsensus.node.stop()

    set_slashing(initial_validator, slashing_node_list, node, economic, "101")
    assert_set_validator_list(node, initial_validator)

    economic.wait_consensus_blocknum(node, 1)
    log.info("Current block height: {}".format(node.eth.blockNumber))
    current_validator = get_pledge_list(node.ppos.getValidatorList)
    log.info("current validator: {}".format(current_validator))
    wait_slashing_list = node.debug.getWaitSlashingNodeList()
    log.info("Zero block record table：{}".format(wait_slashing_list))
    assert get_slash_count(wait_slashing_list, client_noconsensus.node.node_id) == to_int(101)
    assert_not_slashing(node.ppos.getCandidateInfo(client_noconsensus.node.node_id), economic.create_staking_limit)

    economic.wait_consensus_blocknum(node)
    log.info("Current block height: {}".format(node.eth.blockNumber))
    current_validator = get_pledge_list(node.ppos.getValidatorList)
    log.info("current validator: {}".format(current_validator))
    wait_slashing_list = node.debug.getWaitSlashingNodeList()
    log.info("Zero block record table：{}".format(wait_slashing_list))
    assert get_slash_count(wait_slashing_list, client_noconsensus.node.node_id) == to_int(101)
    assert_not_slashing(node.ppos.getCandidateInfo(client_noconsensus.node.node_id), economic.create_staking_limit)


def test_ZB_NP_59(new_genesis_env, client_consensus, client_noconsensus):
    """
    零出块率次数小值变大值提案生效，未到达共识轮位数，未达到零出块率次数，等待共识轮且当选共识验证人出块不处罚（不存在零出块记录）
    """
    update_zero_produce(new_genesis_env, 4, 2)
    economic = client_noconsensus.economic
    node = economic.env.get_consensus_node_by_index(0)

    # Voting failed and expired
    param_governance_verify(client_consensus, 'slashing', 'zeroProduceNumberThreshold', '3', True)

    # Start execution use case
    create_pledge_node_information(client_noconsensus)
    initial_validator, slashing_node_list = gen_validator_list(economic.env.consensus_node_id_list(), client_noconsensus.node.node_id)
    log.info("verifier nodeid list: {}".format(initial_validator))
    client_noconsensus.node.stop()

    set_slashing(initial_validator, slashing_node_list, node, economic, "101")
    client_noconsensus.node.start(False)
    assert_set_validator_list(node, slashing_node_list)

    economic.wait_consensus_blocknum(node, 1)
    log.info("Current block height: {}".format(node.eth.blockNumber))
    current_validator = get_pledge_list(node.ppos.getValidatorList)
    log.info("current validator: {}".format(current_validator))
    wait_slashing_list = node.debug.getWaitSlashingNodeList()
    log.info("Zero block record table：{}".format(wait_slashing_list))
    assert get_slash_count(wait_slashing_list, client_noconsensus.node.node_id) == 0
    assert_not_slashing(node.ppos.getCandidateInfo(client_noconsensus.node.node_id), economic.create_staking_limit)


def test_ZB_NP_60(new_genesis_env, client_consensus, client_noconsensus):
    """
    零出块率次数小值变大值提案生效，未到达共识轮位数，未达到零出块率次数，等待下共识轮且当选共识验证人未出块（不存在零出块记录）
    """
    update_zero_produce(new_genesis_env, 4, 2)
    economic = client_noconsensus.economic
    node = economic.env.get_consensus_node_by_index(0)

    # Voting failed and expired
    param_governance_verify(client_consensus, 'slashing', 'zeroProduceNumberThreshold', '3', True)

    # Start execution use case
    create_pledge_node_information(client_noconsensus)
    initial_validator, slashing_node_list = gen_validator_list(economic.env.consensus_node_id_list(), client_noconsensus.node.node_id)
    log.info("verifier nodeid list: {}".format(initial_validator))
    client_noconsensus.node.stop()

    set_slashing(initial_validator, slashing_node_list, node, economic, "101")
    assert_set_validator_list(node, slashing_node_list)

    economic.wait_consensus_blocknum(node, 1)
    log.info("Current block height: {}".format(node.eth.blockNumber))
    current_validator = get_pledge_list(node.ppos.getValidatorList)
    log.info("current validator: {}".format(current_validator))
    wait_slashing_list = node.debug.getWaitSlashingNodeList()
    log.info("Zero block record table：{}".format(wait_slashing_list))
    assert get_slash_count(wait_slashing_list, client_noconsensus.node.node_id) == to_int(101)
    assert_not_slashing(node.ppos.getCandidateInfo(client_noconsensus.node.node_id), economic.create_staking_limit)

    economic.wait_consensus_blocknum(node)
    log.info("Current block height: {}".format(node.eth.blockNumber))
    current_validator = get_pledge_list(node.ppos.getValidatorList)
    log.info("current validator: {}".format(current_validator))
    wait_slashing_list = node.debug.getWaitSlashingNodeList()
    log.info("Zero block record table：{}".format(wait_slashing_list))
    assert get_slash_count(wait_slashing_list, client_noconsensus.node.node_id) == 0
    assert_slashing(node.ppos.getCandidateInfo(client_noconsensus.node.node_id), economic.create_staking_limit)


def test_ZB_NP_61(new_genesis_env, client_consensus, client_noconsensus):
    """
    零出块率次数小值变大值提案生效，达到共识轮位数，未达到零出块率次数（不存在零出块记录）
    """
    update_zero_produce(new_genesis_env)
    economic = client_noconsensus.economic
    node = economic.env.get_consensus_node_by_index(0)

    # Voting failed and expired
    param_governance_verify(client_consensus, 'slashing', 'zeroProduceNumberThreshold', '4', True)

    # Start execution use case
    create_pledge_node_information(client_noconsensus)
    initial_validator, slashing_node_list = gen_validator_list(economic.env.consensus_node_id_list(), client_noconsensus.node.node_id)
    log.info("verifier nodeid list: {}".format(initial_validator))
    client_noconsensus.node.stop()

    set_slashing(initial_validator, slashing_node_list, node, economic, "1101")

    economic.wait_consensus_blocknum(node, 1)
    log.info("Current block height: {}".format(node.eth.blockNumber))
    current_validator = get_pledge_list(node.ppos.getValidatorList)
    log.info("current validator: {}".format(current_validator))
    wait_slashing_list = node.debug.getWaitSlashingNodeList()
    log.info("Zero block record table：{}".format(wait_slashing_list))
    assert get_slash_count(wait_slashing_list, client_noconsensus.node.node_id) == to_int(1011)
    assert_not_slashing(node.ppos.getCandidateInfo(client_noconsensus.node.node_id), economic.create_staking_limit)


def test_ZB_NP_62_63(new_genesis_env, client_consensus, client_noconsensus):
    """
    ZB_NP_62:零出块率次数小值变大值提案生效，达到共识轮位数，未达到零出块率次数（存在零出块记录）
    ZB_NP_63:零出块率次数小值变大值提案生效，未到达共识轮位数，等待下共识轮且不当选共识验证人达到零出块率次数（存在零出块记录）
    """
    update_zero_produce(new_genesis_env, 4, 2)
    economic = client_noconsensus.economic
    node = economic.env.get_consensus_node_by_index(0)

    # Voting failed and expired
    param_governance_verify_before_endblock(client_consensus, 'slashing', 'zeroProduceNumberThreshold', '3', True)

    # Start execution use case
    create_pledge_node_information(client_noconsensus)
    initial_validator, slashing_node_list = gen_validator_list(economic.env.consensus_node_id_list(), client_noconsensus.node.node_id)
    log.info("verifier nodeid list: {}".format(initial_validator))
    client_noconsensus.node.stop()

    set_slashing(initial_validator, slashing_node_list, node, economic, "011")
    assert_set_validator_list(node, slashing_node_list)

    economic.wait_consensus_blocknum(node)
    log.info("Current block height: {}".format(node.eth.blockNumber))
    current_validator = get_pledge_list(node.ppos.getValidatorList)
    log.info("current validator: {}".format(current_validator))
    assert_set_validator_list(node, initial_validator)

    economic.wait_consensus_blocknum(node)
    log.info("Current block height: {}".format(node.eth.blockNumber))
    current_validator = get_pledge_list(node.ppos.getValidatorList)
    log.info("current validator: {}".format(current_validator))
    wait_slashing_list = node.debug.getWaitSlashingNodeList()
    log.info("Zero block record table：{}".format(wait_slashing_list))
    assert get_slash_count(wait_slashing_list, client_noconsensus.node.node_id) == to_int(11)
    assert_not_slashing(node.ppos.getCandidateInfo(client_noconsensus.node.node_id), economic.create_staking_limit)

    economic.wait_consensus_blocknum(node)
    log.info("Current block height: {}".format(node.eth.blockNumber))
    current_validator = get_pledge_list(node.ppos.getValidatorList)
    log.info("current validator: {}".format(current_validator))
    wait_slashing_list = node.debug.getWaitSlashingNodeList()
    log.info("Zero block record table：{}".format(wait_slashing_list))
    assert get_slash_count(wait_slashing_list, client_noconsensus.node.node_id) == to_int(111)
    assert_not_slashing(node.ppos.getCandidateInfo(client_noconsensus.node.node_id), economic.create_staking_limit)

    economic.wait_consensus_blocknum(node)
    log.info("Current block height: {}".format(node.eth.blockNumber))
    current_validator = get_pledge_list(node.ppos.getValidatorList)
    log.info("current validator: {}".format(current_validator))
    wait_slashing_list = node.debug.getWaitSlashingNodeList()
    log.info("Zero block record table：{}".format(wait_slashing_list))
    assert get_slash_count(wait_slashing_list, client_noconsensus.node.node_id) == 0
    assert_slashing(node.ppos.getCandidateInfo(client_noconsensus.node.node_id), economic.create_staking_limit)


def test_ZB_NP_64(new_genesis_env, client_consensus, client_noconsensus):
    """
    零出块率次数小值变大值提案生效，未到达共识轮位数，等待下共识轮且不当选共识验证人未达到零出块率次数（存在零出块记录）
    """
    update_zero_produce(new_genesis_env)
    economic = client_noconsensus.economic
    node = economic.env.get_consensus_node_by_index(0)

    # Voting failed and expired
    param_governance_verify_before_endblock(client_consensus, 'slashing', 'zeroProduceNumberThreshold', '4', True)

    # Start execution use case
    create_pledge_node_information(client_noconsensus)
    initial_validator, slashing_node_list = gen_validator_list(economic.env.consensus_node_id_list(), client_noconsensus.node.node_id)
    log.info("verifier nodeid list: {}".format(initial_validator))
    client_noconsensus.node.stop()

    set_slashing(initial_validator, slashing_node_list, node, economic, "011")
    assert_set_validator_list(node, slashing_node_list)

    economic.wait_consensus_blocknum(node)
    log.info("Current block height: {}".format(node.eth.blockNumber))
    current_validator = get_pledge_list(node.ppos.getValidatorList)
    log.info("current validator: {}".format(current_validator))
    assert_set_validator_list(node, initial_validator)

    economic.wait_consensus_blocknum(node)
    log.info("Current block height: {}".format(node.eth.blockNumber))
    current_validator = get_pledge_list(node.ppos.getValidatorList)
    log.info("current validator: {}".format(current_validator))
    wait_slashing_list = node.debug.getWaitSlashingNodeList()
    log.info("Zero block record table：{}".format(wait_slashing_list))
    assert get_slash_count(wait_slashing_list, client_noconsensus.node.node_id) == to_int(11)
    assert_not_slashing(node.ppos.getCandidateInfo(client_noconsensus.node.node_id), economic.create_staking_limit)

    economic.wait_consensus_blocknum(node)
    log.info("Current block height: {}".format(node.eth.blockNumber))
    current_validator = get_pledge_list(node.ppos.getValidatorList)
    log.info("current validator: {}".format(current_validator))
    wait_slashing_list = node.debug.getWaitSlashingNodeList()
    log.info("Zero block record table：{}".format(wait_slashing_list))
    assert get_slash_count(wait_slashing_list, client_noconsensus.node.node_id) == to_int(111)
    assert_not_slashing(node.ppos.getCandidateInfo(client_noconsensus.node.node_id), economic.create_staking_limit)

    economic.wait_consensus_blocknum(node)
    log.info("Current block height: {}".format(node.eth.blockNumber))
    current_validator = get_pledge_list(node.ppos.getValidatorList)
    log.info("current validator: {}".format(current_validator))
    wait_slashing_list = node.debug.getWaitSlashingNodeList()
    log.info("Zero block record table：{}".format(wait_slashing_list))
    assert get_slash_count(wait_slashing_list, client_noconsensus.node.node_id) == to_int(111)
    assert_not_slashing(node.ppos.getCandidateInfo(client_noconsensus.node.node_id), economic.create_staking_limit)


def test_ZB_NP_65(new_genesis_env, client_consensus, client_noconsensus):
    """
    零出块率次数小值变大值提案生效，未到达共识轮位数，等待下共识轮且当选共识验证人出块（存在零出块记录）
    """
    update_zero_produce(new_genesis_env)
    economic = client_noconsensus.economic
    node = economic.env.get_consensus_node_by_index(0)

    # Voting failed and expired
    param_governance_verify_before_endblock(client_consensus, 'slashing', 'zeroProduceNumberThreshold', '4', True)

    # Start execution use case
    create_pledge_node_information(client_noconsensus)
    initial_validator, slashing_node_list = gen_validator_list(economic.env.consensus_node_id_list(), client_noconsensus.node.node_id)
    log.info("verifier nodeid list: {}".format(initial_validator))
    client_noconsensus.node.stop()

    set_slashing(initial_validator, slashing_node_list, node, economic, "011")
    assert_set_validator_list(node, slashing_node_list)
    client_noconsensus.node.start(False)

    economic.wait_consensus_blocknum(node)
    log.info("Current block height: {}".format(node.eth.blockNumber))
    current_validator = get_pledge_list(node.ppos.getValidatorList)
    log.info("current validator: {}".format(current_validator))
    assert_set_validator_list(node, initial_validator)

    economic.wait_consensus_blocknum(node)
    log.info("Current block height: {}".format(node.eth.blockNumber))
    current_validator = get_pledge_list(node.ppos.getValidatorList)
    log.info("current validator: {}".format(current_validator))
    wait_slashing_list = node.debug.getWaitSlashingNodeList()
    log.info("Zero block record table：{}".format(wait_slashing_list))
    assert get_slash_count(wait_slashing_list, client_noconsensus.node.node_id) == 0
    assert_not_slashing(node.ppos.getCandidateInfo(client_noconsensus.node.node_id), economic.create_staking_limit)


def test_ZB_NP_66(new_genesis_env, client_consensus, client_noconsensus):
    """
    零出块率次数小值变大值提案生效，未到达共识轮位数，等待下共识轮且当选共识验证人未出块未达到零出块率次数（存在零出块记录）
    """
    update_zero_produce(new_genesis_env, 4, 2)
    economic = client_noconsensus.economic
    node = economic.env.get_consensus_node_by_index(0)

    # Voting failed and expired
    param_governance_verify_before_endblock(client_consensus, 'slashing', 'zeroProduceNumberThreshold', '3', True)

    # Start execution use case
    create_pledge_node_information(client_noconsensus)
    initial_validator, slashing_node_list = gen_validator_list(economic.env.consensus_node_id_list(), client_noconsensus.node.node_id)
    log.info("verifier nodeid list: {}".format(initial_validator))
    client_noconsensus.node.stop()

    set_slashing(initial_validator, slashing_node_list, node, economic, "0110")
    assert_set_validator_list(node, slashing_node_list)

    economic.wait_consensus_blocknum(node, 1)
    log.info("Current block height: {}".format(node.eth.blockNumber))
    current_validator = get_pledge_list(node.ppos.getValidatorList)
    log.info("current validator: {}".format(current_validator))
    wait_slashing_list = node.debug.getWaitSlashingNodeList()
    log.info("Zero block record table：{}".format(wait_slashing_list))
    assert get_slash_count(wait_slashing_list, client_noconsensus.node.node_id) == to_int(11)
    assert_not_slashing(node.ppos.getCandidateInfo(client_noconsensus.node.node_id), economic.create_staking_limit)

    economic.wait_consensus_blocknum(node)
    log.info("Current block height: {}".format(node.eth.blockNumber))
    current_validator = get_pledge_list(node.ppos.getValidatorList)
    log.info("current validator: {}".format(current_validator))
    wait_slashing_list = node.debug.getWaitSlashingNodeList()
    log.info("Zero block record table：{}".format(wait_slashing_list))
    assert get_slash_count(wait_slashing_list, client_noconsensus.node.node_id) == 0
    assert_slashing(node.ppos.getCandidateInfo(client_noconsensus.node.node_id), economic.create_staking_limit)


def test_ZB_NP_67(new_genesis_env, client_consensus, client_noconsensus):
    """
    零出块率次数小值变大值提案生效，达到共识轮位数，未达到零出块率次数（存在零出块记录）
    """
    update_zero_produce(new_genesis_env)
    economic = client_noconsensus.economic
    node = economic.env.get_consensus_node_by_index(0)

    # Voting failed and expired
    param_governance_verify_before_endblock(client_consensus, 'slashing', 'zeroProduceNumberThreshold', '4', True)

    # Start execution use case
    create_pledge_node_information(client_noconsensus)
    initial_validator, slashing_node_list = gen_validator_list(economic.env.consensus_node_id_list(), client_noconsensus.node.node_id)
    log.info("verifier nodeid list: {}".format(initial_validator))
    client_noconsensus.node.stop()

    set_slashing(initial_validator, slashing_node_list, node, economic, "111")
    assert_set_validator_list(node, initial_validator)

    economic.wait_consensus_blocknum(node, 1)
    log.info("Current block height: {}".format(node.eth.blockNumber))
    current_validator = get_pledge_list(node.ppos.getValidatorList)
    log.info("current validator: {}".format(current_validator))
    wait_slashing_list = node.debug.getWaitSlashingNodeList()
    log.info("Zero block record table：{}".format(wait_slashing_list))
    assert get_slash_count(wait_slashing_list, client_noconsensus.node.node_id) == to_int(111)
    assert_not_slashing(node.ppos.getCandidateInfo(client_noconsensus.node.node_id), economic.create_staking_limit)

    economic.wait_consensus_blocknum(node)
    log.info("Current block height: {}".format(node.eth.blockNumber))
    current_validator = get_pledge_list(node.ppos.getValidatorList)
    log.info("current validator: {}".format(current_validator))
    wait_slashing_list = node.debug.getWaitSlashingNodeList()
    log.info("Zero block record table：{}".format(wait_slashing_list))
    assert get_slash_count(wait_slashing_list, client_noconsensus.node.node_id) == to_int(111)
    assert_not_slashing(node.ppos.getCandidateInfo(client_noconsensus.node.node_id), economic.create_staking_limit)


def test_ZB_NP_68(new_genesis_env, clients_consensus):
    """
    当前触发共识轮230块前出块，不零出块处罚
    """
    update_zero_produce(new_genesis_env, 2, 2)
    economic = clients_consensus[0].economic
    first_node = economic.env.get_consensus_node_by_index(0)
    log.info("first_node id :{}".format(first_node.node_id))
    second_node = economic.env.get_consensus_node_by_index(1)
    log.info("first_node id :{}".format(second_node.node_id))

    # Start execution use case
    current_validator = get_pledge_list(first_node.ppos.getValidatorList)
    log.info("current validator: {}".format(current_validator))
    time.sleep(2)
    log.info("Current block height: {}".format(first_node.eth.blockNumber))
    second_node.stop()
    economic.wait_consensus_blocknum(first_node, 1)
    log.info("Current block height: {}".format(first_node.eth.blockNumber))
    wait_slashing_list = first_node.debug.getWaitSlashingNodeList()
    log.info("Zero block record table：{}".format(wait_slashing_list))
    assert get_slash_count(wait_slashing_list, second_node.node_id) == to_int(1)
    assert_not_slashing(first_node.ppos.getCandidateInfo(second_node.node_id), 1500000000000000000000000)
    second_node.start(False)
    wait_block_number(first_node, 111)
    wait_slashing_list = first_node.debug.getWaitSlashingNodeList()
    assert get_slash_count(wait_slashing_list, second_node.node_id) == 0
    assert_not_slashing(first_node.ppos.getCandidateInfo(second_node.node_id), 1500000000000000000000000)


def test_ZB_NP_69(new_genesis_env, clients_consensus):
    """
    当前触发共识轮230块后出块，零出块处罚
    """
    update_zero_produce(new_genesis_env, 2, 2)
    economic = clients_consensus[0].economic
    first_node = economic.env.get_consensus_node_by_index(0)
    log.info("first_node id :{}".format(first_node.node_id))
    second_node = economic.env.get_consensus_node_by_index(3)
    log.info("second_node id :{}".format(second_node.node_id))

    # Start execution use case
    current_validator = get_pledge_list(first_node.ppos.getValidatorList)
    log.info("current validator: {}".format(current_validator))
    time.sleep(2)
    log.info("Current block height: {}".format(first_node.eth.blockNumber))
    second_node.stop()
    economic.wait_consensus_blocknum(first_node, 1)
    log.info("Current block height: {}".format(first_node.eth.blockNumber))
    current_validator = get_pledge_list(first_node.ppos.getValidatorList)
    log.info("current validator: {}".format(current_validator))
    wait_slashing_list = first_node.debug.getWaitSlashingNodeList()
    log.info("Zero block record table：{}".format(wait_slashing_list))
    assert get_slash_count(wait_slashing_list, second_node.node_id) == to_int(1)
    assert_not_slashing(first_node.ppos.getCandidateInfo(second_node.node_id), 1500000000000000000000000)
    wait_block_number(first_node, 100)
    second_node.start(False)
    wait_slashing_list = first_node.debug.getWaitSlashingNodeList()
    assert get_slash_count(wait_slashing_list, second_node.node_id) == 0
    assert_slashing(first_node.ppos.getCandidateInfo(second_node.node_id), 1500000000000000000000000)


# def test_ZB_NP_70(new_genesis_env, clients_noconsensus):
#     """
#     同时触发多个节点零出块
#     """
#     genesis = from_dict(data_class=Genesis, data=new_genesis_env.genesis_config)
#     genesis.economicModel.common.maxEpochMinutes = 5
#     genesis.economicModel.common.maxConsensusVals = 7
#     genesis.economicModel.staking.maxValidators = 7
#     genesis.economicModel.slashing.zeroProduceNumberThreshold = 2
#     genesis.economicModel.slashing.zeroProduceCumulativeTime = 4
#     new_file = new_genesis_env.cfg.env_tmp + "/genesis_0.13.0.json"
#     genesis.to_file(new_file)
#     new_genesis_env.deploy_all(new_file)
#     economic_old = clients_noconsensus[0].economic
#     economic = copy.copy(economic_old)
#     economic.validator_count = 7
#     economic.expected_minutes = 5
#     node = economic.env.get_consensus_node_by_index(0)
#     # Start execution use case
#     for i in range(len(clients_noconsensus)):
#         client = clients_noconsensus[i]
#         log.info("Current connection node：{}".format(client.node.node_mark))
#         log.info("Start creating a pledge account Pledge_address")
#         staking_amount = von_amount(client.economic.create_staking_limit, 2)
#         staking_address, pri_key = client.economic.account.generate_account(client.node.web3, staking_amount)
#         log.info("Created, account address: {} Amount: {}".format(staking_address, staking_amount))
#         log.info("Start applying for a pledge node")
#         result = client.staking.create_staking(0, staking_address, staking_address)
#         assert_code(result, 0)
#         result = client.node.ppos.getCandidateInfo(client.node.node_id)
#         log.info("Candidate Info:{}".format(result))
    # economic.wait_settlement_blocknum(node)
    # log.info("Current block height: {}".format(node.eth.blockNumber))
    # result = node.ppos.getValidatorList()
    # log.info("current Validator List：{}".format(result))
    # clients_noconsensus[1].node.stop()
    # clients_noconsensus[2].node.stop()
    #
    # economic.wait_consensus_blocknum(node)
    # log.info("Current block height: {}".format(node.eth.blockNumber))
    # log.info("Validator List:{}".format(node.ppos.getValidatorList()))
    # wait_slashing_list = node.debug.getWaitSlashingNodeList()
    # log.info("Zero block record table：{}".format(wait_slashing_list))
    # assert get_slash_count(wait_slashing_list, clients_noconsensus[1].node.node_id) == 0
    # assert get_slash_count(wait_slashing_list, clients_noconsensus[2].node.node_id) == 0
    #
    # economic.wait_consensus_blocknum(node)
    # log.info("Current block height: {}".format(node.eth.blockNumber))
    # log.info("Validator List:{}".format(node.ppos.getValidatorList()))
    # wait_slashing_list = node.debug.getWaitSlashingNodeList()
    # log.info("Zero block record table：{}".format(wait_slashing_list))
    # assert get_slash_count(wait_slashing_list, clients_noconsensus[1].node.node_id) == 0
    # assert get_slash_count(wait_slashing_list, clients_noconsensus[2].node.node_id) == 0
    #
    # economic.wait_consensus_blocknum(node)
    # log.info("Current block height: {}".format(node.eth.blockNumber))
    # log.info("Validator List:{}".format(node.ppos.getValidatorList()))
    # wait_slashing_list = node.debug.getWaitSlashingNodeList()
    # log.info("Zero block record table：{}".format(wait_slashing_list))
    # assert get_slash_count(wait_slashing_list, clients_noconsensus[1].node.node_id) == to_int(1)
    # assert get_slash_count(wait_slashing_list, clients_noconsensus[2].node.node_id) == to_int(1)
    #
    # economic.wait_consensus_blocknum(node)
    # log.info("Current block height: {}".format(node.eth.blockNumber))
    # log.info("Validator List:{}".format(node.ppos.getValidatorList()))
    # wait_slashing_list = node.debug.getWaitSlashingNodeList()
    # log.info("Zero block record table：{}".format(wait_slashing_list))
    # assert get_slash_count(wait_slashing_list, clients_noconsensus[1].node.node_id) == to_int(11)
    # assert get_slash_count(wait_slashing_list, clients_noconsensus[2].node.node_id) == to_int(11)
    #
    # economic.wait_consensus_blocknum(node)
    # log.info("Current block height: {}".format(node.eth.blockNumber))
    # log.info("Validator List:{}".format(node.ppos.getValidatorList()))
    # wait_slashing_list = node.debug.getWaitSlashingNodeList()
    # log.info("Zero block record table：{}".format(wait_slashing_list))
    # assert get_slash_count(wait_slashing_list, clients_noconsensus[1].node.node_id) == to_int(111)
    # assert get_slash_count(wait_slashing_list, clients_noconsensus[2].node.node_id) == to_int(111)
    #
    # economic.wait_consensus_blocknum(node)
    # log.info("Current block height: {}".format(node.eth.blockNumber))
    # log.info("Validator List:{}".format(node.ppos.getValidatorList()))
    # wait_slashing_list = node.debug.getWaitSlashingNodeList()
    # log.info("Zero block record table：{}".format(wait_slashing_list))
    # assert get_slash_count(wait_slashing_list, clients_noconsensus[1].node.node_id) == 0
    # assert get_slash_count(wait_slashing_list, clients_noconsensus[2].node.node_id) == 0
    # assert_slashing(node.ppos.getCandidateInfo(clients_noconsensus[1].node.node_id), economic.create_staking_limit)
    # assert_slashing(node.ppos.getCandidateInfo(clients_noconsensus[2].node.node_id), economic.create_staking_limit)
    # economic.wait_consensus_blocknum(node)
    # log.info("Current block height: {}".format(node.eth.blockNumber))
    # economic.wait_consensus_blocknum(node)
    # log.info("Current block height: {}".format(node.eth.blockNumber))

#
# def test_71(new_genesis_env, clients_consensus, clients_noconsensus):
#     """
#     随机生成零出块
#     """
#     genesis = from_dict(data_class=Genesis, data=new_genesis_env.genesis_config)
#     genesis.economicModel.staking.maxValidators = 7
#     genesis.economicModel.slashing.zeroProduceNumberThreshold = 3
#     genesis.economicModel.slashing.zeroProduceCumulativeTime = 4
#     new_file = new_genesis_env.cfg.env_tmp + "/genesis_0.13.0.json"
#     genesis.to_file(new_file)
#     new_genesis_env.deploy_all(new_file)
#     economic = clients_consensus[0].economic
#     node = economic.env.get_consensus_node_by_index(0)
#
#     # Start execution use case
#     for i in range(len(clients_noconsensus)):
#         client = clients_noconsensus[i]
#         log.info("Current connection node：{}".format(client.node.node_mark))
#         log.info("Start creating a pledge account Pledge_address")
#         staking_amount = von_amount(client.economic.create_staking_limit, 6)
#         staking_address, pri_key = client.economic.account.generate_account(client.node.web3, staking_amount)
#         log.info("Created, account address: {} Amount: {}".format(staking_address, staking_amount))
#         log.info("Start applying for a pledge node")
#         result = client.staking.create_staking(0, staking_address, staking_address, amount=von_amount(economic.create_staking_limit, 2))
#         assert_code(result, 0)
#         result = client.node.ppos.getCandidateInfo(client.node.node_id)
#         log.info("Candidate Info:{}".format(result))
#     economic.wait_settlement_blocknum(node)
#     log.info("Current block height: {}".format(node.eth.blockNumber))
#     # initial_validator, slashing_node_list = gen_validator_list(economic.env.consensus_node_id_list(), client_noconsensus.node.node_id)
#     # log.info("verifier nodeid list: {}".format(initial_validator))
#     clients_consensus[2].node.stop()
#     clients_noconsensus[1].node.stop()
#
# # economic.wait_consensus_blocknum(node, 1)
#     # log.info("Current block height: {}".format(node.eth.blockNumber))
#     # client_noconsensus.node.start(False)
#
#     for i in range(8):
#         economic.wait_consensus_blocknum(node)
#         log.info("Current block height: {}".format(node.eth.blockNumber))
#         current_validator = node.ppos.getValidatorList()
#         log.info("current validator: {}".format(current_validator))
#         wait_slashing_list = node.debug.getWaitSlashingNodeList()
#         log.info("Zero block record table：{}".format(wait_slashing_list))
#         result = node.ppos.getCandidateInfo(clients_consensus[2].node.node_id)
#         log.info("Candidate Info:{}".format(result))
#         result = node.ppos.getCandidateInfo(clients_noconsensus[1].node.node_id)
#         log.info("Candidate Info:{}".format(result))
