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
    wait_block_number


def pledge_punishment(client_con_list_obj):
    """
    :return:
    """
    client1 = client_con_list_obj[0]
    client2 = client_con_list_obj[1]
    log.info("Current block height: {}".format(client1.node.eth.blockNumber))
    # stop node
    client1.node.stop()
    # Waiting for a settlement round
    client2.economic.wait_consensus_blocknum(client2.node, 2)
    log.info("Current block height: {}".format(client2.node.eth.blockNumber))
    # view verifier list
    verifier_list = client2.ppos.getVerifierList()
    log.info("verifier_list: {}".format(verifier_list))
    candidate_info = client2.ppos.getCandidateInfo(client1.node.node_id)
    log.info("Pledge node information： {}".format(candidate_info))
    return candidate_info


def information_before_slash_blocks(client_obj):
    node = client_obj.node
    # view Consensus Amount of pledge
    candidate_info1 = client_obj.ppos.getCandidateInfo(node.node_id)
    pledge_amount1 = candidate_info1['Ret']['Released']
    # view block_reward
    log.info("block: {}".format(node.eth.blockNumber))
    block_reward, staking_reward = client_obj.economic.get_current_year_reward(node)
    log.info("block_reward: {} staking_reward: {}".format(block_reward, staking_reward))
    # Get governable parameters
    slash_blocks1 = get_governable_parameter_value(client_obj, 'slashBlocksReward')
    return pledge_amount1, block_reward, slash_blocks1


def Verify_changed_parameters(client_con_list_obj, pledge_amount1, block_reward, slash_blocks):
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
def test_PIP_PVF_001(client_con_list_obj, reset_environment):
    """
    治理修改低0出块率扣除验证人自有质押金块数投票失败
    :param client_con_list_obj:
    :return:
    """
    # Initialize environment
    client_con_list_obj[0].economic.env.deploy_all()
    time.sleep(3)
    # get pledge amount1 and block reward
    pledge_amount1, block_reward, slash_blocks1 = information_before_slash_blocks(client_con_list_obj[0])
    # create Parametric proposal
    param_governance_verify_before_endblock(client_con_list_obj[0], 'slashing', 'slashBlocksReward', '0', False)
    # Get governable parameters again
    slash_blocks2 = get_governable_parameter_value(client_con_list_obj[0], 'slashBlocksReward')
    assert slash_blocks1 == slash_blocks2, "ErrMsg:slash blocks {}".format(slash_blocks2)
    # Verify changed parameters
    Verify_changed_parameters(client_con_list_obj, pledge_amount1, block_reward, slash_blocks1)


@pytest.mark.P1
def test_PIP_PVF_002(client_con_list_obj, reset_environment):
    """
    理修改低0出块率扣除验证人自有质押金块数成功处于未生效期
    :param client_con_list_obj:
    :return:
    """
    # get pledge amount1 and block reward
    pledge_amount1, block_reward, slash_blocks1 = information_before_slash_blocks(client_con_list_obj[0])
    # create Parametric proposal
    param_governance_verify_before_endblock(client_con_list_obj[0], 'slashing', 'slashBlocksReward', '0')
    # Get governable parameters again
    slash_blocks2 = get_governable_parameter_value(client_con_list_obj[0], 'slashBlocksReward')
    assert slash_blocks1 == slash_blocks2, "ErrMsg:slash blocks {}".format(slash_blocks2)
    # Verify changed parameters
    Verify_changed_parameters(client_con_list_obj, pledge_amount1, block_reward, slash_blocks1)


@pytest.mark.P1
def test_PIP_PVF_003(client_con_list_obj, reset_environment):
    """
    治理修改低0出块率扣除验证人自有质押金块数成功处于已生效期
    :param client_con_list_obj:
    :param reset_environment:
    :return:
    """
    # get pledge amount1 and block reward
    pledge_amount1, block_reward, slash_blocks1 = information_before_slash_blocks(client_con_list_obj[0])
    # create Parametric proposal
    param_governance_verify(client_con_list_obj[0], 'slashing', 'slashBlocksReward', '0')
    log.info("Current block height: {}".format(client_con_list_obj[0].node.eth.blockNumber))
    # Get governable parameters again
    slash_blocks2 = get_governable_parameter_value(client_con_list_obj[0], 'slashBlocksReward')
    assert slash_blocks2 == '0', "ErrMsg:Change parameters {}".format(slash_blocks2)
    # Verify changed parameters
    Verify_changed_parameters(client_con_list_obj, pledge_amount1, block_reward, slash_blocks2)


