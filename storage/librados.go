package storage

// 基础库

/*
#cgo LDFLAGS: -lrados
#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <rados/librados.h>
*/
import "C"

import (
	"bytes"
	"encoding/binary"
	"errors"
	"unsafe"
)

type LibRados interface {
	Rados_version(major *int, minor *int, extra *int) string
	Rados_create2(flags uint32) error
	Rados_conf_read_file(path []byte) error
	Rados_connect() error
	Rados_ioctx_create(pool_name []byte, io C.rados_ioctx_t) error
	Rados_write(key []byte, value []byte, offset uint) error
	Rados_ioctx_destroy()
	Rados_shutdown()
	Rados_setxattr(object_name []byte, attr_name []byte, value []byte) error
	Rados_getxattr(object_name []byte, attr_name []byte, size uint) (interface{}, error)
}

type libRados struct {
	cluster C.rados_t // 集群句柄

	cluster_name []byte // 集群名称
	user_name []byte    // 用户名

	pool_name []byte    // 对象池

	io C.rados_ioctx_t
}

func NewLibRados(cluster_name []byte, user_name []byte) *libRados {
	return &libRados{
		cluster_name: cluster_name,
		user_name: user_name,
	}
}

// 创建集群句柄
func (lib *libRados) Rados_create2(flags uint32) error {
	err := C.rados_create2(&lib.cluster, (*C.char)(unsafe.Pointer(&lib.cluster_name)), (*C.char)(unsafe.Pointer(&lib.user_name)), (C.ulong)(flags))
	if err < 0 {
		return errors.New("Couldn't create the ceph cluster handle! ")
	}
	return nil
}

// 读取配置文件
func (lib *libRados) Rados_conf_read_file(path []byte) error {
	err := C.rados_conf_read_file(lib.cluster, (*C.char)(unsafe.Pointer(&path)))
	if err < 0 {
		return errors.New("Cannot read config file: " + string(path))
	}
	return nil
}

// 连接
func (lib *libRados) Rados_connect() error {
	err := C.rados_connect(lib.cluster)
	if err < 0 {
		return errors.New("cannot connect to cluster")
	}

	return nil
}

// 创建io上下文
func (lib *libRados) Rados_ioctx_create(pool_name []byte, io C.rados_ioctx_t) error {
	lib.pool_name = pool_name
	err := C.rados_ioctx_create(lib.cluster, (*C.char)(unsafe.Pointer(&lib.pool_name)), &lib.io)
	if err < 0 {
		return errors.New("cannot open rados pool[" + string(pool_name) + "]")
	}

	return nil
}

// 写入数据
func (lib *libRados) Rados_write(key []byte, value []byte, offset uint) error {
	buf := new(bytes.Buffer)
	_ = binary.Write(buf, binary.LittleEndian, value)
	err := C.rados_write(lib.io, (*C.char)(unsafe.Pointer(&key)), (*C.char)(unsafe.Pointer(&buf)), (C.ulong)(len(buf.Bytes())), (C.size_t)(offset))
	if err < 0 {
		return errors.New("cannot write object to pool[" + string(lib.pool_name) + "]")
	}

	return nil
}

// 销毁io上下文
func (lib *libRados) Rados_ioctx_destroy() {
	C.rados_ioctx_destroy(lib.io)
}

// 关闭集群句柄
func (lib *libRados) Rados_shutdown() {
	C.rados_shutdown(lib.cluster)
}

// 设置属性值
func (lib *libRados) Rados_setxattr(object_name []byte, attr_name []byte, value []byte) error {
	buf := new(bytes.Buffer)
	_ = binary.Write(buf, binary.LittleEndian, value)
	err := C.rados_setxattr(lib.io, (*C.char)(unsafe.Pointer(&object_name)), (*C.char)(unsafe.Pointer(&attr_name)), (*C.char)(unsafe.Pointer(&buf)), (C.ulong)(len(buf.Bytes())))
	if err < 0 {
		return errors.New("cannot set extended attribute on object[" + string(object_name) + "]")
	}
	return nil
}

func (lib *libRados) Rados_getxattr(object_name []byte, attr_name []byte, size uint) (interface{}, error) {
	buf := new(bytes.Buffer)
	err := C.rados_getxattr(lib.io, (*C.char)(unsafe.Pointer(&object_name)), (*C.char)(unsafe.Pointer(&attr_name)), (*C.char)(unsafe.Pointer(&buf)), (C.size_t)(size))
	if err < 0 {
		return nil, errors.New("cannot get extended attribute on object[" + string(object_name) + "]")
	}
	return buf, nil
}

// 获取版本号
func (lib *libRados) Rados_version(major *int, minor *int, extra *int) string {
	return ""
}
