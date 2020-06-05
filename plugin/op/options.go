package op

type Option func(*Options)

type Options struct {
	Es   Es   `json:"es"`   // write to elastic
	Disk Disk `json:"disk"` //write to disk
}

type Disk struct {
	Path string `json:"path"`
}

type Es struct {
	Addr  string `json:"addr"`
	Index string `json:"index"`
}

func SetEs(es Es) Option {
	return func(args *Options) {
		args.Es = es
	}
}

func SetDisk(disk Disk) Option {
	return func(args *Options) {
		args.Disk = disk
	}
}

func (o *Options) Init(ops ...Option) error {
	for _, op := range ops {
		op(o)
	}
	return nil
}
