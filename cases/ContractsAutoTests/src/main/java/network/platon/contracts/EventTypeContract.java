package network.platon.contracts;

import java.math.BigInteger;
import java.util.ArrayList;
import java.util.Arrays;
import java.util.Collections;
import java.util.List;
import org.web3j.abi.EventEncoder;
import org.web3j.abi.TypeReference;
import org.web3j.abi.datatypes.DynamicArray;
import org.web3j.abi.datatypes.Event;
import org.web3j.abi.datatypes.Function;
import org.web3j.abi.datatypes.Type;
import org.web3j.abi.datatypes.Utf8String;
import org.web3j.abi.datatypes.generated.StaticArray2;
import org.web3j.abi.datatypes.generated.Uint256;
import org.web3j.abi.datatypes.generated.Uint8;
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
 * <p>Generated with web3j version 0.13.0.7.
 */
public class EventTypeContract extends Contract {
    private static final String BINARY = "60806040526040518060a00160405280600060ff168152602001600160ff168152602001600260ff168152602001600360ff168152602001600460ff16815250600090600561004f929190610177565b506040518060c001604052806040518060400160405280600060ff168152602001600060ff1681525081526020016040518060400160405280600060ff168152602001600160ff1681525081526020016040518060400160405280600060ff168152602001600260ff1681525081526020016040518060400160405280600060ff168152602001600360ff1681525081526020016040518060400160405280600060ff168152602001600460ff1681525081526020016040518060400160405280600060ff168152602001600560ff1681525081525060019060066101359291906101c9565b5060405180608001604052806058815260200161061b6058913960029080519060200190610164929190610227565b5034801561017157600080fd5b5061034b565b8280548282559060005260206000209081019282156101b8579160200282015b828111156101b7578251829060ff16905591602001919060010190610197565b5b5090506101c591906102a7565b5090565b828054828255906000526020600020906002028101928215610216579160200282015b82811115610215578251829060026102059291906102cc565b50916020019190600201906101ec565b5b5090506102239190610311565b5090565b828054600181600116156101000203166002900490600052602060002090601f016020900481019282601f1061026857805160ff1916838001178555610296565b82800160010185558215610296579182015b8281111561029557825182559160200191906001019061027a565b5b5090506102a391906102a7565b5090565b6102c991905b808211156102c55760008160009055506001016102ad565b5090565b90565b8260028101928215610300579160200282015b828111156102ff578251829060ff169055916020019190600101906102df565b5b50905061030d91906102a7565b5090565b61033a91905b80821115610336576000818161032d919061033d565b50600201610317565b5090565b90565b506000815560010160009055565b6102c18061035a6000396000f3fe608060405234801561001057600080fd5b506004361061004c5760003560e01c80633ac559931461005157806343ae41d81461005b578063b05dfcf814610065578063bbd847af1461006f575b600080fd5b610059610079565b005b610063610135565b005b61006d61017d565b005b6100776101f5565b005b7f617cf8a4400dd7963ed519ebe655a16e8da1282bb8fea36a21f634af912f54ab600260405180806020018281038252838181546001816001161561010002031660029004815260200191508054600181600116156101000203166002900480156101255780601f106100fa57610100808354040283529160200191610125565b820191906000526020600020905b81548152906001019060200180831161010857829003601f168201915b50509250505060405180910390a1565b7fde7a62815e0b38238b6211179d7d98017a99227a90823b0f44227e81dd3ad9c260006040518082600381111561016857fe5b60ff16815260200191505060405180910390a1565b7f38a323fa24260bbb8b86f61cd1d8c1900024088af6d08eda9e2d793da33c1b586000604051808060200182810382528381815481526020019150805480156101e557602002820191906000526020600020905b8154815260200190600101908083116101d1575b50509250505060405180910390a1565b7f406715adbc90cbc793dcd5707190ad1390229b2a75cf5b5ca228b518ae52de9a60016040518080602001828103825283818154815260200191508054801561027c57602002820191906000526020600020905b816002801561026d576020028201915b815481526020019060010190808311610259575b50509060020190808311610249575b50509250505060405180910390a156fea265627a7a72315820859e60af90ac286a9684a075e692aaa0fa1aa31f98dc1289d596f9e55e66b8ab64736f6c634300050d003231323334353637383930303937383635343332313132333435363738393030393837363534333231313233343536373839303039373634333534363636363633323432343434343434343434343735383331353436383536";

