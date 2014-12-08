// Read and write JSON (RFC 4627) encoded values token by token.
//
// Based on the Android JSON utilities:
//   http://developer.android.com/reference/android/util/JsonReader.html
//   http://developer.android.com/reference/android/util/JsonToken.html
//   http://developer.android.com/reference/android/util/JsonWriter.html
package rgo

import (
	"bufio"
	"bytes"
	"errors"
	"io"
	"math"
	"strconv"
)

// A structure, name or value type in a JSON-encoded string.
type Token int

const (
	NO_TOKEN = iota
	// The opening of a JSON array.
	BEGIN_ARRAY
	// The opening of a JSON object.
	BEGIN_OBJECT
	// A JSON true or false.
	BOOLEAN
	// The closing of a JSON array.
	END_ARRAY
	// The end of the JSON stream.
	END_DOCUMENT
	// The closing of a JSON object.
	END_OBJECT
	// A JSON property name.
	NAME
	// A JSON null.
	NULL
	// A JSON number.
	NUMBER
	// A JSON string.
	STRING
)

var (
	IllegalState    = errors.New("rgo: Illegal state")
	IllegalArgument = errors.New("rgo: Illegal argument")
	InvalidInput    = errors.New("rgo: Invalid input")
)

// Write a JSON (RFC 4627) encoded value to a Writer, one token at a time.
type Writer struct {
	w            io.Writer
	pendingComma bool
	state        []Token
}

// Create a new instance that writes a JSON-encoded stream to w.
func NewWriter(w io.Writer) *Writer {
	return &Writer{w: w}
}

func (w *Writer) beginValue() error {
	if len(w.state) > 0 {
		switch w.state[len(w.state)-1] {
		case BEGIN_ARRAY:
			// nothing to be done
		case NAME:
			w.state = w.state[:len(w.state)-1]
			if len(w.state) == 0 || w.state[len(w.state)-1] != BEGIN_OBJECT {
				panic("rgo: Internal error")
			}
		case BEGIN_OBJECT:
			return IllegalState
		default:
			panic("rgo: Internal error")
		}
	}
	if w.pendingComma {
		if _, err := io.WriteString(w.w, ","); err != nil {
			return err
		}
	} else {
		w.pendingComma = true
	}
	return nil
}

// Begin encoding a new array.
func (w *Writer) BeginArray() error {
	if err := w.beginValue(); err != nil {
		return err
	}
	w.pendingComma = false
	w.state = append(w.state, BEGIN_ARRAY)
	if _, err := io.WriteString(w.w, "["); err != nil {
		return err
	}
	return nil
}

// Begin encoding a new object.
func (w *Writer) BeginObject() error {
	if err := w.beginValue(); err != nil {
		return err
	}
	w.pendingComma = false
	w.state = append(w.state, BEGIN_OBJECT)
	if _, err := io.WriteString(w.w, "{"); err != nil {
		return err
	}
	return nil
}

// End encoding the current array.
func (w *Writer) EndArray() error {
	if len(w.state) == 0 || w.state[len(w.state)-1] != BEGIN_ARRAY {
		return IllegalState
	}
	w.pendingComma = true
	w.state = w.state[:len(w.state)-1]
	if _, err := io.WriteString(w.w, "]"); err != nil {
		return err
	}
	return nil
}

// End encoding the current object.
func (w *Writer) EndObject() error {
	if len(w.state) == 0 || w.state[len(w.state)-1] != BEGIN_OBJECT {
		return IllegalState
	}
	w.pendingComma = true
	w.state = w.state[:len(w.state)-1]
	if _, err := io.WriteString(w.w, "}"); err != nil {
		return err
	}
	return nil
}

