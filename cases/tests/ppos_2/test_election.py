# -*- coding: utf-8 -*-
from tests.lib.utils import *
import pytest
from tests.lib.config import EconomicConfig
from common.key import mock_duplicate_sign


@pytest.mark.P1
def test_CS_CL_001(clients_new_node):
    """
    The longer the tenure, the easier it is to replace
    :param client_new_node:
    :return:
    """
    client = clients_new_node[0]
    node_id_1 = client.node.node_id
    log.info(client.node.node_id)
    address, _ = client.economic.account.generate_account(client.node.web3,
                                                          10 ** 18 * 10000000)
    pledge_amount = client.economic.create_staking_limit * 2
    result = client.staking.create_staking(0, address, address, amount=pledge_amount)
    assert_code(result, 0)
    log.info("Next settlement period")
    client.economic.wait_settlement_blocknum(client.node)
    log.info("The next consensus cycle")
    client.economic.wait_consensus_blocknum(client.node)
    validatorlist1 = get_pledge_list(client.ppos.getValidatorList)
    log.info("validatorlist:{}".format(validatorlist1))

    log.info("替换的第一轮共识轮")
    log.info(client.ppos.getValidatorList())

    max_term_nodeid_1 = get_validator_term(client.node)
    log.info("Maximum tenure node list:{}".format(max_term_nodeid_1))
    assert node_id_1 not in max_term_nodeid_1
    assert node_id_1 in validatorlist1

    # The next consensus cycle
    log.info("The next consensus cycle")
    client.economic.wait_consensus_blocknum(client.node)

    validatorlist2 = get_pledge_list(client.ppos.getValidatorList)
    log.info("validatorlist:{}".format(validatorlist2))

    log.info("替换的第二轮共识轮")
    log.info(client.ppos.getValidatorList())

    max_term_nodeid_2 = get_validator_term(client.node)
    log.info("Maximum tenure node:{}".format(max_term_nodeid_2))
    assert node_id_1 not in max_term_nodeid_2
    assert node_id_1 in validatorlist2

    # The next consensus cycle
    log.info("The next consensus cycle")
    client.economic.wait_consensus_blocknum(client.node)

    validatorlist3 = get_pledge_list(client.ppos.getValidatorList)
    log.info("validatorlist:{}".format(validatorlist3))

    log.info("替换的第三轮共识轮")
    log.info(client.ppos.getValidatorList())
    max_term_nodeid_3 = get_validator_term(client.node)
    log.info("Maximum tenure node:{}".format(max_term_nodeid_3))
    assert node_id_1 not in max_term_nodeid_3
    assert node_id_1 in validatorlist3


@pytest.mark.P1
def test_CS_CL_002(clients_new_node):
    """
    The higher the consensus verifier list index is replaced
    :param client_new_node:
    :return:
    """
    client = clients_new_node[0]

    node = get_max_staking_tx_index(client.node)
    log.info("The node with the largest trade index:{}".format(node))

    address, _ = client.economic.account.generate_account(client.node.web3,
                                                          10 ** 18 * 10000000)
    pledge_amount = client.economic.create_staking_limit * 2
    result = client.staking.create_staking(0, address, address, amount=pledge_amount)
    assert_code(result, 0)
    # Next settlement period
    client.economic.wait_settlement_blocknum(client.node)
    # The next consensus cycle
    client.economic.wait_consensus_blocknum(client.node)
    validatorlist = get_pledge_list(client.ppos.getValidatorList)
    log.info("validatorlist:{}".format(validatorlist))
    assert node not in validatorlist


