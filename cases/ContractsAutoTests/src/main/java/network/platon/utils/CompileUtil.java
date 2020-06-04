package network.platon.utils;

import com.example.contract.Solc;

import java.io.BufferedReader;
import java.io.InputStreamReader;
import java.util.concurrent.Semaphore;
import java.util.concurrent.TimeUnit;

/**
 * @title CompileUtil
 * @description 编译工具类
 * @author qcxiao
 * @updateTime 2019/12/27 14:39
 */
public class CompileUtil {

    public void evmCompile(String file, String buildPath) throws Exception {
        try {
            Solc.compile(file, buildPath);
        } catch (Exception e) {
            e.printStackTrace();
            throw new Exception(e);
        }
    }

    public void wasmCompile(String file, String buildPath) throws Exception {
        try {
            String[] args = new String[]{"/bin/bash", "-c", "/usr/local/bin/platon-cpp" + " " + file + " " + "-o" + " " + buildPath};
            execGenerate(args);
        } catch (Exception e) {
            e.printStackTrace();
            throw new Exception(e);
        }
    }

    public static void execGenerate(String[] args) throws Exception {
        Process ps = null;
        StringBuffer sb = null;
        BufferedReader br = null;
        try {
            ps = Runtime.getRuntime().exec(args);
            ps.waitFor(2, TimeUnit.SECONDS);
            br = new BufferedReader(new InputStreamReader(ps.getInputStream()));
            sb = new StringBuffer();

            String line;
            while ((line = br.readLine()) != null) {
                sb.append(line).append("\n");
            }

            String result = sb.toString();
            System.out.println(result);
        } catch (Exception e) {
            throw new Exception();
        } finally {
            br.close();
            ps.destroy();
        }
    }
}