    public static final String FUNC_TESTENUM = "testEnum";

    public static final String FUNC_TESTONEDIMENSIONALARRAY = "testOneDimensionalArray";

    public static final String FUNC_TESTSTR = "testStr";

    public static final String FUNC_TESTTWODIMENSIONALARRAY = "testTwoDimensionalArray";

    public static final Event ENUMEVENT_EVENT = new Event("EnumEvent", 
            Arrays.<TypeReference<?>>asList(new TypeReference<Uint8>() {}));
    ;

    public static final Event ONEDIMENSIONALARRAYEVENT_EVENT = new Event("OneDimensionalArrayEvent", 
            Arrays.<TypeReference<?>>asList(new TypeReference<DynamicArray<Uint256>>() {}));
    ;

    public static final Event STRINGEVENT_EVENT = new Event("StringEvent", 
            Arrays.<TypeReference<?>>asList(new TypeReference<Utf8String>() {}));
    ;

    public static final Event TWODIMENSIONALARRAYEVENT_EVENT = new Event("TwoDimensionalArrayEvent", 
            Arrays.<TypeReference<?>>asList(new TypeReference<DynamicArray<StaticArray2<Uint256>>>() {}));
    ;

    protected EventTypeContract(String contractAddress, Web3j web3j, Credentials credentials, GasProvider contractGasProvider, Long chainId) {
        super(BINARY, contractAddress, web3j, credentials, contractGasProvider, chainId);
    }

    protected EventTypeContract(String contractAddress, Web3j web3j, TransactionManager transactionManager, GasProvider contractGasProvider, Long chainId) {
        super(BINARY, contractAddress, web3j, transactionManager, contractGasProvider, chainId);
    }

    public List<EnumEventEventResponse> getEnumEventEvents(TransactionReceipt transactionReceipt) {
        List<Contract.EventValuesWithLog> valueList = extractEventParametersWithLog(ENUMEVENT_EVENT, transactionReceipt);
        ArrayList<EnumEventEventResponse> responses = new ArrayList<EnumEventEventResponse>(valueList.size());
        for (Contract.EventValuesWithLog eventValues : valueList) {
            EnumEventEventResponse typedResponse = new EnumEventEventResponse();
            typedResponse.log = eventValues.getLog();
            typedResponse.choices = (BigInteger) eventValues.getNonIndexedValues().get(0).getValue();
            responses.add(typedResponse);
        }
        return responses;
    }

    public Observable<EnumEventEventResponse> enumEventEventObservable(PlatonFilter filter) {
        return web3j.platonLogObservable(filter).map(new Func1<Log, EnumEventEventResponse>() {
            @Override
            public EnumEventEventResponse call(Log log) {
                Contract.EventValuesWithLog eventValues = extractEventParametersWithLog(ENUMEVENT_EVENT, log);
                EnumEventEventResponse typedResponse = new EnumEventEventResponse();
                typedResponse.log = log;
                typedResponse.choices = (BigInteger) eventValues.getNonIndexedValues().get(0).getValue();
                return typedResponse;
            }
        });
    }

    public Observable<EnumEventEventResponse> enumEventEventObservable(DefaultBlockParameter startBlock, DefaultBlockParameter endBlock) {
        PlatonFilter filter = new PlatonFilter(startBlock, endBlock, getContractAddress());
        filter.addSingleTopic(EventEncoder.encode(ENUMEVENT_EVENT));
        return enumEventEventObservable(filter);
    }

