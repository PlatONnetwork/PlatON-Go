package network.platon.contracts;

import java.math.BigInteger;
import java.util.ArrayList;
import java.util.Arrays;
import java.util.Collections;
import java.util.List;
import org.web3j.abi.EventEncoder;
import org.web3j.abi.FunctionEncoder;
import org.web3j.abi.TypeReference;
import org.web3j.abi.datatypes.Address;
import org.web3j.abi.datatypes.Bool;
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
public class Guessing extends Contract {
    private static final String BINARY = "60806040526000600160006101000a81548160ff021916908315150217905550674563918244f400006002556000600655604051610aad380380610aad8339818101604052602081101561005257600080fd5b810190808051906020019092919050505033600860006101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff16021790555080600081905550506109f2806100bb6000396000f3fe6080604052600436106100a75760003560e01c8063ae44794111610064578063ae44794114610347578063b69ef8a8146103ac578063c2930f91146103d7578063c7f40b6414610402578063ecbde5e614610459578063ed7a4e0b14610484576100a7565b8063062d6a9814610226578063083c6323146102515780630eecae211461027c57806313eaca4314610293578063220bc55e146102c2578063632dbe90146102cc575b600034116100b457600080fd5b6100c0600354346104db565b6003819055506000349050600060025482816100d857fe5b04905060008090505b8181101561015c573360056000600654815260200190815260200160002060006101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff16021790555060066000815480929190600101919050555080806001019150506100e1565b5081600460003373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff168152602001908152602001600020600082825401925050819055507fe842aea7a5f1b01049d752008c53c52890b1a6daf660cf39e8eec506112bbdf633836001604051808473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200183815260200182151515158152602001935050505060405180910390a15050005b34801561023257600080fd5b5061023b6104f7565b6040518082815260200191505060405180910390f35b34801561025d57600080fd5b506102666104fd565b6040518082815260200191505060405180910390f35b34801561028857600080fd5b50610291610503565b005b34801561029f57600080fd5b506102a8610757565b604051808215151515815260200191505060405180910390f35b6102ca61076a565b005b3480156102d857600080fd5b50610305600480360360208110156102ef57600080fd5b81019080803590602001909291905050506108fb565b604051808273ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200191505060405180910390f35b34801561035357600080fd5b506103966004803603602081101561036a57600080fd5b81019080803573ffffffffffffffffffffffffffffffffffffffff16906020019092919050505061092e565b6040518082815260200191505060405180910390f35b3480156103b857600080fd5b506103c1610946565b6040518082815260200191505060405180910390f35b3480156103e357600080fd5b506103ec61094c565b6040518082815260200191505060405180910390f35b34801561040e57600080fd5b50610417610952565b604051808273ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200191505060405180910390f35b34801561046557600080fd5b5061046e610978565b6040518082815260200191505060405180910390f35b34801561049057600080fd5b50610499610997565b604051808273ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200191505060405180910390f35b6000808284019050838110156104ed57fe5b8091505092915050565b60065481565b60005481565b60005443111561075557600160009054906101000a900460ff1615801561057757503373ffffffffffffffffffffffffffffffffffffffff16600860009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16145b15610754576000805440604051602001808281526020019150506040516020818303038152906040528051906020012060001c9050600060065482816105b957fe5b0690506005600082815260200190815260200160002060009054906101000a900473ffffffffffffffffffffffffffffffffffffffff16600760006101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff160217905550600760009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff166108fc6003549081150290604051600060405180830381858888f1935050505015801561069a573d6000803e3d6000fd5b507fe842aea7a5f1b01049d752008c53c52890b1a6daf660cf39e8eec506112bbdf6600760009054906101000a900473ffffffffffffffffffffffffffffffffffffffff166003546000604051808473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200183815260200182151515158152602001935050505060405180910390a160018060006101000a81548160ff02191690831515021790555050505b5b565b600160009054906101000a900460ff1681565b43600054106108f95760025434106108f8576000341161078957600080fd5b60003490506000600254828161079b57fe5b04905060008090505b8181101561081f573360056000600654815260200190815260200160002060006101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff16021790555060066000815480929190600101919050555080806001019150506107a4565b5081600460003373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060008282540192505081905550816003600082825401925050819055507fe842aea7a5f1b01049d752008c53c52890b1a6daf660cf39e8eec506112bbdf633836001604051808473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200183815260200182151515158152602001935050505060405180910390a150505b5b565b60056020528060005260406000206000915054906101000a900473ffffffffffffffffffffffffffffffffffffffff1681565b60046020528060005260406000206000915090505481565b60035481565b60025481565b600860009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1681565b60003073ffffffffffffffffffffffffffffffffffffffff1631905090565b600760009054906101000a900473ffffffffffffffffffffffffffffffffffffffff168156fea265627a7a723158206b4d347784f71d6715e299138d5a58a7c45aa937a4a5367ba1eed0eaf54f387864736f6c634300050d0032";

