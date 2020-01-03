package network.platon.contracts;

import org.web3j.abi.TypeReference;
import org.web3j.abi.datatypes.Function;
import org.web3j.abi.datatypes.Type;
import org.web3j.abi.datatypes.generated.*;
import org.web3j.crypto.Credentials;
import org.web3j.protocol.Web3j;
import org.web3j.protocol.core.RemoteCall;
import org.web3j.tuples.generated.Tuple2;
import org.web3j.tx.Contract;
import org.web3j.tx.TransactionManager;
import org.web3j.tx.gas.GasProvider;

import java.math.BigInteger;
import java.util.Arrays;
import java.util.List;
import java.util.concurrent.Callable;

/**
 * <p>Auto generated code.
 * <p><strong>Do not modify!</strong>
 * <p>Please use the <a href="https://docs.web3j.io/command_line.html">web3j command line tools</a>,
 * or the org.web3j.codegen.SolidityFunctionWrapperGenerator in the 
 * <a href="https://github.com/web3j/web3j/tree/master/codegen">codegen module</a> to update.
 *
 * <p>Generated with web3j version 0.7.5.0.
 */
public class TypeConversionContractTest extends Contract {
    private static final String BINARY = "608060405234801561001057600080fd5b50610399806100206000396000f3fe608060405234801561001057600080fd5b506004361061007d5760003560e01c806399a909621161005b57806399a909621461013b578063a1360967146101aa578063ad42221214610206578063dcefd42f1461022a5761007d565b8063744708f814610082578063853255cc146100f15780639311ca6914610115575b600080fd5b61008a61028c565b604051808363ffffffff1663ffffffff168152602001827bffffffffffffffffffffffffffffffffffffffffffffffffffffffff19167bffffffffffffffffffffffffffffffffffffffffffffffffffffffff191681526020019250505060405180910390f35b6100f96102ae565b604051808260010b60010b815260200191505060405180910390f35b61011d6102c8565b604051808261ffff1661ffff16815260200191505060405180910390f35b6101436102df565b604051808361ffff1661ffff168152602001827dffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff19167dffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff191681526020019250505060405180910390f35b6101b26102ff565b60405180827bffffffffffffffffffffffffffffffffffffffffffffffffffffffff19167bffffffffffffffffffffffffffffffffffffffffffffffffffffffff1916815260200191505060405180910390f35b61020e610338565b604051808260000b60000b815260200191505060405180910390f35b61023261034c565b60405180827effffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff19167effffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff1916815260200191505060405180910390f35b6000806000611234905060008161ffff169050808160e01b9350935050509091565b60008060029050600060649050808260000b019250505090565b600080600a905060008160ff169050809250505090565b6000806000631234567890506000819050808160f01b9350935050509091565b60008061123460f01b90506000817dffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff19169050809250505090565b600080600190506000819050809250505090565b60008061123460f01b9050600081905080925050509056fea265627a7a723158203768035fe05e97cc0897de8591db41bd022fb9b8339162bd5d798ea0793969f564736f6c634300050d0032";

    public static final String FUNC_CONVERSION = "conversion";

    public static final String FUNC_DISPLAYCONVERSION = "displayConversion";

    public static final String FUNC_DISPLAYCONVERSION1 = "displayConversion1";

    public static final String FUNC_DISPLAYCONVERSION2 = "displayConversion2";

    public static final String FUNC_DISPLAYCONVERSION3 = "displayConversion3";

    public static final String FUNC_DISPLAYCONVERSION4 = "displayConversion4";

    public static final String FUNC_SUM = "sum";

    @Deprecated
    protected TypeConversionContractTest(String contractAddress, Web3j web3j, Credentials credentials, BigInteger gasPrice, BigInteger gasLimit) {
        super(BINARY, contractAddress, web3j, credentials, gasPrice, gasLimit);
    }

    protected TypeConversionContractTest(String contractAddress, Web3j web3j, Credentials credentials, GasProvider contractGasProvider) {
        super(BINARY, contractAddress, web3j, credentials, contractGasProvider);
    }

    @Deprecated
    protected TypeConversionContractTest(String contractAddress, Web3j web3j, TransactionManager transactionManager, BigInteger gasPrice, BigInteger gasLimit) {
        super(BINARY, contractAddress, web3j, transactionManager, gasPrice, gasLimit);
    }

    protected TypeConversionContractTest(String contractAddress, Web3j web3j, TransactionManager transactionManager, GasProvider contractGasProvider) {
        super(BINARY, contractAddress, web3j, transactionManager, contractGasProvider);
    }

    public RemoteCall<BigInteger> conversion() {
        final Function function = new Function(FUNC_CONVERSION,
                Arrays.<Type>asList(),
                Arrays.<TypeReference<?>>asList(new TypeReference<Uint16>() {}));
        return executeRemoteCallSingleValueReturn(function, BigInteger.class);
    }

    public RemoteCall<BigInteger> displayConversion() {
        final Function function = new Function(FUNC_DISPLAYCONVERSION,
                Arrays.<Type>asList(),
                Arrays.<TypeReference<?>>asList(new TypeReference<Int8>() {}));
        return executeRemoteCallSingleValueReturn(function, BigInteger.class);
    }

