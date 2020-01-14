package network.platon.contracts;

import java.math.BigInteger;
import java.util.Arrays;
import java.util.Collections;
import org.web3j.abi.TypeReference;
import org.web3j.abi.datatypes.Function;
import org.web3j.abi.datatypes.Type;
import org.web3j.abi.datatypes.generated.Uint256;
import org.web3j.crypto.Credentials;
import org.web3j.protocol.Web3j;
import org.web3j.protocol.core.RemoteCall;
import org.web3j.protocol.core.methods.response.TransactionReceipt;
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
 * <p>Generated with web3j version 0.7.5.0.
 */
public class PlatONToken extends Contract {
    private static final String BINARY = "608060405234801561001057600080fd5b50610162806100206000396000f3fe608060405234801561001057600080fd5b50600436106100575760003560e01c8063249bb3731461005c57806375efc40d1461007a578063b69ef8a814610098578063c951fdf6146100b6578063eecb9ce9146100d4575b600080fd5b6100646100f2565b6040518082815260200191505060405180910390f35b610082610102565b6040518082815260200191505060405180910390f35b6100a061010f565b6040518082815260200191505060405180910390f35b6100be610115565b6040518082815260200191505060405180910390f35b6100dc61011e565b6040518082815260200191505060405180910390f35b6000670de0b6b3a7640000905090565b600064e8d4a51000905090565b60005481565b60006001905090565b600066038d7ea4c6800090509056fea265627a7a72315820c1ecb2ecb3d4f6c7f4694ba0fcce16e035ae1aa130caace20a9228a87185d8b064736f6c634300050c0032";

    public static final String FUNC_PFINNEY = "Pfinney";

    public static final String FUNC_PLAT = "Plat";

    public static final String FUNC_PSZABO = "Pszabo";

    public static final String FUNC_PVON = "Pvon";

    public static final String FUNC_BALANCE = "balance";

    @Deprecated
    protected PlatONToken(String contractAddress, Web3j web3j, Credentials credentials, BigInteger gasPrice, BigInteger gasLimit) {
        super(BINARY, contractAddress, web3j, credentials, gasPrice, gasLimit);
    }

    protected PlatONToken(String contractAddress, Web3j web3j, Credentials credentials, GasProvider contractGasProvider) {
        super(BINARY, contractAddress, web3j, credentials, contractGasProvider);
    }

    @Deprecated
    protected PlatONToken(String contractAddress, Web3j web3j, TransactionManager transactionManager, BigInteger gasPrice, BigInteger gasLimit) {
        super(BINARY, contractAddress, web3j, transactionManager, gasPrice, gasLimit);
    }

    protected PlatONToken(String contractAddress, Web3j web3j, TransactionManager transactionManager, GasProvider contractGasProvider) {
        super(BINARY, contractAddress, web3j, transactionManager, contractGasProvider);
    }

    public RemoteCall<TransactionReceipt> Pfinney() {
        final Function function = new Function(
                FUNC_PFINNEY, 
                Arrays.<Type>asList(), 
                Collections.<TypeReference<?>>emptyList());
        return executeRemoteCallTransaction(function);
    }

    public RemoteCall<TransactionReceipt> Plat() {
        final Function function = new Function(
                FUNC_PLAT, 
                Arrays.<Type>asList(), 
                Collections.<TypeReference<?>>emptyList());
        return executeRemoteCallTransaction(function);
    }

    public RemoteCall<TransactionReceipt> Pszabo() {
        final Function function = new Function(
                FUNC_PSZABO, 
                Arrays.<Type>asList(), 
                Collections.<TypeReference<?>>emptyList());
        return executeRemoteCallTransaction(function);
    }

    public RemoteCall<TransactionReceipt> Pvon() {
        final Function function = new Function(
                FUNC_PVON, 
                Arrays.<Type>asList(), 
                Collections.<TypeReference<?>>emptyList());
        return executeRemoteCallTransaction(function);
    }

    public RemoteCall<BigInteger> balance() {
        final Function function = new Function(FUNC_BALANCE, 
                Arrays.<Type>asList(), 
                Arrays.<TypeReference<?>>asList(new TypeReference<Uint256>() {}));
        return executeRemoteCallSingleValueReturn(function, BigInteger.class);
    }

    public static RemoteCall<PlatONToken> deploy(Web3j web3j, Credentials credentials, GasProvider contractGasProvider) {
        return deployRemoteCall(PlatONToken.class, web3j, credentials, contractGasProvider, BINARY, "");
    }

    @Deprecated
    public static RemoteCall<PlatONToken> deploy(Web3j web3j, Credentials credentials, BigInteger gasPrice, BigInteger gasLimit) {
        return deployRemoteCall(PlatONToken.class, web3j, credentials, gasPrice, gasLimit, BINARY, "");
    }

    public static RemoteCall<PlatONToken> deploy(Web3j web3j, TransactionManager transactionManager, GasProvider contractGasProvider) {
        return deployRemoteCall(PlatONToken.class, web3j, transactionManager, contractGasProvider, BINARY, "");
    }

    @Deprecated
    public static RemoteCall<PlatONToken> deploy(Web3j web3j, TransactionManager transactionManager, BigInteger gasPrice, BigInteger gasLimit) {
        return deployRemoteCall(PlatONToken.class, web3j, transactionManager, gasPrice, gasLimit, BINARY, "");
    }

    @Deprecated
    public static PlatONToken load(String contractAddress, Web3j web3j, Credentials credentials, BigInteger gasPrice, BigInteger gasLimit) {
        return new PlatONToken(contractAddress, web3j, credentials, gasPrice, gasLimit);
    }

    @Deprecated
    public static PlatONToken load(String contractAddress, Web3j web3j, TransactionManager transactionManager, BigInteger gasPrice, BigInteger gasLimit) {
        return new PlatONToken(contractAddress, web3j, transactionManager, gasPrice, gasLimit);
    }

    public static PlatONToken load(String contractAddress, Web3j web3j, Credentials credentials, GasProvider contractGasProvider) {
        return new PlatONToken(contractAddress, web3j, credentials, contractGasProvider);
    }

    public static PlatONToken load(String contractAddress, Web3j web3j, TransactionManager transactionManager, GasProvider contractGasProvider) {
        return new PlatONToken(contractAddress, web3j, transactionManager, contractGasProvider);
    }
}
