from common.log import log
from tests.lib.utils import upload_platon, get_pledge_list, assert_code
from tests.lib.client import get_client_obj_list
import pytest
import time

@pytest.fixture()
def noproposal_pipobj_list(global_test_env, client_list_obj):
    '''
    获取验证节点Client对象列表
    :param global_test_env:
    :return:
    '''
    if client_list_obj[0].pip.is_exist_effective_proposal():
        log.info('There is effective proposal, Restart the chain')
        global_test_env.deploy_all()
    verifier_list = get_pledge_list(client_list_obj[0].ppos.getVerifierList)
    log.info('verifierlist{}'.format(verifier_list))
    client_obj_list = get_client_obj_list(verifier_list, client_list_obj)
    return [client_obj.pip for client_obj in client_obj_list]

@pytest.fixture()
def proposal_pipobj_list(global_test_env, client_verifier_obj, client_list_obj):
    '''
    get verifier Client object list
    :param global_test_env:
    :return:
    '''
    pip_obj = client_verifier_obj.pip
    if pip_obj.is_exist_effective_proposal():
        proposalinfo = pip_obj.get_effect_proposal_info_of_vote()
        log.info('proprosalinfo : {}'.format(proposalinfo))
        if proposalinfo.get('EndVotingBlock') - pip_obj.node.block_number > 2 * pip_obj.economic.consensus_size \
                and proposalinfo.get('NewVersion') == pip_obj.cfg.version5:
            verifier_list = get_pledge_list(client_list_obj[0].ppos.getVerifierList)
            log.info('verifierlist{}'.format(verifier_list))
            client_obj_list = get_client_obj_list(verifier_list, client_list_obj)
            return [client_obj.pip for client_obj in client_obj_list]
        else:
            global_test_env.deploy_all()
    result = pip_obj.submitVersion(pip_obj.node.node_id, str(time.time_ns()), pip_obj.cfg.version5, 10,
                                   pip_obj.node.staking_address, transaction_cfg=pip_obj.cfg.transaction_cfg)
    log.info('version proposal result :{}'.format(result))
    verifier_list = get_pledge_list(client_list_obj[0].ppos.getVerifierList)
    log.info('verifierlist{}'.format(verifier_list))
    client_obj_list = get_client_obj_list(verifier_list, client_list_obj)
    return [client_obj.pip for client_obj in client_obj_list]

@pytest.fixture()
def big_version_proposal_pipobj_list(global_test_env, client_verifier_obj, client_list_obj):
    '''
    get verifier Client object list
    :param global_test_env:
    :return:
    '''
    pip_obj = client_verifier_obj.pip
    if pip_obj.is_exist_effective_proposal():
        proposalinfo = pip_obj.get_effect_proposal_info_of_vote()
        log.info('proprosalinfo : {}'.format(proposalinfo))
        if proposalinfo.get('EndVotingBlock') - pip_obj.node.block_number > 2 * pip_obj.economic.consensus_size \
                and proposalinfo.get('NewVersion') == pip_obj.cfg.version8:
            verifier_list = get_pledge_list(client_list_obj[0].ppos.getVerifierList)
            log.info('verifierlist{}'.format(verifier_list))
            client_obj_list = get_client_obj_list(verifier_list, client_list_obj)
            return [client_obj.pip for client_obj in client_obj_list]
        else:
            global_test_env.deploy_all()
    result = pip_obj.submitVersion(pip_obj.node.node_id, str(time.time_ns()), pip_obj.cfg.version8, 10,
                                   pip_obj.node.staking_address, transaction_cfg=pip_obj.cfg.transaction_cfg)
    log.info('version proposal result :{}'.format(result))
    verifier_list = get_pledge_list(client_list_obj[0].ppos.getVerifierList)
    log.info('verifierlist{}'.format(verifier_list))
    client_obj_list = get_client_obj_list(verifier_list, client_list_obj)
    return [client_obj.pip for client_obj in client_obj_list]

