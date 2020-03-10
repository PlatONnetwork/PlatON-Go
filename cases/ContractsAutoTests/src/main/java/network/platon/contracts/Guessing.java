package network.platon.contracts;

import java.math.BigInteger;
import java.util.ArrayList;
import java.util.Arrays;
import java.util.Collections;
import java.util.List;
import java.util.concurrent.Callable;
import org.web3j.abi.EventEncoder;
import org.web3j.abi.FunctionEncoder;
import org.web3j.abi.TypeReference;
import org.web3j.abi.datatypes.Address;
import org.web3j.abi.datatypes.Bool;
import org.web3j.abi.datatypes.DynamicArray;
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
    private static final String BINARY = "60806040526000600160006101000a81548160ff021916908315150217905550674563918244f40000600255600060085534801561003c57600080fd5b506040516112df3803806112df8339818101604052602081101561005f57600080fd5b810190808051906020019092919050505033600a60006101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff1602179055508060008190555050611217806100c86000396000f3fe6080604052600436106100f35760003560e01c80636e5ab6711161008a578063b69ef8a811610059578063b69ef8a81461056d578063c2930f9114610598578063c7f40b64146105c3578063ecbde5e61461061a576100f3565b80636e5ab6711461044d57806382d333a0146104b2578063a46c3637146104dd578063ae44794114610508576100f3565b806313eaca43116100c657806313eaca431461032d5780631ef9c56f1461035c578063220bc55e146103c8578063632dbe90146103d2576100f3565b8063045f9c9714610245578063062d6a98146102c0578063083c6323146102eb5780630eecae2114610316575b60025434101561010257600080fd5b60003490503360066000600854815260200190815260200160002060006101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff16021790555060086000815480929190600101919050555080600560003373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060008282540192505081905550806003600082825401925050819055507fe842aea7a5f1b01049d752008c53c52890b1a6daf660cf39e8eec506112bbdf633826001604051808473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200183815260200182151515158152602001935050505060405180910390a150005b34801561025157600080fd5b5061027e6004803603602081101561026857600080fd5b8101908080359060200190929190505050610645565b604051808273ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200191505060405180910390f35b3480156102cc57600080fd5b506102d5610681565b6040518082815260200191505060405180910390f35b3480156102f757600080fd5b50610300610687565b6040518082815260200191505060405180910390f35b34801561032257600080fd5b5061032b61068d565b005b34801561033957600080fd5b50610342610f23565b604051808215151515815260200191505060405180910390f35b34801561036857600080fd5b50610371610f36565b6040518080602001828103825283818151815260200191508051906020019060200280838360005b838110156103b4578082015181840152602081019050610399565b505050509050019250505060405180910390f35b6103d0610fc4565b005b3480156103de57600080fd5b5061040b600480360360208110156103f557600080fd5b810190808035906020019092919050505061111b565b604051808273ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200191505060405180910390f35b34801561045957600080fd5b5061049c6004803603602081101561047057600080fd5b81019080803573ffffffffffffffffffffffffffffffffffffffff16906020019092919050505061114e565b6040518082815260200191505060405180910390f35b3480156104be57600080fd5b506104c7611166565b6040518082815260200191505060405180910390f35b3480156104e957600080fd5b506104f261116c565b6040518082815260200191505060405180910390f35b34801561051457600080fd5b506105576004803603602081101561052b57600080fd5b81019080803573ffffffffffffffffffffffffffffffffffffffff169060200190929190505050611179565b6040518082815260200191505060405180910390f35b34801561057957600080fd5b50610582611191565b6040518082815260200191505060405180910390f35b3480156105a457600080fd5b506105ad611197565b6040518082815260200191505060405180910390f35b3480156105cf57600080fd5b506105d861119d565b604051808273ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200191505060405180910390f35b34801561062657600080fd5b5061062f6111c3565b6040518082815260200191505060405180910390f35b6009818154811061065257fe5b906000526020600020016000915054906101000a900473ffffffffffffffffffffffffffffffffffffffff1681565b60085481565b60005481565b600054431115610f2157600160009054906101000a900460ff1615801561070157503373ffffffffffffffffffffffffffffffffffffffff16600a60009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16145b801561070f57506000600854115b15610f20576000805440604051602001808281526020019150506040516020818303038152906040528051906020012060001c90506000600854828161075157fe5b06905060006064600854101561093357600a828161076b57fe5b06905060008114156108495760008090505b600854811015610843576000600a8383038161079557fe5b0614156108365760096006600083815260200190815260200160002060009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1690806001815401808255809150509060018203906000526020600020016000909192909190916101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff160217905550505b808060010191505061077d565b5061092e565b60008090505b60085481101561092c576000600a828161086557fe5b061415801561088057506000600a8383038161087d57fe5b06145b1561091f5760096006600083815260200190815260200160002060009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1690806001815401808255809150509060018203906000526020600020016000909192909190916101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff160217905550505b808060010191505061084f565b505b610ce4565b6127106008541015610b11576064828161094957fe5b0690506000811415610a275760008090505b600854811015610a2157600060648383038161097357fe5b061415610a145760096006600083815260200190815260200160002060009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1690806001815401808255809150509060018203906000526020600020016000909192909190916101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff160217905550505b808060010191505061095b565b50610b0c565b60008090505b600854811015610b0a57600060648281610a4357fe5b0614158015610a5e57506000606483830381610a5b57fe5b06145b15610afd5760096006600083815260200190815260200160002060009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1690806001815401808255809150509060018203906000526020600020016000909192909190916101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff160217905550505b8080600101915050610a2d565b505b610ce3565b6103e88281610b1c57fe5b0690506000811415610bfb5760008090505b600854811015610bf55760006103e883830381610b4757fe5b061415610be85760096006600083815260200190815260200160002060009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1690806001815401808255809150509060018203906000526020600020016000909192909190916101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff160217905550505b8080600101915050610b2e565b50610ce2565b60008090505b600854811015610ce05760006103e88281610c1857fe5b0614158015610c34575060006103e883830381610c3157fe5b06145b15610cd35760096006600083815260200190815260200160002060009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1690806001815401808255809150509060018203906000526020600020016000909192909190916101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff160217905550505b8080600101915050610c01565b505b5b5b60098054905060035481610cf457fe5b04600481905550600080600090505b600980549050811015610f005760016007600060098481548110610d2357fe5b9060005260206000200160009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002054016007600060098481548110610d9b57fe5b9060005260206000200160009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002081905550600073ffffffffffffffffffffffffffffffffffffffff1660098281548110610e2957fe5b9060005260206000200160009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1614610ef35760098181548110610e7c57fe5b9060005260206000200160009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1691508173ffffffffffffffffffffffffffffffffffffffff166108fc6004549081150290604051600060405180830381858888f19350505050158015610ef1573d6000803e3d6000fd5b505b8080600101915050610d03565b5060018060006101000a81548160ff021916908315150217905550505050505b5b565b600160009054906101000a900460ff1681565b60606009805480602002602001604051908101604052809291908181526020018280548015610fba57602002820191906000526020600020905b8160009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019060010190808311610f70575b5050505050905090565b43600054106111195760025434106111185760003490503360066000600854815260200190815260200160002060006101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff16021790555060086000815480929190600101919050555080600560003373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060008282540192505081905550806003600082825401925050819055507fe842aea7a5f1b01049d752008c53c52890b1a6daf660cf39e8eec506112bbdf633826001604051808473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200183815260200182151515158152602001935050505060405180910390a1505b5b565b60066020528060005260406000206000915054906101000a900473ffffffffffffffffffffffffffffffffffffffff1681565b60076020528060005260406000206000915090505481565b60045481565b6000600980549050905090565b60056020528060005260406000206000915090505481565b60035481565b60025481565b600a60009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1681565b60003073ffffffffffffffffffffffffffffffffffffffff163190509056fea265627a7a72315820ef9ef03c9593b0cdbe23e3bde279b357ddb28d57f57938546c6777d65de9321a64736f6c634300050d0032";

    public static final String FUNC_INDEXOFGUSSINGER = "IndexOfgussinger";

    public static final String FUNC_AVERAGEAMOUNT = "averageAmount";

    public static final String FUNC_BALANCE = "balance";

    public static final String FUNC_BASEUNIT = "baseUnit";

    public static final String FUNC_CREATEADDRESS = "createAddress";

    public static final String FUNC_DRAW = "draw";

    public static final String FUNC_ENDBLOCK = "endBlock";

    public static final String FUNC_GETBALANCEOF = "getBalanceOf";

    public static final String FUNC_GETWINNERADDRESSES = "getWinnerAddresses";

    public static final String FUNC_GETWINNERCOUNT = "getWinnerCount";

    public static final String FUNC_GUESSINGCLOSED = "guessingClosed";

    public static final String FUNC_GUESSINGWITHLAT = "guessingWithLat";

    public static final String FUNC_GUSSINGERLAT = "gussingerLat";

    public static final String FUNC_INDEXKEY = "indexKey";

    public static final String FUNC_WINNERADDRESSES = "winnerAddresses";

    public static final String FUNC_WINNERMAP = "winnerMap";

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

    public static RemoteCall<Guessing> deploy(Web3j web3j, Credentials credentials, GasProvider contractGasProvider, BigInteger _endBlock) {
        String encodedConstructor = FunctionEncoder.encodeConstructor(Arrays.<Type>asList(new org.web3j.abi.datatypes.generated.Uint256(_endBlock)));
        return deployRemoteCall(Guessing.class, web3j, credentials, contractGasProvider, BINARY, encodedConstructor);
    }

    public static RemoteCall<Guessing> deploy(Web3j web3j, TransactionManager transactionManager, GasProvider contractGasProvider, BigInteger _endBlock) {
        String encodedConstructor = FunctionEncoder.encodeConstructor(Arrays.<Type>asList(new org.web3j.abi.datatypes.generated.Uint256(_endBlock)));
        return deployRemoteCall(Guessing.class, web3j, transactionManager, contractGasProvider, BINARY, encodedConstructor);
    }

    @Deprecated
    public static RemoteCall<Guessing> deploy(Web3j web3j, Credentials credentials, BigInteger gasPrice, BigInteger gasLimit, BigInteger _endBlock) {
        String encodedConstructor = FunctionEncoder.encodeConstructor(Arrays.<Type>asList(new org.web3j.abi.datatypes.generated.Uint256(_endBlock)));
        return deployRemoteCall(Guessing.class, web3j, credentials, gasPrice, gasLimit, BINARY, encodedConstructor);
    }

    @Deprecated
    public static RemoteCall<Guessing> deploy(Web3j web3j, TransactionManager transactionManager, BigInteger gasPrice, BigInteger gasLimit, BigInteger _endBlock) {
        String encodedConstructor = FunctionEncoder.encodeConstructor(Arrays.<Type>asList(new org.web3j.abi.datatypes.generated.Uint256(_endBlock)));
        return deployRemoteCall(Guessing.class, web3j, transactionManager, gasPrice, gasLimit, BINARY, encodedConstructor);
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

    public RemoteCall<BigInteger> averageAmount() {
        final Function function = new Function(FUNC_AVERAGEAMOUNT, 
                Arrays.<Type>asList(), 
                Arrays.<TypeReference<?>>asList(new TypeReference<Uint256>() {}));
        return executeRemoteCallSingleValueReturn(function, BigInteger.class);
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

    public RemoteCall<List> getWinnerAddresses() {
        final Function function = new Function(FUNC_GETWINNERADDRESSES, 
                Arrays.<Type>asList(), 
                Arrays.<TypeReference<?>>asList(new TypeReference<DynamicArray<Address>>() {}));
        return new RemoteCall<List>(
                new Callable<List>() {
                    @Override
                    @SuppressWarnings("unchecked")
                    public List call() throws Exception {
                        List<Type> result = (List<Type>) executeCallSingleValueReturn(function, List.class);
                        return convertToNative(result);
                    }
                });
    }

    public RemoteCall<BigInteger> getWinnerCount() {
        final Function function = new Function(FUNC_GETWINNERCOUNT, 
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

    public RemoteCall<String> winnerAddresses(BigInteger param0) {
        final Function function = new Function(FUNC_WINNERADDRESSES, 
                Arrays.<Type>asList(new org.web3j.abi.datatypes.generated.Uint256(param0)), 
                Arrays.<TypeReference<?>>asList(new TypeReference<Address>() {}));
        return executeRemoteCallSingleValueReturn(function, String.class);
    }

    public RemoteCall<BigInteger> winnerMap(String param0) {
        final Function function = new Function(FUNC_WINNERMAP, 
                Arrays.<Type>asList(new org.web3j.abi.datatypes.Address(param0)), 
                Arrays.<TypeReference<?>>asList(new TypeReference<Uint256>() {}));
        return executeRemoteCallSingleValueReturn(function, BigInteger.class);
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
