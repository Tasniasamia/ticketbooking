package user

import (
	"net/http"
	"ticketBooking/internal/user/dto"
    "ticketBooking/internal/httpresponse"
	"github.com/labstack/echo/v5"
)

type handler struct {
	service *Service
}

func NewHandler(s *Service) *handler {
	return &handler{
		service: s,
	}
}

func (h *handler) CreateUser(c *echo.Context) error {


	

	var req dto.CreateRequest;
	if err :=c.Bind(&req); err!=nil{
		return c.JSON(http.StatusBadRequest,httpresponse.Error{
			Code:http.StatusBadRequest,
			Message: "Invalid Request Method",
			Details: err.Error(),
		})


		
	}

	if err :=c.Validate(&req); err!=nil{
		return c.JSON(http.StatusInternalServerError,httpresponse.Error{
			Code:http.StatusInternalServerError,
			Message: "Validation failed",
			Details: err.Error(),
		})


		
	}

	
// 	err :=h.service.repo.CreateUser(&User{
// Name: req.Name,
// Email: req.Email,
// Password: req.Password,
//  })

res,err :=h.service.CreateUser(&req);


	if(err != nil){
		return c.JSON(http.StatusInternalServerError,httpresponse.Error{
			Code:http.StatusInternalServerError,
			Message: "Validation failed",
			Details: err.Error(),
		})
	}

	// return c.JSON(http.StatusCreated,"ok");

		return c.JSON(http.StatusCreated,res);


}