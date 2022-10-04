module github.com/blablatov/scada4modbus2dbase

go 1.16

require (
	github.com/blablatov/scada4modbus2dbase/gmodbus2tcp v0.0.0-00010101000000-000000000000
	github.com/blablatov/scada4modbus2dbase/modbus2mgo v0.0.0-00010101000000-000000000000
	github.com/blablatov/scada4modbus2dbase/modbus2rtu v0.0.0-00010101000000-000000000000
	github.com/blablatov/scada4modbus2dbase/modbus2tcp v0.0.0-00010101000000-000000000000
)

replace github.com/blablatov/scada4modbus2dbase/gmodbus2tcp => ./gmodbus2tcp

replace github.com/blablatov/scada4modbus2dbase/chatbotclient => ./chatbotclient

replace github.com/blablatov/scada4modbus2dbase/modbus2rtu => ./modbus2rtu

replace github.com/blablatov/scada4modbus2dbase/modbus2mgo => ./modbus2mgo

replace github.com/blablatov/scada4modbus2dbase/modbus2tcp => ./modbus2tcp
