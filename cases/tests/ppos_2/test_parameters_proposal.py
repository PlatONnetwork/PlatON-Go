# -*- coding: utf-8 -*-
from tests.conftest import param_governance_verify_before_endblock
from tests.lib.utils import *
import pytest
from dacite import from_dict
from common.log import log
from tests.lib import Genesis


@pytest.mark.P2
def test_POP_001_003(client_consensus_obj, client_new_node_obj):
    """
    Increase the threshold of pledge
    :param client_consensus_obj:
    :param get_generate_account:
    :param client_new_node_obj:
    :return:
    """
    client_consensus_obj.economic.env.deploy_all()
    block = param_governance_verify_before_endblock(client_consensus_obj, "staking", "stakeThreshold",
                                                    "1800000000000000000000000")
    log.info(block)
    address, _ = client_new_node_obj.economic.account.generate_account(client_new_node_obj.node.web3,
                                                                       10 ** 18 * 10000000)
    wait_block_number(client_new_node_obj.node, block)
    result = client_new_node_obj.staking.create_staking(0, address, address, amount=179000000000000000000000)
    log.info(result)
    assert_code(result, 301100)
    result = client_new_node_obj.staking.create_staking(0, address, address, amount=1800000000000000000000000)
    log.info(result)
    assert_code(result, 0)
    verifier_list = get_pledge_list(client_new_node_obj.ppos.getVerifierList)
    log.info(verifier_list)
    client_new_node_obj.economic.wait_settlement_blocknum(client_new_node_obj.node)
    verifier_list = get_pledge_list(client_new_node_obj.ppos.getVerifierList)
    assert client_new_node_obj.node.node_id in verifier_list


@pytest.mark.P2
def test_POP_002(client_consensus_obj, client_new_node_obj, new_genesis_env):
    """
    Minimum pledge reduced pledge threshold
    :param client_consensus_obj:
    :param client_new_node_obj:
    :param new_genesis_env:
    :return:
    """
    genesis = from_dict(data_class=Genesis, data=new_genesis_env.genesis_config)
    genesis.economicModel.staking.stakeThreshold = 1500000000000000000000000
    new_genesis_env.set_genesis(genesis.to_dict())
    new_genesis_env.deploy_all()

    block = param_governance_verify_before_endblock(client_consensus_obj, "staking", "stakeThreshold",
                                                    "1000000000000000000000000")
    log.info(block)
    wait_block_number(client_new_node_obj.node, block)
    address, _ = client_new_node_obj.economic.account.generate_account(client_new_node_obj.node.web3,
                                                                       10 ** 18 * 10000000)
    result = client_new_node_obj.staking.create_staking(0, address, address, amount=990000000000000000000000)
    log.info(result)
    assert_code(result, 301100)
    result = client_new_node_obj.staking.create_staking(0, address, address, amount=1000000000000000000000000)
    log.info(result)
    assert_code(result, 0)


@pytest.mark.P2
def test_POP_005(client_consensus_obj, client_new_node_obj, new_genesis_env):
    """
    The amendment of the threshold of pledge shall not take effect
    :param client_consensus_obj:
    :param get_generate_account:
    :param client_new_node_obj:
    :return:
    """
    genesis = from_dict(data_class=Genesis, data=new_genesis_env.genesis_config)
    genesis.economicModel.staking.stakeThreshold = 1000000000000000000000000
    new_genesis_env.set_genesis(genesis.to_dict())
    new_genesis_env.deploy_all()
    block = param_governance_verify_before_endblock(client_consensus_obj, "staking", "stakeThreshold",
                                                    "1500000000000000000000000", effectiveflag=False)
    log.info(block)
    address, _ = client_new_node_obj.economic.account.generate_account(client_new_node_obj.node.web3,
                                                                       10 ** 18 * 10000000)
    wait_block_number(client_new_node_obj.node, block)
    result = client_new_node_obj.staking.create_staking(0, address, address, amount=1000000000000000000000000)
    assert_code(result, 0)


