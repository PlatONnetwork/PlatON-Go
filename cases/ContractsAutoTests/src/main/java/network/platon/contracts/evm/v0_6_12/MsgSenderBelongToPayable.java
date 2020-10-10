package network.platon.contracts.evm.v0_6_12;

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
public class MsgSenderBelongToPayable extends Contract {
    private static final String BINARY = "608060405234801561001057600080fd5b5060e68061001f6000396000f3fe6080604052348015600f57600080fd5b506004361060325760003560e01c8063485a584b1460375780639e455939146069575b600080fd5b603d609b565b604051808273ffffffffffffffffffffffffffffffffffffffff16815260200191505060405180910390f35b606f60a3565b604051808273ffffffffffffffffffffffffffffffffffffffff16815260200191505060405180910390f35b600033905090565b600080339050309150509056fea2646970667358221220af2f485297320242fe832df13fe6661f3860c4e250ff46536b71430219d52abf64736f6c634300060c0033";

    public static final String FUNC_GETCONTRACTADDR = "getContractAddr";

    public static final String FUNC_GETMSGSENDERADDR = "getMsgSenderAddr";

    protected MsgSenderBelongToPayable(String contractAddress, Web3j web3j, Credentials credentials, GasProvider contractGasProvider, Long chainId) {
        super(BINARY, contractAddress, web3j, credentials, contractGasProvider, chainId);
    }

    protected MsgSenderBelongToPayable(String contractAddress, Web3j web3j, TransactionManager transactionManager, GasProvider contractGasProvider, Long chainId) {
        super(BINARY, contractAddress, web3j, transactionManager, contractGasProvider, chainId);
    }

    public RemoteCall<TransactionReceipt> getContractAddr() {
        final Function function = new Function(
                FUNC_GETCONTRACTADDR, 
                Arrays.<Type>asList(), 
                Collections.<TypeReference<?>>emptyList());
        return executeRemoteCallTransaction(function);
    }

    public RemoteCall<TransactionReceipt> getMsgSenderAddr() {
        final Function function = new Function(
                FUNC_GETMSGSENDERADDR, 
                Arrays.<Type>asList(), 
                Collections.<TypeReference<?>>emptyList());
        return executeRemoteCallTransaction(function);
    }

    public static RemoteCall<MsgSenderBelongToPayable> deploy(Web3j web3j, Credentials credentials, GasProvider contractGasProvider, Long chainId) {
        return deployRemoteCall(MsgSenderBelongToPayable.class, web3j, credentials, contractGasProvider, BINARY,  "", chainId);
    }

    public static RemoteCall<MsgSenderBelongToPayable> deploy(Web3j web3j, TransactionManager transactionManager, GasProvider contractGasProvider, Long chainId) {
        return deployRemoteCall(MsgSenderBelongToPayable.class, web3j, transactionManager, contractGasProvider, BINARY,  "", chainId);
    }

    public static MsgSenderBelongToPayable load(String contractAddress, Web3j web3j, Credentials credentials, GasProvider contractGasProvider, Long chainId) {
        return new MsgSenderBelongToPayable(contractAddress, web3j, credentials, contractGasProvider, chainId);
    }

    public static MsgSenderBelongToPayable load(String contractAddress, Web3j web3j, TransactionManager transactionManager, GasProvider contractGasProvider, Long chainId) {
        return new MsgSenderBelongToPayable(contractAddress, web3j, transactionManager, contractGasProvider, chainId);
    }
}
