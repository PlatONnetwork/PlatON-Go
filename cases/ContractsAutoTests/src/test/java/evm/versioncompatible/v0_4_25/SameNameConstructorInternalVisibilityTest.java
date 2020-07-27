package evm.versioncompatible.v0_4_25;

import evm.beforetest.ContractPrepareTest;
import network.platon.autotest.junit.annotations.DataSource;
import network.platon.autotest.junit.enums.DataSourceType;
import network.platon.contracts.SameNameConstructorInternalVisibility;
import org.junit.Test;
import org.web3j.tuples.generated.Tuple3;

import java.math.BigInteger;
/**
 * @title 构造函数和可见性测试
 * 1. 0.4.25版本同名函数构造函数定义，声明internal可见性验证；
 * 2. 0.4.25版本接口(interface)函数支持external和public两种可见性，可见性声明非必须验证
 * （1）默认可见性（默认public）函数声明
 * （2）public可见性函数声明
 * （3）external可见性声明
 * 3. 0.4.25版本支持，但0.5.x已弃用变量验证
 * (1)0.4.25版本允许声明0长度的定长数组类型
 * (2)0.4.25版本允许声明0结构体成员的结构体类型
 * (3)0.4.25版本允许定义非编译期常量的 constant常量
 * (4)0.4.25版本允许使用空元组组件
 * (5)0.4.25版本允许声明未初始化的storage变量
 * (6)0.4.25版本允许使用var
 * @description:
 * @author: albedo
 * @create: 2019/12/28
 */
public class SameNameConstructorInternalVisibilityTest extends ContractPrepareTest {
    SameNameConstructorInternalVisibility visibility;


    @Test
    @DataSource(type = DataSourceType.EXCEL, file = "test.xls", sheetName = "testDiscardVariable",
            author = "albedo", showName = "evm.SameNameConstructorInternalVisibilityTest-弃用字面量及后缀整体覆盖", sourcePrefix = "evm")
    public void testDiscardVariable() {
        try {
            prepare();
            visibility=SameNameConstructorInternalVisibility.deploy(web3j,transactionManager,provider,chainId).send();
            Tuple3<BigInteger,BigInteger,BigInteger> result = visibility.discardVariable().send();
            Tuple3<BigInteger,BigInteger,BigInteger> expect=new Tuple3(new BigInteger("1"),new BigInteger("0"),new BigInteger("1"));
            collector.assertEqual(result, expect, "checkout visibility assignment result");
        } catch (Exception e) {
            collector.logStepFail("SameNameConstructorInternalVisibility testDiscardVariable failure,exception msg:" , e.getMessage());
            e.printStackTrace();
        }
    }

}
