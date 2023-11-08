[![Go](https://github.com/blablatov/scada4modbus2dbase/actions/workflows/scada-main4modbus-action.yml/badge.svg)](https://github.com/blablatov/scada4modbus2dbase/actions/workflows/scada-main4modbus-action.yml)
### RU

Демо код обмена данными между master-устройством и slave-устройством по `Modbus TCP (RTU)`.  
Основной модуль-вебсервер `main4modbus` содержит демонстрационный код для парсинга `https`-запросов от панели управления, реализует функции ведущего устройства, вызов методов обмена данными по протоколу `Modbus` с ведомым устройством по `TCP` или `RTU` [Diagslave](https://www.modbusdriver.com/diagslave.html)       
Записывает полученные данные от slave-устройства в `MongoDB`.  
Чат-бот `chatbotserver` отправляет полученные по `modbus` tls-данные, чат-клиентам scada.  


***Схема обмена данными (scheme exchange of data):***

```mermaid
graph TB

  SubGraph1 --> SubGraph1Flow
  subgraph "PLC"
  SubGraph1Flow(Slave modbus client)
  end
  
  SubGraph2 --> SubGraph2Flow
  subgraph "MongoDB"
  SubGraph2Flow(Tables of data SCADA in MongoDB) 
  end

  SubGraph3 --> SubGraph3Flow
  subgraph "SCADA Chat"
  SubGraph3Flow(Chatbot Server)
  SubGraph3Flow -- message with modbus a data --> chatbotclient_operator-1
  SubGraph3Flow -- message with modbus a data --> chatbotclient_operator-2
  SubGraph3Flow -- message with modbus a data --> chatbotclient_operator-3
  end

  subgraph "SCADA"
  Node1[REST-SSL-request `imitation from web browser`] --> Node2[Module webserver-parser-profibus_master `main4modbus` ]
  Node2 --> SubGraph1[Go-request methods of interface or goroutine `modbus2tcp`]
  Node2 --> SubGraph2[Insert-method of interface `modbus2mgo` to MongoDB]
  Node2 --> SubGraph3[Method of interface `chatbotserver`]
  SubGraph1Flow -- response modbus a data --> Node2
end
``` 
 			
Для проверки, запустить модуль `main4sensors` и чат-модуль `chatbotserver`, из строки браузера создать запрос:

	https://localhost:8443/modbus_tcp:ReadCoils:16

или

	https://localhost:8443/modbus_rtu:ReadCoils:16

### EN

Demo code of data exchange between the master device and the slave device via `Modbus TCP (RTU)`.  
The main webserver module `main4modbus` contains a demo code for parsing `https`-requests from the control panel, implements the functions of a master device, calls methods for communicating via the `Modbus` protocol with a slave device via `TCP` or `RTU`.  
Writes received data from the slave to `MongoDB`.  
Chatbot `chatbotserver` sends tls-data received via `modbus` to scada chat clients.  

To check, run the `main4sensors` module and `chatbotserver` chat-module, create a request from the browser line:

	https://localhost:8443/modbus_tcp:ReadCoils:16

or

	https://localhost:8443/modbus_rtu:ReadCoils:16




