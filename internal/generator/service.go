package generator

func New(config Config) *Service {
	return &Service{
		Config: config,
	}
}

type Service struct {
	Config Config
}
