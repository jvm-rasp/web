package model

import (
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type RaspConfig struct {
	gorm.Model
	RowGuid    string `gorm:"type:varchar(255);comment:'唯一标识guid'" json:"rowGuid"`
	Name       string `gorm:"type:varchar(255);comment:'策略名称'" json:"name"`
	Version    int    `gorm:"type:int;comment:'版本号'" json:"version"`
	Desc       string `gorm:"type:varchar(255);comment:'策略说明'" json:"desc"`
	Status     bool   `gorm:"type:tinyint(1);default:1;comment:'是否禁用(0,禁用/1,启用, 默认启用)'" json:"status"`
	Creator    string `gorm:"type:varchar(20);comment:'创建人'" json:"creator"`
	Operator   string `gorm:"type:varchar(20);comment:'操作人'" json:"operator"`
	IsDefault  bool   `gorm:"type:tinyint(1);default:0;comment:'默认策略'" json:"isDefault"`
	IsModified bool   `gorm:"type:tinyint(1);default:0;comment:'是否已修改'" json:"isModified"`
}

type RaspFinalConfig struct {
	AgentMode        string         `json:"agentMode"`
	ConfigId         uint           `json:"configId"`
	ModuleAutoUpdate bool           `json:"moduleAutoUpdate"`
	LogPath          string         `json:"logPath"`
	RemoteHosts      string         `json:"remoteHosts"`
	EnableMdns       bool           `json:"enableMdns"`
	AgentConfigs     AgentConfig    `json:"agentConfigs"`
	RaspLibConfigs   ZipFileInfo    `json:"raspLibConfigs"`
	RaspBinConfigs   ZipFileInfo    `json:"raspBinConfigs"`
	ModuleConfigs    []ModuleConfig `json:"moduleConfigs"`
}

type RaspExportConfig struct {
	RaspFileInfo          []RaspFile        `json:"raspFileInfo"`
	RaspModuleInfo        []RaspModule      `json:"raspModuleInfo"`
	RaspComponentInfo     []RaspComponent   `json:"raspComponentInfo"`
	RaspConfigInfo        RaspConfig        `json:"raspConfigInfo"`
	RaspConfigHistoryInfo RaspConfigHistory `json:"raspConfigHistoryInfo"`
}

type AgentConfig struct {
	CheckDisable     bool   `json:"check_disable"`
	RedirectUrl      string `json:"redirect_url"`
	BlockStatusCode  int    `json:"block_status_code"`
	JsonBlockContent string `json:"json_block_content"`
	XmlBlockContent  string `json:"xml_block_content"`
	HtmlBlockContent string `json:"html_block_content"`
}

type ModuleConfig struct {
	ModuleName    string                 `json:"moduleName"`
	ModuleVersion string                 `json:"moduleVersion"`
	DownLoadUrl   string                 `json:"downLoadURL"`
	Parameters    map[string]interface{} `json:"parameters"`
	Md5           string                 `json:"md5"`
}

type ComponentConfig struct {
	ComponentName    string         `json:"componentName"`
	ComponentType    uint           `json:"componentType"`
	ComponentVersion string         `json:"componentVersion"`
	DownLoadURL      string         `json:"downLoadURL"`
	Md5              string         `json:"md5"`
	Parameters       datatypes.JSON `json:"parameters"`
}

type ZipFileInfo struct {
	FileName    string        `json:"fileName"`
	DownloadUrl string        `json:"downloadUrl"`
	Md5         string        `json:"md5"`
	ItemsInfo   []ZipItemInfo `json:"itemsInfo"`
}

type ZipItemInfo struct {
	FileName string `json:"fileName"`
	Md5      string `json:"md5"`
}
