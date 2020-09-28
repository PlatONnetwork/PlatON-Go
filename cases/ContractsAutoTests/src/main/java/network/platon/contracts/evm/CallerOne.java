package network.platon.contracts.evm;

import com.alaya.abi.solidity.TypeReference;
import com.alaya.abi.solidity.datatypes.Function;
import com.alaya.abi.solidity.datatypes.Type;
import com.alaya.abi.solidity.datatypes.generated.Uint256;
import com.alaya.crypto.Credentials;
import com.alaya.protocol.Web3j;
import com.alaya.protocol.core.RemoteCall;
import com.alaya.protocol.core.methods.response.TransactionReceipt;
import com.alaya.tx.Contract;
import com.alaya.tx.TransactionManager;
import com.alaya.tx.gas.GasProvider;
import java.math.BigInteger;
import java.util.Arrays;
import java.util.Collections;

/**
 * <p>Auto generated code.
 * <p><strong>Do not modify!</strong>
 * <p>Please use the <a href="https://github.com/PlatONnetwork/client-sdk-java/releases">platon-web3j command line tools</a>,
 * or the com.alaya.codegen.SolidityFunctionWrapperGenerator in the 
 * <a href="https://github.com/PlatONnetwork/client-sdk-java/tree/master/codegen">codegen module</a> to update.
 *
 * <p>Generated with web3j version 0.13.2.0.
 */
public class CallerOne extends Contract {
    private static final String BINARY = "608060405234801561001057600080fd5b50610625806100206000396000f3fe608060405234801561001057600080fd5b50600436106100415760003560e01c80630c55699c146100465780637b8ed01814610064578063f88bef8c14610082575b600080fd5b61004e61008c565b6040518082815260200191505060405180910390f35b61006c610092565b6040518082815260200191505060405180910390f35b61008a61009b565b005b60005481565b60008054905090565b60006040516100a9906101d2565b604051809103906000f0801580156100c5573d6000803e3d6000fd5b5090508073ffffffffffffffffffffffffffffffffffffffff1660405180807f696e63282900000000000000000000000000000000000000000000000000000081525060050190506040518091039020604051602001808281526020019150506040516020818303038152906040526040518082805190602001908083835b602083106101675780518252602082019150602081019050602083039250610144565b6001836020036101000a038019825116818451168082178552505050505050905001915050600060405180830381855af49150503d80600081146101c7576040519150601f19603f3d011682016040523d82523d6000602084013e6101cc565b606091505b50505050565b610411806101e08339019056fe608060405234801561001057600080fd5b506103f1806100206000396000f3fe608060405234801561001057600080fd5b50600436106100415760003560e01c80630c55699c14610046578063371303c0146100645780635a3617561461006e575b600080fd5b61004e61008c565b6040518082815260200191505060405180910390f35b61006c610092565b005b610076610236565b6040518082815260200191505060405180910390f35b60005481565b60006040516100a09061023f565b604051809103906000f0801580156100bc573d6000803e3d6000fd5b5090508073ffffffffffffffffffffffffffffffffffffffff1660405180807f696e63282900000000000000000000000000000000000000000000000000000081525060050190506040518091039020604051602001808281526020019150506040516020818303038152906040526040518082805190602001908083835b6020831061015e578051825260208201915060208101905060208303925061013b565b6001836020036101000a038019825116818451168082178552505050505050905001915050600060405180830381855af49150503d80600081146101be576040519150601f19603f3d011682016040523d82523d6000602084013e6101c3565b606091505b5050507fb0333e0e3a6b99318e4e2e0d7e5e5f93646f9cbf62da1587955a4092bf7df6e733600054604051808373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020018281526020019250505060405180910390a150565b60008054905090565b6101708061024d8339019056fe608060405234801561001057600080fd5b50610150806100206000396000f3fe608060405234801561001057600080fd5b50600436106100415760003560e01c80630c55699c1461004657806317f936fb14610064578063371303c014610082575b600080fd5b61004e61008c565b6040518082815260200191505060405180910390f35b61006c610092565b6040518082815260200191505060405180910390f35b61008a61009b565b005b60005481565b60008054905090565b60008081548092919060010191905055507fb0333e0e3a6b99318e4e2e0d7e5e5f93646f9cbf62da1587955a4092bf7df6e733600054604051808373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020018281526020019250505060405180910390a156fea265627a7a723158201ac0d5b5fa0956df35d2c00540b867937ccc882c2d650ce73607c0c57031e91564736f6c63430005110032a265627a7a7231582074807c8df2ca81f4bba5a7b7db2fd1f1196c38a13bd07bd61b9fcc941365017564736f6c63430005110032a265627a7a72315820907d5902e2980e99a82f4d3a790aba43914e8a7d79324537322cad4eb81b5a6464736f6c63430005110032";

    public static final String FUNC_GETCALLERX = "getCallerX";

    public static final String FUNC_INC_DELEGATECALL = "inc_delegatecall";

    public static final String FUNC_X = "x";

    protected CallerOne(String contractAddress, Web3j web3j, Credentials credentials, GasProvider contractGasProvider, Long chainId) {
        super(BINARY, contractAddress, web3j, credentials, contractGasProvider, chainId);
    }

    protected CallerOne(String contractAddress, Web3j web3j, TransactionManager transactionManager, GasProvider contractGasProvider, Long chainId) {
        super(BINARY, contractAddress, web3j, transactionManager, contractGasProvider, chainId);
    }

    public RemoteCall<BigInteger> getCallerX() {
        final Function function = new Function(FUNC_GETCALLERX, 
                Arrays.<Type>asList(), 
                Arrays.<TypeReference<?>>asList(new TypeReference<Uint256>() {}));
        return executeRemoteCallSingleValueReturn(function, BigInteger.class);
    }

    public RemoteCall<TransactionReceipt> inc_delegatecall() {
        final Function function = new Function(
                FUNC_INC_DELEGATECALL, 
                Arrays.<Type>asList(), 
                Collections.<TypeReference<?>>emptyList());
        return executeRemoteCallTransaction(function);
    }

    public RemoteCall<BigInteger> x() {
        final Function function = new Function(FUNC_X, 
                Arrays.<Type>asList(), 
                Arrays.<TypeReference<?>>asList(new TypeReference<Uint256>() {}));
        return executeRemoteCallSingleValueReturn(function, BigInteger.class);
    }

    public static RemoteCall<CallerOne> deploy(Web3j web3j, Credentials credentials, GasProvider contractGasProvider, Long chainId) {
        return deployRemoteCall(CallerOne.class, web3j, credentials, contractGasProvider, BINARY,  "", chainId);
    }

    public static RemoteCall<CallerOne> deploy(Web3j web3j, TransactionManager transactionManager, GasProvider contractGasProvider, Long chainId) {
        return deployRemoteCall(CallerOne.class, web3j, transactionManager, contractGasProvider, BINARY,  "", chainId);
    }

    public static CallerOne load(String contractAddress, Web3j web3j, Credentials credentials, GasProvider contractGasProvider, Long chainId) {
        return new CallerOne(contractAddress, web3j, credentials, contractGasProvider, chainId);
    }

    public static CallerOne load(String contractAddress, Web3j web3j, TransactionManager transactionManager, GasProvider contractGasProvider, Long chainId) {
        return new CallerOne(contractAddress, web3j, transactionManager, contractGasProvider, chainId);
    }
}
