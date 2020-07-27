package evm.versioncompatible.v0_4_25;

import evm.beforetest.ContractPrepareTest;
import com.alibaba.fastjson.JSONObject;
import network.platon.autotest.junit.annotations.DataSource;
import network.platon.autotest.junit.enums.DataSourceType;
import network.platon.contracts.ContractAndAddressFunction;
import org.apache.commons.lang.StringUtils;
import org.junit.Test;
import org.web3j.protocol.core.DefaultBlockParameterName;
import org.web3j.protocol.core.methods.response.PlatonGetBalance;
import org.web3j.protocol.core.methods.response.TransactionReceipt;
import org.web3j.tuples.generated.Tuple3;
import org.web3j.tx.Transfer;
import org.web3j.utils.Convert;

import java.math.BigDecimal;
import java.math.BigInteger;

/**
 * @title 0.4.25版本合约和地址成员变量/函数测试
 * 1.0.4.25版本contract合约类型包括 address类型的成员函数，可以直接使用 send()成员函数验证
 * 2.0.4.25版本contract合约类型包括 address类型的成员函数，可以直接使用 transfer()成员函数验证
 * 3.0.4.25版本contract合约类型包括 address类型的成员函数，可以直接使用 balance成员变量验证
 * 4.0.4.25版本msg.sender类型所属验证
 * @description:
 * @author: albedo
 * @create: 2019/12/28
 */
public class ContractAndAddressFunctionTest extends ContractPrepareTest {


    @Test
    @DataSource(type = DataSourceType.EXCEL, file = "test.xls", sheetName = "testAddressCheck",
            author = "albedo", showName = "evm.ContractAndAddressTest-合约和地址成员变量(函数)", sourcePrefix = "evm")
    public void testAddressCheck() {
        try {
            prepare();
            ContractAndAddressFunction contractAndAddress = ContractAndAddressFunction.deploy(web3j, transactionManager, provider, chainId).send();
            String contractAddress = contractAndAddress.getContractAddress();
            String transactionHash = contractAndAddress.getTransactionReceipt().get().getTransactionHash();
            collector.logStepPass("ContractAndAddress issued successfully.contractAddress:" + contractAddress + ", hash:" + transactionHash);
            collector.logStepPass("deploy gas used:" + contractAndAddress.getTransactionReceipt().get().getGasUsed());

            Transfer transfer = new Transfer(web3j, transactionManager);
            TransactionReceipt receipt = transfer.sendFunds(contractAddress, BigDecimal.valueOf(1.00), Convert.Unit.LAT, new BigInteger(gasPrice), new BigInteger(gasLimit)).send();
            if (StringUtils.equals(receipt.getStatus(), "0x1")) {
                PlatonGetBalance balance = web3j.platonGetBalance(contractAddress, DefaultBlockParameterName.LATEST).send();
                collector.assertEqual(balance.getBalance(), new BigInteger("1000000000000000000"), "checkout to contract account transfer");
            } else {
                collector.logStepFail("transfer contract account is failure.contractAddress:", contractAddress);
            }
            Tuple3<String, BigInteger, BigInteger> result = contractAndAddress.addressCheck().send();
            Tuple3<String, BigInteger, BigInteger> expert = new Tuple3<>(receipt.getFrom(), new BigInteger("999999999999999980"), new BigInteger("20"));
            collector.assertEqual(JSONObject.toJSONString(result), JSONObject.toJSONString(expert), "checkout contract address function");
        } catch (Exception e) {
            collector.logStepFail("ContractAndAddressFunctionTest testAddressCheck failure,exception msg:", e.getMessage());
            e.printStackTrace();
        }
    }

