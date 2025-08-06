package fabric

import (
	"bytes"
	"crypto/x509"
	"encoding/json"
	"fmt"
	"github.com/hyperledger/fabric-gateway/pkg/client"
	"github.com/hyperledger/fabric-gateway/pkg/hash"
	"github.com/hyperledger/fabric-gateway/pkg/identity"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"os"
	"path"
	"time"
)

const (
	mspID        = "Org1MSP"
	certPath     = "./sdk/wallet/appUser/msp/signcerts"
	keyPath      = "./sdk/wallet/appUser/msp/keystore"
	tlsCertPath  = "./sdk/wallet/tls/peer0-org1-ca.crt"
	peerEndpoint = "dns:///192.168.154.130:7051"
	gatewayPeer  = "peer0.org1.example.com"
	channelName  = "mychannel"

	// 不再固定使用一个链码名称, 改为定义两个常量
	herbalChaincodeName = "herbalcc"
	goodsChaincodeName  = "goodscc"
)

// FabricClient 同时存储对 herbalcc 和 goodscc 的合约引用
type FabricClient struct {
	Gateway        *client.Gateway // 保留以便需要获取更多信息
	HerbalContract *client.Contract
	GoodsContract  *client.Contract
}

// NewFabricClient 初始化连接，并分别获取对 herbalcc 和 goodscc 的合约引用
func NewFabricClient() (*FabricClient, error) {
	clientConnection := newGrpcConnection()
	id := newIdentity()
	sign := newSign()

	gw, err := client.Connect(
		id,
		client.WithSign(sign),
		client.WithHash(hash.SHA256),
		client.WithClientConnection(clientConnection),
		client.WithEvaluateTimeout(5*time.Second),
		client.WithEndorseTimeout(15*time.Second),
		client.WithSubmitTimeout(5*time.Second),
		client.WithCommitStatusTimeout(1*time.Minute),
	)
	if err != nil {
		return nil, err
	}

	// 获取 channel
	network := gw.GetNetwork(channelName)

	// 分别获取 herbalcc 和 goodscc 的合约引用
	herbalContract := network.GetContract(herbalChaincodeName)
	goodsContract := network.GetContract(goodsChaincodeName)

	return &FabricClient{
		Gateway:        gw,
		HerbalContract: herbalContract,
		GoodsContract:  goodsContract,
	}, nil
}

// ------------------ 以下是对 herbalcc 的操作 -------------------

// SubmitBatchData 提交草药信息 (herbalcc)
func (f *FabricClient) SubmitBatchData(batchNo, cropType, plantDate, location, transformHerbal, description string) (string, error) {
	_, err := f.HerbalContract.SubmitTransaction("StoreBatch", batchNo, cropType, plantDate, location, transformHerbal, description)
	if err != nil {
		return "", fmt.Errorf("failed to submit transaction(StoreBatch): %w", err)
	}
	return "Transaction submitted successfully for herbalcc", nil
}

// QueryBatchData 查询草药信息 (herbalcc)
func (f *FabricClient) QueryBatchData(batchNo string) (string, error) {
	evaluateResult, err := f.HerbalContract.EvaluateTransaction("QueryBatch", batchNo)
	if err != nil {
		return "", fmt.Errorf("failed to evaluate transaction(QueryBatch): %w", err)
	}
	return string(evaluateResult), nil
}

// ------------------ 以下是对 goodscc 的操作 -------------------

// SubmitGoodsData 提交药品信息 (goodscc)
func (f *FabricClient) SubmitGoodsData(batchNo, productName, productionDate, location, transformGoods, description string) (string, error) {
	_, err := f.GoodsContract.SubmitTransaction("AddGoods", batchNo, productName, productionDate, location, transformGoods, description)
	if err != nil {
		return "", fmt.Errorf("failed to submit transaction(AddGoods): %w", err)
	}
	return "Transaction submitted successfully for goodscc", nil
}

