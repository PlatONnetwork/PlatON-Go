from common.log import log
from tests.lib.utils import upload_platon, assert_code, wait_block_number, get_pledge_list
import pytest, time
from tests.govern.test_voting_statistics import submitvpandvote

def replace_version_declare(pip_obj, platon_bin, versiontag):
    upload_platon(pip_obj.node, platon_bin)
    log.info('Replace the platon of the node {} version{}'.format(pip_obj.node.node_id, versiontag))
    pip_obj.node.restart()
    log.info('Restart the node{}'.format(pip_obj.node.node_id))
    assert pip_obj.node.program_version == versiontag
    log.info('assert the version of the node is {}'.format(versiontag))
    log.info("staking: {}".format(pip_obj.node.staking_address))
    log.info("account:{}".format(pip_obj.economic.account.accounts))
    result = pip_obj.declareVersion(pip_obj.node.node_id, pip_obj.node.staking_address,
                                    transaction_cfg=pip_obj.cfg.transaction_cfg)
    log.info('declareversion {} result: {}'.format(pip_obj.node.program_version, result))
    return result

def wrong_verisonsign_declare(pip_obj, pip_obj_test):
    result = pip_obj.declareVersion(pip_obj.node.node_id, pip_obj.node.staking_address,
                                    version_sign=pip_obj_test.node.program_version_sign,
                                    transaction_cfg=pip_obj.cfg.transaction_cfg)
    log.info('wrong program version sign, declareVersion result : {}'.format(result))
    return result

def wrong_verison_declare(pip_obj, version=None):
    if not version:
        proposalinfo = pip_obj.get_effect_proposal_info_of_vote()
        version = proposalinfo.get('NewVersion')
        log.info('The new version of the proposal: {}'.format(version))
    result = pip_obj.declareVersion(pip_obj.node.node_id, pip_obj.node.staking_address,
                                    program_version=version,
                                    transaction_cfg=pip_obj.cfg.transaction_cfg)
    log.info('wrong program version, declareVersion: {} result : {}'.format(version, result))
    return result

def test_DE_DE_001(client_verifier_obj):
    pip_obj = client_verifier_obj.pip
    address, _ = pip_obj.economic.account.generate_account(pip_obj.node.web3, 10**18 * 10000)
    result = pip_obj.declareVersion(pip_obj.node.node_id, address, transaction_cfg=pip_obj.cfg.transaction_cfg)
    log.info('declareVersion result: {}'.format(result))
    assert_code(result, 302021)

class TestNoProposalVE():
    @pytest.mark.compatibility
    def test_DE_VE_001(self, noproposal_pipobj_list):
        pip_obj = noproposal_pipobj_list[0]
        result = replace_version_declare(pip_obj, pip_obj.cfg.PLATON_NEW_BIN2, pip_obj.cfg.version2)
        assert_code(result, 0)

        result = wrong_verisonsign_declare(pip_obj, noproposal_pipobj_list[1])
        assert_code(result, 302024)

        result = wrong_verison_declare(pip_obj, pip_obj.chain_version)
        assert_code(result, 302024)

    def test_DE_VE_002(self, noproposal_pipobj_list):
        pip_obj = noproposal_pipobj_list[0]
        result = replace_version_declare(pip_obj, pip_obj.cfg.PLATON_NEW_BIN1, pip_obj.cfg.version1)
        assert_code(result, 302028)

        result = wrong_verisonsign_declare(pip_obj, noproposal_pipobj_list[1])
        assert_code(result, 302024)

        result = wrong_verison_declare(pip_obj, pip_obj.chain_version)
        assert_code(result, 302024)

    @pytest.mark.P0
    def test_DE_VE_004(self, noproposal_pipobj_list):
        pip_obj = noproposal_pipobj_list[0]
        result = replace_version_declare(pip_obj, pip_obj.cfg.PLATON_NEW_BIN3, pip_obj.cfg.version3)
        assert_code(result, 0)

        result = wrong_verisonsign_declare(pip_obj, noproposal_pipobj_list[1])
        assert_code(result, 302024)

        result = wrong_verison_declare(pip_obj, pip_obj.chain_version)
        assert_code(result, 302024)

    def test_DE_VE_005(self, noproposal_pipobj_list):
        pip_obj = noproposal_pipobj_list[0]
        result = replace_version_declare(pip_obj, pip_obj.cfg.PLATON_NEW_BIN, pip_obj.cfg.version5)
        assert_code(result, 302028)

        result = wrong_verisonsign_declare(pip_obj, noproposal_pipobj_list[1])
        assert_code(result, 302024)

        result = wrong_verison_declare(pip_obj, pip_obj.chain_version)
        assert_code(result, 302024)

    def test_DE_VE_006(self, noproposal_pipobj_list):
        pip_obj = noproposal_pipobj_list[0]
        result = replace_version_declare(pip_obj, pip_obj.cfg.PLATON_NEW_BIN8, pip_obj.cfg.version8)
        assert_code(result, 302028)

        result = wrong_verisonsign_declare(pip_obj, noproposal_pipobj_list[1])
        assert_code(result, 302024)

        result = wrong_verison_declare(pip_obj, pip_obj.chain_version)
        assert_code(result, 302024)

    def test_DE_VE_007(self, noproposal_pipobj_list):
        pip_obj = noproposal_pipobj_list[0]
        result = replace_version_declare(pip_obj, pip_obj.cfg.PLATON_NEW_BIN0, pip_obj.cfg.version0)
        assert_code(result, 0)

        result = wrong_verisonsign_declare(pip_obj, noproposal_pipobj_list[1])
        assert_code(result, 302024)

        result = wrong_verison_declare(pip_obj, pip_obj.cfg.version3)
        assert_code(result, 302024)

        result = wrong_verison_declare(pip_obj, pip_obj.cfg.version2)
        assert_code(result, 302024)

