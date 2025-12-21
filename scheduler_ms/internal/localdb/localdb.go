package localdb

import "scheduler/internal/config"

type LocalDb struct {
	opts *config.LocalDbOptions
}

func NewLocalDb(opts *config.LocalDbOptions) *LocalDb {
	return &LocalDb{
		opts: opts,
	}
}
