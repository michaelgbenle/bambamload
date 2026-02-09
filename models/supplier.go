package models

type SubmitBusinessProfileRequest struct {
	AccountType         string `json:"account_type"`
	BusinessDescription string `json:"business_description"`
	YearFounded         string `json:"year_founded"`
	WebsiteUrl          string `json:"website_url"`
	LinkedInProfile     string `json:"linked_in_profile"`
	Country             string `json:"country"`
	State               string `json:"state"`
	Address             string `json:"address"`
	RegionsServed       string `json:"regions_served"`
}
