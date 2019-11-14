# -*- coding: utf-8 -*-
import pytest





@pytest.fixture()
def get_generate_account(client_consensus_obj):
    """
    Get the new wallet and private key
    :param client_consensus_obj:
    :return:
    """
    account = client_consensus_obj.economic.account
    node = client_consensus_obj.node
    address, prikey = account.generate_account(node.web3, 10 ** 18 * 20000000)
    return address, prikey


@pytest.fixture()
def greater_than_staking_amount(global_test_env):
    """
    Gets the Shares pledge amount greater than all verifiers
    :param global_test_env:
    :return:
    """
    node = global_test_env.get_rand_node()
    msg = node.ppos.getVerifierList()
    info_dict = msg["Ret"]
    amount_list = [info["Shares"] for info in info_dict]
    print(amount_list)
    return max(amount_list) + 10 ** 18 * 100000
