package api


import (
	"bytes"
	"net/http"

	. "gopkg.in/check.v1"

	"github.com/gorilla/mux"
)

type FooStruct struct {
	Foo string `json:"foo" mux:"foo" query:"foo" binding:"required" xml:"foo"`
	Bar int    `json:"bar" binding:"max=3" xml:"bar"`
}

type BindingSuite struct {
}

var _ = Suite(&BindingSuite{})

func (s *BindingSuite) TestBindJSONOk(c *C) {
	req, _ := http.NewRequest("POST", "/foo", bytes.NewBufferString(`{"foo":"value", "bar": 2}`))
	req.Header.Set("Content-Type", "application/json")

	var foo FooStruct

	err := Bind(req, &foo)
	c.Assert(err, IsNil)

	c.Check(foo.Foo, Equals, "value")
	c.Check(foo.Bar, Equals, 2)
}

func (s *BindingSuite) TestBindJSONMalformed(c *C) {
	req, _ := http.NewRequest("POST", "/foo", bytes.NewBufferString(`{"foo":"value"`))
	req.Header.Set("Content-Type", "application/json")

	var foo FooStruct

	err := BindJSON(req, &foo)
	c.Assert(err, ErrorMatches, ".*malformed json: unexpected EOF")
}

func (s *BindingSuite) TestBindNoHeader(c *C) {
	req, _ := http.NewRequest("POST", "/foo", bytes.NewBufferString(`{"foo":"value"}`))

	var foo FooStruct

	err := Bind(req, &foo)
	c.Assert(err, ErrorMatches, ".*Field validation for 'Foo' failed on the 'required' tag")
}

func (s *BindingSuite) TestBindXMLOk(c *C) {
	req, _ := http.NewRequest("POST", "/foo", bytes.NewBufferString(`<root><foo>value</foo><bar>2</bar></root>`))

	req.Header.Set("Content-Type", "application/xml")

	var foo FooStruct

	err := Bind(req, &foo)
	c.Assert(err, IsNil)

	c.Check(foo.Foo, Equals, "value")
	c.Check(foo.Bar, Equals, 2)
}

func (s *BindingSuite) TestBindXMLMalformed(c *C) {
	req, _ := http.NewRequest("POST", "/foo", bytes.NewBufferString(`<root><foo>value</foo><bar>2</bar></root`))

	req.Header.Set("Content-Type", "application/xml")

	var foo FooStruct

	err := BindXML(req, &foo)
	c.Assert(err, ErrorMatches, ".*MalformedXML")
}

func (s *BindingSuite) TestBindValidation1(c *C) {
	req, _ := http.NewRequest("POST", "/foo", bytes.NewBufferString(`{"bar": 1}`))
	req.Header.Set("Content-Type", "application/json")

	var foo FooStruct

	err := Bind(req, &foo)
	c.Assert(err, ErrorMatches, ".*Field validation for 'Foo' failed on the 'required' tag")
}

func (s *BindingSuite) TestBindValidation2(c *C) {
	req, _ := http.NewRequest("POST", "/foo", bytes.NewBufferString(`{"foo": "value", "bar": 42}`))
	req.Header.Set("Content-Type", "application/json")

	var foo FooStruct

	err := Bind(req, &foo)
	c.Assert(err, ErrorMatches, ".*Field validation for 'Bar' failed on the 'max' tag")
}

func (s *BindingSuite) TestBindMuxOk(c *C) {
	c.Skip("no easy way to fool mux vars into request")

	r := mux.NewRouter()
	r.Handle("/item/{foo}", http.Handler(nil))

	req, _ := http.NewRequest("GET", "/item/yahoo", nil)

	var match mux.RouteMatch
	c.Assert(r.Match(req, &match), Equals, true)

	var foo FooStruct

	err := Bind(req, &foo)
	c.Assert(err, IsNil)

	c.Check(foo.Foo, Equals, "yahoo")
	c.Check(foo.Bar, Equals, 0)
}

func (s *BindingSuite) TestBindURLOk(c *C) {
	req, _ := http.NewRequest("GET", "/item/yahoo?foo=super", nil)

	var foo FooStruct

	err := Bind(req, &foo)
	c.Assert(err, IsNil)

	c.Check(foo.Foo, Equals, "super")
	c.Check(foo.Bar, Equals, 0)
}