func writeQuotedString(w *Writer, s string) error {
	if _, err := io.WriteString(w.w, `"`); err != nil {
		return err
	}
	for _, ch := range s {
		switch {
		case ch == 0x08:
			if _, err := io.WriteString(w.w, "\\b"); err != nil {
				return err
			}
		case ch == 0x09:
			if _, err := io.WriteString(w.w, "\\t"); err != nil {
				return err
			}
		case ch == 0x0a:
			if _, err := io.WriteString(w.w, "\\n"); err != nil {
				return err
			}
		case ch == 0x0c:
			if _, err := io.WriteString(w.w, "\\f"); err != nil {
				return err
			}
		case ch == 0x0d:
			if _, err := io.WriteString(w.w, "\\r"); err != nil {
				return err
			}
		case ch == 0x22:
			if _, err := io.WriteString(w.w, "\\\""); err != nil {
				return err
			}
		case ch == 0x5c:
			if _, err := io.WriteString(w.w, "\\\\"); err != nil {
				return err
			}
		case ch < 0x10:
			if _, err := io.WriteString(w.w, "\\u000"); err != nil {
				return err
			}
			if _, err := io.WriteString(w.w, strconv.FormatInt(int64(ch), 16)); err != nil {
				return err
			}
		case ch < 0x20:
			if _, err := io.WriteString(w.w, "\\u00"); err != nil {
				return err
			}
			if _, err := io.WriteString(w.w, strconv.FormatInt(int64(ch), 16)); err != nil {
				return err
			}
		case ch < 0x110000:
			if _, err := io.WriteString(w.w, string(ch)); err != nil {
				return err
			}
		default:
			panic("rgo: Not implemented")
		}
	}
	if _, err := io.WriteString(w.w, `"`); err != nil {
		return err
	}
	return nil
}

// Encode the property name.
func (w *Writer) Name(name string) error {
	if len(w.state) == 0 || w.state[len(w.state)-1] != BEGIN_OBJECT {
		return IllegalState
	}
	w.state = append(w.state, NAME)
	if w.pendingComma {
		if _, err := io.WriteString(w.w, ","); err != nil {
			return err
		}
		w.pendingComma = false
	}
	if err := writeQuotedString(w, name); err != nil {
		return err
	}
	if _, err := io.WriteString(w.w, ":"); err != nil {
		return err
	}
	return nil
}

// Encode null.
func (w *Writer) NullValue() error {
	if err := w.beginValue(); err != nil {
		return err
	}
	if _, err := io.WriteString(w.w, "null"); err != nil {
		return err
	}
	return nil
}

// Encode value.
func (w *Writer) IntValue(value int) error {
	return w.Int64Value(int64(value))
}

// Encode value.
func (w *Writer) Int8Value(value int8) error {
	return w.Int64Value(int64(value))
}

// Encode value.
func (w *Writer) Int16Value(value int16) error {
	return w.Int64Value(int64(value))
}

// Encode value.
func (w *Writer) Int32Value(value int32) error {
	return w.Int64Value(int64(value))
}

// Encode value.
func (w *Writer) Int64Value(value int64) error {
	if err := w.beginValue(); err != nil {
		return err
	}
	if _, err := io.WriteString(w.w, strconv.FormatInt(value, 10)); err != nil {
		return err
	}
	return nil
}

// Encode value.
func (w *Writer) UintValue(value uint) error {
	return w.Uint64Value(uint64(value))
}

// Encode value.
func (w *Writer) Uint8Value(value uint8) error {
	return w.Uint64Value(uint64(value))
}

// Encode value.
func (w *Writer) Uint16Value(value uint16) error {
	return w.Uint64Value(uint64(value))
}

// Encode value.
func (w *Writer) Uint32Value(value uint32) error {
	return w.Uint64Value(uint64(value))
}

// Encode value.
func (w *Writer) Uint64Value(value uint64) error {
	if err := w.beginValue(); err != nil {
		return err
	}
	if _, err := io.WriteString(w.w, strconv.FormatUint(value, 10)); err != nil {
		return err
	}
	return nil
}

// Encode value.
func (w *Writer) Float32Value(value float32) error {
	return w.floatValue(float64(value), 32)
}

// Encode value.
func (w *Writer) Float64Value(value float64) error {
	return w.floatValue(value, 64)
}

