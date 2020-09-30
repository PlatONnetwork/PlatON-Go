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
public class InheritContractAMutipleClass extends Contract {
    private static final String BINARY = "608060405234801561001057600080fd5b5060ba8061001f6000396000f3fe6080604052348015600f57600080fd5b506004361060325760003560e01c8063430fe9c11460375780635fcceceb146053575b600080fd5b603d606f565b6040518082815260200191505060405180910390f35b60596078565b6040518082815260200191505060405180910390f35b60006001905090565b60006080606f565b90509056fea265627a7a723158200d8cae82526d56badd6126f66b46fa8745b0b6846d757cecbb372a3378c81b0964736f6c634300050d0032";

    public static final String FUNC_CALLGETDATEA = "callGetDateA";

    public static final String FUNC_GETDATE = "getDate";

    protected InheritContractAMutipleClass(String contractAddress, Web3j web3j, Credentials credentials, GasProvider contractGasProvider, Long chainId) {
        super(BINARY, contractAddress, web3j, credentials, contractGasProvider, chainId);
    }

    protected InheritContractAMutipleClass(String contractAddress, Web3j web3j, TransactionManager transactionManager, GasProvider contractGasProvider, Long chainId) {
        super(BINARY, contractAddress, web3j, transactionManager, contractGasProvider, chainId);
    }

    public RemoteCall<BigInteger> callGetDateA() {
        final Function function = new Function(FUNC_CALLGETDATEA, 
                Arrays.<Type>asList(), 
                Arrays.<TypeReference<?>>asList(new TypeReference<Uint256>() {}));
        return executeRemoteCallSingleValueReturn(function, BigInteger.class);
    }

    public RemoteCall<BigInteger> getDate() {
        final Function function = new Function(FUNC_GETDATE, 
                Arrays.<Type>asList(), 
                Arrays.<TypeReference<?>>asList(new TypeReference<Uint256>() {}));
        return executeRemoteCallSingleValueReturn(function, BigInteger.class);
    }

    public static RemoteCall<InheritContractAMutipleClass> deploy(Web3j web3j, Credentials credentials, GasProvider contractGasProvider, Long chainId) {
        return deployRemoteCall(InheritContractAMutipleClass.class, web3j, credentials, contractGasProvider, BINARY,  "", chainId);
    }

    public static RemoteCall<InheritContractAMutipleClass> deploy(Web3j web3j, TransactionManager transactionManager, GasProvider contractGasProvider, Long chainId) {
        return deployRemoteCall(InheritContractAMutipleClass.class, web3j, transactionManager, contractGasProvider, BINARY,  "", chainId);
    }

    public static InheritContractAMutipleClass load(String contractAddress, Web3j web3j, Credentials credentials, GasProvider contractGasProvider, Long chainId) {
        return new InheritContractAMutipleClass(contractAddress, web3j, credentials, contractGasProvider, chainId);
    }

    public static InheritContractAMutipleClass load(String contractAddress, Web3j web3j, TransactionManager transactionManager, GasProvider contractGasProvider, Long chainId) {
        return new InheritContractAMutipleClass(contractAddress, web3j, transactionManager, contractGasProvider, chainId);
    }
}
