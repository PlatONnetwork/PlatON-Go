#!/usr/bin/env python3
# -*- coding: utf-8 -*-
#   @Time    : 2020/1/2 10:41
#   @Author  : PlatON-Developer
#   @Site    : https://github.com/PlatONnetwork/
import os
import pytest
from decimal import Decimal
from tests.lib.genesis import to_genesis
from tests.lib.client import Client
from tests.lib.utils import assert_code


def calculate(big_int, mul):
    return int(Decimal(str(big_int)) * Decimal(mul))


@pytest.fixture()
def staking_node_client(client_new_node):
    reward = 1000
    amount = calculate(client_new_node.economic.create_staking_limit, 5)
    staking_amount = calculate(client_new_node.economic.create_staking_limit, 2)
    staking_address, _ = client_new_node.economic.account.generate_account(client_new_node.node.web3, amount)
    delegate_address, _ = client_new_node.economic.account.generate_account(client_new_node.node.web3,
                                                                            client_new_node.economic.add_staking_limit * 2)
    client_new_node.staking.create_staking(0, staking_address, staking_address, amount=staking_amount, reward_per=reward)
    setattr(client_new_node, "staking_address", staking_address)
    setattr(client_new_node, "delegate_address", delegate_address)
    setattr(client_new_node, "amount", amount)
    setattr(client_new_node, "staking_amount", staking_amount)
    setattr(client_new_node, "reward", reward)
    yield client_new_node


@pytest.fixture()
def delegate_node_client(client_new_node):
    reward = 1000
    amount = calculate(client_new_node.economic.create_staking_limit, 5)
    staking_amount = calculate(client_new_node.economic.create_staking_limit, 2)
    staking_address, _ = client_new_node.economic.account.generate_account(client_new_node.node.web3, amount)
    delegate_address, _ = client_new_node.economic.account.generate_account(client_new_node.node.web3,
                                                                            client_new_node.economic.add_staking_limit * 2)
    client_new_node.staking.create_staking(0, staking_address, staking_address, amount=staking_amount, reward_per=reward)
    client_new_node.delegate.delegate(0, delegate_address)
    setattr(client_new_node, "staking_address", staking_address)
    setattr(client_new_node, "delegate_address", delegate_address)
    setattr(client_new_node, "amount", amount)
    setattr(client_new_node, "staking_amount", staking_amount)
    setattr(client_new_node, "reward", reward)
    yield client_new_node


def test_DG_TR_001(client_consensus):
    reward = 1000
    economic = client_consensus.economic
    node = client_consensus.node
    candidate_info = client_consensus.ppos.getCandidateInfo(node.node_id)
    assert candidate_info["RewardPer"] == 0
    result = client_consensus.staking.edit_candidate(economic.cfg.REMAIN_ACCOUNT_ADDRESS, economic.cfg.INCENTIVEPOOL_ADDRESS, reward_per=reward)
    assert_code(result, 0)
    candidate_info = client_consensus.ppos.getCandidateInfo(node.node_id)
    assert candidate_info["RewardPer"] == reward
    edit_after_balance = node.eth.getBalance(economic.cfg.INCENTIVEPOOL_ADDRESS)
    client_consensus.economic.wait_settlement_blocknum(node)
    candidate_info = client_consensus.ppos.getCandidateInfo(node.node_id)
    assert candidate_info["DelegateRewardTotal"] == 0
    assert node.eth.getBalance(economic.cfg.INCENTIVEPOOL_ADDRESS) - edit_after_balance == 0


def test_DG_TR_002(client_consensus):
    reward = 1000
    economic = client_consensus.economic
    node = client_consensus.node
    candidate_info = client_consensus.ppos.getCandidateInfo(node.node_id)
    assert candidate_info["RewardPer"] == 0
    result = client_consensus.staking.edit_candidate(economic.cfg.REMAIN_ACCOUNT_ADDRESS, economic.cfg.DEVELOPER_FOUNDATAION_ADDRESS,
                                                     reward_per=reward)
    assert_code(result, 0)
    candidate_info = client_consensus.ppos.getCandidateInfo(node.node_id)
    assert candidate_info["RewardPer"] == reward
    assert candidate_info["BenefitAddress"] == economic.cfg.INCENTIVEPOOL_ADDRESS
    economic.wait_settlement_blocknum(node)
    candidate_info = client_consensus.ppos.getCandidateInfo(node.node_id)
    assert candidate_info["DelegateRewardTotal"] == 0


