package actions

import (
	"github.com/Popov-Dmitriy-Ivanovich/Diplom_cmd/models"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

// Get
// @Summary      Create action
// @Description  Создает action. Доступен только админу
// @Tags         User
// @Param        data    body     models.Action  true  "Данные пользователя для создания"
// @Produce      json
// @Success      200  {object}   map[string]models.Action
// @Failure      404  {object}   string
// @Router       /actions [post]
func (a *Action) Create() gin.HandlerFunc {
	return func(c *gin.Context) {
		action := models.Action{}
		if err := c.ShouldBindJSON(&action); err != nil {
			c.AbortWithError(422, err)
			return
		}
		action.StatusID = 3
		validate := validator.New(validator.WithRequiredStructEnabled())
		if err := validate.Struct(action); err != nil {
			c.AbortWithError(422, err)
			return
		}
		db := models.GetDb()
		if err := db.Create(&action).Error; err != nil {
			c.AbortWithError(422, err)
			return
		}
		c.JSON(200, gin.H{"action":action})
	}
}

// Get
// @Summary      Update action
// @Description  Обновляет action. Доступен только админу
// @Tags         User
// @Param        data    body     models.Action  true  "Данные пользователя для action"
// @Produce      json
// @Success      200  {object}   map[string]models.Action
// @Failure      404  {object}   string
// @Router       /user/{id} [put]
func (a *Action) Update() gin.HandlerFunc {
	return func(c *gin.Context) {
		body := models.Action{}
		if err := c.ShouldBindJSON(&body); err != nil {
			c.AbortWithError(422, err)
			return
		}
		action := models.Action{}
		db := models.GetDb()
		if err := db.First(&action, c.Param("id")).Error; err != nil {
			c.AbortWithError(404, err)
			return
		}
		action.LastLaunch = body.LastLaunch
		action.Name = body.Name
		action.Description = body.Description
		action.ShortDesc = body.ShortDesc
		action.Cmd = body.Cmd
		if err := db.Save(&action).Error; err != nil {
			c.AbortWithError(500, err)
			return
		}
		c.JSON(200, gin.H{"action": action})
	}
}

// Get
// @Summary      Delete user
// @Description  Удаляет пользователя. Доступен только админу
// @Tags         User
// @Param        id    path     int  true  "id User"
// @Produce      json
// @Success      200  {object}   map[string]models.Action
// @Failure      404  {object}   string
// @Router       /user/{id} [delete]
func (u *Action) Delete() gin.HandlerFunc {
	return func(c *gin.Context) {
		action := models.Action{}
		db := models.GetDb()

		if err := db.First(&action, c.Param("id")).Error; err != nil {
			c.AbortWithError(404, err)
			return
		}

		if err := db.Delete(&action).Error; err != nil {
			c.AbortWithError(500, err)
			return
		}

		c.JSON(200, gin.H{"action": action})
	}
}
