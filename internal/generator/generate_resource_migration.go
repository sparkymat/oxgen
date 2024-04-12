package generator

import "context"

func (s *Service) generateResourceMigration(ctx context.Context, name string) error {
	if err := s.ensureFolderExists("migrations"); err != nil {
		return err
	}

	return nil
}
