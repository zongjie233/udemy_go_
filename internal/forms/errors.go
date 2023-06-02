package forms

type errors map[string][]string

// Add 给特定的表单字段添加错误信息
func (e errors) Add(field, message string) {
	e[field] = append(e[field], message)
}

// Get 返回第一个错误信息
func (e errors) Get(field string) string {
	es := e[field]
	if len(es) == 0 {
		return ""
	}
	return es[0]
}
