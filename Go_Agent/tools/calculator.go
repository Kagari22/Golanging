package tools

import (
	"context"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"strconv"
)

type CalculatorTool struct{}

func NewCalculatorTool() *CalculatorTool {
	return &CalculatorTool{}
}

func (c *CalculatorTool) Name() string {
	return "calculator"
}

func (c *CalculatorTool) Description() string {
	return "Evaluates arithmetic expressions like 12*(3+4)-5."
}

func (c *CalculatorTool) Run(_ context.Context, input string) (string, error) {
	expr, err := parser.ParseExpr(input)
	if err != nil {
		return "", fmt.Errorf("invalid expression: %w", err)
	}

	value, err := evalExpr(expr)
	if err != nil {
		return "", err
	}

	return strconv.FormatFloat(value, 'f', -1, 64), nil
}

func evalExpr(expr ast.Expr) (float64, error) {
	switch node := expr.(type) {
	case *ast.BasicLit:
		if node.Kind != token.INT && node.Kind != token.FLOAT {
			return 0, fmt.Errorf("unsupported literal: %s", node.Value)
		}
		return strconv.ParseFloat(node.Value, 64)
	case *ast.BinaryExpr:
		left, err := evalExpr(node.X)
		if err != nil {
			return 0, err
		}
		right, err := evalExpr(node.Y)
		if err != nil {
			return 0, err
		}

		switch node.Op {
		case token.ADD:
			return left + right, nil
		case token.SUB:
			return left - right, nil
		case token.MUL:
			return left * right, nil
		case token.QUO:
			if right == 0 {
				return 0, fmt.Errorf("division by zero")
			}
			return left / right, nil
		default:
			return 0, fmt.Errorf("unsupported operator: %s", node.Op)
		}
	case *ast.ParenExpr:
		return evalExpr(node.X)
	case *ast.UnaryExpr:
		value, err := evalExpr(node.X)
		if err != nil {
			return 0, err
		}
		switch node.Op {
		case token.ADD:
			return value, nil
		case token.SUB:
			return -value, nil
		default:
			return 0, fmt.Errorf("unsupported unary operator: %s", node.Op)
		}
	default:
		return 0, fmt.Errorf("unsupported expression type %T", expr)
	}
}
