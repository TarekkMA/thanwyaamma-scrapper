# ThanwyAamma Scrapper

Go Scrapper for scrapping ThanwyAamma results.

### Installation

```bash
go get -u github.com/TarekkMA/thanwyaamma-scrapper
```
Excutable should be available at ``$GOBIN``

### Usage

```bash
thanwyaamma-scrapper -s 222131 -n 1000 -c 100
```

#### Flags
| flag| Description                           | Required |
| :-: |:-------------------------------------| :-------:|
| s   | Starting seat number              |  ``✔``   |
| n   | Number of seats numbers to scrap      |   ``✔``    |
| c   | Number of concurrent workers          | ``✔`` |
