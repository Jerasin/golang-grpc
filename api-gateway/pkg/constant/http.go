package constant

type ResponseStatus int
type Headers int
type General int

// Constant Api
const (
	Success ResponseStatus = iota + 1
	DataNotFound
	UnknownError
	InvalidRequest
	Unauthorized
	Duplicated
	ValidateError
	BadRequest
	RequiredQuery
	DataIsExit
)

func (r ResponseStatus) GetResponseStatusCode() int {
	return [...]int{200, 404, 500, 400, 401, 409, 422, 400, 400, 409}[r-1]
}

func (r ResponseStatus) GetResponseStatus() string {
	return [...]string{"SUCCESS", "DATA_NOT_FOUND", "UNKNOWN_ERROR", "INVALID_REQUEST", "UNAUTHORIZED", "DUPLICATED", "VALIDATE_ERROR", "BAD_REQUEST",
		"REQUIRED_QUERY", "DATA_IS_EXIT"}[r-1]
}

func (r ResponseStatus) GetResponseMessage() string {
	return [...]string{"Success", "Data Not Found", "Unknown Error", "Invalid Request", "Unauthorized", "Duplicated", "ValidateError", "BadRequest", "RequiredQuery", "DataIsExit"}[r-1]
}
