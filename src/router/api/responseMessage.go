package api

type SuccessResponse struct {
    Success bool
    Response any
}

type ErrorResponse struct {
    Success bool
    Type string
    Message string
}