@pytest.mark.P1
def test_CS_CL_003(clients_new_node, clients_consensus):
    """
    The higher the consensus verifier list block, the higher it is replaced
    :param client_new_node:
    :return:
    """
    client_noconsensus1 = clients_new_node[0]
    client_noconsensus2 = clients_new_node[1]
    client_noconsensus3 = clients_new_node[2]
    client_noconsensus4 = clients_new_node[3]
    client_consensus1 = clients_consensus[0]
    client_consensus2 = clients_consensus[1]
    client_consensus3 = clients_consensus[2]
    validatorlist = get_pledge_list(client_consensus1.ppos.getValidatorList)
    log.info("validatorlist:{}".format(validatorlist))

    log.info("The next consensus cycle")
    client_consensus1.economic.wait_consensus_blocknum(client_consensus1.node, number=1)
    blocknumber = client_consensus1.node.eth.blockNumber
    log.info("To report the double sign")
    report_information1 = mock_duplicate_sign(1, client_consensus1.node.nodekey,
                                              client_consensus1.node.blsprikey,
                                              blocknumber)
    log.info("Report information: {}".format(report_information1))

    report_information2 = mock_duplicate_sign(1, client_consensus2.node.nodekey,
                                              client_consensus2.node.blsprikey,
                                              blocknumber)
    log.info("Report information: {}".format(report_information2))

    report_information3 = mock_duplicate_sign(1, client_consensus3.node.nodekey,
                                              client_consensus3.node.blsprikey,
                                              blocknumber)
    log.info("Report information: {}".format(report_information3))

    address_1, _ = client_noconsensus1.economic.account.generate_account(client_noconsensus1.node.web3,
                                                                         10 ** 18 * 10000000)
    address_2, _ = client_noconsensus1.economic.account.generate_account(client_noconsensus1.node.web3,
                                                                         10 ** 18 * 10000000)
    address_3, _ = client_noconsensus1.economic.account.generate_account(client_noconsensus1.node.web3,
                                                                         10 ** 18 * 10000000)
    result = client_consensus1.duplicatesign.reportDuplicateSign(1, report_information1, address_1)
    log.info(result)
    result = client_consensus2.duplicatesign.reportDuplicateSign(1, report_information2, address_2)
    log.info(result)
    result = client_consensus3.duplicatesign.reportDuplicateSign(1, report_information3, address_3)
    log.info(result)
    log.info("The next  periods")
    client_noconsensus3.economic.wait_settlement_blocknum(client_noconsensus3.node)
    validatorlist = get_pledge_list(client_consensus1.ppos.getValidatorList)
    log.info("After being reported validatorlist:{}".format(validatorlist))

    staking_address1, _ = client_noconsensus1.economic.account.generate_account(client_consensus1.node.web3,
                                                                                10 ** 18 * 10000000)
    staking_address2, _ = client_noconsensus1.economic.account.generate_account(client_consensus1.node.web3,
                                                                                10 ** 18 * 10000000)
    staking_address3, _ = client_noconsensus1.economic.account.generate_account(client_consensus1.node.web3,
                                                                                10 ** 18 * 10000000)
    staking_address4, _ = client_noconsensus1.economic.account.generate_account(client_consensus1.node.web3,
                                                                                10 ** 18 * 10000000)
    pledge_amount = client_consensus1.economic.create_staking_limit * 2
    log.info("New pledge 4 verifiers")
    result = client_noconsensus1.staking.create_staking(0, staking_address1, staking_address1, amount=pledge_amount)
    assert_code(result, 0)
    result = client_noconsensus2.staking.create_staking(0, staking_address2, staking_address2, amount=pledge_amount)
    assert_code(result, 0)
    result = client_noconsensus3.staking.create_staking(0, staking_address3, staking_address3, amount=pledge_amount)
    assert_code(result, 0)
    result = client_noconsensus4.staking.create_staking(0, staking_address4, staking_address4, amount=pledge_amount)
    assert_code(result, 0)
    log.info("Next settlement period")
    client_noconsensus4.economic.wait_settlement_blocknum(client_noconsensus4.node)
    log.info("The next consensus cycle")
    client_noconsensus4.economic.wait_consensus_blocknum(client_noconsensus4.node)
    verifierList = get_pledge_list(client_noconsensus4.ppos.getVerifierList)
    log.info("verifierList:{}".format(verifierList))
    validatorlist1 = get_pledge_list(client_noconsensus4.ppos.getValidatorList)
    log.info("validatorlist:{}".format(validatorlist1))
    assert client_noconsensus1.node.node_id in validatorlist1

    client_noconsensus4.economic.wait_consensus_blocknum(client_noconsensus4.node)
    verifierList = get_pledge_list(client_noconsensus4.ppos.getVerifierList)
    log.info("verifierList:{}".format(verifierList))
    validatorlist2 = get_pledge_list(client_noconsensus4.ppos.getValidatorList)
    log.info("validatorlist:{}".format(validatorlist2))
    assert client_noconsensus1.node.node_id in validatorlist2

    client_noconsensus4.economic.wait_consensus_blocknum(client_noconsensus4.node)
    verifierList = get_pledge_list(client_noconsensus4.ppos.getVerifierList)
    log.info("verifierList:{}".format(verifierList))
    validatorlist3 = get_pledge_list(client_noconsensus4.ppos.getValidatorList)
    log.info("validatorlist:{}".format(validatorlist3))
    assert client_noconsensus1.node.node_id in validatorlist3


