import json
import time
import pytest
import allure
from dacite import from_dict
from common.key import mock_duplicate_sign
from common.log import log
from decimal import Decimal
from tests.conftest import param_governance_verify, param_governance_verify_before_endblock
from tests.lib import EconomicConfig, Genesis, check_node_in_list, assert_code, get_governable_parameter_value, \
    wait_block_number, von_amount


def pledge_punishment(client_con_list_obj):
    """
    :return:
    """
    first_index = 0
    second_index = 1
    first_client = client_con_list_obj[first_index]
    second_client = client_con_list_obj[second_index]
    log.info("Current block height: {}".format(first_client.node.eth.blockNumber))
    # stop node
    first_client.node.stop()
    # Waiting for a settlement round
    second_client.economic.wait_consensus_blocknum(second_client.node, 2)
    log.info("Current block height: {}".format(second_client.node.eth.blockNumber))
    # view verifier list
    verifier_list = second_client.ppos.getVerifierList()
    log.info("verifier_list: {}".format(verifier_list))
    candidate_info = second_client.ppos.getCandidateInfo(second_client.node.node_id)
    log.info("Pledge node information： {}".format(candidate_info))
    return candidate_info


def information_before_slash_blocks(client):
    node = client.node
    # view Consensus Amount of pledge
    first_candidate_info = client.ppos.getCandidateInfo(node.node_id)
    first_pledge_amount = first_candidate_info['Ret']['Released']
    # view block_reward
    log.info("block: {}".format(node.eth.blockNumber))
    block_reward, staking_reward = client.economic.get_current_year_reward(node)
    log.info("block_reward: {} staking_reward: {}".format(block_reward, staking_reward))
    # Get governable parameters
    first_slash_blocks = get_governable_parameter_value(client, 'slashBlocksReward')
    return first_pledge_amount, block_reward, first_slash_blocks


def verify_changed_parameters(client_con_list_obj, first_pledge_amount, block_reward, slash_blocks):
    # Verify changed parameters
    candidate_info = pledge_punishment(client_con_list_obj)
    second_pledge_amount = candidate_info['Ret']['Released']
    punishment_amonut = int(Decimal(str(block_reward)) * Decimal(str(slash_blocks)))
    if punishment_amonut < first_pledge_amount:
        assert second_pledge_amount == first_pledge_amount - punishment_amonut, "ErrMsg:Consensus Amount of pledge {}".format(
            second_pledge_amount)
    else:
        assert second_pledge_amount == 0, "ErrMsg:Consensus Amount of pledge {}".format(second_pledge_amount)


@pytest.mark.P1
@pytest.mark.parametrize('mark', [False, True])
def test_PIP_PVF_001_002(clients_consensus, mark, reset_environment):
    """
    PIP_PVF_001:治理修改低0出块率扣除验证人自有质押金块数投票失败
    PIP_PVF_002:理修改低0出块率扣除验证人自有质押金块数成功处于未生效期
    :param clients_consensus:
    :return:
    """
    index = 0
    first_client = clients_consensus[index]
    # get pledge amount1 and block reward
    first_pledge_amount, block_reward, first_slash_blocks = information_before_slash_blocks(first_client)
    # create Parametric proposal
    param_governance_verify_before_endblock(first_client, 'slashing', 'slashBlocksReward', '0', mark)
    # Get governable parameters again
    second_slash_blocks = get_governable_parameter_value(first_client, 'slashBlocksReward')
    assert first_slash_blocks == second_slash_blocks, "ErrMsg:slash blocks {}".format(second_slash_blocks)
    # Verify changed parameters
    verify_changed_parameters(clients_consensus, first_pledge_amount, block_reward, first_slash_blocks)


@pytest.mark.P1
def test_PIP_PVF_003(clients_consensus, reset_environment):
    """
    治理修改低0出块率扣除验证人自有质押金块数成功处于已生效期
    :param client_con_list_obj:
    :param reset_environment:
    :return:
    """
    first_index = 0
    first_client = clients_consensus[first_index]
    log.info("当前连接节点：{}".format(first_client.node.node_mark))
    node = first_client.node
    # get pledge amount1 and block reward
    first_pledge_amount, block_reward, first_slash_blocks = information_before_slash_blocks(first_client)
    # create Parametric proposal
    param_governance_verify(first_client, 'slashing', 'slashBlocksReward', '0')
    log.info("Current block height: {}".format(node.eth.blockNumber))
    # Get governable parameters again
    second_slash_blocks = get_governable_parameter_value(first_client, 'slashBlocksReward')
    assert second_slash_blocks == '0', "ErrMsg:Change parameters {}".format(second_slash_blocks)
    # Verify changed parameters
    verify_changed_parameters(clients_consensus, first_pledge_amount, block_reward, second_slash_blocks)


