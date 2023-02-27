package driver

import (
	"github.com/edgexfoundry/device-sdk-go/v2/pkg/service"
	"time"
)

func GetCredentials() (secretData map[string]string, err error) {
	//get creadentials mikrotik
	ds := service.RunningService()
	secretData, err = ds.SecretProvider.GetSecret("Mikrotik-Router", "ip", "username", "password")
	return secretData, err
}

func (s *SimpleDriver) UpdateCredentials() {

	for {
		secretData, err := GetCredentials()
		if err == nil {
			s.mikrotik.Address = secretData["ip"]
			s.mikrotik.Username = secretData["username"]
			s.mikrotik.Password = secretData["password"]
		}
		time.Sleep(1 * time.Minute)

	}
}
