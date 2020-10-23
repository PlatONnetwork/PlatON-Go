package network.platon.contracts.evm;

import java.math.BigInteger;
import java.util.Arrays;
import java.util.Collections;
import java.util.List;
import java.util.concurrent.Callable;
import org.web3j.abi.TypeReference;
import org.web3j.abi.datatypes.Address;
import org.web3j.abi.datatypes.DynamicBytes;
import org.web3j.abi.datatypes.Function;
import org.web3j.abi.datatypes.Type;
import org.web3j.abi.datatypes.generated.Bytes32;
import org.web3j.abi.datatypes.generated.StaticArray2;
import org.web3j.abi.datatypes.generated.Uint256;
import org.web3j.crypto.Credentials;
import org.web3j.protocol.Web3j;
import org.web3j.protocol.core.RemoteCall;
import org.web3j.protocol.core.methods.response.TransactionReceipt;
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
public class Precompiled extends Contract {
    private static final String BINARY = "608060405234801561001057600080fd5b50610f0e806100206000396000f3fe608060405234801561001057600080fd5b50600436106100cf5760003560e01c8063af13657c1161008c578063c3e6b01811610066578063c3e6b018146104d4578063caa260321461052a578063dd0678f11461065e578063ec8b466a1461072d576100cf565b8063af13657c1461035f578063b2acd5091461037d578063be3105401461044c576100cf565b806301f56b78146100d457806341be3d52146101635780636f29e2d7146101a95780637e59b08a146101c7578063897bf0401461024a5780638af606d414610290575b600080fd5b610121600480360360808110156100ea57600080fd5b8101908080359060200190929190803560ff16906020019092919080359060200190929190803590602001909291905050506107ab565b604051808273ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200191505060405180910390f35b61016b61081e565b6040518082600260200280838360005b8381101561019657808201518184015260208101905061017b565b5050505090500191505060405180910390f35b6101b1610868565b6040518082815260200191505060405180910390f35b6101cf610872565b6040518080602001828103825283818151815260200191508051906020019080838360005b8381101561020f5780820151818401526020810190506101f4565b50505050905090810190601f16801561023c5780820380516001836020036101000a031916815260200191505b509250505060405180910390f35b610252610914565b6040518082600260200280838360005b8381101561027d578082015181840152602081019050610262565b5050505090500191505060405180910390f35b610349600480360360208110156102a657600080fd5b81019080803590602001906401000000008111156102c357600080fd5b8201836020820111156102d557600080fd5b803590602001918460018302840111640100000000831117156102f757600080fd5b91908080601f016020809104026020016040519081016040528093929190818152602001838380828437600081840152601f19601f82011690508083019250505050505050919291929050505061095f565b6040518082815260200191505060405180910390f35b6103676109f9565b6040518082815260200191505060405180910390f35b6104366004803603602081101561039357600080fd5b81019080803590602001906401000000008111156103b057600080fd5b8201836020820111156103c257600080fd5b803590602001918460018302840111640100000000831117156103e457600080fd5b91908080601f016020809104026020016040519081016040528093929190818152602001838380828437600081840152601f19601f820116905080830192505050505050509192919290505050610a03565b6040518082815260200191505060405180910390f35b6104966004803603608081101561046257600080fd5b8101908080359060200190929190803590602001909291908035906020019092919080359060200190929190505050610a59565b6040518082600260200280838360005b838110156104c15780820151818401526020810190506104a6565b5050505090500191505060405180910390f35b610514600480360360608110156104ea57600080fd5b81019080803590602001909291908035906020019092919080359060200190929190505050610b01565b6040518082815260200191505060405180910390f35b6105e36004803603602081101561054057600080fd5b810190808035906020019064010000000081111561055d57600080fd5b82018360208201111561056f57600080fd5b8035906020019184600183028401116401000000008311171561059157600080fd5b91908080601f016020809104026020016040519081016040528093929190818152602001838380828437600081840152601f19601f820116905080830192505050505050509192919290505050610b5d565b6040518080602001828103825283818151815260200191508051906020019080838360005b83811015610623578082015181840152602081019050610608565b50505050905090810190601f1680156106505780820380516001836020036101000a031916815260200191505b509250505060405180910390f35b6107176004803603602081101561067457600080fd5b810190808035906020019064010000000081111561069157600080fd5b8201836020820111156106a357600080fd5b803590602001918460018302840111640100000000831117156106c557600080fd5b91908080601f016020809104026020016040519081016040528093929190818152602001838380828437600081840152601f19601f820116905080830192505050505050509192919290505050610bd0565b6040518082815260200191505060405180910390f35b61076d6004803603606081101561074357600080fd5b81019080803590602001909291908035906020019092919080359060200190929190505050610c77565b6040518082600260200280838360005b8381101561079857808201518184015260208101905061077d565b5050505090500191505060405180910390f35b600060018585858560405160008152602001604052604051808581526020018460ff1660ff1681526020018381526020018281526020019450505050506020604051602081039080840390855afa15801561080a573d6000803e3d6000fd5b505050602060405103519050949350505050565b610826610d07565b6002808060200260405190810160405280929190826002801561085e576020028201915b81548152602001906001019080831161084a575b5050505050905090565b6000600154905090565b606060008054600181600116156101000203166002900480601f01602080910402602001604051908101604052809291908181526020018280546001816001161561010002031660029004801561090a5780601f106108df5761010080835404028352916020019161090a565b820191906000526020600020905b8154815290600101906020018083116108ed57829003601f168201915b5050505050905090565b61091c610d29565b6004600280602002604051908101604052809291908260028015610955576020028201915b815481526020019060010190808311610941575b5050505050905090565b60006003826040518082805190602001908083835b602083106109975780518252602082019150602081019050602083039250610974565b6001836020036101000a038019825116818451168082178552505050505050905001915050602060405180830381855afa1580156109d9573d6000803e3d6000fd5b5050506040515160601b6bffffffffffffffffffffffff19169050919050565b6000600654905090565b60008082519050600060c08281610a1657fe5b0614610a2157600080fd5b6040516020818360208701600060085af18060008114610a445782519450610a49565b600080fd5b5050508160068190555050919050565b610a61610d07565b610a69610d4b565b8581600060048110610a7757fe5b6020020181815250508481600160048110610a8e57fe5b6020020181815250508381600260048110610aa557fe5b6020020181815250508281600360048110610abc57fe5b602002018181525050604082608083600060065af18060008114610adf57610ae4565b600080fd5b5050816002906002610af7929190610d6d565b5050949350505050565b600060405160208152602080820152602060408201528460608201528360808201528260a082015260208160c083600060055af18060008114610b475782519350610b4c565b600080fd5b505050806001819055509392505050565b60608082516040519080825280601f01601f191660200182016040528015610b945781602001600182028038833980820191505090505b509050825180602083018260208701600060045af1610baf57fe5b508060009080519060200190610bc6929190610dad565b5080915050919050565b60006002826040518082805190602001908083835b60208310610c085780518252602082019150602081019050602083039250610be5565b6001836020036101000a038019825116818451168082178552505050505050905001915050602060405180830381855afa158015610c4a573d6000803e3d6000fd5b5050506040513d6020811015610c5f57600080fd5b81019080805190602001909291905050509050919050565b610c7f610d29565b610c87610e2d565b8481600060038110610c9557fe5b6020020181815250508381600160038110610cac57fe5b6020020181815250508281600260038110610cc357fe5b602002018181525050604082606083600060075af18060008114610ce657610ceb565b600080fd5b5050816004906002610cfe929190610e4f565b50509392505050565b6040518060400160405280600290602082028038833980820191505090505090565b6040518060400160405280600290602082028038833980820191505090505090565b6040518060800160405280600490602082028038833980820191505090505090565b8260028101928215610d9c579160200282015b82811115610d9b578251825591602001919060010190610d80565b5b509050610da99190610e8f565b5090565b828054600181600116156101000203166002900490600052602060002090601f016020900481019282601f10610dee57805160ff1916838001178555610e1c565b82800160010185558215610e1c579182015b82811115610e1b578251825591602001919060010190610e00565b5b509050610e299190610e8f565b5090565b6040518060600160405280600390602082028038833980820191505090505090565b8260028101928215610e7e579160200282015b82811115610e7d578251825591602001919060010190610e62565b5b509050610e8b9190610eb4565b5090565b610eb191905b80821115610ead576000816000905550600101610e95565b5090565b90565b610ed691905b80821115610ed2576000816000905550600101610eba565b5090565b9056fea265627a7a72315820ba57dcd95e7ad1a566c40bb468db77342519377b0e7935be0b56739534205e4b64736f6c634300050d0032";

