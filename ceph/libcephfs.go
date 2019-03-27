package ceph

// 文件系统
type Libcephfs interface {
	LibRados
}

type libcephfs struct {
	libRados
}

func NewLibCephfs() *libcephfs {
	return &libcephfs{}
}
