package evm.versioncompatible.v0_4_25;

import evm.beforetest.ContractPrepareTest;
import com.alibaba.fastjson.JSONObject;
import network.platon.autotest.junit.annotations.DataSource;
import network.platon.autotest.junit.enums.DataSourceType;
import network.platon.contracts.TypeCast;
import org.junit.Test;
import org.web3j.tuples.generated.Tuple3;

import java.math.BigInteger;
/**
 * @title 类型转换测试
 * 1. 0.4.25版本验证使用constrictor关键字定义构造函数，但是不强制声明可见性(默认为public可见性）
 * 2. 0.4.25版本同一继承层次结构中允许多次指定基类构造函数参数验证:
 * (1) 允许合约直接声明构造函数 ———— is Base(7)
 *（2）子合约构造函数继承父合约构造函数———— constructor(uint _y) Base(_y * _y)
 * 两种引用构造函数方式共存时，合约优先选择（2）方式
 * @description:
 * @author: albedo
 * @create: 2019/12/28
 */
public class TypeCastTest extends ContractPrepareTest {

    @Test
    @DataSource(type = DataSourceType.EXCEL, file = "test.xls", sheetName = "testTypeCast",
            author = "albedo", showName = "evm.v0_4_25.TypeCastTest-类型转换", sourcePrefix = "evm")
    public void testTypeCast() {
        try {
            prepare();
            TypeCast typeCast = TypeCast.deploy(web3j, transactionManager, provider, chainId).send();
            String contractAddress = typeCast.getContractAddress();
            String transactionHash = typeCast.getTransactionReceipt().get().getTransactionHash();
            collector.logStepPass("TypeCast issued successfully.contractAddress:" + contractAddress + ", hash:" + transactionHash);
            collector.logStepPass("deploy gas used:" + typeCast.getTransactionReceipt().get().getGasUsed());
            Tuple3<BigInteger, byte[], byte[]> result = typeCast.typeCast().send();
            Tuple3<BigInteger, byte[], byte[]> expect =
                    new Tuple3(new BigInteger("18"), new byte[]{0,0,4,-46}, new byte[]{0,0,18,52});
            collector.assertEqual(JSONObject.toJSONString(result), JSONObject.toJSONString(expect), "checkout type cast result");
        } catch (Exception e) {
            collector.logStepFail("TypeCastTest testTypeCast failure,exception msg:" , e.getMessage());
            e.printStackTrace();
        }
    }

}
