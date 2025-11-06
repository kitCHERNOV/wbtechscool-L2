package main

import (
	"fmt"
	"log"
	"net"
)

func handleConnection(conn net.Conn) {
	defer conn.Close()

	//reader := bufio.NewReader(conn)
	buffer := make([]byte, 1024)
	for {
		//message, err := reader.ReadString('\n')
		n, err := conn.Read(buffer)
		if err != nil {
			log.Printf("Ошибка чтения: %v", err)
			return
		}
		message := buffer[:n]
		fmt.Printf("Получено сообщение: %s", message)

		// Отправляем сообщение обратно клиенту
		_, err = conn.Write([]byte(message))
		if err != nil {
			log.Printf("Ошибка записи: %v", err)
			return
		}
		fmt.Printf("Сообщение %s отправлено\n", message)
	}
}

func main() {
	// Указываем адрес и порт для сервера
	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatal(err)
	}
	defer listener.Close()

	log.Println("Сервер запущен на :8080")

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("Ошибка принятия соединения: %v", err)
			continue
		}
		log.Printf("Новое соединение от %s", conn.RemoteAddr().String())

		// Обрабатываем соединение в новой горутине
		go handleConnection(conn)
	}
}
