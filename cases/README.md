# PlatON-Tests
This project is an automated test project of the PaltON-Go. see: https://github.com/PlatONnetwork/PlatON-Go

# 安装运行依赖
安装python3.6以上环境，并配置pip，然后执行以下命令安装依赖库：

pip install -r requirements.txt 

# Run test:

## 执行所有用例
pytest test_start.py --nodeFile "deploy/4_node.yml" --accountFile "deploy/accounts.yml" --initChain

## 多套环境并发执行
pytest "case_path" --nodeFile "node_file_1,node_file_2" --accountFile "deploy/accounts.yml" --initChain -n 2

备注：节点配置文件数必须等于线程数，多个节点配置文件用英文","分隔

# pytest 命令行参数

--nodeFile "deploy/4_node.yml":  指定节点配置文件

--accountFile "deploy/accounts.yml": 指定测试用的账号文件

--initChain：出现此选项，表示要初始化链数据；如果没有此选项，表示不初始化链数据

--installDependency：表示节点需要安装必需的依赖，一般第一次部署时使用；如果没有此选项，则不再安装

--installSuperVisor：表示节点是否安装supervisor服务，一般第一次部署时使用；如果没有此选项，则不再安装

# 注意事项
目前仅支持Ubuntu环境部署
文件存放要求：
    accounts.yml文件，放入deploy目录，platon二进制文件放入deploy/bin，nodeFile放入到deploy/node
    其它文件，放入deploy/template模板目录

@pytest.mark.P1
def test_case_3():
    print("begin: test_case_3")
    SomeTxAPI("test_case_3")
    print("end: test_case_3")
