# -*- coding: utf-8 -*-

from tests.lib.utils import *
import pytest
from tests.lib.config import EconomicConfig
import allure


@allure.title("Query delegate parameter validation")
@pytest.mark.P1
@pytest.mark.compatibility
def test_DI_001_009(client_new_node):
    """
    001:Query delegate parameter validation
    009：The money entrusted is equal to the low threshold entrusted
    """
    address, pri_key = client_new_node.economic.account.generate_account(client_new_node.node.web3,
                                                                         10 ** 18 * 10000000)
    client_new_node.staking.create_staking(0, address, address)
    address1, _ = client_new_node.economic.account.generate_account(client_new_node.node.web3,
                                                                    10 ** 18 * 10000000)
    result = client_new_node.delegate.delegate(0, address1)
    assert_code(result, 0)
    msg = client_new_node.ppos.getCandidateInfo(client_new_node.node.node_id)
    staking_blocknum = msg["Ret"]["StakingBlockNum"]
    msg = client_new_node.ppos.getDelegateInfo(staking_blocknum, address1, client_new_node.node.node_id)
    log.info(msg)
    assert client_new_node.node.web3.toChecksumAddress(msg["Ret"]["Addr"]) == address1
    assert msg["Ret"]["NodeId"] == client_new_node.node.node_id
    assert msg["Ret"]["ReleasedHes"] == client_new_node.economic.delegate_limit


@allure.title("Delegate to different people")
@pytest.mark.P1
def test_DI_002_003_004(clients_new_node):
    """
    002:Delegate to candidate
    003:Delegate to verifier
    004:Delegate to consensus verifier
    """
    client1 = clients_new_node[0]
    client2 = clients_new_node[1]

    staking_amount = client1.economic.create_staking_limit
    address, pri_key = client1.economic.account.generate_account(client1.node.web3, 10 ** 18 * 10000000)
    client1.staking.create_staking(0, address, address, amount=staking_amount)

    address, pri_key = client2.economic.account.generate_account(client2.node.web3, 10 ** 18 * 10000000)
    client2.staking.create_staking(0, address, address, amount=staking_amount * 2)

    client2.economic.wait_settlement_blocknum(client2.node)

    nodeid_list = get_pledge_list(client2.ppos.getVerifierList)
    log.info("The billing cycle validates the list of people{}".format(nodeid_list))
    assert client1.node.node_id not in nodeid_list

    address1, _ = client1.economic.account.generate_account(client1.node.web3, 10 ** 18 * 10000000)
    log.info("The candidate delegate")
    result = client1.delegate.delegate(0, address1)
    assert_code(result, 0)

    assert client2.node.node_id in nodeid_list
    address2, _ = client2.economic.account.generate_account(client2.node.web3,
                                                            10 ** 18 * 10000000)
    log.info("The verifier delegates")
    result = client2.delegate.delegate(0, address2)
    assert_code(result, 0)

    client2.economic.wait_consensus_blocknum(client2.node)
    nodeid_list = get_pledge_list(client2.ppos.getValidatorList)
    log.info("Consensus validator list:{}".format(nodeid_list))
    assert client2.node.node_id in nodeid_list
    address3, _ = client2.economic.account.generate_account(client2.node.web3,
                                                            10 ** 18 * 10000000)
    log.info("Consensus verifier delegates")
    result = client2.delegate.delegate(0, address3)
    assert_code(result, 0)


@allure.title("The amount entrusted by the client is less than the threshold")
@pytest.mark.P3
def test_DI_005(client_consensus):
    """
    :param client_consensus_obj:
    :return:
    """
    address, _ = client_consensus.economic.account.generate_account(client_consensus.node.web3,
                                                                    10 ** 18 * 10000000)

    result = client_consensus.delegate.delegate(0, address)
    log.info(result)
    assert_code(result, 301107)


