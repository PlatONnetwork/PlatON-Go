from common.log import log
from tests.lib.utils import upload_platon, assert_code, wait_block_number, get_pledge_list
import pytest
import allure
import time, struct
from tests.govern.test_voting_statistics import submitvpandvote, createstaking, version_proposal_vote
from tests.lib import Genesis
from tests.lib.client import get_client_by_nodeid, get_clients_by_nodeid
from dacite import from_dict
from tests.govern.conftest import verifier_node_version


@pytest.fixture()
def large_version_proposal_pips(all_clients):
    '''
    get verifier Client object list
    :param global_test_env:
    :return:
    '''
    verifier_list = get_pledge_list(all_clients[0].ppos.getVerifierList)
    log.info('verifier list {}'.format(verifier_list))
    pip = get_client_by_nodeid(verifier_list[0], all_clients).pip
    if pip.chain_version != pip.cfg.version0:
        pip.economic.env.deploy_all()
    if pip.is_exist_effective_proposal():
        proposalinfo = pip.get_effect_proposal_info_of_vote()
        log.info('proprosalinfo : {}'.format(proposalinfo))
        if proposalinfo.get('EndVotingBlock') - pip.node.block_number > 2 * pip.economic.consensus_size \
                and proposalinfo.get('NewVersion') == pip.cfg.version8:
            verifier_list = get_pledge_list(all_clients[0].ppos.getVerifierList)
            log.info('verifierlist{}'.format(verifier_list))
            clients = get_clients_by_nodeid(verifier_list, all_clients)
            return [client.pip for client in clients]
        else:
            pip.economic.env.deploy_all()
    result = pip.submitVersion(pip.node.node_id, str(time.time()), pip.cfg.version8, 10,
                               pip.node.staking_address, transaction_cfg=pip.cfg.transaction_cfg)
    log.info('version proposal result :{}'.format(result))
    verifier_list = get_pledge_list(all_clients[0].ppos.getVerifierList)
    log.info('verifierlist{}'.format(verifier_list))
    clients = get_clients_by_nodeid(verifier_list, all_clients)
    return [client.pip for client in clients]


@pytest.fixture()
def proposal_candidate_pips(all_clients):
    '''
    There is voting stage proposal, get candidate list pip object
    :param global_test_env:
    :return:
    '''
    pip = all_clients[0].pip
    if pip.chain_version != pip.cfg.version0 or (pip.is_exist_effective_proposal()
                                                 and not pip.is_exist_effective_proposal_for_vote()):
        log.info('The chain has been upgraded or there is preactive proposal,restart!')
        pip.economic.env.deploy_all()
    nodeid_list = pip.get_candidate_list_not_verifier()
    if nodeid_list:
        if pip.get_effect_proposal_info_of_vote():
            proposalinfo = pip.get_effect_proposal_info_of_vote()
            log.info('get version proposalinfo : {}'.format(proposalinfo))
            if proposalinfo.get('NewVersion') == pip.cfg.version8:
                pip.economic.env.deploy_all()
            else:
                if proposalinfo.get('EndVotingBlock') - pip.node.block_number > pip.economic.consensus_size:
                    client_candidates = get_clients_by_nodeid(nodeid_list, all_clients)
                    return [client.pip for client in client_candidates]

    candidate_list = get_pledge_list(all_clients[0].ppos.getCandidateList)
    log.info('candidate_list{}'.format(candidate_list))
    for client in all_clients:
        if client.node.node_id not in candidate_list:
            address, _ = client.economic.account.generate_account(client.node.web3, 10**18 * 10000000)
            result = client.staking.create_staking(0, address, address)
            log.info('node {} staking result {}'.format(client.node.node_id, result))
    verifier_list = get_pledge_list(all_clients[0].ppos.getVerifierList)
    log.info('Verifier list {}'.format(verifier_list))
    verifier_pip = get_client_by_nodeid(verifier_list[0], all_clients).pip
    result = verifier_pip.submitVersion(verifier_pip.node.node_id, str(time.time()), verifier_pip.cfg.version5,
                                        10, verifier_pip.node.staking_address,
                                        transaction_cfg=verifier_pip.cfg.transaction_cfg)
    log.info('Submit version proposal result {}'.format(result))
    assert_code(result, 0)
    nodeid_list = all_clients[0].pip.get_candidate_list_not_verifier()
    if not nodeid_list:
        raise Exception('get candidate not verifier failed')
    client_candiates = get_clients_by_nodeid(nodeid_list, all_clients)
    return [client_candiate.pip for client_candiate in client_candiates]


@pytest.fixture()
def large_version_proposal_candidate_pips(all_clients):
    '''
    There is voting stage proposal, get candidate list pip object
    :param global_test_env:
    :return:
    '''
    pip = all_clients[0].pip
    if pip.chain_version != pip.cfg.version0 or (pip.is_exist_effective_proposal()
                                                 and not pip.is_exist_effective_proposal_for_vote()):
        log.info('The chain has been upgraded or there is preactive proposal,restart!')
        pip.economic.env.deploy_all()
    nodeid_list = pip.get_candidate_list_not_verifier()
    if nodeid_list:
        if pip.get_effect_proposal_info_of_vote():
            proposalinfo = pip.get_effect_proposal_info_of_vote()
            log.info('get version proposalinfo : {}'.format(proposalinfo))
            if proposalinfo.get('NewVersion') == pip.cfg.version5:
                pip.economic.env.deploy_all()
            else:
                if proposalinfo.get('EndVotingBlock') - pip.node.block_number > pip.economic.consensus_size:
                    client_candiates = get_clients_by_nodeid(nodeid_list, all_clients)
                    return [client.pip for client in client_candiates]

    candidate_list = get_pledge_list(all_clients[0].ppos.getCandidateList)
    log.info('candidate_list{}'.format(candidate_list))
    for client in all_clients:
        if client.node.node_id not in candidate_list:
            address, _ = client.economic.account.generate_account(client.node.web3, 10**18 * 10000000)
            result = client.staking.create_staking(0, address, address)
            log.info('node {} staking result {}'.format(client.node.node_id, result))
    verifier_list = get_pledge_list(all_clients[0].ppos.getVerifierList)
    log.info('Verifier list {}'.format(verifier_list))
    verifier_pip = get_client_by_nodeid(verifier_list[0], all_clients).pip
    result = verifier_pip.submitVersion(verifier_pip.node.node_id, str(time.time()),
                                        verifier_pip.cfg.version8,
                                        10, verifier_pip.node.staking_address,
                                        transaction_cfg=verifier_pip.cfg.transaction_cfg)
    log.info('Submit version proposal result {}'.format(result))
    assert_code(result, 0)
    nodeid_list = all_clients[0].pip.get_candidate_list_not_verifier()
    if not nodeid_list:
        raise Exception('get candidate not verifier failed')
    client_candidates = get_clients_by_nodeid(nodeid_list, all_clients)
    return [client_candidate.pip for client_candidate in client_candidates]


@pytest.fixture()
def proposal_voted_pips(all_clients):
    '''
    get verifier Client object list
    :param global_test_env:
    :return:
    '''
    verifier_list = get_pledge_list(all_clients[0].ppos.getVerifierList)
    log.info('verifier list {}'.format(verifier_list))
    pip = get_client_by_nodeid(verifier_list[0], all_clients).pip
    pip.economic.env.deploy_all()
    result = pip.submitVersion(pip.node.node_id, str(time.time()), pip.cfg.version5, 10,
                               pip.node.staking_address, transaction_cfg=pip.cfg.transaction_cfg)
    log.info('version proposal result :{}'.format(result))
    assert_code(result, 0)
    verifier_list = get_pledge_list(all_clients[0].ppos.getVerifierList)
    log.info('verifierlist{}'.format(verifier_list))
    clients_verifier = get_clients_by_nodeid(verifier_list, all_clients)
    result = version_proposal_vote(clients_verifier[0].pip)
    assert_code(result, 0)
    return [client.pip for client in clients_verifier]


@pytest.fixture()
def large_version_proposal_voted_pips(all_clients):
    '''
    get verifier Client object list
    :param global_test_env:
    :return:
    '''
    verifier_list = get_pledge_list(all_clients[0].ppos.getVerifierList)
    log.info('verifier list {}'.format(verifier_list))
    pip = get_client_by_nodeid(verifier_list[0], all_clients).pip
    pip.economic.env.deploy_all()
    result = pip.submitVersion(pip.node.node_id, str(time.time()), pip.cfg.version8, 10,
                               pip.node.staking_address, transaction_cfg=pip.cfg.transaction_cfg)
    log.info('version proposal result :{}'.format(result))
    verifier_list = get_pledge_list(all_clients[0].ppos.getVerifierList)
    log.info('verifierlist{}'.format(verifier_list))
    clients_verifier = get_clients_by_nodeid(verifier_list, all_clients)
    result = version_proposal_vote(clients_verifier[0].pip)
    assert_code(result, 0)
    return [client.pip for client in clients_verifier]


