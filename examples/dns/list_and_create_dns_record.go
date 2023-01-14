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
	client := active24.New("123456qwerty-ok", active24.ApiEndpoint("https://sandboxapi.active24.com"))

	dns := client.Dns()
	//list all domains
	list, err := dns.List()
	if err != nil {
		panic(err)
	}
	for _, d := range list {
		print(d)
	}

	//list DNS records in domain example.com
	recs, err := dns.With("example.com").List()
	if err != nil {
		panic(err)
	}
	for _, rec := range recs {
		fmt.Printf("rec[type:%s, name:%s, ttl:%d]\n", *rec.Type, rec.Name, rec.Ttl)
	}

	//create CNAME record
	alias := "host1"
	err = dns.With("example.com").Create(active24.DnsRecordTypeA, &active24.DnsRecord{
		Alias: &alias,
	})
	if err != nil {
		panic(err)
	}

}
