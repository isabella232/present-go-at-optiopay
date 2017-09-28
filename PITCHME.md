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

## Mircro Service Framework

- Abstracts away reading of Kafka topics |
- Keeps track of client position within Kafka topic |
- Helps writing to Kafka topics |
- Composes HTTP handlers |
- Service implementations have to implement interfaces | 
- More than 40 services built on top of Micro |
- Most of them store current state in a PostgreSQL DB |

---?code=go/main.go

@[1-7](the main file has to be in package main and needs some imports)
@[9-10](the version is set by a build script)
@[12-13](in main we first create an instance of our internal service)
@[14-23](then we create a configuration for it)
@[25](now we finally create the official micro service)
@[26-29](and run it hopefully without any error)

---

## Further Usages Of Go

- Some command line tools |
- Most create or read CSV files |
