/*
 * Copyright (C) 2015 Red Hat, Inc.
 *
 * Licensed to the Apache Software Foundation (ASF) under one
 * or more contributor license agreements.  See the NOTICE file
 * distributed with this work for additional information
 * regarding copyright ownership.  The ASF licenses this file
 * to you under the Apache License, Version 2.0 (the
 * "License"); you may not use this file except in compliance
 * with the License.  You may obtain a copy of the License at
 *
 *  http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on an
 * "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
 * KIND, either express or implied.  See the License for the
 * specific language governing permissions and limitations
 * under the License.
 *
 */

package storage

import (
	"errors"
	"fmt"

	"github.com/skydive-project/skydive/common"
	"github.com/skydive-project/skydive/config"
	"github.com/skydive-project/skydive/etcd"
	"github.com/skydive-project/skydive/filters"
	"github.com/skydive-project/skydive/flow"
	"github.com/skydive-project/skydive/flow/storage/elasticsearch"
	"github.com/skydive-project/skydive/flow/storage/orientdb"
	"github.com/skydive-project/skydive/logging"
)

// ErrNoStorageConfigured error no storage has been configured
var (
	ErrNoStorageConfigured = errors.New("No storage backend has been configured")
)

// Storage interface a flow storage mechanism
type Storage interface {
	Start()
	StoreFlows(flows []*flow.Flow) error
	SearchFlows(fsq filters.SearchQuery) (*flow.FlowSet, error)
	SearchMetrics(fsq filters.SearchQuery, metricFilter *filters.Filter) (map[string][]common.Metric, error)
	SearchRawPackets(fsq filters.SearchQuery, packetFilter *filters.Filter) (map[string]*flow.RawPackets, error)
	Stop()
}

// NewStorage creates a new flow storage based on the backend
func NewStorage(backend string, etcdClient *etcd.Client) (s Storage, err error) {
	driver := config.GetString("storage." + backend + ".driver")
	switch driver {
	case "elasticsearch":
		s, err = elasticsearch.New(backend, etcdClient)
		if err != nil {
			err = fmt.Errorf("Can't connect to ElasticSearch server: %v", err)
			return
		}
	case "orientdb":
		s, err = orientdb.New(backend)
		if err != nil {
			err = fmt.Errorf("Can't connect to OrientDB server: %v", err)
			return
		}
	case "memory":
		return
	default:
		err = fmt.Errorf("Flow backend driver '%s' not supported", driver)
		return
	}

	logging.GetLogger().Infof("Using %s as storage", backend)
	return
}

// NewStorageFromConfig creates a new storage based configuration
func NewStorageFromConfig(etcdClient *etcd.Client) (s Storage, err error) {
	return NewStorage(config.GetString("analyzer.flow.backend"), etcdClient)
}
