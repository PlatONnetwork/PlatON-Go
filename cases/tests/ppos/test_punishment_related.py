import time
import pytest
import rlp

from common.key import mock_duplicate_sign, generate_key
from common.log import log
from tests.lib import (
    EconomicConfig,
    StakingConfig,
    check_node_in_list,
    assert_code, von_amount,
    get_governable_parameter_value,
    Client,
    update_param_by_dict,
    get_param_by_dict,
    get_the_dynamic_parameter_gas_fee
    )


def penalty_proportion_and_income(client):
    # view Pledge amount
    candidate_info1 = client.ppos.getCandidateInfo(client.node.node_id)
    pledge_amount1 = candidate_info1['Ret']['Released']
    # view Parameter value before treatment
    penalty_ratio = get_governable_parameter_value(client, 'slashFractionDuplicateSign')
    proportion_ratio = get_governable_parameter_value(client, 'duplicateSignReportReward')
    return pledge_amount1, int(penalty_ratio), int(proportion_ratio)


def verification_duplicate_sign(client, evidence_type, reporting_type, report_address, report_block):
    if report_block < 41:
        report_block = 41
    # Obtain evidence of violation
    report_information = mock_duplicate_sign(evidence_type, client.node.nodekey,
                                             client.node.blsprikey,
                                             report_block)
    log.info("Report information: {}".format(report_information))
    result = client.duplicatesign.reportDuplicateSign(reporting_type, report_information, report_address)
    return result


@pytest.mark.P0
@pytest.mark.parametrize('repor_type', [1, 2, 3])
def test_VP_PV_001_to_003(client_consensus, repor_type, reset_environment):
    """
    举报验证人区块双签:VP_PV_001 prepareBlock类型
                    VP_PV_002 prepareVote类型
                    VP_PV_003 viewChange类型
    :param client_consensus:
    :param repor_type:
    :param reset_environment:
    :return:
    """
    client = client_consensus
    economic = client.economic
    node = client.node
    client.economic.env.deploy_all()
    # Obtain penalty proportion and income
    pledge_amount1, penalty_ratio, proportion_ratio = penalty_proportion_and_income(client)
    # create report address
    report_address, _ = economic.account.generate_account(node.web3, node.web3.toWei(1000, 'ether'))
    # view report amount
    report_amount1 = node.eth.getBalance(report_address)
    log.info("report account amount:{} ".format(report_amount1))
    # view Incentive pool account
    incentive_pool_account1 = node.eth.getBalance(EconomicConfig.INCENTIVEPOOL_ADDRESS)
    log.info("incentive pool account1 amount:{} ".format(incentive_pool_account1))
    # Wait for the consensus round to end
    economic.wait_consensus_blocknum(node)
    # Get current block height
    current_block = node.eth.blockNumber
    log.info("Current block height: {}".format(current_block))
    result = verification_duplicate_sign(client, repor_type, repor_type, report_address, current_block)
    assert_code(result, 0)
    # view Amount of penalty
    proportion_reward, incentive_pool_reward = economic.get_report_reward(pledge_amount1, penalty_ratio,
                                                                          proportion_ratio)
    # view report amount again
    report_amount2 = node.eth.getBalance(report_address)
    log.info("report account amount:{} ".format(report_amount2))
    # view Incentive pool account again
    incentive_pool_account2 = node.eth.getBalance(EconomicConfig.INCENTIVEPOOL_ADDRESS)
    log.info("incentive pool account1 amount:{} ".format(incentive_pool_account2))
    # assert account reward
    assert report_amount1 + proportion_reward - report_amount2 < node.web3.toWei(1,'ether'), "ErrMsg:report amount {}".format(
        report_amount2)
    assert incentive_pool_account2 == incentive_pool_account1 + incentive_pool_reward + (report_amount1 + proportion_reward - report_amount2), "ErrMsg:Incentive pool account {}".format(
        incentive_pool_account2)


@pytest.fixture(scope='class')
def initial_report(global_test_env):
    """
    Report a non consensus node prepareBlock
    :param global_test_env:
    :return:
    """
    cfg = StakingConfig("11111", "faker", "www.baidu.com", "how much")
    clients_consensus = []
    consensus_node_obj_list = global_test_env.consensus_node_list
    for node_obj in consensus_node_obj_list:
        clients_consensus.append(Client(global_test_env, node_obj, cfg))
    client = clients_consensus[0]
    economic = client.economic
    node = client.node
    # create report address
    report_address, _ = economic.account.generate_account(node.web3, node.web3.toWei(1000, 'ether'))
    # Wait for the consensus round to end
    economic.wait_consensus_blocknum(node, 2)
    # Get current block height
    current_block = node.eth.blockNumber
    log.info("Current block height: {}".format(current_block))
    result = verification_duplicate_sign(client, 1, 1, report_address, current_block)
    assert_code(result, 0)
    yield clients_consensus, economic, node, report_address, current_block
    log.info("case execution completed")
    global_test_env.deploy_all()
    time.sleep(3)


