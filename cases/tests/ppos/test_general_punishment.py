import time
import pytest
import allure
from dacite import from_dict
from common.key import get_pub_key, mock_duplicate_sign, generate_key
from common.log import log
from client_sdk_python import Web3
from decimal import Decimal

from tests.conftest import get_client_noconsensus_list
from tests.lib import EconomicConfig, Genesis, StakingConfig, Staking, check_node_in_list, assert_code, von_amount, \
    get_governable_parameter_value, Client, update_param_by_dict, get_param_by_dict


def get_out_block_penalty_parameters(client1, node, amount_type):
    # view Consensus Amount of pledge
    candidate_info = client1.ppos.getCandidateInfo(node.node_id)
    log.info("Pledge node information: {}".format(candidate_info))
    pledge_amount1 = candidate_info['Ret'][amount_type]
    # view block_reward
    log.info("block: {}".format(node.eth.blockNumber))
    block_reward, staking_reward = client1.economic.get_current_year_reward(node)
    log.info("block_reward: {} staking_reward: {}".format(block_reward, staking_reward))
    # Get governable parameters
    slash_blocks = get_governable_parameter_value(client1, 'slashBlocksReward')
    return pledge_amount1, block_reward, slash_blocks


def penalty_proportion_and_income(client_obj):
    # view Pledge amount
    candidate_info1 = client_obj.ppos.getCandidateInfo(client_obj.node.node_id)
    pledge_amount1 = candidate_info1['Ret']['Released']
    # view Parameter value before treatment
    penalty_ratio = get_governable_parameter_value(client_obj, 'slashFractionDuplicateSign')
    proportion_ratio = get_governable_parameter_value(client_obj, 'duplicateSignReportReward')
    return pledge_amount1, int(penalty_ratio), int(proportion_ratio)


@pytest.fixture()
def client_new_node_obj_list_reset(global_test_env, staking_cfg):
    """
    Get new node Client object list
    """
    global_test_env.deploy_all()
    yield get_client_noconsensus_list(global_test_env, staking_cfg)
    # global_test_env.deploy_all()


def VP_GPFV_001_002(client_new_node_obj_list_reset):
    """
    VP_GPFV_001:共识轮里一个节点出块两次，第一次出够10个块，第二次只出2个块
    VP_GPFV_002:出块数大于0小于每轮出块数
    :param client_new_node_obj_list_reset:
    :return:
    """
    client1 = client_new_node_obj_list_reset[0]
    log.info("Current connection node 1：{}".format(client1.node.node_mark))
    economic = client1.economic
    node = client1.node
    # create account
    address, _ = economic.account.generate_account(node.web3, von_amount(economic.create_staking_limit, 2))
    # create staking
    result = client1.staking.create_staking(0, address, address)
    assert_code(result, 0)
    # Wait for the settlement round to end
    economic.wait_settlement_blocknum(node)
    for i in range(4):
        result = check_node_in_list(node.node_id, client1.ppos.getValidatorList)
        log.info("Current node in consensus list status：{}".format(result))
        if result:
            log.info("Current block height: {}".format(client1.node.eth.blockNumber))
            # Wait for the 3 consensus round to end
            economic.wait_consensus_blocknum(node)
            # Get the number of pledge nodes out
            block_num = economic.get_block_count_number(node, 5)
            log.info("Number of pledge nodes：{}".format(block_num))
            candidate_info = client1.ppos.getCandidateInfo(node.node_id)
            log.info("Pledged node information：{}".format(candidate_info))
            info = candidate_info['Ret']
            assert info['Released'] == economic.create_staking_limit, "ErrMsg:Pledged Amount {}".format(info['Released'])
            validator_list = client1.ppos.getValidatorList()
            log.info("validator_list: {}".format(validator_list))
            assert len(validator_list['Ret']) == 4, "ErrMsg: Number of verification {}".format(len(validator_list))
        else:
            # wait consensus block
            economic.wait_consensus_blocknum(node)


@pytest.mark.P0
@pytest.mark.compatibility
def test_VP_GPFV_003(client_new_node_obj_list_reset):
    """
    低出快率最高处罚标准
    :param client_new_node_obj_list_reset:
    :return:
    """
    client1 = client_new_node_obj_list_reset[0]
    log.info("Current connection node1: {}".format(client1.node.node_mark))
    client2 = client_new_node_obj_list_reset[1]
    log.info("Current connection node2: {}".format(client2.node.node_mark))
    economic = client1.economic
    node = client1.node
    economic.env.deploy_all()
    # create account
    address, _ = economic.account.generate_account(node.web3, von_amount(economic.create_staking_limit, 2))
    # create staking
    result = client1.staking.create_staking(0, address, address)
    assert_code(result, 0)
    # Wait for the settlement round to end
    economic.wait_settlement_blocknum(node)
    # get pledge amount1 and block reward
    pledge_amount1, block_reward, slash_blocks = get_out_block_penalty_parameters(client1, node, 'Released')
    log.info("Current block height: {}".format(client1.node.eth.blockNumber))
    # stop node
    client1.node.stop()
    # Waiting for a settlement round
    client2.economic.wait_consensus_blocknum(client2.node, 3)
    log.info("Current block height: {}".format(client2.node.eth.blockNumber))
    # view verifier list
    verifier_list = client2.ppos.getVerifierList()
    log.info("verifier_list: {}".format(verifier_list))
    candidate_info = client2.ppos.getCandidateInfo(client1.node.node_id)
    log.info("Pledge node information： {}".format(candidate_info))
    pledge_amount2 = candidate_info['Ret']['Released']
    punishment_amonut = int(Decimal(str(block_reward)) * Decimal(str(slash_blocks)))
    log.info("punishment_amonut: {}".format(punishment_amonut))
    if punishment_amonut < pledge_amount1:
        assert pledge_amount2 == pledge_amount1 - punishment_amonut, "ErrMsg:Consensus Amount of pledge {}".format(
            pledge_amount2)
    else:
        assert pledge_amount2 == 0, "ErrMsg:Consensus Amount of pledge {}".format(pledge_amount2)