class TestVotingProposalVE():
    def test_DE_VE_008(self, proposal_pipobj_list):
        pip_obj = proposal_pipobj_list[0]
        result = replace_version_declare(pip_obj, pip_obj.cfg.PLATON_NEW_BIN2, pip_obj.cfg.version2)
        assert_code(result, 0)

        result = wrong_verisonsign_declare(pip_obj, proposal_pipobj_list[1])
        assert_code(result, 302024)

        result = wrong_verison_declare(pip_obj, pip_obj.chain_version)
        assert_code(result, 302024)

        result = wrong_verison_declare(pip_obj)
        assert_code(result, 302024)

    def test_DE_VE_010(self, proposal_pipobj_list):
        pip_obj = proposal_pipobj_list[0]
        result = replace_version_declare(pip_obj, pip_obj.cfg.PLATON_NEW_BIN1, pip_obj.cfg.version1)
        assert_code(result, 302028)

        result = wrong_verisonsign_declare(pip_obj, proposal_pipobj_list[1])
        assert_code(result, 302024)

        result = wrong_verison_declare(pip_obj, pip_obj.chain_version)
        assert_code(result, 302024)

        result = wrong_verison_declare(pip_obj)
        assert_code(result, 302024)

    def test_DE_VE_014(self, bv_proposal_pipobj_list):
        pip_obj = bv_proposal_pipobj_list[0]
        result = replace_version_declare(pip_obj, pip_obj.cfg.PLATON_NEW_BIN2, pip_obj.cfg.version2)
        assert_code(result, 0)

        result = wrong_verisonsign_declare(pip_obj, bv_proposal_pipobj_list[1])
        assert_code(result, 302024)

        result = wrong_verison_declare(pip_obj, pip_obj.chain_version)
        assert_code(result, 302024)

        result = wrong_verison_declare(pip_obj)
        assert_code(result, 302024)

    def test_DE_VE_025(self, bv_proposal_pipobj_list):
        pip_obj = bv_proposal_pipobj_list[0]
        result = replace_version_declare(pip_obj, pip_obj.cfg.PLATON_NEW_BIN1, pip_obj.cfg.version1)
        assert_code(result, 302028)

        result = wrong_verisonsign_declare(pip_obj, bv_proposal_pipobj_list[1])
        assert_code(result, 302024)

        result = wrong_verison_declare(pip_obj, pip_obj.chain_version)
        assert_code(result, 302024)

        result = wrong_verison_declare(pip_obj)
        assert_code(result, 302024)

    def test_DE_VE_032(self, proposal_pipobj_list):
        pip_obj = proposal_pipobj_list[0]
        result = replace_version_declare(pip_obj, pip_obj.cfg.PLATON_NEW_BIN0, pip_obj.cfg.version0)
        assert_code(result, 0)

        result = wrong_verisonsign_declare(pip_obj, proposal_pipobj_list[1])
        assert_code(result, 302024)

        result = wrong_verison_declare(pip_obj)
        assert_code(result, 302024)

    def test_DE_VE_034(self, bv_proposal_pipobj_list):
        pip_obj = bv_proposal_pipobj_list[0]
        result = replace_version_declare(pip_obj, pip_obj.cfg.PLATON_NEW_BIN0, pip_obj.cfg.version0)
        assert_code(result, 0)

        result = wrong_verisonsign_declare(pip_obj, bv_proposal_pipobj_list[1])
        assert_code(result, 302024)

        result = wrong_verison_declare(pip_obj)
        assert_code(result, 302024)

    def test_DE_VE_036(self, proposal_pipobj_list):
        pip_obj = proposal_pipobj_list[0]
        result = replace_version_declare(pip_obj, pip_obj.cfg.PLATON_NEW_BIN3, pip_obj.cfg.version3)
        assert_code(result, 0)

        result = wrong_verisonsign_declare(pip_obj, proposal_pipobj_list[1])
        assert_code(result, 302024)

        result = wrong_verison_declare(pip_obj)
        assert_code(result, 302024)

        result = wrong_verison_declare(pip_obj, pip_obj.chain_version)
        assert_code(result, 302024)

    def test_DE_VE_038(self, proposal_pipobj_list):
        pip_obj = proposal_pipobj_list[0]
        result = replace_version_declare(pip_obj, pip_obj.cfg.PLATON_NEW_BIN3, pip_obj.cfg.version3)
        assert_code(result, 0)

        result = wrong_verisonsign_declare(pip_obj, proposal_pipobj_list[1])
        assert_code(result, 302024)

        result = wrong_verison_declare(pip_obj, pip_obj.chain_version)
        assert_code(result, 302024)

    def test_DE_VE_040(self, proposal_pipobj_list):
        pip_obj = proposal_pipobj_list[0]
        result = replace_version_declare(pip_obj, pip_obj.cfg.PLATON_NEW_BIN4, pip_obj.cfg.version4)
        assert_code(result, 0)

        result = wrong_verisonsign_declare(pip_obj, proposal_pipobj_list[1])
        assert_code(result, 302024)

        result = wrong_verison_declare(pip_obj)
        assert_code(result, 302024)

        result = wrong_verison_declare(pip_obj, pip_obj.chain_version)
        assert_code(result, 302024)

    def test_DE_VE_042(self, proposal_pipobj_list):
        pip_obj = proposal_pipobj_list[0]
        result = replace_version_declare(pip_obj, pip_obj.cfg.PLATON_NEW_BIN6, pip_obj.cfg.version6)
        assert_code(result, 0)

        result = wrong_verisonsign_declare(pip_obj, proposal_pipobj_list[1])
        assert_code(result, 302024)

        result = wrong_verison_declare(pip_obj)
        assert_code(result, 302024)

        result = wrong_verison_declare(pip_obj, pip_obj.chain_version)
        assert_code(result, 302024)

    def test_DE_VE_044(self, proposal_pipobj_list):
        pip_obj = proposal_pipobj_list[0]
        result = replace_version_declare(pip_obj, pip_obj.cfg.PLATON_NEW_BIN8, pip_obj.cfg.version8)
        assert_code(result, 302028)

        result = wrong_verisonsign_declare(pip_obj, proposal_pipobj_list[1])
        assert_code(result, 302024)

        result = wrong_verison_declare(pip_obj)
        assert_code(result, 302024)

        result = wrong_verison_declare(pip_obj, pip_obj.chain_version)
        assert_code(result, 302024)

    def test_DE_VE_046(self, bv_proposal_pipobj_list):
        pip_obj = bv_proposal_pipobj_list[0]
        result = replace_version_declare(pip_obj, pip_obj.cfg.PLATON_NEW_BIN3, pip_obj.cfg.version3)
        assert_code(result, 0)

        result = wrong_verisonsign_declare(pip_obj, bv_proposal_pipobj_list[1])
        assert_code(result, 302024)

        result = wrong_verison_declare(pip_obj)
        assert_code(result, 302024)

        result = wrong_verison_declare(pip_obj, pip_obj.chain_version)
        assert_code(result, 302024)

    def test_DE_VE_048(self, bv_proposal_pipobj_list):
        pip_obj = bv_proposal_pipobj_list[0]
        result = replace_version_declare(pip_obj, pip_obj.cfg.PLATON_NEW_BIN, pip_obj.cfg.version5)
        assert_code(result, 302028)

        result = wrong_verisonsign_declare(pip_obj, bv_proposal_pipobj_list[1])
        assert_code(result, 302024)

        result = wrong_verison_declare(pip_obj)
        assert_code(result, 302024)

        result = wrong_verison_declare(pip_obj, pip_obj.chain_version)
        assert_code(result, 302024)

    def test_DE_VE_050(self, bv_proposal_pipobj_list):
        pip_obj = bv_proposal_pipobj_list[0]
        result = replace_version_declare(pip_obj, pip_obj.cfg.PLATON_NEW_BIN4, pip_obj.cfg.version4)
        assert_code(result, 302028)

        result = wrong_verisonsign_declare(pip_obj, bv_proposal_pipobj_list[1])
        assert_code(result, 302024)

        result = wrong_verison_declare(pip_obj)
        assert_code(result, 302024)

        result = wrong_verison_declare(pip_obj, pip_obj.chain_version)
        assert_code(result, 302024)

    def test_DE_VE_052(self, bv_proposal_pipobj_list):
        pip_obj = bv_proposal_pipobj_list[0]
        result = replace_version_declare(pip_obj, pip_obj.cfg.PLATON_NEW_BIN6, pip_obj.cfg.version6)
        assert_code(result, 302028)

        result = wrong_verisonsign_declare(pip_obj, bv_proposal_pipobj_list[1])
        assert_code(result, 302024)

        result = wrong_verison_declare(pip_obj)
        assert_code(result, 302024)

        result = wrong_verison_declare(pip_obj, pip_obj.chain_version)
        assert_code(result, 302024)

    def test_DE_VE_054(self, bv_proposal_pipobj_list):
        pip_obj = bv_proposal_pipobj_list[0]
        result = replace_version_declare(pip_obj, pip_obj.cfg.PLATON_NEW_BIN8, pip_obj.cfg.version8)
        assert_code(result, 0)

        result = wrong_verisonsign_declare(pip_obj, bv_proposal_pipobj_list[1])
        assert_code(result, 302024)

        result = wrong_verison_declare(pip_obj, pip_obj.chain_version)
        assert_code(result, 302024)