@pytest.mark.P1
def test_PIP_PVF_004(client_con_list_obj, client_new_node_obj_list, reset_environment):
    """
    治理修改低0出块率扣除验证人自有质押金块数成功扣除区块奖励块数60100-自由金额质押
    :param client_con_list_obj:
    :param reset_environment:
    :return:
    """
    # get pledge amount1 and block reward
    pledge_amount1, block_reward, slash_blocks1 = information_before_slash_blocks(client_con_list_obj[0])
    # create Parametric proposal
    param_governance_verify(client_con_list_obj[0], 'slashing', 'slashBlocksReward', '60100')
    log.info("Current block height: {}".format(client_con_list_obj[0].node.eth.blockNumber))
    # Get governable parameters again
    slash_blocks2 = get_governable_parameter_value(client_con_list_obj[0], 'slashBlocksReward')
    assert slash_blocks2 == '60100', "ErrMsg:Change parameters {}".format(slash_blocks2)
    # create account
    address, _ = client_con_list_obj[0].economic.account.generate_account(client_con_list_obj[0].node.web3,
                                                                          client_con_list_obj[
                                                                              0].economic.create_staking_limit * 2)
    # create staking
    result = client_new_node_obj_list[0].staking.create_staking(0, address, address)
    assert_code(result, 0)
    # wait settlement block
    client_new_node_obj_list[0].economic.wait_settlement_blocknum(client_new_node_obj_list[0].node)
    for i in range(4):
        result = check_node_in_list(client_con_list_obj[0].node.node_id, client_con_list_obj[0].ppos.getValidatorList)
        log.info("Current node in consensus list status：{}".format(result))
        if result:
            # Verify changed parameters
            Verify_changed_parameters(client_con_list_obj, pledge_amount1, block_reward, slash_blocks2)
            break
        else:
            # wait consensus block
            client_con_list_obj[0].economic.wait_consensus_blocknum(client_con_list_obj[0].node)


@pytest.mark.P1
def test_PIP_PVF_005(client_con_list_obj, client_new_node_obj_list, reset_environment):
    """
    治理修改低出块率扣除验证人自有质押金比例扣除区块奖励块数60100-锁仓金额质押
    :param client_con_list_obj:
    :param client_new_node_obj_list:
    :param reset_environment:
    :return:
    """
    # get pledge amount1 and block reward
    pledge_amount1, block_reward, slash_blocks1 = information_before_slash_blocks(client_con_list_obj[0])
    # create Parametric proposal
    param_governance_verify(client_con_list_obj[0], 'slashing', 'slashBlocksReward', '60100')
    log.info("Current block height: {}".format(client_con_list_obj[0].node.eth.blockNumber))
    # Get governable parameters
    slash_blocks2 = get_governable_parameter_value(client_con_list_obj[0], 'slashBlocksReward')
    assert slash_blocks2 == '60100', "ErrMsg:Change parameters {}".format(slash_blocks2)
    # create account
    address, _ = client_con_list_obj[0].economic.account.generate_account(client_con_list_obj[0].node.web3,
                                                                          client_con_list_obj[
                                                                              0].economic.create_staking_limit * 2)
    address1, _ = client_con_list_obj[0].economic.account.generate_account(client_con_list_obj[0].node.web3,
                                                                           client_con_list_obj[0].node.web3.toWei(1000,
                                                                                                                  'ether'))
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
        # log.info("Current node in consensus list status：{}".format(result))
        if result:
            # Verify changed parameters
            Verify_changed_parameters(client_con_list_obj, pledge_amount1, block_reward, slash_blocks2)
            break
        else:
            # wait consensus block
            client_con_list_obj[0].economic.wait_consensus_blocknum(client_con_list_obj[0].node)


def adjust_initial_parameters(new_genesis_env):
    # Change configuration parameters
    genesis = from_dict(data_class=Genesis, data=new_genesis_env.genesis_config)
    genesis.economicModel.staking.unStakeFreezeDuration = 3
    genesis.economicModel.slashing.maxEvidenceAge = 2
    new_file = new_genesis_env.cfg.env_tmp + "/genesis.json"
    genesis.to_file(new_file)
    new_genesis_env.deploy_all(new_file)


