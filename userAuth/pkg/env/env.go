package env

import (
	"bookingSystem/common"
	"os"
	"strconv"
	"time"

	log "github.com/sirupsen/logrus"
)

const (
	EnvDBHost                  = "DB_HOST"
	EnvDBName                  = "DB_NAME"
	EnvDBUname                 = "DB_UNAME"
	EnvDBPwd                   = "DB_PWD"
	EnvGRPCListenPortName      = "USERAUTH_GRPC_LISTEN_PORT"
	EnvUserAccessTokenValidity = "USERAUTH_TOKEN_VALIDITY_SECONDS"
)

type UserEnvValues struct {
	dbHost            string
	dbName            string
	dbUname           string
	dbPwd             string
	grpcListenPort    uint32
	userTokenValidity time.Duration
}

// Get the DB path.
func (ue *UserEnvValues) DBRemote() string {

	return ue.dbHost + ":" + strconv.Itoa(common.DBPORT)
}

func (ue *UserEnvValues) DBName() string {
	return ue.dbName
}

func (ue *UserEnvValues) DBUNameAndPwd() (string, string) {
	return ue.dbUname, ue.dbPwd
}

func (ue *UserEnvValues) GRPCListenPort() uint32 {
	return ue.grpcListenPort
}

func (ue *UserEnvValues) TokenValidity() time.Duration {
	return ue.userTokenValidity
}

func (ue *UserEnvValues) readEnvValues() error {
	ue.dbName = os.Getenv(EnvDBName)
	ue.dbHost = os.Getenv(EnvDBHost)
	ue.dbUname = os.Getenv(EnvDBUname)
	ue.dbPwd = os.Getenv(EnvDBPwd)

	validity := os.Getenv(EnvUserAccessTokenValidity)
	validitySeconds, err := strconv.ParseUint(validity, 10, 0)
	if err != nil {
		log.WithFields(log.Fields{
			"validity": validitySeconds,
		}).WithError(err).Error("failed to read environment variable configuration!")
		return err
	}
	ue.userTokenValidity = time.Duration(validitySeconds) * time.Second

	port := os.Getenv(EnvGRPCListenPortName)
	intPort, err := strconv.ParseUint(port, 10, 0)
	if err != nil {
		log.WithFields(log.Fields{
			"port": intPort,
		}).WithError(err).Error("failed to read environment variable configuration!")
		return err
	}
	ue.grpcListenPort = uint32(intPort)
	return nil
}

func NewUserEnvironment() *UserEnvValues {
	ue := &UserEnvValues{}
	err := ue.readEnvValues()
	if err != nil {
		return nil
	}
	return ue
}
