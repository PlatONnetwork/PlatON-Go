# -*- coding: utf-8 -*-
from tests.lib.utils import *
import pytest
from common.key import mock_duplicate_sign
from tests.ppos_2.conftest import check_receipt


@pytest.fixture()
def set_not_need_analyze(client_new_node):
    client_new_node.ppos.need_analyze = False
    yield client_new_node
    client_new_node.ppos.need_analyze = True


@pytest.mark.P3
def test_staking_receipt(set_not_need_analyze):
    client = set_not_need_analyze
    node = client.node
    economic = client.economic
    address, _ = client.economic.account.generate_account(node.web3, economic.create_staking_limit * 2)
    hash = client.staking.create_staking(0, address, address)
    log.info(hash)
    key = "topics"
    expected_result = []
    check_receipt(node, hash, key, expected_result)


@pytest.mark.P3
def test_delegate_receipt(set_not_need_analyze):
    client = set_not_need_analyze
    node = client.node
    economic = client.economic
    address, _ = economic.account.generate_account(node.web3, economic.create_staking_limit * 2)
    address, pri_key = economic.account.generate_account(node.web3, economic.delegate_limit * 10)
    hash = client.staking.create_staking(0, address, address)
    node.eth.waitForTransactionReceipt(hash)
    hash = client.delegate.delegate(0, address)
    log.info(hash)
    key = "topics"
    expected_result = []
    check_receipt(node, hash, key, expected_result)


@pytest.mark.P3
def test_withdrewDelegate_receipt(client_new_node):
    client = client_new_node
    node = client.node
    economic = client.economic
    staking_address, _ = economic.account.generate_account(node.web3, economic.create_staking_limit * 2)
    delegate_address, pri_key = economic.account.generate_account(node.web3, economic.delegate_limit * 10)
    client.ppos.need_analyze = True
    client.staking.create_staking(0, staking_address, staking_address)
    client.delegate.delegate(0, delegate_address)
    msg = client.ppos.getCandidateInfo(node.node_id)
    staking_blocknum = msg["Ret"]["StakingBlockNum"]
    client.ppos.need_analyze = False
    hash = client.delegate.withdrew_delegate(staking_blocknum, delegate_address)
    key = "topics"
    expected_result = []
    check_receipt(node, hash, key, expected_result)
    client.ppos.need_analyze = True


@pytest.mark.P3
def test_increase_staking_receipt(set_not_need_analyze):
    client = set_not_need_analyze
    node = client.node
    economic = client.economic
    address, _ = economic.account.generate_account(node.web3, economic.create_staking_limit * 2)
    log.info(client.ppos.need_analyze)

    hash = client.staking.create_staking(0, address, address)
    node.eth.waitForTransactionReceipt(hash)
    hash = client.staking.increase_staking(0, address)
    log.info(hash)
    key = "topics"
    expected_result = []
    check_receipt(node, hash, key, expected_result)


@pytest.mark.P3
def test_edit_candidate_receipt(set_not_need_analyze):
    client = set_not_need_analyze
    node = client.node
    economic = client.economic
    address, pri_key = economic.account.generate_account(node.web3, economic.create_staking_limit * 2)
    hash = client.staking.create_staking(0, address, address)
    node.eth.waitForTransactionReceipt(hash)
    hash = client.staking.edit_candidate(address, address)
    key = "topics"
    expected_result = []
    check_receipt(node, hash, key, expected_result)


@pytest.mark.P3
def test_withdrew_staking_receipt(set_not_need_analyze):
    client = set_not_need_analyze
    node = client.node
    economic = client.economic
    address, _ = economic.account.generate_account(node.web3, economic.create_staking_limit * 2)

    hash = client.staking.create_staking(0, address, address)
    node.eth.waitForTransactionReceipt(hash)
    hash = client.staking.withdrew_staking(address)
    key = "topics"
    expected_result = []
    check_receipt(node, hash, key, expected_result)


@pytest.mark.P3
def test_createRestrictingPlan_receipt(set_not_need_analyze):
    client = set_not_need_analyze
    node = client.node
    economic = client.economic
    address, _ = economic.account.generate_account(node.web3, economic.create_staking_limit * 2)
    lockup_amount = client.node.web3.toWei(20, "ether")
    plan = [{'Epoch': 1, 'Amount': lockup_amount}]
    # Create a lock plan
    hash = client.restricting.createRestrictingPlan(address, plan, address)
    key = "topics"
    expected_result = []
    check_receipt(node, hash, key, expected_result)


@pytest.mark.P3
def test_reportDuplicateSign_receipt(set_not_need_analyze):
    client = set_not_need_analyze
    node = client.node
    economic = client.economic
    address, _ = economic.account.generate_account(node.web3, economic.create_staking_limit * 2)
    number = client.node.eth.blockNumber
    report_information = mock_duplicate_sign(1, client.node.nodekey, client.node.blsprikey, number)
    address_, _ = economic.account.generate_account(node.web3, economic.create_staking_limit * 2)
    hash = client.duplicatesign.reportDuplicateSign(1, report_information, address_)
    key = "topics"
    expected_result = []
    check_receipt(node, hash, key, expected_result)
