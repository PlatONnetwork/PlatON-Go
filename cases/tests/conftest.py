import pytest
from tests.lib import StakingConfig
from common.log import log
from tests.lib.client import Client, get_client_obj
from tests.lib.utils import get_pledge_list


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
def client_new_node_obj(global_test_env, client_consensus_obj, client_list_obj):
    '''
    获取单个新节点Client对象
    :param global_test_env:
    :return:
    '''
    candidate_list = get_pledge_list(client_consensus_obj.ppos.getCandidateList)
    log.info('candidatelist{}'.format(candidate_list))
    for noconsensus_node_obj in global_test_env.normal_node_list:
        if noconsensus_node_obj.node_id not in candidate_list:
            return get_client_obj(noconsensus_node_obj.node_id, client_list_obj)
    log.info('非共识节点已全部质押，重新启链')
    global_test_env.deploy_all()
    return get_client_obj(global_test_env.get_a_normal_node.node_id, client_list_obj)


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
