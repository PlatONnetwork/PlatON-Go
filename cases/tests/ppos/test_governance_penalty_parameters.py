import pytest
from dacite import from_dict
from common.key import mock_duplicate_sign
from common.log import log
from decimal import Decimal
from tests.conftest import (param_governance_verify,
                            param_governance_verify_before_endblock,
                            get_client_consensus,
                            staking_cfg
                            )
from tests.lib import (EconomicConfig,
                       Genesis, check_node_in_list,
                       assert_code, get_governable_parameter_value,
                       wait_block_number, von_amount
                       )


def pledge_punishment(clients):
    """
    :return:
    """
    first_index = 0
    second_index = 1
    first_client = clients[first_index]
    second_client = clients[second_index]
    log.info("Current block height: {}".format(first_client.node.eth.blockNumber))
    # stop node
    first_client.node.stop()
    # Waiting for a settlement round
    second_client.economic.wait_consensus_blocknum(second_client.node, 4)
    log.info("Current block height: {}".format(second_client.node.eth.blockNumber))
    # view verifier list
    verifier_list = second_client.ppos.getVerifierList()
    log.info("verifier_list: {}".format(verifier_list))
    candidate_info = second_client.ppos.getCandidateInfo(first_client.node.node_id)
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


def verify_changed_parameters(clients, first_pledge_amount, block_reward, slash_blocks):
    # Verify changed parameters
    candidate_info = pledge_punishment(clients)
    second_pledge_amount = candidate_info['Ret']['Released']
    punishment_amonut = int(Decimal(str(block_reward)) * Decimal(str(slash_blocks)))
    if punishment_amonut < first_pledge_amount:
        assert (second_pledge_amount == first_pledge_amount - punishment_amonut) or (second_pledge_amount == first_pledge_amount - punishment_amonut * 2), "ErrMsg:Consensus Amount of pledge {}".format(
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
    :param clients_consensus:
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
    治理修改低0出块率扣除验证人自有质押金块数成功扣除区块奖励块数49999-自由金额质押
    :param client_consensus:
    :param clients_noconsensus:
    :param reset_environment:
    :return:
    """
    consensus_client = client_consensus
    log.info("Current connection consensus node".format(consensus_client.node.node_mark))
    first_index = 0
    first_client = clients_noconsensus[first_index]
    log.info("Current connection non-consensus node：{}".format(first_client.node.node_mark))
    economic = consensus_client.economic
    node = consensus_client.node
    change_parameter_value = '49999'
    # get pledge amount1 and block reward
    consensus_pledge_amount, block_reward, first_slash_blocks = information_before_slash_blocks(consensus_client)
    # create Parametric proposal
    param_governance_verify(consensus_client, 'slashing', 'slashBlocksReward', change_parameter_value)
    log.info("Current block height: {}".format(node.eth.blockNumber))
    # Get governable parameters again
    second_slash_blocks = get_governable_parameter_value(consensus_client, 'slashBlocksReward')
    assert second_slash_blocks == change_parameter_value, "ErrMsg:Change parameters {}".format(second_slash_blocks)
    # create account
    address, _ = economic.account.generate_account(node.web3, von_amount(economic.create_staking_limit, 2))
    # create staking
    result = first_client.staking.create_staking(0, address, address)
    assert_code(result, 0)
    # wait settlement block
    economic.wait_settlement_blocknum(node)
    candidate_info = consensus_client.ppos.getCandidateInfo(first_client.node.node_id)
    first_pledge_amount = candidate_info['Ret']['Released']
    log.info("Current pledge node amount：{}".format(first_pledge_amount))
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
def test_PIP_PVF_005(client_consensus, clients_noconsensus, reset_environment):
    """
    治理修改低出块率扣除验证人自有质押金比例扣除区块奖励块数49999-锁仓金额质押
    :param client_consensus:
    :param clients_noconsensus:
    :param reset_environment:
    :return:
    """
    consensus_client = client_consensus
    log.info("Current connection consensus node".format(consensus_client.node.node_mark))
    first_index = 0
    first_client = clients_noconsensus[first_index]
    log.info("Current connection non-consensus node：{}".format(first_client.node.node_mark))
    economic = consensus_client.economic
    node = consensus_client.node
    change_parameter_value = '4999'
    # get pledge amount1 and block reward
    consensus_pledge_amount, block_reward, first_slash_blocks = information_before_slash_blocks(consensus_client)
    # create Parametric proposal
    param_governance_verify(consensus_client, 'slashing', 'slashBlocksReward', change_parameter_value)
    log.info("Current block height: {}".format(node.eth.blockNumber))
    # Get governable parameters
    second_slash_blocks = get_governable_parameter_value(consensus_client, 'slashBlocksReward')
    assert second_slash_blocks == change_parameter_value, "ErrMsg:Change parameters {}".format(second_slash_blocks)
    # create account
    address, _ = economic.account.generate_account(node.web3, von_amount(economic.create_staking_limit, 2))
    address1, _ = economic.account.generate_account(node.web3, node.web3.toWei(1000, 'ether'))
    # Create restricting plan
    plan = [{'Epoch': 1, 'Amount': economic.create_staking_limit}]
    result = consensus_client.restricting.createRestrictingPlan(address1, plan, address)
    assert_code(result, 0)
    # create staking
    result = first_client.staking.create_staking(1, address1, address1)
    assert_code(result, 0)
    # wait settlement block
    economic.wait_settlement_blocknum(node)
    candidate_info = consensus_client.ppos.getCandidateInfo(first_client.node.node_id)
    first_pledge_amount = candidate_info['Ret']['RestrictingPlan']
    log.info("Current pledge node amount：{}".format(first_pledge_amount))
    for i in range(4):
        result = check_node_in_list(node.node_id, consensus_client.ppos.getValidatorList)
        # log.info("Current node in consensus list status：{}".format(result))
        if result:
            # Verify changed parameters
            verify_changed_parameters(clients_noconsensus, first_pledge_amount, block_reward, second_slash_blocks)
            break
        else:
            # wait consensus block
            economic.wait_consensus_blocknum(node)


def adjust_initial_parameters(new_genesis_env):
    # Change configuration parameters
    genesis = from_dict(data_class=Genesis, data=new_genesis_env.genesis_config)
    genesis.economicModel.staking.unStakeFreezeDuration = 3
    genesis.economicModel.slashing.maxEvidenceAge = 2
    new_file = new_genesis_env.cfg.env_tmp + "/genesis_0.13.0.json"
    genesis.to_file(new_file)
    new_genesis_env.deploy_all(new_file)


@pytest.mark.P1
@pytest.mark.parametrize('mark', [False, True])
def test_PIP_PVF_006_007(new_genesis_env, mark):
    """
    治理修改区块双签-证据有效期投票失败
    :param new_genesis_env:
    :return:
    """
    # Change configuration parameters
    adjust_initial_parameters(new_genesis_env)
    first_client = get_client_consensus(new_genesis_env, staking_cfg)
    log.info("Current connection non-consensus node：{}".format(first_client.node.node_mark))
    economic = first_client.economic
    node = first_client.node
    # view Parameter value before treatment
    first_slash_blocks = get_governable_parameter_value(first_client, 'maxEvidenceAge')
    # create Parametric proposal
    param_governance_verify_before_endblock(first_client, 'slashing', 'maxEvidenceAge', '1', mark)
    # view Parameter value before treatment again
    second_slash_blocks = get_governable_parameter_value(first_client, 'maxEvidenceAge')
    assert second_slash_blocks == first_slash_blocks, "ErrMsg:Parameter value after treatment {}".format(
        second_slash_blocks)
    report_address, _ = first_client.economic.account.generate_account(node.web3, node.web3.toWei(1000, 'ether'))
    # wait consensus block
    economic.wait_consensus_blocknum(node)
    # Verify changed parameters
    effective_block = economic.get_front_settlement_switchpoint(node, 2)
    if effective_block < economic.consensus_size:
        log.info("Current block: {}".format(node.eth.blockNumber))
        effective_block = economic.consensus_size + 1
    log.info("Effective block height: {}".format(effective_block))
    # Report prepareblock signature
    report_information = mock_duplicate_sign(1, node.nodekey, node.blsprikey, effective_block)
    log.info("Report information: {}".format(report_information))
    result = first_client.duplicatesign.reportDuplicateSign(1, report_information, report_address)
    assert_code(result, 0)


@pytest.mark.P1
def test_PIP_PVF_008(new_genesis_env):
    """
    治理修改区块双签-证据有效期处于已生效期
    :param new_genesis_env:
    :return:
    """
    # Change configuration parameters
    adjust_initial_parameters(new_genesis_env)
    first_client = get_client_consensus(new_genesis_env, staking_cfg)
    log.info("Current connection non-consensus node：{}".format(first_client.node.node_mark))
    economic = first_client.economic
    node = first_client.node
    change_parameter_value = '1'
    # view Parameter value before treatment
    first_slash_blocks = get_governable_parameter_value(first_client, 'maxEvidenceAge')
    # create Parametric proposal
    param_governance_verify(first_client, 'slashing', 'maxEvidenceAge', change_parameter_value)
    # view Parameter value before treatment again
    second_slash_blocks = get_governable_parameter_value(first_client, 'maxEvidenceAge')
    assert second_slash_blocks != first_slash_blocks, "ErrMsg:Parameter value after treatment {}".format(
        second_slash_blocks)
    assert second_slash_blocks == change_parameter_value, "ErrMsg:Parameter value after treatment {}".format(
        second_slash_blocks)
    report_address, _ = economic.account.generate_account(node.web3, node.web3.toWei(1000, 'ether'))
    # wait consensus block
    economic.wait_consensus_blocknum(node)
    # Verify changed parameters
    effective_block = economic.get_front_settlement_switchpoint(node, int(second_slash_blocks))
    log.info("effective_block block height: {}".format(effective_block))
    # wait consensus block
    economic.wait_consensus_blocknum(node)
    # first Report prepareblock signature
    report_information = mock_duplicate_sign(1, node.nodekey, node.blsprikey, economic.consensus_size + 1)
    log.info("Report information: {}".format(report_information))
    result = first_client.duplicatesign.reportDuplicateSign(1, report_information, report_address)
    assert_code(result, 303003)
    # second Report prepareblock signature
    report_information = mock_duplicate_sign(1, node.nodekey, node.blsprikey, effective_block)
    log.info("Report information: {}".format(report_information))
    result = first_client.duplicatesign.reportDuplicateSign(1, report_information, report_address)
    assert_code(result, 0)


@pytest.mark.P1
def test_PIP_PVF_009(new_genesis_env):
    """
    治理修改区块双签-证据有效期（节点质押退回锁定周期-1）
    :param new_genesis_env:
    :return:
    """
    # Change configuration parameters
    genesis = from_dict(data_class=Genesis, data=new_genesis_env.genesis_config)
    genesis.economicModel.staking.unStakeFreezeDuration = 3
    genesis.economicModel.slashing.maxEvidenceAge = 1
    new_file = new_genesis_env.cfg.env_tmp + "/genesis_0.13.0.json"
    genesis.to_file(new_file)
    new_genesis_env.deploy_all(new_file)

    first_client = get_client_consensus(new_genesis_env, staking_cfg)
    log.info("Current connection non-consensus node：{}".format(first_client.node.node_mark))
    economic = first_client.economic
    node = first_client.node
    change_parameter_value = '2'
    # view Parameter value before treatment
    first_slash_blocks = get_governable_parameter_value(first_client, 'maxEvidenceAge')
    # create Parametric proposal
    param_governance_verify(first_client, 'slashing', 'maxEvidenceAge', change_parameter_value)
    # view Parameter value before treatment again
    second_slash_blocks = get_governable_parameter_value(first_client, 'maxEvidenceAge')
    assert second_slash_blocks != first_slash_blocks, "ErrMsg:Parameter value after treatment {}".format(
        second_slash_blocks)
    assert second_slash_blocks == change_parameter_value, "ErrMsg:Parameter value after treatment {}".format(
        second_slash_blocks)
    report_address, _ = economic.account.generate_account(node.web3, node.web3.toWei(1000, 'ether'))
    # wait consensus block
    economic.wait_consensus_blocknum(node)
    # Verify changed parameters
    first_effective_block = economic.get_front_settlement_switchpoint(node, int(first_slash_blocks))
    log.info("first effective block height: {}".format(first_effective_block))
    second_effective_block = economic.get_front_settlement_switchpoint(node, int(second_slash_blocks))
    log.info("second effective block height: {}".format(second_effective_block))
    # first Report prepareblock signature
    report_information = mock_duplicate_sign(1, node.nodekey, node.blsprikey, first_effective_block)
    log.info("Report information: {}".format(report_information))
    result = first_client.duplicatesign.reportDuplicateSign(1, report_information, report_address)
    assert_code(result, 0)
    # second Report prepareblock signature
    report_information = mock_duplicate_sign(1, node.nodekey, node.blsprikey, second_effective_block)
    log.info("Report information: {}".format(report_information))
    result = first_client.duplicatesign.reportDuplicateSign(1, report_information, report_address)
    assert_code(result, 303000)


@pytest.mark.P1
def test_PIP_PVF_010(new_genesis_env, clients_consensus):
    """
    治理修改区块双签-证据有效期（超出有效期）
    :param new_genesis_env:
    :return:
    """
    # Change configuration parameters
    adjust_initial_parameters(new_genesis_env)

    first_client = get_client_consensus(new_genesis_env, staking_cfg)
    log.info("Current connection non-consensus node：{}".format(first_client.node.node_mark))
    economic = first_client.economic
    node = first_client.node
    change_parameter_value = '1'
    # view Parameter value before treatment
    first_slash_blocks = get_governable_parameter_value(first_client, 'maxEvidenceAge')
    assert first_slash_blocks == '2', "ErrMsg:Parameter value before treatment {}".format(first_slash_blocks)
    # create Parametric proposal
    param_governance_verify(first_client, 'slashing', 'maxEvidenceAge', change_parameter_value)
    # view Parameter value before treatment again
    second_slash_blocks = get_governable_parameter_value(first_client, 'maxEvidenceAge')
    assert second_slash_blocks == change_parameter_value, "ErrMsg:Parameter value after treatment {}".format(
        second_slash_blocks)
    report_address, _ = economic.account.generate_account(node.web3, node.web3.toWei(1000, 'ether'))
    # wait consensus block
    economic.wait_consensus_blocknum(node)
    # Verify changed parameters
    effective_block = economic.get_front_settlement_switchpoint(node, int(first_slash_blocks))
    log.info("Effective1 block height: {}".format(effective_block))
    # Report1 prepareblock signature
    report_information = mock_duplicate_sign(1, node.nodekey, node.blsprikey, effective_block)
    log.info("Report information: {}".format(report_information))
    result = first_client.duplicatesign.reportDuplicateSign(1, report_information, report_address)
    assert_code(result, 303003)


def information_before_penalty_ratio(client):
    # view Pledge amount
    first_candidate_info = client.ppos.getCandidateInfo(client.node.node_id)
    first_pledge_amount = first_candidate_info['Ret']['Released']
    # view Parameter value before treatment
    first_penalty_ratio = get_governable_parameter_value(client, 'slashFractionDuplicateSign')
    return first_pledge_amount, first_penalty_ratio


def duplicate_sign(client, report_address, report_block):
    if report_block <= client.economic.consensus_size:
        report_block = client.economic.consensus_size + 1
    # Report1 prepareblock signature
    report_information = mock_duplicate_sign(1, client.node.nodekey, client.node.blsprikey, report_block)
    log.info("Report information: {}".format(report_information))
    result = client.duplicatesign.reportDuplicateSign(1, report_information, report_address)
    assert_code(result, 0)


def assret_penalty_amount(clients, first_pledge_amount, penalty_ratio=None):
    # view Pledge amount after punishment
    proportion_reward, incentive_pool_reward = clients[1].economic.get_report_reward(first_pledge_amount, penalty_ratio)
    log.info("Whistleblower benefits：{} Incentive pool income：{}".format(proportion_reward, incentive_pool_reward))
    # view Pledge amount again
    candidate_info = clients[1].ppos.getCandidateInfo(clients[0].node.node_id)
    second_pledge_amount = candidate_info['Ret']['Released']
    assert second_pledge_amount == first_pledge_amount - proportion_reward - incentive_pool_reward, "ErrMsg:Pledge amount {}".format(
        second_pledge_amount)


@pytest.mark.P1
@pytest.mark.parametrize('mark', [False, True])
def test_PIP_PVF_011_012(clients_consensus, mark, reset_environment):
    """
    治理修改区块双签-最高处罚比例投票失败
    :param clients_consensus:
    :param mark:
    :param reset_environment:
    :return:
    """
    first_index = 0
    first_client = clients_consensus[first_index]
    log.info("Current connection non-consensus node：{}".format(first_client.node.node_mark))
    economic = first_client.economic
    node = first_client.node
    change_parameter_value = '1000'
    # view Pledge amount and Parameter value before treatment
    first_pledge_amount, first_penalty_ratio = information_before_penalty_ratio(first_client)
    # create Parametric proposal
    param_governance_verify_before_endblock(first_client, 'slashing', 'slashFractionDuplicateSign',
                                            change_parameter_value, mark)
    # view Parameter value after treatment again
    second_penalty_ratio = get_governable_parameter_value(first_client, 'slashFractionDuplicateSign')
    assert second_penalty_ratio == first_penalty_ratio, "ErrMsg:Parameter value after treatment {}".format(
        second_penalty_ratio)
    # wait consensus block
    economic.wait_consensus_blocknum(node)
    # create account
    report_address, _ = economic.account.generate_account(node.web3, node.web3.toWei(1000, 'ether'))
    # Verify changed parameters
    current_block = node.eth.blockNumber
    log.info("Current block height: {}".format(current_block))
    # Verify changed parameters report
    duplicate_sign(first_client, report_address, current_block)
    # assret penalty amount
    assret_penalty_amount(clients_consensus, first_pledge_amount)


@pytest.mark.P1
def test_PIP_PVF_013(clients_consensus, reset_environment):
    """
    治理修改区块双签-最高处罚比例处于已生效期
    :param clients_consensus:
    :param reset_environment:
    :return:
    """
    first_index = 0
    first_client = clients_consensus[first_index]
    log.info("Current connection non-consensus node：{}".format(first_client.node.node_mark))
    economic = first_client.economic
    node = first_client.node
    change_parameter_value = '1000'
    # view Pledge amount and Parameter value before treatment
    first_pledge_amount, first_penalty_ratio = information_before_penalty_ratio(first_client)
    # create Parametric proposal
    param_governance_verify(first_client, 'slashing', 'slashFractionDuplicateSign', change_parameter_value)
    # view Parameter value before treatment again
    second_penalty_ratio = get_governable_parameter_value(first_client, 'slashFractionDuplicateSign')
    assert second_penalty_ratio == change_parameter_value, "ErrMsg:Parameter value after treatment {}".format(
        second_penalty_ratio)
    # create account
    report_address, _ = economic.account.generate_account(node.web3, node.web3.toWei(1000, 'ether'))
    # Verify changed parameters
    current_block = node.eth.blockNumber
    log.info("Current block height: {}".format(current_block))
    # Verify changed parameters report
    duplicate_sign(first_client, report_address, current_block)
    # assret penalty amount
    assret_penalty_amount(clients_consensus, first_pledge_amount, int(change_parameter_value))


@pytest.mark.P1
def test_PIP_PVF_014(clients_consensus, reset_environment):
    """
    治理修改区块双签-最高处罚比例为10000‱
    :param clients_consensus:
    :param reset_environment:
    :return:
    """
    first_index = 0
    first_client = clients_consensus[first_index]
    log.info("Current connection non-consensus node：{}".format(first_client.node.node_mark))
    economic = first_client.economic
    node = first_client.node
    change_parameter_value = '10000'
    # view Pledge amount and Parameter value before treatment
    first_pledge_amount, first_penalty_ratio = information_before_penalty_ratio(first_client)
    # create Parametric proposal
    param_governance_verify(first_client, 'slashing', 'slashFractionDuplicateSign', change_parameter_value)
    # view Parameter value before treatment again
    second_penalty_ratio = get_governable_parameter_value(first_client, 'slashFractionDuplicateSign')
    assert second_penalty_ratio == change_parameter_value, "ErrMsg:Parameter value after treatment {}".format(
        second_penalty_ratio)
    # create account
    report_address, _ = economic.account.generate_account(node.web3, node.web3.toWei(1000, 'ether'))
    # Verify changed parameters
    current_block = node.eth.blockNumber
    log.info("Current block height: {}".format(current_block))
    # Verify changed parameters report
    duplicate_sign(first_client, report_address, current_block)
    # assret penalty amount
    assret_penalty_amount(clients_consensus, first_pledge_amount, int(change_parameter_value))


@pytest.mark.P1
def test_PIP_PVF_015(clients_consensus, reset_environment):
    """
    治理修改区块双签-最高处罚比例为1‱
    :param clients_consensus:
    :param reset_environment:
    :return:
    """
    first_index = 0
    first_client = clients_consensus[first_index]
    log.info("Current connection non-consensus node：{}".format(first_client.node.node_mark))
    economic = first_client.economic
    node = first_client.node
    change_parameter_value = '1'
    # view Pledge amount and Parameter value before treatment
    first_pledge_amount, first_penalty_ratio = information_before_penalty_ratio(first_client)
    # create Parametric proposal
    param_governance_verify(first_client, 'slashing', 'slashFractionDuplicateSign', change_parameter_value)
    # view Parameter value before treatment again
    second_penalty_ratio = get_governable_parameter_value(first_client, 'slashFractionDuplicateSign')
    assert second_penalty_ratio == change_parameter_value, "ErrMsg:Parameter value after treatment {}".format(
        second_penalty_ratio)
    # create account
    report_address, _ = economic.account.generate_account(node.web3, node.web3.toWei(1000, 'ether'))
    # Verify changed parameters
    current_block = node.eth.blockNumber
    log.info("Current block height: {}".format(current_block))
    # Verify changed parameters report
    duplicate_sign(first_client, report_address, current_block)
    # assret penalty amount
    assret_penalty_amount(clients_consensus, first_pledge_amount, int(change_parameter_value))


def information_before_report_reward(client):
    # view Pledge amount
    first_candidate_info = client.ppos.getCandidateInfo(client.node.node_id)
    first_pledge_amount = first_candidate_info['Ret']['Released']
    # view Parameter value before treatment
    first_report_reward = get_governable_parameter_value(client, 'duplicateSignReportReward')
    return first_pledge_amount, first_report_reward


def get_account_amount(client):
    # create report account
    report_address, _ = client.economic.account.generate_account(client.node.web3, client.node.web3.toWei(
        1000, 'ether'))
    # view report amount
    first_report_amount = client.node.eth.getBalance(report_address)
    # view Incentive pool account
    first_incentive_pool_account = client.node.eth.getBalance(EconomicConfig.INCENTIVEPOOL_ADDRESS)
    return report_address, first_report_amount, first_incentive_pool_account


def asster_income_account_amount(client, first_report_amount, first_incentive_pool_account, report_address,
                                 proportion_reward,
                                 incentive_pool_reward):
    # view report amount
    second_report_amount = client.node.eth.getBalance(report_address)
    # view Incentive pool account
    second_incentive_pool_account = client.node.eth.getBalance(EconomicConfig.INCENTIVEPOOL_ADDRESS)
    # asster amount reward
    log.info("first_report_amount {} ,proportion_reward {} , second_report_amount {}".format(first_report_amount,
                                                                                             proportion_reward,
                                                                                             second_report_amount))
    assert first_report_amount + proportion_reward - second_report_amount < client.node.web3.toWei(1,
                                                                                                   'ether'), "ErrMsg:report amount {}".format(
        second_report_amount)
    log.info("second_incentive_pool_account {} ,first_incentive_pool_account {} , incentive_pool_reward {}".format(
        second_incentive_pool_account, first_incentive_pool_account, incentive_pool_reward))
    assert second_incentive_pool_account == first_incentive_pool_account + incentive_pool_reward + (
        first_report_amount + proportion_reward - second_report_amount), "ErrMsg:Incentive pool account {}".format(
        second_incentive_pool_account)


@pytest.mark.P1
@pytest.mark.parametrize('mark', [False, True])
def test_PIP_PVF_016_017(clients_consensus, mark, reset_environment):
    """
    PIP_PVF_016:治理修改区块双签-举报奖励比例投票失败
    PIP_PVF_017:治理修改区块双签-举报奖励比例处于未生效期
    :param clients_consensus:
    :param reset_environment:
    :return:
    """
    first_index = 0
    second_index = 1
    first_client = clients_consensus[first_index]
    log.info("Current connection non-consensus first node：{}".format(first_client.node.node_mark))
    second_client = clients_consensus[second_index]
    log.info("Current connection non-consensus second node：{}".format(second_client.node.node_mark))
    economic = first_client.economic
    node = first_client.node
    change_parameter_value = '60'
    first_pledge_amount, first_report_reward = information_before_report_reward(first_client)
    # create Parametric proposal
    param_governance_verify_before_endblock(first_client, 'slashing', 'duplicateSignReportReward', change_parameter_value, mark)
    # view Parameter value after treatment
    second_report_reward = get_governable_parameter_value(first_client, 'duplicateSignReportReward')
    assert first_report_reward == second_report_reward, "ErrMsg:Parameter value after treatment {}".format(second_report_reward)
    # wait consensus block
    economic.wait_consensus_blocknum(node)
    # get account amount
    report_address, first_report_amount, first_incentive_pool_account = get_account_amount(first_client)
    # Verify changed parameters
    current_block = node.eth.blockNumber
    log.info("Current block height: {}".format(current_block))
    # Verify changed parameters report
    duplicate_sign(first_client, report_address, current_block)
    # view Pledge amount after punishment
    proportion_reward, incentive_pool_reward = economic.get_report_reward(first_pledge_amount)
    log.info("Whistleblower benefits：{} Incentive pool income：{}".format(proportion_reward, incentive_pool_reward))
    # asster account amount
    asster_income_account_amount(second_client, first_report_amount, first_incentive_pool_account,
                                 report_address, proportion_reward, incentive_pool_reward)


@pytest.mark.P1
def testt(client_consensus):
    a = client_consensus.node.eth.getTransactionCount(client_consensus.economic.account.account_with_money['address'])
    print(a)


@pytest.mark.P1
def test_PIP_PVF_018(clients_consensus, reset_environment):
    """
    治理修改区块双签-举报奖励比例处于已生效期
    :param clients_consensus:
    :param reset_environment:
    :return:
    """
    first_index = 0
    second_index = 1
    first_client = clients_consensus[first_index]
    log.info("Current connection non-consensus first node：{}".format(first_client.node.node_mark))
    second_client = clients_consensus[second_index]
    log.info("Current connection non-consensus second node：{}".format(second_client.node.node_mark))
    node = first_client.node
    change_parameter_value = '60'
    # get first_pledge_amount first_report_amount first_incentive_pool_account first_report_reward
    first_pledge_amount, first_report_reward = information_before_report_reward(first_client)
    # create Parametric proposal
    param_governance_verify(first_client, 'slashing', 'duplicateSignReportReward', change_parameter_value)
    # view Parameter value after treatment
    second_report_reward = get_governable_parameter_value(first_client, 'duplicateSignReportReward')
    assert second_report_reward == change_parameter_value, "ErrMsg:Parameter value after treatment {}".format(second_report_reward)
    # get account amount
    report_address, first_report_amount, first_incentive_pool_account = get_account_amount(first_client)
    # Verify changed parameters
    current_block = node.eth.blockNumber
    log.info("Current block height: {}".format(current_block))
    # Verify changed parameters report
    duplicate_sign(first_client, report_address, current_block)
    # view Pledge amount after punishment
    proportion_reward, incentive_pool_reward = second_client.economic.get_report_reward(first_pledge_amount, None, int(change_parameter_value))
    log.info("Whistleblower benefits：{} Incentive pool income：{}".format(proportion_reward, incentive_pool_reward))
    # asster account amount
    asster_income_account_amount(second_client, first_report_amount, first_incentive_pool_account,
                                 report_address, proportion_reward, incentive_pool_reward)


@pytest.mark.P1
def test_PIP_PVF_019(clients_consensus, reset_environment):
    """
    治理修改区块双签-举报奖励比例为80%
    :param clients_consensus:
    :param reset_environment:
    :return:
    """
    first_index = 0
    second_index = 1
    first_client = clients_consensus[first_index]
    log.info("Current connection non-consensus first node：{}".format(first_client.node.node_mark))
    second_client = clients_consensus[second_index]
    log.info("Current connection non-consensus second node：{}".format(second_client.node.node_mark))
    economic = first_client.economic
    node = first_client.node
    change_parameter_value = '80'
    # get first_pledge_amount first_report_reward
    first_pledge_amount, first_report_reward = information_before_report_reward(first_client)
    # create Parametric proposal
    param_governance_verify(first_client, 'slashing', 'duplicateSignReportReward', change_parameter_value)
    # view Parameter value after treatment
    second_report_reward = get_governable_parameter_value(first_client, 'duplicateSignReportReward')
    assert second_report_reward == change_parameter_value, "ErrMsg:Parameter value after treatment {}".format(second_report_reward)
    # get account amount
    report_address, first_report_amount, first_incentive_pool_account = get_account_amount(first_client)
    # Verify changed parameters
    current_block = node.eth.blockNumber
    log.info("Current block height: {}".format(current_block))
    # Verify changed parameters report
    duplicate_sign(first_client, report_address, current_block)
    # view Pledge amount after punishment
    proportion_reward, incentive_pool_reward = economic.get_report_reward(first_pledge_amount, None, int(change_parameter_value))
    log.info("Whistleblower benefits：{} Incentive pool income：{}".format(proportion_reward, incentive_pool_reward))
    # asster account amount
    asster_income_account_amount(second_client, first_report_amount, first_incentive_pool_account,
                                 report_address, proportion_reward, incentive_pool_reward)


@pytest.mark.P1
def test_PIP_PVF_020(clients_consensus, reset_environment):
    """
    治理修改区块双签-举报奖励比例为1%
    :param clients_consensus:
    :param reset_environment:
    :return:
    """
    first_index = 0
    second_index = 1
    first_client = clients_consensus[first_index]
    log.info("Current connection non-consensus first node：{}".format(first_client.node.node_mark))
    second_client = clients_consensus[second_index]
    log.info("Current connection non-consensus second node：{}".format(second_client.node.node_mark))
    economic = first_client.economic
    node = first_client.node
    change_parameter_value = '1'
    # get first_pledge_amount first_report_reward
    first_pledge_amount, first_report_reward = information_before_report_reward(first_client)
    # create Parametric proposal
    param_governance_verify(first_client, 'slashing', 'duplicateSignReportReward', change_parameter_value)
    # view Parameter value after treatment
    second_report_reward = get_governable_parameter_value(first_client, 'duplicateSignReportReward')
    assert second_report_reward == change_parameter_value, "ErrMsg:Parameter value after treatment {}".format(second_report_reward)
    # get account amount
    report_address, first_report_amount, first_incentive_pool_account = get_account_amount(first_client)
    # Verify changed parameters
    current_block = node.eth.blockNumber
    log.info("Current block height: {}".format(current_block))
    # Verify changed parameters report
    duplicate_sign(first_client, report_address, current_block)
    # view Pledge amount after punishment
    proportion_reward, incentive_pool_reward = economic.get_report_reward(first_pledge_amount, None, int(change_parameter_value))
    log.info("Whistleblower benefits：{} Incentive pool income：{}".format(proportion_reward, incentive_pool_reward))
    # asster account amount
    asster_income_account_amount(second_client, first_report_amount, first_incentive_pool_account,
                                 report_address, proportion_reward, incentive_pool_reward)


@pytest.mark.parametrize('mark', [False, True])
@pytest.mark.P1
def test_PIP_MG_001_002(client_consensus, mark, reset_environment):
    """
    PIP_MG_001:治理修改默认每个区块的最大Gas 投票失败
    PIP_MG_002:治理修改默认每个区块的最大Gas 处于未生效期
    :param client_consensus:
    :param reset_environment:
    :return:
    """
    first_client = client_consensus
    log.info("Current connection non-consensus first node：{}".format(first_client.node.node_mark))
    node = first_client.node
    # view Parameter value before treatment
    first_max_gas_limit = get_governable_parameter_value(first_client, 'maxBlockGasLimit')
    # create Parametric proposal
    block = param_governance_verify_before_endblock(first_client, 'block', 'maxBlockGasLimit', '4712389', mark)
    # view Parameter value after treatment
    second_max_gas_limit = get_governable_parameter_value(first_client, 'maxBlockGasLimit')
    # wait block
    wait_block_number(node, block)
    assert second_max_gas_limit == first_max_gas_limit, "ErrMsg:Parameter value after treatment {}".format(second_max_gas_limit)


@pytest.mark.P1
def test_PIP_MG_003(client_consensus, reset_environment):
    """
    治理修改默认每个区块的最大Gas 处于已生效期
    :param client_consensus:
    :param reset_environment:
    :return:
    """
    first_client = client_consensus
    log.info("Current connection non-consensus first node：{}".format(first_client.node.node_mark))
    node = first_client.node
    change_parameter_value = '4712389'
    # view Parameter value before treatment
    first_max_gas_limit = get_governable_parameter_value(first_client, 'maxBlockGasLimit')
    # create Parametric proposal
    param_governance_verify(first_client, 'block', 'maxBlockGasLimit', change_parameter_value)
    # view Parameter value after treatment
    second_max_gas_limit = get_governable_parameter_value(first_client, 'maxBlockGasLimit')
    assert second_max_gas_limit != first_max_gas_limit, "ErrMsg:Parameter value after treatment {}".format(second_max_gas_limit)
    assert second_max_gas_limit == change_parameter_value, "ErrMsg:Parameter value after treatment {}".format(second_max_gas_limit)
