package modules

type Module interface {
	RunString(string) (string, error)
	Handler(string) error
}
type Results struct {
	Res []Result
}
type Result struct {
	Service string `json:"service"`
	Status  string `json:"status"`
	Result  string `json:"result"`
	Alarm   bool   `json:"alarm"`
	Tooltip string `json:"tooltip"`
}

type Dock struct {
	Name   string `json:"name"`
	Status string `json:"status"`
}
type El struct {
	Cluster_name                     string  `json:"cluster_name"`
	Status                           string  `json:"status"`
	Timed_out                        bool    `json:"time_out"`
	Number_of_nodes                  int     `json:"number_of_nodes"`
	Number_of_data_nodes             int     `json:"number_of_data_nodes"`
	Active_primary_shards            int     `json:"active_primary_shards"`
	Active_shards                    int     `json:"active_shards"`
	Relocating_shards                int     `json:"relocating_shards"`
	Initializing_shards              int     `json:"initializing_shards"`
	Unassigned_shards                int     `json:"unassigned_shards"`
	Delayed_unassigned_shards        int     `json:"delayed_unassigned_shards"`
	Number_of_pending_tasks          int     `json:"number_of_pending_tasks"`
	Number_of_in_flight_fetch        int     `json:"number_of_in_flight_fetch"`
	Task_max_waiting_in_queue_millis int     `json:"task_max_waiting_in_queue_millis"`
	Active_shards_percent_as_number  float64 `json:"active_shards_percent_as_number"`
}
type Hazel struct {
	Status string `json:"status"`
	State  string `json:"state"`
}
