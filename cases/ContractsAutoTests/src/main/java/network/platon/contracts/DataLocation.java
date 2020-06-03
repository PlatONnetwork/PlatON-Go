package network.platon.contracts;

import java.math.BigInteger;
import java.util.Arrays;
import java.util.Collections;
import java.util.List;
import java.util.concurrent.Callable;
import org.web3j.abi.TypeReference;
import org.web3j.abi.datatypes.DynamicBytes;
import org.web3j.abi.datatypes.Function;
import org.web3j.abi.datatypes.Type;
import org.web3j.abi.datatypes.Utf8String;
import org.web3j.abi.datatypes.generated.Uint256;
import org.web3j.crypto.Credentials;
import org.web3j.protocol.Web3j;
import org.web3j.protocol.core.RemoteCall;
import org.web3j.protocol.core.methods.response.TransactionReceipt;
import org.web3j.tuples.generated.Tuple2;
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
 * <p>Generated with web3j version 0.13.0.7.
 */
public class DataLocation extends Contract {
    private static final String BINARY = "608060405234801561001057600080fd5b506107b5806100206000396000f3fe608060405260043610610062576000357c0100000000000000000000000000000000000000000000000000000000900463ffffffff1680630bcd3b3314610067578063246982c4146100f75780633ca8b1a7146101b2578063a1715deb146102b1575b600080fd5b34801561007357600080fd5b5061007c6103a5565b6040518080602001828103825283818151815260200191508051906020019080838360005b838110156100bc5780820151818401526020810190506100a1565b50505050905090810190601f1680156100e95780820380516001836020036101000a031916815260200191505b509250505060405180910390f35b34801561010357600080fd5b506101306004803603602081101561011a57600080fd5b8101908080359060200190929190505050610447565b6040518080602001838152602001828103825284818151815260200191508051906020019080838360005b8381101561017657808201518184015260208101905061015b565b50505050905090810190601f1680156101a35780820380516001836020036101000a031916815260200191505b50935050505060405180910390f35b3480156101be57600080fd5b50610236600480360360208110156101d557600080fd5b81019080803590602001906401000000008111156101f257600080fd5b82018360208201111561020457600080fd5b8035906020019184600183028401116401000000008311171561022657600080fd5b909192939192939050505061051b565b6040518080602001828103825283818151815260200191508051906020019080838360005b8381101561027657808201518184015260208101905061025b565b50505050905090810190601f1680156102a35780820380516001836020036101000a031916815260200191505b509250505060405180910390f35b3480156102bd57600080fd5b5061038b600480360360608110156102d457600080fd5b8101908080359060200190929190803590602001906401000000008111156102fb57600080fd5b82018360208201111561030d57600080fd5b8035906020019184600183028401116401000000008311171561032f57600080fd5b91908080601f016020809104026020016040519081016040528093929190818152602001838380828437600081840152601f19601f820116905080830192505050505050509192919290803590602001909291905050506105d2565b604051808215151515815260200191505060405180910390f35b606060018054600181600116156101000203166002900480601f01602080910402602001604051908101604052809291908181526020018280546001816001161561010002031660029004801561043d5780601f106104125761010080835404028352916020019161043d565b820191906000526020600020905b81548152906001019060200180831161042057829003601f168201915b5050505050905090565b6060600080600084815260200190815260200160002060000160008085815260200190815260200160002060010154818054600181600116156101000203166002900480601f01602080910402602001604051908101604052809291908181526020018280546001816001161561010002031660029004801561050b5780601f106104e05761010080835404028352916020019161050b565b820191906000526020600020905b8154815290600101906020018083116104ee57829003601f168201915b5050505050915091509150915091565b606082826001919061052e92919061064a565b5060018054600181600116156101000203166002900480601f0160208091040260200160405190810160405280929190818152602001828054600181600116156101000203166002900480156105c55780601f1061059a576101008083540402835291602001916105c5565b820191906000526020600020905b8154815290600101906020018083116105a857829003601f168201915b5050505050905092915050565b60006105dc6106ca565b60408051908101604052808581526020018481525090506105fd8186610609565b60019150509392505050565b8160008083815260200190815260200160002060008201518160000190805190602001906106389291906106e4565b50602082015181600101559050505050565b828054600181600116156101000203166002900490600052602060002090601f016020900481019282601f1061068b57803560ff19168380011785556106b9565b828001600101855582156106b9579182015b828111156106b857823582559160200191906001019061069d565b5b5090506106c69190610764565b5090565b604080519081016040528060608152602001600081525090565b828054600181600116156101000203166002900490600052602060002090601f016020900481019282601f1061072557805160ff1916838001178555610753565b82800160010185558215610753579182015b82811115610752578251825591602001919060010190610737565b5b5090506107609190610764565b5090565b61078691905b8082111561078257600081600090555060010161076a565b5090565b9056fea165627a7a723058200bac4d1bbc74a688607ca0e899c056206da7e259b114fd4e45f190bff5a889ac0029";

    public static final String FUNC_GETBYTES = "getBytes";

    public static final String FUNC_GETPERSON = "getPerson";

    public static final String FUNC_TESTBYTES = "testBytes";

    public static final String FUNC_SAVEPERSON = "savePerson";

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
                Arrays.<Type>asList(new org.web3j.abi.datatypes.generated.Uint256(_id)), 
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

    public RemoteCall<TransactionReceipt> testBytes(byte[] _data) {
        final Function function = new Function(
                FUNC_TESTBYTES, 
                Arrays.<Type>asList(new org.web3j.abi.datatypes.DynamicBytes(_data)), 
                Collections.<TypeReference<?>>emptyList());
        return executeRemoteCallTransaction(function);
    }

    public RemoteCall<TransactionReceipt> savePerson(BigInteger _id, String _name, BigInteger _age) {
        final Function function = new Function(
                FUNC_SAVEPERSON, 
                Arrays.<Type>asList(new org.web3j.abi.datatypes.generated.Uint256(_id), 
                new org.web3j.abi.datatypes.Utf8String(_name), 
                new org.web3j.abi.datatypes.generated.Uint256(_age)), 
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
