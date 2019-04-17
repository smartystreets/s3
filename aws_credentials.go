package s3

import (
	"bufio"
	"bytes"
	"encoding/json"
	"net"
	"net/http"
	"os"
	"time"
)

type awsCredentials struct {
	AccessKeyID     string
	SecretAccessKey string
	SecurityToken   string `json:"Token"`
	Expiration      time.Time
}

// expired checks to see if the temporary credentials from an IAM role are
// within 4 minutes of expiration (The IAM documentation says that new keys
// will be provisioned 5 minutes before the old keys expire). Credentials
// that do not have an Expiration cannot expire.
func (this *awsCredentials) expired() bool {
	if this.Expiration.IsZero() {
		// Credentials with no expiration can't expire
		return false
	}
	expireTime := this.Expiration.Add(-4 * time.Minute)
	// if t - 4 mins is before now, true
	if expireTime.Before(time.Now()) {
		return true
	} else {
		return false
	}
}

const (
	envAccessKey       = "AWS_ACCESS_KEY"
	envAccessKeyID     = "AWS_ACCESS_KEY_ID"
	envSecretKey       = "AWS_SECRET_KEY"
	envSecretAccessKey = "AWS_SECRET_ACCESS_KEY"
	envSecurityToken   = "AWS_SECURITY_TOKEN"
)

// ambientCredentials produces a set of credentials based on the environment
func ambientCredentials() awsCredentials {
	// First use credentials from environment variables
	newCredentials := loadCredentialsFromEnvironment()

	// If there is no Access Key and you are on EC2, get the key from the role
	if (newCredentials.AccessKeyID == "" || newCredentials.SecretAccessKey == "") && onEC2() {
		newCredentials = getIAMRoleCredentials()
	}

	// If the key is expiring, get a new key
	if newCredentials.expired() && onEC2() {
		newCredentials = getIAMRoleCredentials()
	}

	return newCredentials
}

func loadCredentialsFromEnvironment() (credentials awsCredentials) {
	credentials.AccessKeyID = os.Getenv(envAccessKeyID)
	if credentials.AccessKeyID == "" {
		credentials.AccessKeyID = os.Getenv(envAccessKey)
	}
	credentials.SecretAccessKey = os.Getenv(envSecretAccessKey)
	if credentials.SecretAccessKey == "" {
		credentials.SecretAccessKey = os.Getenv(envSecretKey)
	}
	credentials.SecurityToken = os.Getenv(envSecurityToken)
	return credentials
}

// onEC2 checks to see if the program is running on an EC2 instance.
// It does this by looking for the EC2 metadata service.
// This caches that information in a struct so that it doesn't waste time.
func onEC2() bool {
	if location == nil {
		location = &awsLocation{}
	}
	if !(location.checked) {
		c, err := net.DialTimeout("tcp", "169.254.169.254:80", time.Millisecond*100)

		if err != nil {
			location.ec2 = false
		} else {
			_ = c.Close()
			location.ec2 = true
		}
		location.checked = true
	}

	return location.ec2
}

type awsLocation struct {
	ec2     bool
	checked bool
}

var location *awsLocation

// getIAMRoleList gets a list of the roles that are available to this instance
func getIAMRoleList() []string {
	var roles []string
	address := "http://169.254.169.254/latest/meta-data/iam/security-credentials/"

	client := &http.Client{}

	request, err := http.NewRequest("GET", address, nil)

	if err != nil {
		return roles
	}

	response, err := client.Do(request)

	if err != nil {
		return roles
	}
	defer func() { _ = response.Body.Close() }()

	scanner := bufio.NewScanner(response.Body)
	for scanner.Scan() {
		roles = append(roles, scanner.Text())
	}
	return roles
}

func getIAMRoleCredentials() awsCredentials {
	roles := getIAMRoleList()

	if len(roles) == 0 {
		return awsCredentials{}
	}

	// Use the first role in the list
	role := roles[0]

	address := "http://169.254.169.254/latest/meta-data/iam/security-credentials/"

	// Create the full URL of the role
	var buffer bytes.Buffer
	buffer.WriteString(address)
	buffer.WriteString(role)
	roleURL := buffer.String()

	// Get the role
	roleRequest, err := http.NewRequest("GET", roleURL, nil)

	if err != nil {
		return awsCredentials{}
	}

	client := &http.Client{}
	roleResponse, err := client.Do(roleRequest)

	if err != nil {
		return awsCredentials{}
	}
	defer func() { _ = roleResponse.Body.Close() }()

	roleBuffer := new(bytes.Buffer)
	_, _ = roleBuffer.ReadFrom(roleResponse.Body)

	credentials := awsCredentials{}

	err = json.Unmarshal(roleBuffer.Bytes(), &credentials)

	if err != nil {
		return awsCredentials{}
	}

	return credentials
}
