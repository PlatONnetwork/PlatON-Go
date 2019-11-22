# -*- coding: utf-8 -*-
from tests.lib.utils import *
import pytest
from tests.lib.config import EconomicConfig
import allure


@allure.title("Verify the validity of human parameters")
@pytest.mark.P1
def test_IV_001_002_010(global_test_env, client_consensus_obj):
    """
    001:Verify the validity of human parameters
    002:The built-in account is found with the verifier list query
    010:The initial number of nodes is consistent with the number of verifier consensus nodes set
    """
    global_test_env.deploy_all()
    node_info = client_consensus_obj.ppos.getValidatorList()
    log.info(node_info)
    nodeid_list = []
    for node in node_info.get("Ret"):
        nodeid_list.append(node.get("NodeId"))
        StakingAddress = node.get("StakingAddress")
        log.info(StakingAddress)
        assert client_consensus_obj.node.web3.toChecksumAddress(StakingAddress) == \
               client_consensus_obj.economic.cfg.DEVELOPER_FOUNDATAION_ADDRESS
    log.info(nodeid_list)
    consensus_node_list = global_test_env.consensus_node_list
    nodeid_list_ = [node.node_id for node in consensus_node_list]
    log.info(nodeid_list_)
    assert len(nodeid_list_) == len(consensus_node_list)
    for nodeid in nodeid_list_:
        assert nodeid in nodeid_list


@allure.title("Verify the validity of human parameters")
@pytest.mark.P1
def test_IV_003(client_consensus_obj):
    StakingAddress = EconomicConfig.DEVELOPER_FOUNDATAION_ADDRESS
    result = client_consensus_obj.staking.create_staking(0, StakingAddress, StakingAddress)
    log.info("Staking result:{}".format(result))
    assert_code(result, 301101)


@allure.title("The initial verifier accepts the delegate")
@pytest.mark.P1
def test_IV_004(client_consensus_obj):
    address, _ = client_consensus_obj.economic.account.generate_account(client_consensus_obj.node.web3,
                                                                        10 ** 18 * 10000000)
    result = client_consensus_obj.delegate.delegate(0, address)
    log.info(result)
    assert_code(result, 301107)


@allure.title("The initial verifier holds an additional pledge")
@pytest.mark.P1
def test_IV_005(client_consensus_obj):
    StakingAddress = EconomicConfig.DEVELOPER_FOUNDATAION_ADDRESS
    result = client_consensus_obj.staking.increase_staking(0, StakingAddress)
    assert_code(result, 0)


@allure.title("Initial verifier exits")
@pytest.mark.P1
def test_IV_006_007_008(client_consensus_obj):
    """
    006:Initial verifier exits
    007:The original verifier exits and re-pledges the pledge
    008:After the initial verifier quits, re-pledge and accept the entrustment
    """
    StakingAddress = client_consensus_obj.economic.cfg.DEVELOPER_FOUNDATAION_ADDRESS
    result = client_consensus_obj.staking.withdrew_staking(StakingAddress)
    log.info(result)
    result = client_consensus_obj.ppos.getCandidateInfo(client_consensus_obj.node.node_id)
    log.info(result)
    log.info("Let's go to the next three cycles")
    client_consensus_obj.economic.wait_settlement_blocknum(client_consensus_obj.node, number=2)
    msg = client_consensus_obj.ppos.getCandidateInfo(client_consensus_obj.node.node_id)
    log.info(msg)
    assert msg["Code"] == 301204, "预期验证人已退出"
    result = client_consensus_obj.staking.create_staking(0, StakingAddress, StakingAddress)
    assert_code(result, 0)
    address, _ = client_consensus_obj.economic.account.generate_account(client_consensus_obj.node.web3,
                                                                        10 ** 18 * 10000000)
    result = client_consensus_obj.delegate.delegate(0, address)
    log.info(result)
    assert_code(result, 0)


@allure.title("Modify the initial verifier information")
@pytest.mark.P3
def test_IV_009(client_consensus_obj):
    address1, _ = client_consensus_obj.economic.account.generate_account(client_consensus_obj.node.web3,
                                                                         10 ** 18 * 10000000)
    StakingAddress = client_consensus_obj.economic.cfg.DEVELOPER_FOUNDATAION_ADDRESS
    result = client_consensus_obj.staking.edit_candidate(StakingAddress, address1)
    log.info(result)
    assert_code(result, 0)


@allure.title("Normal pledge、Repeat pledge")
@pytest.mark.P1
@pytest.mark.compatibility
def test_IV_014_015_019_024(client_new_node_obj):
    """
    014：Normal pledge
    015：Repeat pledg
    019：The amount pledged by free account reaches the threshold of pledge
    024：Use the correct version signature
    """
    address, _ = client_new_node_obj.economic.account.generate_account(client_new_node_obj.node.web3,
                                                                       10 ** 18 * 10000000)
    log.info("Generate address:{}".format(address))
    result = client_new_node_obj.staking.create_staking(0, address, address)
    log.info(result)
    assert_code(result, 0)
    result = client_new_node_obj.staking.create_staking(0, address, address)
    log.info(result)
    assert_code(result, 301101)


