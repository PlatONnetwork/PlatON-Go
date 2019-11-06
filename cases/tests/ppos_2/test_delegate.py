# -*- coding: utf-8 -*-

from tests.lib.utils import *
import pytest
import allure

def test_DI_001(client_new_node_obj, get_generate_account):
    """
    :param client_new_node_obj:
    :param get_generate_account:
    :return:
    """

    address, pri_key = get_generate_account
    result = client_new_node_obj.staking.create_staking(0, address, address)
    assert result.get('Code') == 0






