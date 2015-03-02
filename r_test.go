package rgo

import (
	"bytes"
	"io"
	"math"
	"testing"
)

func TestWriteArray(t *testing.T) {
	buf := bytes.Buffer{}
	w := NewWriter(&buf)
	if err := w.BeginArray(); err != nil {
		t.Errorf("TestWriteArray:BeginArray:err=%s", err.Error())
		return
	}
	if err := w.EndArray(); err != nil {
		t.Errorf("TestWriteArray:EndArray:err=%s", err.Error())
		return
	}
	if s := buf.String(); s != "[]" {
		t.Errorf("TestWriteArray:expected=[],s=%s", s)
		return
	}
	buf.Reset()
	w = NewWriter(&buf)
	if err := w.BeginArray(); err != nil {
		t.Errorf("TestWriteArray:BeginArray:err=%s", err.Error())
		return
	}
	for i := 0; i < 2; i++ {
		if err := w.NullValue(); err != nil {
			t.Errorf("TestWriteArray:NullValue:err=%s", err.Error())
			return
		}
	}
	if err := w.EndArray(); err != nil {
		t.Errorf("TestWriteArray:EndArray:err=%s", err.Error())
		return
	}
	if s := buf.String(); s != "[null,null]" {
		t.Errorf("TestWriteArray:expected=[null,null],s=%s", s)
		return
	}
}

func TestWriteObject(t *testing.T) {
	buf := bytes.Buffer{}
	w := NewWriter(&buf)
	if err := w.BeginObject(); err != nil {
		t.Errorf("TestWriteObject:BeginObject:err=%s", err.Error())
		return
	}
	if err := w.EndObject(); err != nil {
		t.Errorf("TestWriteObject:EndObject:err=%s", err.Error())
		return
	}
	if s := buf.String(); s != "{}" {
		t.Errorf("TestWriteObject:expected={},s==%s", s)
		return
	}
}