// QueryGoodsData 查询药品信息 (goodscc)
func (f *FabricClient) QueryGoodsData(batchNo string) (string, error) {
	evaluateResult, err := f.GoodsContract.EvaluateTransaction("QueryGoods", batchNo)
	if err != nil {
		return "", fmt.Errorf("failed to evaluate transaction(QueryGoods): %w", err)
	}
	return string(evaluateResult), nil
}

// DeleteBatchData 调用链码方法 DeleteBatch
func (f *FabricClient) DeleteBatchData(batchNo string) (string, error) {
	_, err := f.HerbalContract.SubmitTransaction("DeleteBatch", batchNo)
	if err != nil {
		return "", fmt.Errorf("failed to submit transaction(DeleteBatch): %w", err)
	}
	return "Transaction submitted successfully", nil
}

func (f *FabricClient) DeleteGoodsData(batchNo string) (string, error) {
	_, err := f.GoodsContract.SubmitTransaction("DeleteGoods", batchNo)
	if err != nil {
		return "", fmt.Errorf("failed to submit transaction(DeleteGoods): %w", err)
	}
	return "Transaction submitted successfully", nil
}

// Helper functions for connection, identity, etc.
// newGrpcConnection creates a gRPC connection to the Gateway server.
func newGrpcConnection() *grpc.ClientConn {
	certificatePEM, err := os.ReadFile(tlsCertPath)
	if err != nil {
		panic(fmt.Errorf("failed to read TLS certificate file: %w", err))
	}

	certificate, err := identity.CertificateFromPEM(certificatePEM)
	if err != nil {
		panic(err)
	}

	certPool := x509.NewCertPool()
	certPool.AddCert(certificate)
	transportCredentials := credentials.NewClientTLSFromCert(certPool, gatewayPeer)

	// 注意：在 v0.7.x 或更早的 fabric-gateway 版本中，可能需要使用 `grpc.Dial` 自己封装
	// 最新版本中提供了 client.NewConnection(...) 或类似API
	connection, err := grpc.Dial(
		peerEndpoint,
		grpc.WithTransportCredentials(transportCredentials),
	)
	if err != nil {
		panic(fmt.Errorf("failed to create gRPC connection: %w", err))
	}

	return connection
}

// newIdentity creates a client identity for this Gateway connection using an X.509 certificate.
func newIdentity() *identity.X509Identity {
	certificatePEM, err := readFirstFile(certPath)
	if err != nil {
		panic(fmt.Errorf("failed to read certificate file: %w", err))
	}

	certificate, err := identity.CertificateFromPEM(certificatePEM)
	if err != nil {
		panic(err)
	}

	id, err := identity.NewX509Identity(mspID, certificate)
	if err != nil {
		panic(err)
	}

	return id
}

// newSign creates a function that generates a digital signature from a message digest using a private key.
func newSign() identity.Sign {
	privateKeyPEM, err := readFirstFile(keyPath)
	if err != nil {
		panic(fmt.Errorf("failed to read private key file: %w", err))
	}

	privateKey, err := identity.PrivateKeyFromPEM(privateKeyPEM)
	if err != nil {
		panic(err)
	}

	sign, err := identity.NewPrivateKeySign(privateKey)
	if err != nil {
		panic(err)
	}

	return sign
}

// readFirstFile 读取给定目录下的第一个文件（cert.pem、private.key等）
func readFirstFile(dirPath string) ([]byte, error) {
	dir, err := os.Open(dirPath)
	if err != nil {
		return nil, err
	}
	defer dir.Close()

	fileNames, err := dir.Readdirnames(1)
	if err != nil {
		return nil, err
	}

	return os.ReadFile(path.Join(dirPath, fileNames[0]))
}

// formatJSON 格式化 JSON 输出
func formatJSON(data []byte) string {
	var prettyJSON bytes.Buffer
	if err := json.Indent(&prettyJSON, data, "", "  "); err != nil {
		panic(fmt.Errorf("failed to parse JSON: %w", err))
	}
	return prettyJSON.String()
}
