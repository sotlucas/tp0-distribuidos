# Informe TP0

## Introducción

Los ejercicios de docker se encuentran en el branch `master`. A partir del ejercicio 5, se utilizó un branch separado para cada ejercicio (siguiendo la nomenclatura `ejercicio-<nro>`). Para ver el código de cada ejercicio, se debe cambiar al branch correspondiente. Cada branch de ejercicio sale del ejercicio anterior por lo cual se puede ver la evolución del código de un ejercicio a otro. El branch `ejercicio-8` contiene el codigo final.

El protocolo utilizado fue evolucionando de manera incremental a medida que se avanzaba con los ejercicios.

Para levantar la aplicación y ver los logs ejecutar lo siguiente (es lo mismo para todos los ejercicios, excepto en los casos que se pidió la creación de un script):

```bash
make docker-compose-up
make docker-compose-logs
```

## Parte 1: Docker

### Ejercicio 1.1

Se creó el script `generate_dc.py` que genera un archivo `docker-compose-dev.yaml` con la configuración de los servicios. El script recibe como parámetro la cantidad de clientes que se quieren levantar.

El mismo se puede ejecutar con el siguiente comando (para 5 clientes):

```bash
$ python3 ./utils/generate_dc.py 5
```

### Ejercicio 3

El script `check_server.sh` se puede utilizar para corroborar que el servidor está funcionando correctamente.

Para ello, levanta un contenedor con alpine linux, el cuál se conecta a la misma red que el servidor. Luego, utiliza netcat para conectarse al servidor y enviarle un mensaje de prueba. Si el servidor responde con el mismo mensaje, entonces se considera que el servidor está funcionando correctamente (considerando que a este punto el servidor funciona como un echo server).

El mismo se puede ejecutar con el siguiente comando:

```bash
$ ./utils/check_server.sh
```

**Nota:** este script funciona hasta antes del ejercicio 5, ya que hasta ese punto el mismo se comporta como un echo server.

## Parte 2: Comunicación

### Ejercicio 5

Hasta acá el protocolo es solamente la cantidad de bytes a enviar y los datos, que es la apuesta individual, cuyos datos se separan por el caracter `:`. Además del `OK` como ack del servidor.

**Nota**: quedó el loop del servidor original para probar el envío de varios mensajes del mismo cliente (aunque siempre es el mismo para una agencia en particular, porque los datos se toman de variables de entorno).

### Ejercicio 6

Se envían las apuestas por batches. Para esto se lee el archivo csv secuencialmente y cada X cantidad de lineas (configurables) se envían al server.

Agregados al protocolo:

- las apuestas van separadas por `;`
- el servidor responde `OK` cuando todo el batch fue procesado

### Ejercicio 7

Lista de booleanos `agencies_done` que indica qué agencias terminaron de mandar todas las apuestas.

Agregados al protocolo:

- se agrega el concepto de `ACCION`, para saber de que operación se está tratando (`BET`, `FINISH`, `WINNER`)
- respuesta `WINNER` para resultado de sorteo y respuesta `WINNERWAIT` para comunicar que todavia no está listo el sorteo (y el tiempo a esperar para volver a consultar)

### Protocolo Final

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

**Nota**: el protocolo se podría optimizar utilizando solo 2 bytes fijos para la acción. Además, una simplificación posible es que el servidor responda el mensaje `OK` armado de la misma manera que el resto de mensajes (4 bytes indicando el tamaño del mensaje + el mensaje en sí).

#### Ejemplo

Un mensaje de apuestas con 2 apuestas:

```
0000004CBET::1:Juan:Perez:12345678:1980-01-01:1;1:Fulano:Gomez:12345678:1980-01-01:2
```

## Parte 3: Concurrencia

### Ejercicio 8

Para la concurrencia se utilizó el módulo _multithreading_ en el servidor. Por cada cliente que se conecta se crea un nuevo thread que se encarga de procesar los mensajes del cliente. El thread principal se encarga de aceptar las conexiones entrantes y crear los threads de los clientes. Para que el envío de las apuestas no requiera la creación de multiples sockets (como estaba implementado hasta el ejercicio 7), se crea un socket para el envío de todos los batches de apuestas, hasta el momento en que se recibe el mensaje `FINISH`. La consulta de ganadores sí crea un nuevo socket para cada consulta.

Luego, cuando se requiere acceder al archivo de apuestas, se utiliza un _lock_ para evitar que dos threads accedan al mismo tiempo al archivo (tanto al momento de lectura como de escritura).

La sincronización para obtener los resultados del sorteo se realiza haciendo _polling_ desde el cliente. Si el servidor todavía no recibió las apuestas de todas las agencias, responde con un mensaje indicando el tiempo de espera que debe esperar el cliente para volver a consultar. Si el servidor tiene los resultados, responde con un mensaje con los ganadores. El tiempo de espera es fijo y configurable en el servidor.

Teniendo en cuenta que hay tiempo de I/O se pueden utilizar threads sin que el Global Interpreter Lock (GIL) sea un gran problema. De todas formas, se podría utilizar el módulo _multiprocessing_ para evitar este inconveniente creando procesos en lugar de threads.

**Nota**: en el ejercicio 7 se utilizó una lista de booleanos para tener un registro de las agencias que ya enviaron todas sus apuestas. Desde el ejercicio 8 se utiliza un contador (un int) para este mismo propósito.
