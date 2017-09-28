# Go @ OptioPay

<h5 class="fragment">Ole Bulbuk: Senior Software Engineer</h5>
<h6 class="fragment">Back end guy since the ninties; Gopher for 2 years</h6>


---?image=assets/OptioPay-Business-Model.png

---

## OptioPays Special Requirements

- We manage money and private data |
- Lack of auditing could kill us |
- Lack of data privacy could kill us |
- Security is very important |
- A dynamic language like JavaScript doesn't feel right |

---

## Event Sourcing With Apache Kafka

- All data is stored as events |
- Events can't be modified but current state can change due to new events |
- Enables auditing |
- Decouples services |

---

## OptioPay Mircro Service Framework: Micro

- Abstracts away reading of Kafka topics |
- Keeps track of client position within Kafka topic |
- Helps writing to Kafka topics |
- Composes HTTP handlers |
- Service implementations have to implement interfaces | 
- More than 40 services built on top of Micro |
- Most of them store current state in a PostgreSQL DB |

---

## Example Go Service that Uses Micro

Show code from my onboarding woodblock service.

---

## Further Usages Of Go

- Some command line tools |
- Most create or read CSV files |
