import time
import pytest
import allure
from dacite import from_dict
from common.key import get_pub_key, mock_duplicate_sign
from common.log import log
from client_sdk_python import Web3
from decimal import Decimal
from tests.lib import EconomicConfig, Genesis, StakingConfig, Staking, check_node_in_list, assert_code, von_amount, \
    get_governable_parameter_value, Client, update_param_by_dict, get_param_by_dict


def penalty_proportion_and_income(client_obj):
    # view Pledge amount
    candidate_info1 = client_obj.ppos.getCandidateInfo(client_obj.node.node_id)
    pledge_amount1 = candidate_info1['Ret']['Released']
    # view Parameter value before treatment
    penalty_ratio = get_governable_parameter_value(client_obj, 'SlashFractionDuplicateSign')
    proportion_ratio = get_governable_parameter_value(client_obj, 'DuplicateSignReportReward')
    return pledge_amount1, int(penalty_ratio), int(proportion_ratio)


def verification_duplicate_sign(client_obj, evidence_type, reporting_type, report_address, report_block):
    if report_block < 41:
        report_block = 41
    # Report1 prepareblock signature
    report_information = mock_duplicate_sign(evidence_type, client_obj.node.nodekey,
                                             client_obj.node.blsprikey,
                                             report_block)
    log.info("Report information: {}".format(report_information))
    result = client_obj.duplicatesign.reportDuplicateSign(reporting_type, report_information, report_address)
    return result


@pytest.mark.P0
@pytest.mark.parametrize('repor_type', [1, 2, 3])
def test_VP_PV_001_to_003(client_consensus_obj, repor_type, reset_environment):
    """
    举报验证人区块双签:VP_PV_001 prepareBlock类型
                    VP_PV_002 举报验证人区块双签prepareVote类型
                    VP_PV_003 举报验证人区块双签viewChange类型
    :param client_consensus_obj:
    :param repor_type:
    :param reset_environment:
    :return:
    """
    client = client_consensus_obj
    economic = client.economic
    node = client.node
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
    assert incentive_pool_account2 == incentive_pool_account1 + incentive_pool_reward + (
            report_amount1 + proportion_reward - report_amount2), "ErrMsg:Incentive pool account {}".format(
        incentive_pool_account2)


@pytest.fixture(scope='class')
def initial_report(global_test_env):
    """
    Report a non consensus node prepareBlock
    :param global_test_env:
    :return:
    """
    cfg = StakingConfig("11111", "faker", "www.baidu.com", "how much")
    client_con_list_obj = []
    consensus_node_obj_list = global_test_env.consensus_node_list
    for node_obj in consensus_node_obj_list:
        client_con_list_obj.append(Client(global_test_env, node_obj, cfg))
    client = client_con_list_obj[0]
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
    yield client_con_list_obj, economic, node, report_address, current_block
    log.info("case execution completed")
    global_test_env.deploy_all()
    time.sleep(3)
    # # create pledge address
    # pledge_address, _ = economic.account.generate_account(node.web3, von_amount(economic.create_staking_limit, 2))
    # # create staking
    # result = client.staking.create_staking(0, pledge_address, pledge_address)
    # assert_code(result, 0)
    # # Wait for the settlement round to end
    # economic.wait_settlement_blocknum(node)
    # for i in range(4):
    #     result = check_node_in_list(node.node_id, client.ppos.getValidatorList)
    #     log.info("Current node in consensus list status：{}".format(result))
    #     if result:
    #         # Get current block height
    #         current_block = node.eth.blockNumber
    #         log.info("Current block height: {}".format(current_block))
    #         result = verification_duplicate_sign(client, 1, 1, report_address, current_block)
    #         assert_code(result, 0)
    #         yield client, economic, node, report_address, current_block
    #         log.info("case execution completed")
    #         global_test_env.deploy_all()
    #         time.sleep(3)
    #     else:
    #         # wait consensus block
    #         economic.wait_consensus_blocknum(node)


