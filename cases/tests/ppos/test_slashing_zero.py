import allure
import pytest
from dacite import from_dict
from common.log import log
from tests.lib.utils import wait_block_number, get_pledge_list
from tests.ppos.conftest import create_staking

from tests.lib.genesis import Genesis


def assert_set_validator_list(node, validator_list):
    result = node.debug.setValidatorList(validator_list)
    if result is not None:
        raise Exception("setValidatorList:{}".format(result))


def get_slash_count(wait_slashing_node_list, node_id):
    for slashing in wait_slashing_node_list:
        if slashing["NodeId"] == node_id:
            return slashing["CountBit"]
    return 0


def to_int(value):
    return int(str(value), 2)


@pytest.fixture()
def update_zero_produce(global_test_env, client_noconsensus):
    genesis_config = global_test_env.genesis_config
    genesis = from_dict(data_class=Genesis, data=genesis_config)
    genesis.economicModel.slashing.zeroProduceCumulativeTime = 4
    genesis.economicModel.slashing.zeroProduceNumberThreshold = 3
    genesis_path = global_test_env.cfg.env_tmp + "/genesis.json"

    genesis.to_file(genesis_path)

    global_test_env.deploy_all(genesis_path)
    staking_address, _ = create_staking(client_noconsensus, 10)
    log.info("use node: {}-node id: {}".format(client_noconsensus.node.url, client_noconsensus.node.node_id))
    setattr(client_noconsensus, "staking_address", staking_address)
    yield client_noconsensus
    global_test_env.deploy_all()


@pytest.fixture()
def new_validator_client(update_zero_produce):
    update_zero_produce.economic.wait_settlement_blocknum(update_zero_produce.node)
    yield update_zero_produce


# 出现次数为循环次数减2
def test_ZB_NP_11(new_validator_client):
    new_validator_client.node.stop()
    economic = new_validator_client.economic
    node = economic.env.get_consensus_node_by_index(0)

    pri = economic.env.account.find_pri_key(new_validator_client.staking_address)
    node.ppos.withdrewStaking(new_validator_client.node.node_id, pri)

    consensus_list = economic.env.consensus_node_id_list()
    consensus_list.append(new_validator_client.node.node_id)
    consensus_list.pop(0)

    for i in range(5):
        node.debug.setValidatorList(consensus_list)
        current_validator = get_pledge_list(node.ppos.getValidatorList)
        log.info("current validator: {}".format(get_pledge_list(node.ppos.getValidatorList)))
        if i != 0 and new_validator_client.node.node_id not in current_validator:
            raise Exception("node not in to validator")
        economic.wait_consensus_blocknum(node)
        wait_slashing_list = node.debug.getWaitSlashingNodeList()
        print(wait_slashing_list)

    wait_slashing_list = node.debug.getWaitSlashingNodeList()
    assert get_slash_count(wait_slashing_list, new_validator_client.node.node_id) == to_int(111)

    log.info("current validator: {}".format(get_pledge_list(node.ppos.getValidatorList)))
    economic.wait_consensus_blocknum(node)
    wait_slashing_list = node.debug.getWaitSlashingNodeList()
    print(wait_slashing_list)


def test_ZB_NP_12(new_validator_client):
    new_validator_client.node.stop()
    economic = new_validator_client.economic
    node = economic.env.get_consensus_node_by_index(0)

    pri = economic.env.account.find_pri_key(new_validator_client.staking_address)
    node.ppos.withdrewStaking(new_validator_client.node.node_id, pri)

    consensus_list = economic.env.consensus_node_id_list()
    consensus_list.append(new_validator_client.node.node_id)
    consensus_list.pop(0)

    for i in range(3):
        node.debug.setValidatorList(consensus_list)
        current_validator = get_pledge_list(node.ppos.getValidatorList)
        log.info("current validator: {}".format(get_pledge_list(node.ppos.getValidatorList)))
        if i != 0 and new_validator_client.node.node_id not in current_validator:
            raise Exception("node not in to validator")
        economic.wait_consensus_blocknum(node)

    wait_slashing_list = node.debug.getWaitSlashingNodeList()
    assert get_slash_count(wait_slashing_list, new_validator_client.node.node_id) == to_int(1)


