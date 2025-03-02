package kafka

import (
	"errors"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/IBM/sarama"
	"github.com/Popov-Dmitriy-Ivanovich/Diplom_cmd/models"
	tgApi "github.com/Popov-Dmitriy-Ivanovich/Diplom_telegram/api"
)

type BashStatus struct {

}


func ServeStatusMessages() error {
	consumer, err := sarama.NewConsumer([]string{os.Getenv("KAFKA_URL")}, nil)
	if err != nil {
		log.Printf("Error in creating consumer %v", err)
		panic(err)
		return err
	}
	defer consumer.Close()

	statusBashConsumer, err := consumer.ConsumePartition("BashStatus", 0, sarama.OffsetNewest)
	if err != nil {
		log.Printf("Error in consuming partition BashStatus %v", err)
		panic(err)
		return err
	}
	defer statusBashConsumer.Close()
	db := models.GetDb()
	
	for {
		select {
		// (обработка входящего сообщения и отправка ответа в Kafka)
		case msg, ok := <-statusBashConsumer.Messages():
			if !ok {
				log.Printf("connection closed")
				return errors.New("Connection closed")
			}
			log.Println("Got new message from driver")
			
			key := msg.Key

			id, err := strconv.ParseUint(string(key),16,64)
			if err == nil {
				action := models.Action{}
				if err := db.First(&action, id).Error; err == nil {
					status := string(msg.Value)
					log.Println(status)
					if status == "Launched" {
						action.StatusID = 1
						action.LastLaunch = &models.DateOnly{Time: time.Now()}
						tgApi.Notify("Действие " + action.Name + " переведено в статус Запущено")
					} else if status=="Stoped"{
						action.StatusID = 3
						tgApi.Notify("Действие " + action.Name + " переведено в статус Остановлено")
					}else {
						action.StatusID = 4
						log.Println("Действие не было запущено или остановлено: " + status)
					}
					db.Save(&action)
				} else {
					log.Println(err)
				}
			} else {
				log.Println(err)
			}
		}
	}
}