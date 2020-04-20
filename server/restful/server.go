/*
 * Copyright (C) 2020 The ontology Authors
 * This file is part of The ontology library.
 *
 * The ontology is free software: you can redistribute it and/or modify
 * it under the terms of the GNU Lesser General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 *
 * The ontology is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU Lesser General Public License for more details.
 *
 * You should have received a copy of the GNU Lesser General Public License
 * along with The ontology.  If not, see <http://www.gnu.org/licenses/>.
 */

package restful

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"github.com/ontio/layer2/server/config"
	"github.com/ontio/layer2/server/core"
	"github.com/ontio/layer2/server/log"
	"golang.org/x/net/netutil"
	"io"
	"net"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"
)

const TLS_PORT int = 443
const MAX_REQUEST_BODY_SIZE = 1 << 20

type handler func(map[string]interface{}) map[string]interface{}
type Action struct {
	sync.RWMutex
	name    string
	handler handler
}
type restServer struct {
	router   *Router
	listener net.Listener
	server   *http.Server
	postMap  map[string]Action //post method map
	getMap   map[string]Action //get method map
}

const (
	GET_LAYER2TX    = "/api/v1/getlayer2tx/:address"
	GET_LAYER2DEPOSIT    = "/api/v1/getlayer2deposit/:address"
	GET_LAYER2WITHDRAW    = "/api/v1/getlayer2withdraw/:address"
)

//init restful server
func InitRestServer() *restServer {
	rt := &restServer{}

	rt.router = NewRouter()
	rt.registryMethod()
	rt.initGetHandler()
	rt.initPostHandler()
	return rt
}

//start server
func (this *restServer) Start() error {
	retPort := int(config.DefConfig.RestPort)
	if retPort == 0 {
		log.Fatal("Not configure HttpRestPort port ")
		return nil
	}
	log.Infof("restful port: %d", retPort)
	tlsFlag := false
	if tlsFlag || retPort%1000 == TLS_PORT {
		var err error
		this.listener, err = this.initTlsListen()
		if err != nil {
			log.Error("Https Cert: ", err.Error())
			return err
		}
	} else {
		var err error
		this.listener, err = net.Listen("tcp", ":"+strconv.Itoa(retPort))
		if err != nil {
			log.Fatal("net.Listen: ", err.Error())
			return err
		}
	}
	this.server = &http.Server{Handler: this.router}
	//set LimitListener number
	if config.DefConfig.HttpMaxConnections > 0 {
		this.listener = netutil.LimitListener(this.listener, config.DefConfig.HttpMaxConnections)
	}
	err := this.server.Serve(this.listener)

	if err != nil {
		log.Fatal("ListenAndServe: ", err.Error())
		return err
	}

	return nil
}

//resigtry handler method
func (this *restServer) registryMethod() {
	getMethodMap := map[string]Action{
		GET_LAYER2TX:  {name: "getlayer2tx", handler: GetLayer2Tx},
		GET_LAYER2DEPOSIT:  {name: "getlayer2deposit", handler: GetLayer2Deposit},
		GET_LAYER2WITHDRAW:  {name: "getlayer2withdraw", handler: GetLayer2Withdraw},
	}

	// todo
	postMethodMap := map[string]Action{
	}
	this.postMap = postMethodMap
	this.getMap = getMethodMap
}

func (this *restServer) getPath(url string) string {
	if strings.Contains(url, strings.TrimRight(GET_LAYER2TX, ":address")) {
		return GET_LAYER2TX
	}
	if strings.Contains(url, strings.TrimRight(GET_LAYER2DEPOSIT, ":address")) {
		return GET_LAYER2DEPOSIT
	}
	if strings.Contains(url, strings.TrimRight(GET_LAYER2WITHDRAW, ":address")) {
		return GET_LAYER2WITHDRAW
	}
	return url
}

//get request params
func (this *restServer) getParams(r *http.Request, url string, req map[string]interface{}) map[string]interface{} {
	switch url {
	case GET_LAYER2TX:
		req["address"] = getParam(r, "address")
	case GET_LAYER2DEPOSIT:
		req["address"] = getParam(r, "address")
	case GET_LAYER2WITHDRAW:
		req["address"] = getParam(r, "address")
	default:
	}
	return req
}

//init get handler
func (this *restServer) initGetHandler() {

	for k, _ := range this.getMap {
		this.router.Get(k, func(w http.ResponseWriter, r *http.Request) {

			var req = make(map[string]interface{})
			var resp map[string]interface{}

			url := this.getPath(r.URL.Path)
			if h, ok := this.getMap[url]; ok {
				req = this.getParams(r, url, req)
				resp = h.handler(req)
				resp["action"] = h.name
			} else {
				resp = ResponsePack(core.REST_METHOD_INVALID)
			}
			this.response(w, resp)
		})
	}
}

//init post handler
func (this *restServer) initPostHandler() {
	for k, _ := range this.postMap {
		this.router.Post(k, func(w http.ResponseWriter, r *http.Request) {
			decoder := json.NewDecoder(io.LimitReader(r.Body, MAX_REQUEST_BODY_SIZE))
			defer r.Body.Close()
			var req = make(map[string]interface{})
			var resp map[string]interface{}
			url := this.getPath(r.URL.Path)
			if h, ok := this.postMap[url]; ok {
				if err := decoder.Decode(&req); err == nil {
					req = this.getParams(r, url, req)
					req["ip"] = r.RemoteAddr
					resp = h.handler(req)
					resp["action"] = h.name
				} else {
					resp = ResponsePack(core.REST_ILLEGAL_DATAFORMAT)
					resp["action"] = h.name
				}
			} else {
				resp = ResponsePack(core.REST_METHOD_INVALID)
			}
			this.response(w, resp)
		})
	}
	//Options
	for k, _ := range this.postMap {
		this.router.Options(k, func(w http.ResponseWriter, r *http.Request) {
			this.write(w, []byte{})
		})
	}

}
func (this *restServer) write(w http.ResponseWriter, data []byte) {
	w.Header().Add("Access-Control-Allow-Headers", "Content-Type")
	w.Header().Set("content-type", "application/json;charset=utf-8")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Write(data)
}

//response
func (this *restServer) response(w http.ResponseWriter, resp map[string]interface{}) {
	data, err := json.Marshal(resp)
	if err != nil {
		log.Fatal("HTTP Handle - json.Marshal: %v", err)
		return
	}
	this.write(w, data)
}

//stop restful server
func (this *restServer) Stop() {
	if this.server != nil {
		this.server.Shutdown(context.Background())
		log.Error("Close restful ")
	}
}

//restart server
func (this *restServer) Restart(cmd map[string]interface{}) map[string]interface{} {
	go func() {
		time.Sleep(time.Second)
		this.Stop()
		time.Sleep(time.Second)
		go this.Start()
	}()

	var resp = ResponsePack(core.SUCCESS)
	return resp
}

//init tls
func (this *restServer) initTlsListen() (net.Listener, error) {

	certPath := config.DefConfig.HttpCertPath
	keyPath := config.DefConfig.HttpKeyPath

	// load cert
	cert, err := tls.LoadX509KeyPair(certPath, keyPath)
	if err != nil {
		log.Error("load keys fail", err)
		return nil, err
	}

	tlsConfig := &tls.Config{
		Certificates: []tls.Certificate{cert},
	}

	restPort := strconv.Itoa(int(config.DefConfig.RestPort))
	log.Info("TLS listen port is ", restPort)
	listener, err := tls.Listen("tcp", ":"+restPort, tlsConfig)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	return listener, nil
}