    public List<OneDimensionalArrayEventEventResponse> getOneDimensionalArrayEventEvents(TransactionReceipt transactionReceipt) {
        List<Contract.EventValuesWithLog> valueList = extractEventParametersWithLog(ONEDIMENSIONALARRAYEVENT_EVENT, transactionReceipt);
        ArrayList<OneDimensionalArrayEventEventResponse> responses = new ArrayList<OneDimensionalArrayEventEventResponse>(valueList.size());
        for (Contract.EventValuesWithLog eventValues : valueList) {
            OneDimensionalArrayEventEventResponse typedResponse = new OneDimensionalArrayEventEventResponse();
            typedResponse.log = eventValues.getLog();
            typedResponse.array = (List<BigInteger>) eventValues.getNonIndexedValues().get(0).getValue();
            responses.add(typedResponse);
        }
        return responses;
    }

    public Observable<OneDimensionalArrayEventEventResponse> oneDimensionalArrayEventEventObservable(PlatonFilter filter) {
        return web3j.platonLogObservable(filter).map(new Func1<Log, OneDimensionalArrayEventEventResponse>() {
            @Override
            public OneDimensionalArrayEventEventResponse call(Log log) {
                Contract.EventValuesWithLog eventValues = extractEventParametersWithLog(ONEDIMENSIONALARRAYEVENT_EVENT, log);
                OneDimensionalArrayEventEventResponse typedResponse = new OneDimensionalArrayEventEventResponse();
                typedResponse.log = log;
                typedResponse.array = (List<BigInteger>) eventValues.getNonIndexedValues().get(0).getValue();
                return typedResponse;
            }
        });
    }

    public Observable<OneDimensionalArrayEventEventResponse> oneDimensionalArrayEventEventObservable(DefaultBlockParameter startBlock, DefaultBlockParameter endBlock) {
        PlatonFilter filter = new PlatonFilter(startBlock, endBlock, getContractAddress());
        filter.addSingleTopic(EventEncoder.encode(ONEDIMENSIONALARRAYEVENT_EVENT));
        return oneDimensionalArrayEventEventObservable(filter);
    }

    public List<StringEventEventResponse> getStringEventEvents(TransactionReceipt transactionReceipt) {
        List<Contract.EventValuesWithLog> valueList = extractEventParametersWithLog(STRINGEVENT_EVENT, transactionReceipt);
        ArrayList<StringEventEventResponse> responses = new ArrayList<StringEventEventResponse>(valueList.size());
        for (Contract.EventValuesWithLog eventValues : valueList) {
            StringEventEventResponse typedResponse = new StringEventEventResponse();
            typedResponse.log = eventValues.getLog();
            typedResponse.str = (String) eventValues.getNonIndexedValues().get(0).getValue();
            responses.add(typedResponse);
        }
        return responses;
    }

    public Observable<StringEventEventResponse> stringEventEventObservable(PlatonFilter filter) {
        return web3j.platonLogObservable(filter).map(new Func1<Log, StringEventEventResponse>() {
            @Override
            public StringEventEventResponse call(Log log) {
                Contract.EventValuesWithLog eventValues = extractEventParametersWithLog(STRINGEVENT_EVENT, log);
                StringEventEventResponse typedResponse = new StringEventEventResponse();
                typedResponse.log = log;
                typedResponse.str = (String) eventValues.getNonIndexedValues().get(0).getValue();
                return typedResponse;
            }
        });
    }

    public Observable<StringEventEventResponse> stringEventEventObservable(DefaultBlockParameter startBlock, DefaultBlockParameter endBlock) {
        PlatonFilter filter = new PlatonFilter(startBlock, endBlock, getContractAddress());
        filter.addSingleTopic(EventEncoder.encode(STRINGEVENT_EVENT));
        return stringEventEventObservable(filter);
    }

    public List<TwoDimensionalArrayEventEventResponse> getTwoDimensionalArrayEventEvents(TransactionReceipt transactionReceipt) {
        List<Contract.EventValuesWithLog> valueList = extractEventParametersWithLog(TWODIMENSIONALARRAYEVENT_EVENT, transactionReceipt);
        ArrayList<TwoDimensionalArrayEventEventResponse> responses = new ArrayList<TwoDimensionalArrayEventEventResponse>(valueList.size());
        for (Contract.EventValuesWithLog eventValues : valueList) {
            TwoDimensionalArrayEventEventResponse typedResponse = new TwoDimensionalArrayEventEventResponse();
            typedResponse.log = eventValues.getLog();
            typedResponse.array = (List<List<BigInteger>>) eventValues.getNonIndexedValues().get(0).getValue();
            responses.add(typedResponse);
        }
        return responses;
    }

