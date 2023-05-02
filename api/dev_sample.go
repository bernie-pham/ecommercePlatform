package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (server *Server) sampleData(ctx *gin.Context) {
	// ids := []string{
	// 	"4c69b61db1fc16e7013b43fc926e502d",
	// 	"66d49bbed043f5be260fa9f7fbff5957",
	// 	"2c55cae269aebf53838484b0d7dd931a",
	// 	"18018b6bc416dab347b1b7db79994afa",
	// 	"e04b990e95bf73bbe6a3fa09785d7cd0",
	// 	"40d3cd16b41970ae6872e914aecf2c8e",
	// 	"bc178f33a04dbccefa95b165f8b56830",
	// 	"cc2083338a16c3fe2f7895289d2e98fe",
	// 	"7b0746d8afc8462ba17f8a763d9d5f1e",
	// 	"c5f4c94653a3befd8dd16adf2914c04e",
	// 	"615ba903c134f439eaf8cdd1678ceb5c",
	// }
	err := server.taskDistributor.DistributeSyncAllTagDataTask(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, "failed to sample first 1000 tag rows")
		return
	}
	err = server.taskDistributor.DistributeSyncAllProductDataTask(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, "failed to sample first 1000 product rows")
		return
	}
}
