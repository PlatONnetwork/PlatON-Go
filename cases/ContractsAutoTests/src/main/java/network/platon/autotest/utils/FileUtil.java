package network.platon.autotest.utils;

import java.io.BufferedInputStream;
import java.io.BufferedReader;
import java.io.File;
import java.io.FileInputStream;
import java.io.FileNotFoundException;
import java.io.FileOutputStream;
import java.io.FileReader;
import java.io.FileWriter;
import java.io.IOException;
import java.io.InputStream;
import java.io.InputStreamReader;
import java.io.RandomAccessFile;
import java.io.Reader;
import java.util.ArrayList;
import java.util.List;
import java.util.Map;
import java.util.Properties;
import java.util.regex.Matcher;
import java.util.regex.Pattern;

import javax.imageio.stream.FileImageInputStream;

import org.apache.commons.logging.Log;
import org.apache.commons.logging.LogFactory;

public class FileUtil {
	protected static final Log logger = LogFactory.getLog(FileUtil.class);
	private static final String TEST_PROPERTIES = "/test.properties";

	//获取项目的根路径
	public final static String classPath;

	static {
		//获取的是classpath路径，适用于读取resources下资源
		classPath = Thread.currentThread().getContextClassLoader().getResource("").getPath();
	}

	/**
	 * 自定义追加路径
	 */
	public static String getCompilePath(String u_path) {
		return pathOptimization(classPath + u_path);
	}

	public static String pathOptimization(String path) {
		//windows下
		if ("\\".equals(File.separator)) {
			path = path.replaceAll("/", "\\\\");
			if (path.substring(0, 1).equals("\\") || path.substring(0, 1).equals("/")) {
				path = path.substring(1);
			}
		}
		//linux下
		if ("/".equals(File.separator)) {
			path = path.replaceAll("\\\\", "/");
		}
		return path;
	}


	public static void main(String[] args) {

		System.out.println(getCompilePath("templates/caseResult.vm"));
		String fileName = "C:/temp/newTemp.txt";
		String content = "new append!";

		readFileByBytes(fileName);
		readFileByChars(fileName);
		readFileByLines(fileName);
		readFileByRandomAccess(fileName);

		// 按方法A追加文件
		appendMethodA(fileName, content);
		appendMethodA(fileName, "append end. \n");
		// 显示文件内容
		readFileByLines(fileName);
		// 按方法B追加文件
		appendMethodB(fileName, content);
		appendMethodB(fileName, "append end. \n");
		// 显示文件内容
		readFileByLines(fileName);
	}


	/**
	 * 
	 * @param imagepath
	 * @return
	 */
	public static byte[] getImageByteArray(String imagepath) {
		byte[] image = null;
		// 直接通过文件获取
		try {
			FileImageInputStream fiis = new FileImageInputStream(new File(
					imagepath));
			image = new byte[(int) fiis.length()];
			fiis.read(image);
		} catch (FileNotFoundException e) {
			// TODO Auto-generated catch block
			e.printStackTrace();
		} catch (IOException e) {
			// TODO Auto-generated catch block
			e.printStackTrace();
		}
		// 通过数据库中获取的byte[]字节数组来获得Image对象，这里用fiis模拟数据库中的byte[].如果只是从文件中获取，那么可以跳过这几行代码，最终的目的都是为了获得BufferedImage对象
		return image;
	}

	/**
	 * 
	 * 取得调用者所在类路径中的文件的绝对路径 如getFilePath("/xxx.txt")表示取得classes根目录下xxx.txt的绝对路径
	 * 所有放在src/main/resources目录下的资源文件都会自动复制到classes目录下
	 * 
	 * @param filePath
	 * @return
	 */
	@SuppressWarnings("deprecation")
	public static String getFilePath(String filePath) {
		/*
		 try {
			return new File(sun.reflect.Reflection.getCallerClass(2).getResource(filePath).toURI()).getAbsolutePath();
		} catch (URISyntaxException ex) {
			throw new RuntimeException(ex);
		} catch (NullPointerException ex) {
			logger.warn("The file main/resources" + filePath + " is not exist!");
		}
		 */
		return null;
	}

