#!/bin/bash

# -- > Create DynamoDb Table
echo Creating  DynamoDb \'payments\' table ...
echo $(aws --endpoint-url=http://localhost:4566 dynamodb create-table --cli-input-json '{"TableName":"payments", "KeySchema":[{"AttributeName":"order_id","KeyType":"HASH"}], "AttributeDefinitions":[ {"AttributeName":"order_id","AttributeType":"S"}], "BillingMode":"PAY_PER_REQUEST"}' --profile test-profile --region us-east-1 --output table | cat)

# --> List DynamoDb Tables
echo Listing tables ...
echo $(aws --endpoint-url=http://localhost:4566 dynamodb list-tables --profile test-profile --region us-east-1 --output table | cat)

# -- > Create SNS Topic
echo Creating  SNS \'update_order_status\' topic ...
echo $(aws --endpoint-url=http://localhost:4566 sns create-topic --name update_order_status --region us-east-1 --profile test-profile --output table | cat)

# -- > List SNS Topics
echo Listing topics ...
echo $(aws --endpoint-url=http://localhost:4566 sns list-topics --profile test-profile --region us-east-1 --output table | cat)
