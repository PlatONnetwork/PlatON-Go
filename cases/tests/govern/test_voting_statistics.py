import pytest
from common.log import log
from common.key import mock_duplicate_sign
import time
from tests.lib.utils import assert_code, wait_block_number, upload_platon
from tests.lib.client import Client
from tests.lib import Genesis
from dacite import from_dict
from tests.govern.conftest import param_proposal_vote

class TestVotingStatistics():
    def submitppandvote(self, client_list_obj, *args):
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
        genesis.EconomicModel.Gov.ParamProposalVote_DurationSeconds = 0
        new_genesis_env.set_genesis(genesis.to_dict())
        new_genesis_env.deploy_all()
        self.submitppandvote(client_con_list_obj[0:-1], 1, 2, 3)
        proposalinfo = client_con_list_obj[0].pip.get_effect_proposal_info_of_vote(client_con_list_obj[0].pip.cfg.param_proposal)
        log.info('Param proposal information {}'.format(proposalinfo))
        self.createstaking(client_noc_list_obj[:3])
        # client_con_list_obj[0].economic.wait_settlement_blocknum(client_con_list_obj[0].node)
        result = client_con_list_obj[0].pip.get_accuverifiers_count(proposalinfo.get('ProposalID'))
        log.info('Get proposal vote infomation {}'.format(result))
        assert result == [4, 1, 1, 1]

    def test_VS_EP_005(self, new_genesis_env, client_con_list_obj, client_noc_list_obj):
        genesis = from_dict(data_class=Genesis, data=new_genesis_env.genesis_config)
        genesis.EconomicModel.Gov.ParamProposalVote_DurationSeconds = 160
        new_genesis_env.set_genesis(genesis.to_dict())
        new_genesis_env.deploy_all()
        pip_obj = client_con_list_obj[0].pip
        self.submitppandvote(client_con_list_obj[:2], 1, 2)
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
        genesis.EconomicModel.Gov.ParamProposalVote_DurationSeconds = 320
        new_genesis_env.set_genesis(genesis.to_dict())
        new_genesis_env.deploy_all()
        self.submitppandvote(client_con_list_obj[:1], 1)
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
        genesis.EconomicModel.Gov.ParamProposalVote_DurationSeconds = 0
        genesis.EconomicModel.Gov.ParamProposal_SupportRate = 0.5
        genesis.EconomicModel.Gov.ParamProposal_VoteRate = 0.5
        new_genesis_env.set_genesis(genesis.to_dict())
        new_genesis_env.deploy_all()
        self.submitppandvote(client_con_list_obj[:2], 1, 1)
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
        genesis.EconomicModel.Gov.ParamProposalVote_DurationSeconds = 0
        genesis.EconomicModel.Gov.ParamProposal_SupportRate = 0.5
        genesis.EconomicModel.Gov.ParamProposal_VoteRate = 0.5
        new_genesis_env.set_genesis(genesis.to_dict())
        new_genesis_env.deploy_all()
        self.submitppandvote(client_con_list_obj[:2], 1, 1)
        pip_obj = client_con_list_obj[0].pip
        proposalinfo = pip_obj.get_effect_proposal_info_of_vote(pip_obj.cfg.param_proposal)
        log.info('Get param proposal information {}'.format(proposalinfo))
        # wait_block_number(pip_obj.node, proposalinfo.get('EndVotingBlock'))
        report_information = mock_duplicate_sign(1, pip_obj.node.nodekey, pip_obj.node.blsprikey,
                                                 proposalinfo.get('EndVotingBlock') - 10)
        log.info("Report information: {}".format(report_information))
        address, _ = pip_obj.economic.account.generate_account(pip_obj.node.web3, 10**18 * 1000)
        result = client_con_list_obj[0].duplicatesign.reportDuplicateSign(1, report_information, address)
        assert_code(result, 0)
        assert client_con_list_obj[1].pip.get_yeas_of_proposal(proposalinfo.get('ProposalID')) == 2