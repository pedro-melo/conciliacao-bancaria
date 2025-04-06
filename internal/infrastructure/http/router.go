package http

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"

	"conciliacao-bancaria/internal/infrastructure/http/handler"
	"conciliacao-bancaria/internal/infrastructure/http/middleware"
)

// SetupRouter configura todas as rotas da API e retorna o router
func SetupRouter(
	billetHandler *handler.BilletHandler,
	paymentHandler *handler.PaymentHandler,
	reconciliationHandler *handler.ReconciliationHandler) *gin.Engine {

	// Inicializa o router Gin com o modo definido
	r := gin.Default()

	// Middleware para logging de requisições
	r.Use(middleware.Logger())

	// Middleware para recuperação de pânico
	r.Use(gin.Recovery())

	// Rota básica para verificação de saúde da API
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "up",
		})
	})

	// Configuração da versão da API
	v1 := r.Group("/api/v1")
	{
		// Rotas para boletos
		billets := v1.Group("/billets")
		{
			billets.POST("", billetHandler.CreateBillet)
			billets.POST("/batch", billetHandler.CreateBilletBatch)
			billets.GET("", billetHandler.ListBillets)
			billets.GET("/:id", billetHandler.GetBillet)
			billets.PUT("/:id", billetHandler.UpdateBillet)
			billets.DELETE("/:id", billetHandler.DeleteBillet)
		}

		// Rotas para pagamentos
		payments := v1.Group("/payments")
		{
			payments.POST("", paymentHandler.CreatePayment)
			payments.POST("/batch", paymentHandler.CreatePaymentBatch)
			payments.GET("", paymentHandler.ListPayments)
			payments.GET("/:id", paymentHandler.GetPayment)
			payments.PUT("/:id", paymentHandler.UpdatePayment)
			payments.DELETE("/:id", paymentHandler.DeletePayment)
		}

		// Rotas para conciliação
		reconciliations := v1.Group("/reconciliations")
		{
			// Rota para iniciar uma nova conciliação
			reconciliations.POST("", reconciliationHandler.CreateReconciliation)

			// Rota para conciliar boletos e pagamentos específicos
			reconciliations.POST("/specific", reconciliationHandler.ReconcileSpecific)

			// Rota para listar todas as conciliações
			reconciliations.GET("", reconciliationHandler.ListReconciliations)

			// Rota para obter detalhes de uma conciliação específica
			reconciliations.GET("/:id", reconciliationHandler.GetReconciliation)

			// Rota para obter histórico de conciliações de um boleto
			reconciliations.GET("/billet/:id", reconciliationHandler.GetBilletReconciliationHistory)

			// Rota para obter histórico de conciliações de um pagamento
			reconciliations.GET("/payment/:id", reconciliationHandler.GetPaymentReconciliationHistory)
		}
	}

	// Rota para documentação da API (Swagger se implementado)
	r.GET("/swagger/*any", gin.WrapH(http.StripPrefix("/swagger", http.FileServer(http.Dir("./swagger")))))

	log.Println("Router configurado com sucesso")
	return r
}
