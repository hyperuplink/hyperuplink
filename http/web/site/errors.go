package site

func (s *Site) SetErrors(errsmap map[string]error) {
	s.errsmap = errsmap
}

func (s *Site) GetErrors() (errsmap map[string]error) {
	return s.errsmap
}

func (s *Site) GetErrorClassFor(field string) string {
	if _, exists := s.errsmap[field]; exists {
		return "error"
	}

	return ""
}
