import time
import pytest
import allure
from dacite import from_dict
from common.key import get_pub_key, mock_duplicate_sign
from common.log import log
from client_sdk_python import Web3
from decimal import Decimal
from tests.lib import EconomicConfig, Genesis, StakingConfig, Staking, check_node_in_list, assert_code, von_amount, \
    get_governable_parameter_value


def penalty_proportion_and_income(client_obj):
    # view Pledge amount
    candidate_info1 = client_obj.ppos.getCandidateInfo(client_obj.node.node_id)
    pledge_amount1 = candidate_info1['Ret']['Released']
    # view Parameter value before treatment
    penalty_ratio = get_governable_parameter_value(client_obj, 'SlashFractionDuplicateSign')
    proportion_ratio = get_governable_parameter_value(client_obj, 'DuplicateSignReportReward')
    return pledge_amount1, penalty_ratio, proportion_ratio


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
def test_VP_PV_001(client_consensus_obj, reset_environment):
    """
    举报验证人区块双签prepareBlock类型
    :param client_consensus_obj:
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
    result = verification_duplicate_sign(client, 1, 1, report_address, current_block)
    assert_code(result, 0)
    # view Amount of penalty
    proportion_reward, incentive_pool_reward = economic.get_report_reward(pledge_amount1, penalty_ratio, proportion_ratio)
    # view report amount again
    report_amount2 = node.eth.getBalance(report_address)
    log.info("report account amount:{} ".format(report_amount2))
    # view Incentive pool account again
    incentive_pool_account2 = node.eth.getBalance(EconomicConfig.INCENTIVEPOOL_ADDRESS)
    log.info("incentive pool account1 amount:{} ".format(incentive_pool_account2))
    # assert account reward
    assert report_amount1 + proportion_reward - report_amount2 < node.web3.toWei(1, 'ether'), "ErrMsg:report amount {}".format(
        report_amount2)
    assert incentive_pool_account2 == incentive_pool_account1 + incentive_pool_reward, "ErrMsg:incentive pool account amount {}".format(
        incentive_pool_account2)






