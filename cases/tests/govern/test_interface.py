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
from tests.govern.conftest import version_proposal_vote


class TestgetProposal():
    @pytest.mark.P0
    def test_GP_IF_001(self, submit_cancel_param):
        pip_obj = submit_cancel_param
        proposalinfo = pip_obj.get_effect_proposal_info_of_vote(pip_obj.cfg.cancel_proposal)
        log.info('Cancel proposal information : {}'.format(proposalinfo))
        result = pip_obj.pip.getProposal(proposalinfo.get('ProposalID'))
        log.info('Interface getProposal result : {}'.format(result))
        assert result.get('Ret').get('Proposer') == pip_obj.node.node_id
        assert result.get('Ret').get('ProposalType') == pip_obj.cfg.cancel_proposal
        assert result.get('Ret').get('PIPID') == proposalinfo.get('PIPID')
        assert result.get('Ret').get('SubmitBlock') == proposalinfo.get('SubmitBlock')
        assert result.get('Ret').get('EndVotingBlock') == proposalinfo.get('EndVotingBlock')

    @pytest.mark.P0
    def test_GP_IF_002(self, submit_param):
        pip_obj = submit_param
        proposalinfo = pip_obj.get_effect_proposal_info_of_vote(pip_obj.cfg.param_proposal)
        log.info('Param proposal information : {}'.format(proposalinfo))
        result = pip_obj.pip.getProposal(proposalinfo.get('ProposalID'))
        log.info('Interface getProposal result : {}'.format(result))
        assert result.get('Ret').get('Proposer') == pip_obj.node.node_id
        assert result.get('Ret').get('ProposalType') == pip_obj.cfg.param_proposal
        assert result.get('Ret').get('PIPID') == proposalinfo.get('PIPID')
        assert result.get('Ret').get('SubmitBlock') == proposalinfo.get('SubmitBlock')
        assert result.get('Ret').get('EndVotingBlock') == proposalinfo.get('EndVotingBlock')

    @pytest.mark.compatibility
    @pytest.mark.P0
    def test_PR_IN_001_002(self, no_vp_proposal):
        pip_obj = no_vp_proposal
        pip_id = str(time.time())
        result = pip_obj.submitVersion(pip_obj.node.node_id, pip_id, pip_obj.cfg.version8, 3, pip_obj.node.staking_address,
                                       transaction_cfg=pip_obj.cfg.transaction_cfg)
        log.info('Submit version proposal result : {}'.format(result))
        assert_code(result, 0)
        proposalinfo_version = pip_obj.get_effect_proposal_info_of_vote(pip_obj.cfg.version_proposal)
        log.info('Get version proposal information : {}'.format(proposalinfo_version))
        pip_id_cancel = str(time.time())
        result = pip_obj.submitCancel(pip_obj.node.node_id, pip_id_cancel, 1, proposalinfo_version.get('ProposalID'),
                                      pip_obj.node.staking_address, transaction_cfg=pip_obj .cfg.transaction_cfg)
        log.info('Submit cancel proposal result : {}'.format(result))
        assert_code(result, 0)

        proposalinfo_cancel = pip_obj.get_effect_proposal_info_of_vote(pip_obj.cfg.cancel_proposal)
        log.info('Get cancel proposal information : {}'.format(proposalinfo_cancel))

        result_version = pip_obj.pip.getProposal(proposalinfo_version.get('ProposalID'))
        log.info('Interface getProposal-version result : {}'.format(result_version))

        result_cancel = pip_obj.pip.getProposal(proposalinfo_cancel.get('ProposalID'))
        log.info('Interface getProposal-cancel result : {}'.format(result_cancel))

        assert result_version.get('Ret').get('Proposer') == pip_obj.node.node_id
        assert result_version.get('Ret').get('ProposalType') == pip_obj.cfg.version_proposal
        assert result_version.get('Ret').get('PIPID') == pip_id
        assert result_version.get('Ret').get('SubmitBlock') == proposalinfo_version.get('SubmitBlock')
        caculated_endvotingblock = math.ceil(proposalinfo_version.get('SubmitBlock') / pip_obj.economic.consensus_size +
                                             3) * pip_obj.economic.consensus_size - 20
        assert result_version.get('Ret').get('EndVotingBlock') == caculated_endvotingblock

        assert result_cancel.get('Ret').get('Proposer') == pip_obj.node.node_id
        assert result_cancel.get('Ret').get('ProposalType') == pip_obj.cfg.cancel_proposal
        assert result_cancel.get('Ret').get('PIPID') == pip_id_cancel
        assert result_cancel.get('Ret').get('SubmitBlock') == proposalinfo_cancel.get('SubmitBlock')
        caculated_endvotingblock = math.ceil(proposalinfo_cancel.get('SubmitBlock') / pip_obj.economic.consensus_size) * \
            pip_obj.economic.consensus_size + 20
        assert result_cancel.get('Ret').get('EndVotingBlock') == caculated_endvotingblock

    @pytest.mark.P0
    def test_PR_IN_003(self, client_verifier_obj):
        pip_obj = client_verifier_obj.pip
        pip_id = str(time.time())
        result = pip_obj.submitText(pip_obj.node.node_id, pip_id, pip_obj.node.staking_address,
                                    transaction_cfg=pip_obj.cfg.transaction_cfg)
        log.info('Submit text proposal result : {}'.format(result))
        assert_code(result, 0)
        proposalinfo_text = pip_obj.get_effect_proposal_info_of_vote(pip_obj.cfg.text_proposal)
        log.info('Get text proposal information : {}'.format(proposalinfo_text))

        result_text = pip_obj.pip.getProposal(proposalinfo_text.get('ProposalID'))
        log.info('Interface getProposal-text result : {}'.format(result_text))

        assert result_text.get('Ret').get('Proposer') == pip_obj.node.node_id
        assert result_text.get('Ret').get('ProposalType') == pip_obj.cfg.text_proposal
        assert result_text.get('Ret').get('PIPID') == pip_id
        assert result_text.get('Ret').get('SubmitBlock') == proposalinfo_text.get('SubmitBlock')
        log.info(pip_obj.economic.tp_vote_settlement_wheel)
        caculated_endvotingblock = math.ceil(proposalinfo_text.get('SubmitBlock') / pip_obj.economic.consensus_size +
                                             pip_obj.economic.tp_vote_settlement_wheel) * pip_obj.economic.consensus_size - 20
        assert result_text.get('Ret').get('EndVotingBlock') == caculated_endvotingblock

    @pytest.mark.P1
    def test_PR_IN_004(self, client_noconsensus_obj):
        pip_obj = client_noconsensus_obj.pip
        result = pip_obj.pip.getProposal('0xa89162be0bd0d081c50a5160f412c4926b3ae9ea96cf792935564357ddd11111')
        log.info('Interface getProposal-version result : {}'.format(result))
        assert_code(result, 302006)


