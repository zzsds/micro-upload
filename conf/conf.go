package conf

import (
	"encoding/json"
	"os"
	"path/filepath"
	"runtime"
	"sync"

	"flag"

	"github.com/BurntSushi/toml"
	"github.com/nacos-group/nacos-sdk-go/clients"
	"github.com/nacos-group/nacos-sdk-go/common/constant"
	"github.com/nacos-group/nacos-sdk-go/vo"
)

var (
	Conf       Config
	once       sync.Once
	ConfigPath string

	// ******* acm config ***********************
	CFG_ENDPOINT    = GetEnv("CFG_ENDPOINT", "acm.aliyun.com")
	CFG_NAMESPACEID = GetEnv("CFG_NAMESPACEID", "a23c93cd-491c-44dd-be30-fb1df6e6ddaf")
	CFG_ACCESSKEY   = GetEnv("CFG_ACCESSKEY", "LTAI4FgL4Ew4kGTSEWQ8gSbo")
	CFG_SECRETKEY   = GetEnv("CFG_SECRETKEY", "ZElyfnMQ4E4tE8QKJeXdZmgJ54Mgea")
	CFG_CLUSTER     = GetEnv("CFG_CLUSTER", "test")
)

const (
	DataIdDefault = "api.upload"
)

type Config struct {
	Aliyun *Aliyun
}

type Aliyun struct {
	Endpoint        string
	AccessKeyId     string
	AccessKeySecret string
	BucketName      string
	BucketHost      string
}

// 初始化配置
func InitConfig() {
	flag.StringVar(&ConfigPath, "c", GetFileConfFile(), "this default local conf.toml")
	flag.Parse()

	if isExists(ConfigPath) {
		local()
	} else {
		load()
	}
}

func CheckErr(err error) {
	if err != nil {
		panic(err)
	}
}

// 加载本地配置
func local() {
	once.Do(func() {
		_, err := toml.DecodeFile(ConfigPath, &Conf)
		CheckErr(err)
	})
}

// 加载线上配置
func load() {
	clientConfig := constant.ClientConfig{
		//
		Endpoint:       CFG_ENDPOINT + ":8080",
		NamespaceId:    CFG_NAMESPACEID,
		AccessKey:      CFG_ACCESSKEY,
		SecretKey:      CFG_SECRETKEY,
		TimeoutMs:      5 * 1000,
		ListenInterval: 30 * 1000,
	}
	configClient, err := clients.CreateConfigClient(map[string]interface{}{
		"clientConfig": clientConfig,
	})

	CheckErr(err)

	content, err := configClient.GetConfig(vo.ConfigParam{
		DataId: DataIdDefault,
		Group:  CFG_CLUSTER})
	CheckErr(err)
	CheckErr(json.Unmarshal([]byte(content), &Conf))
}

func GetEnv(key, value string) string {
	newValue := os.Getenv(key)
	if newValue == "" {
		return value
	}
	return newValue
}

func isExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil && !os.IsNotExist(err)
}

func GetFileConfFile() string {
	_, f, _, _ := runtime.Caller(1)
	return filepath.Join(filepath.Dir(f), "/config.toml")
}
