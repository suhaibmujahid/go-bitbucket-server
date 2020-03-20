package bitbucket

type Project struct {
	Key         string     `json:"key,omitempty"`
	Id          int        `json:"id,omitempty"`
	Name        string     `json:"name,omitempty"`
	Description string     `json:"description,omitempty"`
	Public      bool       `json:"public,omitempty"`
	Type        string     `json:"type,omitempty"`
	Owner       *User      `json:"owner,omitempty"` // this populated only for personal projects
	Links       *SelfLinks `json:"links,omitempty"`
}