class TestMultipleReports:
    @pytest.mark.P1
    def test_VP_PV_004(self, initial_report):
        """
        举报双签-同一验证人同一块高不同类型
        :param initial_report:
        :return:
        """
        clients_consensus, economic, node, report_address, current_block = initial_report
        # duplicate sign
        result = verification_duplicate_sign(clients_consensus[0], 2, 2, report_address, current_block)
        assert_code(result, 0)

    @pytest.mark.P1
    def test_VP_PV_005(self, initial_report):
        """
        举报双签-同一验证人不同块高同一类型
        :param initial_report:
        :return:
        """
        clients_consensus, economic, node, report_address, current_block = initial_report
        # duplicate sign
        result = verification_duplicate_sign(clients_consensus[0], 1, 1, report_address, current_block - 1)
        assert_code(result, 0)

    @pytest.mark.P1
    def test_VP_PV_006(self, initial_report):
        """
        举报双签-同一验证人不同块高不同类型
        :param initial_report:
        :return:
        """
        clients_consensus, economic, node, report_address, current_block = initial_report
        # duplicate sign
        result = verification_duplicate_sign(clients_consensus[0], 2, 2, report_address, current_block - 1)
        assert_code(result, 0)

    @pytest.mark.P1
    def test_VP_PV_007(self, initial_report):
        """
        举报双签-不同验证人同一块高同一类型
        :param initial_report:
        :return:
        """
        clients_consensus, economic, node, report_address, current_block = initial_report
        # create account
        report_address2, _ = economic.account.generate_account(node.web3, node.web3.toWei(1000, 'ether'))
        # duplicate sign
        result = verification_duplicate_sign(clients_consensus[1], 1, 1, report_address, current_block)
        assert_code(result, 0)

    @pytest.mark.P1
    def test_VP_PV_008(self, initial_report):
        """
        举报双签-不同验证人同一块高不同类型
        :param initial_report:
        :return:
        """
        clients_consensus, economic, node, report_address, current_block = initial_report
        # create account
        report_address2, _ = economic.account.generate_account(node.web3, node.web3.toWei(1000, 'ether'))
        # duplicate sign
        result = verification_duplicate_sign(clients_consensus[1], 2, 2, report_address, current_block)
        assert_code(result, 0)

    @pytest.mark.P1
    def test_VP_PV_009(self, initial_report):
        """
        举报双签-不同验证人不同块高不同类型
        :param initial_report:
        :return:
        """
        clients_consensus, economic, node, report_address, current_block = initial_report
        # create account
        report_address2, _ = economic.account.generate_account(node.web3, node.web3.toWei(1000, 'ether'))
        # duplicate sign
        result = verification_duplicate_sign(clients_consensus[1], 2, 2, report_address, current_block - 1)
        assert_code(result, 0)

    @pytest.mark.P1
    def test_VP_PR_001(self, initial_report):
        """
        重复举报-同一举报人
        :param initial_report:
        :return:
        """
        clients_consensus, economic, node, report_address, current_block = initial_report
        # duplicate sign
        result = verification_duplicate_sign(clients_consensus[0], 1, 1, report_address, current_block)
        assert_code(result, 303001)

    @pytest.mark.P1
    def test_VP_PR_002(self, initial_report):
        """
        重复举报 - 不同举报人
        :param initial_report:
        :return:
        """
        clients_consensus, economic, node, report_address, current_block = initial_report
        # create account
        report_address2, _ = economic.account.generate_account(node.web3, node.web3.toWei(1000, 'ether'))
        # duplicate sign
        result = verification_duplicate_sign(clients_consensus[0], 1, 1, report_address2, current_block)
        assert_code(result, 303001)


def obtaining_evidence_information(economic, node):
    # Wait for the consensus round to end
    economic.wait_consensus_blocknum(node, 1)
    # Get current block height
    current_block = node.eth.blockNumber
    log.info("Current block height: {}".format(current_block))
    if current_block < 41:
        current_block = 41
    # Obtain evidence of violation
    report_information = mock_duplicate_sign(1, node.nodekey, node.blsprikey, current_block)
    log.info("Report information: {}".format(report_information))
    return report_information, current_block


@pytest.mark.P1
@pytest.mark.parametrize('first_key, second_key, value',
                         [('epoch', None, 1), ('viewNumber', None, 1), ('blockIndex', None, 1),
                          ('validateNode', 'index', 1)])
def test_VP_PV_010_011_014_015(client_consensus, first_key, second_key, value):
    """
    VP_PV_010:举报双签-双签证据epoch不一致
    VP_PV_011:举报双签-双签证据view_number不一致
    VP_PV_014:举报双签-双签证据block_index不一致
    VP_PV_015:举报双签-双签证据validate_node-index不一致
    :param client_consensus:
    :param first_key:
    :param second_key:
    :param value:
    :return:
    """
    client = client_consensus
    economic = client.economic
    node = client.node
    # create report address
    report_address, _ = economic.account.generate_account(node.web3, node.web3.toWei(1000, 'ether'))
    # Obtain information of report evidence
    report_information, current_block = obtaining_evidence_information(economic, node)
    # Modification of evidence
    jsondata = update_param_by_dict(report_information, 'prepareA', first_key, second_key, value)
    log.info("Evidence information: {}".format(jsondata))
    # Report verifier Duplicate Sign
    result = client.duplicatesign.reportDuplicateSign(1, jsondata, report_address)
    return result


@pytest.mark.P1
def test_VP_PV_012(client_consensus):
    """
    举报双签-双签证据block_number不一致
    :param client_consensus:
    :return:
    """
    client = client_consensus
    economic = client.economic
    node = client.node
    # create report address
    report_address, _ = economic.account.generate_account(node.web3, node.web3.toWei(1000, 'ether'))
    # Obtain information of report evidence
    report_information, current_block = obtaining_evidence_information(economic, node)
    # Modification of evidence
    jsondata = update_param_by_dict(report_information, 'prepareA', 'blockNumber', None, current_block - 1)
    log.info("Evidence information: {}".format(jsondata))
    # Report verifier Duplicate Sign
    result = client.duplicatesign.reportDuplicateSign(1, jsondata, report_address)
    assert_code(result, 303000)


@pytest.mark.P1
def test_VP_PV_013(client_consensus):
    """
    举报双签-双签证据block_hash一致
    :param client_consensus:
    :return:
    """
    client = client_consensus
    economic = client.economic
    node = client.node
    # create report address
    report_address, _ = economic.account.generate_account(node.web3, node.web3.toWei(1000, 'ether'))
    # Obtain information of report evidence
    report_information, current_block = obtaining_evidence_information(economic, node)
    # Modification of evidence
    evidence_parameter = get_param_by_dict(report_information, 'prepareB', 'blockHash')
    jsondata = update_param_by_dict(report_information, 'prepareA', 'blockHash', None, evidence_parameter)
    log.info("Evidence information: {}".format(jsondata))
    # Report verifier Duplicate Sign
    result = client.duplicatesign.reportDuplicateSign(1, jsondata, report_address)
    assert_code(result, 303000)


