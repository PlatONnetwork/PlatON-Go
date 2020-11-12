package network.platon.autotest.exception;

public class RemoteWebException extends RuntimeException {
	private static final long serialVersionUID = 3341764449059630537L;

	public RemoteWebException() {
		super();
	}

	public RemoteWebException(String message, Throwable cause) {
		super(message, cause);
	}

	public RemoteWebException(String message) {
		super(message);
	}

	public RemoteWebException(Throwable cause) {
		super(cause);
	}
}
