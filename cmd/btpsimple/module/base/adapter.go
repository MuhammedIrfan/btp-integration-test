/*
 * Copyright 2021 ICON Foundation
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package base

import "fmt"

var Clients = map[string]Client{}

func RegisterClients(networks []string, client Client) {
	for _, network := range networks {
		Clients[network] = client
	}
}

func GetClient(network string) (Client, error) {
	if c := Clients[network]; c != nil {
		return c, nil
	}
	return nil, fmt.Errorf("not supported blockchain:%s", network)
}
