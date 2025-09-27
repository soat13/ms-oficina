package error

const (
	HTTPStatusBadRequest    = "400"
	HTTPStatusConflict      = "409"
	HTTPStatusNotFound      = "404"
	HTTPStatusUnprocessable = "422"
)

type Info struct {
	PrivateCode string
	PublicCode  string
}

type Resolver struct {
	errorMap map[error]Info
}

func NewErrorResolver() *Resolver {
	return &Resolver{
		errorMap: map[error]Info{},
	}
}

func (e *Resolver) RegisterError(err error, privateCode string) {
	e.errorMap[err] = Info{
		PrivateCode: privateCode,
		PublicCode:  err.Error(),
	}
}

func (e *Resolver) RegisterHTTPBadRequestError(err error) {
	e.RegisterError(err, HTTPStatusBadRequest)
}

func (e *Resolver) RegisterHTTPConflictError(err error) {
	e.RegisterError(err, HTTPStatusConflict)
}

func (e *Resolver) RegisterHTTPNotFoundError(err error) {
	e.RegisterError(err, HTTPStatusNotFound)
}

func (e *Resolver) RegisterHTTPUnprocessableError(err error) {
	e.RegisterError(err, HTTPStatusUnprocessable)
}

func (e *Resolver) Resolve(error error) (Info, bool) {
	info, success := e.errorMap[error]
	return info, success
}

func (i *Info) Message() string {
	return i.PrivateCode // todo: will be changed to message
}