    @Test
    @DataSource(type = DataSourceType.EXCEL, file = "test.xls", sheetName = "testVon",
            author = "albedo", showName = "ContractAndAddressTest-VON转账测试", sourcePrefix = "evm")
    public void testVon() {
        try {
            prepare();
            ContractAndAddressFunction contractAndAddress = ContractAndAddressFunction.deploy(web3j, transactionManager, provider, chainId).send();
            String contractAddress = contractAndAddress.getContractAddress();
            String transactionHash = contractAndAddress.getTransactionReceipt().get().getTransactionHash();
            collector.logStepPass("ContractAndAddress issued successfully.contractAddress:" + contractAddress + ", hash:" + transactionHash);
            collector.logStepPass("deploy gas used:" + contractAndAddress.getTransactionReceipt().get().getGasUsed());

            Transfer transfer = new Transfer(web3j, transactionManager);
            TransactionReceipt receipt = transfer.sendFunds(contractAddress, BigDecimal.valueOf(30.00), Convert.Unit.VON, new BigInteger(gasPrice), new BigInteger(gasLimit)).send();
            if (StringUtils.equals(receipt.getStatus(), "0x1")) {
                PlatonGetBalance balance = web3j.platonGetBalance(contractAddress, DefaultBlockParameterName.LATEST).send();
                collector.assertEqual(balance.getBalance(), new BigInteger("30"), "checkout to contract account transfer");
            } else {
                collector.logStepFail("transfer contract account is failure.contractAddress:", contractAddress);
            }
            Tuple3<String, BigInteger, BigInteger> result = contractAndAddress.addressCheck().send();
            Tuple3<String, BigInteger, BigInteger> expert = new Tuple3<>(receipt.getFrom(), new BigInteger("10"), new BigInteger("20"));
            collector.assertEqual(JSONObject.toJSONString(result), JSONObject.toJSONString(expert), "checkout VON transfer function");
        } catch (Exception e) {
            collector.logStepFail("ContractAndAddressFunctionTest testAddressCheck failure,exception msg:", e.getMessage());
            e.printStackTrace();
        }
    }

    @Test
    @DataSource(type = DataSourceType.EXCEL, file = "test.xls", sheetName = "testKVon",
            author = "albedo", showName = "ContractAndAddressTest-KVON转账测试", sourcePrefix = "evm")
    public void testKVon() {
        try {
            prepare();
            ContractAndAddressFunction contractAndAddress = ContractAndAddressFunction.deploy(web3j, transactionManager, provider, chainId).send();
            String contractAddress = contractAndAddress.getContractAddress();
            String transactionHash = contractAndAddress.getTransactionReceipt().get().getTransactionHash();
            collector.logStepPass("ContractAndAddress issued successfully.contractAddress:" + contractAddress + ", hash:" + transactionHash);
            collector.logStepPass("deploy gas used:" + contractAndAddress.getTransactionReceipt().get().getGasUsed());

            Transfer transfer = new Transfer(web3j, transactionManager);
            TransactionReceipt receipt = transfer.sendFunds(contractAddress, BigDecimal.valueOf(1.00), Convert.Unit.KVON, new BigInteger(gasPrice), new BigInteger(gasLimit)).send();
            if (StringUtils.equals(receipt.getStatus(), "0x1")) {
                PlatonGetBalance balance = web3j.platonGetBalance(contractAddress, DefaultBlockParameterName.LATEST).send();
                collector.assertEqual(balance.getBalance(), new BigInteger("1000"), "checkout to contract account transfer");
            } else {
                collector.logStepFail("transfer contract account is failure.contractAddress:", contractAddress);
            }
            Tuple3<String, BigInteger, BigInteger> result = contractAndAddress.addressCheck().send();
            Tuple3<String, BigInteger, BigInteger> expert = new Tuple3<>(receipt.getFrom(), new BigInteger("980"), new BigInteger("20"));
            collector.assertEqual(JSONObject.toJSONString(result), JSONObject.toJSONString(expert), "checkout KVON transfer function");
        } catch (Exception e) {
            collector.logStepFail("ContractAndAddressFunctionTest testAddressCheck failure,exception msg:", e.getMessage());
            e.printStackTrace();
        }
    }

