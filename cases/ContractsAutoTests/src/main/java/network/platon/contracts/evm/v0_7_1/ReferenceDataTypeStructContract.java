package network.platon.contracts.evm.v0_7_1;

import com.alaya.abi.solidity.TypeReference;
import com.alaya.abi.solidity.datatypes.Address;
import com.alaya.abi.solidity.datatypes.Bool;
import com.alaya.abi.solidity.datatypes.Function;
import com.alaya.abi.solidity.datatypes.Type;
import com.alaya.abi.solidity.datatypes.generated.Uint256;
import com.alaya.crypto.Credentials;
import com.alaya.protocol.Web3j;
import com.alaya.protocol.core.RemoteCall;
import com.alaya.tuples.generated.Tuple4;
import com.alaya.tx.Contract;
import com.alaya.tx.TransactionManager;
import com.alaya.tx.gas.GasProvider;
import java.math.BigInteger;
import java.util.Arrays;
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
public class ReferenceDataTypeStructContract extends Contract {
    private static final String BINARY = "608060405234801561001057600080fd5b506102f6806100206000396000f3fe608060405234801561001057600080fd5b50600436106100415760003560e01c80630b1a6855146100465780638e5f3b9c14610091578063cdd9acc1146100dc575b600080fd5b61004e610127565b604051808581526020018473ffffffffffffffffffffffffffffffffffffffff168152602001838152602001821515815260200194505050505060405180910390f35b6100996101a2565b604051808581526020018473ffffffffffffffffffffffffffffffffffffffff168152602001838152602001821515815260200194505050505060405180910390f35b6100e4610211565b604051808581526020018473ffffffffffffffffffffffffffffffffffffffff168152602001838152602001821515815260200194505050505060405180910390f35b600080600080610135610280565b61013d610280565b6040518060800160405280600281526020013373ffffffffffffffffffffffffffffffffffffffff1681526020016019815260200160011515815250905080915081600001518260200151836040015184606001519550955095509550505090919293565b6000806000806101b0610280565b6040518060800160405280600281526020013373ffffffffffffffffffffffffffffffffffffffff16815260200160198152602001600115158152509050806000015181602001518260400151836060015194509450945094505090919293565b60008060008061021f610280565b6040518060800160405280600281526020013373ffffffffffffffffffffffffffffffffffffffff16815260200160198152602001600115158152509050806000015181602001518260400151836060015194509450945094505090919293565b604051806080016040528060008152602001600073ffffffffffffffffffffffffffffffffffffffff16815260200160008152602001600015158152509056fea264697066735822122022432af556adfa14202612780726e74f595cafe9a55f0835f27fa09f6b1b408464736f6c63430007010033";

    public static final String FUNC_INITDATASTRUCTA = "initDataStructA";

    public static final String FUNC_INITDATASTRUCTB = "initDataStructB";

    public static final String FUNC_INITDATASTRUCTC = "initDataStructC";

    protected ReferenceDataTypeStructContract(String contractAddress, Web3j web3j, Credentials credentials, GasProvider contractGasProvider, Long chainId) {
        super(BINARY, contractAddress, web3j, credentials, contractGasProvider, chainId);
    }

    protected ReferenceDataTypeStructContract(String contractAddress, Web3j web3j, TransactionManager transactionManager, GasProvider contractGasProvider, Long chainId) {
        super(BINARY, contractAddress, web3j, transactionManager, contractGasProvider, chainId);
    }

