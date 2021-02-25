package network.platon.contracts.evm;

import java.math.BigInteger;
import java.util.Arrays;
import java.util.Collections;
import org.web3j.abi.TypeReference;
import org.web3j.abi.datatypes.Function;
import org.web3j.abi.datatypes.Type;
import org.web3j.abi.datatypes.generated.Uint256;
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
 * <p>Please use the <a href="https://github.com/PlatONnetwork/client-sdk-java/releases">platon-web3j command line tools</a>,
 * or the org.web3j.codegen.SolidityFunctionWrapperGenerator in the 
 * <a href="https://github.com/PlatONnetwork/client-sdk-java/tree/master/codegen">codegen module</a> to update.
 *
 * <p>Generated with web3j version 0.13.1.5.
 */
public class ConstantViewPure extends Contract {
    private static final String BINARY = "608060405234801561001057600080fd5b5061029f806100206000396000f30060806040526004361061006d576000357c0100000000000000000000000000000000000000000000000000000000900463ffffffff168063262a9dff146100725780632671a08b1461009d5780635df79a41146100b4578063a1b1090e146100df578063ec36d67b1461010a575b600080fd5b34801561007e57600080fd5b50610087610135565b6040518082815260200191505060405180910390f35b3480156100a957600080fd5b506100b261013b565b005b3480156100c057600080fd5b506100c9610191565b6040518082815260200191505060405180910390f35b3480156100eb57600080fd5b506100f461019a565b6040518082815260200191505060405180910390f35b34801561011657600080fd5b5061011f6101b4565b6040518082815260200191505060405180910390f35b60015481565b6040805190810160405280600781526020017f66616e7869616e00000000000000000000000000000000000000000000000000815250600090805190602001906101869291906101ce565b506013600181905550565b60006001905090565b600060018060008282540192505081905550600154905090565b600060018060008282540192505081905550600154905090565b828054600181600116156101000203166002900490600052602060002090601f016020900481019282601f1061020f57805160ff191683800117855561023d565b8280016001018555821561023d579182015b8281111561023c578251825591602001919060010190610221565b5b50905061024a919061024e565b5090565b61027091905b8082111561026c576000816000905550600101610254565b5090565b905600a165627a7a723058201c7b01b08e2ab2008a82d3737192f21229d3f36abef2d7c4958edcbdb991d2390029";

    public static final String FUNC_AGE = "age";

    public static final String FUNC_CONSTANTVIEWPURE = "constantViewPure";

    public static final String FUNC_GETAGEBYPURE = "getAgeByPure";

    public static final String FUNC_GETAGEBYVIEW = "getAgeByView";

    public static final String FUNC_GETAGEBYCONSTANT = "getAgeByConstant";

    protected ConstantViewPure(String contractAddress, Web3j web3j, Credentials credentials, GasProvider contractGasProvider, Long chainId) {
        super(BINARY, contractAddress, web3j, credentials, contractGasProvider, chainId);
    }

    protected ConstantViewPure(String contractAddress, Web3j web3j, TransactionManager transactionManager, GasProvider contractGasProvider, Long chainId) {
        super(BINARY, contractAddress, web3j, transactionManager, contractGasProvider, chainId);
    }

    public RemoteCall<BigInteger> age() {
        final Function function = new Function(FUNC_AGE, 
                Arrays.<Type>asList(), 
                Arrays.<TypeReference<?>>asList(new TypeReference<Uint256>() {}));
        return executeRemoteCallSingleValueReturn(function, BigInteger.class);
    }

    public RemoteCall<TransactionReceipt> constantViewPure() {
        final Function function = new Function(
                FUNC_CONSTANTVIEWPURE, 
                Arrays.<Type>asList(), 
                Collections.<TypeReference<?>>emptyList());
        return executeRemoteCallTransaction(function);
    }

    public RemoteCall<BigInteger> getAgeByPure() {
        final Function function = new Function(FUNC_GETAGEBYPURE, 
                Arrays.<Type>asList(), 
                Arrays.<TypeReference<?>>asList(new TypeReference<Uint256>() {}));
        return executeRemoteCallSingleValueReturn(function, BigInteger.class);
    }

    public RemoteCall<BigInteger> getAgeByView() {
        final Function function = new Function(FUNC_GETAGEBYVIEW, 
                Arrays.<Type>asList(), 
                Arrays.<TypeReference<?>>asList(new TypeReference<Uint256>() {}));
        return executeRemoteCallSingleValueReturn(function, BigInteger.class);
    }

    public RemoteCall<BigInteger> getAgeByConstant() {
        final Function function = new Function(FUNC_GETAGEBYCONSTANT, 
                Arrays.<Type>asList(), 
                Arrays.<TypeReference<?>>asList(new TypeReference<Uint256>() {}));
        return executeRemoteCallSingleValueReturn(function, BigInteger.class);
    }

    public static RemoteCall<ConstantViewPure> deploy(Web3j web3j, Credentials credentials, GasProvider contractGasProvider, Long chainId) {
        return deployRemoteCall(ConstantViewPure.class, web3j, credentials, contractGasProvider, BINARY,  "", chainId);
    }

    public static RemoteCall<ConstantViewPure> deploy(Web3j web3j, TransactionManager transactionManager, GasProvider contractGasProvider, Long chainId) {
        return deployRemoteCall(ConstantViewPure.class, web3j, transactionManager, contractGasProvider, BINARY,  "", chainId);
    }

    public static ConstantViewPure load(String contractAddress, Web3j web3j, Credentials credentials, GasProvider contractGasProvider, Long chainId) {
        return new ConstantViewPure(contractAddress, web3j, credentials, contractGasProvider, chainId);
    }

    public static ConstantViewPure load(String contractAddress, Web3j web3j, TransactionManager transactionManager, GasProvider contractGasProvider, Long chainId) {
        return new ConstantViewPure(contractAddress, web3j, transactionManager, contractGasProvider, chainId);
    }
}
