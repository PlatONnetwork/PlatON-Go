package network.platon.contracts.evm.v0_7_1;

import com.alaya.abi.solidity.TypeReference;
import com.alaya.abi.solidity.datatypes.Address;
import com.alaya.abi.solidity.datatypes.DynamicBytes;
import com.alaya.abi.solidity.datatypes.Function;
import com.alaya.abi.solidity.datatypes.Type;
import com.alaya.abi.solidity.datatypes.generated.Bytes32;
import com.alaya.abi.solidity.datatypes.generated.StaticArray2;
import com.alaya.abi.solidity.datatypes.generated.Uint256;
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
public class Precompiled extends Contract {
    private static final String BINARY = "608060405234801561001057600080fd5b50610efd806100206000396000f3fe608060405234801561001057600080fd5b50600436106100cf5760003560e01c8063af13657c1161008c578063c3e6b01811610066578063c3e6b018146104be578063caa2603214610514578063dd0678f114610648578063ec8b466a14610717576100cf565b8063af13657c14610349578063b2acd50914610367578063be31054014610436576100cf565b806301f56b78146100d457806341be3d521461014d5780636f29e2d7146101935780637e59b08a146101b1578063897bf040146102345780638af606d41461027a575b600080fd5b610121600480360360808110156100ea57600080fd5b8101908080359060200190929190803560ff1690602001909291908035906020019092919080359060200190929190505050610795565b604051808273ffffffffffffffffffffffffffffffffffffffff16815260200191505060405180910390f35b610155610805565b6040518082600260200280838360005b83811015610180578082015181840152602081019050610165565b5050505090500191505060405180910390f35b61019b61084f565b6040518082815260200191505060405180910390f35b6101b9610859565b6040518080602001828103825283818151815260200191508051906020019080838360005b838110156101f95780820151818401526020810190506101de565b50505050905090810190601f1680156102265780820380516001836020036101000a031916815260200191505b509250505060405180910390f35b61023c6108fb565b6040518082600260200280838360005b8381101561026757808201518184015260208101905061024c565b5050505090500191505060405180910390f35b6103336004803603602081101561029057600080fd5b81019080803590602001906401000000008111156102ad57600080fd5b8201836020820111156102bf57600080fd5b803590602001918460018302840111640100000000831117156102e157600080fd5b91908080601f016020809104026020016040519081016040528093929190818152602001838380828437600081840152601f19601f820116905080830192505050505050509192919290505050610946565b6040518082815260200191505060405180910390f35b6103516109e0565b6040518082815260200191505060405180910390f35b6104206004803603602081101561037d57600080fd5b810190808035906020019064010000000081111561039a57600080fd5b8201836020820111156103ac57600080fd5b803590602001918460018302840111640100000000831117156103ce57600080fd5b91908080601f016020809104026020016040519081016040528093929190818152602001838380828437600081840152601f19601f8201169050808301925050505050505091929192905050506109ea565b6040518082815260200191505060405180910390f35b6104806004803603608081101561044c57600080fd5b8101908080359060200190929190803590602001909291908035906020019092919080359060200190929190505050610a40565b6040518082600260200280838360005b838110156104ab578082015181840152602081019050610490565b5050505090500191505060405180910390f35b6104fe600480360360608110156104d457600080fd5b81019080803590602001909291908035906020019092919080359060200190929190505050610ae8565b6040518082815260200191505060405180910390f35b6105cd6004803603602081101561052a57600080fd5b810190808035906020019064010000000081111561054757600080fd5b82018360208201111561055957600080fd5b8035906020019184600183028401116401000000008311171561057b57600080fd5b91908080601f016020809104026020016040519081016040528093929190818152602001838380828437600081840152601f19601f820116905080830192505050505050509192919290505050610b44565b6040518080602001828103825283818151815260200191508051906020019080838360005b8381101561060d5780820151818401526020810190506105f2565b50505050905090810190601f16801561063a5780820380516001836020036101000a031916815260200191505b509250505060405180910390f35b6107016004803603602081101561065e57600080fd5b810190808035906020019064010000000081111561067b57600080fd5b82018360208201111561068d57600080fd5b803590602001918460018302840111640100000000831117156106af57600080fd5b91908080601f016020809104026020016040519081016040528093929190818152602001838380828437600081840152601f19601f820116905080830192505050505050509192919290505050610bce565b6040518082815260200191505060405180910390f35b6107576004803603606081101561072d57600080fd5b81019080803590602001909291908035906020019092919080359060200190929190505050610c75565b6040518082600260200280838360005b83811015610782578082015181840152602081019050610767565b5050505090500191505060405180910390f35b600060018585858560405160008152602001604052604051808581526020018460ff1681526020018381526020018281526020019450505050506020604051602081039080840390855afa1580156107f1573d6000803e3d6000fd5b505050602060405103519050949350505050565b61080d610d05565b60028080602002604051908101604052809291908260028015610845576020028201915b815481526020019060010190808311610831575b5050505050905090565b6000600154905090565b606060008054600181600116156101000203166002900480601f0160208091040260200160405190810160405280929190818152602001828054600181600116156101000203166002900480156108f15780601f106108c6576101008083540402835291602001916108f1565b820191906000526020600020905b8154815290600101906020018083116108d457829003601f168201915b5050505050905090565b610903610d27565b600460028060200260405190810160405280929190826002801561093c576020028201915b815481526020019060010190808311610928575b5050505050905090565b60006003826040518082805190602001908083835b6020831061097e578051825260208201915060208101905060208303925061095b565b6001836020036101000a038019825116818451168082178552505050505050905001915050602060405180830381855afa1580156109c0573d6000803e3d6000fd5b5050506040515160601b6bffffffffffffffffffffffff19169050919050565b6000600654905090565b60008082519050600060c082816109fd57fe5b0614610a0857600080fd5b6040516020818360208701600060085af18060008114610a2b5782519450610a30565b600080fd5b5050508160068190555050919050565b610a48610d05565b610a50610d49565b8581600060048110610a5e57fe5b6020020181815250508481600160048110610a7557fe5b6020020181815250508381600260048110610a8c57fe5b6020020181815250508281600360048110610aa357fe5b602002018181525050604082608083600060065af18060008114610ac657610acb565b600080fd5b5050816002906002610ade929190610d6b565b5050949350505050565b600060405160208152602080820152602060408201528460608201528360808201528260a082015260208160c083600060055af18060008114610b2e5782519350610b33565b600080fd5b505050806001819055509392505050565b606080825167ffffffffffffffff81118015610b5f57600080fd5b506040519080825280601f01601f191660200182016040528015610b925781602001600182028036833780820191505090505b509050825180602083018260208701600060045af1610bad57fe5b508060009080519060200190610bc4929190610dab565b5080915050919050565b60006002826040518082805190602001908083835b60208310610c065780518252602082019150602081019050602083039250610be3565b6001836020036101000a038019825116818451168082178552505050505050905001915050602060405180830381855afa158015610c48573d6000803e3d6000fd5b5050506040513d6020811015610c5d57600080fd5b81019080805190602001909291905050509050919050565b610c7d610d27565b610c85610e2b565b8481600060038110610c9357fe5b6020020181815250508381600160038110610caa57fe5b6020020181815250508281600260038110610cc157fe5b602002018181525050604082606083600060075af18060008114610ce457610ce9565b600080fd5b5050816004906002610cfc929190610e4d565b50509392505050565b6040518060400160405280600290602082028036833780820191505090505090565b6040518060400160405280600290602082028036833780820191505090505090565b6040518060800160405280600490602082028036833780820191505090505090565b8260028101928215610d9a579160200282015b82811115610d99578251825591602001919060010190610d7e565b5b509050610da79190610e8d565b5090565b828054600181600116156101000203166002900490600052602060002090601f016020900481019282601f10610dec57805160ff1916838001178555610e1a565b82800160010185558215610e1a579182015b82811115610e19578251825591602001919060010190610dfe565b5b509050610e279190610e8d565b5090565b6040518060600160405280600390602082028036833780820191505090505090565b8260028101928215610e7c579160200282015b82811115610e7b578251825591602001919060010190610e60565b5b509050610e899190610eaa565b5090565b5b80821115610ea6576000816000905550600101610e8e565b5090565b5b80821115610ec3576000816000905550600101610eab565b509056fea2646970667358221220ffbee8978a3d4b5782db57a60080a2265ac06561505c235f2e19ecd31cc4b93764736f6c63430007010033";

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
                new com.alaya.abi.solidity.datatypes.generated.Uint8(v), 
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
