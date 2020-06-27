from common.log import log
from tests.lib.utils import assert_code, wait_block_number, upload_platon, get_pledge_list
from tests.lib import Genesis
from dacite import from_dict
from tests.govern.test_voting_statistics import submittpandvote, submitcppandvote, \
    submitppandvote, submitcvpandvote, submitvpandvote
import time
from tests.govern.conftest import verifier_node_version
import pytest
import allure
from tests.govern.test_declare_version import replace_version_declare


def verify_proposal_status(clients, proposaltype, status):
    pip = clients[0].pip
    proposalinfo = pip.get_effect_proposal_info_of_vote(proposaltype)
    log.info('Get proposal information {}'.format(proposalinfo))
    assert pip.get_accuverifiers_count(proposalinfo.get('ProposalID')) == [4, 1, 1, 1]
    wait_block_number(clients[1].node, proposalinfo.get('EndVotingBlock'))
    assert_code(pip.get_yeas_of_proposal(proposalinfo.get('ProposalID')), 1)
    assert_code(pip.get_nays_of_proposal(proposalinfo.get('ProposalID')), 1)
    assert_code(pip.get_abstentions_of_proposal(proposalinfo.get('ProposalID')), 1)
    assert_code(pip.get_accu_verifiers_of_proposal(proposalinfo.get('ProposalID')), len(clients))
    assert_code(pip.get_status_of_proposal(proposalinfo.get('ProposalID')), status)


def update_setting_rate(new_genesis_env, proposaltype, *args):
    genesis = from_dict(data_class=Genesis, data=new_genesis_env.genesis_config)
    if proposaltype == 3:
        genesis.economicModel.gov.paramProposalVoteDurationSeconds = args[0]
        genesis.economicModel.gov.paramProposalSupportRate = args[1]
        genesis.economicModel.gov.paramProposalVoteRate = args[2]

    elif proposaltype == 4:
        genesis.economicModel.gov.cancelProposalSupportRate = args[0]
        genesis.economicModel.gov.cancelProposalVoteRate = args[1]

    elif proposaltype == 1:
        genesis.economicModel.gov.textProposalVoteDurationSeconds = args[0]
        genesis.economicModel.gov.textProposalSupportRate = args[1]
        genesis.economicModel.gov.textProposalVoteRate = args[2]

    elif proposaltype == 2:
        genesis.economicModel.gov.versionProposalSupportRate = args[0]
        genesis.economicModel.slashing.slashBlocksReward = 0
    else:
        raise ValueError('Prposal type error')
    new_genesis_env.set_genesis(genesis.to_dict())
    new_genesis_env.deploy_all()


