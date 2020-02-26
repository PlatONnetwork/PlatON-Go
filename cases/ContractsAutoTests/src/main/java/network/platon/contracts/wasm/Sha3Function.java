package network.platon.contracts.wasm;

import com.platon.rlp.datatypes.Uint32;
import java.util.Arrays;
import org.web3j.abi.WasmFunctionEncoder;
import org.web3j.abi.datatypes.WasmFunction;
import org.web3j.crypto.Credentials;
import org.web3j.protocol.Web3j;
import org.web3j.protocol.core.RemoteCall;
import org.web3j.tx.TransactionManager;
import org.web3j.tx.WasmContract;
import org.web3j.tx.gas.GasProvider;

/**
 * <p>Auto generated code.
 * <p><strong>Do not modify!</strong>
 * <p>Please use the <a href="https://docs.web3j.io/command_line.html">web3j command line tools</a>,
 * or the org.web3j.codegen.SolidityFunctionWrapperGenerator in the 
 * <a href="https://github.com/web3j/web3j/tree/master/codegen">codegen module</a> to update.
 *
 * <p>Generated with web3j version 0.7.5.8-SNAPSHOT.
 */
public class Sha3Function extends WasmContract {
    private static String BINARY_0 = "0x0061736d0100000001410b60000060017f0060027f7f0060047f7f7f7f006000017f60037f7f7f017f60017f017f60027f7f017f60037f7f7f0060057f7f7f7f7f017f60047f7f7f7f017f0282010603656e760c706c61746f6e5f6465627567000203656e760c706c61746f6e5f70616e6963000003656e760b706c61746f6e5f73686133000303656e7617706c61746f6e5f6765745f696e7075745f6c656e677468000403656e7610706c61746f6e5f6765745f696e707574000103656e760d706c61746f6e5f72657475726e0002037b7a000005050606060501070100060105050101010805060905070100050207060706010601070105060606060302080601050a02070a05060303050602020702070207020802080502020706070606060808070707070702030607020202030102070705030202070a0206070301050207080700060101060007000405017001040405030100020615037f0141e08c040b7f0041e08c040b7f0041df0c0b075406066d656d6f72790200115f5f7761736d5f63616c6c5f63746f727300060b5f5f686561705f6261736503010a5f5f646174615f656e6403020f5f5f66756e63735f6f6e5f65786974002006696e766f6b65007d0909010041010b03272b7a0a8b6b7a080010111078107f0b02000bc60a010b7f2002410f6a210341002104410020026b21052002410e6a2106410120026b21072002410d6a2108410220026b210902400340200020046a210b200120046a210a20022004460d01200a410371450d01200b200a2d00003a00002003417f6a2103200541016a21052006417f6a2106200741016a21072008417f6a2108200941016a2109200441016a21040c000b0b200220046b210c02400240024002400240200b410371220d450d00200c4120490d03200d4101460d01200d4102460d02200d4103470d03200b200120046a28020022063a0000200041016a210c200220046b417f6a21092004210b0240034020094113490d01200c200b6a220a2001200b6a220741046a2802002208411874200641087672360200200a41046a200741086a2802002206411874200841087672360200200a41086a2007410c6a2802002208411874200641087672360200200a410c6a200741106a2802002206411874200841087672360200200b41106a210b200941706a21090c000b0b2002417f6a2005416d2005416d4b1b20036a4170716b20046b210c2001200b6a41016a210a2000200b6a41016a210b0c030b200c210a02400340200a4110490d01200020046a220b200120046a2207290200370200200b41086a200741086a290200370200200441106a2104200a41706a210a0c000b0b02400240200c4108710d00200120046a210a200020046a21040c010b200020046a220b200120046a2204290200370200200441086a210a200b41086a21040b0240200c410471450d002004200a280200360200200a41046a210a200441046a21040b0240200c410271450d002004200a2f00003b0000200441026a2104200a41026a210a0b200c410171450d032004200a2d00003a000020000f0b200b200120046a220a28020022063a0000200b41016a200a41016a2f00003b0000200041036a210c200220046b417d6a21052004210b0240034020054111490d01200c200b6a220a2001200b6a220741046a2802002203410874200641187672360200200a41046a200741086a2802002206410874200341187672360200200a41086a2007410c6a2802002203410874200641187672360200200a410c6a200741106a2802002206410874200341187672360200200b41106a210b200541706a21050c000b0b2002417d6a2009416f2009416f4b1b20086a4170716b20046b210c2001200b6a41036a210a2000200b6a41036a210b0c010b200b200120046a220a28020022083a0000200b41016a200a41016a2d00003a0000200041026a210c200220046b417e6a21052004210b0240034020054112490d01200c200b6a220a2001200b6a220941046a2802002203411074200841107672360200200a41046a200941086a2802002208411074200341107672360200200a41086a2009410c6a2802002203411074200841107672360200200a410c6a200941106a2802002208411074200341107672360200200b41106a210b200541706a21050c000b0b2002417e6a2007416e2007416e4b1b20066a4170716b20046b210c2001200b6a41026a210a2000200b6a41026a210b0b0240200c411071450d00200b200a2d00003a0000200b200a280001360001200b200a290005370005200b200a2f000d3b000d200b200a2d000f3a000f200b41106a210b200a41106a210a0b0240200c410871450d00200b200a290000370000200b41086a210b200a41086a210a0b0240200c410471450d00200b200a280000360000200b41046a210b200a41046a210a0b0240200c410271450d00200b200a2f00003b0000200b41026a210b200a41026a210a0b200c410171450d00200b200a2d00003a00000b20000bfb0202027f017e02402002450d00200020013a0000200020026a2203417f6a20013a000020024103490d00200020013a0002200020013a00012003417d6a20013a00002003417e6a20013a000020024107490d00200020013a00032003417c6a20013a000020024109490d002000410020006b41037122046a2203200141ff017141818284086c22013602002003200220046b417c7122046a2202417c6a200136020020044109490d002003200136020820032001360204200241786a2001360200200241746a200136020020044119490d002003200136021820032001360214200320013602102003200136020c200241706a20013602002002416c6a2001360200200241686a2001360200200241646a20013602002001ad220542208620058421052004200341047141187222016b2102200320016a2101034020024120490d0120012005370300200141186a2005370300200141106a2005370300200141086a2005370300200141206a2101200241606a21020c000b0b20000b7a01027f200021010240024003402001410371450d0120012d0000450d02200141016a21010c000b0b2001417c6a21010340200141046a22012802002202417f73200241fffdfb776a7141808182847871450d000b0340200241ff0171450d01200141016a2d00002102200141016a21010c000b0b200120006b0b3a01017f23808080800041106b220141e08c84800036020c2000200128020c41076a41787122013602042000200136020020003f0036020c20000b120041808880800020004108108d808080000bc70101067f23808080800041106b22032480808080002003200136020c024002402001450d002000200028020c200241036a410020026b220471220520016a220641107622076a220836020c200020022000280204220120066a6a417f6a20047122023602040240200841107420024b0d002000410c6a200841016a360200200741016a21070b0240200740000d0041c88a808000108e808080000b20012003410c6a41041088808080001a200120056a21000c010b410021000b200341106a24808080800020000b180020002000108a808080001080808080001081808080000b2e000240418088808000200120006c22004108108d808080002201450d002001410020001089808080001a0b20010b02000b0f00418088808000108b808080001a0b3a01027f2000410120001b2101024003402001108c8080800022020d014100280290888080002200450d012000118080808000000c000b0b20020b0a0020001090808080000bce0301067f024020002001460d000240024002400240200120006b20026b410020024101746b4d0d0020012000734103712103200020014f0d012003450d02200021030c030b2000200120021088808080000f0b024020030d002001417f6a210402400340200020026a2203410371450d012002450d052003417f6a200420026a2d00003a00002002417f6a21020c000b0b2000417c6a21032001417c6a2104034020024104490d01200320026a200420026a2802003602002002417c6a21020c000b0b2001417f6a210103402002450d03200020026a417f6a200120026a2d00003a00002002417f6a21020c000b0b200241046a21052002417f7321064100210402400340200120046a2107200020046a2208410371450d0120022004460d03200820072d00003a00002005417f6a2105200641016a2106200441016a21040c000b0b200220046b2101410021030240034020014104490d01200820036a200720036a280200360200200341046a21032001417c6a21010c000b0b200720036a2101200820036a210320022006417c2006417c4b1b20056a417c716b20046b21020b03402002450d01200320012d00003a00002002417f6a2102200341016a2103200141016a21010c000b0b20000b4201027f0240024003402002450d0120002d0000220320012d00002204470d02200141016a2101200041016a21002002417f6a21020c000b0b41000f0b200320046b0b0900200041013602000b0900200041003602000b0900108780808000000b7701027f0240200241704f0d00024002402002410a4b0d00200020024101743a0000200041016a21030c010b200241106a417071220410928080800021032000200236020420002004410172360200200020033602080b200320012002109a808080001a200320026a41003a00000f0b108780808000000b1a0002402002450d0020002001200210888080800021000b20000b1d00024020002d0000410171450d0020002802081093808080000b20000b9a0101027f0240024020002d0000220541017122060d00200541017621050c010b200028020421050b02402004417f460d0020052001490d00200520016b2205200220052002491b21020240024020060d00200041016a21000c010b200028020821000b0240200020016a200320042002200220044b22001b109d808080002201450d0020010f0b417f200020022004491b0f0b108780808000000b190002402002450d002000200120021095808080000f0b41000b270020004200370200200041086a4100360200200020012001108a8080800010998080800020000b0900108780808000000bb60101037f4194888080001096808080004100280298888080002100024003402000450d01024003404100410028029c888080002202417f6a220136029c8880800020024101480d01200020014102746a22004184016a2802002102200041046a2802002100419488808000109780808000200220001181808080000041948880800010968080800041002802988880800021000c000b0b4100412036029c88808000410020002802002200360298888080000c000b0b0bcd0101027f419488808000109680808000024041002802988880800022030d0041a0888080002103410041a088808000360298888080000b02400240410028029c8880800022044120470d004184024101108f808080002203450d0141002104200341002802988880800036020041002003360298888080004100410036029c888080000b4100200441016a36029c88808000200320044102746a22034184016a2001360200200341046a200036020041948880800010978080800041000f0b419488808000109780808000417f0b6001017f23808080800041206b2202248080808000200241186a420037030020024200370310200242003703082000200241086a200110a38080800010a48080800010a5808080001a200241086a10a6808080001a200241206a2480808080000b4101017f23808080800041106b2202248080808000200020022001109e80808000220110e68080800021002001109b808080001a200241106a24808080800020000b23000240200028020c200041106a280200460d0041e98b808000108e808080000b20000b4e01017f20004200370200200041003602080240200128020420012802006b2202450d002000200210dc80808000200041086a2001280200200141046a280200200041046a10dd808080000b20000b19002000410c6a10de808080001a200010a8808080001a20000b0f0041a48a80800010a8808080001a0b2201017f024020002802002201450d002000200136020420011093808080000b20000b4701027f23808080800041206b22012480808080002000200141086a410010aa80808000220210a48080800010a5808080001a200210a6808080001a200141206a2480808080000b24002000420037020820004200370200200041106a42003702002000200110c3808080000b0f0041b08a80800010a8808080001a0b95010020004200370210200042ffffffff0f3702082000200129020037020002402002410871450d00200010ad8080800020012802044f0d00024020024104710d00200042003702000c010b41e28a808000108e808080000b024002402002411071450d00200010ad8080800020012802044d0d0020024104710d01200042003702000b20000f0b41f08a808000108e8080800020000b3400024002402000280204450d0020002802002c0000417f4c0d0141010f0b41000f0b200010ae80808000200010af808080006a0b280002402000280204450d0020002802002c0000417f4c0d0041000f0b200010b48080800041016a0b980401047f0240024002402000280204450d00200010b5808080004101210120002802002c00002202417f4c0d010c020b41000f0b0240200241ff0171220141b7014b0d00200141807f6a0f0b02400240200241ff0171220241bf014b0d000240200041046a22032802002202200141c97e6a22044b0d0041ff8a808000108e80808000200328020021020b024020024102490d0020002802002d00010d0041ff8a808000108e808080000b024020044105490d0041f08a808000108e808080000b024020002802002d00010d0041ff8a808000108e808080000b41002101410021020240034020042002460d012001410874200028020020026a41016a2d0000722101200241016a21020c000b0b200141384f0d0141ff8a808000108e8080800020010f0b0240200241f7014b0d00200141c07e6a0f0b0240200041046a22032802002202200141897e6a22044b0d0041ff8a808000108e80808000200328020021020b024020024102490d0020002802002d00010d0041ff8a808000108e808080000b024020044105490d0041f08a808000108e808080000b024020002802002d00010d0041ff8a808000108e808080000b41002101410021020240034020042002460d012001410874200028020020026a41016a2d0000722101200241016a21020c000b0b200141384f0d0041ff8a808000108e8080800020010f0b200141ff7d490d0041f08a808000108e8080800020010f0b20010b5102017f017e23808080800041306b220124808080800020012000290200220237031020012002370308200141186a200141086a411410ac8080800010ad808080002100200141306a24808080800020000b6a01037f02400240024020012802002204450d0041002105200320026a200128020422064b0d0120062002490d014100210120062003490d02200620026b20032003417f461b2101200420026a21050c020b410021050b410021010b20002001360204200020053602000b3901017f0240200110af80808000220220012802044d0d0041808c808000108e808080000b20002001200110ae80808000200210b1808080000bd003020a7f017e23808080800041c0006b220324808080800002402001280208220420024d0d00200341386a200110b280808000200320032903383703182001200341186a10b08080800036020c200341306a200110b280808000410021044100210541002106024020032802302207450d00410021054100210620032802342208200128020c2209490d00200820092009417f461b2105200721060b20012006360210200141146a2005360200200141086a41003602000b200141106a2106200141146a21092001410c6a2107200141086a210802400340200420024f0d012009280200450d01200341306a200110b28080800041002104024002402003280230220a450d00410021052003280234220b2007280200220c490d01200a200c6a2105200b200c6b21040c010b410021050b20092004360200200620053602002003200436022c2003200536022820032003290328370310200341306a20064100200341106a10b08080800010b18080800020062003290330220d37020020072007280200200d422088a76a3602002008200828020041016a22043602000c000b0b20032006290200220d3703202003200d3703082000200341086a411410ac808080001a200341c0006a2480808080000b4701017f4100210102402000280204450d00024020002802002d0000220041bf014b0d00200041b801490d01200041c97e6a0f0b200041f801490d00200041897e6a21010b20010b6601017f024020002802040d0041ff8a808000108e808080000b0240200028020022012d0000418101470d000240200041046a28020041014b0d0041ff8a808000108e80808000200028020021010b20012c00014100480d0041ff8a808000108e808080000b0b2d01017f2000200028020420012802002203200320012802046a10b7808080001a2000200210b88080800020000b970201057f23808080800041206b22042480808080000240200320026b22054101480d00024020052000280208200028020422066b4c0d00200441086a2000200520066a20002802006b10b980808000200120002802006b200041086a10ba8080800021060240034020032002460d01200641086a220528020020022d00003a00002005200528020041016a360200200241016a21020c000b0b20002006200110bb808080002101200610bc808080001a0c010b024002402005200620016b22074c0d00200041086a200220076a22082003200041046a10bd80808000200741014e0d010c020b200321080b200020012006200120056a10be8080800020022008200110bf808080001a0b200441206a24808080800020010bd00201087f02402001450d002000410c6a2102200041106a2103200041046a21040340200328020022052002280200460d010240200541786a28020020014f0d0041878b808000108e80808000200328020021050b200541786a2206200628020020016b220136020020010d0120032006360200200428020020002802006b2005417c6a28020022016b220510c08080800021062000200428020020002802006b22074101200641016a20054138491b22086a10c180808000200120002802006a220920086a2009200720016b1094808080001a02400240200541374b0d00200028020020016a200541406a3a00000c010b0240200641f7016a220741ff014b0d00200028020020016a20073a00002000280200200620016a6a210103402005450d02200120053a0000200541087621052001417f6a21010c000b0b419b8b808000108e808080000b410121010c000b0b0b4c01017f02402001417f4c0d0041ffffffff0721020240200028020820002802006b220041feffffff034b0d0020012000410174220020002001491b21020b20020f0b2000109f80808000000b5401017f410021042000410036020c200041106a200336020002402001450d00200110928080800021040b200020043602002000200420026a22023602082000410c6a200420016a3602002000200236020420000b8c0101027f20012802042103200041086a220420002802002002200141046a10e380808000200420022000280204200141086a10e980808000200028020021022000200128020436020020012002360204200028020421022000200128020836020420012002360208200028020821022000200128020c3602082001200236020c2001200128020436020020030b2301017f200010e480808000024020002802002201450d0020011093808080000b20000b2e000240200220016b22024101480d002003280200200120021088808080001a2003200328020020026a3602000b0b5c01037f200041046a21042000280204220521062001200520036b6a2203210002400340200020024f0d01200620002d00003a00002004200428020041016a2206360200200041016a21000c000b0b20012003200510e8808080001a0b21000240200120006b2201450d002002200020011094808080001a0b200220016a0b2501017f41002101024003402000450d0120004108762100200141016a21010c000b0b20010b4001027f02402000280204200028020022026b220320014f0d002000200120036b10c2808080000f0b0240200320014d0d00200041046a200220016a3602000b0b920101027f23808080800041206b2202248080808000024002402000280208200028020422036b20014f0d00200241086a2000200320016a20002802006b10b980808000200041046a28020020002802006b200041086a10ba808080002203200110ea808080002000200310e280808000200310bc808080001a0c010b2000200110eb808080000b200241206a2480808080000b7501017f23808080800041106b2202248080808000024002402001450d00200220013602002002200028020420002802006b3602042000410c6a200210c4808080000c010b20024100360208200242003703002000200210c5808080001a200210a8808080001a0b200241106a24808080800020000b3d01017f02402000280204220220002802084f0d0020022001290200370200200041046a2200200028020041086a3602000f0b2000200110c6808080000b5101027f23808080800041106b22022480808080002002200128020022033602082002200128020420036b36020c200220022903083703002000200210c7808080002101200241106a24808080800020010b840101027f23808080800041206b2202248080808000200241086a2000200028020420002802006b41037541016a10ec80808000200028020420002802006b410375200041086a10ed80808000220328020820012902003702002003200328020841086a3602082000200310ee80808000200310ef808080001a200241206a2480808080000b800102027f017e23808080800041206b2202248080808000024002402001280204220341374b0d002002200341406a3a001f20002002411f6a10c8808080000c010b2000200341f70110c9808080000b200220012902002204370310200220043703082000200241086a410110b6808080002100200241206a24808080800020000b3d01017f02402000280204220220002802084f0d00200220012d00003a0000200041046a2200200028020041016a3602000f0b2000200110ca808080000b6401027f23808080800041106b22032480808080000240200110c080808000220420026a2202418002480d0041d18b808000108e808080000b200320023a000f20002003410f6a10c88080800020002001200410cb80808000200341106a2480808080000b7e01027f23808080800041206b2202248080808000200241086a2000200028020441016a20002802006b10b980808000200028020420002802006b200041086a10ba80808000220328020820012d00003a00002003200328020841016a3602082000200310e280808000200310bc808080001a200241206a2480808080000b44002000200028020420026a20002802006b10c1808080002000280204417f6a2100024003402001450d01200020013a00002000417f6a2100200141087621010c000b0b0bfc0101037f23808080800041206b22032480808080002001280200210420012802042105024002402002450d004100210102400340200420016a2102200120054f0d0120022d00000d01200141016a21010c000b0b200520016b21050c010b200421020b0240024002400240024020054101470d0020022c00004100480d012000200210cd808080000c040b200541374b0d010b20032005418001733a001f20002003411f6a10c8808080000c010b2000200541b70110c9808080000b2003200536021420032002360210200320032903103703082000200341086a410010b6808080001a0b2000410110b880808000200341206a24808080800020000b3d01017f0240200028020422022000280208460d00200220012d00003a0000200041046a2200200028020041016a3602000f0b2000200110ce808080000b7e01027f23808080800041206b2202248080808000200241086a2000200028020441016a20002802006b10b980808000200028020420002802006b200041086a10ba80808000220328020820012d00003a00002003200328020841016a3602082000200310e280808000200310bc808080001a200241206a2480808080000bdb0201047f2380808080004190026b2202248080808000024002400240200110d080808000450d00200141b78b80800010d180808000450d012002200110d2808080003a008f0220002002418f026a10c8808080000c020b200041b78b80800010cd808080000c010b200241c8016a200141c0001088808080001a200241c8006a200241c8016a41c0001088808080001a02400240200241c8006a10d380808000220341374b0d0020022003418001733a008f0220002002418f026a10c8808080000c010b0240200310d480808000220441b7016a2205418002490d0041b88b808000108e808080000b200220053a008f0220002002418f026a10c88080800020002003200410d5808080000b20024188016a200141c0001088808080001a200241086a20024188016a41c0001088808080001a2000200241086a200310d6808080000b2000410110b88080800020024190026a24808080800020000b3b01017f23808080800041106b22012480808080002001410036020c20002001410c6a10d7808080002100200141106a24808080800020004101730b5101017f2380808080004180016b2202248080808000200241c0006a200041c0001088808080001a200241c0006a200220012d000010d88080800010d980808000210120024180016a24808080800020010b3701027f410021012000413f6a210241012100024003402000410171450d01200120022d0000722101410021000c000b0b200141ff01710b5701027f23808080800041106b220124808080800041002102024003402001410036020c20002001410c6a10da80808000450d012000410810db808080001a200241016a21020c000b0b200141106a24808080800020020b2501017f41002101024003402000450d0120004108762100200141016a21010c000b0b20010b44002000200028020420026a20002802006b10c1808080002000280204417f6a2100024003402001450d01200020013a00002000417f6a2100200141087621010c000b0b0b54002000200028020420026a20002802006b10c1808080002000280204417f6a210002400340200110d080808000450d012000200110d2808080003a00002001410810db808080001a2000417f6a21000c000b0b0b6d01047f23808080800041c0006b22022480808080002002200128020010df808080001a4100210141012103024003402001413f4b0d01200020016a2104200220016a2105200141016a210120042d000020052d0000460d000b410021030b200241c0006a24808080800020030b1b002000410041c0001089808080002200200110e18080800020000b7501047f23808080800041c0006b22022480808080002002200141c00010888080800021034100210441002101024003402001413f4b0d01200020016a2102200320016a2105200141016a210120022d0000220220052d00002205460d000b200220054921040b200341c0006a24808080800020040b5101017f2380808080004180016b2202248080808000200241c0006a200041c0001088808080001a200241c0006a2002200128020010df8080800010f580808000210120024180016a24808080800020010b3f01017f23808080800041c0006b220224808080800020022000200110f6808080002000200241c0001088808080002100200241c0006a24808080800020000b3801017f02402001417f4c0d00200020011092808080002202360200200020023602042000200220016a3602080f0b2000109f80808000000b2e000240200220016b22024101480d002003280200200120021088808080001a2003200328020020026a3602000b0b2201017f024020002802002201450d002000200136020420011093808080000b20000b1b002000410041c0001089808080002200200110e08080800020000b6802017f027e2000413f6a21022001ac2103420021040240034020044220510d01200220032004873c00002002417f6a2102200442087c21040c000b0b2001411f752101413b2102024003402002417f460d01200020026a20013a00002002417f6a21020c000b0b0b2d00200020013a003f413e2101024003402001417f460d01200020016a41003a00002001417f6a21010c000b0b0b7001017f200041086a20002802002000280204200141046a10e380808000200028020021022000200128020436020020012002360204200028020421022000200128020836020420012002360208200028020821022000200128020c3602082001200236020c200120012802043602000b2f01017f20032003280200200220016b22026b2204360200024020024101480d002004200120021088808080001a0b0b0f002000200028020410e5808080000b2d01017f20002802082102200041086a21000240034020012002460d0120002002417f6a22023602000c000b0b0b4501017f23808080800041106b22022480808080002002200241086a200110e78080800029020037030020002002410010cc808080002100200241106a24808080800020000b360020002001280208200141016a20012d00004101711b3602002000200128020420012d0000220141017620014101711b36020420000b23000240200120006b2201450d00200220016b2202200020011094808080001a0b20020b2e000240200220016b22024101480d002003280200200120021088808080001a2003200328020020026a3602000b0b3401017f20002802082102200041086a21000340200241003a00002000200028020041016a22023602002001417f6a22010d000b0b3401017f20002802042102200041046a21000340200241003a00002000200028020041016a22023602002001417f6a22010d000b0b5301017f024020014180808080024f0d0041ffffffff0121020240200028020820002802006b220041037541feffffff004b0d0020012000410275220020002001491b21020b20020f0b2000109f80808000000b5c01017f410021042000410036020c200041106a200336020002402001450d002003200110f08080800021040b200020043602002000200420024103746a22033602082000410c6a200420014103746a3602002000200336020420000b7001017f200041086a20002802002000280204200141046a10f180808000200028020021022000200128020436020020012002360204200028020421022000200128020836020420012002360208200028020821022000200128020c3602082001200236020c200120012802043602000b2301017f200010f280808000024020002802002201450d0020011093808080000b20000b0e0020002001410010f3808080000b2f01017f20032003280200200220016b22026b2204360200024020024101480d002004200120021088808080001a0b0b0f002000200028020410f4808080000b2300024020014180808080024f0d0020014103741092808080000f0b108780808000000b2d01017f20002802082102200041086a21000240034020012002460d012000200241786a22023602000c000b0b0b0f002000200110f7808080004101730bc10201087f23808080800041c0006b2203248080808000024002402002418004490d002000410010df808080001a0c010b02402002450d002003200141c00010888080800021040240024020024107712205450d00200420042d003f20057622063a003f410820056b210741002101024003402001413e6a4100480d01200420016a2208413e6a220920092d00002209200576220a3a00002008413f6a20062009200774723a00002001417f6a2101200a21060c000b0b200220056b2202450d010b2004200241086d22066b2108413f21010240034020012006480d01200420016a200820016a2d00003a00002001417f6a21010c000b0b034020014100480d01200420016a41003a00002001417f6a21010c000b0b2000200441c0001088808080001a0c010b2000200141c0001088808080001a0b200341c0006a2480808080000b6e01047f23808080800041c0006b22022480808080002002200141c00010888080800021034100210141012104024003402001413f4b0d01200020016a2102200320016a2105200141016a210120022d000020052d0000460d000b410021040b200341c0006a24808080800020040b4a0041a48a80800041e18a80800010a280808000418180808000410041808880800010a1808080001a41b08a80800010a980808000418280808000410041808880800010a1808080001a0b2201017f024020002802002201450d002000200136020420011093808080000b20000b0f0041bc8a808000109b808080001a0b180020002000108a808080001080808080001081808080000b7b01027f23808080800041306b22012480808080002001412c6a41002d008d8c8080003a0000200141186a4200370300200141106a42003703002001420037030820014200370300200141002800898c808000360228200141286a41052001412010828080800020012802002102200141306a24808080800020020b8d0d02077f027e23808080800041f0016b220024808080800010868080800020004190016a410036020020004200370388014100210102404100410c460d00034020004188016a20016a4100360200200141046a2201410c470d000b0b20004100360280012000420037037802400240024002401083808080002201450d002001417f4c0d0220004180016a2001109280808000220241002001108980808000220320016a22013602002000200136027c20002003360278200321030c010b4100210141002102410021030b2003108480808000200020023602b0012000200120026b3602b401200020002903b001370340200041c8006a200041e0006a200041c0006a411c10ac80808000410010b3808080000240024002400240024002400240024002400240200028024c450d0020002802482d000041c0014f0d000240200041c8006a10af808080002202200028024c22014d0d0041d68c80800010fb80808000200028024c21010b200028024821042001450d014100210320042c00002205417f4a0d04200541ff0171220641bf014b0d0241002103200541ff017141b801490d03200641c97e6a21030c030b200041b8016a4100360200200042003703b001410021014100410c460d080340200041b0016a20016a4100360200200141046a2201410c470d000c090b0b4101210320040d020c030b41002103200541ff017141f801490d00200641897e6a21030b200341016a21030b200320026a20014b0d0020012002490d0020012003490d00200120036b20022002417f461b2202200041c8006a10af808080002201490d01200041b8016a4100360200200042003703b001200220012001417f461b220541704f0d06200420036a220220056a21032005410a4d0d02200541106a41707122041092808080002101200020053602b401200020044101723602b001200020013602b8010c030b200041c8006a10af808080001a0b41002105200041b8016a4100360200200042003703b00141002102410021030b200020054101743a00b001200041b0016a41017221010b024020032002460d000340200120022d00003a0000200141016a21012003200241016a2202470d000b0b200141003a00000b0240024020002d0088014101710d00200041003b0188010c010b20002802900141003a00002000410036028c0120002d008801410171450d0020004190016a28020010938080800020004100360288010b20004188016a41086a200041b0016a41086a280200360200200020002903b001370388014100210102404100410c460d000340200041b0016a20016a4100360200200141046a2201410c470d000b0b200041b0016a109b808080001a024002400240200028028c0120002d008801220141017620014101711b450d0020004188016a419c8c80800010fe808080000d0220004188016a41a18c80800010fe80808000450d0120004197016a10fc8080800021024200210720004198016a41106a420037030020004198016a41086a42003703002000420037039801200041e8016a4200370300200041e0016a4200370300200041d8016a4200370300200041d0016a4200370300200041c8016a4200370300200041b0016a41106a4200370300200041b0016a41086a4200370300200042003703b001200041ef016a21012002ad2108024042004220510d000340200120082007883c00002001417f6a2101200742087c22074220520d000b0b413b21010240413b417f460d000340200041b0016a20016a41003a00002001417f6a2201417f470d000b0b200041386a200041b0016a41386a290300370300200041306a200041b0016a41306a290300370300200041286a200041b0016a41286a290300370300200041206a200041b0016a41206a290300370300200041186a200041b0016a41186a290300370300200041106a200041b0016a41106a290300370300200041086a200041b0016a41086a290300370300200020002903b00137030020004198016a200010cf808080001a024020002802a40120004198016a41106a280200460d0041bf8c80800010fb808080000b2000280298012201200028029c0120016b1085808080000240200041a4016a2802002201450d00200041a8016a200136020020011093808080000b20004198016a10f9808080001a0c020b418e8c80800010fb808080000c010b41ac8c80800010fb808080000b200041f8006a10f9808080001a20004188016a109b808080001a200041f0016a2480808080000f0b200041f8006a109f80808000000b200041b0016a109880808000000b4201037f4100210202402001108a808080002203200028020420002d0000220441017620044101711b470d0020004100417f20012003109c808080004521020b20020b5501017f410042003702bc8a808000410041003602c48a8080004174210002404174450d000340200041c88a8080006a4100360200200041046a22000d000b0b418380808000410041808880800010a1808080001a0b0bee0402004180080bc802000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000041c80a0b97026661696c656420746f20616c6c6f6361746520706167657300006f7665722073697a6520726c7000756e6465722073697a6520726c700062616420726c70006974656d436f756e7420746f6f206c61726765006974656d436f756e7420746f6f206c6172676520666f7220524c5000804e756d62657220746f6f206c6172676520666f7220524c5000436f756e7420746f6f206c6172676520666f7220524c50006c697374537461636b206973206e6f7420656d70747900626164206361737400010710203076616c6964206d6574686f640a00696e69740053686133526573756c74006e6f206d6574686f6420746f2063616c6c0a006c697374537461636b206973206e6f7420656d70747900626164206361737400";

