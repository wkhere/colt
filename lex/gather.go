package lex

func (ts TokenStream) Gather() (res []Token, err error) {
	for tok := range ts {
		if tok.Type == TokenError {
			return res, tok.Err
		}
		res = append(res, tok)
	}
	return
}
