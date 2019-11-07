import json
import time

import pytest
import allure

from dacite import from_dict

from common.key import get_pub_key, mock_duplicate_sign
from common.log import log
from client_sdk_python import Web3
from decimal import Decimal

from tests.conftest import param_governance_verify, param_governance_verify_before_endblock
from tests.lib import EconomicConfig, Genesis, StakingConfig, Staking, check_node_in_list, assert_code, von_amount, \
    get_governable_parameter_value


def pledge_punishment(client_con_list_obj):
    """
    :return:
    """
    log.info("Current block height: {}".format(client_con_list_obj[1].node.eth.blockNumber))
    # stop node
    client_con_list_obj[0].node.stop()
    # Waiting for a settlement round
    client_con_list_obj[1].economic.wait_consensus_blocknum(client_con_list_obj[1].node, 2)
    log.info("Current block height: {}".format(client_con_list_obj[1].node.eth.blockNumber))
    # view verifier list
    verifier_list = client_con_list_obj[1].ppos.getVerifierList()
    log.info("verifier_list: {}".format(verifier_list))
    candidate_info = client_con_list_obj[1].ppos.getCandidateInfo(client_con_list_obj[0].node.node_id)
    log.info("Pledge node information： {}".format(candidate_info))
    return candidate_info


def Information_before_governance(client_obj):
    # view Consensus Amount of pledge
    candidate_info1 = client_obj.ppos.getCandidateInfo(client_obj.node.node_id)
    pledge_amount1 = candidate_info1['Ret']['Released']
    # view block_reward
    log.info("block: {}".format(client_obj.node.eth.blockNumber))
    block_reward, staking_reward = client_obj.economic.get_current_year_reward(
        client_obj.node)
    log.info("block_reward: {} staking_reward: {}".format(block_reward, staking_reward))
    slash_blocks = get_governable_parameter_value(client_obj, 'SlashBlocksReward')


@pytest.mark.P1
def test_PIP_PVF_001(client_con_list_obj, reset_environment):
    """
    治理修改低出块率扣除验证人自有质押金比例投票失败
    :param client_con_list_obj:
    :return:
    """
    # Initialize environment
    client_con_list_obj[0].economic.env.deploy_all()
    time.sleep(3)
    # view Consensus Amount of pledge
    candidate_info1 = client_con_list_obj[0].ppos.getCandidateInfo(client_con_list_obj[0].node.node_id)
    pledge_amount1 = candidate_info1['Ret']['Released']
    # view block_reward
    log.info("block: {}".format(client_con_list_obj[0].node.eth.blockNumber))
    block_reward, staking_reward = client_con_list_obj[0].economic.get_current_year_reward(
        client_con_list_obj[0].node)
    log.info("block_reward: {} staking_reward: {}".format(block_reward, staking_reward))
    slash_blocks = get_governable_parameter_value(client_con_list_obj[0], 'Slashing', 'SlashBlocksReward')
    # create Parametric proposal
    param_governance_verify(client_con_list_obj[0], 'Slashing', 'SlashBlocksReward', '0', False)
    # Verify changed parameters
    candidate_info2 = pledge_punishment(client_con_list_obj)
    pledge_amount2 = candidate_info2['Ret']['Released']
    punishment_amonut = int(Decimal(str(block_reward)) * Decimal(str(slash_blocks)))
    if punishment_amonut < pledge_amount1:
        assert pledge_amount2 == pledge_amount1 - punishment_amonut, "ErrMsg:Consensus Amount of pledge {}".format(
            pledge_amount2)
    else:
        assert pledge_amount2 == 0, "ErrMsg:Consensus Amount of pledge {}".format(pledge_amount2)


