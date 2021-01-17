package corporation

type options struct {
	skipCache bool
}

type OptionFunc func(*options) *options

func NewOptionFuncs(optionFuncs ...OptionFunc) []OptionFunc {
	return optionFuncs
}

func (s *service) options(optionFuncs []OptionFunc) *options {
	options := &options{}

	for _, optionFunc := range optionFuncs {
		options = optionFunc(options)
	}

	return options
}

func SkipCache() OptionFunc {
	return func(o *options) *options {
		o.skipCache = true
		return o
	}
}
