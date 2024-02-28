package git

import gitpkg "github.com/go-git/go-git/v5"

func New() (*Service, error) {
	repo, err := gitpkg.PlainOpen(".")
	if err != nil {
		return nil, err //nolint:wrapcheck
	}

	return &Service{
		repo: repo,
	}, nil
}

type Service struct {
	repo *gitpkg.Repository
}
