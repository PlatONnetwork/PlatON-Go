import math
import time
import pytest
import allure
from dacite import from_dict
from common.key import get_pub_key, mock_duplicate_sign, generate_key
from common.log import log
from client_sdk_python import Web3
from decimal import Decimal

from tests.conftest import get_client_consensus_obj
from tests.lib import EconomicConfig, Genesis, StakingConfig, Staking, check_node_in_list, assert_code, von_amount, \
    get_governable_parameter_value, Client, update_param_by_dict, get_param_by_dict


@pytest.mark.P1
def AL_FI_001_to_003(new_genesis_env, staking_cfg):
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
        "1000000000000000000000000000000000000003": {
            "balance": "200000000000000000000000000"
        },
        "0x2e95E3ce0a54951eB9A99152A6d5827872dFB4FD": {
            "balance": surplus_amount
        }
    }
    new_file = new_genesis_env.cfg.env_tmp + "/genesis.json"
    genesis.to_file(new_file)
    new_genesis_env.deploy_all(new_file)

    client = get_client_consensus_obj(new_genesis_env, staking_cfg)
    economic = client.economic
    node = client.node
    # Query the initial amount of incentive pool
    init_incentive_pool = 262215742000000000000000000
    # Query the initial amount of a warehouse lock plan
    init_foundationlockup = 259096239000000000000000000
    # Issued token amount
    init_token = 10000000000000000000000000000
    # Query developer foundation initial amount
    developer_foundation = 0
    # Query the initial amount of the foundation
    FOUNDATION = 0
    # Additional amount
    init_tt = 0
    # Annual issuance
    for i in range(10):
        if i == 0:
            INCENTIVEPOOL = init_incentive_pool
            FOUNDATIONLOCKUP = init_foundationlockup
            init_tt = int(init_token + Decimal(str(init_token)) / Decimal(str(40)))

            # Query the current annual incentive pool amount
            current_annual_incentive_pool_amount = node.eth.getBalance(EconomicConfig.INCENTIVEPOOL_ADDRESS, 0)
            # Query current annual developer foundation amount
            current_annual_developer_foundation_amount = node.eth.getBalance(
                node.web3.toChecksumAddress(EconomicConfig.DEVELOPER_FOUNDATAION_ADDRESS), 0)
            # Query current annual fund amount
            current_annual_foundation_amount = node.eth.getBalance(
                node.web3.toChecksumAddress(EconomicConfig.FOUNDATION_ADDRESS), 0)
            log.info(
                "{} Year Incentive Pool Address: {} Balance: {}".format(i + 1, EconomicConfig.INCENTIVEPOOL_ADDRESS,
                                                                        INCENTIVEPOOL))
            log.info('{} Year Foundation Address: {} Balance: {}'.format(i + 1, EconomicConfig.FOUNDATION_ADDRESS,
                                                                         FOUNDATION))
            log.info("{} Year Developer Foundation Address:{} Balance:{}".format(i + 1,
                                                                                 EconomicConfig.DEVELOPER_FOUNDATAION_ADDRESS,
                                                                                 developer_foundation))
            log.info("{} Year Foundation Locking Address: {} Balance: {}".format(i + 1,
                                                                                 EconomicConfig.FOUNDATION_LOCKUP_ADDRESS,
                                                                                 FOUNDATIONLOCKUP))
            assert current_annual_incentive_pool_amount == init_incentive_pool, "{} Year Incentive Pool Address: {} Balance: {}".format(
                i + 1, EconomicConfig.INCENTIVEPOOL_ADDRESS, INCENTIVEPOOL)
            assert current_annual_developer_foundation_amount == developer_foundation, "{} Year Developer Foundation Address:{} Balance:{}".format(
                i + 1,
                EconomicConfig.DEVELOPER_FOUNDATAION_ADDRESS,
                developer_foundation)
            assert current_annual_foundation_amount == FOUNDATION, "{} Year Developer Foundation Address:{} Balance:{}".format(
                i + 1,
                EconomicConfig.DEVELOPER_FOUNDATAION_ADDRESS,
                developer_foundation)
            log.info("{} Year additional Balance:{}".format(i + 1, init_tt))
            # Waiting for the end of the annual issuance cycle
            economic.wait_annual_blocknum(node)
        elif 0 < i < 9:
            FOUNDATION = 0
            # Current annual total issuance
            additional_amount = int(Decimal(str(init_tt)) / Decimal(str(40)))
            # Incentive pool additional amount
            incentive_pool_additional_amount = int(Decimal(str(additional_amount)) * Decimal(str((80 / 100))))
            # developer foundation s additional amount
            developer_foundation_s_additional_amount = additional_amount - incentive_pool_additional_amount
            # Total amount of additional issuance
            init_tt = init_tt + additional_amount
            # Current annual incentive pool amount
            init_incentive_pool = init_incentive_pool + incentive_pool_additional_amount + EconomicConfig.release_info[i - 1]['amount']
            # Current annual Developer Fund Amount
            developer_foundation = developer_foundation + developer_foundation_s_additional_amount
            # Query the current annual incentive pool amount
            current_annual_incentive_pool_amount = node.eth.getBalance(EconomicConfig.INCENTIVEPOOL_ADDRESS)
            # Query current annual developer foundation amount
            current_annual_developer_foundation_amount = node.eth.getBalance(
                node.web3.toChecksumAddress(EconomicConfig.DEVELOPER_FOUNDATAION_ADDRESS))
            # Query current annual fund amount
            current_annual_foundation_amount = node.eth.getBalance(
                node.web3.toChecksumAddress(EconomicConfig.FOUNDATION_ADDRESS))
            log.info("{} year initialization incentive pool address: {} balance: {}".format(i + 1,
                                                                                            EconomicConfig.INCENTIVEPOOL_ADDRESS,
                                                                                            init_incentive_pool))
            log.info('{} Year Initialization Foundation Address: {} balance: {}'.format(i + 1,
                                                                                        EconomicConfig.FOUNDATION_ADDRESS,
                                                                                        FOUNDATION))
            log.info("{} Year Developer Fund Address: {} balance: {}".format(i + 1,
                                                                             EconomicConfig.DEVELOPER_FOUNDATAION_ADDRESS,
                                                                             developer_foundation))
            log.info("{} Year additional balance:{}".format(i + 1, additional_amount))
            assert current_annual_incentive_pool_amount == init_incentive_pool, "{} year initialization incentive pool address: {} balance: {}".format(
                i + 1, EconomicConfig.INCENTIVEPOOL_ADDRESS, init_incentive_pool)
            assert current_annual_developer_foundation_amount == developer_foundation, "{} Year Developer Fund Address: {} balance: {}".format(
                i + 1, EconomicConfig.DEVELOPER_FOUNDATAION_ADDRESS, developer_foundation)
            assert current_annual_foundation_amount == 0, '{} Year Initialization Foundation Address: {} balance: {}'.format(
                i + 1, EconomicConfig.FOUNDATION_ADDRESS, FOUNDATION)
            # Waiting for the end of the annual issuance cycle
            economic.wait_annual_blocknum(node)
        else:
            # Current annual total issuance
            additional_amount = int(Decimal(str(init_tt)) / Decimal(str(40)))
            # Incentive pool additional amount
            incentive_pool_additional_amount = int(Decimal(str(additional_amount)) * Decimal(str((80 / 100))))
            # developer foundation s additional amount
            developer_foundation_s_additional_amount = int(
                Decimal(str(additional_amount - incentive_pool_additional_amount)) * Decimal(str((50 / 100))))
            # Foundation grant additional amount
            foundation_grant_amount = additional_amount - incentive_pool_additional_amount - developer_foundation_s_additional_amount
            # Total amount of additional issuance
            init_tt = init_tt + additional_amount
            # Current annual incentive pool amount
            init_incentive_pool = init_incentive_pool + incentive_pool_additional_amount
            # Current annual Developer Fund Amount
            developer_foundation = developer_foundation + developer_foundation_s_additional_amount
            # Current annual fund amount
            FOUNDATION = FOUNDATION + foundation_grant_amount
            # Query the current annual incentive pool amount
            current_annual_incentive_pool_amount = node.eth.getBalance(EconomicConfig.INCENTIVEPOOL_ADDRESS)
            # Query current annual developer foundation amount
            current_annual_developer_foundation_amount = node.eth.getBalance(
                node.web3.toChecksumAddress(EconomicConfig.DEVELOPER_FOUNDATAION_ADDRESS))
            # Query current annual fund amount
            current_annual_foundation_amount = node.eth.getBalance(
                node.web3.toChecksumAddress(EconomicConfig.FOUNDATION_ADDRESS))
            log.info("{} year initialization incentive pool address: {} balance: {}".format(i + 1,
                                                                                            EconomicConfig.INCENTIVEPOOL_ADDRESS,
                                                                                            init_incentive_pool))
            log.info('{} Year Initialization Foundation Address: {} balance: {}'.format(i + 1,
                                                                                        EconomicConfig.FOUNDATION_ADDRESS,
                                                                                        FOUNDATION))
            log.info("{} Year Developer Fund Address: {} balance: {}".format(i + 1,
                                                                             EconomicConfig.DEVELOPER_FOUNDATAION_ADDRESS,
                                                                             developer_foundation))
            log.info("{} Year additional balance:{}".format(i + 1, additional_amount))
            assert current_annual_incentive_pool_amount == init_incentive_pool, "{} year initialization incentive pool address: {} balance: {}".format(
                i + 1, EconomicConfig.INCENTIVEPOOL_ADDRESS, init_incentive_pool)
            assert current_annual_developer_foundation_amount == developer_foundation, "{} Year Developer Fund Address: {} balance: {}".format(
                i + 1, EconomicConfig.DEVELOPER_FOUNDATAION_ADDRESS, developer_foundation)
            assert current_annual_foundation_amount == FOUNDATION, '{} Year Initialization Foundation Address: {} balance: {}'.format(
                i + 1, EconomicConfig.FOUNDATION_ADDRESS, FOUNDATION)


