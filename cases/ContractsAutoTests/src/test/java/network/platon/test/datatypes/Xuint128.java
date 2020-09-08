package network.platon.test.datatypes;

import com.platon.rlp.datatypes.Uint128;

import java.math.BigInteger;

/**
 * @author wuyang
 */
public class Xuint128 extends Uint128{

    public Xuint128(String stringValue){
        super(new BigInteger(stringValue));
    }

    public Uint128 add(BigInteger bigInteger){
        return Uint128.of(this.getValue().add(bigInteger));
    }

    public Uint128 add(Uint128 uint128){
        return Uint128.of(this.getValue().add(uint128.getValue()));
    }

    public static final Uint128 ZERO = Uint128.of(BigInteger.ZERO);

    public static final Uint128 ONE = Uint128.of(BigInteger.ONE);

    public static final Uint128 TEN = Uint128.of(BigInteger.TEN);

}
