import configparser
import os

from common.connect import run_ssh, connect_linux
from environment.config import TestConfig
from common.log import log
from environment.mock import mock_connect_linux


class Server:
    def __init__(self, server_conf, cfg: TestConfig):
        self.cfg = cfg
        self.host = server_conf["host"]
        self.username = server_conf["username"]
        self.password = server_conf["password"]
        self.ssh_port = server_conf.get("sshport", 22)
        if self.cfg.can_deploy:
            self.ssh, self.sftp, self.t = connect_linux(self.host, self.username, self.password, self.ssh_port)
        else:
            self.ssh, self.sftp, self.t = mock_connect_linux()
        self.remote_supervisor_conf = "{}/supervisord.conf".format(self.cfg.remote_supervisor_tmp)

    def run_ssh(self, cmd, need_password=False):
        if need_password:
            return run_ssh(self.ssh, cmd, self.password)
        return run_ssh(self.ssh, cmd)

    def clean_supervisor_conf(self):
        self.run_ssh("sudo -S -p '' rm -rf /etc/supervisor/conf.d/node-*", True)

    def put_compression(self):
        try:
            ls = self.run_ssh("cd {};ls".format(self.cfg.remote_compression_tmp_path))
            gz_name = self.cfg.env_id + ".tar.gz"
            local_gz = os.path.join(self.cfg.env_tmp, gz_name)
            if (gz_name + "\n") in ls:
                return True, "need not upload"
            self.run_ssh("rm -rf {};mkdir -p {}".format(self.cfg.remote_compression_tmp_path,
                                                        self.cfg.remote_compression_tmp_path))
            self.sftp.put(local_gz, self.cfg.remote_compression_tmp_path + "/" + os.path.basename(local_gz))
            self.run_ssh("tar -zxvf {}/{}.tar.gz -C {}".format(self.cfg.remote_compression_tmp_path, self.cfg.env_id,
                                                               self.cfg.remote_compression_tmp_path))
        except Exception as e:
            return False, "{}-upload compression failed:{}".format(self.host, e)
        return True, "upload compression success"

    def install_dependency(self):
        try:
            self.run_ssh("sudo -S -p '' ntpdate 0.centos.pool.ntp.org", True)
            self.run_ssh("sudo -S -p '' apt install llvm g++ libgmp-dev libssl-dev -y", True)
        except Exception as e:
            return False, "{}-install dependency failed:{}".format(self.host, e)
        return True, "install dependency success"

    def install_supervisor(self):
        try:
            test_name = "test-node"
            result = self.run_ssh("sudo -S -p '' supervisorctl stop {}".format(test_name), True)
            if len(result) == 0 or test_name not in result[0]:
                tmp_dir = os.path.join(self.cfg.server_tmp, self.host)
                if not os.path.exists(tmp_dir):
                    os.makedirs(tmp_dir)
                tmp = os.path.join(tmp_dir, "supervisord.conf")
                self.__rewrite_supervisor_conf(tmp)
                self.run_ssh("mkdir -p {}".format(self.cfg.remote_supervisor_tmp))
                self.sftp.put(tmp, self.remote_supervisor_conf)
                supervisor_pid_str = self.run_ssh("ps -ef|grep supervisord|grep -v grep|awk {'print $2'}")
                if len(supervisor_pid_str) > 0:
                    self.__reload_supervisor(supervisor_pid_str)
                else:
                    self.run_ssh("sudo -S -p '' apt update", True)
                    self.run_ssh("sudo -S -p '' apt install -y supervisor", True)
                    self.run_ssh("sudo -S -p '' cp {} /etc/supervisor/".format(self.remote_supervisor_conf), True)
                    supervisor_pid_str = self.run_ssh("ps -ef|grep supervisord|grep -v grep|awk {'print $2'}")
                    if len(supervisor_pid_str) > 0:
                        self.__reload_supervisor(supervisor_pid_str)
                    else:
                        self.run_ssh("sudo -S -p '' /etc/init.d/supervisor start", True)
        except Exception as e:
            return False, "{}-install supervisor failed:{}".format(self.host, e)
        return True, "install supervisor success"

    def __reload_supervisor(self, supervisor_pid_str):
        supervisor_pid = supervisor_pid_str[0].strip("\n")
        self.run_ssh("sudo -S -p '' kill {}".format(supervisor_pid), True)
        self.run_ssh("sudo -S -p '' killall supervisord", True)
        self.run_ssh("sudo -S -p '' sudo apt remove supervisor -y", True)
        self.run_ssh("sudo -S -p '' apt update", True)
        self.run_ssh("sudo -S -p '' apt install -y supervisor", True)
        self.run_ssh("sudo -S -p '' cp {} /etc/supervisor/".format(self.remote_supervisor_conf), True)
        self.run_ssh("sudo -S -p '' /etc/init.d/supervisor start", True)

    def __rewrite_supervisor_conf(self, sup_tmp):
        con = configparser.ConfigParser()
        con.read(self.cfg.supervisor_file)
        con.set("inet_http_server", "username", self.username)
        con.set("inet_http_server", "password", self.password)
        con.set("supervisorctl", "username", self.username)
        con.set("supervisorctl", "password", self.password)
        with open(sup_tmp, "w") as file:
            con.write(file)