class TestMultipleReports:
    @pytest.mark.P1
    def test_VP_PV_004(self, initial_report):
        """
        举报双签-同一验证人同一块高不同类型
        :param initial_report:
        :return:
        """
        client_con_list_obj, economic, node, report_address, current_block = initial_report
        # duplicate sign
        result = verification_duplicate_sign(client_con_list_obj[0], 2, 2, report_address, current_block)
        assert_code(result, 0)

    @pytest.mark.P1
    def test_VP_PV_005(self, initial_report):
        """
        举报双签-同一验证人不同块高同一类型
        :param initial_report:
        :return:
        """
        client_con_list_obj, economic, node, report_address, current_block = initial_report
        # duplicate sign
        result = verification_duplicate_sign(client_con_list_obj[0], 1, 1, report_address, current_block - 1)
        assert_code(result, 0)

    @pytest.mark.P1
    def test_VP_PV_006(self, initial_report):
        """
        举报双签-同一验证人不同块高不同类型
        :param initial_report:
        :return:
        """
        client_con_list_obj, economic, node, report_address, current_block = initial_report
        # duplicate sign
        result = verification_duplicate_sign(client_con_list_obj[0], 2, 2, report_address, current_block - 1)
        assert_code(result, 0)

    @pytest.mark.P1
    def test_VP_PV_007(self, initial_report):
        """
        举报双签-不同验证人同一块高同一类型
        :param initial_report:
        :return:
        """
        client_con_list_obj, economic, node, report_address, current_block = initial_report
        # create account
        report_address2, _ = economic.account.generate_account(node.web3, node.web3.toWei(1000, 'ether'))
        # duplicate sign
        result = verification_duplicate_sign(client_con_list_obj[1], 1, 1, report_address, current_block)
        assert_code(result, 0)

    @pytest.mark.P1
    def test_VP_PV_008(self, initial_report):
        """
        举报双签-不同验证人同一块高不同类型
        :param initial_report:
        :return:
        """
        client_con_list_obj, economic, node, report_address, current_block = initial_report
        # create account
        report_address2, _ = economic.account.generate_account(node.web3, node.web3.toWei(1000, 'ether'))
        # duplicate sign
        result = verification_duplicate_sign(client_con_list_obj[1], 2, 2, report_address, current_block)
        assert_code(result, 0)

    @pytest.mark.P1
    def test_VP_PV_009(self, initial_report):
        """
        举报双签-不同验证人不同块高不同类型
        :param initial_report:
        :return:
        """
        client_con_list_obj, economic, node, report_address, current_block = initial_report
        # create account
        report_address2, _ = economic.account.generate_account(node.web3, node.web3.toWei(1000, 'ether'))
        # duplicate sign
        result = verification_duplicate_sign(client_con_list_obj[1], 2, 2, report_address, current_block - 1)
        assert_code(result, 0)


def obtaining_evidence_information(economic, node):
    # Wait for the consensus round to end
    economic.wait_consensus_blocknum(node, 1)
    # Get current block height
    current_block = node.eth.blockNumber
    log.info("Current block height: {}".format(current_block))
    if current_block < 41:
        current_block = 41
    # Report1 prepareblock signature
    report_information = mock_duplicate_sign(1, node.nodekey, node.blsprikey, current_block)
    log.info("Report information: {}".format(report_information))
    return report_information, current_block


@pytest.mark.P1
def test_VP_PV_010(client_consensus_obj):
    """
    举报双签-双签证据epoch不一致
    :param client_consensus_obj:
    :return:
    """
    client = client_consensus_obj
    economic = client.economic
    node = client.node
    # create report address
    report_address, _ = economic.account.generate_account(node.web3, node.web3.toWei(1000, 'ether'))
    # Obtain information of report evidence
    report_information, current_block = obtaining_evidence_information(economic, node)
    # Modification of evidence
    jsondata = update_param_by_dict(report_information, 'prepareA', 'epoch', None, 1)
    log.info("Evidence information: {}".format(jsondata))
    # Report verifier Duplicate Sign
    result = client.duplicatesign.reportDuplicateSign(1, jsondata, report_address)
    assert_code(result, 303000)


