package network.platon.contracts.evm.v0_6_12;

import com.alaya.abi.solidity.TypeReference;
import com.alaya.abi.solidity.datatypes.Function;
import com.alaya.abi.solidity.datatypes.Type;
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
 * <p>Generated with web3j version 0.13.2.0.
 */
public class MappingDataTypeContract extends Contract {
    private static final String BINARY = "608060405260405180606001604052806040518060400160405280600481526020017f4c7563790000000000000000000000000000000000000000000000000000000081525081526020016040518060400160405280600481526020017f456c6c610000000000000000000000000000000000000000000000000000000081525081526020016040518060400160405280600481526020017f4c696c7900000000000000000000000000000000000000000000000000000000815250815250600d906003620000d0929190620000e5565b50348015620000de57600080fd5b5062000266565b82805482825590600052602060002090810192821562000139579160200282015b8281111562000138578251829080519060200190620001279291906200014c565b509160200191906001019062000106565b5b509050620001489190620001d3565b5090565b828054600181600116156101000203166002900490600052602060002090601f016020900481019282601f106200018f57805160ff1916838001178555620001c0565b82800160010185558215620001c0579182015b82811115620001bf578251825591602001919060010190620001a2565b5b509050620001cf9190620001fb565b5090565b5b80821115620001f75760008181620001ed91906200021a565b50600101620001d4565b5090565b5b8082111562000216576000816000905550600101620001fc565b5090565b50805460018160011615610100020316600290046000825580601f1062000242575062000263565b601f016020900490600052602060002090810190620002629190620001fb565b5b50565b610c3080620002766000396000f3fe608060405234801561001057600080fd5b50600436106100b35760003560e01c806387c9644b1161007157806387c9644b14610455578063aad3951214610499578063bbf722a214610562578063bebf70f3146105ba578063f480ae4914610611578063fdf8b008146106e0576100b3565b806252f6ee146100b8578063566ab1fd146101875780636257a397146101d45780636b8ff5741461029157806376fc59c31461033857806382bc11491461039e575b600080fd5b610171600480360360208110156100ce57600080fd5b81019080803590602001906401000000008111156100eb57600080fd5b8201836020820111156100fd57600080fd5b8035906020019184600183028401116401000000008311171561011f57600080fd5b91908080601f016020809104026020016040519081016040528093929190818152602001838380828437600081840152601f19601f8201169050808301925050505050505091929192905050506106ea565b6040518082815260200191505060405180910390f35b6101b36004803603602081101561019d57600080fd5b8101908080359060200190929190505050610718565b604051808260028111156101c357fe5b815260200191505060405180910390f35b610216600480360360208110156101ea57600080fd5b81019080803573ffffffffffffffffffffffffffffffffffffffff169060200190929190505050610738565b6040518080602001828103825283818151815260200191508051906020019080838360005b8381101561025657808201518184015260208101905061023b565b50505050905090810190601f1680156102835780820380516001836020036101000a031916815260200191505b509250505060405180910390f35b6102bd600480360360208110156102a757600080fd5b81019080803590602001909291905050506107e8565b6040518080602001828103825283818151815260200191508051906020019080838360005b838110156102fd5780820151818401526020810190506102e2565b50505050905090810190601f16801561032a5780820380516001836020036101000a031916815260200191505b509250505060405180910390f35b6103666004803603602081101561034e57600080fd5b8101908080351515906020019092919050505061089d565b60405180827effffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff1916815260200191505060405180910390f35b6103ca600480360360208110156103b457600080fd5b81019080803590602001909291905050506108bd565b60405180848152602001806020018315158152602001828103825284818151815260200191508051906020019080838360005b838110156104185780820151818401526020810190506103fd565b50505050905090810190601f1680156104455780820380516001836020036101000a031916815260200191505b5094505050505060405180910390f35b6104816004803603602081101561046b57600080fd5b810190808035906020019092919050505061098c565b60405180821515815260200191505060405180910390f35b6104e7600480360360208110156104af57600080fd5b8101908080357effffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff191690602001909291905050506109ac565b6040518080602001828103825283818151815260200191508051906020019080838360005b8381101561052757808201518184015260208101905061050c565b50505050905090810190601f1680156105545780820380516001836020036101000a031916815260200191505b509250505060405180910390f35b61058e6004803603602081101561057857600080fd5b8101908080359060200190929190505050610a5c565b604051808273ffffffffffffffffffffffffffffffffffffffff16815260200191505060405180910390f35b6105f0600480360360408110156105d057600080fd5b810190808035906020019092919080359060200190929190505050610a8f565b6040518082600281111561060057fe5b815260200191505060405180910390f35b6106ca6004803603602081101561062757600080fd5b810190808035906020019064010000000081111561064457600080fd5b82018360208201111561065657600080fd5b8035906020019184600183028401116401000000008311171561067857600080fd5b91908080601f016020809104026020016040519081016040528093929190818152602001838380828437600081840152601f19601f820116905080830192505050505050509192919290505050610abe565b6040518082815260200191505060405180910390f35b6106e8610aec565b005b6006818051602081018201805184825260208301602085012081835280955050505050506000915090505481565b600b6020528060005260406000206000915054906101000a900460ff1681565b60086020528060005260406000206000915090508054600181600116156101000203166002900480601f0160208091040260200160405190810160405280929190818152602001828054600181600116156101000203166002900480156107e05780601f106107b5576101008083540402835291602001916107e0565b820191906000526020600020905b8154815290600101906020018083116107c357829003601f168201915b505050505081565b6060600c60008381526020019081526020016000208054600181600116156101000203166002900480601f0160208091040260200160405190810160405280929190818152602001828054600181600116156101000203166002900480156108915780601f1061086657610100808354040283529160200191610891565b820191906000526020600020905b81548152906001019060200180831161087457829003601f168201915b50505050509050919050565b60046020528060005260406000206000915054906101000a900460f81b81565b600a602052806000526040600020600091509050806000015490806001018054600181600116156101000203166002900480601f01602080910402602001604051908101604052809291908181526020018280546001816001161561010002031660029004801561096f5780601f106109445761010080835404028352916020019161096f565b820191906000526020600020905b81548152906001019060200180831161095257829003601f168201915b5050505050908060020160009054906101000a900460ff16905083565b60036020528060005260406000206000915054906101000a900460ff1681565b60056020528060005260406000206000915090508054600181600116156101000203166002900480601f016020809104026020016040519081016040528092919081815260200182805460018160011615610100020316600290048015610a545780601f10610a2957610100808354040283529160200191610a54565b820191906000526020600020905b815481529060010190602001808311610a3757829003601f168201915b505050505081565b60026020528060005260406000206000915054906101000a900473ffffffffffffffffffffffffffffffffffffffff1681565b60096020528160005260406000206020528060005260406000206000915091509054906101000a900460ff1681565b6007818051602081018201805184825260208301602085012081835280955050505050506000915090505481565b60005b600d80549050811015610b5357600d8181548110610b0957fe5b90600052602060002001600c60008381526020019081526020016000209080546001816001161561010002031660029004610b45929190610b56565b508080600101915050610aef565b50565b828054600181600116156101000203166002900490600052602060002090601f016020900481019282601f10610b8f5780548555610bcc565b82800160010185558215610bcc57600052602060002091601f016020900482015b82811115610bcb578254825591600101919060010190610bb0565b5b509050610bd99190610bdd565b5090565b5b80821115610bf6576000816000905550600101610bde565b509056fea264697066735822122027c13c41c2679383c2610d8bfb9c6c37733734d5a3f61751e0cecb7fb98d1c1164736f6c634300060c0033";

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

