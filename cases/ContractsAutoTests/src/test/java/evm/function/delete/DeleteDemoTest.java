package evm.function.delete;

import com.platon.sdk.utlis.Bech32;
import com.platon.sdk.utlis.NetworkParameters;
import evm.beforetest.ContractPrepareTest;
import network.platon.autotest.junit.annotations.DataSource;
import network.platon.autotest.junit.enums.DataSourceType;
import network.platon.contracts.DeleteDemo;
import network.platon.utils.DataChangeUtil;
import org.junit.Before;
import org.junit.Test;
import org.web3j.protocol.core.methods.response.TransactionReceipt;
import org.web3j.tuples.generated.Tuple2;

import java.math.BigInteger;

/**
 * 验证delete关键字，验证各种类型的delete,包括bool,uint,address,bytes,string,enum
 * @author liweic
 * @dev 2020/01/09 16:10
 */

public class DeleteDemoTest extends ContractPrepareTest {

    private String i;
    private String addr;

    @Before
    public void before() {
        this.prepare();
        i = driverService.param.get("i");
        addr = driverService.param.get("addr");

    }

    @Test
    @DataSource(type = DataSourceType.EXCEL, file = "test.xls", sheetName = "Sheet1",
            author = "liweic", showName = "function.DeleteDemoTest-delete操作测试", sourcePrefix = "evm")
    public void Deletedemo() {
        try {
            DeleteDemo deletedemo = DeleteDemo.deploy(web3j, transactionManager, provider, chainId).send();

            String contractAddress = deletedemo.getContractAddress();
            TransactionReceipt tx = deletedemo.getTransactionReceipt().get();
            collector.logStepPass("DeleteDemo deploy successfully.contractAddress:" + contractAddress + ", hash:" + tx.getTransactionHash());
            collector.logStepPass("DeleteDemo deploy gasUsed:" + deletedemo.getTransactionReceipt().get().getGasUsed());

            TransactionReceipt result = deletedemo.deleteAttr().send();
            collector.logStepPass("打印result交易Hash：" + result.getTransactionHash());

            //验证delete bool
            boolean delb = deletedemo.getbool().send();
            collector.logStepPass("delete bool返回值：" + delb);
            collector.assertEqual(false ,delb);

            //验证delete uint
            BigInteger deli = deletedemo.getunit().send();
            collector.logStepPass("delete uint返回值：" + deli);
            collector.assertEqual(i ,deli.toString());

            //验证delete addr
            String deladdr = deletedemo.getaddress().send();
            collector.logStepPass("delete addr返回值：" + deladdr);
            addr = Bech32.addressEncode(NetworkParameters.TestNetParams.getHrp(), addr);
            collector.assertEqual(addr ,deladdr);

            //delete bytes
            byte[] delbytes = deletedemo.getbytes().send();
            String a = DataChangeUtil.bytesToHex(delbytes);
            collector.logStepPass("delete bytes返回值：" + a);
            collector.assertEqual("" ,a);

            //delete str
            String delstr = deletedemo.getstr().send();
            collector.logStepPass("delete str返回值：" + delstr);
            collector.assertEqual("" ,delstr);

            //delete enum
            BigInteger delenum = deletedemo.getenum().send();
            collector.logStepPass("delete enum返回值：" + delenum);
            collector.assertEqual("0" ,delenum.toString());

            //delete dynamicarray
            BigInteger deldynamic = deletedemo.delDynamicArray().send();
            collector.logStepPass("delete dynamicarray返回值：" + deldynamic);
            collector.assertEqual("0" ,delenum.toString());

            //delete struct
            TransactionReceipt st = deletedemo.delStruct(new BigInteger("1")).send();

            Tuple2 delstruct = deletedemo.getstruct().send();
            collector.logStepPass("struct执行delete后第一个值：" + delstruct.getValue1());
            collector.assertEqual("0" ,delstruct.getValue1().toString());
            collector.logStepPass("struct执行delete后第二个值：" + delstruct.getValue2());
            collector.assertEqual("" ,delstruct.getValue2().toString());

            //delete struct mapping
            TransactionReceipt sta = deletedemo.delMapping(new BigInteger("1")).send();

            Tuple2 delstructmapping = deletedemo.getdelMapping().send();
            collector.logStepPass("struct执行delete后mapping值：" + delstructmapping.getValue1());
            collector.assertEqual("2000" ,delstructmapping.getValue1().toString());

            collector.logStepPass("struct执行delete后非mapping值：" + delstructmapping.getValue2());
            collector.assertEqual("0" ,delstructmapping.getValue2().toString());


        } catch (Exception e) {
            collector.logStepFail("DeleteDemoContract Calling Method fail.", e.toString());
            e.printStackTrace();
        }
    }
}