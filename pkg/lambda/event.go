package lambda

// Repository is a simple type for Github
type Repository struct {
	Name     string `json:"name"`
	Fullname string `json:"full_name"`
}

// GithubRepoHead
type GithubRepoHead struct {
	Ref        string     `json:"ref"`
	Repository Repository `json:"repo"`
}

// GithubRepoBase
type GithubRepoBase struct {
	Repository Repository `json:"repo"`
}

// GithubPullRequest is the struct for Github
type GithubPullRequest struct {
	State  string         `json:"state"`
	Head   GithubRepoHead `json:"head"`
	Base   GithubRepoBase `json:"base"`
	Merged bool           `json:"merged"`
}

// GithubWebhookEvent is what we're expecting to parse from Github
type GithubWebhookEvent struct {
	Action      string            `json:"action"`
	PullRequest GithubPullRequest `json:"pull_request"`
	Repository  Repository        `json:"repository"`
}

// RepositoryAndBranch is a bitbucket data object
type RepositoryAndBranch struct {
	Repository Repository      `json:"repository"`
	Branch     BitbucketBranch `json:"branch"`
}

// BitbucketBranch is a kv branch
type BitbucketBranch struct {
	Name string `json:"name"`
}

// BitbucketPullRequest is the struct for Bitbucket
type BitbucketPullRequest struct {
	Destination RepositoryAndBranch `json:"destination"`
	Source      RepositoryAndBranch `json:"source"`
	State       string              `json:"state"` // enum of MERGED || DECLINED
}

// BitbucketWebhookEvent is what we are expecting to parse from bitbucket
type BitbucketWebhookEvent struct {
	Repository  Repository             `json:"repository"`
	PullRequest BitbucketPullRequest   `json:"pullrequest"`
	Actor       map[string]interface{} `json:"actor"`
}