@pytest.mark.P1
def test_PIP_PVF_004(client_consensus, clients_noconsensus, reset_environment):
    """
    治理修改低0出块率扣除验证人自有质押金块数成功扣除区块奖励块数60100-自由金额质押
    :param client_consensus_obj:
    :param client_noc_list_obj:
    :param reset_environment:
    :return:
    """
    consensus_client = client_consensus
    log.info("Current connection consensus node".format(consensus_client.node.node_mark))
    first_index = 0
    second_client = clients_noconsensus[first_index]
    log.info("Current connection non-consensus node：{}".format(second_client.node.node_mark))
    economic = consensus_client.economic
    node = consensus_client.node
    change_parameter_value = '60100'
    # get pledge amount1 and block reward
    first_pledge_amount, block_reward, first_slash_blocks = information_before_slash_blocks(consensus_client)
    # create Parametric proposal
    param_governance_verify(consensus_client, 'slashing', 'slashBlocksReward', change_parameter_value)
    log.info("Current block height: {}".format(node.eth.blockNumber))
    # Get governable parameters again
    second_slash_blocks = get_governable_parameter_value(consensus_client, 'slashBlocksReward')
    assert second_slash_blocks == change_parameter_value, "ErrMsg:Change parameters {}".format(second_slash_blocks)
    # create account
    address, _ = economic.account.generate_account(node.web3, von_amount(economic.create_staking_limit, 2))
    # create staking
    result = second_client.staking.create_staking(0, address, address)
    assert_code(result, 0)
    # wait settlement block
    economic.wait_settlement_blocknum(node)
    for i in range(4):
        result = check_node_in_list(node.node_id, consensus_client.ppos.getValidatorList)
        log.info("Current node in consensus list status：{}".format(result))
        if result:
            # Verify changed parameters
            verify_changed_parameters(clients_noconsensus, first_pledge_amount, block_reward, second_slash_blocks)
            break
        else:
            # wait consensus block
            economic.wait_consensus_blocknum(node)


@pytest.mark.P1
def test_PIP_PVF_005(clients_new_node, reset_environment):
    """
    治理修改低出块率扣除验证人自有质押金比例扣除区块奖励块数60100-锁仓金额质押
    :param client_new_node_obj_list:
    :param clients_new_node:
    :param reset_environment:
    :return:
    """
    first_index = 0
    first_client = clients_new_node[first_index]
    log.info("当前连接节点：{}".format(first_client.node.node_mark))
    economic = first_client.economic
    node = first_client.node
    change_parameter_value = '60100'
    # get pledge amount1 and block reward
    first_pledge_amount, block_reward, first_slash_blocks = information_before_slash_blocks(first_client)
    # create Parametric proposal
    param_governance_verify(first_client, 'slashing', 'slashBlocksReward', change_parameter_value)
    log.info("Current block height: {}".format(node.eth.blockNumber))
    # Get governable parameters
    second_slash_blocks = get_governable_parameter_value(first_client, 'slashBlocksReward')
    assert second_slash_blocks == change_parameter_value, "ErrMsg:Change parameters {}".format(second_slash_blocks)
    # create account
    address, _ = economic.account.generate_account(node.web3, von_amount(economic.create_staking_limit, 2))
    address1, _ = economic.account.generate_account(node.web3, node.web3.toWei(1000, 'ether'))
    # Create restricting plan
    plan = [{'Epoch': 1, 'Amount': economic.create_staking_limit}]
    result = first_client.restricting.createRestrictingPlan(address1, plan, address)
    assert_code(result, 0)
    # create staking
    result = first_client.staking.create_staking(1, address1, address1)
    assert_code(result, 0)
    # wait settlement block
    economic.wait_settlement_blocknum(node)
    for i in range(4):
        result = check_node_in_list(node.node_id, first_client.ppos.getValidatorList)
        # log.info("Current node in consensus list status：{}".format(result))
        if result:
            # Verify changed parameters
            verify_changed_parameters(clients_new_node, first_pledge_amount, block_reward, second_slash_blocks)
            break
        else:
            # wait consensus block
            economic.wait_consensus_blocknum(node)


def adjust_initial_parameters(new_genesis_env):
    # Change configuration parameters
    genesis = from_dict(data_class=Genesis, data=new_genesis_env.genesis_config)
    genesis.economicModel.staking.unStakeFreezeDuration = 3
    genesis.economicModel.slashing.maxEvidenceAge = 2
    new_file = new_genesis_env.cfg.env_tmp + "/genesis.json"
    genesis.to_file(new_file)
    new_genesis_env.deploy_all(new_file)


@pytest.mark.P1
def test_PIP_PVF_006(new_genesis_env, clients_consensus):
    """
    治理修改区块双签-证据有效期投票失败
    :param client_con_list_obj:
    :return:
    """
    # Change configuration parameters
    adjust_initial_parameters(new_genesis_env)
    # view Parameter value before treatment
    first_slash_blocks = get_governable_parameter_value(clients_consensus[0], 'maxEvidenceAge')
    # create Parametric proposal
    param_governance_verify_before_endblock(clients_consensus[0], 'slashing', 'maxEvidenceAge', '1', False)
    # view Parameter value before treatment again
    second_slash_blocks = get_governable_parameter_value(clients_consensus[0], 'maxEvidenceAge')
    assert second_slash_blocks == first_slash_blocks, "ErrMsg:Parameter value after treatment {}".format(second_slash_blocks)
    report_address, _ = clients_consensus[0].economic.account.generate_account(clients_consensus[0].node.web3,
                                                                               clients_consensus[0].node.web3.toWei(
                                                                                     1000, 'ether'))
    # wait consensus block
    clients_consensus[0].economic.wait_consensus_blocknum(clients_consensus[0].node)
    # Verify changed parameters
    effective_block = clients_consensus[0].economic.get_front_settlement_switchpoint(clients_consensus[0].node, 2)
    if effective_block < 41:
        log.info("Current block: {}".format(clients_consensus[0].node.eth.blockNumber))
        effective_block = 41
    log.info("Effective block height: {}".format(effective_block))
    # Report prepareblock signature
    report_information = mock_duplicate_sign(1, clients_consensus[0].node.nodekey,
                                             clients_consensus[0].node.blsprikey,
                                             effective_block)
    log.info("Report information: {}".format(report_information))
    result = clients_consensus[0].duplicatesign.reportDuplicateSign(1, report_information, report_address)
    assert_code(result, 0)