@allure.title("The amount entrusted by the client is less than the threshold")
@pytest.mark.P1
def test_DI_006(client_new_node):
    """
    :param client_new_node_obj:
    :return:
    """
    address, _ = client_new_node.economic.account.generate_account(client_new_node.node.web3,
                                                                   10 ** 18 * 10000000)

    client_new_node.staking.create_staking(0, address, address)
    address1, _ = client_new_node.economic.account.generate_account(client_new_node.node.web3,
                                                                    10 ** 18 * 10000000)
    delegate_limit = client_new_node.economic.delegate_limit
    result = client_new_node.delegate.delegate(0, address1, amount=delegate_limit - 1)
    log.info(result)
    assert_code(result, 301105)


@allure.title("gas Insufficient entrustment")
@pytest.mark.P1
def test_DI_007(client_new_node):
    """
    :param client_new_node_obj:
    :return:
    """
    address, _ = client_new_node.economic.account.generate_account(client_new_node.node.web3,
                                                                   10 ** 18 * 10000000)

    client_new_node.staking.create_staking(0, address, address)
    address1, _ = client_new_node.economic.account.generate_account(client_new_node.node.web3,
                                                                    10 ** 18 * 10000000)

    fig = {"gas": 1}
    status = 0
    try:
        result = client_new_node.delegate.delegate(0, address1, transaction_cfg=fig)
        log.info(result)
    except BaseException:
        status = 1
    assert status == 1


@allure.title("not sufficient funds")
@pytest.mark.P1
def test_DI_008(client_new_node):
    """
    :param client_new_node_obj:
    :return:
    """
    address, _ = client_new_node.economic.account.generate_account(client_new_node.node.web3,
                                                                   10 ** 18 * 10000000)

    client_new_node.staking.create_staking(0, address, address)
    address1, _ = client_new_node.economic.account.generate_account(client_new_node.node.web3,
                                                                    10)
    status = 0
    try:
        result = client_new_node.delegate.delegate(0, address1)
        log.info(result)
    except BaseException:
        status = 1
    assert status == 1


@allure.title("Delegate to a candidate who doesn't exist")
@pytest.mark.P3
def test_DI_010_020(client_new_node):
    """
    Delegate to a candidate who doesn't exist
    :param client_new_node_obj:
    :return:
    """
    illegal_nodeID = "7ee3276fd6b9c7864eb896310b5393324b6db785a2528c00cc28ca8c" \
                     "3f86fc229a86f138b1f1c8e3a942204c03faeb40e3b22ab11b8983c35dc025de42865990"
    address1, pri_key = client_new_node.economic.account.generate_account(client_new_node.node.web3,
                                                                          10 ** 18 * 10000000)
    result = client_new_node.delegate.delegate(0, address1, node_id=illegal_nodeID)
    log.info(result)
    assert_code(result, 301102)


@allure.title("Delegate to different people{status}")
@pytest.mark.P1
@pytest.mark.parametrize('status', [0, 1, 2, 3])
def test_DI_011_012_013_014(client_new_node, status):
    """
    0:A valid candidate whose commission is still in doubt
    1:The delegate is also a valid candidate at a lockup period
    2:A candidate whose mandate is voluntarily withdrawn but who is still in the freeze period
    3:A candidate whose mandate has been voluntarily withdrawn and whose freeze period has expired
    :param client_new_node_obj:
    :param status:
    :return:
    """

    address, _ = client_new_node.economic.account.generate_account(client_new_node.node.web3,
                                                                   10 ** 18 * 10000000)

    client_new_node.staking.create_staking(0, address, address)
    if status == 0:
        # A valid candidate whose commission is still in doubt
        address1, _ = client_new_node.economic.account.generate_account(client_new_node.node.web3,
                                                                        10 ** 18 * 10000000)
        result = client_new_node.delegate.delegate(0, address1)
        assert_code(result, 0)

    if status == 1:
        # The delegate is also a valid candidate at a lockup period
        address1, _ = client_new_node.economic.account.generate_account(client_new_node.node.web3,
                                                                        10 ** 18 * 10000000)
        client_new_node.economic.wait_settlement_blocknum(client_new_node.node)
        result = client_new_node.delegate.delegate(0, address1)
        assert_code(result, 0)

    if status == 2:
        address1, _ = client_new_node.economic.account.generate_account(client_new_node.node.web3,
                                                                        10 ** 18 * 10000000)
        client_new_node.economic.wait_settlement_blocknum(client_new_node.node)
        result = client_new_node.staking.withdrew_staking(address)
        assert_code(result, 0)
        result = client_new_node.delegate.delegate(0, address1)
        assert_code(result, 301103)

    if status == 3:
        address1, _ = client_new_node.economic.account.generate_account(client_new_node.node.web3,
                                                                        10 ** 18 * 10000000)
        client_new_node.economic.wait_settlement_blocknum(client_new_node.node)
        result = client_new_node.staking.withdrew_staking(address)
        assert_code(result, 0)
        client_new_node.economic.wait_settlement_blocknum(client_new_node.node, number=2)
        result = client_new_node.delegate.delegate(0, address1)
        log.info(result)
        assert_code(result, 301102)