@pytest.mark.P0
def test_VP_GPFV_004(client_new_node_obj_list_reset):
    """
    锁仓质押被惩罚最高处罚标准
    :param client_new_node_obj_list_reset:
    :return:
    """
    client1 = client_new_node_obj_list_reset[0]
    log.info("Current connection node1: {}".format(client1.node.node_mark))
    client2 = client_new_node_obj_list_reset[1]
    log.info("Current connection node2: {}".format(client2.node.node_mark))
    economic = client1.economic
    node = client1.node
    # create account
    address, _ = economic.account.generate_account(node.web3, von_amount(economic.create_staking_limit, 2))
    # create Restricting Plan
    amount = von_amount(economic.create_staking_limit, 1)
    plan = [{'Epoch': 2, 'Amount': amount}]
    result = client1.restricting.createRestrictingPlan(address, plan, address)
    assert_code(result, 0)
    # create staking
    result = client1.staking.create_staking(1, address, address)
    assert_code(result, 0)
    # Wait for the settlement round to end
    economic.wait_settlement_blocknum(node)
    # get pledge amount1 and block reward
    pledge_amount1, block_reward, slash_blocks = get_out_block_penalty_parameters(client1, node, 'RestrictingPlan')
    log.info("Current block height: {}".format(client1.node.eth.blockNumber))
    # stop node
    client1.node.stop()
    # Waiting for a settlement round
    client2.economic.wait_consensus_blocknum(client2.node, 3)
    log.info("Current block height: {}".format(client2.node.eth.blockNumber))
    # view verifier list
    verifier_list = client2.ppos.getVerifierList()
    log.info("verifier_list: {}".format(verifier_list))
    candidate_info = client2.ppos.getCandidateInfo(client1.node.node_id)
    log.info("Pledge node information： {}".format(candidate_info))
    pledge_amount2 = candidate_info['Ret']['Released']
    punishment_amonut = int(Decimal(str(block_reward)) * Decimal(str(slash_blocks)))
    if punishment_amonut < pledge_amount1:
        assert pledge_amount2 == pledge_amount1 - punishment_amonut, "ErrMsg:Consensus Amount of pledge {}".format(
            pledge_amount2)
    else:
        assert pledge_amount2 == 0, "ErrMsg:Consensus Amount of pledge {}".format(pledge_amount2)


@pytest.mark.P2
def test_VP_GPFV_005(client_new_node_obj_list_reset):
    """
    锁仓增持/委托后被惩罚
    :param client_new_node_obj_list_reset:
    :return:
    """
    client1 = client_new_node_obj_list_reset[0]
    log.info("Current connection node1: {}".format(client1.node.node_mark))
    client2 = client_new_node_obj_list_reset[1]
    log.info("Current connection node2: {}".format(client2.node.node_mark))
    economic = client1.economic
    node = client1.node
    # create account
    address1, _ = economic.account.generate_account(node.web3, von_amount(economic.create_staking_limit, 3))
    # create account
    address2, _ = economic.account.generate_account(node.web3, von_amount(economic.create_staking_limit, 1))
    # create Restricting Plan
    amount1 = von_amount(economic.create_staking_limit, 2)
    plan = [{'Epoch': 1, 'Amount': amount1}]
    result = client1.restricting.createRestrictingPlan(address1, plan, address1)
    assert_code(result, 0)
    # create Restricting Plan
    amount2 = von_amount(economic.delegate_limit, 100)
    plan = [{'Epoch': 1, 'Amount': amount2}]
    result = client1.restricting.createRestrictingPlan(address2, plan, address2)
    assert_code(result, 0)
    # create staking
    result = client1.staking.create_staking(1, address1, address1)
    assert_code(result, 0)
    # increase staking
    result = client1.staking.increase_staking(1, address1)
    assert_code(result, 0)
    # Additional pledge
    result = client1.delegate.delegate(1, address2, amount=von_amount(economic.delegate_limit, 100))
    assert_code(result, 0)
    # Wait for the settlement round to end
    economic.wait_settlement_blocknum(node)
    # get pledge amount1 and block reward
    pledge_amount1, block_reward, slash_blocks = get_out_block_penalty_parameters(client1, node, 'RestrictingPlan')
    log.info("Current block height: {}".format(client1.node.eth.blockNumber))
    # stop node
    client1.node.stop()
    # Waiting for a settlement round
    client2.economic.wait_consensus_blocknum(client2.node, 3)
    log.info("Current block height: {}".format(client2.node.eth.blockNumber))
    # view verifier list
    verifier_list = client2.ppos.getVerifierList()
    log.info("verifier_list: {}".format(verifier_list))
    candidate_info = client2.ppos.getCandidateInfo(client1.node.node_id)
    log.info("Pledge node information： {}".format(candidate_info))
    pledge_amount2 = candidate_info['Ret']['Released']
    punishment_amonut = int(Decimal(str(block_reward)) * Decimal(str(slash_blocks)))
    if punishment_amonut < pledge_amount1:
        assert pledge_amount2 == pledge_amount1 - punishment_amonut, "ErrMsg:Consensus Amount of pledge {}".format(
            pledge_amount2)
    else:
        assert pledge_amount2 == 0, "ErrMsg:Consensus Amount of pledge {}".format(pledge_amount2)


@pytest.mark.P2
def test_VP_GPFV_006(client_new_node_obj_list_reset):
    """
    自由金额增持/委托后被惩罚
    :param client_new_node_obj_list_reset:
    :return:
    """
    client1 = client_new_node_obj_list_reset[0]
    log.info("Current connection node1: {}".format(client1.node.node_mark))
    client2 = client_new_node_obj_list_reset[1]
    log.info("Current connection node2: {}".format(client2.node.node_mark))
    economic = client1.economic
    node = client1.node
    # create account
    address1, _ = economic.account.generate_account(node.web3, von_amount(economic.create_staking_limit, 3))
    # create account
    address2, _ = economic.account.generate_account(node.web3, von_amount(economic.create_staking_limit, 1))
    # create staking
    result = client1.staking.create_staking(0, address1, address1)
    assert_code(result, 0)
    # increase staking
    result = client1.staking.increase_staking(0, address1)
    assert_code(result, 0)
    # Additional pledge
    result = client1.delegate.delegate(0, address2, amount=von_amount(economic.delegate_limit, 100))
    assert_code(result, 0)
    # Wait for the settlement round to end
    economic.wait_settlement_blocknum(node)
    # get pledge amount1 and block reward
    pledge_amount1, block_reward, slash_blocks = get_out_block_penalty_parameters(client1, node, 'Released')
    log.info("Current block height: {}".format(client1.node.eth.blockNumber))
    # stop node
    client1.node.stop()
    # Waiting for a settlement round
    client2.economic.wait_consensus_blocknum(client2.node, 3)
    log.info("Current block height: {}".format(client2.node.eth.blockNumber))
    # view verifier list
    verifier_list = client2.ppos.getVerifierList()
    log.info("verifier_list: {}".format(verifier_list))
    candidate_info = client2.ppos.getCandidateInfo(client1.node.node_id)
    log.info("Pledge node information： {}".format(candidate_info))
    pledge_amount2 = candidate_info['Ret']['Released']
    punishment_amonut = int(Decimal(str(block_reward)) * Decimal(str(slash_blocks)))
    if punishment_amonut < pledge_amount1:
        assert pledge_amount2 == pledge_amount1 - punishment_amonut, "ErrMsg:Consensus Amount of pledge {}".format(
            pledge_amount2)
    else:
        assert pledge_amount2 == 0, "ErrMsg:Consensus Amount of pledge {}".format(pledge_amount2)


