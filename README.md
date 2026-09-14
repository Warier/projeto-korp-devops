# Projeto Korp

API HTTP em Go/Gin que retorna o nome do projeto e o horário atual em UTC.

```http
GET /projeto-korp
```

```json
{"nome":"Projeto Korp","horario":"2026-09-14T12:00:00Z"}
```

O horário é calculado a cada requisição. A aplicação escuta na porta 8080, sem publicação direta no host.

## Arquitetura

- **NGINX:** entrada HTTP na porta 80 e proxy para a API.
- **API Gin:** endpoint HTTP e métricas Prometheus em `/metrics`, bloqueadas na entrada pública do NGINX.
- **NGINX log exporter:** transforma logs de respostas da API em métricas por status e duração.
- **Prometheus:** coleta e armazena métricas por até 7 dias, com retenção por tamanho de 2 GB para os blocos. Não publica porta no host.
- **Grafana:** datasource e dashboard provisionados por arquivos; interface em `127.0.0.1:3000`.

O Compose cria duas redes bridge: `korp` para NGINX/API e `observabilidade` para as coletas. API e NGINX participam das duas. O exporter compartilha a rede do container NGINX (`network_mode: service:nginx`): recebe syslog UDP em `127.0.0.1:5531` e expõe métricas em `nginx:4040`, sem publicar essas portas no host.

As imagens da API usam build multi-stage e runtime `scratch` sem root. Os serviços têm filesystem base somente para leitura; os dados do Prometheus e Grafana ficam em volumes persistentes.


## Provisionamento com Ansible

Exemplo de destino: VM com Ubuntu Server 24.04 LTS, Python 3, SSH, acesso à internet e usuário com sudo.

Na máquina que executará o Ansible, crie `ansible/inventories/server.yml` a partir de `server.yml.example` e ajuste o IP e o usuário SSH.

Execute na raiz do projeto:

```bash
ansible-playbook -i ansible/inventories/server.yml ansible/playbook.yml -K
```

Esse único comando instala o Docker, transfere os arquivos, faz o build, cria as redes e inicia API, NGINX, Prometheus, exporter e Grafana com o dashboard provisionado. Ao final, valida a API e exibe o JSON no console.

### Execução local

Para Linux ou WSL com Docker já instalado e acessível pelo usuário:

```bash
ansible-playbook -i ansible/inventories/local.yml ansible/playbook.yml
```


### Dependência do Ansible

O controller precisa ter Ansible Core 2.17+ e a coleção `community.docker`. Para instalar a versão usada pelo projeto, execute uma vez:

```bash
ansible-galaxy collection install -r ansible/requirements.yml
```


## Execução rápida

Pré-requisitos: Docker Engine ativo e Docker Compose v2.18 ou superior. Na raiz do checkout:

```bash
docker compose up --build --wait --wait-timeout 300
```

Não é necessário criar `.env`. Teste a API:

```bash
curl http://localhost:80/projeto-korp
```

## Grafana e senha

- Dashboard local: [Projeto Korp](http://localhost:3000/d/projeto-korp).
- Usuário: `admin`.
- Sem senha personalizada (`.env` ausente ou variável vazia): senha inicial **`admin`**.
- Com senha personalizada: defina `GRAFANA_ADMIN_PASSWORD` no `.env` antes da primeira inicialização. `.env.example` mostra o formato.


## Métricas e dashboard

O dashboard apresenta volume de requisições, respostas por status, taxa de falhas 5xx no NGINX, latências p50/p95/p99 no proxy e no Gin e saúde das duas coletas.

- O NGINX registra apenas `/projeto-korp` para o exporter e exclui o healthcheck originado em seu loopback. Não envia query strings, IPs ou corpos nesse fluxo de métricas.
- `502` e `504` aparecem nas métricas do proxy mesmo quando a requisição não chega ao Gin. Respostas `4xx` aparecem por status, mas não são classificadas como indisponibilidade do servidor.
- As métricas do Gin incluem o healthcheck do NGINX e excluem `/metrics`; por isso os volumes não precisam coincidir.
- `up` indica sucesso da coleta, não disponibilidade ponta a ponta. Ausência de tráfego não comprova disponibilidade e pode produzir painéis sem dados.
- Syslog UDP local evita um arquivo de logs crescente, mas pode perder eventos durante reinícios ou sobrecarga. As métricas do exporter reiniciam com seu processo; não são um registro de auditoria.

As coletas ocorrem a cada 15 segundos. Faça algumas requisições pelo host e aguarde pelo menos duas coletas para visualizar taxas. Os percentis são estimados pelos buckets, e o volume calculado com `increase()` é uma estimativa no período selecionado.

## Qualidade e testes

Para desenvolvimento, use Go 1.27.1 (a versão do builder) e GNU Make:

```bash
make check
```

Esse comando verifica formatação (`gofmt`), executa `go vet`, Staticcheck em versão fixada e testes com cobertura. Qualquer falha retorna código diferente de zero, permitindo utilizar o mesmo comando em uma CI futura. A primeira execução baixa o Staticcheck para o cache Go; ele não entra na imagem da API.

```bash
make fmt                 # Corrige a formatação
make test                # Executa apenas os testes
go tool cover -func=coverage.out
```

Os testes cobrem o contrato JSON/UTC e a instrumentação de respostas, erros recuperados e exclusão das coletas. `coverage.out` não é versionado. Execute essas verificações na máquina de desenvolvimento; elas não fazem parte do provisionamento da VM.

