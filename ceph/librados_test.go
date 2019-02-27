package ceph

import "testing"

func TestNewLibRados(t *testing.T) {
	librados := NewLibRados([]byte("ceph"), []byte("client.admin"))
	err := librados.Rados_create2(0)
	if err != nil {
		t.Fatalf(err.Error())
	}
	t.Log("create a cluster handle")
	err = librados.Rados_conf_read_file([]byte("/etc/ceph/ceph.conf"))
	if err != nil {
		t.Fatalf(err.Error())
	}
	t.Log("read config file")
	err = librados.Rados_connect()
	if err != nil {
		t.Fatalf(err.Error())
	}
	t.Log("connect to the cluster")

}
