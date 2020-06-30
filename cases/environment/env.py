import json
import os
import time
import random
import shutil
import tarfile
from concurrent.futures import ThreadPoolExecutor, wait, ALL_COMPLETED
import copy
from common.load_file import get_f
from ruamel import yaml
from environment.node import Node
from environment.server import Server
from common.abspath import abspath
from common.key import generate_key, generate_blskey
from common.load_file import LoadFile, calc_hash
from common.log import log
from environment.account import Account
from environment.config import TestConfig
from conf.settings import DEFAULT_CONF_TMP_DIR, ConfTmpDir
from typing import List


def check_file_exists(*args):
    """
    Check if local files exist
    :param args:
    """
    for arg in args:
        if not os.path.exists(os.path.abspath(arg)):
            raise Exception("file:{} does not exist".format(arg))


class TestEnvironment:
    def __init__(self, cfg: TestConfig):
        # env config
        self.cfg = cfg

        # these file must be exist
        check_file_exists(self.cfg.platon_bin_file, self.cfg.genesis_file, self.cfg.supervisor_file,
                          self.cfg.node_file, self.cfg.address_file)
        if not os.path.exists(self.cfg.root_tmp):
            os.mkdir(self.cfg.root_tmp)

        # node config
        self.__is_update_node_file = False
        self.node_config = LoadFile(self.cfg.node_file).get_data()
        self.consensus_node_config_list = self.node_config.get("consensus", [])
        self.noconsensus_node_config_list = self.node_config.get("noconsensus", [])
        self.node_config_list = self.consensus_node_config_list + self.noconsensus_node_config_list
        self.__rewrite_node_file()

        # node obj list
        self.__consensus_node_list = []
        self.__normal_node_list = []

        # env info
        self.cfg.env_id = self.__reset_env()

        # genesis
        self.genesis_config = LoadFile(self.cfg.genesis_file).get_data()

        # servers
        self.server_list = self.__parse_servers()

        # node
        self.__parse_node()

        # accounts
        self.account = Account(self.cfg.account_file, self.genesis_config["config"]["chainId"])

        self.rewrite_genesis_file()

    @property
    def consensus_node_list(self) -> List[Node]:
        return self.__consensus_node_list

    @property
    def normal_node_list(self) -> List[Node]:
        return self.__normal_node_list

    @property
    def chain_id(self):
        return self.genesis_config["config"]["chainId"]

    @property
    def amount(self):
        return self.genesis_config["config"]["cbft"]["amount"]

    @property
    def period(self):
        return self.genesis_config["config"]["cbft"]["period"]

    @property
    def validatorMode(self):
        return self.genesis_config["config"]["cbft"]["validatorMode"]

    @property
    def version(self):
        return "0.13.0"

    @property
    def running(self) -> bool:
        """
        Determine if all nodes are running
        :return: bool
        """
        for node in self.get_all_nodes():
            if not node.running:
                return False
        return True

    @property
    def max_byzantium(self) -> int:
        """
        Maximum number of Byzantine nodes
        """
        return get_f(self.consensus_node_config_list)

    @property
    def block_interval(self) -> int:
        """
        Block interval
        """
        period = self.genesis_config["config"]["cbft"].get("period")
        amount = self.genesis_config["config"]["cbft"].get("amount")
        return int(period / 1000 / amount)

    def consensus_node_id_list(self) -> List[str]:
        return [node.node_id for node in self.consensus_node_list]

    def find_node_by_node_id(self, node_id):
        for node in self.get_all_nodes():
            if node_id == node.node_id:
                return node
        raise Exception("can't find node")

    def copy_env(self):
        """
        Copy environment
        """
        return copy.copy(self)

    def set_cfg(self, cfg: TestConfig):
        """
        Set the configuration file and modify the node's cfg
        :param cfg:
        """
        self.cfg = cfg
        genesis_config = LoadFile(self.cfg.genesis_file).get_data()
        self.rewrite_genesis_file()
        self.set_genesis(genesis_config)
        for node in self.get_all_nodes():
            node.cfg = cfg

    def set_genesis(self, genesis_config: dict):
        """
        Set the genesis and modify the genesis of the node.
        :param genesis_config:
        """
        self.genesis_config = genesis_config
        self.account.chain_id = self.chain_id
        for node in self.get_all_nodes():
            node.chain_id = self.chain_id

    def __reset_env(self) -> str:
        """
        Determine whether you need to re-create a new environment
        based on the platon binary information and the node configuration file.
        :return: env_id
        """
        env_tmp_file = os.path.join(self.cfg.env_tmp, "env.yml")
        if os.path.exists(self.cfg.env_tmp):
            if os.path.exists(env_tmp_file):
                env_data = LoadFile(env_tmp_file).get_data()
                if env_data["bin_hash"] == calc_hash(self.cfg.platon_bin_file) \
                        and env_data["node_hash"] == calc_hash(self.cfg.node_file):
                    return env_data["env_id"]

            shutil.rmtree(self.cfg.env_tmp)
        os.makedirs(self.cfg.env_tmp)
        new_env_data = {"bin_hash": calc_hash(self.cfg.platon_bin_file), "node_hash": calc_hash(self.cfg.node_file)}
        env_id = new_env_data["bin_hash"] + new_env_data["node_hash"]
        new_env_data["env_id"] = env_id
        with open(env_tmp_file, "w", encoding="utf-8") as f:
            yaml.dump(new_env_data, f, Dumper=yaml.RoundTripDumper)
        return env_id

    def get_init_nodes(self) -> List[dict]:
        """
        Get the list of init nodes
        :return: list
        """
        init_node_list = []
        for node in self.consensus_node_list:
            init_node_list.append({"node": node.enode, "blsPubKey": node.blspubkey})
        return init_node_list

    def get_static_nodes(self) -> list:
        """
        Get static node enode list
        :return: list
        """
        static_node_list = []
        for node in self.get_all_nodes():
            static_node_list.append(node.enode)
        return static_node_list

    def get_all_nodes(self) -> List[Node]:
        """
        Get all node objects
        :return: Node object
        """
        return self.__consensus_node_list + self.__normal_node_list

    def get_rand_node(self) -> Node:
        """
        Randomly obtain a consensus node
        :return: Node object
        """
        return random.choice(self.consensus_node_list)

    def get_consensus_node_by_index(self, index) -> Node:
        """
        Get a consensus node based on the index
        :param index:
        :return: Node object
        """
        return self.__consensus_node_list[index]

    def get_normal_node_by_index(self, index) -> Node:
        """
        Get a normal node based on the index
        :param index:
        :return: Node object
        """
        return self.__normal_node_list[index]

    def get_a_normal_node(self) -> Node:
        """
        Get the first normal node
        :return: Node object
        """
        return self.__normal_node_list[0]

    def executor(self, func, data_list, *args) -> bool:
        with ThreadPoolExecutor(max_workers=self.cfg.max_worker) as exe:
            futures = [exe.submit(func, pair, *args) for pair in data_list]
            done, unfinished = wait(futures, timeout=30, return_when=ALL_COMPLETED)
        result = []
        for d in done:
            is_success, msg = d.result()
            if not is_success:
                result.append(msg)
        if len(result) > 0:
            raise Exception("executor {} failed:{}".format(func.__name__, result))
        return True

    def deploy_all(self, genesis_file=None):
        """
        Deploy all nodes and start
        :param genesis_file: Specify genesis, do not pass the default generated using tmp
        """
        self.account.reset()
        self.prepare_all()
        if genesis_file is None:
            genesis_file = self.cfg.genesis_tmp
        log.info("deploy all node")
        self.deploy_nodes(self.get_all_nodes(), genesis_file)
        log.info("deploy success")

    def prepare_all(self):
        """
        Prepare environmental data
        """
        self.rewrite_genesis_file()
        self.rewrite_static_nodes()
        self.rewrite_config_json()
        self.__compression()
        if self.cfg.install_supervisor:
            self.install_all_supervisor()
            self.cfg.install_supervisor = False
        if self.cfg.install_dependency:
            self.install_all_dependency()
            self.cfg.install_dependency = False
        self.put_all_compression()

    def start_all(self):
        """
        Start all nodes, judge whether to initialize according to the value of cfg init_chain
        """
        log.info("start all node")
        self.start_nodes(self.get_all_nodes(), self.cfg.init_chain)

    def stop_all(self):
        """
        Stop all nodes
        """
        log.info("stop all node")
        self.stop_nodes(self.get_all_nodes())

    def reset_all(self):
        """
        Restart all nodes
        """
        log.info("restart all node")
        self.reset_nodes(self.get_all_nodes())

    def clean_all(self):
        """
        Close all nodes and delete the directory of the deployment node
        """
        log.info("clean all node")
        self.clean_nodes(self.get_all_nodes())

    def clean_db_all(self):
        """
        Close all nodes and delete the database
        """
        log.info("clean db all node")
        self.clean_db_nodes(self.get_all_nodes())

    def shutdown(self):
        """
        Close all nodes and delete the node deployment directory, supervisor node configuration
        """
        log.info("shutdown and clean all nodes")

        def close(node: Node):
            return node.close()

        return self.executor(close, self.get_all_nodes())

    def clean_supervisor_confs(self):
        def clean(server: Server):
            return server.clean_supervisor_conf()
        return self.executor(clean, self.server_list)

    def start_nodes(self, node_list: List[Node], init_chain=True):
        """
        Boot node
        :param node_list:
        :param init_chain:
        """
        def start(node: Node, need_init_chain):
            return node.start(need_init_chain)

        return self.executor(start, node_list, init_chain)

    def deploy_nodes(self, node_list: List[Node], genesis_file):
        """
        Deployment node
        Choose whether to empty the environment depending on whether initialization is required
        Upload all node files
        :param node_list:
        :param genesis_file:
        """
        log.info("deploy node")
        if self.cfg.init_chain:
            self.clean_nodes(node_list)

        self.put_file_nodes(node_list, genesis_file)
        return self.start_nodes(node_list, self.cfg.init_chain)

    def put_file_nodes(self, node_list: List[Node], genesis_file):
        """
        Upload all files
        :param node_list:
        :param genesis_file:
        """
        def prepare(node: Node):
            return node.put_all_file(genesis_file)

        return self.executor(prepare, node_list)

    def stop_nodes(self, node_list: List[Node]):
        """
        Close node
        :param node_list:
        """
        def stop(node: Node):
            return node.stop()

        return self.executor(stop, node_list)

    def reset_nodes(self, node_list: List[Node]):
        """
        Restart node
        :param node_list:
        """
        def restart(node: Node):
            return node.restart()

        return self.executor(restart, node_list)

    def clean_nodes(self, node_list: List[Node]):
        """
        Close the node and delete the node data
        :param node_list:
        :return:
        """
        def clean(node: Node):
            return node.clean()

        return self.executor(clean, node_list)

    def clean_db_nodes(self, node_list: List[Node]):
        """
        Close the node and clear the node database
        :param node_list:
        """
        def clean_db(node: Node):
            return node.clean_db()

        return self.executor(clean_db, node_list)

    def __parse_node(self):
        """
        Instantiate all nodes
        """
        def init(node_config):
            return Node(node_config, self.cfg, self.chain_id)

        log.info("parse node to node object")
        with ThreadPoolExecutor(max_workers=self.cfg.max_worker) as executor:
            futures = [executor.submit(init, pair) for pair in self.consensus_node_config_list]
            done, unfinished = wait(futures, timeout=30, return_when=ALL_COMPLETED)
        for do in done:
            self.__consensus_node_list.append(do.result())

        if self.noconsensus_node_config_list:
            with ThreadPoolExecutor(max_workers=self.cfg.max_worker) as executor:
                futures = [executor.submit(init, pair) for pair in self.noconsensus_node_config_list]
                done, unfinished = wait(futures, timeout=30, return_when=ALL_COMPLETED)
            for do in done:
                self.__normal_node_list.append(do.result())

    def put_all_compression(self):
        """
        Upload compressed file
        """
        log.info("upload compression")

        def uploads(server: Server):
            return server.put_compression()

        return self.executor(uploads, self.server_list)

    def install_all_dependency(self):
        """
        Installation dependence
        """
        log.info("install rely")

        def install(server: Server):
            return server.install_dependency()

        return self.executor(install, self.server_list)

    def install_all_supervisor(self):
        """
        Install supervisor
        """
        log.info("install supervisor")

        def install(server: Server):
            return server.install_supervisor()

        return self.executor(install, self.server_list)

    def __parse_servers(self) -> List[Server]:
        """
        Instantiate all servers
        """
        server_config_list, server_list = [], []

        def check_in(_ip, nodes):
            for n in nodes:
                if _ip == n["host"]:
                    return True
            return False

        for node_config in self.node_config_list:
            ip = node_config["host"]
            if check_in(ip, server_config_list):
                continue
            server_config_list.append(node_config)

        def init(config):
            return Server(config, self.cfg)

        with ThreadPoolExecutor(max_workers=self.cfg.max_worker) as executor:
            futures = [executor.submit(init, pair) for pair in server_config_list]
            done, unfinished = wait(futures, timeout=30, return_when=ALL_COMPLETED)
        for do in done:
            server_list.append(do.result())
        return server_list

    def block_numbers(self, node_list: List[Node] = None) -> dict:
        """
        Get the block height of the incoming node
        :param node_list:
        """
        if node_list is None:
            node_list = self.get_all_nodes()
        result = {}
        for node in node_list:
            result[node.node_mark] = node.block_number
        return result

    def check_block(self, need_number=10, multiple=3, node_list: List[Node] = None):
        """
        Verify the highest block in the current chain
        :param need_number:
        :param multiple:
        :param node_list:
        """
        if node_list is None:
            node_list = self.get_all_nodes()
        use_time = int(need_number * self.block_interval * multiple)
        while use_time:
            if max(self.block_numbers(node_list).values()) < need_number:
                time.sleep(1)
                use_time -= 1
                continue
            return
        raise Exception("The environment is not working properly")

    def backup_all_logs(self, case_name: str):
        """
        Download all node logs
        """
        return self.backup_logs(self.get_all_nodes(), case_name)

    def backup_logs(self, node_list: List[Node], case_name):
        """
        Backup log
        :param node_list:
        :param case_name:
        """
        self.__check_log_path()

        def backup(node: Node):
            return node.backup_log()

        self.executor(backup, node_list)
        return self.__zip_all_log(case_name)

    def __check_log_path(self):
        if not os.path.exists(self.cfg.tmp_log):
            os.mkdir(self.cfg.tmp_log)
        else:
            shutil.rmtree(self.cfg.tmp_log)
            os.mkdir(self.cfg.tmp_log)
        if not os.path.exists(self.cfg.bug_log):
            os.mkdir(self.cfg.bug_log)

    def __zip_all_log(self, case_name):
        log.info("Start compressing.....")
        t = time.strftime("%Y%m%d%H%M%S", time.localtime())
        tar_name = "{}/{}_{}.tar.gz".format(self.cfg.bug_log, case_name, t)
        tar = tarfile.open(tar_name, "w:gz")
        tar.add(self.cfg.tmp_log, arcname=os.path.basename(self.cfg.tmp_log))
        tar.close()
        log.info("Compression completed")
        log.info("Start deleting the cache.....")
        shutil.rmtree(self.cfg.tmp_log)
        log.info("Delete cache complete")
        return os.path.basename(tar_name)

    def rewrite_genesis_file(self):
        """
        Rewrite genesis
        """
        log.info("rewrite genesis.json")
        self.genesis_config['config']['cbft']["initialNodes"] = self.get_init_nodes()
        # with open(self.cfg.address_file, "r", encoding="UTF-8") as f:
        #     key_dict = json.load(f)
        # account = key_dict["address"]
        # self.genesis_config['alloc'][account] = {"balance": str(99999999999999999999999999)}
        accounts = self.account.get_all_accounts()
        for account in accounts:
            self.genesis_config['alloc'][account['address']] = {"balance": str(account['balance'])}
        with open(self.cfg.genesis_tmp, 'w', encoding='utf-8') as f:
            f.write(json.dumps(self.genesis_config, indent=4))

    def rewrite_static_nodes(self):
        """
        Rewrite static
        """
        log.info("rewrite static-nodes.json")
        static_nodes = self.get_static_nodes()
        with open(self.cfg.static_node_tmp, 'w', encoding='utf-8') as f:
            f.write(json.dumps(static_nodes, indent=4))

    def rewrite_config_json(self):
        """
        Rewrite config
        :return:
        """
        log.info("rewrite config.json")
        config_data = LoadFile(self.cfg.config_json_file).get_data()
        # config_data['node']['P2P']["BootstrapNodes"] = self.get_static_nodes()
        with open(self.cfg.config_json_tmp, 'w', encoding='utf-8') as f:
            f.write(json.dumps(config_data, indent=4))

    def __fill_node_config(self, node_config: dict):
        """
        Fill in the node file with some necessary values
        :param node_config:
        """
        if not node_config.get("id") or not node_config.get("nodekey"):
            self.__is_update_node_file = True
            node_config["nodekey"], node_config["id"] = generate_key()
        if not node_config.get("blsprikey") or not node_config.get("blspubkey"):
            self.__is_update_node_file = True
            node_config["blsprikey"], node_config["blspubkey"] = generate_blskey()
        if not node_config.get("port"):
            self.__is_update_node_file = True
            node_config["port"] = 16789
        if not node_config.get("rpcport"):
            self.__is_update_node_file = True
            node_config["rpcport"] = 6789
        if not node_config.get("url"):
            self.__is_update_node_file = True
            node_config["url"] = "http://{}:{}".format(node_config["host"], node_config["rpcport"])
        if node_config.get("wsport"):
            self.__is_update_node_file = True
            node_config["wsurl"] = "ws://{}:{}".format(node_config["host"], node_config["wsport"])
        return node_config

    def __rewrite_node_file(self):
        log.info("rewrite node file")
        result, result_consensus_list, result_noconsensus_list = {}, [], []
        if len(self.consensus_node_config_list) >= 1:
            for node_config in self.consensus_node_config_list:
                result_consensus_list.append(self.__fill_node_config(node_config))
            result["consensus"] = result_consensus_list
        if self.noconsensus_node_config_list and len(self.noconsensus_node_config_list) >= 1:
            for node_config in self.noconsensus_node_config_list:
                result_noconsensus_list.append(self.__fill_node_config(node_config))
            result["noconsensus"] = result_noconsensus_list
        if self.__is_update_node_file:
            self.consensus_node_config_list = result_consensus_list
            self.noconsensus_node_config_list = result_noconsensus_list
            with open(self.cfg.node_file, encoding="utf-8", mode="w") as f:
                yaml.dump(result, f, Dumper=yaml.RoundTripDumper)

    def __compression(self):
        """
        Compressed file
        """
        log.info("compression data")
        env_gz = os.path.join(self.cfg.env_tmp, self.cfg.env_id)
        if os.path.exists(env_gz):
            return
        os.makedirs(env_gz)
        data_dir = os.path.join(env_gz, "data")
        os.makedirs(data_dir)
        keystore_dir = os.path.join(data_dir, "keystore")
        os.makedirs(keystore_dir)
        keystore = os.path.join(keystore_dir, os.path.basename(self.cfg.address_file))
        shutil.copyfile(self.cfg.address_file, keystore)
        shutil.copyfile(self.cfg.platon_bin_file, os.path.join(env_gz, "platon"))
        shutil.copyfile(self.cfg.config_json_tmp, os.path.join(env_gz, "config.json"))
        t = tarfile.open(env_gz + ".tar.gz", "w:gz")
        t.add(env_gz, arcname=os.path.basename(env_gz))
        t.close()


