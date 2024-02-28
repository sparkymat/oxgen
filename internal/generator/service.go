package generator

func New(config Config) *Service {
	lookupTable := generateLookupTable(config)

	return &Service{
		Config:      config,
		LookupTable: lookupTable,
	}
}

type Service struct {
	Config      Config
	LookupTable map[string]string
}
