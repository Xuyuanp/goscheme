package lexer

import "fmt"

type precedure func(tokens ...*Token) (*Token, error)

var precedureMap = map[string]precedure{
	"+": add,
}

func doPrecedure(tokens ...*Token) (*Token, error) {
	token := tokens[0]
	if token.typ != typePrecedure {

	}
	result, err := token.val.(precedure)(tokens[1:]...)
	return result, err
}

func add(args ...*Token) (*Token, error) {
	result := &Token{
		typ: typeInteger,
		val: 0.0,
	}

	for _, arg := range args {
		if arg.typ == typeInteger {
			result.val = result.val.(float64) + float64(arg.val.(int))
		} else if arg.typ == typeFloat {
			result.val = result.val.(float64) + arg.val.(float64)
			result.typ = typeFloat
		} else {
			return nil, fmt.Errorf("wrong type arg [%v:%s] for +", arg.raw, arg.typ)
		}
	}
	if result.typ == typeInteger {
		result.val = int(result.val.(float64))
	}

	return result, nil
}
