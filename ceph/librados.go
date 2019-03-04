package ceph

// 基础库

/*
#cgo LDFLAGS: -lrados
#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <rados/librados.h>
#include <rados/rados_types.h>
*/
import "C"

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"unsafe"
)

type LibRados interface {
	// setup and teardown
	Rados_create2(flags uint64) error
	Rados_create() error
	//Rados_ping_monitor(monId string, out *bytes.Buffer) error
	Rados_connect() error
	Rados_shutdown()

	// configure
	Rados_conf_read_file(path string) error


	Rados_write_full(object_name string, value []byte) error
	Rados_write(key string, value []byte, offset uint) error

	Rados_setxattr(object_name string, attr_name string, value []byte) error
	Rados_getxattr(object_name string, attr_name string, size uint) (interface{}, error)
	Rados_rmxattr(object_name string, attr_name string) error


	// 异步IO
	Rados_aio_create_completion() error
	Rados_aio_write(key string, value []byte, offset uint64) error
	Rados_aio_read() (interface{}, error)
	Rados_aio_write_full() error
	Rados_aio_append() error
	Rados_aio_release()
	Rados_aio_wait_for_complete() // in memory
	Rados_aio_wait_for_safe() // on disk
	Rados_aio_flush() error
	Rados_aio_flush_async() error
	Rados_aio_cancel() error
	Rados_aio_is_complete() error

	// Mon/OSD/PG commands
	Rados_mon_command(cmd string, params []byte) (out []byte, err error)
	//Rados_mgr_command(cmd string, params []byte) (out []byte, err error)
	Rados_buffer_free()
	Rados_osd_command(osdId int, cmd string, params []byte) (out []byte, err error)
	Rados_pg_command(pgstr string, cmd string, params []byte) (out []byte, err error)
	Rados_monitor_log() error
	Rados_monitor_log2() error

	// Pools
	Rados_pool_list() (out []byte, err error)
	Rados_pool_delete(pool_name string) error
	Rados_ioctx_create(pool_name string) error
	Rados_ioctx_destroy()

	// Snapshots
	Rados_ioctx_snap_create(snapname string) error
	Rados_ioctx_snap_remove(snapname string) error
	Rados_ioctx_snap_rollback(snapname string, key string) error
	Rados_rollback(snapname string, key string) error
	Rados_ioctx_snap_set_read(snap C.rados_snap_t) error
	Rados_ioctx_snap_list(snaps *C.rados_snap_t) (out []byte, err error)

	// Functions
	Rados_cluster_stat() error
	Rados_version() (major int, minor int, extra int)
	Rados_stat()
	Rados_lock_shared() error
	Raods_unlock() error
}

type libRados struct {
	cluster C.rados_t // 集群句柄

	cluster_name string // 集群名称
	user_name string    // 用户名

	config C.rados_config_t // 上下文配置

	pool_name string    // 对象池

	snapshot C.rados_snap_t // 快照

	io C.rados_ioctx_t // 同步IO上下文

	comp C.rados_completion_t // 异步IO

	Stat C.struct_rados_cluster_stat_t //
}

func NewLibRados(cluster_name string, user_name string) *libRados {
	return &libRados{
		cluster_name: cluster_name,
		user_name: user_name,
	}
}

// 创建集群句柄
func (lib *libRados) Rados_create2(flags uint64) error {
	err := C.rados_create2(&lib.cluster, (*C.char)(C.CString(lib.cluster_name)), (*C.char)(C.CString(lib.user_name)), (C.uint64_t)(flags))
	if int32(err) < 0 {
		return errors.New("Couldn't create the ceph cluster handle! " + fmt.Sprintf("%v", err))
	}
	return nil
}

func (lib *libRados) Rados_create() error  {
	err := C.rados_create(&lib.cluster, (*C.char)(nil))
	if int32(err) < 0 {
		return errors.New("Couldn't create the ceph cluster handle! " + fmt.Sprintf("%v", err))
	}
	return nil
}

