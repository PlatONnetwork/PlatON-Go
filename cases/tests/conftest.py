import time
from typing import List
import pytest
from copy import copy
from tests.lib import StakingConfig
from common.log import log
from tests.lib.client import Client, get_client_by_nodeid, get_clients_by_nodeid
from tests.lib.utils import get_pledge_list, wait_block_number, assert_code, upload_platon


@pytest.fixture()
def global_running_env(global_test_env):
    cfg = global_test_env.cfg
    genesis = global_test_env.genesis_config
    backup_cfg = copy(cfg)
    id_cfg = id(cfg)
    if not global_test_env.running:
        log.info("The environment is not running, redeploying the environment")
        global_test_env.deploy_all()
    yield global_test_env
    if id_cfg != id(global_test_env.cfg) or id(genesis) != id(global_test_env.genesis_config):
        log.info("Environment configuration changes, restore configuration files and redeploy")
        global_test_env.set_cfg(backup_cfg)
        global_test_env.deploy_all()


@pytest.fixture()
def staking_cfg():
    cfg = StakingConfig("externalId", "nodeName", "website", "details")
    return cfg


def get_clients(env, cfg=None):
    if cfg is None:
        cfg = StakingConfig("externalId", "nodeName", "website", "details")
    all_clients = []
    all_nodes = env.get_all_nodes()
    for node_obj in all_nodes:
        all_clients.append(Client(env, node_obj, cfg))
    return all_clients


@pytest.fixture()
def all_clients(global_running_env, staking_cfg) -> List[Client]:
    """
    Get all node  Node object list
    """
    return get_clients(global_running_env, staking_cfg)


def get_consensus_clients(env, cfg):
    clients_consensus = []
    consensus_nodes = env.consensus_node_list
    for node in consensus_nodes:
        clients_consensus.append(Client(env, node, cfg))
    return clients_consensus


@pytest.fixture()
def clients_consensus(global_running_env, staking_cfg) -> List[Client]:
    """
    Get all consensus node  Client object list
    """
    return get_consensus_clients(global_running_env, staking_cfg)


def get_clients_noconsensus(env, cfg):
    client_noconsensus = []
    noconsensus_nodes = env.normal_node_list
    for node_obj in noconsensus_nodes:
        client_noconsensus.append(Client(env, node_obj, cfg))
    return client_noconsensus


@pytest.fixture()
def clients_noconsensus(global_running_env, staking_cfg) -> List[Client]:
    """
    Get all noconsensus node  Client object list
    """
    return get_clients_noconsensus(global_running_env, staking_cfg)


def get_client_consensus(env, cfg):
    consensus_node = env.get_rand_node()
    client_consensus = Client(env, consensus_node, cfg)
    return client_consensus


@pytest.fixture()
def client_consensus(global_running_env, staking_cfg) -> Client:
    """
    Get a consensus node  Client object
    """
    return get_client_consensus(global_running_env, staking_cfg)


@pytest.fixture()
def client_noconsensus(global_running_env, staking_cfg) -> Client:
    """
    Get a noconsensus node  Client object
    """
    noconsensus_node = global_running_env.get_a_normal_node()
    client_noconsensus = Client(global_running_env, noconsensus_node, staking_cfg)
    return client_noconsensus


@pytest.fixture()
def client_verifier(global_running_env, staking_cfg) -> Client:
    """
    Get a verifier node  Client object
    """
    all_clients = get_clients(global_running_env, staking_cfg)
    verifier_list = get_pledge_list(all_clients[0].ppos.getVerifierList)
    log.info('verifierlist{}'.format(verifier_list))
    for client in all_clients:
        if client.node.node_id in verifier_list:
            return client
    raise Exception('Get a verifier node  Client object ')


@pytest.fixture()
def clients_verifier(global_running_env, staking_cfg) -> List[Client]:
    """
    Get verifier node  Client object list
    """
    all_clients = get_clients(global_running_env, staking_cfg)
    verifier_list = get_pledge_list(all_clients[0].ppos.getVerifierList)
    log.info('verifierlist{}'.format(verifier_list))
    return get_clients_by_nodeid(verifier_list, all_clients)