@allure.title("Node ID pledge not added to the chain")
@pytest.mark.P3
def test_IV_016(client_new_node_obj):
    """
    :param client_new_node_obj:
    :return:
    """
    illegal_nodeID = "7ee3276fd6b9c7864eb896310b5393324b6db785a2528c00cc28ca8c" \
                     "3f86fc229a86f138b1f1c8e3a942204c03faeb40e3b22ab11b8983c35dc025de42865990"
    address, _ = client_new_node_obj.economic.account.generate_account(client_new_node_obj.node.web3,
                                                                       10 ** 18 * 10000000)
    result = client_new_node_obj.staking.create_staking(0, address, address, node_id=illegal_nodeID)
    log.info(result)
    assert_code(result, 301003)


@allure.title("The beneficiary address is the excitation pool address")
@pytest.mark.P3
def test_IV_017(client_new_node_obj):
    """
    :param client_new_node_obj:
    :return:
    """
    INCENTPEPOOL_ADDRESS = EconomicConfig.INCENTIVEPOOL_ADDRESS
    address, _ = client_new_node_obj.economic.account.generate_account(client_new_node_obj.node.web3,
                                                                       10 ** 18 * 10000000)
    result = client_new_node_obj.staking.create_staking(0, INCENTPEPOOL_ADDRESS, address)
    assert_code(result, 0)


@allure.title("Beneficiary address is the foundation address")
@pytest.mark.P3
def test_IV_018(client_new_node_obj):
    """
    :param client_new_node_obj:
    :param get_generate_account:
    :return:
    """
    FOUNDATION_ADDRESS = client_new_node_obj.economic.cfg.FOUNDATION_ADDRESS
    address, _ = client_new_node_obj.economic.account.generate_account(client_new_node_obj.node.web3,
                                                                       10 ** 18 * 10000000)
    result = client_new_node_obj.staking.create_staking(0, FOUNDATION_ADDRESS, address)
    assert_code(result, 0)


@allure.title("The amount pledged by free account is less than the threshold of pledge, and gas is insufficient")
@pytest.mark.P2
def test_IV_020_21(client_new_node_obj):
    """
    020:The amount pledged by free account is less than the threshold of pledge
    021:gas is insufficient
    """
    address, _ = client_new_node_obj.economic.account.generate_account(client_new_node_obj.node.web3,
                                                                       10 ** 18 * 10000000)
    amount = client_new_node_obj.economic.create_staking_limit
    result = client_new_node_obj.staking.create_staking(0, address, address, amount=amount - 1)
    log.info(result)
    assert_code(result, 301100)
    cfg = {"gas": 1}
    status = 0
    try:
        result = client_new_node_obj.staking.create_staking(0, address, address, transaction_cfg=cfg)
        log.info(result)
    except BaseException:
        status = 1
    assert status == 1


@allure.title("Incorrect version signature used")
@pytest.mark.P3
def test_IV_025(client_new_node_obj, client_consensus_obj):
    """
    :param client_new_node_obj:
    :param client_consensus_obj:
    :return:
    """
    address, _ = client_new_node_obj.economic.account.generate_account(client_new_node_obj.node.web3,
                                                                       10 ** 18 * 10000000)
    program_version_sign = client_consensus_obj.node.program_version_sign
    result = client_new_node_obj.staking.create_staking(0, address, address, program_version_sign=program_version_sign)
    log.info(result)
    assert_code(result, 301003)


@allure.title("BlsPublicKey is too long")
@pytest.mark.P2
def test_IV_026_01(client_new_node_obj):
    address, _ = client_new_node_obj.economic.account.generate_account(client_new_node_obj.node.web3,
                                                                       10 ** 18 * 10000000)
    blspubkey = client_new_node_obj.node.blspubkey + "00000000"
    log.info(blspubkey)
    status = 0
    try:
        result = client_new_node_obj.staking.create_staking(0, address, address, bls_pubkey=blspubkey)
        log.info(result)
    except BaseException:
        status = 1
    assert status == 1


@allure.title("BlsPublicKey is too short")
@pytest.mark.P2
def test_IV_026_02(client_new_node_obj):
    address, _ = client_new_node_obj.economic.account.generate_account(client_new_node_obj.node.web3,
                                                                       10 ** 18 * 10000000)
    blspubkey = client_new_node_obj.node.blspubkey[0:10]
    log.info(blspubkey)
    status = 0
    try:
        result = client_new_node_obj.staking.create_staking(0, address, address, bls_pubkey=blspubkey)
        log.info(result)
    except BaseException:
        status = 1
    assert status == 1


