package api

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"mxshop-api/userWeb/forms"
	"mxshop-api/userWeb/global/reponse"
	"mxshop-api/userWeb/proto"
	"net/http"
	"strconv"
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
	userConn, err := grpc.Dial(fmt.Sprintf("%s:%d", "127.0.0.1", 50051), grpc.WithInsecure())
	if err != nil {
		zap.S().Errorw("[GetUserList] 连接用户失败", "msg", err.Error())
	}
	userSrvClient := proto.NewUserClient(userConn)
	pn := ctx.DefaultQuery("pn", "0")
	pnInt, _ := strconv.Atoi(pn)
	pSize := ctx.DefaultQuery("pszie", "10")
	pSizeInt, _ := strconv.Atoi(pSize)
	list, err := userSrvClient.GetUserList(context.Background(), &proto.PageInfo{
		Pn:    uint32(pnInt),
		PSize: uint32(pSizeInt),
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

// PassWordLogin 登陆
func PassWordLogin(ctx *gin.Context) {
	passwordLogin := forms.PassWordLoginForm{}
	if err := ctx.ShouldBindJSON(&passwordLogin); err != nil {
		ctx.JSON(http.StatusOK, gin.H{
			"msg": err.Error(),
		})
		return
	}
	userConn, err := grpc.Dial("127.0.0.1:50051", grpc.WithInsecure())
	if err != nil {
		zap.S().Error(err.Error())
	}
	userSerClient := proto.NewUserClient(userConn)
	mobileRsp, err := userSerClient.GetUserByMobile(context.Background(), &proto.MobileRequest{
		Mobile: passwordLogin.Mobile,
	})
	if err != nil {
		if e, ok := status.FromError(err); ok {
			switch e.Code() {
			case codes.NotFound:
				ctx.JSON(http.StatusBadRequest, gin.H{
					"mobile": "用户不存在",
				})
			default:
				ctx.JSON(http.StatusInternalServerError, gin.H{
					"mobile": "登录失败",
				})
			}
			return
		}
	} else {
		passRsp, err := userSerClient.CheckPassWord(context.Background(), &proto.PasswordCheckInfo{
			Password:          passwordLogin.PassWord,
			EncryptedPassword: mobileRsp.Password,
		})
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"mobile": "登录失败1",
			})
		} else {
			if passRsp.Success {
				ctx.JSON(http.StatusOK, gin.H{
					"msg": "登录成功",
				})
			} else {
				ctx.JSON(http.StatusBadRequest, gin.H{
					"mobile": "登录失败2",
				})
			}
		}
	}
}
