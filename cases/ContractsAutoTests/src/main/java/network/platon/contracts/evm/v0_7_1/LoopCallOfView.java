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
public class LoopCallOfView extends Contract {
    private static final String BINARY = "608060405234801561001057600080fd5b5060cd8061001f6000396000f3fe6080604052348015600f57600080fd5b506004361060285760003560e01c80633fde082714602d575b600080fd5b605660048036036020811015604157600080fd5b8101908080359060200190929190505050606c565b6040518082815260200191505060405180910390f35b60008060005b83811015608d57818060010192505080806001019150506072565b508091505091905056fea2646970667358221220e0a114a52703a54f85a4437cc2e28c40b30b344b01063c2ea2e36b00bb0c1f4964736f6c63430007010033";

    public static final String FUNC_LOOPCALLTEST = "loopCallTest";

    protected LoopCallOfView(String contractAddress, Web3j web3j, Credentials credentials, GasProvider contractGasProvider, Long chainId) {
        super(BINARY, contractAddress, web3j, credentials, contractGasProvider, chainId);
    }

    protected LoopCallOfView(String contractAddress, Web3j web3j, TransactionManager transactionManager, GasProvider contractGasProvider, Long chainId) {
        super(BINARY, contractAddress, web3j, transactionManager, contractGasProvider, chainId);
    }

    public RemoteCall<TransactionReceipt> loopCallTest(BigInteger n) {
        final Function function = new Function(
                FUNC_LOOPCALLTEST, 
                Arrays.<Type>asList(new com.alaya.abi.solidity.datatypes.generated.Uint256(n)), 
                Collections.<TypeReference<?>>emptyList());
        return executeRemoteCallTransaction(function);
    }

    public static RemoteCall<LoopCallOfView> deploy(Web3j web3j, Credentials credentials, GasProvider contractGasProvider, Long chainId) {
        return deployRemoteCall(LoopCallOfView.class, web3j, credentials, contractGasProvider, BINARY,  "", chainId);
    }

    public static RemoteCall<LoopCallOfView> deploy(Web3j web3j, TransactionManager transactionManager, GasProvider contractGasProvider, Long chainId) {
        return deployRemoteCall(LoopCallOfView.class, web3j, transactionManager, contractGasProvider, BINARY,  "", chainId);
    }

    public static LoopCallOfView load(String contractAddress, Web3j web3j, Credentials credentials, GasProvider contractGasProvider, Long chainId) {
        return new LoopCallOfView(contractAddress, web3j, credentials, contractGasProvider, chainId);
    }

    public static LoopCallOfView load(String contractAddress, Web3j web3j, TransactionManager transactionManager, GasProvider contractGasProvider, Long chainId) {
        return new LoopCallOfView(contractAddress, web3j, transactionManager, contractGasProvider, chainId);
    }
}
