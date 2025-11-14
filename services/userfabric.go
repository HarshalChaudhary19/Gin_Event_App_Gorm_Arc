package services

import (
	"Gin_Event_App_Arc/fabric"
	"Gin_Event_App_Arc/models"
	"Gin_Event_App_Arc/repositories"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"github.com/hyperledger/fabric-gateway/pkg/client"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
)

type UserFabricServiceI interface {
	GetUserFabricService(int) (*fabric.UserReturn, error)
	UpdateUserFabricService(fabric.UserUpdate, int) error
	DeleteUserFabricService(int) error
	InsertUserFabricService(fabric.User) error
	GetAllUsersFabricService() ([]fabric.UserReturn, error)
	GetFullHistoryService() (map[string][]fabric.UserAssetHistory, error)
	InitLedgerService() error
	FabricConnectionHelper(string) (*client.Network, string, error)
	RegisterUser(*models.UserSentFabric) (*models.FabricUser, error)
	GetAllRegisteredUsersfromFabric() ([]models.FabricUser, error)
}

type UserFabricService struct { //Basically Dusre Package ke Interface ko point krna
	Network     *client.Network
	Gateway     *client.Gateway
	Chaincode   string
	UserFabRepo repositories.UserFabricRepoI
}

func NewUserFabricService(network *client.Network, gateway *client.Gateway, chaincode string, userfabrepo repositories.UserFabricRepoI) UserFabricServiceI {
	return &UserFabricService{Network: network, Gateway: gateway, Chaincode: chaincode, UserFabRepo: userfabrepo}
}

func (serve *UserFabricService) GetUserFabricService(id int) (*fabric.UserReturn, error) {
	//Changed code
	network, chaincode, err := serve.FabricConnectionHelper("userfirst3Nov")
	if err != nil {
		return nil, err
	}
	logrus.Println("✅ Connected to Hyperledger Fabric successfully")
	//End changed code
	result, err := fabric.EvaluateTransaction(network, chaincode, "ReadAsset", strconv.Itoa(id))
	fmt.Println("bytes wala data", string(result))
	if err != nil {
		return nil, err
	}

	var userFabric fabric.UserReturn
	err = json.Unmarshal(result, &userFabric)
	if err != nil {
		return nil, err
	}
	return &userFabric, nil
}

func (serve *UserFabricService) UpdateUserFabricService(userFabric fabric.UserUpdate, id int) error {
	//Changed code
	network, chaincode, err := serve.FabricConnectionHelper("userfirst3Nov")
	if err != nil {
		return err
	}
	logrus.Println("✅ Connected to Hyperledger Fabric successfully")
	//End changed code
	result, err := fabric.SubmitTransaction(network, chaincode, "UpdateAsset", strconv.Itoa(id), userFabric.UserName, userFabric.Email, userFabric.Password, strconv.Itoa(userFabric.Age))
	if err != nil {
		return err
	}
	fmt.Println("result", string(result))
	return nil
}

func (serve *UserFabricService) DeleteUserFabricService(id int) error {
	//Changed code
	network, chaincode, err := serve.FabricConnectionHelper("userfirst3Nov")
	if err != nil {
		return err
	}
	logrus.Println("✅ Connected to Hyperledger Fabric successfully")
	//End changed code
	_, err2 := fabric.SubmitTransaction(network, chaincode, "DeleteAsset", strconv.Itoa(id))
	if err2 != nil {
		return err
	}
	return nil
}

func (serve *UserFabricService) InsertUserFabricService(user fabric.User) error {
	//Changed code
	network, chaincode, err := serve.FabricConnectionHelper("userfirst3Nov")
	if err != nil {
		fmt.Println("Yahan se error aarha hai")
		return err
	}
	logrus.Println("✅ Connected to Hyperledger Fabric successfully")
	//End changed code
	hashedPass, err1 := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err1 != nil {
		return err1
	}
	user.Password = string(hashedPass)
	_, err2 := fabric.SubmitTransaction(network, chaincode, "CreateAsset", strconv.Itoa(user.ID), user.UserName, user.Email, user.Password, strconv.Itoa(user.Age))
	if err2 != nil {
		return err
	}
	return nil
}

