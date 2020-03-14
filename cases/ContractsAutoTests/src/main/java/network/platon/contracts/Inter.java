package network.platon.contracts;

import java.math.BigInteger;
import java.util.Arrays;
import org.web3j.abi.TypeReference;
import org.web3j.abi.datatypes.Function;
import org.web3j.abi.datatypes.Type;
import org.web3j.abi.datatypes.generated.Uint256;
import org.web3j.crypto.Credentials;
import org.web3j.protocol.Web3j;
import org.web3j.protocol.core.RemoteCall;
import org.web3j.tx.Contract;
import org.web3j.tx.TransactionManager;
import org.web3j.tx.gas.GasProvider;

/**
 * <p>Auto generated code.
 * <p><strong>Do not modify!</strong>
 * <p>Please use the <a href="https://docs.web3j.io/command_line.html">web3j command line tools</a>,
 * or the org.web3j.codegen.SolidityFunctionWrapperGenerator in the 
 * <a href="https://github.com/web3j/web3j/tree/master/codegen">codegen module</a> to update.
 *
 * <p>Generated with web3j version 0.7.5.3-SNAPSHOT.
 */
public class Inter extends Contract {
    private static final String BINARY = "608060405234801561001057600080fd5b50610157806100206000396000f3fe608060405234801561001057600080fd5b50600436106100415760003560e01c8063b8b1feb414610046578063ca77156f14610088578063e2179b8e146100ca575b600080fd5b6100726004803603602081101561005c57600080fd5b81019080803590602001909291905050506100e8565b6040518082815260200191505060405180910390f35b6100b46004803603602081101561009e57600080fd5b81019080803590602001909291905050506100f5565b6040518082815260200191505060405180910390f35b6100d2610102565b6040518082815260200191505060405180910390f35b6000600382019050919050565b6000600282019050919050565b600061011060016002610115565b905090565b600081830190509291505056fea265627a7a7231582083c9395283d0b0c96c06f6de14ae14bf76305c09355fc873fb358e1d7ac35d5764736f6c634300050d0032";

    public static final String FUNC_FE = "fe";

    public static final String FUNC_FPUB = "fpub";

    public static final String FUNC_G = "g";

    @Deprecated
    protected Inter(String contractAddress, Web3j web3j, Credentials credentials, BigInteger gasPrice, BigInteger gasLimit) {
        super(BINARY, contractAddress, web3j, credentials, gasPrice, gasLimit);
    }

    protected Inter(String contractAddress, Web3j web3j, Credentials credentials, GasProvider contractGasProvider) {
        super(BINARY, contractAddress, web3j, credentials, contractGasProvider);
    }

    @Deprecated
    protected Inter(String contractAddress, Web3j web3j, TransactionManager transactionManager, BigInteger gasPrice, BigInteger gasLimit) {
        super(BINARY, contractAddress, web3j, transactionManager, gasPrice, gasLimit);
    }

    protected Inter(String contractAddress, Web3j web3j, TransactionManager transactionManager, GasProvider contractGasProvider) {
        super(BINARY, contractAddress, web3j, transactionManager, contractGasProvider);
    }

    public RemoteCall<BigInteger> fe(BigInteger a) {
        final Function function = new Function(FUNC_FE, 
                Arrays.<Type>asList(new org.web3j.abi.datatypes.generated.Uint256(a)), 
                Arrays.<TypeReference<?>>asList(new TypeReference<Uint256>() {}));
        return executeRemoteCallSingleValueReturn(function, BigInteger.class);
    }

    public RemoteCall<BigInteger> fpub(BigInteger a) {
        final Function function = new Function(FUNC_FPUB, 
                Arrays.<Type>asList(new org.web3j.abi.datatypes.generated.Uint256(a)), 
                Arrays.<TypeReference<?>>asList(new TypeReference<Uint256>() {}));
        return executeRemoteCallSingleValueReturn(function, BigInteger.class);
    }

    public RemoteCall<BigInteger> g() {
        final Function function = new Function(FUNC_G, 
                Arrays.<Type>asList(), 
                Arrays.<TypeReference<?>>asList(new TypeReference<Uint256>() {}));
        return executeRemoteCallSingleValueReturn(function, BigInteger.class);
    }

    public static RemoteCall<Inter> deploy(Web3j web3j, Credentials credentials, GasProvider contractGasProvider) {
        return deployRemoteCall(Inter.class, web3j, credentials, contractGasProvider, BINARY, "");
    }

    @Deprecated
    public static RemoteCall<Inter> deploy(Web3j web3j, Credentials credentials, BigInteger gasPrice, BigInteger gasLimit) {
        return deployRemoteCall(Inter.class, web3j, credentials, gasPrice, gasLimit, BINARY, "");
    }

    public static RemoteCall<Inter> deploy(Web3j web3j, TransactionManager transactionManager, GasProvider contractGasProvider) {
        return deployRemoteCall(Inter.class, web3j, transactionManager, contractGasProvider, BINARY, "");
    }

    @Deprecated
    public static RemoteCall<Inter> deploy(Web3j web3j, TransactionManager transactionManager, BigInteger gasPrice, BigInteger gasLimit) {
        return deployRemoteCall(Inter.class, web3j, transactionManager, gasPrice, gasLimit, BINARY, "");
    }

    @Deprecated
    public static Inter load(String contractAddress, Web3j web3j, Credentials credentials, BigInteger gasPrice, BigInteger gasLimit) {
        return new Inter(contractAddress, web3j, credentials, gasPrice, gasLimit);
    }

    @Deprecated
    public static Inter load(String contractAddress, Web3j web3j, TransactionManager transactionManager, BigInteger gasPrice, BigInteger gasLimit) {
        return new Inter(contractAddress, web3j, transactionManager, gasPrice, gasLimit);
    }

    public static Inter load(String contractAddress, Web3j web3j, Credentials credentials, GasProvider contractGasProvider) {
        return new Inter(contractAddress, web3j, credentials, contractGasProvider);
    }

    public static Inter load(String contractAddress, Web3j web3j, TransactionManager transactionManager, GasProvider contractGasProvider) {
        return new Inter(contractAddress, web3j, transactionManager, contractGasProvider);
    }
}
