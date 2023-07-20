# PortScanner

Instalacija:
Klonirati repozitorij u lokalni direktorij.
>git clone https://github.com/kpavlovic123/PortScanner.git


Ako je go instaliran ("sudo apt install golang-go" na linuxu) onda je moguće pokrenuti na 3 načina:

a) Go run <naziv.go> <argumenti>
    Sam izvodi compile i pokreće program.
    Naziv.go jest datoteka definirana s package main (glavna datoteka). U ovome slučaju portScanner.go.
    > go run portScanner.go <argumenti>

b) Go build
    Izvodi compile i stvara izvršnu datoteku (portScanner).
    > go build 
    > portScanner <argumenti>

c) Go install
    Izvodi compile i sprema izvršnu datoteku u GOPATH/bin. Time je moguće pokrenuti datoteku kao naredbu (i u drugome
    direktoriju),ukoliko je GOPATH postavljen u environment varijablu ($PATH). npr.
    > go install
    > cd ..
    > portScanner <argumenti>
    
    Putanju do GOPATH je moguće vidjeti s naredbom 
    >go env GOPATH



Format argumenata(slično nmap-u):

    portScanner [-p <format...>] [-p-] [-sU] [-sT] [-o] !destinationAddress

destinationAddress : Mora biti na kraju naredbe, označuje destinaciju.

-p <format> : Označuje ciljne portove. Format mora biti u obliku "broj,broj...". Nije dopušten razmak. Svi
    elementi moraju biti odvojeni zarezom. Također, moguće je i označiti raspon (npr. 1-1000, oboje je uključeno).
    Po defaultu (dakle bez poziva -p ili -p-) se gledaju portovi 1-1023. Primjer poziva:
    > portScanner -p 1,100,101-1000 localhost

-p- : Označuje sve portove, ekvivalentno s "...-p 1-65535...".

-sU / -sT: Označuje tip scana. Po defaultu se obavljaju i tcp i udp scan. Da bi se izveo samo udp, onda -sU, da bi se 
    izveo samo tcp, onda -sT. Moguće je pozvati s oba flaga, no onda nema efekta.

