import os

# PlatON-Tests path
BASE_DIR = os.path.dirname(os.path.dirname(os.path.abspath(__file__)))

# server deployment node path
DEPLOY_PATH = r"trantor_test"

# script run log level
RUN_LOG_LEVEL = "info"

# local must file
PLATON_BIN_FILE = os.path.abspath(os.path.join(BASE_DIR, "deploy/bin/platon"))
GENESIS_FILE = os.path.abspath(os.path.join(BASE_DIR, "deploy/template/genesis_template.json"))
CONFIG_JSON_FILE = os.path.abspath(os.path.join(BASE_DIR, "deploy/template/config_template.json"))
STATIC_NODE_FILE = os.path.abspath(os.path.join(BASE_DIR, 'deploy/template/static-nodes.json'))
SUPERVISOR_FILE = os.path.abspath(os.path.join(BASE_DIR, "deploy/template/supervisor_template.conf"))
ADDRESS_FILE = os.path.abspath(
    os.path.join(BASE_DIR, 'deploy/template/UTC--2019-08-23T12-33-18.192329788Z--2e95e3ce0a54951eb9a99152a6d5827872dfb4fd'))
ACCOUNT_FILE = os.path.abspath(os.path.join(BASE_DIR, "deploy/accounts.yml"))
LOG_FILE = os.path.abspath(os.path.join(BASE_DIR, "deploy/tmp/platon.log"))

# cache file root directory
LOCAL_TMP_FILE_ROOT_DIR = os.path.abspath(os.path.join(BASE_DIR, "deploy/tmp"))
# LOCAL_TMP_FILE_FOR_NODE = os.path.abspath(os.path.join(LOCAL_TMP_FILE_ROOT_DIR, "node"))
# LOCAL_TMP_FILE_FOR_SERVER = os.path.abspath(os.path.join(LOCAL_TMP_FILE_ROOT_DIR, "server"))
# LOCAL_TMP_FILE_FOR_ENV = os.path.abspath(os.path.join(LOCAL_TMP_FILE_ROOT_DIR, "env"))
# LOCAL_TMP_FILE_FOR_GNENESIS = os.path.abspath(os.path.join(LOCAL_TMP_FILE_ROOT_DIR, "genesis.json"))


# cache directory configuration
class ConfTmpDir:
    def __init__(self, dir):
        self.dir = dir
        self.tmp_root_path = os.path.abspath(os.path.join(LOCAL_TMP_FILE_ROOT_DIR, dir))
        if not os.path.exists(self.tmp_root_path):
            os.makedirs(self.tmp_root_path)
        self.GENESIS_FILE = os.path.abspath(os.path.join(self.tmp_root_path, "genesis.json"))
        self.CONFIG_JSON_FILE = os.path.abspath(os.path.join(self.tmp_root_path, "config.json"))
        self.STATIC_NODE_FILE = os.path.abspath(os.path.join(self.tmp_root_path, 'static-nodes.json'))
        self.LOCAL_TMP_FILE_FOR_NODE = os.path.abspath(os.path.join(self.tmp_root_path, 'node'))
        self.LOCAL_TMP_FILE_FOR_SERVER = os.path.abspath(os.path.join(self.tmp_root_path, 'server'))
        self.LOCAL_TMP_FILE_FOR_ENV = os.path.abspath(os.path.join(self.tmp_root_path, 'env'))


# directory cache configuration
DEFAULT_CONF_TMP_DIR = ConfTmpDir("global")

NODE_FILE = os.path.abspath(os.path.join(BASE_DIR, "deploy/node/ppos_4_4_wyq1.yml"))
