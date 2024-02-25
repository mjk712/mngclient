package mngclient

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"net"
	"time"
)

const maxPacketSize = 20480

//MNG тип для работы обращения к серверу
type MNG struct {
	conn net.Conn
}

//Connect соединиться с addr, например "elmeh.ru:6777"
func (mng *MNG) Connect(addr string, timeout time.Duration) (err error) {
	if mng == nil {
		return errors.New("Connect(): nil ptr")
	}
	d := net.Dialer{Timeout: timeout}
	mng.conn, err = d.Dial("tcp", addr)
	return
}

//Disconnect закрыть соединение
func (mng *MNG) Disconnect() {
	if mng != nil {
		mng.conn.Close()
	}
}

//обработать принятые данные и извлечь из них список
func (mng MNG) getList(data []byte) (result []MGPListItem) {
	b := bytes.NewReader(data)

	var count uint16

	err := binary.Read(b, binary.LittleEndian, &count)
	if err != nil || count == 0 {
		return
	}

	for i := 0; i < int(count); i++ {
		var newItem MGPListItem

		err1 := binary.Read(b, binary.LittleEndian, &newItem.IMEI)
		err2 := binary.Read(b, binary.LittleEndian, &newItem.LocoType)
		err3 := binary.Read(b, binary.LittleEndian, &newItem.LocoNo)

		readOk := (err1 == nil) && (err2 == nil) && (err3 == nil)
		if readOk {
			result = append(result, newItem)
		} else {
			break
		}
	}

	return
}

func (mng MNG) getRawData(cmd string, param string) (rawData []byte, err error) {
	mng.conn.SetWriteDeadline(time.Now().Add(time.Second * 5))
	_, err = fmt.Fprintf(mng.conn, "%s", cmd+param) //посылаем команду
	if err != nil {
		err = errors.New("getRawData() write error:" + err.Error())
		return
	}

	rcvbuf := make([]byte, maxPacketSize)
	mng.conn.SetReadDeadline(time.Now().Add(time.Second * 5))
	bytesRead, readError := bufio.NewReader(mng.conn).Read(rcvbuf)

	if readError == nil {
		if bytesRead > len(cmd) {
			//cmdBuf := rcvbuf[:len(cmd)] //команда в ответ
			rawData = rcvbuf[len(cmd):] //данные
		}
	} else {
		err = errors.New("receive err")
	}

	return
}

//GetList возвращает список локомотивов с сервера. В случае ошибки возвращает пустой слайс (0 элементов)
func (mng MNG) GetList() (result []MGPListItem) {
	data, err := mng.getRawData("GET_NUMBER", "")
	if err == nil {
		result = mng.getList(data) //обработка пришедших данных
	}
	return
}

/*//StopSend закрыть канал связи
func (mng MNG) StopSend() (err error) {
	// data, err := mng.getRawData("STOP_SEND", "")
	mng.conn.SetWriteDeadline(time.Now().Add(time.Second * 5))
	_, err = fmt.Fprintf(mng.conn, "%s", "STOP_SEND") //посылаем команду
	if err == nil {
		fmt.Println("Канал связи закрыт")
	} else {
		fmt.Printf("StartSend() write error: %v", err)
		err = errors.New("StopSend() write error:" + err.Error())
	}
	return
}

//StartSend открыть канал связи
func (mng MNG) StartSend() (err error) {
	// data, err := mng.getRawData("START_SEND", "")
	mng.conn.SetWriteDeadline(time.Now().Add(time.Second * 5))
	_, err = fmt.Fprintf(mng.conn, "%s", "START_SEND") //посылаем команду
	if err == nil {
		fmt.Println("Канал связи открыт")
	} else {
		fmt.Printf("StartSend() write error: %v", err)
		err = errors.New("StartSend() write error:" + err.Error())
	}
	return
}*/