class TestVotingProposlaVotedVE():
    def test_DE_VE_009(self, proposal_voted_pipobj_list):
        pip_obj = proposal_voted_pipobj_list[0]
        result = replace_version_declare(pip_obj, pip_obj.cfg.PLATON_NEW_BIN2, pip_obj.cfg.version2)
        assert_code(result, 302028)

        result = wrong_verisonsign_declare(pip_obj, proposal_voted_pipobj_list[1])
        assert_code(result, 302024)

        result = wrong_verison_declare(pip_obj)
        assert_code(result, 302024)

        result = wrong_verison_declare(pip_obj, pip_obj.chain_version)
        assert_code(result, 302024)

    def test_DE_VE_011(self, proposal_voted_pipobj_list):
        pip_obj = proposal_voted_pipobj_list[0]
        result = replace_version_declare(pip_obj, pip_obj.cfg.PLATON_NEW_BIN1, pip_obj.cfg.version1)
        assert_code(result, 302028)

        result = wrong_verisonsign_declare(pip_obj, proposal_voted_pipobj_list[1])
        assert_code(result, 302024)

        result = wrong_verison_declare(pip_obj)
        assert_code(result, 302024)

        result = wrong_verison_declare(pip_obj, pip_obj.chain_version)
        assert_code(result, 302024)

    def test_DE_VE_021(self, bv_proposal_voted_pipobj_list):
        pip_obj = bv_proposal_voted_pipobj_list[0]
        result = replace_version_declare(pip_obj, pip_obj.cfg.PLATON_NEW_BIN2, pip_obj.cfg.version2)
        assert_code(result, 302028)

        result = wrong_verisonsign_declare(pip_obj, bv_proposal_voted_pipobj_list[1])
        assert_code(result, 302024)

        result = wrong_verison_declare(pip_obj)
        assert_code(result, 302024)

        result = wrong_verison_declare(pip_obj, pip_obj.chain_version)
        assert_code(result, 302024)

    def test_DE_VE_026(self, bv_proposal_voted_pipobj_list):
        pip_obj = bv_proposal_voted_pipobj_list[0]
        result = replace_version_declare(pip_obj, pip_obj.cfg.PLATON_NEW_BIN1, pip_obj.cfg.version1)
        assert_code(result, 302028)

        result = wrong_verisonsign_declare(pip_obj, bv_proposal_voted_pipobj_list[1])
        assert_code(result, 302024)

        result = wrong_verison_declare(pip_obj)
        assert_code(result, 302024)

        result = wrong_verison_declare(pip_obj, pip_obj.chain_version)
        assert_code(result, 302024)

    def test_DE_VE_033(self, proposal_voted_pipobj_list):
        pip_obj = proposal_voted_pipobj_list[0]
        result = replace_version_declare(pip_obj, pip_obj.cfg.PLATON_NEW_BIN0, pip_obj.cfg.version0)
        assert_code(result, 302028)

        result = wrong_verisonsign_declare(pip_obj, proposal_voted_pipobj_list[1])
        assert_code(result, 302024)

        result = wrong_verison_declare(pip_obj)
        assert_code(result, 302024)

    def test_DE_VE_035(self, bv_proposal_voted_pipobj_list):
        pip_obj = bv_proposal_voted_pipobj_list[0]
        result = replace_version_declare(pip_obj, pip_obj.cfg.PLATON_NEW_BIN0, pip_obj.cfg.version0)
        assert_code(result, 302028)

        result = wrong_verisonsign_declare(pip_obj, bv_proposal_voted_pipobj_list[1])
        assert_code(result, 302024)

        result = wrong_verison_declare(pip_obj)
        assert_code(result, 302024)

    def test_DE_VE_037(self, proposal_voted_pipobj_list):
        pip_obj = proposal_voted_pipobj_list[0]
        result = replace_version_declare(pip_obj, pip_obj.cfg.PLATON_NEW_BIN3, pip_obj.cfg.version3)
        assert_code(result, 302028)

        result = wrong_verisonsign_declare(pip_obj, proposal_voted_pipobj_list[1])
        assert_code(result, 302024)

        result = wrong_verison_declare(pip_obj)
        assert_code(result, 302024)

        result = wrong_verison_declare(pip_obj, pip_obj.chain_version)
        assert_code(result, 302024)

    def test_DE_VE_039(self, proposal_voted_pipobj_list):
        pip_obj = proposal_voted_pipobj_list[0]
        result = replace_version_declare(pip_obj, pip_obj.cfg.PLATON_NEW_BIN, pip_obj.cfg.version5)
        assert_code(result, 0)

        result = wrong_verisonsign_declare(pip_obj, proposal_voted_pipobj_list[1])
        assert_code(result, 302024)

        result = wrong_verison_declare(pip_obj, pip_obj.cfg.version4)
        assert_code(result, 302024)

        result = wrong_verison_declare(pip_obj, pip_obj.chain_version)
        assert_code(result, 302024)

    def test_DE_VE_041(self, proposal_voted_pipobj_list):
        pip_obj = proposal_voted_pipobj_list[0]
        result = replace_version_declare(pip_obj, pip_obj.cfg.PLATON_NEW_BIN4, pip_obj.cfg.version4)
        assert_code(result, 0)

        result = wrong_verisonsign_declare(pip_obj, proposal_voted_pipobj_list[1])
        assert_code(result, 302024)

        result = wrong_verison_declare(pip_obj, pip_obj.cfg.version5)
        assert_code(result, 302024)

        result = wrong_verison_declare(pip_obj, pip_obj.chain_version)
        assert_code(result, 302024)

    def test_DE_VE_043(self, proposal_voted_pipobj_list):
        pip_obj = proposal_voted_pipobj_list[0]
        result = replace_version_declare(pip_obj, pip_obj.cfg.PLATON_NEW_BIN6, pip_obj.cfg.version6)
        assert_code(result, 0)

        result = wrong_verisonsign_declare(pip_obj, proposal_voted_pipobj_list[1])
        assert_code(result, 302024)

        result = wrong_verison_declare(pip_obj, pip_obj.cfg.version5)
        assert_code(result, 302024)

        result = wrong_verison_declare(pip_obj, pip_obj.chain_version)
        assert_code(result, 302024)

    def test_DE_VE_045(self, proposal_voted_pipobj_list):
        pip_obj = proposal_voted_pipobj_list[0]
        result = replace_version_declare(pip_obj, pip_obj.cfg.PLATON_NEW_BIN8, pip_obj.cfg.version8)
        assert_code(result, 302028)

        result = wrong_verisonsign_declare(pip_obj, proposal_voted_pipobj_list[1])
        assert_code(result, 302024)

        result = wrong_verison_declare(pip_obj, pip_obj.cfg.version5)
        assert_code(result, 302024)

        result = wrong_verison_declare(pip_obj, pip_obj.chain_version)
        assert_code(result, 302024)

    def test_DE_VE_047(self, bv_proposal_voted_pipobj_list):
        pip_obj = bv_proposal_voted_pipobj_list[0]
        result = replace_version_declare(pip_obj, pip_obj.cfg.PLATON_NEW_BIN3, pip_obj.cfg.version3)
        assert_code(result, 302028)

        result = wrong_verisonsign_declare(pip_obj, bv_proposal_voted_pipobj_list[1])
        assert_code(result, 302024)

        result = wrong_verison_declare(pip_obj, pip_obj.cfg.version5)
        assert_code(result, 302024)

        result = wrong_verison_declare(pip_obj, pip_obj.chain_version)
        assert_code(result, 302024)

    def test_DE_VE_049(self, bv_proposal_voted_pipobj_list):
        pip_obj = bv_proposal_voted_pipobj_list[0]
        result = replace_version_declare(pip_obj, pip_obj.cfg.PLATON_NEW_BIN, pip_obj.cfg.version5)
        assert_code(result, 302028)

        result = wrong_verisonsign_declare(pip_obj, bv_proposal_voted_pipobj_list[1])
        assert_code(result, 302024)

        result = wrong_verison_declare(pip_obj, pip_obj.cfg.version8)
        assert_code(result, 302024)

        result = wrong_verison_declare(pip_obj, pip_obj.chain_version)
        assert_code(result, 302024)

    def test_DE_VE_051(self, bv_proposal_voted_pipobj_list):
        pip_obj = bv_proposal_voted_pipobj_list[0]
        result = replace_version_declare(pip_obj, pip_obj.cfg.PLATON_NEW_BIN4, pip_obj.cfg.version4)
        assert_code(result, 302028)

        result = wrong_verisonsign_declare(pip_obj, bv_proposal_voted_pipobj_list[1])
        assert_code(result, 302024)

        result = wrong_verison_declare(pip_obj, pip_obj.cfg.version8)
        assert_code(result, 302024)

        result = wrong_verison_declare(pip_obj, pip_obj.chain_version)
        assert_code(result, 302024)

    def test_DE_VE_053(self, bv_proposal_voted_pipobj_list):
        pip_obj = bv_proposal_voted_pipobj_list[0]
        result = replace_version_declare(pip_obj, pip_obj.cfg.PLATON_NEW_BIN6, pip_obj.cfg.version6)
        assert_code(result, 302028)

        result = wrong_verisonsign_declare(pip_obj, bv_proposal_voted_pipobj_list[1])
        assert_code(result, 302024)

        result = wrong_verison_declare(pip_obj, pip_obj.cfg.version8)
        assert_code(result, 302024)

        result = wrong_verison_declare(pip_obj, pip_obj.chain_version)
        assert_code(result, 302024)

    def test_DE_VE_055(self, bv_proposal_voted_pipobj_list):
        pip_obj = bv_proposal_voted_pipobj_list[0]
        result = replace_version_declare(pip_obj, pip_obj.cfg.PLATON_NEW_BIN8, pip_obj.cfg.version8)
        assert_code(result, 0)

        result = wrong_verisonsign_declare(pip_obj, bv_proposal_voted_pipobj_list[1])
        assert_code(result, 302024)

        result = wrong_verison_declare(pip_obj, pip_obj.chain_version)
        assert_code(result, 302024)

