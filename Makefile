api:
	@air -c air/.api.toml
acc:
	@air -c air/.account.toml

pb:
ifdef s
	$(MAKE) gen SERVICE=$(s)
else
	$(MAKE) gen SERVICE=account
endif
gen:
	@protoc \
		--proto_path=proto "./proto/$(SERVICE).proto" \
		--go_out=pkg/pb \
		--go_opt=paths=source_relative \
		--go-grpc_out=pkg/pb \
		--go-grpc_opt=paths=source_relative