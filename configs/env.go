package configs

type Env string

func (e Env) IsProd() bool {
	return e == "production"
}

func (e Env) IsDev() bool {
	return e == "development"
}

func (e Env) IsLocal() bool {
	return e == "local"
}

func (e Env) String() string {
	return string(e)
}