func (w *Writer) floatValue(value float64, bitSize int) error {
	if math.IsInf(value, 0) || math.IsNaN(value) {
		return IllegalArgument
	}
	if err := w.beginValue(); err != nil {
		return err
	}
	if _, err := io.WriteString(w.w, strconv.FormatFloat(value, 'g', -1, bitSize)); err != nil {
		return err
	}
	return nil
}

// Encode value.
func (w *Writer) BoolValue(value bool) error {
	if err := w.beginValue(); err != nil {
		return err
	}
	if value {
		if _, err := io.WriteString(w.w, "true"); err != nil {
			return err
		}
	} else {
		if _, err := io.WriteString(w.w, "false"); err != nil {
			return err
		}
	}
	return nil
}

// Encode value.
func (w *Writer) StringValue(value string) error {
	if err := w.beginValue(); err != nil {
		return err
	}
	if err := writeQuotedString(w, value); err != nil {
		return err
	}
	return nil
}

// Encode value.
func (w *Writer) Value(value interface{}) error {
	switch v := value.(type) {
	case nil:
		return w.NullValue()
	case int:
		return w.IntValue(v)
	case int8:
		return w.Int8Value(v)
	case int16:
		return w.Int16Value(v)
	case int32:
		return w.Int32Value(v)
	case int64:
		return w.Int64Value(v)
	case uint:
		return w.UintValue(v)
	case uint8:
		return w.Uint8Value(v)
	case uint16:
		return w.Uint16Value(v)
	case uint32:
		return w.Uint32Value(v)
	case uint64:
		return w.Uint64Value(v)
	case float32:
		return w.Float32Value(v)
	case float64:
		return w.Float64Value(v)
	case bool:
		return w.BoolValue(v)
	case string:
		return w.StringValue(v)
	default:
		return IllegalArgument
	}
}

// Read a JSON (RFC 4627) encoded value as a stream of tokens.
type Reader struct {
	r       *bufio.Reader
	token   Token
	value   bytes.Buffer
	hasNext bool
}

// Create a new instance that reads a JSON-encoded stream from r.
func NewReader(r io.Reader) *Reader {
	return &Reader{r: bufio.NewReaderSize(r, 6)}
}

func (r *Reader) skipWhitespace() error {
	for {
		b, err := r.r.ReadByte()
		if err != nil {
			return err
		}
		switch b {
		case 0x20, 0x09, 0x0a, 0x0d:
		default:
			return r.r.UnreadByte()
		}
	}
}

func (r *Reader) skipLiteral(literal string) error {
	for _, ch := range []byte(literal) {
		if b, err := r.r.ReadByte(); err != nil {
			return err
		} else if b != ch {
			return InvalidInput
		}
	}
	b, err := r.r.ReadByte()
	if err != nil {
		if err == io.EOF {
			r.hasNext = false
			return nil
		}
		return err
	}
	switch b {
	case 0x20, 0x09, 0x0a, 0x0d:
		return r.readTokenEnd()
	case ',':
		r.hasNext = true
		return nil
	case ']', '}':
		r.hasNext = false
		return r.r.UnreadByte()
	default:
		return InvalidInput
	}
}

func (r *Reader) readTokenEnd() error {
	for {
		b, err := r.r.ReadByte()
		if err != nil {
			if err == io.EOF {
				r.hasNext = false
				return nil
			}
			return err
		}
		switch b {
		case 0x20, 0x09, 0x0a, 0x0d:
		case ',':
			r.hasNext = true
			return nil
		default:
			r.hasNext = false
			return r.r.UnreadByte()
		}
	}
}

func (r *Reader) readContainerStart(containerEnd byte) error {
	for {
		b, err := r.r.ReadByte()
		if err != nil {
			if err == io.EOF {
				return io.ErrUnexpectedEOF
			}
			return err
		}
		switch b {
		case 0x20, 0x09, 0x0a, 0x0d:
		default:
			r.hasNext = b != containerEnd
			return r.r.UnreadByte()
		}
	}
}

