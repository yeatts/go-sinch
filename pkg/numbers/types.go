package numbers

type RequestSMSConfiguration struct {
	ServicePlanID string `json:"servicePlanId"` // required
	CampaignID    string `json:"campaignId"`
}

type RequestVoiceConfiguration struct {
	AppID string `json:"appId"`
}

type ResponseSMSConfiguration struct {
	ServicePlanID         string                         `json:"servicePlanId"` // required
	ScheduledProvisioning *ResponseScheduledProvisioning `json:"scheduledProvisioning"`
	CampaignID            string                         `json:"campaignId"`
}

type ResponseScheduledProvisioning struct {
	ServicePlanID   string   `json:"servicePlanId"`
	Status          string   `json:"status"`
	LastUpdatedTime string   `json:"lastUpdatedTime"`
	CampaignID      string   `json:"campaignId"`
	ErrorCodes      []string `json:"errorCodes"`
}

type ResponseVoiceConfiguration struct {
	AppID                      string                              `json:"appId"`
	ScheduledVoiceProvisioning *ResponseScheduledVoiceProvisioning `json:"scheduledVoiceProvisioning"`
	LastUpdatedTime            string                              `json:"lastUpdatedTime"`
}

type ResponseScheduledVoiceProvisioning struct {
	AppID           string `json:"appId"`
	Status          string `json:"status"`
	LastUpdatedTime string `json:"lastUpdatedTime"`
}