@pytest.mark.P1
def test_PIP_PVF_002(client_con_list_obj, reset_environment):
    """
    治理修改低出块率扣除验证人自有质押金比例处于未生效期
    :param client_con_list_obj:
    :return:
    """
    # view Consensus Amount of pledge
    candidate_info1 = client_con_list_obj[0].ppos.getCandidateInfo(client_con_list_obj[0].node.node_id)
    pledge_amount1 = candidate_info1['Ret']['Released']
    # view block_reward
    block_reward, staking_reward = client_con_list_obj[0].economic.get_current_year_reward(
        client_con_list_obj[0].node)
    log.info("block_reward: {} staking_reward: {}".format(block_reward, staking_reward))
    slash_blocks = get_governable_parameter_value(client_con_list_obj[0], 'Slashing', 'SlashBlocksReward')
    # create Parametric proposal
    End_voting_block = param_governance_verify_before_endblock(client_con_list_obj[0], 'Slashing', 'SlashBlocksReward',
                                                               '0')
    # Verify changed parameters
    candidate_info2 = pledge_punishment(client_con_list_obj)
    pledge_amount2 = candidate_info2['Ret']['Released']
    punishment_amonut = int(Decimal(str(block_reward)) * Decimal(str(slash_blocks)))
    if punishment_amonut < pledge_amount1:
        assert pledge_amount2 == pledge_amount1 - punishment_amonut, "ErrMsg:Consensus Amount of pledge {}".format(
            pledge_amount2)
    else:
        assert pledge_amount2 == 0, "ErrMsg:Consensus Amount of pledge {}".format(pledge_amount2)


@pytest.mark.P1
def test_PIP_PVF_003(client_con_list_obj, reset_environment):
    """
    治理修改低出块率扣除验证人自有质押金比例处于已生效期
    :param client_con_list_obj:
    :param reset_environment:
    :return:
    """
    # view Consensus Amount of pledge
    candidate_info1 = client_con_list_obj[0].ppos.getCandidateInfo(client_con_list_obj[0].node.node_id)
    pledge_amount1 = candidate_info1['Ret']['Released']
    # view block_reward
    block_reward, staking_reward = client_con_list_obj[0].economic.get_current_year_reward(
        client_con_list_obj[0].node)
    log.info("block_reward: {} staking_reward: {}".format(block_reward, staking_reward))
    # Get governable parameters
    slash_blocks1 = get_governable_parameter_value(client_con_list_obj[0], 'Slashing', 'SlashBlocksReward')
    # create Parametric proposal
    param_governance_verify(client_con_list_obj[0], 'Slashing', 'SlashBlocksReward', '0')
    log.info("Current block height: {}".format(client_con_list_obj[0].node.eth.blockNumber))
    # Get governable parameters
    slash_blocks2 = get_governable_parameter_value(client_con_list_obj[0], 'Slashing', 'SlashBlocksReward')
    assert slash_blocks2 == '0', "ErrMsg:Change parameters {}".format(slash_blocks2)
    # Verify changed parameters
    candidate_info2 = pledge_punishment(client_con_list_obj)
    pledge_amount2 = candidate_info2['Ret']['Released']
    punishment_amonut = int(Decimal(str(block_reward)) * Decimal(str(slash_blocks2)))
    if punishment_amonut < pledge_amount1:
        assert pledge_amount2 == pledge_amount1 - punishment_amonut, "ErrMsg:Consensus Amount of pledge {}".format(
            pledge_amount2)
    else:
        assert pledge_amount2 == 0, "ErrMsg:Consensus Amount of pledge {}".format(pledge_amount2)


