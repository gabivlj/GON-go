/// This implementation of https://github.com/TylerGlaiel/GON
/// is super duper simple, not expected for real use case, but to show that
/// we can be so much faster than json decoding on Go without effort.
/// It doesn't support JSON syntax as GON does promise.
/// Also error checks are limited in favour of performance
package gon

import (
	"errors"
	"fmt"
	"strconv"
)

func SerializeString(values map[string]any) string {
	panic("UNIMPLEMENTED: 🤨 are you sure chief?")
}

func DeserializeString(gon string) (map[string]any, error) {
	parser := parser{buffer: gon, index: 0}
	return parser.parseGonObject(), parser.err
}

type parser struct {
	index  int
	buffer string
	err    error
}

type tokenType = int

const (
	integer tokenType = iota
	float
	str
	leftBrace
	rightBrace
	leftBracket
	rightBracket
	eof
	trueKeyword
	falseKeyword
)

var eofToken = token{literal: "", kind: eof}
var leftBracketToken = token{literal: "[", kind: leftBracket}
var rightBracketToken = token{literal: "]", kind: rightBracket}
var leftBraceToken = token{literal: "{", kind: leftBrace}
var rightBraceToken = token{literal: "}", kind: rightBrace}

type token struct {
	kind    tokenType
	literal string
}

func (p *parser) error(s string) {
	if p.err != nil {
		return
	}
	p.err = errors.New(s)
}

func (p *parser) isWhitespace() bool {
	return p.buffer[p.index] == ' ' || p.buffer[p.index] == '\n' || p.buffer[p.index] == '\t' || p.buffer[p.index] == '\r'
}

func (p *parser) isAlphanumeric() bool {
	return !p.isWhitespace() && p.buffer[p.index] != '{' && p.buffer[p.index] != '}' && p.buffer[p.index] != '[' && p.buffer[p.index] != ']'
}

func (p *parser) handleTokenKeyword(index int, offset int, keyword string, kind tokenType) token {
	keyword, keywordIndex := keyword[offset:], 0
	for p.index < len(p.buffer) && p.isAlphanumeric() {
		if keywordIndex < len(keyword) && p.buffer[p.index] == keyword[keywordIndex] {
			keywordIndex++
		} else {
			keywordIndex = 100
		}
		p.index++
	}

	if keywordIndex == len(keyword) {
		return token{literal: keyword, kind: kind}
	}
	return token{literal: p.buffer[index:p.index]}
}

func (p *parser) nextToken() token {
	// skip whitespace
	for p.index < len(p.buffer) && p.isWhitespace() {
		p.index++
	}

	if len(p.buffer) <= p.index {
		return eofToken
	}

	index := p.index
	p.index++

	switch p.buffer[index] {
	case '[':
		return leftBracketToken
	case '{':
		return leftBraceToken
	case ']':
		return rightBracketToken
	case '}':
		return rightBraceToken
	case '"':
		for p.index < len(p.buffer) && p.buffer[p.index] != '"' {
			p.index++
		}
		p.index++
		return token{literal: p.buffer[index+1 : p.index-1], kind: str}
	case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
		// this trick will error out on the caller, it's intended, if someone wants a string literal integer they should use ""
		kind := integer
		for p.index < len(p.buffer) && ((p.buffer[p.index] >= '0' && p.buffer[p.index] <= '9') || p.buffer[p.index] == '.') {
			if p.buffer[p.index] == '.' {
				kind = float
			}
			p.index++
		}
		return token{literal: p.buffer[index:p.index], kind: kind}

	// consider this as an in-memory trie
	case 't':
		return p.handleTokenKeyword(index, 1, "true", trueKeyword)

	case 'f':
		return p.handleTokenKeyword(index, 1, "false", falseKeyword)

	default:
		for p.index < len(p.buffer) && p.isAlphanumeric() {
			p.index++
		}
	}

	return token{literal: p.buffer[index:p.index], kind: str}
}

func (p *parser) parseSingleObject(context token) any {
	var value any
	var err error
	next := p.nextToken()
	switch next.kind {
	case trueKeyword:
		value = true
	case falseKeyword:
		value = false
	case rightBrace:
		return nil
	case rightBracket:
		return nil
	case str:
		value = next.literal
	case leftBrace:
		value = p.parseGonObject()
	case leftBracket:
		value = p.parseGonList()
	case integer:
		value, err = strconv.ParseInt(next.literal, 10, 64)
		if err != nil {
			p.error("error parsing integer")
		}
	case float:
		value, err = strconv.ParseFloat(next.literal, 64)
		if err != nil {
			p.error("error parsing float")
		}
	default:
		p.error(fmt.Sprintf("unexpected token '%v' on context '%v'", next.literal, context.literal))
	}

	return value
}

func (p *parser) parseGonObject() map[string]any {
	object := map[string]any{}

	for p.index < len(p.buffer) {
		literal := p.nextToken()
		switch literal.kind {
		case eof:
			return object
		case rightBrace:
			return object
		case rightBracket:
			return nil
		}
		object[literal.literal] = p.parseSingleObject(literal)
	}

	return object
}

func (p *parser) parseGonList() []any {
	elements := make([]any, 0, 5)
	for p.index < len(p.buffer) {
		object := p.parseSingleObject(leftBracketToken)
		if object == nil {
			break
		}
		elements = append(elements, object)
	}
	return elements
}
