package repairorder

type Status string

const (
	StatusReceived    Status = "received"
	StatusInDiagnosis Status = "in_diagnosis"
	StatusAnalyzed    Status = "analyzed"
	StatusApproved    Status = "approved"
	StatusRejected    Status = "rejected"
	StatusClosed      Status = "closed"
)
