package ceph

import "testing"

func TestNewLibRados(t *testing.T) {
	librados := NewLibRados("ceph", "client.admin")
	err := librados.Rados_create2(0)
	if err != nil {
		t.Fatalf(err.Error())
	}
	t.Log("created a cluster handle")
	err = librados.Rados_conf_read_file("/etc/ceph/ceph.conf")
	if err != nil {
		t.Fatalf(err.Error())
	}
	t.Log("read config file")
	err = librados.Rados_connect()
	if err != nil {
		t.Fatalf(err.Error())
	}
	t.Log("connect to the cluster")

	err = librados.Rados_ioctx_create("data")
	if err != nil {
		t.Fatalf(err.Error())
	}
	t.Log("create io context")
	err = librados.Rados_write_full("greeting", []byte("hello"))
	if err != nil {
		t.Fatalf(err.Error())
	}


	t.Fatalf("com")
}
