import time
import pytest
from dacite import from_dict
from common.log import log
from tests.lib.utils import assert_code
from tests.lib.utils import wait_block_number, get_pledge_list
from tests.ppos.conftest import create_staking
from tests.ppos.conftest import calculate
from tests.lib.genesis import Genesis

pri_key = ""


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


@pytest.fixture()
def update_zero_produce_env(global_test_env):
    global pri_key
    pri_key = global_test_env.account.account_with_money["prikey"]
    genesis_config = global_test_env.genesis_config
    genesis = from_dict(data_class=Genesis, data=genesis_config)
    genesis.economicModel.slashing.zeroProduceCumulativeTime = 4
    genesis.economicModel.slashing.zeroProduceNumberThreshold = 3
    genesis_path = global_test_env.cfg.env_tmp + "/genesis_0.13.0.json"

    genesis.to_file(genesis_path)

    global_test_env.deploy_all(genesis_path)
    yield global_test_env
    global_test_env.deploy_all()


@pytest.fixture()
def new_validator_client(update_zero_produce_env, client_noconsensus):
    staking_address, _ = create_staking(client_noconsensus, 10)
    log.info("use node: {} node id: {}".format(client_noconsensus.node.url, client_noconsensus.node.node_id))
    setattr(client_noconsensus, "staking_address", staking_address)
    client_noconsensus.economic.wait_settlement_blocknum(client_noconsensus.node)
    yield client_noconsensus


# def test_case(new_validator_client):
#     new_validator_client.node.stop()
#     economic = new_validator_client.economic
#     node = economic.env.get_consensus_node_by_index(0)
#
#     def slashing(candidate_info):
#         ret = candidate_info["Ret"]
#         if ret["Status"] == 0 and ret["Released"] == new_validator_client.staking_amount:
#             return False
#         return True
#
#     for i in range(1000):
#         if i/2 == 0 and not new_validator_client.node.running:
#             new_validator_client.node.start(False)
#             economic.wait_consensus_blocknum(node)
#             new_validator_client.node.stop()
#         else:
#             economic.wait_consensus_blocknum(node)
#         print(i, get_slash_count(node.debug.getWaitSlashingNodeList(), new_validator_client.node.node_id))
#         if slashing(node.ppos.getCandidateInfo(new_validator_client.node.node_id)):
#             print("slashing")
#             break


def test_ZB_NP_11(new_validator_client):
    new_validator_client.node.stop()
    economic = new_validator_client.economic
    node = economic.env.get_consensus_node_by_index(0)

    pri = economic.env.account.find_pri_key(new_validator_client.staking_address)
    node.ppos.withdrewStaking(new_validator_client.node.node_id, pri)

    initial_validator, slashing_node_list = gen_validator_list(economic.env.consensus_node_id_list(), new_validator_client.node.node_id)

    set_slashing(initial_validator, slashing_node_list, node, economic, "111")
    economic.wait_consensus_blocknum(node, 1)
    wait_slashing_list = node.debug.getWaitSlashingNodeList()
    assert get_slash_count(wait_slashing_list, new_validator_client.node.node_id) == to_int(111)

    log.info("current validator: {}".format(get_pledge_list(node.ppos.getValidatorList)))
    economic.wait_consensus_blocknum(node)
    wait_slashing_list = node.debug.getWaitSlashingNodeList()
    assert len(wait_slashing_list) == 0
    assert_slashing(node.ppos.getCandidateInfo(new_validator_client.node.node_id), new_validator_client.staking_amount)


def test_ZB_NP_12(new_validator_client):
    new_validator_client.node.stop()
    economic = new_validator_client.economic
    node = economic.env.get_consensus_node_by_index(0)

    pri = economic.env.account.find_pri_key(new_validator_client.staking_address)
    node.ppos.withdrewStaking(new_validator_client.node.node_id, pri)

    initial_validator, slashing_node_list = gen_validator_list(economic.env.consensus_node_id_list(), new_validator_client.node.node_id)

    set_slashing(initial_validator, slashing_node_list, node, economic, "1")
    economic.wait_consensus_blocknum(node, 1)

    wait_slashing_list = node.debug.getWaitSlashingNodeList()
    assert get_slash_count(wait_slashing_list, new_validator_client.node.node_id) == to_int(1)
    assert node.ppos.getCandidateInfo(new_validator_client.node.node_id)["Ret"]["Released"] == new_validator_client.staking_amount


