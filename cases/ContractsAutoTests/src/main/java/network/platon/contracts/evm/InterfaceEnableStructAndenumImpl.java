package network.platon.contracts.evm;

import java.math.BigInteger;
import java.util.Arrays;
import org.web3j.abi.TypeReference;
import org.web3j.abi.datatypes.Function;
import org.web3j.abi.datatypes.Type;
import org.web3j.abi.datatypes.generated.Uint8;
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
public class InterfaceEnableStructAndenumImpl extends Contract {
    private static final String BINARY = "608060405234801561001057600080fd5b5060f28061001f6000396000f3fe608060405260043610603f576000357c0100000000000000000000000000000000000000000000000000000000900463ffffffff1680636ff65603146044575b600080fd5b348015604f57600080fd5b5060566079565b60405180826001811115606557fe5b60ff16815260200191505060405180910390f35b6000608160a8565b602060405190810160405280600180811115609857fe5b8152509050806000015191505090565b6020604051908101604052806000600181111560c057fe5b8152509056fea165627a7a72305820278f38a7d78f49bf527afafd5ab2833172e4df97408e85c0c71581df56c566150029";

    public static final String FUNC_GETPRODUCTCONDITION = "getProductCondition";

    protected InterfaceEnableStructAndenumImpl(String contractAddress, Web3j web3j, Credentials credentials, GasProvider contractGasProvider, Long chainId) {
        super(BINARY, contractAddress, web3j, credentials, contractGasProvider, chainId);
    }

    protected InterfaceEnableStructAndenumImpl(String contractAddress, Web3j web3j, TransactionManager transactionManager, GasProvider contractGasProvider, Long chainId) {
        super(BINARY, contractAddress, web3j, transactionManager, contractGasProvider, chainId);
    }

    public RemoteCall<BigInteger> getProductCondition() {
        final Function function = new Function(FUNC_GETPRODUCTCONDITION, 
                Arrays.<Type>asList(), 
                Arrays.<TypeReference<?>>asList(new TypeReference<Uint8>() {}));
        return executeRemoteCallSingleValueReturn(function, BigInteger.class);
    }

    public static RemoteCall<InterfaceEnableStructAndenumImpl> deploy(Web3j web3j, Credentials credentials, GasProvider contractGasProvider, Long chainId) {
        return deployRemoteCall(InterfaceEnableStructAndenumImpl.class, web3j, credentials, contractGasProvider, BINARY,  "", chainId);
    }

    public static RemoteCall<InterfaceEnableStructAndenumImpl> deploy(Web3j web3j, TransactionManager transactionManager, GasProvider contractGasProvider, Long chainId) {
        return deployRemoteCall(InterfaceEnableStructAndenumImpl.class, web3j, transactionManager, contractGasProvider, BINARY,  "", chainId);
    }

    public static InterfaceEnableStructAndenumImpl load(String contractAddress, Web3j web3j, Credentials credentials, GasProvider contractGasProvider, Long chainId) {
        return new InterfaceEnableStructAndenumImpl(contractAddress, web3j, credentials, contractGasProvider, chainId);
    }

    public static InterfaceEnableStructAndenumImpl load(String contractAddress, Web3j web3j, TransactionManager transactionManager, GasProvider contractGasProvider, Long chainId) {
        return new InterfaceEnableStructAndenumImpl(contractAddress, web3j, transactionManager, contractGasProvider, chainId);
    }
}