    public Observable<TwoDimensionalArrayEventEventResponse> twoDimensionalArrayEventEventObservable(PlatonFilter filter) {
        return web3j.platonLogObservable(filter).map(new Func1<Log, TwoDimensionalArrayEventEventResponse>() {
            @Override
            public TwoDimensionalArrayEventEventResponse call(Log log) {
                Contract.EventValuesWithLog eventValues = extractEventParametersWithLog(TWODIMENSIONALARRAYEVENT_EVENT, log);
                TwoDimensionalArrayEventEventResponse typedResponse = new TwoDimensionalArrayEventEventResponse();
                typedResponse.log = log;
                typedResponse.array = (List<List<BigInteger>>) eventValues.getNonIndexedValues().get(0).getValue();
                return typedResponse;
            }
        });
    }

    public Observable<TwoDimensionalArrayEventEventResponse> twoDimensionalArrayEventEventObservable(DefaultBlockParameter startBlock, DefaultBlockParameter endBlock) {
        PlatonFilter filter = new PlatonFilter(startBlock, endBlock, getContractAddress());
        filter.addSingleTopic(EventEncoder.encode(TWODIMENSIONALARRAYEVENT_EVENT));
        return twoDimensionalArrayEventEventObservable(filter);
    }

    public RemoteCall<TransactionReceipt> testEnum() {
        final Function function = new Function(
                FUNC_TESTENUM, 
                Arrays.<Type>asList(), 
                Collections.<TypeReference<?>>emptyList());
        return executeRemoteCallTransaction(function);
    }

    public RemoteCall<TransactionReceipt> testOneDimensionalArray() {
        final Function function = new Function(
                FUNC_TESTONEDIMENSIONALARRAY, 
                Arrays.<Type>asList(), 
                Collections.<TypeReference<?>>emptyList());
        return executeRemoteCallTransaction(function);
    }

    public RemoteCall<TransactionReceipt> testStr() {
        final Function function = new Function(
                FUNC_TESTSTR, 
                Arrays.<Type>asList(), 
                Collections.<TypeReference<?>>emptyList());
        return executeRemoteCallTransaction(function);
    }

    public RemoteCall<TransactionReceipt> testTwoDimensionalArray() {
        final Function function = new Function(
                FUNC_TESTTWODIMENSIONALARRAY, 
                Arrays.<Type>asList(), 
                Collections.<TypeReference<?>>emptyList());
        return executeRemoteCallTransaction(function);
    }

    public static RemoteCall<EventTypeContract> deploy(Web3j web3j, Credentials credentials, GasProvider contractGasProvider, Long chainId) {
        return deployRemoteCall(EventTypeContract.class, web3j, credentials, contractGasProvider, BINARY,  "", chainId);
    }

    public static RemoteCall<EventTypeContract> deploy(Web3j web3j, TransactionManager transactionManager, GasProvider contractGasProvider, Long chainId) {
        return deployRemoteCall(EventTypeContract.class, web3j, transactionManager, contractGasProvider, BINARY,  "", chainId);
    }

    public static EventTypeContract load(String contractAddress, Web3j web3j, Credentials credentials, GasProvider contractGasProvider, Long chainId) {
        return new EventTypeContract(contractAddress, web3j, credentials, contractGasProvider, chainId);
    }

    public static EventTypeContract load(String contractAddress, Web3j web3j, TransactionManager transactionManager, GasProvider contractGasProvider, Long chainId) {
        return new EventTypeContract(contractAddress, web3j, transactionManager, contractGasProvider, chainId);
    }

    public static class EnumEventEventResponse {
        public Log log;

        public BigInteger choices;
    }

    public static class OneDimensionalArrayEventEventResponse {
        public Log log;

        public List<BigInteger> array;
    }

    public static class StringEventEventResponse {
        public Log log;

        public String str;
    }

    public static class TwoDimensionalArrayEventEventResponse {
        public Log log;

        public List<List<BigInteger>> array;
    }
}