@pytest.mark.P1
def test_CS_CL_004(clients_new_node, client_consensus):
    """
    The lower the total Shares, the easier it is to be replaced
    :param client_consensus_obj:
    :param client_new_node:
    :return:
    """
    client = clients_new_node[0]
    StakingAddress = EconomicConfig.DEVELOPER_FOUNDATAION_ADDRESS
    value = client_consensus.node.web3.toWei(1000000, "ether")
    result = client_consensus.staking.increase_staking(0, StakingAddress, amount=value)
    assert_code(result, 0)
    msg = client_consensus.ppos.getCandidateInfo(client_consensus.node.node_id)
    log.info(msg)
    address, _ = client.economic.account.generate_account(client.node.web3,
                                                          10 ** 18 * 10000000)
    value = client.economic.create_staking_limit * 2
    result = client.staking.create_staking(0, address, address, amount=value)
    assert_code(result, 0)
    log.info(client.node.node_id)
    log.info("Next settlement period")
    client_consensus.economic.wait_settlement_blocknum(client_consensus.node)
    log.info("The next consensus cycle")
    client_consensus.economic.wait_consensus_blocknum(client_consensus.node)
    validatorlist1 = get_pledge_list(client_consensus.ppos.getValidatorList)
    log.info("validatorlist:{}".format(validatorlist1))

    msg = client_consensus.ppos.getValidatorList()
    log.info("Consensus validates the person's situation{}".format(msg))
    assert client_consensus.node.node_id in validatorlist1

    client.economic.wait_consensus_blocknum(client.node)
    validatorlist2 = get_pledge_list(client_consensus.ppos.getValidatorList)
    log.info("validatorlist:{}".format(validatorlist2))
    msg = client_consensus.ppos.getValidatorList()
    log.info("Consensus validates the person's situation{}".format(msg))
    assert client_consensus.node.node_id in validatorlist2


