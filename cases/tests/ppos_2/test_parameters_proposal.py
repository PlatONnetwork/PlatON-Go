# -*- coding: utf-8 -*-
from tests.conftest import param_governance_verify_before_endblock
from tests.lib.utils import *
import pytest
from dacite import from_dict
from common.log import log
from tests.lib import Genesis


@pytest.mark.P2
def test_POP_001_003(client_consensus, client_new_node):
    """
    Increase the threshold of pledge
    :param client_consensus:
    :param get_generate_account:
    :param client_new_node:
    :return:
    """
    client_consensus.economic.env.deploy_all()
    old_amount = client_consensus.economic.create_staking_limit
    new_amount = old_amount + client_consensus.node.web3.toWei(1, "ether")
    block = param_governance_verify_before_endblock(client_consensus, "staking", "stakeThreshold",
                                                    str(new_amount))
    log.info(block)
    address, _ = client_new_node.economic.account.generate_account(client_new_node.node.web3,
                                                                   10 ** 18 * 10000000)
    wait_block_number(client_new_node.node, block)
    result = client_new_node.staking.create_staking(0, address, address, amount=old_amount)
    log.info(result)
    assert_code(result, 301100)
    result = client_new_node.staking.create_staking(0, address, address, amount=new_amount)
    log.info(result)
    assert_code(result, 0)
    verifier_list = get_pledge_list(client_new_node.ppos.getVerifierList)
    log.info(verifier_list)
    client_new_node.economic.wait_settlement_blocknum(client_new_node.node)
    verifier_list = get_pledge_list(client_new_node.ppos.getVerifierList)
    assert client_new_node.node.node_id in verifier_list


@pytest.mark.P2
def test_POP_002(client_consensus, client_new_node, new_genesis_env):
    """
    Minimum pledge reduced pledge threshold
    :param client_consensus:
    :param client_new_node:
    :param new_genesis_env:
    :return:
    """

    old_amount = client_consensus.economic.create_staking_limit + client_consensus.node.web3.toWei(2, "ether")
    new_amount = client_consensus.economic.create_staking_limit + client_consensus.node.web3.toWei(1, "ether")

    genesis = from_dict(data_class=Genesis, data=new_genesis_env.genesis_config)
    genesis.economicModel.staking.stakeThreshold = old_amount
    new_genesis_env.set_genesis(genesis.to_dict())
    new_genesis_env.deploy_all()

    block = param_governance_verify_before_endblock(client_consensus, "staking", "stakeThreshold",
                                                    str(new_amount))
    log.info(block)
    wait_block_number(client_new_node.node, block)
    address, _ = client_new_node.economic.account.generate_account(client_new_node.node.web3,
                                                                   10 ** 18 * 10000000)
    result = client_new_node.staking.create_staking(0, address, address,
                                                    amount=new_amount - client_consensus.node.web3.toWei(1,
                                                                                                         "ether"))
    log.info(result)
    assert_code(result, 301100)
    result = client_new_node.staking.create_staking(0, address, address, amount=new_amount)
    log.info(result)
    assert_code(result, 0)


@pytest.mark.P2
def test_POP_005(client_consensus, client_new_node):
    """
    The amendment of the threshold of pledge shall not take effect
    :param client_consensus:
    :param get_generate_account:
    :param client_new_node:
    :return:
    """
    client_consensus.economic.env.deploy_all()
    old_amount = client_consensus.economic.create_staking_limit
    new_amount = old_amount + client_consensus.node.web3.toWei(1, "ether")

    block = param_governance_verify_before_endblock(client_consensus, "staking", "stakeThreshold",
                                                    str(new_amount), effectiveflag=False)
    log.info(block)
    address, _ = client_new_node.economic.account.generate_account(client_new_node.node.web3,
                                                                   10 ** 18 * 10000000)
    wait_block_number(client_new_node.node, block)
    result = client_new_node.staking.create_staking(0, address, address, amount=old_amount)
    assert_code(result, 0)


