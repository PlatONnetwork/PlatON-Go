from common.log import log
from tests.lib.utils import assert_code, wait_block_number, get_blockhash, get_the_dynamic_parameter_gas_fee
from dacite import from_dict
from tests.lib import Genesis
from common.key import mock_duplicate_sign
from tests.govern.test_voting_statistics import submitppandvote, submitcvpandvote, submitvpandvote, submittpandvote, submitcppandvote

import time
import math
import rlp
import pytest
import allure
from tests.govern.conftest import version_proposal_vote

cancelby = "0x0000000000000000000000000000000000000000000000000000000000000000"


class TestgetProposal:
    @pytest.mark.P0
    @allure.title('Interface getProposal function verification--cancel proposal')
    def test_GP_IF_001(self, submit_cancel_param):
        pip = submit_cancel_param
        proposalinfo = pip.get_effect_proposal_info_of_vote(pip.cfg.cancel_proposal)
        log.info('Cancel proposal information : {}'.format(proposalinfo))
        result = pip.pip.getProposal(proposalinfo.get('ProposalID'))
        log.info('Interface getProposal result : {}'.format(result))
        assert result.get('Ret').get('Proposer') == pip.node.node_id
        assert result.get('Ret').get('ProposalType') == pip.cfg.cancel_proposal
        assert result.get('Ret').get('PIPID') == proposalinfo.get('PIPID')
        assert result.get('Ret').get('SubmitBlock') == proposalinfo.get('SubmitBlock')
        assert result.get('Ret').get('EndVotingBlock') == proposalinfo.get('EndVotingBlock')

    @pytest.mark.P0
    @allure.title('Interface getProposal function verification--parammeter proposal')
    def test_GP_IF_002(self, submit_param):
        pip = submit_param
        proposalinfo = pip.get_effect_proposal_info_of_vote(pip.cfg.param_proposal)
        log.info('Param proposal information : {}'.format(proposalinfo))
        result = pip.pip.getProposal(proposalinfo.get('ProposalID'))
        log.info('Interface getProposal result : {}'.format(result))
        assert result.get('Ret').get('Proposer') == pip.node.node_id
        assert result.get('Ret').get('ProposalType') == pip.cfg.param_proposal
        assert result.get('Ret').get('PIPID') == proposalinfo.get('PIPID')
        assert result.get('Ret').get('SubmitBlock') == proposalinfo.get('SubmitBlock')
        assert result.get('Ret').get('EndVotingBlock') == proposalinfo.get('EndVotingBlock')

    @pytest.mark.compatibility
    @pytest.mark.P0
    @allure.title('Interface getProposal function verification--version proposal')
    def test_PR_IN_001_002(self, no_vp_proposal):
        pip = no_vp_proposal
        pip_id = str(time.time())
        result = pip.submitVersion(pip.node.node_id, pip_id, pip.cfg.version8, 3, pip.node.staking_address,
                                   transaction_cfg=pip.cfg.transaction_cfg)
        log.info('Submit version proposal result : {}'.format(result))
        assert_code(result, 0)
        proposalinfo_version = pip.get_effect_proposal_info_of_vote(pip.cfg.version_proposal)
        log.info('Get version proposal information : {}'.format(proposalinfo_version))
        pip_id_cancel = str(time.time())
        result = pip.submitCancel(pip.node.node_id, pip_id_cancel, 1, proposalinfo_version.get('ProposalID'),
                                  pip.node.staking_address, transaction_cfg=pip .cfg.transaction_cfg)
        log.info('Submit cancel proposal result : {}'.format(result))
        assert_code(result, 0)

        proposalinfo_cancel = pip.get_effect_proposal_info_of_vote(pip.cfg.cancel_proposal)
        log.info('Get cancel proposal information : {}'.format(proposalinfo_cancel))

        result_version = pip.pip.getProposal(proposalinfo_version.get('ProposalID'))
        log.info('Interface getProposal-version result : {}'.format(result_version))

        result_cancel = pip.pip.getProposal(proposalinfo_cancel.get('ProposalID'))
        log.info('Interface getProposal-cancel result : {}'.format(result_cancel))

        assert result_version.get('Ret').get('Proposer') == pip.node.node_id
        assert result_version.get('Ret').get('ProposalType') == pip.cfg.version_proposal
        assert result_version.get('Ret').get('PIPID') == pip_id
        assert result_version.get('Ret').get('SubmitBlock') == proposalinfo_version.get('SubmitBlock')
        caculated_endvotingblock = math.ceil(proposalinfo_version.get('SubmitBlock') / pip.economic.consensus_size
                                             + 3) * pip.economic.consensus_size - 20
        assert result_version.get('Ret').get('EndVotingBlock') == caculated_endvotingblock

        assert result_cancel.get('Ret').get('Proposer') == pip.node.node_id
        assert result_cancel.get('Ret').get('ProposalType') == pip.cfg.cancel_proposal
        assert result_cancel.get('Ret').get('PIPID') == pip_id_cancel
        assert result_cancel.get('Ret').get('SubmitBlock') == proposalinfo_cancel.get('SubmitBlock')
        caculated_endvotingblock = math.ceil(proposalinfo_cancel.get('SubmitBlock') / pip.economic.consensus_size) * \
            pip.economic.consensus_size + 20
        assert result_cancel.get('Ret').get('EndVotingBlock') == caculated_endvotingblock

    @pytest.mark.P0
    @allure.title('Interface getProposal function verification--text proposal')
    def test_PR_IN_003(self, client_verifier):
        pip = client_verifier.pip
        pip_id = str(time.time())
        result = pip.submitText(pip.node.node_id, pip_id, pip.node.staking_address,
                                transaction_cfg=pip.cfg.transaction_cfg)
        log.info('Submit text proposal result : {}'.format(result))
        assert_code(result, 0)
        proposalinfo_text = pip.get_effect_proposal_info_of_vote(pip.cfg.text_proposal)
        log.info('Get text proposal information : {}'.format(proposalinfo_text))

        result_text = pip.pip.getProposal(proposalinfo_text.get('ProposalID'))
        log.info('Interface getProposal-text result : {}'.format(result_text))

        assert result_text.get('Ret').get('Proposer') == pip.node.node_id
        assert result_text.get('Ret').get('ProposalType') == pip.cfg.text_proposal
        assert result_text.get('Ret').get('PIPID') == pip_id
        assert result_text.get('Ret').get('SubmitBlock') == proposalinfo_text.get('SubmitBlock')
        log.info(pip.economic.tp_vote_settlement_wheel)
        caculated_endvotingblock = math.ceil(proposalinfo_text.get('SubmitBlock') / pip.economic.consensus_size
                                             + pip.economic.tp_vote_settlement_wheel) * pip.economic.consensus_size - 20
        assert result_text.get('Ret').get('EndVotingBlock') == caculated_endvotingblock

    @pytest.mark.P1
    @allure.title('Interface getProposal function verification--ineffective proposal id')
    def test_PR_IN_004(self, client_noconsensus):
        pip = client_noconsensus.pip
        result = pip.pip.getProposal('0xa89162be0bd0d081c50a5160f412c4926b3ae9ea96cf792935564357ddd11111')
        log.info('Interface getProposal-version result : {}'.format(result))
        assert_code(result, 302006)