class TestSupportRateVoteRatePP:
    @pytest.mark.P0
    @pytest.mark.compatibility
    @allure.title('Parameter proposal statistical function verification')
    def test_UP_PA_001_VS_EP_002(self, new_genesis_env, clients_consensus):
        update_setting_rate(new_genesis_env, 3, 0, 3320, 7510)
        submitppandvote(clients_consensus[:3], 1, 2, 3)
        verify_proposal_status(clients_consensus, proposaltype=3, status=3)

    @pytest.mark.P1
    @allure.title('Parameter proposal statistical function verification')
    def test_UP_PA_002(self, new_genesis_env, clients_consensus):
        update_setting_rate(new_genesis_env, 3, 0, 3340, 7490)
        submitppandvote(clients_consensus[:3], 1, 2, 3)
        verify_proposal_status(clients_consensus, proposaltype=3, status=3)

    @pytest.mark.P1
    @allure.title('Parameter proposal statistical function verification')
    def test_UP_PA_003(self, new_genesis_env, clients_consensus):
        update_setting_rate(new_genesis_env, 3, 0, 3330, 7510)
        submitppandvote(clients_consensus[:3], 1, 2, 3)
        verify_proposal_status(clients_consensus, proposaltype=3, status=3)

    @pytest.mark.P1
    @allure.title('Parameter proposal statistical function verification')
    def test_UP_PA_004(self, new_genesis_env, clients_consensus):
        update_setting_rate(new_genesis_env, 3, 0, 3340, 7500)
        submitppandvote(clients_consensus[:3], 1, 2, 3)
        verify_proposal_status(clients_consensus, proposaltype=3, status=3)

    @pytest.mark.compatibility
    @pytest.mark.P0
    @allure.title('Parameter proposal statistical function verification')
    def test_UP_PA_005(self, new_genesis_env, clients_consensus):
        update_setting_rate(new_genesis_env, 3, 0, 3320, 7490)
        submitppandvote(clients_consensus[:3], 1, 2, 3)
        verify_proposal_status(clients_consensus, proposaltype=3, status=2)

    @pytest.mark.P1
    @allure.title('Parameter proposal statistical function verification')
    def test_UP_PA_006(self, new_genesis_env, clients_consensus):
        update_setting_rate(new_genesis_env, 3, 0, 3330, 7490)
        submitppandvote(clients_consensus[:3], 1, 2, 3)
        verify_proposal_status(clients_consensus, proposaltype=3, status=2)

    @pytest.mark.P1
    @allure.title('Parameter proposal statistical function verification')
    def test_UP_PA_007(self, new_genesis_env, clients_consensus):
        update_setting_rate(new_genesis_env, 3, 0, 3320, 7500)
        submitppandvote(clients_consensus[:3], 1, 2, 3)
        verify_proposal_status(clients_consensus, proposaltype=3, status=3)


class TestSupportRateVoteRateCPP:
    @pytest.mark.P1
    @allure.title('Cancel parameter proposal statistical function verification')
    def test_UC_CP_001(self, new_genesis_env, clients_consensus):
        update_setting_rate(new_genesis_env, 4, 3320, 7510)
        submitcppandvote(clients_consensus[:3], [1, 2, 3])
        verify_proposal_status(clients_consensus, proposaltype=4, status=3)

    @pytest.mark.P1
    @allure.title('Cancel parameter proposal statistical function verification')
    def test_UC_CP_002(self, new_genesis_env, clients_consensus):
        update_setting_rate(new_genesis_env, 4, 3340, 7490)
        submitcppandvote(clients_consensus[:3], [1, 2, 3])
        verify_proposal_status(clients_consensus, proposaltype=4, status=3)

    @pytest.mark.P1
    @allure.title('Cancel parameter proposal statistical function verification')
    def test_UC_CP_003(self, new_genesis_env, clients_consensus):
        update_setting_rate(new_genesis_env, 4, 3330, 7510)
        submitcppandvote(clients_consensus[:3], [1, 2, 3])
        verify_proposal_status(clients_consensus, proposaltype=4, status=3)

    @pytest.mark.P1
    @allure.title('Cancel parameter proposal statistical function verification')
    def test_UC_CP_004(self, new_genesis_env, clients_consensus):
        update_setting_rate(new_genesis_env, 4, 3340, 7500)
        submitcppandvote(clients_consensus[:3], [1, 2, 3])
        verify_proposal_status(clients_consensus, proposaltype=4, status=3)

    @pytest.mark.compatibility
    @pytest.mark.P0
    @allure.title('Cancel parameter proposal statistical function verification')
    def test_UC_CP_005(self, new_genesis_env, clients_consensus):
        update_setting_rate(new_genesis_env, 4, 3320, 7490)
        submitcppandvote(clients_consensus[:3], [1, 2, 3])
        verify_proposal_status(clients_consensus, proposaltype=4, status=2)

    @pytest.mark.P1
    @allure.title('Cancel parameter proposal statistical function verification')
    def test_UC_CP_006(self, new_genesis_env, clients_consensus):
        update_setting_rate(new_genesis_env, 4, 3330, 7490)
        submitcppandvote(clients_consensus[:3], [1, 2, 3])
        verify_proposal_status(clients_consensus, proposaltype=4, status=2)

    @pytest.mark.P1
    @allure.title('Cancel parameter proposal statistical function verification')
    def test_UC_CP_007(self, new_genesis_env, clients_consensus):
        update_setting_rate(new_genesis_env, 4, 3320, 7500)
        submitcppandvote(clients_consensus[:3], [1, 2, 3])
        verify_proposal_status(clients_consensus, proposaltype=4, status=3)


