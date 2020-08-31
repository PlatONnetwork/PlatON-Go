from dacite import from_dict
from tests.lib import Genesis
from tests.lib.utils import assert_code, wait_block_number, upload_platon, get_pledge_list
from tests.lib.client import get_client_by_nodeid
import pytest
import time, os
from tests.govern.test_voting_statistics import submitvpandvote
from common.log import log
from hexbytes import HexBytes
from common.connect import connect_web3


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
    def test_VE_DE_004_VE_DE_011(self, new_genesis_env, client_noconsensus):
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
    def test_VE_DE_007(self, new_genesis_env):
        genesis = from_dict(data_class=Genesis, data=new_genesis_env.genesis_config)
        genesis.economicModel.slashing.zeroProduceCumulativeTime = 5
        new_genesis_env.set_genesis(genesis.to_dict())
        try:
            new_genesis_env.deploy_all()
        except Exception as e:
            log.info('Deploy failed error measage {}'.format(e.args[0]))
            index = e.args[0].find('ZeroProduceCumulativeTime')
            assert e.args[0][index:index + 40] == r'ZeroProduceCumulativeTime must be [1, 4]'

    @pytest.mark.P2
    def test_VE_DE_008_VE_DE_009(self, new_genesis_env, client_noconsensus):
        genesis = from_dict(data_class=Genesis, data=new_genesis_env.genesis_config)
        genesis.economicModel.slashing.zeroProduceCumulativeTime = 4
        genesis.economicModel.slashing.zeroProduceNumberThreshold = 4
        new_genesis_env.set_genesis(genesis.to_dict())
        new_genesis_env.deploy_all()
        time.sleep(3)
        assert client_noconsensus.node.block_number > 0

    @pytest.mark.P2
    def test_VE_DE_008_VE_DE_009(self, new_genesis_env, client_noconsensus):
        genesis = from_dict(data_class=Genesis, data=new_genesis_env.genesis_config)
        genesis.economicModel.slashing.zeroProduceCumulativeTime = 4
        genesis.economicModel.slashing.zeroProduceNumberThreshold = 4
        new_genesis_env.set_genesis(genesis.to_dict())
        new_genesis_env.deploy_all()
        time.sleep(10)
        assert client_noconsensus.node.block_number > 0

    @pytest.mark.P2
    def test_VE_DE_010(self, new_genesis_env):
        genesis = from_dict(data_class=Genesis, data=new_genesis_env.genesis_config)
        genesis.economicModel.slashing.zeroProduceCumulativeTime = 3
        genesis.economicModel.slashing.zeroProduceNumberThreshold = 4
        new_genesis_env.set_genesis(genesis.to_dict())
        try:
            new_genesis_env.deploy_all()
        except Exception as e:
            log.info('Deploy failed error measage {}'.format(e.args[0]))
            index = e.args[0].find('ZeroProduceNumberThreshold')

            assert e.args[0][index:index + 41] == r'ZeroProduceNumberThreshold must be [1, 3]'

    @pytest.mark.P2
    def test_VE_DE_010(self, new_genesis_env):
        genesis = from_dict(data_class=Genesis, data=new_genesis_env.genesis_config)
        genesis.economicModel.slashing.zeroProduceNumberThreshold = 0
        new_genesis_env.set_genesis(genesis.to_dict())
        try:
            new_genesis_env.deploy_all()
        except Exception as e:
            log.info('Deploy failed error measage {}'.format(e.args[0]))
            index = e.args[0].find('ZeroProduceNumberThreshold')
            assert e.args[0][index:index + 41] == r'ZeroProduceNumberThreshold must be [1, 1]'

    @pytest.mark.P2
    def test_VE_DE_014_017(self, new_genesis_env):
        genesis = from_dict(data_class=Genesis, data=new_genesis_env.genesis_config)
        genesis.economicModel.staking.rewardPerMaxChangeRange = 0
        new_genesis_env.set_genesis(genesis.to_dict())
        try:
            new_genesis_env.deploy_all()
        except Exception as e:
            log.info('Deploy failed error measage {}'.format(e.args[0]))
            index = e.args[0].find('RewardPerMaxChangeRange')
            assert e.args[0][index:index + 41] == r'RewardPerMaxChangeRange must be [1, 2000]'

        genesis.economicModel.staking.rewardPerMaxChangeRange = 2001
        new_genesis_env.set_genesis(genesis.to_dict())
        try:
            new_genesis_env.deploy_all()
        except Exception as e:
            log.info('Deploy failed error measage {}'.format(e.args[0]))
            index = e.args[0].find('RewardPerMaxChangeRange')
            assert e.args[0][index:index + 41] == r'RewardPerMaxChangeRange must be [1, 2000]'

    @pytest.mark.P2
    def test_VE_DE_015_016(self, new_genesis_env, client_consensus):
        genesis = from_dict(data_class=Genesis, data=new_genesis_env.genesis_config)
        genesis.economicModel.staking.rewardPerMaxChangeRange = 1
        new_genesis_env.set_genesis(genesis.to_dict())
        new_genesis_env.deploy_all()
        time.sleep(5)
        assert client_consensus.node.block_number != 0

        genesis.economicModel.staking.rewardPerMaxChangeRange = 2000
        new_genesis_env.set_genesis(genesis.to_dict())
        new_genesis_env.deploy_all()
        time.sleep(5)
        assert client_consensus.node.block_number != 0

    @pytest.mark.P2
    def test_VE_DE_018_021(self, new_genesis_env):
        genesis = from_dict(data_class=Genesis, data=new_genesis_env.genesis_config)
        genesis.economicModel.staking.rewardPerChangeInterval = 1
        new_genesis_env.set_genesis(genesis.to_dict())
        try:
            new_genesis_env.deploy_all()
        except Exception as e:
            log.info(type(e.args[0]), 'Deploy failed error measage {}'.format(e.args[0]))
            index = e.args[0].find('RewardPerChangeInterval must be [2, 28]')
            assert index != -1

        genesis.economicModel.staking.rewardPerMaxChangeRange = 2001
        new_genesis_env.set_genesis(genesis.to_dict())
        try:
            new_genesis_env.deploy_all()
        except Exception as e:
            log.info('Deploy failed error measage {}'.format(e.args[0]))
            index = e.args[0].find('RewardPerMaxChangeRange must be [1, 2000]')
            assert index != -1

    @pytest.mark.P2
    def test_VE_DE_019_020(self, new_genesis_env, client_consensus):
        genesis = from_dict(data_class=Genesis, data=new_genesis_env.genesis_config)
        genesis.economicModel.staking.rewardPerChangeInterval = 2
        new_genesis_env.set_genesis(genesis.to_dict())
        new_genesis_env.deploy_all()
        time.sleep(5)
        assert client_consensus.node.block_number != 0

        settlement_count = client_consensus.economic.additional_cycle_time * 60//(
                client_consensus.economic.settlement_size * client_consensus.economic.interval)
        genesis.economicModel.staking.rewardPerChangeInterval = settlement_count
        new_genesis_env.set_genesis(genesis.to_dict())
        new_genesis_env.deploy_all()
        time.sleep(5)
        assert client_consensus.node.block_number != 0

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