@pytest.mark.P2
def test_POP_006(client_consensus, client_new_node):
    """
    (hesitation period) increase - entrustment overweight threshold
    :param client_consensus:
    :param client_new_node:
    :param new_genesis_env:
    :return:
    """
    client_consensus.economic.env.deploy_all()
    old_amount = client_consensus.economic.delegate_limit
    new_amount = old_amount + client_consensus.node.web3.toWei(1, "ether")

    block = param_governance_verify_before_endblock(client_consensus, "staking", "operatingThreshold",
                                                    str(new_amount))
    log.info(block)
    address1, _ = client_new_node.economic.account.generate_account(client_new_node.node.web3,
                                                                    10 ** 18 * 10000000)
    wait_block_number(client_new_node.node, block)
    result = client_new_node.staking.create_staking(0, address1, address1)
    assert_code(result, 0)
    address2, _ = client_new_node.economic.account.generate_account(client_new_node.node.web3,
                                                                    10 ** 18 * 10000000)
    result = client_new_node.delegate.delegate(0, address2, amount=old_amount)
    assert_code(result, 301105)
    result = client_new_node.delegate.delegate(0, address2, amount=new_amount)
    assert_code(result, 0)
    result = client_new_node.staking.increase_staking(0, address1, amount=old_amount)
    assert_code(result, 301104)
    result = client_new_node.staking.increase_staking(0, address1, amount=new_amount)
    assert_code(result, 0)


@pytest.mark.P2
def test_POP_007(client_consensus, client_new_node):
    """
    (lockup period) increase - entrustment overweight threshold
    :param client_consensus:
    :param client_new_node:
    :param new_genesis_env:
    :return:
    """
    client_consensus.economic.env.deploy_all()
    old_amount = client_consensus.economic.delegate_limit
    new_amount = old_amount + client_consensus.node.web3.toWei(1, "ether")

    block = param_governance_verify_before_endblock(client_consensus, "staking", "operatingThreshold",
                                                    str(new_amount))
    log.info(block)
    address1, _ = client_new_node.economic.account.generate_account(client_new_node.node.web3,
                                                                    10 ** 18 * 10000000)
    wait_block_number(client_new_node.node, block)
    result = client_new_node.staking.create_staking(0, address1, address1)
    assert_code(result, 0)
    address2, _ = client_new_node.economic.account.generate_account(client_new_node.node.web3,
                                                                    10 ** 18 * 10000000)
    client_new_node.economic.wait_settlement_blocknum(client_new_node.node)

    result = client_new_node.delegate.delegate(0, address2, amount=old_amount)
    assert_code(result, 301105)
    result = client_new_node.delegate.delegate(0, address2, amount=new_amount)
    assert_code(result, 0)
    result = client_new_node.staking.increase_staking(0, address1, amount=old_amount)
    assert_code(result, 301104)
    result = client_new_node.staking.increase_staking(0, address1, amount=new_amount)
    assert_code(result, 0)


@pytest.mark.P2
def test_POP_008(client_consensus, client_new_node, new_genesis_env):
    """
    (hesitation period) reduce the entrustment overweight threshold - test
    :param client_consensus:
    :param client_new_node:
    :param new_genesis_env:
    :return:
    """

    old_amount = client_consensus.economic.delegate_limit + client_consensus.node.web3.toWei(2, "ether")
    new_amount = client_consensus.economic.delegate_limit + client_consensus.node.web3.toWei(1, "ether")

    genesis = from_dict(data_class=Genesis, data=new_genesis_env.genesis_config)
    genesis.economicModel.staking.operatingThreshold = old_amount
    new_genesis_env.set_genesis(genesis.to_dict())
    new_genesis_env.deploy_all()

    block = param_governance_verify_before_endblock(client_consensus, "staking", "operatingThreshold",
                                                    str(new_amount))
    log.info(block)
    address1, _ = client_new_node.economic.account.generate_account(client_new_node.node.web3,
                                                                    10 ** 18 * 10000000)
    wait_block_number(client_new_node.node, block)
    result = client_new_node.staking.create_staking(0, address1, address1)
    assert_code(result, 0)
    address2, _ = client_new_node.economic.account.generate_account(client_new_node.node.web3,
                                                                    10 ** 18 * 10000000)
    result = client_new_node.delegate.delegate(0, address2,
                                               amount=new_amount - client_consensus.node.web3.toWei(1, "ether"))
    assert_code(result, 301105)
    result = client_new_node.delegate.delegate(0, address2, amount=new_amount)
    assert_code(result, 0)
    result = client_new_node.staking.increase_staking(0, address1,
                                                      amount=new_amount - client_consensus.node.web3.toWei(1, "ether"))
    assert_code(result, 301104)
    result = client_new_node.staking.increase_staking(0, address1, amount=new_amount)
    assert_code(result, 0)


