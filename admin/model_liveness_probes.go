/*
 * Paygate Admin API
 *
 * PayGate is a RESTful API enabling first-party Automated Clearing House ([ACH](https://en.wikipedia.org/wiki/Automated_Clearing_House)) transfers to be created without a deep understanding of a full NACHA file specification. First-party transfers initiate at an Originating Depository Financial Institution (ODFI) and are sent off to other Financial Institutions.
 *
 * API version: v1
 * Generated by: OpenAPI Generator (https://openapi-generator.tech)
 */

package admin

// LivenessProbes struct for LivenessProbes
type LivenessProbes struct {
	// Either an error from checking Customers or good as a string.
	Customers string `json:"customers,omitempty"`
}
