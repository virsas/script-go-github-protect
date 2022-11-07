# script-go-github-protect

Golang script to protect all repositories in single Organization

## .env configuration

``` bash
GITHUB_ORG="example"
GITHUB_TOKEN="ghp_xxxxxxxxxxxxxxxxxxxxxxxxxx"
```

## run

go run main.go

## configuration

``` config
DismissStaleReviews:            True
RequireCodeOwnerReviews:        True
RequiredApprovingReviewCount:   1
RequiredConversationResolution: True
EnforceAdmins:                  True
```