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
public class RecursionCall extends Contract {
    private static final String BINARY = "608060405234801561001057600080fd5b5060fd8061001f6000396000f3fe60806040526004361060265760003560e01c8063191a62d414602b57806357e98139146053575b600080fd5b348015603657600080fd5b50603d6092565b6040518082815260200191505060405180910390f35b607c60048036036020811015606757600080fd5b8101908080359060200190929190505050609b565b6040518082815260200191505060405180910390f35b60008054905090565b600081600054101560be5760008081546001019190508190555060bc82609b565b505b600054905091905056fea265627a7a7231582000d7a1173bda032f4e44dd2bd0c65a9d985c9086547f31faf22ae36c0b3e091864736f6c634300050d0032";

    public static final String FUNC_GET_TOTAL = "get_total";

    public static final String FUNC_RECURSIONCALLTEST = "recursionCallTest";

    protected RecursionCall(String contractAddress, Web3j web3j, Credentials credentials, GasProvider contractGasProvider, Long chainId) {
        super(BINARY, contractAddress, web3j, credentials, contractGasProvider, chainId);
    }

    protected RecursionCall(String contractAddress, Web3j web3j, TransactionManager transactionManager, GasProvider contractGasProvider, Long chainId) {
        super(BINARY, contractAddress, web3j, transactionManager, contractGasProvider, chainId);
    }

    public RemoteCall<BigInteger> get_total() {
        final Function function = new Function(FUNC_GET_TOTAL, 
                Arrays.<Type>asList(), 
                Arrays.<TypeReference<?>>asList(new TypeReference<Uint256>() {}));
        return executeRemoteCallSingleValueReturn(function, BigInteger.class);
    }

    public RemoteCall<TransactionReceipt> recursionCallTest(BigInteger n, BigInteger vonValue) {
        final Function function = new Function(
                FUNC_RECURSIONCALLTEST, 
                Arrays.<Type>asList(new Uint256(n)),
                Collections.<TypeReference<?>>emptyList());
        return executeRemoteCallTransaction(function, vonValue);
    }

    public static RemoteCall<RecursionCall> deploy(Web3j web3j, Credentials credentials, GasProvider contractGasProvider, Long chainId) {
        return deployRemoteCall(RecursionCall.class, web3j, credentials, contractGasProvider, BINARY,  "", chainId);
    }

    public static RemoteCall<RecursionCall> deploy(Web3j web3j, TransactionManager transactionManager, GasProvider contractGasProvider, Long chainId) {
        return deployRemoteCall(RecursionCall.class, web3j, transactionManager, contractGasProvider, BINARY,  "", chainId);
    }

    public static RecursionCall load(String contractAddress, Web3j web3j, Credentials credentials, GasProvider contractGasProvider, Long chainId) {
        return new RecursionCall(contractAddress, web3j, credentials, contractGasProvider, chainId);
    }

    public static RecursionCall load(String contractAddress, Web3j web3j, TransactionManager transactionManager, GasProvider contractGasProvider, Long chainId) {
        return new RecursionCall(contractAddress, web3j, transactionManager, contractGasProvider, chainId);
    }
}
