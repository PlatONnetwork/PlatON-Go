package network.platon.contracts.evm;

import java.math.BigInteger;
import java.util.Arrays;
import java.util.Collections;
import java.util.List;
import java.util.concurrent.Callable;
import org.web3j.abi.TypeReference;
import org.web3j.abi.datatypes.Address;
import org.web3j.abi.datatypes.Bool;
import org.web3j.abi.datatypes.DynamicBytes;
import org.web3j.abi.datatypes.Function;
import org.web3j.abi.datatypes.Type;
import org.web3j.abi.datatypes.Utf8String;
import org.web3j.abi.datatypes.generated.Bytes1;
import org.web3j.abi.datatypes.generated.Int256;
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
public class MappingDataTypeContract extends Contract {
    private static final String BINARY = "608060405260405180606001604052806040518060400160405280600481526020017f4c7563790000000000000000000000000000000000000000000000000000000081525081526020016040518060400160405280600481526020017f456c6c610000000000000000000000000000000000000000000000000000000081525081526020016040518060400160405280600481526020017f4c696c7900000000000000000000000000000000000000000000000000000000815250815250600d906003620000d0929190620000e5565b50348015620000de57600080fd5b5062000278565b82805482825590600052602060002090810192821562000139579160200282015b8281111562000138578251829080519060200190620001279291906200014c565b509160200191906001019062000106565b5b509050620001489190620001d3565b5090565b828054600181600116156101000203166002900490600052602060002090601f016020900481019282601f106200018f57805160ff1916838001178555620001c0565b82800160010185558215620001c0579182015b82811115620001bf578251825591602001919060010190620001a2565b5b509050620001cf919062000204565b5090565b6200020191905b80821115620001fd5760008181620001f391906200022c565b50600101620001da565b5090565b90565b6200022991905b80821115620002255760008160009055506001016200020b565b5090565b90565b50805460018160011615610100020316600290046000825580601f1062000254575062000275565b601f01602090049060005260206000209081019062000274919062000204565b5b50565b610c7c80620002886000396000f3fe608060405234801561001057600080fd5b50600436106100b35760003560e01c806387c9644b1161007157806387c9644b1461047c578063aad39512146104c2578063bbf722a21461058b578063bebf70f3146105f9578063f480ae4914610653578063fdf8b00814610722576100b3565b806252f6ee146100b8578063566ab1fd146101875780636257a397146101d75780636b8ff5741461029457806376fc59c31461033b57806382bc1149146103c3575b600080fd5b610171600480360360208110156100ce57600080fd5b81019080803590602001906401000000008111156100eb57600080fd5b8201836020820111156100fd57600080fd5b8035906020019184600183028401116401000000008311171561011f57600080fd5b91908080601f016020809104026020016040519081016040528093929190818152602001838380828437600081840152601f19601f82011690508083019250505050505050919291929050505061072c565b6040518082815260200191505060405180910390f35b6101b36004803603602081101561019d57600080fd5b810190808035906020019092919050505061075a565b604051808260028111156101c357fe5b60ff16815260200191505060405180910390f35b610219600480360360208110156101ed57600080fd5b81019080803573ffffffffffffffffffffffffffffffffffffffff16906020019092919050505061077a565b6040518080602001828103825283818151815260200191508051906020019080838360005b8381101561025957808201518184015260208101905061023e565b50505050905090810190601f1680156102865780820380516001836020036101000a031916815260200191505b509250505060405180910390f35b6102c0600480360360208110156102aa57600080fd5b810190808035906020019092919050505061082a565b6040518080602001828103825283818151815260200191508051906020019080838360005b838110156103005780820151818401526020810190506102e5565b50505050905090810190601f16801561032d5780820380516001836020036101000a031916815260200191505b509250505060405180910390f35b6103696004803603602081101561035157600080fd5b810190808035151590602001909291905050506108df565b60405180827effffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff19167effffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff1916815260200191505060405180910390f35b6103ef600480360360208110156103d957600080fd5b81019080803590602001909291905050506108ff565b604051808481526020018060200183151515158152602001828103825284818151815260200191508051906020019080838360005b8381101561043f578082015181840152602081019050610424565b50505050905090810190601f16801561046c5780820380516001836020036101000a031916815260200191505b5094505050505060405180910390f35b6104a86004803603602081101561049257600080fd5b81019080803590602001909291905050506109ce565b604051808215151515815260200191505060405180910390f35b610510600480360360208110156104d857600080fd5b8101908080357effffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff191690602001909291905050506109ee565b6040518080602001828103825283818151815260200191508051906020019080838360005b83811015610550578082015181840152602081019050610535565b50505050905090810190601f16801561057d5780820380516001836020036101000a031916815260200191505b509250505060405180910390f35b6105b7600480360360208110156105a157600080fd5b8101908080359060200190929190505050610a9e565b604051808273ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200191505060405180910390f35b61062f6004803603604081101561060f57600080fd5b810190808035906020019092919080359060200190929190505050610ad1565b6040518082600281111561063f57fe5b60ff16815260200191505060405180910390f35b61070c6004803603602081101561066957600080fd5b810190808035906020019064010000000081111561068657600080fd5b82018360208201111561069857600080fd5b803590602001918460018302840111640100000000831117156106ba57600080fd5b91908080601f016020809104026020016040519081016040528093929190818152602001838380828437600081840152601f19601f820116905080830192505050505050509192919290505050610b00565b6040518082815260200191505060405180910390f35b61072a610b2e565b005b6006818051602081018201805184825260208301602085012081835280955050505050506000915090505481565b600b6020528060005260406000206000915054906101000a900460ff1681565b60086020528060005260406000206000915090508054600181600116156101000203166002900480601f0160208091040260200160405190810160405280929190818152602001828054600181600116156101000203166002900480156108225780601f106107f757610100808354040283529160200191610822565b820191906000526020600020905b81548152906001019060200180831161080557829003601f168201915b505050505081565b6060600c60008381526020019081526020016000208054600181600116156101000203166002900480601f0160208091040260200160405190810160405280929190818152602001828054600181600116156101000203166002900480156108d35780601f106108a8576101008083540402835291602001916108d3565b820191906000526020600020905b8154815290600101906020018083116108b657829003601f168201915b50505050509050919050565b60046020528060005260406000206000915054906101000a900460f81b81565b600a602052806000526040600020600091509050806000015490806001018054600181600116156101000203166002900480601f0160208091040260200160405190810160405280929190818152602001828054600181600116156101000203166002900480156109b15780601f10610986576101008083540402835291602001916109b1565b820191906000526020600020905b81548152906001019060200180831161099457829003601f168201915b5050505050908060020160009054906101000a900460ff16905083565b60036020528060005260406000206000915054906101000a900460ff1681565b60056020528060005260406000206000915090508054600181600116156101000203166002900480601f016020809104026020016040519081016040528092919081815260200182805460018160011615610100020316600290048015610a965780601f10610a6b57610100808354040283529160200191610a96565b820191906000526020600020905b815481529060010190602001808311610a7957829003601f168201915b505050505081565b60026020528060005260406000206000915054906101000a900473ffffffffffffffffffffffffffffffffffffffff1681565b60096020528160005260406000206020528060005260406000206000915091509054906101000a900460ff1681565b6007818051602081018201805184825260208301602085012081835280955050505050506000915090505481565b60008090505b600d80549050811015610b9857600d8181548110610b4e57fe5b90600052602060002001600c60008381526020019081526020016000209080546001816001161561010002031660029004610b8a929190610b9b565b508080600101915050610b34565b50565b828054600181600116156101000203166002900490600052602060002090601f016020900481019282601f10610bd45780548555610c11565b82800160010185558215610c1157600052602060002091601f016020900482015b82811115610c10578254825591600101919060010190610bf5565b5b509050610c1e9190610c22565b5090565b610c4491905b80821115610c40576000816000905550600101610c28565b5090565b9056fea265627a7a72315820aa48fec572535dda66b9c78b9d387d1a56f84b8c5f566c7b5f35aecda746907864736f6c634300050d0032";

