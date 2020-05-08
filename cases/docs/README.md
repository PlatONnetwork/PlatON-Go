# PlatON-Tests

### Test execution

1. Dependent installation 

```js
pip install -r requirements.txt
```

2. Execute the test

Follow the pytest command to execute, add a new parameter reference[conftest.py](../conftest.py)。

example: 

```js
py.test 'tests/chain/test_chain_deploy.py' --nodeFile "deploy/node/test_chaininfo.yml" --accountFile "deploy/accounts.yml"
```

More please check[readme](../README.md)

Remarks:
There are many combinations of project management and execution test commands. 
You need to input a large number of commands and consider introducing a Makefile.

### Project structure description

1. Root directory:

+ pytest.ini pytestTest the configuration file. For details, please refer to the official file.
[pytest.ini](https://docs.pytest.org/en/latest/reference.html#configuration-options)。

+ conftest.py-pytestPlug-ins, mainly including the introduction of command line parameters, 
test environment generation, test failure use case log download and other functions.

2.common：

+ Project public methods, including logging, configuration file reading, connecting to linux server, 
establishing web3 connection, binary package download, etc.

3.conf：

+ Local file path configuration when deploying the environment, please see the file notes for details.

4.deploy：

+ Deployment file storage path, used with conf

5.environment：

+ Deployment environment script, consisting of five classes: TestConfig, TestEnvironment, Account, Server, Node。
+ TestConfig: Deployment environment configuration, test environment global configuration。
+ Account: Test account management class for testing, account generation, transfer, etc.。
+ Server: Server management class for deployment, management server dependencies, and compressed file management.
+ TestEnvironment: Test environment class, used to manage the test environment, start the node, close the node, restart the node, get the node, get the account, and so on.
+ Node: Node class, used to manage nodes, with basic node information, node connections, node running status, node start and stop, and so on.

6.tests：

+ Use case storage directory

7.tests/lib:

+ Utils: use case public method
+ Configuration: StakingConfig, pledge basic configuration; PipConfig, pip basic configuration; DefaultEconomicConfig economic model basic configuration
+ genesi：Genesis Genesis.json file convertible object
+ The rest of the classes are secondary packages for sdk


### Example

[staking example](../tests/example/test_staking.py)

[Custom node deployment](../tests/example/test_customize_deploy.py)

### Feature description

+ Abandon the deployment form in the old framework, introduce test environment objects and node objects, used to manage the test environment and nodes, including a large number of node attributes in the node, eliminating the need to retrieve a large number of node information

+ The w3 connection establishment introduces a timeout mechanism to avoid the problem that the node just starts to connect and causes an error.

+ Remove the setup, introduce the fixture, design the use case, the common steps need to be designed as a step method, remove the duplicate step code

+ Deployment optimization, sacrificing environment preparation time, and greatly improving the speed of secondary deployment