class TestgetTallyResult:
    @pytest.mark.P0
    @allure.title('Interface getTallyResult function verification--cancel version proposal')
    def test_TR_IN_002_TR_IN_003(self, no_vp_proposal, clients_verifier):
        pip = no_vp_proposal
        submitcvpandvote(clients_verifier, 1, 1, 1, 2)
        proposalinfo_cancel = pip.get_effect_proposal_info_of_vote(pip.cfg.cancel_proposal)
        log.info('Cancel proposal information {}'.format(proposalinfo_cancel))
        proposalinfo_version = pip.get_effect_proposal_info_of_vote(pip.cfg.version_proposal)
        log.info('Version proposal information {}'.format(proposalinfo_version))
        wait_block_number(pip.node, proposalinfo_version.get('EndVotingBlock'))
        assert pip.get_canceledby_of_proposal(proposalinfo_cancel.get('ProposalID')) == cancelby
        assert pip.get_status_of_proposal(proposalinfo_cancel.get('ProposalID')) == 2
        assert pip.get_yeas_of_proposal(proposalinfo_cancel.get('ProposalID')) == 3
        assert pip.get_nays_of_proposal(proposalinfo_cancel.get('ProposalID')) == 1
        assert pip.get_abstentions_of_proposal(proposalinfo_cancel.get('ProposalID')) == 0
        assert pip.get_accu_verifiers_of_proposal(proposalinfo_cancel.get('ProposalID')) == len(clients_verifier)

        assert pip.get_canceledby_of_proposal(proposalinfo_version.get('ProposalID')) == proposalinfo_cancel.get('ProposalID')
        assert pip.get_status_of_proposal(proposalinfo_version.get('ProposalID')) == 6
        assert pip.get_yeas_of_proposal(proposalinfo_version.get('ProposalID')) == 0
        assert pip.get_nays_of_proposal(proposalinfo_version.get('ProposalID')) == 0
        assert pip.get_abstentions_of_proposal(proposalinfo_version.get('ProposalID')) == 0
        assert pip.get_accu_verifiers_of_proposal(proposalinfo_version.get('ProposalID')) == len(clients_verifier)

    @pytest.mark.P0
    @pytest.mark.compatibility
    @allure.title('Interface getTallyResult function verification--cancel version proposal')
    def test_TR_IN_001(self, no_vp_proposal, clients_verifier):
        pip = no_vp_proposal
        submitcvpandvote(clients_verifier, 1, 2, 3, 3)
        proposalinfo_cancel = pip.get_effect_proposal_info_of_vote(pip.cfg.cancel_proposal)
        log.info('Cancel proposal information {}'.format(proposalinfo_cancel))
        proposalinfo_version = pip.get_effect_proposal_info_of_vote(pip.cfg.version_proposal)
        log.info('Version proposal information {}'.format(proposalinfo_version))
        wait_block_number(pip.node, proposalinfo_version.get('EndVotingBlock'))
        assert pip.get_canceledby_of_proposal(proposalinfo_cancel.get('ProposalID')) == cancelby
        assert pip.get_status_of_proposal(proposalinfo_cancel.get('ProposalID')) == 3
        assert pip.get_yeas_of_proposal(proposalinfo_cancel.get('ProposalID')) == 1
        assert pip.get_nays_of_proposal(proposalinfo_cancel.get('ProposalID')) == 1
        assert pip.get_abstentions_of_proposal(proposalinfo_cancel.get('ProposalID')) == 2
        assert pip.get_accu_verifiers_of_proposal(proposalinfo_cancel.get('ProposalID')) == len(clients_verifier)

        assert pip.get_canceledby_of_proposal(proposalinfo_version.get('ProposalID')) == cancelby
        assert pip.get_status_of_proposal(proposalinfo_version.get('ProposalID')) == 3
        assert pip.get_yeas_of_proposal(proposalinfo_version.get('ProposalID')) == 0
        assert pip.get_nays_of_proposal(proposalinfo_version.get('ProposalID')) == 0
        assert pip.get_abstentions_of_proposal(proposalinfo_version.get('ProposalID')) == 0
        assert pip.get_accu_verifiers_of_proposal(proposalinfo_version.get('ProposalID')) == len(clients_verifier)

    @pytest.mark.P0
    @allure.title('Interface getTallyResult function verification--parammeter proposal')
    def test_TR_IN_010_005(self, new_genesis_env, clients_consensus):
        genesis = from_dict(data_class=Genesis, data=new_genesis_env.genesis_config)
        genesis.economicModel.gov.paramProposalVoteDurationSeconds = 0
        new_genesis_env.set_genesis(genesis.to_dict())
        new_genesis_env.deploy_all()
        pip = clients_consensus[0].pip
        submitppandvote(clients_consensus[0:-1], 1, 2, 3)
        proposalinfo = clients_consensus[0].pip.get_effect_proposal_info_of_vote(clients_consensus[0].pip.cfg.param_proposal)
        log.info('Param proposal information {}'.format(proposalinfo))
        log.info('listparam {}'.format(clients_consensus[0].pip.pip.listGovernParam()))
        result = pip.pip.getTallyResult(proposalinfo.get('ProposalID'))
        log.info('Interface getTallyResult info : {}'.format(result))
        assert_code(result, 302030)
        wait_block_number(pip.node, proposalinfo.get('EndVotingBlock'))

        assert pip.get_canceledby_of_proposal(proposalinfo.get('ProposalID')) == cancelby
        assert pip.get_status_of_proposal(proposalinfo.get('ProposalID')) == 3
        assert pip.get_yeas_of_proposal(proposalinfo.get('ProposalID')) == 1
        assert pip.get_nays_of_proposal(proposalinfo.get('ProposalID')) == 1
        assert pip.get_abstentions_of_proposal(proposalinfo.get('ProposalID')) == 1
        assert pip.get_accu_verifiers_of_proposal(proposalinfo.get('ProposalID')) == len(clients_consensus)

    @pytest.mark.P0
    @allure.title('Interface getTallyResult function verification--cancel parammeter proposal')
    def test_TR_IN_011_TR_IN_012(self, new_genesis_env, clients_consensus):
        genesis = from_dict(data_class=Genesis, data=new_genesis_env.genesis_config)
        genesis.economicModel.gov.paramProposalVoteDurationSeconds = 0
        new_genesis_env.set_genesis(genesis.to_dict())
        new_genesis_env.deploy_all()
        pip = clients_consensus[0].pip
        submitcppandvote(clients_consensus, [1, 1, 1, 3])
        proposalinfo_param = pip.get_effect_proposal_info_of_vote(pip.cfg.param_proposal)
        log.info('Param proposal information {}'.format(proposalinfo_param))
        proposalinfo_cancel = pip.get_effect_proposal_info_of_vote(pip.cfg.cancel_proposal)
        log.info('Cancel proposal information {}'.format(proposalinfo_cancel))

        wait_block_number(pip.node, proposalinfo_cancel.get('EndVotingBlock'))

        assert pip.get_canceledby_of_proposal(proposalinfo_cancel.get('ProposalID')) == cancelby
        assert pip.get_canceledby_of_proposal(proposalinfo_param.get('ProposalID')) == proposalinfo_cancel.get('ProposalID')

        assert pip.get_status_of_proposal(proposalinfo_cancel.get('ProposalID')) == 2
        assert pip.get_yeas_of_proposal(proposalinfo_cancel.get('ProposalID')) == 3
        assert pip.get_nays_of_proposal(proposalinfo_cancel.get('ProposalID')) == 0
        assert pip.get_abstentions_of_proposal(proposalinfo_cancel.get('ProposalID')) == 1
        assert pip.get_accu_verifiers_of_proposal(proposalinfo_cancel.get('ProposalID')) == len(clients_consensus)

        assert pip.get_status_of_proposal(proposalinfo_param.get('ProposalID')) == 6
        assert pip.get_yeas_of_proposal(proposalinfo_param.get('ProposalID')) == 0
        assert pip.get_nays_of_proposal(proposalinfo_param.get('ProposalID')) == 0
        assert pip.get_abstentions_of_proposal(proposalinfo_param.get('ProposalID')) == 0
        assert pip.get_accu_verifiers_of_proposal(proposalinfo_param.get('ProposalID')) == len(clients_consensus)

    @pytest.mark.P1
    @allure.title('Interface getTallyResult function verification--ineffective proposal id')
    def test_TR_IN_006(self, client_verifier):
        pip = client_verifier.pip
        result = pip.pip.getTallyResult('0x9992d1f843fe8f376884d871f87605dda02da0722fd6b350bbf683518f73f111')
        log.info('Ineffective proposalID, interface getTallyResult return : {}'.format(result))
        assert_code(result, 302030)