@pytest.mark.P1
def test_PIP_PVF_006(new_genesis_env, client_con_list_obj):
    """
    治理修改区块双签-证据有效期投票失败
    :param client_con_list_obj:
    :return:
    """
    # Change configuration parameters
    adjust_initial_parameters(new_genesis_env)
    # view Parameter value before treatment
    slash_blocks1 = get_governable_parameter_value(client_con_list_obj[0], 'maxEvidenceAge')
    # create Parametric proposal
    param_governance_verify_before_endblock(client_con_list_obj[0], 'slashing', 'maxEvidenceAge', '1', False)
    # view Parameter value before treatment again
    slash_blocks2 = get_governable_parameter_value(client_con_list_obj[0], 'maxEvidenceAge')
    assert slash_blocks2 == slash_blocks1, "ErrMsg:Parameter value after treatment {}".format(slash_blocks2)
    report_address, _ = client_con_list_obj[0].economic.account.generate_account(client_con_list_obj[0].node.web3,
                                                                                 client_con_list_obj[0].node.web3.toWei(
                                                                                     1000, 'ether'))
    # wait consensus block
    client_con_list_obj[0].economic.wait_consensus_blocknum(client_con_list_obj[0].node)
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
def test_PIP_PVF_007(new_genesis_env, client_con_list_obj):
    """
    治理修改区块双签-证据有效期处于未生效期
    :param new_genesis_env:
    :param client_con_list_obj:
    :return:
    """
    # Change configuration parameters
    adjust_initial_parameters(new_genesis_env)
    # view Parameter value before treatment
    slash_blocks1 = get_governable_parameter_value(client_con_list_obj[0], 'maxEvidenceAge')
    # create Parametric proposal
    param_governance_verify_before_endblock(client_con_list_obj[0], 'slashing', 'maxEvidenceAge', '1')
    # view Parameter value before treatment again
    slash_blocks2 = get_governable_parameter_value(client_con_list_obj[0], 'maxEvidenceAge')
    assert slash_blocks2 == slash_blocks1, "ErrMsg:Parameter value after treatment {}".format(slash_blocks2)
    report_address, _ = client_con_list_obj[0].economic.account.generate_account(client_con_list_obj[0].node.web3,
                                                                                 client_con_list_obj[0].node.web3.toWei(
                                                                                     1000, 'ether'))
    # wait consensus block
    client_con_list_obj[0].economic.wait_consensus_blocknum(client_con_list_obj[0].node)
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
def test_PIP_PVF_008(new_genesis_env, client_con_list_obj):
    """
    治理修改区块双签-证据有效期处于已生效期
    :param new_genesis_env:
    :param client_con_list_obj:
    :return:
    """
    # Change configuration parameters
    adjust_initial_parameters(new_genesis_env)
    # view Parameter value before treatment
    slash_blocks1 = get_governable_parameter_value(client_con_list_obj[0], 'maxEvidenceAge')
    # create Parametric proposal
    param_governance_verify(client_con_list_obj[0], 'slashing', 'maxEvidenceAge', '1')
    # view Parameter value before treatment again
    slash_blocks2 = get_governable_parameter_value(client_con_list_obj[0], 'maxEvidenceAge')
    assert slash_blocks2 == '1', "ErrMsg:Parameter value after treatment {}".format(slash_blocks2)
    report_address, _ = client_con_list_obj[0].economic.account.generate_account(client_con_list_obj[0].node.web3,
                                                                                 client_con_list_obj[0].node.web3.toWei(
                                                                                     1000, 'ether'))
    # wait consensus block
    client_con_list_obj[0].economic.wait_consensus_blocknum(client_con_list_obj[0].node)
    # Verify changed parameters
    effective_block2 = client_con_list_obj[0].economic.get_front_settlement_switchpoint(client_con_list_obj[0].node,
                                                                                        int(slash_blocks2))
    log.info("Effective2 block height: {}".format(effective_block2))
    # wait consensus block
    client_con_list_obj[0].economic.wait_consensus_blocknum(client_con_list_obj[0].node)
    # Report2 prepareblock signature
    report_information = mock_duplicate_sign(1, client_con_list_obj[0].node.nodekey,
                                             client_con_list_obj[0].node.blsprikey,
                                             effective_block2)
    log.info("Report information: {}".format(report_information))
    result = client_con_list_obj[0].duplicatesign.reportDuplicateSign(1, report_information, report_address)
    assert_code(result, 0)


@pytest.mark.P1
def test_PIP_PVF_009(new_genesis_env, client_con_list_obj):
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
    slash_blocks1 = get_governable_parameter_value(client_con_list_obj[0], 'maxEvidenceAge')
    # create Parametric proposal
    param_governance_verify(client_con_list_obj[0], 'slashing', 'maxEvidenceAge', '2')
    # view Parameter value before treatment again
    slash_blocks2 = get_governable_parameter_value(client_con_list_obj[0], 'maxEvidenceAge')
    assert slash_blocks2 == '2', "ErrMsg:Parameter value after treatment {}".format(slash_blocks2)
    report_address, _ = client_con_list_obj[0].economic.account.generate_account(client_con_list_obj[0].node.web3,
                                                                                 client_con_list_obj[0].node.web3.toWei(
                                                                                     1000, 'ether'))
    # wait consensus block
    client_con_list_obj[0].economic.wait_consensus_blocknum(client_con_list_obj[0].node)
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


