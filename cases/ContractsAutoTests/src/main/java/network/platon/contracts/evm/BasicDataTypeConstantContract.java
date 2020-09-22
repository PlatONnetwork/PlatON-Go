package network.platon.contracts.evm;

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
 * <p>Please use the <a href="https://github.com/PlatONnetwork/client-sdk-java/releases">platon-web3j command line tools</a>,
 * or the org.web3j.codegen.SolidityFunctionWrapperGenerator in the 
 * <a href="https://github.com/PlatONnetwork/client-sdk-java/tree/master/codegen">codegen module</a> to update.
 *
 * <p>Generated with web3j version 0.13.1.5.
 */
public class BasicDataTypeConstantContract extends Contract {
    private static final String BINARY = "60806040526040805190810160405280600581526020017f68656c6c6f0000000000000000000000000000000000000000000000000000008152506000908051906020019061004f9291906100ae565b506040805190810160405280600581526020017f776f726c640000000000000000000000000000000000000000000000000000008152506001908051906020019061009b9291906100ae565b503480156100a857600080fd5b50610153565b828054600181600116156101000203166002900490600052602060002090601f016020900481019282601f106100ef57805160ff191683800117855561011d565b8280016001018555821561011d579182015b8281111561011c578251825591602001919060010190610101565b5b50905061012a919061012e565b5090565b61015091905b8082111561014c576000816000905550600101610134565b5090565b90565b610bbb806101626000396000f3fe60806040526004361061011b576000357c01000000000000000000000000000000000000000000000000000000009004806362738998116100b2578063a574971011610081578063a5749710146105a3578063a650b683146105ce578063ccb344151461065e578063ea07cdfa146106cd578063f8b2cb4f146106f85761011b565b8063627389981461038257806369c76917146103ad5780637c66959d1461043d578063833f17d5146105405761011b565b806338cc4831116100ee57806338cc48311461025a57806343e33562146102b15780634dbe1d0f1461031e578063515899f0146103575761011b565b8063166aa6e614610120578063209652551461018057806332c7a283146101dd57806338023fb914610216575b600080fd5b34801561012c57600080fd5b5061015c6004803603602081101561014357600080fd5b81019080803560ff16906020019092919050505061075d565b6040518082600381111561016c57fe5b60ff16815260200191505060405180910390f35b34801561018c57600080fd5b50610195610767565b60405180846fffffffffffffffffffffffffffffffff166fffffffffffffffffffffffffffffffff168152602001838152602001828152602001935050505060405180910390f35b3480156101e957600080fd5b506101f2610798565b6040518082600381111561020257fe5b60ff16815260200191505060405180910390f35b6102586004803603602081101561022c57600080fd5b81019080803573ffffffffffffffffffffffffffffffffffffffff1690602001909291905050506107b0565b005b34801561026657600080fd5b5061026f6107fa565b604051808273ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200191505060405180910390f35b3480156102bd57600080fd5b506102c661081b565b60405180827dffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff19167dffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff1916815260200191505060405180910390f35b34801561032a57600080fd5b50610333610848565b6040518082600381111561034357fe5b60ff16815260200191505060405180910390f35b34801561036357600080fd5b5061036c610859565b6040518082815260200191505060405180910390f35b34801561038e57600080fd5b50610397610876565b6040518082815260200191505060405180910390f35b3480156103b957600080fd5b506103c2610884565b6040518080602001828103825283818151815260200191508051906020019080838360005b838110156104025780820151818401526020810190506103e7565b50505050905090810190601f16801561042f5780820380516001836020036101000a031916815260200191505b509250505060405180910390f35b34801561044957600080fd5b5061045261098b565b60405180847dffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff19167dffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff19168152602001837effffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff19167effffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff19168152602001827effffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff19167effffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff19168152602001935050505060405180910390f35b6105826004803603602081101561055657600080fd5b81019080803573ffffffffffffffffffffffffffffffffffffffff169060200190929190505050610a25565b60405180838152602001821515151581526020019250505060405180910390f35b3480156105af57600080fd5b506105b8610a67565b6040518082815260200191505060405180910390f35b3480156105da57600080fd5b506105e3610a86565b6040518080602001828103825283818151815260200191508051906020019080838360005b83811015610623578082015181840152602081019050610608565b50505050905090810190601f1680156106505780820380516001836020036101000a031916815260200191505b509250505060405180910390f35b34801561066a57600080fd5b50610673610b28565b60405180827effffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff19167effffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff1916815260200191505060405180910390f35b3480156106d957600080fd5b506106e2610b55565b6040518082815260200191505060405180910390f35b34801561070457600080fd5b506107476004803603602081101561071b57600080fd5b81019080803573ffffffffffffffffffffffffffffffffffffffff169060200190929190505050610b6e565b6040518082815260200191505060405180910390f35b6000819050919050565b6000806000806003905060006404a817c80090506000660aa87bee5380009050828282955095509550505050909192565b6000806003808111156107a757fe5b90508091505090565b8073ffffffffffffffffffffffffffffffffffffffff166108fc349081150290604051600060405180830381858888f193505050501580156107f6573d6000803e3d6000fd5b5050565b6000807372ad2b713faa14c2c4cd2d7affe5d8f538968f5a90508091505090565b6000807f01f400000000000000000000000000000000000000000000000000000000000090508091505090565b6000610854600161075d565b905090565b600080805460018160011615610100020316600290049050905090565b600080600690508091505090565b60608060008054600181600116156101000203166002900480601f01602080910402602001604051908101604052809291908181526020018280546001816001161561010002031660029004801561091d5780601f106108f25761010080835404028352916020019161091d565b820191906000526020600020905b81548152906001019060200180831161090057829003601f168201915b505050505090507f610000000000000000000000000000000000000000000000000000000000000081600081518110151561095457fe5b9060200101907effffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff1916908160001a9053508091505090565b6000806000807f01f40000000000000000000000000000000000000000000000000000000000009050808160006002811015156109c457fe5b1a7f0100000000000000000000000000000000000000000000000000000000000000028260016002811015156109f657fe5b1a7f01000000000000000000000000000000000000000000000000000000000000000293509350935050909192565b600080348373ffffffffffffffffffffffffffffffffffffffff166108fc349081150290604051600060405180830381858888f1935050505091509150915091565b60003073ffffffffffffffffffffffffffffffffffffffff1631905090565b606060008054600181600116156101000203166002900480601f016020809104026020016040519081016040528092919081815260200182805460018160011615610100020316600290048015610b1e5780601f10610af357610100808354040283529160200191610b1e565b820191906000526020600020905b815481529060010190602001808311610b0157829003601f168201915b5050505050905090565b6000807fc80000000000000000000000000000000000000000000000000000000000000090508091505090565b60008060006003811115610b6557fe5b90508091505090565b60008173ffffffffffffffffffffffffffffffffffffffff1631905091905056fea165627a7a72305820ef6110aa4b9bee27cc59f5322af06bc276e35776f868cac109ee9db260cfeb660029";

