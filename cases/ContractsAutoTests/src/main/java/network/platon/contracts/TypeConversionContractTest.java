package network.platon.contracts;

import java.math.BigInteger;
import java.util.Arrays;
import java.util.List;
import java.util.concurrent.Callable;
import org.web3j.abi.TypeReference;
import org.web3j.abi.datatypes.Function;
import org.web3j.abi.datatypes.Type;
import org.web3j.abi.datatypes.generated.Bytes1;
import org.web3j.abi.datatypes.generated.Bytes2;
import org.web3j.abi.datatypes.generated.Bytes4;
import org.web3j.abi.datatypes.generated.Int16;
import org.web3j.abi.datatypes.generated.Int8;
import org.web3j.abi.datatypes.generated.Uint16;
import org.web3j.abi.datatypes.generated.Uint32;
import org.web3j.abi.datatypes.generated.Uint8;
import org.web3j.crypto.Credentials;
import org.web3j.protocol.Web3j;
import org.web3j.protocol.core.RemoteCall;
import org.web3j.tuples.generated.Tuple2;
import org.web3j.tuples.generated.Tuple4;
import org.web3j.tx.Contract;
import org.web3j.tx.TransactionManager;
import org.web3j.tx.gas.GasProvider;

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
    private static final String BINARY = "608060405260026000806101000a81548160ff021916908360000b60ff1602179055506064600060016101000a81548161ffff021916908360010b61ffff16021790555034801561004f57600080fd5b506104768061005f6000396000f3fe608060405234801561001057600080fd5b50600436106100885760003560e01c806399a909621161005b57806399a909621461019d578063a13609671461020c578063ad42221214610268578063dcefd42f1461028c57610088565b80633610b66f1461008d578063744708f8146100e4578063853255cc146101535780639311ca6914610177575b600080fd5b6100956102ee565b604051808563ffffffff1663ffffffff1681526020018463ffffffff1663ffffffff1681526020018360ff1660ff1681526020018260ff1660ff16815260200194505050505060405180910390f35b6100ec61035a565b604051808363ffffffff1663ffffffff168152602001827bffffffffffffffffffffffffffffffffffffffffffffffffffffffff19167bffffffffffffffffffffffffffffffffffffffffffffffffffffffff191681526020019250505060405180910390f35b61015b61037c565b604051808260010b60010b815260200191505060405180910390f35b61017f6103a5565b604051808261ffff1661ffff16815260200191505060405180910390f35b6101a56103bc565b604051808361ffff1661ffff168152602001827dffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff19167dffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff191681526020019250505060405180910390f35b6102146103dc565b60405180827bffffffffffffffffffffffffffffffffffffffffffffffffffffffff19167bffffffffffffffffffffffffffffffffffffffffffffffffffffffff1916815260200191505060405180910390f35b610270610415565b604051808260000b60000b815260200191505060405180910390f35b610294610429565b60405180827effffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff19167effffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff1916815260200191505060405180910390f35b600080600080600061123460f01b905060008160f01c61ffff1690506000827dffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff191660e01c905060008360f01c905060008460f81c9050838383839850985098509850505050505090919293565b6000806000611234905060008161ffff169050808160e01b9350935050509091565b60008060019054906101000a900460010b6000809054906101000a900460000b60000b01905090565b600080600a905060008160ff169050809250505090565b6000806000631234567890506000819050808160f01b9350935050509091565b60008061123460f01b90506000817dffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff19169050809250505090565b600080600190506000819050809250505090565b60008061123460f01b9050600081905080925050509056fea265627a7a7231582035682a625d03aaa32a5dc0eda8aea42533b6e915aa775c25fe13735191ee274a64736f6c634300050d0032";

    public static final String FUNC_CONVERSION = "conversion";

    public static final String FUNC_DISPLAYCONVERSION = "displayConversion";

    public static final String FUNC_DISPLAYCONVERSION1 = "displayConversion1";

    public static final String FUNC_DISPLAYCONVERSION2 = "displayConversion2";

    public static final String FUNC_DISPLAYCONVERSION3 = "displayConversion3";

    public static final String FUNC_DISPLAYCONVERSION4 = "displayConversion4";

    public static final String FUNC_DISPLAYCONVERSION5 = "displayConversion5";

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

    public RemoteCall<Tuple4<BigInteger, BigInteger, BigInteger, BigInteger>> displayConversion5() {
        final Function function = new Function(FUNC_DISPLAYCONVERSION5, 
                Arrays.<Type>asList(), 
                Arrays.<TypeReference<?>>asList(new TypeReference<Uint32>() {}, new TypeReference<Uint32>() {}, new TypeReference<Uint8>() {}, new TypeReference<Uint8>() {}));
        return new RemoteCall<Tuple4<BigInteger, BigInteger, BigInteger, BigInteger>>(
                new Callable<Tuple4<BigInteger, BigInteger, BigInteger, BigInteger>>() {
                    @Override
                    public Tuple4<BigInteger, BigInteger, BigInteger, BigInteger> call() throws Exception {
                        List<Type> results = executeCallMultipleValueReturn(function);
                        return new Tuple4<BigInteger, BigInteger, BigInteger, BigInteger>(
                                (BigInteger) results.get(0).getValue(), 
                                (BigInteger) results.get(1).getValue(), 
                                (BigInteger) results.get(2).getValue(), 
                                (BigInteger) results.get(3).getValue());
                    }
                });
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
