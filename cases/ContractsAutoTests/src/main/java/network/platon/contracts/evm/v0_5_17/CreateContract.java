package network.platon.contracts.evm.v0_5_17;

import com.alaya.abi.solidity.TypeReference;
import com.alaya.abi.solidity.datatypes.Function;
import com.alaya.abi.solidity.datatypes.Type;
import com.alaya.abi.solidity.datatypes.generated.Uint256;
import com.alaya.crypto.Credentials;
import com.alaya.protocol.Web3j;
import com.alaya.protocol.core.RemoteCall;
import com.alaya.tuples.generated.Tuple2;
import com.alaya.tx.Contract;
import com.alaya.tx.TransactionManager;
import com.alaya.tx.gas.GasProvider;
import java.math.BigInteger;
import java.util.Arrays;
import java.util.List;
import java.util.concurrent.Callable;

/**
 * <p>Auto generated code.
 * <p><strong>Do not modify!</strong>
 * <p>Please use the <a href="https://github.com/PlatONnetwork/client-sdk-java/releases">platon-web3j command line tools</a>,
 * or the com.alaya.codegen.SolidityFunctionWrapperGenerator in the 
 * <a href="https://github.com/PlatONnetwork/client-sdk-java/tree/master/codegen">codegen module</a> to update.
 *
 * <p>Generated with web3j version 0.13.2.1.
 */
public class CreateContract extends Contract {
    private static final String BINARY = "60806040526103e86040516100139061008b565b80828152602001915050604051809103906000f080158015610039573d6000803e3d6000fd5b506000806101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff16021790555034801561008557600080fd5b50610098565b6101048061028083390190565b6101d9806100a76000396000f3fe608060405234801561001057600080fd5b506004361061002b5760003560e01c8063f26ca91c14610030575b600080fd5b610038610055565b604051808381526020018281526020019250505060405180910390f35b6000806000809054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16630c55699c6040518163ffffffff1660e01b815260040160206040518083038186803b1580156100bf57600080fd5b505afa1580156100d3573d6000803e3d6000fd5b505050506040513d60208110156100e957600080fd5b81019080805190602001909291905050506000809054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1663aa8c217c6040518163ffffffff1660e01b815260040160206040518083038186803b15801561016157600080fd5b505afa158015610175573d6000803e3d6000fd5b505050506040513d602081101561018b57600080fd5b810190808051906020019092919050505091509150909156fea265627a7a72315820ab26f5422a6d7ec0f95da994255441a798b33f0cf6e6b3b74e85077de0be555464736f6c6343000511003260806040526040516101043803806101048339818101604052602081101561002657600080fd5b810190808051906020019092919050505080600081905550346001819055505060b0806100546000396000f3fe6080604052348015600f57600080fd5b506004361060325760003560e01c80630c55699c146037578063aa8c217c146053575b600080fd5b603d606f565b6040518082815260200191505060405180910390f35b60596075565b6040518082815260200191505060405180910390f35b60005481565b6001548156fea265627a7a723158202206163c67423ec8f497273d5293928ce55f70265c4296a4f30bf7b736ad554f64736f6c63430005110032";

    public static final String FUNC_GETTARGETCREATECONTRACTDATA = "getTargetCreateContractData";

    protected CreateContract(String contractAddress, Web3j web3j, Credentials credentials, GasProvider contractGasProvider, Long chainId) {
        super(BINARY, contractAddress, web3j, credentials, contractGasProvider, chainId);
    }

    protected CreateContract(String contractAddress, Web3j web3j, TransactionManager transactionManager, GasProvider contractGasProvider, Long chainId) {
        super(BINARY, contractAddress, web3j, transactionManager, contractGasProvider, chainId);
    }

    public RemoteCall<Tuple2<BigInteger, BigInteger>> getTargetCreateContractData() {
        final Function function = new Function(FUNC_GETTARGETCREATECONTRACTDATA, 
                Arrays.<Type>asList(), 
                Arrays.<TypeReference<?>>asList(new TypeReference<Uint256>() {}, new TypeReference<Uint256>() {}));
        return new RemoteCall<Tuple2<BigInteger, BigInteger>>(
                new Callable<Tuple2<BigInteger, BigInteger>>() {
                    @Override
                    public Tuple2<BigInteger, BigInteger> call() throws Exception {
                        List<Type> results = executeCallMultipleValueReturn(function);
                        return new Tuple2<BigInteger, BigInteger>(
                                (BigInteger) results.get(0).getValue(), 
                                (BigInteger) results.get(1).getValue());
                    }
                });
    }

    public static RemoteCall<CreateContract> deploy(Web3j web3j, Credentials credentials, GasProvider contractGasProvider, Long chainId) {
        return deployRemoteCall(CreateContract.class, web3j, credentials, contractGasProvider, BINARY,  "", chainId);
    }

    public static RemoteCall<CreateContract> deploy(Web3j web3j, TransactionManager transactionManager, GasProvider contractGasProvider, Long chainId) {
        return deployRemoteCall(CreateContract.class, web3j, transactionManager, contractGasProvider, BINARY,  "", chainId);
    }

    public static CreateContract load(String contractAddress, Web3j web3j, Credentials credentials, GasProvider contractGasProvider, Long chainId) {
        return new CreateContract(contractAddress, web3j, credentials, contractGasProvider, chainId);
    }

    public static CreateContract load(String contractAddress, Web3j web3j, TransactionManager transactionManager, GasProvider contractGasProvider, Long chainId) {
        return new CreateContract(contractAddress, web3j, transactionManager, contractGasProvider, chainId);
    }
}