func TestWriteValue(t *testing.T) {
	buf := bytes.Buffer{}
	w := NewWriter(&buf)
	if err := w.BeginArray(); err != nil {
		t.Errorf("TestWriteValue:BeginArray:err=%s", err.Error())
		return
	}
	if err := w.NullValue(); err != nil {
		t.Errorf("TestWriteValue:NullValue:err=%s", err.Error())
		return
	}
	if err := w.BoolValue(true); err != nil {
		t.Errorf("TestWriteValue:BoolValue:err=%s", err.Error())
		return
	}
	if err := w.BoolValue(false); err != nil {
		t.Errorf("TestWriteValue:BoolValue:err=%s", err.Error())
		return
	}
	if err := w.IntValue(0); err != nil {
		t.Errorf("TestWriteValue:IntValue:err=%s", err.Error())
		return
	}
	if err := w.Int8Value(-128); err != nil {
		t.Errorf("TestWriteValue:Int8Value:err=%s", err.Error())
		return
	}
	if err := w.Int16Value(-32768); err != nil {
		t.Errorf("TestWriteValue:Int16Value:err=%s", err.Error())
		return
	}
	if err := w.Int32Value(-32769); err != nil {
		t.Errorf("TestWriteValue:Int32Value:err=%s", err.Error())
		return
	}
	if err := w.Int64Value(-8000000000); err != nil {
		t.Errorf("TestWriteValue:Int64Value:err=%s", err.Error())
		return
	}
	if err := w.UintValue(0); err != nil {
		t.Errorf("TestWriteValue:UintValue:err=%s", err.Error())
		return
	}
	if err := w.Uint8Value(255); err != nil {
		t.Errorf("TestWriteValue:Uint8Value:err=%s", err.Error())
		return
	}
	if err := w.Uint16Value(65535); err != nil {
		t.Errorf("TestWriteValue:Uint16Value:err=%s", err.Error())
		return
	}
	if err := w.Uint32Value(65536); err != nil {
		t.Errorf("TestWriteValue:Uint32Value:err=%s", err.Error())
		return
	}
	if err := w.Uint64Value(8000000000); err != nil {
		t.Errorf("TestWriteValue:Uint64Value:err=%s", err.Error())
		return
	}
	if err := w.Float32Value(0.5); err != nil {
		t.Errorf("TestWriteValue:Float32Value:err=%s", err.Error())
		return
	}
	if err := w.Float64Value(1e-10); err != nil {
		t.Errorf("TestWriteValue:Float64Value:err=%s", err.Error())
		return
	}
	if err := w.Float64Value(math.Inf(-1)); err != IllegalArgument {
		if err == nil {
			t.Errorf("TestWriteValue:Float64Value:err=nil")
		} else {
			t.Errorf("TestWriteValue:Float64Value:err=%s", err.Error())
		}
		return
	}
	if err := w.Float64Value(math.NaN()); err != IllegalArgument {
		if err == nil {
			t.Errorf("TestWriteValue:Float64Value:err=nil")
		} else {
			t.Errorf("TestWriteValue:Float64Value:err=%s", err.Error())
		}
		return
	}
	if err := w.StringValue("a\\\"a"); err != nil {
		t.Errorf("TestWriteValue:StringValue:err=%s", err.Error())
		return
	}
	if err := w.BeginArray(); err != nil {
		t.Errorf("TestWriteValue:BeginArray:err=%s", err.Error())
		return
	}
	if err := w.EndArray(); err != nil {
		t.Errorf("TestWriteValue:EndArray:err=%s", err.Error())
		return
	}
	if err := w.BeginObject(); err != nil {
		t.Errorf("TestWriteValue:BeginObject:err=%s", err.Error())
		return
	}
	if err := w.EndObject(); err != nil {
		t.Errorf("TestWriteValue:EndObject:err=%s", err.Error())
		return
	}
	if err := w.EndArray(); err != nil {
		t.Errorf("TestWriteValue:EndArray:err=%s", err.Error())
		return
	}
	if s := buf.String(); s != "[null,true,false,0,-128,-32768,-32769,-8000000000,0,255,65535,65536,8000000000,0.5,1e-10,\"a\\\\\\\"a\",[],{}]" {
		t.Errorf("TestWriteValue:expected=[null,true,false,0,-128,-32768,-32769,-8000000000,0,255,65535,65536,8000000000,0.5,1e-10,\"a\\\\\\\"a\",[],{}],s==%s", s)
		return
	}
	buf.Reset()
	w = NewWriter(&buf)
	if err := w.BeginObject(); err != nil {
		t.Errorf("TestWriteValue:BeginObject:err=%s", err.Error())
		return
	}
	if err := w.Name("null"); err != nil {
		t.Errorf("TestWriteValue:Name:err=%s", err.Error())
		return
	}
	if err := w.Value(nil); err != nil {
		t.Errorf("TestWriteValue:Value:err=%s", err.Error())
		return
	}
	if err := w.Name("bool"); err != nil {
		t.Errorf("TestWriteValue:Name:err=%s", err.Error())
		return
	}
	if err := w.Value(true); err != nil {
		t.Errorf("TestWriteValue:Value:err=%s", err.Error())
		return
	}
	if err := w.Name("int"); err != nil {
		t.Errorf("TestWriteValue:Name:err=%s", err.Error())
		return
	}
	if err := w.Value(-1); err != nil {
		t.Errorf("TestWriteValue:Value:err=%s", err.Error())
		return
	}
	if err := w.Name("float"); err != nil {
		t.Errorf("TestWriteValue:Name:err=%s", err.Error())
		return
	}
	if err := w.Value(-1.1); err != nil {
		t.Errorf("TestWriteValue:Value:err=%s", err.Error())
		return
	}
	if err := w.Name("string"); err != nil {
		t.Errorf("TestWriteValue:Name:err=%s", err.Error())
		return
	}
	if err := w.Value("string"); err != nil {
		t.Errorf("TestWriteValue:Value:err=%s", err.Error())
		return
	}
	if err := w.Name("array"); err != nil {
		t.Errorf("TestWriteValue:Name:err=%s", err.Error())
		return
	}
	if err := w.BeginArray(); err != nil {
		t.Errorf("TestWriteValue:BeginArray:err=%s", err.Error())
		return
	}
	if err := w.Value(int(-1)); err != nil {
		t.Errorf("TestWriteValue:Value:err=%s", err.Error())
		return
	}
	if err := w.Value(int8(-8)); err != nil {
		t.Errorf("TestWriteValue:Value:err=%s", err.Error())
		return
	}
	if err := w.Value(int16(-16)); err != nil {
		t.Errorf("TestWriteValue:Value:err=%s", err.Error())
		return
	}
	if err := w.Value(int32(-32)); err != nil {
		t.Errorf("TestWriteValue:Value:err=%s", err.Error())
		return
	}
	if err := w.Value(int64(-64)); err != nil {
		t.Errorf("TestWriteValue:Value:err=%s", err.Error())
		return
	}
	if err := w.Value(uint(1)); err != nil {
		t.Errorf("TestWriteValue:Value:err=%s", err.Error())
		return
	}
	if err := w.Value(uint8(8)); err != nil {
		t.Errorf("TestWriteValue:Value:err=%s", err.Error())
		return
	}
	if err := w.Value(uint16(16)); err != nil {
		t.Errorf("TestWriteValue:Value:err=%s", err.Error())
		return
	}
	if err := w.Value(uint32(32)); err != nil {
		t.Errorf("TestWriteValue:Value:err=%s", err.Error())
		return
	}
	if err := w.Value(uint64(64)); err != nil {
		t.Errorf("TestWriteValue:Value:err=%s", err.Error())
		return
	}
	if err := w.Value(float32(0.32)); err != nil {
		t.Errorf("TestWriteValue:Value:err=%s", err.Error())
		return
	}
	if err := w.Value(float64(0.64)); err != nil {
		t.Errorf("TestWriteValue:Value:err=%s", err.Error())
		return
	}
	if err := w.Value([]int{1}); err != IllegalArgument {
		if err == nil {
			t.Errorf("TestWriteValue:Name:err=nil")
		} else {
			t.Errorf("TestWriteValue:Name:err=%s", err.Error())
		}
		return
	}
	if err := w.EndArray(); err != nil {
		t.Errorf("TestWriteValue:EndArray:err=%s", err.Error())
		return
	}
	if err := w.Name("object"); err != nil {
		t.Errorf("TestWriteValue:Name:err=%s", err.Error())
		return
	}
	if err := w.BeginObject(); err != nil {
		t.Errorf("TestWriteValue:BeginObject:err=%s", err.Error())
		return
	}
	if err := w.EndObject(); err != nil {
		t.Errorf("TestWriteValue:EndObject:err=%s", err.Error())
		return
	}
	if err := w.EndObject(); err != nil {
		t.Errorf("TestWriteValue:EndObject:err=%s", err.Error())
		return
	}
	if s := buf.String(); s != "{\"null\":null,\"bool\":true,\"int\":-1,\"float\":-1.1,\"string\":\"string\",\"array\":[-1,-8,-16,-32,-64,1,8,16,32,64,0.32,0.64],\"object\":{}}" {
		t.Errorf("TestWriteValue:expected={\"null\":null,\"bool\":true,\"int\":-1,\"float\":-1.1,\"string\":\"string\",\"array\":[-1,-8,-16,-32,-64,1,8,16,32,64,0.32,0.64],\"object\":{}},s==%s", s)
		return
	}
	buf.Reset()
	w = NewWriter(&buf)
	if err := w.StringValue(" \x00 \x08 \x09 \x0a \x0c \x0d \x0f \x10 \x1f \\ \""); err != nil {
		t.Errorf("TestWriteValue:EndObject:err=%s", err.Error())
		return
	}
	if s := buf.String(); s != `" \u0000 \b \t \n \f \r \u000f \u0010 \u001f \\ \""` {
		t.Errorf("TestWriteValue:s==%s", s)
		return
	}
}