// 读取配置文件
func (lib *libRados) Rados_conf_read_file(path string) error {
	err := C.rados_conf_read_file(lib.cluster, (*C.char)(C.CString(path)))
	if int32(err) < 0 {
		return errors.New("Cannot read config file: " + string(path) + " " + fmt.Sprintf("%v", err))
	}
	return nil
}

// 连接
func (lib *libRados) Rados_connect() error {
	err := C.rados_connect(lib.cluster)
	if int32(err) < 0 {
		return errors.New("cannot connect to cluster " + fmt.Sprintf("%v", err))
	}

	return nil
}


// 写入数据
func (lib *libRados) Rados_write(key string, value []byte, offset uint) error {
	buf := new(bytes.Buffer)
	_ = binary.Write(buf, binary.LittleEndian, value)
	err := C.rados_write(lib.io, (*C.char)(C.CString(key)), (*C.char)(unsafe.Pointer(&buf)), (C.uint64_t)(len(buf.Bytes())), (C.size_t)(offset))
	if int32(err) < 0 {
		return errors.New("cannot write object to pool[" + string(lib.pool_name) + "] " + fmt.Sprintf("%v", err))
	}

	return nil
}

// 关闭集群句柄
func (lib *libRados) Rados_shutdown() {
	C.rados_shutdown(lib.cluster)
}

// 创建io上下文
func (lib *libRados) Rados_ioctx_create(pool_name string) error {
	lib.pool_name = pool_name
	err := C.rados_ioctx_create(lib.cluster, (*C.char)(C.CString(lib.pool_name)), &lib.io)
	if int32(err) < 0 {
		return errors.New("cannot open rados pool[" + string(pool_name) + "] " + fmt.Sprintf("%v", err))
	}

	return nil
}

func (lib *libRados) Rados_write_full(object_name string, value []byte) error {
	buf := new(bytes.Buffer)
	_ = binary.Write(buf, binary.LittleEndian, value)
	err := C.rados_write_full(lib.io, (*C.char)(C.CString(object_name)), (*C.char)(unsafe.Pointer(&buf)), (C.ulong)(len(buf.Bytes())))
	if int32(err) < 0 {
		lib.Rados_ioctx_destroy()
		lib.Rados_shutdown()
		return errors.New("cannot write pool[" + string(object_name) + "] " + fmt.Sprintf("%v", err))
	}
	return nil
}

// 销毁io上下文
func (lib *libRados) Rados_ioctx_destroy() {
	C.rados_ioctx_destroy(lib.io)
}



// 设置属性值
func (lib *libRados) Rados_setxattr(object_name string, attr_name string, value []byte) error {
	buf := new(bytes.Buffer)
	_ = binary.Write(buf, binary.LittleEndian, value)
	err := C.rados_setxattr(lib.io, (*C.char)(C.CString(object_name)), (*C.char)(C.CString(attr_name)), (*C.char)(unsafe.Pointer(&buf)), (C.ulong)(len(buf.Bytes())))
	if int32(err) < 0 {
		return errors.New("cannot set extended attribute on object[" + string(object_name) + "] " + fmt.Sprintf("%v", err))
	}
	return nil
}

func (lib *libRados) Rados_getxattr(object_name string, attr_name string, size uint) (interface{}, error) {
	buf := new(bytes.Buffer)
	err := C.rados_getxattr(lib.io, (*C.char)(C.CString(object_name)), (*C.char)(C.CString(attr_name)), (*C.char)(unsafe.Pointer(&buf)), (C.size_t)(size))
	if int32(err) < 0 {
		return nil, errors.New("cannot get extended attribute on object[" + string(object_name) + "] " + fmt.Sprintf("%v", err))
	}
	return buf, nil
}