class TestgetAccuVerifiersCount:
    @pytest.mark.P0
    @allure.title('Interface getTallyResult function verification--ineffective proposal id')
    def test_AC_IN_018_to_025(self, new_genesis_env, clients_consensus):
        genesis = from_dict(data_class=Genesis, data=new_genesis_env.genesis_config)
        genesis.economicModel.gov.paramProposalVoteDurationSeconds = 0
        new_genesis_env.set_genesis(genesis.to_dict())
        new_genesis_env.deploy_all()
        pip = clients_consensus[0].pip
        pip_test = clients_consensus[-1].pip
        result = pip.submitParam(pip.node.node_id, str(time.time()), 'slashing', 'slashBlocksReward', '999',
                                 pip.node.staking_address, transaction_cfg=pip.cfg.transaction_cfg)
        log.info('Node submit param proposal result : {}'.format(result))
        assert_code(result, 0)
        proposalinfo_param = pip.get_effect_proposal_info_of_vote(pip.cfg.param_proposal)
        log.info('Get param proposal information : {}'.format(proposalinfo_param))
        result = pip.submitCancel(pip.node.node_id, str(time.time()), 3, proposalinfo_param.get('ProposalID'),
                                  pip.node.staking_address, transaction_cfg=pip.cfg.transaction_cfg)
        log.info('Node submit cancel proposal result : {}'.format(result))
        proposalinfo_cancel = pip.get_effect_proposal_info_of_vote(pip.cfg.cancel_proposal)
        log.info('Get cancel proposal information : {}'.format(proposalinfo_cancel))
        for index in range(3):
            client = clients_consensus[index]
            result = client.pip.vote(client.node.node_id, proposalinfo_param.get('ProposalID'), index + 1,
                                     client.node.staking_address, transaction_cfg=pip.cfg.transaction_cfg)
            log.info('Node {} vote param proposal result : {}'.format(client.node.node_id, result))
            assert_code(result, 0)
            result = client.pip.vote(client.node.node_id, proposalinfo_cancel.get('ProposalID'), index + 1,
                                     client.node.staking_address, transaction_cfg=pip.cfg.transaction_cfg)
            log.info('Node {} vote cancel proposal result : {}'.format(client.node.node_id, result))
            assert_code(result, 0)
        assert pip.get_accuverifiers_count(proposalinfo_param.get('ProposalID')) == [4, 1, 1, 1]
        assert pip.get_accuverifiers_count(proposalinfo_cancel.get('ProposalID')) == [4, 1, 1, 1]
        log.info('Stop the node {}'.format(pip.node.node_id))
        pip.node.stop()
        pip_test.economic.wait_consensus_blocknum(pip_test.node, 2)
        assert pip_test.get_accuverifiers_count(proposalinfo_param.get('ProposalID')) == [4, 0, 1, 1]
        assert pip_test.get_accuverifiers_count(proposalinfo_cancel.get('ProposalID')) == [4, 0, 1, 1]

        report_information = mock_duplicate_sign(1, clients_consensus[1].node.nodekey,
                                                 clients_consensus[1].node.blsprikey, 41)
        log.info("Report information: {}".format(report_information))
        address, _ = pip_test.economic.account.generate_account(pip_test.node.web3, 10 ** 18 * 1000)
        result = clients_consensus[-1].duplicatesign.reportDuplicateSign(1, report_information, address)
        log.info('Node duplicate block result : {}'.format(result))
        assert_code(result, 0)
        time.sleep(2)
        assert pip_test.get_accuverifiers_count(proposalinfo_param.get('ProposalID')) == [4, 0, 0, 1]
        assert pip_test.get_accuverifiers_count(proposalinfo_cancel.get('ProposalID')) == [4, 0, 0, 1]

        report_information = mock_duplicate_sign(2, clients_consensus[2].node.nodekey,
                                                 clients_consensus[2].node.blsprikey, 41)
        log.info("Report information: {}".format(report_information))
        address, _ = pip_test.economic.account.generate_account(pip_test.node.web3, 10 ** 18 * 1000)
        result = clients_consensus[-1].duplicatesign.reportDuplicateSign(2, report_information, address)
        log.info('Node duplicate block result : {}'.format(result))
        assert_code(result, 0)
        assert pip_test.get_accuverifiers_count(proposalinfo_param.get('ProposalID')) == [4, 0, 0, 0]
        assert pip_test.get_accuverifiers_count(proposalinfo_cancel.get('ProposalID')) == [4, 0, 0, 0]

    @pytest.mark.compatibility
    @pytest.mark.P0
    @allure.title('Interface getAccuVerifiersCount function verification')
    def test_AC_IN_001_002_004_to_006_012_to_014(self, new_genesis_env, clients_consensus):
        new_genesis_env.deploy_all()
        pip = clients_consensus[-1].pip
        result = pip.submitVersion(pip.node.node_id, str(time.time()), pip.cfg.version5, 5, pip.node.staking_address,
                                   transaction_cfg=pip.cfg.transaction_cfg)
        log.info('Submit version proposal result : {}'.format(result))
        assert_code(result, 0)
        proposalinfo_version = pip.get_effect_proposal_info_of_vote()
        log.info('Get version proposal information : {}'.format(proposalinfo_version))

        result = pip.submitCancel(pip.node.node_id, str(time.time()), 4, proposalinfo_version.get('ProposalID'),
                                  pip.node.staking_address, transaction_cfg=pip.cfg.transaction_cfg)
        log.info('Submit cancel proposal result : {}'.format(result))
        assert_code(result, 0)
        proposalinfo_cancel = pip.get_effect_proposal_info_of_vote(pip.cfg.cancel_proposal)
        log.info('Get cancel proposal information : {}'.format(proposalinfo_version))

        for index in range(3):
            client = clients_consensus[index]
            result = version_proposal_vote(client.pip)
            assert_code(result, 0)
            result = client.pip.vote(client.node.node_id, proposalinfo_cancel.get('ProposalID'), index + 1,
                                     client.node.staking_address, transaction_cfg=pip.cfg.transaction_cfg)
            log.info('Node {} vote cancel proposal result : {}'.format(client.node.node_id, result))
            assert_code(result, 0)

        assert pip.get_accuverifiers_count(proposalinfo_version.get('ProposalID')) == [4, 3, 0, 0]
        assert pip.get_accuverifiers_count(proposalinfo_cancel.get('ProposalID')) == [4, 1, 1, 1]
        log.info('Stop the node {}'.format(clients_consensus[0].node.node_id))
        clients_consensus[0].node.stop()
        pip.economic.wait_consensus_blocknum(pip.node, 2)
        log.info(pip.node.debug.getWaitSlashingNodeList())
        log.info(pip.pip.listGovernParam())
        log.info(clients_consensus[1].ppos.getCandidateInfo(pip.node.node_id, pip.node.staking_address))
        assert pip.get_accuverifiers_count(proposalinfo_version.get('ProposalID')) == [4, 2, 0, 0]
        assert pip.get_accuverifiers_count(proposalinfo_cancel.get('ProposalID')) == [4, 0, 1, 1]

        report_information = mock_duplicate_sign(1, clients_consensus[1].node.nodekey,
                                                 clients_consensus[1].node.blsprikey, 41)
        log.info("Report information: {}".format(report_information))
        address, _ = pip.economic.account.generate_account(pip.node.web3, 10 ** 18 * 1000)
        result = clients_consensus[-1].duplicatesign.reportDuplicateSign(1, report_information, address)
        log.info('Node duplicate block result : {}'.format(result))
        assert_code(result, 0)
        assert pip.get_accuverifiers_count(proposalinfo_version.get('ProposalID')) == [4, 1, 0, 0]
        assert pip.get_accuverifiers_count(proposalinfo_cancel.get('ProposalID')) == [4, 0, 0, 1]

        report_information = mock_duplicate_sign(2, clients_consensus[2].node.nodekey,
                                                 clients_consensus[2].node.blsprikey, 41)
        log.info("Report information: {}".format(report_information))
        address, _ = pip.economic.account.generate_account(pip.node.web3, 10 ** 18 * 1000)
        result = clients_consensus[-1].duplicatesign.reportDuplicateSign(2, report_information, address)
        log.info('Node duplicate block result : {}'.format(result))
        assert_code(result, 0)
        assert pip.get_accuverifiers_count(proposalinfo_version.get('ProposalID')) == [4, 0, 0, 0]
        assert pip.get_accuverifiers_count(proposalinfo_cancel.get('ProposalID')) == [4, 0, 0, 0]

    @pytest.mark.P0
    @allure.title('Interface getAccuVerifiersCount function verification')
    def test_AC_IN_003_008_010(self, new_genesis_env, clients_consensus):
        genesis = from_dict(data_class=Genesis, data=new_genesis_env.genesis_config)
        genesis.economicModel.gov.textProposalVoteDurationSeconds = 120
        new_genesis_env.set_genesis(genesis.to_dict())
        new_genesis_env.deploy_all()
        pip = clients_consensus[-1].pip
        submittpandvote(clients_consensus, 1, 2, 3, 1)
        proposalinfo = pip.get_effect_proposal_info_of_vote(pip.cfg.text_proposal)
        log.info('Get text proposal information : {}'.format(proposalinfo))
        assert pip.get_accuverifiers_count(proposalinfo.get('ProposalID')) == [4, 2, 1, 1]
        log.info('Stop the node {}'.format(clients_consensus[0].node.node_id))
        clients_consensus[0].node.stop()
        pip.economic.wait_consensus_blocknum(pip.node, 2)
        assert pip.get_accuverifiers_count(proposalinfo.get('ProposalID')) == [4, 1, 1, 1]

        report_information = mock_duplicate_sign(1, clients_consensus[1].node.nodekey,
                                                 clients_consensus[1].node.blsprikey, 41)
        log.info("Report information: {}".format(report_information))
        address, _ = pip.economic.account.generate_account(pip.node.web3, 10 ** 18 * 1000)
        result = clients_consensus[-1].duplicatesign.reportDuplicateSign(1, report_information, address)
        log.info('Node duplicate block result : {}'.format(result))
        assert_code(result, 0)
        assert pip.get_accuverifiers_count(proposalinfo.get('ProposalID')) == [4, 1, 0, 1]

        report_information = mock_duplicate_sign(2, clients_consensus[2].node.nodekey,
                                                 clients_consensus[2].node.blsprikey, 41)
        log.info("Report information: {}".format(report_information))
        address, _ = pip.economic.account.generate_account(pip.node.web3, 10 ** 18 * 1000)
        result = clients_consensus[-1].duplicatesign.reportDuplicateSign(2, report_information, address)
        log.info('Node duplicate block result : {}'.format(result))
        assert_code(result, 0)
        assert pip.get_accuverifiers_count(proposalinfo.get('ProposalID')) == [4, 1, 0, 0]

    @pytest.mark.P2
    @allure.title('Interface getAccuVerifiersCount function verification')
    def test_AC_IN_016_to_018(self, client_verifier):
        pip = client_verifier.pip
        result = pip.submitText(pip.node.node_id, str(time.time()), pip.node.staking_address,
                                transaction_cfg=pip.cfg.transaction_cfg)
        log.info('Submit text proposal result : {}'.format(result))
        assert_code(result, 0)
        proposalinfo = pip.get_effect_proposal_info_of_vote(pip.cfg.text_proposal)
        log.info('Get text proposal information : {}'.format(proposalinfo))
        result = pip.pip.getAccuVerifiersCount('0x0c04f578466ead2208dbb15b927ecb27041881e8c16c17cd0db6b3df422e1111',
                                               block_hash=get_blockhash(pip.node))
        log.info('Interface getAccuVerifiersCount result : {}'.format(result))
        assert_code(result, 302006)
        log.info('{}'.format(get_blockhash(pip.node)))

        result = pip.pip.getAccuVerifiersCount(proposalinfo.get('ProposalID'), block_hash='')
        log.info('Interface getAccuVerifiersCount result : {}'.format(result))
        assert_code(result, 3)

        result = pip.pip.getAccuVerifiersCount(proposalinfo.get('ProposalID'),
                                               block_hash='0x5941605fe43ab32fbaf9c6e08dc0970eae50efb7da4248a9a8941f0e50711111')
        log.info('Interface getAccuVerifiersCount result : {}'.format(result))
        assert_code(result, 0)


