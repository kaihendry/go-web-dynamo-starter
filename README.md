Goal: Develop a AWS dynamodb Web application locally

Runs under an AWS_PROFILE called "mine", you will have to change that to yours when deploying to your AWS account.

Start dynamodb server

    ./scripts/local-dynamodb.sh
    ./scripts/create-table.sh

Start Go Web server

    ./scripts/start-local-server.sh

If you like this, check out https://github.com/kaihendry/local-audio which builds on this.
