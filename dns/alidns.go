package dns

import (
    "fmt"
    "github.com/aliyun/alibaba-cloud-sdk-go/services/alidns"
)

// Alidns 阿里云dns实现
type Alidns struct{}

func (alidns *Alidns) addRecord() (ipv4 bool, ipv6 bool) {
	client, err := alidns.NewClientWithAccessKey("cn-hangzhou", "<accessKeyId>", "<accessSecret>")

	request := alidns.CreateAddDomainRecordRequest()
	request.Scheme = "https"

	request.Value = "3.0.3.0"
	request.Type = "A"
	request.RR = "apitest1"
	request.DomainName = "dns-example.com"

	response, err := client.AddDomainRecord(request)
	if err != nil {
		fmt.Print(err.Error())
	}
	fmt.Printf("response is %#v\n", response)
}
