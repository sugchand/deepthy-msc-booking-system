package env

import (
	"bookingSystem/common"
	"os"
	"strconv"

	log "github.com/sirupsen/logrus"
)

const (
	EnvDBHost             = "DB_HOST"
	EnvDBName             = "DB_NAME"
	EnvDBUname            = "DB_UNAME"
	EnvDBPwd              = "DB_PWD"
	EnvGRPCListenPortName = "ROOMINVENTORY_GRPC_LISTEN_PORT"
)

type RoomEnvValues struct {
	dbHost         string
	dbName         string
	dbUname        string
	dbPwd          string
	grpcListenPort uint32
}

// Get the DB path.
func (re *RoomEnvValues) DBRemote() string {

	return re.dbHost + ":" + strconv.Itoa(common.DBPORT)
}

func (re *RoomEnvValues) DBName() string {
	return re.dbName
}

func (re *RoomEnvValues) DBUNameAndPwd() (string, string) {
	return re.dbUname, re.dbPwd
}

func (re *RoomEnvValues) GRPCListenPort() uint32 {
	return re.grpcListenPort
}

func (re *RoomEnvValues) readEnvValues() error {
	re.dbName = os.Getenv(EnvDBName)
	re.dbHost = os.Getenv(EnvDBHost)
	re.dbUname = os.Getenv(EnvDBUname)
	re.dbPwd = os.Getenv(EnvDBPwd)

	port := os.Getenv(EnvGRPCListenPortName)
	intPort, err := strconv.ParseUint(port, 10, 0)
	if err != nil {
		log.WithFields(log.Fields{
			"port": intPort,
		}).WithError(err).Error("failed to read environment variable configuration!")
		return err
	}
	re.grpcListenPort = uint32(intPort)
	return nil
}

func NewRoomInventoryEnv() *RoomEnvValues {
	re := &RoomEnvValues{}
	err := re.readEnvValues()
	if err != nil {
		return nil
	}
	return re
}
