package sugar

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/qingfeng-lian/gosugar/logger"
	"io"
	"net/http"
	"time"

	"go.uber.org/zap"
)

type HttpPostClientConf struct {
	Timeout time.Duration
}

func HttpPost(ctx context.Context, url string, reqData interface{}, respData interface{}, headers map[string]string, clientConf *HttpPostClientConf) (int, error) {
	if ctx == nil {
		ctx = context.Background()
	}
	// 将数据编码为JSON格式
	jsonData, err := json.Marshal(reqData)
	if err != nil {
		return http.StatusBadRequest, err
	}

	// 创建HTTP客户端
	client := &http.Client{}
	if clientConf != nil {
		if clientConf.Timeout > 0 {
			client.Timeout = clientConf.Timeout
		}
	}

	// 构建请求体
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return http.StatusBadRequest, err
	}

	// 设置请求头，表明发送的是JSON格式的数据
	req.Header.Set("Content-Type", "application/json")
	if len(headers) > 0 {
		for k, v := range headers {
			req.Header.Set(k, v)
		}
	}

	// 发送请求
	resp, err := client.Do(req)
	if err != nil {
		return http.StatusBadRequest, err
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	// 读取并打印响应内容
	body, err := io.ReadAll(resp.Body)
	logger.Info(ctx, "http post",
		zap.String("url", url),
		zap.String("jsonData", string(jsonData)),
		zap.String("body", string(body)),
		zap.Any("header", headers),
		zap.Error(err),
	)
	if err != nil {
		return resp.StatusCode, err
	}
	err = json.Unmarshal(body, respData)
	if err != nil {
		logger.Info(ctx, "unmarshal body error",
			zap.String("url", url),
			zap.String("jsonData", string(jsonData)),
			zap.String("body", string(body)),
			zap.Error(err))
		return resp.StatusCode, err
	}
	return resp.StatusCode, err
}