@pytest.mark.P2
def test_POP_006(client_consensus_obj, client_new_node_obj, new_genesis_env):
    """
    (hesitation period) increase - entrustment overweight threshold
    :param client_consensus_obj:
    :param client_new_node_obj:
    :param new_genesis_env:
    :return:
    """
    genesis = from_dict(data_class=Genesis, data=new_genesis_env.genesis_config)
    genesis.economicModel.staking.operatingThreshold = 10000000000000000000
    new_genesis_env.set_genesis(genesis.to_dict())
    new_genesis_env.deploy_all()

    block = param_governance_verify_before_endblock(client_consensus_obj, "staking", "operatingThreshold",
                                                    "11000000000000000000")
    log.info(block)
    address1, _ = client_new_node_obj.economic.account.generate_account(client_new_node_obj.node.web3,
                                                                        10 ** 18 * 10000000)
    wait_block_number(client_new_node_obj.node, block)
    result = client_new_node_obj.staking.create_staking(0, address1, address1)
    assert_code(result, 0)
    address2, _ = client_new_node_obj.economic.account.generate_account(client_new_node_obj.node.web3,
                                                                        10 ** 18 * 10000000)
    result = client_new_node_obj.delegate.delegate(0, address2, amount=10000000000000000000)
    assert_code(result, 301105)
    result = client_new_node_obj.delegate.delegate(0, address2, amount=11000000000000000000)
    assert_code(result, 0)
    result = client_new_node_obj.staking.increase_staking(0, address1, amount=10000000000000000000)
    assert_code(result, 301104)
    result = client_new_node_obj.staking.increase_staking(0, address1, amount=11000000000000000000)
    assert_code(result, 0)


@pytest.mark.P2
def test_POP_007(client_consensus_obj, client_new_node_obj, new_genesis_env):
    """
    (lockup period) increase - entrustment overweight threshold
    :param client_consensus_obj:
    :param client_new_node_obj:
    :param new_genesis_env:
    :return:
    """
    genesis = from_dict(data_class=Genesis, data=new_genesis_env.genesis_config)
    genesis.economicModel.staking.operatingThreshold = 10000000000000000000
    new_genesis_env.set_genesis(genesis.to_dict())
    new_genesis_env.deploy_all()

    block = param_governance_verify_before_endblock(client_consensus_obj, "staking", "operatingThreshold",
                                                    "11000000000000000000")
    log.info(block)
    address1, _ = client_new_node_obj.economic.account.generate_account(client_new_node_obj.node.web3,
                                                                        10 ** 18 * 10000000)
    wait_block_number(client_new_node_obj.node, block)
    result = client_new_node_obj.staking.create_staking(0, address1, address1)
    assert_code(result, 0)
    address2, _ = client_new_node_obj.economic.account.generate_account(client_new_node_obj.node.web3,
                                                                        10 ** 18 * 10000000)
    client_new_node_obj.economic.wait_settlement_blocknum(client_new_node_obj.node)

    result = client_new_node_obj.delegate.delegate(0, address2, amount=10000000000000000000)
    assert_code(result, 301105)
    result = client_new_node_obj.delegate.delegate(0, address2, amount=11000000000000000000)
    assert_code(result, 0)
    result = client_new_node_obj.staking.increase_staking(0, address1, amount=10000000000000000000)
    assert_code(result, 301104)
    result = client_new_node_obj.staking.increase_staking(0, address1, amount=11000000000000000000)
    assert_code(result, 0)


@pytest.mark.P2
def test_POP_008(client_consensus_obj, client_new_node_obj, new_genesis_env):
    """
    (hesitation period) reduce the entrustment overweight threshold - test
    :param client_consensus_obj:
    :param client_new_node_obj:
    :param new_genesis_env:
    :return:
    """
    genesis = from_dict(data_class=Genesis, data=new_genesis_env.genesis_config)
    genesis.economicModel.staking.operatingThreshold = 15000000000000000000
    new_genesis_env.set_genesis(genesis.to_dict())
    new_genesis_env.deploy_all()

    block = param_governance_verify_before_endblock(client_consensus_obj, "staking", "operatingThreshold",
                                                    "14000000000000000000")
    log.info(block)
    address1, _ = client_new_node_obj.economic.account.generate_account(client_new_node_obj.node.web3,
                                                                        10 ** 18 * 10000000)
    wait_block_number(client_new_node_obj.node, block)
    result = client_new_node_obj.staking.create_staking(0, address1, address1)
    assert_code(result, 0)
    address2, _ = client_new_node_obj.economic.account.generate_account(client_new_node_obj.node.web3,
                                                                        10 ** 18 * 10000000)
    result = client_new_node_obj.delegate.delegate(0, address2, amount=13000000000000000000)
    assert_code(result, 301105)
    result = client_new_node_obj.delegate.delegate(0, address2, amount=14000000000000000000)
    assert_code(result, 0)
    result = client_new_node_obj.staking.increase_staking(0, address1, amount=13000000000000000000)
    assert_code(result, 301104)
    result = client_new_node_obj.staking.increase_staking(0, address1, amount=14000000000000000000)
    assert_code(result, 0)