def test_ZB_NP_13(new_validator_client):
    new_validator_client.node.stop()
    economic = new_validator_client.economic
    node = economic.env.get_consensus_node_by_index(0)

    initial_validator, slashing_node_list = gen_validator_list(economic.env.consensus_node_id_list(), new_validator_client.node.node_id)

    set_slashing(initial_validator, slashing_node_list, node, economic, "111")
    economic.wait_consensus_blocknum(node, 1)

    wait_slashing_list = node.debug.getWaitSlashingNodeList()
    assert get_slash_count(wait_slashing_list, new_validator_client.node.node_id) == to_int(111)
    pri = economic.env.account.find_pri_key(new_validator_client.staking_address)
    node.ppos.withdrewStaking(new_validator_client.node.node_id, pri)
    log.info("current validator: {}".format(get_pledge_list(node.ppos.getValidatorList)))
    economic.wait_consensus_blocknum(node)
    wait_slashing_list = node.debug.getWaitSlashingNodeList()
    assert len(wait_slashing_list) == 0
    assert_slashing(node.ppos.getCandidateInfo(new_validator_client.node.node_id), new_validator_client.staking_amount)


def test_ZB_NP_14_15(new_validator_client):
    new_validator_client.node.stop()
    economic = new_validator_client.economic
    node = economic.env.get_consensus_node_by_index(0)

    initial_validator, slashing_node_list = gen_validator_list(economic.env.consensus_node_id_list(), new_validator_client.node.node_id)

    assert_set_validator_list(node, slashing_node_list)
    current_validator = get_pledge_list(node.ppos.getValidatorList)
    log.info("current validator: {}".format(current_validator))
    economic.wait_consensus_blocknum(node)
    wait_slashing_list = node.debug.getWaitSlashingNodeList()
    assert len(wait_slashing_list) == 0
    assert_set_validator_list(node, initial_validator)
    current_validator = get_pledge_list(node.ppos.getValidatorList)
    log.info("current validator: {}".format(current_validator))
    assert new_validator_client.node.node_id in current_validator
    economic.wait_consensus_blocknum(node)
    # 需要在等一个共识轮才能查看到待处罚信息
    assert_set_validator_list(node, slashing_node_list)
    economic.wait_consensus_blocknum(node)
    wait_slashing_list = node.debug.getWaitSlashingNodeList()
    assert get_slash_count(wait_slashing_list, new_validator_client.node.node_id) == to_int(1)

    assert_set_validator_list(node, initial_validator)
    economic.wait_consensus_blocknum(node)
    wait_slashing_list = node.debug.getWaitSlashingNodeList()
    assert get_slash_count(wait_slashing_list, new_validator_client.node.node_id) == to_int(1)

    assert_set_validator_list(node, initial_validator)
    economic.wait_consensus_blocknum(node)
    wait_slashing_list = node.debug.getWaitSlashingNodeList()
    assert get_slash_count(wait_slashing_list, new_validator_client.node.node_id) == to_int(101)

    assert_set_validator_list(node, initial_validator)
    economic.wait_consensus_blocknum(node)
    wait_slashing_list = node.debug.getWaitSlashingNodeList()
    assert get_slash_count(wait_slashing_list, new_validator_client.node.node_id) == to_int(101)

    economic.wait_consensus_blocknum(node)
    wait_slashing_list = node.debug.getWaitSlashingNodeList()
    assert get_slash_count(wait_slashing_list, new_validator_client.node.node_id) == to_int(1)
    assert_not_slashing(node.ppos.getCandidateInfo(new_validator_client.node.node_id), new_validator_client.staking_amount)


def test_ZB_NP_16(new_validator_client):
    new_validator_client.node.stop()
    economic = new_validator_client.economic
    node = economic.env.get_consensus_node_by_index(0)

    initial_validator, slashing_node_list = gen_validator_list(economic.env.consensus_node_id_list(), new_validator_client.node.node_id)

    set_slashing(initial_validator, slashing_node_list, node, economic, "111")
    economic.wait_consensus_blocknum(node, 1)

    wait_slashing_list = node.debug.getWaitSlashingNodeList()
    assert get_slash_count(wait_slashing_list, new_validator_client.node.node_id) == to_int(111)
    assert_not_slashing(node.ppos.getCandidateInfo(new_validator_client.node.node_id), new_validator_client.staking_amount)


