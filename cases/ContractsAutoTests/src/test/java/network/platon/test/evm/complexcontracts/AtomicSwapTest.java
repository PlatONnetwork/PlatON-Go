package network.platon.test.evm.complexcontracts;

import network.platon.test.evm.beforetest.ContractPrepareTest;
import network.platon.autotest.junit.annotations.DataSource;
import network.platon.autotest.junit.enums.DataSourceType;
import network.platon.contracts.evm.AtomicSwap;
import org.junit.Before;
import org.junit.Test;
import org.web3j.protocol.core.methods.response.TransactionReceipt;

import java.math.BigInteger;
import java.util.List;

public class AtomicSwapTest extends ContractPrepareTest {
    @Before
    public void before() {
        this.prepare();
    }

    @Test
    @DataSource(type = DataSourceType.EXCEL, file = "test.xls", sheetName = "Sheet1",
            author = "qcxiao", showName = "complexcontracts.AtomicSwapTest", sourcePrefix = "evm")
    public void test() {
        try {
            AtomicSwap atomicSwap = AtomicSwap.deploy(web3j, transactionManager, provider, chainId).send();
            String contractAddress = atomicSwap.getContractAddress();
            TransactionReceipt tx = atomicSwap.getTransactionReceipt().get();

            collector.logStepPass(
                    "Token issued successfully.contractAddress:" + contractAddress +
                            ", hash:" + tx.getTransactionHash());
            collector.logStepPass("deploy gas used:" + atomicSwap.getTransactionReceipt().get().getGasUsed());


            /*
             * 方法作用：初始化一个map并且给定这个map的key转入多少金额
             * 参数：key、参与者地址(目前0.12.1的版本需要改为0x...)、时间戳(必须大于当前区块时间)、初始减去多少、初始转多少VON
             */
            tx = atomicSwap.initiate(stringToBytes32("hello"), "lax1uqug0zq7rcxddndleq4ux2ft3tv6dqljphydrl", new BigInteger("1690063146613"), BigInteger.valueOf(1), BigInteger.valueOf(100000)).send();

            List<AtomicSwap.InitiatedEventResponse> emitEventData = atomicSwap.getInitiatedEvents(tx);

            System.out.println(emitEventData.get(0)._value);

            // 销毁合约
            // atomicSwap.destruct().send();
        } catch (Exception e) {
            e.printStackTrace();
        }

    }

    public static byte[] stringToBytes32(String string) {
        byte[] byteValue = string.getBytes();
        byte[] byteValueLen32 = new byte[32];
        System.arraycopy(byteValue, 0, byteValueLen32, 0, byteValue.length);
        return byteValueLen32;
    }
}