@pytest.mark.P1
def test_PIP_PVF_004(client_con_list_obj, client_new_node_obj_list, reset_environment):
    """

    :param client_con_list_obj:
    :param reset_environment:
    :return:
    """
    # view Consensus Amount of pledge
    candidate_info1 = client_con_list_obj[0].ppos.getCandidateInfo(client_con_list_obj[0].node.node_id)
    pledge_amount1 = candidate_info1['Ret']['Released']
    # view block_reward
    block_reward, staking_reward = client_con_list_obj[0].economic.get_current_year_reward(
        client_con_list_obj[0].node)
    log.info("block_reward: {} staking_reward: {}".format(block_reward, staking_reward))
    # Get governable parameters
    slash_blocks1 = get_governable_parameter_value(client_con_list_obj[0], 'Slashing', 'SlashBlocksReward')
    # create Parametric proposal
    param_governance_verify(client_con_list_obj[0], 'Slashing', 'SlashBlocksReward', '60100')
    log.info("Current block height: {}".format(client_con_list_obj[0].node.eth.blockNumber))
    # Get governable parameters
    slash_blocks2 = get_governable_parameter_value(client_con_list_obj[0], 'Slashing', 'SlashBlocksReward')
    assert slash_blocks2 == '60100', "ErrMsg:Change parameters {}".format(slash_blocks2)
    # create account
    address, _ = client_con_list_obj[0].economic.account.generate_account(client_con_list_obj[0].node.web3,
                                                                          client_con_list_obj[
                                                                              0].economic.create_staking_limit * 2)
    # create staking
    result = client_new_node_obj_list[0].staking.create_staking(0, address, address)
    assert_code(result, 0)
    # wait settlement block
    log.info(client_new_node_obj_list[0].node)
    client_new_node_obj_list[0].economic.wait_settlement_blocknum(client_new_node_obj_list[0].node)
    for i in range(4):
        result = check_node_in_list(client_con_list_obj[0].node.node_id, client_con_list_obj[0].ppos.getValidatorList)
        log.info("Current node in consensus list status：{}".format(result))
        if not result:
            # Verify changed parameters
            candidate_info2 = pledge_punishment(client_new_node_obj_list)
            pledge_amount2 = candidate_info2['Ret']['Released']
            punishment_amonut = int(Decimal(str(block_reward)) * Decimal(str(slash_blocks2)))
            if punishment_amonut < pledge_amount1:
                assert pledge_amount2 == pledge_amount1 - punishment_amonut, "ErrMsg:Consensus Amount of pledge {}".format(
                    pledge_amount2)
            else:
                assert pledge_amount2 == 0, "ErrMsg:Consensus Amount of pledge {}".format(pledge_amount2)
        else:
            # wait consensus block
            client_con_list_obj[0].economic.wait_consensus_blocknum(client_con_list_obj[0].node)


@pytest.mark.P1
def test_PIP_PVF_005(client_con_list_obj, client_new_node_obj_list, reset_environment):
    """
    治理修改低出块率扣除验证人自有质押金比例扣除区块奖励块数60100-锁仓金额质押
    :param client_con_list_obj:
    :param client_noc_list_obj:
    :param reset_environment:
    :return:
    """
    # view Consensus Amount of pledge
    candidate_info1 = client_con_list_obj[0].ppos.getCandidateInfo(client_con_list_obj[0].node.node_id)
    pledge_amount1 = candidate_info1['Ret']['Released']
    # view block_reward
    block_reward, staking_reward = client_con_list_obj[0].economic.get_current_year_reward(
        client_con_list_obj[0].node)
    log.info("block_reward: {} staking_reward: {}".format(block_reward, staking_reward))
    # Get governable parameters
    slash_blocks1 = get_governable_parameter_value(client_con_list_obj[0], 'Slashing', 'SlashBlocksReward')
    # create Parametric proposal
    param_governance_verify(client_con_list_obj[0], 'Slashing', 'SlashBlocksReward', '60100')
    log.info("Current block height: {}".format(client_con_list_obj[0].node.eth.blockNumber))
    # Get governable parameters
    slash_blocks2 = get_governable_parameter_value(client_con_list_obj[0], 'Slashing', 'SlashBlocksReward')
    assert slash_blocks2 == '60100', "ErrMsg:Change parameters {}".format(slash_blocks2)
    # create account
    address, _ = client_con_list_obj[0].economic.account.generate_account(client_con_list_obj[0].node.web3,
                                                                          client_con_list_obj[
                                                                              0].economic.create_staking_limit * 2)
    address1, _ = client_con_list_obj[0].economic.account.generate_account(client_con_list_obj[0].node.web3, client_con_list_obj[0].node.web3.toWei(1000, 'ether'))
    # Create restricting plan
    plan = [{'Epoch': 1, 'Amount': client_new_node_obj_list[0].economic.create_staking_limit}]
    result = client_new_node_obj_list[0].restricting.createRestrictingPlan(address1, plan, address)
    assert_code(result, 0)
    # create staking
    result = client_new_node_obj_list[0].staking.create_staking(1, address1, address1)
    assert_code(result, 0)
    # wait settlement block
    client_new_node_obj_list[0].economic.wait_settlement_blocknum(client_new_node_obj_list[0].node)
    for i in range(4):
        result = check_node_in_list(client_con_list_obj[0].node.node_id, client_con_list_obj[0].ppos.getValidatorList)
        log.info("Current node in consensus list status：{}".format(result))
        if not result:
            # Verify changed parameters
            candidate_info2 = pledge_punishment(client_new_node_obj_list)
            pledge_amount2 = candidate_info2['Ret']['RestrictingPlan']
            punishment_amonut = int(Decimal(str(block_reward)) * Decimal(str(slash_blocks2)))
            if punishment_amonut < pledge_amount1:
                assert pledge_amount2 == pledge_amount1 - punishment_amonut, "ErrMsg:Consensus Amount of pledge {}".format(
                    pledge_amount2)
            else:
                assert pledge_amount2 == 0, "ErrMsg:Consensus Amount of pledge {}".format(pledge_amount2)
        else:
            # wait consensus block
            client_con_list_obj[0].economic.wait_consensus_blocknum(client_con_list_obj[0].node)