    public RemoteCall<TransactionReceipt> SizeEnumMap(BigInteger param0) {
        final Function function = new Function(
                FUNC_SIZEENUMMAP, 
                Arrays.<Type>asList(new com.alaya.abi.solidity.datatypes.generated.Int256(param0)), 
                Collections.<TypeReference<?>>emptyList());
        return executeRemoteCallTransaction(function);
    }

    public RemoteCall<TransactionReceipt> addName() {
        final Function function = new Function(
                FUNC_ADDNAME, 
                Arrays.<Type>asList(), 
                Collections.<TypeReference<?>>emptyList());
        return executeRemoteCallTransaction(function);
    }

    public RemoteCall<TransactionReceipt> addressMap(BigInteger param0) {
        final Function function = new Function(
                FUNC_ADDRESSMAP, 
                Arrays.<Type>asList(new com.alaya.abi.solidity.datatypes.generated.Uint256(param0)), 
                Collections.<TypeReference<?>>emptyList());
        return executeRemoteCallTransaction(function);
    }

    public RemoteCall<TransactionReceipt> boolMap(BigInteger param0) {
        final Function function = new Function(
                FUNC_BOOLMAP, 
                Arrays.<Type>asList(new com.alaya.abi.solidity.datatypes.generated.Int256(param0)), 
                Collections.<TypeReference<?>>emptyList());
        return executeRemoteCallTransaction(function);
    }

