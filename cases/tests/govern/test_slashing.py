import pytest
from tests.govern.conftest import proposal_vote, preactive_proposal_pips
from tests.lib.config import PipConfig
from loguru import logger
from tests.lib import check_node_in_list, assert_code, von_amount, get_governable_parameter_value, get_pledge_list
from tests.ppos.test_general_punishment import verify_low_block_rate_penalty, get_out_block_penalty_parameters
from tests.lib.client import get_clients_by_nodeid


@pytest.fixture()
def config():
    return PipConfig()


def create_0mb_slash(slash_node, check_node):
    """
    构造低出块处罚
    :param client_new_node_obj_list_reset:
    :return:
    """
    economic = slash_node.economic
    node = slash_node.node
    logger.info("Start creating a pledge account address")
    # address, _ = economic.account.generate_account(node.web3, von_amount(economic.create_staking_limit, 2))
    # logger.info("Start applying for a pledge node")
    # result = slash_node.staking.create_staking(0, address, address)
    # assert_code(result, 0)
    # logger.info("Pledge completed, waiting for the end of the current billing cycle")
    # economic.wait_settlement_blocknum(node)
    logger.info("Get the current pledge node amount and the low block rate penalty block number and the block reward")
    pledge_amount1, block_reward, slash_blocks = get_out_block_penalty_parameters(slash_node, node, 'Released')
    logger.info(
        "Current node deposit amount: {} Current year block reward: {} Current low block rate penalty block: {}".format(
            pledge_amount1, block_reward, slash_blocks))
    logger.info("Current block height: {}".format(slash_node.node.eth.blockNumber))
    logger.info("Start verification penalty amount")
    verify_low_block_rate_penalty(slash_node, check_node, block_reward, slash_blocks, pledge_amount1, 'Released')
    logger.info("Check amount completed")
    result = check_node.ppos.getCandidateInfo(slash_node.node.node_id)
    logger.info("Candidate Info：{}".format(result))
    result = check_node_in_list(slash_node.node.node_id, check_node.ppos.getCandidateList)
    assert result is False, "error: Node not kicked out CandidateList"
    result = check_node_in_list(slash_node.node.node_id, check_node.ppos.getVerifierList)
    assert result is False, "error: Node not kicked out VerifierList"
    result = check_node_in_list(slash_node.node.node_id, check_node.ppos.getValidatorList)
    assert result is False, "error: Node not kicked out ValidatorList"


class TestSlashing:

    def test_0mb_freeze_after_version_vote(self, proposal_pips, submit_version):
        """
        @describe: 版本升级提案投票后，预生效期节点零出块冻结，投票有效，提案可正常生效
        @step:
        - 1. 提交版本提案并投票通过，使投票预生效
        - 2. 停止一个节点，等待节点被零出块处罚
        - 3. 检查提案和投票信息是否正确
        @expect:
        - 1. 节点被处罚后，投票有效，累积验证人含被处罚节点
        - 2. 节点被处罚后，提案可正常生效
        - 3. 所有相关查询接口，返回提案信息正确
        """
        pip = submit_version
        # proposal_vote(pip, config.vote_option_yeas, config.version_proposal)
        pips = preactive_proposal_pips
        create_0mb_slash(clients_consensus[0], clients_consensus[1])
        pip = pips[1]
        proposal_id = pip.get_effect_proposal_info_of_preactive()
        all_verifiers = pip.get_accu_verifiers_of_proposal(proposal_id)
        all_yeas = pip.get_yeas_of_proposal(proposal_id)
        assert all_verifiers == len(pips)
        assert all_yeas == len(pips)

    def test_0mb_freeze_after_param_vote(self, ):
        """
        @describe: 参数提案投票后，预生效期节点零出块冻结，投票有效，提案可正常生效
        @step:
        - 1. 提交版本提案并投票通过，使投票预生效
        - 2. 停止一个节点，等待节点被零出块处罚
        - 3. 检查提案和投票信息是否正确
        @expect:
        - 1. 节点被处罚后，投票有效，累积验证人含被处罚节点
        - 2. 节点被处罚后，提案可正常生效
        - 3. 所有相关查询接口，返回提案信息正确
        """


    def test_0mb_freeze_after_txt_vote(self):
        """
        文本投票后零出块冻结
        """

    def test_0mb_freeze_after_cancel_vote(self):
        """
        取消投票后零出块冻结
        """

    def test_version_vote_at_0mb_freezing(self):
        """
        零出块冻结期版本投票
        """

    def test_param_vote_at_0mb_freezing(self):
        """
        零出块冻结期参数投票
        """

    def test_txt_vote_at_0mb_freezing(self):
        """
        零出块冻结期文本投票
        """

    def test_cancel_vote_at_0mb_freezing(self):
        """
        零出块冻结期取消投票
        """

    def test_vote_after_0mb_unfreeze(self):
        """
        零出块解冻后投票
        """

    def test_send_proposal_at_0mb_freezing(self):
        """
        零出块冻结期提案
        """

    def test_send_declare_at_0mb_freezing(self):
        """
        零出块冻结期声明
        """

    @pytest.mark.parametrize('', [])
    def test_modify_0mb_frzzez_time_param(self):
        """
        修改零出块冻结时长参数-成功
        """

    @pytest.mark.parametrize('', [])
    def test_modify_0mb_frzzez_time_param_fail(self):
        """
        修改零出块冻结时长参数-失败
        """


        """
        解质押和处罚并行
        """

        """
        测创世文件
        """