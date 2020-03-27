package bitbucket

import "time"

// todo: the following are the missed event types:
//  * mirror:repo_synchronized
//  * pr:comment:added
//  * pr:comment:edited
//  * pr:comment:deleted
//  * repo:comment:added
//  * repo:comment:edited
//  * repo:comment:deleted

const (
	EventKeyRepositoryPush              = "repo:refs_changed"
	EventKeyRepositoryModified          = "repo:modified"
	EventKeyRepositoryForked            = "repo:forked"
	EventKeyPullRequestOpened           = "pr:opened"
	EventKeyPullRequestReviewersUpdated = "pr:reviewer:updated"
	EventKeyPullRequestModified         = "pr:modified"
	EventKeyPullRequestBranchUpdated    = "pr:from_ref_updated"
	EventKeyPullRequestApproved         = "pr:reviewer:approved"
	EventKeyPullRequestUnapproved       = "pr:reviewer:unapproved"
	EventKeyPullRequestNeedsWork        = "pr:reviewer:needs_work"
	EventKeyPullRequestMerged           = "pr:merged"
	EventKeyPullRequestDeclined         = "pr:declined"
	EventKeyPullRequestDeleted          = "pr:deleted"
)

// PushEvent is triggered when a user pushes one or more commits, branch created or deleted, or tag created or deleted.
// This payload has a event key of repo:refs_changed
//
// Doc: https://confluence.atlassian.com/bitbucketserver070/event-payload-996644369.html#Eventpayload-Push
type PushEvent struct {
	EventKey   string            `json:"eventKey"`
	Date       time.Time         `json:"date"`
	Actor      *User             `json:"actor"`
	Repository *Repository       `json:"repository"`
	Changes    []PushEventChange `json:"changes"`
}

type PushEventChange struct {
	Ref      Ref    `json:"ref"`
	RefID    string `json:"refId"`
	FromHash string `json:"fromHash"`
	ToHash   string `json:"toHash"`
	Type     string `json:"type"`
}

type Ref struct {
	ID        string `json:"id"`
	DisplayID string `json:"displayId"`
	Type      string `json:"type"`
}

// RepositoryModifiedEvent is triggered when a repository is renamed or moved.
// This payload has a event key of repo:modified
//
// Doc: https://confluence.atlassian.com/bitbucketserver070/event-payload-996644369.html#Eventpayload-Modified
type RepositoryModifiedEvent struct {
	EventKey string      `json:"eventKey"`
	Date     string      `json:"date"`
	Actor    *User       `json:"actor"`
	Old      *Repository `json:"old"`
	New      *Repository `json:"new"`
}

// RepositoryModifiedEvent is triggered when a repository is forked.
// This payload has a event key of repo:forked
//
// Doc: https://confluence.atlassian.com/bitbucketserver070/event-payload-996644369.html#Eventpayload-Fork
type RepositoryForkedEvent struct {
	EventKey   string      `json:"eventKey"`
	Date       string      `json:"date"`
	Actor      *User       `json:"actor"`
	Repository *Repository `json:"repository"`
}

// PullRequestModifiedEvent is triggered when a pull request's description, title, or target branch is changed.
// This payload has a event key of pr:modified
//
// Doc: https://confluence.atlassian.com/bitbucketserver070/event-payload-996644369.html#Eventpayload-Modified.1
type PullRequestModifiedEvent struct {
	EventKey            string             `json:"eventKey"`
	Date                time.Time          `json:"date"`
	Actor               *User              `json:"actor"`
	PullRequest         *PullRequest       `json:"pullRequest"`
	PreviousTitle       string             `json:"previousTitle"`
	PreviousDescription string             `json:"previousDescription"`
	PreviousTarget      *PullRequestTarget `json:"previousTarget"`
}

type PullRequestTarget struct {
	ID              string `json:"id"`
	DisplayID       string `json:"displayId"`
	Type            string `json:"type"`
	LatestCommit    string `json:"latestCommit"`
	LatestChangeset string `json:"latestChangeset"`
}

