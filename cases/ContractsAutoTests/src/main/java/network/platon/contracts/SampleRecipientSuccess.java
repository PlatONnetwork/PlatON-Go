package network.platon.contracts;

import java.math.BigInteger;
import java.util.ArrayList;
import java.util.Arrays;
import java.util.Collections;
import java.util.List;
import org.web3j.abi.EventEncoder;
import org.web3j.abi.TypeReference;
import org.web3j.abi.datatypes.Address;
import org.web3j.abi.datatypes.DynamicBytes;
import org.web3j.abi.datatypes.Event;
import org.web3j.abi.datatypes.Function;
import org.web3j.abi.datatypes.Type;
import org.web3j.abi.datatypes.generated.Uint256;
import org.web3j.crypto.Credentials;
import org.web3j.protocol.Web3j;
import org.web3j.protocol.core.DefaultBlockParameter;
import org.web3j.protocol.core.RemoteCall;
import org.web3j.protocol.core.methods.request.PlatonFilter;
import org.web3j.protocol.core.methods.response.Log;
import org.web3j.protocol.core.methods.response.TransactionReceipt;
import org.web3j.tx.Contract;
import org.web3j.tx.TransactionManager;
import org.web3j.tx.gas.GasProvider;
import rx.Observable;
import rx.functions.Func1;

/**
 * <p>Auto generated code.
 * <p><strong>Do not modify!</strong>
 * <p>Please use the <a href="https://docs.web3j.io/command_line.html">web3j command line tools</a>,
 * or the org.web3j.codegen.SolidityFunctionWrapperGenerator in the 
 * <a href="https://github.com/web3j/web3j/tree/master/codegen">codegen module</a> to update.
 *
 * <p>Generated with web3j version 0.7.5.0.
 */
public class SampleRecipientSuccess extends Contract {
    private static final String BINARY = "608060405234801561001057600080fd5b5061052a806100206000396000f30060806040526004361061006d576000357c0100000000000000000000000000000000000000000000000000000000900463ffffffff1680633fa4f2451461007257806355a373d61461009d578063609d3334146100f45780638f4ffcb114610184578063d5ce338914610237575b600080fd5b34801561007e57600080fd5b5061008761028e565b6040518082815260200191505060405180910390f35b3480156100a957600080fd5b506100b2610294565b604051808273ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200191505060405180910390f35b34801561010057600080fd5b506101096102ba565b6040518080602001828103825283818151815260200191508051906020019080838360005b8381101561014957808201518184015260208101905061012e565b50505050905090810190601f1680156101765780820380516001836020036101000a031916815260200191505b509250505060405180910390f35b34801561019057600080fd5b50610235600480360381019080803573ffffffffffffffffffffffffffffffffffffffff16906020019092919080359060200190929190803573ffffffffffffffffffffffffffffffffffffffff169060200190929190803590602001908201803590602001908080601f0160208091040260200160405190810160405280939291908181526020018383808284378201915050505050509192919290505050610358565b005b34801561024357600080fd5b5061024c610434565b604051808273ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200191505060405180910390f35b60015481565b600260009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1681565b60038054600181600116156101000203166002900480601f0160208091040260200160405190810160405280929190818152602001828054600181600116156101000203166002900480156103505780601f1061032557610100808354040283529160200191610350565b820191906000526020600020905b81548152906001019060200180831161033357829003601f168201915b505050505081565b836000806101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff1602179055508260018190555081600260006101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff16021790555080600390805190602001906103f6929190610459565b507f2db24179b782aab7c5ab64add7f84d4f6c845d0779695371f29be1f658d043cd836040518082815260200191505060405180910390a150505050565b6000809054906101000a900473ffffffffffffffffffffffffffffffffffffffff1681565b828054600181600116156101000203166002900490600052602060002090601f016020900481019282601f1061049a57805160ff19168380011785556104c8565b828001600101855582156104c8579182015b828111156104c75782518255916020019190600101906104ac565b5b5090506104d591906104d9565b5090565b6104fb91905b808211156104f75760008160009055506001016104df565b5090565b905600a165627a7a723058203675e11052c1d1b0d0dcb12f6700a218a009a7ad2e250ad590bd8e193acdde810029";

