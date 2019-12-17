from decimal import Decimal
import pytest
from tests.lib.utils import assert_code


def calculate(big_int, mul):
    return int(Decimal(str(big_int)) * Decimal(mul))


@pytest.fixture()
def create_staking_client(client_new_node):
    amount = calculate(client_new_node.economic.create_staking_limit, 5)
    staking_amount = calculate(client_new_node.economic.create_staking_limit, 2)
    staking_address, _ = client_new_node.economic.account.generate_account(client_new_node.node.web3, amount)
    delegate_address, _ = client_new_node.economic.account.generate_account(client_new_node.node.web3,
                                                                            client_new_node.economic.add_staking_limit * 2)
    client_new_node.staking.create_staking(0, staking_address, staking_address, amount=staking_amount)
    setattr(client_new_node, "staking_address", staking_address)
    setattr(client_new_node, "delegate_address", delegate_address)
    setattr(client_new_node, "amount", amount)
    setattr(client_new_node, "staking_amount", staking_amount)
    yield client_new_node


@pytest.fixture()
def staking_delegate_client(client_new_node):
    staking_amount = client_new_node.economic.create_staking_limit
    delegate_amount = client_new_node.economic.add_staking_limit
    staking_address, _ = client_new_node.economic.account.generate_account(client_new_node.node.web3,
                                                                           staking_amount * 2)
    delegate_address, _ = client_new_node.economic.account.generate_account(client_new_node.node.web3,
                                                                            staking_amount * 2)
    result = client_new_node.staking.create_staking(0, staking_address, staking_address)
    assert_code(result, 0)
    result = client_new_node.delegate.delegate(0, delegate_address, amount=delegate_amount * 2)
    assert_code(result, 0)
    msg = client_new_node.ppos.getCandidateInfo(client_new_node.node.node_id)
    staking_blocknum = msg["Ret"]["StakingBlockNum"]
    setattr(client_new_node, "staking_address", staking_address)
    setattr(client_new_node, "delegate_address", delegate_address)
    setattr(client_new_node, "delegate_amount", delegate_amount)
    setattr(client_new_node, "staking_blocknum", staking_blocknum)
    yield client_new_node


@pytest.fixture()
def free_locked_delegate_client(client_new_node):
    staking_amount = client_new_node.economic.create_staking_limit
    delegate_amount = client_new_node.economic.add_staking_limit
    staking_address, _ = client_new_node.economic.account.generate_account(client_new_node.node.web3,
                                                                           staking_amount * 2)
    delegate_address, _ = client_new_node.economic.account.generate_account(client_new_node.node.web3,
                                                                            staking_amount * 2)
    result = client_new_node.staking.create_staking(0, staking_address, staking_address)
    assert_code(result, 0)
    result = client_new_node.delegate.delegate(0, delegate_address, amount=delegate_amount * 2)
    assert_code(result, 0)

    lockup_amount = client_new_node.node.web3.toWei(50, "ether")
    plan = [{'Epoch': 2, 'Amount': lockup_amount}]
    # Create a lock plan
    result = client_new_node.restricting.createRestrictingPlan(delegate_address, plan, delegate_address)
    assert_code(result, 0)
    result = client_new_node.delegate.delegate(1, delegate_address)
    assert_code(result, 0)
    msg = client_new_node.ppos.getCandidateInfo(client_new_node.node.node_id)
    staking_blocknum = msg["Ret"]["StakingBlockNum"]
    setattr(client_new_node, "staking_address", staking_address)
    setattr(client_new_node, "delegate_address", delegate_address)
    setattr(client_new_node, "delegate_amount", delegate_amount)
    setattr(client_new_node, "staking_blocknum", staking_blocknum)
    yield client_new_node


def check_receipt(node, hash, key, expected_result):
    result = node.eth.waitForTransactionReceipt(hash)
    if result["logs"] and key in result["logs"][0]:
        value = result["logs"][0][key]
        assert value == expected_result, "Value contrast error"
    else:
        assert result[key] == expected_result, "Value contrast error"