@pytest.mark.P2
def test_POP_009(client_consensus_obj, client_new_node_obj, new_genesis_env):
    """
    (lockup period) reduce the entrustment increase threshold - test
    :param client_consensus_obj:
    :param client_new_node_obj:
    :param new_genesis_env:
    :return:
    """
    genesis = from_dict(data_class=Genesis, data=new_genesis_env.genesis_config)
    genesis.economicModel.staking.operatingThreshold = 15000000000000000000
    new_genesis_env.set_genesis(genesis.to_dict())
    new_genesis_env.deploy_all()

    block = param_governance_verify_before_endblock(client_consensus_obj, "staking", "operatingThreshold",
                                                    "14000000000000000000")
    log.info(block)
    address1, _ = client_new_node_obj.economic.account.generate_account(client_new_node_obj.node.web3,
                                                                        10 ** 18 * 10000000)
    wait_block_number(client_new_node_obj.node, block)
    result = client_new_node_obj.staking.create_staking(0, address1, address1)
    assert_code(result, 0)
    address2, _ = client_new_node_obj.economic.account.generate_account(client_new_node_obj.node.web3,
                                                                        10 ** 18 * 10000000)
    log.info("The next cycle")
    client_new_node_obj.economic.wait_settlement_blocknum(client_new_node_obj.node)
    result = client_new_node_obj.delegate.delegate(0, address2, amount=13000000000000000000)
    assert_code(result, 301105)
    result = client_new_node_obj.delegate.delegate(0, address2, amount=14000000000000000000)
    assert_code(result, 0)
    result = client_new_node_obj.staking.increase_staking(0, address1, amount=13000000000000000000)
    assert_code(result, 301104)
    result = client_new_node_obj.staking.increase_staking(0, address1, amount=14000000000000000000)
    assert_code(result, 0)


@pytest.mark.P2
def test_POP_010_011(client_consensus_obj, client_new_node_obj, new_genesis_env):
    """
    (hesitation period, lockup period) free amount initiate revocation entrustment
    :param client_consensus_obj:
    :param client_new_node_obj:
    :param new_genesis_env:
    :return:
    """
    genesis = from_dict(data_class=Genesis, data=new_genesis_env.genesis_config)
    genesis.economicModel.staking.operatingThreshold = 10000000000000000000
    new_genesis_env.set_genesis(genesis.to_dict())
    new_genesis_env.deploy_all()

    address1, _ = client_new_node_obj.economic.account.generate_account(client_new_node_obj.node.web3,
                                                                        10 ** 18 * 10000000)
    address_delegate_1, _ = client_new_node_obj.economic.account.generate_account(client_new_node_obj.node.web3,
                                                                                  10 ** 18 * 10000000)
    result = client_new_node_obj.staking.create_staking(0, address1, address1)
    assert_code(result, 0)
    delegate_amount = 50000000000000000000
    result = client_new_node_obj.delegate.delegate(0, address_delegate_1, amount=delegate_amount)
    assert_code(result, 0)
    amount1_before = client_new_node_obj.node.eth.getBalance(address_delegate_1)
    log.info("The wallet balance:{}".format(amount1_before))
    block = param_governance_verify_before_endblock(client_consensus_obj, "staking", "operatingThreshold",
                                                    "20000000000000000000")
    wait_block_number(client_new_node_obj.node, block)
    log.info("The delegate is initiated after the parameter takes effect")
    address_delegate_2, _ = client_new_node_obj.economic.account.generate_account(client_new_node_obj.node.web3,
                                                                                  10 ** 18 * 10000000)
    result = client_new_node_obj.delegate.delegate(0, address_delegate_2, amount=delegate_amount)
    assert_code(result, 0)
    amount2_before = client_new_node_obj.node.eth.getBalance(address_delegate_1)
    log.info("The wallet balance:{}".format(amount2_before))

    msg = client_new_node_obj.ppos.getCandidateInfo(client_new_node_obj.node.node_id)
    staking_blocknum = msg["Ret"]["StakingBlockNum"]

    result = client_new_node_obj.delegate.withdrew_delegate(staking_blocknum, address_delegate_1,
                                                            amount=31000000000000000000)
    assert_code(result, 0)
    amount1_after = client_new_node_obj.node.eth.getBalance(address_delegate_1)
    log.info("The wallet balance:{}".format(amount1_after))
    assert amount1_before - amount1_after < client_new_node_obj.node.web3.toWei(1, "ether")

    result = client_new_node_obj.delegate.withdrew_delegate(staking_blocknum, address_delegate_2,
                                                            amount=31000000000000000000)
    assert_code(result, 0)
    amount2_after = client_new_node_obj.node.eth.getBalance(address_delegate_1)
    log.info("The wallet balance:{}".format(amount2_after))
    assert amount1_before - amount1_after < client_new_node_obj.node.web3.toWei(1, "ether")