@pytest.mark.P1
def test_VP_PV_011(client_consensus_obj):
    """
    举报双签-双签证据view_number不一致
    :param client_consensus_obj:
    :return:
    """
    client = client_consensus_obj
    economic = client.economic
    node = client.node
    # create report address
    report_address, _ = economic.account.generate_account(node.web3, node.web3.toWei(1000, 'ether'))
    # Obtain information of report evidence
    report_information, current_block = obtaining_evidence_information(economic, node)
    # Modification of evidence
    jsondata = update_param_by_dict(report_information, 'prepareA', 'viewNumber', None, 1)
    log.info("Evidence information: {}".format(jsondata))
    # Report verifier Duplicate Sign
    result = client.duplicatesign.reportDuplicateSign(1, jsondata, report_address)
    assert_code(result, 303000)


@pytest.mark.P1
def test_VP_PV_012(client_consensus_obj):
    """
    举报双签-双签证据block_number不一致
    :param client_consensus_obj:
    :return:
    """
    client = client_consensus_obj
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
def test_VP_PV_013(client_consensus_obj):
    """
    举报双签-双签证据block_hash一致
    :param client_consensus_obj:
    :return:
    """
    client = client_consensus_obj
    economic = client.economic
    node = client.node
    # create report address
    report_address, _ = economic.account.generate_account(node.web3, node.web3.toWei(1000, 'ether'))
    # Obtain information of report evidence
    report_information, current_block = obtaining_evidence_information(economic, node)
    # Modification of evidence
    evidence_parameter = get_param_by_dict(report_information, 'prepareB', 'blockHash', None)
    jsondata = update_param_by_dict(report_information, 'prepareA', 'blockHash', None, evidence_parameter)
    log.info("Evidence information: {}".format(jsondata))
    # Report verifier Duplicate Sign
    result = client.duplicatesign.reportDuplicateSign(1, jsondata, report_address)
    assert_code(result, 303000)


@pytest.mark.P1
def test_VP_PV_014(client_consensus_obj):
    """
    举报双签-双签证据block_index不一致
    :param client_consensus_obj:
    :return:
    """
    client = client_consensus_obj
    economic = client.economic
    node = client.node
    # create report address
    report_address, _ = economic.account.generate_account(node.web3, node.web3.toWei(1000, 'ether'))
    # Obtain information of report evidence
    report_information, current_block = obtaining_evidence_information(economic, node)
    # Modification of evidence
    jsondata = update_param_by_dict(report_information, 'prepareA', 'blockIndex', None, 1)
    log.info("Evidence information: {}".format(jsondata))
    # Report verifier Duplicate Sign
    result = client.duplicatesign.reportDuplicateSign(1, jsondata, report_address)
    assert_code(result, 303000)


@pytest.mark.P1
def test_VP_PV_015(client_consensus_obj):
    """
    举报双签-双签证据validate_node-index不一致
    :param client_consensus_obj:
    :return:
    """
    client = client_consensus_obj
    economic = client.economic
    node = client.node
    # create report address
    report_address, _ = economic.account.generate_account(node.web3, node.web3.toWei(1000, 'ether'))
    # Obtain information of report evidence
    report_information, current_block = obtaining_evidence_information(economic, node)
    # Modification of evidence
    jsondata = update_param_by_dict(report_information, 'prepareA', 'validateNode', 'index', 1)
    log.info("Evidence information: {}".format(jsondata))
    # Report verifier Duplicate Sign
    result = client.duplicatesign.reportDuplicateSign(1, jsondata, report_address)
    assert_code(result, 303000)
