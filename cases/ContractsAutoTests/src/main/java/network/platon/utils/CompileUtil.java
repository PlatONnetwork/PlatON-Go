package network.platon.utils;

import lombok.extern.slf4j.Slf4j;
import network.platon.autotest.utils.FileUtil;

import java.io.BufferedReader;
import java.io.File;
import java.io.InputStreamReader;
import java.nio.file.Paths;
import java.util.concurrent.TimeUnit;

/**
 * @title CompileUtil
 * @description 编译工具类
 * @author qcxiao
 * @updateTime 2019/12/27 14:39
 */
@Slf4j
public class CompileUtil {

    public void evmCompile(String filename) throws Exception {
        File file = new File(filename);
        String compilerVersion = file.getPath().replaceAll("(.*)(0\\..\\d*\\.\\d*)(.*$)", "$2");
        String buildPath = file.getPath().replaceAll("(.*)(0\\.\\d*\\.\\d*)(.*$)", "$1");
        if(System.getProperty("os.name").toLowerCase().startsWith("windows")){
            buildPath += "build\\" + compilerVersion + "\\";
        }else {
            buildPath += "build/" + compilerVersion + "/";
        }
        log.info(compilerVersion);
        log.info(buildPath);
        File buildPathFile = new File(buildPath);
        if (!buildPathFile.exists() || !buildPathFile.isDirectory()) {
            buildPathFile.mkdirs();
        }

        File[] list = new File(buildPath).listFiles();
        if (null != list) {
            for (File f : list) {
                f.delete();
            }
        }
        try {
            Solc.compile(filename, buildPath);
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
        } catch (Exception e) {
            throw new Exception();
        } finally {
            br.close();
            ps.destroy();
        }
    }
}