    public RemoteCall<Tuple2<BigInteger, byte[]>> displayConversion1() {
        final Function function = new Function(FUNC_DISPLAYCONVERSION1,
                Arrays.<Type>asList(),
                Arrays.<TypeReference<?>>asList(new TypeReference<Uint16>() {}, new TypeReference<Bytes2>() {}));
        return new RemoteCall<Tuple2<BigInteger, byte[]>>(
                new Callable<Tuple2<BigInteger, byte[]>>() {
                    @Override
                    public Tuple2<BigInteger, byte[]> call() throws Exception {
                        List<Type> results = executeCallMultipleValueReturn(function);
                        return new Tuple2<BigInteger, byte[]>(
                                (BigInteger) results.get(0).getValue(), 
                                (byte[]) results.get(1).getValue());
                    }
                });
    }

    public RemoteCall<Tuple2<BigInteger, byte[]>> displayConversion2() {
        final Function function = new Function(FUNC_DISPLAYCONVERSION2,
                Arrays.<Type>asList(),
                Arrays.<TypeReference<?>>asList(new TypeReference<Uint32>() {}, new TypeReference<Bytes4>() {}));
        return new RemoteCall<Tuple2<BigInteger, byte[]>>(
                new Callable<Tuple2<BigInteger, byte[]>>() {
                    @Override
                    public Tuple2<BigInteger, byte[]> call() throws Exception {
                        List<Type> results = executeCallMultipleValueReturn(function);
                        return new Tuple2<BigInteger, byte[]>(
                                (BigInteger) results.get(0).getValue(), 
                                (byte[]) results.get(1).getValue());
                    }
                });
    }

    public RemoteCall<byte[]> displayConversion3() {
        final Function function = new Function(FUNC_DISPLAYCONVERSION3,
                Arrays.<Type>asList(),
                Arrays.<TypeReference<?>>asList(new TypeReference<Bytes1>() {}));
        return executeRemoteCallSingleValueReturn(function, byte[].class);
    }

    public RemoteCall<byte[]> displayConversion4() {
        final Function function = new Function(FUNC_DISPLAYCONVERSION4,
                Arrays.<Type>asList(),
                Arrays.<TypeReference<?>>asList(new TypeReference<Bytes4>() {}));
        return executeRemoteCallSingleValueReturn(function, byte[].class);
    }

    public RemoteCall<BigInteger> sum() {
        final Function function = new Function(FUNC_SUM,
                Arrays.<Type>asList(),
                Arrays.<TypeReference<?>>asList(new TypeReference<Int16>() {}));
        return executeRemoteCallSingleValueReturn(function, BigInteger.class);
    }

    public static RemoteCall<TypeConversionContractTest> deploy(Web3j web3j, Credentials credentials, GasProvider contractGasProvider) {
        return deployRemoteCall(TypeConversionContractTest.class, web3j, credentials, contractGasProvider, BINARY, "");
    }

    @Deprecated
    public static RemoteCall<TypeConversionContractTest> deploy(Web3j web3j, Credentials credentials, BigInteger gasPrice, BigInteger gasLimit) {
        return deployRemoteCall(TypeConversionContractTest.class, web3j, credentials, gasPrice, gasLimit, BINARY, "");
    }

    public static RemoteCall<TypeConversionContractTest> deploy(Web3j web3j, TransactionManager transactionManager, GasProvider contractGasProvider) {
        return deployRemoteCall(TypeConversionContractTest.class, web3j, transactionManager, contractGasProvider, BINARY, "");
    }

    @Deprecated
    public static RemoteCall<TypeConversionContractTest> deploy(Web3j web3j, TransactionManager transactionManager, BigInteger gasPrice, BigInteger gasLimit) {
        return deployRemoteCall(TypeConversionContractTest.class, web3j, transactionManager, gasPrice, gasLimit, BINARY, "");
    }

    @Deprecated
    public static TypeConversionContractTest load(String contractAddress, Web3j web3j, Credentials credentials, BigInteger gasPrice, BigInteger gasLimit) {
        return new TypeConversionContractTest(contractAddress, web3j, credentials, gasPrice, gasLimit);
    }

    @Deprecated
    public static TypeConversionContractTest load(String contractAddress, Web3j web3j, TransactionManager transactionManager, BigInteger gasPrice, BigInteger gasLimit) {
        return new TypeConversionContractTest(contractAddress, web3j, transactionManager, gasPrice, gasLimit);
    }

    public static TypeConversionContractTest load(String contractAddress, Web3j web3j, Credentials credentials, GasProvider contractGasProvider) {
        return new TypeConversionContractTest(contractAddress, web3j, credentials, contractGasProvider);
    }

    public static TypeConversionContractTest load(String contractAddress, Web3j web3j, TransactionManager transactionManager, GasProvider contractGasProvider) {
        return new TypeConversionContractTest(contractAddress, web3j, transactionManager, contractGasProvider);
    }
}
