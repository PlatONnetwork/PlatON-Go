from common.log import log
from tests.lib.utils import assert_code, wait_block_number, get_pledge_list
from dacite import from_dict
from tests.lib import Genesis
from common.key import mock_duplicate_sign
from tests.govern.test_voting_statistics import submitppandvote
import json, time

class TestgetProposal():
    def test_GP_IF_001(self, submit_cancel_param):
        pip_obj = submit_cancel_param
        proposalinfo = pip_obj.get_effect_proposal_info_of_vote(pip_obj.cfg.cancel_proposal)
        log.info('Cancel proposal information : {}'.format(proposalinfo))
        result = pip_obj.pip.getProposal(proposalinfo.get('ProposalID'))
        log.info('Interface getProposal result : {}'.format(result))
        assert json.loads(result.get('Ret')).get('Proposer') == pip_obj.node.node_id
        assert json.loads(result.get('Ret')).get('ProposalType') == pip_obj.cfg.cancel_proposal
        assert json.loads(result.get('Ret')).get('PIPID') == proposalinfo.get('PIPID')
        assert json.loads(result.get('Ret')).get('SubmitBlock') == proposalinfo.get('SubmitBlock')
        assert json.loads(result.get('Ret')).get('EndVotingBlock') == proposalinfo.get('EndVotingBlock')

    def test_GP_IF_002(self, submit_param):
        pip_obj = submit_param
        proposalinfo = pip_obj.get_effect_proposal_info_of_vote(pip_obj.cfg.param_proposal)
        log.info('Param proposal information : {}'.format(proposalinfo))
        result = pip_obj.pip.getProposal(proposalinfo.get('ProposalID'))
        log.info('Interface getProposal result : {}'.format(result))
        assert json.loads(result.get('Ret')).get('Proposer') == pip_obj.node.node_id
        assert json.loads(result.get('Ret')).get('ProposalType') == pip_obj.cfg.param_proposal
        assert json.loads(result.get('Ret')).get('PIPID') == proposalinfo.get('PIPID')
        assert json.loads(result.get('Ret')).get('SubmitBlock') == proposalinfo.get('SubmitBlock')
        assert json.loads(result.get('Ret')).get('EndVotingBlock') == proposalinfo.get('EndVotingBlock')

