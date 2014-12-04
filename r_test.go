package rgo

import (
	"bytes"
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
	buf.Reset()
	w = NewWriter(&buf)
	if err := w.BeginObject(); err != nil {
		t.Errorf("TestWriteObject:BeginObject:err=%s", err.Error())
		return
	}
	if err := w.BeginArray(); err != IllegalState {
		if err == nil {
			t.Errorf("TestWriteObject:BeginArray:err=nil")
		} else {
			t.Errorf("TestWriteObject:BeginArray:err=%s", err.Error())
		}
		return
	}
	if err := w.NullValue(); err != IllegalState {
		if err == nil {
			t.Errorf("TestWriteObject:NullValue:err=nil")
		} else {
			t.Errorf("TestWriteObject:NullValue:err=%s", err.Error())
		}
		return
	}
	if err := w.BoolValue(true); err != IllegalState {
		if err == nil {
			t.Errorf("TestWriteObject:BoolValue:err=nil")
		} else {
			t.Errorf("TestWriteObject:BoolValue:err=%s", err.Error())
		}
		return
	}
	if err := w.IntValue(0); err != IllegalState {
		if err == nil {
			t.Errorf("TestWriteObject:IntValue:err=nil")
		} else {
			t.Errorf("TestWriteObject:IntValue:err=%s", err.Error())
		}
		return
	}
	if err := w.Float64Value(0); err != IllegalState {
		if err == nil {
			t.Errorf("TestWriteObject:Float64Value:err=nil")
		} else {
			t.Errorf("TestWriteObject:Float64Value:err=%s", err.Error())
		}
		return
	}
	if err := w.BeginObject(); err != IllegalState {
		if err == nil {
			t.Errorf("TestWriteObject:BeginObject:err=nil")
		} else {
			t.Errorf("TestWriteObject:BeginObject:err=%s", err.Error())
		}
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
	if err := w.Float64Value(1e-10); err != nil {
		t.Errorf("TestWriteValue:Float64Value:err=%s", err.Error())
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
	if s := buf.String(); s != "[null,true,false,0,1e-10,\"a\\\\\\\"a\",[],{}]" {
		t.Errorf("TestWriteValue:expected=[null,true,false,0,1e-10,\"a\\\\\\\"a\",[],{}],s==%s", s)
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
	if s := buf.String(); s != "{\"null\":null,\"bool\":true,\"int\":-1,\"float\":-1.1,\"string\":\"string\",\"array\":[],\"object\":{}}" {
		t.Errorf("TestWriteValue:expected={\"null\":null,\"bool\":true,\"int\":-1,\"float\":-1.1,\"string\":\"string\",\"array\":[],\"object\":{}},s==%s", s)
		return
	}
}

func TestReadArray(t *testing.T) {
	r := NewReader(bytes.NewBufferString("[]"))
	if token, err := r.Peek(); err != nil {
		t.Errorf("TestReadArray:Peek:err=%s", err.Error())
		return
	} else if token != BEGIN_ARRAY {
		t.Errorf("TestReadArray:Peek:token=%s", token)
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
}

func TestReadObject(t *testing.T) {
	r := NewReader(bytes.NewBufferString("{}"))
	if token, err := r.Peek(); err != nil {
		t.Errorf("TestReadObject:Peek:err=%s", err.Error())
		return
	} else if token != BEGIN_OBJECT {
		t.Errorf("TestReadObject:Peek:token=%s", token)
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
}
