package network.platon.autotest.utils;

import java.io.IOException;
import java.net.InetAddress;
import java.net.Socket;
import java.net.UnknownHostException;
import java.util.HashMap;
import java.util.Map;

import network.platon.autotest.junit.enums.BrowserProcessType;
import network.platon.autotest.junit.enums.BrowserType;

/** 
 * SystemUtil 系统公用方法类
 * @Description: TODO(这里用一句话描述这个类的作用) 
 * @author phuang_ckg
 * @date 2015年11月18日 下午6:53:11 
 *  
 */
public class SystemUtil {
	public static Map<BrowserType, BrowserProcessType> browserProcessMap = new HashMap<BrowserType, BrowserProcessType>();

	public static void killProcessOld(String processName) {
		Runtime rt = Runtime.getRuntime();
		try {
			rt.exec("cmd.exe /C start wmic process where name='" + processName + "' call terminate");
		} catch (IOException e) {
			e.printStackTrace();
		}
	}

	public static void killProcess(String processName) {
		Runtime rt = Runtime.getRuntime();
		try {
			rt.exec("cmd.exe /C taskkill /f /im " + processName);
			// tskill iexplore
		} catch (IOException e) {
			e.printStackTrace();
		}
	}

	public static void openProcess(String processName) {
		Runtime rt = Runtime.getRuntime();
		try {
			rt.exec("cmd.exe /k start " + processName);
		} catch (IOException e) {
			e.printStackTrace();
		}
	}

	public static void killBrowserProcess() {
		// killProcess()
	}
	
	@SuppressWarnings({ "unused" })
	public static boolean isPortUsing(String host,int port) {  
        boolean flag = false;  
        InetAddress theAddress = null;
		try {
			theAddress = InetAddress.getByName(host);
		} catch (UnknownHostException e1) {
			e1.printStackTrace();
		}  
        try {  
            Socket socket = new Socket(theAddress,port);  
            flag = true;  
        } catch (Exception e) {  
              
        }
        return flag;  
    }  
}
