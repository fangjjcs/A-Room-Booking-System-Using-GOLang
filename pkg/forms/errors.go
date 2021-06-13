package forms

type errors map[string] []string

// Add an error message
func (e errors) Add(field, msg string){
	e[field] = append(e[field], msg)
}

// check if certain field is filled or not
func (e errors) Get(field string) string{
	es := e[field]
	if len(es) == 0 {
		return ""
	}
	return es[0]
}

