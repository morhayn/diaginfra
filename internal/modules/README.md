Modules for check work programm

New Module add new file
```
package modules

type PROGRAMM struct{}

//Return string command for running on remote server 
func (t *PROGRAMM) RunString(arg ...string) (string, error)

//Return array Result 
// Hadler for parse response run command
func (t *PROGRAMM) Handler(in string) ([]Result, error) 

```