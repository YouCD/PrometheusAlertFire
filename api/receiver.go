package api

import (
	"PrometheusAlertFire/model"
	"PrometheusAlertFire/pkg/dao"
	"PrometheusAlertFire/pkg/log"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
)

type response struct {
	Msg  string      `json:"msg"`
	Flag bool        `json:"flag"`
	Data interface{} `json:"data"`
}

func successResponse(msg string, flag bool, data interface{}) string {
	var r = response{
		Msg:  msg,
		Flag: flag,
		Data: data,
	}
	if data == nil {
		r.Data = struct{}{}
	}

	bytes, _ := json.Marshal(r)
	return string(bytes)
}

func createReceiver(w http.ResponseWriter, r *http.Request) {
	data, err := io.ReadAll(r.Body)
	if err != nil {
		log.Warn(err)
		fmt.Fprint(w, successResponse("body读取失败", false, err))
		return
	}
	var receiver model.Receiver
	err = json.Unmarshal(data, &receiver)
	if err != nil {
		log.Warn(err)
		fmt.Fprint(w, successResponse("json反序列化失败", false, err))
		return
	}

	err = dao.NewReceiver().Create(&receiver)
	if err != nil {
		log.Warn(err)
		fmt.Fprint(w, successResponse("添加失败", false, err))
		return
	}

	fmt.Fprint(w, successResponse("添加成功", true, nil))
	return
}

func delReceiver(w http.ResponseWriter, r *http.Request) {
	TmpID := r.FormValue("id")
	ID, err := strconv.Atoi(TmpID)
	if err != nil {
		log.Warn(err)
		fmt.Fprint(w, successResponse("请指定要删除的id", false, err))
		return
	}
	if err = dao.NewReceiver().Delete(int32(ID)); err != nil {
		log.Warn(err)
		fmt.Fprint(w, successResponse("删除失败", false, err))
	}
	fmt.Fprint(w, successResponse("删除成功", true, nil))
	return

}
func updateReceiver(w http.ResponseWriter, r *http.Request) {
	data, err := io.ReadAll(r.Body)
	if err != nil {
		log.Warn(err)
		fmt.Fprint(w, successResponse("body读取失败", false, err))
		return
	}
	var receiver model.Receiver
	err = json.Unmarshal(data, &receiver)
	if err != nil {
		log.Warn(err)
		fmt.Fprint(w, successResponse("json反序列化失败", false, err))
		return
	}
	err = dao.NewReceiver().Update(&receiver)
	if err != nil {
		log.Warn(err)
		fmt.Fprint(w, successResponse("修改失败", false, err))
		return
	}
	fmt.Fprint(w, successResponse("更新成功", true, nil))
	return
}
func listReceiver(w http.ResponseWriter, r *http.Request) {
	TmpPageIndex := r.FormValue("page_index")
	PageIndex, err := strconv.Atoi(TmpPageIndex)
	if err != nil {
		log.Warn(err)
		fmt.Fprint(w, successResponse("请指定pageIndex", false, err))
		return
	}

	TmpPageSize := r.FormValue("page_size")
	PageSize, err := strconv.Atoi(TmpPageSize)
	if err != nil {
		log.Warn(err)
		fmt.Fprint(w, successResponse("请指定pageIndex", false, err))
		return
	}

	result, count := dao.NewReceiver().Pager(PageIndex, PageSize)

	var a = struct {
		Data  interface{} `json:"data"`
		Count int64       `json:"count"`
	}{
		result,
		count,
	}

	fmt.Fprint(w, successResponse("获取成功", true, a))
	return
}

func searchReceiver(w http.ResponseWriter, r *http.Request) {
	name := r.FormValue("name")

	result, err := dao.NewReceiver().Search(name)
	if err != nil {
		log.Warn(err)
		fmt.Fprint(w, successResponse("搜索失败", false, err))
		return
	}
	var a = struct {
		Data interface{} `json:"data"`
	}{
		result,
	}

	fmt.Fprint(w, successResponse("获取成功", true, a))
	return

}
