package network.platon.autotest.junit.rules;

import java.util.List;
import java.util.concurrent.Callable;
import java.util.concurrent.ExecutionException;
import java.util.concurrent.ExecutorService;
import java.util.concurrent.Executors;
import java.util.concurrent.Future;

/**
 * 并发处理器 适用于如下场景（举例）： 一个任务队列，
 * 有150个任务需要并发处理，使用此对象，可以每次并发执行20次（可设置），则总共串行执行8次并发，可获取执行结果
 * 
 * @param <T>
 *            类型T限制为任务Callable使用的
 */
public class ConcurrentExcutor<T> {
	/** 非空，所有任务数组 */
	private Callable<T>[] tasks;

	/** 非空，每次并发需要处理的任务数 */
	private int numb;

	/**
	 * 可选，存放返回结果，这里有个限制 T必须为Callable返回的类型T
	 */
	private List<T> result;

	/**
	 * 无参构造
	 */
	public ConcurrentExcutor() {
		super();
	}

	/**
	 * 不需要返回结果的任务用此创建对象
	 * 
	 * @param tasks
	 * @param numb
	 */
	public ConcurrentExcutor(Callable<T>[] tasks, int numb) {
		super();
		this.tasks = tasks;
		this.numb = numb;
	}

	/**
	 * 需要结果集用此方法创建对象
	 * 
	 * @param tasks
	 * @param numb
	 * @param result
	 */
	public ConcurrentExcutor(Callable<T>[] tasks, int numb, List<T> result) {
		super();
		this.tasks = tasks;
		this.numb = numb;
		this.result = result;
	}

	public void excute() {
		// 参数校验
		if (tasks == null || numb < 1) {
			return;
		}

		// 待处理的任务数
		int num = tasks.length;
		if (num == 0) {
			return;
		}

		// 第一层循环，每numb条数据作为一次并发
		for (int i = 0; i < (int) Math.floor(num / numb) + 1; i++) {
			// 用于记录此次numb条任务的处理结果
			Future[] futureArray;
			if (numb > num) {
				futureArray = new Future[num];
			} else {
				futureArray = new Future[numb];
			}
			// 创建线程容器
			ExecutorService es = Executors.newCachedThreadPool();

			// 第二层循环，针对这numb条数据进行处理
			for (int j = i * numb; j < (i + 1) * numb; j++) {
				// 如果超出数组长度，退出循环
				if (j + 1 > num) {
					break;
				}
				// 执行任务，并设置Future到数组中
				futureArray[j % numb] = es.submit(tasks[j]);
			}

			// 将结果放入result中
			if (result != null) {
				for (int j = 0; j < futureArray.length; j++) {
					try {
						if (futureArray[j] != null) {
							Object o = futureArray[j].get();
							result.add((T) o);
						}
					} catch (InterruptedException e) {
						System.out.println("处理Future时发生InterruptedException异常，目标Future为： " + futureArray[j].toString());
						e.printStackTrace();
					} catch (ExecutionException e) {
						System.out.println("处理Future时发生ExecutionException异常，目标Future为： " + futureArray[j].toString());
						e.printStackTrace();
					}
				}
			}
			es.shutdown();
		}
	}
}