class TestPreactiveProposalVE():
    def test_DE_VE_056(self, preactive_proposal_pipobj_list):
        pip_obj = preactive_proposal_pipobj_list[0]
        result = replace_version_declare(pip_obj, pip_obj.cfg.PLATON_NEW_BIN2, pip_obj.cfg.version2)
        assert_code(result, 302028)

        result = wrong_verisonsign_declare(pip_obj, preactive_proposal_pipobj_list[1])
        assert_code(result, 302024)

        result = wrong_verison_declare(pip_obj, pip_obj.cfg.version5)
        assert_code(result, 302024)

        result = wrong_verison_declare(pip_obj, pip_obj.chain_version)
        assert_code(result, 302024)

    def test_DE_VE_057(self, preactive_proposal_pipobj_list):
        pip_obj = preactive_proposal_pipobj_list[0]
        result = replace_version_declare(pip_obj, pip_obj.cfg.PLATON_NEW_BIN1, pip_obj.cfg.version1)
        assert_code(result, 302028)

        result = wrong_verisonsign_declare(pip_obj, preactive_proposal_pipobj_list[1])
        assert_code(result, 302024)

        result = wrong_verison_declare(pip_obj, pip_obj.cfg.version5)
        assert_code(result, 302024)

        result = wrong_verison_declare(pip_obj, pip_obj.chain_version)
        assert_code(result, 302024)

    def test_DE_VE_059(self, preactive_proposal_pipobj_list):
        pip_obj = preactive_proposal_pipobj_list[0]
        result = replace_version_declare(pip_obj, pip_obj.cfg.PLATON_NEW_BIN2, pip_obj.cfg.version2)
        assert_code(result, 302028)

        result = wrong_verisonsign_declare(pip_obj, preactive_proposal_pipobj_list[1])
        assert_code(result, 302024)

        result = wrong_verison_declare(pip_obj, pip_obj.cfg.version5)
        assert_code(result, 302024)

        result = wrong_verison_declare(pip_obj, pip_obj.chain_version)
        assert_code(result, 302024)

    def test_DE_VE_060(self, preactive_bv_proposal_pipobj_list):
        pip_obj = preactive_bv_proposal_pipobj_list[0]
        result = replace_version_declare(pip_obj, pip_obj.cfg.PLATON_NEW_BIN1, pip_obj.cfg.version1)
        assert_code(result, 302028)

        result = wrong_verisonsign_declare(pip_obj, preactive_bv_proposal_pipobj_list[1])
        assert_code(result, 302024)

        result = wrong_verison_declare(pip_obj, pip_obj.cfg.version8)
        assert_code(result, 302024)

        result = wrong_verison_declare(pip_obj, pip_obj.chain_version)
        assert_code(result, 302024)

    def test_DE_VE_062(self, preactive_proposal_pipobj_list):
        pip_obj = preactive_proposal_pipobj_list[0]
        result = replace_version_declare(pip_obj, pip_obj.cfg.PLATON_NEW_BIN0, pip_obj.cfg.version0)
        assert_code(result, 302028)

        result = wrong_verisonsign_declare(pip_obj, preactive_proposal_pipobj_list[1])
        assert_code(result, 302024)

        result = wrong_verison_declare(pip_obj, pip_obj.cfg.version5)
        assert_code(result, 302024)

    def test_DE_VE_063(self, preactive_bv_proposal_pipobj_list):
        pip_obj = preactive_bv_proposal_pipobj_list[0]
        result = replace_version_declare(pip_obj, pip_obj.cfg.PLATON_NEW_BIN0, pip_obj.cfg.version0)
        assert_code(result, 302028)

        result = wrong_verisonsign_declare(pip_obj, preactive_bv_proposal_pipobj_list[1])
        assert_code(result, 302024)

        result = wrong_verison_declare(pip_obj, pip_obj.cfg.version8)
        assert_code(result, 302024)

    def test_DE_VE_064(self, preactive_proposal_pipobj_list):
        pip_obj = preactive_proposal_pipobj_list[0]
        result = replace_version_declare(pip_obj, pip_obj.cfg.PLATON_NEW_BIN3, pip_obj.cfg.version3)
        assert_code(result, 302028)

        result = wrong_verisonsign_declare(pip_obj, preactive_proposal_pipobj_list[1])
        assert_code(result, 302024)

        result = wrong_verison_declare(pip_obj, pip_obj.cfg.version5)
        assert_code(result, 302024)

        result = wrong_verison_declare(pip_obj, pip_obj.chain_version)
        assert_code(result, 302024)

    def test_DE_VE_065(self, preactive_proposal_pipobj_list):
        pip_obj = preactive_proposal_pipobj_list[0]
        result = replace_version_declare(pip_obj, pip_obj.cfg.PLATON_NEW_BIN, pip_obj.cfg.version5)
        assert_code(result, 0)

        result = wrong_verisonsign_declare(pip_obj, preactive_proposal_pipobj_list[1])
        assert_code(result, 302024)

        result = wrong_verison_declare(pip_obj, pip_obj.cfg.version0)
        assert_code(result, 302024)

    def test_DE_VE_066(self, preactive_proposal_pipobj_list):
        pip_obj = preactive_proposal_pipobj_list[0]
        result = replace_version_declare(pip_obj, pip_obj.cfg.PLATON_NEW_BIN4, pip_obj.cfg.version4)
        assert_code(result, 0)

        result = wrong_verisonsign_declare(pip_obj, preactive_proposal_pipobj_list[1])
        assert_code(result, 302024)

        result = wrong_verison_declare(pip_obj, pip_obj.cfg.version5)
        assert_code(result, 302024)

        result = wrong_verison_declare(pip_obj, pip_obj.chain_version)
        assert_code(result, 302024)

    def test_DE_VE_067(self, preactive_proposal_pipobj_list):
        pip_obj = preactive_proposal_pipobj_list[0]
        result = replace_version_declare(pip_obj, pip_obj.cfg.PLATON_NEW_BIN6, pip_obj.cfg.version6)
        assert_code(result, 0)

        result = wrong_verisonsign_declare(pip_obj, preactive_proposal_pipobj_list[1])
        assert_code(result, 302024)

        result = wrong_verison_declare(pip_obj, pip_obj.cfg.version5)
        assert_code(result, 302024)

        result = wrong_verison_declare(pip_obj, pip_obj.chain_version)
        assert_code(result, 302024)

    def test_DE_VE_068(self, preactive_proposal_pipobj_list):
        pip_obj = preactive_proposal_pipobj_list[0]
        result = replace_version_declare(pip_obj, pip_obj.cfg.PLATON_NEW_BIN8, pip_obj.cfg.version8)
        assert_code(result, 302028)

        result = wrong_verisonsign_declare(pip_obj, preactive_proposal_pipobj_list[1])
        assert_code(result, 302024)

        result = wrong_verison_declare(pip_obj, pip_obj.cfg.version5)
        assert_code(result, 302024)

        result = wrong_verison_declare(pip_obj, pip_obj.chain_version)
        assert_code(result, 302024)

    def test_DE_VE_069(self, preactive_bv_proposal_pipobj_list):
        pip_obj = preactive_bv_proposal_pipobj_list[0]
        result = replace_version_declare(pip_obj, pip_obj.cfg.PLATON_NEW_BIN3, pip_obj.cfg.version3)
        assert_code(result, 302028)

        result = wrong_verisonsign_declare(pip_obj, preactive_bv_proposal_pipobj_list[1])
        assert_code(result, 302024)

        result = wrong_verison_declare(pip_obj, pip_obj.cfg.version8)
        assert_code(result, 302024)

        result = wrong_verison_declare(pip_obj, pip_obj.chain_version)
        assert_code(result, 302024)

    def test_DE_VE_070(self, preactive_bv_proposal_pipobj_list):
        pip_obj = preactive_bv_proposal_pipobj_list[0]
        result = replace_version_declare(pip_obj, pip_obj.cfg.PLATON_NEW_BIN, pip_obj.cfg.version5)
        assert_code(result, 302028)

        result = wrong_verisonsign_declare(pip_obj, preactive_bv_proposal_pipobj_list[1])
        assert_code(result, 302024)

        result = wrong_verison_declare(pip_obj, pip_obj.cfg.version8)
        assert_code(result, 302024)

        result = wrong_verison_declare(pip_obj, pip_obj.chain_version)
        assert_code(result, 302024)

    def test_DE_VE_071(self, preactive_bv_proposal_pipobj_list):
        pip_obj = preactive_bv_proposal_pipobj_list[0]
        result = replace_version_declare(pip_obj, pip_obj.cfg.PLATON_NEW_BIN4, pip_obj.cfg.version4)
        assert_code(result, 302028)

        result = wrong_verisonsign_declare(pip_obj, preactive_bv_proposal_pipobj_list[1])
        assert_code(result, 302024)

        result = wrong_verison_declare(pip_obj, pip_obj.cfg.version8)
        assert_code(result, 302024)

        result = wrong_verison_declare(pip_obj, pip_obj.chain_version)
        assert_code(result, 302024)

    def test_DE_VE_072(self, preactive_bv_proposal_pipobj_list):
        pip_obj = preactive_bv_proposal_pipobj_list[0]
        result = replace_version_declare(pip_obj, pip_obj.cfg.PLATON_NEW_BIN6, pip_obj.cfg.version6)
        assert_code(result, 302028)

        result = wrong_verisonsign_declare(pip_obj, preactive_bv_proposal_pipobj_list[1])
        assert_code(result, 302024)

        result = wrong_verison_declare(pip_obj, pip_obj.cfg.version8)
        assert_code(result, 302024)

        result = wrong_verison_declare(pip_obj, pip_obj.chain_version)
        assert_code(result, 302024)

    def test_DE_VE_073(self, preactive_bv_proposal_pipobj_list):
        pip_obj = preactive_bv_proposal_pipobj_list[0]
        result = replace_version_declare(pip_obj, pip_obj.cfg.PLATON_NEW_BIN8, pip_obj.cfg.version8)
        assert_code(result, 0)

        result = wrong_verisonsign_declare(pip_obj, preactive_bv_proposal_pipobj_list[1])
        assert_code(result, 302024)

        result = wrong_verison_declare(pip_obj, pip_obj.cfg.version0)
        assert_code(result, 302024)

