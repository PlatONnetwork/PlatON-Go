import pytest
import allure
from common.log import log
from common.key import mock_duplicate_sign
import time
from tests.lib.utils import assert_code, wait_block_number, upload_platon
from tests.lib.client import Client
from tests.lib import Genesis
from dacite import from_dict
from tests.govern.conftest import proposal_vote, version_proposal_vote


def submitvpandvote(clients, votingrounds=3, version=None):
    pip = clients[0].pip
    if version is None:
        version = pip.cfg.version5
    result = pip.submitVersion(pip.node.node_id, str(time.time()), version, votingrounds, pip.node.staking_address,
                               transaction_cfg=pip.cfg.transaction_cfg)
    log.info('Submit version proposal result : {}'.format(result))
    assert_code(result, 0)
    proposalinfo = pip.get_effect_proposal_info_of_vote()
    log.info('Version proposal info {}'.format(proposalinfo))
    for index in range(len(clients)):
        pip = clients[index].pip
        result = version_proposal_vote(pip, vote_option=pip.cfg.vote_option_yeas)
        log.info('Node {} vote result {}'.format(pip.node.node_id, result))
        assert_code(result, 0)


def createstaking(obj, platon_bin=None, reward_per=0):
    if isinstance(obj, Client):
        objs = []
        objs.append(obj)
        obj = objs
    for client in obj:
        if platon_bin:
            log.info('Need replace the platon of the node')
            upload_platon(client.node, platon_bin)
            client.node.restart()

        address, _ = client.economic.account.generate_account(client.node.web3, 10 ** 18 * 10000000)
        result = client.staking.create_staking(0, address, address, amount=10 ** 18 * 2000000,
                                               transaction_cfg=client.pip.cfg.transaction_cfg, reward_per=reward_per)
        log.info('Node {} staking result : {}'.format(client.node.node_id, result))
        assert_code(result, 0)


def submitppandvote(clients, *args):
    pip = clients[0].pip
    result = pip.submitParam(pip.node.node_id, str(time.time()), 'slashing', 'slashBlocksReward', '83',
                             pip.node.staking_address, transaction_cfg=pip.cfg.transaction_cfg)
    log.info('Submit param proposal result : {}'.format(result))
    assert_code(result, 0)
    proposalinfo = pip.get_effect_proposal_info_of_vote(pip.cfg.param_proposal)
    log.info('Param proposal info {}'.format(proposalinfo))
    for index in range(len(clients)):
        pip = clients[index].pip
        log.info('{}'.format(args[index]))
        result = proposal_vote(pip, vote_option=args[index])
        assert_code(result, 0)


def submitcppandvote(clients, list, voting_rounds=2):
    pip = clients[0].pip
    result = pip.submitParam(pip.node.node_id, str(time.time()), 'slashing', 'slashBlocksReward', '83',
                             pip.node.staking_address, transaction_cfg=pip.cfg.transaction_cfg)
    log.info('Submit param proposal result : {}'.format(result))
    assert_code(result, 0)
    proposalinfo_param = pip.get_effect_proposal_info_of_vote(pip.cfg.param_proposal)
    log.info('Param proposal info {}'.format(proposalinfo_param))

    result = pip.submitCancel(pip.node.node_id, str(time.time()), voting_rounds, proposalinfo_param.get('ProposalID'),
                              pip.node.staking_address, transaction_cfg=pip.cfg.transaction_cfg)
    log.info('Submit cancel proposal result : {}'.format(result))
    assert_code(result, 0)
    proposalinfo_cancel = pip.get_effect_proposal_info_of_vote(pip.cfg.cancel_proposal)
    log.info('Cancel proposal info {}'.format(proposalinfo_cancel))
    for index in range(len(clients)):
        pip = clients[index].pip
        log.info('Vote option {}'.format(list[index]))
        result = pip.vote(pip.node.node_id, proposalinfo_cancel.get('ProposalID'), list[index],
                          pip.node.staking_address, transaction_cfg=pip.cfg.transaction_cfg)
        log.info('Node {} vote cancel proposal result : {}'.format(pip.node.node_id, result))
        assert_code(result, 0)


def submitcvpandvote(clients, *args):
    pip = clients[0].pip
    result = pip.submitVersion(pip.node.node_id, str(time.time()), pip.cfg.version5, 3,
                               pip.node.staking_address, transaction_cfg=pip.cfg.transaction_cfg)
    log.info('Submit version proposal result : {}'.format(result))
    assert_code(result, 0)
    proposalinfo_version = pip.get_effect_proposal_info_of_vote(pip.cfg.version_proposal)
    log.info('Version proposal info {}'.format(proposalinfo_version))

    result = pip.submitCancel(pip.node.node_id, str(time.time()), 1, proposalinfo_version.get('ProposalID'),
                              pip.node.staking_address, transaction_cfg=pip.cfg.transaction_cfg)
    log.info('Submit cancel proposal result : {}'.format(result))
    assert_code(result, 0)
    proposalinfo_cancel = pip.get_effect_proposal_info_of_vote(pip.cfg.cancel_proposal)
    log.info('Cancel proposal info {}'.format(proposalinfo_cancel))
    for index in range(len(clients)):
        pip = clients[index].pip
        log.info('{}'.format(args[index]))
        result = pip.vote(pip.node.node_id, proposalinfo_cancel.get('ProposalID'), args[index],
                          pip.node.staking_address, transaction_cfg=pip.cfg.transaction_cfg)
        log.info('Node {} vote cancel proposal result : {}'.format(pip.node.node_id, result))
        assert_code(result, 0)


def submittpandvote(clients, *args):
    pip = clients[0].pip
    result = pip.submitText(pip.node.node_id, str(time.time()), pip.node.staking_address,
                            transaction_cfg=pip.cfg.transaction_cfg)
    log.info('Submit text proposal result : {}'.format(result))
    assert_code(result, 0)
    proposalinfo_text = pip.get_effect_proposal_info_of_vote(pip.cfg.text_proposal)
    log.info('Text proposal info {}'.format(proposalinfo_text))
    result = pip.pip.listProposal()
    print(result)
    for index in range(len(clients)):
        pip = clients[index].pip
        log.info('{}'.format(args[index]))
        result = pip.vote(pip.node.node_id, proposalinfo_text.get('ProposalID'), args[index],
                          pip.node.staking_address, transaction_cfg=pip.cfg.transaction_cfg)
        log.info('Node {} vote text proposal result : {}'.format(pip.node.node_id, result))
        assert_code(result, 0)


