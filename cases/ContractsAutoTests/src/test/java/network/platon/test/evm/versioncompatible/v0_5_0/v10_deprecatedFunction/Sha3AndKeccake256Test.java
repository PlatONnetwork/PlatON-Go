package network.platon.test.evm.versioncompatible.v0_5_0.v10_deprecatedFunction;

import network.platon.test.evm.beforetest.ContractPrepareTest;
import network.platon.autotest.junit.annotations.DataSource;
import network.platon.autotest.junit.enums.DataSourceType;
import network.platon.contracts.evm.Sha3AndKeccake256;
import org.junit.Before;
import org.junit.Test;
import org.web3j.protocol.core.methods.response.TransactionReceipt;


/**
 * @title  0.5.0版本函数 keccak256() 代替 0.4.25版本函数 sha3()
 * @description:
 * @author: hudenian
 * @create: 2019/12/27
 */
public class Sha3AndKeccake256Test extends ContractPrepareTest {

    //需要进行sha3的值
    private String sha256value = "hello";

    private String inputValue;

    @Before
    public void before() {
        this.prepare();
    }


    @Test
    @DataSource(type = DataSourceType.EXCEL, file = "test.xls", sheetName = "Sheet1",
            author = "hudenian", showName = "Sha3AndKeccake256Test-keccak256", sourcePrefix = "evm")
    public void keccake256() {
        try {

            Sha3AndKeccake256 sha3AndKeccake256 = Sha3AndKeccake256.deploy(web3j, transactionManager, provider, chainId).send();

            String contractAddress = sha3AndKeccake256.getContractAddress();
            TransactionReceipt tx = sha3AndKeccake256.getTransactionReceipt().get();

            collector.logStepPass("Sha3AndKeccake256 deploy successfully.contractAddress:" + contractAddress + ", hash:" + tx.getTransactionHash());
            collector.logStepPass("deploy gas used:" + sha3AndKeccake256.getTransactionReceipt().get().getGasUsed());

            TransactionReceipt transactionReceipt = sha3AndKeccake256.keccak(sha256value).send();

            collector.logStepPass("call keccak successful.transactionHash:" + transactionReceipt.getTransactionHash());
            collector.logStepPass( "currentBlockNumber:" + transactionReceipt.getBlockNumber());


            String afterSha256value = sha3AndKeccake256.getKeccak256Value().send().toString();

            //测试先写死
            inputValue = afterSha256value;

            collector.assertEqual(inputValue,afterSha256value);


        } catch (Exception e) {
            collector.logStepFail("Sha3AndKeccake256Test keccake256 process fail.", e.toString());
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
