package beans

type TestPostRequestData struct {
	Name        string `json:"name" binding:"required"`
	DisplayName string `json:"display_name" binding:"required"`
	Description string `json:"description" binding:"required"`
}
