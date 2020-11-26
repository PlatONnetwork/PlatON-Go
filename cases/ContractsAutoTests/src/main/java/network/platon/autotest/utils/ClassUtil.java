package network.platon.autotest.utils;

/**   
 * @Title: Classutil.java
 * @Package network.platon.autotest.utils
 * @Description: TODO(用一句话描述该文件做什么)
 * @author qcxiao   
 * @date 2013-9-18 下午02:20:50  
 */

import java.io.ByteArrayOutputStream;
import java.io.File;
import java.io.FileFilter;
import java.io.FileInputStream;
import java.io.IOException;
import java.net.JarURLConnection;
import java.net.URL;
import java.net.URLDecoder;
import java.util.*;
import java.util.jar.JarEntry;
import java.util.jar.JarFile;
import org.apache.commons.logging.Log;
import org.apache.commons.logging.LogFactory;

public class ClassUtil {
	private static final Log log = LogFactory.getLog(ClassUtil.class);
	@SuppressWarnings("unused")
	private static String classPath = "";

	public static void main(String[] args) throws Exception {
		getClasses("D:\\workspaces\\contracts_workspaces\\PlatON-Go\\cases\\ContractsAutoTests\\target\\test-classes\\");
	}

	public static Set<Class<?>> getClasses(String path) throws IOException {
		classPath = new StringBuffer().append(path).toString();
		// 第一个class类的集合
		Set<Class<?>> classes = new LinkedHashSet<Class<?>>();
		findAndAddClassesInPath(path, classes);
		return classes;
	}

	public static void findAndAddClassesInPath(String Path, Set<Class<?>> classes) throws IOException {
		// 获取此包的目录 建立一个File
		File dir = new File(Path);
		// 如果不存在或者 也不是目录就直接返回
		if (!dir.exists() || !dir.isDirectory()) {
			return;
		}
		// 如果存在 就获取包下的所有文件 包括目录
		File[] dirfiles = dir.listFiles(new FileFilter() {
			// 自定义过滤规则 如果可以循环(包含子目录) 或则是以.class结尾的文件(编译好的java类文件)
			public boolean accept(File file) {
				return (file.isDirectory()) || (file.getName().endsWith(".class"));
			}
		});
		// 循环所有文件
		for (File file : dirfiles) {
			// 如果是目录 则继续扫描
			if (file.isDirectory()) {
				findAndAddClassesInPath(file.getAbsolutePath(), classes);
			} else {
				// 如果是java类文件 去掉后面的.class 只留下类名
				// String className = file.getAbsolutePath().replace(classPath,
				// "").replace("\\", ".").replace(".class", "");

				FileInputStream classIs = new FileInputStream(file);
				ByteArrayOutputStream baos = new ByteArrayOutputStream();
				byte buf[] = new byte[4];
				// 读取文件流
				for (int i = 0; (i = classIs.read(buf)) != -1;) {
					baos.write(buf, 0, i);
				}
				// 创建新的类对象
				byte[] data = baos.toByteArray();
				String className = classNameAnalyzer(data);
				// String className = file.getName().substring(0,
				// file.getName().length() - 6);
				try {
					// 添加到集合中去
					// classes.add(Class.forName(packageName + '.' +
					// className));
					classes.add(Class.forName(className));
				} catch (ClassNotFoundException e) {
					log.error("添加用户自定义视图类错误 找不到此类的.class文件");
					e.printStackTrace();
				}
			}
		}
	}