class TestgetTallyResult():
    def test_TR_IN_010(self, new_genesis_env, client_con_list_obj):
        genesis = from_dict(data_class=Genesis, data=new_genesis_env.genesis_config)
        genesis.EconomicModel.Gov.ParamProposalVote_DurationSeconds = 0
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
        result = pip_obj.pip.getTallyResult(proposalinfo.get('ProposalID'))
        log.info('Interface getTallyResult info : {}'.format(result))
        assert_code(result, 0)
        assert json.loads(result.get('Ret')).get('canceledBy'
                                                 ) == "0x0000000000000000000000000000000000000000000000000000000000000000"
        assert pip_obj.get_status_of_proposal(proposalinfo.get('ProposalID')) == 3
        assert pip_obj.get_yeas_of_proposal(proposalinfo.get('ProposalID')) == 1
        assert pip_obj.get_nays_of_proposal(proposalinfo.get('ProposalID')) == 1
        assert pip_obj.get_abstentions_of_proposal(proposalinfo.get('ProposalID')) == 1
        assert pip_obj.get_accu_verifiers_of_proposal(proposalinfo.get('ProposalID')) == len(client_con_list_obj)


    def test_TR_IN_011_TR_IN_012(self, new_genesis_env, client_con_list_obj):
        genesis = from_dict(data_class=Genesis, data=new_genesis_env.genesis_config)
        genesis.EconomicModel.Gov.ParamProposalVote_DurationSeconds = 0
        new_genesis_env.set_genesis(genesis.to_dict())
        new_genesis_env.deploy_all()
        pip_obj = client_con_list_obj[0].pip
        result = pip_obj.submitParam(pip_obj.node.node_id, str(time.time()), 'Slashing', 'SlashBlocksReward',
                            '101', pip_obj.node.staking_address, transaction_cfg=pip_obj.cfg.transaction_cfg)
        log.info('Submit param proposal result {}'.format(result))
        assert_code(result, 0)
        proposalinfo_param = pip_obj.get_effect_proposal_info_of_vote(pip_obj.cfg.param_proposal)
        log.info('Param proposal information {}'.format(proposalinfo_param))
        result = pip_obj.submitCancel(pip_obj.node.node_id, str(time.time()), 2, proposalinfo_param.get('ProposalID'),
                             pip_obj.node.staking_address, transaction_cfg=pip_obj.cfg.transaction_cfg)
        log.info('Submit cancel proposal result {}'.format(result))
        assert_code(result, 0)
        proposalinfo_cancel = pip_obj.get_effect_proposal_info_of_vote(pip_obj.cfg.cancel_proposal)
        log.info('Cancel proposal information {}'.format(proposalinfo_cancel))

        for client_obj in client_con_list_obj:
            result = client_obj.pip.vote(client_obj.node.node_id, proposalinfo_cancel.get('proposalID'),
                                         pip_obj.cfg.vote_option_yeas, client_obj.node.staking_address,
                                         transaction_cfg=pip_obj.cfg.transaction_cfg)
            log.info('Node {} vote cancel proposal result : {}'.format(client_obj.node.node_id, result))
            assert_code(result, 0)
        wait_block_number(client_obj.node, proposalinfo_cancel.get('EndVotingBlock'))
        result_cancel = client_obj.pip.getTallyResult(proposalinfo_cancel.get('ProposalID'))
        result_param = client_obj.pip.getTallyResult(proposalinfo_param.get('ProposalID'))
        log.info('Interface getTallyResult of cancel proposal info : {}'.format(result_cancel))
        log.info('Interface getTallyResult of param proposal info : {}'.format(result_param))
        assert_code(result_cancel, 0)
        assert_code(result_param, 0)
        assert json.loads(result_cancel.get('Ret')).get('canceledBy'
                                                 ) == "0x0000000000000000000000000000000000000000000000000000000000000000"
        assert json.loads(result_param.get('Ret')).get('canceledBy') == proposalinfo_cancel.get('ProposalID')

        assert client_obj.pip.get_status_of_proposal(proposalinfo_cancel.get('ProposalID')) == 2
        assert client_obj.pip.get_yeas_of_proposal(proposalinfo_cancel.get('ProposalID')) == len(client_con_list_obj)
        assert client_obj.pip.get_nays_of_proposal(proposalinfo_cancel.get('ProposalID')) == 0
        assert client_obj.pip.get_abstentions_of_proposal(proposalinfo_cancel.get('ProposalID')) == 0
        assert client_obj.pip.get_accu_verifiers_of_proposal(proposalinfo_cancel.get('ProposalID')) == len(client_con_list_obj)

        assert client_obj.pip.get_status_of_proposal(proposalinfo_param.get('ProposalID')) == 0
        assert client_obj.pip.get_yeas_of_proposal(proposalinfo_param.get('ProposalID')) == len(client_con_list_obj)
        assert client_obj.pip.get_nays_of_proposal(proposalinfo_param.get('ProposalID')) == 0
        assert client_obj.pip.get_abstentions_of_proposal(proposalinfo_param.get('ProposalID')) == 0
        assert client_obj.pip.get_accu_verifiers_of_proposal(proposalinfo_param.get('ProposalID')) == len(client_con_list_obj)