@allure.title("Delegate to candidates whose penalties have lapsed (freeze period and after freeze period)")
@pytest.mark.P1
def test_DI_015_016(client_new_node, client_consensus):
    """
    :param client_new_node_obj:
    :param client_consensus_obj:
    :return:
    """
    client = client_new_node
    node = client.node
    other_node = client_consensus.node
    economic = client.economic
    address, _ = economic.account.generate_account(client_new_node.node.web3,
                                                         10 ** 18 * 10000000)
    address_delegate, _ = economic.account.generate_account(client_new_node.node.web3,
                                                         10 ** 18 * 10000000)
    value = economic.create_staking_limit * 2
    result = client.staking.create_staking(0, address, address, amount=value)
    assert_code(result, 0)
    economic.wait_consensus_blocknum(other_node, number=4)
    validator_list = get_pledge_list(other_node.ppos.getValidatorList)
    assert node.node_id in validator_list
    candidate_info = other_node.ppos.getCandidateInfo(node.node_id)
    log.info(candidate_info)
    log.info("Close one node")
    node.stop()
    for i in range(4):
        economic.wait_consensus_blocknum(other_node, number=i)
        candidate_info = other_node.ppos.getCandidateInfo(node.node_id)
        log.info(candidate_info)
        if candidate_info["Ret"]["Released"] < value:
            break
        log.info("Node exceptions are not penalized")
    log.info("Restart the node")
    client_new_node.node.start()
    result = client.delegate.delegate(0, address_delegate)
    log.info(result)
    assert_code(result, 301103)
    log.info("Next settlement period")
    client_new_node.economic.wait_settlement_blocknum(node,number=2)
    result = client.delegate.delegate(0, address_delegate)
    assert_code(result, 301102)


@allure.title("Use the pledge account as the entrustment")
@pytest.mark.P1
def test_DI_017(client_new_node):
    """
    Use the pledge account as the entrustment
    :param client_new_node_obj:
    :return:
    """
    address, _ = client_new_node.economic.account.generate_account(client_new_node.node.web3,
                                                                   10 ** 18 * 10000000)

    result = client_new_node.staking.create_staking(0, address, address)
    assert_code(result, 0)
    result = client_new_node.delegate.delegate(0, address)
    log.info(result)
    assert_code(result, 301106)


@allure.title(
    "The verification section receives the delegate, exits, becomes the verification node, and receives the delegate")
@pytest.mark.P1
def test_DI_019(client_new_node):
    """
    :param client_new_node_obj:
    :return:
    """
    address, _ = client_new_node.economic.account.generate_account(client_new_node.node.web3,
                                                                   10 ** 18 * 10000000)
    result = client_new_node.staking.create_staking(0, address, address)
    assert_code(result, 0)

    address1, _ = client_new_node.economic.account.generate_account(client_new_node.node.web3,
                                                                    10 ** 18 * 10000000)
    result = client_new_node.delegate.delegate(0, address1)
    log.info(result)
    # Exit the pledge
    result = client_new_node.staking.withdrew_staking(address)
    assert_code(result, 0)
    # Repeat pledge
    result = client_new_node.staking.create_staking(0, address, address)
    assert_code(result, 0)
    result = client_new_node.delegate.delegate(0, address1)
    log.info(result)
    # Recheck wallet associations
    msg = client_new_node.ppos.getRelatedListByDelAddr(address1)
    log.info(msg)
    print(len(msg["Ret"]))
    assert len(msg["Ret"]) == 2
    for i in msg["Ret"]:
        assert client_new_node.node.web3.toChecksumAddress(i["Addr"]) == address1
        assert i["NodeId"] == client_new_node.node.node_id


