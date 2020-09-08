package network.platon.test.evm.beforetest;

import com.platon.sdk.utlis.Bech32;
import network.platon.autotest.junit.annotations.DataSource;
import network.platon.autotest.junit.enums.DataSourceType;
import network.platon.autotest.utils.FileUtil;
import network.platon.utils.CompileUtil;
import network.platon.utils.DataChangeUtil;
import network.platon.utils.GeneratorUtil;
import network.platon.utils.OneselfFileUtil;
import org.junit.Before;
import org.junit.Test;
import org.web3j.protocol.core.methods.response.TransactionReceipt;

import java.io.*;
import java.nio.file.Paths;
import java.util.*;
import java.util.concurrent.CountDownLatch;
import java.util.concurrent.ExecutorService;
import java.util.concurrent.Executors;
import java.util.concurrent.Semaphore;

/**
 * @title 1.将sol编译成二进制和ABI文件，2.通过合约二进制和ABI文件生成包裝类
 * @author: qcxiao
 * @create: 2019/12/18 11:27
 **/
public class GeneratorPreTest extends ContractPrepareTest {

    private String contractAndLibrarys;

    @Before
    public void before() {
        this.prepare();
        contractAndLibrarys = driverService.param.get("contractAndLibrarys") == null ? "" : driverService.param.get("contractAndLibrarys").toString();
    }


    @Test
    @DataSource(type = DataSourceType.EXCEL, file = "test.xls", author = "qcxiao"
            , showName = "GeneratorPreTest-编译并生成包装类", sourcePrefix = "evm")
    public void compileAndGenerator() {
        Date compileStartDate = new Date();
        try {
            // 1.将sol编译成二进制和ABI文件
            compile();
            Date compileEndDate = new Date();
            long ms = compileEndDate.getTime() - compileStartDate.getTime();
            collector.logStepPass("compile time:" + ms + "ms");

            //2.将含有library库的合约中的引用替换为library库对对合约地址
            String[] contractAndLibrarysArr = contractAndLibrarys.split(";");
            if (contractAndLibrarysArr.length > 0) {
                for (int i = 0; i < contractAndLibrarysArr.length; i++) {
                    System.out.println("contractAndLibrarysArr[i] is:" + contractAndLibrarysArr[i]);
                    String[] singleContractLib = contractAndLibrarysArr[i].split("&");
                    deployLibrary(singleContractLib[0], singleContractLib[1]);
                }
            }

            Date generatorWrapperStartDate = new Date();
            // 3.通过合约二进制和ABI文件生成包裝类
            generatorEVMWrapper();
            Date generatorWrapperEndDate = new Date();

            ms = generatorWrapperEndDate.getTime() - generatorWrapperStartDate.getTime();
            collector.logStepPass("generator wrapper time:" + ms + "ms");
        } catch (InterruptedException e) {
            e.printStackTrace();
        }
    }


    /**
     * @title 将sol编译成二进制和ABI文件
     * @description:
     * @author: qcxiao
     * @create: 2019/12/24 14:44
     **/
    public void compile() throws InterruptedException {
        String resourcePath = FileUtil.pathOptimization(Paths.get("src", "test", "resources", "contracts", "evm").toUri().getPath());
        String buildPath = FileUtil.pathOptimization(Paths.get("src", "test", "resources", "contracts", "evm", "build").toUri().getPath());

        File buildPathFile = new File(buildPath);
        if (!buildPathFile.exists() || !buildPathFile.isDirectory()) {
            buildPathFile.mkdirs();
        }

        File[] list = new File(buildPath).listFiles();
        if (null != list) {
            for (File file : list) {
                file.delete();
            }
        }
        // 获取所有sol源文件
        List<String> files = new OneselfFileUtil().getResourcesFile(resourcePath, 0);
        int size = files.size();

        ExecutorService executorService = Executors.newCachedThreadPool();
        // 同时并发执行的线程数
        final Semaphore semaphore = new Semaphore(20);
        // 请求总数与文件数定义一致size
        CountDownLatch countDownLatch = new CountDownLatch(size);
        CompileUtil compileUtil = new CompileUtil();

        for (String file : files) {
            //collector.logStepPass("staring compile:" + file);
            executorService.execute(() -> {
                try {
                    semaphore.acquire();
                    compileUtil.evmCompile(file, buildPath);
                    collector.logStepPass("compile success:" + file);
                } catch (Exception e) {
                    collector.logStepFail("compile fail:" + file, e.toString());
                } finally {
                    semaphore.release();
                    countDownLatch.countDown();
                }

            });
        }

        countDownLatch.await();
        executorService.shutdown();
    }


