# sdcc-project

## For Testing

Server
```
go test -run TestServer

```

Client
```
go test -run TestClient

```


### TO DO
- [x] Se crusha il nodo leader il producer deve rimandare in broadcast il comando e aspettare la risposta del nuovo leader per poi comunicare con lui
- [x] Implementare e testare un modello di consistenza di tipo sequenziale
- [ ] Memorizzare nel cloud ( dynamoDB ) solo i valori di grandi dimensioni e/o scarsamente acceduti
- [x] Tolleranza crush nodi edge
- [ ] Persistere mappa chiave-valore nei nodi edge anche dopo la distruzione dei container
- [ ] Stabilire giusta sequenza di startup dei container
- [ ] Gestire errori di DynamoDB ( creazione tabella ecc. )
- [ ] Chiamare funzioni PUT GET APPEND DELETE da edge per dynamoDB ( Quando? Inserire un TTL invece del Timestamp in *Data)
- [ ] Testing di: 85 Get, 15 Put; 40 Put, 20 Append, 40 Get
- [ ] Usare Docker-swarm invece di docker-compose
- [ ] Implementare semantica errori di tipo at-least-once
