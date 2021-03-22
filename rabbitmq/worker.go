package rabbitmq

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"strconv"

	"github.com/hanguangbaihuo/sparrow_cloud_go/restclient"
	"github.com/streadway/amqp"
)

type Func func(args []interface{}, kwargs map[string]interface{}) error

type Worker struct {
	BrokerHost        string
	BrokerPort        string
	BrokerUserName    string
	BrokerPassword    string
	BrokerVirtualHost string
	BackendSvc        string
	BackendApi        string
	RetryTimes        int64
	IntervalTime      int64
	QueueName         string
	FuncMap           map[string]Func
}

func New(queueName string, funcMap map[string]Func) *Worker {
	if os.Getenv("SC_CONSUMER_RETRY_TIMES") == "" {
		log.Fatalf("SC_CONSUMER_RETRY_TIMES environment not configured")
	}
	retryTimes, err := strconv.ParseInt(os.Getenv("SC_CONSUMER_RETRY_TIMES"), 10, 64)
	if err != nil {
		failOnError(err, "parse SC_CONSUMER_RETRY_TIMES to int type")
	}
	if os.Getenv("SC_CONSUMER_INTERVAL_TIME") == "" {
		log.Fatalf("SC_CONSUMER_INTERVAL_TIME environment not configured")
	}
	intervalTime, err := strconv.ParseInt(os.Getenv("SC_CONSUMER_INTERVAL_TIME"), 10, 64)
	if err != nil {
		failOnError(err, "parse SC_CONSUMER_INTERVAL_TIME to int type")
	}
	return &Worker{
		BrokerHost:        os.Getenv("SC_BROKER_SERVICE_HOST"),
		BrokerPort:        os.Getenv("SC_BROKER_SERVICE_PORT"),
		BrokerUserName:    os.Getenv("SC_BROKER_USERNAME"),
		BrokerPassword:    os.Getenv("SC_BROKER_PASSWORD"),
		BrokerVirtualHost: os.Getenv("SC_BROKER_VIRTUAL_HOST"),
		BackendSvc:        os.Getenv("SC_BACKEND_SERVICE_SVC"),
		BackendApi:        os.Getenv("SC_BACKEND_SERVICE_API"),
		RetryTimes:        retryTimes,
		IntervalTime:      intervalTime,
		QueueName:         queueName,
		FuncMap:           funcMap,
	}
}

type MessageBoby struct {
	Name   string                 `json:"name"`
	Args   []interface{}          `json:"args"`
	Kwargs map[string]interface{} `json:"kwargs"`
}

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}

func setParentOptions(taskid interface{}, taskname string) {
	parentOptions := struct {
		ID   interface{} `json:"id"`
		Code string      `json:"code"`
	}{
		ID:   taskid,
		Code: taskname,
	}
	data, err := json.Marshal(parentOptions)
	if err != nil {
		log.Printf("set parent options error %s\n", err)
	}
	os.Setenv("SPARROW_TASK_PARENT_OPTIONS", string(data))
}

func getTaskInfo(headers amqp.Table, body MessageBoby) string {
	deliveryInfo := headers["delivery_info"].(amqp.Table)
	var parentID, parentCode interface{}
	parentOptions, ok := headers["parent_options"].(amqp.Table)
	if ok {
		parentID = parentOptions["id"]
		parentCode = parentOptions["code"]
	}

	taskInfo := struct {
		ID          interface{}            `json:"id"`
		Name        string                 `json:"name"`
		TaskArgs    []interface{}          `json:"task_args"`
		TaskKwargs  map[string]interface{} `json:"task_kwargs"`
		Origin      interface{}            `json:"origin"`
		CreatedTime interface{}            `json:"created_time"`
		Exchange    interface{}            `json:"exchange"`
		RoutingKey  interface{}            `json:"routing_key"`
		IsSent      bool                   `json:"is_sent"`
		ParentID    interface{}            `json:"parent_id"`
		ParentCode  interface{}            `json:"parent_code"`
	}{
		ID:          headers["task_id"],
		Name:        body.Name,
		TaskArgs:    body.Args,
		TaskKwargs:  body.Kwargs,
		Origin:      headers["origin"],
		CreatedTime: headers["created_time"],
		Exchange:    deliveryInfo["exchange"],
		RoutingKey:  deliveryInfo["routing_key"],
		IsSent:      true,
		ParentID:    parentID,
		ParentCode:  parentCode,
	}
	data, err := json.Marshal(taskInfo)
	if err != nil {
		log.Printf("marshal taskInfo error %s\n", err)
	}
	return string(data)
}

