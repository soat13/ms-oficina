package repairorder

type Status string

const (
	StatusReceived         Status = "received"
	StatusAnalyzed         Status = "analyzed"
	StatusAwaitingApproval Status = "awaiting_approval"
	StatusApproved         Status = "approved"
	StatusInDiagnosis      Status = "in_diagnosis"
	StatusRejected         Status = "rejected"
	StatusClosed           Status = "closed"
)
