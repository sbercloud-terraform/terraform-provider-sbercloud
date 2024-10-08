// Copyright 2020 Huawei Technologies Co.,Ltd.
//
// Licensed to the Apache Software Foundation (ASF) under one
// or more contributor license agreements.  See the NOTICE file
// distributed with this work for additional information
// regarding copyright ownership.  The ASF licenses this file
// to you under the Apache License, Version 2.0 (the
// "License"); you may not use this file except in compliance
// with the License.  You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing,
// software distributed under the License is distributed on an
// "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
// KIND, either express or implied.  See the License for the
// specific language governing permissions and limitations
// under the License.

package httphandler

import (
	"net/http"
	"time"
)

type MonitorMetric struct {
	Host          string
	Path          string
	Method        string
	Raw           string
	UserAgent     string
	RequestId     string
	StatusCode    int
	ContentLength int64
	Latency       time.Duration
	Attributes    map[string]interface{}
}

type HttpHandler struct {
	RequestHandlers  func(http.Request)
	ResponseHandlers func(http.Response)
	MonitorHandlers  func(*MonitorMetric)
}

func NewHttpHandler() *HttpHandler {
	handler := HttpHandler{}
	return &handler
}

func (handler *HttpHandler) AddRequestHandler(requestHandler func(http.Request)) *HttpHandler {
	handler.RequestHandlers = requestHandler
	return handler
}

func (handler *HttpHandler) AddResponseHandler(responseHandler func(response http.Response)) *HttpHandler {
	handler.ResponseHandlers = responseHandler
	return handler
}

func (handler *HttpHandler) AddMonitorHandler(monitorHandler func(*MonitorMetric)) *HttpHandler {
	handler.MonitorHandlers = monitorHandler
	return handler
}
