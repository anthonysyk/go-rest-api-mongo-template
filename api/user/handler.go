package user

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"io"
	"log"
	"net/http"
)

type Handler struct {
	service *Service
}

func NewHandler(userService *Service) *Handler {
	return &Handler{service: userService}
}

func message(m string) gin.H {
	return gin.H{"message": m}
}

func (h Handler) Login(secret string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		req := struct {
			ID       string `json:"id"`
			Password string `json:"password"`
		}{}
		err := ctx.BindJSON(&req)
		if err != nil {
			ctx.AbortWithStatus(http.StatusBadRequest)
			return
		}
		if req.ID == "" || req.Password == "" {
			ctx.AbortWithStatus(http.StatusBadRequest)
			return
		}
		token, err := h.service.Login(ctx, req.ID, req.Password, secret)
		if err != nil {
			log.Print(err.Error())
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		ctx.JSON(http.StatusOK, gin.H{"token": token})
	}
}

func (h Handler) CreateUsers(ctx *gin.Context) {
	file, err := ctx.FormFile("users")
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, message("could not get form file"))
		return
	}
	f, err := file.Open()
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, message("could not open file"))
		return
	}
	defer f.Close()

	bytes, err := io.ReadAll(f)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, message("could not open file"))
		return
	}

	var users []*User
	err = json.Unmarshal(bytes, &users)
	if err != nil {
		log.Print(err.Error())
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, message("could not parse json file"))
		return
	}

	inserted := h.service.CreateUsers(ctx, users)
	ctx.JSON(http.StatusOK, gin.H{"inserted": inserted})
}

func (h Handler) DeleteUser(ctx *gin.Context) {
	id := ctx.Param("id")
	if id == "" {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, message("id not specified"))
		return
	}
	err := h.service.DeleteUser(ctx, id)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, message("could not delete user"))
		return
	}

	ctx.Status(http.StatusOK)
}

func (h Handler) GetUser(ctx *gin.Context) {
	id := ctx.Param("id")
	if id == "" {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, message("id not specified"))
		return
	}
	user, err := h.service.GetUser(ctx, id)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusNotFound, message("user not found"))
		return
	}
	ctx.JSON(http.StatusOK, user)
}

func (h Handler) ListUsers(ctx *gin.Context) {
	users, err := h.service.ListUsers(ctx)
	if err != nil {
		log.Print(err)
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, message("could not list users"))
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"users": users})
}

func (h Handler) UpdateUser(ctx *gin.Context) {
	id := ctx.Param("id")
	if id == "" {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, message("id not specified"))
		return
	}

	var user User
	err := ctx.BindJSON(&user)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, message("could not parse user"))
		return
	}

	err = h.service.UpdateUser(ctx, id, &user)
	if err != nil {
		log.Print(err.Error())
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, message("could not update user"))
		return
	}

	ctx.Status(http.StatusOK)
}
