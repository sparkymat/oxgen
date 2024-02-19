package template

import (
	"context"

	"github.com/samber/lo"
)

func (s *Service) Generate(ctx context.Context, cfg Config, cmd string, args map[string]string) (map[string][]byte, error) {
	cmd, found := lo.Find(cfg.Commands, func(c Command) bool { return c.Name == cmd })
}
