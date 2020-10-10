package network.platon.contracts.evm.v0_5_17;

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
 * <p>Generated with web3j version 0.13.2.0.
 */
public class DataLocation extends Contract {
    private static final String BINARY = "608060405234801561001057600080fd5b50610774806100206000396000f3fe608060405234801561001057600080fd5b506004361061004c5760003560e01c80630bcd3b3314610051578063246982c4146100d45780633ca8b1a714610182578063a1715deb14610274575b600080fd5b61005961035b565b6040518080602001828103825283818151815260200191508051906020019080838360005b8381101561009957808201518184015260208101905061007e565b50505050905090810190601f1680156100c65780820380516001836020036101000a031916815260200191505b509250505060405180910390f35b610100600480360360208110156100ea57600080fd5b81019080803590602001909291905050506103fd565b6040518080602001838152602001828103825284818151815260200191508051906020019080838360005b8381101561014657808201518184015260208101905061012b565b50505050905090810190601f1680156101735780820380516001836020036101000a031916815260200191505b50935050505060405180910390f35b6101f96004803603602081101561019857600080fd5b81019080803590602001906401000000008111156101b557600080fd5b8201836020820111156101c757600080fd5b803590602001918460018302840111640100000000831117156101e957600080fd5b90919293919293905050506104d1565b6040518080602001828103825283818151815260200191508051906020019080838360005b8381101561023957808201518184015260208101905061021e565b50505050905090810190601f1680156102665780820380516001836020036101000a031916815260200191505b509250505060405180910390f35b6103416004803603606081101561028a57600080fd5b8101908080359060200190929190803590602001906401000000008111156102b157600080fd5b8201836020820111156102c357600080fd5b803590602001918460018302840111640100000000831117156102e557600080fd5b91908080601f016020809104026020016040519081016040528093929190818152602001838380828437600081840152601f19601f82011690508083019250505050505050919291929080359060200190929190505050610588565b604051808215151515815260200191505060405180910390f35b606060018054600181600116156101000203166002900480601f0160208091040260200160405190810160405280929190818152602001828054600181600116156101000203166002900480156103f35780601f106103c8576101008083540402835291602001916103f3565b820191906000526020600020905b8154815290600101906020018083116103d657829003601f168201915b5050505050905090565b6060600080600084815260200190815260200160002060000160008085815260200190815260200160002060010154818054600181600116156101000203166002900480601f0160208091040260200160405190810160405280929190818152602001828054600181600116156101000203166002900480156104c15780601f10610496576101008083540402835291602001916104c1565b820191906000526020600020905b8154815290600101906020018083116104a457829003601f168201915b5050505050915091509150915091565b60608282600191906104e4929190610600565b5060018054600181600116156101000203166002900480601f01602080910402602001604051908101604052809291908181526020018280546001816001161561010002031660029004801561057b5780601f106105505761010080835404028352916020019161057b565b820191906000526020600020905b81548152906001019060200180831161055e57829003601f168201915b5050505050905092915050565b6000610592610680565b60405180604001604052808581526020018481525090506105b381866105bf565b60019150509392505050565b8160008083815260200190815260200160002060008201518160000190805190602001906105ee92919061069a565b50602082015181600101559050505050565b828054600181600116156101000203166002900490600052602060002090601f016020900481019282601f1061064157803560ff191683800117855561066f565b8280016001018555821561066f579182015b8281111561066e578235825591602001919060010190610653565b5b50905061067c919061071a565b5090565b604051806040016040528060608152602001600081525090565b828054600181600116156101000203166002900490600052602060002090601f016020900481019282601f106106db57805160ff1916838001178555610709565b82800160010185558215610709579182015b828111156107085782518255916020019190600101906106ed565b5b509050610716919061071a565b5090565b61073c91905b80821115610738576000816000905550600101610720565b5090565b9056fea265627a7a723158208c8e0dc7eb9c3d5a6588ae835056ea2bbfc5578aacceae531d4a9992fd21ee7464736f6c63430005110032";

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
