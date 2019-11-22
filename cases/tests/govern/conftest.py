import pytest
from common.log import log
import time
import math
from tests.lib.client import get_client_obj, get_client_obj_list
from tests.lib.utils import get_pledge_list, upload_platon, wait_block_number, assert_code, get_governable_parameter_value
from typing import List
from tests.lib import Pip

def get_refund_to_account_block(pip, blocknumber=None):
    '''
    Get refund to account block
    :param pip:
    :return:
    '''
    if blocknumber is None:
        blocknumber = pip.node.block_number
    return math.ceil(blocknumber / pip.economic.settlement_size + pip.economic.unstaking_freeze_ratio
                     ) * pip.economic.settlement_size


def version_proposal_vote(pip, vote_option=None):
    proposalinfo = pip.get_effect_proposal_info_of_vote()
    log.info('Version proposalinfo: {}'.format(proposalinfo))
    if not proposalinfo:
        raise Exception('there is no voting version proposal')
    if proposalinfo.get('NewVersion') == pip.cfg.version5:
        upload_platon(pip.node, pip.cfg.PLATON_NEW_BIN)
        log.info('Replace the node {} version to {}'.format(pip.node.node_id, pip.cfg.version5))
    elif proposalinfo.get('NewVersion') == pip.cfg.version8:
        upload_platon(pip.node, pip.cfg.PLATON_NEW_BIN8)
        log.info('Replace the node {} version to {}'.format(pip.node.node_id, pip.cfg.version8))
    elif proposalinfo.get('NewVersion') == pip.cfg.version9:
        upload_platon(pip.node, pip.cfg.PLATON_NEW_BIN9)
        log.info('Replace the node {} version to {}'.format(pip.node.node_id, pip.cfg.version9))
    else:
        raise Exception('The new version of the proposal is{}'.format(proposalinfo.get('NewVersion')))
    pip.node.restart()
    log.info('Restart the node {}'.format(pip.node.node_id))
    if not vote_option:
        vote_option = pip.cfg.vote_option_yeas
    result = pip.vote(pip.node.node_id, proposalinfo.get('ProposalID'), vote_option,
                          pip.node.staking_address, transaction_cfg=pip.cfg.transaction_cfg)
    log.info('The node {} vote result {}'.format(pip.node.node_id, result))
    return result


def proposal_vote(pip, vote_option=None, proposaltype=3):
    if vote_option is None:
        vote_option = pip.cfg.vote_option_yeas
    proposalinfo = pip.get_effect_proposal_info_of_vote(proposaltype)
    log.info('proposalinfo: {}'.format(proposalinfo))
    result = pip.vote(pip.node.node_id, proposalinfo.get('ProposalID'), vote_option,
                          pip.node.staking_address, transaction_cfg=pip.cfg.transaction_cfg)
    log.info('Node {} vote param proposal result {}'.format(pip.node.node_id, result))
    return result


@pytest.fixture()
def no_vp_proposal(global_test_env, client_verifier_obj):
    pip = client_verifier_obj.pip
    if pip.is_exist_effective_proposal() or pip.chain_version != pip.cfg.version0 \
            or pip.is_exist_effective_proposal(pip.cfg.param_proposal):
        log.info('There is effective proposal,restart the chain')
        global_test_env.deploy_all()
    return pip


@pytest.fixture()
def submit_version(no_vp_proposal):
    pip = no_vp_proposal
    result = pip.submitVersion(pip.node.node_id, str(time.time()), pip.cfg.version5, 10,
                                   pip.node.staking_address,
                                   transaction_cfg=pip.cfg.transaction_cfg)
    log.info('submit version result : {}'.format(result))
    assert_code(result, 0)
    return pip