class TestVotingStatisticsVP:
    @pytest.mark.compatibility
    @pytest.mark.P0
    @allure.title('Version proposal statistics function verification')
    def test_VS_EXV_001_VS_BL_1(self, new_genesis_env, clients_consensus, clients_noconsensus):
        new_genesis_env.deploy_all()
        submitvpandvote(clients_consensus[0:-1])
        proposalinfo = clients_consensus[0].pip.get_effect_proposal_info_of_vote(clients_consensus[0].pip.cfg.version_proposal)
        log.info('Version proposal information {}'.format(proposalinfo))
        createstaking(clients_noconsensus[:3])
        wait_block_number(clients_consensus[0].node, proposalinfo.get('EndVotingBlock'))
        result = clients_consensus[0].pip.get_accuverifiers_count(proposalinfo.get('ProposalID'))
        log.info('Get proposal vote infomation {}'.format(result))
        assert result == [4, 3, 0, 0]
        assert clients_consensus[0].pip.get_status_of_proposal(proposalinfo.get('ProposalID')) == 4
        assert clients_consensus[0].pip.get_yeas_of_proposal(proposalinfo.get('ProposalID')) == 3
        assert clients_consensus[0].pip.get_nays_of_proposal(proposalinfo.get('ProposalID')) == 0
        assert clients_consensus[0].pip.get_abstentions_of_proposal(proposalinfo.get('ProposalID')) == 0
        assert clients_consensus[0].pip.get_accu_verifiers_of_proposal(proposalinfo.get('ProposalID')) == len(
            clients_consensus)

    @pytest.mark.P1
    @allure.title('Version proposal statistics function verification')
    def test_VS_EXV_002(self, new_genesis_env, clients_consensus, clients_noconsensus):
        new_genesis_env.deploy_all()
        pip = clients_consensus[0].pip
        submitvpandvote(clients_consensus[:2], votingrounds=5)
        proposalinfo = pip.get_effect_proposal_info_of_vote(pip.cfg.version_proposal)
        log.info('Version proposal info {}'.format(proposalinfo))
        log.info('{}'.format(clients_consensus[:2]))
        createstaking(clients_noconsensus[:2])
        pip.economic.wait_settlement_blocknum(pip.node)
        result = pip.get_accuverifiers_count(proposalinfo.get('ProposalID'))
        log.info('Get proposal vote infomation {}'.format(result))
        assert result == [6, 2, 0, 0]

        result = version_proposal_vote(clients_consensus[2].pip)
        assert_code(result, 0)
        createstaking(clients_noconsensus[2])
        wait_block_number(pip.node, proposalinfo.get('EndVotingBlock'))
        result = pip.get_accuverifiers_count(proposalinfo.get('ProposalID'))
        log.info('Get proposal vote infomation {}'.format(result))
        assert result == [6, 3, 0, 0]
        assert clients_consensus[0].pip.get_status_of_proposal(proposalinfo.get('ProposalID')) == 3
        assert clients_consensus[0].pip.get_yeas_of_proposal(proposalinfo.get('ProposalID')) == 3
        assert clients_consensus[0].pip.get_nays_of_proposal(proposalinfo.get('ProposalID')) == 0
        assert clients_consensus[0].pip.get_abstentions_of_proposal(proposalinfo.get('ProposalID')) == 0
        assert clients_consensus[0].pip.get_accu_verifiers_of_proposal(proposalinfo.get('ProposalID')) == 6

    @pytest.mark.P2
    @allure.title('Version proposal statistics function verification')
    def test_VS_EXV_003(self, new_genesis_env, clients_consensus, clients_noconsensus):
        new_genesis_env.deploy_all()
        submitvpandvote(clients_consensus[:1], votingrounds=9)
        proposalinfo = clients_consensus[0].pip.get_effect_proposal_info_of_vote(clients_consensus[0].pip.cfg.version_proposal)
        log.info('Version proposal info {}'.format(proposalinfo))
        createstaking(clients_noconsensus[0])
        clients_consensus[0].pip.economic.wait_settlement_blocknum(clients_consensus[0].pip.node)
        result = clients_consensus[0].pip.get_accuverifiers_count(proposalinfo.get('ProposalID'))
        log.info('Get proposal vote infomation {}'.format(result))
        assert result == [5, 1, 0, 0]

        result = version_proposal_vote(clients_consensus[1].pip, clients_consensus[0].pip.cfg.vote_option_yeas)
        assert_code(result, 0)
        createstaking(clients_noconsensus[1])
        clients_consensus[0].pip.economic.wait_settlement_blocknum(clients_consensus[0].pip.node)
        result = clients_consensus[0].pip.get_accuverifiers_count(proposalinfo.get('ProposalID'))
        log.info('Get proposal vote infomation {}'.format(result))
        assert result == [6, 2, 0, 0]

        result = version_proposal_vote(clients_consensus[2].pip, clients_consensus[0].pip.cfg.vote_option_yeas)
        assert_code(result, 0)
        createstaking(clients_noconsensus[2])
        wait_block_number(clients_consensus[0].pip.node, proposalinfo.get('EndVotingBlock'))

        result = clients_consensus[0].pip.get_accuverifiers_count(proposalinfo.get('ProposalID'))
        log.info('Get proposal vote infomation {}'.format(result))
        assert result == [6, 3, 0, 0]
        assert clients_consensus[0].pip.get_status_of_proposal(proposalinfo.get('ProposalID')) == 3
        assert clients_consensus[0].pip.get_yeas_of_proposal(proposalinfo.get('ProposalID')) == 3
        assert clients_consensus[0].pip.get_nays_of_proposal(proposalinfo.get('ProposalID')) == 0
        assert clients_consensus[0].pip.get_abstentions_of_proposal(proposalinfo.get('ProposalID')) == 0
        assert clients_consensus[0].pip.get_accu_verifiers_of_proposal(proposalinfo.get('ProposalID')) == 6

    @pytest.mark.P2
    @allure.title('Version proposal statistics function verification')
    def test_VS_EXV_004(self, new_genesis_env, clients_consensus):
        genesis = from_dict(data_class=Genesis, data=new_genesis_env.genesis_config)
        genesis.economicModel.gov.versionProposalSupportRate = 2500
        new_genesis_env.set_genesis(genesis.to_dict())
        new_genesis_env.deploy_all()
        submitvpandvote(clients_consensus[:2])
        proposalinfo = clients_consensus[0].pip.get_effect_proposal_info_of_vote(clients_consensus[0].pip.cfg.version_proposal)
        log.info('Get version proposal information {}'.format(proposalinfo))
        result = clients_consensus[0].staking.withdrew_staking(clients_consensus[0].node.staking_address)
        log.info('Node {} withdrew staking result {}'.format(clients_consensus[0].node.node_id, result))
        assert_code(result, 0)
        wait_block_number(clients_consensus[0].node, proposalinfo.get('EndVotingBlock'))
        result = clients_consensus[0].pip.get_accuverifiers_count(proposalinfo.get('ProposalID'))
        log.info('Get proposal vote infomation {}'.format(result))
        assert result == [4, 2, 0, 0]
        assert clients_consensus[0].pip.get_yeas_of_proposal(proposalinfo.get('ProposalID')) == 2
        assert clients_consensus[0].pip.get_status_of_proposal(proposalinfo.get('ProposalID')) == 4

    @pytest.mark.P2
    @allure.title('Version proposal statistics function verification')
    def test_VS_EXV_005(self, new_genesis_env, clients_consensus):
        genesis = from_dict(data_class=Genesis, data=new_genesis_env.genesis_config)
        genesis.economicModel.gov.versionProposalSupportRate = 5000
        new_genesis_env.set_genesis(genesis.to_dict())
        new_genesis_env.deploy_all()
        submitvpandvote(clients_consensus[:2])
        pip = clients_consensus[0].pip
        proposalinfo = pip.get_effect_proposal_info_of_vote(pip.cfg.version_proposal)
        log.info('Get version proposal information {}'.format(proposalinfo))
        wait_block_number(pip.node, proposalinfo.get('EndVotingBlock'))
        report_information = mock_duplicate_sign(1, pip.node.nodekey, pip.node.blsprikey,
                                                 proposalinfo.get('EndVotingBlock') - 10)
        log.info("Report information: {}".format(report_information))
        address, _ = pip.economic.account.generate_account(pip.node.web3, 10**18 * 1000)
        result = clients_consensus[0].duplicatesign.reportDuplicateSign(1, report_information, address)
        assert_code(result, 0)
        assert clients_consensus[1].pip.get_yeas_of_proposal(proposalinfo.get('ProposalID')) == 2

    @pytest.mark.P2
    @allure.title('Version proposal statistics function verification')
    def test_VS_EXV_006(self, new_genesis_env, clients_consensus):
        genesis = from_dict(data_class=Genesis, data=new_genesis_env.genesis_config)
        genesis.economicModel.gov.versionProposalSupportRate = 5000
        new_genesis_env.set_genesis(genesis.to_dict())
        new_genesis_env.deploy_all()
        submitvpandvote(clients_consensus[:2])
        pip = clients_consensus[0].pip
        proposalinfo = pip.get_effect_proposal_info_of_vote(pip.cfg.version_proposal)
        log.info('Get version proposal information {}'.format(proposalinfo))
        wait_block_number(pip.node, 80)
        report_information = mock_duplicate_sign(1, pip.node.nodekey, pip.node.blsprikey, 70)
        log.info("Report information: {}".format(report_information))
        address, _ = pip.economic.account.generate_account(pip.node.web3, 10**18 * 1000)
        result = clients_consensus[0].duplicatesign.reportDuplicateSign(1, report_information, address)
        assert_code(result, 0)
        wait_block_number(pip.node, proposalinfo.get('EndVotingBlock'))
        assert clients_consensus[0].pip.get_yeas_of_proposal(proposalinfo.get('ProposalID')) == 1
        assert clients_consensus[0].pip.get_status_of_proposal(proposalinfo.get('ProposalID')) == 3

    @pytest.mark.P2
    @allure.title('Version proposal statistics function verification')
    def test_VS_EXV_007(self, new_genesis_env, clients_consensus):
        genesis = from_dict(data_class=Genesis, data=new_genesis_env.genesis_config)
        genesis.economicModel.gov.versionProposalSupportRate = 5000
        new_genesis_env.set_genesis(genesis.to_dict())
        new_genesis_env.deploy_all()
        submitvpandvote(clients_consensus[:2])
        pip = clients_consensus[0].pip
        proposalinfo = pip.get_effect_proposal_info_of_vote(pip.cfg.version_proposal)
        log.info('Get version proposal information {}'.format(proposalinfo))
        wait_block_number(pip.node, proposalinfo.get('EndVotingBlock'))
        report_information = mock_duplicate_sign(2, pip.node.nodekey, pip.node.blsprikey,
                                                 proposalinfo.get('EndVotingBlock') - 10)
        log.info("Report information: {}".format(report_information))
        address, _ = pip.economic.account.generate_account(pip.node.web3, 10**18 * 1000)
        result = clients_consensus[0].duplicatesign.reportDuplicateSign(2, report_information, address)
        assert_code(result, 0)
        assert clients_consensus[1].pip.get_yeas_of_proposal(proposalinfo.get('ProposalID')) == 2

    @pytest.mark.P2
    @allure.title('Version proposal statistics function verification')
    def test_VS_EXV_008(self, new_genesis_env, clients_consensus):
        genesis = from_dict(data_class=Genesis, data=new_genesis_env.genesis_config)
        genesis.economicModel.gov.versionProposalSupportRate = 5000
        new_genesis_env.set_genesis(genesis.to_dict())
        new_genesis_env.deploy_all()
        submitvpandvote(clients_consensus[:2])
        pip = clients_consensus[0].pip
        proposalinfo = pip.get_effect_proposal_info_of_vote(pip.cfg.version_proposal)
        log.info('Get version proposal information {}'.format(proposalinfo))
        wait_block_number(pip.node, 80)
        report_information = mock_duplicate_sign(2, pip.node.nodekey, pip.node.blsprikey, 70)
        log.info("Report information: {}".format(report_information))
        address, _ = pip.economic.account.generate_account(pip.node.web3, 10**18 * 1000)
        result = clients_consensus[0].duplicatesign.reportDuplicateSign(2, report_information, address)
        assert_code(result, 0)
        wait_block_number(pip.node, proposalinfo.get('EndVotingBlock'))
        assert clients_consensus[0].pip.get_yeas_of_proposal(proposalinfo.get('ProposalID')) == 1
        assert clients_consensus[0].pip.get_status_of_proposal(proposalinfo.get('ProposalID')) == 3

    @pytest.mark.P2
    @allure.title('Version proposal statistics function verification')
    def test_VS_EXV_009(self, new_genesis_env, clients_consensus):
        genesis = from_dict(data_class=Genesis, data=new_genesis_env.genesis_config)
        genesis.economicModel.gov.versionProposalSupportRate = 5000
        new_genesis_env.set_genesis(genesis.to_dict())
        new_genesis_env.deploy_all()
        submitvpandvote(clients_consensus[:2], votingrounds=3)
        pip = clients_consensus[0].pip
        proposalinfo = pip.get_effect_proposal_info_of_vote(pip.cfg.version_proposal)
        log.info('Get version proposal information {}'.format(proposalinfo))
        assert pip.get_accuverifiers_count(proposalinfo.get('ProposalID')) == [4, 2, 0, 0]
        log.info('Stop the node {}'.format(pip.node.node_id))
        pip.node.stop()
        wait_block_number(clients_consensus[1].node, proposalinfo.get('EndVotingBlock'))
        assert clients_consensus[1].pip.get_yeas_of_proposal(proposalinfo.get('ProposalID')) == 1
        assert clients_consensus[1].pip.get_status_of_proposal(proposalinfo.get('ProposalID')) == 3

    @pytest.mark.P2
    @allure.title('Version proposal statistics function verification')
    def test_VS_EXV_010(self, new_genesis_env, clients_consensus):
        genesis = from_dict(data_class=Genesis, data=new_genesis_env.genesis_config)
        genesis.economicModel.gov.versionProposalSupportRate = 5000
        new_genesis_env.set_genesis(genesis.to_dict())
        new_genesis_env.deploy_all()
        submitvpandvote(clients_consensus[:2], votingrounds=3)
        pip = clients_consensus[0].pip
        proposalinfo = pip.get_effect_proposal_info_of_vote(pip.cfg.version_proposal)
        log.info('Get version proposal information {}'.format(proposalinfo))
        assert pip.get_accuverifiers_count(proposalinfo.get('ProposalID')) == [4, 2, 0, 0]
        log.info('Stop the node {}'.format(pip.node.node_id))
        wait_block_number(clients_consensus[1].node, proposalinfo.get('EndVotingBlock') - 10)
        pip.node.stop()
        wait_block_number(clients_consensus[1].node, proposalinfo.get('EndVotingBlock'))
        assert clients_consensus[1].pip.get_yeas_of_proposal(proposalinfo.get('ProposalID')) == 2
        assert clients_consensus[1].pip.get_status_of_proposal(proposalinfo.get('ProposalID')) == 4