class TestgetAccuVerifiersCount():
    def test_AC_IN_018_to_025(self, new_genesis_env, client_con_list_obj):
        genesis = from_dict(data_class=Genesis, data=new_genesis_env.genesis_config)
        genesis.EconomicModel.Gov.ParamProposalVote_DurationSeconds = 0
        new_genesis_env.set_genesis(genesis.to_dict())
        new_genesis_env.deploy_all()
        pip_obj = client_con_list_obj[0].pip
        pip_obj_test = client_con_list_obj[-1].pip
        result = pip_obj.submitParam(pip_obj.node.node_id, str(time.time()), 'Slashing', 'SlashBlocksReward', '999',
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
                                                 client_con_list_obj[1].node.blsprikey,
                                                 41)
        log.info("Report information: {}".format(report_information))
        address, _ = pip_obj_test.economic.account.generate_account(pip_obj_test.node.web3, 10 ** 18 * 1000)
        result = client_con_list_obj[-1].duplicatesign.reportDuplicateSign(1, report_information, address)
        log.info('Node duplicate block result : {}'.format(result))
        assert_code(result, 0)
        time.sleep(2)
        assert pip_obj_test.get_accuverifiers_count(proposalinfo_param.get('ProposalID')) == [4, 0, 0, 1]
        assert pip_obj_test.get_accuverifiers_count(proposalinfo_cancel.get('ProposalID')) == [4, 0, 0, 1]

        report_information = mock_duplicate_sign(2, client_con_list_obj[2].node.nodekey,
                                                 client_con_list_obj[2].node.blsprikey,
                                                 41)
        log.info("Report information: {}".format(report_information))
        address, _ = pip_obj_test.economic.account.generate_account(pip_obj_test.node.web3, 10 ** 18 * 1000)
        result = client_con_list_obj[-1].duplicatesign.reportDuplicateSign(2, report_information, address)
        log.info('Node duplicate block result : {}'.format(result))
        assert_code(result, 0)
        assert pip_obj_test.get_accuverifiers_count(proposalinfo_param.get('ProposalID')) == [4, 0, 0, 0]
        assert pip_obj_test.get_accuverifiers_count(proposalinfo_cancel.get('ProposalID')) == [4, 0, 0, 0]

class TestListGovernParam():
    def get_govern_param(self, client_obj, module=None):
        result = client_obj.pip.pip.listGovernParam(module)
        log.info('Interface listGovernParam result {}'.format(result))
        assert_code(result, 0)
        resultinfo = json.loads(result.get('Ret'))
        module = []
        name = []
        for param in resultinfo:
            module.append(param.get('ParamItem').get('Module'))
            name.append(param.get('ParamItem').get('Name'))
        return name, module

    def test_IN_LG_001(self, client_noconsensus_obj):
        name, module = self.get_govern_param(client_noconsensus_obj)
        assert set(name) == {'MaxValidators', 'UnStakeFreezeDuration', 'OperatingThreshold', 'SlashBlocksReward',
                             'StakeThreshold', 'MaxBlockGasLimit', 'DuplicateSignReportReward', 'MaxEvidenceAge', 'SlashFractionDuplicateSign'}
        assert set(module) == {'Block', 'Slashing', 'Staking'}

    def test_IN_LG_002(self, client_noconsensus_obj):
        name, module = self.get_govern_param(client_noconsensus_obj, 'Staking')
        assert set(name) == {'MaxValidators', 'UnStakeFreezeDuration', 'OperatingThreshold', 'StakeThreshold'}
        assert set(module) == {'Staking'}

    def test_IN_LG_003(self, client_noconsensus_obj):
        name, module = self.get_govern_param(client_noconsensus_obj, 'Slashing')
        assert set(name) == {'SlashBlocksReward', 'DuplicateSignReportReward', 'MaxEvidenceAge', 'SlashFractionDuplicateSign'}
        assert set(module) == {'Slashing'}

    def test_IN_LG_004(self, client_noconsensus_obj):
        name, module = self.get_govern_param(client_noconsensus_obj, 'Block')
        assert set(name) == {'MaxBlockGasLimit'}
        assert set(module) == {'Block'}

    def test_IN_LG_005(self, client_noconsensus_obj):
        result = client_noconsensus_obj.pip.pip.listGovernParam('Txpool')
        log.info('Interface listGovernParam result {}'.format(result))

