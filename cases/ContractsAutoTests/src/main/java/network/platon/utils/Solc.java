//
// Source code recreated from a .class file by IntelliJ IDEA
// (powered by Fernflower decompiler)
//

package network.platon.utils;

import lombok.extern.slf4j.Slf4j;

import java.io.*;
import java.net.HttpURLConnection;
import java.net.URL;
import java.nio.file.Files;
import java.nio.file.LinkOption;
import java.nio.file.Path;
import java.nio.file.Paths;
import java.util.ArrayList;
import java.util.Arrays;
import java.util.Iterator;
import java.util.concurrent.TimeUnit;
import java.util.zip.ZipEntry;
import java.util.zip.ZipInputStream;

@Slf4j
public class Solc {

    public static ArrayList<String> solcList = new ArrayList(Arrays.asList("0.4.12", "0.4.13", "0.4.14", "0.4.15", "0.4.16", "0.4.17", "0.4.18", "0.4.19", "0.4.20", "0.4.21", "0.4.22", "0.4.23", "0.4.24", "0.4.25", "0.4.26", "0.5.0", "0.5.1", "0.5.2", "0.5.3", "0.5.4", "0.5.5", "0.5.6", "0.5.7", "0.5.8", "0.5.9", "0.5.10", "0.5.11", "0.5.12", "0.5.13", "0.5.14", "0.5.15", "0.5.17", "0.6.0", "0.6.12", "0.7.1"));

    public Solc() {
    }

    public static String pullSolcBin(String contractPath) throws Exception {
        String versionContent = readFile(contractPath);
        String version = parseSolcVersion(versionContent);
        log.info("Found solc compiler version: " + version);
        if (version.equals("")) {
            log.error("Contract not provider compiler version, please check source file");
            throw new Exception("Contract not provider compiler version");
        } else if (!checkSolcVersion(version)) {
            log.error("Not supported solc version, must be 0.4.12 or later");
            throw new Exception("Not supported solc version, must be 0.4.12 or later");
        } else {
            String os = System.getProperty("os.name");
            String solcUrl = "";
            String saveFileName = "";
            String absolutePath = "";
            if (!os.startsWith("Linux") && !os.startsWith("Mac OS")) {
                if (os.startsWith("Windows")) {
                    solcUrl = "https://github.com/PlatONnetwork/solidity/releases/download/platon_v" + version + "/solidity-windows-" + version + ".zip";
                    saveFileName = "solc-windows-" + version + ".zip";
                    absolutePath = ".\\solc\\";
                } else {
                    log.error("Unsupported operate system platform");
                }
            } else {
                solcUrl = "https://github.com/PlatONnetwork/solidity/releases/download/platon_v" + version + "/solc";
                saveFileName = "solc-" + version;
                absolutePath = "./solc/";
            }

            Path path = Paths.get(absolutePath + saveFileName);
            if (Files.exists(path, new LinkOption[0])) {
                return version;
            } else {
                log.info("Pulling solc binary, waiting a moment .....");
                String path1 = downRemoteFile(solcUrl, saveFileName, absolutePath);
                if (!path1.equals("")) {
                    if (os.startsWith("Windows")) {
                        String targetDir = "solc-windows-" + version;
                        unzip(absolutePath + saveFileName, absolutePath + targetDir);
                    }

                    return version;
                } else {
                    return "";
                }
            }
        }
    }

    public static void compile(String contractPath, String targetPath) throws Exception {
        String version = null;

        try {
            version = pullSolcBin(contractPath);
        } catch (Exception var13) {
            throw var13;
        }

        if (!version.equals("")) {
            new ProcessBuilder(new String[0]);
            String os = System.getProperty("os.name");
            String[] cmd = new String[0];
            if (!os.startsWith("Linux") && !os.startsWith("Mac OS")) {
                if (os.startsWith("Windows")) {
                    cmd = new String[]{"cmd", "/C", "cd .\\scripts && compile.bat " + version + " " + contractPath + " " + targetPath};
                } else {
                    log.error("Not supported operate system platform");
                }
            } else {
                cmd = new String[]{"/bin/bash", "-c", "pwd && cd ./scripts && ./compile.sh " + version + " " + contractPath + " " + targetPath};
            }

            try {
                Process ps = Runtime.getRuntime().exec(cmd);
                ps.waitFor(5L, TimeUnit.SECONDS);
                BufferedReader br = new BufferedReader(new InputStreamReader(ps.getInputStream()));
                BufferedReader bufrError = new BufferedReader(new InputStreamReader(ps.getErrorStream(), "UTF-8"));
                StringBuilder sb = new StringBuilder();
                StringBuilder err = new StringBuilder();

                String line;
                while(br.ready() && (line = br.readLine()) != null) {
                    sb.append(line).append("\n");
                }

                for(; bufrError.ready() && (line = bufrError.readLine()) != null; sb.append(line).append('\n')) {
                    if (line.contains("Error:")) {
                        err.append(line).append('\n');
                    }
                }

                String result = sb.toString();
                br.close();
                bufrError.close();
                if (ps != null) {
                    ps.destroy();
                }

                if (!err.toString().equals("")) {
                    throw new Exception(err.toString());
                }
            } catch (Exception var14) {
                var14.printStackTrace();
                throw var14;
            }
        }

    }

