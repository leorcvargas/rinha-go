package database

import (
	"github.com/leorcvargas/rinha-2023-q3/internal/app/infra/database/peopledb"
	"go.uber.org/fx"
)

var Module = fx.Options(
	fx.Provide(NewPostgresDatabase),
	peopledb.Module,
)
