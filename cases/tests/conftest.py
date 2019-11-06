from copy import copy

import pytest
from copy import copy
from tests.lib import StakingConfig
from common.log import log
from tests.lib.client import Client, get_client_obj, get_client_obj_list
from tests.lib.utils import get_pledge_list, wait_block_number

@pytest.fixture()
def global_running_env(global_test_env):
    cfg = global_test_env.cfg
    genesis = global_test_env.genesis_config
    backup_cfg = copy(cfg)
    id_cfg = id(cfg)
    if not global_test_env.running:
        global_test_env.deploy_all()
    yield global_test_env
    if id_cfg != id(global_test_env.cfg) or id(genesis) != id(global_test_env.genesis_config):
        global_test_env.set_cfg(backup_cfg)
        global_test_env.deploy_all()


@pytest.fixture()
def staking_cfg():
    cfg = StakingConfig("externalId", "nodeName", "website", "details")
    return cfg


@pytest.fixture()
def client_list_obj(global_test_env, staking_cfg):
    '''
    获取所有Node对象列表
    :param global_test_env:
    :return:
    '''
    client_list_obj = []
    node_obj_list = global_test_env.get_all_nodes()
    for node_obj in node_obj_list:
        client_list_obj.append(Client(global_test_env, node_obj, staking_cfg))
    return client_list_obj


@pytest.fixture()
def client_con_list_obj(global_test_env, staking_cfg):
    '''
    获取共识Client对象列表
    :param global_test_env:
    :return:
    '''
    client_con_list_obj = []
    consensus_node_obj_list = global_test_env.consensus_node_list
    for node_obj in consensus_node_obj_list:
        client_con_list_obj.append(Client(global_test_env, node_obj, staking_cfg))
    return client_con_list_obj


@pytest.fixture()
def client_noc_list_obj(global_test_env, staking_cfg):
    '''
    获取非共识Client对象列表
    :param global_test_env:
    :return:
    '''
    client_noc_list_obj = []
    noconsensus_node_obj_list = global_test_env.normal_node_list
    for node_obj in noconsensus_node_obj_list:
        client_noc_list_obj.append(Client(global_test_env, node_obj, staking_cfg))
    return client_noc_list_obj


@pytest.fixture()
def client_consensus_obj(global_test_env, staking_cfg):
    '''
    随机获取单个共识Client对象
    :param global_test_env:
    :return:
    '''
    consensus_node_obj = global_test_env.get_rand_node()
    client_consensus_obj = Client(global_test_env, consensus_node_obj, staking_cfg)
    return client_consensus_obj


@pytest.fixture()
def client_noconsensus_obj(global_test_env, staking_cfg):
    '''
    随机获取单个非共识Client对象
    :param global_test_env:
    :return:
    '''
    noconsensus_node_obj = global_test_env.get_a_normal_node()
    client_noconsensus_obj = Client(global_test_env, noconsensus_node_obj, staking_cfg)
    return client_noconsensus_obj


@pytest.fixture()
def client_verifier_obj(global_test_env, client_consensus_obj, client_list_obj):
    '''
    获取单个验证节点Client对象
    :param global_test_env:
    :return:
    '''
    verifier_list = get_pledge_list(client_consensus_obj.ppos.getVerifierList)
    log.info('verifierlist{}'.format(verifier_list))
    nodeid = ""
    for nodeobj in global_test_env.consensus_node_list:
        if nodeobj.node_id in verifier_list:
            nodeid = nodeobj.node_id
            break
    if not nodeid:
        raise Exception('获取验证节点Client对象失败')
    client_obj = get_client_obj(nodeid, client_list_obj)
    return client_obj


@pytest.fixture()
def client_new_node_obj(client_noconsensus_obj, client_list_obj, client_noc_list_obj):
    '''
    获取单个未被质押节点Client对象
    :param global_test_env:
    :return:
    '''
    candidate_list = get_pledge_list(client_noconsensus_obj.ppos.getCandidateList)
    log.info('candidatelist{}'.format(candidate_list))
    for noconsensus_node_obj in client_noc_list_obj:
        if noconsensus_node_obj.node.node_id not in candidate_list:
            return noconsensus_node_obj
    log.info('非共识节点已全部质押，重新启链')
    client_noconsensus_obj.economic.env.deploy_all()
    return client_noconsensus_obj


@pytest.fixture()
def client_new_node_obj_list(global_test_env, client_noc_list_obj):
    '''
    获取新节点Client列表对象
    :param global_test_env:
    :return:
    '''
    global_test_env.deploy_all()
    return client_noc_list_obj


@pytest.fixture()
def client_candidate_obj(global_test_env, client_consensus_obj, client_list_obj):
    '''
    获取单个候选节点Client对象
    :param global_test_env:
    :return:
    '''
    address = client_consensus_obj.node.staking_address
    if not client_consensus_obj.staking.get_candidate_list_not_verifier():
        log.info('不存在候选节点，需要对节点进行质押')
        candidate_list = get_pledge_list(client_consensus_obj.node.ppos.getCandidateList)
        for normal_node_obj in global_test_env.normal_node_list:
            if normal_node_obj.node_id not in candidate_list:
                client_obj = get_client_obj(normal_node_obj.node_id, client_list_obj)
                log.info('对节点{}进行质押操作'.format(normal_node_obj.node_id))
                result = client_obj.staking.create_staking(0, address, address)
                log.info('节点{}质押结果为{}'.format(normal_node_obj.node_id, result))
                assert result.get('Code') == 0
        client_consensus_obj.economic.wait_settlement_blocknum(client_consensus_obj.node)
    node_id_list = client_consensus_obj.staking.get_candidate_list_not_verifier()
    log.info('候选非验证人列表为{}'.format(node_id_list))
    if not node_id_list:
        raise Exception('获取候选人失败')
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

def param_governance_verify(client_obj, module, name, newvalue, effectiveflag=None):
    '''
    effectiveflag indicates whether it takes effect
    :param pip_obj:
    :param module:
    :param name:
    :param newvalue:
    :param effectiveflag:
    :param number:
    :return:
    '''
    if isinstance(client_obj, Client):
        pip_obj = client_obj.pip
    if pip_obj.is_exist_effective_proposal_for_vote(pip_obj.cfg.param_proposal) or \
            pip_obj.is_exist_effective_proposal_for_vote(pip_obj.cfg.version_proposal):
        raise Exception('There is effective param proposal or version proposal')
    result = pip_obj.submitParam(pip_obj.node.node_id, str(time.time()), module, name, newvalue, pip_obj.node.staking_address,
                                 transaction_cfg=pip_obj.cfg.transaction_cfg)
    log.info('submit param proposal result : {}'.format(result))
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
    if not effectiveflag:
        for client_obj in client_verifier_obj_list:
            result = client_obj.pip.vote(client_obj.node.node_id, proposalinfo.get('ProposalID'),
                                         client_obj.pip.cfg.vote_option_yeas,
                                         client_obj.node.staking_address, transaction_cfg=client_obj.pip.cfg.transaction_cfg)
            log.info('Node {} vote proposal result : {}'.format(client_obj.node.node_id, result))
    wait_block_number(pip_obj.node, proposalinfo.get('EndVotingBlock'))
    if not effectiveflag:
        assert pip_obj.get_status_of_proposal(proposalinfo.get('ProposalID')) == 2
        log.info("blockNumber {}, the {} has become {}".format(proposalinfo.get('EndVotingBlock'), name, newvalue))
    else:
        assert pip_obj.get_status_of_proposal(proposalinfo.get('ProposalID')) == 3
        log.info("{} retains the original value".format(name))