@allure.title("The entrusted verifier is penalized to verify the entrusted principal")
@pytest.mark.P3
def test_DI_021(client_new_node, client_consensus):
    """
    :param client_new_node_obj:
    :return:
    """
    address, _ = client_new_node.economic.account.generate_account(client_new_node.node.web3,
                                                                   10 ** 18 * 10000000)
    address1, _ = client_new_node.economic.account.generate_account(client_new_node.node.web3,
                                                                    10 ** 18 * 10000000)
    value = client_new_node.economic.create_staking_limit * 2
    result = client_new_node.staking.create_staking(0, address, address, amount=value)
    assert_code(result, 0)
    result = client_new_node.delegate.delegate(0, address1)
    log.info(result)
    msg = client_new_node.ppos.getCandidateInfo(client_new_node.node.node_id)
    staking_blocknum = msg["Ret"]["StakingBlockNum"]
    log.info("Close one node")
    client_new_node.node.stop()
    node = client_consensus.node
    log.info("The next two periods")
    client_new_node.economic.wait_settlement_blocknum(node, number=2)
    log.info("Restart the node")
    client_new_node.node.start()
    msg = client_consensus.ppos.getDelegateInfo(staking_blocknum, address1, client_new_node.node.node_id)
    log.info(msg)
    assert msg["Ret"]["Released"] == client_new_node.economic.delegate_limit


@allure.title("Free amount in different periods when additional entrustment is made")
@pytest.mark.P2
@pytest.mark.parametrize('status', [0, 1, 2])
def test_DI_022_023_024(client_new_node, status):
    """
    022:There is only the free amount of hesitation period when additional entrusting
    023:Only the free amount of the lockup period exists when the delegate is added
    024:The amount of both hesitation period and lockup period exists when additional entrustment is made
    :param client_new_node_obj:
    :param status:
    :return:
    """
    client_new_node.economic.env.deploy_all()
    address, _ = client_new_node.economic.account.generate_account(client_new_node.node.web3,
                                                                   10 ** 18 * 10000000)

    client_new_node.staking.create_staking(0, address, address)
    address1, _ = client_new_node.economic.account.generate_account(client_new_node.node.web3,
                                                                    10 ** 18 * 10000000)
    result = client_new_node.delegate.delegate(0, address1)
    log.info(result)
    msg = client_new_node.ppos.getCandidateInfo(client_new_node.node.node_id)
    staking_blocknum = msg["Ret"]["StakingBlockNum"]

    if status == 0:
        result = client_new_node.delegate.delegate(0, address1)
        log.info(result)
        msg = client_new_node.ppos.getDelegateInfo(staking_blocknum, address1, client_new_node.node.node_id)
        log.info(msg)
        assert msg["Ret"]["ReleasedHes"] == client_new_node.economic.delegate_limit * 2

    if status == 1:
        client_new_node.economic.wait_settlement_blocknum(client_new_node.node)
        result = client_new_node.delegate.delegate(0, address1)
        log.info(result)
        msg = client_new_node.ppos.getDelegateInfo(staking_blocknum, address1, client_new_node.node.node_id)
        log.info(msg)
        assert msg["Ret"]["ReleasedHes"] == client_new_node.economic.delegate_limit
        assert msg["Ret"]["Released"] == client_new_node.economic.delegate_limit

    if status == 2:
        client_new_node.economic.wait_settlement_blocknum(client_new_node.node)
        result = client_new_node.delegate.delegate(0, address1)
        log.info(result)
        result = client_new_node.delegate.delegate(0, address1)
        log.info(result)
        msg = client_new_node.ppos.getDelegateInfo(staking_blocknum, address1, client_new_node.node.node_id)
        log.info(msg)
        assert msg["Ret"]["ReleasedHes"] == client_new_node.economic.delegate_limit * 2
        assert msg["Ret"]["Released"] == client_new_node.economic.delegate_limit