def test_DG_TR_003(staking_node_client):
    economic = staking_node_client.economic
    node = staking_node_client.node
    result = staking_node_client.delegate.delegate(0, staking_node_client.delegate_address)
    assert_code(result, 0)
    reward = 2000
    result = staking_node_client.staking.edit_candidate(economic.cfg.REMAIN_ACCOUNT_ADDRESS,
                                                        economic.cfg.INCENTIVEPOOL_ADDRESS, reward_per=reward)
    assert_code(result, 0)
    candidate_info = staking_node_client.ppos.getCandidateInfo(node.node_id)
    assert candidate_info["RewardPer"] == reward
    economic.wait_settlement_blocknum(node)
    candidate_info = staking_node_client.ppos.getCandidateInfo(node.node_id)
    total_reward = staking_node_client.ppos.getStakingReward() + staking_node_client.ppos.getPackageReward()
    assert candidate_info["DelegateRewardTotal"] == total_reward * reward / 10000


def test_DG_TR_004(staking_node_client):
    node = staking_node_client.node
    candidate_info = staking_node_client.ppos.getCandidateInfo(node.node_id)
    assert candidate_info["DelegateRewardTotal"] == 0


def test_DG_TR_005(global_test_env, reset_environment, staking_cfg):
    genesis_cfg = global_test_env.genesis_config
    genesis = to_genesis(genesis_cfg)
    genesis.economicModel.staking.maxValidators = 4
    genesis_file = os.path.join(global_test_env.cfg.env_tmp, "dg_tr_005_genesis.json")
    genesis.to_file(genesis_file)
    global_test_env.deploy_all(genesis_file)
    node = global_test_env.get_a_normal_node()
    client = Client(global_test_env, node, staking_cfg)
    amount = calculate(client.economic.create_staking_limit, 5)
    staking_amount = calculate(client.economic.create_staking_limit, 2)
    staking_address, _ = client.economic.account.generate_account(client.node.web3, amount)
    delegate_address, _ = client.economic.account.generate_account(client.node.web3, client.economic.add_staking_limit * 2)
    result = client.staking.create_staking(0, staking_address, staking_address, amount=staking_amount, reward_per=1000)
    assert_code(result, 0)
    client.economic.wait_settlement_blocknum(node)
    candidate_info = client.ppos.getCandidateInfo(node.node_id)
    assert candidate_info["DelegateRewardTotal"] == 0


def test_DG_TR_006(staking_node_client):
    pass


def test_DG_TR_007():
    pass


def test_DG_TR_008():
    pass


def test_DG_TR_009(delegate_node_client):
    node = delegate_node_client.node
    economic = delegate_node_client.economic
    start_balance = node.eth.getBalance(delegate_node_client.staking_address)
    economic.wait_settlement_blocknum(node)
    total = delegate_node_client.ppos.getStakingReward() + delegate_node_client.ppos.getPackageReward()
    candidate_info = delegate_node_client.ppos.getCandidateInfo(node.node_id)
    delegate_reward_total = total * delegate_node_client.reward / 10000
    assert candidate_info["DelegateRewardTotal"] == delegate_reward_total
    end_balance = node.eth.getBalance(delegate_node_client.staking_address)
    assert end_balance - start_balance == total - delegate_reward_total


def test_DG_TR_010():
    pass


def test_DG_TR_011():
    pass


def test_DG_TR_012():
    pass


def test_DG_TR_013():
    pass


def test_DG_TR_014():
    pass


def test_DG_TR_015():
    pass


def test_DG_TR_016():
    pass


def test_DG_TR_017():
    pass


def test_DG_TR_018():
    pass


def test_DG_TR_019():
    pass


def test_DG_TR_020():
    pass


def test_DG_TR_021():
    pass


def test_DG_TR_022():
    pass


def test_DG_TR_023():
    pass


def test_DG_TR_024():
    pass


def test_DG_TR_025():
    pass


def test_DG_TR_026():
    pass


def test_DG_TR_027():
    pass


def test_DG_TR_028():
    pass


def test_DG_TR_029():
    pass


def test_DG_TR_030():
    pass