class TestGetGovernParam():
    def test_IN_GG_001(self, client_noconsensus_obj):
        client_noconsensus_obj.economic.env.deploy_all()
        genesis = from_dict(data_class=Genesis, data=client_noconsensus_obj.economic.env.genesis_config)
        pip_obj = client_noconsensus_obj.pip.pip
        result = pip_obj.getGovernParamValue('Slashing', 'SlashBlocksReward')
        log.info('Interface getGovernParamValue result : {}'.format(result))
        assert genesis.EconomicModel.Slashing.SlashBlocksReward == int(result.get('Ret'))

        result = pip_obj.getGovernParamValue('Slashing', 'MaxEvidenceAge')
        log.info('Interface getGovernParamValue result : {}'.format(result))
        assert genesis.EconomicModel.Slashing.MaxEvidenceAge == int(result.get('Ret'))

        result = pip_obj.getGovernParamValue('Slashing', 'SlashFractionDuplicateSign')
        log.info('Interface getGovernParamValue result : {}'.format(result))
        assert genesis.EconomicModel.Slashing.SlashFractionDuplicateSign == int(result.get('Ret'))

        result = pip_obj.getGovernParamValue('Slashing', 'DuplicateSignReportReward')
        log.info('Interface getGovernParamValue result : {}'.format(result))
        assert genesis.EconomicModel.Slashing.DuplicateSignReportReward == int(result.get('Ret'))

        result = pip_obj.getGovernParamValue('Staking', 'StakeThreshold')
        log.info('Interface getGovernParamValue result : {}'.format(result))
        assert genesis.EconomicModel.Staking.StakeThreshold == int(result.get('Ret'))

        result = pip_obj.getGovernParamValue('Staking', 'OperatingThreshold')
        log.info('Interface getGovernParamValue result : {}'.format(result))
        assert genesis.EconomicModel.Staking.OperatingThreshold == int(result.get('Ret'))

        result = pip_obj.getGovernParamValue('Staking', 'UnStakeFreezeDuration')
        log.info('Interface getGovernParamValue result : {}'.format(result))
        assert genesis.EconomicModel.Staking.UnStakeFreezeDuration == int(result.get('Ret'))

        result = pip_obj.getGovernParamValue('Staking', 'MaxValidators')
        log.info('Interface getGovernParamValue result : {}'.format(result))
        assert genesis.EconomicModel.Staking.MaxValidators == int(result.get('Ret'))

        result = pip_obj.getGovernParamValue('Block', 'MaxBlockGasLimit')
        log.info('Interface getGovernParamValue result : {}'.format(result))

    def test_IN_GG_002(self, client_noconsensus_obj):
        pip_obj = client_noconsensus_obj.pip.pip
        result = pip_obj.getGovernParamValue('staking', 'MaxValidators')
        assert_code(result, 302031)
        pip_obj = client_noconsensus_obj.pip.pip
        result = pip_obj.getGovernParamValue('slashing', 'SlashBlocksReward')
        assert_code(result, 302031)
        pip_obj = client_noconsensus_obj.pip.pip
        result = pip_obj.getGovernParamValue('block', 'MaxBlockGasLimit')
        assert_code(result, 302031)

    def test_IN_GG_003(self, client_noconsensus_obj):
        pip_obj = client_noconsensus_obj.pip.pip
        result = pip_obj.getGovernParamValue('Staking', 'maxValidators')
        assert_code(result, 302031)
        pip_obj = client_noconsensus_obj.pip.pip
        result = pip_obj.getGovernParamValue('Slashing', 'slashBlocksReward')
        assert_code(result, 302031)
        pip_obj = client_noconsensus_obj.pip.pip
        result = pip_obj.getGovernParamValue('Block', 'maxValidators')
        assert_code(result, 302031)


        