@allure.title("uncommitted")
@pytest.mark.P2
def test_DI_025(client_new_node):
    """
    :param client_new_node_obj:
    :return:
    """
    address_delegate, _ = client_new_node.economic.account.generate_account(client_new_node.node.web3,
                                                                            10 ** 18 * 10000000)

    result = client_new_node.ppos.getRelatedListByDelAddr(address_delegate)
    log.info(result)
    assert_code(result, 301203)


@allure.title("The entrusted candidate is valid")
@pytest.mark.P2
def test_DI_026(client_new_node):
    """
    :param client_new_node_obj:
    :return:
    """
    address, _ = client_new_node.economic.account.generate_account(client_new_node.node.web3,
                                                                   10 ** 18 * 10000000)
    address_delegate, _ = client_new_node.economic.account.generate_account(client_new_node.node.web3,
                                                                            10 ** 18 * 10000000)
    result = client_new_node.staking.create_staking(0, address, address)
    assert_code(result, 0)

    result = client_new_node.delegate.delegate(0, address_delegate)
    assert_code(result, 0)

    result = client_new_node.ppos.getRelatedListByDelAddr(address_delegate)
    log.info(result)
    assert result["Code"] == 0
    assert client_new_node.node.web3.toChecksumAddress(result["Ret"][0]["Addr"]) == address_delegate
    assert result["Ret"][0]["NodeId"] == client_new_node.node.node_id


@allure.title("The entrusted candidate does not exist")
@pytest.mark.P2
def test_DI_027(client_new_node):
    """
    The entrusted candidate does not exist
    :param client_new_node_obj:
    :return:
    """
    address, _ = client_new_node.economic.account.generate_account(client_new_node.node.web3,
                                                                   10 ** 18 * 10000000)
    address_delegate, _ = client_new_node.economic.account.generate_account(client_new_node.node.web3,
                                                                            10 ** 18 * 10000000)
    illegal_nodeID = "7ee3276fd6b9c7864eb896310b5393324b6db785a2528c00cc28ca8c" \
                     "3f86fc229a86f138b1f1c8e3a942204c03faeb40e3b22ab11b8983c35dc025de42865990"

    result = client_new_node.delegate.delegate(0, address_delegate, node_id=illegal_nodeID)
    log.info(result)
    result = client_new_node.ppos.getRelatedListByDelAddr(address_delegate)
    log.info(result)
    assert_code(result, 301203)


@allure.title("The entrusted candidate is invalid")
@pytest.mark.P2
def test_DI_028(client_new_node):
    """
    The entrusted candidate is invalid
    """
    address, _ = client_new_node.economic.account.generate_account(client_new_node.node.web3,
                                                                   10 ** 18 * 10000000)
    address_delegate, _ = client_new_node.economic.account.generate_account(client_new_node.node.web3,
                                                                            10 ** 18 * 10000000)
    result = client_new_node.staking.create_staking(0, address, address)
    assert_code(result, 0)

    result = client_new_node.delegate.delegate(0, address_delegate)
    assert_code(result, 0)

    # Exit the pledge
    result = client_new_node.staking.withdrew_staking(address)
    assert_code(result, 0)
    result = client_new_node.ppos.getRelatedListByDelAddr(address_delegate)
    assert result["Code"] == 0
    assert client_new_node.node.web3.toChecksumAddress(result["Ret"][0]["Addr"]) == address_delegate
    assert result["Ret"][0]["NodeId"] == client_new_node.node.node_id