class TestNoProposalCA:
    def test_DE_CA_001(self, noproposal_ca_pipobj_list, client_verifier_obj):
        pip_obj = noproposal_ca_pipobj_list[0]
        result = replace_version_declare(pip_obj, pip_obj.cfg.PLATON_NEW_BIN2, pip_obj.cfg.version2)
        assert_code(result, 0)

        result = wrong_verisonsign_declare(pip_obj, client_verifier_obj.pip)
        assert_code(result, 302024)

        result = wrong_verison_declare(pip_obj, pip_obj.chain_version)
        assert_code(result, 302024)

    def test_DE_CA_002(self, noproposal_ca_pipobj_list, client_verifier_obj):
        pip_obj = noproposal_ca_pipobj_list[0]

        result = replace_version_declare(pip_obj, pip_obj.cfg.PLATON_NEW_BIN1, pip_obj.cfg.version1)
        assert_code(result, 302028)

        result = wrong_verisonsign_declare(pip_obj, client_verifier_obj.pip)
        assert_code(result, 302024)

        result = wrong_verison_declare(pip_obj, pip_obj.chain_version)
        assert_code(result, 302024)


    def test_DE_CA_004(self, noproposal_ca_pipobj_list, client_verifier_obj):
        pip_obj = noproposal_ca_pipobj_list[0]

        result = replace_version_declare(pip_obj, pip_obj.cfg.PLATON_NEW_BIN3, pip_obj.cfg.version3)
        assert_code(result, 0)

        result = wrong_verisonsign_declare(pip_obj, client_verifier_obj.pip)
        assert_code(result, 302024)

        result = wrong_verison_declare(pip_obj, pip_obj.chain_version)
        assert_code(result, 302024)


    def test_DE_CA_005(self, noproposal_ca_pipobj_list, client_verifier_obj):
        pip_obj = noproposal_ca_pipobj_list[0]

        result = replace_version_declare(pip_obj, pip_obj.cfg.PLATON_NEW_BIN, pip_obj.cfg.version5)
        assert_code(result, 302028)

        result = wrong_verisonsign_declare(pip_obj, client_verifier_obj.pip)
        assert_code(result, 302024)

        result = wrong_verison_declare(pip_obj, pip_obj.chain_version)
        assert_code(result, 302024)

    def test_DE_CA_006(self, noproposal_ca_pipobj_list, client_verifier_obj):
        pip_obj = noproposal_ca_pipobj_list[0]

        result = replace_version_declare(pip_obj, pip_obj.cfg.PLATON_NEW_BIN8, pip_obj.cfg.version8)
        assert_code(result, 302028)

        result = wrong_verisonsign_declare(pip_obj, client_verifier_obj.pip)
        assert_code(result, 302024)

        result = wrong_verison_declare(pip_obj, pip_obj.chain_version)
        assert_code(result, 302024)

    def test_DE_CA_007(self, noproposal_ca_pipobj_list, client_verifier_obj):
        pip_obj = noproposal_ca_pipobj_list[0]

        result = replace_version_declare(pip_obj, pip_obj.cfg.PLATON_NEW_BIN0, pip_obj.cfg.version0)
        assert_code(result, 0)

        result = wrong_verisonsign_declare(pip_obj, client_verifier_obj.pip)
        assert_code(result, 302024)

        result = wrong_verison_declare(pip_obj, pip_obj.cfg.version2)
        assert_code(result, 302024)

        result = wrong_verison_declare(pip_obj, pip_obj.cfg.version3)
        assert_code(result, 302024)

    def test_DE_CA_008(self, proposal_ca_pipobj_list, client_verifier_obj):
        pip_obj = proposal_ca_pipobj_list[0]

        result = replace_version_declare(pip_obj, pip_obj.cfg.PLATON_NEW_BIN2, pip_obj.cfg.version2)
        assert_code(result, 0)

        result = wrong_verisonsign_declare(pip_obj, client_verifier_obj.pip)
        assert_code(result, 302024)

        result = wrong_verison_declare(pip_obj)
        assert_code(result, 302024)

        result = wrong_verison_declare(pip_obj, pip_obj.chain_version)
        assert_code(result, 302024)

    def test_DE_CA_010(self, proposal_ca_pipobj_list, client_verifier_obj):
        pip_obj = proposal_ca_pipobj_list[0]

        result = replace_version_declare(pip_obj, pip_obj.cfg.PLATON_NEW_BIN1, pip_obj.cfg.version1)
        assert_code(result, 302028)

        result = wrong_verisonsign_declare(pip_obj, client_verifier_obj.pip)
        assert_code(result, 302024)

        result = wrong_verison_declare(pip_obj)
        assert_code(result, 302024)

        result = wrong_verison_declare(pip_obj, pip_obj.chain_version)
        assert_code(result, 302024)

    def test_DE_CA_014(self, bv_proposal_ca_pipobj_list, client_verifier_obj):
        pip_obj = bv_proposal_ca_pipobj_list[0]

        result = replace_version_declare(pip_obj, pip_obj.cfg.PLATON_NEW_BIN2, pip_obj.cfg.version2)
        assert_code(result, 0)

        result = wrong_verisonsign_declare(pip_obj, client_verifier_obj.pip)
        assert_code(result, 302024)

        result = wrong_verison_declare(pip_obj)
        assert_code(result, 302024)

        result = wrong_verison_declare(pip_obj, pip_obj.chain_version)
        assert_code(result, 302024)

    def test_DE_CA_025(self, bv_proposal_ca_pipobj_list, client_verifier_obj):
        pip_obj = bv_proposal_ca_pipobj_list[0]

        result = replace_version_declare(pip_obj, pip_obj.cfg.PLATON_NEW_BIN1, pip_obj.cfg.version1)
        assert_code(result, 302028)

        result = wrong_verisonsign_declare(pip_obj, client_verifier_obj.pip)
        assert_code(result, 302024)

        result = wrong_verison_declare(pip_obj)
        assert_code(result, 302024)

        result = wrong_verison_declare(pip_obj, pip_obj.chain_version)
        assert_code(result, 302024)

    def test_DE_CA_032(self, proposal_ca_pipobj_list, client_verifier_obj):
        pip_obj = proposal_ca_pipobj_list[0]

        result = replace_version_declare(pip_obj, pip_obj.cfg.PLATON_NEW_BIN0, pip_obj.cfg.version0)
        assert_code(result, 0)

        result = wrong_verisonsign_declare(pip_obj, client_verifier_obj.pip)
        assert_code(result, 302024)

        result = wrong_verison_declare(pip_obj)
        assert_code(result, 302024)

    def test_DE_CA_034(self, bv_proposal_ca_pipobj_list, client_verifier_obj):
        pip_obj = bv_proposal_ca_pipobj_list[0]

        result = replace_version_declare(pip_obj, pip_obj.cfg.PLATON_NEW_BIN0, pip_obj.cfg.version0)
        assert_code(result, 0)

        result = wrong_verisonsign_declare(pip_obj, client_verifier_obj.pip)
        assert_code(result, 302024)

        result = wrong_verison_declare(pip_obj)
        assert_code(result, 302024)

    def test_DE_CA_036(self, proposal_ca_pipobj_list, client_verifier_obj):
        pip_obj = proposal_ca_pipobj_list[0]

        result = replace_version_declare(pip_obj, pip_obj.cfg.PLATON_NEW_BIN3, pip_obj.cfg.version3)
        assert_code(result, 0)

        result = wrong_verisonsign_declare(pip_obj, client_verifier_obj.pip)
        assert_code(result, 302024)

        result = wrong_verison_declare(pip_obj)
        assert_code(result, 302024)

        result = wrong_verison_declare(pip_obj, pip_obj.chain_version)
        assert_code(result, 302024)

    def test_DE_CA_038(self, proposal_ca_pipobj_list, client_verifier_obj):
        pip_obj = proposal_ca_pipobj_list[0]

        result = replace_version_declare(pip_obj, pip_obj.cfg.PLATON_NEW_BIN, pip_obj.cfg.version5)
        assert_code(result, 0)

        result = wrong_verisonsign_declare(pip_obj, client_verifier_obj.pip)
        assert_code(result, 302024)

        result = wrong_verison_declare(pip_obj, pip_obj.chain_version)
        assert_code(result, 302024)

    def test_DE_CA_040(self, proposal_ca_pipobj_list, client_verifier_obj):
        pip_obj = proposal_ca_pipobj_list[0]

        result = replace_version_declare(pip_obj, pip_obj.cfg.PLATON_NEW_BIN4, pip_obj.cfg.version4)
        assert_code(result, 0)

        result = wrong_verisonsign_declare(pip_obj, client_verifier_obj.pip)
        assert_code(result, 302024)

        result = wrong_verison_declare(pip_obj)
        assert_code(result, 302024)

        result = wrong_verison_declare(pip_obj, pip_obj.chain_version)
        assert_code(result, 302024)

    def test_DE_CA_042(self, proposal_ca_pipobj_list, client_verifier_obj):
        pip_obj = proposal_ca_pipobj_list[0]

        result = replace_version_declare(pip_obj, pip_obj.cfg.PLATON_NEW_BIN6, pip_obj.cfg.version6)
        assert_code(result, 0)

        result = wrong_verisonsign_declare(pip_obj, client_verifier_obj.pip)
        assert_code(result, 302024)

        result = wrong_verison_declare(pip_obj)
        assert_code(result, 302024)

        result = wrong_verison_declare(pip_obj, pip_obj.chain_version)
        assert_code(result, 302024)

    def test_DE_CA_044(self, proposal_ca_pipobj_list, client_verifier_obj):
        pip_obj = proposal_ca_pipobj_list[0]

        result = replace_version_declare(pip_obj, pip_obj.cfg.PLATON_NEW_BIN8, pip_obj.cfg.version8)
        assert_code(result, 302028)

        result = wrong_verisonsign_declare(pip_obj, client_verifier_obj.pip)
        assert_code(result, 302024)

        result = wrong_verison_declare(pip_obj)
        assert_code(result, 302024)

        result = wrong_verison_declare(pip_obj, pip_obj.chain_version)
        assert_code(result, 302024)

    def test_DE_CA_046(self, bv_proposal_ca_pipobj_list, client_verifier_obj):
        pip_obj = bv_proposal_ca_pipobj_list[0]

        result = replace_version_declare(pip_obj, pip_obj.cfg.PLATON_NEW_BIN3, pip_obj.cfg.version3)
        assert_code(result, 0)

        result = wrong_verisonsign_declare(pip_obj, client_verifier_obj.pip)
        assert_code(result, 302024)

        result = wrong_verison_declare(pip_obj)
        assert_code(result, 302024)

        result = wrong_verison_declare(pip_obj, pip_obj.chain_version)
        assert_code(result, 302024)

    def test_DE_CA_048(self, bv_proposal_ca_pipobj_list, client_verifier_obj):
        pip_obj = bv_proposal_ca_pipobj_list[0]

        result = replace_version_declare(pip_obj, pip_obj.cfg.PLATON_NEW_BIN, pip_obj.cfg.version5)
        assert_code(result, 302028)

        result = wrong_verisonsign_declare(pip_obj, client_verifier_obj.pip)
        assert_code(result, 302024)

        result = wrong_verison_declare(pip_obj)
        assert_code(result, 302024)

        result = wrong_verison_declare(pip_obj, pip_obj.chain_version)
        assert_code(result, 302024)

    def test_DE_CA_050(self, bv_proposal_ca_pipobj_list, client_verifier_obj):
        pip_obj = bv_proposal_ca_pipobj_list[0]

        result = replace_version_declare(pip_obj, pip_obj.cfg.PLATON_NEW_BIN4, pip_obj.cfg.version4)
        assert_code(result, 302028)

        result = wrong_verisonsign_declare(pip_obj, client_verifier_obj.pip)
        assert_code(result, 302024)

        result = wrong_verison_declare(pip_obj)
        assert_code(result, 302024)

        result = wrong_verison_declare(pip_obj, pip_obj.chain_version)
        assert_code(result, 302024)

    def test_DE_CA_052(self, bv_proposal_ca_pipobj_list, client_verifier_obj):
        pip_obj = bv_proposal_ca_pipobj_list[0]

        result = replace_version_declare(pip_obj, pip_obj.cfg.PLATON_NEW_BIN6, pip_obj.cfg.version6)
        assert_code(result, 302028)

        result = wrong_verisonsign_declare(pip_obj, client_verifier_obj.pip)
        assert_code(result, 302024)

        result = wrong_verison_declare(pip_obj)
        assert_code(result, 302024)

        result = wrong_verison_declare(pip_obj, pip_obj.chain_version)
        assert_code(result, 302024)

    def test_DE_CA_054(self, bv_proposal_ca_pipobj_list, client_verifier_obj):
        pip_obj = bv_proposal_ca_pipobj_list[0]

        result = replace_version_declare(pip_obj, pip_obj.cfg.PLATON_NEW_BIN8, pip_obj.cfg.version8)
        assert_code(result, 0)

        result = wrong_verisonsign_declare(pip_obj, client_verifier_obj.pip)
        assert_code(result, 302024)

        result = wrong_verison_declare(pip_obj, pip_obj.chain_version)
        assert_code(result, 302024)

