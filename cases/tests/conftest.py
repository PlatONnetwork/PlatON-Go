import time

import pytest
from copy import copy
from tests.lib import StakingConfig
from common.log import log
from tests.lib.client import Client, get_client_obj, get_client_obj_list
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


def get_client_list_obj(env, cfg):
    client_list_obj = []
    node_obj_list = env.get_all_nodes()
    for node_obj in node_obj_list:
        client_list_obj.append(Client(env, node_obj, cfg))
    return client_list_obj


@pytest.fixture()
def client_list_obj(global_running_env, staking_cfg):
    """
    Get all node  Node object list
    """
    return get_client_list_obj(global_running_env, staking_cfg)


def get_con_list_list(env, cfg):
    client_con_list_obj = []
    consensus_node_obj_list = env.consensus_node_list
    for node_obj in consensus_node_obj_list:
        client_con_list_obj.append(Client(env, node_obj, cfg))
    return client_con_list_obj


@pytest.fixture()
def client_con_list_obj(global_running_env, staking_cfg):
    """
    Get all consensus node  Client object list
    """
    return get_con_list_list(global_running_env, staking_cfg)


def get_client_noconsensus_list(env, cfg):
    client_noc_list_obj = []
    noconsensus_node_obj_list = env.normal_node_list
    for node_obj in noconsensus_node_obj_list:
        client_noc_list_obj.append(Client(env, node_obj, cfg))
    return client_noc_list_obj


@pytest.fixture()
def client_noc_list_obj(global_running_env, staking_cfg):
    """
    Get all noconsensus node  Client object list
    """
    return get_client_noconsensus_list(global_running_env, staking_cfg)


def get_client_consensus_obj(env, cfg):
    consensus_node_obj = env.get_rand_node()
    client_consensus_obj = Client(env, consensus_node_obj, cfg)
    return client_consensus_obj


@pytest.fixture()
def client_consensus_obj(global_running_env, staking_cfg):
    """
    Get a consensus node  Client object
    """
    return get_client_consensus_obj(global_running_env, staking_cfg)


@pytest.fixture()
def client_noconsensus_obj(global_running_env, staking_cfg):
    """
    Get a noconsensus node  Client object
    """
    noconsensus_node_obj = global_running_env.get_a_normal_node()
    client_noconsensus_obj = Client(global_running_env, noconsensus_node_obj, staking_cfg)
    return client_noconsensus_obj


@pytest.fixture()
def client_verifier_obj(global_running_env, staking_cfg):
    """
    Get a verifier node  Client object
    """
    client_list_obj = get_client_list_obj(global_running_env, staking_cfg)
    verifier_list = get_pledge_list(client_list_obj[0].ppos.getVerifierList)
    log.info('verifierlist{}'.format(verifier_list))
    for client in client_list_obj:
        if client.node.node_id in verifier_list:
            return client
    raise Exception('Get a verifier node  Client object ')


@pytest.fixture()
def client_verifier_obj_list(global_running_env, staking_cfg):
    """
    Get verifier node  Client object list
    """
    client_list_obj = get_client_list_obj(global_running_env, staking_cfg)
    verifier_list = get_pledge_list(client_list_obj[0].ppos.getVerifierList)
    log.info('verifierlist{}'.format(verifier_list))
    return get_client_obj_list(verifier_list, client_list_obj)


@pytest.fixture()
def client_new_node_obj(global_running_env, staking_cfg):
    """
    Get a new node  Client object list
    """
    normal_node = global_running_env.get_a_normal_node()
    for noconsensus_node in global_running_env.normal_node_list:
        msg = noconsensus_node.ppos.getCandidateInfo(noconsensus_node.node_id)
        log.info(noconsensus_node.node_id)
        if msg["Code"] == 301204:
            log.info("Current linked node: {}".format(noconsensus_node.node_mark))
            return Client(global_running_env, noconsensus_node, staking_cfg)
    log.info('noconsensus node has been staked, restart the chain')
    global_running_env.deploy_all()
    log.info("Current linked node: {}".format(normal_node.node_mark))
    return Client(global_running_env, normal_node, staking_cfg)


@pytest.fixture()
def client_new_node_obj_list(global_test_env, staking_cfg):
    """
    Get new node Client object list
    """
    global_test_env.deploy_all()
    return get_client_noconsensus_list(global_test_env, staking_cfg)


@pytest.fixture()
def client_candidate_obj(global_running_env, staking_cfg):
    """
    Get a candidate node Client object
    """
    client_consensus_obj = get_client_consensus_obj(global_running_env, staking_cfg)
    client_list_obj = get_client_list_obj(global_running_env, staking_cfg)
    client_noconsensus_obj_list = get_client_noconsensus_list(global_running_env, staking_cfg)
    if not client_consensus_obj.staking.get_candidate_list_not_verifier():
        log.info('There is no candidate, node stake')
        candidate_list = get_pledge_list(client_consensus_obj.node.ppos.getCandidateList)
        for client in client_noconsensus_obj_list:
            if client.node.node_id not in candidate_list:
                if client.node.program_version != client.pip.cfg.version0:
                    upload_platon(client.node, client.pip.cfg.PLATON_NEW_BIN0)
                    client.node.restart()
                log.info('Node {} staking'.format(client.node.node_id))
                address, _ = client.economic.account.generate_account(client.node.web3, client.economic.create_staking_limit * 5)
                result = client.staking.create_staking(0, address, address)
                log.info('Node {} staking result :{}'.format(client.node.node_id, result))
                assert_code(result, 0)
        client_consensus_obj.economic.wait_settlement_blocknum(client_consensus_obj.node)
    node_id_list = client_consensus_obj.staking.get_candidate_list_not_verifier()
    log.info('Get candidate list no verifier {}'.format(node_id_list))
    if len(node_id_list) == 0:
        raise Exception('Get candidate list no verifier failed')
    return get_client_obj(node_id_list[0], client_list_obj)


