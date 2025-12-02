package errors

type Errors struct {
	errsmap map[string]error
	err     error
}

func New() *Errors {
	e := new(Errors)
	e.errsmap = make(map[string]error)
	e.err = nil

	return e
}

func (e *Errors) Has() bool {
	if e.err != nil || (e.errsmap != nil && len(e.errsmap) > 0) {
		return true
	}

	return false
}

func (e *Errors) Set(err error) {
	e.err = err
}

func (e *Errors) Get() error {
	return e.err
}

func (e *Errors) SetMap(errsmap map[string]error) {
	e.errsmap = errsmap
}

func (e *Errors) GetMap() (errsmap map[string]error) {
	return e.errsmap
}

func (e *Errors) ClassFor(field string) string {
	if _, exists := e.errsmap[field]; exists {
		return "error"
	}

	return ""
}
