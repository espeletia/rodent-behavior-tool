package config

type ServiceConfig struct {
	ID          string
	Environment string
	Name        string
	Deployment  string
}

func loadServiceConfig(serviceName string) ServiceConfig {
	serviceConfig := &ServiceConfig{}
	v := configViper("service", serviceName)
	err := v.BindEnv("Name", "SERVICE_NAME")
	if err != nil {
		panic(err)
	}
	err = v.BindEnv("Deployment", "SERVICE_DEPLOYMENT")
	if err != nil {
		panic(err)
	}

	err = v.BindEnv("Environment")
	if err != nil {
		panic(err)
	}
	err = v.BindEnv("ID")
	if err != nil {
		panic(err)
	}
	err = v.ReadInConfig()
	if err != nil {
		panic(err)
	}
	err = v.Unmarshal(serviceConfig)
	if err != nil {
		panic(err)
	}
	return *serviceConfig
}