	public static String classNameAnalyzer(byte[] data) {
		// 常量索引
		HashMap<Short, Short> constIndex = new HashMap<Short, Short>();

		// UTF-8 -8 string 索引
		HashMap<Short, String> stringIndex = new HashMap<Short, String>();

		// 常量池解析
		short dataIndex = 10;
		byte constType = data[dataIndex];
		for (short index = 1; index < getShort(new byte[] { data[8], data[9] }, false); index++)
			switch (constType) {
			case 1: // UTF-8 string 常量
				short d = getShort(new byte[] { data[dataIndex + 1], data[dataIndex + 2] }, false);
				stringIndex.put(index, new String(data, dataIndex + 3, d));
				dataIndex += d + 3;
				constType = data[dataIndex];
				break;
			case 3 : // integer 常量
				dataIndex += 5;
				constType = data[dataIndex];
				break;
			case 4: // flat 常量
				dataIndex += 5;
				constType = data[dataIndex];
				break;
			case 5: // long 常量
				index++;
				dataIndex += 9;
				constType = data[dataIndex];
				break;
			case 6: // double 常量
				index++;
				dataIndex += 9;
				constType = data[dataIndex];
				break;
			case 7: // class or interface reference
				constIndex.put(index, getShort(new byte[] { data[dataIndex + 1], data[dataIndex + 2] }, false));
				dataIndex += 3;
				constType = data[dataIndex];
				break;
			case 8: // string 常量
				constIndex.put(index, getShort(new byte[] { data[dataIndex + 1], data[dataIndex + 2] }, false));
				dataIndex += 3;
				constType = data[dataIndex];
				break;
			case 9: // field reference
				dataIndex += 5;
				constType = data[dataIndex];
				break;
			case 10: // method reference
				dataIndex += 5;
				constType = data[dataIndex];
				break;
			case 11: // interface method reference
				dataIndex += 5;
				constType = data[dataIndex];
				break;
			case 12: // name and type reference
				dataIndex += 5;
				constType = data[dataIndex];
				break;
			case 15: // MethodHandle
				dataIndex += 4;
				constType = data[dataIndex];
				break;
			case 16: // MethodType
				dataIndex += 4;
				constType = data[dataIndex];
				break;
			case 18: // InvokeDynamic
				dataIndex += 5;
				constType = data[dataIndex];
				break;

			default:
				throw new RuntimeException("Invalid constant pool flag: " + constType);
			}

		// 获取当前class的全限定名索引
		Short indexOfThisClass = getShort(new byte[] { data[dataIndex + 2], data[dataIndex + 3] }, false);

		if (!constIndex.containsKey(indexOfThisClass)) {
			throw new RuntimeException("class文件解析错误,获取当前类全限定名index错误");
		}

		// 获取当前class的全限定名String索引
		short resultIndex = constIndex.get(indexOfThisClass);
		if (!stringIndex.containsKey(resultIndex)) {
			throw new RuntimeException("class文件解析错误，,获取当前类全限定名Stringindex错误");
		}

		return stringIndex.get(resultIndex).replace("/", ".");
	}

	public static short getShort(byte[] buf, boolean asc)

	{
		if (buf == null) {
			throw new IllegalArgumentException("byte array is null!");
		}
		if (buf.length > 2) {
			throw new IllegalArgumentException("byte array size > 2 !");
		}
		short r = 0;
		if (asc)
			for (int i = buf.length - 1; i >= 0; i--) {
				r <<= 8;
				r |= (buf[i] & 0x00ff);
			}
		else
			for (int i = 0; i < buf.length; i++) {
				r <<= 8;
				r |= (buf[i] & 0x00ff);
			}
		return r;

	}

	/*
	 * 取得某一类所在包的所有类名 不含迭代
	 */
	public static String[] getPackageAllClassName(String classLocation, String packageName) {
		// 将packageName分解
		String[] packagePathSplit = packageName.split("[.]");
		String realClassLocation = classLocation;
		int packageLength = packagePathSplit.length;
		for (int i = 0; i < packageLength; i++) {
			realClassLocation = realClassLocation + File.separator + packagePathSplit[i];
		}
		File packeageDir = new File(realClassLocation);
		if (packeageDir.isDirectory()) {
			String[] allClassName = packeageDir.list();
			return allClassName;
		}
		return null;
	}

