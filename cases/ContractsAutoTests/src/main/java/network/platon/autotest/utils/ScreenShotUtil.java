package network.platon.autotest.utils;

import java.awt.Dimension;
import java.awt.Rectangle;
import java.awt.Robot;
import java.awt.Toolkit;
import java.awt.image.BufferedImage;
import java.io.File;
import java.util.Date;
import java.util.Properties;

import javax.imageio.ImageIO;

import org.apache.commons.io.FileUtils;
import org.openqa.selenium.OutputType;
import org.openqa.selenium.TakesScreenshot;
import org.openqa.selenium.WebDriver;

public class ScreenShotUtil {
	private static String BASE_DIR = "C:\\autosky_log\\";

	public static String screenShotByDriver(WebDriver driver, String suitePath) {
		String picPath = "";
		//
		// Properties properties = FileUtil.getProperties();
		// BASE_DIR = properties.getProperty("logDir", BASE_DIR) == null ?
		// BASE_DIR : properties.getProperty("logPath", BASE_DIR);
		// if (!BASE_DIR.endsWith("/")) {
		// BASE_DIR = BASE_DIR + "/";
		// }
		// String screenShotPath = BASE_DIR;
		// // Boolean reportMerged =
		// // Boolean.valueOf(System.getProperty("report.merged") == null ?
		// "false"
		// // : System.getProperty("report.merged"));
		// // screenShotPath = screenShotPath + suitePath + "/screenshot";
		// screenShotPath = suitePath + "screenshot";
		// // 这里定义了截图存放目录名
		// if (!(new File(screenShotPath).isDirectory())) { // 判断是否存在该目录
		// new File(screenShotPath).mkdir(); // 如果不存在则新建一个目录
		// }
		String screenShotPath = createScreenShotPath(suitePath);
		String time = DateUtil.dateToStr(new Date(), "yyyyMMdd-HHmmss");
		try {
			File source_file = ((TakesScreenshot) driver).getScreenshotAs(OutputType.FILE); // 关键代码，执行屏幕截图，默认会把截图保存到temp目录
			File file = new File(screenShotPath + File.separator + time + ".png");
			FileUtils.copyFile(source_file, file); // 这里将截图另存到我们需要保存的目录，例如screenshot\20120406-165210.png
			picPath = file.getAbsolutePath();
			return picPath;
		} catch (Exception e) {
			return screenShotByDesktop(suitePath);
		}
	}

	// 如果通过webdriver截图失败的话，就采用桌面截图，后续添加上
	public static String screenShotByDesktop(String suitePath) {
		String picPath = "";
		try {
	    	Dimension d = Toolkit.getDefaultToolkit().getScreenSize();
			String time = DateUtil.dateToStr(new Date(), "yyyyMMdd-HHmmss");
			String screenShotPath = createScreenShotPath(suitePath);
			BufferedImage screen = (new Robot()).createScreenCapture(new Rectangle(0, 0, (int) d.getWidth(), (int) d.getHeight()));
			String name = screenShotPath + "/" + time + ".png";
			File file = new File(name);
			ImageIO.write(screen, "png", file);
			picPath = file.getAbsolutePath();
		} catch (Exception e) {
			// TODO: handle exception
			System.out.println("截图失败！！！\n" + e.getMessage());
		}
		return picPath;
	}

	private static String createScreenShotPath(String suitePath) {
		Properties properties = FileUtil.getProperties();
		BASE_DIR = properties.getProperty("logDir", BASE_DIR) == null ? BASE_DIR : properties.getProperty("logPath", BASE_DIR);
		if (!BASE_DIR.endsWith("/")) {
			BASE_DIR = BASE_DIR + "/";
		}
		String screenShotPath = BASE_DIR;
		screenShotPath = suitePath + "Link/screenshot";
		// 这里定义了截图存放目录名
		if (!(new File(screenShotPath).isDirectory())) { // 判断是否存在该目录
			new File(screenShotPath).mkdir(); // 如果不存在则新建一个目录
		}
		return screenShotPath;
	}
	// 截图前先滚动下拉框,解决懒加载
	// browser.execute_script("""
	// (function () {
	// var y = 0;
	// var step = 100;
	// window.scroll(0, 0);
	//
	// function f() {
	// if (y < document.body.scrollHeight) {
	// y += step;
	// window.scroll(0, y);
	// setTimeout(f, 50);
	// } else {
	// window.scroll(0, 0);
	// document.title += "scroll-done";
	// }
	// }
	//
	// setTimeout(f, 1000);
	// })();
	// """)
	//
	// for i in xrange(30):
	// if "scroll-done" in browser.title:
	// break
	// time.sleep(1)
	//
	// browser.save_screenshot(save_fn)
	// browser.close()

}
