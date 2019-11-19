# -*- coding: utf-8 -*-
from tests.lib.utils import *
import pytest
from tests.lib.config import EconomicConfig
from common.key import mock_duplicate_sign


@pytest.mark.P1
def test_CS_CL_001(global_test_env, client_new_node_obj):
    """
    The longer the tenure, the easier it is to replace
    :param client_new_node_obj:
    :return:
    """
    global_test_env.deploy_all()
    address, _ = client_new_node_obj.economic.account.generate_account(client_new_node_obj.node.web3,
                                                                       10 ** 18 * 10000000)
    pledge_amount = client_new_node_obj.economic.create_staking_limit * 2
    result = client_new_node_obj.staking.create_staking(0, address, address, amount=pledge_amount)
    assert_code(result, 0)

    # Next settlement period
    client_new_node_obj.economic.wait_settlement_blocknum(client_new_node_obj.node)
    # The next consensus cycle
    client_new_node_obj.economic.wait_consensus_blocknum(client_new_node_obj.node)
    validatorlist1 = get_pledge_list(client_new_node_obj.ppos.getValidatorList)
    log.info("validatorlist:{}".format(validatorlist1))
    node_1 = get_validator_term(client_new_node_obj.node)
    log.info("Maximum tenure node:{}".format(node_1))
    # The next consensus cycle
    client_new_node_obj.economic.wait_consensus_blocknum(client_new_node_obj.node)

    validatorlist2 = get_pledge_list(client_new_node_obj.ppos.getValidatorList)
    log.info("validatorlist:{}".format(validatorlist2))
    node_2 = get_validator_term(client_new_node_obj.node)
    log.info("Maximum tenure node:{}".format(node_2))
    assert node_1 not in validatorlist2

    # The next consensus cycle
    client_new_node_obj.economic.wait_consensus_blocknum(client_new_node_obj.node)
    validatorlist3 = get_pledge_list(client_new_node_obj.ppos.getValidatorList)
    log.info("validatorlist:{}".format(validatorlist3))
    node_3 = get_validator_term(client_new_node_obj.node)
    log.info("Maximum tenure node:{}".format(node_3))
    assert node_2 not in validatorlist3


@pytest.mark.P1
def test_CS_CL_002(global_test_env, client_new_node_obj):
    """
    The higher the consensus verifier list index is replaced
    :param client_new_node_obj:
    :return:
    """
    global_test_env.deploy_all()

    node = get_max_staking_tx_index(client_new_node_obj.node)
    log.info("The node with the largest trade index:{}".format(node))

    address, _ = client_new_node_obj.economic.account.generate_account(client_new_node_obj.node.web3,
                                                                       10 ** 18 * 10000000)
    pledge_amount = client_new_node_obj.economic.create_staking_limit * 2
    result = client_new_node_obj.staking.create_staking(0, address, address, amount=pledge_amount)
    assert_code(result, 0)

    # Next settlement period
    client_new_node_obj.economic.wait_settlement_blocknum(client_new_node_obj.node)
    # The next consensus cycle
    client_new_node_obj.economic.wait_consensus_blocknum(client_new_node_obj.node)
    validatorlist = get_pledge_list(client_new_node_obj.ppos.getValidatorList)
    log.info("validatorlist:{}".format(validatorlist))
    assert node not in validatorlist