def test_ZB_NP_17_20(new_validator_client):
    new_validator_client.node.stop()
    economic = new_validator_client.economic
    node = economic.env.get_consensus_node_by_index(0)

    initial_validator, slashing_node_list = gen_validator_list(economic.env.consensus_node_id_list(), new_validator_client.node.node_id)

    set_slashing(initial_validator, slashing_node_list, node, economic, "1000")

    assert_set_validator_list(node, initial_validator)
    economic.wait_consensus_blocknum(node, 1)

    wait_slashing_list = node.debug.getWaitSlashingNodeList()
    assert get_slash_count(wait_slashing_list, new_validator_client.node.node_id) == to_int(1)

    economic.wait_consensus_blocknum(node)
    wait_slashing_list = node.debug.getWaitSlashingNodeList()
    assert len(wait_slashing_list) == 0
    assert_not_slashing(node.ppos.getCandidateInfo(new_validator_client.node.node_id), new_validator_client.staking_amount)


def test_ZB_NP_18(new_validator_client):
    new_validator_client.node.stop()
    economic = new_validator_client.economic
    node = economic.env.get_consensus_node_by_index(0)

    initial_validator, slashing_node_list = gen_validator_list(economic.env.consensus_node_id_list(), new_validator_client.node.node_id)

    set_slashing(initial_validator, slashing_node_list, node, economic, "1000")
    assert_set_validator_list(node, slashing_node_list)
    economic.wait_consensus_blocknum(node, 1)

    wait_slashing_list = node.debug.getWaitSlashingNodeList()
    assert get_slash_count(wait_slashing_list, new_validator_client.node.node_id) == to_int(1)

    economic.wait_consensus_blocknum(node)
    wait_slashing_list = node.debug.getWaitSlashingNodeList()
    assert get_slash_count(wait_slashing_list, new_validator_client.node.node_id) == to_int(1)
    assert_not_slashing(node.ppos.getCandidateInfo(new_validator_client.node.node_id), new_validator_client.staking_amount)


def test_ZB_NP_19(new_validator_client):
    new_validator_client.node.stop()
    economic = new_validator_client.economic
    node = economic.env.get_consensus_node_by_index(0)

    initial_validator, slashing_node_list = gen_validator_list(economic.env.consensus_node_id_list(), new_validator_client.node.node_id)

    set_slashing(initial_validator, slashing_node_list, node, economic, "1000")
    assert_set_validator_list(node, slashing_node_list)
    economic.wait_consensus_blocknum(node)
    new_validator_client.node.start(False)
    economic.wait_consensus_blocknum(node)

    wait_slashing_list = node.debug.getWaitSlashingNodeList()
    assert get_slash_count(wait_slashing_list, new_validator_client.node.node_id) == to_int(1)

    economic.wait_consensus_blocknum(node)
    wait_slashing_list = node.debug.getWaitSlashingNodeList()
    assert len(wait_slashing_list) == 0
    assert_not_slashing(node.ppos.getCandidateInfo(new_validator_client.node.node_id), new_validator_client.staking_amount)


def test_ZB_NP_21_22(new_validator_client):
    new_validator_client.node.stop()
    economic = new_validator_client.economic
    node = economic.env.get_consensus_node_by_index(0)

    initial_validator, slashing_node_list = gen_validator_list(economic.env.consensus_node_id_list(), new_validator_client.node.node_id)

    set_slashing(initial_validator, slashing_node_list, node, economic, "1001")
    assert_set_validator_list(node, slashing_node_list)
    economic.wait_consensus_blocknum(node, 1)

    wait_slashing_list = node.debug.getWaitSlashingNodeList()
    assert get_slash_count(wait_slashing_list, new_validator_client.node.node_id) == to_int(1001)

    economic.wait_consensus_blocknum(node)
    wait_slashing_list = node.debug.getWaitSlashingNodeList()
    assert get_slash_count(wait_slashing_list, new_validator_client.node.node_id) == to_int(11)
    assert_not_slashing(node.ppos.getCandidateInfo(new_validator_client.node.node_id), new_validator_client.staking_amount)