func TestReadArray(t *testing.T) {
	r := NewReader(bytes.NewBufferString("[]"))
	if token, err := r.Peek(); err != nil {
		t.Errorf("TestReadArray:Peek:err=%s", err.Error())
		return
	} else if token != BEGIN_ARRAY {
		t.Errorf("TestReadArray:Peek:token=%d", token)
		return
	}
	if err := r.BeginArray(); err != nil {
		t.Errorf("TestReadArray:BeginArray:err=%s", err.Error())
		return
	}
	if hasNext, err := r.HasNext(); err != nil {
		t.Errorf("TestReadArray:HasNext:err=%s", err.Error())
		return
	} else if hasNext {
		t.Errorf("TestReadArray:HasNext=true")
		return
	}
	if err := r.EndArray(); err != nil {
		t.Errorf("TestReadArray:EndArray:err=%s", err.Error())
		return
	}
	r = NewReader(bytes.NewBufferString(" [ null , true , -2, \"string\" ] "))
	if err := r.BeginArray(); err != nil {
		t.Errorf("TestReadArray:BeginArray:err=%s", err.Error())
		return
	}
	if hasNext, err := r.HasNext(); err != nil {
		t.Errorf("TestReadArray:HasNext:err=%s", err.Error())
		return
	} else if !hasNext {
		t.Errorf("TestReadArray:HasNext=false")
		return
	}
	if err := r.NextNull(); err != nil {
		t.Errorf("TestReadArray:NextNull:err=%s", err.Error())
		return
	}
	if hasNext, err := r.HasNext(); err != nil {
		t.Errorf("TestReadArray:HasNext:err=%s", err.Error())
		return
	} else if !hasNext {
		t.Errorf("TestReadArray:HasNext=false")
		return
	}
	if value, err := r.NextBoolean(); err != nil {
		t.Errorf("TestReadArray:NextBoolean:err=%s", err.Error())
		return
	} else if !value {
		t.Errorf("TestReadArray:NextBoolean=false")
		return
	}
	if value, err := r.NextInt(); err != nil {
		t.Errorf("TestReadArray:NextInt:err=%s", err.Error())
		return
	} else if value != -2 {
		t.Errorf("TestReadArray:NextInt:value=%d", value)
		return
	}
	if value, err := r.NextString(); err != nil {
		t.Errorf("TestReadArray:NextString:err=%s", err.Error())
		return
	} else if value != "string" {
		t.Errorf("TestReadArray:NextString:value=%s", value)
		return
	}
	if err := r.EndArray(); err != nil {
		t.Errorf("TestReadArray:EndArray:err=%s", err.Error())
		return
	}
	r = NewReader(bytes.NewBufferString(" ["))
	if err := r.BeginArray(); err != io.ErrUnexpectedEOF {
		if err == nil {
			t.Errorf("TestReadArray:BeginArray:err=nil")
		} else {
			t.Errorf("TestReadArray:BeginArray:err=%s", err.Error())
		}
		return
	}
	r = NewReader(bytes.NewBufferString("{}"))
	if err := r.BeginArray(); err != IllegalState {
		if err == nil {
			t.Errorf("TestReadArray:BeginArray:err=nil")
		} else {
			t.Errorf("TestReadArray:BeginArray:err=%s", err.Error())
		}
		return
	}
	r = NewReader(bytes.NewBufferString("[null]"))
	if err := r.BeginArray(); err != nil {
		t.Errorf("TestReadArray:BeginArray:err=%s", err.Error())
		return
	}
	if err := r.EndArray(); err != IllegalState {
		if err == nil {
			t.Errorf("TestReadArray:EndArray:err=nil")
		} else {
			t.Errorf("TestReadArray:EndArray:err=%s", err.Error())
		}
		return
	}
}