func (r *Reader) readToken(skipValue bool) error {
	if r.token != NO_TOKEN {
		panic("rgo: Internal error")
	}
	if err := r.skipWhitespace(); err != nil {
		return err
	}
	b, err := r.r.ReadByte()
	if err != nil {
		if err == io.EOF {
			r.token = END_DOCUMENT
			return nil
		}
		return err
	}
	r.value.Reset()
	r.hasNext = false
	switch b {
	case '[':
		r.token = BEGIN_ARRAY
		return r.readContainerStart(']')
	case ']':
		r.token = END_ARRAY
		return r.readTokenEnd()
	case '{':
		r.token = BEGIN_OBJECT
		return r.readContainerStart('}')
	case '}':
		r.token = END_OBJECT
		return r.readTokenEnd()
	case 'f':
		r.token = BOOLEAN
		if !skipValue {
			r.value.WriteByte(0)
		}
		return r.skipLiteral("alse")
	case 't':
		r.token = BOOLEAN
		if !skipValue {
			r.value.WriteByte(1)
		}
		return r.skipLiteral("rue")
	case 'n':
		r.token = NULL
		return r.skipLiteral("ull")
	case '"':
		return r.readStringOrName(skipValue)
	case '-':
		r.token = NUMBER
		if err := r.value.WriteByte(b); err != nil {
			return err
		}
		return r.readNumber(skipValue, true, false, false)
	case '0':
		r.token = NUMBER
		if err := r.value.WriteByte(b); err != nil {
			return err
		}
		return r.readNumber(skipValue, false, true, true)
	case '1', '2', '3', '4', '5', '6', '7', '8', '9':
		r.token = NUMBER
		if err := r.value.WriteByte(b); err != nil {
			return err
		}
		return r.readNumber(skipValue, false, false, false)
	default:
		return InvalidInput
	}
}

func (r *Reader) readStringOrName(skipValue bool) error {
loop:
	for {
		b, err := r.r.ReadByte()
		if err != nil {
			if err == io.EOF {
				return io.ErrUnexpectedEOF
			}
			return err
		}
		switch b {
		case '"':
			break loop
		case '\\':
			b, err = r.r.ReadByte()
			if err != nil {
				if err == io.EOF {
					return io.ErrUnexpectedEOF
				}
				return err
			}
			switch b {
			case '"', '\\', '/':
				if !skipValue {
					if err := r.value.WriteByte(b); err != nil {
						return err
					}
				}
			case 'b':
				if !skipValue {
					if err := r.value.WriteByte(8); err != nil {
						return err
					}
				}
			case 'f':
				if !skipValue {
					if err := r.value.WriteByte(12); err != nil {
						return err
					}
				}
			case 'n':
				if !skipValue {
					if err := r.value.WriteByte(10); err != nil {
						return err
					}
				}
			case 'r':
				if !skipValue {
					if err := r.value.WriteByte(13); err != nil {
						return err
					}
				}
			case 't':
				if !skipValue {
					if err := r.value.WriteByte(9); err != nil {
						return err
					}
				}
			case 'u':
				var buf [4]byte
				if _, err := r.r.Read(buf[:]); err != nil {
					return err
				}
				codePoint, err := strconv.ParseUint(string(buf[:]), 16, 16)
				if err != nil {
					return err
				}
				if codePoint >= 0xdc00 && codePoint < 0xe000 {
					return InvalidInput
				} else if codePoint >= 0xd800 && codePoint < 0xdc00 {
					if b, err := r.r.ReadByte(); err != nil {
						return err
					} else if b != '\\' {
						return InvalidInput
					}
					if b, err := r.r.ReadByte(); err != nil {
						return err
					} else if b != 'u' {
						return InvalidInput
					}
					if _, err := r.r.Read(buf[:]); err != nil {
						return err
					}
					lowSurrogate, err := strconv.ParseUint(string(buf[:]), 16, 16)
					if err != nil {
						return err
					}
					codePoint = 0x10000 + ((codePoint & 0x3ff) << 10) + (lowSurrogate & 0x3ff)
				}
				if !skipValue {
					if _, err := r.value.WriteRune(rune(codePoint)); err != nil {
						return err
					}
				}
			}
		default:
			if b < 0x20 {
				return InvalidInput
			}
			if !skipValue {
				if err := r.value.WriteByte(b); err != nil {
					return err
				}
			}
		}
	}
	for {
		b, err := r.r.ReadByte()
		if err != nil {
			if err == io.EOF {
				r.token = STRING
				r.hasNext = false
				return nil
			}
			return err
		}
		switch b {
		case 0x20, 0x09, 0x0a, 0x0d:
		case ',':
			r.token = STRING
			r.hasNext = true
			return nil
		case ':':
			r.token = NAME
			r.hasNext = false
			return nil
		case ']', '}':
			r.token = STRING
			r.hasNext = false
			return r.r.UnreadByte()
		default:
			return InvalidInput
		}
	}
}

