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

package active24

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type DomainInfo struct {
	Name           string `json:"name"`
	Status         string `json:"status"`
	ExpirationDate uint64 `json:"expirationDate"`
	HolderFullName string `json:"holderFullName"`
	Payable        bool   `json:"payable"`
	PayedTo        uint64 `json:"payedTo"`
}

type DomainConfigurationField struct {
	Error           string      `json:"error"`
	Name            string      `json:"name"`
	Order           int         `json:"order"`
	Required        bool        `json:"required"`
	Type            string      `json:"type"`
	Value           interface{} `json:"value"`
	ValidationRegex string      `json:"validationRegex"`
	WriteOnly       bool        `json:"writeOnly"`
}

type DomainConfigurationPart struct {
	Fields       []DomainConfigurationField `json:"fields"`
	Updatable    bool                       `json:"updatable"`
	ProcessState *string                    `json:"processState"`
}

type DomainConfiguration struct {
	Holder       DomainConfigurationPart `json:"domainHolderConfiguration"`
	AdminContact DomainConfigurationPart `json:"domainAdminContactConfiguration"`
	Nameserver   DomainConfigurationPart `json:"domainNameserverConfiguration"`
}

type DomainList []DomainInfo

type Domains interface {
	//List returns a list of domains where user is owner or payer.
	List() (DomainList, ApiError)
	// Expiration gets domain expiration date
	Expiration(domain string) (*time.Time, ApiError)
	//Status get domain status - checks if domain can be registered.
	Status(domain string) (*string, ApiError)
}

type domains struct {
	h helper
}

func (d *domains) List() (DomainList, ApiError) {
	resp, err := d.h.do(http.MethodGet, "domains/v1", nil)
	if err != nil {
		return nil, apiErr(nil, err)
	}
	defer func(b io.ReadCloser) {
		_ = b.Close()
	}(resp.Body)
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, apiErr(resp, err)
	}
	doms := make([]DomainInfo, 0)
	err = json.Unmarshal(body, &doms)
	if err != nil {
		return nil, apiErr(resp, err)
	}
	return doms, apiErr(resp, nil)
}

func (d *domains) Expiration(domain string) (*time.Time, ApiError) {
	resp, err := d.h.do(http.MethodGet, fmt.Sprintf("domains/%s/expiration/v1", domain), nil)
	if err != nil {
		return nil, apiErr(nil, err)
	}
	defer func(b io.ReadCloser) {
		_ = b.Close()
	}(resp.Body)
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, apiErr(resp, err)
	}
	ret, err := time.Parse(time.RFC3339, string(body))
	if err != nil {
		return nil, apiErr(resp, err)
	}
	return &ret, apiErr(resp, nil)
}

func (d *domains) Status(domain string) (*string, ApiError) {
	resp, err := d.h.do(http.MethodGet, fmt.Sprintf("domains/%s/expiration/v1", domain), nil)
	if err != nil {
		return nil, apiErr(nil, err)
	}
	defer func(b io.ReadCloser) {
		_ = b.Close()
	}(resp.Body)
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, apiErr(resp, err)
	}
	ret := string(body)
	return &ret, apiErr(resp, nil)
}