class TestVotingStatisticsTPCP:
    @pytest.mark.P1
    @allure.title('Cancel proposal and text proposal statistics function verification')
    def test_VS_EXT_001_VS_EXC_001(self, new_genesis_env, clients_consensus, clients_noconsensus):
        genesis = from_dict(data_class=Genesis, data=new_genesis_env.genesis_config)
        genesis.economicModel.gov.textProposalVoteDurationSeconds = 120
        new_genesis_env.set_genesis(genesis.to_dict())
        new_genesis_env.deploy_all()
        submittpandvote(clients_consensus[0:-1], 1, 2, 3)
        submitcppandvote(clients_consensus[0:-1], [1, 2, 3])
        proposalinfo_text = clients_consensus[0].pip.get_effect_proposal_info_of_vote(
            clients_consensus[0].pip.cfg.text_proposal)
        log.info('Text proposal information {}'.format(proposalinfo_text))
        proposalinfo_cancel = clients_consensus[0].pip.get_effect_proposal_info_of_vote(
            clients_consensus[0].pip.cfg.cancel_proposal)
        log.info('Cancel proposal information {}'.format(proposalinfo_cancel))
        createstaking(clients_noconsensus[:3])
        result_text = clients_consensus[0].pip.get_accuverifiers_count(proposalinfo_text.get('ProposalID'))
        log.info('Get text proposal vote infomation {}'.format(result_text))
        result_cancel = clients_consensus[0].pip.get_accuverifiers_count(proposalinfo_cancel.get('ProposalID'))
        log.info('Get cancel proposal vote infomation {}'.format(result_cancel))
        assert result_text == [4, 1, 1, 1]
        assert result_cancel == [4, 1, 1, 1]

    @pytest.mark.P1
    @allure.title('Cancel proposal and text proposal statistics function verification')
    def test_VS_EXT_002_VS_EXC_002(self, new_genesis_env, clients_consensus, clients_noconsensus):
        genesis = from_dict(data_class=Genesis, data=new_genesis_env.genesis_config)
        genesis.economicModel.gov.textProposalVoteDurationSeconds = 200
        genesis.economicModel.gov.paramProposalVoteDurationSeconds = 1000
        new_genesis_env.set_genesis(genesis.to_dict())
        new_genesis_env.deploy_all()
        pip = clients_consensus[0].pip
        submittpandvote(clients_consensus[:2], 1, 2)
        submitcppandvote(clients_consensus[:2], [1, 2], voting_rounds=5)
        proposalinfo_text = pip.get_effect_proposal_info_of_vote(pip.cfg.text_proposal)
        log.info('Text proposal info {}'.format(proposalinfo_text))
        proposalinfo_cancel = pip.get_effect_proposal_info_of_vote(pip.cfg.cancel_proposal)
        log.info('Cancel proposal info {}'.format(proposalinfo_cancel))
        createstaking(clients_noconsensus[:2])
        pip.economic.wait_settlement_blocknum(pip.node)
        result_text = pip.get_accuverifiers_count(proposalinfo_text.get('ProposalID'))
        log.info('Get text proposal vote infomation {}'.format(result_text))
        result_cancel = pip.get_accuverifiers_count(proposalinfo_cancel.get('ProposalID'))
        log.info('Get cancel proposal vote infomation {}'.format(result_cancel))
        assert result_text == [6, 1, 1, 0]
        assert result_cancel == [6, 1, 1, 0]

        result_cancel = proposal_vote(clients_consensus[2].pip, vote_option=pip.cfg.vote_option_Abstentions,
                                      proposaltype=pip.cfg.cancel_proposal)
        assert_code(result_cancel, 0)
        result_text = proposal_vote(clients_consensus[2].pip, vote_option=pip.cfg.vote_option_Abstentions,
                                    proposaltype=pip.cfg.text_proposal)
        assert_code(result_text, 0)
        createstaking(clients_noconsensus[2])
        wait_block_number(pip.node, proposalinfo_text.get('EndVotingBlock'))

        result_cancel = pip.get_accuverifiers_count(proposalinfo_cancel.get('ProposalID'))
        log.info('Get proposal vote infomation {}'.format(result_cancel))
        result_text = pip.get_accuverifiers_count(proposalinfo_text.get('ProposalID'))
        log.info('Get proposal vote infomation {}'.format(result_text))

        assert result_text == [6, 1, 1, 1]
        assert result_cancel == [6, 1, 1, 1]

    @pytest.mark.P2
    @allure.title('Cancel proposal and text proposal statistics function verification')
    def test_VS_EXT_003_VS_EXC_003(self, new_genesis_env, clients_consensus, clients_noconsensus):
        genesis = from_dict(data_class=Genesis, data=new_genesis_env.genesis_config)
        genesis.economicModel.gov.textProposalVoteDurationSeconds = 360
        genesis.economicModel.gov.paramProposalVoteDurationSeconds = 600
        new_genesis_env.set_genesis(genesis.to_dict())
        new_genesis_env.deploy_all()
        pip = clients_consensus[0].pip
        submittpandvote(clients_consensus[:1], 1)
        submitcppandvote(clients_consensus[:1], [1], voting_rounds=9)
        proposalinfo_text = pip.get_effect_proposal_info_of_vote(pip.cfg.text_proposal)
        log.info('Text proposal info {}'.format(proposalinfo_text))

        proposalinfo_cancel = pip.get_effect_proposal_info_of_vote(pip.cfg.cancel_proposal)
        log.info('Cancel proposal info {}'.format(proposalinfo_cancel))
        createstaking(clients_noconsensus[0])
        pip.economic.wait_settlement_blocknum(pip.node)
        result_text = pip.get_accuverifiers_count(proposalinfo_text.get('ProposalID'))
        log.info('Get proposal vote infomation {}'.format(result_text))
        result_cancel = pip.get_accuverifiers_count(proposalinfo_cancel.get('ProposalID'))
        log.info('Get proposal vote infomation {}'.format(result_cancel))
        assert result_text == [5, 1, 0, 0]
        assert result_cancel == [5, 1, 0, 0]

        result = proposal_vote(clients_consensus[1].pip, vote_option=pip.cfg.vote_option_nays,
                               proposaltype=pip.cfg.cancel_proposal)
        assert_code(result, 0)
        result = proposal_vote(clients_consensus[1].pip, vote_option=pip.cfg.vote_option_nays,
                               proposaltype=pip.cfg.text_proposal)
        assert_code(result, 0)
        createstaking(clients_noconsensus[1])
        pip.economic.wait_settlement_blocknum(pip.node)
        result_cancel = pip.get_accuverifiers_count(proposalinfo_cancel.get('ProposalID'))
        log.info('Get proposal vote infomation {}'.format(result_cancel))
        result_text = pip.get_accuverifiers_count(proposalinfo_text.get('ProposalID'))
        log.info('Get proposal vote infomation {}'.format(result_text))
        assert result_text == [6, 1, 1, 0]
        assert result_cancel == [6, 1, 1, 0]

        result = proposal_vote(clients_consensus[2].pip, vote_option=pip.cfg.vote_option_Abstentions,
                               proposaltype=pip.cfg.text_proposal)
        assert_code(result, 0)
        result = proposal_vote(clients_consensus[2].pip, vote_option=pip.cfg.vote_option_Abstentions,
                               proposaltype=pip.cfg.cancel_proposal)
        assert_code(result, 0)
        createstaking(clients_noconsensus[2])
        # wait_block_number(clients_consensus[0].pip.node, proposalinfo.get('EndVotingBlock'))

        result_text = clients_consensus[0].pip.get_accuverifiers_count(proposalinfo_text.get('ProposalID'))
        log.info('Get proposal vote infomation {}'.format(result_text))
        result_cancel = clients_consensus[0].pip.get_accuverifiers_count(proposalinfo_cancel.get('ProposalID'))
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

    def update_setting(self, new_genesis_env, *args):
        genesis = from_dict(data_class=Genesis, data=new_genesis_env.genesis_config)
        genesis.economicModel.gov.paramProposalVoteDurationSeconds = args[0]
        genesis.economicModel.gov.textProposalVoteDurationSeconds = args[1]
        genesis.economicModel.gov.textProposalSupportRate = args[2]
        genesis.economicModel.gov.textProposalVoteRate = args[3]
        genesis.economicModel.gov.cancelProposalSupportRate = args[4]
        genesis.economicModel.gov.cancelProposalVoteRate = args[5]
        new_genesis_env.set_genesis(genesis.to_dict())
        new_genesis_env.deploy_all()

    @pytest.mark.P2
    @allure.title('Cancel proposal and text proposal statistics function verification')
    def test_VS_EXT_004_VS_EXC_004(self, new_genesis_env, clients_consensus):
        self.update_setting(new_genesis_env, 500, 80, 10000, 4990, 10000, 4990)
        pip = clients_consensus[0].pip
        submittpandvote(clients_consensus[:2], 1, 1)
        submitcppandvote(clients_consensus[:2], [1, 1])
        proposalinfo_text = pip.get_effect_proposal_info_of_vote(pip.cfg.text_proposal)
        log.info('Get text proposal information {}'.format(proposalinfo_text))
        proposalinfo_cancel = pip.get_effect_proposal_info_of_vote(pip.cfg.cancel_proposal)
        log.info('Get cancel proposal information {}'.format(proposalinfo_cancel))
        result = clients_consensus[0].staking.withdrew_staking(clients_consensus[0].node.staking_address)
        log.info('Node {} withdrew staking result {}'.format(clients_consensus[0].node.node_id, result))
        assert_code(result, 0)
        wait_block_number(clients_consensus[0].node, self.get_block(proposalinfo_cancel, proposalinfo_text))
        assert_code(pip.get_yeas_of_proposal(proposalinfo_text.get('ProposalID')), 2)
        assert_code(pip.get_status_of_proposal(proposalinfo_text.get('ProposalID')), 2)
        assert_code(pip.get_yeas_of_proposal(proposalinfo_cancel.get('ProposalID')), 2)
        assert_code(pip.get_status_of_proposal(proposalinfo_cancel.get('ProposalID')), 2)

    @pytest.mark.P2
    @allure.title('Cancel proposal and text proposal statistics function verification')
    def test_VS_EXT_005_VS_EXC_005(self, new_genesis_env, clients_consensus):
        self.update_setting(new_genesis_env, 500, 80, 10000, 4990, 10000, 4990)
        pip = clients_consensus[0].pip
        submittpandvote(clients_consensus[:2], 1, 1)
        submitcppandvote(clients_consensus[:2], [1, 1])
        proposalinfo_cancel = pip.get_effect_proposal_info_of_vote(pip.cfg.cancel_proposal)
        log.info('Get cancel proposal information {}'.format(proposalinfo_cancel))
        proposalinfo_text = pip.get_effect_proposal_info_of_vote(pip.cfg.text_proposal)
        log.info('Get text proposal information {}'.format(proposalinfo_text))
        wait_block_number(pip.node, self.get_block(proposalinfo_text, proposalinfo_cancel))
        report_information = mock_duplicate_sign(1, pip.node.nodekey, pip.node.blsprikey,
                                                 proposalinfo_cancel.get('EndVotingBlock') - 10)
        log.info("Report information: {}".format(report_information))
        address, _ = pip.economic.account.generate_account(pip.node.web3, 10 ** 18 * 1000)
        result = clients_consensus[0].duplicatesign.reportDuplicateSign(1, report_information, address)
        assert_code(result, 0)
        assert_code(pip.get_yeas_of_proposal(proposalinfo_cancel.get('ProposalID')), 2)
        assert_code(pip.get_yeas_of_proposal(proposalinfo_text.get('ProposalID')), 2)
        assert_code(pip.get_status_of_proposal(proposalinfo_cancel.get('ProposalID')), 2)
        assert_code(pip.get_status_of_proposal(proposalinfo_text.get('ProposalID')), 2)

    @pytest.mark.P2
    @allure.title('Cancel proposal and text proposal statistics function verification')
    def test_VS_EXT_006_VS_EXC_006(self, new_genesis_env, clients_consensus):
        self.update_setting(new_genesis_env, 500, 80, 10000, 4990, 10000, 4990)
        pip = clients_consensus[0].pip
        submittpandvote(clients_consensus[:2], 1, 1)
        submitcppandvote(clients_consensus[:2], [1, 1])
        proposalinfo_cancel = pip.get_effect_proposal_info_of_vote(pip.cfg.cancel_proposal)
        log.info('Get cancel proposal information {}'.format(proposalinfo_cancel))
        proposalinfo_text = pip.get_effect_proposal_info_of_vote(pip.cfg.text_proposal)
        log.info('Get text proposal information {}'.format(proposalinfo_text))
        wait_block_number(pip.node, 50)
        report_information = mock_duplicate_sign(1, pip.node.nodekey, pip.node.blsprikey, 45)
        log.info("Report information: {}".format(report_information))
        address, _ = pip.economic.account.generate_account(pip.node.web3, 10 ** 18 * 1000)
        result = clients_consensus[0].duplicatesign.reportDuplicateSign(1, report_information, address)
        assert_code(result, 0)
        wait_block_number(pip.node, self.get_block(proposalinfo_cancel, proposalinfo_text))
        assert_code(pip.get_yeas_of_proposal(proposalinfo_cancel.get('ProposalID')), 1)
        assert_code(pip.get_yeas_of_proposal(proposalinfo_text.get('ProposalID')), 1)
        assert_code(pip.get_status_of_proposal(proposalinfo_cancel.get('ProposalID')), 3)
        assert_code(pip.get_status_of_proposal(proposalinfo_text.get('ProposalID')), 3)

    @pytest.mark.P2
    @allure.title('Cancel proposal and text proposal statistics function verification')
    def test_VS_EXT_007_VS_EXC_007(self, new_genesis_env, clients_consensus):
        self.update_setting(new_genesis_env, 500, 80, 10000, 2490, 10000, 2490)
        pip = clients_consensus[0].pip
        submittpandvote(clients_consensus[:2], 2, 1)
        submitcppandvote(clients_consensus[:2], [2, 1])
        proposalinfo_cancel = pip.get_effect_proposal_info_of_vote(pip.cfg.cancel_proposal)
        log.info('Get cancel proposal information {}'.format(proposalinfo_cancel))
        proposalinfo_text = pip.get_effect_proposal_info_of_vote(pip.cfg.text_proposal)
        log.info('Get text proposal information {}'.format(proposalinfo_text))
        wait_block_number(pip.node, 50)
        report_information = mock_duplicate_sign(1, pip.node.nodekey, pip.node.blsprikey, 45)
        log.info("Report information: {}".format(report_information))
        address, _ = pip.economic.account.generate_account(pip.node.web3, 10 ** 18 * 1000)
        result = clients_consensus[0].duplicatesign.reportDuplicateSign(1, report_information, address)
        assert_code(result, 0)
        wait_block_number(pip.node, self.get_block(proposalinfo_cancel, proposalinfo_text))
        assert pip.get_accuverifiers_count(proposalinfo_text.get('ProposalID')) == [4, 1, 0, 0]
        assert pip.get_accuverifiers_count(proposalinfo_cancel.get('ProposalID')) == [4, 1, 0, 0]
        assert_code(pip.get_status_of_proposal(proposalinfo_cancel.get('ProposalID')), 2)
        assert_code(pip.get_status_of_proposal(proposalinfo_text.get('ProposalID')), 2)

    @pytest.mark.P2
    @allure.title('Cancel proposal and text proposal statistics function verification')
    def test_VS_EXT_008_VS_EXC_008(self, new_genesis_env, clients_consensus):
        self.update_setting(new_genesis_env, 500, 80, 10000, 2490, 10000, 2490)
        pip = clients_consensus[0].pip
        submittpandvote(clients_consensus[:2], 3, 1)
        submitcppandvote(clients_consensus[:2], [3, 1])
        proposalinfo_cancel = pip.get_effect_proposal_info_of_vote(pip.cfg.cancel_proposal)
        log.info('Get cancel proposal information {}'.format(proposalinfo_cancel))
        proposalinfo_text = pip.get_effect_proposal_info_of_vote(pip.cfg.text_proposal)
        log.info('Get text proposal information {}'.format(proposalinfo_text))
        wait_block_number(pip.node, 50)
        report_information = mock_duplicate_sign(1, pip.node.nodekey, pip.node.blsprikey, 45)
        log.info("Report information: {}".format(report_information))
        address, _ = pip.economic.account.generate_account(pip.node.web3, 10 ** 18 * 1000)
        result = clients_consensus[0].duplicatesign.reportDuplicateSign(1, report_information, address)
        assert_code(result, 0)
        wait_block_number(pip.node, self.get_block(proposalinfo_cancel, proposalinfo_text))
        assert pip.get_accuverifiers_count(proposalinfo_text.get('ProposalID')) == [4, 1, 0, 0]
        assert pip.get_accuverifiers_count(proposalinfo_cancel.get('ProposalID')) == [4, 1, 0, 0]
        assert_code(pip.get_status_of_proposal(proposalinfo_cancel.get('ProposalID')), 2)
        assert_code(pip.get_status_of_proposal(proposalinfo_text.get('ProposalID')), 2)

    @pytest.mark.P2
    @allure.title('Cancel proposal and text proposal statistics function verification')
    def test_VS_EXT_009_VS_EXC_009(self, new_genesis_env, clients_consensus):
        self.update_setting(new_genesis_env, 500, 80, 10000, 2490, 10000, 2490)
        pip = clients_consensus[0].pip
        submittpandvote(clients_consensus[:2], 2, 1)
        submitcppandvote(clients_consensus[:2], [2, 1])
        proposalinfo_cancel = pip.get_effect_proposal_info_of_vote(pip.cfg.cancel_proposal)
        log.info('Get cancel proposal information {}'.format(proposalinfo_cancel))
        proposalinfo_text = pip.get_effect_proposal_info_of_vote(pip.cfg.text_proposal)
        log.info('Get text proposal information {}'.format(proposalinfo_text))
        wait_block_number(pip.node, self.get_block(proposalinfo_cancel, proposalinfo_text))
        report_information = mock_duplicate_sign(1, pip.node.nodekey, pip.node.blsprikey, 45)
        log.info("Report information: {}".format(report_information))
        address, _ = pip.economic.account.generate_account(pip.node.web3, 10 ** 18 * 1000)
        result = clients_consensus[0].duplicatesign.reportDuplicateSign(1, report_information, address)
        assert_code(result, 0)
        assert pip.get_accuverifiers_count(proposalinfo_text.get('ProposalID')) == [4, 1, 1, 0]
        assert pip.get_accuverifiers_count(proposalinfo_cancel.get('ProposalID')) == [4, 1, 1, 0]
        assert_code(pip.get_status_of_proposal(proposalinfo_cancel.get('ProposalID')), 3)
        assert_code(pip.get_status_of_proposal(proposalinfo_text.get('ProposalID')), 3)

    @pytest.mark.P2
    @allure.title('Cancel proposal and text proposal statistics function verification')
    def test_VS_EXT_010_VS_EXC_010(self, new_genesis_env, clients_consensus):
        self.update_setting(new_genesis_env, 500, 80, 10000, 4990, 10000, 4990)
        pip = clients_consensus[0].pip
        submittpandvote(clients_consensus[:2], 1, 1)
        submitcppandvote(clients_consensus[:2], [1, 1])
        proposalinfo_cancel = pip.get_effect_proposal_info_of_vote(pip.cfg.cancel_proposal)
        log.info('Get cancel proposal information {}'.format(proposalinfo_cancel))
        proposalinfo_text = pip.get_effect_proposal_info_of_vote(pip.cfg.text_proposal)
        log.info('Get text proposal information {}'.format(proposalinfo_text))
        wait_block_number(pip.node, 50)
        report_information = mock_duplicate_sign(2, pip.node.nodekey, pip.node.blsprikey, 45)
        log.info("Report information: {}".format(report_information))
        address, _ = pip.economic.account.generate_account(pip.node.web3, 10 ** 18 * 1000)
        result = clients_consensus[0].duplicatesign.reportDuplicateSign(2, report_information, address)
        assert_code(result, 0)
        wait_block_number(pip.node, self.get_block(proposalinfo_cancel, proposalinfo_text))
        assert pip.get_accuverifiers_count(proposalinfo_text.get('ProposalID')) == [4, 1, 0, 0]
        assert pip.get_accuverifiers_count(proposalinfo_cancel.get('ProposalID')) == [4, 1, 0, 0]
        assert_code(pip.get_status_of_proposal(proposalinfo_cancel.get('ProposalID')), 3)
        assert_code(pip.get_status_of_proposal(proposalinfo_text.get('ProposalID')), 3)

    @pytest.mark.P2
    @allure.title('Cancel proposal and text proposal statistics function verification')
    def test_VS_EXT_011_VS_EXC_011(self, new_genesis_env, clients_consensus):
        self.update_setting(new_genesis_env, 500, 80, 10000, 2499, 10000, 2499)
        pip = clients_consensus[0].pip
        submittpandvote(clients_consensus[:2], 2, 1)
        submitcppandvote(clients_consensus[:2], [2, 1])
        proposalinfo_cancel = pip.get_effect_proposal_info_of_vote(pip.cfg.cancel_proposal)
        log.info('Get cancel proposal information {}'.format(proposalinfo_cancel))
        proposalinfo_text = pip.get_effect_proposal_info_of_vote(pip.cfg.text_proposal)
        log.info('Get text proposal information {}'.format(proposalinfo_text))
        wait_block_number(pip.node, 50)
        report_information = mock_duplicate_sign(2, pip.node.nodekey, pip.node.blsprikey, 45)
        log.info("Report information: {}".format(report_information))
        address, _ = pip.economic.account.generate_account(pip.node.web3, 10 ** 18 * 1000)
        result = clients_consensus[0].duplicatesign.reportDuplicateSign(2, report_information, address)
        assert_code(result, 0)
        wait_block_number(pip.node, self.get_block(proposalinfo_cancel, proposalinfo_text))
        assert pip.get_accuverifiers_count(proposalinfo_text.get('ProposalID')) == [4, 1, 0, 0]
        assert pip.get_accuverifiers_count(proposalinfo_cancel.get('ProposalID')) == [4, 1, 0, 0]
        assert_code(pip.get_status_of_proposal(proposalinfo_cancel.get('ProposalID')), 2)
        assert_code(pip.get_status_of_proposal(proposalinfo_text.get('ProposalID')), 2)

    @pytest.mark.P2
    @allure.title('Cancel proposal and text proposal statistics function verification')
    def test_VS_EXT_012_VS_EXC_012(self, new_genesis_env, clients_consensus):
        self.update_setting(new_genesis_env, 500, 80, 10000, 2499, 10000, 2499)
        pip = clients_consensus[0].pip
        submittpandvote(clients_consensus[:2], 3, 1)
        submitcppandvote(clients_consensus[:2], [3, 1])
        proposalinfo_cancel = pip.get_effect_proposal_info_of_vote(pip.cfg.cancel_proposal)
        log.info('Get cancel proposal information {}'.format(proposalinfo_cancel))
        proposalinfo_text = pip.get_effect_proposal_info_of_vote(pip.cfg.text_proposal)
        log.info('Get text proposal information {}'.format(proposalinfo_text))
        wait_block_number(pip.node, 50)
        report_information = mock_duplicate_sign(2, pip.node.nodekey, pip.node.blsprikey, 45)
        log.info("Report information: {}".format(report_information))
        address, _ = pip.economic.account.generate_account(pip.node.web3, 10 ** 18 * 1000)
        result = clients_consensus[0].duplicatesign.reportDuplicateSign(2, report_information, address)
        assert_code(result, 0)
        wait_block_number(pip.node, self.get_block(proposalinfo_cancel, proposalinfo_text))
        assert pip.get_accuverifiers_count(proposalinfo_text.get('ProposalID')) == [4, 1, 0, 0]
        assert pip.get_accuverifiers_count(proposalinfo_cancel.get('ProposalID')) == [4, 1, 0, 0]
        assert_code(pip.get_status_of_proposal(proposalinfo_cancel.get('ProposalID')), 2)
        assert_code(pip.get_status_of_proposal(proposalinfo_text.get('ProposalID')), 2)

    @pytest.mark.P2
    @allure.title('Cancel proposal and text proposal statistics function verification')
    def test_VS_EXT_013_VS_EXC_013(self, new_genesis_env, clients_consensus):
        self.update_setting(new_genesis_env, 500, 120, 10000, 4999, 10000, 4999)
        pip = clients_consensus[0].pip
        pip_test = clients_consensus[1].pip
        submittpandvote(clients_consensus[:2], 1, 1)
        submitcppandvote(clients_consensus[:2], [1, 1], voting_rounds=3)
        proposalinfo_cancel = pip.get_effect_proposal_info_of_vote(pip.cfg.cancel_proposal)
        log.info('Get cancel proposal information {}'.format(proposalinfo_cancel))
        proposalinfo_text = pip.get_effect_proposal_info_of_vote(pip.cfg.text_proposal)
        log.info('Get text proposal information {}'.format(proposalinfo_text))
        pip.node.stop()
        wait_block_number(pip_test.node, self.get_block(proposalinfo_cancel, proposalinfo_text))
        assert pip_test.get_accuverifiers_count(proposalinfo_text.get('ProposalID')) == [4, 1, 0, 0]
        assert pip_test.get_accuverifiers_count(proposalinfo_cancel.get('ProposalID')) == [4, 1, 0, 0]
        assert_code(pip_test.get_status_of_proposal(proposalinfo_cancel.get('ProposalID')), 3)
        assert_code(pip_test.get_status_of_proposal(proposalinfo_text.get('ProposalID')), 3)

    @pytest.mark.P2
    @allure.title('Cancel proposal and text proposal statistics function verification')
    def test_VS_EXT_014_VS_EXC_014(self, new_genesis_env, clients_consensus):
        self.update_setting(new_genesis_env, 500, 120, 10000, 4999, 10000, 4999)
        pip = clients_consensus[0].pip
        pip_test = clients_consensus[1].pip
        submittpandvote(clients_consensus[:2], 2, 1)
        submitcppandvote(clients_consensus[:2], [2, 1], voting_rounds=3)
        proposalinfo_cancel = pip.get_effect_proposal_info_of_vote(pip.cfg.cancel_proposal)
        log.info('Get cancel proposal information {}'.format(proposalinfo_cancel))
        proposalinfo_text = pip.get_effect_proposal_info_of_vote(pip.cfg.text_proposal)
        log.info('Get text proposal information {}'.format(proposalinfo_text))
        pip.node.stop()
        wait_block_number(pip_test.node, self.get_block(proposalinfo_cancel, proposalinfo_text))
        assert pip_test.get_accuverifiers_count(proposalinfo_text.get('ProposalID')) == [4, 1, 0, 0]
        assert pip_test.get_accuverifiers_count(proposalinfo_cancel.get('ProposalID')) == [4, 1, 0, 0]
        assert_code(pip_test.get_status_of_proposal(proposalinfo_cancel.get('ProposalID')), 3)
        assert_code(pip_test.get_status_of_proposal(proposalinfo_text.get('ProposalID')), 3)

    @pytest.mark.P2
    @allure.title('Cancel proposal and text proposal statistics function verification')
    def test_VS_EXT_015_VS_EXC_015(self, new_genesis_env, clients_consensus):
        self.update_setting(new_genesis_env, 500, 120, 10000, 4999, 10000, 4999)
        pip = clients_consensus[0].pip
        pip_test = clients_consensus[1].pip
        submittpandvote(clients_consensus[:2], 3, 1)
        submitcppandvote(clients_consensus[:2], [3, 1], voting_rounds=3)
        proposalinfo_cancel = pip.get_effect_proposal_info_of_vote(pip.cfg.cancel_proposal)
        log.info('Get cancel proposal information {}'.format(proposalinfo_cancel))
        proposalinfo_text = pip.get_effect_proposal_info_of_vote(pip.cfg.text_proposal)
        log.info('Get text proposal information {}'.format(proposalinfo_text))
        pip.node.stop()
        wait_block_number(pip_test.node, self.get_block(proposalinfo_cancel, proposalinfo_text))
        assert pip_test.get_accuverifiers_count(proposalinfo_text.get('ProposalID')) == [4, 1, 0, 0]
        assert pip_test.get_accuverifiers_count(proposalinfo_cancel.get('ProposalID')) == [4, 1, 0, 0]
        assert_code(pip_test.get_status_of_proposal(proposalinfo_cancel.get('ProposalID')), 3)
        assert_code(pip_test.get_status_of_proposal(proposalinfo_text.get('ProposalID')), 3)

    @pytest.mark.P2
    @allure.title('Cancel proposal and text proposal statistics function verification')
    def test_VS_EXT_016_VS_EXC_016(self, new_genesis_env, clients_consensus):
        self.update_setting(new_genesis_env, 500, 80, 10000, 4999, 10000, 4999)
        pip = clients_consensus[0].pip
        pip_test = clients_consensus[1].pip
        submittpandvote(clients_consensus[:2], 1, 1)
        submitcppandvote(clients_consensus[:2], [1, 1], voting_rounds=2)
        proposalinfo_cancel = pip.get_effect_proposal_info_of_vote(pip.cfg.cancel_proposal)
        log.info('Get cancel proposal information {}'.format(proposalinfo_cancel))
        proposalinfo_text = pip.get_effect_proposal_info_of_vote(pip.cfg.text_proposal)
        log.info('Get text proposal information {}'.format(proposalinfo_text))
        wait_block_number(pip_test.node, self.get_block(proposalinfo_cancel, proposalinfo_text))
        pip.node.stop()
        assert pip_test.get_accuverifiers_count(proposalinfo_text.get('ProposalID')) == [4, 2, 0, 0]
        assert pip_test.get_accuverifiers_count(proposalinfo_cancel.get('ProposalID')) == [4, 2, 0, 0]
        assert_code(pip_test.get_status_of_proposal(proposalinfo_cancel.get('ProposalID')), 2)
        assert_code(pip_test.get_status_of_proposal(proposalinfo_text.get('ProposalID')), 2)


