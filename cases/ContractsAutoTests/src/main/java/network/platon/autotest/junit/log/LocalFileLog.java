package network.platon.autotest.junit.log;

import java.io.File;
import java.io.FileNotFoundException;
import java.io.FileOutputStream;
import java.io.PrintWriter;
import java.io.StringWriter;
import java.util.ArrayList;
import java.util.HashMap;
import java.util.List;
import java.util.Map;
import java.util.Properties;
import java.util.stream.Collectors;

import network.platon.autotest.junit.core.LogModule;
import network.platon.autotest.junit.core.Observer;
import network.platon.autotest.junit.enums.RunResult;
import network.platon.autotest.junit.modules.CaseInfo;
import network.platon.autotest.junit.modules.LogStepInfo;
import network.platon.autotest.junit.modules.ModuleInfo;
import network.platon.autotest.junit.modules.SuiteInfo;
import network.platon.autotest.junit.rules.DriverService;
import org.apache.velocity.Template;
import org.apache.velocity.VelocityContext;
import org.apache.velocity.app.Velocity;

import network.platon.autotest.utils.FileUtil;
import network.platon.autotest.utils.ZipUtil;

public class LocalFileLog implements Observer {
	private String templatesDir = DriverService.PROPERTIES.getProperty("templatesDir", "src/main/resources/templates/");
	// 为在增加在数据推送项目中的错误报文打印到相应的用例日志路径下，修改为公共方法
	public static String moduleDir;
	private String encode = DriverService.ENCODE;
	private Properties properties = DriverService.PROPERTIES;

	public Map<String, String> param = new HashMap<String, String>();
	public List<LogStepInfo> logStepInfoList = new ArrayList<LogStepInfo>();

	/* (non-Javadoc)
	 * @see network.platon.autotest.junit.core.Observer#suiteRunStart(network.platon.autotest.junit.modules.SuiteInfo)
	 * 主要作用是当maven执行时修改日志的编码
	 */
	@Override
	public void suiteRunStart(SuiteInfo suiteInfo) {
		// 判断是否maven执行,决定编码方式，主要是兼容windows窗口控制台输出的gbk编码，后续做改进
		StackTraceElement stack[] = Thread.currentThread().getStackTrace();
		for (StackTraceElement ste : stack) {
			if ((ste.getClassName().contains("maven")||ste.getClassName().contains("mvn"))&&(System.getProperties().getProperty("os.name").contains("Windows"))) {
				encode = "GBK";
				break;
			}
		}
	}

	@Override
	public void suiteRunStop(SuiteInfo suiteInfo) {
		for (ModuleInfo moduleInfo : suiteInfo.getModuleInfoList()) {
			if (0 == moduleInfo.getModuleCaseNum()) {
				suiteInfo.setSuiteModuleNum(suiteInfo.getSuiteModuleNum() - 1);
			}
		}
		suiteInfo.getModuleInfoList().removeIf(item -> item.getModuleCaseNum() == 0);
		Properties p = new Properties();
		p.setProperty(Velocity.INPUT_ENCODING, encode);
		p.setProperty(Velocity.OUTPUT_ENCODING, encode);
		Velocity.init(p);
		// 取得velocity上下文
		VelocityContext context = new VelocityContext();
		String project = (String) (properties.get("project") == null ? "" : properties.get("project"));
		String runner = (String) (properties.get("runner") == null ? "" : properties.get("runner"));
		context.put("suiteInfo", suiteInfo);
		context.put("project", project);
		context.put("runner", runner);
		context.put("encode", encode);
		Template template = Velocity.getTemplate(templatesDir + "suiteResult.vm");
		StringWriter writer = new StringWriter();
		template.merge(context, writer);
		PrintWriter filewriter;
		try {
			filewriter = new PrintWriter(new FileOutputStream(LogModule.SUITE_PATH + suiteInfo.getSuiteName() + ".html"), true);
			filewriter.println(writer.toString());
			filewriter.close();
		} catch (FileNotFoundException e) {
			e.printStackTrace();
		}
		moduleDir = LogModule.SUITE_PATH + "Link/js/";
		File directory = new File(moduleDir);
		if (!directory.exists()) {
			directory.mkdirs();
		}
		FileUtil.copyFolder(templatesDir + "js/", moduleDir);
		// 压缩日志文件
		if (properties.getProperty("reportZip") != null && !properties.getProperty("reportZip").toLowerCase().trim().equals("false")) {
			ZipUtil.zip(LogModule.SUITE_PATH, "report.zip");
		}

	}

