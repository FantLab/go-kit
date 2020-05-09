package sqlapi

type NoRows struct {
	Err error
}

func (r NoRows) Error() error {
	return r.Err
}

func (r NoRows) Scan(output interface{}) error {
	return r.Err
}