@pytest.mark.P1
def test_PIP_PVF_010(new_genesis_env, client_con_list_obj):
    """
    治理修改区块双签-证据有效期（超出有效期）
    :param new_genesis_env:
    :param client_con_list_obj:
    :return:
    """
    # Change configuration parameters
    adjust_initial_parameters(new_genesis_env)
    # view Parameter value before treatment
    slash_blocks1 = get_governable_parameter_value(client_con_list_obj[0], 'maxEvidenceAge')
    assert slash_blocks1 == '2', "ErrMsg:Parameter value before treatment {}".format(slash_blocks1)
    # create Parametric proposal
    param_governance_verify(client_con_list_obj[0], 'slashing', 'maxEvidenceAge', '1')
    # view Parameter value before treatment again
    slash_blocks2 = get_governable_parameter_value(client_con_list_obj[0], 'maxEvidenceAge')
    assert slash_blocks2 == '1', "ErrMsg:Parameter value after treatment {}".format(slash_blocks2)
    report_address, _ = client_con_list_obj[0].economic.account.generate_account(client_con_list_obj[0].node.web3,
                                                                                 client_con_list_obj[0].node.web3.toWei(
                                                                                     1000, 'ether'))
    # wait consensus block
    client_con_list_obj[0].economic.wait_consensus_blocknum(client_con_list_obj[0].node)
    # Verify changed parameters
    effective_block1 = client_con_list_obj[0].economic.get_front_settlement_switchpoint(client_con_list_obj[0].node,
                                                                                        int(slash_blocks1))
    log.info("Effective1 block height: {}".format(effective_block1))
    # Report1 prepareblock signature
    report_information = mock_duplicate_sign(1, client_con_list_obj[0].node.nodekey,
                                             client_con_list_obj[0].node.blsprikey,
                                             effective_block1)
    log.info("Report information: {}".format(report_information))
    result = client_con_list_obj[0].duplicatesign.reportDuplicateSign(1, report_information, report_address)
    assert_code(result, 303003)


def information_before_penalty_ratio(client_obj):
    # view Pledge amount
    candidate_info1 = client_obj.ppos.getCandidateInfo(client_obj.node.node_id)
    pledge_amount1 = candidate_info1['Ret']['Released']
    # view Parameter value before treatment
    penalty_ratio1 = get_governable_parameter_value(client_obj, 'slashFractionDuplicateSign')
    return pledge_amount1, penalty_ratio1


def duplicate_sign(client_obj, report_address, report_block):
    if report_block < 41:
        report_block = 41
    # Report1 prepareblock signature
    report_information = mock_duplicate_sign(1, client_obj.node.nodekey,
                                             client_obj.node.blsprikey,
                                             report_block)
    log.info("Report information: {}".format(report_information))
    result = client_obj.duplicatesign.reportDuplicateSign(1, report_information, report_address)
    assert_code(result, 0)


def assret_penalty_amount(client_con_list_obj, pledge_amount1, penalty_ratio=None):
    # view Pledge amount after punishment
    proportion_reward, incentive_pool_reward = client_con_list_obj[1].economic.get_report_reward(pledge_amount1,
                                                                                                 penalty_ratio)
    log.info("Whistleblower benefits：{} Incentive pool income：{}".format(proportion_reward, incentive_pool_reward))
    # view Pledge amount again
    candidate_info2 = client_con_list_obj[1].ppos.getCandidateInfo(client_con_list_obj[0].node.node_id)
    pledge_amount2 = candidate_info2['Ret']['Released']
    assert pledge_amount2 == pledge_amount1 - proportion_reward - incentive_pool_reward, "ErrMsg:Pledge amount {}".format(
        pledge_amount2)


@pytest.mark.P1
def test_PIP_PVF_011(client_con_list_obj, reset_environment):
    """
    治理修改区块双签-最高处罚比例投票失败
    :param client_con_list_obj:
    :param reset_environment:
    :return:
    """
    # view Pledge amount and Parameter value before treatment
    pledge_amount1, penalty_ratio1 = information_before_penalty_ratio(client_con_list_obj[0])
    # create Parametric proposal
    param_governance_verify_before_endblock(client_con_list_obj[0], 'slashing', 'slashFractionDuplicateSign', '1000',
                                            False)
    # view Parameter value after treatment again
    penalty_ratio2 = get_governable_parameter_value(client_con_list_obj[0], 'slashFractionDuplicateSign')
    assert penalty_ratio1 == penalty_ratio2, "ErrMsg:Parameter value after treatment {}".format(penalty_ratio2)
    # wait consensus block
    client_con_list_obj[0].economic.wait_consensus_blocknum(client_con_list_obj[0].node)
    # create account
    report_address, _ = client_con_list_obj[0].economic.account.generate_account(client_con_list_obj[0].node.web3,
                                                                                 client_con_list_obj[0].node.web3.toWei(
                                                                                     1000, 'ether'))
    # Verify changed parameters
    current_block = client_con_list_obj[0].node.eth.blockNumber
    log.info("Current block height: {}".format(current_block))
    # Verify changed parameters report
    duplicate_sign(client_con_list_obj[0], report_address, current_block)
    # assret penalty amount
    assret_penalty_amount(client_con_list_obj, pledge_amount1)


@pytest.mark.P1
def test_PIP_PVF_012(client_con_list_obj, reset_environment):
    """
    治理修改区块双签-最高处罚比例处于未生效期
    :param client_con_list_obj:
    :param reset_environment:
    :return:
    """
    # view Pledge amount and Parameter value before treatment
    pledge_amount1, penalty_ratio1 = information_before_penalty_ratio(client_con_list_obj[0])
    # create Parametric proposal
    param_governance_verify_before_endblock(client_con_list_obj[0], 'slashing', 'slashFractionDuplicateSign', '1000')
    # view Parameter value before treatment again
    penalty_ratio2 = get_governable_parameter_value(client_con_list_obj[0], 'slashFractionDuplicateSign')
    assert penalty_ratio1 == penalty_ratio2, "ErrMsg:Parameter value after treatment {}".format(penalty_ratio2)
    # wait consensus block
    client_con_list_obj[0].economic.wait_consensus_blocknum(client_con_list_obj[0].node)
    # create account
    report_address, _ = client_con_list_obj[0].economic.account.generate_account(client_con_list_obj[0].node.web3,
                                                                                 client_con_list_obj[0].node.web3.toWei(
                                                                                     1000, 'ether'))
    # Verify changed parameters
    current_block = client_con_list_obj[0].node.eth.blockNumber
    log.info("Current block height: {}".format(current_block))
    # Verify changed parameters report
    duplicate_sign(client_con_list_obj[0], report_address, current_block)
    # assret penalty amount
    assret_penalty_amount(client_con_list_obj, pledge_amount1)


