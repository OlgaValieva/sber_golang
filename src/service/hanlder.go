package service

import (
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	"log"
	"net/http"
	"strconv"
	"time"
)

const (
	getRep     = "/api/v1/get_report"
	setRep     = "/api/v1/set_report"
	observTime = "/api/v1/get_observation_time"
)

func GetReportHandler(router *gin.Engine, db *sqlx.DB) {
	// Обработчик для эндпоинта GET /api/v1/get_report
	router.GET(getRep, func(c *gin.Context) {
		senderIP := c.ClientIP()
		// Время пришедшего запроса
		log.Println("Request received from", senderIP, getRep)

		reportIdSrt := c.Query("report_id")
		reportId, _ := strconv.Atoi(reportIdSrt)

		var reportInfo string
		err := db.QueryRow("SELECT report_info FROM reports WHERE report_id = $1", reportId).Scan(&reportInfo)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error_msg":   "Не удалось получить отчет",
				"report_info": "",
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"error_msg":   "",
			"report_info": reportInfo,
		})
	})
}

func SetReportHandler(router *gin.Engine, db *sqlx.DB) {
	// Обработчик для эндпоинта POST /api/v1/set_report
	router.POST(setRep, func(c *gin.Context) {
		// Время пришедшего запроса
		senderIP := c.ClientIP()
		log.Println("Request received from", senderIP, setRep)
		var requestBody struct {
			ReportInfo string `json:"report_info"`
			ModelId    int    `json:"model_id"`
		}

		if err := c.ShouldBindJSON(&requestBody); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error_msg": "Ошибка валидации данных",
			})
			return
		}

		creationTime := time.Now()
		_, err := db.Exec("INSERT INTO reports (report_id,report_info, creation_time, model_id) VALUES ((select report_id from reports order by report_id desc limit 1)+1, $1, $2, $3)", requestBody.ReportInfo, creationTime, requestBody.ModelId)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error_msg": "Не удалось записать отчет",
				"report_id": 0,
				"err":       err,
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"error_msg": "",
		})
	})
}

func GetObservTimeHandler(r *gin.Engine, db *sqlx.DB) {
	// Обработчик для эндпоинта GET /api/v1/get_observation_time
	r.GET(observTime, func(c *gin.Context) {
		// Время пришедшего запроса
		senderIP := c.ClientIP()
		log.Println("Request received from", senderIP, observTime)
		modelID := c.Query("model_id")

		var maxObservationPeriod string
		// err := db.QueryRow("SELECT MAX(creation_time) FROM reports WHERE model_id = $1", modelID).Scan(&maxObservationPeriod)
		err := db.QueryRow("select max (day_diff) from (select creation_time - LAG(creation_time) over (order by creation_time) as day_diff FROM (select * from reports where model_id = $1) as modal order by day_diff ) as diff_table", modelID).Scan(&maxObservationPeriod)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error_msg":              err,
				"max_observation_period": "",
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"error_msg":              "",
			"max_observation_period": maxObservationPeriod, //strconv.Itoa(observationPeriod),
		})
	})
}
