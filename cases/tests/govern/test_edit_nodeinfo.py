import pytest, time
from tests.conftest import param_governance_verify_before_endblock
from tests.lib.client import get_client_by_nodeid
from tests.lib.utils import get_governable_parameter_value, assert_code, wait_block_number
from common.log import log
from dacite import from_dict
from tests.lib.genesis import Genesis

@pytest.mark.P2
def test_UP_RE_005_012_013(new_genesis_env, client_noconsensus, client_consensus):
    genesis = from_dict(data_class=Genesis, data=new_genesis_env.genesis_config)
    genesis.economicModel.staking.rewardPerChangeInterval = 2
    genesis.economicModel.staking.rewardPerMaxChangeRange = 10
    new_genesis_env.set_genesis(genesis.to_dict())
    new_genesis_env.deploy_all()
    client = client_noconsensus
    client_2 = client_consensus
    address, _ = client.economic.account.generate_account(client.node.web3, 10**18*10000000)
    result = client.staking.create_staking(0, address, address, reward_per=100)
    log.info('nodeid {} staking result : {}'.format(client.node.node_id, result))
    assert_code(result, 0)
    result = client_2.staking.edit_candidate(client_2.node.staking_address, address, reward_per=1)
    log.info('edit nodeinfo result : {}'.format(result))
    assert_code(result, 301008)

    result = client.staking.edit_candidate(client.node.staking_address, address, reward_per=101)
    log.info('edit nodeinfo result : {}'.format(result))
    assert_code(result, 301008)

    client.economic.wait_settlement_blocknum(client.node, 1)
    result = client.staking.edit_candidate(client.node.staking_address, address, reward_per=109)
    log.info('edit nodeinfo result : {}'.format(result))
    assert_code(result, 0)
    result = client_2.staking.edit_candidate(client_2.node.staking_address, address, reward_per=9)
    log.info('edit nodeinfo result : {}'.format(result))
    assert_code(result, 0)

    client.economic.wait_settlement_blocknum(client.node, 1)
    result = client.staking.edit_candidate(client.node.staking_address, address, reward_per=100)
    log.info('edit nodeinfo result : {}'.format(result))
    assert_code(result, 0)
    result = client_2.staking.edit_candidate(client_2.node.staking_address, address, reward_per=0)
    log.info('edit nodeinfo result : {}'.format(result))
    assert_code(result, 0)

@pytest.mark.P2
def test_UP_RE_002_004(new_genesis_env, clients_consensus, client_noconsensus):
    genesis = from_dict(data_class=Genesis, data=new_genesis_env.genesis_config)
    genesis.economicModel.staking.rewardPerChangeInterval = 2
    genesis.economicModel.staking.rewardPerMaxChangeRange = 10
    genesis.economicModel.gov.paramProposalVoteDurationSeconds = 0
    new_genesis_env.set_genesis(genesis.to_dict())
    new_genesis_env.deploy_all()
    client = client_noconsensus
    client_2 = clients_consensus[0]
    address, _ = client.economic.account.generate_account(client.node.web3, 10**18*10000000)
    result = client.staking.create_staking(0, address, address, reward_per=100)
    log.info('nodeid {} staking result : {}'.format(client.node.node_id, result))
    assert_code(result, 0)
    result = client_2.staking.edit_candidate(client_2.node.staking_address, address, reward_per=0)
    log.info('edit nodeinfo result : {}'.format(result))
    assert_code(result, 0)

    param_governance_verify_before_endblock(client_2, 'staking', 'rewardPerMaxChangeRange', '11')
    client_2.economic.wait_settlement_blocknum(client_2.node, 1)
    assert '11' == get_governable_parameter_value(client_2, 'rewardPerMaxChangeRange')
    result = client_2.staking.edit_candidate(client_2.node.staking_address, address, reward_per=12)
    log.info('edit nodeinfo result : {}'.format(result))
    assert_code(result, 301009)

    result = client.staking.edit_candidate(client.node.staking_address, address, reward_per=88)
    log.info('edit nodeinfo result : {}'.format(result))
    assert_code(result, 301009)

    result = client_2.staking.edit_candidate(client_2.node.staking_address, address, reward_per=11)
    log.info('edit nodeinfo result : {}'.format(result))
    assert_code(result, 0)

    result = client.staking.edit_candidate(client.node.staking_address, address, reward_per=89)
    log.info('edit nodeinfo result : {}'.format(result))
    assert_code(result, 0)

    param_governance_verify_before_endblock(client_2, 'staking', 'rewardPerMaxChangeRange', '10')
    client_2.economic.wait_settlement_blocknum(client_2.node, 1)
    assert '10' == get_governable_parameter_value(client_2, 'rewardPerMaxChangeRange')
    result = client_2.staking.edit_candidate(client_2.node.staking_address, address, reward_per=22)
    log.info('edit nodeinfo result : {}'.format(result))
    assert_code(result, 301009)

    result = client.staking.edit_candidate(client.node.staking_address, address, reward_per=78)
    log.info('edit nodeinfo result : {}'.format(result))
    assert_code(result, 301009)

    result = client_2.staking.edit_candidate(client_2.node.staking_address, address, reward_per=21)
    log.info('edit nodeinfo result : {}'.format(result))
    assert_code(result, 0)

    result = client.staking.edit_candidate(client.node.staking_address, address, reward_per=79)
    log.info('edit nodeinfo result : {}'.format(result))
    assert_code(result, 0)