    public static final String FUNC_SIZEENUMMAP = "SizeEnumMap";

    public static final String FUNC_ADDNAME = "addName";

    public static final String FUNC_ADDRESSMAP = "addressMap";

    public static final String FUNC_BOOLMAP = "boolMap";

    public static final String FUNC_BYTEMAP = "byteMap";

    public static final String FUNC_BYTESMAP = "bytesMap";

    public static final String FUNC_GETNAME = "getName";

    public static final String FUNC_INTMAP = "intMap";

    public static final String FUNC_PEOPLEMAP = "peopleMap";

    public static final String FUNC_SIZEMAP = "sizeMap";

    public static final String FUNC_STRINGMAP = "stringMap";

    public static final String FUNC_UINTMAP = "uintMap";

    protected MappingDataTypeContract(String contractAddress, Web3j web3j, Credentials credentials, GasProvider contractGasProvider, Long chainId) {
        super(BINARY, contractAddress, web3j, credentials, contractGasProvider, chainId);
    }

    protected MappingDataTypeContract(String contractAddress, Web3j web3j, TransactionManager transactionManager, GasProvider contractGasProvider, Long chainId) {
        super(BINARY, contractAddress, web3j, transactionManager, contractGasProvider, chainId);
    }

    public RemoteCall<BigInteger> SizeEnumMap(BigInteger param0) {
        final Function function = new Function(FUNC_SIZEENUMMAP, 
                Arrays.<Type>asList(new Int256(param0)),
                Arrays.<TypeReference<?>>asList(new TypeReference<Uint8>() {}));
        return executeRemoteCallSingleValueReturn(function, BigInteger.class);
    }

    public RemoteCall<TransactionReceipt> addName() {
        final Function function = new Function(
                FUNC_ADDNAME, 
                Arrays.<Type>asList(), 
                Collections.<TypeReference<?>>emptyList());
        return executeRemoteCallTransaction(function);
    }

    public RemoteCall<String> addressMap(BigInteger param0) {
        final Function function = new Function(FUNC_ADDRESSMAP, 
                Arrays.<Type>asList(new Uint256(param0)),
                Arrays.<TypeReference<?>>asList(new TypeReference<Address>() {}));
        return executeRemoteCallSingleValueReturn(function, String.class);
    }

