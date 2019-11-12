# PlatON-Tests
This project is an automated test project of the PaltON-Go. see: https://github.com/PlatONnetwork/PlatON-Go

# 安装运行依赖
安装python3.6以上环境，并配置pip，然后执行以下命令安装依赖库：

pip install -r requirements.txt 

# Run test:
## 以并发方式，执行所有用例
py.test test_start.py -s --concmode=asyncnet --nodeFile "deploy/4_node.yml" --accountFile "deploy/accounts.yml" --initChain

## 以同步方式，执行所有用例
py.test test_start.py -s --nodeFile "deploy/4_node.yml" --accountFile "deploy/accounts.yml" --initChain


## 以并发方式，执行所有并发的用例
py.test test_start.py -s -m "not SYNC" --concmode=asyncnet --nodeFile "deploy/4_node.yml" --accountFile "deploy/accounts.yml" --initChain 

## 以同步方式，执行所有同步的用例
py.test test_start.py -s -m "SYNC" --nodeFile "deploy/4_node.yml" --accountFile "deploy/accounts.yml" --initChain

# py.test 命令行参数
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

用例书写：
如果用例不支持并发方式运行，则用@pytest.mark.SYNC标注，如果不家标注，则默认是可以并发运行的，如：

@pytest.mark.P1
@pytest.mark.SYNC
def test_case_3():
    print("begin: test_case_3")
    SomeTxAPI("test_case_3")
    print("end: test_case_3")
