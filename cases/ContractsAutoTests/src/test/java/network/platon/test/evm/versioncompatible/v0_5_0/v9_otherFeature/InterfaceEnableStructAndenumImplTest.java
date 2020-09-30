package network.platon.test.evm.versioncompatible.v0_5_0.v9_otherFeature;

import network.platon.test.evm.beforetest.ContractPrepareTest;
import network.platon.autotest.junit.annotations.DataSource;
import network.platon.autotest.junit.enums.DataSourceType;
import network.platon.contracts.evm.InterfaceEnableStructAndenumImpl;
import org.junit.Before;
import org.junit.Test;
import org.web3j.protocol.core.methods.response.TransactionReceipt;

import java.math.BigInteger;


/**
 * @title  09-其它
 * 10-验证0.5.0接口允许定义结构体与枚举(0.4.25版本会报错)
 * @description:
 * @author: hudenian
 * @create: 2019/12/27
 */
public class InterfaceEnableStructAndenumImplTest extends ContractPrepareTest {

    @Before
    public void before() {
        this.prepare();
    }

    @Test
    @DataSource(type = DataSourceType.EXCEL, file = "test.xls", sheetName = "Sheet1",
            author = "hudenian", showName = "InterfaceEnableStructAndenumImplTest-接口允许定义结构体与枚举", sourcePrefix = "evm")
    public void callEnum() {
        try {

            InterfaceEnableStructAndenumImpl interfaceEnableStructAndenumImpl = InterfaceEnableStructAndenumImpl.deploy(web3j, transactionManager, provider, chainId).send();

            String contractAddress = interfaceEnableStructAndenumImpl.getContractAddress();
            TransactionReceipt tx = interfaceEnableStructAndenumImpl.getTransactionReceipt().get();

            collector.logStepPass("interfaceEnableStructAndenumImpl deploy successfully.contractAddress:" + contractAddress + ", hash:" + tx.getTransactionHash());
            collector.logStepPass("deploy gas used:" + interfaceEnableStructAndenumImpl.getTransactionReceipt().get().getGasUsed());

            BigInteger gender = interfaceEnableStructAndenumImpl.getProductCondition().send();

            //合约取枚举第二个值，翻译过来就是1
            collector.assertEqual("1",gender.toString());

        } catch (Exception e) {
            collector.logStepFail("InterfaceEnableStructAndenumImplTest callEnum process fail.", e.toString());
            e.printStackTrace();
        }
    }


    /**
     * @title do...while结果值
     * @description:
     * @author: hudenian
     * @create: 2019/12/27
     */
    public static String dowhile(int x){
        int y = x+10;
        int z = x+9;
        do{
            x+=1;
            if(x>z) continue;
        }while (x<y);
        return  String.valueOf(x);
    }

    /**
     * @title for 循环后的结果值
     * @description:
     * @author: hudenian
     * @create: 2019/12/27
     */
    public static String forsum(int x){
        int forSum = 0;
        for(int i=0;i<x;i++){
            forSum = forSum +i;
        }
        return  String.valueOf(forSum);
    }

}