@pytest.mark.P1
def test_CS_CL_005_006_008(clients_new_node):
    """
    :param client_consensus:
    :param client_new_node:
    :return:
    """
    client_noconsensus1 = clients_new_node[0]
    client_noconsensus2 = clients_new_node[1]
    client_noconsensus1.economic.env.deploy_all()
    address1, _ = client_noconsensus1.economic.account.generate_account(client_noconsensus1.node.web3,
                                                                        10 ** 18 * 10000000)
    address2, _ = client_noconsensus1.economic.account.generate_account(client_noconsensus1.node.web3,
                                                                        10 ** 18 * 10000000)
    value = client_noconsensus1.economic.create_staking_limit * 2
    result = client_noconsensus1.staking.create_staking(0, address1, address1, amount=value)
    assert_code(result, 0)
    result = client_noconsensus2.staking.create_staking(0, address2, address2,
                                                        amount=value + 1300000000000000000000000)
    assert_code(result, 0)

    # Next settlement period
    client_noconsensus2.economic.wait_settlement_blocknum(client_noconsensus2.node)

    verifierlist = get_pledge_list(client_noconsensus2.ppos.getVerifierList)
    log.info("verifierlist:{}".format(verifierlist))

    msg = client_noconsensus1.ppos.getCandidateInfo(client_noconsensus1.node.node_id)
    log.info(msg)
    msg = client_noconsensus2.ppos.getCandidateInfo(client_noconsensus2.node.node_id)
    log.info(msg)
    assert client_noconsensus2.node.node_id in verifierlist
    assert client_noconsensus2.node.node_id == verifierlist[0]

    address3, _ = client_noconsensus1.economic.account.generate_account(client_noconsensus1.node.web3,
                                                                        10 ** 18 * 10000000)

    result = client_noconsensus1.delegate.delegate(0, address3, amount=700000000000000000000000)
    log.info(result)
    result = client_noconsensus1.staking.increase_staking(0, address1, amount=610000000000000000000000)
    log.info(result)

    # Next settlement period
    client_noconsensus2.economic.wait_settlement_blocknum(client_noconsensus2.node)
    msg = client_noconsensus1.ppos.getCandidateInfo(client_noconsensus1.node.node_id)
    log.info(msg)
    msg = client_noconsensus2.ppos.getCandidateInfo(client_noconsensus2.node.node_id)
    log.info(msg)
    verifierlist = get_pledge_list(client_noconsensus2.ppos.getVerifierList)
    log.info("verifierlist:{}".format(verifierlist))
    assert verifierlist[0] == client_noconsensus1.node.node_id
    assert verifierlist[1] == client_noconsensus2.node.node_id


@pytest.mark.P1
def test_CS_CL_007(clients_new_node):
    """
    :param clients_noconsensus:
    :return:
    """
    client_noconsensus1 = clients_new_node[0]
    client_noconsensus2 = clients_new_node[1]
    address1, _ = client_noconsensus1.economic.account.generate_account(client_noconsensus1.node.web3,
                                                                        10 ** 18 * 10000000)
    address2, _ = client_noconsensus1.economic.account.generate_account(client_noconsensus1.node.web3,
                                                                        10 ** 18 * 10000000)
    value = client_noconsensus1.economic.create_staking_limit * 2
    result = client_noconsensus1.staking.create_staking(0, address1, address1, amount=value)
    assert_code(result, 0)
    result = client_noconsensus2.staking.create_staking(0, address2, address2, amount=value)
    assert_code(result, 0)
    # Next settlement period
    client_noconsensus1.economic.wait_settlement_blocknum(client_noconsensus2.node)

    verifierlist = get_pledge_list(client_noconsensus2.ppos.getVerifierList)
    log.info("verifierlist:{}".format(verifierlist))
    log.info("node:{}".format(client_noconsensus1.node.node_id))
    assert verifierlist[0] == client_noconsensus1.node.node_id


@pytest.mark.P1
def test_CS_CL_010_030(clients_new_node):
    """
    :param global_test_env:
    :param client_new_node:
    :return:
    """
    client = clients_new_node[0]
    address, _ = client.economic.account.generate_account(client.node.web3,
                                                          10 ** 18 * 10000000)
    value = client.economic.create_staking_limit * 2
    result = client.staking.create_staking(0, address, address, amount=value)
    assert_code(result, 0)
    log.info("The next  periods")
    client.economic.wait_settlement_blocknum(client.node)
    log.info("The next consensus cycle")
    client.economic.wait_consensus_blocknum(client.node)
    number = client.node.eth.blockNumber
    log.info("To report the double sign")
    report_information = mock_duplicate_sign(1, client.node.nodekey, client.node.blsprikey,
                                             number)

    log.info("Report information: {}".format(report_information))
    address, _ = client.economic.account.generate_account(client.node.web3,
                                                          10 ** 18 * 10000000)
    result = client.duplicatesign.reportDuplicateSign(1, report_information, address)
    log.info(result)

    log.info("The next  periods")
    client.economic.wait_settlement_blocknum(client.node)
    verifierlist = get_pledge_list(client.ppos.getVerifierList)
    log.info("verifierlist:{}".format(verifierlist))
    assert client.node.node_id not in verifierlist


