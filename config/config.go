package config

import (
	"DetectiveMasterServer/util"
	"encoding/json"
	"io/ioutil"
	"os"
	"sync"
	"time"
)

const CONFIG_CATEGORY = "Config"
const CONFIG_PATH = "config.json"

var (
	config     *Config
	configLock = new(sync.RWMutex)
)

type Config struct {
	WebSocketHost           string   `json:"web_socket_host"`
	WebSocketPath           string   `json:"web_socket_path"`
	ServerOpenId            string   `json:"server_open_id"`
	ServerHost              string   `json:"server_host"`
	AppId                   string   `json:"app_id"`
	AppSecret               string   `json:"app_secret"`
	DBTimeout               int      `json:"db_timeout"`
	DBAddr                  string   `json:"db_addr"`
	DBHost                  string   `json:"db_host"`
	MaxQueue                int      `json:"max_queue"`
	Limit                   int      `json:"limit"`
	TokenTimeout            int      `json:"token_timeout"`
	ImSdkAppId              int      `json:"im_sdk_appid"`
	ImKey                   string   `json:"im_key"`
	ImIdent                 string   `json:"im_identifier"`
	NotifyUrl               string   `json:"notify_url"`
	MchId                   string   `json:"mch_id"`
	MchKey                  string   `json:"mch_key"`
	PayIp                   string   `json:"pay_ip"`
	RedisAddr               string   `json:"redis_addr"`
	RedisAuth               string   `json:"redis_auth"`
	RedisDb                 int      `json:"redis_db"`
	RedisExp                int      `json:"redis_exp"`
	RedisMaxIdle            int      `json:"redis_max_idle"`
	RedisMaxActive          int      `json:"redis_max_active"`
	RedisIdleTimeout        int      `json:"redis_idle_timeout"`
	RedisDialConnectTimeout int      `json:"redis_dial_connect_timeout"`
	RedisDialReadTimeout    int      `json:"redis_dial_read_timeout"`
	RedisDialWriteTimeout   int      `json:"redis_dial_write_timeout"`
	RedisRecvChannels       []string `json:"redis_recv_channels"`
	RedisSendChannels       []string `json:"redis_send_channels"`
	RedisDataDb             int      `json:"redis_data_db"`
	Version                 int      `json:"version"`
}

// Func: Load Config
func LoadConfig() bool {
	f, err := ioutil.ReadFile(CONFIG_PATH)
	if err != nil {
		util.Error("Load Config Err: %v" + err.Error())
		return false
	}

	temp := new(Config)
	err = json.Unmarshal(f, &temp)
	if err != nil {
		util.Error("Parse Config Err: %v", err.Error())
		return false
	}

	configLock.Lock()
	temp.RedisExp = temp.RedisExp * int(time.Hour.Seconds()) //Redis中key过期时间
	config = temp
	configLock.Unlock()

	return true
}

// Func: Get Config
func GetConfig() *Config {
	configLock.RLock()
	defer configLock.RUnlock()
	return config
}

func init() {
	if !LoadConfig() {
		os.Exit(1)
	}
}
