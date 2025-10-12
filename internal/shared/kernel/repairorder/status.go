package repairorder

type Status string

const (
	StatusReceived         Status = "received"
	StatusInDiagnostics    Status = "in_diagnostics"
	StatusAwaitingApproval Status = "awaiting_approval"
	StatusApproved         Status = "approved"
	StatusInExecution      Status = "in_execution"
	StatusFinished         Status = "finished"
	StatusReleased         Status = "released"
	StatusCanceled         Status = "canceled"
)