@pytest.mark.P1
def test_PIP_PVF_007(new_genesis_env, clients_consensus):
    """
    治理修改区块双签-证据有效期处于未生效期
    :param new_genesis_env:
    :param client_con_list_obj:
    :return:
    """
    # Change configuration parameters
    adjust_initial_parameters(new_genesis_env)
    # view Parameter value before treatment
    first_slash_blocks = get_governable_parameter_value(clients_consensus[0], 'maxEvidenceAge')
    # create Parametric proposal
    param_governance_verify_before_endblock(clients_consensus[0], 'slashing', 'maxEvidenceAge', '1')
    # view Parameter value before treatment again
    second_slash_blocks = get_governable_parameter_value(clients_consensus[0], 'maxEvidenceAge')
    assert second_slash_blocks == first_slash_blocks, "ErrMsg:Parameter value after treatment {}".format(second_slash_blocks)
    report_address, _ = clients_consensus[0].economic.account.generate_account(clients_consensus[0].node.web3,
                                                                               clients_consensus[0].node.web3.toWei(
                                                                                     1000, 'ether'))
    # wait consensus block
    clients_consensus[0].economic.wait_consensus_blocknum(clients_consensus[0].node)
    # Verify changed parameters
    effective_block = clients_consensus[0].economic.get_front_settlement_switchpoint(clients_consensus[0].node, 2)
    if effective_block < 41:
        log.info("Current block: {}".format(clients_consensus[0].node.eth.blockNumber))
        effective_block = 41
    log.info("Effective block height: {}".format(effective_block))
    # Report prepareblock signature
    report_information = mock_duplicate_sign(1, clients_consensus[0].node.nodekey,
                                             clients_consensus[0].node.blsprikey,
                                             effective_block)
    log.info("Report information: {}".format(report_information))
    result = clients_consensus[0].duplicatesign.reportDuplicateSign(1, report_information, report_address)
    assert_code(result, 0)


@pytest.mark.P1
def test_PIP_PVF_008(new_genesis_env, clients_consensus):
    """
    治理修改区块双签-证据有效期处于已生效期
    :param new_genesis_env:
    :param client_con_list_obj:
    :return:
    """
    # Change configuration parameters
    adjust_initial_parameters(new_genesis_env)
    # view Parameter value before treatment
    first_slash_blocks = get_governable_parameter_value(clients_consensus[0], 'maxEvidenceAge')
    # create Parametric proposal
    param_governance_verify(clients_consensus[0], 'slashing', 'maxEvidenceAge', '1')
    # view Parameter value before treatment again
    second_slash_blocks = get_governable_parameter_value(clients_consensus[0], 'maxEvidenceAge')
    assert second_slash_blocks == '1', "ErrMsg:Parameter value after treatment {}".format(second_slash_blocks)
    report_address, _ = clients_consensus[0].economic.account.generate_account(clients_consensus[0].node.web3,
                                                                               clients_consensus[0].node.web3.toWei(
                                                                                     1000, 'ether'))
    # wait consensus block
    clients_consensus[0].economic.wait_consensus_blocknum(clients_consensus[0].node)
    # Verify changed parameters
    effective_block2 = clients_consensus[0].economic.get_front_settlement_switchpoint(clients_consensus[0].node,
                                                                                      int(second_slash_blocks))
    log.info("Effective2 block height: {}".format(effective_block2))
    # wait consensus block
    clients_consensus[0].economic.wait_consensus_blocknum(clients_consensus[0].node)
    # Report2 prepareblock signature
    report_information = mock_duplicate_sign(1, clients_consensus[0].node.nodekey,
                                             clients_consensus[0].node.blsprikey,
                                             effective_block2)
    log.info("Report information: {}".format(report_information))
    result = clients_consensus[0].duplicatesign.reportDuplicateSign(1, report_information, report_address)
    assert_code(result, 0)