    public static final String FUNC_INDEXOFGUSSINGER = "IndexOfgussinger";

    public static final String FUNC_BALANCE = "balance";

    public static final String FUNC_BASEUNIT = "baseUnit";

    public static final String FUNC_CREATEADDRESS = "createAddress";

    public static final String FUNC_DRAW = "draw";

    public static final String FUNC_ENDBLOCK = "endBlock";

    public static final String FUNC_GETBALANCEOF = "getBalanceOf";

    public static final String FUNC_GUESSINGCLOSED = "guessingClosed";

    public static final String FUNC_GUESSINGWITHLAT = "guessingWithLat";

    public static final String FUNC_GUSSINGERLAT = "gussingerLat";

    public static final String FUNC_INDEXKEY = "indexKey";

    public static final String FUNC_WINNERADDRESS = "winnerAddress";

    public static final Event CURRENTBALANCE_EVENT = new Event("CurrentBalance", 
            Arrays.<TypeReference<?>>asList(new TypeReference<Address>() {}, new TypeReference<Uint256>() {}));
    ;

    public static final Event FUNDTRANSFER_EVENT = new Event("FundTransfer", 
            Arrays.<TypeReference<?>>asList(new TypeReference<Address>() {}, new TypeReference<Uint256>() {}, new TypeReference<Bool>() {}));
    ;

    @Deprecated
    protected Guessing(String contractAddress, Web3j web3j, Credentials credentials, BigInteger gasPrice, BigInteger gasLimit) {
        super(BINARY, contractAddress, web3j, credentials, gasPrice, gasLimit);
    }

    protected Guessing(String contractAddress, Web3j web3j, Credentials credentials, GasProvider contractGasProvider) {
        super(BINARY, contractAddress, web3j, credentials, contractGasProvider);
    }

    @Deprecated
    protected Guessing(String contractAddress, Web3j web3j, TransactionManager transactionManager, BigInteger gasPrice, BigInteger gasLimit) {
        super(BINARY, contractAddress, web3j, transactionManager, gasPrice, gasLimit);
    }

    protected Guessing(String contractAddress, Web3j web3j, TransactionManager transactionManager, GasProvider contractGasProvider) {
        super(BINARY, contractAddress, web3j, transactionManager, contractGasProvider);
    }

    public static RemoteCall<Guessing> deploy(Web3j web3j, Credentials credentials, GasProvider contractGasProvider, BigInteger initialWeiValue, BigInteger _endBlock) {
        String encodedConstructor = FunctionEncoder.encodeConstructor(Arrays.<Type>asList(new org.web3j.abi.datatypes.generated.Uint256(_endBlock)));
        return deployRemoteCall(Guessing.class, web3j, credentials, contractGasProvider, BINARY, encodedConstructor, initialWeiValue);
    }

    public static RemoteCall<Guessing> deploy(Web3j web3j, TransactionManager transactionManager, GasProvider contractGasProvider, BigInteger initialWeiValue, BigInteger _endBlock) {
        String encodedConstructor = FunctionEncoder.encodeConstructor(Arrays.<Type>asList(new org.web3j.abi.datatypes.generated.Uint256(_endBlock)));
        return deployRemoteCall(Guessing.class, web3j, transactionManager, contractGasProvider, BINARY, encodedConstructor, initialWeiValue);
    }

    @Deprecated
    public static RemoteCall<Guessing> deploy(Web3j web3j, Credentials credentials, BigInteger gasPrice, BigInteger gasLimit, BigInteger initialWeiValue, BigInteger _endBlock) {
        String encodedConstructor = FunctionEncoder.encodeConstructor(Arrays.<Type>asList(new org.web3j.abi.datatypes.generated.Uint256(_endBlock)));
        return deployRemoteCall(Guessing.class, web3j, credentials, gasPrice, gasLimit, BINARY, encodedConstructor, initialWeiValue);
    }

    @Deprecated
    public static RemoteCall<Guessing> deploy(Web3j web3j, TransactionManager transactionManager, BigInteger gasPrice, BigInteger gasLimit, BigInteger initialWeiValue, BigInteger _endBlock) {
        String encodedConstructor = FunctionEncoder.encodeConstructor(Arrays.<Type>asList(new org.web3j.abi.datatypes.generated.Uint256(_endBlock)));
        return deployRemoteCall(Guessing.class, web3j, transactionManager, gasPrice, gasLimit, BINARY, encodedConstructor, initialWeiValue);
    }

