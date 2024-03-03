package generator

func New(config Config) *Service {
	lookupTable := generateLookupTableForProject(config)

	return &Service{
		Config:      config,
		LookupTable: lookupTable,
	}
}

type Service struct {
	Config      Config
	LookupTable map[string]string
}
