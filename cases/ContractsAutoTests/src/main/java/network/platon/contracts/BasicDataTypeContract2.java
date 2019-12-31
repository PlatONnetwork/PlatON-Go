package network.platon.contracts;

import java.math.BigInteger;
import java.util.Arrays;
import java.util.Collections;
import java.util.List;
import java.util.concurrent.Callable;
import org.web3j.abi.TypeReference;
import org.web3j.abi.datatypes.Function;
import org.web3j.abi.datatypes.Type;
import org.web3j.abi.datatypes.Utf8String;
import org.web3j.abi.datatypes.generated.Bytes1;
import org.web3j.abi.datatypes.generated.Bytes2;
import org.web3j.abi.datatypes.generated.Uint128;
import org.web3j.abi.datatypes.generated.Uint256;
import org.web3j.abi.datatypes.generated.Uint8;
import org.web3j.crypto.Credentials;
import org.web3j.protocol.Web3j;
import org.web3j.protocol.core.RemoteCall;
import org.web3j.protocol.core.methods.response.TransactionReceipt;
import org.web3j.tuples.generated.Tuple3;
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
public class BasicDataTypeContract2 extends Contract {
    private static final String BINARY = "60806040526040518060400160405280600581526020017f68656c6c6f0000000000000000000000000000000000000000000000000000008152506000908051906020019061004f9291906100ae565b506040518060400160405280600581526020017f776f726c640000000000000000000000000000000000000000000000000000008152506001908051906020019061009b9291906100ae565b503480156100a857600080fd5b50610153565b828054600181600116156101000203166002900490600052602060002090601f016020900481019282601f106100ef57805160ff191683800117855561011d565b8280016001018555821561011d579182015b8281111561011c578251825591602001919060010190610101565b5b50905061012a919061012e565b5090565b61015091905b8082111561014c576000816000905550600101610134565b5090565b90565b610b60806101626000396000f3fe6080604052600436106100fe5760003560e01c80639b8751a411610095578063ac917d6511610064578063ac917d651461053d578063e064d480146105ac578063ea07cdfa1461063c578063f8b2cb4f14610667578063f922544a146106cc576100fe565b80639b8751a414610323578063a0ee573514610426578063a574971014610482578063ac805759146104ad576100fe565b806338cc4831116100d157806338cc48311461023d57806359a1cf111461029457806362738998146102bf5780637d82a1e1146102ea576100fe565b80630ebefd6b14610103578063166aa6e61461013c578063209652551461019c57806338023fb9146101f9575b600080fd5b34801561010f57600080fd5b50610118610739565b6040518082600381111561012857fe5b60ff16815260200191505060405180910390f35b34801561014857600080fd5b506101786004803603602081101561015f57600080fd5b81019080803560ff169060200190929190505050610751565b6040518082600381111561018857fe5b60ff16815260200191505060405180910390f35b3480156101a857600080fd5b506101b161075b565b60405180846fffffffffffffffffffffffffffffffff166fffffffffffffffffffffffffffffffff168152602001838152602001828152602001935050505060405180910390f35b61023b6004803603602081101561020f57600080fd5b81019080803573ffffffffffffffffffffffffffffffffffffffff16906020019092919050505061078c565b005b34801561024957600080fd5b506102526107d6565b604051808273ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200191505060405180910390f35b3480156102a057600080fd5b506102a96107f7565b6040518082815260200191505060405180910390f35b3480156102cb57600080fd5b506102d4610814565b6040518082815260200191505060405180910390f35b3480156102f657600080fd5b506102ff610822565b6040518082600381111561030f57fe5b60ff16815260200191505060405180910390f35b34801561032f57600080fd5b50610338610833565b60405180847dffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff19167dffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff19168152602001837effffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff19167effffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff19168152602001827effffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff19167effffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff19168152602001935050505060405180910390f35b61042e61088b565b604051808473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200183815260200182151515158152602001935050505060405180910390f35b34801561048e57600080fd5b506104976108d2565b6040518082815260200191505060405180910390f35b3480156104b957600080fd5b506104c26108f1565b6040518080602001828103825283818151815260200191508051906020019080838360005b838110156105025780820151818401526020810190506104e7565b50505050905090810190601f16801561052f5780820380516001836020036101000a031916815260200191505b509250505060405180910390f35b34801561054957600080fd5b50610552610993565b60405180827effffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff19167effffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff1916815260200191505060405180910390f35b3480156105b857600080fd5b506105c16109c0565b6040518080602001828103825283818151815260200191508051906020019080838360005b838110156106015780820151818401526020810190506105e6565b50505050905090810190601f16801561062e5780820380516001836020036101000a031916815260200191505b509250505060405180910390f35b34801561064857600080fd5b50610651610ac4565b6040518082815260200191505060405180910390f35b34801561067357600080fd5b506106b66004803603602081101561068a57600080fd5b81019080803573ffffffffffffffffffffffffffffffffffffffff169060200190929190505050610add565b6040518082815260200191505060405180910390f35b3480156106d857600080fd5b506106e1610afe565b60405180827dffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff19167dffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff1916815260200191505060405180910390f35b60008060038081111561074857fe5b90508091505090565b6000819050919050565b6000806000806003905060006404a817c80090506000660aa87bee5380009050828282955095509550505050909192565b8073ffffffffffffffffffffffffffffffffffffffff166108fc349081150290604051600060405180830381858888f193505050501580156107d2573d6000803e3d6000fd5b5050565b60008073ca35b7d915458ef540ade6068dfe2f44e8fa733c90508091505090565b600080805460018160011615610100020316600290049050905090565b600080600690508091505090565b600061082e6001610751565b905090565b6000806000807f01f40000000000000000000000000000000000000000000000000000000000009050808160006002811061086a57fe5b1a60f81b8260016002811061087b57fe5b1a60f81b93509350935050909192565b600080600033348473ffffffffffffffffffffffffffffffffffffffff166108fc349081150290604051600060405180830381858888f19350505050925092509250909192565b60003073ffffffffffffffffffffffffffffffffffffffff1631905090565b606060008054600181600116156101000203166002900480601f0160208091040260200160405190810160405280929190818152602001828054600181600116156101000203166002900480156109895780601f1061095e57610100808354040283529160200191610989565b820191906000526020600020905b81548152906001019060200180831161096c57829003601f168201915b5050505050905090565b6000807fc80000000000000000000000000000000000000000000000000000000000000090508091505090565b60608060008054600181600116156101000203166002900480601f016020809104026020016040519081016040528092919081815260200182805460018160011615610100020316600290048015610a595780601f10610a2e57610100808354040283529160200191610a59565b820191906000526020600020905b815481529060010190602001808311610a3c57829003601f168201915b505050505090507f610000000000000000000000000000000000000000000000000000000000000081600081518110610a8e57fe5b60200101907effffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff1916908160001a9053508091505090565b60008060006003811115610ad457fe5b90508091505090565b60008173ffffffffffffffffffffffffffffffffffffffff16319050919050565b6000807f01f40000000000000000000000000000000000000000000000000000000000009050809150509056fea265627a7a72315820b1fc123686af3b8b58153669574b0f7e72b49abb998b93648020d05c5c97c3ef64736f6c634300050d0032";

