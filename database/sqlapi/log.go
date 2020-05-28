package sqlapi

import (
	"context"
	"time"
)

type LogEntry struct {
	Query    func() string
	Rows     int64
	Err      error
	Time     time.Time
	Duration time.Duration
}

type LogFunc func(context.Context, LogEntry)

func Log(db DB, f LogFunc) DB {
	return &logDB{db: db, f: f}
}

// *******************************************************

type logDB struct {
	db DB
	f  LogFunc
}

func (l logDB) Write(ctx context.Context, q Query) Result {
	return logRW{rw: l.db, f: l.f}.Write(ctx, q)
}

func (l logDB) Read(ctx context.Context, q Query, output interface{}) error {
	return logRW{rw: l.db, f: l.f}.Read(ctx, q, output)
}

func (l logDB) InTransaction(perform func(ReaderWriter) error) error {
	return l.db.InTransaction(func(rw ReaderWriter) error {
		return perform(logRW{rw: rw, f: l.f})
	})
}

// *******************************************************

type logRW struct {
	rw ReaderWriter
	f  LogFunc
}

func (l logRW) Write(ctx context.Context, q Query) Result {
	t := time.Now()
	result := l.rw.Write(ctx, q)
	l.f(ctx, LogEntry{
		Query:    q.String,
		Rows:     result.RowsAffected,
		Err:      result.Error,
		Time:     t,
		Duration: time.Since(t),
	})
	return result
}

func (l logRW) Read(ctx context.Context, q Query, output interface{}) error {
	t := time.Now()
	err := l.rw.Read(ctx, q, output)
	l.f(ctx, LogEntry{
		Query:    q.String,
		Rows:     -1,
		Err:      err,
		Time:     t,
		Duration: time.Since(t),
	})
	return err
}
