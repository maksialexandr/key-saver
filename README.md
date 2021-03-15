## Демон для работы с ключами

#### WebSocket
```ws://*****:*****/keys/<action>```

Действия ```<action>```:

```
save - Сохранения
encode - Шифрование
decode - Дешифрование
delete - Удаление
copy - Копирование
```

Payload json:
```
{
    "srcMac" : МАС устройства <string> required,
    "dstMac" : МАС устройства <string>,
    "keys" : [
        {
            "panelCode": Номер квартиры <int> required,
            "value": Значение ключа <string> required
        }
        ...
    ],
    "delete" : Флаг удаления всех ключей(работает при пересохранении) <bool>,
}
```

Response json:
```
{
    "value" : Значение ключа <string>,
    "result" : Ответ <bool>
}
```

#### RabbitMQ
```
EXCHANGE_NAME    "*****"
QUEUE_NAME       "*****"
ROUTING_KEY      "*****"
```

Действия ```<action>```:

```
save - Сохранения 
encode - Шифрование 
decode - Дешифрование 
delete - Удаление 
```

Payload:
```
{
    "srcMac" : МАС устройства <string> required,
    "dstMac" : МАС устройства <string>,
    "keys" : [
        {
            "panelCode": Номер квартиры <int> required,
            "value": Значение ключа <string> required
        }
        ...
    ],
}
```