@pytest.mark.P1
def test_VP_PV_016(client_consensus):
    """
    举报双签-双签证据address不一致
    :param client_consensus:
    :return:
    """
    pass
#     client = client_consensus
#     economic = client.economic
#     node = client.node
#     # create report address
#     report_address, _ = economic.account.generate_account(node.web3, node.web3.toWei(1000, 'ether'))
#     # Obtain information of report evidence
#     report_information, current_block = obtaining_evidence_information(economic, node)
#     # Modification of evidence
#     jsondata = update_param_by_dict(report_information, 'prepareA', 'validateNode', 'address',
#                                     economic.account.account_with_money['address'])
#     log.info("Evidence information: {}".format(jsondata))
#     # Report verifier Duplicate Sign
#     result = client.duplicatesign.reportDuplicateSign(1, jsondata, report_address)
#     assert_code(result, 303000)


@pytest.mark.P1
def test_VP_PV_017(clients_consensus):
    """
    举报双签-NodeID不一致举报双签
    :param clients_consensus:
    :return:
    """
    client = clients_consensus[0]
    economic = client.economic
    node = client.node
    # create report address
    report_address, _ = economic.account.generate_account(node.web3, node.web3.toWei(1000, 'ether'))
    # Obtain information of report evidence
    report_information, current_block = obtaining_evidence_information(economic, node)
    # Modification of evidence
    jsondata = update_param_by_dict(report_information, 'prepareA', 'validateNode', 'nodeId',
                                    clients_consensus[1].node.node_id)
    log.info("Evidence information: {}".format(jsondata))
    # Report verifier Duplicate Sign
    result = client.duplicatesign.reportDuplicateSign(1, jsondata, report_address)
    assert_code(result, 303000)


@pytest.mark.P1
def test_VP_PV_018(clients_consensus):
    """
    举报双签-blsPubKey不一致举报双签
    :param clients_consensus:
    :return:
    """
    client = clients_consensus[0]
    economic = client.economic
    node = client.node
    # create report address
    report_address, _ = economic.account.generate_account(node.web3, node.web3.toWei(1000, 'ether'))
    # Obtain information of report evidence
    report_information, current_block = obtaining_evidence_information(economic, node)
    # Modification of evidence
    jsondata = update_param_by_dict(report_information, 'prepareA', 'validateNode', 'blsPubKey',
                                    clients_consensus[1].node.blspubkey)
    log.info("Evidence information: {}".format(jsondata))
    # Report verifier Duplicate Sign
    result = client.duplicatesign.reportDuplicateSign(1, jsondata, report_address)
    assert_code(result, 303000)


@pytest.mark.P1
def test_VP_PV_019(clients_consensus):
    """
    举报双签-signature一致举报双签
    :param clients_consensus:
    :return:
    """
    client = clients_consensus[0]
    economic = client.economic
    node = client.node
    # create report address
    report_address, _ = economic.account.generate_account(node.web3, node.web3.toWei(1000, 'ether'))
    # Obtain information of report evidence
    report_information, current_block = obtaining_evidence_information(economic, node)
    # Modification of evidence
    evidence_parameter = get_param_by_dict(report_information, 'prepareB', 'signature')
    jsondata = update_param_by_dict(report_information, 'prepareA', 'signature', None, evidence_parameter)
    log.info("Evidence information: {}".format(jsondata))
    # Report verifier Duplicate Sign
    result = client.duplicatesign.reportDuplicateSign(1, jsondata, report_address)
    assert_code(result, 303000)


@pytest.mark.P1
@pytest.mark.parametrize("value", [{"epoch": 1}, {"view_number": 1}, {"block_index": 1}, {"index": 1}])
def test_VP_PV_020_to_023(clients_consensus, value):
    """
    VP_PV_020:举报双签-伪造合法signature情况下伪造epoch
    VP_PV_021:举报双签-伪造合法signature情况下伪造viewNumber
    VP_PV_022:举报双签-伪造合法signature情况下伪造blockIndex
    VP_PV_023:举报双签-伪造合法signature情况下伪造index
    :param clients_consensus:
    :return:
    """
    client = clients_consensus[0]
    economic = client.economic
    node = client.node
    # create report address
    report_address, _ = economic.account.generate_account(node.web3, node.web3.toWei(1000, 'ether'))
    # Obtain information of report evidence
    report_information, current_block = obtaining_evidence_information(economic, node)
    # Obtain evidence of violation
    report_information1 = mock_duplicate_sign(1, node.nodekey, node.blsprikey, current_block, **value)
    log.info("Report information: {}".format(report_information))
    # Modification of evidence
    evidence_parameter = get_param_by_dict(report_information1, 'prepareB', 'signature')
    jsondata = update_param_by_dict(report_information, 'prepareA', 'signature', None, evidence_parameter)
    log.info("Evidence information: {}".format(jsondata))
    # Report verifier Duplicate Sign
    result = client.duplicatesign.reportDuplicateSign(1, jsondata, report_address)
    assert_code(result, 303000)


@pytest.mark.P1
def test_VP_PV_024(clients_consensus):
    """
    举报双签-伪造合法signature情况下伪造blockNumber
    :param clients_consensus:
    :return:
    """
    client = clients_consensus[0]
    economic = client.economic
    node = client.node
    # create report address
    report_address, _ = economic.account.generate_account(node.web3, node.web3.toWei(1000, 'ether'))
    # Obtain information of report evidence
    report_information, current_block = obtaining_evidence_information(economic, node)
    # Obtain evidence of violation
    report_information1 = mock_duplicate_sign(1, node.nodekey, node.blsprikey, current_block + 1)
    log.info("Report information: {}".format(report_information))
    # Modification of evidence
    evidence_parameter = get_param_by_dict(report_information1, 'prepareB', 'signature')
    jsondata = update_param_by_dict(report_information, 'prepareA', 'signature', None, evidence_parameter)
    log.info("Evidence information: {}".format(jsondata))
    # Report verifier Duplicate Sign
    result = client.duplicatesign.reportDuplicateSign(1, jsondata, report_address)
    assert_code(result, 303000)