@pytest.mark.P1
def test_CS_CL_012_032(clients_new_node):
    """
    :param client_new_node:
    :return:
    """
    client = clients_new_node[0]
    address, _ = client.economic.account.generate_account(client.node.web3,
                                                          10 ** 18 * 10000000)
    value = client.economic.create_staking_limit * 2
    result = client.staking.create_staking(0, address, address, amount=value)
    assert_code(result, 0)
    log.info("The next  periods")
    client.economic.wait_settlement_blocknum(client.node)
    verifierlist = get_pledge_list(client.ppos.getVerifierList)
    log.info("verifierlist:{}".format(verifierlist))

    assert client.node.node_id in verifierlist

    result = client.staking.withdrew_staking(address)
    assert_code(result, 0)

    log.info("The next  periods")
    client.economic.wait_settlement_blocknum(client.node, number=1)
    verifierlist = get_pledge_list(client.ppos.getVerifierList)
    log.info("verifierlist:{}".format(verifierlist))

    assert client.node.node_id not in verifierlist

    log.info("The next consensus cycle")
    client.economic.wait_consensus_blocknum(client.node)

    validatorlist = get_pledge_list(client.ppos.getValidatorList)
    log.info("validatorlist:{}".format(validatorlist))

    assert client.node.node_id not in validatorlist


@pytest.mark.P1
@pytest.mark.compatibility
def test_CS_CL_013_031(clients_new_node, client_consensus):
    """
    :param client_new_node:
    :param client_consensus_obj:
    :return:
    """
    client = clients_new_node[0]
    address, _ = client.economic.account.generate_account(client.node.web3,
                                                          10 ** 18 * 10000000)
    value = client.economic.create_staking_limit * 2
    result = client.staking.create_staking(0, address, address, amount=value)
    assert_code(result, 0)
    # Next settlement period
    client.economic.wait_settlement_blocknum(client.node)
    verifierlist = get_pledge_list(client.ppos.getVerifierList)
    log.info("verifierlist:{}".format(verifierlist))
    assert client.node.node_id in verifierlist

    log.info("Close one node")
    client.node.stop()
    node = client_consensus.node
    log.info("The next  periods")
    client.economic.wait_settlement_blocknum(node)
    verifierlist = get_pledge_list(client_consensus.ppos.getVerifierList)
    log.info("verifierlist:{}".format(verifierlist))
    assert client.node.node_id not in verifierlist


