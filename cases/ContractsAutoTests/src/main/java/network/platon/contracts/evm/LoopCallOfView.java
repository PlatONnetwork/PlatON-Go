package network.platon.contracts.evm;

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
 * <p>Please use the <a href="https://github.com/PlatONnetwork/client-sdk-java/releases">platon-web3j command line tools</a>,
 * or the org.web3j.codegen.SolidityFunctionWrapperGenerator in the 
 * <a href="https://github.com/PlatONnetwork/client-sdk-java/tree/master/codegen">codegen module</a> to update.
 *
 * <p>Generated with web3j version 0.13.1.5.
 */
public class LoopCallOfView extends Contract {
    private static final String BINARY = "608060405234801561001057600080fd5b5060cf8061001f6000396000f3fe6080604052348015600f57600080fd5b506004361060285760003560e01c80633fde082714602d575b600080fd5b605660048036036020811015604157600080fd5b8101908080359060200190929190505050606c565b6040518082815260200191505060405180910390f35b60008060008090505b83811015609057818060010192505080806001019150506075565b508091505091905056fea265627a7a72315820beefe9ff7cd6e1d892f11e7026a708deb3b9ece9e77178de5d16039ddb76d65d64736f6c634300050d0032";

    public static final String FUNC_LOOPCALLTEST = "loopCallTest";

    protected LoopCallOfView(String contractAddress, Web3j web3j, Credentials credentials, GasProvider contractGasProvider, Long chainId) {
        super(BINARY, contractAddress, web3j, credentials, contractGasProvider, chainId);
    }

    protected LoopCallOfView(String contractAddress, Web3j web3j, TransactionManager transactionManager, GasProvider contractGasProvider, Long chainId) {
        super(BINARY, contractAddress, web3j, transactionManager, contractGasProvider, chainId);
    }

    public RemoteCall<BigInteger> loopCallTest(BigInteger n) {
        final Function function = new Function(FUNC_LOOPCALLTEST, 
                Arrays.<Type>asList(new Uint256(n)),
                Arrays.<TypeReference<?>>asList(new TypeReference<Uint256>() {}));
        return executeRemoteCallSingleValueReturn(function, BigInteger.class);
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