@pytest.mark.P2
def test_VP_GPFV_007(client_new_node_obj_list_reset):
    """
    在被惩罚前退出质押
    :param client_new_node_obj_list_reset:
    :return:
    """
    client1 = client_new_node_obj_list_reset[0]
    log.info("Current connection node1: {}".format(client1.node.node_mark))
    client2 = client_new_node_obj_list_reset[1]
    log.info("Current connection node2: {}".format(client2.node.node_mark))
    economic = client1.economic
    node = client1.node
    # create account
    address, _ = economic.account.generate_account(node.web3, von_amount(economic.create_staking_limit, 2))
    # create staking
    result = client1.staking.create_staking(0, address, address)
    assert_code(result, 0)
    # Wait for the settlement round to end
    economic.wait_settlement_blocknum(node)
    # get pledge amount1 and block reward
    pledge_amount1, block_reward, slash_blocks = get_out_block_penalty_parameters(client1, node, 'Released')
    log.info("Current block height: {}".format(client1.node.eth.blockNumber))
    # stop node
    client1.node.stop()
    # Application for return of pledge
    result = client2.staking.withdrew_staking(address, node_id=node.node_id)
    assert_code(result, 0)
    # Waiting for a settlement round
    client2.economic.wait_consensus_blocknum(client2.node, 3)
    log.info("Current block height: {}".format(client2.node.eth.blockNumber))
    # view verifier list
    verifier_list = client2.ppos.getVerifierList()
    log.info("verifier_list: {}".format(verifier_list))
    candidate_info = client2.ppos.getCandidateInfo(client1.node.node_id)
    log.info("Pledge node information： {}".format(candidate_info))
    pledge_amount2 = candidate_info['Ret']['Released']
    punishment_amonut = int(Decimal(str(block_reward)) * Decimal(slash_blocks))
    if punishment_amonut < pledge_amount1:
        assert pledge_amount2 == pledge_amount1 - punishment_amonut, "ErrMsg:Consensus Amount of pledge {}".format(
            pledge_amount2)
    else:
        assert pledge_amount2 == 0, "ErrMsg:Consensus Amount of pledge {}".format(pledge_amount2)


@pytest.mark.P2
def test_VP_GPFV_008(client_new_node_obj_list_reset):
    """
    被处罚前进行增持
    :param client_new_node_obj_list_reset:
    :return:
    """
    client1 = client_new_node_obj_list_reset[0]
    log.info("Current connection node1: {}".format(client1.node.node_mark))
    client2 = client_new_node_obj_list_reset[1]
    log.info("Current connection node2: {}".format(client2.node.node_mark))
    economic = client1.economic
    node = client1.node
    # create account
    address, _ = economic.account.generate_account(node.web3, von_amount(economic.create_staking_limit, 3))
    # create staking
    result = client1.staking.create_staking(0, address, address)
    assert_code(result, 0)
    # Wait for the settlement round to end
    economic.wait_settlement_blocknum(node)
    # get pledge amount1 and block reward
    pledge_amount1, block_reward, slash_blocks = get_out_block_penalty_parameters(client1, node, 'Released')
    log.info("Current block height: {}".format(client1.node.eth.blockNumber))
    # stop node
    client1.node.stop()
    # Additional pledge
    result = client2.staking.increase_staking(0, address, node_id=node.node_id, amount=economic.create_staking_limit)
    assert_code(result, 0)
    # Waiting for a settlement round
    client2.economic.wait_consensus_blocknum(client2.node, 3)
    log.info("Current block height: {}".format(client2.node.eth.blockNumber))
    # view verifier list
    verifier_list = client2.ppos.getVerifierList()
    log.info("verifier_list: {}".format(verifier_list))
    candidate_info = client2.ppos.getCandidateInfo(client1.node.node_id)
    log.info("Pledge node information： {}".format(candidate_info))
    pledge_amount2 = candidate_info['Ret']['Released']
    punishment_amonut = int(Decimal(str(block_reward)) * Decimal(str(slash_blocks)))
    if punishment_amonut < pledge_amount1:
        assert pledge_amount2 == pledge_amount1 - punishment_amonut, "ErrMsg:Consensus Amount of pledge {}".format(
            pledge_amount2)
    else:
        assert pledge_amount2 == 0, "ErrMsg:Consensus Amount of pledge {}".format(pledge_amount2)


def test_VP_GPFV_009(client_new_node_obj_list_reset):
    """
    节点被处罚后马上重新质押（高惩罚）
    :param client_new_node_obj_list_reset:
    :return:
    """
    client1 = client_new_node_obj_list_reset[0]
    log.info("Current connection node1: {}".format(client1.node.node_mark))
    client2 = client_new_node_obj_list_reset[1]
    log.info("Current connection node2: {}".format(client2.node.node_mark))
    economic = client1.economic
    node = client1.node
    program_version = node.program_version
    program_version_sign = node.program_version_sign
    bls_pubkey = node.blspubkey
    bls_proof = node.schnorr_NIZK_prove
    # create account
    address, _ = economic.account.generate_account(node.web3, von_amount(economic.create_staking_limit, 3))
    # create staking
    result = client1.staking.create_staking(0, address, address)
    assert_code(result, 0)
    # Wait for the settlement round to end
    economic.wait_settlement_blocknum(node)
    log.info("Current block height: {}".format(client1.node.eth.blockNumber))
    # stop node
    node.stop()
    # Waiting for a settlement round
    client2.economic.wait_consensus_blocknum(client2.node, 3)
    log.info("Current block height: {}".format(client2.node.eth.blockNumber))
    # create staking again
    result = client2.staking.create_staking(0, address, address, node_id=node.node_id, program_version=program_version,
                                            program_version_sign=program_version_sign, bls_pubkey=bls_pubkey, bls_proof=bls_proof)
    assert_code(result, 301101)


