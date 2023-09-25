package api

import (
	"errors"
	db "go-k8s/internal/db/sqlc"
	"go-k8s/internal/workers"
	"net/http"
	"time"

	"github.com/ShadrackAdwera/go-utils/utils"
	"github.com/gin-gonic/gin"
	"github.com/hibiken/asynq"
)

type createTenantArgs struct {
	Username   string `json:"username" binding:"required,alphanum"`
	Email      string `json:"email" binding:"required,email"`
	Password   string `json:"password" binding:"required,min=6"`
	Logo       string `json:"logo"`
	TenantName string `json:"tenant_name" binding:"required,min=6"`
}

func (srv *Server) createTenantTx(ctx *gin.Context) {
	var createTenantArgs createTenantArgs

	if err := ctx.ShouldBindJSON(&createTenantArgs); err != nil {
		ctx.JSON(http.StatusBadRequest, errJSON(err))
		return
	}

	hashPw, _ := utils.HashPassword(createTenantArgs.Password)

	response, err := srv.store.CreateTenantTx(ctx, db.CreateTenantInput{
		Username:   createTenantArgs.Username,
		Email:      createTenantArgs.Email,
		TenantName: createTenantArgs.TenantName,
		Logo:       createTenantArgs.Logo,
		Password:   hashPw,
	})

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errJSON(err))
		return
	}

	ctx.JSON(http.StatusCreated, response)
}

type userResponse struct {
	Username          string    `json:"username"`
	Email             string    `json:"email"`
	PasswordChangedAt time.Time `json:"password_changed_at"`
	CreatedAt         time.Time `json:"created_at"`
}

type loginUserRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

type loginUserResponse struct {
	// SessionID             uuid.UUID    `json:"session_id"`
	AccessToken           string       `json:"access_token"`
	AccessTokenExpiresAt  time.Time    `json:"access_token_expires_at"`
	RefreshToken          string       `json:"refresh_token"`
	RefreshTokenExpiresAt time.Time    `json:"refresh_token_expires_at"`
	User                  userResponse `json:"user"`
}

func (server *Server) login(ctx *gin.Context) {
	var req loginUserRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errJSON(err))
		return
	}

	user, err := server.store.GetUserByEmail(ctx, req.Email)
	if err != nil {
		if errors.Is(err, db.ErrRecordNotFound) {
			ctx.JSON(http.StatusNotFound, errJSON(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errJSON(err))
		return
	}

	err = utils.IsPassword(req.Password, user.Password)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, errJSON(err))
		return
	}

	accessToken, accessPayload, err := server.pasetoMaker.CreateToken(
		user.Username,
		time.Hour,
	)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errJSON(err))
		return
	}

	refreshToken, refreshPayload, err := server.pasetoMaker.CreateToken(
		user.Username,
		time.Hour*24*7,
	)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errJSON(err))
		return
	}

	// session, err := server.store.CreateSession(ctx, db.CreateSessionParams{
	// 	ID:           refreshPayload.ID,
	// 	Username:     user.Username,
	// 	RefreshToken: refreshToken,
	// 	UserAgent:    ctx.Request.UserAgent(),
	// 	ClientIp:     ctx.ClientIP(),
	// 	IsBlocked:    false,
	// 	ExpiresAt:    refreshPayload.ExpiredAt,
	// })
	// if err != nil {
	// 	ctx.JSON(http.StatusInternalServerError, errJSON(err))
	// 	return
	// }

	opts := []asynq.Option{
		asynq.MaxRetry(5),
		asynq.ProcessIn(5 * time.Second),
		asynq.Queue(workers.QueueDefault),
	}

	err = server.distro.DistributeTaskCreateLoginLog(ctx, &workers.PayloadCreateLoginLog{
		Email: user.Email,
	}, opts...)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errJSON(err))
		return
	}

	rsp := loginUserResponse{
		AccessToken:           accessToken,
		AccessTokenExpiresAt:  accessPayload.ExpiredAt,
		RefreshToken:          refreshToken,
		RefreshTokenExpiresAt: refreshPayload.ExpiredAt,
		User: userResponse{
			Username:          user.Username,
			Email:             user.Email,
			PasswordChangedAt: user.PasswordChangedAt,
			CreatedAt:         user.CreatedAt,
		},
	}
	ctx.JSON(http.StatusOK, rsp)
}