class TestInit:
    @pytest.mark.P2
    def test_HA_IN_001(self, new_genesis_env, client_consensus):
        blockhash = HexBytes(client_consensus.node.eth.getBlock(0).get('hash')).hex()
        log.info('init block hash : {}'.format(blockhash))
        genesis = from_dict(data_class=Genesis, data=new_genesis_env.genesis_config)
        genesis.economicModel.staking.stakeThreshold = 2*client_consensus.economic.genesis.economicModel.staking.stakeThreshold
        new_genesis_env.set_genesis(genesis.to_dict())
        new_genesis_env.deploy_all()
        stakingThreshold_blockhash = HexBytes(client_consensus.node.eth.getBlock(0).get('hash')).hex()
        log.info('init block hash stakingThreshold_blockhash : {}'.format(stakingThreshold_blockhash))
        assert blockhash != stakingThreshold_blockhash

        genesis.economicModel.staking.operatingThreshold = 2*client_consensus.economic.genesis.economicModel.staking.operatingThreshold
        new_genesis_env.set_genesis(genesis.to_dict())
        new_genesis_env.deploy_all()
        operatingThreshold_blockhash = HexBytes(client_consensus.node.eth.getBlock(0).get('hash')).hex()
        log.info('init block hash operatingThreshold_blockhash : {}'.format(operatingThreshold_blockhash))
        assert stakingThreshold_blockhash != operatingThreshold_blockhash

        genesis.economicModel.staking.maxValidators = client_consensus.economic.genesis.economicModel.staking.maxValidators + 1
        genesis.economicModel.common.maxEpochMinutes = 4
        new_genesis_env.set_genesis(genesis.to_dict())
        new_genesis_env.deploy_all()
        maxValidators_blockhash = HexBytes(client_consensus.node.eth.getBlock(0).get('hash')).hex()
        log.info('init block hash maxValidators_blockhash : {}'.format(maxValidators_blockhash))
        assert operatingThreshold_blockhash != maxValidators_blockhash

        genesis.economicModel.staking.unStakeFreezeDuration = 2*client_consensus.economic.genesis.economicModel.staking.unStakeFreezeDuration
        new_genesis_env.set_genesis(genesis.to_dict())
        new_genesis_env.deploy_all()
        unStakeFreezeDuration_blockhash = HexBytes(client_consensus.node.eth.getBlock(0).get('hash')).hex()
        log.info('init block hash unStakeFreezeDuration_blockhash : {}'.format(unStakeFreezeDuration_blockhash))
        assert maxValidators_blockhash != unStakeFreezeDuration_blockhash

        genesis.economicModel.slashing.slashBlocksReward = 2*client_consensus.economic.genesis.economicModel.slashing.slashBlocksReward
        new_genesis_env.set_genesis(genesis.to_dict())
        new_genesis_env.deploy_all()
        slashBlocksReward_blockhash = HexBytes(client_consensus.node.eth.getBlock(0).get('hash')).hex()
        log.info('init block hash slashBlocksReward_blockhash : {}'.format(slashBlocksReward_blockhash))
        assert unStakeFreezeDuration_blockhash != slashBlocksReward_blockhash

        genesis.economicModel.slashing.slashFractionDuplicateSign = client_consensus.economic.genesis.economicModel.slashing.slashFractionDuplicateSign - 10
        new_genesis_env.set_genesis(genesis.to_dict())
        new_genesis_env.deploy_all()
        slashFractionDuplicateSign_blockhash = HexBytes(client_consensus.node.eth.getBlock(0).get('hash')).hex()
        log.info('init block hash slashFractionDuplicateSign_blockhash : {}'.format(slashFractionDuplicateSign_blockhash))
        assert slashBlocksReward_blockhash != slashFractionDuplicateSign_blockhash

        genesis.economicModel.slashing.duplicateSignReportReward = client_consensus.economic.genesis.economicModel.slashing.duplicateSignReportReward + 10
        new_genesis_env.set_genesis(genesis.to_dict())
        new_genesis_env.deploy_all()
        duplicateSignReportReward_blockhash = HexBytes(client_consensus.node.eth.getBlock(0).get('hash')).hex()
        log.info('init block hash duplicateSignReportReward_blockhash : {}'.format(duplicateSignReportReward_blockhash))
        assert slashFractionDuplicateSign_blockhash != duplicateSignReportReward_blockhash

        genesis.economicModel.slashing.maxEvidenceAge = client_consensus.economic.genesis.economicModel.slashing.maxEvidenceAge + 1
        new_genesis_env.set_genesis(genesis.to_dict())
        new_genesis_env.deploy_all()
        maxEvidenceAge_blockhash = HexBytes(client_consensus.node.eth.getBlock(0).get('hash')).hex()
        log.info('init block hash maxEvidenceAge_blockhash : {}'.format(maxEvidenceAge_blockhash))
        assert duplicateSignReportReward_blockhash != maxEvidenceAge_blockhash

        genesis.gasLimit = str(int(client_consensus.economic.genesis.gasLimit) + 100)
        new_genesis_env.set_genesis(genesis.to_dict())
        new_genesis_env.deploy_all()
        maxBlockGasLimit_blockhash = HexBytes(client_consensus.node.eth.getBlock(0).get('hash')).hex()
        log.info('init block hash maxBlockGasLimit_blockhash : {}'.format(maxBlockGasLimit_blockhash))
        assert maxEvidenceAge_blockhash != maxBlockGasLimit_blockhash

        genesis.economicModel.gov.versionProposalVoteDurationSeconds = genesis.economicModel.gov.versionProposalVoteDurationSeconds + 100
        new_genesis_env.set_genesis(genesis.to_dict())
        new_genesis_env.deploy_all()
        versionProposalVoteDurationSeconds_blockhash = HexBytes(client_consensus.node.eth.getBlock(0).get('hash')).hex()
        log.info('init block hash versionProposalVoteDurationSeconds_blockhash : {}'.format(
            versionProposalVoteDurationSeconds_blockhash))
        assert maxBlockGasLimit_blockhash != versionProposalVoteDurationSeconds_blockhash

        genesis.economicModel.gov.versionProposalSupportRate = genesis.economicModel.gov.versionProposalSupportRate + 100
        new_genesis_env.set_genesis(genesis.to_dict())
        new_genesis_env.deploy_all()
        versionProposalSupportRate_blockhash = HexBytes(client_consensus.node.eth.getBlock(0).get('hash')).hex()
        log.info('init block hash versionProposalSupportRate_blockhash : {}'.format(
            versionProposalSupportRate_blockhash))
        assert versionProposalVoteDurationSeconds_blockhash != versionProposalSupportRate_blockhash

        genesis.economicModel.gov.textProposalVoteDurationSeconds = genesis.economicModel.gov.textProposalVoteDurationSeconds + 100
        new_genesis_env.set_genesis(genesis.to_dict())
        new_genesis_env.deploy_all()
        textProposalVoteDurationSeconds_blockhash = HexBytes(client_consensus.node.eth.getBlock(0).get('hash')).hex()
        log.info('init block hash textProposalVoteDurationSeconds_blockhash : {}'.format(
            textProposalVoteDurationSeconds_blockhash))
        assert versionProposalSupportRate_blockhash != textProposalVoteDurationSeconds_blockhash

        genesis.economicModel.gov.textProposalVoteRate = genesis.economicModel.gov.textProposalVoteRate + 100
        new_genesis_env.set_genesis(genesis.to_dict())
        new_genesis_env.deploy_all()
        textProposalVoteRate_blockhash = HexBytes(client_consensus.node.eth.getBlock(0).get('hash')).hex()
        log.info('init block hash textProposalVoteRate_blockhash : {}'.format(
            textProposalVoteRate_blockhash))
        assert textProposalVoteDurationSeconds_blockhash != textProposalVoteRate_blockhash

        genesis.economicModel.gov.textProposalSupportRate = genesis.economicModel.gov.textProposalSupportRate + 100
        new_genesis_env.set_genesis(genesis.to_dict())
        new_genesis_env.deploy_all()
        textProposalSupportRate_blockhash = HexBytes(client_consensus.node.eth.getBlock(0).get('hash')).hex()
        log.info('init block hash textProposalSupportRate_blockhash : {}'.format(
            textProposalSupportRate_blockhash))
        assert textProposalVoteRate_blockhash != textProposalSupportRate_blockhash

        genesis.economicModel.gov.cancelProposalVoteRate = genesis.economicModel.gov.cancelProposalVoteRate + 100
        new_genesis_env.set_genesis(genesis.to_dict())
        new_genesis_env.deploy_all()
        cancelProposalVoteRate_blockhash = HexBytes(client_consensus.node.eth.getBlock(0).get('hash')).hex()
        log.info('init block hash cancelProposalVoteRate_blockhash : {}'.format(
            cancelProposalVoteRate_blockhash))
        assert textProposalSupportRate_blockhash != cancelProposalVoteRate_blockhash

        genesis.economicModel.gov.cancelProposalSupportRate = genesis.economicModel.gov.cancelProposalSupportRate + 100
        new_genesis_env.set_genesis(genesis.to_dict())
        new_genesis_env.deploy_all()
        cancelProposalSupportRate_blockhash = HexBytes(client_consensus.node.eth.getBlock(0).get('hash')).hex()
        log.info('init block hash cancelProposalSupportRate_blockhash : {}'.format(
            cancelProposalSupportRate_blockhash))
        assert cancelProposalVoteRate_blockhash != cancelProposalSupportRate_blockhash

        genesis.economicModel.gov.paramProposalVoteDurationSeconds = genesis.economicModel.gov.paramProposalVoteDurationSeconds + 100
        new_genesis_env.set_genesis(genesis.to_dict())
        new_genesis_env.deploy_all()
        paramProposalVoteDurationSeconds_blockhash = HexBytes(client_consensus.node.eth.getBlock(0).get('hash')).hex()
        log.info('init block hash paramProposalVoteDurationSeconds_blockhash : {}'.format(
            paramProposalVoteDurationSeconds_blockhash))
        assert cancelProposalSupportRate_blockhash != paramProposalVoteDurationSeconds_blockhash

        genesis.economicModel.gov.paramProposalVoteRate = genesis.economicModel.gov.paramProposalVoteRate + 100
        new_genesis_env.set_genesis(genesis.to_dict())
        new_genesis_env.deploy_all()
        paramProposalVoteRate_blockhash = HexBytes(client_consensus.node.eth.getBlock(0).get('hash')).hex()
        log.info('init block hash paramProposalVoteRate_blockhash : {}'.format(
            paramProposalVoteRate_blockhash))
        assert paramProposalVoteDurationSeconds_blockhash != paramProposalVoteRate_blockhash

        genesis.economicModel.gov.paramProposalSupportRate = genesis.economicModel.gov.paramProposalSupportRate + 100
        new_genesis_env.set_genesis(genesis.to_dict())
        new_genesis_env.deploy_all()
        paramProposalSupportRate_blockhash = HexBytes(client_consensus.node.eth.getBlock(0).get('hash')).hex()
        log.info('init block hash paramProposalSupportRate_blockhash : {}'.format(
            paramProposalSupportRate_blockhash))
        assert paramProposalVoteRate_blockhash != paramProposalSupportRate_blockhash

        genesis.economicModel.reward.newBlockRate = genesis.economicModel.reward.newBlockRate + 1
        new_genesis_env.set_genesis(genesis.to_dict())
        new_genesis_env.deploy_all()
        newBlockRate_blockhash = HexBytes(client_consensus.node.eth.getBlock(0).get('hash')).hex()
        log.info('init block hash newBlockRate_blockhash : {}'.format(
            newBlockRate_blockhash))
        assert paramProposalSupportRate_blockhash != newBlockRate_blockhash

        genesis.config.chainId = genesis.config.chainId + 1
        new_genesis_env.set_genesis(genesis.to_dict())
        new_genesis_env.deploy_all()
        chainId_blockhash = HexBytes(client_consensus.node.eth.getBlock(0).get('hash')).hex()
        log.info('init block hash chainId_blockhash : {}'.format(
            chainId_blockhash))
        assert newBlockRate_blockhash == chainId_blockhash