@pytest.mark.P1
def test_PIP_PVF_009(new_genesis_env, clients_consensus):
    """
    治理修改区块双签-证据有效期（节点质押退回锁定周期-1）
    :param new_genesis_env:
    :param client_con_list_obj:
    :return:
    """
    # Change configuration parameters
    genesis = from_dict(data_class=Genesis, data=new_genesis_env.genesis_config)
    genesis.economicModel.staking.unStakeFreezeDuration = 3
    genesis.economicModel.slashing.maxEvidenceAge = 1
    new_file = new_genesis_env.cfg.env_tmp + "/genesis.json"
    genesis.to_file(new_file)
    new_genesis_env.deploy_all(new_file)
    # view Parameter value before treatment
    first_slash_blocks = get_governable_parameter_value(clients_consensus[0], 'maxEvidenceAge')
    # create Parametric proposal
    param_governance_verify(clients_consensus[0], 'slashing', 'maxEvidenceAge', '2')
    # view Parameter value before treatment again
    second_slash_blocks = get_governable_parameter_value(clients_consensus[0], 'maxEvidenceAge')
    assert second_slash_blocks == '2', "ErrMsg:Parameter value after treatment {}".format(second_slash_blocks)
    report_address, _ = clients_consensus[0].economic.account.generate_account(clients_consensus[0].node.web3,
                                                                               clients_consensus[0].node.web3.toWei(
                                                                                     1000, 'ether'))
    # wait consensus block
    clients_consensus[0].economic.wait_consensus_blocknum(clients_consensus[0].node)
    # Verify changed parameters
    effective_block1 = clients_consensus[0].economic.get_front_settlement_switchpoint(clients_consensus[0].node,
                                                                                      int(first_slash_blocks))
    effective_block2 = clients_consensus[0].economic.get_front_settlement_switchpoint(clients_consensus[0].node,
                                                                                      int(second_slash_blocks))
    log.info("Effective1 block height: {}".format(effective_block1))
    log.info("Effective2 block height: {}".format(effective_block2))
    # Report1 prepareblock signature
    report_information = mock_duplicate_sign(1, clients_consensus[0].node.nodekey,
                                             clients_consensus[0].node.blsprikey,
                                             effective_block1)
    log.info("Report information: {}".format(report_information))
    result = clients_consensus[0].duplicatesign.reportDuplicateSign(1, report_information, report_address)
    assert_code(result, 0)
    # Report2 prepareblock signature
    report_information = mock_duplicate_sign(1, clients_consensus[0].node.nodekey,
                                             clients_consensus[0].node.blsprikey,
                                             effective_block2)
    log.info("Report information: {}".format(report_information))
    result = clients_consensus[0].duplicatesign.reportDuplicateSign(1, report_information, report_address)
    assert_code(result, 303000)


@pytest.mark.P1
def test_PIP_PVF_010(new_genesis_env, clients_consensus):
    """
    治理修改区块双签-证据有效期（超出有效期）
    :param new_genesis_env:
    :param client_con_list_obj:
    :return:
    """
    # Change configuration parameters
    adjust_initial_parameters(new_genesis_env)
    # view Parameter value before treatment
    first_slash_blocks = get_governable_parameter_value(clients_consensus[0], 'maxEvidenceAge')
    assert first_slash_blocks == '2', "ErrMsg:Parameter value before treatment {}".format(first_slash_blocks)
    # create Parametric proposal
    param_governance_verify(clients_consensus[0], 'slashing', 'maxEvidenceAge', '1')
    # view Parameter value before treatment again
    second_slash_blocks = get_governable_parameter_value(clients_consensus[0], 'maxEvidenceAge')
    assert second_slash_blocks == '1', "ErrMsg:Parameter value after treatment {}".format(second_slash_blocks)
    report_address, _ = clients_consensus[0].economic.account.generate_account(clients_consensus[0].node.web3,
                                                                               clients_consensus[0].node.web3.toWei(
                                                                                     1000, 'ether'))
    # wait consensus block
    clients_consensus[0].economic.wait_consensus_blocknum(clients_consensus[0].node)
    # Verify changed parameters
    effective_block1 = clients_consensus[0].economic.get_front_settlement_switchpoint(clients_consensus[0].node,
                                                                                      int(first_slash_blocks))
    log.info("Effective1 block height: {}".format(effective_block1))
    # Report1 prepareblock signature
    report_information = mock_duplicate_sign(1, clients_consensus[0].node.nodekey,
                                             clients_consensus[0].node.blsprikey,
                                             effective_block1)
    log.info("Report information: {}".format(report_information))
    result = clients_consensus[0].duplicatesign.reportDuplicateSign(1, report_information, report_address)
    assert_code(result, 303003)


def information_before_penalty_ratio(client):
    # view Pledge amount
    first_candidate_info = client.ppos.getCandidateInfo(client.node.node_id)
    first_pledge_amount = first_candidate_info['Ret']['Released']
    # view Parameter value before treatment
    penalty_ratio1 = get_governable_parameter_value(client, 'slashFractionDuplicateSign')
    return first_pledge_amount, penalty_ratio1


def duplicate_sign(client, report_address, report_block):
    if report_block < 41:
        report_block = 41
    # Report1 prepareblock signature
    report_information = mock_duplicate_sign(1, client.node.nodekey,
                                             client.node.blsprikey,
                                             report_block)
    log.info("Report information: {}".format(report_information))
    result = client.duplicatesign.reportDuplicateSign(1, report_information, report_address)
    assert_code(result, 0)


def assret_penalty_amount(client_con_list_obj, first_pledge_amount, penalty_ratio=None):
    # view Pledge amount after punishment
    proportion_reward, incentive_pool_reward = client_con_list_obj[1].economic.get_report_reward(first_pledge_amount,
                                                                                                 penalty_ratio)
    log.info("Whistleblower benefits：{} Incentive pool income：{}".format(proportion_reward, incentive_pool_reward))
    # view Pledge amount again
    candidate_info2 = client_con_list_obj[1].ppos.getCandidateInfo(client_con_list_obj[0].node.node_id)
    second_pledge_amount = candidate_info2['Ret']['Released']
    assert second_pledge_amount == first_pledge_amount - proportion_reward - incentive_pool_reward, "ErrMsg:Pledge amount {}".format(
        second_pledge_amount)


