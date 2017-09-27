# Go @ OptioPay

---

## Show OptioPay triangle: Recipient - Issuer - Advertiser

---

## OptioPays Special Requirements

- We manage money and private data |
- Lack of auditing could kill us |
- Lack of data privacy could kill us |
- Security is *very* important |
- A dynamic language like JavaScript doesn't feel right |

---

## Event Sourcing With Apache Kafka

- Enables auditing |
- Decouples services |

---

## OptioPay Mircro Service Framework: Micro

- Abstracts away reading of Kafka topics |
- Keeps track of client position within Kafka topic |
- Helps writing to Kafka topics |

---

## Go Services Use OptioPay Micro

- XX services built on top of Micro
- Most of them store current state in a PostgreSQL DB

---

## Further Usages Of Go

- Some command line tools |
- Most create or read CSV files |
