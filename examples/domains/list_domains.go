/*
Copyright 2023 Richard Kosegi

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package main

import (
	"fmt"
	"github.com/rkosegi/active24-go/active24"
)

func main() {
	client := active24.New(active24.SandboxToken, active24.ApiEndpoint("https://sandboxapi.active24.com"))
	dom := client.Domains()
	//list all domains
	list, err := dom.List()
	if err != nil {
		panic(err)
	}
	for _, d := range list {
		fmt.Printf("Name: %s, Status: %s, Holder: %s", d.Name, d.Status, d.HolderFullName)
	}
}