@pytest.mark.P1
def test_PIP_PVF_013(client_con_list_obj, reset_environment):
    """
    治理修改区块双签-最高处罚比例处于已生效期
    :param client_con_list_obj:
    :param reset_environment:
    :return:
    """
    # view Pledge amount and Parameter value before treatment
    pledge_amount1, penalty_ratio1 = information_before_penalty_ratio(client_con_list_obj[0])
    # create Parametric proposal
    param_governance_verify(client_con_list_obj[0], 'slashing', 'slashFractionDuplicateSign', '1000')
    # view Parameter value before treatment again
    penalty_ratio2 = get_governable_parameter_value(client_con_list_obj[0], 'slashFractionDuplicateSign')
    assert penalty_ratio2 == '1000', "ErrMsg:Parameter value after treatment {}".format(penalty_ratio2)
    # create account
    report_address, _ = client_con_list_obj[0].economic.account.generate_account(client_con_list_obj[0].node.web3,
                                                                                 client_con_list_obj[0].node.web3.toWei(
                                                                                     1000, 'ether'))
    # Verify changed parameters
    current_block = client_con_list_obj[0].node.eth.blockNumber
    log.info("Current block height: {}".format(current_block))
    # Verify changed parameters report
    duplicate_sign(client_con_list_obj[0], report_address, current_block)
    # assret penalty amount
    assret_penalty_amount(client_con_list_obj, pledge_amount1, 1000)


@pytest.mark.P1
def test_PIP_PVF_014(client_con_list_obj, reset_environment):
    """
    治理修改区块双签-最高处罚比例为10000‱
    :param client_con_list_obj:
    :param reset_environment:
    :return:
    """
    # view Pledge amount and Parameter value before treatment
    pledge_amount1, penalty_ratio1 = information_before_penalty_ratio(client_con_list_obj[0])
    # create Parametric proposal
    param_governance_verify(client_con_list_obj[0], 'slashing', 'slashFractionDuplicateSign', '10000')
    # view Parameter value before treatment again
    penalty_ratio2 = get_governable_parameter_value(client_con_list_obj[0], 'slashFractionDuplicateSign')
    assert penalty_ratio2 == '10000', "ErrMsg:Parameter value after treatment {}".format(penalty_ratio2)
    # create account
    report_address, _ = client_con_list_obj[0].economic.account.generate_account(client_con_list_obj[0].node.web3,
                                                                                 client_con_list_obj[0].node.web3.toWei(
                                                                                     1000, 'ether'))
    # Verify changed parameters
    current_block = client_con_list_obj[0].node.eth.blockNumber
    log.info("Current block height: {}".format(current_block))
    # Verify changed parameters report
    duplicate_sign(client_con_list_obj[0], report_address, current_block)
    # assret penalty amount
    assret_penalty_amount(client_con_list_obj, pledge_amount1, 10000)


@pytest.mark.P1
def test_PIP_PVF_015(client_con_list_obj, reset_environment):
    """
    治理修改区块双签-最高处罚比例为1‱
    :param client_con_list_obj:
    :param reset_environment:
    :return:
    """
    # view Pledge amount and Parameter value before treatment
    pledge_amount1, penalty_ratio1 = information_before_penalty_ratio(client_con_list_obj[0])
    # create Parametric proposal
    param_governance_verify(client_con_list_obj[0], 'slashing', 'slashFractionDuplicateSign', '1')
    # view Parameter value before treatment again
    penalty_ratio2 = get_governable_parameter_value(client_con_list_obj[0], 'slashFractionDuplicateSign')
    assert penalty_ratio2 == '1', "ErrMsg:Parameter value after treatment {}".format(penalty_ratio2)
    # create account
    report_address, _ = client_con_list_obj[0].economic.account.generate_account(client_con_list_obj[0].node.web3,
                                                                                 client_con_list_obj[0].node.web3.toWei(
                                                                                     1000, 'ether'))
    # Verify changed parameters
    current_block = client_con_list_obj[0].node.eth.blockNumber
    log.info("Current block height: {}".format(current_block))
    # Verify changed parameters report
    duplicate_sign(client_con_list_obj[0], report_address, current_block)
    # assret penalty amount
    assret_penalty_amount(client_con_list_obj, pledge_amount1, 1)