class TestNewDeclareVersion():
    def test_DE_NN_001_to_003(self, new_genesis_env, client_con_list_obj, client_noc_list_obj):
        new_genesis_env.deploy_all()
        pip_obj = client_noc_list_obj[0].pip
        address, _ = pip_obj.economic.account.generate_account(pip_obj.node.web3, 10**18 * 10000000)
        result = pip_obj.declareVersion(pip_obj.node.node_id, address, transaction_cfg=pip_obj.cfg.transaction_cfg)
        log.info('New node declare version result : {}'.format(result))
        assert_code(result, 302023)

        submitvpandvote(client_con_list_obj)
        proposalinfo = pip_obj.get_effect_proposal_info_of_vote()
        log.info('Get version proposal information : {}'.format(proposalinfo))
        result = pip_obj.declareVersion(pip_obj.node.node_id, address, transaction_cfg=pip_obj.cfg.transaction_cfg)
        log.info('New node declare version result : {}'.format(result))
        assert_code(result, 302023)

        wait_block_number(pip_obj.node, proposalinfo.get('EndVotingBlock'))
        assert_code(pip_obj.get_status_of_proposal(proposalinfo.get('ProposalID')), 4)
        wait_block_number(pip_obj.node, proposalinfo.get('ActiveBlock'))
        assert_code(pip_obj.get_status_of_proposal(proposalinfo.get('ProposalID')), 5)

        result = pip_obj.declareVersion(pip_obj.node.node_id, address, transaction_cfg=pip_obj.cfg.transaction_cfg)
        log.info('New node declare version result : {}'.format(result))
        assert_code(result, 302023)

