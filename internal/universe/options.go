package universe

type options struct {
	loc bool
	inv bool
	chr bool

	disableProgress bool
}

type OptionFunc func(*options) *options

func (s *service) options(optionFuncs ...OptionFunc) *options {
	options := &options{
		loc: true,
		inv: true,
		chr: true,
	}

	for _, optionFunc := range optionFuncs {
		options = optionFunc(options)
	}

	return options
}

func WithDisableProgress() OptionFunc {
	return func(o *options) *options {
		o.disableProgress = true
		return o
	}
}

func WithoutChr() OptionFunc {
	return func(o *options) *options {
		o.chr = false
		return o
	}
}

func WithoutInv() OptionFunc {
	return func(o *options) *options {
		o.inv = false
		return o
	}
}

func WithoutLoc() OptionFunc {
	return func(o *options) *options {
		o.loc = false
		return o
	}
}
