package conf

import (
	"bytes"
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/micro/go-micro/v2/config"
	"github.com/spf13/viper"
	"github.com/star-table/usercenter/core/consts"
	"github.com/star-table/usercenter/pkg/util/json"
	"github.com/star-table/usercenter/pkg/util/network"
	remote "github.com/yoyofxteam/nacos-viper-remote"
)

var (
	confPath string
	env      string
	host     string
	Cfg      Config = Config{
		Viper: viper.New(),
	}
)

type Config struct {
	Viper       *viper.Viper
	Server      *ServerConfig      `json:"server"`
	Application *ApplicationConfig `json:"application"`
	Redis       *RedisConfig       `json:"redis"`
	Mysql       *MysqlConfig       `json:"mysql"`
	Logger      *LoggerConfig      `json:"logger"`
	Trace       *TraceConfig       `json:"trace"`
	Etcd        *EtcdConfig        `json:"etcd"`
	Snowflake   *SnowflakeConfig   `json:"snowflake"`
	Kafka       *KafkaConfig       `json:"kafka"`
	Emitter     *EmitterConfig     `json:"emitter"`
	Sentry      *SentryConfig      `json:"sentry"`
	Jaeger      *JaegerConfig      `json:"jaeger"`
	Apis        *ApisConfig        `json:"apis"`
	Nacos       *NacosConfig       `json:"nacos"`
	Discovery   *DiscoveryConfig   `json:"discovery"`
	OSS         *OSSConfig         `json:"oss"`
	Resource    *ResourceConfig    `json:"resource"`
	MQTT        *MQTTConfig        `json:"mqtt"`
}

type ServerConfig struct {
	Port int    `json:"port"`
	Host string `json:"host"`
}

type ApplicationConfig struct {
	Name    string `json:"name"`
	Domain  string `json:"domain"`
	RunMode int    `json:"runMode"`
}

type SentryConfig struct {
	Dsn string `json:"dsn"`
}

type RedisConfig struct {
	Host           string `json:"host"`
	Port           int    `json:"port"`
	Pwd            string `json:"pwd"`
	Database       int    `json:"database"`
	MaxIdle        int    `json:"maxIdle"`
	MaxActive      int    `json:"maxActive"`
	MaxIdleTimeout int    `json:"maxIdleTimeout"`
	IsSentinel     bool   `json:"isSentinel"`
	MasterName     string `json:"masterName"`
	SentinelPwd    string `json:"sentinelPwd"`
}

type KafkaConfig struct {
	NameServers    string `json:"nameServers"`
	ReconsumeTimes int    `json:"reconsumeTimes"`
	RePushTimes    int    `json:"rePushTimes"`
}

type EmitterConfig struct {
	Address         string `json:"address"`
	SecretKey       string `json:"secretKey"`
	ConnectPoolSize int    `json:"connectPoolSize"`
	Enable          bool   `json:"enable"`
}

type EtcdConfig struct {
	Addrs []string `json:"addrs"`
}

type SnowflakeConfig struct {
	WorkerId int64 `json:"workerId"`
}

type MysqlConfig struct {
	Host         string `json:"host"`
	Port         int    `json:"port"`
	User         string `json:"user"`
	Pwd          string `json:"pwd"`
	Database     string `json:"database"`
	MaxOpenConns int    `json:"maxOpenConns"`
	MaxIdleConns int    `json:"maxIdleConns"`
	MaxLifetime  int    `json:"maxLifetime"`
}

type LoggerConfig struct {
	Level      string   `json:"level"`
	Outputs    []string `json:"outputs"`
	ErrOutputs []string `json:"errOutputs"`
	MaxSize    int      `json:"maxSize"`
	MaxNum     int      `json:"maxNum"`
	Console    bool     `json:"console"`
}

type TraceConfig struct {
	Zipkin ZipkinConfig `json:"zipkin"`
	Jaeger JaegerConfig `json:"jaeger"`
}