    public static final String FUNC_VALUE = "value";

    public static final String FUNC_TOKENCONTRACT = "tokenContract";

    public static final String FUNC_EXTRADATA = "extraData";

    public static final String FUNC_RECEIVEAPPROVAL = "receiveApproval";

    public static final String FUNC_FROM = "from";

    public static final Event RECEIVEDAPPROVAL_EVENT = new Event("ReceivedApproval", 
            Arrays.<TypeReference<?>>asList(new TypeReference<Uint256>() {}));
    ;

    @Deprecated
    protected SampleRecipientSuccess(String contractAddress, Web3j web3j, Credentials credentials, BigInteger gasPrice, BigInteger gasLimit) {
        super(BINARY, contractAddress, web3j, credentials, gasPrice, gasLimit);
    }

    protected SampleRecipientSuccess(String contractAddress, Web3j web3j, Credentials credentials, GasProvider contractGasProvider) {
        super(BINARY, contractAddress, web3j, credentials, contractGasProvider);
    }

    @Deprecated
    protected SampleRecipientSuccess(String contractAddress, Web3j web3j, TransactionManager transactionManager, BigInteger gasPrice, BigInteger gasLimit) {
        super(BINARY, contractAddress, web3j, transactionManager, gasPrice, gasLimit);
    }

    protected SampleRecipientSuccess(String contractAddress, Web3j web3j, TransactionManager transactionManager, GasProvider contractGasProvider) {
        super(BINARY, contractAddress, web3j, transactionManager, contractGasProvider);
    }

    public RemoteCall<BigInteger> value() {
        final Function function = new Function(FUNC_VALUE, 
                Arrays.<Type>asList(), 
                Arrays.<TypeReference<?>>asList(new TypeReference<Uint256>() {}));
        return executeRemoteCallSingleValueReturn(function, BigInteger.class);
    }

    public RemoteCall<String> tokenContract() {
        final Function function = new Function(FUNC_TOKENCONTRACT, 
                Arrays.<Type>asList(), 
                Arrays.<TypeReference<?>>asList(new TypeReference<Address>() {}));
        return executeRemoteCallSingleValueReturn(function, String.class);
    }

    public RemoteCall<byte[]> extraData() {
        final Function function = new Function(FUNC_EXTRADATA, 
                Arrays.<Type>asList(), 
                Arrays.<TypeReference<?>>asList(new TypeReference<DynamicBytes>() {}));
        return executeRemoteCallSingleValueReturn(function, byte[].class);
    }

    public RemoteCall<TransactionReceipt> receiveApproval(String _from, BigInteger _value, String _tokenContract, byte[] _extraData) {
        final Function function = new Function(
                FUNC_RECEIVEAPPROVAL, 
                Arrays.<Type>asList(new org.web3j.abi.datatypes.Address(_from), 
                new org.web3j.abi.datatypes.generated.Uint256(_value), 
                new org.web3j.abi.datatypes.Address(_tokenContract), 
                new org.web3j.abi.datatypes.DynamicBytes(_extraData)), 
                Collections.<TypeReference<?>>emptyList());
        return executeRemoteCallTransaction(function);
    }

    public RemoteCall<String> from() {
        final Function function = new Function(FUNC_FROM, 
                Arrays.<Type>asList(), 
                Arrays.<TypeReference<?>>asList(new TypeReference<Address>() {}));
        return executeRemoteCallSingleValueReturn(function, String.class);
    }

