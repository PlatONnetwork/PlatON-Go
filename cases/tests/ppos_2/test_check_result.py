# -*- coding: utf-8 -*-
from tests.lib.utils import *
import pytest


@pytest.fixture()
def set_not_need_analyze(client_new_node):
    client_new_node.ppos.need_analyze = False
    yield client_new_node
    client_new_node.ppos.need_analyze = True



@pytest.mark.P3
def test_staking_receipt(set_not_need_analyze):
    external_id = "external_id"
    node_name = "node_name"
    website = "website"
    details = "details"
    node = set_not_need_analyze.node
    economic = set_not_need_analyze.economic
    benifit_address, pri_key = set_not_need_analyze.economic.account.generate_account(node.web3,
                                                                                 economic.create_staking_limit * 2)
    log.info(set_not_need_analyze.ppos.need_analyze)
    result = set_not_need_analyze.ppos.createStaking(0, benifit_address, node.node_id, external_id,
                                                node_name, website,
                                                     details, economic.create_staking_limit,
                                                node.program_version, node.program_version_sign, node.blspubkey,
                                                node.schnorr_NIZK_prove,
                                                pri_key)
    log.info(result)