    public static final String FUNC_CALLBIGMODEXP = "callBigModExp";

    public static final String FUNC_CALLBN256ADD = "callBn256Add";

    public static final String FUNC_CALLBN256PAIRING = "callBn256Pairing";

    public static final String FUNC_CALLBN256SCALARMUL = "callBn256ScalarMul";

    public static final String FUNC_CALLDATACOPY = "callDatacopy";

    public static final String FUNC_CALLECRECOVER = "callEcrecover";

    public static final String FUNC_CALLRIPEMD160 = "callRipemd160";

    public static final String FUNC_CALLSHA256 = "callSha256";

    public static final String FUNC_GETCALLBIGMODEXPVALUE = "getCallBigModExpValue";

    public static final String FUNC_GETCALLBN256ADDVALUES = "getCallBn256AddValues";

    public static final String FUNC_GETCALLBN256PAIRINGVALUE = "getCallBn256PairingValue";

    public static final String FUNC_GETCALLBN256SCALARMULVALUES = "getCallBn256ScalarMulValues";

    public static final String FUNC_GETCALLDATACOPYVALUE = "getCallDatacopyValue";

    protected Precompiled(String contractAddress, Web3j web3j, Credentials credentials, GasProvider contractGasProvider, Long chainId) {
        super(BINARY, contractAddress, web3j, credentials, contractGasProvider, chainId);
    }