class TestSupportRateVoteRateCVP:
    @pytest.mark.compatibility
    @pytest.mark.P0
    @allure.title('Cancel version proposal statistical function verification')
    def test_UP_CA_001_VS_BL_2(self, new_genesis_env, clients_consensus):
        update_setting_rate(new_genesis_env, 4, 3320, 7510)
        submitcvpandvote(clients_consensus[:3], 1, 2, 3)
        verify_proposal_status(clients_consensus, proposaltype=4, status=3)

    @pytest.mark.P1
    @allure.title('Cancel version proposal statistical function verification')
    def test_UP_CA_002(self, new_genesis_env, clients_consensus):
        update_setting_rate(new_genesis_env, 4, 3340, 7490)
        submitcvpandvote(clients_consensus[:3], 1, 2, 3)
        verify_proposal_status(clients_consensus, proposaltype=4, status=3)

    @pytest.mark.P1
    @allure.title('Cancel version proposal statistical function verification')
    def test_UP_CA_003(self, new_genesis_env, clients_consensus):
        update_setting_rate(new_genesis_env, 4, 3330, 7510)
        submitcvpandvote(clients_consensus[:3], 1, 2, 3)
        verify_proposal_status(clients_consensus, proposaltype=4, status=3)

    @pytest.mark.P1
    @allure.title('Cancel version proposal statistical function verification')
    def test_UP_CA_004(self, new_genesis_env, clients_consensus):
        update_setting_rate(new_genesis_env, 4, 3340, 7500)
        submitcvpandvote(clients_consensus[:3], 1, 2, 3)
        verify_proposal_status(clients_consensus, proposaltype=4, status=3)

    @pytest.mark.compatibility
    @pytest.mark.P0
    @allure.title('Cancel version proposal statistical function verification')
    def test_UP_CA_005(self, new_genesis_env, clients_consensus):
        update_setting_rate(new_genesis_env, 4, 3320, 7490)
        submitcvpandvote(clients_consensus[:3], 1, 2, 3)
        verify_proposal_status(clients_consensus, proposaltype=4, status=2)

    @pytest.mark.P1
    @allure.title('Cancel version proposal statistical function verification')
    def test_UP_CA_006(self, new_genesis_env, clients_consensus):
        update_setting_rate(new_genesis_env, 4, 3330, 7490)
        submitcvpandvote(clients_consensus[:3], 1, 2, 3)
        verify_proposal_status(clients_consensus, proposaltype=4, status=2)

    @pytest.mark.P1
    @allure.title('Cancel version proposal statistical function verification')
    def test_UP_CA_007(self, new_genesis_env, clients_consensus):
        update_setting_rate(new_genesis_env, 4, 3320, 7500)
        submitcvpandvote(clients_consensus[:3], 1, 2, 3)
        verify_proposal_status(clients_consensus, proposaltype=4, status=3)