@pytest.mark.P1
def test_PIP_PVF_006(new_genesis_env, client_con_list_obj, reset_environment):
    """
    治理修改区块双签-证据有效期投票失败
    :param client_con_list_obj:
    :return:
    """
    # Change configuration parameters
    genesis = from_dict(data_class=Genesis, data=new_genesis_env.genesis_config)
    genesis.EconomicModel.Staking.UnStakeFreezeDuration = 3
    genesis.EconomicModel.Slashing.MaxEvidenceAge = 2
    new_file = new_genesis_env.cfg.env_tmp + "/genesis.json"
    genesis.to_file(new_file)
    new_genesis_env.deploy_all(new_file)
    # view Parameter value before treatment
    slash_blocks1 = get_governable_parameter_value(client_con_list_obj[0], 'Slashing', 'MaxEvidenceAge')
    # create Parametric proposal
    param_governance_verify(client_con_list_obj[0], 'Slashing', 'MaxEvidenceAge', '1', False)
    # view Parameter value before treatment again
    slash_blocks2 = get_governable_parameter_value(client_con_list_obj[0], 'Slashing', 'MaxEvidenceAge')
    assert slash_blocks2 == slash_blocks1, "ErrMsg:Parameter value after treatment {}".format(slash_blocks2)
    report_address, _ = client_con_list_obj[0].economic.account.generate_account(client_con_list_obj[0].node.web3, client_con_list_obj[0].node.web3.toWei(1000, 'ether'))
    # Verify changed parameters
    effective_block = client_con_list_obj[0].economic.get_front_settlement_switchpoint(client_con_list_obj[0].node,2)
    if effective_block < 41:
        log.info("Current block: {}".format(client_con_list_obj[0].node.eth.blockNumber))
        effective_block = 41
    log.info("Effective block height: {}".format(effective_block))
    # Report prepareblock signature
    report_information = mock_duplicate_sign(1, client_con_list_obj[0].node.nodekey, client_con_list_obj[0].node.blsprikey,
                                             effective_block)
    log.info("Report information: {}".format(report_information))
    result = client_con_list_obj[0].duplicatesign.reportDuplicateSign(1, report_information, report_address)
    assert_code(result, 0)