def test_ZB_NP_23(new_validator_client):
    new_validator_client.node.stop()
    economic = new_validator_client.economic
    node = economic.env.get_consensus_node_by_index(0)

    initial_validator, slashing_node_list = gen_validator_list(economic.env.consensus_node_id_list(), new_validator_client.node.node_id)

    set_slashing(initial_validator, slashing_node_list, node, economic, "1001")
    assert_set_validator_list(node, slashing_node_list)
    economic.wait_consensus_blocknum(node)
    new_validator_client.node.start(False)
    economic.wait_consensus_blocknum(node)

    wait_slashing_list = node.debug.getWaitSlashingNodeList()
    assert get_slash_count(wait_slashing_list, new_validator_client.node.node_id) == to_int(1001)

    economic.wait_consensus_blocknum(node)
    print(new_validator_client.node.block_number)
    wait_slashing_list = node.debug.getWaitSlashingNodeList()
    assert len(wait_slashing_list) == 0
    assert_not_slashing(node.ppos.getCandidateInfo(new_validator_client.node.node_id), new_validator_client.staking_amount)


def test_ZB_NP_24(new_validator_client):
    new_validator_client.node.stop()
    economic = new_validator_client.economic
    node = economic.env.get_consensus_node_by_index(0)

    initial_validator, slashing_node_list = gen_validator_list(economic.env.consensus_node_id_list(), new_validator_client.node.node_id)

    set_slashing(initial_validator, slashing_node_list, node, economic, "1001")
    assert_set_validator_list(node, initial_validator)
    economic.wait_consensus_blocknum(node, 1)

    wait_slashing_list = node.debug.getWaitSlashingNodeList()
    assert get_slash_count(wait_slashing_list, new_validator_client.node.node_id) == to_int(1001)

    economic.wait_consensus_blocknum(node)
    wait_slashing_list = node.debug.getWaitSlashingNodeList()
    assert get_slash_count(wait_slashing_list, new_validator_client.node.node_id) == to_int(1)
    assert_not_slashing(node.ppos.getCandidateInfo(new_validator_client.node.node_id), new_validator_client.staking_amount)


def test_ZB_NP_25(new_validator_client):
    new_validator_client.node.stop()
    economic = new_validator_client.economic
    node = economic.env.get_consensus_node_by_index(0)

    initial_validator, slashing_node_list = gen_validator_list(economic.env.consensus_node_id_list(), new_validator_client.node.node_id)

    set_slashing(initial_validator, slashing_node_list, node, economic, "111")
    assert_set_validator_list(node, economic.env.consensus_node_id_list())
    economic.wait_consensus_blocknum(node, 1)

    wait_slashing_list = node.debug.getWaitSlashingNodeList()
    assert get_slash_count(wait_slashing_list, new_validator_client.node.node_id) == to_int(111)

    economic.wait_consensus_blocknum(node, 1)
    wait_slashing_list = node.debug.getWaitSlashingNodeList()
    assert len(wait_slashing_list) == 0
    assert_slashing(node.ppos.getCandidateInfo(new_validator_client.node.node_id), new_validator_client.staking_amount)


def test_ZB_NP_26(new_validator_client):
    new_validator_client.node.stop()
    economic = new_validator_client.economic
    node = economic.env.get_consensus_node_by_index(0)

    initial_validator, slashing_node_list = gen_validator_list(economic.env.consensus_node_id_list(), new_validator_client.node.node_id)

    set_slashing(initial_validator, slashing_node_list, node, economic, "111")
    assert_set_validator_list(node, slashing_node_list)
    economic.wait_consensus_blocknum(node, 1)

    wait_slashing_list = node.debug.getWaitSlashingNodeList()
    assert get_slash_count(wait_slashing_list, new_validator_client.node.node_id) == to_int(111)
    economic.wait_consensus_blocknum(node)
    wait_slashing_list = node.debug.getWaitSlashingNodeList()
    assert len(wait_slashing_list) == 0
    assert_slashing(node.ppos.getCandidateInfo(new_validator_client.node.node_id), new_validator_client.staking_amount)


def test_ZB_NP_27(new_validator_client):
    new_validator_client.node.stop()
    economic = new_validator_client.economic
    node = economic.env.get_consensus_node_by_index(0)

    initial_validator, slashing_node_list = gen_validator_list(economic.env.consensus_node_id_list(), new_validator_client.node.node_id)

    set_slashing(initial_validator, slashing_node_list, node, economic, "101")
    assert_set_validator_list(node, slashing_node_list)
    economic.wait_consensus_blocknum(node, 1)

    wait_slashing_list = node.debug.getWaitSlashingNodeList()
    assert get_slash_count(wait_slashing_list, new_validator_client.node.node_id) == to_int(101)
    economic.wait_consensus_blocknum(node)
    wait_slashing_list = node.debug.getWaitSlashingNodeList()
    assert len(wait_slashing_list) == 0
    assert_slashing(node.ppos.getCandidateInfo(new_validator_client.node.node_id), new_validator_client.staking_amount)