class TestSupportRateVoteRateTP:
    @pytest.mark.compatibility
    @pytest.mark.P1
    @allure.title('Text proposal statistical function verification')
    def test_UP_TE_001_VS_BL_3(self, new_genesis_env, clients_consensus):
        update_setting_rate(new_genesis_env, 1, 40, 3320, 7500)
        submittpandvote(clients_consensus[:3], 1, 2, 3)
        verify_proposal_status(clients_consensus, proposaltype=1, status=3)

    @pytest.mark.P1
    @allure.title('Text proposal statistical function verification')
    def test_UP_TE_002(self, new_genesis_env, clients_consensus):
        update_setting_rate(new_genesis_env, 1, 40, 3340, 7490)
        submittpandvote(clients_consensus[:3], 1, 2, 3)
        verify_proposal_status(clients_consensus, proposaltype=1, status=3)

    @pytest.mark.P1
    @allure.title('Text proposal statistical function verification')
    def test_UP_TE_003(self, new_genesis_env, clients_consensus):
        update_setting_rate(new_genesis_env, 1, 40, 3330, 7510)
        submittpandvote(clients_consensus[:3], 1, 2, 3)
        verify_proposal_status(clients_consensus, proposaltype=1, status=3)

    @pytest.mark.P1
    @allure.title('Text proposal statistical function verification')
    def test_UP_TE_004(self, new_genesis_env, clients_consensus):
        update_setting_rate(new_genesis_env, 1, 40, 3340, 7500)
        submittpandvote(clients_consensus[:3], 1, 2, 3)
        verify_proposal_status(clients_consensus, proposaltype=1, status=3)

    @pytest.mark.compatibility
    @pytest.mark.P0
    @allure.title('Text proposal statistical function verification')
    def test_UP_TE_005(self, new_genesis_env, clients_consensus):
        update_setting_rate(new_genesis_env, 1, 40, 3320, 7490)
        submittpandvote(clients_consensus[:3], 1, 2, 3)
        verify_proposal_status(clients_consensus, proposaltype=1, status=2)

    @pytest.mark.P1
    @allure.title('Text proposal statistical function verification')
    def test_UP_TE_006(self, new_genesis_env, clients_consensus):
        update_setting_rate(new_genesis_env, 1, 40, 3330, 7490)
        submittpandvote(clients_consensus[:3], 1, 2, 3)
        verify_proposal_status(clients_consensus, proposaltype=1, status=2)

    @pytest.mark.P1
    @allure.title('Text proposal statistical function verification')
    def test_UP_TE_007(self, new_genesis_env, clients_consensus):
        update_setting_rate(new_genesis_env, 1, 40, 3330, 7500)
        submittpandvote(clients_consensus[:3], 1, 2, 3)
        verify_proposal_status(clients_consensus, proposaltype=1, status=3)