def replace_version_declare(pip, platon_bin, versiontag):
    upload_platon(pip.node, platon_bin)
    log.info('Replace the platon of the node {} version{}'.format(pip.node.node_id, versiontag))
    pip.node.restart()
    log.info('Restart the node{}'.format(pip.node.node_id))
    assert pip.node.program_version == versiontag
    log.info('assert the version of the node is {}'.format(versiontag))
    log.info("staking: {}".format(pip.node.staking_address))
    log.info("account:{}".format(pip.economic.account.accounts))
    result = pip.declareVersion(pip.node.node_id, pip.node.staking_address,
                                transaction_cfg=pip.cfg.transaction_cfg)
    log.info('declareversion {} result: {}'.format(pip.node.program_version, result))
    return result


def wrong_verisonsign_declare(pip, pip_test):
    result = pip.declareVersion(pip.node.node_id, pip.node.staking_address,
                                version_sign=pip_test.node.program_version_sign,
                                transaction_cfg=pip.cfg.transaction_cfg)
    log.info('wrong program version sign, declareVersion result : {}'.format(result))
    return result


def wrong_verison_declare(pip, version=None):
    if not version:
        proposalinfo = pip.get_effect_proposal_info_of_vote()
        version = proposalinfo.get('NewVersion')
        log.info('The new version of the proposal: {}'.format(version))
    result = pip.declareVersion(pip.node.node_id, pip.node.staking_address,
                                program_version=version,
                                transaction_cfg=pip.cfg.transaction_cfg)
    log.info('wrong program version, declareVersion: {} result : {}'.format(version, result))
    return result


@pytest.mark.P0
@allure.title('Not staking address declare version')
def test_DE_DE_001(client_verifier):
    pip = client_verifier.pip
    address, _ = pip.economic.account.generate_account(pip.node.web3, 10**18 * 10000)
    result = pip.declareVersion(pip.node.node_id, address, transaction_cfg=pip.cfg.transaction_cfg)
    log.info('declareVersion result: {}'.format(result))
    assert_code(result, 302021)


class TestNoProposalVE:
    @pytest.mark.P0
    @pytest.mark.compatibility
    @allure.title('No effective proposal, verifier declare version')
    def test_DE_VE_001(self, noproposal_pips):
        pip = noproposal_pips[0]
        verison = struct.pack('>I', pip.chain_version)
        if verison[3] != 0:
            result = replace_version_declare(pip, pip.cfg.PLATON_NEW_BIN2, pip.cfg.version2)
            assert_code(result, 0)
            verifier_node_version(pip, pip.cfg.version2)

            result = wrong_verisonsign_declare(pip, noproposal_pips[1])
            assert_code(result, 302024)

            result = wrong_verison_declare(pip, pip.chain_version)
            assert_code(result, 302024)

    @pytest.mark.P3
    @allure.title('No effective proposal, verifier declare version')
    def test_DE_VE_002(self, noproposal_pips, all_clients):
        pip = noproposal_pips[0]
        node_version = verifier_node_version(pip)
        result = replace_version_declare(pip, pip.cfg.PLATON_NEW_BIN1, pip.cfg.version1)
        assert_code(result, 302028)
        verifier_node_version(pip, node_version)

        result = wrong_verisonsign_declare(pip, noproposal_pips[1])
        assert_code(result, 302024)

        result = wrong_verison_declare(pip, pip.chain_version)
        assert_code(result, 302024)

    @pytest.mark.P0
    @allure.title('No effective proposal, verifier declare version')
    def test_DE_VE_004(self, noproposal_pips):
        pip = noproposal_pips[0]
        result = replace_version_declare(pip, pip.cfg.PLATON_NEW_BIN3, pip.cfg.version3)
        assert_code(result, 0)
        verifier_node_version(pip, pip.cfg.version3)

        result = wrong_verisonsign_declare(pip, noproposal_pips[1])
        assert_code(result, 302024)

        result = wrong_verison_declare(pip, pip.chain_version)
        assert_code(result, 302024)

    @pytest.mark.P2
    @allure.title('No effective proposal, verifier declare version')
    def test_DE_VE_005(self, noproposal_pips):
        pip = noproposal_pips[0]
        result = replace_version_declare(pip, pip.cfg.PLATON_NEW_BIN, pip.cfg.version5)
        assert_code(result, 302028)

        result = wrong_verisonsign_declare(pip, noproposal_pips[1])
        assert_code(result, 302024)

        result = wrong_verison_declare(pip, pip.chain_version)
        assert_code(result, 302024)

    @pytest.mark.P2
    @allure.title('No effective proposal, verifier declare version')
    def test_DE_VE_006(self, noproposal_pips):
        pip = noproposal_pips[0]
        result = replace_version_declare(pip, pip.cfg.PLATON_NEW_BIN8, pip.cfg.version8)
        assert_code(result, 302028)

        result = wrong_verisonsign_declare(pip, noproposal_pips[1])
        assert_code(result, 302024)

        result = wrong_verison_declare(pip, pip.chain_version)
        assert_code(result, 302024)

    @pytest.mark.P0
    @allure.title('No effective proposal, verifier declare version')
    def test_DE_VE_007(self, noproposal_pips):
        pip = noproposal_pips[0]
        verison = struct.pack('>I', pip.chain_version)
        if verison[3] != 0:
            result = replace_version_declare(pip, pip.cfg.PLATON_NEW_BIN0, pip.cfg.version0)
            assert_code(result, 0)
            verifier_node_version(pip, pip.cfg.version0)

            result = wrong_verisonsign_declare(pip, noproposal_pips[1])
            assert_code(result, 302024)

            result = wrong_verison_declare(pip, pip.cfg.version3)
            assert_code(result, 302024)

            result = wrong_verison_declare(pip, pip.cfg.version2)
            assert_code(result, 302024)