class TestgetTallyResult():
    @pytest.mark.P0
    def test_TR_IN_002_TR_IN_003(self, no_vp_proposal, client_verifier_obj_list):
        pip_obj = no_vp_proposal
        submitcvpandvote(client_verifier_obj_list, 1, 1, 1, 2)
        proposalinfo_cancel = pip_obj.get_effect_proposal_info_of_vote(pip_obj.cfg.cancel_proposal)
        log.info('Cancel proposal information {}'.format(proposalinfo_cancel))
        proposalinfo_version = pip_obj.get_effect_proposal_info_of_vote(pip_obj.cfg.version_proposal)
        log.info('Version proposal information {}'.format(proposalinfo_version))
        wait_block_number(pip_obj.node, proposalinfo_version.get('EndVotingBlock'))
        assert pip_obj.get_canceledby_of_proposal(proposalinfo_cancel.get('ProposalID')) == \
            "0x0000000000000000000000000000000000000000000000000000000000000000"
        assert pip_obj.get_status_of_proposal(proposalinfo_cancel.get('ProposalID')) == 2
        assert pip_obj.get_yeas_of_proposal(proposalinfo_cancel.get('ProposalID')) == 3
        assert pip_obj.get_nays_of_proposal(proposalinfo_cancel.get('ProposalID')) == 1
        assert pip_obj.get_abstentions_of_proposal(proposalinfo_cancel.get('ProposalID')) == 0
        assert pip_obj.get_accu_verifiers_of_proposal(proposalinfo_cancel.get('ProposalID')) == len(client_verifier_obj_list)

        assert pip_obj.get_canceledby_of_proposal(proposalinfo_version.get('ProposalID')) == proposalinfo_cancel.get('ProposalID')
        assert pip_obj.get_status_of_proposal(proposalinfo_version.get('ProposalID')) == 6
        assert pip_obj.get_yeas_of_proposal(proposalinfo_version.get('ProposalID')) == 0
        assert pip_obj.get_nays_of_proposal(proposalinfo_version.get('ProposalID')) == 0
        assert pip_obj.get_abstentions_of_proposal(proposalinfo_version.get('ProposalID')) == 0
        assert pip_obj.get_accu_verifiers_of_proposal(proposalinfo_version.get('ProposalID')) == len(client_verifier_obj_list)

    @pytest.mark.P0
    @pytest.mark.compatibility
    def test_TR_IN_001(self, no_vp_proposal, client_verifier_obj_list):
        pip_obj = no_vp_proposal
        submitcvpandvote(client_verifier_obj_list, 1, 2, 3, 3)
        proposalinfo_cancel = pip_obj.get_effect_proposal_info_of_vote(pip_obj.cfg.cancel_proposal)
        log.info('Cancel proposal information {}'.format(proposalinfo_cancel))
        proposalinfo_version = pip_obj.get_effect_proposal_info_of_vote(pip_obj.cfg.version_proposal)
        log.info('Version proposal information {}'.format(proposalinfo_version))
        wait_block_number(pip_obj.node, proposalinfo_version.get('EndVotingBlock'))
        assert pip_obj.get_canceledby_of_proposal(proposalinfo_cancel.get('ProposalID')) == \
            "0x0000000000000000000000000000000000000000000000000000000000000000"
        assert pip_obj.get_status_of_proposal(proposalinfo_cancel.get('ProposalID')) == 3
        assert pip_obj.get_yeas_of_proposal(proposalinfo_cancel.get('ProposalID')) == 1
        assert pip_obj.get_nays_of_proposal(proposalinfo_cancel.get('ProposalID')) == 1
        assert pip_obj.get_abstentions_of_proposal(proposalinfo_cancel.get('ProposalID')) == 2
        assert pip_obj.get_accu_verifiers_of_proposal(proposalinfo_cancel.get('ProposalID')) == len(client_verifier_obj_list)

        assert pip_obj.get_canceledby_of_proposal(proposalinfo_version.get('ProposalID')) == \
            "0x0000000000000000000000000000000000000000000000000000000000000000"
        assert pip_obj.get_status_of_proposal(proposalinfo_version.get('ProposalID')) == 3
        assert pip_obj.get_yeas_of_proposal(proposalinfo_version.get('ProposalID')) == 0
        assert pip_obj.get_nays_of_proposal(proposalinfo_version.get('ProposalID')) == 0
        assert pip_obj.get_abstentions_of_proposal(proposalinfo_version.get('ProposalID')) == 0
        assert pip_obj.get_accu_verifiers_of_proposal(proposalinfo_version.get('ProposalID')) == len(client_verifier_obj_list)

    @pytest.mark.P0
    def test_TR_IN_010_005(self, new_genesis_env, client_con_list_obj):
        genesis = from_dict(data_class=Genesis, data=new_genesis_env.genesis_config)
        genesis.economicModel.gov.paramProposalVoteDurationSeconds = 0
        new_genesis_env.set_genesis(genesis.to_dict())
        new_genesis_env.deploy_all()
        pip_obj = client_con_list_obj[0].pip
        submitppandvote(client_con_list_obj[0:-1], 1, 2, 3)
        proposalinfo = client_con_list_obj[0].pip.get_effect_proposal_info_of_vote(client_con_list_obj[0].pip.cfg.param_proposal)
        log.info('Param proposal information {}'.format(proposalinfo))
        result = pip_obj.pip.getTallyResult(proposalinfo.get('ProposalID'))
        log.info('Interface getTallyResult info : {}'.format(result))
        assert_code(result, 302030)
        wait_block_number(pip_obj.node, proposalinfo.get('EndVotingBlock'))

        assert pip_obj.get_canceledby_of_proposal(proposalinfo.get('ProposalID')) == \
            "0x0000000000000000000000000000000000000000000000000000000000000000"
        assert pip_obj.get_status_of_proposal(proposalinfo.get('ProposalID')) == 3
        assert pip_obj.get_yeas_of_proposal(proposalinfo.get('ProposalID')) == 1
        assert pip_obj.get_nays_of_proposal(proposalinfo.get('ProposalID')) == 1
        assert pip_obj.get_abstentions_of_proposal(proposalinfo.get('ProposalID')) == 1
        assert pip_obj.get_accu_verifiers_of_proposal(proposalinfo.get('ProposalID')) == len(client_con_list_obj)

    @pytest.mark.P0
    def test_TR_IN_011_TR_IN_012(self, new_genesis_env, client_con_list_obj):
        genesis = from_dict(data_class=Genesis, data=new_genesis_env.genesis_config)
        genesis.economicModel.gov.paramProposalVoteDurationSeconds = 0
        new_genesis_env.set_genesis(genesis.to_dict())
        new_genesis_env.deploy_all()
        pip_obj = client_con_list_obj[0].pip
        submitcppandvote(client_con_list_obj, [1, 1, 1, 3])
        proposalinfo_param = pip_obj.get_effect_proposal_info_of_vote(pip_obj.cfg.param_proposal)
        log.info('Param proposal information {}'.format(proposalinfo_param))
        proposalinfo_cancel = pip_obj.get_effect_proposal_info_of_vote(pip_obj.cfg.cancel_proposal)
        log.info('Cancel proposal information {}'.format(proposalinfo_cancel))

        wait_block_number(pip_obj.node, proposalinfo_cancel.get('EndVotingBlock'))

        assert pip_obj.get_canceledby_of_proposal(proposalinfo_cancel.get('ProposalID')) == \
            "0x0000000000000000000000000000000000000000000000000000000000000000"
        assert pip_obj.get_canceledby_of_proposal(proposalinfo_param.get('ProposalID')) == proposalinfo_cancel.get('ProposalID')

        assert pip_obj.get_status_of_proposal(proposalinfo_cancel.get('ProposalID')) == 2
        assert pip_obj.get_yeas_of_proposal(proposalinfo_cancel.get('ProposalID')) == 3
        assert pip_obj.get_nays_of_proposal(proposalinfo_cancel.get('ProposalID')) == 0
        assert pip_obj.get_abstentions_of_proposal(proposalinfo_cancel.get('ProposalID')) == 1
        assert pip_obj.get_accu_verifiers_of_proposal(proposalinfo_cancel.get('ProposalID')) == len(client_con_list_obj)

        assert pip_obj.get_status_of_proposal(proposalinfo_param.get('ProposalID')) == 6
        assert pip_obj.get_yeas_of_proposal(proposalinfo_param.get('ProposalID')) == 0
        assert pip_obj.get_nays_of_proposal(proposalinfo_param.get('ProposalID')) == 0
        assert pip_obj.get_abstentions_of_proposal(proposalinfo_param.get('ProposalID')) == 0
        assert pip_obj.get_accu_verifiers_of_proposal(proposalinfo_param.get('ProposalID')) == len(client_con_list_obj)

    @pytest.mark.P1
    def test_TR_IN_006(self, client_verifier_obj):
        pip_obj = client_verifier_obj.pip
        result = pip_obj.pip.getTallyResult('0x9992d1f843fe8f376884d871f87605dda02da0722fd6b350bbf683518f73f111')
        log.info('Ineffective proposalID, interface getTallyResult return : {}'.format(result))
        assert_code(result, 302030)


