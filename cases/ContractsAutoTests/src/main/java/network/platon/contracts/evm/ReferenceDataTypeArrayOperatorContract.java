package network.platon.contracts.evm;

import java.math.BigInteger;
import java.util.Arrays;
import java.util.List;
import java.util.concurrent.Callable;
import org.web3j.abi.TypeReference;
import org.web3j.abi.datatypes.Bool;
import org.web3j.abi.datatypes.Function;
import org.web3j.abi.datatypes.Type;
import org.web3j.abi.datatypes.generated.Bytes1;
import org.web3j.abi.datatypes.generated.Uint8;
import org.web3j.crypto.Credentials;
import org.web3j.protocol.Web3j;
import org.web3j.protocol.core.RemoteCall;
import org.web3j.tuples.generated.Tuple2;
import org.web3j.tuples.generated.Tuple5;
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
public class ReferenceDataTypeArrayOperatorContract extends Contract {
    private static final String BINARY = "608060405234801561001057600080fd5b506106e5806100206000396000f3fe608060405234801561001057600080fd5b506004361061007d5760003560e01c806399a1c3691161005b57806399a1c369146101cf5780639a6e3cb71461023e578063b99033341461028c578063fd081ef8146102fb5761007d565b80631ff0db4014610082578063676cb904146100f157806386d1171014610160575b600080fd5b61008a61036a565b60405180837effffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff19167effffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff191681526020018260ff1660ff1681526020019250505060405180910390f35b6100f961039c565b60405180837effffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff19167effffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff191681526020018260ff1660ff1681526020019250505060405180910390f35b6101686103e8565b60405180837effffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff19167effffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff191681526020018260ff1660ff1681526020019250505060405180910390f35b6101d7610434565b60405180837effffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff19167effffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff191681526020018260ff1660ff1681526020019250505060405180910390f35b610246610466565b6040518086151515158152602001851515151581526020018415151515815260200183151515158152602001821515151581526020019550505050505060405180910390f35b610294610657565b60405180837effffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff19167effffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff191681526020018260ff1660ff1681526020019250505060405180910390f35b610303610689565b60405180837effffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff19167effffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff191681526020018260ff1660ff1681526020019250505060405180910390f35b6000806000608160f81b90506000608060f81b90506000818316905060008160f81c9050818195509550505050509091565b6000806000608160f81b905060006001827effffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff1916901b905060008160f81c90508181945094505050509091565b6000806000608160f81b905060006001827effffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff1916901c905060008160f81c90508181945094505050509091565b6000806000608160f81b90506000608060f81b90506000818318905060008160f81c9050818195509550505050509091565b6000806000806000807f6100000000000000000000000000000000000000000000000000000000000000905060007f620000000000000000000000000000000000000000000000000000000000000090506000606160f81b90506000827effffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff1916847effffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff19161090506000837effffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff1916857effffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff19161190506000837effffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff1916867effffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff19161490506000847effffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff1916867effffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff1916141590506000857effffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff1916887effffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff19161015905084848484849c509c509c509c509c5050505050505050509091929394565b6000806000608160f81b90506000608060f81b90506000818317905060008160f81c9050818195509550505050509091565b6000806000608160f81b905060008119905060008160f81c9050818194509450505050909156fea265627a7a72315820e63914aff4803852d659ad346d20d68b33d91e2652c72f49378f06184b9a197a64736f6c634300050d0032";

    public static final String FUNC_ARRAYBITANDOPERATORS = "arrayBitAndOperators";

    public static final String FUNC_ARRAYBITINVERSEOPERATORS = "arrayBitInverseOperators";

    public static final String FUNC_ARRAYBITLEFTSHIFTPERATORS = "arrayBitLeftShiftperators";

    public static final String FUNC_ARRAYBITOROPERATORS = "arrayBitOrOperators";

    public static final String FUNC_ARRAYBITRIGHTSHIFTPERATORS = "arrayBitRightShiftperators";

    public static final String FUNC_ARRAYBITXOROPERATORS = "arrayBitXOROperators";

