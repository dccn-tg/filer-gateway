package task

// SetProjectResource defines json (un)marshal-able task data for setting project resource.
type SetProjectResource struct {
	ProjectID string `json:"projectID"`
	Storage   Storage
	Members   []Member
}

// SetUserResource defines json (un)marshal-able task data for setting user resource.
type SetUserResource struct {
	UserID  string `json:"userID"`
	Storage Storage
}

// Member defines json (un)marshal-able member information for project resource.
type Member struct {
	UserID string `json:"userID"`
	Role   string `json:"role"`
}

// Storage defines json (un)marshal-able storage information.
type Storage struct {
	System  string `json:"system"`
	QuotaGb int64  `json:"quotaGb"`
}