@pytest.mark.P2
def test_POP_009(client_consensus, client_new_node, new_genesis_env):
    """
    (lockup period) reduce the entrustment increase threshold - test
    :param client_consensus:
    :param client_new_node:
    :param new_genesis_env:
    :return:
    """
    old_amount = client_consensus.economic.delegate_limit + client_consensus.node.web3.toWei(2, "ether")
    new_amount = client_consensus.economic.delegate_limit + client_consensus.node.web3.toWei(1, "ether")

    genesis = from_dict(data_class=Genesis, data=new_genesis_env.genesis_config)
    genesis.economicModel.staking.operatingThreshold = old_amount
    new_genesis_env.set_genesis(genesis.to_dict())
    new_genesis_env.deploy_all()

    block = param_governance_verify_before_endblock(client_consensus, "staking", "operatingThreshold",
                                                    str(new_amount))
    log.info(block)
    address1, _ = client_new_node.economic.account.generate_account(client_new_node.node.web3,
                                                                    10 ** 18 * 10000000)
    wait_block_number(client_new_node.node, block)
    result = client_new_node.staking.create_staking(0, address1, address1)
    assert_code(result, 0)
    address2, _ = client_new_node.economic.account.generate_account(client_new_node.node.web3,
                                                                    10 ** 18 * 10000000)
    log.info("The next cycle")
    client_new_node.economic.wait_settlement_blocknum(client_new_node.node)
    result = client_new_node.delegate.delegate(0, address2,
                                               amount=new_amount - client_consensus.node.web3.toWei(1, "ether"))
    assert_code(result, 301105)
    result = client_new_node.delegate.delegate(0, address2, amount=new_amount)
    assert_code(result, 0)
    result = client_new_node.staking.increase_staking(0, address1,
                                                      amount=new_amount - client_consensus.node.web3.toWei(1, "ether"))
    assert_code(result, 301104)
    result = client_new_node.staking.increase_staking(0, address1, amount=new_amount)
    assert_code(result, 0)


@pytest.mark.P2
def test_POP_010_011(client_consensus, client_new_node):
    """
    (hesitation period, lockup period) free amount initiate revocation entrustment
    :param client_consensus:
    :param client_new_node:
    :param new_genesis_env:
    :return:
    """
    client_consensus.economic.env.deploy_all()
    address1, _ = client_new_node.economic.account.generate_account(client_new_node.node.web3,
                                                                    10 ** 18 * 10000000)
    address_delegate_1, _ = client_new_node.economic.account.generate_account(client_new_node.node.web3,
                                                                              10 ** 18 * 10000000)
    result = client_new_node.staking.create_staking(0, address1, address1)
    assert_code(result, 0)
    delegate_amount = client_consensus.economic.delegate_limit + client_consensus.node.web3.toWei(40, "ether")
    log.info("Entrust the sum{}".format(delegate_amount))
    result = client_new_node.delegate.delegate(0, address_delegate_1, amount=delegate_amount)
    assert_code(result, 0)
    amount1_before = client_new_node.node.eth.getBalance(address_delegate_1)
    log.info("The wallet balance:{}".format(amount1_before))
    parameters_amount = client_consensus.economic.delegate_limit + client_consensus.node.web3.toWei(10, "ether")

    block = param_governance_verify_before_endblock(client_consensus, "staking", "operatingThreshold",
                                                    str(parameters_amount))
    log.info("The value of the proposal parameters{}".format(parameters_amount))
    wait_block_number(client_new_node.node, block)
    log.info("The delegate is initiated after the parameter takes effect")
    address_delegate_2, _ = client_new_node.economic.account.generate_account(client_new_node.node.web3,
                                                                              10 ** 18 * 10000000)
    result = client_new_node.delegate.delegate(0, address_delegate_2, amount=delegate_amount)
    assert_code(result, 0)
    amount2_before = client_new_node.node.eth.getBalance(address_delegate_1)
    log.info("The wallet balance:{}".format(amount2_before))

    msg = client_new_node.ppos.getCandidateInfo(client_new_node.node.node_id)
    staking_blocknum = msg["Ret"]["StakingBlockNum"]
    withdrew_delegate_amount = delegate_amount - parameters_amount + client_consensus.node.web3.toWei(1, "ether")

    result = client_new_node.delegate.withdrew_delegate(staking_blocknum, address_delegate_1,
                                                        amount=withdrew_delegate_amount)
    assert_code(result, 0)
    amount1_after = client_new_node.node.eth.getBalance(address_delegate_1)
    log.info("The wallet balance:{}".format(amount1_after))
    assert amount1_before - amount1_after < client_new_node.node.web3.toWei(1, "ether")

    result = client_new_node.delegate.withdrew_delegate(staking_blocknum, address_delegate_2,
                                                        amount=withdrew_delegate_amount)
    assert_code(result, 0)
    amount2_after = client_new_node.node.eth.getBalance(address_delegate_1)
    log.info("The wallet balance:{}".format(amount2_after))
    assert amount1_before - amount1_after < client_new_node.node.web3.toWei(1, "ether")


