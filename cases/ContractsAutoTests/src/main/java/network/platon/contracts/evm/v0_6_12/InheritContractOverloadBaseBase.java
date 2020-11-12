package network.platon.contracts.evm.v0_6_12;

import com.alaya.abi.solidity.TypeReference;
import com.alaya.abi.solidity.datatypes.Function;
import com.alaya.abi.solidity.datatypes.Type;
import com.alaya.abi.solidity.datatypes.generated.Uint256;
import com.alaya.crypto.Credentials;
import com.alaya.protocol.Web3j;
import com.alaya.protocol.core.RemoteCall;
import com.alaya.protocol.core.methods.response.TransactionReceipt;
import com.alaya.tx.Contract;
import com.alaya.tx.TransactionManager;
import com.alaya.tx.gas.GasProvider;
import java.math.BigInteger;
import java.util.Arrays;
import java.util.Collections;

/**
 * <p>Auto generated code.
 * <p><strong>Do not modify!</strong>
 * <p>Please use the <a href="https://github.com/PlatONnetwork/client-sdk-java/releases">platon-web3j command line tools</a>,
 * or the com.alaya.codegen.SolidityFunctionWrapperGenerator in the 
 * <a href="https://github.com/PlatONnetwork/client-sdk-java/tree/master/codegen">codegen module</a> to update.
 *
 * <p>Generated with web3j version 0.13.2.1.
 */
public class InheritContractOverloadBaseBase extends Contract {
    private static final String BINARY = "608060405234801561001057600080fd5b506101b9806100206000396000f3fe608060405234801561001057600080fd5b50600436106100625760003560e01c80630b7f1665146100675780630c55699c146100855780635197c7aa146100a3578063a56dfe4a146100c1578063a5843f08146100df578063b7b0422d14610117575b600080fd5b61006f610145565b6040518082815260200191505060405180910390f35b61008d61014f565b6040518082815260200191505060405180910390f35b6100ab610155565b6040518082815260200191505060405180910390f35b6100c961015e565b6040518082815260200191505060405180910390f35b610115600480360360408110156100f557600080fd5b810190808035906020019092919080359060200190929190505050610164565b005b6101436004803603602081101561012d57600080fd5b8101908080359060200190929190505050610176565b005b6000600154905090565b60005481565b60008054905090565b60015481565b80600081905550816001819055505050565b600181016000819055505056fea2646970667358221220b492544f27bc0fdb12aaf56e57c87d396d769f932b702623fe35a6e9e9842b4864736f6c634300060c0033";

    public static final String FUNC_GETX = "getX";

    public static final String FUNC_GETY = "getY";

    public static final String FUNC_INIT = "init";

    public static final String FUNC_X = "x";

    public static final String FUNC_Y = "y";

    protected InheritContractOverloadBaseBase(String contractAddress, Web3j web3j, Credentials credentials, GasProvider contractGasProvider, Long chainId) {
        super(BINARY, contractAddress, web3j, credentials, contractGasProvider, chainId);
    }

    protected InheritContractOverloadBaseBase(String contractAddress, Web3j web3j, TransactionManager transactionManager, GasProvider contractGasProvider, Long chainId) {
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

    public static RemoteCall<InheritContractOverloadBaseBase> deploy(Web3j web3j, Credentials credentials, GasProvider contractGasProvider, Long chainId) {
        return deployRemoteCall(InheritContractOverloadBaseBase.class, web3j, credentials, contractGasProvider, BINARY,  "", chainId);
    }

    public static RemoteCall<InheritContractOverloadBaseBase> deploy(Web3j web3j, TransactionManager transactionManager, GasProvider contractGasProvider, Long chainId) {
        return deployRemoteCall(InheritContractOverloadBaseBase.class, web3j, transactionManager, contractGasProvider, BINARY,  "", chainId);
    }

    public static InheritContractOverloadBaseBase load(String contractAddress, Web3j web3j, Credentials credentials, GasProvider contractGasProvider, Long chainId) {
        return new InheritContractOverloadBaseBase(contractAddress, web3j, credentials, contractGasProvider, chainId);
    }

    public static InheritContractOverloadBaseBase load(String contractAddress, Web3j web3j, TransactionManager transactionManager, GasProvider contractGasProvider, Long chainId) {
        return new InheritContractOverloadBaseBase(contractAddress, web3j, transactionManager, contractGasProvider, chainId);
    }
}
