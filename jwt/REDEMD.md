### Generate RSA256 key
```bash
$ openssl genrsa -out rs256-private.pem 2048
$ openssl genpkey -algorithm RSA -out rs256_private.pem -pkeyopt rsa_keygen_bits:2048
```
### Generate RSA256 public key
```bash
// 
$ openssl rsa -in rs256_private.pem -pubout -out rs256_public.pem
```
### Generate ES key
```bash
// p-256
$ openssl ecparam -name prime256v1 -genkey -noout -out es256_private.pem
// p-384
$ openssl ecparam -name secp384r1 -genkey -noout -out es384_private.pem
// p-521
$ openssl ecparam -name secp521r1 -genkey -noout -out es512_private.pem
```
### Generate ES public key
```bash
// p-256 public
$ openssl ec -in es256_private.pem -pubout -out es256_public.pem
```

### wire injection
```go
wire.NewSet(NewEES256JWTFromOptions, wire.Bind(new(IJWT), new(*EES256JWT)))
wire.NewSet(NewEHS256JWTFromOptions, wire.Bind(new(IJWT), new(*EHS256JWT)))
wire.NewSet(NewEHS384JWTFromOptions, wire.Bind(new(IJWT), new(*EHS384JWT)))
wire.NewSet(NewEHS512JWTFromOptions, wire.Bind(new(IJWT), new(*EHS512JWT)))
wire.NewSet(NewERS256JWTFromOptions, wire.Bind(new(IJWT), new(*ERS256JWT)))
wire.NewSet(NewES256JWTFromOptions, wire.Bind(new(IJWT), new(*ES256JWT)))
wire.NewSet(NewHS256JWTFromOptions, wire.Bind(new(IJWT), new(*HS256JWT)))
wire.NewSet(NewHS384JWTFromOptions, wire.Bind(new(IJWT), new(*HS384JWT)))
wire.NewSet(NewRS256JWTFromOptions, wire.Bind(new(IJWT), new(*RS256JWT)))



```