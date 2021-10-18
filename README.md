# SDCC PROGETTO
## Storage persistente di tipo Chiave-Valore

<img src="https://miro.medium.com/max/2000/1*qkPPwhQrMtoL-e18k0TXRA.png" width=500px>

Lo scopo del progetto e' realizzare, nel linguaggio di programmazione Go, unsistema di storage distribuitoper edge computing.
L'applicazione e' stata pensata come un servizio di storage persistente, replicato e di tipo chiave-valore, fornito da nodi edge che comunicano tra loro. I valori sono il risultato delle temperature misurate da un insieme di sensori.
I client possono interagire con i nodi edge attraverso 4 chiamate RPC:
- Put (key,value)
- Get (key)
- Append (key,value)
- Delete (key)

## Requisiti
- <img src="https://blog.seeweb.it/wp-content/uploads/2015/06/homepage-docker-logo.png" width=50px> Docker
- <img src="https://www.docker.com/blog/wp-content/uploads/2020/02/Compose.png" width=50px> Docker-compose
- <img src="https://i.pinimg.com/originals/28/ec/74/28ec7440a57536eebad2931517aa1cce.png" width=50px> Terraform
  

- <img src="https://upload.wikimedia.org/wikipedia/commons/thumb/9/93/Amazon_Web_Services_Logo.svg/1200px-Amazon_Web_Services_Logo.svg.png" width=50px> Credenziali AWS valide


- <img src="https://www.geekandjob.com/uploads/wiki/591f10c4e56bf30f45a4ad0b8956223c04802eac.png" width=50px>Go v1.17.1

## Struttura
> ### Server
> Rappresenta il nodo edge del sistema, responsabile di ricevere chiamate
> RPC dai nodi client, fare storage K-V dei valori associate alle chiavi,
> fare chiamate a DynamoDB e comunicare (tramite chiamate RPC) con gli altri
> nodi edge per mantenere consistenti le repliche.
> L'ultimo punto e' reso possibile tramite l'implementazione dell'algoritmo
> [RAFT](https://raft.github.io/) tra i nodi edge.

> ### Client
> Il nodo che rappresenta un cluster di sensori di temperatura,
> responsabile di generare valori randomici e effettuare le
> chiamate RPC ai nodi edge.

> ### Master
> Il nodo Master e' responsabile della fase di registrazione
> nel sistema tenendo traccia dei nodi che la compongono

## How to Use
Aggiornare il file ~/.aws/credentials con delle credenziali valide da [AWS](https://aws.amazon.com/it/).
_(Attenzione: l'applicazione sara' lanciata con i privilegi massimi, di conseguenza si dovra' modificare il file /root/.aws/credentials in dispositivi Linux.)_

Per interagire con l'applicazione, nella directory _scripts_ eseguire:
Build e run del programma senza un backup, in cold start con:
```
./cold-start.sh
```
Lo script andra' anche a creare l'ambiente con dynamoDB e le lambda attraverso terraform.
Terminare il programma con CTRL+C

Build e run del programma con backup:
```
./start.sh
```
Terminare l'applicazione e distruggere l'ambiente con:
```
./shutdown.sh
```

Durante l'esecuzione del programma si puo' eseguire:
```
./cli.sh
```
per generare una shell di dialogo con il cluster di nodi edge

Per eliminare la tabella di DynamoDB eseguire:
```
./delete-table-DynamoDB.sh
```


#### Authors
- Matteo Chiacchia
- Valerio Cristofori