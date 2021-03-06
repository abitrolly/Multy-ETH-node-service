/*
Copyright 2018 Idealnaya rabota LLC
Licensed under Multy.io license.
See LICENSE for details
*/
package node

import (
	"fmt"
	"net"
	"sync"

	"github.com/KristinaEtc/slf"
	_ "github.com/KristinaEtc/slflog"
	"github.com/Multy-io/Multy-ETH-node-service/eth"
	"github.com/Multy-io/Multy-ETH-node-service/streamer"
	pb "github.com/Multy-io/Multy-back/node-streamer/eth"
	"github.com/Multy-io/Multy-back/store"
	"google.golang.org/grpc"
)

var log = slf.WithContext("NodeClient")

// Multy is a main struct of service

// NodeClient is a main struct of service
type NodeClient struct {
	Config     *Configuration
	Instance   *eth.Client
	GRPCserver *streamer.Server
	Clients    *sync.Map // address to userid
	// BtcApi     *gobcy.API
}

// Init initializes Multy instance
func Init(conf *Configuration) (*NodeClient, error) {
	cli := &NodeClient{
		Config: conf,
	}

	var usersData sync.Map

	usersData.Store("address", store.AddressExtended{
		UserID:       "kek",
		WalletIndex:  1,
		AddressIndex: 2,
	})

	// initail initialization of clients data
	cli.Clients = &usersData
	log.Infof("Users data initialization done")

	// init gRPC server
	lis, err := net.Listen("tcp", conf.GrpcPort)
	if err != nil {
		return nil, fmt.Errorf("failed to listen: %v", err.Error())
	}
	// Creates a new gRPC server

	ethCli := eth.NewClient(&conf.EthConf, cli.Clients)
	if err != nil {
		return nil, fmt.Errorf("eth.NewClient initialization: %s", err.Error())
	}
	log.Infof("ETH client initialization done")

	cli.Instance = ethCli

	s := grpc.NewServer()
	srv := streamer.Server{
		UsersData: cli.Clients,
		M:         &sync.Map{},
		EthCli:    cli.Instance,
		Info:      &conf.ServiceInfo,
	}

	pb.RegisterNodeCommuunicationsServer(s, &srv)
	go s.Serve(lis)

	return cli, nil
}
