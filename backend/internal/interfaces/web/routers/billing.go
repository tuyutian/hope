package routers

import (
	"github.com/gin-gonic/gin"

	"backend/internal/interfaces/web/handler"
	"backend/internal/interfaces/web/middleware"
)

func RegisterBillingRouter(r *gin.RouterGroup, h *handler.BillingHandler, authWare *middleware.AuthWare) {
	billingGroup := r.Group("billing", authWare.CheckLogin())

	billingGroup.POST("/list", h.BillList)
	billingGroup.POST("/details", h.BillDetails)
	billingGroup.GET("/current", h.CurrentPeriod)
}
