Goal: Demonstrate how to develop an AWS hosted dynamodb Web application locally

Runs under an AWS_PROFILE called "mine", you will have to change that to yours when deploying to your AWS account.

Start dynamodb server

    ./scripts/local-dynamodb.sh
    ./scripts/create-table.sh

Start Go Web server

    ./scripts/start-local-server.sh
