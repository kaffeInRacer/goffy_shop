package route

import (
	"context"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/kaffein/goffy/internal/infrastructure/persistence/db_postgre"
	"github.com/kaffein/goffy/internal/interfaces/http/dto"
	"github.com/kaffein/goffy/pkg/adapter/postgre"
	"github.com/lib/pq"
	"github.com/rs/zerolog/log"
)

func SetupRoutes(router *gin.Engine) {
	db := postgre.GetDB()

	// dto
	r := dto.AuthRegister{
		Username: "secret 123",
		Password: "secret",
		Email:    "secret@mail.com",
		Gender:   "male",
	}

	data := r.ParsingData()
	data.BeforeSave()

	userRepo := db_postgre.NewUserRepository(db)
	data, err := userRepo.Create(context.Background(), data)
	if err != nil {
		// Check if the error is a PostgreSQL error
		var pgErr *pq.Error
		if errors.As(err, &pgErr) {
			log.Printf("PostgreSQL error: %s (code: %s)", pgErr.Message, pgErr.Code)
		} else {
			log.Printf("Unknown error: %v", err)
		}
		return
	}
	log.Info().Msg("OK")
}
