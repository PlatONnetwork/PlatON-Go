import allure
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

def test_transac(global_running_env):
    w3 = global_running_env.get_rand_node().eth
    w3.sendRawTransaction("0xf867058405f5e100825208940fc888cb0269f243bb4531edcb80ed5d39cc260d830186a08081eda06b2cd17a039d9cb8dec9a7881e681d00e708d8e07914bcd13f91566ffcbd0b6ba07628b5646c2e413ef0be79aba8edc51435c2435b032c86b4af38873a21a93635")

@allure.title("signed transaction")
@pytest.mark.P0
@pytest.mark.compatibility
def test_TR_TX_001(global_running_env):
    """
    Test signature transfer transaction
    """
    node = global_running_env.get_rand_node()
    value = 100000
    address, _ = global_running_env.account.generate_account(node.web3, value)
    assert value == node.eth.getBalance(address)


@allure.title("new account singed transaction")
@pytest.mark.P0
def test_TR_TX_006(global_running_env):
    """
    After the new account has a balance, transfer it to another account.
    """
    node = global_running_env.get_rand_node()
    value = 10 ** 18
    address, _ = global_running_env.account.generate_account(node.web3, value)
    assert value == node.eth.getBalance(address)
    to_address, _ = global_running_env.account.generate_account(node.web3)
    new_value = 1000
    global_running_env.account.sendTransaction(node.web3, "", address, to_address, node.eth.gasPrice, 21000, new_value)
    assert new_value == node.eth.getBalance(to_address)


@allure.title("unlock sendtransaction")
@pytest.mark.P1
@pytest.mark.parametrize("node_type", ["consensus", "normal"])
def test_TR_TX_004_to_005(global_running_env, node_type):
    """
    Node unlock transfer
    """
    if node_type == "consensus":
        node = global_running_env.get_rand_node()
    else:
        node = global_running_env.get_a_normal_node()
    account = global_running_env.account
    address = account.generate_account_in_node(node, "88888888", 10 ** 18)
    account.unlock_account(node, address)
    to_address, _ = account.generate_account(node.web3)
    value = 1000
    tx_hash = transaction(node.web3, address, to_address, value, 21000, node.eth.gasPrice)
    node.eth.waitForTransactionReceipt(tx_hash)
    assert value == node.eth.getBalance(to_address)


@allure.title("money negative transaction")
@pytest.mark.P1
def test_TR_TX_008(global_running_env):
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


@allure.title("balance insufficient transaction")
@pytest.mark.P2
@pytest.mark.parametrize("balance", [0, 100, 1000])
def test_TR_TX_002_to_003_and_007(global_running_env, balance):
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


@allure.title("gas too low")
@pytest.mark.P2
def test_TR_TX_010(global_running_env):
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


@allure.title("wallet non local")
@pytest.mark.P3
def test_TR_TX_011(global_running_env):
    """
    Transfer to a non-local wallet
    """
    node = global_running_env.get_rand_node()
    value = 100000
    address, _ = global_running_env.account.generate_account(node.web3, value)
    assert value == node.eth.getBalance(address)


@allure.title("to account abnormal")
@pytest.mark.P3
@pytest.mark.parametrize("to_address", ["", "abcdefghigk"])
def test_TR_TX_012_to_013(global_running_env, to_address):
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
