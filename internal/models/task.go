package models


type Operation string


const (
	OperationAdd      Operation = "ADD"      
	OperationSubtract Operation = "SUBTRACT" 
	OperationMultiply Operation = "MULTIPLY" 
	OperationDivide   Operation = "DIVIDE"   
)


type Task struct {
	ID           string    `json:"id"`                      
	Expression   string    `json:"expression,omitempty"`    
	Arg1         string    `json:"arg1,omitempty"`        
	Arg2         string    `json:"arg2,omitempty"`          
	Operation    Operation `json:"operation,omitempty"`    
	OperationTime int      `json:"operation_time,omitempty"`
	Result       *float64  `json:"result,omitempty"`    
	Dependencies []string  `json:"-"`                      
	IsReady      bool      `json:"-"`                       
}

type TaskResponse struct {
	Task *Task `json:"task,omitempty"`
}

type TaskResultRequest struct {
	ID     string  `json:"id"`   
	Result float64 `json:"result"` 
}
