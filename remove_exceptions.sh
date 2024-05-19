#!/bin/bash

files='(/doc|/infra|/src|/application|/contract|/mock|/src|/external|/handler|/http_server|/api)'
go list ./... | egrep -v $files$\