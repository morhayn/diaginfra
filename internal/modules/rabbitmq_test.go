package modules

import (
	"testing"
)

func TestRabbitRunString(t *testing.T) {
	rabbit := Rabbitmq{}
	res, err := rabbit.RunString()
	if err != nil {
		t.Fatal(err)
	}
	if res != "sudo rabbitmqctl status" {
		t.Fatal("answer not right ", res)
	}
}
func TestRabbitLogs(t *testing.T) {}
func TestRabbitHandler(t *testing.T) {
	// in := `
	// Status of node rabbit@scuo-rg-demo-tek ...
	// [{pid,606},
	//  {running_applications,
	//  [{rabbitmq_management,"RabbitMQ Management Console","3.7.8"},
	//   {rabbitmq_web_dispatch,"RabbitMQ Web Dispatcher","3.7.8"},
	//   {rabbitmq_stomp,"RabbitMQ STOMP plugin","3.7.8"},
	//   ....
	//   {kernel,"ERTS  CXC 138 10","6.1"}]},
	//  {os,{unix,linux}},
	//  {erlang_version,
	//  "Erlang/OTP 21 [erts-10.1] [source] [64-bit] [smp:2:2] [ds:2:2:10] [async-threads:64] [hipe]\n"},
	//  {memory,
	//  [{connection_readers,0},
	//   {connection_writers,0},
	//   ....
	//   {total,[{erlang,73005248},{rss,82280448},{allocated,83480576}]}]},
	//  {alarms,[]},
	//  {listeners,
	//  [{clustering,25672,"::"},
	//  ...
	//   {http,15672,"::"}]},
	//  {vm_memory_calculation_strategy,rss},
	//  {vm_memory_high_watermark,0.4},
	//  {vm_memory_limit,2500331110},
	//  {disk_free_limit,50000000},
	//  {disk_free,54065967104},
	//  {file_descriptors,
	//  [{total_limit,32668},
	//   {total_used,3},
	//   {sockets_limit,29399},
	//   {sockets_used,0}]},
	//  {processes,[{limit,1048576},{used,603}]},
	//  {run_queue,1},
	//  {uptime,140170},
	//  {kernel,{net_ticktime,60}}]
	// `
}
