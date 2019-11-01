import os
from conf.settings import BASE_DIR


class StakingConfig:
    def __init__(self, external_id, node_name, website, details):
        self.external_id = external_id
        self.node_name = node_name
        self.website = website
        self.details = details


class PipConfig:
    PLATON_NEW_BIN = os.path.abspath(os.path.join(BASE_DIR, "deploy/bin/newpackage/platon"))
    PLATON_NEW_BIN1 = os.path.abspath(os.path.join(BASE_DIR, "deploy/bin/newpackage/version1/platon"))
    PLATON_NEW_BIN2 = os.path.abspath(os.path.join(BASE_DIR, "deploy/bin/newpackage/version2/platon"))
    PLATON_NEW_BIN3 = os.path.abspath(os.path.join(BASE_DIR, "deploy/bin/newpackage/version3/platon"))
    PLATON_NEW_BIN4 = os.path.abspath(os.path.join(BASE_DIR, "deploy/bin/newpackage/version4/platon"))
    PLATON_NEW_BIN8 = os.path.abspath(os.path.join(BASE_DIR, "deploy/bin/newpackage/version8/platon"))
    PLATON_NEW_BIN9 = os.path.abspath(os.path.join(BASE_DIR, "deploy/bin/newpackage/version9/platon"))
    PLATON_NEW_BIN6 = os.path.abspath(os.path.join(BASE_DIR, "deploy/bin/newpackage/version6/platon"))
    PLATON_NEW_BIN7 = os.path.abspath(os.path.join(BASE_DIR, "deploy/bin/newpackage/version7/platon"))
    PLATON_NEW_BIN0 = os.path.abspath(os.path.join(BASE_DIR, "deploy/bin/platon"))
    PLATON_NEW_BIN120 = os.path.abspath(os.path.join(BASE_DIR, "deploy/bin/newpackage/diffcodeversion/platon"))
    text_proposal = 1
    cancel_proposal = 4
    version_proposal = 2
    param_proposal = 3
    vote_option_yeas = 1
    vote_option_nays = 2
    vote_option_Abstentions = 2
    version1 = 1537
    version2 = 1792
    version3 = 1801
    version4 = 2048
    version5 = 2049
    version6 = 2057
    version7 = 2304
    verison8 = 591617
    version9 = 526081
    transaction_cfg = {"gasPrice": 3000000000000000, "gas": 1000000}

class EconomicConfig:
    RELEASE_ZERO = 62215742.48691650
    TOKEN_TOTAL = 10250000000000000000000000000
    # PlatON Foundation Address
    FOUNDATION_ADDRESS = "0x493301712671ada506ba6ca7891f436d29185821"
    # Lock account account address
    FOUNDATION_LOCKUP_ADDRESS = "0x1000000000000000000000000000000000000001"
    # Pledged contract address
    STAKING_ADDRESS = "0x1000000000000000000000000000000000000002"
    # PlatON incentive pool account
    INCENTIVEPOOL_ADDRESS = "0x1000000000000000000000000000000000000003"
    # Remaining total account
    REMAIN_ACCOUNT_ADDRESS = "0x2e95e3ce0a54951eb9a99152a6d5827872dfb4fd"
    # Developer Foundation Account
    DEVELOPER_FOUNDATAION_ADDRESS = '0x54a7a3c6822eb222c53F76443772a60b0f9A8bab'
    init_token_info = {FOUNDATION_LOCKUP_ADDRESS: 259096240418673500000000000,
                       STAKING_ADDRESS: 40000000000000000000000000,
                       INCENTIVEPOOL_ADDRESS: 262215742486916500000000000,
                       FOUNDATION_ADDRESS: 1638688017094410000000000000,
                       REMAIN_ACCOUNT_ADDRESS: 8000000000000000000000000000,
                       DEVELOPER_FOUNDATAION_ADDRESS: 50000000000000000000000000
                       }
    release_info = [{"blockNumber": 1600, "amount": 55965742486916500000000000},
                    {"blockNumber": 3200, "amount": 49559492486916500000000000},
                    {"blockNumber": 4800, "amount": 42993086236916500000000000},
                    {"blockNumber": 6400, "amount": 36262519830666600000000000},
                    {"blockNumber": 8000, "amount": 29363689264263300000000000},
                    {"blockNumber": 9600, "amount": 22292387933693900000000000},
                    {"blockNumber": 11200, "amount": 15044304069863300000000000},
                    {"blockNumber": 12800, "amount": 7615018109436900000000000}
                    ]