    public List<CurrentBalanceEventResponse> getCurrentBalanceEvents(TransactionReceipt transactionReceipt) {
        List<Contract.EventValuesWithLog> valueList = extractEventParametersWithLog(CURRENTBALANCE_EVENT, transactionReceipt);
        ArrayList<CurrentBalanceEventResponse> responses = new ArrayList<CurrentBalanceEventResponse>(valueList.size());
        for (Contract.EventValuesWithLog eventValues : valueList) {
            CurrentBalanceEventResponse typedResponse = new CurrentBalanceEventResponse();
            typedResponse.log = eventValues.getLog();
            typedResponse._msgSenderAddress = (String) eventValues.getNonIndexedValues().get(0).getValue();
            typedResponse._balance = (BigInteger) eventValues.getNonIndexedValues().get(1).getValue();
            responses.add(typedResponse);
        }
        return responses;
    }

    public Observable<CurrentBalanceEventResponse> currentBalanceEventObservable(PlatonFilter filter) {
        return web3j.platonLogObservable(filter).map(new Func1<Log, CurrentBalanceEventResponse>() {
            @Override
            public CurrentBalanceEventResponse call(Log log) {
                Contract.EventValuesWithLog eventValues = extractEventParametersWithLog(CURRENTBALANCE_EVENT, log);
                CurrentBalanceEventResponse typedResponse = new CurrentBalanceEventResponse();
                typedResponse.log = log;
                typedResponse._msgSenderAddress = (String) eventValues.getNonIndexedValues().get(0).getValue();
                typedResponse._balance = (BigInteger) eventValues.getNonIndexedValues().get(1).getValue();
                return typedResponse;
            }
        });
    }

    public Observable<CurrentBalanceEventResponse> currentBalanceEventObservable(DefaultBlockParameter startBlock, DefaultBlockParameter endBlock) {
        PlatonFilter filter = new PlatonFilter(startBlock, endBlock, getContractAddress());
        filter.addSingleTopic(EventEncoder.encode(CURRENTBALANCE_EVENT));
        return currentBalanceEventObservable(filter);
    }

    public List<FundTransferEventResponse> getFundTransferEvents(TransactionReceipt transactionReceipt) {
        List<Contract.EventValuesWithLog> valueList = extractEventParametersWithLog(FUNDTRANSFER_EVENT, transactionReceipt);
        ArrayList<FundTransferEventResponse> responses = new ArrayList<FundTransferEventResponse>(valueList.size());
        for (Contract.EventValuesWithLog eventValues : valueList) {
            FundTransferEventResponse typedResponse = new FundTransferEventResponse();
            typedResponse.log = eventValues.getLog();
            typedResponse._backer = (String) eventValues.getNonIndexedValues().get(0).getValue();
            typedResponse._amount = (BigInteger) eventValues.getNonIndexedValues().get(1).getValue();
            typedResponse._isSuccess = (Boolean) eventValues.getNonIndexedValues().get(2).getValue();
            responses.add(typedResponse);
        }
        return responses;
    }

    public Observable<FundTransferEventResponse> fundTransferEventObservable(PlatonFilter filter) {
        return web3j.platonLogObservable(filter).map(new Func1<Log, FundTransferEventResponse>() {
            @Override
            public FundTransferEventResponse call(Log log) {
                Contract.EventValuesWithLog eventValues = extractEventParametersWithLog(FUNDTRANSFER_EVENT, log);
                FundTransferEventResponse typedResponse = new FundTransferEventResponse();
                typedResponse.log = log;
                typedResponse._backer = (String) eventValues.getNonIndexedValues().get(0).getValue();
                typedResponse._amount = (BigInteger) eventValues.getNonIndexedValues().get(1).getValue();
                typedResponse._isSuccess = (Boolean) eventValues.getNonIndexedValues().get(2).getValue();
                return typedResponse;
            }
        });
    }

    public Observable<FundTransferEventResponse> fundTransferEventObservable(DefaultBlockParameter startBlock, DefaultBlockParameter endBlock) {
        PlatonFilter filter = new PlatonFilter(startBlock, endBlock, getContractAddress());
        filter.addSingleTopic(EventEncoder.encode(FUNDTRANSFER_EVENT));
        return fundTransferEventObservable(filter);
    }

    public RemoteCall<String> IndexOfgussinger(BigInteger param0) {
        final Function function = new Function(FUNC_INDEXOFGUSSINGER, 
                Arrays.<Type>asList(new org.web3j.abi.datatypes.generated.Uint256(param0)), 
                Arrays.<TypeReference<?>>asList(new TypeReference<Address>() {}));
        return executeRemoteCallSingleValueReturn(function, String.class);
    }