@pytest.mark.P2
def test_VP_GPFV_010(client_new_node_obj_list_reset):
    """
    节点被处罚后马上重新增持质押（高惩罚）
    :param client_new_node_obj_list_reset:
    :return:
    """
    client1 = client_new_node_obj_list_reset[0]
    log.info("Current connection node1: {}".format(client1.node.node_mark))
    client2 = client_new_node_obj_list_reset[1]
    log.info("Current connection node2: {}".format(client2.node.node_mark))
    economic = client1.economic
    node = client1.node
    # create account
    address, _ = economic.account.generate_account(node.web3, von_amount(economic.create_staking_limit, 3))
    # create staking
    result = client1.staking.create_staking(0, address, address)
    assert_code(result, 0)
    # Wait for the settlement round to end
    economic.wait_settlement_blocknum(node)
    # stop node
    client1.node.stop()
    # Waiting for a settlement round
    client2.economic.wait_consensus_blocknum(client2.node, 3)
    log.info("Current block height: {}".format(client2.node.eth.blockNumber))
    # Additional pledge
    result = client2.staking.increase_staking(0, address, node_id=node.node_id)
    assert_code(result, 301103)


@pytest.mark.P2
def test_VP_GPFV_011(client_new_node_obj_list_reset):
    """
    先低出块率再双签
    :param client_new_node_obj_list_reset:
    :return:
    """
    client1 = client_new_node_obj_list_reset[0]
    log.info("Current connection node1: {}".format(client1.node.node_mark))
    client2 = client_new_node_obj_list_reset[1]
    log.info("Current connection node2: {}".format(client2.node.node_mark))
    economic = client1.economic
    node = client1.node
    # create account
    address, _ = economic.account.generate_account(node.web3, von_amount(economic.create_staking_limit, 3))
    # create account
    report_address, _ = economic.account.generate_account(node.web3, node.web3.toWei(1000, 'ether'))
    # create staking
    result = client1.staking.create_staking(0, address, address)
    assert_code(result, 0)
    # Wait for the settlement round to end
    economic.wait_settlement_blocknum(node)
    # get pledge amount1 and block reward
    pledge_amount1, block_reward, slash_blocks = get_out_block_penalty_parameters(client1, node, 'Released')
    for i in range(4):
        result = check_node_in_list(node.node_id, client2.ppos.getValidatorList)
        log.info("Current node in consensus list status：{}".format(result))
        if result:
            # stop node
            client1.node.stop()
            report_block = client2.node.eth.blockNumber
            log.info("Current block height: {}".format(report_block))
            # view Parameter value before treatment
            penalty_ratio = get_governable_parameter_value(client2, 'slashFractionDuplicateSign')
            proportion_ratio = get_governable_parameter_value(client2, 'duplicateSignReportReward')
            # view Amount of penalty
            proportion_reward, incentive_pool_reward = economic.get_report_reward(pledge_amount1, int(penalty_ratio),
                                                                                  int(proportion_ratio))
            # Obtain evidence of violation
            report_information = mock_duplicate_sign(1, client1.node.nodekey, client1.node.blsprikey, report_block)
            log.info("Report information: {}".format(report_information))
            result = client2.duplicatesign.reportDuplicateSign(1, report_information, report_address)
            assert_code(result, 0)
            # Query pledge node information:
            candidate_info = client2.ppos.getCandidateInfo(node.node_id)
            log.info("pledge node information: {}".format(candidate_info))
            info = candidate_info['Ret']
            # block_penalty = Decimal(str(block_reward)) * Decimal(str(slash_blocks))
            duplicateSign_penalty = proportion_reward + incentive_pool_reward
            # total_punish = block_penalty + duplicateSign_penalty
            assert info['Released'] == pledge_amount1 - duplicateSign_penalty, "ErrMsg:pledge node account {}".format(
                info['Released'])
            break
        else:
            # wait consensus block
            economic.wait_consensus_blocknum(node)


@pytest.mark.P2
def test_VP_GPFV_012(client_new_node_obj_list_reset):
    """
    先双签再低出块率
    :param client_new_node_obj_list_reset:
    :return:
    """
    client1 = client_new_node_obj_list_reset[0]
    log.info("Current connection node1: {}".format(client1.node.node_mark))
    client2 = client_new_node_obj_list_reset[1]
    log.info("Current connection node2: {}".format(client2.node.node_mark))
    economic = client1.economic
    node = client1.node
    # create account1
    address, _ = economic.account.generate_account(node.web3, von_amount(economic.create_staking_limit, 3))
    # create account2
    report_address, _ = economic.account.generate_account(node.web3, node.web3.toWei(1000, 'ether'))
    # create staking
    result = client1.staking.create_staking(0, address, address)
    assert_code(result, 0)
    # Wait for the settlement round to end
    economic.wait_settlement_blocknum(node)
    for i in range(4):
        result = check_node_in_list(node.node_id, client1.ppos.getValidatorList)
        log.info("Current node in consensus list status：{}".format(result))
        if result:
            # Query current block height
            report_block = client1.node.eth.blockNumber
            log.info("Current block height: {}".format(report_block))
            # Obtain penalty proportion and income
            pledge_amount1, penalty_ratio, proportion_ratio = penalty_proportion_and_income(client1)
            # view Amount of penalty
            proportion_reward, incentive_pool_reward = economic.get_report_reward(pledge_amount1, penalty_ratio,
                                                                                  proportion_ratio)
            # Obtain evidence of violation
            report_information = mock_duplicate_sign(1, client1.node.nodekey, client1.node.blsprikey, report_block)
            log.info("Report information: {}".format(report_information))
            result = client2.duplicatesign.reportDuplicateSign(1, report_information, report_address)
            assert_code(result, 0)
            # Waiting for a consensus round
            client2.economic.wait_consensus_blocknum(client2.node)
            # stop node
            client1.node.stop()
            # Waiting for 2 consensus round
            client2.economic.wait_consensus_blocknum(client2.node, 3)
            # view block_reward
            block_reward, staking_reward = client2.economic.get_current_year_reward(client2.node)
            log.info("block_reward: {} staking_reward: {}".format(block_reward, staking_reward))
            # Query pledge node information:
            candidate_info = client2.ppos.getCandidateInfo(node.node_id)
            log.info("pledge node information: {}".format(candidate_info))
            info = candidate_info['Ret']
            duplicateSign_penalty = proportion_reward + incentive_pool_reward
            assert info['Released'] == pledge_amount1 - duplicateSign_penalty, "ErrMsg:pledge node account {}".format(
                info['Released'])
            break
        else:
            # wait consensus block
            economic.wait_consensus_blocknum(node)