    public RemoteCall<Tuple4<BigInteger, String, BigInteger, Boolean>> initDataStructA() {
        final Function function = new Function(FUNC_INITDATASTRUCTA, 
                Arrays.<Type>asList(), 
                Arrays.<TypeReference<?>>asList(new TypeReference<Uint256>() {}, new TypeReference<Address>() {}, new TypeReference<Uint256>() {}, new TypeReference<Bool>() {}));
        return new RemoteCall<Tuple4<BigInteger, String, BigInteger, Boolean>>(
                new Callable<Tuple4<BigInteger, String, BigInteger, Boolean>>() {
                    @Override
                    public Tuple4<BigInteger, String, BigInteger, Boolean> call() throws Exception {
                        List<Type> results = executeCallMultipleValueReturn(function);
                        return new Tuple4<BigInteger, String, BigInteger, Boolean>(
                                (BigInteger) results.get(0).getValue(), 
                                (String) results.get(1).getValue(), 
                                (BigInteger) results.get(2).getValue(), 
                                (Boolean) results.get(3).getValue());
                    }
                });
    }

    public RemoteCall<Tuple4<BigInteger, String, BigInteger, Boolean>> initDataStructB() {
        final Function function = new Function(FUNC_INITDATASTRUCTB, 
                Arrays.<Type>asList(), 
                Arrays.<TypeReference<?>>asList(new TypeReference<Uint256>() {}, new TypeReference<Address>() {}, new TypeReference<Uint256>() {}, new TypeReference<Bool>() {}));
        return new RemoteCall<Tuple4<BigInteger, String, BigInteger, Boolean>>(
                new Callable<Tuple4<BigInteger, String, BigInteger, Boolean>>() {
                    @Override
                    public Tuple4<BigInteger, String, BigInteger, Boolean> call() throws Exception {
                        List<Type> results = executeCallMultipleValueReturn(function);
                        return new Tuple4<BigInteger, String, BigInteger, Boolean>(
                                (BigInteger) results.get(0).getValue(), 
                                (String) results.get(1).getValue(), 
                                (BigInteger) results.get(2).getValue(), 
                                (Boolean) results.get(3).getValue());
                    }
                });
    }

    public RemoteCall<Tuple4<BigInteger, String, BigInteger, Boolean>> initDataStructC() {
        final Function function = new Function(FUNC_INITDATASTRUCTC, 
                Arrays.<Type>asList(), 
                Arrays.<TypeReference<?>>asList(new TypeReference<Uint256>() {}, new TypeReference<Address>() {}, new TypeReference<Uint256>() {}, new TypeReference<Bool>() {}));
        return new RemoteCall<Tuple4<BigInteger, String, BigInteger, Boolean>>(
                new Callable<Tuple4<BigInteger, String, BigInteger, Boolean>>() {
                    @Override
                    public Tuple4<BigInteger, String, BigInteger, Boolean> call() throws Exception {
                        List<Type> results = executeCallMultipleValueReturn(function);
                        return new Tuple4<BigInteger, String, BigInteger, Boolean>(
                                (BigInteger) results.get(0).getValue(), 
                                (String) results.get(1).getValue(), 
                                (BigInteger) results.get(2).getValue(), 
                                (Boolean) results.get(3).getValue());
                    }
                });
    }

    public static RemoteCall<ReferenceDataTypeStructContract> deploy(Web3j web3j, Credentials credentials, GasProvider contractGasProvider, Long chainId) {
        return deployRemoteCall(ReferenceDataTypeStructContract.class, web3j, credentials, contractGasProvider, BINARY,  "", chainId);
    }

    public static RemoteCall<ReferenceDataTypeStructContract> deploy(Web3j web3j, TransactionManager transactionManager, GasProvider contractGasProvider, Long chainId) {
        return deployRemoteCall(ReferenceDataTypeStructContract.class, web3j, transactionManager, contractGasProvider, BINARY,  "", chainId);
    }

    public static ReferenceDataTypeStructContract load(String contractAddress, Web3j web3j, Credentials credentials, GasProvider contractGasProvider, Long chainId) {
        return new ReferenceDataTypeStructContract(contractAddress, web3j, credentials, contractGasProvider, chainId);
    }

    public static ReferenceDataTypeStructContract load(String contractAddress, Web3j web3j, TransactionManager transactionManager, GasProvider contractGasProvider, Long chainId) {
        return new ReferenceDataTypeStructContract(contractAddress, web3j, transactionManager, contractGasProvider, chainId);
    }
}
