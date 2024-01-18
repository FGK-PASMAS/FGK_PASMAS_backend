package api

type SuccessResponse struct {
    Success bool            `json:"success"`
    Response any            `json:"response"`
}

type ErrorResponse struct {
    Success bool            `json:"success"`
    Type string             `json:"type"`
    Message string          `json:"body"`
}
