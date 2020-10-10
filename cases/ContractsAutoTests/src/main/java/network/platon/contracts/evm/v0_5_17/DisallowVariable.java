package network.platon.contracts.evm.v0_5_17;

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
public class DisallowVariable extends Contract {
    private static final String BINARY = "60806040526001600255348015601457600080fd5b5060c1806100236000396000f3fe608060405260043610601c5760003560e01c80630f2da424146021575b600080fd5b604a60048036036020811015603557600080fd5b81019080803590602001909291905050506060565b6040518082815260200191505060405180910390f35b60008060006002600391509150600060016000868152602001908152602001600020905050505091905056fea265627a7a72315820bc895278a9d5ea9de0e84a168079fccd7f41071677da16f42b406be51e8b247b64736f6c63430005110032";

    public static final String FUNC_TESEMPTY = "tesEmpty";

    protected DisallowVariable(String contractAddress, Web3j web3j, Credentials credentials, GasProvider contractGasProvider, Long chainId) {
        super(BINARY, contractAddress, web3j, credentials, contractGasProvider, chainId);
    }

    protected DisallowVariable(String contractAddress, Web3j web3j, TransactionManager transactionManager, GasProvider contractGasProvider, Long chainId) {
        super(BINARY, contractAddress, web3j, transactionManager, contractGasProvider, chainId);
    }

    public RemoteCall<TransactionReceipt> tesEmpty(BigInteger _id, BigInteger vonValue) {
        final Function function = new Function(
                FUNC_TESEMPTY, 
                Arrays.<Type>asList(new com.alaya.abi.solidity.datatypes.generated.Uint256(_id)), 
                Collections.<TypeReference<?>>emptyList());
        return executeRemoteCallTransaction(function, vonValue);
    }

    public static RemoteCall<DisallowVariable> deploy(Web3j web3j, Credentials credentials, GasProvider contractGasProvider, Long chainId) {
        return deployRemoteCall(DisallowVariable.class, web3j, credentials, contractGasProvider, BINARY,  "", chainId);
    }

    public static RemoteCall<DisallowVariable> deploy(Web3j web3j, TransactionManager transactionManager, GasProvider contractGasProvider, Long chainId) {
        return deployRemoteCall(DisallowVariable.class, web3j, transactionManager, contractGasProvider, BINARY,  "", chainId);
    }

    public static DisallowVariable load(String contractAddress, Web3j web3j, Credentials credentials, GasProvider contractGasProvider, Long chainId) {
        return new DisallowVariable(contractAddress, web3j, credentials, contractGasProvider, chainId);
    }

    public static DisallowVariable load(String contractAddress, Web3j web3j, TransactionManager transactionManager, GasProvider contractGasProvider, Long chainId) {
        return new DisallowVariable(contractAddress, web3j, transactionManager, contractGasProvider, chainId);
    }
}
