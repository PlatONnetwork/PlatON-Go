package network.platon.contracts;

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
 * <p>Generated with web3j version 0.13.0.7.
 */
public class InheritContractOverloadChild extends Contract {
    private static final String BINARY = "608060405234801561001057600080fd5b50610310806100206000396000f3fe608060405234801561001057600080fd5b506004361061009e5760003560e01c8063a56dfe4a11610066578063a56dfe4a14610163578063a5843f0814610181578063b7b0422d146101b9578063cedf673f146101e7578063fa98b8671461021f5761009e565b80630b7f1665146100a35780630c55699c146100c15780631c5e6b98146100df5780635197c7aa1461010d578063560512c61461012b575b600080fd5b6100ab61024d565b6040518082815260200191505060405180910390f35b6100c9610257565b6040518082815260200191505060405180910390f35b61010b600480360360208110156100f557600080fd5b810190808035906020019092919050505061025d565b005b610115610269565b6040518082815260200191505060405180910390f35b6101616004803603604081101561014157600080fd5b810190808035906020019092919080359060200190929190505050610272565b005b61016b610280565b6040518082815260200191505060405180910390f35b6101b76004803603604081101561019757600080fd5b810190808035906020019092919080359060200190929190505050610286565b005b6101e5600480360360208110156101cf57600080fd5b8101908080359060200190929190505050610298565b005b61021d600480360360408110156101fd57600080fd5b8101908080359060200190929190803590602001909291905050506102a2565b005b61024b6004803603602081101561023557600080fd5b81019080803590602001909291905050506102b0565b005b6000600154905090565b60005481565b610266816102bc565b50565b60008054905090565b61027c8282610286565b5050565b60015481565b81600081905550806001819055505050565b8060008190555050565b6102ac82826102c9565b5050565b6102b981610298565b50565b6001810160008190555050565b8060008190555081600181905550505056fea265627a7a723158205b09cc36dad64daeb008a5977f7430d0ace3fea5faf4e295c631b844323610dd64736f6c634300050d0032";

    public static final String FUNC_GETX = "getX";

    public static final String FUNC_GETY = "getY";

    public static final String FUNC_INIT = "init";

    public static final String FUNC_INITBASE = "initBase";

    public static final String FUNC_INITBASEBASE = "initBaseBase";

    public static final String FUNC_X = "x";

    public static final String FUNC_Y = "y";

    protected InheritContractOverloadChild(String contractAddress, Web3j web3j, Credentials credentials, GasProvider contractGasProvider, Long chainId) {
        super(BINARY, contractAddress, web3j, credentials, contractGasProvider, chainId);
    }

    protected InheritContractOverloadChild(String contractAddress, Web3j web3j, TransactionManager transactionManager, GasProvider contractGasProvider, Long chainId) {
        super(BINARY, contractAddress, web3j, transactionManager, contractGasProvider, chainId);
    }

    public RemoteCall<BigInteger> getX() {
        final Function function = new Function(FUNC_GETX, 
                Arrays.<Type>asList(), 
                Arrays.<TypeReference<?>>asList(new TypeReference<Uint256>() {}));
        return executeRemoteCallSingleValueReturn(function, BigInteger.class);
    }

    public RemoteCall<BigInteger> getY() {
        final Function function = new Function(FUNC_GETY, 
                Arrays.<Type>asList(), 
                Arrays.<TypeReference<?>>asList(new TypeReference<Uint256>() {}));
        return executeRemoteCallSingleValueReturn(function, BigInteger.class);
    }

    public RemoteCall<TransactionReceipt> init(BigInteger a, BigInteger b) {
        final Function function = new Function(
                FUNC_INIT, 
                Arrays.<Type>asList(new org.web3j.abi.datatypes.generated.Uint256(a), 
                new org.web3j.abi.datatypes.generated.Uint256(b)), 
                Collections.<TypeReference<?>>emptyList());
        return executeRemoteCallTransaction(function);
    }

    public RemoteCall<TransactionReceipt> init(BigInteger a) {
        final Function function = new Function(
                FUNC_INIT, 
                Arrays.<Type>asList(new org.web3j.abi.datatypes.generated.Uint256(a)), 
                Collections.<TypeReference<?>>emptyList());
        return executeRemoteCallTransaction(function);
    }

    public RemoteCall<TransactionReceipt> initBase(BigInteger c, BigInteger d) {
        final Function function = new Function(
                FUNC_INITBASE, 
                Arrays.<Type>asList(new org.web3j.abi.datatypes.generated.Uint256(c), 
                new org.web3j.abi.datatypes.generated.Uint256(d)), 
                Collections.<TypeReference<?>>emptyList());
        return executeRemoteCallTransaction(function);
    }

    public RemoteCall<TransactionReceipt> initBase(BigInteger c) {
        final Function function = new Function(
                FUNC_INITBASE, 
                Arrays.<Type>asList(new org.web3j.abi.datatypes.generated.Uint256(c)), 
                Collections.<TypeReference<?>>emptyList());
        return executeRemoteCallTransaction(function);
    }

    public RemoteCall<TransactionReceipt> initBaseBase(BigInteger c) {
        final Function function = new Function(
                FUNC_INITBASEBASE, 
                Arrays.<Type>asList(new org.web3j.abi.datatypes.generated.Uint256(c)), 
                Collections.<TypeReference<?>>emptyList());
        return executeRemoteCallTransaction(function);
    }

    public RemoteCall<TransactionReceipt> initBaseBase(BigInteger c, BigInteger d) {
        final Function function = new Function(
                FUNC_INITBASEBASE, 
                Arrays.<Type>asList(new org.web3j.abi.datatypes.generated.Uint256(c), 
                new org.web3j.abi.datatypes.generated.Uint256(d)), 
                Collections.<TypeReference<?>>emptyList());
        return executeRemoteCallTransaction(function);
    }

    public RemoteCall<BigInteger> x() {
        final Function function = new Function(FUNC_X, 
                Arrays.<Type>asList(), 
                Arrays.<TypeReference<?>>asList(new TypeReference<Uint256>() {}));
        return executeRemoteCallSingleValueReturn(function, BigInteger.class);
    }

    public RemoteCall<BigInteger> y() {
        final Function function = new Function(FUNC_Y, 
                Arrays.<Type>asList(), 
                Arrays.<TypeReference<?>>asList(new TypeReference<Uint256>() {}));
        return executeRemoteCallSingleValueReturn(function, BigInteger.class);
    }

    public static RemoteCall<InheritContractOverloadChild> deploy(Web3j web3j, Credentials credentials, GasProvider contractGasProvider, Long chainId) {
        return deployRemoteCall(InheritContractOverloadChild.class, web3j, credentials, contractGasProvider, BINARY,  "", chainId);
    }

    public static RemoteCall<InheritContractOverloadChild> deploy(Web3j web3j, TransactionManager transactionManager, GasProvider contractGasProvider, Long chainId) {
        return deployRemoteCall(InheritContractOverloadChild.class, web3j, transactionManager, contractGasProvider, BINARY,  "", chainId);
    }

    public static InheritContractOverloadChild load(String contractAddress, Web3j web3j, Credentials credentials, GasProvider contractGasProvider, Long chainId) {
        return new InheritContractOverloadChild(contractAddress, web3j, credentials, contractGasProvider, chainId);
    }

    public static InheritContractOverloadChild load(String contractAddress, Web3j web3j, TransactionManager transactionManager, GasProvider contractGasProvider, Long chainId) {
        return new InheritContractOverloadChild(contractAddress, web3j, transactionManager, contractGasProvider, chainId);
    }
}
