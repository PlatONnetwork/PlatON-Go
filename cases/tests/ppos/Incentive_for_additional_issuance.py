import math
import time

import pytest
from dacite import from_dict
from common.log import log
from client_sdk_python import Web3
from decimal import Decimal
from tests.conftest import get_client_consensus
from tests.lib import EconomicConfig, Genesis, assert_code, von_amount, Client


@pytest.mark.P1
def test_AL_FI_001_to_003(new_genesis_env, staking_cfg):
    """
    AL_FI_001:查看每年释放补贴激励池变化
    AL_FI_002:查看每年固定增发变化
    AL_FI_003:第十年固定增发token分配
    :param new_genesis_env:
    :return:
    """
    # Initialization genesis file Initial amount
    node_count = len(new_genesis_env.consensus_node_list)
    default_pledge_amount = Web3.toWei(node_count * 1500000, 'ether')
    community_amount = default_pledge_amount + 259096239000000000000000000 + 62215742000000000000000000
    genesis = from_dict(data_class=Genesis, data=new_genesis_env.genesis_config)
    genesis.economicModel.innerAcc.cdfBalance = community_amount
    surplus_amount = str(EconomicConfig.TOKEN_TOTAL - community_amount - 200000000000000000000000000)
    genesis.alloc = {
        "lax1zqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqrzpqayr": {
            "balance": "200000000000000000000000000"
        },
        "lax196278ns22j23awdfj9f2d4vz0pedld8au6xelj": {
            "balance": surplus_amount
        }
    }
    new_file = new_genesis_env.cfg.env_tmp + "/genesis_0.13.0.json"
    genesis.to_file(new_file)
    new_genesis_env.deploy_all(new_file)

    client = get_client_consensus(new_genesis_env, staking_cfg)
    economic = client.economic
    node = client.node
    # Query the initial amount of incentive pool
    current_incentive_pool = 262215742000000000000000000
    # Query the initial amount of a warehouse lock plan
    init_foundationlockup = 259096239000000000000000000
    # Issued token amount
    init_token = 10000000000000000000000000000
    # Query developer foundation initial amount
    developer_foundation = 0
    # Query the initial amount of the foundation
    foundation_balance = 0
    # Additional amount
    total_amount_of_issuance = 0
    remaining_settlement_cycle = 0
    end_cycle_timestamp = None
    # Annual issuance
    for i in range(10):
        if i == 0:
            incentive_pool = current_incentive_pool
            log.info("Amount of initial incentive pool： {}".format(incentive_pool))
            foundation_lock_up = init_foundationlockup
            log.info("Initial Lockup Plan Amount： {}".format(foundation_lock_up))
            total_amount_of_issuance = int(init_token + Decimal(str(init_token)) / Decimal(str(40)))
            log.info("Current year Total amount of issuance： {}".format(total_amount_of_issuance))
            # Query the current annual incentive pool amount
            current_annual_incentive_pool_amount = node.eth.getBalance(EconomicConfig.INCENTIVEPOOL_ADDRESS, 0)
            # Query current annual developer foundation amount
            DEVELOPER_FOUNDATAION_ADDRESS = node.web3.toChecksumAddress(EconomicConfig.DEVELOPER_FOUNDATAION_ADDRESS)
            current_annual_developer_foundation_amount = node.eth.getBalance(DEVELOPER_FOUNDATAION_ADDRESS, 0)
            # Query current annual fund amount
            FOUNDATION_ADDRESS = node.web3.toChecksumAddress(EconomicConfig.FOUNDATION_ADDRESS)
            current_annual_foundation_amount = node.eth.getBalance(FOUNDATION_ADDRESS, 0)
            log.info(
                "{} Year Incentive Pool Address: {} Balance: {}".format(i + 1, EconomicConfig.INCENTIVEPOOL_ADDRESS,
                                                                        incentive_pool))
            log.info('{} Year Foundation Address: {} Balance: {}'.format(i + 1, EconomicConfig.FOUNDATION_ADDRESS,
                                                                         foundation_balance))
            log.info("{} Year Developer Foundation Address:{} Balance:{}".format(i + 1,
                                                                                 EconomicConfig.DEVELOPER_FOUNDATAION_ADDRESS,
                                                                                 developer_foundation))
            log.info("{} Year Foundation Locking Address: {} Balance: {}".format(i + 1,
                                                                                 EconomicConfig.FOUNDATION_LOCKUP_ADDRESS,
                                                                                 foundation_lock_up))
            assert current_annual_incentive_pool_amount == incentive_pool, "{} Year Incentive Pool Address: {} Balance: {}".format(
                i + 1, EconomicConfig.INCENTIVEPOOL_ADDRESS, incentive_pool)
            assert current_annual_developer_foundation_amount == developer_foundation, "{} Year Developer Foundation Address:{} Balance:{}".format(
                i + 1,
                EconomicConfig.DEVELOPER_FOUNDATAION_ADDRESS,
                developer_foundation)
            assert current_annual_foundation_amount == foundation_balance, "{} Year Developer Foundation Address:{} Balance:{}".format(
                i + 1,
                EconomicConfig.DEVELOPER_FOUNDATAION_ADDRESS,
                developer_foundation)
            # log.info("{} Year additional Balance:{}".format(i + 1, total_amount_of_issuance))
            time.sleep(5)
            economic.wait_settlement_blocknum(node)
            while remaining_settlement_cycle != 1:
                tmp_current_block = node.eth.blockNumber
                if tmp_current_block % economic.settlement_size == 0:
                    time.sleep(1)
                block_info = node.eth.getBlock(1)
                log.info("block_info：{}".format(block_info))
                first_timestamp = block_info['timestamp']
                log.info("First block timestamp： {}".format(first_timestamp))
                end_cycle_timestamp = first_timestamp + (economic.additional_cycle_time * 60000)
                log.info("End time stamp of current issue cycle： {}".format(end_cycle_timestamp))

                last_settlement_block = (math.ceil(tmp_current_block / economic.settlement_size) - 1) * economic.settlement_size
                log.info("The last block height of the previous settlement period： {}".format(last_settlement_block))
                settlement_block_info = node.eth.getBlock(last_settlement_block)
                settlement_timestamp = settlement_block_info['timestamp']
                log.info("High block timestamp at the end settlement cycle： {}".format(settlement_timestamp))
                remaining_additional_time = end_cycle_timestamp - settlement_timestamp
                log.info("Remaining time of current issuance cycle： {}".format(remaining_additional_time))
                average_interval = (settlement_timestamp - first_timestamp) // (last_settlement_block - 1)
                log.info("Block interval in the last settlement cycle： {}".format(average_interval))
                number_of_remaining_blocks = math.ceil(remaining_additional_time / average_interval)
                log.info("Remaining block height of current issuance cycle： {}".format(number_of_remaining_blocks))
                remaining_settlement_cycle = math.ceil(number_of_remaining_blocks / economic.settlement_size)
                log.info(
                    "remaining settlement cycles in the current issuance cycle： {}".format(remaining_settlement_cycle))
                economic.wait_settlement_blocknum(node)

        elif 0 < i < 9:
            annual_last_block = (math.ceil(node.eth.blockNumber / economic.settlement_size) - 1) * economic.settlement_size
            log.info("The last block height in the last issue cycle: {}".format(annual_last_block))
            # Current annual total issuance
            additional_amount = int(Decimal(str(total_amount_of_issuance)) / Decimal(str(40)))
            log.info("Current annual quota： {}".format(additional_amount))
            # Incentive pool additional amount
            incentive_pool_additional_amount = int(Decimal(str(additional_amount)) * Decimal(str((80 / 100))))
            log.info("Additional quota for the current annual incentive pool: {}".format(incentive_pool_additional_amount))
            # developer foundation s additional amount
            developer_foundation_s_additional_amount = additional_amount - incentive_pool_additional_amount
            log.info("Current annual developer foundation additional quota: {}".format(developer_foundation_s_additional_amount))
            # Total amount of additional issuance
            total_amount_of_issuance = total_amount_of_issuance + additional_amount
            log.info("Total current hairstyle：{}".format(total_amount_of_issuance))
            # Current annual incentive pool amount
            current_incentive_pool = current_incentive_pool + incentive_pool_additional_amount + EconomicConfig.release_info[i - 1]['amount']
            log.info("Balance to be allocated for the current annual incentive pool：{}".format(current_incentive_pool))
            # Current annual Developer Fund Amount
            developer_foundation = developer_foundation + developer_foundation_s_additional_amount
            log.info("Current Annual Developer Foundation Total： {}".format(developer_foundation))
            # Query the current annual incentive pool amount
            current_annual_incentive_pool_amount = node.eth.getBalance(EconomicConfig.INCENTIVEPOOL_ADDRESS, annual_last_block)
            # Query current annual developer foundation amount
            DEVELOPER_FOUNDATAION_ADDRESS = node.web3.toChecksumAddress(EconomicConfig.DEVELOPER_FOUNDATAION_ADDRESS)
            current_annual_developer_foundation_amount = node.eth.getBalance(DEVELOPER_FOUNDATAION_ADDRESS, annual_last_block)
            # Query current annual fund amount
            FOUNDATION_ADDRESS = node.web3.toChecksumAddress(EconomicConfig.FOUNDATION_ADDRESS)
            current_annual_foundation_amount = node.eth.getBalance(FOUNDATION_ADDRESS, annual_last_block)
            log.info("{} year initialization incentive pool address: {} balance: {}".format(i + 1,
                                                                                            EconomicConfig.INCENTIVEPOOL_ADDRESS,
                                                                                            current_incentive_pool))
            log.info('{} Year Initialization Foundation Address: {} balance: {}'.format(i + 1,
                                                                                        EconomicConfig.FOUNDATION_ADDRESS,
                                                                                        foundation_balance))
            log.info("{} Year Developer Fund Address: {} balance: {}".format(i + 1,
                                                                             EconomicConfig.DEVELOPER_FOUNDATAION_ADDRESS,
                                                                             developer_foundation))
            log.info("{} Year additional balance:{}".format(i + 1, additional_amount))
            assert current_annual_incentive_pool_amount == current_incentive_pool, "{} year initialization incentive pool address: {} balance: {}".format(
                i + 1, EconomicConfig.INCENTIVEPOOL_ADDRESS, current_incentive_pool)
            assert current_annual_developer_foundation_amount == developer_foundation, "{} Year Developer Fund Address: {} balance: {}".format(
                i + 1, EconomicConfig.DEVELOPER_FOUNDATAION_ADDRESS, developer_foundation)
            assert current_annual_foundation_amount == 0, '{} Year Initialization Foundation Address: {} balance: {}'.format(
                i + 1, EconomicConfig.FOUNDATION_ADDRESS, foundation_balance)
            # Waiting for the end of the annual issuance cycle
            end_cycle_timestamp = end_cycle_timestamp + (economic.additional_cycle_time * 60000)
            log.info("End time stamp of current issue cycle： {}".format(end_cycle_timestamp))
            settlement_block_info = node.eth.getBlock(annual_last_block)
            settlement_timestamp = settlement_block_info['timestamp']
            log.info("High block timestamp at the end of settlement cycle： {}".format(settlement_timestamp))
            remaining_additional_time = end_cycle_timestamp - settlement_timestamp
            log.info("Remaining time of current issuance cycle： {}".format(remaining_additional_time))
            result = client.ppos.getAvgPackTime()
            average_interval = result['Ret']
            log.info("Block interval on the chain：{}".format(average_interval))
            log.info("Block interval on the chain：{}".format(average_interval))
            number_of_remaining_blocks = math.ceil(remaining_additional_time / average_interval)
            log.info("Remaining block height of current issuance cycle： {}".format(number_of_remaining_blocks))
            remaining_settlement_cycle = math.ceil(number_of_remaining_blocks / economic.settlement_size)
            log.info("remaining settlement cycles issuance cycle： {}".format(remaining_settlement_cycle))
            while remaining_settlement_cycle != 1:
                tmp_current_block = node.eth.blockNumber
                if tmp_current_block % economic.settlement_size == 0:
                    time.sleep(economic.interval)
                tmp_current_block = node.eth.blockNumber
                last_settlement_block = (math.ceil(
                    tmp_current_block / economic.settlement_size) - 1) * economic.settlement_size
                log.info("The last block height of the previous settlement period： {}".format(last_settlement_block))
                settlement_block_info = node.eth.getBlock(last_settlement_block)
                settlement_timestamp = settlement_block_info['timestamp']
                log.info("High block timestamp at the end of settlement cycle： {}".format(settlement_timestamp))
                remaining_additional_time = end_cycle_timestamp - settlement_timestamp
                log.info("Remaining time of current issuance cycle： {}".format(remaining_additional_time))
                result = client.ppos.getAvgPackTime()
                average_interval = result['Ret']
                log.info("Block interval on the chain：{}".format(average_interval))
                number_of_remaining_blocks = math.ceil(remaining_additional_time / average_interval)
                log.info("Remaining block height of current issuance cycle： {}".format(number_of_remaining_blocks))
                remaining_settlement_cycle = math.ceil(number_of_remaining_blocks / economic.settlement_size)
                log.info("remaining settlement cycles issuance cycle： {}".format(remaining_settlement_cycle))
                economic.wait_settlement_blocknum(node)
        else:
            annual_last_block = (math.ceil(node.eth.blockNumber / economic.settlement_size) - 1) * economic.settlement_size
            # Current annual total issuance
            additional_amount = int(Decimal(str(total_amount_of_issuance)) / Decimal(str(40)))
            # Incentive pool additional amount
            incentive_pool_additional_amount = int(Decimal(str(additional_amount)) * Decimal(str((80 / 100))))
            # developer foundation s additional amount
            developer_foundation_s_additional_amount = int(
                Decimal(str(additional_amount - incentive_pool_additional_amount)) * Decimal(str((50 / 100))))
            # Foundation grant additional amount
            foundation_grant_amount = additional_amount - incentive_pool_additional_amount - developer_foundation_s_additional_amount
            # Total amount of additional issuance
            total_amount_of_issuance = total_amount_of_issuance + additional_amount
            # Current annual incentive pool amount
            current_incentive_pool = current_incentive_pool + incentive_pool_additional_amount
            # Current annual Developer Fund Amount
            developer_foundation = developer_foundation + developer_foundation_s_additional_amount
            # Current annual fund amount
            foundation_balance = foundation_balance + foundation_grant_amount
            # Query the current annual incentive pool amount
            current_annual_incentive_pool_amount = node.eth.getBalance(EconomicConfig.INCENTIVEPOOL_ADDRESS, annual_last_block)
            # Query current annual developer foundation amount
            DEVELOPER_FOUNDATAION_ADDRESS = node.web3.toChecksumAddress(EconomicConfig.DEVELOPER_FOUNDATAION_ADDRESS)
            current_annual_developer_foundation_amount = node.eth.getBalance(DEVELOPER_FOUNDATAION_ADDRESS, annual_last_block)
            # Query current annual fund amount
            FOUNDATION_ADDRESS = node.web3.toChecksumAddress(EconomicConfig.FOUNDATION_ADDRESS)
            current_annual_foundation_amount = node.eth.getBalance(FOUNDATION_ADDRESS, annual_last_block)
            log.info("{} year initialization incentive pool address: {} balance: {}".format(i + 1,
                                                                                            EconomicConfig.INCENTIVEPOOL_ADDRESS,
                                                                                            current_incentive_pool))
            log.info('{} Year Initialization Foundation Address: {} balance: {}'.format(i + 1,
                                                                                        EconomicConfig.FOUNDATION_ADDRESS,
                                                                                        foundation_balance))
            log.info("{} Year Developer Fund Address: {} balance: {}".format(i + 1,
                                                                             EconomicConfig.DEVELOPER_FOUNDATAION_ADDRESS,
                                                                             developer_foundation))
            log.info("{} Year additional balance:{}".format(i + 1, additional_amount))
            assert current_annual_incentive_pool_amount == current_incentive_pool, "{} year initialization incentive pool address: {} balance: {}".format(
                i + 1, EconomicConfig.INCENTIVEPOOL_ADDRESS, current_incentive_pool)
            assert current_annual_developer_foundation_amount == developer_foundation, "{} Year Developer Fund Address: {} balance: {}".format(
                i + 1, EconomicConfig.DEVELOPER_FOUNDATAION_ADDRESS, developer_foundation)
            assert current_annual_foundation_amount == foundation_balance, '{} Year Initialization Foundation Address: {} balance: {}'.format(
                i + 1, EconomicConfig.FOUNDATION_ADDRESS, foundation_balance)


