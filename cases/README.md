# PlatON-Tests
This is an automated test project of the PaltON-Go

## Installation and operation dependencies
Install the python 3.7 environment and configure pip, then execute the following command to install the dependency library:
```shell script
pip install -r requirements.txt
```
     
## Run test

### Execute all test cases

```shell script
pytest test_start.py --nodeFile "deploy/4_node.yml" --accountFile "deploy/accounts.yml" --initChain
```

### Execute at Multiple environment

```shell script
pytest "case_path" --nodeFile "node_file_1,node_file_2" --accountFile "deploy/accounts.yml" --initChain -n 2
```

Note: The number of node configuration files must be equal to the number of threads, and multiple node configuration files are separated by English ","

## pytest Command line argument

--nodeFile "deploy/node.yml": Specify the node configuration file

--accountFile "deploy/accounts.yml": Specify the account file for testing

--installDependency：Indicates that the node needs to install the required dependencies, which is generally used during the first deployment; if it is not, it is no longer installed.

--installSuperVisor：Indicates whether the node is installed with the supervisor service. It is usually used for the first deployment. If you do not have this option, it is no longer installed.

## Precautions
Currently only supports Ubuntu environment deployment
File storage requirements:
    Accounts.yml file, put in the deploy directory, platon binary file into deploy/bin, nodeFile into deploy/node
    Other files, put in the deploy/template template directory

### test case example:
```python
import pytest
@pytest.mark.P1
def test_case_001():
    print("begin: test_case_001")
    SomeTxAPI("test_case_001")
    print("end: test_case_001")
```
    