func TestReadObject(t *testing.T) {
	r := NewReader(bytes.NewBufferString("{}"))
	if token, err := r.Peek(); err != nil {
		t.Errorf("TestReadObject:Peek:err=%s", err.Error())
		return
	} else if token != BEGIN_OBJECT {
		t.Errorf("TestReadObject:Peek:token=%d", token)
		return
	}
	if err := r.BeginObject(); err != nil {
		t.Errorf("TestReadObject:BeginObject:err=%s", err.Error())
		return
	}
	if hasNext, err := r.HasNext(); err != nil {
		t.Errorf("TestReadObject:HasNext:err=%s", err.Error())
		return
	} else if hasNext {
		t.Errorf("TestReadObject:HasNext=true")
		return
	}
	if err := r.EndObject(); err != nil {
		t.Errorf("TestReadObject:EndObject:err=%s", err.Error())
		return
	}
	r = NewReader(bytes.NewBufferString(`{"null":null,"boolean":true,"int":0,"float":-0.5,"string":"value"}`))
	if err := r.BeginObject(); err != nil {
		t.Errorf("TestReadObject:BeginObject:err=%s", err.Error())
		return
	}
	if hasNext, err := r.HasNext(); err != nil {
		t.Errorf("TestReadObject:HasNext:err=%s", err.Error())
		return
	} else if !hasNext {
		t.Errorf("TestReadObject:HasNext=false")
		return
	}
	if name, err := r.NextName(); err != nil {
		t.Errorf("TestReadObject:NextName:err=%s", err.Error())
		return
	} else if name != "null" {
		t.Errorf("TestReadObject:NextName:name=%s", name)
		return
	}
	if err := r.NextNull(); err != nil {
		t.Errorf("TestReadObject:NextNull:err=%s", err.Error())
		return
	}
	if hasNext, err := r.HasNext(); err != nil {
		t.Errorf("TestReadObject:HasNext:err=%s", err.Error())
		return
	} else if !hasNext {
		t.Errorf("TestReadObject:HasNext=false")
		return
	}
	if name, err := r.NextName(); err != nil {
		t.Errorf("TestReadObject:NextName:err=%s", err.Error())
		return
	} else if name != "boolean" {
		t.Errorf("TestReadObject:NextName:name=%s", name)
		return
	}
	if value, err := r.NextBoolean(); err != nil {
		t.Errorf("TestReadObject:NextBoolean:err=%s", err.Error())
		return
	} else if !value {
		t.Errorf("TestReadObject:NextBoolean=false")
		return
	}
	if name, err := r.NextName(); err != nil {
		t.Errorf("TestReadObject:NextName:err=%s", err.Error())
		return
	} else if name != "int" {
		t.Errorf("TestReadObject:NextName:name=%s", name)
		return
	}
	if value, err := r.NextInt(); err != nil {
		t.Errorf("TestReadObject:NextInt:err=%s", err.Error())
		return
	} else if value != 0 {
		t.Errorf("TestReadObject:NextInt=%d", value)
		return
	}
	if name, err := r.NextName(); err != nil {
		t.Errorf("TestReadObject:NextName:err=%s", err.Error())
		return
	} else if name != "float" {
		t.Errorf("TestReadObject:NextName:name=%s", name)
		return
	}
	if value, err := r.NextFloat64(); err != nil {
		t.Errorf("TestReadObject:NextFloat64:err=%s", err.Error())
		return
	} else if value != -0.5 {
		t.Errorf("TestReadObject:NextFloat64=%g", value)
		return
	}
	if name, err := r.NextName(); err != nil {
		t.Errorf("TestReadObject:NextName:err=%s", err.Error())
		return
	} else if name != "string" {
		t.Errorf("TestReadObject:NextName:name=%s", name)
		return
	}
	if value, err := r.NextString(); err != nil {
		t.Errorf("TestReadObject:NextString:err=%s", err.Error())
		return
	} else if value != "value" {
		t.Errorf("TestReadObject:NextString=%s", value)
		return
	}
	if err := r.EndObject(); err != nil {
		t.Errorf("TestReadObject:EndObject:err=%s", err.Error())
		return
	}
	r = NewReader(bytes.NewBufferString("[]"))
	if err := r.BeginObject(); err != IllegalState {
		if err == nil {
			t.Errorf("TestReadObject:BeginObject:err=nil")
		} else {
			t.Errorf("TestReadObject:BeginObject:err=%s", err.Error())
		}
		return
	}
	r = NewReader(bytes.NewBufferString("{\"a\":null}"))
	if err := r.BeginObject(); err != nil {
		t.Errorf("TestReadObject:BeginObject:err=%s", err.Error())
		return
	}
	if err := r.EndObject(); err != IllegalState {
		if err == nil {
			t.Errorf("TestReadObject:EndObject:err=nil")
		} else {
			t.Errorf("TestReadObject:EndObject:err=%s", err.Error())
		}
		return
	}
}