    @Test
    @DataSource(type = DataSourceType.EXCEL, file = "test.xls", sheetName = "testMVon",
            author = "albedo", showName = "ContractAndAddressTest-MVON转账测试", sourcePrefix = "evm")
    public void testMVon() {
        try {
            prepare();
            ContractAndAddressFunction contractAndAddress = ContractAndAddressFunction.deploy(web3j, transactionManager, provider, chainId).send();
            String contractAddress = contractAndAddress.getContractAddress();
            String transactionHash = contractAndAddress.getTransactionReceipt().get().getTransactionHash();
            collector.logStepPass("ContractAndAddress issued successfully.contractAddress:" + contractAddress + ", hash:" + transactionHash);
            collector.logStepPass("deploy gas used:" + contractAndAddress.getTransactionReceipt().get().getGasUsed());

            Transfer transfer = new Transfer(web3j, transactionManager);
            TransactionReceipt receipt = transfer.sendFunds(contractAddress, BigDecimal.valueOf(0.001), Convert.Unit.MVON, new BigInteger(gasPrice), new BigInteger(gasLimit)).send();
            if (StringUtils.equals(receipt.getStatus(), "0x1")) {
                PlatonGetBalance balance = web3j.platonGetBalance(contractAddress, DefaultBlockParameterName.LATEST).send();
                collector.assertEqual(balance.getBalance(), new BigInteger("1000"), "checkout to contract account transfer");
            } else {
                collector.logStepFail("transfer contract account is failure.contractAddress:", contractAddress);
            }
            Tuple3<String, BigInteger, BigInteger> result = contractAndAddress.addressCheck().send();
            Tuple3<String, BigInteger, BigInteger> expert = new Tuple3<>(receipt.getFrom(), new BigInteger("980"), new BigInteger("20"));
            collector.assertEqual(JSONObject.toJSONString(result), JSONObject.toJSONString(expert), "checkout MVON transfer function");
        } catch (Exception e) {
            collector.logStepFail("ContractAndAddressFunctionTest testAddressCheck failure,exception msg:", e.getMessage());
            e.printStackTrace();
        }
    }

    @Test
    @DataSource(type = DataSourceType.EXCEL, file = "test.xls", sheetName = "testGVon",
            author = "albedo", showName = "ContractAndAddressTest-GVON转账测试", sourcePrefix = "evm")
    public void testGVon() {
        try {
            prepare();
            ContractAndAddressFunction contractAndAddress = ContractAndAddressFunction.deploy(web3j, transactionManager, provider, chainId).send();
            String contractAddress = contractAndAddress.getContractAddress();
            String transactionHash = contractAndAddress.getTransactionReceipt().get().getTransactionHash();
            collector.logStepPass("ContractAndAddress issued successfully.contractAddress:" + contractAddress + ", hash:" + transactionHash);
            collector.logStepPass("deploy gas used:" + contractAndAddress.getTransactionReceipt().get().getGasUsed());


            Transfer transfer = new Transfer(web3j, transactionManager);
            TransactionReceipt receipt = transfer.sendFunds(contractAddress, BigDecimal.valueOf(0.0000001), Convert.Unit.GVON, new BigInteger(gasPrice), new BigInteger(gasLimit)).send();
            if (StringUtils.equals(receipt.getStatus(), "0x1")) {
                PlatonGetBalance balance = web3j.platonGetBalance(contractAddress, DefaultBlockParameterName.LATEST).send();
                collector.assertEqual(balance.getBalance(), new BigInteger("100"), "checkout to contract account transfer");
            } else {
                collector.logStepFail("transfer contract account is failure.contractAddress:", contractAddress);
            }
            Tuple3<String, BigInteger, BigInteger> result = contractAndAddress.addressCheck().send();
            Tuple3<String, BigInteger, BigInteger> expert = new Tuple3<>(receipt.getFrom(), new BigInteger("80"), new BigInteger("20"));
            collector.assertEqual(JSONObject.toJSONString(result), JSONObject.toJSONString(expert), "checkout GVON transfer function");
        } catch (Exception e) {
            collector.logStepFail("ContractAndAddressFunctionTest testGVon failure,exception msg:", e.getMessage());
            e.printStackTrace();
        }
    }