@pytest.fixture()
def client_new_node(global_running_env, staking_cfg) -> Client:
    """
    Get a new node  Client object list
    """
    normal_node = global_running_env.get_a_normal_node()
    for noconsensus_node in global_running_env.normal_node_list:
        msg = noconsensus_node.ppos.getCandidateInfo(noconsensus_node.node_id)
        log.info(msg)
        if msg["Code"] == 301204:
            log.info("Current linked node: {}".format(noconsensus_node.node_mark))
            return Client(global_running_env, noconsensus_node, staking_cfg)
    log.info('noconsensus node has been staked, restart the chain')
    global_running_env.deploy_all()
    log.info("Current linked node: {}".format(normal_node.node_mark))
    return Client(global_running_env, normal_node, staking_cfg)


@pytest.fixture()
def clients_new_node(global_test_env, staking_cfg) -> List[Client]:
    """
    Get new node Client object list
    """
    global_test_env.deploy_all()
    return get_clients_noconsensus(global_test_env, staking_cfg)


@pytest.fixture()
def client_candidate(global_running_env, staking_cfg):
    """
    Get a candidate node Client object
    """
    client_consensus = get_client_consensus(global_running_env, staking_cfg)
    all_clients = get_clients(global_running_env, staking_cfg)
    clients_noconsensus = get_clients_noconsensus(global_running_env, staking_cfg)
    if not client_consensus.staking.get_candidate_list_not_verifier():
        log.info('There is no candidate, node stake')
        candidate_list = get_pledge_list(client_consensus.node.ppos.getCandidateList)
        for client in clients_noconsensus:
            if client.node.node_id not in candidate_list:
                if client.node.program_version != client.pip.cfg.version0:
                    upload_platon(client.node, client.pip.cfg.PLATON_NEW_BIN0)
                    client.node.restart()
                log.info('Node {} staking'.format(client.node.node_id))
                address, _ = client.economic.account.generate_account(client.node.web3, client.economic.create_staking_limit * 5)
                result = client.staking.create_staking(0, address, address)
                log.info('Node {} staking result :{}'.format(client.node.node_id, result))
                assert_code(result, 0)
        client_consensus.economic.wait_settlement_blocknum(client_consensus.node)
    node_id_list = client_consensus.staking.get_candidate_list_not_verifier()
    log.info('Get candidate list no verifier {}'.format(node_id_list))
    if len(node_id_list) == 0:
        raise Exception('Get candidate list no verifier failed')
    return get_client_by_nodeid(node_id_list[0], all_clients)


@pytest.fixture()
def reset_environment(global_test_env):
    log.info("case execution completed")
    yield
    global_test_env.deploy_all()


@pytest.fixture()
def new_genesis_env(global_test_env):
    cfg = copy(global_test_env.cfg)
    yield global_test_env
    log.info("reset deploy.................")
    global_test_env.set_cfg(cfg)
    global_test_env.deploy_all()


