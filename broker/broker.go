package broker

import (
	"errors"
	"log"
	"net"
	"sync"

	"github.com/eclipse/paho.mqtt.golang/packets"
	"github.com/kundan1301/message-broker/config"
	customHttp "github.com/kundan1301/message-broker/http"
	customRedis "github.com/kundan1301/message-broker/redis"
)

// Broker store broker info
type Broker struct {
	Config  *config.Config
	Clients sync.Map
}

var RunningBroker *Broker

func checkConfig(config *config.Config) error {
	if config.Host == "" {
		return errors.New("Host is required")
	}
	if config.MqttPort == "" {
		return errors.New("MqttPort  is required")
	}
	if config.HttpPort == "" {
		return errors.New("HttpPort is required")
	}
	if config.UseRedisCluster {
		if len(config.RedisClusterOptions.Addrs) == 0 {
			return errors.New("Cluster host is required")
		}
	} else {
		if config.RedisOptions.Addr == "" {
			return errors.New("Redis host is required")
		}
	}
	return nil
}

// NewBroker initialize broker
func NewBroker(config *config.Config) (*Broker, error) {
	err := checkConfig(config)
	if err != nil {
		log.Println("error in creating new broker", err)
		return nil, err
	}
	broker := &Broker{
		Config: config,
	}
	return broker, nil
}

func (b *Broker) Start() {
	startListening(b)
}

func startListening(b *Broker) {
	hostUrl := b.Config.Host + ":" + b.Config.MqttPort
	l, err := net.Listen("tcp", hostUrl)
	if err != nil {
		log.Println("error in listening mqtt", err)
		return
	}
	for {
		conn, err := l.Accept()
		if err != nil {
			log.Println("error in accepting connection", err)
		}
		go handleConnection(b, conn)
	}

}

func handleConnection(b *Broker, conn net.Conn) {
	packet, err := packets.ReadPacket(conn)
	if err != nil {
		log.Println("read connect packet error: ", err)
		return
	}
	if packet == nil {
		log.Println("received nil packet")
		return
	}
	msg, ok := packet.(*packets.ConnectPacket)
	if !ok {
		log.Println("received msg that was not Connect")
		return
	}
	auth := customHttp.CheckAuth(b.Config.AuthUrl, msg.Username, string(msg.Password), msg.ClientIdentifier)
	connack := packets.NewControlPacket(packets.Connack).(*packets.ConnackPacket)
	if !auth {
		connack.ReturnCode = packets.ErrRefusedNotAuthorised
		connack.SessionPresent = false
		connack.Write(conn)
		conn.Close()
		return
	}
	connack.ReturnCode = packets.Accepted
	connack.SessionPresent = msg.CleanSession
	err = connack.Write(conn)
	if err != nil {
		log.Println("error in accepting connection", err)
		return
	}
	info := info{clientID: msg.ClientIdentifier,
		keepalive: msg.Keepalive,
	}
	client := &Client{
		conn,
		info,
	}
	prevBrokerIp := customRedis.CheckPrevConn(info.clientID)
	log.Println("prev", prevBrokerIp)
	if prevBrokerIp != "" {
		ClosePrevConnection(client.info.clientID)
	}
	customRedis.SetNewConnInfo(info.clientID, b.Config.NodeIP+":"+b.Config.HttpPort)
	b.Clients.Store(info.clientID, client)

}

func startTLSListening() {

}

func ClosePrevConnection(clientID string) {
	data, _ := RunningBroker.Clients.Load(clientID)
	RunningBroker.Clients.Delete(clientID)
	if data == nil {
		return
	}
	client, ok := data.(*Client)
	log.Println("closing connection", ok)
	if ok {
		err := client.conn.Close()
		if err != nil {
			log.Println("error in closing prev conn", err)
		}
		client.conn = nil
		client = nil
	}
}