    public static final String FUNC_ARRAYCOMPARE = "arrayCompare";

    protected ReferenceDataTypeArrayOperatorContract(String contractAddress, Web3j web3j, Credentials credentials, GasProvider contractGasProvider, Long chainId) {
        super(BINARY, contractAddress, web3j, credentials, contractGasProvider, chainId);
    }

    protected ReferenceDataTypeArrayOperatorContract(String contractAddress, Web3j web3j, TransactionManager transactionManager, GasProvider contractGasProvider, Long chainId) {
        super(BINARY, contractAddress, web3j, transactionManager, contractGasProvider, chainId);
    }

    public RemoteCall<Tuple2<byte[], BigInteger>> arrayBitAndOperators() {
        final Function function = new Function(FUNC_ARRAYBITANDOPERATORS, 
                Arrays.<Type>asList(), 
                Arrays.<TypeReference<?>>asList(new TypeReference<Bytes1>() {}, new TypeReference<Uint8>() {}));
        return new RemoteCall<Tuple2<byte[], BigInteger>>(
                new Callable<Tuple2<byte[], BigInteger>>() {
                    @Override
                    public Tuple2<byte[], BigInteger> call() throws Exception {
                        List<Type> results = executeCallMultipleValueReturn(function);
                        return new Tuple2<byte[], BigInteger>(
                                (byte[]) results.get(0).getValue(), 
                                (BigInteger) results.get(1).getValue());
                    }
                });
    }

    public RemoteCall<Tuple2<byte[], BigInteger>> arrayBitInverseOperators() {
        final Function function = new Function(FUNC_ARRAYBITINVERSEOPERATORS, 
                Arrays.<Type>asList(), 
                Arrays.<TypeReference<?>>asList(new TypeReference<Bytes1>() {}, new TypeReference<Uint8>() {}));
        return new RemoteCall<Tuple2<byte[], BigInteger>>(
                new Callable<Tuple2<byte[], BigInteger>>() {
                    @Override
                    public Tuple2<byte[], BigInteger> call() throws Exception {
                        List<Type> results = executeCallMultipleValueReturn(function);
                        return new Tuple2<byte[], BigInteger>(
                                (byte[]) results.get(0).getValue(), 
                                (BigInteger) results.get(1).getValue());
                    }
                });
    }

    public RemoteCall<Tuple2<byte[], BigInteger>> arrayBitLeftShiftperators() {
        final Function function = new Function(FUNC_ARRAYBITLEFTSHIFTPERATORS, 
                Arrays.<Type>asList(), 
                Arrays.<TypeReference<?>>asList(new TypeReference<Bytes1>() {}, new TypeReference<Uint8>() {}));
        return new RemoteCall<Tuple2<byte[], BigInteger>>(
                new Callable<Tuple2<byte[], BigInteger>>() {
                    @Override
                    public Tuple2<byte[], BigInteger> call() throws Exception {
                        List<Type> results = executeCallMultipleValueReturn(function);
                        return new Tuple2<byte[], BigInteger>(
                                (byte[]) results.get(0).getValue(), 
                                (BigInteger) results.get(1).getValue());
                    }
                });
    }

    public RemoteCall<Tuple2<byte[], BigInteger>> arrayBitOrOperators() {
        final Function function = new Function(FUNC_ARRAYBITOROPERATORS, 
                Arrays.<Type>asList(), 
                Arrays.<TypeReference<?>>asList(new TypeReference<Bytes1>() {}, new TypeReference<Uint8>() {}));
        return new RemoteCall<Tuple2<byte[], BigInteger>>(
                new Callable<Tuple2<byte[], BigInteger>>() {
                    @Override
                    public Tuple2<byte[], BigInteger> call() throws Exception {
                        List<Type> results = executeCallMultipleValueReturn(function);
                        return new Tuple2<byte[], BigInteger>(
                                (byte[]) results.get(0).getValue(), 
                                (BigInteger) results.get(1).getValue());
                    }
                });
    }