@pytest.mark.P2
def test_POP_012(client_consensus, client_new_node):
    """
    :param client_consensus:
    :param client_new_node:
    :param new_genesis_env:
    :return:
    """

    client_consensus.economic.env.deploy_all()
    address, _ = client_new_node.economic.account.generate_account(client_new_node.node.web3,
                                                                   10 ** 18 * 10000000)
    result = client_new_node.staking.create_staking(0, address, address)

    assert_code(result, 0)
    address_delegate, _ = client_new_node.economic.account.generate_account(client_new_node.node.web3,
                                                                            10 ** 18 * 10000000)
    delegate_amount = client_consensus.economic.delegate_limit + client_consensus.node.web3.toWei(20, "ether")

    log.info("Own funds to initiate the commission")
    result = client_new_node.delegate.delegate(0, address_delegate, amount=delegate_amount)
    assert_code(result, 0)

    parameters_amount = client_consensus.economic.delegate_limit + \
        client_consensus.node.web3.toWei(10, "ether")
    block = param_governance_verify_before_endblock(client_consensus, "staking", "operatingThreshold",
                                                    str(parameters_amount))
    wait_block_number(client_new_node.node, block)

    plan = [{'Epoch': 1, 'Amount': delegate_amount}]
    result = client_new_node.restricting.createRestrictingPlan(address_delegate, plan, address_delegate)
    assert_code(result, 0)

    log.info("Fund of lockup is initiated and entrusted")
    result = client_new_node.delegate.delegate(1, address_delegate, amount=delegate_amount)
    assert_code(result, 0)
    amount_before = client_new_node.node.eth.getBalance(address_delegate)
    log.info("The wallet balance:{}".format(amount_before))

    msg = client_new_node.ppos.getCandidateInfo(client_new_node.node.node_id)
    staking_blocknum = msg["Ret"]["StakingBlockNum"]
    withdrew_delegate = delegate_amount * 2 - parameters_amount + \
        client_consensus.node.web3.toWei(1, "ether")
    result = client_new_node.delegate.withdrew_delegate(staking_blocknum, address_delegate,
                                                        amount=withdrew_delegate)
    assert_code(result, 0)
    amount1_after = client_new_node.node.eth.getBalance(address_delegate)
    log.info("The wallet balance:{}".format(amount1_after))
    amount_dill = amount1_after - amount_before
    assert delegate_amount - amount_dill < client_new_node.node.web3.toWei(1, "ether")
    client_new_node.economic.wait_settlement_blocknum(client_new_node.node)

    amount1_last = client_new_node.node.eth.getBalance(address_delegate)
    log.info("The wallet balance:{}".format(amount1_last))
    assert amount1_last - amount1_after == delegate_amount
    assert delegate_amount * 2 - (amount1_last - amount_before) < client_new_node.node.web3.toWei(1, "ether")


@pytest.mark.P2
def test_POP_013(client_consensus, client_new_node, new_genesis_env):
    """
    :param client_consensus:
    :param client_new_node:
    :param new_genesis_env:
    :return:
    """
    genesis = from_dict(data_class=Genesis, data=new_genesis_env.genesis_config)
    old_amount = client_consensus.economic.delegate_limit * 2
    new_amount = client_consensus.economic.delegate_limit
    genesis.economicModel.staking.operatingThreshold = old_amount
    new_genesis_env.set_genesis(genesis.to_dict())
    new_genesis_env.deploy_all()

    block = param_governance_verify_before_endblock(client_consensus, "staking", "operatingThreshold",
                                                    str(new_amount), effectiveflag=False)
    log.info(block)
    address1, _ = client_new_node.economic.account.generate_account(client_new_node.node.web3,
                                                                    10 ** 18 * 10000000)
    wait_block_number(client_new_node.node, block)
    result = client_new_node.staking.create_staking(0, address1, address1)
    assert_code(result, 0)

    address2, _ = client_new_node.economic.account.generate_account(client_new_node.node.web3,
                                                                    10 ** 18 * 10000000)
    result = client_new_node.delegate.delegate(0, address2, amount=new_amount)
    log.info(result)
    assert_code(result, 301105)
    result = client_new_node.delegate.delegate(0, address2, amount=old_amount)
    assert_code(result, 0)
    result = client_new_node.staking.increase_staking(0, address1, amount=new_amount)
    assert_code(result, 301104)
    result = client_new_node.staking.increase_staking(0, address1, amount=old_amount)
    assert_code(result, 0)


