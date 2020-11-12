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
public class InheritContractOverloadBase extends Contract {
    private static final String BINARY = "608060405234801561001057600080fd5b506101b5806100206000396000f3fe608060405234801561001057600080fd5b50600436106100625760003560e01c80630b7f1665146100675780630c55699c146100855780635197c7aa146100a3578063a56dfe4a146100c1578063a5843f08146100df578063b7b0422d14610117575b600080fd5b61006f610145565b6040518082815260200191505060405180910390f35b61008d61014f565b6040518082815260200191505060405180910390f35b6100ab610155565b6040518082815260200191505060405180910390f35b6100c961015e565b6040518082815260200191505060405180910390f35b610115600480360360408110156100f557600080fd5b810190808035906020019092919080359060200190929190505050610164565b005b6101436004803603602081101561012d57600080fd5b8101908080359060200190929190505050610176565b005b6000600154905090565b60005481565b60008054905090565b60015481565b81600081905550806001819055505050565b806000819055505056fea265627a7a72315820d74df97a9889319fc2c3ea6a5d55b68aac63bc6c08753cc1c46d69fba151f4d464736f6c634300050d0032";

    public static final String FUNC_GETX = "getX";

    public static final String FUNC_GETY = "getY";

    public static final String FUNC_INIT = "init";

    public static final String FUNC_X = "x";

    public static final String FUNC_Y = "y";

    protected InheritContractOverloadBase(String contractAddress, Web3j web3j, Credentials credentials, GasProvider contractGasProvider, Long chainId) {
        super(BINARY, contractAddress, web3j, credentials, contractGasProvider, chainId);
    }

    protected InheritContractOverloadBase(String contractAddress, Web3j web3j, TransactionManager transactionManager, GasProvider contractGasProvider, Long chainId) {
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
                Arrays.<Type>asList(new Uint256(a),
                new Uint256(b)),
                Collections.<TypeReference<?>>emptyList());
        return executeRemoteCallTransaction(function);
    }

    public RemoteCall<TransactionReceipt> init(BigInteger a) {
        final Function function = new Function(
                FUNC_INIT, 
                Arrays.<Type>asList(new Uint256(a)),
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

    public static RemoteCall<InheritContractOverloadBase> deploy(Web3j web3j, Credentials credentials, GasProvider contractGasProvider, Long chainId) {
        return deployRemoteCall(InheritContractOverloadBase.class, web3j, credentials, contractGasProvider, BINARY,  "", chainId);
    }

    public static RemoteCall<InheritContractOverloadBase> deploy(Web3j web3j, TransactionManager transactionManager, GasProvider contractGasProvider, Long chainId) {
        return deployRemoteCall(InheritContractOverloadBase.class, web3j, transactionManager, contractGasProvider, BINARY,  "", chainId);
    }

    public static InheritContractOverloadBase load(String contractAddress, Web3j web3j, Credentials credentials, GasProvider contractGasProvider, Long chainId) {
        return new InheritContractOverloadBase(contractAddress, web3j, credentials, contractGasProvider, chainId);
    }

    public static InheritContractOverloadBase load(String contractAddress, Web3j web3j, TransactionManager transactionManager, GasProvider contractGasProvider, Long chainId) {
        return new InheritContractOverloadBase(contractAddress, web3j, transactionManager, contractGasProvider, chainId);
    }
}
