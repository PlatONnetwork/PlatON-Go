package network.platon.autotest.junit.rules;

import static org.hamcrest.Matchers.containsString;
import static org.hamcrest.Matchers.equalTo;
import static org.hamcrest.Matchers.hasItem;
import static org.hamcrest.Matchers.is;
import static org.hamcrest.Matchers.not;
import static org.hamcrest.collection.IsIn.isIn;
import static org.junit.Assert.assertThat;
import java.util.ArrayList;
import java.util.Collection;
import java.util.List;
import java.util.concurrent.Callable;
import org.hamcrest.Matcher;
import org.junit.rules.Verifier;
import org.junit.runners.model.MultipleFailureException;
import org.springframework.web.util.HtmlUtils;
import network.platon.autotest.junit.core.LogModule;
import network.platon.autotest.junit.enums.RunResult;
import network.platon.autotest.junit.enums.StepType;
import network.platon.autotest.junit.enums.TestSuiteType;

public class AssertCollector extends Verifier {
	private List<Throwable> errors = new ArrayList<Throwable>();
	private Boolean hasError = false;
	// 设置断言失败时是否跳出执行
	private Boolean isBreak = false;

	@Override
	protected void verify() throws Throwable {
		MultipleFailureException.assertEmpty(errors);
	}

	/**
	 * Adds a Throwable to the table. Execution continues, but the test will
	 * fail at the end.
	 */
	private void addError(Throwable error) {
		errors.add(error);
	}

	/**
	 * Adds a failure to the table if {@code matcher} does not match
	 * {@code value}. Execution continues, but the test will fail at the end if
	 * the match fails.
	 */
	// private <T> void checkThat(final T value, final Matcher<T> matcher) {
	// checkThat("", value, matcher);
	// }

	/**
	 * Adds a failure with the given {@code reason} to the table if
	 * {@code matcher} does not match {@code value}. Execution continues, but
	 * the test will fail at the end if the match fails.
	 */

	private <T> void checkThat(final String reason, final T value, final Matcher<T> matcher) {
		if (getIsSkip())
			return;
		hasError = false;
		checkSucceeds(new Callable<Object>() {
			public Object call() throws Exception {
				assertThat(reason, value, matcher);
				return value;
			}
		});
	}

	/**
	 * Adds to the table the exception, if any, thrown from {@code callable}.
	 * Execution continues, but the test will fail at the end if
	 * {@code callable} threw an exception.
	 */
	private Object checkSucceeds(Callable<Object> callable) {
		try {
			return callable.call();
		} catch (Throwable e) {
			hasError = true;
			addError(e);
			return null;
		}
	}

	/**
	 * 验证传入的boolean值是否为true
	 * 
	 * @param actual
	 */
	public void assertTrue(boolean actual) {
		assertTrue(actual, "");
	}


	/**
	 * 验证传入的boolean值是否为true
	 * 
	 * @param actual
	 * @param message
	 */
	public void assertTrue(boolean actual, String message) {
		String methodInfo = "assertTrue(" + String.valueOf(actual) + "," + message + ")";
		checkThat(message, actual, is(true));
		assertLog(methodInfo, message);
	}

	/**
	 * 验证传入的boolean值是否为 false
	 * 
	 * @param actual
	 */
	public void assertFalse(boolean actual) {
		assertFalse(actual, "");
	}


	/**
	 * 验证传入的boolean值是否为 false
	 * 
	 * @param actual
	 * @param message
	 */
	public void assertFalse(boolean actual, String message) {
		String methodInfo = "assertFalse(" + String.valueOf(actual) + "," + message + ")";
		checkThat(message, actual, is(false));
		assertLog(methodInfo, message);
	}

	public <T> void assertEqual(T actual, T expect) {
		assertEqual(actual, expect, "");
	}


	/**
	 * 验证传入实际值和期望值相等
	 * 
	 * @param <T>
	 * @param actual
	 * @param expect
	 * @param message
	 */
	public <T> void assertEqual(T actual, T expect, String message) {
		String methodInfo = "assertEqual(" + String.valueOf(actual) + "," + String.valueOf(expect) + "," + message + ")";
		checkThat(message, actual, equalTo(expect));
		assertLog(methodInfo, String.valueOf(actual), String.valueOf(expect), message);
	}

	/**
	 * 验证传入的字符串是否符合传入的正则表达式
	 * 
	 * @param actual
	 * @param regxp
	 */
	public void assertMatch(String actual, String regxp) {
		assertMatch(actual, regxp, "");
	}


	/**
	 * 验证传入的字符串是否符合传入的正则表达式
	 * 
	 * @param actual
	 * @param regxp
	 * @param message
	 */
	public void assertMatch(String actual, String regxp, String message) {
		checkThat(message, actual.matches(regxp), is(true));
		String methodInfo = "assertMatch(" + String.valueOf(actual) + "," + String.valueOf(regxp) + "," + message + ")";
		assertLog(methodInfo, message);

	}

	/**
	 * 校验实际值是否包含期望值
	 * @param actual 实际值
	 * @param expect 期望值
	 */
	public void assertContains(String actual, String expect) {
		assertContains(actual, expect, "");
	}


