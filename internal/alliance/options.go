package alliance

type options struct{}

type OptionFunc func(*options) *options

// func (s *service) options(optionFuncs []OptionFunc) *options {
// 	options := &options{}

// 	for _, optionFunc := range optionFuncs {
// 		options = optionFunc(options)
// 	}

// 	return options
// }