@pytest.fixture()
def submit_param(no_vp_proposal, client_list_obj):
    pip = no_vp_proposal
    client_obj = get_client_obj(pip.node.node_id, client_list_obj)
    newvalue = '1'
    if int(get_governable_parameter_value(client_obj, 'slashBlocksReward')) == 1:
        newvalue = '2'
    result = pip.submitParam(pip.node.node_id, str(time.time()), 'slashing', 'slashBlocksReward', newvalue,
                                 pip.node.staking_address, transaction_cfg=pip.cfg.transaction_cfg)
    log.info('submit param proposal result : {}'.format(result))
    assert_code(result, 0)
    return pip


@pytest.fixture()
def submit_cancel(submit_version):
    pip = submit_version
    propolsalinfo = pip.get_effect_proposal_info_of_vote()
    log.info('get voting version proposal info :{}'.format(propolsalinfo))
    result = pip.submitCancel(pip.node.node_id, str(time.time()), 2, propolsalinfo.get('ProposalID'),
                                  pip.node.staking_address, transaction_cfg=pip.cfg.transaction_cfg)
    log.info('submit cancel proposal result : {}'.format(result))
    assert_code(result, 0)
    return pip


@pytest.fixture()
def submit_cancel_param(submit_param):
    pip = submit_param
    propolsalinfo = pip.get_effect_proposal_info_of_vote(pip.cfg.param_proposal)
    log.info('Get voting param proposal info :{}'.format(propolsalinfo))
    result = pip.submitCancel(pip.node.node_id, str(time.time()), 2, propolsalinfo.get('ProposalID'),
                                  pip.node.staking_address, transaction_cfg=pip.cfg.transaction_cfg)
    log.info('submit cancel proposal result : {}'.format(result))
    assert_code(result, 0)
    return pip


@pytest.fixture()
def submit_text(client_verifier_obj):
    pip = client_verifier_obj.pip
    result = pip.submitText(pip.node.node_id, str(time.time()), pip.node.staking_address,
                                transaction_cfg=pip.cfg.transaction_cfg)
    log.info('submit text result:'.format(result))
    assert_code(result, 0)
    return pip


@pytest.fixture()
def new_node_has_proposal(client_new_node_obj, client_verifier_obj, client_noconsensus_obj):
    pip = client_verifier_obj.pip
    if pip.chain_version != pip.cfg.version0 or pip.is_exist_effective_proposal(pip.cfg.param_proposal):
        client_new_node_obj.economic.env.deploy_all()
    if pip.is_exist_effective_proposal():
        proposalinfo = pip.get_effect_proposal_info_of_vote()
        log.info('Get version proposal information {}'.format(proposalinfo))
        if proposalinfo.get('EndVotingBlock') - pip.node.block_number < 2 * pip.economic.consensus_size:
            client_new_node_obj.economic.env.deploy_all()
            result = pip.submitVersion(pip.node.node_id, str(time.time()), pip.cfg.version5, 5,
                                           pip.node.staking_address,
                                           transaction_cfg=pip.cfg.transaction_cfg)
            assert_code(result, 0)
            return client_noconsensus_obj.pip
        else:
            return client_new_node_obj.pip
    result = pip.submitVersion(pip.node.node_id, str(time.time()), pip.cfg.version5, 5,
                                   pip.node.staking_address, transaction_cfg=pip.cfg.transaction_cfg)
    assert_code(result, 0)
    return client_new_node_obj.pip