def information_before_report_reward(client_obj):
    # view Pledge amount
    candidate_info1 = client_obj.ppos.getCandidateInfo(client_obj.node.node_id)
    pledge_amount1 = candidate_info1['Ret']['Released']
    # view Parameter value before treatment
    report_reward1 = get_governable_parameter_value(client_obj, 'duplicateSignReportReward')
    return pledge_amount1, report_reward1


def get_account_amount(client_obj):
    # create report account
    report_address, _ = client_obj.economic.account.generate_account(client_obj.node.web3, client_obj.node.web3.toWei(
        1000, 'ether'))
    # view report amount
    report_amount1 = client_obj.node.eth.getBalance(report_address)
    # view Incentive pool account
    incentive_pool_account1 = client_obj.node.eth.getBalance(EconomicConfig.INCENTIVEPOOL_ADDRESS)
    return report_address, report_amount1, incentive_pool_account1


def asster_income_account_amount(client_obj, report_amount1, incentive_pool_account1, report_address, proportion_reward,
                                 incentive_pool_reward):
    # view report amount
    report_amount2 = client_obj.node.eth.getBalance(report_address)
    # view Incentive pool account
    incentive_pool_account2 = client_obj.node.eth.getBalance(EconomicConfig.INCENTIVEPOOL_ADDRESS)
    # asster amount reward
    log.info("report_amount1 {} ,proportion_reward {} , report_amount2 {}".format(report_amount1, proportion_reward,
                                                                                  report_amount2))
    assert report_amount1 + proportion_reward - report_amount2 < client_obj.node.web3.toWei(1,
                                                                                            'ether'), "ErrMsg:report amount {}".format(
        report_amount2)
    log.info("incentive_pool_account2 {} ,incentive_pool_account1 {} , incentive_pool_reward {}".format(
        incentive_pool_account2, incentive_pool_account1, incentive_pool_reward))
    assert incentive_pool_account2 == incentive_pool_account1 + incentive_pool_reward + (
        report_amount1 + proportion_reward - report_amount2), "ErrMsg:Incentive pool account {}".format(
        incentive_pool_account2)


@pytest.mark.P1
def test_PIP_PVF_016(client_con_list_obj, reset_environment):
    """
    治理修改区块双签-举报奖励比例投票失败
    :param client_con_list_obj:
    :param reset_environment:
    :return:
    """
    # get pledge_amount1 report_amount1 incentive_pool_account1 report_reward1
    pledge_amount1, report_reward1 = information_before_report_reward(client_con_list_obj[0])
    # create Parametric proposal
    param_governance_verify_before_endblock(client_con_list_obj[0], 'slashing', 'duplicateSignReportReward', '60',
                                            False)
    # view Parameter value after treatment
    report_reward2 = get_governable_parameter_value(client_con_list_obj[0], 'duplicateSignReportReward')
    assert report_reward1 == report_reward2, "ErrMsg:Parameter value after treatment {}".format(report_reward2)
    # wait consensus block
    client_con_list_obj[0].economic.wait_consensus_blocknum(client_con_list_obj[0].node)
    # get account amount
    report_address, report_amount1, incentive_pool_account1 = get_account_amount(client_con_list_obj[0])
    # Verify changed parameters
    current_block = client_con_list_obj[0].node.eth.blockNumber
    log.info("Current block height: {}".format(current_block))
    # Verify changed parameters report
    duplicate_sign(client_con_list_obj[0], report_address, current_block)
    # view Pledge amount after punishment
    proportion_reward, incentive_pool_reward = client_con_list_obj[1].economic.get_report_reward(pledge_amount1)
    log.info("Whistleblower benefits：{} Incentive pool income：{}".format(proportion_reward, incentive_pool_reward))
    # asster account amount
    asster_income_account_amount(client_con_list_obj[1], report_amount1, incentive_pool_account1,
                                 report_address, proportion_reward, incentive_pool_reward)


@pytest.mark.P1
def test_PIP_PVF_017(client_con_list_obj, reset_environment):
    """
    治理修改区块双签-举报奖励比例处于未生效期
    :param client_con_list_obj:
    :param reset_environment:
    :return:
    """
    # get pledge_amount1 report_amount1 incentive_pool_account1 report_reward1
    pledge_amount1, report_reward1 = information_before_report_reward(client_con_list_obj[0])
    # create Parametric proposal
    param_governance_verify_before_endblock(client_con_list_obj[0], 'slashing', 'duplicateSignReportReward', '60')
    # view Parameter value after treatment
    report_reward2 = get_governable_parameter_value(client_con_list_obj[0], 'duplicateSignReportReward')
    assert report_reward1 == report_reward2, "ErrMsg:Parameter value after treatment {}".format(report_reward2)
    # wait consensus block
    client_con_list_obj[0].economic.wait_consensus_blocknum(client_con_list_obj[0].node)
    # get account amount
    report_address, report_amount1, incentive_pool_account1 = get_account_amount(client_con_list_obj[0])
    # Verify changed parameters
    current_block = client_con_list_obj[0].node.eth.blockNumber
    log.info("Current block height: {}".format(current_block))
    # Verify changed parameters report
    duplicate_sign(client_con_list_obj[0], report_address, current_block)
    # view Pledge amount after punishment
    proportion_reward, incentive_pool_reward = client_con_list_obj[1].economic.get_report_reward(pledge_amount1)
    log.info("Whistleblower benefits：{} Incentive pool income：{}".format(proportion_reward, incentive_pool_reward))
    # asster account amount
    asster_income_account_amount(client_con_list_obj[1], report_amount1, incentive_pool_account1,
                                 report_address, proportion_reward, incentive_pool_reward)