@pytest.mark.p1
def test_AL_FI_004_005(new_genesis_env, staking_cfg):
    """
    AL_FI_004:查看每年区块奖励变化
    AL_FI_005:查看每年质押奖励变化
    :param new_genesis_env:
    :param staking_cfg:
    :return:
    """
    # Initialization genesis file Initial amount
    node_count = len(new_genesis_env.consensus_node_list)
    default_pledge_amount = Web3.toWei(node_count * 1500000, 'ether')
    community_amount = default_pledge_amount + 259096239000000000000000000 + 62215742000000000000000000
    genesis = from_dict(data_class=Genesis, data=new_genesis_env.genesis_config)
    genesis.economicModel.innerAcc.cdfBalance = community_amount
    surplus_amount = str(EconomicConfig.TOKEN_TOTAL - community_amount - 200000000000000000000000000)
    genesis.alloc = {
        "lax1zqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqrzpqayr": {
            "balance": "200000000000000000000000000"
        },
        "lax196278ns22j23awdfj9f2d4vz0pedld8au6xelj": {
            "balance": surplus_amount
        }
    }
    new_file = new_genesis_env.cfg.env_tmp + "/genesis_0.13.0.json"
    genesis.to_file(new_file)
    new_genesis_env.deploy_all(new_file)
    normal_node = new_genesis_env.get_a_normal_node()
    client = Client(new_genesis_env, normal_node, staking_cfg)
    economic = client.economic
    node = client.node
    log.info("Current connection node：{}".format(node.node_mark))
    log.info("Current connection nodeid：{}".format(node.node_id))
    address, _ = economic.account.generate_account(node.web3, von_amount(economic.create_staking_limit, 2))
    log.info("address: {}".format(address))
    address1, _ = economic.account.generate_account(node.web3, 0)
    log.info("address1: {}".format(address1))
    end_cycle_timestamp = None
    for i in range(10):
        result = client.staking.create_staking(0, address1, address)
        assert_code(result, 0)
        # view account amount
        benifit_balance = node.eth.getBalance(address1)
        log.info("benifit_balance: {}".format(benifit_balance))
        # Wait for the settlement round to end
        economic.wait_settlement_blocknum(node)
        # 获取当前结算周期验证人
        verifier_list = node.ppos.getVerifierList()
        log.info("verifier_list: {}".format(verifier_list))
        # view block_reward
        block_reward, staking_reward = economic.get_current_year_reward(node)
        log.info("block_reward: {} staking_reward: {}".format(block_reward, staking_reward))
        # withdrew of pledge
        result = client.staking.withdrew_staking(address)
        assert_code(result, 0)
        # wait settlement block
        economic.wait_settlement_blocknum(node)
        # wait consensus block
        economic.wait_consensus_blocknum(node)
        # count the number of blocks
        blocknumber = economic.get_block_count_number(node, 10)
        log.info("blocknumber: {}".format(blocknumber))
        # view account amount again
        benifit_balance1 = node.eth.getBalance(address1)
        log.info("benifit_balance: {}".format(benifit_balance1))
        reward = int(blocknumber * Decimal(str(block_reward)))
        assert benifit_balance1 == benifit_balance + staking_reward + reward, "ErrMsg:benifit_balance: {}".format(
            benifit_balance1)
        if i == 0:
            block_info = node.eth.getBlock(1)
            log.info("block_info：{}".format(block_info))
            first_timestamp = block_info['timestamp']
            log.info("First block timestamp： {}".format(first_timestamp))
            end_cycle_timestamp = first_timestamp + (economic.additional_cycle_time * 60000)
            log.info("End time stamp of current issue cycle： {}".format(end_cycle_timestamp))
        else:
            # Waiting for the end of the annual issuance cycle
            end_cycle_timestamp = end_cycle_timestamp + (economic.additional_cycle_time * 60000)
            log.info("End time stamp of current issue cycle： {}".format(end_cycle_timestamp))
        annual_last_block = (math.ceil(node.eth.blockNumber / economic.settlement_size) - 1) * economic.settlement_size
        log.info("The last block height in the last issue cycle: {}".format(annual_last_block))
        settlement_block_info = node.eth.getBlock(annual_last_block)
        settlement_timestamp = settlement_block_info['timestamp']
        log.info("High block timestamp at the end of settlement cycle： {}".format(settlement_timestamp))
        remaining_additional_time = end_cycle_timestamp - settlement_timestamp
        log.info("Remaining time of current issuance cycle： {}".format(remaining_additional_time))
        result = client.ppos.getAvgPackTime()
        average_interval = result['Ret']
        log.info("Block interval on the chain：{}".format(average_interval))
        log.info("Block interval on the chain：{}".format(average_interval))
        number_of_remaining_blocks = math.ceil(remaining_additional_time / average_interval)
        log.info("Remaining block height of current issuance cycle： {}".format(number_of_remaining_blocks))
        remaining_settlement_cycle = math.ceil(number_of_remaining_blocks / economic.settlement_size)
        log.info("remaining settlement cycles issuance cycle： {}".format(remaining_settlement_cycle))
        while remaining_settlement_cycle != 1:
            tmp_current_block = node.eth.blockNumber
            if tmp_current_block % economic.settlement_size == 0:
                time.sleep(economic.interval)
            tmp_current_block = node.eth.blockNumber
            last_settlement_block = (math.ceil(
                tmp_current_block / economic.settlement_size) - 1) * economic.settlement_size
            log.info("The last block height of the previous settlement period： {}".format(last_settlement_block))
            settlement_block_info = node.eth.getBlock(last_settlement_block)
            settlement_timestamp = settlement_block_info['timestamp']
            log.info("High block timestamp at the end of settlement cycle： {}".format(settlement_timestamp))
            remaining_additional_time = end_cycle_timestamp - settlement_timestamp
            log.info("Remaining time of current issuance cycle： {}".format(remaining_additional_time))
            result = client.ppos.getAvgPackTime()
            average_interval = result['Ret']
            log.info("Block interval on the chain：{}".format(average_interval))
            number_of_remaining_blocks = math.ceil(remaining_additional_time / average_interval)
            log.info("Remaining block height of current issuance cycle： {}".format(number_of_remaining_blocks))
            remaining_settlement_cycle = math.ceil(number_of_remaining_blocks / economic.settlement_size)
            log.info("remaining settlement cycles issuance cycle： {}".format(remaining_settlement_cycle))
            economic.wait_settlement_blocknum(node)

