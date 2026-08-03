
package httpresponse

type Success struct {
	Success    bool        `json:"success"`
	StatusCode int         `json:"statusCode"`
	Message    string      `json:"message"`
	Data       interface{} `json:"data,omitempty"`
}

type Error struct {
	Success      bool   `json:"success"`
	StatusCode   int    `json:"statusCode"`
	Error        bool   `json:"error"`
	ErrorMessage string `json:"errorMessage"`
	ErrorDetails string `json:"errorDetails,omitempty"`
}