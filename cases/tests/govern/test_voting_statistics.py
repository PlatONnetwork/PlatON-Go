import pytest
from common.log import log
from common.key import mock_duplicate_sign
import time
from tests.lib.utils import assert_code, wait_block_number, upload_platon
from tests.lib.client import Client
from tests.lib import Genesis
from dacite import from_dict
from tests.govern.conftest import proposal_vote, version_proposal_vote


def submitvpandvote(client_list_obj, votingrounds=2, version=None):
    pip_obj = client_list_obj[0].pip
    if version is None:
        version = pip_obj.cfg.version5
    result = pip_obj.submitVersion(pip_obj.node.node_id, str(time.time()), version, votingrounds, pip_obj.node.staking_address,
                                   transaction_cfg=pip_obj.cfg.transaction_cfg)
    log.info('Submit version proposal result : {}'.format(result))
    assert_code(result, 0)
    proposalinfo = pip_obj.get_effect_proposal_info_of_vote()
    log.info('Version proposal info {}'.format(proposalinfo))
    for index in range(len(client_list_obj)):
        pip_obj = client_list_obj[index].pip
        result = version_proposal_vote(pip_obj, vote_option=pip_obj.cfg.vote_option_yeas)
        assert_code(result, 0)


def createstaking(obj, platon_bin=None):
    if isinstance(obj, Client):
        obj_list = []
        obj_list.append(obj)
        obj = obj_list
    for client_obj in obj:
        if platon_bin:
            log.info('Need replace the platon of the node')
            upload_platon(client_obj.node, platon_bin)
            client_obj.node.restart()

        address, _ = client_obj.economic.account.generate_account(client_obj.node.web3, 10 ** 18 * 10000000)
        result = client_obj.staking.create_staking(0, address, address, amount=10 ** 18 * 2000000,
                                                   transaction_cfg=client_obj.pip.cfg.transaction_cfg)
        log.info('Node {} staking result : {}'.format(client_obj.node.node_id, result))
        assert_code(result, 0)


def submitppandvote(client_list_obj, *args):
    pip_obj = client_list_obj[0].pip
    result = pip_obj.submitParam(pip_obj.node.node_id, str(time.time()), 'slashing', 'slashBlocksReward', '83',
                                 pip_obj.node.staking_address, transaction_cfg=pip_obj.cfg.transaction_cfg)
    log.info('Submit param proposal result : {}'.format(result))
    assert_code(result, 0)
    proposalinfo = pip_obj.get_effect_proposal_info_of_vote(pip_obj.cfg.param_proposal)
    log.info('Param proposal info {}'.format(proposalinfo))
    for index in range(len(client_list_obj)):
        pip_obj = client_list_obj[index].pip
        log.info('{}'.format(args[index]))
        result = proposal_vote(pip_obj, vote_option=args[index])
        assert_code(result, 0)


def submitcppandvote(client_list_obj, list, voting_rounds=2):
    pip_obj = client_list_obj[0].pip
    result = pip_obj.submitParam(pip_obj.node.node_id, str(time.time()), 'slashing', 'slashBlocksReward', '83',
                                 pip_obj.node.staking_address, transaction_cfg=pip_obj.cfg.transaction_cfg)
    log.info('Submit param proposal result : {}'.format(result))
    assert_code(result, 0)
    proposalinfo_param = pip_obj.get_effect_proposal_info_of_vote(pip_obj.cfg.param_proposal)
    log.info('Param proposal info {}'.format(proposalinfo_param))

    result = pip_obj.submitCancel(pip_obj.node.node_id, str(time.time()), voting_rounds, proposalinfo_param.get('ProposalID'),
                                  pip_obj.node.staking_address, transaction_cfg=pip_obj.cfg.transaction_cfg)
    log.info('Submit cancel proposal result : {}'.format(result))
    assert_code(result, 0)
    proposalinfo_cancel = pip_obj.get_effect_proposal_info_of_vote(pip_obj.cfg.cancel_proposal)
    log.info('Cancel proposal info {}'.format(proposalinfo_cancel))
    for index in range(len(client_list_obj)):
        pip_obj = client_list_obj[index].pip
        log.info('Vote option {}'.format(list[index]))
        result = pip_obj.vote(pip_obj.node.node_id, proposalinfo_cancel.get('ProposalID'), list[index],
                              pip_obj.node.staking_address, transaction_cfg=pip_obj.cfg.transaction_cfg)
        log.info('Node {} vote cancel proposal result : {}'.format(pip_obj.node.node_id, result))
        assert_code(result, 0)


def submitcvpandvote(client_list_obj, *args):
    pip_obj = client_list_obj[0].pip
    result = pip_obj.submitVersion(pip_obj.node.node_id, str(time.time()), pip_obj.cfg.version5, 3,
                                   pip_obj.node.staking_address, transaction_cfg=pip_obj.cfg.transaction_cfg)
    log.info('Submit version proposal result : {}'.format(result))
    assert_code(result, 0)
    proposalinfo_version = pip_obj.get_effect_proposal_info_of_vote(pip_obj.cfg.version_proposal)
    log.info('Version proposal info {}'.format(proposalinfo_version))

    result = pip_obj.submitCancel(pip_obj.node.node_id, str(time.time()), 1, proposalinfo_version.get('ProposalID'),
                                  pip_obj.node.staking_address, transaction_cfg=pip_obj.cfg.transaction_cfg)
    log.info('Submit cancel proposal result : {}'.format(result))
    assert_code(result, 0)
    proposalinfo_cancel = pip_obj.get_effect_proposal_info_of_vote(pip_obj.cfg.cancel_proposal)
    log.info('Cancel proposal info {}'.format(proposalinfo_cancel))
    for index in range(len(client_list_obj)):
        pip_obj = client_list_obj[index].pip
        log.info('{}'.format(args[index]))
        result = pip_obj.vote(pip_obj.node.node_id, proposalinfo_cancel.get('ProposalID'), args[index],
                              pip_obj.node.staking_address, transaction_cfg=pip_obj.cfg.transaction_cfg)
        log.info('Node {} vote cancel proposal result : {}'.format(pip_obj.node.node_id, result))
        assert_code(result, 0)


def submittpandvote(client_list_obj, *args):
    pip_obj = client_list_obj[0].pip
    result = pip_obj.submitText(pip_obj.node.node_id, str(time.time()), pip_obj.node.staking_address,
                                transaction_cfg=pip_obj.cfg.transaction_cfg)
    log.info('Submit text proposal result : {}'.format(result))
    assert_code(result, 0)
    proposalinfo_text = pip_obj.get_effect_proposal_info_of_vote(pip_obj.cfg.text_proposal)
    log.info('Text proposal info {}'.format(proposalinfo_text))

    for index in range(len(client_list_obj)):
        pip_obj = client_list_obj[index].pip
        log.info('{}'.format(args[index]))
        result = pip_obj.vote(pip_obj.node.node_id, proposalinfo_text.get('ProposalID'), args[index],
                              pip_obj.node.staking_address, transaction_cfg=pip_obj.cfg.transaction_cfg)
        log.info('Node {} vote text proposal result : {}'.format(pip_obj.node.node_id, result))
        assert_code(result, 0)