@pytest.fixture()
def reset_environment(global_test_env):
    log.info("case execution completed")
    yield
    global_test_env.deploy_all()


@pytest.fixture
def new_genesis_env(global_test_env):
    cfg = copy(global_test_env.cfg)
    yield global_test_env
    log.info("reset deploy.................")
    global_test_env.set_cfg(cfg)
    global_test_env.deploy_all()


def param_governance_verify(client_obj, module, name, newvalue, effectiveflag=True):
    """
    effectiveflag indicates whether it takes effect
    """
    if isinstance(client_obj, Client):
        pip_obj = client_obj.pip
    else:
        raise Exception("client must Client class")
    if pip_obj.is_exist_effective_proposal_for_vote(pip_obj.cfg.param_proposal) or \
            pip_obj.is_exist_effective_proposal_for_vote(pip_obj.cfg.version_proposal):
        raise Exception('There is effective param proposal or version proposal')
    result = pip_obj.submitParam(pip_obj.node.node_id, str(time.time()), module, name, newvalue, pip_obj.node.staking_address,
                                 transaction_cfg=pip_obj.cfg.transaction_cfg)
    log.info('submit param proposal result : {}'.format(result))
    assert_code(result, 0)
    proposalinfo = pip_obj.get_effect_proposal_info_of_vote(pip_obj.cfg.param_proposal)
    log.info('param proposalinfo : {}'.format(proposalinfo))
    client_obj_list = []
    for node_obj in pip_obj.economic.env.get_all_nodes():
        client_obj_list.append(Client(pip_obj.economic.env, node_obj,
                                      StakingConfig("externalId", "nodeName", "website", "details")))
    client_obj = get_client_obj(pip_obj.node.node_id, client_obj_list)
    verifier_list = get_pledge_list(client_obj.ppos.getVerifierList)
    log.info('verifierlist : {}'.format(verifier_list))
    client_verifier_obj_list = get_client_obj_list(verifier_list, client_obj_list)
    if effectiveflag:
        for client_obj in client_verifier_obj_list:
            result = client_obj.pip.vote(client_obj.node.node_id, proposalinfo.get('ProposalID'),
                                         client_obj.pip.cfg.vote_option_yeas,
                                         client_obj.node.staking_address, transaction_cfg=client_obj.pip.cfg.transaction_cfg)
            log.info('Node {} vote proposal result : {}'.format(client_obj.node.node_id, result))
    wait_block_number(pip_obj.node, proposalinfo.get('EndVotingBlock'))
    if effectiveflag:
        assert pip_obj.get_status_of_proposal(proposalinfo.get('ProposalID')) == 2
        log.info("blockNumber {}, the {} has become {}".format(proposalinfo.get('EndVotingBlock'), name, newvalue))
    else:
        assert pip_obj.get_status_of_proposal(proposalinfo.get('ProposalID')) == 3
        log.info("{} retains the original value".format(name))


def param_governance_verify_before_endblock(client_obj, module, name, newvalue, effectiveflag=True):
    """
    effectiveflag indicates whether it takes effect
    :param client_obj:
    :param module:
    :param name:
    :param newvalue:
    :param effectiveflag:
    :return: the EndVotingBlock of the param proposal
    """
    if isinstance(client_obj, Client):
        pip_obj = client_obj.pip
    else:
        raise Exception("client must Client class")
    if pip_obj.is_exist_effective_proposal_for_vote(pip_obj.cfg.param_proposal) or \
            pip_obj.is_exist_effective_proposal_for_vote(pip_obj.cfg.version_proposal):
        raise Exception('There is effective param proposal or version proposal')
    result = pip_obj.submitParam(pip_obj.node.node_id, str(time.time()), module, name, newvalue, pip_obj.node.staking_address,
                                 transaction_cfg=pip_obj.cfg.transaction_cfg)
    log.info('submit param proposal result : {}'.format(result))
    assert_code(result, 0)
    proposalinfo = pip_obj.get_effect_proposal_info_of_vote(pip_obj.cfg.param_proposal)
    log.info('param proposalinfo : {}'.format(proposalinfo))
    client_obj_list = []
    for node_obj in pip_obj.economic.env.get_all_nodes():
        client_obj_list.append(Client(pip_obj.economic.env, node_obj,
                                      StakingConfig("externalId", "nodeName", "website", "details")))
    client_obj = get_client_obj(pip_obj.node.node_id, client_obj_list)
    verifier_list = get_pledge_list(client_obj.ppos.getVerifierList)
    log.info('verifierlist : {}'.format(verifier_list))
    client_verifier_obj_list = get_client_obj_list(verifier_list, client_obj_list)
    if effectiveflag:
        for client_obj in client_verifier_obj_list:
            result = client_obj.pip.vote(client_obj.node.node_id, proposalinfo.get('ProposalID'),
                                         client_obj.pip.cfg.vote_option_yeas,
                                         client_obj.node.staking_address, transaction_cfg=client_obj.pip.cfg.transaction_cfg)
            log.info('Node {} vote proposal result : {}'.format(client_obj.node.node_id, result))
    log.info('The proposal endvoting block is {}'.format(proposalinfo.get('EndVotingBlock')))
    return proposalinfo.get('EndVotingBlock')
