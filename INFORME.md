# Informe TP0

## Parte 2: Comunicación

El protocolo de comunicación utilizado es el siguiente:

Los primeros 4 bytes representan un int de 32 bits (big-endian) que indican el tamaño del mensaje (acción + datos) en bytes, luego sigue el mensaje (usando el encoding `UTF-8`) que contiene el tipo de acción a realizar y los datos asociados a la acción. La acción y los datos se separan por `::`.

```
 0 1 2 3
+-+-+-+-+------------------------------------------------------------+
|  LEN  |                 ACTION::DATA                               |
+-+-+-+-+------------------------------------------------------------+
```

Las acciones posibles son:

- `BET`: es un batch de apuestas, los datos son una lista de apuestas separadas por `;` y cada apuesta tiene el formato `agencia:nombre:apellido:documento:nacimiento:numero`.
- `FINISH`: indica que no hay más apuestas, el dato asociado es el id de la agencia.
- `WINNER`:
  - si lo envía el cliente: se consulta los ganadores, el dato asociado es el id de la agencia que consulta.
  - sí lo envía el servidor: se informa los ganadores, los datos son una lista de ganadores separados por `;` y cada ganador tiene el formato `agencia:nombre:apellido:documento:nacimiento:numero`.
- `WINNERWAIT`: respuesta enviada por el servidor cuando se consulta los ganadores y todavía no se obtuvieron las apuestas de todas las agencias. El dato asociado es el tiempo de espera en segundos que se debe esperar para repetir la consulta. Este valor es fijo para todos los requests de consulta de ganadores.

Tanto en `BET` como en `FINISH` el servidor responderá con un mensaje `OK` (usando `\n` como delimitador) si se pudo procesar la acción.

### Ejemplo

Un mensaje de apuestas con 2 apuestas:

```
0000004CBET::1:Juan:Perez:12345678:1980-01-01:1;1:Fulano:Gomez:12345678:1980-01-01:2
```

## Parte 3: Concurrencia

Para la concurrencia se utilizó el módulo _multithreading_ en el servidor. Por cada cliente que se conecta se crea un nuevo thread que se encarga de procesar los mensajes del cliente. El thread principal se encarga de aceptar las conexiones entrantes y crear los threads de los clientes.

Luego, cuando se requiere acceder al archivo de apuestas, se utiliza un _lock_ para evitar que dos threads accedan al mismo tiempo al archivo (tanto al momento de lectura como de escritura).
