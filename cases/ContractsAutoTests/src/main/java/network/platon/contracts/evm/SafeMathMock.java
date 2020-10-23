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
public class SafeMathMock extends Contract {
    private static final String BINARY = "608060405234801561001057600080fd5b50610708806100206000396000f3fe608060405234801561001057600080fd5b50600436106100885760003560e01c8063a391c15b1161005b578063a391c15b146101bd578063b67d77c514610209578063c8a4ac9c14610255578063f43f523a146102a157610088565b80632b7423ab1461008d5780636d5433e6146100d9578063771602f7146101255780637ae2b5c714610171575b600080fd5b6100c3600480360360408110156100a357600080fd5b8101908080359060200190929190803590602001909291905050506102ed565b6040518082815260200191505060405180910390f35b61010f600480360360408110156100ef57600080fd5b810190808035906020019092919080359060200190929190505050610301565b6040518082815260200191505060405180910390f35b61015b6004803603604081101561013b57600080fd5b810190808035906020019092919080359060200190929190505050610315565b6040518082815260200191505060405180910390f35b6101a76004803603604081101561018757600080fd5b810190808035906020019092919080359060200190929190505050610329565b6040518082815260200191505060405180910390f35b6101f3600480360360408110156101d357600080fd5b81019080803590602001909291908035906020019092919050505061033d565b6040518082815260200191505060405180910390f35b61023f6004803603604081101561021f57600080fd5b810190808035906020019092919080359060200190929190505050610351565b6040518082815260200191505060405180910390f35b61028b6004803603604081101561026b57600080fd5b810190808035906020019092919080359060200190929190505050610365565b6040518082815260200191505060405180910390f35b6102d7600480360360408110156102b757600080fd5b810190808035906020019092919080359060200190929190505050610379565b6040518082815260200191505060405180910390f35b60006102f9838361038d565b905092915050565b600061030d83836103cf565b905092915050565b600061032183836103e9565b905092915050565b60006103358383610471565b905092915050565b6000610349838361048a565b905092915050565b600061035d8383610519565b905092915050565b600061037183836105a2565b905092915050565b60006103858383610628565b905092915050565b6000600280838161039a57fe5b06600285816103a557fe5b0601816103ae57fe5b04600283816103b957fe5b04600285816103c457fe5b040101905092915050565b6000818310156103df57816103e1565b825b905092915050565b600080828401905083811015610467576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040180806020018281038252601b8152602001807f536166654d6174683a206164646974696f6e206f766572666c6f77000000000081525060200191505060405180910390fd5b8091505092915050565b60008183106104805781610482565b825b905092915050565b6000808211610501576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040180806020018281038252601a8152602001807f536166654d6174683a206469766973696f6e206279207a65726f00000000000081525060200191505060405180910390fd5b600082848161050c57fe5b0490508091505092915050565b600082821115610591576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040180806020018281038252601e8152602001807f536166654d6174683a207375627472616374696f6e206f766572666c6f77000081525060200191505060405180910390fd5b600082840390508091505092915050565b6000808314156105b55760009050610622565b60008284029050828482816105c657fe5b041461061d576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004018080602001828103825260218152602001806106b36021913960400191505060405180910390fd5b809150505b92915050565b6000808214156106a0576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004018080602001828103825260188152602001807f536166654d6174683a206d6f64756c6f206279207a65726f000000000000000081525060200191505060405180910390fd5b8183816106a957fe5b0690509291505056fe536166654d6174683a206d756c7469706c69636174696f6e206f766572666c6f77a265627a7a723158200aac66050447e47a7b3a1350b2ec89f6dd9f4cff7c0cb575daf0a4e06d32d48164736f6c634300050d0032";

    public static final String FUNC_ADD = "add";

    public static final String FUNC_AVERAGE = "average";

    public static final String FUNC_DIV = "div";

    public static final String FUNC_MAX = "max";

    public static final String FUNC_MIN = "min";

    public static final String FUNC_MOD = "mod";

    public static final String FUNC_MUL = "mul";

    public static final String FUNC_SUB = "sub";

    protected SafeMathMock(String contractAddress, Web3j web3j, Credentials credentials, GasProvider contractGasProvider, Long chainId) {
        super(BINARY, contractAddress, web3j, credentials, contractGasProvider, chainId);
    }

