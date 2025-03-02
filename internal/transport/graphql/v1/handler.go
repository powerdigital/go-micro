package graphqlv1

import (
	"net/http"

	"github.com/99designs/gqlgen/graphql"
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/extension"
	"github.com/99designs/gqlgen/graphql/handler/transport"

	userservice "github.com/powerdigital/go-micro/internal/service/v1/user"
	graphqlv1 "github.com/powerdigital/go-micro/pkg/graphql/v1"
)

//nolint:exhaustruct
func NewGqlHandler(service userservice.UserSrv) http.Handler {
	h := handler.New(graphqlv1.NewExecutableSchema(graphqlv1.Config{
		Resolvers: graphqlv1.NewResolver(service),
	}))

	h.AddTransport(transport.GET{})
	h.AddTransport(transport.POST{})
	h.AddTransport(transport.MultipartForm{})
	h.AddTransport(transport.Options{})

	h.SetRecoverFunc(graphql.DefaultRecover)

	h.Use(extension.Introspection{})

	return h
}