    public static final String FUNC_HEXLITERAL3 = "HexLiteral3";

    public static final String FUNC_GETADDRESS = "getAddress";

    public static final String FUNC_GETBALANCE = "getBalance";

    public static final String FUNC_GETCURRENTBALANCE = "getCurrentBalance";

    public static final String FUNC_GETHEXLITERAL = "getHexLiteral";

    public static final String FUNC_GETHEXLITERALBYTES = "getHexLiteralBytes";

    public static final String FUNC_GETINT = "getInt";

    public static final String FUNC_GETSEASON1 = "getSeason1";

    public static final String FUNC_GETSEASON2 = "getSeason2";

    public static final String FUNC_GETSEASONINDEX = "getSeasonIndex";

    public static final String FUNC_GETSTR1 = "getStr1";

    public static final String FUNC_GETSTR1LENGTH = "getStr1Length";

    public static final String FUNC_GETVALUE = "getValue";

    public static final String FUNC_GOSEND = "goSend";

    public static final String FUNC_GOTRANSFER = "goTransfer";

    public static final String FUNC_PRINTSEASON = "printSeason";

    public static final String FUNC_SETSTR1 = "setStr1";

    @Deprecated
    protected BasicDataTypeContract2(String contractAddress, Web3j web3j, Credentials credentials, BigInteger gasPrice, BigInteger gasLimit) {
        super(BINARY, contractAddress, web3j, credentials, gasPrice, gasLimit);
    }

    protected BasicDataTypeContract2(String contractAddress, Web3j web3j, Credentials credentials, GasProvider contractGasProvider) {
        super(BINARY, contractAddress, web3j, credentials, contractGasProvider);
    }

    @Deprecated
    protected BasicDataTypeContract2(String contractAddress, Web3j web3j, TransactionManager transactionManager, BigInteger gasPrice, BigInteger gasLimit) {
        super(BINARY, contractAddress, web3j, transactionManager, gasPrice, gasLimit);
    }

    protected BasicDataTypeContract2(String contractAddress, Web3j web3j, TransactionManager transactionManager, GasProvider contractGasProvider) {
        super(BINARY, contractAddress, web3j, transactionManager, contractGasProvider);
    }

