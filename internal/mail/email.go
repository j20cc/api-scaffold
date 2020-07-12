package mail

import (
	"context"
	"errors"
	"fmt"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/dm"
	"github.com/lukedever/gvue-scaffold/app/models"
	"github.com/lukedever/gvue-scaffold/internal/helper"
	"github.com/lukedever/gvue-scaffold/internal/log"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"net/url"
	"time"
)

type SendMailType string

const (
	Verify SendMailType = "verify"
	Reset  SendMailType = "reset"
)

var keySignPrefix = "user:sign:%s:%s"

func SendEmail(to, sub, body string) error {
	accessKeyId := viper.GetString("email.id")
	accessKeySecret := viper.GetString("email.secret")
	if accessKeyId == "" || accessKeySecret == "" {
		log.Warn("access key id and secret is empty.")
		return nil
	}
	client, err := dm.NewClientWithAccessKey("cn-hangzhou", accessKeyId, accessKeySecret)
	if err != nil {
		return err
	}
	request := dm.CreateSingleSendMailRequest()
	request.Scheme = "https"

	request.AccountName = viper.GetString("email.from")
	request.AddressType = requests.NewInteger(1)
	request.ReplyToAddress = requests.NewBoolean(false)
	request.ToAddress = to
	request.Subject = sub
	request.FromAlias = viper.GetString("app.name")
	request.HtmlBody = body

	response, err := client.SingleSendMail(request)
	if err != nil {
		return err
	}
	fmt.Printf("response is %#v\n", response)
	return nil
}

func GetSignedUrl(email string, t SendMailType) (string, error) {
	host := viper.GetString("app.url")
	//签名
	signstr := helper.Md5(email)
	values := url.Values{}
	//redis key
	if err := saveSignKey(fmt.Sprintf(keySignPrefix, t, signstr), email); err != nil {
		log.Error("signed key save to redis error.", zap.Error(err))
		return "", err
	}
	switch t {
	case Reset:
		host += "/password/reset"
		values.Set("email", email)
	case Verify:
		host += "/verification"
	default:
		return "", errors.New("签名失败")
	}
	values.Set("sign", signstr)
	query := values.Encode()
	//构造url
	myurlstr := host + "?" + query
	myurl, err := url.Parse(myurlstr)
	if err != nil {
		log.Error("url build error.", zap.Error(err))
		return "", err
	}
	//解析url
	return myurl.String(), nil
}

func DecodeSignUrl(t SendMailType, sign string) (string, error) {
	if t != Reset && t != Verify {
		return "", errors.New("链接错误")
	}
	key := fmt.Sprintf(keySignPrefix, t, sign)
	return models.RedisCli.Get(context.Background(), key).Result()
}

func saveSignKey(key string, value string) error {
	//存到redis
	if _, err := models.RedisCli.Set(context.Background(), key, value, time.Minute*30).Result(); err != nil {
		return err
	}
	return nil
}