@pytest.mark.P1
def test_CS_CL_003(global_test_env, client_con_list_obj, client_noc_list_obj):
    """
    The higher the consensus verifier list block, the higher it is replaced
    :param client_new_node_obj:
    :return:
    """
    global_test_env.deploy_all()
    validatorlist = get_pledge_list(client_con_list_obj[0].ppos.getValidatorList)
    log.info("validatorlist:{}".format(validatorlist))

    log.info("The next consensus cycle")
    client_con_list_obj[0].economic.wait_consensus_blocknum(client_con_list_obj[0].node, number=1)
    blocknumber = client_con_list_obj[0].node.eth.blockNumber
    log.info("To report the double sign")
    report_information1 = mock_duplicate_sign(1, client_con_list_obj[0].node.nodekey,
                                              client_con_list_obj[0].node.blsprikey,
                                              blocknumber)
    log.info("Report information: {}".format(report_information1))

    report_information2 = mock_duplicate_sign(1, client_con_list_obj[1].node.nodekey,
                                              client_con_list_obj[1].node.blsprikey,
                                              blocknumber)
    log.info("Report information: {}".format(report_information2))

    report_information3 = mock_duplicate_sign(1, client_con_list_obj[2].node.nodekey,
                                              client_con_list_obj[2].node.blsprikey,
                                              blocknumber)
    log.info("Report information: {}".format(report_information3))

    address_1, _ = client_noc_list_obj[0].economic.account.generate_account(client_noc_list_obj[0].node.web3,
                                                                            10 ** 18 * 10000000)
    address_2, _ = client_noc_list_obj[0].economic.account.generate_account(client_noc_list_obj[0].node.web3,
                                                                            10 ** 18 * 10000000)
    address_3, _ = client_noc_list_obj[0].economic.account.generate_account(client_noc_list_obj[0].node.web3,
                                                                            10 ** 18 * 10000000)
    result = client_con_list_obj[0].duplicatesign.reportDuplicateSign(1, report_information1, address_1)
    log.info(result)
    result = client_con_list_obj[1].duplicatesign.reportDuplicateSign(1, report_information2, address_2)
    log.info(result)
    result = client_con_list_obj[2].duplicatesign.reportDuplicateSign(1, report_information3, address_3)
    log.info(result)
    log.info("The next  periods")
    client_noc_list_obj[2].economic.wait_settlement_blocknum(client_noc_list_obj[2].node)
    validatorlist = get_pledge_list(client_con_list_obj[0].ppos.getValidatorList)
    log.info("After being reported validatorlist:{}".format(validatorlist))

    staking_address1, _ = client_noc_list_obj[0].economic.account.generate_account(client_con_list_obj[0].node.web3,
                                                                                   10 ** 18 * 10000000)
    staking_address2, _ = client_noc_list_obj[0].economic.account.generate_account(client_con_list_obj[0].node.web3,
                                                                                   10 ** 18 * 10000000)
    staking_address3, _ = client_noc_list_obj[0].economic.account.generate_account(client_con_list_obj[0].node.web3,
                                                                                   10 ** 18 * 10000000)
    staking_address4, _ = client_noc_list_obj[0].economic.account.generate_account(client_con_list_obj[0].node.web3,
                                                                                   10 ** 18 * 10000000)
    pledge_amount = client_con_list_obj[0].economic.create_staking_limit * 2
    log.info("New pledge 4 verifiers")
    result = client_noc_list_obj[0].staking.create_staking(0, staking_address1, staking_address1, amount=pledge_amount)
    assert_code(result, 0)
    result = client_noc_list_obj[1].staking.create_staking(0, staking_address2, staking_address2, amount=pledge_amount)
    assert_code(result, 0)
    result = client_noc_list_obj[2].staking.create_staking(0, staking_address3, staking_address3, amount=pledge_amount)
    assert_code(result, 0)
    result = client_noc_list_obj[3].staking.create_staking(0, staking_address4, staking_address4, amount=pledge_amount)
    assert_code(result, 0)
    log.info("Next settlement period")
    client_noc_list_obj[3].economic.wait_settlement_blocknum(client_noc_list_obj[3].node)
    log.info("The next consensus cycle")
    client_noc_list_obj[3].economic.wait_consensus_blocknum(client_noc_list_obj[3].node)
    verifierList = get_pledge_list(client_noc_list_obj[3].ppos.getVerifierList)
    log.info("verifierList:{}".format(verifierList))
    validatorlist1 = get_pledge_list(client_noc_list_obj[3].ppos.getValidatorList)
    log.info("validatorlist:{}".format(validatorlist1))
    assert client_noc_list_obj[0].node.node_id in validatorlist1

    client_noc_list_obj[3].economic.wait_consensus_blocknum(client_noc_list_obj[3].node)
    verifierList = get_pledge_list(client_noc_list_obj[3].ppos.getVerifierList)
    log.info("verifierList:{}".format(verifierList))
    validatorlist2 = get_pledge_list(client_noc_list_obj[3].ppos.getValidatorList)
    log.info("validatorlist:{}".format(validatorlist2))
    assert client_noc_list_obj[0].node.node_id in validatorlist2

    client_noc_list_obj[3].economic.wait_consensus_blocknum(client_noc_list_obj[3].node)
    verifierList = get_pledge_list(client_noc_list_obj[3].ppos.getVerifierList)
    log.info("verifierList:{}".format(verifierList))
    validatorlist3 = get_pledge_list(client_noc_list_obj[3].ppos.getValidatorList)
    log.info("validatorlist:{}".format(validatorlist3))
    assert client_noc_list_obj[0].node.node_id in validatorlist3