class TestVotingStatisticsVP():
    @pytest.mark.compatibility
    @pytest.mark.P0
    def test_VS_EXV_001_VS_BL_1(self, new_genesis_env, client_con_list_obj, client_noc_list_obj):
        new_genesis_env.deploy_all()
        submitvpandvote(client_con_list_obj[0:-1])
        proposalinfo = client_con_list_obj[0].pip.get_effect_proposal_info_of_vote(client_con_list_obj[0].pip.cfg.version_proposal)
        log.info('Version proposal information {}'.format(proposalinfo))
        createstaking(client_noc_list_obj[:3])
        wait_block_number(client_con_list_obj[0].node, proposalinfo.get('EndVotingBlock'))
        result = client_con_list_obj[0].pip.get_accuverifiers_count(proposalinfo.get('ProposalID'))
        log.info('Get proposal vote infomation {}'.format(result))
        assert result == [4, 3, 0, 0]
        assert client_con_list_obj[0].pip.get_status_of_proposal(proposalinfo.get('ProposalID')) == 4
        assert client_con_list_obj[0].pip.get_yeas_of_proposal(proposalinfo.get('ProposalID')) == 3
        assert client_con_list_obj[0].pip.get_nays_of_proposal(proposalinfo.get('ProposalID')) == 0
        assert client_con_list_obj[0].pip.get_abstentions_of_proposal(proposalinfo.get('ProposalID')) == 0
        assert client_con_list_obj[0].pip.get_accu_verifiers_of_proposal(proposalinfo.get('ProposalID')) == len(
            client_con_list_obj)

    @pytest.mark.P1
    def test_VS_EXV_002(self, new_genesis_env, client_con_list_obj, client_noc_list_obj):
        new_genesis_env.deploy_all()
        pip_obj = client_con_list_obj[0].pip
        submitvpandvote(client_con_list_obj[:2], votingrounds=5)
        proposalinfo = pip_obj.get_effect_proposal_info_of_vote(pip_obj.cfg.version_proposal)
        log.info('Version proposal info {}'.format(proposalinfo))
        log.info('{}'.format(client_con_list_obj[:2]))
        createstaking(client_noc_list_obj[:2])
        pip_obj.economic.wait_settlement_blocknum(pip_obj.node)
        result = pip_obj.get_accuverifiers_count(proposalinfo.get('ProposalID'))
        log.info('Get proposal vote infomation {}'.format(result))
        assert result == [6, 2, 0, 0]

        result = version_proposal_vote(client_con_list_obj[2].pip)
        assert_code(result, 0)
        createstaking(client_noc_list_obj[2])
        wait_block_number(pip_obj.node, proposalinfo.get('EndVotingBlock'))
        result = pip_obj.get_accuverifiers_count(proposalinfo.get('ProposalID'))
        log.info('Get proposal vote infomation {}'.format(result))
        assert result == [6, 3, 0, 0]
        assert client_con_list_obj[0].pip.get_status_of_proposal(proposalinfo.get('ProposalID')) == 3
        assert client_con_list_obj[0].pip.get_yeas_of_proposal(proposalinfo.get('ProposalID')) == 3
        assert client_con_list_obj[0].pip.get_nays_of_proposal(proposalinfo.get('ProposalID')) == 0
        assert client_con_list_obj[0].pip.get_abstentions_of_proposal(proposalinfo.get('ProposalID')) == 0
        assert client_con_list_obj[0].pip.get_accu_verifiers_of_proposal(proposalinfo.get('ProposalID')) == 6

    @pytest.mark.P2
    def test_VS_EXV_003(self, new_genesis_env, client_con_list_obj, client_noc_list_obj):
        new_genesis_env.deploy_all()
        submitvpandvote(client_con_list_obj[:1], votingrounds=9)
        proposalinfo = client_con_list_obj[0].pip.get_effect_proposal_info_of_vote(client_con_list_obj[0].pip.cfg.version_proposal)
        log.info('Version proposal info {}'.format(proposalinfo))
        createstaking(client_noc_list_obj[0])
        client_con_list_obj[0].pip.economic.wait_settlement_blocknum(client_con_list_obj[0].pip.node)
        result = client_con_list_obj[0].pip.get_accuverifiers_count(proposalinfo.get('ProposalID'))
        log.info('Get proposal vote infomation {}'.format(result))
        assert result == [5, 1, 0, 0]

        result = version_proposal_vote(client_con_list_obj[1].pip, client_con_list_obj[0].pip.cfg.vote_option_yeas)
        assert_code(result, 0)
        createstaking(client_noc_list_obj[1])
        client_con_list_obj[0].pip.economic.wait_settlement_blocknum(client_con_list_obj[0].pip.node)
        result = client_con_list_obj[0].pip.get_accuverifiers_count(proposalinfo.get('ProposalID'))
        log.info('Get proposal vote infomation {}'.format(result))
        assert result == [6, 2, 0, 0]

        result = version_proposal_vote(client_con_list_obj[2].pip, client_con_list_obj[0].pip.cfg.vote_option_yeas)
        assert_code(result, 0)
        createstaking(client_noc_list_obj[2])
        wait_block_number(client_con_list_obj[0].pip.node, proposalinfo.get('EndVotingBlock'))

        result = client_con_list_obj[0].pip.get_accuverifiers_count(proposalinfo.get('ProposalID'))
        log.info('Get proposal vote infomation {}'.format(result))
        assert result == [6, 3, 0, 0]
        assert client_con_list_obj[0].pip.get_status_of_proposal(proposalinfo.get('ProposalID')) == 3
        assert client_con_list_obj[0].pip.get_yeas_of_proposal(proposalinfo.get('ProposalID')) == 3
        assert client_con_list_obj[0].pip.get_nays_of_proposal(proposalinfo.get('ProposalID')) == 0
        assert client_con_list_obj[0].pip.get_abstentions_of_proposal(proposalinfo.get('ProposalID')) == 0
        assert client_con_list_obj[0].pip.get_accu_verifiers_of_proposal(proposalinfo.get('ProposalID')) == 6

    @pytest.mark.P2
    def test_VS_EXV_004(self, new_genesis_env, client_con_list_obj):
        genesis = from_dict(data_class=Genesis, data=new_genesis_env.genesis_config)
        genesis.economicModel.gov.versionProposalSupportRate = 0.25
        new_genesis_env.set_genesis(genesis.to_dict())
        new_genesis_env.deploy_all()
        submitvpandvote(client_con_list_obj[:2])
        proposalinfo = client_con_list_obj[0].pip.get_effect_proposal_info_of_vote(client_con_list_obj[0].pip.cfg.version_proposal)
        log.info('Get version proposal information {}'.format(proposalinfo))
        result = client_con_list_obj[0].staking.withdrew_staking(client_con_list_obj[0].node.staking_address)
        log.info('Node {} withdrew staking result {}'.format(client_con_list_obj[0].node.node_id, result))
        assert_code(result, 0)
        wait_block_number(client_con_list_obj[0].node, proposalinfo.get('EndVotingBlock'))
        result = client_con_list_obj[0].pip.get_accuverifiers_count(proposalinfo.get('ProposalID'))
        log.info('Get proposal vote infomation {}'.format(result))
        assert result == [4, 2, 0, 0]
        assert client_con_list_obj[0].pip.get_yeas_of_proposal(proposalinfo.get('ProposalID')) == 2
        assert client_con_list_obj[0].pip.get_status_of_proposal(proposalinfo.get('ProposalID')) == 4

    @pytest.mark.P2
    def test_VS_EXV_005(self, new_genesis_env, client_con_list_obj):
        genesis = from_dict(data_class=Genesis, data=new_genesis_env.genesis_config)
        genesis.economicModel.gov.versionProposalSupportRate = 0.5
        new_genesis_env.set_genesis(genesis.to_dict())
        new_genesis_env.deploy_all()
        submitvpandvote(client_con_list_obj[:2])
        pip_obj = client_con_list_obj[0].pip
        proposalinfo = pip_obj.get_effect_proposal_info_of_vote(pip_obj.cfg.version_proposal)
        log.info('Get version proposal information {}'.format(proposalinfo))
        wait_block_number(pip_obj.node, proposalinfo.get('EndVotingBlock'))
        report_information = mock_duplicate_sign(1, pip_obj.node.nodekey, pip_obj.node.blsprikey,
                                                 proposalinfo.get('EndVotingBlock') - 10)
        log.info("Report information: {}".format(report_information))
        address, _ = pip_obj.economic.account.generate_account(pip_obj.node.web3, 10**18 * 1000)
        result = client_con_list_obj[0].duplicatesign.reportDuplicateSign(1, report_information, address)
        assert_code(result, 0)
        assert client_con_list_obj[1].pip.get_yeas_of_proposal(proposalinfo.get('ProposalID')) == 2

    @pytest.mark.P2
    def test_VS_EXV_006(self, new_genesis_env, client_con_list_obj):
        genesis = from_dict(data_class=Genesis, data=new_genesis_env.genesis_config)
        genesis.economicModel.gov.versionProposalSupportRate = 0.5
        new_genesis_env.set_genesis(genesis.to_dict())
        new_genesis_env.deploy_all()
        submitvpandvote(client_con_list_obj[:2])
        pip_obj = client_con_list_obj[0].pip
        proposalinfo = pip_obj.get_effect_proposal_info_of_vote(pip_obj.cfg.version_proposal)
        log.info('Get version proposal information {}'.format(proposalinfo))
        wait_block_number(pip_obj.node, 80)
        report_information = mock_duplicate_sign(1, pip_obj.node.nodekey, pip_obj.node.blsprikey, 70)
        log.info("Report information: {}".format(report_information))
        address, _ = pip_obj.economic.account.generate_account(pip_obj.node.web3, 10**18 * 1000)
        result = client_con_list_obj[0].duplicatesign.reportDuplicateSign(1, report_information, address)
        assert_code(result, 0)
        wait_block_number(pip_obj.node, proposalinfo.get('EndVotingBlock'))
        assert client_con_list_obj[0].pip.get_yeas_of_proposal(proposalinfo.get('ProposalID')) == 1
        assert client_con_list_obj[0].pip.get_status_of_proposal(proposalinfo.get('ProposalID')) == 3

    @pytest.mark.P2
    def test_VS_EXV_007(self, new_genesis_env, client_con_list_obj):
        genesis = from_dict(data_class=Genesis, data=new_genesis_env.genesis_config)
        genesis.economicModel.gov.versionProposalSupportRate = 0.5
        new_genesis_env.set_genesis(genesis.to_dict())
        new_genesis_env.deploy_all()
        submitvpandvote(client_con_list_obj[:2])
        pip_obj = client_con_list_obj[0].pip
        proposalinfo = pip_obj.get_effect_proposal_info_of_vote(pip_obj.cfg.version_proposal)
        log.info('Get version proposal information {}'.format(proposalinfo))
        wait_block_number(pip_obj.node, proposalinfo.get('EndVotingBlock'))
        report_information = mock_duplicate_sign(2, pip_obj.node.nodekey, pip_obj.node.blsprikey,
                                                 proposalinfo.get('EndVotingBlock') - 10)
        log.info("Report information: {}".format(report_information))
        address, _ = pip_obj.economic.account.generate_account(pip_obj.node.web3, 10**18 * 1000)
        result = client_con_list_obj[0].duplicatesign.reportDuplicateSign(2, report_information, address)
        assert_code(result, 0)
        assert client_con_list_obj[1].pip.get_yeas_of_proposal(proposalinfo.get('ProposalID')) == 2

    @pytest.mark.P2
    def test_VS_EXV_008(self, new_genesis_env, client_con_list_obj):
        genesis = from_dict(data_class=Genesis, data=new_genesis_env.genesis_config)
        genesis.economicModel.gov.versionProposalSupportRate = 0.5
        new_genesis_env.set_genesis(genesis.to_dict())
        new_genesis_env.deploy_all()
        submitvpandvote(client_con_list_obj[:2])
        pip_obj = client_con_list_obj[0].pip
        proposalinfo = pip_obj.get_effect_proposal_info_of_vote(pip_obj.cfg.version_proposal)
        log.info('Get version proposal information {}'.format(proposalinfo))
        wait_block_number(pip_obj.node, 80)
        report_information = mock_duplicate_sign(2, pip_obj.node.nodekey, pip_obj.node.blsprikey, 70)
        log.info("Report information: {}".format(report_information))
        address, _ = pip_obj.economic.account.generate_account(pip_obj.node.web3, 10**18 * 1000)
        result = client_con_list_obj[0].duplicatesign.reportDuplicateSign(2, report_information, address)
        assert_code(result, 0)
        wait_block_number(pip_obj.node, proposalinfo.get('EndVotingBlock'))
        assert client_con_list_obj[0].pip.get_yeas_of_proposal(proposalinfo.get('ProposalID')) == 1
        assert client_con_list_obj[0].pip.get_status_of_proposal(proposalinfo.get('ProposalID')) == 3

    @pytest.mark.P2
    def test_VS_EXV_009(self, new_genesis_env, client_con_list_obj):
        genesis = from_dict(data_class=Genesis, data=new_genesis_env.genesis_config)
        genesis.economicModel.gov.versionProposalSupportRate = 0.5
        new_genesis_env.set_genesis(genesis.to_dict())
        new_genesis_env.deploy_all()
        submitvpandvote(client_con_list_obj[:2], votingrounds=3)
        pip_obj = client_con_list_obj[0].pip
        proposalinfo = pip_obj.get_effect_proposal_info_of_vote(pip_obj.cfg.version_proposal)
        log.info('Get version proposal information {}'.format(proposalinfo))
        assert pip_obj.get_accuverifiers_count(proposalinfo.get('ProposalID')) == [4, 2, 0, 0]
        log.info('Stop the node {}'.format(pip_obj.node.node_id))
        pip_obj.node.stop()
        wait_block_number(client_con_list_obj[1].node, proposalinfo.get('EndVotingBlock'))
        assert client_con_list_obj[1].pip.get_yeas_of_proposal(proposalinfo.get('ProposalID')) == 1
        assert client_con_list_obj[1].pip.get_status_of_proposal(proposalinfo.get('ProposalID')) == 3

    @pytest.mark.P2
    def test_VS_EXV_010(self, new_genesis_env, client_con_list_obj):
        genesis = from_dict(data_class=Genesis, data=new_genesis_env.genesis_config)
        genesis.economicModel.gov.versionProposalSupportRate = 0.5
        new_genesis_env.set_genesis(genesis.to_dict())
        new_genesis_env.deploy_all()
        submitvpandvote(client_con_list_obj[:2], votingrounds=3)
        pip_obj = client_con_list_obj[0].pip
        proposalinfo = pip_obj.get_effect_proposal_info_of_vote(pip_obj.cfg.version_proposal)
        log.info('Get version proposal information {}'.format(proposalinfo))
        assert pip_obj.get_accuverifiers_count(proposalinfo.get('ProposalID')) == [4, 2, 0, 0]
        log.info('Stop the node {}'.format(pip_obj.node.node_id))
        wait_block_number(client_con_list_obj[1].node, proposalinfo.get('EndVotingBlock') - 10)
        pip_obj.node.stop()
        wait_block_number(client_con_list_obj[1].node, proposalinfo.get('EndVotingBlock'))
        assert client_con_list_obj[1].pip.get_yeas_of_proposal(proposalinfo.get('ProposalID')) == 2
        assert client_con_list_obj[1].pip.get_status_of_proposal(proposalinfo.get('ProposalID')) == 4


