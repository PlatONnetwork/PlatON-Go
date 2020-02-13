from dacite import from_dict
from tests.lib import Genesis
from tests.lib.utils import assert_code, wait_block_number, upload_platon
from tests.lib.client import get_client_by_nodeid
import pytest
import time, os
from tests.govern.test_voting_statistics import submitvpandvote
from common.log import log


class TestPlatonVersion:
    @pytest.mark.P2
    def test_VE_DE_001(self, new_genesis_env, client_noconsensus):
        genesis = from_dict(data_class=Genesis, data=new_genesis_env.genesis_config)
        genesis.config.genesisVersion = client_noconsensus.pip.cfg.version2
        new_genesis_env.set_genesis(genesis.to_dict())
        new_genesis_env.deploy_all()
        client = client_noconsensus
        address, _ = client.economic.account.generate_account(client.node.web3, 10 ** 18 * 10000000)
        result = client.staking.create_staking(0, address, address)
        assert_code(result, 0)

    @pytest.mark.P2
    def test_VE_DE_002(self, new_genesis_env, client_noconsensus):
        genesis = from_dict(data_class=Genesis, data=new_genesis_env.genesis_config)
        genesis.config.genesisVersion = client_noconsensus.pip.cfg.version1
        new_genesis_env.set_genesis(genesis.to_dict())
        new_genesis_env.deploy_all()
        client = client_noconsensus
        address, _ = client.economic.account.generate_account(client.node.web3, 10 ** 18 * 10000000)
        result = client.staking.create_staking(0, address, address)
        assert_code(result, 301005)

    @pytest.mark.P2
    def test_VE_DE_004(self, new_genesis_env, client_noconsensus):
        genesis = from_dict(data_class=Genesis, data=new_genesis_env.genesis_config)
        genesis.config.genesisVersion = client_noconsensus.pip.cfg.version3
        new_genesis_env.set_genesis(genesis.to_dict())
        new_genesis_env.deploy_all()
        client = client_noconsensus
        address, _ = client.economic.account.generate_account(client.node.web3, 10 ** 18 * 10000000)
        result = client.staking.create_staking(0, address, address)
        assert_code(result, 0)

    @pytest.mark.P2
    def test_VE_DE_005(self, new_genesis_env, client_noconsensus):
        genesis = from_dict(data_class=Genesis, data=new_genesis_env.genesis_config)
        genesis.config.genesisVersion = client_noconsensus.pip.cfg.version7
        new_genesis_env.set_genesis(genesis.to_dict())
        new_genesis_env.deploy_all()
        client = client_noconsensus
        address, _ = client.economic.account.generate_account(client.node.web3, 10 ** 18 * 10000000)
        result = client.staking.create_staking(0, address, address)
        assert_code(result, 301004)

    @pytest.mark.P2
    def test_VE_DE_006(self, new_genesis_env, client_noconsensus):
        genesis = from_dict(data_class=Genesis, data=new_genesis_env.genesis_config)
        genesis.config.genesisVersion = client_noconsensus.pip.cfg.version8
        new_genesis_env.set_genesis(genesis.to_dict())
        new_genesis_env.deploy_all()
        client = client_noconsensus
        address, _ = client.economic.account.generate_account(client.node.web3, 10 ** 18 * 10000000)
        result = client.staking.create_staking(0, address, address)
        assert_code(result, 301004)

    @pytest.mark.P2
    def test_VE_AD_002(self, new_genesis_env):
        genesis = from_dict(data_class=Genesis, data=new_genesis_env.genesis_config)
        genesis.config.genesisVersion = 1796
        file = os.path.join(os.path.dirname(new_genesis_env.cfg.genesis_tmp), 'genesis_tmp2.json')
        genesis.to_file(file)
        consensus_node = new_genesis_env.get_rand_node()
        test_node = new_genesis_env.get_a_normal_node()
        test_node.clean()
        test_node.deploy_me(file)
        test_node.admin.addPeer(consensus_node.enode)
        time.sleep(5)
        assert test_node.web3.net.peerCount == 0
        assert test_node.block_number == 0

    @pytest.mark.P2
    def test_VE_AD_001(self, new_genesis_env, all_clients):
        consensus_node = new_genesis_env.get_rand_node()
        test_node = new_genesis_env.get_a_normal_node()
        test_node.clean()
        test_node.deploy_me(new_genesis_env.cfg.genesis_tmp)
        test_node.admin.addPeer(consensus_node.enode)
        time.sleep(5)
        assert test_node.web3.net.peerCount > 0, 'Join the chain failed'
        assert test_node.block_number > 0, "Non-consensus node sync block failed, block height: {}".format(test_node.block_number)
        time.sleep(5)
        client = get_client_by_nodeid(test_node.node_id, all_clients)
        address, _ = client.economic.account.generate_account(client.node.web3, 10**18 * 10000000)
        result = client.staking.create_staking(0, address, address)
        assert_code(result, 0)

    @pytest.mark.P2
    def test_VE_AD_004(self, new_genesis_env, clients_consensus):
        submitvpandvote(clients_consensus)
        proposalinfo = clients_consensus[0].pip.get_effect_proposal_info_of_vote()
        log.info('Get version proposal information : {}'.format(proposalinfo))
        wait_block_number(clients_consensus[0].node, proposalinfo.get('ActiveBlock'))
        genesis = from_dict(data_class=Genesis, data=new_genesis_env.genesis_config)
        genesis.config.genesisVersion = 2049
        file = os.path.join(os.path.dirname(new_genesis_env.cfg.genesis_tmp), 'genesis_tmp2.json')
        genesis.to_file(file)
        consensus_node = new_genesis_env.get_rand_node()
        test_node = new_genesis_env.get_a_normal_node()
        test_node.clean()
        test_node.deploy_me(file)
        test_node.admin.addPeer(consensus_node.enode)
        time.sleep(5)
        assert test_node.web3.net.peerCount == 0
        assert test_node.block_number == 0

    @pytest.mark.P2
    def test_VE_AD_003(self, new_genesis_env, clients_consensus, all_clients):
        submitvpandvote(clients_consensus)
        proposalinfo = clients_consensus[0].pip.get_effect_proposal_info_of_vote()
        log.info('Get version proposal information : {}'.format(proposalinfo))
        wait_block_number(clients_consensus[0].node, proposalinfo.get('ActiveBlock'))
        consensus_node = new_genesis_env.get_rand_node()
        test_node = new_genesis_env.get_a_normal_node()
        test_node.clean()
        test_node.deploy_me(new_genesis_env.cfg.genesis_tmp)
        test_node.admin.addPeer(consensus_node.enode)
        time.sleep(5)
        assert test_node.web3.net.peerCount > 0, 'Join the chain failed'
        assert test_node.block_number > 0, "Non-consensus node sync block failed, block height: {}".format(test_node.block_number)
        time.sleep(5)
        client = get_client_by_nodeid(test_node.node_id, all_clients)
        address, _ = client.economic.account.generate_account(client.node.web3, 10**18 * 10000000)
        result = client.staking.create_staking(0, address, address)
        assert_code(result, 301004)
        upload_platon(test_node, client.pip.cfg.PLATON_NEW_BIN)
        test_node.restart()
        result = client.staking.create_staking(0, address, address)
        assert_code(result, 0)
