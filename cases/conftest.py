import pytest
import socket
import allure
import os
from common import download
from environment.env import create_env
from common.log import log


def set_report_env(allure_dir, env):
    node = env.get_rand_node()
    version_info_list = node.run_ssh("{} version".format(node.remote_bin_file))
    version_info = "".join(version_info_list).replace(" ", "").replace("Platon\n", "")
    allure_dir_env = os.path.join(allure_dir, "environment.properties")
    consensus_node = "ConsensusNodes:{}\n".format("|><|".join([node.node_mark for node in env.consensus_node_list]))
    normal_node = "NormalNodes:{}\n".format("|><|".join([node.node_mark for node in env.normal_node_list]))
    env_id = "TestEnvironmentID:{}\n".format(env.cfg.env_id)
    with open(allure_dir_env, "w", encoding="UTF-8")as f:
        f.write(version_info)
        f.write(consensus_node)
        f.write(normal_node)
        f.write(env_id)


@pytest.fixture(scope="module")
def consensus_test_env(global_test_env):
    with open("/etc/passwd") as f:
        yield f.readlines()


def pytest_addoption(parser):
    parser.addoption("--job", action="store", help="job: ci run job id")
    parser.addoption("--tmpDir", action="store", help="tmpDir: tmp dir, default in deploy/tmp/global")
    parser.addoption("--platonUrl", action="store", help="platonUrl: url to download platon bin")
    parser.addoption("--nodeFile", action="store", help="nodeFile: the node config file")
    parser.addoption("--accountFile", action="store", help="accountFile: the accounts file")
    parser.addoption("--initChain", action="store_true", default=True, dest="initChain", help="nodeConfig: default to init chain data")
    parser.addoption("--installDependency", action="store_true", default=False, dest="installDependency", help="installDependency: default do not install dependencies")
    parser.addoption("--installSupervisor", action="store_true", default=False, dest="installSuperVisor", help="installSupervisor: default do not install supervisor service")


# pytest 'tests/example/test_step.py' --nodeFile "deploy/node/debug_4_4.yml" --accountFile "deploy/accounts.yml" --alluredir="report/allure"
# --reruns 3
@pytest.fixture(scope="session", autouse=False)
def global_test_env(request, worker_id):
    log.info("start global_test_env>>>>>>>>>>>>>>")
    tmp_dir = request.config.getoption("--tmpDir")
    node_file = request.config.getoption("--nodeFile")
    account_file = request.config.getoption("--accountFile")
    init_chain = request.config.getoption("--initChain")
    install_dependency = request.config.getoption("--installDependency")
    install_supervisor = request.config.getoption("--installSupervisor")
    platon_url = request.config.getoption("--platonUrl")
    allure_dir = request.config.getoption("--alluredir")
    log.info(node_file)
    if worker_id != "master":
        if not node_file:
            raise Exception("The number of configuration files must be equal to the number of threads")
        node_file_list = node_file.split(",")
        for i in range(len(node_file_list)):
            if str(i) in worker_id:
                node_file = node_file_list[i]
                log.info("Session with nodeFile:{}".format(node_file_list[i]))
                tmp_dir = str(tmp_dir) + worker_id
    if platon_url:
        download.download_platon(platon_url)
    env = create_env(tmp_dir, node_file, account_file, init_chain, install_dependency, install_supervisor)
    # Must choose one, don't use both
    env.deploy_all()
    # env.prepare_all()
    yield env

    if allure_dir:
        set_report_env(allure_dir, env)

    # delete env and close env
    # env.shutdown()


@pytest.hookimpl(tryfirst=True, hookwrapper=True)
def pytest_runtest_makereport(item, call):
    # execute all other hooks to obtain the report object
    outcome = yield
    rep = outcome.get_result()
    # we only look at actual failing test calls, not setup/teardown
    if rep.when == "call" and not rep.passed:

        if 'global_test_env' in item.fixturenames:
            # download log in here
            try:
                log_name = item.funcargs["global_test_env"].backup_all_logs(item.name)
                job = item.funcargs["request"].config.getoption("--job")
                if job is None:
                    log_url = os.path.join(item.funcargs["global_test_env"].cfg.bug_log, log_name)
                else:
                    log_url = "http://{}:8080/job/PlatON/job/run/{}/artifact/logs/{}".format(socket.gethostbyname(socket.gethostname()), job, log_name)
                allure.attach('{}'.format(log_url), 'env log', allure.attachment_type.URI_LIST)
            except Exception as e:
                log.info("exception:{}".format(e))
            # Record block number
            try:
                if item.funcargs["global_test_env"].running:
                    env_status = "node blocks:{}".format(item.funcargs["global_test_env"].block_numbers())

                else:
                    env_status = "node runnings:{}".format(["{}:{}".format(node.node_mark, node.running) for node in
                                                            item.funcargs["global_test_env"].get_all_nodes()])
                log.info(env_status)
                allure.attach(env_status, "env status", allure.attachment_type.TEXT)
            except Exception as e:
                log.info("get block exception:{}".format(e))
        else:
            log.error("This case does not use global_test_env")
