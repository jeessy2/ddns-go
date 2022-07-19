package config

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os/exec"
	"runtime"
)

// PveJson pvesh返回Json
type PveJson struct {
	Result []struct {
		HardwareAddress string `json:"hardware-address"`
		IPAddresses     []struct {
			IPAddress     string `json:"ip-address"`
			IPAddressType string `json:"ip-address-type"`
			Prefix        int    `json:"prefix"`
		} `json:"ip-addresses"`
		Name       string `json:"name"`
		Statistics struct {
			RxBytes   int `json:"rx-bytes"`
			RxDropped int `json:"rx-dropped"`
			RxErrs    int `json:"rx-errs"`
			RxPackets int `json:"rx-packets"`
			TxBytes   int `json:"tx-bytes"`
			TxDropped int `json:"tx-dropped"`
			TxErrs    int `json:"tx-errs"`
			TxPackets int `json:"tx-packets"`
		} `json:"statistics"`
	} `json:"result"`
}

// GetPveInterface 获取PVE客户机网卡地址
// 传入虚拟机vmID，返回ipv4, ipv6地址
func GetPveInterface(vmID string) (ipv4PveInterfaces []NetInterface, ipv6PveInterfaces []NetInterface, err error) {

	//vmID := "100"

	data := PveJson{}

	_, err = exec.LookPath("pvesh")
	if (runtime.GOOS != "linux") || (err != nil) {
		fmt.Println("PVE模式不支持在非Proxmox Virtual Environment系统运行！")
		return ipv4PveInterfaces, ipv6PveInterfaces, err
	}

	pve := exec.Command(
		"pvesh", "get",
		"nodes/pve/qemu/"+vmID+"/agent/network-get-interfaces",
		"--output-format",
		"json",
	)
	var bufferOut, bufferErr bytes.Buffer
	pve.Stdout = &bufferOut
	pve.Stderr = &bufferErr
	err = pve.Run()
	if err != nil {
		fmt.Printf("pvesh运行错误:%v\n", err)
	}

	err = json.Unmarshal(bufferOut.Bytes(), &data)
	if err != nil {
		return ipv4PveInterfaces, ipv6PveInterfaces, err
	}

	for i := 0; i < len(data.Result); i++ {
		var ipv4 []string
		var ipv6 []string

		for j := 0; j < len(data.Result[i].IPAddresses); j++ {
			if data.Result[i].IPAddresses[j].IPAddressType == "ipv4" {
				ipv4 = append(ipv4, data.Result[i].IPAddresses[j].IPAddress)
			}
			if data.Result[i].IPAddresses[j].IPAddressType == "ipv6" {
				ipv6 = append(ipv6, data.Result[i].IPAddresses[j].IPAddress)
			}
		}

		if len(ipv4) != 0 {
			ipv4PveInterfaces = append(
				ipv4PveInterfaces,
				NetInterface{
					Name:    data.Result[i].Name,
					Address: ipv4,
				},
			)
		}

		if len(ipv6) != 0 {
			ipv6PveInterfaces = append(
				ipv6PveInterfaces,
				NetInterface{
					Name:    data.Result[i].Name,
					Address: ipv6,
				},
			)
		}

	}

	//for i := 0; i < len(data.Result); i++ {
	//	println(data.Result[i].Name + ":")
	//	for j := 0; j < len(data.Result[i].IPAddresses); j++ {
	//		if data.Result[i].IPAddresses[j].IPAddressType == "ipv4" {
	//			print("ipv4\t")
	//			println(data.Result[i].IPAddresses[j].IPAddress)
	//		} else if data.Result[i].IPAddresses[j].IPAddressType == "ipv6" {
	//			print("ipv6\t")
	//			println(data.Result[i].IPAddresses[j].IPAddress)
	//		}
	//	}
	//}

	return ipv4PveInterfaces, ipv6PveInterfaces, nil
}