@pytest.mark.P1
def test_VP_PV_025(client_consensus):
    """
    举报接口参数测试：举报人账户错误
    :param client_consensus:
    :return:
    """
    client = client_consensus
    economic = client.economic
    node = client.node
    # create report address
    report_address, _ = economic.account.generate_account(node.web3, node.web3.toWei(1000, 'ether'))
    # create nodekey
    privatekey, _ = generate_key()
    # Wait for the consensus round to end
    economic.wait_consensus_blocknum(node, 1)
    # Get current block height
    current_block = node.eth.blockNumber
    log.info("Current block height: {}".format(current_block))
    # Obtain evidence of violation
    report_information = mock_duplicate_sign(1, privatekey, node.blsprikey, current_block)
    log.info("Report information: {}".format(report_information))
    # Report verifier Duplicate Sign
    result = client.duplicatesign.reportDuplicateSign(1, report_information, report_address)
    assert_code(result, 303004)


@pytest.mark.P1
def test_VP_PV_026(clients_consensus):
    """
    链存在的id,blskey不匹配
    :param clients_consensus:
    :return:
    """
    client = clients_consensus[0]
    economic = client.economic
    node = client.node
    # create report address
    report_address, _ = economic.account.generate_account(node.web3, node.web3.toWei(1000, 'ether'))
    # Wait for the consensus round to end
    economic.wait_consensus_blocknum(node, 1)
    # Get current block height
    current_block = node.eth.blockNumber
    log.info("Current block height: {}".format(current_block))
    # Obtain evidence of violation
    report_information = mock_duplicate_sign(1, node.nodekey, clients_consensus[1].node.blsprikey, current_block)
    log.info("Report information: {}".format(report_information))
    # Report verifier Duplicate Sign
    result = client.duplicatesign.reportDuplicateSign(1, report_information, report_address)
    assert_code(result, 303007)


@pytest.mark.P1
def test_VP_PV_027(client_new_node):
    """
    举报候选人
    :param client_new_node:
    :return:
    """
    client = client_new_node
    economic = client.economic
    node = client.node
    # create pledge address
    pledge_address, _ = economic.account.generate_account(node.web3, von_amount(economic.create_staking_limit, 2))
    # create report address
    report_address, _ = economic.account.generate_account(node.web3, node.web3.toWei(1000, 'ether'))
    # create staking
    result = client.staking.create_staking(0, pledge_address, pledge_address)
    assert_code(result, 0)
    # Wait for the settlement round to end
    economic.wait_settlement_blocknum(node)
    for i in range(4):
        result = check_node_in_list(node.node_id, client.ppos.getValidatorList)
        log.info("Current node in consensus list status：{}".format(result))
        if not result:
            # Get current block height
            current_block = node.eth.blockNumber
            log.info("Current block height: {}".format(current_block))
            # Report verifier Duplicate Sign
            result = verification_duplicate_sign(client, 1, 1, report_address, current_block)
            assert_code(result, 303009)
        else:
            # wait consensus block
            economic.wait_consensus_blocknum(node)


@pytest.mark.P1
def test_VP_PV_028(client_consensus):
    """
    举报有效期之前的双签行为
    :param client_consensus:
    :return:
    """
    client = client_consensus
    economic = client.economic
    node = client.node
    # create report address
    report_address, _ = economic.account.generate_account(node.web3, node.web3.toWei(1000, 'ether'))
    # Waiting for the end of the settlement cycle
    economic.wait_settlement_blocknum(node, 1)
    # Get current block height
    current_block = node.eth.blockNumber
    log.info("Current block height: {}".format(current_block))
    # Obtain evidence of violation
    report_information = mock_duplicate_sign(1, node.nodekey, node.blsprikey, 41)
    log.info("Report information: {}".format(report_information))
    # Report verifier Duplicate Sign
    result = client.duplicatesign.reportDuplicateSign(1, report_information, report_address)
    assert_code(result, 303003)


@pytest.mark.P1
def test_VP_PV_028(client_consensus):
    """
    举报有效期之后的双签行为
    :param client_consensus:
    :return:
    """
    client = client_consensus
    economic = client.economic
    node = client.node
    # create report address
    report_address, _ = economic.account.generate_account(node.web3, node.web3.toWei(1000, 'ether'))
    # Obtain evidence of violation
    report_information = mock_duplicate_sign(1, node.nodekey, node.blsprikey, 1000000)
    log.info("Report information: {}".format(report_information))
    # Report verifier Duplicate Sign
    result = client.duplicatesign.reportDuplicateSign(1, report_information, report_address)
    assert_code(result, 303002)


@pytest.mark.P2
def test_VP_PV_030(client_consensus, reset_environment):
    """
    举报签名Gas费
    :param client_consensus:
    :return:
    """
    client = client_consensus
    economic = client.economic
    node = client.node
    client.economic.env.deploy_all()
    # create report address
    report_address, _ = economic.account.generate_account(node.web3, node.web3.toWei(1000, 'ether'))
    # Obtain information of report evidence
    report_information, current_block = obtaining_evidence_information(economic, node)
    # Obtain penalty proportion and income
    pledge_amount1, penalty_ratio, proportion_ratio = penalty_proportion_and_income(client)
    # view Amount of penalty
    proportion_reward, incentive_pool_reward = economic.get_report_reward(pledge_amount1, penalty_ratio,
                                                                          proportion_ratio)
    data = rlp.encode([rlp.encode(int(3000)), rlp.encode(1), rlp.encode(report_information)])
    dynamic_gas = get_the_dynamic_parameter_gas_fee(data)
    gas_total = 21000 + 21000 + 21000 + 21000 + dynamic_gas
    log.info("Call contract to create a lockout plan consumption contract：{}".format(gas_total))
    balance = node.eth.getBalance(report_address)
    log.info("balance: {}".format(balance))
    # Report verifier
    result = client.duplicatesign.reportDuplicateSign(1, report_information, report_address)
    assert_code(result, 0)
    balance1 = node.eth.getBalance(report_address)
    log.info("balance1: {}".format(balance1))
    log.info("proportion reward: {}".format(proportion_reward))
    transaction_fees = gas_total * node.eth.gasPrice
    assert balance + proportion_reward - balance1 == transaction_fees, "ErrMsg:transaction fees {}".format(
        transaction_fees)


