package plugin

import (
	"math/big"

	"github.com/PlatONnetwork/PlatON-Go/common"
)

type issue1625Accounts struct {
	addr   common.Address
	amount *big.Int
}

func NewIssue1625Accounts() ([]issue1625Accounts, error) {
	var accounts = [][]string{
		{"atx1w5wtf7fud8xppguaxeg68z4g3qe0p99yq863qr", "10000000000000000000000"},
		{"atx1jt7zh76t9xdfnczz7usymwxmmfr9zza79p74jv", "10000000000000000000000"},
		{"atx1999r47ahhuc4xhjdhppvgkyw3z03gu765vtx2u", "10000000000000000000000"},
		{"atx1tuhqwl39xa3lgy5cv732tamxx0efteavk6qys7", "10000000000000000000000"},
		{"atx1sh3xncmqdwlafptt9nsc3vv5zzcpufdue0lsua", "10000000000000000000000"},
		{"atx1ytmtawyxxt0kd44nx772h9qqk3wl28u89hhsv0", "10000000000000000000000"},
		{"atx14v68yv0a7a5jphvul2ehm24jt446lq6j2mmwcc", "10000000000000000000000"},
		{"atx1ayzxme7s9apmaejvf6n83uk63ldjshv38809va", "10000000000000000000000"},
		{"atx1x78f927l260rdp9erk8er3jhp26mfwcpfsm98z", "10000000000000000000000"},
		{"atx1jxmdq2gxetydthddqqfrlr8mtm44afwhms2mmy", "10000000000000000000000"},
		{"atx1xc4nl4s5m583xkcfwq4na32hvs7pz8r4a48q8t", "10000000000000000000000"},
		{"atx1enmarze9cu2tzp37g9lg3fqkkkx752d2cm77ku", "10000000000000000000000"},
		{"atx1g8yqegahap0c5jfxkdezl7nscrn0a0xxydf2ay", "10000000000000000000000"},
		{"atx1tx358vuju4qag2mr92nl20552ysyzuhyxz88ke", "10000000000000000000000"},
		{"atx177fthwhxecq5dn4m7j7v8hs7xprhcny45psjdf", "10000000000000000000000"},
		{"atx14gt343vhqs74uesz06m3ugh6gnhzrethjyktgh", "10000000000000000000000"},
		{"atx14kfvcwzet5gkxr0ur8057ne5wzhul98qg8r4sn", "10000000000000000000000"},
		{"atx1smhd0m7nwkaextuuylf864vhmkl09k53jpg8vs", "10000000000000000000000"},
		{"atx1cywdp07xr2cjhdcr8pp8vgz7j532epdck59atq", "10000000000000000000000"},
		{"atx10pn9uxafqps54zj5r3a7ppfgzzee7w253nuy29", "10000000000000000000000"},
		{"atx18r03ct2azfjsktrmndfxdvam58wlupgsnm7pvu", "10000000000000000000000"},
		{"atx14lnu29vj0aulh3myqethsdq5uvyyuu9uxskgue", "10000000000000000000000"},
		{"atx1ggmrtps9knqdhsm5y7pj844ktk4h68cgmqpcrx", "10000000000000000000000"},
		{"atx1tyuyylmskkase20zmns7zsv3mw89skn2lwqctz", "10000000000000000000000"},
		{"atx14tqp8t2f63x9fwgn4xvq26q5ehyrn0rv5lfp24", "10000000000000000000000"},
		{"atx1glkegnc2l8mtfuet8x676gjk6s58lj9aa5k2e7", "10000000000000000000000"},
		{"atx1ww44fcuk9xtlegfes7quueakcj7upr7h7rzd5x", "10000000000000000000000"},
		{"atx1kv9gac3mpdzrg2qyrajfjlzq6j3defhrdgh8hs", "10000000000000000000000"},
		{"atx1u4fpcvdu2px4kwesug69t70r2053gs7at44f9u", "10000000000000000000000"},
		{"atx1xrnu97t2zehuwkyywrkq84n2ns6uyrrn5h8q2l", "10000000000000000000000"},
	}
	Accounts := make([]issue1625Accounts, 0)
	for _, info := range accounts {
		addr, err := common.Bech32ToAddress(info[0])
		if err != nil {
			return nil, err
		}
		amount, _ := new(big.Int).SetString(info[1], 10)

		Accounts = append(Accounts, issue1625Accounts{
			addr, amount,
		})
	}

	return Accounts, nil
}
