package network.platon.utils;

import com.example.contract.Solc;
import java.util.concurrent.Semaphore;

/**
 * @title CompileUtil
 * @description 编译工具类
 * @author qcxiao
 * @updateTime 2019/12/27 14:39
 */
public class CompileUtil {
    private final Semaphore permit = new Semaphore(100, true);

    public void evmCompile(String file, String buildPath) throws Exception {
        try {
            permit.acquire();
            Solc.compile(file, buildPath);
        } catch (Exception e) {
            e.printStackTrace();
            throw new Exception(e);
        } finally {
            permit.release();
        }
    }
}