    protected Precompiled(String contractAddress, Web3j web3j, TransactionManager transactionManager, GasProvider contractGasProvider, Long chainId) {
        super(BINARY, contractAddress, web3j, transactionManager, contractGasProvider, chainId);
    }

    public RemoteCall<TransactionReceipt> callBigModExp(byte[] base, byte[] exponent, byte[] modulus) {
        final Function function = new Function(
                FUNC_CALLBIGMODEXP, 
                Arrays.<Type>asList(new Bytes32(base),
                new Bytes32(exponent),
                new Bytes32(modulus)),
                Collections.<TypeReference<?>>emptyList());
        return executeRemoteCallTransaction(function);
    }

    public RemoteCall<TransactionReceipt> callBn256Add(BigInteger ax, BigInteger ay, BigInteger bx, BigInteger by) {
        final Function function = new Function(
                FUNC_CALLBN256ADD, 
                Arrays.<Type>asList(new Uint256(ax),
                new Uint256(ay),
                new Uint256(bx),
                new Uint256(by)),
                Collections.<TypeReference<?>>emptyList());
        return executeRemoteCallTransaction(function);
    }

    public RemoteCall<TransactionReceipt> callBn256Pairing(byte[] input) {
        final Function function = new Function(
                FUNC_CALLBN256PAIRING, 
                Arrays.<Type>asList(new DynamicBytes(input)),
                Collections.<TypeReference<?>>emptyList());
        return executeRemoteCallTransaction(function);
    }

    public RemoteCall<TransactionReceipt> callBn256ScalarMul(byte[] x, byte[] y, byte[] scalar) {
        final Function function = new Function(
                FUNC_CALLBN256SCALARMUL, 
                Arrays.<Type>asList(new Bytes32(x),
                new Bytes32(y),
                new Bytes32(scalar)),
                Collections.<TypeReference<?>>emptyList());
        return executeRemoteCallTransaction(function);
    }

    public RemoteCall<TransactionReceipt> callDatacopy(byte[] data) {
        final Function function = new Function(
                FUNC_CALLDATACOPY, 
                Arrays.<Type>asList(new DynamicBytes(data)),
                Collections.<TypeReference<?>>emptyList());
        return executeRemoteCallTransaction(function);
    }

    public RemoteCall<String> callEcrecover(byte[] hash, BigInteger v, byte[] r, byte[] s) {
        final Function function = new Function(FUNC_CALLECRECOVER, 
                Arrays.<Type>asList(new Bytes32(hash),
                new org.web3j.abi.datatypes.generated.Uint8(v), 
                new Bytes32(r),
                new Bytes32(s)),
                Arrays.<TypeReference<?>>asList(new TypeReference<Address>() {}));
        return executeRemoteCallSingleValueReturn(function, String.class);
    }