@pytest.mark.P1
def test_PIP_PVF_011(clients_consensus, reset_environment):
    """
    治理修改区块双签-最高处罚比例投票失败
    :param client_con_list_obj:
    :param reset_environment:
    :return:
    """
    # view Pledge amount and Parameter value before treatment
    first_pledge_amount, penalty_ratio1 = information_before_penalty_ratio(clients_consensus[0])
    # create Parametric proposal
    param_governance_verify_before_endblock(clients_consensus[0], 'slashing', 'slashFractionDuplicateSign', '1000',
                                            False)
    # view Parameter value after treatment again
    penalty_ratio2 = get_governable_parameter_value(clients_consensus[0], 'slashFractionDuplicateSign')
    assert penalty_ratio1 == penalty_ratio2, "ErrMsg:Parameter value after treatment {}".format(penalty_ratio2)
    # wait consensus block
    clients_consensus[0].economic.wait_consensus_blocknum(clients_consensus[0].node)
    # create account
    report_address, _ = clients_consensus[0].economic.account.generate_account(clients_consensus[0].node.web3,
                                                                               clients_consensus[0].node.web3.toWei(
                                                                                     1000, 'ether'))
    # Verify changed parameters
    current_block = clients_consensus[0].node.eth.blockNumber
    log.info("Current block height: {}".format(current_block))
    # Verify changed parameters report
    duplicate_sign(clients_consensus[0], report_address, current_block)
    # assret penalty amount
    assret_penalty_amount(clients_consensus, first_pledge_amount)


@pytest.mark.P1
def test_PIP_PVF_012(clients_consensus, reset_environment):
    """
    治理修改区块双签-最高处罚比例处于未生效期
    :param client_con_list_obj:
    :param reset_environment:
    :return:
    """
    # view Pledge amount and Parameter value before treatment
    first_pledge_amount, penalty_ratio1 = information_before_penalty_ratio(clients_consensus[0])
    # create Parametric proposal
    param_governance_verify_before_endblock(clients_consensus[0], 'slashing', 'slashFractionDuplicateSign', '1000')
    # view Parameter value before treatment again
    penalty_ratio2 = get_governable_parameter_value(clients_consensus[0], 'slashFractionDuplicateSign')
    assert penalty_ratio1 == penalty_ratio2, "ErrMsg:Parameter value after treatment {}".format(penalty_ratio2)
    # wait consensus block
    clients_consensus[0].economic.wait_consensus_blocknum(clients_consensus[0].node)
    # create account
    report_address, _ = clients_consensus[0].economic.account.generate_account(clients_consensus[0].node.web3,
                                                                               clients_consensus[0].node.web3.toWei(
                                                                                     1000, 'ether'))
    # Verify changed parameters
    current_block = clients_consensus[0].node.eth.blockNumber
    log.info("Current block height: {}".format(current_block))
    # Verify changed parameters report
    duplicate_sign(clients_consensus[0], report_address, current_block)
    # assret penalty amount
    assret_penalty_amount(clients_consensus, first_pledge_amount)


@pytest.mark.P1
def test_PIP_PVF_013(clients_consensus, reset_environment):
    """
    治理修改区块双签-最高处罚比例处于已生效期
    :param client_con_list_obj:
    :param reset_environment:
    :return:
    """
    # view Pledge amount and Parameter value before treatment
    first_pledge_amount, penalty_ratio1 = information_before_penalty_ratio(clients_consensus[0])
    # create Parametric proposal
    param_governance_verify(clients_consensus[0], 'slashing', 'slashFractionDuplicateSign', '1000')
    # view Parameter value before treatment again
    penalty_ratio2 = get_governable_parameter_value(clients_consensus[0], 'slashFractionDuplicateSign')
    assert penalty_ratio2 == '1000', "ErrMsg:Parameter value after treatment {}".format(penalty_ratio2)
    # create account
    report_address, _ = clients_consensus[0].economic.account.generate_account(clients_consensus[0].node.web3,
                                                                               clients_consensus[0].node.web3.toWei(
                                                                                     1000, 'ether'))
    # Verify changed parameters
    current_block = clients_consensus[0].node.eth.blockNumber
    log.info("Current block height: {}".format(current_block))
    # Verify changed parameters report
    duplicate_sign(clients_consensus[0], report_address, current_block)
    # assret penalty amount
    assret_penalty_amount(clients_consensus, first_pledge_amount, 1000)


@pytest.mark.P1
def test_PIP_PVF_014(clients_consensus, reset_environment):
    """
    治理修改区块双签-最高处罚比例为10000‱
    :param client_con_list_obj:
    :param reset_environment:
    :return:
    """
    # view Pledge amount and Parameter value before treatment
    first_pledge_amount, penalty_ratio1 = information_before_penalty_ratio(clients_consensus[0])
    # create Parametric proposal
    param_governance_verify(clients_consensus[0], 'slashing', 'slashFractionDuplicateSign', '10000')
    # view Parameter value before treatment again
    penalty_ratio2 = get_governable_parameter_value(clients_consensus[0], 'slashFractionDuplicateSign')
    assert penalty_ratio2 == '10000', "ErrMsg:Parameter value after treatment {}".format(penalty_ratio2)
    # create account
    report_address, _ = clients_consensus[0].economic.account.generate_account(clients_consensus[0].node.web3,
                                                                               clients_consensus[0].node.web3.toWei(
                                                                                     1000, 'ether'))
    # Verify changed parameters
    current_block = clients_consensus[0].node.eth.blockNumber
    log.info("Current block height: {}".format(current_block))
    # Verify changed parameters report
    duplicate_sign(clients_consensus[0], report_address, current_block)
    # assret penalty amount
    assret_penalty_amount(clients_consensus, first_pledge_amount, 10000)