func (r *Reader) readNumber(skipValue, digitNeeded, leadingZero, intDone bool) error {
	fracDone := false
	signPossible := false
	for {
		b, err := r.r.ReadByte()
		if err != nil {
			if err == io.EOF {
				return nil
			}
			return err
		}
		switch b {
		case '0':
			if digitNeeded {
				leadingZero = true
				digitNeeded = false
			} else if leadingZero && !intDone {
				return InvalidInput
			} else {
				leadingZero = false
			}
		case '1', '2', '3', '4', '5', '6', '7', '8', '9':
			leadingZero = false
			digitNeeded = false
			signPossible = false
		case '.':
			if intDone {
				return InvalidInput
			}
			intDone = true
			digitNeeded = true
		case 'e', 'E':
			if fracDone {
				return InvalidInput
			}
			intDone = true
			fracDone = true
			digitNeeded = true
			signPossible = true
		case '+', '-':
			if !signPossible {
				return InvalidInput
			}
			signPossible = false
		default:
			if digitNeeded {
				return InvalidInput
			}
			if err := r.r.UnreadByte(); err != nil {
				return err
			}
			return r.readTokenEnd()
		}
		if !skipValue {
			if err := r.value.WriteByte(b); err != nil {
				return err
			}
		}
	}
}

// Consume the next token from the JSON stream and asserts that it is the
// beginning of a new array.
func (r *Reader) BeginArray() error {
	if r.token == NO_TOKEN {
		if err := r.readToken(false); err != nil {
			return err
		}
	}
	if r.token == BEGIN_ARRAY {
		r.token = NO_TOKEN
		return nil
	}
	return IllegalState
}

// Consume the next token from the JSON stream and asserts that it is the
// beginning of a new object.
func (r *Reader) BeginObject() error {
	if r.token == NO_TOKEN {
		if err := r.readToken(false); err != nil {
			return err
		}
	}
	if r.token == BEGIN_OBJECT {
		r.token = NO_TOKEN
		return nil
	}
	return IllegalState
}

// Consume the next token from the JSON stream and asserts that it is the
// end of the current array.
func (r *Reader) EndArray() error {
	if r.token == NO_TOKEN {
		if err := r.readToken(false); err != nil {
			return err
		}
	}
	if r.token == END_ARRAY {
		r.token = NO_TOKEN
		return nil
	}
	return IllegalState
}

// Consume the next token from the JSON stream and asserts that it is the
// end of the current object.
func (r *Reader) EndObject() error {
	if r.token == NO_TOKEN {
		if err := r.readToken(false); err != nil {
			return err
		}
	}
	if r.token == END_OBJECT {
		r.token = NO_TOKEN
		return nil
	}
	return IllegalState
}

// Return true if the current array or object has another element.
func (r *Reader) HasNext() (bool, error) {
	return r.hasNext, nil
}

// Return the boolean value of the next token, consuming it.
func (r *Reader) NextBoolean() (bool, error) {
	if r.token == NO_TOKEN {
		if err := r.readToken(false); err != nil {
			return false, err
		}
	}
	if r.token == BOOLEAN {
		r.token = NO_TOKEN
		if value, err := r.r.Peek(1); err != nil {
			panic("rgo: Internal error")
		} else {
			return value[0] != 0, nil
		}
	}
	return false, IllegalState
}