	@Override
	public void moduleRunStart(ModuleInfo moduleInfo) {
		//List<Boolean> caseRunList = moduleInfo.getCaseInfoList().stream().map(CaseInfo::getCaseRun).collect(Collectors.toList());
//		List<Boolean> caseRunList = moduleInfo.getCaseInfoList().stream().map(CaseInfo::getCaseRun)
//				.filter(item -> item.booleanValue() == true).collect(Collectors.toList());
//		System.out.println("caseRunList:" + caseRunList);
//		System.out.println("caseRunList.size():" + caseRunList.size());
//		if (null != caseRunList && caseRunList.size() == 0) {
//			return;
//		}
		// moduleInfo = LogModule.moduleInfo;
		moduleDir = LogModule.SUITE_PATH + "Link/测试模块_" + moduleInfo.getModuleName() + "/";
		File directory = new File(moduleDir);
		if (!directory.exists()) {
			directory.mkdirs();
		}
	}

	@Override
	public void moduleRunStop(ModuleInfo moduleInfo) {
		// moduleInfo = LogModule.moduleInfo;
		// 下面是本地日志,后续隔离开
		Properties p = new Properties();
		p.setProperty(Velocity.INPUT_ENCODING, encode);
		p.setProperty(Velocity.OUTPUT_ENCODING, encode);
		Velocity.init(p);
		// 取得velocity上下文
		VelocityContext context = new VelocityContext();
		context.put("moduleInfo", moduleInfo);
		context.put("encode", encode);
		templatesDir = properties.getProperty("templatesDir", templatesDir);
		Template template = Velocity.getTemplate(templatesDir + "moduleResult.vm");
		StringWriter writer = new StringWriter();
		template.merge(context, writer);
		PrintWriter filewriter;
		try {
			filewriter = new PrintWriter(new FileOutputStream(LogModule.SUITE_PATH + "Link/测试模块_" + moduleInfo.getModuleName() + ".html"), true);
			filewriter.println(writer.toString());
			filewriter.close();
		} catch (FileNotFoundException e) {
			// TODO Auto-generated catch block
			e.printStackTrace();
		}

	}

	@Override
	public void caseRunStart(CaseInfo caseInfo) {
		// TODO Auto-generated method stub

	}

	@Override
	public void caseRunStop(CaseInfo caseInfo) {
		// caseInfo = LogModule.caseInfo;
		// 只有最后一次重跑时，才记录本地日志
		if (caseInfo.getCaseResult().equals(RunResult.FAIL) && caseInfo.getCaseRerunNum() >= 0) {
			return;
		}
		Properties p = new Properties();

		p.setProperty(Velocity.INPUT_ENCODING, encode);
		p.setProperty(Velocity.OUTPUT_ENCODING, encode);

		Velocity.init(p);

		// 取得velocity上下文
		VelocityContext context = new VelocityContext();
		context.put("logStepInfoList", LogModule.logStepInfoList);
		context.put("caseInfo", caseInfo);
		context.put("encode", encode);
		templatesDir = properties.getProperty("templatesDir", templatesDir);
		Template template = Velocity.getTemplate(templatesDir + "caseResult.vm");

		StringWriter writer = new StringWriter();
		template.merge(context, writer);

		PrintWriter filewriter;
		try {
			filewriter = new PrintWriter(new FileOutputStream(moduleDir + "测试用例_" + caseInfo.getCaseName() + ".html"), true);
			filewriter.println(writer.toString());
			filewriter.close();
		} catch (FileNotFoundException e) {
			// TODO Auto-generated catch block
			e.printStackTrace();
		}

	}

}
