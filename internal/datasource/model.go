package datasource

import "github.com/google/uuid"

type Status string

const(
	StatusWait	Status = "Wait"
	StatusError	Status = "Error"
	StatusReady	Status = "Ready"
)

type Storage struct{
	Id uuid.UUID
	Files []string
	Archive string
	Status Status
}