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
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type DnsRecordType string

//goland:noinspection GoUnusedConst
const (
	DnsRecordTypeA     = DnsRecordType("a")
	DnsRecordTypeAAAA  = DnsRecordType("aaaa")
	DnsRecordTypeCAA   = DnsRecordType("caa")
	DnsRecordTypeCNAME = DnsRecordType("cname")
	DnsRecordTypeMX    = DnsRecordType("mx")
	DnsRecordTypeNS    = DnsRecordType("ns")
	DnsRecordTypeSRV   = DnsRecordType("srv")
	DnsRecordTypeSSHFP = DnsRecordType("sshfp")
	DnsRecordTypeTLSA  = DnsRecordType("tlsa")
	DnsRecordTypeTXT   = DnsRecordType("txt")
)

type DnsDomainList []string

//Dns provides a way to interact with DNS domains
type Dns interface {
	//List returns a list of domains which have DNS records managed by Active 24.
	List() (DnsDomainList, ApiError)
	//With returns interface to interact with DNS records in given domain
	With(domain string) DnsRecordActions
}

type DnsRecord struct {
	HashId           *string `json:"hashId,omitempty"`
	Type             *string `json:"type,omitempty"`
	Ttl              int     `json:"ttl"`
	Name             string  `json:"name"`
	Port             *int    `json:"port,omitempty"`
	Priority         *int    `json:"priority,omitempty"`
	Target           *string `json:"target,omitempty"`
	Weight           *int    `json:"weight,omitempty"`
	NameServer       *string `json:"nameserver,omitempty"`
	Ip               *string `json:"ip,omitempty"`
	CaaValue         *string `json:"caaValue,omitempty"`
	Flags            *int    `json:"flags,omitempty"`
	Tag              *string `json:"tag,omitempty"`
	Alias            *string `json:"alias,omitempty"`
	MailServer       *string `json:"mailserver,omitempty"`
	Algorithm        *int    `json:"algorithm,omitempty"`
	FingerprintType  *int    `json:"fingerprintType,omitempty"`
	CertificateUsage *int    `json:"certificateUsage,omitempty"`
	Hash             *string `json:"hash,omitempty"`
	MatchingType     *int    `json:"matchingType,omitempty"`
	Selector         *int    `json:"selector,omitempty"`
	Text             *string `json:"text,omitempty"`
}

// DnsRecordActions allows interaction with DNS records
type DnsRecordActions interface {
	//Create creates a new DNS record
	Create(DnsRecordType, *DnsRecord) ApiError
	//List lists all DNS records in this domain.
	List() ([]DnsRecord, ApiError)
	//Update updates an existing DNS record
	Update(DnsRecordType, *DnsRecord) ApiError
	//Delete removes single DNS record based on its hash ID
	Delete(string) ApiError
}

type dns struct {
	h helper
}

func (d *dns) List() (DnsDomainList, ApiError) {
	resp, err := d.h.do(http.MethodGet, "dns/domains/v1", nil)
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
	doms := make([]string, 0)
	err = json.Unmarshal(body, &doms)
	if err != nil {
		return nil, apiErr(resp, err)
	}
	return doms, apiErr(resp, nil)
}

func (d *dns) With(domain string) DnsRecordActions {
	return &domainAction{
		h:   d.h,
		dom: domain,
	}
}

type domainAction struct {
	h   helper
	dom string
}

func (d *domainAction) Create(recType DnsRecordType, r *DnsRecord) ApiError {
	return apiErr(d.change(http.MethodPost, string(recType), r))
}

func (d *domainAction) change(method string, sub string, r *DnsRecord) (*http.Response, error) {
	data, err := json.Marshal(r)
	if err != nil {
		return nil, err
	}
	return d.h.do(method, fmt.Sprintf("dns/%s/%s/v1", d.dom, sub), bytes.NewBuffer(data))
}

func (d *domainAction) List() ([]DnsRecord, ApiError) {
	resp, err := d.h.do(http.MethodGet, fmt.Sprintf("dns/%s/records/v1", d.dom), nil)
	if err != nil {
		return nil, apiErr(nil, err)
	}
	defer func(b io.ReadCloser) {
		_ = b.Close()
	}(resp.Body)
	if resp.StatusCode > 399 && resp.StatusCode < 600 {
		return nil, apiErr(resp, fmt.Errorf("invalid response from api: %d", resp.StatusCode))
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, apiErr(resp, err)
	}
	ret := make([]DnsRecord, 0)
	err = json.Unmarshal(body, &ret)
	if err != nil {
		return nil, apiErr(resp, err)
	}
	return ret, apiErr(resp, nil)
}

func (d *domainAction) Update(recType DnsRecordType, r *DnsRecord) ApiError {
	return apiErr(d.change(http.MethodPut, string(recType), r))
}

func (d *domainAction) Delete(hashId string) ApiError {
	return apiErr(d.h.do(http.MethodDelete, fmt.Sprintf("/dns/%s/%s/v1", d.dom, hashId), nil))
}