    public RemoteCall<TransactionReceipt> byteMap(Boolean param0) {
        final Function function = new Function(
                FUNC_BYTEMAP, 
                Arrays.<Type>asList(new com.alaya.abi.solidity.datatypes.Bool(param0)), 
                Collections.<TypeReference<?>>emptyList());
        return executeRemoteCallTransaction(function);
    }

    public RemoteCall<TransactionReceipt> bytesMap(String param0) {
        final Function function = new Function(
                FUNC_BYTESMAP, 
                Arrays.<Type>asList(new com.alaya.abi.solidity.datatypes.Address(param0)), 
                Collections.<TypeReference<?>>emptyList());
        return executeRemoteCallTransaction(function);
    }

    public RemoteCall<TransactionReceipt> getName(BigInteger index) {
        final Function function = new Function(
                FUNC_GETNAME, 
                Arrays.<Type>asList(new com.alaya.abi.solidity.datatypes.generated.Uint256(index)), 
                Collections.<TypeReference<?>>emptyList());
        return executeRemoteCallTransaction(function);
    }

    public RemoteCall<TransactionReceipt> intMap(String param0) {
        final Function function = new Function(
                FUNC_INTMAP, 
                Arrays.<Type>asList(new com.alaya.abi.solidity.datatypes.Utf8String(param0)), 
                Collections.<TypeReference<?>>emptyList());
        return executeRemoteCallTransaction(function);
    }

    public RemoteCall<TransactionReceipt> peopleMap(BigInteger param0) {
        final Function function = new Function(
                FUNC_PEOPLEMAP, 
                Arrays.<Type>asList(new com.alaya.abi.solidity.datatypes.generated.Int256(param0)), 
                Collections.<TypeReference<?>>emptyList());
        return executeRemoteCallTransaction(function);
    }

    public RemoteCall<TransactionReceipt> sizeMap(BigInteger param0, BigInteger param1) {
        final Function function = new Function(
                FUNC_SIZEMAP, 
                Arrays.<Type>asList(new com.alaya.abi.solidity.datatypes.generated.Int256(param0), 
                new com.alaya.abi.solidity.datatypes.generated.Int256(param1)), 
                Collections.<TypeReference<?>>emptyList());
        return executeRemoteCallTransaction(function);
    }

    public RemoteCall<TransactionReceipt> stringMap(byte[] param0) {
        final Function function = new Function(
                FUNC_STRINGMAP, 
                Arrays.<Type>asList(new com.alaya.abi.solidity.datatypes.generated.Bytes1(param0)), 
                Collections.<TypeReference<?>>emptyList());
        return executeRemoteCallTransaction(function);
    }

    public RemoteCall<TransactionReceipt> uintMap(byte[] param0) {
        final Function function = new Function(
                FUNC_UINTMAP, 
                Arrays.<Type>asList(new com.alaya.abi.solidity.datatypes.DynamicBytes(param0)), 
                Collections.<TypeReference<?>>emptyList());
        return executeRemoteCallTransaction(function);
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
