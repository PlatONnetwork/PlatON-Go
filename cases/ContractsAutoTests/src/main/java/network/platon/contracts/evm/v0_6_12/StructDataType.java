package network.platon.contracts.evm.v0_6_12;

import com.alaya.abi.solidity.TypeReference;
import com.alaya.abi.solidity.datatypes.Function;
import com.alaya.abi.solidity.datatypes.Type;
import com.alaya.abi.solidity.datatypes.generated.Uint256;
import com.alaya.crypto.Credentials;
import com.alaya.protocol.Web3j;
import com.alaya.protocol.core.RemoteCall;
import com.alaya.protocol.core.methods.response.TransactionReceipt;
import com.alaya.tuples.generated.Tuple6;
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
public class StructDataType extends Contract {
    private static final String BINARY = "608060405234801561001057600080fd5b5060026000800160006001815260200190815260200160002081905550600460016000016000600381526020019081526020016000208190555060056001800181905550610350806100636000396000f3fe608060405234801561001057600080fd5b50600436106100365760003560e01c8063c04062261461003b578063f3ffc4511461007c575b600080fd5b6100436100bd565b60405180878152602001868152602001858152602001848152602001838152602001828152602001965050505050505060405180910390f35b6100846102ab565b60405180878152602001868152602001858152602001848152602001838152602001828152602001965050505050505060405180910390f35b6000806000806000806100ce6102da565b6040518060200160405290506100e26102e7565b6040518060200160405280600281525090506100fc6102fa565b604051806020016040529050610110610307565b6040518060200160405280600967ffffffffffffffff8111801561013357600080fd5b5060405190808252806020026020018201604052801561016d57816020015b61015a6102da565b8152602001906001900390816101525790505b50815250905061017b610307565b600460405180602001604052908160008201805480602002602001604051908101604052809291908181526020016000905b828210156101dc57838290600052602060002001604051806000016040529050815260200190600101906101ad565b5050505081525050905084600090505083600160008201518160010155905050826003905050600080016000600181526020019081526020016000205460058190555060018001546006819055508360000151600781905550600360000180549050600881905550816000015151600981905550806000015151600a819055506000800160006001815260200190815260200160002054600180015485600001516003600001805490508560000151518560000151519a509a509a509a509a509a505050505050909192939495565b600080600080600080600554600654600754600854600954600a54955095509550955095509550909192939495565b6040518060200160405290565b6040518060200160405280600081525090565b6040518060200160405290565b604051806020016040528060608152509056fea26469706673582212209a245bf8946ab076f3b159a2e323e67dd062e3590e0b02e061c848a26546dff664736f6c634300060c0033";

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