func TestReadNested(t *testing.T) {
	r := NewReader(bytes.NewBufferString(`{"array":[{}]}`))
	if err := r.BeginObject(); err != nil {
		t.Errorf("TestReadNested:BeginObject:err=%s", err.Error())
		return
	}
	if name, err := r.NextName(); err != nil {
		t.Errorf("TestReadNested:NextName:err=%s", err.Error())
		return
	} else if name != "array" {
		t.Errorf("TestReadNested:NextName:name=%s", name)
		return
	}
	if err := r.BeginArray(); err != nil {
		t.Errorf("TestReadNested:BeginArray:err=%s", err.Error())
		return
	}
	if err := r.BeginObject(); err != nil {
		t.Errorf("TestReadNested:BeginObject:err=%s", err.Error())
		return
	}
	if err := r.EndObject(); err != nil {
		t.Errorf("TestReadNested:EndObject:err=%s", err.Error())
		return
	}
	if err := r.EndArray(); err != nil {
		t.Errorf("TestReadNested:EndArray:err=%s", err.Error())
		return
	}
	if err := r.EndObject(); err != nil {
		t.Errorf("TestReadNested:EndObject:err=%s", err.Error())
		return
	}
}

func TestSkipValue(t *testing.T) {
	r := NewReader(bytes.NewBufferString(`[null,null,false,null,1,null,"string",null,{"array":[{}]},null,[[],[]],null]`))
	if err := r.BeginArray(); err != nil {
		t.Errorf("TestSkipValue:BeginArray:err=%s", err.Error())
		return
	}
	if err := r.SkipValue(); err != nil {
		t.Errorf("TestReadObject:SkipValue:err=%s", err.Error())
		return
	}
	if err := r.NextNull(); err != nil {
		t.Errorf("TestReadObject:NextNull:err=%s", err.Error())
		return
	}
	if err := r.SkipValue(); err != nil {
		t.Errorf("TestReadObject:SkipValue:err=%s", err.Error())
		return
	}
	if err := r.NextNull(); err != nil {
		t.Errorf("TestReadObject:NextNull:err=%s", err.Error())
		return
	}
	if err := r.SkipValue(); err != nil {
		t.Errorf("TestReadObject:SkipValue:err=%s", err.Error())
		return
	}
	if err := r.NextNull(); err != nil {
		t.Errorf("TestReadObject:NextNull:err=%s", err.Error())
		return
	}
	if err := r.SkipValue(); err != nil {
		t.Errorf("TestReadObject:SkipValue:err=%s", err.Error())
		return
	}
	if err := r.NextNull(); err != nil {
		t.Errorf("TestReadObject:NextNull:err=%s", err.Error())
		return
	}
	if err := r.SkipValue(); err != nil {
		t.Errorf("TestReadObject:SkipValue:err=%s", err.Error())
		return
	}
	if err := r.NextNull(); err != nil {
		t.Errorf("TestReadObject:NextNull:err=%s", err.Error())
		return
	}
	if err := r.SkipValue(); err != nil {
		t.Errorf("TestReadObject:SkipValue:err=%s", err.Error())
		return
	}
	if err := r.NextNull(); err != nil {
		t.Errorf("TestReadObject:NextNull:err=%s", err.Error())
		return
	}
	if err := r.EndArray(); err != nil {
		t.Errorf("TestSkipValue:EndArray:err=%s", err.Error())
		return
	}
	r = NewReader(bytes.NewBufferString(`{"name":"value"}`))
	if err := r.BeginObject(); err != nil {
		t.Errorf("TestSkipValue:BeginObject:err=%s", err.Error())
		return
	}
	if err := r.SkipValue(); err != IllegalState {
		if err == nil {
			t.Errorf("TestReadObject:SkipValue:err=nil")
		} else {
			t.Errorf("TestReadObject:SkipValue:err=%s", err.Error())
		}
		return
	}
	r = NewReader(bytes.NewBufferString("{}"))
	if err := r.BeginObject(); err != nil {
		t.Errorf("TestSkipValue:BeginObject:err=%s", err.Error())
		return
	}
	if err := r.SkipValue(); err != InvalidInput {
		if err == nil {
			t.Errorf("TestReadObject:SkipValue:err=nil")
		} else {
			t.Errorf("TestReadObject:SkipValue:err=%s", err.Error())
		}
		return
	}
}