@pytest.mark.P2
def test_POP_014(client_consensus, client_new_node):
    """
    Parameter not in effect - initiate redemption
    :param client_consensus:
    :param client_new_node:
    :param new_genesis_env:
    :return:
    """
    client_consensus.economic.env.deploy_all()
    address1, _ = client_new_node.economic.account.generate_account(client_new_node.node.web3,
                                                                    10 ** 18 * 10000000)
    address_delegate_1, _ = client_new_node.economic.account.generate_account(client_new_node.node.web3,
                                                                              10 ** 18 * 10000000)
    result = client_new_node.staking.create_staking(0, address1, address1)
    assert_code(result, 0)
    delegate_amount = client_consensus.economic.delegate_limit + client_consensus.node.web3.toWei(20, "ether")
    log.info("Amount entrusted{}".format(delegate_amount))

    result = client_new_node.delegate.delegate(0, address_delegate_1, amount=delegate_amount)
    assert_code(result, 0)
    amount1_before = client_new_node.node.eth.getBalance(address_delegate_1)
    log.info("The wallet balance:{}".format(amount1_before))
    param_amount = client_consensus.economic.delegate_limit + client_consensus.node.web3.toWei(10, "ether")

    block = param_governance_verify_before_endblock(client_consensus, "staking", "operatingThreshold",
                                                    str(param_amount), effectiveflag=False)
    wait_block_number(client_new_node.node, block)
    log.info("The delegate is initiated after the parameter takes effect")
    address_delegate_2, _ = client_new_node.economic.account.generate_account(client_new_node.node.web3,
                                                                              10 ** 18 * 10000000)
    result = client_new_node.delegate.delegate(0, address_delegate_2, amount=delegate_amount)
    assert_code(result, 0)
    amount2_before = client_new_node.node.eth.getBalance(address_delegate_1)
    log.info("The wallet balance:{}".format(amount2_before))

    msg = client_new_node.ppos.getCandidateInfo(client_new_node.node.node_id)
    staking_blocknum = msg["Ret"]["StakingBlockNum"]
    withdrew_delegate = delegate_amount - param_amount + client_consensus.node.web3.toWei(1, "ether")
    result = client_new_node.delegate.withdrew_delegate(staking_blocknum, address_delegate_1,
                                                        amount=withdrew_delegate)
    assert_code(result, 0)
    amount1_after = client_new_node.node.eth.getBalance(address_delegate_1)
    log.info("The wallet balance:{}".format(amount1_after))
    amount1_dill = amount1_after - amount1_before
    assert withdrew_delegate - amount1_dill < client_new_node.node.web3.toWei(1, "ether")

    result = client_new_node.delegate.withdrew_delegate(staking_blocknum, address_delegate_2,
                                                        amount=withdrew_delegate)
    assert_code(result, 0)
    amount2_after = client_new_node.node.eth.getBalance(address_delegate_2)
    log.info("The wallet balance:{}".format(amount2_after))
    amount1_dill = amount2_after - amount2_before
    assert withdrew_delegate - amount1_dill < client_new_node.node.web3.toWei(1, "ether")


@pytest.mark.P2
def test_POP_015(client_consensus, clients_noconsensus, new_genesis_env):
    """
    Increase the number of alternative nodes
    :param client_consensus:
    :param client_new_node:
    :param new_genesis_env:
    :return:
    """

    genesis = from_dict(data_class=Genesis, data=new_genesis_env.genesis_config)
    genesis.economicModel.staking.maxValidators = 5
    new_genesis_env.set_genesis(genesis.to_dict())
    new_genesis_env.deploy_all()
    client1 = clients_noconsensus[0]
    client2 = clients_noconsensus[1]
    address, _ = client1.economic.account.generate_account(client1.node.web3,
                                                           10 ** 18 * 10000000)
    staking_amount = client1.economic.create_staking_limit
    result = client1.staking.create_staking(0, address, address, amount=staking_amount * 2)
    assert_code(result, 0)
    param = 6
    block = param_governance_verify_before_endblock(client_consensus, "staking", "maxValidators",
                                                    str(param))
    wait_block_number(client2.node, block)
    address1, _ = client2.economic.account.generate_account(client2.node.web3,
                                                            10 ** 18 * 10000000)
    result = client2.staking.create_staking(0, address1, address1, amount=staking_amount * 3)
    assert_code(result, 0)

    client2.economic.wait_settlement_blocknum(client2.node)
    getVerifierList = get_pledge_list(client2.node.ppos.getVerifierList)
    log.info(getVerifierList)
    assert len(getVerifierList) == 6
    node_id_1 = client1.node.node_id
    node_id_2 = client2.node.node_id
    assert node_id_1 in getVerifierList
    assert node_id_2 in getVerifierList


