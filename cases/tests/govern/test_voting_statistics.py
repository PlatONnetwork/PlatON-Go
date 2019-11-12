import pytest
from common.log import log
from common.key import mock_duplicate_sign
import time
from tests.lib.utils import assert_code, wait_block_number, upload_platon
from tests.lib.client import Client
from tests.lib import Genesis
from dacite import from_dict
from tests.govern.conftest import param_proposal_vote, version_proposal_vote

def submitvpandvote(client_list_obj, *args):
    pip_obj = client_list_obj[0].pip
    result = pip_obj.submitVersion(pip_obj.node.node_id, str(time.time()), pip_obj.cfg.version5, 2, pip_obj.node.staking_address,
                                   transaction_cfg=pip_obj.cfg.transaction_cfg)
    log.info('Submit version proposal result : {}'.format(result))
    assert_code(result, 0)
    proposalinfo = pip_obj.get_effect_proposal_info_of_vote()
    log.info('Version proposal info {}'.format(proposalinfo))
    for index in range(len(client_list_obj)):
        pip_obj = client_list_obj[index].pip
        log.info('{}'.format(args[index]))
        result = version_proposal_vote(pip_obj, vote_option=pip_obj.cfg.vote_option_yeas)
        assert_code(result, 0)

def submitppandvote(client_list_obj, *args):
    pip_obj = client_list_obj[0].pip
    result = pip_obj.submitParam(pip_obj.node.node_id, str(time.time()), 'Slashing', 'SlashBlocksReward', '83',
                                 pip_obj.node.staking_address, transaction_cfg=pip_obj.cfg.transaction_cfg)
    log.info('Submit param proposal result : {}'.format(result))
    assert_code(result, 0)
    proposalinfo = pip_obj.get_effect_proposal_info_of_vote(pip_obj.cfg.param_proposal)
    log.info('Param proposal info {}'.format(proposalinfo))
    for index in range(len(client_list_obj)):
        pip_obj = client_list_obj[index].pip
        log.info('{}'.format(args[index]))
        result = param_proposal_vote(pip_obj, vote_option=args[index])
        assert_code(result, 0)

