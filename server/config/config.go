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

package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

var DefConfig = Config{
	RestPort:           20334,
	Version:            "",
	HttpMaxConnections: 10000,
	HttpCertPath:       "",
	HttpKeyPath:        "",
}

type Config struct {
	LogLevel           int     `json:"log_level"`
	RestPort           uint    `json:"rest_port"`
	Version            string  `json:"version"`
	HttpMaxConnections int     `json:"http_max_connections"`
	HttpCertPath       string  `json:"http_cert_path"`
	HttpKeyPath        string  `json:"http_key_path"`
	ProjectDBUrl       string  `json:"explorerdb_url"`
	ProjectDBUser      string  `json:"explorerdb_user"`
	ProjectDBPassword  string  `json:"explorerdb_password"`
	ProjectDBName      string  `json:"explorerdb_name"`
}

func InitConfig() error {
	file, err := os.Open("./config.json")
	if err != nil {
		return err
	}
	bs, err := ioutil.ReadAll(file)
	if err != nil {
		return err
	}
	cfg := Config{}
	err = json.Unmarshal(bs, &cfg)
	if err != nil {
		return err
	}
	if cfg.RestPort == 0 {
		return fmt.Errorf("not config the rest port")
	}
	DefConfig = cfg
	return nil
}
