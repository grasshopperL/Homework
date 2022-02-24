/**
 * @Author: liubaoshuai3
 * @Date: 2022/2/23 15:20
 * @File: main
 * @Description:
 */

package main

import (
	"homework/myhttpserver"
	"net/http"
	"os"
)

func main() {
	_ = os.Setenv("test_env_one", "test_env_one")
	_ = os.Setenv("test_env_two", "test_env_two")
	ts := myhttpserver.MyServer{Port:":8888"}
	tf := myhttpserver.MyRouteFunc{}
	rm := make(map[string]func(http.ResponseWriter, *http.Request))
	rm["/hello"] = tf.GetBaseInfo
	rm["/healthZ"] = tf.HealthCheck
	ts.CreateMyHttpServer(rm)
}