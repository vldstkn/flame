api:
	@air -c air/.api.toml
acc:
	@air -c air/.account.toml
match:
	@air -c air/.matching.toml
swipes:
	@air -c air/.swipes.toml


pb:
ifdef s
	$(MAKE) gen SERVICE=$(s)
else
	$(MAKE) gen SERVICE=account
	$(MAKE) gen SERVICE=matching
	$(MAKE) gen SERVICE=swipes
endif
gen:
	@protoc \
		--proto_path=proto "./proto/$(SERVICE).proto" \
		--go_out=pkg/pb \
		--go_opt=paths=source_relative \
		--go-grpc_out=pkg/pb \
		--go-grpc_opt=paths=source_relative