class TestListGovernParam:
    def get_govern_param(self, client, module=None):
        result = client.pip.pip.listGovernParam(module)
        log.info('Interface listGovernParam result {}'.format(result))
        assert_code(result, 0)
        resultinfo = result.get('Ret')
        module = []
        name = []
        for param in resultinfo:
            module.append(param.get('ParamItem').get('Module'))
            name.append(param.get('ParamItem').get('Name'))
        return name, module

    @pytest.mark.P0
    @allure.title('Interface listGovernParam function verification')
    def test_IN_LG_001(self, client_noconsensus):
        name, module = self.get_govern_param(client_noconsensus)
        assert set(name) == {'maxValidators', 'unStakeFreezeDuration', 'operatingThreshold', 'slashBlocksReward',
                             'stakeThreshold', 'maxBlockGasLimit', 'duplicateSignReportReward', 'maxEvidenceAge',
                             'slashFractionDuplicateSign', 'zeroProduceCumulativeTime', 'zeroProduceNumberThreshold',
                             'rewardPerMaxChangeRange', 'rewardPerChangeInterval', 'increaseIssuanceRatio'}
        assert set(module) == {'block', 'slashing', 'staking', 'reward'}

    @pytest.mark.P2
    @allure.title('Interface listGovernParam function verification')
    def test_IN_LG_002(self, client_noconsensus):
        name, module = self.get_govern_param(client_noconsensus, 'staking')
        assert set(name) == {'maxValidators', 'unStakeFreezeDuration', 'operatingThreshold', 'stakeThreshold',
                             'rewardPerMaxChangeRange', 'rewardPerChangeInterval'}
        assert set(module) == {'staking'}

    @pytest.mark.P2
    @allure.title('Interface listGovernParam function verification')
    def test_IN_LG_003(self, client_noconsensus):
        name, module = self.get_govern_param(client_noconsensus, 'slashing')
        assert set(name) == {'slashBlocksReward', 'duplicateSignReportReward', 'maxEvidenceAge',
                             'slashFractionDuplicateSign', 'zeroProduceCumulativeTime', 'zeroProduceNumberThreshold'}
        assert set(module) == {'slashing'}

    @pytest.mark.P2
    @allure.title('Interface listGovernParam function verification')
    def test_IN_LG_004(self, client_noconsensus):
        name, module = self.get_govern_param(client_noconsensus, 'block')
        assert set(name) == {'maxBlockGasLimit'}
        assert set(module) == {'block'}

        name, module = self.get_govern_param(client_noconsensus, 'reward')
        assert set(name) == {'increaseIssuanceRatio'}
        assert set(module) == {'reward'}


    @pytest.mark.P2
    @allure.title('Interface listGovernParam function verification')
    def test_IN_LG_005(self, client_noconsensus):
        result = client_noconsensus.pip.pip.listGovernParam('txpool')
        log.info('Interface listGovernParam result {}'.format(result))
        assert_code(result, 2)
        assert result.get('Ret') == "Object not found"


