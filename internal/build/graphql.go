package build

import (
	"context"

	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/gorilla/handlers"

	graphqlv1 "github.com/powerdigital/go-micro/internal/transport/graphql/v1"
)

func (b *Builder) SetGqlHandlers(ctx context.Context) error {
	router := b.httpRouter()
	router.Use(handlers.CORS(
		handlers.AllowedHeaders([]string{
			"content-type",
			"Access-Control-Request-Headers",
			"X-Requested-With",
		}),
		handlers.AllowedOrigins([]string{"*"}),
		handlers.AllowCredentials(),
		handlers.AllowedMethods([]string{"OPTIONS", "POST", "GET"}),
	))

	service, err := b.UserService(ctx)
	if err != nil {
		return err
	}

	handler := graphqlv1.NewGqlHandler(service)
	router.Handle("/query", handler)
	router.Handle("/playground", playground.Handler("Playground", "/query"))

	return nil
}
