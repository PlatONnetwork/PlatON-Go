# -*- coding: utf-8 -*-

from tests.lib.utils import *
import pytest
import allure


def test_MPI_052_053(client_new_node_obj, get_generate_account):
    """
    Modify node information
    :param client_new_node_obj:
    :param get_generate_account:
    :return:
    """
    external_id = "ID1"
    node_name = "LIDA"
    website = "WEBSITE"
    details = "talent"
    address, pri_key = get_generate_account
    result = client_new_node_obj.staking.create_staking(0, address, address)
    assert result.get('Code') == 0
    result = client_new_node_obj.ppos.editCandidate(address, client_new_node_obj.node.node_id, external_id,
                                                    node_name, website, details, pri_key)
    assert result.get('Code') == 0
    result = client_new_node_obj.ppos.getCandidateInfo(client_new_node_obj.node.node_id)
    log.info(result)
    assert result["Data"]["ExternalId"] == external_id
    assert result["Data"]["NodeName"] == node_name
    assert result["Data"]["Website"] == website
    assert result["Data"]["Details"] == details


def test_MPI_054(client_new_node_obj, get_generate_account, greater_than_staking_amount):
    """
    Node becomes consensus validator when modifying revenue address
    :param client_new_node_obj:
    :param get_generate_account:
    :param greater_than_staking_amount:
    :return:
    """
    address, _ = get_generate_account
    result = client_new_node_obj.staking.create_staking(0, address, address, amount=greater_than_staking_amount)
    assert result.get('Code') == 0
    log.info("Next settlement period")
    client_new_node_obj.economic.wait_settlement_blocknum(client_new_node_obj.node)
    log.info("Next consensus cycle")
    client_new_node_obj.economic.wait_consensus_blocknum(client_new_node_obj.node)
    validator_list = get_pledge_list(client_new_node_obj.ppos.getValidatorList)
    log.info(validator_list)
    assert client_new_node_obj.node.node_id in validator_list
    result = client_new_node_obj.staking.edit_candidate(address, address)
    log.info(result)


def test_MPI_055(client_consensus_obj, get_generate_account):
    """
    The original verifier beneficiary's address modifies the ordinary address
    :param client_consensus_obj:
    :param get_generate_account:
    :return:
    """
    address, _ = get_generate_account
    INCENTIVEPOOL_ADDRESS = client_consensus_obj.economic.cfg.INCENTIVEPOOL_ADDRESS
    DEVELOPER_FOUNDATAION_ADDRESS = client_consensus_obj.economic.cfg.DEVELOPER_FOUNDATAION_ADDRESS

    result = client_consensus_obj.staking.edit_candidate(DEVELOPER_FOUNDATAION_ADDRESS, address)
    log.info(result)
    msg = client_consensus_obj.ppos.getCandidateInfo(client_consensus_obj.node.node_id)
    log.info(msg)
    assert msg["Data"]["BenefitAddress"] == INCENTIVEPOOL_ADDRESS










