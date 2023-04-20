package mydemo

import (
	"context"
	"encoding/json"
	"fmt"
	"glint/logger"
	pb "glint/mesonrpc"
	"io"
	"os"
	"testing"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/types/known/structpb"
)

func Test_customjs_2(t *testing.T) {
	const (
		port = "50051"
	)

	file, err := os.Open("../originUrls.json")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	// 读取JSON数据
	decoder := json.NewDecoder(file)
	originUrls := make(map[string]interface{})
	if err := decoder.Decode(&originUrls); err != nil {
		panic(err)
	}

	//var WG sync.WaitGroup //当前与jackdaw等待同步计数
	var opts []grpc.DialOption
	opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))
	address := "127.0.0.1:" + port
	conn, err := grpc.Dial(address, opts...)
	if err != nil {
		logger.Error("fail to dial: %v", err)
	}

	defer conn.Close()
	client := pb.NewRouteGuideClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Minute)
	defer cancel()

	stream, err := client.RouteChat(ctx)
	if err != nil {
		logger.Error("%s", err.Error())
		return
	}
	waitc := make(chan struct{})
	go func() {
		for {
			in, err := stream.Recv()
			if err == io.EOF {
				close(waitc)
				return
			}
			if err != nil {
				close(waitc)
				return
			}
			//log.Printf("Got Taskid %d Targetid:%d Report:%v", in.GetTaskid(), in.GetTargetid(), in.GetReport().Fields)
			if _, ok := in.GetReport().Fields["vuln"]; ok {
				logger.Success("发现漏洞!")
				PluginId := in.GetReport().Fields["vuln"].GetStringValue()
				__url := in.GetReport().Fields["url"].GetStringValue()
				body := in.GetReport().Fields["body"].GetStringValue()
				hostid := in.GetReport().Fields["hostid"].GetNumberValue()
				//保存数据库
				// Result_id, err := t.Dm.SaveScanResult(
				// 	t.TaskId,
				// 	PluginId,
				// 	true,
				// 	__url,
				// 	base64.StdEncoding.EncodeToString([]byte("")),
				// 	base64.StdEncoding.EncodeToString([]byte(body)),
				// 	int(hostid),
				// )
				// if err != nil {
				// 	logger.Error("plugin::error %s", err.Error())
				// 	return
				// }
				// 存在漏洞信息,打印到漏洞信息
				Element := make(map[string]interface{}, 1)
				Element["status"] = 3
				Element["vul"] = PluginId
				Element["request"] = ""    //base64.StdEncoding.EncodeToString([]byte())
				Element["response"] = body //base64.StdEncoding.EncodeToString([]byte())
				Element["deail"] = in.GetReport().Fields["payload"].GetStringValue()
				Element["url"] = __url
				Element["vul_level"] = in.GetReport().Fields["level"].GetStringValue()
				Element["result_id"] = hostid
				//通知socket消息
				//t.PliuginsMsg <- Element

			} else if _, ok := in.GetReport().Fields["state"]; ok {
				// WG.Done()
			}
		}
	}()

	var length = 0
	//对于目标链接传递
	for _, v := range originUrls {
		if value_list, ok := v.([]interface{}); ok {
			for _, v := range value_list {
				logger.Debug("%v", v)
				length++
			}
		}
	}

	//对于目标链接传递
	for _, v := range originUrls {
		if value_list, ok := v.([]interface{}); ok {
			for _, v := range value_list {
				if value, ok := v.(map[string]interface{}); ok {
					value["isFile"] = false
					value["taskid"] = 1
					value["targetLength"] = length
					m, err := structpb.NewValue(value)
					if err != nil {
						logger.Error("client.RouteChat NewValue m failed: %v", err)
					}
					//WG.Add(1)
					data := pb.JsonRequest{Details: m.GetStructValue()}
					if err := stream.Send(&data); err != nil {
						logger.Error("client.RouteChat JsonRequest failed: %v", err)
					}
				}
			}
		}
	}
	<-waitc
	stream.CloseSend()
	//WG.Wait()
	fmt.Println("finish")

}
