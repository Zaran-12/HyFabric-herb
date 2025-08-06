package controller

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"herbal_demo/dbconfig"
	fabric "herbal_demo/farbic"
	"herbal_demo/utils"
	"io"
	"net/http"
	"os"
	"strings"
)

type GoodsData struct {
	BatchNO        string `form:"batchNo" json:"batchNo" binding:"required"`
	GoodsType      string `form:"productName" json:"productName" binding:"required"`
	ProductionDate string `form:"productionDate" json:"productionDate" binding:"required"`
	Location       string `form:"location" json:"location" binding:"required"`
	TransformGoods string `form:"transformGoods" json:"transformGoods" binding:"required"`
	Description    string `form:"description" json:"description" binding:"required"`
}

func UploadGoodsData(context *gin.Context) {
	var data GoodsData

	// 1. 解析表单数据
	err := context.ShouldBind(&data)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{
			"code": "40001",
			"msg":  "参数解析失败",
		})
		return
	}

	// 2. 文件上传处理
	uploadFile, header, err := context.Request.FormFile("uploadFile")
	defer uploadFile.Close()
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{
			"code": "40002",
			"msg":  "文件解析失败",
		})
		return
	}

	// 生成文件保存路径
	fileNameArray := strings.Split(header.Filename, ".")
	fileType := fileNameArray[1]
	// 写入文件内容
	dir := "./uploadfile/"
	hash := utils.SHA256(header.Filename)
	path := dir + hash + "." + fileType

	// 创建文件
	newfile, err := os.Create(path)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{
			"code": "40003",
			"msg":  "文件保存失败",
		})
		return
	}
	defer newfile.Close()

	//b.将用户上传的内容写到新创建的空文件中
	length, err := io.Copy(newfile, uploadFile)
	if err != nil {
		context.JSON(http.StatusOK, gin.H{
			"code": "30004",
			"msg":  "文件上传失败" + err.Error(),
		})
		return
	}
	fmt.Println(length)

	// 4. 写入 MySQL (存储非关键数据)
	_, err = dbconfig.DB.Exec("INSERT INTO goods_data (batch_id, goods_type, production_date, transform_goods, goods_location, file_path, description) VALUES (?, ?, ?, ?, ?, ?, ?)",
		data.BatchNO, data.GoodsType, data.ProductionDate, data.TransformGoods, data.Location, path, data.Description)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{
			"code": "40006",
			"msg":  "保存数据库失败: " + err.Error(),
		})
		return
	}

	// 调用 Fabric SDK 提交数据
	fabricClient, err := fabric.NewFabricClient()
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"code": "50001", "msg": "初始化区块链客户端失败"})
		return
	}

	_, err = fabricClient.SubmitGoodsData(data.BatchNO, data.GoodsType, data.ProductionDate, data.Location, data.TransformGoods, data.Description)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"code": "50002", "msg": "上传数据到区块链失败: " + err.Error()})
		return
	}
	// 5. 返回结果
	context.JSON(http.StatusOK, gin.H{
		"code": "20000",
		"msg":  "上传成功",
	})
}

// 获取批次记录的 Handler
func GetGoodsRecords(context *gin.Context) {
	rows, err := dbconfig.DB.Query("SELECT batch_id, goods_type, production_date, goods_location, transform_goods, description FROM goods_data")
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{
			"code": "40001",
			"msg":  "获取批次记录失败: " + err.Error(),
		})
		return
	}
	defer rows.Close()

	var records []map[string]interface{}
	for rows.Next() {
		var batchID, goodsType, goodsLocation, transformGoods, description string
		var productionDate string
		err := rows.Scan(&batchID, &goodsType, &productionDate, &goodsLocation, &transformGoods, &description)
		if err != nil {
			context.JSON(http.StatusInternalServerError, gin.H{
				"code": "40002",
				"msg":  "解析批次记录失败: " + err.Error(),
			})
			return
		}

		record := map[string]interface{}{
			"batchID":        batchID,
			"goodsType":      goodsType,
			"productionDate": productionDate,
			"goodsLocation":  goodsLocation,
			"transformGoods": transformGoods,
			"description":    description,
		}
		records = append(records, record)
	}

	context.JSON(http.StatusOK, gin.H{
		"code":    "20000",
		"msg":     "获取批次记录成功",
		"records": records,
	})
}

// 删除批次记录的 Handler
func DeleteGoods(context *gin.Context) {
	// 从路径参数获取批次号
	batchID := context.Param("batchID")
	if batchID == "" {
		context.JSON(http.StatusBadRequest, gin.H{
			"code": "40001",
			"msg":  "批次号不能为空",
		})
		return
	}

	// 删除数据库中的记录
	_, err := dbconfig.DB.Exec("DELETE FROM goods_data WHERE batch_id = ?", batchID)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{
			"code": "40002",
			"msg":  "删除记录失败: " + err.Error(),
		})
		return
	}

	// 3. 调用链码删除区块链中的记录
	fabricClient, err := fabric.NewFabricClient()
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"code": "50001", "msg": "初始化区块链客户端失败"})
		return
	}

	_, err = fabricClient.DeleteGoodsData(batchID)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{
			"code": "50002",
			"msg":  "删除区块链数据失败: " + err.Error(),
		})
		return
	}

	// 4. 返回成功结果
	context.JSON(http.StatusOK, gin.H{
		"code": "20000",
		"msg":  "删除记录成功（数据库 & 区块链）",
	})
}
