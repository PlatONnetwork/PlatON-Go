import time
from typing import List
import pytest
from common.log import log
from tests.lib import check_node_in_list, upload_platon, wait_block_number
from tests.ppos.test_general_punishment import verify_low_block_rate_penalty, get_out_block_penalty_parameters
from tests.lib.client import Client
from tests.lib.config import PipConfig


@pytest.fixture()
def verifiers(clients_consensus):
    return clients_consensus


def get_pips(clients: List[Client]):
    return [c.pip for c in clients]


def version_proposal(pip, to_version, voting_rounds):
    result = pip.submitVersion(pip.node.node_id, str(time.time()), to_version, voting_rounds,
                               pip.node.staking_address,
                               transaction_cfg=pip.cfg.transaction_cfg)
    log.info('submit version proposal result : {}'.format(result))
    return get_proposal_result(pip, pip.cfg.version_proposal, result)


def param_proposal(pip, module, name, value):
    result = pip.submitParam(pip.node.node_id, str(time.time()), module, name, value,
                             pip.node.staking_address, transaction_cfg=pip.cfg.transaction_cfg)
    log.info('submit param proposal result : {}'.format(result))
    return get_proposal_result(pip, pip.cfg.param_proposal, result)


def text_proposal(pip):
    result = pip.submitText(pip.node.node_id, str(time.time()), pip.node.staking_address,
                            transaction_cfg=pip.cfg.transaction_cfg)
    log.info('submit text proposal result:'.format(result))
    return get_proposal_result(pip, pip.cfg.text_proposal, result)


def cancel_proposal(pip, pip_id, voting_rounds):
    result = pip.submitCancel(pip.node.node_id, str(time.time()), 2, pip_id,
                              pip.node.staking_address, transaction_cfg=pip.cfg.transaction_cfg)
    log.info('submit cancel proposal result : {}'.format(result))
    return get_proposal_result(pip, pip.cfg.cancel_proposal, result)


def get_proposal_result(pip, proposal_type, code):
    if code == 0:
        pip_info = pip.get_effect_proposal_info_of_vote(proposal_type)
        return pip_info['ProposalID']
    return code


def vote(pip, pip_id, vote_option=PipConfig.vote_option_yeas):
    result = pip.vote(pip.node.node_id, pip_id, vote_option,
                      pip.node.staking_address, transaction_cfg=pip.cfg.transaction_cfg)
    log.info(f'Node {pip.node.node_id} vote param proposal result {result}')
    return result


def votes(pip_id, pips, vote_options):
    assert len(pips) == len(vote_options)
    for pip, vote_option in zip(pips, vote_options):
        assert vote(pip, pip_id, vote_option) == 0
    return True


def version_declare(pip):
    version = pip.node.program_version
    version_sign = pip.node.program_version_sign
    result = pip.declareVersion(pip.node.node_id, pip.node.staking_address, version, version_sign,
                                transaction_cfg=pip.cfg.transaction_cfg)
    log.info(f'Node {pip.node.node_id} declare version result {result}')
    return result


def wait_proposal_Active(pip, pip_id):
    result = pip.pip.getProposal(pip_id)
    # endBlock =
    # wait_block_number()
    return result


def make_0mb_slash(slash_client, check_client):
    """
    构造低出块处罚
    :param client_new_node_obj_list_reset:
    :return:
    """
    slash_node = slash_client.node
    # check_client = get_client_by_nodeid(check_node.node_id)
    log.info("Get the current pledge node amount and the low block rate penalty block number and the block reward")
    pledge_amount1, block_reward, slash_blocks = get_out_block_penalty_parameters(slash_client, slash_node, 'Released')
    log.info(
        "Current node deposit amount: {} Current year block reward: {} Current low block rate penalty block: {}".format(
            pledge_amount1, block_reward, slash_blocks))
    log.info("Current block height: {}".format(slash_node.eth.blockNumber))
    log.info("Start verification penalty amount")
    verify_low_block_rate_penalty(slash_client, check_client, block_reward, slash_blocks, pledge_amount1, 'Released')
    log.info("Check amount completed")
    result = check_client.ppos.getCandidateInfo(slash_node.node_id)
    log.info("Candidate Info：{}".format(result))
    result = check_node_in_list(slash_node.node_id, check_client.ppos.getCandidateList)
    assert result is False, "error: Node not kicked out CandidateList"
    result = check_node_in_list(slash_node.node_id, check_client.ppos.getVerifierList)
    assert result is False, "error: Node not kicked out VerifierList"
    result = check_node_in_list(slash_node.node_id, check_client.ppos.getValidatorList)
    assert result is False, "error: Node not kicked out ValidatorList"


