package config

type InstallConfig struct {
	Namespace  string      `json:"namespace" yaml:"namespace"`
	Repo       *Repo       `json:"repo" yaml:"repo"`
	Host       *Host       `json:"host" yaml:"host"`
	Port       string      `json:"port" yaml:"port"`
	Middleware *Middleware `json:"middleware" yaml:"middleware"`
	Plugins    []string    `json:"plugins" yaml:"plugins"`
}

type Repo struct {
	Url  string `json:"url" yaml:"url"`
	Name string `json:"name" yaml:"name"`
}

type Host struct {
	Admin  string `json:"admin" yaml:"admin"`
	Tenant string `json:"tenant" yaml:"tenant"`
}

type Middleware struct {
	Queue           *Value `json:"queue" yaml:"queue"`
	Database        *Value `json:"database" yaml:"database"`
	Cache           *Value `json:"cache" yaml:"cache"`
	Search          *Value `json:"search" yaml:"search"`
	ServiceRegistry *Value `json:"service_registry" yaml:"service_registry"`
	//TSDB            *Value `json:"tsdb" yaml:"tsdb"`
}

type Value struct {
	Customized bool   `json:"customized" yaml:"customized"`
	Url        string `json:"url" yaml:"url"`
}

func (ic *InstallConfig) GetMiddleware() map[string]*Value {
	config := make(map[string]*Value)
	if ic.Middleware == nil {
		return config
	}
	config["queue"] = ic.Middleware.Queue
	config["database"] = ic.Middleware.Database
	config["cache"] = ic.Middleware.Cache
	config["search"] = ic.Middleware.Search
	config["service_registry"] = ic.Middleware.ServiceRegistry
	return config
}

func (ic *InstallConfig) SetMiddleware(middleware map[string]*Value) {
	ic.Middleware = &Middleware{
		Queue:           middleware["queue"],
		Database:        middleware["database"],
		Cache:           middleware["cache"],
		Search:          middleware["search"],
		ServiceRegistry: middleware["service_registry"],
	}
}