def create_env(conf_tmp=None, node_file=None, account_file=None, init_chain=True,
               install_dependency=False, install_supervisor=False, can_deploy=True) -> TestEnvironment:
    if not conf_tmp:
        conf_tmp = DEFAULT_CONF_TMP_DIR
    else:
        conf_tmp = ConfTmpDir(conf_tmp)
    cfg = TestConfig(conf_tmp=conf_tmp, install_supervisor=install_supervisor, install_dependency=install_dependency, init_chain=init_chain, can_deploy=can_deploy)
    if node_file:
        cfg.node_file = node_file
    if account_file:
        cfg.account_file = account_file
    return TestEnvironment(cfg)


if __name__ == "__main__":
    from tests.lib import get_no_pledge_node, get_no_pledge_node_list, get_pledge_list, check_node_in_list
    node_filename = abspath("deploy/node/debug_4_4.yml")
    env = create_env(node_file=node_filename)
    env.shutdown()
    exit(0)
    # print(os.path.getctime(env.cfg.platon_bin_file))
    # new_cfg = copy.copy(env.cfg)
    # new_cfg.syncmode = "fast"
    # print(env.cfg.syncmode)
    log.info("测试部署")
    env.cfg.syncmode = "fast"
    # env.deploy_all(abspath("deploy/tmp/genesis_0.8.0.json"))
    for node in env.get_all_nodes():
        node.admin.addPeer("enode://d203e37d86f1757ee4bbeafd7a0b0b6f7d4f22afaad3c63337e92d9056251dcca95515cb23d5a040e3416e1ebbb303aa52ea535103b9b8c8c2adc98ea3b41c01@10.10.8.195:16789")
        print(node.web3.net.peerCount)
    print(node.ppos.getCandidateList())
    # env.deploy_all(abspath("deploy/tmp/genesis_0.8.0.json"))
    # env.shutdown()
    # stop_nodes = env.consensus_node_list[:2]
    # for node in stop_nodes:
    #     print(node.url)
    # time.sleep(50)
    # print(env.consensus_node_list[3].url)
    # env.stop_nodes(stop_nodes)
    # node = env.get_consensus_node_by_index(0)
    # print(node.debug.economicConfig())
    # print(type(node.debug.economicConfig()))
    # print(node.node_mark)
    # address, prikey = env.account.generate_account(node.web3, 10**18*100000000000)
    # transaction_cfg = {"gasPrice": 3000000000000000, "gas": 1000000}
    # # print(node.pip.submitParam(node.node_id, "ddd", "Slashing", "SlashBlockReward", "1000", prikey, transaction_cfg))
    # print(node.pip.getGovernParamValue("Slashing", "SlashBlockReward", address))
    # print(node.pip.listGovernParam("Staking"))
    # from tests.lib.genesis import Genesis
    # from dacite import from_dict
    # genesis = from_dict(data_class=Genesis, data=env.genesis_config)
    # print(genesis.EconomicModel.Slashing.MaxEvidenceAge)
    # env.account.generate_account(env.get_a_normal_node().web3, 0)
    # log.info("account:{}".format(env.account.accounts))
    # env.deploy_all()
    # log.info("account:{}".format(env.account.accounts))
    # node = env.get_rand_node()
    # print(node.node_id)
    # print(env.normal_node_list[1].node_id)
    # print(env.get_normal_node_by_index(1).node_id)
    # print(get_no_pledge_node(env.get_all_nodes()))
    # print(get_no_pledge_node_list(env.get_all_nodes()))
    # print(get_pledge_list(node.ppos.getVerifierList))
    # print(check_node_in_list(node.node_id, node.ppos.getVerifierList))
    # print(env.block_numbers(env.normal_node_list))
    # print(env.block_numbers())
    # for node in env.consensus_node_list:
    #     print(node.node_id)
    # time.sleep(3000)
    # env.deploy_all()
    # d = env.block_numbers()
    # print(d)
    # node = env.get_rand_node()
    # node.create_keystore()
    # print(node.node_mark)
    # time.sleep(80)
    # log.info("测试关闭")
    # env.stop_all()
    # time.sleep(30)
    # log.info("测试不初始化启动")
    # env.cfg.init_chain = False
    # env.start_all()
    # time.sleep(60)
    # d = env.block_numbers()
    # print(d)
    # log.info("测试重启")
    # env.reset_all()
    # time.sleep(60)
    # d = env.block_numbers()
    # print(d)
    # log.info("测试删除数据库")
    # env.clean_db_all()
    # log.info("删除数据库成功")
    # time.sleep(60)
    # env.cfg.init_chain = True
    # env.start_all()
    # time.sleep(30)
    # d = env.block_numbers()
    # print(d)
    # log.info("测试删除所有数据")
    # env.clean_all()
    # log.info("删除数据成功")
    # log.info("重新部署")
    # env.deploy_all()
    # d = env.block_numbers()
    # print(d)
    # time.sleep(60)
    # d = env.block_numbers()
    # print(d)
    # env.shutdown()