@pytest.fixture()
def candidate_has_proposal(client_candidate_obj, client_list_obj):
    verifier_list = get_pledge_list(client_candidate_obj.ppos.getVerifierList)
    pip = get_client_obj(verifier_list[0], client_list_obj).pip
    if pip.chain_version != pip.cfg.version0:
        pip.economic.env.deploy_all()
    if pip.is_exist_effective_proposal_for_vote(pip.cfg.param_proposal):
        pip.economic.env.deploy_all()
    if pip.is_exist_effective_proposal():
        proposalinfo = pip.get_effect_proposal_info_of_vote()
        log.info('Get version proposal information : {}'.format(proposalinfo))
        if proposalinfo.get('EndVotingBlock') - pip.node.block_number < 2 * pip.economic.consensus_size:
            pip.economic.env.deploy_all()
            normal_node_obj_list = pip.economic.env.normal_node_list
            for normal_node_obj in normal_node_obj_list:
                client_obj = get_client_obj(normal_node_obj.node_id, client_list_obj)
                address, _ = client_obj.economic.account.generate_account(client_obj.node.web3, 10 ** 18 * 10000000)
                log.info('Node {} staking'.format(normal_node_obj.node_id))
                result = client_obj.staking.create_staking(0, address, address)
                log.info('Node {} staking result: {}'.format(normal_node_obj.node_id, result))
                assert_code(result, 0)
            pip.economic.wait_settlement_blocknum(pip.node)
            node_id_list = pip.get_candidate_list_not_verifier()
            if not node_id_list:
                raise Exception('Get candidate list')
            client_candidate_obj = get_client_obj(node_id_list[0], client_list_obj)
        else:
            return client_candidate_obj.pip
    result = pip.submitVersion(pip.node.node_id, str(time.time()), pip.cfg.version5, 5,
                                   pip.node.staking_address, transaction_cfg=pip.cfg.transaction_cfg)
    log.info('Submit version proposal result : {}'.format(result))
    assert_code(result, 0)
    return client_candidate_obj.pip


@pytest.fixture()
def noproposal_pips(client_list_obj) -> List[Pip]:
    '''
    Get candidate Client object list
    :param global_test_env:
    :return:
    '''
    if client_list_obj[0].pip.is_exist_effective_proposal() or client_list_obj[0].pip.chain_version != \
            client_list_obj[0].pip.cfg.version0:
        log.info('There is effective proposal, Restart the chain')
        client_list_obj[0].economic.env.deploy_all()
    verifier_list = get_pledge_list(client_list_obj[0].ppos.getVerifierList)
    log.info('verifierlist{}'.format(verifier_list))
    clients = get_client_obj_list(verifier_list, client_list_obj)
    return [client_obj.pip for client_obj in clients]


@pytest.fixture()
def noproposal_candidate_pips(client_list_obj) -> List[Pip]:
    '''
    Get verifier Client object list
    :param global_test_env:
    :return:
    '''
    if client_list_obj[0].pip.is_exist_effective_proposal() or client_list_obj[0].pip.chain_version != \
            client_list_obj[0].pip.cfg.version0:
        log.info('There is effective proposal, Restart the chain')
        client_list_obj[0].economic.env.deploy_all()
    nodeid_list = client_list_obj[0].pip.get_candidate_list_not_verifier()
    log.info('candidate not verifier list {}'.format(nodeid_list))
    if not nodeid_list:
        candidate_list = get_pledge_list(client_list_obj[0].ppos.getCandidateList)
        log.info('candidate_list{}'.format(candidate_list))
        for client_obj in client_list_obj:
            if client_obj.node.node_id not in candidate_list:
                address, _ = client_obj.economic.account.generate_account(client_obj.node.web3, 10**18 * 10000000)
                result = client_obj.staking.create_staking(0, address, address)
                log.info('node {} staking result {}'.format(client_obj.node.node_id, result))
        client_obj.economic.wait_settlement_blocknum(client_obj.node)
        nodeid_list = client_list_obj[0].pip.get_candidate_list_not_verifier()
        if not nodeid_list:
            raise Exception('get candidate not verifier failed')
    clients_candidate = get_client_obj_list(nodeid_list, client_list_obj)
    return [client_candidate.pip for client_candidate in clients_candidate]