	/**
	 * 读取配置文件(test.properties)的内容
	 * 
	 * @return
	 */
	public static Properties getProperties() {
		return getProperties(null);
	}

	/**
	 * 读取指定配置文件的内容
	 * 
	 * @return
	 */
	public static Properties getProperties(String filePath) {
		InputStream inputStream = null;
		if (filePath == null || "".equals(filePath))
			inputStream = FileUtil.class.getResourceAsStream(TEST_PROPERTIES);
		else
			inputStream = FileUtil.class.getResourceAsStream(filePath);
		Properties properties = new Properties();
		try {
			properties.load(inputStream);
		} catch (IOException e) {
			logger.error(e.getMessage());
		}
		return properties;
	}

	public static List<String> chooseFile(String downloadFileDir,
			String chooseFiles) {
		List<String> fileList = new ArrayList<String>();
		String[] downFileArry = chooseFiles.split(",");
		if (!downloadFileDir.endsWith(File.separator)) {
			downloadFileDir = downloadFileDir + File.separator;
		}
		File dirFile = new File(downloadFileDir);
		// 如果fileDir对应的文件不存在，或者不是一个目录，则退出
		if (!dirFile.exists() || !dirFile.isDirectory()) {
			logger.debug("要删除的文件夹【" + downloadFileDir + "】不存在，不需要删除操作！");
			return fileList;
		}
		// 删除文件夹下的所有文件(包括子目录)
		File[] files = dirFile.listFiles();
		for (int i = 0; i < files.length; i++) {
			boolean flag = false;
			String fileName = files[i].getName();
			for (int j = 0; j < downFileArry.length; j++) {
				String downloadFileName = downFileArry[j];
				if (fileName.equals(downloadFileName)) {
					flag = true;
					break;
				}
			}
			if (!flag) {
				fileList.add(fileName);
			}
		}
		return fileList;
	}

	/**
	 * 删除目录下的所有文件
	 * 
	 * @param fileDir
	 * @return
	 */
	public static boolean deleteDirectory(String fileDir) {
		// 如果fileDir不以文件分隔符结尾，自动添加文件分隔符
		if (!fileDir.endsWith(File.separator)) {
			fileDir = fileDir + File.separator;
		}
		File dirFile = new File(fileDir);
		// 如果fileDir对应的文件不存在，或者不是一个目录，则退出
		if (!dirFile.exists() || !dirFile.isDirectory()) {
			logger.debug("删除目录失败" + fileDir + "目录不存在！");
			return false;
		}
		boolean flag = true;
		// 删除文件夹下的所有文件(包括子目录)
		File[] files = dirFile.listFiles();
		for (int i = 0; i < files.length; i++) {
			// 删除子文件
			if (files[i].isFile()) {
				flag = deleteFile(files[i].getAbsolutePath());
				if (!flag) {
					break;
				}
			}
			// 删除子目录
			else {
				flag = deleteDirectory(files[i].getAbsolutePath());
				if (!flag) {
					break;
				}
			}
		}
		if (!flag) {
			logger.debug("删除目录失败");
			return false;
		}
		// 删除当前目录
		if (dirFile.delete()) {
			logger.debug("删除目录" + fileDir + "成功！");
			return true;
		} else {
			logger.debug("删除目录" + fileDir + "失败！");
			return false;
		}
	}

	/**
	 * 判断目录存在不存在？如果不存在，则创建之
	 * 
	 * @param directory
	 */
	public static void exist(String directory) {
		try {
			if (!new File(directory).isDirectory()) {
				new File(directory).mkdir();
			}
		} catch (SecurityException e) {
			e.printStackTrace();
		}
	}

	/**
	 * 删除文件
	 * 
	 * @param fileName
	 * @return
	 */
	private static boolean deleteFile(String fileName) {
		File file = new File(fileName);
		if (file.isFile() && file.exists()) {
			file.delete();
			return true;
		} else {
			logger.debug("删除单个文件" + fileName + "失败！");
			return false;
		}
	}