@pytest.mark.P2
def test_POP_012(client_consensus_obj, client_new_node_obj, new_genesis_env):
    """
    :param client_consensus_obj:
    :param client_new_node_obj:
    :param new_genesis_env:
    :return:
    """
    genesis = from_dict(data_class=Genesis, data=new_genesis_env.genesis_config)
    genesis.economicModel.staking.operatingThreshold = 10000000000000000000
    new_genesis_env.set_genesis(genesis.to_dict())
    new_genesis_env.deploy_all()

    address, _ = client_new_node_obj.economic.account.generate_account(client_new_node_obj.node.web3,
                                                                       10 ** 18 * 10000000)
    result = client_new_node_obj.staking.create_staking(0, address, address)

    assert_code(result, 0)
    address_delegate, _ = client_new_node_obj.economic.account.generate_account(client_new_node_obj.node.web3,
                                                                                10 ** 18 * 10000000)
    delegate_amount = 30000000000000000000
    result = client_new_node_obj.delegate.delegate(0, address_delegate, amount=delegate_amount)
    assert_code(result, 0)

    block = param_governance_verify_before_endblock(client_consensus_obj, "staking", "operatingThreshold",
                                                    "20000000000000000000")
    wait_block_number(client_new_node_obj.node, block)

    plan = [{'Epoch': 1, 'Amount': delegate_amount}]
    result = client_new_node_obj.restricting.createRestrictingPlan(address_delegate, plan, address_delegate)
    assert_code(result, 0)

    result = client_new_node_obj.delegate.delegate(1, address_delegate, amount=delegate_amount)
    assert_code(result, 0)
    amount_before = client_new_node_obj.node.eth.getBalance(address_delegate)
    log.info("The wallet balance:{}".format(amount_before))

    msg = client_new_node_obj.ppos.getCandidateInfo(client_new_node_obj.node.node_id)
    staking_blocknum = msg["Ret"]["StakingBlockNum"]
    withdrew_delegate = 41000000000000000000
    result = client_new_node_obj.delegate.withdrew_delegate(staking_blocknum, address_delegate,
                                                            amount=withdrew_delegate)
    assert_code(result, 0)
    amount1_after = client_new_node_obj.node.eth.getBalance(address_delegate)
    log.info("The wallet balance:{}".format(amount1_after))
    amount_dill = amount1_after - amount_before
    assert delegate_amount - amount_dill < client_new_node_obj.node.web3.toWei(1, "ether")
    client_new_node_obj.economic.wait_settlement_blocknum(client_new_node_obj.node)

    amount1_last = client_new_node_obj.node.eth.getBalance(address_delegate)
    log.info("The wallet balance:{}".format(amount1_last))
    assert amount1_last - amount1_after == delegate_amount
    assert delegate_amount * 2 - (amount1_last - amount_before) < client_new_node_obj.node.web3.toWei(1, "ether")


@pytest.mark.P2
def test_POP_013(client_consensus_obj, client_new_node_obj, new_genesis_env):
    """
    :param client_consensus_obj:
    :param client_new_node_obj:
    :param new_genesis_env:
    :return:
    """
    genesis = from_dict(data_class=Genesis, data=new_genesis_env.genesis_config)
    genesis.economicModel.staking.operatingThreshold = 20000000000000000000
    new_genesis_env.set_genesis(genesis.to_dict())
    new_genesis_env.deploy_all()

    block = param_governance_verify_before_endblock(client_consensus_obj, "staking", "operatingThreshold",
                                                    "10000000000000000000", effectiveflag=False)
    log.info(block)
    address1, _ = client_new_node_obj.economic.account.generate_account(client_new_node_obj.node.web3,
                                                                        10 ** 18 * 10000000)
    wait_block_number(client_new_node_obj.node, block)
    result = client_new_node_obj.staking.create_staking(0, address1, address1)
    assert_code(result, 0)

    address2, _ = client_new_node_obj.economic.account.generate_account(client_new_node_obj.node.web3,
                                                                        10 ** 18 * 10000000)
    result = client_new_node_obj.delegate.delegate(0, address2, amount=10000000000000000000)
    log.info(result)
    assert_code(result, 301105)
    result = client_new_node_obj.delegate.delegate(0, address2, amount=20000000000000000000)
    assert_code(result, 0)
    result = client_new_node_obj.staking.increase_staking(0, address1, amount=10000000000000000000)
    assert_code(result, 301104)
    result = client_new_node_obj.staking.increase_staking(0, address1, amount=20000000000000000000)
    assert_code(result, 0)