class TestVotingStatisticsPP:
    def update_setting_param(self, new_genesis_env, *args):
        genesis = from_dict(data_class=Genesis, data=new_genesis_env.genesis_config)
        if len(args) == 1:
            genesis.economicModel.gov.paramProposalVoteDurationSeconds = args[0]
        elif len(args) == 3:
            genesis.economicModel.gov.paramProposalVoteDurationSeconds = args[0]
            genesis.economicModel.gov.paramProposalSupportRate = args[1]
            genesis.economicModel.gov.paramProposalVoteRate = args[2]
        else:
            raise ValueError('args error')
        new_genesis_env.set_genesis(genesis.to_dict())
        new_genesis_env.deploy_all()

    def assert_proposal_result(self, pip, proposalinfo, tally_result):
        assert pip.get_yeas_of_proposal(proposalinfo.get('ProposalID')) == tally_result[0]
        assert pip.get_nays_of_proposal(proposalinfo.get('ProposalID')) == tally_result[1]
        assert pip.get_abstentions_of_proposal(proposalinfo.get('ProposalID')) == tally_result[2]
        assert pip.get_status_of_proposal(proposalinfo.get('ProposalID')) == tally_result[3]

    def assert_accuverifiers_count(self, pip, proposalinfo, accuverifiers_result):
        result = pip.get_accuverifiers_count(proposalinfo.get('ProposalID'))
        log.info('Get proposal vote infomation {}'.format(result))
        assert result == accuverifiers_result

    @pytest.mark.P1
    @allure.title('Parammeter proposal statistics function verification')
    def test_VS_EP_004(self, new_genesis_env, clients_consensus, clients_noconsensus):
        self.update_setting_param(new_genesis_env, 0)
        submitppandvote(clients_consensus[0:-1], 1, 2, 3)
        proposalinfo = clients_consensus[0].pip.get_effect_proposal_info_of_vote(clients_consensus[0].pip.cfg.param_proposal)
        log.info('Param proposal information {}'.format(proposalinfo))
        createstaking(clients_noconsensus[:3])
        self.assert_accuverifiers_count(clients_consensus[0].pip, proposalinfo, [4, 1, 1, 1])

    @pytest.mark.P1
    @allure.title('Parammeter proposal statistics function verification')
    def test_VS_EP_005(self, new_genesis_env, clients_consensus, clients_noconsensus):
        self.update_setting_param(new_genesis_env, 160)
        pip = clients_consensus[0].pip
        submitppandvote(clients_consensus[:2], 1, 2)
        proposalinfo = pip.get_effect_proposal_info_of_vote(pip.cfg.param_proposal)
        log.info('Param proposal info {}'.format(proposalinfo))
        log.info('{}'.format(clients_consensus[:2]))
        createstaking(clients_noconsensus[:2])
        pip.economic.wait_settlement_blocknum(pip.node)
        self.assert_accuverifiers_count(pip, proposalinfo, [6, 1, 1, 0])

        result = proposal_vote(clients_noconsensus[0].pip, vote_option=pip.cfg.vote_option_Abstentions)
        assert_code(result, 0)
        log.info('{}'.format(clients_consensus[2]))
        createstaking(clients_noconsensus[2])
        wait_block_number(pip.node, proposalinfo.get('EndVotingBlock'))

        self.assert_accuverifiers_count(pip, proposalinfo, [6, 1, 1, 1])

    @pytest.mark.P2
    @allure.title('Parammeter proposal statistics function verification')
    def test_VS_EP_006(self, new_genesis_env, clients_consensus, clients_noconsensus):
        self.update_setting_param(new_genesis_env, 320)
        submitppandvote(clients_consensus[:1], 1)
        proposalinfo = clients_consensus[0].pip.get_effect_proposal_info_of_vote(clients_consensus[0].pip.cfg.param_proposal)
        log.info('Param proposal info {}'.format(proposalinfo))
        createstaking(clients_noconsensus[0])
        clients_consensus[0].pip.economic.wait_settlement_blocknum(clients_consensus[0].pip.node)
        self.assert_accuverifiers_count(clients_consensus[0].pip, proposalinfo, [5, 1, 0, 0])

        result = proposal_vote(clients_consensus[1].pip, vote_option=clients_consensus[0].pip.cfg.vote_option_nays)
        assert_code(result, 0)
        createstaking(clients_noconsensus[1])
        clients_consensus[0].pip.economic.wait_settlement_blocknum(clients_consensus[0].pip.node)
        self.assert_accuverifiers_count(clients_consensus[0].pip, proposalinfo, [6, 1, 1, 0])

        result = proposal_vote(clients_consensus[2].pip, vote_option=clients_consensus[0].pip.cfg.vote_option_Abstentions)
        assert_code(result, 0)
        createstaking(clients_noconsensus[2])
        wait_block_number(clients_consensus[0].pip.node, proposalinfo.get('EndVotingBlock'))

        self.assert_accuverifiers_count(clients_consensus[0].pip, proposalinfo, [6, 1, 1, 1])

    @pytest.mark.P0
    @allure.title('Parammeter proposal statistics function verification')
    def test_VS_EP_007_VS_EP_003(self, new_genesis_env, clients_consensus):
        self.update_setting_param(new_genesis_env, 0, 10000, 4900)
        submitppandvote(clients_consensus[:2], 1, 1)
        proposalinfo = clients_consensus[0].pip.get_effect_proposal_info_of_vote(clients_consensus[0].pip.cfg.param_proposal)
        log.info('Get param proposal information {}'.format(proposalinfo))
        result = clients_consensus[0].staking.withdrew_staking(clients_consensus[0].node.staking_address)
        log.info('Node {} withdrew staking result {}'.format(clients_consensus[0].node.node_id, result))
        assert_code(result, 0)
        result = clients_consensus[0].pip.pip.getTallyResult(proposalinfo.get('ProposalID'))
        log.info('Before endvoting block, get Tally resul of the parameter proposal result : {}'.format(result))
        assert_code(result, 302030)
        wait_block_number(clients_consensus[0].node, proposalinfo.get('EndVotingBlock'))
        self.assert_proposal_result(clients_consensus[0].pip, proposalinfo, [2, 0, 0, 2])

    @pytest.mark.P2
    @allure.title('Parammeter proposal statistics function verification')
    def test_VS_EP_008(self, new_genesis_env, clients_consensus):
        self.update_setting_param(new_genesis_env, 0, 5000, 5000)
        submitppandvote(clients_consensus[:2], 1, 1)
        pip = clients_consensus[0].pip
        proposalinfo = pip.get_effect_proposal_info_of_vote(pip.cfg.param_proposal)
        log.info('Get param proposal information {}'.format(proposalinfo))
        wait_block_number(pip.node, proposalinfo.get('EndVotingBlock'))
        report_information = mock_duplicate_sign(1, pip.node.nodekey, pip.node.blsprikey,
                                                 proposalinfo.get('EndVotingBlock') - 10)
        log.info("Report information: {}".format(report_information))
        address, _ = pip.economic.account.generate_account(pip.node.web3, 10**18 * 1000)
        result = clients_consensus[0].duplicatesign.reportDuplicateSign(1, report_information, address)
        assert_code(result, 0)
        assert clients_consensus[1].pip.get_yeas_of_proposal(proposalinfo.get('ProposalID')) == 2

    @pytest.mark.P2
    @allure.title('Parammeter proposal statistics function verification')
    def test_VS_EP_009(self, new_genesis_env, clients_consensus):
        self.update_setting_param(new_genesis_env, 0, 5000, 5000)
        submitppandvote(clients_consensus[:2], 1, 1)
        pip = clients_consensus[0].pip
        proposalinfo = pip.get_effect_proposal_info_of_vote(pip.cfg.param_proposal)
        log.info('Get param proposal information {}'.format(proposalinfo))
        wait_block_number(pip.node, 80)
        report_information = mock_duplicate_sign(1, pip.node.nodekey, pip.node.blsprikey, 70)
        log.info("Report information: {}".format(report_information))
        address, _ = pip.economic.account.generate_account(pip.node.web3, 10**18 * 1000)
        result = clients_consensus[0].duplicatesign.reportDuplicateSign(1, report_information, address)
        assert_code(result, 0)
        wait_block_number(pip.node, proposalinfo.get('EndVotingBlock'))
        assert clients_consensus[1].pip.get_yeas_of_proposal(proposalinfo.get('ProposalID')) == 1

    @pytest.mark.P2
    @allure.title('Parammeter proposal statistics function verification')
    def test_VS_EP_010(self, new_genesis_env, clients_consensus):
        self.update_setting_param(new_genesis_env, 0, 5000, 5000)
        submitppandvote(clients_consensus[:2], 2, 2)
        pip = clients_consensus[0].pip
        proposalinfo = pip.get_effect_proposal_info_of_vote(pip.cfg.param_proposal)
        log.info('Get param proposal information {}'.format(proposalinfo))
        wait_block_number(pip.node, 80)
        report_information = mock_duplicate_sign(1, pip.node.nodekey, pip.node.blsprikey, 70)
        log.info("Report information: {}".format(report_information))
        address, _ = pip.economic.account.generate_account(pip.node.web3, 10**18 * 1000)
        result = clients_consensus[0].duplicatesign.reportDuplicateSign(1, report_information, address)
        assert_code(result, 0)
        wait_block_number(pip.node, proposalinfo.get('EndVotingBlock'))
        assert clients_consensus[1].pip.get_nays_of_proposal(proposalinfo.get('ProposalID')) == 1

    @pytest.mark.P2
    @allure.title('Parammeter proposal statistics function verification')
    def test_VS_EP_011(self, new_genesis_env, clients_consensus):
        self.update_setting_param(new_genesis_env, 0, 5000, 5000)
        submitppandvote(clients_consensus[:2], 3, 3)
        pip = clients_consensus[0].pip
        proposalinfo = pip.get_effect_proposal_info_of_vote(pip.cfg.param_proposal)
        log.info('Get param proposal information {}'.format(proposalinfo))
        wait_block_number(pip.node, 80)
        report_information = mock_duplicate_sign(1, pip.node.nodekey, pip.node.blsprikey, 70)
        log.info("Report information: {}".format(report_information))
        address, _ = pip.economic.account.generate_account(pip.node.web3, 10**18 * 1000)
        result = clients_consensus[0].duplicatesign.reportDuplicateSign(1, report_information, address)
        assert_code(result, 0)
        wait_block_number(pip.node, proposalinfo.get('EndVotingBlock'))
        assert clients_consensus[1].pip.get_abstentions_of_proposal(proposalinfo.get('ProposalID')) == 1

    @pytest.mark.P2
    @allure.title('Parammeter proposal statistics function verification')
    def test_VS_EP_012(self, new_genesis_env, clients_consensus):
        self.update_setting_param(new_genesis_env, 0, 5000, 5000)
        submitppandvote(clients_consensus[:2], 1, 1)
        pip = clients_consensus[0].pip
        proposalinfo = pip.get_effect_proposal_info_of_vote(pip.cfg.param_proposal)
        log.info('Get param proposal information {}'.format(proposalinfo))
        wait_block_number(pip.node, proposalinfo.get('EndVotingBlock'))
        report_information = mock_duplicate_sign(2, pip.node.nodekey, pip.node.blsprikey,
                                                 proposalinfo.get('EndVotingBlock') - 10)
        log.info("Report information: {}".format(report_information))
        address, _ = pip.economic.account.generate_account(pip.node.web3, 10**18 * 1000)
        result = clients_consensus[0].duplicatesign.reportDuplicateSign(2, report_information, address)
        assert_code(result, 0)
        assert clients_consensus[1].pip.get_yeas_of_proposal(proposalinfo.get('ProposalID')) == 2

    @pytest.mark.P2
    @allure.title('Parammeter proposal statistics function verification')
    def test_VS_EP_013(self, new_genesis_env, clients_consensus):
        self.update_setting_param(new_genesis_env, 0, 5000, 5000)
        submitppandvote(clients_consensus[:2], 1, 1)
        pip = clients_consensus[0].pip
        proposalinfo = pip.get_effect_proposal_info_of_vote(pip.cfg.param_proposal)
        log.info('Get param proposal information {}'.format(proposalinfo))
        wait_block_number(pip.node, 80)
        report_information = mock_duplicate_sign(2, pip.node.nodekey, pip.node.blsprikey, 70)
        log.info("Report information: {}".format(report_information))
        address, _ = pip.economic.account.generate_account(pip.node.web3, 10**18 * 1000)
        result = clients_consensus[0].duplicatesign.reportDuplicateSign(2, report_information, address)
        assert_code(result, 0)
        wait_block_number(pip.node, proposalinfo.get('EndVotingBlock'))
        assert clients_consensus[1].pip.get_yeas_of_proposal(proposalinfo.get('ProposalID')) == 1

    @pytest.mark.P2
    @allure.title('Parammeter proposal statistics function verification')
    def test_VS_EP_014(self, new_genesis_env, clients_consensus):
        self.update_setting_param(new_genesis_env, 0, 5000, 5000)
        submitppandvote(clients_consensus[:2], 2, 2)
        pip = clients_consensus[0].pip
        proposalinfo = pip.get_effect_proposal_info_of_vote(pip.cfg.param_proposal)
        log.info('Get param proposal information {}'.format(proposalinfo))
        wait_block_number(pip.node, 80)
        report_information = mock_duplicate_sign(2, pip.node.nodekey, pip.node.blsprikey, 70)
        log.info("Report information: {}".format(report_information))
        address, _ = pip.economic.account.generate_account(pip.node.web3, 10**18 * 1000)
        result = clients_consensus[0].duplicatesign.reportDuplicateSign(2, report_information, address)
        assert_code(result, 0)
        wait_block_number(pip.node, proposalinfo.get('EndVotingBlock'))
        assert clients_consensus[1].pip.get_nays_of_proposal(proposalinfo.get('ProposalID')) == 1

    @pytest.mark.P2
    @allure.title('Parammeter proposal statistics function verification')
    def test_VS_EP_015(self, new_genesis_env, clients_consensus):
        self.update_setting_param(new_genesis_env, 0, 5000, 5000)
        submitppandvote(clients_consensus[:2], 3, 3)
        pip = clients_consensus[0].pip
        proposalinfo = pip.get_effect_proposal_info_of_vote(pip.cfg.param_proposal)
        log.info('Get param proposal information {}'.format(proposalinfo))
        wait_block_number(pip.node, 80)
        report_information = mock_duplicate_sign(2, pip.node.nodekey, pip.node.blsprikey, 70)
        log.info("Report information: {}".format(report_information))
        address, _ = pip.economic.account.generate_account(pip.node.web3, 10**18 * 1000)
        result = clients_consensus[0].duplicatesign.reportDuplicateSign(2, report_information, address)
        assert_code(result, 0)
        wait_block_number(pip.node, proposalinfo.get('EndVotingBlock'))
        assert clients_consensus[1].pip.get_abstentions_of_proposal(proposalinfo.get('ProposalID')) == 1

    @pytest.mark.P2
    @allure.title('Parammeter proposal statistics function verification')
    def test_VS_EP_016(self, new_genesis_env, clients_consensus):
        self.update_setting_param(new_genesis_env, 0, 10000, 2490)
        submitppandvote(clients_consensus[:2], 1, 1)
        pip = clients_consensus[0].pip
        proposalinfo = pip.get_effect_proposal_info_of_vote(pip.cfg.param_proposal)
        log.info('Get param proposal information {}'.format(proposalinfo))
        assert pip.get_accuverifiers_count(proposalinfo.get('ProposalID')) == [4, 2, 0, 0]
        log.info('Stop the node {}'.format(pip.node.node_id))
        pip.node.stop()
        wait_block_number(clients_consensus[1].node, proposalinfo.get('EndVotingBlock'))
        self.assert_proposal_result(clients_consensus[1].pip, proposalinfo, [1, 0, 0, 2])

    @pytest.mark.P2
    @allure.title('Parammeter proposal statistics function verification')
    def test_VS_EP_017(self, new_genesis_env, clients_consensus):
        self.update_setting_param(new_genesis_env, 0, 10000, 2490)
        submitppandvote(clients_consensus[:2], 2, 1)
        pip = clients_consensus[0].pip
        proposalinfo = pip.get_effect_proposal_info_of_vote(pip.cfg.param_proposal)
        log.info('Get param proposal information {}'.format(proposalinfo))
        assert pip.get_accuverifiers_count(proposalinfo.get('ProposalID')) == [4, 1, 1, 0]
        log.info('Stop the node {}'.format(pip.node.node_id))
        pip.node.stop()
        wait_block_number(clients_consensus[1].node, proposalinfo.get('EndVotingBlock'))
        self.assert_proposal_result(clients_consensus[1].pip, proposalinfo, [1, 0, 0, 2])

    @pytest.mark.P2
    @allure.title('Parammeter proposal statistics function verification')
    def test_VS_EP_018(self, new_genesis_env, clients_consensus):
        self.update_setting_param(new_genesis_env, 0, 10000, 2490)
        submitppandvote(clients_consensus[:2], 3, 1)
        pip = clients_consensus[0].pip
        proposalinfo = pip.get_effect_proposal_info_of_vote(pip.cfg.param_proposal)
        log.info('Get param proposal information {}'.format(proposalinfo))
        assert pip.get_accuverifiers_count(proposalinfo.get('ProposalID')) == [4, 1, 0, 1]
        log.info('Stop the node {}'.format(pip.node.node_id))
        pip.node.stop()
        wait_block_number(clients_consensus[1].node, proposalinfo.get('EndVotingBlock'))
        self.assert_proposal_result(clients_consensus[1].pip, proposalinfo, [1, 0, 0, 2])

    @pytest.mark.P2
    @allure.title('Parammeter proposal statistics function verification')
    def test_VS_EP_019(self, new_genesis_env, clients_consensus):
        self.update_setting_param(new_genesis_env, 0, 9900, 2500)
        submitppandvote(clients_consensus[:3], 1, 2, 3)
        pip = clients_consensus[0].pip
        proposalinfo = pip.get_effect_proposal_info_of_vote(pip.cfg.param_proposal)
        log.info('Get param proposal information {}'.format(proposalinfo))
        assert pip.get_accuverifiers_count(proposalinfo.get('ProposalID')) == [4, 1, 1, 1]
        wait_block_number(clients_consensus[1].node, proposalinfo.get('EndVotingBlock'))
        self.assert_proposal_result(clients_consensus[1].pip, proposalinfo, [1, 1, 1, 3])
        log.info('Stop the node {}'.format(pip.node.node_id))
        pip.node.stop()
        wait_block_number(clients_consensus[1].node, proposalinfo.get('EndVotingBlock')
                          + clients_consensus[1].economic.consensus_size)
        self.assert_proposal_result(clients_consensus[1].pip, proposalinfo, [1, 1, 1, 3])