func TestUTF16(t *testing.T) {
	r := NewReader(bytes.NewBufferString(`"\uD834\uDD1E"`))
	if value, err := r.NextString(); err != nil {
		t.Errorf("TestUTF16:NextString:err=%s", err.Error())
		return
	} else if value != "\U0001D11E" {
		t.Errorf("TestUTF16:NextString=%s", value)
	}
	r = NewReader(bytes.NewBufferString(`"\uDD1E\uD834"`))
	if _, err := r.NextString(); err != InvalidInput {
		if err == nil {
			t.Errorf("TestUTF16:NextString:err=nil")
		} else {
			t.Errorf("TestUTF16:NextString:err=%s", err.Error())
		}
		return
	}
	r = NewReader(bytes.NewBufferString(`"\uD834"`))
	if _, err := r.NextString(); err != InvalidInput {
		if err == nil {
			t.Errorf("TestUTF16:NextString:err=nil")
		} else {
			t.Errorf("TestUTF16:NextString:err=%s", err.Error())
		}
		return
	}
	r = NewReader(bytes.NewBufferString(`"\uD834\"`))
	if _, err := r.NextString(); err != InvalidInput {
		if err == nil {
			t.Errorf("TestUTF16:NextString:err=nil")
		} else {
			t.Errorf("TestUTF16:NextString:err=%s", err.Error())
		}
		return
	}
	r = NewReader(bytes.NewBufferString(`"\uD834\uD834"`))
	if _, err := r.NextString(); err != InvalidInput {
		if err == nil {
			t.Errorf("TestUTF16:NextString:err=nil")
		} else {
			t.Errorf("TestUTF16:NextString:err=%s", err.Error())
		}
		return
	}
	r = NewReader(bytes.NewBufferString(`"\uD83x\uDD1E"`))
	if _, err := r.NextString(); err == nil {
		t.Errorf("TestUTF16:NextString:err=nil")
		return
	}
	r = NewReader(bytes.NewBufferString(`"\uD834\uDD1x"`))
	if _, err := r.NextString(); err == nil {
		t.Errorf("TestUTF16:NextString:err=nil")
		return
	}
	r = NewReader(bytes.NewBufferString("\"\x08\""))
	if _, err := r.NextString(); err != InvalidInput {
		if err == nil {
			t.Errorf("TestUTF16:NextString:err=nil")
		} else {
			t.Errorf("TestUTF16:NextString:err=%s", err.Error())
		}
		return
	}
	r = NewReader(bytes.NewBufferString(`"""`))
	if _, err := r.NextString(); err != InvalidInput {
		if err == nil {
			t.Errorf("TestUTF16:NextString:err=nil")
		} else {
			t.Errorf("TestUTF16:NextString:err=%s", err.Error())
		}
		return
	}
}

func TestNull(t *testing.T) {
	r := NewReader(bytes.NewBufferString("null"))
	if err := r.NextNull(); err != nil {
		t.Errorf("TestNull:NextNull:err=%s", err.Error())
		return
	}
	if err := r.NextNull(); err != IllegalState {
		if err == nil {
			t.Errorf("TestNull:NextNull:err=nil")
		} else {
			t.Errorf("TestNull:NextNull:err=%s", err.Error())
		}
		return
	}
	r = NewReader(bytes.NewBufferString("NULL"))
	if err := r.NextNull(); err != InvalidInput {
		if err == nil {
			t.Errorf("TestNull:NextNull:err=nil")
		} else {
			t.Errorf("TestNull:NextNull:err=%s", err.Error())
		}
		return
	}
	r = NewReader(bytes.NewBufferString("nulL"))
	if err := r.NextNull(); err != InvalidInput {
		if err == nil {
			t.Errorf("TestNull:NextNull:err=nil")
		} else {
			t.Errorf("TestNull:NextNull:err=%s", err.Error())
		}
		return
	}
	r = NewReader(bytes.NewBufferString("nullx"))
	if err := r.NextNull(); err != InvalidInput {
		if err == nil {
			t.Errorf("TestNull:NextNull:err=nil")
		} else {
			t.Errorf("TestNull:NextNull:err=%s", err.Error())
		}
		return
	}
}

