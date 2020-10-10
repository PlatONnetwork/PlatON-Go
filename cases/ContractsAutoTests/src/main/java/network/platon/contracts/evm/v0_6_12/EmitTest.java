package network.platon.contracts.evm.v0_6_12;

import com.alaya.abi.solidity.EventEncoder;
import com.alaya.abi.solidity.TypeReference;
import com.alaya.abi.solidity.datatypes.Address;
import com.alaya.abi.solidity.datatypes.Event;
import com.alaya.abi.solidity.datatypes.Function;
import com.alaya.abi.solidity.datatypes.Type;
import com.alaya.abi.solidity.datatypes.generated.Uint256;
import com.alaya.crypto.Credentials;
import com.alaya.protocol.Web3j;
import com.alaya.protocol.core.DefaultBlockParameter;
import com.alaya.protocol.core.RemoteCall;
import com.alaya.protocol.core.methods.request.PlatonFilter;
import com.alaya.protocol.core.methods.response.Log;
import com.alaya.protocol.core.methods.response.TransactionReceipt;
import com.alaya.tx.Contract;
import com.alaya.tx.TransactionManager;
import com.alaya.tx.gas.GasProvider;
import java.math.BigInteger;
import java.util.ArrayList;
import java.util.Arrays;
import java.util.Collections;
import java.util.List;
import rx.Observable;
import rx.functions.Func1;

/**
 * <p>Auto generated code.
 * <p><strong>Do not modify!</strong>
 * <p>Please use the <a href="https://github.com/PlatONnetwork/client-sdk-java/releases">platon-web3j command line tools</a>,
 * or the com.alaya.codegen.SolidityFunctionWrapperGenerator in the 
 * <a href="https://github.com/PlatONnetwork/client-sdk-java/tree/master/codegen">codegen module</a> to update.
 *
 * <p>Generated with web3j version 0.13.2.0.
 */
public class EmitTest extends Contract {
    private static final String BINARY = "6080604052348015600f57600080fd5b5060b68061001e6000396000f3fe608060405260043610601c5760003560e01c80634f9d719e146021575b600080fd5b60276029565b005b7fb0333e0e3a6b99318e4e2e0d7e5e5f93646f9cbf62da1587955a4092bf7df6e73334604051808373ffffffffffffffffffffffffffffffffffffffff1681526020018281526020019250505060405180910390a156fea26469706673582212208deb0b9acd365446c557ad807a859edabcc960913ce4d92e000b04ea5fa2459764736f6c634300060c0033";

    public static final String FUNC_TESTEVENT = "testEvent";

    public static final Event EVENTNAME_EVENT = new Event("EventName", 
            Arrays.<TypeReference<?>>asList(new TypeReference<Address>() {}, new TypeReference<Uint256>() {}));
    ;

    protected EmitTest(String contractAddress, Web3j web3j, Credentials credentials, GasProvider contractGasProvider, Long chainId) {
        super(BINARY, contractAddress, web3j, credentials, contractGasProvider, chainId);
    }

    protected EmitTest(String contractAddress, Web3j web3j, TransactionManager transactionManager, GasProvider contractGasProvider, Long chainId) {
        super(BINARY, contractAddress, web3j, transactionManager, contractGasProvider, chainId);
    }

    public List<EventNameEventResponse> getEventNameEvents(TransactionReceipt transactionReceipt) {
        List<EventValuesWithLog> valueList = extractEventParametersWithLog(EVENTNAME_EVENT, transactionReceipt);
        ArrayList<EventNameEventResponse> responses = new ArrayList<EventNameEventResponse>(valueList.size());
        for (EventValuesWithLog eventValues : valueList) {
            EventNameEventResponse typedResponse = new EventNameEventResponse();
            typedResponse.log = eventValues.getLog();
            typedResponse.bidder = (String) eventValues.getNonIndexedValues().get(0).getValue();
            typedResponse.amount = (BigInteger) eventValues.getNonIndexedValues().get(1).getValue();
            responses.add(typedResponse);
        }
        return responses;
    }

    public Observable<EventNameEventResponse> eventNameEventObservable(PlatonFilter filter) {
        return web3j.platonLogObservable(filter).map(new Func1<Log, EventNameEventResponse>() {
            @Override
            public EventNameEventResponse call(Log log) {
                EventValuesWithLog eventValues = extractEventParametersWithLog(EVENTNAME_EVENT, log);
                EventNameEventResponse typedResponse = new EventNameEventResponse();
                typedResponse.log = log;
                typedResponse.bidder = (String) eventValues.getNonIndexedValues().get(0).getValue();
                typedResponse.amount = (BigInteger) eventValues.getNonIndexedValues().get(1).getValue();
                return typedResponse;
            }
        });
    }

    public Observable<EventNameEventResponse> eventNameEventObservable(DefaultBlockParameter startBlock, DefaultBlockParameter endBlock) {
        PlatonFilter filter = new PlatonFilter(startBlock, endBlock, getContractAddress());
        filter.addSingleTopic(EventEncoder.encode(EVENTNAME_EVENT));
        return eventNameEventObservable(filter);
    }

    public RemoteCall<TransactionReceipt> testEvent() {
        final Function function = new Function(
                FUNC_TESTEVENT, 
                Arrays.<Type>asList(), 
                Collections.<TypeReference<?>>emptyList());
        return executeRemoteCallTransaction(function);
    }

    public static RemoteCall<EmitTest> deploy(Web3j web3j, Credentials credentials, GasProvider contractGasProvider, Long chainId) {
        return deployRemoteCall(EmitTest.class, web3j, credentials, contractGasProvider, BINARY,  "", chainId);
    }

    public static RemoteCall<EmitTest> deploy(Web3j web3j, TransactionManager transactionManager, GasProvider contractGasProvider, Long chainId) {
        return deployRemoteCall(EmitTest.class, web3j, transactionManager, contractGasProvider, BINARY,  "", chainId);
    }

    public static EmitTest load(String contractAddress, Web3j web3j, Credentials credentials, GasProvider contractGasProvider, Long chainId) {
        return new EmitTest(contractAddress, web3j, credentials, contractGasProvider, chainId);
    }

    public static EmitTest load(String contractAddress, Web3j web3j, TransactionManager transactionManager, GasProvider contractGasProvider, Long chainId) {
        return new EmitTest(contractAddress, web3j, transactionManager, contractGasProvider, chainId);
    }

    public static class EventNameEventResponse {
        public Log log;

        public String bidder;

        public BigInteger amount;
    }
}