	/**
	 * 从包package中获取所有的Class
	 * 
	 * @param pack
	 * @return
	 */
	public static Set<Class<?>> getClasses(Package pack) {

		// 第一个class类的集合
		Set<Class<?>> classes = new LinkedHashSet<Class<?>>();
		// 是否循环迭代
		boolean recursive = true;
		// 获取包的名字 并进行替换
		String packageName = pack.getName();
		String packageDirName = packageName.replace('.', '/');
		// 定义一个枚举的集合 并进行循环来处理这个目录下的things
		Enumeration<URL> dirs;
		try {
			dirs = Thread.currentThread().getContextClassLoader().getResources(packageDirName);
			// 循环迭代下去
			// while (dirs.hasMoreElements()) {
			// 获取下一个元素
			URL url = dirs.nextElement();// 暂时只取第一个
			// 得到协议的名称
			String protocol = url.getProtocol();
			// 如果是以文件的形式保存在服务器上
			if ("file".equals(protocol)) {
				// 获取包的物理路径
				String filePath = URLDecoder.decode(url.getFile(), "UTF-8");
				// 以文件的方式扫描整个包下的文件 并添加到集合中
				findAndAddClassesInPackageByFile(packageName, filePath, recursive, classes);
			} else if ("jar".equals(protocol)) {
				// 如果是jar包文件
				// 定义一个JarFile
				JarFile jar;
				try {
					// 获取jar
					jar = ((JarURLConnection) url.openConnection()).getJarFile();
					// 从此jar包 得到一个枚举类
					Enumeration<JarEntry> entries = jar.entries();
					// 同样的进行循环迭代
					while (entries.hasMoreElements()) {
						// 获取jar里的一个实体 可以是目录 和一些jar包里的其他文件 如META-INF等文件
						JarEntry entry = entries.nextElement();
						String name = entry.getName();
						// 如果是以/开头的
						if (name.charAt(0) == '/') {
							// 获取后面的字符串
							name = name.substring(1);
						}
						// 如果前半部分和定义的包名相同
						if (name.startsWith(packageDirName)) {
							int idx = name.lastIndexOf('/');
							// 如果以"/"结尾 是一个包
							if (idx != -1) {
								// 获取包名 把"/"替换成"."
								packageName = name.substring(0, idx).replace('/', '.');
							}
							// 如果可以迭代下去 并且是一个包
							if ((idx != -1) || recursive) {
								// 如果是一个.class文件 而且不是目录
								if (name.endsWith(".class") && !entry.isDirectory()) {
									// 去掉后面的".class" 获取真正的类名
									String className = name.substring(packageName.length() + 1, name.length() - 6);
									try {
										// 添加到classes
										classes.add(Class.forName(packageName + '.' + className));
									} catch (ClassNotFoundException e) {
										log.error("添加用户自定义视图类错误 找不到此类的.class文件");
										e.printStackTrace();
									}
								}
							}
						}
					}
				} catch (IOException e) {
					log.error("在扫描用户定义视图时从jar包获取文件出错");
					e.printStackTrace();
				}
			}
			// }
		} catch (IOException e) {
			e.printStackTrace();
		}

		return classes;
	}

	/**
	 * 以文件的形式来获取包下的所有Class
	 * 
	 * @param packageName
	 * @param packagePath
	 * @param recursive
	 * @param classes
	 */
	public static void findAndAddClassesInPackageByFile(String packageName, String packagePath, final boolean recursive, Set<Class<?>> classes) {
		// 获取此包的目录 建立一个File
		File dir = new File(packagePath);
		// 如果不存在或者 也不是目录就直接返回
		if (!dir.exists() || !dir.isDirectory()) {
			log.warn("用户定义包名 " + packageName + " 下没有任何文件");
			return;
		}
		// 如果存在 就获取包下的所有文件 包括目录
		File[] dirfiles = dir.listFiles(new FileFilter() {
			// 自定义过滤规则 如果可以循环(包含子目录) 或则是以.class结尾的文件(编译好的java类文件)
			public boolean accept(File file) {
				return (recursive && file.isDirectory()) || (file.getName().endsWith(".class"));
			}
		});
		// 循环所有文件
		for (File file : dirfiles) {
			// 如果是目录 则继续扫描
			if (file.isDirectory()) {
				findAndAddClassesInPackageByFile(packageName + "." + file.getName(), file.getAbsolutePath(), recursive, classes);
			} else {
				// 如果是java类文件 去掉后面的.class 只留下类名
				String className = file.getName().substring(0, file.getName().length() - 6);
				try {
					// 添加到集合中去
					classes.add(Class.forName(packageName + '.' + className));
				} catch (ClassNotFoundException e) {
					log.error("添加用户自定义视图类错误 找不到此类的.class文件");
					e.printStackTrace();
				}
			}
		}
	}
}