@pytest.fixture()
def proposal_pips(client_list_obj):
    '''
    get verifier Client object list
    :param global_test_env:
    :return:
    '''
    verifier_list = get_pledge_list(client_list_obj[0].ppos.getVerifierList)
    pip = get_client_obj(verifier_list[0], client_list_obj).pip
    if pip.chain_version != pip.cfg.version0:
        pip.economic.env.deploy_all()
    if pip.is_exist_effective_proposal():
        proposalinfo = pip.get_effect_proposal_info_of_vote()
        log.info('proprosalinfo : {}'.format(proposalinfo))
        if proposalinfo.get('EndVotingBlock') - pip.node.block_number > 2 * pip.economic.consensus_size \
                and proposalinfo.get('NewVersion') == pip.cfg.version5:
            verifier_list = get_pledge_list(client_list_obj[0].ppos.getVerifierList)
            log.info('verifierlist{}'.format(verifier_list))
            client_obj_list = get_client_obj_list(verifier_list, client_list_obj)
            return [client_obj.pip for client_obj in client_obj_list]
        else:
            pip.economic.env.deploy_all()
    result = pip.submitVersion(pip.node.node_id, str(time.time_ns()), pip.cfg.version5, 10,
                                   pip.node.staking_address, transaction_cfg=pip.cfg.transaction_cfg)
    log.info('version proposal result :{}'.format(result))
    verifier_list = get_pledge_list(client_list_obj[0].ppos.getVerifierList)
    log.info('verifierlist{}'.format(verifier_list))
    client_obj_list = get_client_obj_list(verifier_list, client_list_obj)
    return [client_obj.pip for client_obj in client_obj_list]

@pytest.fixture()
def preactive_proposal_pips(client_list_obj):
    if client_list_obj[0].pip.is_exist_effective_proposal() or client_list_obj[0].pip.chain_version != \
            client_list_obj[0].pip.cfg.version0:
        log.info('There is effective version proposal, restart the chain')
        client_list_obj[0].economic.env.deploy_all()
    verifier_list = get_pledge_list(client_list_obj[0].ppos.getVerifierList)
    log.info('verifierlist :{}'.format(verifier_list))
    client_verifiers = get_client_obj_list(verifier_list, client_list_obj)
    pips = [client_verifier.pip for client_verifier in client_verifiers]
    result = pips[0].submitVersion(pips[0].node.node_id, str(time.time_ns()),
                                           pips[0].cfg.version5, 2, pips[0].node.staking_address,
                                           transaction_cfg=pips[0].cfg.transaction_cfg)
    log.info('submit version proposal, result : {}'.format(result))
    proposalinfo = pips[0].get_effect_proposal_info_of_vote()
    log.info('Version proposalinfo: {}'.format(proposalinfo))
    for pip in pips:
        result = version_proposal_vote(pip)
        assert_code(result, 0)
    wait_block_number(pip.node, proposalinfo.get('EndVotingBlock'))
    assert pip.get_status_of_proposal(proposalinfo.get('ProposalID')) == 4
    return pips


@pytest.fixture()
def preactive_large_version_proposal_pips(client_list_obj):
    if client_list_obj[0].pip.is_exist_effective_proposal() or client_list_obj[0].pip.chain_version != \
            client_list_obj[0].pip.cfg.version0:
        log.info('There is effective version proposal, restart the chain')
        client_list_obj[0].economic.env.deploy_all()
    verifier_list = get_pledge_list(client_list_obj[0].ppos.getVerifierList)
    log.info('verifierlist :{}'.format(verifier_list))
    client_verifiers = get_client_obj_list(verifier_list, client_list_obj)
    pips = [client_obj.pip for client_obj in client_verifiers]
    result = pips[0].submitVersion(pips[0].node.node_id, str(time.time_ns()),
                                           pips[0].cfg.version8, 2, pips[0].node.staking_address,
                                           transaction_cfg=pips[0].cfg.transaction_cfg)
    log.info('submit version proposal, result : {}'.format(result))
    proposalinfo = pips[0].get_effect_proposal_info_of_vote()
    log.info('Version proposalinfo: {}'.format(proposalinfo))
    for pip in pips:
        result = version_proposal_vote(pip)
        assert_code(result, 0)
    wait_block_number(pip.node, proposalinfo.get('EndVotingBlock'))
    assert pip.get_status_of_proposal(proposalinfo.get('ProposalID')) == 4
    return pips