class TestgetAccuVerifiersCount():
    @pytest.mark.P0
    def test_AC_IN_018_to_025(self, new_genesis_env, client_con_list_obj):
        genesis = from_dict(data_class=Genesis, data=new_genesis_env.genesis_config)
        genesis.economicModel.gov.paramProposalVoteDurationSeconds = 0
        new_genesis_env.set_genesis(genesis.to_dict())
        new_genesis_env.deploy_all()
        pip_obj = client_con_list_obj[0].pip
        pip_obj_test = client_con_list_obj[-1].pip
        result = pip_obj.submitParam(pip_obj.node.node_id, str(time.time()), 'slashing', 'slashBlocksReward', '999',
                                     pip_obj.node.staking_address, transaction_cfg=pip_obj.cfg.transaction_cfg)
        log.info('Node submit param proposal result : {}'.format(result))
        assert_code(result, 0)
        proposalinfo_param = pip_obj.get_effect_proposal_info_of_vote(pip_obj.cfg.param_proposal)
        log.info('Get param proposal information : {}'.format(proposalinfo_param))
        result = pip_obj.submitCancel(pip_obj.node.node_id, str(time.time()), 3, proposalinfo_param.get('ProposalID'),
                                      pip_obj.node.staking_address, transaction_cfg=pip_obj.cfg.transaction_cfg)
        log.info('Node submit cancel proposal result : {}'.format(result))
        proposalinfo_cancel = pip_obj.get_effect_proposal_info_of_vote(pip_obj.cfg.cancel_proposal)
        log.info('Get cancel proposal information : {}'.format(proposalinfo_cancel))
        for index in range(3):
            client_obj = client_con_list_obj[index]
            result = client_obj.pip.vote(client_obj.node.node_id, proposalinfo_param.get('ProposalID'), index + 1,
                                         client_obj.node.staking_address, transaction_cfg=pip_obj.cfg.transaction_cfg)
            log.info('Node {} vote param proposal result : {}'.format(client_obj.node.node_id, result))
            assert_code(result, 0)
            result = client_obj.pip.vote(client_obj.node.node_id, proposalinfo_cancel.get('ProposalID'), index + 1,
                                         client_obj.node.staking_address, transaction_cfg=pip_obj.cfg.transaction_cfg)
            log.info('Node {} vote cancel proposal result : {}'.format(client_obj.node.node_id, result))
            assert_code(result, 0)
        assert pip_obj.get_accuverifiers_count(proposalinfo_param.get('ProposalID')) == [4, 1, 1, 1]
        assert pip_obj.get_accuverifiers_count(proposalinfo_cancel.get('ProposalID')) == [4, 1, 1, 1]
        log.info('Stop the node {}'.format(pip_obj.node.node_id))
        pip_obj.node.stop()
        pip_obj_test.economic.wait_consensus_blocknum(pip_obj_test.node, 2)
        assert pip_obj_test.get_accuverifiers_count(proposalinfo_param.get('ProposalID')) == [4, 0, 1, 1]
        assert pip_obj_test.get_accuverifiers_count(proposalinfo_cancel.get('ProposalID')) == [4, 0, 1, 1]

        report_information = mock_duplicate_sign(1, client_con_list_obj[1].node.nodekey,
                                                 client_con_list_obj[1].node.blsprikey, 41)
        log.info("Report information: {}".format(report_information))
        address, _ = pip_obj_test.economic.account.generate_account(pip_obj_test.node.web3, 10 ** 18 * 1000)
        result = client_con_list_obj[-1].duplicatesign.reportDuplicateSign(1, report_information, address)
        log.info('Node duplicate block result : {}'.format(result))
        assert_code(result, 0)
        time.sleep(2)
        assert pip_obj_test.get_accuverifiers_count(proposalinfo_param.get('ProposalID')) == [4, 0, 0, 1]
        assert pip_obj_test.get_accuverifiers_count(proposalinfo_cancel.get('ProposalID')) == [4, 0, 0, 1]

        report_information = mock_duplicate_sign(2, client_con_list_obj[2].node.nodekey,
                                                 client_con_list_obj[2].node.blsprikey, 41)
        log.info("Report information: {}".format(report_information))
        address, _ = pip_obj_test.economic.account.generate_account(pip_obj_test.node.web3, 10 ** 18 * 1000)
        result = client_con_list_obj[-1].duplicatesign.reportDuplicateSign(2, report_information, address)
        log.info('Node duplicate block result : {}'.format(result))
        assert_code(result, 0)
        assert pip_obj_test.get_accuverifiers_count(proposalinfo_param.get('ProposalID')) == [4, 0, 0, 0]
        assert pip_obj_test.get_accuverifiers_count(proposalinfo_cancel.get('ProposalID')) == [4, 0, 0, 0]

    @pytest.mark.compatibility
    @pytest.mark.P0
    def test_AC_IN_001_002_004_to_006_012_to_014(self, no_vp_proposal, client_verifier_obj_list):
        pip_obj = client_verifier_obj_list[-1].pip
        result = pip_obj.submitVersion(pip_obj.node.node_id, str(time.time()), pip_obj.cfg.version5, 5, pip_obj.node.staking_address,
                                       transaction_cfg=pip_obj.cfg.transaction_cfg)
        log.info('Submit version proposal result : {}'.format(result))
        assert_code(result, 0)
        proposalinfo_version = pip_obj.get_effect_proposal_info_of_vote()
        log.info('Get version proposal information : {}'.format(proposalinfo_version))

        result = pip_obj.submitCancel(pip_obj.node.node_id, str(time.time()), 4, proposalinfo_version.get('ProposalID'),
                                      pip_obj.node.staking_address, transaction_cfg=pip_obj.cfg.transaction_cfg)
        log.info('Submit cancel proposal result : {}'.format(result))
        assert_code(result, 0)
        proposalinfo_cancel = pip_obj.get_effect_proposal_info_of_vote(pip_obj.cfg.cancel_proposal)
        log.info('Get cancel proposal information : {}'.format(proposalinfo_version))

        for index in range(3):
            client_obj = client_verifier_obj_list[index]
            result = version_proposal_vote(client_obj.pip)
            assert_code(result, 0)
            result = client_obj.pip.vote(client_obj.node.node_id, proposalinfo_cancel.get('ProposalID'), index + 1,
                                         client_obj.node.staking_address, transaction_cfg=pip_obj.cfg.transaction_cfg)
            log.info('Node {} vote cancel proposal result : {}'.format(client_obj.node.node_id, result))
            assert_code(result, 0)

        assert pip_obj.get_accuverifiers_count(proposalinfo_version.get('ProposalID')) == [4, 3, 0, 0]
        assert pip_obj.get_accuverifiers_count(proposalinfo_cancel.get('ProposalID')) == [4, 1, 1, 1]
        log.info('Stop the node {}'.format(client_verifier_obj_list[0].node.node_id))
        client_verifier_obj_list[0].node.stop()
        pip_obj.economic.wait_consensus_blocknum(pip_obj.node, 2)
        assert pip_obj.get_accuverifiers_count(proposalinfo_version.get('ProposalID')) == [4, 2, 0, 0]
        assert pip_obj.get_accuverifiers_count(proposalinfo_cancel.get('ProposalID')) == [4, 0, 1, 1]

        report_information = mock_duplicate_sign(1, client_verifier_obj_list[1].node.nodekey,
                                                 client_verifier_obj_list[1].node.blsprikey, 41)
        log.info("Report information: {}".format(report_information))
        address, _ = pip_obj.economic.account.generate_account(pip_obj.node.web3, 10 ** 18 * 1000)
        result = client_verifier_obj_list[-1].duplicatesign.reportDuplicateSign(1, report_information, address)
        log.info('Node duplicate block result : {}'.format(result))
        assert_code(result, 0)
        assert pip_obj.get_accuverifiers_count(proposalinfo_version.get('ProposalID')) == [4, 1, 0, 0]
        assert pip_obj.get_accuverifiers_count(proposalinfo_cancel.get('ProposalID')) == [4, 0, 0, 1]

        report_information = mock_duplicate_sign(2, client_verifier_obj_list[2].node.nodekey,
                                                 client_verifier_obj_list[2].node.blsprikey, 41)
        log.info("Report information: {}".format(report_information))
        address, _ = pip_obj.economic.account.generate_account(pip_obj.node.web3, 10 ** 18 * 1000)
        result = client_verifier_obj_list[-1].duplicatesign.reportDuplicateSign(2, report_information, address)
        log.info('Node duplicate block result : {}'.format(result))
        assert_code(result, 0)
        assert pip_obj.get_accuverifiers_count(proposalinfo_version.get('ProposalID')) == [4, 0, 0, 0]
        assert pip_obj.get_accuverifiers_count(proposalinfo_cancel.get('ProposalID')) == [4, 0, 0, 0]

    @pytest.mark.P0
    def test_AC_IN_003_008_010(self, new_genesis_env, client_con_list_obj):
        genesis = from_dict(data_class=Genesis, data=new_genesis_env.genesis_config)
        genesis.economicModel.gov.textProposalVoteDurationSeconds = 120
        new_genesis_env.set_genesis(genesis.to_dict())
        new_genesis_env.deploy_all()
        pip_obj = client_con_list_obj[-1].pip
        submittpandvote(client_con_list_obj, 1, 2, 3, 1)
        proposalinfo = pip_obj.get_effect_proposal_info_of_vote(pip_obj.cfg.text_proposal)
        log.info('Get text proposal information : {}'.format(proposalinfo))
        assert pip_obj.get_accuverifiers_count(proposalinfo.get('ProposalID')) == [4, 2, 1, 1]
        log.info('Stop the node {}'.format(client_con_list_obj[0].node.node_id))
        client_con_list_obj[0].node.stop()
        pip_obj.economic.wait_consensus_blocknum(pip_obj.node, 2)
        assert pip_obj.get_accuverifiers_count(proposalinfo.get('ProposalID')) == [4, 1, 1, 1]

        report_information = mock_duplicate_sign(1, client_con_list_obj[1].node.nodekey,
                                                 client_con_list_obj[1].node.blsprikey, 41)
        log.info("Report information: {}".format(report_information))
        address, _ = pip_obj.economic.account.generate_account(pip_obj.node.web3, 10 ** 18 * 1000)
        result = client_con_list_obj[-1].duplicatesign.reportDuplicateSign(1, report_information, address)
        log.info('Node duplicate block result : {}'.format(result))
        assert_code(result, 0)
        assert pip_obj.get_accuverifiers_count(proposalinfo.get('ProposalID')) == [4, 1, 0, 1]

        report_information = mock_duplicate_sign(2, client_con_list_obj[2].node.nodekey,
                                                 client_con_list_obj[2].node.blsprikey, 41)
        log.info("Report information: {}".format(report_information))
        address, _ = pip_obj.economic.account.generate_account(pip_obj.node.web3, 10 ** 18 * 1000)
        result = client_con_list_obj[-1].duplicatesign.reportDuplicateSign(2, report_information, address)
        log.info('Node duplicate block result : {}'.format(result))
        assert_code(result, 0)
        assert pip_obj.get_accuverifiers_count(proposalinfo.get('ProposalID')) == [4, 1, 0, 0]

    @pytest.mark.P2
    def test_AC_IN_016_to_018(self, client_verifier_obj):
        pip_obj = client_verifier_obj.pip
        result = pip_obj.submitText(pip_obj.node.node_id, str(time.time()), pip_obj.node.staking_address,
                                    transaction_cfg=pip_obj.cfg.transaction_cfg)
        log.info('Submit text proposal result : {}'.format(result))
        assert_code(result, 0)
        proposalinfo = pip_obj.get_effect_proposal_info_of_vote(pip_obj.cfg.text_proposal)
        log.info('Get text proposal information : {}'.format(proposalinfo))
        result = pip_obj.pip.getAccuVerifiersCount('0x0c04f578466ead2208dbb15b927ecb27041881e8c16c17cd0db6b3df422e1111',
                                                   block_hash=get_blockhash(pip_obj.node))
        log.info('Interface getAccuVerifiersCount result : {}'.format(result))
        assert_code(result, 302006)
        log.info('{}'.format(get_blockhash(pip_obj.node)))

        result = pip_obj.pip.getAccuVerifiersCount(proposalinfo.get('ProposalID'), block_hash='')
        log.info('Interface getAccuVerifiersCount result : {}'.format(result))
        assert not result

        result = pip_obj.pip.getAccuVerifiersCount(proposalinfo.get('ProposalID'),
                                                   block_hash='0x5941605fe43ab32fbaf9c6e08dc0970eae50efb7da4248a9a8941f0e50711111')
        log.info('Interface getAccuVerifiersCount result : {}'.format(result))
        assert_code(result, 0)


