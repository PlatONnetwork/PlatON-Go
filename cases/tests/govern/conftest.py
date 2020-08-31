import pytest
from common.log import log
import time
import math
from tests.lib.client import get_client_by_nodeid, get_clients_by_nodeid, Client
from tests.conftest import get_clients
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


def verifier_node_version(obj, version=None):
    if not isinstance(obj, Client):
        obj = get_client_by_nodeid(obj.node.node_id, get_clients(obj.economic.env))
    node_version = obj.staking.get_version()
    log.info('Node {} version is {}'.format(obj.node.node_id, node_version))
    if version is None:
        return node_version
    else:
        assert_code(node_version, version)


@pytest.fixture()
def no_vp_proposal(global_test_env, client_verifier):
    pip = client_verifier.pip
    if pip.is_exist_effective_proposal() or pip.chain_version != pip.cfg.version0 \
            or pip.is_exist_effective_proposal_for_vote(pip.cfg.param_proposal):
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
def submit_param(no_vp_proposal, all_clients):
    pip = no_vp_proposal
    client = get_client_by_nodeid(pip.node.node_id, all_clients)
    newvalue = '1'
    if int(get_governable_parameter_value(client, 'slashBlocksReward')) == 1:
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
def submit_text(client_verifier):
    pip = client_verifier.pip
    result = pip.submitText(pip.node.node_id, str(time.time()), pip.node.staking_address,
                            transaction_cfg=pip.cfg.transaction_cfg)
    log.info('submit text result:'.format(result))
    assert_code(result, 0)
    return pip


@pytest.fixture()
def new_node_has_proposal(client_new_node, client_verifier, client_noconsensus):
    pip = client_verifier.pip
    if pip.chain_version != pip.cfg.version0 or pip.is_exist_effective_proposal(pip.cfg.param_proposal):
        client_new_node.economic.env.deploy_all()
    if pip.is_exist_effective_proposal():
        proposalinfo = pip.get_effect_proposal_info_of_vote()
        log.info('Get version proposal information {}'.format(proposalinfo))
        if proposalinfo.get('EndVotingBlock') - pip.node.block_number < 2 * pip.economic.consensus_size:
            client_new_node.economic.env.deploy_all()
            result = pip.submitVersion(pip.node.node_id, str(time.time()), pip.cfg.version5, 5,
                                       pip.node.staking_address,
                                       transaction_cfg=pip.cfg.transaction_cfg)
            assert_code(result, 0)
            return client_noconsensus.pip
        else:
            return client_new_node.pip
    result = pip.submitVersion(pip.node.node_id, str(time.time()), pip.cfg.version5, 5,
                               pip.node.staking_address, transaction_cfg=pip.cfg.transaction_cfg)
    assert_code(result, 0)
    return client_new_node.pip


@pytest.fixture()
def candidate_has_proposal(clients_noconsensus, all_clients):
    clients_noconsensus[0].economic.env.deploy_all()
    for client in clients_noconsensus:
        address, _ = client.economic.account.generate_account(client.node.web3, 10 ** 18 * 10000000)
        log.info('Node {} staking'.format(client.node.node_id))
        result = client.staking.create_staking(0, address, address)
        log.info('Node {} staking result: {}'.format(client.node.node_id, result))
        assert_code(result, 0)
    client.economic.wait_settlement_blocknum(client.node)
    node_id_list = client.pip.get_candidate_list_not_verifier()
    if not node_id_list:
        raise Exception('Get candidate list')
    verifiers = get_pledge_list(client.ppos.getVerifierList)
    log.info('Verifier list : {}'.format(verifiers))
    pip = get_client_by_nodeid(verifiers[0], all_clients).pip
    result = pip.submitVersion(pip.node.node_id, str(time.time()), pip.cfg.version5, 5,
                               pip.node.staking_address, transaction_cfg=pip.cfg.transaction_cfg)
    log.info('Submit version proposal result : {}'.format(result))
    assert_code(result, 0)
    return get_client_by_nodeid(node_id_list[0], all_clients).pip


@pytest.fixture()
def noproposal_pips(all_clients) -> List[Pip]:
    '''
    Get candidate Client object list
    :param global_test_env:
    :return:
    '''
    if all_clients[0].pip.is_exist_effective_proposal() or all_clients[0].pip.chain_version != \
            all_clients[0].pip.cfg.version0:
        log.info('There is effective proposal, Restart the chain')
        all_clients[0].economic.env.deploy_all()
    verifier_list = get_pledge_list(all_clients[0].ppos.getVerifierList)
    log.info('verifierlist{}'.format(verifier_list))
    clients = get_clients_by_nodeid(verifier_list, all_clients)
    return [client.pip for client in clients]


