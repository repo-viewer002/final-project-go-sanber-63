package responses

func GenerateFailMessage(message string) BaseResponse {
	return BaseResponse{
		Status:  "fail",
		Message: message,
		Data:    nil,
	}
}