func (lib *libRados) Rados_rmxattr(object_name string, attr_name string) error {
	err := C.rados_rmxattr(lib.io, (*C.char)(C.CString(object_name)), (*C.char)(C.CString(attr_name)))
	if int32(err) < 0 {
		return errors.New("cannot remove extended attribute on object[" + string(object_name) + "] " + fmt.Sprintf("%v", err))
	}
	return nil
}

func (lib *libRados) Rados_aio_create_completion() error {
	err := C.rados_aio_create_completion(nil, nil, nil, &lib.comp)
	if int32(err) < 0 {
		lib.Rados_ioctx_destroy()
		lib.Rados_shutdown()
		return errors.New("cannot create aio completion " + fmt.Sprintf("%v", err))
	}
	return nil
}

func (lib *libRados) Rados_aio_write(key string, value []byte, offset uint64) error {
	buf := new(bytes.Buffer)
	_ = binary.Write(buf, binary.LittleEndian, value)
	err := C.rados_aio_write(lib.io, (*C.char)(C.CString(key)), lib.comp, (*C.char)(unsafe.Pointer(&buf)), (C.size_t)(len(buf.Bytes())), (C.uint64_t)(offset))
	if int32(err) < 0 {
		lib.Rados_aio_release()
		lib.Rados_ioctx_destroy()
		lib.Rados_shutdown()
		return errors.New("cannot not schedule aio write " + fmt.Sprintf("%v", err))
	}

	return nil
}

func (lib *libRados) Rados_aio_read() (interface{}, error) {

	return nil,nil
}

func (lib *libRados) Rados_aio_write_full() error {

	return nil
}

func (lib *libRados) Rados_aio_append() error {

	return nil
}

func (lib *libRados) Rados_aio_release() {
	C.rados_aio_release(lib.comp)
}

func (lib *libRados) Rados_aio_wait_for_complete() {
	C.rados_aio_wait_for_complete(lib.comp)
}

func (lib *libRados) Rados_aio_wait_for_safe() {
	C.rados_aio_wait_for_safe(lib.comp)
}

func (lib *libRados) Rados_aio_flush() error {
	err := C.rados_aio_flush(lib.io)
	if int32(err) < 0 {
		return errors.New("flush to disk error " + fmt.Sprintf("%v", err))
	}
	return nil
}

func (lib *libRados) Rados_aio_flush_async() error {
	err := C.rados_aio_flush_async(lib.io, lib.comp)
	if int32(err) < 0 {
		return errors.New("cannot not schedule flush to disk " + fmt.Sprintf("%v", err))
	}
	return nil
}

func (lib *libRados) Rados_aio_cancel() error {
	err := C.rados_aio_cancel(lib.io, lib.comp)
	if int32(err) < 0 {
		return errors.New("cannot not cancel async operation " + fmt.Sprintf("%v", err))
	}
	return nil
}

func (lib *libRados) Rados_aio_is_complete() error {
	err := C.rados_aio_is_complete(lib.comp)
	if int32(err) < 0 {
		return errors.New("not complete")
	}

	return nil
}


func (lib *libRados) Rados_mon_command(cmd string, params []byte) (out []byte, err error) {
	var outs (**C.char)
	var outlen (*C.size_t)

	err1 := C.rados_mon_command(lib.cluster,
		(**C.char)(unsafe.Pointer(&cmd)),
		(C.size_t)(len([]byte(cmd))),
		(*C.char)(unsafe.Pointer(&params)),
		(C.size_t)(len(params)),
		(**C.char)(unsafe.Pointer(&out)),
		(*C.size_t)(unsafe.Pointer(uintptr(len(out)))),
		outs,
		outlen)
	if int32(err1) < 0 {
		return out, errors.New("mon command execute fail " + fmt.Sprintf("%v", err1))
	}

	return out, nil
}