@pytest.mark.P2
def test_POP_014(client_consensus_obj, client_new_node_obj, new_genesis_env):
    """
    Parameter not in effect - initiate redemption
    :param client_consensus_obj:
    :param client_new_node_obj:
    :param new_genesis_env:
    :return:
    """
    new_genesis_env.deploy_all()
    address1, _ = client_new_node_obj.economic.account.generate_account(client_new_node_obj.node.web3,
                                                                        10 ** 18 * 10000000)
    address_delegate_1, _ = client_new_node_obj.economic.account.generate_account(client_new_node_obj.node.web3,
                                                                                  10 ** 18 * 10000000)
    result = client_new_node_obj.staking.create_staking(0, address1, address1)
    assert_code(result, 0)
    delegate_amount = 30000000000000000000
    result = client_new_node_obj.delegate.delegate(0, address_delegate_1, amount=delegate_amount)
    assert_code(result, 0)
    amount1_before = client_new_node_obj.node.eth.getBalance(address_delegate_1)
    log.info("The wallet balance:{}".format(amount1_before))
    block = param_governance_verify_before_endblock(client_consensus_obj, "staking", "operatingThreshold",
                                                    "20000000000000000000", effectiveflag=False)
    wait_block_number(client_new_node_obj.node, block)
    log.info("The delegate is initiated after the parameter takes effect")
    address_delegate_2, _ = client_new_node_obj.economic.account.generate_account(client_new_node_obj.node.web3,
                                                                                  10 ** 18 * 10000000)
    result = client_new_node_obj.delegate.delegate(0, address_delegate_2, amount=delegate_amount)
    assert_code(result, 0)
    amount2_before = client_new_node_obj.node.eth.getBalance(address_delegate_1)
    log.info("The wallet balance:{}".format(amount2_before))

    msg = client_new_node_obj.ppos.getCandidateInfo(client_new_node_obj.node.node_id)
    staking_blocknum = msg["Ret"]["StakingBlockNum"]
    withdrew_delegate = 11000000000000000000
    result = client_new_node_obj.delegate.withdrew_delegate(staking_blocknum, address_delegate_1,
                                                            amount=withdrew_delegate)
    assert_code(result, 0)
    amount1_after = client_new_node_obj.node.eth.getBalance(address_delegate_1)
    log.info("The wallet balance:{}".format(amount1_after))
    amount1_dill = amount1_after - amount1_before
    assert withdrew_delegate - amount1_dill < client_new_node_obj.node.web3.toWei(1, "ether")

    result = client_new_node_obj.delegate.withdrew_delegate(staking_blocknum, address_delegate_2,
                                                            amount=withdrew_delegate)
    assert_code(result, 0)
    amount2_after = client_new_node_obj.node.eth.getBalance(address_delegate_2)
    log.info("The wallet balance:{}".format(amount2_after))
    amount1_dill = amount2_after - amount2_before
    assert withdrew_delegate - amount1_dill < client_new_node_obj.node.web3.toWei(1, "ether")


@pytest.mark.P2
def test_POP_015(client_consensus_obj, client_noc_list_obj, new_genesis_env):
    """
    Increase the number of alternative nodes
    :param client_consensus_obj:
    :param client_new_node_obj:
    :param new_genesis_env:
    :return:
    """
    genesis = from_dict(data_class=Genesis, data=new_genesis_env.genesis_config)
    genesis.economicModel.staking.maxValidators = 5
    new_genesis_env.set_genesis(genesis.to_dict())
    new_genesis_env.deploy_all()

    address, _ = client_noc_list_obj[0].economic.account.generate_account(client_noc_list_obj[0].node.web3,
                                                                          10 ** 18 * 10000000)
    result = client_noc_list_obj[0].staking.create_staking(0, address, address, amount=1600000000000000000000000)
    assert_code(result, 0)
    block = param_governance_verify_before_endblock(client_consensus_obj, "staking", "maxValidators",
                                                    "6")
    wait_block_number(client_noc_list_obj[1].node, block)
    address1, _ = client_noc_list_obj[1].economic.account.generate_account(client_noc_list_obj[1].node.web3,
                                                                           10 ** 18 * 10000000)
    result = client_noc_list_obj[1].staking.create_staking(0, address1, address1, amount=1700000000000000000000000)
    assert_code(result, 0)

    client_noc_list_obj[1].economic.wait_settlement_blocknum(client_noc_list_obj[1].node)
    getVerifierList = get_pledge_list(client_noc_list_obj[1].node.ppos.getVerifierList)
    log.info(getVerifierList)
    assert len(getVerifierList) == 6
    node_id_1 = client_noc_list_obj[0].node.node_id
    node_id_2 = client_noc_list_obj[1].node.node_id
    assert node_id_1 in getVerifierList
    assert node_id_2 in getVerifierList


