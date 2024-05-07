

.PHONY: gen
gen:
	@#oapi-codegen -generate gin -o ports/openapi.gen.go -package ports api/openapi/*.yaml
	@oapi-codegen -generate spec,gin -o api/http/v1/openapi.gen.go -package v1 api/openapi/*.yaml
	@oapi-codegen -generate types -o api/http/v1//openapi.types.go -package v1 api/openapi/*.yaml