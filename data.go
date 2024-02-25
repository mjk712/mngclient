package mngclient

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
)

/*packetVersion - Номер версии пакета данных. Может принимать значения: 1,2,3,4 (можно игнорировать это)
см. объявление типов в types.go.
По-умолчанию те поля структуры data, информация для которых не пришла, будут заполнены нулевыми значениями.
*/
func (mng MNG) getData(rawData []byte) (ok bool, data TDataPacket, packetVersion int) {

	ok = true

	readerV4 := bytes.NewReader(rawData)
	errV4 := binary.Read(readerV4, binary.LittleEndian, &data)
	if nil != errV4 {
		readerV3 := bytes.NewReader(rawData)
		errV3 := binary.Read(readerV3, binary.LittleEndian, &data.TDataPacketV3)
		if nil != errV3 {
			readerV2 := bytes.NewReader(rawData)
			errV2 := binary.Read(readerV2, binary.LittleEndian, &data.TDataPacketV2)
			if nil != errV2 {
				readerV1 := bytes.NewReader(rawData)
				errV1 := binary.Read(readerV1, binary.LittleEndian, &data.TDataPacketV1)
				if nil != errV1 {
					ok = false
				} else {
					packetVersion = 1
				}
			} else {
				packetVersion = 2
			}
		} else {
			packetVersion = 3
		}
	} else {
		packetVersion = 4
	}

	return
}

/*GetData возвращает данные одного из локомотивов.
li - Элемент списка локомотивов, полученный ранее через функцию GetList().
data - Пакет данных, полученный с локомотива.
*/
func (mng MNG) GetData(li MGPListItem) (data TDataPacket, err error) {
	rawData, err := mng.getRawData("GET_DATA", string(li.IMEI[:]))
	if err != nil {
		err = errors.New("GetData():" + err.Error())
	} else {
		var ok bool
		ok, data, _ = mng.getData(rawData) //обработка пришедших данных
		if !ok {
			err = errors.New("GetData():не получены данные")
		}
	}
	return
}

func (mng MNG) getErrorCode(rawData []byte) (ok bool, data TErrorPacket) {

	ok = true

	reader := bytes.NewReader(rawData)
	err := binary.Read(reader, binary.LittleEndian, &data)
	if nil != err {
		ok = false
	}

	return
}

// GetErrorCode получение списка ошибок
func (mng MNG) GetErrorCode(li MGPListItem) (data TErrorPacket, err error) {
	rawData, err := mng.getRawData("GET_ERROR", string(li.IMEI[:]))
	if err != nil {
		err = errors.New("GetErrorCode():" + err.Error())
	} else {
		var ok bool
		ok, data = mng.getErrorCode(rawData) //обработка пришедших данных
		if !ok {
			err = errors.New("GetErrorCode():не получены данные")
		}
	}
	return
}

// StopSend закрыть канал связи
func (mng MNG) StopSend(li MGPListItem) (err error) {
	_, err = mng.getRawData("STOP_SEND", string(li.IMEI[:]))
	if err == nil {
		fmt.Println("Канал связи закрыт")
	} else {
		fmt.Printf("StopSend() write error: %v", err)
		err = errors.New("StopSend():" + err.Error())
	}
	return
}

// StartSend открыть канал связи
func (mng MNG) StartSend(li MGPListItem) (err error) {
	_, err = mng.getRawData("START_SEND", string(li.IMEI[:]))
	if err == nil {
		fmt.Println("Канал связи открыт")
	} else {
		fmt.Printf("StartSend() write error: %v", err)
		err = errors.New("StartSend():" + err.Error())
	}
	return
}

// после того как открываешь канал связи (отправляешь команду), сервер вычитывает данные которые накопились
// и ждет пока ему отправят команду возобновить обмен (RESUME_SEND), поэтому время перестает обновляться
// в клиентской программе -- сервер не запрашивает текущие данные (так было сделано давно для того чтобы
// те кто проверяют могли убедиться что данные были считаны из памяти (те в клиентской программе выводились только данные полученные из памяти)

// ResumeSend возобновление запроса текущих данных
func (mng MNG) ResumeSend(li MGPListItem) (err error) {
	_, err = mng.getRawData("RESUME_SEND", string(li.IMEI[:]))
	if err == nil {
		fmt.Println("Возобновление запроса текущих данных")
	} else {
		fmt.Printf("ResumeSend() write error: %v", err)
		err = errors.New("ResumeSend():" + err.Error())
	}
	return
}
