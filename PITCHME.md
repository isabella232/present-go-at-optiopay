# Go @ OptioPay

<h5 class="fragment">Ole Bulbuk: Senior Software Engineer</h5>
<h6 class="fragment">Back end guy since the ninties; Gopher for 2 years</h6>


---?image=assets/OptioPay-Business-Model.png

---

## Special Requirements

- We manage money and private data |
- Lack of auditing could kill us |
- Lack of data privacy could kill us |
- Security is VERY important |
- A dynamic language like JavaScript doesn't feel right |

---

## Event Sourcing: Apache Kafka

- All data is stored as events |
- Events are immutable: NEW events change state |
- This enables auditing |
- Persistent messaging decouples services |

---

## Micro Service Framework

- Abstracts away reading of Kafka topics |
- Keeps track of client position within Kafka topic |
- Helps writing to Kafka topics |
- Composes HTTP handlers |
- Each service implements interfaces for Micro | 
- More than 40 services built on top of Micro |
- Most of them store current state in a PostgreSQL DB |

---?code=go/main.go

## Woodblock service: main.go

@[1-8](the main file has to be in package main and needs some imports)
@[10-11](the version is set by a build script)
@[13-14](in main we first create an instance of our internal service)
@[15-24](then we create a configuration for it)
@[26](now we finally create the official micro service)
@[27-30](and run it hopefully without any error)

---?code=go/service.go

## Woodblock service: service.go

@[1-10](the service is usually in its own package and needs some imports, too)
@[18-24](NewService returnes a service ready to use)
@[26-36](ProcessEvent processes a Kafka event in order to keep the state of the service up to date)
@[38-44](remember new vouchers)
@[46-54](count sold vouchers)
@[56-65](forget vouchers that aren't valid anymore)
@[67-76](return the count as JSON)

---

## Further Uses Of Go

- Some command line tools |
- Most create or read CSV files |