    protected SafeMathMock(String contractAddress, Web3j web3j, TransactionManager transactionManager, GasProvider contractGasProvider, Long chainId) {
        super(BINARY, contractAddress, web3j, transactionManager, contractGasProvider, chainId);
    }

    public RemoteCall<BigInteger> add(BigInteger a, BigInteger b) {
        final Function function = new Function(FUNC_ADD, 
                Arrays.<Type>asList(new Uint256(a),
                new Uint256(b)),
                Arrays.<TypeReference<?>>asList(new TypeReference<Uint256>() {}));
        return executeRemoteCallSingleValueReturn(function, BigInteger.class);
    }

    public RemoteCall<BigInteger> average(BigInteger a, BigInteger b) {
        final Function function = new Function(FUNC_AVERAGE, 
                Arrays.<Type>asList(new Uint256(a),
                new Uint256(b)),
                Arrays.<TypeReference<?>>asList(new TypeReference<Uint256>() {}));
        return executeRemoteCallSingleValueReturn(function, BigInteger.class);
    }

    public RemoteCall<BigInteger> div(BigInteger a, BigInteger b) {
        final Function function = new Function(FUNC_DIV, 
                Arrays.<Type>asList(new Uint256(a),
                new Uint256(b)),
                Arrays.<TypeReference<?>>asList(new TypeReference<Uint256>() {}));
        return executeRemoteCallSingleValueReturn(function, BigInteger.class);
    }

    public RemoteCall<BigInteger> max(BigInteger a, BigInteger b) {
        final Function function = new Function(FUNC_MAX, 
                Arrays.<Type>asList(new Uint256(a),
                new Uint256(b)),
                Arrays.<TypeReference<?>>asList(new TypeReference<Uint256>() {}));
        return executeRemoteCallSingleValueReturn(function, BigInteger.class);
    }

    public RemoteCall<BigInteger> min(BigInteger a, BigInteger b) {
        final Function function = new Function(FUNC_MIN, 
                Arrays.<Type>asList(new Uint256(a),
                new Uint256(b)),
                Arrays.<TypeReference<?>>asList(new TypeReference<Uint256>() {}));
        return executeRemoteCallSingleValueReturn(function, BigInteger.class);
    }

    public RemoteCall<BigInteger> mod(BigInteger a, BigInteger b) {
        final Function function = new Function(FUNC_MOD, 
                Arrays.<Type>asList(new Uint256(a),
                new Uint256(b)),
                Arrays.<TypeReference<?>>asList(new TypeReference<Uint256>() {}));
        return executeRemoteCallSingleValueReturn(function, BigInteger.class);
    }

    public RemoteCall<BigInteger> mul(BigInteger a, BigInteger b) {
        final Function function = new Function(FUNC_MUL, 
                Arrays.<Type>asList(new Uint256(a),
                new Uint256(b)),
                Arrays.<TypeReference<?>>asList(new TypeReference<Uint256>() {}));
        return executeRemoteCallSingleValueReturn(function, BigInteger.class);
    }

    public RemoteCall<BigInteger> sub(BigInteger a, BigInteger b) {
        final Function function = new Function(FUNC_SUB, 
                Arrays.<Type>asList(new Uint256(a),
                new Uint256(b)),
                Arrays.<TypeReference<?>>asList(new TypeReference<Uint256>() {}));
        return executeRemoteCallSingleValueReturn(function, BigInteger.class);
    }

    public static RemoteCall<SafeMathMock> deploy(Web3j web3j, Credentials credentials, GasProvider contractGasProvider, Long chainId) {
        return deployRemoteCall(SafeMathMock.class, web3j, credentials, contractGasProvider, BINARY,  "", chainId);
    }

    public static RemoteCall<SafeMathMock> deploy(Web3j web3j, TransactionManager transactionManager, GasProvider contractGasProvider, Long chainId) {
        return deployRemoteCall(SafeMathMock.class, web3j, transactionManager, contractGasProvider, BINARY,  "", chainId);
    }

    public static SafeMathMock load(String contractAddress, Web3j web3j, Credentials credentials, GasProvider contractGasProvider, Long chainId) {
        return new SafeMathMock(contractAddress, web3j, credentials, contractGasProvider, chainId);
    }

    public static SafeMathMock load(String contractAddress, Web3j web3j, TransactionManager transactionManager, GasProvider contractGasProvider, Long chainId) {
        return new SafeMathMock(contractAddress, web3j, transactionManager, contractGasProvider, chainId);
    }
}