@pytest.mark.P1
def test_CS_CL_004(global_test_env, client_consensus_obj, client_new_node_obj):
    """
    The lower the total Shares, the easier it is to be replaced
    :param client_consensus_obj:
    :param client_new_node_obj:
    :return:
    """
    global_test_env.deploy_all()
    log.info(client_consensus_obj.node.node_id)
    log.info(client_new_node_obj.node.node_id)
    log.info(client_consensus_obj.node.url)
    StakingAddress = EconomicConfig.DEVELOPER_FOUNDATAION_ADDRESS
    value = client_consensus_obj.node.web3.toWei(1000000, "ether")
    result = client_consensus_obj.staking.increase_staking(0, StakingAddress, amount=value)
    assert_code(result, 0)
    msg = client_consensus_obj.ppos.getCandidateInfo(client_consensus_obj.node.node_id)
    log.info(msg)
    address, _ = client_new_node_obj.economic.account.generate_account(client_new_node_obj.node.web3,
                                                                       10 ** 18 * 10000000)
    value = client_new_node_obj.economic.create_staking_limit * 2
    result = client_new_node_obj.staking.create_staking(0, address, address, amount=value)
    assert_code(result, 0)
    log.info(client_new_node_obj.node.node_id)
    log.info("Next settlement period")
    client_consensus_obj.economic.wait_settlement_blocknum(client_consensus_obj.node)
    log.info("The next consensus cycle")
    client_consensus_obj.economic.wait_consensus_blocknum(client_consensus_obj.node)
    validatorlist1 = get_pledge_list(client_consensus_obj.ppos.getValidatorList)
    log.info("validatorlist:{}".format(validatorlist1))

    msg = client_consensus_obj.ppos.getValidatorList()
    log.info("Consensus validates the person's situation{}".format(msg))
    assert client_consensus_obj.node.node_id in validatorlist1

    client_new_node_obj.economic.wait_consensus_blocknum(client_new_node_obj.node)
    validatorlist2 = get_pledge_list(client_consensus_obj.ppos.getValidatorList)
    log.info("validatorlist:{}".format(validatorlist2))
    msg = client_consensus_obj.ppos.getValidatorList()
    log.info("Consensus validates the person's situation{}".format(msg))
    assert client_consensus_obj.node.node_id in validatorlist2


@pytest.mark.P1
def test_CS_CL_005_006_008(global_test_env, client_noc_list_obj):
    """
    :param client_consensus_obj:
    :param client_new_node_obj:
    :return:
    """
    global_test_env.deploy_all()
    client_noc_list_obj[0].economic.env.deploy_all()
    address1, _ = client_noc_list_obj[0].economic.account.generate_account(client_noc_list_obj[0].node.web3,
                                                                           10 ** 18 * 10000000)
    address2, _ = client_noc_list_obj[0].economic.account.generate_account(client_noc_list_obj[0].node.web3,
                                                                           10 ** 18 * 10000000)
    value = client_noc_list_obj[0].economic.create_staking_limit * 2
    result = client_noc_list_obj[0].staking.create_staking(0, address1, address1, amount=value)
    assert_code(result, 0)
    result = client_noc_list_obj[1].staking.create_staking(0, address2, address2,
                                                           amount=value + 1300000000000000000000000)
    assert_code(result, 0)

    # Next settlement period
    client_noc_list_obj[1].economic.wait_settlement_blocknum(client_noc_list_obj[1].node)

    verifierlist = get_pledge_list(client_noc_list_obj[1].ppos.getVerifierList)
    log.info("verifierlist:{}".format(verifierlist))

    msg = client_noc_list_obj[0].ppos.getCandidateInfo(client_noc_list_obj[0].node.node_id)
    log.info(msg)
    msg = client_noc_list_obj[1].ppos.getCandidateInfo(client_noc_list_obj[1].node.node_id)
    log.info(msg)
    assert client_noc_list_obj[1].node.node_id in verifierlist
    assert client_noc_list_obj[1].node.node_id == verifierlist[0]

    address3, _ = client_noc_list_obj[0].economic.account.generate_account(client_noc_list_obj[0].node.web3,
                                                                           10 ** 18 * 10000000)

    result = client_noc_list_obj[0].delegate.delegate(0, address3, amount=700000000000000000000000)
    log.info(result)
    result = client_noc_list_obj[0].staking.increase_staking(0, address1, amount=610000000000000000000000)
    log.info(result)

    # Next settlement period
    client_noc_list_obj[1].economic.wait_settlement_blocknum(client_noc_list_obj[1].node)
    msg = client_noc_list_obj[0].ppos.getCandidateInfo(client_noc_list_obj[0].node.node_id)
    log.info(msg)
    msg = client_noc_list_obj[1].ppos.getCandidateInfo(client_noc_list_obj[1].node.node_id)
    log.info(msg)
    verifierlist = get_pledge_list(client_noc_list_obj[1].ppos.getVerifierList)
    log.info("verifierlist:{}".format(verifierlist))
    assert verifierlist[0] == client_noc_list_obj[0].node.node_id
    assert verifierlist[1] == client_noc_list_obj[1].node.node_id


