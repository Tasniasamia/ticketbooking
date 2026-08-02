package user

import (
	"net/http"
	"strings"
	"ticketBooking/internal/httpresponse"
	"ticketBooking/internal/user/dto"
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
			Message: strings.Split(err.Error(), ",")[2],
			Details: err.Error(),
		})


		
	}


res,err :=h.service.CreateUser(&req);


	if(err != nil){
		return c.JSON(http.StatusInternalServerError,httpresponse.Error{
			Code:http.StatusInternalServerError,
			Message:strings.Split(err.Error(), ":")[0],
			Details: err.Error(),
		})
	}

return c.JSON(http.StatusCreated,res);


}


func (h *handler) LoginUser(c *echo.Context) error {

    var req dto.LoginRequest;
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
			Message: strings.Split(err.Error(), ",")[2],
			Details: err.Error(),
		})


		
	}


res,err :=h.service.LoginUser(&req);


	if(err != nil){
		return c.JSON(http.StatusInternalServerError,httpresponse.Error{
			Code:http.StatusInternalServerError,
			Message:strings.Split(err.Error(), ":")[0],
			Details: err.Error(),
		})
	}

return c.JSON(http.StatusCreated,res);


}