    public List<ReceivedApprovalEventResponse> getReceivedApprovalEvents(TransactionReceipt transactionReceipt) {
        List<Contract.EventValuesWithLog> valueList = extractEventParametersWithLog(RECEIVEDAPPROVAL_EVENT, transactionReceipt);
        ArrayList<ReceivedApprovalEventResponse> responses = new ArrayList<ReceivedApprovalEventResponse>(valueList.size());
        for (Contract.EventValuesWithLog eventValues : valueList) {
            ReceivedApprovalEventResponse typedResponse = new ReceivedApprovalEventResponse();
            typedResponse.log = eventValues.getLog();
            typedResponse._value = (BigInteger) eventValues.getNonIndexedValues().get(0).getValue();
            responses.add(typedResponse);
        }
        return responses;
    }

    public Observable<ReceivedApprovalEventResponse> receivedApprovalEventObservable(PlatonFilter filter) {
        return web3j.platonLogObservable(filter).map(new Func1<Log, ReceivedApprovalEventResponse>() {
            @Override
            public ReceivedApprovalEventResponse call(Log log) {
                Contract.EventValuesWithLog eventValues = extractEventParametersWithLog(RECEIVEDAPPROVAL_EVENT, log);
                ReceivedApprovalEventResponse typedResponse = new ReceivedApprovalEventResponse();
                typedResponse.log = log;
                typedResponse._value = (BigInteger) eventValues.getNonIndexedValues().get(0).getValue();
                return typedResponse;
            }
        });
    }

    public Observable<ReceivedApprovalEventResponse> receivedApprovalEventObservable(DefaultBlockParameter startBlock, DefaultBlockParameter endBlock) {
        PlatonFilter filter = new PlatonFilter(startBlock, endBlock, getContractAddress());
        filter.addSingleTopic(EventEncoder.encode(RECEIVEDAPPROVAL_EVENT));
        return receivedApprovalEventObservable(filter);
    }

    public static RemoteCall<SampleRecipientSuccess> deploy(Web3j web3j, Credentials credentials, GasProvider contractGasProvider) {
        return deployRemoteCall(SampleRecipientSuccess.class, web3j, credentials, contractGasProvider, BINARY, "");
    }

    @Deprecated
    public static RemoteCall<SampleRecipientSuccess> deploy(Web3j web3j, Credentials credentials, BigInteger gasPrice, BigInteger gasLimit) {
        return deployRemoteCall(SampleRecipientSuccess.class, web3j, credentials, gasPrice, gasLimit, BINARY, "");
    }

    public static RemoteCall<SampleRecipientSuccess> deploy(Web3j web3j, TransactionManager transactionManager, GasProvider contractGasProvider) {
        return deployRemoteCall(SampleRecipientSuccess.class, web3j, transactionManager, contractGasProvider, BINARY, "");
    }

    @Deprecated
    public static RemoteCall<SampleRecipientSuccess> deploy(Web3j web3j, TransactionManager transactionManager, BigInteger gasPrice, BigInteger gasLimit) {
        return deployRemoteCall(SampleRecipientSuccess.class, web3j, transactionManager, gasPrice, gasLimit, BINARY, "");
    }

    @Deprecated
    public static SampleRecipientSuccess load(String contractAddress, Web3j web3j, Credentials credentials, BigInteger gasPrice, BigInteger gasLimit) {
        return new SampleRecipientSuccess(contractAddress, web3j, credentials, gasPrice, gasLimit);
    }

    @Deprecated
    public static SampleRecipientSuccess load(String contractAddress, Web3j web3j, TransactionManager transactionManager, BigInteger gasPrice, BigInteger gasLimit) {
        return new SampleRecipientSuccess(contractAddress, web3j, transactionManager, gasPrice, gasLimit);
    }

    public static SampleRecipientSuccess load(String contractAddress, Web3j web3j, Credentials credentials, GasProvider contractGasProvider) {
        return new SampleRecipientSuccess(contractAddress, web3j, credentials, contractGasProvider);
    }

    public static SampleRecipientSuccess load(String contractAddress, Web3j web3j, TransactionManager transactionManager, GasProvider contractGasProvider) {
        return new SampleRecipientSuccess(contractAddress, web3j, transactionManager, contractGasProvider);
    }

    public static class ReceivedApprovalEventResponse {
        public Log log;

        public BigInteger _value;
    }
}