@pytest.mark.P2
def test_POP_016(client_consensus_obj, client_noc_list_obj, new_genesis_env):
    """
    Reduce the number of alternative nodes
    :param client_consensus_obj:
    :param client_new_node_obj:
    :param new_genesis_env:
    :return:
    """
    genesis = from_dict(data_class=Genesis, data=new_genesis_env.genesis_config)
    genesis.economicModel.staking.maxValidators = 6
    new_genesis_env.set_genesis(genesis.to_dict())
    new_genesis_env.deploy_all()

    block = param_governance_verify_before_endblock(client_consensus_obj, "staking", "maxValidators",
                                                    "5")
    wait_block_number(client_noc_list_obj[1].node, block)
    address, _ = client_noc_list_obj[0].economic.account.generate_account(client_noc_list_obj[0].node.web3,
                                                                          10 ** 18 * 10000000)
    result = client_noc_list_obj[0].staking.create_staking(0, address, address, amount=1600000000000000000000000)
    assert_code(result, 0)
    address1, _ = client_noc_list_obj[1].economic.account.generate_account(client_noc_list_obj[1].node.web3,
                                                                           10 ** 18 * 10000000)
    result = client_noc_list_obj[1].staking.create_staking(0, address1, address1, amount=1700000000000000000000000)
    assert_code(result, 0)

    client_noc_list_obj[1].economic.wait_settlement_blocknum(client_noc_list_obj[1].node)
    getVerifierList = get_pledge_list(client_noc_list_obj[1].node.ppos.getVerifierList)
    log.info(getVerifierList)
    assert len(getVerifierList) == 5


@pytest.mark.P2
def test_POP_017(client_consensus_obj, client_noc_list_obj, new_genesis_env):
    """
    Increase the number of node candidates - not active
    :param client_consensus_obj:
    :param client_new_node_obj:
    :param new_genesis_env:
    :return:
    """
    genesis = from_dict(data_class=Genesis, data=new_genesis_env.genesis_config)
    genesis.economicModel.staking.maxValidators = 5
    new_genesis_env.set_genesis(genesis.to_dict())
    new_genesis_env.deploy_all()
    block = param_governance_verify_before_endblock(client_consensus_obj, "staking", "maxValidators",
                                                    "6", effectiveflag=False)
    wait_block_number(client_noc_list_obj[1].node, block)
    address, _ = client_noc_list_obj[0].economic.account.generate_account(client_noc_list_obj[0].node.web3,
                                                                          10 ** 18 * 10000000)
    result = client_noc_list_obj[0].staking.create_staking(0, address, address, amount=1600000000000000000000000)
    assert_code(result, 0)
    address1, _ = client_noc_list_obj[1].economic.account.generate_account(client_noc_list_obj[1].node.web3,
                                                                           10 ** 18 * 10000000)
    result = client_noc_list_obj[1].staking.create_staking(0, address1, address1, amount=1700000000000000000000000)
    assert_code(result, 0)

    client_noc_list_obj[1].economic.wait_settlement_blocknum(client_noc_list_obj[1].node)
    getVerifierList = get_pledge_list(client_noc_list_obj[1].node.ppos.getVerifierList)
    log.info(getVerifierList)
    assert len(getVerifierList) == 5


@pytest.mark.P2
def test_POP_018(client_consensus_obj, client_noc_list_obj, new_genesis_env):
    """
    Reduce the number of node candidates - not in effect
    :param client_consensus_obj:
    :param client_new_node_obj:
    :param new_genesis_env:
    :return:
    """
    genesis = from_dict(data_class=Genesis, data=new_genesis_env.genesis_config)
    genesis.economicModel.staking.maxValidators = 6
    new_genesis_env.set_genesis(genesis.to_dict())
    new_genesis_env.deploy_all()

    address, _ = client_noc_list_obj[0].economic.account.generate_account(client_noc_list_obj[0].node.web3,
                                                                          10 ** 18 * 10000000)
    result = client_noc_list_obj[0].staking.create_staking(0, address, address, amount=1600000000000000000000000)
    assert_code(result, 0)
    block = param_governance_verify_before_endblock(client_consensus_obj, "staking", "maxValidators",
                                                    "5", effectiveflag=False)
    wait_block_number(client_noc_list_obj[1].node, block)
    address1, _ = client_noc_list_obj[1].economic.account.generate_account(client_noc_list_obj[1].node.web3,
                                                                           10 ** 18 * 10000000)
    result = client_noc_list_obj[1].staking.create_staking(0, address1, address1, amount=1700000000000000000000000)
    assert_code(result, 0)

    client_noc_list_obj[1].economic.wait_settlement_blocknum(client_noc_list_obj[1].node)
    getVerifierList = get_pledge_list(client_noc_list_obj[1].node.ppos.getVerifierList)
    log.info(getVerifierList)
    assert len(getVerifierList) == 6
    node_id_1 = client_noc_list_obj[0].node.node_id
    node_id_2 = client_noc_list_obj[1].node.node_id
    log.info(node_id_1, node_id_2)
    assert node_id_1 in getVerifierList
    assert node_id_2 in getVerifierList


