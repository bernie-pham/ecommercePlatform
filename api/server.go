package api

import (
	"fmt"

	"github.com/bernie-pham/ecommercePlatform/async"
	db "github.com/bernie-pham/ecommercePlatform/db/sqlc"
	"github.com/bernie-pham/ecommercePlatform/token"
	ultils "github.com/bernie-pham/ecommercePlatform/ultil"
	"github.com/gin-gonic/gin"
)

var (
	ErrBadRequestParameter = fmt.Errorf("bad request parameter")
	ErrWrongPassword       = fmt.Errorf("wrong password")
	ErrNotFound            = fmt.Errorf("not found user")
	ErrInteralErrServer    = fmt.Errorf("internal error server")
)

type Server struct {
	router          *gin.Engine
	config          ultils.Config
	store           db.Store
	tokenMaker      token.TokenMaker
	taskDistributor async.TaskDistributor
}

func NewServer(
	config ultils.Config,
	store db.Store,
	tokenMaker token.TokenMaker,
	taskDistributor async.TaskDistributor,
) (*Server, error) {
	server := &Server{
		store:           store,
		config:          config,
		tokenMaker:      tokenMaker,
		taskDistributor: taskDistributor,
	}
	server.setupRouter()
	return server, nil
}

func (server *Server) setupRouter() {
	router := gin.Default()
	// router.Use(cors.New(cors.Config{
	// 	AllowAllOrigins: true,
	// }))
	// Serving static file
	public_asset_fs := gin.Dir(fmt.Sprintf("%s/public", server.config.StaticRoot), false)
	private_asset_fs := gin.Dir(fmt.Sprintf("%s/private", server.config.StaticRoot), false)

	router.Use(CORSMiddleware())

	router.POST("/dev/elastic/sample", server.sampleData)

	router.POST("/users", server.createUser)
	router.PATCH("/users", server.UpdateUser)
	router.POST("/users/login", server.loginUser)
	router.GET("/users/forgot_password", server.ForgotPassword)
	router.POST("/users/reset_password", server.ResetPassword)

	router.POST("/merchants", server.CreateMerchant)
	router.PATCH("/merchants", server.UpdateMerchant)
	router.GET("/products/tag/:id", server.ListProductsWithTag)
	router.GET("/product", server.ListProducts)
	router.GET("/product/:id", server.ListProductEntries)
	router.GET("/deal", server.getOrdDeal)

	// Serving public assets
	router.StaticFS("/images", public_asset_fs)
	authRoutes := router.Group("/").Use(authMiddleware(server.tokenMaker))
	// Serving private assets
	authRoutes.StaticFS("/private/images", private_asset_fs)

	authRoutes.POST("/cart/add", server.AddCartItem)
	authRoutes.GET("/cart", server.ListCartItems)
	authRoutes.PUT("/cart/item", server.updateCartItemQty)
	authRoutes.DELETE("/cart/item", server.deleteCartItem)

	authRoutes.POST("/order", server.CreateOrder)

	authRoutes.GET("/notifications", server.ListNofications)

	merchantRoutes := router.Group("/merchant")

	merchantRoutes.Use(authMiddleware(server.tokenMaker))
	merchantRoutes.Use(merchantAuthMiddleware())

	merchantRoutes.GET("/orders/:id", server.GetMerchantOrderDetails)

	// router.POST("/merchants/", server.loginUserForThirdParty)
	// router.POST("/token/refresh_token", server.refreshToken)

	// authRoutes.POST("/accounts", server.createAccount)
	// authRoutes.GET("/accounts/:id", server.getAccount)
	// authRoutes.GET("/accounts", server.listAccount)
	// authRoutes.POST("/transfers", server.createTransfer)

	server.router = router
}

func (server *Server) Start(address string) error {
	return server.router.Run(address)
}
func errorResponse(err error) gin.H {
	return gin.H{"Error": err.Error()}
}

func getOptionalString(str string) (string, bool) {
	if len(str) > 0 {
		return str, true
	}
	return "", false
}

func getOptionalInt(value int) (int, bool) {
	if value > 0 {
		return value, true
	}
	return 0, false
}
