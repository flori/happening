package happening

import (
	"strconv"

	"github.com/labstack/echo"
)

type parameters struct {
	Offset  int
	Limit   int
	Start   int
	Query   string
	Filters map[string]string
}

func parseFilters(c echo.Context) map[string]string {
	filters := make(map[string]string)
	for param, values := range c.QueryParams() {
		if len(param) > 2 && param[0:2] == "f:" {
			filters[param[2:]] = values[len(values)-1]
		}
	}
	return filters
}

func parseParameters(c echo.Context) parameters {
	p := parameters{
		Query:   c.QueryParam("q"),
		Filters: parseFilters(c),
	}
	var err error
	if o := c.QueryParam("o"); o == "" {
		p.Offset = 0
	} else {
		if p.Offset, err = strconv.Atoi(o); err != nil {
			p.Offset = 0
		}
	}
	if l := c.QueryParam("l"); l == "" {
		p.Limit = 50
	} else {
		if l == "*" {
			p.Limit = MaxInt
		} else if p.Limit, err = strconv.Atoi(l); err != nil {
			p.Limit = 50
		}
	}
	if s := c.QueryParam("s"); s != "" {
		if p.Start, err = strconv.Atoi(s); err != nil {
			p.Start = 0
		}
	}
	return p
}