class TestListGovernParam():
    def get_govern_param(self, client_obj, module=None):
        result = client_obj.pip.pip.listGovernParam(module)
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
    def test_IN_LG_001(self, client_noconsensus_obj):
        name, module = self.get_govern_param(client_noconsensus_obj)
        assert set(name) == {'maxValidators', 'unStakeFreezeDuration', 'operatingThreshold', 'slashBlocksReward',
                             'stakeThreshold', 'maxBlockGasLimit', 'duplicateSignReportReward', 'maxEvidenceAge', 'slashFractionDuplicateSign'}
        assert set(module) == {'block', 'slashing', 'staking'}

    @pytest.mark.P2
    def test_IN_LG_002(self, client_noconsensus_obj):
        name, module = self.get_govern_param(client_noconsensus_obj, 'staking')
        assert set(name) == {'maxValidators', 'unStakeFreezeDuration', 'operatingThreshold', 'stakeThreshold'}
        assert set(module) == {'staking'}

    @pytest.mark.P2
    def test_IN_LG_003(self, client_noconsensus_obj):
        name, module = self.get_govern_param(client_noconsensus_obj, 'slashing')
        assert set(name) == {'slashBlocksReward', 'duplicateSignReportReward', 'maxEvidenceAge', 'slashFractionDuplicateSign'}
        assert set(module) == {'slashing'}

    @pytest.mark.P2
    def test_IN_LG_004(self, client_noconsensus_obj):
        name, module = self.get_govern_param(client_noconsensus_obj, 'block')
        assert set(name) == {'maxBlockGasLimit'}
        assert set(module) == {'block'}

    @pytest.mark.P2
    def test_IN_LG_005(self, client_noconsensus_obj):
        result = client_noconsensus_obj.pip.pip.listGovernParam('txpool')
        log.info('Interface listGovernParam result {}'.format(result))


