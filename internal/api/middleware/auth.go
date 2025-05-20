package middleware

import (
	"net/http"

	log "github.com/sirupsen/logrus"

	"github.com/AlexJudin/wallet_java_code/internal/api/common"
	"github.com/AlexJudin/wallet_java_code/internal/service"
)

type AuthMiddleware struct {
	authService service.AuthService
}

func NewAuthMiddleware(authService service.AuthService) *AuthMiddleware {
	return &AuthMiddleware{
		authService: authService,
	}
}

func (a *AuthMiddleware) CheckToken(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		accessToken, err := r.Cookie("accessToken")
		if err != nil {
			common.ApiError(http.StatusUnauthorized, "access token not found", w)
			return
		}

		userLogin, err := a.authService.VerifyUser(accessToken.Value)
		if err != nil {
			common.ApiError(http.StatusUnauthorized, err.Error(), w)
			return
		}

		log.Infof("Пользователь %s сделал запрос %s", userLogin, r.URL.Path)

		next.ServeHTTP(w, r)
	}

	return http.HandlerFunc(fn)
}
