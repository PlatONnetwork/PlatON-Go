import pytest
from common.log import log
from copy import copy
import time
from tests.lib.client import get_client_obj, get_client_obj_list
from tests.lib.utils import get_pledge_list, upload_platon, wait_block_number, assert_code

def version_proposal_vote(pip_obj, vote_option=None):
    proposalinfo = pip_obj.get_effect_proposal_info_of_vote()
    log.info('proposalinfo: {}'.format(proposalinfo))
    if not proposalinfo:
        raise Exception('there is no voting version proposal')
    if proposalinfo.get('NewVersion') == pip_obj.cfg.version5:
        upload_platon(pip_obj.node, pip_obj.cfg.PLATON_NEW_BIN)
        log.info('Replace the node {} version to {}'.format(pip_obj.node.node_id, pip_obj.cfg.version5))
    elif proposalinfo.get('NewVersion') == pip_obj.cfg.version8:
        upload_platon(pip_obj.node, pip_obj.cfg.PLATON_NEW_BIN8)
        log.info('Replace the node {} version to {}'.format(pip_obj.node.node_id, pip_obj.cfg.version8))
    else:
        raise Exception('The new version of the proposal is{}'.format(proposalinfo.get('NewVersion')))
    pip_obj.node.restart()
    log.info('Restart the node {}'.format(pip_obj.node.node_id))
    if not vote_option:
        vote_option = pip_obj.cfg.vote_option_yeas
    result = pip_obj.vote(pip_obj.node.node_id, proposalinfo.get('ProposalID'), vote_option,
                          pip_obj.node.staking_address, transaction_cfg=pip_obj.cfg.transaction_cfg)
    log.info('The node {} vote result {}'.format(pip_obj.node.node_id, result))
    assert_code(result, 0)


@pytest.fixture(scope="class")
def pip_env(global_test_env):
    cfg_copy = copy(global_test_env.cfg)
    yield global_test_env
    # global_test_env.set_cfg(cfg_copy)
    # global_test_env.deploy_all()

@pytest.fixture()
def no_version_proposal(global_test_env, client_verifier_obj):
    pip_obj = client_verifier_obj.pip
    if pip_obj.is_exist_effective_proposal() or pip_obj.chain_version != pip_obj.cfg.version0:
        log.info('There is effective proposal,restart the chain')
        global_test_env.deploy_all()
    return pip_obj

@pytest.fixture()
def submit_version(no_version_proposal):
    pip_obj = no_version_proposal
    result = pip_obj.submitVersion(pip_obj.node.node_id, str(time.time()), pip_obj.cfg.version5, 10,
                                   pip_obj.node.staking_address,
                                   transaction_cfg=pip_obj.cfg.transaction_cfg)
    log.info('submit version result:'.format(result))
    assert_code(result, 0)
    return pip_obj

@pytest.fixture()
def submit_cancel(submit_version):
    pip_obj = submit_version
    propolsalinfo = pip_obj.get_effect_proposal_info_of_vote()
    log.info('获取处于投票期的升级提案信息{}'.format(propolsalinfo))
    result = pip_obj.submitCancel(pip_obj.node.node_id, str(time.time()), 4, propolsalinfo.get('ProposalID'),
                                  pip_obj.node.staking_address, transaction_cfg=pip_obj.cfg.transaction_cfg)
    log.info('发起取消提案结果为{}'.format(result))
    assert_code(result, 0)
    return pip_obj

@pytest.fixture()
def submit_text(client_verifier_obj):
    pip_obj = client_verifier_obj.pip
    result = pip_obj.submitText(pip_obj.node.node_id, str(time.time()), pip_obj.node.staking_address,
                                transaction_cfg=pip_obj.cfg.transaction_cfg)
    log.info('submit text result:'.format(result))
    assert_code(result, 0)
    return pip_obj