	/**
	 * 验证传入的字符串是否包含字符串
	 * @param actual
	 * @param expect
	 * @param message
	 */
	public void assertContains(String actual, String expect, String message) {
		checkThat(message, actual, is(containsString(expect)));
		String methodInfo = "assertContains(" + String.valueOf(actual) + "," + String.valueOf(expect) + "," + message + ")";
		assertLog(methodInfo, message);
	}


	/**
	 * 验证传入的集合是包含指定对象
	 * 
	 * @param actual
	 */
	public <T> void assertContains(Collection<T> actual, T expect, String message) {
		checkThat(message, actual, hasItem(expect));
		String methodInfo = "assertContains(" + String.valueOf(actual) + "," + String.valueOf(expect) + "," + message + ")";
		assertLog(methodInfo, message);
	}

	/**
	 * 校验实际值对象是否存在于期望值的对象集合中
	 * @param actual 实际值
	 * @param expect 期望值
	 * @param message 该断言的描述信息（类似备注功能）
	 * @remark 存在时，校验通过；不存在时，会自动截图，并记录错误信息
	 * @Example collector.assertIn("abc", Arrays.asList("abc", "3", "rdf4", "被包含关系");//验证实际值对象存在于期望值的对象集合中
	 */

	/**
	 * 验证传入的对象不在集合中
	 * 
	 * @param actual 实际值
	 * @param expect 期望值
	 * @param message 该断言的描述信息（类似备注功能）
	 * @return
	 * @remark 不存在时，校验通过；存在时，会自动截图，并记录错误信息
	 * @Example collector.assertOut(3, Arrays.asList(3, 5, 66, 8), "未被包含关系");//验证实际值对象不存在于期望值的对象集合中
	 */
//	public <T> void assertOut(T actual, Collection<T> expect, String message) {
//		checkThat(message, actual, not(isIn(expect)));
//		String methodInfo = "assertOut(" + String.valueOf(actual) + "," + String.valueOf(expect) + "," + message + ")";
//		assertLog(methodInfo, message);
//	}
	/**
	 * 日志中输入期望值与实际值的assertEqual断言
	 * 且把描述信息置为空
	 * @param methodInfo
	 * @param actual
	 * @param expect
	 * @param message
	 */
	private void assertLog(String methodInfo,String actual, String expect, String message) {
		if (getIsSkip())
			return;
		//WebDriver driver = browser == null ? Browser.currentDriver : browser.driver;
		// String stepTime = DateUtil.dateToStr(new Date(), "yyyyMMdd-HHmmss");
		if (hasError) {
			String errorMessage = errors.get(errors.size() - 1).getMessage();
			errorMessage = HtmlUtils.htmlEscape(errorMessage);
			LogModule.logStepFail(StepType.ASSERT, message, actual, expect, RunResult.FAIL, errorMessage, TestSuiteType.WEB_UI);
			if (isBreak) {
				throw new RuntimeException(errorMessage);
			}

		} else {
			LogModule.logStepPass(StepType.ASSERT, message, actual, expect, RunResult.PASS);
		}
	}

	/*
	 * 断言统一日志
	 */
	private void assertLog(String methodInfo, String message) {
		if (getIsSkip())
			return;
		//WebDriver driver = browser == null ? Browser.currentDriver : browser.driver;
		// String stepTime = DateUtil.dateToStr(new Date(), "yyyyMMdd-HHmmss");
		if (hasError) {
			String errorMessage = errors.get(errors.size() - 1).getMessage();
			errorMessage = HtmlUtils.htmlEscape(errorMessage);
			LogModule.logStepFail(StepType.ASSERT, methodInfo, RunResult.FAIL, errorMessage, TestSuiteType.WEB_UI);
			if (isBreak) {
				throw new RuntimeException(errorMessage);
			}

		} else {
			LogModule.logStepPass(StepType.ASSERT, methodInfo, RunResult.PASS);
		}
	}

	/**
	 * 通过设置该值，可以实现断言比较失败时，是否终止当前用例
	 * @param isBreak true|false 是否终止当前用例执行
	 * @Example
	 * <div>如果这个断言验证失败的话，就会终止当前用例的执行</div>
	 * <div>collector.setIsBreak(true);</div>
	 * <div>collector.assertTrue(browser.div("class=>booking").table("class=>tableinfo2").cell("text=>810.0(儿童)").exists(2), "儿童价格");</div>
	 */
	public void setIsBreak(Boolean isBreak) {
		this.isBreak = isBreak;
	}

	public Boolean getIsBreak() {
		return isBreak;
	}

	/**
	 * 执行失败日志
	 * @param stepDesc 步骤的描述
	 * @param failReason 失败的原因
	 * @Example collector.logStepFail("订单编号校验", "编号不规范");
	 */
	public void logStepFail(String stepDesc, String failReason) {
		LogModule.logStepFail(StepType.CUSTOM, stepDesc, RunResult.FAIL, failReason, TestSuiteType.WEB_UI);
	}

	/**
	 * 校验实际值是否为假
	 * @param stepDesc 步骤的描述
	 * @Example collector.logStepPass("订单编号校验");
	 */
	public void logStepPass(String stepDesc) {
		LogModule.logStepPass(StepType.CUSTOM, stepDesc, RunResult.PASS);
	}

	// 跳过断言
	public Boolean getIsSkip() {
		Boolean isSkip = !Boolean.valueOf(DriverService.PROPERTIES.getProperty("caseAssert"));
		return isSkip;
	}
}