def test_ZB_NP_13(new_validator_client):
    new_validator_client.node.stop()
    economic = new_validator_client.economic
    node = economic.env.get_consensus_node_by_index(0)

    consensus_list = economic.env.consensus_node_id_list()
    consensus_list.append(new_validator_client.node.node_id)
    consensus_list.pop(0)

    for i in range(5):
        node.debug.setValidatorList(consensus_list)
        current_validator = get_pledge_list(node.ppos.getValidatorList)
        log.info("current validator: {}".format(get_pledge_list(node.ppos.getValidatorList)))
        if i != 0 and new_validator_client.node.node_id not in current_validator:
            raise Exception("node not in to validator")
        economic.wait_consensus_blocknum(node)
        wait_slashing_list = node.debug.getWaitSlashingNodeList()
        print(wait_slashing_list)

    wait_slashing_list = node.debug.getWaitSlashingNodeList()
    assert get_slash_count(wait_slashing_list, new_validator_client.node.node_id) == to_int(111)
    pri = economic.env.account.find_pri_key(new_validator_client.staking_address)
    node.ppos.withdrewStaking(new_validator_client.node.node_id, pri)
    log.info("current validator: {}".format(get_pledge_list(node.ppos.getValidatorList)))
    economic.wait_consensus_blocknum(node)
    wait_slashing_list = node.debug.getWaitSlashingNodeList()
    print(wait_slashing_list)


@pytest.mark.parametrize('value', [1, 2])
def test_ZB_NP_14_15(new_validator_client, value):
    new_validator_client.node.stop()
    economic = new_validator_client.economic
    node = economic.env.get_consensus_node_by_index(0)

    consensus_list = economic.env.consensus_node_id_list()
    consensus_list.append(new_validator_client.node.node_id)
    consensus_list.pop(0)

    for i in range(5):
        node.debug.setValidatorList(consensus_list)
        current_validator = get_pledge_list(node.ppos.getValidatorList)
        log.info("current validator: {}".format(get_pledge_list(node.ppos.getValidatorList)))
        if i != 0 and new_validator_client.node.node_id not in current_validator:
            raise Exception("node not in to validator")
        economic.wait_consensus_blocknum(node)
        wait_slashing_list = node.debug.getWaitSlashingNodeList()
        print(wait_slashing_list)
        if i == value + 1:
            assert get_slash_count(wait_slashing_list, new_validator_client.node.node_id) == to_int(value)



def test_ZB_NP_16(new_validator_client):
    new_validator_client.node.stop()
    print(new_validator_client.node.node_id)
    economic = new_validator_client.economic
    node = economic.env.get_consensus_node_by_index(0)
    consensus_list = economic.env.consensus_node_id_list()
    consensus_list.append(new_validator_client.node.node_id)
    consensus_list.pop(0)
    for i in range(4):
        node.debug.setValidatorList(consensus_list)
        print(get_pledge_list(node.ppos.getValidatorList))
        economic.wait_consensus_blocknum(node)

        wait_slashing_list = node.debug.getWaitSlashingNodeList()
        print(wait_slashing_list)


def test_ZB_NP_17(new_validator_client):
    new_validator_client.node.stop()
    print(new_validator_client.node.node_id)
    economic = new_validator_client.economic
    node = economic.env.get_consensus_node_by_index(0)
    consensus_list = economic.env.consensus_node_id_list()
    consensus_list.append(new_validator_client.node.node_id)
    consensus_list.pop(0)
    for i in range(4):
        node.debug.setValidatorList(consensus_list)
        print(get_pledge_list(node.ppos.getValidatorList))
        economic.wait_consensus_blocknum(node)

        wait_slashing_list = node.debug.getWaitSlashingNodeList()
        print(wait_slashing_list)


def test_ZB_NP_18(new_validator_client):
    new_validator_client.node.stop()
    print(new_validator_client.node.node_id)
    economic = new_validator_client.economic
    node = economic.env.get_consensus_node_by_index(0)
    consensus_list = economic.env.consensus_node_id_list()
    consensus_list.append(new_validator_client.node.node_id)
    consensus_list.pop(0)
    for i in range(4):
        node.debug.setValidatorList(consensus_list)
        print(get_pledge_list(node.ppos.getValidatorList))
        economic.wait_consensus_blocknum(node)

        wait_slashing_list = node.debug.getWaitSlashingNodeList()
        print(wait_slashing_list)


def test_ZB_NP_19(new_validator_client):
    new_validator_client.node.stop()
    print(new_validator_client.node.node_id)
    economic = new_validator_client.economic
    node = economic.env.get_consensus_node_by_index(0)
    consensus_list = economic.env.consensus_node_id_list()
    consensus_list.append(new_validator_client.node.node_id)
    consensus_list.pop(0)
    for i in range(4):
        node.debug.setValidatorList(consensus_list)
        print(get_pledge_list(node.ppos.getValidatorList))
        economic.wait_consensus_blocknum(node)

        wait_slashing_list = node.debug.getWaitSlashingNodeList()
        print(wait_slashing_list)
    new_validator_client.node.start(False)
    wait_slashing_list = node.debug.getWaitSlashingNodeList()
    print(wait_slashing_list)