def submitcppandvote(client_list_obj, *args):
    pip_obj = client_list_obj[0].pip
    result = pip_obj.submitParam(pip_obj.node.node_id, str(time.time()), 'Slashing', 'SlashBlocksReward', '83',
                                 pip_obj.node.staking_address, transaction_cfg=pip_obj.cfg.transaction_cfg)
    log.info('Submit param proposal result : {}'.format(result))
    assert_code(result, 0)
    proposalinfo_param = pip_obj.get_effect_proposal_info_of_vote(pip_obj.cfg.param_proposal)
    log.info('Param proposal info {}'.format(proposalinfo_param))

    result = pip_obj.submitCancel(pip_obj.node.node_id, str(time.time()), 2, proposalinfo_param.get('ProposalID'),
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

class TestVotingStatistics():
    def createstaking(self, obj):
        if isinstance(obj, Client):
            obj = []
            obj.append(obj)
        for client_obj in obj:
            address, _ = client_obj.economic.account.generate_account(client_obj.node.web3, 10 ** 18 * 10000000)
            result = client_obj.staking.create_staking(0, address, address, amount=10 ** 18 * 2000000,
                                                       transaction_cfg=client_obj.pip.cfg.transaction_cfg)
            log.info('Node {} staking result : {}'.format(client_obj.node.node_id, result))
            assert_code(result, 0)

    def test_VS_EP_004(self, new_genesis_env, client_con_list_obj, client_noc_list_obj):
        genesis = from_dict(data_class=Genesis, data=new_genesis_env.genesis_config)
        genesis.economicModel.gov.paramProposalVoteDurationSeconds = 0
        new_genesis_env.set_genesis(genesis.to_dict())
        new_genesis_env.deploy_all()
        submitppandvote(client_con_list_obj[0:-1], 1, 2, 3)
        proposalinfo = client_con_list_obj[0].pip.get_effect_proposal_info_of_vote(client_con_list_obj[0].pip.cfg.param_proposal)
        log.info('Param proposal information {}'.format(proposalinfo))
        self.createstaking(client_noc_list_obj[:3])
        # client_con_list_obj[0].economic.wait_settlement_blocknum(client_con_list_obj[0].node)
        result = client_con_list_obj[0].pip.get_accuverifiers_count(proposalinfo.get('ProposalID'))
        log.info('Get proposal vote infomation {}'.format(result))
        assert result == [4, 1, 1, 1]

    def test_VS_EP_005(self, new_genesis_env, client_con_list_obj, client_noc_list_obj):
        genesis = from_dict(data_class=Genesis, data=new_genesis_env.genesis_config)
        genesis.economicModel.gov.paramProposalVoteDurationSeconds = 160
        new_genesis_env.set_genesis(genesis.to_dict())
        new_genesis_env.deploy_all()
        pip_obj = client_con_list_obj[0].pip
        submitppandvote(client_con_list_obj[:2], 1, 2)
        proposalinfo = pip_obj.get_effect_proposal_info_of_vote(pip_obj.cfg.param_proposal)
        log.info('Param proposal info {}'.format(proposalinfo))
        self.createstaking(client_noc_list_obj[:2])
        pip_obj.economic.wait_settlement_blocknum(pip_obj.node)
        result = pip_obj.get_accuverifiers_count(proposalinfo.get('ProposalID'))
        log.info('Get proposal vote infomation {}'.format(result))
        assert result == [6, 1, 1, 0]

        result = param_proposal_vote(client_con_list_obj[3].pip, pip_obj.cfg.vote_option_Abstentions)
        assert_code(result, 0)
        self.createstaking(client_noc_list_obj[2])
        wait_block_number(pip_obj.node, proposalinfo.get('EndVotingBlock'))

        result = pip_obj.get_accuverifiers_count(proposalinfo.get('ProposalID'))
        log.info('Get proposal vote infomation {}'.format(result))
        assert result == [6, 1, 1, 0]

    def test_VS_EP_006(self, new_genesis_env, client_con_list_obj, client_noc_list_obj):
        genesis = from_dict(data_class=Genesis, data=new_genesis_env.genesis_config)
        genesis.economicModel.gov.paramProposalVoteDurationSeconds = 320
        new_genesis_env.set_genesis(genesis.to_dict())
        new_genesis_env.deploy_all()
        submitppandvote(client_con_list_obj[:1], 1)
        proposalinfo = client_con_list_obj[0].pip.get_effect_proposal_info_of_vote(client_con_list_obj[0].pip.cfg.param_proposal)
        log.info('Param proposal info {}'.format(proposalinfo))
        self.createstaking(client_noc_list_obj[:1])
        client_con_list_obj[0].pip.economic.wait_settlement_blocknum(client_con_list_obj[0].pip.node)
        result = client_con_list_obj[0].pip.get_accuverifiers_count(proposalinfo.get('ProposalID'))
        log.info('Get proposal vote infomation {}'.format(result))
        assert result == [5, 1, 0, 0]

        result = param_proposal_vote(client_con_list_obj[1].pip, client_con_list_obj[0].pip.cfg.vote_option_nays)
        assert_code(result, 0)
        self.createstaking(client_noc_list_obj[1])
        client_con_list_obj[0].pip.economic.wait_settlement_blocknum(client_con_list_obj[0].pip.node)
        result = client_con_list_obj[0].pip.get_accuverifiers_count(proposalinfo.get('ProposalID'))
        log.info('Get proposal vote infomation {}'.format(result))
        assert result == [6, 1, 1, 0]

        result = param_proposal_vote(client_con_list_obj[2].pip, client_con_list_obj[0].pip.cfg.vote_option_nays)
        assert_code(result, 0)
        self.createstaking(client_noc_list_obj[2])
        wait_block_number(client_con_list_obj[0].pip.node, proposalinfo.get('EndVotingBlock'))

        result = client_con_list_obj[0].pip.get_accuverifiers_count(proposalinfo.get('ProposalID'))
        log.info('Get proposal vote infomation {}'.format(result))
        assert result == [6, 1, 1, 0]

    def test_VS_EP_007(self, new_genesis_env, client_con_list_obj):
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
        wait_block_number(client_con_list_obj[0].node, proposalinfo.get('EndVotingBlock'))
        assert client_con_list_obj[0].pip.get_yeas_of_proposal(proposalinfo.get('ProposalID')) == 2
        assert client_con_list_obj[0].pip.get_status_of_proposal(proposalinfo.get('ProposalID')) == 2

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

    def test_VS_EP_016(self, new_genesis_env, client_con_list_obj):
        genesis = from_dict(data_class=Genesis, data=new_genesis_env.genesis_config)
        genesis.economicModel.gov.paramProposalVoteDurationSeconds = 0
        genesis.economicModel.gov.paramProposalSupportRate = 0.999
        genesis.economicModel.gov.paramProposalVoteRate = 0.25
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

    def test_VS_EP_017(self, new_genesis_env, client_con_list_obj):
        genesis = from_dict(data_class=Genesis, data=new_genesis_env.genesis_config)
        genesis.economicModel.gov.paramProposalVoteDurationSeconds = 0
        genesis.economicModel.gov.paramProposalSupportRate = 0.99
        genesis.economicModel.gov.paramProposalVoteRate = 0.25
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

    def test_VS_EP_018(self, new_genesis_env, client_con_list_obj):
        genesis = from_dict(data_class=Genesis, data=new_genesis_env.genesis_config)
        genesis.economicModel.gov.paramProposalVoteDurationSeconds = 0
        genesis.economicModel.gov.paramProposalSupportRate = 0.99
        genesis.economicModel.gov.paramProposalVoteRate = 0.25
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

class TestSupportRateVoteRatePP():
    def test_UP_PA_001(self, new_genesis_env, client_con_list_obj):
        genesis = from_dict(data_class=Genesis, data=new_genesis_env.genesis_config)
        genesis.economicModel.gov.paramProposalVoteDurationSeconds = 0
        genesis.economicModel.gov.paramProposalSupportRate = 0.332
        genesis.economicModel.gov.paramProposalVoteRate = 0.751
        new_genesis_env.set_genesis(genesis.to_dict())
        new_genesis_env.deploy_all()
        submitppandvote(client_con_list_obj[:3], 1, 2, 3)
        pip_obj = client_con_list_obj[0].pip
        proposalinfo = pip_obj.get_effect_proposal_info_of_vote(pip_obj.cfg.param_proposal)
        log.info('Get param proposal information {}'.format(proposalinfo))
        assert pip_obj.get_accuverifiers_count(proposalinfo.get('ProposalID')) == [4, 1, 1, 1]
        wait_block_number(client_con_list_obj[1].node, proposalinfo.get('EndVotingBlock'))
        assert_code(pip_obj.get_yeas_of_proposal(proposalinfo.get('ProposalID')), 1)
        assert_code(pip_obj.get_nays_of_proposal(proposalinfo.get('ProposalID')), 1)
        assert_code(pip_obj.get_abstentions_of_proposal(proposalinfo.get('ProposalID')), 1)
        assert_code(pip_obj.get_accu_verifiers_of_proposal(proposalinfo.get('ProposalID')), len(client_con_list_obj))
        assert_code(pip_obj.get_status_of_proposal(proposalinfo.get('ProposalID')), 3)

    def test_UP_PA_002(self, new_genesis_env, client_con_list_obj):
        genesis = from_dict(data_class=Genesis, data=new_genesis_env.genesis_config)
        genesis.economicModel.gov.paramProposalVoteDurationSeconds = 0
        genesis.economicModel.gov.paramProposalSupportRate = 0.334
        genesis.economicModel.gov.paramProposalVoteRate = 0.749
        new_genesis_env.set_genesis(genesis.to_dict())
        new_genesis_env.deploy_all()
        submitppandvote(client_con_list_obj[:3], 1, 2, 3)
        pip_obj = client_con_list_obj[0].pip
        proposalinfo = pip_obj.get_effect_proposal_info_of_vote(pip_obj.cfg.param_proposal)
        log.info('Get param proposal information {}'.format(proposalinfo))
        assert pip_obj.get_accuverifiers_count(proposalinfo.get('ProposalID')) == [4, 1, 1, 1]
        wait_block_number(client_con_list_obj[1].node, proposalinfo.get('EndVotingBlock'))
        assert_code(pip_obj.get_yeas_of_proposal(proposalinfo.get('ProposalID')), 1)
        assert_code(pip_obj.get_nays_of_proposal(proposalinfo.get('ProposalID')), 1)
        assert_code(pip_obj.get_abstentions_of_proposal(proposalinfo.get('ProposalID')), 1)
        assert_code(pip_obj.get_accu_verifiers_of_proposal(proposalinfo.get('ProposalID')), len(client_con_list_obj))
        assert_code(pip_obj.get_status_of_proposal(proposalinfo.get('ProposalID')), 3)

    def test_UP_PA_003(self, new_genesis_env, client_con_list_obj):
        genesis = from_dict(data_class=Genesis, data=new_genesis_env.genesis_config)
        genesis.economicModel.gov.paramProposalVoteDurationSeconds = 0
        genesis.economicModel.gov.paramProposalSupportRate = 0.333
        genesis.economicModel.gov.paramProposalVoteRate = 0.751
        new_genesis_env.set_genesis(genesis.to_dict())
        new_genesis_env.deploy_all()
        submitppandvote(client_con_list_obj[:3], 1, 2, 3)
        pip_obj = client_con_list_obj[0].pip
        proposalinfo = pip_obj.get_effect_proposal_info_of_vote(pip_obj.cfg.param_proposal)
        log.info('Get param proposal information {}'.format(proposalinfo))
        assert pip_obj.get_accuverifiers_count(proposalinfo.get('ProposalID')) == [4, 1, 1, 1]
        wait_block_number(client_con_list_obj[1].node, proposalinfo.get('EndVotingBlock'))
        assert_code(pip_obj.get_yeas_of_proposal(proposalinfo.get('ProposalID')), 1)
        assert_code(pip_obj.get_nays_of_proposal(proposalinfo.get('ProposalID')), 1)
        assert_code(pip_obj.get_abstentions_of_proposal(proposalinfo.get('ProposalID')), 1)
        assert_code(pip_obj.get_accu_verifiers_of_proposal(proposalinfo.get('ProposalID')), len(client_con_list_obj))
        assert_code(pip_obj.get_status_of_proposal(proposalinfo.get('ProposalID')), 3)

    def test_UP_PA_004(self, new_genesis_env, client_con_list_obj):
        genesis = from_dict(data_class=Genesis, data=new_genesis_env.genesis_config)
        genesis.economicModel.gov.paramProposalVoteDurationSeconds = 0
        genesis.economicModel.gov.paramProposalSupportRate = 0.334
        genesis.economicModel.gov.paramProposalVoteRate = 0.75
        new_genesis_env.set_genesis(genesis.to_dict())
        new_genesis_env.deploy_all()
        submitppandvote(client_con_list_obj[:3], 1, 2, 3)
        pip_obj = client_con_list_obj[0].pip
        proposalinfo = pip_obj.get_effect_proposal_info_of_vote(pip_obj.cfg.param_proposal)
        log.info('Get param proposal information {}'.format(proposalinfo))
        assert pip_obj.get_accuverifiers_count(proposalinfo.get('ProposalID')) == [4, 1, 1, 1]
        wait_block_number(client_con_list_obj[1].node, proposalinfo.get('EndVotingBlock'))
        assert_code(pip_obj.get_yeas_of_proposal(proposalinfo.get('ProposalID')), 1)
        assert_code(pip_obj.get_nays_of_proposal(proposalinfo.get('ProposalID')), 1)
        assert_code(pip_obj.get_abstentions_of_proposal(proposalinfo.get('ProposalID')), 1)
        assert_code(pip_obj.get_accu_verifiers_of_proposal(proposalinfo.get('ProposalID')), len(client_con_list_obj))
        assert_code(pip_obj.get_status_of_proposal(proposalinfo.get('ProposalID')), 3)

    def test_UP_PA_005(self, new_genesis_env, client_con_list_obj):
        genesis = from_dict(data_class=Genesis, data=new_genesis_env.genesis_config)
        genesis.economicModel.gov.paramProposalVoteDurationSeconds = 0
        genesis.economicModel.gov.paramProposalSupportRate = 0.332
        genesis.economicModel.gov.paramProposalVoteRate = 0.749
        new_genesis_env.set_genesis(genesis.to_dict())
        new_genesis_env.deploy_all()
        submitppandvote(client_con_list_obj[:3], 1, 2, 3)
        pip_obj = client_con_list_obj[0].pip
        proposalinfo = pip_obj.get_effect_proposal_info_of_vote(pip_obj.cfg.param_proposal)
        log.info('Get param proposal information {}'.format(proposalinfo))
        assert pip_obj.get_accuverifiers_count(proposalinfo.get('ProposalID')) == [4, 1, 1, 1]
        wait_block_number(client_con_list_obj[1].node, proposalinfo.get('EndVotingBlock'))
        assert_code(pip_obj.get_yeas_of_proposal(proposalinfo.get('ProposalID')), 1)
        assert_code(pip_obj.get_nays_of_proposal(proposalinfo.get('ProposalID')), 1)
        assert_code(pip_obj.get_abstentions_of_proposal(proposalinfo.get('ProposalID')), 1)
        assert_code(pip_obj.get_accu_verifiers_of_proposal(proposalinfo.get('ProposalID')), len(client_con_list_obj))
        assert_code(pip_obj.get_status_of_proposal(proposalinfo.get('ProposalID')), 2)

    def test_UP_PA_006(self, new_genesis_env, client_con_list_obj):
        genesis = from_dict(data_class=Genesis, data=new_genesis_env.genesis_config)
        genesis.economicModel.gov.paramProposalVoteDurationSeconds = 0
        genesis.economicModel.gov.paramProposalSupportRate = 0.333
        genesis.economicModel.gov.paramProposalVoteRate = 0.749
        new_genesis_env.set_genesis(genesis.to_dict())
        new_genesis_env.deploy_all()
        submitppandvote(client_con_list_obj[:3], 1, 2, 3)
        pip_obj = client_con_list_obj[0].pip
        proposalinfo = pip_obj.get_effect_proposal_info_of_vote(pip_obj.cfg.param_proposal)
        log.info('Get param proposal information {}'.format(proposalinfo))
        assert pip_obj.get_accuverifiers_count(proposalinfo.get('ProposalID')) == [4, 1, 1, 1]
        wait_block_number(client_con_list_obj[1].node, proposalinfo.get('EndVotingBlock'))
        assert_code(pip_obj.get_yeas_of_proposal(proposalinfo.get('ProposalID')), 1)
        assert_code(pip_obj.get_nays_of_proposal(proposalinfo.get('ProposalID')), 1)
        assert_code(pip_obj.get_abstentions_of_proposal(proposalinfo.get('ProposalID')), 1)
        assert_code(pip_obj.get_accu_verifiers_of_proposal(proposalinfo.get('ProposalID')), len(client_con_list_obj))
        assert_code(pip_obj.get_status_of_proposal(proposalinfo.get('ProposalID')), 2)

    def test_UP_PA_007(self, new_genesis_env, client_con_list_obj):
        genesis = from_dict(data_class=Genesis, data=new_genesis_env.genesis_config)
        genesis.economicModel.gov.paramProposalVoteDurationSeconds = 0
        genesis.economicModel.gov.paramProposalSupportRate = 0.332
        genesis.economicModel.gov.paramProposalVoteRate = 0.75
        new_genesis_env.set_genesis(genesis.to_dict())
        new_genesis_env.deploy_all()
        submitppandvote(client_con_list_obj[:3], 1, 2, 3)
        pip_obj = client_con_list_obj[0].pip
        proposalinfo = pip_obj.get_effect_proposal_info_of_vote(pip_obj.cfg.param_proposal)
        log.info('Get param proposal information {}'.format(proposalinfo))
        assert pip_obj.get_accuverifiers_count(proposalinfo.get('ProposalID')) == [4, 1, 1, 1]
        wait_block_number(client_con_list_obj[1].node, proposalinfo.get('EndVotingBlock'))
        assert_code(pip_obj.get_yeas_of_proposal(proposalinfo.get('ProposalID')), 1)
        assert_code(pip_obj.get_nays_of_proposal(proposalinfo.get('ProposalID')), 1)
        assert_code(pip_obj.get_abstentions_of_proposal(proposalinfo.get('ProposalID')), 1)
        assert_code(pip_obj.get_accu_verifiers_of_proposal(proposalinfo.get('ProposalID')), len(client_con_list_obj))
        assert_code(pip_obj.get_status_of_proposal(proposalinfo.get('ProposalID')), 3)

class TestSupportRateVoteRateCPP():
    def test_UC_CP_001(self, new_genesis_env, client_con_list_obj):
        genesis = from_dict(data_class=Genesis, data=new_genesis_env.genesis_config)
        genesis.economicModel.gov.cancelProposalSupportRate = 0.332
        genesis.economicModel.gov.cancelProposalVoteRate = 0.751
        new_genesis_env.set_genesis(genesis.to_dict())
        new_genesis_env.deploy_all()
        submitcppandvote(client_con_list_obj[:3], 1, 2, 3)
        pip_obj = client_con_list_obj[0].pip
        proposalinfo = pip_obj.get_effect_proposal_info_of_vote(pip_obj.cfg.cancel_proposal)
        log.info('Get cancel proposal information {}'.format(proposalinfo))

        assert pip_obj.get_accuverifiers_count(proposalinfo.get('ProposalID')) == [4, 1, 1, 1]
        wait_block_number(client_con_list_obj[1].node, proposalinfo.get('EndVotingBlock'))
        assert_code(pip_obj.get_yeas_of_proposal(proposalinfo.get('ProposalID')), 1)
        assert_code(pip_obj.get_nays_of_proposal(proposalinfo.get('ProposalID')), 1)
        assert_code(pip_obj.get_abstentions_of_proposal(proposalinfo.get('ProposalID')), 1)
        assert_code(pip_obj.get_accu_verifiers_of_proposal(proposalinfo.get('ProposalID')), len(client_con_list_obj))
        assert_code(pip_obj.get_status_of_proposal(proposalinfo.get('ProposalID')), 3)

    def test_UC_CP_002(self, new_genesis_env, client_con_list_obj):
        genesis = from_dict(data_class=Genesis, data=new_genesis_env.genesis_config)
        genesis.economicModel.gov.cancelProposalSupportRate = 0.334
        genesis.economicModel.gov.cancelProposalVoteRate = 0.749
        new_genesis_env.set_genesis(genesis.to_dict())
        new_genesis_env.deploy_all()
        submitcppandvote(client_con_list_obj[:3], 1, 2, 3)
        pip_obj = client_con_list_obj[0].pip
        proposalinfo = pip_obj.get_effect_proposal_info_of_vote(pip_obj.cfg.cancel_proposal)
        log.info('Get cancel proposal information {}'.format(proposalinfo))

        assert pip_obj.get_accuverifiers_count(proposalinfo.get('ProposalID')) == [4, 1, 1, 1]
        wait_block_number(client_con_list_obj[1].node, proposalinfo.get('EndVotingBlock'))
        assert_code(pip_obj.get_yeas_of_proposal(proposalinfo.get('ProposalID')), 1)
        assert_code(pip_obj.get_nays_of_proposal(proposalinfo.get('ProposalID')), 1)
        assert_code(pip_obj.get_abstentions_of_proposal(proposalinfo.get('ProposalID')), 1)
        assert_code(pip_obj.get_accu_verifiers_of_proposal(proposalinfo.get('ProposalID')), len(client_con_list_obj))
        assert_code(pip_obj.get_status_of_proposal(proposalinfo.get('ProposalID')), 3)

    def test_UC_CP_003(self, new_genesis_env, client_con_list_obj):
        genesis = from_dict(data_class=Genesis, data=new_genesis_env.genesis_config)
        genesis.economicModel.gov.cancelProposalSupportRate = 0.333
        genesis.economicModel.gov.cancelProposalVoteRate = 0.751
        new_genesis_env.set_genesis(genesis.to_dict())
        new_genesis_env.deploy_all()
        submitcppandvote(client_con_list_obj[:3], 1, 2, 3)
        pip_obj = client_con_list_obj[0].pip
        proposalinfo = pip_obj.get_effect_proposal_info_of_vote(pip_obj.cfg.cancel_proposal)
        log.info('Get cancel proposal information {}'.format(proposalinfo))

        assert pip_obj.get_accuverifiers_count(proposalinfo.get('ProposalID')) == [4, 1, 1, 1]
        wait_block_number(client_con_list_obj[1].node, proposalinfo.get('EndVotingBlock'))
        assert_code(pip_obj.get_yeas_of_proposal(proposalinfo.get('ProposalID')), 1)
        assert_code(pip_obj.get_nays_of_proposal(proposalinfo.get('ProposalID')), 1)
        assert_code(pip_obj.get_abstentions_of_proposal(proposalinfo.get('ProposalID')), 1)
        assert_code(pip_obj.get_accu_verifiers_of_proposal(proposalinfo.get('ProposalID')), len(client_con_list_obj))
        assert_code(pip_obj.get_status_of_proposal(proposalinfo.get('ProposalID')), 3)

    def test_UC_CP_004(self, new_genesis_env, client_con_list_obj):
        genesis = from_dict(data_class=Genesis, data=new_genesis_env.genesis_config)
        genesis.economicModel.gov.cancelProposalSupportRate = 0.334
        genesis.economicModel.gov.cancelProposalVoteRate = 0.75
        new_genesis_env.set_genesis(genesis.to_dict())
        new_genesis_env.deploy_all()
        submitcppandvote(client_con_list_obj[:3], 1, 2, 3)
        pip_obj = client_con_list_obj[0].pip
        proposalinfo = pip_obj.get_effect_proposal_info_of_vote(pip_obj.cfg.cancel_proposal)
        log.info('Get cancel proposal information {}'.format(proposalinfo))

        assert pip_obj.get_accuverifiers_count(proposalinfo.get('ProposalID')) == [4, 1, 1, 1]
        wait_block_number(client_con_list_obj[1].node, proposalinfo.get('EndVotingBlock'))
        assert_code(pip_obj.get_yeas_of_proposal(proposalinfo.get('ProposalID')), 1)
        assert_code(pip_obj.get_nays_of_proposal(proposalinfo.get('ProposalID')), 1)
        assert_code(pip_obj.get_abstentions_of_proposal(proposalinfo.get('ProposalID')), 1)
        assert_code(pip_obj.get_accu_verifiers_of_proposal(proposalinfo.get('ProposalID')), len(client_con_list_obj))
        assert_code(pip_obj.get_status_of_proposal(proposalinfo.get('ProposalID')), 3)

    def test_UC_CP_005(self, new_genesis_env, client_con_list_obj):
        genesis = from_dict(data_class=Genesis, data=new_genesis_env.genesis_config)
        genesis.economicModel.gov.cancelProposalSupportRate = 0.332
        genesis.economicModel.gov.cancelProposalVoteRate = 0.749
        new_genesis_env.set_genesis(genesis.to_dict())
        new_genesis_env.deploy_all()
        submitcppandvote(client_con_list_obj[:3], 1, 2, 3)
        pip_obj = client_con_list_obj[0].pip
        proposalinfo = pip_obj.get_effect_proposal_info_of_vote(pip_obj.cfg.cancel_proposal)
        log.info('Get cancel proposal information {}'.format(proposalinfo))

        assert pip_obj.get_accuverifiers_count(proposalinfo.get('ProposalID')) == [4, 1, 1, 1]
        wait_block_number(client_con_list_obj[1].node, proposalinfo.get('EndVotingBlock'))
        assert_code(pip_obj.get_yeas_of_proposal(proposalinfo.get('ProposalID')), 1)
        assert_code(pip_obj.get_nays_of_proposal(proposalinfo.get('ProposalID')), 1)
        assert_code(pip_obj.get_abstentions_of_proposal(proposalinfo.get('ProposalID')), 1)
        assert_code(pip_obj.get_accu_verifiers_of_proposal(proposalinfo.get('ProposalID')), len(client_con_list_obj))
        assert_code(pip_obj.get_status_of_proposal(proposalinfo.get('ProposalID')), 2)

    def test_UC_CP_006(self, new_genesis_env, client_con_list_obj):
        genesis = from_dict(data_class=Genesis, data=new_genesis_env.genesis_config)
        genesis.economicModel.gov.cancelProposalSupportRate = 0.333
        genesis.economicModel.gov.cancelProposalVoteRate = 0.749
        new_genesis_env.set_genesis(genesis.to_dict())
        new_genesis_env.deploy_all()
        submitcppandvote(client_con_list_obj[:3], 1, 2, 3)
        pip_obj = client_con_list_obj[0].pip
        proposalinfo = pip_obj.get_effect_proposal_info_of_vote(pip_obj.cfg.cancel_proposal)
        log.info('Get cancel proposal information {}'.format(proposalinfo))

        assert pip_obj.get_accuverifiers_count(proposalinfo.get('ProposalID')) == [4, 1, 1, 1]
        wait_block_number(client_con_list_obj[1].node, proposalinfo.get('EndVotingBlock'))
        assert_code(pip_obj.get_yeas_of_proposal(proposalinfo.get('ProposalID')), 1)
        assert_code(pip_obj.get_nays_of_proposal(proposalinfo.get('ProposalID')), 1)
        assert_code(pip_obj.get_abstentions_of_proposal(proposalinfo.get('ProposalID')), 1)
        assert_code(pip_obj.get_accu_verifiers_of_proposal(proposalinfo.get('ProposalID')), len(client_con_list_obj))
        assert_code(pip_obj.get_status_of_proposal(proposalinfo.get('ProposalID')), 2)

    def test_UC_CP_007(self, new_genesis_env, client_con_list_obj):
        genesis = from_dict(data_class=Genesis, data=new_genesis_env.genesis_config)
        genesis.economicModel.gov.cancelProposalSupportRate = 0.332
        genesis.economicModel.gov.cancelProposalVoteRate = 0.75
        new_genesis_env.set_genesis(genesis.to_dict())
        new_genesis_env.deploy_all()
        submitcppandvote(client_con_list_obj[:3], 1, 2, 3)
        pip_obj = client_con_list_obj[0].pip
        proposalinfo = pip_obj.get_effect_proposal_info_of_vote(pip_obj.cfg.cancel_proposal)
        log.info('Get cancel proposal information {}'.format(proposalinfo))

        assert pip_obj.get_accuverifiers_count(proposalinfo.get('ProposalID')) == [4, 1, 1, 1]
        wait_block_number(client_con_list_obj[1].node, proposalinfo.get('EndVotingBlock'))
        assert_code(pip_obj.get_yeas_of_proposal(proposalinfo.get('ProposalID')), 1)
        assert_code(pip_obj.get_nays_of_proposal(proposalinfo.get('ProposalID')), 1)
        assert_code(pip_obj.get_abstentions_of_proposal(proposalinfo.get('ProposalID')), 1)
        assert_code(pip_obj.get_accu_verifiers_of_proposal(proposalinfo.get('ProposalID')), len(client_con_list_obj))
        assert_code(pip_obj.get_status_of_proposal(proposalinfo.get('ProposalID')), 3)


class TestSupportRateVoteRateCVP():
    def test_UP_CA_001(self, new_genesis_env, client_con_list_obj):
        genesis = from_dict(data_class=Genesis, data=new_genesis_env.genesis_config)
        genesis.economicModel.gov.cancelProposalSupportRate = 0.332
        genesis.economicModel.gov.cancelProposalVoteRate = 0.751
        new_genesis_env.set_genesis(genesis.to_dict())
        new_genesis_env.deploy_all()
        submitcvpandvote(client_con_list_obj[:3], 1, 2, 3)
        pip_obj = client_con_list_obj[0].pip
        proposalinfo = pip_obj.get_effect_proposal_info_of_vote(pip_obj.cfg.cancel_proposal)
        log.info('Get cancel proposal information {}'.format(proposalinfo))

        assert pip_obj.get_accuverifiers_count(proposalinfo.get('ProposalID')) == [4, 1, 1, 1]
        wait_block_number(client_con_list_obj[1].node, proposalinfo.get('EndVotingBlock'))
        assert_code(pip_obj.get_yeas_of_proposal(proposalinfo.get('ProposalID')), 1)
        assert_code(pip_obj.get_nays_of_proposal(proposalinfo.get('ProposalID')), 1)
        assert_code(pip_obj.get_abstentions_of_proposal(proposalinfo.get('ProposalID')), 1)
        assert_code(pip_obj.get_accu_verifiers_of_proposal(proposalinfo.get('ProposalID')), len(client_con_list_obj))
        assert_code(pip_obj.get_status_of_proposal(proposalinfo.get('ProposalID')), 3)

    def test_UP_CA_002(self, new_genesis_env, client_con_list_obj):
        genesis = from_dict(data_class=Genesis, data=new_genesis_env.genesis_config)
        genesis.economicModel.gov.cancelProposalSupportRate = 0.334
        genesis.economicModel.gov.cancelProposalVoteRate = 0.749
        new_genesis_env.set_genesis(genesis.to_dict())
        new_genesis_env.deploy_all()
        submitcvpandvote(client_con_list_obj[:3], 1, 2, 3)
        pip_obj = client_con_list_obj[0].pip
        proposalinfo = pip_obj.get_effect_proposal_info_of_vote(pip_obj.cfg.cancel_proposal)
        log.info('Get cancel proposal information {}'.format(proposalinfo))

        assert pip_obj.get_accuverifiers_count(proposalinfo.get('ProposalID')) == [4, 1, 1, 1]
        wait_block_number(client_con_list_obj[1].node, proposalinfo.get('EndVotingBlock'))
        assert_code(pip_obj.get_yeas_of_proposal(proposalinfo.get('ProposalID')), 1)
        assert_code(pip_obj.get_nays_of_proposal(proposalinfo.get('ProposalID')), 1)
        assert_code(pip_obj.get_abstentions_of_proposal(proposalinfo.get('ProposalID')), 1)
        assert_code(pip_obj.get_accu_verifiers_of_proposal(proposalinfo.get('ProposalID')), len(client_con_list_obj))
        assert_code(pip_obj.get_status_of_proposal(proposalinfo.get('ProposalID')), 3)

    def test_UP_CA_003(self, new_genesis_env, client_con_list_obj):
        genesis = from_dict(data_class=Genesis, data=new_genesis_env.genesis_config)
        genesis.economicModel.gov.cancelProposalSupportRate = 0.333
        genesis.economicModel.gov.cancelProposalVoteRate = 0.751
        new_genesis_env.set_genesis(genesis.to_dict())
        new_genesis_env.deploy_all()
        submitcvpandvote(client_con_list_obj[:3], 1, 2, 3)
        pip_obj = client_con_list_obj[0].pip
        proposalinfo = pip_obj.get_effect_proposal_info_of_vote(pip_obj.cfg.cancel_proposal)
        log.info('Get cancel proposal information {}'.format(proposalinfo))

        assert pip_obj.get_accuverifiers_count(proposalinfo.get('ProposalID')) == [4, 1, 1, 1]
        wait_block_number(client_con_list_obj[1].node, proposalinfo.get('EndVotingBlock'))
        assert_code(pip_obj.get_yeas_of_proposal(proposalinfo.get('ProposalID')), 1)
        assert_code(pip_obj.get_nays_of_proposal(proposalinfo.get('ProposalID')), 1)
        assert_code(pip_obj.get_abstentions_of_proposal(proposalinfo.get('ProposalID')), 1)
        assert_code(pip_obj.get_accu_verifiers_of_proposal(proposalinfo.get('ProposalID')), len(client_con_list_obj))
        assert_code(pip_obj.get_status_of_proposal(proposalinfo.get('ProposalID')), 3)

    def test_UP_CA_004(self, new_genesis_env, client_con_list_obj):
        genesis = from_dict(data_class=Genesis, data=new_genesis_env.genesis_config)
        genesis.economicModel.gov.cancelProposalSupportRate = 0.334
        genesis.economicModel.gov.cancelProposalVoteRate = 0.75
        new_genesis_env.set_genesis(genesis.to_dict())
        new_genesis_env.deploy_all()
        submitcvpandvote(client_con_list_obj[:3], 1, 2, 3)
        pip_obj = client_con_list_obj[0].pip
        proposalinfo = pip_obj.get_effect_proposal_info_of_vote(pip_obj.cfg.cancel_proposal)
        log.info('Get cancel proposal information {}'.format(proposalinfo))

        assert pip_obj.get_accuverifiers_count(proposalinfo.get('ProposalID')) == [4, 1, 1, 1]
        wait_block_number(client_con_list_obj[1].node, proposalinfo.get('EndVotingBlock'))
        assert_code(pip_obj.get_yeas_of_proposal(proposalinfo.get('ProposalID')), 1)
        assert_code(pip_obj.get_nays_of_proposal(proposalinfo.get('ProposalID')), 1)
        assert_code(pip_obj.get_abstentions_of_proposal(proposalinfo.get('ProposalID')), 1)
        assert_code(pip_obj.get_accu_verifiers_of_proposal(proposalinfo.get('ProposalID')), len(client_con_list_obj))
        assert_code(pip_obj.get_status_of_proposal(proposalinfo.get('ProposalID')), 3)

    def test_UP_CA_005(self, new_genesis_env, client_con_list_obj):
        genesis = from_dict(data_class=Genesis, data=new_genesis_env.genesis_config)
        genesis.economicModel.gov.cancelProposalSupportRate = 0.332
        genesis.economicModel.gov.cancelProposalVoteRate = 0.749
        new_genesis_env.set_genesis(genesis.to_dict())
        new_genesis_env.deploy_all()
        submitcvpandvote(client_con_list_obj[:3], 1, 2, 3)
        pip_obj = client_con_list_obj[0].pip
        proposalinfo = pip_obj.get_effect_proposal_info_of_vote(pip_obj.cfg.cancel_proposal)
        log.info('Get cancel proposal information {}'.format(proposalinfo))

        assert pip_obj.get_accuverifiers_count(proposalinfo.get('ProposalID')) == [4, 1, 1, 1]
        wait_block_number(client_con_list_obj[1].node, proposalinfo.get('EndVotingBlock'))
        assert_code(pip_obj.get_yeas_of_proposal(proposalinfo.get('ProposalID')), 1)
        assert_code(pip_obj.get_nays_of_proposal(proposalinfo.get('ProposalID')), 1)
        assert_code(pip_obj.get_abstentions_of_proposal(proposalinfo.get('ProposalID')), 1)
        assert_code(pip_obj.get_accu_verifiers_of_proposal(proposalinfo.get('ProposalID')), len(client_con_list_obj))
        assert_code(pip_obj.get_status_of_proposal(proposalinfo.get('ProposalID')), 2)

    def test_UP_CA_006(self, new_genesis_env, client_con_list_obj):
        genesis = from_dict(data_class=Genesis, data=new_genesis_env.genesis_config)
        genesis.economicModel.gov.cancelProposalSupportRate = 0.332
        genesis.economicModel.gov.cancelProposalVoteRate = 0.749
        new_genesis_env.set_genesis(genesis.to_dict())
        new_genesis_env.deploy_all()
        submitcvpandvote(client_con_list_obj[:3], 1, 2, 3)
        pip_obj = client_con_list_obj[0].pip
        proposalinfo = pip_obj.get_effect_proposal_info_of_vote(pip_obj.cfg.cancel_proposal)
        log.info('Get cancel proposal information {}'.format(proposalinfo))

        assert pip_obj.get_accuverifiers_count(proposalinfo.get('ProposalID')) == [4, 1, 1, 1]
        wait_block_number(client_con_list_obj[1].node, proposalinfo.get('EndVotingBlock'))
        assert_code(pip_obj.get_yeas_of_proposal(proposalinfo.get('ProposalID')), 1)
        assert_code(pip_obj.get_nays_of_proposal(proposalinfo.get('ProposalID')), 1)
        assert_code(pip_obj.get_abstentions_of_proposal(proposalinfo.get('ProposalID')), 1)
        assert_code(pip_obj.get_accu_verifiers_of_proposal(proposalinfo.get('ProposalID')), len(client_con_list_obj))
        assert_code(pip_obj.get_status_of_proposal(proposalinfo.get('ProposalID')), 2)

    def test_UP_CA_007(self, new_genesis_env, client_con_list_obj):
        genesis = from_dict(data_class=Genesis, data=new_genesis_env.genesis_config)
        genesis.economicModel.gov.cancelProposalSupportRate = 0.332
        genesis.economicModel.gov.cancelProposalVoteRate = 0.75
        new_genesis_env.set_genesis(genesis.to_dict())
        new_genesis_env.deploy_all()
        submitcvpandvote(client_con_list_obj[:3], 1, 2, 3)
        pip_obj = client_con_list_obj[0].pip
        proposalinfo = pip_obj.get_effect_proposal_info_of_vote(pip_obj.cfg.cancel_proposal)
        log.info('Get cancel proposal information {}'.format(proposalinfo))

        assert pip_obj.get_accuverifiers_count(proposalinfo.get('ProposalID')) == [4, 1, 1, 1]
        wait_block_number(client_con_list_obj[1].node, proposalinfo.get('EndVotingBlock'))
        assert_code(pip_obj.get_yeas_of_proposal(proposalinfo.get('ProposalID')), 1)
        assert_code(pip_obj.get_nays_of_proposal(proposalinfo.get('ProposalID')), 1)
        assert_code(pip_obj.get_abstentions_of_proposal(proposalinfo.get('ProposalID')), 1)
        assert_code(pip_obj.get_accu_verifiers_of_proposal(proposalinfo.get('ProposalID')), len(client_con_list_obj))
        assert_code(pip_obj.get_status_of_proposal(proposalinfo.get('ProposalID')), 3)


class TestSupportRateVoteRateTP():
    def test_UP_TE_001(self, new_genesis_env, client_con_list_obj):
        genesis = from_dict(data_class=Genesis, data=new_genesis_env.genesis_config)
        genesis.economicModel.gov.textProposalSupportRate = 0.332
        genesis.economicModel.gov.textProposalVoteRate = 0.751
        genesis.economicModel.gov.textProposalVoteDurationSeconds = 40
        new_genesis_env.set_genesis(genesis.to_dict())
        new_genesis_env.deploy_all()
        submittpandvote(client_con_list_obj[:3], 1, 2, 3)
        pip_obj = client_con_list_obj[0].pip
        proposalinfo = pip_obj.get_effect_proposal_info_of_vote(pip_obj.cfg.text_proposal)
        log.info('Get text proposal information {}'.format(proposalinfo))

        assert pip_obj.get_accuverifiers_count(proposalinfo.get('ProposalID')) == [4, 1, 1, 1]
        wait_block_number(client_con_list_obj[1].node, proposalinfo.get('EndVotingBlock'))
        assert_code(pip_obj.get_yeas_of_proposal(proposalinfo.get('ProposalID')), 1)
        assert_code(pip_obj.get_nays_of_proposal(proposalinfo.get('ProposalID')), 1)
        assert_code(pip_obj.get_abstentions_of_proposal(proposalinfo.get('ProposalID')), 1)
        assert_code(pip_obj.get_accu_verifiers_of_proposal(proposalinfo.get('ProposalID')), len(client_con_list_obj))
        assert_code(pip_obj.get_status_of_proposal(proposalinfo.get('ProposalID')), 3)

    def test_UP_TE_002(self, new_genesis_env, client_con_list_obj):
        genesis = from_dict(data_class=Genesis, data=new_genesis_env.genesis_config)
        genesis.economicModel.gov.textProposalSupportRate = 0.334
        genesis.economicModel.gov.textProposalVoteRate = 0.749
        genesis.economicModel.gov.textProposalVoteDurationSeconds = 40
        new_genesis_env.set_genesis(genesis.to_dict())
        new_genesis_env.deploy_all()
        submittpandvote(client_con_list_obj[:3], 1, 2, 3)
        pip_obj = client_con_list_obj[0].pip
        proposalinfo = pip_obj.get_effect_proposal_info_of_vote(pip_obj.cfg.text_proposal)
        log.info('Get text proposal information {}'.format(proposalinfo))

        assert pip_obj.get_accuverifiers_count(proposalinfo.get('ProposalID')) == [4, 1, 1, 1]
        wait_block_number(client_con_list_obj[1].node, proposalinfo.get('EndVotingBlock'))
        assert_code(pip_obj.get_yeas_of_proposal(proposalinfo.get('ProposalID')), 1)
        assert_code(pip_obj.get_nays_of_proposal(proposalinfo.get('ProposalID')), 1)
        assert_code(pip_obj.get_abstentions_of_proposal(proposalinfo.get('ProposalID')), 1)
        assert_code(pip_obj.get_accu_verifiers_of_proposal(proposalinfo.get('ProposalID')), len(client_con_list_obj))
        assert_code(pip_obj.get_status_of_proposal(proposalinfo.get('ProposalID')), 3)

    def test_UP_TE_003(self, new_genesis_env, client_con_list_obj):
        genesis = from_dict(data_class=Genesis, data=new_genesis_env.genesis_config)
        genesis.economicModel.gov.textProposalSupportRate = 0.333
        genesis.economicModel.gov.textProposalVoteRate = 0.751
        genesis.economicModel.gov.textProposalVoteDurationSeconds = 40
        new_genesis_env.set_genesis(genesis.to_dict())
        new_genesis_env.deploy_all()
        submittpandvote(client_con_list_obj[:3], 1, 2, 3)
        pip_obj = client_con_list_obj[0].pip
        proposalinfo = pip_obj.get_effect_proposal_info_of_vote(pip_obj.cfg.text_proposal)
        log.info('Get text proposal information {}'.format(proposalinfo))

        assert pip_obj.get_accuverifiers_count(proposalinfo.get('ProposalID')) == [4, 1, 1, 1]
        wait_block_number(client_con_list_obj[1].node, proposalinfo.get('EndVotingBlock'))
        assert_code(pip_obj.get_yeas_of_proposal(proposalinfo.get('ProposalID')), 1)
        assert_code(pip_obj.get_nays_of_proposal(proposalinfo.get('ProposalID')), 1)
        assert_code(pip_obj.get_abstentions_of_proposal(proposalinfo.get('ProposalID')), 1)
        assert_code(pip_obj.get_accu_verifiers_of_proposal(proposalinfo.get('ProposalID')), len(client_con_list_obj))
        assert_code(pip_obj.get_status_of_proposal(proposalinfo.get('ProposalID')), 3)

    def test_UP_TE_004(self, new_genesis_env, client_con_list_obj):
        genesis = from_dict(data_class=Genesis, data=new_genesis_env.genesis_config)
        genesis.economicModel.gov.textProposalSupportRate = 0.334
        genesis.economicModel.gov.textProposalVoteRate = 0.75
        genesis.economicModel.gov.textProposalVoteDurationSeconds = 40
        new_genesis_env.set_genesis(genesis.to_dict())
        new_genesis_env.deploy_all()
        submittpandvote(client_con_list_obj[:3], 1, 2, 3)
        pip_obj = client_con_list_obj[0].pip
        proposalinfo = pip_obj.get_effect_proposal_info_of_vote(pip_obj.cfg.text_proposal)
        log.info('Get text proposal information {}'.format(proposalinfo))

        assert pip_obj.get_accuverifiers_count(proposalinfo.get('ProposalID')) == [4, 1, 1, 1]
        wait_block_number(client_con_list_obj[1].node, proposalinfo.get('EndVotingBlock'))
        assert_code(pip_obj.get_yeas_of_proposal(proposalinfo.get('ProposalID')), 1)
        assert_code(pip_obj.get_nays_of_proposal(proposalinfo.get('ProposalID')), 1)
        assert_code(pip_obj.get_abstentions_of_proposal(proposalinfo.get('ProposalID')), 1)
        assert_code(pip_obj.get_accu_verifiers_of_proposal(proposalinfo.get('ProposalID')), len(client_con_list_obj))
        assert_code(pip_obj.get_status_of_proposal(proposalinfo.get('ProposalID')), 3)

    def test_UP_TE_005(self, new_genesis_env, client_con_list_obj):
        genesis = from_dict(data_class=Genesis, data=new_genesis_env.genesis_config)
        genesis.economicModel.gov.textProposalSupportRate = 0.332
        genesis.economicModel.gov.textProposalVoteRate = 0.749
        genesis.economicModel.gov.textProposalVoteDurationSeconds = 40
        new_genesis_env.set_genesis(genesis.to_dict())
        new_genesis_env.deploy_all()
        submittpandvote(client_con_list_obj[:3], 1, 2, 3)
        pip_obj = client_con_list_obj[0].pip
        proposalinfo = pip_obj.get_effect_proposal_info_of_vote(pip_obj.cfg.text_proposal)
        log.info('Get text proposal information {}'.format(proposalinfo))

        assert pip_obj.get_accuverifiers_count(proposalinfo.get('ProposalID')) == [4, 1, 1, 1]
        wait_block_number(client_con_list_obj[1].node, proposalinfo.get('EndVotingBlock'))
        assert_code(pip_obj.get_yeas_of_proposal(proposalinfo.get('ProposalID')), 1)
        assert_code(pip_obj.get_nays_of_proposal(proposalinfo.get('ProposalID')), 1)
        assert_code(pip_obj.get_abstentions_of_proposal(proposalinfo.get('ProposalID')), 1)
        assert_code(pip_obj.get_accu_verifiers_of_proposal(proposalinfo.get('ProposalID')), len(client_con_list_obj))
        assert_code(pip_obj.get_status_of_proposal(proposalinfo.get('ProposalID')), 2)

    def test_UP_TE_006(self, new_genesis_env, client_con_list_obj):
        genesis = from_dict(data_class=Genesis, data=new_genesis_env.genesis_config)
        genesis.economicModel.gov.textProposalSupportRate = 0.332
        genesis.economicModel.gov.textProposalVoteRate = 0.749
        genesis.economicModel.gov.textProposalVoteDurationSeconds = 40
        new_genesis_env.set_genesis(genesis.to_dict())
        new_genesis_env.deploy_all()
        submittpandvote(client_con_list_obj[:3], 1, 2, 3)
        pip_obj = client_con_list_obj[0].pip
        proposalinfo = pip_obj.get_effect_proposal_info_of_vote(pip_obj.cfg.text_proposal)
        log.info('Get text proposal information {}'.format(proposalinfo))

        assert pip_obj.get_accuverifiers_count(proposalinfo.get('ProposalID')) == [4, 1, 1, 1]
        wait_block_number(client_con_list_obj[1].node, proposalinfo.get('EndVotingBlock'))
        assert_code(pip_obj.get_yeas_of_proposal(proposalinfo.get('ProposalID')), 1)
        assert_code(pip_obj.get_nays_of_proposal(proposalinfo.get('ProposalID')), 1)
        assert_code(pip_obj.get_abstentions_of_proposal(proposalinfo.get('ProposalID')), 1)
        assert_code(pip_obj.get_accu_verifiers_of_proposal(proposalinfo.get('ProposalID')), len(client_con_list_obj))
        assert_code(pip_obj.get_status_of_proposal(proposalinfo.get('ProposalID')), 2)

    def test_UP_TE_007(self, new_genesis_env, client_con_list_obj):
        genesis = from_dict(data_class=Genesis, data=new_genesis_env.genesis_config)
        genesis.economicModel.gov.textProposalSupportRate = 0.332
        genesis.economicModel.gov.textProposalVoteRate = 0.75
        genesis.economicModel.gov.textProposalVoteDurationSeconds = 40
        new_genesis_env.set_genesis(genesis.to_dict())
        new_genesis_env.deploy_all()
        submittpandvote(client_con_list_obj[:3], 1, 2, 3)
        pip_obj = client_con_list_obj[0].pip
        proposalinfo = pip_obj.get_effect_proposal_info_of_vote(pip_obj.cfg.text_proposal)
        log.info('Get text proposal information {}'.format(proposalinfo))

        assert pip_obj.get_accuverifiers_count(proposalinfo.get('ProposalID')) == [4, 1, 1, 1]
        wait_block_number(client_con_list_obj[1].node, proposalinfo.get('EndVotingBlock'))
        assert_code(pip_obj.get_yeas_of_proposal(proposalinfo.get('ProposalID')), 1)
        assert_code(pip_obj.get_nays_of_proposal(proposalinfo.get('ProposalID')), 1)
        assert_code(pip_obj.get_abstentions_of_proposal(proposalinfo.get('ProposalID')), 1)
        assert_code(pip_obj.get_accu_verifiers_of_proposal(proposalinfo.get('ProposalID')), len(client_con_list_obj))
        assert_code(pip_obj.get_status_of_proposal(proposalinfo.get('ProposalID')), 3)