@pytest.mark.P1
def test_VP_PV_031(client_consensus):
    """
    举报的gas费不足
    :param client_consensus:
    :return:
    """
    client = client_consensus
    economic = client.economic
    node = client.node
    status = True
    # create report address
    report_address, _ = economic.account.generate_account(node.web3, 0)
    # Obtain information of report evidence
    report_information, current_block = obtaining_evidence_information(economic, node)
    try:
        # Report verifier Duplicate Sign
        result = client.duplicatesign.reportDuplicateSign(1, report_information, report_address)
        log.info("result: {}".format(result))
        status = False
    except Exception as e:
        log.info("Use case success, exception information：{} ".format(str(e)))
    assert status, "ErrMsg:Report verifier status {}".format(status)


@pytest.mark.P1
def test_VP_PR_003(client_new_node, reset_environment):
    """
    举报被处罚退出状态中的验证人
    :param client_new_node:
    :return:
    """
    client = client_new_node
    economic = client.economic
    node = client.node
    # create pledge address
    pledge_address, _ = economic.account.generate_account(node.web3, von_amount(economic.create_staking_limit, 2))
    # create report address
    report_address, _ = economic.account.generate_account(node.web3, node.web3.toWei(1000, 'ether'))
    # create staking
    result = client.staking.create_staking(0, pledge_address, pledge_address)
    assert_code(result, 0)
    # view Pledge node information
    candidate_info = client.ppos.getCandidateInfo(node.node_id)
    log.info("Pledge node information: {}".format(candidate_info))
    # Wait for the settlement round to end
    economic.wait_settlement_blocknum(node)
    # Obtain penalty proportion and income
    pledge_amount1, penalty_ratio, proportion_ratio = penalty_proportion_and_income(client)
    # view Amount of penalty
    proportion_reward, incentive_pool_reward = economic.get_report_reward(pledge_amount1, penalty_ratio,
                                                                          proportion_ratio)
    for i in range(4):
        result = check_node_in_list(node.node_id, client.ppos.getValidatorList)
        log.info("Current node in consensus list status：{}".format(result))
        if result:
            # Application for return of pledge
            result = client.staking.withdrew_staking(pledge_address)
            assert_code(result, 0)
            # Get current block height
            current_block = node.eth.blockNumber
            log.info("Current block height: {}".format(current_block))
            # Report verifier Duplicate Sign
            result = verification_duplicate_sign(client, 1, 1, report_address, current_block)
            assert_code(result, 0)
            # view Pledge node information
            candidate_info = client.ppos.getCandidateInfo(node.node_id)
            log.info("Pledge node information: {}".format(candidate_info))
            log.info("Pledge node amount: {}".format(pledge_amount1))
            log.info("proportion_reward + incentive_pool_reward: {}".format(proportion_reward + incentive_pool_reward))
            info = candidate_info['Ret']
            assert info['Released'] == pledge_amount1 - (
                proportion_reward + incentive_pool_reward), "ErrMsg:Pledge amount {}".format(
                info['Released'])
            break
        else:
            # wait consensus block
            economic.wait_consensus_blocknum(node)


@pytest.mark.P1
def test_VP_PR_004(client_new_node):
    """
    举报已完成退出的验证人
    :param client_new_node:
    :return:
    """
    client = client_new_node
    economic = client.economic
    node = client.node
    # create pledge address
    pledge_address, _ = economic.account.generate_account(node.web3, von_amount(economic.create_staking_limit, 2))
    # create report address
    report_address, _ = economic.account.generate_account(node.web3, node.web3.toWei(1000, 'ether'))
    # create staking
    result = client.staking.create_staking(0, pledge_address, pledge_address)
    assert_code(result, 0)
    # Wait for the settlement round to end
    economic.wait_settlement_blocknum(node)
    for i in range(4):
        result = check_node_in_list(node.node_id, client.ppos.getValidatorList)
        log.info("Current node in consensus list status：{}".format(result))
        if result:
            # Application for return of pledge
            result = client.staking.withdrew_staking(pledge_address)
            assert_code(result, 0)
            # Waiting for the end of the 2 settlement cycle
            economic.wait_settlement_blocknum(node, 2)
            # Get current block height
            current_block = node.eth.blockNumber
            log.info("Current block height: {}".format(current_block))
            # Report verifier Duplicate Sign
            result = verification_duplicate_sign(client, 1, 1, report_address, current_block)
            assert_code(result, 303004)
            break
        else:
            # wait consensus block
            economic.wait_consensus_blocknum(node)


@pytest.mark.P1
def test_VP_PR_005(client_new_node, reset_environment):
    """
    举报人和被举报人为同一个人
    :param client_new_node:
    :return:
    """
    client = client_new_node
    economic = client.economic
    node = client.node
    # create pledge address
    pledge_address, _ = economic.account.generate_account(node.web3, von_amount(economic.create_staking_limit, 2))
    # create staking
    result = client.staking.create_staking(0, pledge_address, pledge_address)
    assert_code(result, 0)
    # Wait for the settlement round to end
    economic.wait_settlement_blocknum(node)
    for i in range(4):
        result = check_node_in_list(node.node_id, client.ppos.getValidatorList)
        log.info("Current node in consensus list status：{}".format(result))
        if result:
            # Application for return of pledge
            result = client.staking.withdrew_staking(pledge_address)
            assert_code(result, 0)
            # Get current block height
            current_block = node.eth.blockNumber
            log.info("Current block height: {}".format(current_block))
            # Report verifier Duplicate Sign
            result = verification_duplicate_sign(client, 1, 1, pledge_address, current_block)
            assert_code(result, 303010)
            break
        else:
            # wait consensus block
            economic.wait_consensus_blocknum(node)


