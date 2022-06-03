package pprof

import (
	"github.com/gorilla/mux"
	"net/http/pprof"
	_ "net/http/pprof"
)

func AttachProfiler(router *mux.Router) {
	router.HandleFunc("/debug/pprof/", pprof.Index)
	router.HandleFunc("/debug/pprof/cmdline", pprof.Cmdline)                  //Cmdline用正在运行的程序的命令行响应，参数用NUL字节分隔
	router.HandleFunc("/debug/pprof/profile", pprof.Profile)                  //默认进行 30s 的 CPU Profiling，得到一个分析用的 profile 文件
	router.HandleFunc("/debug/pprof/symbol", pprof.Symbol)                    //Symbol查找请求中列出的程序计数器，并使用表映射程序计数器响应函数名称
	router.Handle("/debug/pprof/goroutine", pprof.Handler("goroutine"))       //查看当前所有运行的 goroutines 堆栈跟踪
	router.Handle("/debug/pprof/heap", pprof.Handler("heap"))                 //查看活动对象的内存分配情况
	router.Handle("/debug/pprof/threadcreate", pprof.Handler("threadcreate")) //查看创建新OS线程的堆栈跟踪
	router.Handle("/debug/pprof/block", pprof.Handler("block"))               //查看导致阻塞同步的堆栈跟踪
}
