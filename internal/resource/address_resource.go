package resource

import (
	"errors"
	"github.com/NavPool/navpool-hq-api/internal/service"
	"github.com/NavPool/navpool-hq-api/internal/service/account"
	"github.com/gin-gonic/gin"
	uuid "github.com/satori/go.uuid"
	"log"
)

type AddressResource struct {
	addresses *service.AddressService
}

func NewAddressResource(addressService *service.AddressService) *AddressResource {
	return &AddressResource{addressService}
}

type AddressDto struct {
	Hash      string `json:"hash" binding:"required,len=34"`
	Signature string `json:"signature" binding:"required"`
}

func (r *AddressResource) CreateAddress(c *gin.Context) {
	userId, _ := c.Get("id")
	user := userId.(account.User)

	log.Printf("Creating new address for user %s", user.ID.String())

	addressDto := AddressDto{}
	if err := c.BindJSON(&addressDto); err != nil {
		return
	}

	a, err := r.addresses.CreateNewAddress(addressDto.Hash, addressDto.Signature, user)
	if err != nil {
		_ = c.Error(err).SetType(gin.ErrorTypePublic)
		return
	}

	c.JSON(200, a)
}

func (r *AddressResource) GetAddresses(c *gin.Context) {
	userId, _ := c.Get("id")
	user := userId.(account.User)

	addresses, err := r.addresses.GetAddresses(user)
	if err != nil {
		_ = c.Error(err).SetType(gin.ErrorTypePublic)
		return
	}

	c.JSON(200, addresses)
}

func (r *AddressResource) GetAddress(c *gin.Context) {
	userId, _ := c.Get("id")
	user := userId.(account.User)

	id, err := uuid.FromString(c.Param("id"))
	if err != nil {
		_ = c.Error(ErrUnableToFindAddress).SetType(gin.ErrorTypePublic)
		return
	}

	a, err := r.addresses.GetAddress(id, user)
	if err != nil {
		_ = c.Error(ErrUnableToFindAddress).SetType(gin.ErrorTypePublic)
		return
	}

	c.JSON(200, a)
}

func (r *AddressResource) DeleteAddress(c *gin.Context) {
	userId, _ := c.Get("id")
	user := userId.(account.User)

	id, err := uuid.FromString(c.Param("id"))
	if err != nil {
		_ = c.Error(ErrUnableToDeleteAddress).SetType(gin.ErrorTypePublic)
		return
	}

	err = r.addresses.DeleteAddress(id, user)
	if err != nil {
		_ = c.Error(err).SetType(gin.ErrorTypePublic)
		return
	}

	c.JSON(200, nil)
}

var (
	ErrUnableToFindAddress   = errors.New("Unable to find the address on your account")
	ErrUnableToDeleteAddress = errors.New("Unable to delete the address")
)
