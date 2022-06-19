package base

type DAOBase interface {
	Create()
	Delete()
	Update()
	FindOne()
	FindMany()
}
