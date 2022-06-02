package lex

func (ts TokenStream) Gather() (res []Token, err error) {
	for tok := range ts {
		res = append(res, tok)
		if tok.Type == TokenError {
			return res, tok.Err
		}
	}
	return
}
