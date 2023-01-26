package driver

import (
	"os"
)

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

func GetCredentials() (address string, username string, password string) {
	//get creadentials mikrotik

	address = os.Getenv("ROUTER_IP")
	username = os.Getenv("ROUTER_USER")
	password = os.Getenv("ROUTER_PASSWD")

	return address, username, password
}