class TestGetGovernParam:
    @pytest.mark.P0
    @allure.title('Interface getGovernParamValue function verification')
    def test_IN_GG_001(self, client_noconsensus):
        client_noconsensus.economic.env.deploy_all()
        genesis = from_dict(data_class=Genesis, data=client_noconsensus.economic.env.genesis_config)
        pip = client_noconsensus.pip.pip
        result = pip.getGovernParamValue('slashing', 'slashBlocksReward')
        log.info('Interface getGovernParamValue result : {}'.format(result))
        assert genesis.economicModel.slashing.slashBlocksReward == int(result.get('Ret'))

        result = pip.getGovernParamValue('slashing', 'maxEvidenceAge')
        log.info('Interface getGovernParamValue result : {}'.format(result))
        assert genesis.economicModel.slashing.maxEvidenceAge == int(result.get('Ret'))

        result = pip.getGovernParamValue('slashing', 'slashFractionDuplicateSign')
        log.info('Interface getGovernParamValue result : {}'.format(result))
        assert genesis.economicModel.slashing.slashFractionDuplicateSign == int(result.get('Ret'))

        result = pip.getGovernParamValue('slashing', 'duplicateSignReportReward')
        log.info('Interface getGovernParamValue result : {}'.format(result))
        assert genesis.economicModel.slashing.duplicateSignReportReward == int(result.get('Ret'))

        result = pip.getGovernParamValue('staking', 'stakeThreshold')
        log.info('Interface getGovernParamValue result : {}'.format(result))
        assert genesis.economicModel.staking.stakeThreshold == int(result.get('Ret'))

        result = pip.getGovernParamValue('staking', 'operatingThreshold')
        log.info('Interface getGovernParamValue result : {}'.format(result))
        assert genesis.economicModel.staking.operatingThreshold == int(result.get('Ret'))

        result = pip.getGovernParamValue('staking', 'unStakeFreezeDuration')
        log.info('Interface getGovernParamValue result : {}'.format(result))
        assert genesis.economicModel.staking.unStakeFreezeDuration == int(result.get('Ret'))

        result = pip.getGovernParamValue('staking', 'maxValidators')
        log.info('Interface getGovernParamValue result : {}'.format(result))
        assert genesis.economicModel.staking.maxValidators == int(result.get('Ret'))

        result = pip.getGovernParamValue('block', 'maxBlockGasLimit')
        log.info('Interface getGovernParamValue result : {}'.format(result))
        assert_code(result, 0)

    @pytest.mark.P2
    @allure.title('Interface getGovernParamValue function verification')
    def test_IN_GG_002(self, client_noconsensus):
        pip = client_noconsensus.pip.pip
        result = pip.getGovernParamValue('Staking', 'maxValidators')
        assert_code(result, 302031)
        pip = client_noconsensus.pip.pip
        result = pip.getGovernParamValue('Slashing', 'slashBlocksReward')
        assert_code(result, 302031)
        pip = client_noconsensus.pip.pip
        result = pip.getGovernParamValue('Block', 'maxBlockGasLimit')
        assert_code(result, 302031)

    @pytest.mark.P2
    @allure.title('Interface getGovernParamValue function verification')
    def test_IN_GG_003(self, client_noconsensus):
        pip = client_noconsensus.pip.pip
        result = pip.getGovernParamValue('staking', 'MaxValidators')
        assert_code(result, 302031)
        pip = client_noconsensus.pip.pip
        result = pip.getGovernParamValue('slashing', 'SlashBlocksReward')
        assert_code(result, 302031)
        pip = client_noconsensus.pip.pip
        result = pip.getGovernParamValue('block', 'MaxValidators')
        assert_code(result, 302031)


