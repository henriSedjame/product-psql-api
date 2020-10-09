package core

type AppProperties struct {
	Datasource DBProperties `json:"datasource"`
	Server ServerProperties `json:"server"`
	Cors CorsProperties `json:"cors"`
	ActiveProfiles string `json:"profiles"`
}

type DBProperties struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Dbname string `json:"database"`
	Host string `json:"host"`
	Port int `json:"port"`
}

type ServerProperties struct {
	Port int `json:"port"`
}

type CorsProperties struct {
	AllowedOrigins string `json:"allowed_origins"`
	AllowedHeaders string `json:"allowed_headers"`
	AllowedMethods string `json:"allowed_methods"`
}

func DefaultAppProperties() AppProperties  {

	return AppProperties{
		Server: ServerProperties{
			Port: 8080,
		},
		Datasource: DBProperties{
			Port: 5432,
			Host: "localhost",
		},
	}
}