	/**
	 * 写入文件
	 * 
	 * @param in
	 * @param filePath
	 */
	public static void writeFile(InputStream in, String filePath) {
		try {
			String path = filePath.substring(0, filePath.lastIndexOf("/"));
			File file = new File(path);
			if (!file.exists()) {
				file.mkdirs();
			}
			FileOutputStream fos = null;
			BufferedInputStream bis = null;
			int BUFFER_SIZE = 1024;
			byte[] buf = new byte[BUFFER_SIZE];
			int size = 0;
			bis = new BufferedInputStream(in);
			fos = new FileOutputStream(filePath, false);
			while ((size = bis.read(buf)) != -1)
				fos.write(buf, 0, size);
			fos.close();
			bis.close();
		} catch (Exception ex) {
			ex.printStackTrace();
		}
	}

	/**
	 * 以字节为单位读取文件，常用于读二进制文件，如图片、声音、影像等文件。
	 */
	public static void readFileByBytes(String fileName) {
		File file = new File(fileName);
		InputStream in = null;
		try {
			System.out.println("以字节为单位读取文件内容，一次读一个字节：");
			// 一次读一个字节
			in = new FileInputStream(file);
			int tempbyte;
			while ((tempbyte = in.read()) != -1) {
				System.out.write(tempbyte);
			}
			in.close();
		} catch (IOException e) {
			e.printStackTrace();
			return;
		}
		try {
			System.out.println("以字节为单位读取文件内容，一次读多个字节：");
			// 一次读多个字节
			byte[] tempbytes = new byte[100];
			int byteread = 0;
			in = new FileInputStream(fileName);
			showAvailableBytes(in);
			// 读入多个字节到字节数组中，byteread为一次读入的字节数
			while ((byteread = in.read(tempbytes)) != -1) {
				System.out.write(tempbytes, 0, byteread);
			}
		} catch (Exception e1) {
			e1.printStackTrace();
		} finally {
			if (in != null) {
				try {
					in.close();
				} catch (IOException e1) {
				}
			}
		}
	}

	/**
	 * 以字符为单位读取文件，常用于读文本，数字等类型的文件
	 */
	public static void readFileByChars(String fileName) {
		File file = new File(fileName);
		Reader reader = null;
		try {
			System.out.println("以字符为单位读取文件内容，一次读一个字节：");
			// 一次读一个字符
			reader = new InputStreamReader(new FileInputStream(file));
			int tempchar;
			while ((tempchar = reader.read()) != -1) {
				// 对于windows下，\r\n这两个字符在一起时，表示一个换行。
				// 但如果这两个字符分开显示时，会换两次行。
				// 因此，屏蔽掉\r，或者屏蔽\n。否则，将会多出很多空行。
				if (((char) tempchar) != '\r') {
					System.out.print((char) tempchar);
				}
			}
			reader.close();
		} catch (Exception e) {
			e.printStackTrace();
		}
		try {
			System.out.println("以字符为单位读取文件内容，一次读多个字节：");
			// 一次读多个字符
			char[] tempchars = new char[30];
			int charread = 0;
			reader = new InputStreamReader(new FileInputStream(fileName));
			// 读入多个字符到字符数组中，charread为一次读取字符数
			while ((charread = reader.read(tempchars)) != -1) {
				// 同样屏蔽掉\r不显示
				if ((charread == tempchars.length)
						&& (tempchars[tempchars.length - 1] != '\r')) {
					System.out.print(tempchars);
				} else {
					for (int i = 0; i < charread; i++) {
						if (tempchars[i] == '\r') {
							continue;
						} else {
							System.out.print(tempchars[i]);
						}
					}
				}
			}

		} catch (Exception e1) {
			e1.printStackTrace();
		} finally {
			if (reader != null) {
				try {
					reader.close();
				} catch (IOException e1) {
				}
			}
		}
	}