class TestVotingStatisticsTPCP():
    @pytest.mark.P1
    def test_VS_EXT_001_VS_EXC_001(self, new_genesis_env, client_con_list_obj, client_noc_list_obj):
        genesis = from_dict(data_class=Genesis, data=new_genesis_env.genesis_config)
        genesis.economicModel.gov.textProposalVoteDurationSeconds = 120
        new_genesis_env.set_genesis(genesis.to_dict())
        new_genesis_env.deploy_all()
        submittpandvote(client_con_list_obj[0:-1], 1, 2, 3)
        submitcppandvote(client_con_list_obj[0:-1], [1, 2, 3])
        proposalinfo_text = client_con_list_obj[0].pip.get_effect_proposal_info_of_vote(
            client_con_list_obj[0].pip.cfg.text_proposal)
        log.info('Text proposal information {}'.format(proposalinfo_text))
        proposalinfo_cancel = client_con_list_obj[0].pip.get_effect_proposal_info_of_vote(
            client_con_list_obj[0].pip.cfg.cancel_proposal)
        log.info('Cancel proposal information {}'.format(proposalinfo_cancel))
        createstaking(client_noc_list_obj[:3])
        # client_con_list_obj[0].economic.wait_settlement_blocknum(client_con_list_obj[0].node)
        result_text = client_con_list_obj[0].pip.get_accuverifiers_count(proposalinfo_text.get('ProposalID'))
        log.info('Get text proposal vote infomation {}'.format(result_text))
        result_cancel = client_con_list_obj[0].pip.get_accuverifiers_count(proposalinfo_cancel.get('ProposalID'))
        log.info('Get cancel proposal vote infomation {}'.format(result_cancel))
        assert result_text == [4, 1, 1, 1]
        assert result_cancel == [4, 1, 1, 1]

    @pytest.mark.P1
    def test_VS_EXT_002_VS_EXC_002(self, new_genesis_env, client_con_list_obj, client_noc_list_obj):
        genesis = from_dict(data_class=Genesis, data=new_genesis_env.genesis_config)
        genesis.economicModel.gov.textProposalVoteDurationSeconds = 200
        genesis.economicModel.gov.paramProposalVoteDurationSeconds = 1000
        new_genesis_env.set_genesis(genesis.to_dict())
        new_genesis_env.deploy_all()
        pip_obj = client_con_list_obj[0].pip
        submittpandvote(client_con_list_obj[:2], 1, 2)
        submitcppandvote(client_con_list_obj[:2], [1, 2], voting_rounds=5)
        proposalinfo_text = pip_obj.get_effect_proposal_info_of_vote(pip_obj.cfg.text_proposal)
        log.info('Text proposal info {}'.format(proposalinfo_text))
        proposalinfo_cancel = pip_obj.get_effect_proposal_info_of_vote(pip_obj.cfg.cancel_proposal)
        log.info('Cancel proposal info {}'.format(proposalinfo_cancel))
        createstaking(client_noc_list_obj[:2])
        pip_obj.economic.wait_settlement_blocknum(pip_obj.node)
        result_text = pip_obj.get_accuverifiers_count(proposalinfo_text.get('ProposalID'))
        log.info('Get text proposal vote infomation {}'.format(result_text))
        result_cancel = pip_obj.get_accuverifiers_count(proposalinfo_cancel.get('ProposalID'))
        log.info('Get cancel proposal vote infomation {}'.format(result_cancel))
        assert result_text == [6, 1, 1, 0]
        assert result_cancel == [6, 1, 1, 0]

        result_cancel = proposal_vote(client_con_list_obj[2].pip, vote_option=pip_obj.cfg.vote_option_Abstentions,
                                      proposaltype=pip_obj.cfg.cancel_proposal)
        assert_code(result_cancel, 0)
        result_text = proposal_vote(client_con_list_obj[2].pip, vote_option=pip_obj.cfg.vote_option_Abstentions,
                                    proposaltype=pip_obj.cfg.text_proposal)
        assert_code(result_text, 0)
        createstaking(client_noc_list_obj[2])
        wait_block_number(pip_obj.node, proposalinfo_text.get('EndVotingBlock'))

        result_cancel = pip_obj.get_accuverifiers_count(proposalinfo_cancel.get('ProposalID'))
        log.info('Get proposal vote infomation {}'.format(result_cancel))
        result_text = pip_obj.get_accuverifiers_count(proposalinfo_text.get('ProposalID'))
        log.info('Get proposal vote infomation {}'.format(result_text))

        assert result_text == [6, 1, 1, 1]
        assert result_cancel == [6, 1, 1, 1]

    @pytest.mark.P2
    def test_VS_EXT_003_VS_EXC_003(self, new_genesis_env, client_con_list_obj, client_noc_list_obj):
        genesis = from_dict(data_class=Genesis, data=new_genesis_env.genesis_config)
        genesis.economicModel.gov.textProposalVoteDurationSeconds = 360
        genesis.economicModel.gov.paramProposalVoteDurationSeconds = 600
        new_genesis_env.set_genesis(genesis.to_dict())
        new_genesis_env.deploy_all()
        pip_obj = client_con_list_obj[0].pip
        submittpandvote(client_con_list_obj[:1], 1)
        submitcppandvote(client_con_list_obj[:1], [1], voting_rounds=9)
        proposalinfo_text = pip_obj.get_effect_proposal_info_of_vote(pip_obj.cfg.text_proposal)
        log.info('Text proposal info {}'.format(proposalinfo_text))

        proposalinfo_cancel = pip_obj.get_effect_proposal_info_of_vote(pip_obj.cfg.cancel_proposal)
        log.info('Cancel proposal info {}'.format(proposalinfo_cancel))
        createstaking(client_noc_list_obj[0])
        pip_obj.economic.wait_settlement_blocknum(pip_obj.node)
        result_text = pip_obj.get_accuverifiers_count(proposalinfo_text.get('ProposalID'))
        log.info('Get proposal vote infomation {}'.format(result_text))
        result_cancel = pip_obj.get_accuverifiers_count(proposalinfo_cancel.get('ProposalID'))
        log.info('Get proposal vote infomation {}'.format(result_cancel))
        assert result_text == [5, 1, 0, 0]
        assert result_cancel == [5, 1, 0, 0]

        result = proposal_vote(client_con_list_obj[1].pip, vote_option=pip_obj.cfg.vote_option_nays,
                               proposaltype=pip_obj.cfg.cancel_proposal)
        assert_code(result, 0)
        result = proposal_vote(client_con_list_obj[1].pip, vote_option=pip_obj.cfg.vote_option_nays,
                               proposaltype=pip_obj.cfg.text_proposal)
        assert_code(result, 0)
        createstaking(client_noc_list_obj[1])
        pip_obj.economic.wait_settlement_blocknum(pip_obj.node)
        result_cancel = pip_obj.get_accuverifiers_count(proposalinfo_cancel.get('ProposalID'))
        log.info('Get proposal vote infomation {}'.format(result_cancel))
        result_text = pip_obj.get_accuverifiers_count(proposalinfo_text.get('ProposalID'))
        log.info('Get proposal vote infomation {}'.format(result_text))
        assert result_text == [6, 1, 1, 0]
        assert result_cancel == [6, 1, 1, 0]

        result = proposal_vote(client_con_list_obj[2].pip, vote_option=pip_obj.cfg.vote_option_Abstentions,
                               proposaltype=pip_obj.cfg.text_proposal)
        assert_code(result, 0)
        result = proposal_vote(client_con_list_obj[2].pip, vote_option=pip_obj.cfg.vote_option_Abstentions,
                               proposaltype=pip_obj.cfg.cancel_proposal)
        assert_code(result, 0)
        createstaking(client_noc_list_obj[2])
        # wait_block_number(client_con_list_obj[0].pip.node, proposalinfo.get('EndVotingBlock'))

        result_text = client_con_list_obj[0].pip.get_accuverifiers_count(proposalinfo_text.get('ProposalID'))
        log.info('Get proposal vote infomation {}'.format(result_text))
        result_cancel = client_con_list_obj[0].pip.get_accuverifiers_count(proposalinfo_cancel.get('ProposalID'))
        log.info('Get proposal vote infomation {}'.format(result_cancel))
        assert result_cancel == [6, 1, 1, 1]
        assert result_text == [6, 1, 1, 1]

    def get_block(self, proposalinfo_text, proposalinfo_cancel):
        block1 = proposalinfo_cancel.get('EndVotingBlock')
        block2 = proposalinfo_text.get('EndVotingBlock')
        if block1 > block2:
            return block1
        else:
            return block2

    @pytest.mark.P2
    def test_VS_EXT_004_VS_EXC_004(self, new_genesis_env, client_con_list_obj):
        genesis = from_dict(data_class=Genesis, data=new_genesis_env.genesis_config)
        genesis.economicModel.gov.paramProposalVoteDurationSeconds = 500
        genesis.economicModel.gov.textProposalVoteDurationSeconds = 80
        genesis.economicModel.gov.textProposalSupportRate = 1
        genesis.economicModel.gov.textProposalVoteRate = 0.499
        genesis.economicModel.gov.cancelProposalSupportRate = 1
        genesis.economicModel.gov.cancelProposalVoteRate = 0.499
        new_genesis_env.set_genesis(genesis.to_dict())
        new_genesis_env.deploy_all()
        pip_obj = client_con_list_obj[0].pip
        submittpandvote(client_con_list_obj[:2], 1, 1)
        submitcppandvote(client_con_list_obj[:2], [1, 1])
        proposalinfo_text = pip_obj.get_effect_proposal_info_of_vote(pip_obj.cfg.text_proposal)
        log.info('Get text proposal information {}'.format(proposalinfo_text))
        proposalinfo_cancel = pip_obj.get_effect_proposal_info_of_vote(pip_obj.cfg.cancel_proposal)
        log.info('Get cancel proposal information {}'.format(proposalinfo_cancel))
        result = client_con_list_obj[0].staking.withdrew_staking(client_con_list_obj[0].node.staking_address)
        log.info('Node {} withdrew staking result {}'.format(client_con_list_obj[0].node.node_id, result))
        assert_code(result, 0)
        wait_block_number(client_con_list_obj[0].node, self.get_block(proposalinfo_cancel, proposalinfo_text))
        assert_code(pip_obj.get_yeas_of_proposal(proposalinfo_text.get('ProposalID')), 2)
        assert_code(pip_obj.get_status_of_proposal(proposalinfo_text.get('ProposalID')), 2)
        assert_code(pip_obj.get_yeas_of_proposal(proposalinfo_cancel.get('ProposalID')), 2)
        assert_code(pip_obj.get_status_of_proposal(proposalinfo_cancel.get('ProposalID')), 2)

    @pytest.mark.P2
    def test_VS_EXT_005_VS_EXC_005(self, new_genesis_env, client_con_list_obj):
        genesis = from_dict(data_class=Genesis, data=new_genesis_env.genesis_config)
        genesis.economicModel.gov.paramProposalVoteDurationSeconds = 500
        genesis.economicModel.gov.textProposalVoteDurationSeconds = 80
        genesis.economicModel.gov.textProposalSupportRate = 1
        genesis.economicModel.gov.textProposalVoteRate = 0.499
        genesis.economicModel.gov.cancelProposalSupportRate = 1
        genesis.economicModel.gov.cancelProposalVoteRate = 0.499
        new_genesis_env.set_genesis(genesis.to_dict())
        new_genesis_env.deploy_all()
        pip_obj = client_con_list_obj[0].pip
        submittpandvote(client_con_list_obj[:2], 1, 1)
        submitcppandvote(client_con_list_obj[:2], [1, 1])
        proposalinfo_cancel = pip_obj.get_effect_proposal_info_of_vote(pip_obj.cfg.cancel_proposal)
        log.info('Get cancel proposal information {}'.format(proposalinfo_cancel))
        proposalinfo_text = pip_obj.get_effect_proposal_info_of_vote(pip_obj.cfg.text_proposal)
        log.info('Get text proposal information {}'.format(proposalinfo_text))
        wait_block_number(pip_obj.node, self.get_block(proposalinfo_text, proposalinfo_cancel))
        report_information = mock_duplicate_sign(1, pip_obj.node.nodekey, pip_obj.node.blsprikey,
                                                 proposalinfo_cancel.get('EndVotingBlock') - 10)
        log.info("Report information: {}".format(report_information))
        address, _ = pip_obj.economic.account.generate_account(pip_obj.node.web3, 10 ** 18 * 1000)
        result = client_con_list_obj[0].duplicatesign.reportDuplicateSign(1, report_information, address)
        assert_code(result, 0)
        assert_code(pip_obj.get_yeas_of_proposal(proposalinfo_cancel.get('ProposalID')), 2)
        assert_code(pip_obj.get_yeas_of_proposal(proposalinfo_text.get('ProposalID')), 2)
        assert_code(pip_obj.get_status_of_proposal(proposalinfo_cancel.get('ProposalID')), 2)
        assert_code(pip_obj.get_status_of_proposal(proposalinfo_text.get('ProposalID')), 2)

    @pytest.mark.P2
    def test_VS_EXT_006_VS_EXC_006(self, new_genesis_env, client_con_list_obj):
        genesis = from_dict(data_class=Genesis, data=new_genesis_env.genesis_config)
        genesis.economicModel.gov.paramProposalVoteDurationSeconds = 500
        genesis.economicModel.gov.textProposalVoteDurationSeconds = 80
        genesis.economicModel.gov.textProposalSupportRate = 1
        genesis.economicModel.gov.textProposalVoteRate = 0.499
        genesis.economicModel.gov.cancelProposalSupportRate = 1
        genesis.economicModel.gov.cancelProposalVoteRate = 0.499
        new_genesis_env.set_genesis(genesis.to_dict())
        new_genesis_env.deploy_all()
        pip_obj = client_con_list_obj[0].pip
        submittpandvote(client_con_list_obj[:2], 1, 1)
        submitcppandvote(client_con_list_obj[:2], [1, 1])
        proposalinfo_cancel = pip_obj.get_effect_proposal_info_of_vote(pip_obj.cfg.cancel_proposal)
        log.info('Get cancel proposal information {}'.format(proposalinfo_cancel))
        proposalinfo_text = pip_obj.get_effect_proposal_info_of_vote(pip_obj.cfg.text_proposal)
        log.info('Get text proposal information {}'.format(proposalinfo_text))
        wait_block_number(pip_obj.node, 50)
        report_information = mock_duplicate_sign(1, pip_obj.node.nodekey, pip_obj.node.blsprikey, 45)
        log.info("Report information: {}".format(report_information))
        address, _ = pip_obj.economic.account.generate_account(pip_obj.node.web3, 10 ** 18 * 1000)
        result = client_con_list_obj[0].duplicatesign.reportDuplicateSign(1, report_information, address)
        assert_code(result, 0)
        wait_block_number(pip_obj.node, self.get_block(proposalinfo_cancel, proposalinfo_text))
        assert_code(pip_obj.get_yeas_of_proposal(proposalinfo_cancel.get('ProposalID')), 1)
        assert_code(pip_obj.get_yeas_of_proposal(proposalinfo_text.get('ProposalID')), 1)
        assert_code(pip_obj.get_status_of_proposal(proposalinfo_cancel.get('ProposalID')), 3)
        assert_code(pip_obj.get_status_of_proposal(proposalinfo_text.get('ProposalID')), 3)

    @pytest.mark.P2
    def test_VS_EXT_007_VS_EXC_007(self, new_genesis_env, client_con_list_obj):
        genesis = from_dict(data_class=Genesis, data=new_genesis_env.genesis_config)
        genesis.economicModel.gov.paramProposalVoteDurationSeconds = 500
        genesis.economicModel.gov.textProposalVoteDurationSeconds = 80
        genesis.economicModel.gov.textProposalSupportRate = 1
        genesis.economicModel.gov.textProposalVoteRate = 0.249
        genesis.economicModel.gov.cancelProposalSupportRate = 1
        genesis.economicModel.gov.cancelProposalVoteRate = 0.249
        new_genesis_env.set_genesis(genesis.to_dict())
        new_genesis_env.deploy_all()
        pip_obj = client_con_list_obj[0].pip
        submittpandvote(client_con_list_obj[:2], 2, 1)
        submitcppandvote(client_con_list_obj[:2], [2, 1])
        proposalinfo_cancel = pip_obj.get_effect_proposal_info_of_vote(pip_obj.cfg.cancel_proposal)
        log.info('Get cancel proposal information {}'.format(proposalinfo_cancel))
        proposalinfo_text = pip_obj.get_effect_proposal_info_of_vote(pip_obj.cfg.text_proposal)
        log.info('Get text proposal information {}'.format(proposalinfo_text))
        wait_block_number(pip_obj.node, 50)
        report_information = mock_duplicate_sign(1, pip_obj.node.nodekey, pip_obj.node.blsprikey, 45)
        log.info("Report information: {}".format(report_information))
        address, _ = pip_obj.economic.account.generate_account(pip_obj.node.web3, 10 ** 18 * 1000)
        result = client_con_list_obj[0].duplicatesign.reportDuplicateSign(1, report_information, address)
        assert_code(result, 0)
        wait_block_number(pip_obj.node, self.get_block(proposalinfo_cancel, proposalinfo_text))
        assert pip_obj.get_accuverifiers_count(proposalinfo_text.get('ProposalID')) == [4, 1, 0, 0]
        assert pip_obj.get_accuverifiers_count(proposalinfo_cancel.get('ProposalID')) == [4, 1, 0, 0]
        assert_code(pip_obj.get_status_of_proposal(proposalinfo_cancel.get('ProposalID')), 2)
        assert_code(pip_obj.get_status_of_proposal(proposalinfo_text.get('ProposalID')), 2)

    @pytest.mark.P2
    def test_VS_EXT_008_VS_EXC_008(self, new_genesis_env, client_con_list_obj):
        genesis = from_dict(data_class=Genesis, data=new_genesis_env.genesis_config)
        genesis.economicModel.gov.paramProposalVoteDurationSeconds = 500
        genesis.economicModel.gov.textProposalVoteDurationSeconds = 80
        genesis.economicModel.gov.textProposalSupportRate = 1
        genesis.economicModel.gov.textProposalVoteRate = 0.249
        genesis.economicModel.gov.cancelProposalSupportRate = 1
        genesis.economicModel.gov.cancelProposalVoteRate = 0.249
        new_genesis_env.set_genesis(genesis.to_dict())
        new_genesis_env.deploy_all()
        pip_obj = client_con_list_obj[0].pip
        submittpandvote(client_con_list_obj[:2], 3, 1)
        submitcppandvote(client_con_list_obj[:2], [3, 1])
        proposalinfo_cancel = pip_obj.get_effect_proposal_info_of_vote(pip_obj.cfg.cancel_proposal)
        log.info('Get cancel proposal information {}'.format(proposalinfo_cancel))
        proposalinfo_text = pip_obj.get_effect_proposal_info_of_vote(pip_obj.cfg.text_proposal)
        log.info('Get text proposal information {}'.format(proposalinfo_text))
        wait_block_number(pip_obj.node, 50)
        report_information = mock_duplicate_sign(1, pip_obj.node.nodekey, pip_obj.node.blsprikey, 45)
        log.info("Report information: {}".format(report_information))
        address, _ = pip_obj.economic.account.generate_account(pip_obj.node.web3, 10 ** 18 * 1000)
        result = client_con_list_obj[0].duplicatesign.reportDuplicateSign(1, report_information, address)
        assert_code(result, 0)
        wait_block_number(pip_obj.node, self.get_block(proposalinfo_cancel, proposalinfo_text))
        assert pip_obj.get_accuverifiers_count(proposalinfo_text.get('ProposalID')) == [4, 1, 0, 0]
        assert pip_obj.get_accuverifiers_count(proposalinfo_cancel.get('ProposalID')) == [4, 1, 0, 0]
        assert_code(pip_obj.get_status_of_proposal(proposalinfo_cancel.get('ProposalID')), 2)
        assert_code(pip_obj.get_status_of_proposal(proposalinfo_text.get('ProposalID')), 2)

    @pytest.mark.P2
    def test_VS_EXT_009_VS_EXC_009(self, new_genesis_env, client_con_list_obj):
        genesis = from_dict(data_class=Genesis, data=new_genesis_env.genesis_config)
        genesis.economicModel.gov.paramProposalVoteDurationSeconds = 500
        genesis.economicModel.gov.textProposalVoteDurationSeconds = 80
        genesis.economicModel.gov.textProposalSupportRate = 1
        genesis.economicModel.gov.textProposalVoteRate = 0.249
        genesis.economicModel.gov.cancelProposalSupportRate = 1
        genesis.economicModel.gov.cancelProposalVoteRate = 0.249
        new_genesis_env.set_genesis(genesis.to_dict())
        new_genesis_env.deploy_all()
        pip_obj = client_con_list_obj[0].pip
        submittpandvote(client_con_list_obj[:2], 2, 1)
        submitcppandvote(client_con_list_obj[:2], [2, 1])
        proposalinfo_cancel = pip_obj.get_effect_proposal_info_of_vote(pip_obj.cfg.cancel_proposal)
        log.info('Get cancel proposal information {}'.format(proposalinfo_cancel))
        proposalinfo_text = pip_obj.get_effect_proposal_info_of_vote(pip_obj.cfg.text_proposal)
        log.info('Get text proposal information {}'.format(proposalinfo_text))
        wait_block_number(pip_obj.node, self.get_block(proposalinfo_cancel, proposalinfo_text))
        report_information = mock_duplicate_sign(1, pip_obj.node.nodekey, pip_obj.node.blsprikey, 45)
        log.info("Report information: {}".format(report_information))
        address, _ = pip_obj.economic.account.generate_account(pip_obj.node.web3, 10 ** 18 * 1000)
        result = client_con_list_obj[0].duplicatesign.reportDuplicateSign(1, report_information, address)
        assert_code(result, 0)
        assert pip_obj.get_accuverifiers_count(proposalinfo_text.get('ProposalID')) == [4, 1, 1, 0]
        assert pip_obj.get_accuverifiers_count(proposalinfo_cancel.get('ProposalID')) == [4, 1, 1, 0]
        assert_code(pip_obj.get_status_of_proposal(proposalinfo_cancel.get('ProposalID')), 3)
        assert_code(pip_obj.get_status_of_proposal(proposalinfo_text.get('ProposalID')), 3)

    @pytest.mark.P2
    def test_VS_EXT_010_VS_EXC_010(self, new_genesis_env, client_con_list_obj):
        genesis = from_dict(data_class=Genesis, data=new_genesis_env.genesis_config)
        genesis.economicModel.gov.paramProposalVoteDurationSeconds = 500
        genesis.economicModel.gov.textProposalVoteDurationSeconds = 80
        genesis.economicModel.gov.textProposalSupportRate = 1
        genesis.economicModel.gov.textProposalVoteRate = 0.499
        genesis.economicModel.gov.cancelProposalSupportRate = 1
        genesis.economicModel.gov.cancelProposalVoteRate = 0.499
        new_genesis_env.set_genesis(genesis.to_dict())
        new_genesis_env.deploy_all()
        pip_obj = client_con_list_obj[0].pip
        submittpandvote(client_con_list_obj[:2], 1, 1)
        submitcppandvote(client_con_list_obj[:2], [1, 1])
        proposalinfo_cancel = pip_obj.get_effect_proposal_info_of_vote(pip_obj.cfg.cancel_proposal)
        log.info('Get cancel proposal information {}'.format(proposalinfo_cancel))
        proposalinfo_text = pip_obj.get_effect_proposal_info_of_vote(pip_obj.cfg.text_proposal)
        log.info('Get text proposal information {}'.format(proposalinfo_text))
        wait_block_number(pip_obj.node, 50)
        report_information = mock_duplicate_sign(2, pip_obj.node.nodekey, pip_obj.node.blsprikey, 45)
        log.info("Report information: {}".format(report_information))
        address, _ = pip_obj.economic.account.generate_account(pip_obj.node.web3, 10 ** 18 * 1000)
        result = client_con_list_obj[0].duplicatesign.reportDuplicateSign(2, report_information, address)
        assert_code(result, 0)
        wait_block_number(pip_obj.node, self.get_block(proposalinfo_cancel, proposalinfo_text))
        assert pip_obj.get_accuverifiers_count(proposalinfo_text.get('ProposalID')) == [4, 1, 0, 0]
        assert pip_obj.get_accuverifiers_count(proposalinfo_cancel.get('ProposalID')) == [4, 1, 0, 0]
        assert_code(pip_obj.get_status_of_proposal(proposalinfo_cancel.get('ProposalID')), 3)
        assert_code(pip_obj.get_status_of_proposal(proposalinfo_text.get('ProposalID')), 3)

    @pytest.mark.P2
    def test_VS_EXT_011(self, new_genesis_env, client_con_list_obj):
        genesis = from_dict(data_class=Genesis, data=new_genesis_env.genesis_config)
        genesis.economicModel.gov.paramProposalVoteDurationSeconds = 500
        genesis.economicModel.gov.textProposalVoteDurationSeconds = 80
        genesis.economicModel.gov.textProposalSupportRate = 1
        genesis.economicModel.gov.textProposalVoteRate = 0.249
        genesis.economicModel.gov.cancelProposalSupportRate = 1
        genesis.economicModel.gov.cancelProposalVoteRate = 0.249
        new_genesis_env.set_genesis(genesis.to_dict())
        new_genesis_env.deploy_all()
        pip_obj = client_con_list_obj[0].pip
        submittpandvote(client_con_list_obj[:2], 2, 1)
        submitcppandvote(client_con_list_obj[:2], [2, 1])
        proposalinfo_cancel = pip_obj.get_effect_proposal_info_of_vote(pip_obj.cfg.cancel_proposal)
        log.info('Get cancel proposal information {}'.format(proposalinfo_cancel))
        proposalinfo_text = pip_obj.get_effect_proposal_info_of_vote(pip_obj.cfg.text_proposal)
        log.info('Get text proposal information {}'.format(proposalinfo_text))
        wait_block_number(pip_obj.node, 50)
        report_information = mock_duplicate_sign(2, pip_obj.node.nodekey, pip_obj.node.blsprikey, 45)
        log.info("Report information: {}".format(report_information))
        address, _ = pip_obj.economic.account.generate_account(pip_obj.node.web3, 10 ** 18 * 1000)
        result = client_con_list_obj[0].duplicatesign.reportDuplicateSign(2, report_information, address)
        assert_code(result, 0)
        wait_block_number(pip_obj.node, self.get_block(proposalinfo_cancel, proposalinfo_text))
        assert pip_obj.get_accuverifiers_count(proposalinfo_text.get('ProposalID')) == [4, 1, 0, 0]
        assert pip_obj.get_accuverifiers_count(proposalinfo_cancel.get('ProposalID')) == [4, 1, 0, 0]
        assert_code(pip_obj.get_status_of_proposal(proposalinfo_cancel.get('ProposalID')), 2)
        assert_code(pip_obj.get_status_of_proposal(proposalinfo_text.get('ProposalID')), 2)

    @pytest.mark.P2
    def test_VS_EXT_012(self, new_genesis_env, client_con_list_obj):
        genesis = from_dict(data_class=Genesis, data=new_genesis_env.genesis_config)
        genesis.economicModel.gov.paramProposalVoteDurationSeconds = 500
        genesis.economicModel.gov.textProposalVoteDurationSeconds = 80
        genesis.economicModel.gov.textProposalSupportRate = 1
        genesis.economicModel.gov.textProposalVoteRate = 0.249
        genesis.economicModel.gov.cancelProposalSupportRate = 1
        genesis.economicModel.gov.cancelProposalVoteRate = 0.249
        new_genesis_env.set_genesis(genesis.to_dict())
        new_genesis_env.deploy_all()
        pip_obj = client_con_list_obj[0].pip
        submittpandvote(client_con_list_obj[:2], 3, 1)
        submitcppandvote(client_con_list_obj[:2], [3, 1])
        proposalinfo_cancel = pip_obj.get_effect_proposal_info_of_vote(pip_obj.cfg.cancel_proposal)
        log.info('Get cancel proposal information {}'.format(proposalinfo_cancel))
        proposalinfo_text = pip_obj.get_effect_proposal_info_of_vote(pip_obj.cfg.text_proposal)
        log.info('Get text proposal information {}'.format(proposalinfo_text))
        wait_block_number(pip_obj.node, 50)
        report_information = mock_duplicate_sign(2, pip_obj.node.nodekey, pip_obj.node.blsprikey, 45)
        log.info("Report information: {}".format(report_information))
        address, _ = pip_obj.economic.account.generate_account(pip_obj.node.web3, 10 ** 18 * 1000)
        result = client_con_list_obj[0].duplicatesign.reportDuplicateSign(2, report_information, address)
        assert_code(result, 0)
        wait_block_number(pip_obj.node, self.get_block(proposalinfo_cancel, proposalinfo_text))
        assert pip_obj.get_accuverifiers_count(proposalinfo_text.get('ProposalID')) == [4, 1, 0, 0]
        assert pip_obj.get_accuverifiers_count(proposalinfo_cancel.get('ProposalID')) == [4, 1, 0, 0]
        assert_code(pip_obj.get_status_of_proposal(proposalinfo_cancel.get('ProposalID')), 2)
        assert_code(pip_obj.get_status_of_proposal(proposalinfo_text.get('ProposalID')), 2)

    @pytest.mark.P2
    def test_VS_EXT_013(self, new_genesis_env, client_con_list_obj):
        genesis = from_dict(data_class=Genesis, data=new_genesis_env.genesis_config)
        genesis.economicModel.gov.paramProposalVoteDurationSeconds = 500
        genesis.economicModel.gov.textProposalVoteDurationSeconds = 120
        genesis.economicModel.gov.textProposalSupportRate = 1
        genesis.economicModel.gov.textProposalVoteRate = 0.499
        genesis.economicModel.gov.cancelProposalSupportRate = 1
        genesis.economicModel.gov.cancelProposalVoteRate = 0.499
        new_genesis_env.set_genesis(genesis.to_dict())
        new_genesis_env.deploy_all()
        pip_obj = client_con_list_obj[0].pip
        pip_obj_test = client_con_list_obj[1].pip
        submittpandvote(client_con_list_obj[:2], 1, 1)
        submitcppandvote(client_con_list_obj[:2], [1, 1], voting_rounds=3)
        proposalinfo_cancel = pip_obj.get_effect_proposal_info_of_vote(pip_obj.cfg.cancel_proposal)
        log.info('Get cancel proposal information {}'.format(proposalinfo_cancel))
        proposalinfo_text = pip_obj.get_effect_proposal_info_of_vote(pip_obj.cfg.text_proposal)
        log.info('Get text proposal information {}'.format(proposalinfo_text))
        pip_obj.node.stop()
        wait_block_number(pip_obj_test.node, self.get_block(proposalinfo_cancel, proposalinfo_text))
        assert pip_obj_test.get_accuverifiers_count(proposalinfo_text.get('ProposalID')) == [4, 1, 0, 0]
        assert pip_obj_test.get_accuverifiers_count(proposalinfo_cancel.get('ProposalID')) == [4, 1, 0, 0]
        assert_code(pip_obj_test.get_status_of_proposal(proposalinfo_cancel.get('ProposalID')), 3)
        assert_code(pip_obj_test.get_status_of_proposal(proposalinfo_text.get('ProposalID')), 3)

    @pytest.mark.P2
    def test_VS_EXT_014(self, new_genesis_env, client_con_list_obj):
        genesis = from_dict(data_class=Genesis, data=new_genesis_env.genesis_config)
        genesis.economicModel.gov.paramProposalVoteDurationSeconds = 500
        genesis.economicModel.gov.textProposalVoteDurationSeconds = 120
        genesis.economicModel.gov.textProposalSupportRate = 1
        genesis.economicModel.gov.textProposalVoteRate = 0.499
        genesis.economicModel.gov.cancelProposalSupportRate = 1
        genesis.economicModel.gov.cancelProposalVoteRate = 0.499
        new_genesis_env.set_genesis(genesis.to_dict())
        new_genesis_env.deploy_all()
        pip_obj = client_con_list_obj[0].pip
        pip_obj_test = client_con_list_obj[1].pip
        submittpandvote(client_con_list_obj[:2], 2, 1)
        submitcppandvote(client_con_list_obj[:2], [2, 1], voting_rounds=3)
        proposalinfo_cancel = pip_obj.get_effect_proposal_info_of_vote(pip_obj.cfg.cancel_proposal)
        log.info('Get cancel proposal information {}'.format(proposalinfo_cancel))
        proposalinfo_text = pip_obj.get_effect_proposal_info_of_vote(pip_obj.cfg.text_proposal)
        log.info('Get text proposal information {}'.format(proposalinfo_text))
        pip_obj.node.stop()
        wait_block_number(pip_obj_test.node, self.get_block(proposalinfo_cancel, proposalinfo_text))
        assert pip_obj_test.get_accuverifiers_count(proposalinfo_text.get('ProposalID')) == [4, 1, 0, 0]
        assert pip_obj_test.get_accuverifiers_count(proposalinfo_cancel.get('ProposalID')) == [4, 1, 0, 0]
        assert_code(pip_obj_test.get_status_of_proposal(proposalinfo_cancel.get('ProposalID')), 3)
        assert_code(pip_obj_test.get_status_of_proposal(proposalinfo_text.get('ProposalID')), 3)

    @pytest.mark.P2
    def test_VS_EXT_015(self, new_genesis_env, client_con_list_obj):
        genesis = from_dict(data_class=Genesis, data=new_genesis_env.genesis_config)
        genesis.economicModel.gov.paramProposalVoteDurationSeconds = 500
        genesis.economicModel.gov.textProposalVoteDurationSeconds = 120
        genesis.economicModel.gov.textProposalSupportRate = 1
        genesis.economicModel.gov.textProposalVoteRate = 0.499
        genesis.economicModel.gov.cancelProposalSupportRate = 1
        genesis.economicModel.gov.cancelProposalVoteRate = 0.499
        new_genesis_env.set_genesis(genesis.to_dict())
        new_genesis_env.deploy_all()
        pip_obj = client_con_list_obj[0].pip
        pip_obj_test = client_con_list_obj[1].pip
        submittpandvote(client_con_list_obj[:2], 3, 1)
        submitcppandvote(client_con_list_obj[:2], [3, 1], voting_rounds=3)
        proposalinfo_cancel = pip_obj.get_effect_proposal_info_of_vote(pip_obj.cfg.cancel_proposal)
        log.info('Get cancel proposal information {}'.format(proposalinfo_cancel))
        proposalinfo_text = pip_obj.get_effect_proposal_info_of_vote(pip_obj.cfg.text_proposal)
        log.info('Get text proposal information {}'.format(proposalinfo_text))
        pip_obj.node.stop()
        wait_block_number(pip_obj_test.node, self.get_block(proposalinfo_cancel, proposalinfo_text))
        assert pip_obj_test.get_accuverifiers_count(proposalinfo_text.get('ProposalID')) == [4, 1, 0, 0]
        assert pip_obj_test.get_accuverifiers_count(proposalinfo_cancel.get('ProposalID')) == [4, 1, 0, 0]
        assert_code(pip_obj_test.get_status_of_proposal(proposalinfo_cancel.get('ProposalID')), 3)
        assert_code(pip_obj_test.get_status_of_proposal(proposalinfo_text.get('ProposalID')), 3)

    @pytest.mark.P2
    def test_VS_EXT_016(self, new_genesis_env, client_con_list_obj):
        genesis = from_dict(data_class=Genesis, data=new_genesis_env.genesis_config)
        genesis.economicModel.gov.paramProposalVoteDurationSeconds = 500
        genesis.economicModel.gov.textProposalVoteDurationSeconds = 80
        genesis.economicModel.gov.textProposalSupportRate = 1
        genesis.economicModel.gov.textProposalVoteRate = 0.499
        genesis.economicModel.gov.cancelProposalSupportRate = 1
        genesis.economicModel.gov.cancelProposalVoteRate = 0.499
        new_genesis_env.set_genesis(genesis.to_dict())
        new_genesis_env.deploy_all()
        pip_obj = client_con_list_obj[0].pip
        pip_obj_test = client_con_list_obj[1].pip
        submittpandvote(client_con_list_obj[:2], 1, 1)
        submitcppandvote(client_con_list_obj[:2], [1, 1], voting_rounds=2)
        proposalinfo_cancel = pip_obj.get_effect_proposal_info_of_vote(pip_obj.cfg.cancel_proposal)
        log.info('Get cancel proposal information {}'.format(proposalinfo_cancel))
        proposalinfo_text = pip_obj.get_effect_proposal_info_of_vote(pip_obj.cfg.text_proposal)
        log.info('Get text proposal information {}'.format(proposalinfo_text))
        wait_block_number(pip_obj_test.node, self.get_block(proposalinfo_cancel, proposalinfo_text))
        pip_obj.node.stop()
        assert pip_obj_test.get_accuverifiers_count(proposalinfo_text.get('ProposalID')) == [4, 2, 0, 0]
        assert pip_obj_test.get_accuverifiers_count(proposalinfo_cancel.get('ProposalID')) == [4, 2, 0, 0]
        assert_code(pip_obj_test.get_status_of_proposal(proposalinfo_cancel.get('ProposalID')), 2)
        assert_code(pip_obj_test.get_status_of_proposal(proposalinfo_text.get('ProposalID')), 2)