@pytest.mark.P1
def test_PIP_PVF_018(client_con_list_obj, reset_environment):
    """
    治理修改区块双签-举报奖励比例处于已生效期
    :param client_con_list_obj:
    :param reset_environment:
    :return:
    """
    # get pledge_amount1 report_amount1 incentive_pool_account1 report_reward1
    pledge_amount1, report_reward1 = information_before_report_reward(client_con_list_obj[0])
    # create Parametric proposal
    param_governance_verify(client_con_list_obj[0], 'slashing', 'duplicateSignReportReward', '60')
    # view Parameter value after treatment
    report_reward2 = get_governable_parameter_value(client_con_list_obj[0], 'duplicateSignReportReward')
    assert report_reward2 == '60', "ErrMsg:Parameter value after treatment {}".format(report_reward2)
    # get account amount
    report_address, report_amount1, incentive_pool_account1 = get_account_amount(client_con_list_obj[0])
    # Verify changed parameters
    current_block = client_con_list_obj[0].node.eth.blockNumber
    log.info("Current block height: {}".format(current_block))
    # Verify changed parameters report
    duplicate_sign(client_con_list_obj[0], report_address, current_block)
    # view Pledge amount after punishment
    proportion_reward, incentive_pool_reward = client_con_list_obj[1].economic.get_report_reward(pledge_amount1, None,
                                                                                                 60)
    log.info("Whistleblower benefits：{} Incentive pool income：{}".format(proportion_reward, incentive_pool_reward))
    # asster account amount
    asster_income_account_amount(client_con_list_obj[1], report_amount1, incentive_pool_account1,
                                 report_address, proportion_reward, incentive_pool_reward)


@pytest.mark.P1
def test_PIP_PVF_019(client_con_list_obj, reset_environment):
    """
    治理修改区块双签-举报奖励比例为80%
    :param client_con_list_obj:
    :param reset_environment:
    :return:
    """
    # get pledge_amount1 report_reward1
    pledge_amount1, report_reward1 = information_before_report_reward(client_con_list_obj[0])
    # create Parametric proposal
    param_governance_verify(client_con_list_obj[0], 'slashing', 'duplicateSignReportReward', '80')
    # view Parameter value after treatment
    report_reward2 = get_governable_parameter_value(client_con_list_obj[0], 'duplicateSignReportReward')
    assert report_reward2 == '80', "ErrMsg:Parameter value after treatment {}".format(report_reward2)
    # get account amount
    report_address, report_amount1, incentive_pool_account1 = get_account_amount(client_con_list_obj[0])
    # Verify changed parameters
    current_block = client_con_list_obj[0].node.eth.blockNumber
    log.info("Current block height: {}".format(current_block))
    # Verify changed parameters report
    duplicate_sign(client_con_list_obj[0], report_address, current_block)
    # view Pledge amount after punishment
    proportion_reward, incentive_pool_reward = client_con_list_obj[1].economic.get_report_reward(pledge_amount1, None,
                                                                                                 80)
    log.info("Whistleblower benefits：{} Incentive pool income：{}".format(proportion_reward, incentive_pool_reward))
    # asster account amount
    asster_income_account_amount(client_con_list_obj[1], report_amount1, incentive_pool_account1,
                                 report_address, proportion_reward, incentive_pool_reward)


@pytest.mark.P1
def test_PIP_PVF_020(client_con_list_obj, reset_environment):
    """
    治理修改区块双签-举报奖励比例为1%
    :param client_con_list_obj:
    :param reset_environment:
    :return:
    """
    # get pledge_amount1 report_reward1
    pledge_amount1, report_reward1 = information_before_report_reward(client_con_list_obj[0])
    # create Parametric proposal
    param_governance_verify(client_con_list_obj[0], 'slashing', 'duplicateSignReportReward', '1')
    # view Parameter value after treatment
    report_reward2 = get_governable_parameter_value(client_con_list_obj[0], 'duplicateSignReportReward')
    assert report_reward2 == '1', "ErrMsg:Parameter value after treatment {}".format(report_reward2)
    # get account amount
    report_address, report_amount1, incentive_pool_account1 = get_account_amount(client_con_list_obj[0])
    # Verify changed parameters
    current_block = client_con_list_obj[0].node.eth.blockNumber
    log.info("Current block height: {}".format(current_block))
    # Verify changed parameters report
    duplicate_sign(client_con_list_obj[0], report_address, current_block)
    # view Pledge amount after punishment
    proportion_reward, incentive_pool_reward = client_con_list_obj[1].economic.get_report_reward(pledge_amount1, None,
                                                                                                 1)
    log.info("Whistleblower benefits：{} Incentive pool income：{}".format(proportion_reward, incentive_pool_reward))
    # asster account amount
    asster_income_account_amount(client_con_list_obj[1], report_amount1, incentive_pool_account1,
                                 report_address, proportion_reward, incentive_pool_reward)


