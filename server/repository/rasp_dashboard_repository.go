package repository

type IRaspDashboardRepository interface {
}

type RaspDashboardRepository struct {
}

func NewRaspDashboardRepository() IRaspDashboardRepository {
	return RaspDashboardRepository{}
}
