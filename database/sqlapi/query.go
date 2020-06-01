package sqlapi

import "fmt"

type Query struct {
	text string
	args []interface{}
}

func NewQuery(text string) *Query {
	return &Query{text: text}
}

func (q *Query) Text() string {
	return q.text
}

func (q *Query) Args() []interface{} {
	return q.args
}

func (q *Query) WithArgs(args ...interface{}) *Query {
	q.args = args
	return q
}

func (q *Query) Inject(values ...interface{}) *Query {
	q.text = fmt.Sprintf(q.text, values...)
	return q
}

func (q *Query) FlatArgs() *Query {
	q.text, q.args = flatQuery(q.text, q.args)
	return q
}

func (q *Query) String() string {
	return formatQuery(q.text, BindVarChar, q.args...)
}
