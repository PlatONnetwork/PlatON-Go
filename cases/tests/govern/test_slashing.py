import time
from typing import List
import pytest
from common.log import log
from tests.lib import check_node_in_list, upload_platon, wait_block_number
from tests.lib.genesis import to_genesis
from tests.ppos.test_general_punishment import verify_low_block_rate_penalty, get_out_block_penalty_parameters
from tests.lib.client import Client
from tests.lib.config import PipConfig


@pytest.fixture()
def verifiers(clients_consensus, new_genesis_env):
    yield clients_consensus


def get_pips(clients: List[Client]):
    return [c.pip for c in clients]


def version_proposal(pip, to_version, voting_rounds):
    result = pip.submitVersion(pip.node.node_id, str(time.time()), to_version, voting_rounds,
                               pip.node.staking_address,
                               transaction_cfg=pip.cfg.transaction_cfg)
    log.info('submit version proposal result : {}'.format(result))
    return get_proposal_result(pip, pip.cfg.version_proposal, result)


def param_proposal(pip, module, name, value):
    result = pip.submitParam(pip.node.node_id, str(time.time()), module, name, value, pip.node.staking_address,
                             transaction_cfg=pip.cfg.transaction_cfg)
    log.info('submit param proposal result : {}'.format(result))
    return get_proposal_result(pip, pip.cfg.param_proposal, result)


def text_proposal(pip):
    result = pip.submitText(pip.node.node_id, str(time.time()), pip.node.staking_address,
                            transaction_cfg=pip.cfg.transaction_cfg)
    log.info('submit text proposal result:'.format(result))
    return get_proposal_result(pip, pip.cfg.text_proposal, result)


def cancel_proposal(pip, pip_id, voting_rounds):
    result = pip.submitCancel(pip.node.node_id, str(time.time()), voting_rounds, pip_id,
                              pip.node.staking_address, transaction_cfg=pip.cfg.transaction_cfg)
    log.info('submit cancel proposal result : {}'.format(result))
    return get_proposal_result(pip, pip.cfg.cancel_proposal, result)


def get_proposal_result(pip, proposal_type, code):
    if code == 0:
        pip_info = pip.get_effect_proposal_info_of_vote(proposal_type)
        log.info(f"proposal id is {pip_info['ProposalID']}")
        return pip_info['ProposalID']
    log.info(f'proposal return an exception code {code}')
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


def wait_proposal_active(pip, pip_id):
    res = pip.pip.getProposal(pip_id)
    end_block = res['Ret']['EndVotingBlock']
    end_block = pip.economic.get_consensus_switchpoint(end_block)
    log.info(f'wait proposal active block : {end_block}')
    wait_block_number(pip.node, end_block)
    return


def make_0mb_slash(slash_client, check_client):
    """
    构造零出块处罚场景
    """
    slash_node = slash_client.node
    pledge_amount, block_reward, slash_blocks = get_out_block_penalty_parameters(slash_client, slash_node, 'Released')
    log.info(f'slashing param : pledge_amount ({pledge_amount}), block_reward ({block_reward}), slash_blocks ({slash_blocks})')
    log.info("make zero produce block")
    start_num, end_num = verify_low_block_rate_penalty(slash_client, check_client, block_reward, slash_blocks, pledge_amount, 'Released')
    log.info('check Verifier Lists')
    assert check_node_in_list(slash_node.node_id, check_client.ppos.getCandidateList) is False
    assert check_node_in_list(slash_node.node_id, check_client.ppos.getVerifierList) is False
    assert check_node_in_list(slash_node.node_id, check_client.ppos.getValidatorList) is False
    slash_client.node.start()
    return start_num, end_num