def test_ZB_NP_20(new_validator_client):
    new_validator_client.node.stop()
    print(new_validator_client.node.node_id)
    economic = new_validator_client.economic
    node = economic.env.get_consensus_node_by_index(0)
    consensus_list = economic.env.consensus_node_id_list()
    consensus_list.append(new_validator_client.node.node_id)
    consensus_list.pop(0)
    for i in range(4):
        node.debug.setValidatorList(consensus_list)
        print(get_pledge_list(node.ppos.getValidatorList))
        economic.wait_consensus_blocknum(node)

        wait_slashing_list = node.debug.getWaitSlashingNodeList()
        print(wait_slashing_list)
    new_validator_client.node.start(False)
    wait_slashing_list = node.debug.getWaitSlashingNodeList()
    print(wait_slashing_list)


def test_ZB_NP_21(new_validator_client):
    new_validator_client.node.stop()
    print(new_validator_client.node.node_id)
    economic = new_validator_client.economic
    node = economic.env.get_consensus_node_by_index(0)
    consensus_list = economic.env.consensus_node_id_list()
    consensus_list.append(new_validator_client.node.node_id)
    consensus_list.pop(0)
    for i in range(4):
        node.debug.setValidatorList(consensus_list)
        print(get_pledge_list(node.ppos.getValidatorList))
        economic.wait_consensus_blocknum(node)

        wait_slashing_list = node.debug.getWaitSlashingNodeList()
        print(wait_slashing_list)
    new_validator_client.node.start(False)
    wait_slashing_list = node.debug.getWaitSlashingNodeList()
    print(wait_slashing_list)


def test_ZB_NP_22(new_validator_client):
    new_validator_client.node.stop()
    print(new_validator_client.node.node_id)
    economic = new_validator_client.economic
    node = economic.env.get_consensus_node_by_index(0)
    consensus_list = economic.env.consensus_node_id_list()
    consensus_list.append(new_validator_client.node.node_id)
    consensus_list.pop(0)
    for i in range(4):
        node.debug.setValidatorList(consensus_list)
        print(get_pledge_list(node.ppos.getValidatorList))
        economic.wait_consensus_blocknum(node)

        wait_slashing_list = node.debug.getWaitSlashingNodeList()
        print(wait_slashing_list)
    new_validator_client.node.start(False)
    wait_slashing_list = node.debug.getWaitSlashingNodeList()
    print(wait_slashing_list)


def test_ZB_NP_23(new_validator_client):
    new_validator_client.node.stop()
    print(new_validator_client.node.node_id)
    economic = new_validator_client.economic
    node = economic.env.get_consensus_node_by_index(0)
    consensus_list = economic.env.consensus_node_id_list()
    consensus_list.append(new_validator_client.node.node_id)
    consensus_list.pop(0)
    for i in range(4):
        node.debug.setValidatorList(consensus_list)
        print(get_pledge_list(node.ppos.getValidatorList))
        economic.wait_consensus_blocknum(node)

        wait_slashing_list = node.debug.getWaitSlashingNodeList()
        print(wait_slashing_list)
    new_validator_client.node.start(False)
    wait_slashing_list = node.debug.getWaitSlashingNodeList()
    print(wait_slashing_list)


def test_ZB_NP_24(new_validator_client):
    new_validator_client.node.stop()
    print(new_validator_client.node.node_id)
    economic = new_validator_client.economic
    node = economic.env.get_consensus_node_by_index(0)
    consensus_list = economic.env.consensus_node_id_list()
    consensus_list.append(new_validator_client.node.node_id)
    consensus_list.pop(0)
    for i in range(4):
        node.debug.setValidatorList(consensus_list)
        print(get_pledge_list(node.ppos.getValidatorList))
        economic.wait_consensus_blocknum(node)

        wait_slashing_list = node.debug.getWaitSlashingNodeList()
        print(wait_slashing_list)
    new_validator_client.node.start(False)
    wait_slashing_list = node.debug.getWaitSlashingNodeList()
    print(wait_slashing_list)


