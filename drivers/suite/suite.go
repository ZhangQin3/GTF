package ts

type Suite struct {
	Scheme      string
	SuiteParams map[string]interface{}
}

func (s *Suite) SetParam(param string, value interface{}) {
	if s.SuiteParams == nil {
		s.SuiteParams = make(map[string]interface{})
	}
	s.SuiteParams[param] = value
}
