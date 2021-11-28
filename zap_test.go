package lax_test

import (
	"github.com/jackc/pgx/v4"
	"go.ectobit.com/lax"
)

var (
	_ lax.Logger = (*lax.ZapAdapter)(nil)
	_ pgx.Logger = (*lax.ZapAdapter)(nil)
)
