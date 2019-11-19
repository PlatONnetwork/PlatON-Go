from common.log import log
from tests.lib.utils import assert_code, wait_block_number, upload_platon, get_pledge_list
from tests.lib import Genesis
from dacite import from_dict
from tests.govern.test_voting_statistics import submittpandvote, submitcppandvote, \
    submitppandvote, submitcvpandvote, submitvpandvote
import time
import pytest
from tests.govern.test_declare_version import replace_version_declare


class TestSupportRateVoteRatePP():
    @pytest.mark.P0
    @pytest.mark.compatibility
    def test_UP_PA_001_VS_EP_002(self, new_genesis_env, client_con_list_obj):
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

    @pytest.mark.P1
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

    @pytest.mark.P1
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

    @pytest.mark.P1
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

    @pytest.mark.compatibility
    @pytest.mark.P0
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

    @pytest.mark.P1
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

    @pytest.mark.P1
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
    @pytest.mark.P1
    def test_UC_CP_001(self, new_genesis_env, client_con_list_obj):
        genesis = from_dict(data_class=Genesis, data=new_genesis_env.genesis_config)
        genesis.economicModel.gov.cancelProposalSupportRate = 0.332
        genesis.economicModel.gov.cancelProposalVoteRate = 0.751
        new_genesis_env.set_genesis(genesis.to_dict())
        new_genesis_env.deploy_all()
        submitcppandvote(client_con_list_obj[:3], [1, 2, 3])
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

    @pytest.mark.P1
    def test_UC_CP_002(self, new_genesis_env, client_con_list_obj):
        genesis = from_dict(data_class=Genesis, data=new_genesis_env.genesis_config)
        genesis.economicModel.gov.cancelProposalSupportRate = 0.334
        genesis.economicModel.gov.cancelProposalVoteRate = 0.749
        new_genesis_env.set_genesis(genesis.to_dict())
        new_genesis_env.deploy_all()
        submitcppandvote(client_con_list_obj[:3], [1, 2, 3])
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

    @pytest.mark.P1
    def test_UC_CP_003(self, new_genesis_env, client_con_list_obj):
        genesis = from_dict(data_class=Genesis, data=new_genesis_env.genesis_config)
        genesis.economicModel.gov.cancelProposalSupportRate = 0.333
        genesis.economicModel.gov.cancelProposalVoteRate = 0.751
        new_genesis_env.set_genesis(genesis.to_dict())
        new_genesis_env.deploy_all()
        submitcppandvote(client_con_list_obj[:3], [1, 2, 3])
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

    @pytest.mark.P1
    def test_UC_CP_004(self, new_genesis_env, client_con_list_obj):
        genesis = from_dict(data_class=Genesis, data=new_genesis_env.genesis_config)
        genesis.economicModel.gov.cancelProposalSupportRate = 0.334
        genesis.economicModel.gov.cancelProposalVoteRate = 0.75
        new_genesis_env.set_genesis(genesis.to_dict())
        new_genesis_env.deploy_all()
        submitcppandvote(client_con_list_obj[:3], [1, 2, 3])
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

    @pytest.mark.compatibility
    @pytest.mark.P0
    def test_UC_CP_005(self, new_genesis_env, client_con_list_obj):
        genesis = from_dict(data_class=Genesis, data=new_genesis_env.genesis_config)
        genesis.economicModel.gov.cancelProposalSupportRate = 0.332
        genesis.economicModel.gov.cancelProposalVoteRate = 0.749
        new_genesis_env.set_genesis(genesis.to_dict())
        new_genesis_env.deploy_all()
        submitcppandvote(client_con_list_obj[:3], [1, 2, 3])
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

    @pytest.mark.P1
    def test_UC_CP_006(self, new_genesis_env, client_con_list_obj):
        genesis = from_dict(data_class=Genesis, data=new_genesis_env.genesis_config)
        genesis.economicModel.gov.cancelProposalSupportRate = 0.333
        genesis.economicModel.gov.cancelProposalVoteRate = 0.749
        new_genesis_env.set_genesis(genesis.to_dict())
        new_genesis_env.deploy_all()
        submitcppandvote(client_con_list_obj[:3], [1, 2, 3])
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

    @pytest.mark.P1
    def test_UC_CP_007(self, new_genesis_env, client_con_list_obj):
        genesis = from_dict(data_class=Genesis, data=new_genesis_env.genesis_config)
        genesis.economicModel.gov.cancelProposalSupportRate = 0.332
        genesis.economicModel.gov.cancelProposalVoteRate = 0.75
        new_genesis_env.set_genesis(genesis.to_dict())
        new_genesis_env.deploy_all()
        submitcppandvote(client_con_list_obj[:3], [1, 2, 3])
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
    @pytest.mark.compatibility
    @pytest.mark.P0
    def test_UP_CA_001_VS_BL_2(self, new_genesis_env, client_con_list_obj):
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

    @pytest.mark.P1
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

    @pytest.mark.P1
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

    @pytest.mark.P1
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

    @pytest.mark.compatibility
    @pytest.mark.P0
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

    @pytest.mark.P1
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

    @pytest.mark.P1
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
    @pytest.mark.compatibility
    @pytest.mark.P1
    def test_UP_TE_001_VS_BL_3(self, new_genesis_env, client_con_list_obj):
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

    @pytest.mark.P1
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

    @pytest.mark.P1
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

    @pytest.mark.P1
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

    @pytest.mark.compatibility
    @pytest.mark.P0
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

    @pytest.mark.P1
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

    @pytest.mark.P1
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