class TestSlashing:

    @pytest.mark.P1
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
        vote_info = pip.get_accuverifiers_count(pip_id)
        assert vote_info[0] == 4  # all verifiers
        assert vote_info[1] == 1  # all yeas vote

    @pytest.mark.P1
    def test_0mb_freeze_after_param_vote(self, verifiers, new_genesis_env):
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
        # init: 修改依赖参数的值，并重新部署环境
        genesis = to_genesis(new_genesis_env.genesis_config)
        genesis.economicModel.staking.unStakeFreezeDuration = 10
        new_genesis_env.set_genesis(genesis.to_dict())
        new_genesis_env.deploy_all()
        # step1：提交参数提案并进行投票
        pips = get_pips(verifiers)
        pip = pips[0]
        pip_id = param_proposal(pip, 'slashing', 'zeroProduceFreezeDuration', '5')
        vote(pip, pip_id)
        # step2：停止节点，等待节点被零出块处罚
        make_0mb_slash(verifiers[0], verifiers[1])
        # step3：检查提案和投票信息是否正确
        pip = pips[1]
        vote_info = pip.get_accuverifiers_count(pip_id)
        assert vote_info[0] == 4  # all verifiers
        assert vote_info[1] == 1  # all yeas vote

    @pytest.mark.P1
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
        vote_info = pip.get_accuverifiers_count(pip_id)
        assert vote_info[0] == 4  # all verifiers
        assert vote_info[1] == 1  # all yeas vote

    @pytest.mark.P1
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
        pip_id = cancel_proposal(pip, pip_id, 2)
        vote(pip, pip_id)
        # step2：停止节点，等待节点被零出块处罚
        make_0mb_slash(verifiers[0], verifiers[1])
        # step3：检查提案和投票信息是否正确
        pip = pips[1]
        vote_info = pip.get_accuverifiers_count(pip_id)
        assert vote_info[0] == 4  # all verifiers
        assert vote_info[1] == 1  # all yeas vote

    @pytest.mark.P1
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
        assert version_proposal(pip, pip.cfg.version5, 5) == 302022
        assert param_proposal(pip, 'slashing', 'slashBlocksReward', '10') == 302022
        assert text_proposal(pip) == 302022
        pip_id = version_proposal(pips[1], pips[1].cfg.version5, 5)
        assert cancel_proposal(pip, pip_id, 2) == 302022

    @pytest.mark.P1
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
        assert vote(pip, pip_id) == 302022
        # step3：检查提案和投票信息是否正确
        # pip = pips[1]
        # all_verifiers = pip.get_accu_verifiers_of_proposal(pip_id)
        # all_yeas = pip.get_yeas_of_proposal(pip_id)
        # assert all_verifiers == 1
        # assert all_yeas == 1

    @pytest.mark.P1
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
        pip_id = param_proposal(pip, 'slashing', 'slashBlocksReward', '10')
        # step2：停止节点，等待节点被零出块处罚
        pip = pips[1]
        upload_platon(pip.node, pip.cfg.PLATON_NEW_BIN)
        make_0mb_slash(verifiers[1], verifiers[0])
        assert vote(pip, pip_id) == 302022

    @pytest.mark.P1
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
        assert vote(pip, pip_id) == 302022

    @pytest.mark.P1
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
        assert vote(pip, pip_id) == 302022

    @pytest.mark.P1
    def test_submit_declare_at_0mb_freezing(self, verifiers, new_genesis_env):
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
        upload_platon(pips[1].node, pips[1].cfg.PLATON_NEW_BIN)
        upload_platon(pips[2].node, pips[2].cfg.PLATON_NEW_BIN)
        upload_platon(pips[3].node, pips[3].cfg.PLATON_NEW_BIN)
        vote(pips[1], pip_id)
        vote(pips[2], pip_id)
        vote(pips[3], pip_id)
        # setp3: 在投票期内，构造节点零出块，并等待提案通过
        start_block, end_block = make_0mb_slash(verifiers[0], verifiers[1])
        wait_proposal_active(pip, pip_id)
        # step4：更新零出块节点二进制，进行版本声明
        upload_platon(pip.node, pip.cfg.PLATON_NEW_BIN)
        assert version_declare(pip) == 302023
        # step5: 等待零出块冻结结束，进行版本声明
        wait_block_number(pip.node, end_block)
        assert version_declare(pip) == 0

    @pytest.mark.P1
    def test_proposal_multiple_voting(self, verifiers):
        """
        @describe: 同一提案，节点处罚前已经投票，处罚完成后再次投票
        @step:
        - 1. 提交版本提案并进行投票
        - 2. 停止节点，等待节点被零出块处罚
        - 3. 处罚完成后，再次进行投票
        @expect:
        - 1. 节点重复投票失败，提案未统计重复投票信息
        """
        # step1：
        pips = get_pips(verifiers)
        pip = pips[0]
        pip_id = version_proposal(pip, pip.cfg.version5, 10)
        upload_platon(pip.node, pip.cfg.PLATON_NEW_BIN)
        vote(pip, pip_id)
        # step2：停止节点，等待节点被零出块处罚
        start_block, end_block = make_0mb_slash(verifiers[0], verifiers[1])
        wait_block_number(pip.node, end_block)
        # step3：检查提案和投票信息是否正确
        assert vote(pip, pip_id) == 302027
        vote_info = pip.get_accuverifiers_count(pip_id)
        assert vote_info[0] == 4  # all verifiers
        assert vote_info[1] == 1  # all yeas vote

    @pytest.mark.P1
    @pytest.mark.parametrize('value', ['1', '3'])
    def test_modify_param_of_0mb_freeze_duration(self, verifiers, new_genesis_env, value):
        """
        @describe: 参数提案修改‘zeroProduceFreezeDuration’的值-正常
        @step:
        - 1. 提交参数提案，修改‘zeroProduceFreezeDuration’的值为正常值
        - 2. 检查参数生效值和效果
        @expect:
        - 1. 提案生效后查询该参数，返回正确
        - 2. 参数在链上生效，影响零出块冻结持续时长
        """
        # init: 修改依赖参数的值，并重新部署环境
        genesis = to_genesis(new_genesis_env.genesis_config)
        genesis.economicModel.staking.unStakeFreezeDuration = 4
        genesis.economicModel.slashing.zeroProduceFreezeDuration = 2
        new_genesis_env.set_genesis(genesis.to_dict())
        new_genesis_env.deploy_all()
        # step1: 发起参数提案，投票使提案生效
        pips = get_pips(verifiers)
        pip = pips[0]
        pip_id = param_proposal(pip, 'slashing', 'zeroProduceFreezeDuration', value)
        votes(pip_id, pips, [1, 1, 1, 1])
        wait_proposal_active(pip, pip_id)
        # step2: 检查参数生效值和效果
        res = pip.pip.getGovernParamValue('slashing', 'zeroProduceFreezeDuration')
        assert res['Ret'] == value
        start_block, end_block = make_0mb_slash(verifiers[0], verifiers[1])
        wait_block_number(pip.node, end_block)
        assert check_node_in_list(pip.node.node_id, pip.node.ppos.getCandidateList) is True
        node_info = pip.node.ppos.getCandidateInfo(pip.node.node_id)
        assert node_info['Ret']['Status'] == 0

    @pytest.mark.P1
    @pytest.mark.parametrize('value, code', [('2', 302034), ('0', 3), ('4', 3), ('', 3), ('T', 3)])
    def test_modify_param_fail_of_0mb_freeze_duration(self, verifiers, new_genesis_env, value, code):
        """
        @describe: 参数提案修改‘zeroProduceFreezeDuration’的值-异常
        @step:
        - 1. 提交参数提案，修改‘zeroProduceFreezeDuration’的值为异常值
        - 2. 检查提案异常返回的code
        @expect:
        - 1. 提案返回错误码正确
        """
        # init: 修改依赖参数的值，并重新部署环境
        genesis = to_genesis(new_genesis_env.genesis_config)
        genesis.economicModel.staking.unStakeFreezeDuration = 4
        genesis.economicModel.slashing.zeroProduceFreezeDuration = 2
        new_genesis_env.set_genesis(genesis.to_dict())
        new_genesis_env.deploy_all()
        # step1: 发起参数提案，投票使提案生效
        pips = get_pips(verifiers)
        pip = pips[0]
        assert param_proposal(pip, 'slashing', 'zeroProduceFreezeDuration', value) == code

    @pytest.mark.P1
    def test_all_process_of_pip_after_slashed(self, verifiers):
        """
        @describe: 在处罚结束之后，执行pip所有流程
        @step:
        - 1. 节点在零出块处罚结束后发起提案
        - 2. 节点在零出块处罚结束后进行投票
        - 3. 节点在零出块处罚结束后进行版本声明
        @expect:
        - 1. pip流程未受到处罚影响
        """
        # step1: 构造零出块处罚，并在处罚结束后发起提案
        pips = get_pips(verifiers)
        pip = pips[0]
        _, end_block = make_0mb_slash(verifiers[0], verifiers[1])
        wait_block_number(pip.node, end_block)
        assert check_node_in_list(pip.node.node_id, pip.node.ppos.getCandidateList) is True
        pip.economic.wait_settlement(pip.node)
        assert check_node_in_list(pip.node.node_id, pip.node.ppos.getValidatorList) is True
        assert check_node_in_list(pip.node.node_id, pip.node.ppos.getVerifierList) is True
        # step2: 构造零出块处罚，并在处罚结束后进行投票
        pip_id = version_proposal(pip, pip.cfg.version5, 10)
        _, end_block = make_0mb_slash(verifiers[1], verifiers[0])
        wait_block_number(pip.node, end_block)
        upload_platon(pips[1].node, pip.cfg.PLATON_NEW_BIN)
        upload_platon(pips[2].node, pip.cfg.PLATON_NEW_BIN)
        upload_platon(pips[3].node, pip.cfg.PLATON_NEW_BIN)
        assert vote(pips[1], pip_id, pips[1].cfg.vote_option_yeas) == 0
        assert vote(pips[2], pip_id, pips[2].cfg.vote_option_yeas) == 0
        assert vote(pips[3], pip_id, pips[3].cfg.vote_option_yeas) == 0
        # step3: 在提案生效后进行版本声明
        wait_proposal_active(pip, pip_id)
        upload_platon(pip.node, pip.cfg.PLATON_NEW_BIN)
        assert version_declare(pip) == 0