    @Test
    @DataSource(type = DataSourceType.EXCEL, file = "test.xls", sheetName = "testSzabo",
            author = "albedo", showName = "ContractAndAddressTest-MICROLAT转账测试", sourcePrefix = "evm")
    public void testMicroLat() {
        try {
            prepare();
            ContractAndAddressFunction contractAndAddress = ContractAndAddressFunction.deploy(web3j, transactionManager, provider, chainId).send();
            String contractAddress = contractAndAddress.getContractAddress();
            String transactionHash = contractAndAddress.getTransactionReceipt().get().getTransactionHash();
            collector.logStepPass("ContractAndAddress issued successfully.contractAddress:" + contractAddress + ", hash:" + transactionHash);
            collector.logStepPass("deploy gas used:" + contractAndAddress.getTransactionReceipt().get().getGasUsed());


            Transfer transfer = new Transfer(web3j, transactionManager);
            TransactionReceipt receipt = transfer.sendFunds(contractAddress, BigDecimal.valueOf(0.0000000001), Convert.Unit.MICROLAT, new BigInteger(gasPrice), new BigInteger(gasLimit)).send();
            if (StringUtils.equals(receipt.getStatus(), "0x1")) {
                PlatonGetBalance balance = web3j.platonGetBalance(contractAddress, DefaultBlockParameterName.LATEST).send();
                collector.assertEqual(balance.getBalance(), new BigInteger("100"), "checkout to contract account transfer");
            } else {
                collector.logStepFail("transfer contract account is failure.contractAddress:", contractAddress);
            }
            Tuple3<String, BigInteger, BigInteger> result = contractAndAddress.addressCheck().send();
            Tuple3<String, BigInteger, BigInteger> expert = new Tuple3<>(receipt.getFrom(), new BigInteger("80"), new BigInteger("20"));
            collector.assertEqual(JSONObject.toJSONString(result), JSONObject.toJSONString(expert), "checkout SZABO transfer function");
        } catch (Exception e) {
            collector.logStepFail("ContractAndAddressFunctionTest testSzabo failure,exception msg:", e.getMessage());
            e.printStackTrace();
        }
    }

    @Test
    @DataSource(type = DataSourceType.EXCEL, file = "test.xls", sheetName = "testFinney",
            author = "albedo", showName = "ContractAndAddressTest-MILLILAT转账测试", sourcePrefix = "evm")
    public void testMilliLat() {
        try {
            prepare();
            ContractAndAddressFunction contractAndAddress = ContractAndAddressFunction.deploy(web3j, transactionManager, provider, chainId).send();
            String contractAddress = contractAndAddress.getContractAddress();
            String transactionHash = contractAndAddress.getTransactionReceipt().get().getTransactionHash();
            collector.logStepPass("ContractAndAddress issued successfully.contractAddress:" + contractAddress + ", hash:" + transactionHash);
            collector.logStepPass("deploy gas used:" + contractAndAddress.getTransactionReceipt().get().getGasUsed());


            Transfer transfer = new Transfer(web3j, transactionManager);
            TransactionReceipt receipt = transfer.sendFunds(contractAddress, BigDecimal.valueOf(0.000000000000021), Convert.Unit.MILLILAT, new BigInteger(gasPrice), new BigInteger(gasLimit)).send();
            if (StringUtils.equals(receipt.getStatus(), "0x1")) {
                PlatonGetBalance balance = web3j.platonGetBalance(contractAddress, DefaultBlockParameterName.LATEST).send();
                collector.assertEqual(balance.getBalance(), new BigInteger("21"), "checkout to contract account transfer");
            } else {
                collector.logStepFail("transfer contract account is failure.contractAddress:", contractAddress);
            }
            Tuple3<String, BigInteger, BigInteger> result = contractAndAddress.addressCheck().send();
            Tuple3<String, BigInteger, BigInteger> expert = new Tuple3<>(receipt.getFrom(), new BigInteger("1"), new BigInteger("20"));
            collector.assertEqual(JSONObject.toJSONString(result), JSONObject.toJSONString(expert), "checkout FINNEY transfer function");
        } catch (Exception e) {
            collector.logStepFail("ContractAndAddressFunctionTest testFinney failure,exception msg:", e.getMessage());
            e.printStackTrace();
        }
    }