@pytest.mark.P2
def test_POP_019(client_consensus_obj, client_new_node_obj, new_genesis_env):
    """
    Increased lock - up threshold
    :param client_consensus_obj:
    :param client_new_node_obj:
    :param new_genesis_env:
    :return:
    """
    new_genesis_env.deploy_all()
    address, _ = client_new_node_obj.economic.account.generate_account(client_new_node_obj.node.web3,
                                                                       10 ** 18 * 10000000)
    result = client_new_node_obj.staking.create_staking(0, address, address)
    assert_code(result, 0)
    block = param_governance_verify_before_endblock(client_consensus_obj, "staking", "unStakeFreezeDuration",
                                                    "3")
    wait_block_number(client_new_node_obj.node, block)

    result = client_new_node_obj.staking.withdrew_staking(address)
    assert_code(result, 0)
    amount1 = client_new_node_obj.node.eth.getBalance(address)
    log.info("The wallet balance:{}".format(amount1))
    client_new_node_obj.economic.wait_settlement_blocknum(client_new_node_obj.node, number=2)
    amount2 = client_new_node_obj.node.eth.getBalance(address)
    log.info("The wallet balance:{}".format(amount2))
    client_new_node_obj.economic.wait_settlement_blocknum(client_new_node_obj.node)
    amount3 = client_new_node_obj.node.eth.getBalance(address)
    log.info("The wallet balance:{}".format(amount3))
    staking_amount = client_new_node_obj.economic.create_staking_limit
    assert amount3 - amount2 == staking_amount


@pytest.mark.P2
def test_POP_020(client_consensus_obj, client_new_node_obj, new_genesis_env):
    """
    Reduce lock - up threshold
    :param client_consensus_obj:
    :param client_new_node_obj:
    :param new_genesis_env:
    :return:
    """
    genesis = from_dict(data_class=Genesis, data=new_genesis_env.genesis_config)
    genesis.economicModel.staking.unStakeFreezeDuration = 3
    new_genesis_env.set_genesis(genesis.to_dict())
    new_genesis_env.deploy_all()
    address, _ = client_new_node_obj.economic.account.generate_account(client_new_node_obj.node.web3,
                                                                       10 ** 18 * 10000000)
    result = client_new_node_obj.staking.create_staking(0, address, address)
    assert_code(result, 0)

    block = param_governance_verify_before_endblock(client_consensus_obj, "staking", "unStakeFreezeDuration",
                                                    "2")
    wait_block_number(client_new_node_obj.node, block)

    result = client_new_node_obj.staking.withdrew_staking(address)
    assert_code(result, 0)
    amount1 = client_new_node_obj.node.eth.getBalance(address)
    log.info("The wallet balance:{}".format(amount1))
    client_new_node_obj.economic.wait_settlement_blocknum(client_new_node_obj.node, number=1)
    amount2 = client_new_node_obj.node.eth.getBalance(address)
    log.info("The wallet balance:{}".format(amount2))
    client_new_node_obj.economic.wait_settlement_blocknum(client_new_node_obj.node)
    amount3 = client_new_node_obj.node.eth.getBalance(address)
    log.info("The wallet balance:{}".format(amount3))
    staking_amount = client_new_node_obj.economic.create_staking_limit
    assert amount3 - amount2 == staking_amount


