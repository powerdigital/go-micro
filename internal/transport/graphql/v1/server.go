package graphqlv1

import (
	"net/http"

	"github.com/99designs/gqlgen/graphql"
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/extension"
	"github.com/99designs/gqlgen/graphql/handler/transport"
)

//nolint:exhaustruct
func NewGqlServerHandler(gqlHandler GqlHandler) http.Handler {
	h := handler.New(NewExecutableSchema(Config{
		Resolvers: &Resolver{gqlHandler},
	}))

	h.AddTransport(transport.GET{})
	h.AddTransport(transport.POST{})
	h.AddTransport(transport.MultipartForm{})
	h.AddTransport(transport.Options{})

	h.SetRecoverFunc(graphql.DefaultRecover)

	h.Use(extension.Introspection{})

	return h
}
