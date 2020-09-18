package network.platon.test.evm.evm_assembly_contract_call;

import network.platon.test.evm.beforetest.ContractPrepareTest;
import network.platon.autotest.junit.annotations.DataSource;
import network.platon.autotest.junit.enums.DataSourceType;
import network.platon.contracts.evm.PlatonInner;
import network.platon.utils.DataChangeUtil;
import org.junit.Before;
import org.junit.Test;
import org.web3j.protocol.core.methods.response.TransactionReceipt;

import java.util.Arrays;
import java.util.HashSet;
import java.util.Set;

/**
 * 添加evm合约调用PPOS合约场景
 * @author hudenian
 * @dev 2020/02/24
 */

public class PlatonInnerTest extends ContractPrepareTest {

    //要调用的ppos合约地址（参考：PlatON内置合约及RPC接口说明）
    private String addr;

    //调用ppos合约请求参数进行rlp编码（参考：http://192.168.18.61:8080/browse/PM-613）
    private String input;

    //交易码
    private String code;

    private String caseName;

    //查询方法没有回执
    private String[] queryCodeArray = new String[]{"4100","1005","1100","1101","1102","1103","1104","1105","1200","1201","1202","2004",
            "2100","2101","2102","2103","2104","2105","2106","3001","5100"};

    private Set<String> queryCodeSet;

    @Before
    public void before() {
        this.prepare();
        addr = driverService.param.get("addr");
        input = driverService.param.get("input");
        code = driverService.param.get("code");
        caseName = driverService.param.get("caseName");
        queryCodeSet = new HashSet<String>(Arrays.asList(queryCodeArray));
    }

    @Test
    @DataSource(type = DataSourceType.EXCEL, file = "test.xls", sheetName = "Sheet1",
            author = "hudenian", showName = "evm_assembly_contract_call.PlatonInnerTest-evm合约调用PPOS合约", sourcePrefix = "evm")
    public void platonInner() {
        try {
            PlatonInner platonInner = PlatonInner.deploy(web3j, transactionManager, provider, chainId).send();

            String contractAddress = platonInner.getContractAddress();
            TransactionReceipt tx = platonInner.getTransactionReceipt().get();
            collector.logStepPass("PlatonInnerTest deploy successfully.contractAddress:" + contractAddress + ", hash:" + tx.getTransactionHash());
            collector.logStepPass("deploy gas used:" + platonInner.getTransactionReceipt().get().getGasUsed());

            tx = platonInner.assemblyCallppos(DataChangeUtil.hexToByteArray(input),addr).send();
            collector.logStepPass("PlatonInnerTest call "+addr+" and code is:"+code+" createRestrictingPlan successfully.contractAddress:" + contractAddress + ", hash:" + tx.getTransactionHash());
            collector.logStepPass("caseName:+caseName>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>"+"code:"+code);
            if( !queryCodeSet.contains(code) && (null!=tx.getLogs().get(0).getData() && "".equals(tx.getLogs().get(0).getData()))){
                collector.logStepPass("str is >>>"+DataChangeUtil.decodeSystemContractRlp(tx.getLogs().get(0).getData(), chainId));
            }


            //获取交易回执
            byte[] resultByte = platonInner.getReturnValue().send();
            collector.logStepPass("call "+addr+" and getReturnValue is:"+new String(resultByte));

        } catch (Exception e) {
            collector.logStepFail("PlatonInnerTest Calling Method fail.", e.toString());
            e.printStackTrace();
        }
    }
}


