# Desafio 04 - Rate Limiter

Para iniciar o redis use o comando docker abaixo:

```
docker-compose up -d
```

Para rodar a aplicação, entre na pasta:

```
cd cmd/server
```

E rode o comando abaixo:

```
go run main.go
```

**Observações necessárias:** Existe um token para fins de teste no arquivo `.env`, que está na pasta `/cmd/server`.
As configurações de limite e tempo também estão neste mesmo arquivo e podem ser modificadas.

Para rodar os testes dentro da pasta `/cmd/server` use o seguinte comando:

```
go test -v
```

Também temos dois arquivos `.http` na pasta `/test` para fazer requisições individuais e testar a efetividade.