@pytest.mark.P1
def test_CS_CL_007(global_test_env, client_noc_list_obj):
    """

    :param client_noc_list_obj:
    :return:
    """
    global_test_env.deploy_all()
    address1, _ = client_noc_list_obj[0].economic.account.generate_account(client_noc_list_obj[0].node.web3,
                                                                           10 ** 18 * 10000000)
    address2, _ = client_noc_list_obj[0].economic.account.generate_account(client_noc_list_obj[0].node.web3,
                                                                           10 ** 18 * 10000000)
    value = client_noc_list_obj[0].economic.create_staking_limit * 2
    result = client_noc_list_obj[0].staking.create_staking(0, address1, address1, amount=value)
    assert_code(result, 0)
    result = client_noc_list_obj[1].staking.create_staking(0, address2, address2, amount=value)
    assert_code(result, 0)
    # Next settlement period
    client_noc_list_obj[0].economic.wait_settlement_blocknum(client_noc_list_obj[1].node)

    verifierlist = get_pledge_list(client_noc_list_obj[1].ppos.getVerifierList)
    log.info("verifierlist:{}".format(verifierlist))
    log.info("node:{}".format(client_noc_list_obj[0].node.node_id))
    assert verifierlist[0] == client_noc_list_obj[0].node.node_id


@pytest.mark.P1
def test_CS_CL_010_030(global_test_env, client_new_node_obj):
    """
    :param global_test_env:
    :param client_new_node_obj:
    :return:
    """
    global_test_env.deploy_all()
    address, _ = client_new_node_obj.economic.account.generate_account(client_new_node_obj.node.web3,
                                                                       10 ** 18 * 10000000)
    value = client_new_node_obj.economic.create_staking_limit * 2
    result = client_new_node_obj.staking.create_staking(0, address, address, amount=value)
    assert_code(result, 0)
    log.info("The next  periods")
    client_new_node_obj.economic.wait_settlement_blocknum(client_new_node_obj.node)
    log.info("The next consensus cycle")
    client_new_node_obj.economic.wait_consensus_blocknum(client_new_node_obj.node)
    number = client_new_node_obj.node.eth.blockNumber
    log.info("To report the double sign")
    report_information = mock_duplicate_sign(1, client_new_node_obj.node.nodekey, client_new_node_obj.node.blsprikey,
                                             number)

    log.info("Report information: {}".format(report_information))
    address, _ = client_new_node_obj.economic.account.generate_account(client_new_node_obj.node.web3,
                                                                       10 ** 18 * 10000000)
    result = client_new_node_obj.duplicatesign.reportDuplicateSign(1, report_information, address)
    log.info(result)

    log.info("The next  periods")
    client_new_node_obj.economic.wait_settlement_blocknum(client_new_node_obj.node)
    verifierlist = get_pledge_list(client_new_node_obj.ppos.getVerifierList)
    log.info("verifierlist:{}".format(verifierlist))
    assert client_new_node_obj.node.node_id not in verifierlist