class TestGetGovernParam():
    @pytest.mark.P0
    def test_IN_GG_001(self, client_noconsensus_obj):
        client_noconsensus_obj.economic.env.deploy_all()
        genesis = from_dict(data_class=Genesis, data=client_noconsensus_obj.economic.env.genesis_config)
        pip_obj = client_noconsensus_obj.pip.pip
        result = pip_obj.getGovernParamValue('slashing', 'slashBlocksReward')
        log.info('Interface getGovernParamValue result : {}'.format(result))
        assert genesis.economicModel.slashing.slashBlocksReward == int(result.get('Ret'))

        result = pip_obj.getGovernParamValue('slashing', 'maxEvidenceAge')
        log.info('Interface getGovernParamValue result : {}'.format(result))
        assert genesis.economicModel.slashing.maxEvidenceAge == int(result.get('Ret'))

        result = pip_obj.getGovernParamValue('slashing', 'slashFractionDuplicateSign')
        log.info('Interface getGovernParamValue result : {}'.format(result))
        assert genesis.economicModel.slashing.slashFractionDuplicateSign == int(result.get('Ret'))

        result = pip_obj.getGovernParamValue('slashing', 'duplicateSignReportReward')
        log.info('Interface getGovernParamValue result : {}'.format(result))
        assert genesis.economicModel.slashing.duplicateSignReportReward == int(result.get('Ret'))

        result = pip_obj.getGovernParamValue('staking', 'stakeThreshold')
        log.info('Interface getGovernParamValue result : {}'.format(result))
        assert genesis.economicModel.staking.stakeThreshold == int(result.get('Ret'))

        result = pip_obj.getGovernParamValue('staking', 'operatingThreshold')
        log.info('Interface getGovernParamValue result : {}'.format(result))
        assert genesis.economicModel.staking.operatingThreshold == int(result.get('Ret'))

        result = pip_obj.getGovernParamValue('staking', 'unStakeFreezeDuration')
        log.info('Interface getGovernParamValue result : {}'.format(result))
        assert genesis.economicModel.staking.unStakeFreezeDuration == int(result.get('Ret'))

        result = pip_obj.getGovernParamValue('staking', 'maxValidators')
        log.info('Interface getGovernParamValue result : {}'.format(result))
        assert genesis.economicModel.staking.maxValidators == int(result.get('Ret'))

        result = pip_obj.getGovernParamValue('block', 'maxBlockGasLimit')
        log.info('Interface getGovernParamValue result : {}'.format(result))

    @pytest.mark.P2
    def test_IN_GG_002(self, client_noconsensus_obj):
        pip_obj = client_noconsensus_obj.pip.pip
        result = pip_obj.getGovernParamValue('Staking', 'maxValidators')
        assert_code(result, 302031)
        pip_obj = client_noconsensus_obj.pip.pip
        result = pip_obj.getGovernParamValue('Slashing', 'slashBlocksReward')
        assert_code(result, 302031)
        pip_obj = client_noconsensus_obj.pip.pip
        result = pip_obj.getGovernParamValue('Block', 'maxBlockGasLimit')
        assert_code(result, 302031)

    @pytest.mark.P2
    def test_IN_GG_003(self, client_noconsensus_obj):
        pip_obj = client_noconsensus_obj.pip.pip
        result = pip_obj.getGovernParamValue('staking', 'MaxValidators')
        assert_code(result, 302031)
        pip_obj = client_noconsensus_obj.pip.pip
        result = pip_obj.getGovernParamValue('slashing', 'SlashBlocksReward')
        assert_code(result, 302031)
        pip_obj = client_noconsensus_obj.pip.pip
        result = pip_obj.getGovernParamValue('block', 'MaxValidators')
        assert_code(result, 302031)