    public void deployLibrary(String contractName, String librarys) {
        String libraryAddressNoPre = "";
        //key值为library库名称，value为library库对应的地址
        Map<String, String> libraryAddressNoPreMap = new HashMap<String, String>();

        //key值为library库名称，value为library库对应的引用地址
        Map<String, String> libraryReplaceMap = new HashMap<String, String>();
        String lineTxt = null;

        BufferedReader bufferedReader = null;
        try {
            String resourcePath = FileUtil.pathOptimization(Paths.get("src", "test", "resources", "contracts", "evm", "build").toUri().getPath());

            String[] libraryArr = librarys.split("\\|\\|");

            for (int i = 0; i < libraryArr.length; i++) {
                String libraryFile = resourcePath + libraryArr[i];
                File file = new File(libraryFile);
                if (!file.exists()) {
                    continue;
                }
                bufferedReader = new BufferedReader(new InputStreamReader(new FileInputStream(file), "UTF-8"));
                while ((lineTxt = bufferedReader.readLine()) != null) {
                    BaseLibrary baseLibrary = new BaseLibrary(credentials, web3j, chainId);
                    TransactionReceipt receipt = baseLibrary.deploy(BaseLibrary.GAS_PRICE, BaseLibrary.GAS_LIMIT, lineTxt);
                    collector.logStepPass("status >>>> " + receipt.getStatus());
                    if(receipt.getStatus().toString().equals("0x1")){
                        collector.logStepPass(libraryArr[i] + "部署成功");
                    }else{
                        collector.logStepPass(libraryArr[i] + "部署失败");
                    }
                    libraryAddressNoPre = receipt.getContractAddress();
                    collector.logStepPass("contract address >>>> " + libraryAddressNoPre);
                    if (libraryAddressNoPre.startsWith("lax") || libraryAddressNoPre.startsWith("lat")) {
                        libraryAddressNoPreMap.put(libraryArr[i].split("\\.")[0], DataChangeUtil.bytesToHex(Bech32.addressDecode(libraryAddressNoPre)));
                        break;
                    }
                }
            }

            //将合约中对应的library库地址进行替换
            String binData = "";
            String replaceStr = "";
            String sha3Str = "";
            String contractNameFile = resourcePath + contractName;
            File file = new File(contractNameFile);
            if (!file.exists()) {
                return;
            }
            bufferedReader = new BufferedReader(new InputStreamReader(new FileInputStream(file), "UTF-8"));
            while ((lineTxt = bufferedReader.readLine()) != null) {
                if (lineTxt.startsWith("60")) {
                    binData = lineTxt;
                    continue;
                } else if (lineTxt.startsWith("//")) {
                    String[] arr = lineTxt.substring(2).trim().split("->");
                    replaceStr = "__" + arr[0].trim() + "__"; //bin文件中需要替换的字符串
                    //replaceStr生成规则
//                    sha3Str = Hash.sha3String(arr[1].trim()).substring(2,36);
                    String[] keyArr = arr[1].split("\\:");
                    libraryReplaceMap.put(keyArr[keyArr.length - 1], replaceStr);
                }
            }
            bufferedReader.close();

            //替换合约bin中的library库引用修改为对应的真实library合约地址
            Iterator<String> it = libraryReplaceMap.keySet().iterator();
            collector.logStepPass("contract oldBinData >>> " + binData);
            while (it.hasNext()) {
                String libraryKey = it.next();
                binData = binData.replace(libraryReplaceMap.get(libraryKey), libraryAddressNoPreMap.get(libraryKey));
            }

            collector.logStepPass("contract newBinData >>> " + binData);

            //将替换好的地址的合约重新写入合约的bin文件中
            FileWriter fw = new FileWriter(contractNameFile);
            fw.write(binData);
            fw.close();

        } catch (Exception e) {
            collector.logStepFail(e.getMessage(), e.toString());
            e.printStackTrace();
        }
    }


    /**
     * @title 通过合约二进制和ABI文件生成包裝类
     * @description:
     * @author: qcxiao
     * @create: 2019/12/24 14:45
     **/
    public void generatorEVMWrapper() throws InterruptedException {
        // 获取已编译后的二进制文件
        List<String> binFileName = new OneselfFileUtil().getBinFileName();
        // 获取合约文件数量
        int size = binFileName.size();

        ExecutorService executorService = Executors.newCachedThreadPool();
        CountDownLatch countDownLatch = new CountDownLatch(size);
        // 信号量
        final Semaphore semaphore = new Semaphore(50);
        GeneratorUtil generatorUtil = new GeneratorUtil();
        collector.logStepPass("staring generator, Total " + size + " contract, please wait...");

        for (String fileName : binFileName) {
            //collector.logStepPass("staring compile:" + fileName);
            executorService.execute(() -> {
                try {
                    semaphore.acquire();
                    generatorUtil.generator(fileName);
                    collector.logStepPass("generator success:" + fileName);
                } catch (Exception e) {
                    collector.logStepFail("generator fail:" + fileName, e.toString());
                } finally {
                    semaphore.release();
                    countDownLatch.countDown();
                }
            });
        }
        countDownLatch.await();
        executorService.shutdown();
    }

}

