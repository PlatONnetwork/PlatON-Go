package network.platon.contracts.evm;

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
 * <p>Please use the <a href="https://github.com/PlatONnetwork/client-sdk-java/releases">platon-web3j command line tools</a>,
 * or the org.web3j.codegen.SolidityFunctionWrapperGenerator in the 
 * <a href="https://github.com/PlatONnetwork/client-sdk-java/tree/master/codegen">codegen module</a> to update.
 *
 * <p>Generated with web3j version 0.13.1.5.
 */
public class Modifiers extends Contract {
    private static final String BINARY = "6080604052600a600055348015601457600080fd5b5060b8806100236000396000f3fe6080604052348015600f57600080fd5b506004361060325760003560e01c806366e41cb71460375780636b59084d146053575b600080fd5b603d605b565b6040518082815260200191505060405180910390f35b60596064565b005b60008054905090565b6000546000819050600080549050600c600081905550506000549050600b600081905550505056fea165627a7a723058207bc0d13daf7c2717ba451a33441906127cb4e90fc790711e273fea34ce53d4660029";

    public static final String FUNC_TEST2 = "test2";

    public static final String FUNC_TEST1 = "test1";

    protected Modifiers(String contractAddress, Web3j web3j, Credentials credentials, GasProvider contractGasProvider, Long chainId) {
        super(BINARY, contractAddress, web3j, credentials, contractGasProvider, chainId);
    }

    protected Modifiers(String contractAddress, Web3j web3j, TransactionManager transactionManager, GasProvider contractGasProvider, Long chainId) {
        super(BINARY, contractAddress, web3j, transactionManager, contractGasProvider, chainId);
    }

    public RemoteCall<BigInteger> test2() {
        final Function function = new Function(FUNC_TEST2, 
                Arrays.<Type>asList(), 
                Arrays.<TypeReference<?>>asList(new TypeReference<Uint256>() {}));
        return executeRemoteCallSingleValueReturn(function, BigInteger.class);
    }

    public RemoteCall<TransactionReceipt> test1() {
        final Function function = new Function(
                FUNC_TEST1, 
                Arrays.<Type>asList(), 
                Collections.<TypeReference<?>>emptyList());
        return executeRemoteCallTransaction(function);
    }

    public static RemoteCall<Modifiers> deploy(Web3j web3j, Credentials credentials, GasProvider contractGasProvider, Long chainId) {
        return deployRemoteCall(Modifiers.class, web3j, credentials, contractGasProvider, BINARY,  "", chainId);
    }

    public static RemoteCall<Modifiers> deploy(Web3j web3j, TransactionManager transactionManager, GasProvider contractGasProvider, Long chainId) {
        return deployRemoteCall(Modifiers.class, web3j, transactionManager, contractGasProvider, BINARY,  "", chainId);
    }

    public static Modifiers load(String contractAddress, Web3j web3j, Credentials credentials, GasProvider contractGasProvider, Long chainId) {
        return new Modifiers(contractAddress, web3j, credentials, contractGasProvider, chainId);
    }

    public static Modifiers load(String contractAddress, Web3j web3j, TransactionManager transactionManager, GasProvider contractGasProvider, Long chainId) {
        return new Modifiers(contractAddress, web3j, transactionManager, contractGasProvider, chainId);
    }
}
