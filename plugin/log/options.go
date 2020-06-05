package log

type Option func(*Options)

type Options struct {
	EsAddr   string `json:"es_addr"`
	EsIndex  string `json:"es_index"`
	DiskPath string `json:"disk_path"`
}

func SetEsIndex(index string) Option {
	return func(args *Options) {
		args.EsIndex = index
	}
}

func SetOptions(opts Options) Option {
	return func(args *Options) {
		args = &opts
	}
}

func (o *Options) Init(ops ...Option) error {
	for _, op := range ops {
		op(o)
	}
	return nil
}