def param_governance_verify(client, module, name, newvalue, effectiveflag=True):
    """
    effectiveflag indicates whether it takes effect
    """
    if isinstance(client, Client):
        pip = client.pip
    else:
        raise Exception("client must Client class")
    if pip.is_exist_effective_proposal_for_vote(pip.cfg.param_proposal) or \
            pip.is_exist_effective_proposal_for_vote(pip.cfg.version_proposal):
        raise Exception('There is effective param proposal or version proposal')
    result = pip.submitParam(pip.node.node_id, str(time.time()), module, name, newvalue, pip.node.staking_address,
                             transaction_cfg=pip.cfg.transaction_cfg)
    log.info('submit param proposal result : {}'.format(result))
    assert_code(result, 0)
    proposalinfo = pip.get_effect_proposal_info_of_vote(pip.cfg.param_proposal)
    log.info('param proposalinfo : {}'.format(proposalinfo))
    all_clients = []
    for node_obj in pip.economic.env.get_all_nodes():
        all_clients.append(Client(pip.economic.env, node_obj,
                                  StakingConfig("externalId", "nodeName", "website", "details")))
    client = get_client_by_nodeid(pip.node.node_id, all_clients)
    verifier_list = get_pledge_list(client.ppos.getVerifierList)
    log.info('verifierlist : {}'.format(verifier_list))
    clients_verifier = get_clients_by_nodeid(verifier_list, all_clients)
    if effectiveflag:
        blocknum = 0
        for client in clients_verifier:
            if client.node.block_number < blocknum and blocknum != 0:
                wait_block_number(client.node, blocknum)
            result = client.pip.vote(client.node.node_id, proposalinfo.get('ProposalID'),
                                     client.pip.cfg.vote_option_yeas,
                                     client.node.staking_address, transaction_cfg=client.pip.cfg.transaction_cfg)
            log.info('Node {} vote proposal result : {}'.format(client.node.node_id, result))
            blocknum = client.node.block_number
    wait_block_number(pip.node, proposalinfo.get('EndVotingBlock'))
    if effectiveflag:
        assert pip.get_status_of_proposal(proposalinfo.get('ProposalID')) == 2
        log.info("blockNumber {}, the {} has become {}".format(proposalinfo.get('EndVotingBlock'), name, newvalue))
    else:
        assert pip.get_status_of_proposal(proposalinfo.get('ProposalID')) == 3
        log.info("{} retains the original value".format(name))


def param_governance_verify_before_endblock(client, module, name, newvalue, effectiveflag=True):
    """
    effectiveflag indicates whether it takes effect
    :param client_obj:
    :param module:
    :param name:
    :param newvalue:
    :param effectiveflag:
    :return: the EndVotingBlock of the param proposal
    """
    if isinstance(client, Client):
        pip = client.pip
    else:
        raise Exception("client must Client class")
    if pip.is_exist_effective_proposal_for_vote(pip.cfg.param_proposal) or \
            pip.is_exist_effective_proposal_for_vote(pip.cfg.version_proposal):
        raise Exception('There is effective param proposal or version proposal')
    result = pip.submitParam(pip.node.node_id, str(time.time()), module, name, newvalue, pip.node.staking_address,
                             transaction_cfg=pip.cfg.transaction_cfg)
    log.info('submit param proposal result : {}'.format(result))
    assert_code(result, 0)
    proposalinfo = pip.get_effect_proposal_info_of_vote(pip.cfg.param_proposal)
    log.info('param proposalinfo : {}'.format(proposalinfo))
    all_clients = []
    for node in pip.economic.env.get_all_nodes():
        all_clients.append(Client(pip.economic.env, node,
                                  StakingConfig("externalId", "nodeName", "website", "details")))
    client = get_client_by_nodeid(pip.node.node_id, all_clients)
    verifier_list = get_pledge_list(client.ppos.getVerifierList)
    log.info('verifierlist : {}'.format(verifier_list))
    clients_verifier = get_clients_by_nodeid(verifier_list, all_clients)
    if effectiveflag:
        blocknum = 0
        for client in clients_verifier:
            if not client.node.running:
                continue
            if client.node.block_number < blocknum and blocknum != 0:
                wait_block_number(client.node, blocknum)
            result = client.pip.vote(client.node.node_id, proposalinfo.get('ProposalID'),
                                     client.pip.cfg.vote_option_yeas,
                                     client.node.staking_address, transaction_cfg=client.pip.cfg.transaction_cfg)
            log.info('Node {} vote proposal result : {}'.format(client.node.node_id, result))
            blocknum = client.node.block_number
    log.info('The proposal endvoting block is {}'.format(proposalinfo.get('EndVotingBlock')))
    return proposalinfo.get('EndVotingBlock')