func (w *Worker) updateTaskResult(taskid interface{}, status string, result string, traceback string, taskInfo string) {
	data := struct {
		TaskID    interface{} `json:"task_id"`
		Consumer  string      `json:"consumer"`
		Status    string      `json:"status"`
		Result    string      `json:"result"`
		TraceBack string      `json:"traceback"`
		TaskInfo  string      `json:"task_info"`
	}{
		TaskID:    taskid,
		Consumer:  w.QueueName,
		Status:    status,
		Result:    result,
		TraceBack: traceback,
		TaskInfo:  taskInfo,
	}
	// log.Printf("request data %v\n", data)
	// TODO: add retry time and interval time
	res, err := restclient.Post(w.BackendSvc, w.BackendApi, data)
	if err != nil {
		log.Printf("Update task database info task_id %v error: %s\n", taskid, err)
		return
	}
	log.Printf("Update task database info task_id %v: %s, %d\n", taskid, res.Body, res.Code)
}

func (w *Worker) Run() {
	amqpURIFormat := "amqp://%s:%s@%s:%s/%s"
	// amqpURI := "amqp://hg_test:jft87JheHe23@39.103.7.185:5672/sparrow_test"
	amqpURI := fmt.Sprintf(amqpURIFormat, w.BrokerUserName, w.BrokerPassword, w.BrokerHost, w.BrokerPort, w.BrokerVirtualHost)
	connection, err := amqp.Dial(amqpURI)
	failOnError(err, "dial amqp error")
	defer connection.Close()

	log.Printf("got Connection, getting Channel")
	channel, err := connection.Channel()
	failOnError(err, "open channel error")
	defer channel.Close()

	// err = channel.ExchangeDeclare(
	// 	"logs_topic", // name
	// 	"topic",      // type
	// 	true,         // durable
	// 	false,        // auto-deleted
	// 	false,        // internal
	// 	false,        // noWait
	// 	nil,          // arguments
	// )
	// failOnError(err, "declare exchange error")

	q, err := channel.QueueDeclare(
		w.QueueName, // name
		true,        // durable
		false,       // delete when unused
		false,       // exclusive
		false,       // no-wait
		nil,         // arguments
	)
	failOnError(err, "Failed to declare a queue")

	// fair dispatch for worker
	err = channel.Qos(
		1,     // prefetch count
		0,     // prefetch size
		false, // global
	)
	failOnError(err, "Failed to set QoS")

	msgs, err := channel.Consume(
		q.Name, // queue
		"",     // consumer
		true,   // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	failOnError(err, "Failed to register a consumer")

	for d := range msgs {
		// logout all data
		// log.Printf("reveive raw msg: %v\n", d)
		// log.Printf("header is %v\n", d.Headers)

		data, err := base64.StdEncoding.DecodeString(string(d.Body))
		if err != nil {
			log.Printf("base64 decode body error: %s", err)
			continue
		}
		// log.Printf("data is %v\n", string(data))

		// logout info for search
		log.Printf("Delivery tag: %d Message body: %s Message Header: %v Task_id: %v\n", d.DeliveryTag, data, d.Headers, d.Headers["task_id"])

		var messageBody MessageBoby
		err = json.Unmarshal(data, &messageBody)
		if err != nil {
			log.Printf("unmarshal message body error: %s\n", err)
			continue
		}
		// log.Printf("message body %v\n", messageBody)

		setParentOptions(d.Headers["task_id"], messageBody.Name)

		// execute message code func
		status := "SUCCESS"
		var traceback string

		funcDo, ok := w.FuncMap[messageBody.Name]
		if !ok {
			log.Printf("not found '%s' function\n", messageBody.Name)
			status = "FAILURE"
			traceback = "not found " + messageBody.Name + "function"
			// continue
		} else {
			err = funcDo(messageBody.Args, messageBody.Kwargs)
			if err != nil {
				status = "FAILURE"
				traceback = err.Error()
				log.Printf("exec %s function error %s\n", messageBody.Name, err)
			}
		}

		taskInfo := getTaskInfo(d.Headers, messageBody)
		w.updateTaskResult(d.Headers["task_id"], status, "", traceback, taskInfo)
		// d.Ack(false)
	}
}

type Result struct {
	Status    string
	TraceBack string
}

func doWork(m MessageBoby, fm map[string]Func, r chan Result) {
	status := "SUCCESS"
	var traceback string
	funcDo, ok := fm[m.Name]
	if !ok {
		log.Printf("not found '%s' function\n", m.Name)
		status = "FAILURE"
		traceback = "not found " + m.Name + "function"
	} else {
		err := funcDo(m.Args, m.Kwargs)
		if err != nil {
			status = "FAILURE"
			traceback = err.Error()
			log.Printf("exec %s function error %s\n", m.Name, err)
		}
	}
	r <- Result{Status: status, TraceBack: traceback}
}

func (w *Worker) sendResult(taskid interface{}, taskInfo string, r chan Result) {
	res := <-r
	w.updateTaskResult(taskid, res.Status, "", res.TraceBack, taskInfo)
}