@pytest.mark.P1
def test_PIP_PVF_007(new_genesis_env, client_con_list_obj, reset_environment):
    """
    治理修改区块双签-证据有效期处于未生效期
    :param new_genesis_env:
    :param client_con_list_obj:
    :param reset_environment:
    :return:
    """
    # Change configuration parameters
    genesis = from_dict(data_class=Genesis, data=new_genesis_env.genesis_config)
    genesis.EconomicModel.Staking.UnStakeFreezeDuration = 3
    genesis.EconomicModel.Slashing.MaxEvidenceAge = 2
    new_file = new_genesis_env.cfg.env_tmp + "/genesis.json"
    genesis.to_file(new_file)
    new_genesis_env.deploy_all(new_file)
    # view Parameter value before treatment
    slash_blocks1 = get_governable_parameter_value(client_con_list_obj[0], 'Slashing', 'MaxEvidenceAge')
    assert slash_blocks1 == '1', "ErrMsg:Parameter value before treatment {}".format(slash_blocks1)
    # create Parametric proposal
    param_governance_verify_before_endblock(client_con_list_obj[0], 'Slashing', 'MaxEvidenceAge', '1')
    # view Parameter value before treatment again
    slash_blocks2 = get_governable_parameter_value(client_con_list_obj[0], 'Slashing', 'MaxEvidenceAge')
    assert slash_blocks2 == slash_blocks1, "ErrMsg:Parameter value after treatment {}".format(slash_blocks2)
    report_address, _ = client_con_list_obj[0].economic.account.generate_account(client_con_list_obj[0].node.web3,
                                                                                 client_con_list_obj[0].node.web3.toWei(
                                                                                     1000, 'ether'))
    # Verify changed parameters
    effective_block = client_con_list_obj[0].economic.get_front_settlement_switchpoint(client_con_list_obj[0].node, 2)
    if effective_block < 41:
        log.info("Current block: {}".format(client_con_list_obj[0].node.eth.blockNumber))
        effective_block = 41
    log.info("Effective block height: {}".format(effective_block))
    # Report prepareblock signature
    report_information = mock_duplicate_sign(1, client_con_list_obj[0].node.nodekey,
                                             client_con_list_obj[0].node.blsprikey,
                                             effective_block)
    log.info("Report information: {}".format(report_information))
    result = client_con_list_obj[0].duplicatesign.reportDuplicateSign(1, report_information, report_address)
    assert_code(result, 0)


@pytest.mark.P1
def test_PIP_PVF_008(new_genesis_env, client_con_list_obj, reset_environment):
    """
    治理修改区块双签-证据有效期处于已生效期
    :param new_genesis_env:
    :param client_con_list_obj:
    :param reset_environment:
    :return:
    """
    # Change configuration parameters
    genesis = from_dict(data_class=Genesis, data=new_genesis_env.genesis_config)
    genesis.EconomicModel.Staking.UnStakeFreezeDuration = 3
    genesis.EconomicModel.Slashing.MaxEvidenceAge = 2
    new_file = new_genesis_env.cfg.env_tmp + "/genesis.json"
    genesis.to_file(new_file)
    new_genesis_env.deploy_all(new_file)
    # view Parameter value before treatment
    slash_blocks1 = get_governable_parameter_value(client_con_list_obj[0], 'Slashing', 'MaxEvidenceAge')
    assert slash_blocks1 == '2', "ErrMsg:Parameter value before treatment {}".format(slash_blocks1)
    # create Parametric proposal
    param_governance_verify(client_con_list_obj[0], 'Slashing', 'MaxEvidenceAge', '1')
    # view Parameter value before treatment again
    slash_blocks2 = get_governable_parameter_value(client_con_list_obj[0], 'Slashing', 'MaxEvidenceAge')
    assert slash_blocks2 == '1', "ErrMsg:Parameter value after treatment {}".format(slash_blocks2)
    report_address, _ = client_con_list_obj[0].economic.account.generate_account(client_con_list_obj[0].node.web3,
                                                                                 client_con_list_obj[0].node.web3.toWei(
                                                                                     1000, 'ether'))
    # Verify changed parameters
    effective_block1 = client_con_list_obj[0].economic.get_front_settlement_switchpoint(client_con_list_obj[0].node, int(slash_blocks1))
    effective_block2 = client_con_list_obj[0].economic.get_front_settlement_switchpoint(client_con_list_obj[0].node, int(slash_blocks2))
    log.info("Effective1 block height: {}".format(effective_block1))
    log.info("Effective2 block height: {}".format(effective_block2))
    # Report1 prepareblock signature
    report_information = mock_duplicate_sign(1, client_con_list_obj[0].node.nodekey,
                                             client_con_list_obj[0].node.blsprikey,
                                             effective_block1)
    log.info("Report information: {}".format(report_information))
    result = client_con_list_obj[0].duplicatesign.reportDuplicateSign(1, report_information, report_address)
    assert_code(result, 303003)
    # Report2 prepareblock signature
    report_information = mock_duplicate_sign(1, client_con_list_obj[0].node.nodekey,
                                             client_con_list_obj[0].node.blsprikey,
                                             effective_block2)
    log.info("Report information: {}".format(report_information))
    result = client_con_list_obj[0].duplicatesign.reportDuplicateSign(1, report_information, report_address)
    assert_code(result, 0)


