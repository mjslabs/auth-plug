package auth

import (
	"time"

	"github.com/caarlos0/env"
)

// Server is the interface that the auth server client library satisfies
type Server interface {
	Close()
	Connect() error
	Authenticate(string, string) (bool, map[string]string, error)
	GetGroupsOfUser(string) ([]string, error)
}

// Configuration defines the main authentication configuration parameters.
// This should contain a connection to your authentication server, plus
// configuration for JWT and the auth server's connection parameters.
type Configuration struct {
	// JWT configuration
	JWTMethod       string        `env:"AUTH_JWT_METHOD" envDefault:"HS512"`
	JWTSecret       string        `env:"AUTH_JWT_SECRET,required"`
	JWTValidMinutes time.Duration `env:"AUTH_JWT_VALID_MIN" envDefault:"30m"`

	// Server attributes defined in
	Serv ServerAttributes
}

// Cfg is the instance of Configuration.
var Cfg Configuration

// User contains data about the user that we care about.
//
// The 'mapstructure' tag maps to LDAP fields, and 'structs' maps to JWT
// claims e.g. Email is stored in ldap as "mail", and in the JWT claims
// as "email". There must be a mapstructure tag for the field to be pulled
// from LDAP and similarly there must be a structs tag for the field to be
// put into the JWT claims
type User struct {
	Email    string `mapstructure:"mail"  structs:"email"`
	Realname string `mapstructure:"gecos" structs:"fullname"`
	Username string `mapstructure:"uid"   structs:"name"`
	// Groups isn't directly pulled from a single LDAP field, nor is it
	// stored in the JWT claims
	Groups []string `structs:",omitempty"`
}

// Initialize the Cfg instance and server config
func Initialize() {
	Cfg = Configuration{}
	env.Parse(&Cfg)

	InitializeServer(&Cfg)
}
