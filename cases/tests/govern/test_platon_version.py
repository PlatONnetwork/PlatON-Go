from dacite import from_dict
from tests.lib import Genesis
from tests.lib.utils import assert_code

class TestPlatonVersion:
    def test_VE_DE_001(self, new_genesis_env, client_noconsensus):
        genesis = from_dict(data_class=Genesis, data=new_genesis_env.genesis_config)
        genesis.config.genesisVersion = client_noconsensus.pip.cfg.version2
        new_genesis_env.set_genesis(genesis.to_dict())
        new_genesis_env.deploy_all()
        client = client_noconsensus
        address, _ = client.economic.account.generate_account(client.node.web3, 10 ** 18 * 10000000)
        result = client.staking.create_staking(0, address, address)
        assert_code(result, 0)

    def test_VE_DE_002(self, new_genesis_env, client_noconsensus):
        genesis = from_dict(data_class=Genesis, data=new_genesis_env.genesis_config)
        genesis.config.genesisVersion = client_noconsensus.pip.cfg.version1
        new_genesis_env.set_genesis(genesis.to_dict())
        new_genesis_env.deploy_all()
        client = client_noconsensus
        address, _ = client.economic.account.generate_account(client.node.web3, 10 ** 18 * 10000000)
        result = client.staking.create_staking(0, address, address)
        assert_code(result, 301005)

    def test_VE_DE_004(self, new_genesis_env, client_noconsensus):
        genesis = from_dict(data_class=Genesis, data=new_genesis_env.genesis_config)
        genesis.config.genesisVersion = client_noconsensus.pip.cfg.version3
        new_genesis_env.set_genesis(genesis.to_dict())
        new_genesis_env.deploy_all()
        client = client_noconsensus
        address, _ = client.economic.account.generate_account(client.node.web3, 10 ** 18 * 10000000)
        result = client.staking.create_staking(0, address, address)
        assert_code(result, 0)

    def test_VE_DE_005(self, new_genesis_env, client_noconsensus):
        genesis = from_dict(data_class=Genesis, data=new_genesis_env.genesis_config)
        genesis.config.genesisVersion = client_noconsensus.pip.cfg.version7
        new_genesis_env.set_genesis(genesis.to_dict())
        new_genesis_env.deploy_all()
        client = client_noconsensus
        address, _ = client.economic.account.generate_account(client.node.web3, 10 ** 18 * 10000000)
        result = client.staking.create_staking(0, address, address)
        assert_code(result, 301004)

    def test_VE_DE_006(self, new_genesis_env, client_noconsensus):
        genesis = from_dict(data_class=Genesis, data=new_genesis_env.genesis_config)
        genesis.config.genesisVersion = client_noconsensus.pip.cfg.version8
        new_genesis_env.set_genesis(genesis.to_dict())
        new_genesis_env.deploy_all()
        client = client_noconsensus
        address, _ = client.economic.account.generate_account(client.node.web3, 10 ** 18 * 10000000)
        result = client.staking.create_staking(0, address, address)
        assert_code(result, 301004)