class TestGetActiveVersion:
    @pytest.mark.compatibility
    @pytest.mark.P0
    @allure.title('Interface getActiveVersion function verification')
    def test_AV_IN_001(self, no_vp_proposal):
        assert_code(no_vp_proposal.chain_version, no_vp_proposal.cfg.version0)

    @pytest.mark.P0
    @allure.title('Interface getActiveVersion function verification')
    def test_AV_IN_002_003(self, clients_verifier):
        pip = clients_verifier[0].pip
        submitvpandvote(clients_verifier)
        proposalinfo = pip.get_effect_proposal_info_of_vote()
        log.info('Version proposal information : {}'.format(proposalinfo))
        wait_block_number(pip.node, proposalinfo.get('EndVotingBlock'))
        assert_code(pip.get_status_of_proposal(proposalinfo.get('ProposalID')), 4)
        assert_code(pip.chain_version, pip.cfg.version0)
        wait_block_number(pip.node, proposalinfo.get('ActiveBlock'))
        assert_code(pip.get_status_of_proposal(proposalinfo.get('ProposalID')), 5)
        assert_code(pip.chain_version, pip.cfg.version5)


class TestListProposal:
    @pytest.mark.P1
    @allure.title('Interface listProposal function verification')
    def test_LP_IN_001_002(self, no_vp_proposal):
        pip = no_vp_proposal
        pip_id = str(time.time())
        result = pip.submitParam(pip.node.node_id, pip_id, 'slashing', 'slashBlocksReward', '456',
                                 pip.node.staking_address, transaction_cfg=pip.cfg.transaction_cfg)
        log.info('Submit param proposal result : {}'.format(result))
        assert_code(result, 0)
        proposalinfo = pip.get_effect_proposal_info_of_vote(pip.cfg.param_proposal)
        log.info('Get param proposal information : {}'.format(proposalinfo))
        assert proposalinfo.get('Proposer') == pip.node.node_id
        assert proposalinfo.get('ProposalType') == pip.cfg.param_proposal
        log.info('{}'.format(pip.economic.pp_vote_settlement_wheel))
        calculated_endvotingblock = math.ceil(proposalinfo.get('SubmitBlock') / pip.economic.settlement_size
                                              + pip.economic.pp_vote_settlement_wheel) * pip.economic.settlement_size
        assert proposalinfo.get('EndVotingBlock') == calculated_endvotingblock

        pip_id = str(time.time())
        result = pip.submitCancel(pip.node.node_id, pip_id, 1, proposalinfo.get('ProposalID'),
                                  pip.node.staking_address, transaction_cfg=pip.cfg.transaction_cfg)
        log.info('Submit cancel proposal result : {}'.format(result))
        assert_code(result, 0)
        proposalinfo_cancel = pip.get_effect_proposal_info_of_vote(pip.cfg.cancel_proposal)
        log.info('Get cancel proposal information : {}'.format(proposalinfo_cancel))

        assert proposalinfo_cancel.get('Proposer') == pip.node.node_id
        assert proposalinfo_cancel.get('ProposalType') == pip.cfg.cancel_proposal
        log.info('{}'.format(pip.economic.pp_vote_settlement_wheel))
        calculated_endvotingblock = math.ceil(proposalinfo_cancel.get('SubmitBlock') / pip.economic.consensus_size
                                              + 1) * pip.economic.consensus_size - 20
        assert proposalinfo_cancel.get('EndVotingBlock') == calculated_endvotingblock

    @pytest.mark.P1
    @allure.title('Interface listProposal function verification')
    def test_LP_IN_003(self, client_consensus):
        client_consensus.economic.env.deploy_all()
        result = client_consensus.pip.pip.listProposal()
        log.info('There is no proposal, interface listProposal return : {}'.format(result))
        assert_code(result, 2)
        assert result.get('Ret') == "Object not found"