@pytest.fixture()
def new_node_has_proposal(global_test_env, client_new_node_obj, client_verifier_obj, client_noconsensus_obj):
    pip_obj = client_verifier_obj.pip
    if pip_obj.chain_version != pip_obj.cfg.version0:
        global_test_env.deploy_all()
    if pip_obj.is_exist_effective_proposal():
        proposalinfo = pip_obj.get_effect_proposal_info_of_vote()
        log.info('升级提案信息为{}'.format(proposalinfo))
        if proposalinfo.get('EndVotingBlock') - pip_obj.node.block_number < 2 * pip_obj.economic.consensus_size:
            global_test_env.deploy_all()
            result = pip_obj.submitVersion(pip_obj.node.node_id, str(time.time()), pip_obj.cfg.version5, 5,
                                           pip_obj.node.staking_address,
                                           transaction_cfg=pip_obj.cfg.transaction_cfg)
            assert_code(result, 0)
            return client_noconsensus_obj.pip
        else:
            return client_new_node_obj.pip
    result = pip_obj.submitVersion(pip_obj.node.node_id, str(time.time()), pip_obj.cfg.version5, 5,
                          pip_obj.node.staking_address, transaction_cfg=pip_obj.cfg.transaction_cfg)
    assert_code(result, 0)
    return client_new_node_obj.pip

@pytest.fixture()
def candidate_has_proposal(global_test_env, client_candidate_obj, client_verifier_obj, client_list_obj):
    pip_obj = client_verifier_obj.pip
    if pip_obj.chain_version != pip_obj.cfg.version0:
        global_test_env.deploy_all()
    if pip_obj.is_exist_effective_proposal():
        proposalinfo = pip_obj.get_effect_proposal_info_of_vote()
        log.info('升级提案信息为{}'.format(proposalinfo))
        if proposalinfo.get('EndVotingBlock') - pip_obj.node.block_number < 2 * pip_obj.economic.consensus_size:
            global_test_env.deploy_all()
            normal_node_obj_list = global_test_env.normal_node_list
            for normal_node_obj in normal_node_obj_list:
                client_obj = get_client_obj(normal_node_obj.node_id, client_list_obj)
                address, _ = pip_obj.economic.account.generate_account(pip_obj.node.web3, 10 ** 18 * 10000000)
                log.info('对节点{}进行质押操作'.format(normal_node_obj.node_id))
                result = client_obj.staking.create_staking(0, address, address)
                log.info('节点{}质押结果为{}'.format(normal_node_obj.node_id, result))
                assert_code(result, 0)
            pip_obj.economic.wait_settlement_blocknum(pip_obj.node)
            node_id_list = pip_obj.get_candidate_list_not_verifier()
            if not node_id_list:
                raise Exception('获取候选人失败')
            client_candidate_obj = get_client_obj(node_id_list[0], client_list_obj)
        else:
            return client_candidate_obj.pip
    result = pip_obj.submitVersion(pip_obj.node.node_id, str(time.time()), pip_obj.cfg.version5, 5,
                          pip_obj.node.staking_address, transaction_cfg=pip_obj.cfg.transaction_cfg)
    log.info('发起升级提案结果为{}'.format(result))
    assert_code(result, 0)
    return client_candidate_obj.pip

@pytest.fixture()
def noproposal_pipobj_list(global_test_env, client_list_obj):
    '''
    获取验证节点Client对象列表
    :param global_test_env:
    :return:
    '''
    if client_list_obj[0].pip.is_exist_effective_proposal() or client_list_obj[0].chain_version != client_list_obj[0].cfg.version0:
        log.info('There is effective proposal, Restart the chain')
        global_test_env.deploy_all()
    verifier_list = get_pledge_list(client_list_obj[0].ppos.getVerifierList)
    log.info('verifierlist{}'.format(verifier_list))
    client_obj_list = get_client_obj_list(verifier_list, client_list_obj)
    return [client_obj.pip for client_obj in client_obj_list]

@pytest.fixture()
def noproposal_ca_pipobj_list(global_test_env, client_list_obj, client_noc_list_obj):
    '''
    获取验证节点Client对象列表
    :param global_test_env:
    :return:
    '''
    if client_list_obj[0].pip.is_exist_effective_proposal() or client_list_obj[0].pip.chain_version != \
            client_list_obj[0].pip.cfg.version0:
        log.info('There is effective proposal, Restart the chain')
        global_test_env.deploy_all()
    nodeid_list = client_list_obj[0].pip.get_candidate_list_not_verifier()
    log.info('candidate not verifier list {}'.format(nodeid_list))
    if not nodeid_list:
        candidate_list = get_pledge_list(client_list_obj[0].ppos.getCandidateList)
        log.info('candidate_list{}'.format(candidate_list))
        for client_obj in client_noc_list_obj:
            if client_obj.node.node_id not in candidate_list:
                address, _ = client_obj.economic.account.generate_account(client_obj.node.web3, 10**18 * 10000000)
                result = client_obj.staking.create_staking(0, address, address)
                log.info('node {} staking result {}'.format(client_obj.node.node_id, result))
        client_obj.economic.wait_settlement_blocknum(client_obj.node)
        nodeid_list = client_list_obj[0].pip.get_candidate_list_not_verifier()
        if not nodeid_list:
            raise Exception('get candidate not verifier failed')
    client_ca_list = get_client_obj_list(nodeid_list, client_list_obj)
    return [client_ca_obj.pip for client_ca_obj in client_ca_list]

