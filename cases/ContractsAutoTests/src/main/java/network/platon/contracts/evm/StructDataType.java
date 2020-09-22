package network.platon.contracts.evm;

import java.math.BigInteger;
import java.util.Arrays;
import java.util.Collections;
import java.util.List;
import java.util.concurrent.Callable;
import org.web3j.abi.TypeReference;
import org.web3j.abi.datatypes.Function;
import org.web3j.abi.datatypes.Type;
import org.web3j.abi.datatypes.generated.Uint256;
import org.web3j.crypto.Credentials;
import org.web3j.protocol.Web3j;
import org.web3j.protocol.core.RemoteCall;
import org.web3j.protocol.core.methods.response.TransactionReceipt;
import org.web3j.tuples.generated.Tuple6;
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
public class StructDataType extends Contract {
    private static final String BINARY = "608060405234801561001057600080fd5b50600260008001600060018152602001908152602001600020819055506004600160000160006003815260200190815260200160002081905550600560018001819055506006600360000181610066919061007f565b50600760046000018161007991906100ab565b506100f8565b8154818355818111156100a6578183600052602060002091820191016100a591906100d7565b5b505050565b8154818355818111156100d2578183600052602060002091820191016100d191906100db565b5b505050565b5090565b6100f591905b808211156100f1576001016100e1565b5090565b90565b610345806101076000396000f3fe608060405234801561001057600080fd5b50600436106100365760003560e01c8063c04062261461003b578063f3ffc4511461007c575b600080fd5b6100436100bd565b60405180878152602001868152602001858152602001848152602001838152602001828152602001965050505050505060405180910390f35b610084610294565b60405180878152602001868152602001858152602001848152602001838152602001828152602001965050505050505060405180910390f35b6000806000806000806100ce6102c3565b6040518060200160405290506100e26102d0565b6040518060200160405280600281525090506100fc6102e3565b6040518060200160405290506101106102f0565b6040518060200160405280600960405190808252806020026020018201604052801561015657816020015b610143610303565b81526020019060019003908161013b5790505b5081525090506101646102f0565b600460405180602001604052908160008201805480602002602001604051908101604052809291908181526020016000905b828210156101c55783829060005260206000200160405180600001604052905081526020019060010190610196565b5050505081525050905084600090505083600160008201518160010155905050826003905050600080016000600181526020019081526020016000205460058190555060018001546006819055508360000151600781905550600360000180549050600881905550816000015151600981905550806000015151600a819055506000800160006001815260200190815260200160002054600180015485600001516003600001805490508560000151518560000151519a509a509a509a509a509a505050505050909192939495565b600080600080600080600554600654600754600854600954600a54955095509550955095509550909192939495565b6040518060200160405290565b6040518060200160405280600081525090565b6040518060200160405290565b6040518060200160405280606081525090565b604051806020016040529056fea265627a7a72315820e52c692e5dd079e7774c0e2d5df379b28c4fcc5cb32f41bd473b9c493fbc074764736f6c634300050d0032";

    public static final String FUNC_GETRUNVALUE = "getRunValue";

    public static final String FUNC_RUN = "run";

    protected StructDataType(String contractAddress, Web3j web3j, Credentials credentials, GasProvider contractGasProvider, Long chainId) {
        super(BINARY, contractAddress, web3j, credentials, contractGasProvider, chainId);
    }

    protected StructDataType(String contractAddress, Web3j web3j, TransactionManager transactionManager, GasProvider contractGasProvider, Long chainId) {
        super(BINARY, contractAddress, web3j, transactionManager, contractGasProvider, chainId);
    }

    public static RemoteCall<StructDataType> deploy(Web3j web3j, Credentials credentials, GasProvider contractGasProvider, Long chainId) {
        return deployRemoteCall(StructDataType.class, web3j, credentials, contractGasProvider, BINARY,  "", chainId);
    }

    public static RemoteCall<StructDataType> deploy(Web3j web3j, TransactionManager transactionManager, GasProvider contractGasProvider, Long chainId) {
        return deployRemoteCall(StructDataType.class, web3j, transactionManager, contractGasProvider, BINARY,  "", chainId);
    }

    public RemoteCall<Tuple6<BigInteger, BigInteger, BigInteger, BigInteger, BigInteger, BigInteger>> getRunValue() {
        final Function function = new Function(FUNC_GETRUNVALUE, 
                Arrays.<Type>asList(), 
                Arrays.<TypeReference<?>>asList(new TypeReference<Uint256>() {}, new TypeReference<Uint256>() {}, new TypeReference<Uint256>() {}, new TypeReference<Uint256>() {}, new TypeReference<Uint256>() {}, new TypeReference<Uint256>() {}));
        return new RemoteCall<Tuple6<BigInteger, BigInteger, BigInteger, BigInteger, BigInteger, BigInteger>>(
                new Callable<Tuple6<BigInteger, BigInteger, BigInteger, BigInteger, BigInteger, BigInteger>>() {
                    @Override
                    public Tuple6<BigInteger, BigInteger, BigInteger, BigInteger, BigInteger, BigInteger> call() throws Exception {
                        List<Type> results = executeCallMultipleValueReturn(function);
                        return new Tuple6<BigInteger, BigInteger, BigInteger, BigInteger, BigInteger, BigInteger>(
                                (BigInteger) results.get(0).getValue(), 
                                (BigInteger) results.get(1).getValue(), 
                                (BigInteger) results.get(2).getValue(), 
                                (BigInteger) results.get(3).getValue(), 
                                (BigInteger) results.get(4).getValue(), 
                                (BigInteger) results.get(5).getValue());
                    }
                });
    }

    public RemoteCall<TransactionReceipt> run() {
        final Function function = new Function(
                FUNC_RUN, 
                Arrays.<Type>asList(), 
                Collections.<TypeReference<?>>emptyList());
        return executeRemoteCallTransaction(function);
    }

    public static StructDataType load(String contractAddress, Web3j web3j, Credentials credentials, GasProvider contractGasProvider, Long chainId) {
        return new StructDataType(contractAddress, web3j, credentials, contractGasProvider, chainId);
    }

    public static StructDataType load(String contractAddress, Web3j web3j, TransactionManager transactionManager, GasProvider contractGasProvider, Long chainId) {
        return new StructDataType(contractAddress, web3j, transactionManager, contractGasProvider, chainId);
    }
}