class TestGasUse:
    def get_balance(self, pip):
        balance = pip.node.eth.getBalance(pip.node.staking_address)
        log.info('address balance : {}'.format(balance))
        return balance

    @pytest.mark.P2
    @allure.title('Verify gas --submittext and vote')
    def test_TP_GA_001(self, client_verifier):
        pip = client_verifier.pip
        pip_id = str(time.time())
        data = rlp.encode([rlp.encode(int(2000)), rlp.encode(bytes.fromhex(pip.node.node_id)), rlp.encode(pip_id)])
        balance_before = self.get_balance(pip)
        result = pip.submitText(pip.node.node_id, pip_id, pip.node.staking_address,
                                transaction_cfg=pip.cfg.transaction_cfg)
        log.info('Submit text proposal result : {}'.format(result))
        assert_code(result, 0)
        proposalinfo = pip.get_effect_proposal_info_of_vote(pip.cfg.text_proposal)
        log.info('Get text proposal information : {}'.format(proposalinfo))
        balance_after = self.get_balance(pip)
        gas = get_the_dynamic_parameter_gas_fee(data)
        log.info('Calculated gas : {}'.format(gas))
        assert_code(balance_before - balance_after, (gas + 350000) * pip.cfg.transaction_cfg.get('gasPrice'))

        proposal_id = proposalinfo.get('ProposalID')[2:]
        version_sign = pip.node.program_version_sign[2:]
        data = rlp.encode([rlp.encode(int(2003)), rlp.encode(bytes.fromhex(pip.node.node_id)),
                           rlp.encode(bytes.fromhex(proposal_id)),
                           rlp.encode(pip.cfg.vote_option_yeas), rlp.encode(int(pip.node.program_version)),
                           rlp.encode(bytes.fromhex(version_sign))])
        result = pip.vote(pip.node.node_id, proposalinfo.get('ProposalID'), pip.cfg.vote_option_yeas,
                          pip.node.staking_address, transaction_cfg=pip.cfg.transaction_cfg)
        log.info('Vote reuslt : {}'.format(result))
        assert_code(result, 0)
        balance_after_vote = pip.node.eth.getBalance(pip.node.staking_address)
        log.info('After vote text proposal, the address balance : {}'.format(balance_after_vote))
        gas = get_the_dynamic_parameter_gas_fee(data)
        log.info('Calculated gas : {}'.format(gas))
        assert_code(balance_after - balance_after_vote, (gas + 32000) * pip.cfg.transaction_cfg.get('gasPrice'))

    @pytest.mark.P2
    @allure.title('Verify gas --submitversion')
    def test_VP_GA_001(self, no_vp_proposal):
        pip = no_vp_proposal
        pip_id = str(time.time())
        balance_before = self.get_balance(pip)
        result = pip.submitVersion(pip.node.node_id, pip_id, pip.cfg.version5, 1, pip.node.staking_address,
                                   transaction_cfg=pip.cfg.transaction_cfg)
        log.info('Submit version proposal result : {}'.format(result))
        assert_code(result, 0)
        balance_after = self.get_balance(pip)
        data = rlp.encode([rlp.encode(int(2001)), rlp.encode(bytes.fromhex(pip.node.node_id)), rlp.encode(pip_id),
                           rlp.encode(int(pip.cfg.version5)), rlp.encode(int(1))])
        gas = get_the_dynamic_parameter_gas_fee(data)
        log.info('Calculated gas : {}'.format(gas))
        assert_code(balance_before - balance_after, (gas + 480000) * pip.cfg.transaction_cfg.get('gasPrice'))

    @pytest.mark.P2
    @allure.title('Verify gas --submitparam_and_cancel')
    def test_PP_GA_001_CP_GA_001(self, no_vp_proposal):
        pip = no_vp_proposal
        pip_id = str(time.time())
        balance_before = self.get_balance(pip)
        result = pip.submitParam(pip.node.node_id, pip_id, 'slashing', 'slashBlocksReward', '123',
                                 pip.node.staking_address, transaction_cfg=pip.cfg.transaction_cfg)
        log.info('Submit param proposal result : {}'.format(result))
        assert_code(result, 0)
        proposalinfor_param = pip.get_effect_proposal_info_of_vote(pip.cfg.param_proposal)
        log.info('Get param proposal information : {}'.format(proposalinfor_param))
        balance_after = self.get_balance(pip)
        data = rlp.encode([rlp.encode(int(2002)), rlp.encode(bytes.fromhex(pip.node.node_id)),
                           rlp.encode(pip_id), rlp.encode('slashing'), rlp.encode('slashBlocksReward'),
                           rlp.encode('123')])
        gas = get_the_dynamic_parameter_gas_fee(data)
        log.info('Calculated gas : {}'.format(gas))
        assert_code(balance_before - balance_after, (gas + 530000) * pip.cfg.transaction_cfg.get('gasPrice'))

        pip_id = str(time.time())
        result = pip.submitCancel(pip.node.node_id, pip_id, 1, proposalinfor_param.get('ProposalID'),
                                  pip.node.staking_address, transaction_cfg=pip.cfg.transaction_cfg)
        log.info('Submit cancel proposal result : {}'.format(result))
        assert_code(result, 0)
        assert_code(balance_before - balance_after, (gas + 530000) * pip.cfg.transaction_cfg.get('gasPrice'))
        balance_after_cancel = pip.node.eth.getBalance(pip.node.staking_address)
        log.info('After submitting cancel proposal, the address balance : {}'.format(balance_after_cancel))
        tobe_canceled_proposal_id = proposalinfor_param.get('ProposalID')[2:]
        data = rlp.encode([rlp.encode(int(2005)), rlp.encode(bytes.fromhex(pip.node.node_id)), rlp.encode(pip_id),
                           rlp.encode(int(1)), rlp.encode(bytes.fromhex(tobe_canceled_proposal_id))])
        gas = get_the_dynamic_parameter_gas_fee(data)
        log.info('Calculated gas : {}'.format(gas))
        assert_code(balance_after - balance_after_cancel, (gas + 530000) * pip.cfg.transaction_cfg.get('gasPrice'))

    @pytest.mark.P2
    @allure.title('Verify gas --declare version')
    def test_declareversion(self, client_verifier):
        pip = client_verifier.pip
        balance_before = self.get_balance(pip)
        result = pip.declareVersion(pip.node.node_id, pip.node.staking_address,
                                    transaction_cfg=pip.cfg.transaction_cfg)
        log.info('Declare version result : {}'.format(result))
        assert_code(result, 0)
        version_sign = pip.node.program_version_sign[2:]
        data = rlp.encode([rlp.encode(int(2004)), rlp.encode(bytes.fromhex(pip.node.node_id)),
                           rlp.encode(int(pip.node.program_version)), rlp.encode(bytes.fromhex(version_sign))])
        gas = get_the_dynamic_parameter_gas_fee(data)
        log.info('Calculated gas : {}'.format(gas))
        balance_after = self.get_balance(pip)
        assert_code(balance_before - balance_after, (gas + 33000) * pip.cfg.transaction_cfg.get('gasPrice'))

if __name__ == '__main__':
    pytest.main(['./tests/govern/', '-s', '-q', '--alluredir', './report/report'])
