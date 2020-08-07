import os
import json
from client_sdk_python import Web3
from client_sdk_python.admin import Admin
from client_sdk_python.debug import Debug
from client_sdk_python.eth import Eth
from client_sdk_python.personal import Personal
from client_sdk_python.ppos import Ppos
from client_sdk_python.pip import Pip

from common.connect import run_ssh, connect_linux, wait_connect_web3
from common.load_file import LoadFile
from environment.config import TestConfig
from environment.mock import mock_connect_linux
from common.log import log


failed_msg = "Node-{} do {} failed:{}"
success_msg = "Node-{} do {} success"


class Node:
    def __init__(self, node_conf, cfg: TestConfig, chain_id):
        self.cfg = cfg
        # node identity parameter
        self.blspubkey = node_conf["blspubkey"]
        self.blsprikey = node_conf["blsprikey"]
        self.node_id = node_conf["id"]
        self.nodekey = node_conf["nodekey"]
        # node startup necessary parameters
        self.p2p_port = str(node_conf["port"])
        self.rpc_port = str(node_conf["rpcport"])
        # node starts non essential parameters
        self.wsport = node_conf.get("wsport")
        self.wsurl = node_conf.get("wsurl")
        self.pprofport = node_conf.get("pprofport")
        self.fail_point = node_conf.get("fail_point")
        # node server information
        self.host = node_conf["host"]
        self.username = node_conf["username"]
        self.password = node_conf["password"]
        self.ssh_port = node_conf.get("sshport", 22)
        if self.cfg.can_deploy:
            self.ssh, self.sftp, self.t = connect_linux(self.host, self.username, self.password, self.ssh_port)
        else:
            self.ssh, self.sftp, self.t = mock_connect_linux()
        # node identification information
        self.url = node_conf["url"]
        self.node_name = "node-" + self.p2p_port
        self.node_mark = self.host + ":" + self.p2p_port
        # node remote directory information
        if os.path.isabs(self.cfg.deploy_path):
            self.remote_node_path = "{}/{}".format(self.cfg.deploy_path, self.node_name)
        else:
            self.remote_node_path = "{}/{}/{}".format(self.pwd, self.cfg.deploy_path, self.node_name)

        self.remote_log_dir = '{}/log'.format(self.remote_node_path)
        self.remote_bin_file = self.remote_node_path + "/platon"
        self.remote_genesis_file = self.remote_node_path + "/genesis.json"
        self.remote_config_file = self.remote_node_path + "/config.json"
        self.remote_data_dir = self.remote_node_path + "/data"

        self.remote_blskey_file = '{}/blskey'.format(self.remote_data_dir)
        self.remote_nodekey_file = '{}/nodekey'.format(self.remote_data_dir)
        self.remote_keystore_dir = '{}/keystore'.format(self.remote_data_dir)
        self.remote_static_nodes_file = '{}/static-nodes.json'.format(self.remote_data_dir)
        self.remote_db_dir = '{}/platon'.format(self.remote_data_dir)

        self.remote_supervisor_node_file = '{}/{}.conf'.format(self.cfg.remote_supervisor_tmp, self.node_name)

        # RPC connection
        self.__is_connected = False
        self.__rpc = None

        self.__is_connected_ppos = False
        self.__ppos = None

        self.__is_connected_pip = False
        self.__pip = None

        self.__is_ws_connected = False
        self.__ws_rpc = None

        # remote directory
        self.make_remote_dir()

        # node local tmp
        self.local_node_tmp = self.gen_node_tmp()

        # self.genesis_config = LoadFile(self.cfg.genesis_file).get_data()
        self.chain_id = chain_id

    @property
    def pwd(self):
        pwd_list = self.run_ssh("pwd")
        pwd = pwd_list[0].strip("\r\n")
        return pwd

    def make_remote_dir(self):
        self.run_ssh("mkdir -p {}".format(self.remote_node_path))
        self.run_ssh('mkdir -p {}/log'.format(self.remote_node_path))
        self.run_ssh("mkdir -p {}".format(self.remote_data_dir))
        self.run_ssh('mkdir -p {}'.format(self.remote_keystore_dir))

    def gen_node_tmp(self):
        """
        generate local node cache directory
        :return:
        """
        tmp = os.path.join(self.cfg.node_tmp, self.host + "_" + self.p2p_port)
        if not os.path.exists(tmp):
            os.makedirs(tmp)
        return tmp

    @property
    def enode(self):
        return r"enode://" + self.node_id + "@" + self.host + ":" + self.p2p_port

    def try_do(self, func):
        try:
            func()
        except Exception as e:
            raise Exception(failed_msg.format(self.node_mark, func.__name__, e))

    def try_do_resturn(self, func):
        try:
            func()
        except Exception as e:
            return False, failed_msg.format(self.node_mark, func.__name__, e)
        return True, success_msg.format(self.node_mark, func.__name__)

    def init(self):
        """
        Initialize
        :return:
        """
        def __init():
            cmd = '{} --datadir {} init {}'.format(self.remote_bin_file, self.remote_data_dir, self.remote_genesis_file)
            result = self.run_ssh(cmd)
            # todo ：fix init complete
            # Adding a query here can only alleviate the problem of starting deployment without initialization.
            self.run_ssh("ls {}".format(self.remote_data_dir))
            if len(result) > 0:
                log.error(failed_msg.format(self.node_mark, "init", result[0]))
                raise Exception("Init failed:{}".format(result[0]))
            log.debug("node-{} init success".format(self.node_mark))
        self.try_do(__init)

    def run_ssh(self, cmd, need_password=False):
        if need_password:
            return run_ssh(self.ssh, cmd, self.password)
        return run_ssh(self.ssh, cmd)

    def clean(self):
        """
        clear node data
        :return:
        """
        log.debug("Clean node:{}".format(self.node_mark))

        def __clean():
            is_success = self.stop()
            if not is_success:
                raise Exception("Stop failed")
            self.run_ssh("sudo -S -p '' rm -rf {};mkdir -p {}".format(self.remote_node_path, self.remote_node_path),
                         True)
            self.run_ssh("ls {}".format(self.remote_node_path))
        return self.try_do_resturn(__clean)

    def clean_db(self):
        """
        clear the node database
        :return:
        """
        def __clean_db():
            is_success = self.stop()
            if not is_success:
                raise Exception("Stop failed")
            self.run_ssh("sudo -S -p '' rm -rf {}".format(self.remote_db_dir), True)
        return self.try_do_resturn(__clean_db)

    def clean_log(self):
        """
        clear node log
        :return:
        """
        def __clean_log():
            is_success = self.stop()
            if not is_success:
                raise Exception("Stop failed")
            self.run_ssh("rm -rf {}".format(self.remote_log_dir))
            self.append_log_file()
        self.try_do(__clean_log)

    def append_log_file(self):
        """
        append log id
        :return:
        """
        def __append_log_file():
            self.run_ssh("mkdir -p {};echo {} >> {}/platon.log".format(self.remote_log_dir, self.cfg.env_id,
                                                                       self.remote_log_dir))
        self.try_do(__append_log_file)

    def stop(self):
        """
        close node
        :return:
        """
        log.debug("Stop node:{}".format(self.node_mark))

        def __stop():
            self.__is_connected = False
            self.__is_ws_connected = False
            if not self.running:
                return True, "{}-node is not running".format(self.node_mark)
            self.run_ssh("sudo -S -p '' supervisorctl stop {}".format(self.node_name), True)
        return self.try_do_resturn(__stop)

    def start(self, is_init=False) -> tuple:
        """
        boot node
        :param is_init:
        :return:
        """
        log.debug("Start node:{}".format(self.node_mark))

        def __start():
            is_success = self.stop()
            if not is_success:
                raise Exception("Stop failed")
            if is_init:
                self.init()
            self.append_log_file()
            result = self.run_ssh("sudo -S -p '' supervisorctl start " + self.node_name, True)
            for r in result:
                if "ERROR" in r or "Command 'supervisorctl' not found" in r:
                    raise Exception("Start failed:{}".format(r.strip("\n")))

        return self.try_do_resturn(__start)

    def restart(self) -> tuple:
        """
        restart node
        :return:
        """
        def __restart():
            self.append_log_file()
            result = self.run_ssh("sudo -S -p '' supervisorctl restart " + self.node_name, True)
            for r in result:
                if "ERROR" in r:
                    raise Exception("restart failed:{}".format(r.strip("\n")))
        return self.try_do_resturn(__restart)

    def update(self) -> tuple:
        """
        update node
        :return:
        """
        def __update():
            # todo fix me
            self.stop()
            self.put_bin()
            self.start()
        return self.try_do_resturn(__update)

    def close(self):
        """
        Close the node, delete the node data,
        delete the node supervisor configuration
        :return:
        """
        is_success = True
        msg = "close success"
        try:
            self.clean()
            self.run_ssh("sudo -S -p '' rm -rf /etc/supervisor/conf.d/{}.conf".format(self.node_name), True)
        except Exception as e:
            is_success = False
            msg = "{}-close failed:{}".format(self.node_mark, e)
        finally:
            self.t.close()
            return is_success, msg

    def put_bin(self):
        """
        upload binary package
        :return:
        """
        def __put_bin():
            self.run_ssh("rm -rf {}".format(self.remote_bin_file))
            self.sftp.put(self.cfg.platon_bin_file, self.remote_node_path)
            self.run_ssh('chmod +x {}'.format(self.remote_bin_file))
        self.try_do(__put_bin)

    def put_nodekey(self):
        """
        upload nodekey
        :return:
        """
        def __put_nodekey():
            nodekey_file = os.path.join(self.local_node_tmp, "nodekey")
            with open(nodekey_file, 'w', encoding="utf-8") as f:
                f.write(self.nodekey)
            self.run_ssh('mkdir -p {}'.format(self.remote_data_dir))
            self.sftp.put(nodekey_file, self.remote_nodekey_file)
        self.try_do(__put_nodekey)

    def put_blskey(self):
        """
        upload blskey
        :return:
        """
        def __put_blskey():
            blskey_file = os.path.join(self.local_node_tmp, "blskey")
            with open(blskey_file, 'w', encoding="utf-8") as f:
                f.write(self.blsprikey)
            self.run_ssh('mkdir -p {}'.format(self.remote_data_dir))
            self.sftp.put(blskey_file, self.remote_blskey_file)
        self.try_do(__put_blskey)

    def create_keystore(self, password="88888888"):
        """
        create a wallet
        :param password:
        :return:
        """
        def __create_keystore():
            cmd = "{} account new --datadir {}".format(self.remote_bin_file, self.remote_data_dir)
            stdin, stdout, _ = self.ssh.exec_command("source /etc/profile;%s" % cmd)
            stdin.write(str(password) + "\n")
            stdin.write(str(password) + "\n")
        self.try_do(__create_keystore)

    def put_genesis(self, genesis_file):
        """
        upload genesis
        :param genesis_file:
        :return:
        """
        def __put_genesis():
            self.run_ssh("rm -rf {}".format(self.remote_genesis_file))
            self.sftp.put(genesis_file, self.remote_genesis_file)
        self.try_do(__put_genesis)

    def put_config(self):
        """
        upload config
        :return:
        """
        def __put_config():
            self.run_ssh("rm -rf {}".format(self.remote_config_file))
            self.sftp.put(self.cfg.config_json_tmp, self.remote_config_file)
        self.try_do(__put_config)

    def put_static(self):
        """
        upload static
        :return:
        """
        def __put_static():
            self.sftp.put(self.cfg.static_node_tmp, self.remote_static_nodes_file)
        self.try_do(__put_static)

    def put_deploy_conf(self):
        """
        upload node deployment supervisor configuration
        :return:
        """
        def __put_deploy_conf():
            log.debug("{}-generate supervisor deploy conf...".format(self.node_mark))
            supervisor_tmp_file = os.path.join(self.local_node_tmp, "{}.conf".format(self.node_name))
            self.__gen_deploy_conf(supervisor_tmp_file)
            log.debug("{}-upload supervisor deploy conf...".format(self.node_mark))
            self.run_ssh("rm -rf {}".format(self.remote_supervisor_node_file))
            self.run_ssh("mkdir -p {}".format(self.cfg.remote_supervisor_tmp))
            self.sftp.put(supervisor_tmp_file, self.remote_supervisor_node_file)
            self.run_ssh("sudo -S -p '' cp {} /etc/supervisor/conf.d".format(self.remote_supervisor_node_file), True)
        self.try_do(__put_deploy_conf)

    def upload_file(self, local_file, remote_file):
        if local_file and os.path.exists(local_file):
            self.sftp.put(local_file, remote_file)
        else:
            log.info("file: {} not found".format(local_file))

    def __gen_deploy_conf(self, sup_tmp_file):
        """
        Generate a supervisor configuration for node deployment
        :param sup_tmp_file:
        :return:
        """
        with open(sup_tmp_file, "w") as fp:
            fp.write("[program:" + self.node_name + "]\n")
            go_fail_point = ""
            if self.fail_point:
                go_fail_point = " GO_FAILPOINTS='{}' ".format(self.fail_point)
            cmd = "{} --identity platon --datadir".format(self.remote_bin_file)
            cmd = cmd + " {} --port ".format(self.remote_data_dir) + self.p2p_port
            cmd = cmd + " --db.nogc"
            cmd = cmd + " --nodekey " + self.remote_nodekey_file
            cmd = cmd + " --cbft.blskey " + self.remote_blskey_file
            cmd = cmd + " --config " + self.remote_config_file
            cmd = cmd + " --syncmode '{}'".format(self.cfg.syncmode)
            cmd = cmd + " --debug --verbosity {}".format(self.cfg.log_level)
            if self.pprofport:
                cmd = cmd + " --pprof --pprofaddr 0.0.0.0 --pprofport " + str(self.pprofport)
            if self.wsport:
                cmd = cmd + " --ws --wsorigins '*' --wsaddr 0.0.0.0 --wsport " + str(self.wsport)
                cmd = cmd + " --wsapi platon,debug,personal,admin,net,web3"
            cmd = cmd + " --rpc --rpcaddr 0.0.0.0 --rpcport " + str(self.rpc_port)
            cmd = cmd + " --rpcapi platon,debug,personal,admin,net,web3,txpool"
            cmd = cmd + " --txpool.nolocals"
            if self.cfg.append_cmd:
                cmd = cmd + " " + self.cfg.append_cmd
            fp.write("command=" + cmd + "\n")
            if go_fail_point:
                fp.write("environment={}\n".format(go_fail_point))
            supervisor_default_conf = "numprocs=1\n" + "autostart=false\n" + "startsecs=1\n" + "startretries=3\n" + \
                                      "autorestart=unexpected\n" + "exitcode=0\n" + "stopsignal=TERM\n" + \
                                      "stopwaitsecs=10\n" + "redirect_stderr=true\n" + \
                                      "stdout_logfile_maxbytes=200MB\n" + "stdout_logfile_backups=20\n"
            fp.write(supervisor_default_conf)
            fp.write("stdout_logfile={}/platon.log\n".format(self.remote_log_dir))

    def deploy_me(self, genesis_file) -> tuple:
        """
        deploy this node
        1. Empty environmental data
        2. According to the node server to determine whether it is necessary to upload files
        3. Determine whether to initialize, choose to upload genesis
        4. Upload the node key file
        5. Upload the inter-node supervisor configuration
        6. Start the node
        :param genesis_file:
        :return:
        """
        log.debug("{}-clean node path...".format(self.node_mark))
        is_success, msg = self.clean()
        if not is_success:
            return is_success, msg
        self.clean_log()
        is_success, msg = self.put_all_file(genesis_file)
        if not is_success:
            return is_success, msg
        return self.start(self.cfg.init_chain)

    def put_all_file(self, genesis_file):
        """
        upload or copy the base file
        :param genesis_file:
        :return:
        """
        def __pre_env():
            ls = self.run_ssh("cd {};ls".format(self.cfg.remote_compression_tmp_path))
            if self.cfg.env_id and (self.cfg.env_id + ".tar.gz\n") in ls:
                log.debug("{}-copy bin...".format(self.remote_node_path))
                cmd = "cp -r {}/{}/* {}".format(self.cfg.remote_compression_tmp_path, self.cfg.env_id,
                                                self.remote_node_path)
                self.run_ssh(cmd)
                self.run_ssh("chmod +x {};mkdir {}".format(self.remote_bin_file, self.remote_log_dir))
            else:
                self.put_bin()
                self.put_config()
                # self.put_static()
                self.create_keystore()
            if self.cfg.init_chain:
                log.debug("{}-upload genesis...".format(self.node_mark))
                self.put_genesis(genesis_file)
            if self.cfg.is_need_static:
                self.put_static()
            log.debug("{}-upload blskey...".format(self.node_mark))
            self.put_blskey()
            log.debug("{}-upload nodekey...".format(self.node_mark))
            self.put_nodekey()
            self.put_deploy_conf()
            self.run_ssh("sudo -S -p '' supervisorctl update " + self.node_name, True)
        return self.try_do_resturn(__pre_env)

    def backup_log(self):
        """
        download log
        :return:
        """
        def __backup_log():
            self.run_ssh("cd {};tar zcvf log.tar.gz ./log".format(self.remote_node_path))
            self.sftp.get("{}/log.tar.gz".format(self.remote_node_path),
                          "{}/{}_{}.tar.gz".format(self.cfg.tmp_log, self.host, self.p2p_port))
            self.run_ssh("cd {};rm -rf ./log.tar.gz".format(self.remote_node_path))
        return self.try_do_resturn(__backup_log)

    @property
    def running(self) -> bool:
        p_id = self.run_ssh("ps -ef|grep platon|grep port|grep %s|grep -v grep|awk {'print $2'}" % self.p2p_port)
        if len(p_id) == 0:
            return False
        return True

    @property
    def web3(self) -> Web3:
        if not self.__is_connected:
            self.__rpc = wait_connect_web3(self.url, self.chain_id)
            self.__is_connected = True
        return self.__rpc

    @property
    def ws_web3(self) -> Web3:
        if not self.__is_ws_connected:
            self.__ws_rpc = wait_connect_web3(self.wsurl, self.chain_id)
            self.__is_ws_connected = True
        return self.__ws_rpc

    @property
    def eth(self) -> Eth:
        return Eth(self.web3)

    @property
    def admin(self) -> Admin:
        return Admin(self.web3)

    @property
    def debug(self) -> Debug:
        return Debug(self.web3)

    @property
    def personal(self) -> Personal:
        return Personal(self.web3)

    @property
    def ppos(self) -> Ppos:
        if not self.__is_connected_ppos:
            self.__ppos = Ppos(self.web3)
            self.__is_connected_ppos = True
        return self.__ppos

    @property
    def pip(self) -> Pip:
        if not self.__is_connected_pip:
            self.__pip = Pip(self.web3)
            self.__is_connected_pip = True
        return self.__pip

    @property
    def block_number(self) -> int:
        return self.eth.blockNumber

    @property
    def program_version(self):
        return self.admin.getProgramVersion()['Version']

    @property
    def program_version_sign(self):
        return self.admin.getProgramVersion()['Sign']

    @property
    def schnorr_NIZK_prove(self):
        return self.admin.getSchnorrNIZKProve()

    @property
    def staking_address(self):
        """
        staking wallet address
        """
        result = self.ppos.getCandidateInfo(self.node_id)
        candidate_info = result.get('Ret', {})
        address = candidate_info.get('StakingAddress')
        return self.web3.toChecksumAddress(address)

    @property
    def benifit_address(self):
        """
        staking wallet address
        """
        result = self.ppos.getCandidateInfo(self.node_id)
        candidate_info = result.get('Ret', {})
        address = candidate_info.get('BenefitAddress')
        return self.web3.toChecksumAddress(address)
