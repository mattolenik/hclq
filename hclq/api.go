package api

import (
	"github.com/hashicorp/hcl2/hcl"
)

// HclqFile is a queryable HCL file
type HclqFile = *hcl.File

// Results contains HCLQ query results
type Results struct {
}

// Query returns results that match a given query string
func (f HclqFile) Query(query string) (*Results, error) {
	return nil, nil
}
