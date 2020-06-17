package helper

import (
	"fmt"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/dm"
	"github.com/spf13/viper"
)

var (
	accessKeyId     = "LTAI4G2AQRG1jJ9C6Li1e959"
	accessKeySecret = "Wr3S5ZHXu9YJVDhS9NrYjw6L8Og6Kp"
)

func SendEmail(to, sub, body string) error {
	client, err := dm.NewClientWithAccessKey("cn-hangzhou", accessKeyId, accessKeySecret)
	if err != nil {
		return err
	}
	request := dm.CreateSingleSendMailRequest()
	request.Scheme = "https"

	request.AccountName = "mail@yi2a.com"
	request.AddressType = requests.NewInteger(1)
	request.ReplyToAddress = requests.NewBoolean(false)
	request.ToAddress = to
	request.Subject = sub
	request.FromAlias = viper.GetString("site.name")
	request.HtmlBody = body

	response, err := client.SingleSendMail(request)
	if err != nil {
		return err
	}
	fmt.Printf("response is %#v\n", response)
	return nil
}