    public static String readFile(String filePath) {
        File file = new File(filePath);
        BufferedReader reader = null;
        String line = null;
        try {
            reader = new BufferedReader(new FileReader(file));
            while((line = reader.readLine()) != null){
                if(line.matches("(\\s)?pragma(\\s)?solidity.*")){
                    break;
                }
            }
            reader.close();
        } catch (IOException var15) {
            var15.printStackTrace();
            line = "";
        } finally {
            if (reader != null) {
                try {
                    reader.close();
                } catch (IOException var14) {
                }
            }

        }
        return line;
    }

    public static String parseSolcVersion(String content) {
        String version = content.replace("pragma solidity ", "").replace(" ", "").replace(";", "");
        if (!version.startsWith("0") && !version.startsWith("v")) {
            String maxVersion;
            if (!version.startsWith(">=") && !version.startsWith("^")) {
                if (version.startsWith("<=")) {
                    maxVersion = version.replace("<=", "");
                    return maxVersion.length() == 3 ? maxVersion + ".0" : maxVersion;
                } else {
                    Iterator iterator;
                    if (!version.startsWith(">") && !version.startsWith("!=")) {
                        if (version.startsWith("<")) {
                            maxVersion = version.replace("<", "");
                            if (maxVersion.length() == 3) {
                                maxVersion = maxVersion + ".0";
                            }

                            iterator = solcList.iterator();

                            for(int index = 0; iterator.hasNext(); ++index) {
                                String solcVersion = (String)iterator.next();
                                if (solcVersion.compareTo(maxVersion) >= 0) {
                                    if (index >= 1) {
                                        return (String)solcList.get(index - 1);
                                    }

                                    return "";
                                }
                            }

                            return "";
                        } else {
                            return "";
                        }
                    } else {
                        maxVersion = version.replace(">", "").replace("!=", "").split("<")[0];
                        if (maxVersion.length() == 3) {
                            maxVersion = maxVersion + ".0";
                        }

                        iterator = solcList.iterator();

                        String solcVersion;
                        do {
                            if (!iterator.hasNext()) {
                                return "";
                            }

                            solcVersion = (String)iterator.next();
                        } while(solcVersion.compareTo(maxVersion) <= 0);

                        return solcVersion;
                    }
                }
            } else {
                maxVersion = version.replace(">=", "").replace("^", "").split("<")[0];
                return maxVersion.length() == 3 ? maxVersion + ".0" : maxVersion;
            }
        } else {
            version = version.replace("v", "");
            if (version.length() > 3) {
                return version;
            } else {
                return version.length() == 3 ? version + ".0" : (String)solcList.get(0);
            }
        }
    }

    public static boolean checkSolcVersion(String version) {
        return solcList.contains(version);
    }

    public static String downRemoteFile(String remoteFileUrl, String saveFileName, String saveDir) {
        HttpURLConnection conn = null;
        OutputStream oputstream = null;
        InputStream iputstream = null;

        try {
            File savePath = new File(saveDir);
            if (!savePath.exists()) {
                savePath.mkdir();
            }

            File file = new File(saveDir + saveFileName);
            if (file != null && !file.exists()) {
                file.createNewFile();
            }

            URL url = new URL(remoteFileUrl);
            conn = (HttpURLConnection)url.openConnection();
            conn.setDoInput(true);
            conn.connect();
            iputstream = conn.getInputStream();
            oputstream = new FileOutputStream(file);
            byte[] buffer = new byte[4096];
            boolean var10 = true;

            int byteRead;
            while((byteRead = iputstream.read(buffer)) != -1) {
                oputstream.write(buffer, 0, byteRead);
            }

            oputstream.flush();
        } catch (Exception var19) {
            var19.printStackTrace();
        } finally {
            try {
                if (iputstream != null) {
                    iputstream.close();
                }

                if (oputstream != null) {
                    oputstream.close();
                }

                if (conn != null) {
                    conn.disconnect();
                }
            } catch (IOException var18) {
                var18.printStackTrace();
            }

        }

        return saveDir + saveFileName;
    }

    public static void unzip(String zipFilePath, String destDir) {
        File dir = new File(destDir);
        if (!dir.exists()) {
            dir.mkdirs();
        }

        byte[] buffer = new byte[1024];

        try {
            FileInputStream fis = new FileInputStream(zipFilePath);
            ZipInputStream zis = new ZipInputStream(fis);

            for(ZipEntry ze = zis.getNextEntry(); ze != null; ze = zis.getNextEntry()) {
                String fileName = ze.getName();
                File newFile = new File(destDir + File.separator + fileName);
                log.info("Unzipping to " + newFile.getAbsolutePath());
                (new File(newFile.getParent())).mkdirs();
                FileOutputStream fos = new FileOutputStream(newFile);

                int len;
                while((len = zis.read(buffer)) > 0) {
                    fos.write(buffer, 0, len);
                }

                fos.close();
                zis.closeEntry();
            }

            zis.closeEntry();
            zis.close();
            fis.close();
        } catch (IOException var11) {
            var11.printStackTrace();
        }

    }
}
