package mngclient

//MGPListItem Структура для списка локомотивов на сервере
type MGPListItem struct {
	IMEI     [15]byte
	LocoType uint16
	LocoNo   uint16
}

//TDataPacketV1 версия 1 пакета данных
type TDataPacketV1 struct {
	MGPListItem
	Hour      byte
	Min       byte
	Sec       byte
	Day       byte
	Month     byte
	TabN      uint16
	Speed     uint16
	Dir       byte
	Pwr       byte
	Mileage   uint32
	Press     byte
	Als       byte
	RC        byte
	IfPeriod  byte
	DayUTC    byte
	MonthUTC  byte
	YearUTC   byte
	HourUTC   byte
	MinUTC    byte
	SecUTC    byte
	Latitude  uint32
	Longitude uint32
	Altitude  uint32
	FuelVol   uint32
	FuelTemp  int8
	FuelDens  int8
	Signal    byte
	FuelMass  uint32
	TrainNo   byte
}

//TDataPacketV2 версия 2 пакета данных
type TDataPacketV2 struct {
	TDataPacketV1
	GenA       uint16
	GenV       uint16
	AkbV       uint16
	GenE       uint16
	FreqDis    uint16
	TempWater  int16
	TempOil    int16
	PressFuel  uint16
	PressOil   uint16
	PressTurbo uint16
}

//TDataPacketV3 версия 3 пакета данных
type TDataPacketV3 struct {
	TDataPacketV2
	TabN8        uint32
	FreightMassa uint32
	CountWagons  uint32
	CountAxis    uint32
	HasErrors    bool
	Km           uint16
	Cock395      byte
	Busdata      uint16
	Press2       byte
	Press3       byte
	Epk          byte
}

//TFuelData данные бака топлива
type TFuelData struct {
	FuelVol   uint32
	FuelVol20 uint32
	FuelTemp  int8
	FuelDens  int8
	FuelMass  uint32
}

//TDataPacketV4 версия 4 пакета данных
type TDataPacketV4 struct {
	TDataPacketV3
	TrainNo4  uint16
	Year      byte
	FuelVol20 uint32
	FuelTanks [4]TFuelData
	Dfreqdis  uint16
	Dopsigns  byte
	Falarm    uint16
	Csq       uint16
	Imsi      uint16
	Bindata   uint16
	Statcmd   byte
}

//TDataPacket актуальная версия пакета данных
type TDataPacket TDataPacketV4

//TErrorPacket ошибки
type TErrorPacket struct {
	IMEI       [15]byte
	CountError uint16
	ErrorCodes [255]uint16
}