def test_ZB_NP_25(new_validator_client):
    new_validator_client.node.stop()
    print(new_validator_client.node.node_id)
    economic = new_validator_client.economic
    node = economic.env.get_consensus_node_by_index(0)
    consensus_list = economic.env.consensus_node_id_list()
    consensus_list.append(new_validator_client.node.node_id)
    consensus_list.pop(0)
    for i in range(4):
        node.debug.setValidatorList(consensus_list)
        print(get_pledge_list(node.ppos.getValidatorList))
        economic.wait_consensus_blocknum(node)

        wait_slashing_list = node.debug.getWaitSlashingNodeList()
        print(wait_slashing_list)
    new_validator_client.node.start(False)
    wait_slashing_list = node.debug.getWaitSlashingNodeList()
    print(wait_slashing_list)


def test_ZB_NP_26(new_validator_client):
    new_validator_client.node.stop()
    print(new_validator_client.node.node_id)
    economic = new_validator_client.economic
    node = economic.env.get_consensus_node_by_index(0)
    consensus_list = economic.env.consensus_node_id_list()
    consensus_list.append(new_validator_client.node.node_id)
    consensus_list.pop(0)
    for i in range(4):
        node.debug.setValidatorList(consensus_list)
        print(get_pledge_list(node.ppos.getValidatorList))
        economic.wait_consensus_blocknum(node)

        wait_slashing_list = node.debug.getWaitSlashingNodeList()
        print(wait_slashing_list)
    new_validator_client.node.start(False)
    wait_slashing_list = node.debug.getWaitSlashingNodeList()
    print(wait_slashing_list)


def test_ZB_NP_27(new_validator_client):
    new_validator_client.node.stop()
    print(new_validator_client.node.node_id)
    economic = new_validator_client.economic
    node = economic.env.get_consensus_node_by_index(0)
    consensus_list = economic.env.consensus_node_id_list()
    consensus_list.append(new_validator_client.node.node_id)
    consensus_list.pop(0)
    for i in range(4):
        node.debug.setValidatorList(consensus_list)
        print(get_pledge_list(node.ppos.getValidatorList))
        economic.wait_consensus_blocknum(node)

        wait_slashing_list = node.debug.getWaitSlashingNodeList()
        print(wait_slashing_list)
    new_validator_client.node.start(False)
    wait_slashing_list = node.debug.getWaitSlashingNodeList()
    print(wait_slashing_list)


def test_ZB_NP_28(new_validator_client):
    new_validator_client.node.stop()
    print(new_validator_client.node.node_id)
    economic = new_validator_client.economic
    node = economic.env.get_consensus_node_by_index(0)
    consensus_list = economic.env.consensus_node_id_list()
    consensus_list.append(new_validator_client.node.node_id)
    consensus_list.pop(0)
    for i in range(4):
        node.debug.setValidatorList(consensus_list)
        print(get_pledge_list(node.ppos.getValidatorList))
        economic.wait_consensus_blocknum(node)

        wait_slashing_list = node.debug.getWaitSlashingNodeList()
        print(wait_slashing_list)
    new_validator_client.node.start(False)
    wait_slashing_list = node.debug.getWaitSlashingNodeList()
    print(wait_slashing_list)


def test_ZB_NP_29(new_validator_client):
    new_validator_client.node.stop()
    print(new_validator_client.node.node_id)
    economic = new_validator_client.economic
    node = economic.env.get_consensus_node_by_index(0)
    consensus_list = economic.env.consensus_node_id_list()
    consensus_list.append(new_validator_client.node.node_id)
    consensus_list.pop(0)
    for i in range(4):
        node.debug.setValidatorList(consensus_list)
        print(get_pledge_list(node.ppos.getValidatorList))
        economic.wait_consensus_blocknum(node)

        wait_slashing_list = node.debug.getWaitSlashingNodeList()
        print(wait_slashing_list)
    new_validator_client.node.start(False)
    wait_slashing_list = node.debug.getWaitSlashingNodeList()
    print(wait_slashing_list)


def test_ZB_NP_30(new_validator_client):
    new_validator_client.node.stop()
    print(new_validator_client.node.node_id)
    economic = new_validator_client.economic
    node = economic.env.get_consensus_node_by_index(0)
    consensus_list = economic.env.consensus_node_id_list()
    consensus_list.append(new_validator_client.node.node_id)
    consensus_list.pop(0)
    for i in range(4):
        node.debug.setValidatorList(consensus_list)
        print(get_pledge_list(node.ppos.getValidatorList))
        economic.wait_consensus_blocknum(node)

        wait_slashing_list = node.debug.getWaitSlashingNodeList()
        print(wait_slashing_list)
    new_validator_client.node.start(False)
    wait_slashing_list = node.debug.getWaitSlashingNodeList()
    print(wait_slashing_list)