    public RemoteCall<Boolean> boolMap(BigInteger param0) {
        final Function function = new Function(FUNC_BOOLMAP, 
                Arrays.<Type>asList(new Int256(param0)),
                Arrays.<TypeReference<?>>asList(new TypeReference<Bool>() {}));
        return executeRemoteCallSingleValueReturn(function, Boolean.class);
    }

    public RemoteCall<byte[]> byteMap(Boolean param0) {
        final Function function = new Function(FUNC_BYTEMAP, 
                Arrays.<Type>asList(new Bool(param0)),
                Arrays.<TypeReference<?>>asList(new TypeReference<Bytes1>() {}));
        return executeRemoteCallSingleValueReturn(function, byte[].class);
    }

    public RemoteCall<byte[]> bytesMap(String param0) {
        final Function function = new Function(FUNC_BYTESMAP, 
                Arrays.<Type>asList(new Address(param0)),
                Arrays.<TypeReference<?>>asList(new TypeReference<DynamicBytes>() {}));
        return executeRemoteCallSingleValueReturn(function, byte[].class);
    }

    public RemoteCall<String> getName(BigInteger index) {
        final Function function = new Function(FUNC_GETNAME, 
                Arrays.<Type>asList(new Uint256(index)),
                Arrays.<TypeReference<?>>asList(new TypeReference<Utf8String>() {}));
        return executeRemoteCallSingleValueReturn(function, String.class);
    }

    public RemoteCall<BigInteger> intMap(String param0) {
        final Function function = new Function(FUNC_INTMAP, 
                Arrays.<Type>asList(new Utf8String(param0)),
                Arrays.<TypeReference<?>>asList(new TypeReference<Int256>() {}));
        return executeRemoteCallSingleValueReturn(function, BigInteger.class);
    }

    public RemoteCall<Tuple3<BigInteger, String, Boolean>> peopleMap(BigInteger param0) {
        final Function function = new Function(FUNC_PEOPLEMAP, 
                Arrays.<Type>asList(new Int256(param0)),
                Arrays.<TypeReference<?>>asList(new TypeReference<Uint256>() {}, new TypeReference<Utf8String>() {}, new TypeReference<Bool>() {}));
        return new RemoteCall<Tuple3<BigInteger, String, Boolean>>(
                new Callable<Tuple3<BigInteger, String, Boolean>>() {
                    @Override
                    public Tuple3<BigInteger, String, Boolean> call() throws Exception {
                        List<Type> results = executeCallMultipleValueReturn(function);
                        return new Tuple3<BigInteger, String, Boolean>(
                                (BigInteger) results.get(0).getValue(), 
                                (String) results.get(1).getValue(), 
                                (Boolean) results.get(2).getValue());
                    }
                });
    }

    public RemoteCall<BigInteger> sizeMap(BigInteger param0, BigInteger param1) {
        final Function function = new Function(FUNC_SIZEMAP, 
                Arrays.<Type>asList(new Int256(param0),
                new Int256(param1)),
                Arrays.<TypeReference<?>>asList(new TypeReference<Uint8>() {}));
        return executeRemoteCallSingleValueReturn(function, BigInteger.class);
    }

    public RemoteCall<String> stringMap(byte[] param0) {
        final Function function = new Function(FUNC_STRINGMAP, 
                Arrays.<Type>asList(new Bytes1(param0)),
                Arrays.<TypeReference<?>>asList(new TypeReference<Utf8String>() {}));
        return executeRemoteCallSingleValueReturn(function, String.class);
    }

    public RemoteCall<BigInteger> uintMap(byte[] param0) {
        final Function function = new Function(FUNC_UINTMAP, 
                Arrays.<Type>asList(new DynamicBytes(param0)),
                Arrays.<TypeReference<?>>asList(new TypeReference<Uint256>() {}));
        return executeRemoteCallSingleValueReturn(function, BigInteger.class);
    }

    public static RemoteCall<MappingDataTypeContract> deploy(Web3j web3j, Credentials credentials, GasProvider contractGasProvider, Long chainId) {
        return deployRemoteCall(MappingDataTypeContract.class, web3j, credentials, contractGasProvider, BINARY,  "", chainId);
    }

    public static RemoteCall<MappingDataTypeContract> deploy(Web3j web3j, TransactionManager transactionManager, GasProvider contractGasProvider, Long chainId) {
        return deployRemoteCall(MappingDataTypeContract.class, web3j, transactionManager, contractGasProvider, BINARY,  "", chainId);
    }

    public static MappingDataTypeContract load(String contractAddress, Web3j web3j, Credentials credentials, GasProvider contractGasProvider, Long chainId) {
        return new MappingDataTypeContract(contractAddress, web3j, credentials, contractGasProvider, chainId);
    }

    public static MappingDataTypeContract load(String contractAddress, Web3j web3j, TransactionManager transactionManager, GasProvider contractGasProvider, Long chainId) {
        return new MappingDataTypeContract(contractAddress, web3j, transactionManager, contractGasProvider, chainId);
    }
}