@pytest.mark.P2
@pytest.mark.parametrize('status', [0, 1, 2])
def test_CS_CL_014_015_016_029(status, clients_new_node, clients_consensus):
    """
    :param status:
    :param global_test_env:
    :param client_con_list_obj:
    :param clients_noconsensus:
    :return:
    """
    client_noconsensus1 = clients_new_node[0]
    client_noconsensus2 = clients_new_node[1]
    client_noconsensus3 = clients_new_node[2]
    client_noconsensus4 = clients_new_node[3]
    client_consensus1 = clients_consensus[0]
    client_consensus2 = clients_consensus[1]
    client_consensus3 = clients_consensus[2]
    client_consensus4 = clients_consensus[3]

    log.info("The next consensus cycle")
    client_consensus1.economic.wait_consensus_blocknum(client_consensus1.node, number=1)
    blocknumber = client_consensus1.node.eth.blockNumber

    log.info("To report the double sign")
    report_information1 = mock_duplicate_sign(1, client_consensus1.node.nodekey,
                                              client_consensus1.node.blsprikey,
                                              blocknumber)
    log.info("Report information: {}".format(report_information1))

    report_information2 = mock_duplicate_sign(1, client_consensus2.node.nodekey,
                                              client_consensus2.node.blsprikey,
                                              blocknumber)
    log.info("Report information: {}".format(report_information2))

    report_information3 = mock_duplicate_sign(1, client_consensus3.node.nodekey,
                                              client_consensus3.node.blsprikey,
                                              blocknumber)
    log.info("Report information: {}".format(report_information3))

    address_1, _ = client_noconsensus1.economic.account.generate_account(client_noconsensus1.node.web3,
                                                                         10 ** 18 * 10000000)
    address_2, _ = client_noconsensus1.economic.account.generate_account(client_noconsensus1.node.web3,
                                                                         10 ** 18 * 10000000)
    address_3, _ = client_noconsensus1.economic.account.generate_account(client_noconsensus1.node.web3,
                                                                         10 ** 18 * 10000000)
    result = client_consensus1.duplicatesign.reportDuplicateSign(1, report_information1, address_1)
    log.info(result)
    result = client_consensus2.duplicatesign.reportDuplicateSign(1, report_information2, address_2)
    log.info(result)
    result = client_consensus3.duplicatesign.reportDuplicateSign(1, report_information3, address_3)
    log.info(result)
    log.info("The next  periods")
    client_noconsensus2.economic.wait_settlement_blocknum(client_noconsensus2.node)
    validatorlist = get_pledge_list(client_consensus1.ppos.getValidatorList)
    log.info("After being reported validatorlist:{}".format(validatorlist))

    if status == 0:
        address1, _ = client_noconsensus1.economic.account.generate_account(client_noconsensus1.node.web3,
                                                                            10 ** 18 * 10000000)
        value = client_noconsensus1.economic.create_staking_limit * 2
        result = client_noconsensus1.staking.create_staking(0, address1, address1, amount=value)
        assert_code(result, 0)

        address2, _ = client_noconsensus2.economic.account.generate_account(client_noconsensus2.node.web3,
                                                                            10 ** 18 * 10000000)
        value = client_noconsensus2.economic.create_staking_limit * 2
        result = client_noconsensus2.staking.create_staking(0, address2, address2, amount=value)
        assert_code(result, 0)

        log.info("The next  periods")
        client_noconsensus2.economic.wait_settlement_blocknum(client_noconsensus2.node)

        log.info("The next consensus cycle")
        client_noconsensus2.economic.wait_consensus_blocknum(client_noconsensus2.node)

        validatorlist = get_pledge_list(client_noconsensus2.ppos.getValidatorList)
        log.info("validatorlist:{}".format(validatorlist))
        log.info("node1:{}".format(client_noconsensus1.node.node_id))
        log.info("node2:{}".format(client_noconsensus2.node.node_id))
        log.info("node3:{}".format(client_consensus4.node.node_id))
        assert client_noconsensus1.node.node_id in validatorlist
        assert client_noconsensus2.node.node_id in validatorlist
        assert client_consensus4.node.node_id in validatorlist

    if status == 1:
        address1, _ = client_noconsensus1.economic.account.generate_account(client_noconsensus1.node.web3,
                                                                            10 ** 18 * 10000000)
        value = client_noconsensus1.economic.create_staking_limit * 2
        result = client_noconsensus1.staking.create_staking(0, address1, address1, amount=value)
        assert_code(result, 0)

        address2, _ = client_noconsensus1.economic.account.generate_account(client_noconsensus1.node.web3,
                                                                            10 ** 18 * 10000000)
        value = client_noconsensus1.economic.create_staking_limit * 2
        result = client_noconsensus2.staking.create_staking(0, address2, address2, amount=value)
        assert_code(result, 0)

        address3, _ = client_noconsensus1.economic.account.generate_account(client_noconsensus1.node.web3,
                                                                            10 ** 18 * 10000000)
        value = client_noconsensus1.economic.create_staking_limit * 2
        result = client_noconsensus3.staking.create_staking(0, address3, address3, amount=value)
        assert_code(result, 0)
        log.info("The next  periods")
        client_noconsensus2.economic.wait_settlement_blocknum(client_noconsensus2.node)

        log.info("The next consensus cycle")
        client_noconsensus2.economic.wait_consensus_blocknum(client_noconsensus2.node)

        validatorlist = get_pledge_list(client_noconsensus2.ppos.getValidatorList)
        log.info("validatorlist:{}".format(validatorlist))
        log.info("node1:{}".format(client_noconsensus1.node.node_id))
        log.info("node2:{}".format(client_noconsensus2.node.node_id))
        log.info("node3:{}".format(client_noconsensus3.node.node_id))
        log.info("node4:{}".format(client_consensus4.node.node_id))
        assert client_noconsensus1.node.node_id in validatorlist
        assert client_noconsensus2.node.node_id in validatorlist
        assert client_noconsensus3.node.node_id in validatorlist
        assert client_consensus4.node.node_id in validatorlist

    if status == 2:
        address1, _ = client_noconsensus1.economic.account.generate_account(client_noconsensus1.node.web3,
                                                                            10 ** 18 * 10000000)
        value = client_noconsensus1.economic.create_staking_limit * 2
        result = client_noconsensus1.staking.create_staking(0, address1, address1, amount=value)
        assert_code(result, 0)

        address2, _ = client_noconsensus1.economic.account.generate_account(client_noconsensus1.node.web3,
                                                                            10 ** 18 * 10000000)
        value = client_noconsensus1.economic.create_staking_limit * 2
        result = client_noconsensus2.staking.create_staking(0, address2, address2, amount=value)
        assert_code(result, 0)

        address3, _ = client_noconsensus1.economic.account.generate_account(client_noconsensus1.node.web3,
                                                                            10 ** 18 * 10000000)
        value = client_noconsensus1.economic.create_staking_limit * 2
        result = client_noconsensus3.staking.create_staking(0, address3, address3, amount=value)
        assert_code(result, 0)

        address4, _ = client_noconsensus1.economic.account.generate_account(client_noconsensus1.node.web3,
                                                                            10 ** 18 * 10000000)
        value = client_noconsensus1.economic.create_staking_limit * 2
        result = client_noconsensus4.staking.create_staking(0, address4, address4, amount=value)
        assert_code(result, 0)

        log.info("The next  periods")
        client_noconsensus2.economic.wait_settlement_blocknum(client_noconsensus2.node)

        log.info("The next consensus cycle")
        client_noconsensus2.economic.wait_consensus_blocknum(client_noconsensus2.node)

        validatorlist = get_pledge_list(client_noconsensus2.ppos.getValidatorList)
        log.info("validatorlist:{}".format(validatorlist))
        log.info("node:{}".format(clients_consensus[3].node.node_id))
        assert client_consensus4.node.node_id in validatorlist


