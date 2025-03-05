package models

type Expression struct {
	ID       string  `json:"id"`                  
	RawExpr  string  `json:"expression,omitempty"` 
	Status   Status  `json:"status"`               
	Result   *string `json:"result,omitempty"`    
	ErrorMsg string  `json:"error,omitempty"`      
}

type ExpressionRequest struct {
	Expression string `json:"expression"` 
}


type ExpressionResponse struct {
	ID string `json:"id"` 
}


type ExpressionsResponse struct {
	Expressions []Expression `json:"expressions"` 
}


type ExpressionDetailResponse struct {
	Expression Expression `json:"expression"` 
}
