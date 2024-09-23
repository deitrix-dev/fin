package fin

type ServiceClient interface {
	Store
}

type Service struct {
	Store
}

func NewService(store Store) *Service {
	return &Service{Store: store}
}
