package auth

import (
	"reflect"
	"strings"
	"time"

	"github.com/c0sco/go-ldap-client"
	"github.com/caarlos0/env"
	"github.com/labstack/gommon/log"
	"github.com/mitchellh/mapstructure"
)

// ServerAttributes contains the fields needed to connect to the auth server
type ServerAttributes struct {
	Conn         Server
	BaseDN       string   `env:"AUTH_LDAP_BASE,required"`
	BindDN       string   `env:"AUTH_LDAP_BIND_DN"`
	BindPW       string   `env:"AUTH_LDAP_BIND_PW"`
	Host         string   `env:"AUTH_LDAP_HOST" envDefault:"localhost"`
	Port         int      `env:"AUTH_LDAP_PORT" envDefault:"389"`
	UseTLS       bool     `env:"AUTH_LDAP_TLS" envDefault:"false"`
	StartTLS     bool     `env:"AUTH_LDAP_START_TLS" envDefault:"false"`
	UIDFieldName string   `env:"AUTH_LDAP_UID_NAME" envDefault:"uid"`
	GIDFieldName string   `env:"AUTH_LDAP_GID_NAME" envDefault:"memberUid"`
	Timeout      int      `env:"AUTH_LDAP_TIMEOUT_SECS" envDefault:"3"`
	Fields       []string // Populated based on the mapstructure tags in User
}

// InitializeServer -
func InitializeServer(c *Configuration) {
	// Parse the config from the 'env' tags in ServerAttributes
	env.Parse(&c.Serv)

	// Populate the fields from the 'mapstructure' tags
	c.Serv.Fields = populateFields("mapstructure", User{})

	// Set up the server connection
	c.Serv.Conn = &ldap.LDAPClient{
		Base:         Cfg.Serv.BaseDN,
		Host:         Cfg.Serv.Host,
		Port:         Cfg.Serv.Port,
		UseSSL:       Cfg.Serv.UseTLS,
		SkipTLS:      !Cfg.Serv.StartTLS,
		BindDN:       Cfg.Serv.BindDN,
		BindPassword: Cfg.Serv.BindPW,
		UserFilter:   "(" + Cfg.Serv.UIDFieldName + "=%s)",
		GroupFilter:  "(" + Cfg.Serv.GIDFieldName + "=%s)",
		Attributes:   Cfg.Serv.Fields,
		Timeout:      time.Duration(Cfg.Serv.Timeout) * time.Second,
	}
	if err := c.Serv.Conn.Connect(); err != nil {
		log.Error(err)
	}
}

// populateFields returns a list of tag values on 'u', matching key 'f'
func populateFields(f string, u interface{}) (fieldList []string) {
	t := reflect.TypeOf(u)
	for i := 0; i < t.NumField(); i++ {
		tag := t.Field(i).Tag.Get(f)
		if tag == "" {
			continue
		}
		fieldList = append(fieldList, tag)
	}
	return fieldList
}

// ValidateLogin validates the username and password against the LDAP server.
// Decodes User.Fields from LDAP into the User struct based on the
// 'mapstructure' tag
func ValidateLogin(username, password string) (ok bool, u User, err error) {
	ok, user, err := Cfg.Serv.Conn.Authenticate(username, password)
	// Authenticate returns an error on any invalid login
	if !ok && err != nil && strings.Contains(err.Error(), "Invalid Credentials") {
		log.Infof("Authentication failed for user %s", username)
		return ok, User{}, nil
	} else if !ok || err != nil {
		log.Warnf("Authentication error for user %s:%s", username, err)
		return ok, User{}, err
	}

	log.Infof("Authentication succeeded for user %s", username)

	// Turn the LDAP attributes into a User struct based on the
	// mapstructure tags
	mapstructure.Decode(user, &u)
	return ok, u, err
}
