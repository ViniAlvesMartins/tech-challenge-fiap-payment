#!/bin/bash

files='(/doc|/infra|/doc/swagger|/application|/contract|/mock|/src|/external|/handler|/http_server|/api|/external/database/dynamodb|/src/application/modules/response/order_service|/src/application/modules/response/payment_service|/src/controller/serializer/input|/src/controller/serializer/output|/src/controller/serializer|src/pkg/uuid)'
go list ./... | egrep -v $files$\