    public RemoteCall<byte[]> callRipemd160(byte[] data) {
        final Function function = new Function(FUNC_CALLRIPEMD160, 
                Arrays.<Type>asList(new DynamicBytes(data)),
                Arrays.<TypeReference<?>>asList(new TypeReference<Bytes32>() {}));
        return executeRemoteCallSingleValueReturn(function, byte[].class);
    }

    public RemoteCall<byte[]> callSha256(byte[] data) {
        final Function function = new Function(FUNC_CALLSHA256, 
                Arrays.<Type>asList(new DynamicBytes(data)),
                Arrays.<TypeReference<?>>asList(new TypeReference<Bytes32>() {}));
        return executeRemoteCallSingleValueReturn(function, byte[].class);
    }

    public RemoteCall<byte[]> getCallBigModExpValue() {
        final Function function = new Function(FUNC_GETCALLBIGMODEXPVALUE, 
                Arrays.<Type>asList(), 
                Arrays.<TypeReference<?>>asList(new TypeReference<Bytes32>() {}));
        return executeRemoteCallSingleValueReturn(function, byte[].class);
    }

    public RemoteCall<List> getCallBn256AddValues() {
        final Function function = new Function(FUNC_GETCALLBN256ADDVALUES, 
                Arrays.<Type>asList(), 
                Arrays.<TypeReference<?>>asList(new TypeReference<StaticArray2<Uint256>>() {}));
        return new RemoteCall<List>(
                new Callable<List>() {
                    @Override
                    @SuppressWarnings("unchecked")
                    public List call() throws Exception {
                        List<Type> result = (List<Type>) executeCallSingleValueReturn(function, List.class);
                        return convertToNative(result);
                    }
                });
    }

    public RemoteCall<byte[]> getCallBn256PairingValue() {
        final Function function = new Function(FUNC_GETCALLBN256PAIRINGVALUE, 
                Arrays.<Type>asList(), 
                Arrays.<TypeReference<?>>asList(new TypeReference<Bytes32>() {}));
        return executeRemoteCallSingleValueReturn(function, byte[].class);
    }

    public RemoteCall<List> getCallBn256ScalarMulValues() {
        final Function function = new Function(FUNC_GETCALLBN256SCALARMULVALUES, 
                Arrays.<Type>asList(), 
                Arrays.<TypeReference<?>>asList(new TypeReference<StaticArray2<Bytes32>>() {}));
        return new RemoteCall<List>(
                new Callable<List>() {
                    @Override
                    @SuppressWarnings("unchecked")
                    public List call() throws Exception {
                        List<Type> result = (List<Type>) executeCallSingleValueReturn(function, List.class);
                        return convertToNative(result);
                    }
                });
    }

    public RemoteCall<byte[]> getCallDatacopyValue() {
        final Function function = new Function(FUNC_GETCALLDATACOPYVALUE, 
                Arrays.<Type>asList(), 
                Arrays.<TypeReference<?>>asList(new TypeReference<DynamicBytes>() {}));
        return executeRemoteCallSingleValueReturn(function, byte[].class);
    }

    public static RemoteCall<Precompiled> deploy(Web3j web3j, Credentials credentials, GasProvider contractGasProvider, Long chainId) {
        return deployRemoteCall(Precompiled.class, web3j, credentials, contractGasProvider, BINARY,  "", chainId);
    }

    public static RemoteCall<Precompiled> deploy(Web3j web3j, TransactionManager transactionManager, GasProvider contractGasProvider, Long chainId) {
        return deployRemoteCall(Precompiled.class, web3j, transactionManager, contractGasProvider, BINARY,  "", chainId);
    }

    public static Precompiled load(String contractAddress, Web3j web3j, Credentials credentials, GasProvider contractGasProvider, Long chainId) {
        return new Precompiled(contractAddress, web3j, credentials, contractGasProvider, chainId);
    }

    public static Precompiled load(String contractAddress, Web3j web3j, TransactionManager transactionManager, GasProvider contractGasProvider, Long chainId) {
        return new Precompiled(contractAddress, web3j, transactionManager, contractGasProvider, chainId);
    }
}