@pytest.mark.P2
def test_VP_GPFV_013(new_genesis_env, client_con_list_obj):
    """
    验证人被处罚后质押金=>创建验证人的最小质押门槛金额K
    :param new_genesis_env:
    :param client_con_list_obj:
    :return:
    """
    # Change configuration parameters
    genesis = from_dict(data_class=Genesis, data=new_genesis_env.genesis_config)
    genesis.economicModel.slashing.slashBlocksReward = 5
    new_file = new_genesis_env.cfg.env_tmp + "/genesis.json"
    genesis.to_file(new_file)
    new_genesis_env.deploy_all(new_file)

    client1 = client_con_list_obj[0]
    log.info("Current connection node1: {}".format(client1.node.node_mark))
    client2 = client_con_list_obj[1]
    log.info("Current connection node2: {}".format(client2.node.node_mark))
    economic = client1.economic
    node = client1.node
    # Wait for the consensus round to end
    economic.wait_consensus_blocknum(node, 1)
    # get pledge amount1 and block reward
    pledge_amount1, block_reward, slash_blocks = get_out_block_penalty_parameters(client1, node, 'Released')
    log.info("Current block height: {}".format(client1.node.eth.blockNumber))
    # stop node
    client1.node.stop()
    # Waiting for a 3 consensus round
    client2.economic.wait_consensus_blocknum(client2.node, 3)
    log.info("Current block height: {}".format(client2.node.eth.blockNumber))
    # view verifier list
    verifier_list = client2.ppos.getVerifierList()
    log.info("verifier_list: {}".format(verifier_list))
    candidate_info = client2.ppos.getCandidateInfo(client1.node.node_id)
    log.info("Pledge node information： {}".format(candidate_info))
    pledge_amount2 = candidate_info['Ret']['Released']
    punishment_amonut = int(Decimal(str(block_reward)) * Decimal(str(slash_blocks)))
    log.info("punishment_amonut: {}".format(punishment_amonut))
    assert pledge_amount2 == pledge_amount1 - punishment_amonut * 2, "ErrMsg:Consensus Amount of pledge {}".format(
        pledge_amount2)


@pytest.mark.P2
def test_VP_GPFV_014(new_genesis_env, client_noc_list_obj):
    """
    低出块率被最高处罚金低于质押金额（自由金额质押）
    :param new_genesis_env:
    :param client_noc_list_obj:
    :return:
    """
    # Change configuration parameters
    genesis = from_dict(data_class=Genesis, data=new_genesis_env.genesis_config)
    genesis.economicModel.slashing.slashBlocksReward = 5
    new_file = new_genesis_env.cfg.env_tmp + "/genesis.json"
    genesis.to_file(new_file)
    new_genesis_env.deploy_all(new_file)

    client1 = client_noc_list_obj[0]
    log.info("Current connection node1: {}".format(client1.node.node_mark))
    client2 = client_noc_list_obj[1]
    log.info("Current connection node2: {}".format(client2.node.node_mark))
    economic = client1.economic
    node = client1.node
    # create account
    address, _ = economic.account.generate_account(node.web3, von_amount(economic.create_staking_limit, 3))
    # create Restricting Plan
    amount = economic.create_staking_limit
    plan = [{'Epoch': 1, 'Amount': amount}]
    result = client1.restricting.createRestrictingPlan(address, plan, address)
    assert_code(result, 0)
    # create staking
    result = client1.staking.create_staking(0, address, address)
    assert_code(result, 0)
    # increase staking
    increase_amount = von_amount(economic.create_staking_limit, 0.5)
    result = client1.staking.increase_staking(1, address, amount=increase_amount)
    assert_code(result, 0)
    # Wait for the settlement round to end
    economic.wait_settlement_blocknum(node)
    # get pledge amount1 and block reward
    pledge_amount1, block_reward, slash_blocks = get_out_block_penalty_parameters(client1, node, 'Released')
    log.info("Current block height: {}".format(client1.node.eth.blockNumber))
    # stop node
    client1.node.stop()
    # Waiting for a settlement round
    client2.economic.wait_consensus_blocknum(client2.node, 3)
    log.info("Current block height: {}".format(client2.node.eth.blockNumber))
    # view verifier list
    verifier_list = client2.ppos.getVerifierList()
    log.info("verifier_list: {}".format(verifier_list))
    candidate_info = client2.ppos.getCandidateInfo(client1.node.node_id)
    log.info("Pledge node information： {}".format(candidate_info))
    info = candidate_info['Ret']
    pledge_amount2 = info['Released']
    pledge_amount3 = info['RestrictingPlan']
    punishment_amonut = int(Decimal(str(block_reward)) * Decimal(str(slash_blocks)))
    log.info("punishment_amonut: {}".format(punishment_amonut))
    assert (pledge_amount2 == pledge_amount1 - punishment_amonut * 2) or (pledge_amount2 == pledge_amount1 - punishment_amonut), "ErrMsg:Pledge Released {}".format(
        pledge_amount2)
    assert pledge_amount3 == increase_amount, "ErrMsg:Pledge RestrictingPlan {}".format(pledge_amount3)


@pytest.mark.P2
def test_VP_GPFV_015(new_genesis_env, client_noc_list_obj):
    """
    低出块率被最高处罚金等于于自由处罚金（自由金额质押）
    :param new_genesis_env:
    :param client_noc_list_obj:
    :return:
    """
    # Change configuration parameters
    genesis = from_dict(data_class=Genesis, data=new_genesis_env.genesis_config)
    genesis.economicModel.slashing.slashBlocksReward = 13
    new_file = new_genesis_env.cfg.env_tmp + "/genesis.json"
    genesis.to_file(new_file)
    new_genesis_env.deploy_all(new_file)

    client1 = client_noc_list_obj[0]
    log.info("Current connection node1: {}".format(client1.node.node_mark))
    client2 = client_noc_list_obj[1]
    log.info("Current connection node2: {}".format(client2.node.node_mark))
    economic = client1.economic
    node = client1.node
    # create account
    address, _ = economic.account.generate_account(node.web3, von_amount(economic.create_staking_limit, 4))
    # create Restricting Plan
    amount = economic.create_staking_limit
    plan = [{'Epoch': 1, 'Amount': amount}]
    result = client1.restricting.createRestrictingPlan(address, plan, address)
    assert_code(result, 0)
    # view block_reward
    block_reward, staking_reward = client1.economic.get_current_year_reward(node)
    log.info("block_reward: {} staking_reward: {}".format(block_reward, staking_reward))
    # Get governable parameters
    slash_blocks = get_governable_parameter_value(client1, 'slashBlocksReward')
    # create staking
    staking_amount = int(Decimal(str(block_reward)) * Decimal(slash_blocks))
    result = client1.staking.create_staking(0, address, address, amount=staking_amount * 2)
    assert_code(result, 0)
    # increase staking
    increase_amount = von_amount(economic.create_staking_limit, 0.5)
    result = client1.staking.increase_staking(1, address, amount=increase_amount)
    assert_code(result, 0)
    # Wait for the settlement round to end
    economic.wait_settlement_blocknum(node)
    # view Consensus Amount of pledge
    candidate_info = client1.ppos.getCandidateInfo(node.node_id)
    log.info("Pledge node information: {}".format(candidate_info))
    pledge_amount1 = candidate_info['Ret']['Released']
    log.info("Current block height: {}".format(client1.node.eth.blockNumber))
    # stop node
    client1.node.stop()
    # Waiting for a settlement round
    client2.economic.wait_consensus_blocknum(client2.node, 3)
    log.info("Current block height: {}".format(client2.node.eth.blockNumber))
    # view verifier list
    verifier_list = client2.ppos.getVerifierList()
    log.info("verifier_list: {}".format(verifier_list))
    candidate_info = client2.ppos.getCandidateInfo(client1.node.node_id)
    log.info("Pledge node information： {}".format(candidate_info))
    info = candidate_info['Ret']
    pledge_amount2 = info['Released']
    pledge_amount3 = info['RestrictingPlan']
    punishment_amonut = int(Decimal(str(block_reward)) * Decimal(str(slash_blocks)))
    assert pledge_amount2 == pledge_amount1 - punishment_amonut * 2, "ErrMsg:Pledge Released {}".format(
        pledge_amount2)
    assert pledge_amount3 == increase_amount, "ErrMsg:Pledge RestrictingPlan {}".format(pledge_amount3)


