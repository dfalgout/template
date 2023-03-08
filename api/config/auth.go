package config

import (
	"context"
	"net/http"
	"strings"

	"entgo.io/ent/entc/integration/multischema/ent"
)

type Auth struct {
	client *ent.Client
}

type contextKey struct {
	name string
}

var userCtxKey = &contextKey{"user"}

func NewAuth(client *ent.Client) *Auth {
	return &Auth{
		client,
	}
}

func UserFromContext(ctx context.Context) (*ent.User, bool) {
	raw, ok := ctx.Value(userCtxKey).(*ent.User)
	return raw, ok
}

func (m *Auth) Init() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			const authBearer = "Bearer"
			authHeader := r.Header.Get("Authorization")

			if !strings.HasPrefix(authHeader, authBearer) {
				next.ServeHTTP(w, r)
				return
			}

			token := authHeader[len(authBearer)+1:]
			if token == "" {
				next.ServeHTTP(w, r)
				return
			}

			// found, err := m.client.User.
			// 	Query().
			// 	Where(
			// 		user.ID(userInfo.PublicAddress),
			// 	).
			// 	Only(r.Context())

			// if err != nil {
			// 	found, err = m.client.User.
			// 		Create().
			// 		SetID(userInfo.PublicAddress).
			// 		SetEmail(userInfo.Email).
			// 		SetNickname(userInfo.Email).
			// 		Save(r.Context())
			// 	if err != nil {
			// 		next.ServeHTTP(w, r)
			// 		return
			// 	}
			// }

			// valuesCtx := context.WithValue(r.Context(), userCtxKey, found)
			// r = r.WithContext(valuesCtx)

			next.ServeHTTP(w, r)
		})
	}
}
