package services

import (
	"errors"
	"fmt"
	"github.com/babilu-online/common/context"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"os"
	"strconv"
)

// @title NFT Aggregator Swap API
// @version 1.1.27
// @description NFT Aggregator Swap API
// @schemes https

// @host https://agg.alphabatem.com
// @BasePath /
// @query.collection.format multi

type HttpService struct {
	context.DefaultService
	BaseURL string
	Port    int

	imgSvc *ImageService
}

var ErrUnauthorized = errors.New("unauthorized")
var DeleteResponseOK = `{"status": 200, "error": ""}`

func (svc HttpService) Id() string {
	return "http"
}

func (svc *HttpService) Configure(ctx *context.Context) error {
	port := os.Getenv("HTTP_PORT")
	portFlag, err := strconv.Atoi(port)
	if err != nil {
		return err
	}

	svc.Port = portFlag

	return svc.DefaultService.Configure(ctx)
}

func (svc *HttpService) Start() error {
	svc.imgSvc = svc.Service(IMG_SVC).(*ImageService)
	r := gin.Default()

	r.Use(gin.Recovery())

	config := cors.DefaultConfig()
	config.AllowAllOrigins = true
	config.AllowCredentials = true
	config.AddAllowHeaders("Authorization")
	r.Use(cors.New(config))

	//r.Static("static", "static")
	//r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))

	//Validation endpoints
	r.GET("/ping", svc.ping)

	v1 := r.Group("/v1")
	//docs.SwaggerInfo.BasePath = "/v1"

	v1.GET("nfts/:id", svc.showNFT)
	v1.GET("nfts/:id/image", svc.showNFTImage)
	v1.GET("nfts/:id/image.gif", svc.showNFTImage)
	v1.GET("nfts/:id/image.png", svc.showNFTImage)
	v1.GET("nfts/:id/image.jpg", svc.showNFTImage)
	v1.GET("nfts/:id/image.jpeg", svc.showNFTImage)
	v1.GET("nfts/:id/media", svc.showNFTMedia)

	r.NoRoute(func(c *gin.Context) {
		c.JSON(404, gin.H{"code": "PAGE_NOT_FOUND", "message": "Page not found"})
	})

	return r.Run(fmt.Sprintf(":%v", svc.Port))
}

type Pong struct {
	Message string `json:"message"`
}

//
// @Summary Ping liquify service
// @Accept  json
// @Produce json
// @Router /ping [get]
func (svc *HttpService) ping(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "pong",
	})
}

//
// @Summary Ping liquify service
// @Accept  json
// @Produce json
// @Router /nfts/{id} [get]
func (svc *HttpService) showNFT(c *gin.Context) {

	media, err := svc.imgSvc.Media(c.Param("id"))
	if err != nil {
		svc.httpError(c, err)
		return
	}

	c.JSON(200, media)
}

//
// @Summary Ping liquify service
// @Accept  json
// @Produce json
// @Router/nfts/{id}/image [get]
func (svc *HttpService) showNFTImage(c *gin.Context) {
	err := svc.imgSvc.ImageFile(c, c.Param("id"))
	if err != nil {
		svc.mediaError(c, err)
		return
	}
}

//
// @Summary Ping liquify service
// @Accept  json
// @Produce json
// @Router /nfts/{id}/media [get]
func (svc *HttpService) showNFTMedia(c *gin.Context) {
	err := svc.imgSvc.MediaFile(c, c.Param("id"))
	if err != nil {
		svc.mediaError(c, err)
		return
	}
}

func (svc *HttpService) paramErr(c *gin.Context, err error) {
	c.JSON(400, gin.H{
		"error": err.Error(),
	})
}

func (svc *HttpService) httpError(c *gin.Context, err error) {
	c.JSON(500, gin.H{
		"error": err.Error(),
	})
}

//TODO Replace with placeholder image
func (svc *HttpService) mediaError(c *gin.Context, err error) {
	c.JSON(200, gin.H{
		"error": err.Error(),
	})
}