package modules

import "fmt"

type Mongodb struct{}

func (t *Mongodb) RunString(arg ...string) (string, error) {
	cmd := `mongo -u %s -p "%s"  --eval 'db.stats()'`
	return fmt.Sprintf(cmd, arg), nil
}

func (t *Mongodb) Handler(in string) ([]Result, error) {
	return []Result{}, nil
}
