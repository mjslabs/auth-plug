auth-plug
===

[![Build Status][travis-badge]][travis]
[![Go Report Card][goreport-badge]][goreport]
[![Test Coverage][coverage]][codeclimate]

[travis-badge]: https://travis-ci.org/mjslabs/auth-plug.svg?branch=master
[travis]: https://travis-ci.org/mjslabs/auth-plug
[goreport-badge]: https://goreportcard.com/badge/github.com/mjslabs/auth-plug
[goreport]: https://goreportcard.com/report/github.com/mjslabs/auth-plug
[coverage]: https://api.codeclimate.com/v1/badges/4976c6d311f5c4ac37c4/test_coverage
[codeclimate]: https://codeclimate.com/github/mjslabs/auth-plug/test_coverage

Small Go service that takes LDAP logins and hands out JWTs. Very useful for adding
authentication to an otherwise unsecured API, and easily adaptable to other authentication methods.

Build and Test
---

Standard Go build and test methods apply. Uses modules (Go 1.11+).

```bash
go test -cover ./...
go build -o auth-plug
```

Configure and Run
---

All configuration is done at run time with the following environment variables.

**AUTH_IP** _default: all interfaces_  
The IP address to bind to.

**AUTH_JWT_METHOD** _default: HS512_  
The signing method to use for creating JWTs.

**AUTH_JWT_SECRET** _required_  
The key to use to sign JWTs.

**AUTH_JWT_VALID_MIN** _default: 30_  
The number of minutes a JWT is valid for.

**AUTH_LDAP_BASE** _required_  
The base DN to use when searching the LDAP server.

**AUTH_LDAP_BIND_DN** _default: \<empty\>_  
The DN to bind to the LDAP server with (i.e. the username).

**AUTH_LDAP_BIND_PW** _default: \<empty\>_  
The password to use when binding to the LDAP server.

**AUTH_LDAP_HOST** _default: localhost_  
The hostname or IP of the LDAP server.

**AUTH_LDAP_PORT** _default: 389_  
The port of the LDAP server.

**AUTH_LDAP_TLS** _default: false_  
Use TLS to connect to the LDAP server.

**AUTH_LDAP_START_TLS** _default: false_  
Use STARTTLS to connect to the LDAP server.

**AUTH_LDAP_GID_NAME** _default: memberUid_  
The LDAP attribute that maps a user to a group.  
This functionality is not yet implemented.

**AUTH_LDAP_TIMEOUT_SECS** _default: 3_  
The number of seconds to wait for the LDAP server to respond.

**AUTH_LDAP_UID_NAME** _default: uid_  
The LDAP attribute that contains the username.

**AUTH_PORT** _required_  
The port to bind to.

**AUTH_PROFILE**  
The `[ip]:<port>` for a pprof web server to listen on.  
This then enables the standard set of `/debug` pprof endpoints.

Usage
---

auth-plug follows a familiar flow.

1. POST a `username` and `password` to `/login`.
2. Retrieve the JWT from the response.
3. Send a GET to `/verify`, setting the JWT from step 2 in the `Authorization` header.
4. If step 3 fails, go back to step 1.

Here is a full list of defined endpoints.

**/login** _POST_  
Takes a `username` and `password` as post data, validates it against the LDAP server, and sends back a JWT.

**/verify** _GET_  
Returns `OK` if a valid `Authorization` header w/ JWT (type `Bearer`) is supplied and the JWT is validated.

**/health** _GET_  
Returns an HTTP 200 on healthy and HTTP 503 if an error is found with the service.  
Always returns a JSON structure with a `status` key.

Healthy

```JSON
{"status":"OK"}
```

Unhealthy (e.g.)

```JSON
{"status":"LDAP Result Code 200 \"Network Error\": dial tcp: lookup bad.examplehost.com: no such host"}
```

I don't want to use LDAP
---

No problem! The authentication code is contained inside [auth/main.go](auth/main.go) and [auth/ldap.go](auth/ldap.go). Swap it out with your authentication server type.

Caveats
---

auth-plug only does authentication, **not** authorization. Future versions will support some form of authz control (e.g. [casbin](https://github.com/casbin/casbin)).

A standard use-case for auth-plug is running on localhost with the LDAP server and nginx. Full TLS support is yet to be implemented.

[go-ldap-client](https://github.com/jtblin/go-ldap-client) is used for the LDAP work. This library seems to be abandoned and should be changed out for something that is actively maintained.