    public static String BINARY = BINARY_0;

    public static final String FUNC_SHA3RESULT = "Sha3Result";

    protected Sha3Function(String contractAddress, Web3j web3j, Credentials credentials, GasProvider contractGasProvider) {
        super(BINARY, contractAddress, web3j, credentials, contractGasProvider);
    }

    protected Sha3Function(String contractAddress, Web3j web3j, TransactionManager transactionManager, GasProvider contractGasProvider) {
        super(BINARY, contractAddress, web3j, transactionManager, contractGasProvider);
    }

    public static RemoteCall<Sha3Function> deploy(Web3j web3j, Credentials credentials, GasProvider contractGasProvider) {
        String encodedConstructor = WasmFunctionEncoder.encodeConstructor(BINARY, Arrays.asList());
        return deployRemoteCall(Sha3Function.class, web3j, credentials, contractGasProvider, encodedConstructor);
    }

    public static RemoteCall<Sha3Function> deploy(Web3j web3j, TransactionManager transactionManager, GasProvider contractGasProvider) {
        String encodedConstructor = WasmFunctionEncoder.encodeConstructor(BINARY, Arrays.asList());
        return deployRemoteCall(Sha3Function.class, web3j, transactionManager, contractGasProvider, encodedConstructor);
    }

    public RemoteCall<Uint32> Sha3Result() {
        final WasmFunction function = new WasmFunction(FUNC_SHA3RESULT, Arrays.asList(), Uint32.class);
        return executeRemoteCall(function, Uint32.class);
    }

    public static Sha3Function load(String contractAddress, Web3j web3j, Credentials credentials, GasProvider contractGasProvider) {
        return new Sha3Function(contractAddress, web3j, credentials, contractGasProvider);
    }

    public static Sha3Function load(String contractAddress, Web3j web3j, TransactionManager transactionManager, GasProvider contractGasProvider) {
        return new Sha3Function(contractAddress, web3j, transactionManager, contractGasProvider);
    }
}