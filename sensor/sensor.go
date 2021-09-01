package sensor

import (
	"bytes"
	"encoding/gob"
	"log"
	"math/rand"
	"strconv"
	"time"

	"github.com/streadway/amqp"

	"gh-assan/rmsq/utils"
)

type Sensor struct {
	Name     string
	Freq     int
	Max      float64
	Min      float64
	StepSize float64
	value    float64
	nom      float64
	r        *rand.Rand
}

var url = "amqp://rabbitmq:rabbitmq@localhost:5672/"

type SensorMessage struct {
	Name      string
	Value     float64
	Timestamp time.Time
}

func NewSensor(
	name string,
	freq int,
	max float64,
	min float64,
	stepSize float64,
) *Sensor {

	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	value := r.Float64()*(max-min) + min
	nom := (max-min)/2 + min

	sensor := Sensor{
		Name:     name,
		Freq:     freq,
		Max:      max,
		Min:      min,
		StepSize: stepSize,

		value: value,
		nom:   nom,
		r:     r,
	}

	return &sensor
}

func (sensor *Sensor) SendReading() {

	connection, channel := utils.GetChannel(url)
	defer connection.Close()
	defer channel.Close()

	dataQueue := utils.GetQueue(sensor.Name, channel)

	dur, _ := time.ParseDuration(strconv.Itoa(1000/int(sensor.Freq)) + "ms")

	signal := time.Tick(dur)

	buf := new(bytes.Buffer)
	enc := gob.NewEncoder(buf)

	for range signal {

		sensor.calcValue()
		reading := SensorMessage{
			Name:      sensor.Name,
			Value:     sensor.value,
			Timestamp: time.Now(),
		}

		buf.Reset()
		enc = gob.NewEncoder(buf)
		enc.Encode(reading)

		msg := amqp.Publishing{
			Body: buf.Bytes(),
		}

		channel.Publish(
			"",             //exchange string,
			dataQueue.Name, //key string,
			false,          //mandatory bool,
			false,          //immediate bool,
			msg)            //msg amqp.Publishing)

		log.Printf("Reading sent. Value: %v\n", sensor.value)
	}

}

func (sensor *Sensor) calcValue() {
	var maxStep, minStep float64

	if sensor.value < sensor.nom {
		maxStep = sensor.StepSize
		minStep = -1 * sensor.StepSize * (sensor.value - sensor.Min) / (sensor.nom - sensor.Min)
	} else {
		maxStep = sensor.StepSize * (sensor.Max - sensor.value) / (sensor.Max - sensor.nom)
		minStep = -1 * sensor.StepSize
	}

	sensor.value += sensor.r.Float64()*(maxStep-minStep) + minStep
}
