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
		{"atp1eshshnxuva6f4zqmwj9xszfj65y5vhalr7nyed", "9991365936677426200750"},
		{"atp166ue9gzupre59qsj9xvdxjwrzdrheentp9xlue", "2000000000000000000000"},
		{"atp13fd8zf6gp8jjn46uvlr5q5ayz73qgpzrzcwfdv", "4900000000000000000000"},
		{"atp15ftdv7s5sn7tswnqwegaxrdzhh9cezyaeuv7us", "9991321838252569163750"},
		{"atp1sp0uwm79fnva0x86kzj979psk7f8zpa9zark8j", "9991349622483553702000"},
		{"atp1std7ff5cjdezwe6pz7eq278jpmcaq2mef6d4e6", "9991349622483553702000"},
		{"atp1uyu68y0ygk5gzhg8402qmtg8hj8qwc4je3a2ym", "9991349622483553702000"},
		{"atp1hgerqd2erw89a6f0n9cch6azdm3mlkn7lnq0gp", "9991349622483553702000"},
		{"atp16wrzwyktqc09kjeydnq4jdtuuu793ced6jwyh2", "10000000000000000000000"},
		{"atp1wt9j640hw046jt86ttp7h0mvv66ctkplazxdxx", "9991367115042614713250"},
		{"atp1gkaqpw6q0syy8865yg7f9zd5fmt9mec4k4kx2z", "9991374220297723672250"},
		{"atp1d8paww9mxjar5rrek4mken8ca448xyp2hxsmx0", "10000000000000000000000"},
		{"atp1u9xc3rejevx4ux5p6fl5ghlgml5sa7mnnrwfkv", "10000000000000000000000"},
		{"atp1t0p4xrxq3rw24709rqvudc4ydph5gv43nfe9gj", "10000000000000000000000"},
		{"atp1fklu5nnuvjwhx50vz44y32upf69glx40yux69x", "9991374234971914133750"},
		{"atp1pv4eep9vzyy4rekyfs5xnve3kv20hhgqkdpnnm", "10000000000000000000000"},
		{"atp1jp7m5g4lwmk36932hxnxeupul93kcdwtlfq4p4", "10000000000000000000000"},
		{"atp1zsjz0zapqxe82amc98208x5992nnclvgynhhe4", "10000000000000000000000"},
		{"atp12jsctngpm8tn30scuc0dmpg9ncsc82rlz9vwjj", "10000000000000000000000"},
		{"atp1tql5puw2kf84d7czgh7lucrpd9vhm77h8fc5j6", "10000000000000000000000"},
		{"atp1g60w9nkqehx0w32hc8lhncyar3u9tgdhtyrpgp", "10000000000000000000000"},
		{"atp143u7s2rxqslvqsw8ehuj362kqfjsjlv8d3vjv8", "10000000000000000000000"},
		{"atp19kls6uhznrfcc2yg2q79sdcjfu780e6fhhg472", "10000000000000000000000"},
		{"atp1fj30cwd58djru06h8wwdrcdwh8w0p8qvjytck3", "10000000000000000000000"},
		{"atp1k5fh43wnd0gw339glgeyusjdzfm8cxdcw0nhgq", "10000000000000000000000"},
		{"atp1j6uj6dl06rjqxw73aknejalp78u04ru4ljwswf", "10000000000000000000000"},
		{"atp1guuc25qpqjghellen7uzhpetfdgevcx038nw4f", "10000000000000000000000"},
		{"atp1nvuk3v90cttx3zt5vzy7rnh4kny7tcy9uzftjh", "10000000000000000000000"},
		{"atp1mu0aws2g65hw0z77gzukzed83j3tcjuz366wn3", "10000000000000000000000"},
		{"atp1cwar0uy0lgz3vmfnu5ffagkrmqslpsh82v0vpc", "10000000000000000000000"},
		{"atp12k5nrdn29msfalhfe6jqw8y66m7fq4r5an376z", "10000000000000000000000"},
		{"atp1rknphx8hdq7rnkphy9j4p3cfa8y604fu5cxzew", "10000000000000000000000"},
		{"atp1m6l2jjr39qwf6kuwea5cfa0mn9nv96yt5kxxzx", "10000000000000000000000"},
		{"atp1krv3qnd08ke4y52ylvurdfxen7aq7wehlzqm9n", "10000000000000000000000"},
		{"atp18xtt0sqrg9qcd83u6659nhhdcs7kxzmyr3f44w", "10000000000000000000000"},
		{"atp1t32d3ellldszx5j240ru9qp568umpve9ps7pnk", "10000000000000000000000"},
		{"atp175z4sfjg33r0svp9cdra8hpasgfuxeug4h8fps", "10000000000000000000000"},
		{"atp1vxc0074lv8y3078v3ms9r0trktqqcgpu36xuz9", "10000000000000000000000"},
		{"atp1x8fv9scyvsmfnkf3utzpszrs52hv7e5z07z8sx", "10000000000000000000000"},
		{"atp18c6mfgzy68pfwqu4cnu4nwtl367qmsjy3j20d5", "10000000000000000000000"},
		{"atp1t5e4qzgygkhtj8tjgn72qtqtfu88jdmemsd59t", "10000000000000000000000"},
		{"atp1x2xgnxd3939h9yh07ay3xm7dh6peunsyae3mtz", "10000000000000000000000"},
		{"atp1y4z8ny5c0nhy6uan4nzx2mzv8jp9ntq7quxvgp", "10000000000000000000000"},
		{"atp1q3ffp495wavsuljmgrg3ds06z9ykpgepqa3ypx", "10000000000000000000000"},
		{"atp1npnfxe5gfcczw2659p4pk2ezwrprgvmjruhzfq", "10000000000000000000000"},
		{"atp1225wvv6ug5y28dp0rrr8w9nxedkm8hcglfnk0d", "10000000000000000000000"},
		{"atp1s3rfmcmm4zdz529rvwuagfvqrq5c44e2ja46kz", "10000000000000000000000"},
		{"atp1fnehxk8jcwaquerf4trx03l9na7a56yyfcjklh", "10000000000000000000000"},
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