// PullRequestBranchUpdatedEvent is triggered when the source branch (FromRef) updated.
// This payload has a event key of pr:from_ref_updated
//
// Doc: https://confluence.atlassian.com/bitbucketserver070/event-payload-996644369.html#Eventpayload-Sourcebranchupdated
type PullRequestBranchUpdatedEvent struct {
	EventKey         string       `json:"eventKey"`
	Date             time.Time    `json:"date"`
	Actor            *User        `json:"actor"`
	PullRequest      *PullRequest `json:"pullRequest"`
	PreviousFromHash string       `json:"previousFromHash"`
}

// PullRequestOpenedEvent is triggered when a pull request is opened or reopened.
// This payload has a event key of pr:opened
//
// Doc: https://confluence.atlassian.com/bitbucketserver070/event-payload-996644369.html#Eventpayload-Opened
type PullRequestOpenedEvent PullRequestEvent

// PullRequestMergedEvent is triggered when a user merges a pull request for a repository.
// This payload has a event key of pr:merged
//
// Doc: https://confluence.atlassian.com/bitbucketserver070/event-payload-996644369.html#Eventpayload-Merged
type PullRequestMergedEvent PullRequestEvent

// PullRequestDeclinedEvent is triggered when a user declines a pull request for a repository.
// This payload has a event key of pr:declined
//
// Doc: https://confluence.atlassian.com/bitbucketserver070/event-payload-996644369.html#Eventpayload-Declined
type PullRequestDeclinedEvent PullRequestEvent

// PullRequestDeletedEvent is triggered when a user deletes a pull request for a repository.
// This payload has a event key of pr:deleted
//
// Doc: https://confluence.atlassian.com/bitbucketserver070/event-payload-996644369.html#Eventpayload-Deleted
type PullRequestDeletedEvent PullRequestEvent

// PullRequestEvent present the payload schema some general pull request
// events suh `PullRequestDeclinedEvent` and `PullRequestOpenedEvent`.
type PullRequestEvent struct {
	EventKey    string       `json:"eventKey"`
	Date        time.Time    `json:"date"`
	Actor       *User        `json:"actor"`
	PullRequest *PullRequest `json:"pullRequest"`
}

// PullRequestReviewersUpdatedEvent is triggered when a pull request's reviewers have been added or removed.
// This payload has a event key of pr:reviewer:updated
//
// Doc: https://confluence.atlassian.com/bitbucketserver070/event-payload-996644369.html#Eventpayload-ReviewersUpdated
type PullRequestReviewersUpdatedEvent struct {
	EventKey         string       `json:"eventKey"`
	Date             string       `json:"date"`
	Actor            *User        `json:"actor"`
	PullRequest      *PullRequest `json:"pullRequest"`
	AddedReviewers   []*User      `json:"addedReviewers"`
	RemovedReviewers []*User      `json:"removedReviewers"`
}

// PullRequestApprovedEvent is triggered when a pull request is marked as approved by a reviewer.
// This payload has a event key of pr:reviewer:approved
//
// Doc: https://confluence.atlassian.com/bitbucketserver070/event-payload-996644369.html#Eventpayload-Approved
type PullRequestApprovedEvent PullRequestReviewerEvent

// PullRequestUnapprovedEvent is triggered when a pull request is unapproved by a reviewer.
// This payload has a event key of pr:reviewer:unapproved
//
// Doc: https://confluence.atlassian.com/bitbucketserver070/event-payload-996644369.html#Eventpayload-Unapproved
type PullRequestUnapprovedEvent PullRequestReviewerEvent

// PullRequestNeedsWorkEvent is triggered when a pull request is marked as needs work by a reviewer.
// This payload has a event key of pr:reviewer:needs_work
//
// Doc: https://confluence.atlassian.com/bitbucketserver070/event-payload-996644369.html#Eventpayload-Needswork
type PullRequestNeedsWorkEvent PullRequestReviewerEvent

// PullRequestReviewerEvent present the payload schema for events related to pull
// request review suh `PullRequestNeedsWorkEvent` and `PullRequestApprovedEvent`.
type PullRequestReviewerEvent struct {
	EventKey       string           `json:"eventKey"`
	Date           time.Time        `json:"date"`
	Actor          *User            `json:"actor"`
	PullRequest    *PullRequest     `json:"pullRequest"`
	Participant    *PullRequestUser `json:"participant"`
	PreviousStatus string           `json:"previousStatus"`
}