    public RemoteCall<byte[]> HexLiteral3() {
        final Function function = new Function(FUNC_HEXLITERAL3, 
                Arrays.<Type>asList(), 
                Arrays.<TypeReference<?>>asList(new TypeReference<Bytes2>() {}));
        return executeRemoteCallSingleValueReturn(function, byte[].class);
    }

    public RemoteCall<TransactionReceipt> getAddress() {
        final Function function = new Function(
                FUNC_GETADDRESS, 
                Arrays.<Type>asList(), 
                Collections.<TypeReference<?>>emptyList());
        return executeRemoteCallTransaction(function);
    }

    public RemoteCall<BigInteger> getBalance(String addr) {
        final Function function = new Function(FUNC_GETBALANCE, 
                Arrays.<Type>asList(new org.web3j.abi.datatypes.Address(addr)), 
                Arrays.<TypeReference<?>>asList(new TypeReference<Uint256>() {}));
        return executeRemoteCallSingleValueReturn(function, BigInteger.class);
    }

    public RemoteCall<BigInteger> getCurrentBalance() {
        final Function function = new Function(FUNC_GETCURRENTBALANCE, 
                Arrays.<Type>asList(), 
                Arrays.<TypeReference<?>>asList(new TypeReference<Uint256>() {}));
        return executeRemoteCallSingleValueReturn(function, BigInteger.class);
    }

    public RemoteCall<byte[]> getHexLiteral() {
        final Function function = new Function(FUNC_GETHEXLITERAL, 
                Arrays.<Type>asList(), 
                Arrays.<TypeReference<?>>asList(new TypeReference<Bytes1>() {}));
        return executeRemoteCallSingleValueReturn(function, byte[].class);
    }

    public RemoteCall<Tuple3<byte[], byte[], byte[]>> getHexLiteralBytes() {
        final Function function = new Function(FUNC_GETHEXLITERALBYTES, 
                Arrays.<Type>asList(), 
                Arrays.<TypeReference<?>>asList(new TypeReference<Bytes2>() {}, new TypeReference<Bytes1>() {}, new TypeReference<Bytes1>() {}));
        return new RemoteCall<Tuple3<byte[], byte[], byte[]>>(
                new Callable<Tuple3<byte[], byte[], byte[]>>() {
                    @Override
                    public Tuple3<byte[], byte[], byte[]> call() throws Exception {
                        List<Type> results = executeCallMultipleValueReturn(function);
                        return new Tuple3<byte[], byte[], byte[]>(
                                (byte[]) results.get(0).getValue(), 
                                (byte[]) results.get(1).getValue(), 
                                (byte[]) results.get(2).getValue());
                    }
                });
    }

    public RemoteCall<TransactionReceipt> getInt() {
        final Function function = new Function(
                FUNC_GETINT, 
                Arrays.<Type>asList(), 
                Collections.<TypeReference<?>>emptyList());
        return executeRemoteCallTransaction(function);
    }

    public RemoteCall<BigInteger> getSeason1() {
        final Function function = new Function(FUNC_GETSEASON1, 
                Arrays.<Type>asList(), 
                Arrays.<TypeReference<?>>asList(new TypeReference<Uint8>() {}));
        return executeRemoteCallSingleValueReturn(function, BigInteger.class);
    }

    public RemoteCall<BigInteger> getSeason2() {
        final Function function = new Function(FUNC_GETSEASON2, 
                Arrays.<Type>asList(), 
                Arrays.<TypeReference<?>>asList(new TypeReference<Uint8>() {}));
        return executeRemoteCallSingleValueReturn(function, BigInteger.class);
    }

    public RemoteCall<BigInteger> getSeasonIndex() {
        final Function function = new Function(FUNC_GETSEASONINDEX, 
                Arrays.<Type>asList(), 
                Arrays.<TypeReference<?>>asList(new TypeReference<Uint256>() {}));
        return executeRemoteCallSingleValueReturn(function, BigInteger.class);
    }

    public RemoteCall<String> getStr1() {
        final Function function = new Function(FUNC_GETSTR1, 
                Arrays.<Type>asList(), 
                Arrays.<TypeReference<?>>asList(new TypeReference<Utf8String>() {}));
        return executeRemoteCallSingleValueReturn(function, String.class);
    }

    public RemoteCall<BigInteger> getStr1Length() {
        final Function function = new Function(FUNC_GETSTR1LENGTH, 
                Arrays.<Type>asList(), 
                Arrays.<TypeReference<?>>asList(new TypeReference<Uint256>() {}));
        return executeRemoteCallSingleValueReturn(function, BigInteger.class);
    }