//func (lib *libRados) Rados_mgr_command(cmd string, params []byte) (out []byte, err error) {
//	var outs (**C.char)
//	var outlen (*C.size_t)
//
//	err1 := C.rados_mgr_command(lib.cluster,
//		(**C.char)(unsafe.Pointer(&cmd)),
//		(C.size_t)(len([]byte(cmd))),
//		(*C.char)(unsafe.Pointer(&params)),
//		(C.size_t)(len(params)),
//		(**C.char)(unsafe.Pointer(&out)),
//		(*C.size_t)(unsafe.Pointer(uintptr(len(out)))),
//		outs,
//		outlen)
//	if int32(err1) < 0 {
//		return out, errors.New("mgr command execute fail " + fmt.Sprintf("%v", err1))
//	}
//
//	return out, nil
//}

func (lib *libRados) Rados_buffer_free() {

}

func (lib *libRados) Rados_osd_command(osdId int, cmd string, params []byte) (out []byte, err error) {
	var outs (**C.char)
	var outlen (*C.size_t)

	err1 := C.rados_osd_command(lib.cluster,
		(C.int)(osdId),
		(**C.char)(unsafe.Pointer(&cmd)),
		(C.size_t)(len([]byte(cmd))),
		(*C.char)(unsafe.Pointer(&params)),
		(C.size_t)(len(params)),
		(**C.char)(unsafe.Pointer(&out)),
		(*C.size_t)(unsafe.Pointer(uintptr(len(out)))),
		outs,
		outlen)
	if int32(err1) < 0 {
		return out, errors.New("osd command execute fail " + fmt.Sprintf("%v", err1))
	}

	return out, nil
}

func (lib *libRados) Rados_pg_command(pgstr string, cmd string, params []byte) (out []byte, err error) {
	var outs (**C.char)
	var outlen (*C.size_t)

	err1 := C.rados_pg_command(lib.cluster,
		(*C.char)(unsafe.Pointer(&pgstr)),
		(**C.char)(unsafe.Pointer(&cmd)),
		(C.size_t)(len([]byte(cmd))),
		(*C.char)(unsafe.Pointer(&params)),
		(C.size_t)(len(params)),
		(**C.char)(unsafe.Pointer(&out)),
		(*C.size_t)(unsafe.Pointer(uintptr(len(out)))),
		outs,
		outlen)
	if int32(err1) < 0 {
		return out, errors.New("pg command execute fail " + fmt.Sprintf("%v", err1))
	}

	return out, nil
}

func (lib *libRados) Rados_monitor_log() error {

	return nil
}

func (lib *libRados) Rados_monitor_log2() error {
	return nil
}


func (lib *libRados) Rados_pool_list() (out []byte, err error) {

	return
}

func (lib *libRados) Rados_pool_delete(pool_name string) error {
	return nil
}

func (lib *libRados) Rados_ioctx_snap_create(snapname string) error {
	return nil
}

func (lib *libRados) Rados_ioctx_snap_remove(snapname string) error {
	return nil
}

func (lib *libRados) Rados_ioctx_snap_rollback(snapname string, key string) error {
	return nil
}

func (lib *libRados) Rados_rollback(snapname string, key string) error {
	return nil
}

func (lib *libRados) Rados_ioctx_snap_set_read(snap C.rados_snap_t) error {
	return nil
}

func (lib *libRados) Rados_ioctx_snap_list(snaps *C.rados_snap_t) (out []byte, err error) {
	return
}

func (lib *libRados) Rados_cluster_stat() error {
	err := C.rados_cluster_stat(lib.cluster, &lib.Stat)
	if int32(err) < 0 {
		return errors.New("read cluster stat error " + fmt.Sprintf("%v", err))
	}
	return nil
}

func (lib *libRados) Rados_version() (major int, minor int, extra int) {
	C.rados_version((*C.int)(unsafe.Pointer(&major)), (*C.int)(unsafe.Pointer(&minor)), (*C.int)(unsafe.Pointer(&extra)))
	return
}

func (lib *libRados) Rados_stat() {

}

func (lib *libRados) Rados_lock_shared() error {
	return nil
}

func (lib *libRados) Raods_unlock() error {
	return nil
}