@pytest.mark.P1
def test_CS_CL_012_032(global_test_env, client_new_node_obj):
    """
    :param client_new_node_obj:
    :return:
    """
    global_test_env.deploy_all()
    address, _ = client_new_node_obj.economic.account.generate_account(client_new_node_obj.node.web3,
                                                                       10 ** 18 * 10000000)
    value = client_new_node_obj.economic.create_staking_limit * 2
    result = client_new_node_obj.staking.create_staking(0, address, address, amount=value)
    assert_code(result, 0)
    log.info("The next  periods")
    client_new_node_obj.economic.wait_settlement_blocknum(client_new_node_obj.node)
    verifierlist = get_pledge_list(client_new_node_obj.ppos.getVerifierList)
    log.info("verifierlist:{}".format(verifierlist))

    assert client_new_node_obj.node.node_id in verifierlist

    result = client_new_node_obj.staking.withdrew_staking(address)
    assert_code(result, 0)

    log.info("The next  periods")
    client_new_node_obj.economic.wait_settlement_blocknum(client_new_node_obj.node, number=1)
    verifierlist = get_pledge_list(client_new_node_obj.ppos.getVerifierList)
    log.info("verifierlist:{}".format(verifierlist))

    assert client_new_node_obj.node.node_id not in verifierlist

    log.info("The next consensus cycle")
    client_new_node_obj.economic.wait_consensus_blocknum(client_new_node_obj.node)

    validatorlist = get_pledge_list(client_new_node_obj.ppos.getValidatorList)
    log.info("validatorlist:{}".format(validatorlist))

    assert client_new_node_obj.node.node_id not in validatorlist


@pytest.mark.P1
@pytest.mark.compatibility
def test_CS_CL_013_031(global_test_env, client_new_node_obj, client_consensus_obj):
    """

    :param client_new_node_obj:
    :param client_consensus_obj:
    :return:
    """
    global_test_env.deploy_all()
    address, _ = client_new_node_obj.economic.account.generate_account(client_new_node_obj.node.web3,
                                                                       10 ** 18 * 10000000)
    value = client_new_node_obj.economic.create_staking_limit * 2
    result = client_new_node_obj.staking.create_staking(0, address, address, amount=value)
    assert_code(result, 0)
    # Next settlement period
    client_new_node_obj.economic.wait_settlement_blocknum(client_new_node_obj.node)
    verifierlist = get_pledge_list(client_new_node_obj.ppos.getVerifierList)
    log.info("verifierlist:{}".format(verifierlist))
    assert client_new_node_obj.node.node_id in verifierlist

    log.info("Close one node")
    client_new_node_obj.node.stop()
    node = client_consensus_obj.node

    log.info("The next  periods")
    client_new_node_obj.economic.wait_settlement_blocknum(node)
    verifierlist = get_pledge_list(client_consensus_obj.ppos.getVerifierList)
    log.info("verifierlist:{}".format(verifierlist))
    assert client_new_node_obj.node.node_id not in verifierlist


