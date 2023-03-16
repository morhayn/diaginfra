Modules for check work programm

New Module add new file
```
package modules

type PROGRAMM struct{}

//Return string command for running on remote server
// arg split ':' string from config file 
func (t PROGRAMM) RunString(arg ...string) (string, error)

//function for get logs from servers
func (t PROGRAMM) Logs(count int, arg ...string) (string, error)

//Return array Result 
// Hadler for parse response run command
// 'in' string stdout command
func (t PROGRAMM) Handler(in string) ([]Result, error) 

```