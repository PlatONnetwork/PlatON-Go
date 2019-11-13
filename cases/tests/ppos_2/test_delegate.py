# -*- coding: utf-8 -*-

from tests.lib.utils import *
import pytest



def test_DI_001_002(client_new_node_obj):
    """
    :param client_new_node_obj:
    :param get_generate_account:
    :return:
    """
    address, pri_key = client_new_node_obj.economic.account.generate_account(client_new_node_obj.node.web3,
                                                                                  10 ** 18 * 10000000)
    client_new_node_obj.staking.create_staking(0, address, address)
    address1, _ = client_new_node_obj.economic.account.generate_account(client_new_node_obj.node.web3,
                                                                                  10 ** 18 * 10000000)
    result = client_new_node_obj.delegate.delegate(0,address1)
    assert_code(result, 0)
    msg = client_new_node_obj.ppos.getCandidateInfo(client_new_node_obj.node.node_id)
    staking_blocknum = msg["Ret"]["StakingBlockNum"]
    msg = client_new_node_obj.ppos.getDelegateInfo(staking_blocknum, address1, client_new_node_obj.node.node_id)
    log.info(msg)
    assert client_new_node_obj.node.web3.toChecksumAddress(msg["Ret"]["Addr"]) == address1
    assert msg["Ret"]["NodeId"] == client_new_node_obj.node.node_id
    assert msg["Ret"]["ReleasedHes"] == client_new_node_obj.economic.delegate_limit


def test_DI_003_004_005(client_new_node_obj_list):
    """
    :param client_new_node_obj_list:
    :return:
    """
    address, pri_key = client_new_node_obj_list[0].economic.account.generate_account(client_new_node_obj_list[0].node.web3,
                                                                                  10 ** 18 * 10000000)
    client_new_node_obj_list[0].staking.create_staking(0, address, address,amount=1500000000000000000000000)


    address, pri_key = client_new_node_obj_list[1].economic.account.generate_account(client_new_node_obj_list[1].node.web3,
                                                                                  10 ** 18 * 10000000)
    client_new_node_obj_list[1].staking.create_staking(1, address, address,amount=2000000000000000000000000)


    address, pri_key = client_new_node_obj_list[2].economic.account.generate_account(client_new_node_obj_list[2].node.web3,
                                                                                  10 ** 18 * 10000000)
    client_new_node_obj_list[2].staking.create_staking(0, address, address,amount=2500000000000000000000000)


    client_new_node_obj_list[2].economic.wait_settlement_blocknum(client_new_node_obj_list[2].node)
    client_new_node_obj_list[2].economic.wait_consensus_blocknum(client_new_node_obj_list[2].node)



    nodeid_list2 = get_pledge_list(client_new_node_obj_list[2].ppos.getVerifierList)
    log.info(nodeid_list2)

    assert client_new_node_obj_list[0].node.node_id not in nodeid_list2

    address1, _ = client_new_node_obj_list[0].economic.account.generate_account(client_new_node_obj_list[0].node.web3,
                                                                                  10 ** 18 * 10000000)
    #The candidate delegate
    result = client_new_node_obj_list[0].delegate.delegate(0,address1)
    assert_code(result,0)

    assert client_new_node_obj_list[2].node.node_id in nodeid_list2
    address2, _ = client_new_node_obj_list[2].economic.account.generate_account(client_new_node_obj_list[2].node.web3,
                                                                                  10 ** 18 * 10000000)
    #The verifier delegates
    result = client_new_node_obj_list[2].delegate.delegate(0,address2)
    assert_code(result, 0)
    nodeid_list3 = get_pledge_list(client_new_node_obj_list[2].ppos.getValidatorList)
    log.info(nodeid_list3)
    assert client_new_node_obj_list[2].node.node_id in nodeid_list3
    address3, _ = client_new_node_obj_list[2].economic.account.generate_account(client_new_node_obj_list[2].node.web3,
                                                                                  10 ** 18 * 10000000)
    #Consensus verifier delegates
    result = client_new_node_obj_list[2].delegate.delegate(0,address3)
    assert_code(result, 0)


def test_DI_006(get_generate_account, client_consensus_obj):
    address, _ = get_generate_account
    result = client_consensus_obj.delegate.delegate(0, address)
    log.info(result)
    assert_code(result,301107)







if __name__ == '__main__':
    pytest.main(['-s', 'test_delegate.py::test_DI_003_004_005'])