@pytest.mark.P2
def test_VP_GPFV_016(new_genesis_env, client_noc_list_obj):
    """
    低出块率被最高处罚金大于自由处罚金（自由金额质押）
    :param new_genesis_env:
    :param client_noc_list_obj:
    :return:
    """
    # Change configuration parameters
    genesis = from_dict(data_class=Genesis, data=new_genesis_env.genesis_config)
    genesis.economicModel.slashing.slashBlocksReward = 13
    new_file = new_genesis_env.cfg.env_tmp + "/genesis.json"
    genesis.to_file(new_file)
    new_genesis_env.deploy_all(new_file)

    client1 = client_noc_list_obj[0]
    log.info("Current connection node1: {}".format(client1.node.node_mark))
    client2 = client_noc_list_obj[1]
    log.info("Current connection node2: {}".format(client2.node.node_mark))
    economic = client1.economic
    node = client1.node
    # create account
    address, _ = economic.account.generate_account(node.web3, von_amount(economic.create_staking_limit, 4))
    # create Restricting Plan
    amount = economic.create_staking_limit
    plan = [{'Epoch': 1, 'Amount': amount}]
    result = client1.restricting.createRestrictingPlan(address, plan, address)
    assert_code(result, 0)
    # view block_reward
    block_reward, staking_reward = client1.economic.get_current_year_reward(node)
    log.info("block_reward: {} staking_reward: {}".format(block_reward, staking_reward))
    # Get governable parameters
    slash_blocks = get_governable_parameter_value(client1, 'slashBlocksReward')
    # create staking
    staking_amount = von_amount(economic.create_staking_limit, 2)
    result = client1.staking.create_staking(0, address, address, amount=staking_amount)
    assert_code(result, 0)
    # increase staking
    increase_amount = von_amount(economic.create_staking_limit, 0.5)
    result = client1.staking.increase_staking(1, address, amount=increase_amount)
    assert_code(result, 0)
    # Wait for the settlement round to end
    economic.wait_settlement_blocknum(node)
    # view Consensus Amount of pledge
    candidate_info = client1.ppos.getCandidateInfo(node.node_id)
    log.info("Pledge node information: {}".format(candidate_info))
    pledge_amount1 = candidate_info['Ret']['Released']
    log.info("Current block height: {}".format(client1.node.eth.blockNumber))
    # stop node
    client1.node.stop()
    # Waiting for a settlement round
    client2.economic.wait_consensus_blocknum(client2.node, 3)
    log.info("Current block height: {}".format(client2.node.eth.blockNumber))
    # view verifier list
    verifier_list = client2.ppos.getVerifierList()
    log.info("verifier_list: {}".format(verifier_list))
    candidate_info = client2.ppos.getCandidateInfo(client1.node.node_id)
    log.info("Pledge node information： {}".format(candidate_info))
    info = candidate_info['Ret']
    pledge_amount2 = info['Released']
    pledge_amount3 = info['RestrictingPlan']
    punishment_amonut = int(Decimal(str(block_reward)) * Decimal(str(slash_blocks)))
    log.info("punishment_amonut: {}".format(punishment_amonut))
    assert (pledge_amount2 == 0) or (pledge_amount2 == pledge_amount1 - punishment_amonut), "ErrMsg:Pledge Released {}".format(
        pledge_amount2)
    assert (pledge_amount3 == increase_amount - (punishment_amonut * 2 - pledge_amount1)) or (pledge_amount3 == 0), "ErrMsg:Pledge RestrictingPlan {}".format(pledge_amount3)


@pytest.mark.P2
def test_VP_GPFV_017(new_genesis_env, client_noc_list_obj):
    """
    低出块率被最高处罚金低于质押金额（锁仓金额质押）
    :param new_genesis_env:
    :param client_noc_list_obj:
    :return:
    """
    # Change configuration parameters
    genesis = from_dict(data_class=Genesis, data=new_genesis_env.genesis_config)
    genesis.economicModel.slashing.slashBlocksReward = 5
    new_file = new_genesis_env.cfg.env_tmp + "/genesis.json"
    genesis.to_file(new_file)
    new_genesis_env.deploy_all(new_file)

    client1 = client_noc_list_obj[0]
    log.info("Current connection node1: {}".format(client1.node.node_mark))
    client2 = client_noc_list_obj[1]
    log.info("Current connection node2: {}".format(client2.node.node_mark))
    economic = client1.economic
    node = client1.node
    # create account
    address, _ = economic.account.generate_account(node.web3, von_amount(economic.create_staking_limit, 3))
    # create Restricting Plan
    amount = economic.create_staking_limit
    plan = [{'Epoch': 1, 'Amount': amount}]
    result = client1.restricting.createRestrictingPlan(address, plan, address)
    assert_code(result, 0)
    # create staking
    result = client1.staking.create_staking(1, address, address)
    assert_code(result, 0)
    # increase staking
    increase_amount = von_amount(economic.create_staking_limit, 0.5)
    result = client1.staking.increase_staking(0, address, amount=increase_amount)
    assert_code(result, 0)
    # Wait for the settlement round to end
    economic.wait_settlement_blocknum(node)
    # get pledge amount1 and block reward
    pledge_amount1, block_reward, slash_blocks = get_out_block_penalty_parameters(client1, node, 'Released')
    log.info("Current block height: {}".format(client1.node.eth.blockNumber))
    # stop node
    client1.node.stop()
    # Waiting for a settlement round
    client2.economic.wait_consensus_blocknum(client2.node, 3)
    log.info("Current block height: {}".format(client2.node.eth.blockNumber))
    # view verifier list
    verifier_list = client2.ppos.getVerifierList()
    log.info("verifier_list: {}".format(verifier_list))
    candidate_info = client2.ppos.getCandidateInfo(client1.node.node_id)
    log.info("Pledge node information： {}".format(candidate_info))
    info = candidate_info['Ret']
    pledge_amount2 = info['Released']
    pledge_amount3 = info['RestrictingPlan']
    punishment_amonut = int(Decimal(str(block_reward)) * Decimal(str(slash_blocks)))
    log.info("punishment_amonut: {}".format(punishment_amonut))
    assert pledge_amount2 == 0, "ErrMsg:Pledge Released {}".format(
        pledge_amount2)
    assert pledge_amount3 == economic.create_staking_limit - (punishment_amonut * 2 - increase_amount), "ErrMsg:Pledge RestrictingPlan {}".format(pledge_amount3)


