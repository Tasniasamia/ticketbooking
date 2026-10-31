package subscribe
import ("net/http"; "ticketBooking/internal/httpresponse"; "ticketBooking/internal/subscribe/dto"; "ticketBooking/internal/utils/query"; "github.com/labstack/echo/v5")
type handler struct{ svc *Service }
func NewHandler(svc *Service) *handler { return &handler{svc: svc} }
func (h *handler) Create(c *echo.Context) error {
	var req dto.CreateSubscribeRequest
	if err := c.Bind(&req); err != nil { return c.JSON(http.StatusBadRequest, httpresponse.Error{Success:false,StatusCode:400,Error:true,ErrorMessage:"Invalid body",ErrorDetails:err.Error()}) }
	if err := c.Validate(&req); err != nil { return c.JSON(http.StatusBadRequest, httpresponse.Error{Success:false,StatusCode:400,Error:true,ErrorMessage:"Validation failed",ErrorDetails:err.Error()}) }
	res, err := h.svc.Create(&req)
	if err != nil { return c.JSON(http.StatusBadRequest, httpresponse.Error{Success:false,StatusCode:400,Error:true,ErrorMessage:err.Error()}) }
	return c.JSON(http.StatusCreated, httpresponse.Success{Success:true,StatusCode:201,Message:"Subscribed successfully",Data:res})
}
func (h *handler) ListAdmin(c *echo.Context) error {
	data, err := h.svc.ListAdmin(query.Parse(c))
	if err != nil { return c.JSON(http.StatusInternalServerError, httpresponse.Error{Success:false,StatusCode:500,Error:true,ErrorMessage:"Failed",ErrorDetails:err.Error()}) }
	return c.JSON(http.StatusOK, httpresponse.Success{Success:true,StatusCode:200,Message:"OK",Data:data})
}
