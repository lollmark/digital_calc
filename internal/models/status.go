package models

type Status string

const (
	StatusPending    Status = "PENDING"    
	StatusProcessing Status = "PROCESSING"
	StatusCompleted  Status = "COMPLETED"  
	StatusError      Status = "ERROR"      
)
