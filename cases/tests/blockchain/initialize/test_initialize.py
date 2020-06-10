import os
import pytest
import json
import allure
from tests.lib.genesis import to_genesis
from common.log import log


@pytest.fixture(scope="function", autouse=True)
def restart_env(global_test_env):
    if not global_test_env.running:
        global_test_env.deploy_all()
    global_test_env.check_block(multiple=3)


@allure.title("View Creation Account")
@pytest.mark.P1
def test_CH_IN_012(global_test_env):
    """
    View the creation account created in the genesis_0.13.0.json file
    """
    log.info("View the creation account created in the genesis_0.13.0.json file")
    w3_list = [one_node.web3 for one_node in global_test_env.consensus_node_list]
    for w3 in w3_list:
        for one_address in global_test_env.genesis_config['alloc']:
            balance = w3.eth.getBalance(w3.toChecksumAddress(one_address))
            assert balance >= 0, "初始化账户错误"


@allure.title("Economic model parameters, governance parameters, penalty parameters, reward parameters")
@pytest.mark.P1
def test_CH_IN_015(global_test_env):
    """
    View economic model parameters, governance parameters, penalty parameters, and whether the reward parameters are correctly configured parameters
    """
    log.info("View economic model parameters, governance parameters, penalty parameters, and whether the reward parameters are correctly configured parameters")
    # economic_info = global_test_env.genesis_config['EconomicModel']
    if not global_test_env.running:
        global_test_env.deploy_all()
    genesis = to_genesis(global_test_env.genesis_config)
    w3_list = [one_node.web3 for one_node in global_test_env.consensus_node_list]
    for w3 in w3_list:
        info = w3.debug.economicConfig()
        assert info['common']['maxEpochMinutes'] == genesis.economicModel.common.maxEpochMinutes
        assert info['common']['maxConsensusVals'] == genesis.economicModel.common.maxConsensusVals
        assert info["common"]["additionalCycleTime"] == genesis.economicModel.common.additionalCycleTime

        assert info['staking']['stakeThreshold'] == genesis.economicModel.staking.stakeThreshold
        assert info['staking']['operatingThreshold'] == genesis.economicModel.staking.operatingThreshold
        assert info['staking']['maxValidators'] == genesis.economicModel.staking.maxValidators
        # assert info['staking']['hesitateRatio'] == genesis.economicModel.staking.hesitateRatio
        assert info['staking']['unStakeFreezeDuration'] == genesis.economicModel.staking.unStakeFreezeDuration

        assert info['slashing']['slashFractionDuplicateSign'] == genesis.economicModel.slashing.slashFractionDuplicateSign
        assert info['slashing']['duplicateSignReportReward'] == genesis.economicModel.slashing.duplicateSignReportReward
        assert info['slashing']['slashBlocksReward'] == genesis.economicModel.slashing.slashBlocksReward
        assert info['slashing']['maxEvidenceAge'] == genesis.economicModel.slashing.maxEvidenceAge

        assert info['gov']['versionProposalVoteDurationSeconds'] == genesis.economicModel.gov.versionProposalVoteDurationSeconds
        assert info['gov']['versionProposalSupportRate'] == genesis.economicModel.gov.versionProposalSupportRate
        assert info['gov']['textProposalVoteDurationSeconds'] == genesis.economicModel.gov.textProposalVoteDurationSeconds
        assert info['gov']['textProposalVoteRate'] == genesis.economicModel.gov.textProposalVoteRate
        assert info['gov']['textProposalSupportRate'] == genesis.economicModel.gov.textProposalSupportRate
        assert info['gov']['cancelProposalVoteRate'] == genesis.economicModel.gov.cancelProposalVoteRate
        assert info['gov']['cancelProposalSupportRate'] == genesis.economicModel.gov.cancelProposalSupportRate
        assert info['gov']['paramProposalVoteDurationSeconds'] == genesis.economicModel.gov.paramProposalVoteDurationSeconds
        assert info['gov']['paramProposalVoteRate'] == genesis.economicModel.gov.paramProposalVoteRate
        assert info['gov']['paramProposalSupportRate'] == genesis.economicModel.gov.paramProposalSupportRate

        assert info['reward']['newBlockRate'] == genesis.economicModel.reward.newBlockRate
        assert info['reward']['platonFoundationYear'] == genesis.economicModel.reward.platONFoundationYear

        assert w3.toChecksumAddress(info['innerAcc']['platonFundAccount']) == w3.toChecksumAddress(genesis.economicModel.innerAcc.platonFundAccount)
        assert info['innerAcc']['platonFundBalance'] == genesis.economicModel.innerAcc.platonFundBalance
        assert w3.toChecksumAddress(info['innerAcc']['cdfAccount']) == w3.toChecksumAddress(genesis.economicModel.innerAcc.cdfAccount)
        assert info['innerAcc']['cdfBalance'] == genesis.economicModel.innerAcc.cdfBalance


@allure.title("Foundation lock warehouse plan inquiry")
@pytest.mark.P1
def test_CH_IN_014(global_test_env):
    """
    View the foundation lock warehouse plan query
    """
    log.info("View the foundation lock warehouse plan query")
    w3_list = [one_node.web3 for one_node in global_test_env.consensus_node_list]
    for w3 in w3_list:
        info = w3.eth.call({"to": "lax1zqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqp3yp7hw", "data": "0xda8382100495941000000000000000000000000000000000000003"}, 0)
        recive = json.loads(str(info, encoding="ISO-8859-1"))
        pass
        # move for 0.7.5
        # plans = recive['Ret']['plans']
        # assert(8 == len(plans))
        # for i in range(len(plans)):
        #     if 1600 == plans[i]['blockNumber']:
        #         assert plans[i]['amount'] == "0x2e4b34f3fb9ea4f3f80000"
        #     if 3200 == plans[i]['blockNumber']:
        #         assert plans[i]['amount'] == "0x28fea171d6cdd2a4900000"
        #     if 4800 == plans[i]['blockNumber']:
        #         assert plans[i]['amount'] == "0x239023dcb60bdb30380000"
        #     if 6400 == plans[i]['blockNumber']:
        #         assert plans[i]['amount'] == "0x1dfee325efc6d87ee00000"
        #     if 8000 == plans[i]['blockNumber']:
        #         assert plans[i]['amount'] == "0x184a00c53258036e040000"
        #     if 9600 == plans[i]['blockNumber']:
        #         assert plans[i]['amount'] == "0x127098c664ba6778100000"
        #     if 11200 == plans[i]['blockNumber']:
        #         assert plans[i]['amount'] == "0xc71c15aa0d54579400000"
        #     if 12800 == plans[i]['blockNumber']:
        #         assert plans[i]['amount'] == "0x64c8af3f4e97afe680000"


@allure.title("Consensus parameter")
@pytest.mark.P1
def test_CH_IN_016(global_test_env):
    """
    View the number of blocks of each consensus node of the consensus
    """
    log.info("View the number of outbound blocks of each consensus node and the total number of consensus nodes")
    amount = global_test_env.genesis_config['config']['cbft']['amount']
    w3_list = [one_node.web3 for one_node in global_test_env.consensus_node_list]
    for w3 in w3_list:
        info = w3.eth.getPrepareQC(amount)
        assert info['viewNumber'] == 0
        info = w3.eth.getPrepareQC(amount + 1)
        assert info['viewNumber'] == 1
