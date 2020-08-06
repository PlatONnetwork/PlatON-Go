package evm.versioncompatible.v0_5_0.v10_deprecatedFunction;

import evm.beforetest.ContractPrepareTest;
import network.platon.autotest.junit.annotations.DataSource;
import network.platon.autotest.junit.enums.DataSourceType;
import network.platon.contracts.ChainFunction;
import org.junit.Before;
import org.junit.Test;
import org.web3j.protocol.core.methods.response.TransactionReceipt;

import java.math.BigInteger;

/*
 * @title  0.5.0版本函数 revert()， require()，assert() 代替 0.4.25版本函数 throw
 * @description:
 * @author: hudenian
 * @create: 2019/12/27
 */
public class ChainFunctionTest extends ContractPrepareTest {

    private boolean isDeceased =true;

    private String less9 = "10";

    @Before
    public void before() {
        this.prepare();
    }


    @Test
    @DataSource(type = DataSourceType.EXCEL, file = "test.xls", sheetName = "Sheet1",
            author = "hudenian", showName = "ChainFunctionTest-revert，require，assert代替throw", sourcePrefix = "evm")
    public void throwExceptionTest() {
        try {

            ChainFunction chainFunction = ChainFunction.deploy(web3j, transactionManager, provider,new BigInteger("1"), chainId).send();

            String contractAddress = chainFunction.getContractAddress();
            TransactionReceipt tx = chainFunction.getTransactionReceipt().get();

            collector.logStepPass("ChainFunction deploy successfully.contractAddress:" + contractAddress + ", hash:" + tx.getTransactionHash());
            collector.logStepPass("deploy gas used:" + chainFunction.getTransactionReceipt().get().getGasUsed());

            String  msgSender = chainFunction.deceased(isDeceased,new BigInteger(less9)).send();

            //合约取枚举第二个值，翻译过来就是1
//            collector.assertEqual("1",gender.toString());

            String  msgSender2 = chainFunction.deceasedWithModify(isDeceased).send();

        } catch (Exception e) {
            collector.logStepPass("params is Inappropriate please input isDeceased =true and less9>= 9");
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
