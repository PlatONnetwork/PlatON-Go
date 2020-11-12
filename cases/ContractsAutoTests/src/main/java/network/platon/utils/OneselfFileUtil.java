package network.platon.utils;

import network.platon.autotest.utils.FileUtil;

import java.io.*;
import java.nio.file.Paths;
import java.util.ArrayList;
import java.util.List;

public class OneselfFileUtil {
    /*
     *收集所有evm和wasm的源文件
     */
    private List<String> evmSourceFileList = new ArrayList<>();
    private List<String> wasmSourceFileList = new ArrayList<>();

    /**
     * @title OneselfFileUtil
     * @description 获取所有sol源文件
     * @author qcxiao
     * @updateTime 2019/12/27 14:22
     */
    public List<String> getResourcesFile(String path, int deep) {
        // 获得指定文件对象
        File file = new File(path);
        // 获得该文件夹内的所有文件
        File[] files = file.listFiles();
        for (int i = 0; i < files.length; i++) {
            if (files[i].isFile()) {
                if (files[i].getName().substring(files[i].getName().lastIndexOf(".") + 1).equals("sol")) {
                    evmSourceFileList.add(files[i].getPath());
                }
            } else if (files[i].isDirectory()) {
                //文件夹需要调用递归 ，深度+1
                getResourcesFile(files[i].getPath(), deep + 1);
            }
        }
        return evmSourceFileList;
    }

    /**
     * @title OneselfFileUtil
     * @description 获取所有二进制文件，并返回文件名称列表
     * @author qcxiao
     * @updateTime 2019/12/27 14:24
     */
    public static List<String> getBinFileName() {
        List<String> files = new ArrayList<>();
        String filePath = FileUtil.pathOptimization(Paths.get("src", "test", "resources", "contracts", "evm", "build").toUri().getPath());
        File file = new File(filePath);
        File[] tempList = file.listFiles();

        for (int i = 0; i < tempList.length; i++) {
            if (tempList[i].isFile()) {
                String fileName = tempList[i].getName();
                if (fileName.substring(fileName.lastIndexOf(".") + 1).equals("bin")) {
                    fileName = fileName.substring(0, fileName.lastIndexOf("."));
                    files.add(fileName);
                }
            }
        }
        return files;
    }

    public List<String> getWasmResourcesFile(String path, int deep) {
        // 获得指定文件对象
        File file = new File(path);
        // 获得该文件夹内的所有文件
        File[] files = file.listFiles();
        for (int i = 0; i < files.length; i++) {
            if (files[i].isFile()) {
                if (files[i].getName().substring(files[i].getName().lastIndexOf(".") + 1).equals("cpp")) {
                    wasmSourceFileList.add(files[i].getPath());
                }
            } else if (files[i].isDirectory()) {
                //文件夹需要调用递归 ，深度+1
                getWasmResourcesFile(files[i].getPath(), deep + 1);
            }
        }
        return wasmSourceFileList;
    }

    public List<String> getWasmFileName() throws Exception {
        List<String> files = new ArrayList<>();
        String filePath = FileUtil.pathOptimization(Paths.get("src", "test", "resources", "contracts", "wasm", "build").toUri().getPath());
        File file = new File(filePath);

        File[] tempList = file.listFiles();
        if (null == tempList || 0 == tempList.length) {
            System.out.println("src/test/resources/contracts/wasm/build路径下无wasm和abi文件，因此请查看编译步骤.");
            throw new Exception("src/test/resources/contracts/wasm/build路径下无wasm和abi文件");
        }
        for (int i = 0; i < tempList.length; i++) {
            if (tempList[i].isFile()) {
                String fileName = tempList[i].getName();
                if (fileName.substring(fileName.lastIndexOf(".") + 1).equals("wasm")) {
                    fileName = fileName.substring(0, fileName.lastIndexOf("."));
                    files.add(fileName);
                }
            }
        }
        return files;
    }

    public static String readFile(String Path){
        BufferedReader reader = null;
        String laststr = "";
        try{
            FileInputStream fileInputStream = new FileInputStream(Path);
            InputStreamReader inputStreamReader = new InputStreamReader(fileInputStream, "UTF-8");
            reader = new BufferedReader(inputStreamReader);
            String tempString = null;
            while((tempString = reader.readLine()) != null){
                laststr += tempString;
            }
            reader.close();
        }catch(IOException e){
            e.printStackTrace();
        }finally{
            if(reader != null){
                try {
                    reader.close();
                } catch (IOException e) {
                    e.printStackTrace();
                }
            }
        }
        return laststr;
    }
}
