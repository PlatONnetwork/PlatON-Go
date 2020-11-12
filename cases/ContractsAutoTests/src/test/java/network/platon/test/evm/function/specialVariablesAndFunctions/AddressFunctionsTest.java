package network.platon.test.evm.function.specialVariablesAndFunctions;

import com.platon.sdk.utlis.Bech32;
import com.platon.sdk.utlis.NetworkParameters;
import network.platon.test.evm.beforetest.ContractPrepareTest;
import network.platon.autotest.junit.annotations.DataSource;
import network.platon.autotest.junit.enums.DataSourceType;
import network.platon.contracts.evm.AddressFunctions;
import org.junit.Before;
import org.junit.Test;
import org.web3j.protocol.core.methods.response.TransactionReceipt;

import java.math.BigInteger;

/**
 * @title 验证地址相关函数
 * @description:
 * @author: liweic
 * @create: 2020/01/02 11:30
 **/

public class AddressFunctionsTest extends ContractPrepareTest {
    private String amount;

    @Before
    public void before() {
        this.prepare();
        amount = driverService.param.get("amount");
    }

    @Test
    @DataSource(type = DataSourceType.EXCEL, file = "test.xls", sheetName = "Sheet1",
            author = "liweic", showName = "function.AddressFunctionsTest-地址相关函数测试", sourcePrefix = "evm")
    public void Addressfunctions() {
        try {
            AddressFunctions addressfunctions = AddressFunctions.deploy(web3j, transactionManager, provider, chainId).send();

            String contractAddress = addressfunctions.getContractAddress();
            TransactionReceipt tx = addressfunctions.getTransactionReceipt().get();
            collector.logStepPass("Addressfunctions deploy successfully.contractAddress:" + contractAddress + ", hash:" + tx.getTransactionHash());
            collector.logStepPass("Addressfunctions deploy gasUsed:" + addressfunctions.getTransactionReceipt().get().getGasUsed());

            //验证balance(地址账户)函数
            String balanceaddr = "0x03f0e0a226f081a5daecfda222cafc959ed7b800";
            balanceaddr = Bech32.addressEncode(NetworkParameters.TestNetParams.getHrp(),balanceaddr);
            BigInteger money = addressfunctions.getBalance(balanceaddr).send();
            collector.logStepPass("地址账户getBalance函数返回值：" + money);
            int num = money.toString().length();
            boolean n = num > 0;
            collector.assertEqual(n,true);

            //验证balance(合约账户)函数
            BigInteger contractmoney = addressfunctions.getBalanceOf().send();
            collector.logStepPass("合约账户getBalance函数返回值：" + contractmoney);
            int num2 = contractmoney.toString().length();
            collector.assertEqual(1,num2);

            //验证transfer函数
            String transferaddr = "0x8a9B36694F1eeeb500c84A19bB34137B05162EC1";
            transferaddr = Bech32.addressEncode(NetworkParameters.TestNetParams.getHrp(),transferaddr);
            BigInteger addresspremoney = addressfunctions.getBalance(transferaddr).send();
            collector.logStepPass("转账前余额：" + addresspremoney);
            TransactionReceipt result = addressfunctions.transfer(transferaddr ,new BigInteger(amount)).send();
            collector.logStepPass("transfer转账结果交易Hash：" + result.getTransactionHash());
            BigInteger addressaftermoney = addressfunctions.getBalance(transferaddr).send();
            collector.logStepPass("转账后余额：" + addressaftermoney);
            int a = Integer.valueOf(addressaftermoney.toString());
            int b = Integer.valueOf(addresspremoney.toString());
            int transfercounts = a - b;
            collector.assertEqual(amount ,String.valueOf(transfercounts));


            //验证send函数
            String sendsddr = "0x8a9B36694F1eeeb500c84A19bB34137B05162EC7";
            sendsddr = Bech32.addressEncode(NetworkParameters.TestNetParams.getHrp(),sendsddr);
            BigInteger sendbefore = addressfunctions.getBalance(sendsddr).send();
            collector.logStepPass("转账前余额：" + sendbefore);
            TransactionReceipt result2 = addressfunctions.send(sendsddr ,new BigInteger("10000")).send();
            collector.logStepPass("send转账结果交易Hash：" + result2.getTransactionHash());
            BigInteger addressaftersend = addressfunctions.getBalance(sendsddr).send();
            collector.logStepPass("转账后余额：" + addressaftersend);
            int c = Integer.valueOf(addressaftersend.toString());
            int d = Integer.valueOf(sendbefore.toString());
            int sendcount = c - d;
            collector.assertEqual(1,sendcount);

        } catch (Exception e) {
            collector.logStepFail("AddressFunctionsContract Calling Method fail.", e.toString());
            e.printStackTrace();
        }
    }
}