@allure.title("Incorrect version signature used")
@pytest.mark.P2
def test_IV_026_03(client_new_node_obj):
    address, _ = client_new_node_obj.economic.account.generate_account(client_new_node_obj.node.web3,
                                                                       10 ** 18 * 10000000)
    program_version = 0000
    result = client_new_node_obj.staking.create_staking(0, address, address, program_version=program_version)
    assert_code(result, 301003)


@allure.title("Field length overflow")
@pytest.mark.P2
def test_IV_027(client_new_node_obj):
    external_id = "11111111111111111111111111111111111111111111111111111111111111111111111111111111111"
    node_name = "1111111111111111111111111111111111111111111111111111111111111111111111111111111111111"
    website = "1111111111111111111111111111111111111111111111111111111111111111111111111111111111111111 "
    details = "1111111111111111111111111111111111111111111111111111111111111111111111111111111111111111 "
    program_version = client_new_node_obj.node.program_version
    program_version_sign = client_new_node_obj.node.program_version_sign
    bls_pubkey = client_new_node_obj.node.blspubkey
    bls_proof = client_new_node_obj.node.schnorr_NIZK_prove
    amount = client_new_node_obj.economic.create_staking_limit
    address, pri_key = client_new_node_obj.economic.account.generate_account(client_new_node_obj.node.web3,
                                                                             10 ** 18 * 10000000)

    result = client_new_node_obj.ppos.createStaking(0, address, client_new_node_obj.node.node_id,
                                                    external_id, node_name, website, details, amount,
                                                    program_version, program_version_sign, bls_pubkey, bls_proof,
                                                    pri_key)
    assert_code(result, 301002)


@allure.title("Pledge has been made before and the candidate has been invalidated (penalized)")
@pytest.mark.P2
def test_IV_028(client_new_node_obj, client_consensus_obj):
    address, pri_key = client_new_node_obj.economic.account.generate_account(client_new_node_obj.node.web3,
                                                                             10 ** 18 * 10000000)
    value = client_new_node_obj.economic.create_staking_limit * 2
    result = client_new_node_obj.staking.create_staking(0, address, address, amount=value)
    assert_code(result, 0)
    log.info("Close one node")
    client_new_node_obj.node.stop()
    node = client_consensus_obj.node
    log.info("The next two periods")
    client_new_node_obj.economic.wait_settlement_blocknum(node, number=2)
    log.info("Restart the node")
    client_new_node_obj.node.start()
    result = client_new_node_obj.staking.edit_candidate(address, address)
    log.info(result)
    assert_code(result, 301103)
    log.info("Next settlement period")
    client_new_node_obj.economic.wait_settlement_blocknum(client_new_node_obj.node)
    time.sleep(20)
    result = client_new_node_obj.staking.create_staking(0, address, address, amount=value)
    assert_code(result, 0)


@allure.title("Pledge has been made before, and the candidate is invalid (voluntarily withdraw)")
@pytest.mark.P1
def test_IV_029(client_new_node_obj):
    """
    :param client_new_node_obj:
    :param get_generate_account:
    :return:
    """
    address, _ = client_new_node_obj.economic.account.generate_account(client_new_node_obj.node.web3,
                                                                       10 ** 18 * 10000000)
    result = client_new_node_obj.staking.create_staking(0, address, address)
    assert_code(result, 0)
    result = client_new_node_obj.staking.withdrew_staking(address)
    log.info(result)
    assert_code(result, 0)
    result = client_new_node_obj.staking.create_staking(0, address, address)
    assert_code(result, 0)


@allure.title("Lockup pledge")
@pytest.mark.P1
def test_IV_030(client_new_node_obj):
    """
    :param client_new_node_obj:
    :param get_generate_account:
    :return:
    """
    address, _ = client_new_node_obj.economic.account.generate_account(client_new_node_obj.node.web3,
                                                                       10 ** 18 * 10000000)
    result = client_new_node_obj.staking.create_staking(0, address, address)
    assert_code(result, 0)
    log.info("Into the next grandchild")
    client_new_node_obj.economic.wait_settlement_blocknum(client_new_node_obj.node)
    result = client_new_node_obj.staking.withdrew_staking(address)
    assert_code(result, 0)
    result = client_new_node_obj.staking.create_staking(0, address, address)
    log.info(result)
    assert_code(result, 301101)


@allure.title("Use a new wallet as pledge")
@pytest.mark.P2
def test_IV_031(client_new_node_obj):
    """
    :param client_new_node_obj:
    :param get_generate_account:
    :return:
    """
    address, _ = client_new_node_obj.economic.account.generate_account(client_new_node_obj.node.web3,
                                                                       10 ** 18 * 10000000)
    result = client_new_node_obj.staking.create_staking(0, address, address)
    assert_code(result, 0)

    address1, _ = client_new_node_obj.economic.account.generate_account(client_new_node_obj.node.web3,
                                                                        10 ** 18 * 10000000)
    result = client_new_node_obj.staking.create_staking(0, address1, address1)
    log.info(result)
    assert_code(result, 301101)
