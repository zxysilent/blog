package page

import (
	"fmt"
	"os"
	"testing"
)

func TestRender(t *testing.T) {
	os.MkdirAll("test", 0775)
	os.WriteFile("test/404.html", []byte(Err404().Render()), 0775)
	os.WriteFile("test/csucc.html", []byte(Succ().Render()), 0775)
	os.WriteFile("test/cinfo.html", []byte(Info().Render()), 0775)
	os.WriteFile("test/cwarn.html", []byte(Warn().Render()), 0775)
	os.WriteFile("test/cerror.html", []byte(Error().Render()), 0775)
	os.WriteFile("test/succ.html", []byte(Succ().SetCircle(false).Render()), 0775)
	os.WriteFile("test/info.html", []byte(Info().SetCircle(false).Render()), 0775)
	os.WriteFile("test/warn.html", []byte(Warn().SetCircle(false).Render()), 0775)
	os.WriteFile("test/error.html", []byte(Error().SetCircle(false).Render()), 0775)
}

func BenchmarkRender(t *testing.B) {
	for i := 0; i < t.N; i++ {
		Err404().SetExecJs(fmt.Sprintf(`localStorage.setItem("zblog-token","%s")`, "tokennn")).Render()
	}
}
