# Active24.cz client in Go

This is client library to interact with [Active24 API](https://faq.active24.com/eng/739445-REST-API-for-developers?l=en-US).
Currently, only subset of API is implemented, but contributions are always welcome.

## Usage

```go
package main

import "github.com/rkosegi/active24-go/active24"

func main() {
	client := active24.New("my-secret-api-token")

	alias := "host1"
	_, err := client.Dns().With("example.com").Create(active24.DnsRecordTypeA, &active24.DnsRecord{
		Alias: &alias,
	})
	if err != nil {
		panic(err)
	}
}
```