#
# def transaction(client_obj, nonce, from_address, to_address, value):
#     account = client_obj.economic.account.accounts[from_address]
#     tmp_to_address = Web3.toChecksumAddress(to_address)
#     tmp_from_address = Web3.toChecksumAddress(from_address)
#
#     transaction_dict = {
#         "to": tmp_to_address,
#         "gasPrice": client_obj.node.eth.gasPrice,
#         "gas": 21000,
#         "nonce": nonce,
#         "data": "",
#         "chainId": client_obj.node.chain_id,
#         "value": value,
#         'from': tmp_from_address,
#     }
#     signedTransactionDict = client_obj.node.eth.account.signTransaction(
#         transaction_dict, account['prikey']
#     )
#
#     data = signedTransactionDict.rawTransaction
#     result = HexBytes(client_obj.node.eth.sendRawTransaction(data)).hex()
#    nonce = client_con_list_obj[0].node.eth.getTransactionCount(
#         client_con_list_obj[0].economic.env.account.account_with_money['address'])
#     from_address = client_con_list_obj[0].economic.env.account.account_with_money['address']
#     to_address, _ = client_con_list_obj[0].economic.account.generate_account(client_con_list_obj[0].node.web3, 0)
#     for i in range(0, 200):
#         # Transfer transaction
#         transaction(client_con_list_obj[0], nonce, from_address, to_address, 10)
#         nonce = nonce + 1
#     time.sleep(15)
#     end_block = client_con_list_obj[0].node.block_number
#     max_tx = {"block_num": 0, "tx_num": 0}
#     for i in range(1, end_block + 1):
#         tx_num = client_con_list_obj[0].node.eth.getBlockTransactionCount(i)
#         if tx_num > max_tx["tx_num"]:
#             max_tx = {"block_num": i, "tx_num": tx_num}
#     block_info = client_con_list_obj[0].node.eth.getBlock(max_tx["block_num"])
#     print(max_tx['tx_num'])
#     print(block_info['gasLimit'])
#     print(block_info['gasUsed'])
#     print("block_info", block_info)
#
# #
@pytest.mark.P1
def test_PIP_MG_001(client_con_list_obj, reset_environment):
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
    # client_con_list_obj[0].economic.env.deploy_all()

    # view Parameter value before treatment
    max_gas_limit1 = get_governable_parameter_value(client_con_list_obj[0], 'maxBlockGasLimit')
    # create Parametric proposal
    block = param_governance_verify_before_endblock(client_con_list_obj[0], 'block', 'maxBlockGasLimit', '4712389',
                                                    False)
    # view Parameter value after treatment
    max_gas_limit2 = get_governable_parameter_value(client_con_list_obj[0], 'maxBlockGasLimit')
    # wait block
    wait_block_number(client_con_list_obj[0].node, block)
    assert max_gas_limit2 == max_gas_limit1, "ErrMsg:Parameter value after treatment {}".format(max_gas_limit2)


@pytest.mark.P1
def test_PIP_MG_002(client_con_list_obj, reset_environment):
    """
    治理修改默认每个区块的最大Gas 处于未生效期
    :param client_con_list_obj:
    :param reset_environment:
    :return:
    """
    # view Parameter value before treatment
    max_gas_limit1 = get_governable_parameter_value(client_con_list_obj[0], 'maxBlockGasLimit')
    # create Parametric proposal
    param_governance_verify_before_endblock(client_con_list_obj[0], 'block', 'maxBlockGasLimit', '4712389')
    # view Parameter value after treatment
    max_gas_limit2 = get_governable_parameter_value(client_con_list_obj[0], 'maxBlockGasLimit')

    assert max_gas_limit2 == max_gas_limit1, "ErrMsg:Parameter value after treatment {}".format(max_gas_limit2)


@pytest.mark.P1
def test_PIP_MG_003(client_con_list_obj, reset_environment):
    """
    治理修改默认每个区块的最大Gas 处于已生效期
    :param client_con_list_obj:
    :param reset_environment:
    :return:
    """
    # view Parameter value before treatment
    max_gas_limit1 = get_governable_parameter_value(client_con_list_obj[0], 'maxBlockGasLimit')
    # create Parametric proposal
    param_governance_verify(client_con_list_obj[0], 'block', 'maxBlockGasLimit', '4712389')
    # view Parameter value after treatment
    max_gas_limit2 = get_governable_parameter_value(client_con_list_obj[0], 'maxBlockGasLimit')
    assert max_gas_limit2 == '4712389', "ErrMsg:Parameter value after treatment {}".format(max_gas_limit2)