class TestSlashing:

    def test_0mb_freeze_after_version_vote(self, verifiers):
        """
        @describe: 版本升级提案投票后，节点零出块冻结，投票有效，提案可正常生效
        @step:
        - 1. 提交版本提案并进行投票
        - 2. 停止节点，等待节点被零出块处罚
        - 3. 检查提案和投票信息是否正确
        @expect:
        - 1. 节点被处罚后，投票有效，累积验证人含被处罚节点
        - 2. 节点被处罚后，提案可正常生效
        - 3. 所有相关查询接口，返回提案信息正确
        """
        # step1：提交版本提案并进行投票
        pips = get_pips(verifiers)
        pip = pips[0]
        pip_id = version_proposal(pip, pip.cfg.version5, 5)
        upload_platon(pip.node, pip.cfg.PLATON_NEW_BIN)
        vote(pip, pip_id)
        # step2：停止节点，等待节点被零出块处罚
        make_0mb_slash(verifiers[0], verifiers[1])
        # step3：检查提案和投票信息是否正确
        pip = pips[1]
        all_verifiers = pip.get_accu_verifiers_of_proposal(pip_id)
        all_yeas = pip.get_yeas_of_proposal(pip_id)
        assert all_verifiers == 1
        assert all_yeas == 1

    def test_0mb_freeze_after_param_vote(self, verifiers):
        """
        @describe: 参数提案投票后，节点零出块冻结，投票有效，提案可正常生效
        @step:
        - 1. 提交参数提案并进行投票
        - 2. 停止节点，等待节点被零出块处罚
        - 3. 检查提案和投票信息是否正确
        @expect:
        - 1. 节点被处罚后，投票有效，累积验证人含被处罚节点
        - 2. 节点被处罚后，提案可正常生效
        - 3. 所有相关查询接口，返回提案信息正确
        """
        # step1：提交参数提案并进行投票
        pips = get_pips(verifiers)
        pip = pips[0]
        pip_id = param_proposal(pip, 'slashing', 'zeroProduceFreezeDuration', 5)
        vote(pip, pip_id)
        # step2：停止节点，等待节点被零出块处罚
        make_0mb_slash(verifiers[0], verifiers[1])
        # step3：检查提案和投票信息是否正确
        pip = pips[1]
        all_verifiers = pip.get_accu_verifiers_of_proposal(pip_id)
        all_yeas = pip.get_yeas_of_proposal(pip_id)
        assert all_verifiers == 1
        assert all_yeas == 1

    def test_0mb_freeze_after_text_vote(self, verifiers):
        """
        @describe: 文本提案投票后，节点零出块冻结，投票有效，提案可正常生效
        @step:
        - 1. 提交文本提案并进行投票
        - 2. 停止节点，等待节点被零出块处罚
        - 3. 检查提案和投票信息是否正确
        @expect:
        - 1. 节点被处罚后，投票有效，累积验证人含被处罚节点
        - 2. 节点被处罚后，提案可正常生效
        - 3. 所有相关查询接口，返回提案信息正确
        """
        # step1：提交文本提案并进行投票
        pips = get_pips(verifiers)
        pip = pips[0]
        pip_id = text_proposal(pip)
        vote(pip, pip_id)
        # step2：停止节点，等待节点被零出块处罚
        make_0mb_slash(verifiers[0], verifiers[1])
        # step3：检查提案和投票信息是否正确
        pip = pips[1]
        all_verifiers = pip.get_accu_verifiers_of_proposal(pip_id)
        all_yeas = pip.get_yeas_of_proposal(pip_id)
        assert all_verifiers == 1
        assert all_yeas == 1

    def test_0mb_freeze_after_cancel_vote(self, verifiers):
        """
        @describe: 取消提案投票后，节点零出块冻结，投票有效，提案可正常生效
        @step:
        - 1. 提交取消提案并进行投票
        - 2. 停止节点，等待节点被零出块处罚
        - 3. 检查提案和投票信息是否正确
        @expect:
        - 1. 节点被处罚后，投票有效，累积验证人含被处罚节点
        - 2. 节点被处罚后，提案可正常生效
        - 3. 所有相关查询接口，返回提案信息正确
        """
        # step1：提交版本提案并进行投票
        pips = get_pips(verifiers)
        pip = pips[0]
        pip_id = version_proposal(pip, pip.cfg.version5, 5)
        upload_platon(pip.node, pip.cfg.PLATON_NEW_BIN)
        vote(pip, pip_id)
        pip_id = cancel_proposal(pip, pip_id, 2)
        # step2：停止节点，等待节点被零出块处罚
        make_0mb_slash(verifiers[0], verifiers[1])
        # step3：检查提案和投票信息是否正确
        pip = pips[1]
        all_verifiers = pip.get_accu_verifiers_of_proposal(pip_id)
        all_yeas = pip.get_yeas_of_proposal(pip_id)
        assert all_verifiers == 1
        assert all_yeas == 1

    def test_submit_proposal_at_0mb_freezing(self, verifiers):
        """
        @describe: 节点零出块冻结期内，进行提案，提案失败
        @step:
        - 1. 节点零出块被冻结处罚
        - 2. 冻结期内发送各种提案，提案失败
        @expect:
        - 1. 节点被处罚冻结期内，提案失败
        - 2. 查询未新增提案信息
        """
        # step1：停止节点，等待节点被零出块处罚
        pips = get_pips(verifiers)
        pip = pips[0]
        make_0mb_slash(verifiers[0], verifiers[1])
        # step2：提交各类提案，提案失败
        assert version_proposal(pip, pip.cfg.version5, 5) == 0
        assert param_proposal(pip, 'slashing', 'zeroProduceFreezeDuration', 5) == 0
        assert text_proposal(pip) == 0
        assert cancel_proposal(pip, 'test', 2) == 0

    def test_version_vote_at_0mb_freezing(self, verifiers):
        """
        @describe: 节点零出块冻结期内，进行版本升级提案投票，投票失败
        @step:
        - 1. 节点零出块被冻结处罚
        - 2. 冻结期内进行版本升级提案投票，投票失败
        @expect:
        - 1. 节点被处罚冻结期内，投票失败
        - 2. 提案投票信息查询正确
        - 3. 可投票验证人统计中，不包含被处罚节点
        """
        # step1：提交版本提案
        pips = get_pips(verifiers)
        pip = pips[0]
        pip_id = version_proposal(pip, pip.cfg.version5, 5)
        # step2：停止节点，等待节点被零出块处罚
        pip = pips[1]
        upload_platon(pip.node, pip.cfg.PLATON_NEW_BIN)
        make_0mb_slash(verifiers[1], verifiers[0])
        assert vote(pip, pip_id) == 0
        # step3：检查提案和投票信息是否正确
        # pip = pips[1]
        # all_verifiers = pip.get_accu_verifiers_of_proposal(pip_id)
        # all_yeas = pip.get_yeas_of_proposal(pip_id)
        # assert all_verifiers == 1
        # assert all_yeas == 1

    def test_param_vote_at_0mb_freezing(self, verifiers):
        """
        @describe: 节点零出块冻结期内，进行参数提案投票，投票失败
        @step:
        - 1. 节点零出块被冻结处罚
        - 2. 冻结期内进行参数提案投票，投票失败
        @expect:
        - 1. 节点被处罚冻结期内，投票失败
        - 2. 提案投票信息查询正确
        - 3. 可投票验证人统计中，不包含被处罚节点
        """
        # step1：提交参数提案
        pips = get_pips(verifiers)
        pip = pips[0]
        pip_id = param_proposal(pip, 'slashing', 'zeroProduceFreezeDuration', 5)
        # step2：停止节点，等待节点被零出块处罚
        pip = pips[1]
        upload_platon(pip.node, pip.cfg.PLATON_NEW_BIN)
        make_0mb_slash(verifiers[1], verifiers[0])
        assert vote(pip, pip_id) == 0

    def test_txt_vote_at_0mb_freezing(self, verifiers):
        """
        @describe: 节点零出块冻结期内，进行文本提案投票，投票失败
        @step:
        - 1. 节点零出块被冻结处罚
        - 2. 冻结期内进行文本提案投票，投票失败
        @expect:
        - 1. 节点被处罚冻结期内，投票失败
        - 2. 提案投票信息查询正确
        - 3. 可投票验证人统计中，不包含被处罚节点
        """
        # step1：提交参数提案
        pips = get_pips(verifiers)
        pip = pips[0]
        pip_id = text_proposal(pip)
        # step2：停止节点，等待节点被零出块处罚
        pip = pips[1]
        upload_platon(pip.node, pip.cfg.PLATON_NEW_BIN)
        make_0mb_slash(verifiers[1], verifiers[0])
        assert vote(pip, pip_id) == 0

    def test_cancel_vote_at_0mb_freezing(self, verifiers):
        """
        @describe: 节点零出块冻结期内，进行取消提案投票，投票失败
        @step:
        - 1. 节点零出块被冻结处罚
        - 2. 冻结期内进行文本提案投票，投票失败
        @expect:
        - 1. 节点被处罚冻结期内，投票失败
        - 2. 提案投票信息查询正确
        - 3. 可投票验证人统计中，不包含被处罚节点
        """
        # step1：提交取消提案
        pips = get_pips(verifiers)
        pip = pips[0]
        pip_id = version_proposal(pip, pip.cfg.version5, 5)
        pip_id = cancel_proposal(pip, pip_id, 2)
        # step2：停止节点，等待节点被零出块处罚
        pip = pips[1]
        upload_platon(pip.node, pip.cfg.PLATON_NEW_BIN)
        make_0mb_slash(verifiers[1], verifiers[0])
        assert vote(pip, pip_id) == 0

    def test_sbumit_declare_at_0mb_freezing(self, verifiers):
        """
        @describe: 节点零出块冻结期内，进行版本声明
        @step:
        - 1. 节点零出块被冻结处罚
        - 2. 冻结期内进行版本声明
        @expect:
        - 1. 节点被处罚冻结期内，可以进行版本声明
        - 2. 冻结期内，已发送版本声明，也不会被选举
        """
        # step1：提交版本升级提案
        pips = get_pips(verifiers)
        pip = pips[0]
        pip_id = version_proposal(pip, pip.cfg.version5, 5)
        # setp2：使用其他节点，对提案进行投票，使提案通过
        upload_platon(pip[1].node, pip.cfg.PLATON_NEW_BIN)
        upload_platon(pip[2].node, pip.cfg.PLATON_NEW_BIN)
        upload_platon(pip[3].node, pip.cfg.PLATON_NEW_BIN)
        vote(pips[1], pip_id)
        vote(pips[2], pip_id)
        vote(pips[3], pip_id)
        wait_proposal_Active(pip_id)
        # step3：停止节点，等待节点被零出块处罚
        make_0mb_slash(verifiers[0], verifiers[1])
        upload_platon(pip.node, pip.cfg.PLATON_NEW_BIN)
        version_declare(pip)
        assert vote(pip, pip_id) == 0

    @pytest.mark.parametrize('value', [2])
    def test_modify_0mb_frzzez_time_param(self, verifiers, value):
        """
        修改零出块冻结时长参数-成功  //边界值、大改小、小改大、改动不影响现存，不影响持续处罚判断
        """
        pips = get_pips(verifiers)
        pip = pips[0]
        pip_id = param_proposal(pip, 'slashing', 'zeroProduceFreezeDuration', value)
        votes(pip_id, pips, [1, 1, 1, 1])
        wait_proposal_Active(pip_id)
        assert pip.pip.getGovernParamValue('slashing', 'zeroProduceFreezeDuration') == value
        # TODO:校验实际效果

    @pytest.mark.parametrize('value, code', [(0, 0), (3, 0)])
    def test_modify_0mb_frzzez_time_param_fail(self, verifiers, value, code):
        """
        修改零出块冻结时长参数-失败
        """
        pips = get_pips(verifiers)
        pip = pips[0]
        result = param_proposal(pip, 'slashing', 'zeroProduceFreezeDuration', value)
        assert result == code

    # TODO：
    # 1、整体流程，冻结解冻
    # 2、投票后解质押，锁定时长用例更新
    # 3、投票后，解质押和处罚并行
