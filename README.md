# Eventsourcing CQRS example in GO

Example event sourcing and CQRS implementation in [golang](https://golang.org)

Project is a re-implementation of [https://github.com/pinballjs/event-sourcing-cqrs] in GO. Initial inspiration
originates from Chris Richardson's [presentation](http://www.infoq.com/presentations/event-microservice-scala-spring-boot)

## How to run

    git clone https://github.com/botchniaque/eventsourcing-cqrs-go.git
    cd eventsourcing-cqrs-go/
    export GOPATH=`pwd`
    go get -t -v example eventsourcing
    go build example
    ./example


to run tests

    go test -v eventsourcing

## Design
Example implements 2 aggregates - Account and Money Transfer

implementation includes number of components:
- `Event Store` - persists events, and allows searching by aggregate ID.
- `Event Handler` - gets notified about persisted events, and produces commands as a result of some events
(eg. `DebitAccountCommandBecauseOfMoneyTransfer` when receives `MoneyTransferCreatedEvent`)
- `Command Handler` - each aggregate has own handler, which produce state changing events,
which then are persisted in the `Event Store`
- Customer facing `Services` - expose usable functions to clients and produce commands
as a product of interaction

Sketched diagram:

     +---------------+                        +--------------+
     | Event Handler |                        |  Event Store |
     |               |                        |   * Update   |
     |               +<---Persisted Events----+   * Find     |
     |               |        Channel         |              |
     +-------+-------+                        +------+-------+
             |                                       ^
          commands                                  events
             |                                       |
             |                             +---------+---------+
             +--------->                   | Command Handler   |
                        \                  |  - Account        |
                         +-----Command---> |  - Money Transfer |
                        /      Channel(s)  |                   |
                +------>                   +-------------------+
                |
             comamnds
                |
     +----------+--------+
     | Service           |
     |  - Account        |
     |    * Open         |
     |    * Credit       |
     |    * Debit        +<---Customer requests---+
     |  - Money Transfer |
     |    * Create new   |
     |                   |
     +-------------------+

### What's missing
#### Message bus
 After persisting, `Event Store` should publish events to a proper `Message Bus` where components
 can subscribe to those events. In the example, for simplification,
 the `Persisted Events Channel` works as a single-subscriber `Message Bus` - since channels are
 first-class citizens in GO, this simplified a lot, but it's not enough for real-life implementations.

 To implement a simple bus allowing multiple subscribes, a 'linked channel' idea could be implemented
 as described [here](https://rogpeppe.wordpress.com/2009/12/01/concurrent-idioms-1-broadcasting-values-in-go-with-linked-channels/)

#### Event (de)serialization
Example uses simple map/slice implementation of `Event store`. For real-life implementations events would
need to be serialized (probably to JSON) for persistence. For simplification that part was omitted.