type ZipkinConfig struct {
	Address string `json:"address"`
}

type JaegerConfig struct {
	UdpAddress   string  `json:"udpAddress"`
	TraceService string  `json:"traceService"`
	SamplerType  string  `json:"samplerType"`
	SamplerParam float64 `json:"samplerParam"`
}

type ApisConfig struct {
	Wechaty string `json:"wechaty"`
}

//oss配置
type OSSConfig struct {
	BucketName      string
	EndPoint        string
	AccessKeyId     string
	AccessKeySecret string
	Policies        OSSPolicyConfig
}

// ResourceConfig
type ResourceConfig struct {
	RootPath    string `json:"root_path"`    // RootPath
	LocalDomain string `json:"local_domain"` // LocalDomain
}

type OSSPolicyConfig struct {
	ProjectCover    OSSPolicyInfo
	IssueResource   OSSPolicyInfo
	IssueInputFile  OSSPolicyInfo
	ProjectResource OSSPolicyInfo
	CompatTest      OSSPolicyInfo
	UserAvatar      OSSPolicyInfo
	Feedback        OSSPolicyInfo
	// 混合资源，可以用于一些非主要场景的文件上传配置
	MixResource OSSPolicyInfo
}

type OSSPolicyInfo struct {
	//容器
	BucketName string
	//有效期
	Expire int64
	//目录
	Dir string
	//最大文件大小
	MaxFileSize int64
	//回调地址
	CallbackUrl string
}

type MQTTConfig struct {
	//地址
	Address string `json:"address"`
	//Host
	Host string `json:"host"`
	//Port
	Port int `json:"port"`
	//key
	SecretKey string `json:"secretKey"`
	//channel
	Channel string `json:"channel"`
	//连接数
	ConnectPoolSize int `json:"connectPoolSize"`
	//是否开启
	Enable bool `json:"enable"`
}

type NacosConfig struct {
	Client NacosClientConfig            `json:"client"`
	Server map[string]NacosServerConfig `json:"server"`
}

type NacosClientConfig struct {
	TimeoutMs            uint64 `json:"timeout_ms"`              // TimeoutMs http请求超时时间，单位毫秒
	ListenInterval       uint64 `json:"listen_interval"`         // ListenInterval 监听间隔时间，单位毫秒（仅在ConfigClient中有效）
	BeatInterval         int64  `json:"beat_interval"`           // BeatInterval         心跳间隔时间，单位毫秒（仅在ServiceClient中有效）
	NamespaceId          string `json:"namespace_id"`            // NamespaceId          nacos命名空间
	Endpoint             string `json:"endpoint"`                // Endpoint             获取nacos节点ip的服务地址
	CacheDir             string `json:"cacheDir"`                // CacheDir             缓存目录
	LogDir               string `json:"logDir"`                  // LogDir               日志目录
	UpdateThreadNum      int    `json:"update_thread_num"`       // UpdateThreadNum      更新服务的线程数
	NotLoadCacheAtStart  bool   `json:"not_load_cache_at_start"` // NotLoadCacheAtStart  在启动时不读取本地缓存数据，true--不读取，false--读取
	UpdateCacheWhenEmpty bool   `json:"update_cache_when_empty"` // UpdateCacheWhenEmpty 当服务列表为空时是否更新本地缓存，true--更新,false--不更新
	Username             string `json:"username"`
	Password             string `json:"password"`
}

type NacosServerConfig struct {
	IpAddr      string `json:"ip_addr"`      // IpAddr      nacos命名空间
	ContextPath string `json:"context_path"` // ContextPath 获取nacos节点ip的服务地址
	Port        uint64 `json:"port"`         // Port        缓存目录
}

// DiscoveryConfig is nacos service config.
type DiscoveryConfig struct {
	GroupName   string  `json:"group_name"`
	ClusterName string  `json:"cluster_name"`
	Weight      float64 `json:"weight"`
	Enable      bool    `json:"enable"`
	Healthy     bool    `json:"healthy"`
	Ephemeral   bool    `json:"ephemeral"`
}