@allure.title("Delegate information in the hesitation period, lock period")
@pytest.mark.P2
def test_DI_029_030(client_new_node):
    """
    029:Hesitation period inquiry entrustment details
    030:Lock periodic query information
    """
    address, _ = client_new_node.economic.account.generate_account(client_new_node.node.web3,
                                                                   10 ** 18 * 10000000)
    address_delegate, _ = client_new_node.economic.account.generate_account(client_new_node.node.web3,
                                                                            10 ** 18 * 10000000)

    client_new_node.staking.create_staking(0, address, address)
    result = client_new_node.delegate.delegate(0, address_delegate)
    assert_code(result, 0)

    # Hesitation period inquiry entrustment details
    result = client_new_node.ppos.getRelatedListByDelAddr(address_delegate)
    log.info(result)
    log.info("The next cycle")
    client_new_node.economic.wait_settlement_blocknum(client_new_node.node)
    result = client_new_node.ppos.getRelatedListByDelAddr(address_delegate)
    assert result["Code"] == 0
    assert client_new_node.node.web3.toChecksumAddress(result["Ret"][0]["Addr"]) == address_delegate
    assert result["Ret"][0]["NodeId"] == client_new_node.node.node_id


@allure.title("The delegate message no longer exists")
@pytest.mark.P2
def test_DI_031(client_new_node):
    """
    The delegate message no longer exists
    """
    address, _ = client_new_node.economic.account.generate_account(client_new_node.node.web3,
                                                                   10 ** 18 * 10000000)
    address_delegate, _ = client_new_node.economic.account.generate_account(client_new_node.node.web3,
                                                                            10 ** 18 * 10000000)

    client_new_node.staking.create_staking(0, address, address)
    result = client_new_node.delegate.delegate(0, address_delegate)
    assert_code(result, 0)

    msg = client_new_node.ppos.getCandidateInfo(client_new_node.node.node_id)
    staking_blocknum = msg["Ret"]["StakingBlockNum"]

    result = client_new_node.delegate.withdrew_delegate(staking_blocknum, address_delegate)
    assert_code(result, 0)
    result = client_new_node.ppos.getDelegateInfo(staking_blocknum, address_delegate,
                                                  client_new_node.node.node_id)
    log.info(result)
    assert_code(result, 301205)


@allure.title("The commission information is still in the hesitation period & The delegate information is still locked")
@pytest.mark.P2
def test_DI_032_033(client_new_node):
    """
    032:The commission information is still in the hesitation period
    033The delegate information is still locked
    """
    address, _ = client_new_node.economic.account.generate_account(client_new_node.node.web3,
                                                                   10 ** 18 * 10000000)
    address_delegate, _ = client_new_node.economic.account.generate_account(client_new_node.node.web3,
                                                                            10 ** 18 * 10000000)

    client_new_node.staking.create_staking(0, address, address)
    result = client_new_node.delegate.delegate(0, address_delegate)
    assert_code(result, 0)

    msg = client_new_node.ppos.getCandidateInfo(client_new_node.node.node_id)
    staking_blocknum = msg["Ret"]["StakingBlockNum"]

    # Hesitation period inquiry entrustment details
    result = client_new_node.ppos.getDelegateInfo(staking_blocknum, address_delegate,
                                                  client_new_node.node.node_id)
    log.info(result)
    assert client_new_node.node.web3.toChecksumAddress(result["Ret"]["Addr"]) == address_delegate
    assert result["Ret"]["NodeId"] == client_new_node.node.node_id
    log.info("The next cycle")
    client_new_node.economic.wait_consensus_blocknum(client_new_node.node)
    result = client_new_node.ppos.getDelegateInfo(staking_blocknum, address_delegate,
                                                  client_new_node.node.node_id)
    log.info(result)
    assert client_new_node.node.web3.toChecksumAddress(result["Ret"]["Addr"]) == address_delegate
    assert result["Ret"]["NodeId"] == client_new_node.node.node_id


@allure.title("The entrusted candidate has withdrawn of his own accord")
@pytest.mark.P2
def test_DI_034(client_new_node):
    """
    The entrusted candidate has withdrawn of his own accord
    """
    address, _ = client_new_node.economic.account.generate_account(client_new_node.node.web3,
                                                                   10 ** 18 * 10000000)
    address_delegate, _ = client_new_node.economic.account.generate_account(client_new_node.node.web3,
                                                                            10 ** 18 * 10000000)

    client_new_node.staking.create_staking(0, address, address)
    result = client_new_node.delegate.delegate(0, address_delegate)
    assert_code(result, 0)

    msg = client_new_node.ppos.getCandidateInfo(client_new_node.node.node_id)
    staking_blocknum = msg["Ret"]["StakingBlockNum"]

    # Exit the pledge
    result = client_new_node.staking.withdrew_staking(address)
    assert_code(result, 0)

    result = client_new_node.ppos.getDelegateInfo(staking_blocknum, address_delegate,
                                                  client_new_node.node.node_id)
    log.info(result)
    assert client_new_node.node.web3.toChecksumAddress(result["Ret"]["Addr"]) == address_delegate
    assert result["Ret"]["NodeId"] == client_new_node.node.node_id