class TestVotingProposalVE:
    @pytest.mark.P0
    @allure.title('No effective proposal, verifier declare version')
    def test_DE_VE_008(self, proposal_pips):
        pip = proposal_pips[0]
        verison = struct.pack('>I', pip.chain_version)
        if verison[3] != 0:
            result = replace_version_declare(pip, pip.cfg.PLATON_NEW_BIN2, pip.cfg.version2)
            assert_code(result, 0)
            verifier_node_version(pip, pip.cfg.version2)

            result = wrong_verisonsign_declare(pip, proposal_pips[1])
            assert_code(result, 302024)

            result = wrong_verison_declare(pip, pip.chain_version)
            assert_code(result, 302024)

            result = wrong_verison_declare(pip)
            assert_code(result, 302024)

    @pytest.mark.P2
    @allure.title('No effective proposal, verifier declare version')
    def test_DE_VE_010(self, proposal_pips):
        pip = proposal_pips[0]
        result = replace_version_declare(pip, pip.cfg.PLATON_NEW_BIN1, pip.cfg.version1)
        assert_code(result, 302028)

        result = wrong_verisonsign_declare(pip, proposal_pips[1])
        assert_code(result, 302024)

        result = wrong_verison_declare(pip, pip.chain_version)
        assert_code(result, 302024)

        result = wrong_verison_declare(pip)
        assert_code(result, 302024)

    @pytest.mark.P1
    @allure.title('No effective proposal, verifier declare version')
    def test_DE_VE_014(self, large_version_proposal_pips):
        pip = large_version_proposal_pips[0]
        verison = struct.pack('>I', pip.chain_version)
        if verison[3] != 0:
            result = replace_version_declare(pip, pip.cfg.PLATON_NEW_BIN2, pip.cfg.version2)
            assert_code(result, 0)
            verifier_node_version(pip, pip.cfg.version2)

            result = wrong_verisonsign_declare(pip, large_version_proposal_pips[1])
            assert_code(result, 302024)

            result = wrong_verison_declare(pip, pip.chain_version)
            assert_code(result, 302024)

            result = wrong_verison_declare(pip)
            assert_code(result, 302024)

    @pytest.mark.P2
    @allure.title('No effective proposal, verifier declare version')
    def test_DE_VE_025(self, large_version_proposal_pips):
        pip = large_version_proposal_pips[0]
        result = replace_version_declare(pip, pip.cfg.PLATON_NEW_BIN1, pip.cfg.version1)
        assert_code(result, 302028)

        result = wrong_verisonsign_declare(pip, large_version_proposal_pips[1])
        assert_code(result, 302024)

        result = wrong_verison_declare(pip, pip.chain_version)
        assert_code(result, 302024)

        result = wrong_verison_declare(pip)
        assert_code(result, 302024)

    @pytest.mark.P0
    @allure.title('No effective proposal, verifier declare version')
    def test_DE_VE_032(self, proposal_pips):
        pip = proposal_pips[0]
        result = replace_version_declare(pip, pip.cfg.PLATON_NEW_BIN0, pip.cfg.version0)
        assert_code(result, 0)
        verifier_node_version(pip, pip.cfg.version0)

        result = wrong_verisonsign_declare(pip, proposal_pips[1])
        assert_code(result, 302024)

        result = wrong_verison_declare(pip)
        assert_code(result, 302024)

    @pytest.mark.P1
    @allure.title('No effective proposal, verifier declare version')
    def test_DE_VE_034(self, large_version_proposal_pips):
        pip = large_version_proposal_pips[0]
        result = replace_version_declare(pip, pip.cfg.PLATON_NEW_BIN0, pip.cfg.version0)
        assert_code(result, 0)
        verifier_node_version(pip, pip.cfg.version0)

        result = wrong_verisonsign_declare(pip, large_version_proposal_pips[1])
        assert_code(result, 302024)

        result = wrong_verison_declare(pip)
        assert_code(result, 302024)

    @pytest.mark.P1
    @allure.title('No effective proposal, verifier declare version')
    def test_DE_VE_036(self, proposal_pips):
        pip = proposal_pips[0]
        result = replace_version_declare(pip, pip.cfg.PLATON_NEW_BIN3, pip.cfg.version3)
        assert_code(result, 0)
        verifier_node_version(pip, pip.cfg.version3)

        result = wrong_verisonsign_declare(pip, proposal_pips[1])
        assert_code(result, 302024)

        result = wrong_verison_declare(pip)
        assert_code(result, 302024)

        result = wrong_verison_declare(pip, pip.chain_version)
        assert_code(result, 302024)

    @pytest.mark.P1
    @allure.title('No effective proposal, verifier declare version')
    def test_DE_VE_038(self, proposal_pips, all_clients):
        pip = proposal_pips[0]
        node_version = verifier_node_version(pip)
        result = replace_version_declare(pip, pip.cfg.PLATON_NEW_BIN3, pip.cfg.version3)
        assert_code(result, 0)
        verifier_node_version(pip, node_version)

        result = wrong_verisonsign_declare(pip, proposal_pips[1])
        assert_code(result, 302024)

        result = wrong_verison_declare(pip, pip.chain_version)
        assert_code(result, 302024)

    @pytest.mark.P2
    @allure.title('No effective proposal, verifier declare version')
    def test_DE_VE_040(self, proposal_pips):
        pip = proposal_pips[0]
        node_version = verifier_node_version(pip)
        result = replace_version_declare(pip, pip.cfg.PLATON_NEW_BIN4, pip.cfg.version4)
        assert_code(result, 0)
        verifier_node_version(pip, node_version)

        result = wrong_verisonsign_declare(pip, proposal_pips[1])
        assert_code(result, 302024)

        result = wrong_verison_declare(pip)
        assert_code(result, 302024)

        result = wrong_verison_declare(pip, pip.chain_version)
        assert_code(result, 302024)

    @pytest.mark.P2
    @allure.title('No effective proposal, verifier declare version')
    def test_DE_VE_042(self, proposal_pips):
        pip = proposal_pips[0]
        node_version = verifier_node_version(pip)
        result = replace_version_declare(pip, pip.cfg.PLATON_NEW_BIN6, pip.cfg.version6)
        assert_code(result, 0)
        verifier_node_version(pip, node_version)

        result = wrong_verisonsign_declare(pip, proposal_pips[1])
        assert_code(result, 302024)

        result = wrong_verison_declare(pip)
        assert_code(result, 302024)

        result = wrong_verison_declare(pip, pip.chain_version)
        assert_code(result, 302024)

    @pytest.mark.P2
    @allure.title('No effective proposal, verifier declare version')
    def test_DE_VE_044(self, proposal_pips):
        pip = proposal_pips[0]
        result = replace_version_declare(pip, pip.cfg.PLATON_NEW_BIN8, pip.cfg.version8)
        assert_code(result, 302028)

        result = wrong_verisonsign_declare(pip, proposal_pips[1])
        assert_code(result, 302024)

        result = wrong_verison_declare(pip)
        assert_code(result, 302024)

        result = wrong_verison_declare(pip, pip.chain_version)
        assert_code(result, 302024)

    @pytest.mark.P1
    @allure.title('No effective proposal, verifier declare version')
    def test_DE_VE_046(self, large_version_proposal_pips):
        pip = large_version_proposal_pips[0]
        result = replace_version_declare(pip, pip.cfg.PLATON_NEW_BIN3, pip.cfg.version3)
        assert_code(result, 0)
        verifier_node_version(pip, pip.cfg.version3)

        result = wrong_verisonsign_declare(pip, large_version_proposal_pips[1])
        assert_code(result, 302024)

        result = wrong_verison_declare(pip)
        assert_code(result, 302024)

        result = wrong_verison_declare(pip, pip.chain_version)
        assert_code(result, 302024)

    @pytest.mark.P1
    @allure.title('No effective proposal, verifier declare version')
    def test_DE_VE_048(self, large_version_proposal_pips):
        pip = large_version_proposal_pips[0]
        result = replace_version_declare(pip, pip.cfg.PLATON_NEW_BIN, pip.cfg.version5)
        assert_code(result, 302028)

        result = wrong_verisonsign_declare(pip, large_version_proposal_pips[1])
        assert_code(result, 302024)

        result = wrong_verison_declare(pip)
        assert_code(result, 302024)

        result = wrong_verison_declare(pip, pip.chain_version)
        assert_code(result, 302024)

    @pytest.mark.P1
    @allure.title('No effective proposal, verifier declare version')
    def test_DE_VE_050(self, large_version_proposal_pips):
        pip = large_version_proposal_pips[0]
        result = replace_version_declare(pip, pip.cfg.PLATON_NEW_BIN4, pip.cfg.version4)
        assert_code(result, 302028)

        result = wrong_verisonsign_declare(pip, large_version_proposal_pips[1])
        assert_code(result, 302024)

        result = wrong_verison_declare(pip)
        assert_code(result, 302024)

        result = wrong_verison_declare(pip, pip.chain_version)
        assert_code(result, 302024)

    @pytest.mark.P1
    @allure.title('No effective proposal, verifier declare version')
    def test_DE_VE_052(self, large_version_proposal_pips):
        pip = large_version_proposal_pips[0]
        result = replace_version_declare(pip, pip.cfg.PLATON_NEW_BIN6, pip.cfg.version6)
        assert_code(result, 302028)

        result = wrong_verisonsign_declare(pip, large_version_proposal_pips[1])
        assert_code(result, 302024)

        result = wrong_verison_declare(pip)
        assert_code(result, 302024)

        result = wrong_verison_declare(pip, pip.chain_version)
        assert_code(result, 302024)

    @pytest.mark.P2
    @allure.title('No effective proposal, verifier declare version')
    def test_DE_VE_054(self, large_version_proposal_pips):
        pip = large_version_proposal_pips[0]
        node_version = verifier_node_version(pip)
        result = replace_version_declare(pip, pip.cfg.PLATON_NEW_BIN8, pip.cfg.version8)
        assert_code(result, 0)
        verifier_node_version(pip, node_version)

        result = wrong_verisonsign_declare(pip, large_version_proposal_pips[1])
        assert_code(result, 302024)

        result = wrong_verison_declare(pip, pip.chain_version)
        assert_code(result, 302024)


