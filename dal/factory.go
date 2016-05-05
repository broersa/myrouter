package dal

type (
	// Factory Interface for Business Layer
	Factory interface {
		GetInstance() (Dal, error)
	}
)