@pytest.fixture()
def proposal_pipobj_list(global_test_env, client_verifier_obj, client_list_obj):
    '''
    get verifier Client object list
    :param global_test_env:
    :return:
    '''
    pip_obj = client_verifier_obj.pip
    if pip_obj.chain_version != pip_obj.cfg.version0:
        global_test_env.deploy_all()
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

def proposal_ca_pipobj_list(global_test_env, client_list_obj, client_noc_list_obj):
    '''
    There is voting stage proposal, get candidate list pip object
    :param global_test_env:
    :return:
    '''
    pip_obj = client_list_obj[0].pip
    if pip_obj.chain_version != pip_obj.cfg.version0:
        log.info('The chain has been upgraded,restart!')
        global_test_env.deploy_all()
    if pip_obj.is_exist_effective_proposal:
        if pip_obj.is_exist_effective_proposal_for_vote():
            proposalinfo = pip_obj.get_effect_proposal_info_of_vote()
            log.info('get version proposalinfo : {}'.format(proposalinfo))
            if proposalinfo.get('EndVotingBlock') - pip_obj.node.block_number > pip_obj.economic.consensus_size * 2:
                return pip_obj

    if client_list_obj[0].pip.is_exist_effective_proposal() or client_list_obj[0].pip.chain_version != \
            client_list_obj[0].pip.cfg.version0:
        log.info('There is effective proposal, Restart the chain')
        global_test_env.deploy_all()
    nodeid_list = client_list_obj[0].pip.get_candidate_list_not_verifier()
    log.info('candidate not verifier list {}'.format(nodeid_list))
    if not nodeid_list:
        candidate_list = get_pledge_list(client_list_obj[0].ppos.getCandidateList)
        log.info('candidate_list{}'.format(candidate_list))
        for client_obj in client_noc_list_obj:
            if client_obj.node.node_id not in candidate_list:
                address, _ = client_obj.economic.account.generate_account(client_obj.node.web3, 10**18 * 10000000)
                result = client_obj.staking.create_staking(0, address, address)
                log.info('node {} staking result {}'.format(client_obj.node.node_id, result))
        client_obj.economic.wait_settlement_blocknum(client_obj.node)
        nodeid_list = client_list_obj[0].pip.get_candidate_list_not_verifier()
        if not nodeid_list:
            raise Exception('get candidate not verifier failed')
    client_ca_list = get_client_obj_list(nodeid_list, client_list_obj)
    return [client_ca_obj.pip for client_ca_obj in client_ca_list]

@pytest.fixture()
def bv_proposal_pipobj_list(global_test_env, client_verifier_obj, client_list_obj):
    '''
    get verifier Client object list
    :param global_test_env:
    :return:
    '''
    pip_obj = client_verifier_obj.pip
    if pip_obj.chain_version != pip_obj.cfg.version0:
        global_test_env.deploy_all()
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

@pytest.fixture()
def proposal_voted_pipobj_list(global_test_env, client_verifier_obj, client_list_obj):
    '''
    get verifier Client object list
    :param global_test_env:
    :return:
    '''
    pip_obj = client_verifier_obj.pip
    global_test_env.deploy_all()
    result = pip_obj.submitVersion(pip_obj.node.node_id, str(time.time_ns()), pip_obj.cfg.version5, 10,
                                   pip_obj.node.staking_address, transaction_cfg=pip_obj.cfg.transaction_cfg)
    log.info('version proposal result :{}'.format(result))
    verifier_list = get_pledge_list(client_list_obj[0].ppos.getVerifierList)
    log.info('verifierlist{}'.format(verifier_list))
    client_obj_list = get_client_obj_list(verifier_list, client_list_obj)
    version_proposal_vote(client_obj_list[0].pip)
    assert_code(result, 0)
    return [client_obj.pip for client_obj in client_obj_list]

