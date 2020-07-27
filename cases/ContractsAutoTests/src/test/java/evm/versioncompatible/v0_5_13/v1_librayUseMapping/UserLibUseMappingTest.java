package evm.versioncompatible.v0_5_13.v1_librayUseMapping;

import evm.beforetest.ContractPrepareTest;
import network.platon.autotest.junit.annotations.DataSource;
import network.platon.autotest.junit.enums.DataSourceType;
import network.platon.contracts.UserMapping;
import org.junit.Before;
import org.junit.Test;
import org.web3j.protocol.core.methods.response.TransactionReceipt;

import java.math.BigInteger;


/**
 * @title   允许public和external库函数的参数和返回变量使用mapping类型
 * @description: 此测试gas超限
 * @author: hudenian
 * @create: 2019/12/27
 */
public class UserLibUseMappingTest extends ContractPrepareTest {

    //mapping中id值
    private String id="12";

    //mapping中age值
    private String age="12";

    @Before
    public void before() {
        this.prepare();
    }


    @Test
    @DataSource(type = DataSourceType.EXCEL, file = "test.xls", sheetName = "Sheet1",
            author = "hudenian", showName = "evm.UserLibUseMappingTest-允许public和external库函数的参数和返回变量使用mapping类型", sourcePrefix = "evm")
    public void testLibMappingParam() {
        try {

            UserMapping userMapping = UserMapping.deploy(web3j, transactionManager, provider, chainId).send();

            String contractAddress = userMapping.getContractAddress();
            TransactionReceipt tx = userMapping.getTransactionReceipt().get();

            collector.logStepPass("UserLibUseMappingTest deploy successfully.contractAddress:" + contractAddress + ", hash:" + tx.getTransactionHash());
            collector.logStepPass("deploy gas used:" + userMapping.getTransactionReceipt().get().getGasUsed());

            TransactionReceipt transactionReceipt = userMapping.setOutUser(new BigInteger(age),new BigInteger(id)).send();

            collector.logStepPass("SuicideAndSelfdestructTest increment successful.transactionHash:" + transactionReceipt.getTransactionHash());
            collector.logStepPass( "currentBlockNumber:" + transactionReceipt.getBlockNumber());

            String chainAge = userMapping.getOutUser(new BigInteger(id)).send().toString();

            collector.logStepPass("获取到的age值为："+chainAge);
            //链上的年龄写死
            collector.assertEqual("23",chainAge);


        } catch (Exception e) {
            collector.logStepFail("UserLibUseMappingTest testLibMappingParam process fail.", e.toString());
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