@pytest.mark.P2
def test_POP_016(client_consensus, clients_noconsensus, new_genesis_env):
    """
    Reduce the number of alternative nodes
    :param client_consensus:
    :param client_new_node:
    :param new_genesis_env:
    :return:
    """
    genesis = from_dict(data_class=Genesis, data=new_genesis_env.genesis_config)
    genesis.economicModel.staking.maxValidators = 6
    new_genesis_env.set_genesis(genesis.to_dict())
    new_genesis_env.deploy_all()
    client1 = clients_noconsensus[0]
    client2 = clients_noconsensus[1]
    param = 5
    block = param_governance_verify_before_endblock(client_consensus, "staking", "maxValidators",
                                                    str(param))
    wait_block_number(client2.node, block)
    address, _ = client1.economic.account.generate_account(client1.node.web3,
                                                           10 ** 18 * 10000000)
    staking_amount = client1.economic.create_staking_limit
    result = client1.staking.create_staking(0, address, address, amount=staking_amount * 2)
    assert_code(result, 0)
    address1, _ = client2.economic.account.generate_account(client2.node.web3,
                                                            10 ** 18 * 10000000)
    result = client2.staking.create_staking(0, address1, address1, amount=staking_amount * 3)
    assert_code(result, 0)

    client2.economic.wait_settlement_blocknum(client2.node)
    getVerifierList = get_pledge_list(client2.node.ppos.getVerifierList)
    log.info(getVerifierList)
    assert len(getVerifierList) == 5


@pytest.mark.P2
def test_POP_017(client_consensus, clients_noconsensus, new_genesis_env):
    """
    Increase the number of node candidates - not active
    :param client_consensus:
    :param client_new_node:
    :param new_genesis_env:
    :return:
    """
    genesis = from_dict(data_class=Genesis, data=new_genesis_env.genesis_config)
    genesis.economicModel.staking.maxValidators = 5
    new_genesis_env.set_genesis(genesis.to_dict())
    new_genesis_env.deploy_all()
    param = 6
    block = param_governance_verify_before_endblock(client_consensus, "staking", "maxValidators",
                                                    str(param), effectiveflag=False)
    client1 = clients_noconsensus[0]
    client2 = clients_noconsensus[1]
    wait_block_number(client2.node, block)
    address, _ = client1.economic.account.generate_account(client1.node.web3,
                                                           10 ** 18 * 10000000)
    staking_amount = client1.economic.create_staking_limit
    result = client1.staking.create_staking(0, address, address, amount=staking_amount * 2)
    assert_code(result, 0)
    address1, _ = client2.economic.account.generate_account(client2.node.web3,
                                                            10 ** 18 * 10000000)
    result = client2.staking.create_staking(0, address1, address1, amount=staking_amount * 3)
    assert_code(result, 0)

    client2.economic.wait_settlement_blocknum(client2.node)
    getVerifierList = get_pledge_list(client2.node.ppos.getVerifierList)
    log.info(getVerifierList)
    assert len(getVerifierList) == 5


@pytest.mark.P2
def test_POP_018(client_consensus, clients_noconsensus, new_genesis_env):
    """
    Reduce the number of node candidates - not in effect
    :param client_consensus:
    :param client_new_node:
    :param new_genesis_env:
    :return:
    """
    genesis = from_dict(data_class=Genesis, data=new_genesis_env.genesis_config)
    genesis.economicModel.staking.maxValidators = 6
    new_genesis_env.set_genesis(genesis.to_dict())
    new_genesis_env.deploy_all()
    client1 = clients_noconsensus[0]
    client2 = clients_noconsensus[1]

    address, _ = client1.economic.account.generate_account(client1.node.web3,
                                                           10 ** 18 * 10000000)
    staking_amount = client1.economic.create_staking_limit

    result = client1.staking.create_staking(0, address, address, amount=staking_amount * 2)
    assert_code(result, 0)
    param = 5
    block = param_governance_verify_before_endblock(client_consensus, "staking", "maxValidators",
                                                    str(param), effectiveflag=False)
    wait_block_number(client2.node, block)
    address1, _ = client2.economic.account.generate_account(client2.node.web3,
                                                            10 ** 18 * 10000000)
    result = client2.staking.create_staking(0, address1, address1, amount=staking_amount * 3)
    assert_code(result, 0)

    client2.economic.wait_settlement_blocknum(client2.node)
    getVerifierList = get_pledge_list(client2.node.ppos.getVerifierList)
    log.info(getVerifierList)
    assert len(getVerifierList) == 6
    node_id_1 = client1.node.node_id
    node_id_2 = client2.node.node_id
    assert node_id_1 in getVerifierList
    assert node_id_2 in getVerifierList


