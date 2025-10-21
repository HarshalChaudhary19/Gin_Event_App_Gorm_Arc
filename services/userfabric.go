package services

import (
	"Gin_Event_App_Arc/fabric"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/hyperledger/fabric-gateway/pkg/client"
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
}

type UserFabricService struct { //Basically Dusre Package ke Interface ko point krna
	Network   *client.Network
	Gateway   *client.Gateway
	Chaincode string
}

func NewUserFabricService(network *client.Network, gateway *client.Gateway, chaincode string) UserFabricServiceI {
	return &UserFabricService{Network: network, Gateway: gateway, Chaincode: chaincode}
}

func (serve *UserFabricService) GetUserFabricService(id int) (*fabric.UserReturn, error) {
	result, err := fabric.EvaluateTransaction(serve.Network, serve.Chaincode, "ReadAsset", strconv.Itoa(id))
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
	result, err := fabric.SubmitTransaction(serve.Network, serve.Chaincode, "UpdateAsset", strconv.Itoa(id), userFabric.UserName, userFabric.Email, userFabric.Password, strconv.Itoa(userFabric.Age))
	if err != nil {
		return err
	}
	fmt.Println("result", string(result))
	return nil
}

func (serve *UserFabricService) DeleteUserFabricService(id int) error {
	_, err := fabric.SubmitTransaction(serve.Network, serve.Chaincode, "DeleteAsset", strconv.Itoa(id))
	if err != nil {
		return err
	}
	return nil
}

func (serve *UserFabricService) InsertUserFabricService(user fabric.User) error {
	hashedPass, err1 := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err1 != nil {
		return err1
	}
	user.Password = string(hashedPass)
	_, err := fabric.SubmitTransaction(serve.Network, serve.Chaincode, "CreateAsset", strconv.Itoa(user.ID), user.UserName, user.Email, user.Password, strconv.Itoa(user.Age))
	if err != nil {
		return err
	}
	return nil
}

func (serve *UserFabricService) GetAllUsersFabricService() ([]fabric.UserReturn, error) {
	result, err := fabric.EvaluateTransaction(serve.Network, serve.Chaincode, "GetAllAssets")
	if err != nil {
		return nil, err
	}
	var usersFabric []fabric.UserReturn
	if err := json.Unmarshal(result, &usersFabric); err != nil {
		return nil, err
	}
	return usersFabric, nil
}

func (serve *UserFabricService) GetFullHistoryService() (map[string][]fabric.UserAssetHistory, error) {
	result, err := fabric.SubmitTransaction(serve.Network, serve.Chaincode, "GetFullUserHistory")
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
	_, err := fabric.SubmitTransaction(serve.Network, serve.Chaincode, "CreateAsset")
	if err != nil {
		return err
	}
	return nil
}