func GetEnv() string {
	return env
}

func GetHost() string {
	return host
}

func init() {
	defEvn := os.Getenv(consts.RunEnv)
	if defEvn == consts.RunEnvGray {
		err := LoadNacosConfigAutoConfiguration("usercenter")
		if err != nil {
			log.Fatal(err)
		}
	} else {
		LoadLocalConfig()
	}
	fmt.Printf("config %s\n", json.ToJsonIgnoreError(Cfg))
}

func LoadLocalConfig() {
	var (
		defEvn  = os.Getenv(consts.RunEnv)
		defHost = network.GetIntranetIp()
	)
	flag.StringVar(&env, "e", defEvn, "run env")
	flag.StringVar(&host, "h", defHost, "server host")
	flag.StringVar(&confPath, "c", "config/application.yaml", "run config path")
	flag.Parse()
	if env != "" {
		confPath = strings.Replace(confPath, ".yaml", "."+env+".yaml", -1)
	}
	// 本地环境的调试，或者本地调用（service 目录下的）测试用例
	// 使用测试用例时，可以将下面注释解开
	//if defEvn == "" {
	//	confPath = "../../" + "config/application.dev.yaml"
	//}
	if err := load(confPath); err != nil {
		panic(err)
	}
}

func load(path string) error {
	_, err := os.Stat(path)
	if err != nil {
		return err
	}
	err = config.LoadFile(confPath)
	if err != nil {
		return err
	}
	return config.Scan(&Cfg)
}

func LoadNacosConfigAutoConfiguration(applicationName string) error {
	host := os.Getenv("REGISTER_HOST")
	if host == "" {
		return errors.New("nacos host is empty")
	}
	portStr := os.Getenv("REGISTER_PORT")
	if portStr == "" {
		return errors.New("nacos port is empty")
	}
	port, err := strconv.ParseUint(portStr, 10, 64)
	if err != nil {
		return err
	}
	namespaceId := os.Getenv("REGISTER_NAMESPACE")
	if namespaceId == "" {
		return errors.New("nacos namespaceId is empty")
	}
	username := os.Getenv("REGISTER_USERNAME")
	password := os.Getenv("REGISTER_PASSWORD")
	return LoadNacosConfig(host, port, namespaceId, "DEFAULT_GROUP", "application."+applicationName+".yaml", username, password)
}

func LoadNacosConfig(host string, port uint64, namespaceId, group, dataId string, username, password string) error {
	var auth *remote.Auth
	if username != "" && password != "" {
		auth = &remote.Auth{
			Enable:   true,
			User:     username,
			Password: password,
		}
	}
	remote.SetOptions(&remote.Option{
		Url:         host,
		Port:        port,
		NamespaceId: namespaceId,
		GroupName:   group,
		Config:      remote.Config{DataId: dataId},
		Auth:        auth,
	})
	if Cfg.Viper == nil {
		Cfg.Viper = viper.New()
	}
	err := Cfg.Viper.AddRemoteProvider("nacos", host, "")
	if err != nil {
		log.Println(err)
		return err
	}
	Cfg.Viper.SetConfigType("yaml")
	err = Cfg.Viper.ReadRemoteConfig()
	if err != nil {
		log.Println(err)
		return err
	}
	remote.NewRemoteProvider("yaml")
	respChan, _ := viper.RemoteConfig.WatchChannel(remote.DefaultRemoteProvider())
	go func(rc <-chan *viper.RemoteResponse) {
		for {
			b := <-rc
			reader := bytes.NewReader(b.Value)
			err := Cfg.Viper.MergeConfig(reader)
			if err != nil {
				log.Println(err)
			}
		}
	}(respChan)
	if err := Cfg.Viper.Unmarshal(&Cfg); err != nil {
		log.Println(err)
		return err
	}
	return nil
}