class TestUpgradedST:
    @pytest.mark.compatibility
    @pytest.mark.P0
    @allure.title('Chain upgrade completed, transaction function verification')
    def test_UV_TR_001_004_to_008_011_to_017_VS_EP_001(self, new_genesis_env, clients_consensus):
        new_genesis_env.deploy_all()
        pip = clients_consensus[0].pip
        submitvpandvote(clients_consensus[:3])
        proposalinfo_version = pip.get_effect_proposal_info_of_vote()
        log.info('Get version proposal information {}'.format(proposalinfo_version))
        wait_block_number(pip.node, proposalinfo_version.get('ActiveBlock'))
        assert pip.get_status_of_proposal(proposalinfo_version.get('ProposalID')) == 5
        assert pip.chain_version == pip.cfg.version5
        assert pip.get_accuverifiers_count(proposalinfo_version.get('ProposalID')) == [4, 3, 0, 0]
        submittpandvote(clients_consensus[:2], 1, 2)
        submitcppandvote(clients_consensus[:2], [1, 2])
        proposalinfo_param = pip.get_effect_proposal_info_of_vote(pip.cfg.param_proposal)
        log.info('Get param proposal information {}'.format(proposalinfo_param))
        result = pip.vote(pip.node.node_id, proposalinfo_param.get('ProposalID'), pip.cfg.vote_option_yeas,
                          pip.node.staking_address, transaction_cfg=pip.cfg.transaction_cfg)
        log.info('Vote param proposal result : {}'.format(result))
        assert_code(result, 0)
        result = replace_version_declare(pip, pip.cfg.PLATON_NEW_BIN0, pip.cfg.version0)
        assert_code(result, 302028)
        result = replace_version_declare(pip, pip.cfg.PLATON_NEW_BIN, pip.cfg.version5)
        assert_code(result, 0)
        verifier_node_version(pip, pip.cfg.version5)
        result = replace_version_declare(pip, pip.cfg.PLATON_NEW_BIN4, pip.cfg.version4)
        assert_code(result, 0)
        verifier_node_version(pip, pip.cfg.version4)
        result = replace_version_declare(pip, pip.cfg.PLATON_NEW_BIN6, pip.cfg.version6)
        assert_code(result, 0)
        verifier_node_version(pip, pip.cfg.version6)
        result = pip.pip.listProposal()
        log.info('Interface listProposal result : {}'.format(result))
        assert_code(result, 0)
        result = pip.pip.getProposal(proposalinfo_version.get('ProposalID'))
        log.info('Interface getProposal result : {}'.format(result))
        assert_code(result, 0)

    @pytest.mark.P2
    @allure.title('Chain upgrade completed, transaction function verification')
    def test_UV_TR_002_003_009_010(self, new_genesis_env, clients_consensus):
        new_genesis_env.deploy_all()
        pip = clients_consensus[0].pip
        submitvpandvote(clients_consensus[:3])
        proposalinfo_version = pip.get_effect_proposal_info_of_vote()
        log.info('Get version proposal information {}'.format(proposalinfo_version))
        wait_block_number(pip.node, proposalinfo_version.get('ActiveBlock'))
        assert pip.get_status_of_proposal(proposalinfo_version.get('ProposalID')) == 5
        assert pip.chain_version == pip.cfg.version5
        assert pip.get_accuverifiers_count(proposalinfo_version.get('ProposalID'))

        result = pip.submitVersion(pip.node.node_id, str(time.time()), pip.cfg.version8, 3,
                                   pip.node.staking_address, transaction_cfg=pip.cfg.transaction_cfg)
        log.info('Submit version proposal result : {}'.format(result))
        assert_code(result, 0)
        proposalinfo_version = pip.get_effect_proposal_info_of_vote()
        log.info('Get version proposal information : {}'.format(proposalinfo_version))
        result = pip.submitCancel(pip.node.node_id, str(time.time()), 1, proposalinfo_version.get('ProposalID'),
                                  pip.node.staking_address, transaction_cfg=pip.cfg.transaction_cfg)
        log.info('Submit cancel proposal result : {}'.format(result))
        assert_code(result, 0)
        proposalinfo_cancel = pip.get_effect_proposal_info_of_vote(pip.cfg.cancel_proposal)
        log.info('Get version proposal information : {}'.format(proposalinfo_cancel))

        upload_platon(pip.node, pip.cfg.PLATON_NEW_BIN8)
        pip.node.restart()

        result = pip.vote(pip.node.node_id, proposalinfo_version.get('ProposalID'), pip.cfg.vote_option_yeas,
                          pip.node.staking_address, transaction_cfg=pip.cfg.transaction_cfg)
        log.info('Vote result : {}'.format(result))
        assert_code(result, 0)
        result = pip.vote(pip.node.node_id, proposalinfo_cancel.get('ProposalID'), pip.cfg.vote_option_yeas,
                          pip.node.staking_address, transaction_cfg=pip.cfg.transaction_cfg)
        assert_code(result, 0)
        log.info('Node {} vote result : {}'.format(pip.node.node_id, result))