@pytest.mark.P2
@pytest.mark.parametrize('status', [0, 1])
def test_CS_CL_017_018_019(status, clients_new_node, clients_consensus):
    """
    :param status:
    :param global_test_env:
    :param client_con_list_obj:
    :param clients_noconsensus:
    :return:
    """
    client_noconsensus1 = clients_new_node[0]
    client_noconsensus2 = clients_new_node[1]
    client_consensus1 = clients_consensus[0]
    client_consensus2 = clients_consensus[1]
    client_consensus3 = clients_consensus[2]
    client_consensus4 = clients_consensus[3]

    log.info("The next consensus cycle")
    client_consensus1.economic.wait_consensus_blocknum(client_consensus1.node, number=1)

    validatorlist = get_pledge_list(client_consensus1.ppos.getValidatorList)
    log.info("initial validatorlist:{}".format(validatorlist))
    blocknumber = client_consensus1.node.eth.blockNumber
    log.info("The thrill of being reported{}".format(blocknumber))

    log.info("To report the double sign")
    report_information1 = mock_duplicate_sign(1, client_consensus1.node.nodekey,
                                              client_consensus1.node.blsprikey,
                                              blocknumber)
    log.info("Report information: {}".format(report_information1))

    address, _ = client_noconsensus1.economic.account.generate_account(client_noconsensus1.node.web3,
                                                                       10 ** 18 * 10000000)
    result = client_consensus1.duplicatesign.reportDuplicateSign(1, report_information1, address)
    log.info(result)
    log.info("The next  periods")
    client_noconsensus2.economic.wait_settlement_blocknum(client_noconsensus2.node)
    validatorlist = get_pledge_list(client_consensus1.ppos.getValidatorList)
    log.info("After being reported validatorlist:{}".format(validatorlist))

    if status == 0:
        address1, _ = client_noconsensus1.economic.account.generate_account(client_noconsensus1.node.web3,
                                                                            10 ** 18 * 10000000)
        value = client_noconsensus1.economic.create_staking_limit * 2
        result = client_noconsensus1.staking.create_staking(0, address1, address1, amount=value)
        assert_code(result, 0)
        log.info("The next  periods")
        client_noconsensus2.economic.wait_settlement_blocknum(client_noconsensus2.node)
        log.info("The next consensus cycle")
        client_noconsensus2.economic.wait_consensus_blocknum(client_noconsensus2.node)
        validatorlist = get_pledge_list(client_noconsensus2.ppos.getValidatorList)
        log.info("validatorlist:{}".format(validatorlist))
        assert client_consensus2.node.node_id in validatorlist
        assert client_consensus3.node.node_id in validatorlist
        assert client_consensus4.node.node_id in validatorlist
        assert client_noconsensus1.node.node_id in validatorlist

    if status == 1:
        address1, _ = client_noconsensus1.economic.account.generate_account(client_noconsensus1.node.web3,
                                                                            10 ** 18 * 10000000)
        value = client_noconsensus1.economic.create_staking_limit * 2
        result = client_noconsensus1.staking.create_staking(0, address1, address1, amount=value)
        assert_code(result, 0)

        address2, _ = client_noconsensus1.economic.account.generate_account(client_noconsensus1.node.web3,
                                                                            10 ** 18 * 10000000)
        value = client_noconsensus1.economic.create_staking_limit * 2
        result = client_noconsensus2.staking.create_staking(0, address2, address2, amount=value)
        assert_code(result, 0)

        log.info("The next  periods")
        client_noconsensus2.economic.wait_settlement_blocknum(client_noconsensus2.node)
        log.info("The next consensus cycle")
        client_noconsensus2.economic.wait_consensus_blocknum(client_noconsensus2.node)

        validatorlist = get_pledge_list(client_noconsensus2.ppos.getValidatorList)
        log.info("validatorlist:{}".format(validatorlist))
        assert client_consensus2.node.node_id in validatorlist
        assert client_consensus3.node.node_id in validatorlist
        assert client_consensus4.node.node_id in validatorlist


