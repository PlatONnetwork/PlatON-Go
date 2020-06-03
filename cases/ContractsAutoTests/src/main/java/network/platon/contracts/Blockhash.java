package network.platon.contracts;

import java.util.Arrays;
import org.web3j.abi.TypeReference;
import org.web3j.abi.datatypes.Function;
import org.web3j.abi.datatypes.Type;
import org.web3j.abi.datatypes.generated.Bytes32;
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
 * <p>Generated with web3j version 0.13.0.7.
 */
public class Blockhash extends Contract {
    private static final String BINARY = "608060405234801561001057600080fd5b5061017b806100206000396000f3fe608060405234801561001057600080fd5b50600436106100575760003560e01c80630f7536281461005c57806366b3eb341461007a578063696d67c3146100985780639e1f194e146100b6578063e1b99d74146100d4575b600080fd5b6100646100f2565b6040518082815260200191505060405180910390f35b610082610103565b6040518082815260200191505060405180910390f35b6100a0610115565b6040518082815260200191505060405180910390f35b6100be610127565b6040518082815260200191505060405180910390f35b6100dc610138565b6040518082815260200191505060405180910390f35b60008060ff43034090508091505090565b60008061010043034090508091505090565b60008061010143034090508091505090565b600080601e43034090508091505090565b60008043409050809150509056fea265627a7a72315820ccdc7f58d7930bf24b80831c1b34a129b2abf1d38e9641e7d04914120948aeb164736f6c634300050d0032";

    public static final String FUNC_GETBLOCKHASHBEFORE0 = "getBlockhashbefore0";

    public static final String FUNC_GETBLOCKHASHBEFORE255 = "getBlockhashbefore255";

    public static final String FUNC_GETBLOCKHASHBEFORE256 = "getBlockhashbefore256";

    public static final String FUNC_GETBLOCKHASHBEFORE257 = "getBlockhashbefore257";

    public static final String FUNC_GETBLOCKHASHBEFORE30 = "getBlockhashbefore30";

    protected Blockhash(String contractAddress, Web3j web3j, Credentials credentials, GasProvider contractGasProvider, Long chainId) {
        super(BINARY, contractAddress, web3j, credentials, contractGasProvider, chainId);
    }

    protected Blockhash(String contractAddress, Web3j web3j, TransactionManager transactionManager, GasProvider contractGasProvider, Long chainId) {
        super(BINARY, contractAddress, web3j, transactionManager, contractGasProvider, chainId);
    }

    public RemoteCall<byte[]> getBlockhashbefore0() {
        final Function function = new Function(FUNC_GETBLOCKHASHBEFORE0, 
                Arrays.<Type>asList(), 
                Arrays.<TypeReference<?>>asList(new TypeReference<Bytes32>() {}));
        return executeRemoteCallSingleValueReturn(function, byte[].class);
    }

    public RemoteCall<byte[]> getBlockhashbefore255() {
        final Function function = new Function(FUNC_GETBLOCKHASHBEFORE255, 
                Arrays.<Type>asList(), 
                Arrays.<TypeReference<?>>asList(new TypeReference<Bytes32>() {}));
        return executeRemoteCallSingleValueReturn(function, byte[].class);
    }

    public RemoteCall<byte[]> getBlockhashbefore256() {
        final Function function = new Function(FUNC_GETBLOCKHASHBEFORE256, 
                Arrays.<Type>asList(), 
                Arrays.<TypeReference<?>>asList(new TypeReference<Bytes32>() {}));
        return executeRemoteCallSingleValueReturn(function, byte[].class);
    }

    public RemoteCall<byte[]> getBlockhashbefore257() {
        final Function function = new Function(FUNC_GETBLOCKHASHBEFORE257, 
                Arrays.<Type>asList(), 
                Arrays.<TypeReference<?>>asList(new TypeReference<Bytes32>() {}));
        return executeRemoteCallSingleValueReturn(function, byte[].class);
    }

    public RemoteCall<byte[]> getBlockhashbefore30() {
        final Function function = new Function(FUNC_GETBLOCKHASHBEFORE30, 
                Arrays.<Type>asList(), 
                Arrays.<TypeReference<?>>asList(new TypeReference<Bytes32>() {}));
        return executeRemoteCallSingleValueReturn(function, byte[].class);
    }

    public static RemoteCall<Blockhash> deploy(Web3j web3j, Credentials credentials, GasProvider contractGasProvider, Long chainId) {
        return deployRemoteCall(Blockhash.class, web3j, credentials, contractGasProvider, BINARY,  "", chainId);
    }

    public static RemoteCall<Blockhash> deploy(Web3j web3j, TransactionManager transactionManager, GasProvider contractGasProvider, Long chainId) {
        return deployRemoteCall(Blockhash.class, web3j, transactionManager, contractGasProvider, BINARY,  "", chainId);
    }

    public static Blockhash load(String contractAddress, Web3j web3j, Credentials credentials, GasProvider contractGasProvider, Long chainId) {
        return new Blockhash(contractAddress, web3j, credentials, contractGasProvider, chainId);
    }

    public static Blockhash load(String contractAddress, Web3j web3j, TransactionManager transactionManager, GasProvider contractGasProvider, Long chainId) {
        return new Blockhash(contractAddress, web3j, transactionManager, contractGasProvider, chainId);
    }
}
