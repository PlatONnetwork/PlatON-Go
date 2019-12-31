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
 * <p>Please use the <a href="https://docs.web3j.io/command_line.html">web3j command line tools</a>,
 * or the org.web3j.codegen.SolidityFunctionWrapperGenerator in the 
 * <a href="https://github.com/web3j/web3j/tree/master/codegen">codegen module</a> to update.
 *
 * <p>Generated with web3j version 0.7.5.0.
 */
public class DisallowVariable extends Contract {
    private static final String BINARY = "6080604052600160025534801561001557600080fd5b5060db806100246000396000f3fe608060405260043610603f576000357c0100000000000000000000000000000000000000000000000000000000900463ffffffff1680630f2da424146044575b600080fd5b606d60048036036020811015605857600080fd5b81019080803590602001909291905050506083565b6040518082815260200191505060405180910390f35b60008060006002600391509150600060016000868152602001908152602001600020905050505091905056fea165627a7a72305820d7483d7eebf4a3d732e3a32b298cf0c4e38fa5390d2260de964ebd1ba08434810029";

    public static final String FUNC_TESEMPTY = "tesEmpty";

    @Deprecated
    protected DisallowVariable(String contractAddress, Web3j web3j, Credentials credentials, BigInteger gasPrice, BigInteger gasLimit) {
        super(BINARY, contractAddress, web3j, credentials, gasPrice, gasLimit);
    }

    protected DisallowVariable(String contractAddress, Web3j web3j, Credentials credentials, GasProvider contractGasProvider) {
        super(BINARY, contractAddress, web3j, credentials, contractGasProvider);
    }

    @Deprecated
    protected DisallowVariable(String contractAddress, Web3j web3j, TransactionManager transactionManager, BigInteger gasPrice, BigInteger gasLimit) {
        super(BINARY, contractAddress, web3j, transactionManager, gasPrice, gasLimit);
    }

    protected DisallowVariable(String contractAddress, Web3j web3j, TransactionManager transactionManager, GasProvider contractGasProvider) {
        super(BINARY, contractAddress, web3j, transactionManager, contractGasProvider);
    }

    public RemoteCall<TransactionReceipt> tesEmpty(BigInteger _id, BigInteger weiValue) {
        final Function function = new Function(
                FUNC_TESEMPTY, 
                Arrays.<Type>asList(new org.web3j.abi.datatypes.generated.Uint256(_id)), 
                Collections.<TypeReference<?>>emptyList());
        return executeRemoteCallTransaction(function, weiValue);
    }

    public static RemoteCall<DisallowVariable> deploy(Web3j web3j, Credentials credentials, GasProvider contractGasProvider) {
        return deployRemoteCall(DisallowVariable.class, web3j, credentials, contractGasProvider, BINARY, "");
    }

    @Deprecated
    public static RemoteCall<DisallowVariable> deploy(Web3j web3j, Credentials credentials, BigInteger gasPrice, BigInteger gasLimit) {
        return deployRemoteCall(DisallowVariable.class, web3j, credentials, gasPrice, gasLimit, BINARY, "");
    }

    public static RemoteCall<DisallowVariable> deploy(Web3j web3j, TransactionManager transactionManager, GasProvider contractGasProvider) {
        return deployRemoteCall(DisallowVariable.class, web3j, transactionManager, contractGasProvider, BINARY, "");
    }

    @Deprecated
    public static RemoteCall<DisallowVariable> deploy(Web3j web3j, TransactionManager transactionManager, BigInteger gasPrice, BigInteger gasLimit) {
        return deployRemoteCall(DisallowVariable.class, web3j, transactionManager, gasPrice, gasLimit, BINARY, "");
    }

    @Deprecated
    public static DisallowVariable load(String contractAddress, Web3j web3j, Credentials credentials, BigInteger gasPrice, BigInteger gasLimit) {
        return new DisallowVariable(contractAddress, web3j, credentials, gasPrice, gasLimit);
    }

    @Deprecated
    public static DisallowVariable load(String contractAddress, Web3j web3j, TransactionManager transactionManager, BigInteger gasPrice, BigInteger gasLimit) {
        return new DisallowVariable(contractAddress, web3j, transactionManager, gasPrice, gasLimit);
    }

    public static DisallowVariable load(String contractAddress, Web3j web3j, Credentials credentials, GasProvider contractGasProvider) {
        return new DisallowVariable(contractAddress, web3j, credentials, contractGasProvider);
    }

    public static DisallowVariable load(String contractAddress, Web3j web3j, TransactionManager transactionManager, GasProvider contractGasProvider) {
        return new DisallowVariable(contractAddress, web3j, transactionManager, contractGasProvider);
    }
}