class TestVotingStatisticsPP():
    @pytest.mark.P1
    def test_VS_EP_004(self, new_genesis_env, client_con_list_obj, client_noc_list_obj):
        genesis = from_dict(data_class=Genesis, data=new_genesis_env.genesis_config)
        genesis.economicModel.gov.paramProposalVoteDurationSeconds = 0
        new_genesis_env.set_genesis(genesis.to_dict())
        new_genesis_env.deploy_all()
        submitppandvote(client_con_list_obj[0:-1], 1, 2, 3)
        proposalinfo = client_con_list_obj[0].pip.get_effect_proposal_info_of_vote(client_con_list_obj[0].pip.cfg.param_proposal)
        log.info('Param proposal information {}'.format(proposalinfo))
        createstaking(client_noc_list_obj[:3])
        # client_con_list_obj[0].economic.wait_settlement_blocknum(client_con_list_obj[0].node)
        result = client_con_list_obj[0].pip.get_accuverifiers_count(proposalinfo.get('ProposalID'))
        log.info('Get proposal vote infomation {}'.format(result))
        assert result == [4, 1, 1, 1]

    @pytest.mark.P1
    def test_VS_EP_005(self, new_genesis_env, client_con_list_obj, client_noc_list_obj):
        genesis = from_dict(data_class=Genesis, data=new_genesis_env.genesis_config)
        genesis.economicModel.gov.paramProposalVoteDurationSeconds = 160
        new_genesis_env.set_genesis(genesis.to_dict())
        new_genesis_env.deploy_all()
        pip_obj = client_con_list_obj[0].pip
        submitppandvote(client_con_list_obj[:2], 1, 2)
        proposalinfo = pip_obj.get_effect_proposal_info_of_vote(pip_obj.cfg.param_proposal)
        log.info('Param proposal info {}'.format(proposalinfo))
        log.info('{}'.format(client_con_list_obj[:2]))
        createstaking(client_noc_list_obj[:2])
        pip_obj.economic.wait_settlement_blocknum(pip_obj.node)
        result = pip_obj.get_accuverifiers_count(proposalinfo.get('ProposalID'))
        log.info('Get proposal vote infomation {}'.format(result))
        assert result == [6, 1, 1, 0]

        result = proposal_vote(client_noc_list_obj[0].pip, vote_option=pip_obj.cfg.vote_option_Abstentions)
        assert_code(result, 0)
        log.info('{}'.format(client_con_list_obj[2]))
        createstaking(client_noc_list_obj[2])
        wait_block_number(pip_obj.node, proposalinfo.get('EndVotingBlock'))

        result = pip_obj.get_accuverifiers_count(proposalinfo.get('ProposalID'))
        log.info('Get proposal vote infomation {}'.format(result))
        assert result == [6, 1, 1, 1]

    @pytest.mark.P2
    def test_VS_EP_006(self, new_genesis_env, client_con_list_obj, client_noc_list_obj):
        genesis = from_dict(data_class=Genesis, data=new_genesis_env.genesis_config)
        genesis.economicModel.gov.paramProposalVoteDurationSeconds = 320
        new_genesis_env.set_genesis(genesis.to_dict())
        new_genesis_env.deploy_all()
        submitppandvote(client_con_list_obj[:1], 1)
        proposalinfo = client_con_list_obj[0].pip.get_effect_proposal_info_of_vote(client_con_list_obj[0].pip.cfg.param_proposal)
        log.info('Param proposal info {}'.format(proposalinfo))
        createstaking(client_noc_list_obj[0])
        client_con_list_obj[0].pip.economic.wait_settlement_blocknum(client_con_list_obj[0].pip.node)
        result = client_con_list_obj[0].pip.get_accuverifiers_count(proposalinfo.get('ProposalID'))
        log.info('Get proposal vote infomation {}'.format(result))
        assert result == [5, 1, 0, 0]

        result = proposal_vote(client_con_list_obj[1].pip, vote_option=client_con_list_obj[0].pip.cfg.vote_option_nays)
        assert_code(result, 0)
        createstaking(client_noc_list_obj[1])
        client_con_list_obj[0].pip.economic.wait_settlement_blocknum(client_con_list_obj[0].pip.node)
        result = client_con_list_obj[0].pip.get_accuverifiers_count(proposalinfo.get('ProposalID'))
        log.info('Get proposal vote infomation {}'.format(result))
        assert result == [6, 1, 1, 0]

        result = proposal_vote(client_con_list_obj[2].pip, vote_option=client_con_list_obj[0].pip.cfg.vote_option_Abstentions)
        assert_code(result, 0)
        createstaking(client_noc_list_obj[2])
        wait_block_number(client_con_list_obj[0].pip.node, proposalinfo.get('EndVotingBlock'))

        result = client_con_list_obj[0].pip.get_accuverifiers_count(proposalinfo.get('ProposalID'))
        log.info('Get proposal vote infomation {}'.format(result))
        assert result == [6, 1, 1, 1]

    @pytest.mark.P0
    def test_VS_EP_007_VS_EP_003(self, new_genesis_env, client_con_list_obj):
        genesis = from_dict(data_class=Genesis, data=new_genesis_env.genesis_config)
        genesis.economicModel.gov.paramProposalVoteDurationSeconds = 0
        genesis.economicModel.gov.paramProposalSupportRate = 1
        genesis.economicModel.gov.paramProposalVoteRate = 0.49
        new_genesis_env.set_genesis(genesis.to_dict())
        new_genesis_env.deploy_all()
        submitppandvote(client_con_list_obj[:2], 1, 1)
        proposalinfo = client_con_list_obj[0].pip.get_effect_proposal_info_of_vote(client_con_list_obj[0].pip.cfg.param_proposal)
        log.info('Get param proposal information {}'.format(proposalinfo))
        result = client_con_list_obj[0].staking.withdrew_staking(client_con_list_obj[0].node.staking_address)
        log.info('Node {} withdrew staking result {}'.format(client_con_list_obj[0].node.node_id, result))
        assert_code(result, 0)
        result = client_con_list_obj[0].pip.pip.getTallyResult(proposalinfo.get('ProposalID'))
        log.info('Before endvoting block, get Tally resul of the parameter proposal result : {}'.format(result))
        assert_code(result, 302030)
        wait_block_number(client_con_list_obj[0].node, proposalinfo.get('EndVotingBlock'))
        assert client_con_list_obj[0].pip.get_yeas_of_proposal(proposalinfo.get('ProposalID')) == 2
        assert client_con_list_obj[0].pip.get_status_of_proposal(proposalinfo.get('ProposalID')) == 2

    @pytest.mark.P2
    def test_VS_EP_008(self, new_genesis_env, client_con_list_obj):
        genesis = from_dict(data_class=Genesis, data=new_genesis_env.genesis_config)
        genesis.economicModel.gov.paramProposalVoteDurationSeconds = 0
        genesis.economicModel.gov.paramProposalSupportRate = 0.5
        genesis.economicModel.gov.paramProposalVoteRate = 0.5
        new_genesis_env.set_genesis(genesis.to_dict())
        new_genesis_env.deploy_all()
        submitppandvote(client_con_list_obj[:2], 1, 1)
        pip_obj = client_con_list_obj[0].pip
        proposalinfo = pip_obj.get_effect_proposal_info_of_vote(pip_obj.cfg.param_proposal)
        log.info('Get param proposal information {}'.format(proposalinfo))
        wait_block_number(pip_obj.node, proposalinfo.get('EndVotingBlock'))
        report_information = mock_duplicate_sign(1, pip_obj.node.nodekey, pip_obj.node.blsprikey,
                                                 proposalinfo.get('EndVotingBlock') - 10)
        log.info("Report information: {}".format(report_information))
        address, _ = pip_obj.economic.account.generate_account(pip_obj.node.web3, 10**18 * 1000)
        result = client_con_list_obj[0].duplicatesign.reportDuplicateSign(1, report_information, address)
        assert_code(result, 0)
        assert client_con_list_obj[1].pip.get_yeas_of_proposal(proposalinfo.get('ProposalID')) == 2

    @pytest.mark.P2
    def test_VS_EP_009(self, new_genesis_env, client_con_list_obj):
        genesis = from_dict(data_class=Genesis, data=new_genesis_env.genesis_config)
        genesis.economicModel.gov.paramProposalVoteDurationSeconds = 0
        genesis.economicModel.gov.paramProposalSupportRate = 0.5
        genesis.economicModel.gov.paramProposalVoteRate = 0.5
        new_genesis_env.set_genesis(genesis.to_dict())
        new_genesis_env.deploy_all()
        submitppandvote(client_con_list_obj[:2], 1, 1)
        pip_obj = client_con_list_obj[0].pip
        proposalinfo = pip_obj.get_effect_proposal_info_of_vote(pip_obj.cfg.param_proposal)
        log.info('Get param proposal information {}'.format(proposalinfo))
        wait_block_number(pip_obj.node, 80)
        report_information = mock_duplicate_sign(1, pip_obj.node.nodekey, pip_obj.node.blsprikey, 70)
        log.info("Report information: {}".format(report_information))
        address, _ = pip_obj.economic.account.generate_account(pip_obj.node.web3, 10**18 * 1000)
        result = client_con_list_obj[0].duplicatesign.reportDuplicateSign(1, report_information, address)
        assert_code(result, 0)
        wait_block_number(pip_obj.node, proposalinfo.get('EndVotingBlock'))
        assert client_con_list_obj[1].pip.get_yeas_of_proposal(proposalinfo.get('ProposalID')) == 1

    @pytest.mark.P2
    def test_VS_EP_010(self, new_genesis_env, client_con_list_obj):
        genesis = from_dict(data_class=Genesis, data=new_genesis_env.genesis_config)
        genesis.economicModel.gov.paramProposalVoteDurationSeconds = 0
        genesis.economicModel.gov.paramProposalSupportRate = 0.5
        genesis.economicModel.gov.paramProposalVoteRate = 0.5
        new_genesis_env.set_genesis(genesis.to_dict())
        new_genesis_env.deploy_all()
        submitppandvote(client_con_list_obj[:2], 2, 2)
        pip_obj = client_con_list_obj[0].pip
        proposalinfo = pip_obj.get_effect_proposal_info_of_vote(pip_obj.cfg.param_proposal)
        log.info('Get param proposal information {}'.format(proposalinfo))
        wait_block_number(pip_obj.node, 80)
        report_information = mock_duplicate_sign(1, pip_obj.node.nodekey, pip_obj.node.blsprikey, 70)
        log.info("Report information: {}".format(report_information))
        address, _ = pip_obj.economic.account.generate_account(pip_obj.node.web3, 10**18 * 1000)
        result = client_con_list_obj[0].duplicatesign.reportDuplicateSign(1, report_information, address)
        assert_code(result, 0)
        wait_block_number(pip_obj.node, proposalinfo.get('EndVotingBlock'))
        assert client_con_list_obj[1].pip.get_nays_of_proposal(proposalinfo.get('ProposalID')) == 1

    @pytest.mark.P2
    def test_VS_EP_011(self, new_genesis_env, client_con_list_obj):
        genesis = from_dict(data_class=Genesis, data=new_genesis_env.genesis_config)
        genesis.economicModel.gov.paramProposalVoteDurationSeconds = 0
        genesis.economicModel.gov.paramProposalSupportRate = 0.5
        genesis.economicModel.gov.paramProposalVoteRate = 0.5
        new_genesis_env.set_genesis(genesis.to_dict())
        new_genesis_env.deploy_all()
        submitppandvote(client_con_list_obj[:2], 3, 3)
        pip_obj = client_con_list_obj[0].pip
        proposalinfo = pip_obj.get_effect_proposal_info_of_vote(pip_obj.cfg.param_proposal)
        log.info('Get param proposal information {}'.format(proposalinfo))
        wait_block_number(pip_obj.node, 80)
        report_information = mock_duplicate_sign(1, pip_obj.node.nodekey, pip_obj.node.blsprikey, 70)
        log.info("Report information: {}".format(report_information))
        address, _ = pip_obj.economic.account.generate_account(pip_obj.node.web3, 10**18 * 1000)
        result = client_con_list_obj[0].duplicatesign.reportDuplicateSign(1, report_information, address)
        assert_code(result, 0)
        wait_block_number(pip_obj.node, proposalinfo.get('EndVotingBlock'))
        assert client_con_list_obj[1].pip.get_abstentions_of_proposal(proposalinfo.get('ProposalID')) == 1

    @pytest.mark.P2
    def test_VS_EP_012(self, new_genesis_env, client_con_list_obj):
        genesis = from_dict(data_class=Genesis, data=new_genesis_env.genesis_config)
        genesis.economicModel.gov.paramProposalVoteDurationSeconds = 0
        genesis.economicModel.gov.paramProposalSupportRate = 0.5
        genesis.economicModel.gov.paramProposalVoteRate = 0.5
        new_genesis_env.set_genesis(genesis.to_dict())
        new_genesis_env.deploy_all()
        submitppandvote(client_con_list_obj[:2], 1, 1)
        pip_obj = client_con_list_obj[0].pip
        proposalinfo = pip_obj.get_effect_proposal_info_of_vote(pip_obj.cfg.param_proposal)
        log.info('Get param proposal information {}'.format(proposalinfo))
        wait_block_number(pip_obj.node, proposalinfo.get('EndVotingBlock'))
        report_information = mock_duplicate_sign(2, pip_obj.node.nodekey, pip_obj.node.blsprikey,
                                                 proposalinfo.get('EndVotingBlock') - 10)
        log.info("Report information: {}".format(report_information))
        address, _ = pip_obj.economic.account.generate_account(pip_obj.node.web3, 10**18 * 1000)
        result = client_con_list_obj[0].duplicatesign.reportDuplicateSign(2, report_information, address)
        assert_code(result, 0)
        assert client_con_list_obj[1].pip.get_yeas_of_proposal(proposalinfo.get('ProposalID')) == 2

    @pytest.mark.P2
    def test_VS_EP_013(self, new_genesis_env, client_con_list_obj):
        genesis = from_dict(data_class=Genesis, data=new_genesis_env.genesis_config)
        genesis.economicModel.gov.paramProposalVoteDurationSeconds = 0
        genesis.economicModel.gov.paramProposalSupportRate = 0.5
        genesis.economicModel.gov.paramProposalVoteRate = 0.5
        new_genesis_env.set_genesis(genesis.to_dict())
        new_genesis_env.deploy_all()
        submitppandvote(client_con_list_obj[:2], 1, 1)
        pip_obj = client_con_list_obj[0].pip
        proposalinfo = pip_obj.get_effect_proposal_info_of_vote(pip_obj.cfg.param_proposal)
        log.info('Get param proposal information {}'.format(proposalinfo))
        wait_block_number(pip_obj.node, 80)
        report_information = mock_duplicate_sign(2, pip_obj.node.nodekey, pip_obj.node.blsprikey, 70)
        log.info("Report information: {}".format(report_information))
        address, _ = pip_obj.economic.account.generate_account(pip_obj.node.web3, 10**18 * 1000)
        result = client_con_list_obj[0].duplicatesign.reportDuplicateSign(2, report_information, address)
        assert_code(result, 0)
        wait_block_number(pip_obj.node, proposalinfo.get('EndVotingBlock'))
        assert client_con_list_obj[1].pip.get_yeas_of_proposal(proposalinfo.get('ProposalID')) == 1

    @pytest.mark.P2
    def test_VS_EP_014(self, new_genesis_env, client_con_list_obj):
        genesis = from_dict(data_class=Genesis, data=new_genesis_env.genesis_config)
        genesis.economicModel.gov.paramProposalVoteDurationSeconds = 0
        genesis.economicModel.gov.paramProposalSupportRate = 0.5
        genesis.economicModel.gov.paramProposalVoteRate = 0.5
        new_genesis_env.set_genesis(genesis.to_dict())
        new_genesis_env.deploy_all()
        submitppandvote(client_con_list_obj[:2], 2, 2)
        pip_obj = client_con_list_obj[0].pip
        proposalinfo = pip_obj.get_effect_proposal_info_of_vote(pip_obj.cfg.param_proposal)
        log.info('Get param proposal information {}'.format(proposalinfo))
        wait_block_number(pip_obj.node, 80)
        report_information = mock_duplicate_sign(2, pip_obj.node.nodekey, pip_obj.node.blsprikey, 70)
        log.info("Report information: {}".format(report_information))
        address, _ = pip_obj.economic.account.generate_account(pip_obj.node.web3, 10**18 * 1000)
        result = client_con_list_obj[0].duplicatesign.reportDuplicateSign(2, report_information, address)
        assert_code(result, 0)
        wait_block_number(pip_obj.node, proposalinfo.get('EndVotingBlock'))
        assert client_con_list_obj[1].pip.get_nays_of_proposal(proposalinfo.get('ProposalID')) == 1

    @pytest.mark.P2
    def test_VS_EP_015(self, new_genesis_env, client_con_list_obj):
        genesis = from_dict(data_class=Genesis, data=new_genesis_env.genesis_config)
        genesis.economicModel.gov.paramProposalVoteDurationSeconds = 0
        genesis.economicModel.gov.paramProposalSupportRate = 0.5
        genesis.economicModel.gov.paramProposalVoteRate = 0.5
        new_genesis_env.set_genesis(genesis.to_dict())
        new_genesis_env.deploy_all()
        submitppandvote(client_con_list_obj[:2], 3, 3)
        pip_obj = client_con_list_obj[0].pip
        proposalinfo = pip_obj.get_effect_proposal_info_of_vote(pip_obj.cfg.param_proposal)
        log.info('Get param proposal information {}'.format(proposalinfo))
        wait_block_number(pip_obj.node, 80)
        report_information = mock_duplicate_sign(2, pip_obj.node.nodekey, pip_obj.node.blsprikey, 70)
        log.info("Report information: {}".format(report_information))
        address, _ = pip_obj.economic.account.generate_account(pip_obj.node.web3, 10**18 * 1000)
        result = client_con_list_obj[0].duplicatesign.reportDuplicateSign(2, report_information, address)
        assert_code(result, 0)
        wait_block_number(pip_obj.node, proposalinfo.get('EndVotingBlock'))
        assert client_con_list_obj[1].pip.get_abstentions_of_proposal(proposalinfo.get('ProposalID')) == 1

    @pytest.mark.P2
    def test_VS_EP_016(self, new_genesis_env, client_con_list_obj):
        genesis = from_dict(data_class=Genesis, data=new_genesis_env.genesis_config)
        genesis.economicModel.gov.paramProposalVoteDurationSeconds = 0
        genesis.economicModel.gov.paramProposalSupportRate = 1
        genesis.economicModel.gov.paramProposalVoteRate = 0.249
        new_genesis_env.set_genesis(genesis.to_dict())
        new_genesis_env.deploy_all()
        submitppandvote(client_con_list_obj[:2], 1, 1)
        pip_obj = client_con_list_obj[0].pip
        proposalinfo = pip_obj.get_effect_proposal_info_of_vote(pip_obj.cfg.param_proposal)
        log.info('Get param proposal information {}'.format(proposalinfo))
        assert pip_obj.get_accuverifiers_count(proposalinfo.get('ProposalID')) == [4, 2, 0, 0]
        log.info('Stop the node {}'.format(pip_obj.node.node_id))
        pip_obj.node.stop()
        wait_block_number(client_con_list_obj[1].node, proposalinfo.get('EndVotingBlock'))
        assert client_con_list_obj[1].pip.get_yeas_of_proposal(proposalinfo.get('ProposalID')) == 1
        assert client_con_list_obj[1].pip.get_status_of_proposal(proposalinfo.get('ProposalID')) == 2

    @pytest.mark.P2
    def test_VS_EP_017(self, new_genesis_env, client_con_list_obj):
        genesis = from_dict(data_class=Genesis, data=new_genesis_env.genesis_config)
        genesis.economicModel.gov.paramProposalVoteDurationSeconds = 0
        genesis.economicModel.gov.paramProposalSupportRate = 1
        genesis.economicModel.gov.paramProposalVoteRate = 0.249
        new_genesis_env.set_genesis(genesis.to_dict())
        new_genesis_env.deploy_all()
        submitppandvote(client_con_list_obj[:2], 2, 1)
        pip_obj = client_con_list_obj[0].pip
        proposalinfo = pip_obj.get_effect_proposal_info_of_vote(pip_obj.cfg.param_proposal)
        log.info('Get param proposal information {}'.format(proposalinfo))
        assert pip_obj.get_accuverifiers_count(proposalinfo.get('ProposalID')) == [4, 1, 1, 0]
        log.info('Stop the node {}'.format(pip_obj.node.node_id))
        pip_obj.node.stop()
        wait_block_number(client_con_list_obj[1].node, proposalinfo.get('EndVotingBlock'))
        assert client_con_list_obj[1].pip.get_nays_of_proposal(proposalinfo.get('ProposalID')) == 0
        assert client_con_list_obj[1].pip.get_yeas_of_proposal(proposalinfo.get('ProposalID')) == 1
        assert client_con_list_obj[1].pip.get_status_of_proposal(proposalinfo.get('ProposalID')) == 2

    @pytest.mark.P2
    def test_VS_EP_018(self, new_genesis_env, client_con_list_obj):
        genesis = from_dict(data_class=Genesis, data=new_genesis_env.genesis_config)
        genesis.economicModel.gov.paramProposalVoteDurationSeconds = 0
        genesis.economicModel.gov.paramProposalSupportRate = 1
        genesis.economicModel.gov.paramProposalVoteRate = 0.249
        new_genesis_env.set_genesis(genesis.to_dict())
        new_genesis_env.deploy_all()
        submitppandvote(client_con_list_obj[:2], 3, 1)
        pip_obj = client_con_list_obj[0].pip
        proposalinfo = pip_obj.get_effect_proposal_info_of_vote(pip_obj.cfg.param_proposal)
        log.info('Get param proposal information {}'.format(proposalinfo))
        assert pip_obj.get_accuverifiers_count(proposalinfo.get('ProposalID')) == [4, 1, 0, 1]
        log.info('Stop the node {}'.format(pip_obj.node.node_id))
        pip_obj.node.stop()
        wait_block_number(client_con_list_obj[1].node, proposalinfo.get('EndVotingBlock'))
        assert client_con_list_obj[1].pip.get_abstentions_of_proposal(proposalinfo.get('ProposalID')) == 0
        assert client_con_list_obj[1].pip.get_yeas_of_proposal(proposalinfo.get('ProposalID')) == 1
        assert client_con_list_obj[1].pip.get_status_of_proposal(proposalinfo.get('ProposalID')) == 2

    @pytest.mark.P2
    def test_VS_EP_019(self, new_genesis_env, client_con_list_obj):
        genesis = from_dict(data_class=Genesis, data=new_genesis_env.genesis_config)
        genesis.economicModel.gov.paramProposalVoteDurationSeconds = 0
        genesis.economicModel.gov.paramProposalSupportRate = 0.99
        genesis.economicModel.gov.paramProposalVoteRate = 0.25
        new_genesis_env.set_genesis(genesis.to_dict())
        new_genesis_env.deploy_all()
        submitppandvote(client_con_list_obj[:3], 1, 2, 3)
        pip_obj = client_con_list_obj[0].pip
        proposalinfo = pip_obj.get_effect_proposal_info_of_vote(pip_obj.cfg.param_proposal)
        log.info('Get param proposal information {}'.format(proposalinfo))
        assert pip_obj.get_accuverifiers_count(proposalinfo.get('ProposalID')) == [4, 1, 1, 1]
        wait_block_number(client_con_list_obj[1].node, proposalinfo.get('EndVotingBlock'))
        assert client_con_list_obj[1].pip.get_abstentions_of_proposal(proposalinfo.get('ProposalID')) == 1
        assert client_con_list_obj[1].pip.get_yeas_of_proposal(proposalinfo.get('ProposalID')) == 1
        assert client_con_list_obj[1].pip.get_nays_of_proposal(proposalinfo.get('ProposalID')) == 1
        assert client_con_list_obj[1].pip.get_status_of_proposal(proposalinfo.get('ProposalID')) == 3
        log.info('Stop the node {}'.format(pip_obj.node.node_id))
        pip_obj.node.stop()
        wait_block_number(client_con_list_obj[1].node, proposalinfo.get('EndVotingBlock') +
                          client_con_list_obj[1].economic.consensus_size)
        assert client_con_list_obj[1].pip.get_abstentions_of_proposal(proposalinfo.get('ProposalID')) == 1
        assert client_con_list_obj[1].pip.get_yeas_of_proposal(proposalinfo.get('ProposalID')) == 1
        assert client_con_list_obj[1].pip.get_nays_of_proposal(proposalinfo.get('ProposalID')) == 1
        assert client_con_list_obj[1].pip.get_status_of_proposal(proposalinfo.get('ProposalID')) == 3
