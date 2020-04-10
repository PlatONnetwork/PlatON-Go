#!/usr/bin/env python3
# -*- coding: utf-8 -*-
#   @Time    : 2020/2/13 20:40
#   @Author  : PlatON-Developer
#   @Site    : https://github.com/PlatONnetwork/
import time

import math

from tests.lib import get_pledge_list, Decimal
from tests.lib.client import Client
from tests.lib.config import StakingConfig


def test_deploy(global_test_env):
    global_test_env.deploy_all("./deploy/tmp/genesis_0.8.0.json")
    node1 = global_test_env.get_all_nodes()[0]
    node2 = global_test_env.get_all_nodes()[1]
    node1.admin.addPeer("enode://80053b99102f118b99006c436b5e63513d405ba560498b3debc0ea038b0338c01ce7c1a0447ec7b41400f20b7706ed68f3267c226170cf406e066e5bbe3445b2@10.10.8.191:16789")
    node2.admin.addPeer("enode://80053b99102f118b99006c436b5e63513d405ba560498b3debc0ea038b0338c01ce7c1a0447ec7b41400f20b7706ed68f3267c226170cf406e066e5bbe3445b2@10.10.8.191:16789")
    node1.admin.addPeer("enode://0bf7e804ad9c518a72d2ae698d453946109d2fa4547aa2c28288a80868aad1d6042fcdcfd7cb38cebb849846da5204226a113176a93a777133ea4340dff0a274@10.10.8.192:16789")
    node2.admin.addPeer("enode://0bf7e804ad9c518a72d2ae698d453946109d2fa4547aa2c28288a80868aad1d6042fcdcfd7cb38cebb849846da5204226a113176a93a777133ea4340dff0a274@10.10.8.192:16789")
    n = 0
    print(node1.admin.nodeInfo)
    print(node1.debug.economicConfig())
    print(node1.eth.getBalance("0x2e95E3ce0a54951eB9A99152A6d5827872dFB4FD"))

    while n < 30:
        print(node1.block_number, "node1")
        print(node2.block_number, "node2")
        time.sleep(1)
        n += 1

def test_stak(global_test_env):
    node1 = global_test_env.get_all_nodes()[0]
    node2 = global_test_env.get_all_nodes()[1]
    address1, _ = global_test_env.account.generate_account(node1.web3, 2000000000000000000000000)
    print(address1, _)
    client1 = Client(global_test_env, node1, StakingConfig("external_id", "node_name221", "website", "details"))
    client1.staking.create_staking(0, address1,address1)
    address2, _ = global_test_env.account.generate_account(node1.web3, 2000000000000000000000000)
    print(address2, _)
    client2 = Client(global_test_env, node2, StakingConfig("external_id", "node_name222", "website", "details"))
    client2.staking.create_staking(0, address2, address2)

def test_00(global_test_env):
    node1 = global_test_env.get_all_nodes()[0]
    # node2 = global_test_env.get_all_nodes()[1]
    client1 = Client(global_test_env, node1, StakingConfig("external_id", "node_name221", "website", "details"))
    # client2 = Client(global_test_env, node2, StakingConfig("external_id", "node_name222", "website", "details"))
    # result = client1.ppos.getCandidateInfo(client1.node.node_id)
    # print(result)
    node_list = client1.economic.env.node_config_list
    nodeid_list = []
    for i in range(len(node_list)):
        id = node_list[i]['id']
        nodeid_list.append(id)
    print(nodeid_list)

def test_01(global_test_env):
    node_id = '3058ac78b0a05637218a417e562daaca2d640afb3d142ada765650cc0bed892d91d6e8128df0a59397ea051a2d91af5b532866f411811f4fd46de068ad0e168d'
    node1 = global_test_env.find_node_by_node_id(node_id)
    print(node1.node_mark)
    # # node2 = global_test_env.get_all_nodes()[1]
    # client1 = Client(global_test_env, node1, StakingConfig("external_id", "node_name221", "website", "details"))
    # # client2 = Client(global_test_env, node2, StakingConfig("external_id", "node_name222", "website", "details"))
    # tmp_current_block = client1.node.eth.blockNumber
    # last_settlement_block = (math.ceil(tmp_current_block / 10750) - 1) * 10750
    # settlement_block_info = client1.node.eth.getBlock(last_settlement_block)
    # settlement_timestamp = settlement_block_info['timestamp']
    # print(settlement_timestamp)
    # block_info = client1.node.eth.getBlock(1)
    # first_timestamp = block_info['timestamp']
    # print('first_timestamp', first_timestamp)
    # issuing_cycle_timestamp = first_timestamp + (525960 * 60000)
    # remaining_additional_time = issuing_cycle_timestamp - settlement_timestamp
    # result = client1.node.ppos.getAvgPackTime()
    # average_interval = result['Ret']
    # number_of_remaining_blocks = math.ceil(remaining_additional_time / average_interval)
    # remaining_settlement_cycle = math.ceil(number_of_remaining_blocks / 10750)
    # block_proportion = str(50 / 100)
    # verifier_list = get_pledge_list(client1.node.ppos.getVerifierList)
    # verifier_num = len(verifier_list)
    # incentive_pool_amount = client1.node.eth.getBalance('0x1000000000000000000000000000000000000003', last_settlement_block)
    # amount_per_settlement = int(Decimal(str(incentive_pool_amount)) / Decimal(str(remaining_settlement_cycle)))
    # total_block_rewards = int(Decimal(str(amount_per_settlement)) * Decimal(str(block_proportion)))
    # per_block_reward = int(Decimal(str(total_block_rewards)) / Decimal(str(10750)))
    # staking_reward_total = amount_per_settlement - total_block_rewards
    # # staking_reward = int(Decimal(str(staking_reward_total)) / Decimal(str(verifier_num)))
    # print(staking_reward_total, per_block_reward)
    # result = client1.node.ppos.getPackageReward()
    # block_reward = result['Ret']
    # result = client1.node.ppos.getStakingReward()
    # staking_reward = result['Ret']
    # print('system block_reward,staking_reward', block_reward, staking_reward)
    result = node1.debug.getWaitSlashingNodeList()
    print(result)