@pytest.mark.P1
def test_PIP_PVF_015(clients_consensus, reset_environment):
    """
    治理修改区块双签-最高处罚比例为1‱
    :param client_con_list_obj:
    :param reset_environment:
    :return:
    """
    # view Pledge amount and Parameter value before treatment
    first_pledge_amount, penalty_ratio1 = information_before_penalty_ratio(clients_consensus[0])
    # create Parametric proposal
    param_governance_verify(clients_consensus[0], 'slashing', 'slashFractionDuplicateSign', '1')
    # view Parameter value before treatment again
    penalty_ratio2 = get_governable_parameter_value(clients_consensus[0], 'slashFractionDuplicateSign')
    assert penalty_ratio2 == '1', "ErrMsg:Parameter value after treatment {}".format(penalty_ratio2)
    # create account
    report_address, _ = clients_consensus[0].economic.account.generate_account(clients_consensus[0].node.web3,
                                                                               clients_consensus[0].node.web3.toWei(
                                                                                     1000, 'ether'))
    # Verify changed parameters
    current_block = clients_consensus[0].node.eth.blockNumber
    log.info("Current block height: {}".format(current_block))
    # Verify changed parameters report
    duplicate_sign(clients_consensus[0], report_address, current_block)
    # assret penalty amount
    assret_penalty_amount(clients_consensus, first_pledge_amount, 1)


def information_before_report_reward(client):
    # view Pledge amount
    first_candidate_info = client.ppos.getCandidateInfo(client.node.node_id)
    first_pledge_amount = first_candidate_info['Ret']['Released']
    # view Parameter value before treatment
    report_reward1 = get_governable_parameter_value(client, 'duplicateSignReportReward')
    return first_pledge_amount, report_reward1


def get_account_amount(client):
    # create report account
    report_address, _ = client.economic.account.generate_account(client.node.web3, client.node.web3.toWei(
        1000, 'ether'))
    # view report amount
    report_amount1 = client.node.eth.getBalance(report_address)
    # view Incentive pool account
    incentive_pool_account1 = client.node.eth.getBalance(EconomicConfig.INCENTIVEPOOL_ADDRESS)
    return report_address, report_amount1, incentive_pool_account1


def asster_income_account_amount(client, report_amount1, incentive_pool_account1, report_address, proportion_reward,
                                 incentive_pool_reward):
    # view report amount
    report_amount2 = client.node.eth.getBalance(report_address)
    # view Incentive pool account
    incentive_pool_account2 = client.node.eth.getBalance(EconomicConfig.INCENTIVEPOOL_ADDRESS)
    # asster amount reward
    log.info("report_amount1 {} ,proportion_reward {} , report_amount2 {}".format(report_amount1, proportion_reward,
                                                                                  report_amount2))
    assert report_amount1 + proportion_reward - report_amount2 < client.node.web3.toWei(1,
                                                                                            'ether'), "ErrMsg:report amount {}".format(
        report_amount2)
    log.info("incentive_pool_account2 {} ,incentive_pool_account1 {} , incentive_pool_reward {}".format(
        incentive_pool_account2, incentive_pool_account1, incentive_pool_reward))
    assert incentive_pool_account2 == incentive_pool_account1 + incentive_pool_reward + (
        report_amount1 + proportion_reward - report_amount2), "ErrMsg:Incentive pool account {}".format(
        incentive_pool_account2)


@pytest.mark.P1
def test_PIP_PVF_016(clients_consensus, reset_environment):
    """
    治理修改区块双签-举报奖励比例投票失败
    :param client_con_list_obj:
    :param reset_environment:
    :return:
    """
    # get first_pledge_amount report_amount1 incentive_pool_account1 report_reward1
    first_pledge_amount, report_reward1 = information_before_report_reward(clients_consensus[0])
    # create Parametric proposal
    param_governance_verify_before_endblock(clients_consensus[0], 'slashing', 'duplicateSignReportReward', '60',
                                            False)
    # view Parameter value after treatment
    report_reward2 = get_governable_parameter_value(clients_consensus[0], 'duplicateSignReportReward')
    assert report_reward1 == report_reward2, "ErrMsg:Parameter value after treatment {}".format(report_reward2)
    # wait consensus block
    clients_consensus[0].economic.wait_consensus_blocknum(clients_consensus[0].node)
    # get account amount
    report_address, report_amount1, incentive_pool_account1 = get_account_amount(clients_consensus[0])
    # Verify changed parameters
    current_block = clients_consensus[0].node.eth.blockNumber
    log.info("Current block height: {}".format(current_block))
    # Verify changed parameters report
    duplicate_sign(clients_consensus[0], report_address, current_block)
    # view Pledge amount after punishment
    proportion_reward, incentive_pool_reward = clients_consensus[1].economic.get_report_reward(first_pledge_amount)
    log.info("Whistleblower benefits：{} Incentive pool income：{}".format(proportion_reward, incentive_pool_reward))
    # asster account amount
    asster_income_account_amount(clients_consensus[1], report_amount1, incentive_pool_account1,
                                 report_address, proportion_reward, incentive_pool_reward)