    public static final String FUNC_PRINTSEASON = "printSeason";

    public static final String FUNC_GETVALUE = "getValue";

    public static final String FUNC_GETSEASONB = "getSeasonB";

    public static final String FUNC_GOTRANSFER = "goTransfer";

    public static final String FUNC_GETADDRESS = "getAddress";

    public static final String FUNC_GETHEXLITERAB = "getHexLiteraB";

    public static final String FUNC_GETSEASONA = "getSeasonA";

    public static final String FUNC_GETSTRALENGTH = "getStrALength";

    public static final String FUNC_GETINT = "getInt";

    public static final String FUNC_SETSTRA = "setStrA";

    public static final String FUNC_GETHEXLITERAC = "getHexLiteraC";

    public static final String FUNC_GOSEND = "goSend";

    public static final String FUNC_GETCURRENTBALANCE = "getCurrentBalance";

    public static final String FUNC_GETSTRA = "getStrA";

    public static final String FUNC_GETHEXLITERAA = "getHexLiteraA";

    public static final String FUNC_GETSEASONINDEX = "getSeasonIndex";

    public static final String FUNC_GETBALANCE = "getBalance";

    protected BasicDataTypeConstantContract(String contractAddress, Web3j web3j, Credentials credentials, GasProvider contractGasProvider, Long chainId) {
        super(BINARY, contractAddress, web3j, credentials, contractGasProvider, chainId);
    }

    protected BasicDataTypeConstantContract(String contractAddress, Web3j web3j, TransactionManager transactionManager, GasProvider contractGasProvider, Long chainId) {
        super(BINARY, contractAddress, web3j, transactionManager, contractGasProvider, chainId);
    }

