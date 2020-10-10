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
public class Destructable extends Contract {
    private static final String BINARY = "608060405234801561001057600080fd5b50336000806101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff1602179055506101e1806100606000396000f3fe608060405234801561001057600080fd5b506004361061002b5760003560e01c80632b68b9c614610030575b600080fd5b61003861003a565b005b6000809054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff16146100fc576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040180806020018281038252600a8152602001807f6f6e6c79206f776e65720000000000000000000000000000000000000000000081525060200191505060405180910390fd5b60004714610172576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004018080602001828103825260138152602001807f62616c616e6365206973206e6f74207a65726f0000000000000000000000000081525060200191505060405180910390fd5b6000809054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16fffea265627a7a72315820c7b9956175fa7d77db29e4908a20107c0185fe2abc4645efa1f91088f30e23ae64736f6c63430005110032";

    public static final String FUNC_DESTRUCT = "destruct";

    protected Destructable(String contractAddress, Web3j web3j, Credentials credentials, GasProvider contractGasProvider, Long chainId) {
        super(BINARY, contractAddress, web3j, credentials, contractGasProvider, chainId);
    }

    protected Destructable(String contractAddress, Web3j web3j, TransactionManager transactionManager, GasProvider contractGasProvider, Long chainId) {
        super(BINARY, contractAddress, web3j, transactionManager, contractGasProvider, chainId);
    }

    public static RemoteCall<Destructable> deploy(Web3j web3j, Credentials credentials, GasProvider contractGasProvider, Long chainId) {
        return deployRemoteCall(Destructable.class, web3j, credentials, contractGasProvider, BINARY,  "", chainId);
    }

    public static RemoteCall<Destructable> deploy(Web3j web3j, TransactionManager transactionManager, GasProvider contractGasProvider, Long chainId) {
        return deployRemoteCall(Destructable.class, web3j, transactionManager, contractGasProvider, BINARY,  "", chainId);
    }

    public RemoteCall<TransactionReceipt> destruct() {
        final Function function = new Function(
                FUNC_DESTRUCT, 
                Arrays.<Type>asList(), 
                Collections.<TypeReference<?>>emptyList());
        return executeRemoteCallTransaction(function);
    }

    public static Destructable load(String contractAddress, Web3j web3j, Credentials credentials, GasProvider contractGasProvider, Long chainId) {
        return new Destructable(contractAddress, web3j, credentials, contractGasProvider, chainId);
    }

    public static Destructable load(String contractAddress, Web3j web3j, TransactionManager transactionManager, GasProvider contractGasProvider, Long chainId) {
        return new Destructable(contractAddress, web3j, transactionManager, contractGasProvider, chainId);
    }
}
