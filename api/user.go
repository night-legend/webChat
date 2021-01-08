package api

import (
	"net/http"
	"we-chat/models"

	"github.com/gin-gonic/gin"
)

func Login(c *gin.Context) {
	var user models.User
	if err := c.ShouldBind(&user); err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}

	token, err := models.ManageEnv.UserManager.Login(&user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	c.Writer.Header().Set("token", token)
	c.JSON(http.StatusOK, nil)
}

func Register(c *gin.Context) {
	var user models.User
	if err := c.ShouldBind(&user); err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}

	err := models.ManageEnv.UserManager.Register(&user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
	}
	c.JSON(http.StatusOK, user)
}

func GetUser(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, "id must not be empty")
		return
	}
	user, err := models.ManageEnv.UserManager.GetUser(id, "id")
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusOK, user)
}

func SearchUsers(c *gin.Context) {
	name := c.Param("name")
	if name == "" {
		c.JSON(http.StatusBadRequest, "name must not be empty")
		return
	}
	users, err := models.ManageEnv.UserManager.SearchUsers(name)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusOK, users)
}

// API Friends
func GetFriends(c *gin.Context) {
	id := c.GetString("userID")
	if id == "" {
		c.JSON(http.StatusBadRequest, nil)
		return
	}
	users, err := models.ManageEnv.UserManager.ListFriends(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusOK, users)
}

func AddFriend(c *gin.Context) {
	id := c.GetString("userID")
	addID := c.Param("id")

	var option models.AddUserOptions
	_ = c.ShouldBind(&option)

	if err := models.ManageEnv.UserManager.AddFriend(id, addID, option.Content); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
	}  else {
		c.JSON(http.StatusOK, nil)
	}
}

func DeleteFriend(c *gin.Context) {
	requestID := c.GetString("userID")
	destinationID := c.Param("id")

	if requestID == "" || destinationID == "" {
		c.JSON(http.StatusBadRequest, nil)
		return
	}
	err := models.ManageEnv.UserManager.DeleteFriend(requestID, destinationID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusOK, nil)
}

