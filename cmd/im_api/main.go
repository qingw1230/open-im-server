package main

import (
	"flag"
	"fmt"
	"strconv"

	"github.com/gin-gonic/gin"
	apiAuth "github.com/qingw1230/studyim/internal/api/auth"
	"github.com/qingw1230/studyim/pkg/common/log"
	"github.com/qingw1230/studyim/pkg/utils"
)

func main() {
	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()
	r.Use(utils.CorsHandler())

	// certificate
	authRouterGroup := r.Group("/auth")
	{
		authRouterGroup.POST("/user_register", apiAuth.UserRegister)
		authRouterGroup.POST("/user_token", apiAuth.UserToken)
	}

	log.NewPrivateLog("api")
	ginPort := flag.Int("port", 10000, "get ginServerPort from cmd, default 10000 as port")
	flag.Parse()
	fmt.Println("api server start...")
	r.Run(":" + strconv.Itoa(*ginPort))
}