    public RemoteCall<BigInteger> balance() {
        final Function function = new Function(FUNC_BALANCE, 
                Arrays.<Type>asList(), 
                Arrays.<TypeReference<?>>asList(new TypeReference<Uint256>() {}));
        return executeRemoteCallSingleValueReturn(function, BigInteger.class);
    }

    public RemoteCall<BigInteger> baseUnit() {
        final Function function = new Function(FUNC_BASEUNIT, 
                Arrays.<Type>asList(), 
                Arrays.<TypeReference<?>>asList(new TypeReference<Uint256>() {}));
        return executeRemoteCallSingleValueReturn(function, BigInteger.class);
    }

    public RemoteCall<String> createAddress() {
        final Function function = new Function(FUNC_CREATEADDRESS, 
                Arrays.<Type>asList(), 
                Arrays.<TypeReference<?>>asList(new TypeReference<Address>() {}));
        return executeRemoteCallSingleValueReturn(function, String.class);
    }

    public RemoteCall<TransactionReceipt> draw() {
        final Function function = new Function(
                FUNC_DRAW, 
                Arrays.<Type>asList(), 
                Collections.<TypeReference<?>>emptyList());
        return executeRemoteCallTransaction(function);
    }

    public RemoteCall<BigInteger> endBlock() {
        final Function function = new Function(FUNC_ENDBLOCK, 
                Arrays.<Type>asList(), 
                Arrays.<TypeReference<?>>asList(new TypeReference<Uint256>() {}));
        return executeRemoteCallSingleValueReturn(function, BigInteger.class);
    }

    public RemoteCall<BigInteger> getBalanceOf() {
        final Function function = new Function(FUNC_GETBALANCEOF, 
                Arrays.<Type>asList(), 
                Arrays.<TypeReference<?>>asList(new TypeReference<Uint256>() {}));
        return executeRemoteCallSingleValueReturn(function, BigInteger.class);
    }

    public RemoteCall<Boolean> guessingClosed() {
        final Function function = new Function(FUNC_GUESSINGCLOSED, 
                Arrays.<Type>asList(), 
                Arrays.<TypeReference<?>>asList(new TypeReference<Bool>() {}));
        return executeRemoteCallSingleValueReturn(function, Boolean.class);
    }

    public RemoteCall<TransactionReceipt> guessingWithLat(BigInteger weiValue) {
        final Function function = new Function(
                FUNC_GUESSINGWITHLAT, 
                Arrays.<Type>asList(), 
                Collections.<TypeReference<?>>emptyList());
        return executeRemoteCallTransaction(function, weiValue);
    }

    public RemoteCall<BigInteger> gussingerLat(String param0) {
        final Function function = new Function(FUNC_GUSSINGERLAT, 
                Arrays.<Type>asList(new org.web3j.abi.datatypes.Address(param0)), 
                Arrays.<TypeReference<?>>asList(new TypeReference<Uint256>() {}));
        return executeRemoteCallSingleValueReturn(function, BigInteger.class);
    }

    public RemoteCall<BigInteger> indexKey() {
        final Function function = new Function(FUNC_INDEXKEY, 
                Arrays.<Type>asList(), 
                Arrays.<TypeReference<?>>asList(new TypeReference<Uint256>() {}));
        return executeRemoteCallSingleValueReturn(function, BigInteger.class);
    }

    public RemoteCall<String> winnerAddress() {
        final Function function = new Function(FUNC_WINNERADDRESS, 
                Arrays.<Type>asList(), 
                Arrays.<TypeReference<?>>asList(new TypeReference<Address>() {}));
        return executeRemoteCallSingleValueReturn(function, String.class);
    }

    @Deprecated
    public static Guessing load(String contractAddress, Web3j web3j, Credentials credentials, BigInteger gasPrice, BigInteger gasLimit) {
        return new Guessing(contractAddress, web3j, credentials, gasPrice, gasLimit);
    }

    @Deprecated
    public static Guessing load(String contractAddress, Web3j web3j, TransactionManager transactionManager, BigInteger gasPrice, BigInteger gasLimit) {
        return new Guessing(contractAddress, web3j, transactionManager, gasPrice, gasLimit);
    }

    public static Guessing load(String contractAddress, Web3j web3j, Credentials credentials, GasProvider contractGasProvider) {
        return new Guessing(contractAddress, web3j, credentials, contractGasProvider);
    }

    public static Guessing load(String contractAddress, Web3j web3j, TransactionManager transactionManager, GasProvider contractGasProvider) {
        return new Guessing(contractAddress, web3j, transactionManager, contractGasProvider);
    }

    public static class CurrentBalanceEventResponse {
        public Log log;

        public String _msgSenderAddress;

        public BigInteger _balance;
    }

    public static class FundTransferEventResponse {
        public Log log;

        public String _backer;

        public BigInteger _amount;

        public Boolean _isSuccess;
    }
}
