package config

import (
	"io/ioutil"
	"os"
	"path"
	"strconv"

	"github.com/ca17/teamsedge/common"
	"gopkg.in/yaml.v2"
)

// SysConfig 系统配置
type SysConfig struct {
	Appid    string `yaml:"appid"`
	Location string `yaml:"location"`
	Workdir  string `yaml:"workdir"`
	Version  string `yaml:"version"`
	Debug    bool   `yaml:"debug"`
}

// TeamsacsConfig TeamsACS 配置
type TeamsacsConfig struct {
	NbiUrl    string `yaml:"nbi_url"`    // Teamsacs 北向接口地址
	PgResturl string `yaml:"pg_resturl"` // T数据库 rest接口地址
	Secret    string `yaml:"secret"`     // JWT 密钥
	SubAddr   string `yaml:"sub_addr"`   // 订阅地址
	PubAddr   string `yaml:"pub_addr"`   // 发布地址
	Debug     bool   `yaml:"debug"`
}



type AppConfig struct {
	System   SysConfig      `yaml:"system"`
	Teamsacs TeamsacsConfig `yaml:"teamsacs"`
}

func (c *AppConfig) GetLogDir() string {
	return path.Join(c.System.Workdir, "logs")
}

func (c *AppConfig) GetDataDir() string {
	return path.Join(c.System.Workdir, "data")
}

func (c *AppConfig) InitDirs() {
	os.MkdirAll(path.Join(c.System.Workdir, "logs"), 0755)
	os.MkdirAll(path.Join(c.System.Workdir, "data"), 0755)
}

func setEnvValue(name string, val *string) {
	var evalue = os.Getenv(name)
	if evalue != "" {
		*val = evalue
	}
}

func setEnvBoolValue(name string, val *bool) {
	var evalue = os.Getenv(name)
	if evalue != "" {
		*val = evalue == "true" || evalue == "1" || evalue == "on"
	}
}

func setEnvIntValue(name string, val *int) {
	var evalue = os.Getenv(name)
	if evalue == "" {
		return
	}

	p, err := strconv.ParseInt(evalue, 10, 64)
	if err == nil {
		*val = int(p)
	}
}

func setEnvInt64Value(name string, val *int64) {
	var evalue = os.Getenv(name)
	if evalue == "" {
		return
	}

	p, err := strconv.ParseInt(evalue, 10, 64)
	if err == nil {
		*val = p
	}
}

var DefaultAppConfig = &AppConfig{
	System: SysConfig{
		Appid:    "TeamsEdge",
		Location: "Asia/Shanghai",
		Workdir:  "/var/teamsedge",
		Version:  "latest",
		Debug:    true,
	},
	Teamsacs: TeamsacsConfig{
		NbiUrl:    "http://127.0.0.1:1879",
		PgResturl: "http://127.0.0.1:3080",
		Secret:    "9b6de5cc-0731-4bf1-zpms-0f568ac9da37",
		PubAddr:   "tcp://127.0.0.1:1935",
		SubAddr:   "tcp://127.0.0.1:1936",
		Debug:     false,
	},
}

func LoadConfig(cfile string) *AppConfig {
	if cfile == "" {
		cfile = "teamsedge.yml"
	}
	if !common.FileExists(cfile) {
		cfile = "/etc/teamsedge.yml"
	}
	cfg := new(AppConfig)
	if common.FileExists(cfile) {
		data := common.Must2(ioutil.ReadFile(cfile))
		common.Must(yaml.Unmarshal(data.([]byte), cfg))
	} else {
		cfg = DefaultAppConfig
	}
	// 系统配置
	setEnvValue("TEAMSEDGE_SYSTEM_WORKER_DIR", &cfg.System.Workdir)
	setEnvBoolValue("TEAMSEDGE_SYSTEM_DEBUG", &cfg.System.Debug)
	setEnvValue("TEAMSACS_NBIURL", &cfg.Teamsacs.NbiUrl)
	setEnvValue("TEAMSACS_PG_RESTURL", &cfg.Teamsacs.PgResturl)
	setEnvValue("TEAMSACS_SECRET", &cfg.Teamsacs.Secret)

	// 订阅发布配置
	setEnvValue("TEAMSACS_PUB_ADDR", &cfg.Teamsacs.PubAddr)
	setEnvValue("TEAMSACS_SUB_ADDR", &cfg.Teamsacs.SubAddr)

	return cfg
}