class TestGetActiveVersion():
    @pytest.mark.compatibility
    @pytest.mark.P0
    def test_AV_IN_001(self, no_vp_proposal):
        assert_code(no_vp_proposal.chain_version, no_vp_proposal.cfg.version0)

    @pytest.mark.P0
    def test_AV_IN_002_003(self, client_verifier_obj_list):
        pip_obj = client_verifier_obj_list[0].pip
        submitvpandvote(client_verifier_obj_list)
        proposalinfo = pip_obj.get_effect_proposal_info_of_vote()
        log.info('Version proposal information : {}'.format(proposalinfo))
        wait_block_number(pip_obj.node, proposalinfo.get('EndVotingBlock'))
        assert_code(pip_obj.get_status_of_proposal(proposalinfo.get('ProposalID')), 4)
        assert_code(pip_obj.chain_version, pip_obj.cfg.version0)
        wait_block_number(pip_obj.node, proposalinfo.get('ActiveBlock'))
        assert_code(pip_obj.get_status_of_proposal(proposalinfo.get('ProposalID')), 5)
        assert_code(pip_obj.chain_version, pip_obj.cfg.version5)


class TestListProposal():
    @pytest.mark.P1
    def test_LP_IN_001_002(self, no_vp_proposal):
        pip_obj = no_vp_proposal
        pip_id = str(time.time())
        result = pip_obj.submitParam(pip_obj.node.node_id, pip_id, 'slashing', 'slashBlocksReward', '456',
                                     pip_obj.node.staking_address, transaction_cfg=pip_obj.cfg.transaction_cfg)
        log.info('Submit param proposal result : {}'.format(result))
        assert_code(result, 0)
        proposalinfo = pip_obj.get_effect_proposal_info_of_vote((pip_obj.cfg.param_proposal))
        log.info('Get param proposal information : {}'.format(proposalinfo))
        assert proposalinfo.get('Proposer') == pip_obj.node.node_id
        assert proposalinfo.get('ProposalType') == pip_obj.cfg.param_proposal
        log.info('{}'.format(pip_obj.economic.pp_vote_settlement_wheel))
        calculated_endvotingblock = math.ceil(proposalinfo.get('SubmitBlock') / pip_obj.economic.settlement_size +
                                              pip_obj.economic.pp_vote_settlement_wheel) * pip_obj.economic.settlement_size
        assert proposalinfo.get('EndVotingBlock') == calculated_endvotingblock

        pip_id = str(time.time())
        result = pip_obj.submitCancel(pip_obj.node.node_id, pip_id, 1, proposalinfo.get('ProposalID'),
                                      pip_obj.node.staking_address, transaction_cfg=pip_obj.cfg.transaction_cfg)
        log.info('Submit cancel proposal result : {}'.format(result))
        assert_code(result, 0)
        proposalinfo_cancel = pip_obj.get_effect_proposal_info_of_vote(pip_obj.cfg.cancel_proposal)
        log.info('Get cancel proposal information : {}'.format(proposalinfo_cancel))

        assert proposalinfo_cancel.get('Proposer') == pip_obj.node.node_id
        assert proposalinfo_cancel.get('ProposalType') == pip_obj.cfg.cancel_proposal
        log.info('{}'.format(pip_obj.economic.pp_vote_settlement_wheel))
        calculated_endvotingblock = math.ceil(proposalinfo_cancel.get('SubmitBlock') / pip_obj.economic.consensus_size +
                                              1) * pip_obj.economic.consensus_size - 20
        assert proposalinfo_cancel.get('EndVotingBlock') == calculated_endvotingblock

    @pytest.mark.P1
    def test_LP_IN_003(self, client_consensus_obj):
        client_consensus_obj.economic.env.deploy_all()
        result = client_consensus_obj.pip.pip.listProposal()
        log.info('There is no proposal, interface listProposal return : {}'.format(result))
        assert_code(result, 2)
        assert result.get('Ret') == "Object not found"