def replace_version_declare(pip_obj, platon_bin, versiontag):
    upload_platon(pip_obj.node, platon_bin)
    log.info('Replace the platon of the node {}'.format(versiontag))
    pip_obj.node.restart()
    log.info('Restart the node{}'.format(pip_obj.node.node_id))
    assert pip_obj.node.program_version == versiontag
    log.info('assert the version of the node is {}'.format(versiontag))
    result = pip_obj.declareVersion(pip_obj.node.node_id, pip_obj.node.staking_address,
                                    transaction_cfg=pip_obj.cfg.transaction_cfg)
    log.info('declareversion {} result: {}'.format(pip_obj.node.program_version, result))
    return result

def wrong_verisonsign_declare(pip_obj, pip_obj_test):
    result = pip_obj.declareVersion(pip_obj.node.node_id, pip_obj.node.staking_address,
                                    version_sign=pip_obj_test.node.program_version_sign,
                                    transaction_cfg=pip_obj.cfg.transaction_cfg)
    log.info('wrong program version sign, declareVersion result : {}'.format(result))
    return result

def wrong_verison_declare(pip_obj, version=None):
    if not version:
        proposalinfo = pip_obj.get_effect_proposal_info_of_vote()
        version = proposalinfo.get('NewVersion')
        log.info('The new version of the proposal: {}'.format(version))
    result = pip_obj.declareVersion(pip_obj.node.node_id, pip_obj.node.staking_address,
                                    program_version=version,
                                    transaction_cfg=pip_obj.cfg.transaction_cfg)
    log.info('wrong program version, declareVersion: {} result : {}'.format(version, result))
    return result

def test_DE_DE_001(client_verifier_obj):
    pip_obj = client_verifier_obj.pip
    address, _ = pip_obj.economic.account.generate_account(pip_obj.node.web3, 10**18 * 10000)
    result = pip_obj.declareVersion(pip_obj.node.node_id, address, transaction_cfg=pip_obj.cfg.transaction_cfg)
    log.info('declareVersion result: {}'.format(result))
    assert result.get('Code') == 302021

class TestNoProposalVE():
    def test_DE_VE_001(self, noproposal_pipobj_list):
        pip_obj = noproposal_pipobj_list[0]
        result = replace_version_declare(pip_obj, pip_obj.cfg.PLATON_NEW_BIN2, pip_obj.cfg.version2)
        assert_code(result, 0)

        result = wrong_verisonsign_declare(pip_obj, noproposal_pipobj_list[1])
        assert_code(result, 302024)

        result = wrong_verison_declare(pip_obj, pip_obj.chain_version)
        assert_code(result, 302024)

    def test_DE_VE_002(self, noproposal_pipobj_list):
        pip_obj = noproposal_pipobj_list[0]
        result = replace_version_declare(pip_obj, pip_obj.cfg.PLATON_NEW_BIN1, pip_obj.cfg.version1)
        assert_code(result, 302028)

        result = wrong_verisonsign_declare(pip_obj, noproposal_pipobj_list[1])
        assert_code(result, 302024)

        result = wrong_verison_declare(pip_obj, pip_obj.chain_version)
        assert_code(result, 302024)

    @pytest.mark.P0
    def test_DE_VE_004(self, noproposal_pipobj_list):
        pip_obj = noproposal_pipobj_list[0]
        result = replace_version_declare(pip_obj, pip_obj.cfg.PLATON_NEW_BIN3, pip_obj.cfg.version3)
        assert_code(result, 0)

        result = wrong_verisonsign_declare(pip_obj, noproposal_pipobj_list[1])
        assert_code(result, 302024)

        result = wrong_verison_declare(pip_obj, pip_obj.chain_version)
        assert_code(result, 302024)

    def test_DE_VE_005(self, noproposal_pipobj_list):
        pip_obj = noproposal_pipobj_list[0]
        result = replace_version_declare(pip_obj, pip_obj.cfg.PLATON_NEW_BIN5, pip_obj.cfg.version5)
        assert_code(result, 302028)

        result = wrong_verisonsign_declare(pip_obj, noproposal_pipobj_list[1])
        assert_code(result, 302024)

        result = wrong_verison_declare(pip_obj, pip_obj.chain_version)
        assert_code(result, 302024)

    def test_DE_VE_006(self, noproposal_pipobj_list):
        pip_obj = noproposal_pipobj_list[0]
        result = replace_version_declare(pip_obj, pip_obj.cfg.PLATON_NEW_BIN8, pip_obj.cfg.version8)
        assert_code(result, 302028)

        result = wrong_verisonsign_declare(pip_obj, noproposal_pipobj_list[1])
        assert_code(result, 302024)

        result = wrong_verison_declare(pip_obj, pip_obj.chain_version)
        assert_code(result, 302024)

    def test_DE_VE_007(self, noproposal_pipobj_list):
        pip_obj = noproposal_pipobj_list[0]
        result = replace_version_declare(pip_obj, pip_obj.cfg.PLATON_NEW_BIN0, pip_obj.cfg.version0)
        assert_code(result, 0)

        result = wrong_verisonsign_declare(pip_obj, noproposal_pipobj_list[1])
        assert_code(result, 302024)

        result = wrong_verison_declare(pip_obj, pip_obj.chain_version)
        assert_code(result, 302024)

