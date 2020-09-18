package network.platon.test.evm.exec_efficiency;


import network.platon.test.evm.beforetest.ContractPrepareTest;
import network.platon.autotest.junit.annotations.DataSource;
import network.platon.autotest.junit.enums.DataSourceType;
import network.platon.contracts.evm.InsertSort;
import org.junit.Before;
import org.junit.Test;
import org.web3j.protocol.core.methods.response.TransactionReceipt;

import java.math.BigInteger;
import java.util.ArrayList;
import java.util.List;


/**
 * @title 插入排序
 * @description:
 * @author: liweic
 * @create: 2020/02/26
 **/
public class InsertSortTest extends ContractPrepareTest {

    private BigInteger numberOfCalls;
    private String contractAddress;

    @Before
    public void before() {
        numberOfCalls = new BigInteger(driverService.param.get("numberOfCalls"));
    }

    @Test
    @DataSource(type = DataSourceType.EXCEL, file = "test.xls", sheetName = "Sheet1",
            author = "liweic", showName = "exec_efficiency.InsertSortTest-插入排序", sourcePrefix = "evm")
    public void test() {
        prepare();
        try {

            Integer numberOfCalls = Integer.valueOf(driverService.param.get("numberOfCalls"));

            InsertSort insertSort = InsertSort.deploy(web3j, transactionManager, provider, chainId).send();
            contractAddress = insertSort.getContractAddress();
            collector.logStepPass("contract deploy successful. contractAddress:" + contractAddress);
            collector.logStepPass("deploy gas used:" + insertSort.getTransactionReceipt().get().getGasUsed());

            List<BigInteger> array = new ArrayList<>(numberOfCalls);

            int min = -1000, max = 2000;

            for (int i = 0; i < numberOfCalls; i++) {
                BigInteger a = BigInteger.valueOf(min + (int) (Math.random() * 3001));
                array.add(a);
            }

//            Integer[] intArray = new Integer[]{131, 662, 585, 1946, 640, 1339, 1099, 1713, -588, -683, 1452, 56, 350, 1984, 1616, 382, 82, 427, -641, 564, -752, 598, 753, -912, -414, 930, -514, 1005, 793, -872, 964, 188, 1468, -266, 1513, 337, -443, 1712, -310, 1133, 723, -589, 748, 1036, -115, -916, -342, 1430, -133, -572, -137, -939, -415, 610, 937, 1117, -420, 1768, 780, 711, -82, 111, 1725, 172, 1514, -530, 1000, 1104, 1830, -523, -513, -27, 1647, 470, 1425, 835, 186, 1729, -95, -309, 1265, 1754, 452, 38, 472, -10, -238, -721, 435, 700, 904, -126, -802, 1579, 720, 412, -416, -845, 60, 213, 1979, 1154, 637, 531, -475, 1578, -141, 1026, 1032, 623, -158, -113, 41, 316, 47, 720, -427, 255, 547, 199, 888, 1865, 827, 473, 321, 209, -699, 1327, 812, 1645, 289, 288, -885, 1409, 165, 97, -52, 1657, -324, 301, 1375, 1486, -984, -839, -304, 1066, 279, 601, 1850, 217, 1158, 1527, 931, 1550, 484, -427, 946, -324, 145, 150, -845, 985, -700, 523, -452, 33, -495, -279, -47, 788, 1370, 1093, -876, 1707, 759, 538, 842, 1994, 1291, 746, 1542, 1327, 1912, 1088, 385, 1892, 769, -688, 873, -313, 346, 1750, 163, 797, 1909, -639, -869, 1661, 771, 1730, 1973, -730, 1669, 1929, 698, 754, -362, 129, 339, -119, 1299, 837, 895, 1770, -776, 777, 1220, 1311, 278, -364, -42, 761, 256, 219, -518, 401, -942, 924, -628, 1242, 781, 1996, 1534, 444, -340, -451, 1309, -667, -316, -911, 185, 807, -100, -26, -357, 315, 1017, -801, 1630, 703, 1807, 793, -218, 769, 1070, 202, 1119, -305, 1513, 1180, 287, 266, 844, 954, -296, -753, 661, 766, -992, -700, 344, 1322, -852, 1397, -873, 1226, 234, 915, 862, 66, -142, 1381, 1596, 563, 65, 1214, 167, -26, 451, 1786, -722, 1871, 2, -271, -726, 377, 1295, 1540, -166, -925, 477, 1316, 1201, -64, 154, 830, -739, 1929, 1814, 42, 887, 337, 1769, 765, 324, 1250, 1836, -734, -800, -722, 985, 1361, -205, 1189, 1112, -140, -324, 886, 1803, -840, -439, 22, 569, 257, 1436, 449, 422, 1586, 1383, -264, 950, 1796, -572, 113, -834, 822, -144, -585, -656, 1550, 932, 1900, 1872, 821, 864, 200, -220, 589, 1966, 222, -889, -249, 561, -691, 1603, -625, 41, -251, 348, 1956};
//
//            List<BigInteger> array = changeIntToBigInteger(intArray);


            BigInteger n = BigInteger.valueOf(array.size());
//            BigInteger n = new BigInteger("370");
            TransactionReceipt transactionReceipt = insertSort.OuputArrays(array, n, new BigInteger("1000")).send();

            BigInteger gasUsed = transactionReceipt.getGasUsed();
            collector.logStepPass("gasUsed:" + gasUsed);
            collector.logStepPass("contract load successful. transactionHash:" + transactionReceipt.getTransactionHash());
            collector.logStepPass("currentBlockNumber:" + transactionReceipt.getBlockNumber());

            List resultarr = insertSort.get_arr().send();
            collector.logStepPass("resultarr:" + resultarr);

        } catch (Exception e) {
            e.printStackTrace();
            collector.logStepFail("The contract fail.", e.toString());
        }
    }

    /**
     * int 数组转 BigInteger
     * @param intArray
     * @return
     */
    private List<BigInteger> changeIntToBigInteger(Integer[] intArray) {
        List<BigInteger> bigIntegerList = new ArrayList<>();
        for (int i = 0; i <intArray.length ; i++) {
            bigIntegerList.add(new BigInteger(String.valueOf(intArray[i])));
        }
        return bigIntegerList;
    }

}

