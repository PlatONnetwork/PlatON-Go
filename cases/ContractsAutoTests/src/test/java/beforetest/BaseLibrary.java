package beforetest;

import org.web3j.crypto.Credentials;
import org.web3j.crypto.RawTransaction;
import org.web3j.crypto.TransactionEncoder;
import org.web3j.protocol.Web3j;
import org.web3j.protocol.Web3jService;
import org.web3j.protocol.core.DefaultBlockParameterName;
import org.web3j.protocol.core.methods.response.PlatonGetTransactionCount;
import org.web3j.protocol.core.methods.response.PlatonSendTransaction;
import org.web3j.protocol.core.methods.response.TransactionReceipt;
import org.web3j.protocol.exceptions.TransactionException;
import org.web3j.protocol.http.HttpService;
import org.web3j.tx.response.PollingTransactionReceiptProcessor;
import org.web3j.utils.Numeric;

import java.io.IOException;
import java.math.BigInteger;

public class BaseLibrary {
	protected static final BigInteger GAS_LIMIT = BigInteger.valueOf(4700000);
	protected static final BigInteger GAS_PRICE = BigInteger.valueOf(1000000000L);

	public static final int DEFAULT_POLLING_ATTEMPTS_PER_TX_HASH = 40;
	public static final long DEFAULT_POLLING_FREQUENCY = 2 * 1000;

	private Credentials credentials;
	private Web3j web3j;
	private long chainId;

	public BaseLibrary(Credentials credentials, Web3j web3j, long chainId) {
		this.credentials = credentials;
		this.web3j = web3j;
		this.chainId = chainId;
	}

	public TransactionReceipt deploy(BigInteger gasPrice, BigInteger gasLimit, String data) throws Exception {
		PlatonGetTransactionCount platonGetTransactionCount = web3j
				.platonGetTransactionCount(credentials.getAddress(), DefaultBlockParameterName.LATEST).send();
		BigInteger nonce = platonGetTransactionCount.getTransactionCount();

		String to = "";
		BigInteger value = BigInteger.valueOf(0L);
		RawTransaction rawTransaction = RawTransaction.createTransaction(nonce, gasPrice, gasLimit, to, value, data);

		byte[] signedMessage = TransactionEncoder.signMessage(rawTransaction, chainId, credentials);
		String hexValue = Numeric.toHexString(signedMessage);
		PlatonSendTransaction platonSendTransaction = web3j.platonSendRawTransaction(hexValue).send();

		return processResponse(platonSendTransaction);
	}

	private TransactionReceipt processResponse(PlatonSendTransaction transactionResponse) throws IOException, TransactionException {
		if (transactionResponse.hasError()) {
			throw new RuntimeException("Error processing transaction request: " + transactionResponse.getError().getMessage());
		}

		String transactionHash = transactionResponse.getTransactionHash();

		return new PollingTransactionReceiptProcessor(web3j, DEFAULT_POLLING_FREQUENCY, DEFAULT_POLLING_ATTEMPTS_PER_TX_HASH)
				.waitForTransactionReceipt(transactionHash);
	}

	public static void main(String[] args) throws Exception {
		long chainId = 100L;
		String nodeUrl = "http://10.10.8.21:8804";
		String privateKey = "11e20dc277fafc4bc008521adda4b79c2a9e403131798c94eacb071005d43532";

		Credentials credentials = Credentials.create(privateKey);
		Web3jService web3jService = new HttpService(nodeUrl);
		Web3j web3j = Web3j.build(web3jService);

		BaseLibrary baseLibrary = new BaseLibrary(credentials, web3j, chainId);
		String data = "608060405234801561001057600080fd5b506101fa806100206000396000f30060806040526004361061006d576000357c0100000000000000000000000000000000000000000000000000000000900463ffffffff1680630bb2b6961461007257806320de797e1461009d57806325b29d84146100d1578063c6d8d657146100fc578063c6f8a3b714610127575b600080fd5b34801561007e57600080fd5b50610087610152565b6040518082815260200191505060405180910390f35b6100bb6004803603810190808035906020019092919050505061015c565b6040518082815260200191505060405180910390f35b3480156100dd57600080fd5b506100e6610196565b6040518082815260200191505060405180910390f35b34801561010857600080fd5b506101116101bb565b6040518082815260200191505060405180910390f35b34801561013357600080fd5b5061013c6101c4565b6040518082815260200191505060405180910390f35b6000600254905090565b60006301e13380600081905550680dd2d5fcf3bc9c000060018190555060ff600281905550680dd2d5fcf3bc9c0000600381905550919050565b6000671bc16d674ec80000801415156101b257600090506101b8565b60015490505b90565b60008054905090565b60006003549050905600a165627a7a7230582019a721a340d19b63e87d4013d34cd7f47f473d98ee2fd48e43affd15939987d80029";
		TransactionReceipt receipt = baseLibrary.deploy(BaseLibrary.GAS_PRICE, BaseLibrary.GAS_LIMIT, data);
		System.err.println("status >>>> " + receipt.getStatus());
		System.err.println("contract address >>>> " + receipt.getContractAddress());
	}
}