    public RemoteCall<Tuple2<byte[], BigInteger>> arrayBitRightShiftperators() {
        final Function function = new Function(FUNC_ARRAYBITRIGHTSHIFTPERATORS, 
                Arrays.<Type>asList(), 
                Arrays.<TypeReference<?>>asList(new TypeReference<Bytes1>() {}, new TypeReference<Uint8>() {}));
        return new RemoteCall<Tuple2<byte[], BigInteger>>(
                new Callable<Tuple2<byte[], BigInteger>>() {
                    @Override
                    public Tuple2<byte[], BigInteger> call() throws Exception {
                        List<Type> results = executeCallMultipleValueReturn(function);
                        return new Tuple2<byte[], BigInteger>(
                                (byte[]) results.get(0).getValue(), 
                                (BigInteger) results.get(1).getValue());
                    }
                });
    }

    public RemoteCall<Tuple2<byte[], BigInteger>> arrayBitXOROperators() {
        final Function function = new Function(FUNC_ARRAYBITXOROPERATORS, 
                Arrays.<Type>asList(), 
                Arrays.<TypeReference<?>>asList(new TypeReference<Bytes1>() {}, new TypeReference<Uint8>() {}));
        return new RemoteCall<Tuple2<byte[], BigInteger>>(
                new Callable<Tuple2<byte[], BigInteger>>() {
                    @Override
                    public Tuple2<byte[], BigInteger> call() throws Exception {
                        List<Type> results = executeCallMultipleValueReturn(function);
                        return new Tuple2<byte[], BigInteger>(
                                (byte[]) results.get(0).getValue(), 
                                (BigInteger) results.get(1).getValue());
                    }
                });
    }

    public RemoteCall<Tuple5<Boolean, Boolean, Boolean, Boolean, Boolean>> arrayCompare() {
        final Function function = new Function(FUNC_ARRAYCOMPARE, 
                Arrays.<Type>asList(), 
                Arrays.<TypeReference<?>>asList(new TypeReference<Bool>() {}, new TypeReference<Bool>() {}, new TypeReference<Bool>() {}, new TypeReference<Bool>() {}, new TypeReference<Bool>() {}));
        return new RemoteCall<Tuple5<Boolean, Boolean, Boolean, Boolean, Boolean>>(
                new Callable<Tuple5<Boolean, Boolean, Boolean, Boolean, Boolean>>() {
                    @Override
                    public Tuple5<Boolean, Boolean, Boolean, Boolean, Boolean> call() throws Exception {
                        List<Type> results = executeCallMultipleValueReturn(function);
                        return new Tuple5<Boolean, Boolean, Boolean, Boolean, Boolean>(
                                (Boolean) results.get(0).getValue(), 
                                (Boolean) results.get(1).getValue(), 
                                (Boolean) results.get(2).getValue(), 
                                (Boolean) results.get(3).getValue(), 
                                (Boolean) results.get(4).getValue());
                    }
                });
    }

    public static RemoteCall<ReferenceDataTypeArrayOperatorContract> deploy(Web3j web3j, Credentials credentials, GasProvider contractGasProvider, Long chainId) {
        return deployRemoteCall(ReferenceDataTypeArrayOperatorContract.class, web3j, credentials, contractGasProvider, BINARY,  "", chainId);
    }

    public static RemoteCall<ReferenceDataTypeArrayOperatorContract> deploy(Web3j web3j, TransactionManager transactionManager, GasProvider contractGasProvider, Long chainId) {
        return deployRemoteCall(ReferenceDataTypeArrayOperatorContract.class, web3j, transactionManager, contractGasProvider, BINARY,  "", chainId);
    }

    public static ReferenceDataTypeArrayOperatorContract load(String contractAddress, Web3j web3j, Credentials credentials, GasProvider contractGasProvider, Long chainId) {
        return new ReferenceDataTypeArrayOperatorContract(contractAddress, web3j, credentials, contractGasProvider, chainId);
    }

    public static ReferenceDataTypeArrayOperatorContract load(String contractAddress, Web3j web3j, TransactionManager transactionManager, GasProvider contractGasProvider, Long chainId) {
        return new ReferenceDataTypeArrayOperatorContract(contractAddress, web3j, transactionManager, contractGasProvider, chainId);
    }
}
