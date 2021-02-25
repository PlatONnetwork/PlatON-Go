package network.platon.contracts.evm;

import java.math.BigInteger;
import java.util.ArrayList;
import java.util.Arrays;
import java.util.Collections;
import java.util.List;
import org.web3j.abi.EventEncoder;
import org.web3j.abi.TypeReference;
import org.web3j.abi.datatypes.Address;
import org.web3j.abi.datatypes.Event;
import org.web3j.abi.datatypes.Function;
import org.web3j.abi.datatypes.Type;
import org.web3j.abi.datatypes.generated.Int8;
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
 * <p>Please use the <a href="https://github.com/PlatONnetwork/client-sdk-java/releases">platon-web3j command line tools</a>,
 * or the org.web3j.codegen.SolidityFunctionWrapperGenerator in the 
 * <a href="https://github.com/PlatONnetwork/client-sdk-java/tree/master/codegen">codegen module</a> to update.
 *
 * <p>Generated with web3j version 0.13.1.5.
 */
public class Eventer extends Contract {
    private static final String BINARY = "6080604052348015600f57600080fd5b5060f48061001e6000396000f300608060405260043610603f576000357c0100000000000000000000000000000000000000000000000000000000900463ffffffff168063bf819c20146044575b600080fd5b348015604f57600080fd5b5060566058565b005b7ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffd7ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe7f8f50d21be7587a4814a9d4c10b7c8d1eea6389adbd44cb59ddaba790fd2ecbbd60405160405180910390a35600a165627a7a72305820d7159162cd1d095e14d10ff44c305d12521e8bd9b9ce43e3c14a3eaf6849d7110029";

    public static final String FUNC_GETEVENT = "getEvent";

    public static final Event TESTINT8_EVENT = new Event("TestInt8", 
            Arrays.<TypeReference<?>>asList(new TypeReference<Int8>(true) {}, new TypeReference<Int8>(true) {}));
    ;

    public static final Event ANONEVENT_EVENT = new Event("AnonEvent", 
            Arrays.<TypeReference<?>>asList(new TypeReference<Address>() {}, new TypeReference<Address>() {}));
    ;

    protected Eventer(String contractAddress, Web3j web3j, Credentials credentials, GasProvider contractGasProvider, Long chainId) {
        super(BINARY, contractAddress, web3j, credentials, contractGasProvider, chainId);
    }

    protected Eventer(String contractAddress, Web3j web3j, TransactionManager transactionManager, GasProvider contractGasProvider, Long chainId) {
        super(BINARY, contractAddress, web3j, transactionManager, contractGasProvider, chainId);
    }

    public RemoteCall<TransactionReceipt> getEvent() {
        final Function function = new Function(
                FUNC_GETEVENT, 
                Arrays.<Type>asList(), 
                Collections.<TypeReference<?>>emptyList());
        return executeRemoteCallTransaction(function);
    }

    public List<TestInt8EventResponse> getTestInt8Events(TransactionReceipt transactionReceipt) {
        List<EventValuesWithLog> valueList = extractEventParametersWithLog(TESTINT8_EVENT, transactionReceipt);
        ArrayList<TestInt8EventResponse> responses = new ArrayList<TestInt8EventResponse>(valueList.size());
        for (EventValuesWithLog eventValues : valueList) {
            TestInt8EventResponse typedResponse = new TestInt8EventResponse();
            typedResponse.log = eventValues.getLog();
            typedResponse.out1 = (BigInteger) eventValues.getIndexedValues().get(0).getValue();
            typedResponse.out2 = (BigInteger) eventValues.getIndexedValues().get(1).getValue();
            responses.add(typedResponse);
        }
        return responses;
    }

