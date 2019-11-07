import pytest
from common.log import log
import time
from tests.lib.utils import assert_code, wait_block_number, upload_platon
from tests.lib.client import get_client_obj
from tests.govern.conftest import version_proposal_vote, get_refund_to_account_block
from tests.lib import Genesis
from dacite import from_dict
from tests.govern.conftest import param_proposal_vote

class TestVotingStatistics():

    def test_VS_EP_004(self, new_genesis_env, client_con_list_obj, client_noc_list_obj):
        genesis = from_dict(data_class=Genesis, data=new_genesis_env.genesis_config)
        genesis.EconomicModel.Gov.ParamProposalVote_DurationSeconds = 0
        new_genesis_env.set_genesis(genesis.to_dict())
        new_genesis_env.deploy_all()
        pip_obj = client_con_list_obj[0].pip
        result = pip_obj.submitParam(pip_obj.node.node_id, str(time.time()), 'Slashing', 'SlashBlocksReward', '83',
                            pip_obj.node.staking_address, transaction_cfg=pip_obj.cfg.transaction_cfg)
        log.info('Submit param proposal result : {}'.format(result))
        assert_code(result, 0)
        proposalinfo = pip_obj.get_effect_proposal_info_of_vote(pip_obj.cfg.param_proposal)
        log.info('Param proposal info {}'.format(proposalinfo))
        for index in range(len(client_con_list_obj) - 1):
            pip_obj = client_con_list_obj[index].pip
            result = param_proposal_vote(pip_obj, index+1)
            assert_code(result, 0)
            address, _ = client_noc_list_obj[index].economic.account.generate_account(client_noc_list_obj[index].node.web3,
                                                                                      10**18 * 10000000)
            result = client_noc_list_obj[index].staking.create_staking(0, address, address,
                                                                       transaction_cfg=pip_obj.cfg.transaction_cfg)
            log.info('Node {} staking result : {}'.format(client_noc_list_obj[index].node.node_id, result))
            assert_code(result, 0)
        pip_obj.economic.wait_settlement_blocknum(pip_obj.node)
        result = pip_obj.get_accuverifiers_count(proposalinfo.get('ProposalID'))
        log.info('Get proposal vote infomation {}'.format(result))
        assert result[3] == 7

    def test_VS_EP_005(self, new_genesis_env, client_con_list_obj, client_noc_list_obj):
        genesis = from_dict(data_class=Genesis, data=new_genesis_env.genesis_config)
        genesis.EconomicModel.Gov.ParamProposalVote_DurationSeconds = 160
        new_genesis_env.set_genesis(genesis.to_dict())
        new_genesis_env.deploy_all()
        pip_obj = client_con_list_obj[0].pip
        result = pip_obj.submitParam(pip_obj.node.node_id, str(time.time()), 'Slashing', 'SlashBlocksReward', '83',
                            pip_obj.node.staking_address, transaction_cfg=pip_obj.cfg.transaction_cfg)
        log.info('Submit param proposal result : {}'.format(result))
        assert_code(result, 0)
        proposalinfo = pip_obj.get_effect_proposal_info_of_vote(pip_obj.cfg.param_proposal)
        log.info('Param proposal info {}'.format(proposalinfo))
        for index in range(len(client_con_list_obj) - 2):
            pip_obj = client_con_list_obj[index].pip
            result = param_proposal_vote(pip_obj, index+1)
            assert_code(result, 0)
            address, _ = client_noc_list_obj[index].economic.account.generate_account(client_noc_list_obj[index].node.web3,
                                                                                      10**18 * 10000000)
            result = client_noc_list_obj[index].staking.create_staking(0, address, address,
                                                                       transaction_cfg=pip_obj.cfg.transaction_cfg)
            log.info('Node {} staking result : {}'.format(client_noc_list_obj[index].node.node_id, result))
            assert_code(result, 0)
        pip_obj.economic.wait_settlement_blocknum(pip_obj.node)

        result = pip_obj.get_accuverifiers_count(proposalinfo.get('ProposalID'))
        log.info('Get proposal vote infomation {}'.format(result))