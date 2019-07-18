package config

import "fmt"

func ConstructAMQPURL(user string, pass string, host string, port string) string {
	return fmt.Sprintf("amqp://%s:%s@%s:%s", user, pass, host, port)
}

func MongoURL(hosts string, port string, username string, password string) string {
	return fmt.Sprintf("mongodb://%s:%s@%s:%s", username, password, hosts, port)
}