class TestGasUse():
    def get_balance(self, pip_obj):
        balance = pip_obj.node.eth.getBalance(pip_obj.node.staking_address)
        log.info('address balance : {}'.format(balance))
        return balance

    def test_submitText(self, client_verifier_obj):
        pip_obj = client_verifier_obj.pip
        pip_id = str(time.time())
        data = rlp.encode([rlp.encode(int(2000)), rlp.encode(bytes.fromhex(pip_obj.node.node_id)), rlp.encode(pip_id)])
        balance_before = self.get_balance(pip_obj)
        result = pip_obj.submitText(pip_obj.node.node_id, pip_id, pip_obj.node.staking_address,
                                    transaction_cfg=pip_obj.cfg.transaction_cfg)
        log.info('Submit text proposal result : {}'.format(result))
        assert_code(result, 0)
        proposalinfo = pip_obj.get_effect_proposal_info_of_vote(pip_obj.cfg.text_proposal)
        log.info('Get text proposal information : {}'.format(proposalinfo))
        balance_after = self.get_balance(pip_obj)
        gas = get_the_dynamic_parameter_gas_fee(data)
        log.info('Calculated gas : {}'.format(gas))
        assert_code(balance_before - balance_after, (gas + 350000) * pip_obj.cfg.transaction_cfg.get('gasPrice'))

        proposal_id = proposalinfo.get('ProposalID')[2:]
        version_sign = pip_obj.node.program_version_sign[2:]
        data = rlp.encode([rlp.encode(int(2003)), rlp.encode(bytes.fromhex(pip_obj.node.node_id)),
                           rlp.encode(bytes.fromhex(proposal_id)),
                           rlp.encode(pip_obj.cfg.vote_option_yeas), rlp.encode(int(pip_obj.node.program_version)),
                           rlp.encode(bytes.fromhex(version_sign))])
        result = pip_obj.vote(pip_obj.node.node_id, proposalinfo.get('ProposalID'), pip_obj.cfg.vote_option_yeas,
                              pip_obj.node.staking_address, transaction_cfg=pip_obj.cfg.transaction_cfg)
        log.info('Vote reuslt : {}'.format(result))
        assert_code(result, 0)
        balance_after_vote = pip_obj.node.eth.getBalance(pip_obj.node.staking_address)
        log.info('After vote text proposal, the address balance : {}'.format(balance_after_vote))
        gas = get_the_dynamic_parameter_gas_fee(data)
        log.info('Calculated gas : {}'.format(gas))
        assert_code(balance_after - balance_after_vote, (gas + 32000) * pip_obj.cfg.transaction_cfg.get('gasPrice'))

    def test_submitversion(self, no_vp_proposal):
        pip_obj = no_vp_proposal
        pip_id = str(time.time())
        balance_before = self.get_balance(pip_obj)
        result = pip_obj.submitVersion(pip_obj.node.node_id, pip_id, pip_obj.cfg.version5, 1, pip_obj.node.staking_address,
                                       transaction_cfg=pip_obj.cfg.transaction_cfg)
        log.info('Submit version proposal result : {}'.format(result))
        assert_code(result, 0)
        balance_after = self.get_balance(pip_obj)
        data = rlp.encode([rlp.encode(int(2001)), rlp.encode(bytes.fromhex(pip_obj.node.node_id)), rlp.encode(pip_id),
                           rlp.encode(int(pip_obj.cfg.version5)), rlp.encode(int(1))])
        gas = get_the_dynamic_parameter_gas_fee(data)
        log.info('Calculated gas : {}'.format(gas))
        assert_code(balance_before - balance_after, (gas + 480000) * pip_obj.cfg.transaction_cfg.get('gasPrice'))

    def test_submitparam_and_cancel(self, no_vp_proposal):
        pip_obj = no_vp_proposal
        pip_id = str(time.time())
        balance_before = self.get_balance(pip_obj)
        result = pip_obj.submitParam(pip_obj.node.node_id, pip_id, 'slashing', 'slashBlocksReward', '123',
                                     pip_obj.node.staking_address, transaction_cfg=pip_obj.cfg.transaction_cfg)
        log.info('Submit param proposal result : {}'.format(result))
        assert_code(result, 0)
        proposalinfor_param = pip_obj.get_effect_proposal_info_of_vote(pip_obj.cfg.param_proposal)
        log.info('Get param proposal information : {}'.format(proposalinfor_param))
        balance_after = self.get_balance(pip_obj)
        data = rlp.encode([rlp.encode(int(2002)), rlp.encode(bytes.fromhex(pip_obj.node.node_id)),
                           rlp.encode(pip_id), rlp.encode('slashing'), rlp.encode('slashBlocksReward'),
                           rlp.encode('123')])
        gas = get_the_dynamic_parameter_gas_fee(data)
        log.info('Calculated gas : {}'.format(gas))
        assert_code(balance_before - balance_after, (gas + 530000) * pip_obj.cfg.transaction_cfg.get('gasPrice'))

        pip_id = str(time.time())
        result = pip_obj.submitCancel(pip_obj.node.node_id, pip_id, 1, proposalinfor_param.get('ProposalID'),
                                      pip_obj.node.staking_address, transaction_cfg=pip_obj.cfg.transaction_cfg)
        log.info('Submit cancel proposal result : {}'.format(result))
        assert_code(result, 0)
        assert_code(balance_before - balance_after, (gas + 530000) * pip_obj.cfg.transaction_cfg.get('gasPrice'))
        balance_after_cancel = pip_obj.node.eth.getBalance(pip_obj.node.staking_address)
        log.info('After submitting cancel proposal, the address balance : {}'.format(balance_after_cancel))
        tobe_canceled_proposal_id = proposalinfor_param.get('ProposalID')[2:]
        data = rlp.encode([rlp.encode(int(2005)), rlp.encode(bytes.fromhex(pip_obj.node.node_id)), rlp.encode(pip_id),
                           rlp.encode(int(1)), rlp.encode(bytes.fromhex(tobe_canceled_proposal_id))])
        gas = get_the_dynamic_parameter_gas_fee(data)
        log.info('Calculated gas : {}'.format(gas))
        assert_code(balance_after - balance_after_cancel, (gas + 530000) * pip_obj.cfg.transaction_cfg.get('gasPrice'))

    def test_declareversion(self, client_verifier_obj):
        pip_obj = client_verifier_obj.pip
        balance_before = self.get_balance(pip_obj)
        result = pip_obj.declareVersion(pip_obj.node.node_id, pip_obj.node.staking_address,
                                        transaction_cfg=pip_obj.cfg.transaction_cfg)
        log.info('Declare version result : {}'.format(result))
        assert_code(result, 0)
        version_sign = pip_obj.node.program_version_sign[2:]
        data = rlp.encode([rlp.encode(int(2004)), rlp.encode(bytes.fromhex(pip_obj.node.node_id)),
                           rlp.encode(int(pip_obj.node.program_version)), rlp.encode(bytes.fromhex(version_sign))])
        gas = get_the_dynamic_parameter_gas_fee(data)
        log.info('Calculated gas : {}'.format(gas))
        balance_after = self.get_balance(pip_obj)
        assert_code(balance_before - balance_after, (gas + 33000) * pip_obj.cfg.transaction_cfg.get('gasPrice'))


if __name__ == '__main__':
    pytest.main(['./tests/govern/', '-s', '-q', '--alluredir', './report/report'])