func (serve *UserFabricService) GetAllUsersFabricService() ([]fabric.UserReturn, error) {
	//Changed code
	network, chaincode, err := serve.FabricConnectionHelper("userfirst3Nov")
	if err != nil {
		return nil, err
	}
	logrus.Println("✅ Connected to Hyperledger Fabric successfully")
	//End changed code
	result, err := fabric.EvaluateTransaction(network, chaincode, "GetAllAssets")
	if err != nil {
		fmt.Println("Idhr aarha hai")
		return nil, err
	}
	var usersFabric []fabric.UserReturn
	if err := json.Unmarshal(result, &usersFabric); err != nil {
		fmt.Println("Nhi Idhr aarha hai")
		return nil, err
	}
	return usersFabric, nil
}

func (serve *UserFabricService) GetFullHistoryService() (map[string][]fabric.UserAssetHistory, error) {
	//Changed code
	network, chaincode, err := serve.FabricConnectionHelper("userfirst3Nov")
	if err != nil {
		return nil, err
	}
	logrus.Println("✅ Connected to Hyperledger Fabric successfully")
	//End changed code
	result, err := fabric.SubmitTransaction(network, chaincode, "GetFullUserHistory")
	if err != nil {
		return nil, err
	}
	var userHistory map[string][]fabric.UserAssetHistory
	err = json.Unmarshal(result, &userHistory)
	if err != nil {
		return nil, err
	}
	return userHistory, nil
}

func (serve *UserFabricService) InitLedgerService() error {
	//Changed code
	network, chaincode, err := serve.FabricConnectionHelper("userfirst3Nov")
	if err != nil {
		return err
	}
	logrus.Println("✅ Connected to Hyperledger Fabric successfully")
	//End changed code
	_, err2 := fabric.SubmitTransaction(network, chaincode, "CreateAsset")
	if err2 != nil {
		return err
	}
	return nil
}

func (serve *UserFabricService) FabricConnectionHelper(username string) (*client.Network, string, error) {
	userFabric, err := serve.UserFabRepo.LoadCerts(username)
	if err != nil {
		fmt.Println("Ab yahan se hai error")
		return nil, "", err
	}
	var peerEndPoint string
	encryptionKey := os.Getenv("FABRIC_ENCRYPTION_KEY")
	encKeyByte := []byte(encryptionKey)
	keyPEM, err := fabric.DecryptAESGCM(encKeyByte, userFabric.KeyPEMEnc)
	// fmt.Println("User from fabric", userFabric)
	network, _, err := fabric.ConnectDyn(userFabric.MSPID, peerEndPoint, "mychannel", string(userFabric.CertPEM), string(keyPEM))
	if err != nil {
		fmt.Println("Nhi yahan se hai")
		return nil, "", err
	}
	return network, "basic", nil
}

func (serve *UserFabricService) RegisterUser(userSent *models.UserSentFabric) (*models.FabricUser, error) {
	apiURL := "http://localhost:8081/userfabric/v1/registerUser"
	reqBody, err := json.Marshal(userSent)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal user request: %v", err)
	}
	resp, err := http.Post(apiURL, "application/json", bytes.NewBuffer(reqBody))
	if err != nil {
		return nil, fmt.Errorf("failed to send POST Request: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("API returned status %d: %s", resp.StatusCode, string(body))
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read API response: %v", err)
	}
	// Step 6: Parse the response into your FabricUser struct
	var user models.FabricUser
	if err := json.Unmarshal(body, &user); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %v", err)
	}
	userFabcreated, err := serve.UserFabRepo.SaveRegisteredUser(&user)
	if err != nil {
		return nil, err
	}
	return userFabcreated, nil
}

func (serve *UserFabricService) GetAllRegisteredUsersfromFabric() ([]models.FabricUser, error) {
	usersListFabric, err := serve.UserFabRepo.GetAllUsersRegistered()
	if err != nil {
		return nil, err
	}
	return usersListFabric, nil
}
