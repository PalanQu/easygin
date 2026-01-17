package db

import (
	"context"
	"database/sql"
	"easygin/pkg/apperror"
	"easygin/pkg/config"
	"easygin/pkg/ent"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/mattn/go-sqlite3"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"

	entsql "entgo.io/ent/dialect/sql"
)

func CreateDBClient(config *config.Config, tracer trace.Tracer) (*ent.Client, error) {
	db, err := sql.Open(config.Database.Driver, config.Database.DSN)
	if err != nil {
		return nil, apperror.InternalError("failed to open database connection", err)
	}
	drv := entsql.OpenDB(config.Database.Driver, db)

	client := ent.NewClient(ent.Driver(drv))
	client.Use(func(next ent.Mutator) ent.Mutator {
		return ent.MutateFunc(func(ctx context.Context, m ent.Mutation) (ent.Value, error) {
			opType := m.Op().String()
			tablesName := m.Type()
			spanName := fmt.Sprintf("%s-%s", tablesName, opType)
			newCtx, span := tracer.Start(ctx, spanName)
			span.SetAttributes(
				attribute.String("table", tablesName),
				attribute.String("operation", opType),
			)
			defer span.End()
			return next.Mutate(newCtx, m)
		})
	})
	return client, nil
}
