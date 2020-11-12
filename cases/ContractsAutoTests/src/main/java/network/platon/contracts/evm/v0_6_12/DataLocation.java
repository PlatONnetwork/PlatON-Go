package network.platon.contracts.evm.v0_6_12;

import com.alaya.abi.solidity.TypeReference;
import com.alaya.abi.solidity.datatypes.DynamicBytes;
import com.alaya.abi.solidity.datatypes.Function;
import com.alaya.abi.solidity.datatypes.Type;
import com.alaya.abi.solidity.datatypes.Utf8String;
import com.alaya.abi.solidity.datatypes.generated.Uint256;
import com.alaya.crypto.Credentials;
import com.alaya.protocol.Web3j;
import com.alaya.protocol.core.RemoteCall;
import com.alaya.protocol.core.methods.response.TransactionReceipt;
import com.alaya.tuples.generated.Tuple2;
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
public class DataLocation extends Contract {
    private static final String BINARY = "608060405234801561001057600080fd5b5061076b806100206000396000f3fe608060405234801561001057600080fd5b506004361061004c5760003560e01c80630bcd3b3314610051578063246982c4146100d45780633ca8b1a714610182578063a1715deb14610274575b600080fd5b610059610359565b6040518080602001828103825283818151815260200191508051906020019080838360005b8381101561009957808201518184015260208101905061007e565b50505050905090810190601f1680156100c65780820380516001836020036101000a031916815260200191505b509250505060405180910390f35b610100600480360360208110156100ea57600080fd5b81019080803590602001909291905050506103fb565b6040518080602001838152602001828103825284818151815260200191508051906020019080838360005b8381101561014657808201518184015260208101905061012b565b50505050905090810190601f1680156101735780820380516001836020036101000a031916815260200191505b50935050505060405180910390f35b6101f96004803603602081101561019857600080fd5b81019080803590602001906401000000008111156101b557600080fd5b8201836020820111156101c757600080fd5b803590602001918460018302840111640100000000831117156101e957600080fd5b90919293919293905050506104cf565b6040518080602001828103825283818151815260200191508051906020019080838360005b8381101561023957808201518184015260208101905061021e565b50505050905090810190601f1680156102665780820380516001836020036101000a031916815260200191505b509250505060405180910390f35b6103416004803603606081101561028a57600080fd5b8101908080359060200190929190803590602001906401000000008111156102b157600080fd5b8201836020820111156102c357600080fd5b803590602001918460018302840111640100000000831117156102e557600080fd5b91908080601f016020809104026020016040519081016040528093929190818152602001838380828437600081840152601f19601f82011690508083019250505050505050919291929080359060200190929190505050610586565b60405180821515815260200191505060405180910390f35b606060018054600181600116156101000203166002900480601f0160208091040260200160405190810160405280929190818152602001828054600181600116156101000203166002900480156103f15780601f106103c6576101008083540402835291602001916103f1565b820191906000526020600020905b8154815290600101906020018083116103d457829003601f168201915b5050505050905090565b6060600080600084815260200190815260200160002060000160008085815260200190815260200160002060010154818054600181600116156101000203166002900480601f0160208091040260200160405190810160405280929190818152602001828054600181600116156101000203166002900480156104bf5780601f10610494576101008083540402835291602001916104bf565b820191906000526020600020905b8154815290600101906020018083116104a257829003601f168201915b5050505050915091509150915091565b60608282600191906104e29291906105fe565b5060018054600181600116156101000203166002900480601f0160208091040260200160405190810160405280929190818152602001828054600181600116156101000203166002900480156105795780601f1061054e57610100808354040283529160200191610579565b820191906000526020600020905b81548152906001019060200180831161055c57829003601f168201915b5050505050905092915050565b600061059061067e565b60405180604001604052808581526020018481525090506105b181866105bd565b60019150509392505050565b8160008083815260200190815260200160002060008201518160000190805190602001906105ec929190610698565b50602082015181600101559050505050565b828054600181600116156101000203166002900490600052602060002090601f016020900481019282601f1061063f57803560ff191683800117855561066d565b8280016001018555821561066d579182015b8281111561066c578235825591602001919060010190610651565b5b50905061067a9190610718565b5090565b604051806040016040528060608152602001600081525090565b828054600181600116156101000203166002900490600052602060002090601f016020900481019282601f106106d957805160ff1916838001178555610707565b82800160010185558215610707579182015b828111156107065782518255916020019190600101906106eb565b5b5090506107149190610718565b5090565b5b80821115610731576000816000905550600101610719565b509056fea2646970667358221220a8d5c9c1b4250e0c6d9809b204e3d0bb301804fb60839f4d63a63564570305b364736f6c634300060c0033";

