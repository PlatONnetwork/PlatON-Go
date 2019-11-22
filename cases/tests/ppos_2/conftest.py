from decimal import Decimal
import pytest
from tests.lib.utils import assert_code


def calculate(big_int, mul):
    return int(Decimal(str(big_int)) * Decimal(mul))


@pytest.fixture()
def staking_client(client_new_node_obj):
    amount = calculate(client_new_node_obj.economic.create_staking_limit, 5)
    staking_amount = calculate(client_new_node_obj.economic.create_staking_limit, 2)
    staking_address, _ = client_new_node_obj.economic.account.generate_account(client_new_node_obj.node.web3, amount)
    delegate_address, _ = client_new_node_obj.economic.account.generate_account(client_new_node_obj.node.web3,
                                                                                client_new_node_obj.economic.add_staking_limit * 2)
    client_new_node_obj.staking.create_staking(0, staking_address, staking_address, amount=staking_amount)
    setattr(client_new_node_obj, "staking_address", staking_address)
    setattr(client_new_node_obj, "delegate_address", delegate_address)
    setattr(client_new_node_obj, "amount", amount)
    setattr(client_new_node_obj, "staking_amount", staking_amount)
    yield client_new_node_obj
    client_new_node_obj.economic.env.deploy_all()


@pytest.fixture()
def staking_delegate_client(client_new_node_obj):
    staking_amount = client_new_node_obj.economic.create_staking_limit
    delegate_amount = client_new_node_obj.economic.add_staking_limit
    staking_address, _ = client_new_node_obj.economic.account.generate_account(client_new_node_obj.node.web3,
                                                                               staking_amount * 2)
    delegate_address, _ = client_new_node_obj.economic.account.generate_account(client_new_node_obj.node.web3,
                                                                                staking_amount * 2)
    result = client_new_node_obj.staking.create_staking(0, staking_address, staking_address)
    assert_code(result, 0)
    result = client_new_node_obj.delegate.delegate(0, delegate_address, amount=delegate_amount * 2)
    assert_code(result, 0)
    msg = client_new_node_obj.ppos.getCandidateInfo(client_new_node_obj.node.node_id)
    staking_blocknum = msg["Ret"]["StakingBlockNum"]
    setattr(client_new_node_obj, "staking_address", staking_address)
    setattr(client_new_node_obj, "delegate_address", delegate_address)
    setattr(client_new_node_obj, "delegate_amount", delegate_amount)
    setattr(client_new_node_obj, "staking_blocknum", staking_blocknum)
    yield client_new_node_obj


@pytest.fixture()
def free_locked_delegate_client(client_new_node_obj):
    staking_amount = client_new_node_obj.economic.create_staking_limit
    delegate_amount = client_new_node_obj.economic.add_staking_limit
    staking_address, _ = client_new_node_obj.economic.account.generate_account(client_new_node_obj.node.web3,
                                                                               staking_amount * 2)
    delegate_address, _ = client_new_node_obj.economic.account.generate_account(client_new_node_obj.node.web3,
                                                                                staking_amount * 2)
    result = client_new_node_obj.staking.create_staking(0, staking_address, staking_address)
    assert_code(result, 0)
    result = client_new_node_obj.delegate.delegate(0, delegate_address, amount=delegate_amount * 2)
    assert_code(result, 0)

    lockup_amount = client_new_node_obj.node.web3.toWei(50, "ether")
    plan = [{'Epoch': 2, 'Amount': lockup_amount}]
    # Create a lock plan
    result = client_new_node_obj.restricting.createRestrictingPlan(delegate_address, plan, delegate_address)
    assert_code(result, 0)
    result = client_new_node_obj.delegate.delegate(1, delegate_address)
    assert_code(result, 0)
    msg = client_new_node_obj.ppos.getCandidateInfo(client_new_node_obj.node.node_id)
    staking_blocknum = msg["Ret"]["StakingBlockNum"]
    setattr(client_new_node_obj, "staking_address", staking_address)
    setattr(client_new_node_obj, "delegate_address", delegate_address)
    setattr(client_new_node_obj, "delegate_amount", delegate_amount)
    setattr(client_new_node_obj, "staking_blocknum", staking_blocknum)
    yield client_new_node_obj
