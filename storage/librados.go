package storage


/*
#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <rados/librados.h>


*/
import "C"

import (
	"errors"
	"bytes"
	"encoding/binary"
)

type LibRados interface {
	Rados_version(major *int, minor *int, extra *int) string
	Rados_create2(flags C.uint64_t) (int, error)
	Rados_conf_read_file(path *C.char) (int, error)
	Rados_connect() (int, error)
	Rados_ioctx_create(pool_name *C.char, io C.rados_ioctx_t) (int, error)
	Rados_write(key string, value string, offset C.unit64_t) (int, error)
	Rados_ioctx_destroy()
	Rados_shutdown()
	Rados_setxattr(object_name *C.char, attr_name *C.char, value string) (int, error)
	Rados_getxattr(object_name *C.char, attr_name *C.char, size C.uint) (interface{}, error)
}

type libRados struct {
	cluster *C.rados_t // 集群句柄

	cluster_name *C.char // 集群名称
	user_name *C.char    // 用户名

	pool_name *C.char    // 对象池

	io C.rados_ioctx_t
}

func NewLibRados(cluster_name string, user_name string) LibRados {
	return &libRados{
		cluster: -1,
		cluster_name: cluster_name,
		user_name: user_name,
	}
}

// 创建集群句柄
func (lib *libRados) Rados_create2(flags uint64) (int, error) {
	err := C.rados_create2(lib.cluster, lib.cluster_name, lib.user_name, flags)
	if err < 0 {
		return err, errors.New("Couldn't create the ceph cluster handle! ")
	}
	return err, nil
}

// 读取配置文件
func (lib *libRados) Rados_conf_read_file(path *C.char) (int, error) {
	err := C.rados_conf_read_file(lib.cluster, path)
	if err < 0 {
		return err, errors.New("Cannot read config file: " + path)
	}
	return err, nil
}

// 连接
func (lib *libRados) Rados_connect() (int, error) {
	err := C.rados_connect(lib.cluster)
	if err < 0 {
		return err, errors.New("cannot connect to cluster")
	}

	return err, nil
}

// 创建io上下文
func (lib *libRados) Rados_ioctx_create(pool_name *C.char, io C.rados_ioctx_t) (int, error) {
	lib.pool_name = pool_name
	err := C.rados_ioctx_create(lib.cluster, pool_name, &lib.io)
	if err < 0 {
		return err, errors.New("cannot open rados pool[" + lib.pool_name + "]")
	}

	return err, nil
}

// 写入数据
func (lib *libRados) Rados_write(key string, value string, offset C.unit64_t) (int, error) {
	buf := new(bytes.Buffer)
	_ = binary.Write(buf, binary.LittleEndian, value)
	err := C.rados_write(lib.io, key, buf, buf.Bytes(), offset)
	if err < 0 {
		return err, errors.New("cannot write object to pool[" + lib.pool_name + "]")
	}

	return err, nil
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
func (lib *libRados) Rados_setxattr(object_name *C.char, attr_name *C.char, value string) (int, error) {
	buf := new(bytes.Buffer)
	_ = binary.Write(buf, binary.LittleEndian, value)
	err := C.rados_setxattr(lib.io, object_name, attr_name, buf, buf.Bytes())
	if err < 0 {
		return err, errors.New("cannot set extended attribute on object[" + object_name + "]")
	}
	return err, nil
}

func (lib *libRados) Rados_getxattr(object_name *C.char, attr_name *C.char, size C.uint) (interface{}, error) {
	buf := new(bytes.Buffer)
	err := C.rados_getxattr(lib.io, object_name, attr_name, buf, size)
	if err < 0 {
		return nil, errors.New("cannot get extended attribute on object[" + object_name + "]")
	}
	return buf, nil
}

// 获取版本号
func (lib *libRados) Rados_version(major *int, minor *int, extra *int) string {
	return ""
}
