/*
 * Licensed to the Apache Software Foundation (ASF) under one or more
 * contributor license agreements.  See the NOTICE file distributed with
 * this work for additional information regarding copyright ownership.
 * The ASF licenses this file to You under the Apache License, Version 2.0
 * (the "License"); you may not use this file except in compliance with
 * the License.  You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package main

import (
	"context"
	"dubbo.apache.org/dubbo-go/v3/common/constant"
	"time"
)

import (
	hessian "github.com/apache/dubbo-go-hessian2"
	"github.com/apache/dubbo-go-samples/multi-zone/go-client/pkg"
	"github.com/dubbogo/gost/log"
)

import (
	_ "dubbo.apache.org/dubbo-go/v3/cluster/cluster_impl"
	_ "dubbo.apache.org/dubbo-go/v3/cluster/loadbalance"
	_ "dubbo.apache.org/dubbo-go/v3/common/proxy/proxy_factory"
	"dubbo.apache.org/dubbo-go/v3/config"
	_ "dubbo.apache.org/dubbo-go/v3/filter/filter_impl"
	_ "dubbo.apache.org/dubbo-go/v3/protocol/dubbo"
	_ "dubbo.apache.org/dubbo-go/v3/registry/protocol"
	_ "dubbo.apache.org/dubbo-go/v3/registry/zookeeper"
)

var userProvider = new(pkg.UserProvider)

func init() {
	config.SetConsumerService(userProvider)
	hessian.RegisterPOJO(&pkg.User{})
}

// need to setup environment variable "CONF_CONSUMER_FILE_PATH" to "conf/client.yml" before run
func main() {
	hessian.RegisterPOJO(&pkg.User{})
	config.Load()
	time.Sleep(3 * time.Second)

	gxlog.CInfo("\n\n\nstart to test dubbo")

	// context zone hangzhou,
	ctx := context.Background()
	// if set zoneForce, must have zone tag
	ctx = context.WithValue(ctx, constant.REGISTRY_KEY+"."+constant.ZONE_FORCE_KEY, true)

	var hz, sh int
	loop := 50
	user := &pkg.User{}
	for i := 0; i < loop; i++ {
		err := userProvider.GetUser(ctx, []interface{}{i}, user)
		if err != nil {
			panic(err)
		}
		if "dev-hz" == user.ID {
			hz++
		}
		if "dev-sh" == user.ID {
			sh++
		}
		gxlog.CInfo("response %d result: %v\n", i, user)
	}

	gxlog.CInfo("loop count : %d, hangzhou count : %d, shanghai count : %d", loop, hz, sh)
}