@pytest.mark.P1
def test_PIP_PVF_017(clients_consensus, reset_environment):
    """
    治理修改区块双签-举报奖励比例处于未生效期
    :param client_con_list_obj:
    :param reset_environment:
    :return:
    """
    # get first_pledge_amount report_amount1 incentive_pool_account1 report_reward1
    first_pledge_amount, report_reward1 = information_before_report_reward(clients_consensus[0])
    # create Parametric proposal
    param_governance_verify_before_endblock(clients_consensus[0], 'slashing', 'duplicateSignReportReward', '60')
    # view Parameter value after treatment
    report_reward2 = get_governable_parameter_value(clients_consensus[0], 'duplicateSignReportReward')
    assert report_reward1 == report_reward2, "ErrMsg:Parameter value after treatment {}".format(report_reward2)
    # wait consensus block
    clients_consensus[0].economic.wait_consensus_blocknum(clients_consensus[0].node)
    # get account amount
    report_address, report_amount1, incentive_pool_account1 = get_account_amount(clients_consensus[0])
    # Verify changed parameters
    current_block = clients_consensus[0].node.eth.blockNumber
    log.info("Current block height: {}".format(current_block))
    # Verify changed parameters report
    duplicate_sign(clients_consensus[0], report_address, current_block)
    # view Pledge amount after punishment
    proportion_reward, incentive_pool_reward = clients_consensus[1].economic.get_report_reward(first_pledge_amount)
    log.info("Whistleblower benefits：{} Incentive pool income：{}".format(proportion_reward, incentive_pool_reward))
    # asster account amount
    asster_income_account_amount(clients_consensus[1], report_amount1, incentive_pool_account1,
                                 report_address, proportion_reward, incentive_pool_reward)


@pytest.mark.P1
def test_PIP_PVF_018(clients_consensus, reset_environment):
    """
    治理修改区块双签-举报奖励比例处于已生效期
    :param client_con_list_obj:
    :param reset_environment:
    :return:
    """
    # get first_pledge_amount report_amount1 incentive_pool_account1 report_reward1
    first_pledge_amount, report_reward1 = information_before_report_reward(clients_consensus[0])
    # create Parametric proposal
    param_governance_verify(clients_consensus[0], 'slashing', 'duplicateSignReportReward', '60')
    # view Parameter value after treatment
    report_reward2 = get_governable_parameter_value(clients_consensus[0], 'duplicateSignReportReward')
    assert report_reward2 == '60', "ErrMsg:Parameter value after treatment {}".format(report_reward2)
    # get account amount
    report_address, report_amount1, incentive_pool_account1 = get_account_amount(clients_consensus[0])
    # Verify changed parameters
    current_block = clients_consensus[0].node.eth.blockNumber
    log.info("Current block height: {}".format(current_block))
    # Verify changed parameters report
    duplicate_sign(clients_consensus[0], report_address, current_block)
    # view Pledge amount after punishment
    proportion_reward, incentive_pool_reward = clients_consensus[1].economic.get_report_reward(first_pledge_amount, None,
                                                                                               60)
    log.info("Whistleblower benefits：{} Incentive pool income：{}".format(proportion_reward, incentive_pool_reward))
    # asster account amount
    asster_income_account_amount(clients_consensus[1], report_amount1, incentive_pool_account1,
                                 report_address, proportion_reward, incentive_pool_reward)


@pytest.mark.P1
def test_PIP_PVF_019(clients_consensus, reset_environment):
    """
    治理修改区块双签-举报奖励比例为80%
    :param client_con_list_obj:
    :param reset_environment:
    :return:
    """
    # get first_pledge_amount report_reward1
    first_pledge_amount, report_reward1 = information_before_report_reward(clients_consensus[0])
    # create Parametric proposal
    param_governance_verify(clients_consensus[0], 'slashing', 'duplicateSignReportReward', '80')
    # view Parameter value after treatment
    report_reward2 = get_governable_parameter_value(clients_consensus[0], 'duplicateSignReportReward')
    assert report_reward2 == '80', "ErrMsg:Parameter value after treatment {}".format(report_reward2)
    # get account amount
    report_address, report_amount1, incentive_pool_account1 = get_account_amount(clients_consensus[0])
    # Verify changed parameters
    current_block = clients_consensus[0].node.eth.blockNumber
    log.info("Current block height: {}".format(current_block))
    # Verify changed parameters report
    duplicate_sign(clients_consensus[0], report_address, current_block)
    # view Pledge amount after punishment
    proportion_reward, incentive_pool_reward = clients_consensus[1].economic.get_report_reward(first_pledge_amount, None,
                                                                                               80)
    log.info("Whistleblower benefits：{} Incentive pool income：{}".format(proportion_reward, incentive_pool_reward))
    # asster account amount
    asster_income_account_amount(clients_consensus[1], report_amount1, incentive_pool_account1,
                                 report_address, proportion_reward, incentive_pool_reward)


@pytest.mark.P1
def test_PIP_PVF_020(clients_consensus, reset_environment):
    """
    治理修改区块双签-举报奖励比例为1%
    :param client_con_list_obj:
    :param reset_environment:
    :return:
    """
    # get first_pledge_amount report_reward1
    first_pledge_amount, report_reward1 = information_before_report_reward(clients_consensus[0])
    # create Parametric proposal
    param_governance_verify(clients_consensus[0], 'slashing', 'duplicateSignReportReward', '1')
    # view Parameter value after treatment
    report_reward2 = get_governable_parameter_value(clients_consensus[0], 'duplicateSignReportReward')
    assert report_reward2 == '1', "ErrMsg:Parameter value after treatment {}".format(report_reward2)
    # get account amount
    report_address, report_amount1, incentive_pool_account1 = get_account_amount(clients_consensus[0])
    # Verify changed parameters
    current_block = clients_consensus[0].node.eth.blockNumber
    log.info("Current block height: {}".format(current_block))
    # Verify changed parameters report
    duplicate_sign(clients_consensus[0], report_address, current_block)
    # view Pledge amount after punishment
    proportion_reward, incentive_pool_reward = clients_consensus[1].economic.get_report_reward(first_pledge_amount, None,
                                                                                               1)
    log.info("Whistleblower benefits：{} Incentive pool income：{}".format(proportion_reward, incentive_pool_reward))
    # asster account amount
    asster_income_account_amount(clients_consensus[1], report_amount1, incentive_pool_account1,
                                 report_address, proportion_reward, incentive_pool_reward)


