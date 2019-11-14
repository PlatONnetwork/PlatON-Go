import pytest
from common.log import log


def transaction(w3, from_address, to_address=None, value=1000000000000000000000, gas=91000000, gasPrice=9000000000):
    params = {
        'to': to_address,
        'from': from_address,
        'gas': gas,
        'gasPrice': gasPrice,
        'value': value
    }
    tx_hash = w3.eth.sendTransaction(params)
    return tx_hash


@pytest.mark.P0
def test_singed_transaction(global_running_env):
    """
    Test signature transfer transaction
    """
    node = global_running_env.get_rand_node()
    value = 100000
    address, _ = global_running_env.account.generate_account(node.web3, value)
    assert value == node.eth.getBalance(address)


@pytest.mark.P0
def test_new_account_singed_transaction(global_running_env):
    """
    After the new account has a balance, transfer it to another account.
    """
    node = global_running_env.get_rand_node()
    value = 10**18
    address, _ = global_running_env.account.generate_account(node.web3, value)
    assert value == node.eth.getBalance(address)
    to_address, _ = global_running_env.account.generate_account(node.web3)
    new_value = 1000
    global_running_env.account.sendTransaction(node.web3, "", address, to_address, node.eth.gasPrice, 21000, new_value)
    assert new_value == node.eth.getBalance(to_address)


@pytest.mark.P1
@pytest.mark.parametrize("node_type", ["consensus", "normal"])
def test_unlock_sendtransaction(global_running_env, node_type):
    """
    Node unlock transfer
    """
    if node_type == "consensus":
        node = global_running_env.get_rand_node()
    else:
        node = global_running_env.get_a_normal_node()
    account = global_running_env.account
    address = account.generate_account_in_node(node, "88888888", 10**18)
    account.unlock_account(node, address)
    to_address, _ = account.generate_account(node.web3)
    value = 1000
    tx_hash = transaction(node.web3, address, to_address, value, 21000, node.eth.gasPrice)
    node.eth.waitForTransactionReceipt(tx_hash)
    assert value == node.eth.getBalance(to_address)


@pytest.mark.P1
def test_money_negative_transaction(global_running_env):
    """
    The transfer amount is negative
    """
    node = global_running_env.get_rand_node()
    account = global_running_env.account
    status = 0
    try:
        account.generate_account(node.web3, -1000)
        status = 1
    except Exception as e:
        log.info("transaction error:{}".format(e))
    assert status == 0


@pytest.mark.P2
@pytest.mark.parametrize("balance", [0, 100])
def test_balance_insufficient_transaction(global_running_env, balance):
    """
    Balance is 0 transfer
    Insufficient balance transfer
    """
    node = global_running_env.get_rand_node()
    account = global_running_env.account
    address, _ = account.generate_account(node.web3, balance)
    to_address, _ = account.generate_account(node.web3)
    value = 1000
    status = 0
    try:
        tx_res = account.sendTransaction(node.web3, "", address, to_address, node.eth.gasPrice, 21000, value)
        status = 1
    except Exception as e:
        log.info("transaction error:{}".format(e))
    assert status == 0


@pytest.mark.P2
def test_gas_to_low(global_running_env):
    """
    Gas is too low
    """
    node = global_running_env.get_rand_node()
    account = global_running_env.account
    address = account.account_with_money["address"]
    to_address, _ = account.generate_account(node.web3)
    value = 1000
    status = 0
    try:
        tx_res = account.sendTransaction(node.web3, "", address, to_address, node.eth.gasPrice, 210, value)
        status = 1
    except Exception as e:
        log.info("transaction error:{}".format(e))
    assert status == 0


@pytest.mark.P3
def test_wallet_non_local(global_running_env):
    """
    Transfer to a non-local wallet
    """
    node = global_running_env.get_rand_node()
    value = 100000
    address, _ = global_running_env.account.generate_account(node.web3, value)
    assert value == node.eth.getBalance(address)


@pytest.mark.P3
@pytest.mark.parametrize("to_address", ["", "abcdefghigk"])
def test_to_account_abnormal(global_running_env, to_address):
    """
    To address abnormal
    """
    node = global_running_env.get_rand_node()
    account = global_running_env.account
    from_address = account.account_with_money["address"]
    status = 0
    try:
        tx_res = account.sendTransaction(node.web3, "", from_address, to_address, node.eth.gasPrice, 21000, 100, False)
        status = 1
    except Exception as e:
        log.info("transaction error:{}".format(e))
    assert status == 0
