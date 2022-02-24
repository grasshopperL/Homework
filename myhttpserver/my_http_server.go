/**
 * @Author: liubaoshuai3
 * @Date: 2022/2/23 15:21
 * @File: my_http_server
 * @Description:
 */

package myhttpserver

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
)

type MyServer struct {
	Port string
	mux *http.ServeMux
}

type MyRouteFunc struct {}

type Res struct {
	Status int
	Data interface{}
}

// create a http server with port and
func (m *MyServer) CreateMyHttpServer(route map[string]func(http.ResponseWriter, *http.Request)) {
	m.mux = http.NewServeMux()
	for k, v := range route {
		m.mux.HandleFunc(k, v)
	}
	err := http.ListenAndServe(m.Port, m.mux)
	if err != nil {
		log.Fatal(err)
	}
}

// get base info about one http
func (m *MyRouteFunc) GetBaseInfo(w http.ResponseWriter, r *http.Request) {
	rHeader := make(map[string]interface{})
	data := make(map[string]interface{})
	resp := Res{
		Status: 0,
		Data:   nil,
	}
	for k, v := range r.Header {
		rHeader[k] = strings.Join(v, ",")
	}
	w.WriteHeader(200)
	//TODO if change add to set, maybe influence the origin http res header or not
	for k, v := range rHeader {
		w.Header().Add(k, v.(string))
	}
	for k, v := range getEnv(){
		w.Header().Add(k, v)
	}
	w.Header().Add("lbs", "lbs")
	fmt.Print(w.Header())
	data["message"] = "Success"
	resp.Status = http.StatusCreated
	resp.Data = data
	jsonRep, err := json.Marshal(resp)
	if err != nil {
		log.Fatal(err)
	}
	_, err = w.Write(jsonRep)
	if err != nil {
		log.Fatal(err)
	}
	serverOutput(&resp, r)
	return
}

// get the health status about application
func(m *MyRouteFunc) HealthCheck(w http.ResponseWriter, r *http.Request)  {
	data := make(map[string]interface{})
	var outBytes bytes.Buffer
	cmdStr := fmt.Sprintf("netstat -ano -p tcp | findstr %d", 8888)
	cmd := exec.Command("cmd", "/c", cmdStr)
	cmd.Stdout = &outBytes
	_ = cmd.Run()
	resStr := outBytes.String()
	re := regexp.MustCompile(`\s\d+\s`).FindAllString(resStr, -1)
	if len(re) > 0 {
		pid, err := strconv.Atoi(strings.TrimSpace(re[0]))
		if err != nil {
			data["port"] = 0
			data["health"] = "Unhealthy"
		} else {
			data["port"] = pid
			data["health"] = "Healthy"
		}
	}
	resp := Res{
		Status: 0,
		Data:   nil,
	}
	resp.Status = 200
	resp.Data = data
	jsonRep, err := json.Marshal(resp)
	if err != nil {
		log.Fatal(err)
	}
	_, err = w.Write(jsonRep)
	if err != nil {
		log.Fatal(err)
	}
	return
}

// output client ip and http code to standard
func serverOutput(res *Res, r *http.Request) {
	rIp := r.Header.Get("X-Real-Ip")
	if rIp == "" {
		rIp = r.Header.Get("X-Forwarded-For")
	}
	if rIp == "" {
		rIp = r.RemoteAddr
	}
	fmt.Print("remote ip is:", rIp)
	fmt.Print("status code is:", res.Status)
}

 // get all system environment variables
func getEnv() map[string]string {
	defer func() {
		_ = os.Unsetenv("test_env_one")
		_ = os.Unsetenv("test_env_two")
	}()
	em := make(map[string]string)
	el := os.Environ()
	for _, e := range el {
		em[e] = e
	}
	return em
}