@pytest.mark.P2
@pytest.mark.parametrize('status', [0, 1, 2])
def test_CS_CL_014_015_016_029(status, global_test_env, client_con_list_obj, client_noc_list_obj):
    """
    :param status:
    :param global_test_env:
    :param client_con_list_obj:
    :param client_noc_list_obj:
    :return:
    """
    global_test_env.deploy_all()

    log.info("The next consensus cycle")
    client_con_list_obj[0].economic.wait_consensus_blocknum(client_con_list_obj[0].node, number=1)
    blocknumber = client_con_list_obj[0].node.eth.blockNumber

    log.info("To report the double sign")
    report_information1 = mock_duplicate_sign(1, client_con_list_obj[0].node.nodekey,
                                              client_con_list_obj[0].node.blsprikey,
                                              blocknumber)
    log.info("Report information: {}".format(report_information1))

    report_information2 = mock_duplicate_sign(1, client_con_list_obj[1].node.nodekey,
                                              client_con_list_obj[1].node.blsprikey,
                                              blocknumber)
    log.info("Report information: {}".format(report_information2))

    report_information3 = mock_duplicate_sign(1, client_con_list_obj[2].node.nodekey,
                                              client_con_list_obj[2].node.blsprikey,
                                              blocknumber)
    log.info("Report information: {}".format(report_information3))

    address_1, _ = client_noc_list_obj[0].economic.account.generate_account(client_noc_list_obj[0].node.web3,
                                                                            10 ** 18 * 10000000)
    address_2, _ = client_noc_list_obj[0].economic.account.generate_account(client_noc_list_obj[0].node.web3,
                                                                            10 ** 18 * 10000000)
    address_3, _ = client_noc_list_obj[0].economic.account.generate_account(client_noc_list_obj[0].node.web3,
                                                                            10 ** 18 * 10000000)
    result = client_con_list_obj[0].duplicatesign.reportDuplicateSign(1, report_information1, address_1)
    log.info(result)
    result = client_con_list_obj[1].duplicatesign.reportDuplicateSign(1, report_information2, address_2)
    log.info(result)
    result = client_con_list_obj[2].duplicatesign.reportDuplicateSign(1, report_information3, address_3)
    log.info(result)
    log.info("The next  periods")
    client_noc_list_obj[1].economic.wait_settlement_blocknum(client_noc_list_obj[1].node)
    validatorlist = get_pledge_list(client_con_list_obj[0].ppos.getValidatorList)
    log.info("After being reported validatorlist:{}".format(validatorlist))

    if status == 0:
        address1, _ = client_noc_list_obj[0].economic.account.generate_account(client_noc_list_obj[0].node.web3,
                                                                               10 ** 18 * 10000000)
        value = client_noc_list_obj[0].economic.create_staking_limit * 2
        result = client_noc_list_obj[0].staking.create_staking(0, address1, address1, amount=value)
        assert_code(result, 0)

        address2, _ = client_noc_list_obj[0].economic.account.generate_account(client_noc_list_obj[0].node.web3,
                                                                               10 ** 18 * 10000000)
        value = client_noc_list_obj[0].economic.create_staking_limit * 2
        result = client_noc_list_obj[1].staking.create_staking(0, address2, address2, amount=value)
        assert_code(result, 0)

        log.info("The next  periods")
        client_noc_list_obj[1].economic.wait_settlement_blocknum(client_noc_list_obj[1].node)

        log.info("The next consensus cycle")
        client_noc_list_obj[1].economic.wait_consensus_blocknum(client_noc_list_obj[1].node)

        validatorlist = get_pledge_list(client_noc_list_obj[1].ppos.getValidatorList)
        log.info("validatorlist:{}".format(validatorlist))
        log.info("node1:{}".format(client_noc_list_obj[0].node.node_id))
        log.info("node2:{}".format(client_noc_list_obj[1].node.node_id))
        log.info("node3:{}".format(client_con_list_obj[3].node.node_id))
        assert client_noc_list_obj[0].node.node_id in validatorlist
        assert client_noc_list_obj[1].node.node_id in validatorlist
        assert client_con_list_obj[3].node.node_id in validatorlist

    if status == 1:
        address1, _ = client_noc_list_obj[0].economic.account.generate_account(client_noc_list_obj[0].node.web3,
                                                                               10 ** 18 * 10000000)
        value = client_noc_list_obj[0].economic.create_staking_limit * 2
        result = client_noc_list_obj[0].staking.create_staking(0, address1, address1, amount=value)
        assert_code(result, 0)

        address2, _ = client_noc_list_obj[0].economic.account.generate_account(client_noc_list_obj[0].node.web3,
                                                                               10 ** 18 * 10000000)
        value = client_noc_list_obj[0].economic.create_staking_limit * 2
        result = client_noc_list_obj[1].staking.create_staking(0, address2, address2, amount=value)
        assert_code(result, 0)

        address3, _ = client_noc_list_obj[0].economic.account.generate_account(client_noc_list_obj[0].node.web3,
                                                                               10 ** 18 * 10000000)
        value = client_noc_list_obj[0].economic.create_staking_limit * 2
        result = client_noc_list_obj[2].staking.create_staking(0, address3, address3, amount=value)
        assert_code(result, 0)
        log.info("The next  periods")
        client_noc_list_obj[1].economic.wait_settlement_blocknum(client_noc_list_obj[1].node)

        log.info("The next consensus cycle")
        client_noc_list_obj[1].economic.wait_consensus_blocknum(client_noc_list_obj[1].node)

        validatorlist = get_pledge_list(client_noc_list_obj[1].ppos.getValidatorList)
        log.info("validatorlist:{}".format(validatorlist))
        log.info("node1:{}".format(client_noc_list_obj[0].node.node_id))
        log.info("node2:{}".format(client_noc_list_obj[1].node.node_id))
        log.info("node3:{}".format(client_noc_list_obj[2].node.node_id))
        log.info("node4:{}".format(client_con_list_obj[3].node.node_id))
        assert client_noc_list_obj[0].node.node_id in validatorlist
        assert client_noc_list_obj[1].node.node_id in validatorlist
        assert client_noc_list_obj[2].node.node_id in validatorlist
        assert client_con_list_obj[3].node.node_id in validatorlist

    if status == 2:
        address1, _ = client_noc_list_obj[0].economic.account.generate_account(client_noc_list_obj[0].node.web3,
                                                                               10 ** 18 * 10000000)
        value = client_noc_list_obj[0].economic.create_staking_limit * 2
        result = client_noc_list_obj[0].staking.create_staking(0, address1, address1, amount=value)
        assert_code(result, 0)

        address2, _ = client_noc_list_obj[0].economic.account.generate_account(client_noc_list_obj[0].node.web3,
                                                                               10 ** 18 * 10000000)
        value = client_noc_list_obj[0].economic.create_staking_limit * 2
        result = client_noc_list_obj[1].staking.create_staking(0, address2, address2, amount=value)
        assert_code(result, 0)

        address3, _ = client_noc_list_obj[0].economic.account.generate_account(client_noc_list_obj[0].node.web3,
                                                                               10 ** 18 * 10000000)
        value = client_noc_list_obj[0].economic.create_staking_limit * 2
        result = client_noc_list_obj[2].staking.create_staking(0, address3, address3, amount=value)
        assert_code(result, 0)

        address4, _ = client_noc_list_obj[0].economic.account.generate_account(client_noc_list_obj[0].node.web3,
                                                                               10 ** 18 * 10000000)
        value = client_noc_list_obj[0].economic.create_staking_limit * 2
        result = client_noc_list_obj[3].staking.create_staking(0, address4, address4, amount=value)
        assert_code(result, 0)

        log.info("The next  periods")
        client_noc_list_obj[1].economic.wait_settlement_blocknum(client_noc_list_obj[1].node)

        log.info("The next consensus cycle")
        client_noc_list_obj[1].economic.wait_consensus_blocknum(client_noc_list_obj[1].node)

        validatorlist = get_pledge_list(client_noc_list_obj[1].ppos.getValidatorList)
        log.info("validatorlist:{}".format(validatorlist))
        log.info("node:{}".format(client_con_list_obj[3].node.node_id))
        assert client_con_list_obj[3].node.node_id in validatorlist


