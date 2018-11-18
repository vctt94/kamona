module github.com/getset0/kamona

replace github.com/vctt94/kamona/database => ./database

require (
	github.com/dgrijalva/jwt-go v3.2.0+incompatible
	github.com/golang/protobuf v1.2.0
	github.com/grpc-ecosystem/grpc-gateway v1.5.1
	github.com/mitchellh/mapstructure v1.1.2
	github.com/vctt94/kamona/database v0.0.0-20181118164052-a1b6aa341049
	golang.org/x/net v0.0.0-20181114220301-adae6a3d119a
	google.golang.org/genproto v0.0.0-20180817151627-c66870c02cf8
	google.golang.org/grpc v1.16.0
)
