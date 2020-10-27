package icon

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/yankghjh/selfhosted_store/cli/pipe"
)

func init() {
	pipe.RegisterInitFunc("icon", InitFunc)
	pipe.RegisterSourceLoader("icon", Loader)
}

// Icon for the app
type Icon struct {
	URL string
}

// InitFunc make dir for icons
func InitFunc(pipe *pipe.Pipe) error {
	os.MkdirAll(pipe.GetDistPath("assets/icons"), os.ModePerm)
	return nil
}

// Loader for app.yml
func Loader(pipe *pipe.Pipe, ctx *pipe.Context) error {
	path := ctx.GetPath("icon.png")

	_, err := os.Stat(path)
	if err != nil {
		if os.IsExist(err) {
			return fmt.Errorf("stat icon %s error: %s", path, err.Error())
		}
		return nil
	}

	distpath := pipe.GetDistPath("assets/icons/" + ctx.Name + ".png")

	err = copyFile(path, distpath)
	if err != nil {
		return fmt.Errorf("copy icon %s error: %s", path, err.Error())
	}

	icon := &Icon{
		URL: pipe.Config.GetString("icon.basepath") + ctx.Name + ".png",
	}

	ctx.Set("icon", icon)

	return nil
}

func copyFile(src, dst string) error {
	input, err := ioutil.ReadFile(src)
	if err != nil {
		return fmt.Errorf("read file %s error: %s", src, err.Error())
	}

	err = ioutil.WriteFile(dst, input, 0644)
	if err != nil {
		return fmt.Errorf("write file %s error: %s", dst, err.Error())
	}

	return nil
}
