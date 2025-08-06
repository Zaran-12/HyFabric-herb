package controller

import (
	"fmt"
	"github.com/gin-gonic/gin"
	fabric "herbal_demo/farbic"
	"net/http"
)

// TraceQuery handles querying batch data from the blockchain
func JGTraceQuery(context *gin.Context) {
	batchNo := context.Query("batchNo") // 获取前端传递的批次号
	if batchNo == "" {
		context.JSON(http.StatusBadRequest, gin.H{
			"code": "40001",
			"msg":  "批次号不能为空",
		})
		return
	}

	// 调用 Fabric SDK 查询区块链数据
	fabricClient, err := fabric.NewFabricClient()
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{
			"code": "50001",
			"msg":  "初始化区块链客户端失败: " + err.Error(),
		})
		return
	}

	queryResult, err := fabricClient.QueryBatchData(batchNo)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{
			"code": "50002",
			"msg":  "查询区块链数据失败: " + err.Error(),
		})
		return
	}

	// 返回查询结果
	context.JSON(http.StatusOK, gin.H{
		"code": "20000",
		"msg":  "查询成功",
		"info": queryResult,
	})
}

// TraceQuery handles querying batch data from the blockchain
func TraceQuery(context *gin.Context) {
	batchNo := context.Query("batchNo") // 获取前端传递的批次号
	if batchNo == "" {
		context.JSON(http.StatusBadRequest, gin.H{
			"code": "40001",
			"msg":  "批次号不能为空",
		})
		return
	}

	// 调用 Fabric SDK 查询区块链数据
	fabricClient, err := fabric.NewFabricClient()
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{
			"code": "50001",
			"msg":  "初始化区块链客户端失败: " + err.Error(),
		})
		return
	}

	queryResult, err := fabricClient.QueryBatchData(batchNo)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{
			"code": "50002",
			"msg":  "查询区块链数据失败: " + err.Error(),
		})
		return
	}

	// 返回查询结果
	context.JSON(http.StatusOK, gin.H{
		"code": "20000",
		"msg":  "查询成功",
		"info": queryResult,
	})
}

// CombinedTraceQuery handles querying both herbal and goods data from the blockchain
func CombinedTraceQuery(context *gin.Context) {
	batchNo := context.Query("batchNo")
	if batchNo == "" {
		context.JSON(http.StatusBadRequest, gin.H{
			"code": "40001",
			"msg":  "批次号不能为空",
		})
		return
	}

	// 1. 初始化 Fabric 客户端
	fabricClient, err := fabric.NewFabricClient()
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{
			"code": "50001",
			"msg":  "初始化区块链客户端失败: " + err.Error(),
		})
		return
	}

	// 2. 分别查询 “草药链码” 和 “商品链码”
	herbalResult, errHerbal := fabricClient.QueryBatchData(batchNo) // herbalcc
	goodsResult, errGoods := fabricClient.QueryGoodsData(batchNo)   // goodscc

	// 3. 根据需求，决定如何处理错误
	//    如果你希望两个都必须成功，否则报错，可以这样：
	if errHerbal != nil && errGoods != nil {
		context.JSON(http.StatusInternalServerError, gin.H{
			"code": "50002",
			"msg":  fmt.Sprintf("查询两个链码都失败: herbalErr=%v, goodsErr=%v", errHerbal, errGoods),
		})
		return
	}

	// 如果你允许一个成功一个失败，那么可以分别记录
	// 这里做示例处理：把错误信息也返回给前端，供前端展示
	responseData := gin.H{}

	if errHerbal == nil {
		responseData["herbal"] = herbalResult
	} else {
		// 仅记录错误信息，也可以直接设为 null
		responseData["herbal"] = fmt.Sprintf("查询herbalcc失败: %v", errHerbal)
	}

	if errGoods == nil {
		responseData["goods"] = goodsResult
	} else {
		responseData["goods"] = fmt.Sprintf("查询goodscc失败: %v", errGoods)
	}

	// 4. 返回合并结果
	context.JSON(http.StatusOK, gin.H{
		"code": "20000",
		"msg":  "查询成功",
		"data": responseData,
	})
}
