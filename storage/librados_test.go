package storage

import "testing"

func TestNewLibRados(t *testing.T) {
	librados := NewLibRados("ceph", "client.admin")
	_, err := librados.Rados_create2(0)
	if err != nil {
		t.Fatalf(err.Error())
	}
	t.Log("create a cluster handle")
	_, err = librados.Rados_conf_read_file("/etc/ceph/ceph.conf")
	if err != nil {
		t.Fatalf(err.Error())
	}
	t.Log("read config file")
	_, err = librados.Rados_connect()
	if err != nil {
		t.Fatalf(err.Error())
	}
	t.Log("connect to the cluster")

}
