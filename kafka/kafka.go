package kafka

import (
	"strconv"

	"github.com/IBM/sarama"
	"github.com/Popov-Dmitriy-Ivanovich/Diplom_cmd/models"
)

func RunAction(action models.Action)(error){
	prod, err := sarama.NewSyncProducer([]string{"kafka:9092"},nil)
	if err != nil {
		return err
	}
	defer prod.Close()
	runReq := &sarama.ProducerMessage{
		Topic: "RunBashAction",
		Key: sarama.StringEncoder(strconv.FormatUint(uint64(action.ID),16)),
		Value: sarama.StringEncoder(action.Cmd),
	}
	_, _, err = prod.SendMessage(runReq)
	if err != nil {
		return err
	}
	return nil
}

func StopAction(action models.Action)(error){
	prod, err := sarama.NewSyncProducer([]string{"kafka:9092"},nil)
	if err != nil {
		return err
	}
	defer prod.Close()
	stopReq := &sarama.ProducerMessage{
		Topic: "StopBashAction",
		Key: sarama.StringEncoder(strconv.FormatUint(uint64(action.ID),16)),
		Value: sarama.StringEncoder(""),
	}
	_, _, err = prod.SendMessage(stopReq)
	if err != nil {
		return err
	}
	return nil
}