def test_ZB_NP_28(new_validator_client):
    new_validator_client.node.stop()
    economic = new_validator_client.economic
    node = economic.env.get_consensus_node_by_index(0)

    initial_validator, slashing_node_list = gen_validator_list(economic.env.consensus_node_id_list(), new_validator_client.node.node_id)

    set_slashing(initial_validator, slashing_node_list, node, economic, "110")
    assert_set_validator_list(node, slashing_node_list)
    economic.wait_consensus_blocknum(node, 1)

    wait_slashing_list = node.debug.getWaitSlashingNodeList()
    assert get_slash_count(wait_slashing_list, new_validator_client.node.node_id) == to_int(11)

    economic.wait_consensus_blocknum(node)
    wait_slashing_list = node.debug.getWaitSlashingNodeList()
    assert len(wait_slashing_list) == 0
    assert_slashing(node.ppos.getCandidateInfo(new_validator_client.node.node_id), new_validator_client.staking_amount)


def test_ZB_NP_29(update_zero_produce_env, clients_noconsensus):
    client_a = clients_noconsensus[0]
    client_b = clients_noconsensus[1]
    amount = calculate(client_a.economic.create_staking_limit, 5)
    staking_amount_a = calculate(client_a.economic.create_staking_limit, 1)
    staking_amount_b = calculate(client_a.economic.create_staking_limit, 2)
    staking_address, _ = client_a.economic.account.generate_account(client_a.node.web3, amount)
    result = client_a.staking.create_staking(0, staking_address, staking_address, amount=staking_amount_a, reward_per=10)
    assert_code(result, 0)
    economic = client_b.economic
    economic.wait_settlement_blocknum(client_b.node)
    client_a.node.stop()

    economic = client_b.economic
    node = economic.env.get_consensus_node_by_index(0)

    initial_validator, slashing_node_list = gen_validator_list(economic.env.consensus_node_id_list(), client_a.node.node_id)

    assert_set_validator_list(node, initial_validator)
    current_validator = get_pledge_list(node.ppos.getValidatorList)
    log.info("current validator: {}".format(current_validator))
    economic.wait_consensus_blocknum(node, 1)
    result = client_b.staking.create_staking(0, staking_address, staking_address, amount=staking_amount_b, reward_per=10)
    assert_code(result, 0)
    assert_set_validator_list(node, slashing_node_list)
    economic.wait_consensus_blocknum(node)
    assert_set_validator_list(node, initial_validator)
    economic.wait_consensus_blocknum(node)
    economic.wait_consensus_blocknum(node, 3)
    wait_slashing_list = node.debug.getWaitSlashingNodeList()
    assert get_slash_count(wait_slashing_list, client_a.node.node_id) == to_int(1)

    economic.wait_consensus_blocknum(node)
    wait_slashing_list = node.debug.getWaitSlashingNodeList()
    assert len(wait_slashing_list) == 0
    assert_not_slashing(node.ppos.getCandidateInfo(client_a.node.node_id), staking_amount_a)


def test_ZB_NP_30(new_validator_client):
    new_validator_client.node.stop()
    economic = new_validator_client.economic
    node = economic.env.get_consensus_node_by_index(0)

    initial_validator, slashing_node_list = gen_validator_list(economic.env.consensus_node_id_list(), new_validator_client.node.node_id)

    assert_set_validator_list(node, slashing_node_list)
    current_validator = get_pledge_list(node.ppos.getValidatorList)
    log.info("current validator: {}".format(current_validator))
    economic.wait_consensus_blocknum(node)

    assert_set_validator_list(node, slashing_node_list)
    economic.wait_consensus_blocknum(node)
    num = economic.get_consensus_switchpoint(node, 1)
    while node.block_number < num:
        new_validator_client.node.start(False)
        time.sleep(4)
        print(new_validator_client.node.block_number, node.block_number)
        new_validator_client.node.stop()
        time.sleep(2)
    wait_slashing_list = node.debug.getWaitSlashingNodeList()
    assert len(wait_slashing_list) == 0
    assert_not_slashing(node.ppos.getCandidateInfo(new_validator_client.node.node_id), new_validator_client.staking_amount)

