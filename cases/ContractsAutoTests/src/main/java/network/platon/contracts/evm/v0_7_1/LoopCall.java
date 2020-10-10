package network.platon.contracts.evm.v0_7_1;

import com.alaya.abi.solidity.TypeReference;
import com.alaya.abi.solidity.datatypes.Function;
import com.alaya.abi.solidity.datatypes.Type;
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
public class LoopCall extends Contract {
    private static final String BINARY = "608060405234801561001057600080fd5b5060b98061001f6000396000f3fe6080604052348015600f57600080fd5b506004361060285760003560e01c80633fde082714602d575b600080fd5b605660048036036020811015604157600080fd5b81019080803590602001909291905050506058565b005b60005b81811015607f5760008081548092919060010191905055508080600101915050605b565b505056fea2646970667358221220b6ed2acf4e9c2f5880d13171c9c1f22afcbe228ad633460cfc59cc850625232a64736f6c63430007010033";

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
                Arrays.<Type>asList(new com.alaya.abi.solidity.datatypes.generated.Uint256(n)), 
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