class TestVotingProposlaVotedVE:
    @pytest.mark.P2
    @allure.title('There is a proposal in voting stage, verifier declare version')
    def test_DE_VE_009(self, proposal_voted_pips):
        pip = proposal_voted_pips[0]
        verison = struct.pack('>I', pip.chain_version)
        if verison[3] != 0:
            result = replace_version_declare(pip, pip.cfg.PLATON_NEW_BIN2, pip.cfg.version2)
            assert_code(result, 302028)

            result = wrong_verisonsign_declare(pip, proposal_voted_pips[1])
            assert_code(result, 302024)

            result = wrong_verison_declare(pip)
            assert_code(result, 302024)

            result = wrong_verison_declare(pip, pip.chain_version)
            assert_code(result, 302024)

    @pytest.mark.P2
    @allure.title('There is a proposal in voting stage, verifier declare version')
    def test_DE_VE_011(self, proposal_voted_pips):
        pip = proposal_voted_pips[0]
        result = replace_version_declare(pip, pip.cfg.PLATON_NEW_BIN1, pip.cfg.version1)
        assert_code(result, 302028)

        result = wrong_verisonsign_declare(pip, proposal_voted_pips[1])
        assert_code(result, 302024)

        result = wrong_verison_declare(pip)
        assert_code(result, 302024)

        result = wrong_verison_declare(pip, pip.chain_version)
        assert_code(result, 302024)

    @pytest.mark.P2
    @allure.title('There is a proposal in voting stage, verifier declare version')
    def test_DE_VE_021(self, large_version_proposal_voted_pips):
        pip = large_version_proposal_voted_pips[0]
        verison = struct.pack('>I', pip.chain_version)
        if verison[3] != 0:
            result = replace_version_declare(pip, pip.cfg.PLATON_NEW_BIN2, pip.cfg.version2)
            assert_code(result, 302028)

            result = wrong_verisonsign_declare(pip, large_version_proposal_voted_pips[1])
            assert_code(result, 302024)

            result = wrong_verison_declare(pip)
            assert_code(result, 302024)

            result = wrong_verison_declare(pip, pip.chain_version)
            assert_code(result, 302024)

    @pytest.mark.P2
    @allure.title('There is a proposal in voting stage, verifier declare version')
    def test_DE_VE_026(self, large_version_proposal_voted_pips):
        pip = large_version_proposal_voted_pips[0]
        result = replace_version_declare(pip, pip.cfg.PLATON_NEW_BIN1, pip.cfg.version1)
        assert_code(result, 302028)

        result = wrong_verisonsign_declare(pip, large_version_proposal_voted_pips[1])
        assert_code(result, 302024)

        result = wrong_verison_declare(pip)
        assert_code(result, 302024)

        result = wrong_verison_declare(pip, pip.chain_version)
        assert_code(result, 302024)

    @pytest.mark.P2
    @allure.title('There is a proposal in voting stage, verifier declare version')
    def test_DE_VE_033(self, proposal_voted_pips):
        pip = proposal_voted_pips[0]
        result = replace_version_declare(pip, pip.cfg.PLATON_NEW_BIN0, pip.cfg.version0)
        assert_code(result, 302028)

        result = wrong_verisonsign_declare(pip, proposal_voted_pips[1])
        assert_code(result, 302024)

        result = wrong_verison_declare(pip)
        assert_code(result, 302024)

    @pytest.mark.P2
    @allure.title('There is a proposal in voting stage, verifier declare version')
    def test_DE_VE_035(self, large_version_proposal_voted_pips):
        pip = large_version_proposal_voted_pips[0]
        result = replace_version_declare(pip, pip.cfg.PLATON_NEW_BIN0, pip.cfg.version0)
        assert_code(result, 302028)

        result = wrong_verisonsign_declare(pip, large_version_proposal_voted_pips[1])
        assert_code(result, 302024)

        result = wrong_verison_declare(pip)
        assert_code(result, 302024)

    @pytest.mark.P2
    @allure.title('There is a proposal in voting stage, verifier declare version')
    def test_DE_VE_037(self, proposal_voted_pips):
        pip = proposal_voted_pips[0]
        result = replace_version_declare(pip, pip.cfg.PLATON_NEW_BIN3, pip.cfg.version3)
        assert_code(result, 302028)

        result = wrong_verisonsign_declare(pip, proposal_voted_pips[1])
        assert_code(result, 302024)

        result = wrong_verison_declare(pip)
        assert_code(result, 302024)

        result = wrong_verison_declare(pip, pip.chain_version)
        assert_code(result, 302024)

    @pytest.mark.P1
    @allure.title('There is a proposal in voting stage, verifier declare version')
    def test_DE_VE_039(self, proposal_voted_pips):
        pip = proposal_voted_pips[0]
        node_version = verifier_node_version(pip)
        result = replace_version_declare(pip, pip.cfg.PLATON_NEW_BIN, pip.cfg.version5)
        assert_code(result, 0)
        verifier_node_version(pip, node_version)

        result = wrong_verisonsign_declare(pip, proposal_voted_pips[1])
        assert_code(result, 302024)

        result = wrong_verison_declare(pip, pip.cfg.version4)
        assert_code(result, 302024)

        result = wrong_verison_declare(pip, pip.chain_version)
        assert_code(result, 302024)

    @pytest.mark.P2
    @allure.title('There is a proposal in voting stage, verifier declare version')
    def test_DE_VE_041(self, proposal_voted_pips):
        pip = proposal_voted_pips[0]
        node_version = verifier_node_version(pip)
        result = replace_version_declare(pip, pip.cfg.PLATON_NEW_BIN4, pip.cfg.version4)
        assert_code(result, 0)
        verifier_node_version(pip, node_version)

        result = wrong_verisonsign_declare(pip, proposal_voted_pips[1])
        assert_code(result, 302024)

        result = wrong_verison_declare(pip, pip.cfg.version5)
        assert_code(result, 302024)

        result = wrong_verison_declare(pip, pip.chain_version)
        assert_code(result, 302024)

    @pytest.mark.P2
    @allure.title('There is a proposal in voting stage, verifier declare version')
    def test_DE_VE_043(self, proposal_voted_pips):
        pip = proposal_voted_pips[0]
        node_version = verifier_node_version(pip)
        result = replace_version_declare(pip, pip.cfg.PLATON_NEW_BIN6, pip.cfg.version6)
        assert_code(result, 0)
        verifier_node_version(pip, node_version)

        result = wrong_verisonsign_declare(pip, proposal_voted_pips[1])
        assert_code(result, 302024)

        result = wrong_verison_declare(pip, pip.cfg.version5)
        assert_code(result, 302024)

        result = wrong_verison_declare(pip, pip.chain_version)
        assert_code(result, 302024)

    @pytest.mark.P2
    @allure.title('There is a proposal in voting stage, verifier declare version')
    def test_DE_VE_045(self, proposal_voted_pips):
        pip = proposal_voted_pips[0]
        result = replace_version_declare(pip, pip.cfg.PLATON_NEW_BIN8, pip.cfg.version8)
        assert_code(result, 302028)

        result = wrong_verisonsign_declare(pip, proposal_voted_pips[1])
        assert_code(result, 302024)

        result = wrong_verison_declare(pip, pip.cfg.version5)
        assert_code(result, 302024)

        result = wrong_verison_declare(pip, pip.chain_version)
        assert_code(result, 302024)

    @pytest.mark.P2
    @allure.title('There is a proposal in voting stage, verifier declare version')
    def test_DE_VE_047(self, large_version_proposal_voted_pips):
        pip = large_version_proposal_voted_pips[0]
        result = replace_version_declare(pip, pip.cfg.PLATON_NEW_BIN3, pip.cfg.version3)
        assert_code(result, 302028)

        result = wrong_verisonsign_declare(pip, large_version_proposal_voted_pips[1])
        assert_code(result, 302024)

        result = wrong_verison_declare(pip, pip.cfg.version5)
        assert_code(result, 302024)

        result = wrong_verison_declare(pip, pip.chain_version)
        assert_code(result, 302024)

    @pytest.mark.P1
    @allure.title('There is a proposal in voting stage, verifier declare version')
    def test_DE_VE_049(self, large_version_proposal_voted_pips):
        pip = large_version_proposal_voted_pips[0]
        result = replace_version_declare(pip, pip.cfg.PLATON_NEW_BIN, pip.cfg.version5)
        assert_code(result, 302028)

        result = wrong_verisonsign_declare(pip, large_version_proposal_voted_pips[1])
        assert_code(result, 302024)

        result = wrong_verison_declare(pip, pip.cfg.version8)
        assert_code(result, 302024)

        result = wrong_verison_declare(pip, pip.chain_version)
        assert_code(result, 302024)

    @pytest.mark.P1
    @allure.title('There is a proposal in voting stage, verifier declare version')
    def test_DE_VE_051(self, large_version_proposal_voted_pips):
        pip = large_version_proposal_voted_pips[0]
        result = replace_version_declare(pip, pip.cfg.PLATON_NEW_BIN4, pip.cfg.version4)
        assert_code(result, 302028)

        result = wrong_verisonsign_declare(pip, large_version_proposal_voted_pips[1])
        assert_code(result, 302024)

        result = wrong_verison_declare(pip, pip.cfg.version8)
        assert_code(result, 302024)

        result = wrong_verison_declare(pip, pip.chain_version)
        assert_code(result, 302024)

    @pytest.mark.P1
    @allure.title('There is a proposal in voting stage, verifier declare version')
    def test_DE_VE_053(self, large_version_proposal_voted_pips):
        pip = large_version_proposal_voted_pips[0]
        result = replace_version_declare(pip, pip.cfg.PLATON_NEW_BIN6, pip.cfg.version6)
        assert_code(result, 302028)

        result = wrong_verisonsign_declare(pip, large_version_proposal_voted_pips[1])
        assert_code(result, 302024)

        result = wrong_verison_declare(pip, pip.cfg.version8)
        assert_code(result, 302024)

        result = wrong_verison_declare(pip, pip.chain_version)
        assert_code(result, 302024)

    @pytest.mark.P2
    @allure.title('There is a proposal in voting stage, verifier declare version')
    def test_DE_VE_055(self, large_version_proposal_voted_pips):
        pip = large_version_proposal_voted_pips[0]
        node_version = verifier_node_version(pip)
        result = replace_version_declare(pip, pip.cfg.PLATON_NEW_BIN8, pip.cfg.version8)
        assert_code(result, 0)
        verifier_node_version(pip, node_version)

        result = wrong_verisonsign_declare(pip, large_version_proposal_voted_pips[1])
        assert_code(result, 302024)

        result = wrong_verison_declare(pip, pip.chain_version)
        assert_code(result, 302024)