@pytest.mark.P2
def test_VP_GPFV_018(new_genesis_env, client_noc_list_obj):
    """
    低出块率被最高处罚金等于质押金额（锁仓金额质押）
    :param new_genesis_env:
    :param client_noc_list_obj:
    :return:
    """
    # Change configuration parameters
    genesis = from_dict(data_class=Genesis, data=new_genesis_env.genesis_config)
    genesis.economicModel.slashing.slashBlocksReward = 13
    new_file = new_genesis_env.cfg.env_tmp + "/genesis.json"
    genesis.to_file(new_file)
    new_genesis_env.deploy_all(new_file)

    client1 = client_noc_list_obj[0]
    log.info("Current connection node1: {}".format(client1.node.node_mark))
    client2 = client_noc_list_obj[1]
    log.info("Current connection node2: {}".format(client2.node.node_mark))
    economic = client1.economic
    node = client1.node
    # create account
    address, _ = economic.account.generate_account(node.web3, von_amount(economic.create_staking_limit, 4))
    # create Restricting Plan
    amount = von_amount(economic.create_staking_limit, 3)
    plan = [{'Epoch': 1, 'Amount': amount}]
    result = client1.restricting.createRestrictingPlan(address, plan, address)
    assert_code(result, 0)
    # view block_reward
    block_reward, staking_reward = client1.economic.get_current_year_reward(node)
    log.info("block_reward: {} staking_reward: {}".format(block_reward, staking_reward))
    # Get governable parameters
    slash_blocks = get_governable_parameter_value(client1, 'slashBlocksReward')
    # create staking
    staking_amount = von_amount(block_reward, 26)
    log.info("staking_amount: {}".format(staking_amount))
    result = client1.staking.create_staking(1, address, address, amount=staking_amount)
    assert_code(result, 0)
    # increase staking
    increase_amount = von_amount(economic.create_staking_limit, 0.5)
    result = client1.staking.increase_staking(0, address, amount=increase_amount)
    assert_code(result, 0)
    # Wait for the settlement round to end
    economic.wait_settlement_blocknum(node)
    # view Consensus Amount of pledge
    candidate_info = client1.ppos.getCandidateInfo(node.node_id)
    log.info("Pledge node information: {}".format(candidate_info))
    pledge_amount1 = candidate_info['Ret']['Released']
    log.info("Current block height: {}".format(client1.node.eth.blockNumber))
    # stop node
    client1.node.stop()
    # Waiting for a settlement round
    client2.economic.wait_consensus_blocknum(client2.node, 3)
    log.info("Current block height: {}".format(client2.node.eth.blockNumber))
    # view verifier list
    verifier_list = client2.ppos.getVerifierList()
    log.info("verifier_list: {}".format(verifier_list))
    candidate_info = client2.ppos.getCandidateInfo(client1.node.node_id)
    log.info("Pledge node information： {}".format(candidate_info))
    info = candidate_info['Ret']
    pledge_amount2 = info['Released']
    pledge_amount3 = info['RestrictingPlan']
    punishment_amonut = int(Decimal(str(block_reward)) * Decimal(str(slash_blocks)))
    log.info("punishment_amonut: {}".format(punishment_amonut))
    assert pledge_amount2 == 0, "ErrMsg:Pledge Released {}".format(pledge_amount2)
    assert (pledge_amount3 == staking_amount - (von_amount(punishment_amonut, 2) - increase_amount)) or (pledge_amount3 == staking_amount - (punishment_amonut - increase_amount)), "ErrMsg:Pledge RestrictingPlan {}".format(pledge_amount3)


@pytest.mark.P2
def test_VP_GPFV_019(new_genesis_env, client_noc_list_obj):
    """
    低出块率被最高处罚金大于质押金额（锁仓金额质押）
    :param new_genesis_env:
    :param client_noc_list_obj:
    :return:
    """
    # Change configuration parameters
    genesis = from_dict(data_class=Genesis, data=new_genesis_env.genesis_config)
    genesis.economicModel.slashing.slashBlocksReward = 13
    new_file = new_genesis_env.cfg.env_tmp + "/genesis.json"
    genesis.to_file(new_file)
    new_genesis_env.deploy_all(new_file)

    client1 = client_noc_list_obj[0]
    log.info("Current connection node1: {}".format(client1.node.node_mark))
    client2 = client_noc_list_obj[1]
    log.info("Current connection node2: {}".format(client2.node.node_mark))
    economic = client1.economic
    node = client1.node
    # create account
    address, _ = economic.account.generate_account(node.web3, von_amount(economic.create_staking_limit, 4))
    # create Restricting Plan
    amount = von_amount(economic.create_staking_limit, 2)
    plan = [{'Epoch': 1, 'Amount': amount}]
    result = client1.restricting.createRestrictingPlan(address, plan, address)
    assert_code(result, 0)
    # view block_reward
    block_reward, staking_reward = client1.economic.get_current_year_reward(node)
    log.info("block_reward: {} staking_reward: {}".format(block_reward, staking_reward))
    # Get governable parameters
    slash_blocks = get_governable_parameter_value(client1, 'slashBlocksReward')
    # create staking
    result = client1.staking.create_staking(1, address, address, amount=amount)
    assert_code(result, 0)
    # increase staking
    increase_amount = von_amount(economic.create_staking_limit, 0.5)
    result = client1.staking.increase_staking(0, address, amount=increase_amount)
    assert_code(result, 0)
    # Wait for the settlement round to end
    economic.wait_settlement_blocknum(node)
    # view Consensus Amount of pledge
    candidate_info = client1.ppos.getCandidateInfo(node.node_id)
    log.info("Pledge node information: {}".format(candidate_info))
    pledge_amount1 = candidate_info['Ret']['Released']
    log.info("Current block height: {}".format(client1.node.eth.blockNumber))
    # stop node
    client1.node.stop()
    # Waiting for a settlement round
    client2.economic.wait_consensus_blocknum(client2.node, 3)
    log.info("Current block height: {}".format(client2.node.eth.blockNumber))
    # view verifier list
    verifier_list = client2.ppos.getVerifierList()
    log.info("verifier_list: {}".format(verifier_list))
    candidate_info = client2.ppos.getCandidateInfo(client1.node.node_id)
    log.info("Pledge node information： {}".format(candidate_info))
    info = candidate_info['Ret']
    pledge_amount2 = info['Released']
    pledge_amount3 = info['RestrictingPlan']
    punishment_amonut = int(Decimal(str(block_reward)) * Decimal(str(slash_blocks)))
    log.info("punishment_amonut: {}".format(punishment_amonut))
    assert pledge_amount2 == 0, "ErrMsg:Pledge Released {}".format(pledge_amount2)
    assert pledge_amount3 == amount - (punishment_amonut * 2 - pledge_amount1), "ErrMsg:Pledge RestrictingPlan {}".format(pledge_amount3)


