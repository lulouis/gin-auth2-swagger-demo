package main

import (
	"net/http"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/lulouis/gin-swagger"
	"github.com/lulouis/gin-swagger/swaggerFiles"
	"github.com/lulouis/gin-auth2-swagger-demo/controller"
	_ "github.com/lulouis/gin-auth2-swagger-demo/docs"
	"github.com/lulouis/gin-auth2-swagger-demo/httputil"

	"gopkg.in/oauth2.v3/manage"
	"gopkg.in/oauth2.v3/models"
	aserver "gopkg.in/oauth2.v3/server"
	"gopkg.in/oauth2.v3/store"
	"github.com/lulouis/gin-auth2-swagger-demo/ginserver"

	"github.com/gin-contrib/cors"
)

// @title Restful Oauth2 Swagger Example API
// @version 2018.11.2
// @description This is a Restful Oauth2 Swagger Example API server.
// @termsOfService http://swagger.io/terms/

// hide@contact.name API Support
// hide@contact.url http://www.swagger.io/support
// hide@contact.email support@swagger.io

// hide@license.name Apache 2.0
// hide@license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:8080
// @BasePath /api/v1

// @securityDefinitions.basic BasicAuth

// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization

// @securitydefinitions.oauth2.application OAuth2Application
// @tokenUrl http://localhost:8080/oauth2/token
// @scope.write Grants write access
// @scope.admin Grants read and write access to administrative information

// // @securitydefinitions.oauth2.application OAuth2Application
// // @tokenUrl http://localhost:8080/oauth2/token?grant_type=client_credentials&scope=admin&client_id=000000&client_secret=999999
// // @scope.write Grants write access
// // @scope.admin Grants read and write access to administrative information

// // @securitydefinitions.oauth2.implicit OAuth2Implicit
// // @authorizationUrl https://example.com/oauth/authorize
// // @scope.write Grants write access
// // @scope.admin Grants read and write access to administrative information

// // @securitydefinitions.oauth2.password OAuth2Password
// // @tokenUrl https://example.com/oauth/token
// // @scope.read Grants read access
// // @scope.write Grants write access
// // @scope.admin Grants read and write access to administrative information

// // @securitydefinitions.oauth2.accessCode OAuth2AccessCode
// // @tokenUrl http://localhost:8080/oauth2/token
// // @authorizationUrl http://localhost:8080/oauth2/authorize
// // @scope.admin Grants read and write access to administrative information


func main() {
	//r := gin.Default()
	// gin.SetMode(gin.ReleaseMode)

	manager := manage.NewDefaultManager()
	// token store
	manager.MustTokenStorage(store.NewMemoryTokenStore())
	// client store
	clientStore := store.NewClientStore()
	clientStore.Set("000000", &models.Client{
		ID:     "000000",
		Secret: "999999",
	})
	manager.MapClientStorage(clientStore)
	// Initialize the oauth2 service
	ginserver.InitServer(manager)
	ginserver.SetAllowGetAccessRequest(true)
	ginserver.SetClientInfoHandler(aserver.ClientFormHandler)
	// ginserver.SetUserAuthorizationHandler(aserver.UserAuthorizationHandler)
	r := gin.Default()
	r.Use(cors.Default())
	// r := gin.New()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())
	oauth2 := r.Group("/oauth2")
	{
		oauth2.GET("/token", ginserver.HandleTokenRequest)
		oauth2.POST("/token", ginserver.HandleTokenRequest)
		oauth2.GET("/authorize", ginserver.HandleAuthorizeRequest)
	}

	
	api := r.Group("/api")
	{
		api.Use(ginserver.HandleTokenVerify())
		api.GET("/test", func(c *gin.Context) {
			ti, exists := c.Get("AccessToken")
			if exists {
				c.JSON(http.StatusOK, ti)
				return
			}
			c.String(http.StatusOK, "not found")
		})
	}

	r.GET("/", func(c *gin.Context) {
        c.Redirect(http.StatusMovedPermanently, "/swagger/index.html")
    })
	
	c := controller.NewController()
	
	v1 := r.Group("/api/v1")
	{
		accountsWithToken := v1.Group("/accounts")
		{
			//accountsWithToken.Use(checkClientToken())
			accountsWithToken.Use(ginserver.HandleTokenVerify())
			accountsWithToken.GET("", c.ListAccounts)
		}
		accountsWithOutToken := v1.Group("/accounts")
		{
			accountsWithOutToken.GET(":id", c.ShowAccount)
			accountsWithOutToken.POST("", c.AddAccount)
			accountsWithOutToken.DELETE(":id", c.DeleteAccount)
			accountsWithOutToken.PATCH(":id", c.UpdateAccount)
			accountsWithOutToken.POST(":id/images", c.UploadAccountImage)
		}

		bottles := v1.Group("/bottles")
		{
			bottles.GET(":id", c.ShowBottle)
			bottles.GET("", c.ListBottles)
		}
		admin := v1.Group("/admin")
		{
			//admin.Use(checkClientToken())
			//admin.Use(ginserver.HandleTokenVerify())
			admin.Use(auth())
			admin.POST("/auth", c.Auth)
		}
		examples := v1.Group("/examples")
		{
			examples.GET("ping", c.PingExample)
			examples.GET("calc", c.CalcExample)
			examples.GET("groups/:group_id/accounts/:account_id", c.PathParamsExample)
			examples.GET("header", c.HeaderExample)
			examples.GET("securities", c.SecuritiesExample)
			examples.GET("attribute", c.AttributeExample)
		}
	}

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	r.Run(":8080")
}


func auth() gin.HandlerFunc {
	return func(c *gin.Context) {
		if len(c.GetHeader("Authorization")) == 0 {
			httputil.NewError(c, http.StatusUnauthorized, errors.New("Authorization is required Header"))
			c.Abort()
		}
		c.Next()
	}
}