class TestHasProposalVE():
    def test_DE_VE_008(self, proposal_pipobj_list):
        pip_obj = proposal_pipobj_list[0]
        result = replace_version_declare(pip_obj, pip_obj.cfg.PLATON_NEW_BIN2, pip_obj.cfg.version2)
        assert_code(result, 0)

        result = wrong_verisonsign_declare(pip_obj, proposal_pipobj_list[1])
        assert_code(result, 302024)

        result = wrong_verison_declare(pip_obj, pip_obj.chain_version)
        assert_code(result, 302024)

        result = wrong_verison_declare(pip_obj)
        assert_code(result, 302024)

    def test_DE_VE_010(self, proposal_pipobj_list):
        pip_obj = proposal_pipobj_list[0]
        result = replace_version_declare(pip_obj, pip_obj.cfg.PLATON_NEW_BIN1, pip_obj.cfg.version1)
        assert_code(result, 302028)

        result = wrong_verisonsign_declare(pip_obj, proposal_pipobj_list[1])
        assert_code(result, 302024)

        result = wrong_verison_declare(pip_obj, pip_obj.chain_version)
        assert_code(result, 302024)

        result = wrong_verison_declare(pip_obj)
        assert_code(result, 302024)

    def test_DE_VE_014(self, big_version_proposal_pipobj_list):
        pip_obj = big_version_proposal_pipobj_list[0]
        result = replace_version_declare(pip_obj, pip_obj.cfg.PLATON_NEW_BIN2, pip_obj.cfg.version2)
        assert_code(result, 0)

        result = wrong_verisonsign_declare(pip_obj, big_version_proposal_pipobj_list[1])
        assert_code(result, 302024)

        result = wrong_verison_declare(pip_obj, pip_obj.chain_version)
        assert_code(result, 302024)

        result = wrong_verison_declare(pip_obj)
        assert_code(result, 302024)

    def test_DE_VE_025(self, big_version_proposal_pipobj_list):
        pip_obj = big_version_proposal_pipobj_list[0]
        result = replace_version_declare(pip_obj, pip_obj.cfg.PLATON_NEW_BIN1, pip_obj.cfg.version1)
        assert_code(result, 302028)

        result = wrong_verisonsign_declare(pip_obj, big_version_proposal_pipobj_list[1])
        assert_code(result, 302024)

        result = wrong_verison_declare(pip_obj, pip_obj.chain_version)
        assert_code(result, 302024)

        result = wrong_verison_declare(pip_obj)
        assert_code(result, 302024)