func TestEOF(t *testing.T) {
	r := NewReader(bytes.NewBufferString(""))
	if token, err := r.Peek(); err != nil {
		t.Errorf("TestEOF:Peek:err=%s", err.Error())
		return
	} else if token != END_DOCUMENT {
		t.Errorf("TestEOF:Peek:token=%d", token)
		return
	}
	r = NewReader(bytes.NewBufferString("\""))
	if _, err := r.NextString(); err != io.ErrUnexpectedEOF {
		if err == nil {
			t.Errorf("TestEOF:NextString:err=nil")
		} else {
			t.Errorf("TestEOF:NextString:err=%s", err.Error())
		}
		return
	}
	r = NewReader(bytes.NewBufferString("\"\\"))
	if _, err := r.NextString(); err != io.ErrUnexpectedEOF {
		if err == nil {
			t.Errorf("TestEOF:NextString:err=nil")
		} else {
			t.Errorf("TestEOF:NextString:err=%s", err.Error())
		}
		return
	}
}

func TestReadBoolean(t *testing.T) {
	r := NewReader(bytes.NewBufferString("[true,false]"))
	if err := r.BeginArray(); err != nil {
		t.Errorf("TestReadBoolean:BeginArray:err=%s", err.Error())
		return
	}
	if value, err := r.NextBoolean(); err != nil {
		t.Errorf("TestReadBoolean:NextBoolean:err=%s", err.Error())
		return
	} else if !value {
		t.Errorf("TestReadBoolean:NextBoolean=false")
		return
	}
	if value, err := r.NextBoolean(); err != nil {
		t.Errorf("TestReadBoolean:NextBoolean:err=%s", err.Error())
		return
	} else if value {
		t.Errorf("TestReadBoolean:NextBoolean=true")
		return
	}
	if err := r.EndArray(); err != nil {
		t.Errorf("TestReadBoolean:EndArray:err=%s", err.Error())
		return
	}
	r = NewReader(bytes.NewBufferString("[true,false]"))
	if _, err := r.NextBoolean(); err != IllegalState {
		if err == nil {
			t.Errorf("TestReadBoolean:NextBoolean:err=nil")
		} else {
			t.Errorf("TestReadBoolean:NextBoolean:err=%s", err.Error())
		}
		return
	}
}