class TestDV():
    # def upgrade_chain(self, ):

    def test_DE_VE_003_DE_VE_009_DE_VE_012_DE_CA_003_DE_CA_009_DE_CA_012(self, new_genesis_env, client_con_list_obj):
        new_genesis_env.deploy_all()
        pip_obj = client_con_list_obj[-1].pip
        submitvpandvote(client_con_list_obj[0:3], votingrounds=3, version=pip_obj.cfg.version9)
        proposalinfo = pip_obj.get_effect_proposal_info_of_vote()
        log.info("Get version proposal information : {}".format(proposalinfo))
        wait_block_number(pip_obj.node, proposalinfo.get('EndVotingBlock'))
        assert_code(pip_obj.get_status_of_proposal(proposalinfo.get('ProposalID')), 4)
        wait_block_number(pip_obj.node, proposalinfo.get('ActiveBlock'))
        assert_code(pip_obj.get_status_of_proposal(proposalinfo.get('ProposalID')), 5)
        assert pip_obj.cfg.version9 == pip_obj.chain_version
        result = pip_obj.declareVersion(pip_obj.node.node_id, pip_obj.node.staking_address,
                                        transaction_cfg=pip_obj.cfg.transaction_cfg)
        log.info('Node {} declare version result {}'.format(pip_obj.node.node_id, result))
        assert_code(result, 302028)
        result = client_con_list_obj[0].pip.submitVersion(client_con_list_obj[0].node.node_id, str(time.time()),
                                                          pip_obj.cfg.version8, 1,
                                                          client_con_list_obj[0].node.staking_address,
                                                          transaction_cfg=pip_obj.cfg.transaction_cfg)
        log.info('Node {} submit version proposal result : {}'.format(client_con_list_obj[0].node.node_id, result))
        assert_code(result, 0)
        result = replace_version_declare(pip_obj, pip_obj.cfg.PLATON_NEW_BIN4, versiontag=pip_obj.cfg.version4)
        assert_code(result, 302028)
        verifier_list = get_pledge_list(client_con_list_obj[-1].ppos.getVerifierList)
        log.info('verifier list : {}'.format(verifier_list))
        assert pip_obj.node.node_id not in verifier_list
        client_con_list_obj[-1].economic.wait_settlement_blocknum(client_con_list_obj[-1].node)
        assert pip_obj.node.node_id not in verifier_list

        result = replace_version_declare(pip_obj, pip_obj.cfg.PLATON_NEW_BIN0, versiontag=pip_obj.cfg.version0)
        assert_code(result, 302028)
        result = client_con_list_obj[0].pip.submitVersion(client_con_list_obj[0].node.node_id, str(time.time()),
                                                          pip_obj.cfg.version8, 1,
                                                          client_con_list_obj[0].node.staking_address,
                                                          transaction_cfg=pip_obj.cfg.transaction_cfg)
        log.info('Node {} submit version proposal result : {}'.format(client_con_list_obj[0].node.node_id, result))
        assert_code(result, 0)
        result = replace_version_declare(pip_obj, pip_obj.cfg.PLATON_NEW_BIN0, versiontag=pip_obj.cfg.version0)
        assert_code(result, 302028)
        result = replace_version_declare(pip_obj, pip_obj.cfg.PLATON_NEW_BIN4, versiontag=pip_obj.cfg.version4)
        assert_code(result, 302028)