	/**
	 * 以行为单位读取文件，常用于读面向行的格式化文件
	 */
	public static List<String> readFileByLines(String fileName) {
		List<String> fileLineList = new ArrayList<String>();
		File file = new File(fileName);
		BufferedReader reader = null;
		try {
			// System.out.println("以行为单位读取文件内容，一次读一整行：");
			reader = new BufferedReader(new FileReader(file));
			String tempString = null;
			// 一次读入一行，直到读入null为文件结束
			while ((tempString = reader.readLine()) != null) {
				// 显示行号
				fileLineList.add(tempString);
			}
			reader.close();
		} catch (IOException e) {
			e.printStackTrace();
		} finally {
			if (reader != null) {
				try {
					reader.close();
				} catch (IOException e1) {
				}
			}

		}
		return fileLineList;
	}

	/**
	 * 随机读取文件内容
	 */
	public static void readFileByRandomAccess(String fileName) {
		RandomAccessFile randomFile = null;
		try {
			System.out.println("随机读取一段文件内容：");
			// 打开一个随机访问文件流，按只读方式
			randomFile = new RandomAccessFile(fileName, "r");
			// 文件长度，字节数
			long fileLength = randomFile.length();
			// 读文件的起始位置
			int beginIndex = (fileLength > 4) ? 4 : 0;
			// 将读文件的开始位置移到beginIndex位置。
			randomFile.seek(beginIndex);
			byte[] bytes = new byte[10];
			int byteread = 0;
			// 一次读10个字节，如果文件内容不足10个字节，则读剩下的字节。
			// 将一次读取的字节数赋给byteread
			while ((byteread = randomFile.read(bytes)) != -1) {
				System.out.write(bytes, 0, byteread);
			}
		} catch (IOException e) {
			e.printStackTrace();
		} finally {
			if (randomFile != null) {
				try {
					randomFile.close();
				} catch (IOException e1) {
				}
			}
		}
	}

	/**
	 * 显示输入流中还剩的字节数
	 */
	private static void showAvailableBytes(InputStream in) {
		try {
			System.out.println("当前字节输入流中的字节数为:" + in.available());
		} catch (IOException e) {
			e.printStackTrace();
		}
	}

	/**
	 * A方法追加文件：使用RandomAccessFile
	 */
	public static void appendMethodA(String fileName, String content) {
		try {
			// 打开一个随机访问文件流，按读写方式
			RandomAccessFile randomFile = new RandomAccessFile(fileName, "rw");
			// 文件长度，字节数
			long fileLength = randomFile.length();
			// 将写文件指针移到文件尾。
			randomFile.seek(fileLength);
			randomFile.writeBytes(content);
			randomFile.close();
		} catch (IOException e) {
			e.printStackTrace();
		}
	}

	/**
	 * 复制整个文件夹内容
	 * 
	 * @param oldPath String 原文件路径 如：c:/fqf
	 * @param newPath String 复制后路径 如：f:/fqf/ff
	 * @return boolean
	 */
	public static void copyFolder(String oldPath, String newPath) {
		try {
			(new File(newPath)).mkdirs(); // 如果文件夹不存在 则建立新文件夹
			File a = new File(oldPath);
			String[] file = a.list();
			File temp = null;
			for (int i = 0; i < file.length; i++) {
				if (oldPath.endsWith(File.separator)) {
					temp = new File(oldPath + file[i]);
				} else {
					temp = new File(oldPath + File.separator + file[i]);
				}

				if (temp.isFile()) {
					FileInputStream input = new FileInputStream(temp);
					FileOutputStream output = new FileOutputStream(newPath
							+ "/" + (temp.getName()).toString());
					byte[] b = new byte[1024 * 5];
					int len;
					while ((len = input.read(b)) != -1) {
						output.write(b, 0, len);
					}
					output.flush();
					output.close();
					input.close();
				}
				if (temp.isDirectory()) {// 如果是子文件夹
					copyFolder(oldPath + "/" + file[i], newPath + "/" + file[i]);
				}
			}
		} catch (Exception e) {
			System.out.println("复制整个文件夹内容操作出错");
			e.printStackTrace();

		}

	}

