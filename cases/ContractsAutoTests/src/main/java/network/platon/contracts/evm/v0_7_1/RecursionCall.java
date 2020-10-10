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
public class RecursionCall extends Contract {
    private static final String BINARY = "608060405234801561001057600080fd5b5060fe8061001f6000396000f3fe60806040526004361060265760003560e01c8063191a62d414602b57806357e98139146053575b600080fd5b348015603657600080fd5b50603d6092565b6040518082815260200191505060405180910390f35b607c60048036036020811015606757600080fd5b8101908080359060200190929190505050609b565b6040518082815260200191505060405180910390f35b60008054905090565b600081600054101560be5760008081546001019190508190555060bc82609b565b505b600054905091905056fea26469706673582212207611d868aaac2ce07e2bf51b440a07085b595e6e2ebc9efb39dd4c746f9421d464736f6c63430007010033";

    public static final String FUNC_GET_TOTAL = "get_total";

    public static final String FUNC_RECURSIONCALLTEST = "recursionCallTest";

    protected RecursionCall(String contractAddress, Web3j web3j, Credentials credentials, GasProvider contractGasProvider, Long chainId) {
        super(BINARY, contractAddress, web3j, credentials, contractGasProvider, chainId);
    }

    protected RecursionCall(String contractAddress, Web3j web3j, TransactionManager transactionManager, GasProvider contractGasProvider, Long chainId) {
        super(BINARY, contractAddress, web3j, transactionManager, contractGasProvider, chainId);
    }

    public RemoteCall<TransactionReceipt> get_total() {
        final Function function = new Function(
                FUNC_GET_TOTAL, 
                Arrays.<Type>asList(), 
                Collections.<TypeReference<?>>emptyList());
        return executeRemoteCallTransaction(function);
    }

    public RemoteCall<TransactionReceipt> recursionCallTest(BigInteger n) {
        final Function function = new Function(
                FUNC_RECURSIONCALLTEST, 
                Arrays.<Type>asList(new com.alaya.abi.solidity.datatypes.generated.Uint256(n)), 
                Collections.<TypeReference<?>>emptyList());
        return executeRemoteCallTransaction(function);
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
