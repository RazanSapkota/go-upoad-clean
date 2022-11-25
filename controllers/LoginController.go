package controllers

import (
	"example/go-api/models"
	"example/go-api/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

type LoginController interface {
	Login(ctx *gin.Context)
}

type loginController struct {
	loginService service.LoginService
	jwtService service.JWTService
	
}

func NewLoginController(loginService service.LoginService,jwtService service.JWTService) LoginController{
	return &loginController{	
		loginService: loginService,
		jwtService: jwtService,
	}
}

func (controller *loginController) 	Login(ctx *gin.Context){
 var loginData models.Login
 err:=ctx.ShouldBind(&loginData)
 if err !=nil{
	ctx.JSON(400, gin.H{
		"error":err.Error(),
	})
return
 }
 isUserAuthenticated:=controller.loginService.LoginUser(loginData.Email,loginData.Password)

 if isUserAuthenticated{
	token:=controller.jwtService.GenerateToken(loginData.Email,isUserAuthenticated)
	ctx.JSON(200, gin.H{
		"token":token,
	})
	return
 }
 ctx.JSON(http.StatusBadRequest, gin.H{
	"message":"login error",
}) 

}

