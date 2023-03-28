package rabbitmq

import (
	"fmt"
	feishu_service "github.com/hxkjason/sgc/services/feifu_service"
	uuid "github.com/satori/go.uuid"
	"github.com/streadway/amqp"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

const (
	ChannelBufferLength = 100
	ReceiverNum         = 5
	AckerNum            = 10

	NotifySuccess = 1
	NotifyFailure = 0
)

var (
	GMQConn = map[string]MqConn{}
)

type (
	MqConn struct {
		Name        string
		Connection  *amqp.Connection
		Channel     *amqp.Channel
		IsConnected bool
	}

	NotifyResponse int

	Message struct {
		ConnName       string         // 连接Mq的名称
		QueueConfig    QueueConfig    // 队列的配置
		AmqpDelivery   *amqp.Delivery // message read from rabbitmq
		NotifyResponse NotifyResponse // notify result from callback url
	}

	QueueConfig struct {
		QueueName       string
		RoutingKey      []string
		BindingExchange string
	}
)

func InitMQPublishConn(connName, connUrl string, i int) {
	closeConnChannel := make(chan *amqp.Error)
	go func() { // 连接突然中断会走这里发起重新连接
		err := <-closeConnChannel
		if err != nil {
			fmt.Println("rabbitmq conn close")
		}
		i++
		InitMQPublishConn(connName, connUrl, i)
	}()
	defer func() { // 每次连接时如果出错，会从这里继续发起连接
		if err := recover(); err != nil {
			errMsg := fmt.Sprintf("%s 第 %d 次连接失败,err: %s , try reconnect\n", connName, i, err)
			fmt.Printf(errMsg)
			time.Sleep(3 * time.Second)
			i++
			feishu_service.SendDevopsMsg(errMsg, "", "")
			InitMQPublishConn(connName, connUrl, i)
		}
	}()

	conn, err := amqp.Dial(connUrl)
	fmt.Println(connUrl)
	if err != nil {
		panic(err)
	}
	fmt.Printf("%s 第 %d 次连接成功！\n", connName, i)
	channel, err := conn.Channel()
	if err != nil {
		fmt.Println(connUrl, "channel.open.err:", err)
		panic(err)
	}
	GMQConn[connName] = MqConn{Name: connName, Connection: conn, Channel: channel, IsConnected: true}
	if i > 1 && err == nil {
		feishu_service.SendDevopsMsg(connName+" rabbitMQ重新连接成功", "", "")
	}
	i = 0 // 连接成功后,将连接次数标识初始化为0,重新连接时会自动加1
	conn.NotifyClose(closeConnChannel)
}

func HandleSignal(done chan<- struct{}) {
	chanSigs := make(chan os.Signal, 1)
	signal.Notify(chanSigs, syscall.SIGQUIT)

	go func() {
		sig := <-chanSigs
		if sig != nil {
			fmt.Println("received a signal:, close done channel", sig)
			close(done)
		}
	}()
}

func (qc QueueConfig) WorkerQueueName() string {
	return qc.QueueName
}

func (qc QueueConfig) WorkerExchangeName() string {
	if qc.BindingExchange == "" {

	}
	return qc.BindingExchange
}

func (qc QueueConfig) DeclareExchange(channel *amqp.Channel) {
	exchanges := []string{
		qc.WorkerExchangeName(),
	}

	for _, e := range exchanges {
		fmt.Printf("declaring exchange: %s\n", e)

		err := channel.ExchangeDeclare(e, "topic", true, false, false, false, nil)
		PanicOnError(err)
	}
}

func (qc QueueConfig) DeclareQueue(channel *amqp.Channel) {
	var err error

	// 定义工作队列
	workerQueueOptions := map[string]interface{}{}
	_, err = channel.QueueDeclare(qc.WorkerQueueName(), true, false, false, false, workerQueueOptions)
	PanicOnError(err)

	for _, key := range qc.RoutingKey {
		err = channel.QueueBind(qc.WorkerQueueName(), key, qc.WorkerExchangeName(), false, nil)
		PanicOnError(err)
	}
}

func ReceiveGhostMessage(connName string, queues []*QueueConfig, done <-chan struct{}) <-chan Message {
	// connName = ConnSyncTripWhTask
	out := make(chan Message, ChannelBufferLength)
	var wg sync.WaitGroup

	receiver := func(qc QueueConfig) {
		defer wg.Done()

	RECONNECT:
		for {
			channel := GMQConn[connName].Channel
			if channel == nil {
				fmt.Println(connName, "channel is nil, RECONNECT")
				time.Sleep(5 * time.Second)
				continue RECONNECT
			}

			msgs, err := channel.Consume(
				qc.WorkerQueueName(), // queue
				"",                   // consumer
				false,                // auto-ack
				false,                // exclusive
				false,                // no-local
				false,                // no-wait
				nil,                  // args
			)
			PanicOnError(err)

			for {
				select {
				case msg, ok := <-msgs:
					if !ok {
						fmt.Println(connName, "receiver: channel is closed, maybe lost connection")
						time.Sleep(5 * time.Second)
						continue RECONNECT
					}
					msg.MessageId = uuid.NewV4().String()
					message := Message{connName, qc, &msg, NotifyFailure}
					out <- message
					//message.Printf("receiver: received msg=====")
				case <-done:
					fmt.Println("receiver: received a done signal")
					return
				}
			}
		}
	}

	for _, queue := range queues {
		wg.Add(ReceiverNum)
		for i := 0; i < ReceiverNum; i++ {
			go receiver(*queue)
		}
	}

	go func() {
		wg.Wait()
		fmt.Println("all receiver is done, closing channel")
		close(out)
	}()

	return out
}

func AckMessage(in <-chan Message) <-chan Message {
	out := make(chan Message)
	var wg sync.WaitGroup

	acker := func() {
		defer wg.Done()

		for m := range in {
			if m.IsNotifySuccess() {
				m.Ack()
				//m.Printf("Ack: true")
			} else {
				m.Reject()
				m.Printf("该消息没有 Ack: true")
			}
		}
	}

	for i := 0; i < AckerNum; i++ {
		wg.Add(1)
		go acker()
	}

	go func() {
		wg.Wait()
		fmt.Println("all acker is done, close out")
		close(out)
	}()

	return out
}

// Publish 发布消息
func (MqConn MqConn) Publish(exchangeName string, routingKey string, data []byte) error {
	return MqConn.Channel.Publish(
		exchangeName, // Exchange
		routingKey,   // Routing key
		false,        // Mandatory
		false,        // Immediate
		amqp.Publishing{
			DeliveryMode: amqp.Persistent,
			ContentType:  "application/json",
			Body:         data,
			Timestamp:    time.Now(),
		},
	)
}

func DeclareQueue(ch *amqp.Channel, queueName string, durable bool) (amqp.Queue, error) {
	queue, err := ch.QueueDeclare(
		queueName, // name
		durable,   // durable
		false,     // delete when unused
		false,     // exclusive
		false,     // no-wait
		nil,       // arguments
	)

	return queue, err
}

func DeclareExchange(ch *amqp.Channel, exchangeName, kind string) error {
	err := ch.ExchangeDeclare(
		exchangeName, // name
		kind,         // type
		true,         // durable
		false,        // auto-deleted
		false,        // internal
		false,        // no-wait
		nil,          // arguments
	)

	return err
}

func DealMQMessage(messageContent []byte, messageId string) {
	fmt.Println("deal msg:", messageId, string(messageContent))
}

func (m Message) IsNotifySuccess() bool {
	return m.NotifyResponse == NotifySuccess
}

func (m Message) Ack() error {
	err := m.AmqpDelivery.Ack(false)
	LogOnError(err)
	return err
}

func (m Message) Reject() error {
	m.Printf("acker: reject message")
	err := m.AmqpDelivery.Reject(false)
	LogOnError(err)
	return err
}

func PanicOnError(err error) {
	if err != nil {
		panic(err)
	}
}

func LogOnError(err error) {
	if err != nil {
		fmt.Printf("ERROR - %s\n", err)
	}
}

func (m Message) Printf(v ...interface{}) {
	msg := m.AmqpDelivery

	var vv []interface{}
	vv = append(vv, msg.MessageId, "queueName:"+m.QueueConfig.QueueName, msg.RoutingKey)
	vv = append(vv, v[1:]...)

	fmt.Printf("[%s] [%s] [%s] "+v[0].(string), vv...)
}
