package middleware

import (
	"encoding/json"
	"net/http"

	log "github.com/sirupsen/logrus"

	"github.com/AlexJudin/wallet_java_code/internal/api/entity"
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
			w.WriteHeader(http.StatusUnauthorized)
			apiError, _ := json.Marshal(entity.ApiError{Error: "access token not found"})
			w.Write(apiError)
			return
		}

		userLogin, err := a.authService.VerifyUser(accessToken.Value)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			apiError, _ := json.Marshal(entity.ApiError{Error: err.Error()})
			w.Write(apiError)
			return
		}

		log.Infof("Пользователь %s сделал запрос %s\n", userLogin, r.URL.Path)

		next.ServeHTTP(w, r)
	}

	return http.HandlerFunc(fn)
}
