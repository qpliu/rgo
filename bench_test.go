package rgo

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"io/ioutil"
	"os"
	"testing"
)

type codeResponse struct {
	Tree     *codeNode `json:"tree"`
	Username string    `json:"username"`
}

type codeNode struct {
	Name     string      `json:"name"`
	Kids     []*codeNode `json:"kids"`
	CLWeight float64     `json:"cl_weight"`
	Touches  int         `json:"touches"`
	MinT     int64       `json:"min_t"`
	MaxT     int64       `json:"max_t"`
	MeanT    int64       `json:"mean_t"`
}

var codeJSON []byte
var codeStruct codeResponse

func decodeResponse(response *codeResponse, r *Reader) error {
	if err := r.BeginObject(); err != nil {
		return err
	}
	for {
		if hasNext, err := r.HasNext(); err != nil {
			return err
		} else if !hasNext {
			break
		}
		name, err := r.NextName()
		if err != nil {
			return err
		}
		switch name {
		case "tree":
			response.Tree = &codeNode{}
			if err := decodeNode(response.Tree, r); err != nil {
				return err
			}
		case "username":
			response.Username, err = r.NextString()
			if err != nil {
				return err
			}
		default:
			if err := r.SkipValue(); err != nil {
				return err
			}
		}
	}
	if err := r.EndObject(); err != nil {
		return err
	}
	return nil
}

func decodeNode(node *codeNode, r *Reader) error {
	if err := r.BeginObject(); err != nil {
		return err
	}
	for {
		if hasNext, err := r.HasNext(); err != nil {
			return err
		} else if !hasNext {
			break
		}
		name, err := r.NextName()
		if err != nil {
			return err
		}
		switch name {
		case "name":
			node.Name, err = r.NextString()
			if err != nil {
				return err
			}
		case "kids":
			if err := r.BeginArray(); err != nil {
				return err
			}
			for {
				if hasNext, err := r.HasNext(); err != nil {
					return err
				} else if !hasNext {
					break
				}
				kid := &codeNode{}
				if err := decodeNode(kid, r); err != nil {
					return err
				}
				node.Kids = append(node.Kids, kid)
			}
			if err := r.EndArray(); err != nil {
				return err
			}
		case "cl_weight":
			node.CLWeight, err = r.NextFloat64()
			if err != nil {
				return err
			}
		case "touches":
			node.Touches, err = r.NextInt()
			if err != nil {
				return err
			}
		case "min_t":
			node.MinT, err = r.NextInt64()
			if err != nil {
				return err
			}
		case "max_t":
			node.MaxT, err = r.NextInt64()
			if err != nil {
				return err
			}
		case "mean_t":
			node.MeanT, err = r.NextInt64()
			if err != nil {
				return err
			}
		default:
			if err := r.SkipValue(); err != nil {
				return err
			}
		}
	}
	if err := r.EndObject(); err != nil {
		return err
	}
	return nil
}

func encodeResponse(response *codeResponse, w *Writer) error {
	if err := w.BeginObject(); err != nil {
		return err
	}
	if err := w.Name("tree"); err != nil {
		return err
	}
	if err := encodeNode(response.Tree, w); err != nil {
		return err
	}
	if err := w.Name("username"); err != nil {
		return err
	}
	if err := w.StringValue(response.Username); err != nil {
		return err
	}
	if err := w.EndObject(); err != nil {
		return err
	}
	return nil
}

func encodeNode(node *codeNode, w *Writer) error {
	if err := w.BeginObject(); err != nil {
		return err
	}
	if err := w.Name("name"); err != nil {
		return err
	}
	if err := w.StringValue(node.Name); err != nil {
		return err
	}
	if err := w.Name("kids"); err != nil {
		return err
	}
	if err := w.BeginArray(); err != nil {
		return err
	}
	for _, kid := range node.Kids {
		if err := encodeNode(kid, w); err != nil {
			return err
		}
	}
	if err := w.EndArray(); err != nil {
		return err
	}
	if err := w.Name("cl_weight"); err != nil {
		return err
	}
	if err := w.Float64Value(node.CLWeight); err != nil {
		return err
	}
	if err := w.Name("touches"); err != nil {
		return err
	}
	if err := w.IntValue(node.Touches); err != nil {
		return err
	}
	if err := w.Name("min_t"); err != nil {
		return err
	}
	if err := w.Int64Value(node.MinT); err != nil {
		return err
	}
	if err := w.Name("max_t"); err != nil {
		return err
	}
	if err := w.Int64Value(node.MaxT); err != nil {
		return err
	}
	if err := w.Name("mean_t"); err != nil {
		return err
	}
	if err := w.Int64Value(node.MeanT); err != nil {
		return err
	}
	if err := w.EndObject(); err != nil {
		return err
	}
	return nil
}

