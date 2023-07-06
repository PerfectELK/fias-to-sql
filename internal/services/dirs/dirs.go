package dirs

import "os"

func InitServiceDirs() error {
	serviceDirs := struct {
		logs    string
		storage string
	}{
		logs:    "/log",
		storage: "/storage",
	}

	p, _ := os.Getwd()

	if _, err := os.Stat(p + serviceDirs.logs); err != nil {
		err = os.MkdirAll(p+serviceDirs.logs, os.ModePerm)
		if err != nil {
			return err
		}
	}

	if _, err := os.Stat(p + serviceDirs.storage); err != nil {
		err = os.MkdirAll(p+serviceDirs.storage, os.ModePerm)
		if err != nil {
			return err
		}
	}

	return nil
}
