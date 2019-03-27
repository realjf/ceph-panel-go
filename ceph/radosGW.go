package ceph

// 对象存储网关
type RadosGW interface {
	LibRados
}

type radosGW struct {
	libRados
}

func NewRadosGW() *radosGW {
	return &radosGW{}
}
