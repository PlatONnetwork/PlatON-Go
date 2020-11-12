package network.platon.test.evm.data_type.BasicDataType;

import network.platon.test.evm.beforetest.ContractPrepareTest;
import network.platon.autotest.junit.annotations.DataSource;
import network.platon.autotest.junit.enums.DataSourceType;
import network.platon.contracts.evm.BasicDataTypeConstantContract;
import org.junit.Before;
import org.junit.Test;
import org.web3j.protocol.core.methods.response.TransactionReceipt;

import java.math.BigInteger;

/**
 * @title 测试：合约字面常量地址Address
 * @description:
 * @author: qudong
 * @create: 2019/12/25 15:09
 **/
public class BasicDataTypeConstantAddressTest extends ContractPrepareTest {

    private String amountValue1;
    private String amountValue2;
    private String amountSumValue;

    @Before
    public void before() {
       this.prepare();
        amountValue1 = driverService.param.get("amountValue1");
        amountValue2 = driverService.param.get("amountValue2");
        amountSumValue = driverService.param.get("amountSumValue");
    }

    @Test
    @DataSource(type = DataSourceType.EXCEL, file = "test.xls", author = "qudong", showName = "BasicDataTypeConstantContractTest.合约地址Address常量",sourcePrefix = "evm")
    public void testBasicDataTypeContract() {

        BasicDataTypeConstantContract basicDataTypeConstantContract = null;
        try {
            //合约部署
            basicDataTypeConstantContract = BasicDataTypeConstantContract.deploy(web3j, transactionManager, provider, chainId).send();
            String contractAddress = basicDataTypeConstantContract.getContractAddress();
            TransactionReceipt tx =  basicDataTypeConstantContract.getTransactionReceipt().get();
            collector.logStepPass("BasicDataTypeConstantContract issued successfully.contractAddress:" + contractAddress
                                    + ", hash:" + tx.getTransactionHash() + ",deploy gas used:" + tx.getGasUsed());
            collector.logStepPass("deployFinishCurrentBlockNumber:" + tx.getBlockNumber());
        } catch (Exception e) {
            collector.logStepFail("BasicDataTypeConstantContract deploy fail.", e.toString());
            e.printStackTrace();
        }

        //调用合约方法
        //1、验证：address类型，转账给指定地址
        try {
            int expectValue = Integer.parseInt(amountSumValue);
            String myContractAddress = "lax132dnv620rmht2qxgfgvmkdqn0vz3vtkfw87hjs";

            //1)、查询地址余额
            BigInteger contractBalance = basicDataTypeConstantContract.getBalance(myContractAddress).send();
            collector.logStepPass("BasicDataTypeConstantContract 执行 1)、查询地址余额 getBalance() successfully.contractBalance:" + contractBalance);

            //2）、查询当前转账人账户余额getBalance()
            BigInteger currentTransferUserBalance = basicDataTypeConstantContract.getBalance(walletAddress).send();
            collector.logStepPass("BasicDataTypeConstantContract 执行 2)、转账人账户余额 getBalance() successfully.currentTransferUserBalance: " + currentTransferUserBalance);

            //3)、给当前地址第一次转账
            TransactionReceipt transactionReceipt = basicDataTypeConstantContract.goTransfer(myContractAddress,new BigInteger(amountValue1)).send();
            collector.logStepPass("BasicDataTypeConstantContract 执行 3)、给当前地址转账  goTransfer() successfully hash:" + transactionReceipt.getTransactionHash());

            //4)、给当前地址第二次转账
            TransactionReceipt transactionReceipt1 = basicDataTypeConstantContract.goSend(myContractAddress,new BigInteger(amountValue2)).send();
            collector.logStepPass("BasicDataTypeConstantContract 执行 4)、给当前地址转账  goSend() successfully hash:" + transactionReceipt1.getTransactionHash());

            //5)、查询转账完毕，地址余额
            BigInteger currentContractBalance = basicDataTypeConstantContract.getBalance(myContractAddress).send();
            collector.logStepPass("BasicDataTypeConstantContract 执行 5)、查询转账完毕，地址余额 getBalance() successfully currentContractBalance:" + currentContractBalance);
            int actualValue = currentContractBalance.intValue() - contractBalance.intValue();

            collector.assertEqual(actualValue,expectValue, "checkout execute success.");

        } catch (Exception e) {
            collector.logStepFail("BasicDataTypeConstantContract Calling Method fail.", e.toString());
            e.printStackTrace();
        }
    }
}