@pytest.mark.P2
def test_POP_019(client_consensus, client_new_node, new_genesis_env):
    """
    Increased lock - up threshold
    :param client_consensus:
    :param client_new_node:
    :param new_genesis_env:
    :return:
    """
    new_genesis_env.deploy_all()
    address, _ = client_new_node.economic.account.generate_account(client_new_node.node.web3,
                                                                   10 ** 18 * 10000000)
    result = client_new_node.staking.create_staking(0, address, address)
    assert_code(result, 0)
    param = 3
    block = param_governance_verify_before_endblock(client_consensus, "staking", "unStakeFreezeDuration",
                                                    str(param))
    wait_block_number(client_new_node.node, block)

    result = client_new_node.staking.withdrew_staking(address)
    assert_code(result, 0)
    amount1 = client_new_node.node.eth.getBalance(address)
    log.info("The wallet balance:{}".format(amount1))
    client_new_node.economic.wait_settlement_blocknum(client_new_node.node, number=2)
    amount2 = client_new_node.node.eth.getBalance(address)
    log.info("The wallet balance:{}".format(amount2))
    client_new_node.economic.wait_settlement_blocknum(client_new_node.node)
    amount3 = client_new_node.node.eth.getBalance(address)
    log.info("The wallet balance:{}".format(amount3))
    staking_amount = client_new_node.economic.create_staking_limit
    assert amount3 - amount2 == staking_amount


@pytest.mark.P2
def test_POP_020(client_consensus, client_new_node, new_genesis_env):
    """
    Reduce lock - up threshold
    :param client_consensus:
    :param client_new_node:
    :param new_genesis_env:
    :return:
    """
    genesis = from_dict(data_class=Genesis, data=new_genesis_env.genesis_config)
    unStakeFreezeDuration = 3
    genesis.economicModel.staking.unStakeFreezeDuration = unStakeFreezeDuration
    new_genesis_env.set_genesis(genesis.to_dict())
    new_genesis_env.deploy_all()
    address, _ = client_new_node.economic.account.generate_account(client_new_node.node.web3,
                                                                   10 ** 18 * 10000000)
    result = client_new_node.staking.create_staking(0, address, address)
    assert_code(result, 0)
    param = 2
    block = param_governance_verify_before_endblock(client_consensus, "staking", "unStakeFreezeDuration",
                                                    str(param))
    wait_block_number(client_new_node.node, block)

    result = client_new_node.staking.withdrew_staking(address)
    assert_code(result, 0)
    amount1 = client_new_node.node.eth.getBalance(address)
    log.info("The wallet balance:{}".format(amount1))
    client_new_node.economic.wait_settlement_blocknum(client_new_node.node, number=1)
    amount2 = client_new_node.node.eth.getBalance(address)
    log.info("The wallet balance:{}".format(amount2))
    client_new_node.economic.wait_settlement_blocknum(client_new_node.node)
    amount3 = client_new_node.node.eth.getBalance(address)
    log.info("The wallet balance:{}".format(amount3))
    staking_amount = client_new_node.economic.create_staking_limit
    assert amount3 - amount2 == staking_amount