@allure.title("Entrusted candidate (penalized in lockup period, penalized out completely)")
@pytest.mark.P2
def test_DI_035_036(clients_new_node, client_consensus):
    """
    The entrusted candidate is still penalized in the lockup period
    The entrusted candidate was penalized to withdraw completely

    """
    client = clients_new_node[0]
    node = client.node
    other_node = client_consensus.node
    economic = client.economic
    address, _ = economic.account.generate_account(node.web3,10 ** 18 * 10000000)

    address_delegate, _ = economic.account.generate_account(node.web3,10 ** 18 * 10000000)

    value = economic.create_staking_limit * 2
    result = client.staking.create_staking(0, address, address)
    assert_code(result, 0)
    result = client.delegate.delegate(0, address_delegate)
    assert_code(result, 0)
    ##The validation node becomes the out-block validation node
    economic.wait_consensus_blocknum(other_node, number=4)
    validator_list = get_pledge_list(other_node.ppos.getValidatorList)
    assert node.node_id in validator_list
    candidate_info = other_node.ppos.getCandidateInfo(node.node_id)
    log.info(candidate_info)
    staking_blocknum = candidate_info["Ret"]["StakingBlockNum"]

    log.info("Close one node")
    node.stop()
    for i in range(4):
        economic.wait_consensus_blocknum(other_node, number=i)
        candidate_info = other_node.ppos.getCandidateInfo(node.node_id)
        log.info(candidate_info)
        if candidate_info["Ret"]["Released"] < value:
            break

    result = other_node.ppos.getDelegateInfo(staking_blocknum, address_delegate,
                                       node.node_id)
    log.info(result)
    assert other_node.web3.toChecksumAddress(result["Ret"]["Addr"]) == address_delegate
    assert result["Ret"]["NodeId"] == node.node_id
    log.info("Restart the node")
    node.start()
    log.info("Next settlement period")
    economic.wait_settlement_blocknum(other_node,number=2)

    result = other_node.ppos.getDelegateInfo(staking_blocknum, address_delegate,
                                       node.node_id)
    log.info(result)
    assert other_node.web3.toChecksumAddress(result["Ret"]["Addr"]) == address_delegate
    assert result["Ret"]["NodeId"] == node.node_id


@allure.title("Query for delegate information in undo")
@pytest.mark.P2
def test_DI_038(client_new_node):
    """
    Query for delegate information in undo
    :param client_new_node_obj:
    :return:
    """
    address, _ = client_new_node.economic.account.generate_account(client_new_node.node.web3,
                                                                   10 ** 18 * 10000000)
    address_delegate, _ = client_new_node.economic.account.generate_account(client_new_node.node.web3,
                                                                            10 ** 18 * 10000000)

    client_new_node.staking.create_staking(0, address, address)
    result = client_new_node.delegate.delegate(0, address_delegate)
    assert_code(result, 0)

    msg = client_new_node.ppos.getCandidateInfo(client_new_node.node.node_id)
    staking_blocknum = msg["Ret"]["StakingBlockNum"]

    log.info("The next cycle")
    client_new_node.economic.wait_consensus_blocknum(client_new_node.node)

    # Exit the pledge
    result = client_new_node.staking.withdrew_staking(address)
    assert_code(result, 0)

    result = client_new_node.ppos.getDelegateInfo(staking_blocknum, address_delegate,
                                                  client_new_node.node.node_id)
    log.info(result)
    assert client_new_node.node.web3.toChecksumAddress(result["Ret"]["Addr"]) == address_delegate
    assert result["Ret"]["NodeId"] == client_new_node.node.node_id