func TestReadNumber(t *testing.T) {
	r := NewReader(bytes.NewBufferString("[0,1,-1,0.5,0.25]"))
	if err := r.BeginArray(); err != nil {
		t.Errorf("TestReadNumber:BeginArray:err=%s", err.Error())
		return
	}
	if value, err := r.NextInt(); err != nil {
		t.Errorf("TestReadNumber:NextInt:err=%s", err.Error())
		return
	} else if value != 0 {
		t.Errorf("TestReadNumber:NextInt=%d", value)
		return
	}
	if value, err := r.NextInt(); err != nil {
		t.Errorf("TestReadNumber:NextInt:err=%s", err.Error())
		return
	} else if value != 1 {
		t.Errorf("TestReadNumber:NextInt=%d", value)
		return
	}
	if value, err := r.NextInt64(); err != nil {
		t.Errorf("TestReadNumber:NextInt64:err=%s", err.Error())
		return
	} else if value != -1 {
		t.Errorf("TestReadNumber:NextInt64=%d", value)
		return
	}
	if value, err := r.NextFloat32(); err != nil {
		t.Errorf("TestReadNumber:NextFloat32:err=%s", err.Error())
		return
	} else if value != 0.5 {
		t.Errorf("TestReadNumber:NextFloat32=%g", value)
		return
	}
	if value, err := r.NextFloat64(); err != nil {
		t.Errorf("TestReadNumber:NextFloat64:err=%s", err.Error())
		return
	} else if value != 0.25 {
		t.Errorf("TestReadNumber:NextFloat64=%g", value)
		return
	}
	if err := r.EndArray(); err != nil {
		t.Errorf("TestReadNumber:EndArray:err=%s", err.Error())
		return
	}
	r = NewReader(bytes.NewBufferString("-1"))
	if value, err := r.NextInt(); err != nil {
		t.Errorf("TestReadNumber:NextInt:err=%s", err.Error())
		return
	} else if value != -1 {
		t.Errorf("TestReadNumber:NextInt=%d", value)
		return
	}
	r = NewReader(bytes.NewBufferString("-"))
	if _, err := r.NextInt(); err != io.ErrUnexpectedEOF {
		if err == nil {
			t.Errorf("TestReadNumber:NextInt:err=nil")
		} else {
			t.Errorf("TestReadNumber:NextInt:err=%s", err.Error())
		}
		return
	}
	r = NewReader(bytes.NewBufferString("0.."))
	if _, err := r.NextFloat64(); err != InvalidInput {
		if err == nil {
			t.Errorf("TestReadNumber:NextFloat64:err=nil")
		} else {
			t.Errorf("TestReadNumber:NextFloat64:err=%s", err.Error())
		}
		return
	}
	r = NewReader(bytes.NewBufferString("1ee"))
	if _, err := r.NextFloat64(); err != InvalidInput {
		if err == nil {
			t.Errorf("TestReadNumber:NextFloat64:err=nil")
		} else {
			t.Errorf("TestReadNumber:NextFloat64:err=%s", err.Error())
		}
		return
	}
	r = NewReader(bytes.NewBufferString("1e++"))
	if _, err := r.NextFloat64(); err != InvalidInput {
		if err == nil {
			t.Errorf("TestReadNumber:NextFloat64:err=nil")
		} else {
			t.Errorf("TestReadNumber:NextFloat64:err=%s", err.Error())
		}
		return
	}
	r = NewReader(bytes.NewBufferString("1+"))
	if _, err := r.NextFloat64(); err != InvalidInput {
		if err == nil {
			t.Errorf("TestReadNumber:NextFloat64:err=nil")
		} else {
			t.Errorf("TestReadNumber:NextFloat64:err=%s", err.Error())
		}
		return
	}
	r = NewReader(bytes.NewBufferString("1e+1"))
	if value, err := r.NextFloat64(); err != nil {
		t.Errorf("TestReadNumber:NextFloat64:err=%s", err.Error())
		return
	} else if value != 10 {
		t.Errorf("TestReadNumber:NextFloat64=%g", value)
		return
	}
	r = NewReader(bytes.NewBufferString("001"))
	if _, err := r.NextInt(); err != InvalidInput {
		if err == nil {
			t.Errorf("TestReadNumber:NextInt:err=nil")
		} else {
			t.Errorf("TestReadNumber:NextInt:err=%s", err.Error())
		}
		return
	}
	r = NewReader(bytes.NewBufferString("1024"))
	if value, err := r.NextInt(); err != nil {
		t.Errorf("TestReadNumber:NextInt:err=%s", err.Error())
		return
	} else if value != 1024 {
		t.Errorf("TestReadNumber:NextInt=%d", value)
		return
	}
	r = NewReader(bytes.NewBufferString("-,"))
	if _, err := r.NextInt(); err != InvalidInput {
		if err == nil {
			t.Errorf("TestReadNumber:NextInt:err=nil")
		} else {
			t.Errorf("TestReadNumber:NextInt:err=%s", err.Error())
		}
		return
	}
	r = NewReader(bytes.NewBufferString("null"))
	if _, err := r.NextInt(); err != IllegalState {
		if err == nil {
			t.Errorf("TestReadNumber:NextInt:err=nil")
		} else {
			t.Errorf("TestReadNumber:NextInt:err=%s", err.Error())
		}
		return
	}
	r = NewReader(bytes.NewBufferString("null"))
	if _, err := r.NextFloat64(); err != IllegalState {
		if err == nil {
			t.Errorf("TestReadNumber:NextFloat64:err=nil")
		} else {
			t.Errorf("TestReadNumber:NextFloat64:err=%s", err.Error())
		}
		return
	}
}

func TestReadString(t *testing.T) {
	r := NewReader(bytes.NewBufferString(`"\\\"\b\t\f\r\n\u0001"`))
	if value, err := r.NextString(); err != nil {
		t.Errorf("TestReadString:NextString:err=%s", err.Error())
		return
	} else if value != "\\\"\b\t\f\r\n\x01" {
		t.Errorf("TestReadString:NextString=%s", value)
		return
	}
	r = NewReader(bytes.NewBufferString(`"\\\"\b\t\f\r\n\u0001"`))
	if _, err := r.NextName(); err != IllegalState {
		if err == nil {
			t.Errorf("TestReadString:NextName:err=nil")
		} else {
			t.Errorf("TestReadString:NextName:err=%s", err.Error())
		}
		return
	}
	r = NewReader(bytes.NewBufferString("null"))
	if _, err := r.NextString(); err != IllegalState {
		if err == nil {
			t.Errorf("TestReadString:NextString:err=nil")
		} else {
			t.Errorf("TestReadString:NextString:err=%s", err.Error())
		}
		return
	}
}
