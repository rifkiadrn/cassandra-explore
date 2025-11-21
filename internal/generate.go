package internal

//go:generate oapi-codegen -config ./model/type.cfg.yaml ../openapi/openapi.yaml
//go:generate oapi-codegen -config ./handler/rest/server.cfg.yaml ../openapi/openapi.yaml