// Return the float32 value of the next token, consuming it.  If the next
// token is a string, this method will attempt to parse it as a float32.
func (r *Reader) NextFloat32() (float32, error) {
	value, err := r.nextFloat(32)
	return float32(value), err
}

// Return the float64 value of the next token, consuming it.  If the next
// token is a string, this method will attempt to parse it as a float64.
func (r *Reader) NextFloat64() (float64, error) {
	return r.nextFloat(64)
}

func (r *Reader) nextFloat(bitSize int) (float64, error) {
	if r.token == NO_TOKEN {
		if err := r.readToken(false); err != nil {
			return 0, err
		}
	}
	switch r.token {
	case STRING, NUMBER:
		r.token = NO_TOKEN
		return strconv.ParseFloat(r.value.String(), bitSize)
	default:
		return 0, IllegalState
	}
}

// Return the int value of the next token, consuming it.  If the next
// token is a string, this method will attempt to parse it as a int.
func (r *Reader) NextInt() (int, error) {
	value, err := r.NextInt64()
	return int(value), err
}

// Return the int64 value of the next token, consuming it.  If the next
// token is a string, this method will attempt to parse it as a int64.
func (r *Reader) NextInt64() (int64, error) {
	if r.token == NO_TOKEN {
		if err := r.readToken(false); err != nil {
			return 0, err
		}
	}
	switch r.token {
	case STRING, NUMBER:
		r.token = NO_TOKEN
		return strconv.ParseInt(r.value.String(), 10, 64)
	default:
		return 0, IllegalState
	}
}

// Return the next token, a property name, consuming it.
func (r *Reader) NextName() (string, error) {
	if r.token == NO_TOKEN {
		if err := r.readToken(false); err != nil {
			return "", err
		}
	}
	if r.token == NAME {
		r.token = NO_TOKEN
		return r.value.String(), nil
	}
	return "", IllegalState
}

// Consume the next token from the JSON stream and assert that it is a
// literal null.
func (r *Reader) NextNull() error {
	if r.token == NO_TOKEN {
		if err := r.readToken(false); err != nil {
			return err
		}
	}
	if r.token == NULL {
		r.token = NO_TOKEN
		return nil
	}
	return IllegalState
}

// Return the string value of the next token, consuming it.  If the next
// token is a number, this method will return its string form.
func (r *Reader) NextString() (string, error) {
	if r.token == NO_TOKEN {
		if err := r.readToken(false); err != nil {
			return "", err
		}
	}
	switch r.token {
	case STRING, NUMBER:
		r.token = NO_TOKEN
		return r.value.String(), nil
	default:
		return "", IllegalState
	}
}

// Return the type of the next token without consuming it.
func (r *Reader) Peek() (Token, error) {
	if r.token == NO_TOKEN {
		if err := r.readToken(false); err != nil {
			return r.token, err
		}
	}
	if r.token == NO_TOKEN {
		panic("rgo: Internal error")
	}
	return r.token, nil
}

// Skip the next value recursively.  If it is an object or array, all
// nested elements are skipped.  This method is intended for use when
// the JSON token stream contains unrecognized or unhandled values.
func (r *Reader) SkipValue() error {
	if r.token == NO_TOKEN {
		if err := r.readToken(true); err != nil {
			return err
		}
	}
	if r.token == NAME {
		return IllegalState
	}
	for nesting := 0; ; {
		switch r.token {
		case NO_TOKEN:
			panic("rgo: Internal error")
		case BEGIN_ARRAY, BEGIN_OBJECT:
			nesting++
		case END_ARRAY, END_OBJECT:
			if nesting <= 0 {
				return InvalidInput
			}
			nesting--
		default:
		}
		r.token = NO_TOKEN
		if nesting <= 0 {
			return nil
		}
		if err := r.readToken(true); err != nil {
			return err
		}
	}
}
