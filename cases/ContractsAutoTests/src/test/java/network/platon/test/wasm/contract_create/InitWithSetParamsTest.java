package network.platon.test.wasm.contract_create;

import network.platon.autotest.junit.annotations.DataSource;
import network.platon.autotest.junit.enums.DataSourceType;
import network.platon.contracts.wasm.InitWithSetParams;
import org.junit.Test;
import network.platon.test.wasm.beforetest.WASMContractPrepareTest;

import java.util.HashSet;
import java.util.Iterator;
import java.util.Set;

/**
 * @title 创建合约init入参包含set
 * @description:
 * @author: hudenian
 * @create: 2020/02/26
 */
public class InitWithSetParamsTest extends WASMContractPrepareTest {
    @Test
    @DataSource(type = DataSourceType.EXCEL, file = "test.xls", sheetName = "Sheet1",
            author = "hudenian", showName = "wasm.contract_create创建合约init入参包含set",sourcePrefix = "wasm")
    public void testMapParams() {

        Set<String> set = new HashSet<String>();
        set.add("setOne");
        set.add("setTwo");
        try {
            prepare();
            InitWithSetParams initWithSetParams = InitWithSetParams.deploy(web3j, transactionManager, provider, chainId,set).send();
            String contractAddress = initWithSetParams.getContractAddress();
            String transactionHash = initWithSetParams.getTransactionReceipt().get().getTransactionHash();
            collector.logStepPass("InitWithSetParamsTest issued successfully.contractAddress:" + contractAddress + ", hash:" + transactionHash);
            collector.logStepPass("deploy gas used:" + initWithSetParams.getTransactionReceipt().get().getGasUsed());

            //查询包含map对象
            Set chainSet = initWithSetParams.get_set().send();
            Iterator it =chainSet.iterator();
            String chainValue = "";
            int i =0;
            if(it.hasNext()){
                chainValue =it.next().toString();
                if(i == 0){
                    collector.assertEqual("setOne",chainValue);
                }else if(i==1){
                    collector.assertEqual("setTwo",chainValue);
                }
                i++;
            }

        } catch (Exception e) {
            collector.logStepFail("InitWithSetParamsTest failure,exception msg:" , e.getMessage());
            e.printStackTrace();
        }
    }
}