class TestUpgradeVP:
    def calculate_version(self, version):
        ver_byte = (version).to_bytes(length=4, byteorder='big', signed=False)
        new_ver3 = (0).to_bytes(length=1, byteorder='big', signed=False)
        new_version_byte = ver_byte[0:3] + new_ver3
        return int.from_bytes(new_version_byte, byteorder='big', signed=False)

    @pytest.mark.compatibility
    @pytest.mark.P0
    @allure.title('Version proposal statistical function verification')
    def test_UV_UPG_1_UV_UPG_2(self, new_genesis_env, clients_consensus, client_noconsensus):
        new_genesis_env.deploy_all()
        pip = clients_consensus[0].pip
        pip_test = client_noconsensus.pip
        address, _ = pip_test.economic.account.generate_account(pip_test.node.web3, 10**18 * 10000000)
        result = client_noconsensus.staking.create_staking(0, address, address, amount=10 ** 18 * 2000000,
                                                           transaction_cfg=pip_test.cfg.transaction_cfg)
        log.info('Node {} staking result : {}'.format(pip_test.node.node_id, result))
        programversion = client_noconsensus.staking.get_version()
        assert_code(programversion, pip.cfg.version0)
        pip_test.economic.wait_settlement_blocknum(pip_test.node)
        log.info(f'blocknem ====== {pip_test.node.eth.blockNumber}')
        verifier_list = get_pledge_list(clients_consensus[0].ppos.getVerifierList)
        log.info('Get verifier list : {}'.format(verifier_list))
        assert pip_test.node.node_id in verifier_list

        submitvpandvote(clients_consensus)
        programversion = clients_consensus[0].staking.get_version()
        assert_code(programversion, pip.cfg.version0)
        proposalinfo = pip.get_effect_proposal_info_of_vote()
        log.info('Get version proposal information : {}'.format(proposalinfo))
        wait_block_number(pip.node, proposalinfo.get('EndVotingBlock') - 1)
        validator_list = get_pledge_list(clients_consensus[0].ppos.getValidatorList)
        log.info('Validator list =====: {}'.format(validator_list))

        wait_block_number(pip.node, proposalinfo.get('EndVotingBlock'))
        assert_code(pip.get_status_of_proposal(proposalinfo.get('ProposalID')), 4)
        validator_list = get_pledge_list(clients_consensus[0].ppos.getValidatorList)
        log.info('Validator list : {}'.format(validator_list))
        wait_block_number(pip.node, proposalinfo.get('ActiveBlock'))
        log.info(f'blocknem ====== {pip_test.node.eth.blockNumber}')

        validator_list = get_pledge_list(clients_consensus[0].ppos.getValidatorList)
        log.info('Validator list : {}'.format(validator_list))
        assert pip_test.node.node_id not in validator_list

        assert_code(pip.get_status_of_proposal(proposalinfo.get('ProposalID')), 5)
        pip.economic.wait_settlement_blocknum(pip.node)
        validator_list = get_pledge_list(clients_consensus[0].ppos.getValidatorList)
        log.info('Validator list : {}'.format(validator_list))
        assert pip_test.node.node_id not in validator_list
        verifier_list = get_pledge_list(clients_consensus[0].ppos.getVerifierList)
        log.info('Get verifier list : {}'.format(verifier_list))
        assert pip_test.node.node_id not in verifier_list
        balance_before = pip.node.eth.getBalance(address, 2 * pip.economic.settlement_size - 1)
        log.info('Block number {} address balace {}'.format(2 * pip.economic.settlement_size - 1, balance_before))
        balance_after = pip.node.eth.getBalance(address, 2 * pip.economic.settlement_size)
        log.info('Block number {} address balace {}'.format(2 * pip.economic.settlement_size, balance_after))
        _, staking_reward = pip_test.economic.get_current_year_reward(pip_test.node, verifier_num=5)
        log.info('Staking reward : {}'.format(staking_reward))
        assert balance_after - balance_before == staking_reward

    @pytest.mark.P0
    def test_UV_UPG_2(self, new_genesis_env, clients_consensus, client_noconsensus):
        new_genesis_env.deploy_all()
        pip = clients_consensus[0].pip
        pip_test = client_noconsensus.pip
        address, _ = pip_test.economic.account.generate_account(pip_test.node.web3, 10 ** 18 * 10000000)
        result = client_noconsensus.staking.create_staking(0, address, address, amount=10 ** 18 * 2000000,
                                                           transaction_cfg=pip_test.cfg.transaction_cfg)
        log.info('Node {} staking result : {}'.format(pip_test.node.node_id, result))
        programversion = client_noconsensus.staking.get_version()
        assert_code(programversion, pip.cfg.version0)
        pip_test.economic.wait_settlement_blocknum(pip_test.node)
        verifier_list = get_pledge_list(clients_consensus[0].ppos.getVerifierList)
        log.info('Get verifier list : {}'.format(verifier_list))
        assert pip_test.node.node_id in verifier_list

        submitvpandvote(clients_consensus)
        programversion = clients_consensus[0].staking.get_version()
        assert_code(programversion, pip.cfg.version0)
        proposalinfo = pip.get_effect_proposal_info_of_vote()
        log.info('Get version proposal information : {}'.format(proposalinfo))
        replace_version_declare(pip_test, pip_test.cfg.PLATON_NEW_BIN, pip_test.cfg.version5)
        assert_code(result, 0)
        programversion = client_noconsensus.staking.get_version()
        assert_code(programversion, pip.cfg.version0)
        wait_block_number(pip.node, proposalinfo.get('EndVotingBlock'))
        verifier_node_version(pip, proposalinfo.get('NewVersion'))
        assert_code(pip.get_status_of_proposal(proposalinfo.get('ProposalID')), 4)
        validator_list = get_pledge_list(clients_consensus[0].ppos.getValidatorList)
        log.info('Validator list : {}'.format(validator_list))
        wait_block_number(pip.node, proposalinfo.get('ActiveBlock'))

        validator_list = get_pledge_list(clients_consensus[0].ppos.getValidatorList)
        log.info('Validator list : {}'.format(validator_list))
        assert pip_test.node.node_id in validator_list

        assert_code(pip.get_status_of_proposal(proposalinfo.get('ProposalID')), 5)
        pip.economic.wait_settlement_blocknum(pip.node)
        validator_list = get_pledge_list(clients_consensus[0].ppos.getValidatorList)
        log.info('Validator list : {}'.format(validator_list))
        assert pip_test.node.node_id in validator_list
        verifier_list = get_pledge_list(clients_consensus[0].ppos.getVerifierList)
        log.info('Get verifier list : {}'.format(verifier_list))
        assert pip_test.node.node_id in verifier_list
        programversion = clients_consensus[0].staking.get_version()
        assert_code(programversion, pip.cfg.version5)
        programversion = client_noconsensus.staking.get_version()
        assert_code(programversion, pip_test.cfg.version5)

    @pytest.mark.P1
    @allure.title('Version proposal statistical function verification')
    def test_UV_NO_1(self, new_genesis_env, clients_consensus):
        update_setting_rate(new_genesis_env, 2, 2501)
        pip = clients_consensus[0].pip
        submitvpandvote([clients_consensus[0]])
        node_version = verifier_node_version(pip)
        proposalinfo = pip.get_effect_proposal_info_of_vote()
        log.info('Get version proposal infomation  {}'.format(proposalinfo))
        wait_block_number(pip.node, proposalinfo.get('EndVotingBlock'))
        verifier_node_version(pip, node_version)
        assert pip.get_accuverifiers_count(proposalinfo.get('ProposalID')) == [4, 1, 0, 0]
        assert pip.get_accu_verifiers_of_proposal(proposalinfo.get('ProposalID')) == len(clients_consensus)
        assert pip.get_yeas_of_proposal(proposalinfo.get('ProposalID')) == 1
        assert pip.get_nays_of_proposal(proposalinfo.get('ProposalID')) == 0
        assert pip.get_abstentions_of_proposal(proposalinfo.get('ProposalID')) == 0
        assert_code(pip.get_status_of_proposal(proposalinfo.get('ProposalID')), 3)

    @pytest.mark.P1
    @allure.title('Version proposal statistical function verification')
    def test_UV_UP_1(self, new_genesis_env, clients_consensus):
        update_setting_rate(new_genesis_env, 2, 2500)
        pip = clients_consensus[0].pip
        submitvpandvote([clients_consensus[0]])
        proposalinfo = pip.get_effect_proposal_info_of_vote()
        log.info('Get version proposal infomation  {}'.format(proposalinfo))
        wait_block_number(pip.node, proposalinfo.get('EndVotingBlock'))
        verifier_node_version(pip, proposalinfo.get('NewVersion'))
        assert pip.get_accuverifiers_count(proposalinfo.get('ProposalID')) == [4, 1, 0, 0]
        assert pip.get_accu_verifiers_of_proposal(proposalinfo.get('ProposalID')) == len(clients_consensus)
        assert pip.get_yeas_of_proposal(proposalinfo.get('ProposalID')) == 1
        assert pip.get_nays_of_proposal(proposalinfo.get('ProposalID')) == 0
        assert pip.get_abstentions_of_proposal(proposalinfo.get('ProposalID')) == 0
        assert_code(pip.get_status_of_proposal(proposalinfo.get('ProposalID')), 4)
        wait_block_number(pip.node, proposalinfo.get('ActiveBlock'))
        assert_code(pip.get_status_of_proposal(proposalinfo.get('ProposalID')), 5)

    def test_1(self, new_genesis_env, clients_consensus):
        pip = clients_consensus[-1].pip
        submitvpandvote(clients_consensus[0:2])
        replace_version_declare(pip, pip.cfg.PLATON_NEW_BIN, pip.cfg.version5)
        proposalinfo = pip.get_effect_proposal_info_of_vote()
        wait_block_number(pip.node, proposalinfo.get('EndVotingBlock'))
        submitvpandvote(clients_consensus[:3])
        proposalinfo = pip.get_effect_proposal_info_of_vote()
        wait_block_number(pip.node, proposalinfo.get('EndVotingBlock'))
        verifier_list = get_pledge_list(clients_consensus[0].ppos.getVerifierList)
        log.info(verifier_list)
        assert pip.node.node_id in verifier_list

    def test_2(self, new_genesis_env, clients_consensus):
        new_genesis_env.deploy_all()
        pip = clients_consensus[0].pip
        submitvpandvote(clients_consensus[0:2])
        # replace_version_declare(pip, pip.cfg.PLATON_NEW_BIN, pip.cfg.version5)
        proposalinfo = pip.get_effect_proposal_info_of_vote()
        wait_block_number(pip.node, proposalinfo.get('EndVotingBlock'))
        assert_code(pip.get_status_of_proposal(proposalinfo.get('ProposalID')), 3)
        submitvpandvote(clients_consensus[1:4])
        proposalinfo = pip.get_effect_proposal_info_of_vote()
        wait_block_number(pip.node, proposalinfo.get('EndVotingBlock'))
        assert_code(pip.get_status_of_proposal(proposalinfo.get('ProposalID')), 4)
        pip.economic.wait_consensus_blocknum(pip.node)
        validator = get_pledge_list(clients_consensus[0].ppos.getValidatorList)
        log.info(validator)
        assert pip.node.node_id not in validator
        programversion = clients_consensus[0].staking.get_version()
        log.info(programversion)

        programversion = clients_consensus[1].staking.get_version()
        log.info(programversion)

        programversion = clients_consensus[2].staking.get_version()
        log.info(programversion)

        programversion = clients_consensus[3].staking.get_version()
        log.info(programversion)

        verifier_list = get_pledge_list(clients_consensus[0].ppos.getVerifierList)
        log.info(verifier_list)