func codeInit() {
	f, err := os.Open("testdata/code.json.gz")
	if err != nil {
		panic(err)
	}
	defer f.Close()
	gz, err := gzip.NewReader(f)
	if err != nil {
		panic(err)
	}
	data, err := ioutil.ReadAll(gz)
	if err != nil {
		panic(err)
	}

	codeJSON = data

	if err := json.Unmarshal(codeJSON, &codeStruct); err != nil {
		panic("unmarshal code.json: " + err.Error())
	}

	if data, err = json.Marshal(&codeStruct); err != nil {
		panic("marshal code.json: " + err.Error())
	}

	if !bytes.Equal(data, codeJSON) {
		println("different lengths", len(data), len(codeJSON))
		for i := 0; i < len(data) && i < len(codeJSON); i++ {
			if data[i] != codeJSON[i] {
				println("re-marshal: changed at byte", i)
				println("orig: ", string(codeJSON[i-10:i+10]))
				println("new: ", string(data[i-10:i+10]))
				break
			}
		}
		panic("re-marshal code.json: different result")
	}

	var response codeResponse
	if err := decodeResponse(&response, NewReader(bytes.NewBuffer(data))); err != nil {
		panic("decodeResponse: " + err.Error())
	}

	var buf bytes.Buffer
	if err := encodeResponse(&response, NewWriter(&buf)); err != nil {
		panic("encodeResponse: " + err.Error())
	}

	data = buf.Bytes()
	if !bytes.Equal(data, codeJSON) {
		println("encodeResponse: different lengths", len(data), len(codeJSON))
		for i := 0; i < len(data) && i < len(codeJSON); i++ {
			if data[i] != codeJSON[i] {
				println("encodeResponse: changed at byte", i)
				println("orig: ", string(codeJSON[i-10:i+10]))
				println("new: ", string(data[i-10:i+10]))
				break
			}
		}
		panic("encodeResponse: different result")
	}
}

func BenchmarkCodeEncoderJson(b *testing.B) {
	if codeJSON == nil {
		b.StopTimer()
		codeInit()
		b.StartTimer()
	}
	enc := json.NewEncoder(ioutil.Discard)
	for i := 0; i < b.N; i++ {
		if err := enc.Encode(&codeStruct); err != nil {
			b.Fatal("Encode:", err)
		}
	}
	b.SetBytes(int64(len(codeJSON)))
}

func BenchmarkCodeEncoderRgo(b *testing.B) {
	if codeJSON == nil {
		b.StopTimer()
		codeInit()
		b.StartTimer()
	}
	w := NewWriter(ioutil.Discard)
	for i := 0; i < b.N; i++ {
		if err := encodeResponse(&codeStruct, w); err != nil {
			b.Fatal("encodeResponse:", err)
		}
	}
	b.SetBytes(int64(len(codeJSON)))
}

func BenchmarkCodeDecoderJson(b *testing.B) {
	if codeJSON == nil {
		b.StopTimer()
		codeInit()
		b.StartTimer()
	}
	var buf bytes.Buffer
	dec := json.NewDecoder(&buf)
	var r codeResponse
	for i := 0; i < b.N; i++ {
		buf.Write(codeJSON)
		// hide EOF
		buf.WriteByte('\n')
		buf.WriteByte('\n')
		buf.WriteByte('\n')
		if err := dec.Decode(&r); err != nil {
			b.Fatal("Decode:", err)
		}
	}
	b.SetBytes(int64(len(codeJSON)))
}

func BenchmarkCodeDecoderRgo(b *testing.B) {
	if codeJSON == nil {
		b.StopTimer()
		codeInit()
		b.StartTimer()
	}
	var buf bytes.Buffer
	r := NewReader(&buf)
	var response codeResponse
	for i := 0; i < b.N; i++ {
		buf.Write(codeJSON)
		// hide EOF
		buf.WriteByte('\n')
		buf.WriteByte('\n')
		buf.WriteByte('\n')
		if err := decodeResponse(&response, r); err != nil {
			b.Fatal("decodeResponse:", err)
		}
	}
	b.SetBytes(int64(len(codeJSON)))
}

func BenchmarkUnmarshalStringJson(b *testing.B) {
	data := []byte(`"hello, world"`)
	var s string
	for i := 0; i < b.N; i++ {
		if err := json.Unmarshal(data, &s); err != nil {
			b.Fatal("Unmarshal:", err)
		}
	}
}

func BenchmarkUnmarshalStringRgo(b *testing.B) {
	data := []byte(`"hello, world"`)
	for i := 0; i < b.N; i++ {
		if _, err := NewReader(bytes.NewBuffer(data)).NextString(); err != nil {
			b.Fatal("NextString:", err)
		}
	}
}

func BenchmarkUnmarshalFloat64Json(b *testing.B) {
	data := []byte(`3.14`)
	var f float64
	for i := 0; i < b.N; i++ {
		if err := json.Unmarshal(data, &f); err != nil {
			b.Fatal("Unmarshal:", err)
		}
	}
}

func BenchmarkUnmarshalFloat64Rgo(b *testing.B) {
	data := []byte(`3.14`)
	for i := 0; i < b.N; i++ {
		if _, err := NewReader(bytes.NewBuffer(data)).NextFloat64(); err != nil {
			b.Fatal("NextFloat64:", err)
		}
	}
}

func BenchmarkUnmarshalInt64Json(b *testing.B) {
	data := []byte(`3`)
	var x int64
	for i := 0; i < b.N; i++ {
		if err := json.Unmarshal(data, &x); err != nil {
			b.Fatal("Unmarshal:", err)
		}
	}
}

func BenchmarkUnmarshalInt64Rgo(b *testing.B) {
	data := []byte(`3`)
	for i := 0; i < b.N; i++ {
		if _, err := NewReader(bytes.NewBuffer(data)).NextInt64(); err != nil {
			b.Fatal("NextInt64:", err)
		}
	}
}
