package network.platon.autotest.utils;

import java.io.IOException;

public class CommandUtil {
	public static void excuteCommand(String command){
		try {
			Runtime.getRuntime().exec(command);
		} catch (IOException e) {
			e.printStackTrace();
		}
	}
	public static void excuteSwipCommand(){
		try {
			Runtime.getRuntime().exec("adb -s 48621121 shell input swipe 300 200 10 200");
		} catch (IOException e) {
			e.printStackTrace();
		}
	}
}