@pytest.mark.P1
def test_VP_PVF_001(client_consensus, reset_environment):
    """
    查询已成功的举报
    :param client_consensus:
    :return:
    """
    client = client_consensus
    economic = client.economic
    node = client.node
    # create report address
    report_address, _ = economic.account.generate_account(node.web3, node.web3.toWei(1000, 'ether'))
    # Wait for the consensus round to end
    economic.wait_consensus_blocknum(node, 1)
    # Get current block height
    current_block = node.eth.blockNumber
    log.info("Current block height: {}".format(current_block))
    # Obtain evidence of violation
    report_information = mock_duplicate_sign(1, node.nodekey, node.blsprikey, current_block)
    log.info("Report information: {}".format(report_information))
    result = client.duplicatesign.reportDuplicateSign(1, report_information, report_address)
    assert_code(result, 0)
    # Query and report violation records
    evidence_parameter = get_param_by_dict(report_information, 'prepareA', 'validateNode', 'nodeId')
    result = client.ppos.checkDuplicateSign(1, evidence_parameter, current_block)
    assert_code(result, 0)
    assert result['Ret'] is not None, "ErrMsg:Query results {}".format(result['Ret'])


@pytest.mark.P1
def test_VP_PVF_002(client_consensus):
    """
    查询未成功的举报记录
    :param client_consensus:
    :return:
    """
    client = client_consensus
    economic = client.economic
    node = client.node
    # create report address
    report_address, _ = economic.account.generate_account(node.web3, node.web3.toWei(1000, 'ether'))
    # Wait for the consensus round to end
    economic.wait_consensus_blocknum(node, 1)
    # Get current block height
    current_block = node.eth.blockNumber
    log.info("Current block height: {}".format(current_block))
    # Obtain evidence of violation
    report_information = mock_duplicate_sign(1, node.nodekey, node.blsprikey, 100000)
    log.info("Report information: {}".format(report_information))
    result = client.duplicatesign.reportDuplicateSign(1, report_information, report_address)
    assert_code(result, 303002)
    # create account
    report_address2, _ = economic.account.generate_account(node.web3, node.web3.toWei(1000, 'ether'))
    # Query and report violation records
    evidence_parameter = get_param_by_dict(report_information, 'prepareA', 'validateNode', 'nodeId')
    result = client.ppos.checkDuplicateSign(1, evidence_parameter, current_block)
    assert_code(result, 0)
    assert result['Ret'] == "", "ErrMsg:Query results {}".format(result['Ret'])


@pytest.mark.P1
def test_VP_PVF_003(client_new_node, reset_environment):
    """
    被系统剔除出验证人与候选人名单，节点可继续完成轮的出块及验证工作
    :param client_new_node:
    :return:
    """
    client = client_new_node
    economic = client.economic
    node = client.node
    time.sleep(5)
    # create pledge address
    pledge_address, _ = economic.account.generate_account(node.web3, von_amount(economic.create_staking_limit, 2))
    # create report address
    report_address, _ = economic.account.generate_account(node.web3, node.web3.toWei(1000, 'ether'))
    # create staking
    result = client.staking.create_staking(0, pledge_address, pledge_address)
    assert_code(result, 0)
    # Wait for the settlement round to end
    economic.wait_settlement_blocknum(node)
    for i in range(4):
        result = check_node_in_list(node.node_id, client.ppos.getValidatorList)
        log.info("Current node in consensus list status：{}".format(result))
        if result:
            # Get current block height
            current_block = node.eth.blockNumber
            log.info("Current block height: {}".format(current_block))
            # Report verifier Duplicate Sign
            result = verification_duplicate_sign(client, 1, 1, report_address, current_block)
            assert_code(result, 0)
            result = check_node_in_list(node.node_id, client.ppos.getValidatorList)
            log.info("Current node in consensus list status：{}".format(result))
            assert result, "ErrMsg:Node current status {}".format(result)
            # Wait for the settlement round to end
            economic.wait_consensus_blocknum(node, 2)
            result = check_node_in_list(node.node_id, client.ppos.getValidatorList)
            log.info("Current node in consensus list status：{}".format(result))
            assert not result, "ErrMsg:Node current status {}".format(result)
            break
        else:
            # wait consensus block
            economic.wait_consensus_blocknum(node)


@pytest.mark.P1
def test_VP_PVF_004(client_new_node, reset_environment):
    """
    验证人在共识轮第230区块前被举报并被处罚
    :param client_new_node:
    :return:
    """
    client = client_new_node
    economic = client.economic
    node = client.node
    # create pledge address
    pledge_address, _ = economic.account.generate_account(node.web3, von_amount(economic.create_staking_limit, 2))
    # create report address
    report_address, _ = economic.account.generate_account(node.web3, node.web3.toWei(1000, 'ether'))
    # create staking
    result = client.staking.create_staking(0, pledge_address, pledge_address)
    assert_code(result, 0)
    # Wait for the settlement round to end
    economic.wait_settlement_blocknum(node)
    endtime = int(time.time()) + 120
    while int(time.time()) < endtime:
        time.sleep(1)
        current_block = node.eth.blockNumber
        log.info("current block: {}".format(current_block))
        block = current_block % economic.consensus_size
        log.info("block: {}".format(block))
        log.info("Current block height: {}, block of current consensus round: {}".format(current_block, block))
        if block < 20:
            break
    for i in range(4):
        result = check_node_in_list(node.node_id, client.ppos.getValidatorList)
        log.info("Current node in consensus list status：{}".format(result))
        if result:
            # Get current block height
            current_block = node.eth.blockNumber
            log.info("Current block height: {}".format(current_block))
            # Report verifier Duplicate Sign
            result = verification_duplicate_sign(client, 1, 1, report_address, current_block)
            assert_code(result, 0)
            result = check_node_in_list(node.node_id, client.ppos.getValidatorList)
            log.info("Current node in consensus list status：{}".format(result))
            assert result, "ErrMsg:Node current status {}".format(result)
            # Wait for the settlement round to end
            economic.wait_consensus_blocknum(node)
            result = check_node_in_list(node.node_id, client.ppos.getValidatorList)
            log.info("Current node in consensus list status：{}".format(result))
            assert not result, "ErrMsg:Node current status {}".format(result)
            break
        else:
            # wait consensus block
            economic.wait_consensus_blocknum(node)


