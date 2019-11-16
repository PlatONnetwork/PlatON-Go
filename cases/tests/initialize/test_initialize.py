import os
import pytest
import json
import allure
from environment.env import create_env
from conf.settings import NODE_FILE
from common.log import log

@pytest.fixture(scope="function", autouse=True)
def restart_env(global_test_env):
    if not global_test_env.running:
        global_test_env.deploy_all()
    global_test_env.check_block(multiple=3)


@allure.title("查看创世账户")
@pytest.mark.P1
def test_initial_account(global_test_env):
    """
    查看存在genesis.json文件中配置的创世账户
    """
    log.info("查看存在genesis.json文件中配置的创世账户")
    w3_list = [one_node.web3 for one_node in global_test_env.collusion_node_list]
    for w3 in w3_list:
        for one_address in global_test_env.genesis_config['alloc']:
            balance = w3.eth.getBalance(w3.toChecksumAddress(one_address))
            assert balance >= 0, "初始化账户错误"


@allure.title("经济模型参数，治理参数，惩罚参数，奖励参数")
@pytest.mark.P1
def test_initial_economic(global_test_env):
    """
    查看经济模型参数，治理参数，惩罚参数，奖励参数是否为正确配置的参数
    """
    log.info("查看经济模型参数，治理参数，惩罚参数，奖励参数是否为正确配置的参数")
    economic_info = global_test_env.genesis_config['EconomicModel']
    w3_list = [one_node.web3 for one_node in global_test_env.collusion_node_list]
    for w3 in w3_list:
        info = json.loads(w3.debug.economicConfig())
        assert economic_info['Common']['ExpectedMinutes'] == info['Common']['ExpectedMinutes']
        assert economic_info['Common']['ValidatorCount'] == info['Common']['ValidatorCount']
        assert economic_info['Common']['AdditionalCycleTime'] == info['Common']['AdditionalCycleTime']
        assert info['Staking']['StakeThreshold'] == economic_info['Staking']['StakeThreshold']
        assert info['Staking']['MinimumThreshold'] == economic_info['Staking']['MinimumThreshold']
        assert info['Staking']['EpochValidatorNum'] == economic_info['Staking']['EpochValidatorNum']
        assert info['Staking']['HesitateRatio'] == economic_info['Staking']['HesitateRatio']
        assert info['Staking']['UnStakeFreezeRatio'] == economic_info['Staking']['UnStakeFreezeRatio']
        assert info['Staking']['ActiveUnDelegateFreezeRatio'] == economic_info['Staking']['ActiveUnDelegateFreezeRatio']
        assert info['Slashing']['DuplicateSignReportReward'] == economic_info['Slashing']['DuplicateSignReportReward']
        assert info['Slashing']['DuplicateSignHighSlashing'] == economic_info['Slashing']['DuplicateSignHighSlashing']
        assert info['Slashing']['NumberOfBlockRewardForSlashing'] == economic_info['Slashing']['NumberOfBlockRewardForSlashing']
        assert info['Slashing']['EvidenceValidEpoch'] == economic_info['Slashing']['EvidenceValidEpoch']
        assert info['Gov']['VersionProposalVote_DurationSeconds'] == economic_info['Gov']['VersionProposalVote_DurationSeconds']
        assert info['Gov']['VersionProposalActive_ConsensusRounds'] == economic_info['Gov']['VersionProposalActive_ConsensusRounds']
        assert info['Gov']['VersionProposal_SupportRate'] == economic_info['Gov']['VersionProposal_SupportRate']
        assert info['Gov']['TextProposalVote_DurationSeconds'] == economic_info['Gov']['TextProposalVote_DurationSeconds']
        assert info['Gov']['TextProposal_VoteRate'] == economic_info['Gov']['TextProposal_VoteRate']
        assert info['Gov']['TextProposal_SupportRate'] == economic_info['Gov']['TextProposal_SupportRate']
        assert info['Gov']['CancelProposal_VoteRate'] == economic_info['Gov']['CancelProposal_VoteRate']
        assert info['Gov']['CancelProposal_SupportRate'] == economic_info['Gov']['CancelProposal_SupportRate']
        assert info['Reward']['NewBlockRate'] == economic_info['Reward']['NewBlockRate']
        assert info['Reward']['PlatONFoundationYear'] == economic_info['Reward']['PlatONFoundationYear']
        assert w3.toChecksumAddress(info['InnerAcc']['PlatONFundAccount']) == w3.toChecksumAddress(economic_info['InnerAcc']['PlatONFundAccount'])
        assert info['InnerAcc']['PlatONFundBalance'] == economic_info['InnerAcc']['PlatONFundBalance']
        assert w3.toChecksumAddress(info['InnerAcc']['CDFAccount']) == w3.toChecksumAddress(economic_info['InnerAcc']['CDFAccount'])
        assert info['InnerAcc']['CDFBalance'] == economic_info['InnerAcc']['CDFBalance']
        
   
@allure.title("基金会锁仓计划查询")
@pytest.mark.P1
def test_initial_plan(global_test_env):
    """
    查看基金会锁仓计划查询
    """
    log.info("查看基金会锁仓计划查询")
    w3_list = [one_node.web3 for one_node in global_test_env.collusion_node_list]
    for w3 in w3_list:
        info = w3.eth.call({"to": "0x1000000000000000000000000000000000000001", "data":"0xda8382100495941000000000000000000000000000000000000003"}, 0)
        recive = json.loads(str(info, encoding="ISO-8859-1"))
        plans = json.loads(recive['Data'])['plans']
        assert(8 == len(plans))
        for i in range(len(plans)):
            if 1600 == plans[i]['blockNumber']:
                assert plans[i]['amount'] == "0x2e4b34fabd7e9f0dbec800"
            if 3200 == plans[i]['blockNumber']:
                assert plans[i]['amount'] == "0x28fea17898adccbe56c800"
            if 4800 == plans[i]['blockNumber']:
                assert plans[i]['amount'] == "0x239023dfffbe286025c800"
            if 6400 == plans[i]['blockNumber']:
                assert plans[i]['amount'] == "0x1dfee323962f03fb441000"
            if 8000 == plans[i]['blockNumber']:
                assert plans[i]['amount'] == "0x184a00c8dd3215a9b02800"
            if 9600 == plans[i]['blockNumber']:
                assert plans[i]['amount'] == "0x127098c579295b871ef800"
            if 11200 == plans[i]['blockNumber']:
                assert plans[i]['amount'] == "0xc71c15b99099294962800"
            if 12800 == plans[i]['blockNumber']:
                assert plans[i]['amount'] == "0x64c8af579b5c2dc39e800"
                

@allure.title("共识参数")
@pytest.mark.P1
def test_initial_consensus(global_test_env):
    """
    查看共识的每个共识节点的出块个数
    """
    log.info("查看共识的每个共识节点的出块个数 和 总共的共识节点的个数")
    amount = global_test_env.genesis_config['config']['cbft']['amount']
    w3_list = [one_node.web3 for one_node in global_test_env.collusion_node_list]
    for w3 in w3_list:
        info = w3.eth.getPrepareQC(amount)
        assert info['viewNumber'] == 0
        info = w3.eth.getPrepareQC(amount+1)
        assert info['viewNumber'] == 1