class TestPreactiveProposalVE:
    @pytest.mark.P2
    @allure.title('There is a preactive proposal, verifier declare version')
    def test_DE_VE_056(self, preactive_proposal_pips):
        pip = preactive_proposal_pips[0]
        result = replace_version_declare(pip, pip.cfg.PLATON_NEW_BIN2, pip.cfg.version2)
        assert_code(result, 302028)

        result = wrong_verisonsign_declare(pip, preactive_proposal_pips[1])
        assert_code(result, 302024)

        result = wrong_verison_declare(pip, pip.cfg.version5)
        assert_code(result, 302024)

        result = wrong_verison_declare(pip, pip.chain_version)
        assert_code(result, 302028)

    @pytest.mark.P2
    @allure.title('There is a preactive proposal, verifier declare version')
    def test_DE_VE_057(self, preactive_proposal_pips):
        pip = preactive_proposal_pips[0]
        result = replace_version_declare(pip, pip.cfg.PLATON_NEW_BIN1, pip.cfg.version1)
        assert_code(result, 302028)

        result = wrong_verisonsign_declare(pip, preactive_proposal_pips[1])
        assert_code(result, 302024)

        result = wrong_verison_declare(pip, pip.cfg.version5)
        assert_code(result, 302024)

        result = wrong_verison_declare(pip, pip.chain_version)
        assert_code(result, 302024)

    @pytest.mark.P2
    @allure.title('There is a preactive proposal, verifier declare version')
    def test_DE_VE_059(self, preactive_proposal_pips):
        pip = preactive_proposal_pips[0]
        result = replace_version_declare(pip, pip.cfg.PLATON_NEW_BIN2, pip.cfg.version2)
        assert_code(result, 302028)

        result = wrong_verisonsign_declare(pip, preactive_proposal_pips[1])
        assert_code(result, 302024)

        result = wrong_verison_declare(pip, pip.cfg.version5)
        assert_code(result, 302024)

        result = wrong_verison_declare(pip, pip.chain_version)
        assert_code(result, 302028)

    @pytest.mark.P2
    @allure.title('There is a preactive proposal, verifier declare version')
    def test_DE_VE_060(self, preactive_large_version_proposal_pips):
        pip = preactive_large_version_proposal_pips[0]
        result = replace_version_declare(pip, pip.cfg.PLATON_NEW_BIN1, pip.cfg.version1)
        assert_code(result, 302028)

        result = wrong_verisonsign_declare(pip, preactive_large_version_proposal_pips[1])
        assert_code(result, 302024)

        result = wrong_verison_declare(pip, pip.cfg.version8)
        assert_code(result, 302024)

        result = wrong_verison_declare(pip, pip.chain_version)
        assert_code(result, 302024)

    @pytest.mark.P1
    @allure.title('There is a preactive proposal, verifier declare version')
    def test_DE_VE_062(self, preactive_proposal_pips):
        pip = preactive_proposal_pips[0]
        result = replace_version_declare(pip, pip.cfg.PLATON_NEW_BIN0, pip.cfg.version0)
        assert_code(result, 302028)

        result = wrong_verisonsign_declare(pip, preactive_proposal_pips[1])
        assert_code(result, 302024)

        result = wrong_verison_declare(pip, pip.cfg.version5)
        assert_code(result, 302024)

    @pytest.mark.P2
    @allure.title('There is a preactive proposal, verifier declare version')
    def test_DE_VE_063(self, preactive_large_version_proposal_pips):
        pip = preactive_large_version_proposal_pips[0]
        result = replace_version_declare(pip, pip.cfg.PLATON_NEW_BIN0, pip.cfg.version0)
        assert_code(result, 302028)

        result = wrong_verisonsign_declare(pip, preactive_large_version_proposal_pips[1])
        assert_code(result, 302024)

        result = wrong_verison_declare(pip, pip.cfg.version8)
        assert_code(result, 302024)

    @pytest.mark.P2
    @allure.title('There is a preactive proposal, verifier declare version')
    def test_DE_VE_064(self, preactive_proposal_pips):
        pip = preactive_proposal_pips[0]
        result = replace_version_declare(pip, pip.cfg.PLATON_NEW_BIN3, pip.cfg.version3)
        assert_code(result, 302028)

        result = wrong_verisonsign_declare(pip, preactive_proposal_pips[1])
        assert_code(result, 302024)

        result = wrong_verison_declare(pip, pip.cfg.version5)
        assert_code(result, 302024)

        result = wrong_verison_declare(pip, pip.chain_version)
        assert_code(result, 302024)

    @pytest.mark.P2
    @allure.title('There is a preactive proposal, verifier declare version')
    def test_DE_VE_065(self, preactive_proposal_pips):
        pip = preactive_proposal_pips[0]
        node_version = verifier_node_version(pip)
        result = replace_version_declare(pip, pip.cfg.PLATON_NEW_BIN, pip.cfg.version5)
        assert_code(result, 0)
        verifier_node_version(pip, node_version)

        result = wrong_verisonsign_declare(pip, preactive_proposal_pips[1])
        assert_code(result, 302024)

        result = wrong_verison_declare(pip, pip.cfg.version0)
        assert_code(result, 302024)

    @pytest.mark.P1
    @allure.title('There is a preactive proposal, verifier declare version')
    def test_DE_VE_066(self, preactive_proposal_pips):
        pip = preactive_proposal_pips[0]
        result = replace_version_declare(pip, pip.cfg.PLATON_NEW_BIN4, pip.cfg.version4)
        assert_code(result, 0)
        verifier_node_version(pip, pip.cfg.version4)

        result = wrong_verisonsign_declare(pip, preactive_proposal_pips[1])
        assert_code(result, 302024)

        result = wrong_verison_declare(pip, pip.cfg.version5)
        assert_code(result, 302024)

        result = wrong_verison_declare(pip, pip.chain_version)
        assert_code(result, 302024)

    @pytest.mark.P1
    @allure.title('There is a preactive proposal, verifier declare version')
    def test_DE_VE_067(self, preactive_proposal_pips):
        pip = preactive_proposal_pips[0]
        result = replace_version_declare(pip, pip.cfg.PLATON_NEW_BIN6, pip.cfg.version6)
        assert_code(result, 0)
        verifier_node_version(pip, pip.cfg.version6)

        result = wrong_verisonsign_declare(pip, preactive_proposal_pips[1])
        assert_code(result, 302024)

        result = wrong_verison_declare(pip, pip.cfg.version5)
        assert_code(result, 302024)

        result = wrong_verison_declare(pip, pip.chain_version)
        assert_code(result, 302024)

    @pytest.mark.P2
    @allure.title('There is a preactive proposal, verifier declare version')
    def test_DE_VE_068(self, preactive_proposal_pips):
        pip = preactive_proposal_pips[0]
        result = replace_version_declare(pip, pip.cfg.PLATON_NEW_BIN8, pip.cfg.version8)
        assert_code(result, 302028)

        result = wrong_verisonsign_declare(pip, preactive_proposal_pips[1])
        assert_code(result, 302024)

        result = wrong_verison_declare(pip, pip.cfg.version5)
        assert_code(result, 302024)

        result = wrong_verison_declare(pip, pip.chain_version)
        assert_code(result, 302024)

    @pytest.mark.P2
    @allure.title('There is a preactive proposal, verifier declare version')
    def test_DE_VE_069(self, preactive_large_version_proposal_pips):
        pip = preactive_large_version_proposal_pips[0]
        result = replace_version_declare(pip, pip.cfg.PLATON_NEW_BIN3, pip.cfg.version3)
        assert_code(result, 302028)

        result = wrong_verisonsign_declare(pip, preactive_large_version_proposal_pips[1])
        assert_code(result, 302024)

        result = wrong_verison_declare(pip, pip.cfg.version8)
        assert_code(result, 302024)

        result = wrong_verison_declare(pip, pip.chain_version)
        assert_code(result, 302024)

    @pytest.mark.P2
    @allure.title('There is a preactive proposal, verifier declare version')
    def test_DE_VE_070(self, preactive_large_version_proposal_pips):
        pip = preactive_large_version_proposal_pips[0]
        result = replace_version_declare(pip, pip.cfg.PLATON_NEW_BIN, pip.cfg.version5)
        assert_code(result, 302028)

        result = wrong_verisonsign_declare(pip, preactive_large_version_proposal_pips[1])
        assert_code(result, 302024)

        result = wrong_verison_declare(pip, pip.cfg.version8)
        assert_code(result, 302024)

        result = wrong_verison_declare(pip, pip.chain_version)
        assert_code(result, 302024)

    @pytest.mark.P2
    @allure.title('There is a preactive proposal, verifier declare version')
    def test_DE_VE_071(self, preactive_large_version_proposal_pips):
        pip = preactive_large_version_proposal_pips[0]
        result = replace_version_declare(pip, pip.cfg.PLATON_NEW_BIN4, pip.cfg.version4)
        assert_code(result, 302028)

        result = wrong_verisonsign_declare(pip, preactive_large_version_proposal_pips[1])
        assert_code(result, 302024)

        result = wrong_verison_declare(pip, pip.cfg.version8)
        assert_code(result, 302024)

        result = wrong_verison_declare(pip, pip.chain_version)
        assert_code(result, 302024)

    @pytest.mark.P2
    @allure.title('There is a preactive proposal, verifier declare version')
    def test_DE_VE_072(self, preactive_large_version_proposal_pips):
        pip = preactive_large_version_proposal_pips[0]
        result = replace_version_declare(pip, pip.cfg.PLATON_NEW_BIN6, pip.cfg.version6)
        assert_code(result, 302028)

        result = wrong_verisonsign_declare(pip, preactive_large_version_proposal_pips[1])
        assert_code(result, 302024)

        result = wrong_verison_declare(pip, pip.cfg.version8)
        assert_code(result, 302024)

        result = wrong_verison_declare(pip, pip.chain_version)
        assert_code(result, 302024)

    @pytest.mark.P2
    @allure.title('There is a preactive proposal, verifier declare version')
    def test_DE_VE_073(self, preactive_large_version_proposal_pips):
        pip = preactive_large_version_proposal_pips[0]
        result = replace_version_declare(pip, pip.cfg.PLATON_NEW_BIN8, pip.cfg.version8)
        assert_code(result, 0)
        verifier_node_version(pip, pip.cfg.version8)

        result = wrong_verisonsign_declare(pip, preactive_large_version_proposal_pips[1])
        assert_code(result, 302024)

        result = wrong_verison_declare(pip, pip.cfg.version0)
        assert_code(result, 302024)


