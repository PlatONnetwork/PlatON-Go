package network.platon.contracts.evm.v0_5_17;

import com.alaya.abi.solidity.TypeReference;
import com.alaya.abi.solidity.datatypes.Function;
import com.alaya.abi.solidity.datatypes.Type;
import com.alaya.abi.solidity.datatypes.generated.Uint256;
import com.alaya.crypto.Credentials;
import com.alaya.protocol.Web3j;
import com.alaya.protocol.core.RemoteCall;
import com.alaya.protocol.core.methods.response.TransactionReceipt;
import com.alaya.tx.Contract;
import com.alaya.tx.TransactionManager;
import com.alaya.tx.gas.GasProvider;
import java.math.BigInteger;
import java.util.Arrays;
import java.util.Collections;

/**
 * <p>Auto generated code.
 * <p><strong>Do not modify!</strong>
 * <p>Please use the <a href="https://github.com/PlatONnetwork/client-sdk-java/releases">platon-web3j command line tools</a>,
 * or the com.alaya.codegen.SolidityFunctionWrapperGenerator in the 
 * <a href="https://github.com/PlatONnetwork/client-sdk-java/tree/master/codegen">codegen module</a> to update.
 *
 * <p>Generated with web3j version 0.13.2.0.
 */
public class Modifiers extends Contract {
    private static final String BINARY = "6080604052600a600055348015601457600080fd5b5060c1806100236000396000f3fe6080604052348015600f57600080fd5b506004361060325760003560e01c806366e41cb71460375780636b59084d146053575b600080fd5b603d605b565b6040518082815260200191505060405180910390f35b60596064565b005b60008054905090565b6000546000819050600080549050600c600081905550506000549050600b600081905550505056fea265627a7a72315820576bc351eea344908e41e531502a56d04f1917b0d89a793e1f80c83f1261278764736f6c63430005110032";

    public static final String FUNC_TEST1 = "test1";

    public static final String FUNC_TEST2 = "test2";

    protected Modifiers(String contractAddress, Web3j web3j, Credentials credentials, GasProvider contractGasProvider, Long chainId) {
        super(BINARY, contractAddress, web3j, credentials, contractGasProvider, chainId);
    }

    protected Modifiers(String contractAddress, Web3j web3j, TransactionManager transactionManager, GasProvider contractGasProvider, Long chainId) {
        super(BINARY, contractAddress, web3j, transactionManager, contractGasProvider, chainId);
    }

    public RemoteCall<TransactionReceipt> test1() {
        final Function function = new Function(
                FUNC_TEST1, 
                Arrays.<Type>asList(), 
                Collections.<TypeReference<?>>emptyList());
        return executeRemoteCallTransaction(function);
    }

    public RemoteCall<BigInteger> test2() {
        final Function function = new Function(FUNC_TEST2, 
                Arrays.<Type>asList(), 
                Arrays.<TypeReference<?>>asList(new TypeReference<Uint256>() {}));
        return executeRemoteCallSingleValueReturn(function, BigInteger.class);
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