    public Observable<TestInt8EventResponse> testInt8EventObservable(PlatonFilter filter) {
        return web3j.platonLogObservable(filter).map(new Func1<Log, TestInt8EventResponse>() {
            @Override
            public TestInt8EventResponse call(Log log) {
                EventValuesWithLog eventValues = extractEventParametersWithLog(TESTINT8_EVENT, log);
                TestInt8EventResponse typedResponse = new TestInt8EventResponse();
                typedResponse.log = log;
                typedResponse.out1 = (BigInteger) eventValues.getIndexedValues().get(0).getValue();
                typedResponse.out2 = (BigInteger) eventValues.getIndexedValues().get(1).getValue();
                return typedResponse;
            }
        });
    }

    public Observable<TestInt8EventResponse> testInt8EventObservable(DefaultBlockParameter startBlock, DefaultBlockParameter endBlock) {
        PlatonFilter filter = new PlatonFilter(startBlock, endBlock, getContractAddress());
        filter.addSingleTopic(EventEncoder.encode(TESTINT8_EVENT));
        return testInt8EventObservable(filter);
    }

    public List<AnonEventEventResponse> getAnonEventEvents(TransactionReceipt transactionReceipt) {
        List<EventValuesWithLog> valueList = extractEventParametersWithLog(ANONEVENT_EVENT, transactionReceipt);
        ArrayList<AnonEventEventResponse> responses = new ArrayList<AnonEventEventResponse>(valueList.size());
        for (EventValuesWithLog eventValues : valueList) {
            AnonEventEventResponse typedResponse = new AnonEventEventResponse();
            typedResponse.log = eventValues.getLog();
            typedResponse.param0 = (String) eventValues.getNonIndexedValues().get(0).getValue();
            typedResponse.param1 = (String) eventValues.getNonIndexedValues().get(1).getValue();
            responses.add(typedResponse);
        }
        return responses;
    }

    public Observable<AnonEventEventResponse> anonEventEventObservable(PlatonFilter filter) {
        return web3j.platonLogObservable(filter).map(new Func1<Log, AnonEventEventResponse>() {
            @Override
            public AnonEventEventResponse call(Log log) {
                EventValuesWithLog eventValues = extractEventParametersWithLog(ANONEVENT_EVENT, log);
                AnonEventEventResponse typedResponse = new AnonEventEventResponse();
                typedResponse.log = log;
                typedResponse.param0 = (String) eventValues.getNonIndexedValues().get(0).getValue();
                typedResponse.param1 = (String) eventValues.getNonIndexedValues().get(1).getValue();
                return typedResponse;
            }
        });
    }

    public Observable<AnonEventEventResponse> anonEventEventObservable(DefaultBlockParameter startBlock, DefaultBlockParameter endBlock) {
        PlatonFilter filter = new PlatonFilter(startBlock, endBlock, getContractAddress());
        filter.addSingleTopic(EventEncoder.encode(ANONEVENT_EVENT));
        return anonEventEventObservable(filter);
    }

    public static RemoteCall<Eventer> deploy(Web3j web3j, Credentials credentials, GasProvider contractGasProvider, Long chainId) {
        return deployRemoteCall(Eventer.class, web3j, credentials, contractGasProvider, BINARY,  "", chainId);
    }

    public static RemoteCall<Eventer> deploy(Web3j web3j, TransactionManager transactionManager, GasProvider contractGasProvider, Long chainId) {
        return deployRemoteCall(Eventer.class, web3j, transactionManager, contractGasProvider, BINARY,  "", chainId);
    }

    public static Eventer load(String contractAddress, Web3j web3j, Credentials credentials, GasProvider contractGasProvider, Long chainId) {
        return new Eventer(contractAddress, web3j, credentials, contractGasProvider, chainId);
    }

    public static Eventer load(String contractAddress, Web3j web3j, TransactionManager transactionManager, GasProvider contractGasProvider, Long chainId) {
        return new Eventer(contractAddress, web3j, transactionManager, contractGasProvider, chainId);
    }

    public static class TestInt8EventResponse {
        public Log log;

        public BigInteger out1;

        public BigInteger out2;
    }

    public static class AnonEventEventResponse {
        public Log log;

        public String param0;

        public String param1;
    }
}
