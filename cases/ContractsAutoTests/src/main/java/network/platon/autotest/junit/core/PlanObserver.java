package network.platon.autotest.junit.core;

import java.io.File;
import java.io.FileNotFoundException;
import java.io.FileOutputStream;
import java.io.FileWriter;
import java.io.IOException;
import java.io.PrintWriter;
import java.io.StringWriter;
import java.util.ArrayList;
import java.util.List;
import java.util.Properties;

import network.platon.autotest.junit.enums.RunResult;
import network.platon.autotest.junit.enums.RunStatus;
import network.platon.autotest.junit.modules.CaseInfo;
import network.platon.autotest.junit.modules.ModuleInfo;
import network.platon.autotest.junit.modules.SuiteInfo;
import network.platon.autotest.junit.rules.DriverService;
import org.apache.velocity.Template;
import org.apache.velocity.VelocityContext;
import org.apache.velocity.app.Velocity;

/**
 * @Title: PlanInitial.java
 * @Package network.platon.autotest.junit.modules
 * @Description: TODO(用一句话描述该文件做什么)
 * @author qcxiao
 * @date 2013-12-16 上午10:28:02
 */
public class PlanObserver implements Observer {
	private Properties properties = DriverService.PROPERTIES;
	private String templatesDir = DriverService.PROPERTIES.getProperty("templatesDir", "src/main/resources/templates/");
	private String sourcesDir = DriverService.PROPERTIES.getProperty("sourcesDir", "src/test/resources/");

	@Override
	public void suiteRunStart(SuiteInfo suiteInfo) {
		// TODO Auto-generated method stub

	}

	@Override
	public void moduleRunStart(ModuleInfo moduleInfo) {
		// TODO Auto-generated method stub
		// 渲染plan.vm
		String plan = properties.getProperty("plan", "create").trim();
		if (plan.equals("create") || plan.contains("{")) {
			createPlan(LogModule.SUITE_INFO);
		}
	}

	@Override
	public void caseRunStart(CaseInfo caseInfo) {
		// TODO Auto-generated method stub

	}

	@Override
	public void suiteRunStop(SuiteInfo suiteInfo) {
		if (suiteInfo.getSuiteStatus() == RunStatus.COMPLETED) {
			errorPlan();
		}
	}

	@Override
	public void moduleRunStop(ModuleInfo moduleInfo) {
		// TODO Auto-generated method stub

	}

	@Override
	public void caseRunStop(CaseInfo caseInfo) {
		// TODO Auto-generated method stub

	}

	private void createPlan(SuiteInfo suiteInfo) {
		Properties p = new Properties();
		p.setProperty(Velocity.INPUT_ENCODING, "UTF-8");
		p.setProperty(Velocity.OUTPUT_ENCODING, "UTF-8");
		Velocity.init(p);
		// 取得velocity上下文
		VelocityContext context = new VelocityContext();
		// TODO 为后续加上项目名称和执行人做接口准备
		//String project = (String) (properties.get("project") == null ? "" : properties.get("project"));
		//String runner = (String) (properties.get("runner") == null ? "" : properties.get("runner"));
		context.put("suiteInfo", suiteInfo);
		// context.put("moduleInfo", LogModule.moduleInfo);
		context.put("encode", DriverService.ENCODE);

		File file = new File(templatesDir + "plan.vm");
		if (!file.exists()) {
			try {
				FileWriter writer = new FileWriter(templatesDir + "plan.vm", true);
				writer.write(DriverService.PLAN_VM_CONTENT);
				writer.close();
			} catch (IOException e) {
				e.printStackTrace();
			}
		}
		Template template = Velocity.getTemplate(templatesDir + "plan.vm");
		StringWriter writer = new StringWriter();
		template.merge(context, writer);
		PrintWriter filewriter;
		try {
			// filewriter = new PrintWriter(new FileOutputStream(LOG_DIR +
			// "测试套件_" + suiteInfo.getSuiteName() + ".html"), true);
			filewriter = new PrintWriter(new FileOutputStream(sourcesDir + "plan.xml"), true);
			filewriter.println(writer.toString());
			filewriter.close();
		} catch (FileNotFoundException e) {
			e.printStackTrace();
		}
	}

	private void errorPlan() {
		Properties p = new Properties();
		p.setProperty(Velocity.INPUT_ENCODING, "UTF-8");
		p.setProperty(Velocity.OUTPUT_ENCODING, "UTF-8");
		Velocity.init(p);
		// 取得velocity上下文
		VelocityContext context = new VelocityContext();
		// TODO 为后续加上项目名称和执行人做接口准备
		//String project = (String) (properties.get("project") == null ? "" : properties.get("project"));
		//String runner = (String) (properties.get("runner") == null ? "" : properties.get("runner"));
		List<ModuleInfo> errorModuleInfoList = new ArrayList<ModuleInfo>();
		for (ModuleInfo moduleInfo : LogModule.SUITE_INFO.getModuleInfoList()) {
			if (!moduleInfo.getModuleResult().equals(RunResult.PASS)) {
				List<CaseInfo> errorCaseInfoList = new ArrayList<CaseInfo>();
				for (CaseInfo caseInfo : moduleInfo.getCaseInfoList()) {
					if (!caseInfo.getCaseResult().equals(RunResult.PASS)) {
						errorCaseInfoList.add(caseInfo);
					}
				}
				moduleInfo.setCaseInfoList(errorCaseInfoList);
				errorModuleInfoList.add(moduleInfo);
			}
		}
		SuiteInfo suiteInfo = new SuiteInfo();
		suiteInfo.setSuiteName(LogModule.SUITE_INFO.getSuiteName());
		suiteInfo.setModuleInfoList(errorModuleInfoList);

		context.put("suiteInfo", suiteInfo);

		context.put("encode", DriverService.ENCODE);
		Template template = Velocity.getTemplate(templatesDir + "plan.vm");
		StringWriter writer = new StringWriter();
		template.merge(context, writer);
		PrintWriter filewriter;
		try {
			// filewriter = new PrintWriter(new FileOutputStream(LOG_DIR +
			// "测试套件_" + suiteInfo.getSuiteName() + ".html"), true);
			filewriter = new PrintWriter(new FileOutputStream(sourcesDir + "plan.error.xml"), true);
			filewriter.println(writer.toString());
			filewriter.close();
		} catch (FileNotFoundException e) {
			// TODO Auto-generated catch block
			e.printStackTrace();
		}
	}
}
