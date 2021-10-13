package api

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"mxshop-api/userWeb/global/reponse"
	"mxshop-api/userWeb/proto"
	"net/http"
	"time"
)

// HandlerGrpcErrorHttp 将grpc的code转换成http状态码
func HandlerGrpcErrorHttp(err error, c *gin.Context) {
	if err != nil {
		if e, ok := status.FromError(err); ok {
			switch e.Code() {
			case codes.NotFound:
				c.JSON(http.StatusNotFound, gin.H{
					"msg": e.Message(),
				})
			case codes.Internal:
				c.JSON(http.StatusInternalServerError, gin.H{
					"msg": "内部错误",
				})
			case codes.InvalidArgument:
				c.JSON(http.StatusBadRequest, gin.H{
					"msg": "参数错误",
				})
			default:
				c.JSON(http.StatusInternalServerError, gin.H{
					"msg": "其他错误",
				})
			}
			return
		}
	}
}

func GetUserList(ctx *gin.Context) {
	ip := "0.0.0.0"
	port := 50051
	userConn, err := grpc.Dial(fmt.Sprintf("%s:%d", ip, port), grpc.WithInsecure())
	if err != nil {
		zap.S().Errorw("[GetUserList] 连接用户失败", "msg", err.Error())
	}
	userSrvClient := proto.NewUserClient(userConn)
	list, err := userSrvClient.GetUserList(context.Background(), &proto.PageInfo{
		Pn:    0,
		PSize: 0,
	})
	if err != nil {
		zap.S().Errorw("[GetUserList] 查询用户列表失败", "msg", err.Error())
		HandlerGrpcErrorHttp(err, ctx)
		return
	}
	result := make([]interface{}, 0)
	for _, val := range list.Data {
		user := reponse.UserResponse{
			Id:       val.Id,
			NickName: val.NickName,
			Birthday: time.Time(time.Unix(int64(val.BirthDay), 0)).Format("2006-01-02"),
			Gender:   val.Gender,
			Mobile:   val.Mobile,
		}
		result = append(result, user)
	}
	ctx.JSON(http.StatusOK, result)
}