@pytest.mark.P1
def test_PIP_PVF_009(new_genesis_env, client_con_list_obj, reset_environment):
    """
    治理修改区块双签-证据有效期（节点质押退回锁定周期-1）
    :param new_genesis_env:
    :param client_con_list_obj:
    :param reset_environment:
    :return:
    """
    # Change configuration parameters
    genesis = from_dict(data_class=Genesis, data=new_genesis_env.genesis_config)
    genesis.EconomicModel.Staking.UnStakeFreezeDuration = 3
    genesis.EconomicModel.Slashing.MaxEvidenceAge = 1
    new_file = new_genesis_env.cfg.env_tmp + "/genesis.json"
    genesis.to_file(new_file)
    new_genesis_env.deploy_all(new_file)
    # view Parameter value before treatment
    slash_blocks1 = get_governable_parameter_value(client_con_list_obj[0], 'Slashing', 'MaxEvidenceAge')
    assert slash_blocks1 == '1', "ErrMsg:Parameter value before treatment {}".format(slash_blocks1)
    # create Parametric proposal
    param_governance_verify(client_con_list_obj[0], 'Slashing', 'MaxEvidenceAge', '2')
    # view Parameter value before treatment again
    slash_blocks2 = get_governable_parameter_value(client_con_list_obj[0], 'Slashing', 'MaxEvidenceAge')
    assert slash_blocks2 == '2', "ErrMsg:Parameter value after treatment {}".format(slash_blocks2)
    report_address, _ = client_con_list_obj[0].economic.account.generate_account(client_con_list_obj[0].node.web3,
                                                                                 client_con_list_obj[0].node.web3.toWei(
                                                                                     1000, 'ether'))
    # Verify changed parameters
    effective_block1 = client_con_list_obj[0].economic.get_front_settlement_switchpoint(client_con_list_obj[0].node,
                                                                                        int(slash_blocks1))
    effective_block2 = client_con_list_obj[0].economic.get_front_settlement_switchpoint(client_con_list_obj[0].node,
                                                                                        int(slash_blocks2))
    log.info("Effective1 block height: {}".format(effective_block1))
    log.info("Effective2 block height: {}".format(effective_block2))
    # Report1 prepareblock signature
    report_information = mock_duplicate_sign(1, client_con_list_obj[0].node.nodekey,
                                             client_con_list_obj[0].node.blsprikey,
                                             effective_block1)
    log.info("Report information: {}".format(report_information))
    result = client_con_list_obj[0].duplicatesign.reportDuplicateSign(1, report_information, report_address)
    assert_code(result, 0)
    # Report2 prepareblock signature
    report_information = mock_duplicate_sign(1, client_con_list_obj[0].node.nodekey,
                                             client_con_list_obj[0].node.blsprikey,
                                             effective_block2)
    log.info("Report information: {}".format(report_information))
    result = client_con_list_obj[0].duplicatesign.reportDuplicateSign(1, report_information, report_address)
    assert_code(result, 303000)