class TestNoProposalCA:
    @pytest.mark.P0
    @allure.title('No effective proposal, candiate declare version')
    def test_DE_CA_001(self, noproposal_candidate_pips, client_verifier):
        pip = noproposal_candidate_pips[0]
        verison = struct.pack('>I', pip.chain_version)
        if verison[3] != 0:
            result = replace_version_declare(pip, pip.cfg.PLATON_NEW_BIN2, pip.cfg.version2)
            assert_code(result, 0)
            verifier_node_version(pip, pip.cfg.version2)

            result = wrong_verisonsign_declare(pip, client_verifier.pip)
            assert_code(result, 302024)

            result = wrong_verison_declare(pip, pip.chain_version)
            assert_code(result, 302024)

    @pytest.mark.P3
    @allure.title('No effective proposal, candiate declare version')
    def test_DE_CA_002(self, noproposal_candidate_pips, client_verifier):
        pip = noproposal_candidate_pips[0]

        result = replace_version_declare(pip, pip.cfg.PLATON_NEW_BIN1, pip.cfg.version1)
        assert_code(result, 302028)

        result = wrong_verisonsign_declare(pip, client_verifier.pip)
        assert_code(result, 302024)

        result = wrong_verison_declare(pip, pip.chain_version)
        assert_code(result, 302024)

    @pytest.mark.P0
    @allure.title('No effective proposal, candiate declare version')
    def test_DE_CA_004(self, noproposal_candidate_pips, client_verifier):
        pip = noproposal_candidate_pips[0]

        result = replace_version_declare(pip, pip.cfg.PLATON_NEW_BIN3, pip.cfg.version3)
        assert_code(result, 0)
        verifier_node_version(pip, pip.cfg.version3)

        result = wrong_verisonsign_declare(pip, client_verifier.pip)
        assert_code(result, 302024)

        result = wrong_verison_declare(pip, pip.chain_version)
        assert_code(result, 302024)

    @pytest.mark.P2
    @allure.title('No effective proposal, candiate declare version')
    def test_DE_CA_005(self, noproposal_candidate_pips, client_verifier):
        pip = noproposal_candidate_pips[0]

        result = replace_version_declare(pip, pip.cfg.PLATON_NEW_BIN, pip.cfg.version5)
        assert_code(result, 302028)

        result = wrong_verisonsign_declare(pip, client_verifier.pip)
        assert_code(result, 302024)

        result = wrong_verison_declare(pip, pip.chain_version)
        assert_code(result, 302024)

    @pytest.mark.P2
    @allure.title('No effective proposal, candiate declare version')
    def test_DE_CA_006(self, noproposal_candidate_pips, client_verifier):
        pip = noproposal_candidate_pips[0]

        result = replace_version_declare(pip, pip.cfg.PLATON_NEW_BIN8, pip.cfg.version8)
        assert_code(result, 302028)

        result = wrong_verisonsign_declare(pip, client_verifier.pip)
        assert_code(result, 302024)

        result = wrong_verison_declare(pip, pip.chain_version)
        assert_code(result, 302024)

    @pytest.mark.P0
    @allure.title('No effective proposal, candiate declare version')
    def test_DE_CA_007(self, noproposal_candidate_pips, client_verifier):
        pip = noproposal_candidate_pips[0]
        verison = struct.pack('>I', pip.chain_version)
        if verison[3] != 0:
            result = replace_version_declare(pip, pip.cfg.PLATON_NEW_BIN0, pip.cfg.version0)
            assert_code(result, 0)
            verifier_node_version(pip, pip.cfg.version0)

            result = wrong_verisonsign_declare(pip, client_verifier.pip)
            assert_code(result, 302024)

            result = wrong_verison_declare(pip, pip.cfg.version2)
            assert_code(result, 302024)

            result = wrong_verison_declare(pip, pip.cfg.version3)
            assert_code(result, 302024)

    @pytest.mark.P0
    @allure.title('No effective proposal, candiate declare version')
    def test_DE_CA_008(self, proposal_candidate_pips, client_verifier):
        pip = proposal_candidate_pips[0]
        verison = struct.pack('>I', pip.chain_version)
        if verison[3] != 0:
            result = replace_version_declare(pip, pip.cfg.PLATON_NEW_BIN2, pip.cfg.version2)
            assert_code(result, 0)
            verifier_node_version(pip, pip.cfg.version2)

            result = wrong_verisonsign_declare(pip, client_verifier.pip)
            assert_code(result, 302024)

            result = wrong_verison_declare(pip)
            assert_code(result, 302024)

            result = wrong_verison_declare(pip, pip.chain_version)
            assert_code(result, 302024)

    @pytest.mark.P2
    @allure.title('No effective proposal, candiate declare version')
    def test_DE_CA_010(self, proposal_candidate_pips, client_verifier):
        pip = proposal_candidate_pips[0]

        result = replace_version_declare(pip, pip.cfg.PLATON_NEW_BIN1, pip.cfg.version1)
        assert_code(result, 302028)

        result = wrong_verisonsign_declare(pip, client_verifier.pip)
        assert_code(result, 302024)

        result = wrong_verison_declare(pip)
        assert_code(result, 302024)

        result = wrong_verison_declare(pip, pip.chain_version)
        assert_code(result, 302024)

    @pytest.mark.P1
    @allure.title('No effective proposal, candiate declare version')
    def test_DE_CA_014(self, large_version_proposal_candidate_pips, client_verifier):
        pip = large_version_proposal_candidate_pips[0]
        verison = struct.pack('>I', pip.chain_version)
        if verison[3] != 0:
            result = replace_version_declare(pip, pip.cfg.PLATON_NEW_BIN2, pip.cfg.version2)
            assert_code(result, 0)
            verifier_node_version(pip, pip.cfg.version2)

            result = wrong_verisonsign_declare(pip, client_verifier.pip)
            assert_code(result, 302024)

            result = wrong_verison_declare(pip)
            assert_code(result, 302024)

            result = wrong_verison_declare(pip, pip.chain_version)
            assert_code(result, 302024)

    @pytest.mark.P2
    @allure.title('No effective proposal, candiate declare version')
    def test_DE_CA_025(self, large_version_proposal_candidate_pips, client_verifier):
        pip = large_version_proposal_candidate_pips[0]

        result = replace_version_declare(pip, pip.cfg.PLATON_NEW_BIN1, pip.cfg.version1)
        assert_code(result, 302028)

        result = wrong_verisonsign_declare(pip, client_verifier.pip)
        assert_code(result, 302024)

        result = wrong_verison_declare(pip)
        assert_code(result, 302024)

        result = wrong_verison_declare(pip, pip.chain_version)
        assert_code(result, 302024)

    @pytest.mark.P0
    @allure.title('No effective proposal, candiate declare version')
    def test_DE_CA_032(self, proposal_candidate_pips, client_verifier):
        pip = proposal_candidate_pips[0]

        result = replace_version_declare(pip, pip.cfg.PLATON_NEW_BIN0, pip.cfg.version0)
        assert_code(result, 0)
        verifier_node_version(pip, pip.cfg.version0)

        result = wrong_verisonsign_declare(pip, client_verifier.pip)
        assert_code(result, 302024)

        result = wrong_verison_declare(pip)
        assert_code(result, 302024)

    @pytest.mark.P2
    @allure.title('No effective proposal, candiate declare version')
    def test_DE_CA_034(self, large_version_proposal_candidate_pips, client_verifier):
        pip = large_version_proposal_candidate_pips[0]

        result = replace_version_declare(pip, pip.cfg.PLATON_NEW_BIN0, pip.cfg.version0)
        assert_code(result, 0)
        verifier_node_version(pip, pip.cfg.version0)

        result = wrong_verisonsign_declare(pip, client_verifier.pip)
        assert_code(result, 302024)

        result = wrong_verison_declare(pip)
        assert_code(result, 302024)

    @pytest.mark.P1
    @allure.title('No effective proposal, candiate declare version')
    def test_DE_CA_036(self, proposal_candidate_pips, client_verifier):
        pip = proposal_candidate_pips[0]

        result = replace_version_declare(pip, pip.cfg.PLATON_NEW_BIN3, pip.cfg.version3)
        assert_code(result, 0)
        verifier_node_version(pip, pip.cfg.version3)

        result = wrong_verisonsign_declare(pip, client_verifier.pip)
        assert_code(result, 302024)

        result = wrong_verison_declare(pip)
        assert_code(result, 302024)

        result = wrong_verison_declare(pip, pip.chain_version)
        assert_code(result, 302024)

    @pytest.mark.P1
    @allure.title('No effective proposal, candiate declare version')
    def test_DE_CA_038(self, proposal_candidate_pips, client_verifier):
        pip = proposal_candidate_pips[0]

        node_version = verifier_node_version(pip)
        result = replace_version_declare(pip, pip.cfg.PLATON_NEW_BIN, pip.cfg.version5)
        assert_code(result, 0)
        verifier_node_version(pip, node_version)

        result = wrong_verisonsign_declare(pip, client_verifier.pip)
        assert_code(result, 302024)

        result = wrong_verison_declare(pip, pip.chain_version)
        assert_code(result, 302024)

    @pytest.mark.P2
    @allure.title('No effective proposal, candiate declare version')
    def test_DE_CA_040(self, proposal_candidate_pips, client_verifier):
        pip = proposal_candidate_pips[0]

        node_version = verifier_node_version(pip)
        result = replace_version_declare(pip, pip.cfg.PLATON_NEW_BIN4, pip.cfg.version4)
        assert_code(result, 0)
        verifier_node_version(pip, node_version)

        result = wrong_verisonsign_declare(pip, client_verifier.pip)
        assert_code(result, 302024)

        result = wrong_verison_declare(pip)
        assert_code(result, 302024)

        result = wrong_verison_declare(pip, pip.chain_version)
        assert_code(result, 302024)

    @pytest.mark.P2
    @allure.title('No effective proposal, candiate declare version')
    def test_DE_CA_042(self, proposal_candidate_pips, client_verifier):
        pip = proposal_candidate_pips[0]

        node_version = verifier_node_version(pip)
        result = replace_version_declare(pip, pip.cfg.PLATON_NEW_BIN6, pip.cfg.version6)
        assert_code(result, 0)
        verifier_node_version(pip, node_version)

        result = wrong_verisonsign_declare(pip, client_verifier.pip)
        assert_code(result, 302024)

        result = wrong_verison_declare(pip)
        assert_code(result, 302024)

        result = wrong_verison_declare(pip, pip.chain_version)
        assert_code(result, 302024)

    @pytest.mark.P2
    @allure.title('No effective proposal, candiate declare version')
    def test_DE_CA_044(self, proposal_candidate_pips, client_verifier):
        pip = proposal_candidate_pips[0]

        result = replace_version_declare(pip, pip.cfg.PLATON_NEW_BIN8, pip.cfg.version8)
        assert_code(result, 302028)

        result = wrong_verisonsign_declare(pip, client_verifier.pip)
        assert_code(result, 302024)

        result = wrong_verison_declare(pip)
        assert_code(result, 302024)

        result = wrong_verison_declare(pip, pip.chain_version)
        assert_code(result, 302024)

    @pytest.mark.P1
    @allure.title('No effective proposal, candiate declare version')
    def test_DE_CA_046(self, large_version_proposal_candidate_pips, client_verifier):
        pip = large_version_proposal_candidate_pips[0]

        result = replace_version_declare(pip, pip.cfg.PLATON_NEW_BIN3, pip.cfg.version3)
        assert_code(result, 0)
        verifier_node_version(pip, pip.cfg.version3)

        result = wrong_verisonsign_declare(pip, client_verifier.pip)
        assert_code(result, 302024)

        result = wrong_verison_declare(pip)
        assert_code(result, 302024)

        result = wrong_verison_declare(pip, pip.chain_version)
        assert_code(result, 302024)

    @pytest.mark.P1
    @allure.title('No effective proposal, candiate declare version')
    def test_DE_CA_048(self, large_version_proposal_candidate_pips, client_verifier):
        pip = large_version_proposal_candidate_pips[0]

        result = replace_version_declare(pip, pip.cfg.PLATON_NEW_BIN, pip.cfg.version5)
        assert_code(result, 302028)

        result = wrong_verisonsign_declare(pip, client_verifier.pip)
        assert_code(result, 302024)

        result = wrong_verison_declare(pip)
        assert_code(result, 302024)

        result = wrong_verison_declare(pip, pip.chain_version)
        assert_code(result, 302024)

    @pytest.mark.P1
    @allure.title('No effective proposal, candiate declare version')
    def test_DE_CA_050(self, large_version_proposal_candidate_pips, client_verifier):
        pip = large_version_proposal_candidate_pips[0]

        result = replace_version_declare(pip, pip.cfg.PLATON_NEW_BIN4, pip.cfg.version4)
        assert_code(result, 302028)

        result = wrong_verisonsign_declare(pip, client_verifier.pip)
        assert_code(result, 302024)

        result = wrong_verison_declare(pip)
        assert_code(result, 302024)

        result = wrong_verison_declare(pip, pip.chain_version)
        assert_code(result, 302024)

    @pytest.mark.P1
    @allure.title('No effective proposal, candiate declare version')
    def test_DE_CA_052(self, large_version_proposal_candidate_pips, client_verifier):
        pip = large_version_proposal_candidate_pips[0]

        result = replace_version_declare(pip, pip.cfg.PLATON_NEW_BIN6, pip.cfg.version6)
        assert_code(result, 302028)

        result = wrong_verisonsign_declare(pip, client_verifier.pip)
        assert_code(result, 302024)

        result = wrong_verison_declare(pip)
        assert_code(result, 302024)

        result = wrong_verison_declare(pip, pip.chain_version)
        assert_code(result, 302024)

    @pytest.mark.P2
    @allure.title('No effective proposal, candiate declare version')
    def test_DE_CA_054(self, large_version_proposal_candidate_pips, client_verifier):
        pip = large_version_proposal_candidate_pips[0]

        node_version = verifier_node_version(pip)
        result = replace_version_declare(pip, pip.cfg.PLATON_NEW_BIN8, pip.cfg.version8)
        assert_code(result, 0)
        verifier_node_version(pip, node_version)

        result = wrong_verisonsign_declare(pip, client_verifier.pip)
        assert_code(result, 302024)

        result = wrong_verison_declare(pip, pip.chain_version)
        assert_code(result, 302024)


