auth-plug
===

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

**AUTH_LDAP_UID_NAME** _default: uid_  
The LDAP attribute that contains the username.

**AUTH_PORT** _required_  
The port to bind to.

Usage
---

These are the endpoints that are defined.

**/login** _POST_  
Takes a `username` and `password` as post data, validates it against the LDAP server, and sends back a JWT.

**/verify** _GET_  
Returns `OK` if a valid `Authorization` header w/ JWT (type `Bearer`) is supplied and the JWT is validated.

I don't want to use LDAP
---

No problem! The authentication code is contained inside [auth/main.go](auth/main.go) and [auth/ldap.go](auth/ldap.go). Swap it out with your authentication server type.

Caveats
---

auth-plug only does authentication, **not** authorization. Future versions will support some form of authz control (e.g. [casbin](https://github.com/casbin/casbin)).

A standard use-case for auth-plug is running on localhost with the LDAP server and nginx. Full TLS support is yet to be implemented.
