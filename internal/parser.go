package application

import (
	"fmt"
	"strconv"
	"strings"
	"unicode"

	"github.com/lollmark/digital_calc/pkg/calculator"
)

type ASTNode struct {
	IsLeaf        bool
	Value         float64
	Operator      string
	Left, Right   *ASTNode
	TaskScheduled bool
}

func ParseAST(expression string) (*ASTNode, error) {
	expr := strings.ReplaceAll(expression, " ", "")
	if expr == "" {
		return nil, fmt.Errorf("empty expression")
	}
	p := &parser{input: expr, pos: 0}
	node, err := p.parseExpression()
	if err != nil {
		return nil, err
	}
	if p.pos < len(p.input) {
		return nil, fmt.Errorf("unexpected token at position %d", p.pos)
	}
	return node, nil
}

type parser struct {
	input string
	pos   int
}

func (p *parser) peek() rune {
	if p.pos < len(p.input) {
		return rune(p.input[p.pos])
	}
	return 0
}

func (p *parser) get() rune {
	ch := p.peek()
	p.pos++
	return ch
}

func (p *parser) parseExpression() (*ASTNode, error) {
	node, err := p.parseTerm()
	if err != nil {
		return nil, err
	}
	for {
		ch := p.peek()
		if ch == '+' || ch == '-' {
			op := string(p.get())
			right, err := p.parseTerm()
			if err != nil {
				return nil, err
			}
			node = &ASTNode{
				IsLeaf:   false,
				Operator: op,
				Left:     node,
				Right:    right,
			}
		} else {
			break
		}
	}
	return node, nil
}

func (p *parser) parseTerm() (*ASTNode, error) {
	node, err := p.parseFactor()
	if err != nil {
		return nil, err
	}
	for {
		ch := p.peek()
		if ch == '*' || ch == '/' {
			op := string(p.get())
			right, err := p.parseFactor()
			if err != nil {
				return nil, err
			}
			node = &ASTNode{
				IsLeaf:   false,
				Operator: op,
				Left:     node,
				Right:    right,
			}
		} else {
			break
		}
	}
	return node, nil
}

func (p *parser) parseFactor() (*ASTNode, error) {
	ch := p.peek()
	if ch == '(' {
		p.get()
		node, err := p.parseExpression()
		if err != nil {
			return nil, err
		}
		if p.peek() != ')' {
			return nil, fmt.Errorf("missing closing parenthesis")
		}
		p.get()
		return node, nil
	}
	start := p.pos
	// Обрабатываем знак: унарный плюс разрешаем только если он стоит в начале или сразу после '('
	if ch == '+' {
		if p.pos > 0 && p.input[p.pos-1] != '(' {
			return nil, fmt.Errorf("unexpected unary plus at position %d", p.pos)
		}
		p.get()
	} else if ch == '-' {
		p.get()
	}
	for {
		ch = p.peek()
		if unicode.IsDigit(ch) || ch == '.' {
			p.get()
		} else {
			break
		}
	}
	token := p.input[start:p.pos]
	if token == "" {
		return nil, fmt.Errorf("expected number at position %d", start)
	}
	value, err := strconv.ParseFloat(token, 64)
	if err != nil {
		return nil, fmt.Errorf("invalid number %s", token)
	}
	return &ASTNode{
		IsLeaf: true,
		Value:  value,
	}, nil
}

// В конце файла ast.go, после ParseAST и парсера:
func EvalAST(node *ASTNode) (float64, error) {
	if node.IsLeaf {
		return node.Value, nil
	}
	left, err := EvalAST(node.Left)
	if err != nil {
		return 0, err
	}
	right, err := EvalAST(node.Right)
	if err != nil {
		return 0, err
	}
	if node.Operator == "/" && right == 0 {
		return 0, fmt.Errorf("division by zero")
	}
	return calculation.Compute(node.Operator, left, right)
}
