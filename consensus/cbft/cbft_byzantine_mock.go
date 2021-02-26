package cbft

import "github.com/pingcap/failpoint"

func (cbft *Cbft) byzantineMock() {
	cbft.byzantinePBHandler()
	cbft.byzantineVTHandler()
	cbft.byzantineVCHandler()
}

func (cbft *Cbft) byzantinePBHandler() {
	node, err := cbft.isCurrentValidator()
	if err != nil || node == nil {
		return
	}
	failpoint.Inject("Byzantine-PB01", func() {
		cbft.MockPB01(node.Index)
	})
	failpoint.Inject("Byzantine-PB02", func() {
		cbft.MockPB02(node.Index)
	})
	failpoint.Inject("Byzantine-PB03", func() {
		cbft.MockPB03(node.Index)
	})
	failpoint.Inject("Byzantine-PB04", func() {
		cbft.MockPB04(node.Index)
	})
	failpoint.Inject("Byzantine-PB06", func(value failpoint.Value) {
		cbft.MockPB06(node.Index, value)
	})
	failpoint.Inject("Byzantine-PB07", func() {
		cbft.MockPB07(node.Index)
	})
	failpoint.Inject("Byzantine-PB08", func() {
		cbft.MockPB08(node.Index)
	})
	failpoint.Inject("Byzantine-PB09", func() {
		cbft.MockPB09(node.Index)
	})
	failpoint.Inject("Byzantine-PB10", func() {
		cbft.MockPB10(node.Index)
	})
	failpoint.Inject("Byzantine-PB11", func() {
		cbft.MockPB11(node.Index)
	})
	failpoint.Inject("Byzantine-PB12", func() {
		cbft.MockPB12(node.Index)
	})
}

func (cbft *Cbft) byzantineVTHandler() {
	node, err := cbft.isCurrentValidator()
	if err != nil || node == nil {
		return
	}
	failpoint.Inject("Byzantine-VT01", func() {
		cbft.MockVT01(node.Index)
	})
	failpoint.Inject("Byzantine-VT02", func() {
		cbft.MockVT02(node.Index)
	})
	failpoint.Inject("Byzantine-VT03", func() {
		cbft.MockVT03(node.Index)
	})
	failpoint.Inject("Byzantine-VT04", func() {
		cbft.MockVT04(node.Index)
	})
	failpoint.Inject("Byzantine-VT05", func() {
		cbft.MockVT05(node.Index)
	})
	failpoint.Inject("Byzantine-VT06", func() {
		cbft.MockVT06(node.Index)
	})
	failpoint.Inject("Byzantine-VT07", func() {
		cbft.MockVT07(node.Index)
	})
	failpoint.Inject("Byzantine-VT08", func() {
		cbft.MockVT08(node.Index)
	})
}

func (cbft *Cbft) byzantineVCHandler() {
	failpoint.Inject("Byzantine-VC01", func() {
		cbft.MockVC01()
	})
	failpoint.Inject("Byzantine-VC02", func() {
		cbft.MockVC02()
	})
	failpoint.Inject("Byzantine-VC03", func() {
		cbft.MockVC03()
	})
	failpoint.Inject("Byzantine-VC04", func() {
		cbft.MockVC04()
	})
	failpoint.Inject("Byzantine-VC05", func() {
		cbft.MockVC05()
	})
}
