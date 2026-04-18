package repairorder

type Status string

const (
	StatusReceived            Status = "received"
	StatusInDiagnostics       Status = "in_diagnostics"
	StatusDiagnosticsFinished Status = "diagnostics_finished"
	StatusAwaitingApproval    Status = "awaiting_approval"
	StatusApproved            Status = "approved"
	StatusInExecution         Status = "in_execution"
	StatusFinished            Status = "finished"
	StatusPaymentCreated      Status = "payment_created"
	StatusPaymentProcessing   Status = "payment_processing"
	StatusPaymentSucceeded    Status = "payment_succeeded"
	StatusPaymentFailed       Status = "payment_failed"
	StatusPaymentError        Status = "payment_error"
	StatusReleased            Status = "released"
	StatusCanceled            Status = "canceled"
)