class TestUpgradedST():
    @pytest.mark.compatibility
    @pytest.mark.P0
    def test_UV_TR_001_004_to_008_011_to_017_VS_EP_001(self, new_genesis_env, client_con_list_obj):
        new_genesis_env.deploy_all()
        pip_obj = client_con_list_obj[0].pip
        submitvpandvote(client_con_list_obj[:3])
        proposalinfo_version = pip_obj.get_effect_proposal_info_of_vote()
        log.info('Get version proposal information {}'.format(proposalinfo_version))
        wait_block_number(pip_obj.node, proposalinfo_version.get('ActiveBlock'))
        assert pip_obj.get_status_of_proposal(proposalinfo_version.get('ProposalID')) == 5
        assert pip_obj.chain_version == pip_obj.cfg.version5
        assert pip_obj.get_accuverifiers_count(proposalinfo_version.get('ProposalID')) == [4, 3, 0, 0]
        submittpandvote(client_con_list_obj[:2], 1, 2)
        submitcppandvote(client_con_list_obj[:2], [1, 2])
        proposalinfo_param = pip_obj.get_effect_proposal_info_of_vote(pip_obj.cfg.param_proposal)
        log.info('Get param proposal information {}'.format(proposalinfo_param))
        result = pip_obj.vote(pip_obj.node.node_id, proposalinfo_param.get('ProposalID'), pip_obj.cfg.vote_option_yeas,
                              pip_obj.node.staking_address, transaction_cfg=pip_obj.cfg.transaction_cfg)
        log.info('Vote param proposal result : {}'.format(result))
        assert_code(result, 0)
        result = replace_version_declare(pip_obj, pip_obj.cfg.PLATON_NEW_BIN0, pip_obj.cfg.version0)
        assert_code(result, 302028)
        result = replace_version_declare(pip_obj, pip_obj.cfg.PLATON_NEW_BIN, pip_obj.cfg.version5)
        assert_code(result, 0)
        result = replace_version_declare(pip_obj, pip_obj.cfg.PLATON_NEW_BIN4, pip_obj.cfg.version4)
        assert_code(result, 0)
        result = replace_version_declare(pip_obj, pip_obj.cfg.PLATON_NEW_BIN6, pip_obj.cfg.version6)
        assert_code(result, 0)
        result = pip_obj.pip.listProposal()
        log.info('Interface listProposal result : {}'.format(result))
        assert_code(result, 0)
        result = pip_obj.pip.getProposal(proposalinfo_version.get('ProposalID'))
        log.info('Interface getProposal result : {}'.format(result))
        assert_code(result, 0)

    @pytest.mark.P2
    def test_UV_TR_002_003_009_010(self, new_genesis_env, client_con_list_obj):
        new_genesis_env.deploy_all()
        pip_obj = client_con_list_obj[0].pip
        submitvpandvote(client_con_list_obj[:3])
        proposalinfo_version = pip_obj.get_effect_proposal_info_of_vote()
        log.info('Get version proposal information {}'.format(proposalinfo_version))
        wait_block_number(pip_obj.node, proposalinfo_version.get('ActiveBlock'))
        assert pip_obj.get_status_of_proposal(proposalinfo_version.get('ProposalID')) == 5
        assert pip_obj.chain_version == pip_obj.cfg.version5
        assert pip_obj.get_accuverifiers_count(proposalinfo_version.get('ProposalID'))

        result = pip_obj.submitVersion(pip_obj.node.node_id, str(time.time()), pip_obj.cfg.version8, 3,
                                       pip_obj.node.staking_address, transaction_cfg=pip_obj.cfg.transaction_cfg)
        log.info('Submit version proposal result : {}'.format(result))
        assert_code(result, 0)
        proposalinfo_version = pip_obj.get_effect_proposal_info_of_vote()
        log.info('Get version proposal information : {}'.format(proposalinfo_version))
        result = pip_obj.submitCancel(pip_obj.node.node_id, str(time.time()), 1, proposalinfo_version.get('ProposalID'),
                                      pip_obj.node.staking_address, transaction_cfg=pip_obj.cfg.transaction_cfg)
        log.info('Submit cancel proposal result : {}'.format(result))
        assert_code(result, 0)
        proposalinfo_cancel = pip_obj.get_effect_proposal_info_of_vote(pip_obj.cfg.cancel_proposal)
        log.info('Get version proposal information : {}'.format(proposalinfo_cancel))

        upload_platon(pip_obj.node, pip_obj.cfg.PLATON_NEW_BIN8)
        pip_obj.node.restart()

        result = pip_obj.vote(pip_obj.node.node_id, proposalinfo_version.get('ProposalID'), pip_obj.cfg.vote_option_yeas,
                              pip_obj.node.staking_address, transaction_cfg=pip_obj.cfg.transaction_cfg)
        log.info('Vote result : {}'.format(result))
        assert_code(result, 0)
        result = pip_obj.vote(pip_obj.node.node_id, proposalinfo_cancel.get('ProposalID'), pip_obj.cfg.vote_option_yeas,
                              pip_obj.node.staking_address, transaction_cfg=pip_obj.cfg.transaction_cfg)
        assert_code(result, 0)
        log.info('Node {} vote result : {}'.format(pip_obj.node.node_id, result))


