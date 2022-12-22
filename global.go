package bargle

// Convenience function to parse a long using builtin unmarshaller. TODO: *Parser should be relaxed
// to a Parser interface when it can be separated out.
func ParseLongBuiltin[U BuiltinUnmarshalerType](p *Parser, value *U, elem string, elems ...string) bool {
	return p.Parse(LongElems(BuiltinUnmarshaler(value), elem, elems...))
}
