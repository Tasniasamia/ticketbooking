package testimonial
import ("net/http"; "strconv"; "ticketBooking/internal/httpresponse"; "ticketBooking/internal/testimonial/dto"; "ticketBooking/internal/utils/query"; "github.com/labstack/echo/v5")
type handler struct{ svc *Service }
func NewHandler(svc *Service) *handler { return &handler{svc: svc} }
func (h *handler) Create(c *echo.Context) error {
	var req dto.CreateTestimonialRequest
	if err := c.Bind(&req); err != nil { return c.JSON(http.StatusBadRequest, httpresponse.Error{Success:false,StatusCode:400,Error:true,ErrorMessage:"Invalid request body",ErrorDetails:err.Error()}) }
	if err := c.Validate(&req); err != nil { return c.JSON(http.StatusBadRequest, httpresponse.Error{Success:false,StatusCode:400,Error:true,ErrorMessage:"Validation failed",ErrorDetails:err.Error()}) }
	res, err := h.svc.Create(&req)
	if err != nil { return c.JSON(http.StatusBadRequest, httpresponse.Error{Success:false,StatusCode:400,Error:true,ErrorMessage:err.Error()}) }
	return c.JSON(http.StatusCreated, httpresponse.Success{Success:true,StatusCode:201,Message:"Testimonial created successfully",Data:res})
}
func (h *handler) Update(c *echo.Context) error {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil { return c.JSON(http.StatusBadRequest, httpresponse.Error{Success:false,StatusCode:400,Error:true,ErrorMessage:"Invalid ID"}) }
	var req dto.UpdateTestimonialRequest
	if err := c.Bind(&req); err != nil { return c.JSON(http.StatusBadRequest, httpresponse.Error{Success:false,StatusCode:400,Error:true,ErrorMessage:"Invalid body",ErrorDetails:err.Error()}) }
	res, err := h.svc.Update(uint(id), &req)
	if err != nil { return c.JSON(http.StatusBadRequest, httpresponse.Error{Success:false,StatusCode:400,Error:true,ErrorMessage:err.Error()}) }
	return c.JSON(http.StatusOK, httpresponse.Success{Success:true,StatusCode:200,Message:"Testimonial updated successfully",Data:res})
}
func (h *handler) Delete(c *echo.Context) error {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil { return c.JSON(http.StatusBadRequest, httpresponse.Error{Success:false,StatusCode:400,Error:true,ErrorMessage:"Invalid ID"}) }
	if err := h.svc.Delete(uint(id)); err != nil { return c.JSON(http.StatusBadRequest, httpresponse.Error{Success:false,StatusCode:400,Error:true,ErrorMessage:err.Error()}) }
	return c.JSON(http.StatusOK, httpresponse.Success{Success:true,StatusCode:200,Message:"Testimonial deleted successfully",Data:nil})
}
func (h *handler) GetByID(c *echo.Context) error {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil { return c.JSON(http.StatusBadRequest, httpresponse.Error{Success:false,StatusCode:400,Error:true,ErrorMessage:"Invalid ID"}) }
	res, err := h.svc.GetByID(uint(id))
	if err != nil { return c.JSON(http.StatusNotFound, httpresponse.Error{Success:false,StatusCode:404,Error:true,ErrorMessage:err.Error()}) }
	return c.JSON(http.StatusOK, httpresponse.Success{Success:true,StatusCode:200,Message:"OK",Data:res})
}
func (h *handler) ListAdmin(c *echo.Context) error {
	data, err := h.svc.ListAdmin(query.Parse(c))
	if err != nil { return c.JSON(http.StatusInternalServerError, httpresponse.Error{Success:false,StatusCode:500,Error:true,ErrorMessage:"Failed",ErrorDetails:err.Error()}) }
	return c.JSON(http.StatusOK, httpresponse.Success{Success:true,StatusCode:200,Message:"OK",Data:data})
}
func (h *handler) ListPublic(c *echo.Context) error {
	data, err := h.svc.ListPublic(query.Parse(c))
	if err != nil { return c.JSON(http.StatusInternalServerError, httpresponse.Error{Success:false,StatusCode:500,Error:true,ErrorMessage:"Failed",ErrorDetails:err.Error()}) }
	return c.JSON(http.StatusOK, httpresponse.Success{Success:true,StatusCode:200,Message:"OK",Data:data})
}
