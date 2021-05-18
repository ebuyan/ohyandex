package device

type ActionResult struct {
	Status  string `json:"status"`
	Code    string `json:"error_code"`
	Message string `json:"error_message"`
}

func NewActionResult(status int, err error) (actionResult ActionResult) {
	if err != nil {
		actionResult.Status = "ERROR"
		actionResult.Code = "INVALID_ACTION"
		actionResult.Message = err.Error()
		return
	}
	switch status {
	case 200:
		actionResult.Status = "DONE"
	case 404:
		actionResult.Status = "ERROR"
		actionResult.Code = "INVALID_ACTION"
		actionResult.Message = "Устройство не найдено"
	case 400:
		actionResult.Status = "ERROR"
		actionResult.Code = "INVALID_ACTION"
		actionResult.Message = "Ошибка выполнения команды"
	}
	return
}