@pytest.mark.P1
def test_VP_PVF_005(client_new_node, reset_environment):
    """
    验证人在共识轮第230区块后举报并被处罚
    :param client_new_node:
    :return:
    """
    client = client_new_node
    economic = client.economic
    node = client.node
    # create pledge address
    pledge_address, _ = economic.account.generate_account(node.web3, von_amount(economic.create_staking_limit, 2))
    # create report address
    report_address, _ = economic.account.generate_account(node.web3, node.web3.toWei(1000, 'ether'))
    # create staking
    result = client.staking.create_staking(0, pledge_address, pledge_address)
    assert_code(result, 0)
    # Wait for the settlement round to end
    economic.wait_settlement_blocknum(node)
    endtime = int(time.time()) + 120
    while int(time.time()) < endtime:
        time.sleep(1)
        current_block = node.eth.blockNumber
        log.info("current block: {}".format(current_block))
        block = current_block % economic.consensus_size
        log.info("block: {}".format(block))
        log.info("Current block height: {}, block of current consensus round: {}".format(current_block, block))
        if block > 19:
            break
    for i in range(4):
        result = check_node_in_list(node.node_id, client.ppos.getValidatorList)
        log.info("Current node in consensus list status：{}".format(result))
        if result:
            # Get current block height
            current_block = node.eth.blockNumber
            log.info("Current block height: {}".format(current_block))
            # Report verifier Duplicate Sign
            result = verification_duplicate_sign(client, 1, 1, report_address, current_block)
            assert_code(result, 0)
            result = check_node_in_list(node.node_id, client.ppos.getValidatorList)
            log.info("Current node in consensus list status：{}".format(result))
            assert result, "ErrMsg:Node current status {}".format(result)
            break
        else:
            # wait consensus block
            economic.wait_consensus_blocknum(node)


@pytest.mark.P2
def test_VP_PVF_006(client_new_node, reset_environment):
    """
    移出PlatON验证人与候选人名单，验证人申请退回质押金
    :param client_new_node:
    :return:
    """
    client = client_new_node
    economic = client.economic
    node = client.node
    # create pledge address
    pledge_address, _ = economic.account.generate_account(node.web3, von_amount(economic.create_staking_limit, 2))
    # create report address
    report_address, _ = economic.account.generate_account(node.web3, node.web3.toWei(1000, 'ether'))
    # create staking
    result = client.staking.create_staking(0, pledge_address, pledge_address)
    assert_code(result, 0)
    # Wait for the settlement round to end
    economic.wait_settlement_blocknum(node)
    for i in range(4):
        result = check_node_in_list(node.node_id, client.ppos.getValidatorList)
        log.info("Current node in consensus list status：{}".format(result))
        if result:
            # Get current block height
            current_block = node.eth.blockNumber
            log.info("Current block height: {}".format(current_block))
            # Report verifier Duplicate Sign
            result = verification_duplicate_sign(client, 1, 1, report_address, current_block)
            assert_code(result, 0)
            # Application for return of pledge
            result = client.staking.withdrew_staking(pledge_address)
            assert_code(result, 301103)
            break
        else:
            # wait consensus block
            economic.wait_consensus_blocknum(node)


@pytest.mark.P2
def test_VP_PVF_007(client_new_node, reset_environment):
    """
    节点被处罚后马上重新质押（双签）
    :param client_new_node:
    :return:
    """
    client = client_new_node
    economic = client.economic
    node = client.node
    # create pledge address
    pledge_address, _ = economic.account.generate_account(node.web3, von_amount(economic.create_staking_limit, 3))
    # create report address
    report_address, _ = economic.account.generate_account(node.web3, node.web3.toWei(1000, 'ether'))
    # create staking
    result = client.staking.create_staking(0, pledge_address, pledge_address)
    assert_code(result, 0)
    # Wait for the settlement round to end
    economic.wait_settlement_blocknum(node)
    for i in range(4):
        result = check_node_in_list(node.node_id, client.ppos.getValidatorList)
        log.info("Current node in consensus list status：{}".format(result))
        if result:
            # Get current block height
            current_block = node.eth.blockNumber
            log.info("Current block height: {}".format(current_block))
            # Report verifier Duplicate Sign
            result = verification_duplicate_sign(client, 1, 1, report_address, current_block)
            assert_code(result, 0)
            # create staking
            result = client.staking.create_staking(0, pledge_address, pledge_address)
            assert_code(result, 301101)
            break
        else:
            # wait consensus block
            economic.wait_consensus_blocknum(node)


@pytest.mark.P2
def test_VP_PVF_008(client_new_node, reset_environment):
    """
    节点被处罚后马上重新增持质押（双签）
    :param client_new_node:
    :return:
    """
    client = client_new_node
    economic = client.economic
    node = client.node
    # create pledge address
    pledge_address, _ = economic.account.generate_account(node.web3, von_amount(economic.create_staking_limit, 3))
    # create report address
    report_address, _ = economic.account.generate_account(node.web3, node.web3.toWei(1000, 'ether'))
    # create staking
    result = client.staking.create_staking(0, pledge_address, pledge_address)
    assert_code(result, 0)
    # Wait for the settlement round to end
    economic.wait_settlement_blocknum(node)
    for i in range(4):
        result = check_node_in_list(node.node_id, client.ppos.getValidatorList)
        log.info("Current node in consensus list status：{}".format(result))
        if result:
            # Get current block height
            current_block = node.eth.blockNumber
            log.info("Current block height: {}".format(current_block))
            # Report verifier Duplicate Sign
            result = verification_duplicate_sign(client, 1, 1, report_address, current_block)
            assert_code(result, 0)
            # increase staking
            result = client.staking.increase_staking(0, pledge_address)
            assert_code(result, 301103)
            break
        else:
            # wait consensus block
            economic.wait_consensus_blocknum(node)


