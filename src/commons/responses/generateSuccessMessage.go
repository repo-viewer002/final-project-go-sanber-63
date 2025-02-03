package responses

func GenerateSuccessMessage(message string) BaseResponse {
	return BaseResponse{
		Status:  "success",
		Message: message,
	}
}

func GenerateSuccessMessageWithData(message string, data interface{}) BaseResponse {
	return BaseResponse{
		Status:  "success",
		Message: message,
		Data:    data,
	}
}
