package network.platon.contracts.evm;

import java.util.Arrays;
import org.web3j.abi.TypeReference;
import org.web3j.abi.datatypes.Address;
import org.web3j.abi.datatypes.Function;
import org.web3j.abi.datatypes.Type;
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
public class MsgSenderBelongToPayable extends Contract {
    private static final String BINARY = "608060405234801561001057600080fd5b50610140806100206000396000f3fe60806040526004361061004c576000357c0100000000000000000000000000000000000000000000000000000000900463ffffffff168063485a584b146100515780639e455939146100a8575b600080fd5b34801561005d57600080fd5b506100666100ff565b604051808273ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200191505060405180910390f35b3480156100b457600080fd5b506100bd610107565b604051808273ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200191505060405180910390f35b600033905090565b600080339050309150509056fea165627a7a72305820d02efef88f5ba9d37ea7d84ae8525e6d232ed59c12b66d0e571ff71b814f30820029";

    public static final String FUNC_GETMSGSENDERADDR = "getMsgSenderAddr";

    public static final String FUNC_GETCONTRACTADDR = "getContractAddr";

    protected MsgSenderBelongToPayable(String contractAddress, Web3j web3j, Credentials credentials, GasProvider contractGasProvider, Long chainId) {
        super(BINARY, contractAddress, web3j, credentials, contractGasProvider, chainId);
    }

    protected MsgSenderBelongToPayable(String contractAddress, Web3j web3j, TransactionManager transactionManager, GasProvider contractGasProvider, Long chainId) {
        super(BINARY, contractAddress, web3j, transactionManager, contractGasProvider, chainId);
    }

    public RemoteCall<String> getMsgSenderAddr() {
        final Function function = new Function(FUNC_GETMSGSENDERADDR, 
                Arrays.<Type>asList(), 
                Arrays.<TypeReference<?>>asList(new TypeReference<Address>() {}));
        return executeRemoteCallSingleValueReturn(function, String.class);
    }

    public RemoteCall<String> getContractAddr() {
        final Function function = new Function(FUNC_GETCONTRACTADDR, 
                Arrays.<Type>asList(), 
                Arrays.<TypeReference<?>>asList(new TypeReference<Address>() {}));
        return executeRemoteCallSingleValueReturn(function, String.class);
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