@pytest.mark.P2
def test_UP_RE_007_009(new_genesis_env, client_noconsensus, client_consensus):
    genesis = from_dict(data_class=Genesis, data=new_genesis_env.genesis_config)
    genesis.economicModel.staking.rewardPerChangeInterval = 2
    genesis.economicModel.staking.rewardPerMaxChangeRange = 10
    genesis.economicModel.gov.paramProposalVoteDurationSeconds = 0
    new_genesis_env.set_genesis(genesis.to_dict())
    new_genesis_env.deploy_all()
    client = client_noconsensus
    client_2 = client_consensus
    address, _ = client.economic.account.generate_account(client.node.web3, 10**18*10000000)
    result = client.staking.create_staking(0, address, address, reward_per=100)
    log.info('nodeid {} staking result : {}'.format(client.node.node_id, result))
    assert_code(result, 0)

    result = client_2.staking.edit_candidate(client_2.node.staking_address, address, reward_per=2)
    log.info('edit nodeinfo result : {}'.format(result))
    assert_code(result, 301008)

    param_governance_verify_before_endblock(client_2, 'staking', 'rewardPerChangeInterval', '3')
    client_2.economic.wait_settlement_blocknum(client_2.node)
    assert 3 == int(get_governable_parameter_value(client_2, 'rewardPerChangeInterval'))

    result = client_2.staking.edit_candidate(client_2.node.staking_address, address, reward_per=2)
    log.info('edit nodeinfo result : {}'.format(result))
    assert_code(result, 301008)
    result = client.staking.edit_candidate(client.node.staking_address, address, reward_per=99)
    log.info('edit nodeinfo result : {}'.format(result))
    assert_code(result, 301008)
    client_2.economic.wait_settlement_blocknum(client_2.node)
    result = client_2.staking.edit_candidate(client_2.node.staking_address, address, reward_per=2)
    log.info('edit nodeinfo result : {}'.format(result))
    assert_code(result, 0)
    result = client.staking.edit_candidate(client.node.staking_address, address, reward_per=99)
    log.info('edit nodeinfo result : {}'.format(result))
    assert_code(result, 301008)

    client.economic.wait_settlement_blocknum(client.node)
    result = client.staking.edit_candidate(client.node.staking_address, address, reward_per=99)
    log.info('edit nodeinfo result : {}'.format(result))
    assert_code(result, 0)

    param_governance_verify_before_endblock(client_2, 'staking', 'rewardPerChangeInterval', '2')
    client_2.economic.wait_settlement_blocknum(client_2.node)
    assert 2 == int(get_governable_parameter_value(client_2, 'rewardPerChangeInterval'))
    result = client.staking.edit_candidate(client.node.staking_address, address, reward_per=98)
    log.info('edit nodeinfo result : {}'.format(result))
    assert_code(result, 301008)

    client.economic.wait_settlement_blocknum(client.node)
    result = client.staking.edit_candidate(client.node.staking_address, address, reward_per=98)
    log.info('edit nodeinfo result : {}'.format(result))
    assert_code(result, 0)

@pytest.mark.P2
def test_UP_RE_001_003_006_008_010_UP_RE_011(new_genesis_env, client_noconsensus, client_consensus):
    genesis = from_dict(data_class=Genesis, data=new_genesis_env.genesis_config)
    genesis.economicModel.staking.rewardPerChangeInterval = 2
    genesis.economicModel.staking.rewardPerMaxChangeRange = 10
    new_genesis_env.set_genesis(genesis.to_dict())
    new_genesis_env.deploy_all()
    client = client_noconsensus
    client_2 = client_consensus
    address, _ = client.economic.account.generate_account(client.node.web3, 10**18*10000000)
    result = client.staking.create_staking(0, address, address, reward_per=100)
    log.info('nodeid {} staking result : {}'.format(client.node.node_id, result))
    assert_code(result, 0)

    wait_block_number(client.node, 2 * client.economic.settlement_size)
    result = client.staking.edit_candidate(client.node.staking_address, address, reward_per=111)
    log.info('edit nodeinfo result : {}'.format(result))
    assert_code(result, 301009)

    result = client.staking.edit_candidate(client.node.staking_address, address, reward_per=89)
    log.info('edit nodeinfo result : {}'.format(result))
    assert_code(result, 301009)

    address, _ = client.economic.account.generate_account(client.node.web3, 10**18*10)
    result = client.staking.edit_candidate(client.node.staking_address, address, reward_per=110)
    log.info('edit nodeinfo result : {}'.format(result))
    assert_code(result, 0)

    wait_block_number(client.node, 2 * client.economic.settlement_size)
    address, _ = client.economic.account.generate_account(client.node.web3, 10**18*10)
    result = client.staking.edit_candidate(client.node.staking_address, address, reward_per=100)
    log.info('edit nodeinfo result : {}'.format(result))
    assert_code(result, 0)

    wait_block_number(client.node, 2*client.economic.settlement_size - 10)
    result = client.staking.edit_candidate(client.node.staking_address, address, reward_per=111)
    log.info('edit nodeinfo result : {}'.format(result))
    assert_code(result, 301008)

    log.info('node {} candidate info {}'.format(client_2.node.node_id, client_2.ppos.getCandidateInfo(
        client_2.node.node_id)))

    result = client_2.staking.edit_candidate(client_2.node.staking_address, address, reward_per=10)
    log.info('edit nodeinfo result : {}'.format(result))
    assert_code(result, 0)

    result = client_2.staking.edit_candidate(client_2.node.staking_address, address, reward_per=31)
    log.info('edit nodeinfo result : {}'.format(result))
    assert_code(result, 301008)





