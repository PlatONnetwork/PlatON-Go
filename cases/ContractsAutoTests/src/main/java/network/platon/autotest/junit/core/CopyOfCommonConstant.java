package network.platon.autotest.junit.core;

import java.util.HashMap;
import java.util.Map;
import java.util.Properties;

import org.junit.runner.Description;
import org.junit.runners.model.Statement;

import network.platon.autotest.utils.FileUtil;

/**
 * @Title: CommonConstant.java
 * @Package network.platon.autotest.junit.core
 * @Description: TODO(用一句话描述该文件做什么)
 * @author qcxiao
 */
public class CopyOfCommonConstant {
	public static String SOURCES_DIR = "src/test/resources/";
	public static String TEMPLATES_DIR = "src/main/resources/templates/";
	public static Properties PROPERTIES = FileUtil.getProperties();
	public static Description DESCRIPTION = null;
	public static Statement STATEMENT = null;
	public static Boolean SUITE_MERGED = true;
	public static String PLAN_VM_CONTENT = "<?xml version=\"1.0\" encoding=\"$!encode\"?>\n<suite name=\"$!suiteInfo.suiteName\">\n#foreach ($!moduleInfo in $!suiteInfo.moduleInfoList)\n<module  run=\"$!moduleInfo.moduleRun\" name=\"$!moduleInfo.moduleName\">\n#foreach ($caseInfo in $moduleInfo.caseInfoList)\n<case  run=\"$!caseInfo.caseRun\" name=\"$!caseInfo.caseName\"> \n</case>\n#end\n </module>\n#end\n</suite>";
	public static String ENCODE = "utf-8";
	public static Map<String, Map<String, String>> PROPERTIES_MAP = new HashMap<String, Map<String, String>>();
}