class TestNewNodeDeclareVersion:
    @pytest.mark.P1
    @allure.title('New node declare version')
    def test_DE_NN_001_to_003(self, new_genesis_env, clients_consensus, clients_noconsensus):
        new_genesis_env.deploy_all()
        pip = clients_noconsensus[0].pip
        address, _ = pip.economic.account.generate_account(pip.node.web3, 10**18 * 10000000)
        result = pip.declareVersion(pip.node.node_id, address, transaction_cfg=pip.cfg.transaction_cfg)
        log.info('New node declare version result : {}'.format(result))
        assert_code(result, 302023)

        submitvpandvote(clients_consensus)
        proposalinfo = pip.get_effect_proposal_info_of_vote()
        log.info('Get version proposal information : {}'.format(proposalinfo))
        result = pip.declareVersion(pip.node.node_id, address, transaction_cfg=pip.cfg.transaction_cfg)
        log.info('New node declare version result : {}'.format(result))
        assert_code(result, 302023)

        wait_block_number(pip.node, proposalinfo.get('EndVotingBlock'))
        assert_code(pip.get_status_of_proposal(proposalinfo.get('ProposalID')), 4)
        wait_block_number(pip.node, proposalinfo.get('ActiveBlock'))
        assert_code(pip.get_status_of_proposal(proposalinfo.get('ProposalID')), 5)

        result = pip.declareVersion(pip.node.node_id, address, transaction_cfg=pip.cfg.transaction_cfg)
        log.info('New node declare version result : {}'.format(result))
        assert_code(result, 302023)


class TestDV:
    @pytest.mark.P3
    @allure.title('Declare version')
    def test_DE_VE_003_DE_VE_012_DE_VE_027_DE_CA_003_DE_CA_012_DE_VE_061_DE_CA_027(self, new_genesis_env,
                                                                                   clients_consensus):
        new_genesis_env.deploy_all()
        pip_ca = clients_consensus[-1].pip
        pip_ve = clients_consensus[0].pip
        submitvpandvote(clients_consensus[0:3], votingrounds=3, version=pip_ca.cfg.version9)
        proposalinfo = pip_ca.get_effect_proposal_info_of_vote()
        log.info("Get version proposal information : {}".format(proposalinfo))
        wait_block_number(pip_ca.node, proposalinfo.get('EndVotingBlock'))
        assert_code(pip_ca.get_status_of_proposal(proposalinfo.get('ProposalID')), 4)
        wait_block_number(pip_ca.node, proposalinfo.get('ActiveBlock'))
        assert_code(pip_ca.get_status_of_proposal(proposalinfo.get('ProposalID')), 5)
        assert pip_ca.cfg.version9 == pip_ca.chain_version

        verifier_list = get_pledge_list(clients_consensus[0].ppos.getVerifierList)
        log.info('verifier list : {}'.format(verifier_list))
        assert pip_ca.node not in verifier_list

        result = replace_version_declare(pip_ve, pip_ve.cfg.PLATON_NEW_BIN0, pip_ve.cfg.version0)
        assert_code(result, 302028)
        result = pip_ca.declareVersion(pip_ca.node.node_id, pip_ca.node.staking_address,
                                       transaction_cfg=pip_ca.cfg.transaction_cfg)
        log.info('Node {} declare version result {}'.format(pip_ca.node.node_id, result))
        assert_code(result, 302028)
        result = clients_consensus[1].pip.submitVersion(clients_consensus[1].node.node_id, str(time.time()),
                                                        pip_ca.cfg.version8, 4,
                                                        clients_consensus[1].node.staking_address,
                                                        transaction_cfg=pip_ca.cfg.transaction_cfg)
        log.info('Node {} submit version proposal result : {}'.format(clients_consensus[1].node.node_id, result))
        assert_code(result, 0)
        result = replace_version_declare(pip_ve, pip_ve.cfg.PLATON_NEW_BIN0, versiontag=pip_ve.cfg.version0)
        assert_code(result, 302028)

        result = replace_version_declare(pip_ca, pip_ve.cfg.PLATON_NEW_BIN0, versiontag=pip_ve.cfg.version0)
        assert_code(result, 302028)

        for client in clients_consensus[:3]:
            version_proposal_vote(client.pip)
        proposalinfo = pip_ve.get_effect_proposal_info_of_vote()
        log.info('Get proposal information : {}'.format(proposalinfo))
        wait_block_number(pip_ve.node, proposalinfo.get('EndVotingBlock'))
        assert_code(pip_ve.get_status_of_proposal(proposalinfo.get('ProposalID')), 4)
        wait_block_number(pip_ve.node, proposalinfo.get('ActiveBlock'))
        assert_code(pip_ve.get_status_of_proposal(proposalinfo.get('ProposalID')), 5)

        result = replace_version_declare(pip_ve, pip_ve.cfg.PLATON_NEW_BIN0, versiontag=pip_ve.cfg.version0)
        assert_code(result, 302028)