@pytest.mark.P2
@pytest.mark.parametrize('status', [0, 1])
def test_CS_CL_017_018_019(status, global_test_env, client_con_list_obj, client_noc_list_obj):
    """
    :param status:
    :param global_test_env:
    :param client_con_list_obj:
    :param client_noc_list_obj:
    :return:
    """
    global_test_env.deploy_all()

    log.info("The next consensus cycle")
    client_con_list_obj[0].economic.wait_consensus_blocknum(client_con_list_obj[0].node, number=1)

    validatorlist = get_pledge_list(client_con_list_obj[0].ppos.getValidatorList)
    log.info("initial validatorlist:{}".format(validatorlist))
    blocknumber = client_con_list_obj[0].node.eth.blockNumber
    log.info("The thrill of being reported{}".format(blocknumber))

    log.info("To report the double sign")
    report_information1 = mock_duplicate_sign(1, client_con_list_obj[0].node.nodekey,
                                              client_con_list_obj[0].node.blsprikey,
                                              blocknumber)
    log.info("Report information: {}".format(report_information1))

    address, _ = client_noc_list_obj[0].economic.account.generate_account(client_noc_list_obj[0].node.web3,
                                                                          10 ** 18 * 10000000)
    result = client_con_list_obj[0].duplicatesign.reportDuplicateSign(1, report_information1, address)
    log.info(result)
    log.info("The next  periods")
    client_noc_list_obj[1].economic.wait_settlement_blocknum(client_noc_list_obj[1].node)
    validatorlist = get_pledge_list(client_con_list_obj[0].ppos.getValidatorList)
    log.info("After being reported validatorlist:{}".format(validatorlist))

    if status == 0:
        address1, _ = client_noc_list_obj[0].economic.account.generate_account(client_noc_list_obj[0].node.web3,
                                                                               10 ** 18 * 10000000)
        value = client_noc_list_obj[0].economic.create_staking_limit * 2
        result = client_noc_list_obj[0].staking.create_staking(0, address1, address1, amount=value)
        assert_code(result, 0)
        log.info("The next  periods")
        client_noc_list_obj[1].economic.wait_settlement_blocknum(client_noc_list_obj[1].node)
        log.info("The next consensus cycle")
        client_noc_list_obj[1].economic.wait_consensus_blocknum(client_noc_list_obj[1].node)

        validatorlist = get_pledge_list(client_noc_list_obj[1].ppos.getValidatorList)
        log.info("validatorlist:{}".format(validatorlist))
        assert client_con_list_obj[1].node.node_id in validatorlist
        assert client_con_list_obj[2].node.node_id in validatorlist
        assert client_con_list_obj[3].node.node_id in validatorlist
        assert client_noc_list_obj[0].node.node_id in validatorlist

    if status == 1:
        address1, _ = client_noc_list_obj[0].economic.account.generate_account(client_noc_list_obj[0].node.web3,
                                                                               10 ** 18 * 10000000)
        value = client_noc_list_obj[0].economic.create_staking_limit * 2
        result = client_noc_list_obj[0].staking.create_staking(0, address1, address1, amount=value)
        assert_code(result, 0)

        address2, _ = client_noc_list_obj[0].economic.account.generate_account(client_noc_list_obj[0].node.web3,
                                                                               10 ** 18 * 10000000)
        value = client_noc_list_obj[0].economic.create_staking_limit * 2
        result = client_noc_list_obj[1].staking.create_staking(0, address2, address2, amount=value)
        assert_code(result, 0)

        log.info("The next  periods")
        client_noc_list_obj[1].economic.wait_settlement_blocknum(client_noc_list_obj[1].node)
        log.info("The next consensus cycle")
        client_noc_list_obj[1].economic.wait_consensus_blocknum(client_noc_list_obj[1].node)

        validatorlist = get_pledge_list(client_noc_list_obj[1].ppos.getValidatorList)
        log.info("validatorlist:{}".format(validatorlist))
        assert client_con_list_obj[1].node.node_id in validatorlist
        assert client_con_list_obj[2].node.node_id in validatorlist
        assert client_con_list_obj[3].node.node_id in validatorlist


