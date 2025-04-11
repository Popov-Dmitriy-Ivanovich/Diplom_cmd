package actions

import (
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/Popov-Dmitriy-Ivanovich/Diplom_auth/middleware"
	authmodels "github.com/Popov-Dmitriy-Ivanovich/Diplom_auth/models"
	"github.com/Popov-Dmitriy-Ivanovich/Diplom_cmd/kafka"
	"github.com/Popov-Dmitriy-Ivanovich/Diplom_cmd/models"
	tgApi "github.com/Popov-Dmitriy-Ivanovich/Diplom_telegram/api"
	"github.com/gin-gonic/gin"
)

type Action struct {
}

var Limit uint = 999999

//AR_VIEW_AND_RUN_ACTION

func (a *Action) WriteRoutes(rg *gin.RouterGroup) {
	actionGroup := rg.Group("/actions")
	actionGroup.Use(middleware.AuthMiddleware(authmodels.AR_VIEW_AND_RUN_ACTION))
	actionGroup.GET("/", a.Get())
	actionGroup.GET("/:id", a.GetId())
	actionGroup.GET("/:id/run", a.Run())
	actionGroup.GET("/:id/status", a.Status())
	actionGroup.GET("/:id/stop", a.Stop())
	actionGroup.Use(middleware.AuthMiddleware(authmodels.AR_CREATE_ACTION))
	actionGroup.POST("/", a.Create())
	actionGroup.PUT("/:id", a.Update())
	actionGroup.DELETE("/:id", a.Delete())
	actionGroup.GET("/limit/:limit", a.Limit())
}

// Get
// @Summary      Get list of actions ids
// @Description  Возращает список id всех доступных action
// @Tags         Actions
// @Produce      json
// @Success      200  {object}   map[string][]uint
// @Failure      500  {object}   string
// @Failure      404  {object}   string
// @Router       /actions [get]
func (a *Action) Get() gin.HandlerFunc {
	return func(c *gin.Context) {
		db := models.GetDb()
		ids := []uint{}
		kafka.GetStats()
		if err := db.Model(models.Action{}).Pluck("id", &ids).Error; err != nil {
			c.AbortWithError(500, err)
			return
		}
		c.JSON(200, gin.H{"ids": ids, "stats": kafka.Stats})
	}
}

// Get
// @Summary      Get concrete Action
// @Description  Возращает Action соответсвующую указанному ID
// @Tags         Actions
// @Param        id    path     int  true  "id Action"
// @Produce      json
// @Success      200  {object}   map[string]models.Action
// @Failure      404  {object}   string
// @Router       /actions/{id} [get]
func (a *Action) GetId() gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		db := models.GetDb()
		action := models.Action{}
		if err := db.Preload("Events").First(&action, id).Error; err != nil {
			c.AbortWithError(404, err)
			return
		}
		c.JSON(200, gin.H{"action": action})
	}
}

// Get
// @Summary      Run concrete Action
// @Description  Запускает Action соответсвующую указанному ID
// @Tags         Actions
// @Param        id    path     int  true  "id Action"
// @Produce      json
// @Success      200  {object}   map[string]any
// @Failure      404  {object}   string
// @Router       /actions/{id}/run [get]
func (a *Action) Run() gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		db := models.GetDb()
		action := models.Action{}

		activeActions := int64(0)
		if err := db.Model(&models.Action{}).Debug().Where("status_id = ? or status_id = ?", 1, 2).Count(&activeActions).Error; err != nil {
			c.AbortWithError(500, err)
			return
		}
		fmt.Println("Active actions ", activeActions)
		fmt.Println("Active limit ", Limit)
		if activeActions >= int64(Limit) {
			c.AbortWithError(422, errors.New("превышен лимит запущенных action"))
			return
		}

		if err := db.First(&action, id).Error; err != nil {
			c.AbortWithError(404, err)
			return
		}
		if err := kafka.RunAction(action); err != nil {
			c.AbortWithError(500, err)
			return
		}
		action.StatusID = 2
		action.Events = append(action.Events, models.Event{Type: "Launch", TimeStamp: models.TimeStamp{time.Now()}})
		if err := db.Save(&action).Error; err != nil {
			c.AbortWithError(500, err)
			return
		}
		tgApi.Notify("Действие " + action.Name + "переведено в статус Запускается")
		c.JSON(200, gin.H{"action": action.Cmd})
	}
}

// Get
// @Summary      Get concrete Action's status
// @Description  Возвращает информацию о статусе Action
// @Tags         Actions
// @Param        id    path     int  true  "id Action"
// @Produce      json
// @Success      200  {object}   map[string]models.ActionStatus
// @Failure      404  {object}   string
// @Router       /actions/{id}/status [get]
func (a *Action) Status() gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		db := models.GetDb()
		action := models.Action{}

		if err := db.Preload("Status").First(&action, id).Error; err != nil {
			c.AbortWithError(404, err)
			return
		}

		c.JSON(200, gin.H{"status": action.Status})
	}
}

// Get
// @Summary      Stops concrete Action
// @Description  Останавливает Action соответсвующую указанному ID
// @Tags         Actions
// @Param        id    path     int  true  "id Action"
// @Produce      json
// @Success      200  {object}   map[string]any
// @Failure      404  {object}   string
// @Router       /actions/{id}/stop [get]
func (a *Action) Stop() gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		db := models.GetDb()
		action := models.Action{}
		if err := db.First(&action, id).Error; err != nil {
			c.AbortWithError(404, err)
			return
		}
		if err := kafka.StopAction(action); err != nil {
			c.AbortWithError(500, err)
			return
		}
		action.StatusID = 5
		action.Events = append(action.Events, models.Event{Type: "Stopped", TimeStamp: models.TimeStamp{time.Now()}})
		if err := db.Save(&action).Error; err != nil {
			c.AbortWithError(500, err)
			return
		}
		tgApi.Notify("Действие " + action.Name + "переведено в статус Останавливается")
		c.JSON(200, gin.H{"action": action.Cmd})
	}
}

// Get
// @Summary      Set launched action
// @Description  Устанавливает лимит запущенных action
// @Tags         Actions
// @Param        limit    path     int  true  "id Action"
// @Produce      json
// @Success      200  {object}   map[string]any
// @Failure      404  {object}   string
// @Router       /actions/limit/{limit} [get]
func (a *Action) Limit() gin.HandlerFunc {
	return func(c *gin.Context) {
		limit := c.Param("limit")
		uintLimit, err := strconv.ParseUint(limit, 10, 64)
		if err != nil {
			c.AbortWithError(422, err)
			return
		}
		Limit = uint(uintLimit)
	}
}
