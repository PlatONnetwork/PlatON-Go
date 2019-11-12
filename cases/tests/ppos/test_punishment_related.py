import time
import pytest
import allure
from dacite import from_dict
from common.key import get_pub_key, mock_duplicate_sign, generate_key
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
    # Obtain evidence of violation
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
    assert report_amount1 + proportion_reward - report_amount2 < node.web3.toWei(1,
                                                                                 'ether'), "ErrMsg:report amount {}".format(
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

    @pytest.mark.P1
    def test_VP_PR_001(self, initial_report):
        """
        重复举报-同一举报人
        :param initial_report:
        :return:
        """
        client_con_list_obj, economic, node, report_address, current_block = initial_report
        # duplicate sign
        result = verification_duplicate_sign(client_con_list_obj[0], 1, 1, report_address, current_block)
        assert_code(result, 303001)

    @pytest.mark.P1
    def test_VP_PR_002(self, initial_report):
        """
        重复举报 - 不同举报人
        :param initial_report:
        :return:
        """
        client_con_list_obj, economic, node, report_address, current_block = initial_report
        # create account
        report_address2, _ = economic.account.generate_account(node.web3, node.web3.toWei(1000, 'ether'))
        # duplicate sign
        result = verification_duplicate_sign(client_con_list_obj[0], 1, 1, report_address2, current_block)
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
    evidence_parameter = get_param_by_dict(report_information, 'prepareB', 'blockHash')
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


@pytest.mark.P1
def test_VP_PV_016(client_consensus_obj):
    """
    举报双签-双签证据address不一致
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
    jsondata = update_param_by_dict(report_information, 'prepareA', 'validateNode', 'address',
                                    economic.account.account_with_money['address'])
    log.info("Evidence information: {}".format(jsondata))
    # Report verifier Duplicate Sign
    result = client.duplicatesign.reportDuplicateSign(1, jsondata, report_address)
    assert_code(result, 303000)


@pytest.mark.P1
def test_VP_PV_017(client_con_list_obj):
    """
    举报双签-NodeID不一致举报双签
    :param client_con_list_obj:
    :return:
    """
    client = client_con_list_obj[0]
    economic = client.economic
    node = client.node
    # create report address
    report_address, _ = economic.account.generate_account(node.web3, node.web3.toWei(1000, 'ether'))
    # Obtain information of report evidence
    report_information, current_block = obtaining_evidence_information(economic, node)
    # Modification of evidence
    jsondata = update_param_by_dict(report_information, 'prepareA', 'validateNode', 'nodeId',
                                    client_con_list_obj[1].node.node_id)
    log.info("Evidence information: {}".format(jsondata))
    # Report verifier Duplicate Sign
    result = client.duplicatesign.reportDuplicateSign(1, jsondata, report_address)
    assert_code(result, 303000)


@pytest.mark.P1
def test_VP_PV_018(client_con_list_obj):
    """
    举报双签-blsPubKey不一致举报双签
    :param client_con_list_obj:
    :return:
    """
    client = client_con_list_obj[0]
    economic = client.economic
    node = client.node
    # create report address
    report_address, _ = economic.account.generate_account(node.web3, node.web3.toWei(1000, 'ether'))
    # Obtain information of report evidence
    report_information, current_block = obtaining_evidence_information(economic, node)
    # Modification of evidence
    jsondata = update_param_by_dict(report_information, 'prepareA', 'validateNode', 'blsPubKey',
                                    client_con_list_obj[1].node.blspubkey)
    log.info("Evidence information: {}".format(jsondata))
    # Report verifier Duplicate Sign
    result = client.duplicatesign.reportDuplicateSign(1, jsondata, report_address)
    assert_code(result, 303000)


@pytest.mark.P1
def test_VP_PV_019(client_con_list_obj):
    """
    举报双签-signature一致举报双签
    :param client_con_list_obj:
    :return:
    """
    client = client_con_list_obj[0]
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
def test_VP_PV_020(client_con_list_obj):
    """
    举报双签-伪造合法signature情况下伪造epoch
    :param client_con_list_obj:
    :return:
    """
    client = client_con_list_obj[0]
    economic = client.economic
    node = client.node
    # create report address
    report_address, _ = economic.account.generate_account(node.web3, node.web3.toWei(1000, 'ether'))
    # Obtain information of report evidence
    report_information, current_block = obtaining_evidence_information(economic, node)
    # Obtain evidence of violation
    report_information1 = mock_duplicate_sign(1, node.nodekey, node.blsprikey, current_block, epoch=1)
    log.info("Report information: {}".format(report_information))
    # Modification of evidence
    evidence_parameter = get_param_by_dict(report_information1, 'prepareB', 'signature')
    jsondata = update_param_by_dict(report_information, 'prepareA', 'signature', None, evidence_parameter)
    log.info("Evidence information: {}".format(jsondata))
    # Report verifier Duplicate Sign
    result = client.duplicatesign.reportDuplicateSign(1, jsondata, report_address)
    assert_code(result, 303000)


@pytest.mark.P1
def test_VP_PV_020(client_con_list_obj):
    """
    举报双签-伪造合法signature情况下伪造epoch
    :param client_con_list_obj:
    :return:
    """
    client = client_con_list_obj[0]
    economic = client.economic
    node = client.node
    # create report address
    report_address, _ = economic.account.generate_account(node.web3, node.web3.toWei(1000, 'ether'))
    # Obtain information of report evidence
    report_information, current_block = obtaining_evidence_information(economic, node)
    # Obtain evidence of violation
    report_information1 = mock_duplicate_sign(1, node.nodekey, node.blsprikey, current_block, epoch=1)
    log.info("Report information: {}".format(report_information))
    # Modification of evidence
    evidence_parameter = get_param_by_dict(report_information1, 'prepareB', 'signature')
    jsondata = update_param_by_dict(report_information, 'prepareA', 'signature', None, evidence_parameter)
    log.info("Evidence information: {}".format(jsondata))
    # Report verifier Duplicate Sign
    result = client.duplicatesign.reportDuplicateSign(1, jsondata, report_address)
    assert_code(result, 303000)


@pytest.mark.P1
def test_VP_PV_021(client_con_list_obj):
    """
    举报双签-伪造合法signature情况下伪造viewNumber
    :param client_con_list_obj:
    :return:
    """
    client = client_con_list_obj[0]
    economic = client.economic
    node = client.node
    # create report address
    report_address, _ = economic.account.generate_account(node.web3, node.web3.toWei(1000, 'ether'))
    # Obtain information of report evidence
    report_information, current_block = obtaining_evidence_information(economic, node)
    # Obtain evidence of violation
    report_information1 = mock_duplicate_sign(1, node.nodekey, node.blsprikey, current_block, view_number=1)
    log.info("Report information: {}".format(report_information))
    # Modification of evidence
    evidence_parameter = get_param_by_dict(report_information1, 'prepareB', 'signature')
    jsondata = update_param_by_dict(report_information, 'prepareA', 'signature', None, evidence_parameter)
    log.info("Evidence information: {}".format(jsondata))
    # Report verifier Duplicate Sign
    result = client.duplicatesign.reportDuplicateSign(1, jsondata, report_address)
    assert_code(result, 303000)


@pytest.mark.P1
def test_VP_PV_022(client_con_list_obj):
    """
    举报双签-伪造合法signature情况下伪造blockIndex
    :param client_con_list_obj:
    :return:
    """
    client = client_con_list_obj[0]
    economic = client.economic
    node = client.node
    # create report address
    report_address, _ = economic.account.generate_account(node.web3, node.web3.toWei(1000, 'ether'))
    # Obtain information of report evidence
    report_information, current_block = obtaining_evidence_information(economic, node)
    # Obtain evidence of violation
    report_information1 = mock_duplicate_sign(1, node.nodekey, node.blsprikey, current_block, block_index=1)
    log.info("Report information: {}".format(report_information))
    # Modification of evidence
    evidence_parameter = get_param_by_dict(report_information1, 'prepareB', 'signature')
    jsondata = update_param_by_dict(report_information, 'prepareA', 'signature', None, evidence_parameter)
    log.info("Evidence information: {}".format(jsondata))
    # Report verifier Duplicate Sign
    result = client.duplicatesign.reportDuplicateSign(1, jsondata, report_address)
    assert_code(result, 303000)


@pytest.mark.P1
def test_VP_PV_023(client_con_list_obj):
    """
    举报双签-伪造合法signature情况下伪造index
    :param client_con_list_obj:
    :return:
    """
    client = client_con_list_obj[0]
    economic = client.economic
    node = client.node
    # create report address
    report_address, _ = economic.account.generate_account(node.web3, node.web3.toWei(1000, 'ether'))
    # Obtain information of report evidence
    report_information, current_block = obtaining_evidence_information(economic, node)
    # Obtain evidence of violation
    report_information1 = mock_duplicate_sign(1, node.nodekey, node.blsprikey, current_block, index=1)
    log.info("Report information: {}".format(report_information))
    # Modification of evidence
    evidence_parameter = get_param_by_dict(report_information1, 'prepareB', 'signature')
    jsondata = update_param_by_dict(report_information, 'prepareA', 'signature', None, evidence_parameter)
    log.info("Evidence information: {}".format(jsondata))
    # Report verifier Duplicate Sign
    result = client.duplicatesign.reportDuplicateSign(1, jsondata, report_address)
    assert_code(result, 303000)


@pytest.mark.P1
def test_VP_PV_024(client_con_list_obj):
    """
    举报双签-伪造合法signature情况下伪造blockNumber
    :param client_con_list_obj:
    :return:
    """
    client = client_con_list_obj[0]
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
def test_VP_PV_025(client_consensus_obj):
    """
    举报接口参数测试：举报人账户错误
    :param client_consensus_obj:
    :return:
    """
    client = client_consensus_obj
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
def test_VP_PV_026(client_con_list_obj):
    """
    链存在的id,blskey不匹配
    :param client_con_list_obj:
    :return:
    """
    client = client_con_list_obj[0]
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
    report_information = mock_duplicate_sign(1, node.nodekey, client_con_list_obj[1].node.blsprikey, current_block)
    log.info("Report information: {}".format(report_information))
    # Report verifier Duplicate Sign
    result = client.duplicatesign.reportDuplicateSign(1, report_information, report_address)
    assert_code(result, 303007)


@pytest.mark.P1
def test_VP_PV_027(client_new_node_obj):
    """
    举报候选人
    :param client_new_node_obj:
    :return:
    """
    client = client_new_node_obj
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
def test_VP_PV_028(client_consensus_obj):
    """
    举报有效期之前的双签行为
    :param client_consensus_obj:
    :return:
    """
    client = client_consensus_obj
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
def test_VP_PV_028(client_consensus_obj):
    """
    举报有效期之前的双签行为
    :param client_consensus_obj:
    :return:
    """
    client = client_consensus_obj
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


@pytest.mark.P1
def test_VP_PV_031(client_consensus_obj):
    """
    举报的gas费不足
    :param client_consensus_obj:
    :return:
    """
    client = client_consensus_obj
    economic = client.economic
    node = client.node
    # create report address
    report_address, _ = economic.account.generate_account(node.web3, 0)
    # Obtain information of report evidence
    report_information, current_block = obtaining_evidence_information(economic, node)
    try:
        # Report verifier Duplicate Sign
        result = client.duplicatesign.reportDuplicateSign(1, report_information, report_address)
        assert_code(result, 0)
    except Exception as e:
        log.info("Use case success, exception information：{} ".format(str(e)))


@pytest.mark.P1
def test_VP_PR_003(client_new_node_obj):
    """
    举报被处罚退出状态中的验证人
    :param client_new_node_obj:
    :return:
    """
    client = client_new_node_obj
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
            # Get current block height
            current_block = node.eth.blockNumber
            log.info("Current block height: {}".format(current_block))
            # Report verifier Duplicate Sign
            result = verification_duplicate_sign(client, 1, 1, report_address, current_block)
            assert_code(result, 0)
            # Obtain penalty proportion and income
            pledge_amount1, penalty_ratio, proportion_ratio = penalty_proportion_and_income(client)
            # view Amount of penalty
            proportion_reward, incentive_pool_reward = economic.get_report_reward(pledge_amount1, penalty_ratio,
                                                                                  proportion_ratio)
            # view Pledge node information
            candidate_info = client.ppos.getCandidateInfo(node.node_id)
            log.info("Pledge node information: {}".format(candidate_info))
            info = candidate_info['Ret']
            assert info['Released'] == pledge_amount1 - (
                        proportion_reward + incentive_pool_reward), "ErrMsg:Pledge amount {}".format(
                info['Released'])
        else:
            # wait consensus block
            economic.wait_consensus_blocknum(node)


@pytest.mark.P1
def test_VP_PR_004(client_new_node_obj):
    """
    举报人和被举报人为同一个人
    :param client_new_node_obj:
    :return:
    """
    client = client_new_node_obj
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
            assert_code(result, 0)
            # Obtain penalty proportion and income
            pledge_amount1, penalty_ratio, proportion_ratio = penalty_proportion_and_income(client)
            # view Amount of penalty
            proportion_reward, incentive_pool_reward = economic.get_report_reward(pledge_amount1, penalty_ratio,
                                                                                  proportion_ratio)
            log.info("Pledge node amount: {}".format(pledge_amount1))
            log.info("proportion_reward + incentive_pool_reward: {}".format(proportion_reward + incentive_pool_reward))
            # view Pledge node information
            candidate_info = client.ppos.getCandidateInfo(node.node_id)
            log.info("Pledge node information: {}".format(candidate_info))
            info = candidate_info['Ret']
            assert info['Released'] == pledge_amount1 - (
                        proportion_reward + incentive_pool_reward), "ErrMsg:Pledge amount {}".format(
                info['Released'])
        else:
            # wait consensus block
            economic.wait_consensus_blocknum(node)
