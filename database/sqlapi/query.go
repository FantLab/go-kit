package sqlapi

import "fmt"

type Query struct {
	text string
	args []interface{}
}

func NewQuery(text string) Query {
	return Query{text: text}
}

func (q Query) Text() string {
	return q.text
}

func (q Query) Args() []interface{} {
	return q.args
}

func (q Query) WithArgs(args ...interface{}) Query {
	return Query{
		text: q.text,
		args: args,
	}
}

func (q Query) Inject(values ...interface{}) Query {
	return Query{
		text: fmt.Sprintf(q.text, values...),
		args: q.args,
	}
}

func (q Query) FlatArgs() Query {
	newArgs, counts := flatArgs(q.args...)
	newQuery := expandQuery(q.text, BindVarChar, counts)

	return Query{
		text: newQuery,
		args: newArgs,
	}
}

func (q Query) String() string {
	return formatQuery(q.text, BindVarChar, q.args...)
}
