package driver

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"io/ioutil"
	"math"
	"net/http"
	"strconv"
	"time"
)

type Mikrotik struct {
	Address  string
	Username string
	Password string
}

func NewMikrotik() *Mikrotik {
	address, username, password := GetCredentials()
	return &Mikrotik{
		Address:  address,
		Username: username,
		Password: password,
	}
}

func (m Mikrotik) DataRate() (tx float64, rx float64, err error) {
	// This function calculates the transmission and
	// reception rate of the LTE interface in KB/s
	period := 5
	sample := 1.0
	for i := 0.0; i < sample; i++ {
		tx1, rx1, err := m.requestTxRx()
		if err != nil {
			return 0, 0, err
		}
		t1 := time.Now()
		time.Sleep(time.Duration(period) * time.Second)
		tx2, rx2, err := m.requestTxRx()
		if err != nil {
			return 0, 0, err
		}
		t2 := time.Now()
		delay := float64(t2.Sub(t1)/time.Nanosecond) / 1e+09
		tx = tx + (tx2-tx1)/delay
		rx = rx + (rx2-rx1)/delay
	}
	// averages the transmission and reception
	// rates and converts them from B/s to KB/s
	tx = (tx / sample) / 1000
	rx = (rx / sample) / 1000
	//round to two decimal places
	ratio := math.Pow(10, float64(2))
	tx = math.Round(tx*ratio) / ratio
	rx = math.Round(rx*ratio) / ratio
	return tx, rx, nil
}

func (m Mikrotik) requestTxRx() (txByte float64, rxByte float64, err error) {
	// This function is responsible for making the request
	// to the router to obtain information about the bytes
	// transmitted and received in the lte interface
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Timeout: 2 * time.Second, Transport: tr}
	req, err := http.NewRequest("GET", "https://"+m.Address+"/rest/interface", nil)
	req.SetBasicAuth(m.Username, m.Password)
	resp, err := client.Do(req)
	if err != nil {
		return 0, 0, err
	}
	body, err := ioutil.ReadAll(resp.Body)
	var interfaces []Interface
	json.Unmarshal(body, &interfaces)

	for i := 0; i < len(interfaces); i++ {
		if interfaces[i].Name == "lte1" {
			txByte, _ := strconv.ParseFloat(interfaces[i].TxByte, 64)
			rxByte, _ := strconv.ParseFloat(interfaces[i].RxByte, 64)
			return txByte, rxByte, nil
		}

	}
	return 0, 0, nil

}

func (m Mikrotik) SignalQuality() (rssi float64, rsrp float64, rscp float64, err error) {
	// This function obtains the parameters that measure
	// the signal quality of the LTE interface
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Timeout: 2 * time.Second, Transport: tr}

	var jsonData = []byte(`{
		"numbers": 0,
		"duration": 1
	}`)

	req, err := http.NewRequest("POST", "https://"+m.Address+"/rest/interface/lte/monitor", bytes.NewBuffer(jsonData))
	if err != nil {
		return 0, 0, 0, err
	}
	req.Header.Set("Content-Type", "application/json; charset=UTF-8")
	req.SetBasicAuth(m.Username, m.Password)

	resp, err := client.Do(req)
	if err != nil {
		return 0, 0, 0, err
	}
	body, err := ioutil.ReadAll(resp.Body)
	var data []LTE
	json.Unmarshal(body, &data)
	if err != nil {
		return 0, 0, 0, err
	}
	rssi, _ = strconv.ParseFloat(data[0].RSSI, 64)
	rsrp, _ = strconv.ParseFloat(data[0].RSRP, 64)
	rscp, _ = strconv.ParseFloat(data[0].RSCP, 64)
	return rssi, rsrp, rscp, nil

}
