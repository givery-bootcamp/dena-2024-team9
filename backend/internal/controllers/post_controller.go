package controllers

import (
	"errors"
	"myapp/internal/usecases"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func GetPosts(ctx *gin.Context, usecase *usecases.ListPostUsecase) {
	result, err := usecase.Execute()
	if err != nil {
		handleError(ctx, http.StatusInternalServerError, err)
	} else if result != nil {
		ctx.JSON(http.StatusOK, result)
	} else {
		handleError(ctx, http.StatusNotFound, errors.New("not found"))
	}
}

func GetPostById(ctx *gin.Context, usecase *usecases.GetPostByIdUsecase) {
	idString := ctx.Param("id")
	id, err := strconv.Atoi(idString)
	if err != nil {
		ctx.String(http.StatusBadRequest, "Invalid ID format")
		return
	}
	// ctx.String(http.StatusOK, greetings[id])

	result, err := usecase.Execute(id)
	if err != nil {
		handleError(ctx, http.StatusInternalServerError, err)
	} else if result != nil {
		ctx.JSON(http.StatusOK, result)
	} else {
		handleError(ctx, http.StatusNotFound, errors.New("not found"))
	}
}

type PostRequest struct {
	Title  string `json:"title"`
	Body   string `json:"body"`
	UserId int    `json:"user_id"`
}

func PostPosts(ctx *gin.Context, usecase *usecases.CreatePostUsecase) {
	var post PostRequest
	if err := ctx.ShouldBindJSON(&post); err != nil {
		handleError(ctx, http.StatusBadRequest, err)
		return
	}
	result, err := usecase.Execute(post.UserId, post.Title, post.Body)
	if err != nil {
		handleError(ctx, http.StatusBadRequest, err)
		return
	}
	ctx.JSON(http.StatusOK, result)
}