@pytest.fixture()
def bv_proposal_voted_pipobj_list(global_test_env, client_verifier_obj, client_list_obj):
    '''
    get verifier Client object list
    :param global_test_env:
    :return:
    '''
    pip_obj = client_verifier_obj.pip
    global_test_env.deploy_all()
    result = pip_obj.submitVersion(pip_obj.node.node_id, str(time.time_ns()), pip_obj.cfg.version8, 10,
                                   pip_obj.node.staking_address, transaction_cfg=pip_obj.cfg.transaction_cfg)
    log.info('version proposal result :{}'.format(result))
    verifier_list = get_pledge_list(client_list_obj[0].ppos.getVerifierList)
    log.info('verifierlist{}'.format(verifier_list))
    client_obj_list = get_client_obj_list(verifier_list, client_list_obj)
    version_proposal_vote(client_obj_list[0].pip)
    return [client_obj.pip for client_obj in client_obj_list]

@pytest.fixture()
def preactive_proposal_pipobj_list(global_test_env, client_list_obj):
    if client_list_obj[0].pip.is_exist_effective_proposal() or client_list_obj[0].pip.chain_version != \
            client_list_obj[0].pip.cfg.version0:
        log.info('There is effective version proposal, restart the chain')
        global_test_env.deploy_all()
    verifier_list = get_pledge_list(client_list_obj[0].ppos.getVerifierList)
    log.info('verifierlist :{}'.format(verifier_list))
    client_verifier_list_obj = get_client_obj_list(verifier_list, client_list_obj)
    pip_list_obj = [client_obj.pip for client_obj in client_verifier_list_obj]
    result = pip_list_obj[0].submitVersion(pip_list_obj[0].node.node_id, str(time.time_ns()),
                                           pip_list_obj[0].cfg.version5, 2, pip_list_obj[0].node.staking_address,
                                  transaction_cfg=pip_list_obj[0].cfg.transaction_cfg)
    log.info('submit version proposal, result : {}'.format(result))
    proposalinfo = pip_list_obj[0].get_effect_proposal_info_of_vote()
    log.info('Version proposalinfo: {}'.format(proposalinfo))
    for pip_obj in pip_list_obj:
        version_proposal_vote(pip_obj)
    wait_block_number(pip_obj.node, proposalinfo.get('EndVotingBlock'))
    assert pip_obj.get_status_of_proposal(proposalinfo.get('ProposalID')) == 4
    return pip_list_obj


@pytest.fixture()
def preactive_bv_proposal_pipobj_list(global_test_env, client_list_obj):
    if client_list_obj[0].pip.is_exist_effective_proposal() or client_list_obj[0].pip.chain_version != \
            client_list_obj[0].pip.cfg.version0:
        log.info('There is effective version proposal, restart the chain')
        global_test_env.deploy_all()
    verifier_list = get_pledge_list(client_list_obj[0].ppos.getVerifierList)
    log.info('verifierlist :{}'.format(verifier_list))
    client_verifier_list_obj = get_client_obj_list(verifier_list, client_list_obj)
    pip_list_obj = [client_obj.pip for client_obj in client_verifier_list_obj]
    result = pip_list_obj[0].submitVersion(pip_list_obj[0].node.node_id, str(time.time_ns()),
                                           pip_list_obj[0].cfg.version8, 2, pip_list_obj[0].node.staking_address,
                                  transaction_cfg=pip_list_obj[0].cfg.transaction_cfg)
    log.info('submit version proposal, result : {}'.format(result))
    proposalinfo = pip_list_obj[0].get_effect_proposal_info_of_vote()
    log.info('Version proposalinfo: {}'.format(proposalinfo))
    for pip_obj in pip_list_obj:
        version_proposal_vote(pip_obj)
    wait_block_number(pip_obj.node, proposalinfo.get('EndVotingBlock'))
    assert pip_obj.get_status_of_proposal(proposalinfo.get('ProposalID')) == 4
    return pip_list_obj