@pytest.mark.p1
def AL_FI_004_005(new_genesis_env, staking_cfg):
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
        "1000000000000000000000000000000000000003": {
            "balance": "200000000000000000000000000"
        },
        "0x2e95E3ce0a54951eB9A99152A6d5827872dFB4FD": {
            "balance": surplus_amount
        }
    }
    new_file = new_genesis_env.cfg.env_tmp + "/genesis.json"
    genesis.to_file(new_file)
    new_genesis_env.deploy_all(new_file)
    normal_node = new_genesis_env.get_a_normal_node()
    client1 = Client(new_genesis_env, normal_node, staking_cfg)
    economic = client1.economic
    node = client1.node
    log.info("Current connection node：{}".format(node.node_mark))
    log.info("Current connection nodeid：{}".format(node.node_id))
    address, _ = client1.economic.account.generate_account(node.web3, von_amount(economic.create_staking_limit, 2))
    log.info("address: {}".format(address))
    address1, _ = client1.economic.account.generate_account(node.web3, 0)
    log.info("address1: {}".format(address1))
    for i in range(10):
        current_block = node.eth.blockNumber
        log.info("Current query block height： {}".format(node.eth.blockNumber))
        annualcycle = (economic.additional_cycle_time * 60) // economic.settlement_size
        annual_size = annualcycle * economic.settlement_size
        starting_block_height = math.floor(current_block / annual_size) * annual_size
        amount = node.eth.getBalance(EconomicConfig.INCENTIVEPOOL_ADDRESS, starting_block_height)
        log.info("Current annual incentive pool amount: {}".format(amount))
        # if i == 0:
        #     current_annual_incentive_pool_amount = 262215742000000000000000000
        # else:
        #     current_annual_incentive_pool_amount = node.eth.getBalance(EconomicConfig.INCENTIVEPOOL_ADDRESS)
        #     log.info("Current annual incentive pool amount: {}".format(current_annual_incentive_pool_amount))
        # Free amount application pledge node
        result = client1.staking.create_staking(0, address1, address)
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
        block_reward, staking_reward = economic.get_current_year_reward(node, amount=amount)
        log.info("block_reward: {} staking_reward: {}".format(block_reward, staking_reward))
        # withdrew of pledge
        result = client1.staking.withdrew_staking(address)
        assert_code(result, 0)
        # wait settlement block
        client1.economic.wait_settlement_blocknum(node)
        # wait consensus block
        client1.economic.wait_consensus_blocknum(node)
        # count the number of blocks
        blocknumber = client1.economic.get_block_count_number(client1.node, 10)
        log.info("blocknumber: {}".format(blocknumber))
        # view account amount again
        benifit_balance1 = node.eth.getBalance(address1)
        log.info("benifit_balance: {}".format(benifit_balance1))
        reward = int(blocknumber * Decimal(str(block_reward)))
        assert benifit_balance1 == benifit_balance + staking_reward + reward, "ErrMsg:benifit_balance: {}".format(
            benifit_balance1)
        # Waiting for the end of the annual increase
        economic.wait_annual_blocknum(node)