@pytest.mark.P2
def test_POP_021(client_consensus_obj, client_new_node_obj, new_genesis_env):
    """
    Increased lock - time threshold - not in effect
    :param client_consensus_obj:
    :param client_new_node_obj:
    :param new_genesis_env:
    :return:
    """
    genesis = from_dict(data_class=Genesis, data=new_genesis_env.genesis_config)
    genesis.economicModel.staking.unStakeFreezeDuration = 2
    new_genesis_env.set_genesis(genesis.to_dict())
    new_genesis_env.deploy_all()
    address, _ = client_new_node_obj.economic.account.generate_account(client_new_node_obj.node.web3,
                                                                       10 ** 18 * 10000000)
    result = client_new_node_obj.staking.create_staking(0, address, address)
    assert_code(result, 0)

    block = param_governance_verify_before_endblock(client_consensus_obj, "staking", "unStakeFreezeDuration",
                                                    "3", effectiveflag=False)
    wait_block_number(client_new_node_obj.node, block)

    result = client_new_node_obj.staking.withdrew_staking(address)
    assert_code(result, 0)
    amount1 = client_new_node_obj.node.eth.getBalance(address)
    log.info("The wallet balance:{}".format(amount1))
    client_new_node_obj.economic.wait_settlement_blocknum(client_new_node_obj.node, number=1)
    amount2 = client_new_node_obj.node.eth.getBalance(address)
    log.info("The wallet balance:{}".format(amount2))
    client_new_node_obj.economic.wait_settlement_blocknum(client_new_node_obj.node)
    amount3 = client_new_node_obj.node.eth.getBalance(address)
    log.info("The wallet balance:{}".format(amount3))
    staking_amount = client_new_node_obj.economic.create_staking_limit
    assert amount3 - amount2 == staking_amount


@pytest.mark.P2
def test_POP_022(client_consensus_obj, client_new_node_obj, new_genesis_env):
    """
    Reduced lockup threshold - not in effect
    :param client_consensus_obj:
    :param client_new_node_obj:
    :param new_genesis_env:
    :return:
    """
    genesis = from_dict(data_class=Genesis, data=new_genesis_env.genesis_config)
    genesis.economicModel.staking.unStakeFreezeDuration = 3
    new_genesis_env.set_genesis(genesis.to_dict())
    new_genesis_env.deploy_all()
    address, _ = client_new_node_obj.economic.account.generate_account(client_new_node_obj.node.web3,
                                                                       10 ** 18 * 10000000)
    result = client_new_node_obj.staking.create_staking(0, address, address)
    assert_code(result, 0)
    block = param_governance_verify_before_endblock(client_consensus_obj, "staking", "unStakeFreezeDuration",
                                                    "2", effectiveflag=False)
    wait_block_number(client_new_node_obj.node, block)

    result = client_new_node_obj.staking.withdrew_staking(address)
    assert_code(result, 0)
    amount1 = client_new_node_obj.node.eth.getBalance(address)
    log.info("The wallet balance:{}".format(amount1))
    client_new_node_obj.economic.wait_settlement_blocknum(client_new_node_obj.node, number=2)
    amount2 = client_new_node_obj.node.eth.getBalance(address)
    log.info("The wallet balance:{}".format(amount2))
    client_new_node_obj.economic.wait_settlement_blocknum(client_new_node_obj.node)
    amount3 = client_new_node_obj.node.eth.getBalance(address)
    log.info("The wallet balance:{}".format(amount3))
    staking_amount = client_new_node_obj.economic.create_staking_limit
    assert amount3 - amount2 == staking_amount


@pytest.mark.P2
def test_POP_023(client_consensus_obj, client_new_node_obj, new_genesis_env):
    """
    Return pledge before lock time parameter takes effect
    :param client_consensus_obj:
    :param client_new_node_obj:
    :param new_genesis_env:
    :return:
    """
    genesis = from_dict(data_class=Genesis, data=new_genesis_env.genesis_config)
    genesis.economicModel.staking.unStakeFreezeDuration = 2
    new_genesis_env.set_genesis(genesis.to_dict())
    new_genesis_env.deploy_all()
    address, _ = client_new_node_obj.economic.account.generate_account(client_new_node_obj.node.web3,
                                                                       10 ** 18 * 10000000)
    result = client_new_node_obj.staking.create_staking(0, address, address)
    assert_code(result, 0)
    log.info("Next settlement period")
    client_new_node_obj.economic.wait_settlement_blocknum(client_new_node_obj.node)
    result = client_new_node_obj.staking.withdrew_staking(address)
    assert_code(result, 0)
    log.info("Withdraw pledge before parameter takes effect")
    block = param_governance_verify_before_endblock(client_consensus_obj, "staking", "unStakeFreezeDuration",
                                                    "3")
    wait_block_number(client_new_node_obj.node, block)

    amount1 = client_new_node_obj.node.eth.getBalance(address)
    log.info("The wallet balance:{}".format(amount1))
    client_new_node_obj.economic.wait_settlement_blocknum(client_new_node_obj.node)
    amount2 = client_new_node_obj.node.eth.getBalance(address)
    log.info("The wallet balance:{}".format(amount2))
    staking_amount = client_new_node_obj.economic.create_staking_limit
    assert amount2 - amount1 == staking_amount
