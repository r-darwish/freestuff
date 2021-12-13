package freestuff

func ResponseOk(statusCode int) bool {
	return statusCode >= 200 && statusCode < 300
}