    @Test
    @DataSource(type = DataSourceType.EXCEL, file = "test.xls", sheetName = "testLat",
            author = "albedo", showName = "ContractAndAddressTest-LAT转账测试", sourcePrefix = "evm")
    public void testLat() {
        try {
            prepare();
            ContractAndAddressFunction contractAndAddress = ContractAndAddressFunction.deploy(web3j, transactionManager, provider, chainId).send();
            String contractAddress = contractAndAddress.getContractAddress();
            String transactionHash = contractAndAddress.getTransactionReceipt().get().getTransactionHash();
            collector.logStepPass("ContractAndAddress issued successfully.contractAddress:" + contractAddress + ", hash:" + transactionHash);
            collector.logStepPass("deploy gas used:" + contractAndAddress.getTransactionReceipt().get().getGasUsed());


            Transfer transfer = new Transfer(web3j, transactionManager);
            TransactionReceipt receipt = transfer.sendFunds(contractAddress, BigDecimal.valueOf(0.00000000000000005), Convert.Unit.LAT, new BigInteger(gasPrice), new BigInteger(gasLimit)).send();
            if (StringUtils.equals(receipt.getStatus(), "0x1")) {
                PlatonGetBalance balance = web3j.platonGetBalance(contractAddress, DefaultBlockParameterName.LATEST).send();
                collector.assertEqual(balance.getBalance(), new BigInteger("50"), "checkout to contract account transfer");
            } else {
                collector.logStepFail("transfer contract account is failure.contractAddress:", contractAddress);
            }
            Tuple3<String, BigInteger, BigInteger> result = contractAndAddress.addressCheck().send();
            Tuple3<String, BigInteger, BigInteger> expert = new Tuple3<>(receipt.getFrom(), new BigInteger("30"), new BigInteger("20"));
            collector.assertEqual(JSONObject.toJSONString(result), JSONObject.toJSONString(expert), "checkout LAT transfer function");
        } catch (Exception e) {
            collector.logStepFail("ContractAndAddressFunctionTest testLat failure,exception msg:", e.getMessage());
            e.printStackTrace();
        }
    }

    @Test
    @DataSource(type = DataSourceType.EXCEL, file = "test.xls", sheetName = "testKLat",
            author = "albedo", showName = "ContractAndAddressTest-KLAT转账测试", sourcePrefix = "evm")
    public void testKLat() {
        try {
            prepare();
            ContractAndAddressFunction contractAndAddress = ContractAndAddressFunction.deploy(web3j, transactionManager, provider, chainId).send();
            String contractAddress = contractAndAddress.getContractAddress();
            String transactionHash = contractAndAddress.getTransactionReceipt().get().getTransactionHash();
            collector.logStepPass("ContractAndAddress issued successfully.contractAddress:" + contractAddress + ", hash:" + transactionHash);
            collector.logStepPass("deploy gas used:" + contractAndAddress.getTransactionReceipt().get().getGasUsed());

            Transfer transfer = new Transfer(web3j, transactionManager);
            TransactionReceipt receipt = transfer.sendFunds(contractAddress, BigDecimal.valueOf(0.0000000000000000001), Convert.Unit.KLAT, new BigInteger(gasPrice), new BigInteger(gasLimit)).send();
            if (StringUtils.equals(receipt.getStatus(), "0x1")) {
                PlatonGetBalance balance = web3j.platonGetBalance(contractAddress, DefaultBlockParameterName.LATEST).send();
                collector.assertEqual(balance.getBalance(), new BigInteger("100"), "checkout to contract account transfer");
            } else {
                collector.logStepFail("transfer contract account is failure.contractAddress:", contractAddress);
            }
            Tuple3<String, BigInteger, BigInteger> result = contractAndAddress.addressCheck().send();
            Tuple3<String, BigInteger, BigInteger> expert = new Tuple3<>(receipt.getFrom(), new BigInteger("80"), new BigInteger("20"));
            collector.assertEqual(JSONObject.toJSONString(result), JSONObject.toJSONString(expert), "checkout KLAT transfer function");
        } catch (Exception e) {
            collector.logStepFail("ContractAndAddressFunctionTest testKLat failure,exception msg:", e.getMessage());
            e.printStackTrace();
        }
    }

