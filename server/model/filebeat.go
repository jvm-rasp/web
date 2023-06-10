package model

type FilebeatConfig struct {
	Filebeat   *Filebeat              `json:"filebeat" mapstructure:"filebeat" structs:"filebeat"`
	Logging    map[string]interface{} `json:"logging" mapstructure:"logging" structs:"logging"`
	Output     map[string]interface{} `json:"output" mapstructure:"output" structs:"output"`
	Processors []interface{}          `json:"processors" mapstructure:"processors" structs:"processors"`
}

type Filebeat struct {
	Config map[string]interface{} `json:"config" mapstructure:"config"  structs:"config"`
	Inputs []*Inputs              `json:"inputs" mapstructure:"inputs"  structs:"inputs"`
}

type Inputs struct {
	Fields  map[string]interface{} `json:"fields" mapstructure:"fields" structs:"fields"`
	Paths   []string               `json:"paths" mapstructure:"paths" structs:"paths"`
	Type    string                 `json:"type" mapstructure:"type" structs:"type"`
	Enabled bool                   `json:"enabled" mapstructure:"enabled" structs:"enabled"`
}