    public RemoteCall<BigInteger> printSeason(BigInteger s) {
        final Function function = new Function(FUNC_PRINTSEASON, 
                Arrays.<Type>asList(new Uint8(s)),
                Arrays.<TypeReference<?>>asList(new TypeReference<Uint8>() {}));
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

    public RemoteCall<BigInteger> getSeasonB() {
        final Function function = new Function(FUNC_GETSEASONB, 
                Arrays.<Type>asList(), 
                Arrays.<TypeReference<?>>asList(new TypeReference<Uint8>() {}));
        return executeRemoteCallSingleValueReturn(function, BigInteger.class);
    }

    public RemoteCall<TransactionReceipt> goTransfer(String addr, BigInteger vonValue) {
        final Function function = new Function(
                FUNC_GOTRANSFER, 
                Arrays.<Type>asList(new org.web3j.abi.datatypes.Address(addr)), 
                Collections.<TypeReference<?>>emptyList());
        return executeRemoteCallTransaction(function, vonValue);
    }

    public RemoteCall<TransactionReceipt> getAddress() {
        final Function function = new Function(
                FUNC_GETADDRESS, 
                Arrays.<Type>asList(), 
                Collections.<TypeReference<?>>emptyList());
        return executeRemoteCallTransaction(function);
    }

    public RemoteCall<byte[]> getHexLiteraB() {
        final Function function = new Function(FUNC_GETHEXLITERAB, 
                Arrays.<Type>asList(), 
                Arrays.<TypeReference<?>>asList(new TypeReference<Bytes2>() {}));
        return executeRemoteCallSingleValueReturn(function, byte[].class);
    }

    public RemoteCall<BigInteger> getSeasonA() {
        final Function function = new Function(FUNC_GETSEASONA, 
                Arrays.<Type>asList(), 
                Arrays.<TypeReference<?>>asList(new TypeReference<Uint8>() {}));
        return executeRemoteCallSingleValueReturn(function, BigInteger.class);
    }

    public RemoteCall<BigInteger> getStrALength() {
        final Function function = new Function(FUNC_GETSTRALENGTH, 
                Arrays.<Type>asList(), 
                Arrays.<TypeReference<?>>asList(new TypeReference<Uint256>() {}));
        return executeRemoteCallSingleValueReturn(function, BigInteger.class);
    }

    public RemoteCall<TransactionReceipt> getInt() {
        final Function function = new Function(
                FUNC_GETINT, 
                Arrays.<Type>asList(), 
                Collections.<TypeReference<?>>emptyList());
        return executeRemoteCallTransaction(function);
    }

    public RemoteCall<String> setStrA() {
        final Function function = new Function(FUNC_SETSTRA, 
                Arrays.<Type>asList(), 
                Arrays.<TypeReference<?>>asList(new TypeReference<Utf8String>() {}));
        return executeRemoteCallSingleValueReturn(function, String.class);
    }

    public RemoteCall<Tuple3<byte[], byte[], byte[]>> getHexLiteraC() {
        final Function function = new Function(FUNC_GETHEXLITERAC, 
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

    public RemoteCall<TransactionReceipt> goSend(String addr, BigInteger vonValue) {
        final Function function = new Function(
                FUNC_GOSEND, 
                Arrays.<Type>asList(new org.web3j.abi.datatypes.Address(addr)), 
                Collections.<TypeReference<?>>emptyList());
        return executeRemoteCallTransaction(function, vonValue);
    }

    public RemoteCall<BigInteger> getCurrentBalance() {
        final Function function = new Function(FUNC_GETCURRENTBALANCE, 
                Arrays.<Type>asList(), 
                Arrays.<TypeReference<?>>asList(new TypeReference<Uint256>() {}));
        return executeRemoteCallSingleValueReturn(function, BigInteger.class);
    }

    public RemoteCall<String> getStrA() {
        final Function function = new Function(FUNC_GETSTRA, 
                Arrays.<Type>asList(), 
                Arrays.<TypeReference<?>>asList(new TypeReference<Utf8String>() {}));
        return executeRemoteCallSingleValueReturn(function, String.class);
    }

    public RemoteCall<byte[]> getHexLiteraA() {
        final Function function = new Function(FUNC_GETHEXLITERAA, 
                Arrays.<Type>asList(), 
                Arrays.<TypeReference<?>>asList(new TypeReference<Bytes1>() {}));
        return executeRemoteCallSingleValueReturn(function, byte[].class);
    }

    public RemoteCall<BigInteger> getSeasonIndex() {
        final Function function = new Function(FUNC_GETSEASONINDEX, 
                Arrays.<Type>asList(), 
                Arrays.<TypeReference<?>>asList(new TypeReference<Uint256>() {}));
        return executeRemoteCallSingleValueReturn(function, BigInteger.class);
    }

    public RemoteCall<BigInteger> getBalance(String addr) {
        final Function function = new Function(FUNC_GETBALANCE, 
                Arrays.<Type>asList(new org.web3j.abi.datatypes.Address(addr)), 
                Arrays.<TypeReference<?>>asList(new TypeReference<Uint256>() {}));
        return executeRemoteCallSingleValueReturn(function, BigInteger.class);
    }

    public static RemoteCall<BasicDataTypeConstantContract> deploy(Web3j web3j, Credentials credentials, GasProvider contractGasProvider, Long chainId) {
        return deployRemoteCall(BasicDataTypeConstantContract.class, web3j, credentials, contractGasProvider, BINARY,  "", chainId);
    }

    public static RemoteCall<BasicDataTypeConstantContract> deploy(Web3j web3j, TransactionManager transactionManager, GasProvider contractGasProvider, Long chainId) {
        return deployRemoteCall(BasicDataTypeConstantContract.class, web3j, transactionManager, contractGasProvider, BINARY,  "", chainId);
    }

    public static BasicDataTypeConstantContract load(String contractAddress, Web3j web3j, Credentials credentials, GasProvider contractGasProvider, Long chainId) {
        return new BasicDataTypeConstantContract(contractAddress, web3j, credentials, contractGasProvider, chainId);
    }

    public static BasicDataTypeConstantContract load(String contractAddress, Web3j web3j, TransactionManager transactionManager, GasProvider contractGasProvider, Long chainId) {
        return new BasicDataTypeConstantContract(contractAddress, web3j, transactionManager, contractGasProvider, chainId);
    }
}