@pytest.mark.P2
def test_CS_CL_027_028(global_test_env, client_noc_list_obj):
    global_test_env.deploy_all()
    address1, _ = client_noc_list_obj[0].economic.account.generate_account(client_noc_list_obj[0].node.web3,
                                                                           10 ** 18 * 10000000)
    address2, _ = client_noc_list_obj[0].economic.account.generate_account(client_noc_list_obj[0].node.web3,
                                                                           10 ** 18 * 10000000)

    result = client_noc_list_obj[0].staking.create_staking(0, address1, address1,
                                                           amount=client_noc_list_obj[
                                                                      0].economic.create_staking_limit * 2)
    assert_code(result, 0)

    result = client_noc_list_obj[1].staking.create_staking(0, address2, address2,
                                                           amount=client_noc_list_obj[1].economic.create_staking_limit)
    assert_code(result, 0)

    log.info("Next settlement period")
    client_noc_list_obj[1].economic.wait_settlement_blocknum(client_noc_list_obj[1].node)
    msg = client_noc_list_obj[1].ppos.getVerifierList()
    log.info(msg)
    verifierlist = get_pledge_list(client_noc_list_obj[1].ppos.getVerifierList)
    log.info("verifierlist:{}".format(verifierlist))
    assert client_noc_list_obj[0].node.node_id in verifierlist
    assert client_noc_list_obj[1].node.node_id not in verifierlist


@pytest.mark.P2
def test_CS_CL_033(global_test_env, client_new_node_obj):
    global_test_env.deploy_all()
    address1, _ = client_new_node_obj.economic.account.generate_account(client_new_node_obj.node.web3,
                                                                        10 ** 18 * 10000000)

    value = client_new_node_obj.economic.create_staking_limit * 2
    result = client_new_node_obj.staking.create_staking(0, address1, address1, amount=value)
    assert_code(result, 0)

    # Next settlement period
    client_new_node_obj.economic.wait_settlement_blocknum(client_new_node_obj.node)
    # Next consensus period
    client_new_node_obj.economic.wait_consensus_blocknum(client_new_node_obj.node)

    verifierlist = get_pledge_list(client_new_node_obj.ppos.getVerifierList)
    log.info("verifierlist:{}".format(verifierlist))
    assert client_new_node_obj.node.node_id in verifierlist
