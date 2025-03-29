package forms

// hold validation errs
type errors map[string][]string

// add err msg to the errors map for a given field
func (e errors) Add(field, message string) {
	e[field] = append(e[field], message)
}
func (e errors) Get(field string) string {
	es := e[field]
	if len(es) == 0 {
		return ""
	}
	return es[0]
}
