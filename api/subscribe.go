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
	"time"
)

type selectUser struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}
type Subscribe struct {
	Rule          string       `json:"rule"`
	SelectUserIds []selectUser `json:"selectUserIds"`
	Label         string       `json:"label"`
}

func createSubscribe(w http.ResponseWriter, r *http.Request) {
	data, err := io.ReadAll(r.Body)
	if err != nil {
		log.Warn(err)
		fmt.Fprint(w, successResponse("body读取失败", false, err))
		return
	}
	var S Subscribe
	if err = json.Unmarshal(data, &S); err != nil {
		log.Warn(err)
		fmt.Fprint(w, successResponse("反序列化失败", false, err))
		return
	}
	bytes, err := json.Marshal(S.SelectUserIds)
	if err != nil {
		log.Warn(err)
		fmt.Fprint(w, successResponse("序列化失败", false, err))
		return
	}

	var sub = &model.Subscribe{
		RuleName:  S.Rule,
		Label:     S.Label,
		Receiver:  string(bytes),
		Timestamp: time.Now().Unix(),
		Enable:    0,
	}
	if err = dao.NewSubscribe().Create(sub); err != nil {
		log.Warn(err)
		fmt.Fprint(w, successResponse("存储失败", false, err))
		return
	}
	fmt.Fprint(w, successResponse("添加爱成功", true, nil))
	return
}
func delSubscribe(w http.ResponseWriter, r *http.Request) {
	TmpID := r.FormValue("id")
	ID, err := strconv.Atoi(TmpID)
	if err != nil {
		log.Warn(err)
		fmt.Fprint(w, successResponse("请指定要删除的id", false, err))
		return
	}
	if err = dao.NewSubscribe().Delete(int32(ID)); err != nil {
		log.Warn(err)
		fmt.Fprint(w, successResponse("删除失败", false, err))
	}
	fmt.Fprint(w, successResponse("删除成功", true, nil))
	return
}
func updateSubscribe(w http.ResponseWriter, r *http.Request) {
	data, err := io.ReadAll(r.Body)
	if err != nil {
		log.Warn(err)
		fmt.Fprint(w, successResponse("body读取失败", false, err))
		return
	}
	var subscribe model.Subscribe
	err = json.Unmarshal(data, &subscribe)
	if err != nil {
		log.Warn(err)
		fmt.Fprint(w, successResponse("json反序列化失败", false, err))
		return
	}
	err = dao.NewSubscribe().Update(&subscribe)
	if err != nil {
		log.Warn(err)
		fmt.Fprint(w, successResponse("修改失败", false, err))
		return
	}
	fmt.Fprint(w, successResponse("更新成功", true, nil))
	return
}
func listSubscribe(w http.ResponseWriter, r *http.Request) {
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

	result, count := dao.NewSubscribe().Pager(PageIndex, PageSize)

	var a = struct {
		Data  interface{} `json:"data"`
		Count int64       `json:"count"`
	}{
		result,
		count,
	}

	fmt.Fprint(w, successResponse("获取成功", true, a))
}
