package ceph

type LibRBD interface {
	LibRados
}

type libRBD struct {
	libRados


}

func NewLibRBD() *libRBD {
	return &libRBD{

	}
}