class TestUpgradeVP():
    @pytest.mark.compatibility
    @pytest.mark.P0
    def test_UV_UPG_1_UV_UPG_2(self, new_genesis_env, client_con_list_obj, client_noconsensus_obj):
        new_genesis_env.deploy_all()
        pip_obj = client_con_list_obj[0].pip
        pip_obj_test = client_noconsensus_obj.pip
        address, _ = pip_obj_test.economic.account.generate_account(pip_obj_test.node.web3, 10**18 * 10000000)
        result = client_noconsensus_obj.staking.create_staking(0, address, address, amount=10**18 * 2000000,
                                                               transaction_cfg=pip_obj_test.cfg.transaction_cfg)
        log.info('Node {} staking result : {}'.format(pip_obj_test.node.node_id, result))
        pip_obj_test.economic.wait_settlement_blocknum(pip_obj_test.node)
        verifier_list = get_pledge_list(client_con_list_obj[0].ppos.getVerifierList)
        log.info('Get verifier list : {}'.format(verifier_list))
        assert pip_obj_test.node.node_id in verifier_list

        submitvpandvote(client_con_list_obj)
        proposalinfo = pip_obj.get_effect_proposal_info_of_vote()
        log.info('Get version proposal information : {}'.format(proposalinfo))
        wait_block_number(pip_obj.node, proposalinfo.get('EndVotingBlock'))
        assert_code(pip_obj.get_status_of_proposal(proposalinfo.get('ProposalID')), 4)
        validator_list = get_pledge_list(client_con_list_obj[0].ppos.getValidatorList)
        log.info('Validator list : {}'.format(validator_list))
        wait_block_number(pip_obj.node, proposalinfo.get('ActiveBlock'))

        validator_list = get_pledge_list(client_con_list_obj[0].ppos.getValidatorList)
        log.info('Validator list : {}'.format(validator_list))
        assert pip_obj_test.node.node_id not in validator_list

        assert_code(pip_obj.get_status_of_proposal(proposalinfo.get('ProposalID')), 5)
        pip_obj.economic.wait_settlement_blocknum(pip_obj.node)
        validator_list = get_pledge_list(client_con_list_obj[0].ppos.getValidatorList)
        log.info('Validator list : {}'.format(validator_list))
        assert pip_obj_test.node.node_id not in validator_list
        verifier_list = get_pledge_list(client_con_list_obj[0].ppos.getVerifierList)
        log.info('Get verifier list : {}'.format(verifier_list))
        assert pip_obj_test.node.node_id not in verifier_list
        balance_before = pip_obj.node.eth.getBalance(address, 2 * pip_obj.economic.settlement_size - 1)
        log.info('Block number {} address balace {}'.format(2 * pip_obj.economic.settlement_size - 1, balance_before))
        balance_after = pip_obj.node.eth.getBalance(address, 2 * pip_obj.economic.settlement_size)
        log.info('Block number {} address balace {}'.format(2 * pip_obj.economic.settlement_size, balance_after))
        _, staking_reward = pip_obj_test.economic.get_current_year_reward(pip_obj_test.node, verifier_num=5)
        assert balance_after - balance_before == staking_reward

    @pytest.mark.P1
    def test_UV_NO_1(self, new_genesis_env, client_con_list_obj):
        genesis = from_dict(data_class=Genesis, data=new_genesis_env.genesis_config)
        genesis.economicModel.gov.versionProposalSupportRate = 0.249
        genesis.economicModel.slashing.slashBlocksReward = 0
        new_genesis_env.set_genesis(genesis.to_dict())
        new_genesis_env.deploy_all()
        pip_obj = client_con_list_obj[0].pip
        submitvpandvote([client_con_list_obj[0]])
        proposalinfo = pip_obj.get_effect_proposal_info_of_vote()
        log.info('Get version proposal infomation  {}'.format(proposalinfo))
        wait_block_number(pip_obj.node, proposalinfo.get('EndVotingBlock'))
        assert pip_obj.get_accuverifiers_count(proposalinfo.get('ProposalID')) == [4, 1, 0, 0]
        assert pip_obj.get_accu_verifiers_of_proposal(proposalinfo.get('ProposalID')) == len(client_con_list_obj)
        assert pip_obj.get_yeas_of_proposal(proposalinfo.get('ProposalID')) == 1
        assert pip_obj.get_nays_of_proposal(proposalinfo.get('ProposalID')) == 0
        assert pip_obj.get_abstentions_of_proposal(proposalinfo.get('ProposalID')) == 0
        assert_code(pip_obj.get_status_of_proposal(proposalinfo.get('ProposalID')), 3)

    @pytest.mark.P1
    def test_UV_UP_1(self, new_genesis_env, client_con_list_obj):
        genesis = from_dict(data_class=Genesis, data=new_genesis_env.genesis_config)
        genesis.economicModel.gov.versionProposalSupportRate = 0.25
        genesis.economicModel.slashing.slashBlocksReward = 0
        new_genesis_env.set_genesis(genesis.to_dict())
        new_genesis_env.deploy_all()
        pip_obj = client_con_list_obj[0].pip
        submitvpandvote([client_con_list_obj[0]])
        proposalinfo = pip_obj.get_effect_proposal_info_of_vote()
        log.info('Get version proposal infomation  {}'.format(proposalinfo))
        wait_block_number(pip_obj.node, proposalinfo.get('EndVotingBlock'))
        assert pip_obj.get_accuverifiers_count(proposalinfo.get('ProposalID')) == [4, 1, 0, 0]
        assert pip_obj.get_accu_verifiers_of_proposal(proposalinfo.get('ProposalID')) == len(client_con_list_obj)
        assert pip_obj.get_yeas_of_proposal(proposalinfo.get('ProposalID')) == 1
        assert pip_obj.get_nays_of_proposal(proposalinfo.get('ProposalID')) == 0
        assert pip_obj.get_abstentions_of_proposal(proposalinfo.get('ProposalID')) == 0
        assert_code(pip_obj.get_status_of_proposal(proposalinfo.get('ProposalID')), 4)
        wait_block_number(pip_obj.node, proposalinfo.get('ActiveBlock'))
        assert_code(pip_obj.get_status_of_proposal(proposalinfo.get('ProposalID')), 5)
