PROFILE=mine
STACK=github-actions-test
GITHUBORG=kaihendry
REPOSITORY=go-web-dynamo-starter

deploy:
	AWS_PROFILE=$(PROFILE) aws cloudformation deploy --stack-name $(STACK) --template-file configure-aws-credentials.yml --parameter-overrides RepositoryName=$(REPOSITORY) GitHubOrg=$(GITHUBORG) --capabilities CAPABILITY_IAM

getrole:
	aws iam get-role --role-name github-actions-test-Role-QVSUYDA0W3QG

# view attached policies
list-attached-role-policies:
	aws iam list-attached-role-policies --role-name github-actions-test-Role-QVSUYDA0W3QG
