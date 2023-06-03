package protoval

type ValidationError struct {
	FieldName string
	Message   string
}

func (validationError ValidationError) GetFieldName() string {
	return validationError.FieldName
}

func (validationError ValidationError) Error() string {
	return validationError.Message
}

func Error(fieldName string, message string) *ValidationError {
	validationError := ValidationError{
		FieldName: fieldName,
		Message:   message,
	}
	return &validationError
}
