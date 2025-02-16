package kafka

import (
	"errors"
	"strconv"

	"github.com/IBM/sarama"
	"github.com/Popov-Dmitriy-Ivanovich/Diplom_cmd/models"
)

type BashStatus struct {

}


func ServeStatusMessages() error {
	consumer, err := sarama.NewConsumer([]string{"kafka:9092"}, nil)
	if err != nil {
		return err
	}
	defer consumer.Close()

	statusBashConsumer, err := consumer.ConsumePartition("BashStatus", 0, sarama.OffsetNewest)
	if err != nil {
		return err
	}
	defer statusBashConsumer.Close()
	db := models.GetDb()
	
	for {
		select {
		// (обработка входящего сообщения и отправка ответа в Kafka)
		case msg, ok := <-statusBashConsumer.Messages():
			if !ok {
				return errors.New("Connection closed")
			}
			
			key := msg.Key

			id, err := strconv.ParseUint(string(key),16,64)
			if err == nil {
				action := models.Action{}
				if err := db.First(&action, id).Error; err == nil {
					status := string(msg.Value)
					if status == "Launched" {
						action.StatusID = 1
					} else if status=="Stoped"{
						action.StatusID = 3
					}else {
						action.StatusID = 4
					}
					db.Save(&action)
				}
			}
		}
	}
}