	/**
	 * B方法追加文件：使用FileWriter
	 */
	public static void appendMethodB(String fileName, String content) {
		try {
			// 打开一个写文件器，构造函数中的第二个参数true表示以追加形式写文件
			FileWriter writer = new FileWriter(fileName, true);
			writer.write(content);
			writer.close();
		} catch (IOException e) {
			e.printStackTrace();
		}
	}



	// 读取文件
	public static String readFile(String fileName) {
		String returnStr = "";
		File file = new File(fileName);
		Reader reader = null;
		try {
			// System.out.println("以字符为单位读取文件内容，一次读一个字节：");
			// 一次读一个字符
			reader = new InputStreamReader(new FileInputStream(file));
			int tempchar;
			while ((tempchar = reader.read()) != -1) {
				// 对于windows下，\r\n这两个字符在一起时，表示一个换行。
				// 但如果这两个字符分开显示时，会换两次行。
				// 因此，屏蔽掉\r，或者屏蔽\n。否则，将会多出很多空行。
				if (((char) tempchar) != '\r') {
					returnStr += (char) tempchar;
					// System.out.print((char) tempchar);
				}
			}
			reader.close();
		} catch (Exception e) {
			e.printStackTrace();
		}
		char[] tempchars = new char[30];
		try {
			// System.out.println("以字符为单位读取文件内容，一次读多个字节：");
			// 一次读多个字符
			int charread = 0;
			reader = new InputStreamReader(new FileInputStream(fileName));
			// 读入多个字符到字符数组中，charread为一次读取字符数
			while ((charread = reader.read(tempchars)) != -1) {
				// 同样屏蔽掉\r不显示
				if ((charread == tempchars.length)
						&& (tempchars[tempchars.length - 1] != '\r')) {
					// System.out.print(tempchars);
				} else {
					for (int i = 0; i < charread; i++) {
						if (tempchars[i] == '\r') {
							continue;
						} else {
							// System.out.print(tempchars[i]);
						}
					}
				}
			}

		} catch (Exception e1) {
			e1.printStackTrace();
		} finally {
			if (reader != null) {
				try {
					reader.close();
				} catch (IOException e1) {
				}
			}
		}
		return returnStr;
	}

	public static String initHttpRequest(String fileName,
			Map<String, String> params) {
		String req = FileUtil.readFile(fileName);
		// 替换请求相关字段
		Pattern p1 = Pattern.compile("\\$\\{\\#\\#.*?\\}");
		Matcher mat1 = p1.matcher(req);
		List<String> elements1 = new ArrayList<String>();

		while (mat1.find()) {
			String element = mat1.group();
			elements1.add(element);
		}

		for (String element : elements1) {
			int i = element.indexOf("${##");
			int j = element.indexOf("}");
			String sub = element.substring(i + 4, j);
			String param = "";
			if (!params.get(sub).equals("null")) {
				param = params.get(sub);
			}
			req = req.substring(0, req.indexOf(element))
					+ param
					+ req.substring(req.indexOf(element) + element.length(),
							req.length());
		}

		// 以下处理注释，去掉/**/
		Pattern p2 = Pattern.compile("/\\*.*?\\*/");
		Matcher mat2 = p2.matcher(req);
		List<String> elements2 = new ArrayList<String>();

		while (mat2.find()) {
			String element = mat2.group();
			elements2.add(element);
		}

		for (String element : elements2) {
			req = req.substring(0, req.indexOf(element))
					+ req.substring(req.indexOf(element) + element.length(),
							req.length());
		}

		return req;

	}

	// 替换指定字段的请求值
	public static String changeHttpRequestByColumn(String req, String column,
			String value) {
		String element = "${#" + column + "}";
		return req.substring(0, req.indexOf(element))
				+ value
				+ req.substring(req.indexOf(element) + element.length(),
						req.length());
	}

}