@pytest.mark.P2
def test_VP_PVF_009(client_new_node, reset_environment):
    """
    移出PlatON验证人与候选人名单，委托人可在处罚所在结算周期，申请赎回全部委托金
    :param client_new_node:
    :return:
    """
    client = client_new_node
    economic = client.economic
    node = client.node
    # create pledge address
    pledge_address, _ = economic.account.generate_account(node.web3, von_amount(economic.create_staking_limit, 3))
    # create report address
    report_address, _ = economic.account.generate_account(node.web3, node.web3.toWei(1000, 'ether'))
    # create staking
    result = client.staking.create_staking(0, pledge_address, pledge_address)
    assert_code(result, 0)
    # Additional pledge
    result = client.delegate.delegate(0, report_address)
    assert_code(result, 0)
    # Wait for the settlement round to end
    economic.wait_settlement_blocknum(node)
    for i in range(4):
        result = check_node_in_list(node.node_id, client.ppos.getValidatorList)
        log.info("Current node in consensus list status：{}".format(result))
        if result:
            # Get current block height
            current_block = node.eth.blockNumber
            log.info("Current block height: {}".format(current_block))
            # Report verifier Duplicate Sign
            result = verification_duplicate_sign(client, 1, 1, report_address, current_block)
            assert_code(result, 0)
            time.sleep(3)
            # Access to pledge information
            candidate_info = client.ppos.getCandidateInfo(node.node_id)
            info = candidate_info['Ret']
            staking_blocknum = info['StakingBlockNum']
            # To view the entrusted account balance
            report_balance = node.eth.getBalance(report_address)
            log.info("report address balance: {}".format(report_balance))
            # withdrew delegate
            result = client.delegate.withdrew_delegate(staking_blocknum, report_address)
            assert_code(result, 0)
            # To view the entrusted account balance
            report_balance1 = node.eth.getBalance(report_address)
            log.info("report address balance: {}".format(report_balance1))
            assert report_balance + economic.delegate_limit - report_balance1 < node.web3.toWei(1,
                                                                                                'ether'), "ErrMsg:Ireport balance {}".format(
                report_balance1)
            break
        else:
            # wait consensus block
            economic.wait_consensus_blocknum(node)


def test_VP_PVF_010(client_consensus):
    """

    """
    client = client_consensus
    economic = client.economic
    node = client.node
    client.economic.env.deploy_all()
    # Obtain penalty proportion and income
    pledge_amount1, penalty_ratio, proportion_ratio = penalty_proportion_and_income(client)
    # create report address
    report_address, _ = economic.account.generate_account(node.web3, node.web3.toWei(1000, 'ether'))
    # view report amount
    report_amount1 = node.eth.getBalance(report_address)
    log.info("report account amount:{} ".format(report_amount1))
    # view Incentive pool account
    incentive_pool_account1 = node.eth.getBalance(EconomicConfig.INCENTIVEPOOL_ADDRESS)
    log.info("incentive pool account1 amount:{} ".format(incentive_pool_account1))
    # Wait for the consensus round to end
    economic.wait_consensus_blocknum(node)
    # Get current block height
    current_block = node.eth.blockNumber
    log.info("Current block height: {}".format(current_block))
    result = verification_duplicate_sign(client, 1, 1, report_address, current_block)
    assert_code(result, 0)
    # view Amount of penalty
    proportion_reward, incentive_pool_reward = economic.get_report_reward(pledge_amount1, penalty_ratio,
                                                                          proportion_ratio)
    # view report amount again
    report_amount2 = node.eth.getBalance(report_address)
    log.info("report account amount:{} ".format(report_amount2))
    # view Incentive pool account again
    incentive_pool_account2 = node.eth.getBalance(EconomicConfig.INCENTIVEPOOL_ADDRESS)
    log.info("incentive pool account1 amount:{} ".format(incentive_pool_account2))
    # assert account reward
    assert report_amount1 + proportion_reward - report_amount2 < node.web3.toWei(1,'ether'), "ErrMsg:report amount {}".format(
        report_amount2)
    assert incentive_pool_account2 == incentive_pool_account1 + incentive_pool_reward + (report_amount1 + proportion_reward - report_amount2), "ErrMsg:Incentive pool account {}".format(
        incentive_pool_account2)

    result = node.ppos.getCandidateInfo(node.node_id)
    log.info("Candidate Info:{} ".format(result))
    pledge_amount2 = result['Ret']['Released']
    result = verification_duplicate_sign(client, 1, 1, report_address, current_block + 1)
    assert_code(result, 0)

    # view Amount of penalty
    proportion_reward2, incentive_pool_reward2 = economic.get_report_reward(pledge_amount2, penalty_ratio,
                                                                          proportion_ratio)

    # view report amount again
    report_amount3 = node.eth.getBalance(report_address)
    log.info("report account amount:{} ".format(report_amount3))

    # view Incentive pool account again
    incentive_pool_account3 = node.eth.getBalance(EconomicConfig.INCENTIVEPOOL_ADDRESS)
    log.info("incentive pool account1 amount:{} ".format(incentive_pool_account3))
    # assert account reward
    assert report_amount2 + proportion_reward2 - report_amount3 < node.web3.toWei(1, 'ether'), "ErrMsg:report amount {}".format(
        report_amount2 + proportion_reward2 - report_amount3)
    assert incentive_pool_account3 == incentive_pool_account2 + incentive_pool_reward2 + (report_amount2 + proportion_reward2 - report_amount3), "ErrMsg:Incentive pool account {}".format(
        incentive_pool_account2)