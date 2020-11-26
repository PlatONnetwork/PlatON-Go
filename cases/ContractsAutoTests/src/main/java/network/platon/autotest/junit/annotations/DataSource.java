package network.platon.autotest.junit.annotations;

import java.lang.annotation.ElementType;
import java.lang.annotation.Retention;
import java.lang.annotation.RetentionPolicy;
import java.lang.annotation.Target;

import network.platon.autotest.junit.enums.DataSourceType;

@Retention(RetentionPolicy.RUNTIME)
@Target({ ElementType.METHOD })
public @interface DataSource {
	public String file() default "";

	public DataSourceType type() default DataSourceType.EXCEL;

	public String sheetName() default "";
	
	public String showName() default "";
	
	public String author() default "";
	
	public String expert() default "";

	public String sourcePrefix() default "";

	public int executionSequence() default 1;

}