    @Test
    @DataSource(type = DataSourceType.EXCEL, file = "test.xls", sheetName = "testMLat",
            author = "albedo", showName = "ContractAndAddressTest-MLAT转账测试", sourcePrefix = "evm")
    public void testMLat() {
        try {
            prepare();
            ContractAndAddressFunction contractAndAddress = ContractAndAddressFunction.deploy(web3j, transactionManager, provider, chainId).send();
            String contractAddress = contractAndAddress.getContractAddress();
            String transactionHash = contractAndAddress.getTransactionReceipt().get().getTransactionHash();
            collector.logStepPass("ContractAndAddress issued successfully.contractAddress:" + contractAddress + ", hash:" + transactionHash);
            collector.logStepPass("deploy gas used:" + contractAndAddress.getTransactionReceipt().get().getGasUsed());


            Transfer transfer = new Transfer(web3j, transactionManager);
            TransactionReceipt receipt = transfer.sendFunds(contractAddress, BigDecimal.valueOf(0.00000000000000000000001), Convert.Unit.MLAT, new BigInteger(gasPrice), new BigInteger(gasLimit)).send();
            if (StringUtils.equals(receipt.getStatus(), "0x1")) {
                PlatonGetBalance balance = web3j.platonGetBalance(contractAddress, DefaultBlockParameterName.LATEST).send();
                collector.assertEqual(balance.getBalance(), new BigInteger("10"), "checkout to contract account transfer");
            } else {
                collector.logStepFail("transfer contract account is failure.contractAddress:", contractAddress);
            }
            Tuple3<String, BigInteger, BigInteger> result = contractAndAddress.addressCheck().send();
            Tuple3<String, BigInteger, BigInteger> expert = new Tuple3<>(receipt.getFrom(), new BigInteger("0"), new BigInteger("10"));
            collector.assertEqual(JSONObject.toJSONString(result), JSONObject.toJSONString(expert), "checkout MLAT transfer function");
        } catch (Exception e) {
            collector.logStepFail("ContractAndAddressFunctionTest testMLat failure,exception msg:", e.getMessage());
            e.printStackTrace();
        }
    }

    @Test
    @DataSource(type = DataSourceType.EXCEL, file = "test.xls", sheetName = "testGLat",
            author = "albedo", showName = "ContractAndAddressTest-GLAT转账测试", sourcePrefix = "evm")
    public void testGLat() {
        try {
            prepare();
            ContractAndAddressFunction contractAndAddress = ContractAndAddressFunction.deploy(web3j, transactionManager, provider, chainId).send();
            String contractAddress = contractAndAddress.getContractAddress();
            String transactionHash = contractAndAddress.getTransactionReceipt().get().getTransactionHash();
            collector.logStepPass("ContractAndAddress issued successfully.contractAddress:" + contractAddress + ", hash:" + transactionHash);
            collector.logStepPass("deploy gas used:" + contractAndAddress.getTransactionReceipt().get().getGasUsed());

            Transfer transfer = new Transfer(web3j, transactionManager);
            TransactionReceipt receipt = transfer.sendFunds(contractAddress, BigDecimal.valueOf(0.0000000000000000000000001), Convert.Unit.GLAT, new BigInteger(gasPrice), new BigInteger(gasLimit)).send();
            if (StringUtils.equals(receipt.getStatus(), "0x1")) {
                PlatonGetBalance balance = web3j.platonGetBalance(contractAddress, DefaultBlockParameterName.LATEST).send();
                collector.assertEqual(balance.getBalance(), new BigInteger("100"), "checkout to contract account transfer");
            } else {
                collector.logStepFail("transfer contract account is failure.contractAddress:", contractAddress);
            }
            Tuple3<String, BigInteger, BigInteger> result = contractAndAddress.addressCheck().send();
            Tuple3<String, BigInteger, BigInteger> expert = new Tuple3<>(receipt.getFrom(), new BigInteger("80"), new BigInteger("20"));
            collector.assertEqual(JSONObject.toJSONString(result), JSONObject.toJSONString(expert), "checkout GLAT transfer function");
        } catch (Exception e) {
            collector.logStepFail("ContractAndAddressFunctionTest testGLat failure,exception msg:", e.getMessage());
            e.printStackTrace();
        }
    }


}
