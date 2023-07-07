package nacos

import (
	"fmt"

	"github.com/nacos-group/nacos-sdk-go/model"
	"github.com/nacos-group/nacos-sdk-go/vo"
	"github.com/pkg/errors"
	"github.com/star-table/usercenter/core/conf"
	"github.com/star-table/usercenter/core/consts"
	"github.com/star-table/usercenter/core/logger"
	"github.com/star-table/usercenter/pkg/nacos"
	"github.com/star-table/usercenter/pkg/util/http"
	"github.com/star-table/usercenter/pkg/util/slice"
)

var (
	nacosClient *nacos.NacosClient
)

func init() {
	env := conf.GetEnv()
	if ok, _ := slice.Contain([]string{consts.RunEnvBjxLocal, consts.RunEnvBjxTest, consts.RunEnvBjxProd}, env); ok {
		return
	}
	client, err := nacos.NewNacosClient(conf.Cfg.Nacos)
	if err != nil {
		panic(err)
	}
	nacosClient = client
	logger.InitNacosLog()
}

func RegisterInstance(param vo.RegisterInstanceParam) (bool, error) {
	return nacosClient.RegisterInstance(param)
}

func DeregisterInstance(param vo.DeregisterInstanceParam) (bool, error) {
	return nacosClient.DeregisterInstance(param)
}

func GetService(param vo.GetServiceParam) (model.Service, error) {
	return nacosClient.GetService(param)
}

func GetAllServicesInfo(param vo.GetAllServiceInfoParam) (model.ServiceList, error) {
	return nacosClient.GetAllServicesInfo(param)
}

func SelectAllInstances(param vo.SelectAllInstancesParam) ([]model.Instance, error) {
	return nacosClient.SelectAllInstances(param)
}

func SelectInstances(param vo.SelectInstancesParam) ([]model.Instance, error) {
	return nacosClient.SelectInstances(param)
}

func SelectOneHealthyInstance(param vo.SelectOneHealthInstanceParam) (*model.Instance, error) {
	return nacosClient.SelectOneHealthyInstance(param)
}

func Subscribe(param *vo.SubscribeParam) error {
	return nacosClient.Subscribe(param)
}

func Unsubscribe(param *vo.SubscribeParam) error {
	return nacosClient.Unsubscribe(param)
}

func DoGet(serviceName, api string, params map[string]interface{}, headerOptions ...http.HeaderOption) (string, int, error) {
	instance, err := nacosClient.SelectOneHealthyInstance(vo.SelectOneHealthInstanceParam{
		ServiceName: serviceName,
	})
	if err != nil {
		return "", 501, err
	}
	if instance == nil {
		return "", 501, errors.New(fmt.Sprintf("service [%s] not found one healthy instance! "))
	}
	url := fmt.Sprintf("http://%s:%d/%s", instance.Ip, instance.Port, api)
	return http.Get(url, params, headerOptions...)
}

func DoPost(serviceName, api string, params map[string]interface{}, body string, headerOptions ...http.HeaderOption) (string, int, error) {
	instance, err := nacosClient.SelectOneHealthyInstance(vo.SelectOneHealthInstanceParam{
		ServiceName: serviceName,
	})
	if err != nil {
		return "", 500, err
	}
	if instance == nil {
		return "", 500, errors.New(fmt.Sprintf("service [%s] not found one healthy instance! "))
	}
	url := fmt.Sprintf("http://%s:%d/%s", instance.Ip, instance.Port, api)
	return http.Post(url, params, body, headerOptions...)
}
