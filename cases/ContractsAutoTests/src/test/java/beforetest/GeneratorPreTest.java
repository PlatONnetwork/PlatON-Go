package beforetest;

import network.platon.autotest.junit.annotations.DataSource;
import network.platon.autotest.junit.enums.DataSourceType;
import network.platon.autotest.junit.rules.AssertCollector;
import network.platon.autotest.junit.rules.DriverService;
import network.platon.autotest.utils.FileUtil;
import network.platon.utils.CompileUtil;
import network.platon.utils.GeneratorUtil;
import network.platon.utils.OneselfFileUtil;
import org.junit.Rule;
import org.junit.Test;
import java.nio.file.Paths;
import java.util.ArrayList;
import java.util.Date;
import java.util.List;
import java.util.concurrent.CountDownLatch;
import java.util.concurrent.ExecutorService;
import java.util.concurrent.Executors;
import java.util.concurrent.Semaphore;

/**
 * @title 1.将sol编译成二进制和ABI文件，2.通过合约二进制和ABI文件生成包裝类
 * @author: qcxiao
 * @create: 2019/12/18 11:27
 **/
public class GeneratorPreTest {
    @Rule
    public DriverService driverService = new DriverService();
    @Rule
    public AssertCollector collector = new AssertCollector();



    @Test
    @DataSource(type = DataSourceType.EXCEL, file = "test.xls", author = "qcxiao", showName = "GeneratorPreTest-编译并生成包装类")
    public void compileAndGenerator() {
        Date compileStartDate = new Date();
        try {
            // 1.将sol编译成二进制和ABI文件
            compile();
            Date compileEndDate = new Date();
            long ms = compileEndDate.getTime() - compileStartDate.getTime();
            collector.logStepPass("compile time:" + ms + "ms");

            Date generatorWrapperStartDate = new Date();
            // 2.通过合约二进制和ABI文件生成包裝类
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
        String resourcePath = FileUtil.pathOptimization(Paths.get("src", "test", "resources", "contracts").toUri().getPath());
        String buildPath = FileUtil.pathOptimization(Paths.get("src", "test", "resources", "contracts", "build").toUri().getPath());
        // 获取所有sol源文件
        List<String> files = new OneselfFileUtil().getResourcesFile(resourcePath, 0);
        int size = files.size();

        ExecutorService executorService = Executors.newCachedThreadPool();
        // 同时并发执行的线程数
        final Semaphore semaphore = new Semaphore(50);
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
                    semaphore.release();
                } catch (Exception e) {
                    collector.logStepFail("compile fail:" + file, e.toString());
                }
                countDownLatch.countDown();
            });
        }

        countDownLatch.await();
        executorService.shutdown();
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
                    semaphore.release();
                } catch (Exception e) {
                    collector.logStepFail("generator fail:" + fileName, e.toString());
                } finally {
                    countDownLatch.countDown();
                }
            });
        }
        countDownLatch.await();
        executorService.shutdown();
    }

}

