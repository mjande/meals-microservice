package utils

import (
	"context"
	"strconv"

	"github.com/go-chi/jwtauth/v5"
)

func ExtractUserIDFromContext(ctx context.Context) string {
	_, claims, _ := jwtauth.FromContext(ctx)
	userId := claims["user_id"].(float64)
	return strconv.Itoa(int(userId))
}