    public static final String FUNC_GETBYTES = "getBytes";

    public static final String FUNC_GETPERSON = "getPerson";

    public static final String FUNC_SAVEPERSON = "savePerson";

    public static final String FUNC_TESTBYTES = "testBytes";

    protected DataLocation(String contractAddress, Web3j web3j, Credentials credentials, GasProvider contractGasProvider, Long chainId) {
        super(BINARY, contractAddress, web3j, credentials, contractGasProvider, chainId);
    }

    protected DataLocation(String contractAddress, Web3j web3j, TransactionManager transactionManager, GasProvider contractGasProvider, Long chainId) {
        super(BINARY, contractAddress, web3j, transactionManager, contractGasProvider, chainId);
    }

    public RemoteCall<byte[]> getBytes() {
        final Function function = new Function(FUNC_GETBYTES, 
                Arrays.<Type>asList(), 
                Arrays.<TypeReference<?>>asList(new TypeReference<DynamicBytes>() {}));
        return executeRemoteCallSingleValueReturn(function, byte[].class);
    }

    public RemoteCall<Tuple2<String, BigInteger>> getPerson(BigInteger _id) {
        final Function function = new Function(FUNC_GETPERSON, 
                Arrays.<Type>asList(new Uint256(_id)),
                Arrays.<TypeReference<?>>asList(new TypeReference<Utf8String>() {}, new TypeReference<Uint256>() {}));
        return new RemoteCall<Tuple2<String, BigInteger>>(
                new Callable<Tuple2<String, BigInteger>>() {
                    @Override
                    public Tuple2<String, BigInteger> call() throws Exception {
                        List<Type> results = executeCallMultipleValueReturn(function);
                        return new Tuple2<String, BigInteger>(
                                (String) results.get(0).getValue(), 
                                (BigInteger) results.get(1).getValue());
                    }
                });
    }

    public RemoteCall<TransactionReceipt> savePerson(BigInteger _id, String _name, BigInteger _age) {
        final Function function = new Function(
                FUNC_SAVEPERSON, 
                Arrays.<Type>asList(new Uint256(_id),
                new Utf8String(_name),
                new Uint256(_age)),
                Collections.<TypeReference<?>>emptyList());
        return executeRemoteCallTransaction(function);
    }

    public RemoteCall<TransactionReceipt> testBytes(byte[] _data) {
        final Function function = new Function(
                FUNC_TESTBYTES, 
                Arrays.<Type>asList(new DynamicBytes(_data)),
                Collections.<TypeReference<?>>emptyList());
        return executeRemoteCallTransaction(function);
    }

    public static RemoteCall<DataLocation> deploy(Web3j web3j, Credentials credentials, GasProvider contractGasProvider, Long chainId) {
        return deployRemoteCall(DataLocation.class, web3j, credentials, contractGasProvider, BINARY,  "", chainId);
    }

    public static RemoteCall<DataLocation> deploy(Web3j web3j, TransactionManager transactionManager, GasProvider contractGasProvider, Long chainId) {
        return deployRemoteCall(DataLocation.class, web3j, transactionManager, contractGasProvider, BINARY,  "", chainId);
    }

    public static DataLocation load(String contractAddress, Web3j web3j, Credentials credentials, GasProvider contractGasProvider, Long chainId) {
        return new DataLocation(contractAddress, web3j, credentials, contractGasProvider, chainId);
    }

    public static DataLocation load(String contractAddress, Web3j web3j, TransactionManager transactionManager, GasProvider contractGasProvider, Long chainId) {
        return new DataLocation(contractAddress, web3j, transactionManager, contractGasProvider, chainId);
    }
}