    public RemoteCall<Tuple3<BigInteger, BigInteger, BigInteger>> getValue() {
        final Function function = new Function(FUNC_GETVALUE, 
                Arrays.<Type>asList(), 
                Arrays.<TypeReference<?>>asList(new TypeReference<Uint128>() {}, new TypeReference<Uint256>() {}, new TypeReference<Uint256>() {}));
        return new RemoteCall<Tuple3<BigInteger, BigInteger, BigInteger>>(
                new Callable<Tuple3<BigInteger, BigInteger, BigInteger>>() {
                    @Override
                    public Tuple3<BigInteger, BigInteger, BigInteger> call() throws Exception {
                        List<Type> results = executeCallMultipleValueReturn(function);
                        return new Tuple3<BigInteger, BigInteger, BigInteger>(
                                (BigInteger) results.get(0).getValue(), 
                                (BigInteger) results.get(1).getValue(), 
                                (BigInteger) results.get(2).getValue());
                    }
                });
    }

    public RemoteCall<TransactionReceipt> goSend(BigInteger weiValue) {
        final Function function = new Function(
                FUNC_GOSEND, 
                Arrays.<Type>asList(), 
                Collections.<TypeReference<?>>emptyList());
        return executeRemoteCallTransaction(function, weiValue);
    }

    public RemoteCall<TransactionReceipt> goTransfer(String addr, BigInteger weiValue) {
        final Function function = new Function(
                FUNC_GOTRANSFER, 
                Arrays.<Type>asList(new org.web3j.abi.datatypes.Address(addr)), 
                Collections.<TypeReference<?>>emptyList());
        return executeRemoteCallTransaction(function, weiValue);
    }

    public RemoteCall<BigInteger> printSeason(BigInteger s) {
        final Function function = new Function(FUNC_PRINTSEASON, 
                Arrays.<Type>asList(new org.web3j.abi.datatypes.generated.Uint8(s)), 
                Arrays.<TypeReference<?>>asList(new TypeReference<Uint8>() {}));
        return executeRemoteCallSingleValueReturn(function, BigInteger.class);
    }

    public RemoteCall<String> setStr1() {
        final Function function = new Function(FUNC_SETSTR1, 
                Arrays.<Type>asList(), 
                Arrays.<TypeReference<?>>asList(new TypeReference<Utf8String>() {}));
        return executeRemoteCallSingleValueReturn(function, String.class);
    }

    public static RemoteCall<BasicDataTypeContract2> deploy(Web3j web3j, Credentials credentials, GasProvider contractGasProvider) {
        return deployRemoteCall(BasicDataTypeContract2.class, web3j, credentials, contractGasProvider, BINARY, "");
    }

    @Deprecated
    public static RemoteCall<BasicDataTypeContract2> deploy(Web3j web3j, Credentials credentials, BigInteger gasPrice, BigInteger gasLimit) {
        return deployRemoteCall(BasicDataTypeContract2.class, web3j, credentials, gasPrice, gasLimit, BINARY, "");
    }

    public static RemoteCall<BasicDataTypeContract2> deploy(Web3j web3j, TransactionManager transactionManager, GasProvider contractGasProvider) {
        return deployRemoteCall(BasicDataTypeContract2.class, web3j, transactionManager, contractGasProvider, BINARY, "");
    }

    @Deprecated
    public static RemoteCall<BasicDataTypeContract2> deploy(Web3j web3j, TransactionManager transactionManager, BigInteger gasPrice, BigInteger gasLimit) {
        return deployRemoteCall(BasicDataTypeContract2.class, web3j, transactionManager, gasPrice, gasLimit, BINARY, "");
    }

    @Deprecated
    public static BasicDataTypeContract2 load(String contractAddress, Web3j web3j, Credentials credentials, BigInteger gasPrice, BigInteger gasLimit) {
        return new BasicDataTypeContract2(contractAddress, web3j, credentials, gasPrice, gasLimit);
    }

    @Deprecated
    public static BasicDataTypeContract2 load(String contractAddress, Web3j web3j, TransactionManager transactionManager, BigInteger gasPrice, BigInteger gasLimit) {
        return new BasicDataTypeContract2(contractAddress, web3j, transactionManager, gasPrice, gasLimit);
    }

    public static BasicDataTypeContract2 load(String contractAddress, Web3j web3j, Credentials credentials, GasProvider contractGasProvider) {
        return new BasicDataTypeContract2(contractAddress, web3j, credentials, contractGasProvider);
    }

    public static BasicDataTypeContract2 load(String contractAddress, Web3j web3j, TransactionManager transactionManager, GasProvider contractGasProvider) {
        return new BasicDataTypeContract2(contractAddress, web3j, transactionManager, contractGasProvider);
    }
}
