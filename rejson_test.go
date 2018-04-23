package rejson

import (
	"reflect"
	"testing"

	"github.com/gomodule/redigo/redis"
)

type TestObject struct {
	Name   string `json:"name"`
	Number int    `json:"number"`
}

func TestJSONSet(t *testing.T) {
	conn, err := redis.Dial("tcp", ":6390")
	if err != nil {
		t.Fatal("Could not connect to redis.")
		return
	}
	defer func() {
		conn.Close()
		conn.Do("FLUSHALL")
	}()

	testObj := TestObject{
		"item#1",
		1,
	}
	type args struct {
		key  string
		path string
		obj  interface{}
		NX   bool
		XX   bool
	}
	tests := []struct {
		name    string
		args    args
		wantRes interface{}
		wantErr bool
	}{
		{
			name: "SimpleString",
			args: args{
				key:  "kstr",
				path: ".",
				obj:  "simplestring",
			},
			wantRes: "OK",
			wantErr: false,
		},
		{
			name: "SimpleInt",
			args: args{
				key:  "kint",
				path: ".",
				obj:  1234,
			},
			wantRes: "OK",
			wantErr: false,
		},
		{
			name: "SimpleStruct",
			args: args{
				key:  "kstruct",
				path: ".",
				obj:  testObj,
			},
			wantRes: "OK",
			wantErr: false,
		},
		{
			name: "SimpleStructFieldOK",
			args: args{
				key:  "kstruct",
				path: "name",
				obj:  "foobar",
			},
			wantRes: "OK",
			wantErr: false,
		},
		// {
		// 	name: "SimpleStringWithNXOK",
		// 	args: args{
		// 		key:  "knxstr",
		// 		path: ".",
		// 		obj:  123,
		// 		NX:   true,
		// 	},
		// 	wantRes: "OK",
		// 	wantErr: false,
		// },
		{
			name: "SimpleStringWithNXNotOK",
			args: args{
				key:  "kstrnx",
				path: ".",
				obj:  "simplestringnx",
				NX:   true,
			},
			wantRes: nil,
			wantErr: false,
		},
		{
			name: "SimpleStringWithXXOK",
			args: args{
				key:  "kstrnx",
				path: ".",
				obj:  "simplestringfoo",
				XX:   true,
			},
			wantRes: "OK",
			wantErr: false,
		},
		{
			name: "SimpleStringWithXXNotOK",
			args: args{
				key:  "kstrxx",
				path: ".",
				obj:  "simplestringfoobar",
				XX:   true,
			},
			wantRes: nil,
			wantErr: false,
		},
		{
			name: "SimpleStringWithXXNX",
			args: args{
				key:  "kstrxxnx",
				path: ".",
				obj:  "simplestringfoobar",
				XX:   true,
				NX:   true,
			},
			wantRes: nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotRes, err := JSONSet(conn, tt.args.key, tt.args.path, tt.args.obj, tt.args.NX, tt.args.XX)
			if (err != nil) != tt.wantErr {
				t.Errorf("JSONSet() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotRes, tt.wantRes) {
				t.Errorf("JSONSet() = %v, want %v", gotRes, tt.wantRes)
			}
		})
	}
}

func TestJSONGet(t *testing.T) {
	conn, err := redis.Dial("tcp", ":6390")
	if err != nil {
		t.Fatal("Could not connect to redis.")
		return
	}
	defer func() {
		conn.Close()
		conn.Do("FLUSHALL")
	}()

	_, err = JSONSet(conn, "kstr", ".", "simplestring", false, false)
	if err != nil {
		return
	}

	_, err = JSONSet(conn, "kint", ".", 123, false, false)
	if err != nil {
		return
	}

	testObj := TestObject{
		Name:   "Item#1",
		Number: 1,
	}

	_, err = JSONSet(conn, "kstruct", ".", testObj, false, false)
	if err != nil {
		return
	}

	type args struct {
		conn redis.Conn
		key  string
		path string
	}
	tests := []struct {
		name    string
		args    args
		wantRes interface{}
		wantErr bool
	}{
		{
			name: "SimpleString",
			args: args{
				conn: conn,
				key:  "kstr",
				path: ".",
			},
			wantRes: []byte("\"simplestring\""),
			wantErr: false,
		},
		{
			name: "SimpleInt",
			args: args{
				conn: conn,
				key:  "kint",
				path: ".",
			},
			wantRes: []byte("123"),
			wantErr: false,
		},
		{
			name: "SimpleStruct",
			args: args{
				conn: conn,
				key:  "kstruct",
				path: ".",
			},
			wantRes: []byte("{\"name\":\"Item#1\",\"number\":1}"),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotRes, err := JSONGet(tt.args.conn, tt.args.key, tt.args.path)
			if (err != nil) != tt.wantErr {
				t.Errorf("JSONGet() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotRes, tt.wantRes) {
				t.Errorf("JSONGet() = %v, want %v", gotRes, tt.wantRes)
			}
		})
	}
}
