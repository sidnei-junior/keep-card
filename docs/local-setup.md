# Local setup

## Prerequisites

- [Docker Desktop](https://www.docker.com/products/docker-desktop/) using the **WSL 2** backend (Windows).

### Recommended: limit WSL resources (Windows)

Docker Desktop runs inside WSL 2, which can use a lot of RAM by default. Create `%UserProfile%\.wslconfig`:

```ini
[wsl2]
memory=4GB
processors=2
```

Then apply it (this stops Docker Desktop; start it again afterwards):

```bash
wsl --shutdown
```

## RabbitMQ

The local message broker is defined in [`docker-compose.yml`](../docker-compose.yml).

| What | Value |
| --- | --- |
| AMQP (services) | `amqp://guest:guest@localhost:5672` |
| Management UI | http://localhost:15672 (user `guest`, password `guest`) |

### Commands

| Action | Command |
| --- | --- |
| Start | `docker compose up -d` |
| Check status | `docker compose ps` (wait for `healthy`) |
| Logs | `docker compose logs -f rabbitmq` |
| Stop (keeps data) | `docker compose down` |
| Reset (deletes all exchanges, queues and messages) | `docker compose down -v` |

## Creating the topology through the UI

Follow the [base architecture](./architecture/README.md#messaging). In the Management UI:

1. **Exchange** — *Exchanges* tab → *Add a new exchange*:
   - Name: `keepcard.events`, Type: `topic`, Durability: `Durable`.
2. **Queues** — *Queues and Streams* tab → *Add a new queue* (Type `Quorum` or `Classic`, durable), one for each:
   - `card.events`, `inventory.purchases`, `notification.all`
3. **Bindings** — open the `keepcard.events` exchange → *Bindings* → *Add binding from this exchange*:

   | To queue | Routing key |
   | --- | --- |
   | `card.events` | `user.created` |
   | `card.events` | `product.*` |
   | `inventory.purchases` | `purchase.registered` |
   | `notification.all` | `user.created` |
   | `notification.all` | `purchase.registered` |
   | `notification.all` | `reward.redeemed` |
   | `notification.all` | `stock.low` |

## Sending a test message

1. Open the `keepcard.events` exchange → *Publish message*.
2. Routing key: `purchase.registered`. Payload, for example:

   ```json
   { "purchaseId": "p-1", "cardId": "c-1", "productId": "prod-1", "quantity": 1 }
   ```

3. Click *Publish message*. It should land in **`inventory.purchases`** and **`notification.all`**, but not in `card.events`.
4. Open a queue → *Get messages* to read it. Use *Ack mode: Reject requeue true* to peek without removing it, or *Automatic ack* to consume it.

Tip: publish with routing key `product.created` and check that only `card.events` receives it — that is the `*` wildcard in action.