@pytest.mark.P2
def test_POP_021(client_consensus, client_new_node, new_genesis_env):
    """
    Increased lock - time threshold - not in effect
    :param client_consensus:
    :param client_new_node:
    :param new_genesis_env:
    :return:
    """
    genesis = from_dict(data_class=Genesis, data=new_genesis_env.genesis_config)
    unStakeFreezeDuration = 2
    genesis.economicModel.staking.unStakeFreezeDuration = unStakeFreezeDuration
    new_genesis_env.set_genesis(genesis.to_dict())
    new_genesis_env.deploy_all()
    address, _ = client_new_node.economic.account.generate_account(client_new_node.node.web3,
                                                                   10 ** 18 * 10000000)
    result = client_new_node.staking.create_staking(0, address, address)
    assert_code(result, 0)

    param = 3
    block = param_governance_verify_before_endblock(client_consensus, "staking", "unStakeFreezeDuration",
                                                    str(param), effectiveflag=False)
    wait_block_number(client_new_node.node, block)

    result = client_new_node.staking.withdrew_staking(address)
    assert_code(result, 0)
    amount1 = client_new_node.node.eth.getBalance(address)
    log.info("The wallet balance:{}".format(amount1))
    client_new_node.economic.wait_settlement_blocknum(client_new_node.node, number=1)
    amount2 = client_new_node.node.eth.getBalance(address)
    log.info("The wallet balance:{}".format(amount2))
    client_new_node.economic.wait_settlement_blocknum(client_new_node.node)
    amount3 = client_new_node.node.eth.getBalance(address)
    log.info("The wallet balance:{}".format(amount3))
    staking_amount = client_new_node.economic.create_staking_limit
    assert amount3 - amount2 == staking_amount


@pytest.mark.P2
def test_POP_022(client_consensus, client_new_node, new_genesis_env):
    """
    Reduced lockup threshold - not in effect
    :param client_consensus:
    :param client_new_node:
    :param new_genesis_env:
    :return:
    """
    genesis = from_dict(data_class=Genesis, data=new_genesis_env.genesis_config)
    unStakeFreezeDuration = 3
    genesis.economicModel.staking.unStakeFreezeDuration = unStakeFreezeDuration
    new_genesis_env.set_genesis(genesis.to_dict())
    new_genesis_env.deploy_all()
    address, _ = client_new_node.economic.account.generate_account(client_new_node.node.web3,
                                                                   10 ** 18 * 10000000)
    result = client_new_node.staking.create_staking(0, address, address)
    assert_code(result, 0)
    param = 2
    block = param_governance_verify_before_endblock(client_consensus, "staking", "unStakeFreezeDuration",
                                                    str(param), effectiveflag=False)
    wait_block_number(client_new_node.node, block)

    result = client_new_node.staking.withdrew_staking(address)
    assert_code(result, 0)
    amount1 = client_new_node.node.eth.getBalance(address)
    log.info("The wallet balance:{}".format(amount1))
    client_new_node.economic.wait_settlement_blocknum(client_new_node.node, number=2)
    amount2 = client_new_node.node.eth.getBalance(address)
    log.info("The wallet balance:{}".format(amount2))
    client_new_node.economic.wait_settlement_blocknum(client_new_node.node)
    amount3 = client_new_node.node.eth.getBalance(address)
    log.info("The wallet balance:{}".format(amount3))
    staking_amount = client_new_node.economic.create_staking_limit
    assert amount3 - amount2 == staking_amount


@pytest.mark.P2
def test_POP_023(client_consensus, client_new_node, new_genesis_env):
    """
    Return pledge before lock time parameter takes effect
    :param client_consensus:
    :param client_new_node:
    :param new_genesis_env:
    :return:
    """
    genesis = from_dict(data_class=Genesis, data=new_genesis_env.genesis_config)
    unStakeFreezeDuration = 2
    genesis.economicModel.staking.unStakeFreezeDuration = unStakeFreezeDuration
    new_genesis_env.set_genesis(genesis.to_dict())
    new_genesis_env.deploy_all()
    address, _ = client_new_node.economic.account.generate_account(client_new_node.node.web3,
                                                                   10 ** 18 * 10000000)
    result = client_new_node.staking.create_staking(0, address, address)
    assert_code(result, 0)
    log.info("Next settlement period")
    client_new_node.economic.wait_settlement_blocknum(client_new_node.node)
    result = client_new_node.staking.withdrew_staking(address)
    assert_code(result, 0)
    log.info("Withdraw pledge before parameter takes effect")
    param = 3
    block = param_governance_verify_before_endblock(client_consensus, "staking", "unStakeFreezeDuration",
                                                    str(param))
    wait_block_number(client_new_node.node, block)

    amount1 = client_new_node.node.eth.getBalance(address)
    log.info("The wallet balance:{}".format(amount1))
    client_new_node.economic.wait_settlement_blocknum(client_new_node.node)
    amount2 = client_new_node.node.eth.getBalance(address)
    log.info("The wallet balance:{}".format(amount2))
    staking_amount = client_new_node.economic.create_staking_limit
    assert amount2 - amount1 == staking_amount