class TestVotedCADV:
    def get_candidate_no_verifier(self, client_list):
        verifier_list = get_pledge_list(client_list[0].ppos.getVerifierList)
        log.info('verifier list : {}'.format(verifier_list))
        candidate_list = get_pledge_list(client_list[0].ppos.getCandidateList)
        log.info('candidate list : {}'.format(candidate_list))
        for nodeid in candidate_list:
            if nodeid not in verifier_list:
                return get_client_by_nodeid(nodeid, client_list)
        raise Exception('There is not candidate no verifier node')

    @pytest.mark.P2
    @allure.title('Voted candidate, Declare version')
    def test_DE_CA_009_011_033_037_039_041_043_045(self, new_genesis_env, clients_consensus, clients_noconsensus,
                                                   all_clients):
        genesis = from_dict(data_class=Genesis, data=new_genesis_env.genesis_config)
        genesis.economicModel.gov.versionProposalVoteDurationSeconds = 2000
        new_genesis_env.set_genesis(genesis.to_dict())
        new_genesis_env.deploy_all()
        submitvpandvote(clients_consensus, votingrounds=40)
        createstaking(clients_noconsensus)
        clients_consensus[0].economic.wait_settlement_blocknum(clients_consensus[0].node)
        client = self.get_candidate_no_verifier(all_clients)
        pip = client.pip
        result = replace_version_declare(pip, pip.cfg.PLATON_NEW_BIN1, pip.cfg.version1)
        assert_code(result, 302028)

        result = wrong_verisonsign_declare(pip, clients_noconsensus[0].pip)
        assert_code(result, 302024)

        result = wrong_verison_declare(pip, pip.chain_version)
        assert_code(result, 302024)

        result = wrong_verison_declare(pip, pip.cfg.version5)
        assert_code(result, 302024)

        result = replace_version_declare(pip, pip.cfg.PLATON_NEW_BIN2, pip.cfg.version2)
        assert_code(result, 302028)

        result = wrong_verisonsign_declare(pip, clients_noconsensus[0].pip)
        assert_code(result, 302024)

        result = replace_version_declare(pip, pip.cfg.PLATON_NEW_BIN0, pip.cfg.version0)
        assert_code(result, 302028)

        result = wrong_verisonsign_declare(pip, clients_noconsensus[0].pip)
        assert_code(result, 302024)

        result = replace_version_declare(pip, pip.cfg.PLATON_NEW_BIN3, pip.cfg.version3)
        assert_code(result, 302028)

        result = wrong_verisonsign_declare(pip, clients_noconsensus[0].pip)
        assert_code(result, 302024)

        result = replace_version_declare(pip, pip.cfg.PLATON_NEW_BIN, pip.cfg.version5)
        assert_code(result, 0)

        result = wrong_verisonsign_declare(pip, clients_noconsensus[0].pip)
        assert_code(result, 302024)

        result = replace_version_declare(pip, pip.cfg.PLATON_NEW_BIN4, pip.cfg.version4)
        assert_code(result, 0)

        result = wrong_verisonsign_declare(pip, clients_noconsensus[0].pip)
        assert_code(result, 302024)

        result = replace_version_declare(pip, pip.cfg.PLATON_NEW_BIN6, pip.cfg.version6)
        assert_code(result, 0)

        result = wrong_verisonsign_declare(pip, clients_noconsensus[0].pip)
        assert_code(result, 302024)

        result = replace_version_declare(pip, pip.cfg.PLATON_NEW_BIN8, pip.cfg.version8)
        assert_code(result, 302028)

        result = wrong_verisonsign_declare(pip, clients_noconsensus[0].pip)
        assert_code(result, 302024)

    @pytest.mark.P2
    @allure.title('Voted candidate, Declare version')
    def test_DE_CA_021_026_035_047_049_051_053_055(self, new_genesis_env, clients_consensus, clients_noconsensus,
                                                   all_clients):
        verison = struct.pack('>I', clients_consensus[0].pip.chain_version)
        if verison[3] != 0:
            genesis = from_dict(data_class=Genesis, data=new_genesis_env.genesis_config)
            genesis.economicModel.gov.versionProposalVoteDurationSeconds = 2000
            new_genesis_env.set_genesis(genesis.to_dict())
            new_genesis_env.deploy_all()
            submitvpandvote(clients_consensus, votingrounds=40, version=clients_noconsensus[0].pip.cfg.version8)
            createstaking(clients_noconsensus)
            clients_consensus[0].economic.wait_settlement_blocknum(clients_consensus[0].node)
            client = self.get_candidate_no_verifier(all_clients)
            pip = client.pip
            result = replace_version_declare(pip, pip.cfg.PLATON_NEW_BIN2, pip.cfg.version2)
            assert_code(result, 302028)

            result = wrong_verisonsign_declare(pip, clients_noconsensus[0].pip)
            assert_code(result, 302024)

            result = wrong_verison_declare(pip, pip.chain_version)
            assert_code(result, 302024)

            result = wrong_verison_declare(pip, pip.cfg.version5)
            assert_code(result, 302024)

            result = replace_version_declare(pip, pip.cfg.PLATON_NEW_BIN1, pip.cfg.version1)
            assert_code(result, 302028)

            result = wrong_verisonsign_declare(pip, clients_noconsensus[0].pip)
            assert_code(result, 302024)

            result = replace_version_declare(pip, pip.cfg.PLATON_NEW_BIN0, pip.cfg.version0)
            assert_code(result, 302028)

            result = wrong_verisonsign_declare(pip, clients_noconsensus[0].pip)
            assert_code(result, 302024)

            result = replace_version_declare(pip, pip.cfg.PLATON_NEW_BIN3, pip.cfg.version3)
            assert_code(result, 302028)

            result = wrong_verisonsign_declare(pip, clients_noconsensus[0].pip)
            assert_code(result, 302024)

            result = replace_version_declare(pip, pip.cfg.PLATON_NEW_BIN, pip.cfg.version5)
            assert_code(result, 302028)

            result = wrong_verisonsign_declare(pip, clients_noconsensus[0].pip)
            assert_code(result, 302024)

            result = replace_version_declare(pip, pip.cfg.PLATON_NEW_BIN4, pip.cfg.version4)
            assert_code(result, 302028)

            result = wrong_verisonsign_declare(pip, clients_noconsensus[0].pip)
            assert_code(result, 302024)

            result = replace_version_declare(pip, pip.cfg.PLATON_NEW_BIN6, pip.cfg.version6)
            assert_code(result, 302028)

            result = wrong_verisonsign_declare(pip, clients_noconsensus[0].pip)
            assert_code(result, 302024)

            result = replace_version_declare(pip, pip.cfg.PLATON_NEW_BIN8, pip.cfg.version8)
            assert_code(result, 0)

            result = wrong_verisonsign_declare(pip, clients_noconsensus[0].pip)
            assert_code(result, 302024)


@pytest.mark.P2
@allure.title('Voted verifier, replace the platon bin and declare version')
def test_DE_VE_074(no_vp_proposal, client_verifier):
    pip = client_verifier.pip
    submitvpandvote([client_verifier], votingrounds=2)
    proposalinfo = pip.get_effect_proposal_info_of_vote()
    log.info('Get proposal information : {}'.format(proposalinfo))
    wait_block_number(pip.node, proposalinfo.get('EndVotingBlock'))
    assert_code(pip.get_status_of_proposal(proposalinfo.get('ProposalID')), 3)
    result = replace_version_declare(pip, pip.cfg.PLATON_NEW_BIN0, pip.cfg.version0)
    assert_code(result, 0)
