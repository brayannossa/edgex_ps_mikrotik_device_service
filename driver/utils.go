package driver

type Interface struct {
	Name   string `json:"name"`
	TxByte string `json:"tx-byte"`
	RxByte string `json:"rx-byte"`
}

type LTE struct {
	RSSI string `json:"rssi"`
	RSRP string `json:"rsrp"`
	RSCP string `json:"rscp"`
}
