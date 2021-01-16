package athena

type Environment string

const (
	ProductionEnvironment Environment = "production"
	StagingEnvironment    Environment = "staging"
	LocalEnvironment      Environment = "local"
)

func (e Environment) String() string {
	return string(e)
}

var AllEnvironments = []Environment{
	ProductionEnvironment, StagingEnvironment, LocalEnvironment,
}

func (e Environment) Validate() bool {
	for _, v := range AllEnvironments {
		if v == e {
			return true
		}
	}

	return false
}
