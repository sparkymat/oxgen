package git

func (s *Service) StatusClean() (bool, error) {
	workTree, err := s.repo.Worktree()
	if err != nil {
		return true, err //nolint:wrapcheck
	}

	status, err := workTree.Status()
	if err != nil {
		return true, err //nolint:wrapcheck
	}

	return status.IsClean(), nil
}
