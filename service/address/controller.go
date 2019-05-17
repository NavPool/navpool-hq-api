package address

import (
	"github.com/NavPool/navpool-hq-api/service/account"
	"github.com/gin-gonic/gin"
	uuid "github.com/satori/go.uuid"
	"log"
)

type Controller struct{}

type AddressDto struct {
	Hash      string `json:"hash" binding:"required,len=34"`
	Signature string `json:"signature" binding:"required"`
}

func (controller *Controller) CreateAddress(c *gin.Context) {
	userId, _ := c.Get("id")
	user := userId.(account.User)

	log.Printf("Creating new address for user %s", user.ID.String())

	addressDto := AddressDto{}
	if err := c.BindJSON(&addressDto); err != nil {
		return
	}

	address, err := CreateNewAddress(addressDto, user)
	if err != nil {
		_ = c.Error(err).SetType(gin.ErrorTypePublic)
		return
	}

	c.JSON(200, address)
}

func (controller *Controller) GetAddresses(c *gin.Context) {
	userId, _ := c.Get("id")
	user := userId.(account.User)

	addresses, err := GetAddresses(user)
	if err != nil {
		_ = c.Error(err).SetType(gin.ErrorTypePublic)
		return
	}

	c.JSON(200, addresses)
}

func (controller *Controller) GetAddress(c *gin.Context) {
	userId, _ := c.Get("id")
	user := userId.(account.User)

	id, err := uuid.FromString(c.Param("id"))
	if err != nil {
		_ = c.Error(ErrorUnableToFindAddress).SetType(gin.ErrorTypePublic)
		return
	}

	address, err := GetAddress(id, user)
	if err != nil {
		_ = c.Error(ErrorUnableToFindAddress).SetType(gin.ErrorTypePublic)
		return
	}

	c.JSON(200, address)
}

func (controller *Controller) DeleteAddress(c *gin.Context) {
	userId, _ := c.Get("id")
	user := userId.(account.User)

	id, err := uuid.FromString(c.Param("id"))
	if err != nil {
		_ = c.Error(ErrorUnableToDeleteAddress).SetType(gin.ErrorTypePublic)
		return
	}

	err = DeleteAddress(id, user)
	if err != nil {
		_ = c.Error(err).SetType(gin.ErrorTypePublic)
		return
	}

	c.JSON(200, nil)
}