@pytest.mark.P2
def test_CS_CL_027_028(clients_new_node):
    client_noconsensus1 = clients_new_node[0]
    client_noconsensus2 = clients_new_node[1]

    address1, _ = client_noconsensus1.economic.account.generate_account(client_noconsensus1.node.web3,
                                                                        10 ** 18 * 10000000)
    address2, _ = client_noconsensus1.economic.account.generate_account(client_noconsensus1.node.web3,
                                                                        10 ** 18 * 10000000)

    result = client_noconsensus1.staking.create_staking(0, address1, address1,
                                                        amount=client_noconsensus1.economic.create_staking_limit * 2)
    assert_code(result, 0)

    result = client_noconsensus2.staking.create_staking(0, address2, address2,
                                                        amount=client_noconsensus2.economic.create_staking_limit)
    assert_code(result, 0)

    log.info("Next settlement period")
    client_noconsensus2.economic.wait_settlement_blocknum(client_noconsensus2.node)
    msg = client_noconsensus2.ppos.getVerifierList()
    log.info(msg)
    verifierlist = get_pledge_list(client_noconsensus2.ppos.getVerifierList)
    log.info("verifierlist:{}".format(verifierlist))
    assert client_noconsensus1.node.node_id in verifierlist
    assert client_noconsensus2.node.node_id not in verifierlist


@pytest.mark.P2
def test_CS_CL_033(clients_new_node):
    client = clients_new_node[0]
    address1, _ = client.economic.account.generate_account(client.node.web3,
                                                           10 ** 18 * 10000000)

    value = client.economic.create_staking_limit * 2
    result = client.staking.create_staking(0, address1, address1, amount=value)
    assert_code(result, 0)

    # Next settlement period
    client.economic.wait_settlement_blocknum(client.node)
    # Next consensus period
    client.economic.wait_consensus_blocknum(client.node)
    verifierlist = get_pledge_list(client.ppos.getVerifierList)
    log.info("verifierlist:{}".format(verifierlist))
    assert client.node.node_id in verifierlist