@pytest.mark.P2
def test_VP_GPFV_020(new_genesis_env, client_noc_list_obj):
    """
    移出PlatON验证人与候选人名单，（扣除以后剩余自有质押金），未申请退回质押金
    :param client_noc_list_obj:
    :return:
    """
    # Change configuration parameters
    genesis = from_dict(data_class=Genesis, data=new_genesis_env.genesis_config)
    genesis.economicModel.slashing.slashBlocksReward = 5
    new_file = new_genesis_env.cfg.env_tmp + "/genesis.json"
    genesis.to_file(new_file)
    new_genesis_env.deploy_all(new_file)

    client1 = client_noc_list_obj[0]
    log.info("Current connection node1: {}".format(client1.node.node_mark))
    client2 = client_noc_list_obj[1]
    log.info("Current connection node2: {}".format(client2.node.node_mark))
    economic = client1.economic
    node = client1.node
    # create account
    address, _ = economic.account.generate_account(node.web3, von_amount(economic.create_staking_limit, 2))
    address1, _ = economic.account.generate_account(node.web3, 0)
    # create staking
    result = client1.staking.create_staking(0, address1, address)
    assert_code(result, 0)
    # Wait for the settlement round to end
    economic.wait_settlement_blocknum(node)
    # get pledge amount1 and block reward
    pledge_amount1, block_reward, slash_blocks = get_out_block_penalty_parameters(client1, node, 'Released')
    log.info("Current block height: {}".format(client1.node.eth.blockNumber))
    # stop node
    client1.node.stop()
    # Waiting for a settlement round
    client2.economic.wait_consensus_blocknum(client2.node, 3)
    log.info("Current block height: {}".format(client2.node.eth.blockNumber))
    # view verifier list
    verifier_list = client2.ppos.getVerifierList()
    log.info("verifier_list: {}".format(verifier_list))
    candidate_info = client2.ppos.getCandidateInfo(client1.node.node_id)
    log.info("Pledge node information： {}".format(candidate_info))
    punishment_amonut = int(Decimal(str(block_reward)) * Decimal(str(slash_blocks)))
    log.info("punishment_amonut: {}".format(punishment_amonut))
    # Query pledge account balance
    balance1 = client2.node.eth.getBalance(address)
    log.info("pledge account balance: {}".format(balance1))
    # Wait for the 2 settlement round to end
    economic.wait_settlement_blocknum(client2.node, 2)
    # Query pledge account balance
    balance2 = client2.node.eth.getBalance(address)
    log.info("pledge account balance: {}".format(balance2))
    assert balance2 == balance1 + (pledge_amount1 - punishment_amonut * 2), "ErrMsg:pledge account balance {}".format(
        balance2)


@pytest.mark.P2
def test_VP_GPFV_021(client_new_node_obj_list_reset):
    """
    移出PlatON验证人与候选人名单，委托人可在处罚所在结算周期，申请赎回全部委托金
    :param client_new_node_obj_list_reset:
    :return:
    """
    client1 = client_new_node_obj_list_reset[0]
    log.info("Current connection node1: {}".format(client1.node.node_mark))
    client2 = client_new_node_obj_list_reset[1]
    log.info("Current connection node2: {}".format(client2.node.node_mark))
    economic = client1.economic
    node = client1.node
    # create pledge address
    pledge_address, _ = economic.account.generate_account(node.web3, von_amount(economic.create_staking_limit, 3))
    # create report address
    delegate_address, _ = economic.account.generate_account(node.web3, node.web3.toWei(1000, 'ether'))
    # create staking
    result = client1.staking.create_staking(0, pledge_address, pledge_address)
    assert_code(result, 0)
    # Additional pledge
    result = client1.delegate.delegate(0, delegate_address)
    assert_code(result, 0)
    # Wait for the settlement round to end
    economic.wait_settlement_blocknum(node)
    # stop node
    client1.node.stop()
    # Waiting for a settlement round
    client2.economic.wait_consensus_blocknum(client2.node, 3)
    log.info("Current block height: {}".format(client2.node.eth.blockNumber))
    # view verifier list
    verifier_list = client2.ppos.getVerifierList()
    log.info("verifier_list: {}".format(verifier_list))
    candidate_info = client2.ppos.getCandidateInfo(client1.node.node_id)
    log.info("Pledge node information： {}".format(candidate_info))
    time.sleep(3)
    # Access to pledge information
    candidate_info = client2.ppos.getCandidateInfo(node.node_id)
    info = candidate_info['Ret']
    staking_blocknum = info['StakingBlockNum']
    # To view the entrusted account balance
    delegate_balance = client2.node.eth.getBalance(delegate_address)
    log.info("report address balance: {}".format(delegate_balance))
    # withdrew delegate
    result = client2.delegate.withdrew_delegate(staking_blocknum, delegate_address, node_id=node.node_id)
    assert_code(result, 0)
    # To view the entrusted account balance
    delegate_balance1 = client2.node.eth.getBalance(delegate_address)
    log.info("report address balance: {}".format(delegate_balance1))
    assert delegate_balance + economic.delegate_limit - delegate_balance1 < client2.node.web3.toWei(1,
                                                                                                    'ether'), "ErrMsg:Ireport balance {}".format(
        delegate_balance1)