#
# def transaction(client, nonce, from_address, to_address, value):
#     account = client.economic.account.accounts[from_address]
#     tmp_to_address = Web3.toChecksumAddress(to_address)
#     tmp_from_address = Web3.toChecksumAddress(from_address)
#
#     transaction_dict = {
#         "to": tmp_to_address,
#         "gasPrice": client.node.eth.gasPrice,
#         "gas": 21000,
#         "nonce": nonce,
#         "data": "",
#         "chainId": client.node.chain_id,
#         "value": value,
#         'from': tmp_from_address,
#     }
#     signedTransactionDict = client.node.eth.account.signTransaction(
#         transaction_dict, account['prikey']
#     )
#
#     data = signedTransactionDict.rawTransaction
#     result = HexBytes(client.node.eth.sendRawTransaction(data)).hex()
#    nonce = clients_consensus[0].node.eth.getTransactionCount(
#         clients_consensus[0].economic.env.account.account_with_money['address'])
#     from_address = clients_consensus[0].economic.env.account.account_with_money['address']
#     to_address, _ = clients_consensus[0].economic.account.generate_account(clients_consensus[0].node.web3, 0)
#     for i in range(0, 200):
#         # Transfer transaction
#         transaction(clients_consensus[0], nonce, from_address, to_address, 10)
#         nonce = nonce + 1
#     time.sleep(15)
#     end_block = clients_consensus[0].node.block_number
#     max_tx = {"block_num": 0, "tx_num": 0}
#     for i in range(1, end_block + 1):
#         tx_num = clients_consensus[0].node.eth.getBlockTransactionCount(i)
#         if tx_num > max_tx["tx_num"]:
#             max_tx = {"block_num": i, "tx_num": tx_num}
#     block_info = clients_consensus[0].node.eth.getBlock(max_tx["block_num"])
#     print(max_tx['tx_num'])
#     print(block_info['gasLimit'])
#     print(block_info['gasUsed'])
#     print("block_info", block_info)
#
# #
@pytest.mark.P1
def test_PIP_MG_001(clients_consensus, reset_environment):
    """
    治理修改默认每个区块的最大Gas 投票失败
    :param client_con_list_obj:
    :param reset_environment:
    :return:
    """
    # # Change configuration parameters
    # genesis = from_dict(data_class=Genesis, data=new_genesis_env.genesis_config)
    # genesis.config.cbft.period = 50000
    # genesis.EconomicModel.Common.MaxEpochMinutes = 14
    # genesis.EconomicModel.Common.AdditionalCycleTime = 55
    # new_file = new_genesis_env.cfg.env_tmp + "/genesis.json"
    # genesis.to_file(new_file)
    # new_genesis_env.deploy_all(new_file)
    # clients_consensus[0].economic.env.deploy_all()

    # view Parameter value before treatment
    max_gas_limit1 = get_governable_parameter_value(clients_consensus[0], 'maxBlockGasLimit')
    # create Parametric proposal
    block = param_governance_verify_before_endblock(clients_consensus[0], 'block', 'maxBlockGasLimit', '4712389',
                                                    False)
    # view Parameter value after treatment
    max_gas_limit2 = get_governable_parameter_value(clients_consensus[0], 'maxBlockGasLimit')
    # wait block
    wait_block_number(clients_consensus[0].node, block)
    assert max_gas_limit2 == max_gas_limit1, "ErrMsg:Parameter value after treatment {}".format(max_gas_limit2)


@pytest.mark.P1
def test_PIP_MG_002(clients_consensus, reset_environment):
    """
    治理修改默认每个区块的最大Gas 处于未生效期
    :param client_con_list_obj:
    :param reset_environment:
    :return:
    """
    # view Parameter value before treatment
    max_gas_limit1 = get_governable_parameter_value(clients_consensus[0], 'maxBlockGasLimit')
    # create Parametric proposal
    param_governance_verify_before_endblock(clients_consensus[0], 'block', 'maxBlockGasLimit', '4712389')
    # view Parameter value after treatment
    max_gas_limit2 = get_governable_parameter_value(clients_consensus[0], 'maxBlockGasLimit')

    assert max_gas_limit2 == max_gas_limit1, "ErrMsg:Parameter value after treatment {}".format(max_gas_limit2)


@pytest.mark.P1
def test_PIP_MG_003(clients_consensus, reset_environment):
    """
    治理修改默认每个区块的最大Gas 处于已生效期
    :param client_con_list_obj:
    :param reset_environment:
    :return:
    """
    # view Parameter value before treatment
    max_gas_limit1 = get_governable_parameter_value(clients_consensus[0], 'maxBlockGasLimit')
    # create Parametric proposal
    param_governance_verify(clients_consensus[0], 'block', 'maxBlockGasLimit', '4712389')
    # view Parameter value after treatment
    max_gas_limit2 = get_governable_parameter_value(clients_consensus[0], 'maxBlockGasLimit')
    assert max_gas_limit2 == '4712389', "ErrMsg:Parameter value after treatment {}".format(max_gas_limit2)
