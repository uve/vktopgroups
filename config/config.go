package config

import (
	"fmt"
	"os"
	"encoding/json"
	
	
	"appengine"

	//"github.com/martini-contrib/oauth2"
	//"github.com/golang/oauth2"
)


const config_prod = "config/config_prod.json"
const config_dev  = "config/config_dev.json"



// and access level. A sample configuration:
type Options struct {
	// ClientID is the OAuth client identifier used when communicating with
	// the configured OAuth provider.
	ClientID string// `json:"client_id"`

	// ClientSecret is the OAuth client secret used when communicating with
	// the configured OAuth provider.
	ClientSecret string// `json:"client_secret"`

	// RedirectURL is the URL to which the user will be returned after
	// granting (or denying) access.
	//RedirectURL string `json:"redirect_url"`

	// Scopes optionally specifies a list of requested permission scopes.
	//Scopes []string `json:"scopes,omitempty"`
}


type OAuthProviders struct{

	Google    Options
	/*
	Facebook  Options
	Github    Options
	Twitter   Options
	*/
	
}


var Config = struct {

	
	CookieName string
	CookieSecret string
	
	OAuthProviders OAuthProviders
	// Your OAuth configuration information for protected user data access.
	//OAuthConfig oauth.Config
	
	// The path in your application to which users will be redirected after they
	// allow or deny permission for your application to access their data.
	RedirectURL string
	// The scheme, hostname and port at which your application can be accessed
	// when running on App Engine.
	RootUrl string
}{}





	

func init() {	
	
	var configPath string
		
	if appengine.IsDevAppServer() {
	
		configPath = config_dev
		
	} else {
	
		configPath = config_prod	
	}
	

	configFile, err := os.Open(configPath)
	if err != nil {
		panic(fmt.Sprintf("Could not open %s: %s", configPath, err))
	}
	defer configFile.Close()


	if err = json.NewDecoder(configFile).Decode(&Config); err != nil {
		panic(fmt.Sprintf("Could not parse %s: %s", configPath, err))
	}


}
