package network.platon.contracts;

import java.math.BigInteger;
import java.util.Arrays;
import java.util.Collections;
import org.web3j.abi.TypeReference;
import org.web3j.abi.datatypes.Function;
import org.web3j.abi.datatypes.Type;
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
 * <p>Generated with web3j version 0.13.0.7.
 */
public class LoopCall extends Contract {
    private static final String BINARY = "608060405234801561001057600080fd5b5060bb8061001f6000396000f3fe6080604052348015600f57600080fd5b506004361060285760003560e01c80633fde082714602d575b600080fd5b605660048036036020811015604157600080fd5b81019080803590602001909291905050506058565b005b60008090505b8181101560825760008081548092919060010191905055508080600101915050605e565b505056fea265627a7a723158203cc1339ca8c387514ae67f752b0b7aa2da87f12c7b245c0cc8ba2faa17402cdb64736f6c634300050d0032";

    public static final String FUNC_LOOPCALLTEST = "loopCallTest";

    protected LoopCall(String contractAddress, Web3j web3j, Credentials credentials, GasProvider contractGasProvider, Long chainId) {
        super(BINARY, contractAddress, web3j, credentials, contractGasProvider, chainId);
    }

    protected LoopCall(String contractAddress, Web3j web3j, TransactionManager transactionManager, GasProvider contractGasProvider, Long chainId) {
        super(BINARY, contractAddress, web3j, transactionManager, contractGasProvider, chainId);
    }

    public RemoteCall<TransactionReceipt> loopCallTest(BigInteger n) {
        final Function function = new Function(
                FUNC_LOOPCALLTEST, 
                Arrays.<Type>asList(new org.web3j.abi.datatypes.generated.Uint256(n)), 
                Collections.<TypeReference<?>>emptyList());
        return executeRemoteCallTransaction(function);
    }

    public static RemoteCall<LoopCall> deploy(Web3j web3j, Credentials credentials, GasProvider contractGasProvider, Long chainId) {
        return deployRemoteCall(LoopCall.class, web3j, credentials, contractGasProvider, BINARY,  "", chainId);
    }

    public static RemoteCall<LoopCall> deploy(Web3j web3j, TransactionManager transactionManager, GasProvider contractGasProvider, Long chainId) {
        return deployRemoteCall(LoopCall.class, web3j, transactionManager, contractGasProvider, BINARY,  "", chainId);
    }

    public static LoopCall load(String contractAddress, Web3j web3j, Credentials credentials, GasProvider contractGasProvider, Long chainId) {
        return new LoopCall(contractAddress, web3j, credentials, contractGasProvider, chainId);
    }

    public static LoopCall load(String contractAddress, Web3j web3j, TransactionManager transactionManager, GasProvider contractGasProvider, Long chainId) {
        return new LoopCall(contractAddress, web3j, transactionManager, contractGasProvider, chainId);
    }
}
