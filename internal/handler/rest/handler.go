package rest

// aggregate struct
type APIHandler struct {
	*GenericHandler
	*UserHandler
	*BlogHandler
}

// constructor
func NewAPIHandler(generic *GenericHandler, user *UserHandler, blog *BlogHandler) *APIHandler {
	return &APIHandler{generic, user, blog}
}
