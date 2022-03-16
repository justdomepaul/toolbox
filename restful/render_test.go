package restful

import (
	"github.com/stretchr/testify/suite"
	"html/template"
	"reflect"
	"testing"
)

type RenderSuite struct {
	suite.Suite
}

func (suite *RenderSuite) TestNewRender() {
	suite.Equal("*restful.Render", reflect.TypeOf(NewRender()).String())
}

func (suite *RenderSuite) TestAddMethod() {
	render := NewRender()
	tmpl, err := template.New("tmpl").Parse(`<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <title>{{ .MainTitle }}</title>
</head>
<body>
	<h1>Hello</h1>
</body>
</html>`)
	suite.NoError(err)
	suite.NoError(render.Add("test", tmpl))
}

func (suite *RenderSuite) TestAddMethodNoName() {
	render := NewRender()
	tmpl, err := template.New("tmpl").Parse(`<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <title>{{ .MainTitle }}</title>
</head>
<body>
	<h1>Hello</h1>
</body>
</html>`)
	suite.NoError(err)
	suite.Error(render.Add("", tmpl))
}

func (suite *RenderSuite) TestAddMethodNoTemplate() {
	render := NewRender()
	suite.Error(render.Add("", nil))
}

func (suite *RenderSuite) TestInstanceMethod() {
	render := NewRender()
	tmpl, err := template.New("tmpl").Parse(`<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <title>{{ .MainTitle }}</title>
</head>
<body>
	<h1>Hello</h1>
</body>
</html>`)
	suite.NoError(err)
	suite.NoError(render.Add("test", tmpl))

	suite.Equal("*render.HTML", reflect.TypeOf(render.Instance("test", map[string]interface{}{})).String())
}

func TestRenderSuite(t *testing.T) {
	suite.Run(t, new(RenderSuite))
}