@pytest.fixture()
def noproposal_candidate_pips(all_clients) -> List[Pip]:
    '''
    Get verifier Client object list
    :param global_test_env:
    :return:
    '''
    if all_clients[0].pip.is_exist_effective_proposal() or all_clients[0].pip.chain_version != \
            all_clients[0].pip.cfg.version0:
        log.info('There is effective proposal, Restart the chain')
        all_clients[0].economic.env.deploy_all()
    nodeid_list = all_clients[0].pip.get_candidate_list_not_verifier()
    log.info('candidate not verifier list {}'.format(nodeid_list))
    if not nodeid_list:
        candidate_list = get_pledge_list(all_clients[0].ppos.getCandidateList)
        log.info('candidate_list{}'.format(candidate_list))
        for client in all_clients:
            if client.node.node_id not in candidate_list:
                address, _ = client.economic.account.generate_account(client.node.web3, 10**18 * 10000000)
                result = client.staking.create_staking(0, address, address)
                log.info('node {} staking result {}'.format(client.node.node_id, result))
        client.economic.wait_settlement_blocknum(client.node)
        nodeid_list = all_clients[0].pip.get_candidate_list_not_verifier()
        if not nodeid_list:
            raise Exception('get candidate not verifier failed')
    clients_candidate = get_clients_by_nodeid(nodeid_list, all_clients)
    return [client_candidate.pip for client_candidate in clients_candidate]


@pytest.fixture()
def proposal_pips(all_clients):
    '''
    get verifier Client object list
    :param global_test_env:
    :return:
    '''
    verifier_list = get_pledge_list(all_clients[0].ppos.getVerifierList)
    pip = get_client_by_nodeid(verifier_list[0], all_clients).pip
    if pip.chain_version != pip.cfg.version0:
        pip.economic.env.deploy_all()
    if pip.is_exist_effective_proposal():
        proposalinfo = pip.get_effect_proposal_info_of_vote()
        log.info('proprosalinfo : {}'.format(proposalinfo))
        if proposalinfo.get('EndVotingBlock') - pip.node.block_number > 2 * pip.economic.consensus_size \
                and proposalinfo.get('NewVersion') == pip.cfg.version5:
            verifier_list = get_pledge_list(all_clients[0].ppos.getVerifierList)
            log.info('verifierlist{}'.format(verifier_list))
            clients_verifier = get_clients_by_nodeid(verifier_list, all_clients)
            return [client.pip for client in clients_verifier]
        else:
            pip.economic.env.deploy_all()
    result = pip.submitVersion(pip.node.node_id, str(time.time()), pip.cfg.version5, 10,
                               pip.node.staking_address, transaction_cfg=pip.cfg.transaction_cfg)
    log.info('version proposal result :{}'.format(result))
    verifier_list = get_pledge_list(all_clients[0].ppos.getVerifierList)
    log.info('verifierlist{}'.format(verifier_list))
    clients_verifier = get_clients_by_nodeid(verifier_list, all_clients)
    return [client.pip for client in clients_verifier]


@pytest.fixture()
def preactive_proposal_pips(all_clients):
    if all_clients[0].pip.is_exist_effective_proposal() or all_clients[0].pip.chain_version != \
            all_clients[0].pip.cfg.version0 or all_clients[0].pip.is_exist_effective_proposal_for_vote(
        all_clients[0].pip.cfg.param_proposal
    ):
        log.info('There is effective version proposal, restart the chain')
        all_clients[0].economic.env.deploy_all()
    verifier_list = get_pledge_list(all_clients[0].ppos.getVerifierList)
    log.info('verifierlist :{}'.format(verifier_list))
    client_verifiers = get_clients_by_nodeid(verifier_list, all_clients)
    pips = [client_verifier.pip for client_verifier in client_verifiers]
    result = pips[0].submitVersion(pips[0].node.node_id, str(time.time()),
                                   pips[0].cfg.version5, 3, pips[0].node.staking_address,
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
def preactive_large_version_proposal_pips(all_clients):
    if all_clients[0].pip.is_exist_effective_proposal() or all_clients[0].pip.chain_version != \
            all_clients[0].pip.cfg.version0:
        log.info('There is effective version proposal, restart the chain')
        all_clients[0].economic.env.deploy_all()
    verifier_list = get_pledge_list(all_clients[0].ppos.getVerifierList)
    log.info('verifierlist :{}'.format(verifier_list))
    client_verifiers = get_clients_by_nodeid(verifier_list, all_clients)
    pips = [client.pip for client in client_verifiers]
    result = pips[0].submitVersion(pips[0].node.node_id, str(time.time()),
                                   pips[0].cfg.version8, 3, pips[0].node.staking_address,
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
