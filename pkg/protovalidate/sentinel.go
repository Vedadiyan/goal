package protoval

type ValidationError struct {
	fieldName string
	message   string
}

func (validationError ValidationError) FieldName() string {
	return validationError.fieldName
}

func (validationError ValidationError) Error() string {
	return validationError.message
}

func Error(fieldName string, message string) *ValidationError {
	validationError := ValidationError{
		fieldName: fieldName,
		message:   message,
	}
	return &validationError
}
